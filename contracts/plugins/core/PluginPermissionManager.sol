// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../interfaces/IPluginPermissionManager.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title PluginPermissionManager
 * @notice Manages permissions for plugins
 * @dev Controls what actions plugins are allowed to perform
 */
contract PluginPermissionManager is IPluginPermissionManager, Ownable {
    // ============ Storage ============

    // pluginId => permission => target => PermissionGrant
    mapping(bytes32 => mapping(Permission => mapping(address => PermissionGrant))) private
        _permissions;

    // pluginId => array of all grants
    mapping(bytes32 => PermissionGrant[]) private _pluginPermissions;

    // permission => target => array of plugin IDs
    mapping(Permission => mapping(address => bytes32[])) private _pluginsWithPermission;

    // ============ Grant/Revoke Functions ============

    function grantPermission(
        bytes32 pluginId,
        Permission permission,
        address target,
        uint256 duration
    ) external override onlyOwner {
        require(permission != Permission.NONE, "Invalid permission");

        uint256 expiresAt = duration > 0 ? block.timestamp + duration : 0;

        PermissionGrant storage grant = _permissions[pluginId][permission][target];

        // If new permission, add to arrays
        if (!grant.isActive) {
            grant.pluginId = pluginId;
            grant.permission = permission;
            grant.target = target;
            grant.grantedAt = block.timestamp;
            grant.isActive = true;

            _pluginPermissions[pluginId].push(grant);
            _pluginsWithPermission[permission][target].push(pluginId);
        }

        grant.expiresAt = expiresAt;

        emit PermissionGranted(pluginId, permission, target, expiresAt);
    }

    function grantPermissions(
        bytes32 pluginId,
        Permission[] calldata permissions,
        address[] calldata targets,
        uint256[] calldata durations
    ) external override onlyOwner {
        require(
            permissions.length == targets.length && targets.length == durations.length,
            "Array length mismatch"
        );

        for (uint256 i = 0; i < permissions.length; i++) {
            this.grantPermission(pluginId, permissions[i], targets[i], durations[i]);
        }
    }

    function revokePermission(bytes32 pluginId, Permission permission, address target)
        external
        override
        onlyOwner
    {
        PermissionGrant storage grant = _permissions[pluginId][permission][target];
        require(grant.isActive, "Permission not granted");

        grant.isActive = false;

        // Remove from pluginsWithPermission array
        bytes32[] storage plugins = _pluginsWithPermission[permission][target];
        for (uint256 i = 0; i < plugins.length; i++) {
            if (plugins[i] == pluginId) {
                plugins[i] = plugins[plugins.length - 1];
                plugins.pop();
                break;
            }
        }

        emit PermissionRevoked(pluginId, permission, target);
    }

    function revokeAllPermissions(bytes32 pluginId) external override onlyOwner {
        PermissionGrant[] storage grants = _pluginPermissions[pluginId];

        for (uint256 i = 0; i < grants.length; i++) {
            if (grants[i].isActive) {
                grants[i].isActive = false;
                emit PermissionRevoked(grants[i].pluginId, grants[i].permission, grants[i].target);
            }
        }
    }

    // ============ Check Functions ============

    function hasPermission(bytes32 pluginId, Permission permission, address target)
        public
        view
        override
        returns (bool)
    {
        PermissionGrant storage grant = _permissions[pluginId][permission][target];

        if (!grant.isActive) return false;

        // Check expiry
        if (grant.expiresAt > 0 && block.timestamp > grant.expiresAt) {
            return false;
        }

        return true;
    }

    function hasPermissions(
        bytes32 pluginId,
        Permission[] calldata permissions,
        address[] calldata targets
    ) external view override returns (bool) {
        require(permissions.length == targets.length, "Array length mismatch");

        for (uint256 i = 0; i < permissions.length; i++) {
            if (!hasPermission(pluginId, permissions[i], targets[i])) {
                return false;
            }
        }

        return true;
    }

    function requirePermission(bytes32 pluginId, Permission permission, address target)
        external
        view
        override
    {
        require(hasPermission(pluginId, permission, target), "Permission denied");
    }

    function isPermissionExpired(bytes32 pluginId, Permission permission, address target)
        external
        view
        override
        returns (bool)
    {
        PermissionGrant storage grant = _permissions[pluginId][permission][target];

        if (!grant.isActive) return true;
        if (grant.expiresAt == 0) return false;

        return block.timestamp > grant.expiresAt;
    }

    // ============ Query Functions ============

    function getPluginPermissions(bytes32 pluginId)
        external
        view
        override
        returns (PermissionGrant[] memory)
    {
        PermissionGrant[] storage grants = _pluginPermissions[pluginId];

        // Count active permissions
        uint256 activeCount = 0;
        for (uint256 i = 0; i < grants.length; i++) {
            if (grants[i].isActive) {
                // Check not expired
                if (grants[i].expiresAt == 0 || block.timestamp <= grants[i].expiresAt) {
                    activeCount++;
                }
            }
        }

        // Create result array
        PermissionGrant[] memory result = new PermissionGrant[](activeCount);
        uint256 index = 0;

        for (uint256 i = 0; i < grants.length; i++) {
            if (grants[i].isActive) {
                if (grants[i].expiresAt == 0 || block.timestamp <= grants[i].expiresAt) {
                    result[index++] = grants[i];
                }
            }
        }

        return result;
    }

    function getPluginsWithPermission(Permission permission, address target)
        external
        view
        override
        returns (bytes32[] memory)
    {
        bytes32[] storage plugins = _pluginsWithPermission[permission][target];

        // Count active permissions
        uint256 activeCount = 0;
        for (uint256 i = 0; i < plugins.length; i++) {
            if (hasPermission(plugins[i], permission, target)) {
                activeCount++;
            }
        }

        // Create result array
        bytes32[] memory result = new bytes32[](activeCount);
        uint256 index = 0;

        for (uint256 i = 0; i < plugins.length; i++) {
            if (hasPermission(plugins[i], permission, target)) {
                result[index++] = plugins[i];
            }
        }

        return result;
    }

    function getPermissionExpiry(bytes32 pluginId, Permission permission, address target)
        external
        view
        override
        returns (uint256)
    {
        return _permissions[pluginId][permission][target].expiresAt;
    }

    // ============ Delegation Functions ============

    function delegatePermission(
        bytes32 fromPlugin,
        bytes32 toPlugin,
        Permission permission,
        address target
    ) external override {
        // Simplified: only owner can delegate
        require(msg.sender == owner(), "Only owner can delegate");
        require(hasPermission(fromPlugin, permission, target), "Source lacks permission");

        // Grant permission to target plugin
        this.grantPermission(toPlugin, permission, target, 0);
    }

    function revokeDelegation(bytes32, bytes32 toPlugin, Permission permission, address target)
        external
        override
        onlyOwner
    {
        this.revokePermission(toPlugin, permission, target);
    }

    // ============ Batch Operations ============

    function cleanupExpiredPermissions(bytes32[] calldata pluginIds) external override {
        bytes32[] memory plugins = pluginIds.length > 0 ? pluginIds : new bytes32[](0);

        // If no specific plugins provided, this would iterate all (expensive)
        // For simplicity, only cleanup specified plugins

        for (uint256 i = 0; i < plugins.length; i++) {
            PermissionGrant[] storage grants = _pluginPermissions[plugins[i]];

            for (uint256 j = 0; j < grants.length; j++) {
                if (grants[j].isActive && grants[j].expiresAt > 0) {
                    if (block.timestamp > grants[j].expiresAt) {
                        grants[j].isActive = false;
                        emit PermissionExpired(
                            grants[j].pluginId, grants[j].permission, grants[j].target
                        );
                    }
                }
            }
        }
    }

    function extendPermission(
        bytes32 pluginId,
        Permission permission,
        address target,
        uint256 additionalDuration
    ) external override onlyOwner {
        PermissionGrant storage grant = _permissions[pluginId][permission][target];
        require(grant.isActive, "Permission not granted");
        require(additionalDuration > 0, "Invalid duration");

        if (grant.expiresAt == 0) {
            // Permanent permission, set expiry
            grant.expiresAt = block.timestamp + additionalDuration;
        } else {
            // Extend existing expiry
            grant.expiresAt += additionalDuration;
        }

        emit PermissionGranted(pluginId, permission, target, grant.expiresAt);
    }
}
