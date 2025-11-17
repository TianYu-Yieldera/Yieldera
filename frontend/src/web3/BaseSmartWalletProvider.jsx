/**
 * Base Smart Wallet Provider
 *
 * Integrates Coinbase Smart Wallet (ERC-4337) for Base chain
 *
 * Features:
 * - Passkey login (no seed phrase)
 * - Social login (Google, Apple, Email)
 * - Gasless transactions via Paymaster
 * - Seamless Web2-like UX
 */

import React, { createContext, useContext, useState, useEffect } from 'react';
import { CoinbaseWalletSDK } from '@coinbase/wallet-sdk';
import { createConfig, WagmiProvider, useAccount, useConnect, useDisconnect } from 'wagmi';
import { base, baseSepolia } from 'wagmi/chains';
import { coinbaseWallet } from 'wagmi/connectors';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { http } from 'viem';

// ============================================================================
// Configuration
// ============================================================================

const PROJECT_NAME = 'Yieldera - DeFi Risk Management';
const PROJECT_LOGO = '/pointfi-logo-mark.svg';
const APP_URL = window.location.origin;

// Determine environment
const IS_TESTNET = import.meta.env.VITE_NETWORK === 'testnet';
const BASE_CHAIN = IS_TESTNET ? baseSepolia : base;

// ============================================================================
// Wagmi Config
// ============================================================================

const queryClient = new QueryClient();

const wagmiConfig = createConfig({
  chains: [BASE_CHAIN],
  transports: {
    [BASE_CHAIN.id]: http()
  },
  connectors: [
    coinbaseWallet({
      appName: PROJECT_NAME,
      appLogoUrl: `${APP_URL}${PROJECT_LOGO}`,
      preference: 'smartWalletOnly', // Force Smart Wallet only
      version: '4',
      enableMobileWalletLink: true,
    })
  ]
});

// ============================================================================
// Context
// ============================================================================

const BaseSmartWalletContext = createContext(null);

export function BaseSmartWalletProvider({ children }) {
  return (
    <WagmiProvider config={wagmiConfig}>
      <QueryClientProvider client={queryClient}>
        <BaseSmartWalletInner>
          {children}
        </BaseSmartWalletInner>
      </QueryClientProvider>
    </WagmiProvider>
  );
}

function BaseSmartWalletInner({ children }) {
  const { address, isConnected, chain } = useAccount();
  const { connect, connectors } = useConnect();
  const { disconnect: wagmiDisconnect } = useDisconnect();

  const [walletInfo, setWalletInfo] = useState({
    address: null,
    chainId: null,
    isSmartWallet: false,
    isConnected: false
  });

  // Update wallet info when account changes
  useEffect(() => {
    if (address && chain) {
      setWalletInfo({
        address: address,
        chainId: chain.id,
        isSmartWallet: true, // Always true for this provider
        isConnected: isConnected
      });
    } else {
      setWalletInfo({
        address: null,
        chainId: null,
        isSmartWallet: false,
        isConnected: false
      });
    }
  }, [address, chain, isConnected]);

  /**
   * Connect to Smart Wallet
   * Opens Coinbase Smart Wallet connection flow
   */
  const connectSmartWallet = async () => {
    try {
      const coinbaseConnector = connectors.find(c => c.id === 'coinbaseWalletSDK');

      if (!coinbaseConnector) {
        throw new Error('Coinbase Wallet connector not found');
      }

      await connect({ connector: coinbaseConnector });

      console.log('✅ Smart Wallet connected:', address);
    } catch (error) {
      console.error('❌ Failed to connect Smart Wallet:', error);
      throw error;
    }
  };

  /**
   * Disconnect Smart Wallet
   */
  const disconnect = async () => {
    try {
      await wagmiDisconnect();
      console.log('✅ Smart Wallet disconnected');
    } catch (error) {
      console.error('❌ Failed to disconnect:', error);
    }
  };

  /**
   * Check if on correct chain (Base)
   */
  const isOnBaseChain = () => {
    if (!walletInfo.chainId) return false;

    if (IS_TESTNET) {
      return walletInfo.chainId === baseSepolia.id; // 84532
    } else {
      return walletInfo.chainId === base.id; // 8453
    }
  };

  /**
   * Get user-friendly error messages
   */
  const getErrorMessage = (error) => {
    if (error.message?.includes('User rejected')) {
      return '用户取消了连接请求';
    }
    if (error.message?.includes('Already processing')) {
      return '请先完成当前连接流程';
    }
    return error.message || '连接失败，请重试';
  };

  const contextValue = {
    // Wallet state
    address: walletInfo.address,
    account: walletInfo.address, // Alias for compatibility
    chainId: walletInfo.chainId,
    isConnected: walletInfo.isConnected,
    isSmartWallet: walletInfo.isSmartWallet,

    // Chain info
    chain: BASE_CHAIN,
    isOnBaseChain: isOnBaseChain(),

    // Actions
    connect: connectSmartWallet,
    disconnect,

    // Utils
    getErrorMessage
  };

  return (
    <BaseSmartWalletContext.Provider value={contextValue}>
      {children}
    </BaseSmartWalletContext.Provider>
  );
}

// ============================================================================
// Hook
// ============================================================================

export function useBaseSmartWallet() {
  const context = useContext(BaseSmartWalletContext);

  if (!context) {
    throw new Error('useBaseSmartWallet must be used within BaseSmartWalletProvider');
  }

  return context;
}

// ============================================================================
// Export
// ============================================================================

export default BaseSmartWalletProvider;
