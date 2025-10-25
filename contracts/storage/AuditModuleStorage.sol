// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../interfaces/core/IAuditModule.sol";

/**
 * @title AuditModuleStorage
 * @notice Diamond Storage library for AuditModule
 * @dev Implements EIP-2535 Diamond Storage pattern
 */
library AuditModuleStorage {
    // Storage position is keccak256("audit.module.storage") - 1
    bytes32 constant STORAGE_POSITION =
        0x4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e;

    struct AuditData {
        // Audit log storage
        IAuditModule.AuditLog[] logs;

        // Indices for efficient querying
        mapping(address => uint256[]) userLogIndices;
        mapping(uint8 => uint256[]) eventTypeIndices;   // eventType -> log indices
        mapping(uint8 => uint256[]) severityIndices;     // severity -> log indices

        // Statistics
        uint256 totalLogs;
        uint256 criticalCount;
        uint256 errorCount;
        uint256 warningCount;
        uint256 infoCount;

        // Configuration
        uint256 maxLogsPerQuery;  // Default: 100
        uint256 retentionPeriod;  // Default: 90 days
        bool autoArchive;          // Default: false

        // Metadata
        uint256 lastLogTime;
        uint256 lastUpdate;

        // Reserved slots for future upgrades
        uint256[50] __gap;
    }

    /**
     * @notice Returns the storage layout
     */
    function layout() internal pure returns (AuditData storage ds) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            ds.slot := position
        }
    }

    /**
     * @notice Initialize storage
     */
    function initialize() internal {
        AuditData storage ds = layout();
        require(ds.lastUpdate == 0, "Already initialized");

        ds.maxLogsPerQuery = 100;
        ds.retentionPeriod = 90 days;
        ds.autoArchive = false;
        ds.lastUpdate = block.timestamp;
    }

    /**
     * @notice Add audit log
     */
    function addLog(IAuditModule.AuditLog memory log) internal returns (uint256) {
        AuditData storage ds = layout();

        uint256 logIndex = ds.logs.length;
        ds.logs.push(log);

        // Index by user
        ds.userLogIndices[log.user].push(logIndex);

        // Index by event type
        ds.eventTypeIndices[uint8(log.eventType)].push(logIndex);

        // Index by severity
        ds.severityIndices[uint8(log.severity)].push(logIndex);

        // Update statistics
        ds.totalLogs++;
        if (log.severity == IAuditModule.Severity.CRITICAL) {
            ds.criticalCount++;
        } else if (log.severity == IAuditModule.Severity.ERROR) {
            ds.errorCount++;
        } else if (log.severity == IAuditModule.Severity.WARNING) {
            ds.warningCount++;
        } else {
            ds.infoCount++;
        }

        ds.lastLogTime = block.timestamp;
        ds.lastUpdate = block.timestamp;

        return logIndex;
    }

    /**
     * @notice Get audit logs
     */
    function getLogs(uint256 offset, uint256 limit)
        internal
        view
        returns (IAuditModule.AuditLog[] memory)
    {
        AuditData storage ds = layout();
        uint256 length = ds.logs.length;

        if (offset >= length) {
            return new IAuditModule.AuditLog[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }

        IAuditModule.AuditLog[] memory result = new IAuditModule.AuditLog[](end - offset);
        for (uint256 i = offset; i < end; i++) {
            result[i - offset] = ds.logs[i];
        }

        return result;
    }

    /**
     * @notice Get logs by user
     */
    function getLogsByUser(address user, uint256 offset, uint256 limit)
        internal
        view
        returns (IAuditModule.AuditLog[] memory)
    {
        AuditData storage ds = layout();
        uint256[] storage indices = ds.userLogIndices[user];
        uint256 length = indices.length;

        if (offset >= length) {
            return new IAuditModule.AuditLog[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }

        IAuditModule.AuditLog[] memory result = new IAuditModule.AuditLog[](end - offset);
        for (uint256 i = offset; i < end; i++) {
            result[i - offset] = ds.logs[indices[i]];
        }

        return result;
    }

    /**
     * @notice Get logs by event type
     */
    function getLogsByEventType(
        IAuditModule.EventType eventType,
        uint256 offset,
        uint256 limit
    ) internal view returns (IAuditModule.AuditLog[] memory) {
        AuditData storage ds = layout();
        uint256[] storage indices = ds.eventTypeIndices[uint8(eventType)];
        uint256 length = indices.length;

        if (offset >= length) {
            return new IAuditModule.AuditLog[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }

        IAuditModule.AuditLog[] memory result = new IAuditModule.AuditLog[](end - offset);
        for (uint256 i = offset; i < end; i++) {
            result[i - offset] = ds.logs[indices[i]];
        }

        return result;
    }

    /**
     * @notice Get logs by severity
     */
    function getLogsBySeverity(
        IAuditModule.Severity severity,
        uint256 offset,
        uint256 limit
    ) internal view returns (IAuditModule.AuditLog[] memory) {
        AuditData storage ds = layout();
        uint256[] storage indices = ds.severityIndices[uint8(severity)];
        uint256 length = indices.length;

        if (offset >= length) {
            return new IAuditModule.AuditLog[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }

        IAuditModule.AuditLog[] memory result = new IAuditModule.AuditLog[](end - offset);
        for (uint256 i = offset; i < end; i++) {
            result[i - offset] = ds.logs[indices[i]];
        }

        return result;
    }

    /**
     * @notice Get logs by time range
     */
    function getLogsByTimeRange(
        uint256 startTime,
        uint256 endTime,
        uint256 offset,
        uint256 limit
    ) internal view returns (IAuditModule.AuditLog[] memory) {
        AuditData storage ds = layout();
        uint256 length = ds.logs.length;

        // Count matching logs
        uint256 matchCount = 0;
        for (uint256 i = 0; i < length && matchCount < offset + limit; i++) {
            if (ds.logs[i].timestamp >= startTime && ds.logs[i].timestamp <= endTime) {
                matchCount++;
            }
        }

        if (matchCount <= offset) {
            return new IAuditModule.AuditLog[](0);
        }

        uint256 resultSize = matchCount - offset;
        if (resultSize > limit) {
            resultSize = limit;
        }

        IAuditModule.AuditLog[] memory result = new IAuditModule.AuditLog[](resultSize);
        uint256 resultIndex = 0;
        uint256 skipCount = 0;

        for (uint256 i = 0; i < length && resultIndex < resultSize; i++) {
            if (ds.logs[i].timestamp >= startTime && ds.logs[i].timestamp <= endTime) {
                if (skipCount >= offset) {
                    result[resultIndex] = ds.logs[i];
                    resultIndex++;
                } else {
                    skipCount++;
                }
            }
        }

        return result;
    }

    /**
     * @notice Get statistics
     */
    function getStatistics()
        internal
        view
        returns (
            uint256 totalLogs,
            uint256 criticalCount,
            uint256 errorCount,
            uint256 warningCount,
            uint256 infoCount
        )
    {
        AuditData storage ds = layout();
        return (
            ds.totalLogs,
            ds.criticalCount,
            ds.errorCount,
            ds.warningCount,
            ds.infoCount
        );
    }

    /**
     * @notice Get configuration
     */
    function getConfig()
        internal
        view
        returns (uint256 maxLogsPerQuery, uint256 retentionPeriod, bool autoArchive)
    {
        AuditData storage ds = layout();
        return (ds.maxLogsPerQuery, ds.retentionPeriod, ds.autoArchive);
    }

    /**
     * @notice Update configuration
     */
    function setConfig(uint256 maxLogs, uint256 retention, bool archive) internal {
        AuditData storage ds = layout();
        ds.maxLogsPerQuery = maxLogs;
        ds.retentionPeriod = retention;
        ds.autoArchive = archive;
        ds.lastUpdate = block.timestamp;
    }

    /**
     * @notice Update last modified timestamp
     */
    function touch() internal {
        layout().lastUpdate = block.timestamp;
    }
}
