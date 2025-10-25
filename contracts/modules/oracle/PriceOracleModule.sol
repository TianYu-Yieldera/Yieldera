// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/core/IPriceOracleModule.sol";
import "../../plugins/core/BaseModule.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title PriceOracleModule
 * @notice Pluggable price oracle module with multi-source support
 * @dev Extends BaseModule and implements IPriceOracleModule
 */
contract PriceOracleModule is IPriceOracleModule, BaseModule, Ownable {
    // ============ Constants ============

    bytes32 public constant MODULE_ID = keccak256("PRICE_ORACLE_MODULE");
    string public constant MODULE_NAME = "PriceOracleModule";
    string public constant MODULE_VERSION = "1.0.0";

    // ============ State Variables ============

    // Price feed storage
    mapping(address => PriceFeed[]) private _tokenPriceFeeds;
    mapping(address => bool) private _hasActiveFeed;

    // Supported tokens
    address[] private _supportedTokens;

    // Configuration
    uint256 public stalenessThreshold = 3600; // 1 hour default
    uint256 public minConfidence = 9000; // 90% minimum confidence
    uint256 public maxDeviation = 500; // 5% max deviation

    // Price update authorization
    mapping(address => bool) public isPriceUpdater;

    // Statistics
    uint256 private _totalTokens;
    uint256 private _activeFeeds;
    uint256 private _stalePrices;
    uint256 private _lastGlobalUpdate;

    // ============ Constructor ============

    constructor() BaseModule(MODULE_ID, MODULE_NAME, MODULE_VERSION) {
        isPriceUpdater[msg.sender] = true;
    }

    // ============ Modifiers ============

    modifier onlyPriceUpdater() {
        require(isPriceUpdater[msg.sender] || msg.sender == owner(), "Not authorized");
        _;
    }

    // ============ BaseModule Overrides ============

    function getDependencies() external pure override returns (bytes32[] memory) {
        // No dependencies for oracle module
        return new bytes32[](0);
    }

    function healthCheck()
        external
        view
        override(IModule, BaseModule)
        returns (bool healthy, string memory message)
    {
        // Check base module health first
        (bool baseHealthy, string memory baseMessage) = BaseModule.healthCheck();
        if (!baseHealthy) {
            return (false, baseMessage);
        }

        // Check if we have any active feeds
        if (_activeFeeds == 0) {
            return (false, "No active price feeds");
        }

        // Check for too many stale prices
        if (_stalePrices > _activeFeeds / 2) {
            return (false, "Too many stale prices");
        }

        return (true, "Oracle healthy");
    }

    // ============ IPriceOracleModule Implementation ============

    function getLatestPrice(address token)
        external
        view
        override
        whenNotPaused
        returns (uint256 price, uint8 decimals, uint256 timestamp)
    {
        require(_hasActiveFeed[token], "No price feed for token");

        PriceFeed[] storage feeds = _tokenPriceFeeds[token];
        require(feeds.length > 0, "No feeds available");

        // Use first active feed for now
        for (uint256 i = 0; i < feeds.length; i++) {
            if (feeds[i].isActive) {
                require(!_isPriceStale(feeds[i].lastUpdateTime), "Price is stale");
                return (feeds[i].lastPrice, feeds[i].decimals, feeds[i].lastUpdateTime);
            }
        }

        revert("No active feed found");
    }

    function getPriceWithConfidence(address token)
        external
        view
        override
        returns (uint256 price, uint256 confidence, uint256 timestamp)
    {
        require(_hasActiveFeed[token], "No price feed for token");

        PriceFeed[] storage feeds = _tokenPriceFeeds[token];
        uint256 activeCount = 0;
        uint256 totalPrice = 0;
        uint256 latestTime = 0;

        // Aggregate from all active feeds
        for (uint256 i = 0; i < feeds.length; i++) {
            if (feeds[i].isActive && !_isPriceStale(feeds[i].lastUpdateTime)) {
                totalPrice += feeds[i].lastPrice;
                activeCount++;
                if (feeds[i].lastUpdateTime > latestTime) {
                    latestTime = feeds[i].lastUpdateTime;
                }
            }
        }

        require(activeCount > 0, "No active feeds");

        price = totalPrice / activeCount;
        confidence = activeCount >= 3 ? 10000 : (activeCount * 5000); // Higher confidence with more sources
        timestamp = latestTime;

        return (price, confidence, timestamp);
    }

    function getAggregatedPrice(address token)
        external
        view
        override
        returns (AggregatedPrice memory)
    {
        require(_hasActiveFeed[token], "No price feed for token");

        PriceFeed[] storage feeds = _tokenPriceFeeds[token];
        uint256[] memory prices = new uint256[](feeds.length);
        uint256 activeCount = 0;
        uint256 latestTime = 0;

        // Collect active prices
        for (uint256 i = 0; i < feeds.length; i++) {
            if (feeds[i].isActive && !_isPriceStale(feeds[i].lastUpdateTime)) {
                prices[activeCount] = feeds[i].lastPrice;
                activeCount++;
                if (feeds[i].lastUpdateTime > latestTime) {
                    latestTime = feeds[i].lastUpdateTime;
                }
            }
        }

        require(activeCount > 0, "No active feeds");

        // Calculate median and deviation
        uint256 median = _calculateMedian(prices, activeCount);
        uint256 maxDev = _calculateMaxDeviation(prices, activeCount, median);

        return AggregatedPrice({
            price: median,
            confidence: activeCount >= 3 ? 10000 : (activeCount * 5000),
            timestamp: latestTime,
            deviation: maxDev,
            sourcesUsed: activeCount
        });
    }

    function getPriceInTermsOf(address tokenA, address tokenB)
        external
        view
        override
        returns (uint256)
    {
        (uint256 priceA, , ) = this.getLatestPrice(tokenA);
        (uint256 priceB, , ) = this.getLatestPrice(tokenB);

        require(priceB > 0, "Invalid quote token price");

        return (priceA * 1e8) / priceB;
    }

    function tryGetPrice(address token)
        external
        view
        override
        returns (bool success, uint256 price, uint256 timestamp)
    {
        if (!_hasActiveFeed[token]) {
            return (false, 0, 0);
        }

        PriceFeed[] storage feeds = _tokenPriceFeeds[token];
        for (uint256 i = 0; i < feeds.length; i++) {
            if (feeds[i].isActive && !_isPriceStale(feeds[i].lastUpdateTime)) {
                return (true, feeds[i].lastPrice, feeds[i].lastUpdateTime);
            }
        }

        return (false, 0, 0);
    }

    // ============ Price Feed Management ============

    function addPriceFeed(
        address token,
        address feedAddress,
        PriceSource source,
        uint256 heartbeat
    ) external override onlyOwner {
        require(token != address(0), "Invalid token");
        require(feedAddress != address(0), "Invalid feed");

        PriceFeed memory newFeed = PriceFeed({
            feedAddress: feedAddress,
            source: source,
            decimals: 8,
            heartbeat: heartbeat,
            lastPrice: 1e8, // Default to $1
            lastUpdateTime: block.timestamp,
            isActive: true,
            minPrice: 0,
            maxPrice: type(uint256).max
        });

        _tokenPriceFeeds[token].push(newFeed);

        if (!_hasActiveFeed[token]) {
            _hasActiveFeed[token] = true;
            _supportedTokens.push(token);
            _totalTokens++;
        }

        _activeFeeds++;

        emit PriceFeedAdded(token, feedAddress, source);
    }

    function removePriceFeed(address token, address feedAddress)
        external
        override
        onlyOwner
    {
        PriceFeed[] storage feeds = _tokenPriceFeeds[token];

        for (uint256 i = 0; i < feeds.length; i++) {
            if (feeds[i].feedAddress == feedAddress) {
                if (feeds[i].isActive) {
                    _activeFeeds--;
                }

                // Remove by swapping with last element
                feeds[i] = feeds[feeds.length - 1];
                feeds.pop();

                emit PriceFeedRemoved(token, feedAddress);
                break;
            }
        }

        // Check if token still has feeds
        if (feeds.length == 0) {
            _hasActiveFeed[token] = false;
        }
    }

    function setPriceFeedStatus(address token, address feedAddress, bool isActive)
        external
        override
        onlyOwner
    {
        PriceFeed[] storage feeds = _tokenPriceFeeds[token];

        for (uint256 i = 0; i < feeds.length; i++) {
            if (feeds[i].feedAddress == feedAddress) {
                if (feeds[i].isActive != isActive) {
                    feeds[i].isActive = isActive;
                    if (isActive) {
                        _activeFeeds++;
                    } else {
                        _activeFeeds--;
                    }
                    emit PriceFeedStatusChanged(token, feedAddress, isActive);
                }
                break;
            }
        }
    }

    function getPriceFeeds(address token)
        external
        view
        override
        returns (PriceFeed[] memory)
    {
        return _tokenPriceFeeds[token];
    }

    function hasPriceFeed(address token)
        external
        view
        override
        returns (bool)
    {
        return _hasActiveFeed[token];
    }

    // ============ Price Updates ============

    function updatePrice(address token, uint256 price)
        external
        override
        onlyPriceUpdater
    {
        require(_hasActiveFeed[token], "No feed for token");
        require(price > 0, "Invalid price");

        PriceFeed[] storage feeds = _tokenPriceFeeds[token];

        // Update first manual/custom feed
        for (uint256 i = 0; i < feeds.length; i++) {
            if ((feeds[i].source == PriceSource.MANUAL ||
                 feeds[i].source == PriceSource.CUSTOM) &&
                feeds[i].isActive) {

                // Validate circuit breaker
                require(
                    price >= feeds[i].minPrice && price <= feeds[i].maxPrice,
                    "Price outside circuit breaker limits"
                );

                feeds[i].lastPrice = price;
                feeds[i].lastUpdateTime = block.timestamp;
                _lastGlobalUpdate = block.timestamp;

                emit PriceUpdated(token, price, block.timestamp, feeds[i].source);
                break;
            }
        }
    }

    function batchUpdatePrices(address[] calldata tokens, uint256[] calldata prices)
        external
        override
        onlyPriceUpdater
    {
        require(tokens.length == prices.length, "Length mismatch");

        for (uint256 i = 0; i < tokens.length; i++) {
            this.updatePrice(tokens[i], prices[i]);
        }
    }

    function refreshPrice(address token) external override {
        // In production, this would trigger oracle refresh
        // For now, just emit event
        require(_hasActiveFeed[token], "No feed for token");
        emit PriceUpdated(token, 0, block.timestamp, PriceSource.CUSTOM);
    }

    // ============ Validation ============

    function isPriceStale(address token)
        external
        view
        override
        returns (bool)
    {
        if (!_hasActiveFeed[token]) {
            return true;
        }

        PriceFeed[] storage feeds = _tokenPriceFeeds[token];
        for (uint256 i = 0; i < feeds.length; i++) {
            if (feeds[i].isActive) {
                if (!_isPriceStale(feeds[i].lastUpdateTime)) {
                    return false; // At least one fresh feed
                }
            }
        }

        return true; // All feeds are stale
    }

    function getTimeSinceUpdate(address token)
        external
        view
        override
        returns (uint256)
    {
        require(_hasActiveFeed[token], "No feed for token");

        PriceFeed[] storage feeds = _tokenPriceFeeds[token];
        uint256 latestUpdate = 0;

        for (uint256 i = 0; i < feeds.length; i++) {
            if (feeds[i].isActive && feeds[i].lastUpdateTime > latestUpdate) {
                latestUpdate = feeds[i].lastUpdateTime;
            }
        }

        return block.timestamp - latestUpdate;
    }

    function validatePrice(address token, uint256 price)
        external
        view
        override
        returns (bool valid)
    {
        if (!_hasActiveFeed[token]) {
            return false;
        }

        PriceFeed[] storage feeds = _tokenPriceFeeds[token];

        for (uint256 i = 0; i < feeds.length; i++) {
            if (feeds[i].isActive) {
                if (price < feeds[i].minPrice || price > feeds[i].maxPrice) {
                    return false;
                }
            }
        }

        return true;
    }

    function setCircuitBreaker(address token, uint256 minPrice, uint256 maxPrice)
        external
        override
        onlyOwner
    {
        require(minPrice < maxPrice, "Invalid range");

        PriceFeed[] storage feeds = _tokenPriceFeeds[token];

        for (uint256 i = 0; i < feeds.length; i++) {
            feeds[i].minPrice = minPrice;
            feeds[i].maxPrice = maxPrice;
        }

        emit CircuitBreakerTriggered(token, 0, minPrice, maxPrice);
    }

    // ============ Historical Data (Stub) ============

    function getHistoricalPrice(address token, uint256 timestamp)
        external
        view
        override
        returns (uint256 price, uint256 actualTimestamp)
    {
        // Stub: return current price
        (price, , actualTimestamp) = this.getLatestPrice(token);
        return (price, actualTimestamp);
    }

    function getPriceChange(address token, uint256 period)
        external
        view
        override
        returns (int256 priceChange, int256 percentChange)
    {
        // Stub: returns 0
        return (0, 0);
    }

    function getTWAP(address token, uint256 period)
        external
        view
        override
        returns (uint256)
    {
        // Stub: return current price
        (uint256 price, , ) = this.getLatestPrice(token);
        return price;
    }

    // ============ Configuration ============

    function setStalenessThreshold(uint256 threshold)
        external
        override
        onlyOwner
    {
        require(threshold > 0, "Invalid threshold");
        stalenessThreshold = threshold;
        emit OracleConfigUpdated("stalenessThreshold", threshold);
    }

    function setMinConfidence(uint256 minConf)
        external
        override
        onlyOwner
    {
        require(minConf <= 10000, "Invalid confidence");
        minConfidence = minConf;
        emit OracleConfigUpdated("minConfidence", minConf);
    }

    function setMaxDeviation(uint256 maxDev)
        external
        override
        onlyOwner
    {
        require(maxDev <= 10000, "Invalid deviation");
        maxDeviation = maxDev;
        emit OracleConfigUpdated("maxDeviation", maxDev);
    }

    function getOracleConfig()
        external
        view
        override
        returns (
            uint256 _stalenessThreshold,
            uint256 _minConfidence,
            uint256 _maxDeviation
        )
    {
        return (stalenessThreshold, minConfidence, maxDeviation);
    }

    // ============ Statistics ============

    function getOracleStats()
        external
        view
        override
        returns (
            uint256 totalTokens,
            uint256 activeFeeds,
            uint256 stalePrices,
            uint256 lastGlobalUpdate
        )
    {
        return (_totalTokens, _activeFeeds, _stalePrices, _lastGlobalUpdate);
    }

    function getSupportedTokens()
        external
        view
        override
        returns (address[] memory)
    {
        return _supportedTokens;
    }

    // ============ Admin Functions ============

    function setPriceUpdater(address updater, bool authorized) external onlyOwner {
        isPriceUpdater[updater] = authorized;
    }

    // ============ Internal Helper Functions ============

    function _isPriceStale(uint256 lastUpdate) internal view returns (bool) {
        return block.timestamp > lastUpdate + stalenessThreshold;
    }

    function _calculateMedian(uint256[] memory prices, uint256 length)
        internal
        pure
        returns (uint256)
    {
        if (length == 1) return prices[0];
        if (length == 2) return (prices[0] + prices[1]) / 2;

        // Simple bubble sort for small arrays
        for (uint256 i = 0; i < length - 1; i++) {
            for (uint256 j = 0; j < length - i - 1; j++) {
                if (prices[j] > prices[j + 1]) {
                    uint256 temp = prices[j];
                    prices[j] = prices[j + 1];
                    prices[j + 1] = temp;
                }
            }
        }

        if (length % 2 == 0) {
            return (prices[length / 2 - 1] + prices[length / 2]) / 2;
        } else {
            return prices[length / 2];
        }
    }

    function _calculateMaxDeviation(
        uint256[] memory prices,
        uint256 length,
        uint256 median
    ) internal pure returns (uint256) {
        uint256 maxDev = 0;

        for (uint256 i = 0; i < length; i++) {
            uint256 dev = prices[i] > median
                ? ((prices[i] - median) * 10000) / median
                : ((median - prices[i]) * 10000) / median;

            if (dev > maxDev) {
                maxDev = dev;
            }
        }

        return maxDev;
    }

    // ============ Override Required Functions ============

    function pause() external override(IModule, BaseModule) onlyOwner {
        BaseModule.pause();
    }

    function unpause() external override(IModule, BaseModule) onlyOwner {
        BaseModule.unpause();
    }
}
