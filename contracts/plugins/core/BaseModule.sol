// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/core/IModule.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

/**
 * @title BaseModule
 * @notice Abstract base contract for all pluggable modules
 * @dev Provides common functionality for module lifecycle management
 */
abstract contract BaseModule is IModule, Pausable {
    // ============ State Variables ============

    bytes32 private immutable _moduleId;
    string private _moduleName;
    string private _moduleVersion;
    ModuleState private _moduleState;
    uint256 private _installedAt;
    uint256 private _lastUpdated;
    bool private _initialized;

    // ============ Constructor ============

    constructor(bytes32 moduleId, string memory name, string memory version) {
        _moduleId = moduleId;
        _moduleName = name;
        _moduleVersion = version;
        _moduleState = ModuleState.UNINITIALIZED;
    }

    // ============ IModule Implementation ============

    function getModuleId() external view virtual override returns (bytes32) {
        return _moduleId;
    }

    function getModuleInfo() external view virtual override returns (ModuleInfo memory) {
        return ModuleInfo({
            moduleId: _moduleId,
            name: _moduleName,
            version: _moduleVersion,
            implementation: address(this),
            state: _moduleState,
            installedAt: _installedAt,
            lastUpdated: _lastUpdated
        });
    }

    function getVersion() external view virtual override returns (string memory) {
        return _moduleVersion;
    }

    function getDependencies() external view virtual override returns (bytes32[] memory) {
        // Default: no dependencies
        // Override in derived contracts to specify dependencies
        return new bytes32[](0);
    }

    function isActive() external view virtual override returns (bool) {
        return _moduleState == ModuleState.ACTIVE && !paused();
    }

    function initialize(bytes calldata data) external virtual override {
        require(!_initialized, "Already initialized");
        require(_moduleState == ModuleState.UNINITIALIZED, "Invalid state");

        _initialized = true;
        _moduleState = ModuleState.ACTIVE;
        _installedAt = block.timestamp;
        _lastUpdated = block.timestamp;

        _initializeModule(data);

        emit ModuleInitialized(_moduleId, address(this), _moduleVersion);
    }

    function pause() external virtual override {
        _requireNotPaused();
        _pause();
        _updateState(ModuleState.PAUSED);
    }

    function unpause() external virtual override {
        _requirePaused();
        _unpause();
        _updateState(ModuleState.ACTIVE);
    }

    function healthCheck() external view virtual override returns (bool healthy, string memory message) {
        // Default health check
        if (!_initialized) {
            return (false, "Module not initialized");
        }

        if (_moduleState != ModuleState.ACTIVE) {
            return (false, "Module not active");
        }

        if (paused()) {
            return (false, "Module paused");
        }

        return (true, "Module healthy");
    }

    // ============ Internal Functions ============

    /**
     * @notice Override this function to implement module-specific initialization
     * @param data Initialization data
     */
    function _initializeModule(bytes calldata data) internal virtual {
        // Override in derived contracts
    }

    /**
     * @notice Update module state
     */
    function _updateState(ModuleState newState) internal {
        ModuleState oldState = _moduleState;
        _moduleState = newState;
        _lastUpdated = block.timestamp;

        emit ModuleStateChanged(_moduleId, oldState, newState);
    }

    /**
     * @notice Mark module as deprecated
     */
    function _deprecate() internal {
        _updateState(ModuleState.DEPRECATED);
    }

    /**
     * @notice Mark module as upgraded
     */
    function _markUpgraded() internal {
        _updateState(ModuleState.UPGRADED);
    }

    /**
     * @notice Require module to be initialized
     */
    function _requireInitialized() internal view {
        require(_initialized, "Module not initialized");
    }

    /**
     * @notice Require module to be active
     */
    function _requireActive() internal view {
        require(_moduleState == ModuleState.ACTIVE, "Module not active");
    }
}
