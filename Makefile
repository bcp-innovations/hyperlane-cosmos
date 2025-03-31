#!/usr/bin/make -f

DOCKER := $(shell which docker)

COMMIT := $(shell git log -1 --format='%H')
VERSION := nightly

ldflags := $(LDFLAGS)
ldflags += -X github.com/cosmos/cosmos-sdk/version.Name=hyperlane \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=hypd \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -s -w

ldflags := $(strip $(ldflags))

BUILD_FLAGS := -ldflags '$(ldflags)' -trimpath -tags 'ledger' -buildvcs=false

all: proto-all format lint test build-simapp

#################
###   Build   ###
#################

build-simapp:
	@echo "--> Building simapp..."
	@go build $(BUILD_FLAGS) -o "$(PWD)/build/" ./tests/hypd
	@echo "--> Completed build!"

release-simapp:
	@echo "--> Release simapp..."
	@for b in darwin:amd64 darwin:arm64 linux:amd64 linux:arm64; do \
    		os=$$(echo $$b | cut -d':' -f1); \
    		arch=$$(echo $$b | cut -d':' -f2); \
    		echo "--> Building "$$os" "$$arch""; \
    		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build $(BUILD_FLAGS) -o release/hypd_"$$os"_"$$arch" ./tests/hypd; \
    done

test:
	@echo "--> Running tests"
	@go test -cover -mod=readonly ./x/... ./util/...


.PHONY: build-simapp release-simapp test

##################
###  Protobuf  ###
##################

protoVer=0.15.3
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -u 0 -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "--> Generating protobuf files..."
	@$(protoImage) sh ./proto/protocgen.sh
	@go mod tidy

proto-format:
	@echo "--> Running protobuf formatter..."
	@$(protoImage) find ./proto -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@echo "--> Running protobuf linter..."
	@$(protoImage) buf lint proto/ --error-format=json

.PHONY: proto-all proto-gen proto-format proto-lint

#################
###  Linting  ###
#################

gofumpt_cmd=mvdan.cc/gofumpt
golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

format:
	@echo "--> Running Go formatter..."
	@go run $(gofumpt_cmd) -l -w .

lint:
	@echo "--> Running Go linter..."
	@go run $(golangci_lint_cmd) run --exclude-dirs scripts --timeout=10m

.PHONY: format lint