// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "./interfaces/IPriceOracle.sol";

/**
 * @title PriceOracle
 * @notice Price oracle implementation supporting multiple price feeds
 * @dev Can integrate with Chainlink, Uniswap TWAP, or custom price feeds
 */
contract PriceOracle is IPriceOracle, Ownable, Pausable {
    // Price feed data structure
    struct PriceFeedData {
        address feedAddress;
        uint256 lastPrice;
        uint256 lastUpdateTime;
        uint8 decimals;
        bool isActive;
    }

    // Mapping from token address to price feed data
    mapping(address => PriceFeedData) public priceFeeds;

    // Stale price threshold (default: 1 hour)
    uint256 public stalePriceThreshold = 3600;

    // Manual price override for testing/emergency
    mapping(address => uint256) public manualPriceOverride;
    mapping(address => bool) public useManualPrice;

    // Constants
    uint8 public constant USD_DECIMALS = 8;
    uint256 public constant PRICE_PRECISION = 1e8;

    // Events (additional to interface events)
    event StalePriceThresholdUpdated(uint256 newThreshold);
    event ManualPriceSet(address indexed token, uint256 price);
    event ManualPriceRemoved(address indexed token);

    constructor() {}

    /**
     * @notice Get the latest price for a token
     * @param token The token address
     * @return price The price in USD with 8 decimals
     * @return decimals The number of decimals
     * @return timestamp The timestamp of the price
     */
    function getLatestPrice(address token)
        external
        view
        override
        whenNotPaused
        returns (uint256 price, uint8 decimals, uint256 timestamp)
    {
        require(token != address(0), "Invalid token address");

        // Check for manual price override
        if (useManualPrice[token]) {
            return (manualPriceOverride[token], USD_DECIMALS, block.timestamp);
        }

        PriceFeedData memory feed = priceFeeds[token];
        require(feed.isActive, "Price feed not available");
        require(!_isPriceStale(feed.lastUpdateTime), "Price data is stale");

        // In production, this would fetch from actual oracle
        // For now, return stored price
        return (feed.lastPrice, feed.decimals, feed.lastUpdateTime);
    }

    /**
     * @notice Get price of tokenA in terms of tokenB
     * @param tokenA Base token
     * @param tokenB Quote token
     * @return price Exchange rate with 8 decimals
     */
    function getPriceInTermsOf(address tokenA, address tokenB)
        external
        view
        override
        whenNotPaused
        returns (uint256)
    {
        require(tokenA != address(0) && tokenB != address(0), "Invalid token address");

        (uint256 priceA, , ) = this.getLatestPrice(tokenA);
        (uint256 priceB, , ) = this.getLatestPrice(tokenB);

        require(priceB > 0, "Invalid quote token price");

        // Calculate cross rate: priceA / priceB
        return (priceA * PRICE_PRECISION) / priceB;
    }

    /**
     * @notice Check if price feed exists for token
     * @param token Token to check
     * @return True if feed exists and is active
     */
    function hasPriceFeed(address token) external view override returns (bool) {
        return priceFeeds[token].isActive || useManualPrice[token];
    }

    /**
     * @notice Set price feed for a token (admin only)
     * @param token Token address
     * @param priceFeed Price feed address
     */
    function setPriceFeed(address token, address priceFeed)
        external
        override
        onlyOwner
    {
        require(token != address(0), "Invalid token address");
        require(priceFeed != address(0), "Invalid feed address");

        priceFeeds[token] = PriceFeedData({
            feedAddress: priceFeed,
            lastPrice: PRICE_PRECISION, // Default to $1
            lastUpdateTime: block.timestamp,
            decimals: USD_DECIMALS,
            isActive: true
        });

        emit PriceFeedUpdated(token, priceFeed);
    }

    /**
     * @notice Get price feed address for token
     * @param token Token address
     * @return Price feed address
     */
    function getPriceFeed(address token) external view override returns (address) {
        return priceFeeds[token].feedAddress;
    }

    /**
     * @notice Check if price is stale
     * @param token Token to check
     * @return True if price is stale
     */
    function isPriceStale(address token) external view override returns (bool) {
        if (useManualPrice[token]) {
            return false; // Manual prices are never stale
        }

        PriceFeedData memory feed = priceFeeds[token];
        if (!feed.isActive) {
            return true;
        }

        return _isPriceStale(feed.lastUpdateTime);
    }

    /**
     * @notice Update price for a token (oracle keeper function)
     * @param token Token address
     * @param price New price
     * @dev In production, this would be called by oracle nodes or keepers
     */
    function updatePrice(address token, uint256 price) external onlyOwner {
        require(priceFeeds[token].isActive, "Price feed not active");
        require(price > 0, "Invalid price");

        priceFeeds[token].lastPrice = price;
        priceFeeds[token].lastUpdateTime = block.timestamp;

        emit PriceUpdated(token, price, block.timestamp);
    }

    /**
     * @notice Set manual price override (emergency/testing)
     * @param token Token address
     * @param price Manual price to set
     */
    function setManualPrice(address token, uint256 price) external onlyOwner {
        require(token != address(0), "Invalid token address");
        require(price > 0, "Invalid price");

        manualPriceOverride[token] = price;
        useManualPrice[token] = true;

        emit ManualPriceSet(token, price);
    }

    /**
     * @notice Remove manual price override
     * @param token Token address
     */
    function removeManualPrice(address token) external onlyOwner {
        useManualPrice[token] = false;
        delete manualPriceOverride[token];

        emit ManualPriceRemoved(token);
    }

    /**
     * @notice Update stale price threshold
     * @param newThreshold New threshold in seconds
     */
    function setStalePriceThreshold(uint256 newThreshold) external onlyOwner {
        require(newThreshold > 0, "Invalid threshold");
        stalePriceThreshold = newThreshold;

        emit StalePriceThresholdUpdated(newThreshold);
    }

    /**
     * @notice Pause oracle
     */
    function pause() external override onlyOwner {
        _pause();
        emit OraclePaused(msg.sender);
    }

    /**
     * @notice Unpause oracle
     */
    function unpause() external override onlyOwner {
        _unpause();
        emit OracleUnpaused(msg.sender);
    }

    /**
     * @notice Internal function to check if price is stale
     * @param lastUpdate Last update timestamp
     * @return True if stale
     */
    function _isPriceStale(uint256 lastUpdate) internal view returns (bool) {
        return block.timestamp > lastUpdate + stalePriceThreshold;
    }

    /**
     * @notice Batch update prices (gas optimization)
     * @param tokens Array of token addresses
     * @param prices Array of prices
     */
    function batchUpdatePrices(
        address[] calldata tokens,
        uint256[] calldata prices
    ) external onlyOwner {
        require(tokens.length == prices.length, "Array length mismatch");

        for (uint256 i = 0; i < tokens.length; i++) {
            if (priceFeeds[tokens[i]].isActive && prices[i] > 0) {
                priceFeeds[tokens[i]].lastPrice = prices[i];
                priceFeeds[tokens[i]].lastUpdateTime = block.timestamp;

                emit PriceUpdated(tokens[i], prices[i], block.timestamp);
            }
        }
    }

    /**
     * @notice Deactivate a price feed
     * @param token Token address
     */
    function deactivatePriceFeed(address token) external onlyOwner {
        priceFeeds[token].isActive = false;
    }

    /**
     * @notice Reactivate a price feed
     * @param token Token address
     */
    function reactivatePriceFeed(address token) external onlyOwner {
        require(priceFeeds[token].feedAddress != address(0), "Feed not configured");
        priceFeeds[token].isActive = true;
    }
}