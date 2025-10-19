import React, { useEffect, useState } from "react";
import { useWallet } from "../web3/WalletContext";
import { useDemoMode } from "../web3/DemoModeContext";
import { config } from "../config/env";

const API_BASE = config.api.baseUrl;

function KPICard({ title, value, sub }){
  return (
    <div className="kpi">
      <div className="title">{title}</div>
      <div className="value">{value}</div>
      {sub ? <div className="muted" style={{marginTop:4}}>{sub}</div> : null}
    </div>
  );
}

export default function HomeView(){
  const { address } = useWallet();
  const { demoMode, demoAddress, demoData } = useDemoMode();
  const [data, setData] = useState({
    metrics:{ netWorth: 0, pnl24h: 0 },
    points:{ total: 0, available: 0 },
    staking:{ tvl: 0, pendingRewards: 0 },
    txs:[], extra:{ badgesCount: 0, rank: null }
  });

  useEffect(() => {
    if (demoMode) {
      // 演示模式：使用模拟数据
      const totalDeposits = demoData.defiDeposits.reduce((sum, d) => sum + d.amount, 0);
      const netWorth = demoData.pfiTokens + demoData.points + totalDeposits + demoData.stakedTokens + demoData.collateral + demoData.stablecoinMinted;

      setData({
        metrics:{ netWorth, pnl24h: 125 },
        points:{ total: demoData.points, available: demoData.points },
        staking:{ tvl: demoData.stakedTokens, pendingRewards: demoData.stakingRewards },
        txs:[],
        extra:{ badgesCount: 3, rank: 42 }
      });
    } else {
      // 真实模式：从API获取数据
      const addr = address || "0x3C07226A3f1488320426eB5FE9976f72E5712346";
      const j = (p) => fetch(`${API_BASE}${p}`).then(r=>r.json()).catch(()=>null);
      Promise.all([ j(`/users/${addr}/balance`), j(`/users/${addr}/points`), j(`/leaderboard`), j(`/users/${addr}/badges`) ])
        .then(([bal, pts, board, badges]) => {
          const rank = Array.isArray(board) ? (board.findIndex(x => String(x.Address||x.address).toLowerCase()===addr.toLowerCase())+1 || null) : null;
          setData({
            metrics:{ netWorth: Number(bal?.balance||0), pnl24h: 0 },
            points:{ total: Number(pts?.points||0), available: Number(pts?.points||0) },
            staking:{ tvl: 0, pendingRewards: 0 },
            txs:[],
            extra:{ badgesCount: Array.isArray(badges)? badges.length: 0, rank }
          });
        });
    }
  }, [address, demoMode, demoData]);

  const displayAddress = demoMode ? demoAddress : address;

  return (
    <div className="container">
      <div className="row" style={{justifyContent:'space-between', marginBottom:12}}>
        <h2>Dashboard {demoMode && <span style={{fontSize: '14px', background: '#10b981', color: 'white', padding: '4px 12px', borderRadius: '6px', marginLeft: '8px'}}>演示模式</span>}</h2>
        <div className="muted">{displayAddress ? `Connected: ${displayAddress.slice(0,6)}...${displayAddress.slice(-4)}` : 'Not connected'}</div>
      </div>
      <div className="grid grid-2" style={{gridTemplateColumns:'repeat(5,1fr)'}}>
        <KPICard title="Net Worth" value={`$${data.metrics.netWorth.toLocaleString()}`} sub="Total portfolio value" />
        {demoMode && <KPICard title="PFI Tokens" value={demoData.pfiTokens.toLocaleString()} sub="可交易代币" />}
        <KPICard title="Points" value={data.points.total.toLocaleString()} sub={`Available: ${data.points.available.toLocaleString()}`} />
        <KPICard title="Staked" value={`${data.staking.tvl.toLocaleString()} PFI`} sub={`Rewards: ${(data.staking.pendingRewards || 0).toFixed(2)} Points`} />
        <KPICard title="Badges" value={`${data.extra.badgesCount}`} sub="Owned badges" />
      </div>
      <div className="card" style={{marginTop:16}}>
        <div style={{fontWeight:600, marginBottom:8}}>Recent Activity</div>
        <div className="muted">接入链上交易后展示</div>
      </div>
    </div>
  );
}
