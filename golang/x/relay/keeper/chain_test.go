package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

func (s *KeeperSuite) TestEmitReorg() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	s.Keeper.emitReorg(s.Context, headers[0].HashLE, headers[1].HashLE, headers[2].HashLE)

	events := s.Context.EventManager().Events()
	e := events[0]
	s.Equal(e.Type, "reorg")
}

func (s *KeeperSuite) TestGetDigestByStoreKey() {
	wrongLenDigest := bytes.Repeat([]byte{0}, 31)
	key := "bad-digest"

	store := s.Keeper.getChainStore(s.Context)
	store.Set([]byte(key), wrongLenDigest)

	_, err := s.Keeper.getDigestByStoreKey(s.Context, key)
	s.Equal(sdk.CodeType(105), err.Code())
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

func (s *KeeperSuite) TestIsMostRecentCommonAncestor() {
	tv := s.Fixtures.ChainTestCases.IsMostRecentCA
	pre := tv.PreRetargetChain
	post := tv.PostRetargetChain

	var postWithOrphan []types.BitcoinHeader
	postWithOrphan = append(postWithOrphan, post[:len(post)-2]...)
	postWithOrphan = append(postWithOrphan, tv.Orphan)

	err := s.Keeper.SetGenesisState(s.Context, tv.Genesis, tv.OldPeriodStart)
	s.SDKNil(err)

	err = s.Keeper.IngestHeaderChain(s.Context, pre)
	s.SDKNil(err)
	err = s.Keeper.IngestDifficultyChange(s.Context, tv.OldPeriodStart.HashLE, post)
	s.SDKNil(err)
	err = s.Keeper.IngestDifficultyChange(s.Context, tv.OldPeriodStart.HashLE, postWithOrphan)
	s.SDKNil(err)

	for i := range tv.TestCases {
		isMostRecent := s.Keeper.IsMostRecentCommonAncestor(
			s.Context,
			tv.TestCases[i].Ancestor,
			tv.TestCases[i].Left,
			tv.TestCases[i].Right,
			tv.TestCases[i].Limit)
		s.Equal(tv.TestCases[i].Output, isMostRecent)
	}
}

func (s *KeeperSuite) TestHeaviestFromAncestor() {
	tv := s.Fixtures.ChainTestCases.HeaviestFromAncestor
	headers := tv.Headers[0:8]
	headersWithMain := tv.Headers[0:9]

	var headersWithOrphan []types.BitcoinHeader
	headersWithOrphan = append(headersWithOrphan, headers...)
	headersWithOrphan = append(headersWithOrphan, tv.Orphan)

	s.Keeper.ingestHeader(s.Context, tv.Genesis)
	err := s.Keeper.IngestHeaderChain(s.Context, headersWithMain)
	s.SDKNil(err)
	err = s.Keeper.IngestHeaderChain(s.Context, headersWithOrphan)
	s.SDKNil(err)

	for i := range tv.TestCases {
		if tv.TestCases[i].Error == 0 {
			heaviest, err := s.Keeper.HeaviestFromAncestor(
				s.Context,
				tv.TestCases[i].Ancestor,
				tv.TestCases[i].CurrentBest,
				tv.TestCases[i].NewBest,
				tv.TestCases[i].Limit)
			s.SDKNil(err)
			s.Equal(heaviest, tv.TestCases[i].Output)
		} else {
			_, err = s.Keeper.HeaviestFromAncestor(
				s.Context,
				tv.TestCases[i].Ancestor,
				tv.TestCases[i].CurrentBest,
				tv.TestCases[i].NewBest,
				tv.TestCases[i].Limit)
			s.Equal(err.Code(), sdk.CodeType(tv.TestCases[i].Error))
		}
	}
}

func (s *KeeperSuite) TestMarkNewHeaviest() {
	tv := s.Fixtures.ChainTestCases.IsMostRecentCA
	tc := s.Fixtures.ChainTestCases.MarkNewHeaviest
	pre := tv.PreRetargetChain
	post := tv.PostRetargetChain
	var postWithOrphan []types.BitcoinHeader
	postWithOrphan = append(postWithOrphan, post[:len(post)-2]...)
	postWithOrphan = append(postWithOrphan, tv.Orphan)

	err := s.Keeper.SetGenesisState(s.Context, tv.Genesis, tv.OldPeriodStart)
	s.SDKNil(err)

	err = s.Keeper.MarkNewHeaviest(
		s.Context,
		tv.Genesis.HashLE,
		pre[0].Raw,
		pre[1].Raw,
		10,
	)
	s.EqualError(err, 103)

	err = s.Keeper.IngestHeaderChain(s.Context, pre)
	s.SDKNil(err)
	err = s.Keeper.IngestDifficultyChange(s.Context, tv.OldPeriodStart.HashLE, post)
	s.SDKNil(err)
	err = s.Keeper.IngestDifficultyChange(s.Context, tv.OldPeriodStart.HashLE, postWithOrphan)
	s.SDKNil(err)

	for i := range tc {
		curBestDigest := btcspv.Hash256(tc[i].CurrentBest[:])
		s.Keeper.setBestKnownDigest(s.Context, curBestDigest)
		if tc[i].Error == 0 {
			// updates the best known and emits an event
			err = s.Keeper.MarkNewHeaviest(
				s.Context,
				tc[i].Ancestor,
				tc[i].CurrentBest,
				tc[i].NewBest,
				tc[i].Limit,
			)
			s.SDKNil(err)
			events := s.Context.EventManager().Events()
			e := events[0]
			s.Equal(e.Type, tc[i].Output)
		} else {
			err = s.Keeper.MarkNewHeaviest(
				s.Context,
				tc[i].Ancestor,
				tc[i].CurrentBest,
				tc[i].NewBest,
				tc[i].Limit,
			)
			s.Equal(err.Code(), sdk.CodeType(tc[i].Error))
		}
	}

	// errors if the ancestor is not the heaviest common ancestor
	err = s.Keeper.MarkNewHeaviest(
		s.Context,
		tv.Genesis.HashLE,
		tv.Genesis.Raw,
		pre[0].Raw,
		10,
	)
	s.SDKNil(err)
	err = s.Keeper.MarkNewHeaviest(
		s.Context,
		tv.Genesis.HashLE,
		pre[0].Raw,
		pre[1].Raw,
		10,
	)
	s.Equal(err.Code(), sdk.CodeType(404))
}
