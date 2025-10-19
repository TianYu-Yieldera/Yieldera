/**
 * Web3 Utilities for Frontend
 * Handles MetaMask connection, contract interactions, and transaction management
 */

import { config } from '../config/env';

// Check if MetaMask is installed
export const isMetaMaskInstalled = () => {
  return typeof window !== 'undefined' && typeof window.ethereum !== 'undefined';
};

// Connect to MetaMask
export const connectWallet = async () => {
  if (!isMetaMaskInstalled()) {
    throw new Error('MetaMask is not installed. Please install MetaMask to continue.');
  }

  try {
    const accounts = await window.ethereum.request({
      method: 'eth_requestAccounts',
    });

    if (accounts.length === 0) {
      throw new Error('No accounts found. Please unlock MetaMask.');
    }

    const address = accounts[0];
    const chainId = await window.ethereum.request({ method: 'eth_chainId' });

    return {
      address,
      chainId: parseInt(chainId, 16),
    };
  } catch (error) {
    console.error('Failed to connect wallet:', error);
    throw error;
  }
};

// Get current account
export const getCurrentAccount = async () => {
  if (!isMetaMaskInstalled()) {
    return null;
  }

  try {
    const accounts = await window.ethereum.request({
      method: 'eth_accounts',
    });
    return accounts[0] || null;
  } catch (error) {
    console.error('Failed to get current account:', error);
    return null;
  }
};

// Switch network
export const switchNetwork = async (chainId) => {
  if (!isMetaMaskInstalled()) {
    throw new Error('MetaMask is not installed');
  }

  try {
    await window.ethereum.request({
      method: 'wallet_switchEthereumChain',
      params: [{ chainId: `0x${chainId.toString(16)}` }],
    });
  } catch (error) {
    // This error code indicates that the chain has not been added to MetaMask
    if (error.code === 4902) {
      await addNetwork(chainId);
    } else {
      throw error;
    }
  }
};

// Add network to MetaMask
export const addNetwork = async (chainId) => {
  const networks = {
    11155111: {
      chainId: '0xaa36a7',
      chainName: 'Sepolia Testnet',
      nativeCurrency: {
        name: 'Sepolia ETH',
        symbol: 'ETH',
        decimals: 18,
      },
      rpcUrls: [config.blockchain.rpcUrl],
      blockExplorerUrls: ['https://sepolia.etherscan.io'],
    },
  };

  const network = networks[chainId];
  if (!network) {
    throw new Error(`Network ${chainId} not supported`);
  }

  try {
    await window.ethereum.request({
      method: 'wallet_addEthereumChain',
      params: [network],
    });
  } catch (error) {
    console.error('Failed to add network:', error);
    throw error;
  }
};

// Get balance
export const getBalance = async (address) => {
  if (!isMetaMaskInstalled()) {
    throw new Error('MetaMask is not installed');
  }

  try {
    const balance = await window.ethereum.request({
      method: 'eth_getBalance',
      params: [address, 'latest'],
    });
    return parseInt(balance, 16);
  } catch (error) {
    console.error('Failed to get balance:', error);
    throw error;
  }
};

// Sign message
export const signMessage = async (message, address) => {
  if (!isMetaMaskInstalled()) {
    throw new Error('MetaMask is not installed');
  }

  try {
    const signature = await window.ethereum.request({
      method: 'personal_sign',
      params: [message, address],
    });
    return signature;
  } catch (error) {
    console.error('Failed to sign message:', error);
    throw error;
  }
};

// Format address for display
export const formatAddress = (address) => {
  if (!address) return '';
  return `${address.substring(0, 6)}...${address.substring(address.length - 4)}`;
};

// Format amount with decimals
export const formatAmount = (amount, decimals = 6) => {
  if (!amount) return '0';
  const divisor = Math.pow(10, decimals);
  return (parseInt(amount) / divisor).toFixed(2);
};

// Parse amount to wei
export const parseAmount = (amount, decimals = 6) => {
  const multiplier = Math.pow(10, decimals);
  return Math.floor(parseFloat(amount) * multiplier).toString();
};

