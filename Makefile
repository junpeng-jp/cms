# based on https://tech.davis-hansson.com/p/make/
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
.DEFAULT_GOAL := help
.DELETE_ON_ERROR:

## variables
PROTO_PATH=$(shell cd .. && pwd)/protos
NODE_PROTO=node.proto
FILE_PROTO=file.proto

GO_MODULE_PREFIX=github.com/junpeng-jp/blog
GO_PB_MODULE_PREFIX=${GO_MODULE_PREFIX}/internal/file/filepb

.PHONY: help
help:  ## print help message
	@grep -E '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## protbuf
.PHONY: gen-pb
gen-pb: ## generate the protobuf files used by the content
	protoc \
	--go_out=. \
	--go_opt=module=${GO_MODULE_PREFIX} \
	--go_opt=M${NODE_PROTO}=${GO_PB_MODULE_PREFIX} \
	--go_opt=M${FILE_PROTO}=${GO_PB_MODULE_PREFIX} \
	--go-vtproto_out=. \
	--go-vtproto_opt=module=${GO_MODULE_PREFIX} \
	--go-vtproto_opt=M${NODE_PROTO}=${GO_PB_MODULE_PREFIX} \
	--go-vtproto_opt=M${FILE_PROTO}=${GO_PB_MODULE_PREFIX} \
	--go-vtproto_opt=features=marshal+unmarshal+size \
	--go-vtproto_opt=features=marshal+unmarshal+size \
	--proto_path=${PROTO_PATH} \
	${PROTO_PATH}/${FILE_PROTO} ${PROTO_PATH}/${NODE_PROTO} 

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
	go build -ldflags="-s -w" -o build/toolkit cmd/toolkit/*.go

.PHONY: build-protoc-gen-tinygo
build-protoc-gen-tinygo:
	go build -ldflags="-s -w" -o build/protoc-gen-tinygo cmd/protoc-gen-tinygo/main.go
	