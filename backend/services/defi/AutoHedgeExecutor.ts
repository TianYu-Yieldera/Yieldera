/**
 * Auto Hedge Executor Service
 *
 * Automatically executes hedging strategies to protect user positions from liquidation
 * Connects AI risk predictions with on-chain execution
 *
 * Features:
 * - Real-time risk monitoring
 * - Automated position adjustment
 * - Multi-protocol hedge execution
 * - Gas optimization
 * - Slippage protection
 * - Emergency intervention
 */

import { ethers } from 'ethers';
import { Pool } from 'pg';
import axios from 'axios';

// Import protocol ABIs
import AavePoolABI from '../../abis/AavePool.json';
import CompoundCometABI from '../../abis/CompoundComet.json';
import UniswapRouterABI from '../../abis/UniswapRouter.json';
import GMXRouterABI from '../../abis/GMXRouter.json';

export enum HedgeStrategy {
  ADD_COLLATERAL = 'add_collateral',        // Add more collateral
  REDUCE_DEBT = 'reduce_debt',              // Repay debt
  CLOSE_POSITION = 'close_position',        // Full position closure
  PARTIAL_DELEVERAGE = 'partial_deleverage', // Reduce leverage
  SWAP_COLLATERAL = 'swap_collateral',      // Swap to safer asset
  OPEN_HEDGE_POSITION = 'open_hedge_position', // Open opposing position
}

export enum ExecutionStatus {
  PENDING = 'pending',
  SIMULATING = 'simulating',
  EXECUTING = 'executing',
  COMPLETED = 'completed',
  FAILED = 'failed',
  CANCELLED = 'cancelled',
}

export interface RiskSignal {
  userId: string;
  protocol: 'aave' | 'compound' | 'uniswap' | 'gmx';
  positionId: string;
  currentHealthFactor: number;
  predictedHealthFactor: number;
  hoursUntilLiquidation: number;
  riskScore: number; // 0-100
  recommendedStrategy: HedgeStrategy;
  urgency: 'low' | 'medium' | 'high' | 'critical';
}

export interface HedgeExecution {
  id?: string;
  userId: string;
  riskSignal: RiskSignal;
  strategy: HedgeStrategy;
  status: ExecutionStatus;
  simulatedOutcome?: SimulationResult;
  txHash?: string;
  gasUsed?: bigint;
  executedAt?: Date;
  completedAt?: Date;
  error?: string;
  resultingHealthFactor?: number;
}

export interface SimulationResult {
  success: boolean;
  estimatedGas: bigint;
  estimatedCost: number; // USD
  newHealthFactor: number;
  slippage: number;
  recommendation: string;
}

export class AutoHedgeExecutor {
  private provider: ethers.Provider;
  private wallet: ethers.Wallet;
  private db: Pool;
  private aiServiceUrl: string;

  // Protocol contracts
  private aavePool: ethers.Contract;
  private compoundComet: ethers.Contract;
  private uniswapRouter: ethers.Contract;
  private gmxRouter: ethers.Contract;

  // Configuration
  private readonly MIN_HEALTH_FACTOR = 1.5;
  private readonly TARGET_HEALTH_FACTOR = 2.0;
  private readonly MAX_SLIPPAGE = 0.02; // 2%
  private readonly MAX_GAS_PRICE = ethers.parseUnits('100', 'gwei');

  // Execution queue
  private executionQueue: HedgeExecution[] = [];
  private isProcessing = false;

  constructor(
    rpcUrl: string,
    privateKey: string,
    dbPool: Pool,
    aiServiceUrl = 'http://localhost:8084'
  ) {
    this.provider = new ethers.JsonRpcProvider(rpcUrl);
    this.wallet = new ethers.Wallet(privateKey, this.provider);
    this.db = dbPool;
    this.aiServiceUrl = aiServiceUrl;

    // Initialize protocol contracts with signer
    this.aavePool = new ethers.Contract(
      process.env.AAVE_POOL_ADDRESS || '0x87870Bca3F3fD6335C3F4ce8392D69350B4fA4E2',
      AavePoolABI,
      this.wallet
    );

    this.compoundComet = new ethers.Contract(
      process.env.COMPOUND_COMET_ADDRESS || '0xc3d688B66703497DAA19211EEdff47f25384cdc3',
      CompoundCometABI,
      this.wallet
    );

    this.uniswapRouter = new ethers.Contract(
      process.env.UNISWAP_ROUTER_ADDRESS || '0xE592427A0AEce92De3Edee1F18E0157C05861564',
      UniswapRouterABI,
      this.wallet
    );

    this.gmxRouter = new ethers.Contract(
      process.env.GMX_ROUTER_ADDRESS || '0x...',
      GMXRouterABI,
      this.wallet
    );

    // Start monitoring loop
    this.startMonitoring();
  }

