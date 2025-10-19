/**
 * Environment configuration for the application
 * All environment variables are prefixed with VITE_ and accessed via import.meta.env
 */

export const config = {
  // API Configuration
  api: {
    baseUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080',
    timeout: parseInt(import.meta.env.VITE_API_TIMEOUT) || 30000,
    graphqlUrl: import.meta.env.VITE_GRAPHQL_URL || 'http://localhost:8080/graphql',
  },

  // Blockchain Configuration
  blockchain: {
    chainId: parseInt(import.meta.env.VITE_CHAIN_ID) || 11155111,
    chainName: import.meta.env.VITE_CHAIN_NAME || 'Sepolia',
    rpcUrl: import.meta.env.VITE_RPC_URL || 'https://sepolia.infura.io/v3/YOUR_INFURA_KEY',
  },

  // Smart Contract Addresses
  contracts: {
    lusd: import.meta.env.VITE_LUSD_ADDRESS || '0x0000000000000000000000000000000000000000',
    vault: import.meta.env.VITE_VAULT_ADDRESS || '0x0000000000000000000000000000000000000000',
    liquidation: import.meta.env.VITE_LIQUIDATION_ADDRESS || '0x0000000000000000000000000000000000000000',
    badgeNFT: import.meta.env.VITE_BADGE_NFT_ADDRESS || '0x0000000000000000000000000000000000000000',
  },

  // DeFi Pool Addresses
  pools: {
    uniswap: import.meta.env.VITE_UNISWAP_POOL_ADDRESS || '0x0000000000000000000000000000000000000000',
    aave: import.meta.env.VITE_AAVE_POOL_ADDRESS || '0x0000000000000000000000000000000000000000',
    staking: import.meta.env.VITE_STAKING_POOL_ADDRESS || '0x0000000000000000000000000000000000000000',
  },

  // Feature Flags
  features: {
    enableTestnet: import.meta.env.VITE_ENABLE_TESTNET === 'true',
    enableDebug: import.meta.env.VITE_ENABLE_DEBUG === 'true',
  },
};

// Validation function
export function validateConfig() {
  const errors = [];

  if (!config.api.baseUrl) {
    errors.push('API base URL is not configured');
  }

  if (!config.blockchain.rpcUrl || config.blockchain.rpcUrl.includes('YOUR_')) {
    errors.push('Blockchain RPC URL is not properly configured');
  }

  if (config.features.enableDebug) {
    console.log('🔧 Configuration:', config);
  }

  if (errors.length > 0) {
    console.warn('⚠️ Configuration warnings:', errors);
  }

  return errors.length === 0;
}

// Initialize configuration on module load
if (config.features.enableDebug) {
  console.log('🚀 Application configuration loaded');
  validateConfig();
}

export default config;
