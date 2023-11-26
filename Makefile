PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COSMOS_PKG_VERSION := $(shell go list -m github.com/cosmos/cosmos-sdk | sed 's:.* ::')
COMMIT := $(shell git log -1 --format='%H')

export GO111MODULE = on

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=decentr \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=decentrd \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

build: check-go-version go.sum
ifeq ($(OS), Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/$(shell go env GOOS)/decentrd.exe ./cmd/decentrd
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/$(shell go env GOOS)/decentrd ./cmd/decentrd
endif

install: check-go-version go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/decentrd

linux: check-go-version go.sum
	GOOS=linux GOARCH=amd64 $(MAKE) build

# Add check to make sure we are using the proper Go version before proceeding with anything
check-go-version:
	@if ! go version | grep -q "go1.19"; then \
		echo "\033[0;31mERROR:\033[0m Go version 1.19 is required for compiling decentrd. It looks like you are using" "$(shell go version) \nThere are potential consensus-breaking changes that can occur when running binaries compiled with different versions of Go. Please download Go version 1.19 and retry. Thank you!"; \
		exit 1; \
	fi

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
docker_buf := docker run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf
buildtools=decentr/buildtools:v0.1
containerProtoGen=decentr-buildtools-protogen
containerProtoSwaggerGen=decentr-buildtools-protoswaggergen
containerProtoFmt=decentr-buildtools-protofmt

proto-all: check-go-version proto-lint proto-gen

proto-gen: check-go-version
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(buildtools) \
		sh ./scripts/protocgen.sh; fi

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoSwaggerGen}$$"; then docker start -a $(containerProtoSwaggerGen); else docker run --name $(containerProtoSwaggerGen) -v $(CURDIR):/workspace --workdir /workspace $(buildtools) \
		sh ./scripts/protoc-swagger-gen.sh; fi

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoFmt}$$"; then docker start -a $(containerProtoFmt); else docker run --name $(containerProtoFmt) -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./ -path "./proto/decentr/*" -name *.proto -exec clang-format -style=file -i {} \; ; fi

proto-lint:
	@$(docker_buf) lint --error-format=json

proto-check-breaking:
	@$(docker_buf) breaking --against $(HTTPS_GIT)

proto-update-deps:
	rm -rf proto/3rdparty
	mkdir proto/3rdparty
	git clone --depth 1 --branch $(COSMOS_PKG_VERSION) git@github.com:cosmos/cosmos-sdk.git
	mv cosmos-sdk/third_party/proto/* ./proto/3rdparty/
	mv cosmos-sdk/proto/cosmos ./proto/3rdparty/cosmos
	rm -rf cosmos-sdk

.PHONY: build vendor proto-all proto-gen proto-swagger-gen proto-format proto-lint proto-check-breaking proto-update-deps