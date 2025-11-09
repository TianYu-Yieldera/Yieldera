/**
 * Onchain Data Collector
 *
 * Collects historical and real-time data from blockchain for AI risk engine
 *
 * Phase 1: Data Layer
 * Created: 2025-11-09
 */

import { ethers } from 'ethers';
import { Pool } from 'pg';
import axios from 'axios';

// ============================================================
// TYPES
// ============================================================

export interface PricePoint {
  time: Date;
  asset: string;
  price: string;
  source: 'chainlink' | 'uniswap_twap' | 'coingecko';
  volume?: string;
  marketCap?: string;
}

export interface LiquidationEvent {
  time: Date;
  protocol: 'aave' | 'compound' | 'gmx';
  userAddress: string;
  collateralAsset: string;
  debtAsset?: string;
  collateralAmount: string;
  debtAmount?: string;
  liquidationPrice: string;
  healthFactor?: number;
  gasPrice: string;
  txHash: string;
  blockNumber: number;
}

export interface UserMetrics {
  userAddress: string;
  riskScore: number;
  avgLeverage: number;
  liquidationCount: number;
  totalPositions: number;
  maxPositionSize: string;
  preferredAssets: string[];
  preferredProtocols: string[];
  avgHealthFactor: number;
  totalVolumeUsd: string;
  lastActivity: Date;
}

export interface MarketDepth {
  time: Date;
  protocol: string;
  market: string;
  totalSupply: string;
  totalBorrow: string;
  utilizationRate: number;
  supplyAPY: number;
  borrowAPY: number;
  availableLiquidity: string;
  totalReserves?: string;
}

export interface PositionSnapshot {
  time: Date;
  userAddress: string;
  protocol: string;
  positionId?: string;
  collateralAsset: string;
  collateralAmount: string;
  debtAsset?: string;
  debtAmount?: string;
  healthFactor?: number;
  ltv?: number;
  leverage?: number;
  liquidationPrice?: string;
  collateralValueUsd?: string;
  debtValueUsd?: string;
}

export interface CollectionConfig {
  batchSize: number;
  maxRetries: number;
  retryDelay: number;
  rateLimit: number;
}

// ============================================================
// ONCHAIN DATA COLLECTOR
// ============================================================

export class OnchainDataCollector {
  private provider: ethers.providers.Provider;
  private db: Pool;
  private config: CollectionConfig;
  private chainlinkFeeds: Map<string, string>;
  private isCollecting: boolean = false;

  constructor(
    provider: ethers.providers.Provider,
    db: Pool,
    config?: Partial<CollectionConfig>
  ) {
    this.provider = provider;
    this.db = db;
    this.config = {
      batchSize: config?.batchSize || 10000,
      maxRetries: config?.maxRetries || 3,
      retryDelay: config?.retryDelay || 1000,
      rateLimit: config?.rateLimit || 100, // requests per second
    };

    // Initialize Chainlink price feeds (Arbitrum Sepolia)
    this.chainlinkFeeds = new Map([
      // Will be populated from config
    ]);
  }

  // ============================================================
  // PRICE DATA COLLECTION
  // ============================================================

  /**
   * Collect historical price data from Chainlink
   */
  async collectPriceHistory(
    asset: string,
    fromBlock: number,
    toBlock: number
  ): Promise<PricePoint[]> {
    const pricePoints: PricePoint[] = [];
    const feedAddress = this.chainlinkFeeds.get(asset);

    if (!feedAddress) {
      console.warn(`No Chainlink feed found for ${asset}`);
      return pricePoints;
    }

    try {
      // Chainlink Aggregator ABI (minimal)
      const aggregatorABI = [
        'function latestRoundData() view returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)',
        'event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)',
      ];

      const aggregator = new ethers.Contract(feedAddress, aggregatorABI, this.provider);

      // Get AnswerUpdated events in batches
      for (let start = fromBlock; start <= toBlock; start += this.config.batchSize) {
        const end = Math.min(start + this.config.batchSize - 1, toBlock);

        const events = await this.retryOperation(async () => {
          return await aggregator.queryFilter('AnswerUpdated', start, end);
        });

        for (const event of events) {
          const block = await this.retryOperation(async () => {
            return await this.provider.getBlock(event.blockNumber);
          });

          pricePoints.push({
            time: new Date(block.timestamp * 1000),
            asset,
            price: ethers.utils.formatUnits(event.args!.current, 8), // Chainlink uses 8 decimals
            source: 'chainlink',
          });
        }

        // Rate limiting
        await this.sleep(1000 / this.config.rateLimit);
      }

      // Store in database
      await this.storePricePoints(pricePoints);

      console.log(`Collected ${pricePoints.length} price points for ${asset}`);
      return pricePoints;
    } catch (error) {
      console.error(`Error collecting price history for ${asset}:`, error);
      throw error;
    }
  }

