/* global artifacts contract before after describe it assert web3 */
const BN = require('bn.js');
const constants = require('./OnDemandSPVHelpers.json');
const REGULAR_CHAIN = require('./headers.json');

const DummyConsumer = artifacts.require('DummyConsumer');
const DummyOnDemandSPV = artifacts.require('DummyOnDemandSPV');

contract('OnDemandSPV', async (accounts) => {
  let instance;
  let consumer;

  const [deployer, requestOwner, outsideCaller] = accounts;
  const { genesis } = REGULAR_CHAIN;
  const BYTES32_0 = '0x0000000000000000000000000000000000000000000000000000000000000000';

  before(async () => {
    instance = await DummyOnDemandSPV.new(
      genesis.hex,
      genesis.height,
      genesis.digest_le,
      0,
      { from: deployer }
    );
    consumer = await DummyConsumer.new({ from: deployer });
  });

  describe('#cancelRequest', async () => {
    const [sub1Id, sub2Id] = [1, 2];

    before(async () => {
      await instance.requestTest(
        sub1Id,
        `0x${'00'.repeat(36)}`,
        '0x',
        0,
        consumer.address,
        1,
        0,
        { from: requestOwner }
      );
      await instance.requestTest(
        sub2Id,
        `0x${'00'.repeat(36)}`,
        '0x',
        0,
        consumer.address,
        1,
        0,
        { from: requestOwner }
      );
    });

    it('errors if not active', async () => {
      try {
        await instance.cancelRequest(3);
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Request not active');
      }
    });

    it('cannot be cancelled by an outside caller', async () => {
      try {
        await instance.cancelRequest(sub1Id, { from: outsideCaller });
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Can only be cancelled by owner or consumer');
      }
    });

    it('can be canceled by the owner', async () => {
      await instance.cancelRequest(sub1Id, { from: requestOwner });
      const request = await instance.getRequest(sub1Id);
      assert(request[3].eq(new BN('2', 10)));
    });

    it('can be canceled by the consumer', async () => {
      await consumer.cancel(sub2Id, instance.address);
      const request = await instance.getRequest(sub2Id);
      assert(request[3].eq(new BN('2', 10)));
    });
  });

  describe('#getRequest', async () => {
    const sub3Id = 3;

    before(async () => {
      await instance.requestTest(
        sub3Id,
        `0x${'11'.repeat(36)}`,
        '0x',
        100,
        consumer.address,
        0,
        0,
        { from: requestOwner }
      );
    });

    it('retrieves request information', async () => {
      const res = await instance.getRequest(sub3Id);

      // this is the keccak256 of `0x${'11'.repeat(36)}`
      assert.strictEqual(res[0], '0x600e7bfdb8c3cc85df9cd058022100f260c17e7c58603758d2c1ac92c63469a6');
      assert.strictEqual(res[1], BYTES32_0);
      assert(res[2].eq(new BN('100', 10)));
      assert(res[3].eq(new BN('1', 10)));
      assert.strictEqual(res[4], consumer.address);
      assert.strictEqual(res[5], requestOwner);
    });
  });

  describe('#request', async () => {
    const sub4Id = 4;

    it('errors if the outpoint is present and not 36 bytes', async () => {
      try {
        await instance.requestTest(sub4Id, '0xff', '0x', 0, consumer.address, 1, 0, { from: requestOwner });
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Not a valid UTXO');
      }
    });

    it('errors if the spk is valid, but nonstandard', async () => {
      try {
        await instance.requestTest(sub4Id, '0x', '0x0d00000000000000000000000000', 0, consumer.address, 1, 0, { from: requestOwner });
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Not a standard output type');
      }
    });

    it('errors if spends and pays are both 0', async () => {
      try {
        await instance.requestTest(sub4Id, '0x', '0x', 0, consumer.address, 1, 0, { from: requestOwner });
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'No request specified');
      }
    });

    it('stores a request and emits an event', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number;

      await instance.requestTest(sub4Id, `0x${'11'.repeat(36)}`, '0x', 100, consumer.address, 1, 0, { from: requestOwner });
      const res = await instance.getRequest(sub4Id);

      // this is the keccak256 of `0x${'11'.repeat(36)}`
      assert.strictEqual(res[0], '0x600e7bfdb8c3cc85df9cd058022100f260c17e7c58603758d2c1ac92c63469a6');
      assert.strictEqual(res[1], BYTES32_0);
      assert(res[2].eq(new BN('100', 10)));
      assert(res[3].eq(new BN('1', 10)));
      assert.strictEqual(res[4], consumer.address);
      assert.strictEqual(res[5], requestOwner);


      const eventList = await instance.getPastEvents(
        'NewProofRequest',
        { fromBlock: blockNumber, toBlock: 'latest' }
      );
      /* eslint-disable-next-line no-underscore-dangle */
      assert.strictEqual(eventList[0].returnValues._requester, requestOwner);
      /* eslint-disable-next-line no-underscore-dangle */
      assert.strictEqual(parseInt(eventList[0].returnValues._requestID, 10), 4);
    });

    it('stores bytes32(0) for unspecified spends', async () => {
      await instance.requestTest(
        88,
        '0x',
        constants.OP_RETURN_PAYS_1,
        0,
        consumer.address,
        1,
        0,
        { from: requestOwner }
      );
      const res = await instance.getRequest(88);
      assert.strictEqual(res[0], BYTES32_0);
    });
  });

  describe('#provideProof', async () => {
    const sub5Id = 5;

    before(async () => {
      await instance.requestTest(
        sub5Id,
        constants.OP_RETURN_SPENDS_0,
        constants.OP_RETURN_PAYS_1,
        0,
        consumer.address,
        1,
        0,
        { from: requestOwner }
      );
      await instance.setCallResult(true);
    });

    after(async () => {
      await instance.setCallResult(false);
    });

    it('runs succesfully, and sets validatedTxns and latestValidatedTx', async () => {
      assert.isOk(await instance.provideProof(
        constants.OP_RETURN_HEADER,
        constants.OP_RETURN_PROOF,
        constants.OP_RETURN_VERSION,
        constants.OP_RETURN_LOCKTIME,
        constants.OP_RETURN_INDEX,
        '0x0001', // requestIndices
        constants.OP_RETURN_VIN,
        constants.OP_RETURN_VOUT,
        sub5Id
      ));

      let res = await instance.getValidatedTx(constants.OP_RETURN_TX_ID_LE);
      assert.isTrue(res);

      res = await instance.getLatestValidatedTx.call();
      assert.equal(res, constants.OP_RETURN_TX_ID_LE);
    });

    it('shortcuts inclusion validatins for already-seen txns', async () => {
      if (!await instance.getValidatedTx.call(constants.OP_RETURN_TX_ID_LE)) {
        await instance.setValidatedTx(constants.OP_RETURN_TX_ID_LE);
      }
      assert.isOk(await instance.provideProof(
        '0x',
        '0x',
        constants.OP_RETURN_VERSION,
        constants.OP_RETURN_LOCKTIME,
        0,
        '0x0001', // requestIndices
        constants.OP_RETURN_VIN,
        constants.OP_RETURN_VOUT,
        sub5Id
      ));
    });
  });

  describe('#callCallback', async () => {
    const sub6Id = 6;

    async function submitValidProof() {
      return instance.provideProof(
        constants.OP_RETURN_HEADER,
        constants.OP_RETURN_PROOF,
        constants.OP_RETURN_VERSION,
        constants.OP_RETURN_LOCKTIME,
        constants.OP_RETURN_INDEX,
        '0x0001', // requestIndices
        constants.OP_RETURN_VIN,
        constants.OP_RETURN_VOUT,
        sub6Id,
        { gas: 4000000 }
      );
    }

    before(async () => {
      await instance.requestTest(
        sub6Id,
        constants.OP_RETURN_SPENDS_0,
        constants.OP_RETURN_PAYS_1,
        0,
        consumer.address,
        1,
        0,
        { from: requestOwner }
      );
      await instance.setCallResult(true);
    });

    it('calls the consumer with 500,000 gas', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number;

      await submitValidProof();

      const eventList = await consumer.getPastEvents(
        'Consumed',
        { fromBlock: blockNumber, toBlock: 'latest' }
      );
      /* eslint-disable-next-line no-underscore-dangle */
      assert(new BN(eventList[0].returnValues._gasLeft, 10).ltn(500000));
    });

    it('functions even if the remote contract reverts', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number;

      await consumer.setBroken(true);
      assert.isOk(await submitValidProof()); // doesn't revert
      await consumer.setBroken(false);

      // should be no events, because consumer reverted
      const eventList = await consumer.getPastEvents(
        'Consumed',
        { fromBlock: blockNumber, toBlock: 'latest' }
      );
      assert.strictEqual(eventList.length, 0);
    });
  });

  describe('#checkInclusion', async () => {
    it('errors on a bad inclusion proof', async () => {
      try {
        await instance.checkInclusion(
          constants.OP_RETURN_HEADER,
          '0x',
          constants.OP_RETURN_INDEX,
          constants.OP_RETURN_TX_ID_LE,
          1,
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Bad inclusion proof');
      }
    });

    it('errors if isAncestor fails', async () => {
      await instance.setCallResult(false);
      try {
        await instance.checkInclusion(
          constants.OP_RETURN_HEADER,
          constants.OP_RETURN_PROOF,
          constants.OP_RETURN_INDEX,
          constants.OP_RETURN_TX_ID_LE,
          1
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'GCD does not confirm header');
      }
      await instance.setCallResult(true);
    });

    it('errors if insufficinet confirmations', async () => {
      const sub = 838;
      await instance.requestTest(
        sub,
        `0x${'00'.repeat(36)}`,
        '0x',
        0,
        consumer.address,
        240, // very large conf requirement
        0,
        { from: requestOwner }
      );
      try {
        await instance.checkInclusion(
          constants.OP_RETURN_HEADER,
          constants.OP_RETURN_PROOF,
          constants.OP_RETURN_INDEX,
          constants.OP_RETURN_TX_ID_LE,
          sub
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Insufficient confirmations');
      }
      await instance.setCallResult(true);
    });

    it('succeeds', async () => {
      await instance.setCallResult(true);
      const res = await instance.checkInclusion.call(
        constants.OP_RETURN_HEADER,
        constants.OP_RETURN_PROOF,
        constants.OP_RETURN_INDEX,
        constants.OP_RETURN_TX_ID_LE,
        1
      );
      assert.isTrue(res);
    });
  });

  describe('#checkRequests', async () => {
    const sub7Id = 7;
    const sub8Id = 8;
    const sub9Id = 9;

    before(async () => {
      await instance.requestTest( // both
        sub7Id,
        constants.OP_RETURN_SPENDS_0,
        constants.OP_RETURN_PAYS_1,
        0,
        consumer.address,
        1,
        0,
        { from: requestOwner }
      );
      await instance.requestTest( // only spends
        sub8Id,
        constants.OP_RETURN_SPENDS_0,
        '0x',
        0,
        consumer.address,
        1,
        0,
        { from: requestOwner }
      );
      await instance.requestTest( // only pays
        sub9Id,
        '0x',
        constants.OP_RETURN_PAYS_1,
        0,
        consumer.address,
        1,
        0,
        { from: requestOwner }
      );
      await instance.setCallResult(true);
    });

    it('errors if the vin is malformatted', async () => {
      try {
        await instance.checkRequests(
          '0x0001',
          `0x${'01'.repeat(66)}`,
          constants.OP_RETURN_VOUT,
          sub7Id
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Vin is malformatted');
      }
    });

    it('errors if the vout is malformatted', async () => {
      try {
        await instance.checkRequests(
          '0x0001',
          constants.OP_RETURN_VIN,
          `0x${'01'.repeat(66)}`,
          sub7Id
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Vout is malformatted');
      }
    });

    it('errors if the request is not active', async () => {
      try {
        await instance.checkRequests(
          '0x0001',
          constants.OP_RETURN_VIN,
          constants.OP_RETURN_VOUT,
          11 // fragile: not yet active. breaks if we add more cases above this;
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Request is not active');
      }
    });

    it('errors if the specified output does not match the pays request', async () => {
      try {
        await instance.checkRequests(
          '0x0000', // first output instead of second
          constants.OP_RETURN_VIN,
          constants.OP_RETURN_VOUT,
          sub7Id
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Does not match pays request');
      }

      try {
        await instance.checkRequests(
          '0xFF00', // first output instead of second
          constants.OP_RETURN_VIN,
          constants.OP_RETURN_VOUT,
          sub9Id
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Does not match pays request');
      }
    });

    it('errors if the specified output does not match the paysValue request', async () => {
      await instance.requestTest(
        333,
        constants.OP_RETURN_SPENDS_0,
        constants.OP_RETURN_PAYS_1,
        1000,
        consumer.address,
        1,
        0,
        { from: requestOwner }
      );
      try {
        await instance.checkRequests(
          '0x0001',
          constants.OP_RETURN_VIN,
          constants.OP_RETURN_VOUT,
          333
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Does not match value request');
      }
    });

    it('errors if the specified input does not match the spends request', async () => {
      await instance.requestTest(
        444,
        `0x${'33'.repeat(36)}`,
        constants.OP_RETURN_PAYS_1,
        0,
        consumer.address,
        1,
        0,
        { from: requestOwner }
      );
      try {
        await instance.checkRequests(
          '0x0001',
          constants.OP_RETURN_VIN,
          constants.OP_RETURN_VOUT,
          444
        );
        assert(false, 'expected an error');
      } catch (e) {
        assert.include(e.message, 'Does not match spends request');
      }
    });

    it('suceeds', async () => {
      assert.ok(await instance.checkRequests(
        '0x0001',
        constants.OP_RETURN_VIN,
        constants.OP_RETURN_VOUT,
        sub7Id
      ));
      assert.ok(await instance.checkRequests(
        '0x0001',
        constants.OP_RETURN_VIN,
        constants.OP_RETURN_VOUT,
        sub8Id
      ));
      assert.ok(await instance.checkRequests(
        '0x0001',
        constants.OP_RETURN_VIN,
        constants.OP_RETURN_VOUT,
        sub9Id
      ));
    });

    it('allows 0xFF for unchecked', async () => {
      assert.ok(await instance.checkRequests(
        '0x0001',
        constants.OP_RETURN_VIN,
        constants.OP_RETURN_VOUT,
        sub7Id
      ));
      assert.ok(await instance.checkRequests(
        '0x00FF',
        constants.OP_RETURN_VIN,
        constants.OP_RETURN_VOUT,
        sub8Id
      ));
      assert.ok(await instance.checkRequests(
        '0xFF01',
        constants.OP_RETURN_VIN,
        constants.OP_RETURN_VOUT,
        sub9Id
      ));
    });
  });

  describe('#_getConfs', async () => {
    it('should return the number of confirmations', async () => {
      const res = await instance.getConfsTest.call();
      assert(res.eq(new BN('0', 10))); // <- fragile: best and LCA are the same here
    });
  });
});
