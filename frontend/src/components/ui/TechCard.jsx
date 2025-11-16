import React from 'react';

/**
 * TechCard - Unified dark tech-style card component
 * Used across all pages for consistent design
 */
export default function TechCard({
  icon: Icon,
  title,
  value,
  subtitle,
  trend,
  iconColor = 'rgb(34, 211, 238)',
  children,
  onClick,
  className = ''
}) {
  return (
    <div
      style={{
        background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
        borderRadius: 12,
        padding: 24,
        border: '1px solid rgba(34, 211, 238, 0.2)',
        transition: 'all 0.3s ease',
        cursor: onClick ? 'pointer' : 'default',
        position: 'relative',
        overflow: 'hidden'
      }}
      className={className}
      onClick={onClick}
      onMouseEnter={(e) => {
        e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.4)';
        e.currentTarget.style.transform = 'translateY(-4px)';
        e.currentTarget.style.boxShadow = '0 8px 24px rgba(0,0,0,0.3)';
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.2)';
        e.currentTarget.style.transform = 'translateY(0)';
        e.currentTarget.style.boxShadow = 'none';
      }}
    >
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

      {/* Gradient accent */}
      <div style={{
        position: 'absolute',
        top: -50,
        right: -50,
        width: 150,
        height: 150,
        background: 'radial-gradient(circle, rgba(34, 211, 238, 0.15) 0%, transparent 70%)',
        opacity: 0.6
      }} />

      <div style={{ position: 'relative', zIndex: 1 }}>
        {/* Custom content */}
        {children ? children : (
          <>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
              <div style={{
                fontSize: 12,
                fontWeight: 600,
                color: 'rgba(203, 213, 225, 0.7)',
                textTransform: 'uppercase',
                letterSpacing: 1
              }}>
                {title}
              </div>
              {Icon && (
                <div style={{
                  width: 36,
                  height: 36,
                  borderRadius: 8,
                  background: 'rgba(34, 211, 238, 0.15)',
                  border: '1px solid rgba(34, 211, 238, 0.3)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center'
                }}>
                  <Icon style={{ width: 18, height: 18, color: iconColor }} />
                </div>
              )}
            </div>
            <div style={{
              fontSize: 28,
              fontWeight: 700,
              color: 'white',
              marginBottom: 8,
              textShadow: '0 0 20px rgba(255, 255, 255, 0.1)'
            }}>
              {value}
            </div>
            {subtitle && (
              <div style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.8)', display: 'flex', alignItems: 'center', gap: 6 }}>
                {trend !== undefined && (
                  <span style={{
                    color: trend > 0 ? 'rgb(34, 197, 94)' : trend < 0 ? 'rgb(239, 68, 68)' : 'rgba(203, 213, 225, 0.6)',
                    fontWeight: 700,
                    fontSize: 14
                  }}>
                    {trend > 0 ? '↑' : trend < 0 ? '↓' : '•'} {Math.abs(trend)}%
                  </span>
                )}
                <span>{subtitle}</span>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
