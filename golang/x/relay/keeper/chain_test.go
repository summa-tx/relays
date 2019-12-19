package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

	isMostRecent := s.Keeper.IsMostRecentCommonAncestor(s.Context, post[2].HashLE, post[3].HashLE, post[2].HashLE, 5)
	s.Equal(true, isMostRecent)

	isMostRecent = s.Keeper.IsMostRecentCommonAncestor(s.Context, post[5].HashLE, post[6].HashLE, tv.Orphan.HashLE, 5)
	s.Equal(true, isMostRecent)

	isMostRecent = s.Keeper.IsMostRecentCommonAncestor(s.Context, post[3].HashLE, post[3].HashLE, post[3].HashLE, 5)
	s.Equal(true, isMostRecent)

	isMostRecent = s.Keeper.IsMostRecentCommonAncestor(s.Context, post[0].HashLE, post[3].HashLE, post[2].HashLE, 5)
	s.Equal(false, isMostRecent)

	isMostRecent = s.Keeper.IsMostRecentCommonAncestor(s.Context, post[1].HashLE, post[3].HashLE, post[2].HashLE, 1)
	s.Equal(false, isMostRecent)
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

	_, err = s.Keeper.HeaviestFromAncestor(s.Context, tv.Headers[10].HashLE, headers[3].HashLE, headers[4].HashLE, 20)
	s.Equal(err.Code(), sdk.CodeType(103))

	_, err = s.Keeper.HeaviestFromAncestor(s.Context, headers[3].HashLE, tv.Headers[10].HashLE, headers[4].HashLE, 20)
	s.Equal(err.Code(), sdk.CodeType(103))

	_, err = s.Keeper.HeaviestFromAncestor(s.Context, headers[3].HashLE, headers[4].HashLE, tv.Headers[10].HashLE, 20)
	s.Equal(err.Code(), sdk.CodeType(103))

	_, err = s.Keeper.HeaviestFromAncestor(s.Context, headers[3].HashLE, headers[2].HashLE, headers[4].HashLE, 20)
	s.Equal(err.Code(), sdk.CodeType(104))

	_, err = s.Keeper.HeaviestFromAncestor(s.Context, headers[3].HashLE, headers[4].HashLE, headers[2].HashLE, 20)
	s.Equal(err.Code(), sdk.CodeType(104))

	heaviest, err := s.Keeper.HeaviestFromAncestor(s.Context, headers[3].HashLE, headers[5].HashLE, headers[4].HashLE, 20)
	s.SDKNil(err)
	s.Equal(heaviest, headers[5].HashLE)

	heaviest, err = s.Keeper.HeaviestFromAncestor(s.Context, headers[3].HashLE, headers[4].HashLE, headers[5].HashLE, 20)
	s.SDKNil(err)
	s.Equal(heaviest, headers[5].HashLE)

	heaviest, err = s.Keeper.HeaviestFromAncestor(s.Context, headers[3].HashLE, tv.Headers[8].HashLE, tv.Orphan.HashLE, 20)
	s.SDKNil(err)
	s.Equal(heaviest, tv.Headers[8].HashLE)

	heaviest, err = s.Keeper.HeaviestFromAncestor(s.Context, headers[3].HashLE, tv.Orphan.HashLE, tv.Headers[8].HashLE, 20)
	s.SDKNil(err)
	s.Equal(heaviest, tv.Orphan.HashLE)
}
