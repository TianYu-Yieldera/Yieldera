import React, { useState } from 'react';
import { ShoppingCart, DollarSign } from 'lucide-react';
import { useWallet } from '../web3/WalletContext';
import treasuryService from '../services/treasuryService';

export default function TradingForm({ asset, onOrderCreated }) {
  const { account, isConnected, signer } = useWallet();
  const [orderType, setOrderType] = useState('BUY');
  const [amount, setAmount] = useState('');
  const [price, setPrice] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);

  const formatCurrency = (value) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(value);
  };

  const calculateTotal = () => {
    if (!amount || !price) return 0;
    return parseFloat(amount) * parseFloat(price);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);

    if (!isConnected) {
      setError('Please connect your wallet');
      return;
    }

    if (!signer) {
      setError('Signer not available');
      return;
    }

    if (!amount || !price) {
      setError('Please enter amount and price');
      return;
    }

    if (parseFloat(amount) <= 0 || parseFloat(price) <= 0) {
      setError('Amount and price must be greater than 0');
      return;
    }

    try {
      setLoading(true);

      // Create EIP-712 typed data for signature
      const domain = {
        name: 'Treasury Marketplace',
        version: '1',
        chainId: await signer.provider.getNetwork().then(n => n.chainId),
        verifyingContract: import.meta.env.VITE_TREASURY_MARKETPLACE_ADDRESS || '0x90708d3663C3BE0DF3002dC293Bb06c45b67a334',
      };

      const types = {
        Order: [
          { name: 'assetId', type: 'uint256' },
          { name: 'userAddress', type: 'address' },
          { name: 'orderType', type: 'string' },
          { name: 'tokenAmount', type: 'uint256' },
          { name: 'pricePerToken', type: 'uint256' },
          { name: 'timestamp', type: 'uint256' },
        ],
      };

      const value = {
        assetId: asset.asset_id,
        userAddress: account,
        orderType: orderType,
        tokenAmount: Math.floor(parseFloat(amount) * 1e18).toString(), // Convert to wei
        pricePerToken: Math.floor(parseFloat(price) * 1e18).toString(), // Convert to wei
        timestamp: Math.floor(Date.now() / 1000),
      };

      // Sign the typed data
      const signature = await signer.signTypedData(domain, types, value);

      const orderData = {
        asset_id: asset.asset_id,
        user_address: account,
        order_type: orderType,
        token_amount: parseFloat(amount),
        price_per_token: parseFloat(price),
        signature: signature,
      };

      await treasuryService.createOrder(orderData);

      setSuccess(`${orderType} order created successfully!`);
      setAmount('');
      setPrice('');

      // Notify parent component to refresh data
      if (onOrderCreated) {
        onOrderCreated();
      }

      // Clear success message after 3 seconds
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      console.error('Order creation error:', err);
      setError(err.message || 'Failed to create order');
    } finally {
      setLoading(false);
    }
  };

  const setMarketPrice = () => {
    const currentPrice = asset.current_price || asset.face_value;
    setPrice(currentPrice.toString());
  };

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center gap-2 mb-6">
        <ShoppingCart className="h-5 w-5 text-blue-600" />
        <h3 className="text-lg font-semibold">Place Order</h3>
      </div>

      {/* Order Type Selector */}
      <div className="flex gap-2 mb-6">
        <button
          onClick={() => setOrderType('BUY')}
          className={`flex-1 py-2 rounded-lg font-medium transition-colors ${
            orderType === 'BUY'
              ? 'bg-green-600 text-white'
              : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
          }`}
        >
          Buy
        </button>
        <button
          onClick={() => setOrderType('SELL')}
          className={`flex-1 py-2 rounded-lg font-medium transition-colors ${
            orderType === 'SELL'
              ? 'bg-red-600 text-white'
              : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
          }`}
        >
          Sell
        </button>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Amount Input */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Amount (Tokens)
          </label>
          <input
            type="number"
            step="0.01"
            min="0"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            placeholder="0.00"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            disabled={loading}
          />
        </div>

        {/* Price Input */}
        <div>
          <div className="flex justify-between items-center mb-2">
            <label className="block text-sm font-medium text-gray-700">
              Price per Token
            </label>
            <button
              type="button"
              onClick={setMarketPrice}
              className="text-xs text-blue-600 hover:text-blue-700"
            >
              Use Market Price
            </button>
          </div>
          <div className="relative">
            <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
            <input
              type="number"
              step="0.01"
              min="0"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              placeholder="0.00"
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              disabled={loading}
            />
          </div>
        </div>

        {/* Order Summary */}
        <div className="bg-gray-50 rounded-lg p-4 space-y-2">
          <div className="flex justify-between text-sm">
            <span className="text-gray-600">Amount:</span>
            <span className="font-medium">{amount || '0.00'} tokens</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-gray-600">Price:</span>
            <span className="font-medium">{formatCurrency(price || 0)}</span>
          </div>
          <div className="border-t pt-2 flex justify-between">
            <span className="font-medium text-gray-900">Total:</span>
            <span className="font-bold text-gray-900">
              {formatCurrency(calculateTotal())}
            </span>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-3 text-sm text-red-600">
            {error}
          </div>
        )}

        {/* Success Message */}
        {success && (
          <div className="bg-green-50 border border-green-200 rounded-lg p-3 text-sm text-green-600">
            {success}
          </div>
        )}

        {/* Submit Button */}
        <button
          type="submit"
          disabled={loading || !isConnected}
          className={`w-full py-3 rounded-lg font-semibold transition-colors ${
            orderType === 'BUY'
              ? 'bg-green-600 hover:bg-green-700 text-white'
              : 'bg-red-600 hover:bg-red-700 text-white'
          } disabled:bg-gray-300 disabled:cursor-not-allowed`}
        >
          {loading ? (
            <span className="flex items-center justify-center">
              <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
              Processing...
            </span>
          ) : !isConnected ? (
            'Connect Wallet'
          ) : (
            `Place ${orderType} Order`
          )}
        </button>
      </form>

      {/* Market Info */}
      <div className="mt-6 pt-6 border-t">
        <div className="space-y-2 text-sm">
          <div className="flex justify-between">
            <span className="text-gray-600">Current Market Price:</span>
            <span className="font-medium">
              {formatCurrency(asset.current_price || asset.face_value)}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-600">Face Value:</span>
            <span className="font-medium">{formatCurrency(asset.face_value)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-600">Current Yield:</span>
            <span className="font-medium">
              {((asset.current_yield || asset.coupon_rate) * 100).toFixed(2)}%
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
