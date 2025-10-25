// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IMarketDataProvider
 * @notice Interface for market statistics and data
 * @dev Provides market metrics and analytics
 */
interface IMarketDataProvider {
    // ============ Structs ============

    struct MarketStats {
        uint256 totalVolume;
        uint256 totalTrades;
        uint256 lastPrice;
        uint256 highPrice24h;
        uint256 lowPrice24h;
        uint256 volumeLast24h;
        uint256 activeOrders;
    }

    struct PricePoint {
        uint256 price;
        uint256 timestamp;
        uint256 volume;
    }

    struct OrderBookLevel {
        uint256 price;
        uint256 totalAmount;
        uint256 orderCount;
    }

    // ============ Events ============

    event PriceUpdated(uint256 newPrice, uint256 timestamp);
    event VolumeUpdated(uint256 volume24h);
    event StatsRefreshed(uint256 timestamp);

    // ============ Functions ============

    /**
     * @notice Update market statistics after a trade
     * @param price Trade price
     * @param volume Trade volume
     */
    function updateStats(uint256 price, uint256 volume) external;

    /**
     * @notice Get current market statistics
     * @return stats Market statistics
     */
    function getMarketStats() external view returns (MarketStats memory stats);

    /**
     * @notice Get price history
     * @param count Number of price points to return
     * @return priceHistory Array of price points
     */
    function getPriceHistory(uint256 count) external view returns (PricePoint[] memory priceHistory);

    /**
     * @notice Get current order book (aggregated by price level)
     * @param levels Number of price levels to return
     * @return buyLevels Buy side order book
     * @return sellLevels Sell side order book
     */
    function getOrderBook(uint256 levels)
        external
        view
        returns (OrderBookLevel[] memory buyLevels, OrderBookLevel[] memory sellLevels);

    /**
     * @notice Get 24-hour trading volume
     * @return volume24h Volume in last 24 hours
     */
    function get24hVolume() external view returns (uint256 volume24h);

    /**
     * @notice Get 24-hour price change
     * @return priceChange24h Price change percentage (in basis points)
     */
    function get24hPriceChange() external view returns (int256 priceChange24h);

    /**
     * @notice Get current spread
     * @return spread Difference between best bid and ask
     */
    function getSpread() external view returns (uint256 spread);

    /**
     * @notice Get market depth
     * @param priceRange Price range to analyze (in basis points from mid price)
     * @return buyDepth Total buy volume within range
     * @return sellDepth Total sell volume within range
     */
    function getMarketDepth(uint256 priceRange)
        external
        view
        returns (uint256 buyDepth, uint256 sellDepth);

    /**
     * @notice Get time-weighted average price
     * @param duration Duration in seconds
     * @return twap Time-weighted average price
     */
    function getTWAP(uint256 duration) external view returns (uint256 twap);
}
