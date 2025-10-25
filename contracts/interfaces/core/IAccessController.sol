// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IAccessController
 * @notice Unified access control system for the pluggable architecture
 * @dev Combines role-based and permission-based access control
 *
 * Features:
 * - Role-based access control (RBAC)
 * - Fine-grained permissions
 * - Module-specific access control
 * - Time-locked operations
 * - Emergency pause capabilities
 */
interface IAccessController {
    /**
     * @notice Permission structure for fine-grained access control
     */
    struct Permission {
        bytes32 permissionId;       // Unique permission identifier
        string name;                // Human-readable name
        string description;         // Permission description
        bool requiresTimelock;      // Whether this permission requires timelock
        uint256 timelockDuration;   // Timelock duration in seconds
        bool isActive;              // Whether permission is currently active
    }

    /**
     * @notice Role structure
     */
    struct Role {
        bytes32 roleId;             // Unique role identifier
        string name;                // Human-readable name
        bytes32[] permissions;      // Permissions granted to this role
        bytes32 adminRole;          // Role that can grant/revoke this role
        bool isActive;              // Whether role is currently active
    }

    /**
     * @notice Time-locked operation
     */
    struct TimeLockOperation {
        bytes32 operationId;        // Unique operation identifier
        address executor;           // Address that initiated the operation
        bytes32 permission;         // Required permission
        bytes callData;             // Encoded function call
        uint256 scheduledTime;      // When operation was scheduled
        uint256 executeAfter;       // Earliest execution time
        bool executed;              // Whether operation has been executed
        bool cancelled;             // Whether operation was cancelled
    }

    // ============ Events ============

    event RoleGranted(bytes32 indexed roleId, address indexed account, address indexed sender);
    event RoleRevoked(bytes32 indexed roleId, address indexed account, address indexed sender);
    event RoleCreated(bytes32 indexed roleId, string name, bytes32 indexed adminRole);

    event PermissionGranted(bytes32 indexed permissionId, bytes32 indexed roleId, address indexed sender);
    event PermissionRevoked(bytes32 indexed permissionId, bytes32 indexed roleId, address indexed sender);
    event PermissionCreated(bytes32 indexed permissionId, string name, bool requiresTimelock);

    event OperationScheduled(bytes32 indexed operationId, address indexed executor, uint256 executeAfter);
    event OperationExecuted(bytes32 indexed operationId, address indexed executor);
    event OperationCancelled(bytes32 indexed operationId, address indexed canceller);

    event EmergencyPauseActivated(address indexed activator);
    event EmergencyPauseDeactivated(address indexed deactivator);
    event ModuleAccessRevoked(bytes32 indexed moduleId, address indexed revoker);

    // ============ Role Management ============

    /**
     * @notice Create a new role
     * @param roleId Unique identifier for the role
     * @param name Human-readable role name
     * @param adminRole Role that can grant/revoke this role
     * @dev Only callable by DEFAULT_ADMIN_ROLE
     */
    function createRole(bytes32 roleId, string calldata name, bytes32 adminRole) external;

    /**
     * @notice Grant a role to an account
     * @param roleId Role identifier
     * @param account Address to grant role to
     * @dev Caller must have the role's admin role
     */
    function grantRole(bytes32 roleId, address account) external;

    /**
     * @notice Revoke a role from an account
     * @param roleId Role identifier
     * @param account Address to revoke role from
     * @dev Caller must have the role's admin role
     */
    function revokeRole(bytes32 roleId, address account) external;

    /**
     * @notice Renounce a role (caller gives up their own role)
     * @param roleId Role identifier
     */
    function renounceRole(bytes32 roleId) external;

    /**
     * @notice Check if account has a specific role
     * @param roleId Role identifier
     * @param account Address to check
     * @return True if account has the role
     */
    function hasRole(bytes32 roleId, address account) external view returns (bool);

    /**
     * @notice Get role information
     * @param roleId Role identifier
     * @return Role structure
     */
    function getRole(bytes32 roleId) external view returns (Role memory);

    /**
     * @notice Get all roles for an account
     * @param account Address to check
     * @return Array of role IDs
     */
    function getRolesForAccount(address account) external view returns (bytes32[] memory);

    // ============ Permission Management ============

    /**
     * @notice Create a new permission
     * @param permissionId Unique identifier
     * @param name Human-readable name
     * @param description Permission description
     * @param requiresTimelock Whether timelock is required
     * @param timelockDuration Timelock duration in seconds (if applicable)
     */
    function createPermission(
        bytes32 permissionId,
        string calldata name,
        string calldata description,
        bool requiresTimelock,
        uint256 timelockDuration
    ) external;

