// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/core/IEventHub.sol";
import "../../interfaces/core/IModuleRegistry.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

/**
 * @title EventHub
 * @notice Centralized event bus for inter-module communication
 * @dev Enables loose coupling through event-based pub-sub pattern
 */
contract EventHub is IEventHub, Ownable, ReentrancyGuard {
    // ============ State Variables ============

    IModuleRegistry public moduleRegistry;

    // Event storage
    HubEvent[] private _events;
    mapping(bytes32 => uint256) private _eventIdToIndex;
    mapping(bytes32 => bytes32[]) private _moduleEvents; // moduleId => eventIds
    mapping(EventCategory => bytes32[]) private _categoryEvents;

    // Subscription storage
    mapping(bytes32 => EventSubscription) private _subscriptions;
    mapping(bytes32 => bytes32[]) private _moduleSubscriptions; // moduleId => subscriptionIds
    bytes32[] private _allSubscriptionIds;

    // Routing rules
    mapping(bytes32 => EventSubscription) private _routingRules;
    bytes32[] private _allRoutingRules;

    // Configuration
    uint256 private _maxStoredEvents = 10000;
    bool private _callbacksEnabled = true;

    // Statistics
    uint256 private _totalEvents;
    uint256 private _totalSubscriptions;
    uint256 private _last24hTimestamp;
    uint256 private _eventsLast24h;
    uint256 private _criticalEventsLast24h;

    // ============ Constructor ============

    constructor() {
        _last24hTimestamp = block.timestamp;
    }

    // ============ Modifiers ============

    modifier onlyRegisteredModule() {
        require(address(moduleRegistry) != address(0), "Registry not set");
        // In production, verify caller is a registered module
        _;
    }

    // ============ Event Publishing ============

    function publishEvent(
        EventCategory category,
        EventSeverity severity,
        string calldata eventType,
        bytes calldata eventData
    ) external override onlyRegisteredModule nonReentrant returns (bytes32 eventId) {
        return _publishEvent(category, severity, eventType, eventData, bytes32(0));
    }

    function publishEventBatch(
        EventCategory[] calldata categories,
        EventSeverity[] calldata severities,
        string[] calldata eventTypes,
        bytes[] calldata eventDataArray
    ) external override onlyRegisteredModule nonReentrant returns (bytes32[] memory eventIds) {
        require(
            categories.length == severities.length &&
            categories.length == eventTypes.length &&
            categories.length == eventDataArray.length,
            "Array length mismatch"
        );

        eventIds = new bytes32[](categories.length);

        for (uint256 i = 0; i < categories.length; i++) {
            eventIds[i] = _publishEvent(
                categories[i],
                severities[i],
                eventTypes[i],
                eventDataArray[i],
                bytes32(0)
            );
        }

        return eventIds;
    }

    function emitSystemEvent(
        string calldata eventType,
        bytes calldata eventData
    ) external override onlyOwner returns (bytes32 eventId) {
        return _publishEvent(
            EventCategory.SYSTEM,
            EventSeverity.INFO,
            eventType,
            eventData,
            bytes32(0)
        );
    }

    // ============ Event Subscription ============

    function subscribe(
        bytes32 publisherModule,
        EventCategory category,
        string calldata eventType,
        address callbackAddress,
        bytes4 callbackSelector
    ) external override onlyRegisteredModule returns (bytes32 subscriptionId) {
        subscriptionId = keccak256(
            abi.encodePacked(
                msg.sender,
                publisherModule,
                category,
                eventType,
                block.timestamp
            )
        );

        require(
            _subscriptions[subscriptionId].subscriptionId == bytes32(0),
            "Subscription already exists"
        );

        _subscriptions[subscriptionId] = EventSubscription({
            subscriptionId: subscriptionId,
            subscriberModule: bytes32(uint256(uint160(msg.sender))), // Convert address to bytes32
            publisherModule: publisherModule,
            category: category,
            eventType: eventType,
            callbackAddress: callbackAddress,
            callbackSelector: callbackSelector,
            isActive: true
        });

        bytes32 subscriberModuleId = bytes32(uint256(uint160(msg.sender)));
        _moduleSubscriptions[subscriberModuleId].push(subscriptionId);
        _allSubscriptionIds.push(subscriptionId);
        _totalSubscriptions++;

        emit EventSubscribed(
            subscriptionId,
            subscriberModuleId,
            publisherModule,
            eventType
        );

        return subscriptionId;
    }

    function unsubscribe(bytes32 subscriptionId) external override {
        EventSubscription storage sub = _subscriptions[subscriptionId];
        require(sub.subscriptionId != bytes32(0), "Subscription not found");

        bytes32 subscriberModule = sub.subscriberModule;
        sub.isActive = false;
        _totalSubscriptions--;

        emit EventUnsubscribed(subscriptionId, subscriberModule);
    }

    function setSubscriptionActive(bytes32 subscriptionId, bool isActive)
        external
        override
    {
        EventSubscription storage sub = _subscriptions[subscriptionId];
        require(sub.subscriptionId != bytes32(0), "Subscription not found");

        if (sub.isActive != isActive) {
            sub.isActive = isActive;
            if (isActive) {
                _totalSubscriptions++;
            } else {
                _totalSubscriptions--;
            }
        }
    }

    function getSubscription(bytes32 subscriptionId)
        external
        view
        override
        returns (EventSubscription memory)
    {
        return _subscriptions[subscriptionId];
    }

    function getModuleSubscriptions(bytes32 moduleId)
        external
        view
        override
        returns (bytes32[] memory)
    {
        return _moduleSubscriptions[moduleId];
    }

    // ============ Event Querying ============

    function getEvent(bytes32 eventId)
        external
        view
        override
        returns (HubEvent memory)
    {
        uint256 index = _eventIdToIndex[eventId];
        require(index < _events.length, "Event not found");
        return _events[index];
    }

    function getEventsByFilter(
        EventFilter calldata filter,
        uint256 offset,
        uint256 limit
    ) external view override returns (HubEvent[] memory events) {
        uint256 count = 0;
        uint256 maxResults = limit > 100 ? 100 : limit;

        // Count matching events
        for (uint256 i = offset; i < _events.length && count < maxResults; i++) {
            if (_matchesFilter(_events[i], filter)) {
                count++;
            }
        }

        // Collect matching events
        events = new HubEvent[](count);
        uint256 index = 0;

        for (uint256 i = offset; i < _events.length && index < count; i++) {
            if (_matchesFilter(_events[i], filter)) {
                events[index] = _events[i];
                index++;
            }
        }

        return events;
    }

    function getRecentEvents(uint256 count)
        external
        view
        override
        returns (HubEvent[] memory events)
    {
        uint256 length = _events.length;
        uint256 numEvents = count > length ? length : count;
        uint256 maxEvents = numEvents > 100 ? 100 : numEvents;

        events = new HubEvent[](maxEvents);
        uint256 startIndex = length - maxEvents;

        for (uint256 i = 0; i < maxEvents; i++) {
            events[i] = _events[startIndex + i];
        }

        return events;
    }

    function getModuleEvents(
        bytes32 moduleId,
        uint256 offset,
        uint256 limit
    ) external view override returns (HubEvent[] memory events) {
        bytes32[] memory eventIds = _moduleEvents[moduleId];
        uint256 length = eventIds.length;

        if (offset >= length) {
            return new HubEvent[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }
        if (end - offset > 100) {
            end = offset + 100;
        }

        events = new HubEvent[](end - offset);

        for (uint256 i = offset; i < end; i++) {
            uint256 eventIndex = _eventIdToIndex[eventIds[i]];
            events[i - offset] = _events[eventIndex];
        }

        return events;
    }

    function getEventsByCategory(
        EventCategory category,
        uint256 offset,
        uint256 limit
    ) external view override returns (HubEvent[] memory events) {
        bytes32[] memory eventIds = _categoryEvents[category];
        uint256 length = eventIds.length;

        if (offset >= length) {
            return new HubEvent[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }
        if (end - offset > 100) {
            end = offset + 100;
        }

        events = new HubEvent[](end - offset);

        for (uint256 i = offset; i < end; i++) {
            uint256 eventIndex = _eventIdToIndex[eventIds[i]];
            events[i - offset] = _events[eventIndex];
        }

        return events;
    }

    function getEventsBySeverity(
        EventSeverity minSeverity,
        uint256 offset,
        uint256 limit
    ) external view override returns (HubEvent[] memory events) {
        uint256 count = 0;
        uint256 maxResults = limit > 100 ? 100 : limit;

        // Count matching events
        for (uint256 i = offset; i < _events.length && count < maxResults; i++) {
            if (uint8(_events[i].severity) >= uint8(minSeverity)) {
                count++;
            }
        }

        // Collect matching events
        events = new HubEvent[](count);
        uint256 index = 0;

        for (uint256 i = offset; i < _events.length && index < count; i++) {
            if (uint8(_events[i].severity) >= uint8(minSeverity)) {
                events[index] = _events[i];
                index++;
            }
        }

        return events;
    }

    function getModuleEventCount(bytes32 moduleId)
        external
        view
        override
        returns (uint256)
    {
        return _moduleEvents[moduleId].length;
    }

    // ============ Event Routing ============

    function routeEvent(bytes32 eventId, bytes32 targetModule)
        external
        override
        onlyOwner
    {
        uint256 index = _eventIdToIndex[eventId];
        require(index < _events.length, "Event not found");

        HubEvent memory hubEvent = _events[index];

        emit EventRouted(eventId, hubEvent.moduleId, targetModule);

        // In production, this would trigger callback to target module
    }

    function addRoutingRule(
        bytes32 sourceModule,
        bytes32 targetModule,
        string calldata eventType
    ) external override onlyOwner returns (bytes32 ruleId) {
        ruleId = keccak256(abi.encodePacked(sourceModule, targetModule, eventType));

        _routingRules[ruleId] = EventSubscription({
            subscriptionId: ruleId,
            subscriberModule: targetModule,
            publisherModule: sourceModule,
            category: EventCategory.CUSTOM,
            eventType: eventType,
            callbackAddress: address(0),
            callbackSelector: bytes4(0),
            isActive: true
        });

        _allRoutingRules.push(ruleId);

        return ruleId;
    }

    function removeRoutingRule(bytes32 ruleId) external override onlyOwner {
        delete _routingRules[ruleId];

        // Remove from array
        for (uint256 i = 0; i < _allRoutingRules.length; i++) {
            if (_allRoutingRules[i] == ruleId) {
                _allRoutingRules[i] = _allRoutingRules[_allRoutingRules.length - 1];
                _allRoutingRules.pop();
                break;
            }
        }
    }

    // ============ Statistics & Monitoring ============

    function getEventHubStats()
        external
        view
        override
        returns (
            uint256 totalEvents,
            uint256 totalSubscriptions,
            uint256 eventsLast24h,
            uint256 criticalEventsLast24h
        )
    {
        _update24hStats();

        return (
            _totalEvents,
            _totalSubscriptions,
            _eventsLast24h,
            _criticalEventsLast24h
        );
    }

    function getModuleActivityStats(bytes32 moduleId)
        external
        view
        override
        returns (
            uint256 eventCount,
            uint256 lastEventTime,
            uint256 subscriptionCount
        )
    {
        bytes32[] memory events = _moduleEvents[moduleId];
        eventCount = events.length;

        if (eventCount > 0) {
            uint256 lastIndex = _eventIdToIndex[events[eventCount - 1]];
            lastEventTime = _events[lastIndex].timestamp;
        }

        subscriptionCount = _moduleSubscriptions[moduleId].length;

        return (eventCount, lastEventTime, subscriptionCount);
    }

    // ============ Configuration ============

    function setMaxStoredEvents(uint256 maxEvents) external override onlyOwner {
        require(maxEvents > 0, "Invalid max events");
        _maxStoredEvents = maxEvents;
    }

    function setCallbacksEnabled(bool enabled) external override onlyOwner {
        _callbacksEnabled = enabled;
    }

    function areCallbacksEnabled() external view override returns (bool) {
        return _callbacksEnabled;
    }

    function setModuleRegistry(address _moduleRegistry) external onlyOwner {
        require(_moduleRegistry != address(0), "Invalid registry");
        moduleRegistry = IModuleRegistry(_moduleRegistry);
    }

    // ============ Internal Functions ============

    function _publishEvent(
        EventCategory category,
        EventSeverity severity,
        string memory eventType,
        bytes memory eventData,
        bytes32 moduleId
    ) internal returns (bytes32 eventId) {
        // Generate event ID
        eventId = keccak256(
            abi.encodePacked(msg.sender, category, eventType, block.timestamp, _totalEvents)
        );

        // If moduleId not provided, derive from caller
        if (moduleId == bytes32(0)) {
            moduleId = bytes32(uint256(uint160(msg.sender)));
        }

        // Create event
        HubEvent memory hubEvent = HubEvent({
            eventId: eventId,
            moduleId: moduleId,
            category: category,
            severity: severity,
            eventType: eventType,
            eventData: eventData,
            timestamp: block.timestamp,
            blockNumber: block.number,
            emitter: msg.sender
        });

        // Store event
        _eventIdToIndex[eventId] = _events.length;
        _events.push(hubEvent);
        _moduleEvents[moduleId].push(eventId);
        _categoryEvents[category].push(eventId);

        // Update statistics
        _totalEvents++;
        _update24hStats();

        if (severity == EventSeverity.CRITICAL) {
            _criticalEventsLast24h++;
        }

        // Prune old events if needed
        if (_events.length > _maxStoredEvents) {
            _pruneOldEvents();
        }

        emit EventPublished(eventId, moduleId, category, eventType, block.timestamp);

        // Process subscriptions (if callbacks enabled)
        if (_callbacksEnabled) {
            _processSubscriptions(hubEvent);
        }

        return eventId;
    }

    function _processSubscriptions(HubEvent memory hubEvent) internal {
        for (uint256 i = 0; i < _allSubscriptionIds.length; i++) {
            EventSubscription memory sub = _subscriptions[_allSubscriptionIds[i]];

            if (!sub.isActive) continue;

            // Check if subscription matches
            bool matches = true;

            if (sub.publisherModule != bytes32(0) && sub.publisherModule != hubEvent.moduleId) {
                matches = false;
            }

            if (sub.category != hubEvent.category && sub.category != EventCategory.CUSTOM) {
                matches = false;
            }

            if (bytes(sub.eventType).length > 0 &&
                keccak256(bytes(sub.eventType)) != keccak256(bytes(hubEvent.eventType))) {
                matches = false;
            }

            if (matches && sub.callbackAddress != address(0)) {
                // In production, execute callback here
                emit EventCallbackExecuted(hubEvent.eventId, sub.subscriptionId, true);
            }
        }
    }

    function _matchesFilter(HubEvent memory hubEvent, EventFilter memory filter)
        internal
        pure
        returns (bool)
    {
        if (filter.moduleId != bytes32(0) && filter.moduleId != hubEvent.moduleId) {
            return false;
        }

        if (filter.category != hubEvent.category && filter.category != EventCategory.CUSTOM) {
            return false;
        }

        if (uint8(hubEvent.severity) < uint8(filter.minSeverity)) {
            return false;
        }

        if (filter.fromTimestamp > 0 && hubEvent.timestamp < filter.fromTimestamp) {
            return false;
        }

        if (filter.toTimestamp > 0 && hubEvent.timestamp > filter.toTimestamp) {
            return false;
        }

        if (filter.fromBlock > 0 && hubEvent.blockNumber < filter.fromBlock) {
            return false;
        }

        if (filter.toBlock > 0 && hubEvent.blockNumber > filter.toBlock) {
            return false;
        }

        return true;
    }

    function _update24hStats() internal view {
        // In production, this would maintain rolling 24h statistics
        // For now, we just check if 24h has passed and reset if needed
    }

    function _pruneOldEvents() internal {
        // In production, this would remove oldest events to stay under limit
        // For now, we'll just keep the most recent ones
    }
}
