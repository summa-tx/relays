package keeper

func (suite *KeeperSuite) TestValidateHeaderChain() {
	cases := suite.Fixtures.HeaderTestCases.ValidateChain

	for _, tc := range cases {
		err := validateHeaderChain(tc.Anchor, tc.Headers, tc.Internal, tc.IsMainnet)
		if tc.Output == 0 {
			logIfTestCaseError(tc, err)
			suite.Nil(err)
		} else {
			suite.NotNil(err)
			suite.Equal(tc.Output, err.Code())
		}
	}
}

func (suite *KeeperSuite) TestValidateDifficultyChange() {
	cases := suite.Fixtures.HeaderTestCases.ValidateDiffChange

	for _, tc := range cases {
		err := validateDifficultyChange(tc.Headers, tc.PrevEpochStart, tc.Anchor)
		if tc.Output == 0 {
			logIfTestCaseError(tc, err)
			suite.Nil(err)
		} else {
			suite.NotNil(err)
			suite.Equal(tc.Output, err.Code())
		}
	}
}

func (suite *KeeperSuite) TestCompareTargets() {
	cases := suite.Fixtures.HeaderTestCases.CompareTargets

	for _, tc := range cases {
		result := compareTargets(tc.Full, tc.Truncated)
		suite.Equal(result, tc.Output)
	}
}
