import React, { useState, useEffect } from "react";
import { Gift, CheckCircle, Clock, AlertCircle, ChevronRight } from "lucide-react";
import { useWallet } from "../web3/WalletContext";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

export default function AirdropView() {
  const { address, signer } = useWallet();
  const [campaigns, setCampaigns] = useState([]);
  const [eligibility, setEligibility] = useState({});
  const [loading, setLoading] = useState(false);
  const [claimingId, setClaimingId] = useState(null);

  // Fetch campaigns
  const fetchCampaigns = async () => {
    try {
      const res = await fetch(`${API_URL}/api/airdrop/campaigns`);
      const data = await res.json();
      setCampaigns(data.campaigns || []);
    } catch (err) {
      console.error("Failed to fetch campaigns:", err);
    }
  };

  // Check eligibility for all campaigns
  const checkAllEligibility = async () => {
    if (!address || campaigns.length === 0) return;

    const eligibilityMap = {};
    for (const campaign of campaigns) {
      try {
        const res = await fetch(
          `${API_URL}/api/airdrop/campaigns/${campaign.id}/eligibility?address=${address.toLowerCase()}`
        );
        const data = await res.json();
        eligibilityMap[campaign.id] = data;
      } catch (err) {
        console.error(`Failed to check eligibility for campaign ${campaign.id}:`, err);
      }
    }
    setEligibility(eligibilityMap);
  };

  useEffect(() => {
    fetchCampaigns();
  }, []);

  useEffect(() => {
    if (address && campaigns.length > 0) {
      checkAllEligibility();
    }
  }, [address, campaigns]);

  // Claim airdrop
  const handleClaim = async (campaign) => {
    if (!signer || !address) {
      alert("Please connect your wallet first");
      return;
    }

    const eligible = eligibility[campaign.id];
    if (!eligible?.eligible) {
      alert("You are not eligible for this airdrop");
      return;
    }

    if (eligible.claimed) {
      alert("You have already claimed this airdrop");
      return;
    }

    setClaimingId(campaign.id);
    setLoading(true);

    try {
      // Generate nonce
      const nonce = Date.now().toString();

      // Create message to sign
      const message = `Claim airdrop from campaign ${campaign.id} with nonce ${nonce}`;

      // Sign message
      const signature = await signer.signMessage(message);

      // Submit claim
      const res = await fetch(`${API_URL}/api/airdrop/campaigns/${campaign.id}/claim`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          address: address.toLowerCase(),
          nonce,
          signature
        })
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.error || "Failed to claim airdrop");
      }

      const data = await res.json();
      alert(`Success! You claimed ${data.amount} points!`);

      // Refresh eligibility
      checkAllEligibility();
    } catch (err) {
      console.error("Claim error:", err);
      alert(`Error: ${err.message}`);
    } finally {
      setLoading(false);
      setClaimingId(null);
    }
  };

  // Map status to display
  const getStatusDisplay = (campaign) => {
    const now = new Date();
    const startTime = new Date(campaign.start_time);
    const endTime = new Date(campaign.end_time);

    if (campaign.status === "draft") return { label: "草稿", color: "#666", icon: <Clock size={20} color="#666" /> };
    if (campaign.status === "scheduled") return { label: "即将开始", color: "#6366F1", icon: <AlertCircle size={20} color="#6366F1" /> };
    if (campaign.status === "active") return { label: "进行中", color: "#10b981", icon: <Gift size={20} color="#10b981" /> };
    if (campaign.status === "claimable") return { label: "可领取", color: "#F59E0B", icon: <Clock size={20} color="#F59E0B" /> };
    if (campaign.status === "closed") return { label: "已结束", color: "#666", icon: <CheckCircle size={20} color="#666" /> };
    return { label: campaign.status, color: "#999", icon: null };
  };

  // Calculate stats
  const availableCampaigns = campaigns.filter(c => c.status === "active" || c.status === "claimable").length;
  const eligibleCount = Object.values(eligibility).filter(e => e.eligible && !e.claimed).length;
  const claimedCount = Object.values(eligibility).filter(e => e.claimed).length;

  return (
    <div className="container">
      <div style={{ marginBottom: 24 }}>
        <h1 style={{ margin: 0, fontSize: 32, display: 'flex', alignItems: 'center', gap: 12 }}>
          <Gift size={36} color="#10b981" />
          空投中心
        </h1>
        <p className="muted" style={{ marginTop: 8 }}>查看并领取你的空投奖励</p>
      </div>

      {/* 统计卡片 */}
      <div className="grid" style={{ gridTemplateColumns: 'repeat(4, 1fr)', gap: 16, marginBottom: 24 }}>
        <div className="kpi">
          <div className="title">进行中</div>
          <div className="value" style={{ color: "#10b981" }}>{availableCampaigns}</div>
        </div>
        <div className="kpi">
          <div className="title">可领取</div>
          <div className="value" style={{ color: "#F59E0B" }}>{address ? eligibleCount : "-"}</div>
        </div>
        <div className="kpi">
          <div className="title">已领取</div>
          <div className="value" style={{ color: "#6366F1" }}>{address ? claimedCount : "-"}</div>
        </div>
        <div className="kpi">
          <div className="title">总活动</div>
          <div className="value" style={{ color: "#A855F7" }}>{campaigns.length}</div>
        </div>
      </div>

      {/* 空投列表 */}
      {!address ? (
        <div className="card" style={{ padding: 48, textAlign: 'center', background: 'rgba(245, 158, 11, .1)', borderColor: '#F59E0B' }}>
          <Gift size={48} color="#F59E0B" style={{ margin: '0 auto 16px' }} />
          <div style={{ fontWeight: 700, marginBottom: 8 }}>请先连接钱包</div>
          <div className="muted">连接钱包后即可查看你的空投</div>
        </div>
      ) : campaigns.length === 0 ? (
        <div className="card" style={{ padding: 48, textAlign: 'center' }}>
          <p className="muted">暂无空投活动</p>
        </div>
      ) : (
        <div className="grid" style={{ gap: 16 }}>
          {campaigns.map((campaign) => {
            const statusDisplay = getStatusDisplay(campaign);
            const eligible = eligibility[campaign.id];
            const canClaim = eligible?.eligible && !eligible?.claimed &&
                             (campaign.status === "active" || campaign.status === "claimable");
            const isClaimed = eligible?.claimed;

            return (
              <div
                key={campaign.id}
                className="card"
                style={{
                  padding: 24,
                  display: 'flex',
                  alignItems: 'center',
                  gap: 20,
                  opacity: campaign.status === "closed" ? 0.6 : 1,
                  borderColor: canClaim ? "#10b981" : 'rgba(255,255,255,.1)',
                  position: 'relative',
                  overflow: 'hidden'
                }}
              >
                {/* 背景装饰 */}
                <div style={{
                  position: 'absolute',
                  top: -30,
                  right: -30,
                  fontSize: 120,
                  opacity: 0.05
                }}>
                  🎁
                </div>

                {/* 图标 */}
                <div style={{
                  fontSize: 60,
                  width: 80,
                  height: 80,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  background: `${statusDisplay.color}22`,
                  borderRadius: 16,
                  flexShrink: 0,
                  position: 'relative',
                  zIndex: 1
                }}>
                  🎁
                </div>

                {/* 信息 */}
                <div style={{ flex: 1, position: 'relative', zIndex: 1 }}>
                  <div className="row" style={{ justifyContent: 'space-between', marginBottom: 8 }}>
                    <div style={{ fontWeight: 700, fontSize: 18 }}>{campaign.name}</div>
                    <div className="row" style={{ gap: 8 }}>
                      {statusDisplay.icon}
                      <span style={{ fontSize: 13, color: statusDisplay.color }}>
                        {statusDisplay.label}
                      </span>
                    </div>
                  </div>

                  <div className="muted" style={{ fontSize: 14, marginBottom: 12 }}>
                    {campaign.description || "No description"}
                  </div>

                  <div className="row" style={{ gap: 16, flexWrap: 'wrap' }}>
                    {eligible?.amount && (
                      <div>
                        <span className="muted" style={{ fontSize: 12 }}>你的份额: </span>
                        <span style={{ fontWeight: 700, color: "#10b981" }}>{eligible.amount} points</span>
                      </div>
                    )}
                    <div>
                      <span className="muted" style={{ fontSize: 12 }}>总预算: </span>
                      <span>{campaign.total_budget}</span>
                    </div>
                    <div>
                      <span className="muted" style={{ fontSize: 12 }}>已领取: </span>
                      <span>{campaign.claimed_amount}</span>
                    </div>
                    <div>
                      <span className="muted" style={{ fontSize: 12 }}>参与人数: </span>
                      <span>{campaign.participant_count}</span>
                    </div>
                    {eligible && !eligible.eligible && !eligible.claimed && (
                      <div style={{ color: '#F59E0B', fontSize: 12 }}>
                        ⚠️ {eligible.reason || "不符合领取条件"}
                      </div>
                    )}
                  </div>
                </div>

                {/* 操作按钮 */}
                <div style={{ position: 'relative', zIndex: 1 }}>
                  {isClaimed ? (
                    <div style={{
                      padding: '12px 24px',
                      borderRadius: 12,
                      background: 'rgba(16, 185, 129, .2)',
                      color: '#10b981',
                      display: 'flex',
                      alignItems: 'center',
                      gap: 8,
                      fontWeight: 600
                    }}>
                      <CheckCircle size={20} />
                      已领取
                    </div>
                  ) : canClaim ? (
                    <button
                      className="btn"
                      style={{
                        background: "#10b981",
                        display: 'flex',
                        alignItems: 'center',
                        gap: 8,
                        padding: '12px 24px'
                      }}
                      onClick={() => handleClaim(campaign)}
                      disabled={loading && claimingId === campaign.id}
                    >
                      {loading && claimingId === campaign.id ? "领取中..." : "领取"}
                      <ChevronRight size={16} />
                    </button>
                  ) : (
                    <button
                      className="btn"
                      style={{
                        background: 'rgba(255,255,255,.1)',
                        cursor: 'not-allowed',
                        opacity: 0.5
                      }}
                      disabled
                    >
                      {statusDisplay.label}
                    </button>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* 说明 */}
      <div className="card" style={{ marginTop: 24, padding: 20, background: 'rgba(99, 102, 241, .1)', borderColor: '#6366F1' }}>
        <div style={{ fontWeight: 700, marginBottom: 8, display: 'flex', alignItems: 'center', gap: 8 }}>
          <AlertCircle size={20} color="#6366F1" />
          空投说明
        </div>
        <div className="muted" style={{ fontSize: 14, lineHeight: 1.6 }}>
          • 连接钱包即可查看你的空投资格<br/>
          • 每个空投活动都有独立的白名单<br/>
          • 领取时需要钱包签名验证（不上链，无gas费）<br/>
          • 领取的points将直接添加到你的账户<br/>
          • 注意活动截止日期，过期后将无法领取
        </div>
      </div>
    </div>
  );
}
