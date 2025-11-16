/**
 * Analytics Query Service
 *
 * Fast access to aggregated data for AI/ML models
 * Uses TimescaleDB continuous aggregates and optimized queries
 *
 * Phase 1: Data Layer
 * Created: 2025-11-09
 */

import { Pool, QueryResult } from 'pg';
import { LiquidationEvent } from './OnchainDataCollector';

// ============================================================
// TYPES
// ============================================================

export interface PriceStats {
  asset: string;
  period: string; // '24h', '7d', '30d'
  open: number;
  high: number;
  low: number;
  close: number;
  avgPrice: number;
  volatility: number;
  priceChange: number;
  priceChangePercent: number;
}

export interface MarketMetrics {
  protocol: string;
  market: string;
  utilizationRate: number;
  supplyAPY: number;
  borrowAPY: number;
  availableLiquidity: string;
  totalSupply: string;
  totalBorrow: string;
  trend: 'increasing' | 'decreasing' | 'stable';
}

export interface LiquidationStats {
  protocol: string;
  asset?: string;
  period: string;
  liquidationCount: number;
  totalCollateralLiquidated: string;
  totalDebtLiquidated: string;
  avgHealthFactor: number;
  medianHealthFactor: number;
  minHealthFactor: number;
  maxHealthFactor: number;
}

export interface UserRiskData {
  userAddress: string;
  currentRiskScore: number;
  riskTrend: 'improving' | 'worsening' | 'stable';
  riskHistory: Array<{ time: Date; score: number }>;
  liquidationHistory: LiquidationEvent[];
  avgHealthFactor: number;
  currentPositions: number;
  riskRanking: number;
  riskPercentile: number;
}

export interface CorrelationMatrix {
  assets: string[];
  correlations: number[][];
  period: string;
}

export interface VolatilityMetrics {
  asset: string;
  volatility1d: number;
  volatility7d: number;
  volatility30d: number;
  volatility90d: number;
  trend: 'increasing' | 'decreasing' | 'stable';
}

// ============================================================
// ANALYTICS QUERY SERVICE
// ============================================================

export class AnalyticsQuery {
  private db: Pool;

  constructor(db: Pool) {
    this.db = db;
  }

  // ============================================================
  // PRICE ANALYTICS
  // ============================================================

  /**
   * Get price statistics for an asset
   */
  async getPriceStats(asset: string, period: '24h' | '7d' | '30d' = '30d'): Promise<PriceStats> {
    const interval = period === '24h' ? '1 day' : period === '7d' ? '7 days' : '30 days';

    const query = `
      WITH price_data AS (
        SELECT
          FIRST(price, time) OVER (ORDER BY time) AS first_price,
          LAST(price, time) OVER (ORDER BY time) AS last_price,
          MAX(price) OVER () AS high,
          MIN(price) OVER () AS low,
          AVG(price) OVER () AS avg_price,
          STDDEV(price) OVER () AS stddev_price,
          time
        FROM price_history
        WHERE asset = $1
          AND source = 'chainlink'
          AND time > NOW() - INTERVAL '${interval}'
      )
      SELECT
        MAX(first_price) AS open,
        MAX(high) AS high,
        MAX(low) AS low,
        MAX(last_price) AS close,
        MAX(avg_price) AS avg_price,
        MAX(stddev_price) / NULLIF(MAX(avg_price), 0) AS volatility
      FROM price_data
    `;

    const result = await this.db.query(query, [asset]);
    const row = result.rows[0];

    if (!row || !row.close) {
      throw new Error(`No price data found for ${asset}`);
    }

    const open = parseFloat(row.open);
    const close = parseFloat(row.close);
    const priceChange = close - open;
    const priceChangePercent = (priceChange / open) * 100;

    return {
      asset,
      period,
      open,
      high: parseFloat(row.high),
      low: parseFloat(row.low),
      close,
      avgPrice: parseFloat(row.avg_price),
      volatility: parseFloat(row.volatility) || 0,
      priceChange,
      priceChangePercent,
    };
  }

