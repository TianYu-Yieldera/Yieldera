import React, { Suspense, lazy } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { WalletProvider } from "./web3/WalletContext";
import { DemoModeProvider } from "./web3/DemoModeContext";
import Header from "./components/Header";

// Lazy load all views for better performance
const Landing = lazy(() => import("./views/Landing"));
const HomeView = lazy(() => import("./views/HomeView"));
const VaultView = lazy(() => import("./views/VaultView"));
const TreasuryMarketView = lazy(() => import("./views/TreasuryMarketView"));
const TreasuryDetailView = lazy(() => import("./views/TreasuryDetailView"));
const TreasuryHoldingsView = lazy(() => import("./views/TreasuryHoldingsView"));
const TutorialView = lazy(() => import("./views/TutorialView"));
const MonitoringView = lazy(() => import("./views/MonitoringView"));

// Loading component with better UX
function LoadingFallback() {
  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      background: 'linear-gradient(135deg, rgb(15, 23, 42) 0%, rgb(30, 41, 59) 100%)',
      gap: 24
    }}>
      <style>{`
        @keyframes spin {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
        @keyframes pulse-ring {
          0%, 100% { transform: scale(1); opacity: 0.5; }
          50% { transform: scale(1.1); opacity: 0.8; }
        }
      `}</style>

      {/* Animated loading spinner */}
      <div style={{ position: 'relative' }}>
        <div style={{
          width: 80,
          height: 80,
          border: '3px solid rgba(34, 211, 238, 0.2)',
          borderTop: '3px solid rgb(34, 211, 238)',
          borderRadius: '50%',
          animation: 'spin 1s linear infinite'
        }} />
        <div style={{
          position: 'absolute',
          top: -10,
          left: -10,
          width: 100,
          height: 100,
          border: '2px solid rgba(34, 211, 238, 0.1)',
          borderRadius: '50%',
          animation: 'pulse-ring 2s ease-in-out infinite'
        }} />
      </div>

      <div style={{
        fontSize: 18,
        fontWeight: 600,
        color: 'white',
        letterSpacing: 0.5
      }}>
        Loading...
      </div>

      <div style={{
        fontSize: 14,
        color: 'rgba(203, 213, 225, 0.7)'
      }}>
        Preparing your dashboard
      </div>
    </div>
  );
}

function App(){
  return (
    <BrowserRouter>
      <DemoModeProvider>
        <WalletProvider>
          <Header />
          <Suspense fallback={<LoadingFallback />}>
            <Routes>
              <Route path='/' element={<Landing />} />
              <Route path='/dashboard' element={<HomeView />} />
              <Route path='/vault' element={<VaultView />} />
              <Route path='/treasury' element={<TreasuryMarketView />} />
              <Route path='/treasury/:assetId' element={<TreasuryDetailView />} />
              <Route path='/treasury/holdings' element={<TreasuryHoldingsView />} />
              <Route path='/monitoring' element={<MonitoringView />} />
              <Route path='/tutorial' element={<TutorialView />} />
              {/* 删除的路由重定向到首页 */}
              <Route path='*' element={<Navigate to="/dashboard" replace />} />
            </Routes>
          </Suspense>
        </WalletProvider>
      </DemoModeProvider>
    </BrowserRouter>
  );
}

createRoot(document.getElementById("root")).render(<App />);
