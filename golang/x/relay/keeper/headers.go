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
	store := k.getHeaderStore(ctx)

	buf := k.cdc.MustMarshalBinaryBare(header)
	store.Set(header.HashLE[:], buf)
}

// GetHeader retrieves a header from the store using its LE diges
func (k Keeper) GetHeader(ctx sdk.Context, digestLE types.Hash256Digest) types.BitcoinHeader {
	var header types.BitcoinHeader
	store := k.getHeaderStore(ctx)

	buf := store.Get(digestLE[:])
	k.cdc.MustUnmarshalBinaryBare(buf, &header)

	return header
}

func (k Keeper) ingestHeaders(ctx sdk.Context, anchor types.BitcoinHeader, headers []types.BitcoinHeader) sdk.Error {
	target := btcspv.ExtractTarget(anchor.Raw)
	// TODO: validation
	// // Make the raw chain
	// raw := []byte{}
	// for i, header := range headers {
	// 	_, err := header.Validate()
	// 	if err != nil {
	// 		return types.FromBTCSPVError(types.DefaultCodespace, err)
	// 	}
	// 	/// Check that height is consistent throughout
	// 	if i > 0 {
	// 		if headers[i-1].Height != header.Height-1 {
	// 			return types.ErrBadHeight(types.DefaultCodespace)
	// 		}
	// 		if btcspv.ExtractTarget(header.Raw) != target {
	// 			return types.ErrUnexptectedRetarget(types.DefaultCodespace)
	// 		}
	// 	}
	// 	raw = append(raw, header.Raw[:]...)
	// }
	//
	// // Then validate the chain
	// _, err := btcspv.ValidateHeaderChain(raw)
	// if err != nil {
	// 	return types.FromBTCSPVError(types.DefaultCodespace, err)
	// }
	//
	// for _, header := range headers {
	// 	k.setLink(ctx, header)
	// 	k.ingestHeader(ctx, header)
	// }
	return nil
}

func (k Keeper) ingestDiffChange(ctx sdk.Context, anchor types.BitcoinHeader, headers []types.BitcoinHeader) sdk.Error {
	// TODO
	return nil
}

// IngestHeaderChain ingests a chain of headers
func (k Keeper) IngestHeaderChain(ctx sdk.Context, headers []types.BitcoinHeader) sdk.Error {

	// Find the anchor in our store
	anchor := k.GetHeader(ctx, headers[0].PrevHashLE)

	if anchor.Height%2016 == 0 {
		return k.ingestDiffChange(ctx, anchor, headers)
	}
	return k.ingestHeaders(ctx, anchor, headers)
}
