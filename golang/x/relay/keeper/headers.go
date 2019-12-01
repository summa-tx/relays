package keeper

import (
	"math/big"

	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) getHeaderStore(ctx sdk.Context) sdk.KVStore {
	return k.getPrefixStore(ctx, types.HeaderStorePrefix)
}

func (k Keeper) emitExtension(ctx sdk.Context, first, last types.BitcoinHeader) {
	ctx.EventManager().EmitEvent(types.NewExtensionEvent(first, last))
}

// HasHeader checks if a header is in the store
func (k Keeper) HasHeader(ctx sdk.Context, digestLE types.Hash256Digest) bool {
	return k.getHeaderStore(ctx).Has(digestLE[:])
}

// GetHeader retrieves a header from the store using its LE diges
func (k Keeper) GetHeader(ctx sdk.Context, digestLE types.Hash256Digest) types.BitcoinHeader {
	var header types.BitcoinHeader
	store := k.getHeaderStore(ctx)

	buf := store.Get(digestLE[:])
	k.cdc.MustUnmarshalBinaryBare(buf, &header)

	return header
}

// compareTargets compares Bitcoin truncated and full-length targets
func compareTargets(full, truncated sdk.Uint) bool {
	// dirty hacks. sdk.Uint doesn't give us easy access to the underlying
	// will be fixed in future sdk version
	f, _ := full.MarshalAmino()
	t, _ := truncated.MarshalAmino()
	fullBI := new(big.Int)
	fullBI.SetString(f, 0)
	truncatedBI := new(big.Int)
	truncatedBI.SetString(t, 0)

	res := new(big.Int)
	res.And(fullBI, truncatedBI)

	return truncatedBI.Cmp(res) == 0
}

func (k Keeper) ingestHeader(ctx sdk.Context, header types.BitcoinHeader) {
	store := k.getHeaderStore(ctx)

	buf := k.cdc.MustMarshalBinaryBare(header)
	store.Set(header.HashLE[:], buf)
}

func validateHeaderChain(anchor types.BitcoinHeader, headers []types.BitcoinHeader, internal, isMainnet bool) sdk.Error {
	prev := anchor // scratchpad, we change this later

	// On internal call, use the header chain target
	expectedTarget := btcspv.ExtractTarget(anchor.Raw)
	if internal {
		expectedTarget = btcspv.ExtractTarget(headers[0].Raw)
	}
	// allocate memory for raw anchor + all headers
	raw := make([]byte, 80*(len(headers)+1))
	copy(raw[0:80], anchor.Raw[:])

	// Make the raw chain
	for i, header := range headers {
		_, err := header.Validate()
		if err != nil {
			return types.FromBTCSPVError(types.DefaultCodespace, err)
		}

		// ensure height changes as expected
		if prev.Height != header.Height-1 {
			return types.ErrBadHeight(types.DefaultCodespace)
		}

		// ensure expectedTarget doesn't change
		// it's allowed to change if the relay is in testnet mode
		if isMainnet && !btcspv.ExtractTarget(header.Raw).Equal(expectedTarget) {
			return types.ErrUnexpectedRetarget(types.DefaultCodespace)
		}

		// copy header raw into a bytearray
		offset := 80 * (i + 1)
		copy(raw[offset:offset+80], header.Raw[:])
		prev = header
	}

	// Then validate the chain
	_, err := btcspv.ValidateHeaderChain(raw)
	if err != nil {
		return types.FromBTCSPVError(types.DefaultCodespace, err)
	}

	return nil
}

func (k Keeper) ingestHeaders(ctx sdk.Context, headers []types.BitcoinHeader, internal bool) sdk.Error {
	if !k.HasHeader(ctx, headers[0].PrevHashLE) {
		return types.ErrUnknownBlock(types.DefaultCodespace)
	}

	anchor := k.GetHeader(ctx, headers[0].PrevHashLE)

	err := validateHeaderChain(anchor, headers, internal, k.IsMainNet)
	if err != nil {
		return err
	}

	for _, header := range headers {
		k.setLink(ctx, header)
		k.ingestHeader(ctx, header)
	}
	return nil
}

func validateDifficultyChange(headers []types.BitcoinHeader, prevEpochStart, anchor types.BitcoinHeader) sdk.Error {
	if anchor.Height%2016 != 2015 {
		return types.ErrWrongEnd(types.DefaultCodespace)
	}
	if anchor.Height != prevEpochStart.Height+2015 || anchor.Height < prevEpochStart.Height {
		return types.ErrWrongStart(types.DefaultCodespace)
	}
	if !btcspv.ExtractDifficulty(anchor.Raw).Equal(btcspv.ExtractDifficulty(prevEpochStart.Raw)) {
		return types.ErrPeriodMismatch(types.DefaultCodespace)
	}

	// calculated target
	expectedTarget := btcspv.RetargetAlgorithm(
		btcspv.ExtractTarget(prevEpochStart.Raw),
		btcspv.ExtractTimestamp(prevEpochStart.Raw),
		btcspv.ExtractTimestamp(anchor.Raw))

	// Observed target in the new period start header
	actualTarget := btcspv.ExtractTarget(headers[0].Raw)

	if !compareTargets(expectedTarget, actualTarget) {
		return types.ErrBadRetarget(types.DefaultCodespace)
	}
	return nil
}

func (k Keeper) ingestDifficultyChange(ctx sdk.Context, prevEpochStartLE types.Hash256Digest, headers []types.BitcoinHeader) sdk.Error {
	if !k.HasHeader(ctx, prevEpochStartLE) {
		return types.ErrUnknownBlock(types.DefaultCodespace)
	}
	if !k.HasHeader(ctx, headers[0].PrevHashLE) {
		return types.ErrUnknownBlock(types.DefaultCodespace)
	}

	// Find the anchor in our store
	prevEpochStart := k.GetHeader(ctx, prevEpochStartLE)
	anchor := k.GetHeader(ctx, headers[0].PrevHashLE)

	err := validateDifficultyChange(headers, prevEpochStart, anchor)
	if err != nil {
		return err
	}

	return k.ingestHeaders(ctx, headers, true)
}

// IngestHeaderChain ingests a chain of headers
func (k Keeper) IngestHeaderChain(ctx sdk.Context, headers []types.BitcoinHeader) sdk.Error {
	// Find the anchor in our store
	return k.ingestHeaders(ctx, headers, false)
}

// IngestDifficultyChange ingests a chain of headers
func (k Keeper) IngestDifficultyChange(ctx sdk.Context, prevEpochStartLE types.Hash256Digest, headers []types.BitcoinHeader) sdk.Error {
	return k.ingestDifficultyChange(ctx, prevEpochStartLE, headers)
}
