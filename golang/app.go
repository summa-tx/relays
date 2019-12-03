package relay

import (
	"encoding/json"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/summa-tx/relays/golang/x/relay"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const appName = "relay"

var (
	// DefaultCLIHome is default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.summa/cosmosrelay")

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.summa/cosmosrelay/config")

	// ModuleBasics is  basic module elemnets
	ModuleBasics = module.NewBasicManager(
		relay.AppModule{},
	)
)

type relayApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the substores
	keys map[string]*sdk.KVStoreKey

	relayKeeper relay.Keeper

	// Module Manager
	mm *module.Manager
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
func NewRelayApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *relayApp {
	// First define the top level codec that will be shared by the different modules. Note: Codec will be explained later
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)

	keys := sdk.NewKVStoreKeys(relay.StoreKey)

	// Here you initialize your application with the store keys it requires
	var app = &relayApp{
		BaseApp: bApp,
		cdc:     cdc,
		keys:    keys,
	}

	// The RelayKeeper is the Keeper from the module for this tutorial
	// It handles interactions with the store
	app.relayKeeper = relay.NewKeeper(
		keys[relay.StoreKey],
		app.cdc,
		true, // TODO: pass this in somehow
	)

	app.mm = module.NewManager(
		relay.NewAppModule(app.relayKeeper),
	)

	// app.mm.SetOrderBeginBlockers(distr.ModuleName, slashing.ModuleName)
	// app.mm.SetOrderEndBlockers(staking.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	app.mm.SetOrderInitGenesis(
		relay.ModuleName,
	)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// initialize stores
	app.MountKVStores(keys)

	err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState returns a new default genesis
func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func (app *relayApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *relayApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *relayApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *relayApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *relayApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	return modAccAddrs
}

//_________________________________________________________

func (app *relayApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string,
) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	genState := app.mm.ExportGenesis(ctx)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, nil
}
