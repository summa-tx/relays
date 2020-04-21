/* global artifacts */

const nets = require('./networkInfo');

const Relay = artifacts.require('OnDemandSPV');
const TestnetRelay = artifacts.require('TestnetRelay');

function sleep(milliseconds) {
  return new Promise(resolve => setTimeout(resolve, milliseconds));
}

module.exports = async (deployer, network) => {
  if (['test', 'development', 'soliditycoverage'].includes(network)) {
    // never run deployments on development. We deploy in tests
    return;
  }

  const isBitcoinTestnet = network.includes('_test');
  const contract = isBitcoinTestnet ? TestnetRelay : Relay;

  // dry runs are postfixed with '-fork'
  const strippedNetwork = network.split('-')[0];

  const deployInfo = nets[strippedNetwork];

  const { firstID } = deployInfo;
  const { genesis, height, epochStart } = deployInfo.bitcoin;

  /* eslint-disable */
  console.log('');
  console.log(`network is ${network}`);
  console.log(`First request ID is ${firstID}`);
  console.log('Press Ctrl+C to cancel');
  console.log('');
  await sleep(7500);
  /* eslint-enable */

  deployer.deploy(contract, genesis, height, epochStart, firstID);
};
