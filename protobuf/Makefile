

all: clean queries_go protoc_go

clean:
	rm -rf relay
	mkdir -p relay
	rm -f queries_go
	rm -f protoc_go

queries_go: relay_proto/queries.proto
	mkdir -p relay
	protoc $$PROTO_PATH --go_out=relay relay_proto/queries.proto
	@touch queries_go

protoc_go: example/addressbook.proto
	mkdir -p example_output # make directory for go package
	protoc $$PROTO_PATH --go_out=example_output example/addressbook.proto
	@touch protoc_go

