import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Wallet, TrendingUp, TrendingDown, DollarSign, Award } from 'lucide-react';
import { useWallet } from '../web3/WalletContext';
import treasuryService from '../services/treasuryService';

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
      // Create EIP-712 typed data for signature
      const domain = {
        name: 'Treasury Yield Distributor',
        version: '1',
        chainId: await signer.provider.getNetwork().then(n => n.chainId),
        verifyingContract: '0x0000000000000000000000000000000000000000', // TODO: Replace with actual yield distributor contract address
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

      // Sign the typed data
      const signature = await signer.signTypedData(domain, types, value);

      await treasuryService.claimYield({
        asset_id: assetId,
        user_address: account,
        signature: signature,
      });

      alert('Yield claimed successfully!');
      fetchHoldings(); // Refresh after claim
    } catch (err) {
      console.error('Claim yield error:', err);
      alert('Failed to claim yield: ' + err.message);
    }
  };

  if (!isConnected) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="bg-white rounded-lg shadow-lg p-8 max-w-md text-center">
          <Wallet className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-gray-900 mb-2">
            Connect Your Wallet
          </h2>
          <p className="text-gray-600 mb-6">
            Please connect your wallet to view your Treasury holdings
          </p>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading your holdings...</p>
        </div>
      </div>
    );
  }

  const totals = calculateTotals();
  const isProfit = totals.totalGainLoss >= 0;

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">My Treasury Holdings</h1>
          <p className="mt-2 text-gray-600">
            Track your Treasury portfolio and yield earnings
          </p>
        </div>

        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <DollarSign className="h-8 w-8 text-blue-600" />
              <div className="ml-4">
                <p className="text-sm text-gray-500">Total Invested</p>
                <p className="text-2xl font-bold text-gray-900">
                  {formatCurrency(totals.totalInvested)}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <Wallet className="h-8 w-8 text-green-600" />
              <div className="ml-4">
                <p className="text-sm text-gray-500">Current Value</p>
                <p className="text-2xl font-bold text-gray-900">
                  {formatCurrency(totals.currentValue)}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              {isProfit ? (
                <TrendingUp className="h-8 w-8 text-green-600" />
              ) : (
                <TrendingDown className="h-8 w-8 text-red-600" />
              )}
              <div className="ml-4">
                <p className="text-sm text-gray-500">Gain/Loss</p>
                <p
                  className={`text-2xl font-bold ${
                    isProfit ? 'text-green-600' : 'text-red-600'
                  }`}
                >
                  {isProfit ? '+' : ''}
                  {formatCurrency(totals.totalGainLoss)}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <Award className="h-8 w-8 text-purple-600" />
              <div className="ml-4">
                <p className="text-sm text-gray-500">Accrued Yield</p>
                <p className="text-2xl font-bold text-gray-900">
                  {formatCurrency(totals.totalYield)}
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Holdings Table */}
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">Your Holdings</h2>
          </div>

          {holdings.length === 0 ? (
            <div className="p-12 text-center">
              <Wallet className="h-16 w-16 text-gray-300 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">
                No Holdings Yet
              </h3>
              <p className="text-gray-500 mb-6">
                Start investing in US Treasury securities
              </p>
              <Link
                to="/treasury"
                className="inline-block bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700"
              >
                Browse Treasury Market
              </Link>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Asset
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Tokens Held
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Avg Price
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Current Value
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Gain/Loss
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Accrued Yield
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {holdings.map((holding) => {
                    const gainLoss =
                      parseFloat(holding.current_value || 0) -
                      parseFloat(holding.total_invested || 0);
                    const isProfit = gainLoss >= 0;

                    return (
                      <tr key={holding.asset_id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div>
                            <div className="text-sm font-medium text-gray-900">
                              {holding.treasury_type}
                            </div>
                            <div className="text-sm text-gray-500">
                              {holding.maturity_term} | {holding.cusip}
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {parseFloat(holding.tokens_held).toFixed(2)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {formatCurrency(holding.avg_purchase_price)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {formatCurrency(holding.current_value)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span
                            className={`text-sm font-medium ${
                              isProfit ? 'text-green-600' : 'text-red-600'
                            }`}
                          >
                            {isProfit ? '+' : ''}
                            {formatCurrency(gainLoss)}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {formatCurrency(holding.accrued_interest)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                          <Link
                            to={`/treasury/${holding.asset_id}`}
                            className="text-blue-600 hover:text-blue-900"
                          >
                            Trade
                          </Link>
                          {parseFloat(holding.accrued_interest) > 0 && (
                            <button
                              onClick={() => handleClaimYield(holding.asset_id)}
                              className="text-green-600 hover:text-green-900"
                            >
                              Claim
                            </button>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </div>

        {/* Yield History */}
        {yieldData && yieldData.distributions_count > 0 && (
          <div className="mt-8 bg-white rounded-lg shadow">
            <div className="px-6 py-4 border-b border-gray-200">
              <h2 className="text-lg font-semibold text-gray-900">
                Yield Distribution History
              </h2>
            </div>
            <div className="p-6">
              <div className="space-y-3">
                {yieldData.distributions.slice(0, 5).map((dist) => (
                  <div
                    key={dist.id}
                    className="flex justify-between items-center p-4 border rounded hover:bg-gray-50"
                  >
                    <div>
                      <p className="font-medium text-gray-900">
                        {dist.treasury_type} - {dist.cusip}
                      </p>
                      <p className="text-sm text-gray-500">
                        {dist.distribution_type} on {dist.distribution_date}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-bold text-green-600">
                        +{formatCurrency(dist.user_yield)}
                      </p>
                      <p className="text-sm text-gray-500">{dist.status}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
