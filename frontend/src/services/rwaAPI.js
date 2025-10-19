// RWA API Service with Demo Mode Support

import { API_ENDPOINTS, apiCall } from './apiConfig';

class RWAAPI {
  // Get list of RWA assets
  async getAssets(type = '', page = 1, limit = 20) {
    const params = new URLSearchParams();
    if (type) params.append('type', type);
    params.append('page', page);
    params.append('limit', limit);

    return apiCall(`${API_ENDPOINTS.rwa.assets}?${params.toString()}`);
  }

  // Get asset detail
  async getAssetDetail(ticker) {
    return apiCall(API_ENDPOINTS.rwa.assetDetail(ticker));
  }

  // Create buy/sell order
  async createOrder(userAddress, ticker, orderType, orderStyle, amount, limitPrice) {
    return apiCall(API_ENDPOINTS.rwa.createOrder, {
      method: 'POST',
      body: JSON.stringify({
        user_address: userAddress,
        ticker,
        type: orderType,       // 'buy' or 'sell'
        style: orderStyle,     // 'market' or 'limit'
        amount,
        limitPrice,
      }),
    });
  }

  // Get user's orders
  async getUserOrders(address, status = '') {
    const params = status ? `?status=${status}` : '';
    return apiCall(`${API_ENDPOINTS.rwa.userOrders(address)}${params}`);
  }

  // Cancel an order
  async cancelOrder(orderId) {
    return apiCall(API_ENDPOINTS.rwa.cancelOrder(orderId), {
      method: 'DELETE',
    });
  }

  // Get user's holdings
  async getUserHoldings(address) {
    return apiCall(API_ENDPOINTS.rwa.holdings(address));
  }

  // Get price history
  async getPriceHistory(ticker, period = '7d') {
    return apiCall(`${API_ENDPOINTS.rwa.prices(ticker)}?period=${period}`);
  }

  // Calculate order value
  calculateOrderValue(amount, price) {
    return Math.round(amount * price * 100) / 100;
  }

  // Calculate profit/loss
  calculatePnL(currentPrice, averageCost, amount) {
    const currentValue = currentPrice * amount;
    const costBasis = averageCost * amount;
    const pnl = currentValue - costBasis;
    const pnlPercentage = ((currentValue / costBasis) - 1) * 100;

    return {
      pnl: Math.round(pnl * 100) / 100,
      pnlPercentage: Math.round(pnlPercentage * 100) / 100,
      currentValue: Math.round(currentValue * 100) / 100,
      costBasis: Math.round(costBasis * 100) / 100,
    };
  }

