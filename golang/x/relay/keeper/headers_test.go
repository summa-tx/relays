package keeper

func (s *KeeperSuite) TestValidateHeaderChain() {
	cases := s.Fixtures.HeaderTestCases.ValidateChain

	for _, tc := range cases {
		err := validateHeaderChain(tc.Anchor, tc.Headers, tc.Internal, tc.IsMainnet)
		if tc.Output == 0 {
			logIfTestCaseError(tc, err)
			s.Nil(err)
		} else {
			s.NotNil(err)
			s.Equal(tc.Output, err.Code())
		}
	}
}

func (s *KeeperSuite) TestValidateDifficultyChange() {
	cases := s.Fixtures.HeaderTestCases.ValidateDiffChange

	for _, tc := range cases {
		err := validateDifficultyChange(tc.Headers, tc.PrevEpochStart, tc.Anchor)
		if tc.Output == 0 {
			logIfTestCaseError(tc, err)
			s.Nil(err)
		} else {
			s.NotNil(err)
			s.Equal(tc.Output, err.Code())
		}
	}
}

func (s *KeeperSuite) TestCompareTargets() {
	cases := s.Fixtures.HeaderTestCases.CompareTargets

	for _, tc := range cases {
		result := compareTargets(tc.Full, tc.Truncated)
		s.Equal(result, tc.Output)
	}
}
