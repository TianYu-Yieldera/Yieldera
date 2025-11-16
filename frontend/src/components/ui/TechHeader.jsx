import React from 'react';

/**
 * TechHeader - Unified dark tech-style page header
 */
export default function TechHeader({ icon: Icon, title, subtitle, children }) {
  return (
    <div style={{ marginBottom: 32 }}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 8 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          {Icon && (
            <div style={{
              width: 48,
              height: 48,
              borderRadius: 12,
              background: 'rgba(34, 211, 238, 0.15)',
              border: '1px solid rgba(34, 211, 238, 0.3)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
              <Icon style={{ width: 28, height: 28, color: 'rgb(34, 211, 238)' }} />
            </div>
          )}
          <div>
            <h1 style={{
              fontSize: 30,
              fontWeight: 700,
              color: 'white',
              margin: 0,
              textShadow: '0 0 30px rgba(34, 211, 238, 0.3)'
            }}>
              {title}
            </h1>
            {subtitle && (
              <p style={{
                color: 'rgba(203, 213, 225, 0.7)',
                margin: '4px 0 0 0',
                fontSize: 15
              }}>
                {subtitle}
              </p>
            )}
          </div>
        </div>
        {children}
      </div>
    </div>
  );
}
