package keeper

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper instantiates a new keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (k Keeper) setLink(ctx sdk.Context, header []byte) {
	println(hex.EncodeToString(header))
	digest := btcspv.Hash256(header)
	println(hex.EncodeToString(digest))
	parent := btcspv.ExtractPrevBlockHashLE(header)
	println(hex.EncodeToString(parent))
	println(k.storeKey)
	store := ctx.KVStore(k.storeKey)
	println("store got")
	store.Set(digest[:], parent[:])
}

// SetLink sets a header parent
func (k Keeper) SetLink(ctx sdk.Context, header []byte) {
	// TODO: Remove this
	k.setLink(ctx, header)
}

// GetLink gets headers links
func (k Keeper) GetLink(ctx sdk.Context, digest [32]byte) []byte {
	println(hex.EncodeToString(digest[:]))
	println(k.storeKey)
	store := ctx.KVStore(k.storeKey)
	println(1)
	return store.Get(digest[:])
}

// func (k Keeper) SetRelayGenesis(ctx sdk.Context, digest [32]byte) []byte {
// 	store := ctx.KVStore(k.storeKey)
// 	store.Set([32]byte())
// }
