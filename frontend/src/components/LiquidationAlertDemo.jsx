/**
 * Liquidation Alert Demo Page
 * Demonstrates the LiquidationAlert component with various scenarios
 */

import React, { useState } from 'react';
import { RefreshCw, Zap, Play } from 'lucide-react';
import LiquidationAlert from './LiquidationAlert';

const LiquidationAlertDemo = () => {
  const [selectedScenario, setSelectedScenario] = useState('healthy');
  const [refreshKey, setRefreshKey] = useState(0);

  // Demo scenarios with different risk levels
  const scenarios = {
    healthy: {
      name: 'Healthy Position',
      description: 'Health factor > 2.0, low liquidation risk',
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
    <div style={{ minHeight: '100vh', background: '#f9fafb', padding: 24 }}>
      <div style={{ maxWidth: 1200, margin: '0 auto' }}>
        {/* Header */}
        <div style={{ marginBottom: 32 }}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
            <div>
              <h1 style={{ fontSize: 32, fontWeight: 700, color: '#111827', margin: '0 0 8px 0' }}>
                Liquidation Alert Demo
              </h1>
              <p style={{ fontSize: 16, color: '#6b7280', margin: 0 }}>
                Interactive demonstration of the LiquidationAlert component with FastAPI integration
              </p>
            </div>
            <button
              onClick={handleRefresh}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: 8,
                padding: '10px 18px',
                background: '#3b82f6',
                border: 'none',
                borderRadius: 8,
                color: 'white',
                fontSize: 14,
                fontWeight: 600,
                cursor: 'pointer',
                transition: 'all 0.2s'
              }}
              onMouseEnter={(e) => e.currentTarget.style.background = '#2563eb'}
              onMouseLeave={(e) => e.currentTarget.style.background = '#3b82f6'}
            >
              <RefreshCw style={{ width: 16, height: 16 }} />
              <span>Refresh</span>
            </button>
          </div>

          {/* FastAPI Status */}
          <div style={{
            background: '#eff6ff',
            border: '1px solid #bfdbfe',
            borderRadius: 8,
            padding: 16
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8, color: '#1e40af' }}>
              <Zap style={{ width: 20, height: 20 }} />
              <div>
                <strong>FastAPI Integration:</strong> The component will automatically attempt to use
                the FastAPI service for advanced risk calculations (VaR, CVaR, Sharpe ratio).
                If the service is unavailable, it will fall back to basic calculations.
              </div>
            </div>
            <div style={{ marginTop: 8, fontSize: 13, color: '#1e40af' }}>
              <strong>To enable FastAPI:</strong>
              <code style={{ marginLeft: 8, background: '#dbeafe', padding: '2px 8px', borderRadius: 4, fontSize: 12 }}>
                cd services/ai && ./start-risk-api.sh
              </code>
            </div>
          </div>
        </div>

        {/* Scenario Selector */}
        <div style={{ marginBottom: 24 }}>
          <label style={{ display: 'block', fontSize: 14, fontWeight: 500, color: '#374151', marginBottom: 12 }}>
            Select Risk Scenario:
          </label>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 12 }}>
            {Object.entries(scenarios).map(([key, scenario]) => (
              <button
                key={key}
                onClick={() => setSelectedScenario(key)}
                style={{
                  padding: 16,
                  borderRadius: 8,
                  border: selectedScenario === key ? '2px solid rgb(59, 130, 246)' : '2px solid rgb(229, 231, 235)',
                  background: selectedScenario === key ? 'rgb(239, 246, 255)' : 'white',
                  textAlign: 'left',
                  transition: 'all 0.2s',
                  cursor: 'pointer',
                  boxShadow: selectedScenario === key ? '0 4px 12px rgba(59, 130, 246, 0.3)' : 'none'
                }}
                onMouseEnter={(e) => {
                  if (selectedScenario !== key) {
                    e.currentTarget.style.borderColor = 'rgb(191, 219, 254)';
                    e.currentTarget.style.boxShadow = '0 2px 8px rgba(0, 0, 0, 0.1)';
                  }
                }}
                onMouseLeave={(e) => {
                  if (selectedScenario !== key) {
                    e.currentTarget.style.borderColor = 'rgb(229, 231, 235)';
                    e.currentTarget.style.boxShadow = 'none';
                  }
                }}
              >
                <div style={{ fontWeight: 700, color: 'rgb(17, 24, 39)', marginBottom: 4, fontSize: 14 }}>
                  {scenario.name}
                </div>
                <div style={{ fontSize: 11, color: 'rgb(107, 114, 128)', marginBottom: 8, lineHeight: 1.4 }}>
                  {scenario.description}
                </div>
                <div style={{ fontSize: 16, fontWeight: 700, color: selectedScenario === key ? 'rgb(59, 130, 246)' : 'rgb(17, 24, 39)' }}>
                  HF: {scenario.position.healthFactor.toFixed(2)}
                </div>
              </button>
            ))}
          </div>
        </div>

        {/* Current Scenario Details */}
        <div style={{ marginBottom: 24, background: 'white', borderRadius: 8, border: '1px solid rgb(229, 231, 235)', padding: 16 }}>
          <h3 style={{ fontWeight: 700, color: 'rgb(17, 24, 39)', marginBottom: 12, fontSize: 16 }}>Current Scenario: {currentScenario.name}</h3>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(120px, 1fr))', gap: 16, fontSize: 14 }}>
            <div>
              <div style={{ color: 'rgb(75, 85, 99)', marginBottom: 4 }}>Protocol</div>
              <div style={{ fontWeight: 600 }}>{currentScenario.position.protocol}</div>
            </div>
            <div>
              <div style={{ color: 'rgb(75, 85, 99)', marginBottom: 4 }}>Collateral</div>
              <div style={{ fontWeight: 600 }}>
                ${currentScenario.position.collateralValueUSD.toLocaleString()}
              </div>
            </div>
            <div>
              <div style={{ color: 'rgb(75, 85, 99)', marginBottom: 4 }}>Debt</div>
              <div style={{ fontWeight: 600 }}>
                ${currentScenario.position.borrowValueUSD.toLocaleString()}
              </div>
            </div>
            <div>
              <div style={{ color: 'rgb(75, 85, 99)', marginBottom: 4 }}>LTV</div>
              <div style={{ fontWeight: 600 }}>
                {(currentScenario.position.ltv * 100).toFixed(1)}%
              </div>
            </div>
          </div>
        </div>

        {/* Liquidation Alert Component */}
        <div style={{ marginBottom: 24 }}>
          <h2 style={{ fontSize: 20, fontWeight: 700, color: 'rgb(17, 24, 39)', marginBottom: 16 }}>Component Output:</h2>
          <LiquidationAlert
            key={`${selectedScenario}-${refreshKey}`}
            position={currentScenario.position}
            marketData={currentScenario.marketData}
            historicalPrices={currentScenario.historicalPrices}
            onAddCollateral={handleAddCollateral}
            refreshInterval={0} // Disable auto-refresh in demo
          />
        </div>

        {/* Technical Details */}
        <div style={{ background: 'rgb(31, 41, 55)', borderRadius: 8, padding: 24, color: 'white' }}>
          <h3 style={{ fontSize: 18, fontWeight: 700, marginBottom: 16, display: 'flex', alignItems: 'center' }}>
            <Play style={{ width: 20, height: 20, marginRight: 8 }} />
            Technical Features Demonstrated
          </h3>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: 16, fontSize: 14 }}>
            <div>
              <h4 style={{ fontWeight: 600, color: 'rgb(74, 222, 128)', marginBottom: 8 }}>✓ Implemented Features:</h4>
              <ul style={{ listStyle: 'none', padding: 0, margin: 0, color: 'rgb(209, 213, 219)' }}>
                <li style={{ marginBottom: 4 }}>• Color-coded health factor display</li>
                <li style={{ marginBottom: 4 }}>• Real-time liquidation probability</li>
                <li style={{ marginBottom: 4 }}>• Distance to liquidation percentage</li>
                <li style={{ marginBottom: 4 }}>• VaR/CVaR metrics (FastAPI)</li>
                <li style={{ marginBottom: 4 }}>• Sharpe ratio (FastAPI)</li>
                <li style={{ marginBottom: 4 }}>• Expandable advanced metrics</li>
                <li style={{ marginBottom: 4 }}>• One-click add collateral button</li>
                <li style={{ marginBottom: 4 }}>• Auto-refresh capability</li>
                <li style={{ marginBottom: 4 }}>• FastAPI fallback mechanism</li>
                <li style={{ marginBottom: 4 }}>• Responsive design</li>
              </ul>
            </div>
            <div>
              <h4 style={{ fontWeight: 600, color: 'rgb(96, 165, 250)', marginBottom: 8 }}>🎨 Design Elements:</h4>
              <ul style={{ listStyle: 'none', padding: 0, margin: 0, color: 'rgb(209, 213, 219)' }}>
                <li style={{ marginBottom: 4 }}>• Health Factor: Green &gt; 2.0, Yellow &gt; 1.5, Orange &gt; 1.2, Red &lt; 1.2</li>
                <li style={{ marginBottom: 4 }}>• Animated pulse on critical alerts</li>
                <li style={{ marginBottom: 4 }}>• Contextual warning messages</li>
                <li style={{ marginBottom: 4 }}>• Data source indicator (FastAPI/Basic)</li>
                <li style={{ marginBottom: 4 }}>• Loading states</li>
                <li style={{ marginBottom: 4 }}>• Error handling display</li>
                <li style={{ marginBottom: 4 }}>• Accessibility-friendly icons</li>
                <li style={{ marginBottom: 4 }}>• Responsive grid layouts</li>
              </ul>
            </div>
          </div>

          <div style={{ marginTop: 24, paddingTop: 16, borderTop: '1px solid rgb(55, 65, 81)' }}>
            <h4 style={{ fontWeight: 600, color: 'rgb(250, 204, 21)', marginBottom: 8 }}>📊 Data Flow:</h4>
            <div style={{ color: 'rgb(209, 213, 219)', fontSize: 12, fontFamily: 'monospace', background: 'rgb(17, 24, 39)', padding: 12, borderRadius: 6, whiteSpace: 'pre-wrap' }}>
              Position Data → LiquidationAlert → checkFastAPIAvailability() →{'\n'}
              [FastAPI Available] → calculateRisk() → Display VaR/CVaR/Sharpe{'\n'}
              [FastAPI Unavailable] → Basic Calculations → Display Health Factor/Liq Prob{'\n'}
              Auto-refresh every refreshInterval ms (configurable)
            </div>
          </div>
        </div>

        {/* Usage Instructions */}
        <div style={{ marginTop: 24, background: 'white', borderRadius: 8, border: '1px solid rgb(229, 231, 235)', padding: 24 }}>
          <h3 style={{ fontWeight: 700, color: 'rgb(17, 24, 39)', marginBottom: 12, fontSize: 16 }}>Usage in Production:</h3>
          <div style={{ fontSize: 14, color: 'rgb(55, 65, 81)' }}>
            <p style={{ marginBottom: 12 }}>
              <strong>1. Import the component:</strong>
            </p>
            <pre style={{ background: 'rgb(249, 250, 251)', padding: 12, borderRadius: 6, overflowX: 'auto', marginBottom: 16 }}>
              <code>{`import LiquidationAlert from './components/LiquidationAlert';`}</code>
            </pre>

            <p style={{ marginBottom: 12 }}>
              <strong>2. Use in your dashboard:</strong>
            </p>
            <pre style={{ background: 'rgb(249, 250, 251)', padding: 12, borderRadius: 6, overflowX: 'auto', marginBottom: 16 }}>
              <code>{`<LiquidationAlert
  position={userPosition}
  marketData={currentMarketData}
  historicalPrices={priceHistory}
  onAddCollateral={handleAddCollateral}
  refreshInterval={60000} // 60 seconds
/>`}</code>
            </pre>

            <p style={{ marginBottom: 12 }}>
              <strong>3. Position data format:</strong>
            </p>
            <pre style={{ background: 'rgb(249, 250, 251)', padding: 12, borderRadius: 6, overflowX: 'auto', fontSize: 12 }}>
              <code>{`{
  id: "position-123",
  protocol: "Aave V3",
  chain: "Arbitrum",
  collateralAsset: "ETH",
  collateralValueUSD: 50000,
  borrowValueUSD: 25000,
  healthFactor: 1.85,
  currentPrice: 2500,
  liquidationPrice: 2000
}`}</code>
            </pre>
          </div>
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
    // Random walk with drift
    const change = (Math.random() - 0.5) * 2 * volatility * price;
    price = Math.max(price + change, basePrice * 0.5); // Don't go below 50% of base

    prices.push({
      timestamp: new Date(Date.now() - i * 3600000).toISOString(), // Hourly data
      price: parseFloat(price.toFixed(2))
    });
  }

  return prices;
}

export default LiquidationAlertDemo;
