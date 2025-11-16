/**
 * Hedge Recommendation Component - Professional Financial Tech Style
 * AI-powered hedging suggestions with institutional-grade visualization
 */

import React, { useState, useEffect } from 'react';
import { Shield, TrendingDown, AlertTriangle, CheckCircle, Info, ChevronRight, Activity, Target, Zap } from 'lucide-react';
import hedgeService from '../services/hedgeService';
import { useWallet } from '../web3/WalletContext';

export default function HedgeRecommendation({ riskScore, portfolioValue, volatility }) {
  const { address } = useWallet();
  const [hedgeHistory, setHedgeHistory] = useState([]);
  const [hedgeSettings, setHedgeSettings] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (address) {
      loadHedgeData();
    } else {
      setLoading(false);
    }
  }, [address]);

  const loadHedgeData = async () => {
    try {
      setLoading(true);
      const [history, settings] = await Promise.all([
        hedgeService.getHedgeHistory(address, 5).catch(() => ({ hedges: [] })),
        hedgeService.getUserSettings(address).catch(() => ({
          auto_hedge_enabled: false,
          risk_threshold: 60,
          hedge_percentage: 30,
          hedge_strategy: 'short_position'
        }))
      ]);

      setHedgeHistory(history.hedges || []);
      setHedgeSettings(settings);
    } catch (error) {
      console.error('Failed to load hedge data:', error);
    } finally {
      setLoading(false);
    }
  };

  const recommendation = hedgeService.calculateHedgeRecommendation({
    value: portfolioValue || 50000,
    risk_score: riskScore || 35,
    volatility: volatility || 0.15
  });

  const strategyDetails = hedgeService.getStrategyDetails(recommendation.strategy);

  const getUrgencyColor = (urgency) => {
    switch(urgency) {
      case 'high': return 'rgb(239, 68, 68)';
      case 'medium': return 'rgb(234, 179, 8)';
      case 'low': return 'rgb(34, 197, 94)';
      default: return 'rgb(148, 163, 184)';
    }
  };

  if (loading) {
    return (
      <div style={{
        background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
        borderRadius: 16,
        padding: 32,
        border: '1px solid rgba(34, 211, 238, 0.2)',
        position: 'relative',
        overflow: 'hidden'
      }}>
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
        <div style={{ position: 'relative', zIndex: 1, display: 'flex', alignItems: 'center', gap: 16 }}>
          <Shield style={{ width: 32, height: 32, color: 'rgb(34, 211, 238)' }} />
          <div>
            <h3 style={{ fontSize: 20, fontWeight: 700, color: 'white', margin: 0 }}>
              Hedge Recommendations
            </h3>
            <p style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.8)', margin: '4px 0 0 0' }}>
              Analyzing risk mitigation strategies...
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
      <style>{`
        @keyframes pulse-glow {
          0%, 100% {
            box-shadow: 0 0 20px rgba(34, 211, 238, 0.2), 0 0 40px rgba(34, 211, 238, 0.1);
          }
          50% {
            box-shadow: 0 0 30px rgba(34, 211, 238, 0.3), 0 0 60px rgba(34, 211, 238, 0.15);
          }
        }
        @keyframes slide-in {
          from { opacity: 0; transform: translateX(-20px); }
          to { opacity: 1; transform: translateX(0); }
        }
        @keyframes pulse-ring {
          0% { transform: scale(0.95); opacity: 1; }
          50% { transform: scale(1.05); opacity: 0.7; }
          100% { transform: scale(0.95); opacity: 1; }
        }
      `}</style>

      {/* Main Container */}
      <div style={{
        background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 50%, rgb(15, 23, 42) 100%)',
        borderRadius: 16,
        padding: 32,
        border: '1px solid rgba(34, 211, 238, 0.2)',
        position: 'relative',
        overflow: 'hidden',
        animation: 'pulse-glow 4s ease-in-out infinite'
      }}>
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
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 32 }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
              <div style={{
                width: 56,
                height: 56,
                borderRadius: 12,
                background: 'linear-gradient(135deg, rgba(139, 92, 246, 0.2) 0%, rgba(79, 70, 229, 0.2) 100%)',
                border: '1px solid rgba(139, 92, 246, 0.3)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                position: 'relative'
              }}>
                <div style={{
                  position: 'absolute',
                  width: '100%',
                  height: '100%',
                  border: '2px solid rgba(139, 92, 246, 0.3)',
                  borderRadius: 12,
                  animation: 'pulse-ring 2s ease-in-out infinite'
                }} />
                <Shield style={{ width: 28, height: 28, color: 'rgb(167, 139, 250)' }} />
              </div>
              <div>
                <h3 style={{ fontSize: 24, fontWeight: 700, color: 'white', margin: 0, marginBottom: 4 }}>
                  Hedge Recommendations
                </h3>
                <p style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.8)', margin: 0 }}>
                  AI-powered risk mitigation strategies
                </p>
              </div>
            </div>

            {hedgeSettings && (
              <div style={{
                display: 'flex',
                alignItems: 'center',
                gap: 8,
                padding: '10px 16px',
                borderRadius: 10,
                background: hedgeSettings.auto_hedge_enabled
                  ? 'rgba(34, 197, 94, 0.15)'
                  : 'rgba(148, 163, 184, 0.15)',
                border: `1px solid ${hedgeSettings.auto_hedge_enabled ? 'rgba(34, 197, 94, 0.3)' : 'rgba(148, 163, 184, 0.3)'}`
              }}>
                <div style={{
                  width: 8,
                  height: 8,
                  borderRadius: '50%',
                  background: hedgeSettings.auto_hedge_enabled ? 'rgb(34, 197, 94)' : 'rgb(148, 163, 184)',
                  boxShadow: hedgeSettings.auto_hedge_enabled ? '0 0 10px rgb(34, 197, 94)' : 'none'
                }} />
                <span style={{
                  fontSize: 12,
                  fontWeight: 600,
                  color: hedgeSettings.auto_hedge_enabled ? 'rgb(134, 239, 172)' : 'rgba(203, 213, 225, 0.8)',
                  textTransform: 'uppercase',
                  letterSpacing: 0.5
                }}>
                  Auto-Hedge {hedgeSettings.auto_hedge_enabled ? 'Active' : 'Disabled'}
                </span>
              </div>
            )}
          </div>

          {/* Main Recommendation */}
          {recommendation.recommended ? (
            <div style={{
              background: `linear-gradient(135deg, ${getUrgencyColor(recommendation.urgency)}15 0%, ${getUrgencyColor(recommendation.urgency)}05 100%)`,
              border: `1px solid ${getUrgencyColor(recommendation.urgency)}40`,
              borderRadius: 12,
              padding: 24,
              marginBottom: 24,
              position: 'relative',
              overflow: 'hidden'
            }}>
              {/* Warning stripe animation */}
              <div style={{
                position: 'absolute',
                top: 0,
                left: 0,
                right: 0,
                height: 2,
                background: `linear-gradient(90deg, transparent, ${getUrgencyColor(recommendation.urgency)}, transparent)`,
                animation: 'slide-in 2s ease-in-out infinite'
              }} />

              <div style={{ display: 'flex', alignItems: 'flex-start', gap: 20 }}>
                <div style={{
                  width: 48,
                  height: 48,
                  borderRadius: 12,
                  background: `${getUrgencyColor(recommendation.urgency)}20`,
                  border: `1px solid ${getUrgencyColor(recommendation.urgency)}40`,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  flexShrink: 0
                }}>
                  <TrendingDown style={{ width: 24, height: 24, color: getUrgencyColor(recommendation.urgency) }} />
                </div>

                <div style={{ flex: 1 }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 12 }}>
                    <h4 style={{ fontSize: 20, fontWeight: 700, color: 'white', margin: 0 }}>
                      {strategyDetails.name}
                    </h4>
                    <span style={{
                      padding: '4px 12px',
                      borderRadius: 6,
                      fontSize: 11,
                      fontWeight: 700,
                      color: 'white',
                      background: getUrgencyColor(recommendation.urgency),
                      textTransform: 'uppercase',
                      letterSpacing: '0.5px',
                      boxShadow: `0 0 20px ${getUrgencyColor(recommendation.urgency)}40`
                    }}>
                      {recommendation.urgency} Priority
                    </span>
                  </div>

                  <p style={{
                    fontSize: 14,
                    color: 'rgba(203, 213, 225, 0.9)',
                    margin: '0 0 20px 0',
                    lineHeight: 1.6
                  }}>
                    {recommendation.reason}
                  </p>

                  {/* Metrics Grid */}
                  <div style={{
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))',
                    gap: 16,
                    marginBottom: 20
                  }}>
                    {[
                      { label: 'Hedge Amount', value: hedgeService.formatAmount(recommendation.hedge_amount), icon: Target },
                      { label: 'Percentage', value: `${recommendation.hedge_percentage}%`, icon: Activity },
                      { label: 'Est. Cost', value: hedgeService.formatAmount(recommendation.estimated_cost), icon: TrendingDown }
                    ].map((metric, i) => {
                      const Icon = metric.icon;
                      return (
                        <div key={i} style={{
                          padding: 16,
                          background: 'rgba(255, 255, 255, 0.05)',
                          borderRadius: 10,
                          border: '1px solid rgba(255, 255, 255, 0.1)'
                        }}>
                          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 8 }}>
                            <p style={{
                              fontSize: 11,
                              color: 'rgba(203, 213, 225, 0.7)',
                              margin: 0,
                              fontWeight: 600,
                              textTransform: 'uppercase',
                              letterSpacing: 0.5
                            }}>
                              {metric.label}
                            </p>
                            <Icon style={{ width: 16, height: 16, color: 'rgb(34, 211, 238)', opacity: 0.6 }} />
                          </div>
                          <p style={{ fontSize: 22, fontWeight: 700, color: 'white', margin: 0 }}>
                            {metric.value}
                          </p>
                        </div>
                      );
                    })}
                  </div>

                  {/* Action Button */}
                  <button style={{
                    padding: '14px 28px',
                    background: `linear-gradient(135deg, ${getUrgencyColor(recommendation.urgency)} 0%, ${getUrgencyColor(recommendation.urgency)}dd 100%)`,
                    color: 'white',
                    border: 'none',
                    borderRadius: 10,
                    fontSize: 14,
                    fontWeight: 600,
                    cursor: 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    gap: 10,
                    transition: 'all 0.3s ease',
                    boxShadow: `0 4px 16px ${getUrgencyColor(recommendation.urgency)}40`
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.transform = 'translateY(-2px)';
                    e.currentTarget.style.boxShadow = `0 6px 20px ${getUrgencyColor(recommendation.urgency)}60`;
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.transform = 'translateY(0)';
                    e.currentTarget.style.boxShadow = `0 4px 16px ${getUrgencyColor(recommendation.urgency)}40`;
                  }}>
                    <Zap style={{ width: 18, height: 18 }} />
                    Execute Hedge Strategy
                    <ChevronRight style={{ width: 18, height: 18 }} />
                  </button>
                </div>
              </div>
            </div>
          ) : (
            <div style={{
              background: 'linear-gradient(135deg, rgba(34, 197, 94, 0.15) 0%, rgba(34, 197, 94, 0.05) 100%)',
              border: '1px solid rgba(34, 197, 94, 0.3)',
              borderRadius: 12,
              padding: 24,
              marginBottom: 24,
              display: 'flex',
              alignItems: 'center',
              gap: 20
            }}>
              <div style={{
                width: 56,
                height: 56,
                borderRadius: 12,
                background: 'rgba(34, 197, 94, 0.2)',
                border: '1px solid rgba(34, 197, 94, 0.4)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                flexShrink: 0
              }}>
                <CheckCircle style={{ width: 28, height: 28, color: 'rgb(134, 239, 172)' }} />
              </div>
              <div>
                <h4 style={{ fontSize: 18, fontWeight: 700, color: 'white', margin: '0 0 6px 0' }}>
                  No Hedging Required
                </h4>
                <p style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.9)', margin: 0, lineHeight: 1.5 }}>
                  {recommendation.reason}
                </p>
              </div>
            </div>
          )}

          {/* Strategy Details */}
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(3, 1fr)',
            gap: 16,
            marginBottom: 24
          }}>
            {[
              { label: 'Cost', value: strategyDetails.cost, color: 'rgb(234, 179, 8)' },
              { label: 'Effectiveness', value: strategyDetails.effectiveness, color: 'rgb(34, 211, 238)' },
              { label: 'Complexity', value: strategyDetails.complexity, color: 'rgb(167, 139, 250)' }
            ].map((item, i) => (
              <div key={i} style={{
                padding: 20,
                background: 'rgba(255, 255, 255, 0.03)',
                borderRadius: 12,
                border: '1px solid rgba(255, 255, 255, 0.08)',
                textAlign: 'center',
                transition: 'all 0.3s ease'
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.background = 'rgba(255, 255, 255, 0.06)';
                e.currentTarget.style.borderColor = `${item.color}40`;
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.background = 'rgba(255, 255, 255, 0.03)';
                e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.08)';
              }}>
                <p style={{
                  fontSize: 11,
                  color: 'rgba(203, 213, 225, 0.7)',
                  margin: '0 0 8px 0',
                  fontWeight: 600,
                  textTransform: 'uppercase',
                  letterSpacing: 0.5
                }}>
                  {item.label}
                </p>
                <p style={{
                  fontSize: 18,
                  fontWeight: 700,
                  color: item.color,
                  margin: 0,
                  textShadow: `0 0 20px ${item.color}40`
                }}>
                  {item.value}
                </p>
              </div>
            ))}
          </div>

          {/* Recent Hedge History */}
          {hedgeHistory.length > 0 && (
            <div style={{
              background: 'rgba(255, 255, 255, 0.02)',
              borderRadius: 12,
              padding: 20,
              border: '1px solid rgba(255, 255, 255, 0.08)'
            }}>
              <h4 style={{ fontSize: 16, fontWeight: 700, color: 'white', marginBottom: 16, display: 'flex', alignItems: 'center', gap: 8 }}>
                <Activity style={{ width: 18, height: 18, color: 'rgb(34, 211, 238)' }} />
                Recent Hedges
              </h4>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
                {hedgeHistory.slice(0, 3).map((hedge, index) => (
                  <div key={index} style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    padding: 14,
                    background: 'rgba(255, 255, 255, 0.03)',
                    borderRadius: 10,
                    border: '1px solid rgba(255, 255, 255, 0.08)',
                    transition: 'all 0.2s ease'
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.background = 'rgba(255, 255, 255, 0.06)';
                    e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.3)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.background = 'rgba(255, 255, 255, 0.03)';
                    e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.08)';
                  }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                      <div style={{
                        width: 36,
                        height: 36,
                        borderRadius: 10,
                        background: 'rgba(34, 197, 94, 0.15)',
                        border: '1px solid rgba(34, 197, 94, 0.3)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center'
                      }}>
                        <Shield style={{ width: 16, height: 16, color: 'rgb(134, 239, 172)' }} />
                      </div>
                      <div>
                        <p style={{ fontSize: 14, fontWeight: 600, color: 'white', margin: 0 }}>
                          {hedge.strategy || 'Short Position'}
                        </p>
                        <p style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', margin: 0 }}>
                          {hedge.timestamp ? new Date(hedge.timestamp).toLocaleDateString() : '2 days ago'}
                        </p>
                      </div>
                    </div>
                    <p style={{ fontSize: 16, fontWeight: 700, color: 'rgb(34, 211, 238)', margin: 0 }}>
                      {hedgeService.formatAmount(hedge.amount || 5000)}
                    </p>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Info Footer */}
      <div style={{
        background: 'linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(37, 99, 235, 0.05) 100%)',
        border: '1px solid rgba(59, 130, 246, 0.2)',
        borderRadius: 12,
        padding: 16,
        display: 'flex',
        alignItems: 'flex-start',
        gap: 12
      }}>
        <Info style={{ width: 20, height: 20, color: 'rgb(96, 165, 250)', flexShrink: 0, marginTop: 2 }} />
        <p style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.8)', margin: 0, lineHeight: 1.6 }}>
          Hedge recommendations are AI-powered suggestions based on your current risk profile and market conditions.
          Always conduct your own research before implementing any hedging strategy.
        </p>
      </div>
    </div>
  );
}
