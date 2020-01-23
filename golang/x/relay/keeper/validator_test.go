package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

// TODO: Fix IncrementID, ID is not incrementing when a new request is set
func (s *KeeperSuite) TestValidateProof() {
	// TODO: Add Request ID and NumConfs check to test
	proofCases := s.Fixtures.ValidatorTestCases.ValidateProof
	proof := proofCases[0].Proof

	// errors if LCA is not found
	err := s.Keeper.validateProof(s.Context, proof, types.RequestID{})
	s.Equal(sdk.CodeType(105), err.Code())

	// errors if link is not found
	s.Keeper.setLastReorgLCA(s.Context, proofCases[0].LCA)

	err = s.Keeper.validateProof(s.Context, proof, types.RequestID{})
	s.Equal(sdk.CodeType(610), err.Code())

	// errors if request is not found
	s.Keeper.ingestHeader(s.Context, proof.ConfirmingHeader)
	s.Keeper.setLink(s.Context, proof.ConfirmingHeader)

	err = s.Keeper.validateProof(s.Context, proof, types.RequestID{})
	s.Equal(sdk.CodeType(601), err.Code())

	// errors if Best Known Digest is not found
	requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 4)
	s.Nil(requestErr)

	err = s.Keeper.validateProof(s.Context, proof, types.RequestID{})
	s.Equal(sdk.CodeType(105), err.Code())

	// errors if Best Known Digest header is not found
	s.Keeper.setBestKnownDigest(s.Context, proofCases[0].BestKnown.HashLE)

	err = s.Keeper.validateProof(s.Context, proof, types.RequestID{})
	s.Equal(sdk.CodeType(103), err.Code())

	for i := range proofCases {
		requestID := byte(i + 1)
		// Store lots of stuff
		s.Keeper.setLastReorgLCA(s.Context, proofCases[i].LCA)
		s.Keeper.ingestHeader(s.Context, proofCases[i].Proof.ConfirmingHeader)
		s.Keeper.setLink(s.Context, proofCases[i].Proof.ConfirmingHeader)
		s.Keeper.ingestHeader(s.Context, proofCases[i].BestKnown)
		s.Keeper.setBestKnownDigest(s.Context, proofCases[i].BestKnown.HashLE)
		requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 4)
		s.Nil(requestErr)

		if proofCases[i].Error != 0 {
			err := s.Keeper.validateProof(s.Context, proofCases[i].Proof, types.RequestID{0, 0, 0, 0, 0, 0, 0, requestID})
			s.Equal(sdk.CodeType(proofCases[i].Error), err.Code())
		} else {
			err := s.Keeper.validateProof(s.Context, proofCases[i].Proof, types.RequestID{0, 0, 0, 0, 0, 0, 0, requestID})
			s.Nil(err)
		}
	}
}

func (s *KeeperSuite) TestCheckRequestsFilled() {
	tc := s.Fixtures.ValidatorTestCases.CheckRequestsFilled
	validProof := s.Fixtures.ValidatorTestCases.ValidateProof[0]

	s.Keeper.setLastReorgLCA(s.Context, validProof.LCA)
	s.Keeper.ingestHeader(s.Context, validProof.Proof.ConfirmingHeader)
	s.Keeper.setLink(s.Context, validProof.Proof.ConfirmingHeader)
	s.Keeper.ingestHeader(s.Context, validProof.BestKnown)
	requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 4)
	s.Nil(requestErr)

	// errors if getConfs fails
	// TODO: this is failing inside the ValidateProof check
	err := s.Keeper.checkRequestsFilled(s.Context, tc[0].FilledRequests)
	s.Equal(sdk.CodeType(105), err.Code())

	s.Keeper.setBestKnownDigest(s.Context, validProof.BestKnown.HashLE)

	// errors if checkRequest errors
	// deactivate request
	activeErr := s.Keeper.setRequestState(s.Context, types.RequestID{}, false)
	s.SDKNil(activeErr)

	err = s.Keeper.checkRequestsFilled(s.Context, tc[0].FilledRequests)
	s.Equal(sdk.CodeType(606), err.Code())

	// reactivate request
	activeErr = s.Keeper.setRequestState(s.Context, types.RequestID{}, true)
	s.SDKNil(activeErr)

	for i := range tc {
		err := s.Keeper.checkRequestsFilled(s.Context, tc[i].FilledRequests)
		if tc[i].Error != 0 {
			s.Equal(sdk.CodeType(tc[i].Error), err.Code())
		} else {
			s.SDKNil(err)
		}
	}

	// errors if number of confirmations is less than the number of confirmations on the request
	requestErr = s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 5)
	s.Nil(requestErr)

	copiedRequest := tc[0].FilledRequests
	copiedRequest.Requests[1].ID = types.RequestID{0, 0, 0, 0, 0, 0, 0, 1}
	err = s.Keeper.checkRequestsFilled(s.Context, copiedRequest)
	s.Equal(sdk.CodeType(611), err.Code())
}
