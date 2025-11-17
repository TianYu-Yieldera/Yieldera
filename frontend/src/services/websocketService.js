/**
 * WebSocket Service
 * Real-time risk monitoring and alerts
 */

import { config } from '../config/env';

class WebSocketService {
  constructor() {
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 3000;
    this.listeners = new Map();
    this.connected = false;
    this.lastAlertId = 0;

    // Mock alert generation for demo
    this.mockAlertInterval = null;
  }

  /**
   * Connect to WebSocket server
   */
  connect(address) {
    // For now, use mock WebSocket since backend WS endpoint may not be ready
    // In production: const wsUrl = config.api.baseUrl.replace('http', 'ws') + '/ws/risk';

    console.log('[WebSocket] Connecting for address:', address);
    this.connected = true;
    this.startMockAlertGeneration(address);

    // Notify listeners of connection
    this.notifyListeners('connection', { connected: true });

    return Promise.resolve();
  }

  /**
   * Disconnect from WebSocket
   */
  disconnect() {
    console.log('[WebSocket] Disconnecting');

    if (this.mockAlertInterval) {
      clearInterval(this.mockAlertInterval);
      this.mockAlertInterval = null;
    }

    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }

    this.connected = false;
    this.notifyListeners('connection', { connected: false });
  }

  /**
   * Subscribe to events
   */
  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, []);
    }
    this.listeners.get(event).push(callback);
  }

  /**
   * Unsubscribe from events
   */
  off(event, callback) {
    if (!this.listeners.has(event)) return;

    const callbacks = this.listeners.get(event);
    const index = callbacks.indexOf(callback);
    if (index > -1) {
      callbacks.splice(index, 1);
    }
  }

  /**
   * Notify all listeners of an event
   */
  notifyListeners(event, data) {
    if (!this.listeners.has(event)) return;

    this.listeners.get(event).forEach(callback => {
      try {
        callback(data);
      } catch (error) {
        console.error(`[WebSocket] Error in ${event} listener:`, error);
      }
    });
  }

  /**
   * Start generating mock alerts for demo
   */
  startMockAlertGeneration(address) {
    // Clear any existing interval
    if (this.mockAlertInterval) {
      clearInterval(this.mockAlertInterval);
    }

    // Generate initial alert
    setTimeout(() => {
      this.generateMockAlert(address);
    }, 2000);

    // Generate periodic alerts (every 15-30 seconds)
    this.mockAlertInterval = setInterval(() => {
      if (Math.random() > 0.5) { // 50% chance
        this.generateMockAlert(address);
      }
    }, 20000 + Math.random() * 10000);
  }

  /**
   * Generate a mock alert
   */
  generateMockAlert(address) {
    const alertTypes = [
      {
        level: 'warning',
        type: 'price_movement',
        title: 'Price Volatility Detected',
        message: 'ETH price dropped 3.2% in the last hour',
        severity: 'medium',
        recommendation: 'Monitor your positions closely'
      },
      {
        level: 'danger',
        type: 'liquidation_risk',
        title: 'Health Factor Declining',
        message: 'Your Aave position health factor dropped to 1.65',
        severity: 'high',
        recommendation: 'Consider adding collateral to improve health factor'
      },
      {
        level: 'critical',
        type: 'liquidation_imminent',
        title: 'Liquidation Risk Critical',
        message: 'Position at severe risk - Health Factor: 1.15',
        severity: 'critical',
        recommendation: 'Add collateral immediately to avoid liquidation'
      },
      {
        level: 'info',
        type: 'opportunity',
        title: 'Yield Optimization Available',
        message: 'Higher APY detected: Compound V3 USDC at 8.5%',
        severity: 'low',
        recommendation: 'Consider rebalancing to optimize yield'
      },
      {
        level: 'warning',
        type: 'gas_spike',
        title: 'High Gas Prices',
        message: 'Network congestion detected - gas at 120 gwei',
        severity: 'medium',
        recommendation: 'Delay non-urgent transactions'
      },
      {
        level: 'danger',
        type: 'market_crash',
        title: 'Market Downturn Alert',
        message: 'Crypto market down 8% - potential liquidation risk',
        severity: 'high',
        recommendation: 'Review all leveraged positions'
      }
    ];

    const alert = alertTypes[Math.floor(Math.random() * alertTypes.length)];

    const alertData = {
      id: `alert-${++this.lastAlertId}`,
      level: alert.level,
      type: alert.type,
      title: alert.title,
      message: alert.message,
      severity: alert.severity,
      recommendation: alert.recommendation,
      timestamp: new Date().toISOString(),
      address: address,
      metadata: {
        position_id: Math.random() > 0.5 ? 'demo-position-1' : null,
        protocol: ['Aave V3', 'Compound V3', 'GMX V2'][Math.floor(Math.random() * 3)],
        asset: ['ETH', 'WBTC', 'USDC'][Math.floor(Math.random() * 3)],
        health_factor: alert.level === 'critical' ? 1.15 : alert.level === 'danger' ? 1.65 : 2.1
      }
    };

    console.log('[WebSocket] Mock alert generated:', alertData);
    this.notifyListeners('alert', alertData);

    // Also send metrics update
    this.sendMockMetrics();
  }

  /**
   * Send mock monitoring metrics
   */
  sendMockMetrics() {
    const metrics = {
      timestamp: new Date().toISOString(),
      positions_monitored: 3,
      alerts_triggered: this.lastAlertId,
      alerts_by_level: {
        info: Math.floor(this.lastAlertId * 0.3),
        warning: Math.floor(this.lastAlertId * 0.4),
        danger: Math.floor(this.lastAlertId * 0.2),
        critical: Math.floor(this.lastAlertId * 0.1)
      },
      avg_health_factor: 1.8 + Math.random() * 0.4,
      at_risk_positions: Math.random() > 0.7 ? 1 : 0,
      total_value_at_risk: Math.random() * 50000,
      system_health: 'operational',
      response_time_ms: 50 + Math.random() * 150
    };

    this.notifyListeners('metrics', metrics);
  }

  /**
   * Check if connected
   */
  isConnected() {
    return this.connected;
  }

  /**
   * Request browser notification permission
   */
  async requestNotificationPermission() {
    if (!('Notification' in window)) {
      console.warn('[WebSocket] Browser notifications not supported');
      return false;
    }

    if (Notification.permission === 'granted') {
      return true;
    }

    if (Notification.permission !== 'denied') {
      const permission = await Notification.requestPermission();
      return permission === 'granted';
    }

    return false;
  }

  /**
   * Show browser notification
   */
  showNotification(alert) {
    if (Notification.permission !== 'granted') {
      return;
    }

    const icon = alert.level === 'critical' ? '🚨' :
                 alert.level === 'danger' ? '⚠️' :
                 alert.level === 'warning' ? '⚡' : 'ℹ️';

    const notification = new Notification(`${icon} ${alert.title}`, {
      body: alert.message,
      icon: '/pointfi-logo-mark.svg',
      badge: '/pointfi-logo-mark.svg',
      tag: alert.id,
      requireInteraction: alert.level === 'critical',
      vibrate: alert.level === 'critical' ? [200, 100, 200] : undefined
    });

    notification.onclick = () => {
      window.focus();
      notification.close();
    };
  }
}

// Export singleton instance
export const websocketService = new WebSocketService();
export default websocketService;
