import React, { useState, useEffect } from "react";
import { Wallet, TrendingUp, Zap, DollarSign, Lock, Unlock, Info } from "lucide-react";
import { useWallet } from "../web3/WalletContext";
import { useDemoMode } from "../web3/DemoModeContext";
import { config } from "../config/env";

const API_BASE = config.api.baseUrl;

export default function DeFiPoolView() {
  const { address } = useWallet();
  const { demoMode, demoAddress, demoData, demoDeposit, demoWithdraw, demoStake, demoUnstake, demoMintStablecoin, demoRedeemStablecoin, updateKey } = useDemoMode();
  const [userPoints, setUserPoints] = useState(0);
  const [userTokens, setUserTokens] = useState(0);
  const [selectedPool, setSelectedPool] = useState(null);
  const [amounts, setAmounts] = useState({}); // 为每个池子维护独立的输入值

  useEffect(() => {
    if (demoMode) {
      // 演示模式：使用模拟数据
      setUserPoints(demoData.points);
      setUserTokens(demoData.pfiTokens);
    } else if (address) {
      // 真实模式：从API获取
      Promise.all([
        fetch(`${API_BASE}/users/${address}/points`).then(r => r.json()),
        fetch(`${API_BASE}/users/${address}/balance`).then(r => r.json())
      ]).then(([points, balance]) => {
        setUserPoints(parseFloat(points?.points || 0));
        setUserTokens(parseFloat(balance?.balance || 0));
      }).catch(console.error);
    }
  }, [address, demoMode, demoData, updateKey]); // 添加 updateKey 到依赖

  // 根据演示模式数据查找用户存款
  const getUserDeposit = (protocol) => {
    if (!demoMode) return null;
    return demoData.defiDeposits.find(d => d.protocol === protocol);
  };

  // DeFi 协议池配置
  const poolsConfig = [
    {
      id: "uniswap",
      protocol: "Uniswap V3",
      name: "积分兑换 USDC",
      description: "用积分兑换 USDC，兑换的资金提供流动性赚取手续费",
      icon: "🦄",
      color: "#FF007A",
      type: "积分兑换",
      apr: 45.2,
      tradingFeeApr: "12.8%",
      pointsRewardApr: "32.4%",
      tvl: "$2,450,000",
      pointsRequired: 100,
      exchangeRate: 0.01, // 1 积分 = 0.01 USDC
      pointsMultiplier: 1.5,
      risk: "中等",
      canDeposit: true,
      canWithdraw: true,
      features: ["积分兑换 USDC", "自动提供流动性", "赚取交易手续费", "自动复投收益"]
    },
    {
      id: "aave",
      protocol: "Aave V3",
      name: "积分兑换 aUSDC",
      description: "用积分兑换 Aave 生息代币 aUSDC，自动赚取存款利息",
      icon: "👻",
      color: "#B6509E",
      type: "积分兑换",
      apr: 8.5,
      baseApr: "5.2%",
      pointsBoostApr: "3.3%",
      tvl: "$8,920,000",
      pointsRequired: 100,
      exchangeRate: 0.01, // 1 积分 = 0.01 USDC
      pointsMultiplier: 1.2,
      risk: "低",
      canDeposit: true,
      canWithdraw: true,
      features: ["积分兑换 aUSDC", "自动赚取利息", "即时赎回积分", "本金保护"]
    },
    {
      id: "stablecoin",
      protocol: "LoyaltyUSD",
      name: "LUSD 稳定币铸造",
      description: "按 100 积分 = 1 LUSD 的比例铸造稳定币（需 150% 抵押率）",
      icon: "💵",
      color: "#10b981",
      type: "稳定币铸造",
      apr: 0,
      collateralRatio: "150%",
      liquidationThreshold: "120%",
      mintRatio: 100, // 100 积分 = 1 LUSD
      tvl: "$18,450,000",
      pointsRequired: 0,
      pointsMultiplier: 1.0,
      risk: "低",
      canDeposit: true,
      canWithdraw: true,
      features: ["1:1 美元挂钩", "150% 超额抵押", "按比例铸造", "随时赎回抵押物"]
    },
    {
      id: "staking",
      protocol: "LoyaltyX Protocol",
      name: "PFI 质押池",
      description: "质押 PFI 代币赚取积分奖励，积分可提升 DeFi 收益",
      icon: "🔒",
      color: "#6366F1",
      type: "代币质押",
      apr: 125,
      stakingApr: "125%",
      pointsRewardApr: "125%",
      tvl: "$5,680,000",
      pointsRequired: 0,
      pointsMultiplier: 1.0,
      risk: "低",
      canDeposit: true,
      canWithdraw: true,
      features: ["高额积分奖励", "质押赚积分", "随时解除质押", "无锁定期"],
      usesTokens: true  // 标记这个池子使用代币而非积分
    }
  ];

  // 为每个池子添加用户数据
  const defiPools = poolsConfig.map(pool => {
    // 质押池特殊处理
    if (pool.id === "staking" && demoMode) {
      return {
        ...pool,
        userDeposited: demoData.stakedTokens.toString(),
        userEarned: "0 PFI", // 质押池不产生代币收益
        userPointsEarned: demoData.stakingRewards > 0 ? `${demoData.stakingRewards.toFixed(2)} Points` : "0 Points",
        apr: pool.apr > 0 ? `${pool.apr}%` : "N/A"
      };
    }

    // 稳定币池特殊处理
    if (pool.id === "stablecoin" && demoMode) {
      return {
        ...pool,
        userDeposited: demoData.collateral.toString(), // 抵押的积分
        userEarned: demoData.stablecoinMinted > 0 ? `${demoData.stablecoinMinted.toFixed(2)} LUSD` : "0 LUSD", // 铸造的稳定币
        userPointsEarned: demoData.collateral > 0 ? `${((demoData.collateral / demoData.stablecoinMinted || 0) * 100).toFixed(0)}%` : "0%", // 抵押率
        apr: pool.apr > 0 ? `${pool.apr}%` : "N/A"
      };
    }

    // 其他池子使用存款数据（积分兑换）
    const userDeposit = getUserDeposit(pool.protocol);
    const depositAmount = userDeposit ? userDeposit.amount : 0;
    const earnedAmount = userDeposit ? userDeposit.earned : 0;
    const exchangedAmount = userDeposit ? userDeposit.exchangedAmount : 0;

    return {
      ...pool,
      userDeposited: depositAmount.toString(),
      userEarned: exchangedAmount > 0 ? `${(exchangedAmount + earnedAmount).toFixed(2)} ${userDeposit.asset}` : "0 USDC",
      userPointsEarned: earnedAmount > 0 ? `${pool.apr}%` : `${pool.apr}%`, // 显示收益率
      apr: pool.apr > 0 ? `${pool.apr}%` : "N/A"
    };
  });

  // 用户总资产统计
  const totalDeposited = demoMode
    ? demoData.defiDeposits.reduce((sum, d) => sum + d.amount, 0)
    : 0;
  const averageApr = poolsConfig.reduce((sum, pool) => sum + pool.apr, 0) / poolsConfig.length;

  const handleDeposit = (pool) => {
    const addr = demoMode ? demoAddress : address;
    if (!addr) {
      alert("请先连接钱包或启用演示模式");
      return;
    }

    const amount = amounts[pool.id] || "";
    if (!amount || parseFloat(amount) <= 0) {
      alert("请输入有效金额");
      return;
    }

    if (demoMode) {
      // 质押池特殊处理：质押代币
      if (pool.id === "staking") {
        if (parseFloat(amount) > userTokens) {
          alert("PFI 代币余额不足");
          return;
        }
        const result = demoStake(parseFloat(amount));
        if (result.success) {
          alert(`✅ ${result.message}\n\n您的 PFI 代币已被质押，开始赚取积分奖励！`);
          setAmounts({...amounts, [pool.id]: ""});
        } else {
          alert(`❌ ${result.error}`);
        }
        return;
      }

      // 稳定币池特殊处理：铸造稳定币
      if (pool.id === "stablecoin") {
        const result = demoMintStablecoin(parseFloat(amount));
        if (result.success) {
          alert(`✅ ${result.message}`);
          setAmounts({...amounts, [pool.id]: ""});
        } else {
          alert(`❌ ${result.error}`);
        }
        return;
      }

      // 其他池子：需要积分才能参与
      if (userPoints < pool.pointsRequired) {
        alert(`需要至少 ${pool.pointsRequired} 积分才能参与此池`);
        return;
      }

      // 演示模式：使用模拟操作兑换积分
      const result = demoDeposit(pool.protocol, parseFloat(amount));
      if (result.success) {
        alert(`✅ ${result.message}`);
        setAmounts({...amounts, [pool.id]: ""});
      } else {
        alert(`❌ ${result.error}`);
      }
    } else {
      // 真实模式：调用智能合约（待开发）
      const action = pool.id === "staking" ? "质押" : pool.id === "stablecoin" ? "铸造" : "兑换";
      alert(`${action} ${amount} 到 ${pool.name}\n\n智能合约功能即将推出...`);
      setAmounts({...amounts, [pool.id]: ""});
    }
  };

  const handleWithdraw = (pool) => {
    const addr = demoMode ? demoAddress : address;
    if (!addr) {
      alert("请先连接钱包或启用演示模式");
      return;
    }

    if (demoMode) {
      // 质押池特殊处理：解除质押
      if (pool.id === "staking") {
        if (demoData.stakedTokens <= 0) {
          alert("您没有质押的代币");
          return;
        }
        // 取出全部质押
        const result = demoUnstake(demoData.stakedTokens);
        if (result.success) {
          alert(`✅ ${result.message}`);
        } else {
          alert(`❌ ${result.error}`);
        }
        return;
      }

      // 稳定币池特殊处理：赎回抵押物
      if (pool.id === "stablecoin") {
        if (demoData.stablecoinMinted <= 0) {
          alert("您没有铸造的稳定币");
          return;
        }
        // 赎回全部稳定币
        const result = demoRedeemStablecoin(demoData.stablecoinMinted);
        if (result.success) {
          alert(`✅ ${result.message}`);
        } else {
          alert(`❌ ${result.error}`);
        }
        return;
      }

      // 其他池子：查找并取出兑换的资产
      const depositIndex = demoData.defiDeposits.findIndex(d => d.protocol === pool.protocol);
      if (depositIndex === -1) {
        alert("您在此协议中没有兑换记录");
        return;
      }

      const result = demoWithdraw(depositIndex);
      if (result.success) {
        alert(`✅ ${result.message}`);
      } else {
        alert(`❌ ${result.error}`);
      }
    } else {
      // 真实模式：调用智能合约（待开发）
      const action = pool.id === "staking" ? "解除质押" : pool.id === "stablecoin" ? "赎回" : "取回";
      alert(`${action}从 ${pool.name}\n\n智能合约功能即将推出...`);
    }
  };

  const handleClaim = (pool) => {
    const addr = demoMode ? demoAddress : address;
    if (!addr) {
      alert("请先连接钱包或启用演示模式");
      return;
    }

    if (demoMode) {
      alert(`✅ 已领取 ${pool.name} 的收益！\n\n（演示模式：收益已自动复投）`);
    } else {
      alert(`领取 ${pool.name} 的收益\n\n智能合约功能即将推出...`);
    }
  };

  return (
    <div className="container">
      {/* 标题 */}
      <div style={{ marginBottom: 24 }}>
        <h1 style={{ margin: 0, fontSize: 32, display: 'flex', alignItems: 'center', gap: 12 }}>
          <Zap size={36} color="#6366F1" />
          DeFi 协议池
        </h1>
        <p className="muted" style={{ marginTop: 8 }}>
          用积分参与 DeFi 协议，获得额外收益加成
        </p>
      </div>

      {/* 用户资产概览 */}
      <div className="grid" style={{ gridTemplateColumns: 'repeat(4, 1fr)', gap: 16, marginBottom: 24 }}>
        <div className="kpi">
          <div className="title">PFI 代币</div>
          <div className="value" style={{ color: '#10b981' }}>
            {userTokens.toLocaleString('en-US', { maximumFractionDigits: 2 })}
          </div>
          <div className="muted" style={{ marginTop: 4 }}>可交易代币</div>
        </div>
        <div className="kpi">
          <div className="title">我的积分</div>
          <div className="value" style={{ color: '#6366F1' }}>
            {userPoints.toLocaleString('en-US', { maximumFractionDigits: 2 })}
          </div>
          <div className="muted" style={{ marginTop: 4 }}>权益凭证</div>
        </div>
        <div className="kpi">
          <div className="title">总存入</div>
          <div className="value" style={{ color: '#A855F7' }}>
            ${totalDeposited.toLocaleString()}
          </div>
          <div className="muted" style={{ marginTop: 4 }}>所有协议</div>
        </div>
        <div className="kpi">
          <div className="title">平均 APR</div>
          <div className="value" style={{ color: '#F59E0B' }}>
            {averageApr.toFixed(1)}%
          </div>
          <div className="muted" style={{ marginTop: 4 }}>年化收益率</div>
        </div>
      </div>

      {/* 积分权益说明 */}
      <div className="card" style={{
        padding: 20,
        marginBottom: 24,
        background: 'linear-gradient(135deg, rgba(99, 102, 241, .1), rgba(168, 85, 247, .1))',
        borderColor: '#6366F1'
      }}>
        <div className="row" style={{ gap: 12, marginBottom: 12 }}>
          <Info size={20} color="#6366F1" />
          <div style={{ fontWeight: 700 }}>LoyaltyX 积分兑换系统</div>
        </div>
        <div className="muted" style={{ fontSize: 14, lineHeight: 1.6 }}>
          • <strong>积分兑换资产：</strong>用积分兑换 USDC、aUSDC 等 DeFi 资产<br/>
          • <strong>自动赚取收益：</strong>兑换的资产自动参与 DeFi 协议赚取收益<br/>
          • <strong>随时赎回：</strong>可随时赎回兑换的资产，获取本金和收益<br/>
          • <strong>稳定币铸造：</strong>按 100 积分 = 1 LUSD 比例铸造稳定币
        </div>
      </div>

      {/* DeFi 池列表 */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
        {defiPools.map((pool) => (
          <div
            key={pool.id}
            className="card"
            style={{
              padding: 0,
              overflow: 'hidden',
              borderColor: selectedPool?.id === pool.id ? pool.color : 'rgba(255,255,255,.1)',
              transition: 'all 0.3s'
            }}
          >
            {/* 池子头部 */}
            <div style={{
              padding: 24,
              background: `linear-gradient(135deg, ${pool.color}15, ${pool.color}05)`,
              borderBottom: '1px solid rgba(255,255,255,.05)'
            }}>
              <div className="row" style={{ justifyContent: 'space-between', alignItems: 'flex-start' }}>
                {/* 左侧：协议信息 */}
                <div className="row" style={{ gap: 16, flex: 1 }}>
                  <div style={{
                    fontSize: 48,
                    width: 80,
                    height: 80,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    background: `${pool.color}22`,
                    borderRadius: 16
                  }}>
                    {pool.icon}
                  </div>

                  <div style={{ flex: 1 }}>
                    <div className="row" style={{ gap: 12, marginBottom: 4 }}>
                      <div style={{ fontWeight: 700, fontSize: 20 }}>{pool.name}</div>
                      <div style={{
                        background: `${pool.color}22`,
                        color: pool.color,
                        padding: '4px 10px',
                        borderRadius: 6,
                        fontSize: 12,
                        fontWeight: 700
                      }}>
                        {pool.protocol}
                      </div>
                    </div>
                    <div className="muted" style={{ fontSize: 14, marginBottom: 12 }}>
                      {pool.description}
                    </div>
                    <div className="row" style={{ gap: 12, flexWrap: 'wrap' }}>
                      <div style={{
                        padding: '6px 12px',
                        background: 'rgba(255,255,255,.05)',
                        borderRadius: 8,
                        fontSize: 13
                      }}>
                        <span className="muted">类型: </span>
                        <span>{pool.type}</span>
                      </div>
                      <div style={{
                        padding: '6px 12px',
                        background: 'rgba(255,255,255,.05)',
                        borderRadius: 8,
                        fontSize: 13
                      }}>
                        <span className="muted">风险: </span>
                        <span style={{
                          color: pool.risk === "低" ? "#10b981" : pool.risk === "中等" ? "#F59E0B" : "#EF4444"
                        }}>{pool.risk}</span>
                      </div>
                      {pool.pointsRequired > 0 && (
                        <div style={{
                          padding: '6px 12px',
                          background: `${pool.color}22`,
                          borderRadius: 8,
                          fontSize: 13,
                          color: pool.color,
                          fontWeight: 600
                        }}>
                          需要 {pool.pointsRequired} 积分
                        </div>
                      )}
                    </div>
                  </div>
                </div>

                {/* 右侧：APR 展示 */}
                <div style={{
                  background: `${pool.color}22`,
                  padding: 20,
                  borderRadius: 16,
                  minWidth: 200,
                  textAlign: 'center'
                }}>
                  <div className="muted" style={{ fontSize: 12, marginBottom: 4 }}>年化收益率</div>
                  <div style={{ fontSize: 36, fontWeight: 900, color: pool.color }}>
                    {pool.apr}
                  </div>
                  <div className="muted" style={{ fontSize: 11, marginTop: 8 }}>
                    TVL: {pool.tvl}
                  </div>
                </div>
              </div>

              {/* 功能特性 */}
              <div className="row" style={{ gap: 12, marginTop: 16, flexWrap: 'wrap' }}>
                {pool.features.map((feature, idx) => (
                  <div key={idx} style={{
                    padding: '6px 12px',
                    background: 'rgba(255,255,255,.05)',
                    borderRadius: 8,
                    fontSize: 12,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 6
                  }}>
                    <TrendingUp size={14} color={pool.color} />
                    {feature}
                  </div>
                ))}
              </div>
            </div>

            {/* 用户数据和操作 */}
            <div style={{ padding: 24 }}>
              <div className="grid" style={{ gridTemplateColumns: 'repeat(3, 1fr) 2fr', gap: 16 }}>
                {/* 用户数据 */}
                <div>
                  <div className="muted" style={{ fontSize: 12, marginBottom: 4 }}>
                    {pool.id === "staking" ? "质押的代币" : pool.id === "stablecoin" ? "积分抵押" : "兑换积分"}
                  </div>
                  <div style={{ fontWeight: 700, fontSize: 16 }}>
                    {pool.userDeposited} {pool.id === "staking" ? "PFI" : "Points"}
                  </div>
                </div>
                <div>
                  <div className="muted" style={{ fontSize: 12, marginBottom: 4 }}>
                    {pool.id === "stablecoin" ? "已铸造" : "已赚取"}
                  </div>
                  <div style={{ fontWeight: 700, fontSize: 16, color: '#10b981' }}>{pool.userEarned}</div>
                </div>
                <div>
                  <div className="muted" style={{ fontSize: 12, marginBottom: 4 }}>
                    {pool.id === "staking" ? "积分奖励" : pool.id === "stablecoin" ? "抵押率" : "收益率"}
                  </div>
                  <div style={{ fontWeight: 700, fontSize: 16, color: pool.color }}>{pool.userPointsEarned}</div>
                </div>

                {/* 操作区域 */}
                <div className="row" style={{ gap: 8 }}>
                  <input
                    type="number"
                    placeholder="输入金额"
                    value={amounts[pool.id] || ""}
                    onChange={(e) => setAmounts({...amounts, [pool.id]: e.target.value})}
                    style={{
                      flex: 1,
                      padding: '10px 12px',
                      background: 'rgba(0,0,0,.3)',
                      border: '1px solid rgba(255,255,255,.1)',
                      borderRadius: 8,
                      color: '#fff',
                      fontSize: 14
                    }}
                  />
                  <button
                    className="btn"
                    style={{ background: pool.color, minWidth: 80 }}
                    onClick={() => handleDeposit(pool)}
                  >
                    <Lock size={16} style={{ marginRight: 4 }} />
                    {pool.id === "staking" ? "质押" : pool.id === "stablecoin" ? "铸造" : "兑换"}
                  </button>
                  <button
                    className="btn"
                    style={{ background: 'rgba(255,255,255,.1)', minWidth: 80 }}
                    onClick={() => handleWithdraw(pool)}
                  >
                    <Unlock size={16} style={{ marginRight: 4 }} />
                    {pool.id === "staking" ? "解押" : pool.id === "stablecoin" ? "赎回" : "取回"}
                  </button>
                  <button
                    className="btn"
                    style={{ background: '#10b981', minWidth: 80 }}
                    onClick={() => handleClaim(pool)}
                  >
                    领取
                  </button>
                </div>
              </div>

              {/* 收益明细 */}
              {pool.id === "uniswap" && (
                <div style={{
                  marginTop: 16,
                  padding: 16,
                  background: 'rgba(255,255,255,.02)',
                  borderRadius: 12
                }}>
                  <div style={{ fontWeight: 600, marginBottom: 12, fontSize: 14 }}>收益构成</div>
                  <div className="row" style={{ gap: 24, fontSize: 13 }}>
                    <div>
                      <span className="muted">交易手续费: </span>
                      <span style={{ color: '#10b981' }}>{pool.tradingFeeApr}</span>
                    </div>
                    <div>
                      <span className="muted">积分奖励: </span>
                      <span style={{ color: pool.color }}>{pool.pointsRewardApr}</span>
                    </div>
                    <div>
                      <span className="muted">积分加成: </span>
                      <span style={{ color: '#F59E0B' }}>{pool.pointsMultiplier}x</span>
                    </div>
                  </div>
                </div>
              )}

              {pool.id === "aave" && (
                <div style={{
                  marginTop: 16,
                  padding: 16,
                  background: 'rgba(255,255,255,.02)',
                  borderRadius: 12
                }}>
                  <div style={{ fontWeight: 600, marginBottom: 12, fontSize: 14 }}>收益构成</div>
                  <div className="row" style={{ gap: 24, fontSize: 13 }}>
                    <div>
                      <span className="muted">基础利率: </span>
                      <span style={{ color: '#10b981' }}>{pool.baseApr}</span>
                    </div>
                    <div>
                      <span className="muted">积分加速: </span>
                      <span style={{ color: pool.color }}>{pool.pointsBoostApr}</span>
                    </div>
                    <div>
                      <span className="muted">加成倍数: </span>
                      <span style={{ color: '#F59E0B' }}>{pool.pointsMultiplier}x</span>
                    </div>
                  </div>
                </div>
              )}

              {pool.id === "stablecoin" && (
                <div style={{
                  marginTop: 16,
                  padding: 16,
                  background: 'rgba(255,255,255,.02)',
                  borderRadius: 12
                }}>
                  <div style={{ fontWeight: 600, marginBottom: 12, fontSize: 14 }}>协议参数</div>
                  <div className="row" style={{ gap: 24, fontSize: 13 }}>
                    <div>
                      <span className="muted">最低抵押率: </span>
                      <span style={{ color: '#10b981' }}>{pool.collateralRatio}</span>
                    </div>
                    <div>
                      <span className="muted">清算阈值: </span>
                      <span style={{ color: '#F59E0B' }}>{pool.liquidationThreshold}</span>
                    </div>
                    <div>
                      <span className="muted">LUSD 价格: </span>
                      <span style={{ color: pool.color }}>$1.00</span>
                    </div>
                  </div>
                </div>
              )}

              {pool.id === "staking" && (
                <div style={{
                  marginTop: 16,
                  padding: 16,
                  background: 'rgba(255,255,255,.02)',
                  borderRadius: 12
                }}>
                  <div style={{ fontWeight: 600, marginBottom: 12, fontSize: 14 }}>质押详情</div>
                  <div className="row" style={{ gap: 24, fontSize: 13 }}>
                    <div>
                      <span className="muted">质押奖励率: </span>
                      <span style={{ color: pool.color }}>{pool.stakingApr}</span>
                    </div>
                    <div>
                      <span className="muted">积分产出: </span>
                      <span style={{ color: '#10b981' }}>持续累积</span>
                    </div>
                    <div>
                      <span className="muted">锁定期: </span>
                      <span style={{ color: '#F59E0B' }}>无锁定</span>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>

      {/* 未连接钱包提示 */}
      {!address && (
        <div className="card" style={{
          marginTop: 24,
          padding: 24,
          textAlign: 'center',
          background: 'rgba(245, 158, 11, .1)',
          borderColor: '#F59E0B'
        }}>
          <Wallet size={48} color="#F59E0B" style={{ margin: '0 auto 16px' }} />
          <div style={{ fontWeight: 700, marginBottom: 8 }}>请先连接钱包</div>
          <div className="muted">连接钱包后即可参与 DeFi 协议池</div>
        </div>
      )}
    </div>
  );
}
