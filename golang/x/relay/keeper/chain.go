package keeper

import (
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
)

func (k Keeper) getChainStore(ctx sdk.Context) sdk.KVStore {
	return k.getPrefixStore(ctx, types.ChainStorePrefix)
}

func (k Keeper) emitReorg(ctx sdk.Context, prev, new, lca types.BitcoinHeader) {
	ctx.EventManager().EmitEvent(types.NewReorgEvent(prev, new, lca))
}

// SetGenesisState sets the genesis state
func (k Keeper) SetGenesisState(ctx sdk.Context, genesis, epochStart btcspv.BitcoinHeader) sdk.Error {
	store := k.getChainStore(ctx)

	if store.Has([]byte(types.RelayGenesisStorage)) {
		return types.ErrAlreadyInit(types.DefaultCodespace)
	}

	k.ingestHeader(ctx, genesis)
	k.ingestHeader(ctx, epochStart)

	k.setRelayGenesis(ctx, genesis.HashLE)
	k.setBestKnownDigest(ctx, genesis.HashLE)

	return nil

}

// GetRelayGenesis returns the first digest in the relay
func (k Keeper) GetRelayGenesis(ctx sdk.Context) (types.Hash256Digest, sdk.Error) {
	store := k.getChainStore(ctx)
	result := store.Get([]byte(types.RelayGenesisStorage))

	digest, err := btcspv.NewHash256Digest(result)
	if err != nil {
		return digest, types.ErrBadHash256Digest(types.DefaultCodespace)
	}
	return digest, nil
}

// setRelayGenesis sets the first digest in the relay
func (k Keeper) setRelayGenesis(ctx sdk.Context, relayGenesis types.Hash256Digest) {
	// TODO: Remove this in favor of Genesis state
	store := k.getChainStore(ctx)
	store.Set([]byte(types.RelayGenesisStorage), relayGenesis[:])
}

// setBestKnownDigest sets the first digest in the relay
func (k Keeper) setBestKnownDigest(ctx sdk.Context, bestKnown types.Hash256Digest) {
	// TODO: Remove this in favor of Genesis state
	store := k.getChainStore(ctx)
	store.Set([]byte(types.BestKnownDigestStorage), bestKnown[:])
}

// GetBestKnownDigest returns the best known digest in the relay
func (k Keeper) GetBestKnownDigest(ctx sdk.Context, digest types.Hash256Digest) (types.Hash256Digest, sdk.Error) {
	store := k.getChainStore(ctx)
	result := store.Get([]byte(types.BestKnownDigestStorage))

	digest, err := btcspv.NewHash256Digest(result)
	if err != nil {
		return digest, types.ErrBadHash256Digest(types.DefaultCodespace)
	}
	return digest, nil
}
