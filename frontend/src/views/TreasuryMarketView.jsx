import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { TrendingUp, TrendingDown, DollarSign, Calendar, Percent } from 'lucide-react';
import treasuryService from '../services/treasuryService';

export default function TreasuryMarketView() {
  const [assets, setAssets] = useState([]);
  const [filteredAssets, setFilteredAssets] = useState([]);
  const [stats, setStats] = useState(null);
  const [selectedType, setSelectedType] = useState('ALL');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

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
      const [assetsData, statsData] = await Promise.all([
        treasuryService.getAssets(),
        treasuryService.getStats(),
      ]);

      setAssets(assetsData.assets || []);
      setFilteredAssets(assetsData.assets || []);
      setStats(statsData);
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
      minimumFractionDigits: 2,
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

  const getTypeColor = (type) => {
    switch (type) {
      case 'T-BILL':
        return 'bg-blue-100 text-blue-800';
      case 'T-NOTE':
        return 'bg-green-100 text-green-800';
      case 'T-BOND':
        return 'bg-purple-100 text-purple-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const calculatePriceChange = (asset) => {
    // Mock calculation - in production, compare with previous price
    const change = Math.random() * 2 - 1; // -1% to +1%
    return change;
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading Treasury Market...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 max-w-md">
          <h3 className="text-red-800 font-semibold mb-2">Error Loading Data</h3>
          <p className="text-red-600">{error}</p>
          <button
            onClick={fetchData}
            className="mt-4 bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">US Treasury Market</h1>
          <p className="mt-2 text-gray-600">
            Trade tokenized US Treasury securities on-chain
          </p>
        </div>

        {/* Stats Cards */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <DollarSign className="h-8 w-8 text-blue-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm text-gray-500">Total TVL</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {formatCurrency(stats.total_tvl)}
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <TrendingUp className="h-8 w-8 text-green-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm text-gray-500">Total Assets</p>
                  <p className="text-2xl font-bold text-gray-900">{stats.total_assets}</p>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <Calendar className="h-8 w-8 text-purple-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm text-gray-500">24h Volume</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {formatCurrency(stats.volume_24h || 0)}
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <Percent className="h-8 w-8 text-orange-600" />
                </div>
                <div className="ml-4">
                  <p className="text-sm text-gray-500">Active Orders</p>
                  <p className="text-2xl font-bold text-gray-900">{stats.active_orders}</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Filter Tabs */}
        <div className="bg-white rounded-lg shadow mb-6">
          <div className="border-b border-gray-200">
            <nav className="flex -mb-px">
              {['ALL', 'T-BILL', 'T-NOTE', 'T-BOND'].map((type) => (
                <button
                  key={type}
                  onClick={() => setSelectedType(type)}
                  className={`${
                    selectedType === type
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  } whitespace-nowrap py-4 px-8 border-b-2 font-medium text-sm`}
                >
                  {type === 'ALL' ? 'All Securities' : type}
                </button>
              ))}
            </nav>
          </div>
        </div>

        {/* Assets Grid */}
        {filteredAssets.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-12 text-center">
            <p className="text-gray-500">No treasury assets available</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-6">
            {filteredAssets.map((asset) => {
              const priceChange = calculatePriceChange(asset);
              const isPositive = priceChange >= 0;

              return (
                <Link
                  key={asset.asset_id}
                  to={`/treasury/${asset.asset_id}`}
                  className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow duration-200 p-6"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-2">
                        <span
                          className={`px-3 py-1 rounded-full text-xs font-semibold ${getTypeColor(
                            asset.treasury_type
                          )}`}
                        >
                          {asset.treasury_type}
                        </span>
                        <span className="text-sm text-gray-500">
                          {asset.maturity_term}
                        </span>
                        <span className="text-xs text-gray-400">
                          CUSIP: {asset.cusip}
                        </span>
                      </div>

                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-4">
                        <div>
                          <p className="text-sm text-gray-500">Current Price</p>
                          <p className="text-lg font-bold text-gray-900">
                            {formatCurrency(asset.current_price || asset.face_value)}
                          </p>
                          <div className="flex items-center mt-1">
                            {isPositive ? (
                              <TrendingUp className="h-4 w-4 text-green-500" />
                            ) : (
                              <TrendingDown className="h-4 w-4 text-red-500" />
                            )}
                            <span
                              className={`text-xs ml-1 ${
                                isPositive ? 'text-green-600' : 'text-red-600'
                              }`}
                            >
                              {isPositive ? '+' : ''}
                              {priceChange.toFixed(2)}%
                            </span>
                          </div>
                        </div>

                        <div>
                          <p className="text-sm text-gray-500">Current Yield</p>
                          <p className="text-lg font-bold text-gray-900">
                            {formatPercent(asset.current_yield || asset.coupon_rate)}
                          </p>
                        </div>

                        <div>
                          <p className="text-sm text-gray-500">Face Value</p>
                          <p className="text-lg font-bold text-gray-900">
                            {formatCurrency(asset.face_value)}
                          </p>
                        </div>

                        <div>
                          <p className="text-sm text-gray-500">Maturity Date</p>
                          <p className="text-lg font-bold text-gray-900">
                            {formatDate(asset.maturity_date)}
                          </p>
                        </div>
                      </div>
                    </div>

                    <div className="ml-6 text-right">
                      <button className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition-colors">
                        Trade
                      </button>
                    </div>
                  </div>
                </Link>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
