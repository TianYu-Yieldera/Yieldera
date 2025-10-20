// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IPriceOracle
 * @notice Interface for price oracle implementations
 * @dev Supports multiple price feed sources and token pairs
 */
interface IPriceOracle {
    /**
     * @notice Get the latest price for a token
     * @param token The token address
     * @return price The price in USD with 8 decimals
     * @return decimals The number of decimals in the price
     * @return timestamp The timestamp of the price update
     */
    function getLatestPrice(address token)
        external
        view
        returns (uint256 price, uint8 decimals, uint256 timestamp);

    /**
     * @notice Get the price of tokenA in terms of tokenB
     * @param tokenA The base token
     * @param tokenB The quote token
     * @return price The exchange rate with 8 decimals
     */
    function getPriceInTermsOf(address tokenA, address tokenB)
        external
        view
        returns (uint256 price);

    /**
     * @notice Check if a price feed is available for a token
     * @param token The token to check
     * @return True if price feed exists
     */
    function hasPriceFeed(address token) external view returns (bool);

    /**
     * @notice Update price feed address for a token (admin only)
     * @param token The token address
     * @param priceFeed The price feed address
     */
    function setPriceFeed(address token, address priceFeed) external;

    /**
     * @notice Get the address of the price feed for a token
     * @param token The token address
     * @return The price feed address
     */
    function getPriceFeed(address token) external view returns (address);

    /**
     * @notice Check if the price data is stale
     * @param token The token to check
     * @return True if price is stale (older than acceptable threshold)
     */
    function isPriceStale(address token) external view returns (bool);

    /**
     * @notice Emergency pause oracle updates
     */
    function pause() external;

    /**
     * @notice Resume oracle updates
     */
    function unpause() external;

    // Events
    event PriceFeedUpdated(address indexed token, address indexed priceFeed);
    event PriceUpdated(address indexed token, uint256 price, uint256 timestamp);
    event OraclePaused(address indexed by);
    event OracleUnpaused(address indexed by);
}