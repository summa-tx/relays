package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

// NewHandler returns a handler for relay type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgIngestHeaderChain:
			return handleMsgIngestHeaderChain(ctx, keeper, msg)
		case types.MsgIngestDifficultyChange:
			return handleMsgIngestDifficultyChange(ctx, keeper, msg)
		case types.MsgMarkNewHeaviest:
			return handleMsgMarkNewHeaviest(ctx, keeper, msg)
		case types.MsgNewRequest:
			return handleMsgNewRequest(ctx, keeper, msg)
		case types.MsgProvideProof:
			return handleMsgProvideProof(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized relay Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgIngestHeaderChain(ctx sdk.Context, keeper Keeper, msg types.MsgIngestHeaderChain) sdk.Result {
	err := keeper.IngestHeaderChain(ctx, msg.Headers)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgIngestDifficultyChange(ctx sdk.Context, keeper Keeper, msg types.MsgIngestDifficultyChange) sdk.Result {
	err := keeper.IngestDifficultyChange(ctx, msg.Start, msg.Headers)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgMarkNewHeaviest(ctx sdk.Context, keeper Keeper, msg types.MsgMarkNewHeaviest) sdk.Result {
	err := keeper.MarkNewHeaviest(ctx, msg.Ancestor, msg.CurrentBest, msg.NewBest, msg.Limit)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgNewRequest(ctx sdk.Context, keeper Keeper, msg types.MsgNewRequest) sdk.Result {
	// Validate message
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	// TODO: Add more complex permissioning
	// Set request
	err = keeper.setRequest(ctx, msg.Spends, msg.Pays, msg.PaysValue, msg.NumConfs, msg.Origin, msg.Action)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgProvideProof(ctx sdk.Context, keeper Keeper, msg types.MsgProvideProof) sdk.Result {
	err := keeper.checkRequestsFilled(ctx, msg.Filled)
	if err != nil {
		return err.Result()
	}

	// TODO: Add "hooks"
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
