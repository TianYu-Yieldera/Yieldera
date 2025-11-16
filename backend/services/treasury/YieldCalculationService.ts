/**
 * Treasury Yield Calculation Service
 *
 * Calculates and distributes daily yields for tokenized US Treasury bonds
 * Supports T-Bills, T-Notes, and T-Bonds with automated compounding
 *
 * Key Features:
 * - Real-time yield calculation based on current rates
 * - Automated daily accrual
 * - Compound interest support
 * - Historical yield tracking
 * - Tax reporting data generation
 */

import { ethers } from 'ethers';
import { Pool } from 'pg';
import axios from 'axios';

// Treasury bond types
export enum BondType {
  TBILL_3M = 'TBILL_3M',   // 3-month Treasury Bill
  TBILL_6M = 'TBILL_6M',   // 6-month Treasury Bill
  TNOTE_2Y = 'TNOTE_2Y',   // 2-year Treasury Note
  TNOTE_5Y = 'TNOTE_5Y',   // 5-year Treasury Note
  TNOTE_10Y = 'TNOTE_10Y', // 10-year Treasury Note
  TBOND_20Y = 'TBOND_20Y', // 20-year Treasury Bond
  TBOND_30Y = 'TBOND_30Y', // 30-year Treasury Bond
}

// Yield rates (updated from external APIs)
export interface YieldRate {
  bondType: BondType;
  annualYield: number;    // e.g., 0.045 for 4.5%
  effectiveDate: Date;
  source: string;
}

// User holding information
export interface UserHolding {
  userId: string;
  bondType: BondType;
  tokenAmount: bigint;    // Amount of tokenized bonds
  principalUsd: number;   // Original investment in USD
  purchaseDate: Date;
  lastYieldDate: Date;
  totalYieldEarned: number;
  compoundingEnabled: boolean;
}

// Daily yield calculation result
export interface YieldAccrual {
  userId: string;
  bondType: BondType;
  date: Date;
  principalAmount: number;
  dailyYieldRate: number;
  yieldAmountUsd: number;
  cumulativeYield: number;
  compounded: boolean;
}

export class YieldCalculationService {
  private provider: ethers.Provider;
  private db: Pool;
  private treasuryApiUrl: string;
  private yieldRates: Map<BondType, YieldRate> = new Map();
  private lastRateUpdate: Date = new Date(0);

  // Contract addresses (will be set from env)
  private treasuryTokenContracts: Map<BondType, ethers.Contract> = new Map();

  constructor(
    rpcUrl: string,
    dbPool: Pool,
    treasuryApiUrl = 'https://api.fiscaldata.treasury.gov/services/api/fiscal_service'
  ) {
    this.provider = new ethers.JsonRpcProvider(rpcUrl);
    this.db = dbPool;
    this.treasuryApiUrl = treasuryApiUrl;

    this.initializeTreasuryContracts();
  }

  /**
   * Initialize treasury token contracts
   */
  private initializeTreasuryContracts(): void {
    const abi = [
      'function balanceOf(address owner) view returns (uint256)',
      'function mint(address to, uint256 amount)',
      'function totalSupply() view returns (uint256)',
    ];

    // Initialize contracts for each bond type
    Object.values(BondType).forEach(bondType => {
      const address = process.env[`${bondType}_TOKEN_ADDRESS`] || ethers.ZeroAddress;
      this.treasuryTokenContracts.set(
        bondType,
        new ethers.Contract(address, abi, this.provider)
      );
    });
  }

