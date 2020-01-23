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
	valid, err := s.Keeper.validateProof(s.Context, proof, types.RequestID{})
	s.Equal(err.Code(), sdk.CodeType(105))
	s.Equal(false, valid)

	// errors if link is not found
	s.Keeper.setLastReorgLCA(s.Context, proofCases[0].LCA)

	valid, err = s.Keeper.validateProof(s.Context, proof, types.RequestID{})
	s.Equal(err.Code(), sdk.CodeType(610))
	s.Equal(false, valid)

	// errors if request is not found
	s.Keeper.ingestHeader(s.Context, proof.ConfirmingHeader)
	s.Keeper.setLink(s.Context, proof.ConfirmingHeader)

	valid, err = s.Keeper.validateProof(s.Context, proof, types.RequestID{})
	s.Equal(err.Code(), sdk.CodeType(601))
	s.Equal(false, valid)

	// errors if Best Known Digest is not found
	requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 4)
	s.Nil(requestErr)

	valid, err = s.Keeper.validateProof(s.Context, proof, types.RequestID{})
	s.Equal(err.Code(), sdk.CodeType(105))
	s.Equal(false, valid)

	// errors if Best Known Digest header is not found
	s.Keeper.setBestKnownDigest(s.Context, proofCases[0].BestKnown.HashLE)

	valid, err = s.Keeper.validateProof(s.Context, proof, types.RequestID{})
	s.Equal(err.Code(), sdk.CodeType(103))
	s.Equal(false, valid)

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
			valid, err := s.Keeper.validateProof(s.Context, proofCases[i].Proof, types.RequestID{0, 0, 0, 0, 0, 0, 0, requestID})
			s.Equal(sdk.CodeType(proofCases[i].Error), err.Code())
			s.Equal(proofCases[i].Output, valid)
		} else {
			valid, err := s.Keeper.validateProof(s.Context, proofCases[i].Proof, types.RequestID{0, 0, 0, 0, 0, 0, 0, requestID})
			s.Nil(err)
			s.Equal(proofCases[i].Output, valid)
		}
	}
}

func (s *KeeperSuite) TestCheckRequestsFilled() {
	tc := s.Fixtures.ValidatorTestCases.CheckRequestsFilled
	validProof := s.Fixtures.ValidatorTestCases.ValidateProof[0]

	// // Success
	// out, _ := btcspv.ExtractOutputAtIndex(v.Vout, v.OutputIdx)
	// outputScript := out[8:]

	// in := btcspv.ExtractInputAtIndex(v.Vin, v.InputIdx)
	// outpoint := btcspv.ExtractOutpoint(in)

	// requestErr := s.Keeper.setRequest(s.Context, outpoint, outputScript, 10, 255)
	// s.SDKNil(requestErr)
	// valid, err := s.Keeper.checkRequests(
	// 	s.Context,
	// 	v.InputIdx,
	// 	v.OutputIdx,
	// 	v.Vin,
	// 	v.Vout,
	// 	types.RequestID{})
	// s.SDKNil(err)
	// s.Equal(true, valid)

	s.Keeper.setLastReorgLCA(s.Context, validProof.LCA)
	s.Keeper.ingestHeader(s.Context, validProof.Proof.ConfirmingHeader)
	s.Keeper.setLink(s.Context, validProof.Proof.ConfirmingHeader)
	s.Keeper.ingestHeader(s.Context, validProof.BestKnown)
	s.Keeper.setBestKnownDigest(s.Context, validProof.BestKnown.HashLE)
	requestErr := s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 4)
	s.Nil(requestErr)

	for i := range tc {
		valid, err := s.Keeper.checkRequestsFilled(s.Context, tc[i].FilledRequests)
		if tc[i].Error != 0 {
			s.Equal(false, valid)
			s.Equal(sdk.CodeType(tc[i].Error), err.Code())
		} else {
			s.SDKNil(err)
			s.Equal(true, valid)
		}
	}
}
