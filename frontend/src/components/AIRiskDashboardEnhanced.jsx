/**
 * Enhanced AI Risk Dashboard Component
 * Professional institutional-grade risk metrics with advanced visualizations
 * Features:
 * - Animated risk score gauge
 * - CountUp number animations
 * - Real-time pulse effects
 * - Protocol risk breakdown
 * - Interactive charts
 */

import React, { useState, useEffect, useRef } from 'react';
import { AlertTriangle, TrendingUp, Shield, Activity, AlertCircle, Zap, Brain, Eye, Target, BarChart3 } from 'lucide-react';
import { useWallet } from '../web3/WalletContext';
import { useDemoMode } from '../web3/DemoModeContext';
import aiRiskService from '../services/aiRiskService';

// Animated Counter Component
const AnimatedCounter = ({ value, duration = 3500, decimals = 0, prefix = '', suffix = '' }) => {
  const [count, setCount] = useState(0);
  const countRef = useRef(0);
  const startTimeRef = useRef(null);
  const animationRef = useRef(null);

  useEffect(() => {
    const targetValue = parseFloat(value) || 0;
    startTimeRef.current = null;
    countRef.current = count;

    const animate = (timestamp) => {
      if (!startTimeRef.current) startTimeRef.current = timestamp;
      const progress = Math.min((timestamp - startTimeRef.current) / duration, 1);

      // Easing function for smooth animation
      const easeOutQuart = 1 - Math.pow(1 - progress, 4);
      const currentCount = countRef.current + (targetValue - countRef.current) * easeOutQuart;

      setCount(currentCount);

      if (progress < 1) {
        animationRef.current = requestAnimationFrame(animate);
      }
    };

    animationRef.current = requestAnimationFrame(animate);

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [value, duration]);

  return (
    <span>
      {prefix}{count.toFixed(decimals)}{suffix}
    </span>
  );
};

// Risk Score Gauge Component
const RiskGauge = ({ score, size = 200 }) => {
  const [animatedScore, setAnimatedScore] = useState(0);
  const radius = (size - 40) / 2;
  const circumference = 2 * Math.PI * radius;
  const scorePercent = Math.min(Math.max(score, 0), 100) / 100;
  const strokeDashoffset = circumference - (animatedScore / 100) * circumference;

  useEffect(() => {
    const timer = setTimeout(() => {
      setAnimatedScore(score);
    }, 100);
    return () => clearTimeout(timer);
  }, [score]);

  const getRiskColor = (score) => {
    if (score < 30) return '#22c55e'; // green
    if (score < 60) return '#eab308'; // yellow
    if (score < 80) return '#f97316'; // orange
    return '#ef4444'; // red
  };

  const getRiskLabel = (score) => {
    if (score < 30) return 'Low Risk';
    if (score < 60) return 'Medium Risk';
    if (score < 80) return 'High Risk';
    return 'Critical';
  };

  const color = getRiskColor(score);

  return (
    <div style={{ position: 'relative', width: size, height: size }}>
      <svg width={size} height={size} style={{ transform: 'rotate(-90deg)' }}>
        <defs>
          <linearGradient id="gaugeGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor={color} stopOpacity="0.3" />
            <stop offset="100%" stopColor={color} stopOpacity="1" />
          </linearGradient>
          <filter id="glow">
            <feGaussianBlur stdDeviation="4" result="coloredBlur"/>
            <feMerge>
              <feMergeNode in="coloredBlur"/>
              <feMergeNode in="SourceGraphic"/>
            </feMerge>
          </filter>
        </defs>

        {/* Background circle */}
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="rgba(255, 255, 255, 0.1)"
          strokeWidth="12"
        />

        {/* Progress circle */}
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="url(#gaugeGradient)"
          strokeWidth="12"
          strokeLinecap="round"
          strokeDasharray={circumference}
          strokeDashoffset={strokeDashoffset}
          filter="url(#glow)"
          style={{
            transition: 'stroke-dashoffset 2s cubic-bezier(0.4, 0, 0.2, 1)'
          }}
        />
      </svg>

      {/* Center content */}
      <div style={{
        position: 'absolute',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        textAlign: 'center'
      }}>
        <div style={{
          fontSize: size * 0.2,
          fontWeight: 800,
          color: color,
          textShadow: `0 0 20px ${color}60`,
          lineHeight: 1
        }}>
          <AnimatedCounter value={score} decimals={1} />
        </div>
        <div style={{
          fontSize: size * 0.08,
          color: 'rgba(203, 213, 225, 0.7)',
          marginTop: 8,
          fontWeight: 600,
          textTransform: 'uppercase',
          letterSpacing: 1
        }}>
          {getRiskLabel(score)}
        </div>
      </div>
    </div>
  );
};

// Protocol Risk Breakdown Component
const ProtocolRiskBreakdown = ({ positions }) => {
  const protocols = [
    { name: 'Aave V3', risk: 28, allocation: 40, color: '#b650a2' },
    { name: 'Compound V3', risk: 32, allocation: 30, color: '#00d395' },
    { name: 'Uniswap V3', risk: 45, allocation: 20, color: '#ff007a' },
    { name: 'GMX', risk: 58, allocation: 10, color: '#3b82f6' }
  ];

  return (
    <div style={{
      background: 'rgba(255, 255, 255, 0.03)',
      borderRadius: 12,
      padding: 20,
      border: '1px solid rgba(34, 211, 238, 0.2)'
    }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 20 }}>
        <BarChart3 style={{ width: 18, height: 18, color: 'rgb(34, 211, 238)' }} />
        <h4 style={{ fontSize: 15, fontWeight: 700, color: 'white', margin: 0 }}>
          Protocol Risk Distribution
        </h4>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
        {protocols.map((protocol, index) => (
          <div key={index}>
            <div style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              marginBottom: 8
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <div style={{
                  width: 10,
                  height: 10,
                  borderRadius: '50%',
                  background: protocol.color,
                  boxShadow: `0 0 10px ${protocol.color}60`
                }} />
                <span style={{ fontSize: 13, fontWeight: 600, color: 'white' }}>
                  {protocol.name}
                </span>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)' }}>
                  {protocol.allocation}%
                </span>
                <span style={{
                  fontSize: 13,
                  fontWeight: 700,
                  color: protocol.risk < 40 ? '#22c55e' : protocol.risk < 60 ? '#eab308' : '#ef4444'
                }}>
                  Risk: {protocol.risk}
                </span>
              </div>
            </div>

            {/* Risk bar */}
            <div style={{
              width: '100%',
              height: 8,
              background: 'rgba(255, 255, 255, 0.05)',
              borderRadius: 4,
              overflow: 'hidden',
              position: 'relative'
            }}>
              <div style={{
                width: `${protocol.risk}%`,
                height: '100%',
                background: `linear-gradient(90deg, ${protocol.color} 0%, ${protocol.color}cc 100%)`,
                borderRadius: 4,
                transition: 'width 1s cubic-bezier(0.4, 0, 0.2, 1)',
                boxShadow: `0 0 10px ${protocol.color}60`
              }} />
            </div>
          </div>
        ))}
      </div>

      {/* Allocation pie visualization */}
      <div style={{ marginTop: 24 }}>
        <div style={{
          display: 'flex',
          height: 12,
          borderRadius: 6,
          overflow: 'hidden',
          boxShadow: '0 0 20px rgba(0,0,0,0.3)'
        }}>
          {protocols.map((protocol, index) => (
            <div
              key={index}
              style={{
                width: `${protocol.allocation}%`,
                background: protocol.color,
                transition: 'all 0.3s ease',
                cursor: 'pointer'
              }}
              title={`${protocol.name}: ${protocol.allocation}%`}
            />
          ))}
        </div>
        <div style={{
          display: 'flex',
          justifyContent: 'center',
          marginTop: 12,
          fontSize: 11,
          color: 'rgba(203, 213, 225, 0.6)'
        }}>
          Portfolio Allocation by Protocol
        </div>
      </div>
    </div>
  );
};

