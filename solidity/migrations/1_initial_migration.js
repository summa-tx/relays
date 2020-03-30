/* global artifacts */
const Migrations = artifacts.require('Migrations');

module.exports = async (deployer) => {
  deployer.deploy(Migrations);
};