  /**
   * Get current price from Chainlink
   */
  async getCurrentPrice(asset: string): Promise<PricePoint | null> {
    const feedAddress = this.chainlinkFeeds.get(asset);
    if (!feedAddress) return null;

    try {
      const aggregatorABI = [
        'function latestRoundData() view returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)',
      ];

      const aggregator = new ethers.Contract(feedAddress, aggregatorABI, this.provider);
      const roundData = await aggregator.latestRoundData();

      const pricePoint: PricePoint = {
        time: new Date(roundData.updatedAt.toNumber() * 1000),
        asset,
        price: ethers.utils.formatUnits(roundData.answer, 8),
        source: 'chainlink',
      };

      // Store immediately
      await this.storePricePoints([pricePoint]);

      return pricePoint;
    } catch (error) {
      console.error(`Error getting current price for ${asset}:`, error);
      return null;
    }
  }

  /**
   * Collect price data from CoinGecko (backup/historical)
   */
  async collectCoinGeckoHistory(
    coinId: string,
    asset: string,
    days: number = 30
  ): Promise<PricePoint[]> {
    try {
      const response = await axios.get(
        `https://api.coingecko.com/api/v3/coins/${coinId}/market_chart`,
        {
          params: {
            vs_currency: 'usd',
            days: days,
            interval: 'hourly',
          },
        }
      );

      const pricePoints: PricePoint[] = response.data.prices.map(
        ([timestamp, price]: [number, number]) => ({
          time: new Date(timestamp),
          asset,
          price: price.toString(),
          source: 'coingecko' as const,
          volume: undefined,
          marketCap: undefined,
        })
      );

      await this.storePricePoints(pricePoints);

      console.log(`Collected ${pricePoints.length} CoinGecko price points for ${asset}`);
      return pricePoints;
    } catch (error) {
      console.error(`Error collecting CoinGecko data for ${asset}:`, error);
      return [];
    }
  }

  // ============================================================
  // LIQUIDATION DATA COLLECTION
  // ============================================================

  /**
   * Collect Aave V3 liquidations
   */
  async collectAaveLiquidations(
    fromBlock: number,
    toBlock: number
  ): Promise<LiquidationEvent[]> {
    const liquidations: LiquidationEvent[] = [];

    try {
      // Aave V3 Pool ABI (minimal)
      const poolABI = [
        'event LiquidationCall(address indexed collateralAsset, address indexed debtAsset, address indexed user, uint256 debtToCover, uint256 liquidatedCollateralAmount, address liquidator, bool receiveAToken)',
      ];

      // Get Aave pool address from config
      const poolAddress = process.env.AAVE_POOL_ADDRESS;
      if (!poolAddress) {
        console.warn('Aave pool address not configured');
        return liquidations;
      }

      const pool = new ethers.Contract(poolAddress, poolABI, this.provider);

      for (let start = fromBlock; start <= toBlock; start += this.config.batchSize) {
        const end = Math.min(start + this.config.batchSize - 1, toBlock);

        const events = await this.retryOperation(async () => {
          return await pool.queryFilter('LiquidationCall', start, end);
        });

        for (const event of events) {
          const block = await this.retryOperation(async () => {
            return await this.provider.getBlock(event.blockNumber);
          });

          const tx = await this.retryOperation(async () => {
            return await this.provider.getTransaction(event.transactionHash);
          });

          liquidations.push({
            time: new Date(block.timestamp * 1000),
            protocol: 'aave',
            userAddress: event.args!.user,
            collateralAsset: event.args!.collateralAsset,
            debtAsset: event.args!.debtAsset,
            collateralAmount: ethers.utils.formatEther(event.args!.liquidatedCollateralAmount),
            debtAmount: ethers.utils.formatEther(event.args!.debtToCover),
            liquidationPrice: '0', // Calculate from price feeds
            healthFactor: undefined, // Not available in event
            gasPrice: tx.gasPrice?.toString() || '0',
            txHash: event.transactionHash,
            blockNumber: event.blockNumber,
          });
        }

        await this.sleep(1000 / this.config.rateLimit);
      }

      await this.storeLiquidations(liquidations);

      console.log(`Collected ${liquidations.length} Aave liquidations`);
      return liquidations;
    } catch (error) {
      console.error('Error collecting Aave liquidations:', error);
      throw error;
    }
  }

