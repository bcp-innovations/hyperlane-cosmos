#!/usr/bin/make -f

DOCKER := $(shell which docker)

#################
###   Build   ###
#################

build-simapp:
	@echo "--> Building simapp..."
	@go build $(BUILD_FLAGS) -o "$(PWD)/build/" ./tests/cmd/hypd
	@echo "--> Completed build!"

test:
	@echo "--> Running tests"
	@go test -cover -mod=readonly ./x/...

.PHONY: build-simapp test

##################
###  Protobuf  ###
##################

protoVer=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating protobuf files..."
	@$(protoImage) sh ./scripts/protocgen.sh
	@go mod tidy

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint proto/ --error-format=json

.PHONY: proto-all proto-gen proto-format proto-lint

#################
###  Linting  ###
#################

gofumpt_cmd=mvdan.cc/gofumpt
golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

format:
	@echo "--> Running formatter"
	@go run $(gofumpt_cmd) -l -w .

lint:
	@echo "--> Running linter..."
	# TODO temporarily disabled govet
	@go run $(golangci_lint_cmd) run --exclude-dirs scripts --timeout=10m -D govet

.PHONY: format lint