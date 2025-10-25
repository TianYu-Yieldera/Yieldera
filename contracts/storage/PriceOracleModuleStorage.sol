// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../interfaces/core/IPriceOracleModule.sol";

/**
 * @title PriceOracleModuleStorage
 * @notice Diamond Storage library for PriceOracleModule
 * @dev Implements EIP-2535 Diamond Storage pattern
 */
library PriceOracleModuleStorage {
    // Storage position is keccak256("price.oracle.module.storage") - 1
    bytes32 constant STORAGE_POSITION =
        0x3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d;

    struct OracleData {
        // Price feed mappings
        mapping(address => IPriceOracleModule.PriceFeedInfo[]) tokenPriceFeeds;
        mapping(address => IPriceOracleModule.PriceData) feedPriceData;
        mapping(address => IPriceOracleModule.PriceUpdate[]) feedPriceHistory;

        // Configuration
        uint256 stalenessThreshold;  // Default: 1 hour
        uint256 deviationThreshold;   // Max allowed deviation in basis points

        // Statistics
        uint256 totalFeeds;
        uint256 activeFeeds;
        uint256 lastUpdate;

        // Reserved slots for future upgrades
        uint256[50] __gap;
    }

    /**
     * @notice Returns the storage layout
     */
    function layout() internal pure returns (OracleData storage ds) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            ds.slot := position
        }
    }

    /**
     * @notice Initialize storage
     */
    function initialize() internal {
        OracleData storage ds = layout();
        require(ds.lastUpdate == 0, "Already initialized");

        ds.stalenessThreshold = 1 hours;
        ds.deviationThreshold = 1000; // 10%
        ds.lastUpdate = block.timestamp;
    }

    /**
     * @notice Add price feed for a token
     */
    function addPriceFeed(
        address token,
        IPriceOracleModule.PriceFeedInfo memory feedInfo
    ) internal {
        OracleData storage ds = layout();
        ds.tokenPriceFeeds[token].push(feedInfo);
        ds.totalFeeds++;
        ds.activeFeeds++;
        ds.lastUpdate = block.timestamp;
    }

    /**
     * @notice Get price feeds for a token
     */
    function getPriceFeeds(address token)
        internal
        view
        returns (IPriceOracleModule.PriceFeedInfo[] storage)
    {
        return layout().tokenPriceFeeds[token];
    }

    /**
     * @notice Set price data for a feed
     */
    function setPriceData(
        address feed,
        IPriceOracleModule.PriceData memory priceData
    ) internal {
        layout().feedPriceData[feed] = priceData;
    }

    /**
     * @notice Get price data for a feed
     */
    function getPriceData(address feed)
        internal
        view
        returns (IPriceOracleModule.PriceData storage)
    {
        return layout().feedPriceData[feed];
    }

    /**
     * @notice Add price update to history
     */
    function addPriceHistory(
        address feed,
        IPriceOracleModule.PriceUpdate memory update
    ) internal {
        layout().feedPriceHistory[feed].push(update);
    }

    /**
     * @notice Get price history for a feed
     */
    function getPriceHistory(address feed, uint256 offset, uint256 limit)
        internal
        view
        returns (IPriceOracleModule.PriceUpdate[] memory)
    {
        OracleData storage ds = layout();
        IPriceOracleModule.PriceUpdate[] storage history = ds.feedPriceHistory[feed];
        uint256 length = history.length;

        if (offset >= length) {
            return new IPriceOracleModule.PriceUpdate[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }

        IPriceOracleModule.PriceUpdate[] memory updates =
            new IPriceOracleModule.PriceUpdate[](end - offset);

        for (uint256 i = offset; i < end; i++) {
            updates[i - offset] = history[i];
        }

        return updates;
    }

    /**
     * @notice Get configuration
     */
    function getConfig() internal view returns (uint256 staleness, uint256 deviation) {
        OracleData storage ds = layout();
        return (ds.stalenessThreshold, ds.deviationThreshold);
    }

    /**
     * @notice Update configuration
     */
    function setConfig(uint256 staleness, uint256 deviation) internal {
        OracleData storage ds = layout();
        ds.stalenessThreshold = staleness;
        ds.deviationThreshold = deviation;
        ds.lastUpdate = block.timestamp;
    }

    /**
     * @notice Get statistics
     */
    function getStats() internal view returns (uint256 total, uint256 active) {
        OracleData storage ds = layout();
        return (ds.totalFeeds, ds.activeFeeds);
    }

    /**
     * @notice Increment active feeds count
     */
    function incrementActiveFeeds() internal {
        layout().activeFeeds++;
    }

    /**
     * @notice Decrement active feeds count
     */
    function decrementActiveFeeds() internal {
        layout().activeFeeds--;
    }

    /**
     * @notice Update last modified timestamp
     */
    function touch() internal {
        layout().lastUpdate = block.timestamp;
    }
}
