// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "../interfaces/treasury/ITreasuryAsset.sol";

/**
 * @title TreasuryYieldDistributor
 * @notice Distribute coupon payments and maturity proceeds to treasury token holders
 * @dev Supports both proportional distribution and claiming mechanisms
 *
 * Key Features:
 * - Periodic coupon payment distribution
 * - Maturity redemption
 * - Yield tracking per user
 * - Claim mechanism for users
 * - Batch distribution for gas efficiency
 */
contract TreasuryYieldDistributor is AccessControl, ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;

    bytes32 public constant YIELD_MANAGER_ROLE = keccak256("YIELD_MANAGER_ROLE");
    bytes32 public constant DISTRIBUTOR_ROLE = keccak256("DISTRIBUTOR_ROLE");

    /// @notice Yield distribution record
    struct YieldDistribution {
        uint256 distributionId;
        uint256 assetId;
        uint256 totalYield;          // Total USD amount
        uint256 yieldPerToken;       // Yield per token (18 decimals)
        uint256 totalTokens;         // Total tokens at snapshot
        uint256 recipientsCount;
        uint256 distributedAmount;
        bool isCompleted;
        uint256 timestamp;
        string distributionType;     // "COUPON" or "MATURITY"
    }

    /// @notice User yield tracking
    struct UserYield {
        uint256 totalEarned;
        uint256 totalClaimed;
        uint256 pendingYield;
        uint256 lastClaimTime;
    }

    /// @notice Distribution counter
    uint256 private distributionIdCounter;

    /// @notice Storage
    mapping(uint256 => YieldDistribution) public distributions; // distributionId => distribution
    mapping(uint256 => uint256[]) public assetDistributions;    // assetId => distributionIds
    mapping(uint256 => mapping(address => uint256)) public userLastDistribution; // assetId => user => lastDistId
    mapping(uint256 => mapping(address => UserYield)) public userYields; // assetId => user => yield

    /// @notice Integration contracts
    ITreasuryAsset public immutable assetFactory;
    address public immutable yieldToken; // USDC or stablecoin for yield payments

    /// @notice Statistics
    uint256 public totalDistributions;
    uint256 public totalYieldDistributed;

    /// @notice Events
    event YieldDeposited(
        uint256 indexed distributionId,
        uint256 indexed assetId,
        uint256 totalYield,
        uint256 yieldPerToken,
        string distributionType
    );

    event YieldClaimed(
        address indexed user,
        uint256 indexed assetId,
        uint256 amount,
        uint256 distributionId
    );

    event BatchDistributed(
        uint256 indexed distributionId,
        uint256 indexed assetId,
        uint256 recipientsCount,
        uint256 totalAmount
    );

    /**
     * @notice Constructor
     * @param admin Admin address
     * @param factory TreasuryAssetFactory address
     * @param yieldToken_ Yield payment token (USDC)
     */
    constructor(
        address admin,
        address factory,
        address yieldToken_
    ) {
        require(admin != address(0), "Invalid admin");
        require(factory != address(0), "Invalid factory");
        require(yieldToken_ != address(0), "Invalid yield token");

        assetFactory = ITreasuryAsset(factory);
        yieldToken = yieldToken_;

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(YIELD_MANAGER_ROLE, admin);
        _grantRole(DISTRIBUTOR_ROLE, admin);
    }

    // =============================================================
    //                     YIELD DEPOSIT
    // =============================================================

    /**
     * @notice Deposit yield for distribution
     * @param assetId Treasury asset ID
     * @param totalYield Total yield amount in USD
     * @param distributionType Type of distribution ("COUPON" or "MATURITY")
     * @return distributionId Distribution identifier
     */
    function depositYield(
        uint256 assetId,
        uint256 totalYield,
        string memory distributionType
    ) external onlyRole(YIELD_MANAGER_ROLE) whenNotPaused nonReentrant returns (uint256 distributionId) {
        require(assetFactory.isAssetActive(assetId), "Asset not active");
        require(totalYield > 0, "Invalid yield amount");

        // Get treasury token
        address tokenAddress = assetFactory.getTokenAddress(assetId);
        require(tokenAddress != address(0), "Token not found");

        IERC20 treasuryToken = IERC20(tokenAddress);
        uint256 totalTokens = treasuryToken.totalSupply();
        require(totalTokens > 0, "No tokens issued");

        // Transfer yield tokens to contract
        IERC20(yieldToken).safeTransferFrom(msg.sender, address(this), totalYield);

        // Calculate yield per token
        uint256 yieldPerToken = (totalYield * 1e18) / totalTokens;

        // Create distribution record
        distributionIdCounter++;
        distributionId = distributionIdCounter;

        YieldDistribution storage dist = distributions[distributionId];
        dist.distributionId = distributionId;
        dist.assetId = assetId;
        dist.totalYield = totalYield;
        dist.yieldPerToken = yieldPerToken;
        dist.totalTokens = totalTokens;
        dist.recipientsCount = 0;
        dist.distributedAmount = 0;
        dist.isCompleted = false;
        dist.timestamp = block.timestamp;
        dist.distributionType = distributionType;

        assetDistributions[assetId].push(distributionId);
        totalDistributions++;

        emit YieldDeposited(distributionId, assetId, totalYield, yieldPerToken, distributionType);
    }

    // =============================================================
    //                     YIELD CLAIMING
    // =============================================================

    /**
     * @notice Claim pending yield for user
     * @param assetId Treasury asset ID
     * @return claimedAmount Amount claimed
     */
    function claimYield(uint256 assetId) external whenNotPaused nonReentrant returns (uint256 claimedAmount) {
        address user = msg.sender;

        // Calculate pending yield
        uint256 pending = calculatePendingYield(assetId, user);
        require(pending > 0, "No pending yield");

        // Get treasury token balance
        address tokenAddress = assetFactory.getTokenAddress(assetId);
        IERC20 treasuryToken = IERC20(tokenAddress);
        uint256 userBalance = treasuryToken.balanceOf(user);
        require(userBalance > 0, "No token balance");

        // Update user yield tracking
        UserYield storage userYield = userYields[assetId][user];
        userYield.totalClaimed += pending;
        userYield.pendingYield = 0;
        userYield.lastClaimTime = block.timestamp;

        // Update last processed distribution
        uint256[] storage distIds = assetDistributions[assetId];
        if (distIds.length > 0) {
            userLastDistribution[assetId][user] = distIds[distIds.length - 1];
        }

        // Transfer yield
        IERC20(yieldToken).safeTransfer(user, pending);

        totalYieldDistributed += pending;
        claimedAmount = pending;

        emit YieldClaimed(user, assetId, claimedAmount, userLastDistribution[assetId][user]);
    }

    /**
     * @notice Calculate pending yield for user
     * @param assetId Treasury asset ID
     * @param user User address
     * @return Pending yield amount
     */
    function calculatePendingYield(
        uint256 assetId,
        address user
    ) public view returns (uint256) {
        address tokenAddress = assetFactory.getTokenAddress(assetId);
        if (tokenAddress == address(0)) return 0;

        IERC20 treasuryToken = IERC20(tokenAddress);
        uint256 userBalance = treasuryToken.balanceOf(user);
        if (userBalance == 0) return 0;

        uint256[] storage distIds = assetDistributions[assetId];
        if (distIds.length == 0) return 0;

        uint256 lastProcessed = userLastDistribution[assetId][user];
        uint256 pendingYield = 0;

        // Sum up yield from all distributions after last processed
        for (uint256 i = 0; i < distIds.length; i++) {
            uint256 distId = distIds[i];
            if (distId <= lastProcessed) continue;

            YieldDistribution storage dist = distributions[distId];
            uint256 userYieldAmount = (userBalance * dist.yieldPerToken) / 1e18;
            pendingYield += userYieldAmount;
        }

        return pendingYield;
    }

    // =============================================================
    //                     BATCH DISTRIBUTION
    // =============================================================

    /**
     * @notice Batch distribute yield to multiple recipients
     * @param distributionId Distribution identifier
     * @param recipients Array of recipient addresses
     * @param amounts Array of amounts to distribute
     */
    function batchDistribute(
        uint256 distributionId,
        address[] calldata recipients,
        uint256[] calldata amounts
    ) external onlyRole(DISTRIBUTOR_ROLE) whenNotPaused nonReentrant {
        require(recipients.length == amounts.length, "Length mismatch");
        require(recipients.length > 0, "Empty recipients");

        YieldDistribution storage dist = distributions[distributionId];
        require(dist.distributionId != 0, "Distribution not found");
        require(!dist.isCompleted, "Distribution completed");

        uint256 totalAmount = 0;

        for (uint256 i = 0; i < recipients.length; i++) {
            address recipient = recipients[i];
            uint256 amount = amounts[i];

            require(recipient != address(0), "Invalid recipient");
            require(amount > 0, "Invalid amount");

            // Transfer yield
            IERC20(yieldToken).safeTransfer(recipient, amount);

            // Update user yield tracking
            UserYield storage userYield = userYields[dist.assetId][recipient];
            userYield.totalEarned += amount;
            userYield.totalClaimed += amount;
            userYield.lastClaimTime = block.timestamp;

            totalAmount += amount;

            emit YieldClaimed(recipient, dist.assetId, amount, distributionId);
        }

        // Update distribution
        dist.recipientsCount += recipients.length;
        dist.distributedAmount += totalAmount;

        // Check if distribution is complete
        if (dist.distributedAmount >= dist.totalYield) {
            dist.isCompleted = true;
        }

        totalYieldDistributed += totalAmount;

        emit BatchDistributed(distributionId, dist.assetId, recipients.length, totalAmount);
    }

    // =============================================================
    //                      VIEW FUNCTIONS
    // =============================================================

    /**
     * @notice Get distribution details
     * @param distributionId Distribution identifier
     * @return Distribution struct
     */
    function getDistribution(uint256 distributionId)
        external
        view
        returns (YieldDistribution memory)
    {
        return distributions[distributionId];
    }

    /**
     * @notice Get all distributions for asset
     * @param assetId Asset identifier
     * @return Array of distribution IDs
     */
    function getAssetDistributions(uint256 assetId)
        external
        view
        returns (uint256[] memory)
    {
        return assetDistributions[assetId];
    }

    /**
     * @notice Get user yield info
     * @param assetId Asset identifier
     * @param user User address
     * @return User yield struct
     */
    function getUserYieldInfo(uint256 assetId, address user)
        external
        view
        returns (UserYield memory)
    {
        return userYields[assetId][user];
    }

    /**
     * @notice Get user's total earnings and pending yield
     * @param assetId Asset identifier
     * @param user User address
     * @return totalEarned Total earned (historical)
     * @return pendingYield Current pending yield
     * @return totalClaimed Total claimed
     */
    function getUserYieldSummary(uint256 assetId, address user)
        external
        view
        returns (
            uint256 totalEarned,
            uint256 pendingYield,
            uint256 totalClaimed
        )
    {
        UserYield storage userYield = userYields[assetId][user];
        uint256 pending = calculatePendingYield(assetId, user);

        return (
            userYield.totalEarned + pending,
            pending,
            userYield.totalClaimed
        );
    }

    // =============================================================
    //                     ADMIN FUNCTIONS
    // =============================================================

    /**
     * @notice Mark distribution as completed (emergency)
     * @param distributionId Distribution identifier
     */
    function markDistributionCompleted(uint256 distributionId)
        external
        onlyRole(YIELD_MANAGER_ROLE)
    {
        YieldDistribution storage dist = distributions[distributionId];
        require(dist.distributionId != 0, "Distribution not found");
        dist.isCompleted = true;
    }

    /**
     * @notice Pause contract
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause contract
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }

    /**
     * @notice Emergency withdraw (only undistributed yield)
     * @param token Token address
     * @param amount Amount to withdraw
     * @param recipient Recipient address
     */
    function emergencyWithdraw(
        address token,
        uint256 amount,
        address recipient
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(recipient != address(0), "Invalid recipient");
        IERC20(token).safeTransfer(recipient, amount);
    }
}
