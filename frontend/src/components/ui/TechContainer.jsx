import React from 'react';

/**
 * TechContainer - Unified dark tech-style page container
 */
export default function TechContainer({ children }) {
  return (
    <div style={{
      minHeight: '100vh',
      background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
      padding: '32px 0',
      position: 'relative',
      overflow: 'hidden'
    }}>
      {/* Animated background effects */}
      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundImage: `
          radial-gradient(circle at 20% 30%, rgba(34, 211, 238, 0.08) 0%, transparent 50%),
          radial-gradient(circle at 80% 70%, rgba(59, 130, 246, 0.08) 0%, transparent 50%)
        `,
        pointerEvents: 'none'
      }} />

      {/* Grid overlay */}
      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundImage: `
          linear-gradient(rgba(34, 211, 238, 0.02) 1px, transparent 1px),
          linear-gradient(90deg, rgba(34, 211, 238, 0.02) 1px, transparent 1px)
        `,
        backgroundSize: '50px 50px',
        opacity: 0.3,
        pointerEvents: 'none'
      }} />

      <div style={{ maxWidth: 1280, margin: '0 auto', padding: '0 24px', position: 'relative', zIndex: 1 }}>
        {children}
      </div>
    </div>
  );
}
