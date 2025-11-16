import React from 'react';

/**
 * TechButton - Unified dark tech-style button
 */
export default function TechButton({
  children,
  onClick,
  variant = 'primary',
  disabled = false,
  loading = false,
  icon: Icon,
  fullWidth = false
}) {
  const getVariantStyles = () => {
    switch (variant) {
      case 'primary':
        return {
          background: 'linear-gradient(135deg, rgb(34, 211, 238) 0%, rgb(59, 130, 246) 100%)',
          color: 'rgb(15, 23, 42)',
          border: 'none',
          hoverBg: 'linear-gradient(135deg, rgb(34, 211, 238) 0%, rgb(29, 78, 216) 100%)'
        };
      case 'secondary':
        return {
          background: 'rgba(34, 211, 238, 0.1)',
          color: 'rgb(34, 211, 238)',
          border: '1px solid rgba(34, 211, 238, 0.3)',
          hoverBg: 'rgba(34, 211, 238, 0.2)'
        };
      case 'danger':
        return {
          background: 'linear-gradient(135deg, rgb(239, 68, 68) 0%, rgb(220, 38, 38) 100%)',
          color: 'white',
          border: 'none',
          hoverBg: 'linear-gradient(135deg, rgb(220, 38, 38) 0%, rgb(185, 28, 28) 100%)'
        };
      default:
        return {
          background: 'rgba(100, 116, 139, 0.2)',
          color: 'rgba(203, 213, 225, 0.9)',
          border: '1px solid rgba(100, 116, 139, 0.3)',
          hoverBg: 'rgba(100, 116, 139, 0.3)'
        };
    }
  };

  const styles = getVariantStyles();

  return (
    <button
      onClick={onClick}
      disabled={disabled || loading}
      style={{
        background: styles.background,
        color: styles.color,
        border: styles.border,
        padding: '12px 24px',
        borderRadius: 8,
        fontSize: 14,
        fontWeight: 600,
        cursor: disabled || loading ? 'not-allowed' : 'pointer',
        opacity: disabled || loading ? 0.5 : 1,
        transition: 'all 0.3s ease',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: 8,
        width: fullWidth ? '100%' : 'auto',
        textTransform: 'uppercase',
        letterSpacing: 0.5,
        boxShadow: variant === 'primary' ? '0 4px 12px rgba(34, 211, 238, 0.3)' : 'none'
      }}
      onMouseEnter={(e) => {
        if (!disabled && !loading) {
          e.currentTarget.style.background = styles.hoverBg;
          e.currentTarget.style.transform = 'translateY(-2px)';
          if (variant === 'primary') {
            e.currentTarget.style.boxShadow = '0 6px 16px rgba(34, 211, 238, 0.4)';
          }
        }
      }}
      onMouseLeave={(e) => {
        if (!disabled && !loading) {
          e.currentTarget.style.background = styles.background;
          e.currentTarget.style.transform = 'translateY(0)';
          if (variant === 'primary') {
            e.currentTarget.style.boxShadow = '0 4px 12px rgba(34, 211, 238, 0.3)';
          }
        }
      }}
    >
      {loading ? (
        <>
          <div style={{
            width: 16,
            height: 16,
            border: '2px solid currentColor',
            borderTopColor: 'transparent',
            borderRadius: '50%',
            animation: 'spin 1s linear infinite'
          }} />
          Processing...
        </>
      ) : (
        <>
          {Icon && <Icon style={{ width: 18, height: 18 }} />}
          {children}
        </>
      )}
    </button>
  );
}
