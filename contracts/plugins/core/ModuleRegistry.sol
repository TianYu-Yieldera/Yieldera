// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/core/IModuleRegistry.sol";
import "../../interfaces/core/IModule.sol";
import "../../interfaces/core/IAccessController.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

/**
 * @title ModuleRegistry
 * @notice Central registry for managing all pluggable modules
 * @dev Implements IModuleRegistry with full dependency management
 */
contract ModuleRegistry is IModuleRegistry, Ownable, Pausable, ReentrancyGuard {
    // ============ State Variables ============

    // Access controller reference
    IAccessController public accessController;

    // Module registrations
    mapping(bytes32 => ModuleRegistration) private _modules;
    mapping(address => bytes32) private _moduleAddressToId;

    // All registered module IDs
    bytes32[] private _allModuleIds;

    // Module states
    mapping(bytes32 => IModule.ModuleState) private _moduleStates;

    // Dependency tracking
    mapping(bytes32 => bytes32[]) private _moduleDependents; // moduleId => dependent module IDs

    // Statistics
    uint256 private _totalRegistered;
    uint256 private _totalEnabled;
    uint256 private _totalDeprecated;

    // ============ Modifiers ============

    modifier onlyRegistered(bytes32 moduleId) {
        require(_modules[moduleId].isRegistered, "Module not registered");
        _;
    }

    modifier onlyValidModule(address moduleAddress) {
        require(moduleAddress != address(0), "Invalid module address");
        require(_supportsInterface(moduleAddress), "Does not implement IModule");
        _;
    }

    // ============ Constructor ============

    constructor() {
        // Owner is set by Ownable constructor
    }

    // ============ Module Management Functions ============

    /**
     * @notice Register a new module
     * @param moduleAddress Address of the module contract
     * @return moduleId The assigned module ID
     */
    function registerModule(address moduleAddress)
        external
        override
        onlyOwner
        whenNotPaused
        nonReentrant
        onlyValidModule(moduleAddress)
        returns (bytes32 moduleId)
    {
        IModule module = IModule(moduleAddress);

        // Get module ID from the module itself
        moduleId = module.getModuleId();

        // Check if already registered
        require(!_modules[moduleId].isRegistered, "Module already registered");
        require(_moduleAddressToId[moduleAddress] == bytes32(0), "Address already registered");

        // Get module info
        IModule.ModuleInfo memory info = module.getModuleInfo();
        bytes32[] memory dependencies = module.getDependencies();

        // Validate dependencies exist
        for (uint256 i = 0; i < dependencies.length; i++) {
            require(
                _modules[dependencies[i]].isRegistered,
                "Dependency not registered"
            );
            // Track this module as a dependent
            _moduleDependents[dependencies[i]].push(moduleId);
        }

        // Create registration
        _modules[moduleId] = ModuleRegistration({
            moduleAddress: moduleAddress,
            moduleId: moduleId,
            name: info.name,
            version: info.version,
            dependencies: dependencies,
            isRegistered: true,
            registeredAt: block.timestamp
        });

        _moduleAddressToId[moduleAddress] = moduleId;
        _allModuleIds.push(moduleId);
        _moduleStates[moduleId] = IModule.ModuleState.UNINITIALIZED;
        _totalRegistered++;

        emit ModuleRegistered(moduleId, moduleAddress, info.name, info.version);

        return moduleId;
    }

    /**
     * @notice Unregister a module
     * @param moduleId Module identifier
     */
    function unregisterModule(bytes32 moduleId)
        external
        override
        onlyOwner
        onlyRegistered(moduleId)
        nonReentrant
    {
        // Check if other modules depend on this one
        require(
            _moduleDependents[moduleId].length == 0,
            "Cannot unregister: other modules depend on it"
        );

        ModuleRegistration storage module = _modules[moduleId];
        address moduleAddress = module.moduleAddress;

        // Remove from dependencies of other modules
        bytes32[] memory deps = module.dependencies;
        for (uint256 i = 0; i < deps.length; i++) {
            _removeDependency(deps[i], moduleId);
        }

        // Clean up
        delete _moduleAddressToId[moduleAddress];
        delete _modules[moduleId];
        delete _moduleStates[moduleId];

        // Remove from all module IDs array
        _removeFromArray(_allModuleIds, moduleId);

        _totalRegistered--;

        emit ModuleUnregistered(moduleId, moduleAddress);
    }

    /**
     * @notice Enable a registered module
     * @param moduleId Module identifier
     */
    function enableModule(bytes32 moduleId)
        external
        override
        onlyOwner
        onlyRegistered(moduleId)
    {
        require(
            _moduleStates[moduleId] == IModule.ModuleState.UNINITIALIZED ||
            _moduleStates[moduleId] == IModule.ModuleState.PAUSED,
            "Invalid state for enabling"
        );

        // Validate dependencies are enabled
        require(validateDependencies(moduleId), "Dependencies not satisfied");

        address moduleAddress = _modules[moduleId].moduleAddress;
        _moduleStates[moduleId] = IModule.ModuleState.ACTIVE;
        _totalEnabled++;

        emit ModuleEnabled(moduleId, moduleAddress);
    }

    /**
     * @notice Disable a module
     * @param moduleId Module identifier
     */
    function disableModule(bytes32 moduleId)
        external
        override
        onlyOwner
        onlyRegistered(moduleId)
    {
        require(
            _moduleStates[moduleId] == IModule.ModuleState.ACTIVE,
            "Module not active"
        );

        address moduleAddress = _modules[moduleId].moduleAddress;
        _moduleStates[moduleId] = IModule.ModuleState.PAUSED;
        _totalEnabled--;

        emit ModuleDisabled(moduleId, moduleAddress);
    }

    /**
     * @notice Upgrade a module to a new implementation
     * @param moduleId Module identifier
     * @param newImplementation Address of new implementation
     */
    function upgradeModule(bytes32 moduleId, address newImplementation)
        external
        override
        onlyOwner
        onlyRegistered(moduleId)
        onlyValidModule(newImplementation)
        nonReentrant
    {
        ModuleRegistration storage module = _modules[moduleId];
        address oldAddress = module.moduleAddress;

        // Verify new implementation has same module ID
        IModule newModule = IModule(newImplementation);
        require(newModule.getModuleId() == moduleId, "Module ID mismatch");

        // Get new version
        IModule.ModuleInfo memory newInfo = newModule.getModuleInfo();

        // Update registration
        module.moduleAddress = newImplementation;
        module.version = newInfo.version;

        // Update address mapping
        delete _moduleAddressToId[oldAddress];
        _moduleAddressToId[newImplementation] = moduleId;

        // Mark old state as upgraded
        _moduleStates[moduleId] = IModule.ModuleState.UPGRADED;

        emit ModuleUpgraded(moduleId, oldAddress, newImplementation, newInfo.version);
    }

    // ============ Module Discovery Functions ============

    /**
     * @notice Get module address by ID
     * @param moduleId Module identifier
     * @return Module contract address
     */
    function getModule(bytes32 moduleId)
        external
        view
        override
        onlyRegistered(moduleId)
        returns (address)
    {
        return _modules[moduleId].moduleAddress;
    }

    /**
     * @notice Get module registration details
     * @param moduleId Module identifier
     * @return Module registration data
     */
    function getModuleRegistration(bytes32 moduleId)
        external
        view
        override
        onlyRegistered(moduleId)
        returns (ModuleRegistration memory)
    {
        return _modules[moduleId];
    }

    /**
     * @notice Check if a module is registered
     * @param moduleId Module identifier
     * @return True if module is registered
     */
    function isModuleRegistered(bytes32 moduleId)
        external
        view
        override
        returns (bool)
    {
        return _modules[moduleId].isRegistered;
    }

    /**
     * @notice Check if a module is enabled
     * @param moduleId Module identifier
     * @return True if module is enabled and active
     */
    function isModuleEnabled(bytes32 moduleId)
        external
        view
        override
        returns (bool)
    {
        return _moduleStates[moduleId] == IModule.ModuleState.ACTIVE;
    }

    /**
     * @notice Get all registered module IDs
     * @return Array of module identifiers
     */
    function getAllModuleIds()
        external
        view
        override
        returns (bytes32[] memory)
    {
        return _allModuleIds;
    }

    /**
     * @notice Get modules by state
     * @param state Module state to filter by
     * @return Array of module IDs in the specified state
     */
    function getModulesByState(IModule.ModuleState state)
        external
        view
        override
        returns (bytes32[] memory)
    {
        // Count modules in state
        uint256 count = 0;
        for (uint256 i = 0; i < _allModuleIds.length; i++) {
            if (_moduleStates[_allModuleIds[i]] == state) {
                count++;
            }
        }

        // Collect module IDs
        bytes32[] memory result = new bytes32[](count);
        uint256 index = 0;
        for (uint256 i = 0; i < _allModuleIds.length; i++) {
            if (_moduleStates[_allModuleIds[i]] == state) {
                result[index] = _allModuleIds[i];
                index++;
            }
        }

        return result;
    }

    // ============ Dependency Management Functions ============

    /**
     * @notice Validate module dependencies
     * @param moduleId Module identifier
     * @return True if all dependencies are satisfied
     */
    function validateDependencies(bytes32 moduleId)
        public
        view
        override
        onlyRegistered(moduleId)
        returns (bool)
    {
        bytes32[] memory deps = _modules[moduleId].dependencies;

        for (uint256 i = 0; i < deps.length; i++) {
            bytes32 depId = deps[i];

            // Check if dependency is registered
            if (!_modules[depId].isRegistered) {
                return false;
            }

            // Check if dependency is active
            if (_moduleStates[depId] != IModule.ModuleState.ACTIVE) {
                return false;
            }

            emit DependencyValidated(moduleId, depId, true);
        }

        return true;
    }

    /**
     * @notice Get modules that depend on a specific module
     * @param moduleId Module identifier
     * @return Array of dependent module IDs
     */
    function getDependents(bytes32 moduleId)
        external
        view
        override
        onlyRegistered(moduleId)
        returns (bytes32[] memory)
    {
        return _moduleDependents[moduleId];
    }

    /**
     * @notice Check if upgrading/removing a module would break dependencies
     * @param moduleId Module identifier
     * @return safe True if operation is safe
     * @return affectedModules Array of module IDs that would be affected
     */
    function checkDependencySafety(bytes32 moduleId)
        external
        view
        override
        onlyRegistered(moduleId)
        returns (bool safe, bytes32[] memory affectedModules)
    {
        affectedModules = _moduleDependents[moduleId];
        safe = affectedModules.length == 0;
        return (safe, affectedModules);
    }

    // ============ System Health Functions ============

    /**
     * @notice Perform health check on all active modules
     * @return healthyModules Number of healthy modules
     * @return totalModules Total number of active modules
     * @return unhealthyModuleIds Array of unhealthy module IDs
     */
    function systemHealthCheck()
        external
        view
        override
        returns (
            uint256 healthyModules,
            uint256 totalModules,
            bytes32[] memory unhealthyModuleIds
        )
    {
        // Count active modules
        totalModules = 0;
        for (uint256 i = 0; i < _allModuleIds.length; i++) {
            if (_moduleStates[_allModuleIds[i]] == IModule.ModuleState.ACTIVE) {
                totalModules++;
            }
        }

        // Temporary array to store unhealthy module IDs
        bytes32[] memory tempUnhealthy = new bytes32[](totalModules);
        uint256 unhealthyCount = 0;

        // Check health of each active module
        for (uint256 i = 0; i < _allModuleIds.length; i++) {
            bytes32 moduleId = _allModuleIds[i];

            if (_moduleStates[moduleId] == IModule.ModuleState.ACTIVE) {
                IModule module = IModule(_modules[moduleId].moduleAddress);
                (bool healthy, ) = module.healthCheck();

                if (healthy) {
                    healthyModules++;
                } else {
                    tempUnhealthy[unhealthyCount] = moduleId;
                    unhealthyCount++;
                }
            }
        }

        // Create properly sized array for unhealthy modules
        unhealthyModuleIds = new bytes32[](unhealthyCount);
        for (uint256 i = 0; i < unhealthyCount; i++) {
            unhealthyModuleIds[i] = tempUnhealthy[i];
        }

        return (healthyModules, totalModules, unhealthyModuleIds);
    }

    /**
     * @notice Get registry statistics
     * @return totalRegistered Total registered modules
     * @return totalEnabled Total enabled modules
     * @return totalDeprecated Total deprecated modules
     */
    function getRegistryStats()
        external
        view
        override
        returns (
            uint256 totalRegistered,
            uint256 totalEnabled,
            uint256 totalDeprecated
        )
    {
        return (_totalRegistered, _totalEnabled, _totalDeprecated);
    }

    // ============ Admin Functions ============

    /**
     * @notice Set access controller
     * @param _accessController Access controller address
     */
    function setAccessController(address _accessController) external onlyOwner {
        require(_accessController != address(0), "Invalid access controller");
        accessController = IAccessController(_accessController);
    }

    /**
     * @notice Pause registry operations
     */
    function pause() external onlyOwner {
        _pause();
    }

    /**
     * @notice Unpause registry operations
     */
    function unpause() external onlyOwner {
        _unpause();
    }

    // ============ Internal Helper Functions ============

    /**
     * @notice Check if address implements IModule interface
     */
    function _supportsInterface(address moduleAddress) private view returns (bool) {
        try IModule(moduleAddress).getModuleId() returns (bytes32) {
            return true;
        } catch {
            return false;
        }
    }

    /**
     * @notice Remove module from dependency list
     */
    function _removeDependency(bytes32 dependencyId, bytes32 moduleId) private {
        bytes32[] storage dependents = _moduleDependents[dependencyId];
        for (uint256 i = 0; i < dependents.length; i++) {
            if (dependents[i] == moduleId) {
                dependents[i] = dependents[dependents.length - 1];
                dependents.pop();
                break;
            }
        }
    }

    /**
     * @notice Remove element from array
     */
    function _removeFromArray(bytes32[] storage array, bytes32 element) private {
        for (uint256 i = 0; i < array.length; i++) {
            if (array[i] == element) {
                array[i] = array[array.length - 1];
                array.pop();
                break;
            }
        }
    }
}
