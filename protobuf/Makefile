protoc_go: example/addressbook.proto
	mkdir -p example_output # make directory for go package
	protoc $$PROTO_PATH --go_out=example_output example/addressbook.proto
	@touch protoc_go
