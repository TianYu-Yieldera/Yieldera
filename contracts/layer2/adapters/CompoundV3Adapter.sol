// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "../interfaces/external/ICompoundV3.sol";
import "../aggregator/L2StateAggregator.sol";

/**
 * @title CompoundV3Adapter
 * @notice Adapter for Compound V3 (Comet) protocol
 * @dev Integrates with Compound V3 for lending and borrowing
 *
 * Architecture:
 * - Simplified single-market design (one base asset)
 * - Multiple collateral assets supported
 * - Direct interaction with Comet contract
 * - State aggregation for L1 reporting
 *
 * Key Features:
 * - Supply base asset (e.g., USDC) to earn interest
 * - Supply collateral to enable borrowing
 * - Borrow base asset against collateral
 * - Repay borrowed base asset
 * - Withdraw supplied assets
 * - COMP rewards tracking (optional)
 *
 * Security:
 * - ReentrancyGuard on all external functions
 * - Pausable for emergency stops
 * - AccessControl for admin functions
 */
contract CompoundV3Adapter is AccessControl, ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;

    bytes32 public constant MANAGER_ROLE = keccak256("MANAGER_ROLE");

    // Compound V3 contracts
    IComet public immutable comet;
    ICometRewards public immutable cometRewards;
    address public immutable baseAsset;

    // State aggregator integration
    L2StateAggregator public stateAggregator;

    // Statistics tracking
    uint256 public totalSupplied;      // Total base asset supplied
    uint256 public totalBorrowed;      // Total base asset borrowed
    uint256 public totalCollateral;    // Total collateral value (USD)
    uint256 public activePositions;    // Number of active positions

    // User position tracking
    struct Position {
        uint256 baseSupplied;    // Base asset supplied
        uint256 baseBorrowed;    // Base asset borrowed
        uint256 lastUpdate;      // Last interaction timestamp
        bool active;             // Position active flag
    }

    mapping(address => Position) public positions;

    // Events
    event BaseSupplied(address indexed user, uint256 amount);
    event BaseWithdrawn(address indexed user, uint256 amount);
    event CollateralSupplied(address indexed user, address indexed asset, uint256 amount);
    event CollateralWithdrawn(address indexed user, address indexed asset, uint256 amount);
    event Borrowed(address indexed user, uint256 amount);
    event Repaid(address indexed user, uint256 amount);
    event RewardsClaimed(address indexed user, uint256 amount);

    /**
     * @notice Constructor
     * @param _comet Compound V3 Comet contract
     * @param _cometRewards Compound V3 rewards contract
     * @param _admin Admin address
     */
    constructor(
        address _comet,
        address _cometRewards,
        address _admin
    ) {
        require(_comet != address(0), "Invalid comet");
        require(_admin != address(0), "Invalid admin");

        comet = IComet(_comet);
        cometRewards = ICometRewards(_cometRewards);
        baseAsset = comet.baseToken();

        _grantRole(DEFAULT_ADMIN_ROLE, _admin);
        _grantRole(MANAGER_ROLE, _admin);
    }

    // =============================================================
    //                      SUPPLY FUNCTIONS
    // =============================================================

    /**
     * @notice Supply base asset to earn interest
     * @param amount Amount of base asset to supply
     */
    function supplyBase(uint256 amount) external nonReentrant whenNotPaused {
        require(amount > 0, "Amount must be greater than zero");

        // Transfer base asset from user
        IERC20(baseAsset).safeTransferFrom(msg.sender, address(this), amount);

        // Approve Comet
        IERC20(baseAsset).forceApprove(address(comet), amount);

        // Supply to Compound
        comet.supply(baseAsset, amount);

        // Update position
        Position storage pos = positions[msg.sender];
        if (!pos.active) {
            pos.active = true;
            activePositions++;
        }
        pos.baseSupplied += amount;
        pos.lastUpdate = block.timestamp;

        // Update global stats
        totalSupplied += amount;
        _updateAggregator();

        emit BaseSupplied(msg.sender, amount);
    }

    /**
     * @notice Supply collateral to enable borrowing
     * @param asset Collateral asset address
     * @param amount Amount to supply
     */
    function supplyCollateral(address asset, uint256 amount)
        external
        nonReentrant
        whenNotPaused
    {
        require(amount > 0, "Amount must be greater than zero");
        require(asset != baseAsset, "Use supplyBase for base asset");

        // Transfer collateral from user
        IERC20(asset).safeTransferFrom(msg.sender, address(this), amount);

        // Approve Comet
        IERC20(asset).forceApprove(address(comet), amount);

        // Supply to Compound
        comet.supply(asset, amount);

        // Update position
        Position storage pos = positions[msg.sender];
        if (!pos.active) {
            pos.active = true;
            activePositions++;
        }
        pos.lastUpdate = block.timestamp;

        // Update global stats (simplified - not tracking individual collateral)
        _updateAggregator();

        emit CollateralSupplied(msg.sender, asset, amount);
    }

    // =============================================================
    //                    WITHDRAW FUNCTIONS
    // =============================================================

    /**
     * @notice Withdraw supplied base asset
     * @param amount Amount to withdraw
     */
    function withdrawBase(uint256 amount) external nonReentrant whenNotPaused {
        require(amount > 0, "Amount must be greater than zero");

        Position storage pos = positions[msg.sender];
        require(pos.baseSupplied >= amount, "Insufficient supplied balance");

        // Withdraw from Compound to user
        comet.withdrawTo(msg.sender, baseAsset, amount);

        // Update position
        pos.baseSupplied -= amount;
        pos.lastUpdate = block.timestamp;

        if (pos.baseSupplied == 0 && pos.baseBorrowed == 0) {
            pos.active = false;
            activePositions--;
        }

        // Update global stats
        totalSupplied -= amount;
        _updateAggregator();

        emit BaseWithdrawn(msg.sender, amount);
    }

    /**
     * @notice Withdraw collateral
     * @param asset Collateral asset address
     * @param amount Amount to withdraw
     */
    function withdrawCollateral(address asset, uint256 amount)
        external
        nonReentrant
        whenNotPaused
    {
        require(amount > 0, "Amount must be greater than zero");
        require(asset != baseAsset, "Use withdrawBase for base asset");

        // Withdraw from Compound to user
        comet.withdrawTo(msg.sender, asset, amount);

        // Update position timestamp
        Position storage pos = positions[msg.sender];
        pos.lastUpdate = block.timestamp;

        _updateAggregator();

        emit CollateralWithdrawn(msg.sender, asset, amount);
    }

    // =============================================================
    //                     BORROW FUNCTIONS
    // =============================================================

    /**
     * @notice Borrow base asset against collateral
     * @param amount Amount to borrow
     */
    function borrow(uint256 amount) external nonReentrant whenNotPaused {
        require(amount > 0, "Amount must be greater than zero");

        // Withdraw acts as borrow in Compound V3
        comet.withdrawTo(msg.sender, baseAsset, amount);

        // Update position
        Position storage pos = positions[msg.sender];
        if (!pos.active) {
            pos.active = true;
            activePositions++;
        }
        pos.baseBorrowed += amount;
        pos.lastUpdate = block.timestamp;

        // Update global stats
        totalBorrowed += amount;
        _updateAggregator();

        emit Borrowed(msg.sender, amount);
    }

    /**
     * @notice Repay borrowed base asset
     * @param amount Amount to repay
     */
    function repay(uint256 amount) external nonReentrant whenNotPaused {
        require(amount > 0, "Amount must be greater than zero");

        Position storage pos = positions[msg.sender];
        require(pos.baseBorrowed > 0, "No debt to repay");

        uint256 repayAmount = amount > pos.baseBorrowed ? pos.baseBorrowed : amount;

        // Transfer base asset from user
        IERC20(baseAsset).safeTransferFrom(msg.sender, address(this), repayAmount);

        // Approve Comet
        IERC20(baseAsset).forceApprove(address(comet), repayAmount);

        // Supply to repay debt
        comet.supply(baseAsset, repayAmount);

        // Update position
        pos.baseBorrowed -= repayAmount;
        pos.lastUpdate = block.timestamp;

        if (pos.baseSupplied == 0 && pos.baseBorrowed == 0) {
            pos.active = false;
            activePositions--;
        }

        // Update global stats
        totalBorrowed -= repayAmount;
        _updateAggregator();

        emit Repaid(msg.sender, repayAmount);
    }

    // =============================================================
    //                      REWARDS FUNCTIONS
    // =============================================================

    /**
     * @notice Claim COMP rewards
     * @param shouldAccrue Whether to accrue before claiming
     */
    function claimRewards(bool shouldAccrue) external nonReentrant {
        if (address(cometRewards) == address(0)) {
            revert("Rewards not configured");
        }

        cometRewards.claim(address(comet), msg.sender, shouldAccrue);

        // Note: We don't track exact amount here as it's emitted by CometRewards
        emit RewardsClaimed(msg.sender, 0);
    }

    // =============================================================
    //                       VIEW FUNCTIONS
    // =============================================================

    /**
     * @notice Get user position details
     * @param user User address
     * @return baseSupplied Amount of base asset supplied
     * @return baseBorrowed Amount of base asset borrowed
     * @return collateralValue Total collateral value (not implemented in MVP)
     * @return lastUpdate Last update timestamp
     */
    function getPosition(address user)
        external
        view
        returns (
            uint256 baseSupplied,
            uint256 baseBorrowed,
            uint256 collateralValue,
            uint256 lastUpdate
        )
    {
        Position memory pos = positions[user];
        return (
            pos.baseSupplied,
            pos.baseBorrowed,
            0, // Collateral value calculation not implemented in MVP
            pos.lastUpdate
        );
    }

    /**
     * @notice Get current supply and borrow rates
     * @return supplyRate Current supply APR
     * @return borrowRate Current borrow APR
     */
    function getRates() external view returns (uint256 supplyRate, uint256 borrowRate) {
        uint256 utilization = comet.getUtilization();
        supplyRate = comet.getSupplyRate(utilization);
        borrowRate = comet.getBorrowRate(utilization);
    }

    /**
     * @notice Check if user position is healthy
     * @param user User address
     * @return isHealthy True if position is healthy
     */
    function isPositionHealthy(address user) external view returns (bool) {
        return !comet.isLiquidatable(user);
    }

    // =============================================================
    //                    INTERNAL FUNCTIONS
    // =============================================================

    /**
     * @notice Update state aggregator
     */
    function _updateAggregator() internal {
        if (address(stateAggregator) != address(0)) {
            stateAggregator.updateSystemState(
                totalCollateral,
                totalBorrowed,
                activePositions,
                totalSupplied // Using totalSupplied as totalOrders
            );
        }
    }

    // =============================================================
    //                      ADMIN FUNCTIONS
    // =============================================================

    /**
     * @notice Set state aggregator
     * @param _stateAggregator State aggregator address
     */
    function setStateAggregator(address _stateAggregator)
        external
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        stateAggregator = L2StateAggregator(_stateAggregator);
    }

    /**
     * @notice Pause contract
     */
    function pause() external onlyRole(MANAGER_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause contract
     */
    function unpause() external onlyRole(MANAGER_ROLE) {
        _unpause();
    }

    /**
     * @notice Emergency token withdrawal
     * @param token Token address
     * @param amount Amount to withdraw
     */
    function emergencyWithdraw(address token, uint256 amount)
        external
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        IERC20(token).safeTransfer(msg.sender, amount);
    }

    /**
     * @notice Receive ETH
     */
    receive() external payable {}
}
