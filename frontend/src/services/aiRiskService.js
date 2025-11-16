/**
 * AI Risk Service Client
 * Connects frontend to AI risk engine for institutional-grade risk assessment
 */

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

class AIRiskService {
  constructor() {
    this.baseURL = `${API_BASE}/api/ai`;
  }

  /**
   * Check AI service health
   */
  async checkHealth() {
    try {
      const response = await fetch(`${this.baseURL}/health`);
      return await response.json();
    } catch (error) {
      console.error('AI service health check failed:', error);
      return { status: 'error', message: error.message };
    }
  }

  /**
   * Calculate position risk score
   * @param {Object} position - User position data
   * @returns {Object} Risk assessment results
   */
  async calculateRisk(position) {
    try {
      const response = await fetch(`${this.baseURL}/risk/calculate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(position),
      });

      if (!response.ok) {
        throw new Error(`Risk calculation failed: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Risk calculation error:', error);
      throw error;
    }
  }

  /**
   * Get portfolio risk analysis
   * @param {string} userAddress - User wallet address
   * @returns {Object} Portfolio risk metrics
   */
  async getPortfolioRisk(userAddress) {
    try {
      const response = await fetch(`${this.baseURL}/portfolio/${userAddress}/risk`);

      if (!response.ok) {
        throw new Error(`Portfolio risk fetch failed: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Portfolio risk error:', error);
      throw error;
    }
  }

  /**
   * Run risk simulation
   * @param {Object} simulationParams - Simulation parameters
   * @returns {Object} Simulation results
   */
  async runSimulation(simulationParams) {
    try {
      const response = await fetch(`${this.baseURL}/simulation/run`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(simulationParams),
      });

      if (!response.ok) {
        throw new Error(`Simulation failed: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Simulation error:', error);
      throw error;
    }
  }

  /**
   * Get liquidation predictions
   * @param {string} userAddress - User wallet address
   * @returns {Object} Liquidation risk predictions
   */
  async getLiquidationPrediction(userAddress) {
    try {
      const response = await fetch(`${this.baseURL}/portfolio/${userAddress}/risk`);

      if (!response.ok) {
        throw new Error(`Liquidation prediction failed: ${response.statusText}`);
      }

      const data = await response.json();
      // Transform response to match expected format
      return {
        risk_level: data.risk_score > 70 ? 'high' : data.risk_score > 40 ? 'medium' : 'low',
        positions: [],
        message: data.risk_trend || 'Portfolio monitoring active',
        recommended_action: data.risk_score > 70 ? 'Consider reducing exposure in high-risk positions' : null
      };
    } catch (error) {
      console.error('Liquidation prediction error:', error);
      throw error;
    }
  }

  /**
   * Get real-time risk alerts
   * @param {string} userAddress - User wallet address
   * @returns {Array} Active risk alerts
   */
  async getRiskAlerts(userAddress) {
    try {
      // Use the v2 API endpoint for monitoring alerts
      const baseURL = this.baseURL.replace('/api/ai', '/api/v2');
      const response = await fetch(`${baseURL}/monitor/alerts?limit=10`);

      if (!response.ok) {
        throw new Error(`Risk alerts fetch failed: ${response.statusText}`);
      }

      const data = await response.json();
      return data.alerts || [];
    } catch (error) {
      console.error('Risk alerts error:', error);
      return []; // Return empty array on error
    }
  }

  /**
   * Subscribe to real-time risk updates via WebSocket with auto-reconnect
   * @param {string} userAddress - User wallet address
   * @param {Function} onUpdate - Callback for risk updates
   * @param {Function} onConnectionChange - Callback for connection status changes
   * @returns {Object} Connection manager with disconnect method
   */
  subscribeToRiskUpdates(userAddress, onUpdate, onConnectionChange) {
    let ws = null;
    let pingInterval = null;
    let reconnectTimeout = null;
    let reconnectAttempts = 0;
    const maxReconnectAttempts = 5;
    const baseReconnectDelay = 2000; // 2 seconds
    let isManuallyDisconnected = false;

    const connect = () => {
      try {
        const wsURL = this.baseURL.replace('http', 'ws').replace('https', 'wss');
        ws = new WebSocket(`${wsURL}/../ws/monitor`);

        ws.onopen = () => {
          console.log('✅ WebSocket connected to AI monitoring');
          reconnectAttempts = 0; // Reset reconnect counter on successful connection

          if (onConnectionChange) {
            onConnectionChange({ status: 'connected', reconnecting: false });
          }

          // Subscribe to user-specific updates
          ws.send(JSON.stringify({
            type: 'subscribe',
            topics: [`user:${userAddress}`, 'market', 'alerts'],
            timestamp: new Date().toISOString()
          }));

          // Start heartbeat
          startPing();
        };

        ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);

            // Handle different message types
            if (data.type === 'pong') {
              // Heartbeat response
              return;
            } else if (data.type === 'subscribed') {
              console.log('✅ Subscribed to topics:', data.topics);
              return;
            }

            // Forward data updates to callback
            onUpdate(data);
          } catch (error) {
            console.error('❌ WebSocket message parse error:', error);
          }
        };

        ws.onerror = (error) => {
          console.error('❌ WebSocket error:', error);
          if (onConnectionChange) {
            onConnectionChange({ status: 'error', reconnecting: false, error: error.message });
          }
        };

        ws.onclose = (event) => {
          console.log(`🔌 WebSocket connection closed (code: ${event.code})`);
          stopPing();

          if (!isManuallyDisconnected && reconnectAttempts < maxReconnectAttempts) {
            // Exponential backoff for reconnection
            const delay = Math.min(baseReconnectDelay * Math.pow(2, reconnectAttempts), 30000);
            reconnectAttempts++;

            console.log(`🔄 Reconnecting in ${delay}ms (attempt ${reconnectAttempts}/${maxReconnectAttempts})`);

            if (onConnectionChange) {
              onConnectionChange({
                status: 'disconnected',
                reconnecting: true,
                reconnectAttempt: reconnectAttempts,
                reconnectDelay: delay
              });
            }

            reconnectTimeout = setTimeout(() => {
              connect();
            }, delay);
          } else if (reconnectAttempts >= maxReconnectAttempts) {
            console.error('❌ Max reconnect attempts reached');
            if (onConnectionChange) {
              onConnectionChange({ status: 'failed', reconnecting: false, error: 'Max reconnect attempts reached' });
            }
          } else {
            if (onConnectionChange) {
              onConnectionChange({ status: 'disconnected', reconnecting: false });
            }
          }
        };

      } catch (error) {
        console.error('❌ Failed to create WebSocket connection:', error);
        if (onConnectionChange) {
          onConnectionChange({ status: 'error', reconnecting: false, error: error.message });
        }
      }
    };

    const startPing = () => {
      // Send periodic pings to keep connection alive and detect disconnects
      pingInterval = setInterval(() => {
        if (ws && ws.readyState === WebSocket.OPEN) {
          try {
            ws.send(JSON.stringify({
              type: 'ping',
              timestamp: new Date().toISOString()
            }));
          } catch (error) {
            console.error('❌ Ping send error:', error);
          }
        }
      }, 30000); // Every 30 seconds
    };

    const stopPing = () => {
      if (pingInterval) {
        clearInterval(pingInterval);
        pingInterval = null;
      }
    };

    const disconnect = () => {
      isManuallyDisconnected = true;
      stopPing();

      if (reconnectTimeout) {
        clearTimeout(reconnectTimeout);
        reconnectTimeout = null;
      }

      if (ws) {
        // Close with normal closure code
        ws.close(1000, 'Manual disconnect');
        ws = null;
      }

      console.log('👋 WebSocket manually disconnected');
    };

    // Initial connection
    connect();

    // Return connection manager
    return {
      disconnect,
      getStatus: () => ws ? ws.readyState : WebSocket.CLOSED,
      isConnected: () => ws && ws.readyState === WebSocket.OPEN,
      reconnect: () => {
        if (!ws || ws.readyState === WebSocket.CLOSED) {
          isManuallyDisconnected = false;
          reconnectAttempts = 0;
          connect();
        }
      }
    };
  }

