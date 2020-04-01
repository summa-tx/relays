package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

func (s *KeeperSuite) TestEmitProofRequest() {
	s.Keeper.emitProofRequest(s.Context, []byte{0}, []byte{0}, 0, types.RequestID{}, types.Local)

	events := s.Context.EventManager().Events()
	e := events[0]
	s.Equal("proof_request", e.Type)
}

// tests getNextID and incrementID
func (s *KeeperSuite) TestIncrementID() {
	id, err := s.Keeper.getNextID(s.Context)
	s.SDKNil(err)
	s.Equal(types.RequestID{}, id)

	err = s.Keeper.incrementID(s.Context)
	s.SDKNil(err)

	id, err = s.Keeper.getNextID(s.Context)
	s.SDKNil(err)
	s.Equal(types.RequestID{0, 0, 0, 0, 0, 0, 0, 1}, id)

	// errors if it cannot get next ID
	store := s.Keeper.getRequestStore(s.Context)
	idTag := []byte(types.RequestIDTag)
	store.Set(idTag, bytes.Repeat([]byte{9}, 9))

	err = s.Keeper.incrementID(s.Context)
	s.Equal(sdk.CodeType(107), err.Code())
}

func (s *KeeperSuite) TestHasRequest() {
	hasRequest := s.Keeper.hasRequest(s.Context, types.RequestID{})
	s.Equal(false, hasRequest)
	requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 4, types.Local, nil)
	s.Nil(requestErr)
	hasRequest = s.Keeper.hasRequest(s.Context, types.RequestID{})
	s.Equal(true, hasRequest)
}

func (s *KeeperSuite) TestSetRequest() {
	store := s.Keeper.getRequestStore(s.Context)
	idTag := []byte(types.RequestIDTag)
	store.Set(idTag, bytes.Repeat([]byte{9}, 9))

	err := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 0, types.Local, nil)
	s.Equal(sdk.CodeType(107), err.Code())
}

func (s *KeeperSuite) TestSetRequestState() {
	// errors if request is not found
	activeErr := s.Keeper.setRequestState(s.Context, types.RequestID{}, false)
	s.Equal(sdk.CodeType(601), activeErr.Code())

	// set request
	requestErr := s.Keeper.setRequest(s.Context, []byte{1}, []byte{1}, 0, 0, types.Local, nil)
	s.Nil(requestErr)
	// change active state to false
	activeErr = s.Keeper.setRequestState(s.Context, types.RequestID{}, false)
	s.Nil(activeErr)

	deactivatedRequest, deactivatedRequestErr := s.Keeper.getRequest(s.Context, types.RequestID{})
	s.Nil(deactivatedRequestErr)
	s.Equal(false, deactivatedRequest.ActiveState)
}

func (s *KeeperSuite) TestGetRequest() {
	requestRes := s.Fixtures.RequestTestCases.EmptyRequest
	request, err := s.Keeper.getRequest(s.Context, types.RequestID{})
	s.Equal(sdk.CodeType(601), err.Code())
	s.Equal(types.ProofRequest{}, request)

	requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 0, types.Local, nil)
	s.Nil(requestErr)

	request, err = s.Keeper.getRequest(s.Context, types.RequestID{})
	s.Nil(err)
	s.Equal(requestRes, request)
}

func (s *KeeperSuite) TestCheckRequests() {
	tc := s.Fixtures.RequestTestCases.CheckRequests
	v := tc[0]

	// Errors if request is not found
	err := s.Keeper.checkRequests(
		s.Context,
		v.InputIdx,
		v.OutputIdx,
		v.Vin,
		v.Vout,
		v.RequestID)
	s.Equal(sdk.CodeType(601), err.Code())

	// set request
	requestErr := s.Keeper.setRequest(s.Context, []byte{1}, []byte{1}, 0, 0, types.Local, nil)
	s.Nil(requestErr)
	// change active state to false
	activeErr := s.Keeper.setRequestState(s.Context, types.RequestID{}, false)
	s.Nil(activeErr)
	// errors if request is not active
	err = s.Keeper.checkRequests(
		s.Context,
		v.InputIdx,
		v.OutputIdx,
		v.Vin,
		v.Vout,
		v.RequestID)
	s.Equal(sdk.CodeType(606), err.Code())

	// change active state to false
	activeErr = s.Keeper.setRequestState(s.Context, types.RequestID{}, true)
	s.Nil(activeErr)
	// errors if request pays is not equal to output
	err = s.Keeper.checkRequests(
		s.Context,
		v.InputIdx,
		v.OutputIdx,
		v.Vin,
		v.Vout,
		v.RequestID)
	s.Equal(sdk.CodeType(607), err.Code())

	// Errors if output value is less than pays value
	out, outErr := btcspv.ExtractOutputAtIndex(v.Vout, uint(v.OutputIdx))
	s.Nil(outErr)
	// out[8:] extracts the output script which we use to set the request
	requestErr = s.Keeper.setRequest(s.Context, []byte{0}, out[8:], 1000, 0, types.Local, nil)
	s.SDKNil(requestErr)
	err = s.Keeper.checkRequests(
		s.Context,
		v.InputIdx,
		v.OutputIdx,
		v.Vin,
		v.Vout,
		types.RequestID{0, 0, 0, 0, 0, 0, 0, 1})
	s.Equal(sdk.CodeType(608), err.Code())

	// Errors if input value does not equal spends value
	requestErr = s.Keeper.setRequest(s.Context, []byte{1}, []byte{}, 0, 255, types.Local, nil)
	s.SDKNil(requestErr)
	err = s.Keeper.checkRequests(
		s.Context,
		v.InputIdx,
		v.OutputIdx,
		v.Vin,
		v.Vout,
		types.RequestID{0, 0, 0, 0, 0, 0, 0, 2})
	s.Equal(sdk.CodeType(609), err.Code())

	// Success
	in, extractErr := btcspv.ExtractInputAtIndex(v.Vin, uint(v.InputIdx))
	s.Nil(extractErr)
	outpoint := btcspv.ExtractOutpoint(in)
	// out[8:] extracts the output script which we use to set the request
	requestErr = s.Keeper.setRequest(s.Context, outpoint, out[8:], 10, 255, types.Local, nil)
	s.SDKNil(requestErr)
	err = s.Keeper.checkRequests(
		s.Context,
		v.InputIdx,
		v.OutputIdx,
		v.Vin,
		v.Vout,
		types.RequestID{0, 0, 0, 0, 0, 0, 0, 3})
	s.SDKNil(err)

	for i := 1; i < len(tc); i++ {
		err := s.Keeper.checkRequests(
			s.Context,
			tc[i].InputIdx,
			tc[i].OutputIdx,
			tc[i].Vin,
			tc[i].Vout,
			tc[i].RequestID)
		if tc[i].Error == 0 {
			s.SDKNil(err)
		} else {
			s.Equal(sdk.CodeType(tc[i].Error), err.Code())
		}
	}
}
