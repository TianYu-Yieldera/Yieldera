/**
 * AI Risk Dashboard Pro
 * Professional risk monitoring dashboard with real functionality
 * Integrates FastAPI risk engine and liquidation alerts
 */

import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import {
  AlertTriangle,
  Shield,
  Activity,
  TrendingUp,
  Brain,
  Zap,
  RefreshCw,
  Eye,
  ChevronRight,
  Users,
  TrendingDown
} from 'lucide-react';
import { useWallet } from '../web3/WalletContext';
import { useDemoMode } from '../web3/DemoModeContext';
import LiquidationAlert from './LiquidationAlert';
import AgentSimulationPanel from './AgentSimulationPanel';
import StressTestPanel from './StressTestPanel';
import MLPredictionPanel from './MLPredictionPanel';
import RealtimeMonitorPanel from './RealtimeMonitorPanel';
import fastAPIRiskService from '../services/fastAPIRiskService';
import aiRiskService from '../services/aiRiskService';

const AIRiskDashboardPro = () => {
  const [searchParams] = useSearchParams();
  const { address } = useWallet();
  const { demoMode, demoData } = useDemoMode();
  const [positions, setPositions] = useState([]);
  const [riskSummary, setRiskSummary] = useState(null);
  const [loading, setLoading] = useState(true);
  const [fastAPIAvailable, setFastAPIAvailable] = useState(false);
  const [selectedTab, setSelectedTab] = useState(() => {
    // Check URL params for initial tab
    return searchParams.get('tab') === 'analytics' ? 'analytics' : 'positions';
  }); // 'positions' | 'analytics'
  const [analyticsSubTab, setAnalyticsSubTab] = useState('simulation'); // 'simulation' | 'stress' | 'prediction' | 'realtime'

  useEffect(() => {
    loadDashboardData();
    checkFastAPIStatus();
  }, [address, demoMode, demoData]);

  const checkFastAPIStatus = async () => {
    try {
      const health = await fastAPIRiskService.checkHealth();
      setFastAPIAvailable(health.available);
    } catch (error) {
      setFastAPIAvailable(false);
    }
  };

  const loadDashboardData = async () => {
    try {
      setLoading(true);

      if (demoMode) {
        // Load demo positions from demoData
        const demoPositions = generateDemoPositions(demoData);
        setPositions(demoPositions);

        // Calculate risk summary
        const summary = calculateRiskSummary(demoPositions);
        setRiskSummary(summary);
      } else {
        // Try to load real data
        try {
          const portfolioRisk = await aiRiskService.getPortfolioRisk(address);
          const alerts = await aiRiskService.getRiskAlerts(address);

          // Convert to positions format
          const realPositions = convertToPositions(portfolioRisk, alerts);
          setPositions(realPositions);

          const summary = calculateRiskSummary(realPositions);
          setRiskSummary(summary);
        } catch (err) {
          console.warn('Failed to load real data, using demo:', err);
          const demoPositions = generateDemoPositions(demoData);
          setPositions(demoPositions);
          const summary = calculateRiskSummary(demoPositions);
          setRiskSummary(summary);
        }
      }
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  const generateDemoPositions = (demoData) => {
    const defiDeposits = demoData?.defiDeposits || [];

    // If no deposits exist, generate default demo positions
    if (defiDeposits.length === 0) {
      return [
        {
          id: 'demo-default-1',
          protocol: 'Aave V3',
          chain: 'Arbitrum',
          collateralAsset: 'ETH',
          collateralValueUSD: 50000,
          borrowValueUSD: 20000,
          healthFactor: 2.5,
          leverage: 1.4,
          ltv: 0.4,
          currentPrice: 2500,
          liquidationPrice: 1950,
          position_age_days: 45
        },
        {
          id: 'demo-default-2',
          protocol: 'Compound V3',
          chain: 'Base',
          collateralAsset: 'WBTC',
          collateralValueUSD: 75000,
          borrowValueUSD: 45000,
          healthFactor: 1.67,
          leverage: 1.6,
          ltv: 0.6,
          currentPrice: 45000,
          liquidationPrice: 40500,
          position_age_days: 30
        },
        {
          id: 'demo-default-3',
          protocol: 'GMX V2',
          chain: 'Arbitrum',
          collateralAsset: 'ETH',
          collateralValueUSD: 40000,
          borrowValueUSD: 28000,
          healthFactor: 1.43,
          leverage: 1.7,
          ltv: 0.7,
          currentPrice: 2500,
          liquidationPrice: 2200,
          position_age_days: 60
        }
      ];
    }

    return defiDeposits.map((deposit, index) => ({
      id: `demo-${index}`,
      protocol: deposit.protocol || 'Aave V3',
      chain: deposit.protocol?.includes('GMX') ? 'Arbitrum' : 'Base',
      collateralAsset: deposit.token || 'ETH',
      collateralValueUSD: deposit.amount || 1000,
      borrowValueUSD: deposit.amount * 0.6 || 600, // 60% LTV
      healthFactor: 2.5 - (index * 0.3), // Declining health factors
      leverage: 1 + (index * 0.2),
      ltv: 0.4 + (index * 0.1),
      currentPrice: deposit.token === 'ETH' ? 2500 : deposit.token === 'WBTC' ? 45000 : 1,
      liquidationPrice: deposit.token === 'ETH' ? 2000 : deposit.token === 'WBTC' ? 38000 : 0.95,
      position_age_days: 30 + index * 15
    }));
  };

  const convertToPositions = (portfolioRisk, alerts) => {
    // Convert API data to positions format
    // This is a placeholder - adjust based on actual API response
    return [];
  };

  const calculateRiskSummary = (positions) => {
    if (!positions || positions.length === 0) {
      return {
        totalCollateral: 0,
        totalDebt: 0,
        avgHealthFactor: 0,
        positionsAtRisk: 0,
        overallRiskScore: 0,
        riskLevel: 'None'
      };
    }

    const totalCollateral = positions.reduce((sum, p) => sum + p.collateralValueUSD, 0);
    const totalDebt = positions.reduce((sum, p) => sum + p.borrowValueUSD, 0);
    const avgHealthFactor = positions.reduce((sum, p) => sum + p.healthFactor, 0) / positions.length;
    const positionsAtRisk = positions.filter(p => p.healthFactor < 1.5).length;

    // Calculate overall risk score (0-100)
    let riskScore = 0;
    riskScore += (100 - (avgHealthFactor * 30)); // Health factor impact
    riskScore += (positionsAtRisk * 15); // At-risk positions impact
    riskScore = Math.max(0, Math.min(100, riskScore));

    const riskLevel =
      riskScore < 30 ? 'Low' :
      riskScore < 50 ? 'Medium' :
      riskScore < 70 ? 'High' : 'Critical';

    return {
      totalCollateral,
      totalDebt,
      avgHealthFactor,
      positionsAtRisk,
      overallRiskScore: riskScore,
      riskLevel
    };
  };

  const handleAddCollateral = (position) => {
    alert(`Add collateral to ${position.protocol} position\n\nThis would open a modal to add collateral and improve health factor.`);
  };

  const handleRefresh = () => {
    loadDashboardData();
    checkFastAPIStatus();
  };

  const getRiskLevelColor = (level) => {
    switch (level) {
      case 'Low': return { bg: 'rgba(34, 197, 94, 0.1)', text: 'rgb(34, 197, 94)', border: 'rgba(34, 197, 94, 0.3)' };
      case 'Medium': return { bg: 'rgba(234, 179, 8, 0.1)', text: 'rgb(234, 179, 8)', border: 'rgba(234, 179, 8, 0.3)' };
      case 'High': return { bg: 'rgba(249, 115, 22, 0.1)', text: 'rgb(249, 115, 22)', border: 'rgba(249, 115, 22, 0.3)' };
      case 'Critical': return { bg: 'rgba(239, 68, 68, 0.1)', text: 'rgb(239, 68, 68)', border: 'rgba(239, 68, 68, 0.3)' };
      default: return { bg: 'rgba(148, 163, 184, 0.1)', text: 'rgb(148, 163, 184)', border: 'rgba(148, 163, 184, 0.3)' };
    }
  };

  if (loading) {
    return (
      <div className="p-8 text-center">
        <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p className="mt-4 text-gray-600">Loading risk dashboard...</p>
      </div>
    );
  }

  const riskColors = riskSummary ? getRiskLevelColor(riskSummary.riskLevel) : getRiskLevelColor('None');

  return (
    <div style={{
      background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
      borderRadius: 16,
      padding: 32,
      border: '1px solid rgba(34, 211, 238, 0.2)',
      position: 'relative',
      overflow: 'hidden'
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
        backgroundSize: '20px 20px',
        opacity: 0.5
      }} />

      <div style={{ position: 'relative', zIndex: 1 }}>
        {/* Header */}
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 24 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <div style={{
              width: 48,
              height: 48,
              borderRadius: 12,
              background: 'rgba(34, 211, 238, 0.15)',
              border: '1px solid rgba(34, 211, 238, 0.3)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
              <Brain style={{ width: 24, height: 24, color: 'rgb(34, 211, 238)' }} />
            </div>
            <div>
              <h2 style={{ fontSize: 24, fontWeight: 700, color: 'white', margin: 0 }}>
                AI Risk Engine
              </h2>
              <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.7)', margin: '4px 0 0 0', display: 'flex', alignItems: 'center', gap: 6 }}>
                {fastAPIAvailable ? (
                  <>
                    <Zap size={14} style={{ color: 'rgb(34, 197, 94)' }} />
                    <span>Advanced analytics powered by FastAPI</span>
                  </>
                ) : (
                  <>
                    <Activity size={14} style={{ color: 'rgb(148, 163, 184)' }} />
                    <span>Basic risk monitoring</span>
                  </>
                )}
              </p>
            </div>
          </div>

          <button
            onClick={handleRefresh}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 8,
              padding: '10px 18px',
              background: 'rgba(34, 211, 238, 0.15)',
              border: '1px solid rgba(34, 211, 238, 0.3)',
              borderRadius: 8,
              color: 'rgb(34, 211, 238)',
              fontSize: 14,
              fontWeight: 600,
              cursor: 'pointer',
              transition: 'all 0.2s'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = 'rgba(34, 211, 238, 0.25)';
              e.currentTarget.style.transform = 'scale(1.05)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = 'rgba(34, 211, 238, 0.15)';
              e.currentTarget.style.transform = 'scale(1)';
            }}
          >
            <RefreshCw size={16} />
            Refresh
          </button>
        </div>

        {/* Risk Summary Cards */}
        {riskSummary && (
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
            gap: 16,
            marginBottom: 24
          }}>
            <div style={{
              background: 'rgba(255, 255, 255, 0.05)',
              borderRadius: 12,
              padding: 20,
              border: '1px solid rgba(255, 255, 255, 0.1)'
            }}>
              <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 600 }}>
                Risk Level
              </div>
              <div style={{
                display: 'flex',
                alignItems: 'center',
                gap: 8
              }}>
                <div style={{
                  fontSize: 28,
                  fontWeight: 700,
                  color: riskColors.text
                }}>
                  {riskSummary.riskLevel}
                </div>
                <div style={{
                  padding: '4px 10px',
                  background: riskColors.bg,
                  border: `1px solid ${riskColors.border}`,
                  borderRadius: 6,
                  fontSize: 12,
                  fontWeight: 700,
                  color: riskColors.text
                }}>
                  {riskSummary.overallRiskScore.toFixed(0)}
                </div>
              </div>
            </div>

            <div style={{
              background: 'rgba(255, 255, 255, 0.05)',
              borderRadius: 12,
              padding: 20,
              border: '1px solid rgba(255, 255, 255, 0.1)'
            }}>
              <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 600 }}>
                Avg Health Factor
              </div>
              <div style={{ fontSize: 28, fontWeight: 700, color: 'white' }}>
                {riskSummary.avgHealthFactor.toFixed(2)}
              </div>
            </div>

            <div style={{
              background: 'rgba(255, 255, 255, 0.05)',
              borderRadius: 12,
              padding: 20,
              border: '1px solid rgba(255, 255, 255, 0.1)'
            }}>
              <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 600 }}>
                Total Collateral
              </div>
              <div style={{ fontSize: 28, fontWeight: 700, color: 'white' }}>
                ${riskSummary.totalCollateral.toLocaleString()}
              </div>
            </div>

            <div style={{
              background: 'rgba(255, 255, 255, 0.05)',
              borderRadius: 12,
              padding: 20,
              border: '1px solid rgba(255, 255, 255, 0.1)'
            }}>
              <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 600 }}>
                Positions at Risk
              </div>
              <div style={{
                fontSize: 28,
                fontWeight: 700,
                color: riskSummary.positionsAtRisk > 0 ? 'rgb(239, 68, 68)' : 'rgb(34, 197, 94)'
              }}>
                {riskSummary.positionsAtRisk} / {positions.length}
              </div>
            </div>
          </div>
        )}

        {/* Tabs */}
        <div style={{ display: 'flex', gap: 12, marginBottom: 24, borderBottom: '1px solid rgba(255, 255, 255, 0.1)' }}>
          <button
            onClick={() => setSelectedTab('positions')}
            style={{
              padding: '12px 20px',
              background: selectedTab === 'positions' ? 'rgba(34, 211, 238, 0.15)' : 'transparent',
              border: 'none',
              borderBottom: selectedTab === 'positions' ? '2px solid rgb(34, 211, 238)' : '2px solid transparent',
              color: selectedTab === 'positions' ? 'rgb(34, 211, 238)' : 'rgba(203, 213, 225, 0.7)',
              fontSize: 14,
              fontWeight: 600,
              cursor: 'pointer',
              transition: 'all 0.2s'
            }}
          >
            <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
              <Shield size={16} />
              Positions ({positions.length})
            </div>
          </button>
          <button
            onClick={() => setSelectedTab('analytics')}
            style={{
              padding: '12px 20px',
              background: selectedTab === 'analytics' ? 'rgba(34, 211, 238, 0.15)' : 'transparent',
              border: 'none',
              borderBottom: selectedTab === 'analytics' ? '2px solid rgb(34, 211, 238)' : '2px solid transparent',
              color: selectedTab === 'analytics' ? 'rgb(34, 211, 238)' : 'rgba(203, 213, 225, 0.7)',
              fontSize: 14,
              fontWeight: 600,
              cursor: 'pointer',
              transition: 'all 0.2s'
            }}
          >
            <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
              <TrendingUp size={16} />
              Analytics
            </div>
          </button>
        </div>

        {/* Content */}
        {selectedTab === 'positions' ? (
          <div>
            {positions.length === 0 ? (
              <div style={{
                padding: 40,
                textAlign: 'center',
                background: 'rgba(255, 255, 255, 0.03)',
                borderRadius: 12,
                border: '1px solid rgba(255, 255, 255, 0.08)'
              }}>
                <Shield size={48} style={{ color: 'rgba(203, 213, 225, 0.5)', marginBottom: 16 }} />
                <p style={{ fontSize: 16, color: 'rgba(203, 213, 225, 0.8)', margin: '0 0 8px 0', fontWeight: 600 }}>
                  No Active Positions
                </p>
                <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.6)', margin: 0 }}>
                  Open a position in the Vault to start monitoring risk
                </p>
              </div>
            ) : (
              <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
                {positions.map((position) => (
                  <LiquidationAlert
                    key={position.id}
                    position={position}
                    marketData={{
                      asset: position.collateralAsset,
                      price: position.currentPrice,
                      volatility: 0.15
                    }}
                    historicalPrices={[]} // Will use basic calculations
                    onAddCollateral={handleAddCollateral}
                    refreshInterval={0} // Disable auto-refresh in dashboard
                  />
                ))}
              </div>
            )}
          </div>
        ) : (
          <div>
            {/* Analytics Sub-tabs */}
            <div style={{ display: 'flex', gap: 12, marginBottom: 24, borderBottom: '1px solid rgba(255, 255, 255, 0.1)' }}>
              <button
                onClick={() => setAnalyticsSubTab('simulation')}
                style={{
                  padding: '12px 20px',
                  background: analyticsSubTab === 'simulation' ? 'rgba(34, 211, 238, 0.15)' : 'transparent',
                  border: 'none',
                  borderBottom: analyticsSubTab === 'simulation' ? '2px solid rgb(34, 211, 238)' : '2px solid transparent',
                  color: analyticsSubTab === 'simulation' ? 'rgb(34, 211, 238)' : 'rgba(203, 213, 225, 0.7)',
                  fontSize: 14,
                  fontWeight: 600,
                  cursor: 'pointer',
                  transition: 'all 0.2s'
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <Users size={16} />
                  Agent Simulation
                </div>
              </button>
              <button
                onClick={() => setAnalyticsSubTab('stress')}
                style={{
                  padding: '12px 20px',
                  background: analyticsSubTab === 'stress' ? 'rgba(34, 211, 238, 0.15)' : 'transparent',
                  border: 'none',
                  borderBottom: analyticsSubTab === 'stress' ? '2px solid rgb(34, 211, 238)' : '2px solid transparent',
                  color: analyticsSubTab === 'stress' ? 'rgb(34, 211, 238)' : 'rgba(203, 213, 225, 0.7)',
                  fontSize: 14,
                  fontWeight: 600,
                  cursor: 'pointer',
                  transition: 'all 0.2s'
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <TrendingDown size={16} />
                  Stress Testing
                </div>
              </button>
              <button
                onClick={() => setAnalyticsSubTab('prediction')}
                style={{
                  padding: '12px 20px',
                  background: analyticsSubTab === 'prediction' ? 'rgba(34, 211, 238, 0.15)' : 'transparent',
                  border: 'none',
                  borderBottom: analyticsSubTab === 'prediction' ? '2px solid rgb(34, 211, 238)' : '2px solid transparent',
                  color: analyticsSubTab === 'prediction' ? 'rgb(34, 211, 238)' : 'rgba(203, 213, 225, 0.7)',
                  fontSize: 14,
                  fontWeight: 600,
                  cursor: 'pointer',
                  transition: 'all 0.2s'
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <Brain size={16} />
                  ML Prediction
                </div>
              </button>
              <button
                onClick={() => setAnalyticsSubTab('realtime')}
                style={{
                  padding: '12px 20px',
                  background: analyticsSubTab === 'realtime' ? 'rgba(34, 211, 238, 0.15)' : 'transparent',
                  border: 'none',
                  borderBottom: analyticsSubTab === 'realtime' ? '2px solid rgb(34, 211, 238)' : '2px solid transparent',
                  color: analyticsSubTab === 'realtime' ? 'rgb(34, 211, 238)' : 'rgba(203, 213, 225, 0.7)',
                  fontSize: 14,
                  fontWeight: 600,
                  cursor: 'pointer',
                  transition: 'all 0.2s'
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <Activity size={16} />
                  Live Monitor
                </div>
              </button>
            </div>

            {/* Sub-tab Content */}
            {analyticsSubTab === 'simulation' ? (
              <AgentSimulationPanel positions={positions} />
            ) : analyticsSubTab === 'stress' ? (
              <StressTestPanel positions={positions} />
            ) : analyticsSubTab === 'prediction' ? (
              <MLPredictionPanel positions={positions} />
            ) : (
              <RealtimeMonitorPanel />
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default AIRiskDashboardPro;