  // Format price with appropriate decimals
  formatPrice(price, assetType = 'stock') {
    if (assetType === 'commodity' && price > 1000) {
      return price.toLocaleString('en-US', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    }
    return price.toFixed(2);
  }

  // Format price change
  formatPriceChange(change) {
    const sign = change >= 0 ? '+' : '';
    return `${sign}${change.toFixed(2)}%`;
  }

  // Get asset type icon
  getAssetTypeIcon(assetType) {
    const icons = {
      stock: '📈',
      bond: '🏛️',
      commodity: '🥇',
      realestate: '🏠',
      other: '💰',
    };
    return icons[assetType] || '📊';
  }

  // Get asset type color
  getAssetTypeColor(assetType) {
    const colors = {
      stock: '#667eea',      // purple
      bond: '#10b981',       // green
      commodity: '#f59e0b',  // yellow
      realestate: '#3b82f6', // blue
      other: '#6b7280',      // gray
    };
    return colors[assetType] || '#6b7280';
  }

  // Validate order
  validateOrder(orderType, amount, balance, holdings) {
    if (orderType === 'buy') {
      if (amount > balance) {
        return { valid: false, error: 'Insufficient balance' };
      }
    } else if (orderType === 'sell') {
      if (amount > holdings) {
        return { valid: false, error: 'Insufficient holdings' };
      }
    }

    if (amount <= 0) {
      return { valid: false, error: 'Amount must be positive' };
    }

    return { valid: true };
  }

  // Get market status
  getMarketStatus(assetType) {
    const now = new Date();
    const hour = now.getUTCHours();
    const day = now.getUTCDay();

    if (assetType === 'stock') {
      // NYSE trading hours (9:30 AM - 4:00 PM ET)
      // Convert to UTC: 2:30 PM - 9:00 PM UTC
      const isWeekday = day >= 1 && day <= 5;
      const isTradingHours = hour >= 14.5 && hour < 21;

      return {
        isOpen: isWeekday && isTradingHours,
        message: isWeekday && isTradingHours
          ? 'Market Open'
          : 'Market Closed',
        nextOpen: this.getNextMarketOpen(day, hour),
      };
    }

    // Crypto/commodity markets are 24/7
    return {
      isOpen: true,
      message: '24/7 Trading',
      nextOpen: null,
    };
  }

  // Get next market open time
  getNextMarketOpen(currentDay, currentHour) {
    // Calculate next weekday at 14:30 UTC
    let daysUntilOpen = 0;

    if (currentDay === 0) {
      // Sunday
      daysUntilOpen = 1;
    } else if (currentDay === 6) {
      // Saturday
      daysUntilOpen = 2;
    } else if (currentHour >= 21) {
      // After market close
      daysUntilOpen = currentDay === 5 ? 3 : 1;
    }

    const nextOpen = new Date();
    nextOpen.setUTCDate(nextOpen.getUTCDate() + daysUntilOpen);
    nextOpen.setUTCHours(14, 30, 0, 0);

    return nextOpen;
  }

  // Filter assets by criteria
  filterAssets(assets, filters) {
    let filtered = [...assets];

    if (filters.minPrice) {
      filtered = filtered.filter(a => a.current_price >= filters.minPrice);
    }
    if (filters.maxPrice) {
      filtered = filtered.filter(a => a.current_price <= filters.maxPrice);
    }
    if (filters.issuer) {
      filtered = filtered.filter(a => a.issuer === filters.issuer);
    }
    if (filters.minChange) {
      filtered = filtered.filter(a => a.price_change_24h >= filters.minChange);
    }

    return filtered;
  }

  // Sort assets
  sortAssets(assets, sortBy = 'market_cap', order = 'desc') {
    const sorted = [...assets].sort((a, b) => {
      let aValue = a[sortBy];
      let bValue = b[sortBy];

      if (typeof aValue === 'string') {
        aValue = aValue.toLowerCase();
        bValue = bValue.toLowerCase();
      }

      if (order === 'asc') {
        return aValue > bValue ? 1 : -1;
      } else {
        return aValue < bValue ? 1 : -1;
      }
    });

    return sorted;
  }

  // Get portfolio summary
  calculatePortfolioSummary(holdings) {
    const summary = {
      totalValue: 0,
      totalPnL: 0,
      totalPnLPercentage: 0,
      assetCount: holdings.length,
      byType: {},
    };

    holdings.forEach(holding => {
      summary.totalValue += holding.current_value;
      summary.totalPnL += holding.pnl;

      // Group by asset type
      if (!summary.byType[holding.asset_type]) {
        summary.byType[holding.asset_type] = {
          count: 0,
          value: 0,
          pnl: 0,
        };
      }

      summary.byType[holding.asset_type].count++;
      summary.byType[holding.asset_type].value += holding.current_value;
      summary.byType[holding.asset_type].pnl += holding.pnl;
    });

    // Calculate overall PnL percentage
    const totalCost = summary.totalValue - summary.totalPnL;
    summary.totalPnLPercentage = totalCost > 0
      ? (summary.totalPnL / totalCost) * 100
      : 0;

    return summary;
  }
}

// Export singleton instance
const rwaAPI = new RWAAPI();
export default rwaAPI;