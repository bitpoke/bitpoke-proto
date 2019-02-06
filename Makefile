BINDIR                   ?= $(CURDIR)/bin
PROTOTOOL_VERSION        := 1.3.0
PROTOC_VERSION           := 3.6.1
PROTOC_GEN_GOGO_VERSION  := 1.2.0
GRPC_VERSION             := 1.17.0
PROTOC_GEN_LINT_VERSION  := 0.2.1

PATH := $(BINDIR):$(CURDIR)/node_modules/.bin:$(PATH)
SHELL := env 'PATH=$(PATH)' /bin/sh

.PHONY: build
build:
	$(BINDIR)/prototool generate src

.PHONY: clean
clean:
	rm -rf build

.PHONY: lint
lint:
	$(BINDIR)/prototool format -l src
	$(BINDIR)/prototool lint src

.PHONY: test
test: build
	test -e "$(PWD)/build/go/src/github.com/presslabs/dashboard-go/vendor" || ln -s "$(PWD)/vendor" "$(PWD)/build/go/src/github.com/presslabs/dashboard-go"
	GOPATH="$(PWD)/build/go" go test ./build/go/src/github.com/presslabs/dashboard-go/pkg/...

.PHONY: dependencies
dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)

	npm install ts-protoc-gen@0.8.0

	curl -sSL https://github.com/uber/prototool/releases/download/v$(PROTOTOOL_VERSION)/prototool-$(shell uname -s)-$(shell uname -m) \
		-o $(BINDIR)/prototool && \
		chmod +x $(BINDIR)/prototool

	rm -rf /tmp/prototool-bootstrap
	mkdir -p /tmp/prototool-bootstrap
	echo 'protoc:'            >  /tmp/prototool-bootstrap/prototool.yaml
	echo '  version: 3.6.1'   >> /tmp/prototool-bootstrap/prototool.yaml
	echo 'syntax = "proto3";' >  /tmp/prototool-bootstrap/tmp.proto
	cat /tmp/prototool-bootstrap/prototool.yaml
	$(BINDIR)/prototool compile /tmp/prototool-bootstrap
	rm -rf /tmp/prototool-bootstrap

	GOBIN=$(BINDIR) go install ./vendor/github.com/ckaznocha/protoc-gen-lint
	GOBIN=$(BINDIR) go install ./vendor/github.com/gogo/protobuf/proto
	GOBIN=$(BINDIR) go install ./vendor/github.com/gogo/protobuf/jsonpb
	GOBIN=$(BINDIR) go install ./vendor/github.com/gogo/protobuf/protoc-gen-gogo
	GOBIN=$(BINDIR) go install ./vendor/github.com/gogo/protobuf/gogoproto
	GOBIN=$(BINDIR) go install ./vendor/google.golang.org/grpc
