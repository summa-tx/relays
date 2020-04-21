/* eslint-disable */
require('dotenv').config();

const Kit = require('@celo/contractkit')

const HDWalletProvider = require('truffle-hdwallet-provider');
const infuraKey = process.env.SUMMA_RELAY_INFURA_KEY;
const mnemonic = process.env.MNEMONIC;


const ropsten = {
  provider: () => new HDWalletProvider(mnemonic, `https://ropsten.infura.io/v3/${infuraKey}`),
  network_id: 3,
  gas: 5500000,
  confirmations: 2,
  timeoutBlocks: 200
}

const kovan = {
  provider: () => new HDWalletProvider(mnemonic, `https://kovan.infura.io/v3/${infuraKey}`),
  network_id: 42,
  gas: 5500000,
  confirmations: 2,
  timeoutBlocks: 200
}

const alfajores = {
  provider: () => {
    const provider = new HDWalletProvider(mnemonic, 'http://127.0.0.1:9999'); // sinkhole any requests
    // slip44
    const celoBIP44 = "m/44'/52752'/0'/0/0";
    const hdkey = provider.hdwallet.derivePath(celoBIP44);
    // Get the privkey and hand it to the kit
    const privkey = hdkey._hdkey.privateKey.toString('hex');
    const kit = Kit.newKit('https://alfajores-forno.celo-testnet.org');
    kit.addAccount(privkey);
    return kit.web3.currentProvider;
  },
  network_id: 44786,
  gas: 5500000,
  confirmations: 2,
  timeoutBlocks: 200
}

module.exports = {
  api_keys: {
    etherscan: process.env.ETHERSCAN_KEY
  },
  plugins: [
    'solidity-coverage',
    'truffle-plugin-verify'
  ],
  networks: {
    coverage: {
      host: "localhost",
      network_id: "*",
      port: 8555,
      gas: 0xfffffffffff,
      gasPrice: 0x01
    },

    ropsten: ropsten,
    ropsten_test: ropsten,

    kovan: kovan,
    kovan_test: kovan,

    alfajores: alfajores,
    alfajores_test: alfajores,
  },

  // mocha: {
  // },

  compilers: {
    solc: {
      version: "0.5.10",
      settings: {
        optimizer: {
          enabled: true,
          runs: 200
        }
      }
    }
  }
};
