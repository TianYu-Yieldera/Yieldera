/**
 * Hedge Service Client
 * Manages auto-hedging strategies and settings
 */

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

class HedgeService {
  constructor() {
    this.baseURL = `${API_BASE}/api/v1/hedge`;
  }

  /**
   * Get user's hedge history
   * @param {string} userId - User wallet address
   * @param {number} limit - Number of records to fetch
   * @returns {Promise<Object>} Hedge execution history
   */
  async getHedgeHistory(userId, limit = 50) {
    try {
      const response = await fetch(`${this.baseURL}/history/${userId}?limit=${limit}`);

      if (!response.ok) {
        throw new Error(`Failed to fetch hedge history: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Hedge history fetch error:', error);
      throw error;
    }
  }

  /**
   * Get hedge statistics
   * @returns {Promise<Object>} System-wide hedge statistics
   */
  async getHedgeStats() {
    try {
      const response = await fetch(`${this.baseURL}/stats`);

      if (!response.ok) {
        throw new Error(`Failed to fetch hedge stats: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Hedge stats fetch error:', error);
      throw error;
    }
  }

  /**
   * Get user's hedge settings
   * @param {string} userId - User wallet address
   * @returns {Promise<Object>} User hedge settings
   */
  async getUserSettings(userId) {
    try {
      const response = await fetch(`${this.baseURL}/settings/${userId}`);

      if (!response.ok) {
        throw new Error(`Failed to fetch hedge settings: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Hedge settings fetch error:', error);
      throw error;
    }
  }

  /**
   * Update user's hedge settings
   * @param {string} userId - User wallet address
   * @param {Object} settings - New hedge settings
   * @param {boolean} settings.auto_hedge_enabled - Enable/disable auto-hedging
   * @param {number} settings.risk_threshold - Risk threshold for triggering hedge
   * @param {number} settings.hedge_percentage - Percentage of position to hedge
   * @param {string} settings.hedge_strategy - Hedge strategy type
   * @returns {Promise<Object>} Update result
   */
  async updateUserSettings(userId, settings) {
    try {
      const response = await fetch(`${this.baseURL}/settings/${userId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(settings),
      });

      if (!response.ok) {
        throw new Error(`Failed to update hedge settings: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Hedge settings update error:', error);
      throw error;
    }
  }

  /**
   * Calculate hedge recommendation
   * @param {Object} position - Position data
   * @param {number} position.value - Position value in USD
   * @param {number} position.risk_score - Current risk score
   * @param {number} position.volatility - Asset volatility
   * @returns {Object} Hedge recommendation
   */
  calculateHedgeRecommendation(position) {
    const { value, risk_score, volatility } = position;

    // Simple hedge calculation logic
    let hedgePercentage = 0;
    let strategy = 'none';
    let urgency = 'low';

    if (risk_score > 70) {
      hedgePercentage = 50;
      strategy = 'put_options';
      urgency = 'high';
    } else if (risk_score > 50) {
      hedgePercentage = 30;
      strategy = 'short_position';
      urgency = 'medium';
    } else if (risk_score > 30) {
      hedgePercentage = 15;
      strategy = 'stablecoin_allocation';
      urgency = 'low';
    }

    const hedgeAmount = value * (hedgePercentage / 100);
    const estimatedCost = hedgeAmount * 0.02; // Assume 2% cost

    return {
      recommended: hedgePercentage > 0,
      hedge_percentage: hedgePercentage,
      hedge_amount: hedgeAmount,
      estimated_cost: estimatedCost,
      strategy,
      urgency,
      reason: this.getHedgeReason(risk_score, volatility),
    };
  }

  /**
   * Get reason for hedge recommendation
   * @param {number} riskScore - Current risk score
   * @param {number} volatility - Asset volatility
   * @returns {string} Human-readable reason
   */
  getHedgeReason(riskScore, volatility) {
    if (riskScore > 70) {
      return 'High risk exposure detected. Immediate hedging recommended to protect capital.';
    } else if (riskScore > 50) {
      return 'Moderate risk level. Consider partial hedging to reduce downside exposure.';
    } else if (riskScore > 30) {
      return 'Elevated market volatility. Light hedging may provide insurance.';
    }
    return 'Current risk levels are acceptable. No hedging required.';
  }

  /**
   * Get hedge strategy details
   * @param {string} strategy - Strategy name
   * @returns {Object} Strategy details
   */
  getStrategyDetails(strategy) {
    const strategies = {
      'put_options': {
        name: 'Put Options',
        description: 'Purchase put options to protect against downside',
        cost: 'Medium',
        effectiveness: 'High',
        complexity: 'High',
      },
      'short_position': {
        name: 'Short Position',
        description: 'Open short positions to offset long exposure',
        cost: 'Variable',
        effectiveness: 'Medium',
        complexity: 'Medium',
      },
      'stablecoin_allocation': {
        name: 'Stablecoin Allocation',
        description: 'Rotate portion of portfolio into stablecoins',
        cost: 'Low',
        effectiveness: 'Low',
        complexity: 'Low',
      },
      'inverse_etf': {
        name: 'Inverse ETF',
        description: 'Invest in inverse ETFs for market downturns',
        cost: 'Medium',
        effectiveness: 'Medium',
        complexity: 'Low',
      },
    };

    return strategies[strategy] || {
      name: 'Unknown Strategy',
      description: 'Strategy details not available',
      cost: 'N/A',
      effectiveness: 'N/A',
      complexity: 'N/A',
    };
  }

  /**
   * Format hedge amount as currency
   * @param {number} amount - Amount to format
   * @returns {string} Formatted currency string
   */
  formatAmount(amount) {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(amount);
  }
}

// Export singleton instance
const hedgeService = new HedgeService();
export default hedgeService;
