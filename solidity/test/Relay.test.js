/* global artifacts contract describe before it assert web3 */
const utils = require('./utils.js');
const REGULAR_CHAIN = require('./headers.json');
const RETARGET_CHAIN = require('./headersWithRetarget.json');
const REORG_AND_RETARGET_CHAIN = require('./headersReorgAndRetarget.json');

const Relay = artifacts.require('Relay');

contract('Relay', async () => {
  let instance;

  describe('#constructor', async () => {
    /* eslint-disable-next-line camelcase */
    const { genesis, orphan_562630 } = REGULAR_CHAIN;

    before(async () => {
      instance = await Relay.new(
        genesis.hex,
        genesis.height,
        orphan_562630.digest_le,
      );
    });

    it('errors if the caller is being an idiot', async () => {
      try {
        await Relay.new(
          '0x00',
          genesis.height,
          genesis.digest_le,
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Stop being dumb');
      }
    });

    it('errors if the period start is in wrong byte order', async () => {
      try {
        await Relay.new(
          genesis.hex,
          genesis.height,
          orphan_562630.digest
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Hint: wrong byte order?');
      }
    });

    it('stores genesis block info', async () => {
      let res = await instance.getRelayGenesis.call();
      assert.equal(res, genesis.digest_le);

      res = await instance.getBestKnownDigest.call();
      assert.equal(res, genesis.digest_le);

      res = await instance.getLastReorgCommonAncestor.call();
      assert.equal(res, genesis.digest_le);

      res = await instance.findAncestor.call(genesis.digest_le, 0);
      assert.equal(res, genesis.digest_le);

      res = await instance.findHeight.call(genesis.digest_le);
      assert(res.eqn(genesis.height));
    });
  });

  describe('#addHeaders', async () => {
    /* eslint-disable-next-line camelcase */
    const { chain, genesis, orphan_562630 } = REGULAR_CHAIN;
    const headerHex = chain.map(header => header.hex);

    const headers = utils.concatenateHexStrings(headerHex.slice(0, 6));

    before(async () => {
      instance = await Relay.new(
        genesis.hex,
        genesis.height,
        orphan_562630.digest_le
      );
    });

    it('errors if the anchor is unknown', async () => {
      try {
        await instance.addHeaders('0x00', headers);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Unknown block');
      }
    });

    it('errors if it encounters a retarget on an external call', async () => {
      try {
        const badHeaders = '0x0000002073bd2184edd9c4fc76642ea6754ee40136970efc10c4190000000000000000000296ef123ea96da5cf695f22bf7d94be87d49db1ad7ac371ac43c4da4161c8c216349c5ba11928170d38782b0000002073bd2184edd9c4fc76642ea6754ee40136970efc10c4190000000000000000005af53b865c27c6e9b5e5db4c3ea8e024f8329178a79ddb39f7727ea2fe6e6825d1349c5ba1192817e2d951590000002073bd2184edd9c4fc76642ea6754ee40136970efc10c419000000000000000000c63a8848a448a43c9e4402bd893f701cd11856e14cbbe026699e8fdc445b35a8d93c9c5ba1192817b945dc6c00000020f402c0b551b944665332466753f1eebb846a64ef24c71700000000000000000033fc68e070964e908d961cd11033896fa6c9b8b76f64a2db7ea928afa7e304257d3f9c5ba11928176164145d0000ff3f63d40efa46403afd71a254b54f2b495b7b0164991c2d22000000000000000000f046dc1b71560b7d0786cfbdb25ae320bd9644c98d5c7c77bf9df05cbe96212758419c5ba1192817a2bb2caa00000020e2d4f0edd5edd80bdcb880535443747c6b22b48fb6200d0000000000000000001d3799aa3eb8d18916f46bf2cf807cb89a9b1b4c56c3f2693711bf1064d9a32435429c5ba1192817752e49ae0000002022dba41dff28b337ee3463bf1ab1acf0e57443e0f7ab1d000000000000000000c3aadcc8def003ecbd1ba514592a18baddddcd3a287ccf74f584b04c5c10044e97479c5ba1192817c341f595';
        await instance.addHeaders(genesis.hex, badHeaders);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Unexpected retarget on external call');
      }
    });

    it('errors if the header array is not a multiple of 80 bytes', async () => {
      try {
        // 3 extra bytes on the end
        const badHeaders = headers.substring(0, 8 + 5 * 160);
        await instance.addHeaders(genesis.hex, badHeaders);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Header array length must be divisible by 80');
      }
    });

    it('errors if a header work is too low', async () => {
      try {
        const badHeaders = `${headers}${'00'.repeat(80)}`;
        await instance.addHeaders(genesis.hex, badHeaders);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Header work is insufficient');
      }
    });

    it('errors if the target changes mid-chain', async () => {
      try {
        const badHeaders = utils.concatenateHexStrings([headers, REGULAR_CHAIN.badHeader.hex]);
        await instance.addHeaders(genesis.hex, badHeaders);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Target changed unexpectedly');
      }
    });

    it('errors if a prevhash link is broken', async () => {
      try {
        const badHeaders = utils.concatenateHexStrings([headers, chain[15].hex]);
        await instance.addHeaders(genesis.hex, badHeaders);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Headers do not form a consistent chain');
      }
    });

    it('appends new links to the chain and fires an event', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number;

      await instance.addHeaders(genesis.hex, headers);

      const res = await instance.findHeight.call(chain[0].digest_le);
      assert(res.eqn(genesis.height + 1));

      const eventList = await instance.getPastEvents(
        'Extension',
        { fromBlock: blockNumber, toBlock: 'latest' }
      );
      /* eslint-disable-next-line no-underscore-dangle */
      assert.equal(eventList[0].returnValues._last, chain[5].digest_le);
    });

    it('skips some validation steps for known blocks', async () => {
      const oneMoreHeader = utils.concatenateHexStrings([headers, headerHex[6]]);
      await instance.addHeaders(genesis.hex, oneMoreHeader);
    });
  });

  describe('#addHeadersWithRetarget', async () => {
    const { chain } = RETARGET_CHAIN;
    const headerHex = chain.map(header => header.hex);
    const genesis = chain[1];

    const firstHeader = RETARGET_CHAIN.oldPeriodStart;
    const lastHeader = chain[8];
    const preChange = utils.concatenateHexStrings(headerHex.slice(2, 9));
    const headers = utils.concatenateHexStrings(headerHex.slice(9, 15));

    // let btcutils

    before(async () => {
      // btcutils = await BTCUtils.new()
      instance = await Relay.new(
        genesis.hex,
        genesis.height,
        firstHeader.digest_le
      );
      await instance.addHeaders(genesis.hex, preChange);
    });

    it('errors if the old period start header is unknown', async () => {
      try {
        await instance.addHeadersWithRetarget('0x00', lastHeader.hex, headers);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Unknown block');
      }
    });

    it('errors if the old period end header is unknown', async () => {
      try {
        await instance.addHeadersWithRetarget(firstHeader.hex, chain[15].hex, headers);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Unknown block');
      }
    });

    it('errors if the provided last header does not match records', async () => {
      try {
        await instance.addHeadersWithRetarget(firstHeader.hex, firstHeader.hex, headers);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Must provide the last header of the closing difficulty period');
      }
    });

    it('errors if the start and end headers are not exactly 2015 blocks apart', async () => {
      try {
        await instance.addHeadersWithRetarget(lastHeader.hex, lastHeader.hex, headers);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Must provide exactly 1 difficulty period');
      }
    });

    it('errors if the retarget is performed incorrectly', async () => {
      const tmpInstance = await Relay.new(
        genesis.hex,
        lastHeader.height, // This is a lie
        firstHeader.digest_le
      );
      try {
        await tmpInstance.addHeadersWithRetarget(firstHeader.hex, genesis.hex, headers);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Invalid retarget provided');
      }
    });

    it('appends new links to the chain', async () => {
      await instance.addHeadersWithRetarget(
        firstHeader.hex,
        lastHeader.hex,
        headers
      );

      const res = await instance.findHeight.call(chain[10].digest_le);
      assert(res.eqn(lastHeader.height + 2));
    });
  });

  describe('#findHeight', async () => {
    const { chain, genesis } = REGULAR_CHAIN;
    const headerHex = chain.map(header => header.hex);
    const headers = utils.concatenateHexStrings(headerHex.slice(0, 6));

    before(async () => {
      instance = await Relay.new(
        genesis.hex,
        genesis.height,
        genesis.digest_le
      );
      await instance.addHeaders(genesis.hex, headers);
    });

    it('errors on unknown blocks', async () => {
      try {
        await instance.findHeight(`0x${'00'.repeat(32)}`);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Unknown block');
      }
    });

    it('Finds height of known blocks', async () => {
      for (let i; i < chain.length; i += 1) {
        /* eslint-disable-next-line camelcase */
        const { digest_le, height } = chain[i];
        /* eslint-disable-next-line no-await-in-loop */
        const res = await instance.findHeight(digest_le);
        assert(res.eqn(height), `incorrect height returned ${height}`);
      }
    });
  });

  describe('#findAncestor', async () => {
    const { chain, genesis } = REGULAR_CHAIN;
    const headerHex = chain.map(header => header.hex);
    const headers = utils.concatenateHexStrings(headerHex.slice(0, 6));

    before(async () => {
      instance = await Relay.new(
        genesis.hex,
        genesis.height,
        genesis.digest_le
      );
      await instance.addHeaders(genesis.hex, headers);
    });

    it('errors on unknown blocks', async () => {
      try {
        await instance.findAncestor(`0x${'00'.repeat(32)}`, 3);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Unknown ancestor');
      }
    });

    it('Finds known ancestors based on on offsets', async () => {
      for (let i; i < chain.length; i += 1) {
        /* eslint-disable-next-line camelcase */
        const { digest_le } = chain[i];
        /* eslint-disable-next-line no-await-in-loop */
        let res = await instance.findAncestor(digest_le, 0);
        assert.equal(res, digest_le);
        if (i > 0) {
          /* eslint-disable-next-line no-await-in-loop */
          res = await instance.findAncestor(digest_le, 1);
          assert.equal(res, chain[i - 1].digest_le);
        }
      }
    });
  });

  describe('#isAncestor', async () => {
    const { chain, genesis } = REGULAR_CHAIN;
    const headerHex = chain.map(header => header.hex);
    const headers = utils.concatenateHexStrings(headerHex.slice(0, 6));

    before(async () => {
      instance = await Relay.new(
        genesis.hex,
        genesis.height,
        genesis.digest_le
      );
      await instance.addHeaders(genesis.hex, headers);
    });

    it('returns false if it exceeds the limit', async () => {
      const res = await instance.isAncestor.call(genesis.digest_le, chain[3].digest_le, 1);
      assert.isFalse(res);
    });

    it('finds the ancestor if within the limit', async () => {
      const res = await instance.isAncestor.call(genesis.digest_le, chain[3].digest_le, 5);
      assert.isTrue(res);
    });
  });

  describe('#heaviestFromAncestor', async () => {
    const { chain, genesis } = REGULAR_CHAIN;
    const headerHex = chain.map(header => header.hex);
    const headers = utils.concatenateHexStrings(headerHex.slice(0, 8));
    const headersWithMain = utils.concatenateHexStrings([headers, chain[8].hex]);
    const headersWithOrphan = utils.concatenateHexStrings(
      [headers, REGULAR_CHAIN.orphan_562630.hex]
    );

    before(async () => {
      instance = await Relay.new(
        genesis.hex,
        genesis.height,
        genesis.digest_le
      );
      await instance.addHeaders(genesis.hex, headersWithMain);
      await instance.addHeaders(genesis.hex, headersWithOrphan);
    });

    it('errors if ancestor is unknown', async () => {
      try {
        await instance.heaviestFromAncestor(
          chain[10].digest_le,
          headerHex[3],
          headerHex[4]
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Unknown block');
      }
    });

    it('errors if left is unknown', async () => {
      try {
        await instance.heaviestFromAncestor(
          chain[3].digest_le,
          chain[10].hex,
          headerHex[4]
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Unknown block');
      }
    });

    it('errors if right is unknown', async () => {
      try {
        await instance.heaviestFromAncestor(
          chain[3].digest_le,
          headerHex[4],
          chain[10].hex
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Unknown block');
      }
    });

    it('errors if either block is below the ancestor', async () => {
      try {
        await instance.heaviestFromAncestor(
          chain[3].digest_le,
          headerHex[2],
          headerHex[4]
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'A descendant height is below the ancestor height');
      }

      try {
        await instance.heaviestFromAncestor(
          chain[3].digest_le,
          headerHex[4],
          headerHex[2]
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'A descendant height is below the ancestor height');
      }
    });

    it('returns left if left is heavier', async () => {
      const res = await instance.heaviestFromAncestor(
        chain[3].digest_le,
        headerHex[5],
        headerHex[4]
      );
      assert.equal(res, chain[5].digest_le);
    });

    it('returns right if right is heavier', async () => {
      const res = await instance.heaviestFromAncestor(
        chain[3].digest_le,
        headerHex[4],
        headerHex[5]
      );
      assert.equal(res, chain[5].digest_le);
    });

    it('returns left if the weights are equal', async () => {
      let res = await instance.heaviestFromAncestor(
        chain[3].digest_le,
        chain[8].hex,
        REGULAR_CHAIN.orphan_562630.hex
      );
      assert.equal(res, chain[8].digest_le);

      res = await instance.heaviestFromAncestor(
        chain[3].digest_le,
        REGULAR_CHAIN.orphan_562630.hex,
        chain[8].hex
      );
      assert.equal(res, REGULAR_CHAIN.orphan_562630.digest_le);
    });
  });

  describe('#heaviestFromAncestor (with retarget)', async () => {
    const PRE_CHAIN = REORG_AND_RETARGET_CHAIN.preRetargetChain;
    const POST_CHAIN = REORG_AND_RETARGET_CHAIN.postRetargetChain;

    const orphan = REORG_AND_RETARGET_CHAIN.orphan_437478;
    const pre = utils.concatenateHeadersHexes(PRE_CHAIN)
    const post = utils.concatenateHeadersHexes(POST_CHAIN)
    const shortPost = utils.concatenateHeadersHexes(POST_CHAIN.slice(0, POST_CHAIN.length - 2))
    const postWithOrphan = utils.concatenateHexStrings([shortPost, orphan.hex]);

    before(async () => {
      instance = await Relay.new(
        REORG_AND_RETARGET_CHAIN.genesis.hex,
        REORG_AND_RETARGET_CHAIN.genesis.height,
        REORG_AND_RETARGET_CHAIN.oldPeriodStart.digest_le
      );
      await instance.addHeaders(
        REORG_AND_RETARGET_CHAIN.genesis.hex,
        pre
      );
      await instance.addHeadersWithRetarget(
        REORG_AND_RETARGET_CHAIN.oldPeriodStart.hex,
        PRE_CHAIN[PRE_CHAIN.length - 1].hex,
        post
      );
      await instance.addHeadersWithRetarget(
        REORG_AND_RETARGET_CHAIN.oldPeriodStart.hex,
        PRE_CHAIN[PRE_CHAIN.length - 1].hex,
        postWithOrphan
      );
    });

    it('handles descendants in different difficulty periods', async () => {
      let res = await instance.heaviestFromAncestor.call(
        REORG_AND_RETARGET_CHAIN.genesis.digest_le,
        orphan.hex,
        PRE_CHAIN[3].hex
      );
      assert.equal(res, orphan.digest_le);

      res = await instance.heaviestFromAncestor.call(
        REORG_AND_RETARGET_CHAIN.genesis.digest_le,
        PRE_CHAIN[3].hex,
        orphan.hex
      );
      assert.equal(res, orphan.digest_le);
    });

    it('handles descendants when both are in a new difficulty period', async () => {
      let res = await instance.heaviestFromAncestor.call(
        REORG_AND_RETARGET_CHAIN.genesis.digest_le,
        orphan.hex,
        POST_CHAIN[3].hex
      );
      assert.equal(res, orphan.digest_le);

      res = await instance.heaviestFromAncestor.call(
        REORG_AND_RETARGET_CHAIN.genesis.digest_le,
        POST_CHAIN[3].hex,
        orphan.hex
      );
      assert.equal(res, orphan.digest_le);
    });
  });

  describe('#isMostRecentAncestor', async () => {
    const PRE_CHAIN = REORG_AND_RETARGET_CHAIN.preRetargetChain;
    const POST_CHAIN = REORG_AND_RETARGET_CHAIN.postRetargetChain;

    const orphan = REORG_AND_RETARGET_CHAIN.orphan_437478;
    const pre = utils.concatenateHeadersHexes(PRE_CHAIN)
    const post = utils.concatenateHeadersHexes(POST_CHAIN)
    const shortPost = utils.concatenateHeadersHexes(POST_CHAIN.slice(0, POST_CHAIN.length - 2))
    const postWithOrphan = utils.concatenateHexStrings([shortPost, orphan.hex]);

    before(async () => {
      instance = await Relay.new(
        REORG_AND_RETARGET_CHAIN.genesis.hex,
        REORG_AND_RETARGET_CHAIN.genesis.height,
        REORG_AND_RETARGET_CHAIN.oldPeriodStart.digest_le
      );
      await instance.addHeaders(
        REORG_AND_RETARGET_CHAIN.genesis.hex,
        pre
      );
      await instance.addHeadersWithRetarget(
        REORG_AND_RETARGET_CHAIN.oldPeriodStart.hex,
        PRE_CHAIN[PRE_CHAIN.length - 1].hex,
        post
      );
      await instance.addHeadersWithRetarget(
        REORG_AND_RETARGET_CHAIN.oldPeriodStart.hex,
        PRE_CHAIN[PRE_CHAIN.length - 1].hex,
        postWithOrphan
      );
    });

    it('returns false if it found a more recent ancestor', async () => {
      const res = await instance.isMostRecentAncestor(
        POST_CHAIN[0].digest_le,
        POST_CHAIN[3].digest_le,
        POST_CHAIN[2].digest_le,
        5
      );
      assert.isFalse(res);
    });

    it('returns false if it did not find the specified common ancestor within the limit', async () => {
      const res = await instance.isMostRecentAncestor(
        POST_CHAIN[1].digest_le,
        POST_CHAIN[3].digest_le,
        POST_CHAIN[2].digest_le,
        1
      );
      assert.isFalse(res);
    });

    it('returns true if the provided digest is the most recent common ancestor', async () => {
      let res = await instance.isMostRecentAncestor(
        POST_CHAIN[2].digest_le,
        POST_CHAIN[3].digest_le,
        POST_CHAIN[2].digest_le,
        5
      );
      assert.isTrue(res);

      res = await instance.isMostRecentAncestor(
        POST_CHAIN[5].digest_le,
        POST_CHAIN[6].digest_le,
        orphan.digest_le,
        5
      );
      assert.isTrue(res);
    });

    it('shortcuts the trivial case (ancestor is left is right)', async () => {
      const res = await instance.isMostRecentAncestor(
        POST_CHAIN[3].digest_le,
        POST_CHAIN[3].digest_le,
        POST_CHAIN[3].digest_le,
        5
      );
      assert.isTrue(res);
    });
  });

  describe('#markNewHeaviest', async () => {
    const PRE_CHAIN = REORG_AND_RETARGET_CHAIN.preRetargetChain;
    const POST_CHAIN = REORG_AND_RETARGET_CHAIN.postRetargetChain;

    const orphan = REORG_AND_RETARGET_CHAIN.orphan_437478;
    const pre = utils.concatenateHeadersHexes(PRE_CHAIN)
    const post = utils.concatenateHeadersHexes(POST_CHAIN)
    const shortPost = utils.concatenateHeadersHexes(POST_CHAIN.slice(0, POST_CHAIN.length - 2))
    const postWithOrphan = utils.concatenateHexStrings([shortPost, orphan.hex]);

    before(async () => {
      instance = await Relay.new(
        REORG_AND_RETARGET_CHAIN.genesis.hex,
        REORG_AND_RETARGET_CHAIN.genesis.height,
        REORG_AND_RETARGET_CHAIN.oldPeriodStart.digest_le
      );
      await instance.addHeaders(
        REORG_AND_RETARGET_CHAIN.genesis.hex,
        pre
      );
      await instance.addHeadersWithRetarget(
        REORG_AND_RETARGET_CHAIN.oldPeriodStart.hex,
        PRE_CHAIN[PRE_CHAIN.length - 1].hex,
        post
      );
      await instance.addHeadersWithRetarget(
        REORG_AND_RETARGET_CHAIN.oldPeriodStart.hex,
        PRE_CHAIN[PRE_CHAIN.length - 1].hex,
        postWithOrphan
      );
    });

    it('errors if the passed in best is not the best known', async () => {
      try {
        await instance.markNewHeaviest(
          REORG_AND_RETARGET_CHAIN.oldPeriodStart.digest_le,
          REORG_AND_RETARGET_CHAIN.oldPeriodStart.hex,
          REORG_AND_RETARGET_CHAIN.oldPeriodStart.hex,
          10
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Passed in best is not best known');
      }
    });

    it('errors if the new best is not already known', async () => {
      try {
        await instance.markNewHeaviest(
          REORG_AND_RETARGET_CHAIN.genesis.digest_le,
          REORG_AND_RETARGET_CHAIN.genesis.hex,
          `0x${'99'.repeat(80)}`,
          10
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'New best is unknown');
      }
    });

    it('errors if the ancestor is not the heaviest common ancestor', async () => {
      await instance.markNewHeaviest(
        REORG_AND_RETARGET_CHAIN.genesis.digest_le,
        REORG_AND_RETARGET_CHAIN.genesis.hex,
        PRE_CHAIN[0].hex,
        10
      );
      try {
        await instance.markNewHeaviest(
          REORG_AND_RETARGET_CHAIN.genesis.digest_le,
          PRE_CHAIN[0].hex,
          PRE_CHAIN[1].hex,
          10
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Ancestor must be heaviest common ancestor');
      }
    });

    it('updates the best known and emits a reorg event', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number;
      await instance.markNewHeaviest(
        PRE_CHAIN[0].digest_le,
        PRE_CHAIN[0].hex,
        orphan.hex,
        20
      );
      const eventList = await instance.getPastEvents(
        'Reorg',
        { fromBlock: blockNumber, toBlock: 'latest' }
      );
      /* eslint-disable no-underscore-dangle */
      assert.equal(eventList[0].returnValues._to, orphan.digest_le);
      assert.equal(eventList[0].returnValues._from, PRE_CHAIN[0].digest_le);
      assert.equal(eventList[0].returnValues._gcd, PRE_CHAIN[0].digest_le);
      /* eslint-enable no-underscore-dangle */
    });

    it('errors if the new best hash is not better', async () => {
      try {
        await instance.markNewHeaviest(
          POST_CHAIN.slice(-3)[0].digest_le, // the main chain before the split
          orphan.hex,
          POST_CHAIN.slice(-2)[0].hex, // the main chain competing with the split
          10
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'New best hash does not have more work than previous');
      }
    });
  });
});