  /**
   * Start continuous risk monitoring
   */
  private startMonitoring(): void {
    console.log('🔍 Auto hedge monitoring started');

    // Check for high-risk positions every 60 seconds
    setInterval(async () => {
      try {
        await this.checkAllPositions();
        await this.processQueue();
      } catch (error) {
        console.error('Monitoring error:', error);
      }
    }, 60000);
  }

  /**
   * Check all positions for risk signals
   */
  private async checkAllPositions(): Promise<void> {
    // Get all active positions from database
    const result = await this.db.query(
      `SELECT DISTINCT user_id, protocol, position_id, health_factor, risk_score
       FROM defi_positions
       WHERE health_factor < $1 OR risk_score > $2`,
      [this.MIN_HEALTH_FACTOR, 70]
    );

    for (const row of result.rows) {
      // Get AI risk assessment
      const riskSignal = await this.getAIRiskSignal(
        row.user_id,
        row.protocol,
        row.position_id
      );

      if (riskSignal && this.shouldExecuteHedge(riskSignal)) {
        await this.queueHedgeExecution(riskSignal);
      }
    }
  }

  /**
   * Get AI risk signal for a position
   */
  private async getAIRiskSignal(
    userId: string,
    protocol: string,
    positionId: string
  ): Promise<RiskSignal | null> {
    try {
      const response = await axios.post(`${this.aiServiceUrl}/risk/signal`, {
        userId,
        protocol,
        positionId,
      });

      const data = response.data;

      return {
        userId,
        protocol: protocol as any,
        positionId,
        currentHealthFactor: data.current_health_factor,
        predictedHealthFactor: data.predicted_health_factor,
        hoursUntilLiquidation: data.hours_until_liquidation,
        riskScore: data.risk_score,
        recommendedStrategy: data.recommended_strategy as HedgeStrategy,
        urgency: data.urgency,
      };
    } catch (error) {
      console.error('Failed to get AI risk signal:', error);
      return null;
    }
  }

  /**
   * Determine if hedge should be executed
   */
  private shouldExecuteHedge(signal: RiskSignal): boolean {
    // Check if user has auto-hedge enabled
    // Check if urgency warrants immediate action
    // Check if we haven't already queued this position

    if (signal.urgency === 'critical' && signal.hoursUntilLiquidation < 24) {
      return true;
    }

    if (signal.urgency === 'high' && signal.currentHealthFactor < 1.3) {
      return true;
    }

    return false;
  }

  /**
   * Queue hedge execution
   */
  private async queueHedgeExecution(signal: RiskSignal): Promise<void> {
    // Check if user has auto-hedge enabled
    const autoHedgeEnabled = await this.isAutoHedgeEnabled(signal.userId);

    if (!autoHedgeEnabled) {
      console.log(`Auto-hedge disabled for user ${signal.userId}`);
      return;
    }

    // Check if already queued
    const alreadyQueued = this.executionQueue.some(
      exec => exec.userId === signal.userId && exec.riskSignal.positionId === signal.positionId
    );

    if (alreadyQueued) {
      console.log(`Position ${signal.positionId} already queued`);
      return;
    }

    const execution: HedgeExecution = {
      userId: signal.userId,
      riskSignal: signal,
      strategy: signal.recommendedStrategy,
      status: ExecutionStatus.PENDING,
    };

    this.executionQueue.push(execution);

    // Save to database
    await this.saveExecution(execution);

    console.log(`✓ Queued hedge execution for ${signal.positionId}`);
  }

  /**
   * Process execution queue
   */
  private async processQueue(): Promise<void> {
    if (this.isProcessing || this.executionQueue.length === 0) {
      return;
    }

    this.isProcessing = true;

    try {
      // Sort by urgency
      this.executionQueue.sort((a, b) => {
        const urgencyOrder = { critical: 0, high: 1, medium: 2, low: 3 };
        return urgencyOrder[a.riskSignal.urgency] - urgencyOrder[b.riskSignal.urgency];
      });

      // Process one execution at a time
      const execution = this.executionQueue.shift();

      if (execution) {
        await this.executeHedge(execution);
      }
    } finally {
      this.isProcessing = false;
    }
  }

