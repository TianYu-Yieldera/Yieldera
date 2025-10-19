// API Configuration with Demo Mode Support

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// Check if demo mode is enabled from DemoModeContext
export const getApiMode = () => {
  const demoMode = localStorage.getItem('demoMode') === 'true';
  return demoMode ? 'demo' : 'real';
};

// API endpoints configuration
export const API_ENDPOINTS = {
  // Vault endpoints
  vault: {
    deposit: '/api/vault/deposit',
    withdraw: '/api/vault/withdraw',
    balance: (address) => `/api/vault/balance/${address}`,
    strategies: '/api/vault/strategies',
    stake: '/api/vault/stake',
    unstake: '/api/vault/unstake',
    earnings: (address) => `/api/vault/earnings/${address}`,
    positions: (address) => `/api/vault/positions/${address}`,
    protocols: '/api/vault/protocols',
  },
  // RWA endpoints
  rwa: {
    assets: '/api/rwa/assets',
    assetDetail: (ticker) => `/api/rwa/assets/${ticker}`,
    createOrder: '/api/rwa/orders',
    userOrders: (address) => `/api/rwa/orders/${address}`,
    cancelOrder: (orderId) => `/api/rwa/orders/${orderId}`,
    holdings: (address) => `/api/rwa/holdings/${address}`,
    prices: (ticker) => `/api/rwa/prices/${ticker}`,
  },
  // Oracle endpoints
  oracle: {
    price: (ticker) => `/api/oracle/price/${ticker}`,
    prices: '/api/oracle/prices',
    apy: (protocol) => `/api/oracle/apy/${protocol}`,
    apys: '/api/oracle/apys',
    stats: '/api/oracle/stats',
  },
};

// Helper function to make API calls with proper error handling
export const apiCall = async (endpoint, options = {}) => {
  const mode = getApiMode();

  // If in demo mode and not explicitly requesting real API
  if (mode === 'demo' && !options.forceReal) {
    // Return mock data for demo mode
    return getMockResponse(endpoint, options);
  }

  // Real API call
  try {
    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error('API call failed:', error);

    // Fallback to mock data if real API fails and demo mode is available
    if (!options.noFallback) {
      console.log('Falling back to mock data...');
      return getMockResponse(endpoint, options);
    }

    throw error;
  }
};

