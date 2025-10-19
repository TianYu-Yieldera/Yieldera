import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ChevronRight, ChevronLeft, Check, Home, Play } from 'lucide-react';
import { useDemoMode } from '../web3/DemoModeContext';

// 视频配置 - 修改这里来更换视频
const VIDEO_CONFIG = {
  type: 'placeholder', // 'youtube' | 'local' | 'placeholder'
  youtubeId: '', // YouTube 视频 ID (例如: 'dQw4w9WgXcQ')
  localPath: '/videos/tutorial.mp4', // 本地视频路径
  posterImage: '/pointfi-logo-mark.svg' // 视频封面图
};

// 视频播放器组件
function VideoPlayer({ config }) {
  const [showVideo, setShowVideo] = useState(true);

  if (!showVideo) return null;

  // YouTube 视频
  if (config.type === 'youtube' && config.youtubeId) {
    return (
      <div style={{ marginBottom: '24px', position: 'relative' }}>
        <div style={{ position: 'relative', paddingBottom: '56.25%', height: 0, borderRadius: '12px', overflow: 'hidden', boxShadow: '0 10px 25px rgba(0,0,0,0.3)' }}>
          <iframe
            style={{ position: 'absolute', top: 0, left: 0, width: '100%', height: '100%' }}
            src={`https://www.youtube.com/embed/${config.youtubeId}?rel=0&modestbranding=1`}
            title="Yieldera Tutorial"
            frameBorder="0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            allowFullScreen
          />
        </div>
        <button
          onClick={() => setShowVideo(false)}
          style={{
            marginTop: '12px',
            padding: '8px 16px',
            background: 'rgba(255,255,255,0.1)',
            border: '1px solid rgba(255,255,255,0.2)',
            borderRadius: '6px',
            color: '#9ca3af',
            fontSize: '13px',
            cursor: 'pointer',
            width: '100%'
          }}
        >
          隐藏视频
        </button>
      </div>
    );
  }

  // 本地视频
  if (config.type === 'local' && config.localPath) {
    return (
      <div style={{ marginBottom: '24px' }}>
        <video
          controls
          poster={config.posterImage}
          style={{ width: '100%', borderRadius: '12px', boxShadow: '0 10px 25px rgba(0,0,0,0.3)' }}
        >
          <source src={config.localPath} type="video/mp4" />
          您的浏览器不支持视频播放
        </video>
        <button
          onClick={() => setShowVideo(false)}
          style={{
            marginTop: '12px',
            padding: '8px 16px',
            background: 'rgba(255,255,255,0.1)',
            border: '1px solid rgba(255,255,255,0.2)',
            borderRadius: '6px',
            color: '#9ca3af',
            fontSize: '13px',
            cursor: 'pointer',
            width: '100%'
          }}
        >
          隐藏视频
        </button>
      </div>
    );
  }

  // 占位符 - 暂无视频
  return (
    <div style={{
      marginBottom: '24px',
      background: 'linear-gradient(135deg, #1e293b 0%, #0f172a 100%)',
      borderRadius: '12px',
      padding: '48px 32px',
      textAlign: 'center',
      border: '2px dashed rgba(255,255,255,0.1)',
      position: 'relative',
      overflow: 'hidden'
    }}>
      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        background: 'radial-gradient(circle at 50% 50%, rgba(102, 126, 234, 0.1) 0%, transparent 70%)',
        pointerEvents: 'none'
      }} />
      <div style={{ position: 'relative', zIndex: 1 }}>
        <div style={{
          width: '80px',
          height: '80px',
          margin: '0 auto 20px',
          background: 'rgba(102, 126, 234, 0.2)',
          borderRadius: '50%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          border: '3px solid rgba(102, 126, 234, 0.3)'
        }}>
          <Play size={36} style={{ color: '#667eea', marginLeft: '4px' }} />
        </div>
        <h3 style={{ color: '#fff', fontSize: '20px', marginBottom: '12px', fontWeight: '600' }}>
          📹 系统演示视频
        </h3>
        <p style={{ color: '#9ca3af', fontSize: '14px', lineHeight: '1.6', marginBottom: '20px', maxWidth: '500px', margin: '0 auto 20px' }}>
          我们正在为您准备专业的系统演示视频！<br/>
          视频将展示完整的使用流程，包括理财金库、RWA 资产购买等核心功能。
        </p>
        <div style={{
          display: 'inline-block',
          padding: '10px 20px',
          background: 'rgba(102, 126, 234, 0.15)',
          borderRadius: '8px',
          border: '1px solid rgba(102, 126, 234, 0.3)',
          color: '#a5b4fc',
          fontSize: '13px'
        }}>
          <strong>预计时长:</strong> 6-8 分钟 | <strong>语言:</strong> 英文配音 + 中文字幕
        </div>
      </div>
    </div>
  );
}

