import React, { useMemo } from 'react';
import { TrendingUp, TrendingDown, Activity } from 'lucide-react';

export default function PriceChart({ data, asset }) {
  const formatCurrency = (value) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
    }).format(value);
  };

  const formatDate = (timestamp) => {
    return new Date(timestamp).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const chartStats = useMemo(() => {
    if (!data || data.length === 0) {
      return {
        min: 0,
        max: 0,
        change: 0,
        changePercent: 0,
        latest: 0,
        earliest: 0,
      };
    }

    const prices = data.map((d) => parseFloat(d.price));
    const min = Math.min(...prices);
    const max = Math.max(...prices);
    const latest = prices[prices.length - 1];
    const earliest = prices[0];
    const change = latest - earliest;
    const changePercent = (change / earliest) * 100;

    return { min, max, change, changePercent, latest, earliest };
  }, [data]);

  const normalizeData = useMemo(() => {
    if (!data || data.length === 0) return [];

    const prices = data.map((d) => parseFloat(d.price));
    const min = Math.min(...prices);
    const max = Math.max(...prices);
    const range = max - min || 1;

    return data.map((point, index) => ({
      ...point,
      normalized: ((parseFloat(point.price) - min) / range) * 100,
      index,
    }));
  }, [data]);

  const createSVGPath = () => {
    if (normalizeData.length === 0) return '';

    const width = 100;
    const height = 100;
    const points = normalizeData.length;
    const stepX = width / (points - 1 || 1);

    let path = `M 0 ${height - normalizeData[0].normalized}`;

    normalizeData.forEach((point, i) => {
      if (i > 0) {
        const x = i * stepX;
        const y = height - point.normalized;
        path += ` L ${x} ${y}`;
      }
    });

    return path;
  };

  const isPositive = chartStats.change >= 0;

  if (!data || data.length === 0) {
    return (
      <div className="space-y-4">
        <h3 className="text-lg font-semibold">Price Chart</h3>
        <div className="bg-gray-50 rounded-lg p-12 text-center">
          <Activity className="h-12 w-12 text-gray-300 mx-auto mb-3" />
          <p className="text-gray-500">No price history available</p>
          <p className="text-sm text-gray-400 mt-1">
            Price data will appear once trading begins
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-start justify-between">
        <div>
          <h3 className="text-lg font-semibold mb-2">Price Chart</h3>
          <div className="flex items-baseline gap-3">
            <span className="text-3xl font-bold text-gray-900">
              {formatCurrency(chartStats.latest)}
            </span>
            <div
              className={`flex items-center gap-1 ${
                isPositive ? 'text-green-600' : 'text-red-600'
              }`}
            >
              {isPositive ? (
                <TrendingUp className="h-5 w-5" />
              ) : (
                <TrendingDown className="h-5 w-5" />
              )}
              <span className="font-semibold">
                {isPositive ? '+' : ''}
                {formatCurrency(chartStats.change)}
              </span>
              <span className="text-sm">
                ({isPositive ? '+' : ''}
                {chartStats.changePercent.toFixed(2)}%)
              </span>
            </div>
          </div>
        </div>

        <div className="text-right text-sm">
          <div className="text-gray-600">High: {formatCurrency(chartStats.max)}</div>
          <div className="text-gray-600">Low: {formatCurrency(chartStats.min)}</div>
        </div>
      </div>

      {/* SVG Chart */}
      <div className="relative bg-gradient-to-b from-blue-50 to-white rounded-lg p-4 border">
        <svg
          viewBox="0 0 100 100"
          preserveAspectRatio="none"
          className="w-full h-64"
          style={{ overflow: 'visible' }}
        >
          {/* Grid lines */}
          <line
            x1="0"
            y1="25"
            x2="100"
            y2="25"
            stroke="#e5e7eb"
            strokeWidth="0.2"
          />
          <line
            x1="0"
            y1="50"
            x2="100"
            y2="50"
            stroke="#e5e7eb"
            strokeWidth="0.2"
          />
          <line
            x1="0"
            y1="75"
            x2="100"
            y2="75"
            stroke="#e5e7eb"
            strokeWidth="0.2"
          />

          {/* Price line */}
          <path
            d={createSVGPath()}
            fill="none"
            stroke={isPositive ? '#16a34a' : '#dc2626'}
            strokeWidth="0.5"
            vectorEffect="non-scaling-stroke"
          />

          {/* Area fill */}
          <path
            d={`${createSVGPath()} L 100 100 L 0 100 Z`}
            fill={isPositive ? 'rgba(34, 197, 94, 0.1)' : 'rgba(239, 68, 68, 0.1)'}
          />

          {/* Data points */}
          {normalizeData.map((point, i) => {
            const x = (i / (normalizeData.length - 1 || 1)) * 100;
            const y = 100 - point.normalized;
            return (
              <circle
                key={i}
                cx={x}
                cy={y}
                r="0.8"
                fill={isPositive ? '#16a34a' : '#dc2626'}
                className="hover:r-2 transition-all cursor-pointer"
              >
                <title>
                  {formatDate(point.timestamp)}: {formatCurrency(point.price)}
                </title>
              </circle>
            );
          })}
        </svg>
      </div>

      {/* Timeline */}
      <div className="flex justify-between text-xs text-gray-500">
        <span>{data[0] ? formatDate(data[0].timestamp) : ''}</span>
        <span>
          {data[Math.floor(data.length / 2)]
            ? formatDate(data[Math.floor(data.length / 2)].timestamp)
            : ''}
        </span>
        <span>{data[data.length - 1] ? formatDate(data[data.length - 1].timestamp) : ''}</span>
      </div>

      {/* Price Statistics */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t">
        <div>
          <p className="text-xs text-gray-500 mb-1">Open</p>
          <p className="font-semibold">{formatCurrency(chartStats.earliest)}</p>
        </div>
        <div>
          <p className="text-xs text-gray-500 mb-1">Close</p>
          <p className="font-semibold">{formatCurrency(chartStats.latest)}</p>
        </div>
        <div>
          <p className="text-xs text-gray-500 mb-1">High</p>
          <p className="font-semibold text-green-600">
            {formatCurrency(chartStats.max)}
          </p>
        </div>
        <div>
          <p className="text-xs text-gray-500 mb-1">Low</p>
          <p className="font-semibold text-red-600">
            {formatCurrency(chartStats.min)}
          </p>
        </div>
      </div>
    </div>
  );
}
