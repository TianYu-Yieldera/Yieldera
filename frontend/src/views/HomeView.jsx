import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useWallet } from "../web3/WalletContext";
import { useDemoMode } from "../web3/DemoModeContext";
import { config } from "../config/env";
import AIRiskDashboardPro from "../components/AIRiskDashboardPro";
import HedgeRecommendation from "../components/HedgeRecommendation";
import RiskTrendChart from "../components/RiskTrendChart";
import { DollarSign, Shield, TrendingUp, Wallet, Target, Layers, Activity } from "lucide-react";

const API_BASE = config.api.baseUrl;

function MetricCard({ icon: Icon, title, value, subtitle, trend, onClick, clickable }) {
  return (
    <div style={{
      background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
      borderRadius: 12,
      padding: 24,
      border: '1px solid rgba(34, 211, 238, 0.2)',
      transition: 'all 0.3s ease',
      cursor: clickable ? 'pointer' : 'default',
      position: 'relative',
      overflow: 'hidden'
    }}
    onClick={onClick}
    onMouseEnter={(e) => {
      e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.4)';
      e.currentTarget.style.transform = 'translateY(-4px)';
      e.currentTarget.style.boxShadow = clickable ? '0 8px 24px rgba(34, 211, 238, 0.3)' : '0 8px 24px rgba(0,0,0,0.3)';
    }}
    onMouseLeave={(e) => {
      e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.2)';
      e.currentTarget.style.transform = 'translateY(0)';
      e.currentTarget.style.boxShadow = 'none';
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

      {/* Gradient accent */}
      <div style={{
        position: 'absolute',
        top: -50,
        right: -50,
        width: 150,
        height: 150,
        background: 'radial-gradient(circle, rgba(34, 211, 238, 0.15) 0%, transparent 70%)',
        opacity: 0.6
      }} />

      <div style={{ position: 'relative', zIndex: 1 }}>
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
          <div style={{
            fontSize: 12,
            fontWeight: 600,
            color: 'rgba(203, 213, 225, 0.7)',
            textTransform: 'uppercase',
            letterSpacing: 1
          }}>
            {title}
          </div>
          <div style={{
            width: 36,
            height: 36,
            borderRadius: 8,
            background: 'rgba(34, 211, 238, 0.15)',
            border: '1px solid rgba(34, 211, 238, 0.3)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
          }}>
            <Icon style={{ width: 18, height: 18, color: 'rgb(34, 211, 238)' }} />
          </div>
        </div>
        <div style={{
          fontSize: 28,
          fontWeight: 700,
          color: 'white',
          marginBottom: 8,
          textShadow: '0 0 20px rgba(255, 255, 255, 0.1)'
        }}>
          {value}
        </div>
        {subtitle && (
          <div style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.8)', display: 'flex', alignItems: 'center', gap: 6 }}>
            {trend && (
              <span style={{
                color: trend > 0 ? 'rgb(34, 197, 94)' : trend < 0 ? 'rgb(239, 68, 68)' : 'rgba(203, 213, 225, 0.6)',
                fontWeight: 700,
                fontSize: 14
              }}>
                {trend > 0 ? '↑' : trend < 0 ? '↓' : '•'} {Math.abs(trend)}%
              </span>
            )}
            <span>{subtitle}</span>
          </div>
        )}
      </div>
    </div>
  );
}

