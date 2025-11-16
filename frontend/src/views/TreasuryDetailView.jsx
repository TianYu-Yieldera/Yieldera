import React, { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import {
  ArrowLeft,
  TrendingUp,
  Users,
  Calendar,
  ShoppingCart,
  Activity,
  DollarSign,
  BarChart3,
  Clock
} from 'lucide-react';
import TechContainer from "../components/ui/TechContainer";
import TechCard from "../components/ui/TechCard";
import TechButton from "../components/ui/TechButton";

export default function TreasuryDetailView() {
  const { assetId } = useParams();
  const navigate = useNavigate();
  const [asset, setAsset] = useState(null);
  const [activeTab, setActiveTab] = useState('info');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchAssetData();
  }, [assetId]);

  const fetchAssetData = async () => {
    try {
      setLoading(true);

      // Mock data as fallback
      const mockAssets = {
        '1': {
          asset_id: 1,
          asset_code: 'T-BILL-3M',
          treasury_type: 'T-BILL',
          maturity_date: '2025-02-15',
          maturity_term: '3 Months',
          current_price: 99.25,
          face_value: 100,
          current_yield: 0.0450,
          yield_rate: 0.0450,
          coupon_rate: 0.0435,
          total_supply: 10000000,
          available_supply: 7500000,
          total_holders: 156,
          volume_24h: 185000,
          cusip: 'US912796XX12'
        },
        '2': {
          asset_id: 2,
          asset_code: 'T-BILL-6M',
          treasury_type: 'T-BILL',
          maturity_date: '2025-05-15',
          maturity_term: '6 Months',
          current_price: 98.50,
          face_value: 100,
          current_yield: 0.0480,
          yield_rate: 0.0480,
          coupon_rate: 0.0465,
          total_supply: 15000000,
          available_supply: 12000000,
          total_holders: 234,
          volume_24h: 245000,
          cusip: 'US912796YY23'
        },
        '3': {
          asset_id: 3,
          asset_code: 'T-NOTE-2Y',
          treasury_type: 'T-NOTE',
          maturity_date: '2027-11-15',
          maturity_term: '2 Years',
          current_price: 97.80,
          face_value: 100,
          current_yield: 0.0500,
          yield_rate: 0.0500,
          coupon_rate: 0.0485,
          total_supply: 25000000,
          available_supply: 18000000,
          total_holders: 412,
          volume_24h: 520000,
          cusip: 'US912828ZZ34'
        },
        '4': {
          asset_id: 4,
          asset_code: 'T-NOTE-5Y',
          treasury_type: 'T-NOTE',
          maturity_date: '2030-11-15',
          maturity_term: '5 Years',
          current_price: 96.50,
          face_value: 100,
          current_yield: 0.0535,
          yield_rate: 0.0535,
          coupon_rate: 0.0515,
          total_supply: 35000000,
          available_supply: 25000000,
          total_holders: 689,
          volume_24h: 780000,
          cusip: 'US912828AA45'
        },
        '5': {
          asset_id: 5,
          asset_code: 'T-BOND-10Y',
          treasury_type: 'T-BOND',
          maturity_date: '2035-11-15',
          maturity_term: '10 Years',
          current_price: 95.20,
          face_value: 100,
          current_yield: 0.0565,
          yield_rate: 0.0565,
          coupon_rate: 0.0540,
          total_supply: 50000000,
          available_supply: 35000000,
          total_holders: 1245,
          volume_24h: 1200000,
          cusip: 'US912810BB56'
        },
        '6': {
          asset_id: 6,
          asset_code: 'T-BOND-30Y',
          treasury_type: 'T-BOND',
          maturity_date: '2055-11-15',
          maturity_term: '30 Years',
          current_price: 94.00,
          face_value: 100,
          current_yield: 0.0590,
          yield_rate: 0.0590,
          coupon_rate: 0.0560,
          total_supply: 75000000,
          available_supply: 50000000,
          total_holders: 2134,
          volume_24h: 1850000,
          cusip: 'US912810CC67'
        }
      };

      const mockAsset = mockAssets[assetId];

      if (mockAsset) {
        setAsset(mockAsset);
      } else {
        setAsset(null);
      }
    } catch (err) {
      console.error('Error fetching asset:', err);
      setAsset(null);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (value) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
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
            Loading asset details...
          </p>
        </div>
      </TechContainer>
    );
  }

  if (!asset) {
    return (
      <TechContainer>
        <div style={{ textAlign: 'center', padding: '100px 0' }}>
          <div style={{
            padding: 40,
            background: 'rgba(239, 68, 68, 0.1)',
            border: '1px solid rgba(239, 68, 68, 0.3)',
            borderRadius: 12,
            maxWidth: 480,
            margin: '0 auto'
          }}>
            <div style={{
              width: 64,
              height: 64,
              margin: '0 auto 20px',
              borderRadius: '50%',
              background: 'rgba(239, 68, 68, 0.2)',
              border: '2px solid rgba(239, 68, 68, 0.4)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
              <Activity style={{ width: 32, height: 32, color: 'rgb(239, 68, 68)' }} />
            </div>
            <h3 style={{ color: 'rgb(239, 68, 68)', fontWeight: 700, marginBottom: 12, fontSize: 20 }}>Asset Not Found</h3>
            <p style={{ color: 'rgba(203, 213, 225, 0.8)', marginBottom: 24, fontSize: 15 }}>
              The treasury asset you're looking for doesn't exist or has been removed.
            </p>
            <TechButton onClick={() => navigate('/treasury')} variant="danger">
              <ArrowLeft size={16} />
              Back to Market
            </TechButton>
          </div>
        </div>
      </TechContainer>
    );
  }

  const getTypeColor = (type) => {
    switch (type) {
      case 'T-BILL': return { bg: 'rgba(34, 211, 238, 0.2)', color: 'rgb(34, 211, 238)', border: 'rgba(34, 211, 238, 0.4)' };
      case 'T-NOTE': return { bg: 'rgba(59, 130, 246, 0.2)', color: 'rgb(59, 130, 246)', border: 'rgba(59, 130, 246, 0.4)' };
      case 'T-BOND': return { bg: 'rgba(251, 191, 36, 0.2)', color: 'rgb(251, 191, 36)', border: 'rgba(251, 191, 36, 0.4)' };
      default: return { bg: 'rgba(100, 116, 139, 0.2)', color: 'rgba(203, 213, 225, 0.8)', border: 'rgba(100, 116, 139, 0.4)' };
    }
  };

  const typeStyle = getTypeColor(asset.treasury_type);

  return (
    <TechContainer>
      <div style={{ marginBottom: 24 }}>
        <Link
          to="/treasury"
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: 8,
            color: 'rgb(34, 211, 238)',
            textDecoration: 'none',
            fontSize: 14,
            fontWeight: 600,
            transition: 'all 0.3s'
          }}
          onMouseEnter={(e) => e.currentTarget.style.gap = '12px'}
          onMouseLeave={(e) => e.currentTarget.style.gap = '8px'}
        >
          <ArrowLeft size={16} />
          Back to Market
        </Link>
      </div>

      {/* Asset Header */}
      <TechCard style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', flexWrap: 'wrap', gap: 24 }}>
          <div style={{ flex: 1 }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 16, flexWrap: 'wrap' }}>
              <div style={{
                padding: '6px 16px',
                borderRadius: 12,
                fontSize: 12,
                fontWeight: 600,
                textTransform: 'uppercase',
                background: typeStyle.bg,
                color: typeStyle.color,
                border: `1px solid ${typeStyle.border}`
              }}>
                {asset.treasury_type}
              </div>
              <span style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.7)' }}>{asset.maturity_term}</span>
              {asset.cusip && <span style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.6)' }}>CUSIP: {asset.cusip}</span>}
            </div>
            <h1 style={{ fontSize: 32, fontWeight: 700, color: 'white', marginBottom: 24 }}>
              {asset.asset_code}
            </h1>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))', gap: 20 }}>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginBottom: 6 }}>
                  <DollarSign size={14} style={{ color: 'rgba(203, 213, 225, 0.6)' }} />
                  <p style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', margin: 0, textTransform: 'uppercase', letterSpacing: 0.5 }}>Current Price</p>
                </div>
                <p style={{ fontSize: 22, fontWeight: 700, color: 'white', margin: 0 }}>
                  {formatCurrency(asset.current_price)}
                </p>
              </div>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginBottom: 6 }}>
                  <TrendingUp size={14} style={{ color: 'rgba(203, 213, 225, 0.6)' }} />
                  <p style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', margin: 0, textTransform: 'uppercase', letterSpacing: 0.5 }}>Coupon Rate</p>
                </div>
                <p style={{ fontSize: 22, fontWeight: 700, color: 'rgb(34, 211, 238)', margin: 0 }}>
                  {formatPercent(asset.coupon_rate)}
                </p>
              </div>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginBottom: 6 }}>
                  <Clock size={14} style={{ color: 'rgba(203, 213, 225, 0.6)' }} />
                  <p style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', margin: 0, textTransform: 'uppercase', letterSpacing: 0.5 }}>Maturity Date</p>
                </div>
                <p style={{ fontSize: 15, fontWeight: 600, color: 'white', margin: 0 }}>
                  {formatDate(asset.maturity_date)}
                </p>
              </div>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginBottom: 6 }}>
                  <DollarSign size={14} style={{ color: 'rgba(203, 213, 225, 0.6)' }} />
                  <p style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', margin: 0, textTransform: 'uppercase', letterSpacing: 0.5 }}>Face Value</p>
                </div>
                <p style={{ fontSize: 15, fontWeight: 600, color: 'white', margin: 0 }}>
                  {formatCurrency(asset.face_value)}
                </p>
              </div>
            </div>
          </div>

          <div style={{ textAlign: 'right', minWidth: 200 }}>
            <div style={{ fontSize: 48, fontWeight: 700, color: 'rgb(34, 211, 238)', textShadow: '0 0 30px rgba(34, 211, 238, 0.6)', marginBottom: 8 }}>
              {formatPercent(asset.current_yield)}
            </div>
            <div style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.6)', textTransform: 'uppercase', letterSpacing: 1, fontWeight: 600 }}>
              Annual Yield
            </div>
          </div>
        </div>
      </TechCard>

      {/* Stats Row */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: 16, marginBottom: 24 }}>
        <TechCard
          icon={BarChart3}
          title="Total Supply"
          value={formatCurrency(asset.total_supply)}
          subtitle="Total issued"
          iconColor="rgb(34, 211, 238)"
        />
        <TechCard
          icon={DollarSign}
          title="Available Supply"
          value={formatCurrency(asset.available_supply)}
          subtitle={`${((asset.available_supply / asset.total_supply) * 100).toFixed(0)}% available`}
          iconColor="rgb(59, 130, 246)"
        />
        <TechCard
          icon={Users}
          title="Total Holders"
          value={asset.total_holders.toLocaleString()}
          subtitle="Active investors"
          iconColor="rgb(34, 197, 94)"
        />
        <TechCard
          icon={Activity}
          title="24h Volume"
          value={formatCurrency(asset.volume_24h)}
          subtitle="Trading volume"
          iconColor="rgb(251, 191, 36)"
        />
      </div>

      {/* Info Section */}
      <TechCard>
        <h2 style={{ fontSize: 20, fontWeight: 600, color: 'white', margin: '0 0 20px 0' }}>
          Asset Information
        </h2>

        <div style={{
          padding: 24,
          background: 'rgba(34, 211, 238, 0.05)',
          borderRadius: 8,
          border: '1px solid rgba(34, 211, 238, 0.2)'
        }}>
          <div style={{ display: 'grid', gap: 20 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', paddingBottom: 12, borderBottom: '1px solid rgba(34, 211, 238, 0.1)' }}>
              <span style={{ color: 'rgba(203, 213, 225, 0.7)', fontSize: 14 }}>Asset Code</span>
              <span style={{ color: 'white', fontSize: 14, fontWeight: 600 }}>{asset.asset_code}</span>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', paddingBottom: 12, borderBottom: '1px solid rgba(34, 211, 238, 0.1)' }}>
              <span style={{ color: 'rgba(203, 213, 225, 0.7)', fontSize: 14 }}>Treasury Type</span>
              <span style={{ color: 'white', fontSize: 14, fontWeight: 600 }}>{asset.treasury_type}</span>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', paddingBottom: 12, borderBottom: '1px solid rgba(34, 211, 238, 0.1)' }}>
              <span style={{ color: 'rgba(203, 213, 225, 0.7)', fontSize: 14 }}>Maturity Term</span>
              <span style={{ color: 'white', fontSize: 14, fontWeight: 600 }}>{asset.maturity_term}</span>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', paddingBottom: 12, borderBottom: '1px solid rgba(34, 211, 238, 0.1)' }}>
              <span style={{ color: 'rgba(203, 213, 225, 0.7)', fontSize: 14 }}>Current Price</span>
              <span style={{ color: 'rgb(34, 211, 238)', fontSize: 14, fontWeight: 700 }}>{formatCurrency(asset.current_price)}</span>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <span style={{ color: 'rgba(203, 213, 225, 0.7)', fontSize: 14 }}>Annual Yield</span>
              <span style={{ color: 'rgb(34, 197, 94)', fontSize: 14, fontWeight: 700 }}>{formatPercent(asset.current_yield)}</span>
            </div>
          </div>
        </div>

        <div style={{ marginTop: 24 }}>
          <TechButton variant="primary" fullWidth onClick={() => alert('Trading functionality coming soon!')}>
            <ShoppingCart size={18} />
            Buy {asset.asset_code}
          </TechButton>
        </div>
      </TechCard>
    </TechContainer>
  );
}
