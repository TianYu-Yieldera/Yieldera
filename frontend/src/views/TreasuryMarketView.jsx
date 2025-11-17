import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { TrendingUp, DollarSign, Activity, BarChart3, Shield, Clock, Target, ExternalLink, Wallet } from 'lucide-react';
import TechContainer from "../components/ui/TechContainer";
import TechHeader from "../components/ui/TechHeader";
import TechCard from "../components/ui/TechCard";
import TechButton from "../components/ui/TechButton";
import { useBaseSmartWallet } from "../web3/BaseSmartWalletProvider";
import BaseLoginCard from "../components/BaseLoginCard";

export default function TreasuryMarketView() {
  const navigate = useNavigate();
  const { isConnected, address, chain } = useBaseSmartWallet();
  const [assets, setAssets] = useState([]);
  const [filteredAssets, setFilteredAssets] = useState([]);
  const [stats, setStats] = useState(null);
  const [selectedType, setSelectedType] = useState('ALL');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchData();
  }, []);

  useEffect(() => {
    if (selectedType === 'ALL') {
      setFilteredAssets(assets);
    } else {
      setFilteredAssets(assets.filter(asset => asset.treasury_type === selectedType));
    }
  }, [selectedType, assets]);

  const fetchData = async () => {
    try {
      setLoading(true);

      const mockAssets = [
        {
          asset_id: 1,
          asset_code: 'T-BILL-3M',
          treasury_type: 'T-BILL',
          maturity_date: '2025-02-15',
          coupon_rate: 0.0435,
          current_price: 99.25,
          par_value: 100,
          apy: 0.0450,
          total_supply: 10000000,
          available_supply: 7500000,
          min_investment: 1000,
          duration: '3 Months'
        },
        {
          asset_id: 2,
          asset_code: 'T-BILL-6M',
          treasury_type: 'T-BILL',
          maturity_date: '2025-05-15',
          coupon_rate: 0.0465,
          current_price: 98.50,
          par_value: 100,
          apy: 0.0480,
          total_supply: 15000000,
          available_supply: 12000000,
          min_investment: 1000,
          duration: '6 Months'
        },
        {
          asset_id: 3,
          asset_code: 'T-NOTE-2Y',
          treasury_type: 'T-NOTE',
          maturity_date: '2027-11-15',
          coupon_rate: 0.0485,
          current_price: 97.80,
          par_value: 100,
          apy: 0.0500,
          total_supply: 25000000,
          available_supply: 18000000,
          min_investment: 5000,
          duration: '2 Years'
        },
        {
          asset_id: 4,
          asset_code: 'T-NOTE-5Y',
          treasury_type: 'T-NOTE',
          maturity_date: '2030-11-15',
          coupon_rate: 0.0515,
          current_price: 96.50,
          par_value: 100,
          apy: 0.0535,
          total_supply: 35000000,
          available_supply: 25000000,
          min_investment: 5000,
          duration: '5 Years'
        },
        {
          asset_id: 5,
          asset_code: 'T-BOND-10Y',
          treasury_type: 'T-BOND',
          maturity_date: '2035-11-15',
          coupon_rate: 0.0540,
          current_price: 95.20,
          par_value: 100,
          apy: 0.0565,
          total_supply: 50000000,
          available_supply: 35000000,
          min_investment: 10000,
          duration: '10 Years'
        },
        {
          asset_id: 6,
          asset_code: 'T-BOND-30Y',
          treasury_type: 'T-BOND',
          maturity_date: '2055-11-15',
          coupon_rate: 0.0560,
          current_price: 94.00,
          par_value: 100,
          apy: 0.0590,
          total_supply: 75000000,
          available_supply: 50000000,
          min_investment: 10000,
          duration: '30 Years'
        }
      ];

      const mockStats = {
        total_tvl: 210000000,
        avg_apy: 0.0520,
        total_volume_24h: 5240000,
        total_holders: 15234
      };

      setAssets(mockAssets);
      setFilteredAssets(mockAssets);
      setStats(mockStats);
    } catch (err) {
      console.error('Error loading data:', err);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (value) => {
    if (value >= 1000000) {
      return `$${(value / 1000000).toFixed(2)}M`;
    }
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(value);
  };

  const formatPercent = (value) => {
    return `${(parseFloat(value) * 100).toFixed(2)}%`;
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const getTypeStyles = (type) => {
    switch (type) {
      case 'T-BILL':
        return {
          bg: 'rgba(34, 211, 238, 0.2)',
          text: 'rgb(34, 211, 238)',
          border: 'rgba(34, 211, 238, 0.4)'
        };
      case 'T-NOTE':
        return {
          bg: 'rgba(59, 130, 246, 0.2)',
          text: 'rgb(59, 130, 246)',
          border: 'rgba(59, 130, 246, 0.4)'
        };
      case 'T-BOND':
        return {
          bg: 'rgba(251, 191, 36, 0.2)',
          text: 'rgb(251, 191, 36)',
          border: 'rgba(251, 191, 36, 0.4)'
        };
      default:
        return {
          bg: 'rgba(100, 116, 139, 0.2)',
          text: 'rgba(203, 213, 225, 0.8)',
          border: 'rgba(100, 116, 139, 0.4)'
        };
    }
  };

  // Check if user is connected with Smart Wallet
  if (!isConnected) {
    return (
      <TechContainer>
        <div style={{ maxWidth: 600, margin: '40px auto', padding: 20 }}>
          {/* Chain Badge */}
          <div style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: 8,
            padding: '8px 16px',
            background: 'rgba(99, 102, 241, 0.15)',
            border: '1px solid rgba(99, 102, 241, 0.3)',
            borderRadius: 20,
            marginBottom: 24,
            fontSize: 13,
            fontWeight: 600,
            color: 'rgb(99, 102, 241)'
          }}>
            <Shield size={16} />
            Base Chain Exclusive
          </div>

          <BaseLoginCard onLoginSuccess={() => window.location.reload()} />

          {/* Info Box */}
          <div style={{
            marginTop: 32,
            padding: 20,
            background: 'rgba(34, 211, 238, 0.05)',
            border: '1px solid rgba(34, 211, 238, 0.2)',
            borderRadius: 12,
            fontSize: 14,
            color: 'rgba(203, 213, 225, 0.9)',
            lineHeight: 1.6
          }}>
            <strong style={{ display: 'block', marginBottom: 8, color: 'rgb(34, 211, 238)' }}>
              💡 为什么只能在 Base 链购买国债？
            </strong>
            <ul style={{ margin: '8px 0', paddingLeft: 20 }}>
              <li>Base 链提供 Smart Wallet，无需助记词</li>
              <li>Gas 费由平台赞助，降低使用门槛</li>
              <li>支持信用卡入金，更适合传统投资者</li>
              <li>美国国债属于稳健资产，匹配 Base 链的定位</li>
            </ul>
          </div>
        </div>
      </TechContainer>
    );
  }

  if (loading) {
    return (
      <TechContainer>
        <div style={{ textAlign: 'center', padding: '100px 0' }}>
          <div style={{
            width: 48,
            height: 48,
            margin: '0 auto',
            border: '3px solid rgb(34, 211, 238)',
            borderTopColor: 'transparent',
            borderRadius: '50%',
            animation: 'spin 1s linear infinite'
          }}></div>
          <p style={{ marginTop: 16, color: 'rgba(203, 213, 225, 0.8)', fontWeight: 500 }}>
            Loading Treasury Market...
          </p>
        </div>
      </TechContainer>
    );
  }

  return (
    <TechContainer>
      <TechHeader
        icon={Shield}
        title="US Treasury Securities"
        subtitle="Investment-grade tokenized US government bonds"
      >
        <TechButton
          variant="secondary"
          icon={Wallet}
          onClick={() => navigate('/treasury/holdings')}
        >
          My Holdings
        </TechButton>
      </TechHeader>

      {/* Stats Grid */}
      {stats && (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: 16, marginBottom: 32 }}>
          <TechCard
            icon={DollarSign}
            title="Total Value Locked"
            value={formatCurrency(stats.total_tvl)}
            subtitle="Across all maturities"
            iconColor="rgb(34, 211, 238)"
          />
          <TechCard
            icon={TrendingUp}
            title="Average APY"
            value={formatPercent(stats.avg_apy)}
            subtitle="Weighted average yield"
            iconColor="rgb(34, 197, 94)"
          />
          <TechCard
            icon={BarChart3}
            title="24h Volume"
            value={formatCurrency(stats.total_volume_24h)}
            subtitle="Trading volume"
            iconColor="rgb(59, 130, 246)"
          />
          <TechCard
            icon={Activity}
            title="Total Holders"
            value={stats.total_holders.toLocaleString()}
            subtitle="Active investors"
            iconColor="rgb(251, 191, 36)"
          />
        </div>
      )}

      {/* Filter Tabs */}
      <TechCard style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap', alignItems: 'center', justifyContent: 'space-between' }}>
          <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
            {['ALL', 'T-BILL', 'T-NOTE', 'T-BOND'].map((type) => (
            <button
              key={type}
              onClick={() => setSelectedType(type)}
              style={{
                padding: '10px 24px',
                border: selectedType === type
                  ? '1px solid rgba(34, 211, 238, 0.5)'
                  : '1px solid rgba(100, 116, 139, 0.3)',
                borderRadius: 8,
                background: selectedType === type
                  ? 'rgba(34, 211, 238, 0.15)'
                  : 'rgba(15, 23, 42, 0.5)',
                color: selectedType === type
                  ? 'rgb(34, 211, 238)'
                  : 'rgba(203, 213, 225, 0.7)',
                fontSize: 14,
                fontWeight: 600,
                cursor: 'pointer',
                transition: 'all 0.3s ease',
                textTransform: 'uppercase',
                letterSpacing: 0.5
              }}
              onMouseEnter={(e) => {
                if (selectedType !== type) {
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.4)';
                  e.currentTarget.style.background = 'rgba(34, 211, 238, 0.1)';
                }
              }}
              onMouseLeave={(e) => {
                if (selectedType !== type) {
                  e.currentTarget.style.borderColor = 'rgba(100, 116, 139, 0.3)';
                  e.currentTarget.style.background = 'rgba(15, 23, 42, 0.5)';
                }
              }}
            >
              {type === 'ALL' ? 'All Securities' : type}
            </button>
          ))}
          </div>

          <div
            onClick={() => navigate('/treasury/holdings')}
            style={{
              padding: '10px 20px',
              background: 'linear-gradient(135deg, rgba(34, 211, 238, 0.15) 0%, rgba(59, 130, 246, 0.15) 100%)',
              border: '1px solid rgba(34, 211, 238, 0.4)',
              borderRadius: 8,
              color: 'rgb(34, 211, 238)',
              fontSize: 14,
              fontWeight: 600,
              cursor: 'pointer',
              transition: 'all 0.3s ease',
              display: 'flex',
              alignItems: 'center',
              gap: 8,
              whiteSpace: 'nowrap'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = 'linear-gradient(135deg, rgba(34, 211, 238, 0.25) 0%, rgba(59, 130, 246, 0.25) 100%)';
              e.currentTarget.style.transform = 'translateY(-2px)';
              e.currentTarget.style.boxShadow = '0 4px 12px rgba(34, 211, 238, 0.3)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = 'linear-gradient(135deg, rgba(34, 211, 238, 0.15) 0%, rgba(59, 130, 246, 0.15) 100%)';
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = 'none';
            }}
          >
            <Wallet size={16} />
            View My Holdings
          </div>
        </div>
      </TechCard>

      {/* Assets List */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
        {filteredAssets.map((asset) => {
          const typeStyles = getTypeStyles(asset.treasury_type);
          const availabilityPercent = (asset.available_supply / asset.total_supply) * 100;

          return (
            <Link
              key={asset.asset_id}
              to={`/treasury/${asset.asset_id}`}
              style={{ textDecoration: 'none' }}
            >
              <div style={{
                background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
                borderRadius: 12,
                padding: 24,
                border: '1px solid rgba(34, 211, 238, 0.2)',
                transition: 'all 0.3s ease',
                cursor: 'pointer',
                position: 'relative',
                overflow: 'hidden'
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.5)';
                e.currentTarget.style.transform = 'translateY(-2px)';
                e.currentTarget.style.boxShadow = '0 8px 24px rgba(0,0,0,0.3)';
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.2)';
                e.currentTarget.style.transform = 'translateY(0)';
                e.currentTarget.style.boxShadow = 'none';
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
                  backgroundSize: '20px 20px',
                  opacity: 0.5
                }} />

                <div style={{ position: 'relative', zIndex: 1 }}>
                  {/* Header Row */}
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: 20, flexWrap: 'wrap', gap: 16 }}>
                    <div>
                      <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 8 }}>
                        <h3 style={{ fontSize: 20, fontWeight: 700, color: 'white', margin: 0 }}>
                          {asset.asset_code}
                        </h3>
                        <div style={{
                          padding: '4px 12px',
                          borderRadius: 12,
                          fontSize: 11,
                          fontWeight: 600,
                          textTransform: 'uppercase',
                          background: typeStyles.bg,
                          color: typeStyles.text,
                          border: `1px solid ${typeStyles.border}`
                        }}>
                          {asset.treasury_type}
                        </div>
                      </div>
                      <div style={{ display: 'flex', alignItems: 'center', gap: 8, color: 'rgba(203, 213, 225, 0.7)', fontSize: 14 }}>
                        <Clock style={{ width: 16, height: 16 }} />
                        <span>Maturity: {formatDate(asset.maturity_date)}</span>
                        <span>•</span>
                        <span>{asset.duration}</span>
                      </div>
                    </div>
                    <div style={{ textAlign: 'right' }}>
                      <div style={{ fontSize: 32, fontWeight: 700, color: 'rgb(34, 211, 238)', textShadow: '0 0 20px rgba(34, 211, 238, 0.5)' }}>
                        {formatPercent(asset.apy)}
                      </div>
                      <div style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.6)', textTransform: 'uppercase', letterSpacing: 0.5 }}>
                        APY
                      </div>
                    </div>
                  </div>

                  {/* Stats Grid */}
                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))', gap: 16, marginBottom: 16 }}>
                    <div>
                      <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', textTransform: 'uppercase', letterSpacing: 0.5, marginBottom: 4 }}>
                        Current Price
                      </div>
                      <div style={{ fontSize: 18, fontWeight: 600, color: 'white' }}>
                        ${asset.current_price.toFixed(2)}
                      </div>
                    </div>
                    <div>
                      <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', textTransform: 'uppercase', letterSpacing: 0.5, marginBottom: 4 }}>
                        Coupon Rate
                      </div>
                      <div style={{ fontSize: 18, fontWeight: 600, color: 'white' }}>
                        {formatPercent(asset.coupon_rate)}
                      </div>
                    </div>
                    <div>
                      <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', textTransform: 'uppercase', letterSpacing: 0.5, marginBottom: 4 }}>
                        Min Investment
                      </div>
                      <div style={{ fontSize: 18, fontWeight: 600, color: 'white' }}>
                        ${asset.min_investment.toLocaleString()}
                      </div>
                    </div>
                    <div>
                      <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', textTransform: 'uppercase', letterSpacing: 0.5, marginBottom: 4 }}>
                        Available Supply
                      </div>
                      <div style={{ fontSize: 18, fontWeight: 600, color: 'white' }}>
                        {formatCurrency(asset.available_supply)}
                      </div>
                    </div>
                  </div>

                  {/* Availability Bar */}
                  <div>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                      <span style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', textTransform: 'uppercase', letterSpacing: 0.5 }}>
                        Availability
                      </span>
                      <span style={{ fontSize: 12, fontWeight: 600, color: availabilityPercent > 50 ? 'rgb(34, 197, 94)' : 'rgb(251, 191, 36)' }}>
                        {availabilityPercent.toFixed(0)}% Available
                      </span>
                    </div>
                    <div style={{ width: '100%', height: 8, background: 'rgba(15, 23, 42, 0.5)', borderRadius: 4, overflow: 'hidden', border: '1px solid rgba(34, 211, 238, 0.2)' }}>
                      <div style={{
                        width: `${availabilityPercent}%`,
                        height: '100%',
                        background: availabilityPercent > 50
                          ? 'linear-gradient(90deg, rgb(34, 197, 94) 0%, rgb(74, 222, 128) 100%)'
                          : 'linear-gradient(90deg, rgb(251, 191, 36) 0%, rgb(251, 146, 60) 100%)',
                        transition: 'width 0.3s ease',
                        boxShadow: `0 0 10px ${availabilityPercent > 50 ? 'rgba(34, 197, 94, 0.5)' : 'rgba(251, 191, 36, 0.5)'}`
                      }}></div>
                    </div>
                  </div>

                  {/* View Details Link */}
                  <div style={{ marginTop: 16, display: 'flex', alignItems: 'center', gap: 6, color: 'rgb(34, 211, 238)', fontSize: 14, fontWeight: 600 }}>
                    <span>View Details</span>
                    <ExternalLink style={{ width: 16, height: 16 }} />
                  </div>
                </div>
              </div>
            </Link>
          );
        })}
      </div>

      {filteredAssets.length === 0 && (
        <TechCard>
          <div style={{ textAlign: 'center', padding: 48 }}>
            <Target style={{ width: 48, height: 48, color: 'rgba(203, 213, 225, 0.5)', margin: '0 auto 16px' }} />
            <p style={{ color: 'rgba(203, 213, 225, 0.7)', fontSize: 16 }}>
              No assets found for this category
            </p>
          </div>
        </TechCard>
      )}
    </TechContainer>
  );
}
