import React from 'react';
import { TrendingUp, TrendingDown } from 'lucide-react';

export default function OrderBook({ buyOrders, sellOrders }) {
  const formatCurrency = (value) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
    }).format(value);
  };

  const formatAmount = (value) => {
    return parseFloat(value).toFixed(2);
  };

  const calculateTotal = (price, amount) => {
    return parseFloat(price) * parseFloat(amount);
  };

  return (
    <div className="space-y-6">
      <h3 className="text-lg font-semibold">Order Book</h3>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Buy Orders */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <TrendingUp className="h-5 w-5 text-green-600" />
            <h4 className="font-semibold text-green-600">Buy Orders</h4>
          </div>

          <div className="border rounded-lg overflow-hidden">
            <div className="bg-gray-50 grid grid-cols-3 gap-2 p-2 text-xs font-medium text-gray-600">
              <div>Price</div>
              <div className="text-right">Amount</div>
              <div className="text-right">Total</div>
            </div>

            <div className="divide-y max-h-96 overflow-y-auto">
              {buyOrders && buyOrders.length > 0 ? (
                buyOrders.map((order) => (
                  <div
                    key={order.order_id}
                    className="grid grid-cols-3 gap-2 p-2 text-sm hover:bg-green-50 transition-colors"
                  >
                    <div className="text-green-600 font-medium">
                      {formatCurrency(order.price_per_token)}
                    </div>
                    <div className="text-right text-gray-900">
                      {formatAmount(order.token_amount)}
                    </div>
                    <div className="text-right text-gray-600">
                      {formatCurrency(
                        calculateTotal(order.price_per_token, order.token_amount)
                      )}
                    </div>
                  </div>
                ))
              ) : (
                <div className="p-8 text-center text-gray-500 text-sm">
                  No buy orders
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Sell Orders */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <TrendingDown className="h-5 w-5 text-red-600" />
            <h4 className="font-semibold text-red-600">Sell Orders</h4>
          </div>

          <div className="border rounded-lg overflow-hidden">
            <div className="bg-gray-50 grid grid-cols-3 gap-2 p-2 text-xs font-medium text-gray-600">
              <div>Price</div>
              <div className="text-right">Amount</div>
              <div className="text-right">Total</div>
            </div>

            <div className="divide-y max-h-96 overflow-y-auto">
              {sellOrders && sellOrders.length > 0 ? (
                sellOrders.map((order) => (
                  <div
                    key={order.order_id}
                    className="grid grid-cols-3 gap-2 p-2 text-sm hover:bg-red-50 transition-colors"
                  >
                    <div className="text-red-600 font-medium">
                      {formatCurrency(order.price_per_token)}
                    </div>
                    <div className="text-right text-gray-900">
                      {formatAmount(order.token_amount)}
                    </div>
                    <div className="text-right text-gray-600">
                      {formatCurrency(
                        calculateTotal(order.price_per_token, order.token_amount)
                      )}
                    </div>
                  </div>
                ))
              ) : (
                <div className="p-8 text-center text-gray-500 text-sm">
                  No sell orders
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Market Summary */}
      {(buyOrders?.length > 0 || sellOrders?.length > 0) && (
        <div className="bg-gray-50 rounded-lg p-4">
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span className="text-gray-600">Best Bid:</span>
              <span className="ml-2 font-semibold text-green-600">
                {buyOrders?.[0]
                  ? formatCurrency(buyOrders[0].price_per_token)
                  : 'N/A'}
              </span>
            </div>
            <div>
              <span className="text-gray-600">Best Ask:</span>
              <span className="ml-2 font-semibold text-red-600">
                {sellOrders?.[0]
                  ? formatCurrency(sellOrders[0].price_per_token)
                  : 'N/A'}
              </span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
