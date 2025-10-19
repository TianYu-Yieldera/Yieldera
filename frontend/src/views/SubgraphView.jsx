import React, { useState } from "react";
import { Network, Search, TrendingUp, Activity, Database, Code } from "lucide-react";

export default function SubgraphView() {
  const [selectedChain, setSelectedChain] = useState("sepolia");
  const [query, setQuery] = useState("");

  // 模拟链数据
  const chains = [
    { id: "sepolia", name: "Sepolia", color: "#6366F1", status: "active" },
    { id: "mainnet", name: "Ethereum Mainnet", color: "#627EEA", status: "coming" },
    { id: "arbitrum", name: "Arbitrum", color: "#28A0F0", status: "coming" },
    { id: "optimism", name: "Optimism", color: "#FF0420", status: "coming" }
  ];

  // 模拟链上数据统计
  const chainStats = {
    sepolia: {
      totalTransactions: "12,458",
      totalUsers: "1,234",
      totalVolume: "1,458,920 PFI",
      avgGasPrice: "2.5 Gwei",
      lastBlock: "5,234,567"
    }
  };

  // 示例查询
  const exampleQueries = [
    {
      title: "获取前10名用户",
      description: "查询积分最高的10个地址",
      query: `{
  users(first: 10, orderBy: points, orderDirection: desc) {
    id
    address
    points
    balance
  }
}`
    },
    {
      title: "用户交易历史",
      description: "查询特定用户的所有交易",
      query: `{
  balanceEvents(where: { user: "0x..." }) {
    id
    amount
    eventType
    txHash
    blockNumber
  }
}`
    },
    {
      title: "质押数据",
      description: "查询所有质押事件",
      query: `{
  stakingEvents(first: 100) {
    id
    user
    amount
    timestamp
    type
  }
}`
    }
  ];

  // 最近活动
  const recentActivity = [
    { type: "Transfer", from: "0x1234...5678", to: "0xabcd...efgh", amount: "1,000 PFI", time: "2分钟前" },
    { type: "Stake", from: "0x9876...5432", to: "Staking Pool", amount: "5,000 PFI", time: "5分钟前" },
    { type: "Claim", from: "Reward Pool", to: "0x2468...1357", amount: "250 PFI", time: "8分钟前" },
    { type: "Transfer", from: "0x1111...2222", to: "0x3333...4444", amount: "500 PFI", time: "12分钟前" },
    { type: "Unstake", from: "Staking Pool", to: "0xaaaa...bbbb", amount: "3,000 PFI", time: "15分钟前" }
  ];

  const getTypeColor = (type) => {
    switch(type) {
      case "Transfer": return "#6366F1";
      case "Stake": return "#10b981";
      case "Unstake": return "#F59E0B";
      case "Claim": return "#A855F7";
      default: return "#666";
    }
  };

  return (
    <div className="container">
      <div style={{ marginBottom: 24 }}>
        <h1 style={{ margin: 0, fontSize: 32, display: 'flex', alignItems: 'center', gap: 12 }}>
          <Network size={36} color="#6366F1" />
          链上数据索引
        </h1>
        <p className="muted" style={{ marginTop: 8 }}>实时查询链上数据和事件</p>
      </div>

      {/* 链选择器 */}
      <div className="card" style={{ padding: 20, marginBottom: 24 }}>
        <div style={{ fontWeight: 700, marginBottom: 12 }}>选择网络</div>
        <div className="row" style={{ gap: 12, flexWrap: 'wrap' }}>
          {chains.map(chain => (
            <button
              key={chain.id}
              className="btn"
              style={{
                background: selectedChain === chain.id ? chain.color : 'rgba(255,255,255,.05)',
                borderColor: selectedChain === chain.id ? chain.color : 'rgba(255,255,255,.1)',
                opacity: chain.status === "coming" ? 0.5 : 1,
                cursor: chain.status === "coming" ? 'not-allowed' : 'pointer',
                padding: '12px 20px'
              }}
              onClick={() => chain.status === "active" && setSelectedChain(chain.id)}
              disabled={chain.status === "coming"}
            >
              {chain.name}
              {chain.status === "coming" && <span className="muted" style={{ marginLeft: 8, fontSize: 11 }}>(即将上线)</span>}
            </button>
          ))}
        </div>
      </div>

      {/* 链数据统计 */}
      <div className="grid" style={{ gridTemplateColumns: 'repeat(5, 1fr)', gap: 16, marginBottom: 24 }}>
        {Object.entries(chainStats[selectedChain] || {}).map(([key, value], index) => (
          <div key={index} className="kpi">
            <div className="title">{
              key === "totalTransactions" ? "总交易数" :
              key === "totalUsers" ? "总用户数" :
              key === "totalVolume" ? "总交易量" :
              key === "avgGasPrice" ? "平均Gas" :
              "最新区块"
            }</div>
            <div className="value" style={{ fontSize: 18 }}>{value}</div>
          </div>
        ))}
      </div>

      <div className="grid grid-2" style={{ gap: 16, marginBottom: 24 }}>
        {/* GraphQL 查询编辑器 */}
        <div className="card" style={{ padding: 24 }}>
          <div className="row" style={{ justifyContent: 'space-between', marginBottom: 16 }}>
            <div style={{ fontWeight: 700, display: 'flex', alignItems: 'center', gap: 8 }}>
              <Code size={20} color="#6366F1" />
              GraphQL 查询
            </div>
            <button
              className="btn"
              style={{ background: '#10b981', padding: '8px 16px' }}
              onClick={() => alert('查询功能开发中...')}
            >
              执行查询
            </button>
          </div>

          <textarea
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="输入 GraphQL 查询..."
            style={{
              width: '100%',
              minHeight: 200,
              background: 'rgba(0,0,0,.3)',
              border: '1px solid rgba(255,255,255,.1)',
              borderRadius: 8,
              padding: 16,
              color: '#fff',
              fontFamily: 'monospace',
              fontSize: 13,
              resize: 'vertical'
            }}
          />

          <div style={{ marginTop: 16 }}>
            <div style={{ fontWeight: 600, marginBottom: 8, fontSize: 14 }}>示例查询</div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
              {exampleQueries.map((example, index) => (
                <button
                  key={index}
                  style={{
                    background: 'rgba(255,255,255,.05)',
                    border: '1px solid rgba(255,255,255,.1)',
                    borderRadius: 8,
                    padding: 12,
                    textAlign: 'left',
                    cursor: 'pointer',
                    transition: 'all 0.2s'
                  }}
                  onClick={() => setQuery(example.query)}
                  onMouseEnter={(e) => e.currentTarget.style.background = 'rgba(255,255,255,.08)'}
                  onMouseLeave={(e) => e.currentTarget.style.background = 'rgba(255,255,255,.05)'}
                >
                  <div style={{ fontWeight: 600, fontSize: 13, marginBottom: 4 }}>{example.title}</div>
                  <div className="muted" style={{ fontSize: 12 }}>{example.description}</div>
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* 最近活动 */}
        <div className="card" style={{ padding: 24 }}>
          <div style={{ fontWeight: 700, marginBottom: 16, display: 'flex', alignItems: 'center', gap: 8 }}>
            <Activity size={20} color="#A855F7" />
            实时活动
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
            {recentActivity.map((activity, index) => (
              <div
                key={index}
                style={{
                  padding: 16,
                  background: 'rgba(255,255,255,.02)',
                  border: '1px solid rgba(255,255,255,.05)',
                  borderRadius: 12,
                  transition: 'all 0.2s'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.background = 'rgba(255,255,255,.05)';
                  e.currentTarget.style.borderColor = getTypeColor(activity.type);
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = 'rgba(255,255,255,.02)';
                  e.currentTarget.style.borderColor = 'rgba(255,255,255,.05)';
                }}
              >
                <div className="row" style={{ justifyContent: 'space-between', marginBottom: 8 }}>
                  <div style={{
                    background: `${getTypeColor(activity.type)}22`,
                    color: getTypeColor(activity.type),
                    padding: '4px 8px',
                    borderRadius: 6,
                    fontSize: 11,
                    fontWeight: 700
                  }}>
                    {activity.type}
                  </div>
                  <div className="muted" style={{ fontSize: 12 }}>{activity.time}</div>
                </div>
                <div style={{ fontSize: 13, marginBottom: 4 }}>
                  <span className="muted">From: </span>
                  <span style={{ fontFamily: 'monospace' }}>{activity.from}</span>
                </div>
                <div style={{ fontSize: 13, marginBottom: 4 }}>
                  <span className="muted">To: </span>
                  <span style={{ fontFamily: 'monospace' }}>{activity.to}</span>
                </div>
                <div style={{ fontWeight: 700, color: getTypeColor(activity.type) }}>
                  {activity.amount}
                </div>
              </div>
            ))}
          </div>

          <button
            className="btn"
            style={{
              width: '100%',
              marginTop: 16,
              background: 'rgba(255,255,255,.05)'
            }}
            onClick={() => alert('查看更多功能开发中...')}
          >
            查看更多活动
          </button>
        </div>
      </div>

      {/* API 端点 */}
      <div className="card" style={{ padding: 24, background: 'rgba(99, 102, 241, .1)', borderColor: '#6366F1' }}>
        <div style={{ fontWeight: 700, marginBottom: 12, display: 'flex', alignItems: 'center', gap: 8 }}>
          <Database size={20} color="#6366F1" />
          GraphQL 端点
        </div>
        <div style={{
          background: 'rgba(0,0,0,.3)',
          padding: 16,
          borderRadius: 8,
          fontFamily: 'monospace',
          fontSize: 14,
          marginBottom: 12
        }}>
          https://api.loyaltyx.io/subgraph/{selectedChain}
        </div>
        <div className="muted" style={{ fontSize: 13, lineHeight: 1.6 }}>
          使用此端点可以查询链上所有数据，包括用户、交易、质押、积分等信息。
          支持标准的 GraphQL 查询语法。
        </div>
      </div>
    </div>
  );
}
