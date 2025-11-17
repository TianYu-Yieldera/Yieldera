/**
 * Realtime Monitor Panel
 * 24/7 WebSocket-powered risk monitoring with live alerts
 */

import React, { useState, useEffect, useRef } from 'react';
import { Bell, BellOff, Activity, AlertTriangle, Info, CheckCircle, TrendingDown, Zap, Wifi, WifiOff, Clock, Shield } from 'lucide-react';
import { useWallet } from '../web3/WalletContext';
import websocketService from '../services/websocketService';

const RealtimeMonitorPanel = () => {
  const { address } = useWallet();
  const [connected, setConnected] = useState(false);
  const [alerts, setAlerts] = useState([]);
  const [metrics, setMetrics] = useState(null);
  const [notificationsEnabled, setNotificationsEnabled] = useState(false);
  const [loading, setLoading] = useState(false);
  const alertsEndRef = useRef(null);

  useEffect(() => {
    if (!address) return;

    // Connect to WebSocket
    startMonitoring();

    return () => {
      stopMonitoring();
    };
  }, [address]);

  useEffect(() => {
    // Auto-scroll to latest alert
    if (alertsEndRef.current) {
      alertsEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [alerts]);

  const startMonitoring = async () => {
    setLoading(true);

    try {
      // Connect WebSocket
      await websocketService.connect(address);

      // Subscribe to events
      websocketService.on('connection', handleConnection);
      websocketService.on('alert', handleAlert);
      websocketService.on('metrics', handleMetrics);

      setConnected(true);
    } catch (error) {
      console.error('[Monitor] Failed to start monitoring:', error);
    } finally {
      setLoading(false);
    }
  };

  const stopMonitoring = () => {
    websocketService.disconnect();
    websocketService.off('connection', handleConnection);
    websocketService.off('alert', handleAlert);
    websocketService.off('metrics', handleMetrics);
    setConnected(false);
  };

  const handleConnection = (data) => {
    setConnected(data.connected);
  };

  const handleAlert = (alert) => {
    console.log('[Monitor] New alert received:', alert);

    // Add to alerts list
    setAlerts(prev => [alert, ...prev].slice(0, 50)); // Keep last 50 alerts

    // Show browser notification if enabled
    if (notificationsEnabled && (alert.level === 'critical' || alert.level === 'danger')) {
      websocketService.showNotification(alert);
    }

    // Play sound for critical alerts
    if (alert.level === 'critical') {
      playAlertSound();
    }
  };

  const handleMetrics = (data) => {
    setMetrics(data);
  };

  const playAlertSound = () => {
    // Simple beep sound using Web Audio API
    try {
      const audioContext = new (window.AudioContext || window.webkitAudioContext)();
      const oscillator = audioContext.createOscillator();
      const gainNode = audioContext.createGain();

      oscillator.connect(gainNode);
      gainNode.connect(audioContext.destination);

      oscillator.frequency.value = 800;
      oscillator.type = 'sine';

      gainNode.gain.setValueAtTime(0.3, audioContext.currentTime);
      gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.5);

      oscillator.start(audioContext.currentTime);
      oscillator.stop(audioContext.currentTime + 0.5);
    } catch (error) {
      console.warn('[Monitor] Could not play alert sound:', error);
    }
  };

  const toggleNotifications = async () => {
    if (!notificationsEnabled) {
      const granted = await websocketService.requestNotificationPermission();
      if (granted) {
        setNotificationsEnabled(true);
      } else {
        alert('Please enable notifications in your browser settings to receive alerts.');
      }
    } else {
      setNotificationsEnabled(false);
    }
  };

  const clearAlerts = () => {
    setAlerts([]);
  };

  const getAlertIcon = (level) => {
    switch (level) {
      case 'critical': return <AlertTriangle size={20} style={{ color: 'rgb(239, 68, 68)' }} />;
      case 'danger': return <TrendingDown size={20} style={{ color: 'rgb(249, 115, 22)' }} />;
      case 'warning': return <Zap size={20} style={{ color: 'rgb(245, 158, 11)' }} />;
      default: return <Info size={20} style={{ color: 'rgb(34, 211, 238)' }} />;
    }
  };

  const getAlertStyle = (level) => {
    switch (level) {
      case 'critical':
        return {
          bg: 'rgba(239, 68, 68, 0.1)',
          border: 'rgba(239, 68, 68, 0.4)',
          color: 'rgb(239, 68, 68)'
        };
      case 'danger':
        return {
          bg: 'rgba(249, 115, 22, 0.1)',
          border: 'rgba(249, 115, 22, 0.4)',
          color: 'rgb(249, 115, 22)'
        };
      case 'warning':
        return {
          bg: 'rgba(245, 158, 11, 0.1)',
          border: 'rgba(245, 158, 11, 0.4)',
          color: 'rgb(245, 158, 11)'
        };
      default:
        return {
          bg: 'rgba(34, 211, 238, 0.1)',
          border: 'rgba(34, 211, 238, 0.4)',
          color: 'rgb(34, 211, 238)'
        };
    }
  };

  const formatTimestamp = (timestamp) => {
    const date = new Date(timestamp);
    const now = new Date();
    const diffSeconds = Math.floor((now - date) / 1000);

    if (diffSeconds < 60) return `${diffSeconds}s ago`;
    if (diffSeconds < 3600) return `${Math.floor(diffSeconds / 60)}m ago`;
    if (diffSeconds < 86400) return `${Math.floor(diffSeconds / 3600)}h ago`;
    return date.toLocaleString();
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
      {/* Header */}
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', flexWrap: 'wrap', gap: 16 }}>
        <div>
          <h3 style={{ fontSize: 18, fontWeight: 700, color: 'white', margin: '0 0 8px 0', display: 'flex', alignItems: 'center', gap: 12 }}>
            <Activity size={24} style={{ color: connected ? 'rgb(34, 197, 94)' : 'rgb(148, 163, 184)' }} />
            Real-time Risk Monitoring
          </h3>
          <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.7)', margin: 0, display: 'flex', alignItems: 'center', gap: 8 }}>
            {connected ? (
              <>
                <Wifi size={14} style={{ color: 'rgb(34, 197, 94)' }} />
                <span>Connected - monitoring {metrics?.positions_monitored || 0} positions</span>
              </>
            ) : (
              <>
                <WifiOff size={14} style={{ color: 'rgb(148, 163, 184)' }} />
                <span>Disconnected</span>
              </>
            )}
          </p>
        </div>

        <div style={{ display: 'flex', gap: 12 }}>
          <button
            onClick={toggleNotifications}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 8,
              padding: '10px 20px',
              background: notificationsEnabled ? 'rgba(34, 197, 94, 0.2)' : 'rgba(148, 163, 184, 0.2)',
              border: notificationsEnabled ? '1px solid rgba(34, 197, 94, 0.4)' : '1px solid rgba(148, 163, 184, 0.4)',
              borderRadius: 8,
              color: notificationsEnabled ? 'rgb(34, 197, 94)' : 'rgb(148, 163, 184)',
              fontSize: 14,
              fontWeight: 600,
              cursor: 'pointer',
              transition: 'all 0.3s'
            }}
          >
            {notificationsEnabled ? <Bell size={16} /> : <BellOff size={16} />}
            {notificationsEnabled ? 'Notifications On' : 'Enable Notifications'}
          </button>

          {alerts.length > 0 && (
            <button
              onClick={clearAlerts}
              style={{
                padding: '10px 20px',
                background: 'rgba(148, 163, 184, 0.2)',
                border: '1px solid rgba(148, 163, 184, 0.4)',
                borderRadius: 8,
                color: 'rgb(148, 163, 184)',
                fontSize: 14,
                fontWeight: 600,
                cursor: 'pointer',
                transition: 'all 0.3s'
              }}
            >
              Clear Alerts
            </button>
          )}
        </div>
      </div>

      {/* Monitoring Metrics */}
      {metrics && (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: 16 }}>
          <div style={{
            background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
            borderRadius: 12,
            padding: 20,
            border: '1px solid rgba(34, 211, 238, 0.2)'
          }}>
            <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 500 }}>
              Total Alerts
            </div>
            <div style={{ fontSize: 32, fontWeight: 800, color: 'white' }}>
              {metrics.alerts_triggered}
            </div>
          </div>

          <div style={{
            background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
            borderRadius: 12,
            padding: 20,
            border: '1px solid rgba(34, 211, 238, 0.2)'
          }}>
            <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 500 }}>
              Avg Health Factor
            </div>
            <div style={{ fontSize: 32, fontWeight: 800, color: metrics.avg_health_factor < 1.5 ? 'rgb(239, 68, 68)' : 'rgb(34, 197, 94)' }}>
              {metrics.avg_health_factor.toFixed(2)}
            </div>
          </div>

          <div style={{
            background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
            borderRadius: 12,
            padding: 20,
            border: '1px solid rgba(34, 211, 238, 0.2)'
          }}>
            <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 500 }}>
              At-Risk Positions
            </div>
            <div style={{ fontSize: 32, fontWeight: 800, color: metrics.at_risk_positions > 0 ? 'rgb(239, 68, 68)' : 'rgb(34, 197, 94)' }}>
              {metrics.at_risk_positions}
            </div>
          </div>

          <div style={{
            background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
            borderRadius: 12,
            padding: 20,
            border: '1px solid rgba(34, 211, 238, 0.2)'
          }}>
            <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 500 }}>
              Response Time
            </div>
            <div style={{ fontSize: 32, fontWeight: 800, color: 'white' }}>
              {metrics.response_time_ms.toFixed(0)}ms
            </div>
          </div>
        </div>
      )}

      {/* Alert Feed */}
      <div style={{
        background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
        borderRadius: 16,
        padding: 24,
        border: '1px solid rgba(34, 211, 238, 0.2)',
        minHeight: 400,
        maxHeight: 600,
        overflow: 'hidden',
        display: 'flex',
        flexDirection: 'column'
      }}>
        <h4 style={{ fontSize: 16, fontWeight: 700, color: 'white', marginBottom: 20, display: 'flex', alignItems: 'center', gap: 8 }}>
          <Bell size={20} style={{ color: 'rgb(34, 211, 238)' }} />
          Live Alert Feed
          {alerts.length > 0 && (
            <span style={{
              fontSize: 11,
              padding: '4px 10px',
              background: 'rgba(34, 211, 238, 0.2)',
              borderRadius: 6,
              fontWeight: 700,
              color: 'rgb(34, 211, 238)'
            }}>
              {alerts.length}
            </span>
          )}
        </h4>

        <div style={{ flex: 1, overflowY: 'auto', display: 'flex', flexDirection: 'column', gap: 12 }}>
          {loading ? (
            <div style={{ textAlign: 'center', padding: 40 }}>
              <div style={{
                width: 40,
                height: 40,
                margin: '0 auto 16px',
                borderRadius: '50%',
                border: '3px solid rgba(34, 211, 238, 0.3)',
                borderTopColor: 'rgb(34, 211, 238)',
                animation: 'spin 1s linear infinite'
              }} />
              <style>{`
                @keyframes spin {
                  to { transform: rotate(360deg); }
                }
              `}</style>
              <p style={{ color: 'rgba(203, 213, 225, 0.7)', fontSize: 14 }}>
                Connecting to monitoring service...
              </p>
            </div>
          ) : alerts.length === 0 ? (
            <div style={{ textAlign: 'center', padding: 40 }}>
              <Shield size={48} style={{ color: 'rgba(203, 213, 225, 0.3)', marginBottom: 16 }} />
              <p style={{ fontSize: 16, color: 'rgba(203, 213, 225, 0.8)', margin: '0 0 8px 0', fontWeight: 600 }}>
                No Alerts
              </p>
              <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.6)', margin: 0 }}>
                All positions are healthy. Monitoring is active.
              </p>
            </div>
          ) : (
            <>
              {alerts.map((alert, index) => {
                const style = getAlertStyle(alert.level);
                return (
                  <div
                    key={alert.id}
                    style={{
                      padding: 16,
                      background: style.bg,
                      border: `1px solid ${style.border}`,
                      borderRadius: 12,
                      animation: index === 0 ? 'slideIn 0.3s ease-out' : 'none'
                    }}
                  >
                    <style>{`
                      @keyframes slideIn {
                        from {
                          opacity: 0;
                          transform: translateY(-10px);
                        }
                        to {
                          opacity: 1;
                          transform: translateY(0);
                        }
                      }
                    `}</style>

                    <div style={{ display: 'flex', alignItems: 'flex-start', gap: 12 }}>
                      {getAlertIcon(alert.level)}
                      <div style={{ flex: 1 }}>
                        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 8 }}>
                          <div style={{ fontSize: 14, fontWeight: 700, color: 'white' }}>
                            {alert.title}
                          </div>
                          <div style={{ display: 'flex', alignItems: 'center', gap: 4, fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                            <Clock size={12} />
                            {formatTimestamp(alert.timestamp)}
                          </div>
                        </div>

                        <p style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.9)', margin: '0 0 12px 0', lineHeight: 1.5 }}>
                          {alert.message}
                        </p>

                        {alert.recommendation && (
                          <div style={{
                            padding: 10,
                            background: 'rgba(15, 23, 42, 0.5)',
                            borderRadius: 8,
                            fontSize: 12,
                            color: 'rgba(203, 213, 225, 0.9)',
                            display: 'flex',
                            alignItems: 'flex-start',
                            gap: 8
                          }}>
                            <CheckCircle size={14} style={{ color: style.color, marginTop: 2, flexShrink: 0 }} />
                            <span><strong>Recommendation:</strong> {alert.recommendation}</span>
                          </div>
                        )}

                        {alert.metadata && (
                          <div style={{ marginTop: 12, display: 'flex', gap: 12, flexWrap: 'wrap' }}>
                            {alert.metadata.protocol && (
                              <span style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                                Protocol: <strong style={{ color: 'white' }}>{alert.metadata.protocol}</strong>
                              </span>
                            )}
                            {alert.metadata.asset && (
                              <span style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                                Asset: <strong style={{ color: 'white' }}>{alert.metadata.asset}</strong>
                              </span>
                            )}
                            {alert.metadata.health_factor && (
                              <span style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                                HF: <strong style={{ color: style.color }}>{alert.metadata.health_factor.toFixed(2)}</strong>
                              </span>
                            )}
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                );
              })}
              <div ref={alertsEndRef} />
            </>
          )}
        </div>
      </div>

      {/* Info Banner */}
      <div style={{
        padding: 16,
        background: 'rgba(34, 211, 238, 0.1)',
        border: '1px solid rgba(34, 211, 238, 0.3)',
        borderRadius: 10,
        fontSize: 13,
        color: 'rgba(203, 213, 225, 0.9)'
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
          <Activity size={16} style={{ color: 'rgb(34, 211, 238)' }} />
          <strong style={{ color: 'white' }}>24/7 Monitoring System</strong>
        </div>
        <p style={{ margin: 0, lineHeight: 1.6 }}>
          实时监控系统通过WebSocket连接持续跟踪您的仓位风险。
          当检测到价格波动、健康因子下降或其他风险信号时，系统会立即发送警报。
          启用浏览器通知可在后台接收关键警报。
        </p>
      </div>
    </div>
  );
};

export default RealtimeMonitorPanel;
