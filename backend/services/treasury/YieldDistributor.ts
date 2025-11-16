/**
 * Yield Distribution Service
 *
 * Automated daily distribution of Treasury bond yields to user accounts
 * Handles on-chain minting and gas optimization through batching
 *
 * Features:
 * - Automated daily distribution
 * - Gas-optimized batch processing
 * - Merkle tree proofs for verification
 * - Audit trail and compliance
 * - Emergency pause mechanism
 * - Retry logic for failed distributions
 */

import { ethers } from 'ethers';
import { Pool } from 'pg';
import { MerkleTree } from 'merkletreejs';
import keccak256 from 'keccak256';

// Treasury token contract interface
const TREASURY_TOKEN_ABI = [
  'function mint(address to, uint256 amount) external',
  'function batchMint(address[] calldata recipients, uint256[] calldata amounts) external',
  'function balanceOf(address owner) view returns (uint256)',
  'function totalSupply() view returns (uint256)',
  'event YieldDistributed(address indexed recipient, uint256 amount, uint256 timestamp)',
];

export interface YieldDistribution {
  id?: string;
  userId: string;
  bondType: string;
  yieldAmount: number; // USD
  tokenAmount: bigint;  // On-chain token amount
  distributionDate: Date;
  txHash?: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  retryCount: number;
  error?: string;
}

export interface BatchDistribution {
  id?: string;
  batchNumber: number;
  distributionDate: Date;
  totalRecipients: number;
  totalYieldUsd: number;
  totalTokenAmount: bigint;
  merkleRoot: string;
  txHash?: string;
  gasUsed?: bigint;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  executedAt?: Date;
}

export interface DistributionStats {
  date: Date;
  totalDistributed: number;
  recipientCount: number;
  successRate: number;
  averageYieldPerUser: number;
  totalGasCost: number;
}

export class YieldDistributor {
  private provider: ethers.Provider;
  private wallet: ethers.Wallet;
  private db: Pool;

  // Treasury token contracts by bond type
  private tokenContracts: Map<string, ethers.Contract> = new Map();

  // Configuration
  private readonly BATCH_SIZE = 100; // Process 100 users per batch
  private readonly MAX_RETRIES = 3;
  private readonly RETRY_DELAY = 60000; // 1 minute
  private readonly MIN_DISTRIBUTION_AMOUNT = 0.01; // Minimum $0.01

  // State
  private isDistributing = false;
  private isPaused = false;

  constructor(
    rpcUrl: string,
    privateKey: string,
    dbPool: Pool
  ) {
    this.provider = new ethers.JsonRpcProvider(rpcUrl);
    this.wallet = new ethers.Wallet(privateKey, this.provider);
    this.db = dbPool;

    this.initializeTokenContracts();

    // Schedule daily distribution at 00:00 UTC
    this.scheduleDailyDistribution();
  }

  /**
   * Initialize treasury token contracts
   */
  private initializeTokenContracts(): void {
    const bondTypes = [
      'TBILL_3M', 'TBILL_6M',
      'TNOTE_2Y', 'TNOTE_5Y', 'TNOTE_10Y',
      'TBOND_20Y', 'TBOND_30Y',
    ];

    for (const bondType of bondTypes) {
      const address = process.env[`${bondType}_TOKEN_ADDRESS`];
      if (address) {
        const contract = new ethers.Contract(
          address,
          TREASURY_TOKEN_ABI,
          this.wallet
        );
        this.tokenContracts.set(bondType, contract);
      }
    }
  }

  /**
   * Schedule daily distribution
   */
  private scheduleDailyDistribution(): void {
    // Run every hour to check if it's time to distribute
    setInterval(async () => {
      const now = new Date();
      const hour = now.getUTCHours();
      const minute = now.getUTCMinutes();

      // Run at 00:00 UTC (can be configured)
      if (hour === 0 && minute < 10 && !this.isDistributing) {
        console.log('🕐 Starting scheduled daily distribution...');
        await this.distributeDailyYields();
      }
    }, 600000); // Check every 10 minutes
  }

  /**
   * Main distribution function - called daily
   */
  async distributeDailyYields(): Promise<void> {
    if (this.isDistributing) {
      console.log('Distribution already in progress');
      return;
    }

    if (this.isPaused) {
      console.log('Distribution is paused');
      return;
    }

    this.isDistributing = true;

    try {
      console.log('📊 Starting daily yield distribution...');

      // Get all pending distributions from yesterday's calculations
      const distributions = await this.getPendingDistributions();

      if (distributions.length === 0) {
        console.log('No pending distributions');
        return;
      }

      console.log(`Found ${distributions.length} pending distributions`);

      // Group by bond type for efficient batch processing
      const byBondType = this.groupByBondType(distributions);

      // Process each bond type
      for (const [bondType, userDistributions] of byBondType.entries()) {
        await this.distributeBondType(bondType, userDistributions);
      }

      // Generate distribution report
      await this.generateDistributionReport();

      console.log('✓ Daily distribution completed successfully');
    } catch (error) {
      console.error('Distribution failed:', error);
      await this.handleDistributionError(error);
    } finally {
      this.isDistributing = false;
    }
  }

