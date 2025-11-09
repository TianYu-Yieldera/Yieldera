const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

/**
 * GMX V2 Adapter 部署脚本
 *
 * 支持的网络:
 * - Arbitrum One (主网)
 * - Arbitrum Sepolia (测试网)
 * - 本地测试网 (使用 Mock 合约)
 */

// GMX V2 合约地址配置
const GMX_ADDRESSES = {
  // Arbitrum One (主网)
  arbitrumOne: {
    exchangeRouter: "0x7C68C7866A64FA2160F78EEaE12217FFbf871fa8",
    reader: "0xf60becbba223EEA9495Da3f606753867eC10d139",
    dataStore: "0xFD70de6b91282D8017aA4E741e9Ae325CAb992d8",
    markets: {
      // ETH/USD 市场
      ETH_USD: "0x70d95587d40A2caf56bd97485aB3Eec10Bee6336",
      // BTC/USD 市场
      BTC_USD: "0x47c031236e19d024b42f8AE6780E44A573170703",
      // 更多市场可以在这里添加
    },
    collateral: {
      USDC: "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
      USDT: "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
      WETH: "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
      WBTC: "0x2f2a2543B76A4166549F7aaB2e75Bef0aefC5B0f",
    }
  },

  // Arbitrum Sepolia (测试网)
  arbitrumSepolia: {
    // 注意: GMX V2 在 Sepolia 可能没有官方部署
    // 这里使用占位符，实际部署前需要确认
    exchangeRouter: "0x0000000000000000000000000000000000000000",
    reader: "0x0000000000000000000000000000000000000000",
    dataStore: "0x0000000000000000000000000000000000000000",
    markets: {
      ETH_USD: "0x0000000000000000000000000000000000000000",
    },
    collateral: {
      USDC: "0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d", // Sepolia USDC
    }
  }
};

