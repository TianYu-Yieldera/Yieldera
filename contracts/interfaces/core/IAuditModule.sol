// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IModule.sol";

/**
 * @title IAuditModule
 * @notice Standard interface for audit and logging modules
 * @dev Extends IModule with comprehensive audit functionality
 *
 * Purpose:
 * - Record all critical system events
 * - Provide audit trails for compliance
 * - Support forensic analysis
 * - Enable monitoring and alerting
 */
interface IAuditModule is IModule {
    /**
     * @notice Event types for categorization
     */
    enum AuditEventType {
        DEPOSIT,
        WITHDRAWAL,
        MINT,
        BURN,
        TRANSFER,
        LIQUIDATION,
        DEBT_CHANGE,
        COLLATERAL_CHANGE,
        PRICE_UPDATE,
        CONFIG_CHANGE,
        ROLE_CHANGE,
        EMERGENCY_ACTION,
        GOVERNANCE,
        MODULE_LIFECYCLE,
        SYSTEM_EVENT
    }

    /**
     * @notice Event severity levels
     */
    enum Severity {
        INFO,               // Informational
        WARNING,            // Warning
        ERROR,              // Error condition
        CRITICAL            // Critical event
    }

    /**
     * @notice Audit log entry
     */
    struct AuditEntry {
        uint256 logId;              // Unique log identifier
        uint256 timestamp;          // Event timestamp
        uint256 blockNumber;        // Block number
        AuditEventType eventType;   // Type of event
        Severity severity;          // Event severity
        bytes32 moduleId;           // Source module
        address actor;              // Address that triggered event
        address target;             // Target address (if applicable)
        uint256 value;              // Numeric value (amount, etc.)
        bytes data;                 // Additional encoded data
        string metadata;            // Human-readable metadata
        bytes32 transactionHash;    // Transaction hash
    }

    /**
     * @notice Query filter for log retrieval
     */
    struct AuditFilter {
        AuditEventType eventType;   // Filter by event type (optional)
        Severity minSeverity;       // Minimum severity level
        bytes32 moduleId;           // Filter by module (bytes32(0) = all)
        address actor;              // Filter by actor (address(0) = all)
        uint256 fromTimestamp;      // Start time
        uint256 toTimestamp;        // End time
        uint256 fromBlock;          // Start block
        uint256 toBlock;            // End block
    }

    /**
     * @notice Statistics for compliance reporting
     */
    struct AuditStats {
        uint256 totalLogs;          // Total number of logs
        uint256 logsLast24h;        // Logs in last 24 hours
        uint256 criticalEvents;     // Number of critical events
        uint256 uniqueActors;       // Number of unique actors
        uint256 lastLogTime;        // Timestamp of last log
    }

    // ============ Events ============

    event AuditLogCreated(
        uint256 indexed logId,
        AuditEventType indexed eventType,
        address indexed actor,
        uint256 timestamp
    );

    event AuditEventLogged(
        bytes32 indexed moduleId,
        AuditEventType indexed eventType,
        Severity severity,
        address actor
    );

    event CriticalEventDetected(
        uint256 indexed logId,
        AuditEventType indexed eventType,
        address indexed actor,
        string reason
    );

    event AuditExported(
        address indexed requester,
        uint256 fromTimestamp,
        uint256 toTimestamp,
        uint256 recordCount
    );

    event RetentionPolicyUpdated(uint256 newRetentionPeriod);

    // ============ Logging Functions ============

    /**
     * @notice Log an audit event
     * @param eventType Type of event
     * @param severity Event severity
     * @param actor Address performing action
     * @param target Target address (if applicable)
     * @param value Numeric value
     * @param metadata Human-readable description
     * @return logId Unique log identifier
     */
    function logEvent(
        AuditEventType eventType,
        Severity severity,
        address actor,
        address target,
        uint256 value,
        string calldata metadata
    ) external returns (uint256 logId);

    /**
     * @notice Log event with encoded data
     * @param eventType Type of event
     * @param severity Event severity
     * @param actor Address performing action
     * @param target Target address
     * @param value Numeric value
     * @param data ABI-encoded additional data
     * @param metadata Human-readable description
     * @return logId Unique log identifier
     */
    function logEventWithData(
        AuditEventType eventType,
        Severity severity,
        address actor,
        address target,
        uint256 value,
        bytes calldata data,
        string calldata metadata
    ) external returns (uint256 logId);

    /**
     * @notice Batch log multiple events
     * @param eventTypes Array of event types
     * @param severities Array of severities
     * @param actors Array of actor addresses
     * @param targets Array of target addresses
     * @param values Array of values
     * @param metadataArray Array of metadata strings
     * @return logIds Array of log identifiers
     */
    function logEventBatch(
        AuditEventType[] calldata eventTypes,
        Severity[] calldata severities,
        address[] calldata actors,
        address[] calldata targets,
        uint256[] calldata values,
        string[] calldata metadataArray
    ) external returns (uint256[] memory logIds);

    // ============ Specialized Logging ============

    /**
     * @notice Log a deposit event
     * @param user User address
     * @param token Token address
     * @param amount Deposit amount
     */
    function logDeposit(address user, address token, uint256 amount) external;

    /**
     * @notice Log a withdrawal event
     * @param user User address
     * @param token Token address
     * @param amount Withdrawal amount
     */
    function logWithdrawal(address user, address token, uint256 amount) external;

    /**
     * @notice Log a liquidation event
     * @param liquidator Liquidator address
     * @param liquidated User being liquidated
     * @param debtCovered Debt amount covered
     * @param collateralSeized Collateral seized
     */
    function logLiquidation(
        address liquidator,
        address liquidated,
        uint256 debtCovered,
        uint256 collateralSeized
    ) external;

