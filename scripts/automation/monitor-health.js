/**
 * System Health Monitoring Script
 * Monitors L1 and L2 system health with alerting
 *
 * Usage:
 *   node scripts/automation/monitor-health.js
 *
 * Features:
 *   - Monitors system collateral ratio
 *   - Tracks liquidatable positions
 *   - Checks state freshness
 *   - Alerts on critical conditions
 */

import hre from "hardhat";
import * as dotenv from "dotenv";

dotenv.config();

// Alert thresholds
const THRESHOLDS = {
  MIN_COLLATERAL_RATIO: 150, // 150%
  CRITICAL_COLLATERAL_RATIO: 120, // 120%
  MAX_STATE_AGE_HOURS: 2,
  MAX_PENDING_LIQUIDATIONS: 5,
};

// Configuration
const CHECK_INTERVAL = parseInt(process.env.CHECK_INTERVAL || "5") * 60 * 1000; // 5 minutes default
const ALERT_WEBHOOK = process.env.ALERT_WEBHOOK_URL; // Discord/Slack webhook

class HealthMonitor {
  constructor() {
    this.alerts = [];
    this.lastAlertTime = new Map();
    this.ALERT_COOLDOWN = 15 * 60 * 1000; // 15 minutes between same alerts
  }

  async initialize() {
    console.log("====================================");
    console.log("🏥 System Health Monitor Starting...");
    console.log("====================================\n");

    this.l2Provider = new hre.ethers.JsonRpcProvider(
      process.env.ARBITRUM_SEPOLIA_RPC_URL || "https://sepolia-rollup.arbitrum.io/rpc"
    );

    this.l1Provider = new hre.ethers.JsonRpcProvider(process.env.SEPOLIA_RPC_URL);

    console.log("✅ Providers initialized\n");
  }

  async checkSystemHealth() {
    console.log(`\n${"=".repeat(60)}`);
    console.log(`🏥 Health Check - ${new Date().toISOString()}`);
    console.log("=".repeat(60) + "\n");

    const healthStatus = {
      timestamp: new Date().toISOString(),
      l1: {},
      l2: {},
      alerts: [],
      overall: "HEALTHY",
    };

    // Check L1 components
    await this.checkL1Health(healthStatus);

    // Check L2 components
    await this.checkL2Health(healthStatus);

    // Determine overall health
    if (healthStatus.alerts.some(a => a.severity === "CRITICAL")) {
      healthStatus.overall = "CRITICAL";
    } else if (healthStatus.alerts.some(a => a.severity === "WARNING")) {
      healthStatus.overall = "WARNING";
    }

    // Send alerts if needed
    if (healthStatus.alerts.length > 0) {
      await this.sendAlerts(healthStatus.alerts);
    }

    // Print summary
    this.printHealthSummary(healthStatus);

    return healthStatus;
  }

