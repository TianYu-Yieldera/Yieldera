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
        bg: 'bg-green-50',
        border: 'border-green-200',
        text: 'text-green-800',
        icon: 'text-green-600',
        badge: 'bg-green-100 text-green-800'
      },
      warning: {
        bg: 'bg-yellow-50',
        border: 'border-yellow-200',
        text: 'text-yellow-800',
        icon: 'text-yellow-600',
        badge: 'bg-yellow-100 text-yellow-800'
      },
      danger: {
        bg: 'bg-orange-50',
        border: 'border-orange-200',
        text: 'text-orange-800',
        icon: 'text-orange-600',
        badge: 'bg-orange-100 text-orange-800'
      },
      critical: {
        bg: 'bg-red-50',
        border: 'border-red-200',
        text: 'text-red-800',
        icon: 'text-red-600',
        badge: 'bg-red-100 text-red-800'
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
    <div className={`rounded-lg border-2 ${colors.border} ${colors.bg} p-4 mb-4`}>
      {/* Header Section */}
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center space-x-3">
          <div className={`${colors.icon}`}>
            {alertLevel === 'critical' ? (
              <AlertTriangle className="h-6 w-6 animate-pulse" />
            ) : alertLevel === 'healthy' ? (
              <Shield className="h-6 w-6" />
            ) : (
              <TrendingDown className="h-6 w-6" />
            )}
          </div>

          <div>
            <div className="flex items-center space-x-2">
              <h3 className={`font-bold text-lg ${colors.text}`}>
                {position.protocol || 'Position'} - {position.collateralAsset || 'ETH'}
              </h3>
              <span className={`text-xs px-2 py-0.5 rounded-full ${colors.badge}`}>
                {alertLevel.toUpperCase()}
              </span>
            </div>
            <p className={`text-sm ${colors.text} opacity-75`}>
              Chain: {position.chain || 'Arbitrum'}
              {lastUpdate && (
                <span className="ml-2">
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
            className="flex items-center space-x-1 px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors"
          >
            <Plus className="h-4 w-4" />
            <span>Add Collateral</span>
          </button>
        )}
      </div>

      {/* Main Metrics */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-3">
        {/* Health Factor */}
        <div className="bg-white rounded-lg p-3 border border-gray-200">
          <div className="text-xs text-gray-600 mb-1">Health Factor</div>
          <div className={`text-2xl font-bold ${colors.text}`}>
            {riskMetrics.healthFactor.toFixed(2)}
          </div>
        </div>

        {/* Distance to Liquidation */}
        <div className="bg-white rounded-lg p-3 border border-gray-200">
          <div className="text-xs text-gray-600 mb-1">To Liquidation</div>
          <div className={`text-2xl font-bold ${colors.text}`}>
            {riskMetrics.distanceToLiquidation.toFixed(1)}%
          </div>
        </div>

        {/* Liquidation Probability */}
        <div className="bg-white rounded-lg p-3 border border-gray-200">
          <div className="text-xs text-gray-600 mb-1">Liq. Probability (24h)</div>
          <div className={`text-2xl font-bold ${colors.text}`}>
            {(riskMetrics.liquidationProb * 100).toFixed(1)}%
          </div>
        </div>

        {/* Risk Score */}
        {riskMetrics.source === 'fastapi' && (
          <div className="bg-white rounded-lg p-3 border border-gray-200">
            <div className="text-xs text-gray-600 mb-1 flex items-center">
              Risk Score
              <Zap className="h-3 w-3 ml-1 text-yellow-500" title="FastAPI Powered" />
            </div>
            <div className={`text-2xl font-bold ${colors.text}`}>
              {riskMetrics.overallRiskScore.toFixed(1)}
            </div>
          </div>
        )}
      </div>

      {/* Expandable Details */}
      <div>
        <button
          onClick={() => setExpanded(!expanded)}
          className={`flex items-center justify-between w-full text-sm font-medium ${colors.text} hover:opacity-75 transition-opacity`}
        >
          <span>Advanced Metrics</span>
          {expanded ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
        </button>

        {expanded && (
          <div className="mt-3 pt-3 border-t border-gray-200">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {/* VaR (99%) */}
              {riskMetrics.source === 'fastapi' && riskMetrics.var99 > 0 && (
                <div className="bg-white rounded-lg p-3 border border-gray-200">
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-xs text-gray-600">VaR (99%)</span>
                    <Activity className="h-3 w-3 text-gray-400" />
                  </div>
                  <div className="text-lg font-bold text-gray-800">
                    ${riskMetrics.var99.toLocaleString(undefined, { maximumFractionDigits: 2 })}
                  </div>
                  <div className="text-xs text-gray-500 mt-1">
                    Maximum expected loss (1% probability)
                  </div>
                </div>
              )}

              {/* CVaR / Expected Shortfall */}
              {riskMetrics.source === 'fastapi' && riskMetrics.cvar > 0 && (
                <div className="bg-white rounded-lg p-3 border border-gray-200">
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-xs text-gray-600">CVaR / Expected Shortfall</span>
                    <Activity className="h-3 w-3 text-gray-400" />
                  </div>
                  <div className="text-lg font-bold text-gray-800">
                    ${riskMetrics.cvar.toLocaleString(undefined, { maximumFractionDigits: 2 })}
                  </div>
                  <div className="text-xs text-gray-500 mt-1">
                    Average loss when VaR is exceeded
                  </div>
                </div>
              )}

              {/* Liquidation Price */}
              {riskMetrics.liquidationPrice > 0 && (
                <div className="bg-white rounded-lg p-3 border border-gray-200">
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-xs text-gray-600">Liquidation Price</span>
                  </div>
                  <div className="text-lg font-bold text-gray-800">
                    ${riskMetrics.liquidationPrice.toLocaleString(undefined, { maximumFractionDigits: 2 })}
                  </div>
                  <div className="text-xs text-gray-500 mt-1">
                    Current: ${riskMetrics.currentPrice.toLocaleString(undefined, { maximumFractionDigits: 2 })}
                  </div>
                </div>
              )}

              {/* Sharpe Ratio */}
              {riskMetrics.source === 'fastapi' && riskMetrics.sharpeRatio !== 0 && (
                <div className="bg-white rounded-lg p-3 border border-gray-200">
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-xs text-gray-600">Sharpe Ratio</span>
                  </div>
                  <div className="text-lg font-bold text-gray-800">
                    {riskMetrics.sharpeRatio.toFixed(2)}
                  </div>
                  <div className="text-xs text-gray-500 mt-1">
                    Risk-adjusted return metric
                  </div>
                </div>
              )}
            </div>

            {/* Data Source Badge */}
            <div className="mt-3 flex items-center justify-between text-xs text-gray-500">
              <div className="flex items-center space-x-2">
                <span>Data Source:</span>
                {riskMetrics.source === 'fastapi' ? (
                  <span className="flex items-center text-green-600 font-medium">
                    <Zap className="h-3 w-3 mr-1" />
                    FastAPI (Advanced)
                  </span>
                ) : (
                  <span className="text-gray-600">Basic Calculations</span>
                )}
              </div>
              {loading && (
                <span className="flex items-center text-blue-600">
                  <Activity className="h-3 w-3 mr-1 animate-pulse" />
                  Calculating...
                </span>
              )}
            </div>

            {/* Error Display */}
            {riskMetrics.error && (
              <div className="mt-2 text-xs text-orange-600 bg-orange-50 rounded p-2">
                Warning: {riskMetrics.error}
              </div>
            )}
          </div>
        )}
      </div>

      {/* Warning Messages */}
      {alertLevel === 'critical' && (
        <div className="mt-3 bg-red-100 border border-red-300 rounded-lg p-3">
          <div className="flex items-start space-x-2">
            <AlertTriangle className="h-5 w-5 text-red-600 flex-shrink-0 mt-0.5" />
            <div className="text-sm text-red-800">
              <strong>Critical Risk!</strong> Your position is at high risk of liquidation.
              Add collateral immediately to improve your health factor.
            </div>
          </div>
        </div>
      )}

      {alertLevel === 'danger' && (
        <div className="mt-3 bg-orange-100 border border-orange-300 rounded-lg p-3">
          <div className="flex items-start space-x-2">
            <AlertTriangle className="h-5 w-5 text-orange-600 flex-shrink-0 mt-0.5" />
            <div className="text-sm text-orange-800">
              <strong>High Risk:</strong> Your health factor is below the recommended threshold.
              Consider adding more collateral to reduce liquidation risk.
            </div>
          </div>
        </div>
      )}

      {alertLevel === 'warning' && (
        <div className="mt-3 bg-yellow-100 border border-yellow-300 rounded-lg p-3">
          <div className="text-sm text-yellow-800">
            <strong>Moderate Risk:</strong> Monitor your position closely.
            Market volatility could impact your health factor.
          </div>
        </div>
      )}
    </div>
  );
};

export default LiquidationAlert;
