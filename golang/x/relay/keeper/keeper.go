package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
