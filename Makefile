# based on https://tech.davis-hansson.com/p/make/
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
.DEFAULT_GOAL := help
.DELETE_ON_ERROR:

## variables
GO_MODULE_PREFIX=github.com/junpeng-jp/blog

.PHONY: help
help:  ## print help message
	@grep -E '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## protbuf
.PHONY: gen-blocks
gen-blocks: ## generate the protobuf files used for blocks
	protoc \
	--go_out=. \
	--go-vtproto_out=. \
	--go_opt=module=${GO_MODULE_PREFIX} \
	--go-vtproto_opt=module=${GO_MODULE_PREFIX} \
	--go_opt=Mprotos/blocks.proto=${GO_MODULE_PREFIX}/internal/pb/block \
	--go-vtproto_opt=Mprotos/blocks.proto=${GO_MODULE_PREFIX}/internal/pb/block \
	--go-vtproto_opt=features=marshal+unmarshal+size \
	./protos/blocks.proto

## checks
.PHONY: format
format: ## run lint and auto code formatting
	echo "TODO: gofmt"

.PHONY: lint
lint: ## run lint on staged changes
	echo "TODO: golint"

.PHONY: test
test: ## run unit tests
	echo "TODO: gotest"

## productionise
.PHONY: build
build:
	go build -ldflags="-s -w" -o build/toolkit cmd/toolkit/*
	