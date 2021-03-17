module github.com/summa-tx/relays/golang

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.37.11
	github.com/gorilla/mux v1.7.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.5.1
	github.com/summa-tx/bitcoin-spv/golang v1.4.0
	github.com/summa-tx/relays/proto v0.0.0-00010101000000-000000000000
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.32.10
	github.com/tendermint/tm-db v0.2.0
)

replace github.com/summa-tx/relays/proto => ../proto
