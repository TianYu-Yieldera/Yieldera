/**
 * Base Login Card
 *
 * 专为 Base 链设计的登录界面
 * 支持 Coinbase Smart Wallet 的多种登录方式
 */

import React, { useState } from 'react';
import { Wallet, CreditCard, Mail, Key, Zap, Shield, ArrowRight } from 'lucide-react';
import { useBaseSmartWallet } from '../web3/BaseSmartWalletProvider';

export function BaseLoginCard({ onLoginSuccess }) {
  const { connect, isConnected, address, getErrorMessage } = useBaseSmartWallet();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleConnect = async () => {
    setLoading(true);
    setError(null);

    try {
      await connect();

      // Success callback
      if (onLoginSuccess) {
        onLoginSuccess();
      }
    } catch (err) {
      console.error('Login failed:', err);
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  if (isConnected) {
    return (
      <div style={{
        padding: 24,
        background: 'linear-gradient(135deg, rgba(34, 197, 94, 0.1) 0%, rgba(16, 185, 129, 0.1) 100%)',
        border: '1px solid rgba(34, 197, 94, 0.3)',
        borderRadius: 12
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 12 }}>
          <Shield size={24} style={{ color: 'rgb(34, 197, 94)' }} />
          <div>
            <div style={{ fontSize: 16, fontWeight: 700, color: 'white' }}>
              Smart Wallet Connected
            </div>
            <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.7)', fontFamily: 'monospace' }}>
              {address?.slice(0, 6)}...{address?.slice(-4)}
            </div>
          </div>
        </div>
        <div style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.8)' }}>
          ✅ 您现在可以购买美国国债 (免 gas 费)
        </div>
      </div>
    );
  }

  return (
    <div style={{
      background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
      borderRadius: 16,
      padding: 32,
      border: '1px solid rgba(99, 102, 241, 0.2)',
      maxWidth: 480,
      margin: '0 auto'
    }}>
      {/* Header */}
      <div style={{ textAlign: 'center', marginBottom: 32 }}>
        <div style={{
          width: 64,
          height: 64,
          borderRadius: 16,
          background: 'linear-gradient(135deg, rgb(99, 102, 241) 0%, rgb(139, 92, 246) 100%)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          margin: '0 auto 16px'
        }}>
          <Wallet size={32} style={{ color: 'white' }} />
        </div>

        <h2 style={{
          fontSize: 24,
          fontWeight: 700,
          color: 'white',
          margin: '0 0 8px 0'
        }}>
          购买美国国债
        </h2>

        <p style={{
          fontSize: 14,
          color: 'rgba(203, 213, 225, 0.7)',
          margin: 0
        }}>
          使用 Base Smart Wallet 登��
        </p>
      </div>

      {/* Features */}
      <div style={{
        display: 'grid',
        gap: 12,
        marginBottom: 24
      }}>
        <FeatureItem
          icon={<Zap size={18} />}
          title="免 Gas 费"
          description="平台赞助所有交易费用"
        />
        <FeatureItem
          icon={<Shield size={18} />}
          title="安全可靠"
          description="Coinbase 官方 Smart Wallet"
        />
        <FeatureItem
          icon={<CreditCard size={18} />}
          title="支持信用卡"
          description="无需持有加密货币"
        />
      </div>

      {/* Login Button */}
      <button
        onClick={handleConnect}
        disabled={loading}
        style={{
          width: '100%',
          padding: '16px 24px',
          background: loading
            ? 'rgba(99, 102, 241, 0.5)'
            : 'linear-gradient(135deg, rgb(99, 102, 241) 0%, rgb(139, 92, 246) 100%)',
          border: 'none',
          borderRadius: 12,
          color: 'white',
          fontSize: 16,
          fontWeight: 600,
          cursor: loading ? 'not-allowed' : 'pointer',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          gap: 12,
          transition: 'all 0.3s',
          opacity: loading ? 0.7 : 1
        }}
        onMouseEnter={(e) => {
          if (!loading) {
            e.currentTarget.style.transform = 'translateY(-2px)';
            e.currentTarget.style.boxShadow = '0 8px 20px rgba(99, 102, 241, 0.4)';
          }
        }}
        onMouseLeave={(e) => {
          if (!loading) {
            e.currentTarget.style.transform = 'translateY(0)';
            e.currentTarget.style.boxShadow = 'none';
          }
        }}
      >
        {loading ? (
          <>
            <div style={{
              width: 20,
              height: 20,
              border: '2px solid rgba(255, 255, 255, 0.3)',
              borderTopColor: 'white',
              borderRadius: '50%',
              animation: 'spin 1s linear infinite'
            }} />
            <style>{`
              @keyframes spin {
                to { transform: rotate(360deg); }
              }
            `}</style>
            连接中...
          </>
        ) : (
          <>
            <Key size={20} />
            登录 Smart Wallet
            <ArrowRight size={20} />
          </>
        )}
      </button>

      {/* Error Message */}
      {error && (
        <div style={{
          marginTop: 16,
          padding: 12,
          background: 'rgba(239, 68, 68, 0.1)',
          border: '1px solid rgba(239, 68, 68, 0.3)',
          borderRadius: 8,
          fontSize: 13,
          color: 'rgb(239, 68, 68)'
        }}>
          ❌ {error}
        </div>
      )}

      {/* Info */}
      <div style={{
        marginTop: 24,
        padding: 16,
        background: 'rgba(34, 211, 238, 0.1)',
        border: '1px solid rgba(34, 211, 238, 0.2)',
        borderRadius: 10,
        fontSize: 13,
        color: 'rgba(203, 213, 225, 0.8)',
        lineHeight: 1.6
      }}>
        <strong style={{ color: 'rgb(34, 211, 238)' }}>💡 什么是 Smart Wallet?</strong>
        <br />
        无需助记词，使用 Passkey、Google 或邮箱即可登录。
        所有交易费用由平台赞助，体验如同 Web2 应用。
      </div>

      {/* Restricted Notice */}
      <div style={{
        marginTop: 16,
        padding: 12,
        background: 'rgba(245, 158, 11, 0.1)',
        border: '1px solid rgba(245, 158, 11, 0.3)',
        borderRadius: 8,
        fontSize: 12,
        color: 'rgb(245, 158, 11)',
        textAlign: 'center'
      }}>
        ⚠️ 国债购买仅支持 Base 链 Smart Wallet
      </div>
    </div>
  );
}

// ============================================================================
// Feature Item Component
// ============================================================================

function FeatureItem({ icon, title, description }) {
  return (
    <div style={{
      display: 'flex',
      alignItems: 'flex-start',
      gap: 12,
      padding: 12,
      background: 'rgba(255, 255, 255, 0.03)',
      borderRadius: 8,
      border: '1px solid rgba(255, 255, 255, 0.08)'
    }}>
      <div style={{
        width: 36,
        height: 36,
        borderRadius: 8,
        background: 'rgba(99, 102, 241, 0.15)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        color: 'rgb(99, 102, 241)',
        flexShrink: 0
      }}>
        {icon}
      </div>
      <div style={{ flex: 1 }}>
        <div style={{
          fontSize: 14,
          fontWeight: 600,
          color: 'white',
          marginBottom: 4
        }}>
          {title}
        </div>
        <div style={{
          fontSize: 12,
          color: 'rgba(203, 213, 225, 0.7)'
        }}>
          {description}
        </div>
      </div>
    </div>
  );
}

export default BaseLoginCard;
