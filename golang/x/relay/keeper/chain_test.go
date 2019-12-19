package keeper

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
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

func (s *KeeperSuite) TestIsMostRecentCommonAncestor() {
	tv := s.Fixtures.ChainTestCases.IsMostRecentCA
	pre := tv.PreRetargetChain
	post := tv.PostRetargetChain
	postWithOrphan := append(post[len(post)-2:], []btcspv.BitcoinHeader{tv.Orphan}...)

	err := s.Keeper.SetGenesisState(s.Context, tv.Genesis, tv.OldPeriodStart)
	s.SDKNil(err)

	err = s.Keeper.IngestHeaderChain(s.Context, pre)
	s.SDKNil(err)
	err = s.Keeper.IngestDifficultyChange(s.Context, tv.OldPeriodStart.HashLE, post)
	s.SDKNil(err)
	err = s.Keeper.IngestDifficultyChange(s.Context, tv.OldPeriodStart.HashLE, postWithOrphan)
	s.SDKNil(err)
	// err = s.Keeper.IngestHeaderChain(s.Context, postWithoutOrphan)
	// s.Nil(err)
	// err = s.Keeper.IngestHeaderChain(s.Context, postWithOrphan)
	// s.Nil(err)

	// // Not passing
	// isMostRecent := s.Keeper.IsMostRecentCommonAncestor(s.Context, post[2].HashLE, post[3].HashLE, post[2].HashLE, 5)
	// s.Equal(true, isMostRecent)

	// // Passing
	// isMostRecent = s.Keeper.IsMostRecentCommonAncestor(s.Context, post[5].HashLE, post[6].HashLE, tv.Orphan.HashLE, 5)
	// s.Equal(true, isMostRecent)

	// // Passing
	// isMostRecent = s.Keeper.IsMostRecentCommonAncestor(s.Context, post[3].HashLE, post[3].HashLE, post[3].HashLE, 5)
	// s.Equal(true, isMostRecent)

	// // Not passing
	// isMostRecent := s.Keeper.IsMostRecentCommonAncestor(s.Context, post[0].HashLE, post[3].HashLE, post[2].HashLE, 5)
	// s.Equal(false, isMostRecent)

	// // Not passing
	// isMostRecent := s.Keeper.IsMostRecentCommonAncestor(s.Context, post[1].HashLE, post[3].HashLE, post[2].HashLE, 1)
	// s.Equal(false, isMostRecent)
}

// err = s.Keeper.IngestDifficultyChange(s.Context, tv.OldPeriodStart.HashLE, tv.PreRetargetChain)
// s.Nil(err)

func (s *KeeperSuite) TestHeaviestFromAncestor() {

}

// describe('#heaviestFromAncestor', async () => {
// 	const { chain, genesis } = REGULAR_CHAIN;
// 	const headerHex = chain.map(header => header.hex);
// 	const headers = utils.concatenateHexStrings(headerHex.slice(0, 8));
// 	const headersWithMain = utils.concatenateHexStrings([headers, chain[8].hex]);
// 	const headersWithOrphan = utils.concatenateHexStrings(
// 		[headers, REGULAR_CHAIN.orphan_562630.hex]
// 	);

// 	before(async () => {
// 		instance = await Relay.new(
// 			genesis.hex,
// 			genesis.height,
// 			genesis.digest_le
// 		);
// 		await instance.addHeaders(genesis.hex, headersWithMain);
// 		await instance.addHeaders(genesis.hex, headersWithOrphan);
// 	});

// 	it('errors if ancestor is unknown', async () => {
// 		try {
// 			await instance.heaviestFromAncestor(
// 				chain[10].digest_le,
// 				headerHex[3],
// 				headerHex[4]
// 			);
// 			assert(false, 'expected an error');
// 		} catch (e) {
// 			assert.include(e.message, 'Unknown block');
// 		}
// 	});

// 	it('errors if left is unknown', async () => {
// 		try {
// 			await instance.heaviestFromAncestor(
// 				chain[3].digest_le,
// 				chain[10].hex,
// 				headerHex[4]
// 			);
// 			assert(false, 'expected an error');
// 		} catch (e) {
// 			assert.include(e.message, 'Unknown block');
// 		}
// 	});

// 	it('errors if right is unknown', async () => {
// 		try {
// 			await instance.heaviestFromAncestor(
// 				chain[3].digest_le,
// 				headerHex[4],
// 				chain[10].hex
// 			);
// 			assert(false, 'expected an error');
// 		} catch (e) {
// 			assert.include(e.message, 'Unknown block');
// 		}
// 	});

// 	it('errors if either block is below the ancestor', async () => {
// 		try {
// 			await instance.heaviestFromAncestor(
// 				chain[3].digest_le,
// 				headerHex[2],
// 				headerHex[4]
// 			);
// 			assert(false, 'expected an error');
// 		} catch (e) {
// 			assert.include(e.message, 'A descendant height is below the ancestor height');
// 		}

// 		try {
// 			await instance.heaviestFromAncestor(
// 				chain[3].digest_le,
// 				headerHex[4],
// 				headerHex[2]
// 			);
// 			assert(false, 'expected an error');
// 		} catch (e) {
// 			assert.include(e.message, 'A descendant height is below the ancestor height');
// 		}
// 	});

// 	it('returns left if left is heavier', async () => {
// 		const res = await instance.heaviestFromAncestor(
// 			chain[3].digest_le,
// 			headerHex[5],
// 			headerHex[4]
// 		);
// 		assert.equal(res, chain[5].digest_le);
// 	});

// 	it('returns right if right is heavier', async () => {
// 		const res = await instance.heaviestFromAncestor(
// 			chain[3].digest_le,
// 			headerHex[4],
// 			headerHex[5]
// 		);
// 		assert.equal(res, chain[5].digest_le);
// 	});

// 	it('returns left if the weights are equal', async () => {
// 		let res = await instance.heaviestFromAncestor(
// 			chain[3].digest_le,
// 			chain[8].hex,
// 			REGULAR_CHAIN.orphan_562630.hex
// 		);
// 		assert.equal(res, chain[8].digest_le);

// 		res = await instance.heaviestFromAncestor(
// 			chain[3].digest_le,
// 			REGULAR_CHAIN.orphan_562630.hex,
// 			chain[8].hex
// 		);
// 		assert.equal(res, REGULAR_CHAIN.orphan_562630.digest_le);
// 	});
// });
