package keeper

import (
	"github.com/summa-tx/relays/golang/x/relay/types"

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
	digest := btcspv.Hash256(header)
	parent := btcspv.ExtractPrevBlockHashLE(header)
	store := ctx.KVStore(k.storeKey)
	store.Set(digest[:], parent[:])
}

// SetLink sets a header parent
func (k Keeper) SetLink(ctx sdk.Context, header []byte) {
	// TODO: Remove this in favor of fully validating add headers
	k.setLink(ctx, header)
}

// GetLink gets headers links
func (k Keeper) GetLink(ctx sdk.Context, digest types.Hash256Digest) (types.Hash256Digest, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	result := store.Get(digest[:])
	digest, err := types.NewHash256Digest(result)
	if err != nil {
		return digest, types.ErrBadHash256Digest(types.DefaultCodespace)
	}
	return digest, nil
}

// GetRelayGenesis returns the first digest in the relay
func (k Keeper) GetRelayGenesis(ctx sdk.Context) (types.Hash256Digest, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	result := store.Get([]byte(types.RelayGenesisStorage))

	digest, err := types.NewHash256Digest(result)
	if err != nil {
		return digest, types.ErrBadHash256Digest(types.DefaultCodespace)
	}
	return digest, nil
}

// SetRelayGenesis sets the first digest in the relay
func (k Keeper) SetRelayGenesis(ctx sdk.Context, relayGenesis types.Hash256Digest) {
	// TODO: Remove this in favor of Genesis state
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.RelayGenesisStorage), relayGenesis[:])
}

// GetBestKnownDigest returns the best known digest in the relay
func (k Keeper) GetBestKnownDigest(ctx sdk.Context, digest types.Hash256Digest) (types.Hash256Digest, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	result := store.Get([]byte(types.BestKnownDigestStorage))

	digest, err := types.NewHash256Digest(result)
	if err != nil {
		return digest, types.ErrBadHash256Digest(types.DefaultCodespace)
	}
	return digest, nil
}

// func (k Keeper) SetBestKnownDigest(ctx sdk.Context, bestKnownDigest types.Hash256Digest, digest types.Hash256Digest) {
// 	link := k.GetLink(ctx, digest)
// 	link.BestKnownDigest = bestKnownDigest
// 	k.SetLink(ctx, bestKnownDigest)
// }
//
// func (k Keeper) GetLastReorgCommonAncestor(ctx sdk.Context, digest types.Hash256Digest) types.Hash256Digest {
// 	return k.GetLink(ctx, digest).LastReorgCommonAncestor
// }
//
// func (k Keeper) SetLastReorgCommonAncestor(ctx sdk.Context, lastReorgCommonAncestor types.Hash256Digest, digest types.Hash256Digest) {
// 	link := k.GetLink(ctx, digest)
// 	link.LastReorgCommonAncestor = lastReorgCommonAncestor
// 	k.SetLink(ctx, lastReorgCommonAncestor)
// }
//
// func (k Keeper) GetPreviousBlock(ctx sdk.Context, digest types.Hash256Digest) types.Hash256Digest {
// 	return k.GetLink(ctx, digest).PreviousBlock
// }
//
// func (k Keeper) SetPreviousBlock(ctx sdk.Context, previousBlock types.Hash256Digest, digest types.Hash256Digest) {
// 	link := k.GetLink(ctx, digest)
// 	link.PreviousBlock = previousBlock
// 	k.SetLink(ctx, previousBlock)
// }
//
// func (k Keeper) GetBlockHeight(ctx sdk.Context, digest types.Hash256Digest) types.Hash256Digest {
// 	return k.GetLink(ctx, digest).BlockHeight
// }
//
// func (k Keeper) SetBlockHeight(ctx sdk.Context, blockHeight types.Hash256Digest, digest types.Hash256Digest) {
// 	link := k.GetLink(ctx, digest)
// 	link.BlockHeight = blockHeight
// 	k.SetLink(ctx, blockHeight)
// }
