// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../interfaces/IPluginRegistry.sol";
import "../interfaces/IPlugin.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title PluginRegistry
 * @notice Central registry for all plugins
 * @dev Manages plugin registration, discovery, and statistics
 */
contract PluginRegistry is IPluginRegistry, Ownable {
    // ============ Storage ============

    mapping(bytes32 => PluginRegistration) private _plugins;
    mapping(bytes32 => PluginStats) private _pluginStats;
    mapping(address => bytes32[]) private _pluginsByAuthor;
    mapping(IPlugin.PluginType => bytes32[]) private _pluginsByType;

    bytes32[] private _allPlugins;
    bytes32[] private _activePlugins;
    bytes32[] private _verifiedPlugins;

    mapping(address => bool) private _verifiers;

    uint256 private _totalPluginCount;
    uint256 private _activePluginCount;
    uint256 private _verifiedPluginCount;

    // ============ Constructor ============

    constructor() {
        _verifiers[msg.sender] = true; // Owner is default verifier
    }

    // ============ Admin Functions ============

    function addVerifier(address verifier) external onlyOwner {
        require(verifier != address(0), "Invalid verifier");
        _verifiers[verifier] = true;
    }

    function removeVerifier(address verifier) external onlyOwner {
        _verifiers[verifier] = false;
    }

    modifier onlyVerifier() {
        require(_verifiers[msg.sender], "Not a verifier");
        _;
    }

    // ============ Registration Functions ============

    function registerPlugin(address plugin) external override returns (bytes32 pluginId) {
        require(plugin != address(0), "Invalid plugin address");
        require(plugin.code.length > 0, "Not a contract");

        IPlugin pluginContract = IPlugin(plugin);

        // Get plugin info
        pluginId = pluginContract.getPluginId();
        require(!_plugins[pluginId].isActive, "Plugin already registered");

        string memory version = pluginContract.getPluginVersion();
        IPlugin.PluginType pluginType = pluginContract.getPluginType();
        address author = pluginContract.getAuthor();

        // Create registration
        _plugins[pluginId] = PluginRegistration({
            pluginId: pluginId,
            implementation: plugin,
            pluginType: pluginType,
            version: version,
            author: author,
            registeredAt: block.timestamp,
            isVerified: false,
            isActive: true
        });

        // Initialize stats
        _pluginStats[pluginId] = PluginStats({
            totalInstalls: 1,
            activeInstalls: 1,
            totalExecutions: 0,
            failedExecutions: 0,
            lastExecuted: 0
        });

        // Add to indices
        _allPlugins.push(pluginId);
        _activePlugins.push(pluginId);
        _pluginsByAuthor[author].push(pluginId);
        _pluginsByType[pluginType].push(pluginId);

        // Update counters
        _totalPluginCount++;
        _activePluginCount++;

        emit PluginRegistered(pluginId, plugin, author, version);

        return pluginId;
    }

    function unregisterPlugin(bytes32 pluginId) external override {
        require(_plugins[pluginId].isActive, "Plugin not registered");
        require(
            _plugins[pluginId].author == msg.sender || msg.sender == owner(),
            "Not authorized"
        );

        _plugins[pluginId].isActive = false;

        // Remove from active plugins
        _removeFromArray(_activePlugins, pluginId);

        // Remove from verified if present
        if (_plugins[pluginId].isVerified) {
            _removeFromArray(_verifiedPlugins, pluginId);
            _verifiedPluginCount--;
        }

        _activePluginCount--;

        emit PluginUnregistered(pluginId, _plugins[pluginId].implementation);
    }

    function verifyPlugin(bytes32 pluginId) external override onlyVerifier {
        require(_plugins[pluginId].isActive, "Plugin not registered");
        require(!_plugins[pluginId].isVerified, "Already verified");

        _plugins[pluginId].isVerified = true;
        _verifiedPlugins.push(pluginId);
        _verifiedPluginCount++;

        emit PluginVerified(pluginId, msg.sender);
    }

    function activatePlugin(bytes32 pluginId) external override {
        require(_plugins[pluginId].pluginId != bytes32(0), "Plugin not registered");
        require(!_plugins[pluginId].isActive, "Already active");
        require(
            _plugins[pluginId].author == msg.sender || msg.sender == owner(),
            "Not authorized"
        );

        _plugins[pluginId].isActive = true;
        _activePlugins.push(pluginId);
        _activePluginCount++;

        emit PluginActivated(pluginId);
    }

    function deactivatePlugin(bytes32 pluginId) external override {
        require(_plugins[pluginId].isActive, "Not active");
        require(
            _plugins[pluginId].author == msg.sender || msg.sender == owner(),
            "Not authorized"
        );

        _plugins[pluginId].isActive = false;
        _removeFromArray(_activePlugins, pluginId);
        _activePluginCount--;

        emit PluginDeactivated(pluginId);
    }

    // ============ Query Functions ============

    function isPluginRegistered(bytes32 pluginId) external view override returns (bool) {
        return _plugins[pluginId].pluginId != bytes32(0);
    }

    function isPluginActive(bytes32 pluginId) external view override returns (bool) {
        return _plugins[pluginId].isActive;
    }

    function isPluginVerified(bytes32 pluginId) external view override returns (bool) {
        return _plugins[pluginId].isVerified;
    }

    function getPluginRegistration(bytes32 pluginId)
        external
        view
        override
        returns (PluginRegistration memory registration)
    {
        return _plugins[pluginId];
    }

    function getPluginImplementation(bytes32 pluginId) external view override returns (address) {
        return _plugins[pluginId].implementation;
    }

    function getPluginStats(bytes32 pluginId)
        external
        view
        override
        returns (PluginStats memory stats)
    {
        return _pluginStats[pluginId];
    }

    // ============ Discovery Functions ============

    function getAllPlugins() external view override returns (bytes32[] memory) {
        return _allPlugins;
    }

    function getPluginsByType(IPlugin.PluginType pluginType)
        external
        view
        override
        returns (bytes32[] memory)
    {
        return _pluginsByType[pluginType];
    }

    function getPluginsByAuthor(address author) external view override returns (bytes32[] memory) {
        return _pluginsByAuthor[author];
    }

    function getActivePlugins() external view override returns (bytes32[] memory) {
        return _activePlugins;
    }

    function getVerifiedPlugins() external view override returns (bytes32[] memory) {
        return _verifiedPlugins;
    }

    function searchPluginsByName(string calldata namePattern)
        external
        view
        override
        returns (bytes32[] memory)
    {
        // Simple implementation - returns all plugins
        // In production, would implement actual pattern matching
        return _allPlugins;
    }

    // ============ Statistics Functions ============

    function recordExecution(bytes32 pluginId, bool success) external override {
        require(_plugins[pluginId].isActive, "Plugin not active");

        _pluginStats[pluginId].totalExecutions++;
        if (!success) {
            _pluginStats[pluginId].failedExecutions++;
        }
        _pluginStats[pluginId].lastExecuted = block.timestamp;

        emit PluginExecutionRecorded(pluginId, success);
    }

    function getTotalPluginCount() external view override returns (uint256) {
        return _totalPluginCount;
    }

    function getActivePluginCount() external view override returns (uint256) {
        return _activePluginCount;
    }

    function getVerifiedPluginCount() external view override returns (uint256) {
        return _verifiedPluginCount;
    }

    // ============ Compatibility Functions ============

    function checkCompatibility(bytes32 pluginId, string calldata systemVersion)
        external
        view
        override
        returns (bool compatible)
    {
        address implementation = _plugins[pluginId].implementation;
        if (implementation == address(0)) return false;

        IPlugin plugin = IPlugin(implementation);
        return plugin.isCompatible(systemVersion);
    }

    function checkDependencies(bytes32 pluginId)
        external
        view
        override
        returns (bool satisfied, bytes32[] memory missingDeps)
    {
        address implementation = _plugins[pluginId].implementation;
        require(implementation != address(0), "Plugin not registered");

        IPlugin plugin = IPlugin(implementation);
        bytes32[] memory dependencies = plugin.getDependencies();

        // Count missing dependencies
        uint256 missingCount = 0;
        for (uint256 i = 0; i < dependencies.length; i++) {
            if (!_plugins[dependencies[i]].isActive) {
                missingCount++;
            }
        }

        if (missingCount == 0) {
            return (true, new bytes32[](0));
        }

        // Populate missing dependencies
        missingDeps = new bytes32[](missingCount);
        uint256 index = 0;
        for (uint256 i = 0; i < dependencies.length; i++) {
            if (!_plugins[dependencies[i]].isActive) {
                missingDeps[index++] = dependencies[i];
            }
        }

        return (false, missingDeps);
    }

    // ============ Internal Helper Functions ============

    function _removeFromArray(bytes32[] storage array, bytes32 value) internal {
        for (uint256 i = 0; i < array.length; i++) {
            if (array[i] == value) {
                array[i] = array[array.length - 1];
                array.pop();
                break;
            }
        }
    }
}
