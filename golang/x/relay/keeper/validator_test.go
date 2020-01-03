package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (s *KeeperSuite) TestValidateProof() {
	proofCases := s.Fixtures.ValidatorTestCases.ValidateProof

	// errors if LCA is not found
	valid, err := validateProof(s.Context, s.Keeper, proofCases[0].Proof)
	s.Equal(err.Code(), sdk.CodeType(105))
	s.Equal(false, valid)

	// errors if link is not found
	s.Keeper.setLastReorgLCA(s.Context, proofCases[0].LCA)
	valid, err = validateProof(s.Context, s.Keeper, proofCases[0].Proof)
	s.Nil(err)
	s.Equal(false, valid)

	for i := range proofCases {
		s.Keeper.setLastReorgLCA(s.Context, proofCases[i].LCA)
		s.Keeper.ingestHeader(s.Context, proofCases[i].Proof.ConfirmingHeader)
		s.Keeper.setLink(s.Context, proofCases[i].Proof.ConfirmingHeader)

		if proofCases[i].Error != 0 {
			valid, err := validateProof(s.Context, s.Keeper, proofCases[i].Proof)
			s.Equal(err.Code(), sdk.CodeType(proofCases[i].Error))
			s.Equal(proofCases[i].Output, valid)
		} else {
			valid, err := validateProof(s.Context, s.Keeper, proofCases[i].Proof)
			s.Nil(err)
			s.Equal(proofCases[i].Output, valid)
		}
	}
}