export default function HomeView() {
  const navigate = useNavigate();
  const { address } = useWallet();
  const { demoMode, demoAddress, demoData } = useDemoMode();
  const [data, setData] = useState({
    portfolio: { totalValue: 0, dailyChange: 0, changePercent: 0 },
    treasury: { holdings: 0, apy: 0 },
    defi: { tvl: 0, positions: 0, avgApy: 0 },
    risk: { score: 0, level: 'Low' },
    protocols: { active: 0, total: 0 }
  });

  useEffect(() => {
    if (demoMode) {
      // Demo mode: calculate from mock data
      const defiTVL = demoData.defiDeposits.reduce((sum, d) => sum + d.amount, 0);
      const treasuryHoldings = demoData.collateral || 0; // Assuming collateral represents treasury holdings
      const totalValue = demoData.pfiTokens + defiTVL + treasuryHoldings + demoData.stakedTokens;

      // Calculate weighted average APY
      const totalApy = demoData.defiDeposits.reduce((sum, d) => sum + (d.apy * d.amount), 0);
      const avgApy = defiTVL > 0 ? (totalApy / defiTVL) : 0;

      // Mock risk score based on positions
      const riskScore = Math.min(75, 20 + (demoData.defiDeposits.length * 10));

      setData({
        portfolio: { totalValue, dailyChange: 245.50, changePercent: 0.8 },
        treasury: { holdings: treasuryHoldings, apy: 4.5 },
        defi: { tvl: defiTVL, positions: demoData.defiDeposits.length, avgApy },
        risk: { score: riskScore, level: riskScore > 70 ? 'Medium' : 'Low' },
        protocols: { active: demoData.defiDeposits.length, total: 8 }
      });
    } else {
      // Real mode: fetch from API
      const addr = address || "0x3C07226A3f1488320426eB5FE9976f72E5712346";
      const j = (p) => fetch(`${API_BASE}${p}`).then(r => r.json()).catch(() => null);

      Promise.all([
        j(`/users/${addr}/balance`)
      ]).then(([bal]) => {
        const balance = Number(bal?.balance || 0);

        setData({
          portfolio: { totalValue: balance, dailyChange: 0, changePercent: 0 },
          treasury: { holdings: 0, apy: 4.5 },
          defi: { tvl: 0, positions: 0, avgApy: 0 },
          risk: { score: 0, level: 'Low' },
          protocols: { active: 0, total: 8 }
        });
      });
    }
  }, [address, demoMode, demoData]);

  const displayAddress = demoMode ? demoAddress : address;

  return (
    <div className="container">
      {/* Simplified Status Bar */}
      {demoMode && (
        <div style={{
          display: 'flex',
          justifyContent: 'flex-end',
          marginBottom: 20
        }}>
          <span style={{
            fontSize: '12px',
            background: 'linear-gradient(135deg, rgb(34, 197, 94) 0%, rgb(22, 163, 74) 100%)',
            color: 'white',
            padding: '6px 14px',
            borderRadius: '6px',
            fontWeight: 600,
            letterSpacing: 0.5,
            boxShadow: '0 0 15px rgba(34, 197, 94, 0.3)'
          }}>
            DEMO MODE
          </span>
        </div>
      )}

      {/* Metrics Grid */}
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
        gap: 20,
        marginBottom: 32
      }}>
        <MetricCard
          icon={DollarSign}
          title="Portfolio Value"
          value={`$${data.portfolio.totalValue.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`}
          subtitle="Total assets under management"
          trend={data.portfolio.changePercent}
        />
        <MetricCard
          icon={Shield}
          title="Treasury Holdings"
          value={`$${data.treasury.holdings.toLocaleString()}`}
          subtitle={`${data.treasury.apy}% APY · Zero-risk yield · Click to view`}
          onClick={() => navigate('/treasury/holdings')}
          clickable={true}
        />
        <MetricCard
          icon={Wallet}
          title="DeFi TVL"
          value={`$${data.defi.tvl.toLocaleString()}`}
          subtitle={data.defi.positions > 0 ? `${data.defi.positions} positions · ${data.defi.avgApy.toFixed(2)}% avg APY` : 'No active positions'}
        />
        <MetricCard
          icon={Target}
          title="AI Risk Score"
          value={data.risk.score.toFixed(0)}
          subtitle={`${data.risk.level} risk level · Real-time monitoring`}
        />
        <MetricCard
          icon={TrendingUp}
          title="Portfolio APY"
          value={data.defi.avgApy > 0 ? `${data.defi.avgApy.toFixed(2)}%` : '0%'}
          subtitle="Weighted average yield"
        />
        <MetricCard
          icon={Layers}
          title="Active Protocols"
          value={`${data.protocols.active} / ${data.protocols.total}`}
          subtitle="Multi-protocol diversification"
        />
      </div>

      {/* AI Risk Dashboard - Professional Version */}
      <div style={{ marginBottom: 32 }}>
        <AIRiskDashboardPro />
      </div>

      {/* Risk & Yield Trends - Original Version */}
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(400px, 1fr))',
        gap: 20,
        marginBottom: 32
      }}>
        <RiskTrendChart type="risk" />
        <RiskTrendChart type="yield" />
      </div>

      {/* Hedge Recommendations */}
      <div style={{ marginBottom: 32 }}>
        <HedgeRecommendation
          riskScore={data.risk.score}
          portfolioValue={data.portfolio.totalValue}
          volatility={demoMode ? 24 : 0}
        />
      </div>

      {/* Recent Activity */}
      <div style={{
        marginTop: 32,
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
          <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 16 }}>
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
              <Activity style={{ width: 20, height: 20, color: 'rgb(34, 211, 238)' }} />
            </div>
            <div>
              <h3 style={{ fontSize: 20, fontWeight: 700, color: 'white', margin: 0 }}>
                Recent Activity
              </h3>
              <p style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.7)', margin: '2px 0 0 0' }}>
                On-chain transaction history
              </p>
            </div>
          </div>

          <div style={{
            padding: 24,
            background: 'rgba(255, 255, 255, 0.03)',
            borderRadius: 12,
            border: '1px solid rgba(255, 255, 255, 0.08)',
            textAlign: 'center'
          }}>
            <p style={{
              fontSize: 14,
              color: 'rgba(203, 213, 225, 0.8)',
              margin: 0
            }}>
              On-chain transaction tracking coming soon
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
