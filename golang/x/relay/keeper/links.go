package keeper

import (
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) getLinkStore(ctx sdk.Context) sdk.KVStore {
	return k.getPrefixStore(ctx, types.LinkStorePrefix)
}

func (k Keeper) setLink(ctx sdk.Context, header types.BitcoinHeader) {
	store := k.getLinkStore(ctx)
	store.Set(header.HashLE[:], header.PrevHash[:])
}