  /**
   * Distribute yields for a specific bond type
   */
  private async distributeBondType(
    bondType: string,
    distributions: YieldDistribution[]
  ): Promise<void> {
    console.log(`Processing ${distributions.length} distributions for ${bondType}`);

    const contract = this.tokenContracts.get(bondType);
    if (!contract) {
      throw new Error(`No contract found for ${bondType}`);
    }

    // Split into batches
    const batches = this.createBatches(distributions, this.BATCH_SIZE);

    for (let i = 0; i < batches.length; i++) {
      const batch = batches[i];
      console.log(`Processing batch ${i + 1}/${batches.length}`);

      await this.processBatch(bondType, batch, i + 1);

      // Small delay between batches to avoid overwhelming the network
      if (i < batches.length - 1) {
        await this.delay(5000); // 5 second delay
      }
    }
  }

  /**
   * Process a single batch of distributions
   */
  private async processBatch(
    bondType: string,
    distributions: YieldDistribution[],
    batchNumber: number
  ): Promise<void> {
    const contract = this.tokenContracts.get(bondType);
    if (!contract) {
      throw new Error(`No contract found for ${bondType}`);
    }

    // Prepare batch data
    const recipients: string[] = [];
    const amounts: bigint[] = [];
    let totalAmount = 0n;

    for (const dist of distributions) {
      recipients.push(dist.userId);
      amounts.push(dist.tokenAmount);
      totalAmount += dist.tokenAmount;
    }

    // Create Merkle tree for verification
    const merkleTree = this.createMerkleTree(recipients, amounts);
    const merkleRoot = merkleTree.getHexRoot();

    // Create batch record
    const batchRecord: BatchDistribution = {
      batchNumber,
      distributionDate: new Date(),
      totalRecipients: recipients.length,
      totalYieldUsd: distributions.reduce((sum, d) => sum + d.yieldAmount, 0),
      totalTokenAmount: totalAmount,
      merkleRoot,
      status: 'processing',
    };

    await this.saveBatchRecord(batchRecord);

    try {
      // Execute batch mint on-chain
      console.log(`Minting ${ethers.formatEther(totalAmount)} tokens to ${recipients.length} recipients`);

      const tx = await contract.batchMint(recipients, amounts, {
        gasLimit: 500000 + (recipients.length * 50000), // Dynamic gas estimation
      });

      console.log(`Transaction submitted: ${tx.hash}`);
      batchRecord.txHash = tx.hash;
      await this.updateBatchRecord(batchRecord);

      // Wait for confirmation
      const receipt = await tx.wait();

      if (receipt && receipt.status === 1) {
        batchRecord.status = 'completed';
        batchRecord.gasUsed = receipt.gasUsed;
        batchRecord.executedAt = new Date();

        // Mark individual distributions as completed
        for (const dist of distributions) {
          dist.status = 'completed';
          dist.txHash = tx.hash;
          await this.updateDistribution(dist);
        }

        console.log(`✓ Batch ${batchNumber} completed. Gas used: ${receipt.gasUsed}`);
      } else {
        throw new Error('Transaction failed');
      }
    } catch (error: any) {
      console.error(`Batch ${batchNumber} failed:`, error);
      batchRecord.status = 'failed';

      // Mark distributions for retry
      for (const dist of distributions) {
        dist.status = 'failed';
        dist.error = error.message;
        dist.retryCount++;
        await this.updateDistribution(dist);
      }
    } finally {
      await this.updateBatchRecord(batchRecord);
    }
  }

  /**
   * Retry failed distributions
   */
  async retryFailedDistributions(): Promise<void> {
    console.log('Retrying failed distributions...');

    const result = await this.db.query(
      `SELECT * FROM yield_distributions
       WHERE status = 'failed' AND retry_count < $1`,
      [this.MAX_RETRIES]
    );

    const distributions = result.rows.map(row => this.mapRowToDistribution(row));

    if (distributions.length === 0) {
      console.log('No failed distributions to retry');
      return;
    }

    const byBondType = this.groupByBondType(distributions);

    for (const [bondType, dists] of byBondType.entries()) {
      await this.distributeBondType(bondType, dists);
    }
  }

