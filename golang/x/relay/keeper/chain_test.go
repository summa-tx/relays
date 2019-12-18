package keeper

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/types"
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
	s.Equal(types.CodeType(105), err.Code())
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

// func (s *KeeperSuite) TestIsMostRecentCommonAncestor() {
// 	tv := s.Fixtures.ChainTestCases.IsMostRecentCA

// 	err := s.Keeper.SetGenesisState(s.Context, tv.Genesis, tv.OldPeriodStart)
// 	s.Nil(err)

// 	err = s.Keeper.IngestDifficultyChange(s.Context, tv.OldPeriodStart.HashLE, tv.PreRetargetChain)
// 	s.Nil(err)
// }

// describe('#isMostRecentAncestor', async () => {
// 	const PRE_CHAIN = REORG_AND_RETARGET_CHAIN.preRetargetChain;
// 	const POST_CHAIN = REORG_AND_RETARGET_CHAIN.postRetargetChain;

// 	const orphan = REORG_AND_RETARGET_CHAIN.orphan_437478;
// 	const preHex = PRE_CHAIN.map(header => header.hex);
// 	const pre = utils.concatenateHexStrings(preHex);
// 	const postHex = POST_CHAIN.map(header => header.hex);
// 	const post = utils.concatenateHexStrings(postHex.slice(0, -2));
// 	const postWithOrphan = utils.concatenateHexStrings([post, orphan.hex]);
// 	const lastTwo = POST_CHAIN.slice(-2);
// 	const postWithoutOrphan = utils.concatenateHexStrings([post, lastTwo[0].hex, lastTwo[1].hex]);

// 	before(async () => {
// 		instance = await Relay.new(
// 			REORG_AND_RETARGET_CHAIN.genesis.hex,
// 			REORG_AND_RETARGET_CHAIN.genesis.height,
// 			REORG_AND_RETARGET_CHAIN.oldPeriodStart.digest_le
// 		);
// 		await instance.addHeaders(
// 			REORG_AND_RETARGET_CHAIN.genesis.hex,
// 			pre
// 		);
// 		await instance.addHeadersWithRetarget(
// 			REORG_AND_RETARGET_CHAIN.oldPeriodStart.hex,
// 			preHex.slice(-1)[0],
// 			postWithoutOrphan
// 		);
// 		await instance.addHeadersWithRetarget(
// 			REORG_AND_RETARGET_CHAIN.oldPeriodStart.hex,
// 			preHex.slice(-1)[0],
// 			postWithOrphan
// 		);
// 	});

// 	it('returns false if it found a more recent ancestor', async () => {
// 		const res = await instance.isMostRecentAncestor(
// 			POST_CHAIN[0].digest_le,
// 			POST_CHAIN[3].digest_le,
// 			POST_CHAIN[2].digest_le,
// 			5
// 		);
// 		assert.isFalse(res);
// 	});

// 	it('returns false if it did not find the specified common ancestor within the limit', async () => {
// 		const res = await instance.isMostRecentAncestor(
// 			POST_CHAIN[1].digest_le,
// 			POST_CHAIN[3].digest_le,
// 			POST_CHAIN[2].digest_le,
// 			1
// 		);
// 		assert.isFalse(res);
// 	});

// 	it('returns true if the provided digest is the most recent common ancestor', async () => {
// 		let res = await instance.isMostRecentAncestor(
// 			POST_CHAIN[2].digest_le,
// 			POST_CHAIN[3].digest_le,
// 			POST_CHAIN[2].digest_le,
// 			5
// 		);
// 		assert.isTrue(res);

// 		res = await instance.isMostRecentAncestor(
// 			POST_CHAIN[5].digest_le,
// 			POST_CHAIN[6].digest_le,
// 			orphan.digest_le,
// 			5
// 		);
// 		assert.isTrue(res);
// 	});

// 	it('shortcuts the trivial case (ancestor is left is right)', async () => {
// 		const res = await instance.isMostRecentAncestor(
// 			POST_CHAIN[3].digest_le,
// 			POST_CHAIN[3].digest_le,
// 			POST_CHAIN[3].digest_le,
// 			5
// 		);
// 		assert.isTrue(res);
// 	});
// });
