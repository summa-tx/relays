package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

func (k Keeper) getConfs(ctx sdk.Context, header types.BitcoinHeader) (uint32, sdk.Error) {
	bestKnown, err := k.GetBestKnownDigest(ctx)
	if err != nil {
		return 0, err
	}
	bestKnownHeader, err := k.GetHeader(ctx, bestKnown)
	if err != nil {
		return 0, err
	}
	return bestKnownHeader.Height - header.Height, nil
}

func (k Keeper) validateProof(ctx sdk.Context, proof types.SPVProof, requestID uint64) (bool, sdk.Error) {
	// If it is not valid, it will return an error
	_, err := proof.Validate()
	if err != nil {
		return false, types.FromBTCSPVError(types.DefaultCodespace, err)
	}

	lca, lcaErr := k.GetLastReorgLCA(ctx)
	if lcaErr != nil {
		return false, lcaErr
	}
	isAncestor := k.IsAncestor(ctx, proof.ConfirmingHeader.HashLE, lca, 240)
	if !isAncestor {
		return false, types.ErrNotAncestor(types.DefaultCodespace)
	}

	request, getErr := k.getRequest(ctx, requestID)
	if getErr != nil {
		return false, getErr
	}
	confs, confsErr := k.getConfs(ctx, proof.ConfirmingHeader)
	if confsErr != nil {
		return false, confsErr
	}
	if confs < uint32(request.NumConfs) {
		return false, types.ErrNotEnoughConfs(types.DefaultCodespace)
	}

	return true, nil
}

func (k Keeper) checkAll(ctx sdk.Context, proof types.SPVProof, inputIndex, outputIndex uint8, vin, vout []byte, requestIDs []uint64) ([]uint64, sdk.Error) {
	var matching []uint64
	for i := range requestIDs {
		// Validate Proof
		validProof, err := k.validateProof(ctx, proof, requestIDs[i])
		if err != nil {
			return []uint64{}, err
		}
		// Check request
		matchingRequest, err := k.checkRequests(ctx, inputIndex, outputIndex, vin, vout, requestIDs[i])
		if err != nil {
			return []uint64{}, err
		}
		if validProof && matchingRequest {
			matching = append(matching, requestIDs[i])
		}
	}
	return matching, nil
}

func (k Keeper) checkRequestsFilled(ctx sdk.Context, r types.FilledRequests) (bool, sdk.Error) {
	// Validate Proof once.  If proof errors it is not valid.
	_, err := k.validateProof(ctx, r.Proof, r.Requests[0].ID)
	if err != nil {
		return false, err
	}

	// get confs
	confs, confsErr := k.getConfs(ctx, r.Proof.ConfirmingHeader)
	if confsErr != nil {
		return false, confsErr
	}

	// check subsequent proofs
	for i := 1; i < len(r.Requests); i++ {
		// get request
		request, getErr := k.getRequest(ctx, r.Requests[i].ID)
		if getErr != nil {
			return false, getErr
		}
		// check confirmations
		if confs < uint32(request.NumConfs) {
			return false, types.ErrNotEnoughConfs(types.DefaultCodespace)
		}
	}

	for i := range r.Requests {
		// check request
		matchingRequest, err := k.checkRequests(
			ctx,
			r.Requests[i].InputIndex,
			r.Requests[i].OutputIndex,
			r.Proof.Vin,
			r.Proof.Vout,
			r.Requests[i].ID)
		if err != nil {
			return false, err
		}
		if !matchingRequest {
			return false, types.ErrNotMatching(types.DefaultCodespace)
		}
	}
	return true, nil
}
