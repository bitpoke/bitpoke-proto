NODE_ENV ?= production

NODE_MODULES ?= $(CURDIR)/node_modules
BINDIR       := $(NODE_MODULES)/.bin
TSLINT       := $(BINDIR)/tslint
REACT        := $(BINDIR)/react-scripts
SERVER       := $(NODE_MODULES)/.bin/http-server

PROTO_OUT_DIR ?= ./src/proto

PATH := $(BINDIR):$(PATH)
SHELL := env 'PATH=$(PATH)' /bin/sh

dependencies:
	yarn install --modules-folder $(NODE_MODULES)

run:
	NODE_ENV=development BROWSER=none $(REACT) start

.PHONY: build
build:
	NODE_ENV=$(NODE_ENV) $(REACT) build

test:
	NODE_ENV=$(NODE_ENV) $(REACT) test --env=jsdom --passWithNoTests

run-build: build
	NODE_ENV=$(NODE_ENV) $(SERVER) ./build -p 3000

lint:
	$(TSLINT) --project $(CURDIR)

clean:
	rm -rf build
	rm -rf $(NODE_MODULES)

.PHONY: dependencies run test lint clean