  /**
   * Collect GMX V2 liquidations
   */
  async collectGMXLiquidations(
    fromBlock: number,
    toBlock: number
  ): Promise<LiquidationEvent[]> {
    const liquidations: LiquidationEvent[] = [];

    // Implementation similar to Aave, using GMX adapter events
    // Will be implemented based on GMX adapter

    return liquidations;
  }

  // ============================================================
  // USER BEHAVIOR COLLECTION
  // ============================================================

  /**
   * Calculate and update user behavior metrics
   */
  async collectUserBehavior(userAddress: string): Promise<UserMetrics> {
    try {
      // Query user's historical data
      const query = `
        WITH user_liquidations AS (
          SELECT COUNT(*) as count
          FROM liquidation_history
          WHERE user_address = $1
        ),
        user_positions AS (
          SELECT
            COUNT(DISTINCT position_id) as total_positions,
            AVG(leverage) as avg_leverage,
            MAX(collateral_value_usd) as max_position_size,
            AVG(health_factor) as avg_health_factor,
            MAX(time) as last_activity,
            SUM(collateral_value_usd) as total_volume
          FROM position_snapshots
          WHERE user_address = $1
            AND time > NOW() - INTERVAL '90 days'
        ),
        user_assets AS (
          SELECT
            jsonb_agg(DISTINCT collateral_asset) as assets,
            jsonb_agg(DISTINCT protocol) as protocols
          FROM position_snapshots
          WHERE user_address = $1
            AND time > NOW() - INTERVAL '30 days'
        )
        SELECT
          ul.count as liquidation_count,
          up.total_positions,
          up.avg_leverage,
          up.max_position_size,
          up.avg_health_factor,
          up.last_activity,
          up.total_volume,
          ua.assets as preferred_assets,
          ua.protocols as preferred_protocols
        FROM user_liquidations ul
        CROSS JOIN user_positions up
        CROSS JOIN user_assets ua
      `;

      const result = await this.db.query(query, [userAddress]);
      const row = result.rows[0];

      // Calculate risk score (0-100)
      const riskScore = this.calculateRiskScore({
        liquidationCount: row.liquidation_count || 0,
        avgLeverage: parseFloat(row.avg_leverage) || 1,
        avgHealthFactor: parseFloat(row.avg_health_factor) || 2,
      });

      const metrics: UserMetrics = {
        userAddress,
        riskScore,
        avgLeverage: parseFloat(row.avg_leverage) || 0,
        liquidationCount: row.liquidation_count || 0,
        totalPositions: row.total_positions || 0,
        maxPositionSize: row.max_position_size || '0',
        preferredAssets: row.preferred_assets || [],
        preferredProtocols: row.preferred_protocols || [],
        avgHealthFactor: parseFloat(row.avg_health_factor) || 0,
        totalVolumeUsd: row.total_volume || '0',
        lastActivity: row.last_activity || new Date(),
      };

      // Update user_risk_profiles table
      await this.updateUserRiskProfile(metrics);

      return metrics;
    } catch (error) {
      console.error(`Error collecting user behavior for ${userAddress}:`, error);
      throw error;
    }
  }

  /**
   * Calculate user risk score (0-100)
   */
  private calculateRiskScore(params: {
    liquidationCount: number;
    avgLeverage: number;
    avgHealthFactor: number;
  }): number {
    const { liquidationCount, avgLeverage, avgHealthFactor } = params;

    let score = 50; // Start at medium risk

    // Liquidation history penalty
    score += Math.min(liquidationCount * 10, 30);

    // Leverage penalty
    if (avgLeverage > 5) score += 15;
    else if (avgLeverage > 3) score += 10;
    else if (avgLeverage > 2) score += 5;

    // Health factor adjustment
    if (avgHealthFactor < 1.2) score += 20;
    else if (avgHealthFactor < 1.5) score += 10;
    else if (avgHealthFactor > 2.5) score -= 10;
    else if (avgHealthFactor > 3.0) score -= 20;

    return Math.max(0, Math.min(100, score));
  }

  // ============================================================
  // MARKET DEPTH COLLECTION
  // ============================================================