// Wait for transaction
export const waitForTransaction = async (txHash, onUpdate) => {
  if (!isMetaMaskInstalled()) {
    throw new Error('MetaMask is not installed');
  }

  let attempts = 0;
  const maxAttempts = 60; // 2 minutes with 2 second intervals

  return new Promise((resolve, reject) => {
    const checkTransaction = async () => {
      try {
        const receipt = await window.ethereum.request({
          method: 'eth_getTransactionReceipt',
          params: [txHash],
        });

        if (receipt) {
          if (receipt.status === '0x1') {
            onUpdate?.({ status: 'success', receipt });
            resolve(receipt);
          } else {
            onUpdate?.({ status: 'failed', receipt });
            reject(new Error('Transaction failed'));
          }
        } else {
          attempts++;
          if (attempts >= maxAttempts) {
            onUpdate?.({ status: 'timeout' });
            reject(new Error('Transaction timeout'));
          } else {
            onUpdate?.({ status: 'pending', attempts });
            setTimeout(checkTransaction, 2000);
          }
        }
      } catch (error) {
        onUpdate?.({ status: 'error', error });
        reject(error);
      }
    };

    checkTransaction();
  });
};

// Listen for account changes
export const onAccountsChanged = (callback) => {
  if (!isMetaMaskInstalled()) return;

  window.ethereum.on('accountsChanged', (accounts) => {
    callback(accounts[0] || null);
  });
};

// Listen for chain changes
export const onChainChanged = (callback) => {
  if (!isMetaMaskInstalled()) return;

  window.ethereum.on('chainChanged', (chainId) => {
    callback(parseInt(chainId, 16));
  });
};

// Disconnect listeners
export const removeListeners = () => {
  if (!isMetaMaskInstalled()) return;

  window.ethereum.removeAllListeners('accountsChanged');
  window.ethereum.removeAllListeners('chainChanged');
};

// Contract interaction helper
export const callContract = async (contractAddress, method, params = []) => {
  if (!isMetaMaskInstalled()) {
    throw new Error('MetaMask is not installed');
  }

  try {
    const result = await window.ethereum.request({
      method: 'eth_call',
      params: [{
        to: contractAddress,
        data: method + params.join(''),
      }, 'latest'],
    });
    return result;
  } catch (error) {
    console.error('Contract call failed:', error);
    throw error;
  }
};

// Send transaction
export const sendTransaction = async (tx) => {
  if (!isMetaMaskInstalled()) {
    throw new Error('MetaMask is not installed');
  }

  try {
    const txHash = await window.ethereum.request({
      method: 'eth_sendTransaction',
      params: [tx],
    });
    return txHash;
  } catch (error) {
    console.error('Transaction failed:', error);
    throw error;
  }
};

// Estimate gas
export const estimateGas = async (tx) => {
  if (!isMetaMaskInstalled()) {
    throw new Error('MetaMask is not installed');
  }

  try {
    const gas = await window.ethereum.request({
      method: 'eth_estimateGas',
      params: [tx],
    });
    return parseInt(gas, 16);
  } catch (error) {
    console.error('Gas estimation failed:', error);
    throw error;
  }
};

// Get gas price
export const getGasPrice = async () => {
  if (!isMetaMaskInstalled()) {
    throw new Error('MetaMask is not installed');
  }

  try {
    const gasPrice = await window.ethereum.request({
      method: 'eth_gasPrice',
    });
    return parseInt(gasPrice, 16);
  } catch (error) {
    console.error('Failed to get gas price:', error);
    throw error;
  }
};

// Validate Ethereum address
export const isValidAddress = (address) => {
  return /^0x[a-fA-F0-9]{40}$/.test(address);
};

// Get block number
export const getBlockNumber = async () => {
  if (!isMetaMaskInstalled()) {
    throw new Error('MetaMask is not installed');
  }

  try {
    const blockNumber = await window.ethereum.request({
      method: 'eth_blockNumber',
    });
    return parseInt(blockNumber, 16);
  } catch (error) {
    console.error('Failed to get block number:', error);
    throw error;
  }
};

export default {
  isMetaMaskInstalled,
  connectWallet,
  getCurrentAccount,
  switchNetwork,
  addNetwork,
  getBalance,
  signMessage,
  formatAddress,
  formatAmount,
  parseAmount,
  waitForTransaction,
  onAccountsChanged,
  onChainChanged,
  removeListeners,
  callContract,
  sendTransaction,
  estimateGas,
  getGasPrice,
  isValidAddress,
  getBlockNumber,
};
