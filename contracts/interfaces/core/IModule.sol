// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IModule
 * @notice Standard interface that all pluggable modules must implement
 * @dev This is the foundation of the pluggable architecture
 *
 * Design Principles:
 * - Every module must have a unique identifier
 * - Modules can be enabled/disabled without redeployment
 * - Modules declare their version for upgrade tracking
 * - Modules expose dependencies for validation
 */
interface IModule {
    /**
     * @notice Module lifecycle states
     */
    enum ModuleState {
        UNINITIALIZED,  // Module deployed but not initialized
        ACTIVE,         // Module is active and operational
        PAUSED,         // Module is paused (temporary)
        DEPRECATED,     // Module is deprecated (permanent)
        UPGRADED        // Module has been upgraded to a new version
    }

    /**
     * @notice Module metadata structure
     */
    struct ModuleInfo {
        bytes32 moduleId;           // Unique module identifier
        string name;                // Human-readable module name
        string version;             // Semantic version (e.g., "1.0.0")
        address implementation;     // Implementation contract address
        ModuleState state;          // Current module state
        uint256 installedAt;        // Installation timestamp
        uint256 lastUpdated;        // Last update timestamp
    }

    // ============ Events ============

    /**
     * @notice Emitted when module is initialized
     */
    event ModuleInitialized(bytes32 indexed moduleId, address indexed implementation, string version);

    /**
     * @notice Emitted when module state changes
     */
    event ModuleStateChanged(bytes32 indexed moduleId, ModuleState oldState, ModuleState newState);

    /**
     * @notice Emitted when module is upgraded
     */
    event ModuleUpgraded(bytes32 indexed moduleId, address oldImplementation, address newImplementation, string newVersion);

    // ============ Core Functions ============

    /**
     * @notice Get the unique identifier for this module
     * @return The module identifier (e.g., keccak256("VAULT_MODULE"))
     */
    function getModuleId() external pure returns (bytes32);

    /**
     * @notice Get comprehensive module information
     * @return Module metadata structure
     */
    function getModuleInfo() external view returns (ModuleInfo memory);

    /**
     * @notice Get module version
     * @return Semantic version string
     */
    function getVersion() external pure returns (string memory);

    /**
     * @notice Get module dependencies
     * @return Array of module IDs this module depends on
     * @dev Used by ModuleRegistry to validate dependency graph
     */
    function getDependencies() external pure returns (bytes32[] memory);

    /**
     * @notice Check if module is in active state
     * @return True if module is active and operational
     */
    function isActive() external view returns (bool);

    /**
     * @notice Initialize the module
     * @param data Initialization parameters (ABI encoded)
     * @dev Can only be called once, usually by ModuleRegistry
     */
    function initialize(bytes calldata data) external;

    /**
     * @notice Pause the module
     * @dev Should be called by authorized roles only
     */
    function pause() external;

    /**
     * @notice Resume the module
     * @dev Should be called by authorized roles only
     */
    function unpause() external;

    /**
     * @notice Perform health check on the module
     * @return healthy True if module is functioning correctly
     * @return message Status message or error description
     */
    function healthCheck() external view returns (bool healthy, string memory message);
}
