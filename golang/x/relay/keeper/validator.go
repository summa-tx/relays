package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

func validateProof(proof types.SPVProof, keeper Keeper, ctx sdk.Context) (bool, error) {
	valid, err := proof.Validate()
	if err != nil {
		return false, types.FromBTCSPVError(types.DefaultCodespace, err)
	}
	if !valid {
		return false, nil
	}

	lca, err := keeper.GetLastReorgLCA(ctx)
	if err != nil {
		return false, err
	}
	isAncestor := keeper.IsAncestor(ctx, proof.ConfirmingHeader.Hash, lca, 240)
	if !isAncestor {
		return false, nil
	}
	return true, nil
}
