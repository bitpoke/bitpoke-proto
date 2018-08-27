
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
KUBEBUILDER_VERSION ?= 1.0.0

all: test dashboard

# Run tests
test: generate manifests
	go test ./pkg/... ./cmd/... -coverprofile cover.out

# Build dashboard binary
dashboard: generate fmt vet
	go build -o bin/dashboard github.com/presslabs/dashboard/cmd/dashboard

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/dashboard/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crds

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crds
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all

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
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

lint: vet
	gometalinter.v2 --disable-all --deadline 5m \
	--enable=vetshadow \
	--enable=misspell \
	--enable=structcheck \
	--enable=golint \
	--enable=deadcode \
	--enable=goimports \
	--enable=errcheck \
	--enable=varcheck \
	--enable=goconst \
	--enable=gas \
	--enable=unparam \
	--enable=ineffassign \
	--enable=nakedret \
	--enable=interfacer \
	--enable=misspell \
	--enable=gocyclo \
	--line-length=170 \
	--enable=lll \
	--dupl-threshold=400 \
	--enable=dupl \
	--enable=maligned \
	./pkg/... ./cmd/...

dependencies:
	go get -u gopkg.in/alecthomas/gometalinter.v2
	gometalinter.v2 --install

	# install Kubebuilder
	curl -L -O https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VERSION}/kubebuilder_${KUBEBUILDER_VERSION}_linux_amd64.tar.gz
	tar -zxvf kubebuilder_${KUBEBUILDER_VERSION}_linux_amd64.tar.gz
	mv kubebuilder_${KUBEBUILDER_VERSION}_linux_amd64 -T /usr/local/kubebuilder
	export PATH=$PATH:/usr/local/kubebuilder/bin
