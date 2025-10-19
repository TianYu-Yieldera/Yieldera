import React from "react";
import { TrendingUp, Gem, Zap, Shield, CheckCircle, ArrowRight, Users, BarChart3, Lock } from "lucide-react";
import { Link } from "react-router-dom";

export default function Landing(){
  return (
    <div className="container">
      {/* Hero Section */}
      <div className="card" style={{
        padding:'60px 40px',
        background:'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        textAlign: 'center',
        position: 'relative',
        overflow: 'hidden'
      }}>
        {/* 背景装饰 */}
        <div style={{
          position: 'absolute',
          top: '-50%',
          right: '-10%',
          width: '400px',
          height: '400px',
          background: 'rgba(255,255,255,0.1)',
          borderRadius: '50%',
          filter: 'blur(80px)'
        }} />
        <div style={{
          position: 'absolute',
          bottom: '-30%',
          left: '-5%',
          width: '300px',
          height: '300px',
          background: 'rgba(255,255,255,0.08)',
          borderRadius: '50%',
          filter: 'blur(60px)'
        }} />

        <div style={{position: 'relative', zIndex: 1}}>
          <div style={{
            display: 'inline-block',
            padding: '6px 18px',
            background: 'rgba(255,255,255,0.2)',
            borderRadius: '20px',
            fontSize: '13px',
            fontWeight: '600',
            color: '#fff',
            marginBottom: '20px',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(255,255,255,0.3)'
          }}>
            🚀 DeFi 收益聚合 × RWA 资产平台
          </div>

          <h1 style={{
            margin:'0',
            fontSize: '56px',
            fontWeight: '900',
            color: '#fff',
            letterSpacing: '-2px',
            lineHeight: '1.1'
          }}>
            Yieldera
          </h1>

          {/* 核心口号 - 震撼效果 */}
          <div style={{
            marginTop: '24px',
            position: 'relative',
            display: 'inline-block'
          }}>
            <style>{`
              @keyframes glow-pulse {
                0%, 100% {
                  text-shadow:
                    0 0 10px rgba(255,255,255,0.8),
                    0 0 20px rgba(255,255,255,0.6),
                    0 0 30px rgba(255,255,255,0.4),
                    0 0 40px rgba(102,126,234,0.6),
                    0 0 70px rgba(102,126,234,0.4),
                    0 0 80px rgba(102,126,234,0.2);
                }
                50% {
                  text-shadow:
                    0 0 20px rgba(255,255,255,1),
                    0 0 30px rgba(255,255,255,0.8),
                    0 0 40px rgba(255,255,255,0.6),
                    0 0 50px rgba(102,126,234,0.8),
                    0 0 80px rgba(102,126,234,0.6),
                    0 0 100px rgba(102,126,234,0.4);
                }
              }

              @keyframes shimmer {
                0% {
                  background-position: -1000px 0;
                }
                100% {
                  background-position: 1000px 0;
                }
              }

              .tagline-text {
                background: linear-gradient(
                  90deg,
                  rgba(255,255,255,0.6) 0%,
                  rgba(255,255,255,1) 25%,
                  rgba(255,255,255,1) 50%,
                  rgba(255,255,255,1) 75%,
                  rgba(255,255,255,0.6) 100%
                );
                background-size: 1000px 100%;
                -webkit-background-clip: text;
                -webkit-text-fill-color: transparent;
                background-clip: text;
                animation: shimmer 3s linear infinite, glow-pulse 2s ease-in-out infinite;
              }
            `}</style>

            <div style={{
              fontSize: '36px',
              fontWeight: '900',
              letterSpacing: '3px',
              textTransform: 'uppercase',
              position: 'relative',
              padding: '16px 0'
            }}>
              <div className="tagline-text">
                Enter the Yieldera
              </div>
            </div>

            {/* 装饰线条 */}
            <div style={{
              position: 'absolute',
              bottom: '0',
              left: '50%',
              transform: 'translateX(-50%)',
              width: '60%',
              height: '2px',
              background: 'linear-gradient(90deg, transparent, rgba(255,255,255,0.8), transparent)',
              boxShadow: '0 0 10px rgba(255,255,255,0.6)'
            }} />

            {/* 左右装饰点 */}
            <div style={{
              position: 'absolute',
              left: '-20px',
              top: '50%',
              transform: 'translateY(-50%)',
              width: '8px',
              height: '8px',
              borderRadius: '50%',
              background: '#fff',
              boxShadow: '0 0 20px rgba(255,255,255,0.8), 0 0 40px rgba(102,126,234,0.6)'
            }} />
            <div style={{
              position: 'absolute',
              right: '-20px',
              top: '50%',
              transform: 'translateY(-50%)',
              width: '8px',
              height: '8px',
              borderRadius: '50%',
              background: '#fff',
              boxShadow: '0 0 20px rgba(255,255,255,0.8), 0 0 40px rgba(102,126,234,0.6)'
            }} />
          </div>

          <p style={{
            marginTop:'24px',
            fontSize: '17px',
            color: 'rgba(255,255,255,0.9)',
            maxWidth: '680px',
            margin: '24px auto 0',
            lineHeight: '1.7',
            fontWeight: '400'
          }}>
            自动聚合多个 DeFi 协议，智能优化收益策略，用赚取的收益购买真实世界资产。<br/>
            一站式解决方案，让您的资产自动增值。
          </p>

          {/* CTA 按钮 */}
          <div style={{marginTop: '40px', display: 'flex', gap: '16px', justifyContent: 'center', flexWrap: 'wrap'}}>
            <Link
              to="/vault"
              style={{
                padding: '16px 36px',
                background: '#fff',
                color: '#667eea',
                borderRadius: '10px',
                fontWeight: '700',
                fontSize: '16px',
                textDecoration: 'none',
                transition: 'all 0.3s',
                boxShadow: '0 8px 24px rgba(0,0,0,0.25)',
                display: 'flex',
                alignItems: 'center',
                gap: '8px'
              }}
            >
              开始赚取收益 <ArrowRight size={18} />
            </Link>

            <Link
              to="/rwa-market"
              style={{
                padding: '16px 36px',
                background: 'rgba(255,255,255,0.15)',
                color: '#fff',
                border: '2px solid rgba(255,255,255,0.4)',
                borderRadius: '10px',
                fontWeight: '700',
                fontSize: '16px',
                textDecoration: 'none',
                backdropFilter: 'blur(10px)',
                transition: 'all 0.3s'
              }}
            >
              浏览 RWA 资产
            </Link>
          </div>
        </div>
      </div>

      {/* 核心数据展示 */}
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))',
        gap: '20px',
        marginTop: '32px'
      }}>
        <div className="card" style={{
          padding: '28px 24px',
          textAlign: 'center',
          background: 'linear-gradient(135deg, rgba(102,126,234,0.1), rgba(118,75,162,0.05))',
          border: '1px solid rgba(102,126,234,0.2)'
        }}>
          <div style={{fontSize: '40px', fontWeight: '900', color: '#667eea', marginBottom: '4px'}}>8.45%</div>
          <div className="muted" style={{fontSize: '14px', fontWeight: '500'}}>平均年化收益</div>
          <div style={{fontSize: '12px', color: '#43e97b', marginTop: '8px', fontWeight: '600'}}>↑ 12.3% 本月</div>
        </div>

        <div className="card" style={{
          padding: '28px 24px',
          textAlign: 'center',
          background: 'linear-gradient(135deg, rgba(240,147,251,0.1), rgba(245,87,108,0.05))',
          border: '1px solid rgba(240,147,251,0.2)'
        }}>
          <div style={{fontSize: '40px', fontWeight: '900', color: '#f093fb', marginBottom: '4px'}}>$125.6M</div>
          <div className="muted" style={{fontSize: '14px', fontWeight: '500'}}>总锁仓价值 (TVL)</div>
          <div style={{fontSize: '12px', color: '#43e97b', marginTop: '8px', fontWeight: '600'}}>↑ 18.7% 本周</div>
        </div>

        <div className="card" style={{
          padding: '28px 24px',
          textAlign: 'center',
          background: 'linear-gradient(135deg, rgba(67,233,123,0.1), rgba(56,178,172,0.05))',
          border: '1px solid rgba(67,233,123,0.2)'
        }}>
          <div style={{fontSize: '40px', fontWeight: '900', color: '#43e97b', marginBottom: '4px'}}>42,341</div>
          <div className="muted" style={{fontSize: '14px', fontWeight: '500'}}>活跃用户</div>
          <div style={{fontSize: '12px', color: '#43e97b', marginTop: '8px', fontWeight: '600'}}>+2,156 本周</div>
        </div>

        <div className="card" style={{
          padding: '28px 24px',
          textAlign: 'center',
          background: 'linear-gradient(135deg, rgba(79,172,254,0.1), rgba(0,242,254,0.05))',
          border: '1px solid rgba(79,172,254,0.2)'
        }}>
          <div style={{fontSize: '40px', fontWeight: '900', color: '#4facfe', marginBottom: '4px'}}>24/7</div>
          <div className="muted" style={{fontSize: '14px', fontWeight: '500'}}>自动再平衡</div>
          <div style={{fontSize: '12px', color: '#9ca3af', marginTop: '8px', fontWeight: '600'}}>智能优化</div>
        </div>
      </div>

      {/* 核心价值主张 */}
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))',
        gap: '24px',
        marginTop: '48px'
      }}>
        <div className="card" style={{padding: '32px'}}>
          <div style={{
            width: '56px',
            height: '56px',
            borderRadius: '14px',
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            marginBottom: '20px',
            boxShadow: '0 8px 16px rgba(102,126,234,0.3)'
          }}>
            <TrendingUp size={28} color="#fff" strokeWidth={2.5} />
          </div>
          <h3 style={{marginTop: 0, marginBottom: '12px', fontSize: '20px', fontWeight: '700'}}>智能收益优化</h3>
          <p className="muted" style={{fontSize: '15px', lineHeight: '1.7', marginBottom: '16px'}}>
            自动监控 Aave、Compound、Uniswap V3、GMX 等多个协议的 APY，每 24 小时智能再平衡，确保您始终获得最优收益。
          </p>
          <div style={{display: 'flex', alignItems: 'center', gap: '8px', color: '#667eea', fontSize: '14px', fontWeight: '600'}}>
            <CheckCircle size={16} />
            平均年化 8.45%
          </div>
        </div>

        <div className="card" style={{padding: '32px'}}>
          <div style={{
            width: '56px',
            height: '56px',
            borderRadius: '14px',
            background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            marginBottom: '20px',
            boxShadow: '0 8px 16px rgba(240,147,251,0.3)'
          }}>
            <Gem size={28} color="#fff" strokeWidth={2.5} />
          </div>
          <h3 style={{marginTop: 0, marginBottom: '12px', fontSize: '20px', fontWeight: '700'}}>真实资产投资</h3>
          <p className="muted" style={{fontSize: '15px', lineHeight: '1.7', marginBottom: '16px'}}>
            用收益直接购买 Apple、Tesla、美国国债、黄金等 RWA 资产，实现 DeFi 收益到真实世界资产的无缝转换。
          </p>
          <div style={{display: 'flex', alignItems: 'center', gap: '8px', color: '#f093fb', fontSize: '14px', fontWeight: '600'}}>
            <CheckCircle size={16} />
            12+ 资产类型
          </div>
        </div>

        <div className="card" style={{padding: '32px'}}>
          <div style={{
            width: '56px',
            height: '56px',
            borderRadius: '14px',
            background: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            marginBottom: '20px',
            boxShadow: '0 8px 16px rgba(79,172,254,0.3)'
          }}>
            <Zap size={28} color="#fff" strokeWidth={2.5} />
          </div>
          <h3 style={{marginTop: 0, marginBottom: '12px', fontSize: '20px', fontWeight: '700'}}>零操作成本</h3>
          <p className="muted" style={{fontSize: '15px', lineHeight: '1.7', marginBottom: '16px'}}>
            存入后完全自动化，无需手动管理。系统自动收割、复投、再平衡，您只需要享受收益增长。
          </p>
          <div style={{display: 'flex', alignItems: 'center', gap: '8px', color: '#4facfe', fontSize: '14px', fontWeight: '600'}}>
            <CheckCircle size={16} />
            全自动运行
          </div>
        </div>
      </div>

      {/* 集成协议展示 */}
      <div className="card" style={{marginTop: '48px', padding: '40px', background: '#1a1d29'}}>
        <h2 style={{
          marginTop: 0,
          marginBottom: '16px',
          fontSize: '28px',
          fontWeight: '800',
          textAlign: 'center'
        }}>
          集成顶级 DeFi 协议
        </h2>
        <p className="muted" style={{
          textAlign: 'center',
          fontSize: '16px',
          marginBottom: '40px',
          maxWidth: '600px',
          margin: '0 auto 40px'
        }}>
          与行业领先的 DeFi 协议深度集成，为您提供最安全、最高效的收益策略
        </p>

        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
          gap: '20px'
        }}>
          {[
            { name: 'Aave', apy: '7.2%', tvl: '$12.3B', color: '#B6509E' },
            { name: 'Compound', apy: '6.8%', tvl: '$3.1B', color: '#00D395' },
            { name: 'Uniswap V3', apy: '12.5%', tvl: '$4.2B', color: '#FF007A' },
            { name: 'GMX', apy: '15.3%', tvl: '$582M', color: '#3B82F6' }
          ].map((protocol) => (
            <div key={protocol.name} style={{
              padding: '24px',
              background: 'rgba(255,255,255,0.03)',
              borderRadius: '12px',
              border: '1px solid rgba(255,255,255,0.08)',
              textAlign: 'center',
              transition: 'all 0.3s'
            }}>
              <div style={{
                width: '48px',
                height: '48px',
                borderRadius: '12px',
                background: protocol.color,
                margin: '0 auto 16px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontSize: '20px',
                fontWeight: 'bold',
                color: '#fff'
              }}>
                {protocol.name[0]}
              </div>
              <div style={{fontSize: '18px', fontWeight: '700', marginBottom: '8px'}}>{protocol.name}</div>
              <div style={{fontSize: '14px', color: '#43e97b', fontWeight: '600', marginBottom: '4px'}}>
                APY: {protocol.apy}
              </div>
              <div className="muted" style={{fontSize: '13px'}}>TVL: {protocol.tvl}</div>
            </div>
          ))}
        </div>
      </div>

      {/* 用户数据和信任指标 */}
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
        gap: '24px',
        marginTop: '48px'
      }}>
        <div className="card" style={{padding: '32px', textAlign: 'center'}}>
          <div style={{
            width: '64px',
            height: '64px',
            borderRadius: '50%',
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            margin: '0 auto 20px',
            boxShadow: '0 8px 24px rgba(102,126,234,0.4)'
          }}>
            <Users size={32} color="#fff" strokeWidth={2.5} />
          </div>
          <div style={{fontSize: '36px', fontWeight: '900', color: '#fff', marginBottom: '8px'}}>42,341</div>
          <div className="muted" style={{fontSize: '15px', fontWeight: '500'}}>全球活跃用户</div>
          <div style={{marginTop: '12px', fontSize: '13px', color: '#43e97b', fontWeight: '600'}}>
            +2,156 本周新增
          </div>
        </div>

        <div className="card" style={{padding: '32px', textAlign: 'center'}}>
          <div style={{
            width: '64px',
            height: '64px',
            borderRadius: '50%',
            background: 'linear-gradient(135deg, #43e97b 0%, #38b2ac 100%)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            margin: '0 auto 20px',
            boxShadow: '0 8px 24px rgba(67,233,123,0.4)'
          }}>
            <BarChart3 size={32} color="#fff" strokeWidth={2.5} />
          </div>
          <div style={{fontSize: '36px', fontWeight: '900', color: '#fff', marginBottom: '8px'}}>$125.6M</div>
          <div className="muted" style={{fontSize: '15px', fontWeight: '500'}}>总锁仓价值</div>
          <div style={{marginTop: '12px', fontSize: '13px', color: '#43e97b', fontWeight: '600'}}>
            +18.7% 本周增长
          </div>
        </div>

        <div className="card" style={{padding: '32px', textAlign: 'center'}}>
          <div style={{
            width: '64px',
            height: '64px',
            borderRadius: '50%',
            background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            margin: '0 auto 20px',
            boxShadow: '0 8px 24px rgba(240,147,251,0.4)'
          }}>
            <Lock size={32} color="#fff" strokeWidth={2.5} />
          </div>
          <div style={{fontSize: '36px', fontWeight: '900', color: '#fff', marginBottom: '8px'}}>100%</div>
          <div className="muted" style={{fontSize: '15px', fontWeight: '500'}}>安全保障</div>
          <div style={{marginTop: '12px', fontSize: '13px', color: '#9ca3af', fontWeight: '600'}}>
            多重审计认证
          </div>
        </div>
      </div>

      {/* CTA Section */}
      <div className="card" style={{
        marginTop: '64px',
        padding: '48px 40px',
        background: 'linear-gradient(135deg, rgba(102,126,234,0.15), rgba(118,75,162,0.15))',
        border: '2px solid rgba(102,126,234,0.3)',
        textAlign: 'center'
      }}>
        <h2 style={{
          marginTop: 0,
          fontSize: '32px',
          fontWeight: '800',
          marginBottom: '16px'
        }}>
          准备好开始了吗？
        </h2>
        <p className="muted" style={{
          fontSize: '17px',
          marginBottom: '32px',
          maxWidth: '600px',
          margin: '0 auto 32px',
          lineHeight: '1.6'
        }}>
          连接钱包，立即开始赚取收益。无需复杂操作，几分钟即可完成设置。
        </p>
        <Link
          to="/vault"
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: '10px',
            padding: '18px 40px',
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            color: '#fff',
            borderRadius: '12px',
            fontWeight: '700',
            fontSize: '18px',
            textDecoration: 'none',
            boxShadow: '0 8px 24px rgba(102,126,234,0.4)',
            transition: 'all 0.3s'
          }}
        >
          立即开始 <ArrowRight size={20} />
        </Link>
      </div>
    </div>
  );
}