  /**
   * Calculate price volatility
   */
  async getPriceVolatility(asset: string, days: number = 30): Promise<number> {
    const query = `
      SELECT STDDEV(price) / NULLIF(AVG(price), 0) AS volatility
      FROM price_history
      WHERE asset = $1
        AND source = 'chainlink'
        AND time > NOW() - INTERVAL '${days} days'
    `;

    const result = await this.db.query(query, [asset]);
    return parseFloat(result.rows[0]?.volatility) || 0;
  }

  /**
   * Get multi-period volatility metrics
   */
  async getVolatilityMetrics(asset: string): Promise<VolatilityMetrics> {
    const [vol1d, vol7d, vol30d, vol90d] = await Promise.all([
      this.getPriceVolatility(asset, 1),
      this.getPriceVolatility(asset, 7),
      this.getPriceVolatility(asset, 30),
      this.getPriceVolatility(asset, 90),
    ]);

    // Determine trend
    let trend: 'increasing' | 'decreasing' | 'stable';
    if (vol7d > vol30d * 1.2) {
      trend = 'increasing';
    } else if (vol7d < vol30d * 0.8) {
      trend = 'decreasing';
    } else {
      trend = 'stable';
    }

    return {
      asset,
      volatility1d: vol1d,
      volatility7d: vol7d,
      volatility30d: vol30d,
      volatility90d: vol90d,
      trend,
    };
  }

  /**
   * Calculate correlation between two assets
   */
  async getAssetCorrelation(asset1: string, asset2: string, days: number = 30): Promise<number> {
    const query = `
      WITH prices AS (
        SELECT
          time_bucket('1 hour', time) AS hour,
          asset,
          AVG(price) AS avg_price
        FROM price_history
        WHERE asset IN ($1, $2)
          AND source = 'chainlink'
          AND time > NOW() - INTERVAL '${days} days'
        GROUP BY hour, asset
      ),
      pivoted AS (
        SELECT
          hour,
          MAX(CASE WHEN asset = $1 THEN avg_price END) AS price1,
          MAX(CASE WHEN asset = $2 THEN avg_price END) AS price2
        FROM prices
        GROUP BY hour
      )
      SELECT CORR(price1, price2) AS correlation
      FROM pivoted
      WHERE price1 IS NOT NULL AND price2 IS NOT NULL
    `;

    const result = await this.db.query(query, [asset1, asset2]);
    return parseFloat(result.rows[0]?.correlation) || 0;
  }

  /**
   * Get correlation matrix for multiple assets
   */
  async getCorrelationMatrix(assets: string[], days: number = 30): Promise<CorrelationMatrix> {
    const n = assets.length;
    const correlations: number[][] = Array(n)
      .fill(0)
      .map(() => Array(n).fill(0));

    // Diagonal is always 1
    for (let i = 0; i < n; i++) {
      correlations[i][i] = 1;
    }

    // Calculate correlations for each pair
    for (let i = 0; i < n; i++) {
      for (let j = i + 1; j < n; j++) {
        const corr = await this.getAssetCorrelation(assets[i], assets[j], days);
        correlations[i][j] = corr;
        correlations[j][i] = corr; // Symmetric
      }
    }

    return {
      assets,
      correlations,
      period: `${days}d`,
    };
  }

  // ============================================================
  // MARKET ANALYTICS
  // ============================================================

