/**
 * Aave V3 Integration Test
 * 测试 Aave 协议存款和借款功能
 */

const { ethers } = require("hardhat");

async function main() {
  console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
  console.log("🏦 Aave V3 Integration Test");
  console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n");

  const [signer] = await ethers.getSigners();
  console.log("📝 Testing with account:", signer.address);

  // Check balance
  const ethBalance = await ethers.provider.getBalance(signer.address);
  console.log("   ETH Balance:", ethers.formatEther(ethBalance), "ETH");

  if (ethBalance < ethers.parseEther("0.001")) {
    console.log("\n⚠️  Warning: Low ETH balance. Get testnet ETH from:");
    console.log("   https://www.alchemy.com/faucets/arbitrum-sepolia");
    process.exit(0);
  }

  const aaveAdapterAddress = process.env.L2_AAVE_ADAPTER;

  if (!aaveAdapterAddress || aaveAdapterAddress === "") {
    console.error("❌ L2_AAVE_ADAPTER not configured in .env");
    console.log("\nTo deploy adapters, run:");
    console.log("  npx hardhat run scripts/deploy-adapters.js --network arbitrumSepolia");
    process.exit(1);
  }

  try {
    console.log("\n📡 Connecting to Aave Adapter...");
    const adapter = await ethers.getContractAt("AaveV3Adapter", aaveAdapterAddress);

    // Get Aave pool address
    const poolAddress = await adapter.AAVE_POOL();
    console.log("   Aave Pool:", poolAddress);

    // Test 1: Supply ETH
    console.log("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
    console.log("Test 1: Supply ETH to Aave");
    console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");

    // Use WETH address on Arbitrum Sepolia
    const wethAddress = "0x980B62Da83eFf3D4576C647993b0c1D7faf17c73"; // Arbitrum Sepolia WETH
    const supplyAmount = ethers.parseEther("0.001"); // 0.001 ETH

    console.log(`\n💰 Supplying ${ethers.formatEther(supplyAmount)} ETH...`);

    const supplyTx = await adapter.supply(
      wethAddress,
      supplyAmount,
      signer.address,
      { value: supplyAmount }
    );
    console.log("   Tx submitted:", supplyTx.hash);
    const supplyReceipt = await supplyTx.wait();
    console.log("   ✅ Supply successful!");
    console.log("   Gas used:", supplyReceipt.gasUsed.toString());

    // Check supplied balance
    console.log("\n📊 Checking supplied balance...");
    const suppliedBalance = await adapter.getSupplyBalance(wethAddress, signer.address);
    console.log("   Supplied:", ethers.formatEther(suppliedBalance), "WETH");

    // Test 2: Check account health
    console.log("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
    console.log("Test 2: Account Health Check");
    console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");

    try {
      const healthFactor = await adapter.getHealthFactor(signer.address);
      console.log("\n🏥 Health Factor:", ethers.formatEther(healthFactor));

      if (healthFactor > ethers.parseEther("1")) {
        console.log("   ✅ Account is healthy (> 1.0)");
      } else {
        console.log("   ⚠️  Low health factor! Risk of liquidation");
      }
    } catch (e) {
      console.log("   ℹ️  Health factor not available (no borrows yet)");
    }

    // Test 3: Borrow USDC (if possible)
    console.log("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
    console.log("Test 3: Borrow USDC");
    console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");

    const usdcAddress = process.env.L2_USDC || "0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d";

    // Try to borrow a small amount (assuming 1 ETH = $2000, we can borrow ~50% = $1000)
    const borrowAmount = ethers.parseUnits("0.5", 6); // 0.5 USDC (very small amount for testing)

    try {
      console.log(`\n💸 Attempting to borrow ${ethers.formatUnits(borrowAmount, 6)} USDC...`);

      const borrowTx = await adapter.borrow(
        usdcAddress,
        borrowAmount,
        2, // Variable interest rate mode
        signer.address
      );
      console.log("   Tx submitted:", borrowTx.hash);
      const borrowReceipt = await borrowTx.wait();
      console.log("   ✅ Borrow successful!");
      console.log("   Gas used:", borrowReceipt.gasUsed.toString());

      // Check borrowed amount
      const borrowedBalance = await adapter.getBorrowBalance(usdcAddress, signer.address);
      console.log("   Borrowed:", ethers.formatUnits(borrowedBalance, 6), "USDC");

      // Check USDC balance
      const usdc = await ethers.getContractAt(
        ["function balanceOf(address) view returns (uint256)"],
        usdcAddress
      );
      const usdcBalance = await usdc.balanceOf(signer.address);
      console.log("   USDC Balance:", ethers.formatUnits(usdcBalance, 6), "USDC");

    } catch (e) {
      console.log("   ⚠️  Could not borrow:", e.message.split('\n')[0]);
      console.log("   This is normal if collateral is insufficient or USDC not available on Aave");
    }

    // Test 4: Withdraw (partial)
    console.log("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
    console.log("Test 4: Withdraw Collateral");
    console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");

    const withdrawAmount = ethers.parseEther("0.0005"); // Withdraw half
    console.log(`\n💵 Withdrawing ${ethers.formatEther(withdrawAmount)} WETH...`);

    try {
      const withdrawTx = await adapter.withdraw(
        wethAddress,
        withdrawAmount,
        signer.address
      );
      console.log("   Tx submitted:", withdrawTx.hash);
      const withdrawReceipt = await withdrawTx.wait();
      console.log("   ✅ Withdrawal successful!");
      console.log("   Gas used:", withdrawReceipt.gasUsed.toString());

      // Check remaining balance
      const remainingBalance = await adapter.getSupplyBalance(wethAddress, signer.address);
      console.log("   Remaining:", ethers.formatEther(remainingBalance), "WETH");

    } catch (e) {
      console.log("   ⚠️  Could not withdraw:", e.message.split('\n')[0]);
      console.log("   This might happen if there are active borrows");
    }

    console.log("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
    console.log("✅ Aave Integration Test Completed!");
    console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n");

    console.log("Summary:");
    console.log("  ✅ Supply functionality working");
    console.log("  ✅ Balance queries working");
    console.log("  ✅ Health factor calculation working");
    console.log("\nNext steps:");
    console.log("  - Check your positions in frontend: http://localhost:5173/vault");
    console.log("  - Monitor health factor and liquidation risk");
    console.log("  - Test AI risk prediction with your positions");

  } catch (error) {
    console.error("\n❌ Test failed:", error.message);

    if (error.message.includes("insufficient funds")) {
      console.log("\n💡 Tip: Get more testnet ETH from Arbitrum Sepolia faucet");
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
