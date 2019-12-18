package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

func getAccAddress() sdk.AccAddress {
	address, _ := sdk.AccAddressFromBech32("cosmos1ay37rp2pc3kjarg7a322vu3sa8j9puah8msyfw")
	return address
}

// Create a bad sdk.msg to pass into TestNewHandler
type MsgBadMessage struct {
	Signer sdk.AccAddress `json:"signer"`
}

func (msg MsgBadMessage) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
func (msg MsgBadMessage) Type() string             { return "bad_message" }
func (msg MsgBadMessage) ValidateBasic() sdk.Error { return nil }
func (msg MsgBadMessage) GetSignBytes() []byte {
	return sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgBadMessage) Route() string { return types.RouterKey }

func (s *KeeperSuite) TestNewHandler() {
	handler := NewHandler(s.Keeper)

	badMsg := MsgBadMessage{
		Signer: getAccAddress(),
	}

	res := handler(s.Context, badMsg)
	s.Equal(res.Log, "{\"codespace\":\"sdk\",\"code\":6,\"message\":\"Unrecognized relay Msg type: bad_message\"}")
}

func (s *KeeperSuite) TestHandleMsgIngestHeaderChain() {
	testCases := s.Fixtures.HeaderTestCases.ValidateChain
	handler := NewHandler(s.Keeper)

	newMsg := types.NewMsgIngestHeaderChain(getAccAddress(), testCases[0].Headers)

	res := handler(s.Context, newMsg)
	s.Equal(res.Code, sdk.CodeType(103))

	s.Keeper.ingestHeader(s.Context, testCases[0].Anchor)
	res = handler(s.Context, newMsg)
	s.Equal(res.Events[0].Type, "extension")
}

func (s *KeeperSuite) TestHandleMsgIngestDifficultyChange() {
	testCases := s.Fixtures.HeaderTestCases.ValidateDiffChange
	handler := NewHandler(s.Keeper)

	newMsg := types.NewMsgIngestDifficultyChange(getAccAddress(), testCases[0].PrevEpochStart.HashLE, testCases[0].Headers)

	res := handler(s.Context, newMsg)
	s.Equal(res.Code, sdk.CodeType(103))

	s.Keeper.ingestHeader(s.Context, testCases[0].PrevEpochStart)
	s.Keeper.ingestHeader(s.Context, testCases[0].Anchor)
	res = handler(s.Context, newMsg)
	s.Equal(res.Events[0].Type, "extension")
}

func (s *KeeperSuite) TestHandleMsgMarkNewHeaviest() {
	testCases := s.Fixtures.HeaderTestCases.ValidateDiffChange
	handler := NewHandler(s.Keeper)

	s.Keeper.ingestHeader(s.Context, testCases[0].PrevEpochStart)
	s.Keeper.ingestHeader(s.Context, testCases[0].Anchor)
	newMsg := types.NewMsgIngestDifficultyChange(getAccAddress(), testCases[0].PrevEpochStart.HashLE, testCases[0].Headers)
	res := handler(s.Context, newMsg)
	s.Equal(res.Events[0].Type, "extension")
}
