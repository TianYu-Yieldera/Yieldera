/**
 * Liquidation Alert Component
 * Real-time liquidation risk monitoring with FastAPI-powered VaR/CVaR calculations
 * Color-coded health factor display with one-click collateral deposit
 */

import React, { useState, useEffect, useCallback } from 'react';
import {
  AlertTriangle,
  Shield,
  TrendingDown,
  Plus,
  ChevronDown,
  ChevronUp,
  Activity,
  Zap
} from 'lucide-react';
import fastAPIRiskService from '../services/fastAPIRiskService';
import { useWallet } from '../web3/WalletContext';

const LiquidationAlert = ({
  position,
  marketData,
  historicalPrices = [],
  onAddCollateral,
  refreshInterval = 60000 // 60 seconds default
}) => {
  const { address } = useWallet();
  const [riskMetrics, setRiskMetrics] = useState(null);
  const [loading, setLoading] = useState(false);
  const [expanded, setExpanded] = useState(false);
  const [fastAPIAvailable, setFastAPIAvailable] = useState(false);
  const [lastUpdate, setLastUpdate] = useState(null);

  // Calculate alert level based on health factor
  const getAlertLevel = (healthFactor) => {
    if (healthFactor >= 2.0) return 'healthy';
    if (healthFactor >= 1.5) return 'warning';
    if (healthFactor >= 1.2) return 'danger';
    return 'critical';
  };

  // Get color scheme for alert level
  const getColorScheme = (level) => {
    const schemes = {
      healthy: {
        primary: 'rgb(34, 197, 94)',
        bg: 'rgba(34, 197, 94, 0.1)',
        border: 'rgba(34, 197, 94, 0.3)',
        text: 'rgb(34, 197, 94)'
      },
      warning: {
        primary: 'rgb(245, 158, 11)',
        bg: 'rgba(245, 158, 11, 0.1)',
        border: 'rgba(245, 158, 11, 0.3)',
        text: 'rgb(245, 158, 11)'
      },
      danger: {
        primary: 'rgb(249, 115, 22)',
        bg: 'rgba(249, 115, 22, 0.1)',
        border: 'rgba(249, 115, 22, 0.3)',
        text: 'rgb(249, 115, 22)'
      },
      critical: {
        primary: 'rgb(239, 68, 68)',
        bg: 'rgba(239, 68, 68, 0.1)',
        border: 'rgba(239, 68, 68, 0.3)',
        text: 'rgb(239, 68, 68)'
      }
    };
    return schemes[level] || schemes.healthy;
  };

  // Check FastAPI availability
  const checkFastAPIAvailability = useCallback(async () => {
    try {
      const health = await fastAPIRiskService.checkHealth();
      setFastAPIAvailable(health.available);
      return health.available;
    } catch (error) {
      console.warn('[LiquidationAlert] FastAPI check failed:', error);
      setFastAPIAvailable(false);
      return false;
    }
  }, []);

  // Calculate risk metrics
  const calculateRisk = useCallback(async () => {
    if (!position) return;

    setLoading(true);

    try {
      // Check if FastAPI is available
      const available = await checkFastAPIAvailability();

      if (available && historicalPrices.length > 0) {
        // Use FastAPI for advanced risk calculations
        const metrics = await fastAPIRiskService.calculateRisk(
          position,
          marketData || {},
          historicalPrices
        );

        setRiskMetrics({
          healthFactor: metrics.health_factor || position.healthFactor || 2.0,
          liquidationProb: metrics.liquidation_probability || 0,
          var99: metrics.value_at_risk || 0,
          cvar: metrics.expected_shortfall || 0,
          sharpeRatio: metrics.sharpe_ratio || 0,
          overallRiskScore: metrics.overall_risk_score || 0,
          liquidationPrice: metrics.liquidation_price || position.liquidationPrice || 0,
          currentPrice: metrics.current_price || position.currentPrice || 0,
          distanceToLiquidation: metrics.distance_to_liquidation_pct || 0,
          source: 'fastapi'
        });
      } else {
        // Fallback to basic calculations
        const healthFactor = position.healthFactor || 2.0;
        const currentPrice = position.currentPrice || 0;
        const liquidationPrice = position.liquidationPrice || 0;

        const distanceToLiquidation = currentPrice > 0 && liquidationPrice > 0
          ? ((currentPrice - liquidationPrice) / currentPrice) * 100
          : 0;

        // Estimate liquidation probability based on health factor
        let liquidationProb = 0;
        if (healthFactor < 1.1) liquidationProb = 0.8;
        else if (healthFactor < 1.2) liquidationProb = 0.5;
        else if (healthFactor < 1.5) liquidationProb = 0.2;
        else if (healthFactor < 2.0) liquidationProb = 0.05;

        setRiskMetrics({
          healthFactor,
          liquidationProb,
          liquidationPrice,
          currentPrice,
          distanceToLiquidation,
          var99: 0,
          cvar: 0,
          sharpeRatio: 0,
          overallRiskScore: (100 - healthFactor * 40).toFixed(1),
          source: 'basic'
        });
      }

      setLastUpdate(new Date());
    } catch (error) {
      console.error('[LiquidationAlert] Risk calculation failed:', error);

      // Fallback to position data
      setRiskMetrics({
        healthFactor: position.healthFactor || 2.0,
        liquidationProb: 0,
        liquidationPrice: position.liquidationPrice || 0,
        currentPrice: position.currentPrice || 0,
        distanceToLiquidation: 0,
        var99: 0,
        cvar: 0,
        sharpeRatio: 0,
        overallRiskScore: 0,
        source: 'fallback',
        error: error.message
      });
    } finally {
      setLoading(false);
    }
  }, [position, marketData, historicalPrices, checkFastAPIAvailability]);

  // Initial calculation
  useEffect(() => {
    calculateRisk();
  }, [calculateRisk]);

  // Auto-refresh
  useEffect(() => {
    if (!refreshInterval) return;

    const interval = setInterval(() => {
      calculateRisk();
    }, refreshInterval);

    return () => clearInterval(interval);
  }, [refreshInterval, calculateRisk]);

  if (!position || !riskMetrics) {
    return null;
  }

  const alertLevel = getAlertLevel(riskMetrics.healthFactor);
  const colors = getColorScheme(alertLevel);

  const handleAddCollateral = () => {
    if (onAddCollateral) {
      onAddCollateral(position);
    }
  };

  return (
    <div style={{
      background: `linear-gradient(135deg, ${colors.bg} 0%, rgba(15, 23, 42, 0.5) 100%)`,
      border: `2px solid ${colors.border}`,
      borderRadius: 16,
      padding: 24,
      marginBottom: 16
    }}>
      {/* Header Section */}
      <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', marginBottom: 24 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
          <div style={{
            width: 48,
            height: 48,
            borderRadius: 12,
            background: colors.bg,
            border: `2px solid ${colors.border}`,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: colors.primary
          }}>
            {alertLevel === 'critical' ? (
              <AlertTriangle size={24} style={{ animation: 'pulse 2s ease-in-out infinite' }} />
            ) : alertLevel === 'healthy' ? (
              <Shield size={24} />
            ) : (
              <TrendingDown size={24} />
            )}
          </div>

          <div>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 6 }}>
              <h3 style={{ fontSize: 20, fontWeight: 700, color: 'white', margin: 0 }}>
                {position.protocol || 'Position'} - {position.collateralAsset || 'ETH'}
              </h3>
              <span style={{
                fontSize: 11,
                padding: '4px 10px',
                background: colors.bg,
                border: `1px solid ${colors.border}`,
                borderRadius: 6,
                fontWeight: 700,
                color: colors.primary,
                letterSpacing: 0.5
              }}>
                {alertLevel.toUpperCase()}
              </span>
            </div>
            <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.8)', margin: 0 }}>
              Chain: {position.chain || 'Arbitrum'}
              {lastUpdate && (
                <span style={{ marginLeft: 12 }}>
                  • Updated {Math.round((Date.now() - lastUpdate) / 1000)}s ago
                </span>
              )}
            </p>
          </div>
        </div>

        {/* Add Collateral Button */}
        {alertLevel !== 'healthy' && (
          <button
            onClick={handleAddCollateral}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 8,
              padding: '10px 20px',
              background: 'rgb(59, 130, 246)',
              border: 'none',
              borderRadius: 8,
              color: 'white',
              fontSize: 14,
              fontWeight: 600,
              cursor: 'pointer',
              transition: 'all 0.3s'
            }}
            onMouseEnter={(e) => e.currentTarget.style.background = 'rgb(37, 99, 235)'}
            onMouseLeave={(e) => e.currentTarget.style.background = 'rgb(59, 130, 246)'}
          >
            <Plus size={16} />
            Add Collateral
          </button>
        )}
      </div>

      {/* Main Metrics */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: 16, marginBottom: 20 }}>
        {/* Health Factor */}
        <div style={{
          background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
          borderRadius: 12,
          padding: 20,
          border: '1px solid rgba(34, 211, 238, 0.2)'
        }}>
          <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 500 }}>
            Health Factor
          </div>
          <div style={{ fontSize: 32, fontWeight: 800, color: colors.primary }}>
            {riskMetrics.healthFactor.toFixed(2)}
          </div>
        </div>

        {/* Distance to Liquidation */}
        <div style={{
          background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
          borderRadius: 12,
          padding: 20,
          border: '1px solid rgba(34, 211, 238, 0.2)'
        }}>
          <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 500 }}>
            Distance to Liquidation
          </div>
          <div style={{ fontSize: 32, fontWeight: 800, color: colors.primary }}>
            {riskMetrics.distanceToLiquidation.toFixed(1)}%
          </div>
        </div>

        {/* Liquidation Probability */}
        <div style={{
          background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
          borderRadius: 12,
          padding: 20,
          border: '1px solid rgba(34, 211, 238, 0.2)'
        }}>
          <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 500 }}>
            Liquidation Probability (24h)
          </div>
          <div style={{ fontSize: 32, fontWeight: 800, color: colors.primary }}>
            {(riskMetrics.liquidationProb * 100).toFixed(1)}%
          </div>
        </div>

        {/* Risk Score (FastAPI only) */}
        {riskMetrics.source === 'fastapi' && (
          <div style={{
            background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
            borderRadius: 12,
            padding: 20,
            border: '1px solid rgba(34, 211, 238, 0.2)'
          }}>
            <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 500, display: 'flex', alignItems: 'center', gap: 6 }}>
              Risk Score
              <Zap size={14} style={{ color: 'rgb(245, 158, 11)' }} title="FastAPI Powered" />
            </div>
            <div style={{ fontSize: 32, fontWeight: 800, color: colors.primary }}>
              {riskMetrics.overallRiskScore.toFixed(1)}
            </div>
          </div>
        )}
      </div>

      {/* Expandable Details */}
      <div>
        <button
          onClick={() => setExpanded(!expanded)}
          style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            width: '100%',
            padding: '12px 16px',
            background: 'rgba(15, 23, 42, 0.5)',
            border: '1px solid rgba(34, 211, 238, 0.2)',
            borderRadius: 8,
            fontSize: 14,
            fontWeight: 600,
            color: 'white',
            cursor: 'pointer',
            transition: 'all 0.3s'
          }}
          onMouseEnter={(e) => {
            e.currentTarget.style.background = 'rgba(15, 23, 42, 0.8)';
            e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.4)';
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.background = 'rgba(15, 23, 42, 0.5)';
            e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.2)';
          }}
        >
          <span>Advanced Metrics</span>
          {expanded ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
        </button>

        {expanded && (
          <div style={{ marginTop: 16, paddingTop: 16, borderTop: '1px solid rgba(34, 211, 238, 0.2)' }}>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: 16 }}>
              {/* VaR (99%) */}
              {riskMetrics.source === 'fastapi' && riskMetrics.var99 > 0 && (
                <div style={{
                  background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
                  borderRadius: 12,
                  padding: 16,
                  border: '1px solid rgba(34, 211, 238, 0.2)'
                }}>
                  <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 8 }}>
                    <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 500 }}>VaR (99%)</span>
                    <Activity size={14} style={{ color: 'rgba(203, 213, 225, 0.5)' }} />
                  </div>
                  <div style={{ fontSize: 20, fontWeight: 700, color: 'white', marginBottom: 6 }}>
                    ${riskMetrics.var99.toLocaleString(undefined, { maximumFractionDigits: 2 })}
                  </div>
                  <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                    Maximum expected loss (1% probability)
                  </div>
                </div>
              )}

              {/* CVaR / Expected Shortfall */}
              {riskMetrics.source === 'fastapi' && riskMetrics.cvar > 0 && (
                <div style={{
                  background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
                  borderRadius: 12,
                  padding: 16,
                  border: '1px solid rgba(34, 211, 238, 0.2)'
                }}>
                  <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 8 }}>
                    <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontWeight: 500 }}>CVaR / Expected Shortfall</span>
                    <Activity size={14} style={{ color: 'rgba(203, 213, 225, 0.5)' }} />
                  </div>
                  <div style={{ fontSize: 20, fontWeight: 700, color: 'white', marginBottom: 6 }}>
                    ${riskMetrics.cvar.toLocaleString(undefined, { maximumFractionDigits: 2 })}
                  </div>
                  <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                    Average loss when VaR is exceeded
                  </div>
                </div>
              )}

              {/* Liquidation Price */}
              {riskMetrics.liquidationPrice > 0 && (
                <div style={{
                  background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
                  borderRadius: 12,
                  padding: 16,
                  border: '1px solid rgba(34, 211, 238, 0.2)'
                }}>
                  <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 500 }}>
                    Liquidation Price
                  </div>
                  <div style={{ fontSize: 20, fontWeight: 700, color: 'white', marginBottom: 6 }}>
                    ${riskMetrics.liquidationPrice.toLocaleString(undefined, { maximumFractionDigits: 2 })}
                  </div>
                  <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                    Current: ${riskMetrics.currentPrice.toLocaleString(undefined, { maximumFractionDigits: 2 })}
                  </div>
                </div>
              )}

              {/* Sharpe Ratio */}
              {riskMetrics.source === 'fastapi' && riskMetrics.sharpeRatio !== 0 && (
                <div style={{
                  background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
                  borderRadius: 12,
                  padding: 16,
                  border: '1px solid rgba(34, 211, 238, 0.2)'
                }}>
                  <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8, fontWeight: 500 }}>
                    Sharpe Ratio
                  </div>
                  <div style={{ fontSize: 20, fontWeight: 700, color: 'white', marginBottom: 6 }}>
                    {riskMetrics.sharpeRatio.toFixed(2)}
                  </div>
                  <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)' }}>
                    Risk-adjusted return metric
                  </div>
                </div>
              )}
            </div>

            {/* Data Source Badge */}
            <div style={{ marginTop: 16, display: 'flex', alignItems: 'center', justifyContent: 'space-between', fontSize: 12 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, color: 'rgba(203, 213, 225, 0.7)' }}>
                <span>Data Source:</span>
                {riskMetrics.source === 'fastapi' ? (
                  <span style={{ display: 'flex', alignItems: 'center', color: 'rgb(34, 197, 94)', fontWeight: 600 }}>
                    <Zap size={14} style={{ marginRight: 4 }} />
                    FastAPI (Advanced)
                  </span>
                ) : (
                  <span style={{ color: 'rgba(203, 213, 225, 0.8)', fontWeight: 600 }}>Basic Calculations</span>
                )}
              </div>
              {loading && (
                <span style={{ display: 'flex', alignItems: 'center', color: 'rgb(59, 130, 246)' }}>
                  <Activity size={14} style={{ marginRight: 4, animation: 'pulse 2s ease-in-out infinite' }} />
                  Calculating...
                </span>
              )}
            </div>

            {/* Error Display */}
            {riskMetrics.error && (
              <div style={{
                marginTop: 12,
                padding: 12,
                background: 'rgba(249, 115, 22, 0.1)',
                border: '1px solid rgba(249, 115, 22, 0.3)',
                borderRadius: 8,
                fontSize: 12,
                color: 'rgb(249, 115, 22)'
              }}>
                Warning: {riskMetrics.error}
              </div>
            )}
          </div>
        )}
      </div>

      {/* Warning Messages */}
      {alertLevel === 'critical' && (
        <div style={{
          marginTop: 20,
          padding: 16,
          background: 'rgba(239, 68, 68, 0.15)',
          border: '1px solid rgba(239, 68, 68, 0.4)',
          borderRadius: 12
        }}>
          <div style={{ display: 'flex', alignItems: 'flex-start', gap: 12 }}>
            <AlertTriangle size={20} style={{ color: 'rgb(239, 68, 68)', flexShrink: 0, marginTop: 2 }} />
            <div style={{ fontSize: 14, color: 'rgba(255, 255, 255, 0.95)' }}>
              <strong>Critical Risk!</strong> Your position is at high risk of liquidation.
              Add collateral immediately to improve your health factor.
            </div>
          </div>
        </div>
      )}

      {alertLevel === 'danger' && (
        <div style={{
          marginTop: 20,
          padding: 16,
          background: 'rgba(249, 115, 22, 0.15)',
          border: '1px solid rgba(249, 115, 22, 0.4)',
          borderRadius: 12
        }}>
          <div style={{ display: 'flex', alignItems: 'flex-start', gap: 12 }}>
            <AlertTriangle size={20} style={{ color: 'rgb(249, 115, 22)', flexShrink: 0, marginTop: 2 }} />
            <div style={{ fontSize: 14, color: 'rgba(255, 255, 255, 0.95)' }}>
              <strong>High Risk:</strong> Your health factor is below the recommended threshold.
              Consider adding more collateral to reduce liquidation risk.
            </div>
          </div>
        </div>
      )}

      {alertLevel === 'warning' && (
        <div style={{
          marginTop: 20,
          padding: 16,
          background: 'rgba(245, 158, 11, 0.15)',
          border: '1px solid rgba(245, 158, 11, 0.4)',
          borderRadius: 12
        }}>
          <div style={{ fontSize: 14, color: 'rgba(255, 255, 255, 0.95)' }}>
            <strong>Moderate Risk:</strong> Monitor your position closely.
            Market volatility could impact your health factor.
          </div>
        </div>
      )}
    </div>
  );
};

export default LiquidationAlert;
