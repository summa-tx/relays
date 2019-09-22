package keeper

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSetLink:
			return handleMsgSetLink(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleMsgSetLink(ctx sdk.Context, keeper Keeper, msg types.MsgSetLink) sdk.Result {
	b, _ := hex.DecodeString(msg.Header)
	keeper.SetLink(ctx, b) // If so, set the name to the value specified in the msg.//s
	return sdk.Result{}    // return
}