async function main() {
  const [deployer] = await ethers.getSigners();
  const network = hre.network.name;

  console.log("============================================");
  console.log("GMX V2 Adapter 部署脚本");
  console.log("============================================");
  console.log("Network:", network);
  console.log("Deployer:", deployer.address);
  console.log("Balance:", ethers.formatEther(await ethers.provider.getBalance(deployer.address)), "ETH");
  console.log("============================================\n");

  let exchangeRouter, reader, dataStore;
  let useMockContracts = false;

  // 根据网络选择合约地址
  if (network === "arbitrumOne") {
    console.log("📍 Using Arbitrum One GMX V2 contracts\n");
    const addresses = GMX_ADDRESSES.arbitrumOne;
    exchangeRouter = addresses.exchangeRouter;
    reader = addresses.reader;
    dataStore = addresses.dataStore;
  } else if (network === "arbitrumSepolia") {
    console.log("📍 Using Arbitrum Sepolia GMX V2 contracts\n");
    const addresses = GMX_ADDRESSES.arbitrumSepolia;

    // 检查是否为占位符
    if (addresses.exchangeRouter === "0x0000000000000000000000000000000000000000") {
      console.log("⚠️  Warning: GMX V2 may not be deployed on Sepolia");
      console.log("⚠️  Deploying Mock GMX V2 contracts instead...\n");
      useMockContracts = true;
    } else {
      exchangeRouter = addresses.exchangeRouter;
      reader = addresses.reader;
      dataStore = addresses.dataStore;
    }
  } else {
    console.log("📍 Local network detected - deploying Mock GMX V2 contracts\n");
    useMockContracts = true;
  }

  // 部署或使用 Mock GMX V2 合约
  if (useMockContracts) {
    console.log("🚀 Deploying Mock GMX V2 contracts...");

    const MockExchangeRouter = await ethers.getContractFactory("MockGMXV2ExchangeRouter");
    const mockExchangeRouter = await MockExchangeRouter.deploy();
    await mockExchangeRouter.waitForDeployment();
    exchangeRouter = await mockExchangeRouter.getAddress();
    console.log("✅ MockGMXV2ExchangeRouter:", exchangeRouter);

    const MockReader = await ethers.getContractFactory("MockGMXV2Reader");
    const mockReader = await MockReader.deploy();
    await mockReader.waitForDeployment();
    reader = await mockReader.getAddress();
    console.log("✅ MockGMXV2Reader:", reader);

    const MockDataStore = await ethers.getContractFactory("MockGMXV2DataStore");
    const mockDataStore = await MockDataStore.deploy();
    await mockDataStore.waitForDeployment();
    dataStore = await mockDataStore.getAddress();
    console.log("✅ MockGMXV2DataStore:", dataStore);
    console.log();
  }

  // 部署 GMXV2Adapter
  console.log("🚀 Deploying GMXV2Adapter...");
  const GMXV2Adapter = await ethers.getContractFactory("GMXV2Adapter");
  const adapter = await GMXV2Adapter.deploy(
    exchangeRouter,
    reader,
    dataStore
  );
  await adapter.waitForDeployment();
  const adapterAddress = await adapter.getAddress();

  console.log("✅ GMXV2Adapter deployed to:", adapterAddress);
  console.log();

  // 配置 Adapter
  console.log("⚙️  Configuring GMXV2Adapter...");

  // 添加市场
  const markets = network === "arbitrumOne"
    ? GMX_ADDRESSES.arbitrumOne.markets
    : { ETH_USD: useMockContracts ? "0x70d95587d40A2caf56bd97485aB3Eec10Bee6336" : GMX_ADDRESSES.arbitrumSepolia.markets.ETH_USD };

  for (const [name, address] of Object.entries(markets)) {
    if (address !== "0x0000000000000000000000000000000000000000") {
      console.log(`  Adding market ${name}:`, address);
      const tx = await adapter.addMarket(address);
      await tx.wait();
    }
  }

  // 添加抵押品
  const collateral = network === "arbitrumOne"
    ? GMX_ADDRESSES.arbitrumOne.collateral
    : GMX_ADDRESSES.arbitrumSepolia.collateral;

  for (const [name, address] of Object.entries(collateral)) {
    if (address !== "0x0000000000000000000000000000000000000000") {
      console.log(`  Adding collateral ${name}:`, address);
      const tx = await adapter.addCollateral(address);
      await tx.wait();
    }
  }

  console.log("✅ Configuration complete\n");

  // 保存部署信息
  const deploymentInfo = {
    network: network,
    chainId: (await ethers.provider.getNetwork()).chainId.toString(),
    deployer: deployer.address,
    timestamp: new Date().toISOString(),
    contracts: {
      GMXV2Adapter: adapterAddress,
      GMXExchangeRouter: exchangeRouter,
      GMXReader: reader,
      GMXDataStore: dataStore,
    },
    markets: markets,
    collateral: collateral,
    useMockContracts: useMockContracts,
  };

  const deploymentsDir = path.join(__dirname, "../deployments");
  if (!fs.existsSync(deploymentsDir)) {
    fs.mkdirSync(deploymentsDir, { recursive: true });
  }

  const deploymentFile = path.join(deploymentsDir, `gmx-adapter-${network}.json`);
  fs.writeFileSync(deploymentFile, JSON.stringify(deploymentInfo, null, 2));

  console.log("📄 Deployment info saved to:", deploymentFile);
  console.log();

  // 验证合约 (仅在主网和测试网)
  if (network !== "hardhat" && network !== "localhost" && !useMockContracts) {
    console.log("🔍 Verifying contract on Etherscan...");
    console.log("⏳ Waiting 30 seconds before verification...");
    await new Promise(resolve => setTimeout(resolve, 30000));

    try {
      await hre.run("verify:verify", {
        address: adapterAddress,
        constructorArguments: [exchangeRouter, reader, dataStore],
      });
      console.log("✅ Contract verified successfully");
    } catch (error) {
      console.log("⚠️  Verification failed:", error.message);
      console.log("You can verify manually later with:");
      console.log(`npx hardhat verify --network ${network} ${adapterAddress} ${exchangeRouter} ${reader} ${dataStore}`);
    }
  }

  console.log();
  console.log("============================================");
  console.log("✅ Deployment Complete!");
  console.log("============================================");
  console.log("GMXV2Adapter:", adapterAddress);
  console.log();
  console.log("Next steps:");
  console.log("1. Update .env with GMXV2_ADAPTER_ADDRESS=" + adapterAddress);
  console.log("2. Grant RISK_MANAGER_ROLE to your backend service");
  console.log("3. Start the GMX monitoring listeners");
  console.log("4. Test with a small position");
  console.log("============================================");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
