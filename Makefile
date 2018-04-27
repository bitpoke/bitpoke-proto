PACKAGE_NAME := github.com/presslabs/dashboard
REGISTRY := gcr.io/pl-infra
APP_NAME := presslabs-dashboard
IMAGE_TAGS := canary
GOPATH ?= $HOME/go
HACK_DIR ?= hack
BUILD_TAG := build

# Get a list of all binaries to be built
CMDS := $(shell find ./cmd/ -maxdepth 1 -type d -exec basename {} \; | grep -v cmd)

ifeq ($(APP_VERSION),)
APP_VERSION := $(shell git describe --abbrev=4 --dirty --tags --always)
endif

GIT_COMMIT ?= $(shell git rev-parse HEAD)

ifeq ($(shell git status --porcelain),)
	GIT_STATE ?= clean
else
	GIT_STATE ?= dirty
endif

# Go build flags
GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH ?= amd64
GOLDFLAGS := -ldflags "-X $(PACKAGE_NAME)/pkg/version.AppGitState=${GIT_STATE} -X $(PACKAGE_NAME)/pkg/version.AppGitCommit=${GIT_COMMIT} -X $(PACKAGE_NAME)/pkg/version.AppVersion=${APP_VERSION}"

# Get a list of all binaries to be built
CMDS := $(shell find ./cmd/ -maxdepth 1 -type d -exec basename {} \; | grep -v cmd)
SRC_CMDS := $(patsubst %, cmd/%, $(CMDS))
BIN_CMDS := $(patsubst %, bin/dashboard-%_$(GOOS)_$(GOARCH), $(CMDS))

.DEFAULT_GOAL := bin/dashboard-controller_$(GOOS)_$(GOARCH)

.PHONY: run
run: bin/dashboard-controller_$(GOOS)_$(GOARCH)
	./bin/dashboard-controller_$(GOOS)_$(GOARCH)

# Code building targets
#######################

.PHONY: test
test:
	go test -v \
	    -race \
		$$(go list ./... | \
			grep -v '/vendor/' | \
			grep -v '/test/e2e' | \
			grep -v '/pkg/client' \
		)

.PHONY: full-test
full-test: generate_verify test

.PHONY: lint
lint:
	@set -e; \
	GO_FMT=$$(git ls-files *.go | grep -v 'vendor/' | xargs gofmt -d); \
	if [ -n "$${GO_FMT}" ] ; then \
		echo "Please run go fmt"; \
		echo "$$GO_FMT"; \
		exit 1; \
	fi

.PHONY: build
build: $(BIN_CMDS)

.PHONY: $(SRC_CMDS)
bin/dashboard-%_darwin_amd64: cmd/%
	test -d bin || mkdir bin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v \
		-tags netgo \
		-o $@ \
		$(GOLDFLAGS) \
		./$<
bin/dashboard-%_linux_amd64: cmd/%
	test -d bin || mkdir bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v \
		-tags netgo \
		-o $@ \
		$(GOLDFLAGS) \
		./$<


# Code generation targets
#########################

.PHONY: generate generate_verify
generate:
	$(HACK_DIR)/update-codegen.sh

generate_verify:
	$(HACK_DIR)/verify-codegen.sh


# Docker image targets
######################
images: bin/dashboard-controller_linux_amd64
	docker build \
		--build-arg VCS_REF=$(GIT_COMMIT) \
		-t $(REGISTRY)/$(APP_NAME):$(BUILD_TAG) \
		-f ./Dockerfile .
	set -e; \
		for tag in $(IMAGE_TAGS); do \
			docker tag $(REGISTRY)/$(APP_NAME):$(BUILD_TAG) $(REGISTRY)/$(APP_NAME):$${tag} ; \
	done

publish: images
	set -e; \
		for tag in $(IMAGE_TAGS); do \
		gcloud docker -- push $(REGISTRY)/$(APP_NAME):$${tag}; \
	done
