NODE_ENV ?= production

NODE_MODULES ?= $(PWD)/node_modules
REACT        := $(NODE_MODULES)/.bin/react-scripts-ts
TSLINT       := $(NODE_MODULES)/.bin/tslint
PROTOC_GEN   := $(NODE_MODULES)/.bin/protoc-gen-ts

PROTO_OUT_DIR ?= ./src/proto

deps:
	yarn install --modules-folder $(NODE_MODULES)

run:
	NODE_ENV=development BROWSER=none $(REACT) start

build:
	NODE_ENV=$(NODE_ENV) $(REACT) build

generate:
	mkdir -p $(PROTO_OUT_DIR)
	protoc -I ../proto \
		--plugin="protoc-gen-ts=${PROTOC_GEN}" \
		--js_out="import_style=commonjs,binary:${PROTO_OUT_DIR}" \
		--ts_out="service=true:${PROTO_OUT_DIR}" \
		presslabs/dashboard/core/v1/project.proto

test:
	NODE_ENV=$(NODE_ENV) $(REACT) test --env=jsdom

lint:
	$(TSLINT) --project .

clean:
	rm -rf build
	rm -rf $(NODE_MODULES)

.PHONY: deps run test lint clean
