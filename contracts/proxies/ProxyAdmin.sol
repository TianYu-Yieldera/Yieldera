// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "../interfaces/IUpgradeable.sol";

/**
 * @title ProxyAdmin
 * @notice Centralized admin for managing multiple module proxies
 * @dev Manages upgrade permissions and coordinates multi-proxy upgrades
 *
 * Features:
 * - Manage multiple proxies from single contract
 * - Batch upgrades
 * - Timelock support
 * - Emergency pause
 */
contract ProxyAdmin is Ownable {
    using Address for address;

    // ============ Events ============

    event ProxyRegistered(address indexed proxy, string moduleName);
    event ProxyUnregistered(address indexed proxy);
    event UpgradeScheduled(
        address indexed proxy,
        address indexed newImplementation,
        uint256 executeAfter
    );
    event UpgradeExecuted(address indexed proxy, address indexed newImplementation);
    event UpgradeCancelled(address indexed proxy, address indexed newImplementation);
    event TimelockUpdated(uint256 oldTimelock, uint256 newTimelock);
    event EmergencyPaused(address indexed by);
    event EmergencyUnpaused(address indexed by);

    // ============ Structures ============

    struct ProxyInfo {
        string moduleName;
        address currentImplementation;
        bool isRegistered;
        uint256 lastUpgrade;
    }

    struct ScheduledUpgrade {
        address newImplementation;
        uint256 executeAfter;
        bool executed;
        bool cancelled;
    }

    // ============ State Variables ============

    // Registered proxies
    mapping(address => ProxyInfo) public proxies;
    address[] public proxyList;

    // Scheduled upgrades
    mapping(address => ScheduledUpgrade) public scheduledUpgrades;

    // Timelock duration (default: 2 days)
    uint256 public timelockDuration = 2 days;

    // Emergency pause
    bool public paused;

    // Upgrade authorizers
    mapping(address => bool) public isUpgradeAuthorizer;

    // ============ Modifiers ============

    modifier onlyAuthorizer() {
        require(
            isUpgradeAuthorizer[msg.sender] || msg.sender == owner(),
            "Not authorized"
        );
        _;
    }

    modifier whenNotPaused() {
        require(!paused, "Paused");
        _;
    }

    modifier onlyRegistered(address proxy) {
        require(proxies[proxy].isRegistered, "Proxy not registered");
        _;
    }

    // ============ Constructor ============

    constructor() {
        isUpgradeAuthorizer[msg.sender] = true;
    }

    // ============ Proxy Management ============

    /**
     * @notice Register a new proxy
     * @param proxy Proxy address
     * @param moduleName Module name
     */
    function registerProxy(address proxy, string memory moduleName) external onlyOwner {
        require(!proxies[proxy].isRegistered, "Already registered");
        require(proxy.isContract(), "Not a contract");

        proxies[proxy] = ProxyInfo({
            moduleName: moduleName,
            currentImplementation: address(0), // Will be set on first upgrade
            isRegistered: true,
            lastUpgrade: block.timestamp
        });

        proxyList.push(proxy);

        emit ProxyRegistered(proxy, moduleName);
    }

    /**
     * @notice Unregister a proxy
     * @param proxy Proxy address
     */
    function unregisterProxy(address proxy) external onlyOwner onlyRegistered(proxy) {
        proxies[proxy].isRegistered = false;

        emit ProxyUnregistered(proxy);
    }

    /**
     * @notice Get all registered proxies
     * @return Array of proxy addresses
     */
    function getProxies() external view returns (address[] memory) {
        return proxyList;
    }

    /**
     * @notice Get proxy information
     * @param proxy Proxy address
     * @return info Proxy information
     */
    function getProxyInfo(address proxy) external view returns (ProxyInfo memory) {
        return proxies[proxy];
    }

    // ============ Upgrade Management ============

    /**
     * @notice Schedule an upgrade
     * @param proxy Proxy address
     * @param newImplementation New implementation address
     */
    function scheduleUpgrade(address proxy, address newImplementation)
        external
        onlyAuthorizer
        whenNotPaused
        onlyRegistered(proxy)
    {
        require(newImplementation.isContract(), "Not a contract");
        require(
            scheduledUpgrades[proxy].executeAfter == 0 ||
            scheduledUpgrades[proxy].executed ||
            scheduledUpgrades[proxy].cancelled,
            "Upgrade already scheduled"
        );

        uint256 executeAfter = block.timestamp + timelockDuration;

        scheduledUpgrades[proxy] = ScheduledUpgrade({
            newImplementation: newImplementation,
            executeAfter: executeAfter,
            executed: false,
            cancelled: false
        });

        emit UpgradeScheduled(proxy, newImplementation, executeAfter);
    }

    /**
     * @notice Execute a scheduled upgrade
     * @param proxy Proxy address
     */
    function executeUpgrade(address proxy)
        external
        onlyAuthorizer
        whenNotPaused
        onlyRegistered(proxy)
    {
        ScheduledUpgrade storage upgrade = scheduledUpgrades[proxy];

        require(upgrade.executeAfter > 0, "No upgrade scheduled");
        require(!upgrade.executed, "Already executed");
        require(!upgrade.cancelled, "Upgrade cancelled");
        require(block.timestamp >= upgrade.executeAfter, "Timelock not expired");

        // Perform the upgrade
        _upgradeProxy(proxy, upgrade.newImplementation);

        // Mark as executed
        upgrade.executed = true;

        // Update proxy info
        proxies[proxy].currentImplementation = upgrade.newImplementation;
        proxies[proxy].lastUpgrade = block.timestamp;

        emit UpgradeExecuted(proxy, upgrade.newImplementation);
    }

    /**
     * @notice Cancel a scheduled upgrade
     * @param proxy Proxy address
     */
    function cancelUpgrade(address proxy) external onlyOwner onlyRegistered(proxy) {
        ScheduledUpgrade storage upgrade = scheduledUpgrades[proxy];

        require(upgrade.executeAfter > 0, "No upgrade scheduled");
        require(!upgrade.executed, "Already executed");
        require(!upgrade.cancelled, "Already cancelled");

        upgrade.cancelled = true;

        emit UpgradeCancelled(proxy, upgrade.newImplementation);
    }

    /**
     * @notice Immediate upgrade (bypasses timelock) - Emergency only
     * @param proxy Proxy address
     * @param newImplementation New implementation address
     * @dev Only owner can call, for emergency situations
     */
    function emergencyUpgrade(address proxy, address newImplementation)
        external
        onlyOwner
        onlyRegistered(proxy)
    {
        require(newImplementation.isContract(), "Not a contract");

        _upgradeProxy(proxy, newImplementation);

        proxies[proxy].currentImplementation = newImplementation;
        proxies[proxy].lastUpgrade = block.timestamp;

        emit UpgradeExecuted(proxy, newImplementation);
    }

    /**
     * @notice Batch upgrade multiple proxies
     * @param proxyAddresses Array of proxy addresses
     * @param implementations Array of implementation addresses
     */
    function batchScheduleUpgrade(
        address[] calldata proxyAddresses,
        address[] calldata implementations
    ) external onlyAuthorizer whenNotPaused {
        require(proxyAddresses.length == implementations.length, "Length mismatch");

        for (uint256 i = 0; i < proxyAddresses.length; i++) {
            require(proxies[proxyAddresses[i]].isRegistered, "Proxy not registered");
            require(implementations[i].isContract(), "Not a contract");

            uint256 executeAfter = block.timestamp + timelockDuration;

            scheduledUpgrades[proxyAddresses[i]] = ScheduledUpgrade({
                newImplementation: implementations[i],
                executeAfter: executeAfter,
                executed: false,
                cancelled: false
            });

            emit UpgradeScheduled(proxyAddresses[i], implementations[i], executeAfter);
        }
    }

    // ============ Internal Functions ============

    /**
     * @notice Internal function to upgrade a proxy
     * @param proxy Proxy address
     * @param newImplementation New implementation address
     */
    function _upgradeProxy(address proxy, address newImplementation) internal {
        // Call upgradeTo on the proxy
        // This will call the upgradeTo function in the UUPS implementation
        (bool success, ) = proxy.call(
            abi.encodeWithSignature("upgradeTo(address)", newImplementation)
        );
        require(success, "Upgrade failed");
    }

    // ============ Configuration ============

    /**
     * @notice Update timelock duration
     * @param newDuration New duration in seconds
     */
    function setTimelockDuration(uint256 newDuration) external onlyOwner {
        require(newDuration >= 1 hours, "Duration too short");
        require(newDuration <= 30 days, "Duration too long");

        uint256 oldDuration = timelockDuration;
        timelockDuration = newDuration;

        emit TimelockUpdated(oldDuration, newDuration);
    }

    /**
     * @notice Add upgrade authorizer
     * @param authorizer Address to authorize
     */
    function addAuthorizer(address authorizer) external onlyOwner {
        isUpgradeAuthorizer[authorizer] = true;
    }

    /**
     * @notice Remove upgrade authorizer
     * @param authorizer Address to revoke
     */
    function removeAuthorizer(address authorizer) external onlyOwner {
        isUpgradeAuthorizer[authorizer] = false;
    }

    // ============ Emergency Functions ============

    /**
     * @notice Emergency pause all upgrades
     */
    function pause() external onlyOwner {
        paused = true;
        emit EmergencyPaused(msg.sender);
    }

    /**
     * @notice Resume upgrades
     */
    function unpause() external onlyOwner {
        paused = false;
        emit EmergencyUnpaused(msg.sender);
    }

    // ============ View Functions ============

    /**
     * @notice Check if an upgrade can be executed
     * @param proxy Proxy address
     * @return canExecute True if upgrade can be executed
     * @return reason Reason if cannot execute
     */
    function canExecuteUpgrade(address proxy)
        external
        view
        returns (bool canExecute, string memory reason)
    {
        if (!proxies[proxy].isRegistered) {
            return (false, "Proxy not registered");
        }

        ScheduledUpgrade storage upgrade = scheduledUpgrades[proxy];

        if (upgrade.executeAfter == 0) {
            return (false, "No upgrade scheduled");
        }

        if (upgrade.executed) {
            return (false, "Already executed");
        }

        if (upgrade.cancelled) {
            return (false, "Upgrade cancelled");
        }

        if (block.timestamp < upgrade.executeAfter) {
            return (false, "Timelock not expired");
        }

        if (paused) {
            return (false, "System paused");
        }

        return (true, "");
    }

    /**
     * @notice Get time remaining until upgrade can be executed
     * @param proxy Proxy address
     * @return timeRemaining Seconds until execution (0 if ready)
     */
    function getUpgradeTimeRemaining(address proxy) external view returns (uint256) {
        ScheduledUpgrade storage upgrade = scheduledUpgrades[proxy];

        if (upgrade.executeAfter == 0 || block.timestamp >= upgrade.executeAfter) {
            return 0;
        }

        return upgrade.executeAfter - block.timestamp;
    }
}
