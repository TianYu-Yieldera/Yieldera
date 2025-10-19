import React, { useState, useEffect } from "react";
import { useWallet } from "../web3/WalletContext";
import { useDemoMode } from "../web3/DemoModeContext";

export default function VaultView() {
  const { address } = useWallet();
  const { demoMode, demoData, demoDeposit, demoWithdraw, updateKey } = useDemoMode();

  // 模式切换：auto（智能）或 manual（手动）
  const [mode, setMode] = useState("auto");

  // Mock 数据
  const [stats, setStats] = useState({
    availablePoints: "0",
    lockedPoints: "0",
    totalEarnings: "0",
    currentAPY: "8.45"
  });

  const [depositAmount, setDepositAmount] = useState("");
  const [isDepositing, setIsDepositing] = useState(false);

  // 手动模式的质押记录
  const [manualStakes, setManualStakes] = useState([]);

  // 根据演示模式或钱包连接状态更新数据
  useEffect(() => {
    if (demoMode) {
      // 演示模式：使用 DemoModeContext 数据
      const lockedAmount = demoData.defiDeposits.reduce((sum, d) => sum + d.amount, 0);
      const totalEarned = demoData.defiDeposits.reduce((sum, d) => sum + (d.earned || 0), 0);

      setStats({
        availablePoints: demoData.points.toLocaleString(),
        lockedPoints: lockedAmount.toLocaleString(),
        totalEarnings: totalEarned.toFixed(2),
        currentAPY: "8.45"
      });

      // 手动模式：从 defiDeposits 生成质押记录
      const stakes = demoData.defiDeposits.map(d => ({
        protocol: d.protocol,
        amount: d.amount.toLocaleString(),
        apy: d.apy.toFixed(2),
        color: "#667eea"
      }));
      setManualStakes(stakes);
    } else if (address) {
      // 非演示模式：使用 Mock 数据
      setStats({
        availablePoints: "15,000",
        lockedPoints: "50,000",
        totalEarnings: "3,247.80",
        currentAPY: "8.45"
      });

      setManualStakes([
        { protocol: "Aave V3", amount: "20,000", apy: "3.52", color: "#9945FF" },
        { protocol: "Compound V3", amount: "15,000", apy: "4.18", color: "#00D395" }
      ]);
    }
  }, [address, demoMode, demoData, updateKey]);

  // 协议列表
  const protocols = [
    { id: "aave", name: "Aave V3", apy: "3.52", risk: "低", color: "#9945FF", desc: "最安全的借贷协议" },
    { id: "compound", name: "Compound V3", apy: "4.18", risk: "低", color: "#00D395", desc: "经典借贷协议" },
    { id: "uniswap", name: "Uniswap V3 LP", apy: "12.85", risk: "中", color: "#FF007A", desc: "流动性挖矿" },
    { id: "gmx", name: "GMX", apy: "22.30", risk: "高", color: "#3772FF", desc: "永续合约交易" }
  ];

  // Mock 策略分配数据（智能模式用）
  const strategies = [
    { protocol: "Aave V3", allocation: 40, apy: 3.52, amount: "20,000", color: "#9945FF" },
    { protocol: "Compound V3", allocation: 30, apy: 4.18, amount: "15,000", color: "#00D395" },
    { protocol: "Uniswap V3 LP", allocation: 20, apy: 18.75, amount: "10,000", color: "#FF007A" },
    { protocol: "GMX", allocation: 10, apy: 22.30, amount: "5,000", color: "#3772FF" }
  ];

  // Mock 收益历史
  const earningsHistory = [
    { date: "2025-10-18", amount: "+142.50", source: "Uniswap V3 LP" },
    { date: "2025-10-17", amount: "+89.20", source: "GMX" },
    { date: "2025-10-16", amount: "+67.30", source: "Aave V3" },
    { date: "2025-10-15", amount: "+95.80", source: "Compound V3" },
    { date: "2025-10-14", amount: "+123.40", source: "Uniswap V3 LP" }
  ];

  const handleAutoDeposit = () => {
    if (!address && !demoMode) {
      alert("请先连接钱包");
      return;
    }

    if (!depositAmount || parseFloat(depositAmount) <= 0) {
      alert("请输入有效的积分数量");
      return;
    }

    setIsDepositing(true);

    if (demoMode) {
      // 演示模式：调用 DemoModeContext
      const result = demoDeposit('Smart Vault', parseFloat(depositAmount));
      setIsDepositing(false);

      if (result.success) {
        setDepositAmount("");
        alert("✅ 存入成功！系统已自动分配到最优协议组合\n" + result.message);
      } else {
        alert("❌ " + result.error);
      }
    } else {
      // Mock: 模拟存款过程
      setTimeout(() => {
        const newLocked = parseInt(stats.lockedPoints.replace(/,/g, "")) + parseFloat(depositAmount);
        const newAvailable = parseInt(stats.availablePoints.replace(/,/g, "")) - parseFloat(depositAmount);

        setStats({
          ...stats,
          lockedPoints: newLocked.toLocaleString(),
          availablePoints: Math.max(0, newAvailable).toLocaleString()
        });

        setDepositAmount("");
        setIsDepositing(false);
        alert("✅ 存入成功！系统已自动分配到最优协议组合");
      }, 1500);
    }
  };

  const handleManualStake = (protocol) => {
    if (!address && !demoMode) {
      alert("请先连接钱包");
      return;
    }

    const amount = prompt(`质押到 ${protocol.name}\n\n当前 APY: ${protocol.apy}%\n风险等级: ${protocol.risk}\n\n请输入质押金额（积分）:`);

    if (!amount || parseFloat(amount) <= 0) {
      return;
    }

    if (demoMode) {
      // 演示模式：调用 DemoModeContext
      const result = demoDeposit(protocol.name, parseFloat(amount));

      if (result.success) {
        alert(`✅ 成功质押 ${amount} 积分到 ${protocol.name}\n` + result.message);
      } else {
        alert("❌ " + result.error);
      }
    } else {
      // 非演示模式：使用本地 Mock
      const available = parseInt(stats.availablePoints.replace(/,/g, ""));
      if (parseFloat(amount) > available) {
        alert("可用积分不足！");
        return;
      }

      setTimeout(() => {
        const newAvailable = available - parseFloat(amount);
        const newLocked = parseInt(stats.lockedPoints.replace(/,/g, "")) + parseFloat(amount);

        setStats({
          ...stats,
          availablePoints: newAvailable.toLocaleString(),
          lockedPoints: newLocked.toLocaleString()
        });

        const existing = manualStakes.find(s => s.protocol === protocol.name);
        if (existing) {
          setManualStakes(manualStakes.map(s =>
            s.protocol === protocol.name
              ? { ...s, amount: (parseFloat(s.amount.replace(/,/g, "")) + parseFloat(amount)).toLocaleString() }
              : s
          ));
        } else {
          setManualStakes([...manualStakes, {
            protocol: protocol.name,
            amount: parseFloat(amount).toLocaleString(),
            apy: protocol.apy,
            color: protocol.color
          }]);
        }

        alert(`✅ 成功质押 ${amount} 积分到 ${protocol.name}`);
      }, 500);
    }
  };

  const handleUnstake = (stake) => {
    if (confirm(`确定要从 ${stake.protocol} 取消质押吗？\n\n金额: ${stake.amount} 积分`)) {
      if (demoMode) {
        // 演示模式：找到对应的 deposit 并调用 demoWithdraw
        const depositIndex = demoData.defiDeposits.findIndex(d => d.protocol === stake.protocol);
        if (depositIndex >= 0) {
          const result = demoWithdraw(depositIndex);
          if (result.success) {
            alert(`✅ 已取消质押\n` + result.message);
          } else {
            alert("❌ " + result.error);
          }
        } else {
          alert("❌ 未找到对应的质押记录");
        }
      } else {
        // Mock: 取消质押
        const amount = parseFloat(stake.amount.replace(/,/g, ""));
        const newAvailable = parseInt(stats.availablePoints.replace(/,/g, "")) + amount;
        const newLocked = parseInt(stats.lockedPoints.replace(/,/g, "")) - amount;

        setStats({
          ...stats,
          availablePoints: newAvailable.toLocaleString(),
          lockedPoints: newLocked.toLocaleString()
        });

        setManualStakes(manualStakes.filter(s => s.protocol !== stake.protocol));

        alert(`✅ 已取消质押，${stake.amount} 积分已返回账户`);
      }
    }
  };

  const handleWithdraw = () => {
    alert("提现功能开发中... \n\n将支持：\n• 随时提现本金和收益\n• T+1 到账\n• 0.5% 手续费");
  };

  if (!address && !demoMode) {
    return (
      <div style={{ padding: "60px 20px", textAlign: "center", maxWidth: "600px", margin: "0 auto" }}>
        <div style={{ fontSize: "64px", marginBottom: "20px" }}>🔒</div>
        <h2 style={{ color: "#fff", marginBottom: "15px" }}>请先连接钱包或启用演示模式</h2>
        <p style={{ color: "#9ca3af", lineHeight: "1.6" }}>
          连接钱包后，您可以将积分存入理财金库，选择智能模式自动优化收益，或手动模式自主控制投资。
        </p>
        <button
          onClick={() => window.location.href = "/tutorial"}
          style={{
            marginTop: "20px",
            padding: "12px 30px",
            background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
            border: "none",
            borderRadius: "8px",
            color: "#fff",
            fontSize: "16px",
            fontWeight: "bold",
            cursor: "pointer"
          }}
        >
          前往教程启用演示模式
        </button>
      </div>
    );
  }

  return (
    <div style={{ padding: "40px 20px", maxWidth: "1200px", margin: "0 auto" }}>
      {/* 页面标题 */}
      <div style={{ marginBottom: "30px" }}>
        <h1 style={{ color: "#fff", fontSize: "32px", marginBottom: "10px", display: "flex", alignItems: "center", gap: "12px" }}>
          💰 理财金库
          {demoMode && (
            <span style={{
              fontSize: "14px",
              background: "#10b981",
              color: "white",
              padding: "4px 12px",
              borderRadius: "6px",
              fontWeight: "600"
            }}>
              演示模式
            </span>
          )}
        </h1>
        <p style={{ color: "#9ca3af", fontSize: "16px" }}>
          智能模式：系统自动优化收益 | 手动模式：自主选择协议投资
        </p>
      </div>

      {/* 统计卡片 */}
      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(240px, 1fr))", gap: "20px", marginBottom: "30px" }}>
        <div style={{ background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)", padding: "25px", borderRadius: "12px", boxShadow: "0 4px 6px rgba(0,0,0,0.3)" }}>
          <p style={{ color: "rgba(255,255,255,0.8)", margin: 0, fontSize: "14px" }}>可用积分</p>
          <h2 style={{ color: "#fff", margin: "10px 0 0 0", fontSize: "28px", fontWeight: "bold" }}>
            {stats.availablePoints}
          </h2>
        </div>

        <div style={{ background: "linear-gradient(135deg, #f093fb 0%, #f5576c 100%)", padding: "25px", borderRadius: "12px", boxShadow: "0 4px 6px rgba(0,0,0,0.3)" }}>
          <p style={{ color: "rgba(255,255,255,0.8)", margin: 0, fontSize: "14px" }}>已锁定积分</p>
          <h2 style={{ color: "#fff", margin: "10px 0 0 0", fontSize: "28px", fontWeight: "bold" }}>
            {stats.lockedPoints}
          </h2>
        </div>

        <div style={{ background: "linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)", padding: "25px", borderRadius: "12px", boxShadow: "0 4px 6px rgba(0,0,0,0.3)" }}>
          <p style={{ color: "rgba(255,255,255,0.8)", margin: 0, fontSize: "14px" }}>累计收益 (USDC)</p>
          <h2 style={{ color: "#fff", margin: "10px 0 0 0", fontSize: "28px", fontWeight: "bold" }}>
            ${stats.totalEarnings}
          </h2>
        </div>

        <div style={{ background: "linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)", padding: "25px", borderRadius: "12px", boxShadow: "0 4px 6px rgba(0,0,0,0.3)" }}>
          <p style={{ color: "rgba(255,255,255,0.8)", margin: 0, fontSize: "14px" }}>当前年化收益</p>
          <h2 style={{ color: "#fff", margin: "10px 0 0 0", fontSize: "28px", fontWeight: "bold" }}>
            {stats.currentAPY}%
          </h2>
        </div>
      </div>

      {/* 模式切换 */}
      <div style={{ background: "#1f2937", padding: "30px", borderRadius: "12px", border: "1px solid #374151", marginBottom: "30px" }}>
        <h3 style={{ color: "#fff", marginTop: 0, fontSize: "20px", marginBottom: "20px" }}>
          选择投资模式
        </h3>

        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "20px" }}>
          {/* 智能模式 */}
          <button
            onClick={() => setMode("auto")}
            style={{
              padding: "30px",
              background: mode === "auto"
                ? "linear-gradient(135deg, rgba(102,126,234,.2), rgba(118,75,162,.2))"
                : "#374151",
              border: mode === "auto" ? "2px solid #667eea" : "2px solid transparent",
              borderRadius: "12px",
              cursor: "pointer",
              textAlign: "left",
              transition: "all 0.3s"
            }}
          >
            <div style={{ fontSize: "32px", marginBottom: "15px" }}>🤖</div>
            <div style={{ color: "#fff", fontSize: "20px", fontWeight: "bold", marginBottom: "8px" }}>
              智能模式
              {mode === "auto" && <span style={{
                marginLeft: "10px",
                padding: "2px 8px",
                background: "#667eea",
                borderRadius: "4px",
                fontSize: "12px"
              }}>当前</span>}
            </div>
            <div style={{ color: "#9ca3af", fontSize: "14px", lineHeight: "1.5" }}>
              系统自动分配到收益最高的协议组合，每 24 小时智能再平衡。适合新手用户，无需操作。
            </div>
          </button>

          {/* 手动模式 */}
          <button
            onClick={() => setMode("manual")}
            style={{
              padding: "30px",
              background: mode === "manual"
                ? "linear-gradient(135deg, rgba(102,126,234,.2), rgba(118,75,162,.2))"
                : "#374151",
              border: mode === "manual" ? "2px solid #667eea" : "2px solid transparent",
              borderRadius: "12px",
              cursor: "pointer",
              textAlign: "left",
              transition: "all 0.3s"
            }}
          >
            <div style={{ fontSize: "32px", marginBottom: "15px" }}>🎯</div>
            <div style={{ color: "#fff", fontSize: "20px", fontWeight: "bold", marginBottom: "8px" }}>
              手动模式
              {mode === "manual" && <span style={{
                marginLeft: "10px",
                padding: "2px 8px",
                background: "#667eea",
                borderRadius: "4px",
                fontSize: "12px"
              }}>当前</span>}
              <span style={{
                marginLeft: "10px",
                padding: "2px 8px",
                background: "#f59e0b",
                borderRadius: "4px",
                fontSize: "12px"
              }}>高级</span>
            </div>
            <div style={{ color: "#9ca3af", fontSize: "14px", lineHeight: "1.5" }}>
              自己选择协议和投资比例，完全掌控资金分配。适合有经验的用户，追求个性化策略。
            </div>
          </button>
        </div>
      </div>

      {/* 根据模式显示不同内容 */}
      {mode === "auto" ? (
        <>
          {/* 智能模式：自动投资 */}
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "30px", marginBottom: "30px" }}>
            {/* 存入 */}
            <div style={{ background: "#1f2937", padding: "30px", borderRadius: "12px", border: "1px solid #374151" }}>
              <h3 style={{ color: "#fff", marginTop: 0, fontSize: "20px" }}>存入积分</h3>
              <p style={{ color: "#9ca3af", fontSize: "14px", marginBottom: "20px" }}>
                存入后，系统将自动分配到下方的最优协议组合
              </p>

              <div style={{ marginBottom: "15px" }}>
                <label style={{ color: "#9ca3af", fontSize: "14px", display: "block", marginBottom: "8px" }}>
                  存入数量
                </label>
                <input
                  type="number"
                  placeholder="输入积分数量"
                  value={depositAmount}
                  onChange={(e) => setDepositAmount(e.target.value)}
                  style={{
                    width: "100%",
                    padding: "12px",
                    background: "#374151",
                    border: "1px solid #4b5563",
                    borderRadius: "8px",
                    color: "#fff",
                    fontSize: "16px"
                  }}
                />
                <div style={{ marginTop: "8px", display: "flex", gap: "8px" }}>
                  {["25%", "50%", "75%", "100%"].map(pct => (
                    <button
                      key={pct}
                      onClick={() => {
                        const available = parseInt(stats.availablePoints.replace(/,/g, ""));
                        const amount = Math.floor(available * parseInt(pct) / 100);
                        setDepositAmount(amount.toString());
                      }}
                      style={{
                        flex: 1,
                        padding: "6px",
                        background: "#4b5563",
                        border: "none",
                        borderRadius: "6px",
                        color: "#fff",
                        fontSize: "12px",
                        cursor: "pointer"
                      }}
                    >
                      {pct}
                    </button>
                  ))}
                </div>
              </div>

              <button
                onClick={handleAutoDeposit}
                disabled={isDepositing}
                style={{
                  width: "100%",
                  padding: "14px",
                  background: isDepositing ? "#4b5563" : "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
                  border: "none",
                  borderRadius: "8px",
                  color: "#fff",
                  fontSize: "16px",
                  fontWeight: "bold",
                  cursor: isDepositing ? "not-allowed" : "pointer",
                  transition: "all 0.3s"
                }}
              >
                {isDepositing ? "处理中..." : "智能存入"}
              </button>
            </div>

            {/* 提现 */}
            <div style={{ background: "#1f2937", padding: "30px", borderRadius: "12px", border: "1px solid #374151" }}>
              <h3 style={{ color: "#fff", marginTop: 0, fontSize: "20px" }}>提现</h3>
              <p style={{ color: "#9ca3af", fontSize: "14px", marginBottom: "20px" }}>
                随时提现您的本金和收益，T+1 到账
              </p>

              <div style={{ background: "#374151", padding: "20px", borderRadius: "8px", marginBottom: "20px" }}>
                <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "12px" }}>
                  <span style={{ color: "#9ca3af", fontSize: "14px" }}>可提现本金</span>
                  <span style={{ color: "#fff", fontWeight: "bold" }}>{stats.lockedPoints} 积分</span>
                </div>
                <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "12px" }}>
                  <span style={{ color: "#9ca3af", fontSize: "14px" }}>可提现收益</span>
                  <span style={{ color: "#10b981", fontWeight: "bold" }}>${stats.totalEarnings}</span>
                </div>
                <div style={{ borderTop: "1px solid #4b5563", marginTop: "12px", paddingTop: "12px", display: "flex", justifyContent: "space-between" }}>
                  <span style={{ color: "#9ca3af", fontSize: "14px" }}>总计</span>
                  <span style={{ color: "#fff", fontSize: "18px", fontWeight: "bold" }}>
                    {stats.lockedPoints} 积分 + ${stats.totalEarnings}
                  </span>
                </div>
              </div>

              <button
                onClick={handleWithdraw}
                style={{
                  width: "100%",
                  padding: "14px",
                  background: "linear-gradient(135deg, #f093fb 0%, #f5576c 100%)",
                  border: "none",
                  borderRadius: "8px",
                  color: "#fff",
                  fontSize: "16px",
                  fontWeight: "bold",
                  cursor: "pointer"
                }}
              >
                申请提现
              </button>

              <p style={{ color: "#9ca3af", fontSize: "12px", marginTop: "12px", marginBottom: 0 }}>
                ⚠️ 提现需要 24 小时处理，手续费 0.5%
              </p>
            </div>
          </div>

          {/* 当前投资策略 */}
          <div style={{ background: "#1f2937", padding: "30px", borderRadius: "12px", border: "1px solid #374151", marginBottom: "30px" }}>
            <h3 style={{ color: "#fff", marginTop: 0, fontSize: "20px", marginBottom: "20px" }}>
              📊 系统自动分配策略
            </h3>

            <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(250px, 1fr))", gap: "15px" }}>
              {strategies.map((strategy, idx) => (
                <div
                  key={idx}
                  style={{
                    background: "#374151",
                    padding: "20px",
                    borderRadius: "8px",
                    border: `2px solid ${strategy.color}`
                  }}
                >
                  <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "12px" }}>
                    <h4 style={{ color: "#fff", margin: 0, fontSize: "16px" }}>{strategy.protocol}</h4>
                    <span style={{ color: strategy.color, fontWeight: "bold", fontSize: "14px" }}>
                      {strategy.allocation}%
                    </span>
                  </div>

                  <div style={{ marginBottom: "8px" }}>
                    <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "4px" }}>
                      <span style={{ color: "#9ca3af", fontSize: "13px" }}>投资金额</span>
                      <span style={{ color: "#fff", fontSize: "13px" }}>{strategy.amount} USDC</span>
                    </div>
                    <div style={{ display: "flex", justifyContent: "space-between" }}>
                      <span style={{ color: "#9ca3af", fontSize: "13px" }}>当前 APY</span>
                      <span style={{ color: "#10b981", fontSize: "13px", fontWeight: "bold" }}>
                        {strategy.apy}%
                      </span>
                    </div>
                  </div>

                  <div style={{ background: "#1f2937", borderRadius: "4px", height: "6px", overflow: "hidden" }}>
                    <div
                      style={{
                        background: strategy.color,
                        width: `${strategy.allocation}%`,
                        height: "100%"
                      }}
                    />
                  </div>
                </div>
              ))}
            </div>

            <div style={{ marginTop: "20px", padding: "15px", background: "#374151", borderRadius: "8px" }}>
              <p style={{ color: "#9ca3af", fontSize: "14px", margin: 0 }}>
                💡 <strong style={{ color: "#fff" }}>智能再平衡：</strong>
                系统每 24 小时自动调整策略，确保您获得最优收益。无需手动操作，完全自动化。
              </p>
            </div>
          </div>
        </>
      ) : (
        <>
          {/* 手动模式：自选协议 */}
          <div style={{ marginBottom: "30px" }}>
            <div style={{ background: "#1f2937", padding: "30px", borderRadius: "12px", border: "1px solid #374151" }}>
              <h3 style={{ color: "#fff", marginTop: 0, fontSize: "20px", marginBottom: "20px" }}>
                🎯 选择协议进行质押
              </h3>

              <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(280px, 1fr))", gap: "20px" }}>
                {protocols.map((protocol) => (
                  <div
                    key={protocol.id}
                    style={{
                      background: "#374151",
                      padding: "24px",
                      borderRadius: "12px",
                      border: "2px solid #4b5563",
                      transition: "all 0.3s"
                    }}
                  >
                    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", marginBottom: "15px" }}>
                      <div>
                        <h4 style={{ color: "#fff", margin: "0 0 8px 0", fontSize: "18px" }}>
                          {protocol.name}
                        </h4>
                        <p style={{ color: "#9ca3af", fontSize: "13px", margin: 0 }}>
                          {protocol.desc}
                        </p>
                      </div>
                      <div style={{
                        padding: "4px 10px",
                        background: protocol.risk === "低" ? "#10b981" : protocol.risk === "中" ? "#f59e0b" : "#ef4444",
                        borderRadius: "6px",
                        fontSize: "12px",
                        fontWeight: "bold",
                        color: "#fff"
                      }}>
                        {protocol.risk}风险
                      </div>
                    </div>

                    <div style={{
                      background: "#1f2937",
                      padding: "15px",
                      borderRadius: "8px",
                      marginBottom: "15px"
                    }}>
                      <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "8px" }}>
                        <span style={{ color: "#9ca3af", fontSize: "14px" }}>当前 APY</span>
                        <span style={{ color: "#10b981", fontSize: "18px", fontWeight: "bold" }}>
                          {protocol.apy}%
                        </span>
                      </div>

                      {/* 显示已质押金额（如果有） */}
                      {manualStakes.find(s => s.protocol === protocol.name) && (
                        <div style={{ display: "flex", justifyContent: "space-between", paddingTop: "8px", borderTop: "1px solid #374151" }}>
                          <span style={{ color: "#9ca3af", fontSize: "14px" }}>已质押</span>
                          <span style={{ color: "#fff", fontSize: "14px", fontWeight: "bold" }}>
                            {manualStakes.find(s => s.protocol === protocol.name).amount}
                          </span>
                        </div>
                      )}
                    </div>

                    <button
                      onClick={() => handleManualStake(protocol)}
                      style={{
                        width: "100%",
                        padding: "12px",
                        background: `linear-gradient(135deg, ${protocol.color}dd 0%, ${protocol.color} 100%)`,
                        border: "none",
                        borderRadius: "8px",
                        color: "#fff",
                        fontSize: "14px",
                        fontWeight: "bold",
                        cursor: "pointer"
                      }}
                    >
                      质押到此协议
                    </button>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* 我的质押记录 */}
          {manualStakes.length > 0 && (
            <div style={{ background: "#1f2937", padding: "30px", borderRadius: "12px", border: "1px solid #374151", marginBottom: "30px" }}>
              <h3 style={{ color: "#fff", marginTop: 0, fontSize: "20px", marginBottom: "20px" }}>
                📋 我的质押记录
              </h3>

              <div style={{ display: "grid", gap: "15px" }}>
                {manualStakes.map((stake, idx) => (
                  <div
                    key={idx}
                    style={{
                      background: "#374151",
                      padding: "20px",
                      borderRadius: "8px",
                      display: "flex",
                      justifyContent: "space-between",
                      alignItems: "center"
                    }}
                  >
                    <div style={{ flex: 1 }}>
                      <h4 style={{ color: "#fff", margin: "0 0 8px 0", fontSize: "16px" }}>
                        {stake.protocol}
                      </h4>
                      <div style={{ display: "flex", gap: "20px" }}>
                        <div>
                          <span style={{ color: "#9ca3af", fontSize: "13px" }}>质押金额: </span>
                          <span style={{ color: "#fff", fontSize: "13px", fontWeight: "bold" }}>{stake.amount}</span>
                        </div>
                        <div>
                          <span style={{ color: "#9ca3af", fontSize: "13px" }}>APY: </span>
                          <span style={{ color: "#10b981", fontSize: "13px", fontWeight: "bold" }}>{stake.apy}%</span>
                        </div>
                      </div>
                    </div>

                    <button
                      onClick={() => handleUnstake(stake)}
                      style={{
                        padding: "8px 20px",
                        background: "transparent",
                        border: "2px solid #ef4444",
                        borderRadius: "6px",
                        color: "#ef4444",
                        fontSize: "14px",
                        fontWeight: "bold",
                        cursor: "pointer"
                      }}
                    >
                      取消质押
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}
        </>
      )}

      {/* 收益历史（两种模式共享） */}
      <div style={{ background: "#1f2937", padding: "30px", borderRadius: "12px", border: "1px solid #374151" }}>
        <h3 style={{ color: "#fff", marginTop: 0, fontSize: "20px", marginBottom: "20px" }}>
          📈 收益历史
        </h3>

        <div style={{ overflowX: "auto" }}>
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr style={{ borderBottom: "2px solid #374151" }}>
                <th style={{ textAlign: "left", padding: "12px", color: "#9ca3af", fontSize: "14px", fontWeight: "500" }}>
                  日期
                </th>
                <th style={{ textAlign: "left", padding: "12px", color: "#9ca3af", fontSize: "14px", fontWeight: "500" }}>
                  来源
                </th>
                <th style={{ textAlign: "right", padding: "12px", color: "#9ca3af", fontSize: "14px", fontWeight: "500" }}>
                  收益 (USDC)
                </th>
              </tr>
            </thead>
            <tbody>
              {earningsHistory.map((item, idx) => (
                <tr key={idx} style={{ borderBottom: "1px solid #374151" }}>
                  <td style={{ padding: "12px", color: "#fff", fontSize: "14px" }}>{item.date}</td>
                  <td style={{ padding: "12px", color: "#9ca3af", fontSize: "14px" }}>{item.source}</td>
                  <td style={{ padding: "12px", color: "#10b981", fontSize: "14px", fontWeight: "bold", textAlign: "right" }}>
                    {item.amount}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* CTA: 去 RWA 商城 */}
      <div
        style={{
          marginTop: "30px",
          padding: "30px",
          background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
          borderRadius: "12px",
          textAlign: "center"
        }}
      >
        <h3 style={{ color: "#fff", marginTop: 0, fontSize: "24px", marginBottom: "10px" }}>
          💎 用收益购买真实资产
        </h3>
        <p style={{ color: "rgba(255,255,255,0.9)", marginBottom: "20px" }}>
          您已赚取 ${stats.totalEarnings}，可以用来购买股票、黄金、美债等真实资产
        </p>
        <button
          onClick={() => window.location.href = "/rwa-market"}
          style={{
            padding: "14px 40px",
            background: "#fff",
            border: "none",
            borderRadius: "8px",
            color: "#667eea",
            fontSize: "16px",
            fontWeight: "bold",
            cursor: "pointer"
          }}
        >
          前往 RWA 商城 →
        </button>
      </div>
    </div>
  );
}
