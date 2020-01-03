package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

func validateProof(ctx sdk.Context, keeper Keeper, proof types.SPVProof) (bool, sdk.Error) {
	valid, err := proof.Validate()
	if err != nil {
		return false, types.FromBTCSPVError(types.DefaultCodespace, err)
	}
	if !valid {
		return false, nil
	}

	lca, lcaErr := keeper.GetLastReorgLCA(ctx)
	if lcaErr != nil {
		return false, lcaErr
	}
	isAncestor := keeper.IsAncestor(ctx, proof.ConfirmingHeader.HashLE, lca, 240)
	if !isAncestor {
		return false, nil
	}
	return true, nil
}
