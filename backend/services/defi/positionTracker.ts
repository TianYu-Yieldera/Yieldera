/**
 * DeFi Position Tracking Service
 * Aggregates user positions across all integrated DeFi protocols
 * Critical for delivering institutional risk management to retail users
 */

import { ethers } from 'ethers';
import { Pool } from 'pg';
import axios from 'axios';

// Protocol ABIs
import AavePoolABI from '../../abis/AavePool.json';
import CompoundCometABI from '../../abis/CompoundComet.json';
import UniswapV3PositionABI from '../../abis/UniswapV3Position.json';
import GMXReaderABI from '../../abis/GMXReader.json';

export interface Position {
  userId: string;
  protocol: 'aave' | 'compound' | 'uniswap' | 'gmx';
  positionId: string;
  collateralAsset?: string;
  collateralAmount?: bigint;
  debtAsset?: string;
  debtAmount?: bigint;
  healthFactor?: number;
  liquidationPrice?: number;
  apy?: number;
  value: number; // USD value
  risk: number; // 0-100 risk score
  lastUpdated: Date;
}

export interface PortfolioSummary {
  userId: string;
  totalValue: number;
  totalCollateral: number;
  totalDebt: number;
  averageHealthFactor: number;
  overallRisk: number;
  positions: Position[];
  recommendations: string[];
}

export class PositionTracker {
  private provider: ethers.Provider;
  private db: Pool;
  private aiServiceUrl: string;
  private priceOracle: Map<string, number> = new Map();

  // Protocol contracts
  private aavePool: ethers.Contract;
  private compoundComet: ethers.Contract;
  private uniswapPositionManager: ethers.Contract;
  private gmxReader: ethers.Contract;

  constructor(
    rpcUrl: string,
    dbPool: Pool,
    aiServiceUrl = 'http://localhost:8084'
  ) {
    this.provider = new ethers.JsonRpcProvider(rpcUrl);
    this.db = dbPool;
    this.aiServiceUrl = aiServiceUrl;

    // Initialize protocol contracts
    this.aavePool = new ethers.Contract(
      process.env.AAVE_POOL_ADDRESS || '0x87870Bca3F3fD6335C3F4ce8392D69350B4fA4E2',
      AavePoolABI,
      this.provider
    );

    this.compoundComet = new ethers.Contract(
      process.env.COMPOUND_COMET_ADDRESS || '0xc3d688B66703497DAA19211EEdff47f25384cdc3',
      CompoundCometABI,
      this.provider
    );

    this.uniswapPositionManager = new ethers.Contract(
      process.env.UNISWAP_POSITION_MANAGER || '0xC36442b4a4522E871399CD717aBDD847Ab11FE88',
      UniswapV3PositionABI,
      this.provider
    );

    this.gmxReader = new ethers.Contract(
      process.env.GMX_READER_ADDRESS || '0x...',
      GMXReaderABI,
      this.provider
    );
  }

  /**
   * Track all positions for a user across protocols
   */
  async trackUserPositions(userId: string): Promise<PortfolioSummary> {
    const positions: Position[] = [];

    // Fetch positions from each protocol in parallel
    const [aavePositions, compoundPositions, uniswapPositions, gmxPositions] =
      await Promise.all([
        this.getAavePositions(userId),
        this.getCompoundPositions(userId),
        this.getUniswapPositions(userId),
        this.getGMXPositions(userId),
      ]);

    positions.push(...aavePositions, ...compoundPositions, ...uniswapPositions, ...gmxPositions);

    // Calculate portfolio metrics
    const portfolio = this.calculatePortfolioMetrics(userId, positions);

    // Get AI risk assessment
    const aiRisk = await this.getAIRiskAssessment(portfolio);
    portfolio.overallRisk = aiRisk.riskScore;
    portfolio.recommendations = aiRisk.recommendations;

    // Store in database
    await this.savePositions(positions);
    await this.savePortfolioSnapshot(portfolio);

    return portfolio;
  }

  /**
   * Get Aave positions
   */
  private async getAavePositions(userId: string): Promise<Position[]> {
    try {
      const accountData = await this.aavePool.getUserAccountData(userId);

      if (accountData.totalCollateralBase.toString() === '0') {
        return [];
      }

      const healthFactor = Number(accountData.healthFactor) / 1e18;
      const ltv = Number(accountData.ltv) / 10000;

      const position: Position = {
        userId,
        protocol: 'aave',
        positionId: `aave-${userId}`,
        collateralAmount: accountData.totalCollateralBase,
        debtAmount: accountData.totalDebtBase,
        healthFactor,
        liquidationPrice: this.calculateLiquidationPrice(
          Number(accountData.totalCollateralBase) / 1e8,
          Number(accountData.totalDebtBase) / 1e8,
          ltv
        ),
        value: Number(accountData.totalCollateralBase) / 1e8,
        risk: this.calculateRiskScore(healthFactor),
        lastUpdated: new Date(),
      };

      return [position];
    } catch (error) {
      console.error('Error fetching Aave positions:', error);
      return [];
    }
  }

