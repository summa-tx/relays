

all: clean relay

.PHONY: clean relay

clean:
	rm -rf relay
	mkdir -p relay
	rm -f queries_go
	rm -f protoc_go
	mkdir -p relay

relay: types/shared_go types/queries_go types/query_services_go

types/shared_go: relay_proto/x/relay/types/shared.proto
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/x/relay/types/shared.proto

types/queries_go: relay_proto/x/relay/types/shared.pb.go relay_proto/x/relay/types/queries.proto
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/x/relay/types/queries.proto

types/query_services_go: relay_proto/x/relay/types/queries.pb.go relay_proto/x/relay/types/queries.proto
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/x/relay/types/query_services.proto
