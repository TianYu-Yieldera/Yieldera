// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

/**
 * @title AuditLogger
 * @notice Comprehensive audit logging system for DeFi protocol operations
 * @dev Records all critical actions with timestamps and metadata for compliance and analysis
 */
contract AuditLogger is AccessControl, Pausable {
    // Roles
    bytes32 public constant LOGGER_ROLE = keccak256("LOGGER_ROLE");
    bytes32 public constant AUDITOR_ROLE = keccak256("AUDITOR_ROLE");

    // Event types enumeration
    enum EventType {
        DEPOSIT,
        WITHDRAWAL,
        MINT,
        BURN,
        TRANSFER,
        LIQUIDATION,
        DEBT_INCREASE,
        DEBT_DECREASE,
        COLLATERAL_DEPOSIT,
        COLLATERAL_WITHDRAWAL,
        PRICE_UPDATE,
        ROLE_CHANGE,
        EMERGENCY_ACTION,
        CONFIGURATION_CHANGE
    }

    // Audit log entry structure
    struct LogEntry {
        uint256 timestamp;
        EventType eventType;
        address actor;
        address target;
        uint256 value;
        string metadata;
        bytes32 transactionHash;
        uint256 blockNumber;
    }

    // Storage
    LogEntry[] public auditLogs;
    mapping(address => uint256[]) public userLogs; // User -> log indices
    mapping(EventType => uint256[]) public eventTypeLogs; // EventType -> log indices
    mapping(bytes32 => uint256) public transactionLogs; // TxHash -> log index

    // Statistics
    uint256 public totalLogs;
    mapping(EventType => uint256) public eventCounts;
    mapping(address => uint256) public userActionCounts;

    // Configuration
    uint256 public maxLogsPerQuery = 1000;
    bool public detailedLogging = true;

    // Events
    event LogCreated(
        uint256 indexed logId,
        EventType indexed eventType,
        address indexed actor,
        uint256 timestamp
    );
    event LoggingPaused();
    event LoggingResumed();
    event ConfigurationUpdated(string parameter, uint256 value);

    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(LOGGER_ROLE, msg.sender);
        _grantRole(AUDITOR_ROLE, msg.sender);
    }

    /**
     * @notice Log a protocol event
     * @param eventType Type of event
     * @param actor Address performing the action
     * @param target Target address (if applicable)
     * @param value Numeric value associated with the event
     * @param metadata Additional information
     */
    function logEvent(
        EventType eventType,
        address actor,
        address target,
        uint256 value,
        string calldata metadata
    ) external onlyRole(LOGGER_ROLE) whenNotPaused {
        uint256 logId = auditLogs.length;

        LogEntry memory newLog = LogEntry({
            timestamp: block.timestamp,
            eventType: eventType,
            actor: actor,
            target: target,
            value: value,
            metadata: metadata,
            transactionHash: keccak256(abi.encodePacked(tx.origin, block.timestamp, logId)),
            blockNumber: block.number
        });

        auditLogs.push(newLog);
        userLogs[actor].push(logId);
        eventTypeLogs[eventType].push(logId);
        transactionLogs[newLog.transactionHash] = logId;

        // Update statistics
        totalLogs++;
        eventCounts[eventType]++;
        userActionCounts[actor]++;

        emit LogCreated(logId, eventType, actor, block.timestamp);
    }

    /**
     * @notice Log a deposit event
     */
    function logDeposit(
        address user,
        address token,
        uint256 amount
    ) external onlyRole(LOGGER_ROLE) whenNotPaused {
        string memory metadata = string(
            abi.encodePacked("Deposit: ", _addressToString(token))
        );
        logEvent(EventType.DEPOSIT, user, token, amount, metadata);
    }

    /**
     * @notice Log a withdrawal event
     */
    function logWithdrawal(
        address user,
        address token,
        uint256 amount
    ) external onlyRole(LOGGER_ROLE) whenNotPaused {
        string memory metadata = string(
            abi.encodePacked("Withdrawal: ", _addressToString(token))
        );
        logEvent(EventType.WITHDRAWAL, user, token, amount, metadata);
    }

    /**
     * @notice Log a liquidation event
     */
    function logLiquidation(
        address liquidator,
        address liquidated,
        uint256 debtCovered,
        uint256 collateralSeized
    ) external onlyRole(LOGGER_ROLE) whenNotPaused {
        string memory metadata = string(
            abi.encodePacked(
                "Liquidation: debt=",
                _uint256ToString(debtCovered),
                ", collateral=",
                _uint256ToString(collateralSeized)
            )
        );
        logEvent(EventType.LIQUIDATION, liquidator, liquidated, debtCovered, metadata);
    }

    /**
     * @notice Get logs for a specific user
     * @param user User address
     * @param offset Starting index
     * @param limit Number of logs to return
     */
    function getUserLogs(
        address user,
        uint256 offset,
        uint256 limit
    ) external view returns (LogEntry[] memory) {
        uint256[] memory logIndices = userLogs[user];
        uint256 length = logIndices.length;

        if (offset >= length) {
            return new LogEntry[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }
        if (end - offset > maxLogsPerQuery) {
            end = offset + maxLogsPerQuery;
        }

        LogEntry[] memory result = new LogEntry[](end - offset);
        for (uint256 i = offset; i < end; i++) {
            result[i - offset] = auditLogs[logIndices[i]];
        }

        return result;
    }

    /**
     * @notice Get logs by event type
     * @param eventType Type of event
     * @param offset Starting index
     * @param limit Number of logs to return
     */
    function getEventTypeLogs(
        EventType eventType,
        uint256 offset,
        uint256 limit
    ) external view returns (LogEntry[] memory) {
        uint256[] memory logIndices = eventTypeLogs[eventType];
        uint256 length = logIndices.length;

        if (offset >= length) {
            return new LogEntry[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }
        if (end - offset > maxLogsPerQuery) {
            end = offset + maxLogsPerQuery;
        }

        LogEntry[] memory result = new LogEntry[](end - offset);
        for (uint256 i = offset; i < end; i++) {
            result[i - offset] = auditLogs[logIndices[i]];
        }

        return result;
    }

    /**
     * @notice Get logs within a time range
     * @param startTime Start timestamp
     * @param endTime End timestamp
     * @param limit Maximum number of logs to return
     */
    function getLogsInTimeRange(
        uint256 startTime,
        uint256 endTime,
        uint256 limit
    ) external view returns (LogEntry[] memory) {
        require(startTime <= endTime, "Invalid time range");

        uint256 count = 0;
        uint256 maxLimit = limit > maxLogsPerQuery ? maxLogsPerQuery : limit;

        // First pass: count matching logs
        for (uint256 i = 0; i < auditLogs.length && count < maxLimit; i++) {
            if (auditLogs[i].timestamp >= startTime && auditLogs[i].timestamp <= endTime) {
                count++;
            }
        }

        // Second pass: collect matching logs
        LogEntry[] memory result = new LogEntry[](count);
        uint256 index = 0;

        for (uint256 i = 0; i < auditLogs.length && index < count; i++) {
            if (auditLogs[i].timestamp >= startTime && auditLogs[i].timestamp <= endTime) {
                result[index] = auditLogs[i];
                index++;
            }
        }

        return result;
    }

    /**
     * @notice Get recent logs
     * @param count Number of recent logs to return
     */
    function getRecentLogs(uint256 count) external view returns (LogEntry[] memory) {
        uint256 length = auditLogs.length;
        if (count > length) {
            count = length;
        }
        if (count > maxLogsPerQuery) {
            count = maxLogsPerQuery;
        }

        LogEntry[] memory result = new LogEntry[](count);
        uint256 startIndex = length - count;

        for (uint256 i = 0; i < count; i++) {
            result[i] = auditLogs[startIndex + i];
        }

        return result;
    }

    /**
     * @notice Get statistics for an event type
     */
    function getEventStatistics(EventType eventType)
        external
        view
        returns (uint256 count, uint256 lastOccurrence)
    {
        count = eventCounts[eventType];
        uint256[] memory indices = eventTypeLogs[eventType];

        if (indices.length > 0) {
            lastOccurrence = auditLogs[indices[indices.length - 1]].timestamp;
        }
    }

    /**
     * @notice Get user activity statistics
     */
    function getUserStatistics(address user)
        external
        view
        returns (uint256 actionCount, uint256 lastAction)
    {
        actionCount = userActionCounts[user];
        uint256[] memory indices = userLogs[user];

        if (indices.length > 0) {
            lastAction = auditLogs[indices[indices.length - 1]].timestamp;
        }
    }

    /**
     * @notice Generate compliance report
     * @param startTime Start of reporting period
     * @param endTime End of reporting period
     */
    function generateComplianceReport(uint256 startTime, uint256 endTime)
        external
        view
        onlyRole(AUDITOR_ROLE)
        returns (
            uint256 totalEvents,
            uint256 deposits,
            uint256 withdrawals,
            uint256 liquidations,
            uint256 uniqueUsers
        )
    {
        for (uint256 i = 0; i < auditLogs.length; i++) {
            if (auditLogs[i].timestamp >= startTime && auditLogs[i].timestamp <= endTime) {
                totalEvents++;

                if (auditLogs[i].eventType == EventType.DEPOSIT) deposits++;
                else if (auditLogs[i].eventType == EventType.WITHDRAWAL) withdrawals++;
                else if (auditLogs[i].eventType == EventType.LIQUIDATION) liquidations++;
            }
        }

        // Note: Unique users count would require additional tracking
        uniqueUsers = 0; // Placeholder
    }

    /**
     * @notice Export logs to CSV format (off-chain processing)
     * @param logIds Array of log IDs to export
     */
    function exportLogs(uint256[] calldata logIds)
        external
        view
        onlyRole(AUDITOR_ROLE)
        returns (LogEntry[] memory)
    {
        require(logIds.length <= maxLogsPerQuery, "Too many logs requested");

        LogEntry[] memory result = new LogEntry[](logIds.length);

        for (uint256 i = 0; i < logIds.length; i++) {
            require(logIds[i] < auditLogs.length, "Invalid log ID");
            result[i] = auditLogs[logIds[i]];
        }

        return result;
    }

    /**
     * @notice Update configuration
     */
    function setMaxLogsPerQuery(uint256 _max) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(_max > 0 && _max <= 10000, "Invalid max logs");
        maxLogsPerQuery = _max;
        emit ConfigurationUpdated("maxLogsPerQuery", _max);
    }

    function setDetailedLogging(bool _enabled) external onlyRole(DEFAULT_ADMIN_ROLE) {
        detailedLogging = _enabled;
        emit ConfigurationUpdated("detailedLogging", _enabled ? 1 : 0);
    }

    /**
     * @notice Pause logging
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
        emit LoggingPaused();
    }

    /**
     * @notice Resume logging
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
        emit LoggingResumed();
    }

    // Helper functions
    function _addressToString(address _addr) internal pure returns (string memory) {
        bytes32 value = bytes32(uint256(uint160(_addr)));
        bytes memory alphabet = "0123456789abcdef";
        bytes memory str = new bytes(42);

        str[0] = '0';
        str[1] = 'x';

        for (uint256 i = 0; i < 20; i++) {
            str[2+i*2] = alphabet[uint8(value[i + 12] >> 4)];
            str[3+i*2] = alphabet[uint8(value[i + 12] & 0x0f)];
        }

        return string(str);
    }

    function _uint256ToString(uint256 value) internal pure returns (string memory) {
        if (value == 0) {
            return "0";
        }

        uint256 temp = value;
        uint256 digits;

        while (temp != 0) {
            digits++;
            temp /= 10;
        }

        bytes memory buffer = new bytes(digits);

        while (value != 0) {
            digits -= 1;
            buffer[digits] = bytes1(uint8(48 + uint256(value % 10)));
            value /= 10;
        }

        return string(buffer);
    }
}