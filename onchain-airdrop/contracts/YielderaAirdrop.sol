// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";

/**
 * @title YielderaAirdrop
 * @notice Yieldera 链上空投分发合约
 * @dev 支持 Merkle Tree 验证，Gas 优化，批量操作
 */
contract YielderaAirdrop is Ownable, ReentrancyGuard {

    // ==================== 状态变量 ====================

    /// @notice 空投代币（Yieldera Token）
    IERC20 public immutable token;

    /// @notice 活动ID计数器
    uint256 public campaignIdCounter;

    /// @notice 活动信息
    struct Campaign {
        string name;                // 活动名称
        string description;         // 活动描述
        bytes32 merkleRoot;         // Merkle 树根
        uint256 totalBudget;        // 总预算
        uint256 claimedAmount;      // 已领取金额
        uint256 startTime;          // 开始时间
        uint256 endTime;            // 结束时间
        bool isActive;              // 是否激活
        uint256 participantCount;   // 参与人数
        address creator;            // 创建者
    }

    /// @notice 活动ID => 活动信息
    mapping(uint256 => Campaign) public campaigns;

    /// @notice 活动ID => 用户地址 => 是否已领取
    mapping(uint256 => mapping(address => bool)) public hasClaimed;

    /// @notice 活动ID => 用户地址 => 领取金额
    mapping(uint256 => mapping(address => uint256)) public claimedAmounts;

    // ==================== 事件 ====================

    /// @notice 活动创建事件
    event CampaignCreated(
        uint256 indexed campaignId,
        string name,
        uint256 totalBudget,
        uint256 startTime,
        uint256 endTime,
        address indexed creator
    );

    /// @notice 领取成功事件（Subgraph 监听这个）
    event Claimed(
        uint256 indexed campaignId,
        address indexed user,
        uint256 amount,
        uint256 timestamp
    );

    /// @notice 活动状态更新事件
    event CampaignStatusUpdated(
        uint256 indexed campaignId,
        bool isActive
    );

    /// @notice 紧急提取事件
    event EmergencyWithdraw(
        uint256 indexed campaignId,
        uint256 amount,
        address indexed to
    );

    // ==================== 修饰符 ====================

    /// @notice 检查活动是否存在
    modifier campaignExists(uint256 campaignId) {
        require(campaignId < campaignIdCounter, "Campaign does not exist");
        _;
    }

    /// @notice 检查活动是否激活且在有效期内
    modifier campaignActive(uint256 campaignId) {
        Campaign storage campaign = campaigns[campaignId];
        require(campaign.isActive, "Campaign not active");
        require(block.timestamp >= campaign.startTime, "Campaign not started");
        require(block.timestamp <= campaign.endTime, "Campaign ended");
        _;
    }

    // ==================== 构造函数 ====================

    constructor(address _token) {
        require(_token != address(0), "Invalid token address");
        token = IERC20(_token);
    }

    // ==================== 管理员函数 ====================

    /**
     * @notice 创建空投活动
     * @param name 活动名称
     * @param description 活动描述
     * @param merkleRoot Merkle树根
     * @param totalBudget 总预算
     * @param startTime 开始时间
     * @param endTime 结束时间
     */
    function createCampaign(
        string memory name,
        string memory description,
        bytes32 merkleRoot,
        uint256 totalBudget,
        uint256 startTime,
        uint256 endTime
    ) external onlyOwner returns (uint256) {
        require(totalBudget > 0, "Budget must be > 0");
        require(endTime > startTime, "Invalid time range");
        require(merkleRoot != bytes32(0), "Invalid merkle root");

        uint256 campaignId = campaignIdCounter++;

        campaigns[campaignId] = Campaign({
            name: name,
            description: description,
            merkleRoot: merkleRoot,
            totalBudget: totalBudget,
            claimedAmount: 0,
            startTime: startTime,
            endTime: endTime,
            isActive: true,
            participantCount: 0,
            creator: msg.sender
        });

        // 转入代币到合约
        require(
            token.transferFrom(msg.sender, address(this), totalBudget),
            "Token transfer failed"
        );

        emit CampaignCreated(
            campaignId,
            name,
            totalBudget,
            startTime,
            endTime,
            msg.sender
        );

        return campaignId;
    }

    /**
     * @notice 更新活动状态
     * @param campaignId 活动ID
     * @param isActive 是否激活
     */
    function updateCampaignStatus(
        uint256 campaignId,
        bool isActive
    ) external onlyOwner campaignExists(campaignId) {
        campaigns[campaignId].isActive = isActive;
        emit CampaignStatusUpdated(campaignId, isActive);
    }

    /**
     * @notice 紧急提取剩余代币
     * @param campaignId 活动ID
     */
    function emergencyWithdraw(
        uint256 campaignId
    ) external onlyOwner campaignExists(campaignId) {
        Campaign storage campaign = campaigns[campaignId];

        uint256 remainingAmount = campaign.totalBudget - campaign.claimedAmount;
        require(remainingAmount > 0, "No tokens to withdraw");

        campaign.isActive = false;

        require(
            token.transfer(owner(), remainingAmount),
            "Token transfer failed"
        );

        emit EmergencyWithdraw(campaignId, remainingAmount, owner());
    }

    // ==================== 用户函数 ====================

    /**
     * @notice 领取空投（单个活动）
     * @param campaignId 活动ID
     * @param amount 领取金额
     * @param merkleProof Merkle证明
     */
    function claim(
        uint256 campaignId,
        uint256 amount,
        bytes32[] calldata merkleProof
    )
        external
        nonReentrant
        campaignExists(campaignId)
        campaignActive(campaignId)
    {
        require(!hasClaimed[campaignId][msg.sender], "Already claimed");
        require(amount > 0, "Amount must be > 0");

        Campaign storage campaign = campaigns[campaignId];

        // 验证 Merkle Proof
        bytes32 leaf = keccak256(abi.encodePacked(msg.sender, amount));
        require(
            MerkleProof.verify(merkleProof, campaign.merkleRoot, leaf),
            "Invalid merkle proof"
        );

        // 检查预算
        require(
            campaign.claimedAmount + amount <= campaign.totalBudget,
            "Exceeds budget"
        );

        // 更新状态
        hasClaimed[campaignId][msg.sender] = true;
        claimedAmounts[campaignId][msg.sender] = amount;
        campaign.claimedAmount += amount;
        campaign.participantCount += 1;

        // 转账代币
        require(
            token.transfer(msg.sender, amount),
            "Token transfer failed"
        );

        // 触发事件（Subgraph 监听）
        emit Claimed(campaignId, msg.sender, amount, block.timestamp);
    }

    /**
     * @notice 批量领取空投（多个活动）
     * @param campaignIds 活动ID数组
     * @param amounts 金额数组
     * @param merkleProofs Merkle证明数组
     */
    function claimMultiple(
        uint256[] calldata campaignIds,
        uint256[] calldata amounts,
        bytes32[][] calldata merkleProofs
    ) external nonReentrant {
        require(
            campaignIds.length == amounts.length &&
            amounts.length == merkleProofs.length,
            "Array length mismatch"
        );

        for (uint256 i = 0; i < campaignIds.length; i++) {
            uint256 campaignId = campaignIds[i];

            // 检查活动状态
            if (campaignId >= campaignIdCounter) continue;
            if (hasClaimed[campaignId][msg.sender]) continue;

            Campaign storage campaign = campaigns[campaignId];
            if (!campaign.isActive) continue;
            if (block.timestamp < campaign.startTime) continue;
            if (block.timestamp > campaign.endTime) continue;

            // 验证 Merkle Proof
            bytes32 leaf = keccak256(abi.encodePacked(msg.sender, amounts[i]));
            if (!MerkleProof.verify(merkleProofs[i], campaign.merkleRoot, leaf)) {
                continue;
            }

            // 检查预算
            if (campaign.claimedAmount + amounts[i] > campaign.totalBudget) {
                continue;
            }

            // 更新状态
            hasClaimed[campaignId][msg.sender] = true;
            claimedAmounts[campaignId][msg.sender] = amounts[i];
            campaign.claimedAmount += amounts[i];
            campaign.participantCount += 1;

            // 转账代币
            require(
                token.transfer(msg.sender, amounts[i]),
                "Token transfer failed"
            );

            // 触发事件
            emit Claimed(campaignId, msg.sender, amounts[i], block.timestamp);
        }
    }

    // ==================== 查询函数 ====================

    /**
     * @notice 检查用户是否已领取
     * @param campaignId 活动ID
     * @param user 用户地址
     */
    function hasUserClaimed(
        uint256 campaignId,
        address user
    ) external view returns (bool) {
        return hasClaimed[campaignId][user];
    }

    /**
     * @notice 获取活动信息
     * @param campaignId 活动ID
     */
    function getCampaign(
        uint256 campaignId
    ) external view campaignExists(campaignId) returns (Campaign memory) {
        return campaigns[campaignId];
    }

    /**
     * @notice 获取活动剩余预算
     * @param campaignId 活动ID
     */
    function getRemainingBudget(
        uint256 campaignId
    ) external view campaignExists(campaignId) returns (uint256) {
        Campaign storage campaign = campaigns[campaignId];
        return campaign.totalBudget - campaign.claimedAmount;
    }

    /**
     * @notice 获取用户领取金额
     * @param campaignId 活动ID
     * @param user 用户地址
     */
    function getUserClaimedAmount(
        uint256 campaignId,
        address user
    ) external view returns (uint256) {
        return claimedAmounts[campaignId][user];
    }
}