  async checkL1Health(healthStatus) {
    console.log("📍 Checking L1 Components...");

    try {
      const l1RegistryAddress = process.env.L1_STATE_REGISTRY;
      const l1VaultAddress = process.env.L1_COLLATERAL_VAULT;

      if (!l1RegistryAddress) {
        console.log("⚠️  L1_STATE_REGISTRY not configured, skipping L1 checks");
        return;
      }

      // Check State Registry
      const L1StateRegistry = await hre.ethers.getContractFactory("L1StateRegistry");
      const registry = L1StateRegistry.attach(l1RegistryAddress).connect(this.l1Provider);

      const [latestRoot, latestBlock, timestamp] = await registry.getLatestState();
      const isStateFresh = await registry.isStateFresh();
      const timeSinceLast = await registry.timeSinceLastSubmission();
      const isPaused = await registry.emergencyPaused();

      healthStatus.l1.stateRegistry = {
        latestBlock: latestBlock.toString(),
        latestRoot,
        timestamp: new Date(Number(timestamp) * 1000).toISOString(),
        isStateFresh,
        timeSinceLastSubmission: Number(timeSinceLast),
        emergencyPaused: isPaused,
      };

      console.log("├─ State Registry:");
      console.log("│  ├─ Latest Block:", latestBlock.toString());
      console.log("│  ├─ Is Fresh:", isStateFresh);
      console.log("│  ├─ Time Since Last:", timeSinceLast.toString(), "seconds");
      console.log("│  └─ Emergency Paused:", isPaused);

      // Alert if state is stale
      if (!isStateFresh) {
        this.addAlert(healthStatus, {
          component: "L1StateRegistry",
          severity: "WARNING",
          message: `State is stale (> ${THRESHOLDS.MAX_STATE_AGE_HOURS} hours old)`,
          data: { timeSinceLast: Number(timeSinceLast) },
        });
      }

      // Alert if emergency paused
      if (isPaused) {
        this.addAlert(healthStatus, {
          component: "L1StateRegistry",
          severity: "CRITICAL",
          message: "System is in EMERGENCY PAUSE mode",
        });
      }

      // Check Collateral Vault if address provided
      if (l1VaultAddress) {
        const CollateralVaultL1 = await hre.ethers.getContractFactory("CollateralVaultL1");
        const vault = CollateralVaultL1.attach(l1VaultAddress).connect(this.l1Provider);

        const [totalLocked, contractBalance] = await vault.getVaultStats();
        const vaultPaused = await vault.emergencyPaused();

        healthStatus.l1.collateralVault = {
          totalLocked: hre.ethers.formatEther(totalLocked),
          contractBalance: hre.ethers.formatEther(contractBalance),
          emergencyPaused: vaultPaused,
        };

        console.log("├─ Collateral Vault:");
        console.log("│  ├─ Total Locked:", hre.ethers.formatEther(totalLocked), "tokens");
        console.log("│  ├─ Contract Balance:", hre.ethers.formatEther(contractBalance), "tokens");
        console.log("│  └─ Emergency Paused:", vaultPaused);

        // Alert if balances don't match
        if (totalLocked > contractBalance) {
          this.addAlert(healthStatus, {
            component: "CollateralVaultL1",
            severity: "CRITICAL",
            message: "Vault accounting mismatch - locked > balance",
            data: {
              totalLocked: hre.ethers.formatEther(totalLocked),
              contractBalance: hre.ethers.formatEther(contractBalance),
            },
          });
        }

        if (vaultPaused) {
          this.addAlert(healthStatus, {
            component: "CollateralVaultL1",
            severity: "CRITICAL",
            message: "Vault is in EMERGENCY PAUSE mode",
          });
        }
      }

      console.log("└─ ✅ L1 checks complete\n");
    } catch (error) {
      console.error("❌ Error checking L1 health:", error.message);
      this.addAlert(healthStatus, {
        component: "L1",
        severity: "CRITICAL",
        message: `L1 health check failed: ${error.message}`,
      });
    }
  }

