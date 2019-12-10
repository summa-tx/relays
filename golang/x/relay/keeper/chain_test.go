package keeper

func (s *KeeperSuite) TestEmitReorg() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	s.Keeper.emitReorg(s.Context, headers[0].HashLE, headers[1].HashLE, headers[2].HashLE)

	events := s.Context.EventManager().Events()
	e := events[0]
	s.Equal(e.Type, "reorg")
}

func (s *KeeperSuite) TestGetBestKnownDigest() {
	digest := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers[0].HashLE
	s.Keeper.setBestKnownDigest(s.Context, digest)
	bestKnown, _ := s.Keeper.GetBestKnownDigest(s.Context)
	s.Equal(digest, bestKnown)
}

func (s *KeeperSuite) TestGetLastReorgLCA() {
	digest := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers[0].HashLE
	s.Keeper.setLastReorgLCA(s.Context, digest)
	lca, _ := s.Keeper.GetLastReorgLCA(s.Context)
	s.Equal(digest, lca)
}