  /**
   * Collect Aave market depth
   */
  async collectAaveMarketDepth(): Promise<MarketDepth[]> {
    const markets: MarketDepth[] = [];

    try {
      // Aave V3 Pool ABI
      const poolABI = [
        'function getReserveData(address asset) view returns (tuple(uint256 configuration, uint128 liquidityIndex, uint128 currentLiquidityRate, uint128 variableBorrowIndex, uint128 currentVariableBorrowRate, uint128 currentStableBorrowRate, uint40 lastUpdateTimestamp, uint16 id, address aTokenAddress, address stableDebtTokenAddress, address variableDebtTokenAddress, address interestRateStrategyAddress, uint128 accruedToTreasury, uint128 unbacked, uint128 isolationModeTotalDebt))',
      ];

      const poolAddress = process.env.AAVE_POOL_ADDRESS;
      if (!poolAddress) return markets;

      const pool = new ethers.Contract(poolAddress, poolABI, this.provider);

      // List of assets to monitor
      const assets = [
        // Will be populated from config
      ];

      for (const asset of assets) {
        const reserveData = await pool.getReserveData(asset);

        // Get total supply and borrow from aToken
        const aTokenABI = ['function totalSupply() view returns (uint256)'];
        const aToken = new ethers.Contract(
          reserveData.aTokenAddress,
          aTokenABI,
          this.provider
        );
        const totalSupply = await aToken.totalSupply();

        const supplyAPY = parseFloat(ethers.utils.formatUnits(reserveData.currentLiquidityRate, 27)) * 100;
        const borrowAPY = parseFloat(ethers.utils.formatUnits(reserveData.currentVariableBorrowRate, 27)) * 100;

        markets.push({
          time: new Date(),
          protocol: 'aave',
          market: asset,
          totalSupply: ethers.utils.formatEther(totalSupply),
          totalBorrow: '0', // Calculate from debt tokens
          utilizationRate: 0, // Calculate
          supplyAPY,
          borrowAPY,
          availableLiquidity: '0', // Calculate
        });
      }

      await this.storeMarketDepth(markets);

      return markets;
    } catch (error) {
      console.error('Error collecting Aave market depth:', error);
      return [];
    }
  }

  // ============================================================
  // POSITION SNAPSHOTS
  // ============================================================

  /**
   * Take snapshot of all user positions
   */
  async capturePositionSnapshots(): Promise<number> {
    let count = 0;

    try {
      // This would iterate through all active positions
      // and capture their current state
      // Implementation depends on how positions are tracked

      console.log(`Captured ${count} position snapshots`);
      return count;
    } catch (error) {
      console.error('Error capturing position snapshots:', error);
      return 0;
    }
  }

  // ============================================================
  // REAL-TIME COLLECTION
  // ============================================================

  /**
   * Start real-time data collection
   */
  async startRealtimeCollection(): Promise<void> {
    if (this.isCollecting) {
      console.log('Already collecting data');
      return;
    }

    this.isCollecting = true;
    console.log('Starting real-time data collection...');

    // Set up listeners for price updates
    this.setupPriceListeners();

    // Set up listeners for liquidations
    this.setupLiquidationListeners();

    // Schedule periodic snapshots
    this.scheduleSnapshots();

    console.log('Real-time data collection started');
  }

  /**
   * Stop real-time data collection
   */
  async stopRealtimeCollection(): Promise<void> {
    this.isCollecting = false;
    // Clean up listeners
    console.log('Real-time data collection stopped');
  }

  private setupPriceListeners(): void {
    // Set up Chainlink price feed listeners
    // Implementation
  }

  private setupLiquidationListeners(): void {
    // Set up protocol event listeners
    // Implementation
  }

  private scheduleSnapshots(): void {
    // Schedule position and market snapshots
    setInterval(async () => {
      if (!this.isCollecting) return;

      await this.capturePositionSnapshots();
      await this.collectAaveMarketDepth();
    }, 60000); // Every minute
  }

  // ============================================================
  // DATABASE OPERATIONS
  // ============================================================

