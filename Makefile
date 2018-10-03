# Image URL to use all building/pushing image targets
APP_VERSION ?= $(shell git describe --abbrev=5 --dirty --tags --always)
REGISTRY := quay.io/presslabs
IMAGE_NAME := dashboard
BUILD_TAG := build
IMAGE_TAGS := $(APP_VERSION)
KUBEBUILDER_VERSION ?= 1.0.4
PROTOC_VERSION ?= 3.6.1
BINDIR ?= $(PWD)/bin
CHARTDIR ?= $(PWD)/chart/dashboard

GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH ?= amd64

PATH := $(BINDIR):$(PATH)
SHELL := env 'PATH=$(PATH)' /bin/sh

all: test dashboard

# Run tests
test: generate manifests
	KUBEBUILDER_ASSETS=$(BINDIR) ginkgo \
		--randomizeAllSpecs --randomizeSuites --failOnPending \
		--cover --coverprofile cover.out --trace --race -v \
		./pkg/... ./cmd/...

# Build dashboard binary
build: generate fmt vet
	go build -o bin/dashboard github.com/presslabs/dashboard/cmd/dashboard

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/dashboard/main.go

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all
	# CRDs
	awk 'FNR==1 && NR!=1 {print "---"}{print}' config/crds/*.yaml > $(CHARTDIR)/templates/_crds.yaml
	yq m -d'*' -i $(CHARTDIR)/templates/_crds.yaml hack/chart-metadata.yaml
	yq w -d'*' -i $(CHARTDIR)/templates/_crds.yaml 'metadata.annotations[helm.sh/hook]' crd-install
	yq d -d'*' -i $(CHARTDIR)/templates/_crds.yaml metadata.creationTimestamp
	yq d -d'*' -i $(CHARTDIR)/templates/_crds.yaml status metadata.creationTimestamp
	echo '{{- if .Values.crd.install }}' > $(CHARTDIR)/templates/crds.yaml
	cat $(CHARTDIR)/templates/_crds.yaml >> $(CHARTDIR)/templates/crds.yaml
	echo '{{- end }}' >> $(CHARTDIR)/templates/crds.yaml
	rm $(CHARTDIR)/templates/_crds.yaml
	# RBAC
	cp config/rbac/rbac_role.yaml $(CHARTDIR)/templates/_rbac.yaml
	yq m -d'*' -i $(CHARTDIR)/templates/_rbac.yaml hack/chart-metadata.yaml
	yq d -d'*' -i $(CHARTDIR)/templates/_rbac.yaml metadata.creationTimestamp
	yq w -d'*' -i $(CHARTDIR)/templates/_rbac.yaml metadata.name '{{ template "dashboard.fullname" . }}-controller'
	echo '{{- if .Values.rbac.create }}' > $(CHARTDIR)/templates/controller-clusterrole.yaml
	cat $(CHARTDIR)/templates/_rbac.yaml >> $(CHARTDIR)/templates/controller-clusterrole.yaml
	echo '{{- end }}' >> $(CHARTDIR)/templates/controller-clusterrole.yaml
	rm $(CHARTDIR)/templates/_rbac.yaml

.PHONY: chart
chart:
	yq w -i $(CHARTDIR)/Chart.yaml version "$(APP_VERSION)"
	yq w -i $(CHARTDIR)/Chart.yaml appVersion "$(APP_VERSION)"
	mv $(CHARTDIR)/values.yaml $(CHARTDIR)/_values.yaml
	sed 's#$(REGISTRY)/$(IMAGE_NAME):latest#$(REGISTRY)/$(IMAGE_NAME):$(APP_VERSION)#g' $(CHARTDIR)/_values.yaml > $(CHARTDIR)/values.yaml
	rm $(CHARTDIR)/_values.yaml

.PHONY: bundle
bundle:
	$(BINDIR)/packr -i $(PWD)/pkg -v

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate:
	go generate ./pkg/... ./cmd/...
	protoc -I ./proto \
	--go_out=plugins=grpc:$(GOPATH)/src \
	presslabs/dashboard/core/v1/project.proto

# Build the docker image
.PHONY: images
images: bundle
	docker build . -t $(REGISTRY)/$(IMAGE_NAME):$(BUILD_TAG)
	set -e; \
		for tag in $(IMAGE_TAGS); do \
			docker tag $(REGISTRY)/$(IMAGE_NAME):$(BUILD_TAG) $(REGISTRY)/$(IMAGE_NAME):$${tag}; \
	done
	$(BINDIR)/packr clean -v


# Push the docker image
.PHONY: publish
publish: images
	set -e; \
		for tag in $(IMAGE_TAGS); do \
		docker push $(REGISTRY)/$(IMAGE_NAME):$${tag}; \
	done

lint:
	$(BINDIR)/golangci-lint run ./pkg/... ./cmd/...

dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)
	GOBIN=$(BINDIR) go install ./vendor/github.com/onsi/ginkgo/ginkgo
	GOBIN=$(BINDIR) go install ./vendor/github.com/golang/protobuf/protoc-gen-go
	GOBIN=$(BINDIR) go install ./vendor/github.com/gobuffalo/packr/packr
	which unzip || (apt-get update && apt-get install --no-install-recommends -y unzip)
ifeq ($(GOOS),darwin)
	curl -sfL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-osx-x86_64.zip -o protoc.zip
else
	curl -sfL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(GOOS)-x86_64.zip -o protoc.zip
endif
	unzip -u protoc.zip bin/protoc
	rm protoc.zip
	curl -sfL https://github.com/mikefarah/yq/releases/download/2.1.1/yq_$(GOOS)_$(GOARCH) -o $(BINDIR)/yq
	chmod +x $(BINDIR)/yq
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $(BINDIR) v1.10.2
	curl -sfL https://github.com/kubernetes-sigs/kubebuilder/releases/download/v$(KUBEBUILDER_VERSION)/kubebuilder_$(KUBEBUILDER_VERSION)_$(GOOS)_$(GOARCH).tar.gz | \
		tar -zx -C $(BINDIR) --strip-components=2
