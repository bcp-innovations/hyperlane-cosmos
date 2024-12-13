#!/usr/bin/make -f

COMMIT := $(shell git log -1 --format='%H')
VERSION := 0.1

ldflags = -X main.AppName=hypertools \
		  -X main.Version=$(VERSION) \
		  -X main.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)' -trimpath -mod=readonly

.PHONY: build

all: build

###############################################################################
###                                  Build                                  ###
###############################################################################

build:
	@echo "🤖 Building Hypertools ..."
	@go build $(BUILD_FLAGS) -o "$(PWD)/build/" ./
	@echo "✅ Completed build!"
