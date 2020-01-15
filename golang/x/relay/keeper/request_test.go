package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

func (s *KeeperSuite) TestEmitProofRequest() {
	s.Keeper.emitProofRequest(s.Context, []byte{0}, []byte{0}, 0, 3)

	events := s.Context.EventManager().Events()
	e := events[0]
	s.Equal(e.Type, "proof_request")
}

// tests getNextID and incrementID
func (s *KeeperSuite) TestIncrementID() {
	id := s.Keeper.getNextID(s.Context)
	s.Equal(id, []byte{0, 0, 0, 0, 0, 0, 0, 0})
	s.Keeper.incrementID(s.Context)
	id = s.Keeper.getNextID(s.Context)
	s.Equal(id, []byte{0, 0, 0, 0, 0, 0, 0, 1})
}

func (s *KeeperSuite) TestHasRequest() {
	hasRequest := s.Keeper.hasRequest(s.Context, 0)
	s.Equal(false, hasRequest)
	requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 4)
	s.Nil(requestErr)
	hasRequest = s.Keeper.hasRequest(s.Context, 0)
	s.Equal(true, hasRequest)
}

func (s *KeeperSuite) TestGetRequest() {
	requestRes := types.ProofRequest{
		Spends:      types.Hash256Digest{0x14, 0x6, 0xe0, 0x58, 0x81, 0xe2, 0x99, 0x36, 0x77, 0x66, 0xd3, 0x13, 0xe2, 0x6c, 0x5, 0x56, 0x4e, 0xc9, 0x1b, 0xf7, 0x21, 0xd3, 0x17, 0x26, 0xbd, 0x6e, 0x46, 0xe6, 0x6, 0x89, 0x53, 0x9a},
		Pays:        types.Hash256Digest{0x14, 0x6, 0xe0, 0x58, 0x81, 0xe2, 0x99, 0x36, 0x77, 0x66, 0xd3, 0x13, 0xe2, 0x6c, 0x5, 0x56, 0x4e, 0xc9, 0x1b, 0xf7, 0x21, 0xd3, 0x17, 0x26, 0xbd, 0x6e, 0x46, 0xe6, 0x6, 0x89, 0x53, 0x9a},
		PaysValue:   0,
		ActiveState: true,
		NumConfs:    0,
	}
	request, err := s.Keeper.getRequest(s.Context, 0)
	s.Equal(err.Code(), sdk.CodeType(601))
	s.Equal(types.ProofRequest{}, request)

	requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 0)
	s.Nil(requestErr)

	request, err = s.Keeper.getRequest(s.Context, 0)
	s.Nil(err)
	s.Equal(requestRes, request)
}
