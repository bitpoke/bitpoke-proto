NODE_ENV ?= production

NODE_MODULES     ?= $(CURDIR)/node_modules
REACT            := $(NODE_MODULES)/.bin/react-scripts-ts
TSLINT           := $(NODE_MODULES)/.bin/tslint

BINDIR ?= $(NODE_MODULES)/.bin

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

generate:
	$(MAKE) -C ../proto

test:
	# NODE_ENV=$(NODE_ENV) $(REACT) test --env=jsdom

lint:
	$(TSLINT) --project $(CURDIR)

clean:
	rm -rf build
	rm -rf $(NODE_MODULES)

.PHONY: dependencies run test lint clean
