/**
 * Automated State Synchronization Script
 * Monitors L2 state and submits to L1 periodically
 *
 * Usage:
 *   node scripts/automation/state-sync.js
 *
 * Environment Variables:
 *   L2_STATE_AGGREGATOR - L2StateAggregator contract address
 *   L1_STATE_REGISTRY - L1StateRegistry contract address
 *   SYNC_INTERVAL - Sync interval in minutes (default: 60)
 *   AUTO_SUBMIT - Enable auto submission (default: true)
 */

import hre from "hardhat";
import * as dotenv from "dotenv";

dotenv.config();

// Configuration
const SYNC_INTERVAL = parseInt(process.env.SYNC_INTERVAL || "60") * 60 * 1000; // Convert to ms
const AUTO_SUBMIT = process.env.AUTO_SUBMIT !== "false";
const L2_RPC = process.env.ARBITRUM_SEPOLIA_RPC_URL || "https://sepolia-rollup.arbitrum.io/rpc";
const L1_RPC = process.env.SEPOLIA_RPC_URL;

class StateSyncMonitor {
  constructor() {
    this.running = false;
    this.l2Provider = null;
    this.l1Provider = null;
    this.l2Aggregator = null;
    this.l1Registry = null;
  }

  async initialize() {
    console.log("====================================");
    console.log("🔄 State Sync Monitor Starting...");
    console.log("====================================\n");

    // Connect to L2
    this.l2Provider = new hre.ethers.JsonRpcProvider(L2_RPC);
    const l2Signer = new hre.ethers.Wallet(process.env.PRIVATE_KEY, this.l2Provider);

    // Connect to L1
    if (L1_RPC) {
      this.l1Provider = new hre.ethers.JsonRpcProvider(L1_RPC);
    }

    // Get contract addresses
    const l2AggregatorAddress = process.env.L2_STATE_AGGREGATOR;
    const l1RegistryAddress = process.env.L1_STATE_REGISTRY;

    if (!l2AggregatorAddress) {
      throw new Error("L2_STATE_AGGREGATOR address not set in .env");
    }

    console.log("📍 Configuration:");
    console.log("├─ L2 Aggregator:", l2AggregatorAddress);
    console.log("├─ L1 Registry:", l1RegistryAddress || "Not monitoring");
    console.log("├─ Sync Interval:", SYNC_INTERVAL / 60000, "minutes");
    console.log("└─ Auto Submit:", AUTO_SUBMIT ? "Enabled" : "Disabled");
    console.log("");

    // Load contracts
    const L2StateAggregator = await hre.ethers.getContractFactory("L2StateAggregator");
    this.l2Aggregator = L2StateAggregator.attach(l2AggregatorAddress).connect(l2Signer);

    if (l1RegistryAddress && this.l1Provider) {
      const L1StateRegistry = await hre.ethers.getContractFactory("L1StateRegistry");
      this.l1Registry = L1StateRegistry.attach(l1RegistryAddress).connect(this.l1Provider);
    }

    console.log("✅ Contracts loaded successfully\n");
  }

  async checkL2State() {
    try {
      console.log("📊 Checking L2 State...");

      // Get current system state
      const state = await this.l2Aggregator.getSystemState();
      const currentRoot = await this.l2Aggregator.currentStateRoot();
      const canSubmit = await this.l2Aggregator.canSubmitToL1();
      const timeUntilNext = await this.l2Aggregator.timeUntilNextSubmission();

      console.log("├─ Block Number:", state.blockNumber.toString());
      console.log("├─ Total Collateral:", hre.ethers.formatEther(state.totalCollateral), "tokens");
      console.log("├─ Total Debt:", hre.ethers.formatUnits(state.totalDebt, 6), "LUSD");
      console.log("├─ Active Positions:", state.activePositions.toString());
      console.log("├─ Total Orders:", state.totalOrders.toString());
      console.log("├─ Current State Root:", currentRoot);
      console.log("├─ Can Submit:", canSubmit);
      console.log("└─ Time Until Next:", timeUntilNext.toString(), "seconds");
      console.log("");

      return { state, currentRoot, canSubmit, timeUntilNext };
    } catch (error) {
      console.error("❌ Error checking L2 state:", error.message);
      return null;
    }
  }

