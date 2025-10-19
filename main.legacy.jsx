import React, { useState, useEffect } from 'react';
import { createRoot } from 'react-dom/client';
import { Wallet, TrendingUp, Award, Users, DollarSign, Activity, Lock, Flame, ExternalLink, RefreshCw, LogOut, CheckCircle, AlertCircle, Info } from 'lucide-react';

const API_BASE = 'http://localhost:8080';
const CHAIN_ID = '0xaa36a7';
const TOKEN_ADDRESS = '0x9e667a4ce092086C63c667e1Ea575B9Aa2a4762B';
const STAKING_ADDRESS = '0x3dBF997f45b7AF6EA274A7c4e9BaaBAC79989eed';

const ERC20_ABI = [{"constant": false, "inputs": [{"name": "_spender", "type": "address"}, {"name": "_value", "type": "uint256"}], "name": "approve", "outputs": [{"name": "", "type": "bool"}], "type": "function"}, {"constant": true, "inputs": [{"name": "_owner", "type": "address"}], "name": "balanceOf", "outputs": [{"name": "balance", "type": "uint256"}], "type": "function"}];

const STAKING_ABI = [{"inputs": [{"name": "amount", "type": "uint256"}], "name": "stake", "outputs": [], "stateMutability": "nonpayable", "type": "function"}, {"inputs": [{"name": "amount", "type": "uint256"}], "name": "withdraw", "outputs": [], "stateMutability": "nonpayable", "type": "function"}, {"inputs": [], "name": "claimRewards", "outputs": [], "stateMutability": "nonpayable", "type": "function"}, {"inputs": [{"name": "user", "type": "address"}], "name": "getUserStake", "outputs": [{"name": "", "type": "uint256"}], "stateMutability": "view", "type": "function"}];