  /**
   * Get Compound positions
   */
  private async getCompoundPositions(userId: string): Promise<Position[]> {
    try {
      const borrowBalance = await this.compoundComet.borrowBalanceOf(userId);
      const collateralBalance = await this.compoundComet.balanceOf(userId);

      if (borrowBalance.toString() === '0' && collateralBalance.toString() === '0') {
        return [];
      }

      const position: Position = {
        userId,
        protocol: 'compound',
        positionId: `compound-${userId}`,
        collateralAmount: collateralBalance,
        debtAmount: borrowBalance,
        healthFactor: this.calculateHealthFactor(
          Number(collateralBalance) / 1e6,
          Number(borrowBalance) / 1e6
        ),
        value: Number(collateralBalance) / 1e6,
        risk: 30, // Will be calculated by AI
        lastUpdated: new Date(),
      };

      return [position];
    } catch (error) {
      console.error('Error fetching Compound positions:', error);
      return [];
    }
  }

  /**
   * Get Uniswap V3 positions
   */
  private async getUniswapPositions(userId: string): Promise<Position[]> {
    try {
      const balance = await this.uniswapPositionManager.balanceOf(userId);
      const positions: Position[] = [];

      for (let i = 0; i < Number(balance); i++) {
        const tokenId = await this.uniswapPositionManager.tokenOfOwnerByIndex(userId, i);
        const positionData = await this.uniswapPositionManager.positions(tokenId);

        // Calculate position value (simplified)
        const liquidity = Number(positionData.liquidity);
        const value = liquidity / 1e18 * 2; // Rough estimation

        positions.push({
          userId,
          protocol: 'uniswap',
          positionId: `uniswap-${tokenId}`,
          value,
          risk: 20, // LP positions generally lower risk
          lastUpdated: new Date(),
        });
      }

      return positions;
    } catch (error) {
      console.error('Error fetching Uniswap positions:', error);
      return [];
    }
  }

  /**
   * Get GMX positions
   */
  private async getGMXPositions(userId: string): Promise<Position[]> {
    // GMX V2 position tracking implementation
    // This would connect to GMX Reader contract
    return [];
  }

  /**
   * Calculate portfolio-level metrics
   */
  private calculatePortfolioMetrics(userId: string, positions: Position[]): PortfolioSummary {
    const totalValue = positions.reduce((sum, p) => sum + p.value, 0);
    const totalCollateral = positions.reduce((sum, p) =>
      sum + (p.collateralAmount ? Number(p.collateralAmount) / 1e18 : 0), 0);
    const totalDebt = positions.reduce((sum, p) =>
      sum + (p.debtAmount ? Number(p.debtAmount) / 1e18 : 0), 0);

    const healthFactors = positions
      .filter(p => p.healthFactor !== undefined)
      .map(p => p.healthFactor!);

    const averageHealthFactor = healthFactors.length > 0
      ? healthFactors.reduce((sum, hf) => sum + hf, 0) / healthFactors.length
      : 999;

    return {
      userId,
      totalValue,
      totalCollateral,
      totalDebt,
      averageHealthFactor,
      overallRisk: 0, // Will be set by AI
      positions,
      recommendations: [],
    };
  }

  /**
   * Get AI risk assessment for portfolio
   */
  private async getAIRiskAssessment(portfolio: PortfolioSummary): Promise<{
    riskScore: number;
    recommendations: string[];
  }> {
    try {
      const response = await axios.post(`${this.aiServiceUrl}/risk/portfolio`, {
        userId: portfolio.userId,
        totalValue: portfolio.totalValue,
        totalCollateral: portfolio.totalCollateral,
        totalDebt: portfolio.totalDebt,
        healthFactor: portfolio.averageHealthFactor,
        positions: portfolio.positions.map(p => ({
          protocol: p.protocol,
          value: p.value,
          healthFactor: p.healthFactor,
        })),
      });

      return {
        riskScore: response.data.overall_risk_score,
        recommendations: response.data.recommendations || [],
      };
    } catch (error) {
      console.error('AI risk assessment failed:', error);
      return {
        riskScore: this.calculateBasicRiskScore(portfolio),
        recommendations: this.generateBasicRecommendations(portfolio),
      };
    }
  }