  /**
   * Get current market metrics
   */
  async getMarketMetrics(protocol: string, market: string): Promise<MarketMetrics> {
    const query = `
      WITH current AS (
        SELECT *
        FROM market_depth_snapshots
        WHERE protocol = $1 AND market = $2
        ORDER BY time DESC
        LIMIT 1
      ),
      previous AS (
        SELECT avg_utilization, avg_supply_apy
        FROM hourly_market_metrics
        WHERE protocol = $1 AND market = $2
          AND hour > NOW() - INTERVAL '24 hours'
        ORDER BY hour DESC
        LIMIT 1
      )
      SELECT
        c.*,
        p.avg_utilization AS prev_utilization,
        p.avg_supply_apy AS prev_supply_apy
      FROM current c
      LEFT JOIN previous p ON true
    `;

    const result = await this.db.query(query, [protocol, market]);
    const row = result.rows[0];

    if (!row) {
      throw new Error(`No market data found for ${protocol}/${market}`);
    }

    // Determine trend
    let trend: 'increasing' | 'decreasing' | 'stable';
    const currentUtil = parseFloat(row.utilization_rate);
    const prevUtil = parseFloat(row.prev_utilization);

    if (prevUtil && currentUtil > prevUtil * 1.1) {
      trend = 'increasing';
    } else if (prevUtil && currentUtil < prevUtil * 0.9) {
      trend = 'decreasing';
    } else {
      trend = 'stable';
    }

    return {
      protocol,
      market,
      utilizationRate: currentUtil,
      supplyAPY: parseFloat(row.supply_apy),
      borrowAPY: parseFloat(row.borrow_apy),
      availableLiquidity: row.available_liquidity,
      totalSupply: row.total_supply,
      totalBorrow: row.total_borrow,
      trend,
    };
  }

  /**
   * Get all markets for a protocol
   */
  async getProtocolMarkets(protocol: string): Promise<MarketMetrics[]> {
    const query = `
      SELECT DISTINCT market
      FROM market_depth_snapshots
      WHERE protocol = $1
        AND time > NOW() - INTERVAL '1 hour'
    `;

    const result = await this.db.query(query, [protocol]);
    const markets = await Promise.all(
      result.rows.map(row => this.getMarketMetrics(protocol, row.market))
    );

    return markets;
  }

  // ============================================================
  // LIQUIDATION ANALYTICS
  // ============================================================

  /**
   * Get liquidation statistics
   */
  async getLiquidationStats(
    protocol: string,
    asset?: string,
    period: '24h' | '7d' | '30d' = '30d'
  ): Promise<LiquidationStats> {
    const interval = period === '24h' ? '1 day' : period === '7d' ? '7 days' : '30 days';

    const query = `
      SELECT
        COUNT(*) AS liquidation_count,
        SUM(collateral_amount) AS total_collateral,
        SUM(debt_amount) AS total_debt,
        AVG(health_factor) AS avg_health_factor,
        PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY health_factor) AS median_health_factor,
        MIN(health_factor) AS min_health_factor,
        MAX(health_factor) AS max_health_factor
      FROM liquidation_history
      WHERE protocol = $1
        ${asset ? 'AND collateral_asset = $2' : ''}
        AND time > NOW() - INTERVAL '${interval}'
        AND health_factor IS NOT NULL
    `;

    const params = asset ? [protocol, asset] : [protocol];
    const result = await this.db.query(query, params);
    const row = result.rows[0];

    return {
      protocol,
      asset,
      period,
      liquidationCount: parseInt(row.liquidation_count) || 0,
      totalCollateralLiquidated: row.total_collateral || '0',
      totalDebtLiquidated: row.total_debt || '0',
      avgHealthFactor: parseFloat(row.avg_health_factor) || 0,
      medianHealthFactor: parseFloat(row.median_health_factor) || 0,
      minHealthFactor: parseFloat(row.min_health_factor) || 0,
      maxHealthFactor: parseFloat(row.max_health_factor) || 0,
    };
  }

  /**
   * Get liquidation rate (liquidations per active position)
   */
  async getLiquidationRate(protocol: string, days: number = 30): Promise<number> {
    const query = `
      WITH liquidations AS (
        SELECT COUNT(*) AS count
        FROM liquidation_history
        WHERE protocol = $1
          AND time > NOW() - INTERVAL '${days} days'
      ),
      positions AS (
        SELECT COUNT(DISTINCT user_address) AS count
        FROM position_snapshots
        WHERE protocol = $1
          AND time > NOW() - INTERVAL '${days} days'
      )
      SELECT
        CAST(l.count AS FLOAT) / NULLIF(p.count, 0) AS rate
      FROM liquidations l
      CROSS JOIN positions p
    `;

    const result = await this.db.query(query, [protocol]);
    return parseFloat(result.rows[0]?.rate) || 0;
  }

