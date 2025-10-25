// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IPluginPermissionManager
 * @notice Interface for managing plugin permissions
 * @dev Controls what actions plugins are allowed to perform
 */
interface IPluginPermissionManager {
    // ============ Enums ============

    enum Permission {
        NONE,
        READ,
        WRITE,
        ADMIN,
        TRANSFER,
        MINT,
        BURN,
        UPGRADE,
        DELEGATE
    }

    // ============ Structs ============

    struct PermissionGrant {
        bytes32 pluginId;
        Permission permission;
        address target; // Target module/contract
        uint256 grantedAt;
        uint256 expiresAt; // 0 = no expiry
        bool isActive;
    }

    // ============ Events ============

    event PermissionGranted(
        bytes32 indexed pluginId,
        Permission permission,
        address indexed target,
        uint256 expiresAt
    );
    event PermissionRevoked(bytes32 indexed pluginId, Permission permission, address indexed target);
    event PermissionExpired(bytes32 indexed pluginId, Permission permission, address indexed target);

    // ============ Grant/Revoke Functions ============

    /**
     * @notice Grant permission to plugin
     * @param pluginId Plugin ID
     * @param permission Permission to grant
     * @param target Target contract (address(0) for global)
     * @param duration Permission duration in seconds (0 for permanent)
     */
    function grantPermission(
        bytes32 pluginId,
        Permission permission,
        address target,
        uint256 duration
    ) external;

    /**
     * @notice Grant multiple permissions to plugin
     * @param pluginId Plugin ID
     * @param permissions Array of permissions
     * @param targets Array of target contracts
     * @param durations Array of durations
     */
    function grantPermissions(
        bytes32 pluginId,
        Permission[] calldata permissions,
        address[] calldata targets,
        uint256[] calldata durations
    ) external;

    /**
     * @notice Revoke permission from plugin
     * @param pluginId Plugin ID
     * @param permission Permission to revoke
     * @param target Target contract
     */
    function revokePermission(bytes32 pluginId, Permission permission, address target) external;

    /**
     * @notice Revoke all permissions from plugin
     * @param pluginId Plugin ID
     */
    function revokeAllPermissions(bytes32 pluginId) external;

    // ============ Check Functions ============

    /**
     * @notice Check if plugin has permission
     * @param pluginId Plugin ID
     * @param permission Permission to check
     * @param target Target contract
     * @return True if plugin has permission
     */
    function hasPermission(bytes32 pluginId, Permission permission, address target)
        external
        view
        returns (bool);

    /**
     * @notice Check if plugin has all required permissions
     * @param pluginId Plugin ID
     * @param permissions Array of required permissions
     * @param targets Array of target contracts
     * @return True if plugin has all permissions
     */
    function hasPermissions(
        bytes32 pluginId,
        Permission[] calldata permissions,
        address[] calldata targets
    ) external view returns (bool);

    /**
     * @notice Require permission or revert
     * @param pluginId Plugin ID
     * @param permission Required permission
     * @param target Target contract
     */
    function requirePermission(bytes32 pluginId, Permission permission, address target)
        external
        view;

    /**
     * @notice Check if permission has expired
     * @param pluginId Plugin ID
     * @param permission Permission to check
     * @param target Target contract
     * @return True if expired
     */
    function isPermissionExpired(bytes32 pluginId, Permission permission, address target)
        external
        view
        returns (bool);

    // ============ Query Functions ============

    /**
     * @notice Get all permissions for plugin
     * @param pluginId Plugin ID
     * @return Array of permission grants
     */
    function getPluginPermissions(bytes32 pluginId)
        external
        view
        returns (PermissionGrant[] memory);

    /**
     * @notice Get plugins with specific permission
     * @param permission Permission type
     * @param target Target contract
     * @return Array of plugin IDs
     */
    function getPluginsWithPermission(Permission permission, address target)
        external
        view
        returns (bytes32[] memory);

    /**
     * @notice Get permission expiry time
     * @param pluginId Plugin ID
     * @param permission Permission type
     * @param target Target contract
     * @return Expiry timestamp (0 if permanent or not granted)
     */
    function getPermissionExpiry(bytes32 pluginId, Permission permission, address target)
        external
        view
        returns (uint256);

    // ============ Delegation Functions ============

    /**
     * @notice Allow plugin to delegate permission to another plugin
     * @param fromPlugin Source plugin ID
     * @param toPlugin Target plugin ID
     * @param permission Permission to delegate
     * @param target Target contract
     */
    function delegatePermission(
        bytes32 fromPlugin,
        bytes32 toPlugin,
        Permission permission,
        address target
    ) external;

    /**
     * @notice Revoke delegated permission
     * @param fromPlugin Source plugin ID
     * @param toPlugin Target plugin ID
     * @param permission Permission that was delegated
     * @param target Target contract
     */
    function revokeDelegation(
        bytes32 fromPlugin,
        bytes32 toPlugin,
        Permission permission,
        address target
    ) external;

    // ============ Batch Operations ============

    /**
     * @notice Clean up expired permissions
     * @param pluginIds Array of plugin IDs to clean (empty for all)
     */
    function cleanupExpiredPermissions(bytes32[] calldata pluginIds) external;

    /**
     * @notice Extend permission duration
     * @param pluginId Plugin ID
     * @param permission Permission type
     * @param target Target contract
     * @param additionalDuration Additional duration in seconds
     */
    function extendPermission(
        bytes32 pluginId,
        Permission permission,
        address target,
        uint256 additionalDuration
    ) external;
}
