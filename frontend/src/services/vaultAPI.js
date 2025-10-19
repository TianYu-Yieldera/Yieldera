// Vault API Service with Demo Mode Support

import { API_ENDPOINTS, apiCall } from './apiConfig';

class VaultAPI {
  // Get user's vault balance
  async getBalance(address) {
    return apiCall(API_ENDPOINTS.vault.balance(address));
  }

  // Get available strategies
  async getStrategies() {
    return apiCall(API_ENDPOINTS.vault.strategies);
  }

  // Get DeFi protocols
  async getProtocols() {
    return apiCall(API_ENDPOINTS.vault.protocols);
  }

  // Deposit to vault
  async deposit(userAddress, amount, mode = 'smart', strategy = 'balanced') {
    return apiCall(API_ENDPOINTS.vault.deposit, {
      method: 'POST',
      body: JSON.stringify({
        user_address: userAddress,
        amount,
        mode,
        strategy,
      }),
    });
  }

  // Withdraw from vault
  async withdraw(userAddress, amount, emergency = false) {
    return apiCall(API_ENDPOINTS.vault.withdraw, {
      method: 'POST',
      body: JSON.stringify({
        user_address: userAddress,
        amount,
        emergency,
      }),
    });
  }

  // Manual stake to specific protocol
  async stakeToProtocol(userAddress, protocol, amount) {
    return apiCall(API_ENDPOINTS.vault.stake, {
      method: 'POST',
      body: JSON.stringify({
        user_address: userAddress,
        protocol,
        amount,
      }),
    });
  }

  // Manual unstake from specific protocol
  async unstakeFromProtocol(userAddress, protocol, amount) {
    return apiCall(API_ENDPOINTS.vault.unstake, {
      method: 'POST',
      body: JSON.stringify({
        user_address: userAddress,
        protocol,
        amount,
      }),
    });
  }

  // Get user's earnings history
  async getEarnings(address) {
    return apiCall(API_ENDPOINTS.vault.earnings(address));
  }

  // Get user's current positions
  async getPositions(address) {
    return apiCall(API_ENDPOINTS.vault.positions(address));
  }

  // Calculate estimated earnings
  calculateEstimatedEarnings(amount, apy, days = 365) {
    const dailyRate = apy / 100 / 365;
    const earnings = amount * dailyRate * days;
    return Math.round(earnings * 100) / 100;
  }

  // Format APY display
  formatAPY(apy) {
    return `${apy.toFixed(2)}%`;
  }

  // Format amount with commas
  formatAmount(amount) {
    return amount.toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });
  }

  // Get risk color based on level
  getRiskColor(riskLevel) {
    switch (riskLevel) {
      case 'low':
        return '#10b981';  // green
      case 'medium':
        return '#f59e0b';  // yellow
      case 'high':
        return '#ef4444';  // red
      default:
        return '#6b7280';  // gray
    }
  }

  // Get strategy description
  getStrategyDescription(strategy) {
    const descriptions = {
      conservative: 'Low risk, stable returns. Focuses on established lending protocols.',
      balanced: 'Moderate risk and returns. Mix of lending, DEX, and yield protocols.',
      aggressive: 'High risk, high returns. Includes derivatives and leveraged positions.',
    };
    return descriptions[strategy] || 'Custom strategy';
  }

  // Validate deposit amount
  validateDepositAmount(amount, minAmount = 100) {
    if (isNaN(amount) || amount <= 0) {
      return { valid: false, error: 'Amount must be positive' };
    }
    if (amount < minAmount) {
      return { valid: false, error: `Minimum amount is ${minAmount}` };
    }
    return { valid: true };
  }

  // Calculate fee for withdrawal
  calculateWithdrawalFee(amount, emergency = false) {
    const feeRate = emergency ? 0.02 : 0.005;  // 2% for emergency, 0.5% for normal
    const fee = amount * feeRate;
    const netAmount = amount - fee;
    return {
      fee: Math.round(fee * 100) / 100,
      netAmount: Math.round(netAmount * 100) / 100,
      feeRate: `${feeRate * 100}%`,
    };
  }

  // Get optimal strategy based on user profile
  getRecommendedStrategy(riskTolerance, investmentHorizon, amount) {
    if (riskTolerance === 'low' || investmentHorizon < 30) {
      return 'conservative';
    } else if (riskTolerance === 'high' && investmentHorizon > 180) {
      return 'aggressive';
    } else {
      return 'balanced';
    }
  }

  // Calculate portfolio allocation
  calculatePortfolioAllocation(amount, allocations) {
    const result = {};
    for (const [protocol, percentage] of Object.entries(allocations)) {
      result[protocol] = Math.round((amount * percentage / 100) * 100) / 100;
    }
    return result;
  }

  // Monitor position performance
  async monitorPerformance(address, interval = 60000) {
    // Set up polling for real-time updates
    const fetchData = async () => {
      const [balance, positions] = await Promise.all([
        this.getBalance(address),
        this.getPositions(address),
      ]);
      return { balance, positions, timestamp: new Date() };
    };

    // Initial fetch
    const data = await fetchData();

    // Set up interval for updates (optional, can be managed by React component)
    if (interval > 0) {
      setInterval(fetchData, interval);
    }

    return data;
  }
}

// Export singleton instance
const vaultAPI = new VaultAPI();
export default vaultAPI;