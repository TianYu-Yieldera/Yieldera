import React, { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
const API = (path) => `http://localhost:8080${path}`;

function App() {
  const [addr, setAddr] = useState("0x3C07226A3f1488320426eB5FE9976f72E5712346");
  const [bal, setBal] = useState("0");
  const [pts, setPts] = useState("0");
  const [board, setBoard] = useState([]);
  const [badges, setBadges] = useState([]);

  const refresh = async () => {
    try {
      const b = await fetch(API(`/users/${addr}/balance`));
      setBal(b.ok ? (await b.json()).balance : "0");
      const p = await fetch(API(`/users/${addr}/points`));
      setPts((await p.json()).points ?? "0");
      const l = await fetch(API(`/leaderboard`));
      setBoard((await l.json()).items ?? []);
      const bd = await fetch(API(`/users/${addr}/badges`));
      setBadges((await bd.json()).badges ?? []);
    } catch (e) {}
  };
  useEffect(() => { const t = setInterval(refresh, 3000); return () => clearInterval(t); }, [addr]);

  return (
    <div style={{fontFamily:"system-ui, -apple-system, Segoe UI, Roboto", padding:"24px", maxWidth:900, margin:"0 auto"}}>
      <h1>⭐ Loyalty Points — Real Chain</h1>
      <p>订阅 ERC-20 Transfer & Staking（6 确认），积分滚动与徽章发放。</p>
      <div style={{display:"flex", gap:12, alignItems:"center"}}>
        <input style={{flex:1, padding:"8px 12px"}} value={addr} onChange={e=>setAddr(e.target.value)} />
        <button onClick={refresh} style={{padding:"8px 12px"}}>Refresh</button>
      </div>
      <div style={{display:"grid", gridTemplateColumns:"1fr 1fr", gap:16, marginTop:16}}>
        <Card title="Balance" value={bal} />
        <Card title="Points" value={pts} />
      </div>
      <h3 style={{marginTop:24}}>Badges</h3>
      <div style={{display:"flex", gap:8, flexWrap:"wrap"}}>
        {badges.length===0? <span>—</span> : badges.map((b,i)=><Badge key={i} code={b}/>)}
      </div>
      <h3 style={{marginTop:24}}>Leaderboard (Top 20)</h3>
      <table style={{width:"100%", borderCollapse:"collapse"}}>
        <thead><tr><th style={{textAlign:"left"}}>Address</th><th style={{textAlign:"right"}}>Points</th></tr></thead>
        <tbody>{board.map((it, idx)=>(<tr key={idx} style={{borderTop:"1px solid #eee"}}><td>{it.Address}</td><td style={{textAlign:"right"}}>{it.Points}</td></tr>))}</tbody>
      </table>
    </div>
  );
}
function Card({title, value}){ return (<div style={{border:"1px solid #eee", borderRadius:12, padding:16}}><h3>{title}</h3><div style={{fontSize:28}}>{value}</div></div>); }
function Badge({code}){ return <span style={{padding:"6px 10px", border:"1px solid #ddd", borderRadius:999, fontWeight:600}}>{code}</span>; }
createRoot(document.getElementById("root")).render(<App />);
