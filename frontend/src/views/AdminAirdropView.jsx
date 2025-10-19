import React, { useState, useEffect } from "react";
import { Settings, Upload, Play, XCircle, BarChart3, Plus } from "lucide-react";
import { useWallet } from "../web3/WalletContext";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

export default function AdminAirdropView() {
  const { address } = useWallet();
  const [campaigns, setCampaigns] = useState([]);
  const [loading, setLoading] = useState(false);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [selectedCampaign, setSelectedCampaign] = useState(null);
  const [stats, setStats] = useState(null);

  // Create campaign form state
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    start_time: "",
    end_time: "",
    total_budget: "",
    is_demo: true
  });

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

  useEffect(() => {
    fetchCampaigns();
  }, []);

  // Fetch campaign stats
  const fetchStats = async (campaignId) => {
    try {
      const res = await fetch(
        `${API_URL}/api/admin/airdrop/campaigns/${campaignId}/stats`,
        {
          headers: { Authorization: `Bearer ${address?.toLowerCase()}` }
        }
      );
      const data = await res.json();
      setStats(data);
    } catch (err) {
      console.error("Failed to fetch stats:", err);
    }
  };

  // Create campaign
  const handleCreate = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const res = await fetch(`${API_URL}/api/admin/airdrop/campaigns`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${address?.toLowerCase()}`
        },
        body: JSON.stringify(formData)
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.error || "Failed to create campaign");
      }

      const data = await res.json();
      alert(`Campaign created successfully! ID: ${data.id}`);
      setShowCreateForm(false);
      fetchCampaigns();

      // Reset form
      setFormData({
        name: "",
        description: "",
        start_time: "",
        end_time: "",
        total_budget: "",
        is_demo: true
      });
    } catch (err) {
      alert(`Error: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Upload CSV
  const handleUploadCSV = async (campaignId, file) => {
    setLoading(true);
    const formData = new FormData();
    formData.append("file", file);

    try {
      const res = await fetch(
        `${API_URL}/api/admin/airdrop/campaigns/${campaignId}/allocations/import`,
        {
          method: "POST",
          headers: { Authorization: `Bearer ${address?.toLowerCase()}` },
          body: formData
        }
      );

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.error || "Failed to upload CSV");
      }

      const data = await res.json();
      alert(`Success! Imported ${data.count} allocations`);
      fetchCampaigns();
    } catch (err) {
      alert(`Error: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Activate campaign
  const handleActivate = async (campaignId) => {
    if (!confirm("Are you sure you want to activate this campaign?")) return;

    setLoading(true);
    try {
      const res = await fetch(
        `${API_URL}/api/admin/airdrop/campaigns/${campaignId}/activate`,
        {
          method: "POST",
          headers: { Authorization: `Bearer ${address?.toLowerCase()}` }
        }
      );

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.error || "Failed to activate");
      }

      alert("Campaign activated successfully!");
      fetchCampaigns();
    } catch (err) {
      alert(`Error: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  // Close campaign
  const handleClose = async (campaignId) => {
    if (!confirm("Are you sure you want to close this campaign?")) return;

    setLoading(true);
    try {
      const res = await fetch(
        `${API_URL}/api/admin/airdrop/campaigns/${campaignId}/close`,
        {
          method: "POST",
          headers: { Authorization: `Bearer ${address?.toLowerCase()}` }
        }
      );

      if (!res.ok) throw new Error("Failed to close campaign");

      alert("Campaign closed successfully!");
      fetchCampaigns();
    } catch (err) {
      alert(`Error: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  if (!address) {
    return (
      <div className="container">
        <div className="card" style={{ padding: 48, textAlign: "center" }}>
          <Settings size={48} color="#F59E0B" style={{ margin: "0 auto 16px" }} />
          <h2>Admin Access Required</h2>
          <p className="muted">Please connect your admin wallet to manage airdrops</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container">
      <div className="row" style={{ justifyContent: "space-between", marginBottom: 24 }}>
        <div>
          <h1 style={{ margin: 0, fontSize: 32, display: "flex", alignItems: "center", gap: 12 }}>
            <Settings size={36} color="#6366F1" />
            Airdrop Admin
          </h1>
          <p className="muted" style={{ marginTop: 8 }}>Manage campaigns and allocations</p>
        </div>
        <button
          className="btn"
          onClick={() => setShowCreateForm(!showCreateForm)}
          style={{ background: "#10b981", display: "flex", alignItems: "center", gap: 8 }}
        >
          <Plus size={20} />
          Create Campaign
        </button>
      </div>

      {/* Create Form */}
      {showCreateForm && (
        <div className="card" style={{ padding: 24, marginBottom: 24 }}>
          <h3 style={{ marginBottom: 16 }}>Create New Campaign</h3>
          <form onSubmit={handleCreate} style={{ display: "grid", gap: 16 }}>
            <div>
              <label style={{ display: "block", marginBottom: 8, fontWeight: 600 }}>Campaign Name</label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                required
                style={{
                  width: "100%",
                  padding: 12,
                  borderRadius: 8,
                  border: "1px solid rgba(255,255,255,.1)",
                  background: "rgba(0,0,0,.3)",
                  color: "#fff"
                }}
              />
            </div>

            <div>
              <label style={{ display: "block", marginBottom: 8, fontWeight: 600 }}>Description</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                rows={3}
                style={{
                  width: "100%",
                  padding: 12,
                  borderRadius: 8,
                  border: "1px solid rgba(255,255,255,.1)",
                  background: "rgba(0,0,0,.3)",
                  color: "#fff"
                }}
              />
            </div>

            <div className="grid" style={{ gridTemplateColumns: "1fr 1fr", gap: 16 }}>
              <div>
                <label style={{ display: "block", marginBottom: 8, fontWeight: 600 }}>Start Time</label>
                <input
                  type="datetime-local"
                  value={formData.start_time}
                  onChange={(e) => setFormData({ ...formData, start_time: e.target.value + ":00Z" })}
                  required
                  style={{
                    width: "100%",
                    padding: 12,
                    borderRadius: 8,
                    border: "1px solid rgba(255,255,255,.1)",
                    background: "rgba(0,0,0,.3)",
                    color: "#fff"
                  }}
                />
              </div>

              <div>
                <label style={{ display: "block", marginBottom: 8, fontWeight: 600 }}>End Time</label>
                <input
                  type="datetime-local"
                  value={formData.end_time}
                  onChange={(e) => setFormData({ ...formData, end_time: e.target.value + ":00Z" })}
                  required
                  style={{
                    width: "100%",
                    padding: 12,
                    borderRadius: 8,
                    border: "1px solid rgba(255,255,255,.1)",
                    background: "rgba(0,0,0,.3)",
                    color: "#fff"
                  }}
                />
              </div>
            </div>

            <div>
              <label style={{ display: "block", marginBottom: 8, fontWeight: 600 }}>Total Budget (Points)</label>
              <input
                type="number"
                value={formData.total_budget}
                onChange={(e) => setFormData({ ...formData, total_budget: e.target.value })}
                required
                min="0"
                step="0.01"
                style={{
                  width: "100%",
                  padding: 12,
                  borderRadius: 8,
                  border: "1px solid rgba(255,255,255,.1)",
                  background: "rgba(0,0,0,.3)",
                  color: "#fff"
                }}
              />
            </div>

            <div className="row" style={{ gap: 12 }}>
              <button type="submit" className="btn" disabled={loading} style={{ background: "#10b981" }}>
                {loading ? "Creating..." : "Create Campaign"}
              </button>
              <button
                type="button"
                className="btn"
                onClick={() => setShowCreateForm(false)}
                style={{ background: "rgba(255,255,255,.1)" }}
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Campaigns List */}
      <div className="grid" style={{ gap: 16 }}>
        {campaigns.map((campaign) => (
          <div key={campaign.id} className="card" style={{ padding: 24 }}>
            <div className="row" style={{ justifyContent: "space-between", marginBottom: 16 }}>
              <div>
                <h3 style={{ margin: 0 }}>{campaign.name}</h3>
                <p className="muted" style={{ marginTop: 4 }}>{campaign.description}</p>
              </div>
              <div
                style={{
                  padding: "6px 12px",
                  borderRadius: 6,
                  background: campaign.status === "active" ? "#10b98122" : "rgba(255,255,255,.1)",
                  color: campaign.status === "active" ? "#10b981" : "#999",
                  fontSize: 12,
                  fontWeight: 600,
                  height: "fit-content"
                }}
              >
                {campaign.status.toUpperCase()}
              </div>
            </div>

            <div className="grid" style={{ gridTemplateColumns: "repeat(4, 1fr)", gap: 12, marginBottom: 16 }}>
              <div className="kpi" style={{ padding: 12 }}>
                <div className="title" style={{ fontSize: 11 }}>Budget</div>
                <div className="value" style={{ fontSize: 16 }}>{campaign.total_budget}</div>
              </div>
              <div className="kpi" style={{ padding: 12 }}>
                <div className="title" style={{ fontSize: 11 }}>Claimed</div>
                <div className="value" style={{ fontSize: 16 }}>{campaign.claimed_amount}</div>
              </div>
              <div className="kpi" style={{ padding: 12 }}>
                <div className="title" style={{ fontSize: 11 }}>Participants</div>
                <div className="value" style={{ fontSize: 16 }}>{campaign.participant_count}</div>
              </div>
              <div className="kpi" style={{ padding: 12 }}>
                <div className="title" style={{ fontSize: 11 }}>Created By</div>
                <div className="value" style={{ fontSize: 10 }}>
                  {campaign.created_by.slice(0, 6)}...{campaign.created_by.slice(-4)}
                </div>
              </div>
            </div>

            <div className="row" style={{ gap: 8, flexWrap: "wrap" }}>
              {campaign.status === "draft" && (
                <>
                  <label className="btn" style={{ background: "#6366F1", cursor: "pointer" }}>
                    <Upload size={16} style={{ marginRight: 8 }} />
                    Upload CSV
                    <input
                      type="file"
                      accept=".csv"
                      onChange={(e) => handleUploadCSV(campaign.id, e.target.files[0])}
                      style={{ display: "none" }}
                    />
                  </label>
                  <button
                    className="btn"
                    onClick={() => handleActivate(campaign.id)}
                    style={{ background: "#10b981" }}
                  >
                    <Play size={16} style={{ marginRight: 8 }} />
                    Activate
                  </button>
                </>
              )}

              {(campaign.status === "active" || campaign.status === "claimable") && (
                <button
                  className="btn"
                  onClick={() => handleClose(campaign.id)}
                  style={{ background: "#EF4444" }}
                >
                  <XCircle size={16} style={{ marginRight: 8 }} />
                  Close
                </button>
              )}

              <button
                className="btn"
                onClick={() => {
                  setSelectedCampaign(campaign.id);
                  fetchStats(campaign.id);
                }}
                style={{ background: "rgba(255,255,255,.1)" }}
              >
                <BarChart3 size={16} style={{ marginRight: 8 }} />
                View Stats
              </button>
            </div>

            {/* Stats Section */}
            {selectedCampaign === campaign.id && stats && (
              <div style={{ marginTop: 16, padding: 16, background: "rgba(99,102,241,.1)", borderRadius: 8 }}>
                <h4 style={{ marginBottom: 12 }}>Campaign Statistics</h4>
                <div className="grid" style={{ gridTemplateColumns: "repeat(3, 1fr)", gap: 12 }}>
                  <div>
                    <div className="muted" style={{ fontSize: 12 }}>Remaining</div>
                    <div style={{ fontWeight: 700 }}>{stats.remaining_amount}</div>
                  </div>
                  <div>
                    <div className="muted" style={{ fontSize: 12 }}>Total Allocations</div>
                    <div style={{ fontWeight: 700 }}>{stats.total_allocations}</div>
                  </div>
                  <div>
                    <div className="muted" style={{ fontSize: 12 }}>Claim Rate</div>
                    <div style={{ fontWeight: 700 }}>{stats.claim_rate}</div>
                  </div>
                </div>
              </div>
            )}
          </div>
        ))}
      </div>

      {campaigns.length === 0 && (
        <div className="card" style={{ padding: 48, textAlign: "center" }}>
          <p className="muted">No campaigns yet. Create your first one!</p>
        </div>
      )}
    </div>
  );
}
