/**
 * Unified Wallet Manager
 *
 * 双钱包策略管理器:
 * - Arbitrum: MetaMask (EOA)
 * - Base: Smart Wallet (ERC-4337)
 *
 * 根据链自动切换钱包类型
 */

import React, { createContext, useContext, useState, useEffect } from 'react';
import { WalletProvider as EOAWalletProvider, useWallet as useEOAWallet } from './WalletContext';
import { BaseSmartWalletProvider, useBaseSmartWallet } from './BaseSmartWalletProvider';

// ============================================================================
// Chain IDs
// ============================================================================

const CHAINS = {
  // Arbitrum
  ARBITRUM_SEPOLIA: 421614,
  ARBITRUM_ONE: 42161,

  // Base
  BASE_SEPOLIA: 84532,
  BASE_MAINNET: 8453,

  // Ethereum
  SEPOLIA: 11155111,
  MAINNET: 1
};

// ============================================================================
// Wallet Type Detection
// ============================================================================

function getWalletTypeForChain(chainId) {
  if (!chainId) return null;

  // Base chains → Smart Wallet
  if (chainId === CHAINS.BASE_SEPOLIA || chainId === CHAINS.BASE_MAINNET) {
    return 'smartwallet';
  }

  // Arbitrum chains → EOA (MetaMask)
  if (chainId === CHAINS.ARBITRUM_SEPOLIA || chainId === CHAINS.ARBITRUM_ONE) {
    return 'eoa';
  }

  // Default to EOA for other chains
  return 'eoa';
}

function getChainName(chainId) {
  switch (chainId) {
    case CHAINS.ARBITRUM_SEPOLIA: return 'Arbitrum Sepolia';
    case CHAINS.ARBITRUM_ONE: return 'Arbitrum One';
    case CHAINS.BASE_SEPOLIA: return 'Base Sepolia';
    case CHAINS.BASE_MAINNET: return 'Base';
    case CHAINS.SEPOLIA: return 'Sepolia';
    case CHAINS.MAINNET: return 'Ethereum';
    default: return `Chain ${chainId}`;
  }
}

function isBaseChain(chainId) {
  return chainId === CHAINS.BASE_SEPOLIA || chainId === CHAINS.BASE_MAINNET;
}

function isArbitrumChain(chainId) {
  return chainId === CHAINS.ARBITRUM_SEPOLIA || chainId === CHAINS.ARBITRUM_ONE;
}

// ============================================================================
// Unified Context
// ============================================================================

const UnifiedWalletContext = createContext(null);

export function UnifiedWalletProvider({ children }) {
  const [targetChain, setTargetChain] = useState(null);
  const [walletMode, setWalletMode] = useState('auto'); // 'auto' | 'eoa' | 'smartwallet'

  // Determine wallet type based on target chain
  const activeWalletType = walletMode === 'auto'
    ? getWalletTypeForChain(targetChain)
    : walletMode;

  return (
    <UnifiedWalletContext.Provider value={{ targetChain, setTargetChain, walletMode, setWalletMode }}>
      {activeWalletType === 'smartwallet' ? (
        <BaseSmartWalletProvider>
          <UnifiedWalletAdapter walletType="smartwallet">
            {children}
          </UnifiedWalletAdapter>
        </BaseSmartWalletProvider>
      ) : (
        <EOAWalletProvider>
          <UnifiedWalletAdapter walletType="eoa">
            {children}
          </UnifiedWalletAdapter>
        </EOAWalletProvider>
      )}
    </UnifiedWalletContext.Provider>
  );
}

// ============================================================================
// Unified Adapter
// ============================================================================

function UnifiedWalletAdapter({ walletType, children }) {
  const { targetChain, setTargetChain } = useContext(UnifiedWalletContext);

  // Get wallet context based on type
  const eoaWallet = walletType === 'eoa' ? useEOAWallet() : null;
  const smartWallet = walletType === 'smartwallet' ? useBaseSmartWallet() : null;

  const activeWallet = walletType === 'eoa' ? eoaWallet : smartWallet;

  // Monitor chain changes
  useEffect(() => {
    if (activeWallet?.chainId && activeWallet.chainId !== targetChain) {
      setTargetChain(activeWallet.chainId);
    }
  }, [activeWallet?.chainId, targetChain, setTargetChain]);

  // Create unified interface
  const unifiedContext = {
    // Core wallet info
    address: activeWallet?.address || null,
    account: activeWallet?.address || null, // Alias
    chainId: activeWallet?.chainId || null,
    isConnected: activeWallet?.isConnected || false,

    // Wallet type info
    walletType: walletType,
    isSmartWallet: walletType === 'smartwallet',
    isEOA: walletType === 'eoa',

    // Chain utilities
    chainName: getChainName(activeWallet?.chainId),
    isBaseChain: isBaseChain(activeWallet?.chainId),
    isArbitrumChain: isArbitrumChain(activeWallet?.chainId),

    // Actions
    connect: activeWallet?.connect,
    disconnect: activeWallet?.disconnect,
    signer: eoaWallet?.signer || null, // Only available for EOA

    // Chain switching (EOA only)
    switchToSepolia: eoaWallet?.switchToSepolia,

    // Utils
    CHAINS,
    getChainName,
    isBaseChain,
    isArbitrumChain
  };

  return (
    <UnifiedWalletContext.Provider value={{ ...unifiedContext, targetChain, setTargetChain }}>
      {children}
    </UnifiedWalletContext.Provider>
  );
}

// ============================================================================
// Hook
// ============================================================================

export function useUnifiedWallet() {
  const context = useContext(UnifiedWalletContext);

  if (!context) {
    throw new Error('useUnifiedWallet must be used within UnifiedWalletProvider');
  }

  return context;
}

// ============================================================================
// Convenience Hooks
// ============================================================================

/**
 * Hook to check if user can purchase US Treasury Bonds
 * Only allowed on Base chain with Smart Wallet
 */
export function useCanPurchaseTreasury() {
  const { chainId, isSmartWallet } = useUnifiedWallet();

  return isBaseChain(chainId) && isSmartWallet;
}

/**
 * Hook to enforce Base chain for Treasury operations
 * Automatically shows warning if not on Base
 */
export function useEnforceBaseChain() {
  const { chainId, isBaseChain, connect } = useUnifiedWallet();
  const [showWarning, setShowWarning] = useState(false);

  useEffect(() => {
    if (chainId && !isBaseChain) {
      setShowWarning(true);
    } else {
      setShowWarning(false);
    }
  }, [chainId, isBaseChain]);

  const switchToBase = () => {
    // This will trigger UnifiedWalletProvider to switch to Smart Wallet
    setShowWarning(false);
    connect();
  };

  return {
    isOnBase: isBaseChain,
    showWarning,
    switchToBase
  };
}

// ============================================================================
// Export
// ============================================================================

export default UnifiedWalletProvider;
export { CHAINS };
