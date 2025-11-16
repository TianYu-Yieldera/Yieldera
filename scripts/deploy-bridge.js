/**
 * Cross-Chain Bridge Deployment Script
 * Deploys bridge contracts to multiple chains
 */

const hre = require("hardhat");
const { ethers } = require("hardhat");
const fs = require("fs");
const path = require("path");

// Network configurations
const NETWORKS = {
  ethereum: {
    chainId: 1,
    rpcUrl: process.env.ETHEREUM_RPC || "https://eth.llamarpc.com",
    explorer: "https://etherscan.io",
    validators: [
      process.env.VALIDATOR_1 || "0x...",
      process.env.VALIDATOR_2 || "0x...",
      process.env.VALIDATOR_3 || "0x...",
    ],
  },
  arbitrum: {
    chainId: 42161,
    rpcUrl: process.env.ARBITRUM_RPC || "https://arb1.arbitrum.io/rpc",
    explorer: "https://arbiscan.io",
    validators: [
      process.env.VALIDATOR_1 || "0x...",
      process.env.VALIDATOR_2 || "0x...",
      process.env.VALIDATOR_3 || "0x...",
    ],
  },
  optimism: {
    chainId: 10,
    rpcUrl: process.env.OPTIMISM_RPC || "https://mainnet.optimism.io",
    explorer: "https://optimistic.etherscan.io",
    validators: [
      process.env.VALIDATOR_1 || "0x...",
      process.env.VALIDATOR_2 || "0x...",
      process.env.VALIDATOR_3 || "0x...",
    ],
  },
  base: {
    chainId: 8453,
    rpcUrl: process.env.BASE_RPC || "https://mainnet.base.org",
    explorer: "https://basescan.org",
    validators: [
      process.env.VALIDATOR_1 || "0x...",
      process.env.VALIDATOR_2 || "0x...",
      process.env.VALIDATOR_3 || "0x...",
    ],
  },
  // Testnets
  sepolia: {
    chainId: 11155111,
    rpcUrl: process.env.SEPOLIA_RPC || "https://rpc.sepolia.org",
    explorer: "https://sepolia.etherscan.io",
    validators: [
      "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1", // Test validator 1
      "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", // Test validator 2
      "0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2", // Test validator 3
    ],
  },
  arbitrumSepolia: {
    chainId: 421614,
    rpcUrl: process.env.ARBITRUM_SEPOLIA_RPC || "https://sepolia-rollup.arbitrum.io/rpc",
    explorer: "https://sepolia.arbiscan.io",
    validators: [
      "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1",
      "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4",
      "0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2",
    ],
  },
};

// Token configurations
const TOKENS = {
  USDC: {
    ethereum: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    arbitrum: "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
    optimism: "0x0b2C639c533813f4Aa9D7837CAf62653d097Ff85",
    base: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
    sepolia: "0x...", // Deploy mock USDC on testnet
    arbitrumSepolia: "0x...", // Deploy mock USDC on testnet
  },
  USDT: {
    ethereum: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
    arbitrum: "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
    optimism: "0x94b008aA00579c1307B0EF2c499aD98a8ce58e58",
    base: "0xfde4C96c8593536E31F229EA8f37b2ADa2699bb2",
  },
  WETH: {
    ethereum: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
    arbitrum: "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
    optimism: "0x4200000000000000000000000000000000000006",
    base: "0x4200000000000000000000000000000000000006",
  },
};

