BINDIR                   ?= $(CURDIR)/bin
PROTOTOOL_VERSION        := 1.3.0
PROTOC_VERSION           := 3.6.1
PROTOC_GEN_GO_VERSION    := 1.2.0
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
test:
	$(BINDIR)/prototool compile src

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

	go get -d -u github.com/ckaznocha/protoc-gen-lint
	git -C "$(GOPATH)"/src/github.com/ckaznocha/protoc-gen-lint checkout v$(PROTOC_GEN_LINT_VERSION) --quiet
	GOBIN=$(BINDIR) go install github.com/ckaznocha/protoc-gen-lint
	git -C "$(GOPATH)"/src/github.com/ckaznocha/protoc-gen-lint checkout master --quiet

	go get -d -u github.com/golang/protobuf/protoc-gen-go
	git -C "$(GOPATH)"/src/github.com/golang/protobuf checkout v$(PROTOC_GEN_GO_VERSION) --quiet
	GOBIN=$(BINDIR) go install github.com/golang/protobuf/protoc-gen-go
	git -C "$(GOPATH)"/src/github.com/golang/protobuf checkout master --quiet

	go get -d -u google.golang.org/grpc
	git -C "$(GOPATH)"/src/google.golang.org/grpc checkout v$(GRPC_VERSION) --quiet
	GOBIN=$(BINDIR) go install google.golang.org/grpc
	git -C "$(GOPATH)"/src/google.golang.org/grpc checkout master --quiet