  private async storePricePoints(points: PricePoint[]): Promise<void> {
    if (points.length === 0) return;

    const query = `
      INSERT INTO price_history (time, asset, price, source, volume, market_cap)
      VALUES ($1, $2, $3, $4, $5, $6)
      ON CONFLICT (time, asset, source) DO UPDATE SET
        price = EXCLUDED.price,
        volume = EXCLUDED.volume,
        market_cap = EXCLUDED.market_cap
    `;

    const client = await this.db.connect();
    try {
      await client.query('BEGIN');

      for (const point of points) {
        await client.query(query, [
          point.time,
          point.asset,
          point.price,
          point.source,
          point.volume || null,
          point.marketCap || null,
        ]);
      }

      await client.query('COMMIT');
    } catch (error) {
      await client.query('ROLLBACK');
      throw error;
    } finally {
      client.release();
    }
  }

  private async storeLiquidations(liquidations: LiquidationEvent[]): Promise<void> {
    if (liquidations.length === 0) return;

    const query = `
      INSERT INTO liquidation_history (
        time, protocol, user_address, collateral_asset, debt_asset,
        collateral_amount, debt_amount, liquidation_price, health_factor,
        gas_price, tx_hash, block_number
      )
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
      ON CONFLICT (tx_hash) DO NOTHING
    `;

    const client = await this.db.connect();
    try {
      await client.query('BEGIN');

      for (const liq of liquidations) {
        await client.query(query, [
          liq.time,
          liq.protocol,
          liq.userAddress,
          liq.collateralAsset,
          liq.debtAsset || null,
          liq.collateralAmount,
          liq.debtAmount || null,
          liq.liquidationPrice,
          liq.healthFactor || null,
          liq.gasPrice,
          liq.txHash,
          liq.blockNumber,
        ]);
      }

      await client.query('COMMIT');
    } catch (error) {
      await client.query('ROLLBACK');
      throw error;
    } finally {
      client.release();
    }
  }

  private async updateUserRiskProfile(metrics: UserMetrics): Promise<void> {
    const query = `
      INSERT INTO user_risk_profiles (
        user_address, risk_score, avg_leverage, liquidation_count,
        total_positions, max_position_size, preferred_assets,
        preferred_protocols, avg_health_factor, total_volume_usd, last_activity
      )
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
      ON CONFLICT (user_address) DO UPDATE SET
        risk_score = EXCLUDED.risk_score,
        avg_leverage = EXCLUDED.avg_leverage,
        liquidation_count = EXCLUDED.liquidation_count,
        total_positions = EXCLUDED.total_positions,
        max_position_size = EXCLUDED.max_position_size,
        preferred_assets = EXCLUDED.preferred_assets,
        preferred_protocols = EXCLUDED.preferred_protocols,
        avg_health_factor = EXCLUDED.avg_health_factor,
        total_volume_usd = EXCLUDED.total_volume_usd,
        last_activity = EXCLUDED.last_activity,
        last_updated = NOW()
    `;

    await this.db.query(query, [
      metrics.userAddress,
      metrics.riskScore,
      metrics.avgLeverage,
      metrics.liquidationCount,
      metrics.totalPositions,
      metrics.maxPositionSize,
      JSON.stringify(metrics.preferredAssets),
      JSON.stringify(metrics.preferredProtocols),
      metrics.avgHealthFactor,
      metrics.totalVolumeUsd,
      metrics.lastActivity,
    ]);
  }

  private async storeMarketDepth(markets: MarketDepth[]): Promise<void> {
    if (markets.length === 0) return;

    const query = `
      INSERT INTO market_depth_snapshots (
        time, protocol, market, total_supply, total_borrow,
        utilization_rate, supply_apy, borrow_apy, available_liquidity, total_reserves
      )
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `;

    const client = await this.db.connect();
    try {
      await client.query('BEGIN');

      for (const market of markets) {
        await client.query(query, [
          market.time,
          market.protocol,
          market.market,
          market.totalSupply,
          market.totalBorrow,
          market.utilizationRate,
          market.supplyAPY,
          market.borrowAPY,
          market.availableLiquidity,
          market.totalReserves || null,
        ]);
      }

      await client.query('COMMIT');
    } catch (error) {
      await client.query('ROLLBACK');
      throw error;
    } finally {
      client.release();
    }
  }

  // ============================================================
  // UTILITY FUNCTIONS
  // ============================================================

  private async retryOperation<T>(
    operation: () => Promise<T>,
    retries: number = this.config.maxRetries
  ): Promise<T> {
    for (let i = 0; i < retries; i++) {
      try {
        return await operation();
      } catch (error) {
        if (i === retries - 1) throw error;
        await this.sleep(this.config.retryDelay * Math.pow(2, i));
      }
    }
    throw new Error('Max retries exceeded');
  }

  private sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}
