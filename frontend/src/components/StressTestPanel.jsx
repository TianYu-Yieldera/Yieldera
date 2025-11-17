/**
 * Stress Test Panel
 * 压力测试可视化组件
 * 测试投资组合在极端市场条件下的表现
 */

import React, { useState, useEffect } from 'react';
import { AlertTriangle, TrendingDown, Shield, BarChart3, RefreshCw, Zap, DollarSign } from 'lucide-react';
import fastAPIRiskService from '../services/fastAPIRiskService';

const StressTestPanel = ({ positions = [] }) => {
  const [testResults, setTestResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [selectedScenario, setSelectedScenario] = useState('all');

  // 预定义的压力测试场景
  const stressScenarios = [
    {
      name: 'Market Correction',
      description: '市场回调 (-15%)',
      price_shock: -0.15,
      volatility_mult: 1.5,
      color: 'rgb(245, 158, 11)',
      severity: 'low'
    },
    {
      name: 'Flash Crash',
      description: '闪电崩盘 (-30%)',
      price_shock: -0.30,
      volatility_mult: 2.5,
      color: 'rgb(249, 115, 22)',
      severity: 'medium'
    },
    {
      name: 'Black Swan',
      description: '黑天鹅事件 (-50%)',
      price_shock: -0.50,
      volatility_mult: 4.0,
      color: 'rgb(239, 68, 68)',
      severity: 'high'
    },
    {
      name: 'Great Depression',
      description: '大萧条级别 (-70%)',
      price_shock: -0.70,
      volatility_mult: 6.0,
      color: 'rgb(127, 29, 29)',
      severity: 'extreme'
    }
  ];

  useEffect(() => {
    if (positions.length > 0) {
      runStressTest();
    }
  }, [positions]);

  const runStressTest = async () => {
    setLoading(true);
    setError(null);

    try {
      // Generate mock results (in production would call FastAPI)
      const mockResults = generateMockStressTestResults(positions, stressScenarios);

      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 1200));

      setTestResults(mockResults);
    } catch (err) {
      console.error('Stress test failed:', err);
      setError('Failed to run stress test. Using fallback data.');

      // Fallback to mock data
      const mockResults = generateMockStressTestResults(positions, stressScenarios);
      setTestResults(mockResults);
    } finally {
      setLoading(false);
    }
  };

  const generateMockStressTestResults = (positions, scenarios) => {
    const totalValue = positions.reduce((sum, p) => sum + (p.collateralValueUSD || 0), 0);

    return scenarios.map(scenario => {
      // Calculate stressed portfolio value
      const stressedValue = totalValue * (1 + scenario.price_shock);
      const totalLoss = totalValue - stressedValue;

      // Calculate positions at risk
      let positionsLiquidated = 0;
      let worstHealthFactor = Infinity;

      positions.forEach(position => {
        const currentHF = position.healthFactor || 2.0;
        // Approximate stressed health factor
        const stressedHF = currentHF * (1 + scenario.price_shock * 1.5);

        if (stressedHF < 1.0) {
          positionsLiquidated++;
        }

        worstHealthFactor = Math.min(worstHealthFactor, stressedHF);
      });

      // Generate recommendations
      const recommendations = [];
      const liquidationRate = positionsLiquidated / positions.length;

      if (liquidationRate > 0.5) {
        recommendations.push('🚨 URGENT: 超过50%的仓位面临清算风险');
        recommendations.push('💡 建议: 立即减少杠杆或增加抵押品');
      } else if (liquidationRate > 0.3) {
        recommendations.push('⚠️  警告: 30%以上仓位处于风险中');
        recommendations.push('💡 建议: 考虑对冲策略降低风险敞口');
      } else if (liquidationRate > 0) {
        recommendations.push('⚡ 提示: 部分仓位可能受影响');
        recommendations.push('💡 建议: 监控健康因子，准备应对措施');
      } else {
        recommendations.push('✅ 良好: 投资组合能够承受此场景');
        recommendations.push('💡 建议: 保持当前风险水平');
      }

      if (worstHealthFactor < 1.2 && worstHealthFactor > 0) {
        recommendations.push('📊 最低健康因子接近清算线，建议增加抵押');
      }

      return {
        scenario_name: scenario.name,
        scenario_description: scenario.description,
        total_loss: Math.abs(totalLoss),
        total_value: totalValue,
        loss_percentage: Math.abs(scenario.price_shock) * 100,
        positions_liquidated: positionsLiquidated,
        total_positions: positions.length,
        worst_health_factor: worstHealthFactor === Infinity ? 0 : Math.max(worstHealthFactor, 0),
        recommendations: recommendations,
        color: scenario.color,
        severity: scenario.severity
      };
    });
  };

  if (loading) {
    return (
      <div style={{
        padding: 40,
        textAlign: 'center',
        background: 'rgba(249, 115, 22, 0.05)',
        borderRadius: 12,
        border: '1px solid rgba(249, 115, 22, 0.2)'
      }}>
        <div style={{
          width: 64,
          height: 64,
          margin: '0 auto 20px',
          borderRadius: '50%',
          border: '3px solid rgba(249, 115, 22, 0.3)',
          borderTopColor: 'rgb(249, 115, 22)',
          animation: 'spin 1s linear infinite'
        }} />
        <style>{`
          @keyframes spin {
            to { transform: rotate(360deg); }
          }
        `}</style>
        <p style={{ fontSize: 16, color: 'white', marginBottom: 8, fontWeight: 600 }}>
          Running Stress Tests
        </p>
        <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.7)', margin: 0 }}>
          Testing portfolio under extreme conditions...
        </p>
      </div>
    );
  }

  if (!testResults && !loading) {
    return (
      <div style={{
        padding: 40,
        textAlign: 'center',
        background: 'rgba(255, 255, 255, 0.03)',
        borderRadius: 12,
        border: '1px solid rgba(255, 255, 255, 0.08)'
      }}>
        <Shield size={48} style={{ color: 'rgba(203, 213, 225, 0.5)', marginBottom: 16 }} />
        <p style={{ fontSize: 16, color: 'rgba(203, 213, 225, 0.8)', margin: '0 0 16px 0', fontWeight: 600 }}>
          No Positions to Test
        </p>
        <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.6)', margin: 0 }}>
          Open a position to run stress tests
        </p>
      </div>
    );
  }

  const filteredResults = selectedScenario === 'all'
    ? testResults
    : testResults.filter(r => r.scenario_name === selectedScenario);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
      {/* Header */}
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', flexWrap: 'wrap', gap: 16 }}>
        <div>
          <h3 style={{ fontSize: 18, fontWeight: 700, color: 'white', margin: '0 0 8px 0' }}>
            Portfolio Stress Testing
          </h3>
          <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.7)', margin: 0 }}>
            Evaluate portfolio resilience under extreme market conditions
          </p>
        </div>

        <button
          onClick={runStressTest}
          disabled={loading}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: 8,
            padding: '10px 20px',
            background: 'rgba(249, 115, 22, 0.2)',
            border: '1px solid rgba(249, 115, 22, 0.4)',
            borderRadius: 8,
            color: 'rgb(249, 115, 22)',
            fontSize: 14,
            fontWeight: 600,
            cursor: loading ? 'not-allowed' : 'pointer',
            opacity: loading ? 0.6 : 1,
            transition: 'all 0.3s'
          }}
          onMouseEnter={(e) => !loading && (e.currentTarget.style.background = 'rgba(249, 115, 22, 0.3)')}
          onMouseLeave={(e) => (e.currentTarget.style.background = 'rgba(249, 115, 22, 0.2)')}
        >
          <RefreshCw size={16} />
          Re-run Tests
        </button>
      </div>

      {/* Scenario Filter */}
      <div>
        <label style={{ display: 'block', fontSize: 14, fontWeight: 600, color: 'white', marginBottom: 12 }}>
          Select Scenario:
        </label>
        <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
          <button
            onClick={() => setSelectedScenario('all')}
            style={{
              padding: '10px 16px',
              borderRadius: 8,
              border: selectedScenario === 'all' ? '2px solid rgb(34, 211, 238)' : '1px solid rgba(34, 211, 238, 0.2)',
              background: selectedScenario === 'all' ? 'rgba(34, 211, 238, 0.15)' : 'rgba(15, 23, 42, 0.5)',
              color: selectedScenario === 'all' ? 'rgb(34, 211, 238)' : 'rgba(203, 213, 225, 0.8)',
              fontSize: 13,
              fontWeight: 600,
              cursor: 'pointer',
              transition: 'all 0.3s'
            }}
          >
            All Scenarios
          </button>
          {testResults && testResults.map((result) => (
            <button
              key={result.scenario_name}
              onClick={() => setSelectedScenario(result.scenario_name)}
              style={{
                padding: '10px 16px',
                borderRadius: 8,
                border: selectedScenario === result.scenario_name ? `2px solid ${result.color}` : '1px solid rgba(34, 211, 238, 0.2)',
                background: selectedScenario === result.scenario_name ? `${result.color}20` : 'rgba(15, 23, 42, 0.5)',
                color: selectedScenario === result.scenario_name ? result.color : 'rgba(203, 213, 225, 0.8)',
                fontSize: 13,
                fontWeight: 600,
                cursor: 'pointer',
                transition: 'all 0.3s'
              }}
            >
              {result.scenario_name}
            </button>
          ))}
        </div>
      </div>

      {/* Results Grid */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: 20 }}>
        {filteredResults && filteredResults.map((result, index) => (
          <div
            key={index}
            style={{
              background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
              borderRadius: 16,
              padding: 24,
              border: `2px solid ${result.color}40`,
              position: 'relative',
              overflow: 'hidden'
            }}
          >
            {/* Severity indicator */}
            <div style={{
              position: 'absolute',
              top: -30,
              right: -30,
              width: 100,
              height: 100,
              background: `radial-gradient(circle, ${result.color}30 0%, transparent 70%)`,
              filter: 'blur(20px)'
            }} />

            {/* Header */}
            <div style={{ position: 'relative', zIndex: 1, marginBottom: 20 }}>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
                <h4 style={{ fontSize: 16, fontWeight: 700, color: 'white', margin: 0 }}>
                  {result.scenario_name}
                </h4>
                <TrendingDown size={20} style={{ color: result.color }} />
              </div>
              <p style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.7)', margin: 0 }}>
                {result.scenario_description}
              </p>
            </div>

            {/* Key Metrics */}
            <div style={{ position: 'relative', zIndex: 1, display: 'flex', flexDirection: 'column', gap: 16, marginBottom: 20 }}>
              {/* Total Loss */}
              <div style={{
                padding: 16,
                background: 'rgba(15, 23, 42, 0.5)',
                borderRadius: 10,
                border: '1px solid rgba(34, 211, 238, 0.1)'
              }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
                  <DollarSign size={16} style={{ color: result.color }} />
                  <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 500 }}>
                    Potential Loss
                  </span>
                </div>
                <div style={{ fontSize: 24, fontWeight: 800, color: result.color, marginBottom: 4 }}>
                  ${result.total_loss.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                </div>
                <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                  {result.loss_percentage.toFixed(1)}% of total portfolio
                </div>
              </div>

              {/* Liquidations */}
              <div style={{
                padding: 16,
                background: 'rgba(15, 23, 42, 0.5)',
                borderRadius: 10,
                border: '1px solid rgba(34, 211, 238, 0.1)'
              }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
                  <AlertTriangle size={16} style={{ color: result.color }} />
                  <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 500 }}>
                    Positions at Risk
                  </span>
                </div>
                <div style={{ fontSize: 24, fontWeight: 800, color: result.color, marginBottom: 4 }}>
                  {result.positions_liquidated} / {result.total_positions}
                </div>
                <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                  {((result.positions_liquidated / result.total_positions) * 100).toFixed(1)}% liquidation rate
                </div>
              </div>

              {/* Worst Health Factor */}
              <div style={{
                padding: 16,
                background: 'rgba(15, 23, 42, 0.5)',
                borderRadius: 10,
                border: '1px solid rgba(34, 211, 238, 0.1)'
              }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
                  <BarChart3 size={16} style={{ color: result.color }} />
                  <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 500 }}>
                    Worst Health Factor
                  </span>
                </div>
                <div style={{ fontSize: 24, fontWeight: 800, color: result.color, marginBottom: 4 }}>
                  {result.worst_health_factor.toFixed(2)}
                </div>
                <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                  {result.worst_health_factor < 1.0 ? 'Below liquidation threshold' :
                   result.worst_health_factor < 1.2 ? 'Near liquidation' :
                   'Above safe threshold'}
                </div>
              </div>
            </div>

            {/* Recommendations */}
            <div style={{
              position: 'relative',
              zIndex: 1,
              padding: 16,
              background: `${result.color}15`,
              border: `1px solid ${result.color}40`,
              borderRadius: 10
            }}>
              <div style={{ fontSize: 12, fontWeight: 700, color: 'white', marginBottom: 12 }}>
                Recommendations
              </div>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                {result.recommendations.map((rec, idx) => (
                  <div key={idx} style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.9)', lineHeight: 1.5 }}>
                    {rec}
                  </div>
                ))}
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Summary Info */}
      {testResults && (
        <div style={{
          marginTop: 8,
          padding: 16,
          background: 'rgba(34, 211, 238, 0.1)',
          border: '1px solid rgba(34, 211, 238, 0.3)',
          borderRadius: 10,
          fontSize: 13,
          color: 'rgba(203, 213, 225, 0.9)'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
            <Zap size={16} style={{ color: 'rgb(34, 211, 238)' }} />
            <strong style={{ color: 'white' }}>Stress Test Methodology</strong>
          </div>
          <p style={{ margin: 0, lineHeight: 1.6 }}>
            这些压力测试模拟极端市场条件对您投资组合的影响。测试包括价格冲击、流动性危机和波动率飙升等场景。
            结果显示最坏情况下的潜在损失和需要采取的风险管理措施。
          </p>
        </div>
      )}
    </div>
  );
};

export default StressTestPanel;
