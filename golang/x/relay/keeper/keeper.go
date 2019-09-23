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

func (k Keeper) GetRelayGenesis(ctx sdk.Context, digest [32]byte) [32]byte {
	return k.GetLink(ctx, digest).RelayGenesis
}

func (k Keeper) SetRelayGenesis(ctx sdk.Context, relayGenesis [32]byte, digest [32]byte) {
	link := k.GetLink(ctx, digest)
	link.RelayGenesis = relayGenesis
	k.SetLink(ctx, relayGenesis)
}

func (k Keeper) GetBestKnownDigest(ctx sdk.Context, digest [32]byte) [32]byte {
	return k.GetLink(ctx, digest).BestKnownDigest
}

func (k Keeper) SetBestKnownDigest(ctx sdk.Context, bestKnownDigest [32]byte, digest [32]byte) {
	link := k.GetLink(ctx, digest)
	link.BestKnownDigest = bestKnownDigest
	k.SetLink(ctx, bestKnownDigest)
}

func (k Keeper) GetLastReorgCommonAncestor(ctx sdk.Context, digest [32]byte) [32]byte {
	return k.GetLink(ctx, digest).LastReorgCommonAncestor
}

func (k Keeper) SetLastReorgCommonAncestor(ctx sdk.Context, lastReorgCommonAncestor [32]byte, digest [32]byte) {
	link := k.GetLink(ctx, digest)
	link.LastReorgCommonAncestor = lastReorgCommonAncestor
	k.SetLink(ctx, lastReorgCommonAncestor)
}

func (k Keeper) GetPreviousBlock(ctx sdk.Context, digest [32]byte) [32]byte {
	return k.GetLink(ctx, digest).PreviousBlock
}

func (k Keeper) SetPreviousBlock(ctx sdk.Context, previousBlock [32]byte, digest [32]byte) {
	link := k.GetLink(ctx, digest)
	link.PreviousBlock = previousBlock
	k.SetLink(ctx, previousBlock)
}

func (k Keeper) GetBlockHeight(ctx sdk.Context, digest [32]byte) [32]byte {
	return k.GetLink(ctx, digest).BlockHeight
}

func (k Keeper) SetBlockHeight(ctx sdk.Context, blockHeight [32]byte, digest [32]byte) {
	link := k.GetLink(ctx, digest)
	link.BlockHeight = blockHeight
	k.SetLink(ctx, blockHeight)
}
