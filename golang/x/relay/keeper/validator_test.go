package keeper

func (s *KeeperSuite) TestValidateProof() {
	// headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	// anchor := s.Fixtures.HeaderTestCases.ValidateChain[0].Anchor
	// tc := s.Fixtures.LinkTestCases.FindAncestor.TestCases

	// // errors if link is not found
	// _, err := s.Keeper.FindAncestor(s.Context, tc[0].Digest, tc[0].Offset)
	// s.Equal(err.Code(), types.CodeType(tc[0].Error))

	// s.Keeper.ingestHeader(s.Context, anchor)
	// s.Keeper.IngestHeaderChain(s.Context, headers)

	// for i := 1; i < len(tc); i++ {
	// 	if tc[i].Error == 0 {
	// 		// successfully retrieves ancestor
	// 		ancestor, err := s.Keeper.FindAncestor(s.Context, tc[i].Digest, tc[i].Offset)
	// 		s.SDKNil(err)
	// 		s.Equal(tc[i].Output, ancestor)
	// 	} else {
	// 		// errors if link is not found
	// 		_, err = s.Keeper.FindAncestor(s.Context, tc[i].Digest, tc[i].Offset)
	// 		s.Equal(err.Code(), types.CodeType(tc[i].Error))
	// 	}
	// }
}
