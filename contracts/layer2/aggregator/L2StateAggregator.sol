// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";

/**
 * @title L2StateAggregator
 * @notice Layer 2 state aggregator for periodic L1 submissions
 * @dev Aggregates L2 state and submits merkle roots to L1
 *
 * Architecture:
 * - Collects state from all L2 modules
 * - Computes merkle root of system state
 * - Submits to L1 StateRegistry every hour
 * - Uses Arbitrum's native ArbSys for L1 messaging
 */

// Arbitrum ArbSys interface
interface IArbSys {
    function sendTxToL1(address destination, bytes calldata data)
        external
        payable
        returns (uint256);
}

contract L2StateAggregator is Ownable {
    // ============ Constants ============

    // Arbitrum ArbSys precompile address
    address constant ARBSYS_ADDR = address(100);

    // Minimum time between submissions (1 hour)
    uint256 public constant MIN_SUBMISSION_INTERVAL = 1 hours;

    // ============ State Variables ============

    // Current state root
    bytes32 public currentStateRoot;

    // Last submission timestamp
    uint256 public lastSubmission;

    // L1 State Registry address
    address public l1StateRegistry;

    // Module state tracking
    struct ModuleState {
        bytes32 stateHash;
        uint256 lastUpdate;
        bool active;
    }

    // Module registry
    mapping(bytes32 => ModuleState) public moduleStates;
    bytes32[] public moduleIds;

    // Aggregated system state
    struct SystemState {
        uint256 totalCollateral;
        uint256 totalDebt;
        uint256 activePositions;
        uint256 totalOrders;
        uint256 blockNumber;
        uint256 timestamp;
    }

    SystemState public systemState;

    // Historical state roots
    bytes32[] public stateHistory;
    mapping(bytes32 => uint256) public stateRootToBlock;

    // Authorized modules (can update state)
    mapping(address => bool) public authorizedModules;

    // ============ Events ============

    event StateRootComputed(
        bytes32 indexed stateRoot,
        uint256 blockNumber,
        uint256 timestamp
    );

    event StateSubmittedToL1(
        bytes32 indexed stateRoot,
        uint256 indexed l2Block,
        uint256 indexed l1TxId,
        uint256 totalCollateral,
        uint256 totalDebt
    );

    event ModuleStateUpdated(
        bytes32 indexed moduleId,
        bytes32 stateHash,
        address updatedBy
    );

    event ModuleRegistered(bytes32 indexed moduleId, address moduleAddress);
    event ModuleDeactivated(bytes32 indexed moduleId);
    event L1RegistryUpdated(address indexed oldRegistry, address indexed newRegistry);

    // ============ Modifiers ============

    modifier onlyAuthorizedModule() {
        require(authorizedModules[msg.sender], "Not authorized module");
        _;
    }

    modifier canSubmit() {
        require(
            block.timestamp >= lastSubmission + MIN_SUBMISSION_INTERVAL,
            "Submission too frequent"
        );
        _;
    }

    // ============ Constructor ============

    constructor(address _l1StateRegistry) Ownable(msg.sender) {
        require(_l1StateRegistry != address(0), "Invalid L1 registry");
        l1StateRegistry = _l1StateRegistry;
        lastSubmission = block.timestamp;
    }

    // ============ Core Functions ============

    /**
     * @notice Update module state
     * @param moduleId Module identifier
     * @param stateHash Hash of module state
     */
    function updateModuleState(bytes32 moduleId, bytes32 stateHash)
        external
        onlyAuthorizedModule
    {
        require(moduleId != bytes32(0), "Invalid module ID");
        require(stateHash != bytes32(0), "Invalid state hash");

        // Update or create module state
        if (!moduleStates[moduleId].active) {
            moduleIds.push(moduleId);
        }

        moduleStates[moduleId] = ModuleState({
            stateHash: stateHash,
            lastUpdate: block.timestamp,
            active: true
        });

        emit ModuleStateUpdated(moduleId, stateHash, msg.sender);
    }

    /**
     * @notice Update system-level state
     * @param totalCollateral Total collateral in system
     * @param totalDebt Total debt in system
     * @param activePositions Number of active positions
     * @param totalOrders Number of active orders
     */
    function updateSystemState(
        uint256 totalCollateral,
        uint256 totalDebt,
        uint256 activePositions,
        uint256 totalOrders
    ) external onlyAuthorizedModule {
        systemState = SystemState({
            totalCollateral: totalCollateral,
            totalDebt: totalDebt,
            activePositions: activePositions,
            totalOrders: totalOrders,
            blockNumber: block.number,
            timestamp: block.timestamp
        });
    }

    /**
     * @notice Calculate current state root
     * @return State root (merkle root of all module states)
     */
    function calculateStateRoot() public view returns (bytes32) {
        // Create array of module state hashes
        bytes32[] memory hashes = new bytes32[](moduleIds.length + 1);

        // Add module states
        for (uint256 i = 0; i < moduleIds.length; i++) {
            bytes32 moduleId = moduleIds[i];
            if (moduleStates[moduleId].active) {
                hashes[i] = keccak256(
                    abi.encodePacked(moduleId, moduleStates[moduleId].stateHash)
                );
            }
        }

        // Add system state
        hashes[moduleIds.length] = keccak256(
            abi.encodePacked(
                systemState.totalCollateral,
                systemState.totalDebt,
                systemState.activePositions,
                systemState.totalOrders,
                systemState.blockNumber,
                systemState.timestamp
            )
        );

        // Compute merkle root
        return _computeMerkleRoot(hashes);
    }

    /**
     * @notice Submit state to L1
     * @dev Can be called by anyone after MIN_SUBMISSION_INTERVAL
     */
    function submitToL1() external canSubmit {
        // Calculate state root
        bytes32 stateRoot = calculateStateRoot();
        currentStateRoot = stateRoot;

        // Record in history
        stateHistory.push(stateRoot);
        stateRootToBlock[stateRoot] = block.number;

        // Prepare L1 message
        bytes memory data = abi.encodeWithSignature(
            "receiveStateRoot(bytes32,uint256,uint256,uint256)",
            stateRoot,
            block.number,
            systemState.totalCollateral,
            systemState.totalDebt
        );

        // Send to L1 via ArbSys
        uint256 l1TxId = IArbSys(ARBSYS_ADDR).sendTxToL1(
            l1StateRegistry,
            data
        );

        // Update last submission time
        lastSubmission = block.timestamp;

        emit StateRootComputed(stateRoot, block.number, block.timestamp);
        emit StateSubmittedToL1(
            stateRoot,
            block.number,
            l1TxId,
            systemState.totalCollateral,
            systemState.totalDebt
        );
    }

    /**
     * @notice Force submit to L1 (owner only)
     * @dev Bypasses time restriction for emergencies
     */
    function forceSubmitToL1() external onlyOwner {
        bytes32 stateRoot = calculateStateRoot();
        currentStateRoot = stateRoot;

        stateHistory.push(stateRoot);
        stateRootToBlock[stateRoot] = block.number;

        bytes memory data = abi.encodeWithSignature(
            "receiveStateRoot(bytes32,uint256,uint256,uint256)",
            stateRoot,
            block.number,
            systemState.totalCollateral,
            systemState.totalDebt
        );

        uint256 l1TxId = IArbSys(ARBSYS_ADDR).sendTxToL1(
            l1StateRegistry,
            data
        );

        lastSubmission = block.timestamp;

        emit StateRootComputed(stateRoot, block.number, block.timestamp);
        emit StateSubmittedToL1(
            stateRoot,
            block.number,
            l1TxId,
            systemState.totalCollateral,
            systemState.totalDebt
        );
    }

    /**
     * @notice Compute merkle root from array of hashes
     * @dev Internal helper function
     */
    function _computeMerkleRoot(bytes32[] memory hashes)
        internal
        pure
        returns (bytes32)
    {
        uint256 n = hashes.length;
        if (n == 0) return bytes32(0);
        if (n == 1) return hashes[0];

        // Build merkle tree bottom-up
        while (n > 1) {
            uint256 j = 0;
            for (uint256 i = 0; i < n - 1; i += 2) {
                hashes[j] = keccak256(abi.encodePacked(hashes[i], hashes[i + 1]));
                j++;
            }
            if (n % 2 == 1) {
                hashes[j] = hashes[n - 1];
                j++;
            }
            n = j;
        }

        return hashes[0];
    }

    // ============ Module Management ============

    /**
     * @notice Register a new module
     * @param moduleId Module identifier
     * @param moduleAddress Module contract address
     */
    function registerModule(bytes32 moduleId, address moduleAddress)
        external
        onlyOwner
    {
        require(moduleId != bytes32(0), "Invalid module ID");
        require(moduleAddress != address(0), "Invalid module address");

        authorizedModules[moduleAddress] = true;

        emit ModuleRegistered(moduleId, moduleAddress);
    }

    /**
     * @notice Deactivate a module
     * @param moduleId Module identifier
     */
    function deactivateModule(bytes32 moduleId) external onlyOwner {
        require(moduleStates[moduleId].active, "Module not active");
        moduleStates[moduleId].active = false;
        emit ModuleDeactivated(moduleId);
    }

    /**
     * @notice Authorize a module address
     * @param moduleAddress Module contract address
     */
    function authorizeModule(address moduleAddress) external onlyOwner {
        require(moduleAddress != address(0), "Invalid address");
        authorizedModules[moduleAddress] = true;
    }

    /**
     * @notice Revoke module authorization
     * @param moduleAddress Module contract address
     */
    function revokeModuleAuthorization(address moduleAddress) external onlyOwner {
        authorizedModules[moduleAddress] = false;
    }

    // ============ Admin Functions ============

    /**
     * @notice Update L1 registry address
     * @param _l1StateRegistry New L1 registry address
     */
    function setL1StateRegistry(address _l1StateRegistry) external onlyOwner {
        require(_l1StateRegistry != address(0), "Invalid address");
        address oldRegistry = l1StateRegistry;
        l1StateRegistry = _l1StateRegistry;
        emit L1RegistryUpdated(oldRegistry, _l1StateRegistry);
    }

    // ============ View Functions ============

    /**
     * @notice Get module state
     * @param moduleId Module identifier
     */
    function getModuleState(bytes32 moduleId)
        external
        view
        returns (bytes32 stateHash, uint256 lastUpdate, bool active)
    {
        ModuleState memory state = moduleStates[moduleId];
        return (state.stateHash, state.lastUpdate, state.active);
    }

    /**
     * @notice Get number of registered modules
     */
    function getModuleCount() external view returns (uint256) {
        return moduleIds.length;
    }

    /**
     * @notice Get all module IDs
     */
    function getAllModuleIds() external view returns (bytes32[] memory) {
        return moduleIds;
    }

    /**
     * @notice Get state history length
     */
    function getStateHistoryLength() external view returns (uint256) {
        return stateHistory.length;
    }

    /**
     * @notice Get state root by index
     */
    function getStateRootByIndex(uint256 index) external view returns (bytes32) {
        require(index < stateHistory.length, "Invalid index");
        return stateHistory[index];
    }

    /**
     * @notice Check if can submit to L1
     */
    function canSubmitToL1() external view returns (bool) {
        return block.timestamp >= lastSubmission + MIN_SUBMISSION_INTERVAL;
    }

    /**
     * @notice Get time until next submission
     */
    function timeUntilNextSubmission() external view returns (uint256) {
        uint256 nextSubmissionTime = lastSubmission + MIN_SUBMISSION_INTERVAL;
        if (block.timestamp >= nextSubmissionTime) {
            return 0;
        }
        return nextSubmissionTime - block.timestamp;
    }

    /**
     * @notice Get current system state
     */
    function getSystemState()
        external
        view
        returns (
            uint256 totalCollateral,
            uint256 totalDebt,
            uint256 activePositions,
            uint256 totalOrders,
            uint256 blockNumber,
            uint256 timestamp
        )
    {
        return (
            systemState.totalCollateral,
            systemState.totalDebt,
            systemState.activePositions,
            systemState.totalOrders,
            systemState.blockNumber,
            systemState.timestamp
        );
    }
}
