# Image URL to use all building/pushing image targets
APP_VERSION      ?= $(shell git describe --abbrev=5 --dirty --tags --always)
REGISTRY         := quay.io/presslabs
IMAGE_NAME       := dashboard
BUILD_TAG        := build
IMAGE_TAGS       := $(APP_VERSION)
BINDIR           ?= $(CURDIR)/bin
CHARTDIR         ?= $(CURDIR)/chart/dashboard

GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH ?= amd64

APP_BIN_PATH := $(CURDIR)/app/node_modules/.bin
PATH := $(APP_BIN_PATH):$(BINDIR):$(PATH)
SHELL := env 'PATH=$(PATH)' /bin/sh

# Run tests
test: generate manifests
	ginkgo \
		--randomizeAllSpecs --randomizeSuites --failOnPending \
		--cover --coverprofile cover.out --trace --race \
		./pkg/... ./cmd/...

# Run apiserver tests
apiserver-test:
	ginkgo \
	--randomizeAllSpecs --randomizeSuites --failOnPending \
	--cover --coverprofile cover.out --trace --race -v \
	./pkg/apiserver/...

# Build dashboard binary
build: generate fmt vet
	go build -o bin/dashboard github.com/presslabs/dashboard/cmd/dashboard

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/dashboard/main.go

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go webhook

	# CRDs
	awk 'FNR==1 && NR!=1 {print "---"}{print}' config/crds/*.yaml > $(CHARTDIR)/templates/_crds.yaml
	yq m -d'*' -i $(CHARTDIR)/templates/_crds.yaml hack/chart-metadata.yaml
	yq w -d'*' -i $(CHARTDIR)/templates/_crds.yaml 'metadata.annotations[helm.sh/hook]' crd-install
	yq d -d'*' -i $(CHARTDIR)/templates/_crds.yaml metadata.creationTimestamp
	yq d -d'*' -i $(CHARTDIR)/templates/_crds.yaml status metadata.creationTimestamp
	cat hack/generated-by-makefile.yaml > $(CHARTDIR)/templates/crds.yaml
	echo '{{- if .Values.crd.install }}' >> $(CHARTDIR)/templates/crds.yaml
	cat $(CHARTDIR)/templates/_crds.yaml >> $(CHARTDIR)/templates/crds.yaml
	echo '{{- end }}' >> $(CHARTDIR)/templates/crds.yaml
	rm $(CHARTDIR)/templates/_crds.yaml
	# RBAC
	cp config/rbac/rbac_role.yaml $(CHARTDIR)/templates/_rbac.yaml
	yq m -d'*' -i $(CHARTDIR)/templates/_rbac.yaml hack/chart-metadata.yaml
	yq d -d'*' -i $(CHARTDIR)/templates/_rbac.yaml metadata.creationTimestamp
	yq w -d'*' -i $(CHARTDIR)/templates/_rbac.yaml metadata.name '{{ template "dashboard.fullname" . }}-controller'
	cat hack/generated-by-makefile.yaml > $(CHARTDIR)/templates/controller-clusterrole.yaml
	echo '{{- if .Values.rbac.create }}' >> $(CHARTDIR)/templates/controller-clusterrole.yaml
	cat $(CHARTDIR)/templates/_rbac.yaml >> $(CHARTDIR)/templates/controller-clusterrole.yaml
	echo '{{- end }}' >> $(CHARTDIR)/templates/controller-clusterrole.yaml
	rm $(CHARTDIR)/templates/_rbac.yaml
	# Webhooks
	cp config/webhook/webhook.yaml chart/dashboard/templates/webhook.yaml
	yq m -d'*' -i $(CHARTDIR)/templates/webhook.yaml hack/chart-metadata.yaml
	yq d -d'*' -i $(CHARTDIR)/templates/webhook.yaml metadata.creationTimestamp
	yq d -d'*' -i $(CHARTDIR)/templates/webhook.yaml metadata.namespace
	#   secret
	yq w -d0 -i $(CHARTDIR)/templates/webhook.yaml 'data[ca-cert.pem]' '{{ $$cert.Cert | b64enc }}'
	yq w -d0 -i $(CHARTDIR)/templates/webhook.yaml 'data[key.pem]' '{{ $$cert.Key | b64enc }}'
	yq w -d0 -i $(CHARTDIR)/templates/webhook.yaml 'data[cert.pem]' '{{ $$cert.Cert | b64enc }}'
	yq d -d0 -i $(CHARTDIR)/templates/webhook.yaml 'data[ca-key.pem]'
	yq w -d0 -i $(CHARTDIR)/templates/webhook.yaml metadata.name '{{ include "dashboard.webhook.secretName" . }}'
	#   service
	yq w -d1 -i $(CHARTDIR)/templates/webhook.yaml metadata.name '{{ include "dashboard.fullname" . }}-webhook'
	yq d -d1 -i $(CHARTDIR)/templates/webhook.yaml metadata.namespace '{{ .Release.Namespace }}'
	yq w -d1 -i $(CHARTDIR)/templates/webhook.yaml spec.type '{{ .Values.webhook.service.type }}'
	yq w -d1 -i $(CHARTDIR)/templates/webhook.yaml spec.type '{{ .Values.webhook.service.type }}'
	yq w -d1 -i $(CHARTDIR)/templates/webhook.yaml 'spec.ports[0].port' '{{ .Values.webhook.service.port }}'
	yq w -d1 -i $(CHARTDIR)/templates/webhook.yaml 'spec.ports[0].targetPort' '{{ .Values.webhook.service.targetPort }}'
	yq w -d1 -i $(CHARTDIR)/templates/webhook.yaml 'spec.ports[0].protocol' TCP
	yq w -d1 -i $(CHARTDIR)/templates/webhook.yaml 'spec.ports[0].name' http
	yq w -d1 -i $(CHARTDIR)/templates/webhook.yaml spec.selector.app '{{ include "dashboard.name" . }}'
	yq w -d1 -i $(CHARTDIR)/templates/webhook.yaml spec.selector.release '{{ .Release.Name }}'
	#   configurations
	number=0 ; while [ "$$(yq r -d2 config/webhook/webhook.yaml webhooks[$$number])" != "null" ] ; do \
		yq w -d2 -i $(CHARTDIR)/templates/webhook.yaml webhooks[$$number].clientConfig.caBundle '{{ $$cert.Cert | b64enc }}' ;\
		yq w -d2 -i $(CHARTDIR)/templates/webhook.yaml webhooks[$$number].clientConfig.service.name '{{ include "dashboard.fullname" . }}-webhook' ;\
		yq w -d2 -i $(CHARTDIR)/templates/webhook.yaml webhooks[$$number].clientConfig.service.namespace '{{ .Release.Namespace }}' ;\
		number=`expr $$number + 1` ; \
	done
	cat hack/generated-by-makefile.yaml \
		hack/certificate-globals.yaml \
		chart/dashboard/templates/webhook.yaml \
		> chart/dashboard/templates/_webhook.yaml
	sed -i "s/'{{ \.Values\.webhook\.service\.targetPort }}'/{{ \.Values\.webhook\.service\.targetPort }}/g" chart/dashboard/templates/_webhook.yaml
	sed -i "s/'{{ \.Values\.webhook\.service\.port }}'/{{ \.Values\.webhook\.service\.port }}/g" chart/dashboard/templates/_webhook.yaml
	mv chart/dashboard/templates/_webhook.yaml chart/dashboard/templates/webhook.yaml

.PHONY: chart
chart:
	yq w -i $(CHARTDIR)/Chart.yaml version "$(APP_VERSION)"
	yq w -i $(CHARTDIR)/Chart.yaml appVersion "$(APP_VERSION)"
	mv $(CHARTDIR)/values.yaml $(CHARTDIR)/_values.yaml
	sed 's#$(REGISTRY)/$(IMAGE_NAME):latest#$(REGISTRY)/$(IMAGE_NAME):$(APP_VERSION)#g' $(CHARTDIR)/_values.yaml > $(CHARTDIR)/values.yaml
	rm $(CHARTDIR)/_values.yaml

.PHONY: bundle
bundle:
	$(BINDIR)/packr -i $(CURDIR)/pkg -v

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate:
	go generate ./pkg/... ./cmd/...

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
	GOBIN=$(BINDIR) go install ./vendor/github.com/gobuffalo/packr/packr
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $(BINDIR) v1.10.2
	$(MAKE) -C app dependencies
