/**
 * Deploy RWA Marketplace to Arbitrum (Layer 2)
 *
 * 这个脚本展示了如何将RWA合约部署到Layer 2 (Arbitrum)
 * Layer 2 是一条完整的区块链，合约需要部署到Arbitrum网络上才能执行
 */

import hre from "hardhat";

async function main() {
  console.log("=".repeat(60));
  console.log("🚀 开始部署 RWA Marketplace 到 Arbitrum (Layer 2)");
  console.log("=".repeat(60));

  // 获取部署者账户
  const [deployer] = await hre.ethers.getSigners();

  console.log("\n📋 部署信息:");
  console.log("├── 部署账户:", deployer.address);
  console.log("├── 账户余额:", hre.ethers.formatEther(await hre.ethers.provider.getBalance(deployer.address)), "ETH");

  // 获取网络信息
  const network = await hre.ethers.provider.getNetwork();
  console.log("├── 目标网络:", network.name);
  console.log("├── Chain ID:", network.chainId.toString());

  // 验证是否在正确的网络
  if (network.chainId === 42161n) {
    console.log("└── ✅ Arbitrum One (主网)");
  } else if (network.chainId === 421614n) {
    console.log("└── ✅ Arbitrum Sepolia (测试网)");
  } else if (network.chainId === 1337n) {
    console.log("└── ⚠️  本地Hardhat网络 (模拟Arbitrum)");
  } else {
    console.log("└── ❌ 错误: 不是Arbitrum网络!");
    console.log("\n请使用以下命令部署到正确的网络:");
    console.log("  npx hardhat run scripts/deploy-l2-rwa.js --network arbitrumOne");
    console.log("  npx hardhat run scripts/deploy-l2-rwa.js --network arbitrumSepolia");
    return;
  }

  console.log("\n" + "=".repeat(60));
  console.log("📦 第1步: 部署基础设施合约");
  console.log("=".repeat(60));

  // 1. 部署 RWACompliance (KYC/AML合规)
  console.log("\n[1/6] 部署 RWACompliance...");
  const RWACompliance = await hre.ethers.getContractFactory("RWACompliance");
  const compliance = await RWACompliance.deploy(deployer.address);
  await compliance.waitForDeployment();
  const complianceAddress = await compliance.getAddress();
  console.log("  ✅ RWACompliance 部署到:", complianceAddress);
  console.log("     (这个合约现在运行在 Arbitrum 链上!)");

  // 2. 部署 RWAValuation (资产估值)
  console.log("\n[2/6] 部署 RWAValuation...");
  const RWAValuation = await hre.ethers.getContractFactory("RWAValuation");
  const valuation = await RWAValuation.deploy(deployer.address);
  await valuation.waitForDeployment();
  const valuationAddress = await valuation.getAddress();
  console.log("  ✅ RWAValuation 部署到:", valuationAddress);

  // 3. 部署 RWAAssetFactory (资产工厂)
  console.log("\n[3/6] 部署 RWAAssetFactory...");
  const RWAAssetFactory = await hre.ethers.getContractFactory("RWAAssetFactory");
  const assetFactory = await RWAAssetFactory.deploy(
    deployer.address,
    complianceAddress,
    valuationAddress
  );
  await assetFactory.waitForDeployment();
  const factoryAddress = await assetFactory.getAddress();
  console.log("  ✅ RWAAssetFactory 部署到:", factoryAddress);

  // 4. 部署 RWAMarketplace (交易市场)
  console.log("\n[4/6] 部署 RWAMarketplace...");
  const RWAMarketplace = await hre.ethers.getContractFactory("RWAMarketplace");
  const marketplace = await RWAMarketplace.deploy(
    deployer.address,
    factoryAddress,
    complianceAddress,
    valuationAddress,
    deployer.address // Fee collector
  );
  await marketplace.waitForDeployment();
  const marketplaceAddress = await marketplace.getAddress();
  console.log("  ✅ RWAMarketplace 部署到:", marketplaceAddress);

  // 5. 部署 RWAYieldDistributor (收益分配)
  console.log("\n[5/6] 部署 RWAYieldDistributor...");
  const RWAYieldDistributor = await hre.ethers.getContractFactory("RWAYieldDistributor");
  const yieldDistributor = await RWAYieldDistributor.deploy(
    deployer.address,
    factoryAddress
  );
  await yieldDistributor.waitForDeployment();
  const yieldDistributorAddress = await yieldDistributor.getAddress();
  console.log("  ✅ RWAYieldDistributor 部署到:", yieldDistributorAddress);

  // 6. 部署 RWAGovernance (治理)
  console.log("\n[6/6] 部署 RWAGovernance...");
  const RWAGovernance = await hre.ethers.getContractFactory("RWAGovernance");
  const governance = await RWAGovernance.deploy(
    deployer.address,
    factoryAddress
  );
  await governance.waitForDeployment();
  const governanceAddress = await governance.getAddress();
  console.log("  ✅ RWAGovernance 部署到:", governanceAddress);

  console.log("\n" + "=".repeat(60));
  console.log("📋 部署总结");
  console.log("=".repeat(60));

  const deploymentInfo = {
    network: network.name,
    chainId: network.chainId.toString(),
    deployer: deployer.address,
    contracts: {
      RWACompliance: complianceAddress,
      RWAValuation: valuationAddress,
      RWAAssetFactory: factoryAddress,
      RWAMarketplace: marketplaceAddress,
      RWAYieldDistributor: yieldDistributorAddress,
      RWAGovernance: governanceAddress
    }
  };

  console.log("\n所有合约已成功部署到 Arbitrum 链上:");
  console.log(JSON.stringify(deploymentInfo, null, 2));

  // 保存部署地址
  const fs = await import('fs');
  const deploymentPath = `./deployments/arbitrum-${network.chainId}.json`;
  fs.writeFileSync(deploymentPath, JSON.stringify(deploymentInfo, null, 2));
  console.log("\n✅ 部署地址已保存到:", deploymentPath);

  console.log("\n" + "=".repeat(60));
  console.log("🔍 验证合约 (可选)");
  console.log("=".repeat(60));

  if (network.chainId !== 1337n) {
    console.log("\n在区块浏览器上验证合约:");
    console.log("npx hardhat verify --network", network.name, complianceAddress, deployer.address);
    console.log("npx hardhat verify --network", network.name, valuationAddress, deployer.address);
    console.log("npx hardhat verify --network", network.name, factoryAddress, deployer.address, complianceAddress, valuationAddress);

    console.log("\n查看合约:");
    if (network.chainId === 42161n) {
      console.log("Arbitrum One 浏览器: https://arbiscan.io/address/" + factoryAddress);
    } else if (network.chainId === 421614n) {
      console.log("Arbitrum Sepolia 浏览器: https://sepolia.arbiscan.io/address/" + factoryAddress);
    }
  }

  console.log("\n" + "=".repeat(60));
  console.log("✨ 部署完成!");
  console.log("=".repeat(60));

  console.log("\n📚 重要概念:");
  console.log("├── 这些合约现在运行在 Arbitrum 区块链上");
  console.log("├── Arbitrum 是一条完整的 Layer 2 区块链");
  console.log("├── 用户通过连接到 Arbitrum 网络来使用这些合约");
  console.log("├── 所有交易在 Arbitrum 链上确认 (0.25秒)");
  console.log("└── Gas费使用 Arbitrum 上的 ETH 支付 (比主网便宜99%+)");

  console.log("\n🔗 下一步:");
  console.log("1. 设置权限和角色");
  console.log("2. 创建第一个RWA资产进行测试");
  console.log("3. 配置与L1的跨链通信 (如果需要)");
  console.log("4. 部署前端并连接到这些合约地址\n");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("\n❌ 部署失败:", error);
    process.exit(1);
  });