  /**
   * Fetch latest yield rates from US Treasury API
   */
  async updateYieldRates(): Promise<void> {
    try {
      console.log('Fetching latest Treasury yield rates...');

      // US Treasury provides daily yield curve rates
      // API: https://fiscaldata.treasury.gov/datasets/daily-treasury-rates/
      const response = await axios.get(
        `${this.treasuryApiUrl}/v2/accounting/od/avg_interest_rates`,
        {
          params: {
            fields: 'record_date,security_desc,avg_interest_rate_amt',
            filter: `record_date:eq:${this.getLatestBusinessDay()}`,
            sort: '-record_date',
          },
        }
      );

      const data = response.data.data;

      // Map treasury descriptions to our BondType enum
      const mappings = [
        { desc: 'Treasury Bills', bondType: BondType.TBILL_3M, rate: 0.045 },
        { desc: 'Treasury Bills', bondType: BondType.TBILL_6M, rate: 0.047 },
        { desc: 'Treasury Notes', bondType: BondType.TNOTE_2Y, rate: 0.048 },
        { desc: 'Treasury Notes', bondType: BondType.TNOTE_5Y, rate: 0.050 },
        { desc: 'Treasury Notes', bondType: BondType.TNOTE_10Y, rate: 0.052 },
        { desc: 'Treasury Bonds', bondType: BondType.TBOND_20Y, rate: 0.053 },
        { desc: 'Treasury Bonds', bondType: BondType.TBOND_30Y, rate: 0.055 },
      ];

      // Update rates (using fallback if API fails)
      for (const mapping of mappings) {
        const apiData = data.find((d: any) =>
          d.security_desc?.includes(mapping.desc)
        );

        const yieldRate: YieldRate = {
          bondType: mapping.bondType,
          annualYield: apiData?.avg_interest_rate_amt
            ? parseFloat(apiData.avg_interest_rate_amt) / 100
            : mapping.rate,
          effectiveDate: new Date(),
          source: apiData ? 'US Treasury API' : 'Fallback',
        };

        this.yieldRates.set(mapping.bondType, yieldRate);
      }

      this.lastRateUpdate = new Date();
      await this.saveYieldRatesToDB();

      console.log('✓ Yield rates updated successfully');
    } catch (error) {
      console.error('Error updating yield rates:', error);
      // Use fallback rates if API fails
      this.loadFallbackRates();
    }
  }

  /**
   * Load fallback rates if API is unavailable
   */
  private loadFallbackRates(): void {
    const fallbackRates = [
      { bondType: BondType.TBILL_3M, rate: 0.045 },
      { bondType: BondType.TBILL_6M, rate: 0.047 },
      { bondType: BondType.TNOTE_2Y, rate: 0.048 },
      { bondType: BondType.TNOTE_5Y, rate: 0.050 },
      { bondType: BondType.TNOTE_10Y, rate: 0.052 },
      { bondType: BondType.TBOND_20Y, rate: 0.053 },
      { bondType: BondType.TBOND_30Y, rate: 0.055 },
    ];

    fallbackRates.forEach(({ bondType, rate }) => {
      this.yieldRates.set(bondType, {
        bondType,
        annualYield: rate,
        effectiveDate: new Date(),
        source: 'Fallback',
      });
    });
  }

  /**
   * Calculate daily yield for all users
   * This should be run once per day via cron job
   */
  async calculateDailyYields(): Promise<YieldAccrual[]> {
    console.log('Starting daily yield calculation...');

    // Ensure we have latest rates
    if (this.shouldUpdateRates()) {
      await this.updateYieldRates();
    }

    // Get all active holdings
    const holdings = await this.getAllActiveHoldings();
    const accruals: YieldAccrual[] = [];

    for (const holding of holdings) {
      const accrual = await this.calculateYieldForHolding(holding);
      accruals.push(accrual);
    }

    // Save all accruals to database
    await this.saveYieldAccruals(accruals);

    console.log(`✓ Calculated yields for ${accruals.length} holdings`);
    return accruals;
  }

  /**
   * Calculate yield for a single user holding
   */
  private async calculateYieldForHolding(holding: UserHolding): Promise<YieldAccrual> {
    const yieldRate = this.yieldRates.get(holding.bondType);

    if (!yieldRate) {
      throw new Error(`No yield rate found for ${holding.bondType}`);
    }

    // Calculate daily yield
    const dailyRate = yieldRate.annualYield / 365;

    // Principal is either original amount or compounded amount
    const principal = holding.compoundingEnabled
      ? holding.principalUsd + holding.totalYieldEarned
      : holding.principalUsd;

    const dailyYield = principal * dailyRate;

    const accrual: YieldAccrual = {
      userId: holding.userId,
      bondType: holding.bondType,
      date: new Date(),
      principalAmount: principal,
      dailyYieldRate: dailyRate,
      yieldAmountUsd: dailyYield,
      cumulativeYield: holding.totalYieldEarned + dailyYield,
      compounded: holding.compoundingEnabled,
    };

    return accrual;
  }

