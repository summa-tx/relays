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

func (s *KeeperSuite) TestIngestHeaders() {
	cases := s.Fixtures.HeaderTestCases.ValidateChain

	for _, tc := range cases {
		s.InitTestContext(tc.IsMainnet, false)
		s.Keeper.ingestHeader(s.Context, tc.Anchor)
		err := s.Keeper.ingestHeaders(s.Context, tc.Headers, tc.Internal)
		if tc.Output == 0 {
			logIfTestCaseError(tc, err)
			s.Nil(err)
		} else {
			s.NotNil(err)
			s.Equal(tc.Output, err.Code())
		}
	}
}

func (s *KeeperSuite) TestIngestHeaderChain() {
	cases := s.Fixtures.HeaderTestCases.ValidateChain

	for _, tc := range cases {
		if tc.Internal == false {
			s.InitTestContext(tc.IsMainnet, false)
			s.Keeper.ingestHeader(s.Context, tc.Anchor)
			err := s.Keeper.IngestHeaderChain(s.Context, tc.Headers)
			if tc.Output == 0 {
				logIfTestCaseError(tc, err)
				s.Nil(err)
			} else {
				s.NotNil(err)
				s.Equal(tc.Output, err.Code())
			}
		}
	}
}

// TestIngestHeader tests ingestHeader, HasHeader, and GetHeader
func (s *KeeperSuite) TestIngestHeader() {
	cases := s.Fixtures.HeaderTestCases.ValidateChain

	for _, tc := range cases {
		s.Keeper.ingestHeader(s.Context, tc.Headers[0])
		hasHeader := s.Keeper.HasHeader(s.Context, tc.Headers[0].HashLE)
		s.Equal(hasHeader, true)
		header, err := s.Keeper.GetHeader(s.Context, tc.Headers[0].HashLE)
		s.Nil(err)
		s.Equal(header, tc.Headers[0])
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

func (s *KeeperSuite) TestIngestDifficultyChange() {
	cases := s.Fixtures.HeaderTestCases.ValidateDiffChange

	for _, tc := range cases {
		s.Keeper.ingestHeader(s.Context, tc.PrevEpochStart)
		s.Keeper.ingestHeader(s.Context, tc.Anchor)
		err := s.Keeper.IngestDifficultyChange(s.Context, tc.PrevEpochStart.HashLE, tc.Headers)
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
