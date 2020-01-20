package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

func (s *KeeperSuite) TestEmitProofRequest() {
	s.Keeper.emitProofRequest(s.Context, []byte{0}, []byte{0}, 0, types.RequestID{})

	events := s.Context.EventManager().Events()
	e := events[0]
	s.Equal("proof_request", e.Type)
}

// tests getNextID and incrementID
func (s *KeeperSuite) TestIncrementID() {
	id, err := s.Keeper.getNextID(s.Context)
	s.SDKNil(err)
	s.Equal(id, types.RequestID{})

	err = s.Keeper.incrementID(s.Context)
	s.SDKNil(err)

	id, err = s.Keeper.getNextID(s.Context)
	s.SDKNil(err)
	s.Equal(types.RequestID{0, 0, 0, 0, 0, 0, 0, 1}, id)

	// errors if it cannot get next ID
	store := s.Keeper.getRequestStore(s.Context)
	idTag := []byte(types.RequestIdTag)
	store.Set(idTag, bytes.Repeat([]byte{9}, 9))

	err = s.Keeper.incrementID(s.Context)
	s.Equal(sdk.CodeType(107), err.Code())
}

func (s *KeeperSuite) TestHasRequest() {
	hasRequest := s.Keeper.hasRequest(s.Context, types.RequestID{})
	s.Equal(false, hasRequest)
	requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 4)
	s.Nil(requestErr)
	hasRequest = s.Keeper.hasRequest(s.Context, types.RequestID{})
	s.Equal(true, hasRequest)
}

func (s *KeeperSuite) TestSetRequest() {
	store := s.Keeper.getRequestStore(s.Context)
	idTag := []byte(types.RequestIdTag)
	store.Set(idTag, bytes.Repeat([]byte{9}, 9))

	err := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 0)
	s.Equal(sdk.CodeType(107), err.Code())
}

func (s *KeeperSuite) TestGetRequest() {
	requestRes := s.Fixtures.RequestTestCases.EmptyRequest
	request, err := s.Keeper.getRequest(s.Context, types.RequestID{})
	s.Equal(sdk.CodeType(601), err.Code())
	s.Equal(types.ProofRequest{}, request)

	requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 0)
	s.Nil(requestErr)

	request, err = s.Keeper.getRequest(s.Context, types.RequestID{})
	s.Nil(err)
	s.Equal(requestRes, request)
}

func (s *KeeperSuite) TestCheckRequests() {

}
