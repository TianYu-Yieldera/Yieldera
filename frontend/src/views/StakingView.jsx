import React, { useState } from "react";
import { Lock, Unlock, TrendingUp, Coins, Clock, Zap } from "lucide-react";
import { useWallet } from "../web3/WalletContext";

export default function StakingView() {
  const { address } = useWallet();
  const [selectedPool, setSelectedPool] = useState(null);

  // 模拟的质押池数据
  const pools = [
    {
      id: 1,
      name: "LoyaltyX Token",
      symbol: "PFI",
      apr: "125.5",
      tvl: "12,450,000",
      yourStake: "0",
      rewards: "0",
      lockPeriod: "无锁定",
      icon: "🎯",
      color: "#6366F1"
    },
    {
      id: 2,
      name: "ETH-PFI LP",
      symbol: "LP",
      apr: "245.8",
      tvl: "8,920,000",
      yourStake: "0",
      rewards: "0",
      lockPeriod: "7 天",
      icon: "💎",
      color: "#A855F7"
    },
    {
      id: 3,
      name: "USDC Pool",
      symbol: "USDC",
      apr: "85.2",
      tvl: "25,680,000",
      yourStake: "0",
      rewards: "0",
      lockPeriod: "无锁定",
      icon: "💵",
      color: "#22D3EE"
    },
    {
      id: 4,
      name: "Long Term Vault",
      symbol: "PFI",
      apr: "380.0",
      tvl: "5,120,000",
      yourStake: "0",
      rewards: "0",
      lockPeriod: "90 天",
      icon: "🔒",
      color: "#F59E0B"
    }
  ];

  const stats = [
    { label: "总锁仓价值", value: "$52.17M", icon: Lock, color: "#6366F1" },
    { label: "我的质押", value: "$0", icon: Coins, color: "#A855F7" },
    { label: "待领取奖励", value: "0 PFI", icon: TrendingUp, color: "#22D3EE" },
    { label: "年化收益", value: "183.6%", icon: Zap, color: "#F59E0B" }
  ];

  return (
    <div className="container">
      <div style={{ marginBottom: 24 }}>
        <h1 style={{ margin: 0, fontSize: 32, display: 'flex', alignItems: 'center', gap: 12 }}>
          <Lock size={36} color="#6366F1" />
          质押池
        </h1>
        <p className="muted" style={{ marginTop: 8 }}>质押代币赚取收益和积分奖励</p>
      </div>

      {/* 统计卡片 */}
      <div className="grid" style={{ gridTemplateColumns: 'repeat(4, 1fr)', gap: 16, marginBottom: 24 }}>
        {stats.map((stat, index) => (
          <div key={index} className="kpi" style={{ position: 'relative', overflow: 'hidden' }}>
            <div style={{ position: 'absolute', top: -10, right: -10, opacity: 0.1 }}>
              <stat.icon size={80} color={stat.color} />
            </div>
            <div style={{ position: 'relative', zIndex: 1 }}>
              <div className="title">{stat.label}</div>
              <div className="value" style={{ color: stat.color }}>{stat.value}</div>
            </div>
          </div>
        ))}
      </div>

      {/* 质押池列表 */}
      <div className="grid" style={{ gridTemplateColumns: 'repeat(2, 1fr)', gap: 16 }}>
        {pools.map((pool) => (
          <div
            key={pool.id}
            className="card"
            style={{
              padding: 24,
              cursor: 'pointer',
              transition: 'all 0.3s',
              borderColor: selectedPool?.id === pool.id ? pool.color : 'rgba(255,255,255,.1)',
              borderWidth: 2,
              position: 'relative',
              overflow: 'hidden'
            }}
            onClick={() => setSelectedPool(pool)}
            onMouseEnter={(e) => {
              e.currentTarget.style.transform = 'translateY(-4px)';
              e.currentTarget.style.boxShadow = `0 20px 40px rgba(0,0,0,.3)`;
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = '0 10px 24px rgba(0,0,0,.2)';
            }}
          >
            <div style={{ position: 'absolute', top: -20, right: -20, fontSize: 120, opacity: 0.05 }}>
              {pool.icon}
            </div>

            <div style={{ position: 'relative', zIndex: 1 }}>
              {/* 池子标题 */}
              <div className="row" style={{ justifyContent: 'space-between', marginBottom: 16 }}>
                <div className="row" style={{ gap: 12 }}>
                  <div style={{
                    fontSize: 40,
                    width: 60,
                    height: 60,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    borderRadius: 12,
                    background: `${pool.color}22`
                  }}>
                    {pool.icon}
                  </div>
                  <div>
                    <div style={{ fontWeight: 700, fontSize: 18 }}>{pool.name}</div>
                    <div className="muted">{pool.symbol}</div>
                  </div>
                </div>
                <div style={{
                  background: `${pool.color}22`,
                  color: pool.color,
                  padding: '4px 12px',
                  borderRadius: 8,
                  fontSize: 12,
                  fontWeight: 700,
                  height: 'fit-content'
                }}>
                  <Clock size={12} style={{ display: 'inline', marginRight: 4 }} />
                  {pool.lockPeriod}
                </div>
              </div>

              {/* APR 展示 */}
              <div style={{
                background: `linear-gradient(135deg, ${pool.color}33, ${pool.color}11)`,
                padding: 16,
                borderRadius: 12,
                marginBottom: 16
              }}>
                <div className="row" style={{ justifyContent: 'space-between', alignItems: 'flex-end' }}>
                  <div>
                    <div className="muted" style={{ fontSize: 12 }}>年化收益率</div>
                    <div style={{ fontSize: 32, fontWeight: 900, color: pool.color }}>
                      {pool.apr}%
                    </div>
                  </div>
                  <div style={{ textAlign: 'right' }}>
                    <div className="muted" style={{ fontSize: 12 }}>TVL</div>
                    <div style={{ fontSize: 16, fontWeight: 700 }}>${pool.tvl}</div>
                  </div>
                </div>
              </div>

              {/* 用户数据 */}
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12, marginBottom: 16 }}>
                <div>
                  <div className="muted" style={{ fontSize: 12 }}>我的质押</div>
                  <div style={{ fontWeight: 700 }}>{pool.yourStake} {pool.symbol}</div>
                </div>
                <div>
                  <div className="muted" style={{ fontSize: 12 }}>待领取奖励</div>
                  <div style={{ fontWeight: 700, color: '#10b981' }}>{pool.rewards} {pool.symbol}</div>
                </div>
              </div>

              {/* 操作按钮 */}
              <div className="row" style={{ gap: 8 }}>
                <button
                  className="btn"
                  style={{ flex: 1, background: pool.color }}
                  onClick={(e) => {
                    e.stopPropagation();
                    alert('质押功能开发中...');
                  }}
                >
                  <Lock size={16} style={{ display: 'inline', marginRight: 4 }} />
                  质押
                </button>
                <button
                  className="btn"
                  style={{ flex: 1, background: 'rgba(255,255,255,.1)' }}
                  onClick={(e) => {
                    e.stopPropagation();
                    alert('解除质押功能开发中...');
                  }}
                >
                  <Unlock size={16} style={{ display: 'inline', marginRight: 4 }} />
                  解除
                </button>
                <button
                  className="btn"
                  style={{ background: '#10b981' }}
                  onClick={(e) => {
                    e.stopPropagation();
                    alert('领取奖励功能开发中...');
                  }}
                >
                  领取
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* 连接钱包提示 */}
      {!address && (
        <div className="card" style={{ marginTop: 24, padding: 24, textAlign: 'center', background: 'rgba(245, 158, 11, .1)', borderColor: '#F59E0B' }}>
          <Lock size={32} color="#F59E0B" style={{ margin: '0 auto 12px' }} />
          <div style={{ fontWeight: 700, marginBottom: 8 }}>请先连接钱包</div>
          <div className="muted">连接钱包后即可查看并管理你的质押</div>
        </div>
      )}
    </div>
  );
}
