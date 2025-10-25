// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IEventHub
 * @notice Centralized event bus for inter-module communication
 * @dev Enables loose coupling between modules through event-based communication
 *
 * Design Principles:
 * - Modules publish events to the hub
 * - Modules can subscribe to events from other modules
 * - Prevents direct dependencies between modules
 * - Supports event filtering and routing
 * - Maintains event history for auditing
 */
interface IEventHub {
    /**
     * @notice Event types for categorization
     */
    enum EventCategory {
        SYSTEM,         // System-level events
        MODULE,         // Module lifecycle events
        TRANSACTION,    // Transaction-related events
        GOVERNANCE,     // Governance events
        ORACLE,         // Price oracle events
        AUDIT,          // Audit events
        CUSTOM          // Custom module events
    }

    /**
     * @notice Event severity levels
     */
    enum EventSeverity {
        INFO,           // Informational
        WARNING,        // Warning
        ERROR,          // Error
        CRITICAL        // Critical
    }

    /**
     * @notice Standardized event structure
     */
    struct HubEvent {
        bytes32 eventId;            // Unique event identifier
        bytes32 moduleId;           // Source module identifier
        EventCategory category;     // Event category
        EventSeverity severity;     // Event severity
        string eventType;           // Event type name (e.g., "DEPOSIT", "WITHDRAWAL")
        bytes eventData;            // ABI-encoded event data
        uint256 timestamp;          // Event timestamp
        uint256 blockNumber;        // Block number
        address emitter;            // Address that emitted the event
    }

    /**
     * @notice Event subscription
     */
    struct EventSubscription {
        bytes32 subscriptionId;     // Unique subscription identifier
        bytes32 subscriberModule;   // Subscribing module ID
        bytes32 publisherModule;    // Publisher module ID (bytes32(0) = all)
        EventCategory category;     // Category filter (optional)
        string eventType;           // Event type filter (empty = all)
        address callbackAddress;    // Address to call on event
        bytes4 callbackSelector;    // Function selector for callback
        bool isActive;              // Subscription status
    }

    /**
     * @notice Event filter criteria
     */
    struct EventFilter {
        bytes32 moduleId;           // Filter by module (bytes32(0) = all)
        EventCategory category;     // Filter by category
        EventSeverity minSeverity;  // Minimum severity level
        uint256 fromTimestamp;      // Start timestamp
        uint256 toTimestamp;        // End timestamp
        uint256 fromBlock;          // Start block
        uint256 toBlock;            // End block
    }

    // ============ Events ============

    event EventPublished(
        bytes32 indexed eventId,
        bytes32 indexed moduleId,
        EventCategory indexed category,
        string eventType,
        uint256 timestamp
    );

    event EventSubscribed(
        bytes32 indexed subscriptionId,
        bytes32 indexed subscriberModule,
        bytes32 indexed publisherModule,
        string eventType
    );

    event EventUnsubscribed(
        bytes32 indexed subscriptionId,
        bytes32 indexed subscriberModule
    );

    event EventCallbackExecuted(
        bytes32 indexed eventId,
        bytes32 indexed subscriptionId,
        bool success
    );

    event EventCallbackFailed(
        bytes32 indexed eventId,
        bytes32 indexed subscriptionId,
        string reason
    );

    event EventRouted(
        bytes32 indexed eventId,
        bytes32 indexed fromModule,
        bytes32 indexed toModule
    );

    // ============ Event Publishing ============

    /**
     * @notice Publish an event to the hub
     * @param category Event category
     * @param severity Event severity
     * @param eventType Event type name
     * @param eventData ABI-encoded event data
     * @return eventId Unique event identifier
     * @dev Only callable by registered modules
     */
    function publishEvent(
        EventCategory category,
        EventSeverity severity,
        string calldata eventType,
        bytes calldata eventData
    ) external returns (bytes32 eventId);

    /**
     * @notice Publish a batch of events
     * @param categories Array of event categories
     * @param severities Array of event severities
     * @param eventTypes Array of event type names
     * @param eventDataArray Array of ABI-encoded event data
     * @return eventIds Array of event identifiers
     * @dev Gas-optimized batch operation
     */
    function publishEventBatch(
        EventCategory[] calldata categories,
        EventSeverity[] calldata severities,
        string[] calldata eventTypes,
        bytes[] calldata eventDataArray
    ) external returns (bytes32[] memory eventIds);

    /**
     * @notice Emit a pre-defined system event
     * @param eventType System event type
     * @param eventData Event data
     * @return eventId Event identifier
     * @dev Convenience function for common system events
     */
    function emitSystemEvent(
        string calldata eventType,
        bytes calldata eventData
    ) external returns (bytes32 eventId);

    // ============ Event Subscription ============

    /**
     * @notice Subscribe to events from a specific module
     * @param publisherModule Module to subscribe to (bytes32(0) for all modules)
     * @param category Event category filter
     * @param eventType Event type filter (empty string for all types)
     * @param callbackAddress Address to call when event occurs
     * @param callbackSelector Function selector for callback
     * @return subscriptionId Unique subscription identifier
     * @dev Callback function should have signature: callback(bytes32 eventId, bytes calldata data)
     */
    function subscribe(
        bytes32 publisherModule,
        EventCategory category,
        string calldata eventType,
        address callbackAddress,
        bytes4 callbackSelector
    ) external returns (bytes32 subscriptionId);