  async checkL2Health(healthStatus) {
    console.log("📍 Checking L2 Components...");

    try {
      const l2AggregatorAddress = process.env.L2_STATE_AGGREGATOR;

      if (!l2AggregatorAddress) {
        console.log("⚠️  L2_STATE_AGGREGATOR not configured, skipping L2 checks");
        return;
      }

      const L2StateAggregator = await hre.ethers.getContractFactory("L2StateAggregator");
      const aggregator = L2StateAggregator.attach(l2AggregatorAddress).connect(this.l2Provider);

      // Get system state
      const state = await aggregator.getSystemState();
      const canSubmit = await aggregator.canSubmitToL1();

      healthStatus.l2.stateAggregator = {
        totalCollateral: hre.ethers.formatEther(state.totalCollateral),
        totalDebt: hre.ethers.formatUnits(state.totalDebt, 6),
        activePositions: state.activePositions.toString(),
        totalOrders: state.totalOrders.toString(),
        blockNumber: state.blockNumber.toString(),
        canSubmit,
      };

      console.log("├─ State Aggregator:");
      console.log("│  ├─ Total Collateral:", hre.ethers.formatEther(state.totalCollateral), "tokens");
      console.log("│  ├─ Total Debt:", hre.ethers.formatUnits(state.totalDebt, 6), "LUSD");
      console.log("│  ├─ Active Positions:", state.activePositions.toString());
      console.log("│  ├─ Total Orders:", state.totalOrders.toString());
      console.log("│  └─ Can Submit:", canSubmit);

      // Calculate system collateral ratio
      if (state.totalDebt > 0n) {
        const collateralRatio = (Number(state.totalCollateral) / Number(state.totalDebt)) * 100;

        healthStatus.l2.collateralRatio = collateralRatio.toFixed(2);

        console.log("├─ System Metrics:");
        console.log("│  └─ Collateral Ratio:", collateralRatio.toFixed(2) + "%");

        // Alert if ratio is low
        if (collateralRatio < THRESHOLDS.CRITICAL_COLLATERAL_RATIO) {
          this.addAlert(healthStatus, {
            component: "L2System",
            severity: "CRITICAL",
            message: `System collateral ratio critically low: ${collateralRatio.toFixed(2)}%`,
            data: { ratio: collateralRatio },
          });
        } else if (collateralRatio < THRESHOLDS.MIN_COLLATERAL_RATIO) {
          this.addAlert(healthStatus, {
            component: "L2System",
            severity: "WARNING",
            message: `System collateral ratio below minimum: ${collateralRatio.toFixed(2)}%`,
            data: { ratio: collateralRatio },
          });
        }
      }

      console.log("└─ ✅ L2 checks complete\n");
    } catch (error) {
      console.error("❌ Error checking L2 health:", error.message);
      this.addAlert(healthStatus, {
        component: "L2",
        severity: "CRITICAL",
        message: `L2 health check failed: ${error.message}`,
      });
    }
  }

  addAlert(healthStatus, alert) {
    const alertKey = `${alert.component}-${alert.message}`;
    const lastAlert = this.lastAlertTime.get(alertKey);

    // Check cooldown
    if (lastAlert && Date.now() - lastAlert < this.ALERT_COOLDOWN) {
      return; // Skip duplicate alert
    }

    healthStatus.alerts.push(alert);
    this.lastAlertTime.set(alertKey, Date.now());
  }

  async sendAlerts(alerts) {
    if (!ALERT_WEBHOOK) {
      console.log("⚠️  ALERT_WEBHOOK_URL not configured, skipping webhook alerts");
      return;
    }

    for (const alert of alerts) {
      try {
        const message = {
          content: `**[${alert.severity}]** ${alert.component}: ${alert.message}`,
          embeds: alert.data ? [{
            fields: Object.entries(alert.data).map(([key, value]) => ({
              name: key,
              value: String(value),
              inline: true,
            })),
          }] : [],
        };

        // Send to webhook (Discord/Slack)
        await fetch(ALERT_WEBHOOK, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(message),
        });

        console.log(`📢 Alert sent: [${alert.severity}] ${alert.component}`);
      } catch (error) {
        console.error("❌ Failed to send alert:", error.message);
      }
    }
  }

  printHealthSummary(healthStatus) {
    console.log("\n" + "=".repeat(60));
    console.log("📊 HEALTH SUMMARY");
    console.log("=".repeat(60));

    const statusEmoji = {
      HEALTHY: "✅",
      WARNING: "⚠️",
      CRITICAL: "🚨",
    };

    console.log(`\nOverall Status: ${statusEmoji[healthStatus.overall]} ${healthStatus.overall}`);

    if (healthStatus.alerts.length > 0) {
      console.log("\n🔔 Active Alerts:");
      healthStatus.alerts.forEach((alert, i) => {
        console.log(`  ${i + 1}. [${alert.severity}] ${alert.component}: ${alert.message}`);
      });
    } else {
      console.log("\n✅ No alerts - All systems operational");
    }

    console.log("\n" + "=".repeat(60) + "\n");
  }

  sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  async monitorLoop() {
    while (true) {
      try {
        await this.checkSystemHealth();
      } catch (error) {
        console.error("💥 Error in monitor loop:", error);
      }

      console.log(`⏸️  Next check in ${CHECK_INTERVAL / 60000} minutes...\n`);
      await this.sleep(CHECK_INTERVAL);
    }
  }
}

// Main execution
async function main() {
  const monitor = new HealthMonitor();

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

export default HealthMonitor;