  /**
   * Manual distribution for specific user (for testing or special cases)
   */
  async manualDistribute(
    userId: string,
    bondType: string,
    yieldAmount: number
  ): Promise<string> {
    if (yieldAmount < this.MIN_DISTRIBUTION_AMOUNT) {
      throw new Error('Amount below minimum distribution threshold');
    }

    const contract = this.tokenContracts.get(bondType);
    if (!contract) {
      throw new Error(`No contract found for ${bondType}`);
    }

    // Convert USD to token amount (1:1 for simplicity)
    const tokenAmount = ethers.parseUnits(yieldAmount.toString(), 18);

    const distribution: YieldDistribution = {
      userId,
      bondType,
      yieldAmount,
      tokenAmount,
      distributionDate: new Date(),
      status: 'processing',
      retryCount: 0,
    };

    await this.saveDistribution(distribution);

    try {
      const tx = await contract.mint(userId, tokenAmount);
      distribution.txHash = tx.hash;

      const receipt = await tx.wait();

      if (receipt && receipt.status === 1) {
        distribution.status = 'completed';
        console.log(`✓ Manual distribution completed: ${tx.hash}`);
      } else {
        throw new Error('Transaction failed');
      }
    } catch (error: any) {
      distribution.status = 'failed';
      distribution.error = error.message;
      throw error;
    } finally {
      await this.updateDistribution(distribution);
    }

    return distribution.txHash!;
  }

  /**
   * Get distribution statistics
   */
  async getDistributionStats(
    startDate?: Date,
    endDate?: Date
  ): Promise<DistributionStats[]> {
    const query = `
      SELECT
        DATE(distribution_date) as date,
        SUM(yield_amount) as total_distributed,
        COUNT(DISTINCT user_id) as recipient_count,
        COUNT(CASE WHEN status = 'completed' THEN 1 END)::float / COUNT(*)::float as success_rate,
        AVG(yield_amount) as average_yield
      FROM yield_distributions
      WHERE ($1::date IS NULL OR distribution_date >= $1)
        AND ($2::date IS NULL OR distribution_date <= $2)
      GROUP BY DATE(distribution_date)
      ORDER BY date DESC
    `;

    const result = await this.db.query(query, [startDate, endDate]);

    // Get total gas cost from batches
    const gasQuery = `
      SELECT
        DATE(distribution_date) as date,
        SUM(gas_used * gas_price) as total_gas_cost
      FROM batch_distributions
      WHERE status = 'completed'
        AND ($1::date IS NULL OR distribution_date >= $1)
        AND ($2::date IS NULL OR distribution_date <= $2)
      GROUP BY DATE(distribution_date)
    `;

    const gasResult = await this.db.query(gasQuery, [startDate, endDate]);
    const gasCostMap = new Map(
      gasResult.rows.map(row => [row.date, Number(row.total_gas_cost)])
    );

    return result.rows.map(row => ({
      date: new Date(row.date),
      totalDistributed: parseFloat(row.total_distributed),
      recipientCount: parseInt(row.recipient_count),
      successRate: parseFloat(row.success_rate),
      averageYieldPerUser: parseFloat(row.average_yield),
      totalGasCost: gasCostMap.get(row.date) || 0,
    }));
  }

  /**
   * Pause/resume distribution
   */
  setPaused(paused: boolean): void {
    this.isPaused = paused;
    console.log(`Distribution ${paused ? 'paused' : 'resumed'}`);
  }

  /**
   * Helper methods
   */

  private async getPendingDistributions(): Promise<YieldDistribution[]> {
    const result = await this.db.query(
      `SELECT * FROM yield_distributions
       WHERE status = 'pending' AND yield_amount >= $1
       ORDER BY bond_type, user_id`,
      [this.MIN_DISTRIBUTION_AMOUNT]
    );

    return result.rows.map(row => this.mapRowToDistribution(row));
  }

  private groupByBondType(
    distributions: YieldDistribution[]
  ): Map<string, YieldDistribution[]> {
    const grouped = new Map<string, YieldDistribution[]>();

    for (const dist of distributions) {
      if (!grouped.has(dist.bondType)) {
        grouped.set(dist.bondType, []);
      }
      grouped.get(dist.bondType)!.push(dist);
    }

    return grouped;
  }

  private createBatches<T>(items: T[], batchSize: number): T[][] {
    const batches: T[][] = [];
    for (let i = 0; i < items.length; i += batchSize) {
      batches.push(items.slice(i, i + batchSize));
    }
    return batches;
  }

  private createMerkleTree(recipients: string[], amounts: bigint[]): MerkleTree {
    const leaves = recipients.map((recipient, i) =>
      keccak256(
        ethers.solidityPacked(
          ['address', 'uint256'],
          [recipient, amounts[i]]
        )
      )
    );

    return new MerkleTree(leaves, keccak256, { sortPairs: true });
  }

