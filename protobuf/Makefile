

all: clean relay

.Phony: clean relay

clean:
	rm -rf relay_proto
	mkdir -p relay_proto

relay: types/shared_go types/queries_go types/query_services_go
# relay: types/shared_go types/queries_go types/query_services_go types/msgs_go

types/shared_go: relay_proto/x/relay/types/shared.proto
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/x/relay/types/shared.proto

types/queries_go: relay_proto/x/relay/types/shared.pb.go relay_proto/x/relay/types/queries.proto
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/x/relay/types/queries.proto

relay/query_services_go: relay/queries.pb.go relay_proto/queries.proto
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/query_services.proto

types/query_services_go: relay_proto/x/relay/types/queries.pb.go relay_proto/x/relay/types/queries.proto
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/x/relay/types/query_services.proto

types/msgs_go: relay_proto/x/relay/types/msgs.proto
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/x/relay/types/msgs.proto
