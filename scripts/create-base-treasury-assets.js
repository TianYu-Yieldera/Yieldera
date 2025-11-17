const hre = require("hardhat");

async function main() {
  console.log("\n📊 Creating Base Treasury Sample Assets\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("Account:", deployer.address);

  const balance = await deployer.provider.getBalance(deployer.address);
  console.log("Balance:", hre.ethers.formatEther(balance), "ETH");

  // Use deployed factory address from Base Sepolia
  const FACTORY_ADDRESS = process.env.BASE_TREASURY_FACTORY || "0x9e667a4ce092086C63c667e1Ea575B9Aa2a4762B";
  const ORACLE_ADDRESS = process.env.BASE_TREASURY_ORACLE || "0xB478ca7F5f03f2700BfC56613bb22546D6D10681";

  console.log("\nUsing contracts:");
  console.log("Factory:", FACTORY_ADDRESS);
  console.log("Oracle:", ORACLE_ADDRESS);

  const TreasuryAssetFactory = await hre.ethers.getContractAt("TreasuryAssetFactory", FACTORY_ADDRESS);
  const TreasuryPriceOracle = await hre.ethers.getContractAt("TreasuryPriceOracle", ORACLE_ADDRESS);

  const currentTime = Math.floor(Date.now() / 1000);
  const oneDay = 24 * 60 * 60;

  // Conservative portfolio: More short-term and mid-term treasuries
  // Note: Assets 1-3 already exist (4W, 13W, 26W T-Bills), creating T-Notes only
  const assets = [
    { type: 1, term: "2Y", cusip: "91282CHX6", maturity: 730, coupon: 450, name: "T-Note 2Y" },
    { type: 1, term: "5Y", cusip: "912828YL2", maturity: 1825, coupon: 435, name: "T-Note 5Y" },
    { type: 1, term: "10Y", cusip: "912828YK4", maturity: 3650, coupon: 425, name: "T-Note 10Y" }
  ];

  console.log("\n📝 Creating 3 remaining T-Notes (T-Bills 1-3 already exist)...\n");

  const deployments = { assets: [] };

  for (let i = 0; i < assets.length; i++) {
    const asset = assets[i];
    console.log(`[${i+1}/3] Creating ${asset.name}...`);

    try {
      // Estimate gas first
      const gasEstimate = await TreasuryAssetFactory.createTreasuryAsset.estimateGas(
        asset.type,
        asset.term,
        asset.cusip,
        currentTime,
        currentTime + (asset.maturity * oneDay),
        hre.ethers.parseEther("10000"),
        asset.coupon
      );

      const tx = await TreasuryAssetFactory.createTreasuryAsset(
        asset.type,
        asset.term,
        asset.cusip,
        currentTime,
        currentTime + (asset.maturity * oneDay),
        hre.ethers.parseEther("10000"), // Higher liquidity for conservative assets
        asset.coupon,
        { gasLimit: gasEstimate * 120n / 100n } // 20% buffer
      );

      const receipt = await tx.wait();
      console.log(`  ✅ Created (tx: ${receipt.hash.substring(0, 10)}...)`);

      deployments.assets.push({
        name: asset.name,
        cusip: asset.cusip,
        term: asset.term,
        type: asset.type === 0 ? "Treasury Bill" : "Treasury Note"
      });
    } catch (error) {
      console.error(`  ❌ Failed: ${error.message}`);
    }
  }

  console.log("\n📊 Setting initial prices...\n");

  const prices = [
    { assetId: 1, price: "995", yield: 535 },  // 4W T-Bill (near par)
    { assetId: 2, price: "980", yield: 540 },  // 13W T-Bill
    { assetId: 3, price: "975", yield: 535 },  // 26W T-Bill
    { assetId: 4, price: "985", yield: 465 },  // 2Y T-Note
    { assetId: 5, price: "970", yield: 455 },  // 5Y T-Note
    { assetId: 6, price: "950", yield: 475 }   // 10Y T-Note
  ];

  for (const p of prices) {
    try {
      console.log(`Setting price for Asset ${p.assetId}...`);

      // Estimate gas first
      const gasEstimate = await TreasuryPriceOracle.updatePrice.estimateGas(
        p.assetId,
        hre.ethers.parseEther(p.price),
        p.yield,
        "BASE_INITIAL_SETUP"
      );

      const tx = await TreasuryPriceOracle.updatePrice(
        p.assetId,
        hre.ethers.parseEther(p.price),
        p.yield,
        "BASE_INITIAL_SETUP",
        { gasLimit: gasEstimate * 120n / 100n } // 20% buffer
      );
      await tx.wait();
      console.log(`  ✅ Price set`);
    } catch (error) {
      console.error(`  ❌ Failed: ${error.message}`);
    }
  }

  console.log("\n" + "=".repeat(60));
  console.log("✅ Base Treasury Assets Created!");
  console.log("=".repeat(60));
  console.log("\nAssets created:");
  deployments.assets.forEach((a, i) => {
    console.log(`${i+1}. ${a.name} (CUSIP: ${a.cusip}, Type: ${a.type})`);
  });

  console.log("\n🌐 View on BaseScan:");
  console.log(`   Factory: https://sepolia.basescan.org/address/${FACTORY_ADDRESS}`);
  console.log(`   Oracle:  https://sepolia.basescan.org/address/${ORACLE_ADDRESS}`);
  console.log();
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("\n❌ Failed:");
    console.error(error);
    process.exit(1);
  });
