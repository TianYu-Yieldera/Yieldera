// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IPlugin
 * @notice Core interface that all plugins must implement
 * @dev Provides standardized plugin functionality and metadata
 */
interface IPlugin {
    // ============ Enums ============

    enum PluginType {
        VAULT,
        RWA,
        ORACLE,
        AUDIT,
        REWARDS,
        ANALYTICS,
        INTEGRATION,
        CUSTOM
    }

    enum PluginStatus {
        INACTIVE,
        ACTIVE,
        PAUSED,
        DEPRECATED
    }

    // ============ Structs ============

    struct PluginInfo {
        bytes32 pluginId;
        string name;
        string version;
        PluginType pluginType;
        PluginStatus status;
        address implementation;
        uint256 installedAt;
        uint256 lastUpdated;
    }

    // ============ Events ============

    event PluginExecuted(bytes32 indexed pluginId, address indexed caller, bytes data, bytes result);
    event PluginStatusChanged(bytes32 indexed pluginId, PluginStatus oldStatus, PluginStatus newStatus);

    // ============ Core Functions ============

    /**
     * @notice Get unique plugin identifier
     * @return Plugin ID (keccak256 hash)
     */
    function getPluginId() external pure returns (bytes32);

    /**
     * @notice Get plugin name
     * @return Human-readable plugin name
     */
    function getPluginName() external pure returns (string memory);

    /**
     * @notice Get plugin version
     * @return Semantic version string (e.g., "1.0.0")
     */
    function getPluginVersion() external pure returns (string memory);

    /**
     * @notice Get plugin type category
     * @return Plugin type enum
     */
    function getPluginType() external pure returns (PluginType);

    /**
     * @notice Get required permissions
     * @return Array of permission identifiers
     */
    function getRequiredPermissions() external pure returns (bytes32[] memory);

    /**
     * @notice Get plugin dependencies
     * @return Array of required module/plugin IDs
     */
    function getDependencies() external pure returns (bytes32[] memory);

    /**
     * @notice Initialize plugin with configuration data
     * @param data Initialization parameters (ABI encoded)
     */
    function initialize(bytes calldata data) external;

    /**
     * @notice Main plugin execution function
     * @param data Execution parameters (ABI encoded)
     * @return result Execution result (ABI encoded)
     */
    function execute(bytes calldata data) external returns (bytes memory result);

    /**
     * @notice Check plugin health status
     * @return healthy True if plugin is functioning correctly
     * @return message Status message
     */
    function healthCheck() external view returns (bool healthy, string memory message);

    /**
     * @notice Get current plugin status
     * @return Current status enum
     */
    function getStatus() external view returns (PluginStatus);

    /**
     * @notice Get plugin information
     * @return Plugin info struct
     */
    function getPluginInfo() external view returns (PluginInfo memory);

    // ============ Lifecycle Functions ============

    /**
     * @notice Enable/activate the plugin
     * @dev Only callable by authorized addresses
     */
    function enable() external;

    /**
     * @notice Disable/deactivate the plugin
     * @dev Only callable by authorized addresses
     */
    function disable() external;

    /**
     * @notice Pause plugin execution
     * @dev Only callable by authorized addresses
     */
    function pause() external;

    /**
     * @notice Unpause plugin execution
     * @dev Only callable by authorized addresses
     */
    function unpause() external;

    /**
     * @notice Upgrade plugin to new version
     * @param newImplementation New implementation address
     * @dev Only callable by authorized addresses
     */
    function upgrade(address newImplementation) external;

    // ============ Metadata Functions ============

    /**
     * @notice Get plugin author
     * @return Author address
     */
    function getAuthor() external view returns (address);

    /**
     * @notice Get plugin description
     * @return Description string
     */
    function getDescription() external pure returns (string memory);

    /**
     * @notice Get documentation URI
     * @return URI to plugin documentation
     */
    function getDocumentationURI() external pure returns (string memory);

    /**
     * @notice Get source code URI
     * @return URI to source code repository
     */
    function getSourceURI() external pure returns (string memory);

    /**
     * @notice Check if plugin is compatible with system version
     * @param systemVersion System version string
     * @return compatible True if compatible
     */
    function isCompatible(string calldata systemVersion) external pure returns (bool compatible);
}
