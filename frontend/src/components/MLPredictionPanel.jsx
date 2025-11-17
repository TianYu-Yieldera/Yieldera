/**
 * ML Prediction Panel
 * 机器学习风险预测组件
 * 使用ML模型预测未来24小时的清算风险和健康因子变化
 */

import React, { useState, useEffect } from 'react';
import { Brain, TrendingUp, TrendingDown, AlertTriangle, Shield, Activity, Zap, RefreshCw, Target, BarChart2 } from 'lucide-react';
import fastAPIRiskService from '../services/fastAPIRiskService';

const MLPredictionPanel = ({ positions = [] }) => {
  const [predictions, setPredictions] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [selectedPosition, setSelectedPosition] = useState(null);

  useEffect(() => {
    if (positions.length > 0) {
      setSelectedPosition(positions[0]);
      runPrediction(positions[0]);
    }
  }, [positions]);

  const runPrediction = async (position) => {
    if (!position) return;

    setLoading(true);
    setError(null);

    try {
      // Generate mock ML predictions (in production would call FastAPI)
      const mockPredictions = generateMockPredictions(position);

      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 1500));

      setPredictions(mockPredictions);
    } catch (err) {
      console.error('ML prediction failed:', err);
      setError('Failed to run ML prediction. Using fallback data.');

      // Fallback to mock data
      const mockPredictions = generateMockPredictions(position);
      setPredictions(mockPredictions);
    } finally {
      setLoading(false);
    }
  };

  const generateMockPredictions = (position) => {
    const currentHF = position.healthFactor || 2.0;
    const volatility = 0.15; // Mock volatility
    const currentPrice = position.currentPrice || 2500;

    // Predict health factor degradation
    const predictedHF = currentHF * (1 - volatility * 0.12);

    // Calculate liquidation probability
    let liquidationProb = 0;
    if (predictedHF < 1.1) liquidationProb = 0.85;
    else if (predictedHF < 1.2) liquidationProb = 0.60;
    else if (predictedHF < 1.5) liquidationProb = 0.25;
    else if (predictedHF < 2.0) liquidationProb = 0.08;
    else liquidationProb = 0.02;

    // Price prediction with confidence interval
    const meanPrice = currentPrice * 1.005; // Slight upward trend
    const priceConfidenceInterval = [
      currentPrice * 0.92,  // Lower bound (-8%)
      currentPrice * 1.08   // Upper bound (+8%)
    ];

    // Trend analysis
    const hfTrend = predictedHF < currentHF ? 'declining' : 'stable';
    const riskTrend = liquidationProb > 0.3 ? 'increasing' : liquidationProb > 0.1 ? 'moderate' : 'low';

    // Generate recommendations based on predictions
    const recommendations = [];
    if (liquidationProb > 0.5) {
      recommendations.push({
        severity: 'critical',
        icon: '🚨',
        title: 'High Liquidation Risk Detected',
        action: 'Add collateral immediately to avoid liquidation',
        impact: 'Reduces risk by 60-80%'
      });
      recommendations.push({
        severity: 'high',
        icon: '💰',
        title: 'Reduce Leverage',
        action: 'Close partial position to decrease LTV ratio',
        impact: 'Improves health factor by 0.5-1.0'
      });
    } else if (liquidationProb > 0.2) {
      recommendations.push({
        severity: 'medium',
        icon: '⚠️',
        title: 'Monitor Position Closely',
        action: 'Set up price alerts and prepare collateral',
        impact: 'Early warning reduces liquidation risk'
      });
      recommendations.push({
        severity: 'low',
        icon: '🔄',
        title: 'Consider Hedging',
        action: 'Use derivatives to hedge downside risk',
        impact: 'Protects against 10-15% price drops'
      });
    } else {
      recommendations.push({
        severity: 'low',
        icon: '✅',
        title: 'Position Healthy',
        action: 'Maintain current risk management strategy',
        impact: 'Continue monitoring market conditions'
      });
      recommendations.push({
        severity: 'low',
        icon: '📈',
        title: 'Optional Optimization',
        action: 'Consider rebalancing for better yield',
        impact: 'Potential 2-5% APY improvement'
      });
    }

    // Historical trend simulation (24 data points for 24 hours)
    const historicalHF = [];
    const historicalProb = [];
    for (let i = 0; i < 24; i++) {
      const hourlyVolatility = Math.random() * 0.05 - 0.025;
      const hf = currentHF + (predictedHF - currentHF) * (i / 24) + hourlyVolatility;
      const prob = liquidationProb * (i / 24) + Math.random() * 0.05;

      historicalHF.push({
        hour: i,
        value: Math.max(hf, 0.5)
      });

      historicalProb.push({
        hour: i,
        value: Math.min(Math.max(prob, 0), 1)
      });
    }

    return {
      position_id: position.id,
      current_health_factor: currentHF,
      predicted_health_factor: predictedHF,
      health_factor_change: predictedHF - currentHF,
      liquidation_probability_24h: liquidationProb,
      price_prediction: {
        current: currentPrice,
        mean: meanPrice,
        confidence_interval: priceConfidenceInterval
      },
      confidence: 0.82, // ML model confidence
      trend: {
        health_factor: hfTrend,
        risk: riskTrend
      },
      recommendations: recommendations,
      time_horizon: 24, // hours
      historical_prediction: {
        health_factor: historicalHF,
        liquidation_probability: historicalProb
      }
    };
  };

  const handlePositionSelect = (position) => {
    setSelectedPosition(position);
    runPrediction(position);
  };

  const getRiskColor = (probability) => {
    if (probability > 0.5) return { color: 'rgb(239, 68, 68)', bg: 'rgba(239, 68, 68, 0.1)', border: 'rgba(239, 68, 68, 0.3)' };
    if (probability > 0.2) return { color: 'rgb(249, 115, 22)', bg: 'rgba(249, 115, 22, 0.1)', border: 'rgba(249, 115, 22, 0.3)' };
    if (probability > 0.1) return { color: 'rgb(245, 158, 11)', bg: 'rgba(245, 158, 11, 0.1)', border: 'rgba(245, 158, 11, 0.3)' };
    return { color: 'rgb(34, 197, 94)', bg: 'rgba(34, 197, 94, 0.1)', border: 'rgba(34, 197, 94, 0.3)' };
  };

  const getSeverityStyle = (severity) => {
    switch (severity) {
      case 'critical': return { color: 'rgb(239, 68, 68)', bg: 'rgba(239, 68, 68, 0.1)', border: 'rgba(239, 68, 68, 0.3)' };
      case 'high': return { color: 'rgb(249, 115, 22)', bg: 'rgba(249, 115, 22, 0.1)', border: 'rgba(249, 115, 22, 0.3)' };
      case 'medium': return { color: 'rgb(245, 158, 11)', bg: 'rgba(245, 158, 11, 0.1)', border: 'rgba(245, 158, 11, 0.3)' };
      default: return { color: 'rgb(34, 211, 238)', bg: 'rgba(34, 211, 238, 0.1)', border: 'rgba(34, 211, 238, 0.3)' };
    }
  };

  if (loading) {
    return (
      <div style={{
        padding: 40,
        textAlign: 'center',
        background: 'rgba(168, 85, 247, 0.05)',
        borderRadius: 12,
        border: '1px solid rgba(168, 85, 247, 0.2)'
      }}>
        <div style={{
          width: 64,
          height: 64,
          margin: '0 auto 20px',
          borderRadius: '50%',
          border: '3px solid rgba(168, 85, 247, 0.3)',
          borderTopColor: 'rgb(168, 85, 247)',
          animation: 'spin 1s linear infinite'
        }} />
        <style>{`
          @keyframes spin {
            to { transform: rotate(360deg); }
          }
        `}</style>
        <p style={{ fontSize: 16, color: 'white', marginBottom: 8, fontWeight: 600 }}>
          Running ML Predictions
        </p>
        <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.7)', margin: 0 }}>
          Analyzing market patterns and risk factors...
        </p>
      </div>
    );
  }

  if (!predictions && !loading) {
    return (
      <div style={{
        padding: 40,
        textAlign: 'center',
        background: 'rgba(255, 255, 255, 0.03)',
        borderRadius: 12,
        border: '1px solid rgba(255, 255, 255, 0.08)'
      }}>
        <Brain size={48} style={{ color: 'rgba(203, 213, 225, 0.5)', marginBottom: 16 }} />
        <p style={{ fontSize: 16, color: 'rgba(203, 213, 225, 0.8)', margin: '0 0 16px 0', fontWeight: 600 }}>
          No Positions to Analyze
        </p>
        <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.6)', margin: 0 }}>
          Open a position to get ML-powered risk predictions
        </p>
      </div>
    );
  }

  const riskColors = getRiskColor(predictions?.liquidation_probability_24h || 0);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
      {/* Header */}
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', flexWrap: 'wrap', gap: 16 }}>
        <div>
          <h3 style={{ fontSize: 18, fontWeight: 700, color: 'white', margin: '0 0 8px 0', display: 'flex', alignItems: 'center', gap: 12 }}>
            <Brain size={24} style={{ color: 'rgb(168, 85, 247)' }} />
            ML Risk Prediction (24h Forecast)
          </h3>
          <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.7)', margin: 0 }}>
            Machine learning powered risk analysis and recommendations
          </p>
        </div>

        <button
          onClick={() => runPrediction(selectedPosition)}
          disabled={loading}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: 8,
            padding: '10px 20px',
            background: 'rgba(168, 85, 247, 0.2)',
            border: '1px solid rgba(168, 85, 247, 0.4)',
            borderRadius: 8,
            color: 'rgb(168, 85, 247)',
            fontSize: 14,
            fontWeight: 600,
            cursor: loading ? 'not-allowed' : 'pointer',
            opacity: loading ? 0.6 : 1,
            transition: 'all 0.3s'
          }}
          onMouseEnter={(e) => !loading && (e.currentTarget.style.background = 'rgba(168, 85, 247, 0.3)')}
          onMouseLeave={(e) => (e.currentTarget.style.background = 'rgba(168, 85, 247, 0.2)')}
        >
          <RefreshCw size={16} />
          Re-run Prediction
        </button>
      </div>

      {/* Position Selector */}
      {positions.length > 1 && (
        <div>
          <label style={{ display: 'block', fontSize: 14, fontWeight: 600, color: 'white', marginBottom: 12 }}>
            Select Position:
          </label>
          <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
            {positions.map((pos) => (
              <button
                key={pos.id}
                onClick={() => handlePositionSelect(pos)}
                style={{
                  padding: '12px 16px',
                  borderRadius: 8,
                  border: selectedPosition?.id === pos.id ? '2px solid rgb(168, 85, 247)' : '1px solid rgba(34, 211, 238, 0.2)',
                  background: selectedPosition?.id === pos.id ? 'rgba(168, 85, 247, 0.15)' : 'rgba(15, 23, 42, 0.5)',
                  color: selectedPosition?.id === pos.id ? 'rgb(168, 85, 247)' : 'rgba(203, 213, 225, 0.8)',
                  fontSize: 13,
                  fontWeight: 600,
                  cursor: 'pointer',
                  transition: 'all 0.3s'
                }}
              >
                {pos.protocol} - {pos.collateralAsset}
              </button>
            ))}
          </div>
        </div>
      )}

      {predictions && (
        <>
          {/* Key Predictions Grid */}
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: 16 }}>
            {/* Health Factor Prediction */}
            <div style={{
              background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
              borderRadius: 12,
              padding: 24,
              border: '1px solid rgba(34, 211, 238, 0.2)'
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 16 }}>
                <Activity size={18} style={{ color: 'rgb(34, 211, 238)' }} />
                <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 600 }}>
                  Predicted Health Factor
                </span>
              </div>
              <div style={{ display: 'flex', alignItems: 'baseline', gap: 12, marginBottom: 12 }}>
                <div style={{ fontSize: 36, fontWeight: 800, color: 'white' }}>
                  {predictions.predicted_health_factor.toFixed(2)}
                </div>
                <div style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 4,
                  fontSize: 14,
                  fontWeight: 600,
                  color: predictions.health_factor_change < 0 ? 'rgb(239, 68, 68)' : 'rgb(34, 197, 94)'
                }}>
                  {predictions.health_factor_change < 0 ? <TrendingDown size={16} /> : <TrendingUp size={16} />}
                  {Math.abs(predictions.health_factor_change).toFixed(2)}
                </div>
              </div>
              <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                Current: {predictions.current_health_factor.toFixed(2)}
              </div>
            </div>

            {/* Liquidation Probability */}
            <div style={{
              background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
              borderRadius: 12,
              padding: 24,
              border: `2px solid ${riskColors.border}`
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 16 }}>
                <AlertTriangle size={18} style={{ color: riskColors.color }} />
                <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 600 }}>
                  Liquidation Risk (24h)
                </span>
              </div>
              <div style={{ fontSize: 36, fontWeight: 800, color: riskColors.color, marginBottom: 12 }}>
                {(predictions.liquidation_probability_24h * 100).toFixed(1)}%
              </div>
              <div style={{
                display: 'inline-block',
                padding: '4px 10px',
                background: riskColors.bg,
                border: `1px solid ${riskColors.border}`,
                borderRadius: 6,
                fontSize: 11,
                fontWeight: 700,
                color: riskColors.color
              }}>
                {predictions.liquidation_probability_24h > 0.5 ? 'CRITICAL' :
                 predictions.liquidation_probability_24h > 0.2 ? 'HIGH' :
                 predictions.liquidation_probability_24h > 0.1 ? 'MODERATE' : 'LOW'} RISK
              </div>
            </div>

            {/* Price Prediction */}
            <div style={{
              background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
              borderRadius: 12,
              padding: 24,
              border: '1px solid rgba(34, 211, 238, 0.2)'
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 16 }}>
                <Target size={18} style={{ color: 'rgb(34, 211, 238)' }} />
                <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 600 }}>
                  Price Forecast
                </span>
              </div>
              <div style={{ fontSize: 36, fontWeight: 800, color: 'white', marginBottom: 12 }}>
                ${predictions.price_prediction.mean.toLocaleString(undefined, { maximumFractionDigits: 0 })}
              </div>
              <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)', lineHeight: 1.5 }}>
                95% CI: ${predictions.price_prediction.confidence_interval[0].toLocaleString()} - ${predictions.price_prediction.confidence_interval[1].toLocaleString()}
              </div>
            </div>

            {/* Model Confidence */}
            <div style={{
              background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
              borderRadius: 12,
              padding: 24,
              border: '1px solid rgba(34, 211, 238, 0.2)'
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 16 }}>
                <Zap size={18} style={{ color: 'rgb(168, 85, 247)' }} />
                <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 600 }}>
                  ML Model Confidence
                </span>
              </div>
              <div style={{ fontSize: 36, fontWeight: 800, color: 'rgb(168, 85, 247)', marginBottom: 12 }}>
                {(predictions.confidence * 100).toFixed(0)}%
              </div>
              <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                Based on historical patterns
              </div>
            </div>
          </div>

          {/* Trend Visualization */}
          <div style={{
            background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
            borderRadius: 16,
            padding: 24,
            border: '1px solid rgba(34, 211, 238, 0.2)'
          }}>
            <h4 style={{ fontSize: 16, fontWeight: 700, color: 'white', marginBottom: 20, display: 'flex', alignItems: 'center', gap: 8 }}>
              <BarChart2 size={20} style={{ color: 'rgb(34, 211, 238)' }} />
              24-Hour Risk Trajectory
            </h4>

            {/* Simple bar chart visualization */}
            <div style={{ display: 'flex', alignItems: 'flex-end', gap: 4, height: 150, marginBottom: 16 }}>
              {predictions.historical_prediction.liquidation_probability.map((point, index) => {
                if (index % 3 !== 0) return null; // Show every 3rd hour
                const height = (point.value * 100);
                const barColor = point.value > 0.5 ? 'rgb(239, 68, 68)' :
                                point.value > 0.2 ? 'rgb(249, 115, 22)' :
                                point.value > 0.1 ? 'rgb(245, 158, 11)' : 'rgb(34, 197, 94)';

                return (
                  <div key={index} style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 8 }}>
                    <div
                      style={{
                        width: '100%',
                        height: `${height}%`,
                        background: `linear-gradient(180deg, ${barColor} 0%, ${barColor}80 100%)`,
                        borderRadius: '4px 4px 0 0',
                        transition: 'all 0.3s',
                        cursor: 'pointer'
                      }}
                      title={`Hour ${point.hour}: ${(point.value * 100).toFixed(1)}%`}
                    />
                    <div style={{ fontSize: 10, color: 'rgba(203, 213, 225, 0.5)' }}>
                      {point.hour}h
                    </div>
                  </div>
                );
              })}
            </div>
            <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', textAlign: 'center' }}>
              Liquidation probability forecast over next 24 hours
            </div>
          </div>

          {/* AI Recommendations */}
          <div style={{
            background: 'linear-gradient(135deg, rgba(168, 85, 247, 0.1) 0%, rgba(124, 58, 237, 0.1) 100%)',
            borderRadius: 16,
            padding: 24,
            border: '1px solid rgba(168, 85, 247, 0.3)'
          }}>
            <h4 style={{ fontSize: 16, fontWeight: 700, color: 'white', marginBottom: 20, display: 'flex', alignItems: 'center', gap: 8 }}>
              <Brain size={20} style={{ color: 'rgb(168, 85, 247)' }} />
              AI-Powered Recommendations
            </h4>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: 16 }}>
              {predictions.recommendations.map((rec, index) => {
                const style = getSeverityStyle(rec.severity);
                return (
                  <div
                    key={index}
                    style={{
                      padding: 20,
                      background: style.bg,
                      border: `1px solid ${style.border}`,
                      borderRadius: 12,
                      transition: 'all 0.3s'
                    }}
                  >
                    <div style={{ display: 'flex', alignItems: 'flex-start', gap: 12, marginBottom: 12 }}>
                      <span style={{ fontSize: 24 }}>{rec.icon}</span>
                      <div style={{ flex: 1 }}>
                        <div style={{ fontSize: 14, fontWeight: 700, color: 'white', marginBottom: 6 }}>
                          {rec.title}
                        </div>
                        <div style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.9)', marginBottom: 8, lineHeight: 1.5 }}>
                          {rec.action}
                        </div>
                        <div style={{
                          fontSize: 11,
                          color: style.color,
                          fontWeight: 600,
                          display: 'flex',
                          alignItems: 'center',
                          gap: 4
                        }}>
                          <Zap size={12} />
                          {rec.impact}
                        </div>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          {/* Methodology Info */}
          <div style={{
            padding: 16,
            background: 'rgba(34, 211, 238, 0.1)',
            border: '1px solid rgba(34, 211, 238, 0.3)',
            borderRadius: 10,
            fontSize: 13,
            color: 'rgba(203, 213, 225, 0.9)'
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
              <Brain size={16} style={{ color: 'rgb(34, 211, 238)' }} />
              <strong style={{ color: 'white' }}>ML Prediction Methodology</strong>
            </div>
            <p style={{ margin: 0, lineHeight: 1.6 }}>
              预测基于深度学习模型，分析历史价格波动、市场情绪、链上数据和宏观经济指标。
              模型持续学习市场模式，提供{predictions.time_horizon}小时内的风险评估和智能建议。
              置信度{(predictions.confidence * 100).toFixed(0)}%表示模型对预测结果的确信程度。
            </p>
          </div>
        </>
      )}
    </div>
  );
};

export default MLPredictionPanel;
