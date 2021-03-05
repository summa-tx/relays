

all: clean relay

.PHONY: clean relay

clean:
	rm -rf relay
	mkdir -p relay
	rm -f queries_go
	rm -f protoc_go
	mkdir -p relay

relay: relay/shared.pb.go relay/queries.pb.go relay/query_services.pb.go

relay/shared.pb.go: relay_proto/shared.proto
	protoc $$PROTO_PATH --go_out=relay relay_proto/shared.proto

relay/queries.pb.go: relay/shared.pb.go relay_proto/queries.proto
	protoc $$PROTO_PATH --go_out=relay relay_proto/queries.proto

relay/query_services.pb.go: relay/queries.pb.go relay_proto/queries.proto
	protoc $$PROTO_PATH --go_out=relay relay_proto/query_services.proto