  /**
   * Calculate projected yield for a given investment
   * Used for UI display before investment
   */
  async calculateProjectedYield(
    bondType: BondType,
    principalUsd: number,
    durationDays: number,
    compounding: boolean = true
  ): Promise<{
    totalYield: number;
    effectiveAPY: number;
    dailyYield: number;
    projectedValue: number;
  }> {
    const yieldRate = this.yieldRates.get(bondType);

    if (!yieldRate) {
      throw new Error(`No yield rate found for ${bondType}`);
    }

    let totalYield: number;
    let finalValue: number;

    if (compounding) {
      // Compound interest formula: A = P(1 + r/n)^(nt)
      // Daily compounding: n = 365
      const rate = yieldRate.annualYield;
      const years = durationDays / 365;
      finalValue = principalUsd * Math.pow(1 + rate / 365, 365 * years);
      totalYield = finalValue - principalUsd;
    } else {
      // Simple interest: I = P * r * t
      const rate = yieldRate.annualYield;
      const years = durationDays / 365;
      totalYield = principalUsd * rate * years;
      finalValue = principalUsd + totalYield;
    }

    const dailyYield = (yieldRate.annualYield / 365) * principalUsd;
    const effectiveAPY = ((finalValue / principalUsd) ** (365 / durationDays) - 1);

    return {
      totalYield,
      effectiveAPY,
      dailyYield,
      projectedValue: finalValue,
    };
  }

  /**
   * Get user's yield history
   */
  async getUserYieldHistory(
    userId: string,
    bondType?: BondType,
    startDate?: Date,
    endDate?: Date
  ): Promise<YieldAccrual[]> {
    let query = `
      SELECT * FROM treasury_yield_accruals
      WHERE user_id = $1
    `;
    const params: any[] = [userId];
    let paramIndex = 2;

    if (bondType) {
      query += ` AND bond_type = $${paramIndex}`;
      params.push(bondType);
      paramIndex++;
    }

    if (startDate) {
      query += ` AND date >= $${paramIndex}`;
      params.push(startDate);
      paramIndex++;
    }

    if (endDate) {
      query += ` AND date <= $${paramIndex}`;
      params.push(endDate);
      paramIndex++;
    }

    query += ' ORDER BY date DESC';

    const result = await this.db.query(query, params);
    return result.rows.map(row => ({
      userId: row.user_id,
      bondType: row.bond_type,
      date: new Date(row.date),
      principalAmount: parseFloat(row.principal_amount),
      dailyYieldRate: parseFloat(row.daily_yield_rate),
      yieldAmountUsd: parseFloat(row.yield_amount_usd),
      cumulativeYield: parseFloat(row.cumulative_yield),
      compounded: row.compounded,
    }));
  }

  /**
   * Get total yield earned by user
   */
  async getUserTotalYield(userId: string): Promise<{
    totalYield: number;
    yieldByBondType: Map<BondType, number>;
    lastCalculationDate: Date;
  }> {
    const result = await this.db.query(
      `SELECT
        bond_type,
        SUM(yield_amount_usd) as total_yield,
        MAX(date) as last_date
       FROM treasury_yield_accruals
       WHERE user_id = $1
       GROUP BY bond_type`,
      [userId]
    );

    const yieldByBondType = new Map<BondType, number>();
    let totalYield = 0;
    let lastDate = new Date(0);

    for (const row of result.rows) {
      const bondYield = parseFloat(row.total_yield);
      yieldByBondType.set(row.bond_type as BondType, bondYield);
      totalYield += bondYield;

      const rowDate = new Date(row.last_date);
      if (rowDate > lastDate) {
        lastDate = rowDate;
      }
    }

    return {
      totalYield,
      yieldByBondType,
      lastCalculationDate: lastDate,
    };
  }

  /**
   * Enable/disable compounding for user
   */
  async setCompounding(
    userId: string,
    bondType: BondType,
    enabled: boolean
  ): Promise<void> {
    await this.db.query(
      `UPDATE treasury_holdings
       SET compounding_enabled = $1, updated_at = NOW()
       WHERE user_id = $2 AND bond_type = $3`,
      [enabled, userId, bondType]
    );

    console.log(`✓ Compounding ${enabled ? 'enabled' : 'disabled'} for user ${userId}`);
  }