function App() {
  const [account, setAccount] = useState(null);
  const [balance, setBalance] = useState('0');
  const [points, setPoints] = useState('0');
  const [stakedAmount, setStakedAmount] = useState('0');
  const [badges, setBadges] = useState([]);
  const [leaderboard, setLeaderboard] = useState([]);
  const [loading, setLoading] = useState(false);
  const [txStatus, setTxStatus] = useState(null);
  const [activeTab, setActiveTab] = useState('home');

  useEffect(() => {
    if (account) {
      fetchData();
      const interval = setInterval(fetchData, 5000);
      return () => clearInterval(interval);
    }
  }, [account]);

  const fetchData = async () => {
    if (!account) return;
    try {
      const [balRes, ptsRes, badgesRes, lbRes] = await Promise.all([
        fetch(`${API_BASE}/users/${account}/balance`).catch(() => ({ ok: false })),
        fetch(`${API_BASE}/users/${account}/points`).catch(() => ({ ok: false })),
        fetch(`${API_BASE}/users/${account}/badges`).catch(() => ({ ok: false })),
        fetch(`${API_BASE}/leaderboard`).catch(() => ({ ok: false }))
      ]);
      if (balRes.ok) { const balData = await balRes.json(); setBalance(balData.balance || '0'); }
      if (ptsRes.ok) { const ptsData = await ptsRes.json(); setPoints(ptsData.points || '0'); }
      if (badgesRes.ok) { const badgesData = await badgesRes.json(); setBadges(badgesData.badges || []); }
      if (lbRes.ok) { const lbData = await lbRes.json(); setLeaderboard(lbData.items || []); }
      if (window.ethereum) {
        try {
          const provider = new window.ethers.providers.Web3Provider(window.ethereum);
          const contract = new window.ethers.Contract(STAKING_ADDRESS, STAKING_ABI, provider);
          const staked = await contract.getUserStake(account);
          setStakedAmount(window.ethers.utils.formatEther(staked));
        } catch (e) { console.error('Error fetching staked amount:', e); }
      }
    } catch (error) { console.error('Error fetching data:', error); }
  };

  const connectWallet = async () => {
    if (!window.ethereum) { alert('Please install MetaMask!'); return; }
    try {
      setLoading(true);
      const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
      const chainId = await window.ethereum.request({ method: 'eth_chainId' });
      if (chainId !== CHAIN_ID) {
        try {
          await window.ethereum.request({ method: 'wallet_switchEthereumChain', params: [{ chainId: CHAIN_ID }] });
        } catch (switchError) { alert('Please switch to Sepolia network'); setLoading(false); return; }
      }
      setAccount(accounts[0]);
      setTxStatus({ type: 'success', message: 'Wallet connected!' });
      setTimeout(() => setTxStatus(null), 3000);
    } catch (error) {
      console.error('Connection error:', error);
      setTxStatus({ type: 'error', message: 'Failed to connect wallet' });
      setTimeout(() => setTxStatus(null), 3000);
    } finally { setLoading(false); }
  };

  const disconnectWallet = () => { setAccount(null); setBalance('0'); setPoints('0'); setStakedAmount('0'); setBadges([]); };

  const handleStake = async (amount) => {
    if (!amount || parseFloat(amount) <= 0) { setTxStatus({ type: 'error', message: 'Enter valid amount' }); setTimeout(() => setTxStatus(null), 3000); return; }
    try {
      setLoading(true);
      const provider = new window.ethers.providers.Web3Provider(window.ethereum);
      const signer = provider.getSigner();
      const tokenContract = new window.ethers.Contract(TOKEN_ADDRESS, ERC20_ABI, signer);
      const stakingContract = new window.ethers.Contract(STAKING_ADDRESS, STAKING_ABI, signer);
      const amountWei = window.ethers.utils.parseEther(amount);
      setTxStatus({ type: 'info', message: 'Approving tokens...' });
      const approveTx = await tokenContract.approve(STAKING_ADDRESS, amountWei);
      await approveTx.wait();
      setTxStatus({ type: 'info', message: 'Staking tokens...' });
      const stakeTx = await stakingContract.stake(amountWei);
      await stakeTx.wait();
      setTxStatus({ type: 'success', message: 'Staked successfully!' });
      fetchData();
    } catch (error) {
      console.error('Stake error:', error);
      setTxStatus({ type: 'error', message: error.message || 'Staking failed' });
    } finally { setLoading(false); setTimeout(() => setTxStatus(null), 5000); }
  };

  const handleUnstake = async (amount) => {
    if (!amount || parseFloat(amount) <= 0) { setTxStatus({ type: 'error', message: 'Enter valid amount' }); setTimeout(() => setTxStatus(null), 3000); return; }
    try {
      setLoading(true);
      const provider = new window.ethers.providers.Web3Provider(window.ethereum);
      const signer = provider.getSigner();
      const stakingContract = new window.ethers.Contract(STAKING_ADDRESS, STAKING_ABI, signer);
      const amountWei = window.ethers.utils.parseEther(amount);
      setTxStatus({ type: 'info', message: 'Withdrawing tokens...' });
      const withdrawTx = await stakingContract.withdraw(amountWei);
      await withdrawTx.wait();
      setTxStatus({ type: 'success', message: 'Withdrawn successfully!' });
      fetchData();
    } catch (error) {
      console.error('Unstake error:', error);
      setTxStatus({ type: 'error', message: error.message || 'Withdrawal failed' });
    } finally { setLoading(false); setTimeout(() => setTxStatus(null), 5000); }
  };

  const handleClaimRewards = async () => {
    try {
      setLoading(true);
      const provider = new window.ethers.providers.Web3Provider(window.ethereum);
      const signer = provider.getSigner();
      const stakingContract = new window.ethers.Contract(STAKING_ADDRESS, STAKING_ABI, signer);
      setTxStatus({ type: 'info', message: 'Claiming rewards...' });
      const claimTx = await stakingContract.claimRewards();
      await claimTx.wait();
      setTxStatus({ type: 'success', message: 'Rewards claimed!' });
      fetchData();
    } catch (error) {
      console.error('Claim error:', error);
      setTxStatus({ type: 'error', message: error.message || 'Claim failed' });
    } finally { setLoading(false); setTimeout(() => setTxStatus(null), 5000); }
  };

  const formatAddress = (addr) => `${addr.slice(0, 6)}...${addr.slice(-4)}`;
  const formatNumber = (num) => parseFloat(num).toFixed(4);
  const userRank = leaderboard.findIndex(item => item.Address.toLowerCase() === account?.toLowerCase()) + 1;
  const tabs = [{ id: 'home', label: '首页', icon: Activity }, { id: 'defi', label: 'DeFi池', icon: DollarSign }, { id: 'leaderboard', label: '排行榜', icon: Users }, { id: 'badges', label: '徽章', icon: Award }, { id: 'airdrop', label: '空投', icon: TrendingUp }, { id: 'health', label: '健康', icon: Info }];

  return (<div style={{minHeight: '100vh', background: 'linear-gradient(to bottom right, #0f172a, #1e293b, #0f172a)', color: 'white'}}><header style={{borderBottom: '1px solid rgba(71, 85, 105, 0.5)', background: 'rgba(15, 23, 42, 0.8)', backdropFilter: 'blur(12px)'}}><div style={{maxWidth: '1280px', margin: '0 auto', padding: '16px 24px'}}><div style={{display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}><div style={{display: 'flex', alignItems: 'center', gap: '16px'}}><div style={{width: '40px', height: '40px', background: 'linear-gradient(to bottom right, #3b82f6, #9333ea)', borderRadius: '12px', display: 'flex', alignItems: 'center', justifyContent: 'center'}}><TrendingUp style={{width: '24px', height: '24px'}} /></div><div><h1 style={{fontSize: '24px', fontWeight: 'bold', margin: 0}}>币股交易</h1><p style={{fontSize: '12px', color: '#94a3b8', margin: 0}}>Loyalty DeFi Protocol</p></div></div>{account ? (<div style={{display: 'flex', alignItems: 'center', gap: '16px'}}><button onClick={fetchData} disabled={loading} style={{padding: '10px', background: 'transparent', border: 'none', borderRadius: '12px', cursor: 'pointer', transition: 'all 0.2s'}}><RefreshCw style={{width: '20px', height: '20px', color: 'white', animation: loading ? 'spin 1s linear infinite' : 'none'}} /></button><div style={{padding: '10px 16px', background: 'linear-gradient(to right, rgba(59, 130, 246, 0.1), rgba(147, 51, 234, 0.1))', border: '1px solid rgba(59, 130, 246, 0.2)', borderRadius: '12px'}}><div style={{fontSize: '10px', color: '#94a3b8', marginBottom: '2px'}}>已连接</div><div style={{fontFamily: 'monospace', fontSize: '14px', fontWeight: '600'}}>{formatAddress(account)}</div></div><button onClick={disconnectWallet} style={{padding: '10px', background: 'transparent', border: 'none', borderRadius: '12px', cursor: 'pointer', transition: 'all 0.2s'}}><LogOut style={{width: '20px', height: '20px', color: '#f87171'}} /></button></div>) : (<button onClick={connectWallet} disabled={loading} style={{padding: '12px 24px', background: 'linear-gradient(to right, #2563eb, #9333ea)', border: 'none', borderRadius: '12px', color: 'white', fontWeight: '600', cursor: 'pointer', transition: 'all 0.2s', display: 'flex', alignItems: 'center', gap: '8px', opacity: loading ? 0.5 : 1}}><Wallet style={{width: '20px', height: '20px'}} /><span>{loading ? 'Connecting...' : 'Connect Wallet'}</span></button>)}</div>{account && (<nav style={{display: 'flex', gap: '4px', marginTop: '24px', overflowX: 'auto'}}>{tabs.map(tab => (<button key={tab.id} onClick={() => setActiveTab(tab.id)} style={{padding: '10px 24px', borderRadius: '12px 12px 0 0', fontWeight: '500', transition: 'all 0.2s', display: 'flex', alignItems: 'center', gap: '8px', whiteSpace: 'nowrap', background: activeTab === tab.id ? '#1e293b' : 'transparent', color: activeTab === tab.id ? 'white' : '#94a3b8', border: 'none', borderTop: activeTab === tab.id ? '2px solid #3b82f6' : 'none', cursor: 'pointer'}}><tab.icon style={{width: '16px', height: '16px'}} /><span>{tab.label}</span></button>))}</nav>)}</div></header>{txStatus && (<div style={{position: 'fixed', top: '96px', right: '24px', zIndex: 50}}><div style={{padding: '16px 24px', borderRadius: '12px', backdropFilter: 'blur(12px)', border: '1px solid', boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1)', display: 'flex', alignItems: 'center', gap: '12px', background: txStatus.type === 'success' ? 'rgba(34, 197, 94, 0.2)' : txStatus.type === 'error' ? 'rgba(239, 68, 68, 0.2)' : 'rgba(59, 130, 246, 0.2)', borderColor: txStatus.type === 'success' ? 'rgba(34, 197, 94, 0.3)' : txStatus.type === 'error' ? 'rgba(239, 68, 68, 0.3)' : 'rgba(59, 130, 246, 0.3)'}}>{txStatus.type === 'success' && <CheckCircle style={{width: '20px', height: '20px', color: '#4ade80'}} />}{txStatus.type === 'error' && <AlertCircle style={{width: '20px', height: '20px', color: '#f87171'}} />}{txStatus.type === 'info' && <RefreshCw style={{width: '20px', height: '20px', color: '#60a5fa'}} />}<span style={{fontSize: '14px', fontWeight: '500'}}>{txStatus.message}</span></div></div>)}{!account ? (<div style={{display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: 'calc(100vh - 80px)'}}><div style={{textAlign: 'center', maxWidth: '448px', padding: '16px'}}><div style={{width: '96px', height: '96px', background: 'linear-gradient(to bottom right, #3b82f6, #9333ea)', borderRadius: '16px', margin: '0 auto 24px', display: 'flex', alignItems: 'center', justifyContent: 'center'}}><Wallet style={{width: '48px', height: '48px'}} /></div><h2 style={{fontSize: '30px', fontWeight: 'bold', marginBottom: '16px'}}>Welcome to Loyalty DeFi</h2><p style={{color: '#94a3b8', marginBottom: '24px'}}>Connect your wallet to access staking pools and earn rewards</p><button onClick={connectWallet} style={{padding: '16px 32px', background: 'linear-gradient(to right, #2563eb, #9333ea)', border: 'none', borderRadius: '12px', color: 'white', fontWeight: '600', cursor: 'pointer', transition: 'all 0.2s', display: 'inline-flex', alignItems: 'center', gap: '8px'}}><Wallet style={{width: '20px', height: '20px'}} /><span>Connect MetaMask</span></button></div></div>) : (<main style={{maxWidth: '1280px', margin: '0 auto', padding: '32px 24px'}}>{activeTab === 'home' && <div>HOME TAB - Add HomeView component here</div>}{activeTab === 'defi' && <div>DEFI TAB - Add DefiView component here</div>}{activeTab === 'leaderboard' && <div>LEADERBOARD TAB - Add LeaderboardView component here</div>}{activeTab === 'badges' && <div>BADGES TAB - Add BadgesView component here</div>}{activeTab === 'airdrop' && <div>AIRDROP TAB - Add AirdropView component here</div>}{activeTab === 'health' && <div>HEALTH TAB - Add HealthView component here</div>}</main>)}</div>);
}

createRoot(document.getElementById('root')).render(<App />);
