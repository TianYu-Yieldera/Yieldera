// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title CollateralVaultL1
 * @notice Layer 1 optimized version - Only asset custody
 * @dev Simplified vault for L1, computation moved to L2
 *
 * L1 Responsibilities:
 * - Lock/unlock collateral tokens
 * - Final custody of assets
 * - Emergency withdrawals
 * - Bridge integration
 *
 * L2 handles:
 * - Collateral ratio calculations
 * - Debt tracking
 * - Interest calculations
 * - Position management
 * - Liquidation logic
 */
contract CollateralVaultL1 is ReentrancyGuard, Ownable {
    using SafeERC20 for IERC20;

    // ============ State Variables ============

    // Loyalty Points token address
    IERC20 public immutable collateralToken;

    // User collateral balances (locked on L1)
    mapping(address => uint256) public lockedCollateral;

    // Total collateral locked in vault
    uint256 public totalLocked;

    // L2 Bridge address (authorized to lock/unlock)
    address public l2Bridge;

    // L1 State Registry (for emergency operations)
    address public stateRegistry;

    // Emergency pause flag
    bool public emergencyPaused;

    // Emergency withdrawal tracking (7-day delay)
    struct EmergencyWithdrawal {
        uint256 amount;
        uint256 requestTime;
        bool executed;
    }

    mapping(address => EmergencyWithdrawal) public emergencyWithdrawals;

    uint256 public constant EMERGENCY_DELAY = 7 days;

    // Rate limiting for security
    uint256 public constant MAX_LOCK_PER_TX = 1_000_000 * 1e18; // 1M tokens
    uint256 public dailyLockLimit = 10_000_000 * 1e18; // 10M tokens per day
    uint256 public dailyLockedAmount;
    uint256 public lastLockResetTime;

    // ============ Events ============

    event CollateralLocked(
        address indexed user,
        uint256 amount,
        uint256 totalUserLocked,
        bytes32 indexed l2TxHash
    );

    event CollateralUnlocked(
        address indexed user,
        uint256 amount,
        uint256 remaining,
        bytes32 indexed l2TxHash
    );

    event EmergencyWithdrawalRequested(
        address indexed user,
        uint256 amount,
        uint256 unlockTime
    );

    event EmergencyWithdrawalExecuted(
        address indexed user,
        uint256 amount
    );

    event L2BridgeUpdated(address indexed oldBridge, address indexed newBridge);
    event StateRegistryUpdated(address indexed oldRegistry, address indexed newRegistry);
    event EmergencyPauseTriggered(address indexed triggeredBy);
    event EmergencyResumed(address indexed resumedBy);

    // ============ Modifiers ============

    modifier onlyBridge() {
        require(msg.sender == l2Bridge, "Only L2 bridge");
        _;
    }

    modifier whenNotPaused() {
        require(!emergencyPaused, "Emergency paused");
        _;
    }

    // ============ Constructor ============

    constructor(address _collateralToken) Ownable(msg.sender) {
        require(_collateralToken != address(0), "Invalid token address");
        collateralToken = IERC20(_collateralToken);
        lastLockResetTime = block.timestamp;
    }

    // ============ Core Functions ============

    /**
     * @notice Lock collateral from user (called by bridge)
     * @param user User address
     * @param amount Amount to lock
     * @param l2TxHash L2 transaction hash for tracking
     */
    function lockCollateral(address user, uint256 amount, bytes32 l2TxHash)
        external
        onlyBridge
        whenNotPaused
        nonReentrant
    {
        require(user != address(0), "Invalid user");
        require(amount > 0, "Amount must be > 0");
        require(amount <= MAX_LOCK_PER_TX, "Exceeds max per tx");

        // Check daily limit
        _checkAndUpdateDailyLimit(amount);

        // Transfer tokens from user to vault
        collateralToken.safeTransferFrom(user, address(this), amount);

        // Update balances
        lockedCollateral[user] += amount;
        totalLocked += amount;

        emit CollateralLocked(user, amount, lockedCollateral[user], l2TxHash);
    }

    /**
     * @notice Unlock collateral to user (called by bridge)
     * @param user User address
     * @param amount Amount to unlock
     * @param l2TxHash L2 transaction hash for tracking
     */
    function unlockCollateral(address user, uint256 amount, bytes32 l2TxHash)
        external
        onlyBridge
        whenNotPaused
        nonReentrant
    {
        require(user != address(0), "Invalid user");
        require(amount > 0, "Amount must be > 0");
        require(lockedCollateral[user] >= amount, "Insufficient locked collateral");

        // Update balances
        lockedCollateral[user] -= amount;
        totalLocked -= amount;

        // Transfer tokens back to user
        collateralToken.safeTransfer(user, amount);

        emit CollateralUnlocked(user, amount, lockedCollateral[user], l2TxHash);
    }

    /**
     * @notice Check and update daily lock limit
     * @dev Internal function to prevent over-locking
     */
    function _checkAndUpdateDailyLimit(uint256 amount) internal {
        // Reset daily counter if 24 hours passed
        if (block.timestamp >= lastLockResetTime + 1 days) {
            dailyLockedAmount = 0;
            lastLockResetTime = block.timestamp;
        }

        require(
            dailyLockedAmount + amount <= dailyLockLimit,
            "Exceeds daily lock limit"
        );

        dailyLockedAmount += amount;
    }

    // ============ Emergency Functions ============

    /**
     * @notice Request emergency withdrawal (7-day delay)
     * @param amount Amount to withdraw
     */
    function requestEmergencyWithdrawal(uint256 amount) external nonReentrant {
        require(emergencyPaused, "Not in emergency mode");
        require(lockedCollateral[msg.sender] >= amount, "Insufficient collateral");
        require(!emergencyWithdrawals[msg.sender].executed, "Already requested");
        require(emergencyWithdrawals[msg.sender].requestTime == 0, "Pending request exists");

        emergencyWithdrawals[msg.sender] = EmergencyWithdrawal({
            amount: amount,
            requestTime: block.timestamp,
            executed: false
        });

        emit EmergencyWithdrawalRequested(
            msg.sender,
            amount,
            block.timestamp + EMERGENCY_DELAY
        );
    }

    /**
     * @notice Execute emergency withdrawal after delay
     */
    function executeEmergencyWithdrawal() external nonReentrant {
        EmergencyWithdrawal storage withdrawal = emergencyWithdrawals[msg.sender];

        require(withdrawal.requestTime > 0, "No withdrawal request");
        require(!withdrawal.executed, "Already executed");
        require(
            block.timestamp >= withdrawal.requestTime + EMERGENCY_DELAY,
            "Delay period not passed"
        );
        require(lockedCollateral[msg.sender] >= withdrawal.amount, "Insufficient collateral");

        // Mark as executed
        withdrawal.executed = true;

        // Update balances
        lockedCollateral[msg.sender] -= withdrawal.amount;
        totalLocked -= withdrawal.amount;

        // Transfer tokens
        collateralToken.safeTransfer(msg.sender, withdrawal.amount);

        emit EmergencyWithdrawalExecuted(msg.sender, withdrawal.amount);
    }

    /**
     * @notice Trigger emergency pause
     */
    function triggerEmergencyPause() external onlyOwner {
        emergencyPaused = true;
        emit EmergencyPauseTriggered(msg.sender);
    }

    /**
     * @notice Resume from emergency
     */
    function resumeFromEmergency() external onlyOwner {
        emergencyPaused = false;
        emit EmergencyResumed(msg.sender);
    }

    // ============ Admin Functions ============

    /**
     * @notice Set L2 bridge address
     * @param _l2Bridge L2 bridge address
     */
    function setL2Bridge(address _l2Bridge) external onlyOwner {
        require(_l2Bridge != address(0), "Invalid bridge address");
        address oldBridge = l2Bridge;
        l2Bridge = _l2Bridge;
        emit L2BridgeUpdated(oldBridge, _l2Bridge);
    }

    /**
     * @notice Set state registry address
     * @param _stateRegistry State registry address
     */
    function setStateRegistry(address _stateRegistry) external onlyOwner {
        require(_stateRegistry != address(0), "Invalid registry address");
        address oldRegistry = stateRegistry;
        stateRegistry = _stateRegistry;
        emit StateRegistryUpdated(oldRegistry, _stateRegistry);
    }

    /**
     * @notice Update daily lock limit
     * @param newLimit New daily limit
     */
    function setDailyLockLimit(uint256 newLimit) external onlyOwner {
        dailyLockLimit = newLimit;
    }

    // ============ View Functions ============

    /**
     * @notice Get user's locked collateral
     * @param user User address
     */
    function getLockedCollateral(address user) external view returns (uint256) {
        return lockedCollateral[user];
    }

    /**
     * @notice Get total locked collateral
     */
    function getTotalLocked() external view returns (uint256) {
        return totalLocked;
    }

    /**
     * @notice Get vault statistics
     */
    function getVaultStats()
        external
        view
        returns (uint256 _totalLocked, uint256 contractBalance)
    {
        return (totalLocked, collateralToken.balanceOf(address(this)));
    }

    /**
     * @notice Get emergency withdrawal info
     * @param user User address
     */
    function getEmergencyWithdrawal(address user)
        external
        view
        returns (uint256 amount, uint256 requestTime, bool executed, uint256 unlockTime)
    {
        EmergencyWithdrawal memory withdrawal = emergencyWithdrawals[user];
        return (
            withdrawal.amount,
            withdrawal.requestTime,
            withdrawal.executed,
            withdrawal.requestTime > 0 ? withdrawal.requestTime + EMERGENCY_DELAY : 0
        );
    }

    /**
     * @notice Check if emergency withdrawal can be executed
     * @param user User address
     */
    function canExecuteEmergencyWithdrawal(address user) external view returns (bool) {
        EmergencyWithdrawal memory withdrawal = emergencyWithdrawals[user];
        return
            withdrawal.requestTime > 0 &&
            !withdrawal.executed &&
            block.timestamp >= withdrawal.requestTime + EMERGENCY_DELAY;
    }

    /**
     * @notice Get remaining daily lock capacity
     */
    function getRemainingDailyLockCapacity() external view returns (uint256) {
        // Check if we need to reset
        if (block.timestamp >= lastLockResetTime + 1 days) {
            return dailyLockLimit;
        }

        if (dailyLockedAmount >= dailyLockLimit) {
            return 0;
        }

        return dailyLockLimit - dailyLockedAmount;
    }

    /**
     * @notice Get time until daily limit resets
     */
    function getTimeUntilLimitReset() external view returns (uint256) {
        uint256 resetTime = lastLockResetTime + 1 days;
        if (block.timestamp >= resetTime) {
            return 0;
        }
        return resetTime - block.timestamp;
    }
}
