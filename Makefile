PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COSMOS_PKG_VERSION := $(shell go list -m github.com/cosmos/cosmos-sdk | sed 's:.* ::')
COMMIT := $(shell git log -1 --format='%H')

export GO111MODULE = on

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=decentr \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=decentr \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

build: go.sum
ifeq ($(OS), Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/$(shell go env GOOS)/decentr.exe ./cmd/decentr
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/$(shell go env GOOS)/decentr ./cmd/decentr
endif

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/decentr

build-linux: go.sum
	GOOS=linux GOARCH=amd64 $(MAKE) build

### tools ###

clean:
	rm -rf build/

go.sum:
	@echo "--> Ensure dependencies have not been modified"
	go mod verify

vendor:
	go mod tidy
	go mod vendor

### qa ###

test:
	@echo "--> Running tests"
	go test -mod=readonly $(PACKAGES)


# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify

start:
	bash init.sh

### proto ###

proto-all: proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	./scripts/protocgen.sh

proto-lint:
	@buf lint --error-format=json | jq

proto-check-breaking:
	@buf breaking --against $(HTTPS_GIT)#branch=master

proto-update-deps:
	rm -rf proto/3rdparty
	mkdir proto/3rdparty
	git clone --depth 1 --branch $(COSMOS_PKG_VERSION) git@github.com:cosmos/cosmos-sdk.git
	mv cosmos-sdk/third_party/proto/* ./proto/3rdparty/
	mv cosmos-sdk/proto/cosmos ./proto/3rdparty/cosmos
	rm -rf cosmos-sdk

.PHONY: proto-all proto-gen proto-gen-any proto-swagger-gen proto-format proto-lint proto-check-breaking proto-update-deps