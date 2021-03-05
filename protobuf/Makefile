

all: clean relay

.PHONY: clean relay

clean:
	rm -rf relay
	mkdir -p relay
	rm -f queries_go
	rm -f protoc_go
	mkdir -p relay

# relay_proto: relay_proto/shared.pb.go relay_proto/x/relay/types/queries.pb.go relay_proto/x/relay/types/query_services.pb.go

types/shared: relay_proto/x/relay/types/shared.proto
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/x/relay/types/shared.proto

types/queries: relay_proto/x/relay/types/shared.pb.go relay_proto/x/relay/types/queries.proto
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/x/relay/types/queries.proto

types/query_services: relay_proto/x/relay/types/queries.pb.go relay_proto/x/relay/types/queries.proto
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/x/relay/types/query_services.proto
