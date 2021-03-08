package keeper

import "github.com/cosmos/cosmos-sdk/types"

func (s *KeeperSuite) TestGetLink() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	parent := headers[0]
	child := headers[1]

	// stores and retrieves link
	s.Keeper.setLink(s.Context, child)
	hasHeader := s.Keeper.hasLink(s.Context, child.Hash)
	s.Equal(true, hasHeader)
	getHeader := s.Keeper.getLink(s.Context, child.Hash)
	s.Equal(parent.Hash, getHeader)
}

func (s *KeeperSuite) TestFindAncestor() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	anchor := s.Fixtures.HeaderTestCases.ValidateChain[0].Anchor
	tc := s.Fixtures.LinkTestCases.FindAncestor.TestCases

	// errors if link is not found
	_, err := s.Keeper.FindAncestor(s.Context, tc[0].Digest, tc[0].Offset)
	s.Equal(types.CodeType(tc[0].Error), err.Code())

	s.Keeper.ingestHeader(s.Context, anchor)
	err = s.Keeper.IngestHeaderChain(s.Context, headers)
	s.SDKNil(err)

	for i := 1; i < len(tc); i++ {
		ancestor, err := s.Keeper.FindAncestor(s.Context, tc[i].Digest, tc[i].Offset)
		if tc[i].Error == 0 {
			s.SDKNil(err)
			s.Equal(tc[i].Output, ancestor)
		} else {
			s.Equal(types.CodeType(tc[i].Error), err.Code())
		}
	}
}

func (s *KeeperSuite) TestIsAncestor() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	anchor := s.Fixtures.HeaderTestCases.ValidateChain[0].Anchor
	tc := s.Fixtures.LinkTestCases.IsAncestor.TestCases

	s.Keeper.ingestHeader(s.Context, anchor)
	err := s.Keeper.IngestHeaderChain(s.Context, headers)
	s.SDKNil(err)

	for i := range tc {
		isAncestor := s.Keeper.IsAncestor(s.Context, tc[i].Digest, tc[i].Ancestor, tc[i].Limit)
		s.Equal(tc[i].Output, isAncestor)
	}
}
