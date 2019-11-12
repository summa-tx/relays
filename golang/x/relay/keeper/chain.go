package keeper

import (
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
)

// GetRelayGenesis returns the first digest in the relay
func (k Keeper) GetRelayGenesis(ctx sdk.Context) (types.Hash256Digest, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	result := store.Get([]byte(types.RelayGenesisStorage))

	digest, err := btcspv.NewHash256Digest(result)
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

	digest, err := btcspv.NewHash256Digest(result)
	if err != nil {
		return digest, types.ErrBadHash256Digest(types.DefaultCodespace)
	}
	return digest, nil
}
