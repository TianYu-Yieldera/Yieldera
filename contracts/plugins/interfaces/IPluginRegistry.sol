// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IPlugin.sol";

/**
 * @title IPluginRegistry
 * @notice Interface for plugin registry and discovery
 * @dev Central registry for all installed plugins
 */
interface IPluginRegistry {
    // ============ Structs ============

    struct PluginRegistration {
        bytes32 pluginId;
        address implementation;
        IPlugin.PluginType pluginType;
        string version;
        address author;
        uint256 registeredAt;
        bool isVerified;
        bool isActive;
    }

    struct PluginStats {
        uint256 totalInstalls;
        uint256 activeInstalls;
        uint256 totalExecutions;
        uint256 failedExecutions;
        uint256 lastExecuted;
    }

    // ============ Events ============

    event PluginRegistered(
        bytes32 indexed pluginId,
        address indexed implementation,
        address indexed author,
        string version
    );
    event PluginUnregistered(bytes32 indexed pluginId, address implementation);
    event PluginVerified(bytes32 indexed pluginId, address verifier);
    event PluginActivated(bytes32 indexed pluginId);
    event PluginDeactivated(bytes32 indexed pluginId);
    event PluginExecutionRecorded(bytes32 indexed pluginId, bool success);

    // ============ Registration Functions ============

    /**
     * @notice Register a new plugin
     * @param plugin Plugin address
     * @return pluginId Registered plugin ID
     */
    function registerPlugin(address plugin) external returns (bytes32 pluginId);

    /**
     * @notice Unregister a plugin
     * @param pluginId Plugin ID to unregister
     */
    function unregisterPlugin(bytes32 pluginId) external;

    /**
     * @notice Mark plugin as verified
     * @param pluginId Plugin ID
     * @dev Only callable by authorized verifiers
     */
    function verifyPlugin(bytes32 pluginId) external;

    /**
     * @notice Activate a registered plugin
     * @param pluginId Plugin ID
     */
    function activatePlugin(bytes32 pluginId) external;

    /**
     * @notice Deactivate a plugin
     * @param pluginId Plugin ID
     */
    function deactivatePlugin(bytes32 pluginId) external;

    // ============ Query Functions ============

    /**
     * @notice Check if plugin is registered
     * @param pluginId Plugin ID
     * @return True if registered
     */
    function isPluginRegistered(bytes32 pluginId) external view returns (bool);

    /**
     * @notice Check if plugin is active
     * @param pluginId Plugin ID
     * @return True if active
     */
    function isPluginActive(bytes32 pluginId) external view returns (bool);

    /**
     * @notice Check if plugin is verified
     * @param pluginId Plugin ID
     * @return True if verified
     */
    function isPluginVerified(bytes32 pluginId) external view returns (bool);

    /**
     * @notice Get plugin registration details
     * @param pluginId Plugin ID
     * @return registration Plugin registration data
     */
    function getPluginRegistration(bytes32 pluginId)
        external
        view
        returns (PluginRegistration memory registration);

    /**
     * @notice Get plugin implementation address
     * @param pluginId Plugin ID
     * @return Plugin contract address
     */
    function getPluginImplementation(bytes32 pluginId) external view returns (address);

    /**
     * @notice Get plugin statistics
     * @param pluginId Plugin ID
     * @return stats Plugin usage statistics
     */
    function getPluginStats(bytes32 pluginId) external view returns (PluginStats memory stats);

    // ============ Discovery Functions ============

    /**
     * @notice Get all registered plugin IDs
     * @return Array of plugin IDs
     */
    function getAllPlugins() external view returns (bytes32[] memory);

    /**
     * @notice Get plugins by type
     * @param pluginType Plugin type filter
     * @return Array of plugin IDs
     */
    function getPluginsByType(IPlugin.PluginType pluginType)
        external
        view
        returns (bytes32[] memory);

    /**
     * @notice Get plugins by author
     * @param author Author address
     * @return Array of plugin IDs
     */
    function getPluginsByAuthor(address author) external view returns (bytes32[] memory);

    /**
     * @notice Get active plugins
     * @return Array of active plugin IDs
     */
    function getActivePlugins() external view returns (bytes32[] memory);

    /**
     * @notice Get verified plugins
     * @return Array of verified plugin IDs
     */
    function getVerifiedPlugins() external view returns (bytes32[] memory);

    /**
     * @notice Search plugins by name pattern
     * @param namePattern Name pattern to match
     * @return Array of matching plugin IDs
     */
    function searchPluginsByName(string calldata namePattern)
        external
        view
        returns (bytes32[] memory);

    // ============ Statistics Functions ============

    /**
     * @notice Record plugin execution
     * @param pluginId Plugin ID
     * @param success Whether execution was successful
     */
    function recordExecution(bytes32 pluginId, bool success) external;

    /**
     * @notice Get total number of registered plugins
     * @return Total count
     */
    function getTotalPluginCount() external view returns (uint256);

    /**
     * @notice Get number of active plugins
     * @return Active count
     */
    function getActivePluginCount() external view returns (uint256);

    /**
     * @notice Get number of verified plugins
     * @return Verified count
     */
    function getVerifiedPluginCount() external view returns (uint256);

    // ============ Compatibility Functions ============

    /**
     * @notice Check if plugin is compatible with system
     * @param pluginId Plugin ID
     * @param systemVersion System version
     * @return compatible True if compatible
     */
    function checkCompatibility(bytes32 pluginId, string calldata systemVersion)
        external
        view
        returns (bool compatible);

    /**
     * @notice Check if plugin dependencies are satisfied
     * @param pluginId Plugin ID
     * @return satisfied True if all dependencies available
     * @return missingDeps Array of missing dependency IDs
     */
    function checkDependencies(bytes32 pluginId)
        external
        view
        returns (bool satisfied, bytes32[] memory missingDeps);
}
