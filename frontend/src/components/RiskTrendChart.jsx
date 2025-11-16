/**
 * Risk Trend Chart Component
 * Visualizes risk score trends and yield distributions
 */

import React, { useState, useEffect } from 'react';
import { TrendingUp, TrendingDown, BarChart3, Activity } from 'lucide-react';

export default function RiskTrendChart({ data, type = 'risk' }) {
  const [hoveredIndex, setHoveredIndex] = useState(null);

  // Generate sample data if not provided
  const chartData = data || generateSampleData(type);

  // Calculate chart dimensions
  const width = 800;
  const height = 300;
  const padding = { top: 20, right: 20, bottom: 40, left: 60 };
  const chartWidth = width - padding.left - padding.right;
  const chartHeight = height - padding.top - padding.bottom;

  // Get min/max values
  const values = chartData.map(d => d.value);
  const minValue = Math.min(...values);
  const maxValue = Math.max(...values);
  const valueRange = maxValue - minValue;

  // Scale functions
  const scaleX = (index) => {
    return (index / (chartData.length - 1)) * chartWidth + padding.left;
  };

  const scaleY = (value) => {
    return height - padding.bottom - ((value - minValue) / valueRange) * chartHeight;
  };

  // Generate path for line chart
  const linePath = chartData.map((d, i) => {
    const x = scaleX(i);
    const y = scaleY(d.value);
    return i === 0 ? `M ${x} ${y}` : `L ${x} ${y}`;
  }).join(' ');

  // Generate area path (for gradient fill)
  const areaPath = linePath +
    ` L ${scaleX(chartData.length - 1)} ${height - padding.bottom}` +
    ` L ${padding.left} ${height - padding.bottom} Z`;

  // Calculate trend
  const firstValue = chartData[0].value;
  const lastValue = chartData[chartData.length - 1].value;
  const trend = lastValue > firstValue ? 'up' : lastValue < firstValue ? 'down' : 'stable';
  const trendPercent = ((lastValue - firstValue) / firstValue * 100).toFixed(1);

  // Get color based on type and trend
  const getColor = () => {
    if (type === 'risk') {
      return trend === 'up' ? 'rgb(239, 68, 68)' : 'rgb(34, 197, 94)';
    } else {
      return trend === 'up' ? 'rgb(34, 197, 94)' : 'rgb(239, 68, 68)';
    }
  };

  const color = getColor();

  return (
    <div style={{
      background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
      borderRadius: 16,
      padding: 24,
      border: '1px solid rgba(34, 211, 238, 0.2)',
      position: 'relative',
      overflow: 'hidden'
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
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 24 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <div style={{
              width: 40,
              height: 40,
              borderRadius: 10,
              background: `${color}20`,
              border: `1px solid ${color}40`,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
              <BarChart3 style={{ width: 20, height: 20, color }} />
            </div>
            <div>
              <h3 style={{ fontSize: 18, fontWeight: 700, color: 'white', margin: 0 }}>
                {type === 'risk' ? 'Risk Score Trend' : 'Yield Performance'}
              </h3>
              <p style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.7)', margin: '2px 0 0 0' }}>
                Last {chartData.length} days
              </p>
            </div>
          </div>

          {/* Trend indicator */}
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: 8,
            padding: '8px 16px',
            borderRadius: 8,
            background: `${color}15`,
            border: `1px solid ${color}30`
          }}>
            {trend === 'up' ? (
              <TrendingUp style={{ width: 18, height: 18, color }} />
            ) : (
              <TrendingDown style={{ width: 18, height: 18, color }} />
            )}
            <span style={{
              fontSize: 16,
              fontWeight: 700,
              color
            }}>
              {trendPercent > 0 ? '+' : ''}{trendPercent}%
            </span>
          </div>
        </div>

        {/* Chart */}
        <div style={{ position: 'relative' }}>
          <svg width="100%" height={height} viewBox={`0 0 ${width} ${height}`} style={{ overflow: 'visible' }}>
            <defs>
              {/* Gradient for area fill */}
              <linearGradient id={`gradient-${type}`} x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" stopColor={color} stopOpacity="0.3" />
                <stop offset="100%" stopColor={color} stopOpacity="0.05" />
              </linearGradient>

              {/* Glow filter */}
              <filter id={`glow-${type}`}>
                <feGaussianBlur stdDeviation="3" result="coloredBlur"/>
                <feMerge>
                  <feMergeNode in="coloredBlur"/>
                  <feMergeNode in="SourceGraphic"/>
                </feMerge>
              </filter>
            </defs>

            {/* Grid lines */}
            {[0, 1, 2, 3, 4].map((i) => {
              const y = padding.top + (chartHeight / 4) * i;
              const value = maxValue - (valueRange / 4) * i;
              return (
                <g key={i}>
                  <line
                    x1={padding.left}
                    y1={y}
                    x2={width - padding.right}
                    y2={y}
                    stroke="rgba(255, 255, 255, 0.1)"
                    strokeWidth="1"
                  />
                  <text
                    x={padding.left - 10}
                    y={y + 4}
                    textAnchor="end"
                    fill="rgba(203, 213, 225, 0.6)"
                    fontSize="12"
                  >
                    {value.toFixed(0)}
                  </text>
                </g>
              );
            })}

            {/* Area fill */}
            <path
              d={areaPath}
              fill={`url(#gradient-${type})`}
            />

            {/* Line */}
            <path
              d={linePath}
              fill="none"
              stroke={color}
              strokeWidth="3"
              strokeLinecap="round"
              strokeLinejoin="round"
              filter={`url(#glow-${type})`}
            />

            {/* Data points */}
            {chartData.map((d, i) => {
              const x = scaleX(i);
              const y = scaleY(d.value);
              const isHovered = hoveredIndex === i;

              return (
                <g key={i}>
                  {/* Outer circle (hover effect) */}
                  {isHovered && (
                    <circle
                      cx={x}
                      cy={y}
                      r="8"
                      fill="none"
                      stroke={color}
                      strokeWidth="2"
                      opacity="0.3"
                    >
                      <animate
                        attributeName="r"
                        from="8"
                        to="12"
                        dur="1s"
                        repeatCount="indefinite"
                      />
                      <animate
                        attributeName="opacity"
                        from="0.3"
                        to="0"
                        dur="1s"
                        repeatCount="indefinite"
                      />
                    </circle>
                  )}

                  {/* Data point */}
                  <circle
                    cx={x}
                    cy={y}
                    r={isHovered ? "5" : "4"}
                    fill="white"
                    stroke={color}
                    strokeWidth="2"
                    style={{ cursor: 'pointer', transition: 'all 0.2s ease' }}
                    onMouseEnter={() => setHoveredIndex(i)}
                    onMouseLeave={() => setHoveredIndex(null)}
                  />

                  {/* Tooltip */}
                  {isHovered && (
                    <g>
                      <rect
                        x={x - 60}
                        y={y - 50}
                        width="120"
                        height="40"
                        rx="8"
                        fill="rgba(15, 23, 42, 0.95)"
                        stroke={color}
                        strokeWidth="1"
                      />
                      <text
                        x={x}
                        y={y - 32}
                        textAnchor="middle"
                        fill="white"
                        fontSize="12"
                        fontWeight="600"
                      >
                        {d.label}
                      </text>
                      <text
                        x={x}
                        y={y - 18}
                        textAnchor="middle"
                        fill={color}
                        fontSize="14"
                        fontWeight="700"
                      >
                        {d.value.toFixed(1)}
                      </text>
                    </g>
                  )}
                </g>
              );
            })}

            {/* X-axis labels */}
            {chartData.filter((_, i) => i % Math.ceil(chartData.length / 6) === 0).map((d, i) => {
              const originalIndex = chartData.indexOf(d);
              const x = scaleX(originalIndex);
              return (
                <text
                  key={i}
                  x={x}
                  y={height - padding.bottom + 25}
                  textAnchor="middle"
                  fill="rgba(203, 213, 225, 0.6)"
                  fontSize="12"
                >
                  {d.label}
                </text>
              );
            })}
          </svg>
        </div>

        {/* Stats */}
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))',
          gap: 16,
          marginTop: 24
        }}>
          {[
            { label: 'Current', value: lastValue.toFixed(1), color: color },
            { label: 'Average', value: (values.reduce((a, b) => a + b) / values.length).toFixed(1), color: 'rgb(34, 211, 238)' },
            { label: 'Peak', value: maxValue.toFixed(1), color: 'rgb(167, 139, 250)' },
            { label: 'Lowest', value: minValue.toFixed(1), color: 'rgb(234, 179, 8)' }
          ].map((stat, i) => (
            <div key={i} style={{
              padding: 12,
              background: 'rgba(255, 255, 255, 0.03)',
              borderRadius: 8,
              border: '1px solid rgba(255, 255, 255, 0.08)'
            }}>
              <p style={{
                fontSize: 11,
                color: 'rgba(203, 213, 225, 0.7)',
                margin: '0 0 6px 0',
                textTransform: 'uppercase',
                letterSpacing: 0.5
              }}>
                {stat.label}
              </p>
              <p style={{
                fontSize: 20,
                fontWeight: 700,
                color: stat.color,
                margin: 0
              }}>
                {stat.value}
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

// Generate sample data for demonstration
function generateSampleData(type) {
  const days = 30;
  const data = [];
  const now = new Date();

  // Base value depends on type
  const baseValue = type === 'risk' ? 35 : 12;
  const volatility = type === 'risk' ? 8 : 3;

  for (let i = days - 1; i >= 0; i--) {
    const date = new Date(now);
    date.setDate(date.getDate() - i);

    // Generate realistic trending data with some noise
    const trend = (days - i) / days; // 0 to 1
    const noise = (Math.random() - 0.5) * volatility;
    const seasonal = Math.sin((i / 7) * Math.PI) * (volatility / 2);

    let value;
    if (type === 'risk') {
      // Risk tends to decrease over time (good)
      value = baseValue - (trend * 5) + noise + seasonal;
    } else {
      // Yield tends to increase over time (good)
      value = baseValue + (trend * 2) + noise + seasonal;
    }

    data.push({
      label: date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
      value: Math.max(0, value),
      date: date.toISOString()
    });
  }

  return data;
}