    /**
     * @notice Log a price update
     * @param token Token address
     * @param oldPrice Previous price
     * @param newPrice New price
     */
    function logPriceUpdate(address token, uint256 oldPrice, uint256 newPrice) external;

    /**
     * @notice Log a configuration change
     * @param parameter Parameter name
     * @param oldValue Previous value
     * @param newValue New value
     */
    function logConfigChange(string calldata parameter, uint256 oldValue, uint256 newValue) external;

    /**
     * @notice Log a critical emergency action
     * @param action Description of action
     * @param initiator Address that initiated action
     */
    function logEmergencyAction(string calldata action, address initiator) external;

    // ============ Query Functions ============

    /**
     * @notice Get audit entry by ID
     * @param logId Log identifier
     * @return Audit entry
     */
    function getAuditEntry(uint256 logId) external view returns (AuditEntry memory);

    /**
     * @notice Get logs by filter
     * @param filter Query filter
     * @param offset Starting index
     * @param limit Maximum results
     * @return Array of audit entries
     */
    function getAuditLogs(
        AuditFilter calldata filter,
        uint256 offset,
        uint256 limit
    ) external view returns (AuditEntry[] memory);

    /**
     * @notice Get recent logs
     * @param count Number of recent logs
     * @return Array of recent audit entries
     */
    function getRecentLogs(uint256 count) external view returns (AuditEntry[] memory);

    /**
     * @notice Get logs for a specific user
     * @param user User address
     * @param offset Starting index
     * @param limit Maximum results
     * @return Array of user's audit entries
     */
    function getUserLogs(address user, uint256 offset, uint256 limit)
        external
        view
        returns (AuditEntry[] memory);

    /**
     * @notice Get logs by event type
     * @param eventType Event type to filter
     * @param offset Starting index
     * @param limit Maximum results
     * @return Array of audit entries
     */
    function getLogsByEventType(AuditEventType eventType, uint256 offset, uint256 limit)
        external
        view
        returns (AuditEntry[] memory);

    /**
     * @notice Get logs within time range
     * @param fromTimestamp Start time
     * @param toTimestamp End time
     * @param limit Maximum results
     * @return Array of audit entries
     */
    function getLogsInTimeRange(uint256 fromTimestamp, uint256 toTimestamp, uint256 limit)
        external
        view
        returns (AuditEntry[] memory);

    /**
     * @notice Get logs by module
     * @param moduleId Module identifier
     * @param offset Starting index
     * @param limit Maximum results
     * @return Array of module's audit entries
     */
    function getModuleLogs(bytes32 moduleId, uint256 offset, uint256 limit)
        external
        view
        returns (AuditEntry[] memory);

    // ============ Statistics & Reporting ============

    /**
     * @notice Get audit statistics
     * @return Audit statistics structure
     */
    function getAuditStats() external view returns (AuditStats memory);

    /**
     * @notice Get event count by type
     * @param eventType Event type
     * @return Number of events of this type
     */
    function getEventCount(AuditEventType eventType) external view returns (uint256);

    /**
     * @notice Get user activity count
     * @param user User address
     * @return Number of logged actions by user
     */
    function getUserActivityCount(address user) external view returns (uint256);

    /**
     * @notice Generate compliance report
     * @param fromTimestamp Report start time
     * @param toTimestamp Report end time
     * @return totalEvents Total events in period
     * @return criticalEvents Critical events in period
     * @return uniqueActors Number of unique actors
     * @return eventsByType Array of counts per event type
     */
    function generateComplianceReport(uint256 fromTimestamp, uint256 toTimestamp)
        external
        view
        returns (
            uint256 totalEvents,
            uint256 criticalEvents,
            uint256 uniqueActors,
            uint256[] memory eventsByType
        );

    /**
     * @notice Get critical events in time range
     * @param fromTimestamp Start time
     * @param toTimestamp End time
     * @return Array of critical audit entries
     */
    function getCriticalEvents(uint256 fromTimestamp, uint256 toTimestamp)
        external
        view
        returns (AuditEntry[] memory);

    // ============ Export & Archive ============

    /**
     * @notice Export logs for external analysis
     * @param filter Query filter
     * @param offset Starting index
     * @param limit Maximum results
     * @return Array of audit entries for export
     * @dev Only callable by authorized auditors
     */
    function exportAuditLogs(
        AuditFilter calldata filter,
        uint256 offset,
        uint256 limit
    ) external view returns (AuditEntry[] memory);

    /**
     * @notice Archive old logs
     * @param beforeTimestamp Archive logs before this timestamp
     * @return archivedCount Number of logs archived
     * @dev Moves logs to archive storage to save gas
     */
    function archiveLogs(uint256 beforeTimestamp) external returns (uint256 archivedCount);

    // ============ Configuration ============

    /**
     * @notice Set maximum logs to keep in active storage
     * @param maxLogs Maximum number of logs
     */
    function setMaxActiveLogs(uint256 maxLogs) external;

    /**
     * @notice Set log retention period
     * @param retentionPeriod Period in seconds
     */
    function setRetentionPeriod(uint256 retentionPeriod) external;

    /**
     * @notice Enable/disable automatic archival
     * @param enabled Whether auto-archive is enabled
     */
    function setAutoArchive(bool enabled) external;

    /**
     * @notice Get audit module configuration
     * @return maxActiveLogs Maximum logs in active storage
     * @return retentionPeriod Retention period in seconds
     * @return autoArchiveEnabled Whether auto-archive is enabled
     */
    function getAuditConfig()
        external
        view
        returns (
            uint256 maxActiveLogs,
            uint256 retentionPeriod,
            bool autoArchiveEnabled
        );
}
