PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=Relay \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=relayd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=relaycli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

all: clean relay lint install

.Phony: clean relay

clean:
	rm x/relay/types/shared.pb.go
	rm x/relay/types/querier.pb.go
	rm x/relay/types/query_services.pb.go
	rm x/relay/types/msgs.pb.go
	rm x/relay/types/msg_services.pb.go

relay: types/shared types/querier types/query_services types/msgs types/msg_services

types/shared: x/relay/types/shared.proto
	protoc $$PROTO_PATH --go_out=x x/relay/types/shared.proto

types/querier: x/relay/types/shared.pb.go x/relay/types/querier.proto
	protoc $$PROTO_PATH --go_out=x x/relay/types/querier.proto

types/query_services: x/relay/types/querier.pb.go x/relay/types/querier.proto
	protoc $$PROTO_PATH --go_out=x x/relay/types/query_services.proto

types/msgs: x/relay/types/msgs.proto
	protoc $$PROTO_PATH --go_out=x x/relay/types/msgs.proto

types/msg_services: x/relay/types/msg_services.proto
	protoc $$PROTO_PATH --go_out=x x/relay/types/msg_services.proto

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
