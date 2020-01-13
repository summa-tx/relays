package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperSuite) TestValidateProof() {
	// TODO: Add Request ID and NumConfs check to test
	proofCases := s.Fixtures.ValidatorTestCases.ValidateProof

	// errors if LCA is not found
	valid, err := s.Keeper.validateProof(s.Context, proofCases[0].Proof, 0)
	s.Equal(err.Code(), sdk.CodeType(105))
	s.Equal(false, valid)

	// errors if link is not found
	s.Keeper.setLastReorgLCA(s.Context, proofCases[0].LCA)
	valid, err = s.Keeper.validateProof(s.Context, proofCases[0].Proof, 0)
	s.Equal(err.Code(), sdk.CodeType(610))
	s.Equal(false, valid)

	for i := range proofCases {
		s.Keeper.setLastReorgLCA(s.Context, proofCases[i].LCA)
		s.Keeper.ingestHeader(s.Context, proofCases[i].Proof.ConfirmingHeader)
		s.Keeper.setLink(s.Context, proofCases[i].Proof.ConfirmingHeader)
		// Store request
		// TODO: update request with correct numConfs
		requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 4)
		s.Nil(requestErr)

		if proofCases[i].Error != 0 {
			valid, err := s.Keeper.validateProof(s.Context, proofCases[i].Proof, uint64(i))
			s.Equal(err.Code(), sdk.CodeType(proofCases[i].Error))
			s.Equal(proofCases[i].Output, valid)
		} else {
			valid, err := s.Keeper.validateProof(s.Context, proofCases[i].Proof, uint64(i))
			s.Nil(err)
			s.Equal(proofCases[i].Output, valid)
		}
	}
}
