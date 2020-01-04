package keeper

import (
	btcspv "github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) getRequestStore(ctx sdk.Context) sdk.KVStore {
	return k.getPrefixStore(ctx, types.RequestStorePrefix)
}

func (k Keeper) hasRequest(ctx sdk.Context, digestLE types.Hash256Digest) bool {
	store := k.getRequestStore(ctx)
	return store.Has(digestLE[:])
}

func (k Keeper) setRequest(ctx sdk.Context, header types.BitcoinHeader) {
	store := k.getRequestStore(ctx)
	store.Set(header.HashLE[:], header.PrevHashLE[:])
}

func (k Keeper) getRequest(ctx sdk.Context, digestLE types.Hash256Digest) types.Hash256Digest {
	store := k.getRequestStore(ctx)
	buf := store.Get(digestLE[:])
	// Can only fail if data store is corrupt
	parentHashLE, _ := btcspv.NewHash256Digest(buf)
	return parentHashLE
}
