BINDIR                   ?= $(CURDIR)/bin
APP_BINDIR               ?= $(CURDIR)/../app/bin
PROTOTOOL_VERSION        := 1.3.0
PROTOC_VERSION           := 3.6.1
PROTOC_GEN_GO_VERSION    := 1.2.0
GRPC_VERSION             := 1.17.0
PROTOC_GEN_LINT_VERSION  := 0.2.1

PATH := $(BINDIR):$(APP_BINDIR):$(PATH)
SHELL := env 'PATH=$(PATH)' /bin/sh

build:
	$(BINDIR)/prototool generate

lint:
	$(BINDIR)/prototool format -l
	$(BINDIR)/prototool lint

clean:
	find $(CURDIR)/../pkg -name '*.pb.go' -delete
	find $(CURDIR)/../app/src -name '*_pb*.js' -delete

test:
	$(BINDIR)/prototool compile

dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)
	go get -d -u github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
	GOBIN=$(BINDIR) go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators

	curl -sSL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m).zip \
		-o protoc.zip \
		&& unzip -u -o protoc.zip -d protoc \
		&& mv protoc/bin/protoc $(BINDIR) \
		&& chmod +x $(BINDIR)/protoc \
		&& rm -rf protoc

	curl -sSL https://github.com/uber/prototool/releases/download/v$(PROTOTOOL_VERSION)/prototool-$(shell uname -s)-$(shell uname -m) \
		-o $(BINDIR)/prototool && \
		chmod +x $(BINDIR)/prototool

	go get -d -u github.com/golang/protobuf/protoc-gen-go
	git -C "$(GOPATH)"/src/github.com/golang/protobuf checkout v$(PROTOC_GEN_GO_VERSION) --quiet
	GOBIN=$(BINDIR) go install github.com/golang/protobuf/protoc-gen-go
	git -C "$(GOPATH)"/src/github.com/golang/protobuf checkout master --quiet

	go get -d -u google.golang.org/grpc
	git -C "$(GOPATH)"/src/google.golang.org/grpc checkout v$(GRPC_VERSION) --quiet
	GOBIN=$(BINDIR) go install google.golang.org/grpc
	git -C "$(GOPATH)"/src/google.golang.org/grpc checkout master --quiet

	go get -d -u github.com/ckaznocha/protoc-gen-lint
	git -C "$(GOPATH)"/src/github.com/ckaznocha/protoc-gen-lint checkout v$(PROTOC_GEN_LINT_VERSION) --quiet
	GOBIN=$(BINDIR) go install github.com/ckaznocha/protoc-gen-lint
	git -C "$(GOPATH)"/src/github.com/ckaznocha/protoc-gen-lint checkout master --quiet
