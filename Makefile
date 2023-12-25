.SUFFIXES:
.DELETE_ON_ERROR:

MAKEFLAGS += --no-builtin-rules
SHELL := /bin/sh -e -c

all:

.PHONY: generate
generate:
	go generate -x ./...

.PHONY: webui-shell
webui-shell:
	./contrib/run-webdev sh

.PHONY: webui-build
webui-build:
	./contrib/run-webdev npm run build