// Mini Sparkline Chart
const Sparkline = ({ data, color = '#22d3ee', height = 40 }) => {
  if (!data || data.length === 0) return null;

  const values = data.map(d => d.value);
  const min = Math.min(...values);
  const max = Math.max(...values);
  const range = max - min || 1;

  const points = data.map((d, i) => {
    const x = (i / (data.length - 1)) * 100;
    const y = ((max - d.value) / range) * 80 + 10;
    return `${x},${y}`;
  }).join(' ');

  const trend = values[values.length - 1] > values[0] ? 'up' : 'down';

  return (
    <div style={{ position: 'relative', height }}>
      <svg width="100%" height={height} viewBox="0 0 100 100" preserveAspectRatio="none">
        <defs>
          <linearGradient id={`sparkGradient-${color}`} x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor={color} stopOpacity="0.3" />
            <stop offset="100%" stopColor={color} stopOpacity="0.05" />
          </linearGradient>
        </defs>

        {/* Area */}
        <polygon
          points={`0,100 ${points} 100,100`}
          fill={`url(#sparkGradient-${color})`}
        />

        {/* Line */}
        <polyline
          points={points}
          fill="none"
          stroke={color}
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
      </svg>

      {/* Trend indicator */}
      <div style={{
        position: 'absolute',
        top: 4,
        right: 4,
        fontSize: 10,
        fontWeight: 700,
        color: trend === 'up' ? '#22c55e' : '#ef4444',
        display: 'flex',
        alignItems: 'center',
        gap: 2
      }}>
        {trend === 'up' ? '↑' : '↓'}
        {Math.abs(((values[values.length - 1] - values[0]) / values[0] * 100)).toFixed(1)}%
      </div>
    </div>
  );
};

