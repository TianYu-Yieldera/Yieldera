import React, { useEffect, useState } from "react";
import { TrendingUp, Gem, Shield, CheckCircle, ArrowRight, Zap, Lock, Globe, BarChart3, Target, Repeat, DollarSign } from "lucide-react";
import { Link } from "react-router-dom";
import ProtocolNetwork3D from "../components/ProtocolNetwork3D";

export default function Landing(){
  const [tvl, setTvl] = useState(125.6);
  const [users, setUsers] = useState(42341);

  // Simulate real-time data updates
  useEffect(() => {
    const interval = setInterval(() => {
      setTvl(prev => prev + (Math.random() - 0.5) * 0.1);
      setUsers(prev => prev + Math.floor(Math.random() * 3));
    }, 3000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div style={{ minHeight: '100vh', background: 'rgb(248, 250, 252)' }}>
      <style>{`
        @keyframes float {
          0%, 100% { transform: translateY(0px); }
          50% { transform: translateY(-10px); }
        }

        @keyframes pulse-glow {
          0%, 100% { box-shadow: 0 0 20px rgba(34, 211, 238, 0.3); }
          50% { box-shadow: 0 0 40px rgba(34, 211, 238, 0.6); }
        }

        @keyframes counter-up {
          from { opacity: 0.5; transform: scale(0.95); }
          to { opacity: 1; transform: scale(1); }
        }

        @keyframes star-twinkle {
          0%, 100% { opacity: 0.2; transform: scale(0.8); }
          50% { opacity: 1; transform: scale(1); }
        }

        @keyframes page-fade-in {
          from {
            opacity: 0;
            transform: translateY(20px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }

        .page-container {
          animation: page-fade-in 0.8s ease-out;
        }

        .stat-card {
          position: relative;
          transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
          background: linear-gradient(135deg, #ffffff 0%, #fafbfc 100%);
          box-shadow:
            0 4px 20px rgba(0, 0, 0, 0.05),
            0 0 0 1px rgba(34, 211, 238, 0.08),
            inset 0 1px 0 rgba(255, 255, 255, 0.1);
        }

        .stat-card::before {
          content: '';
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          height: 3px;
          background: linear-gradient(90deg,
            rgba(34, 211, 238, 0) 0%,
            rgba(34, 211, 238, 0.5) 50%,
            rgba(34, 211, 238, 0) 100%
          );
          opacity: 0;
          transition: opacity 0.3s ease;
        }

        .stat-card:hover {
          transform: translateY(-6px);
          box-shadow:
            0 12px 32px rgba(29, 78, 216, 0.12),
            0 0 0 1px rgba(34, 211, 238, 0.2),
            inset 0 1px 0 rgba(255, 255, 255, 0.2),
            0 0 30px rgba(34, 211, 238, 0.15);
        }

        .stat-card:hover::before {
          opacity: 1;
        }

        .feature-card {
          transition: all 0.3s ease;
        }

        .feature-card:hover {
          transform: translateY(-8px);
          box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
        }

        .protocol-card {
          transition: all 0.3s ease;
        }

        .protocol-card:hover {
          transform: translateY(-4px) scale(1.05);
          box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
        }

        .counter {
          animation: counter-up 0.5s ease-out;
        }
      `}</style>

      <div className="page-container" style={{ maxWidth: 1280, margin: '0 auto', padding: '24px 24px' }}>
        {/* Hero Section - Enhanced with 3D Visualization */}
        <div style={{
          background: 'radial-gradient(ellipse at center, rgb(10, 25, 47) 0%, rgb(5, 10, 20) 50%, rgb(0, 0, 0) 100%)',
          borderRadius: 16,
          padding: '24px 40px 32px',
          textAlign: 'center',
          position: 'relative',
          overflow: 'hidden',
          boxShadow: '0 25px 80px rgba(0, 0, 0, 0.6)',
          border: '1px solid rgba(34, 211, 238, 0.2)'
        }}>
          {/* Deep space stars layer 1 - far away small stars */}
          {[...Array(80)].map((_, i) => (
            <div
              key={`star-hero-far-${i}`}
              style={{
                position: 'absolute',
                width: Math.random() * 1.5 + 0.5 + 'px',
                height: Math.random() * 1.5 + 0.5 + 'px',
                background: 'rgba(255, 255, 255, 0.4)',
                borderRadius: '50%',
                left: Math.random() * 100 + '%',
                top: Math.random() * 100 + '%',
                boxShadow: `0 0 ${Math.random() * 2 + 1}px rgba(255, 255, 255, 0.3)`,
                animation: `star-twinkle ${Math.random() * 4 + 3}s ease-in-out infinite ${Math.random() * 4}s`,
                pointerEvents: 'none'
              }}
            />
          ))}

          {/* Medium distance stars - slightly bigger and brighter */}
          {[...Array(40)].map((_, i) => (
            <div
              key={`star-hero-mid-${i}`}
              style={{
                position: 'absolute',
                width: Math.random() * 2 + 1 + 'px',
                height: Math.random() * 2 + 1 + 'px',
                background: i % 3 === 0
                  ? 'rgba(173, 216, 255, 0.8)' // Bluish stars
                  : 'rgba(255, 255, 255, 0.7)',
                borderRadius: '50%',
                left: Math.random() * 100 + '%',
                top: Math.random() * 100 + '%',
                boxShadow: `0 0 ${Math.random() * 3 + 2}px currentColor`,
                animation: `star-twinkle ${Math.random() * 3 + 2}s ease-in-out infinite ${Math.random() * 3}s`,
                pointerEvents: 'none'
              }}
            />
          ))}

          {/* Close bright stars - largest and most prominent */}
          {[...Array(20)].map((_, i) => (
            <div
              key={`star-hero-close-${i}`}
              style={{
                position: 'absolute',
                width: Math.random() * 3 + 2 + 'px',
                height: Math.random() * 3 + 2 + 'px',
                background: i % 4 === 0
                  ? 'rgba(34, 211, 238, 1)' // Cyan accent stars
                  : i % 4 === 1
                  ? 'rgba(173, 216, 255, 1)' // Blue stars
                  : 'rgba(255, 255, 255, 1)',
                borderRadius: '50%',
                left: Math.random() * 100 + '%',
                top: Math.random() * 100 + '%',
                boxShadow: `0 0 ${Math.random() * 6 + 4}px currentColor`,
                animation: `star-twinkle ${Math.random() * 2.5 + 1.5}s ease-in-out infinite ${Math.random() * 2}s`,
                pointerEvents: 'none'
              }}
            />
          ))}

          {/* Grid background - subtle */}
          <div style={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundImage: `
              linear-gradient(rgba(34, 211, 238, 0.03) 1px, transparent 1px),
              linear-gradient(90deg, rgba(34, 211, 238, 0.03) 1px, transparent 1px)
            `,
            backgroundSize: '40px 40px',
            opacity: 0.3
          }} />

          {/* Nebula-like cyan glow effects */}
          <div style={{
            position: 'absolute',
            top: '-50%',
            right: '10%',
            width: '600px',
            height: '600px',
            background: 'radial-gradient(circle, rgba(34, 211, 238, 0.12) 0%, rgba(59, 130, 246, 0.08) 40%, transparent 70%)',
            filter: 'blur(80px)',
            animation: 'float 6s ease-in-out infinite'
          }} />
          <div style={{
            position: 'absolute',
            bottom: '-50%',
            left: '10%',
            width: '600px',
            height: '600px',
            background: 'radial-gradient(circle, rgba(59, 130, 246, 0.1) 0%, rgba(147, 51, 234, 0.06) 40%, transparent 70%)',
            filter: 'blur(80px)',
            animation: 'float 8s ease-in-out infinite'
          }} />

          <div style={{position: 'relative', zIndex: 1}}>
            <h1 style={{
              margin:'0 0 12px 0',
              fontSize: '56px',
              fontWeight: '900',
              color: '#fff',
              letterSpacing: '-2px',
              lineHeight: '1.1',
              textShadow: `
                0 0 30px rgba(34, 211, 238, 0.6),
                0 0 60px rgba(34, 211, 238, 0.4),
                0 0 90px rgba(34, 211, 238, 0.3),
                0 0 120px rgba(34, 211, 238, 0.2),
                0 4px 16px rgba(0, 0, 0, 0.5)
              `
            }}>
              Yieldera
            </h1>

            {/* Bold Professional Tagline - Equal importance to brand */}
            <h2 style={{
              fontSize: '32px',
              fontWeight: '600',
              letterSpacing: '-0.5px',
              margin: '0 0 16px 0',
              color: 'rgba(255, 255, 255, 0.95)',
              textShadow: `
                0 0 25px rgba(34, 211, 238, 0.5),
                0 0 50px rgba(34, 211, 238, 0.35),
                0 0 75px rgba(34, 211, 238, 0.25),
                0 0 100px rgba(34, 211, 238, 0.15),
                0 2px 8px rgba(0, 0, 0, 0.3)
              `,
              lineHeight: '1.2'
            }}>
              Multi-Chain DeFi + RWA Treasury
            </h2>

            <p style={{
              fontSize: '13px',
              color: 'rgba(203, 213, 225, 0.6)',
              maxWidth: '700px',
              margin: '0 auto 24px',
              lineHeight: '1.5',
              fontWeight: '400',
              opacity: 0.8
            }}>
              <strong style={{ color: 'rgba(34, 211, 238, 0.8)' }}>Arbitrum</strong> DeFi (Aave, GMX, Uniswap) + <strong style={{ color: 'rgba(34, 211, 238, 0.8)' }}>Base</strong> Treasury + <strong style={{ color: 'rgba(34, 211, 238, 0.8)' }}>AI Risk</strong> (10K agents, VaR/CVaR)
            </p>

            {/* 3D Protocol Network Visualization */}
            <div style={{ margin: '0 auto', maxWidth: '900px' }}>
              <ProtocolNetwork3D />
            </div>
          </div>
        </div>

        {/* Live Stats with animated counters */}
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
          gap: '16px',
          marginTop: '24px'
        }}>
          <div className="stat-card counter" style={{
            padding: '28px 20px',
            textAlign: 'center',
            borderRadius: 12
          }}>
            <div style={{fontSize: '32px', fontWeight: '700', color: 'rgb(29, 78, 216)', marginBottom: '8px'}}>2</div>
            <div style={{fontSize: '14px', fontWeight: '600', color: 'rgb(100, 116, 139)'}}>L2 Chains</div>
            <div style={{fontSize: '12px', color: 'rgb(22, 163, 74)', marginTop: '8px', fontWeight: '600'}}>Arbitrum + Base</div>
          </div>

          <div className="stat-card counter" style={{
            padding: '28px 20px',
            textAlign: 'center',
            borderRadius: 12
          }}>
            <div style={{fontSize: '32px', fontWeight: '700', color: 'rgb(29, 78, 216)', marginBottom: '8px'}}>6+</div>
            <div style={{fontSize: '14px', fontWeight: '600', color: 'rgb(100, 116, 139)'}}>DeFi Protocols</div>
            <div style={{fontSize: '12px', color: 'rgb(22, 163, 74)', marginTop: '8px', fontWeight: '600'}}>Aave, Compound, GMX...</div>
          </div>

          <div className="stat-card counter" style={{
            padding: '28px 20px',
            textAlign: 'center',
            borderRadius: 12
          }}>
            <div style={{fontSize: '32px', fontWeight: '700', color: 'rgb(29, 78, 216)', marginBottom: '8px'}}>5</div>
            <div style={{fontSize: '14px', fontWeight: '600', color: 'rgb(100, 116, 139)'}}>AI Risk Functions</div>
            <div style={{fontSize: '12px', color: 'rgb(22, 163, 74)', marginTop: '8px', fontWeight: '600'}}>VaR, CVaR, ML, Agent Sim</div>
          </div>

          <div className="stat-card counter" style={{
            padding: '28px 20px',
            textAlign: 'center',
            borderRadius: 12
          }}>
            <div style={{fontSize: '32px', fontWeight: '700', color: 'rgb(29, 78, 216)', marginBottom: '8px'}}>10K+</div>
            <div style={{fontSize: '14px', fontWeight: '600', color: 'rgb(100, 116, 139)'}}>Simulated Agents</div>
            <div style={{fontSize: '12px', color: 'rgb(22, 163, 74)', marginTop: '8px', fontWeight: '600'}}>Monte Carlo Risk Analysis</div>
          </div>

          <div className="stat-card counter" style={{
            padding: '28px 20px',
            textAlign: 'center',
            borderRadius: 12
          }}>
            <div style={{fontSize: '32px', fontWeight: '700', color: 'rgb(29, 78, 216)', marginBottom: '8px'}}>24/7</div>
            <div style={{fontSize: '14px', fontWeight: '600', color: 'rgb(100, 116, 139)'}}>Liquidation Monitoring</div>
            <div style={{fontSize: '12px', color: 'rgb(22, 163, 74)', marginTop: '8px', fontWeight: '600'}}>Real-time Health Factor</div>
          </div>
        </div>

        {/* Why Yieldera Section */}
        <div style={{
          marginTop: '80px',
          padding: '64px 40px',
          background: 'radial-gradient(ellipse at top, rgb(10, 25, 47) 0%, rgb(5, 10, 20) 40%, rgb(0, 0, 0) 100%)',
          borderRadius: 16,
          position: 'relative',
          overflow: 'hidden',
          border: '1px solid rgba(34, 211, 238, 0.2)'
        }}>
          {/* Deep space stars - far layer */}
          {[...Array(60)].map((_, i) => (
            <div
              key={`star-why-far-${i}`}
              style={{
                position: 'absolute',
                width: Math.random() * 1.5 + 0.5 + 'px',
                height: Math.random() * 1.5 + 0.5 + 'px',
                background: 'rgba(255, 255, 255, 0.4)',
                borderRadius: '50%',
                left: Math.random() * 100 + '%',
                top: Math.random() * 100 + '%',
                boxShadow: `0 0 ${Math.random() * 2 + 1}px rgba(255, 255, 255, 0.3)`,
                animation: `star-twinkle ${Math.random() * 4 + 3}s ease-in-out infinite ${Math.random() * 4}s`,
                pointerEvents: 'none'
              }}
            />
          ))}

          {/* Medium distance bright stars */}
          {[...Array(35)].map((_, i) => (
            <div
              key={`star-why-mid-${i}`}
              style={{
                position: 'absolute',
                width: Math.random() * 2.5 + 1 + 'px',
                height: Math.random() * 2.5 + 1 + 'px',
                background: i % 4 === 0
                  ? 'rgba(34, 211, 238, 0.9)' // Cyan stars
                  : i % 4 === 1
                  ? 'rgba(173, 216, 255, 0.8)' // Blue stars
                  : i % 4 === 2
                  ? 'rgba(147, 197, 253, 0.7)' // Light blue
                  : 'rgba(255, 255, 255, 0.7)',
                borderRadius: '50%',
                left: Math.random() * 100 + '%',
                top: Math.random() * 100 + '%',
                boxShadow: `0 0 ${Math.random() * 4 + 2}px currentColor`,
                animation: `star-twinkle ${Math.random() * 3 + 2}s ease-in-out infinite ${Math.random() * 3}s`,
                pointerEvents: 'none'
              }}
            />
          ))}

          {/* Close prominent stars */}
          {[...Array(15)].map((_, i) => (
            <div
              key={`star-why-close-${i}`}
              style={{
                position: 'absolute',
                width: Math.random() * 3 + 2 + 'px',
                height: Math.random() * 3 + 2 + 'px',
                background: i % 3 === 0
                  ? 'rgba(34, 211, 238, 1)' // Bright cyan
                  : i % 3 === 1
                  ? 'rgba(255, 255, 255, 1)' // Pure white
                  : 'rgba(173, 216, 255, 1)', // Bright blue
                borderRadius: '50%',
                left: Math.random() * 100 + '%',
                top: Math.random() * 100 + '%',
                boxShadow: `0 0 ${Math.random() * 6 + 4}px currentColor`,
                animation: `star-twinkle ${Math.random() * 2.5 + 1.5}s ease-in-out infinite ${Math.random() * 2}s`,
                pointerEvents: 'none'
              }}
            />
          ))}

          {/* Grid overlay - subtle tech feel */}
          <div style={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundImage: `
              linear-gradient(rgba(34, 211, 238, 0.03) 1px, transparent 1px),
              linear-gradient(90deg, rgba(34, 211, 238, 0.03) 1px, transparent 1px)
            `,
            backgroundSize: '40px 40px',
            opacity: 0.3
          }} />

          {/* Nebula glow effects */}
          <div style={{
            position: 'absolute',
            top: '-30%',
            right: '15%',
            width: '500px',
            height: '500px',
            background: 'radial-gradient(circle, rgba(34, 211, 238, 0.08) 0%, rgba(59, 130, 246, 0.05) 40%, transparent 70%)',
            filter: 'blur(70px)',
            animation: 'float 7s ease-in-out infinite'
          }} />
          <div style={{
            position: 'absolute',
            bottom: '-30%',
            left: '15%',
            width: '500px',
            height: '500px',
            background: 'radial-gradient(circle, rgba(147, 51, 234, 0.06) 0%, rgba(59, 130, 246, 0.04) 40%, transparent 70%)',
            filter: 'blur(70px)',
            animation: 'float 9s ease-in-out infinite'
          }} />

          <div style={{ position: 'relative', zIndex: 1 }}>
            <h2 style={{
              fontSize: '42px',
              fontWeight: '800',
              color: 'white',
              textAlign: 'center',
              marginBottom: '16px',
              marginTop: 0
            }}>
              Why Yieldera
            </h2>
            <p style={{
              fontSize: '17px',
              color: 'rgba(255, 255, 255, 0.85)',
              textAlign: 'center',
              maxWidth: '650px',
              margin: '0 auto 48px',
              lineHeight: '1.6'
            }}>
              The only platform combining DeFi, RWA, and institutional AI in one interface
            </p>

            <div style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
              gap: '24px'
            }}>
              {[
                {
                  icon: <Globe size={32} />,
                  title: 'Dual-Chain Architecture',
                  description: 'Arbitrum for high-yield DeFi. Base for stable US Treasuries. Best of both worlds.'
                },
                {
                  icon: <Shield size={32} />,
                  title: 'AI Risk Engine',
                  description: '10,000-agent Monte Carlo simulations. VaR/CVaR metrics. Real-time liquidation alerts.'
                },
                {
                  icon: <TrendingUp size={32} />,
                  title: 'Auto-Optimization',
                  description: '24/7 rebalancing across 6+ protocols. Intelligent batching. Ultra-low gas costs.'
                },
                {
                  icon: <Gem size={32} />,
                  title: 'Treasury Bonds',
                  description: 'Tokenized US Treasuries on Base. 5% APY. Multiple maturities (1M-12M).'
                },
                {
                  icon: <Lock size={32} />,
                  title: 'Enterprise Security',
                  description: 'Multi-sig wallets. Audited contracts. Role-based access control.'
                },
                {
                  icon: <BarChart3 size={32} />,
                  title: 'Pro Analytics',
                  description: 'Liquidation probability. Sharpe ratios. Performance tracking. CSV/PDF export.'
                },
                {
                  icon: <Target size={32} />,
                  title: 'Aerodrome Integration',
                  description: 'Efficient Base chain liquidity. Minimal slippage for treasury purchases.'
                },
                {
                  icon: <Zap size={32} />,
                  title: 'L2 Speed & Cost',
                  description: '90%+ gas savings vs Ethereum. Lightning-fast transactions on Arbitrum & Base.'
                }
              ].map((item, i) => (
                <div key={i} style={{
                  background: 'rgba(255, 255, 255, 0.1)',
                  backdropFilter: 'blur(10px)',
                  border: '1px solid rgba(255, 255, 255, 0.2)',
                  borderRadius: 12,
                  padding: '28px',
                  color: 'white'
                }}>
                  <div style={{
                    width: '56px',
                    height: '56px',
                    borderRadius: '12px',
                    background: 'rgba(34, 211, 238, 0.2)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    marginBottom: '16px',
                    border: '1px solid rgba(34, 211, 238, 0.4)'
                  }}>
                    {item.icon}
                  </div>
                  <h3 style={{
                    fontSize: '18px',
                    fontWeight: '700',
                    marginTop: 0,
                    marginBottom: '8px'
                  }}>
                    {item.title}
                  </h3>
                  <p style={{
                    fontSize: '14px',
                    lineHeight: '1.6',
                    opacity: 0.9,
                    margin: 0
                  }}>
                    {item.description}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Core Features - Enhanced cards */}
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))',
          gap: '24px',
          marginTop: '64px'
        }}>
          <div className="feature-card" style={{
            padding: '36px',
            background: 'white',
            border: '1px solid rgb(226, 232, 240)',
            borderRadius: 12
          }}>
            <div style={{
              width: '64px',
              height: '64px',
              borderRadius: '16px',
              background: 'linear-gradient(135deg, rgb(59, 130, 246), rgb(37, 99, 235))',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginBottom: '24px',
              boxShadow: '0 8px 16px rgba(59, 130, 246, 0.3)'
            }}>
              <TrendingUp size={32} color="#fff" strokeWidth={2.5} />
            </div>
            <h3 style={{marginTop: 0, marginBottom: '12px', fontSize: '22px', fontWeight: '700', color: 'rgb(15, 23, 42)'}}>
              Arbitrum DeFi
            </h3>
            <p style={{fontSize: '15px', lineHeight: '1.7', marginBottom: '20px', color: 'rgb(71, 85, 105)'}}>
              High-yield strategies across Aave V3, Compound V3, Uniswap V3, and GMX V2. Auto-rebalancing every 24 hours.
            </p>
            <div style={{display: 'flex', alignItems: 'center', gap: '8px', color: 'rgb(59, 130, 246)', fontSize: '14px', fontWeight: '600'}}>
              <CheckCircle size={18} />
              6+ Protocols Integrated
            </div>
          </div>

          <div className="feature-card" style={{
            padding: '36px',
            background: 'white',
            border: '1px solid rgb(226, 232, 240)',
            borderRadius: 12
          }}>
            <div style={{
              width: '64px',
              height: '64px',
              borderRadius: '16px',
              background: 'linear-gradient(135deg, rgb(16, 185, 129), rgb(5, 150, 105))',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginBottom: '24px',
              boxShadow: '0 8px 16px rgba(16, 185, 129, 0.3)'
            }}>
              <Gem size={32} color="#fff" strokeWidth={2.5} />
            </div>
            <h3 style={{marginTop: 0, marginBottom: '12px', fontSize: '22px', fontWeight: '700', color: 'rgb(15, 23, 42)'}}>
              Base Treasury
            </h3>
            <p style={{fontSize: '15px', lineHeight: '1.7', marginBottom: '20px', color: 'rgb(71, 85, 105)'}}>
              Tokenized US Treasury bonds. Stable 5% APY backed by U.S. government. Choose from 1M, 3M, 6M, or 12M maturities.
            </p>
            <div style={{display: 'flex', alignItems: 'center', gap: '8px', color: 'rgb(16, 185, 129)', fontSize: '14px', fontWeight: '600'}}>
              <CheckCircle size={18} />
              Government-Backed RWA
            </div>
          </div>

          <div className="feature-card" style={{
            padding: '36px',
            background: 'white',
            border: '1px solid rgb(226, 232, 240)',
            borderRadius: 12
          }}>
            <div style={{
              width: '64px',
              height: '64px',
              borderRadius: '16px',
              background: 'linear-gradient(135deg, rgb(245, 158, 11), rgb(217, 119, 6))',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginBottom: '24px',
              boxShadow: '0 8px 16px rgba(245, 158, 11, 0.3)'
            }}>
              <Shield size={32} color="#fff" strokeWidth={2.5} />
            </div>
            <h3 style={{marginTop: 0, marginBottom: '12px', fontSize: '22px', fontWeight: '700', color: 'rgb(15, 23, 42)'}}>
              AI Risk Engine
            </h3>
            <p style={{fontSize: '15px', lineHeight: '1.7', marginBottom: '20px', color: 'rgb(71, 85, 105)'}}>
              10,000-agent simulations. VaR/CVaR calculations. ML predictions. Real-time liquidation alerts protect your positions.
            </p>
            <div style={{display: 'flex', alignItems: 'center', gap: '8px', color: 'rgb(245, 158, 11)', fontSize: '14px', fontWeight: '600'}}>
              <CheckCircle size={18} />
              5 AI Functions (FastAPI)
            </div>
          </div>
        </div>

        {/* Protocol Integration */}
        <div style={{
          marginTop: '80px',
          padding: '56px 40px',
          background: 'white',
          border: '1px solid rgb(226, 232, 240)',
          borderRadius: 16
        }}>
          <h2 style={{
            marginTop: 0,
            marginBottom: '48px',
            fontSize: '36px',
            fontWeight: '800',
            textAlign: 'center',
            color: 'rgb(15, 23, 42)'
          }}>
            Integrated Protocols
          </h2>

          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
            gap: '20px'
          }}>
            {[
              { name: 'Aave V3', chain: 'Arbitrum', apy: '7.2%', tvl: '$12.3B', color: 'rgb(182, 80, 158)' },
              { name: 'Compound V3', chain: 'Arbitrum', apy: '6.8%', tvl: '$3.1B', color: 'rgb(0, 211, 149)' },
              { name: 'Uniswap V3', chain: 'Arbitrum', apy: '12.5%', tvl: '$4.2B', color: 'rgb(255, 0, 122)' },
              { name: 'GMX V2', chain: 'Arbitrum', apy: '15.3%', tvl: '$582M', color: 'rgb(59, 130, 246)' },
              { name: 'Aerodrome', chain: 'Base', apy: '8.7%', tvl: '$450M', color: 'rgb(168, 85, 247)' },
              { name: 'US Treasury', chain: 'Base', apy: '5.0%', tvl: 'RWA', color: 'rgb(16, 185, 129)' }
            ].map((protocol) => (
              <div key={protocol.name} className="protocol-card" style={{
                padding: '24px',
                background: 'rgb(248, 250, 252)',
                borderRadius: '12px',
                border: '1px solid rgb(226, 232, 240)',
                textAlign: 'center'
              }}>
                <div style={{
                  width: '52px',
                  height: '52px',
                  borderRadius: '12px',
                  background: protocol.color,
                  margin: '0 auto 16px',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: '22px',
                  fontWeight: 'bold',
                  color: '#fff',
                  boxShadow: `0 8px 16px ${protocol.color}40`
                }}>
                  {protocol.name[0]}
                </div>
                <div style={{fontSize: '18px', fontWeight: '700', marginBottom: '6px', color: 'rgb(15, 23, 42)'}}>
                  {protocol.name}
                </div>
                <div style={{fontSize: '12px', color: 'rgb(100, 116, 139)', marginBottom: '10px', fontWeight: '600'}}>
                  {protocol.chain}
                </div>
                <div style={{fontSize: '15px', color: 'rgb(22, 163, 74)', fontWeight: '700', marginBottom: '4px'}}>
                  APY: {protocol.apy}
                </div>
                <div style={{fontSize: '12px', color: 'rgb(100, 116, 139)', fontWeight: '500'}}>
                  TVL: {protocol.tvl}
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* CTA Section with cyan accent */}
        <div style={{
          marginTop: '80px',
          padding: '64px 40px',
          background: 'radial-gradient(ellipse at center, rgb(10, 25, 47) 0%, rgb(5, 10, 20) 50%, rgb(0, 0, 0) 100%)',
          border: '2px solid rgba(34, 211, 238, 0.3)',
          borderRadius: 16,
          textAlign: 'center',
          position: 'relative',
          overflow: 'hidden'
        }}>
          {/* Far distant stars */}
          {[...Array(50)].map((_, i) => (
            <div
              key={`star-cta-far-${i}`}
              style={{
                position: 'absolute',
                width: Math.random() * 1.5 + 0.5 + 'px',
                height: Math.random() * 1.5 + 0.5 + 'px',
                background: 'rgba(255, 255, 255, 0.4)',
                borderRadius: '50%',
                left: Math.random() * 100 + '%',
                top: Math.random() * 100 + '%',
                boxShadow: `0 0 ${Math.random() * 2 + 1}px rgba(255, 255, 255, 0.3)`,
                animation: `star-twinkle ${Math.random() * 4 + 3}s ease-in-out infinite ${Math.random() * 4}s`,
                pointerEvents: 'none'
              }}
            />
          ))}

          {/* Medium bright stars with colors */}
          {[...Array(30)].map((_, i) => (
            <div
              key={`star-cta-mid-${i}`}
              style={{
                position: 'absolute',
                width: Math.random() * 2.5 + 1 + 'px',
                height: Math.random() * 2.5 + 1 + 'px',
                background: i % 5 === 0
                  ? 'rgba(34, 211, 238, 0.9)' // Cyan
                  : i % 5 === 1
                  ? 'rgba(173, 216, 255, 0.8)' // Blue
                  : i % 5 === 2
                  ? 'rgba(147, 197, 253, 0.7)' // Light blue
                  : 'rgba(255, 255, 255, 0.7)',
                borderRadius: '50%',
                left: Math.random() * 100 + '%',
                top: Math.random() * 100 + '%',
                boxShadow: `0 0 ${Math.random() * 4 + 2}px currentColor`,
                animation: `star-twinkle ${Math.random() * 3 + 2}s ease-in-out infinite ${Math.random() * 3}s`,
                pointerEvents: 'none'
              }}
            />
          ))}

          {/* Close prominent stars */}
          {[...Array(12)].map((_, i) => (
            <div
              key={`star-cta-close-${i}`}
              style={{
                position: 'absolute',
                width: Math.random() * 3 + 2 + 'px',
                height: Math.random() * 3 + 2 + 'px',
                background: i % 3 === 0
                  ? 'rgba(34, 211, 238, 1)' // Bright cyan
                  : i % 3 === 1
                  ? 'rgba(255, 255, 255, 1)' // Pure white
                  : 'rgba(173, 216, 255, 1)', // Bright blue
                borderRadius: '50%',
                left: Math.random() * 100 + '%',
                top: Math.random() * 100 + '%',
                boxShadow: `0 0 ${Math.random() * 6 + 4}px currentColor`,
                animation: `star-twinkle ${Math.random() * 2.5 + 1.5}s ease-in-out infinite ${Math.random() * 2}s`,
                pointerEvents: 'none'
              }}
            />
          ))}

          {/* Central nebula glow */}
          <div style={{
            position: 'absolute',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
            width: '600px',
            height: '600px',
            background: 'radial-gradient(circle, rgba(34, 211, 238, 0.12) 0%, rgba(59, 130, 246, 0.08) 40%, transparent 70%)',
            filter: 'blur(70px)'
          }} />

          {/* Additional nebula glows */}
          <div style={{
            position: 'absolute',
            top: '-20%',
            right: '10%',
            width: '400px',
            height: '400px',
            background: 'radial-gradient(circle, rgba(147, 51, 234, 0.08) 0%, transparent 70%)',
            filter: 'blur(60px)',
            animation: 'float 8s ease-in-out infinite'
          }} />
          <div style={{
            position: 'absolute',
            bottom: '-20%',
            left: '10%',
            width: '400px',
            height: '400px',
            background: 'radial-gradient(circle, rgba(59, 130, 246, 0.1) 0%, transparent 70%)',
            filter: 'blur(60px)',
            animation: 'float 6s ease-in-out infinite'
          }} />

          <div style={{ position: 'relative', zIndex: 1 }}>
            <h2 style={{
              marginTop: 0,
              fontSize: '40px',
              fontWeight: '800',
              marginBottom: '36px',
              color: 'white'
            }}>
              Start Earning Now
            </h2>
            <Link
              to="/vault"
              style={{
                display: 'inline-flex',
                alignItems: 'center',
                gap: '12px',
                padding: '20px 56px',
                background: 'rgb(34, 211, 238)',
                color: 'rgb(15, 23, 42)',
                borderRadius: '12px',
                fontWeight: '700',
                fontSize: '20px',
                textDecoration: 'none',
                boxShadow: '0 8px 32px rgba(34, 211, 238, 0.5)',
                transition: 'all 0.3s',
                border: 'none'
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.transform = 'translateY(-3px) scale(1.02)';
                e.currentTarget.style.boxShadow = '0 16px 48px rgba(34, 211, 238, 0.6)';
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.transform = 'translateY(0) scale(1)';
                e.currentTarget.style.boxShadow = '0 8px 32px rgba(34, 211, 238, 0.5)';
              }}
            >
              Launch App <ArrowRight size={24} />
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