  async submitToL1() {
    try {
      console.log("🚀 Submitting state to L1...");

      const canSubmit = await this.l2Aggregator.canSubmitToL1();

      if (!canSubmit) {
        const timeUntilNext = await this.l2Aggregator.timeUntilNextSubmission();
        console.log("⏳ Cannot submit yet. Time remaining:", timeUntilNext.toString(), "seconds");
        return false;
      }

      // Submit to L1
      const tx = await this.l2Aggregator.submitToL1();
      console.log("├─ Transaction sent:", tx.hash);

      const receipt = await tx.wait();
      console.log("├─ Transaction confirmed");
      console.log("├─ Gas used:", receipt.gasUsed.toString());
      console.log("└─ Block number:", receipt.blockNumber);
      console.log("");

      return true;
    } catch (error) {
      console.error("❌ Error submitting to L1:", error.message);
      return false;
    }
  }

  async checkL1State() {
    if (!this.l1Registry) {
      return null;
    }

    try {
      console.log("🔍 Checking L1 State Registry...");

      const [latestRoot, latestBlock, timestamp] = await this.l1Registry.getLatestState();
      const isStateFresh = await this.l1Registry.isStateFresh();
      const timeSinceLast = await this.l1Registry.timeSinceLastSubmission();

      console.log("├─ Latest L2 Block:", latestBlock.toString());
      console.log("├─ Latest State Root:", latestRoot);
      console.log("├─ Timestamp:", new Date(Number(timestamp) * 1000).toISOString());
      console.log("├─ Is Fresh:", isStateFresh);
      console.log("└─ Time Since Last:", timeSinceLast.toString(), "seconds");
      console.log("");

      return { latestRoot, latestBlock, timestamp, isStateFresh };
    } catch (error) {
      console.error("❌ Error checking L1 state:", error.message);
      return null;
    }
  }

  async monitorLoop() {
    this.running = true;
    let iteration = 1;

    while (this.running) {
      console.log(`\n${"=".repeat(60)}`);
      console.log(`📈 Monitoring Iteration #${iteration}`);
      console.log(`⏰ Time: ${new Date().toISOString()}`);
      console.log("=".repeat(60) + "\n");

      // Check L2 state
      const l2State = await this.checkL2State();

      // Check L1 state (if available)
      const l1State = await this.checkL1State();

      // Auto-submit if enabled and ready
      if (AUTO_SUBMIT && l2State && l2State.canSubmit) {
        console.log("✅ Auto-submit enabled and ready");
        await this.submitToL1();
      }

      // Verify L1-L2 consistency
      if (l1State && l2State) {
        this.verifyConsistency(l1State, l2State);
      }

      iteration++;

      // Wait for next iteration
      console.log(`⏸️  Waiting ${SYNC_INTERVAL / 60000} minutes until next check...\n`);
      await this.sleep(SYNC_INTERVAL);
    }
  }

  verifyConsistency(l1State, l2State) {
    console.log("🔎 Verifying L1-L2 Consistency...");

    if (l1State.latestRoot === l2State.currentRoot) {
      console.log("✅ State roots match - System synchronized");
    } else {
      console.log("⚠️  State roots differ - Pending synchronization");
      console.log("├─ L1 Root:", l1State.latestRoot);
      console.log("└─ L2 Root:", l2State.currentRoot);
    }

    if (!l1State.isStateFresh) {
      console.log("⚠️  L1 state is stale (> 2 hours old)");
    }

    console.log("");
  }

  sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  async stop() {
    console.log("\n🛑 Stopping monitor...");
    this.running = false;
  }
}

// Main execution
async function main() {
  const monitor = new StateSyncMonitor();

  // Handle graceful shutdown
  process.on("SIGINT", async () => {
    await monitor.stop();
    process.exit(0);
  });

  process.on("SIGTERM", async () => {
    await monitor.stop();
    process.exit(0);
  });

  try {
    await monitor.initialize();
    await monitor.monitorLoop();
  } catch (error) {
    console.error("💥 Fatal error:", error);
    process.exit(1);
  }
}

// Run if called directly
if (import.meta.url === `file://${process.argv[1]}`) {
  main().catch(console.error);
}

export default StateSyncMonitor;
