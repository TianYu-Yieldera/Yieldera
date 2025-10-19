import React, { useEffect, useState } from "react";
import { useWallet } from "../web3/WalletContext";
import {
  DollarSign,
  TrendingUp,
  Shield,
  Droplet,
  ArrowDownUp,
  AlertCircle,
  CheckCircle2,
  Info,
  Lock,
  Unlock,
  Zap,
  Activity,
  Target
} from "lucide-react";

import { config } from "../config/env";

const API_BASE = config.api.baseUrl;

export default function StablecoinView() {
  const { address } = useWallet();
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState("mint"); // mint, redeem, liquidity

  // User data
  const [userPoints, setUserPoints] = useState(0);
  const [userBalance, setUserBalance] = useState(0);

  // Stablecoin position
  const [position, setPosition] = useState({
    collateral: 0,
    debt: 0,
    collateralRatio: 0,
    maxMintable: 0,
    liquidationPrice: 0,
    healthStatus: "safe" // safe, warning, danger
  });

  // Stablecoin stats
  const [stats, setStats] = useState({
    totalSupply: "1,245,678",
    lusdPrice: "1.00",
    deviation: "+0.02%",
    totalCollateral: "2,890,456",
    avgCollateralRatio: "185%"
  });

  // LP position
  const [lpPosition, setLpPosition] = useState({
    lpTokens: 0,
    shareOfPool: "0%",
    rewardsEarned: 0,
    estimatedAPR: "18.5%"
  });

  // Form state
  const [mintAmount, setMintAmount] = useState("");
  const [redeemAmount, setRedeemAmount] = useState("");
  const [liquidityAmount, setLiquidityAmount] = useState("");

  useEffect(() => {
    if (!address) {
      setLoading(false);
      return;
    }

    // Fetch user data
    Promise.all([
      fetch(`${API_BASE}/users/${address}/points`).then(r => r.json()),
      fetch(`${API_BASE}/users/${address}/balance`).then(r => r.json())
    ])
      .then(([pointsData, balanceData]) => {
        setUserPoints(pointsData.points || 0);
        setUserBalance(parseFloat(balanceData.balance) || 0);
        setLoading(false);
      })
      .catch(() => {
        setLoading(false);
      });
  }, [address]);

  const calculateMintPreview = () => {
    if (!mintAmount || parseFloat(mintAmount) <= 0) return null;

    const amount = parseFloat(mintAmount);
    const fee = amount * 0.002; // 0.2% fee
    const netAmount = amount - fee;
    const newDebt = position.debt + amount;
    const newRatio = position.collateral > 0 ? (position.collateral * 100 / newDebt).toFixed(0) : 0;

    return {
      gross: amount,
      fee: fee.toFixed(2),
      net: netAmount.toFixed(2),
      newRatio: newRatio
    };
  };

  const calculateRedeemPreview = () => {
    if (!redeemAmount || parseFloat(redeemAmount) <= 0) return null;

    const amount = parseFloat(redeemAmount);
    const fee = amount * 0.002; // 0.2% fee
    const collateralReturned = amount - fee;
    const newDebt = Math.max(0, position.debt - amount);
    const newRatio = newDebt > 0 ? (position.collateral * 100 / newDebt).toFixed(0) : 0;

    return {
      gross: amount,
      fee: fee.toFixed(2),
      collateralReturned: collateralReturned.toFixed(2),
      newRatio: newRatio
    };
  };

  const getHealthColor = (status) => {
    switch (status) {
      case "safe": return "#10b981";
      case "warning": return "#F59E0B";
      case "danger": return "#EF4444";
      default: return "#666";
    }
  };

  const getHealthStatus = (ratio) => {
    if (ratio >= 150) return "safe";
    if (ratio >= 120) return "warning";
    return "danger";
  };

  const handleMint = () => {
    if (!address) {
      alert("请先连接钱包");
      return;
    }

    const preview = calculateMintPreview();
    if (!preview) {
      alert("请输入有效金额");
      return;
    }

    alert(`铸造 ${preview.net} LUSD\n手续费: ${preview.fee} LUSD\n新抵押率: ${preview.newRatio}%\n\n智能合约功能开发中...`);
  };

  const handleRedeem = () => {
    if (!address) {
      alert("请先连接钱包");
      return;
    }

    const preview = calculateRedeemPreview();
    if (!preview) {
      alert("请输入有效金额");
      return;
    }

    alert(`赎回 ${preview.collateralReturned} LP\n手续费: ${preview.fee} LUSD\n新抵押率: ${preview.newRatio}%\n\n智能合约功能开发中...`);
  };

  if (loading) {
    return (
      <div className="container" style={{ textAlign: 'center', padding: 60 }}>
        <div className="muted">加载中...</div>
      </div>
    );
  }

  if (!address) {
    return (
      <div className="container" style={{ textAlign: 'center', padding: 60 }}>
        <DollarSign size={48} color="#10b981" style={{ marginBottom: 16 }} />
        <h2>LoyaltyUSD 稳定币</h2>
        <p className="muted">连接钱包后即可铸造稳定币</p>
      </div>
    );
  }

  return (
    <div className="container">
      {/* Header */}
      <div style={{ marginBottom: 24 }}>
        <h1 style={{ margin: 0, fontSize: 32, display: 'flex', alignItems: 'center', gap: 12 }}>
          <DollarSign size={36} style={{ color: '#10b981' }} />
          LoyaltyUSD 稳定币
        </h1>
        <p className="muted" style={{ marginTop: 8 }}>
          用 Loyalty Points 作为抵押物铸造与美元 1:1 挂钩的稳定币
        </p>
      </div>

      {/* Global Stats */}
      <div className="grid" style={{ gridTemplateColumns: 'repeat(5, 1fr)', gap: 16, marginBottom: 24 }}>
        <div className="kpi">
          <div className="title">LUSD 价格</div>
          <div className="value" style={{ color: '#10b981' }}>${stats.lusdPrice}</div>
          <div style={{ fontSize: 12, color: stats.deviation.startsWith('+') ? '#10b981' : '#EF4444', marginTop: 4 }}>
            {stats.deviation}
          </div>
        </div>
        <div className="kpi">
          <div className="title">总供应量</div>
          <div className="value">{stats.totalSupply}</div>
          <div className="muted" style={{ marginTop: 4 }}>LUSD</div>
        </div>
        <div className="kpi">
          <div className="title">总抵押物</div>
          <div className="value" style={{ color: '#6366F1' }}>{stats.totalCollateral}</div>
          <div className="muted" style={{ marginTop: 4 }}>LP Tokens</div>
        </div>
        <div className="kpi">
          <div className="title">平均抵押率</div>
          <div className="value" style={{ color: '#A855F7' }}>{stats.avgCollateralRatio}</div>
          <div className="muted" style={{ marginTop: 4 }}>全网平均</div>
        </div>
        <div className="kpi">
          <div className="title">我的积分</div>
          <div className="value" style={{ color: '#F59E0B' }}>{userPoints.toLocaleString()}</div>
          <div className="muted" style={{ marginTop: 4 }}>可用抵押</div>
        </div>
      </div>

      <div className="grid grid-2" style={{ gap: 24, marginBottom: 24 }}>
        {/* Left: Mint/Redeem/Liquidity */}
        <div className="card" style={{ padding: 24 }}>
          {/* Tabs */}
          <div className="row" style={{ gap: 8, marginBottom: 24 }}>
            <button
              className="btn"
              style={{
                flex: 1,
                background: activeTab === "mint" ? '#10b981' : 'rgba(255,255,255,.05)',
                borderColor: activeTab === "mint" ? '#10b981' : 'rgba(255,255,255,.1)'
              }}
              onClick={() => setActiveTab("mint")}
            >
              <Lock size={16} style={{ marginRight: 6 }} />
              铸造 LUSD
            </button>
            <button
              className="btn"
              style={{
                flex: 1,
                background: activeTab === "redeem" ? '#6366F1' : 'rgba(255,255,255,.05)',
                borderColor: activeTab === "redeem" ? '#6366F1' : 'rgba(255,255,255,.1)'
              }}
              onClick={() => setActiveTab("redeem")}
            >
              <Unlock size={16} style={{ marginRight: 6 }} />
              赎回抵押
            </button>
            <button
              className="btn"
              style={{
                flex: 1,
                background: activeTab === "liquidity" ? '#A855F7' : 'rgba(255,255,255,.05)',
                borderColor: activeTab === "liquidity" ? '#A855F7' : 'rgba(255,255,255,.1)'
              }}
              onClick={() => setActiveTab("liquidity")}
            >
              <Droplet size={16} style={{ marginRight: 6 }} />
              流动性挖矿
            </button>
          </div>

          {/* Mint Tab */}
          {activeTab === "mint" && (
            <div>
              <div style={{ marginBottom: 16 }}>
                <div className="row" style={{ justifyContent: 'space-between', marginBottom: 8 }}>
                  <label style={{ fontSize: 14, fontWeight: 600 }}>抵押 Loyalty Points</label>
                  <span className="muted" style={{ fontSize: 12 }}>
                    可用: {userPoints.toLocaleString()} LP
                  </span>
                </div>
                <input
                  type="number"
                  placeholder="输入抵押数量"
                  value={mintAmount}
                  onChange={(e) => setMintAmount(e.target.value)}
                  style={{
                    width: '100%',
                    padding: '12px 16px',
                    background: 'rgba(0,0,0,.3)',
                    border: '1px solid rgba(255,255,255,.1)',
                    borderRadius: 12,
                    color: '#fff',
                    fontSize: 16
                  }}
                />
              </div>

              {calculateMintPreview() && (
                <div style={{
                  padding: 16,
                  background: 'rgba(16, 185, 129, .1)',
                  border: '1px solid rgba(16, 185, 129, .3)',
                  borderRadius: 12,
                  marginBottom: 16
                }}>
                  <div style={{ fontWeight: 700, marginBottom: 12, display: 'flex', alignItems: 'center', gap: 8 }}>
                    <Info size={16} color="#10b981" />
                    铸造预览
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 8, fontSize: 14 }}>
                    <div className="row" style={{ justifyContent: 'space-between' }}>
                      <span className="muted">铸造金额:</span>
                      <span>{calculateMintPreview().gross} LUSD</span>
                    </div>
                    <div className="row" style={{ justifyContent: 'space-between' }}>
                      <span className="muted">手续费 (0.2%):</span>
                      <span style={{ color: '#F59E0B' }}>-{calculateMintPreview().fee} LUSD</span>
                    </div>
                    <div className="row" style={{ justifyContent: 'space-between' }}>
                      <span style={{ fontWeight: 700 }}>实际收到:</span>
                      <span style={{ fontWeight: 700, color: '#10b981' }}>{calculateMintPreview().net} LUSD</span>
                    </div>
                    <div style={{ borderTop: '1px solid rgba(255,255,255,.1)', paddingTop: 8, marginTop: 4 }} />
                    <div className="row" style={{ justifyContent: 'space-between' }}>
                      <span className="muted">新抵押率:</span>
                      <span style={{ fontWeight: 700, color: calculateMintPreview().newRatio >= 150 ? '#10b981' : '#EF4444' }}>
                        {calculateMintPreview().newRatio}%
                      </span>
                    </div>
                  </div>
                </div>
              )}

              <button
                className="btn"
                style={{
                  width: '100%',
                  background: '#10b981',
                  fontSize: 16,
                  padding: '14px 20px'
                }}
                onClick={handleMint}
              >
                <Lock size={18} style={{ marginRight: 8 }} />
                铸造 LUSD
              </button>

              <div style={{
                marginTop: 16,
                padding: 12,
                background: 'rgba(99, 102, 241, .1)',
                borderRadius: 8,
                fontSize: 13
              }}>
                <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
                  <Info size={14} color="#6366F1" style={{ flexShrink: 0, marginTop: 2 }} />
                  <div>
                    <div style={{ fontWeight: 600, marginBottom: 4 }}>如何铸造?</div>
                    <ul style={{ margin: 0, paddingLeft: 16, lineHeight: 1.6 }} className="muted">
                      <li>存入 Loyalty Points 作为抵押</li>
                      <li>保持 ≥150% 抵押率</li>
                      <li>支付 0.2% 铸造费用</li>
                      <li>获得等值 LUSD 稳定币</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Redeem Tab */}
          {activeTab === "redeem" && (
            <div>
              <div style={{ marginBottom: 16 }}>
                <div className="row" style={{ justifyContent: 'space-between', marginBottom: 8 }}>
                  <label style={{ fontSize: 14, fontWeight: 600 }}>赎回 LUSD</label>
                  <span className="muted" style={{ fontSize: 12 }}>
                    债务: {position.debt.toLocaleString()} LUSD
                  </span>
                </div>
                <input
                  type="number"
                  placeholder="输入赎回数量"
                  value={redeemAmount}
                  onChange={(e) => setRedeemAmount(e.target.value)}
                  style={{
                    width: '100%',
                    padding: '12px 16px',
                    background: 'rgba(0,0,0,.3)',
                    border: '1px solid rgba(255,255,255,.1)',
                    borderRadius: 12,
                    color: '#fff',
                    fontSize: 16
                  }}
                />
              </div>

              {calculateRedeemPreview() && (
                <div style={{
                  padding: 16,
                  background: 'rgba(99, 102, 241, .1)',
                  border: '1px solid rgba(99, 102, 241, .3)',
                  borderRadius: 12,
                  marginBottom: 16
                }}>
                  <div style={{ fontWeight: 700, marginBottom: 12, display: 'flex', alignItems: 'center', gap: 8 }}>
                    <Info size={16} color="#6366F1" />
                    赎回预览
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 8, fontSize: 14 }}>
                    <div className="row" style={{ justifyContent: 'space-between' }}>
                      <span className="muted">赎回金额:</span>
                      <span>{calculateRedeemPreview().gross} LUSD</span>
                    </div>
                    <div className="row" style={{ justifyContent: 'space-between' }}>
                      <span className="muted">手续费 (0.2%):</span>
                      <span style={{ color: '#F59E0B' }}>-{calculateRedeemPreview().fee} LUSD</span>
                    </div>
                    <div className="row" style={{ justifyContent: 'space-between' }}>
                      <span style={{ fontWeight: 700 }}>返还抵押物:</span>
                      <span style={{ fontWeight: 700, color: '#10b981' }}>{calculateRedeemPreview().collateralReturned} LP</span>
                    </div>
                    <div style={{ borderTop: '1px solid rgba(255,255,255,.1)', paddingTop: 8, marginTop: 4 }} />
                    <div className="row" style={{ justifyContent: 'space-between' }}>
                      <span className="muted">新抵押率:</span>
                      <span style={{ fontWeight: 700, color: '#10b981' }}>
                        {calculateRedeemPreview().newRatio > 0 ? `${calculateRedeemPreview().newRatio}%` : 'N/A'}
                      </span>
                    </div>
                  </div>
                </div>
              )}

              <button
                className="btn"
                style={{
                  width: '100%',
                  background: '#6366F1',
                  fontSize: 16,
                  padding: '14px 20px'
                }}
                onClick={handleRedeem}
              >
                <Unlock size={18} style={{ marginRight: 8 }} />
                赎回抵押物
              </button>

              <div style={{
                marginTop: 16,
                padding: 12,
                background: 'rgba(99, 102, 241, .1)',
                borderRadius: 8,
                fontSize: 13
              }}>
                <div style={{ display: 'flex', gap: 8 }}>
                  <Info size={14} color="#6366F1" style={{ flexShrink: 0, marginTop: 2 }} />
                  <div>
                    <div style={{ fontWeight: 600, marginBottom: 4 }}>赎回说明</div>
                    <ul style={{ margin: 0, paddingLeft: 16, lineHeight: 1.6 }} className="muted">
                      <li>燃烧 LUSD 以取回抵押的 LP</li>
                      <li>支付 0.2% 赎回费用</li>
                      <li>按比例返还抵押物</li>
                      <li>可部分或全额赎回</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Liquidity Tab */}
          {activeTab === "liquidity" && (
            <div>
              <div style={{
                padding: 16,
                background: 'linear-gradient(135deg, rgba(168, 85, 247, .15), rgba(236, 72, 153, .15))',
                borderRadius: 12,
                marginBottom: 20
              }}>
                <div className="row" style={{ gap: 12, marginBottom: 12 }}>
                  <Droplet size={24} color="#A855F7" />
                  <div>
                    <div style={{ fontWeight: 700, fontSize: 16 }}>LUSD-USDC 流动性池</div>
                    <div className="muted" style={{ fontSize: 13 }}>Uniswap V3 Pool</div>
                  </div>
                </div>
                <div className="row" style={{ gap: 24, marginTop: 16 }}>
                  <div>
                    <div className="muted" style={{ fontSize: 12 }}>年化收益率</div>
                    <div style={{ fontSize: 24, fontWeight: 800, color: '#A855F7' }}>{lpPosition.estimatedAPR}</div>
                  </div>
                  <div>
                    <div className="muted" style={{ fontSize: 12 }}>池子 TVL</div>
                    <div style={{ fontSize: 18, fontWeight: 700 }}>$2.4M</div>
                  </div>
                  <div>
                    <div className="muted" style={{ fontSize: 12 }}>24h 交易量</div>
                    <div style={{ fontSize: 18, fontWeight: 700 }}>$156K</div>
                  </div>
                </div>
              </div>

              <div style={{ marginBottom: 16 }}>
                <div className="row" style={{ justifyContent: 'space-between', marginBottom: 8 }}>
                  <label style={{ fontSize: 14, fontWeight: 600 }}>添加流动性</label>
                  <span className="muted" style={{ fontSize: 12 }}>
                    比例: 50% LUSD + 50% USDC
                  </span>
                </div>
                <input
                  type="number"
                  placeholder="输入 LUSD 数量"
                  value={liquidityAmount}
                  onChange={(e) => setLiquidityAmount(e.target.value)}
                  style={{
                    width: '100%',
                    padding: '12px 16px',
                    background: 'rgba(0,0,0,.3)',
                    border: '1px solid rgba(255,255,255,.1)',
                    borderRadius: 12,
                    color: '#fff',
                    fontSize: 16,
                    marginBottom: 8
                  }}
                />
                <div className="muted" style={{ fontSize: 12, marginLeft: 4 }}>
                  需要同时提供 {liquidityAmount || "0"} USDC
                </div>
              </div>

              <div className="row" style={{ gap: 8, marginBottom: 16 }}>
                <button
                  className="btn"
                  style={{
                    flex: 1,
                    background: '#A855F7',
                    fontSize: 14
                  }}
                >
                  <Droplet size={16} style={{ marginRight: 6 }} />
                  添加流动性
                </button>
                <button
                  className="btn"
                  style={{
                    flex: 1,
                    background: 'rgba(255,255,255,.1)',
                    fontSize: 14
                  }}
                >
                  移除流动性
                </button>
              </div>

              {/* LP Position Info */}
              <div style={{
                padding: 16,
                background: 'rgba(255,255,255,.02)',
                borderRadius: 12,
                marginBottom: 16
              }}>
                <div style={{ fontWeight: 600, marginBottom: 12 }}>我的 LP 仓位</div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: 8, fontSize: 14 }}>
                  <div className="row" style={{ justifyContent: 'space-between' }}>
                    <span className="muted">LP Tokens:</span>
                    <span>{lpPosition.lpTokens.toLocaleString()}</span>
                  </div>
                  <div className="row" style={{ justifyContent: 'space-between' }}>
                    <span className="muted">池子份额:</span>
                    <span>{lpPosition.shareOfPool}</span>
                  </div>
                  <div className="row" style={{ justifyContent: 'space-between' }}>
                    <span className="muted">已赚取奖励:</span>
                    <span style={{ color: '#10b981' }}>{lpPosition.rewardsEarned.toFixed(2)} PFI</span>
                  </div>
                </div>
                <button
                  className="btn"
                  style={{
                    width: '100%',
                    marginTop: 12,
                    background: '#10b981'
                  }}
                >
                  <Zap size={16} style={{ marginRight: 6 }} />
                  领取奖励
                </button>
              </div>

              <div style={{
                padding: 12,
                background: 'rgba(168, 85, 247, .1)',
                borderRadius: 8,
                fontSize: 13
              }}>
                <div style={{ display: 'flex', gap: 8 }}>
                  <Info size={14} color="#A855F7" style={{ flexShrink: 0, marginTop: 2 }} />
                  <div>
                    <div style={{ fontWeight: 600, marginBottom: 4 }}>流动性挖矿收益</div>
                    <ul style={{ margin: 0, paddingLeft: 16, lineHeight: 1.6 }} className="muted">
                      <li>交易手续费分成 (0.3%)</li>
                      <li>流动性挖矿奖励</li>
                      <li>PFI 代币激励</li>
                      <li>随时添加或移除</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Right: Position Details */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
          {/* Health Status */}
          <div className="card" style={{
            padding: 24,
            background: `linear-gradient(135deg, ${getHealthColor(position.healthStatus)}15, rgba(255,255,255,.02))`,
            borderColor: `${getHealthColor(position.healthStatus)}55`
          }}>
            <div style={{ fontWeight: 700, marginBottom: 16, fontSize: 18, display: 'flex', alignItems: 'center', gap: 8 }}>
              <Shield size={24} color={getHealthColor(position.healthStatus)} />
              仓位健康度
            </div>

            {/* Health Meter */}
            <div style={{ marginBottom: 20 }}>
              <div className="row" style={{ justifyContent: 'space-between', marginBottom: 8 }}>
                <span className="muted" style={{ fontSize: 13 }}>抵押率</span>
                <span style={{ fontSize: 20, fontWeight: 800, color: getHealthColor(position.healthStatus) }}>
                  {position.collateralRatio > 0 ? `${position.collateralRatio}%` : 'N/A'}
                </span>
              </div>
              <div style={{
                width: '100%',
                height: 12,
                background: 'rgba(0,0,0,.3)',
                borderRadius: 6,
                overflow: 'hidden',
                position: 'relative'
              }}>
                <div style={{
                  width: `${Math.min(100, (position.collateralRatio / 300) * 100)}%`,
                  height: '100%',
                  background: `linear-gradient(90deg, ${getHealthColor(position.healthStatus)}, ${getHealthColor(position.healthStatus)}88)`,
                  transition: 'width 0.3s'
                }} />
                {/* Markers */}
                <div style={{
                  position: 'absolute',
                  left: '40%',
                  top: 0,
                  bottom: 0,
                  width: 2,
                  background: '#F59E0B',
                  opacity: 0.5
                }} />
                <div style={{
                  position: 'absolute',
                  left: '50%',
                  top: 0,
                  bottom: 0,
                  width: 2,
                  background: '#10b981',
                  opacity: 0.5
                }} />
              </div>
              <div className="row" style={{ justifyContent: 'space-between', marginTop: 6, fontSize: 11 }}>
                <span className="muted">0%</span>
                <span style={{ color: '#EF4444' }}>120%</span>
                <span style={{ color: '#10b981' }}>150%</span>
                <span className="muted">300%+</span>
              </div>
            </div>

            <div style={{
              padding: 12,
              background: 'rgba(255,255,255,.05)',
              borderRadius: 8,
              marginBottom: 12
            }}>
              {position.collateralRatio >= 150 ? (
                <div className="row" style={{ gap: 8 }}>
                  <CheckCircle2 size={18} color="#10b981" />
                  <div style={{ fontSize: 14 }}>
                    <div style={{ fontWeight: 600, color: '#10b981' }}>仓位安全</div>
                    <div className="muted" style={{ fontSize: 12 }}>抵押率高于最低要求</div>
                  </div>
                </div>
              ) : position.collateralRatio >= 120 ? (
                <div className="row" style={{ gap: 8 }}>
                  <AlertCircle size={18} color="#F59E0B" />
                  <div style={{ fontSize: 14 }}>
                    <div style={{ fontWeight: 600, color: '#F59E0B' }}>需要补充抵押</div>
                    <div className="muted" style={{ fontSize: 12 }}>接近清算线，请补充抵押物</div>
                  </div>
                </div>
              ) : (
                <div className="row" style={{ gap: 8 }}>
                  <AlertCircle size={18} color="#EF4444" />
                  <div style={{ fontSize: 14 }}>
                    <div style={{ fontWeight: 600, color: '#EF4444' }}>风险警告</div>
                    <div className="muted" style={{ fontSize: 12 }}>抵押率过低，可能被清算</div>
                  </div>
                </div>
              )}
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: 10, fontSize: 14 }}>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <span className="muted">抵押物:</span>
                <span style={{ fontWeight: 600 }}>{position.collateral.toLocaleString()} LP</span>
              </div>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <span className="muted">债务:</span>
                <span style={{ fontWeight: 600 }}>{position.debt.toLocaleString()} LUSD</span>
              </div>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <span className="muted">最大可铸造:</span>
                <span style={{ fontWeight: 600, color: '#10b981' }}>{position.maxMintable.toLocaleString()} LUSD</span>
              </div>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <span className="muted">清算价格:</span>
                <span style={{ fontWeight: 600, color: '#EF4444' }}>
                  ${position.liquidationPrice > 0 ? position.liquidationPrice.toFixed(4) : 'N/A'}
                </span>
              </div>
            </div>
          </div>

          {/* Protocol Info */}
          <div className="card" style={{ padding: 20, background: 'rgba(99, 102, 241, .1)', borderColor: '#6366F1' }}>
            <div style={{ fontWeight: 700, marginBottom: 12, display: 'flex', alignItems: 'center', gap: 8 }}>
              <Target size={20} color="#6366F1" />
              协议参数
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8, fontSize: 13 }}>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <span className="muted">最低抵押率:</span>
                <span>150%</span>
              </div>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <span className="muted">清算阈值:</span>
                <span style={{ color: '#F59E0B' }}>120%</span>
              </div>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <span className="muted">铸造费:</span>
                <span>0.2%</span>
              </div>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <span className="muted">赎回费:</span>
                <span>0.2%</span>
              </div>
              <div className="row" style={{ justifyContent: 'space-between' }}>
                <span className="muted">稳定费率:</span>
                <span>2% 年化</span>
              </div>
            </div>
          </div>

          {/* Activity */}
          <div className="card" style={{ padding: 20 }}>
            <div style={{ fontWeight: 700, marginBottom: 12, display: 'flex', alignItems: 'center', gap: 8 }}>
              <Activity size={20} color="#A855F7" />
              最近活动
            </div>
            <div style={{ fontSize: 13, color: '#666' }}>
              暂无交易记录
            </div>
          </div>
        </div>
      </div>

      {/* Info Banner */}
      <div className="card" style={{
        padding: 20,
        background: 'linear-gradient(135deg, rgba(16, 185, 129, .1), rgba(99, 102, 241, .1))',
        borderColor: '#10b981'
      }}>
        <div style={{ fontWeight: 700, marginBottom: 8, display: 'flex', alignItems: 'center', gap: 8 }}>
          <Info size={20} color="#10b981" />
          关于 LoyaltyUSD
        </div>
        <div className="muted" style={{ fontSize: 14, lineHeight: 1.6 }}>
          LoyaltyUSD (LUSD) 是与美元 1:1 挂钩的去中心化稳定币，由 Loyalty Points 超额抵押支持。
          用户可以通过质押 LP 代币铸造 LUSD，享受稳定的购买力和 DeFi 生态的流动性。
          系统采用自动化清算机制和价格预言机确保稳定性，所有操作完全透明且无需许可。
        </div>
      </div>
    </div>
  );
}
