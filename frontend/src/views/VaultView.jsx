import React, { useState, useEffect } from "react";
import { useWallet } from "../web3/WalletContext";
import { useDemoMode } from "../web3/DemoModeContext";
import { Vault, TrendingUp, DollarSign, Shield, Zap, Activity, ArrowDownCircle, ArrowUpCircle } from 'lucide-react';
import TechContainer from "../components/ui/TechContainer";
import TechHeader from "../components/ui/TechHeader";
import TechCard from "../components/ui/TechCard";
import TechButton from "../components/ui/TechButton";

export default function VaultView() {
  const { address } = useWallet();
  const { demoMode, demoData, demoDeposit, demoWithdraw, updateKey } = useDemoMode();

  const [stats, setStats] = useState({
    availablePoints: 0,
    lockedPoints: 0,
    totalEarnings: 0,
    currentAPY: 8.45
  });

  const [depositAmount, setDepositAmount] = useState("");
  const [isDepositing, setIsDepositing] = useState(false);
  const [activePositions, setActivePositions] = useState([]);

  useEffect(() => {
    if (demoMode) {
      const lockedAmount = demoData.defiDeposits.reduce((sum, d) => sum + d.amount, 0);
      const totalEarned = demoData.defiDeposits.reduce((sum, d) => sum + (d.earned || 0), 0);

      setStats({
        availablePoints: demoData.points,
        lockedPoints: lockedAmount,
        totalEarnings: totalEarned,
        currentAPY: 8.45
      });

      setActivePositions(demoData.defiDeposits.map(d => ({
        protocol: d.protocol,
        amount: d.amount,
        apy: d.apy,
        earned: d.earned || 0,
        index: demoData.defiDeposits.indexOf(d)
      })));
    } else {
      setStats({
        availablePoints: 15000,
        lockedPoints: 50000,
        totalEarnings: 3247.80,
        currentAPY: 8.45
      });

      setActivePositions([
        { protocol: "Aave V3", amount: 20000, apy: 3.52, earned: 1250.50 },
        { protocol: "Compound V3", amount: 15000, apy: 4.18, earned: 890.20 },
        { protocol: "Uniswap V3 LP", amount: 10000, apy: 12.85, earned: 850.40 },
        { protocol: "GMX", amount: 5000, apy: 22.30, earned: 256.70 }
      ]);
    }
  }, [address, demoMode, demoData, updateKey]);

  const handleDeposit = () => {
    if (!depositAmount || parseFloat(depositAmount) <= 0) {
      alert("请输入有效的存入金额");
      return;
    }

    setIsDepositing(true);

    if (demoMode) {
      const result = demoDeposit('Smart Vault', parseFloat(depositAmount));
      setIsDepositing(false);

      if (result.success) {
        setDepositAmount("");
        alert("✅ 存入成功！系统已自动分配到最优协议组合\n" + result.message);
      } else {
        alert("❌ " + result.error);
      }
    } else {
      setTimeout(() => {
        setStats(prev => ({
          ...prev,
          lockedPoints: prev.lockedPoints + parseFloat(depositAmount),
          availablePoints: Math.max(0, prev.availablePoints - parseFloat(depositAmount))
        }));
        setDepositAmount("");
        setIsDepositing(false);
        alert("✅ 存入成功！");
      }, 1500);
    }
  };

  const handleWithdraw = (position) => {
    if (confirm(`确定要取出 ${position.amount.toLocaleString()} Points 吗？`)) {
      if (demoMode && position.index !== undefined) {
        const result = demoWithdraw(position.index);
        if (result.success) {
          alert(`✅ 取出成功\n` + result.message);
        } else {
          alert("❌ " + result.error);
        }
      } else {
        setStats(prev => ({
          ...prev,
          lockedPoints: prev.lockedPoints - position.amount,
          availablePoints: prev.availablePoints + position.amount + position.earned
        }));
        alert("✅ 取出成功！");
      }
    }
  };

  const protocols = [
    { name: "Aave V3", apy: 3.52, risk: "Low", color: "rgb(34, 211, 238)" },
    { name: "Compound V3", apy: 4.18, risk: "Low", color: "rgb(34, 197, 94)" },
    { name: "Uniswap V3 LP", apy: 12.85, risk: "Medium", color: "rgb(251, 191, 36)" },
    { name: "GMX", apy: 22.30, risk: "High", color: "rgb(239, 68, 68)" }
  ];

  return (
    <TechContainer>
      <TechHeader
        icon={Vault}
        title="Smart Vault"
        subtitle="AI-powered yield optimization across DeFi protocols"
      />

      {/* Stats Grid */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: 16, marginBottom: 32 }}>
        <TechCard
          icon={DollarSign}
          title="Available Balance"
          value={stats.availablePoints.toLocaleString()}
          subtitle="Points ready to deploy"
          iconColor="rgb(34, 211, 238)"
        />
        <TechCard
          icon={Shield}
          title="Total Deposited"
          value={stats.lockedPoints.toLocaleString()}
          subtitle={`Across ${activePositions.length} positions`}
          iconColor="rgb(59, 130, 246)"
        />
        <TechCard
          icon={TrendingUp}
          title="Total Earnings"
          value={stats.totalEarnings.toLocaleString()}
          subtitle="Lifetime earnings"
          iconColor="rgb(34, 197, 94)"
        />
        <TechCard
          icon={Zap}
          title="Current APY"
          value={`${stats.currentAPY}%`}
          subtitle="Weighted average yield"
          iconColor="rgb(251, 191, 36)"
        />
      </div>

      {/* Deposit Section */}
      <TechCard>
        <h2 style={{ fontSize: 18, fontWeight: 600, color: 'white', margin: '0 0 20px 0' }}>
          Deposit to Smart Vault
        </h2>
        <div style={{ display: 'flex', gap: 12, alignItems: 'flex-end', flexWrap: 'wrap' }}>
          <div style={{ flex: 1, minWidth: 250 }}>
            <label style={{ display: 'block', fontSize: 14, fontWeight: 500, color: 'rgba(203, 213, 225, 0.8)', marginBottom: 8 }}>
              Amount (Points)
            </label>
            <input
              type="number"
              value={depositAmount}
              onChange={(e) => setDepositAmount(e.target.value)}
              placeholder="Enter amount to deposit"
              style={{
                width: '100%',
                padding: '12px 16px',
                border: '1px solid rgba(34, 211, 238, 0.3)',
                borderRadius: 8,
                fontSize: 16,
                fontWeight: 500,
                color: 'white',
                background: 'rgba(15, 23, 42, 0.5)',
                outline: 'none'
              }}
              onFocus={(e) => e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.6)'}
              onBlur={(e) => e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.3)'}
            />
          </div>
          <TechButton
            onClick={handleDeposit}
            disabled={isDepositing}
            loading={isDepositing}
            icon={ArrowDownCircle}
            variant="primary"
          >
            Deposit
          </TechButton>
        </div>
        <p style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.7)', marginTop: 16, display: 'flex', alignItems: 'center', gap: 8 }}>
          <Activity style={{ width: 16, height: 16 }} />
          Your funds will be automatically allocated across optimal DeFi protocols to maximize yield
        </p>
      </TechCard>

      {/* Active Positions */}
      <div style={{ marginTop: 24 }}>
        <TechCard>
          <h2 style={{ fontSize: 18, fontWeight: 600, color: 'white', margin: '0 0 20px 0' }}>
            Active Positions
          </h2>
          {activePositions.length === 0 ? (
            <div style={{
              padding: 48,
              textAlign: 'center',
              background: 'rgba(15, 23, 42, 0.5)',
              borderRadius: 8,
              border: '1px dashed rgba(34, 211, 238, 0.3)'
            }}>
              <Activity style={{ width: 40, height: 40, color: 'rgba(203, 213, 225, 0.5)', margin: '0 auto 12px' }} />
              <p style={{ color: 'rgba(203, 213, 225, 0.7)', fontSize: 14 }}>
                No active positions. Deposit funds to start earning yield.
              </p>
            </div>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
              {activePositions.map((position, index) => (
                <div
                  key={index}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    padding: 20,
                    background: 'rgba(15, 23, 42, 0.5)',
                    borderRadius: 8,
                    border: '1px solid rgba(34, 211, 238, 0.2)',
                    transition: 'all 0.3s ease',
                    flexWrap: 'wrap',
                    gap: 16
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.4)';
                    e.currentTarget.style.background = 'rgba(15, 23, 42, 0.7)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.2)';
                    e.currentTarget.style.background = 'rgba(15, 23, 42, 0.5)';
                  }}
                >
                  <div style={{ flex: 1, minWidth: 200 }}>
                    <div style={{ fontSize: 16, fontWeight: 600, color: 'white', marginBottom: 6 }}>
                      {position.protocol}
                    </div>
                    <div style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.7)', display: 'flex', alignItems: 'center', gap: 12, flexWrap: 'wrap' }}>
                      <span style={{ color: 'rgb(34, 211, 238)', fontWeight: 600 }}>
                        APY: {position.apy.toFixed(2)}%
                      </span>
                      <span>•</span>
                      <span style={{ color: 'rgb(34, 197, 94)', fontWeight: 600 }}>
                        Earned: +{position.earned.toLocaleString()} Points
                      </span>
                    </div>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
                    <div style={{ textAlign: 'right' }}>
                      <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', textTransform: 'uppercase', letterSpacing: 0.5 }}>
                        Deposited
                      </div>
                      <div style={{ fontSize: 20, fontWeight: 700, color: 'white' }}>
                        {position.amount.toLocaleString()}
                      </div>
                    </div>
                    <TechButton
                      onClick={() => handleWithdraw(position)}
                      variant="danger"
                      icon={ArrowUpCircle}
                    >
                      Withdraw
                    </TechButton>
                  </div>
                </div>
              ))}
            </div>
          )}
        </TechCard>
      </div>

      {/* Available Protocols */}
      <div style={{ marginTop: 24 }}>
        <TechCard>
          <h2 style={{ fontSize: 18, fontWeight: 600, color: 'white', margin: '0 0 20px 0' }}>
            Supported Protocols
          </h2>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: 16 }}>
            {protocols.map((protocol, index) => (
              <div
                key={index}
                style={{
                  padding: 20,
                  background: 'rgba(15, 23, 42, 0.5)',
                  borderRadius: 8,
                  border: '1px solid rgba(34, 211, 238, 0.2)',
                  transition: 'all 0.3s ease'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.borderColor = protocol.color;
                  e.currentTarget.style.boxShadow = `0 4px 20px ${protocol.color}40`;
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.2)';
                  e.currentTarget.style.boxShadow = 'none';
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
                  <div style={{ fontSize: 16, fontWeight: 600, color: 'white' }}>
                    {protocol.name}
                  </div>
                  <div style={{
                    padding: '4px 12px',
                    borderRadius: 12,
                    fontSize: 11,
                    fontWeight: 600,
                    textTransform: 'uppercase',
                    background: protocol.risk === 'Low'
                      ? 'rgba(34, 197, 94, 0.2)'
                      : protocol.risk === 'Medium'
                      ? 'rgba(251, 191, 36, 0.2)'
                      : 'rgba(239, 68, 68, 0.2)',
                    color: protocol.risk === 'Low'
                      ? 'rgb(34, 197, 94)'
                      : protocol.risk === 'Medium'
                      ? 'rgb(251, 191, 36)'
                      : 'rgb(239, 68, 68)',
                    border: protocol.risk === 'Low'
                      ? '1px solid rgba(34, 197, 94, 0.4)'
                      : protocol.risk === 'Medium'
                      ? '1px solid rgba(251, 191, 36, 0.4)'
                      : '1px solid rgba(239, 68, 68, 0.4)'
                  }}>
                    {protocol.risk}
                  </div>
                </div>
                <div style={{
                  fontSize: 24,
                  fontWeight: 700,
                  color: protocol.color,
                  marginBottom: 4,
                  textShadow: `0 0 20px ${protocol.color}60`
                }}>
                  {protocol.apy}% APY
                </div>
                <div style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.6)' }}>
                  Annual Percentage Yield
                </div>
              </div>
            ))}
          </div>
        </TechCard>
      </div>
    </TechContainer>
  );
}
