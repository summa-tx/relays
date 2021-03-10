# Location of .proto files
INDIR = ./proto

# Options for protobuf generation command
PROTOC_OPTS = --proto_path=proto \
	--go_out=plugins=grpc:. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative

all: clean relay

.PHONY: clean relay

relay:
	protoc $(PROTOC_OPTS) $(INDIR)/*.proto
	go build ./...

clean:
	rm -rf *.pb.go
