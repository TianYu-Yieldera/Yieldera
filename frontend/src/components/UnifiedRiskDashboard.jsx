/**
 * Unified Risk Dashboard
 * Displays institutional-grade risk metrics across DeFi, RWA, and AI insights
 */

import React, { useState, useEffect, useContext } from 'react';
import {
  AlertTriangle,
  TrendingUp,
  Shield,
  Activity,
  DollarSign,
  PieChart,
  Bell,
  RefreshCw,
  ChevronRight,
  AlertCircle,
} from 'lucide-react';
import { WalletContext } from '../web3/WalletContext';
import { Line, Doughnut } from 'react-chartjs-2';

const UnifiedRiskDashboard = () => {
  const { address } = useContext(WalletContext);
  const [portfolio, setPortfolio] = useState(null);
  const [alerts, setAlerts] = useState([]);
  const [riskHistory, setRiskHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [autoRefresh, setAutoRefresh] = useState(true);

  useEffect(() => {
    if (!address) {
      setLoading(false);
      return;
    }

    loadDashboardData();

    // Set up auto-refresh
    const interval = autoRefresh ? setInterval(loadDashboardData, 10000) : null;

    // Set up WebSocket for real-time updates
    const ws = new WebSocket(`ws://localhost:8080/api/defi/stream/${address}`);

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      handleRealtimeUpdate(data);
    };

    return () => {
      if (interval) clearInterval(interval);
      ws.close();
    };
  }, [address, autoRefresh]);

  const loadDashboardData = async () => {
    try {
      setLoading(true);

      // Fetch all data in parallel
      const [positionsRes, alertsRes, historyRes] = await Promise.all([
        fetch(`/api/defi/positions/${address}`),
        fetch(`/api/defi/alerts/${address}`),
        fetch(`/api/defi/history/${address}?period=7d`),
      ]);

      const positionsData = await positionsRes.json();
      const alertsData = await alertsRes.json();
      const historyData = await historyRes.json();

      if (positionsData.success) {
        setPortfolio(positionsData.data.current);
      }

      if (alertsData.success) {
        setAlerts(alertsData.data);
      }

      if (historyData.success) {
        setRiskHistory(historyData.data);
      }
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleRealtimeUpdate = (update) => {
    if (update.type === 'risk_update' && portfolio) {
      setPortfolio(prev => ({
        ...prev,
        overallRisk: update.data.overall_risk,
        averageHealthFactor: update.data.average_health_factor,
        totalValue: update.data.total_value,
      }));
    } else if (update.type === 'new_alerts') {
      setAlerts(prev => [...update.data, ...prev]);
    }
  };

  const getRiskColor = (score) => {
    if (score < 30) return '#10b981'; // green
    if (score < 60) return '#f59e0b'; // yellow
    if (score < 80) return '#f97316'; // orange
    return '#ef4444'; // red
  };

  const getRiskLevel = (score) => {
    if (score < 30) return 'Low Risk';
    if (score < 60) return 'Medium Risk';
    if (score < 80) return 'High Risk';
    return 'Critical Risk';
  };

  const formatCurrency = (value) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  if (!address) {
    return (
      <div className="bg-gray-800 rounded-lg p-8 text-center">
        <Shield className="h-16 w-16 text-gray-600 mx-auto mb-4" />
        <h2 className="text-xl font-semibold text-white mb-2">
          Connect Wallet to View Risk Dashboard
        </h2>
        <p className="text-gray-400">
          Access institutional-grade risk management for your DeFi positions
        </p>
      </div>
    );
  }

  if (loading && !portfolio) {
    return (
      <div className="bg-gray-800 rounded-lg p-8">
        <div className="flex items-center justify-center">
          <Activity className="h-8 w-8 animate-pulse text-blue-500 mr-3" />
          <span className="text-gray-400 text-lg">
            Analyzing your positions across all protocols...
          </span>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header with Key Metrics */}
      <div className="bg-gradient-to-r from-blue-900 to-purple-900 rounded-lg p-6">
        <div className="flex justify-between items-start mb-4">
          <div>
            <h1 className="text-3xl font-bold text-white flex items-center">
              <Shield className="h-10 w-10 mr-3" />
              Institutional Risk Dashboard
            </h1>
            <p className="text-blue-200 mt-2">
              AI-powered risk management across DeFi & RWA positions
            </p>
          </div>
          <button
            onClick={() => setAutoRefresh(!autoRefresh)}
            className={`px-4 py-2 rounded-lg flex items-center gap-2 transition-colors ${
              autoRefresh
                ? 'bg-green-600 hover:bg-green-700 text-white'
                : 'bg-gray-700 hover:bg-gray-600 text-gray-300'
            }`}
          >
            <RefreshCw className={`h-4 w-4 ${autoRefresh ? 'animate-spin' : ''}`} />
            {autoRefresh ? 'Live' : 'Paused'}
          </button>
        </div>

        {portfolio && (
          <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
            {/* Total Portfolio Value */}
            <div className="bg-black/30 rounded-lg p-4">
              <div className="text-sm text-gray-300 mb-1">Portfolio Value</div>
              <div className="text-2xl font-bold text-white">
                {formatCurrency(portfolio.totalValue)}
              </div>
              <div className="text-xs text-gray-400 mt-1">
                Across {portfolio.positions?.length || 0} positions
              </div>
            </div>

            {/* Overall Risk Score */}
            <div className="bg-black/30 rounded-lg p-4">
              <div className="text-sm text-gray-300 mb-1">Risk Score</div>
              <div
                className="text-2xl font-bold"
                style={{ color: getRiskColor(portfolio.overallRisk) }}
              >
                {portfolio.overallRisk}/100
              </div>
              <div className="text-xs text-gray-400 mt-1">
                {getRiskLevel(portfolio.overallRisk)}
              </div>
            </div>

            {/* Health Factor */}
            <div className="bg-black/30 rounded-lg p-4">
              <div className="text-sm text-gray-300 mb-1">Health Factor</div>
              <div
                className={`text-2xl font-bold ${
                  portfolio.averageHealthFactor > 1.5 ? 'text-green-400' : 'text-orange-400'
                }`}
              >
                {portfolio.averageHealthFactor?.toFixed(2)}
              </div>
              <div className="text-xs text-gray-400 mt-1">
                {portfolio.averageHealthFactor > 1.5 ? 'Safe' : 'Monitor Closely'}
              </div>
            </div>

            {/* Total Collateral */}
            <div className="bg-black/30 rounded-lg p-4">
              <div className="text-sm text-gray-300 mb-1">Collateral</div>
              <div className="text-2xl font-bold text-white">
                {formatCurrency(portfolio.totalCollateral)}
              </div>
              <div className="text-xs text-gray-400 mt-1">Securing positions</div>
            </div>

            {/* Total Debt */}
            <div className="bg-black/30 rounded-lg p-4">
              <div className="text-sm text-gray-300 mb-1">Debt</div>
              <div className="text-2xl font-bold text-white">
                {formatCurrency(portfolio.totalDebt)}
              </div>
              <div className="text-xs text-gray-400 mt-1">
                LTV: {((portfolio.totalDebt / portfolio.totalCollateral) * 100).toFixed(1)}%
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Critical Alerts Section */}
      {alerts.length > 0 && (
        <div className="bg-red-900/20 border border-red-600 rounded-lg p-6">
          <div className="flex items-start mb-4">
            <Bell className="h-6 w-6 text-red-500 mr-3 mt-0.5" />
            <div className="flex-1">
              <h2 className="text-lg font-semibold text-white mb-3">
                Active Risk Alerts ({alerts.length})
              </h2>
              <div className="space-y-3">
                {alerts.slice(0, 3).map((alert, index) => (
                  <div
                    key={index}
                    className="bg-gray-800/50 rounded-lg p-4 flex items-start justify-between"
                  >
                    <div className="flex items-start flex-1">
                      <AlertTriangle
                        className={`h-5 w-5 mr-3 mt-0.5 ${
                          alert.alert_type === 'critical' ? 'text-red-500' : 'text-yellow-500'
                        }`}
                      />
                      <div>
                        <div className="font-semibold text-white">
                          {alert.alert_type === 'critical'
                            ? 'Critical Liquidation Risk'
                            : 'Liquidation Warning'}
                        </div>
                        <div className="text-sm text-gray-300 mt-1">
                          Health Factor: {alert.health_factor?.toFixed(2)} • Risk Score:{' '}
                          {alert.risk_score}
                        </div>
                        {alert.predicted_liquidation_time && (
                          <div className="text-sm text-gray-400 mt-1">
                            Est. liquidation in:{' '}
                            {new Date(alert.predicted_liquidation_time).toLocaleString()}
                          </div>
                        )}
                        {alert.recommended_action && (
                          <div className="bg-gray-900/50 rounded p-2 mt-2">
                            <div className="text-xs text-gray-400">Recommended Action:</div>
                            <div className="text-sm text-white mt-1">
                              {alert.recommended_action}
                            </div>
                            {alert.required_collateral > 0 && (
                              <div className="text-sm text-blue-400 mt-1">
                                Add {formatCurrency(alert.required_collateral)} collateral
                              </div>
                            )}
                          </div>
                        )}
                      </div>
                    </div>
                    <button className="px-3 py-1 bg-blue-600 hover:bg-blue-700 rounded text-sm text-white transition-colors">
                      Take Action
                    </button>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Positions by Protocol */}
      {portfolio?.positions && portfolio.positions.length > 0 && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Protocol Distribution */}
          <div className="bg-gray-800 rounded-lg p-6">
            <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
              <PieChart className="h-5 w-5 mr-2 text-blue-500" />
              Protocol Distribution
            </h3>
            <div className="space-y-3">
              {Object.entries(
                portfolio.positions.reduce((acc, pos) => {
                  acc[pos.protocol] = (acc[pos.protocol] || 0) + pos.value;
                  return acc;
                }, {})
              ).map(([protocol, value]) => (
                <div key={protocol} className="flex items-center justify-between">
                  <div className="flex items-center">
                    <div
                      className="w-3 h-3 rounded-full mr-3"
                      style={{
                        backgroundColor:
                          protocol === 'aave'
                            ? '#B152A0'
                            : protocol === 'compound'
                            ? '#00D395'
                            : protocol === 'uniswap'
                            ? '#FF007A'
                            : '#4B5FFA',
                      }}
                    />
                    <span className="text-gray-300 capitalize">{protocol}</span>
                  </div>
                  <div className="text-right">
                    <div className="text-white font-semibold">{formatCurrency(value)}</div>
                    <div className="text-xs text-gray-400">
                      {((value / portfolio.totalValue) * 100).toFixed(1)}%
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* AI Recommendations */}
          <div className="bg-gray-800 rounded-lg p-6">
            <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
              <AlertCircle className="h-5 w-5 mr-2 text-purple-500" />
              AI Recommendations
            </h3>
            {portfolio.recommendations && portfolio.recommendations.length > 0 ? (
              <div className="space-y-3">
                {portfolio.recommendations.map((rec, index) => (
                  <div key={index} className="flex items-start">
                    <ChevronRight className="h-4 w-4 text-blue-500 mt-0.5 mr-2 flex-shrink-0" />
                    <p className="text-sm text-gray-300">{rec}</p>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-gray-400">
                Your portfolio is well-balanced. No immediate actions required.
              </p>
            )}
          </div>
        </div>
      )}

      {/* Risk History Chart */}
      {riskHistory.length > 0 && (
        <div className="bg-gray-800 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
            <TrendingUp className="h-5 w-5 mr-2 text-green-500" />
            7-Day Risk Trend
          </h3>
          <div className="h-64">
            <Line
              data={{
                labels: riskHistory.map(h =>
                  new Date(h.time).toLocaleDateString('en-US', {
                    month: 'short',
                    day: 'numeric',
                    hour: '2-digit',
                  })
                ),
                datasets: [
                  {
                    label: 'Risk Score',
                    data: riskHistory.map(h => h.risk),
                    borderColor: '#ef4444',
                    backgroundColor: 'rgba(239, 68, 68, 0.1)',
                    tension: 0.4,
                  },
                  {
                    label: 'Portfolio Value ($K)',
                    data: riskHistory.map(h => h.value / 1000),
                    borderColor: '#10b981',
                    backgroundColor: 'rgba(16, 185, 129, 0.1)',
                    tension: 0.4,
                    yAxisID: 'y1',
                  },
                ],
              }}
              options={{
                responsive: true,
                maintainAspectRatio: false,
                interaction: {
                  mode: 'index',
                  intersect: false,
                },
                scales: {
                  y: {
                    type: 'linear',
                    display: true,
                    position: 'left',
                    title: {
                      display: true,
                      text: 'Risk Score',
                      color: '#ef4444',
                    },
                  },
                  y1: {
                    type: 'linear',
                    display: true,
                    position: 'right',
                    title: {
                      display: true,
                      text: 'Value ($K)',
                      color: '#10b981',
                    },
                    grid: {
                      drawOnChartArea: false,
                    },
                  },
                },
              }}
            />
          </div>
        </div>
      )}

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <button className="bg-blue-600 hover:bg-blue-700 rounded-lg p-4 text-white font-semibold transition-colors">
          Generate Risk Report
        </button>
        <button className="bg-purple-600 hover:bg-purple-700 rounded-lg p-4 text-white font-semibold transition-colors">
          Configure Auto-Hedging
        </button>
        <button className="bg-green-600 hover:bg-green-700 rounded-lg p-4 text-white font-semibold transition-colors">
          Export Portfolio Data
        </button>
      </div>
    </div>
  );
};

export default UnifiedRiskDashboard;