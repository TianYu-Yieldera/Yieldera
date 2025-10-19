// Custom hook for API integration with Demo Mode support

import { useState, useEffect, useCallback } from 'react';
import { useDemoMode } from '../web3/DemoModeContext';
import { useWallet } from '../web3/WalletContext';

/**
 * Custom hook that handles API calls with automatic demo/real mode switching
 * @param {Function} apiCall - The API function to call
 * @param {Array} dependencies - Dependencies for re-fetching
 * @param {Object} options - Additional options
 * @returns {Object} - { data, loading, error, refetch }
 */
export function useApiWithDemo(apiCall, dependencies = [], options = {}) {
  const { demoMode } = useDemoMode();
  const { address } = useWallet();
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchData = useCallback(async () => {
    if (!address && !demoMode && options.requireAuth !== false) {
      setLoading(false);
      setError('Please connect wallet or enable demo mode');
      return;
    }

    try {
      setLoading(true);
      setError(null);

      // Call the API function
      const result = await apiCall();
      setData(result);
    } catch (err) {
      console.error('API call failed:', err);
      setError(err.message || 'Failed to fetch data');
      setData(null);
    } finally {
      setLoading(false);
    }
  }, [address, demoMode, ...dependencies]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    data,
    loading,
    error,
    refetch: fetchData,
  };
}

/**
 * Custom hook for handling API mutations (POST, PUT, DELETE)
 * @param {Function} apiCall - The API function to call
 * @returns {Object} - { mutate, data, loading, error }
 */
export function useApiMutation(apiCall) {
  const { demoMode } = useDemoMode();
  const { address } = useWallet();
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const mutate = useCallback(async (...args) => {
    if (!address && !demoMode) {
      setError('Please connect wallet or enable demo mode');
      return { success: false, error: 'Not authenticated' };
    }

    try {
      setLoading(true);
      setError(null);

      // Call the API function with arguments
      const result = await apiCall(...args);
      setData(result);

      return { success: true, data: result };
    } catch (err) {
      console.error('API mutation failed:', err);
      const errorMessage = err.message || 'Operation failed';
      setError(errorMessage);

      return { success: false, error: errorMessage };
    } finally {
      setLoading(false);
    }
  }, [address, demoMode, apiCall]);

  return {
    mutate,
    data,
    loading,
    error,
    reset: () => {
      setData(null);
      setError(null);
    },
  };
}

/**
 * Custom hook for polling data at regular intervals
 * @param {Function} apiCall - The API function to call
 * @param {number} interval - Polling interval in milliseconds
 * @param {boolean} enabled - Whether polling is enabled
 * @returns {Object} - { data, loading, error, refetch }
 */
export function useApiPolling(apiCall, interval = 60000, enabled = true) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchData = useCallback(async () => {
    try {
      const result = await apiCall();
      setData(result);
      setError(null);
    } catch (err) {
      console.error('Polling failed:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [apiCall]);

  useEffect(() => {
    if (!enabled) return;

    // Initial fetch
    fetchData();

    // Set up polling
    const intervalId = setInterval(fetchData, interval);

    return () => clearInterval(intervalId);
  }, [fetchData, interval, enabled]);

  return { data, loading, error, refetch: fetchData };
}

/**
 * Custom hook for WebSocket connections with demo fallback
 * @param {string} url - WebSocket URL
 * @param {Object} options - Connection options
 * @returns {Object} - { connected, data, send, close }
 */
export function useWebSocketWithDemo(url, options = {}) {
  const { demoMode } = useDemoMode();
  const [connected, setConnected] = useState(false);
  const [data, setData] = useState(null);
  const [ws, setWs] = useState(null);

  useEffect(() => {
    if (demoMode) {
      // In demo mode, simulate WebSocket with mock data
      setConnected(true);

      // Simulate periodic data updates
      const interval = setInterval(() => {
        setData({
          type: 'price_update',
          data: {
            timestamp: new Date().toISOString(),
            prices: {
              'bAAPL': 178.52 + (Math.random() - 0.5) * 2,
              'bTSLA': 242.18 + (Math.random() - 0.5) * 5,
              'PAXG': 2042.50 + (Math.random() - 0.5) * 10,
            },
          },
        });
      }, 5000);

      return () => clearInterval(interval);
    }

    // Real WebSocket connection
    try {
      const websocket = new WebSocket(url);

      websocket.onopen = () => {
        setConnected(true);
        if (options.onOpen) options.onOpen();
      };

      websocket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        setData(message);
        if (options.onMessage) options.onMessage(message);
      };

      websocket.onclose = () => {
        setConnected(false);
        if (options.onClose) options.onClose();
      };

      websocket.onerror = (error) => {
        console.error('WebSocket error:', error);
        if (options.onError) options.onError(error);
      };

      setWs(websocket);

      return () => {
        websocket.close();
      };
    } catch (err) {
      console.error('Failed to connect WebSocket:', err);
      setConnected(false);
    }
  }, [url, demoMode]);

  const send = useCallback((message) => {
    if (demoMode) {
      console.log('Demo mode: Simulating message send:', message);
      return;
    }

    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(message));
    }
  }, [ws, demoMode]);

  const close = useCallback(() => {
    if (ws) {
      ws.close();
    }
  }, [ws]);

  return { connected, data, send, close };
}

export default {
  useApiWithDemo,
  useApiMutation,
  useApiPolling,
  useWebSocketWithDemo,
};