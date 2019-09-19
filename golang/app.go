package relay

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/libs/log"

	// "github.com/cosmos/cosmos-sdk/x/auth/genaccounts"
	//
	// "github.com/cosmos/cosmos-sdk/x/bank"
	// distr "github.com/cosmos/cosmos-sdk/x/distribution"
	// "github.com/cosmos/cosmos-sdk/x/params"
	// "github.com/cosmos/cosmos-sdk/x/staking"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// abci "github.com/tendermint/tendermint/abci/types"
	// cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tm-db"
)

const appName = "summa-cosmos-relay"

var (
	// DefaultCLIHome is default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.summa/cosmosrelay")

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.summa/cosmosrelay/config")

	// ModuleBasics is  basic module elemnets
	ModuleBasics = module.NewBasicManager()
)

type relayApp struct {
	*bam.BaseApp
	cdc *codec.Codec
}

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

// NewRelayApp instantiates a new Bitcoin relay app
func NewRelayApp(logger log.Logger, db dbm.DB) *relayApp {
	// First define the top level codec that will be shared by the different modules. Note: Codec will be explained later
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	var app = &relayApp{
		BaseApp: bApp,
		cdc:     cdc,
	}

	return app
}
