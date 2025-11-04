// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../interfaces/IPriceOracle.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title ChainlinkPriceOracle
 * @notice Chainlink-based price oracle implementation
 * @dev Aggregates prices from Chainlink data feeds
 */
contract ChainlinkPriceOracle is IPriceOracle, Ownable {

    // Mapping from asset address to Chainlink price feed
    mapping(address => AggregatorV3Interface) public priceFeeds;

    // Price staleness threshold (1 hour)
    uint256 public constant STALENESS_THRESHOLD = 3600;

    event PriceFeedUpdated(address indexed asset, address indexed feed);

    constructor(address initialOwner) Ownable(initialOwner) {}

    /**
     * @notice Set Chainlink price feed for an asset
     * @param asset Address of the asset
     * @param feed Address of Chainlink price feed
     */
    function setPriceFeed(address asset, address feed) external onlyOwner {
        require(asset != address(0), "Invalid asset");
        require(feed != address(0), "Invalid feed");

        priceFeeds[asset] = AggregatorV3Interface(feed);
        emit PriceFeedUpdated(asset, feed);
    }

    /**
     * @notice Get the latest price of an asset
     * @param asset Address of the asset
     * @return price Price in USD with 8 decimals
     */
    function getAssetPrice(address asset) external view override returns (uint256) {
        AggregatorV3Interface feed = priceFeeds[asset];
        require(address(feed) != address(0), "Price feed not set");

        (
            uint80 roundId,
            int256 answer,
            ,
            uint256 updatedAt,
            uint80 answeredInRound
        ) = feed.latestRoundData();

        require(answer > 0, "Invalid price");
        require(updatedAt > 0, "Invalid timestamp");
        require(answeredInRound >= roundId, "Stale price");
        require(block.timestamp - updatedAt <= STALENESS_THRESHOLD, "Price too old");

        return uint256(answer);
    }

    /**
     * @notice Get decimals for price data
     * @return decimals Always returns 8 (Chainlink standard)
     */
    function decimals() external pure override returns (uint8) {
        return 8;
    }

    /**
     * @notice Check if price feed exists for an asset
     * @param asset Address of the asset
     * @return exists True if feed is configured
     */
    function hasPriceFeed(address asset) external view returns (bool) {
        return address(priceFeeds[asset]) != address(0);
    }
}
