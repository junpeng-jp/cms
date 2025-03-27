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
GO_CONTENT_MODULE_PREFIX=${GO_MODULE_PREFIX}/internal/pb/contentpb
GO_DOCUMENT_MODULE_PREFIX=${GO_MODULE_PREFIX}/internal/pb/docpb
GO_FILE_MODULE_PREFIX=${GO_MODULE_PREFIX}/internal/pb/filepb

.PHONY: help
help:  ## print help message
	@grep -E '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## protbuf
.PHONY: gen-content
gen-content: ## generate the protobuf files used by the content
	protoc \
	--go_out=. \
	--go-vtproto_out=. \
	--go_opt=module=${GO_MODULE_PREFIX} \
	--go-vtproto_opt=module=${GO_MODULE_PREFIX} \
	--go_opt=Mprotos/content.proto=${GO_CONTENT_MODULE_PREFIX} \
	--go-vtproto_opt=Mprotos/content.proto=${GO_CONTENT_MODULE_PREFIX} \
	--go-vtproto_opt=features=marshal+unmarshal+size \
	./protos/content.proto

.PHONY: gen-file
gen-file: ## generate the protobuf files used for file encoding
	protoc \
	--go_out=. \
	--go-vtproto_out=. \
	--go_opt=module=${GO_MODULE_PREFIX} \
	--go-vtproto_opt=module=${GO_MODULE_PREFIX} \
	--go_opt=Mprotos/file_metadata.proto=${GO_FILE_MODULE_PREFIX} \
	--go-vtproto_opt=Mprotos/file_metadata.proto=${GO_FILE_MODULE_PREFIX} \
	--go-vtproto_opt=features=marshal+unmarshal+size \
	./protos/file_metadata.proto

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
	