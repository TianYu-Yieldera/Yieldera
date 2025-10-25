// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/modules/rwa/IMarketDataProvider.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title MarketDataProvider
 * @notice Provides market statistics and analytics
 * @dev Tracks price history, volume, and order book data
 */
contract MarketDataProvider is IMarketDataProvider, Ownable {
    // ============ Constants ============

    uint256 private constant TWENTY_FOUR_HOURS = 24 hours;
    uint256 private constant BASIS_POINTS = 10000;

    // ============ State Variables ============

    address public rwaModule; // Main coordinator contract

    // Storage for market data (using Diamond Storage pattern)
    bytes32 private constant MARKET_DATA_STORAGE_POSITION = keccak256("market.data.storage");

    struct MarketDataStorage {
        MarketStats stats;
        PricePoint[] priceHistory;
        mapping(uint256 => OrderBookLevel) buyLevels;
        mapping(uint256 => OrderBookLevel) sellLevels;
        uint256[] buyPriceLevels;
        uint256[] sellPriceLevels;
        uint256 maxPriceHistory;
        uint256 price24hAgo;
        uint256 timestamp24hAgo;
    }

    // ============ Modifiers ============

    modifier onlyRWAModule() {
        require(msg.sender == rwaModule, "Only RWA module");
        _;
    }

    // ============ Constructor ============

    constructor() {
        MarketDataStorage storage mds = _getStorage();
        mds.maxPriceHistory = 1000; // Keep last 1000 price points
    }

    // ============ Admin Functions ============

    function setRWAModule(address _rwaModule) external onlyOwner {
        require(_rwaModule != address(0), "Invalid address");
        rwaModule = _rwaModule;
    }

    // ============ Internal Storage ============

    function _getStorage() private pure returns (MarketDataStorage storage mds) {
        bytes32 position = MARKET_DATA_STORAGE_POSITION;
        assembly {
            mds.slot := position
        }
    }

    // ============ IMarketDataProvider Implementation ============

    function updateStats(uint256 price, uint256 volume) external override onlyRWAModule {
        MarketDataStorage storage mds = _getStorage();

        // Update last price
        uint256 oldPrice = mds.stats.lastPrice;
        mds.stats.lastPrice = price;

        // Update 24h high/low
        if (price > mds.stats.highPrice24h || mds.stats.highPrice24h == 0) {
            mds.stats.highPrice24h = price;
        }
        if (price < mds.stats.lowPrice24h || mds.stats.lowPrice24h == 0) {
            mds.stats.lowPrice24h = price;
        }

        // Update volume
        mds.stats.totalVolume += volume;
        mds.stats.totalTrades++;
        mds.stats.volumeLast24h += volume;

        // Add to price history
        mds.priceHistory.push(
            PricePoint({price: price, timestamp: block.timestamp, volume: volume})
        );

        // Trim price history if needed
        if (mds.priceHistory.length > mds.maxPriceHistory) {
            // Remove oldest entry
            for (uint256 i = 0; i < mds.priceHistory.length - 1; i++) {
                mds.priceHistory[i] = mds.priceHistory[i + 1];
            }
            mds.priceHistory.pop();
        }

        // Update 24h reference if needed
        if (block.timestamp - mds.timestamp24hAgo >= TWENTY_FOUR_HOURS) {
            mds.price24hAgo = oldPrice;
            mds.timestamp24hAgo = block.timestamp;
        }

        emit PriceUpdated(price, block.timestamp);
        emit VolumeUpdated(mds.stats.volumeLast24h);
    }

    function getMarketStats() external view override returns (MarketStats memory stats) {
        MarketDataStorage storage mds = _getStorage();
        return mds.stats;
    }

    function getPriceHistory(uint256 count)
        external
        view
        override
        returns (PricePoint[] memory priceHistory)
    {
        MarketDataStorage storage mds = _getStorage();

        uint256 length = mds.priceHistory.length;
        if (count > length) count = length;

        priceHistory = new PricePoint[](count);
        for (uint256 i = 0; i < count; i++) {
            priceHistory[i] = mds.priceHistory[length - count + i];
        }

        return priceHistory;
    }

    function getOrderBook(uint256 levels)
        external
        view
        override
        returns (OrderBookLevel[] memory buyLevels, OrderBookLevel[] memory sellLevels)
    {
        MarketDataStorage storage mds = _getStorage();

        uint256 buyCount = mds.buyPriceLevels.length;
        uint256 sellCount = mds.sellPriceLevels.length;

        if (levels > buyCount) levels = buyCount;
        if (levels > sellCount) levels = sellCount;

        buyLevels = new OrderBookLevel[](levels);
        sellLevels = new OrderBookLevel[](levels);

        for (uint256 i = 0; i < levels; i++) {
            buyLevels[i] = mds.buyLevels[mds.buyPriceLevels[i]];
            sellLevels[i] = mds.sellLevels[mds.sellPriceLevels[i]];
        }

        return (buyLevels, sellLevels);
    }

    function get24hVolume() external view override returns (uint256 volume24h) {
        MarketDataStorage storage mds = _getStorage();
        return mds.stats.volumeLast24h;
    }

    function get24hPriceChange() external view override returns (int256 priceChange24h) {
        MarketDataStorage storage mds = _getStorage();

        if (mds.price24hAgo == 0) return 0;

        // Calculate percentage change in basis points
        int256 priceDiff = int256(mds.stats.lastPrice) - int256(mds.price24hAgo);
        priceChange24h = (priceDiff * int256(BASIS_POINTS)) / int256(mds.price24hAgo);

        return priceChange24h;
    }

    function getSpread() external view override returns (uint256 spread) {
        MarketDataStorage storage mds = _getStorage();

        if (mds.buyPriceLevels.length == 0 || mds.sellPriceLevels.length == 0) {
            return 0;
        }

        // Get best bid and ask
        uint256 bestBid = mds.buyLevels[mds.buyPriceLevels[0]].price;
        uint256 bestAsk = mds.sellLevels[mds.sellPriceLevels[0]].price;

        return bestAsk > bestBid ? bestAsk - bestBid : 0;
    }

    function getMarketDepth(uint256 priceRange)
        external
        view
        override
        returns (uint256 buyDepth, uint256 sellDepth)
    {
        MarketDataStorage storage mds = _getStorage();

        if (mds.stats.lastPrice == 0) return (0, 0);

        uint256 minPrice = (mds.stats.lastPrice * (BASIS_POINTS - priceRange)) / BASIS_POINTS;
        uint256 maxPrice = (mds.stats.lastPrice * (BASIS_POINTS + priceRange)) / BASIS_POINTS;

        // Sum buy volume within range
        for (uint256 i = 0; i < mds.buyPriceLevels.length; i++) {
            uint256 price = mds.buyPriceLevels[i];
            if (price >= minPrice && price <= maxPrice) {
                buyDepth += mds.buyLevels[price].totalAmount;
            }
        }

        // Sum sell volume within range
        for (uint256 i = 0; i < mds.sellPriceLevels.length; i++) {
            uint256 price = mds.sellPriceLevels[i];
            if (price >= minPrice && price <= maxPrice) {
                sellDepth += mds.sellLevels[price].totalAmount;
            }
        }

        return (buyDepth, sellDepth);
    }

    function getTWAP(uint256 duration) external view override returns (uint256 twap) {
        MarketDataStorage storage mds = _getStorage();

        if (mds.priceHistory.length == 0) return 0;
        if (duration == 0) return mds.stats.lastPrice;

        uint256 cutoffTime = block.timestamp - duration;
        uint256 weightedSum = 0;
        uint256 totalTime = 0;

        for (uint256 i = mds.priceHistory.length; i > 0; i--) {
            PricePoint memory point = mds.priceHistory[i - 1];

            if (point.timestamp < cutoffTime) break;

            uint256 timeDiff;
            if (i == mds.priceHistory.length) {
                timeDiff = block.timestamp - point.timestamp;
            } else {
                timeDiff = mds.priceHistory[i].timestamp - point.timestamp;
            }

            weightedSum += point.price * timeDiff;
            totalTime += timeDiff;
        }

        return totalTime > 0 ? weightedSum / totalTime : mds.stats.lastPrice;
    }

    // ============ Additional Functions ============

    /**
     * @notice Update order book level
     * @param price Price level
     * @param isBuy True for buy side
     * @param totalAmount Total amount at this price
     * @param orderCount Number of orders at this price
     */
    function updateOrderBookLevel(uint256 price, bool isBuy, uint256 totalAmount, uint256 orderCount)
        external
        onlyRWAModule
    {
        MarketDataStorage storage mds = _getStorage();

        if (isBuy) {
            OrderBookLevel storage level = mds.buyLevels[price];
            bool isNew = level.price == 0;

            level.price = price;
            level.totalAmount = totalAmount;
            level.orderCount = orderCount;

            if (isNew && totalAmount > 0) {
                mds.buyPriceLevels.push(price);
                _sortPriceLevels(mds.buyPriceLevels, true);
            }
        } else {
            OrderBookLevel storage level = mds.sellLevels[price];
            bool isNew = level.price == 0;

            level.price = price;
            level.totalAmount = totalAmount;
            level.orderCount = orderCount;

            if (isNew && totalAmount > 0) {
                mds.sellPriceLevels.push(price);
                _sortPriceLevels(mds.sellPriceLevels, false);
            }
        }
    }

    /**
     * @notice Update active orders count
     * @param count New active orders count
     */
    function updateActiveOrders(uint256 count) external onlyRWAModule {
        MarketDataStorage storage mds = _getStorage();
        mds.stats.activeOrders = count;
    }

    /**
     * @notice Refresh statistics (cleanup old data)
     */
    function refreshStats() external onlyRWAModule {
        MarketDataStorage storage mds = _getStorage();

        // Reset 24h volume if window has passed
        if (block.timestamp - mds.timestamp24hAgo >= TWENTY_FOUR_HOURS) {
            mds.stats.volumeLast24h = 0;
            mds.stats.highPrice24h = mds.stats.lastPrice;
            mds.stats.lowPrice24h = mds.stats.lastPrice;
        }

        emit StatsRefreshed(block.timestamp);
    }

    // ============ Internal Helper Functions ============

    function _sortPriceLevels(uint256[] storage levels, bool descending) internal {
        // Simple bubble sort (gas intensive, but ok for small arrays)
        for (uint256 i = 0; i < levels.length; i++) {
            for (uint256 j = i + 1; j < levels.length; j++) {
                if (descending) {
                    if (levels[i] < levels[j]) {
                        (levels[i], levels[j]) = (levels[j], levels[i]);
                    }
                } else {
                    if (levels[i] > levels[j]) {
                        (levels[i], levels[j]) = (levels[j], levels[i]);
                    }
                }
            }
        }
    }

    function setMaxPriceHistory(uint256 max) external onlyOwner {
        MarketDataStorage storage mds = _getStorage();
        mds.maxPriceHistory = max;
    }
}
