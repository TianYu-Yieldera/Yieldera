/**
 * Yields Service Client
 * Manages treasury yield calculations and history tracking
 */

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

class YieldsService {
  constructor() {
    this.baseURL = `${API_BASE}/api/v1/yields`;
  }

  /**
   * Get current treasury rates
   * @returns {Promise<Object>} Current yield rates across assets
   */
  async getRates() {
    try {
      const response = await fetch(`${this.baseURL}/rates`);

      if (!response.ok) {
        throw new Error(`Failed to fetch rates: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Yields rates fetch error:', error);
      throw error;
    }
  }

  /**
   * Get user's yield history
   * @param {string} userId - User wallet address
   * @param {number} limit - Number of records to fetch
   * @returns {Promise<Object>} User yield history
   */
  async getUserYieldHistory(userId, limit = 50) {
    try {
      const response = await fetch(`${this.baseURL}/history/${userId}?limit=${limit}`);

      if (!response.ok) {
        throw new Error(`Failed to fetch yield history: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Yield history fetch error:', error);
      throw error;
    }
  }

  /**
   * Get user's total yield earned
   * @param {string} userId - User wallet address
   * @returns {Promise<Object>} Total yield statistics
   */
  async getUserTotalYield(userId) {
    try {
      const response = await fetch(`${this.baseURL}/total/${userId}`);

      if (!response.ok) {
        throw new Error(`Failed to fetch total yield: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Total yield fetch error:', error);
      throw error;
    }
  }

  /**
   * Project future yield based on current holdings
   * @param {Object} projectionData - Projection parameters
   * @param {string} projectionData.user_address - User wallet address
   * @param {number} projectionData.time_horizon_days - Days to project
   * @param {Object} projectionData.holdings - Current holdings
   * @returns {Promise<Object>} Yield projections
   */
  async projectYield(projectionData) {
    try {
      const response = await fetch(`${this.baseURL}/project`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(projectionData),
      });

      if (!response.ok) {
        throw new Error(`Failed to project yield: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Yield projection error:', error);
      throw error;
    }
  }

  /**
   * Calculate APY for a specific investment
   * @param {number} principal - Initial investment amount
   * @param {number} currentValue - Current value
   * @param {number} daysHeld - Number of days held
   * @returns {number} Annualized percentage yield
   */
  calculateAPY(principal, currentValue, daysHeld) {
    if (daysHeld <= 0 || principal <= 0) return 0;

    const gain = currentValue - principal;
    const dailyReturn = gain / principal / daysHeld;
    const apy = Math.pow(1 + dailyReturn, 365) - 1;

    return apy * 100;
  }

  /**
   * Format yield as currency
   * @param {number} value - Yield value
   * @returns {string} Formatted currency string
   */
  formatYield(value) {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(value);
  }

  /**
   * Format APY percentage
   * @param {number} apy - APY value (e.g., 0.0435 = 4.35%)
   * @returns {string} Formatted percentage
   */
  formatAPY(apy) {
    return `${(apy * 100).toFixed(2)}%`;
  }

  /**
   * Generate realistic mock treasury rates
   * @returns {Object} Mock treasury yield rates
   */
  getMockRates() {
    return {
      timestamp: new Date().toISOString(),
      rates: [
        {
          asset: 'US Treasury 3-Month',
          symbol: 'T3M',
          apy: 0.0528,
          yield_7day: 0.0531,
          price: 99.12,
          change_24h: 0.0012,
          liquidity: 125000000
        },
        {
          asset: 'US Treasury 6-Month',
          symbol: 'T6M',
          apy: 0.0545,
          yield_7day: 0.0548,
          price: 98.54,
          change_24h: 0.0008,
          liquidity: 98000000
        },
        {
          asset: 'US Treasury 1-Year',
          symbol: 'T1Y',
          apy: 0.0562,
          yield_7day: 0.0565,
          price: 97.23,
          change_24h: -0.0003,
          liquidity: 215000000
        },
        {
          asset: 'US Treasury 2-Year',
          symbol: 'T2Y',
          apy: 0.0478,
          yield_7day: 0.0481,
          price: 95.67,
          change_24h: 0.0015,
          liquidity: 180000000
        },
        {
          asset: 'US Treasury 5-Year',
          symbol: 'T5Y',
          apy: 0.0445,
          yield_7day: 0.0448,
          price: 92.34,
          change_24h: 0.0021,
          liquidity: 320000000
        },
        {
          asset: 'US Treasury 10-Year',
          symbol: 'T10Y',
          apy: 0.0438,
          yield_7day: 0.0441,
          price: 89.12,
          change_24h: -0.0007,
          liquidity: 450000000
        }
      ]
    };
  }

  /**
   * Generate realistic mock yield history
   * @param {string} userId - User wallet address
   * @param {number} days - Number of days of history
   * @returns {Object} Mock yield history
   */
  getMockYieldHistory(userId, days = 30) {
    const history = [];
    const now = new Date();

    // Generate daily yield entries
    for (let i = days - 1; i >= 0; i--) {
      const date = new Date(now);
      date.setDate(date.getDate() - i);

      // Simulate varying daily yields with realistic patterns
      const baseYield = 12.5;
      const volatility = Math.sin(i / 7) * 3; // Weekly patterns
      const random = (Math.random() - 0.5) * 2;
      const dailyYield = baseYield + volatility + random;

      history.push({
        date: date.toISOString().split('T')[0],
        timestamp: date.toISOString(),
        daily_yield: Math.max(0, dailyYield),
        cumulative_yield: baseYield * (days - i) + (volatility * (days - i) / 2),
        apy: 0.0528 + (Math.random() * 0.01 - 0.005), // Slight APY variations
        principal_balance: 50000 + (Math.random() * 1000 - 500),
        assets: [
          {
            symbol: 'T3M',
            amount: 20000,
            yield: dailyYield * 0.4
          },
          {
            symbol: 'T6M',
            amount: 30000,
            yield: dailyYield * 0.6
          }
        ]
      });
    }

    return {
      user_address: userId,
      period_days: days,
      total_entries: history.length,
      history
    };
  }

  /**
   * Get mock total yield statistics
   * @param {string} userId - User wallet address
   * @returns {Object} Mock total yield stats
   */
  getMockTotalYield(userId) {
    return {
      user_address: userId,
      total_yield_earned: 3847.52,
      total_principal: 50000,
      average_apy: 0.0538,
      current_apy: 0.0545,
      yield_by_asset: [
        {
          symbol: 'T3M',
          principal: 20000,
          yield_earned: 1234.56,
          apy: 0.0528,
          days_held: 87
        },
        {
          symbol: 'T6M',
          principal: 30000,
          yield_earned: 2612.96,
          apy: 0.0545,
          days_held: 87
        }
      ],
      yield_by_month: [
        { month: '2025-01', yield: 892.34 },
        { month: '2025-02', yield: 1234.56 },
        { month: '2025-03', yield: 1720.62 }
      ],
      projections: {
        next_30_days: 458.32,
        next_90_days: 1374.96,
        next_365_days: 5475.00
      },
      performance_metrics: {
        best_day: { date: '2025-03-15', yield: 23.45 },
        worst_day: { date: '2025-02-08', yield: 8.12 },
        avg_daily_yield: 15.32,
        consistency_score: 0.87 // 0-1 scale
      }
    };
  }

  /**
   * Enhanced getRates with fallback to mock data
   */
  async getRatesWithFallback() {
    try {
      return await this.getRates();
    } catch (error) {
      console.warn('Using mock rates data:', error);
      return this.getMockRates();
    }
  }

  /**
   * Enhanced getUserYieldHistory with fallback
   */
  async getUserYieldHistoryWithFallback(userId, limit = 50) {
    try {
      return await this.getUserYieldHistory(userId, limit);
    } catch (error) {
      console.warn('Using mock yield history:', error);
      return this.getMockYieldHistory(userId, Math.min(limit, 90));
    }
  }

  /**
   * Enhanced getUserTotalYield with fallback
   */
  async getUserTotalYieldWithFallback(userId) {
    try {
      return await this.getUserTotalYield(userId);
    } catch (error) {
      console.warn('Using mock total yield:', error);
      return this.getMockTotalYield(userId);
    }
  }

  /**
   * Calculate yield trend (increasing/decreasing)
   * @param {Array} history - Yield history array
   * @returns {Object} Trend analysis
   */
  calculateYieldTrend(history) {
    if (!history || history.length < 2) {
      return { direction: 'stable', change: 0, percentage: 0 };
    }

    const recent = history.slice(-7); // Last 7 days
    const older = history.slice(-14, -7); // Previous 7 days

    const recentAvg = recent.reduce((sum, h) => sum + h.daily_yield, 0) / recent.length;
    const olderAvg = older.reduce((sum, h) => sum + h.daily_yield, 0) / older.length;

    const change = recentAvg - olderAvg;
    const percentage = olderAvg > 0 ? (change / olderAvg) * 100 : 0;

    return {
      direction: change > 0 ? 'increasing' : change < 0 ? 'decreasing' : 'stable',
      change: Math.abs(change),
      percentage: percentage,
      recent_avg: recentAvg,
      previous_avg: olderAvg
    };
  }

  /**
   * Get yield performance summary
   * @param {Object} yieldData - Total yield data
   * @returns {Object} Performance summary
   */
  getPerformanceSummary(yieldData) {
    if (!yieldData) return null;

    const roi = (yieldData.total_yield_earned / yieldData.total_principal) * 100;
    const dailyRate = yieldData.average_apy / 365;
    const compoundedValue = yieldData.total_principal * Math.pow(1 + yieldData.average_apy, 1);

    return {
      roi: roi,
      roi_formatted: `${roi.toFixed(2)}%`,
      daily_rate: dailyRate,
      daily_earnings: yieldData.total_principal * dailyRate,
      compounded_value: compoundedValue,
      total_return: yieldData.total_yield_earned,
      performance_rating: this.getPerformanceRating(roi)
    };
  }

  /**
   * Get performance rating based on ROI
   * @param {number} roi - Return on investment percentage
   * @returns {string} Performance rating
   */
  getPerformanceRating(roi) {
    if (roi >= 10) return 'Excellent';
    if (roi >= 7) return 'Very Good';
    if (roi >= 5) return 'Good';
    if (roi >= 3) return 'Fair';
    return 'Poor';
  }
}

// Export singleton instance
const yieldsService = new YieldsService();
export default yieldsService;
