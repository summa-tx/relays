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

func (k Keeper) getDigestByStoreKey(ctx sdk.Context, key string) (types.Hash256Digest, sdk.Error) {
	store := k.getChainStore(ctx)
	result := store.Get([]byte(key))

	digest, err := btcspv.NewHash256Digest(result)
	if err != nil {
		return types.Hash256Digest{}, types.ErrBadHash256Digest(types.DefaultCodespace)
	}
	return digest, nil
}

func (k Keeper) setDigestByStoreKey(ctx sdk.Context, key string, digest types.Hash256Digest) {
	// TODO: Remove this in favor of Genesis state
	store := k.getChainStore(ctx)
	store.Set([]byte(key), digest[:])
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
	k.setLastReorgLCA(ctx, genesis.HashLE)

	return nil
}

// setRelayGenesis sets the first digest in the relay
func (k Keeper) setRelayGenesis(ctx sdk.Context, relayGenesis types.Hash256Digest) {
	k.setDigestByStoreKey(ctx, types.RelayGenesisStorage, relayGenesis)
}

// GetRelayGenesis returns the first digest in the relay
func (k Keeper) GetRelayGenesis(ctx sdk.Context) (types.Hash256Digest, sdk.Error) {
	return k.getDigestByStoreKey(ctx, types.RelayGenesisStorage)
}

// setBestKnownDigest sets the best known chain tip
func (k Keeper) setBestKnownDigest(ctx sdk.Context, bestKnown types.Hash256Digest) {
	k.setDigestByStoreKey(ctx, types.BestKnownDigestStorage, bestKnown)
}

// GetBestKnownDigest returns the best known digest in the relay
func (k Keeper) GetBestKnownDigest(ctx sdk.Context) (types.Hash256Digest, sdk.Error) {
	return k.getDigestByStoreKey(ctx, types.BestKnownDigestStorage)
}

// setLastReorgLCA sets the latest common ancestor of the last reorg
func (k Keeper) setLastReorgLCA(ctx sdk.Context, bestKnown types.Hash256Digest) {
	k.setDigestByStoreKey(ctx, types.LastReorgLCAStorage, bestKnown)
}

// GetLastReorgLCA returns the best known digest in the relay
func (k Keeper) GetLastReorgLCA(ctx sdk.Context) (types.Hash256Digest, sdk.Error) {
	return k.getDigestByStoreKey(ctx, types.LastReorgLCAStorage)
}
