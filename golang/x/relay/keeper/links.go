package keeper

import (
	btcspv "github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) getLinkStore(ctx sdk.Context) sdk.KVStore {
	return k.getPrefixStore(ctx, types.LinkStorePrefix)
}

func (k Keeper) setLink(ctx sdk.Context, header types.BitcoinHeader) {
	store := k.getLinkStore(ctx)
	store.Set(header.HashLE[:], header.PrevHashLE[:])
}

func (k Keeper) getLink(ctx sdk.Context, digestLE types.Hash256Digest) types.Hash256Digest {
	store := k.getLinkStore(ctx)
	buf := store.Get(digestLE[:])
	// Can only fail if data store is corrupt
	parentHashLE, _ := btcspv.NewHash256Digest(buf)
	return parentHashLE
}
