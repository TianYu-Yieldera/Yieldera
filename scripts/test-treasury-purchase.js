/**
 * Treasury Token Purchase Test
 * 测试国债代币购买流程
 */

const { ethers } = require("hardhat");

async function main() {
  console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
  console.log("🏦 Treasury Token Purchase Test");
  console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n");

  const [signer] = await ethers.getSigners();
  console.log("📝 Testing with account:", signer.address);

  // Load contract addresses from environment
  const factoryAddress = process.env.L2_TREASURY_FACTORY;
  const usdcAddress = process.env.L2_USDC;

  if (!factoryAddress || factoryAddress === "") {
    console.error("❌ L2_TREASURY_FACTORY not configured in .env");
    console.log("\nTo deploy treasury contracts, run:");
    console.log("  npx hardhat run scripts/deploy-treasury.js --network arbitrumSepolia");
    process.exit(1);
  }

  if (!usdcAddress || usdcAddress === "") {
    console.error("❌ L2_USDC not configured in .env");
    console.log("\nPlease add USDC contract address to .env");
    process.exit(1);
  }

  try {
    // Connect to contracts
    console.log("\n📡 Connecting to contracts...");
    const factory = await ethers.getContractAt("TreasuryAssetFactory", factoryAddress);
    const usdc = await ethers.getContractAt(
      ["function balanceOf(address) view returns (uint256)",
       "function approve(address,uint256) returns (bool)",
       "function decimals() view returns (uint8)"],
      usdcAddress
    );

    // Check USDC balance
    console.log("\n💰 Checking USDC balance...");
    const balance = await usdc.balanceOf(signer.address);
    const decimals = await usdc.decimals();
    const balanceFormatted = ethers.formatUnits(balance, decimals);
    console.log("   Balance:", balanceFormatted, "USDC");

    if (balance === 0n) {
      console.log("\n⚠️  Warning: You have 0 USDC");
      console.log("   Get testnet USDC from: https://faucet.circle.com/");
      console.log("   Or bridge from Sepolia");
      process.exit(0);
    }

    // Get available treasury assets
    console.log("\n📊 Fetching available treasury assets...");
    const assets = await factory.getAllAssets();

    if (assets.length === 0) {
      console.log("❌ No treasury assets found!");
      console.log("\nCreate assets using the admin panel or run:");
      console.log("  npx hardhat run scripts/create-treasury-assets.js --network arbitrumSepolia");
      process.exit(1);
    }

    console.log(`   Found ${assets.length} asset(s):`);
    for (let i = 0; i < Math.min(assets.length, 3); i++) {
      console.log(`   ${i + 1}. ${assets[i]}`);
    }

    // Select first asset for testing
    const assetAddress = assets[0];
    console.log(`\n🎯 Testing with asset: ${assetAddress}`);

    // Get asset details
    const asset = await ethers.getContractAt("TreasuryToken", assetAddress);
    const name = await asset.name();
    const symbol = await asset.symbol();
    const minInvestment = await asset.minInvestment();

    console.log("   Name:", name);
    console.log("   Symbol:", symbol);
    console.log("   Min Investment:", ethers.formatUnits(minInvestment, decimals), "USDC");

    // Determine purchase amount (min investment or $10, whichever is higher)
    const purchaseAmount = balance < ethers.parseUnits("10", decimals)
      ? minInvestment
      : ethers.parseUnits("10", decimals);

    const purchaseFormatted = ethers.formatUnits(purchaseAmount, decimals);
    console.log(`\n💵 Purchasing ${purchaseFormatted} USDC worth of tokens...`);

    // Step 1: Approve USDC
    console.log("\n1️⃣  Approving USDC...");
    const approveTx = await usdc.approve(assetAddress, purchaseAmount);
    console.log("   Tx submitted:", approveTx.hash);
    await approveTx.wait();
    console.log("   ✅ Approved");

    // Step 2: Purchase tokens
    console.log("\n2️⃣  Purchasing treasury tokens...");
    const purchaseTx = await asset.purchase(purchaseAmount);
    console.log("   Tx submitted:", purchaseTx.hash);
    const receipt = await purchaseTx.wait();
    console.log("   ✅ Purchase successful!");
    console.log("   Gas used:", receipt.gasUsed.toString());

    // Step 3: Verify balance
    console.log("\n3️⃣  Verifying balances...");
    const tokenBalance = await asset.balanceOf(signer.address);
    const newUsdcBalance = await usdc.balanceOf(signer.address);

    console.log("   Treasury Token Balance:", ethers.formatEther(tokenBalance));
    console.log("   Remaining USDC:", ethers.formatUnits(newUsdcBalance, decimals));

    // Step 4: Check yield info
    console.log("\n4️⃣  Checking yield information...");
    try {
      const apy = await asset.currentAPY();
      const nextDistribution = await asset.nextDistributionTime();

      console.log("   Current APY:", (Number(apy) / 100).toFixed(2) + "%");

      if (nextDistribution > 0) {
        const date = new Date(Number(nextDistribution) * 1000);
        console.log("   Next Distribution:", date.toLocaleString());
      }
    } catch (e) {
      console.log("   ⚠️  Yield info not available yet");
    }

    console.log("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
    console.log("✅ Test completed successfully!");
    console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n");

    console.log("Next steps:");
    console.log("  - Check your holdings in the frontend: http://localhost:5173/treasury/holdings");
    console.log("  - View market: http://localhost:5173/treasury");
    console.log("  - Monitor yields: Wait for daily distribution");

  } catch (error) {
    console.error("\n❌ Test failed:", error.message);

    if (error.message.includes("insufficient funds")) {
      console.log("\n💡 Tip: Get more testnet USDC from https://faucet.circle.com/");
    } else if (error.message.includes("ERC20: insufficient allowance")) {
      console.log("\n💡 Tip: Make sure USDC approval succeeded");
    }

    process.exit(1);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
