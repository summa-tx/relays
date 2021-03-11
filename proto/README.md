## Setup

1. Install Go:

    - [link](https://golang.org/doc/install)
    - Add to path:

        ```sh
        echo 'export GOPATH=$HOME/Go' >> $HOME/.zshrc
        source $HOME/.zshrc
        ```

    - Check it installed properly: `go version`

2. Install Protobuf:

    - `brew install protobuf`
    - Add to path:

        ```sh
        echo 'export PATH=$PATH:$GOPATH/bin' >> $HOME/.zshrc
        source $HOME/.zshrc
        ```

    - Install the Go protocol buffers plugin:

        ```sh
        go get google.golang.org/protobuf/cmd/protoc-gen-go \
               google.golang.org/grpc/cmd/protoc-gen-go-grpc
        ```

## Building protobuf files

To compile the Go code, simply run this command with the relevant file paths:

```sh
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/$SRC_FILE.proto
```

- `SRC_DIR`: the protobuf directory to compile
- `DST_DIR`: the destination directory where the Go code will be deposited
- `$SRC_FILE`: the specific file you wish to compile

For example, this command will take `example/addressbook.proto` and generate
`addressbook.pb.go` inside `/example_output`:

```sh
protoc -I=example --go_out=example_output example/addressbook.proto
```
