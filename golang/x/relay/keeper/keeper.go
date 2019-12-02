package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"

	"github.com/summa-tx/relays/golang/x/relay/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey  sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc       *codec.Codec // The wire codec for binary encoding/decoding.
	IsMainNet bool
}

// NewKeeper instantiates a new keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, mainnet bool) Keeper {
	return Keeper{
		storeKey:  storeKey,
		cdc:       cdc,
		IsMainNet: mainnet,
	}
}

func (k Keeper) getPrefixStore(ctx sdk.Context, namespace string) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(namespace))
}

func (k Keeper) hasRelayGenesis(ctx sdk.Context) bool {
	store := k.getChainStore(ctx)
	return store.Has([]byte(types.RelayGenesisStorage))
}

// setRelayGenesis sets the first digest in the relay
func (k Keeper) setRelayGenesis(ctx sdk.Context, relayGenesis types.Hash256Digest) {
	k.setDigestByStoreKey(ctx, types.RelayGenesisStorage, relayGenesis)
}

// SetGenesisState sets the genesis state
func (k Keeper) SetGenesisState(ctx sdk.Context, genesis, epochStart btcspv.BitcoinHeader) sdk.Error {
	if k.hasRelayGenesis(ctx) {
		return types.ErrAlreadyInit(types.DefaultCodespace)
	}

	k.ingestHeader(ctx, genesis)
	k.ingestHeader(ctx, epochStart)

	k.setRelayGenesis(ctx, genesis.HashLE)
	k.setBestKnownDigest(ctx, genesis.HashLE)
	k.setLastReorgLCA(ctx, genesis.HashLE)

	return nil
}
