const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🚀 开始部署 YielderaAirdrop 合约...\n");

  // 获取部署者账户
  const [deployer] = await hre.ethers.getSigners();
  console.log("📝 部署者地址:", deployer.address);

  const balance = await hre.ethers.provider.getBalance(deployer.address);
  console.log("💰 部署者余额:", hre.ethers.formatEther(balance), "ETH\n");

  // 从环境变量获取代币地址，如果没有则部署测试代币
  let tokenAddress = process.env.TOKEN_ADDRESS;

  if (!tokenAddress) {
    console.log("⚠️  未找到 TOKEN_ADDRESS，部署测试 ERC20 代币...");

    // 部署测试 ERC20 代币
    const TestToken = await hre.ethers.getContractFactory("MockERC20");
    const testToken = await TestToken.deploy(
      "Yieldera Token",
      "YLD",
      hre.ethers.parseEther("1000000000") // 10亿代币
    );
    await testToken.waitForDeployment();
    tokenAddress = await testToken.getAddress();

    console.log("✅ 测试代币部署成功:", tokenAddress);
    console.log("   代币名称: Yieldera Token (YLD)");
    console.log("   总供应量: 1,000,000,000 YLD\n");
  } else {
    console.log("✅ 使用现有代币地址:", tokenAddress, "\n");
  }

  // 部署 YielderaAirdrop 合约
  console.log("📦 部署 YielderaAirdrop 合约...");
  const YielderaAirdrop = await hre.ethers.getContractFactory("YielderaAirdrop");
  const airdrop = await YielderaAirdrop.deploy(tokenAddress);
  await airdrop.waitForDeployment();

  const airdropAddress = await airdrop.getAddress();
  console.log("✅ YielderaAirdrop 部署成功:", airdropAddress, "\n");

  // 获取网络信息
  const network = await hre.ethers.provider.getNetwork();
  const chainId = network.chainId;
  const blockNumber = await hre.ethers.provider.getBlockNumber();

  // 保存部署信息
  const deploymentInfo = {
    network: hre.network.name,
    chainId: Number(chainId),
    deployer: deployer.address,
    timestamp: new Date().toISOString(),
    blockNumber: blockNumber,
    contracts: {
      YielderaAirdrop: {
        address: airdropAddress,
        tokenAddress: tokenAddress,
      },
    },
  };

  // 创建 deployments 目录（如果不存在）
  const deploymentsDir = path.join(__dirname, "../deployments");
  if (!fs.existsSync(deploymentsDir)) {
    fs.mkdirSync(deploymentsDir, { recursive: true });
  }

  // 保存部署信息到文件
  const deploymentFile = path.join(
    deploymentsDir,
    `${hre.network.name}.json`
  );
  fs.writeFileSync(deploymentFile, JSON.stringify(deploymentInfo, null, 2));

  console.log("📄 部署信息已保存到:", deploymentFile, "\n");

  // 打印部署总结
  console.log("=" .repeat(60));
  console.log("🎉 部署完成！");
  console.log("=" .repeat(60));
  console.log("网络:", hre.network.name);
  console.log("链ID:", chainId);
  console.log("区块高度:", blockNumber);
  console.log("代币地址:", tokenAddress);
  console.log("空投合约:", airdropAddress);
  console.log("=" .repeat(60));

  // 如果是测试网，打印验证命令
  if (hre.network.name === "sepolia" || hre.network.name === "mainnet") {
    console.log("\n📝 验证合约命令:");
    console.log(
      `npx hardhat verify --network ${hre.network.name} ${airdropAddress} ${tokenAddress}`
    );
    console.log("\n🔄 更新 Subgraph 配置:");
    console.log(`   1. 编辑 subgraph/subgraph.yaml`);
    console.log(`   2. 设置 address: "${airdropAddress}"`);
    console.log(`   3. 设置 startBlock: ${blockNumber}`);
  }

  console.log("\n");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("❌ 部署失败:", error);
    process.exit(1);
  });
