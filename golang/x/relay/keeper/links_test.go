package keeper

import "github.com/cosmos/cosmos-sdk/types"

func (s *KeeperSuite) TestGetLink() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	parent := headers[0]
	child := headers[1]

	// stores and retrieves link
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

	// errors if link is not found
	_, err := s.Keeper.FindAncestor(s.Context, headers[4].HashLE, 2)
	s.Equal(err.Code(), types.CodeType(103))

	s.Keeper.ingestHeader(s.Context, anchor)
	s.Keeper.IngestHeaderChain(s.Context, headers)

	// successfully retrieves ancestor
	ancestor, err := s.Keeper.FindAncestor(s.Context, headers[4].HashLE, 2)
	s.Nil(err)
	s.Equal(headers[2].HashLE, ancestor)

	// errors if link is not found
	// this occurs when the offset overflows the length of the header chain
	_, err = s.Keeper.FindAncestor(s.Context, headers[1].HashLE, 3)
	s.Equal(err.Code(), types.CodeType(103))
}

func (s *KeeperSuite) TestIsAncestor() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	anchor := s.Fixtures.HeaderTestCases.ValidateChain[0].Anchor

	s.Keeper.ingestHeader(s.Context, anchor)
	s.Keeper.IngestHeaderChain(s.Context, headers)

	// is ancestor
	isAncestor := s.Keeper.IsAncestor(s.Context, headers[4].HashLE, headers[1].HashLE, 15)
	s.Equal(true, isAncestor)

	// is not ancestor
	isAncestor = s.Keeper.IsAncestor(s.Context, headers[1].HashLE, headers[4].HashLE, 15)
	s.Equal(false, isAncestor)

	// is not ancestor
	isAncestor = s.Keeper.IsAncestor(s.Context, headers[1].HashLE, headers[4].HashLE, 0)
	s.Equal(false, isAncestor)
}