  /**
   * Execute hedge strategy
   */
  private async executeHedge(execution: HedgeExecution): Promise<void> {
    console.log(`⚡ Executing hedge for ${execution.riskSignal.positionId}`);

    try {
      // Step 1: Simulate execution
      execution.status = ExecutionStatus.SIMULATING;
      await this.updateExecution(execution);

      const simulation = await this.simulateHedge(execution);
      execution.simulatedOutcome = simulation;

      if (!simulation.success) {
        execution.status = ExecutionStatus.FAILED;
        execution.error = simulation.recommendation;
        await this.updateExecution(execution);
        return;
      }

      // Step 2: Check gas price
      const gasPrice = await this.provider.getFeeData();
      if (gasPrice.gasPrice && gasPrice.gasPrice > this.MAX_GAS_PRICE) {
        console.log('Gas price too high, delaying execution');
        this.executionQueue.unshift(execution); // Put back in queue
        return;
      }

      // Step 3: Execute on-chain
      execution.status = ExecutionStatus.EXECUTING;
      execution.executedAt = new Date();
      await this.updateExecution(execution);

      const txHash = await this.executeStrategy(execution);
      execution.txHash = txHash;

      // Step 4: Wait for confirmation
      const receipt = await this.provider.waitForTransaction(txHash);

      if (receipt && receipt.status === 1) {
        execution.status = ExecutionStatus.COMPLETED;
        execution.gasUsed = receipt.gasUsed;
        execution.completedAt = new Date();

        // Verify new health factor
        const newHealthFactor = await this.getHealthFactor(
          execution.userId,
          execution.riskSignal.protocol
        );
        execution.resultingHealthFactor = newHealthFactor;

        console.log(`✓ Hedge executed successfully. New health factor: ${newHealthFactor}`);
      } else {
        execution.status = ExecutionStatus.FAILED;
        execution.error = 'Transaction reverted';
      }
    } catch (error: any) {
      execution.status = ExecutionStatus.FAILED;
      execution.error = error.message;
      console.error('Hedge execution failed:', error);
    } finally {
      await this.updateExecution(execution);
    }
  }

  /**
   * Simulate hedge execution
   */
  private async simulateHedge(execution: HedgeExecution): Promise<SimulationResult> {
    const { strategy, riskSignal } = execution;

    try {
      switch (strategy) {
        case HedgeStrategy.ADD_COLLATERAL:
          return await this.simulateAddCollateral(riskSignal);
        case HedgeStrategy.REDUCE_DEBT:
          return await this.simulateReduceDebt(riskSignal);
        case HedgeStrategy.PARTIAL_DELEVERAGE:
          return await this.simulatePartialDeleverage(riskSignal);
        case HedgeStrategy.CLOSE_POSITION:
          return await this.simulateClosePosition(riskSignal);
        default:
          return {
            success: false,
            estimatedGas: 0n,
            estimatedCost: 0,
            newHealthFactor: riskSignal.currentHealthFactor,
            slippage: 0,
            recommendation: 'Strategy not implemented',
          };
      }
    } catch (error: any) {
      return {
        success: false,
        estimatedGas: 0n,
        estimatedCost: 0,
        newHealthFactor: riskSignal.currentHealthFactor,
        slippage: 0,
        recommendation: `Simulation failed: ${error.message}`,
      };
    }
  }

  /**
   * Simulate adding collateral
   */
  private async simulateAddCollateral(signal: RiskSignal): Promise<SimulationResult> {
    // Calculate how much collateral needed to reach target health factor
    const currentHF = signal.currentHealthFactor;
    const targetHF = this.TARGET_HEALTH_FACTOR;

    // Simplified calculation (actual implementation would be more complex)
    const collateralNeeded = ((targetHF - currentHF) / currentHF) * 1000; // USD

    // Estimate gas
    const estimatedGas = 200000n; // Typical for supply operation
    const gasPrice = (await this.provider.getFeeData()).gasPrice || 0n;
    const estimatedCost = Number(estimatedGas * gasPrice) / 1e18;

    return {
      success: collateralNeeded < 10000, // Max $10k auto-hedge
      estimatedGas,
      estimatedCost,
      newHealthFactor: targetHF,
      slippage: 0,
      recommendation: `Add $${collateralNeeded.toFixed(2)} collateral`,
    };
  }

