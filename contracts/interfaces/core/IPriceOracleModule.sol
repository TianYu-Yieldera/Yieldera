// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IModule.sol";

/**
 * @title IPriceOracleModule
 * @notice Standard interface for price oracle modules
 * @dev Extends IModule with price feed functionality
 *
 * Purpose:
 * - Provide reliable price data
 * - Support multiple price sources
 * - Aggregation and validation
 * - Staleness detection
 */
interface IPriceOracleModule is IModule {
    /**
     * @notice Price source types
     */
    enum PriceSource {
        CHAINLINK,          // Chainlink price feed
        UNISWAP_V3,        // Uniswap V3 TWAP
        CUSTOM,            // Custom oracle
        MANUAL,            // Manual price setting
        AGGREGATED         // Aggregated from multiple sources
    }

    /**
     * @notice Price feed configuration
     */
    struct PriceFeed {
        address feedAddress;        // Price feed contract address
        PriceSource source;         // Price source type
        uint8 decimals;             // Price decimals
        uint256 heartbeat;          // Maximum time between updates
        uint256 lastPrice;          // Last reported price
        uint256 lastUpdateTime;     // Last update timestamp
        bool isActive;              // Feed status
        uint256 minPrice;           // Minimum acceptable price (circuit breaker)
        uint256 maxPrice;           // Maximum acceptable price (circuit breaker)
    }

    /**
     * @notice Aggregated price data
     */
    struct AggregatedPrice {
        uint256 price;              // Final aggregated price
        uint256 confidence;         // Confidence level (0-10000 basis points)
        uint256 timestamp;          // Price timestamp
        uint256 deviation;          // Max deviation from median
        uint256 sourcesUsed;        // Number of sources used
    }

    // ============ Events ============

    event PriceUpdated(
        address indexed token,
        uint256 price,
        uint256 timestamp,
        PriceSource source
    );

    event PriceFeedAdded(
        address indexed token,
        address indexed feedAddress,
        PriceSource source
    );

    event PriceFeedRemoved(
        address indexed token,
        address indexed feedAddress
    );

    event PriceFeedStatusChanged(
        address indexed token,
        address indexed feedAddress,
        bool isActive
    );

    event StalePriceDetected(
        address indexed token,
        address indexed feedAddress,
        uint256 lastUpdate,
        uint256 heartbeat
    );

    event CircuitBreakerTriggered(
        address indexed token,
        uint256 attemptedPrice,
        uint256 minPrice,
        uint256 maxPrice
    );

    event OracleConfigUpdated(string parameter, uint256 newValue);

    // ============ Price Queries ============

    /**
     * @notice Get latest price for a token
     * @param token Token address
     * @return price Current price
     * @return decimals Price decimals
     * @return timestamp Price timestamp
     * @dev Reverts if price is stale or unavailable
     */
    function getLatestPrice(address token)
        external
        view
        returns (uint256 price, uint8 decimals, uint256 timestamp);

    /**
     * @notice Get price with confidence level
     * @param token Token address
     * @return price Current price
     * @return confidence Confidence level (0-10000)
     * @return timestamp Price timestamp
     */
    function getPriceWithConfidence(address token)
        external
        view
        returns (uint256 price, uint256 confidence, uint256 timestamp);

    /**
     * @notice Get aggregated price from multiple sources
     * @param token Token address
     * @return Aggregated price data
     */
    function getAggregatedPrice(address token) external view returns (AggregatedPrice memory);

    /**
     * @notice Get price of tokenA in terms of tokenB
     * @param tokenA Base token
     * @param tokenB Quote token
     * @return Exchange rate
     */
    function getPriceInTermsOf(address tokenA, address tokenB) external view returns (uint256);

    /**
     * @notice Try to get price (returns false instead of reverting)
     * @param token Token address
     * @return success True if price is available
     * @return price Current price (0 if unavailable)
     * @return timestamp Price timestamp
     */
    function tryGetPrice(address token)
        external
        view
        returns (bool success, uint256 price, uint256 timestamp);

    // ============ Price Feed Management ============

    /**
     * @notice Add a price feed for a token
     * @param token Token address
     * @param feedAddress Price feed contract address
     * @param source Price source type
     * @param heartbeat Maximum time between updates
     * @dev Only callable by authorized roles
     */
    function addPriceFeed(
        address token,
        address feedAddress,
        PriceSource source,
        uint256 heartbeat
    ) external;

    /**
     * @notice Remove a price feed
     * @param token Token address
     * @param feedAddress Feed address to remove
     */
    function removePriceFeed(address token, address feedAddress) external;