  /**
   * Find similar historical liquidations
   */
  async getSimilarLiquidations(
    position: {
      collateralAsset: string;
      healthFactor: number;
      protocol: string;
    },
    limit: number = 100
  ): Promise<LiquidationEvent[]> {
    const query = `
      SELECT
        time,
        protocol,
        user_address,
        collateral_asset,
        debt_asset,
        collateral_amount,
        debt_amount,
        liquidation_price,
        health_factor,
        gas_price,
        tx_hash,
        block_number
      FROM liquidation_history
      WHERE protocol = $1
        AND collateral_asset = $2
        AND health_factor BETWEEN $3 - 0.2 AND $3 + 0.2
      ORDER BY time DESC
      LIMIT $4
    `;

    const result = await this.db.query(query, [
      position.protocol,
      position.collateralAsset,
      position.healthFactor,
      limit,
    ]);

    return result.rows.map(row => ({
      time: row.time,
      protocol: row.protocol,
      userAddress: row.user_address,
      collateralAsset: row.collateral_asset,
      debtAsset: row.debt_asset,
      collateralAmount: row.collateral_amount,
      debtAmount: row.debt_amount,
      liquidationPrice: row.liquidation_price,
      healthFactor: row.health_factor,
      gasPrice: row.gas_price,
      txHash: row.tx_hash,
      blockNumber: row.block_number,
    }));
  }

  // ============================================================
  // USER RISK ANALYTICS
  // ============================================================

  /**
   * Get comprehensive user risk data
   */
  async getUserRiskData(userAddress: string): Promise<UserRiskData> {
    // Get current profile
    const profileQuery = `
      SELECT
        risk_score,
        avg_health_factor,
        total_positions,
        liquidation_count
      FROM user_risk_profiles
      WHERE user_address = $1
    `;

    const profileResult = await this.db.query(profileQuery, [userAddress]);
    const profile = profileResult.rows[0];

    if (!profile) {
      throw new Error(`No risk profile found for ${userAddress}`);
    }

    // Get risk history
    const historyQuery = `
      SELECT day, avg_health_factor AS score
      FROM user_activity_summary
      WHERE user_address = $1
        AND day > NOW() - INTERVAL '30 days'
      ORDER BY day DESC
    `;

    const historyResult = await this.db.query(historyQuery, [userAddress]);
    const riskHistory = historyResult.rows.map(row => ({
      time: row.day,
      score: parseFloat(row.score) || 0,
    }));

    // Calculate trend
    let riskTrend: 'improving' | 'worsening' | 'stable' = 'stable';
    if (riskHistory.length >= 2) {
      const recent = riskHistory[0].score;
      const older = riskHistory[riskHistory.length - 1].score;
      if (recent < older * 0.9) {
        riskTrend = 'improving';
      } else if (recent > older * 1.1) {
        riskTrend = 'worsening';
      }
    }

    // Get liquidation history
    const liqQuery = `
      SELECT *
      FROM liquidation_history
      WHERE user_address = $1
      ORDER BY time DESC
      LIMIT 10
    `;

    const liqResult = await this.db.query(liqQuery, [userAddress]);
    const liquidationHistory: LiquidationEvent[] = liqResult.rows.map(row => ({
      time: row.time,
      protocol: row.protocol,
      userAddress: row.user_address,
      collateralAsset: row.collateral_asset,
      debtAsset: row.debt_asset,
      collateralAmount: row.collateral_amount,
      debtAmount: row.debt_amount,
      liquidationPrice: row.liquidation_price,
      healthFactor: row.health_factor,
      gasPrice: row.gas_price,
      txHash: row.tx_hash,
      blockNumber: row.block_number,
    }));

    // Get risk ranking
    const rankQuery = `
      SELECT risk_rank, risk_percentile
      FROM user_risk_rankings
      WHERE user_address = $1
    `;

    const rankResult = await this.db.query(rankQuery, [userAddress]);
    const ranking = rankResult.rows[0] || { risk_rank: 0, risk_percentile: 0 };

    return {
      userAddress,
      currentRiskScore: parseFloat(profile.risk_score),
      riskTrend,
      riskHistory,
      liquidationHistory,
      avgHealthFactor: parseFloat(profile.avg_health_factor) || 0,
      currentPositions: parseInt(profile.total_positions) || 0,
      riskRanking: parseInt(ranking.risk_rank) || 0,
      riskPercentile: parseFloat(ranking.risk_percentile) || 0,
    };
  }

