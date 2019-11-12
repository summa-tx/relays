package keeper

import (
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) getHeaderStore(ctx sdk.Context) sdk.KVStore {
	return k.getPrefixStore(ctx, types.HeaderStorePrefix)
}

func (k Keeper) ingestHeader(ctx sdk.Context, header types.BitcoinHeader) {
	// store := k.getHeaderStore(ctx)

	// TODO: store the header blob as serialized something
}

// IngestHeaderChain ingests headers
func (k Keeper) IngestHeaderChain(ctx sdk.Context, headers []types.BitcoinHeader) sdk.Error {

	raw := []byte{}
	for _, header := range headers {
		_, err := header.Validate()
		if err != nil {
			return types.FromBTCSPVError(types.DefaultCodespace, err)
		}
		raw = append(raw, header.Raw[:]...)
	}

	_, err := btcspv.ValidateHeaderChain(raw)
	if err != nil {
		return types.FromBTCSPVError(types.DefaultCodespace, err)
	}

	for _, header := range headers {
		k.setLink(ctx, header)
		k.ingestHeader(ctx, header)
	}

	return nil
}
