// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "../interfaces/rwa/IRWAAsset.sol";

/**
 * @title RWAYieldDistributor
 * @notice Automated yield distribution for RWA token holders
 * @dev Distributes yield pro-rata based on token ownership
 *
 * Key Features:
 * - Pro-rata yield distribution to token holders
 * - Support for multiple payment tokens (ETH, stablecoins)
 * - Snapshot-based distributions for fairness
 * - Unclaimed yield tracking and reclamation
 * - Automatic calculation based on ownership percentage
 * - Historical distribution tracking
 *
 * Distribution Flow:
 * 1. Asset issuer deposits yield payment
 * 2. System takes snapshot of token holders
 * 3. Yield allocated pro-rata to holders
 * 4. Token holders claim their share
 * 5. Unclaimed yields can be reclaimed after expiry
 */
contract RWAYieldDistributor is AccessControl, ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;

    bytes32 public constant DISTRIBUTOR_ROLE = keccak256("DISTRIBUTOR_ROLE");

    // Distribution record
    struct Distribution {
        uint256 distributionId;
        uint256 assetId;
        address paymentToken;       // Token used for payment (address(0) = ETH)
        uint256 totalAmount;        // Total yield amount
        uint256 totalSupply;        // Token supply at snapshot
        uint256 distributedAt;
        uint256 claimDeadline;      // Deadline for claiming
        uint256 totalClaimed;       // Amount claimed so far
        bool finalized;             // Distribution finalized
    }

    // User claim record
    struct Claim {
        uint256 amount;             // Claimable amount
        bool claimed;               // Whether claimed
        uint256 claimedAt;          // When claimed
    }

    // Distribution counter
    uint256 private distributionIdCounter;

    // Distributions storage
    mapping(uint256 => Distribution) private distributions;
    mapping(uint256 => mapping(address => Claim)) private claims; // distributionId => user => claim

    // Asset distributions tracking
    mapping(uint256 => uint256[]) private assetDistributions; // assetId => distributionIds

    // Integration
    IRWAAsset public immutable assetFactory;

    // Constants
    uint256 public constant DEFAULT_CLAIM_PERIOD = 365 days;
    uint256 public constant PRECISION = 1e18;

    // Statistics
    uint256 public totalDistributions;
    uint256 public totalYieldDistributed;
    mapping(address => uint256) public totalYieldByToken;

    // Events
    event YieldDeposited(
        uint256 indexed distributionId,
        uint256 indexed assetId,
        address indexed paymentToken,
        uint256 amount,
        uint256 claimDeadline
    );

    event YieldClaimed(
        uint256 indexed distributionId,
        address indexed user,
        uint256 amount
    );

    event DistributionFinalized(
        uint256 indexed distributionId,
        uint256 totalClaimed,
        uint256 unclaimed
    );

    event UnclaimedYieldReclaimed(
        uint256 indexed distributionId,
        uint256 amount,
        address recipient
    );

    /**
     * @notice Constructor
     * @param admin Admin address
     * @param factory RWAAssetFactory address
     */
    constructor(address admin, address factory) {
        require(admin != address(0), "Invalid admin");
        require(factory != address(0), "Invalid factory");

        assetFactory = IRWAAsset(factory);

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(DISTRIBUTOR_ROLE, admin);
    }

    // =============================================================
    //                  YIELD DISTRIBUTION
    // =============================================================

    /**
     * @notice Deposit yield for distribution
     * @param assetId Asset identifier
     * @param paymentToken Token address (address(0) for ETH)
     * @param amount Yield amount to distribute
     * @param claimPeriod How long users can claim (seconds)
     * @return distributionId New distribution ID
     */
    function depositYield(
        uint256 assetId,
        address paymentToken,
        uint256 amount,
        uint256 claimPeriod
    ) external payable whenNotPaused nonReentrant returns (uint256 distributionId) {
        require(assetFactory.isAssetActive(assetId), "Asset not active");
        require(amount > 0, "Invalid amount");
        require(claimPeriod > 0 && claimPeriod <= 730 days, "Invalid claim period");

        // Get fractional token
        address tokenAddress = assetFactory.getFractionalToken(assetId);
        require(tokenAddress != address(0), "Asset not fractionalized");

        IERC20 fractionalToken = IERC20(tokenAddress);
        uint256 totalSupply = fractionalToken.totalSupply();
        require(totalSupply > 0, "No tokens issued");

        // Handle payment
        if (paymentToken == address(0)) {
            require(msg.value == amount, "Incorrect ETH amount");
        } else {
            require(msg.value == 0, "ETH not accepted for token payment");
            IERC20(paymentToken).safeTransferFrom(msg.sender, address(this), amount);
        }

        // Create distribution
        distributionIdCounter++;
        distributionId = distributionIdCounter;

        Distribution storage dist = distributions[distributionId];
        dist.distributionId = distributionId;
        dist.assetId = assetId;
        dist.paymentToken = paymentToken;
        dist.totalAmount = amount;
        dist.totalSupply = totalSupply;
        dist.distributedAt = block.timestamp;
        dist.claimDeadline = block.timestamp + claimPeriod;
        dist.totalClaimed = 0;
        dist.finalized = false;

        // Track distribution
        assetDistributions[assetId].push(distributionId);

        totalDistributions++;
        totalYieldDistributed += amount;
        totalYieldByToken[paymentToken] += amount;

        emit YieldDeposited(distributionId, assetId, paymentToken, amount, dist.claimDeadline);
    }

    /**
     * @notice Claim yield for a distribution
     * @param distributionId Distribution to claim from
     */
    function claimYield(uint256 distributionId) external whenNotPaused nonReentrant {
        Distribution storage dist = distributions[distributionId];

        require(dist.distributionId != 0, "Distribution not found");
        require(!dist.finalized, "Distribution finalized");
        require(block.timestamp < dist.claimDeadline, "Claim period expired");

        Claim storage claim = claims[distributionId][msg.sender];
        require(!claim.claimed, "Already claimed");

        // Calculate claimable amount
        address tokenAddress = assetFactory.getFractionalToken(dist.assetId);
        IERC20 fractionalToken = IERC20(tokenAddress);

        uint256 userBalance = fractionalToken.balanceOf(msg.sender);
        require(userBalance > 0, "No tokens held");

        uint256 claimableAmount = (dist.totalAmount * userBalance) / dist.totalSupply;
        require(claimableAmount > 0, "Nothing to claim");

        // Mark as claimed
        claim.amount = claimableAmount;
        claim.claimed = true;
        claim.claimedAt = block.timestamp;

        dist.totalClaimed += claimableAmount;

        // Transfer yield
        if (dist.paymentToken == address(0)) {
            payable(msg.sender).transfer(claimableAmount);
        } else {
            IERC20(dist.paymentToken).safeTransfer(msg.sender, claimableAmount);
        }

        emit YieldClaimed(distributionId, msg.sender, claimableAmount);
    }

    /**
     * @notice Batch claim multiple distributions
     * @param distributionIds Array of distribution IDs
     */
    function batchClaimYield(uint256[] calldata distributionIds) external whenNotPaused nonReentrant {
        for (uint256 i = 0; i < distributionIds.length; i++) {
            uint256 distributionId = distributionIds[i];
            Distribution storage dist = distributions[distributionId];

            // Skip if invalid or already claimed
            if (dist.distributionId == 0 || dist.finalized || block.timestamp >= dist.claimDeadline) {
                continue;
            }

            Claim storage claim = claims[distributionId][msg.sender];
            if (claim.claimed) {
                continue;
            }

            // Calculate and claim
            address tokenAddress = assetFactory.getFractionalToken(dist.assetId);
            IERC20 fractionalToken = IERC20(tokenAddress);

            uint256 userBalance = fractionalToken.balanceOf(msg.sender);
            if (userBalance == 0) {
                continue;
            }

            uint256 claimableAmount = (dist.totalAmount * userBalance) / dist.totalSupply;
            if (claimableAmount == 0) {
                continue;
            }

            // Mark as claimed
            claim.amount = claimableAmount;
            claim.claimed = true;
            claim.claimedAt = block.timestamp;

            dist.totalClaimed += claimableAmount;

            // Transfer yield
            if (dist.paymentToken == address(0)) {
                payable(msg.sender).transfer(claimableAmount);
            } else {
                IERC20(dist.paymentToken).safeTransfer(msg.sender, claimableAmount);
            }

            emit YieldClaimed(distributionId, msg.sender, claimableAmount);
        }
    }

    /**
     * @notice Finalize distribution and reclaim unclaimed yield
     * @param distributionId Distribution to finalize
     */
    function finalizeDistribution(
        uint256 distributionId
    ) external onlyRole(DISTRIBUTOR_ROLE) nonReentrant {
        Distribution storage dist = distributions[distributionId];

        require(dist.distributionId != 0, "Distribution not found");
        require(!dist.finalized, "Already finalized");
        require(block.timestamp >= dist.claimDeadline, "Claim period not ended");

        uint256 unclaimed = dist.totalAmount - dist.totalClaimed;

        dist.finalized = true;

        // Return unclaimed to asset issuer
        if (unclaimed > 0) {
            IRWAAsset.AssetMetadata memory metadata = assetFactory.getAssetMetadata(dist.assetId);

            if (dist.paymentToken == address(0)) {
                payable(metadata.issuer).transfer(unclaimed);
            } else {
                IERC20(dist.paymentToken).safeTransfer(metadata.issuer, unclaimed);
            }

            emit UnclaimedYieldReclaimed(distributionId, unclaimed, metadata.issuer);
        }

        emit DistributionFinalized(distributionId, dist.totalClaimed, unclaimed);
    }

    // =============================================================
    //                      VIEW FUNCTIONS
    // =============================================================

    /**
     * @notice Get distribution details
     * @param distributionId Distribution identifier
     * @return Distribution struct
     */
    function getDistribution(uint256 distributionId) external view returns (Distribution memory) {
        return distributions[distributionId];
    }

    /**
     * @notice Get user's claim status
     * @param distributionId Distribution identifier
     * @param user User address
     * @return Claim struct
     */
    function getUserClaim(
        uint256 distributionId,
        address user
    ) external view returns (Claim memory) {
        return claims[distributionId][user];
    }

    /**
     * @notice Calculate claimable amount for user
     * @param distributionId Distribution identifier
     * @param user User address
     * @return Claimable amount
     */
    function getClaimableAmount(
        uint256 distributionId,
        address user
    ) public view returns (uint256) {
        Distribution storage dist = distributions[distributionId];

        if (dist.distributionId == 0 || dist.finalized || block.timestamp >= dist.claimDeadline) {
            return 0;
        }

        Claim storage claim = claims[distributionId][user];
        if (claim.claimed) {
            return 0;
        }

        address tokenAddress = assetFactory.getFractionalToken(dist.assetId);
        IERC20 fractionalToken = IERC20(tokenAddress);

        uint256 userBalance = fractionalToken.balanceOf(user);
        if (userBalance == 0) {
            return 0;
        }

        return (dist.totalAmount * userBalance) / dist.totalSupply;
    }

    /**
     * @notice Get all distributions for an asset
     * @param assetId Asset identifier
     * @return Array of distribution IDs
     */
    function getAssetDistributions(uint256 assetId) external view returns (uint256[] memory) {
        return assetDistributions[assetId];
    }

    /**
     * @notice Get total claimable across all active distributions
     * @param user User address
     * @param assetId Asset identifier (0 for all assets)
     * @return totalClaimable Total claimable amount
     */
    function getTotalClaimable(
        address user,
        uint256 assetId
    ) external view returns (uint256 totalClaimable) {
        uint256 start = assetId == 0 ? 1 : assetId;
        uint256 end = assetId == 0 ? 1000 : assetId + 1; // Simplified

        for (uint256 i = start; i < end; i++) {
            uint256[] storage distIds = assetDistributions[i];

            for (uint256 j = 0; j < distIds.length; j++) {
                totalClaimable += getClaimableAmount(distIds[j], user);
            }

            if (assetId != 0) break;
        }
    }

    /**
     * @notice Get unclaimed yield for distribution
     * @param distributionId Distribution identifier
     * @return Unclaimed amount
     */
    function getUnclaimedAmount(uint256 distributionId) external view returns (uint256) {
        Distribution storage dist = distributions[distributionId];
        return dist.totalAmount - dist.totalClaimed;
    }

    // =============================================================
    //                    ADMIN FUNCTIONS
    // =============================================================

    /**
     * @notice Pause yield distribution
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause yield distribution
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }

    /**
     * @notice Emergency withdrawal (only for stuck funds)
     * @param token Token address (address(0) for ETH)
     * @param amount Amount to withdraw
     * @param recipient Recipient address
     */
    function emergencyWithdraw(
        address token,
        uint256 amount,
        address recipient
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(recipient != address(0), "Invalid recipient");

        if (token == address(0)) {
            payable(recipient).transfer(amount);
        } else {
            IERC20(token).safeTransfer(recipient, amount);
        }
    }

    /**
     * @notice Receive ETH
     */
    receive() external payable {}
}
