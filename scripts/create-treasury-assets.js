/**
 * Create Treasury Assets Script
 * 创建国债资产用于测试
 */

const { ethers } = require("hardhat");

async function main() {
  console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
  console.log("📝 Creating Treasury Assets");
  console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n");

  const [signer] = await ethers.getSigners();
  console.log("Creating with account:", signer.address);

  // Load contract addresses
  const factoryAddress = process.env.L2_TREASURY_FACTORY;
  const usdcAddress = process.env.L2_USDC;

  if (!factoryAddress || factoryAddress === "") {
    console.error("❌ L2_TREASURY_FACTORY not configured");
    process.exit(1);
  }

  console.log("\n📡 Connecting to Treasury Factory...");
  console.log("   Factory:", factoryAddress);

  const factory = await ethers.getContractAt("TreasuryAssetFactory", factoryAddress);

  // Define treasury assets to create
  const assets = [
    {
      name: "3-Month T-Bill",
      symbol: "TBILL-3M",
      assetType: 0, // T-Bill
      maturityMonths: 3,
      initialApy: 525, // 5.25%
      minInvestment: ethers.parseUnits("1", 6), // $1 USDC
      cusip: "912796YV9",
      description: "US Treasury Bill - 3 Month"
    },
    {
      name: "6-Month T-Bill",
      symbol: "TBILL-6M",
      assetType: 0,
      maturityMonths: 6,
      initialApy: 540, // 5.40%
      minInvestment: ethers.parseUnits("1", 6),
      cusip: "912796YW7",
      description: "US Treasury Bill - 6 Month"
    },
    {
      name: "2-Year T-Note",
      symbol: "TNOTE-2Y",
      assetType: 1, // T-Note
      maturityMonths: 24,
      initialApy: 475, // 4.75%
      minInvestment: ethers.parseUnits("1", 6),
      cusip: "912828YX8",
      description: "US Treasury Note - 2 Year"
    },
    {
      name: "5-Year T-Note",
      symbol: "TNOTE-5Y",
      assetType: 1,
      maturityMonths: 60,
      initialApy: 450, // 4.50%
      minInvestment: ethers.parseUnits("1", 6),
      cusip: "912828YZ3",
      description: "US Treasury Note - 5 Year"
    },
    {
      name: "10-Year T-Note",
      symbol: "TNOTE-10Y",
      assetType: 1,
      maturityMonths: 120,
      initialApy: 435, // 4.35%
      minInvestment: ethers.parseUnits("1", 6),
      cusip: "912828ZA7",
      description: "US Treasury Note - 10 Year"
    }
  ];

  console.log(`\n🏦 Creating ${assets.length} treasury assets...\n`);

  const createdAssets = [];

  for (let i = 0; i < assets.length; i++) {
    const asset = assets[i];
    console.log(`${i + 1}/${assets.length} Creating: ${asset.name}`);
    console.log(`   Symbol: ${asset.symbol}`);
    console.log(`   APY: ${(asset.initialApy / 100).toFixed(2)}%`);
    console.log(`   Maturity: ${asset.maturityMonths} months`);
    console.log(`   Min Investment: $${ethers.formatUnits(asset.minInvestment, 6)}`);

    try {
      // Create the asset
      const tx = await factory.createTreasuryAsset(
        asset.name,
        asset.symbol,
        asset.assetType,
        asset.maturityMonths,
        asset.initialApy,
        asset.minInvestment,
        usdcAddress || ethers.ZeroAddress, // Payment token
        asset.cusip
      );

      console.log(`   Tx submitted: ${tx.hash}`);
      const receipt = await tx.wait();
      console.log(`   ✅ Created! Gas used: ${receipt.gasUsed.toString()}`);

      // Get the created asset address from events
      const event = receipt.logs.find(
        log => {
          try {
            const parsed = factory.interface.parseLog(log);
            return parsed && parsed.name === 'AssetCreated';
          } catch {
            return false;
          }
        }
      );

      if (event) {
        const parsed = factory.interface.parseLog(event);
        const assetAddress = parsed.args[0];
        console.log(`   Asset Address: ${assetAddress}`);
        createdAssets.push({
          name: asset.name,
          symbol: asset.symbol,
          address: assetAddress
        });
      }

      console.log("");
    } catch (error) {
      console.error(`   ❌ Failed: ${error.message.split('\n')[0]}`);
      console.log("");
    }
  }

  // Summary
  console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
  console.log("✅ Asset Creation Complete");
  console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n");

  console.log(`Created ${createdAssets.length} assets:`);
  createdAssets.forEach((asset, i) => {
    console.log(`${i + 1}. ${asset.name} (${asset.symbol})`);
    console.log(`   ${asset.address}`);
  });

  console.log("\nNext steps:");
  console.log("  - View assets in frontend: http://localhost:5173/treasury");
  console.log("  - Test purchase with: npx hardhat run scripts/test-treasury-purchase.js --network arbitrumSepolia");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
