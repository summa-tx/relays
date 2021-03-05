

all: clean msgs_go protoc_go

clean:
	rm -rf relay
	mkdir -p relay
	rm -f msgs_go
	rm -f protoc_go

msgs_go: relay_proto/x/relay/types/msgs.proto
	mkdir -p relay
	protoc $$PROTO_PATH --go_out=relay_proto relay_proto/x/relay/types/msgs.proto
	@touch msgs_go

protoc_go: example/addressbook.proto
	mkdir -p example_output # make directory for go package
	protoc $$PROTO_PATH --go_out=example_output example/addressbook.proto
	@touch protoc_go
