

all: clean relay

.Phony: clean relay

clean:
	rm -rf relay
	mkdir -p relay
	rm -f msgs_go
	rm -f queries_go
	rm -f protoc_go

relay: relay/shared.pb.go relay/queries.pb.go relay/query_services.pb.go

relay/shared.pb.go: relay_proto/shared.proto
	protoc $$PROTO_PATH --go_out=relay relay_proto/shared.proto

relay/queries.pb.go: relay/shared.pb.go relay_proto/queries.proto
	protoc $$PROTO_PATH --go_out=relay relay_proto/queries.proto

relay/query_services.pb.go: relay/queries.pb.go relay_proto/queries.proto
	protoc $$PROTO_PATH --go_out=relay relay_proto/query_services.proto

msgs_go: relay_proto/x/relay/types/msgs.proto
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/x/relay/types/msgs.proto