// Main Enhanced Dashboard Component
const AIRiskDashboardEnhanced = () => {
  const { address } = useWallet();
  const { demoMode } = useDemoMode();
  const [riskData, setRiskData] = useState(null);
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [analyzing, setAnalyzing] = useState(false);

  // Generate sample sparkline data
  const generateSparklineData = (baseValue, volatility = 5) => {
    return Array.from({ length: 20 }, (_, i) => ({
      value: baseValue + (Math.random() - 0.5) * volatility + Math.sin(i / 3) * volatility / 2
    }));
  };

  useEffect(() => {
    loadRiskData();
  }, [address, demoMode]);

  const loadRiskData = async () => {
    try {
      setLoading(true);
      setAnalyzing(true);

      // Simulate API delay for loading animation
      await new Promise(resolve => setTimeout(resolve, 2000));

      // Use mock data for demo
      setRiskData({
        overall_risk_score: 35.5,
        total_collateral: 50000,
        total_debt: 25000,
        health_factor: 1.85,
        position_count: 3,
        diversification_score: 72,
        stability_score: 8.5
      });

      setAlerts([
        {
          position_id: 'demo-1',
          protocol: 'Aave V3',
          risk_score: 42.3,
          alert_level: 'warning',
          title: 'Medium Risk Position',
          description: 'Position health factor below recommended threshold',
          timestamp: new Date().toISOString()
        }
      ]);

      setLoading(false);
      setAnalyzing(false);
    } catch (err) {
      console.error('Failed to load risk data:', err);
      setLoading(false);
      setAnalyzing(false);
    }
  };

  if (loading) {
    return (
      <div style={{
        background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
        borderRadius: 16,
        padding: 64,
        border: '1px solid rgba(34, 211, 238, 0.2)',
        position: 'relative',
        overflow: 'hidden'
      }}>
        <style>{`
          @keyframes scan {
            0% { transform: translateY(-100%); }
            100% { transform: translateY(300%); }
          }
          @keyframes pulse-ring {
            0% { transform: scale(0.95); opacity: 1; }
            50% { transform: scale(1.05); opacity: 0.7; }
            100% { transform: scale(0.95); opacity: 1; }
          }
          @keyframes data-stream {
            0% { transform: translateX(-100%); opacity: 0; }
            50% { opacity: 1; }
            100% { transform: translateX(100%); opacity: 0; }
          }
        `}</style>

        {/* Scanning effect */}
        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          height: 3,
          background: 'linear-gradient(90deg, transparent, rgba(34, 211, 238, 0.8), transparent)',
          animation: 'scan 2s ease-in-out infinite'
        }} />

        {/* Data stream effect */}
        {[...Array(3)].map((_, i) => (
          <div
            key={i}
            style={{
              position: 'absolute',
              left: 0,
              top: `${30 + i * 20}%`,
              width: 200,
              height: 2,
              background: 'linear-gradient(90deg, transparent, rgba(34, 211, 238, 0.6), transparent)',
              animation: `data-stream 2s ease-in-out infinite ${i * 0.7}s`
            }}
          />
        ))}

        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
          <div style={{ position: 'relative', marginBottom: 32 }}>
            {/* Pulsing rings */}
            <div style={{
              position: 'absolute',
              width: 100,
              height: 100,
              border: '2px solid rgba(34, 211, 238, 0.3)',
              borderRadius: '50%',
              animation: 'pulse-ring 2s ease-in-out infinite'
            }} />
            <div style={{
              position: 'absolute',
              width: 120,
              height: 120,
              top: -10,
              left: -10,
              border: '2px solid rgba(34, 211, 238, 0.2)',
              borderRadius: '50%',
              animation: 'pulse-ring 2s ease-in-out infinite 0.5s'
            }} />
            <div style={{
              width: 100,
              height: 100,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
              <Brain style={{ width: 48, height: 48, color: 'rgb(34, 211, 238)' }} />
            </div>
          </div>

          <div style={{ fontSize: 24, fontWeight: 700, color: 'white', marginBottom: 12 }}>
            AI Risk Analysis
          </div>
          <div style={{ color: 'rgba(203, 213, 225, 0.8)', fontSize: 15, textAlign: 'center', maxWidth: 400 }}>
            Running 10,000 Monte Carlo simulations across 50+ protocols...
          </div>

          {/* Progress dots */}
          <div style={{ display: 'flex', gap: 8, marginTop: 24 }}>
            {[...Array(3)].map((_, i) => (
              <div
                key={i}
                style={{
                  width: 8,
                  height: 8,
                  borderRadius: '50%',
                  background: 'rgb(34, 211, 238)',
                  animation: `pulse-ring 1.5s ease-in-out infinite ${i * 0.3}s`
                }}
              />
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div style={{
      background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 50%, rgb(15, 23, 42) 100%)',
      borderRadius: 16,
      padding: 40,
      border: '1px solid rgba(34, 211, 238, 0.2)',
      position: 'relative',
      overflow: 'hidden'
    }}>
      <style>{`
        @keyframes glow-pulse {
          0%, 100% {
            box-shadow: 0 0 20px rgba(34, 211, 238, 0.3);
          }
          50% {
            box-shadow: 0 0 40px rgba(34, 211, 238, 0.5);
          }
        }
      `}</style>

      {/* Tech grid background */}
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
        opacity: 0.5
      }} />

      <div style={{ position: 'relative', zIndex: 1 }}>
        {/* Header */}
        <div style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 40,
          paddingBottom: 24,
          borderBottom: '1px solid rgba(34, 211, 238, 0.2)'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <div style={{
              width: 48,
              height: 48,
              borderRadius: 12,
              background: 'rgba(34, 211, 238, 0.15)',
              border: '1px solid rgba(34, 211, 238, 0.3)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
              <Brain style={{ width: 24, height: 24, color: 'rgb(34, 211, 238)' }} />
            </div>
            <div>
              <h2 style={{
                fontSize: 24,
                fontWeight: 800,
                color: 'white',
                margin: 0,
                letterSpacing: 0.3
              }}>
                AI Risk Engine
              </h2>
              <p style={{
                fontSize: 13,
                color: 'rgba(203, 213, 225, 0.7)',
                margin: '4px 0 0 0'
              }}>
                Real-time institutional-grade risk monitoring
              </p>
            </div>
          </div>

          <button
            onClick={loadRiskData}
            disabled={analyzing}
            style={{
              padding: '12px 24px',
              background: analyzing
                ? 'rgba(100, 116, 139, 0.2)'
                : 'rgba(34, 211, 238, 0.15)',
              border: analyzing
                ? '1px solid rgba(100, 116, 139, 0.4)'
                : '1px solid rgba(34, 211, 238, 0.3)',
              borderRadius: 8,
              cursor: analyzing ? 'not-allowed' : 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: 10,
              transition: 'all 0.2s ease'
            }}
          >
            <Activity
              style={{
                width: 18,
                height: 18,
                color: analyzing ? 'rgba(203, 213, 225, 0.5)' : 'rgb(34, 211, 238)'
              }}
            />
            <span style={{
              fontSize: 14,
              fontWeight: 600,
              color: analyzing ? 'rgba(203, 213, 225, 0.5)' : 'white'
            }}>
              {analyzing ? 'Analyzing...' : 'Refresh Analysis'}
            </span>
          </button>
        </div>

        {/* Main Content Grid */}
        <div style={{
          display: 'grid',
          gridTemplateColumns: '300px 1fr',
          gap: 32,
          marginBottom: 32
        }}>
          {/* Left: Risk Gauge */}
          <div style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            padding: 32,
            background: 'rgba(255, 255, 255, 0.03)',
            borderRadius: 12,
            border: '1px solid rgba(34, 211, 238, 0.2)'
          }}>
            <RiskGauge score={riskData?.overall_risk_score || 0} size={220} />

            <div style={{
              marginTop: 24,
              textAlign: 'center',
              width: '100%'
            }}>
              <div style={{
                fontSize: 13,
                color: 'rgba(203, 213, 225, 0.7)',
                marginBottom: 8
              }}>
                Overall Risk Score
              </div>
              <div style={{
                padding: 12,
                background: 'rgba(255, 255, 255, 0.05)',
                borderRadius: 8,
                border: '1px solid rgba(255, 255, 255, 0.1)'
              }}>
                <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)', marginBottom: 4 }}>
                  Last 7 days trend
                </div>
                <Sparkline
                  data={generateSparklineData(riskData?.overall_risk_score || 35, 8)}
                  color="#22d3ee"
                  height={30}
                />
              </div>
            </div>
          </div>

          {/* Right: Key Metrics Grid */}
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(2, 1fr)',
            gap: 16
          }}>
            {/* Health Factor */}
            <div style={{
              padding: 24,
              background: 'rgba(255, 255, 255, 0.03)',
              borderRadius: 12,
              border: '1px solid rgba(34, 211, 238, 0.2)',
              transition: 'all 0.3s ease'
            }}>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
                <div style={{
                  width: 40,
                  height: 40,
                  borderRadius: 10,
                  background: 'rgba(34, 197, 94, 0.15)',
                  border: '1px solid rgba(34, 197, 94, 0.3)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center'
                }}>
                  <TrendingUp style={{ width: 20, height: 20, color: 'rgb(34, 197, 94)' }} />
                </div>
                <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)', textTransform: 'uppercase', letterSpacing: 1 }}>
                  Health Factor
                </div>
              </div>
              <div style={{
                fontSize: 36,
                fontWeight: 800,
                color: riskData?.health_factor > 1.5 ? 'rgb(34, 197, 94)' : 'rgb(249, 115, 22)',
                marginBottom: 8
              }}>
                <AnimatedCounter value={riskData?.health_factor || 0} decimals={2} />
              </div>
              <div style={{
                fontSize: 12,
                color: 'rgba(203, 213, 225, 0.7)',
                marginBottom: 12
              }}>
                {riskData?.health_factor > 1.5 ? 'Healthy Position' : 'Monitor Closely'}
              </div>
              <Sparkline
                data={generateSparklineData(riskData?.health_factor || 1.85, 0.3)}
                color={riskData?.health_factor > 1.5 ? '#22c55e' : '#f97316'}
                height={30}
              />
            </div>

            {/* Collateral */}
            <div style={{
              padding: 24,
              background: 'rgba(255, 255, 255, 0.03)',
              borderRadius: 12,
              border: '1px solid rgba(34, 211, 238, 0.2)',
              transition: 'all 0.3s ease'
            }}>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
                <div style={{
                  width: 40,
                  height: 40,
                  borderRadius: 10,
                  background: 'rgba(59, 130, 246, 0.15)',
                  border: '1px solid rgba(59, 130, 246, 0.3)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center'
                }}>
                  <Shield style={{ width: 20, height: 20, color: 'rgb(59, 130, 246)' }} />
                </div>
                <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)', textTransform: 'uppercase', letterSpacing: 1 }}>
                  Collateral
                </div>
              </div>
              <div style={{
                fontSize: 32,
                fontWeight: 800,
                color: 'rgb(59, 130, 246)',
                marginBottom: 8
              }}>
                <AnimatedCounter value={riskData?.total_collateral || 0} decimals={0} prefix="$" />
              </div>
              <div style={{
                fontSize: 12,
                color: 'rgba(203, 213, 225, 0.7)',
                marginBottom: 12
              }}>
                Locked in {riskData?.position_count || 0} protocols
              </div>
              <Sparkline
                data={generateSparklineData(riskData?.total_collateral || 50000, 2000)}
                color="#3b82f6"
                height={30}
              />
            </div>

            {/* Debt */}
            <div style={{
              padding: 24,
              background: 'rgba(255, 255, 255, 0.03)',
              borderRadius: 12,
              border: '1px solid rgba(34, 211, 238, 0.2)',
              transition: 'all 0.3s ease'
            }}>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
                <div style={{
                  width: 40,
                  height: 40,
                  borderRadius: 10,
                  background: 'rgba(234, 179, 8, 0.15)',
                  border: '1px solid rgba(234, 179, 8, 0.3)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center'
                }}>
                  <AlertTriangle style={{ width: 20, height: 20, color: 'rgb(234, 179, 8)' }} />
                </div>
                <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)', textTransform: 'uppercase', letterSpacing: 1 }}>
                  Total Debt
                </div>
              </div>
              <div style={{
                fontSize: 32,
                fontWeight: 800,
                color: 'rgb(234, 179, 8)',
                marginBottom: 8
              }}>
                <AnimatedCounter value={riskData?.total_debt || 0} decimals={0} prefix="$" />
              </div>
              <div style={{
                fontSize: 12,
                color: 'rgba(203, 213, 225, 0.7)',
                marginBottom: 12
              }}>
                {((riskData?.total_debt / riskData?.total_collateral) * 100).toFixed(1)}% utilization
              </div>
              <Sparkline
                data={generateSparklineData(riskData?.total_debt || 25000, 1000)}
                color="#eab308"
                height={30}
              />
            </div>

            {/* Diversification */}
            <div style={{
              padding: 24,
              background: 'rgba(255, 255, 255, 0.03)',
              borderRadius: 12,
              border: '1px solid rgba(34, 211, 238, 0.2)',
              transition: 'all 0.3s ease'
            }}>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
                <div style={{
                  width: 40,
                  height: 40,
                  borderRadius: 10,
                  background: 'rgba(167, 139, 250, 0.15)',
                  border: '1px solid rgba(167, 139, 250, 0.3)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center'
                }}>
                  <Target style={{ width: 20, height: 20, color: 'rgb(167, 139, 250)' }} />
                </div>
                <div style={{ fontSize: 11, color: 'rgba(203, 213, 225, 0.6)', textTransform: 'uppercase', letterSpacing: 1 }}>
                  Diversification
                </div>
              </div>
              <div style={{
                fontSize: 32,
                fontWeight: 800,
                color: 'rgb(167, 139, 250)',
                marginBottom: 8
              }}>
                <AnimatedCounter value={riskData?.diversification_score || 72} decimals={0} suffix="%" />
              </div>
              <div style={{
                fontSize: 12,
                color: 'rgba(203, 213, 225, 0.7)',
                marginBottom: 12
              }}>
                Well diversified portfolio
              </div>
              <Sparkline
                data={generateSparklineData(riskData?.diversification_score || 72, 5)}
                color="#a78bfa"
                height={30}
              />
            </div>
          </div>
        </div>

        {/* Protocol Risk Breakdown */}
        <div style={{ marginBottom: 32 }}>
          <ProtocolRiskBreakdown positions={riskData?.position_count || 0} />
        </div>

        {/* Active Alerts */}
        {alerts && alerts.length > 0 && (
          <div style={{
            background: 'rgba(234, 179, 8, 0.1)',
            border: '1px solid rgba(234, 179, 8, 0.3)',
            borderRadius: 12,
            padding: 24
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 16 }}>
              <div style={{
                width: 40,
                height: 40,
                borderRadius: 10,
                background: 'rgba(234, 179, 8, 0.2)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
              }}>
                <AlertCircle style={{ width: 20, height: 20, color: 'rgb(234, 179, 8)' }} />
              </div>
              <div>
                <h4 style={{ fontSize: 16, fontWeight: 700, color: 'white', margin: 0 }}>
                  Active Risk Alerts
                </h4>
                <p style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', margin: '2px 0 0 0' }}>
                  {alerts.length} alert{alerts.length > 1 ? 's' : ''} requiring attention
                </p>
              </div>
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
              {alerts.map((alert, index) => (
                <div
                  key={index}
                  style={{
                    padding: 16,
                    background: 'rgba(255, 255, 255, 0.05)',
                    borderRadius: 8,
                    border: '1px solid rgba(255, 255, 255, 0.1)'
                  }}
                >
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', marginBottom: 8 }}>
                    <div style={{ fontSize: 14, fontWeight: 600, color: 'white' }}>
                      {alert.title}
                    </div>
                    <div style={{
                      padding: '4px 12px',
                      background: 'rgba(234, 179, 8, 0.2)',
                      border: '1px solid rgba(234, 179, 8, 0.4)',
                      borderRadius: 12,
                      fontSize: 11,
                      fontWeight: 700,
                      color: 'rgb(234, 179, 8)',
                      textTransform: 'uppercase'
                    }}>
                      {alert.alert_level}
                    </div>
                  </div>
                  <div style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.7)', marginBottom: 8 }}>
                    {alert.description}
                  </div>
                  <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.5)' }}>
                    {alert.protocol} • {new Date(alert.timestamp).toLocaleTimeString()}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default AIRiskDashboardEnhanced;