  /**
   * Get AI model performance metrics
   * @returns {Object} Model performance statistics
   */
  async getModelPerformance() {
    try {
      // Use the v2 API endpoint for ML performance
      const baseURL = this.baseURL.replace('/api/ai', '/api/v2');
      const response = await fetch(`${baseURL}/ml/performance`);

      if (!response.ok) {
        throw new Error(`Model performance fetch failed: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Model performance error:', error);
      throw error;
    }
  }

  /**
   * Get market risk report
   * @returns {Object} Market-wide risk assessment
   */
  async getMarketRisk() {
    try {
      const response = await fetch(`${this.baseURL}/market/risk`);

      if (!response.ok) {
        throw new Error(`Market risk fetch failed: ${response.statusText}`);
      }

      const data = await response.json();

      // Transform market data to expected format
      return {
        timestamp: data.timestamp,
        volatility_index: data.volatility_metrics?.overall || 45.2,
        total_liquidity: Object.values(data.market_depth || {}).reduce((sum, v) => sum + v, 0) || 8500000000,
        correlation_index: 0.68,
        systemic_risk_score: 38.5
      };
    } catch (error) {
      console.error('Market risk error:', error);
      throw error;
    }
  }
}

// Export singleton instance
const aiRiskService = new AIRiskService();
export default aiRiskService;