  /**
   * Calculate basic risk score without AI
   */
  private calculateBasicRiskScore(portfolio: PortfolioSummary): number {
    const healthFactorRisk = portfolio.averageHealthFactor < 2
      ? (2 - portfolio.averageHealthFactor) * 50
      : 0;

    const leverageRisk = portfolio.totalDebt > 0
      ? (portfolio.totalDebt / portfolio.totalCollateral) * 30
      : 0;

    return Math.min(100, healthFactorRisk + leverageRisk);
  }

  /**
   * Generate basic recommendations
   */
  private generateBasicRecommendations(portfolio: PortfolioSummary): string[] {
    const recommendations: string[] = [];

    if (portfolio.averageHealthFactor < 1.5) {
      recommendations.push('⚠️ Critical: Add collateral or reduce debt immediately');
    } else if (portfolio.averageHealthFactor < 2) {
      recommendations.push('Warning: Monitor position closely, consider reducing leverage');
    }

    if (portfolio.totalDebt / portfolio.totalCollateral > 0.7) {
      recommendations.push('High leverage detected. Consider diversifying');
    }

    const aavePositions = portfolio.positions.filter(p => p.protocol === 'aave');
    if (aavePositions.length > 0 && aavePositions[0].healthFactor! < 1.5) {
      recommendations.push('Aave position at risk. Consider adding collateral');
    }

    return recommendations;
  }

  /**
   * Save positions to database
   */
  private async savePositions(positions: Position[]): Promise<void> {
    const client = await this.db.connect();
    try {
      await client.query('BEGIN');

      for (const position of positions) {
        await client.query(
          `INSERT INTO defi_positions
           (user_id, protocol, position_id, collateral_amount, debt_amount,
            health_factor, liquidation_price, value_usd, risk_score, last_updated)
           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
           ON CONFLICT (position_id) DO UPDATE SET
             collateral_amount = EXCLUDED.collateral_amount,
             debt_amount = EXCLUDED.debt_amount,
             health_factor = EXCLUDED.health_factor,
             liquidation_price = EXCLUDED.liquidation_price,
             value_usd = EXCLUDED.value_usd,
             risk_score = EXCLUDED.risk_score,
             last_updated = EXCLUDED.last_updated`,
          [
            position.userId,
            position.protocol,
            position.positionId,
            position.collateralAmount?.toString() || '0',
            position.debtAmount?.toString() || '0',
            position.healthFactor || 999,
            position.liquidationPrice || 0,
            position.value,
            position.risk,
            position.lastUpdated,
          ]
        );
      }

      await client.query('COMMIT');
    } catch (error) {
      await client.query('ROLLBACK');
      throw error;
    } finally {
      client.release();
    }
  }

  /**
   * Save portfolio snapshot
   */
  private async savePortfolioSnapshot(portfolio: PortfolioSummary): Promise<void> {
    await this.db.query(
      `INSERT INTO portfolio_snapshots
       (user_id, total_value, total_collateral, total_debt,
        average_health_factor, overall_risk, recommendations, timestamp)
       VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`,
      [
        portfolio.userId,
        portfolio.totalValue,
        portfolio.totalCollateral,
        portfolio.totalDebt,
        portfolio.averageHealthFactor,
        portfolio.overallRisk,
        JSON.stringify(portfolio.recommendations),
      ]
    );
  }

  /**
   * Helper functions
   */
  private calculateHealthFactor(collateral: number, debt: number): number {
    if (debt === 0) return 999;
    return (collateral * 0.8) / debt; // Assuming 80% LTV
  }

  private calculateLiquidationPrice(collateral: number, debt: number, ltv: number): number {
    if (debt === 0) return 0;
    return debt / (collateral * ltv);
  }

  private calculateRiskScore(healthFactor: number): number {
    if (healthFactor > 2) return 10;
    if (healthFactor > 1.5) return 30;
    if (healthFactor > 1.2) return 60;
    if (healthFactor > 1.1) return 80;
    return 95;
  }
}

// Export singleton instance
export const positionTracker = new PositionTracker(
  process.env.RPC_URL || 'https://eth.llamarpc.com',
  new Pool({
    connectionString: process.env.DATABASE_URL,
  }),
  process.env.AI_SERVICE_URL || 'http://localhost:8084'
);