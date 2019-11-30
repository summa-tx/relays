package keeper

import "log"

func (suite *KeeperSuite) TestValidateHeaderChain() {
	cases := suite.Fixtures.HeaderTestCases.ValidateChain

	for _, tc := range cases {
		err := validateHeaderChain(tc.Anchor, tc.Headers, tc.Internal, tc.IsMainnet)
		if tc.Output == 0 {
			if err != nil {
				log.Printf("Unexpected Error\n%s\n", err.Error())
			}
			suite.Nil(err)
		} else {
			suite.Equal(tc.Output, err.Code())
		}
	}
}

func (suite *KeeperSuite) TestValidateDifficultyChange() {
	cases := suite.Fixtures.HeaderTestCases.ValidateDiffChange

	for _, tc := range cases {
		err := validateDifficultyChange(tc.Headers, tc.PrevEpochStart, tc.Anchor)
		if tc.Output == 0 {
			if err != nil {
				log.Printf("Unexpected Error\n%s\n", err.Error())
			}
			suite.Nil(err)
		} else {
			suite.Equal(tc.Output, err.Code())
		}
	}

}
