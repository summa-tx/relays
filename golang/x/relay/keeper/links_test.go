package keeper

func (s *KeeperSuite) TestGetLink() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	parent := headers[0]
	child := headers[1]

	s.Keeper.setLink(s.Context, child)
	hasHeader := s.Keeper.hasLink(s.Context, child.HashLE)
	s.Equal(true, hasHeader)
	getHeader := s.Keeper.getLink(s.Context, child.HashLE)
	s.Equal(getHeader, parent.HashLE)
}

// TODO: Add test cases: add to JSON, loop over
func (s *KeeperSuite) TestFindAncestor() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	anchor := s.Fixtures.HeaderTestCases.ValidateChain[0].Anchor

	s.Keeper.ingestHeader(s.Context, anchor)
	s.Keeper.IngestHeaderChain(s.Context, headers)

	ancestor, err := s.Keeper.FindAncestor(s.Context, headers[4].HashLE, 2)
	s.Nil(err)
	s.Equal(headers[2].HashLE, ancestor)
}

func (s *KeeperSuite) TestIsAncestor() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	anchor := s.Fixtures.HeaderTestCases.ValidateChain[0].Anchor

	s.Keeper.ingestHeader(s.Context, anchor)
	s.Keeper.IngestHeaderChain(s.Context, headers)

	isAncestor := s.Keeper.IsAncestor(s.Context, headers[4].HashLE, headers[1].HashLE, 15)
	s.Equal(true, isAncestor)

	isAncestor = s.Keeper.IsAncestor(s.Context, headers[1].HashLE, headers[4].HashLE, 15)
	s.Equal(false, isAncestor)
}
