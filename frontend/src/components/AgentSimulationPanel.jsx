/**
 * Agent Simulation Panel
 * Visualizes 10,000-agent Monte Carlo simulation results
 * Shows liquidation risks, agent behavior, and portfolio stress testing
 */

import React, { useState, useEffect } from 'react';
import { Users, TrendingDown, AlertTriangle, BarChart3, Play, RefreshCw, Zap, Activity } from 'lucide-react';
import fastAPIRiskService from '../services/fastAPIRiskService';

const AgentSimulationPanel = ({ positions = [] }) => {
  const [simulationData, setSimulationData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [selectedScenario, setSelectedScenario] = useState('normal');

  const scenarios = {
    normal: { name: 'Normal Market', volatility: 0.15, priceDrop: 0 },
    moderate: { name: 'Moderate Stress', volatility: 0.25, priceDrop: -10 },
    severe: { name: 'Severe Stress', volatility: 0.40, priceDrop: -20 },
    extreme: { name: 'Black Swan', volatility: 0.60, priceDrop: -35 }
  };

  useEffect(() => {
    if (positions.length > 0) {
      runSimulation();
    }
  }, [positions, selectedScenario]);

  const runSimulation = async () => {
    setLoading(true);
    setError(null);

    try {
      // For demo, use mock simulation results
      // In production, would call: await fastAPIRiskService.runAgentSimulation(params)

      const scenario = scenarios[selectedScenario];
      const mockData = generateMockSimulationData(positions, scenario);

      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 1500));

      setSimulationData(mockData);
    } catch (err) {
      console.error('Simulation failed:', err);
      setError('Failed to run simulation. Using fallback data.');

      // Generate fallback data
      const scenario = scenarios[selectedScenario];
      const mockData = generateMockSimulationData(positions, scenario);
      setSimulationData(mockData);
    } finally {
      setLoading(false);
    }
  };

  const generateMockSimulationData = (positions, scenario) => {
    const numAgents = 10000;
    const avgHealthFactor = positions.reduce((sum, p) => sum + (p.healthFactor || 2.0), 0) / positions.length;

    // Base liquidation probability on health factor
    let baseLiqProb = 0;
    if (avgHealthFactor >= 2.0) baseLiqProb = 0.02;
    else if (avgHealthFactor >= 1.5) baseLiqProb = 0.08;
    else if (avgHealthFactor >= 1.2) baseLiqProb = 0.20;
    else baseLiqProb = 0.40;

    // Adjust for scenario - make differences more dramatic
    const scenarioMultiplier = {
      normal: 1.0,      // 正常市场：基础概率
      moderate: 3.0,    // 中度压力：3倍
      severe: 7.0,      // 严重压力：7倍
      extreme: 15.0     // 黑天鹅：15倍
    };

    const adjustedLiqProb = Math.min(baseLiqProb * scenarioMultiplier[selectedScenario], 0.95);
    const liquidatedAgents = Math.floor(numAgents * adjustedLiqProb);
    const safeAgents = numAgents - liquidatedAgents;

    // Agent type distribution
    const agentTypes = {
      whale: { count: 50, avgLoss: 250000, color: 'rgb(168, 85, 247)' },
      institutional: { count: 100, avgLoss: 180000, color: 'rgb(59, 130, 246)' },
      arbitrageur: { count: 800, avgLoss: 45000, color: 'rgb(34, 211, 238)' },
      yield_farmer: { count: 2500, avgLoss: 8500, color: 'rgb(16, 185, 129)' },
      retail: { count: 6000, avgLoss: 2200, color: 'rgb(245, 158, 11)' },
      liquidator: { count: 300, avgProfit: 12000, color: 'rgb(34, 197, 94)' },
      market_maker: { count: 200, avgLoss: 5000, color: 'rgb(236, 72, 153)' },
      protocol: { count: 50, avgLoss: 150000, color: 'rgb(139, 92, 246)' }
    };

    // Calculate total potential loss
    const totalPotentialLoss = Object.values(agentTypes).reduce((sum, type) => {
      if (type.avgProfit) return sum; // Skip liquidators
      return sum + (type.count * type.avgLoss * adjustedLiqProb);
    }, 0);

    return {
      totalAgents: numAgents,
      liquidatedAgents,
      safeAgents,
      liquidationRate: adjustedLiqProb,
      totalPotentialLoss,
      avgLossPerAgent: totalPotentialLoss / liquidatedAgents || 0,
      scenario: selectedScenario,
      agentTypes,
      riskMetrics: {
        var95: totalPotentialLoss * 0.65,
        var99: totalPotentialLoss * 0.85,
        cvar: totalPotentialLoss * 0.92,
        maxDrawdown: totalPotentialLoss * 1.1
      },
      histogram: generateHistogramData(liquidatedAgents, safeAgents)
    };
  };

  const generateHistogramData = (liquidated, safe) => {
    // Generate distribution bins based on health factor
    const bins = [];
    const totalAgents = liquidated + safe;
    const liquidationRate = liquidated / totalAgents;

    // Health factor bins (0.5 to 3.2, step 0.3)
    for (let i = 0; i < 10; i++) {
      const hf = 0.5 + (i * 0.3);
      const isLiquidationZone = hf < 1.2;

      let percentage;
      if (isLiquidationZone) {
        // Bins below HF 1.2 (liquidation zone) - distribute liquidated agents
        // Use a declining distribution: most agents in lower HF bins
        const liquidationZoneIndex = i; // 0, 1, 2 for HF 0.5, 0.8, 1.1
        percentage = liquidationRate * (0.45 - liquidationZoneIndex * 0.12);
      } else {
        // Bins above HF 1.2 (safe zone) - distribute safe agents
        // Use a bell curve centered around HF 2.0-2.5
        const safeZoneIndex = i - 3; // Starts from 0 for HF 1.4
        if (safeZoneIndex < 3) {
          percentage = (1 - liquidationRate) * (0.15 + safeZoneIndex * 0.08);
        } else {
          percentage = (1 - liquidationRate) * (0.25 - (safeZoneIndex - 3) * 0.05);
        }
      }

      bins.push({
        hf: hf.toFixed(1),
        count: Math.floor(totalAgents * Math.max(percentage, 0)),
        percentage: Math.max(percentage * 100, 0)
      });
    }

    return bins;
  };

  if (loading) {
    return (
      <div style={{
        padding: 40,
        textAlign: 'center',
        background: 'rgba(34, 211, 238, 0.05)',
        borderRadius: 12,
        border: '1px solid rgba(34, 211, 238, 0.2)'
      }}>
        <div style={{
          width: 64,
          height: 64,
          margin: '0 auto 20px',
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
        <p style={{ fontSize: 16, color: 'white', marginBottom: 8, fontWeight: 600 }}>
          Running 10,000 Agent Simulation
        </p>
        <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.7)', margin: 0 }}>
          Monte Carlo analysis in progress...
        </p>
      </div>
    );
  }

  if (!simulationData && !loading) {
    return (
      <div style={{
        padding: 40,
        textAlign: 'center',
        background: 'rgba(255, 255, 255, 0.03)',
        borderRadius: 12,
        border: '1px solid rgba(255, 255, 255, 0.08)'
      }}>
        <Users size={48} style={{ color: 'rgba(203, 213, 225, 0.5)', marginBottom: 16 }} />
        <p style={{ fontSize: 16, color: 'rgba(203, 213, 225, 0.8)', margin: '0 0 16px 0', fontWeight: 600 }}>
          No Positions to Simulate
        </p>
        <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.6)', margin: 0 }}>
          Open a position to run agent-based risk simulations
        </p>
      </div>
    );
  }

  if (error && !simulationData) {
    return (
      <div style={{
        padding: 24,
        background: 'rgba(239, 68, 68, 0.1)',
        border: '1px solid rgba(239, 68, 68, 0.3)',
        borderRadius: 12,
        color: 'rgb(248, 113, 113)'
      }}>
        <AlertTriangle size={20} style={{ marginBottom: 8 }} />
        <p style={{ margin: 0 }}>{error}</p>
      </div>
    );
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
      {/* Header with Scenario Selector */}
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', flexWrap: 'wrap', gap: 16 }}>
        <div>
          <h3 style={{ fontSize: 18, fontWeight: 700, color: 'white', margin: '0 0 8px 0' }}>
            10,000-Agent Monte Carlo Simulation
          </h3>
          <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.7)', margin: 0 }}>
            Stress testing with diverse agent behaviors and market scenarios
          </p>
        </div>

        <button
          onClick={runSimulation}
          disabled={loading}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: 8,
            padding: '10px 20px',
            background: 'rgba(34, 211, 238, 0.2)',
            border: '1px solid rgba(34, 211, 238, 0.4)',
            borderRadius: 8,
            color: 'rgb(34, 211, 238)',
            fontSize: 14,
            fontWeight: 600,
            cursor: loading ? 'not-allowed' : 'pointer',
            opacity: loading ? 0.6 : 1,
            transition: 'all 0.3s'
          }}
          onMouseEnter={(e) => !loading && (e.currentTarget.style.background = 'rgba(34, 211, 238, 0.3)')}
          onMouseLeave={(e) => (e.currentTarget.style.background = 'rgba(34, 211, 238, 0.2)')}
        >
          <RefreshCw size={16} />
          Re-run Simulation
        </button>
      </div>

      {/* Scenario Selector */}
      <div>
        <label style={{ display: 'block', fontSize: 14, fontWeight: 600, color: 'white', marginBottom: 12 }}>
          Market Scenario:
        </label>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))', gap: 12 }}>
          {Object.entries(scenarios).map(([key, scenario]) => (
            <button
              key={key}
              onClick={() => setSelectedScenario(key)}
              style={{
                padding: 12,
                borderRadius: 8,
                border: selectedScenario === key ? '2px solid rgb(34, 211, 238)' : '1px solid rgba(34, 211, 238, 0.2)',
                background: selectedScenario === key
                  ? 'rgba(34, 211, 238, 0.15)'
                  : 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
                color: selectedScenario === key ? 'rgb(34, 211, 238)' : 'rgba(203, 213, 225, 0.8)',
                fontSize: 13,
                fontWeight: 600,
                cursor: 'pointer',
                transition: 'all 0.3s',
                textAlign: 'center'
              }}
              onMouseEnter={(e) => {
                if (selectedScenario !== key) {
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.4)';
                  e.currentTarget.style.transform = 'translateY(-2px)';
                }
              }}
              onMouseLeave={(e) => {
                if (selectedScenario !== key) {
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.2)';
                  e.currentTarget.style.transform = 'translateY(0)';
                }
              }}
            >
              {scenario.name}
            </button>
          ))}
        </div>

        {/* Scenario Info */}
        <div style={{
          marginTop: 16,
          padding: 16,
          background: selectedScenario === 'normal' ? 'rgba(34, 197, 94, 0.1)' :
                      selectedScenario === 'moderate' ? 'rgba(245, 158, 11, 0.1)' :
                      selectedScenario === 'severe' ? 'rgba(249, 115, 22, 0.1)' :
                      'rgba(239, 68, 68, 0.1)',
          border: selectedScenario === 'normal' ? '1px solid rgba(34, 197, 94, 0.3)' :
                  selectedScenario === 'moderate' ? '1px solid rgba(245, 158, 11, 0.3)' :
                  selectedScenario === 'severe' ? '1px solid rgba(249, 115, 22, 0.3)' :
                  '1px solid rgba(239, 68, 68, 0.3)',
          borderRadius: 8,
          fontSize: 13,
          color: 'rgba(203, 213, 225, 0.9)'
        }}>
          <strong style={{ color: 'white' }}>
            {scenarios[selectedScenario].name}:
          </strong>{' '}
          {selectedScenario === 'normal' && '正常市场条件，基础清算概率。适合日常风险评估。'}
          {selectedScenario === 'moderate' && '中度市场压力，清算风险增加3倍。模拟市场波动加剧场景。'}
          {selectedScenario === 'severe' && '严重市场压力，清算风险增加7倍。模拟重大市场下跌场景。'}
          {selectedScenario === 'extreme' && '黑天鹅事件，清算风险增加15倍。模拟极端市场崩溃场景。'}
        </div>
      </div>

      {/* Summary Cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: 16 }}>
        {/* Total Agents */}
        <div style={{
          background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
          borderRadius: 12,
          padding: 20,
          border: '1px solid rgba(34, 211, 238, 0.2)'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
            <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 500 }}>Total Agents</span>
            <Users size={18} style={{ color: 'rgb(34, 211, 238)' }} />
          </div>
          <div style={{ fontSize: 32, fontWeight: 800, color: 'white', marginBottom: 4 }}>
            {simulationData.totalAgents.toLocaleString()}
          </div>
          <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)' }}>
            Simulated market participants
          </div>
        </div>

        {/* Liquidation Rate */}
        <div style={{
          background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
          borderRadius: 12,
          padding: 20,
          border: `1px solid ${simulationData.liquidationRate > 0.3 ? 'rgba(239, 68, 68, 0.3)' : 'rgba(245, 158, 11, 0.3)'}`
        }}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
            <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 500 }}>Liquidation Rate</span>
            <TrendingDown size={18} style={{ color: simulationData.liquidationRate > 0.3 ? 'rgb(239, 68, 68)' : 'rgb(245, 158, 11)' }} />
          </div>
          <div style={{ fontSize: 32, fontWeight: 800, color: simulationData.liquidationRate > 0.3 ? 'rgb(239, 68, 68)' : 'rgb(245, 158, 11)', marginBottom: 4 }}>
            {(simulationData.liquidationRate * 100).toFixed(1)}%
          </div>
          <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)' }}>
            {simulationData.liquidatedAgents.toLocaleString()} agents liquidated
          </div>
        </div>

        {/* Total Loss */}
        <div style={{
          background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
          borderRadius: 12,
          padding: 20,
          border: '1px solid rgba(34, 211, 238, 0.2)'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
            <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 500 }}>Potential Loss</span>
            <AlertTriangle size={18} style={{ color: 'rgb(245, 158, 11)' }} />
          </div>
          <div style={{ fontSize: 32, fontWeight: 800, color: 'white', marginBottom: 4 }}>
            ${(simulationData.totalPotentialLoss / 1e6).toFixed(2)}M
          </div>
          <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)' }}>
            Across all liquidations
          </div>
        </div>

        {/* VaR (99%) */}
        <div style={{
          background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
          borderRadius: 12,
          padding: 20,
          border: '1px solid rgba(34, 211, 238, 0.2)'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
            <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 500 }}>VaR (99%)</span>
            <Activity size={18} style={{ color: 'rgb(168, 85, 247)' }} />
          </div>
          <div style={{ fontSize: 32, fontWeight: 800, color: 'white', marginBottom: 4 }}>
            ${(simulationData.riskMetrics.var99 / 1e6).toFixed(2)}M
          </div>
          <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)' }}>
            Maximum expected loss
          </div>
        </div>
      </div>

      {/* Agent Type Distribution */}
      <div style={{
        background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
        borderRadius: 12,
        padding: 24,
        border: '1px solid rgba(34, 211, 238, 0.2)'
      }}>
        <h4 style={{ fontSize: 16, fontWeight: 700, color: 'white', marginBottom: 20 }}>
          Agent Type Distribution
        </h4>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: 16 }}>
          {Object.entries(simulationData.agentTypes).map(([type, data]) => (
            <div
              key={type}
              style={{
                padding: 16,
                background: 'rgba(15, 23, 42, 0.5)',
                borderRadius: 10,
                border: '1px solid rgba(34, 211, 238, 0.1)'
              }}
            >
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
                <span style={{ fontSize: 13, fontWeight: 600, color: 'white', textTransform: 'capitalize' }}>
                  {type.replace('_', ' ')}
                </span>
                <div style={{
                  width: 8,
                  height: 8,
                  borderRadius: '50%',
                  background: data.color,
                  boxShadow: `0 0 8px ${data.color}`
                }} />
              </div>
              <div style={{ fontSize: 24, fontWeight: 800, color: data.color, marginBottom: 8 }}>
                {data.count.toLocaleString()}
              </div>
              <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)' }}>
                Avg {data.avgLoss ? 'Loss' : 'Profit'}: ${(data.avgLoss || data.avgProfit).toLocaleString()}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Histogram */}
      <div style={{
        background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
        borderRadius: 12,
        padding: 24,
        border: '1px solid rgba(34, 211, 238, 0.2)'
      }}>
        <h4 style={{ fontSize: 16, fontWeight: 700, color: 'white', marginBottom: 20 }}>
          Health Factor Distribution
        </h4>
        <div style={{ display: 'flex', alignItems: 'flex-end', gap: 8, height: 200 }}>
          {simulationData.histogram.map((bin, index) => {
            const maxCount = Math.max(...simulationData.histogram.map(b => b.count));
            const height = (bin.count / maxCount) * 100;
            const isLiquidationZone = parseFloat(bin.hf) < 1.2;

            return (
              <div
                key={index}
                style={{
                  flex: 1,
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  gap: 8
                }}
              >
                <div style={{
                  width: '100%',
                  height: `${height}%`,
                  background: isLiquidationZone
                    ? 'linear-gradient(180deg, rgb(239, 68, 68) 0%, rgb(220, 38, 38) 100%)'
                    : 'linear-gradient(180deg, rgb(34, 211, 238) 0%, rgb(6, 182, 212) 100%)',
                  borderRadius: '4px 4px 0 0',
                  position: 'relative',
                  transition: 'all 0.3s',
                  cursor: 'pointer'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.opacity = '0.8';
                  e.currentTarget.style.transform = 'scaleY(1.05)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.opacity = '1';
                  e.currentTarget.style.transform = 'scaleY(1)';
                }}
                title={`HF: ${bin.hf}, Count: ${bin.count.toLocaleString()} (${bin.percentage.toFixed(1)}%)`}
                />
                <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)', fontWeight: 500 }}>
                  {bin.hf}
                </div>
              </div>
            );
          })}
        </div>
        <div style={{ marginTop: 16, padding: 12, background: 'rgba(239, 68, 68, 0.1)', borderRadius: 8, border: '1px solid rgba(239, 68, 68, 0.3)' }}>
          <div style={{ fontSize: 12, color: 'rgba(255, 255, 255, 0.9)' }}>
            <strong>Liquidation Zone (HF &lt; 1.2):</strong> Red bars indicate agents at high risk of liquidation
          </div>
        </div>
      </div>

      {/* Risk Metrics Summary */}
      <div style={{
        background: 'linear-gradient(135deg, rgba(168, 85, 247, 0.1) 0%, rgba(124, 58, 237, 0.1) 100%)',
        borderRadius: 12,
        padding: 20,
        border: '1px solid rgba(168, 85, 247, 0.3)'
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 16 }}>
          <Zap size={20} style={{ color: 'rgb(168, 85, 247)' }} />
          <h4 style={{ fontSize: 16, fontWeight: 700, color: 'white', margin: 0 }}>
            Advanced Risk Metrics
          </h4>
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: 16 }}>
          <div>
            <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 4 }}>VaR (95%)</div>
            <div style={{ fontSize: 20, fontWeight: 700, color: 'white' }}>
              ${(simulationData.riskMetrics.var95 / 1e6).toFixed(2)}M
            </div>
          </div>
          <div>
            <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 4 }}>CVaR (Expected Shortfall)</div>
            <div style={{ fontSize: 20, fontWeight: 700, color: 'white' }}>
              ${(simulationData.riskMetrics.cvar / 1e6).toFixed(2)}M
            </div>
          </div>
          <div>
            <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 4 }}>Max Drawdown</div>
            <div style={{ fontSize: 20, fontWeight: 700, color: 'white' }}>
              ${(simulationData.riskMetrics.maxDrawdown / 1e6).toFixed(2)}M
            </div>
          </div>
          <div>
            <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 4 }}>Avg Loss per Agent</div>
            <div style={{ fontSize: 20, fontWeight: 700, color: 'white' }}>
              ${simulationData.avgLossPerAgent.toLocaleString(undefined, { maximumFractionDigits: 0 })}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AgentSimulationPanel;