    /**
     * @notice Grant permission to a role
     * @param permissionId Permission identifier
     * @param roleId Role identifier
     */
    function grantPermissionToRole(bytes32 permissionId, bytes32 roleId) external;

    /**
     * @notice Revoke permission from a role
     * @param permissionId Permission identifier
     * @param roleId Role identifier
     */
    function revokePermissionFromRole(bytes32 permissionId, bytes32 roleId) external;

    /**
     * @notice Check if account has a specific permission
     * @param permissionId Permission identifier
     * @param account Address to check
     * @return True if account has the permission
     */
    function hasPermission(bytes32 permissionId, address account) external view returns (bool);

    /**
     * @notice Check if account has permission for a specific module
     * @param moduleId Module identifier
     * @param permissionId Permission identifier
     * @param account Address to check
     * @return True if account has module-specific permission
     */
    function hasModulePermission(
        bytes32 moduleId,
        bytes32 permissionId,
        address account
    ) external view returns (bool);

    /**
     * @notice Get permission information
     * @param permissionId Permission identifier
     * @return Permission structure
     */
    function getPermission(bytes32 permissionId) external view returns (Permission memory);

    // ============ Time-Locked Operations ============

    /**
     * @notice Schedule a time-locked operation
     * @param permission Required permission
     * @param target Target contract address
     * @param callData Encoded function call
     * @return operationId Unique operation identifier
     */
    function scheduleOperation(
        bytes32 permission,
        address target,
        bytes calldata callData
    ) external returns (bytes32 operationId);

    /**
     * @notice Execute a scheduled operation
     * @param operationId Operation identifier
     * @dev Can only be executed after timelock expires
     */
    function executeOperation(bytes32 operationId) external returns (bytes memory);

    /**
     * @notice Cancel a scheduled operation
     * @param operationId Operation identifier
     * @dev Only callable by operation creator or admin
     */
    function cancelOperation(bytes32 operationId) external;

    /**
     * @notice Get operation details
     * @param operationId Operation identifier
     * @return Operation structure
     */
    function getOperation(bytes32 operationId) external view returns (TimeLockOperation memory);

    /**
     * @notice Check if operation is ready to execute
     * @param operationId Operation identifier
     * @return True if operation can be executed
     */
    function isOperationReady(bytes32 operationId) external view returns (bool);

    /**
     * @notice Get pending operations for an executor
     * @param executor Address to check
     * @return Array of operation IDs
     */
    function getPendingOperations(address executor) external view returns (bytes32[] memory);

    // ============ Module-Specific Access Control ============

    /**
     * @notice Grant module-specific role
     * @param moduleId Module identifier
     * @param roleId Role identifier
     * @param account Address to grant role to
     */
    function grantModuleRole(bytes32 moduleId, bytes32 roleId, address account) external;

    /**
     * @notice Revoke module-specific role
     * @param moduleId Module identifier
     * @param roleId Role identifier
     * @param account Address to revoke role from
     */
    function revokeModuleRole(bytes32 moduleId, bytes32 roleId, address account) external;

    /**
     * @notice Check if account has module-specific role
     * @param moduleId Module identifier
     * @param roleId Role identifier
     * @param account Address to check
     * @return True if account has the module role
     */
    function hasModuleRole(
        bytes32 moduleId,
        bytes32 roleId,
        address account
    ) external view returns (bool);

    // ============ Emergency Controls ============

    /**
     * @notice Activate emergency pause
     * @dev Pauses all operations across all modules
     * @dev Only callable by EMERGENCY_ROLE
     */
    function activateEmergencyPause() external;

    /**
     * @notice Deactivate emergency pause
     * @dev Only callable by DEFAULT_ADMIN_ROLE
     */
    function deactivateEmergencyPause() external;

    /**
     * @notice Check if system is in emergency pause
     * @return True if emergency pause is active
     */
    function isEmergencyPaused() external view returns (bool);

    /**
     * @notice Revoke all access for a specific module
     * @param moduleId Module identifier
     * @dev Emergency function to isolate a compromised module
     */
    function emergencyRevokeModuleAccess(bytes32 moduleId) external;

    // ============ View Functions ============

    /**
     * @notice Get default admin role identifier
     * @return Default admin role ID
     */
    function DEFAULT_ADMIN_ROLE() external pure returns (bytes32);

    /**
     * @notice Get emergency role identifier
     * @return Emergency role ID
     */
    function EMERGENCY_ROLE() external pure returns (bytes32);

    /**
     * @notice Get access control statistics
     * @return totalRoles Total number of roles
     * @return totalPermissions Total number of permissions
     * @return activeRoles Number of active roles
     */
    function getAccessControlStats()
        external
        view
        returns (
            uint256 totalRoles,
            uint256 totalPermissions,
            uint256 activeRoles
        );
}
