package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: Fix IncrementID, ID is not incrementing when a new request is set
func (s *KeeperSuite) TestValidateProof() {
	// TODO: Add Request ID and NumConfs check to test
	proofCases := s.Fixtures.ValidatorTestCases.ValidateProof
	proof := proofCases[0].Proof

	// errors if LCA is not found
	valid, err := s.Keeper.validateProof(s.Context, proof, 0)
	s.Equal(err.Code(), sdk.CodeType(105))
	s.Equal(false, valid)

	// errors if link is not found
	s.Keeper.setLastReorgLCA(s.Context, proofCases[0].LCA)

	valid, err = s.Keeper.validateProof(s.Context, proof, 0)
	s.Equal(err.Code(), sdk.CodeType(610))
	s.Equal(false, valid)

	// errors if request is not found
	s.Keeper.ingestHeader(s.Context, proof.ConfirmingHeader)
	s.Keeper.setLink(s.Context, proof.ConfirmingHeader)

	valid, err = s.Keeper.validateProof(s.Context, proof, 0)
	s.Equal(err.Code(), sdk.CodeType(601))
	s.Equal(false, valid)

	// errors if Best Known Digest is not found
	requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 4)
	s.Nil(requestErr)

	valid, err = s.Keeper.validateProof(s.Context, proof, 0)
	s.Equal(err.Code(), sdk.CodeType(105))
	s.Equal(false, valid)

	// errors if Best Known Digest header is not found
	s.Keeper.setBestKnownDigest(s.Context, proofCases[0].BestKnown.HashLE)

	valid, err = s.Keeper.validateProof(s.Context, proof, 0)
	s.Equal(err.Code(), sdk.CodeType(103))
	s.Equal(false, valid)

	for i := range proofCases {
		requestID := uint64(i + 1)
		// Store lots of stuff
		s.Keeper.setLastReorgLCA(s.Context, proofCases[i].LCA)
		s.Keeper.ingestHeader(s.Context, proofCases[i].Proof.ConfirmingHeader)
		s.Keeper.setLink(s.Context, proofCases[i].Proof.ConfirmingHeader)
		s.Keeper.ingestHeader(s.Context, proofCases[i].BestKnown)
		s.Keeper.setBestKnownDigest(s.Context, proofCases[i].BestKnown.HashLE)
		requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 4)
		s.Nil(requestErr)

		if proofCases[i].Error != 0 {
			valid, err := s.Keeper.validateProof(s.Context, proofCases[i].Proof, requestID)
			s.Equal(err.Code(), sdk.CodeType(proofCases[i].Error))
			s.Equal(proofCases[i].Output, valid)
		} else {
			valid, err := s.Keeper.validateProof(s.Context, proofCases[i].Proof, requestID)
			s.Nil(err)
			s.Equal(proofCases[i].Output, valid)
		}
	}
}
