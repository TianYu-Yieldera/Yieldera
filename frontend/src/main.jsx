import React from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { WalletProvider } from "./web3/WalletContext";
import { DemoModeProvider } from "./web3/DemoModeContext";
import Header from "./components/Header";
import Landing from "./views/Landing";
import HomeView from "./views/HomeView";
import DeFiPoolView from "./views/DeFiPoolView";
import LeaderboardView from "./views/LeaderboardView";
import BadgesView from "./views/BadgesView";
import AirdropView from "./views/AirdropView";
import AdminAirdropView from "./views/AdminAirdropView";
import OnchainAirdropView from "./views/OnchainAirdropView";
import VaultView from "./views/VaultView";
import RWAMarketView from "./views/RWAMarketView";
import StatusView from "./views/StatusView";
import TutorialView from "./views/TutorialView";

function App(){
  return (
    <BrowserRouter>
      <DemoModeProvider>
        <WalletProvider>
          <Header />
          <Routes>
            <Route path='/' element={<Landing />} />
            <Route path='/dashboard' element={<HomeView />} />
            <Route path='/staking' element={<DeFiPoolView />} />
            <Route path='/leaderboard' element={<LeaderboardView />} />
            <Route path='/badges' element={<BadgesView />} />
            <Route path='/airdrop' element={<AirdropView />} />
            <Route path='/airdrop/onchain' element={<OnchainAirdropView />} />
            <Route path='/admin/airdrop' element={<AdminAirdropView />} />
            <Route path='/vault' element={<VaultView />} />
            <Route path='/rwa-market' element={<RWAMarketView />} />
            <Route path='/status' element={<StatusView />} />
            <Route path='/tutorial' element={<TutorialView />} />
          </Routes>
        </WalletProvider>
      </DemoModeProvider>
    </BrowserRouter>
  );
}

createRoot(document.getElementById("root")).render(<App />);