  /**
   * Simulate reducing debt
   */
  private async simulateReduceDebt(signal: RiskSignal): Promise<SimulationResult> {
    const currentHF = signal.currentHealthFactor;
    const targetHF = this.TARGET_HEALTH_FACTOR;

    const debtToRepay = ((targetHF - currentHF) / currentHF) * 800; // USD

    const estimatedGas = 250000n;
    const gasPrice = (await this.provider.getFeeData()).gasPrice || 0n;
    const estimatedCost = Number(estimatedGas * gasPrice) / 1e18;

    return {
      success: debtToRepay < 5000,
      estimatedGas,
      estimatedCost,
      newHealthFactor: targetHF,
      slippage: 0.01,
      recommendation: `Repay $${debtToRepay.toFixed(2)} debt`,
    };
  }

  /**
   * Simulate partial deleverage
   */
  private async simulatePartialDeleverage(signal: RiskSignal): Promise<SimulationResult> {
    const estimatedGas = 400000n; // More complex operation
    const gasPrice = (await this.provider.getFeeData()).gasPrice || 0n;
    const estimatedCost = Number(estimatedGas * gasPrice) / 1e18;

    return {
      success: true,
      estimatedGas,
      estimatedCost,
      newHealthFactor: this.TARGET_HEALTH_FACTOR,
      slippage: 0.015,
      recommendation: 'Reduce position by 30%',
    };
  }

  /**
   * Simulate close position
   */
  private async simulateClosePosition(signal: RiskSignal): Promise<SimulationResult> {
    const estimatedGas = 500000n;
    const gasPrice = (await this.provider.getFeeData()).gasPrice || 0n;
    const estimatedCost = Number(estimatedGas * gasPrice) / 1e18;

    return {
      success: true,
      estimatedGas,
      estimatedCost,
      newHealthFactor: 999, // No position = infinite health factor
      slippage: 0.02,
      recommendation: 'Close entire position',
    };
  }

  /**
   * Execute strategy on-chain
   */
  private async executeStrategy(execution: HedgeExecution): Promise<string> {
    const { strategy, riskSignal } = execution;

    switch (strategy) {
      case HedgeStrategy.ADD_COLLATERAL:
        return await this.executeAddCollateral(riskSignal);
      case HedgeStrategy.REDUCE_DEBT:
        return await this.executeReduceDebt(riskSignal);
      case HedgeStrategy.PARTIAL_DELEVERAGE:
        return await this.executePartialDeleverage(riskSignal);
      case HedgeStrategy.CLOSE_POSITION:
        return await this.executeClosePosition(riskSignal);
      default:
        throw new Error('Strategy not implemented');
    }
  }

  /**
   * Execute add collateral on Aave
   */
  private async executeAddCollateral(signal: RiskSignal): Promise<string> {
    if (signal.protocol !== 'aave') {
      throw new Error('Protocol not supported for this strategy');
    }

    // This is a simplified example
    // In production, you'd need to:
    // 1. Get user's approval to spend their tokens
    // 2. Determine optimal collateral asset
    // 3. Execute supply transaction

    const USDC = '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48'; // Example
    const amount = ethers.parseUnits('1000', 6); // $1000 USDC

    const tx = await this.aavePool.supply(
      USDC,
      amount,
      signal.userId,
      0 // referral code
    );

    return tx.hash;
  }

  /**
   * Execute reduce debt
   */
  private async executeReduceDebt(signal: RiskSignal): Promise<string> {
    if (signal.protocol === 'aave') {
      const USDC = '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48';
      const amount = ethers.parseUnits('500', 6);

      const tx = await this.aavePool.repay(
        USDC,
        amount,
        2, // variable rate mode
        signal.userId
      );

      return tx.hash;
    } else if (signal.protocol === 'compound') {
      const amount = ethers.parseUnits('500', 6);

      const tx = await this.compoundComet.supply(
        signal.userId,
        amount
      );

      return tx.hash;
    }

    throw new Error('Protocol not supported');
  }

  /**
   * Execute partial deleverage
   */
  private async executePartialDeleverage(signal: RiskSignal): Promise<string> {
    // This would involve:
    // 1. Withdraw some collateral
    // 2. Swap to debt token
    // 3. Repay debt
    // 4. Return remaining to user

    // Placeholder for complex multi-step operation
    throw new Error('Not yet implemented');
  }

  /**
   * Execute close position
   */
  private async executeClosePosition(signal: RiskSignal): Promise<string> {
    // This would involve:
    // 1. Withdraw all collateral
    // 2. Swap to debt token
    // 3. Repay all debt
    // 4. Return remaining to user

    // Placeholder for complex operation
    throw new Error('Not yet implemented');
  }

