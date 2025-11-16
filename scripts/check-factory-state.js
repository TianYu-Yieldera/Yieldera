const hre = require("hardhat");

async function main() {
  console.log("\n🔍 Checking Factory State\n");

  const factoryAddr = "0x9e667a4ce092086C63c667e1Ea575B9Aa2a4762B";
  const factory = await hre.ethers.getContractAt("TreasuryAssetFactory", factoryAddr);

  console.log("Checking total assets...");
  const total = await factory.getTotalAssets();
  console.log(`Total assets: ${total}`);

  if (total > 0) {
    console.log("\n📋 Asset List:");
    for (let i = 1; i <= total; i++) {
      try {
        const assetInfo = await factory.getAssetInfo(i);
        console.log(`Asset ${i}:`);
        console.log(`  CUSIP: ${assetInfo.cusip}`);
        console.log(`  Type: ${assetInfo.treasuryType}`);
        console.log(`  Term: ${assetInfo.maturityTerm}`);
        console.log(`  Token: ${assetInfo.tokenAddress}`);
        console.log(`  Status: ${assetInfo.status}`);
      } catch (error) {
        console.log(`Asset ${i}: Failed to get info - ${error.message}`);
      }
    }
  } else {
    console.log("\n❌ No assets created yet!");
    console.log("The createTreasuryAsset transactions must have reverted.");
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
