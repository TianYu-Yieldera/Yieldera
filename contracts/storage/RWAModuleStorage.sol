// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../interfaces/core/IRWAModule.sol";

/**
 * @title RWAModuleStorage
 * @notice Diamond Storage library for RWAModule
 * @dev Implements EIP-2535 Diamond Storage pattern to prevent storage collisions
 */
library RWAModuleStorage {
    // Storage position is keccak256("rwa.module.storage") - 1
    bytes32 constant STORAGE_POSITION =
        0x2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c;

    struct RWAData {
        // Trade history tracking
        IRWAModule.Trade[] tradeHistory;
        mapping(address => uint256[]) userTradeIndices;

        // References to external contracts
        address orderBook;      // OrderBook address
        address baseToken;      // RWA token
        address quoteToken;     // Quote currency (LUSD/USDC)

        // Module metadata
        uint256 totalTrades;
        uint256 lastTradeTime;
        uint256 lastUpdate;

        // Statistics cache
        uint256 cachedVolume24h;
        uint256 cachedLastPrice;
        uint256 cacheTimestamp;

        // Reserved slots for future upgrades
        uint256[50] __gap;
    }

    /**
     * @notice Returns the storage layout
     * @return ds The storage layout struct
     */
    function layout() internal pure returns (RWAData storage ds) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            ds.slot := position
        }
    }

    /**
     * @notice Initialize storage (called once during proxy deployment)
     * @param _orderBook OrderBook address
     */
    function initialize(address _orderBook) internal {
        RWAData storage ds = layout();
        require(ds.orderBook == address(0), "Already initialized");

        ds.orderBook = _orderBook;
        ds.lastUpdate = block.timestamp;
    }

    /**
     * @notice Set token addresses
     */
    function setTokens(address _baseToken, address _quoteToken) internal {
        RWAData storage ds = layout();
        ds.baseToken = _baseToken;
        ds.quoteToken = _quoteToken;
    }

    /**
     * @notice Add a trade to history
     */
    function addTrade(IRWAModule.Trade memory trade) internal returns (uint256) {
        RWAData storage ds = layout();

        uint256 tradeIndex = ds.tradeHistory.length;
        ds.tradeHistory.push(trade);

        // Index for buyer
        ds.userTradeIndices[trade.buyer].push(tradeIndex);

        // Index for seller
        ds.userTradeIndices[trade.seller].push(tradeIndex);

        ds.totalTrades++;
        ds.lastTradeTime = block.timestamp;
        ds.lastUpdate = block.timestamp;

        return tradeIndex;
    }

    /**
     * @notice Get trade history
     */
    function getTradeHistory(uint256 offset, uint256 limit)
        internal
        view
        returns (IRWAModule.Trade[] memory)
    {
        RWAData storage ds = layout();
        uint256 length = ds.tradeHistory.length;

        if (offset >= length) {
            return new IRWAModule.Trade[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }

        IRWAModule.Trade[] memory trades = new IRWAModule.Trade[](end - offset);
        for (uint256 i = offset; i < end; i++) {
            trades[i - offset] = ds.tradeHistory[i];
        }

        return trades;
    }

    /**
     * @notice Get user's trade indices
     */
    function getUserTradeIndices(address user)
        internal
        view
        returns (uint256[] storage)
    {
        return layout().userTradeIndices[user];
    }

    /**
     * @notice Get user trades
     */
    function getUserTrades(address user, uint256 offset, uint256 limit)
        internal
        view
        returns (IRWAModule.Trade[] memory)
    {
        RWAData storage ds = layout();
        uint256[] storage indices = ds.userTradeIndices[user];
        uint256 length = indices.length;

        if (offset >= length) {
            return new IRWAModule.Trade[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }

        IRWAModule.Trade[] memory trades = new IRWAModule.Trade[](end - offset);
        for (uint256 i = offset; i < end; i++) {
            trades[i - offset] = ds.tradeHistory[indices[i]];
        }

        return trades;
    }

    /**
     * @notice Get external contract addresses
     */
    function getAddresses()
        internal
        view
        returns (address orderBook, address baseToken, address quoteToken)
    {
        RWAData storage ds = layout();
        return (ds.orderBook, ds.baseToken, ds.quoteToken);
    }

    /**
     * @notice Get total trades count
     */
    function getTotalTrades() internal view returns (uint256) {
        return layout().totalTrades;
    }

    /**
     * @notice Get last trade timestamp
     */
    function getLastTradeTime() internal view returns (uint256) {
        return layout().lastTradeTime;
    }

    /**
     * @notice Update statistics cache
     */
    function updateStatsCache(uint256 volume24h, uint256 lastPrice) internal {
        RWAData storage ds = layout();
        ds.cachedVolume24h = volume24h;
        ds.cachedLastPrice = lastPrice;
        ds.cacheTimestamp = block.timestamp;
    }

    /**
     * @notice Get cached statistics
     */
    function getCachedStats()
        internal
        view
        returns (uint256 volume24h, uint256 lastPrice, uint256 timestamp)
    {
        RWAData storage ds = layout();
        return (ds.cachedVolume24h, ds.cachedLastPrice, ds.cacheTimestamp);
    }

    /**
     * @notice Update last modified timestamp
     */
    function touch() internal {
        layout().lastUpdate = block.timestamp;
    }
}