  /**
   * Generate tax reporting data (1099-INT equivalent)
   */
  async generateTaxReport(
    userId: string,
    year: number
  ): Promise<{
    totalInterest: number;
    byBondType: Map<BondType, number>;
    reportGenerated: Date;
  }> {
    const startDate = new Date(year, 0, 1);
    const endDate = new Date(year, 11, 31);

    const result = await this.db.query(
      `SELECT
        bond_type,
        SUM(yield_amount_usd) as total_yield
       FROM treasury_yield_accruals
       WHERE user_id = $1 AND date >= $2 AND date <= $3
       GROUP BY bond_type`,
      [userId, startDate, endDate]
    );

    const byBondType = new Map<BondType, number>();
    let totalInterest = 0;

    for (const row of result.rows) {
      const yieldAmount = parseFloat(row.total_yield);
      byBondType.set(row.bond_type as BondType, yieldAmount);
      totalInterest += yieldAmount;
    }

    // Save tax report generation record
    await this.db.query(
      `INSERT INTO treasury_tax_reports
       (user_id, tax_year, total_interest, report_data, generated_at)
       VALUES ($1, $2, $3, $4, NOW())`,
      [userId, year, totalInterest, JSON.stringify(Object.fromEntries(byBondType))]
    );

    return {
      totalInterest,
      byBondType,
      reportGenerated: new Date(),
    };
  }

  /**
   * Private helper methods
   */

  private shouldUpdateRates(): boolean {
    const hoursSinceUpdate = (Date.now() - this.lastRateUpdate.getTime()) / (1000 * 60 * 60);
    return hoursSinceUpdate >= 24; // Update once per day
  }

  private getLatestBusinessDay(): string {
    const today = new Date();
    const day = today.getDay();

    // If Saturday, go back to Friday
    if (day === 6) {
      today.setDate(today.getDate() - 1);
    }
    // If Sunday, go back to Friday
    else if (day === 0) {
      today.setDate(today.getDate() - 2);
    }

    return today.toISOString().split('T')[0];
  }

  private async getAllActiveHoldings(): Promise<UserHolding[]> {
    const result = await this.db.query(
      `SELECT * FROM treasury_holdings WHERE token_amount > 0`
    );

    return result.rows.map(row => ({
      userId: row.user_id,
      bondType: row.bond_type as BondType,
      tokenAmount: BigInt(row.token_amount),
      principalUsd: parseFloat(row.principal_usd),
      purchaseDate: new Date(row.purchase_date),
      lastYieldDate: new Date(row.last_yield_date),
      totalYieldEarned: parseFloat(row.total_yield_earned),
      compoundingEnabled: row.compounding_enabled,
    }));
  }

  private async saveYieldRatesToDB(): Promise<void> {
    const client = await this.db.connect();
    try {
      await client.query('BEGIN');

      for (const [bondType, rate] of this.yieldRates.entries()) {
        await client.query(
          `INSERT INTO treasury_yield_rates
           (bond_type, annual_yield, effective_date, source)
           VALUES ($1, $2, $3, $4)
           ON CONFLICT (bond_type, effective_date) DO UPDATE SET
             annual_yield = EXCLUDED.annual_yield,
             source = EXCLUDED.source`,
          [bondType, rate.annualYield, rate.effectiveDate, rate.source]
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

  private async saveYieldAccruals(accruals: YieldAccrual[]): Promise<void> {
    const client = await this.db.connect();
    try {
      await client.query('BEGIN');

      for (const accrual of accruals) {
        // Save accrual record
        await client.query(
          `INSERT INTO treasury_yield_accruals
           (user_id, bond_type, date, principal_amount, daily_yield_rate,
            yield_amount_usd, cumulative_yield, compounded)
           VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
          [
            accrual.userId,
            accrual.bondType,
            accrual.date,
            accrual.principalAmount,
            accrual.dailyYieldRate,
            accrual.yieldAmountUsd,
            accrual.cumulativeYield,
            accrual.compounded,
          ]
        );

        // Update holding record
        await client.query(
          `UPDATE treasury_holdings
           SET total_yield_earned = $1,
               last_yield_date = $2,
               updated_at = NOW()
           WHERE user_id = $3 AND bond_type = $4`,
          [
            accrual.cumulativeYield,
            accrual.date,
            accrual.userId,
            accrual.bondType,
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
}

// Export singleton instance
export const yieldCalculationService = new YieldCalculationService(
  process.env.RPC_URL || 'https://eth.llamarpc.com',
  new Pool({
    connectionString: process.env.DATABASE_URL,
  })
);
