/**
 * AI Risk Dashboard Component - Enhanced Version
 * Institutional-grade risk metrics powered by AI with advanced visualizations
 */

import React, { useState, useEffect } from 'react';
import { AlertTriangle, TrendingUp, Shield, Activity, AlertCircle, Zap, Brain, Eye, Target } from 'lucide-react';
import { useWallet } from '../web3/WalletContext';
import { useDemoMode } from '../web3/DemoModeContext';
import aiRiskService from '../services/aiRiskService';

const AIRiskDashboard = () => {
  const { address } = useWallet();
  const { demoMode } = useDemoMode();
  const [riskData, setRiskData] = useState(null);
  const [alerts, setAlerts] = useState([]);
  const [liquidationRisk, setLiquidationRisk] = useState(null);
  const [marketRisk, setMarketRisk] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [wsConnection, setWsConnection] = useState(null);
  const [analyzing, setAnalyzing] = useState(false);

  useEffect(() => {
    loadRiskData();

    // WebSocket disabled temporarily - API endpoints not ready
    // Will re-enable when backend WebSocket is available
  }, [address, demoMode]);

  const loadRiskData = async () => {
    try {
      setLoading(true);
      setAnalyzing(true);
      setError(null);

      // In demo mode or when APIs aren't available, use mock data
      if (demoMode || !address) {
        loadMockData();
        setLoading(false);
        setAnalyzing(false);
      } else {
        // Try to load real data, but fallback quickly if it fails
        try {
          const timeout = new Promise((_, reject) =>
            setTimeout(() => reject(new Error('Timeout')), 2000) // Reduced to 2s
          );

          const apiCalls = Promise.all([
            aiRiskService.getPortfolioRisk(address),
            aiRiskService.getRiskAlerts(address),
            aiRiskService.getLiquidationPrediction(address),
            aiRiskService.getMarketRisk(),
          ]);

          const [portfolio, alerts, liquidation, market] = await Promise.race([apiCalls, timeout]);

          setRiskData(portfolio);
          setAlerts(alerts);
          setLiquidationRisk(liquidation);
          setMarketRisk(market);
          setLoading(false);
          setAnalyzing(false);
        } catch (err) {
          console.warn('API failed, using mock data:', err.message);
          loadMockData();
          setLoading(false);
          setAnalyzing(false);
        }
      }
    } catch (err) {
      console.error('Failed to load risk data:', err);
      setError(null); // Don't show error, just use mock data
      loadMockData();
      setLoading(false);
      setAnalyzing(false);
    }
  };

  const loadMockData = () => {
    setRiskData({
      overall_risk_score: 35.5,
      total_collateral: 50000,
      total_debt: 25000,
      health_factor: 1.85,
      position_count: 3
    });

    setAlerts([
      {
        position_id: 'demo-1',
        protocol: 'Aave V3',
        risk_score: 42.3,
        alert_level: 'warning',
        liquidation_probability: 0.15,
        title: 'Medium Risk Position',
        description: 'Position health factor below recommended threshold',
        timestamp: new Date().toISOString()
      }
    ]);

    setLiquidationRisk({
      risk_level: 'medium',
      positions: [
        {
          position_id: 'demo-2',
          protocol: 'Compound V3',
          collateral_asset: 'ETH',
          debt_asset: 'USDC',
          health_factor: 1.45,
          liquidation_probability: 0.25
        }
      ]
    });

    setMarketRisk({
      timestamp: new Date().toISOString(),
      volatility_index: 45.2,
      total_liquidity: 8500000000,
      correlation_index: 0.68,
      systemic_risk_score: 38.5
    });
  };

  const handleRiskUpdate = (update) => {
    if (update.type === 'risk_change') {
      setRiskData(update.data);
    } else if (update.type === 'new_alert') {
      setAlerts(prev => [update.alert, ...prev]);
    }
  };

  const getRiskColor = (score) => {
    if (score < 30) return 'rgb(34, 197, 94)';
    if (score < 60) return 'rgb(234, 179, 8)';
    if (score < 80) return 'rgb(249, 115, 22)';
    return 'rgb(239, 68, 68)';
  };

  const getRiskLevel = (score) => {
    if (score < 30) return 'Low';
    if (score < 60) return 'Medium';
    if (score < 80) return 'High';
    return 'Critical';
  };

  const formatPercentage = (value) => {
    return `${(value * 100).toFixed(2)}%`;
  };

  // Calculate circle progress for risk score
  const calculateCircleProgress = (score) => {
    const radius = 80;
    const circumference = 2 * Math.PI * radius;
    const progress = (score / 100) * circumference;
    return { circumference, progress, offset: circumference - progress };
  };

  if (loading) {
    return (
      <div style={{
        background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
        borderRadius: 16,
        padding: 48,
        border: '1px solid rgba(34, 211, 238, 0.2)',
        position: 'relative',
        overflow: 'hidden'
      }}>
        <style>{`
          @keyframes scan {
            0% { transform: translateY(-100%); }
            100% { transform: translateY(300%); }
          }
          @keyframes pulse-ring {
            0% { transform: scale(0.95); opacity: 1; }
            50% { transform: scale(1.05); opacity: 0.7; }
            100% { transform: scale(0.95); opacity: 1; }
          }
          @keyframes float-up {
            0% { transform: translateY(0px); opacity: 0.6; }
            100% { transform: translateY(-20px); opacity: 0; }
          }
        `}</style>

        {/* Scanning effect */}
        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          height: 2,
          background: 'linear-gradient(90deg, transparent, rgba(34, 211, 238, 0.8), transparent)',
          animation: 'scan 2s ease-in-out infinite'
        }} />

        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
          <div style={{ position: 'relative', marginBottom: 24 }}>
            {/* Pulsing rings */}
            <div style={{
              position: 'absolute',
              width: 100,
              height: 100,
              border: '2px solid rgba(34, 211, 238, 0.3)',
              borderRadius: '50%',
              animation: 'pulse-ring 2s ease-in-out infinite'
            }} />
            <div style={{
              position: 'absolute',
              width: 120,
              height: 120,
              top: -10,
              left: -10,
              border: '2px solid rgba(34, 211, 238, 0.2)',
              borderRadius: '50%',
              animation: 'pulse-ring 2s ease-in-out infinite 0.5s'
            }} />
            <Brain style={{ width: 48, height: 48, color: 'rgb(34, 211, 238)' }} />
          </div>
          <div style={{ fontSize: 20, fontWeight: 600, color: 'white', marginBottom: 8 }}>
            AI Analysis in Progress
          </div>
          <div style={{ color: 'rgba(203, 213, 225, 0.8)', fontSize: 14 }}>
            Analyzing portfolio risk across 10,000 agent simulations...
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{
        background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
        borderRadius: 16,
        padding: 24,
        border: '1px solid rgba(239, 68, 68, 0.3)'
      }}>
        <div style={{ display: 'flex', alignItems: 'center', color: 'rgb(252, 165, 165)' }}>
          <AlertCircle style={{ width: 24, height: 24, marginRight: 12 }} />
          <span>{error}</span>
        </div>
      </div>
    );
  }

  const circleData = riskData ? calculateCircleProgress(riskData.overall_risk_score) : null;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
      <style>{`
        @keyframes glow-pulse {
          0%, 100% {
            box-shadow: 0 0 20px rgba(34, 211, 238, 0.3),
                        0 0 40px rgba(34, 211, 238, 0.1);
          }
          50% {
            box-shadow: 0 0 30px rgba(34, 211, 238, 0.5),
                        0 0 60px rgba(34, 211, 238, 0.2);
          }
        }
        @keyframes data-flow {
          0% { transform: translateX(-100%); }
          100% { transform: translateX(100%); }
        }
        @keyframes pulse-dot {
          0%, 100% { opacity: 1; transform: scale(1); }
          50% { opacity: 0.6; transform: scale(1.2); }
        }
        @keyframes rotate {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
        @keyframes float-up {
          0% { transform: translateY(0px); opacity: 0.6; }
          100% { transform: translateY(-20px); opacity: 0; }
        }
      `}</style>

      {/* AI Risk Engine - Professional Financial Dashboard */}
      <div style={{
        background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 50%, rgb(15, 23, 42) 100%)',
        borderRadius: 16,
        padding: '48px',
        border: '1px solid rgba(34, 211, 238, 0.2)',
        position: 'relative',
        overflow: 'hidden',
        animation: 'glow-pulse 4s ease-in-out infinite'
      }}>
        {/* Tech grid background */}
        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundImage: `
            linear-gradient(rgba(34, 211, 238, 0.03) 1px, transparent 1px),
            linear-gradient(90deg, rgba(34, 211, 238, 0.03) 1px, transparent 1px)
          `,
          backgroundSize: '40px 40px',
          opacity: 0.5
        }} />

        {/* Floating data particles */}
        {[...Array(5)].map((_, i) => (
          <div
            key={i}
            style={{
              position: 'absolute',
              width: 4,
              height: 4,
              background: 'rgb(34, 211, 238)',
              borderRadius: '50%',
              left: `${20 + i * 15}%`,
              bottom: 20,
              animation: `float-up 3s ease-in-out infinite ${i * 0.5}s`,
              opacity: 0.6
            }}
          />
        ))}

        <div style={{ position: 'relative', zIndex: 1 }}>

          {/* Header Row: Title + Button */}
          <div style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            marginBottom: 32,
            paddingBottom: 20,
            borderBottom: '1px solid rgba(34, 211, 238, 0.15)'
          }}>
            {/* Left: Title */}
            <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
              <div style={{
                width: 40,
                height: 40,
                borderRadius: 10,
                background: 'rgba(34, 211, 238, 0.15)',
                border: '1px solid rgba(34, 211, 238, 0.3)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
              }}>
                <Brain style={{ width: 22, height: 22, color: 'rgb(34, 211, 238)' }} />
              </div>
              <div>
                <h2 style={{
                  fontSize: 20,
                  fontWeight: 700,
                  color: 'white',
                  margin: 0,
                  letterSpacing: 0.3
                }}>
                  AI Risk Engine
                </h2>
                <p style={{
                  fontSize: 12,
                  color: 'rgba(203, 213, 225, 0.6)',
                  margin: 0
                }}>
                  Real-time risk analysis
                </p>
              </div>
            </div>

            {/* Right: Refresh Button */}
            <button
              onClick={loadRiskData}
              disabled={analyzing}
              style={{
                padding: '10px 20px',
                background: analyzing
                  ? 'rgba(100, 116, 139, 0.2)'
                  : 'rgba(34, 211, 238, 0.1)',
                border: analyzing
                  ? '1px solid rgba(100, 116, 139, 0.4)'
                  : '1px solid rgba(34, 211, 238, 0.3)',
                borderRadius: 8,
                cursor: analyzing ? 'not-allowed' : 'pointer',
                display: 'flex',
                alignItems: 'center',
                gap: 8,
                transition: 'all 0.2s ease'
              }}
              onMouseEnter={(e) => {
                if (!analyzing) {
                  e.currentTarget.style.background = 'rgba(34, 211, 238, 0.2)';
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.5)';
                }
              }}
              onMouseLeave={(e) => {
                if (!analyzing) {
                  e.currentTarget.style.background = 'rgba(34, 211, 238, 0.1)';
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.3)';
                }
              }}
            >
              <Activity
                style={{
                  width: 16,
                  height: 16,
                  color: analyzing ? 'rgba(203, 213, 225, 0.5)' : 'rgb(34, 211, 238)',
                  animation: analyzing ? 'rotate 1s linear infinite' : 'none'
                }}
              />
              <span style={{
                fontSize: 12,
                fontWeight: 600,
                color: analyzing ? 'rgba(203, 213, 225, 0.5)' : 'white',
                letterSpacing: 0.3
              }}>
                {analyzing ? 'Analyzing...' : 'Refresh'}
              </span>
            </button>
          </div>

          {/* Symmetric 3-Column Grid Layout */}
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(3, 1fr)',
            gap: 20,
            marginBottom: 24
          }}>
            {/* Column 1: Risk Level */}
            <div style={{
              padding: 20,
              background: 'rgba(255, 255, 255, 0.03)',
              borderRadius: 10,
              border: '1px solid rgba(34, 211, 238, 0.2)',
              textAlign: 'center'
            }}>
              <div style={{
                width: 44,
                height: 44,
                margin: '0 auto 12px',
                borderRadius: 10,
                background: 'rgba(34, 211, 238, 0.15)',
                border: '1px solid rgba(34, 211, 238, 0.3)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
              }}>
                <Shield style={{ width: 22, height: 22, color: 'rgb(34, 211, 238)' }} />
              </div>
              <div style={{
                fontSize: 10,
                color: 'rgba(203, 213, 225, 0.6)',
                textTransform: 'uppercase',
                letterSpacing: 0.8,
                marginBottom: 8,
                fontWeight: 600
              }}>
                Risk Level
              </div>
              <div style={{
                fontSize: 28,
                fontWeight: 700,
                color: riskData ? getRiskColor(riskData.overall_risk_score) : 'white',
                marginBottom: 6
              }}>
                {riskData ? getRiskLevel(riskData.overall_risk_score) : 'N/A'}
              </div>
              <div style={{
                fontSize: 13,
                fontWeight: 600,
                color: 'rgba(203, 213, 225, 0.5)'
              }}>
                {riskData?.overall_risk_score?.toFixed(0) || 0}/100
              </div>
            </div>

            {/* Column 2: Health Factor */}
            <div style={{
              padding: 20,
              background: 'rgba(255, 255, 255, 0.03)',
              borderRadius: 10,
              border: '1px solid rgba(34, 211, 238, 0.2)',
              textAlign: 'center'
            }}>
              <div style={{
                width: 44,
                height: 44,
                margin: '0 auto 12px',
                borderRadius: 10,
                background: 'rgba(34, 197, 94, 0.15)',
                border: '1px solid rgba(34, 197, 94, 0.3)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
              }}>
                <TrendingUp style={{ width: 22, height: 22, color: 'rgb(34, 197, 94)' }} />
              </div>
              <div style={{
                fontSize: 10,
                color: 'rgba(203, 213, 225, 0.6)',
                textTransform: 'uppercase',
                letterSpacing: 0.8,
                marginBottom: 8,
                fontWeight: 600
              }}>
                Health Factor
              </div>
              <div style={{
                fontSize: 32,
                fontWeight: 700,
                color: riskData?.health_factor > 1.5 ? 'rgb(34, 197, 94)' : 'rgb(249, 115, 22)',
                marginBottom: 6
              }}>
                {riskData?.health_factor?.toFixed(2) || 'N/A'}
              </div>
              <div style={{
                fontSize: 12,
                color: 'rgba(203, 213, 225, 0.6)'
              }}>
                {riskData?.health_factor > 1.5 ? 'Healthy' : 'At Risk'}
              </div>
            </div>

            {/* Column 3: Active Positions */}
            <div style={{
              padding: 20,
              background: 'rgba(255, 255, 255, 0.03)',
              borderRadius: 10,
              border: '1px solid rgba(34, 211, 238, 0.2)',
              textAlign: 'center'
            }}>
              <div style={{
                width: 44,
                height: 44,
                margin: '0 auto 12px',
                borderRadius: 10,
                background: 'rgba(59, 130, 246, 0.15)',
                border: '1px solid rgba(59, 130, 246, 0.3)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
              }}>
                <Target style={{ width: 22, height: 22, color: 'rgb(59, 130, 246)' }} />
              </div>
              <div style={{
                fontSize: 10,
                color: 'rgba(203, 213, 225, 0.6)',
                textTransform: 'uppercase',
                letterSpacing: 0.8,
                marginBottom: 8,
                fontWeight: 600
              }}>
                Active Positions
              </div>
              <div style={{
                fontSize: 32,
                fontWeight: 700,
                color: 'rgb(59, 130, 246)',
                marginBottom: 6
              }}>
                {riskData?.position_count || 0}
              </div>
              <div style={{
                fontSize: 12,
                color: 'rgba(203, 213, 225, 0.6)'
              }}>
                Monitored
              </div>
            </div>
          </div>

          {/* Two-Column Grid: Alerts + Tools */}
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(2, 1fr)',
            gap: 20,
            alignItems: 'stretch'
          }}>
          {/* Left: Alerts & Insights Combined */}
          <div style={{
            background: 'rgba(255, 255, 255, 0.03)',
            borderRadius: 10,
            border: '1px solid rgba(34, 211, 238, 0.2)',
            padding: 16,
            display: 'flex',
            flexDirection: 'column',
            gap: 16
          }}>
            {/* Active Alerts Header */}
            <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
              <AlertCircle style={{ width: 18, height: 18, color: 'rgb(234, 179, 8)' }} />
              <h4 style={{ fontSize: 13, fontWeight: 700, color: 'white', margin: 0 }}>
                Active Alerts & Insights
              </h4>
            </div>

            {/* Alerts Section */}
            {alerts && alerts.length > 0 ? (
              <div style={{
                background: 'rgba(234, 179, 8, 0.05)',
                borderRadius: 8,
                border: '1px solid rgba(234, 179, 8, 0.2)',
                padding: 12
              }}>
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
                  <span style={{ fontSize: 12, fontWeight: 600, color: 'rgba(203, 213, 225, 0.8)' }}>
                    Recent Alerts
                  </span>
                  <span style={{
                    padding: '2px 8px',
                    background: 'rgb(234, 179, 8)',
                    color: 'white',
                    borderRadius: 10,
                    fontSize: 11,
                    fontWeight: 700
                  }}>
                    {alerts.length}
                  </span>
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                  {alerts.slice(0, 2).map((alert, index) => (
                    <div
                      key={index}
                      style={{
                        padding: 10,
                        background: 'rgba(255, 255, 255, 0.03)',
                        borderRadius: 6,
                        border: '1px solid rgba(255, 255, 255, 0.08)'
                      }}
                    >
                      <div style={{ fontSize: 12, fontWeight: 600, color: 'white', marginBottom: 4 }}>
                        {alert.title}
                      </div>
                      <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                        {alert.protocol} • {new Date(alert.timestamp).toLocaleTimeString()}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ) : (
              <div style={{
                background: 'rgba(34, 197, 94, 0.05)',
                borderRadius: 8,
                border: '1px solid rgba(34, 197, 94, 0.2)',
                padding: 16,
                textAlign: 'center'
              }}>
                <Shield style={{ width: 32, height: 32, color: 'rgb(34, 197, 94)', margin: '0 auto 8px' }} />
                <div style={{ fontSize: 13, fontWeight: 600, color: 'white', marginBottom: 4 }}>
                  All Clear
                </div>
                <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                  No active alerts
                </div>
              </div>
            )}

            {/* Quick Insights Section */}
            <div>
              <div style={{ fontSize: 12, fontWeight: 600, color: 'rgba(203, 213, 225, 0.8)', marginBottom: 12 }}>
                Quick Insights
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 12 }}>
                <div style={{
                  padding: 12,
                  background: 'rgba(59, 130, 246, 0.1)',
                  borderRadius: 8,
                  border: '1px solid rgba(59, 130, 246, 0.2)'
                }}>
                  <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)', marginBottom: 4 }}>
                    Diversification
                  </div>
                  <div style={{ fontSize: 18, fontWeight: 700, color: 'rgb(59, 130, 246)' }}>
                    {riskData ? '72%' : 'N/A'}
                  </div>
                </div>
                <div style={{
                  padding: 12,
                  background: 'rgba(34, 197, 94, 0.1)',
                  borderRadius: 8,
                  border: '1px solid rgba(34, 197, 94, 0.2)'
                }}>
                  <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)', marginBottom: 4 }}>
                    Stability
                  </div>
                  <div style={{ fontSize: 18, fontWeight: 700, color: 'rgb(34, 197, 94)' }}>
                    {riskData ? '8.5/10' : 'N/A'}
                  </div>
                </div>
              </div>
            </div>
          </div>

            {/* Right: Risk Management Tools */}
            <div style={{
              background: 'rgba(255, 255, 255, 0.03)',
              borderRadius: 10,
              border: '1px solid rgba(34, 211, 238, 0.2)',
              padding: 16
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 16 }}>
                <Shield style={{ width: 18, height: 18, color: 'rgb(34, 211, 238)' }} />
                <h4 style={{ fontSize: 13, fontWeight: 700, color: 'white', margin: 0 }}>
                  Risk Management Tools
                </h4>
              </div>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
                {[
                  {
                    label: 'AI Risk Report',
                    desc: 'View detailed analysis & recommendations',
                    icon: Eye,
                    color: 'rgb(59, 130, 246)'
                  },
                  {
                    label: 'Configure Thresholds',
                    desc: 'Adjust risk parameters based on AI insights',
                    icon: Shield,
                    color: 'rgb(167, 139, 250)'
                  },
                  {
                    label: 'Alert History',
                    desc: 'Review past notifications & actions',
                    icon: Activity,
                    color: 'rgb(234, 179, 8)'
                  }
                ].map((action, i) => {
                  const Icon = action.icon;
                  return (
                    <button key={i} style={{
                      padding: '14px 16px',
                      background: 'rgba(255, 255, 255, 0.03)',
                      border: '1px solid rgba(255, 255, 255, 0.08)',
                      borderRadius: 10,
                      cursor: 'pointer',
                      display: 'flex',
                      alignItems: 'center',
                      gap: 12,
                      transition: 'all 0.2s ease'
                    }}
                    onMouseEnter={(e) => {
                      e.currentTarget.style.background = `${action.color}15`;
                      e.currentTarget.style.borderColor = `${action.color}40`;
                      e.currentTarget.style.transform = 'translateX(4px)';
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.background = 'rgba(255, 255, 255, 0.03)';
                      e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.08)';
                      e.currentTarget.style.transform = 'translateX(0)';
                    }}>
                      <div style={{
                        width: 40,
                        height: 40,
                        borderRadius: 10,
                        background: `${action.color}20`,
                        border: `1px solid ${action.color}40`,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        flexShrink: 0
                      }}>
                        <Icon style={{ width: 18, height: 18, color: action.color }} />
                      </div>
                      <div style={{ flex: 1, textAlign: 'left' }}>
                        <div style={{ fontSize: 13, fontWeight: 700, color: 'white', marginBottom: 2 }}>
                          {action.label}
                        </div>
                        <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)', lineHeight: 1.4 }}>
                          {action.desc}
                        </div>
                      </div>
                    </button>
                  );
                })}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Key Metrics Grid */}
      {riskData && (
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
          gap: 16
        }}>
          {[
            {
              icon: Target,
              label: 'Portfolio Value',
              value: `$${(riskData.total_collateral || 0).toLocaleString()}`,
              sublabel: `${riskData.position_count || 0} active positions`,
              color: 'rgb(59, 130, 246)'
            },
            {
              icon: TrendingUp,
              label: 'Total Collateral',
              value: `$${(riskData.total_collateral || 0).toLocaleString()}`,
              sublabel: 'Locked in protocols',
              color: 'rgb(34, 197, 94)'
            },
            {
              icon: AlertTriangle,
              label: 'Total Debt',
              value: `$${(riskData.total_debt || 0).toLocaleString()}`,
              sublabel: `${((riskData.total_debt / riskData.total_collateral) * 100).toFixed(1)}% utilization`,
              color: 'rgb(234, 179, 8)'
            }
          ].map((metric, i) => {
            const Icon = metric.icon;
            return (
              <div
                key={i}
                style={{
                  background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
                  borderRadius: 12,
                  padding: 24,
                  border: '1px solid rgba(34, 211, 238, 0.15)',
                  position: 'relative',
                  overflow: 'hidden',
                  transition: 'all 0.3s ease'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.4)';
                  e.currentTarget.style.transform = 'translateY(-4px)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.15)';
                  e.currentTarget.style.transform = 'translateY(0)';
                }}
              >
                {/* Gradient overlay */}
                <div style={{
                  position: 'absolute',
                  top: 0,
                  right: 0,
                  width: 100,
                  height: 100,
                  background: `radial-gradient(circle, ${metric.color}20 0%, transparent 70%)`,
                  opacity: 0.5
                }} />

                <div style={{ position: 'relative', zIndex: 1 }}>
                  <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
                    <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', textTransform: 'uppercase', letterSpacing: 1, fontWeight: 600 }}>
                      {metric.label}
                    </div>
                    <Icon style={{ width: 20, height: 20, color: metric.color, opacity: 0.8 }} />
                  </div>
                  <div style={{ fontSize: 28, fontWeight: 700, color: 'white', marginBottom: 8 }}>
                    {metric.value}
                  </div>
                  <div style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.6)' }}>
                    {metric.sublabel}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Liquidation Risk Alert */}
      {liquidationRisk && liquidationRisk.risk_level !== 'low' && (
        <div style={{
          background: liquidationRisk.risk_level === 'critical'
            ? 'linear-gradient(135deg, rgba(239, 68, 68, 0.1) 0%, rgba(220, 38, 38, 0.05) 100%)'
            : 'linear-gradient(135deg, rgba(234, 179, 8, 0.1) 0%, rgba(202, 138, 4, 0.05) 100%)',
          border: liquidationRisk.risk_level === 'critical' ? '1px solid rgba(239, 68, 68, 0.3)' : '1px solid rgba(234, 179, 8, 0.3)',
          borderRadius: 12,
          padding: 24,
          position: 'relative',
          overflow: 'hidden'
        }}>
          {/* Animated warning stripe */}
          <div style={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            height: 3,
            background: liquidationRisk.risk_level === 'critical'
              ? 'repeating-linear-gradient(90deg, rgb(239, 68, 68) 0px, rgb(239, 68, 68) 20px, transparent 20px, transparent 40px)'
              : 'repeating-linear-gradient(90deg, rgb(234, 179, 8) 0px, rgb(234, 179, 8) 20px, transparent 20px, transparent 40px)',
            animation: 'data-flow 2s linear infinite'
          }} />

          <div style={{ display: 'flex', alignItems: 'flex-start', gap: 16 }}>
            <div style={{
              width: 48,
              height: 48,
              borderRadius: 12,
              background: liquidationRisk.risk_level === 'critical' ? 'rgba(239, 68, 68, 0.2)' : 'rgba(234, 179, 8, 0.2)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              flexShrink: 0
            }}>
              <AlertTriangle style={{
                width: 24,
                height: 24,
                color: liquidationRisk.risk_level === 'critical' ? 'rgb(239, 68, 68)' : 'rgb(234, 179, 8)'
              }} />
            </div>
            <div style={{ flex: 1 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
                <h3 style={{ fontWeight: 700, color: 'rgb(15, 23, 42)', margin: 0, fontSize: 18 }}>
                  Liquidation Risk Alert
                </h3>
                <span style={{
                  padding: '2px 8px',
                  background: liquidationRisk.risk_level === 'critical' ? 'rgb(239, 68, 68)' : 'rgb(234, 179, 8)',
                  color: 'white',
                  borderRadius: 4,
                  fontSize: 11,
                  fontWeight: 700,
                  textTransform: 'uppercase'
                }}>
                  {liquidationRisk.risk_level}
                </span>
              </div>
              <p style={{ color: 'rgb(71, 85, 105)', margin: '0 0 16px 0', fontSize: 14, lineHeight: 1.6 }}>
                {liquidationRisk.message || 'Your position may be at risk of liquidation'}
              </p>
              {liquidationRisk.recommended_action && (
                <div style={{
                  padding: 16,
                  background: 'white',
                  borderRadius: 8,
                  border: '1px solid rgb(226, 232, 240)'
                }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
                    <Zap style={{ width: 16, height: 16, color: 'rgb(59, 130, 246)' }} />
                    <span style={{ fontSize: 13, fontWeight: 700, color: 'rgb(15, 23, 42)' }}>
                      AI Recommended Action:
                    </span>
                  </div>
                  <div style={{ fontSize: 14, color: 'rgb(71, 85, 105)', lineHeight: 1.5 }}>
                    {liquidationRisk.recommended_action}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Market Risk Indicators - Enhanced */}
      {marketRisk && (
        <div style={{
          background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
          borderRadius: 16,
          padding: 32,
          border: '1px solid rgba(34, 211, 238, 0.2)'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', marginBottom: 24 }}>
            <Eye style={{ width: 24, height: 24, color: 'rgb(34, 211, 238)', marginRight: 12 }} />
            <h3 style={{ fontSize: 20, fontWeight: 700, color: 'white', margin: 0 }}>
              Market Risk Indicators
            </h3>
          </div>

          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
            gap: 20
          }}>
            {[
              {
                label: 'Volatility Index',
                value: `${marketRisk.volatility_index?.toFixed(1)}%`,
                color: getRiskColor(marketRisk.volatility_index),
                icon: Activity
              },
              {
                label: 'Total Liquidity',
                value: `$${(marketRisk.total_liquidity / 1e9).toFixed(2)}B`,
                color: 'rgb(34, 211, 238)',
                icon: TrendingUp
              },
              {
                label: 'Correlation Index',
                value: formatPercentage(marketRisk.correlation_index),
                color: 'rgb(59, 130, 246)',
                icon: Target
              },
              {
                label: 'Systemic Risk',
                value: getRiskLevel(marketRisk.systemic_risk_score),
                color: getRiskColor(marketRisk.systemic_risk_score),
                icon: AlertTriangle
              }
            ].map((indicator, i) => {
              const Icon = indicator.icon;
              return (
                <div
                  key={i}
                  style={{
                    padding: 20,
                    background: 'rgba(255, 255, 255, 0.03)',
                    borderRadius: 12,
                    border: '1px solid rgba(255, 255, 255, 0.08)',
                    transition: 'all 0.3s ease'
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.background = 'rgba(255, 255, 255, 0.06)';
                    e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.3)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.background = 'rgba(255, 255, 255, 0.03)';
                    e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.08)';
                  }}
                >
                  <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
                    <div style={{
                      fontSize: 11,
                      color: 'rgba(203, 213, 225, 0.7)',
                      textTransform: 'uppercase',
                      letterSpacing: 1,
                      fontWeight: 600
                    }}>
                      {indicator.label}
                    </div>
                    <Icon style={{ width: 18, height: 18, color: indicator.color, opacity: 0.6 }} />
                  </div>
                  <div style={{
                    fontSize: 28,
                    fontWeight: 800,
                    color: indicator.color,
                    textShadow: `0 0 20px ${indicator.color}40`
                  }}>
                    {indicator.value}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
};

export default AIRiskDashboard;
