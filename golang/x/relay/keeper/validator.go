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
