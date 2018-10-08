# Image URL to use all building/pushing image targets
APP_VERSION ?= $(shell git describe --abbrev=5 --dirty --tags --always)
IMG ?= quay.io/presslabs/dashboard:$(APP_VERSION)
KUBEBUILDER_VERSION ?= 1.0.4
PROTOC_VERSION ?= 3.6.1
BINDIR ?= $(PWD)/bin

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

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
install: manifests chart
	helm upgrade --install --namespace=presslabs-sys dashboard chart/dashboard --set image=$(IMG)

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all

.PHONY: chart
chart:
	rm -rf chart/dashboard
	cp -r chart/dashboard-src chart/dashboard
	yq w -i chart/dashboard/Chart.yaml version "$(APP_VERSION)"
	yq w -i chart/dashboard/Chart.yaml appVersion "$(APP_VERSION)"
	yq w -i chart/dashboard/values.yaml image "$(IMG)"
	awk 'FNR==1 && NR!=1 {print "---"}{print}' config/crds/*.yaml > chart/dashboard/templates/crds.yaml
	yq m -d'*' -i chart/dashboard/templates/crds.yaml hack/chart-metadata.yaml
	yq w -d'*' -i chart/dashboard/templates/crds.yaml 'metadata.annotations[helm.sh/hook]' crd-install
	yq d -d'*' -i chart/dashboard/templates/crds.yaml metadata.creationTimestamp
	yq d -d'*' -i chart/dashboard/templates/crds.yaml status metadata.creationTimestamp
	cp config/rbac/rbac_role.yaml chart/dashboard/templates/rbac.yaml
	yq m -d'*' -i chart/dashboard/templates/rbac.yaml hack/chart-metadata.yaml
	yq d -d'*' -i chart/dashboard/templates/rbac.yaml metadata.creationTimestamp
	yq w -d'*' -i chart/dashboard/templates/rbac.yaml metadata.name '{{ template "dashboard.fullname" . }}-controller'
	echo '{{- if .Values.rbac.create }}' > chart/dashboard/templates/controller-clusterrole.yaml
	cat chart/dashboard/templates/rbac.yaml >> chart/dashboard/templates/controller-clusterrole.yaml
	echo '{{- end }}' >> chart/dashboard/templates/controller-clusterrole.yaml
	rm chart/dashboard/templates/rbac.yaml

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
images: test
	docker build . -t ${IMG}

# Push the docker image
publish:
	docker push ${IMG}

lint:
	$(BINDIR)/golangci-lint run ./pkg/... ./cmd/...

dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)
	GOBIN=$(BINDIR) go install ./vendor/github.com/onsi/ginkgo/ginkgo
	GOBIN=$(BINDIR) go install ./vendor/github.com/golang/protobuf/protoc-gen-go
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
