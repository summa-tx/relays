require('dotenv').config();
const BN = require('bn.js');

const ID_SPACE_SIZE = new BN('2', 10).pow(new BN('32', 10));

const truffleConf = require('../truffle-config');

const bitcoinMain = {
  genesis: '0x00000020d208b5e50a8d3bd3a87f7a238e3f196621d0f9ffb5f302000000000000000000ee3af51ad3643a8a109935b45d9ca32b1003cda41df39dd75a17e13ba13aff4211aa585d39301c174a8ead73',
  height: 590588,
  epochStart: '0x704de08dc5329269011b878835be108a8202a93a0a2a1c000000000000000000',
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
