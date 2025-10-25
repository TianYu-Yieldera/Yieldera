// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/core/IAccessController.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

/**
 * @title AccessController
 * @notice Unified access control system for the pluggable architecture
 * @dev Implements role-based and permission-based access control with timelock
 */
contract AccessController is IAccessController, Ownable, Pausable, ReentrancyGuard {
    // ============ Constants ============

    bytes32 public constant override DEFAULT_ADMIN_ROLE = 0x00;
    bytes32 public constant override EMERGENCY_ROLE = keccak256("EMERGENCY_ROLE");

    // ============ State Variables ============

    // Role storage
    mapping(bytes32 => Role) private _roles;
    mapping(bytes32 => mapping(address => bool)) private _roleMembers;
    mapping(address => bytes32[]) private _accountRoles;

    // Permission storage
    mapping(bytes32 => Permission) private _permissions;
    mapping(bytes32 => mapping(bytes32 => bool)) private _rolePermissions; // roleId => permissionId => hasPermission

    // Module-specific access control
    mapping(bytes32 => mapping(bytes32 => mapping(address => bool))) private _moduleRoles; // moduleId => roleId => account => hasRole

    // Time-locked operations
    mapping(bytes32 => TimeLockOperation) private _operations;
    mapping(address => bytes32[]) private _executorOperations;
    uint256 private _operationCounter;

    // Emergency pause
    bool private _emergencyPaused;

    // Revoked modules
    mapping(bytes32 => bool) private _revokedModules;

    // Statistics
    uint256 private _totalRoles;
    uint256 private _totalPermissions;
    uint256 private _activeRoles;

    // ============ Constructor ============

    constructor() {
        // Setup default admin role
        _roles[DEFAULT_ADMIN_ROLE] = Role({
            roleId: DEFAULT_ADMIN_ROLE,
            name: "DEFAULT_ADMIN",
            permissions: new bytes32[](0),
            adminRole: DEFAULT_ADMIN_ROLE,
            isActive: true
        });

        // Grant default admin to deployer
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);

        // Setup emergency role
        _createRole(EMERGENCY_ROLE, "EMERGENCY", DEFAULT_ADMIN_ROLE);

        _totalRoles = 2;
        _activeRoles = 2;
    }

    // ============ Modifiers ============

    modifier onlyRole(bytes32 roleId) {
        require(hasRole(roleId, msg.sender), "AccessController: missing role");
        _;
    }

    modifier notEmergencyPaused() {
        require(!_emergencyPaused, "AccessController: emergency paused");
        _;
    }

    // ============ Role Management ============

    function createRole(bytes32 roleId, string calldata name, bytes32 adminRole)
        external
        override
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        _createRole(roleId, name, adminRole);
    }

    function grantRole(bytes32 roleId, address account)
        external
        override
        notEmergencyPaused
    {
        bytes32 adminRole = _roles[roleId].adminRole;
        require(hasRole(adminRole, msg.sender), "AccessController: not admin");
        _grantRole(roleId, account);
    }

    function revokeRole(bytes32 roleId, address account)
        external
        override
        notEmergencyPaused
    {
        bytes32 adminRole = _roles[roleId].adminRole;
        require(hasRole(adminRole, msg.sender), "AccessController: not admin");
        _revokeRole(roleId, account);
    }

    function renounceRole(bytes32 roleId) external override {
        _revokeRole(roleId, msg.sender);
    }

    function hasRole(bytes32 roleId, address account)
        public
        view
        override
        returns (bool)
    {
        return _roleMembers[roleId][account];
    }

    function getRole(bytes32 roleId)
        external
        view
        override
        returns (Role memory)
    {
        return _roles[roleId];
    }

    function getRolesForAccount(address account)
        external
        view
        override
        returns (bytes32[] memory)
    {
        return _accountRoles[account];
    }

    // ============ Permission Management ============

    function createPermission(
        bytes32 permissionId,
        string calldata name,
        string calldata description,
        bool requiresTimelock,
        uint256 timelockDuration
    ) external override onlyRole(DEFAULT_ADMIN_ROLE) {
        require(!_permissions[permissionId].isActive, "Permission already exists");

        _permissions[permissionId] = Permission({
            permissionId: permissionId,
            name: name,
            description: description,
            requiresTimelock: requiresTimelock,
            timelockDuration: timelockDuration,
            isActive: true
        });

        _totalPermissions++;

        emit PermissionCreated(permissionId, name, requiresTimelock);
    }

    function grantPermissionToRole(bytes32 permissionId, bytes32 roleId)
        external
        override
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        require(_permissions[permissionId].isActive, "Permission not active");
        require(_roles[roleId].isActive, "Role not active");
        require(!_rolePermissions[roleId][permissionId], "Permission already granted");

        _rolePermissions[roleId][permissionId] = true;
        _roles[roleId].permissions.push(permissionId);

        emit PermissionGranted(permissionId, roleId, msg.sender);
    }

    function revokePermissionFromRole(bytes32 permissionId, bytes32 roleId)
        external
        override
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        require(_rolePermissions[roleId][permissionId], "Permission not granted");

        _rolePermissions[roleId][permissionId] = false;

        // Remove from permissions array
        bytes32[] storage perms = _roles[roleId].permissions;
        for (uint256 i = 0; i < perms.length; i++) {
            if (perms[i] == permissionId) {
                perms[i] = perms[perms.length - 1];
                perms.pop();
                break;
            }
        }

        emit PermissionRevoked(permissionId, roleId, msg.sender);
    }

    function hasPermission(bytes32 permissionId, address account)
        public
        view
        override
        returns (bool)
    {
        // Check all roles of the account
        bytes32[] memory roles = _accountRoles[account];
        for (uint256 i = 0; i < roles.length; i++) {
            if (_rolePermissions[roles[i]][permissionId]) {
                return true;
            }
        }
        return false;
    }

    function hasModulePermission(
        bytes32 moduleId,
        bytes32 permissionId,
        address account
    ) external view override returns (bool) {
        // Check if module is revoked
        if (_revokedModules[moduleId]) {
            return false;
        }

        // Check module-specific roles
        bytes32[] memory roles = _accountRoles[account];
        for (uint256 i = 0; i < roles.length; i++) {
            if (_moduleRoles[moduleId][roles[i]][account]) {
                if (_rolePermissions[roles[i]][permissionId]) {
                    return true;
                }
            }
        }

        // Fall back to global permission check
        return hasPermission(permissionId, account);
    }

    function getPermission(bytes32 permissionId)
        external
        view
        override
        returns (Permission memory)
    {
        return _permissions[permissionId];
    }

    // ============ Time-Locked Operations ============

    function scheduleOperation(
        bytes32 permission,
        address target,
        bytes calldata callData
    ) external override notEmergencyPaused returns (bytes32 operationId) {
        require(hasPermission(permission, msg.sender), "Missing permission");

        Permission memory perm = _permissions[permission];
        require(perm.requiresTimelock, "Permission does not require timelock");

        operationId = keccak256(
            abi.encodePacked(msg.sender, target, callData, block.timestamp, _operationCounter++)
        );

        uint256 executeAfter = block.timestamp + perm.timelockDuration;

        _operations[operationId] = TimeLockOperation({
            operationId: operationId,
            executor: msg.sender,
            permission: permission,
            callData: callData,
            scheduledTime: block.timestamp,
            executeAfter: executeAfter,
            executed: false,
            cancelled: false
        });

        _executorOperations[msg.sender].push(operationId);

        emit OperationScheduled(operationId, msg.sender, executeAfter);

        return operationId;
    }

    function executeOperation(bytes32 operationId)
        external
        override
        notEmergencyPaused
        nonReentrant
        returns (bytes memory)
    {
        TimeLockOperation storage operation = _operations[operationId];

        require(operation.executor != address(0), "Operation not found");
        require(!operation.executed, "Operation already executed");
        require(!operation.cancelled, "Operation cancelled");
        require(block.timestamp >= operation.executeAfter, "Timelock not expired");
        require(msg.sender == operation.executor, "Not operation executor");

        operation.executed = true;

        emit OperationExecuted(operationId, msg.sender);

        // Note: In production, this would execute the actual call
        // For now, we just return empty bytes
        return "";
    }

    function cancelOperation(bytes32 operationId) external override {
        TimeLockOperation storage operation = _operations[operationId];

        require(operation.executor != address(0), "Operation not found");
        require(!operation.executed, "Operation already executed");
        require(!operation.cancelled, "Operation already cancelled");
        require(
            msg.sender == operation.executor || hasRole(DEFAULT_ADMIN_ROLE, msg.sender),
            "Not authorized to cancel"
        );

        operation.cancelled = true;

        emit OperationCancelled(operationId, msg.sender);
    }

    function getOperation(bytes32 operationId)
        external
        view
        override
        returns (TimeLockOperation memory)
    {
        return _operations[operationId];
    }

    function isOperationReady(bytes32 operationId)
        external
        view
        override
        returns (bool)
    {
        TimeLockOperation memory operation = _operations[operationId];
        return !operation.executed &&
            !operation.cancelled &&
            block.timestamp >= operation.executeAfter;
    }

    function getPendingOperations(address executor)
        external
        view
        override
        returns (bytes32[] memory)
    {
        bytes32[] memory allOps = _executorOperations[executor];
        uint256 pendingCount = 0;

        // Count pending operations
        for (uint256 i = 0; i < allOps.length; i++) {
            TimeLockOperation memory op = _operations[allOps[i]];
            if (!op.executed && !op.cancelled) {
                pendingCount++;
            }
        }

        // Collect pending operations
        bytes32[] memory pending = new bytes32[](pendingCount);
        uint256 index = 0;
        for (uint256 i = 0; i < allOps.length; i++) {
            TimeLockOperation memory op = _operations[allOps[i]];
            if (!op.executed && !op.cancelled) {
                pending[index] = allOps[i];
                index++;
            }
        }

        return pending;
    }

    // ============ Module-Specific Access Control ============

    function grantModuleRole(bytes32 moduleId, bytes32 roleId, address account)
        external
        override
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        require(!_revokedModules[moduleId], "Module access revoked");
        require(_roles[roleId].isActive, "Role not active");

        _moduleRoles[moduleId][roleId][account] = true;
    }

    function revokeModuleRole(bytes32 moduleId, bytes32 roleId, address account)
        external
        override
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        _moduleRoles[moduleId][roleId][account] = false;
    }

    function hasModuleRole(
        bytes32 moduleId,
        bytes32 roleId,
        address account
    ) external view override returns (bool) {
        if (_revokedModules[moduleId]) {
            return false;
        }
        return _moduleRoles[moduleId][roleId][account];
    }

    // ============ Emergency Controls ============

    function activateEmergencyPause()
        external
        override
        onlyRole(EMERGENCY_ROLE)
    {
        _emergencyPaused = true;
        emit EmergencyPauseActivated(msg.sender);
    }

    function deactivateEmergencyPause()
        external
        override
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        _emergencyPaused = false;
        emit EmergencyPauseDeactivated(msg.sender);
    }

    function isEmergencyPaused() external view override returns (bool) {
        return _emergencyPaused;
    }

    function emergencyRevokeModuleAccess(bytes32 moduleId)
        external
        override
        onlyRole(EMERGENCY_ROLE)
    {
        _revokedModules[moduleId] = true;
        emit ModuleAccessRevoked(moduleId, msg.sender);
    }

    // ============ View Functions ============

    function getAccessControlStats()
        external
        view
        override
        returns (
            uint256 totalRoles,
            uint256 totalPermissions,
            uint256 activeRoles
        )
    {
        return (_totalRoles, _totalPermissions, _activeRoles);
    }

    // ============ Internal Functions ============

    function _createRole(bytes32 roleId, string memory name, bytes32 adminRole) internal {
        require(!_roles[roleId].isActive, "Role already exists");
        require(_roles[adminRole].isActive, "Admin role not active");

        _roles[roleId] = Role({
            roleId: roleId,
            name: name,
            permissions: new bytes32[](0),
            adminRole: adminRole,
            isActive: true
        });

        _totalRoles++;
        _activeRoles++;

        emit RoleCreated(roleId, name, adminRole);
    }

    function _grantRole(bytes32 roleId, address account) internal {
        require(_roles[roleId].isActive, "Role not active");
        require(!_roleMembers[roleId][account], "Role already granted");

        _roleMembers[roleId][account] = true;
        _accountRoles[account].push(roleId);

        emit RoleGranted(roleId, account, msg.sender);
    }

    function _revokeRole(bytes32 roleId, address account) internal {
        require(_roleMembers[roleId][account], "Role not granted");

        _roleMembers[roleId][account] = false;

        // Remove from account roles array
        bytes32[] storage roles = _accountRoles[account];
        for (uint256 i = 0; i < roles.length; i++) {
            if (roles[i] == roleId) {
                roles[i] = roles[roles.length - 1];
                roles.pop();
                break;
            }
        }

        emit RoleRevoked(roleId, account, msg.sender);
    }
}
