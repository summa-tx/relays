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
