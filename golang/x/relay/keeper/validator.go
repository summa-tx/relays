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

func (k Keeper) validateProof(ctx sdk.Context, proof types.SPVProof, requestID types.RequestID) sdk.Error {
	// If it is not valid, it will return an error
	_, err := proof.Validate()
	if err != nil {
		return types.FromBTCSPVError(types.DefaultCodespace, err)
	}

	lca, lcaErr := k.GetLastReorgLCA(ctx)
	if lcaErr != nil {
		return lcaErr
	}
	isAncestor := k.IsAncestor(ctx, proof.ConfirmingHeader.HashLE, lca, 240)
	if !isAncestor {
		return types.ErrNotAncestor(types.DefaultCodespace)
	}

	request, getErr := k.getRequest(ctx, requestID)
	if getErr != nil {
		return getErr
	}
	confs, confsErr := k.getConfs(ctx, proof.ConfirmingHeader)
	if confsErr != nil {
		return confsErr
	}
	if confs < uint32(request.NumConfs) {
		return types.ErrNotEnoughConfs(types.DefaultCodespace)
	}

	return nil
}

func (k Keeper) checkRequestsFilled(ctx sdk.Context, r types.FilledRequests) sdk.Error {
	// Validate Proof once
	err := k.validateProof(ctx, r.Proof, r.Requests[0].ID)
	if err != nil {
		return err
	}

	// get confs
	confs, confsErr := k.getConfs(ctx, r.Proof.ConfirmingHeader)
	if confsErr != nil {
		return confsErr
	}

	for i := range r.Requests {
		// check request
		err := k.checkRequests(
			ctx,
			r.Requests[i].InputIndex,
			r.Requests[i].OutputIndex,
			r.Proof.Vin,
			r.Proof.Vout,
			r.Requests[i].ID)
		if err != nil {
			return err
		}

		if i >= 1 {
			// get request
			request, getErr := k.getRequest(ctx, r.Requests[i].ID)
			if getErr != nil {
				return getErr
			}
			// check confirmations
			if confs < uint32(request.NumConfs) {
				return types.ErrNotEnoughConfs(types.DefaultCodespace)
			}
		}
	}
	return nil
}
