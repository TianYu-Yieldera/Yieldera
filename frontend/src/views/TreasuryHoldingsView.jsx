import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Wallet, TrendingUp, DollarSign, Award, ArrowLeft } from 'lucide-react';
import { useWallet } from '../web3/WalletContext';
import treasuryService from '../services/treasuryService';
import TechContainer from "../components/ui/TechContainer";
import TechHeader from "../components/ui/TechHeader";
import TechCard from "../components/ui/TechCard";
import TechButton from "../components/ui/TechButton";

export default function TreasuryHoldingsView() {
  const { account, isConnected, signer } = useWallet();
  const [holdings, setHoldings] = useState([]);
  const [yieldData, setYieldData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (isConnected && account) {
      fetchHoldings();
    }
  }, [account, isConnected]);

  const fetchHoldings = async () => {
    try {
      setLoading(true);
      const [holdingsData, yieldInfo] = await Promise.all([
        treasuryService.getUserHoldings(account),
        treasuryService.getUserYield(account),
      ]);

      setHoldings(holdingsData.holdings || []);
      setYieldData(yieldInfo);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (value) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(value || 0);
  };

  const formatPercent = (value) => {
    return `${(parseFloat(value || 0) * 100).toFixed(2)}%`;
  };

  const calculateTotals = () => {
    if (!holdings || holdings.length === 0) {
      return {
        totalInvested: 0,
        currentValue: 0,
        totalGainLoss: 0,
        totalYield: 0,
      };
    }

    const totalInvested = holdings.reduce(
      (sum, h) => sum + parseFloat(h.total_invested || 0),
      0
    );
    const currentValue = holdings.reduce(
      (sum, h) => sum + parseFloat(h.current_value || 0),
      0
    );
    const totalYield = holdings.reduce(
      (sum, h) => sum + parseFloat(h.accrued_interest || 0),
      0
    );

    return {
      totalInvested,
      currentValue,
      totalGainLoss: currentValue - totalInvested,
      totalYield,
    };
  };

  const handleClaimYield = async (assetId) => {
    if (!signer) {
      alert('Please connect your wallet first');
      return;
    }

    try {
      const domain = {
        name: 'Treasury Yield Distributor',
        version: '1',
        chainId: await signer.provider.getNetwork().then(n => n.chainId),
        verifyingContract: import.meta.env.VITE_TREASURY_YIELD_DISTRIBUTOR_ADDRESS || '0x0BE14D40188FCB5924c36af46630faBD76698A80',
      };

      const types = {
        ClaimYield: [
          { name: 'assetId', type: 'uint256' },
          { name: 'userAddress', type: 'address' },
          { name: 'timestamp', type: 'uint256' },
        ],
      };

      const value = {
        assetId: assetId,
        userAddress: account,
        timestamp: Math.floor(Date.now() / 1000),
      };

      const signature = await signer.signTypedData(domain, types, value);

      await treasuryService.claimYield({
        asset_id: assetId,
        user_address: account,
        signature: signature,
      });

      alert('Yield claimed successfully!');
      fetchHoldings();
    } catch (err) {
      console.error('Claim yield error:', err);
      alert('Failed to claim yield: ' + err.message);
    }
  };

  if (!isConnected) {
    return (
      <TechContainer>
        <div style={{ textAlign: 'center', padding: '100px 0' }}>
          <div style={{
            width: 80,
            height: 80,
            margin: '0 auto 24px',
            borderRadius: '50%',
            background: 'rgba(34, 211, 238, 0.15)',
            border: '2px solid rgba(34, 211, 238, 0.4)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
          }}>
            <Wallet style={{ width: 40, height: 40, color: 'rgb(34, 211, 238)' }} />
          </div>
          <h2 style={{ fontSize: 28, fontWeight: 700, color: 'white', marginBottom: 12 }}>
            Connect Your Wallet
          </h2>
          <p style={{ color: 'rgba(203, 213, 225, 0.7)', marginBottom: 24, fontSize: 16 }}>
            Please connect your wallet to view your Treasury holdings
          </p>
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
            Loading your holdings...
          </p>
        </div>
      </TechContainer>
    );
  }

  const totals = calculateTotals();

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
            fontWeight: 500,
            transition: 'all 0.3s'
          }}
          onMouseEnter={(e) => e.currentTarget.style.gap = '12px'}
          onMouseLeave={(e) => e.currentTarget.style.gap = '8px'}
        >
          <ArrowLeft size={16} />
          Back to Market
        </Link>
      </div>

      <TechHeader
        icon={Wallet}
        title="My Treasury Holdings"
        subtitle="View and manage your US Treasury bond investments"
      />

      {/* Summary Stats */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: 16, marginBottom: 32 }}>
        <TechCard
          icon={DollarSign}
          title="Total Invested"
          value={formatCurrency(totals.totalInvested)}
          subtitle="Principal amount"
          iconColor="rgb(34, 211, 238)"
        />
        <TechCard
          icon={TrendingUp}
          title="Current Value"
          value={formatCurrency(totals.currentValue)}
          subtitle="Including accrued yield"
          iconColor="rgb(59, 130, 246)"
        />
        <TechCard
          icon={Award}
          title="Total Yield"
          value={formatCurrency(totals.totalYield)}
          subtitle="Accrued interest"
          iconColor="rgb(34, 197, 94)"
        />
        <TechCard
          icon={TrendingUp}
          title="Gain/Loss"
          value={formatCurrency(totals.totalGainLoss)}
          subtitle={totals.totalGainLoss >= 0 ? 'profit' : 'loss'}
          iconColor={totals.totalGainLoss >= 0 ? 'rgb(34, 197, 94)' : 'rgb(239, 68, 68)'}
        />
      </div>

      {/* Holdings List */}
      <TechCard>
        <h2 style={{ fontSize: 20, fontWeight: 600, color: 'white', margin: '0 0 20px 0' }}>
          Your Holdings
        </h2>
        {holdings.length === 0 ? (
          <div style={{ textAlign: 'center', padding: 48 }}>
            <p style={{ color: 'rgba(203, 213, 225, 0.7)', fontSize: 16 }}>
              No holdings found. Visit the{' '}
              <Link to="/treasury" style={{ color: 'rgb(34, 211, 238)', textDecoration: 'none', fontWeight: 600 }}>
                Treasury Market
              </Link>
              {' '}to start investing.
            </p>
          </div>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
            {holdings.map((holding, index) => (
              <div
                key={index}
                style={{
                  padding: 24,
                  background: 'rgba(15, 23, 42, 0.5)',
                  borderRadius: 12,
                  border: '1px solid rgba(34, 211, 238, 0.2)',
                  transition: 'all 0.3s'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.4)';
                  e.currentTarget.style.transform = 'translateY(-2px)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.2)';
                  e.currentTarget.style.transform = 'translateY(0)';
                }}
              >
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: 16, flexWrap: 'wrap', gap: 16 }}>
                  <div>
                    <h3 style={{ fontSize: 18, fontWeight: 600, color: 'white', marginBottom: 8 }}>
                      {holding.asset_code}
                    </h3>
                    <div style={{ fontSize: 14, color: 'rgba(203, 213, 225, 0.7)' }}>
                      {holding.treasury_type} • Maturity: {new Date(holding.maturity_date).toLocaleDateString()}
                    </div>
                  </div>
                  <div style={{ textAlign: 'right' }}>
                    <div style={{ fontSize: 24, fontWeight: 700, color: 'rgb(34, 211, 238)' }}>
                      {formatCurrency(holding.current_value)}
                    </div>
                    <div style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.6)' }}>
                      Current Value
                    </div>
                  </div>
                </div>

                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))', gap: 16, marginBottom: 16 }}>
                  <div>
                    <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', marginBottom: 4 }}>Invested</div>
                    <div style={{ fontSize: 16, fontWeight: 600, color: 'white' }}>{formatCurrency(holding.total_invested)}</div>
                  </div>
                  <div>
                    <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', marginBottom: 4 }}>Quantity</div>
                    <div style={{ fontSize: 16, fontWeight: 600, color: 'white' }}>{holding.quantity || 0}</div>
                  </div>
                  <div>
                    <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', marginBottom: 4 }}>Accrued Yield</div>
                    <div style={{ fontSize: 16, fontWeight: 600, color: 'rgb(34, 197, 94)' }}>{formatCurrency(holding.accrued_interest)}</div>
                  </div>
                  <div>
                    <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)', marginBottom: 4 }}>Current APY</div>
                    <div style={{ fontSize: 16, fontWeight: 600, color: 'white' }}>{formatPercent(holding.current_yield)}</div>
                  </div>
                </div>

                {parseFloat(holding.accrued_interest) > 0 && (
                  <TechButton
                    onClick={() => handleClaimYield(holding.asset_id)}
                    variant="primary"
                  >
                    Claim Yield ({formatCurrency(holding.accrued_interest)})
                  </TechButton>
                )}
              </div>
            ))}
          </div>
        )}
      </TechCard>
    </TechContainer>
  );
}