export default function TutorialView() {
  const navigate = useNavigate();
  const [currentStep, setCurrentStep] = useState(0);
  const [completedSteps, setCompletedSteps] = useState([]);
  const { demoMode, demoAddress, enableDemoMode, disableDemoMode } = useDemoMode();

  const handleToggleDemoMode = async () => {
    try {
      if (!demoMode) {
        await enableDemoMode();
      } else {
        await disableDemoMode();
      }
    } catch (err) {
      console.error('Failed to toggle demo mode', err);
      alert('演示模式初始化失败，请稍后再试');
    }
  };

  const handleNext = () => {
    if (currentStep < tutorialSteps.length - 1) {
      setCompletedSteps([...completedSteps, currentStep]);
      setCurrentStep(currentStep + 1);
    } else {
      handleComplete();
    }
  };

  const handlePrevious = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const handleComplete = () => {
    localStorage.setItem('pointfi_tutorial_completed', 'true');
    navigate('/dashboard');
  };

  const handleClose = () => {
    navigate(-1); // 返回上一页
  };

  const tutorialSteps = [
    {
      id: 'welcome',
      title: '🎉 欢迎来到 Yieldera！',
      description: '观看演示视频或阅读下方指南了解系统',
      content: (
        <div style={{display: 'flex', flexDirection: 'column', gap: '16px'}}>
          {/* 视频播放器 */}
          <VideoPlayer config={VIDEO_CONFIG} />

          <div style={{
            background: 'linear-gradient(135deg, #6366F1, #A855F7)',
            padding: '24px',
            borderRadius: '12px',
            color: 'white'
          }}>
            <h3 style={{fontSize: '20px', fontWeight: '700', marginBottom: '16px'}}>🚀 你可以做什么：</h3>
            <ul style={{display: 'flex', flexDirection: 'column', gap: '8px', listStyle: 'none', padding: 0, margin: 0}}>
              {[
                '💰 存入积分赚取被动收益',
                '🪙 铸造 LUSD 稳定币（1:1锚定美元）',
                '🏆 访问 Uniswap、Aave 等真实 DeFi 协议'
              ].map((text, i) => (
                <li key={i} style={{display: 'flex', alignItems: 'center', gap: '8px'}}>
                  <Check size={20} />
                  <span>{text}</span>
                </li>
              ))}
            </ul>
          </div>

          <div style={{
            background: '#FEF3C7',
            border: '2px solid #FCD34D',
            padding: '20px',
            borderRadius: '12px'
          }}>
            <p style={{color: '#92400E', marginBottom: '12px', fontSize: '15px'}}>
              <strong>💡 推荐：启用演示模式</strong>
              <br/>
              <span style={{fontSize: '14px'}}>使用虚拟代币测试系统，无需真实交易！</span>
            </p>
            <button
              onClick={handleToggleDemoMode}
              style={{
                padding: '12px 20px',
                borderRadius: '10px',
                fontWeight: '600',
                border: 'none',
                cursor: 'pointer',
                background: demoMode ? '#10B981' : '#3B82F6',
                color: 'white',
                fontSize: '15px',
                width: '100%',
                transition: 'all 0.2s',
                boxShadow: demoMode ? '0 4px 12px rgba(16, 185, 129, 0.3)' : '0 4px 12px rgba(59, 130, 246, 0.3)'
              }}
            >
              {demoMode ? '✅ 演示模式已激活 - 10,000 测试积分' : '🎮 点击启用演示模式'}
            </button>
            {demoAddress && (
              <p style={{fontSize: '13px', color: '#6B7280', marginTop: '12px', textAlign: 'center'}}>
                演示钱包: <code style={{background: '#F3F4F6', padding: '2px 8px', borderRadius: '4px'}}>{demoAddress.slice(0, 6)}...{demoAddress.slice(-4)}</code>
              </p>
            )}
          </div>
        </div>
      ),
      action: '下一步'
    },
    {
      id: 'explore',
      title: '📊 探索功能',
      description: '看看主页上的7大功能模块',
      content: (
        <div style={{display: 'flex', flexDirection: 'column', gap: '20px'}}>
          <div style={{background: '#EEF2FF', padding: '20px', borderRadius: '12px', border: '1px solid #C7D2FE'}}>
            <h4 style={{margin: '0 0 12px 0', fontSize: '18px', color: '#4338CA'}}>🎯 关键功能区域</h4>
            <div style={{display: 'grid', gap: '12px'}}>
              {[
                { icon: '📈', name: '概览 Dashboard', desc: '查看你的资产总览、收益和积分' },
                { icon: '🏊', name: 'DeFi 池', desc: '存入资金到 Uniswap、Aave 等协议赚取高收益' },
                { icon: '🏆', name: '排行榜', desc: '查看积分排名，与其他用户竞争' }
              ].map((item, i) => (
                <div key={i} style={{
                  background: 'white',
                  padding: '16px',
                  borderRadius: '10px',
                  border: '1px solid #E0E7FF',
                  display: 'flex',
                  gap: '12px',
                  alignItems: 'start'
                }}>
                  <span style={{fontSize: '24px'}}>{item.icon}</span>
                  <div>
                    <div style={{fontWeight: '700', marginBottom: '4px'}}>{item.name}</div>
                    <div style={{fontSize: '13px', color: '#6B7280'}}>{item.desc}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div style={{background: '#D1FAE5', padding: '16px', borderRadius: '10px', border: '1px solid #6EE7B7'}}>
            <p style={{margin: 0, color: '#065F46', fontSize: '14px'}}>
              <strong>💡 提示：</strong> 关闭这个教程后，点击顶部菜单栏的任意功能开始探索！
            </p>
          </div>
        </div>
      ),
      action: '下一步'
    },
    {
      id: 'start',
      title: '🚀 开始使用',
      description: '你已经准备好了！',
      content: (
        <div style={{display: 'flex', flexDirection: 'column', gap: '20px', alignItems: 'center'}}>
          <div style={{
            background: 'linear-gradient(135deg, #10B981, #059669)',
            padding: '32px',
            borderRadius: '16px',
            color: 'white',
            textAlign: 'center',
            width: '100%'
          }}>
            <div style={{fontSize: '60px', marginBottom: '16px'}}>🎊</div>
            <h3 style={{fontSize: '24px', fontWeight: '700', margin: '0 0 12px 0'}}>准备就绪！</h3>
            <p style={{opacity: 0.95, margin: 0, fontSize: '15px'}}>
              {demoMode ? '你已启用演示模式，可以安全地探索所有功能' : '连接钱包后即可开始使用真实资产'}
            </p>
          </div>

          <div style={{background: '#F9FAFB', padding: '24px', borderRadius: '12px', border: '1px solid #E5E7EB', width: '100%'}}>
            <h4 style={{margin: '0 0 16px 0', fontSize: '16px', fontWeight: '700'}}>📝 快速开始指南：</h4>
            <div style={{display: 'flex', flexDirection: 'column', gap: '12px'}}>
              {[
                { num: '1', text: '点击顶部菜单栏的"概览"查看你的资产', highlight: false },
                { num: '2', text: '访问"DeFi 池"开始赚取收益', highlight: true },
                { num: '3', text: '在"排行榜"查看你的排名' }
              ].map((item, i) => (
                <div key={i} style={{
                  display: 'flex',
                  gap: '12px',
                  alignItems: 'center',
                  padding: '12px',
                  background: item.highlight ? '#EFF6FF' : 'white',
                  borderRadius: '8px',
                  border: item.highlight ? '2px solid #3B82F6' : '1px solid #E5E7EB'
                }}>
                  <div style={{
                    background: item.highlight ? '#3B82F6' : '#6B7280',
                    color: 'white',
                    width: '28px',
                    height: '28px',
                    borderRadius: '50%',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontWeight: '700',
                    fontSize: '14px',
                    flexShrink: 0
                  }}>{item.num}</div>
                  <span style={{fontSize: '14px', color: '#374151'}}>{item.text}</span>
                </div>
              ))}
            </div>
          </div>

          <div style={{background: '#DBEAFE', padding: '16px', borderRadius: '10px', border: '1px solid #93C5FD', width: '100%'}}>
            <p style={{margin: 0, color: '#1E40AF', fontSize: '14px', textAlign: 'center'}}>
              <strong>🎉 专业提示：</strong> 随时点击右上角的 <strong>❓ 帮助按钮</strong> 重新查看教程
            </p>
          </div>
        </div>
      ),
      action: '开始探索'
    }
  ];

  // Demo mode state is now managed by DemoModeContext

  const currentStepData = tutorialSteps[currentStep];
  const progress = ((currentStep + 1) / tutorialSteps.length) * 100;

  return (
    <div style={{
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      padding: '40px 20px',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center'
    }}>
      <div style={{
        background: 'white',
        borderRadius: '20px',
        boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.35)',
        maxWidth: '680px',
        width: '100%',
        maxHeight: '90vh',
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden'
      }}>
        {/* Header */}
        <div style={{
          background: 'linear-gradient(135deg, #6366F1 0%, #8B5CF6 50%, #A855F7 100%)',
          padding: '28px 32px',
          color: 'white'
        }}>
          <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'start', marginBottom: '12px'}}>
            <h2 style={{fontSize: '28px', fontWeight: '700', margin: 0, lineHeight: 1.2}}>{currentStepData.title}</h2>
            <button
              onClick={handleClose}
              className="btn"
              style={{
                background: 'rgba(255, 255, 255, 0.2)',
                border: 'none',
                color: 'white',
                cursor: 'pointer',
                padding: '8px 16px',
                display: 'flex',
                alignItems: 'center',
                gap: '6px',
                borderRadius: '8px',
                transition: 'all 0.2s',
                fontSize: '14px'
              }}
            >
              <Home size={16} />
              返回
            </button>
          </div>
          <p style={{opacity: 0.95, fontSize: '15px', margin: 0}}>{currentStepData.description}</p>

          {/* Progress Bar */}
          <div style={{marginTop: '20px'}}>
            <div style={{display: 'flex', justifyContent: 'space-between', fontSize: '13px', marginBottom: '8px', opacity: 0.9}}>
              <span>进度</span>
              <span>{currentStep + 1} / {tutorialSteps.length}</span>
            </div>
            <div style={{
              width: '100%',
              background: 'rgba(255, 255, 255, 0.25)',
              borderRadius: '10px',
              height: '6px',
              overflow: 'hidden'
            }}>
              <div style={{
                background: 'white',
                borderRadius: '10px',
                height: '100%',
                width: `${progress}%`,
                transition: 'width 0.4s cubic-bezier(0.4, 0, 0.2, 1)',
                boxShadow: '0 0 8px rgba(255, 255, 255, 0.5)'
              }} />
            </div>
          </div>
        </div>

        {/* Content */}
        <div style={{
          flex: 1,
          overflowY: 'auto',
          padding: '32px'
        }}>
          {currentStepData.content}
        </div>

        {/* Footer */}
        <div style={{
          borderTop: '1px solid #E5E7EB',
          padding: '24px 32px',
          background: '#F9FAFB',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          gap: '16px'
        }}>
          <button
            onClick={handlePrevious}
            disabled={currentStep === 0}
            className="btn"
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '6px',
              padding: '12px 24px',
              borderRadius: '10px',
              fontWeight: '600',
              fontSize: '15px',
              border: 'none',
              cursor: currentStep === 0 ? 'not-allowed' : 'pointer',
              background: currentStep === 0 ? '#E5E7EB' : 'white',
              color: currentStep === 0 ? '#9CA3AF' : '#374151',
              transition: 'all 0.2s',
              boxShadow: currentStep === 0 ? 'none' : '0 1px 3px rgba(0, 0, 0, 0.1)'
            }}
          >
            <ChevronLeft size={18} />
            上一步
          </button>

          <div style={{display: 'flex', gap: '8px'}}>
            {tutorialSteps.map((_, index) => (
              <div
                key={index}
                style={{
                  width: index === currentStep ? '32px' : '8px',
                  height: '8px',
                  borderRadius: '10px',
                  background: index === currentStep ? '#6366F1' : completedSteps.includes(index) ? '#10B981' : '#D1D5DB',
                  transition: 'all 0.3s ease',
                  boxShadow: index === currentStep ? '0 0 8px rgba(99, 102, 241, 0.5)' : 'none'
                }}
              />
            ))}
          </div>

          <button
            onClick={handleNext}
            className="btn"
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              padding: '12px 28px',
              background: currentStep === tutorialSteps.length - 1
                ? 'linear-gradient(135deg, #10B981, #059669)'
                : 'linear-gradient(135deg, #3B82F6, #6366F1)',
              color: 'white',
              borderRadius: '10px',
              fontWeight: '600',
              fontSize: '15px',
              border: 'none',
              cursor: 'pointer',
              transition: 'all 0.2s',
              boxShadow: '0 4px 12px rgba(59, 130, 246, 0.4)'
            }}
          >
            {currentStep === tutorialSteps.length - 1 ? (
              <>
                <Check size={20} />
                {currentStepData.action}
              </>
            ) : (
              <>
                {currentStepData.action}
                <ChevronRight size={20} />
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
