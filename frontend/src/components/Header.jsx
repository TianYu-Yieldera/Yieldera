import React, { useState, useEffect, useRef } from "react";
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
  Activity,
  Bell,
  X,
  CheckCircle,
  AlertCircle,
  Info
} from "lucide-react";
import { Link, useLocation } from "react-router-dom";
import notificationsService from "../services/notificationsService";

export default function Header(){
  const { address, connect, disconnect, switchToSepolia, chainId } = useWallet();
  const location = useLocation();
  const short = address ? `${address.slice(0,6)}...${address.slice(-4)}` : 'Connect Wallet';

  // Notifications state
  const [showNotifications, setShowNotifications] = useState(false);
  const [notifications, setNotifications] = useState([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const notificationRef = useRef(null);

  // Load notifications when user connects
  useEffect(() => {
    if (address) {
      loadNotifications();
    } else {
      setNotifications([]);
      setUnreadCount(0);
    }
  }, [address]);

  // Close notification panel when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (notificationRef.current && !notificationRef.current.contains(event.target)) {
        setShowNotifications(false);
      }
    };

    if (showNotifications) {
      document.addEventListener('mousedown', handleClickOutside);
      return () => document.removeEventListener('mousedown', handleClickOutside);
    }
  }, [showNotifications]);

  const loadNotifications = async () => {
    try {
      const data = await notificationsService.getUserNotifications(address, 10);
      setNotifications(data.notifications || []);
      setUnreadCount(notificationsService.getUnreadCount(data.notifications || []));
    } catch (error) {
      console.error('Failed to load notifications:', error);
      // Use mock data on error
      setNotifications([
        {
          id: 1,
          type: 'yield',
          severity: 'low',
          title: 'Yield Distributed',
          message: 'You earned $45.32 in treasury yield',
          timestamp: new Date(Date.now() - 3600000).toISOString(),
          read: false
        },
        {
          id: 2,
          type: 'risk',
          severity: 'medium',
          title: 'Risk Alert',
          message: 'Portfolio risk score increased to 42',
          timestamp: new Date(Date.now() - 7200000).toISOString(),
          read: false
        }
      ]);
      setUnreadCount(2);
    }
  };

  const handleMarkAsRead = async (notification) => {
    try {
      await notificationsService.markAsRead(address, notification.timestamp);
      setNotifications(prev =>
        prev.map(n => n.id === notification.id ? { ...n, read: true } : n)
      );
      setUnreadCount(prev => Math.max(0, prev - 1));
    } catch (error) {
      console.error('Failed to mark as read:', error);
    }
  };

  const getNotificationIcon = (type) => {
    switch(type) {
      case 'success': case 'yield': return CheckCircle;
      case 'risk': case 'alert': return AlertCircle;
      default: return Info;
    }
  };

  // Navigation items - core features only
  const navItems = [
    {
      path: '/dashboard',
      label: 'Portfolio',
      icon: LayoutDashboard,
      color: '#9ca3af'
    },
    {
      path: '/vault',
      label: 'Vault',
      icon: TrendingUp,
      color: '#667eea',
      gradient: true,
      featured: true
    },
    {
      path: '/treasury',
      label: 'Treasury',
      icon: Gem,
      color: '#1e40af',
      gradient: true,
      featured: true
    },
    {
      path: '/monitoring',
      label: 'Monitoring',
      icon: Activity,
      color: '#10b981',
      gradient: true,
      featured: true
    }
  ];

  const adminItems = [];

  const isActive = (path) => location.pathname === path;

  const NavLink = ({ item }) => {
    const Icon = item.icon;
    const active = isActive(item.path);

    return (
      <Link
        to={item.path}
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: '8px',
          padding: '8px 16px',
          borderRadius: '6px',
          fontSize: '14px',
          fontWeight: active ? '600' : '500',
          textDecoration: 'none',
          color: active ? 'rgb(29, 78, 216)' : 'rgb(71, 85, 105)',
          background: active ? 'rgb(239, 246, 255)' : 'transparent',
          border: active ? '1px solid rgb(191, 219, 254)' : '1px solid transparent',
          transition: 'all 0.2s ease'
        }}
        onMouseEnter={(e) => {
          if (!active) {
            e.currentTarget.style.background = 'rgb(248, 250, 252)';
            e.currentTarget.style.color = 'rgb(15, 23, 42)';
          }
        }}
        onMouseLeave={(e) => {
          if (!active) {
            e.currentTarget.style.background = 'transparent';
            e.currentTarget.style.color = 'rgb(71, 85, 105)';
          }
        }}
      >
        <Icon size={16} />
        {item.label}
      </Link>
    );
  };

  return (
    <header style={{
      background: 'rgba(255, 255, 255, 0.95)',
      backdropFilter: 'blur(10px)',
      borderBottom: '1px solid rgb(226, 232, 240)',
      boxShadow: '0 1px 3px rgba(0, 0, 0, 0.05)'
    }}>
      <div className="row" style={{ justifyContent: 'space-between', width: '100%' }}>
        {/* Logo */}
        <Link to="/" className="row" style={{gap:8, textDecoration: 'none'}}>
          <img src="/pointfi-logo-mark.svg" width="32" height="32" style={{borderRadius:8}} alt="Yieldera"/>
          <div>
            <div style={{fontWeight:700, fontSize: '18px', color: 'rgb(15, 23, 42)'}}>Yieldera</div>
            <div style={{fontSize: '11px', letterSpacing: '0.5px', color: 'rgb(100, 116, 139)'}}>Professional DeFi Platform</div>
          </div>
        </Link>

        {/* Navigation - centered */}
        <nav style={{display:'flex', alignItems:'center', gap:'6px', position: 'absolute', left: '50%', transform: 'translateX(-50%)'}}>
          {navItems.map((item) => (
            <NavLink key={item.path} item={item} />
          ))}
        </nav>

        {/* Right side buttons */}
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <Link
            to="/tutorial"
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: 36,
              height: 36,
              background: 'white',
              border: '1px solid rgb(226, 232, 240)',
              borderRadius: '6px',
              color: 'rgb(71, 85, 105)',
              textDecoration: 'none',
              transition: 'all 0.2s ease',
              cursor: 'pointer'
            }}
            title="View Tutorial"
            onMouseEnter={(e) => {
              e.currentTarget.style.background = 'rgb(248, 250, 252)';
              e.currentTarget.style.borderColor = 'rgb(203, 213, 225)';
              e.currentTarget.style.color = 'rgb(15, 23, 42)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = 'white';
              e.currentTarget.style.borderColor = 'rgb(226, 232, 240)';
              e.currentTarget.style.color = 'rgb(71, 85, 105)';
            }}
          >
            <HelpCircle size={18} />
          </Link>

          {/* Notifications Bell - Always visible for demo */}
          <div style={{ position: 'relative' }} ref={notificationRef}>
              <button
                onClick={() => setShowNotifications(!showNotifications)}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  width: 36,
                  height: 36,
                  background: showNotifications ? 'rgb(239, 246, 255)' : 'white',
                  border: showNotifications ? '1px solid rgb(191, 219, 254)' : '1px solid rgb(226, 232, 240)',
                  borderRadius: '6px',
                  color: showNotifications ? 'rgb(29, 78, 216)' : 'rgb(71, 85, 105)',
                  transition: 'all 0.2s ease',
                  cursor: 'pointer',
                  position: 'relative'
                }}
                title="Notifications"
                onMouseEnter={(e) => {
                  if (!showNotifications) {
                    e.currentTarget.style.background = 'rgb(248, 250, 252)';
                    e.currentTarget.style.borderColor = 'rgb(203, 213, 225)';
                    e.currentTarget.style.color = 'rgb(15, 23, 42)';
                  }
                }}
                onMouseLeave={(e) => {
                  if (!showNotifications) {
                    e.currentTarget.style.background = 'white';
                    e.currentTarget.style.borderColor = 'rgb(226, 232, 240)';
                    e.currentTarget.style.color = 'rgb(71, 85, 105)';
                  }
                }}
              >
                <Bell size={18} />
                {unreadCount > 0 && (
                  <span style={{
                    position: 'absolute',
                    top: -4,
                    right: -4,
                    background: 'rgb(220, 38, 38)',
                    color: 'white',
                    borderRadius: '50%',
                    width: 18,
                    height: 18,
                    fontSize: 10,
                    fontWeight: 700,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    border: '2px solid white'
                  }}>
                    {unreadCount}
                  </span>
                )}
              </button>

              {/* Notification Panel */}
              {showNotifications && (
                <div style={{
                  position: 'absolute',
                  top: 'calc(100% + 8px)',
                  right: 0,
                  width: 380,
                  maxHeight: 500,
                  background: 'white',
                  border: '1px solid rgb(226, 232, 240)',
                  borderRadius: 12,
                  boxShadow: '0 10px 25px rgba(0,0,0,0.1), 0 0 0 1px rgba(0,0,0,0.05)',
                  overflow: 'hidden',
                  zIndex: 1000
                }}>
                  {/* Header */}
                  <div style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    padding: '16px 20px',
                    borderBottom: '1px solid rgb(226, 232, 240)',
                    background: 'rgb(248, 250, 252)'
                  }}>
                    <div>
                      <h3 style={{ fontSize: 16, fontWeight: 700, color: 'rgb(15, 23, 42)', margin: 0 }}>
                        Notifications
                      </h3>
                      <p style={{ fontSize: 12, color: 'rgb(100, 116, 139)', margin: '2px 0 0 0' }}>
                        {unreadCount} unread
                      </p>
                    </div>
                    <button
                      onClick={() => setShowNotifications(false)}
                      style={{
                        background: 'transparent',
                        border: 'none',
                        cursor: 'pointer',
                        color: 'rgb(100, 116, 139)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        padding: 4
                      }}
                    >
                      <X size={18} />
                    </button>
                  </div>

                  {/* Notifications List */}
                  <div style={{
                    maxHeight: 400,
                    overflowY: 'auto'
                  }}>
                    {notifications.length === 0 ? (
                      <div style={{
                        padding: 48,
                        textAlign: 'center',
                        color: 'rgb(148, 163, 184)'
                      }}>
                        <Bell size={32} style={{ margin: '0 auto 12px', opacity: 0.3 }} />
                        <p style={{ fontSize: 14 }}>No notifications yet</p>
                      </div>
                    ) : (
                      notifications.map((notification) => {
                        const Icon = getNotificationIcon(notification.type);
                        const colors = notificationsService.getNotificationColor(notification.severity);

                        return (
                          <div
                            key={notification.id}
                            onClick={() => !notification.read && handleMarkAsRead(notification)}
                            style={{
                              padding: 16,
                              borderBottom: '1px solid rgb(243, 244, 246)',
                              background: notification.read ? 'white' : 'rgb(239, 246, 255)',
                              cursor: notification.read ? 'default' : 'pointer',
                              transition: 'background 0.2s'
                            }}
                            onMouseEnter={(e) => {
                              if (!notification.read) {
                                e.currentTarget.style.background = 'rgb(219, 234, 254)';
                              }
                            }}
                            onMouseLeave={(e) => {
                              if (!notification.read) {
                                e.currentTarget.style.background = 'rgb(239, 246, 255)';
                              }
                            }}
                          >
                            <div style={{ display: 'flex', gap: 12 }}>
                              <div style={{
                                width: 40,
                                height: 40,
                                borderRadius: 8,
                                background: colors.bg,
                                border: `1px solid ${colors.border}`,
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                flexShrink: 0
                              }}>
                                <Icon size={18} style={{ color: colors.text }} />
                              </div>
                              <div style={{ flex: 1 }}>
                                <div style={{
                                  display: 'flex',
                                  alignItems: 'center',
                                  justifyContent: 'space-between',
                                  marginBottom: 4
                                }}>
                                  <h4 style={{
                                    fontSize: 14,
                                    fontWeight: 600,
                                    color: 'rgb(15, 23, 42)',
                                    margin: 0
                                  }}>
                                    {notification.title}
                                  </h4>
                                  {!notification.read && (
                                    <div style={{
                                      width: 8,
                                      height: 8,
                                      borderRadius: '50%',
                                      background: 'rgb(59, 130, 246)',
                                      flexShrink: 0
                                    }} />
                                  )}
                                </div>
                                <p style={{
                                  fontSize: 13,
                                  color: 'rgb(71, 85, 105)',
                                  margin: '0 0 4px 0',
                                  lineHeight: 1.4
                                }}>
                                  {notification.message}
                                </p>
                                <p style={{
                                  fontSize: 11,
                                  color: 'rgb(148, 163, 184)',
                                  margin: 0
                                }}>
                                  {notificationsService.formatTimeAgo(notification.timestamp)}
                                </p>
                              </div>
                            </div>
                          </div>
                        );
                      })
                    )}
                  </div>
                </div>
              )}
            </div>

          {chainId && chainId !== '0xaa36a7' && (
            <button
              style={{
                padding: '8px 16px',
                background: 'rgb(254, 243, 199)',
                color: 'rgb(146, 64, 14)',
                border: '1px solid rgb(253, 224, 71)',
                borderRadius: '6px',
                fontSize: '14px',
                fontWeight: '500',
                cursor: 'pointer',
                transition: 'all 0.2s ease'
              }}
              onClick={switchToSepolia}
              onMouseEnter={(e) => {
                e.currentTarget.style.background = 'rgb(253, 230, 138)';
                e.currentTarget.style.borderColor = 'rgb(251, 191, 36)';
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.background = 'rgb(254, 243, 199)';
                e.currentTarget.style.borderColor = 'rgb(253, 224, 71)';
              }}
            >
              Switch to Sepolia
            </button>
          )}

          <button
            style={{
              padding: '8px 16px',
              display: 'flex',
              alignItems: 'center',
              gap: 8,
              background: address ? 'white' : 'rgb(15, 23, 42)',
              color: address ? 'rgb(15, 23, 42)' : 'white',
              border: address ? '1px solid rgb(203, 213, 225)' : 'none',
              borderRadius: '6px',
              fontSize: '14px',
              fontWeight: '500',
              cursor: 'pointer',
              transition: 'all 0.2s ease'
            }}
            onClick={address ? disconnect : connect}
            onMouseEnter={(e) => {
              if (address) {
                e.currentTarget.style.background = 'rgb(248, 250, 252)';
                e.currentTarget.style.borderColor = 'rgb(148, 163, 184)';
              } else {
                e.currentTarget.style.background = 'rgb(30, 41, 59)';
              }
            }}
            onMouseLeave={(e) => {
              if (address) {
                e.currentTarget.style.background = 'white';
                e.currentTarget.style.borderColor = 'rgb(203, 213, 225)';
              } else {
                e.currentTarget.style.background = 'rgb(15, 23, 42)';
              }
            }}
          >
            <Wallet size={16} /> {short}
          </button>
        </div>
      </div>
    </header>
  )
}
