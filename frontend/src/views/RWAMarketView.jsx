import React, { useState, useEffect } from "react";
import { useWallet } from "../web3/WalletContext";
import { useDemoMode } from "../web3/DemoModeContext";

export default function RWAMarketView() {
  const { address } = useWallet();
  const { demoMode, demoData, updateKey } = useDemoMode();

  // Mock 用户余额
  const [userBalance, setUserBalance] = useState({
    earnings: "3,247.80", // USDC 收益
    holdings: []
  });

  const [selectedTab, setSelectedTab] = useState("stocks"); // stocks, bonds, commodities

  // Mock RWA 资产数据
  const assets = {
    stocks: [
      {
        ticker: "bAAPL",
        name: "Apple Inc.",
        price: "178.52",
        change: "+2.34%",
        logo: "🍎",
        description: "科技巨头，iPhone 制造商",
        protocol: "Backed Finance"
      },
      {
        ticker: "bTSLA",
        name: "Tesla Inc.",
        price: "242.18",
        change: "+5.67%",
        logo: "⚡",
        description: "电动汽车领导者",
        protocol: "Backed Finance"
      },
      {
        ticker: "bNVDA",
        name: "NVIDIA Corp.",
        price: "485.93",
        change: "+8.92%",
        logo: "🎮",
        description: "AI 芯片龙头",
        protocol: "Backed Finance"
      },
      {
        ticker: "bGOOGL",
        name: "Google",
        price: "139.47",
        change: "+1.23%",
        logo: "🔍",
        description: "搜索引擎巨头",
        protocol: "Backed Finance"
      },
      {
        ticker: "bAMZN",
        name: "Amazon",
        price: "178.25",
        change: "+3.45%",
        logo: "📦",
        description: "电商及云计算巨头",
        protocol: "Backed Finance"
      },
      {
        ticker: "bMSFT",
        name: "Microsoft",
        price: "398.76",
        change: "+2.11%",
        logo: "💻",
        description: "软件及云服务巨头",
        protocol: "Backed Finance"
      }
    ],
    bonds: [
      {
        ticker: "OUSG",
        name: "Ondo Short-Term US Treasuries",
        price: "105.20",
        change: "+0.12%",
        logo: "🏛️",
        description: "短期美国国债代币",
        protocol: "Ondo Finance",
        apy: "4.5%"
      },
      {
        ticker: "USDY",
        name: "Ondo US Dollar Yield",
        price: "100.85",
        change: "+0.08%",
        logo: "💵",
        description: "美元收益代币",
        protocol: "Ondo Finance",
        apy: "5.2%"
      }
    ],
    commodities: [
      {
        ticker: "PAXG",
        name: "Paxos Gold",
        price: "2,042.50",
        change: "+1.05%",
        logo: "🥇",
        description: "1 PAXG = 1 盎司黄金",
        protocol: "Paxos"
      },
      {
        ticker: "XAUT",
        name: "Tether Gold",
        price: "2,041.80",
        change: "+1.02%",
        logo: "🪙",
        description: "1 XAUT = 1 盎司黄金",
        protocol: "Tether"
      }
    ]
  };

  // 根据演示模式或钱包连接状态更新余额
  useEffect(() => {
    if (demoMode) {
      // 演示模式：从 defiDeposits 计算总收益
      const totalEarned = demoData.defiDeposits.reduce((sum, d) => sum + (d.earned || 0), 0);

      // 从 localStorage 读取 RWA 持仓和已花费金额
      const savedHoldings = JSON.parse(localStorage.getItem('demo_rwa_holdings') || '[]');
      const spentAmount = parseFloat(localStorage.getItem('demo_rwa_spent') || '0');

      // 可用收益 = 总收益 - 已花费
      const availableEarnings = Math.max(0, totalEarned - spentAmount);

      setUserBalance({
        earnings: availableEarnings.toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ","),
        holdings: savedHoldings
      });
    } else if (address) {
      // 非演示模式：使用 Mock 数据
      setUserBalance({
        earnings: "3,247.80",
        holdings: [
          { ticker: "bAAPL", amount: "5.2", value: "928.30" },
          { ticker: "PAXG", amount: "0.5", value: "1,021.25" }
        ]
      });
    }
  }, [address, demoMode, demoData, updateKey]);

  const handleBuy = (asset) => {
    if (!address && !demoMode) {
      alert("请先连接钱包或启用演示模式");
      return;
    }

    const amountToBuy = prompt(`购买 ${asset.name}\n\n当前价格: $${asset.price}\n您的可用余额: $${userBalance.earnings}\n\n请输入购买金额（USDC）:`);

    if (!amountToBuy || parseFloat(amountToBuy) <= 0) {
      return;
    }

    const userBalanceNum = parseFloat(userBalance.earnings.replace(/,/g, ""));
    const buyAmount = parseFloat(amountToBuy);

    if (buyAmount > userBalanceNum) {
      alert("余额不足！请先在理财金库赚取更多收益。");
      return;
    }

    setTimeout(() => {
      if (demoMode) {
        // 演示模式：更新 localStorage
        const currentSpent = parseFloat(localStorage.getItem('demo_rwa_spent') || '0');
        localStorage.setItem('demo_rwa_spent', (currentSpent + buyAmount).toString());

        // 更新持仓
        const savedHoldings = JSON.parse(localStorage.getItem('demo_rwa_holdings') || '[]');
        const assetPrice = parseFloat(asset.price.replace(/,/g, ""));
        const quantity = buyAmount / assetPrice;

        const existingIndex = savedHoldings.findIndex(h => h.ticker === asset.ticker);
        if (existingIndex >= 0) {
          savedHoldings[existingIndex] = {
            ticker: asset.ticker,
            amount: (parseFloat(savedHoldings[existingIndex].amount) + quantity).toFixed(4),
            value: (parseFloat(savedHoldings[existingIndex].value.replace(/,/g, "")) + buyAmount).toFixed(2)
          };
        } else {
          savedHoldings.push({
            ticker: asset.ticker,
            amount: quantity.toFixed(4),
            value: buyAmount.toFixed(2)
          });
        }

        localStorage.setItem('demo_rwa_holdings', JSON.stringify(savedHoldings));

        // 触发 UI 更新
        setUserBalance({
          earnings: (userBalanceNum - buyAmount).toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ","),
          holdings: savedHoldings
        });
      } else {
        // 非演示模式：本地状态更新
        setUserBalance({
          ...userBalance,
          earnings: (userBalanceNum - buyAmount).toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ",")
        });
      }

      alert(`✅ 购买成功！\n\n您购买了 $${buyAmount} 的 ${asset.ticker}\n交易将在 1-2 分钟内确认`);
    }, 500);
  };

  const currentAssets = assets[selectedTab] || [];

  if (!address && !demoMode) {
    return (
      <div style={{ padding: "60px 20px", textAlign: "center", maxWidth: "600px", margin: "0 auto" }}>
        <div style={{ fontSize: "64px", marginBottom: "20px" }}>🔒</div>
        <h2 style={{ color: "#fff", marginBottom: "15px" }}>请先连接钱包或启用演示模式</h2>
        <p style={{ color: "#9ca3af", lineHeight: "1.6" }}>
          连接钱包后，您可以使用理财收益购买真实世界资产（股票、债券、黄金等）
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
    <div style={{ padding: "40px 20px", maxWidth: "1400px", margin: "0 auto" }}>
      {/* 页面标题 */}
      <div style={{ marginBottom: "30px" }}>
        <h1 style={{ color: "#fff", fontSize: "32px", marginBottom: "10px", display: "flex", alignItems: "center", gap: "12px" }}>
          💎 RWA 资产市场
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
          用您的收益购买真实世界资产（Real World Assets）
        </p>
      </div>

      {/* 用户余额卡片 */}
      <div style={{ background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)", padding: "30px", borderRadius: "12px", marginBottom: "30px", boxShadow: "0 4px 6px rgba(0,0,0,0.3)" }}>
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", flexWrap: "wrap", gap: "20px" }}>
          <div>
            <p style={{ color: "rgba(255,255,255,0.8)", margin: 0, fontSize: "14px" }}>可用收益</p>
            <h2 style={{ color: "#fff", margin: "8px 0 0 0", fontSize: "36px", fontWeight: "bold" }}>
              ${userBalance.earnings}
            </h2>
          </div>

          <button
            onClick={() => window.location.href = "/vault"}
            style={{
              padding: "12px 30px",
              background: "rgba(255,255,255,0.2)",
              border: "2px solid #fff",
              borderRadius: "8px",
              color: "#fff",
              fontSize: "14px",
              fontWeight: "bold",
              cursor: "pointer",
              backdropFilter: "blur(10px)"
            }}
          >
            前往理财金库赚取更多 →
          </button>
        </div>

        {userBalance.holdings.length > 0 && (
          <div style={{ marginTop: "20px", paddingTop: "20px", borderTop: "1px solid rgba(255,255,255,0.2)" }}>
            <p style={{ color: "rgba(255,255,255,0.8)", margin: "0 0 10px 0", fontSize: "14px" }}>我的持仓</p>
            <div style={{ display: "flex", gap: "15px", flexWrap: "wrap" }}>
              {userBalance.holdings.map((holding, idx) => (
                <div
                  key={idx}
                  style={{
                    background: "rgba(255,255,255,0.15)",
                    padding: "12px 20px",
                    borderRadius: "8px",
                    backdropFilter: "blur(10px)"
                  }}
                >
                  <span style={{ color: "#fff", fontWeight: "bold" }}>{holding.ticker}</span>
                  <span style={{ color: "rgba(255,255,255,0.8)", margin: "0 10px" }}>×</span>
                  <span style={{ color: "#fff" }}>{holding.amount}</span>
                  <span style={{ color: "rgba(255,255,255,0.6)", marginLeft: "10px", fontSize: "12px" }}>
                    (${holding.value})
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* 分类标签 */}
      <div style={{ display: "flex", gap: "10px", marginBottom: "30px", borderBottom: "2px solid #374151" }}>
        {[
          { key: "stocks", label: "📈 股票", count: assets.stocks.length },
          { key: "bonds", label: "🏛️ 债券", count: assets.bonds.length },
          { key: "commodities", label: "🥇 大宗商品", count: assets.commodities.length }
        ].map((tab) => (
          <button
            key={tab.key}
            onClick={() => setSelectedTab(tab.key)}
            style={{
              padding: "15px 25px",
              background: "transparent",
              border: "none",
              borderBottom: selectedTab === tab.key ? "3px solid #667eea" : "3px solid transparent",
              color: selectedTab === tab.key ? "#fff" : "#9ca3af",
              fontSize: "16px",
              fontWeight: selectedTab === tab.key ? "bold" : "normal",
              cursor: "pointer",
              transition: "all 0.3s"
            }}
          >
            {tab.label} ({tab.count})
          </button>
        ))}
      </div>

      {/* 资产列表 */}
      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(320px, 1fr))", gap: "20px" }}>
        {currentAssets.map((asset, idx) => (
          <div
            key={idx}
            style={{
              background: "#1f2937",
              padding: "25px",
              borderRadius: "12px",
              border: "1px solid #374151",
              transition: "all 0.3s",
              cursor: "pointer"
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.borderColor = "#667eea";
              e.currentTarget.style.transform = "translateY(-5px)";
              e.currentTarget.style.boxShadow = "0 10px 20px rgba(0,0,0,0.3)";
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.borderColor = "#374151";
              e.currentTarget.style.transform = "translateY(0)";
              e.currentTarget.style.boxShadow = "none";
            }}
          >
            {/* 资产头部 */}
            <div style={{ display: "flex", alignItems: "center", gap: "15px", marginBottom: "15px" }}>
              <div style={{ fontSize: "48px" }}>{asset.logo}</div>
              <div style={{ flex: 1 }}>
                <h3 style={{ color: "#fff", margin: "0 0 5px 0", fontSize: "18px" }}>{asset.name}</h3>
                <p style={{ color: "#9ca3af", margin: 0, fontSize: "13px" }}>{asset.ticker}</p>
              </div>
            </div>

            {/* 价格 */}
            <div style={{ marginBottom: "15px" }}>
              <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                <span style={{ color: "#fff", fontSize: "28px", fontWeight: "bold" }}>
                  ${asset.price}
                </span>
                <span style={{
                  color: asset.change.startsWith("+") ? "#10b981" : "#ef4444",
                  fontSize: "16px",
                  fontWeight: "bold"
                }}>
                  {asset.change}
                </span>
              </div>
              {asset.apy && (
                <div style={{ marginTop: "8px" }}>
                  <span style={{ color: "#9ca3af", fontSize: "13px" }}>年化收益: </span>
                  <span style={{ color: "#10b981", fontSize: "13px", fontWeight: "bold" }}>{asset.apy}</span>
                </div>
              )}
            </div>

            {/* 描述 */}
            <p style={{ color: "#9ca3af", fontSize: "14px", lineHeight: "1.5", marginBottom: "15px" }}>
              {asset.description}
            </p>

            {/* 协议 */}
            <div style={{ marginBottom: "15px" }}>
              <span style={{
                display: "inline-block",
                padding: "4px 12px",
                background: "#374151",
                borderRadius: "6px",
                color: "#9ca3af",
                fontSize: "12px"
              }}>
                {asset.protocol}
              </span>
            </div>

            {/* 购买按钮 */}
            <button
              onClick={() => handleBuy(asset)}
              style={{
                width: "100%",
                padding: "12px",
                background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
                border: "none",
                borderRadius: "8px",
                color: "#fff",
                fontSize: "14px",
                fontWeight: "bold",
                cursor: "pointer"
              }}
            >
              立即购买
            </button>
          </div>
        ))}
      </div>

      {/* 说明区域 */}
      <div style={{ marginTop: "40px", background: "#1f2937", padding: "30px", borderRadius: "12px", border: "1px solid #374151" }}>
        <h3 style={{ color: "#fff", marginTop: 0, fontSize: "20px", marginBottom: "20px" }}>
          ℹ️ 关于 RWA（真实世界资产）
        </h3>

        <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(300px, 1fr))", gap: "20px" }}>
          <div>
            <h4 style={{ color: "#667eea", marginTop: 0, fontSize: "16px" }}>什么是 RWA？</h4>
            <p style={{ color: "#9ca3af", fontSize: "14px", lineHeight: "1.6" }}>
              RWA 是链上代币化的真实世界资产，如股票、债券、黄金等。持有 RWA 代币等同于持有底层资产的权益。
            </p>
          </div>

          <div>
            <h4 style={{ color: "#667eea", marginTop: 0, fontSize: "16px" }}>安全性如何？</h4>
            <p style={{ color: "#9ca3af", fontSize: "14px", lineHeight: "1.6" }}>
              所有 RWA 资产均由受监管的机构发行，底层资产由托管机构保管，并定期审计。如 Backed Finance、Ondo Finance 等。
            </p>
          </div>

          <div>
            <h4 style={{ color: "#667eea", marginTop: 0, fontSize: "16px" }}>如何交易？</h4>
            <p style={{ color: "#9ca3af", fontSize: "14px", lineHeight: "1.6" }}>
              购买后，RWA 代币存储在您的钱包中，可以随时在 Uniswap 等 DEX 上交易，或在本平台内卖出。
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
