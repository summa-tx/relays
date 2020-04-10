require('dotenv').config();
const BN = require('bn.js');

const ID_SPACE_SIZE = new BN('2', 10).pow(new BN('32', 10));

const truffleConf = require('../truffle-config');

const bitcoinMain = {
  genesis: '0x000040202842774747733a4863b6bbb7b4cfb66baa9287d5ce0d13000000000000000000df550e01d02ee37fce8dd2fbf919a47b8b65684bcb48d4da699078916da2f7decbc7905ebc2013178f58d533',
  height: 625332,
  epochStart: '0x6e9d58fb0ab8d0181b1c9e54614f80b64004c2e04da310000000000000000000',
};

const bitcoinTest = {
  genesis: '0x0000ff3ffc663e3a0b12b4cc2c05a425bdaf51922ce090acd8fa3a8a180300000000000084080b23fc40476d284da49fedaea9f7cee3aba33a8bad1347fa54740a29f02752b4c45dfcff031a279c2b3a',
  height: 1607272,
  epochStart: '0x84a9ec3b82556297ea36d1377901ecaef0bb5a5cf683f9f05103000000000000'
};

module.exports = {
  ropsten: {
    network_id: truffleConf.networks.ropsten.network_id,
    bitcoin: bitcoinMain,
    firstID: ID_SPACE_SIZE.muln(truffleConf.networks.ropsten.network_id)
  },
  ropsten_test: {
    network_id: truffleConf.networks.ropsten_test.network_id,
    bitcoin: bitcoinTest,
    firstID: ID_SPACE_SIZE.muln(truffleConf.networks.ropsten_test.network_id + 0x800000)
  },
  kovan: {
    network_id: truffleConf.networks.kovan.network_id,
    bitcoin: bitcoinMain,
    firstID: ID_SPACE_SIZE.muln(truffleConf.networks.kovan.network_id)
  },
  kovan_test: {
    network_id: truffleConf.networks.kovan_test.network_id,
    bitcoin: bitcoinTest,
    firstID: ID_SPACE_SIZE.muln(truffleConf.networks.kovan_test.network_id + 0x800000)
  },
  alfajores: {
    network_id: truffleConf.networks.alfajores_test.network_id,
    bitcoin: bitcoinMain,
    firstID: ID_SPACE_SIZE.muln(truffleConf.networks.alfajores.network_id)
  },
  alfajores_test: {
    network_id: truffleConf.networks.alfajores_test.network_id,
    bitcoin: bitcoinTest,
    firstID: ID_SPACE_SIZE.muln(truffleConf.networks.alfajores_test.network_id + 0x800000)
  },
};
