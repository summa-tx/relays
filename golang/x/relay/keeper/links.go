package keeper

import (
	btcspv "github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) getLinkStore(ctx sdk.Context) sdk.KVStore {
	return k.getPrefixStore(ctx, types.LinkStorePrefix)
}

func (k Keeper) hasLink(ctx sdk.Context, digestLE types.Hash256Digest) bool {
	store := k.getLinkStore(ctx)
	return store.Has(digestLE[:])
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

// FindAncestor finds the nth ancestor of some digest
func (k Keeper) FindAncestor(ctx sdk.Context, digestLE types.Hash256Digest, offset uint32) (types.Hash256Digest, sdk.Error) {
	current := digestLE
	if !k.hasLink(ctx, current) {
		return types.Hash256Digest{}, types.ErrUnknownBlock(types.DefaultCodespace)
	}

	for i := uint32(0); i < offset; i++ {
		current := k.getLink(ctx, current)
		if !k.hasLink(ctx, current) {
			return types.Hash256Digest{}, types.ErrUnknownBlock(types.DefaultCodespace)
		}
	}

	return current, nil
}

// IsAncestor checks whether
func (k Keeper) IsAncestor(ctx sdk.Context, digestLE, ancestor types.Hash256Digest, limit uint32) bool {
	current := digestLE

	for i := uint32(0); i < limit; i++ {
		if !k.hasLink(ctx, current) {
			return false
		}
		current := k.getLink(ctx, current)
		if current == ancestor {
			return true
		}
	}
	return false
}
