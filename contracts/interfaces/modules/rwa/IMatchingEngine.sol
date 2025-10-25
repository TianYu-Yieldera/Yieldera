// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IOrderManager.sol";

/**
 * @title IMatchingEngine
 * @notice Interface for trade matching engine
 * @dev Handles order matching and trade execution
 */
interface IMatchingEngine {
    // ============ Structs ============

    struct Trade {
        uint256 tradeId;
        uint256 buyOrderId;
        uint256 sellOrderId;
        address buyer;
        address seller;
        uint256 amount;
        uint256 price;
        uint256 timestamp;
    }

    struct MatchResult {
        uint256 tradeId;
        uint256 matchedAmount;
        uint256 executionPrice;
        bool fullyMatched;
    }

    // ============ Events ============

    event TradeExecuted(
        uint256 indexed tradeId,
        uint256 indexed buyOrderId,
        uint256 indexed sellOrderId,
        address buyer,
        address seller,
        uint256 amount,
        uint256 price
    );
    event OrderMatched(uint256 indexed orderId, uint256 matchedAmount);

    // ============ Functions ============

    /**
     * @notice Match a buy order with sell orders
     * @param buyOrderId Buy order ID
     * @return results Array of match results
     */
    function matchBuyOrder(uint256 buyOrderId) external returns (MatchResult[] memory results);

    /**
     * @notice Match a sell order with buy orders
     * @param sellOrderId Sell order ID
     * @return results Array of match results
     */
    function matchSellOrder(uint256 sellOrderId) external returns (MatchResult[] memory results);

    /**
     * @notice Execute a trade between two orders
     * @param buyOrderId Buy order ID
     * @param sellOrderId Sell order ID
     * @param amount Amount to trade
     * @param price Trade price
     * @return tradeId Executed trade ID
     */
    function executeTrade(uint256 buyOrderId, uint256 sellOrderId, uint256 amount, uint256 price)
        external
        returns (uint256 tradeId);

    /**
     * @notice Get trade details
     * @param tradeId Trade ID
     * @return Trade data
     */
    function getTrade(uint256 tradeId) external view returns (Trade memory);

    /**
     * @notice Get user's trade history
     * @param user User address
     * @return Array of trade IDs
     */
    function getUserTrades(address user) external view returns (uint256[] memory);

    /**
     * @notice Get recent trades
     * @param count Number of trades to return
     * @return Array of trade IDs
     */
    function getRecentTrades(uint256 count) external view returns (uint256[] memory);

    /**
     * @notice Find matching orders for a given order
     * @param orderId Order ID
     * @param orderType Order type
     * @return matchingOrderIds Array of matching order IDs
     */
    function findMatchingOrders(uint256 orderId, IOrderManager.OrderType orderType)
        external
        view
        returns (uint256[] memory matchingOrderIds);

    /**
     * @notice Check if two orders can be matched
     * @param buyOrderId Buy order ID
     * @param sellOrderId Sell order ID
     * @return canMatch True if orders can be matched
     * @return matchAmount Amount that can be matched
     */
    function canMatch(uint256 buyOrderId, uint256 sellOrderId)
        external
        view
        returns (bool canMatch, uint256 matchAmount);

    /**
     * @notice Get best available price for order type
     * @param orderType Order type (BUY/SELL)
     * @return Best price available
     */
    function getBestPrice(IOrderManager.OrderType orderType) external view returns (uint256);
}
