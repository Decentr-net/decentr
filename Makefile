PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

CLIENT=decentrcli
SERVER=decentrd

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=decentr \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=$(SERVER) \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=$(CLIENT) \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

build:
		go build -mod=vendor $(BUILD_FLAGS) -o build/$(SERVER) ./cmd/$(SERVER)
		go build -mod=vendor $(BUILD_FLAGS) -o build/$(CLIENT) ./cmd/$(CLIENT)

install: go.sum
		go install -mod=vendor $(BUILD_FLAGS) ./cmd/$(SERVER)
		go install -mod=vendor $(BUILD_FLAGS) ./cmd/$(CLIENT)

go.sum:
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test:
		go test -mod=readonly $(PACKAGES)

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify

start:
	bash init.sh