  /**
   * Helper methods
   */

  private async getHealthFactor(userId: string, protocol: string): Promise<number> {
    if (protocol === 'aave') {
      const accountData = await this.aavePool.getUserAccountData(userId);
      return Number(accountData.healthFactor) / 1e18;
    } else if (protocol === 'compound') {
      const borrowBalance = await this.compoundComet.borrowBalanceOf(userId);
      const collateralBalance = await this.compoundComet.balanceOf(userId);

      if (Number(borrowBalance) === 0) return 999;

      return (Number(collateralBalance) * 0.8) / Number(borrowBalance);
    }

    return 999;
  }

  private async isAutoHedgeEnabled(userId: string): Promise<boolean> {
    const result = await this.db.query(
      'SELECT auto_hedge_enabled FROM user_settings WHERE user_id = $1',
      [userId]
    );
    return result.rows[0]?.auto_hedge_enabled || false;
  }

  /**
   * Database operations
   */

  private async saveExecution(execution: HedgeExecution): Promise<void> {
    const result = await this.db.query(
      `INSERT INTO hedge_executions
       (user_id, risk_signal, strategy, status, created_at)
       VALUES ($1, $2, $3, $4, NOW())
       RETURNING id`,
      [
        execution.userId,
        JSON.stringify(execution.riskSignal),
        execution.strategy,
        execution.status,
      ]
    );

    execution.id = result.rows[0].id;
  }

  private async updateExecution(execution: HedgeExecution): Promise<void> {
    await this.db.query(
      `UPDATE hedge_executions
       SET status = $1,
           simulated_outcome = $2,
           tx_hash = $3,
           gas_used = $4,
           executed_at = $5,
           completed_at = $6,
           error = $7,
           resulting_health_factor = $8,
           updated_at = NOW()
       WHERE id = $9`,
      [
        execution.status,
        JSON.stringify(execution.simulatedOutcome),
        execution.txHash,
        execution.gasUsed?.toString(),
        execution.executedAt,
        execution.completedAt,
        execution.error,
        execution.resultingHealthFactor,
        execution.id,
      ]
    );
  }

  /**
   * Public API for manual hedge execution
   */
  async manualHedge(
    userId: string,
    positionId: string,
    strategy: HedgeStrategy
  ): Promise<HedgeExecution> {
    // Get current position data
    const result = await this.db.query(
      'SELECT * FROM defi_positions WHERE position_id = $1',
      [positionId]
    );

    if (result.rows.length === 0) {
      throw new Error('Position not found');
    }

    const position = result.rows[0];

    // Create manual risk signal
    const riskSignal: RiskSignal = {
      userId,
      protocol: position.protocol,
      positionId,
      currentHealthFactor: position.health_factor,
      predictedHealthFactor: position.health_factor,
      hoursUntilLiquidation: 999,
      riskScore: position.risk_score,
      recommendedStrategy: strategy,
      urgency: 'medium',
    };

    const execution: HedgeExecution = {
      userId,
      riskSignal,
      strategy,
      status: ExecutionStatus.PENDING,
    };

    await this.saveExecution(execution);
    this.executionQueue.push(execution);

    return execution;
  }

  /**
   * Get execution history for user
   */
  async getExecutionHistory(userId: string, limit = 20): Promise<HedgeExecution[]> {
    const result = await this.db.query(
      `SELECT * FROM hedge_executions
       WHERE user_id = $1
       ORDER BY created_at DESC
       LIMIT $2`,
      [userId, limit]
    );

    return result.rows.map(row => ({
      id: row.id,
      userId: row.user_id,
      riskSignal: JSON.parse(row.risk_signal),
      strategy: row.strategy,
      status: row.status,
      simulatedOutcome: row.simulated_outcome ? JSON.parse(row.simulated_outcome) : undefined,
      txHash: row.tx_hash,
      gasUsed: row.gas_used ? BigInt(row.gas_used) : undefined,
      executedAt: row.executed_at,
      completedAt: row.completed_at,
      error: row.error,
      resultingHealthFactor: row.resulting_health_factor,
    }));
  }
}

// Export singleton instance (use with caution - requires private key)
export const autoHedgeExecutor = process.env.HEDGE_PRIVATE_KEY
  ? new AutoHedgeExecutor(
      process.env.RPC_URL || 'https://eth.llamarpc.com',
      process.env.HEDGE_PRIVATE_KEY,
      new Pool({
        connectionString: process.env.DATABASE_URL,
      })
    )
  : null;