  /**
   * Get high-risk users
   */
  async getHighRiskUsers(limit: number = 100): Promise<UserRiskData[]> {
    const query = `
      SELECT user_address
      FROM user_risk_profiles
      WHERE risk_score > 70
      ORDER BY risk_score DESC
      LIMIT $1
    `;

    const result = await this.db.query(query, [limit]);
    const users = await Promise.all(
      result.rows.map(row => this.getUserRiskData(row.user_address))
    );

    return users;
  }

  // ============================================================
  // ADVANCED ANALYTICS
  // ============================================================

  /**
   * Calculate Value at Risk (VaR) for a portfolio
   * Simple historical VaR implementation
   */
  async calculateVaR(
    positions: Array<{ asset: string; amount: number }>,
    confidenceLevel: number = 0.95,
    days: number = 1
  ): Promise<number> {
    // Get historical returns for each asset
    const returns: number[][] = [];

    for (const position of positions) {
      const query = `
        WITH prices AS (
          SELECT
            time,
            price,
            LAG(price) OVER (ORDER BY time) AS prev_price
          FROM price_history
          WHERE asset = $1
            AND source = 'chainlink'
            AND time > NOW() - INTERVAL '90 days'
          ORDER BY time
        )
        SELECT (price - prev_price) / prev_price AS return
        FROM prices
        WHERE prev_price IS NOT NULL
      `;

      const result = await this.db.query(query, [position.asset]);
      returns.push(result.rows.map(row => parseFloat(row.return)));
    }

    // Calculate portfolio returns
    const portfolioReturns: number[] = [];
    const minLength = Math.min(...returns.map(r => r.length));

    for (let i = 0; i < minLength; i++) {
      let portfolioReturn = 0;
      for (let j = 0; j < positions.length; j++) {
        portfolioReturn += returns[j][i] * positions[j].amount;
      }
      portfolioReturns.push(portfolioReturn);
    }

    // Sort and find VaR
    portfolioReturns.sort((a, b) => a - b);
    const varIndex = Math.floor((1 - confidenceLevel) * portfolioReturns.length);
    const var95 = -portfolioReturns[varIndex]; // Negative because we want loss

    return var95;
  }

  /**
   * Get performance metrics for the data pipeline
   */
  async getSystemMetrics(): Promise<{
    totalPricePoints: number;
    totalLiquidations: number;
    totalUsers: number;
    dataFreshness: number; // seconds since last update
    queryPerformance: number; // avg query time in ms
  }> {
    const query = `
      SELECT
        (SELECT COUNT(*) FROM price_history) AS price_points,
        (SELECT COUNT(*) FROM liquidation_history) AS liquidations,
        (SELECT COUNT(*) FROM user_risk_profiles) AS users,
        EXTRACT(EPOCH FROM (NOW() - MAX(time))) AS freshness
      FROM price_history
      WHERE source = 'chainlink'
    `;

    const result = await this.db.query(query);
    const row = result.rows[0];

    return {
      totalPricePoints: parseInt(row.price_points) || 0,
      totalLiquidations: parseInt(row.liquidations) || 0,
      totalUsers: parseInt(row.users) || 0,
      dataFreshness: parseFloat(row.freshness) || 0,
      queryPerformance: 0, // Would need to track separately
    };
  }
}
