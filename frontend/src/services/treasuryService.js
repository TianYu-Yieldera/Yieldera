const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

/**
 * Treasury Service - Handles all Treasury-related API calls
 */
class TreasuryService {
  /**
   * Get all treasury assets
   * @param {string} type - Optional filter by type (T-BILL, T-NOTE, T-BOND)
   * @returns {Promise<Object>}
   */
  async getAssets(type = null) {
    const url = type
      ? `${API_BASE_URL}/api/v1/treasury/assets?type=${type}`
      : `${API_BASE_URL}/api/v1/treasury/assets`;

    const response = await fetch(url);
    if (!response.ok) {
      throw new Error('Failed to fetch treasury assets');
    }
    return response.json();
  }

  /**
   * Get single treasury asset by ID
   * @param {number} assetId
   * @returns {Promise<Object>}
   */
  async getAsset(assetId) {
    const response = await fetch(`${API_BASE_URL}/api/v1/treasury/assets/${assetId}`);
    if (!response.ok) {
      throw new Error('Failed to fetch treasury asset');
    }
    return response.json();
  }

  /**
   * Get price history for an asset
   * @param {number} assetId
   * @param {number} limit - Number of records to return
   * @returns {Promise<Object>}
   */
  async getPriceHistory(assetId, limit = 100) {
    const response = await fetch(
      `${API_BASE_URL}/api/v1/treasury/assets/${assetId}/price-history?limit=${limit}`
    );
    if (!response.ok) {
      throw new Error('Failed to fetch price history');
    }
    return response.json();
  }

  /**
   * Get trade history for an asset
   * @param {number} assetId
   * @param {number} limit
   * @returns {Promise<Object>}
   */
  async getTradeHistory(assetId, limit = 50) {
    const response = await fetch(
      `${API_BASE_URL}/api/v1/treasury/assets/${assetId}/trades?limit=${limit}`
    );
    if (!response.ok) {
      throw new Error('Failed to fetch trade history');
    }
    return response.json();
  }

  /**
   * Get user's treasury holdings
   * @param {string} address - User wallet address
   * @returns {Promise<Object>}
   */
  async getUserHoldings(address) {
    const response = await fetch(
      `${API_BASE_URL}/api/v1/treasury/user/${address}/holdings`
    );
    if (!response.ok) {
      throw new Error('Failed to fetch user holdings');
    }
    return response.json();
  }

  /**
   * Get user's yield information
   * @param {string} address
   * @returns {Promise<Object>}
   */
  async getUserYield(address) {
    const response = await fetch(
      `${API_BASE_URL}/api/v1/treasury/user/${address}/yield`
    );
    if (!response.ok) {
      throw new Error('Failed to fetch user yield');
    }
    return response.json();
  }

  /**
   * Get market orders for an asset
   * @param {number} assetId
   * @param {string} type - Optional filter by order type (BUY, SELL)
   * @returns {Promise<Object>}
   */
  async getMarketOrders(assetId, type = null) {
    const url = type
      ? `${API_BASE_URL}/api/v1/treasury/market/${assetId}/orders?type=${type}`
      : `${API_BASE_URL}/api/v1/treasury/market/${assetId}/orders`;

    const response = await fetch(url);
    if (!response.ok) {
      throw new Error('Failed to fetch market orders');
    }
    return response.json();
  }

  /**
   * Create a buy or sell order
   * @param {Object} orderData
   * @returns {Promise<Object>}
   */
  async createOrder(orderData) {
    const response = await fetch(`${API_BASE_URL}/api/v1/treasury/market/order`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(orderData),
    });

    if (!response.ok) {
      throw new Error('Failed to create order');
    }
    return response.json();
  }

  /**
   * Cancel an order
   * @param {number} orderId
   * @returns {Promise<Object>}
   */
  async cancelOrder(orderId) {
    const response = await fetch(
      `${API_BASE_URL}/api/v1/treasury/market/order/${orderId}`,
      { method: 'DELETE' }
    );

    if (!response.ok) {
      throw new Error('Failed to cancel order');
    }
    return response.json();
  }

  /**
   * Claim yield for an asset
   * @param {Object} claimData
   * @returns {Promise<Object>}
   */
  async claimYield(claimData) {
    const response = await fetch(`${API_BASE_URL}/api/v1/treasury/yield/claim`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(claimData),
    });

    if (!response.ok) {
      throw new Error('Failed to claim yield');
    }
    return response.json();
  }

  /**
   * Get yield distributions
   * @param {number} assetId - Optional filter by asset
   * @param {number} limit
   * @returns {Promise<Object>}
   */
  async getYieldDistributions(assetId = null, limit = 20) {
    const url = assetId
      ? `${API_BASE_URL}/api/v1/treasury/yield/distributions?asset_id=${assetId}&limit=${limit}`
      : `${API_BASE_URL}/api/v1/treasury/yield/distributions?limit=${limit}`;

    const response = await fetch(url);
    if (!response.ok) {
      throw new Error('Failed to fetch yield distributions');
    }
    return response.json();
  }

  /**
   * Get treasury statistics
   * @returns {Promise<Object>}
   */
  async getStats() {
    const response = await fetch(`${API_BASE_URL}/api/v1/treasury/stats`);
    if (!response.ok) {
      throw new Error('Failed to fetch treasury stats');
    }
    return response.json();
  }
}

export default new TreasuryService();
