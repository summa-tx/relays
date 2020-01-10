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
	return bestKnownHeader.Height - header.Height, nil
}

// TODO: Add errors
func (k Keeper) validateProof(ctx sdk.Context, proof types.SPVProof, requestID uint64) (bool, sdk.Error) {
	valid, err := proof.Validate()
	if err != nil {
		return false, types.FromBTCSPVError(types.DefaultCodespace, err)
	}
	if !valid {
		return false, nil
	}

	lca, lcaErr := k.GetLastReorgLCA(ctx)
	if lcaErr != nil {
		return false, lcaErr
	}
	isAncestor := k.IsAncestor(ctx, proof.ConfirmingHeader.HashLE, lca, 240)
	if !isAncestor {
		return false, nil
	}

	request, _ := k.getRequest(ctx, requestID)
	// TODO: fix these errors
	// if err != nil {
	// 	return false, err
	// }
	confs, _ := k.getConfs(ctx, proof.ConfirmingHeader)
	// if err != nil {
	// 	return false, err
	// }
	if confs < uint32(request.NumConfs) {
		return false, nil
	}

	return true, nil
}