    /**
     * @notice Unsubscribe from events
     * @param subscriptionId Subscription identifier
     */
    function unsubscribe(bytes32 subscriptionId) external;

    /**
     * @notice Update subscription status
     * @param subscriptionId Subscription identifier
     * @param isActive New status
     */
    function setSubscriptionActive(bytes32 subscriptionId, bool isActive) external;

    /**
     * @notice Get subscription details
     * @param subscriptionId Subscription identifier
     * @return Subscription structure
     */
    function getSubscription(bytes32 subscriptionId) external view returns (EventSubscription memory);

    /**
     * @notice Get all subscriptions for a module
     * @param moduleId Module identifier
     * @return Array of subscription IDs
     */
    function getModuleSubscriptions(bytes32 moduleId) external view returns (bytes32[] memory);

    // ============ Event Querying ============

    /**
     * @notice Get event by ID
     * @param eventId Event identifier
     * @return Event structure
     */
    function getEvent(bytes32 eventId) external view returns (HubEvent memory);

    /**
     * @notice Get events by filter
     * @param filter Event filter criteria
     * @param offset Starting index
     * @param limit Maximum number of events to return
     * @return events Array of matching events
     */
    function getEventsByFilter(
        EventFilter calldata filter,
        uint256 offset,
        uint256 limit
    ) external view returns (HubEvent[] memory events);

    /**
     * @notice Get recent events
     * @param count Number of recent events to return
     * @return events Array of recent events
     */
    function getRecentEvents(uint256 count) external view returns (HubEvent[] memory events);

    /**
     * @notice Get events for a specific module
     * @param moduleId Module identifier
     * @param offset Starting index
     * @param limit Maximum number of events to return
     * @return events Array of module events
     */
    function getModuleEvents(
        bytes32 moduleId,
        uint256 offset,
        uint256 limit
    ) external view returns (HubEvent[] memory events);

    /**
     * @notice Get events by category
     * @param category Event category
     * @param offset Starting index
     * @param limit Maximum number of events to return
     * @return events Array of events in category
     */
    function getEventsByCategory(
        EventCategory category,
        uint256 offset,
        uint256 limit
    ) external view returns (HubEvent[] memory events);

    /**
     * @notice Get events by severity
     * @param minSeverity Minimum severity level
     * @param offset Starting index
     * @param limit Maximum number of events to return
     * @return events Array of events meeting severity criteria
     */
    function getEventsBySeverity(
        EventSeverity minSeverity,
        uint256 offset,
        uint256 limit
    ) external view returns (HubEvent[] memory events);

    /**
     * @notice Get event count for a module
     * @param moduleId Module identifier
     * @return Total number of events from module
     */
    function getModuleEventCount(bytes32 moduleId) external view returns (uint256);

    // ============ Event Routing ============

    /**
     * @notice Route event to specific module
     * @param eventId Event identifier
     * @param targetModule Target module identifier
     * @dev Allows manual event routing for special cases
     */
    function routeEvent(bytes32 eventId, bytes32 targetModule) external;

    /**
     * @notice Set up automatic event routing rule
     * @param sourceModule Source module identifier
     * @param targetModule Target module identifier
     * @param eventType Event type to route (empty = all)
     * @return ruleId Routing rule identifier
     */
    function addRoutingRule(
        bytes32 sourceModule,
        bytes32 targetModule,
        string calldata eventType
    ) external returns (bytes32 ruleId);

    /**
     * @notice Remove routing rule
     * @param ruleId Routing rule identifier
     */
    function removeRoutingRule(bytes32 ruleId) external;

    // ============ Statistics & Monitoring ============

    /**
     * @notice Get event hub statistics
     * @return totalEvents Total events published
     * @return totalSubscriptions Total active subscriptions
     * @return eventsLast24h Events in last 24 hours
     * @return criticalEventsLast24h Critical events in last 24 hours
     */
    function getEventHubStats()
        external
        view
        returns (
            uint256 totalEvents,
            uint256 totalSubscriptions,
            uint256 eventsLast24h,
            uint256 criticalEventsLast24h
        );

    /**
     * @notice Get module activity statistics
     * @param moduleId Module identifier
     * @return eventCount Total events from module
     * @return lastEventTime Timestamp of last event
     * @return subscriptionCount Number of subscriptions
     */
    function getModuleActivityStats(bytes32 moduleId)
        external
        view
        returns (
            uint256 eventCount,
            uint256 lastEventTime,
            uint256 subscriptionCount
        );

    // ============ Configuration ============

    /**
     * @notice Set maximum events to store
     * @param maxEvents Maximum number of events
     * @dev Older events are pruned when limit is reached
     */
    function setMaxStoredEvents(uint256 maxEvents) external;

    /**
     * @notice Enable/disable event callbacks
     * @param enabled Whether callbacks are enabled
     * @dev Can disable in case of emergency
     */
    function setCallbacksEnabled(bool enabled) external;

    /**
     * @notice Check if callbacks are enabled
     * @return True if callbacks are enabled
     */
    function areCallbacksEnabled() external view returns (bool);
}
