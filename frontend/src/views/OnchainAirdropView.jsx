import React, { useState, useEffect } from "react";
import { useWallet } from "../web3/WalletContext";
import { ethers } from "ethers";

// 合约 ABI（仅需要的函数）
const AIRDROP_ABI = [
  "function claim(uint256 campaignId, uint256 amount, bytes32[] calldata merkleProof) external",
  "function hasClaimed(uint256 campaignId, address user) external view returns (bool)",
  "function campaigns(uint256) external view returns (string name, string description, bytes32 merkleRoot, uint256 totalBudget, uint256 claimedAmount, uint256 startTime, uint256 endTime, bool isActive, uint256 participantCount, address creator)",
  "function campaignIdCounter() external view returns (uint256)",
];

// 合约地址（部署后替换）
const AIRDROP_CONTRACT_ADDRESS = "0x0000000000000000000000000000000000000000";

export default function OnchainAirdropView() {
  const { address, signer } = useWallet();
  const [campaigns, setCampaigns] = useState([]);
  const [loading, setLoading] = useState(true);
  const [claiming, setClaiming] = useState(false);
  const [selectedCampaign, setSelectedCampaign] = useState(null);

  // Mock Merkle Proof 数据（实际应该从后端或 IPFS 获取）
  const [merkleData, setMerkleData] = useState({});

  useEffect(() => {
    loadCampaigns();
  }, [address]);

  // 加载活动列表
  const loadCampaigns = async () => {
    setLoading(true);
    try {
      // TODO: 从 Subgraph 或合约读取活动
      // 这里使用 Mock 数据作为示例
      const mockCampaigns = [
        {
          id: "0",
          name: "Yieldera Genesis Airdrop",
          description: "感谢早期用户的支持！每个地址可领取 1000-5000 YLD",
          totalBudget: ethers.parseEther("100000"),
          claimedAmount: ethers.parseEther("35000"),
          startTime: Date.now() / 1000 - 86400 * 7, // 7天前
          endTime: Date.now() / 1000 + 86400 * 23, // 23天后
          isActive: true,
          participantCount: 234,
          hasClaimed: false,
        },
        {
          id: "1",
          name: "DeFi Vault 用户奖励",
          description: "奖励在理财金库中质押超过 1000 USDC 的用户",
          totalBudget: ethers.parseEther("50000"),
          claimedAmount: ethers.parseEther("12000"),
          startTime: Date.now() / 1000 - 86400 * 3,
          endTime: Date.now() / 1000 + 86400 * 27,
          isActive: true,
          participantCount: 89,
          hasClaimed: false,
        },
        {
          id: "2",
          name: "RWA 投资者空投",
          description: "奖励在 RWA 商城购买过资产的用户",
          totalBudget: ethers.parseEther("30000"),
          claimedAmount: ethers.parseEther("30000"),
          startTime: Date.now() / 1000 - 86400 * 30,
          endTime: Date.now() / 1000 - 86400 * 1, // 已结束
          isActive: false,
          participantCount: 156,
          hasClaimed: false,
        },
      ];

      setCampaigns(mockCampaigns);

      // Mock Merkle Proof 数据（实际应该根据用户地址从后端获取）
      if (address) {
        setMerkleData({
          "0": {
            amount: ethers.parseEther("2500"),
            proof: [
              "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
              "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
            ],
          },
          "1": {
            amount: ethers.parseEther("1000"),
            proof: [
              "0x9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba",
            ],
          },
        });
      }
    } catch (error) {
      console.error("加载活动失败:", error);
    }
    setLoading(false);
  };

  // 领取空投
  const handleClaim = async (campaign) => {
    if (!address || !signer) {
      alert("请先连接钱包");
      return;
    }

    const eligibility = merkleData[campaign.id];
    if (!eligibility) {
      alert("您不在此活动的白名单中");
      return;
    }

    setClaiming(true);
    setSelectedCampaign(campaign.id);

    try {
      // 创建合约实例
      const contract = new ethers.Contract(
        AIRDROP_CONTRACT_ADDRESS,
        AIRDROP_ABI,
        signer
      );

      // 发送领取交易
      const tx = await contract.claim(
        campaign.id,
        eligibility.amount,
        eligibility.proof
      );

      console.log("交易已发送:", tx.hash);
      alert(`交易已提交！\n哈希: ${tx.hash.slice(0, 10)}...`);

      // 等待确认
      await tx.wait();
      alert("领取成功！");

      // 重新加载数据
      await loadCampaigns();
    } catch (error) {
      console.error("领取失败:", error);
      if (error.code === "ACTION_REJECTED") {
        alert("用户取消了交易");
      } else if (error.message.includes("Already claimed")) {
        alert("您已经领取过此活动的空投");
      } else {
        alert(`领取失败: ${error.message}`);
      }
    } finally {
      setClaiming(false);
      setSelectedCampaign(null);
    }
  };

  // 格式化时间
  const formatDate = (timestamp) => {
    return new Date(timestamp * 1000).toLocaleDateString("zh-CN", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
    });
  };

  // 计算进度百分比
  const getProgress = (campaign) => {
    if (!campaign.totalBudget || campaign.totalBudget === 0n) return 0;
    return Number((campaign.claimedAmount * 10000n) / campaign.totalBudget) / 100;
  };

  if (loading) {
    return (
      <div style={{ textAlign: "center", padding: "60px 20px" }}>
        <div style={{ fontSize: "32px" }}>⏳</div>
        <div style={{ marginTop: "16px", color: "#9ca3af" }}>加载中...</div>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: "1200px", margin: "0 auto", padding: "40px 20px" }}>
      {/* 页面标题 */}
      <div style={{ marginBottom: "40px" }}>
        <h1
          style={{
            fontSize: "36px",
            fontWeight: "900",
            margin: "0 0 12px 0",
            background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
            WebkitBackgroundClip: "text",
            WebkitTextFillColor: "transparent",
          }}
        >
          链上空投 🎁
        </h1>
        <p style={{ fontSize: "16px", color: "#9ca3af", margin: 0 }}>
          基于 Merkle Tree 的链上空投系统，透明可验证，去中心化分发
        </p>
      </div>

      {/* 连接钱包提示 */}
      {!address && (
        <div
          style={{
            background: "rgba(251, 191, 36, 0.1)",
            border: "1px solid rgba(251, 191, 36, 0.3)",
            borderRadius: "12px",
            padding: "20px",
            marginBottom: "30px",
            textAlign: "center",
          }}
        >
          <div style={{ fontSize: "24px", marginBottom: "8px" }}>⚠️</div>
          <div style={{ color: "#fbbf24", fontWeight: "600" }}>
            请先连接钱包以查看您的空投资格
          </div>
        </div>
      )}

      {/* 活动列表 */}
      <div style={{ display: "grid", gap: "24px" }}>
        {campaigns.map((campaign) => {
          const isEligible = address && merkleData[campaign.id];
          const progress = getProgress(campaign);
          const isEnded = campaign.endTime * 1000 < Date.now();
          const isNotStarted = campaign.startTime * 1000 > Date.now();

          return (
            <div
              key={campaign.id}
              style={{
                background: "#1f2937",
                borderRadius: "16px",
                padding: "30px",
                border: isEligible
                  ? "2px solid rgba(102, 126, 234, 0.5)"
                  : "1px solid #374151",
              }}
            >
              {/* 活动头部 */}
              <div
                style={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "flex-start",
                  marginBottom: "20px",
                }}
              >
                <div style={{ flex: 1 }}>
                  <h3
                    style={{
                      fontSize: "24px",
                      fontWeight: "700",
                      margin: "0 0 8px 0",
                      color: "#fff",
                    }}
                  >
                    {campaign.name}
                  </h3>
                  <p
                    style={{
                      fontSize: "14px",
                      color: "#9ca3af",
                      margin: "0 0 16px 0",
                    }}
                  >
                    {campaign.description}
                  </p>

                  {/* 状态标签 */}
                  <div style={{ display: "flex", gap: "8px", flexWrap: "wrap" }}>
                    {isEligible && (
                      <span
                        style={{
                          background: "rgba(34, 197, 94, 0.1)",
                          color: "#22c55e",
                          padding: "4px 12px",
                          borderRadius: "6px",
                          fontSize: "12px",
                          fontWeight: "600",
                        }}
                      >
                        ✓ 符合资格
                      </span>
                    )}
                    {campaign.hasClaimed && (
                      <span
                        style={{
                          background: "rgba(156, 163, 175, 0.1)",
                          color: "#9ca3af",
                          padding: "4px 12px",
                          borderRadius: "6px",
                          fontSize: "12px",
                          fontWeight: "600",
                        }}
                      >
                        已领取
                      </span>
                    )}
                    {isEnded && (
                      <span
                        style={{
                          background: "rgba(239, 68, 68, 0.1)",
                          color: "#ef4444",
                          padding: "4px 12px",
                          borderRadius: "6px",
                          fontSize: "12px",
                          fontWeight: "600",
                        }}
                      >
                        已结束
                      </span>
                    )}
                    {isNotStarted && (
                      <span
                        style={{
                          background: "rgba(251, 191, 36, 0.1)",
                          color: "#fbbf24",
                          padding: "4px 12px",
                          borderRadius: "6px",
                          fontSize: "12px",
                          fontWeight: "600",
                        }}
                      >
                        未开始
                      </span>
                    )}
                  </div>
                </div>

                {/* 领取金额 */}
                {isEligible && (
                  <div
                    style={{
                      textAlign: "right",
                      background: "rgba(102, 126, 234, 0.1)",
                      padding: "16px 20px",
                      borderRadius: "12px",
                      minWidth: "180px",
                    }}
                  >
                    <div
                      style={{ fontSize: "12px", color: "#9ca3af", marginBottom: "4px" }}
                    >
                      可领取金额
                    </div>
                    <div
                      style={{
                        fontSize: "28px",
                        fontWeight: "700",
                        color: "#667eea",
                      }}
                    >
                      {ethers.formatEther(merkleData[campaign.id].amount)}
                    </div>
                    <div style={{ fontSize: "12px", color: "#9ca3af" }}>YLD</div>
                  </div>
                )}
              </div>

              {/* 进度条 */}
              <div style={{ marginBottom: "20px" }}>
                <div
                  style={{
                    display: "flex",
                    justifyContent: "space-between",
                    marginBottom: "8px",
                  }}
                >
                  <span style={{ fontSize: "14px", color: "#9ca3af" }}>
                    已分发进度
                  </span>
                  <span style={{ fontSize: "14px", color: "#fff", fontWeight: "600" }}>
                    {progress.toFixed(1)}%
                  </span>
                </div>
                <div
                  style={{
                    background: "#374151",
                    height: "8px",
                    borderRadius: "4px",
                    overflow: "hidden",
                  }}
                >
                  <div
                    style={{
                      background: "linear-gradient(90deg, #667eea 0%, #764ba2 100%)",
                      height: "100%",
                      width: `${progress}%`,
                      transition: "width 0.3s ease",
                    }}
                  />
                </div>
                <div
                  style={{
                    display: "flex",
                    justifyContent: "space-between",
                    marginTop: "8px",
                  }}
                >
                  <span style={{ fontSize: "12px", color: "#6b7280" }}>
                    {ethers.formatEther(campaign.claimedAmount)} YLD 已领取
                  </span>
                  <span style={{ fontSize: "12px", color: "#6b7280" }}>
                    {ethers.formatEther(campaign.totalBudget)} YLD 总额
                  </span>
                </div>
              </div>

              {/* 活动信息 */}
              <div
                style={{
                  display: "grid",
                  gridTemplateColumns: "repeat(3, 1fr)",
                  gap: "16px",
                  marginBottom: "20px",
                }}
              >
                <div>
                  <div style={{ fontSize: "12px", color: "#6b7280", marginBottom: "4px" }}>
                    开始时间
                  </div>
                  <div style={{ fontSize: "14px", color: "#fff" }}>
                    {formatDate(campaign.startTime)}
                  </div>
                </div>
                <div>
                  <div style={{ fontSize: "12px", color: "#6b7280", marginBottom: "4px" }}>
                    结束时间
                  </div>
                  <div style={{ fontSize: "14px", color: "#fff" }}>
                    {formatDate(campaign.endTime)}
                  </div>
                </div>
                <div>
                  <div style={{ fontSize: "12px", color: "#6b7280", marginBottom: "4px" }}>
                    参与人数
                  </div>
                  <div style={{ fontSize: "14px", color: "#fff" }}>
                    {campaign.participantCount} 人
                  </div>
                </div>
              </div>

              {/* 领取按钮 */}
              {isEligible && !campaign.hasClaimed && !isEnded && !isNotStarted && (
                <button
                  onClick={() => handleClaim(campaign)}
                  disabled={claiming && selectedCampaign === campaign.id}
                  className="btn"
                  style={{
                    width: "100%",
                    padding: "14px",
                    fontSize: "16px",
                    fontWeight: "600",
                    background:
                      claiming && selectedCampaign === campaign.id
                        ? "#6b7280"
                        : "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
                    cursor:
                      claiming && selectedCampaign === campaign.id
                        ? "not-allowed"
                        : "pointer",
                  }}
                >
                  {claiming && selectedCampaign === campaign.id
                    ? "领取中..."
                    : `立即领取 ${ethers.formatEther(
                        merkleData[campaign.id].amount
                      )} YLD`}
                </button>
              )}

              {/* 不符合资格提示 */}
              {address && !isEligible && (
                <div
                  style={{
                    background: "rgba(107, 114, 128, 0.1)",
                    border: "1px solid rgba(107, 114, 128, 0.3)",
                    borderRadius: "8px",
                    padding: "12px",
                    textAlign: "center",
                    color: "#9ca3af",
                    fontSize: "14px",
                  }}
                >
                  您不在此活动的白名单中
                </div>
              )}
            </div>
          );
        })}
      </div>

      {/* 空状态 */}
      {campaigns.length === 0 && (
        <div style={{ textAlign: "center", padding: "60px 20px" }}>
          <div style={{ fontSize: "64px", marginBottom: "16px" }}>📭</div>
          <div style={{ fontSize: "18px", color: "#9ca3af" }}>暂无活动</div>
        </div>
      )}
    </div>
  );
}
