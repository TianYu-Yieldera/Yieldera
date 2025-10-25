// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/core/IAuditModule.sol";
import "../../plugins/core/BaseModule.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

/**
 * @title AuditModule
 * @notice Pluggable audit and logging module for compliance tracking
 * @dev Extends BaseModule and implements IAuditModule
 */
contract AuditModule is IAuditModule, BaseModule, AccessControl {
    // ============ Constants ============

    bytes32 public constant MODULE_ID = keccak256("AUDIT_MODULE");
    string public constant MODULE_NAME = "AuditModule";
    string public constant MODULE_VERSION = "1.0.0";

    bytes32 public constant LOGGER_ROLE = keccak256("LOGGER_ROLE");
    bytes32 public constant AUDITOR_ROLE = keccak256("AUDITOR_ROLE");

    // ============ State Variables ============

    // Audit log storage
    AuditEntry[] private _auditLogs;
    mapping(address => uint256[]) private _userLogIndices;
    mapping(AuditEventType => uint256[]) private _eventTypeIndices;
    mapping(bytes32 => uint256[]) private _moduleLogIndices;

    // Configuration
    uint256 public maxActiveLogs = 100000;
    uint256 public retentionPeriod = 365 days;
    bool public autoArchiveEnabled = false;

    // Statistics
    uint256 private _totalLogs;
    uint256 private _logsLast24h;
    uint256 private _criticalEvents;
    uint256 private _lastStatsReset;

    // Event counters
    mapping(AuditEventType => uint256) private _eventCounts;
    mapping(address => uint256) private _userActivityCounts;

    // ============ Constructor ============

    constructor() BaseModule(MODULE_ID, MODULE_NAME, MODULE_VERSION) {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(LOGGER_ROLE, msg.sender);
        _grantRole(AUDITOR_ROLE, msg.sender);
        _lastStatsReset = block.timestamp;
    }

    // ============ BaseModule Overrides ============

    function getDependencies() external pure override returns (bytes32[] memory) {
        // No dependencies
        return new bytes32[](0);
    }

    function healthCheck()
        external
        view
        override(IModule, BaseModule)
        returns (bool healthy, string memory message)
    {
        (bool baseHealthy, string memory baseMessage) = BaseModule.healthCheck();
        if (!baseHealthy) {
            return (false, baseMessage);
        }

        if (_auditLogs.length >= maxActiveLogs) {
            return (false, "Audit log storage full");
        }

        return (true, "Audit module healthy");
    }

    // ============ Logging Functions ============

    function logEvent(
        AuditEventType eventType,
        Severity severity,
        address actor,
        address target,
        uint256 value,
        string calldata metadata
    ) external override onlyRole(LOGGER_ROLE) returns (uint256 logId) {
        return _createLogEntry(eventType, severity, actor, target, value, "", metadata);
    }

    function logEventWithData(
        AuditEventType eventType,
        Severity severity,
        address actor,
        address target,
        uint256 value,
        bytes calldata data,
        string calldata metadata
    ) external override onlyRole(LOGGER_ROLE) returns (uint256 logId) {
        return _createLogEntry(eventType, severity, actor, target, value, data, metadata);
    }

    function logEventBatch(
        AuditEventType[] calldata eventTypes,
        Severity[] calldata severities,
        address[] calldata actors,
        address[] calldata targets,
        uint256[] calldata values,
        string[] calldata metadataArray
    ) external override onlyRole(LOGGER_ROLE) returns (uint256[] memory logIds) {
        require(
            eventTypes.length == severities.length &&
            eventTypes.length == actors.length &&
            eventTypes.length == targets.length &&
            eventTypes.length == values.length &&
            eventTypes.length == metadataArray.length,
            "Array length mismatch"
        );

        logIds = new uint256[](eventTypes.length);

        for (uint256 i = 0; i < eventTypes.length; i++) {
            logIds[i] = _createLogEntry(
                eventTypes[i],
                severities[i],
                actors[i],
                targets[i],
                values[i],
                "",
                metadataArray[i]
            );
        }

        return logIds;
    }

    // ============ Specialized Logging ============

    function logDeposit(address user, address token, uint256 amount)
        external
        override
        onlyRole(LOGGER_ROLE)
    {
        _createLogEntry(
            AuditEventType.DEPOSIT,
            Severity.INFO,
            user,
            token,
            amount,
            "",
            "Deposit event"
        );
    }

    function logWithdrawal(address user, address token, uint256 amount)
        external
        override
        onlyRole(LOGGER_ROLE)
    {
        _createLogEntry(
            AuditEventType.WITHDRAWAL,
            Severity.INFO,
            user,
            token,
            amount,
            "",
            "Withdrawal event"
        );
    }

    function logLiquidation(
        address liquidator,
        address liquidated,
        uint256 debtCovered,
        uint256 collateralSeized
    ) external override onlyRole(LOGGER_ROLE) {
        bytes memory data = abi.encode(debtCovered, collateralSeized);

        _createLogEntry(
            AuditEventType.LIQUIDATION,
            Severity.WARNING,
            liquidator,
            liquidated,
            debtCovered,
            data,
            "Liquidation event"
        );
    }

    function logPriceUpdate(address token, uint256 oldPrice, uint256 newPrice)
        external
        override
        onlyRole(LOGGER_ROLE)
    {
        bytes memory data = abi.encode(oldPrice, newPrice);

        _createLogEntry(
            AuditEventType.PRICE_UPDATE,
            Severity.INFO,
            msg.sender,
            token,
            newPrice,
            data,
            "Price update"
        );
    }

    function logConfigChange(string calldata parameter, uint256 oldValue, uint256 newValue)
        external
        override
        onlyRole(LOGGER_ROLE)
    {
        bytes memory data = abi.encode(parameter, oldValue, newValue);

        _createLogEntry(
            AuditEventType.CONFIG_CHANGE,
            Severity.WARNING,
            msg.sender,
            address(0),
            newValue,
            data,
            parameter
        );
    }

    function logEmergencyAction(string calldata action, address initiator)
        external
        override
        onlyRole(LOGGER_ROLE)
    {
        _createLogEntry(
            AuditEventType.EMERGENCY_ACTION,
            Severity.CRITICAL,
            initiator,
            address(0),
            0,
            "",
            action
        );
    }

    // ============ Query Functions ============

    function getAuditEntry(uint256 logId)
        external
        view
        override
        returns (AuditEntry memory)
    {
        require(logId < _auditLogs.length, "Invalid log ID");
        return _auditLogs[logId];
    }

    function getAuditLogs(
        AuditFilter calldata filter,
        uint256 offset,
        uint256 limit
    ) external view override returns (AuditEntry[] memory) {
        uint256 count = 0;
        uint256 maxResults = limit > 100 ? 100 : limit;

        // Count matching logs
        for (uint256 i = offset; i < _auditLogs.length && count < maxResults; i++) {
            if (_matchesFilter(_auditLogs[i], filter)) {
                count++;
            }
        }

        // Collect matching logs
        AuditEntry[] memory result = new AuditEntry[](count);
        uint256 index = 0;

        for (uint256 i = offset; i < _auditLogs.length && index < count; i++) {
            if (_matchesFilter(_auditLogs[i], filter)) {
                result[index] = _auditLogs[i];
                index++;
            }
        }

        return result;
    }

    function getRecentLogs(uint256 count)
        external
        view
        override
        returns (AuditEntry[] memory)
    {
        uint256 length = _auditLogs.length;
        uint256 numLogs = count > length ? length : count;
        uint256 maxLogs = numLogs > 100 ? 100 : numLogs;

        AuditEntry[] memory result = new AuditEntry[](maxLogs);
        uint256 startIndex = length - maxLogs;

        for (uint256 i = 0; i < maxLogs; i++) {
            result[i] = _auditLogs[startIndex + i];
        }

        return result;
    }

    function getUserLogs(address user, uint256 offset, uint256 limit)
        external
        view
        override
        returns (AuditEntry[] memory)
    {
        uint256[] memory indices = _userLogIndices[user];
        return _getLogsByIndices(indices, offset, limit);
    }

    function getLogsByEventType(AuditEventType eventType, uint256 offset, uint256 limit)
        external
        view
        override
        returns (AuditEntry[] memory)
    {
        uint256[] memory indices = _eventTypeIndices[eventType];
        return _getLogsByIndices(indices, offset, limit);
    }

    function getLogsInTimeRange(uint256 fromTimestamp, uint256 toTimestamp, uint256 limit)
        external
        view
        override
        returns (AuditEntry[] memory)
    {
        require(fromTimestamp <= toTimestamp, "Invalid time range");

        uint256 count = 0;
        uint256 maxResults = limit > 100 ? 100 : limit;

        // Count matching logs
        for (uint256 i = 0; i < _auditLogs.length && count < maxResults; i++) {
            if (_auditLogs[i].timestamp >= fromTimestamp &&
                _auditLogs[i].timestamp <= toTimestamp) {
                count++;
            }
        }

        // Collect matching logs
        AuditEntry[] memory result = new AuditEntry[](count);
        uint256 index = 0;

        for (uint256 i = 0; i < _auditLogs.length && index < count; i++) {
            if (_auditLogs[i].timestamp >= fromTimestamp &&
                _auditLogs[i].timestamp <= toTimestamp) {
                result[index] = _auditLogs[i];
                index++;
            }
        }

        return result;
    }

    function getModuleLogs(bytes32 moduleId, uint256 offset, uint256 limit)
        external
        view
        override
        returns (AuditEntry[] memory)
    {
        uint256[] memory indices = _moduleLogIndices[moduleId];
        return _getLogsByIndices(indices, offset, limit);
    }

    // ============ Statistics & Reporting ============

    function getAuditStats()
        external
        view
        override
        returns (AuditStats memory)
    {
        _update24hStats();

        return AuditStats({
            totalLogs: _totalLogs,
            logsLast24h: _logsLast24h,
            criticalEvents: _criticalEvents,
            uniqueActors: 0, // Would require additional tracking
            lastLogTime: _auditLogs.length > 0 ? _auditLogs[_auditLogs.length - 1].timestamp : 0
        });
    }

    function getEventCount(AuditEventType eventType)
        external
        view
        override
        returns (uint256)
    {
        return _eventCounts[eventType];
    }

    function getUserActivityCount(address user)
        external
        view
        override
        returns (uint256)
    {
        return _userActivityCounts[user];
    }

    function generateComplianceReport(uint256 fromTimestamp, uint256 toTimestamp)
        external
        view
        override
        onlyRole(AUDITOR_ROLE)
        returns (
            uint256 totalEvents,
            uint256 criticalEvents,
            uint256 uniqueActors,
            uint256[] memory eventsByType
        )
    {
        eventsByType = new uint256[](15); // Number of event types

        for (uint256 i = 0; i < _auditLogs.length; i++) {
            if (_auditLogs[i].timestamp >= fromTimestamp &&
                _auditLogs[i].timestamp <= toTimestamp) {
                totalEvents++;

                if (_auditLogs[i].severity == Severity.CRITICAL) {
                    criticalEvents++;
                }

                eventsByType[uint256(_auditLogs[i].eventType)]++;
            }
        }

        uniqueActors = 0; // Would require set tracking

        return (totalEvents, criticalEvents, uniqueActors, eventsByType);
    }

    function getCriticalEvents(uint256 fromTimestamp, uint256 toTimestamp)
        external
        view
        override
        returns (AuditEntry[] memory)
    {
        uint256 count = 0;

        // Count critical events
        for (uint256 i = 0; i < _auditLogs.length; i++) {
            if (_auditLogs[i].timestamp >= fromTimestamp &&
                _auditLogs[i].timestamp <= toTimestamp &&
                _auditLogs[i].severity == Severity.CRITICAL) {
                count++;
            }
        }

        // Collect critical events
        AuditEntry[] memory result = new AuditEntry[](count);
        uint256 index = 0;

        for (uint256 i = 0; i < _auditLogs.length && index < count; i++) {
            if (_auditLogs[i].timestamp >= fromTimestamp &&
                _auditLogs[i].timestamp <= toTimestamp &&
                _auditLogs[i].severity == Severity.CRITICAL) {
                result[index] = _auditLogs[i];
                index++;
            }
        }

        return result;
    }

    // ============ Export & Archive ============

    function exportAuditLogs(
        AuditFilter calldata filter,
        uint256 offset,
        uint256 limit
    ) external view override onlyRole(AUDITOR_ROLE) returns (AuditEntry[] memory) {
        return this.getAuditLogs(filter, offset, limit);
    }

    function archiveLogs(uint256 beforeTimestamp)
        external
        override
        onlyRole(DEFAULT_ADMIN_ROLE)
        returns (uint256 archivedCount)
    {
        // In production, this would move logs to archive storage
        // For now, just count how many would be archived
        for (uint256 i = 0; i < _auditLogs.length; i++) {
            if (_auditLogs[i].timestamp < beforeTimestamp) {
                archivedCount++;
            }
        }

        return archivedCount;
    }

    // ============ Configuration ============

    function setMaxActiveLogs(uint256 maxLogs)
        external
        override
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        require(maxLogs > 0, "Invalid max logs");
        maxActiveLogs = maxLogs;
    }

    function setRetentionPeriod(uint256 period)
        external
        override
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        require(period > 0, "Invalid period");
        retentionPeriod = period;
        emit RetentionPolicyUpdated(period);
    }

    function setAutoArchive(bool enabled)
        external
        override
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        autoArchiveEnabled = enabled;
    }

    function getAuditConfig()
        external
        view
        override
        returns (
            uint256 _maxActiveLogs,
            uint256 _retentionPeriod,
            bool _autoArchiveEnabled
        )
    {
        return (maxActiveLogs, retentionPeriod, autoArchiveEnabled);
    }

    // ============ Internal Functions ============

    function _createLogEntry(
        AuditEventType eventType,
        Severity severity,
        address actor,
        address target,
        uint256 value,
        bytes memory data,
        string memory metadata
    ) internal returns (uint256 logId) {
        logId = _auditLogs.length;

        bytes32 txHash = keccak256(
            abi.encodePacked(tx.origin, block.timestamp, logId)
        );

        AuditEntry memory entry = AuditEntry({
            logId: logId,
            timestamp: block.timestamp,
            blockNumber: block.number,
            eventType: eventType,
            severity: severity,
            moduleId: MODULE_ID,
            actor: actor,
            target: target,
            value: value,
            data: data,
            metadata: metadata,
            transactionHash: txHash
        });

        _auditLogs.push(entry);
        _userLogIndices[actor].push(logId);
        _eventTypeIndices[eventType].push(logId);
        _moduleLogIndices[MODULE_ID].push(logId);

        // Update statistics
        _totalLogs++;
        _eventCounts[eventType]++;
        _userActivityCounts[actor]++;

        if (severity == Severity.CRITICAL) {
            _criticalEvents++;
        }

        _update24hStats();

        emit AuditLogCreated(logId, eventType, actor, block.timestamp);
        emit AuditEventLogged(MODULE_ID, eventType, severity, actor);

        if (severity == Severity.CRITICAL) {
            emit CriticalEventDetected(logId, eventType, actor, metadata);
        }

        return logId;
    }

    function _matchesFilter(AuditEntry memory entry, AuditFilter memory filter)
        internal
        pure
        returns (bool)
    {
        if (filter.moduleId != bytes32(0) && filter.moduleId != entry.moduleId) {
            return false;
        }

        if (filter.actor != address(0) && filter.actor != entry.actor) {
            return false;
        }

        if (uint8(entry.severity) < uint8(filter.minSeverity)) {
            return false;
        }

        if (filter.fromTimestamp > 0 && entry.timestamp < filter.fromTimestamp) {
            return false;
        }

        if (filter.toTimestamp > 0 && entry.timestamp > filter.toTimestamp) {
            return false;
        }

        if (filter.fromBlock > 0 && entry.blockNumber < filter.fromBlock) {
            return false;
        }

        if (filter.toBlock > 0 && entry.blockNumber > filter.toBlock) {
            return false;
        }

        return true;
    }

    function _getLogsByIndices(
        uint256[] memory indices,
        uint256 offset,
        uint256 limit
    ) internal view returns (AuditEntry[] memory) {
        uint256 length = indices.length;

        if (offset >= length) {
            return new AuditEntry[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }
        if (end - offset > 100) {
            end = offset + 100;
        }

        AuditEntry[] memory result = new AuditEntry[](end - offset);

        for (uint256 i = offset; i < end; i++) {
            result[i - offset] = _auditLogs[indices[i]];
        }

        return result;
    }

    function _update24hStats() internal view {
        // In production, maintain rolling 24h window
        // For now, simple implementation
    }

    // ============ Override Required Functions ============

    function pause() external override(IModule, BaseModule) onlyRole(DEFAULT_ADMIN_ROLE) {
        BaseModule.pause();
    }

    function unpause() external override(IModule, BaseModule) onlyRole(DEFAULT_ADMIN_ROLE) {
        BaseModule.unpause();
    }
}