async function deployBridge(networkName) {
  const network = NETWORKS[networkName];
  if (!network) {
    throw new Error(`Network ${networkName} not configured`);
  }

  console.log(`\nDeploying CrossChainBridge to ${networkName}...`);
  console.log(`Chain ID: ${network.chainId}`);
  console.log(`RPC: ${network.rpcUrl}`);

  // Get deployer account
  const [deployer] = await ethers.getSigners();
  console.log(`Deployer: ${deployer.address}`);

  // Deploy treasury (for fees)
  const treasuryAddress = deployer.address; // Use deployer as treasury for now

  // Deploy bridge contract
  const CrossChainBridge = await ethers.getContractFactory("CrossChainBridge");
  const bridge = await CrossChainBridge.deploy(network.chainId, treasuryAddress);
  await bridge.waitForDeployment();

  const bridgeAddress = await bridge.getAddress();
  console.log(`CrossChainBridge deployed to: ${bridgeAddress}`);

  // Setup validators
  console.log("\nSetting up validators...");
  for (const validator of network.validators) {
    if (validator && validator !== "0x...") {
      await bridge.grantRole(await bridge.VALIDATOR_ROLE(), validator);
      console.log(`Added validator: ${validator}`);
    }
  }

  // Setup supported tokens
  console.log("\nSetting up supported tokens...");
  for (const [tokenSymbol, addresses] of Object.entries(TOKENS)) {
    const tokenAddress = addresses[networkName];
    if (tokenAddress && tokenAddress !== "0x...") {
      // Set token as supported with 1M limit
      const limit = ethers.parseUnits("1000000", 6); // Assuming 6 decimals for USDC/USDT
      await bridge.setSupportedToken(tokenAddress, true, limit);
      console.log(`Added ${tokenSymbol}: ${tokenAddress}`);

      // Set token mappings for other chains
      for (const [otherNetwork, otherAddress] of Object.entries(addresses)) {
        if (otherNetwork !== networkName && otherAddress !== "0x...") {
          const otherChainId = NETWORKS[otherNetwork].chainId;
          await bridge.setTokenMapping(tokenAddress, otherChainId, otherAddress);
          console.log(`  Mapped to ${otherNetwork} (${otherChainId}): ${otherAddress}`);
        }
      }
    }
  }

  // Save deployment info
  const deploymentInfo = {
    network: networkName,
    chainId: network.chainId,
    bridge: bridgeAddress,
    treasury: treasuryAddress,
    validators: network.validators.filter(v => v !== "0x..."),
    tokens: Object.entries(TOKENS).reduce((acc, [symbol, addresses]) => {
      if (addresses[networkName] && addresses[networkName] !== "0x...") {
        acc[symbol] = addresses[networkName];
      }
      return acc;
    }, {}),
    deployedAt: new Date().toISOString(),
    deployer: deployer.address,
  };

  const deploymentsDir = path.join(__dirname, "../deployments");
  if (!fs.existsSync(deploymentsDir)) {
    fs.mkdirSync(deploymentsDir, { recursive: true });
  }

  const deploymentFile = path.join(deploymentsDir, `bridge-${networkName}.json`);
  fs.writeFileSync(deploymentFile, JSON.stringify(deploymentInfo, null, 2));

  console.log(`\nDeployment info saved to: ${deploymentFile}`);
  console.log(`Bridge contract: ${bridgeAddress}`);
  console.log(`Explorer: ${network.explorer}/address/${bridgeAddress}`);

  return deploymentInfo;
}

async function deployAllBridges() {
  const deployments = {};

  // Deploy to testnets first for testing
  const testnetNetworks = ["sepolia", "arbitrumSepolia"];

  for (const network of testnetNetworks) {
    try {
      const deployment = await deployBridge(network);
      deployments[network] = deployment;
    } catch (error) {
      console.error(`Failed to deploy to ${network}:`, error.message);
    }
  }

  // Save all deployments
  const allDeploymentsFile = path.join(__dirname, "../deployments/bridges-all.json");
  fs.writeFileSync(allDeploymentsFile, JSON.stringify(deployments, null, 2));

  console.log("\n========================================");
  console.log("Bridge Deployment Complete!");
  console.log("========================================");
  console.log("\nDeployed bridges:");
  for (const [network, deployment] of Object.entries(deployments)) {
    console.log(`  ${network}: ${deployment.bridge}`);
  }
}

// Main execution
async function main() {
  const networkArg = process.argv[2];

  if (networkArg === "all") {
    await deployAllBridges();
  } else if (networkArg) {
    await deployBridge(networkArg);
  } else {
    console.log("Usage: npx hardhat run scripts/deploy-bridge.js [network|all]");
    console.log("Available networks:", Object.keys(NETWORKS).join(", "));
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });