// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title L1StateRegistry
 * @notice Layer 1 state registry for Layer 2 state commitments
 * @dev Stores merkle roots from L2 and enables verification
 *
 * Architecture:
 * - Receives state roots from L2 via Arbitrum bridge
 * - Validates state commitments
 * - Triggers critical L1 operations (liquidations, emergency exits)
 * - Maintains historical state for fraud proofs
 */
contract L1StateRegistry is Ownable, ReentrancyGuard {
    // ============ State Variables ============

    // Mapping: L2 block number => state root
    mapping(uint256 => bytes32) public stateRoots;

    // Mapping: L2 block number => timestamp
    mapping(uint256 => uint256) public stateTimestamps;

    // Latest L2 block number submitted
    uint256 public latestL2Block;

    // Minimum time between state submissions (prevents spam)
    uint256 public constant MIN_SUBMISSION_INTERVAL = 1 hours;

    // Last submission timestamp
    uint256 public lastSubmissionTime;

    // Authorized L2 address (Arbitrum L2 contract)
    address public l2StateAggregator;

    // Emergency pause flag
    bool public emergencyPaused;

    // Critical state thresholds
    struct SystemThresholds {
        uint256 minCollateralRatio;      // 150 = 150%
        uint256 liquidationThreshold;    // 120 = 120%
        uint256 maxDebtCeiling;          // Maximum total debt allowed
    }

    SystemThresholds public thresholds;

    // State snapshot structure
    struct StateSnapshot {
        bytes32 stateRoot;
        uint256 l2BlockNumber;
        uint256 timestamp;
        uint256 totalCollateral;
        uint256 totalDebt;
        bool criticalCondition;
    }

    // Historical snapshots
    StateSnapshot[] public snapshots;

    // Emergency exit tracking
    mapping(address => bool) public emergencyExitClaimed;
    mapping(address => uint256) public emergencyExitAmount;

    // ============ Events ============

    event StateRootReceived(
        bytes32 indexed stateRoot,
        uint256 indexed l2Block,
        uint256 timestamp
    );

    event CriticalConditionDetected(
        bytes32 stateRoot,
        uint256 totalCollateral,
        uint256 totalDebt,
        uint256 ratio
    );

    event EmergencyPauseTriggered(address indexed triggeredBy, string reason);

    event EmergencyExitInitiated(
        address indexed user,
        uint256 amount,
        uint256 blockNumber
    );

    event L2AggregatorUpdated(address indexed oldAggregator, address indexed newAggregator);

    event ThresholdsUpdated(
        uint256 minCollateralRatio,
        uint256 liquidationThreshold,
        uint256 maxDebtCeiling
    );

    // ============ Modifiers ============

    modifier onlyL2Aggregator() {
        require(msg.sender == l2StateAggregator, "L1Registry: Only L2 aggregator");
        _;
    }

    modifier whenNotPaused() {
        require(!emergencyPaused, "L1Registry: Emergency paused");
        _;
    }

    // ============ Constructor ============

    constructor(address _l2StateAggregator) Ownable(msg.sender) {
        require(_l2StateAggregator != address(0), "Invalid L2 aggregator");

        l2StateAggregator = _l2StateAggregator;

        // Set default thresholds
        thresholds = SystemThresholds({
            minCollateralRatio: 150,
            liquidationThreshold: 120,
            maxDebtCeiling: 10_000_000 * 1e6  // 10M LUSD
        });

        lastSubmissionTime = block.timestamp;
    }

    // ============ Core Functions ============

    /**
     * @notice Receive state root from L2
     * @param stateRoot Merkle root of L2 state
     * @param l2Block L2 block number
     * @param totalCollateral Total collateral in system
     * @param totalDebt Total debt in system
     */
    function receiveStateRoot(
        bytes32 stateRoot,
        uint256 l2Block,
        uint256 totalCollateral,
        uint256 totalDebt
    ) external onlyL2Aggregator whenNotPaused nonReentrant {
        require(stateRoot != bytes32(0), "Invalid state root");
        require(l2Block > latestL2Block, "L2 block not newer");
        require(
            block.timestamp >= lastSubmissionTime + MIN_SUBMISSION_INTERVAL,
            "Submission too frequent"
        );

        // Store state root
        stateRoots[l2Block] = stateRoot;
        stateTimestamps[l2Block] = block.timestamp;
        latestL2Block = l2Block;
        lastSubmissionTime = block.timestamp;

        // Check for critical conditions
        bool critical = _checkCriticalConditions(totalCollateral, totalDebt);

        // Create snapshot
        snapshots.push(StateSnapshot({
            stateRoot: stateRoot,
            l2BlockNumber: l2Block,
            timestamp: block.timestamp,
            totalCollateral: totalCollateral,
            totalDebt: totalDebt,
            criticalCondition: critical
        }));

        emit StateRootReceived(stateRoot, l2Block, block.timestamp);
    }

    /**
     * @notice Verify a state commitment exists
     * @param l2Block L2 block number
     * @param expectedRoot Expected state root
     * @return bool True if state root matches
     */
    function verifyStateRoot(uint256 l2Block, bytes32 expectedRoot)
        external
        view
        returns (bool)
    {
        return stateRoots[l2Block] == expectedRoot && expectedRoot != bytes32(0);
    }

    /**
     * @notice Get the latest state root
     * @return stateRoot Latest state root
     * @return l2Block Latest L2 block number
     * @return timestamp Timestamp of submission
     */
    function getLatestState()
        external
        view
        returns (bytes32 stateRoot, uint256 l2Block, uint256 timestamp)
    {
        return (
            stateRoots[latestL2Block],
            latestL2Block,
            stateTimestamps[latestL2Block]
        );
    }

    /**
     * @notice Get state root for specific L2 block
     * @param l2Block L2 block number
     */
    function getStateRoot(uint256 l2Block) external view returns (bytes32) {
        return stateRoots[l2Block];
    }

    /**
     * @notice Get snapshot by index
     * @param index Snapshot index
     */
    function getSnapshot(uint256 index)
        external
        view
        returns (StateSnapshot memory)
    {
        require(index < snapshots.length, "Invalid index");
        return snapshots[index];
    }

    /**
     * @notice Get total number of snapshots
     */
    function getSnapshotCount() external view returns (uint256) {
        return snapshots.length;
    }

    // ============ Critical Condition Management ============

    /**
     * @notice Check if system meets critical conditions
     * @dev Internal function called when receiving state
     */
    function _checkCriticalConditions(
        uint256 totalCollateral,
        uint256 totalDebt
    ) internal returns (bool) {
        if (totalDebt == 0) return false;

        // Calculate system collateral ratio
        uint256 systemRatio = (totalCollateral * 100) / totalDebt;

        bool critical = false;

        // Check if system is under-collateralized
        if (systemRatio < thresholds.minCollateralRatio) {
            critical = true;
            emit CriticalConditionDetected(
                stateRoots[latestL2Block],
                totalCollateral,
                totalDebt,
                systemRatio
            );
        }

        // Check if debt ceiling exceeded
        if (totalDebt > thresholds.maxDebtCeiling) {
            critical = true;
        }

        return critical;
    }

    // ============ Emergency Functions ============

    /**
     * @notice Trigger emergency pause
     * @param reason Reason for emergency pause
     */
    function triggerEmergencyPause(string calldata reason) external onlyOwner {
        emergencyPaused = true;
        emit EmergencyPauseTriggered(msg.sender, reason);
    }

    /**
     * @notice Resume from emergency pause
     */
    function resumeFromEmergency() external onlyOwner {
        emergencyPaused = false;
    }

    /**
     * @notice Initiate emergency exit for user
     * @param user User address
     * @param amount Amount to exit
     * @dev Only callable by owner in emergency situations
     */
    function initiateEmergencyExit(address user, uint256 amount)
        external
        onlyOwner
    {
        require(emergencyPaused, "Not in emergency");
        require(!emergencyExitClaimed[user], "Already claimed");

        emergencyExitAmount[user] = amount;

        emit EmergencyExitInitiated(user, amount, latestL2Block);
    }

    // ============ Admin Functions ============

    /**
     * @notice Update L2 aggregator address
     * @param newAggregator New aggregator address
     */
    function setL2Aggregator(address newAggregator) external onlyOwner {
        require(newAggregator != address(0), "Invalid address");
        address oldAggregator = l2StateAggregator;
        l2StateAggregator = newAggregator;
        emit L2AggregatorUpdated(oldAggregator, newAggregator);
    }

    /**
     * @notice Update system thresholds
     */
    function updateThresholds(
        uint256 _minCollateralRatio,
        uint256 _liquidationThreshold,
        uint256 _maxDebtCeiling
    ) external onlyOwner {
        require(_minCollateralRatio > _liquidationThreshold, "Invalid ratios");
        require(_liquidationThreshold > 100, "Threshold too low");

        thresholds.minCollateralRatio = _minCollateralRatio;
        thresholds.liquidationThreshold = _liquidationThreshold;
        thresholds.maxDebtCeiling = _maxDebtCeiling;

        emit ThresholdsUpdated(
            _minCollateralRatio,
            _liquidationThreshold,
            _maxDebtCeiling
        );
    }

    // ============ View Functions ============

    /**
     * @notice Check if state is fresh (submitted within last 2 hours)
     */
    function isStateFresh() external view returns (bool) {
        if (latestL2Block == 0) return false;
        return block.timestamp - stateTimestamps[latestL2Block] <= 2 hours;
    }

    /**
     * @notice Get time since last submission
     */
    function timeSinceLastSubmission() external view returns (uint256) {
        return block.timestamp - lastSubmissionTime;
    }

    /**
     * @notice Check if system can accept new submission
     */
    function canSubmit() external view returns (bool) {
        return block.timestamp >= lastSubmissionTime + MIN_SUBMISSION_INTERVAL;
    }
}