  private async delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  private mapRowToDistribution(row: any): YieldDistribution {
    return {
      id: row.id,
      userId: row.user_id,
      bondType: row.bond_type,
      yieldAmount: parseFloat(row.yield_amount),
      tokenAmount: BigInt(row.token_amount),
      distributionDate: new Date(row.distribution_date),
      txHash: row.tx_hash,
      status: row.status,
      retryCount: row.retry_count,
      error: row.error,
    };
  }

  /**
   * Database operations
   */

  private async saveDistribution(dist: YieldDistribution): Promise<void> {
    const result = await this.db.query(
      `INSERT INTO yield_distributions
       (user_id, bond_type, yield_amount, token_amount, distribution_date, status, retry_count)
       VALUES ($1, $2, $3, $4, $5, $6, $7)
       RETURNING id`,
      [
        dist.userId,
        dist.bondType,
        dist.yieldAmount,
        dist.tokenAmount.toString(),
        dist.distributionDate,
        dist.status,
        dist.retryCount,
      ]
    );
    dist.id = result.rows[0].id;
  }

  private async updateDistribution(dist: YieldDistribution): Promise<void> {
    await this.db.query(
      `UPDATE yield_distributions
       SET status = $1, tx_hash = $2, error = $3, retry_count = $4, updated_at = NOW()
       WHERE id = $5`,
      [dist.status, dist.txHash, dist.error, dist.retryCount, dist.id]
    );
  }

  private async saveBatchRecord(batch: BatchDistribution): Promise<void> {
    const result = await this.db.query(
      `INSERT INTO batch_distributions
       (batch_number, distribution_date, total_recipients, total_yield_usd,
        total_token_amount, merkle_root, status)
       VALUES ($1, $2, $3, $4, $5, $6, $7)
       RETURNING id`,
      [
        batch.batchNumber,
        batch.distributionDate,
        batch.totalRecipients,
        batch.totalYieldUsd,
        batch.totalTokenAmount.toString(),
        batch.merkleRoot,
        batch.status,
      ]
    );
    batch.id = result.rows[0].id;
  }

  private async updateBatchRecord(batch: BatchDistribution): Promise<void> {
    await this.db.query(
      `UPDATE batch_distributions
       SET status = $1, tx_hash = $2, gas_used = $3, executed_at = $4, updated_at = NOW()
       WHERE id = $5`,
      [
        batch.status,
        batch.txHash,
        batch.gasUsed?.toString(),
        batch.executedAt,
        batch.id,
      ]
    );
  }

  private async generateDistributionReport(): Promise<void> {
    const today = new Date();
    const stats = await this.getDistributionStats(today, today);

    if (stats.length > 0) {
      const stat = stats[0];

      await this.db.query(
        `INSERT INTO distribution_reports
         (date, total_distributed, recipient_count, success_rate,
          average_yield, total_gas_cost, generated_at)
         VALUES ($1, $2, $3, $4, $5, $6, NOW())`,
        [
          stat.date,
          stat.totalDistributed,
          stat.recipientCount,
          stat.successRate,
          stat.averageYieldPerUser,
          stat.totalGasCost,
        ]
      );

      console.log(`📄 Distribution report generated:
        Total: $${stat.totalDistributed.toFixed(2)}
        Recipients: ${stat.recipientCount}
        Success Rate: ${(stat.successRate * 100).toFixed(2)}%
        Gas Cost: $${stat.totalGasCost.toFixed(2)}
      `);
    }
  }

  private async handleDistributionError(error: any): Promise<void> {
    console.error('Critical distribution error:', error);

    // Log to database
    await this.db.query(
      `INSERT INTO distribution_errors
       (error_message, stack_trace, occurred_at)
       VALUES ($1, $2, NOW())`,
      [error.message, error.stack]
    );

    // Send alert (integrate with NotificationService)
    // await notificationService.sendSystemAlert(...)
  }
}

// Export singleton instance
export const yieldDistributor = process.env.DISTRIBUTOR_PRIVATE_KEY
  ? new YieldDistributor(
      process.env.RPC_URL || 'https://eth.llamarpc.com',
      process.env.DISTRIBUTOR_PRIVATE_KEY,
      new Pool({
        connectionString: process.env.DATABASE_URL,
      })
    )
  : null;

/**
 * CLI interface for manual operations
 */
export async function runDistribution() {
  if (!yieldDistributor) {
    console.error('Yield distributor not initialized. Check DISTRIBUTOR_PRIVATE_KEY');
    process.exit(1);
  }

  await yieldDistributor.distributeDailyYields();
}

export async function retryFailed() {
  if (!yieldDistributor) {
    console.error('Yield distributor not initialized');
    process.exit(1);
  }

  await yieldDistributor.retryFailedDistributions();
}
