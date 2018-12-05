BINDIR           ?= $(CURDIR)/bin

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
	go get -u github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
	GOBIN=$(BINDIR) go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
	curl -sSL https://github.com/uber/prototool/releases/download/v1.3.0/prototool-$(shell uname -s)-$(shell uname -m) \
	     -o $(BINDIR)/prototool && \
	     chmod +x $(BINDIR)/prototool
