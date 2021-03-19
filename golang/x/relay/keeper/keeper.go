package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"

	"github.com/summa-tx/relays/golang/x/relay/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey     sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc          *codec.LegacyAmino // The wire codec for binary encoding/decoding.
	IsMainNet    bool
	ProofHandler types.ProofHandler
}

// NewKeeper instantiates a new keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.LegacyAmino, mainnet bool, handler types.ProofHandler) Keeper {
	return Keeper{
		storeKey:     storeKey,
		cdc:          cdc,
		IsMainNet:    mainnet,
		ProofHandler: handler,
	}
}

func (k Keeper) getPrefixStore(ctx sdk.Context, namespace string) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(namespace))
}

func (k Keeper) hasRelayGenesis(ctx sdk.Context) bool {
	store := k.getChainStore(ctx)
	return store.Has([]byte(types.RelayGenesisStorage))
}

// setRelayGenesis sets the first digest in the relay
func (k Keeper) setRelayGenesis(ctx sdk.Context, relayGenesis types.Hash256Digest) {
	k.setDigestByStoreKey(ctx, types.RelayGenesisStorage, relayGenesis)
}

// GetRelayGenesis returns the first digest in the relay
func (k Keeper) GetRelayGenesis(ctx sdk.Context) (types.Hash256Digest, *sdkerrors.Error) {
	return k.getDigestByStoreKey(ctx, types.RelayGenesisStorage)
}

// SetGenesisState sets the genesis state
func (k Keeper) SetGenesisState(ctx sdk.Context, genesis, epochStart btcspv.BitcoinHeader) *sdkerrors.Error {
	if k.hasRelayGenesis(ctx) {
		return types.ErrAlreadyInit(types.DefaultCodespace)
	}

	k.ingestHeader(ctx, genesis)
	k.ingestHeader(ctx, epochStart)

	k.setRelayGenesis(ctx, genesis.Hash)
	k.setBestKnownDigest(ctx, genesis.Hash)
	k.setLastReorgLCA(ctx, genesis.Hash)

	// this will only fail if the genesis state is corrupt
	_ = k.setCurrentEpochDifficulty(ctx, btcspv.ExtractDifficulty(genesis.Raw))

	return nil
}
