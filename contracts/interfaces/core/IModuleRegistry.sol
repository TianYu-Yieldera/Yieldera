// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IModule.sol";

/**
 * @title IModuleRegistry
 * @notice Central registry for managing all pluggable modules
 * @dev Acts as the service locator and dependency injection container
 *
 * Responsibilities:
 * - Register/unregister modules
 * - Validate module dependencies
 * - Track module lifecycle
 * - Provide module discovery
 */
interface IModuleRegistry {
    /**
     * @notice Module registration data
     */
    struct ModuleRegistration {
        address moduleAddress;      // Module contract address
        bytes32 moduleId;           // Unique identifier
        string name;                // Module name
        string version;             // Module version
        bytes32[] dependencies;     // Required module dependencies
        bool isRegistered;          // Registration status
        uint256 registeredAt;       // Registration timestamp
    }

    // ============ Events ============

    event ModuleRegistered(
        bytes32 indexed moduleId,
        address indexed moduleAddress,
        string name,
        string version
    );

    event ModuleUnregistered(
        bytes32 indexed moduleId,
        address indexed moduleAddress
    );

    event ModuleEnabled(
        bytes32 indexed moduleId,
        address indexed moduleAddress
    );

    event ModuleDisabled(
        bytes32 indexed moduleId,
        address indexed moduleAddress
    );

    event ModuleUpgraded(
        bytes32 indexed moduleId,
        address indexed oldAddress,
        address indexed newAddress,
        string newVersion
    );

    event DependencyValidated(
        bytes32 indexed moduleId,
        bytes32 indexed dependencyId,
        bool satisfied
    );

    // ============ Module Management ============

    /**
     * @notice Register a new module
     * @param moduleAddress Address of the module contract
     * @return moduleId The assigned module ID
     * @dev Validates that module implements IModule interface
     * @dev Checks that all dependencies are satisfied
     */
    function registerModule(address moduleAddress) external returns (bytes32 moduleId);

    /**
     * @notice Unregister a module
     * @param moduleId Module identifier
     * @dev Only possible if no other modules depend on this one
     */
    function unregisterModule(bytes32 moduleId) external;

    /**
     * @notice Enable a registered module
     * @param moduleId Module identifier
     * @dev Validates dependencies before enabling
     */
    function enableModule(bytes32 moduleId) external;

    /**
     * @notice Disable a module
     * @param moduleId Module identifier
     * @dev Does not unregister, just marks as inactive
     */
    function disableModule(bytes32 moduleId) external;

    /**
     * @notice Upgrade a module to a new implementation
     * @param moduleId Module identifier
     * @param newImplementation Address of new implementation
     * @dev Validates that new implementation is compatible
     */
    function upgradeModule(bytes32 moduleId, address newImplementation) external;

    // ============ Module Discovery ============

    /**
     * @notice Get module address by ID
     * @param moduleId Module identifier
     * @return Module contract address
     */
    function getModule(bytes32 moduleId) external view returns (address);

    /**
     * @notice Get module registration details
     * @param moduleId Module identifier
     * @return Module registration data
     */
    function getModuleRegistration(bytes32 moduleId) external view returns (ModuleRegistration memory);

    /**
     * @notice Check if a module is registered
     * @param moduleId Module identifier
     * @return True if module is registered
     */
    function isModuleRegistered(bytes32 moduleId) external view returns (bool);

    /**
     * @notice Check if a module is enabled
     * @param moduleId Module identifier
     * @return True if module is enabled and active
     */
    function isModuleEnabled(bytes32 moduleId) external view returns (bool);

    /**
     * @notice Get all registered module IDs
     * @return Array of module identifiers
     */
    function getAllModuleIds() external view returns (bytes32[] memory);

    /**
     * @notice Get modules by state
     * @param state Module state to filter by
     * @return Array of module IDs in the specified state
     */
    function getModulesByState(IModule.ModuleState state) external view returns (bytes32[] memory);

    // ============ Dependency Management ============

    /**
     * @notice Validate module dependencies
     * @param moduleId Module identifier
     * @return True if all dependencies are satisfied
     */
    function validateDependencies(bytes32 moduleId) external view returns (bool);

    /**
     * @notice Get modules that depend on a specific module
     * @param moduleId Module identifier
     * @return Array of dependent module IDs
     */
    function getDependents(bytes32 moduleId) external view returns (bytes32[] memory);

    /**
     * @notice Check if upgrading/removing a module would break dependencies
     * @param moduleId Module identifier
     * @return safe True if operation is safe
     * @return affectedModules Array of module IDs that would be affected
     */
    function checkDependencySafety(bytes32 moduleId)
        external
        view
        returns (bool safe, bytes32[] memory affectedModules);

    // ============ System Health ============

    /**
     * @notice Perform health check on all active modules
     * @return healthyModules Number of healthy modules
     * @return totalModules Total number of active modules
     * @return unhealthyModuleIds Array of unhealthy module IDs
     */
    function systemHealthCheck()
        external
        view
        returns (
            uint256 healthyModules,
            uint256 totalModules,
            bytes32[] memory unhealthyModuleIds
        );

    /**
     * @notice Get registry statistics
     * @return totalRegistered Total registered modules
     * @return totalEnabled Total enabled modules
     * @return totalDeprecated Total deprecated modules
     */
    function getRegistryStats()
        external
        view
        returns (
            uint256 totalRegistered,
            uint256 totalEnabled,
            uint256 totalDeprecated
        );
}
