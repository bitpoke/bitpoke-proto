NODE_ENV ?= production

NODE_MODULES ?= $(PWD)/node_modules
REACT        := $(NODE_MODULES)/.bin/react-scripts-ts
TSLINT       := $(NODE_MODULES)/.bin/tslint

deps:
	yarn install --modules-folder $(NODE_MODULES)

run:
	NODE_ENV=development BROWSER=none $(REACT) start

build:
	NODE_ENV=$(NODE_ENV) $(REACT) build

test:
	NODE_ENV=$(NODE_ENV) $(REACT) test --env=jsdom

lint:
	$(TSLINT) --project .

clean:
	rm -rf build
	rm -rf $(NODE_MODULES)

.PHONY: deps run test lint clean
