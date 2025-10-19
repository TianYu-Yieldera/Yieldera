import React, { useEffect, useState } from "react";
import { Award, Star, Shield, Zap, Crown, Trophy, Medal, Lock } from "lucide-react";
import { useWallet } from "../web3/WalletContext";
import { config } from "../config/env";

const API_BASE = config.api.baseUrl;

export default function BadgesView() {
  const { address } = useWallet();
  const [ownedBadges, setOwnedBadges] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!address) {
      setLoading(false);
      return;
    }
    fetch(`${API_BASE}/users/${address}/badges`)
      .then(r => r.json())
      .then(data => {
        setOwnedBadges(data.badges || []);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [address]);

  // 所有可用徽章
  const allBadges = [
    {
      code: "EARLY_BIRD",
      name: "早鸟用户",
      description: "前100名注册用户",
      icon: Star,
      color: "#FFD700",
      rarity: "传说",
      requirement: "成为前100名用户"
    },
    {
      code: "STAKER_BRONZE",
      name: "青铜质押者",
      description: "质押超过 1,000 代币",
      icon: Shield,
      color: "#CD7F32",
      rarity: "普通",
      requirement: "质押 ≥ 1,000 代币"
    },
    {
      code: "STAKER_SILVER",
      name: "白银质押者",
      description: "质押超过 10,000 代币",
      icon: Shield,
      color: "#C0C0C0",
      rarity: "稀有",
      requirement: "质押 ≥ 10,000 代币"
    },
    {
      code: "STAKER_GOLD",
      name: "黄金质押者",
      description: "质押超过 100,000 代币",
      icon: Crown,
      color: "#FFD700",
      rarity: "史诗",
      requirement: "质押 ≥ 100,000 代币"
    },
    {
      code: "POINTS_1K",
      name: "积分新手",
      description: "累积 1,000 积分",
      icon: Zap,
      color: "#10b981",
      rarity: "普通",
      requirement: "获得 1,000 积分"
    },
    {
      code: "POINTS_10K",
      name: "积分大师",
      description: "累积 10,000 积分",
      icon: Trophy,
      color: "#6366F1",
      rarity: "稀有",
      requirement: "获得 10,000 积分"
    },
    {
      code: "POINTS_100K",
      name: "积分传说",
      description: "累积 100,000 积分",
      icon: Crown,
      color: "#A855F7",
      rarity: "传说",
      requirement: "获得 100,000 积分"
    },
    {
      code: "TOP_10",
      name: "排行榜前十",
      description: "进入积分排行榜前10名",
      icon: Medal,
      color: "#F59E0B",
      rarity: "史诗",
      requirement: "排名进入前10"
    },
    {
      code: "TRADER",
      name: "活跃交易者",
      description: "完成 100 笔交易",
      icon: Zap,
      color: "#22D3EE",
      rarity: "稀有",
      requirement: "完成 100 笔交易"
    },
    {
      code: "LOYAL",
      name: "忠实用户",
      description: "连续活跃 30 天",
      icon: Star,
      color: "#EC4899",
      rarity: "史诗",
      requirement: "连续活跃 30 天"
    }
  ];

  const getRarityColor = (rarity) => {
    switch(rarity) {
      case "传说": return "#FFD700";
      case "史诗": return "#A855F7";
      case "稀有": return "#6366F1";
      default: return "#9CA3AF";
    }
  };

  const isOwned = (code) => ownedBadges.includes(code);

  return (
    <div className="container">
      <div style={{ marginBottom: 24 }}>
        <h1 style={{ margin: 0, fontSize: 32, display: 'flex', alignItems: 'center', gap: 12 }}>
          <Award size={36} color="#FFD700" />
          成就徽章
        </h1>
        <p className="muted" style={{ marginTop: 8 }}>解锁徽章，展示你的成就</p>
      </div>

      {/* 统计 */}
      <div className="grid" style={{ gridTemplateColumns: 'repeat(4, 1fr)', gap: 16, marginBottom: 24 }}>
        <div className="kpi">
          <div className="title">已解锁</div>
          <div className="value" style={{ color: '#10b981' }}>{ownedBadges.length}</div>
        </div>
        <div className="kpi">
          <div className="title">总徽章数</div>
          <div className="value">{allBadges.length}</div>
        </div>
        <div className="kpi">
          <div className="title">完成度</div>
          <div className="value" style={{ color: '#6366F1' }}>
            {allBadges.length > 0 ? Math.floor(ownedBadges.length / allBadges.length * 100) : 0}%
          </div>
        </div>
        <div className="kpi">
          <div className="title">稀有度</div>
          <div className="value" style={{ color: '#A855F7' }}>
            {ownedBadges.filter(code => {
              const badge = allBadges.find(b => b.code === code);
              return badge && (badge.rarity === '传说' || badge.rarity === '史诗');
            }).length}
          </div>
        </div>
      </div>

      {/* 徽章网格 */}
      {loading ? (
        <div className="card" style={{ padding: 48, textAlign: 'center' }}>
          <div className="muted">加载中...</div>
        </div>
      ) : !address ? (
        <div className="card" style={{ padding: 48, textAlign: 'center', background: 'rgba(245, 158, 11, .1)', borderColor: '#F59E0B' }}>
          <Lock size={48} color="#F59E0B" style={{ margin: '0 auto 16px' }} />
          <div style={{ fontWeight: 700, marginBottom: 8 }}>请先连接钱包</div>
          <div className="muted">连接钱包后即可查看你的徽章收藏</div>
        </div>
      ) : (
        <div className="grid" style={{ gridTemplateColumns: 'repeat(3, 1fr)', gap: 16 }}>
          {allBadges.map((badge, index) => {
            const owned = isOwned(badge.code);
            const Icon = badge.icon;

            return (
              <div
                key={index}
                className="card"
                style={{
                  padding: 24,
                  position: 'relative',
                  overflow: 'hidden',
                  opacity: owned ? 1 : 0.5,
                  borderColor: owned ? badge.color : 'rgba(255,255,255,.1)',
                  cursor: 'pointer',
                  transition: 'all 0.3s'
                }}
                onMouseEnter={(e) => {
                  if (owned) {
                    e.currentTarget.style.transform = 'translateY(-4px)';
                    e.currentTarget.style.boxShadow = `0 20px 40px ${badge.color}33`;
                  }
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.transform = 'translateY(0)';
                  e.currentTarget.style.boxShadow = '0 10px 24px rgba(0,0,0,.2)';
                }}
              >
                {/* 稀有度标签 */}
                <div style={{
                  position: 'absolute',
                  top: 12,
                  right: 12,
                  background: getRarityColor(badge.rarity),
                  color: badge.rarity === '普通' ? '#000' : '#fff',
                  padding: '4px 8px',
                  borderRadius: 6,
                  fontSize: 11,
                  fontWeight: 700
                }}>
                  {badge.rarity}
                </div>

                {/* 未解锁遮罩 */}
                {!owned && (
                  <div style={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    right: 0,
                    bottom: 0,
                    background: 'rgba(0,0,0,.6)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    zIndex: 2
                  }}>
                    <Lock size={32} color="#666" />
                  </div>
                )}

                {/* 图标 */}
                <div style={{
                  width: 80,
                  height: 80,
                  margin: '0 auto 16px',
                  background: owned ? `${badge.color}22` : 'rgba(255,255,255,.05)',
                  borderRadius: 16,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  position: 'relative',
                  zIndex: 1
                }}>
                  <Icon size={40} color={owned ? badge.color : '#666'} />
                </div>

                {/* 信息 */}
                <div style={{ textAlign: 'center', position: 'relative', zIndex: 1 }}>
                  <div style={{ fontWeight: 700, fontSize: 16, marginBottom: 4 }}>
                    {badge.name}
                  </div>
                  <div className="muted" style={{ fontSize: 13, marginBottom: 12 }}>
                    {badge.description}
                  </div>
                  <div style={{
                    background: 'rgba(255,255,255,.05)',
                    padding: '8px 12px',
                    borderRadius: 8,
                    fontSize: 12
                  }}>
                    {badge.requirement}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
