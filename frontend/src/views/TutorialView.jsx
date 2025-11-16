import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Play, Home, Monitor, TrendingUp, Gem, Shield, Youtube, Video, ExternalLink, Check, BookOpen, Zap } from 'lucide-react';
import { useDemoMode } from '../web3/DemoModeContext';
import TechContainer from "../components/ui/TechContainer";
import TechHeader from "../components/ui/TechHeader";
import TechCard from "../components/ui/TechCard";
import TechButton from "../components/ui/TechButton";

// Video configuration - modify here to change video source
const VIDEO_CONFIG = {
  type: 'placeholder', // 'youtube' | 'local' | 'placeholder'
  youtubeId: '', // YouTube video ID (example: 'dQw4w9WgXcQ')
  localPath: '/videos/demo.mp4', // Local video path
  posterImage: '/pointfi-logo-mark.svg' // Video poster image
};

export default function TutorialView() {
  const navigate = useNavigate();
  const { demoMode, enableDemoMode, disableDemoMode } = useDemoMode();
  const [videoExpanded, setVideoExpanded] = useState(true);

  const handleToggleDemoMode = async () => {
    try {
      if (!demoMode) {
        await enableDemoMode();
      } else {
        await disableDemoMode();
      }
    } catch (err) {
      console.error('Failed to toggle demo mode', err);
      alert('Failed to initialize demo mode, please try again later');
    }
  };

  return (
    <TechContainer>
      <TechHeader
        icon={BookOpen}
        title="Platform Demo & Quick Start"
        subtitle="Watch the demo video to understand Yieldera's core features, or enable demo mode to start exploring"
      >
        <TechButton
          onClick={() => navigate(-1)}
          variant="secondary"
          icon={Home}
        >
          Back to Home
        </TechButton>
      </TechHeader>

      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: 24 }}>
        {/* Left Column - Video Section */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
          {/* Main Video Player */}
          <TechCard>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 20 }}>
              <h2 style={{ fontSize: 20, fontWeight: 600, color: 'white', margin: 0 }}>
                Yieldera Platform Demo
              </h2>
              {videoExpanded && (
                <button
                  onClick={() => setVideoExpanded(false)}
                  style={{
                    padding: '8px 16px',
                    background: 'rgba(15, 23, 42, 0.5)',
                    border: '1px solid rgba(34, 211, 238, 0.3)',
                    borderRadius: 6,
                    fontSize: 13,
                    color: 'rgba(203, 213, 225, 0.8)',
                    cursor: 'pointer',
                    transition: 'all 0.3s'
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.6)';
                    e.currentTarget.style.color = 'rgb(34, 211, 238)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.3)';
                    e.currentTarget.style.color = 'rgba(203, 213, 225, 0.8)';
                  }}
                >
                  Collapse Video
                </button>
              )}
            </div>

            {videoExpanded ? (
              <>
                {/* YouTube Video */}
                {VIDEO_CONFIG.type === 'youtube' && VIDEO_CONFIG.youtubeId && (
                  <div style={{ position: 'relative', paddingBottom: '56.25%', height: 0, borderRadius: 12, overflow: 'hidden', border: '1px solid rgba(34, 211, 238, 0.3)' }}>
                    <iframe
                      style={{ position: 'absolute', top: 0, left: 0, width: '100%', height: '100%' }}
                      src={`https://www.youtube.com/embed/${VIDEO_CONFIG.youtubeId}?rel=0&modestbranding=1`}
                      title="Yieldera Demo"
                      frameBorder="0"
                      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                      allowFullScreen
                    />
                  </div>
                )}

                {/* Local Video */}
                {VIDEO_CONFIG.type === 'local' && VIDEO_CONFIG.localPath && (
                  <video
                    controls
                    poster={VIDEO_CONFIG.posterImage}
                    style={{ width: '100%', borderRadius: 12, border: '1px solid rgba(34, 211, 238, 0.3)' }}
                  >
                    <source src={VIDEO_CONFIG.localPath} type="video/mp4" />
                    Your browser does not support video playback
                  </video>
                )}

                {/* Placeholder */}
                {VIDEO_CONFIG.type === 'placeholder' && (
                  <div style={{
                    background: 'linear-gradient(135deg, rgba(15, 23, 42, 0.8) 0%, rgba(30, 41, 59, 0.8) 100%)',
                    borderRadius: 12,
                    padding: '80px 32px',
                    textAlign: 'center',
                    border: '2px dashed rgba(34, 211, 238, 0.3)',
                    position: 'relative'
                  }}>
                    <div style={{
                      width: 120,
                      height: 120,
                      margin: '0 auto 24px',
                      background: 'rgba(34, 211, 238, 0.15)',
                      borderRadius: '50%',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      border: '4px solid rgba(34, 211, 238, 0.4)',
                      boxShadow: '0 8px 24px rgba(34, 211, 238, 0.3)'
                    }}>
                      <Play size={48} style={{ color: 'rgb(34, 211, 238)', marginLeft: 6 }} />
                    </div>
                    <h3 style={{ color: 'rgb(34, 211, 238)', fontSize: 24, marginBottom: 12, fontWeight: 700, textShadow: '0 0 20px rgba(34, 211, 238, 0.5)' }}>
                      Demo Video Coming Soon
                    </h3>
                    <p style={{ color: 'rgba(203, 213, 225, 0.8)', fontSize: 16, lineHeight: 1.6, marginBottom: 24, maxWidth: 500, margin: '0 auto 24px' }}>
                      We're creating a professional system demo video!<br/>
                      Showcasing DeFi yield optimization, US Treasury bonds, and AI risk management.
                    </p>
                    <div style={{
                      display: 'inline-flex',
                      alignItems: 'center',
                      gap: 12,
                      padding: '12px 24px',
                      background: 'rgba(15, 23, 42, 0.5)',
                      borderRadius: 8,
                      border: '1px solid rgba(34, 211, 238, 0.3)'
                    }}>
                      <div style={{ display: 'flex', gap: 16, fontSize: 14, color: 'rgba(203, 213, 225, 0.7)' }}>
                        <span><strong style={{ color: 'rgb(34, 211, 238)' }}>Duration:</strong> 5-8 min</span>
                        <span style={{ color: 'rgba(100, 116, 139, 0.5)' }}>|</span>
                        <span><strong style={{ color: 'rgb(34, 211, 238)' }}>Language:</strong> English</span>
                      </div>
                    </div>
                  </div>
                )}

                {/* Video Info */}
                <div style={{
                  marginTop: 20,
                  padding: 20,
                  background: 'rgba(34, 211, 238, 0.1)',
                  borderRadius: 8,
                  border: '1px solid rgba(34, 211, 238, 0.3)'
                }}>
                  <h4 style={{ fontSize: 15, fontWeight: 600, color: 'rgb(34, 211, 238)', margin: '0 0 12px 0' }}>
                    Video Content Overview
                  </h4>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                    {[
                      'Product Positioning & Core Value',
                      'DeFi Yield Aggregation Demo',
                      'US Treasury Bond Investment Flow',
                      'AI Smart Risk Management System',
                      'Complete User Journey Walkthrough'
                    ].map((item, i) => (
                      <div key={i} style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 14, color: 'rgba(203, 213, 225, 0.8)' }}>
                        <Check size={16} style={{ color: 'rgb(34, 197, 94)', flexShrink: 0 }} />
                        <span>{item}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </>
            ) : (
              <button
                onClick={() => setVideoExpanded(true)}
                style={{
                  width: '100%',
                  padding: '48px 32px',
                  background: 'rgba(15, 23, 42, 0.5)',
                  border: '2px dashed rgba(34, 211, 238, 0.3)',
                  borderRadius: 12,
                  cursor: 'pointer',
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  gap: 12,
                  transition: 'all 0.3s'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.6)';
                  e.currentTarget.style.background = 'rgba(15, 23, 42, 0.7)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.3)';
                  e.currentTarget.style.background = 'rgba(15, 23, 42, 0.5)';
                }}
              >
                <Play size={32} style={{ color: 'rgb(34, 211, 238)' }} />
                <span style={{ fontSize: 15, color: 'rgba(203, 213, 225, 0.8)' }}>Click to expand video</span>
              </button>
            )}
          </TechCard>

          {/* Video Links */}
          <TechCard>
            <h3 style={{ fontSize: 18, fontWeight: 600, color: 'white', margin: '0 0 16px 0' }}>
              Related Links
            </h3>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
              <a
                href="#"
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  padding: 16,
                  background: 'rgba(15, 23, 42, 0.5)',
                  border: '1px solid rgba(34, 211, 238, 0.2)',
                  borderRadius: 8,
                  textDecoration: 'none',
                  color: 'rgba(203, 213, 225, 0.8)',
                  transition: 'all 0.3s'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.borderColor = 'rgba(239, 68, 68, 0.5)';
                  e.currentTarget.style.background = 'rgba(239, 68, 68, 0.1)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.2)';
                  e.currentTarget.style.background = 'rgba(15, 23, 42, 0.5)';
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                  <Youtube size={20} style={{ color: 'rgb(239, 68, 68)' }} />
                  <div>
                    <div style={{ fontSize: 14, fontWeight: 500, color: 'white' }}>YouTube Video</div>
                    <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)' }}>HD quality | Subtitles available</div>
                  </div>
                </div>
                <ExternalLink size={16} />
              </a>
              <a
                href="#"
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  padding: 16,
                  background: 'rgba(15, 23, 42, 0.5)',
                  border: '1px solid rgba(34, 211, 238, 0.2)',
                  borderRadius: 8,
                  textDecoration: 'none',
                  color: 'rgba(203, 213, 225, 0.8)',
                  transition: 'all 0.3s'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.borderColor = 'rgba(59, 130, 246, 0.5)';
                  e.currentTarget.style.background = 'rgba(59, 130, 246, 0.1)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.2)';
                  e.currentTarget.style.background = 'rgba(15, 23, 42, 0.5)';
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                  <Video size={20} style={{ color: 'rgb(59, 130, 246)' }} />
                  <div>
                    <div style={{ fontSize: 14, fontWeight: 500, color: 'white' }}>Documentation</div>
                    <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)' }}>Technical documentation | PDF</div>
                  </div>
                </div>
                <ExternalLink size={16} />
              </a>
            </div>
          </TechCard>
        </div>

        {/* Right Column - Quick Start */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
          {/* Demo Mode Card */}
          <TechCard style={{
            background: demoMode
              ? 'linear-gradient(135deg, rgba(34, 197, 94, 0.2) 0%, rgba(22, 163, 74, 0.2) 100%)'
              : 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
            border: demoMode ? '2px solid rgba(34, 197, 94, 0.5)' : '1px solid rgba(34, 211, 238, 0.2)'
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 16 }}>
              <div style={{
                width: 48,
                height: 48,
                borderRadius: 12,
                background: demoMode ? 'rgb(34, 197, 94)' : 'rgb(34, 211, 238)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                boxShadow: demoMode ? '0 0 20px rgba(34, 197, 94, 0.5)' : '0 0 20px rgba(34, 211, 238, 0.5)'
              }}>
                <Monitor size={24} color="white" />
              </div>
              <div style={{ flex: 1 }}>
                <h3 style={{ fontSize: 16, fontWeight: 600, color: 'white', margin: 0 }}>
                  {demoMode ? 'Demo Mode Enabled' : 'Enable Demo Mode'}
                </h3>
                <p style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.7)', margin: '4px 0 0 0' }}>
                  {demoMode ? '10,000 test tokens available' : 'Test safely with virtual tokens'}
                </p>
              </div>
            </div>

            {demoMode && (
              <div style={{
                padding: 16,
                background: 'rgba(34, 197, 94, 0.1)',
                borderRadius: 8,
                marginBottom: 16,
                border: '1px solid rgba(34, 197, 94, 0.3)'
              }}>
                <div style={{ fontSize: 13, color: 'rgb(34, 197, 94)', fontWeight: 600, marginBottom: 8 }}>
                  You can safely:
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
                  {[
                    'Deposit test tokens into DeFi pools',
                    'Purchase virtual US Treasury bonds',
                    'View AI risk assessments'
                  ].map((text, i) => (
                    <div key={i} style={{ fontSize: 13, color: 'rgba(203, 213, 225, 0.8)', display: 'flex', alignItems: 'center', gap: 6 }}>
                      <Check size={14} style={{ color: 'rgb(34, 197, 94)' }} />
                      <span>{text}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}

            <TechButton
              onClick={handleToggleDemoMode}
              variant={demoMode ? 'secondary' : 'primary'}
              fullWidth
              icon={demoMode ? null : Zap}
            >
              {demoMode ? 'Exit Demo Mode' : 'Enable Now'}
            </TechButton>
          </TechCard>

          {/* Quick Start Guide */}
          <TechCard>
            <h3 style={{ fontSize: 18, fontWeight: 600, color: 'white', margin: '0 0 20px 0' }}>
              Quick Start
            </h3>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
              {[
                { icon: TrendingUp, title: 'Yield Vault', desc: 'Auto-optimize DeFi returns', path: '/vault', color: 'rgb(34, 211, 238)' },
                { icon: Gem, title: 'Treasury Bonds', desc: 'Stable RWA investments', path: '/treasury', color: 'rgb(59, 130, 246)' },
                { icon: Shield, title: 'AI Risk Control', desc: 'View risk assessment reports', path: '/dashboard', color: 'rgb(34, 197, 94)' }
              ].map((item, i) => {
                const Icon = item.icon;
                return (
                  <button
                    key={i}
                    onClick={() => navigate(item.path)}
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      gap: 12,
                      padding: 16,
                      background: 'rgba(15, 23, 42, 0.5)',
                      border: '1px solid rgba(34, 211, 238, 0.2)',
                      borderRadius: 8,
                      cursor: 'pointer',
                      textAlign: 'left',
                      transition: 'all 0.3s',
                      width: '100%'
                    }}
                    onMouseEnter={(e) => {
                      e.currentTarget.style.borderColor = item.color;
                      e.currentTarget.style.transform = 'translateX(4px)';
                      e.currentTarget.style.boxShadow = `0 4px 12px ${item.color}40`;
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.borderColor = 'rgba(34, 211, 238, 0.2)';
                      e.currentTarget.style.transform = 'translateX(0)';
                      e.currentTarget.style.boxShadow = 'none';
                    }}
                  >
                    <div style={{
                      width: 40,
                      height: 40,
                      borderRadius: 8,
                      background: item.color,
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      color: 'white',
                      flexShrink: 0,
                      boxShadow: `0 0 15px ${item.color}60`
                    }}>
                      <Icon size={20} />
                    </div>
                    <div style={{ flex: 1 }}>
                      <div style={{ fontSize: 14, fontWeight: 600, color: 'white', marginBottom: 2 }}>
                        {item.title}
                      </div>
                      <div style={{ fontSize: 12, color: 'rgba(203, 213, 225, 0.6)' }}>
                        {item.desc}
                      </div>
                    </div>
                  </button>
                );
              })}
            </div>
          </TechCard>
        </div>
      </div>
    </TechContainer>
  );
}
