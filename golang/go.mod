module github.com/summa-tx/relays/golang

go 1.15

require (
	github.com/bombsimon/wsl v1.2.8 // indirect
	github.com/cosmos/cosmos-sdk v0.41.4
	github.com/go-critic/go-critic v0.4.0 // indirect
	github.com/gogo/protobuf v1.4.3
	github.com/golangci/gocyclo v0.0.0-20180528144436-0a533e8fa43d // indirect
	github.com/golangci/golangci-lint v1.21.0 // indirect
	github.com/golangci/revgrep v0.0.0-20180812185044-276a5c0a1039 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gostaticanalysis/analysisutil v0.0.3 // indirect
	github.com/securego/gosec v0.0.0-20191119104125-df484bfa9e9f // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/stumble/gorocksdb v0.0.3 // indirect
	github.com/summa-tx/bitcoin-spv/golang v1.4.0
	github.com/summa-tx/relays/proto v0.0.0-00010101000000-000000000000
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/iavl v0.12.4 // indirect
	github.com/tendermint/tendermint v0.34.8
	github.com/tendermint/tm-db v0.6.4
	github.com/uudashr/gocognit v1.0.0 // indirect
	mvdan.cc/unparam v0.0.0-20191111180625-960b1ec0f2c2 // indirect
	sourcegraph.com/sqs/pbtypes v1.0.0 // indirect
)

// only way to get it to build?? --> https://github.com/cosmos/cosmos-sdk/blob/v0.40.0-rc6/go.mod#L59
replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

// https://github.com/99designs/keyring/issues/64#issuecomment-742903794
replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

// TODO: remove before merging to master
replace github.com/summa-tx/relays/proto => ../proto
