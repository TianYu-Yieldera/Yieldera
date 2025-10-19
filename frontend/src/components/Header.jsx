import React from "react";
import { useWallet } from "../web3/WalletContext";
import {
  Wallet,
  HelpCircle,
  LayoutDashboard,
  TrendingUp,
  Gem,
  Trophy,
  Gift,
  Shield,
  Activity
} from "lucide-react";
import { Link, useLocation } from "react-router-dom";

export default function Header(){
  const { address, connect, disconnect, switchToSepolia, chainId } = useWallet();
  const location = useLocation();
  const short = address ? `${address.slice(0,6)}...${address.slice(-4)}` : 'Connect Wallet';

  // 导航项配置
  const navItems = [
    {
      path: '/dashboard',
      label: '概览',
      icon: LayoutDashboard,
      color: '#9ca3af'
    },
    {
      path: '/vault',
      label: '理财金库',
      icon: TrendingUp,
      color: '#667eea',
      gradient: true,
      featured: true
    },
    {
      path: '/rwa-market',
      label: 'RWA 商城',
      icon: Gem,
      color: '#f093fb',
      gradient: true,
      featured: true
    },
    {
      path: '/leaderboard',
      label: '排行榜',
      icon: Trophy,
      color: '#9ca3af'
    },
    {
      path: '/airdrop',
      label: '空投',
      icon: Gift,
      color: '#9ca3af'
    }
  ];

  const adminItems = [
    {
      path: '/admin/airdrop',
      label: '空投管理',
      icon: Shield,
      color: '#F59E0B'
    },
    {
      path: '/status',
      label: '状态',
      icon: Activity,
      color: '#9ca3af'
    }
  ];

  const isActive = (path) => location.pathname === path;

  const NavLink = ({ item }) => {
    const Icon = item.icon;
    const active = isActive(item.path);

    if (item.featured) {
      // 核心功能 - 渐变背景
      return (
        <Link
          to={item.path}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '6px',
            padding: '8px 14px',
            borderRadius: '8px',
            fontSize: '14px',
            fontWeight: active ? '700' : '600',
            textDecoration: 'none',
            background: active
              ? `linear-gradient(135deg, ${item.color}dd 0%, ${item.color}99 100%)`
              : `linear-gradient(135deg, ${item.color}33 0%, ${item.color}22 100%)`,
            color: active ? '#fff' : item.color,
            border: active ? 'none' : `1px solid ${item.color}44`,
            transition: 'all 0.2s ease',
            transform: active ? 'translateY(-1px)' : 'none',
            boxShadow: active ? `0 4px 12px ${item.color}44` : 'none'
          }}
          onMouseEnter={(e) => {
            if (!active) {
              e.currentTarget.style.background = `linear-gradient(135deg, ${item.color}55 0%, ${item.color}33 100%)`;
              e.currentTarget.style.transform = 'translateY(-1px)';
              e.currentTarget.style.boxShadow = `0 4px 12px ${item.color}33`;
            }
          }}
          onMouseLeave={(e) => {
            if (!active) {
              e.currentTarget.style.background = `linear-gradient(135deg, ${item.color}33 0%, ${item.color}22 100%)`;
              e.currentTarget.style.transform = 'none';
              e.currentTarget.style.boxShadow = 'none';
            }
          }}
        >
          <Icon size={16} />
          {item.label}
        </Link>
      );
    }

    // 普通功能
    return (
      <Link
        to={item.path}
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: '6px',
          padding: '8px 12px',
          borderRadius: '6px',
          fontSize: '14px',
          fontWeight: active ? '600' : '500',
          textDecoration: 'none',
          color: active ? '#fff' : item.color,
          background: active ? 'rgba(255,255,255,0.1)' : 'transparent',
          transition: 'all 0.2s ease'
        }}
        onMouseEnter={(e) => {
          if (!active) {
            e.currentTarget.style.background = 'rgba(255,255,255,0.05)';
            e.currentTarget.style.color = '#fff';
          }
        }}
        onMouseLeave={(e) => {
          if (!active) {
            e.currentTarget.style.background = 'transparent';
            e.currentTarget.style.color = item.color;
          }
        }}
      >
        <Icon size={16} />
        {item.label}
      </Link>
    );
  };

  return (
    <header>
      <div className="row">
        {/* Logo */}
        <Link to="/" className="row" style={{gap:8}}>
          <img src="/pointfi-logo-mark.svg" width="32" height="32" style={{borderRadius:8}} alt="Yieldera"/>
          <div>
            <div style={{fontWeight:700, background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent'}}>Yieldera</div>
            <div className="muted" style={{fontSize: '11px', letterSpacing: '0.5px'}}>Enter the Yieldera</div>
          </div>
        </Link>

        {/* 主导航 */}
        <nav style={{marginLeft:'auto', display:'flex', alignItems:'center', gap:'6px'}}>
          {navItems.map((item) => (
            <NavLink key={item.path} item={item} />
          ))}

          {/* 分隔线 */}
          <div style={{width:'1px', height:'24px', background:'rgba(255,255,255,0.1)', margin:'0 8px'}} />

          {/* 管理功能 */}
          {adminItems.map((item) => (
            <NavLink key={item.path} item={item} />
          ))}
        </nav>

        {/* 右侧按钮 */}
        <Link
          to="/tutorial"
          className="btn"
          style={{
            marginLeft: 8,
            display: 'flex',
            alignItems: 'center',
            gap: 6,
            background: 'transparent',
            border: '1px solid rgba(255,255,255,0.2)',
            textDecoration: 'none',
            transition: 'all 0.2s ease'
          }}
          title="查看教程"
          onMouseEnter={(e) => {
            e.currentTarget.style.background = 'rgba(255,255,255,0.1)';
            e.currentTarget.style.borderColor = 'rgba(255,255,255,0.3)';
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.background = 'transparent';
            e.currentTarget.style.borderColor = 'rgba(255,255,255,0.2)';
          }}
        >
          <HelpCircle size={16} />
        </Link>

        {chainId && chainId !== '0xaa36a7' && (
          <button
            className="btn"
            style={{marginLeft:8, background:'#f59e0b', transition: 'all 0.2s ease'}}
            onClick={switchToSepolia}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = '#d97706';
              e.currentTarget.style.transform = 'translateY(-1px)';
              e.currentTarget.style.boxShadow = '0 4px 12px rgba(245, 158, 11, 0.4)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = '#f59e0b';
              e.currentTarget.style.transform = 'none';
              e.currentTarget.style.boxShadow = 'none';
            }}
          >
            切到 Sepolia
          </button>
        )}

        <button
          className="btn"
          style={{
            marginLeft: 8,
            display: 'flex',
            alignItems: 'center',
            gap: 6,
            background: address ? 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' : '#374151',
            transition: 'all 0.2s ease'
          }}
          onClick={address ? disconnect : connect}
          onMouseEnter={(e) => {
            if (address) {
              e.currentTarget.style.transform = 'translateY(-1px)';
              e.currentTarget.style.boxShadow = '0 4px 12px rgba(102, 126, 234, 0.4)';
            } else {
              e.currentTarget.style.background = '#4b5563';
            }
          }}
          onMouseLeave={(e) => {
            if (address) {
              e.currentTarget.style.transform = 'none';
              e.currentTarget.style.boxShadow = 'none';
            } else {
              e.currentTarget.style.background = '#374151';
            }
          }}
        >
          <Wallet size={16} /> {short}
        </button>
      </div>
    </header>
  )
}
