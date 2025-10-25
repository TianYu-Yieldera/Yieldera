// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../interfaces/IPlugin.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title BasePlugin
 * @notice Abstract base contract for all plugins
 * @dev Provides common plugin functionality and structure
 */
abstract contract BasePlugin is IPlugin, Ownable {
    // ============ State Variables ============

    PluginStatus private _status;
    address private _author;
    uint256 private _installedAt;
    uint256 private _lastUpdated;

    bool private _initialized;

    // ============ Constructor ============

    constructor(address author_) {
        _author = author_;
        _status = PluginStatus.INACTIVE;
        _installedAt = block.timestamp;
        _lastUpdated = block.timestamp;
    }

    // ============ Modifiers ============

    modifier whenActive() {
        require(_status == PluginStatus.ACTIVE, "Plugin not active");
        _;
    }

    modifier whenNotPaused() {
        require(_status != PluginStatus.PAUSED, "Plugin paused");
        _;
    }

    modifier onlyInitialized() {
        require(_initialized, "Plugin not initialized");
        _;
    }

    // ============ IPlugin Implementation ============

    function initialize(bytes calldata data) external override {
        require(!_initialized, "Already initialized");
        _initialized = true;
        _status = PluginStatus.ACTIVE;
        _onInitialize(data);
    }

    function execute(bytes calldata data)
        external
        override
        whenActive
        onlyInitialized
        returns (bytes memory result)
    {
        result = _execute(data);
        emit PluginExecuted(getPluginId(), msg.sender, data, result);
        return result;
    }

    function healthCheck()
        external
        view
        override
        returns (bool healthy, string memory message)
    {
        if (!_initialized) {
            return (false, "Not initialized");
        }
        if (_status != PluginStatus.ACTIVE) {
            return (false, "Not active");
        }
        return _healthCheck();
    }

    function getStatus() external view override returns (PluginStatus) {
        return _status;
    }

    function getPluginInfo() external view override returns (PluginInfo memory) {
        return PluginInfo({
            pluginId: getPluginId(),
            name: getPluginName(),
            version: getPluginVersion(),
            pluginType: getPluginType(),
            status: _status,
            implementation: address(this),
            installedAt: _installedAt,
            lastUpdated: _lastUpdated
        });
    }

    // ============ Lifecycle Functions ============

    function enable() external override onlyOwner {
        require(_initialized, "Not initialized");
        require(_status != PluginStatus.ACTIVE, "Already active");
        PluginStatus oldStatus = _status;
        _status = PluginStatus.ACTIVE;
        _lastUpdated = block.timestamp;
        emit PluginStatusChanged(getPluginId(), oldStatus, _status);
    }

    function disable() external override onlyOwner {
        require(_status != PluginStatus.INACTIVE, "Already inactive");
        PluginStatus oldStatus = _status;
        _status = PluginStatus.INACTIVE;
        _lastUpdated = block.timestamp;
        emit PluginStatusChanged(getPluginId(), oldStatus, _status);
    }

    function pause() external override onlyOwner {
        require(_status == PluginStatus.ACTIVE, "Not active");
        PluginStatus oldStatus = _status;
        _status = PluginStatus.PAUSED;
        _lastUpdated = block.timestamp;
        emit PluginStatusChanged(getPluginId(), oldStatus, _status);
    }

    function unpause() external override onlyOwner {
        require(_status == PluginStatus.PAUSED, "Not paused");
        PluginStatus oldStatus = _status;
        _status = PluginStatus.ACTIVE;
        _lastUpdated = block.timestamp;
        emit PluginStatusChanged(getPluginId(), oldStatus, _status);
    }

    function upgrade(address newImplementation) external override onlyOwner {
        require(newImplementation != address(0), "Invalid address");
        require(newImplementation != address(this), "Same implementation");
        // Actual upgrade logic would be implemented in proxy pattern
        _lastUpdated = block.timestamp;
    }

    // ============ Metadata Functions ============

    function getAuthor() external view override returns (address) {
        return _author;
    }

    function isCompatible(string calldata systemVersion)
        external
        pure
        virtual
        override
        returns (bool compatible)
    {
        // Default: compatible with all versions
        // Plugins can override for specific version requirements
        systemVersion; // Avoid unused parameter warning
        return true;
    }

    // ============ Abstract Functions (Must Override) ============

    /**
     * @notice Plugin-specific initialization logic
     * @param data Initialization data
     */
    function _onInitialize(bytes calldata data) internal virtual;

    /**
     * @notice Plugin-specific execution logic
     * @param data Execution data
     * @return result Execution result
     */
    function _execute(bytes calldata data) internal virtual returns (bytes memory result);

    /**
     * @notice Plugin-specific health check logic
     * @return healthy True if healthy
     * @return message Status message
     */
    function _healthCheck() internal view virtual returns (bool healthy, string memory message);

    // ============ Helper Functions ============

    /**
     * @notice Update last updated timestamp
     */
    function _updateTimestamp() internal {
        _lastUpdated = block.timestamp;
    }

    /**
     * @notice Check if plugin is initialized
     */
    function _isInitialized() internal view returns (bool) {
        return _initialized;
    }

    /**
     * @notice Get current status
     */
    function _getStatus() internal view returns (PluginStatus) {
        return _status;
    }
}
