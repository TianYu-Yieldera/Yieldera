import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  ArrowLeft,
  TrendingUp,
  Users,
  Calendar,
  ShoppingCart,
  Activity,
} from 'lucide-react';
import treasuryService from '../services/treasuryService';
import OrderBook from '../components/OrderBook';
import TradingForm from '../components/TradingForm';
import PriceChart from '../components/PriceChart';

export default function TreasuryDetailView() {
  const { assetId } = useParams();
  const [asset, setAsset] = useState(null);
  const [priceHistory, setPriceHistory] = useState([]);
  const [trades, setTrades] = useState([]);
  const [orders, setOrders] = useState({ buy_orders: [], sell_orders: [] });
  const [activeTab, setActiveTab] = useState('chart');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchAssetData();
    const interval = setInterval(fetchOrders, 10000); // Refresh orders every 10s
    return () => clearInterval(interval);
  }, [assetId]);

  const fetchAssetData = async () => {
    try {
      setLoading(true);
      const [assetData, historyData, tradesData, ordersData] = await Promise.all([
        treasuryService.getAsset(assetId),
        treasuryService.getPriceHistory(assetId, 100),
        treasuryService.getTradeHistory(assetId, 50),
        treasuryService.getMarketOrders(assetId),
      ]);

      setAsset(assetData.asset);
      setPriceHistory(historyData.history || []);
      setTrades(tradesData.trades || []);
      setOrders(ordersData);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const fetchOrders = async () => {
    try {
      const ordersData = await treasuryService.getMarketOrders(assetId);
      setOrders(ordersData);
    } catch (err) {
      console.error('Failed to refresh orders:', err);
    }
  };

  const handleOrderCreated = () => {
    fetchOrders();
    fetchAssetData();
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
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  if (error || !asset) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 max-w-md">
          <h3 className="text-red-800 font-semibold mb-2">Error</h3>
          <p className="text-red-600">{error || 'Asset not found'}</p>
          <Link
            to="/treasury"
            className="mt-4 inline-block bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700"
          >
            Back to Market
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Back Button */}
        <Link
          to="/treasury"
          className="inline-flex items-center text-blue-600 hover:text-blue-700 mb-6"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Market
        </Link>

        {/* Asset Header */}
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <div className="flex items-start justify-between">
            <div>
              <div className="flex items-center gap-3 mb-2">
                <span className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm font-semibold">
                  {asset.treasury_type}
                </span>
                <span className="text-lg text-gray-700">{asset.maturity_term}</span>
                <span className="text-sm text-gray-500">CUSIP: {asset.cusip}</span>
              </div>
              <h1 className="text-3xl font-bold text-gray-900 mb-4">
                US Treasury {asset.treasury_type}
              </h1>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
                <div>
                  <p className="text-sm text-gray-500">Current Price</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {formatCurrency(asset.current_price || asset.face_value)}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">Current Yield</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {formatPercent(asset.current_yield || asset.coupon_rate)}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">Face Value</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {formatCurrency(asset.face_value)}
                  </p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">Maturity Date</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {formatDate(asset.maturity_date)}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Tabs */}
        <div className="bg-white rounded-lg shadow mb-6">
          <div className="border-b border-gray-200">
            <nav className="flex -mb-px">
              {[
                { key: 'chart', label: 'Price Chart', icon: Activity },
                { key: 'orderbook', label: 'Order Book', icon: ShoppingCart },
                { key: 'trades', label: 'Trade History', icon: Users },
              ].map(({ key, label, icon: Icon }) => (
                <button
                  key={key}
                  onClick={() => setActiveTab(key)}
                  className={`${
                    activeTab === key
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700'
                  } flex items-center gap-2 whitespace-nowrap py-4 px-6 border-b-2 font-medium text-sm`}
                >
                  <Icon className="h-4 w-4" />
                  {label}
                </button>
              ))}
            </nav>
          </div>
        </div>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Column - Chart/OrderBook/Trades */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow p-6">
              {activeTab === 'chart' && (
                <PriceChart data={priceHistory} asset={asset} />
              )}
              {activeTab === 'orderbook' && (
                <OrderBook buyOrders={orders.buy_orders} sellOrders={orders.sell_orders} />
              )}
              {activeTab === 'trades' && (
                <div>
                  <h3 className="text-lg font-semibold mb-4">Recent Trades</h3>
                  {trades.length === 0 ? (
                    <p className="text-gray-500 text-center py-8">No trades yet</p>
                  ) : (
                    <div className="space-y-2">
                      {trades.map((trade) => (
                        <div
                          key={trade.trade_id}
                          className="flex justify-between items-center p-3 border rounded hover:bg-gray-50"
                        >
                          <div>
                            <p className="font-medium">
                              {formatCurrency(trade.price_per_token)}
                            </p>
                            <p className="text-sm text-gray-500">
                              Amount: {parseFloat(trade.token_amount).toFixed(2)}
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="text-sm text-gray-500">
                              {new Date(trade.executed_at).toLocaleString()}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>

          {/* Right Column - Trading Form */}
          <div className="lg:col-span-1">
            <TradingForm
              asset={asset}
              onOrderCreated={handleOrderCreated}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
