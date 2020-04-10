PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=Relay \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=relayd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=relaycli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

all: lint install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/relayd
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/relaycli

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

lint:
	golangci-lint run
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

test:
	@go test -mod=readonly $(PACKAGES)

init:
	./scripts/init_chain.sh
