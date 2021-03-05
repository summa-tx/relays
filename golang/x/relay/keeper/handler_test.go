package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
func (msg MsgBadMessage) ValidateBasic() *sdkerrors.Error { return nil }
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
	s.Equal("{\"codespace\":\"sdk\",\"code\":6,\"message\":\"Unrecognized relay Msg type: bad_message\"}", res.Log)
}

func (s *KeeperSuite) TestHandleMsgIngestHeaderChain() {
	testCases := s.Fixtures.HeaderTestCases.ValidateChain
	handler := NewHandler(s.Keeper)

	newMsg := types.NewMsgIngestHeaderChain(getAccAddress(), testCases[0].Headers)

	res := handler(s.Context, newMsg)
	s.Equal(sdk.CodeType(types.UnknownBlock), res.Code)

	s.Keeper.ingestHeader(s.Context, testCases[0].Anchor)
	res = handler(s.Context, newMsg)
	s.Equal("extension", res.Events[0].Type)
}

func (s *KeeperSuite) TestHandleMsgIngestDifficultyChange() {
	testCases := s.Fixtures.HeaderTestCases.ValidateDiffChange
	handler := NewHandler(s.Keeper)

	newMsg := types.NewMsgIngestDifficultyChange(getAccAddress(), testCases[0].PrevEpochStart.Hash, testCases[0].Headers)

	res := handler(s.Context, newMsg)
	s.Equal(sdk.CodeType(types.UnknownBlock), res.Code)

	s.Keeper.ingestHeader(s.Context, testCases[0].PrevEpochStart)
	s.Keeper.ingestHeader(s.Context, testCases[0].Anchor)
	res = handler(s.Context, newMsg)
	s.Equal("extension", res.Events[0].Type)
}

func (s *KeeperSuite) TestHandleMsgMarkNewHeaviest() {
	testCases := s.Fixtures.HeaderTestCases.ValidateDiffChange
	handler := NewHandler(s.Keeper)

	s.Keeper.ingestHeader(s.Context, testCases[0].PrevEpochStart)
	s.Keeper.ingestHeader(s.Context, testCases[0].Anchor)
	newMsg := types.NewMsgIngestDifficultyChange(getAccAddress(), testCases[0].PrevEpochStart.Hash, testCases[0].Headers)
	res := handler(s.Context, newMsg)
	s.Equal("extension", res.Events[0].Type)
}

func (s *KeeperSuite) TestHandleMarkNewHeaviest() {
	tv := s.Fixtures.ChainTestCases.IsMostRecentCA
	pre := tv.PreRetargetChain
	post := tv.PostRetargetChain
	handler := NewHandler(s.Keeper)

	var postWithOrphan []types.BitcoinHeader
	postWithOrphan = append(postWithOrphan, post[:len(post)-2]...)
	postWithOrphan = append(postWithOrphan, tv.Orphan)

	err := s.Keeper.SetGenesisState(s.Context, tv.Genesis, tv.OldPeriodStart)
	s.SDKNil(err)

	err = s.Keeper.IngestHeaderChain(s.Context, pre)
	s.SDKNil(err)
	err = s.Keeper.IngestDifficultyChange(s.Context, tv.OldPeriodStart.Hash, post)
	s.SDKNil(err)
	err = s.Keeper.IngestDifficultyChange(s.Context, tv.OldPeriodStart.Hash, postWithOrphan)
	s.SDKNil(err)

	// returns correct error
	newMsg := types.NewMsgMarkNewHeaviest(getAccAddress(), tv.OldPeriodStart.Hash, tv.OldPeriodStart.Raw, tv.OldPeriodStart.Raw, 10)
	res := handler(s.Context, newMsg)
	s.Equal(sdk.CodeType(types.NotBestKnown), res.Code)

	// Successfully marks new heaviest
	newMsg = types.NewMsgMarkNewHeaviest(getAccAddress(), tv.Genesis.Hash, tv.Genesis.Raw, pre[0].Raw, 10)
	res = handler(s.Context, newMsg)
	s.Equal("extension", res.Events[0].Type)
}

func (s *KeeperSuite) TestHandleNewRequest() {
	handler := NewHandler(s.Keeper)

	// Success
	newRequest := types.NewMsgNewRequest(getAccAddress(), bytes.Repeat([]byte{0}, 36), []byte{0}, 0, 0, types.Local, nil)
	res := handler(s.Context, newRequest)
	hasRequest := s.Keeper.hasRequest(s.Context, types.RequestID{})
	s.Equal(true, hasRequest)
	s.Equal("proof_request", res.Events[0].Type)

	// Msg validation failed
	newRequest = types.NewMsgNewRequest(getAccAddress(), []byte{0}, []byte{0}, 0, 0, types.Local, nil)
	res = handler(s.Context, newRequest)
	s.Equal(sdk.CodeType(types.SpendsLength), res.Code)

	// setRequest error
	store := s.Keeper.getRequestStore(s.Context)
	store.Set([]byte(types.RequestIDTag), []byte("badID"))

	newRequest = types.NewMsgNewRequest(getAccAddress(), bytes.Repeat([]byte{0}, 36), []byte{0}, 0, 0, types.Local, nil)
	res = handler(s.Context, newRequest)
	s.Equal(sdk.CodeType(types.BadHexLen), res.Code)
}