    /**
     * @notice Activate/deactivate a price feed
     * @param token Token address
     * @param feedAddress Feed address
     * @param isActive New status
     */
    function setPriceFeedStatus(address token, address feedAddress, bool isActive) external;

    /**
     * @notice Get all price feeds for a token
     * @param token Token address
     * @return Array of price feed configurations
     */
    function getPriceFeeds(address token) external view returns (PriceFeed[] memory);

    /**
     * @notice Check if token has an active price feed
     * @param token Token address
     * @return True if active feed exists
     */
    function hasPriceFeed(address token) external view returns (bool);

    // ============ Price Updates ============

    /**
     * @notice Update price manually (for custom/manual sources)
     * @param token Token address
     * @param price New price
     * @dev Only callable by authorized price updaters
     */
    function updatePrice(address token, uint256 price) external;

    /**
     * @notice Batch update prices
     * @param tokens Array of token addresses
     * @param prices Array of prices
     */
    function batchUpdatePrices(address[] calldata tokens, uint256[] calldata prices) external;

    /**
     * @notice Force refresh price from oracle
     * @param token Token address
     * @dev Useful for on-demand updates
     */
    function refreshPrice(address token) external;

    // ============ Staleness & Validation ============

    /**
     * @notice Check if price is stale
     * @param token Token address
     * @return True if price is stale
     */
    function isPriceStale(address token) external view returns (bool);

    /**
     * @notice Get time since last price update
     * @param token Token address
     * @return Seconds since last update
     */
    function getTimeSinceUpdate(address token) external view returns (uint256);

    /**
     * @notice Validate price against circuit breaker limits
     * @param token Token address
     * @param price Price to validate
     * @return valid True if price is within limits
     */
    function validatePrice(address token, uint256 price) external view returns (bool valid);

    /**
     * @notice Set circuit breaker limits
     * @param token Token address
     * @param minPrice Minimum acceptable price
     * @param maxPrice Maximum acceptable price
     */
    function setCircuitBreaker(address token, uint256 minPrice, uint256 maxPrice) external;

    // ============ Historical Data ============

    /**
     * @notice Get historical price
     * @param token Token address
     * @param timestamp Timestamp to query
     * @return price Historical price
     * @return actualTimestamp Actual timestamp of the price
     * @dev Returns closest available price if exact timestamp not found
     */
    function getHistoricalPrice(address token, uint256 timestamp)
        external
        view
        returns (uint256 price, uint256 actualTimestamp);

    /**
     * @notice Get price change over period
     * @param token Token address
     * @param period Time period in seconds
     * @return priceChange Change in price (can be negative)
     * @return percentChange Percentage change in basis points
     */
    function getPriceChange(address token, uint256 period)
        external
        view
        returns (int256 priceChange, int256 percentChange);

    /**
     * @notice Get TWAP (Time-Weighted Average Price)
     * @param token Token address
     * @param period Time period for TWAP calculation
     * @return TWAP over the specified period
     */
    function getTWAP(address token, uint256 period) external view returns (uint256);

    // ============ Configuration ============

    /**
     * @notice Set global staleness threshold
     * @param threshold Maximum acceptable staleness in seconds
     */
    function setStalenessThreshold(uint256 threshold) external;

    /**
     * @notice Set minimum confidence level
     * @param minConfidence Minimum confidence (0-10000 basis points)
     */
    function setMinConfidence(uint256 minConfidence) external;

    /**
     * @notice Set price deviation tolerance for aggregation
     * @param maxDeviation Maximum deviation in basis points
     */
    function setMaxDeviation(uint256 maxDeviation) external;

    /**
     * @notice Get oracle configuration
     * @return stalenessThreshold Current staleness threshold
     * @return minConfidence Minimum confidence level
     * @return maxDeviation Maximum price deviation
     */
    function getOracleConfig()
        external
        view
        returns (
            uint256 stalenessThreshold,
            uint256 minConfidence,
            uint256 maxDeviation
        );

    // ============ Statistics ============

    /**
     * @notice Get oracle statistics
     * @return totalTokens Number of tokens with price feeds
     * @return activeFeeds Number of active price feeds
     * @return stalePrices Number of stale prices
     * @return lastGlobalUpdate Last update across all feeds
     */
    function getOracleStats()
        external
        view
        returns (
            uint256 totalTokens,
            uint256 activeFeeds,
            uint256 stalePrices,
            uint256 lastGlobalUpdate
        );

    /**
     * @notice Get supported tokens
     * @return Array of token addresses with price feeds
     */
    function getSupportedTokens() external view returns (address[] memory);
}
