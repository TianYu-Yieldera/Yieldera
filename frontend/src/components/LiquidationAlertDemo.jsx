/**
 * Liquidation Alert Demo Page
 * Demonstrates the LiquidationAlert component with various scenarios
 */

import React, { useState } from 'react';
import { RefreshCw, ArrowLeft, Info } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import LiquidationAlert from './LiquidationAlert';

const LiquidationAlertDemo = () => {
  const navigate = useNavigate();
  const [selectedScenario, setSelectedScenario] = useState('healthy');
  const [refreshKey, setRefreshKey] = useState(0);

  // Demo scenarios with different risk levels
  const scenarios = {
    healthy: {
      name: 'Healthy Position',
      description: 'Health factor > 2.0, low liquidation risk',
      color: 'rgb(34, 197, 94)',
      position: {
        id: 'demo-healthy-1',
        protocol: 'Aave V3',
        chain: 'Arbitrum',
        collateralAsset: 'ETH',
        collateralValueUSD: 50000,
        borrowValueUSD: 20000,
        healthFactor: 2.35,
        leverage: 1.4,
        ltv: 0.4,
        currentPrice: 2500,
        liquidationPrice: 1800,
        position_age_days: 45
      },
      marketData: {
        asset: 'ETH',
        price: 2500,
        priceChange24h: 2.5,
        volume24h: 15000000000,
        volatility: 0.15,
        liquidity: 5000000000,
        timestamp: new Date().toISOString()
      },
      historicalPrices: generateHistoricalPrices(2500, 0.15, 100)
    },
    warning: {
      name: 'Warning Position',
      description: 'Health factor 1.5-2.0, moderate risk',
      color: 'rgb(245, 158, 11)',
      position: {
        id: 'demo-warning-1',
        protocol: 'Compound V3',
        chain: 'Base',
        collateralAsset: 'WBTC',
        collateralValueUSD: 100000,
        borrowValueUSD: 60000,
        healthFactor: 1.67,
        leverage: 1.6,
        ltv: 0.6,
        currentPrice: 45000,
        liquidationPrice: 40000,
        position_age_days: 30
      },
      marketData: {
        asset: 'WBTC',
        price: 45000,
        priceChange24h: -3.2,
        volume24h: 8000000000,
        volatility: 0.18,
        liquidity: 2000000000,
        timestamp: new Date().toISOString()
      },
      historicalPrices: generateHistoricalPrices(45000, 0.18, 100)
    },
    danger: {
      name: 'Danger Position',
      description: 'Health factor 1.2-1.5, high risk',
      color: 'rgb(249, 115, 22)',
      position: {
        id: 'demo-danger-1',
        protocol: 'Aave V3',
        chain: 'Arbitrum',
        collateralAsset: 'ETH',
        collateralValueUSD: 75000,
        borrowValueUSD: 55000,
        healthFactor: 1.36,
        leverage: 1.73,
        ltv: 0.73,
        currentPrice: 2400,
        liquidationPrice: 2200,
        position_age_days: 60
      },
      marketData: {
        asset: 'ETH',
        price: 2400,
        priceChange24h: -5.8,
        volume24h: 18000000000,
        volatility: 0.25,
        liquidity: 4500000000,
        timestamp: new Date().toISOString()
      },
      historicalPrices: generateHistoricalPrices(2400, 0.25, 100)
    },
    critical: {
      name: 'Critical Position',
      description: 'Health factor < 1.2, imminent liquidation risk',
      color: 'rgb(239, 68, 68)',
      position: {
        id: 'demo-critical-1',
        protocol: 'GMX V2',
        chain: 'Arbitrum',
        collateralAsset: 'USDC',
        collateralValueUSD: 50000,
        borrowValueUSD: 42000,
        healthFactor: 1.15,
        leverage: 1.84,
        ltv: 0.84,
        currentPrice: 2350,
        liquidationPrice: 2300,
        position_age_days: 15
      },
      marketData: {
        asset: 'ETH',
        price: 2350,
        priceChange24h: -8.5,
        volume24h: 20000000000,
        volatility: 0.35,
        liquidity: 4000000000,
        timestamp: new Date().toISOString()
      },
      historicalPrices: generateHistoricalPrices(2350, 0.35, 100)
    }
  };

  const handleAddCollateral = (position) => {
    alert(`Add Collateral clicked for position: ${position.id}\n\nIn production, this would:\n1. Open a modal with collateral deposit form\n2. Calculate required collateral to improve health factor\n3. Execute deposit transaction\n4. Refresh position data`);
  };

  const handleRefresh = () => {
    setRefreshKey(prev => prev + 1);
  };

  const currentScenario = scenarios[selectedScenario];

  return (
    <div style={{ minHeight: '100vh', background: 'rgb(15, 23, 42)' }}>
      <style>{`
        @keyframes star-twinkle {
          0%, 100% { opacity: 0.3; transform: scale(0.8); }
          50% { opacity: 1; transform: scale(1); }
        }

        @keyframes pulse {
          0%, 100% { opacity: 1; }
          50% { opacity: 0.5; }
        }
      `}</style>

      {/* Top Navigation Bar */}
      <div style={{
        background: 'rgba(15, 23, 42, 0.95)',
        borderBottom: '1px solid rgba(34, 211, 238, 0.2)',
        padding: '16px 24px',
        position: 'sticky',
        top: 0,
        zIndex: 100,
        backdropFilter: 'blur(10px)'
      }}>
        <div style={{ maxWidth: 1400, margin: '0 auto', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <button
            onClick={() => navigate('/tutorial')}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 8,
              padding: '8px 16px',
              background: 'rgba(34, 211, 238, 0.1)',
              border: '1px solid rgba(34, 211, 238, 0.3)',
              borderRadius: 8,
              color: 'rgb(34, 211, 238)',
              fontSize: 14,
              fontWeight: 600,
              cursor: 'pointer',
              transition: 'all 0.3s'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = 'rgba(34, 211, 238, 0.2)';
              e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.5)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = 'rgba(34, 211, 238, 0.1)';
              e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.3)';
            }}
          >
            <ArrowLeft size={16} />
            Back to Tutorial
          </button>

          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <h1 style={{ fontSize: 20, fontWeight: 700, color: 'white', margin: 0 }}>
              AI Risk Dashboard Demo
            </h1>
          </div>

          <button
            onClick={handleRefresh}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 8,
              padding: '8px 16px',
              background: 'rgba(168, 85, 247, 0.1)',
              border: '1px solid rgba(168, 85, 247, 0.3)',
              borderRadius: 8,
              color: 'rgb(168, 85, 247)',
              fontSize: 14,
              fontWeight: 600,
              cursor: 'pointer',
              transition: 'all 0.3s'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = 'rgba(168, 85, 247, 0.2)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = 'rgba(168, 85, 247, 0.1)';
            }}
          >
            <RefreshCw size={16} />
            Refresh
          </button>
        </div>
      </div>

      <div style={{ maxWidth: 1400, margin: '0 auto', padding: '32px 24px' }}>
        {/* Info Banner */}
        <div style={{
          marginBottom: 32,
          padding: '16px 20px',
          background: 'linear-gradient(135deg, rgba(168, 85, 247, 0.1) 0%, rgba(124, 58, 237, 0.1) 100%)',
          border: '1px solid rgba(168, 85, 247, 0.3)',
          borderRadius: 12,
          display: 'flex',
          alignItems: 'center',
          gap: 12
        }}>
          <Info size={20} style={{ color: 'rgb(168, 85, 247)', flexShrink: 0 }} />
          <div style={{ flex: 1 }}>
            <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.9)', margin: 0, lineHeight: 1.6 }}>
              <strong style={{ color: 'white' }}>Interactive Demo:</strong> Select a risk scenario below to see how the AI Risk Dashboard monitors different liquidation risk levels in real-time.
            </p>
          </div>
        </div>

        {/* Scenario Selector */}
        <div style={{ marginBottom: 32 }}>
          <h2 style={{ fontSize: 18, fontWeight: 700, color: 'white', marginBottom: 16, display: 'flex', alignItems: 'center', gap: 8 }}>
            Risk Scenarios
            <span style={{
              fontSize: 12,
              padding: '4px 10px',
              background: 'rgba(34, 211, 238, 0.2)',
              borderRadius: 6,
              fontWeight: 600,
              color: 'rgb(34, 211, 238)'
            }}>
              Select One
            </span>
          </h2>

          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: 16 }}>
            {Object.entries(scenarios).map(([key, scenario]) => (
              <button
                key={key}
                onClick={() => setSelectedScenario(key)}
                style={{
                  padding: 20,
                  borderRadius: 12,
                  border: selectedScenario === key ? `2px solid ${scenario.color}` : '1px solid rgba(34, 211, 238, 0.2)',
                  background: selectedScenario === key
                    ? `linear-gradient(135deg, ${scenario.color}20 0%, ${scenario.color}10 100%)`
                    : 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
                  textAlign: 'left',
                  transition: 'all 0.3s',
                  cursor: 'pointer',
                  boxShadow: selectedScenario === key ? `0 8px 24px ${scenario.color}40` : 'none',
                  position: 'relative',
                  overflow: 'hidden'
                }}
                onMouseEnter={(e) => {
                  if (selectedScenario !== key) {
                    e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.5)';
                    e.currentTarget.style.transform = 'translateY(-4px)';
                    e.currentTarget.style.boxShadow = '0 4px 12px rgba(34, 211, 238, 0.2)';
                  }
                }}
                onMouseLeave={(e) => {
                  if (selectedScenario !== key) {
                    e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.2)';
                    e.currentTarget.style.transform = 'translateY(0)';
                    e.currentTarget.style.boxShadow = 'none';
                  }
                }}
              >
                {/* Selected Indicator */}
                {selectedScenario === key && (
                  <div style={{
                    position: 'absolute',
                    top: 12,
                    right: 12,
                    width: 12,
                    height: 12,
                    borderRadius: '50%',
                    background: scenario.color,
                    boxShadow: `0 0 12px ${scenario.color}`,
                    animation: 'pulse 2s ease-in-out infinite'
                  }} />
                )}

                <div style={{ fontWeight: 700, color: 'white', marginBottom: 6, fontSize: 16 }}>
                  {scenario.name}
                </div>
                <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 12, lineHeight: 1.5 }}>
                  {scenario.description}
                </div>
                <div style={{
                  display: 'inline-block',
                  padding: '6px 14px',
                  background: selectedScenario === key ? scenario.color : 'rgba(34, 211, 238, 0.2)',
                  borderRadius: 8,
                  fontSize: 18,
                  fontWeight: 800,
                  color: 'white'
                }}>
                  HF: {scenario.position.healthFactor.toFixed(2)}
                </div>
              </button>
            ))}
          </div>
        </div>

        {/* Current Scenario Info Panel */}
        <div style={{
          marginBottom: 32,
          background: `linear-gradient(135deg, ${currentScenario.color}15 0%, ${currentScenario.color}08 100%)`,
          borderRadius: 12,
          border: `1px solid ${currentScenario.color}40`,
          padding: 24,
          position: 'relative',
          overflow: 'hidden'
        }}>
          {/* Decorative gradient */}
          <div style={{
            position: 'absolute',
            top: -50,
            right: -50,
            width: 200,
            height: 200,
            background: `radial-gradient(circle, ${currentScenario.color}30 0%, transparent 70%)`,
            filter: 'blur(40px)'
          }} />

          <div style={{ position: 'relative', zIndex: 1 }}>
            <h3 style={{ fontWeight: 700, color: 'white', marginBottom: 20, fontSize: 18, display: 'flex', alignItems: 'center', gap: 12 }}>
              <span style={{
                width: 8,
                height: 8,
                borderRadius: '50%',
                background: currentScenario.color,
                boxShadow: `0 0 12px ${currentScenario.color}`
              }} />
              Current Scenario: {currentScenario.name}
            </h3>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))', gap: 20 }}>
              {[
                { label: 'Protocol', value: currentScenario.position.protocol },
                { label: 'Chain', value: currentScenario.position.chain },
                { label: 'Collateral', value: `$${currentScenario.position.collateralValueUSD.toLocaleString()}` },
                { label: 'Debt', value: `$${currentScenario.position.borrowValueUSD.toLocaleString()}` },
                { label: 'LTV Ratio', value: `${(currentScenario.position.ltv * 100).toFixed(1)}%` },
                { label: 'Health Factor', value: currentScenario.position.healthFactor.toFixed(2) }
              ].map((item, idx) => (
                <div key={idx}>
                  <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 6, fontWeight: 500 }}>
                    {item.label}
                  </div>
                  <div style={{ fontSize: 16, fontWeight: 700, color: 'white' }}>
                    {item.value}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Risk Assessment Panel */}
        <div>
          <h2 style={{ fontSize: 18, fontWeight: 700, color: 'white', marginBottom: 16 }}>
            Live Risk Assessment
          </h2>
          <LiquidationAlert
            key={`${selectedScenario}-${refreshKey}`}
            position={currentScenario.position}
            marketData={currentScenario.marketData}
            historicalPrices={currentScenario.historicalPrices}
            onAddCollateral={handleAddCollateral}
            refreshInterval={0}
          />
        </div>
      </div>
    </div>
  );
};

// Helper function to generate mock historical prices
function generateHistoricalPrices(basePrice, volatility, count) {
  const prices = [];
  let price = basePrice;

  for (let i = count; i > 0; i--) {
    const change = (Math.random() - 0.5) * 2 * volatility * price;
    price = Math.max(price + change, basePrice * 0.5);

    prices.push({
      timestamp: new Date(Date.now() - i * 3600000).toISOString(),
      price: parseFloat(price.toFixed(2))
    });
  }

  return prices;
}

export default LiquidationAlertDemo;