// Mock data responses for demo mode
const getMockResponse = (endpoint, options) => {
  // Parse endpoint to determine response

  // Vault balance
  if (endpoint.includes('/api/vault/balance/')) {
    return Promise.resolve({
      available: 15000,
      locked: 50000,
      earned: 3247.80,
      positions: [
        { protocol: 'Aave V3', amount: 20000, earned: 1200, apy: 3.52 },
        { protocol: 'Compound V3', amount: 15000, earned: 800, apy: 4.18 },
        { protocol: 'Uniswap V3', amount: 10000, earned: 1247.80, apy: 12.85 },
      ],
    });
  }

  // Vault strategies
  if (endpoint === '/api/vault/strategies') {
    return Promise.resolve({
      strategies: [
        {
          name: 'Conservative Strategy',
          mode: 'conservative',
          allocations: { 'Aave V3': 40, 'Compound V3': 30, 'Curve': 20, 'Lido': 10 },
          min_amount: 100,
        },
        {
          name: 'Balanced Strategy',
          mode: 'balanced',
          allocations: { 'Aave V3': 25, 'Compound V3': 20, 'Uniswap V3': 20, 'Yearn Finance': 20, 'Curve': 15 },
          min_amount: 100,
        },
        {
          name: 'Aggressive Strategy',
          mode: 'aggressive',
          allocations: { 'Uniswap V3': 25, 'GMX': 20, 'Yearn Finance': 25, 'Aave V3': 15, 'Rocket Pool': 15 },
          min_amount: 100,
        },
      ],
    });
  }

  // Vault protocols
  if (endpoint === '/api/vault/protocols') {
    return Promise.resolve({
      protocols: [
        { name: 'Aave V3', protocol_type: 'lending', risk_level: 'low', current_apy: 3.52, tvl: 5000000000 },
        { name: 'Compound V3', protocol_type: 'lending', risk_level: 'low', current_apy: 4.18, tvl: 3000000000 },
        { name: 'Uniswap V3', protocol_type: 'dex', risk_level: 'medium', current_apy: 12.85, tvl: 7000000000 },
        { name: 'GMX', protocol_type: 'derivatives', risk_level: 'high', current_apy: 22.30, tvl: 500000000 },
      ],
    });
  }

  // RWA assets
  if (endpoint.includes('/api/rwa/assets')) {
    const mockAssets = {
      stocks: [
        {
          ticker: 'bAAPL',
          name: 'Apple Inc.',
          asset_type: 'stock',
          issuer: 'Backed Finance',
          current_price: 178.52,
          price_change_24h: 2.34,
          metadata: { logo: '🍎', description: 'Technology giant' },
        },
        {
          ticker: 'bTSLA',
          name: 'Tesla Inc.',
          asset_type: 'stock',
          issuer: 'Backed Finance',
          current_price: 242.18,
          price_change_24h: 5.67,
          metadata: { logo: '⚡', description: 'Electric vehicle leader' },
        },
        {
          ticker: 'bNVDA',
          name: 'NVIDIA Corp.',
          asset_type: 'stock',
          issuer: 'Backed Finance',
          current_price: 485.93,
          price_change_24h: 8.92,
          metadata: { logo: '🎮', description: 'AI chip leader' },
        },
      ],
      bonds: [
        {
          ticker: 'OUSG',
          name: 'Ondo Short-Term US Treasuries',
          asset_type: 'bond',
          issuer: 'Ondo Finance',
          current_price: 105.20,
          price_change_24h: 0.12,
          metadata: { logo: '🏛️', description: 'Short-term US Treasury token', apy: '4.5%' },
        },
      ],
      commodities: [
        {
          ticker: 'PAXG',
          name: 'Paxos Gold',
          asset_type: 'commodity',
          issuer: 'Paxos',
          current_price: 2042.50,
          price_change_24h: 1.05,
          metadata: { logo: '🥇', description: '1 PAXG = 1 troy ounce of gold' },
        },
      ],
    };

    const url = new URL(endpoint, 'http://example.com');
    const assetType = url.searchParams.get('type');

    let assets = [];
    if (assetType === 'stock') {
      assets = mockAssets.stocks;
    } else if (assetType === 'bond') {
      assets = mockAssets.bonds;
    } else if (assetType === 'commodity') {
      assets = mockAssets.commodities;
    } else {
      assets = [...mockAssets.stocks, ...mockAssets.bonds, ...mockAssets.commodities];
    }

    return Promise.resolve({ assets });
  }

  // RWA holdings
  if (endpoint.includes('/api/rwa/holdings/')) {
    return Promise.resolve({
      holdings: [
        {
          ticker: 'bAAPL',
          name: 'Apple Inc.',
          asset_type: 'stock',
          amount: 5.2,
          average_cost: 175.00,
          current_value: 928.30,
          pnl: 18.30,
          pnl_percentage: 2.01,
          current_price: 178.52,
          price_change_24h: 2.34,
        },
        {
          ticker: 'PAXG',
          name: 'Paxos Gold',
          asset_type: 'commodity',
          amount: 0.5,
          average_cost: 2000.00,
          current_value: 1021.25,
          pnl: 21.25,
          pnl_percentage: 2.13,
          current_price: 2042.50,
          price_change_24h: 1.05,
        },
      ],
      total_value: 1949.55,
      total_pnl: 39.55,
      total_assets: 2,
    });
  }

  // Oracle APYs
  if (endpoint === '/api/oracle/apys') {
    return Promise.resolve({
      protocols: [
        { protocol: 'GMX', protocol_type: 'derivatives', risk_level: 'high', current_apy: 22.30, tvl: 500000000 },
        { protocol: 'Uniswap V3', protocol_type: 'dex', risk_level: 'medium', current_apy: 12.85, tvl: 7000000000 },
        { protocol: 'Yearn Finance', protocol_type: 'yield', risk_level: 'medium', current_apy: 10.50, tvl: 1000000000 },
        { protocol: 'Curve', protocol_type: 'dex', risk_level: 'low', current_apy: 5.20, tvl: 4000000000 },
        { protocol: 'Compound V3', protocol_type: 'lending', risk_level: 'low', current_apy: 4.18, tvl: 3000000000 },
        { protocol: 'Aave V3', protocol_type: 'lending', risk_level: 'low', current_apy: 3.52, tvl: 5000000000 },
      ],
      averages: { low: 4.30, medium: 11.68, high: 22.30 },
      count: 6,
    });
  }

  // Oracle market stats
  if (endpoint === '/api/oracle/stats') {
    return Promise.resolve({
      market_stats: {
        total_rwa_market_cap: '$5.20B',
        total_vault_tvl: '$25.50B',
        average_apy: '8.45%',
        top_gainer: { ticker: 'bNVDA', change: '+8.92%' },
        top_loser: { ticker: 'bGOOGL', change: '-1.23%' },
        active_assets: 12,
        active_protocols: 8,
        total_users: 42341,
        total_transactions: 156789,
      },
    });
  }

  // Default mock response for operations
  if (options.method === 'POST') {
    return Promise.resolve({
      success: true,
      message: 'Demo mode: Operation simulated successfully',
      data: options.body,
    });
  }

  // Default response
  return Promise.resolve({
    success: true,
    message: 'Mock data returned',
    mode: 'demo',
  });
};

export default {
  API_BASE,
  API_ENDPOINTS,
  apiCall,
  getApiMode,
};