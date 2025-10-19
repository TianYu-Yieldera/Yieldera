import React, { useEffect, useState } from "react";
import { Trophy, Medal, Crown, Star, TrendingUp } from "lucide-react";
import { config } from "../config/env";

const API_BASE = config.api.baseUrl;

export default function LeaderboardView() {
  const [leaderboard, setLeaderboard] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${API_BASE}/leaderboard`)
      .then(r => r.json())
      .then(data => {
        setLeaderboard(data.items || []);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

  const getRankIcon = (index) => {
    if (index === 0) return <Crown size={24} color="#FFD700" />;
    if (index === 1) return <Medal size={24} color="#C0C0C0" />;
    if (index === 2) return <Medal size={24} color="#CD7F32" />;
    return <Star size={20} color="#666" />;
  };

  const getRankStyle = (index) => {
    if (index === 0) return { background: 'linear-gradient(135deg, #FFD700 0%, #FFA500 100%)', color: '#000' };
    if (index === 1) return { background: 'linear-gradient(135deg, #C0C0C0 0%, #808080 100%)', color: '#000' };
    if (index === 2) return { background: 'linear-gradient(135deg, #CD7F32 0%, #8B4513 100%)', color: '#fff' };
    return {};
  };

  return (
    <div className="container">
      <div className="row" style={{ justifyContent: 'space-between', marginBottom: 24 }}>
        <div>
          <h1 style={{ margin: 0, fontSize: 32, display: 'flex', alignItems: 'center', gap: 12 }}>
            <Trophy size={36} color="#FFD700" />
            积分排行榜
          </h1>
          <p className="muted" style={{ marginTop: 8 }}>实时更新的用户积分排名</p>
        </div>
        <div className="row" style={{ gap: 16 }}>
          <div className="kpi" style={{ minWidth: 140 }}>
            <div className="title">总用户数</div>
            <div className="value">{leaderboard.length}</div>
          </div>
        </div>
      </div>

      {loading ? (
        <div className="card" style={{ padding: 48, textAlign: 'center' }}>
          <div className="muted">加载中...</div>
        </div>
      ) : leaderboard.length === 0 ? (
        <div className="card" style={{ padding: 48, textAlign: 'center' }}>
          <Trophy size={48} color="#666" style={{ margin: '0 auto 16px' }} />
          <div className="muted">暂无排行数据</div>
        </div>
      ) : (
        <>
          {/* 前三名特殊展示 */}
          <div className="grid" style={{ gridTemplateColumns: 'repeat(3, 1fr)', gap: 16, marginBottom: 24 }}>
            {leaderboard.slice(0, 3).map((item, index) => (
              <div key={index} className="card" style={{
                ...getRankStyle(index),
                padding: 24,
                textAlign: 'center',
                position: 'relative',
                overflow: 'hidden'
              }}>
                <div style={{ position: 'absolute', top: 12, right: 12, opacity: 0.3, fontSize: 64, fontWeight: 900 }}>
                  #{index + 1}
                </div>
                <div style={{ position: 'relative', zIndex: 1 }}>
                  <div style={{ marginBottom: 12 }}>
                    {getRankIcon(index)}
                  </div>
                  <div style={{ fontSize: 14, opacity: 0.8, marginBottom: 4 }}>
                    {item.Address?.slice(0, 6)}...{item.Address?.slice(-4)}
                  </div>
                  <div style={{ fontSize: 28, fontWeight: 800 }}>
                    {parseFloat(item.Points || 0).toLocaleString('en-US', { maximumFractionDigits: 2 })}
                  </div>
                  <div style={{ fontSize: 12, opacity: 0.8, marginTop: 4 }}>积分</div>
                </div>
              </div>
            ))}
          </div>

          {/* 其余排名列表 */}
          {leaderboard.length > 3 && (
            <div className="card">
              <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                <thead>
                  <tr style={{ borderBottom: '2px solid rgba(255,255,255,.1)' }}>
                    <th style={{ padding: 16, textAlign: 'left', width: 80 }}>排名</th>
                    <th style={{ padding: 16, textAlign: 'left' }}>地址</th>
                    <th style={{ padding: 16, textAlign: 'right' }}>积分</th>
                    <th style={{ padding: 16, textAlign: 'right', width: 120 }}>趋势</th>
                  </tr>
                </thead>
                <tbody>
                  {leaderboard.slice(3).map((item, index) => (
                    <tr key={index} style={{
                      borderTop: '1px solid rgba(255,255,255,.05)',
                      transition: 'all 0.2s',
                      cursor: 'pointer'
                    }}
                      onMouseEnter={(e) => e.currentTarget.style.background = 'rgba(255,255,255,.02)'}
                      onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
                    >
                      <td style={{ padding: 16 }}>
                        <div style={{
                          display: 'inline-flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          width: 36,
                          height: 36,
                          borderRadius: 8,
                          background: 'rgba(255,255,255,.05)',
                          fontWeight: 700
                        }}>
                          #{index + 4}
                        </div>
                      </td>
                      <td style={{ padding: 16, fontFamily: 'monospace', fontSize: 14 }}>
                        {item.Address}
                      </td>
                      <td style={{ padding: 16, textAlign: 'right', fontWeight: 700, fontSize: 16 }}>
                        {parseFloat(item.Points || 0).toLocaleString('en-US', { maximumFractionDigits: 2 })}
                      </td>
                      <td style={{ padding: 16, textAlign: 'right' }}>
                        <div className="row" style={{ justifyContent: 'flex-end', gap: 4 }}>
                          <TrendingUp size={16} color="#10b981" />
                          <span style={{ color: '#10b981', fontSize: 14 }}>+{Math.floor(Math.random() * 20)}%</span>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}
    </div>
  );
}
