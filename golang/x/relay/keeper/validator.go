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

func (k Keeper) validateProof(ctx sdk.Context, proof types.SPVProof) sdk.Error {
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

	return nil
}

func (k Keeper) checkRequestsFilled(ctx sdk.Context, filledRequests types.FilledRequests) sdk.Error {
	// Validate Proof once
	err := k.validateProof(ctx, filledRequests.Proof)
	if err != nil {
		return err
	}

	confs, confsErr := k.getConfs(ctx, filledRequests.Proof.ConfirmingHeader)
	if confsErr != nil {
		return confsErr
	}

	for i := range filledRequests.Filled {
		// get request
		request, getErr := k.getRequest(ctx, filledRequests.Filled[i].ID)
		if getErr != nil {
			return getErr
		}
		// check confirmations
		if confs < uint32(request.NumConfs) {
			return types.ErrNotEnoughConfs(types.DefaultCodespace)
		}

		// check request
		err := k.checkRequests(
			ctx,
			filledRequests.Filled[i].InputIndex,
			filledRequests.Filled[i].OutputIndex,
			filledRequests.Proof.Vin,
			filledRequests.Proof.Vout,
			filledRequests.Filled[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}
