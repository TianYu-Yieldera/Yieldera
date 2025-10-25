import "@nomicfoundation/hardhat-toolbox";

/** @type import('hardhat/config').HardhatUserConfig */
export default {
  solidity: {
    version: "0.8.20",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  paths: {
    sources: "./contracts",
    tests: "./contracts/test",
    cache: "./cache",
    artifacts: "./artifacts"
  },
  networks: {
    hardhat: {
      chainId: 1337,
      allowUnlimitedContractSize: true,
    },
  },
  mocha: {
    timeout: 40000
  }
};
