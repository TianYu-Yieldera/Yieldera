// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/modules/rwa/IMatchingEngine.sol";
import "../../interfaces/modules/rwa/IOrderManager.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title MatchingEngine
 * @notice Handles trade matching and execution
 * @dev Matches buy and sell orders based on price-time priority
 */
contract MatchingEngine is IMatchingEngine, Ownable {
    // ============ State Variables ============

    address public rwaModule; // Main coordinator contract
    IOrderManager public orderManager;

    // Storage for trades (using Diamond Storage pattern)
    bytes32 private constant MATCHING_STORAGE_POSITION = keccak256("matching.engine.storage");

    struct MatchingStorage {
        mapping(uint256 => Trade) trades;
        mapping(address => uint256[]) userTrades;
        uint256[] recentTrades;
        uint256 nextTradeId;
        uint256 maxRecentTrades;
    }

    // ============ Modifiers ============

    modifier onlyRWAModule() {
        require(msg.sender == rwaModule, "Only RWA module");
        _;
    }

    // ============ Constructor ============

    constructor(address _orderManager) {
        require(_orderManager != address(0), "Invalid order manager");
        orderManager = IOrderManager(_orderManager);

        MatchingStorage storage ms = _getStorage();
        ms.maxRecentTrades = 100; // Keep last 100 trades
    }

    // ============ Admin Functions ============

    function setRWAModule(address _rwaModule) external onlyOwner {
        require(_rwaModule != address(0), "Invalid address");
        rwaModule = _rwaModule;
    }

    function setOrderManager(address _orderManager) external onlyOwner {
        require(_orderManager != address(0), "Invalid address");
        orderManager = IOrderManager(_orderManager);
    }

    // ============ Internal Storage ============

    function _getStorage() private pure returns (MatchingStorage storage ms) {
        bytes32 position = MATCHING_STORAGE_POSITION;
        assembly {
            ms.slot := position
        }
    }

    // ============ IMatchingEngine Implementation ============

    function matchBuyOrder(uint256 buyOrderId)
        external
        override
        onlyRWAModule
        returns (MatchResult[] memory results)
    {
        IOrderManager.Order memory buyOrder = orderManager.getOrder(buyOrderId);
        require(buyOrder.orderType == IOrderManager.OrderType.BUY, "Not a buy order");
        require(buyOrder.status == IOrderManager.OrderStatus.ACTIVE, "Order not active");

        uint256[] memory sellOrders = orderManager.getOrdersByType(IOrderManager.OrderType.SELL);
        return _matchOrders(buyOrderId, sellOrders, true);
    }

    function matchSellOrder(uint256 sellOrderId)
        external
        override
        onlyRWAModule
        returns (MatchResult[] memory results)
    {
        IOrderManager.Order memory sellOrder = orderManager.getOrder(sellOrderId);
        require(sellOrder.orderType == IOrderManager.OrderType.SELL, "Not a sell order");
        require(sellOrder.status == IOrderManager.OrderStatus.ACTIVE, "Order not active");

        uint256[] memory buyOrders = orderManager.getOrdersByType(IOrderManager.OrderType.BUY);
        return _matchOrders(sellOrderId, buyOrders, false);
    }

    function executeTrade(
        uint256 buyOrderId,
        uint256 sellOrderId,
        uint256 amount,
        uint256 price
    ) external override onlyRWAModule returns (uint256 tradeId) {
        IOrderManager.Order memory buyOrder = orderManager.getOrder(buyOrderId);
        IOrderManager.Order memory sellOrder = orderManager.getOrder(sellOrderId);

        require(buyOrder.status == IOrderManager.OrderStatus.ACTIVE, "Buy order not active");
        require(sellOrder.status == IOrderManager.OrderStatus.ACTIVE, "Sell order not active");
        require(amount > 0, "Invalid amount");
        require(buyOrder.price >= price && sellOrder.price <= price, "Price mismatch");

        MatchingStorage storage ms = _getStorage();
        tradeId = ms.nextTradeId++;

        Trade storage trade = ms.trades[tradeId];
        trade.tradeId = tradeId;
        trade.buyOrderId = buyOrderId;
        trade.sellOrderId = sellOrderId;
        trade.buyer = buyOrder.trader;
        trade.seller = sellOrder.trader;
        trade.amount = amount;
        trade.price = price;
        trade.timestamp = block.timestamp;

        // Update user trades
        ms.userTrades[buyOrder.trader].push(tradeId);
        ms.userTrades[sellOrder.trader].push(tradeId);

        // Add to recent trades
        ms.recentTrades.push(tradeId);
        if (ms.recentTrades.length > ms.maxRecentTrades) {
            // Remove oldest trade (simple implementation)
            for (uint256 i = 0; i < ms.recentTrades.length - 1; i++) {
                ms.recentTrades[i] = ms.recentTrades[i + 1];
            }
            ms.recentTrades.pop();
        }

        // Fill orders through OrderManager (will be called by RWA module)
        emit TradeExecuted(tradeId, buyOrderId, sellOrderId, buyOrder.trader, sellOrder.trader, amount, price);
        emit OrderMatched(buyOrderId, amount);
        emit OrderMatched(sellOrderId, amount);

        return tradeId;
    }

    function getTrade(uint256 tradeId) external view override returns (Trade memory) {
        MatchingStorage storage ms = _getStorage();
        return ms.trades[tradeId];
    }

    function getUserTrades(address user) external view override returns (uint256[] memory) {
        MatchingStorage storage ms = _getStorage();
        return ms.userTrades[user];
    }

    function getRecentTrades(uint256 count) external view override returns (uint256[] memory) {
        MatchingStorage storage ms = _getStorage();

        uint256 length = ms.recentTrades.length;
        if (count > length) count = length;

        uint256[] memory trades = new uint256[](count);
        for (uint256 i = 0; i < count; i++) {
            trades[i] = ms.recentTrades[length - count + i];
        }

        return trades;
    }

    function findMatchingOrders(uint256 orderId, IOrderManager.OrderType orderType)
        external
        view
        override
        returns (uint256[] memory matchingOrderIds)
    {
        IOrderManager.Order memory order = orderManager.getOrder(orderId);
        require(order.status == IOrderManager.OrderStatus.ACTIVE, "Order not active");

        // Get opposite order type
        IOrderManager.OrderType oppositeType = orderType == IOrderManager.OrderType.BUY
            ? IOrderManager.OrderType.SELL
            : IOrderManager.OrderType.BUY;

        uint256[] memory oppositeOrders = orderManager.getOrdersByType(oppositeType);

        // Count matching orders
        uint256 matchCount = 0;
        for (uint256 i = 0; i < oppositeOrders.length; i++) {
            IOrderManager.Order memory oppositeOrder = orderManager.getOrder(oppositeOrders[i]);
            if (_canMatchOrders(order, oppositeOrder, orderType == IOrderManager.OrderType.BUY)) {
                matchCount++;
            }
        }

        // Populate result
        matchingOrderIds = new uint256[](matchCount);
        uint256 index = 0;
        for (uint256 i = 0; i < oppositeOrders.length; i++) {
            IOrderManager.Order memory oppositeOrder = orderManager.getOrder(oppositeOrders[i]);
            if (_canMatchOrders(order, oppositeOrder, orderType == IOrderManager.OrderType.BUY)) {
                matchingOrderIds[index++] = oppositeOrders[i];
            }
        }

        return matchingOrderIds;
    }

    function canMatch(uint256 buyOrderId, uint256 sellOrderId)
        external
        view
        override
        returns (bool, uint256 matchAmount)
    {
        IOrderManager.Order memory buyOrder = orderManager.getOrder(buyOrderId);
        IOrderManager.Order memory sellOrder = orderManager.getOrder(sellOrderId);

        if (!_canMatchOrders(buyOrder, sellOrder, true)) {
            return (false, 0);
        }

        uint256 buyRemaining = buyOrder.amount - buyOrder.filled;
        uint256 sellRemaining = sellOrder.amount - sellOrder.filled;
        matchAmount = buyRemaining < sellRemaining ? buyRemaining : sellRemaining;

        return (true, matchAmount);
    }

    function getBestPrice(IOrderManager.OrderType orderType)
        external
        view
        override
        returns (uint256)
    {
        uint256[] memory orders = orderManager.getOrdersByType(orderType);

        if (orders.length == 0) return 0;

        uint256 bestPrice = 0;
        bool found = false;

        for (uint256 i = 0; i < orders.length; i++) {
            IOrderManager.Order memory order = orderManager.getOrder(orders[i]);
            if (order.status == IOrderManager.OrderStatus.ACTIVE) {
                if (!found) {
                    bestPrice = order.price;
                    found = true;
                } else {
                    if (orderType == IOrderManager.OrderType.BUY) {
                        if (order.price > bestPrice) bestPrice = order.price;
                    } else {
                        if (order.price < bestPrice) bestPrice = order.price;
                    }
                }
            }
        }

        return bestPrice;
    }

    // ============ Internal Functions ============

    function _matchOrders(
        uint256 orderId,
        uint256[] memory oppositeOrders,
        bool isBuyOrder
    ) internal returns (MatchResult[] memory) {
        IOrderManager.Order memory order = orderManager.getOrder(orderId);
        uint256 remainingAmount = order.amount - order.filled;

        // Count potential matches
        uint256 matchCount = 0;
        for (uint256 i = 0; i < oppositeOrders.length && remainingAmount > 0; i++) {
            IOrderManager.Order memory oppositeOrder = orderManager.getOrder(oppositeOrders[i]);
            if (_canMatchOrders(order, oppositeOrder, isBuyOrder)) {
                matchCount++;
                uint256 matchAmount = _calculateMatchAmount(remainingAmount, oppositeOrder);
                remainingAmount -= matchAmount;
            }
        }

        // Create results array
        MatchResult[] memory results = new MatchResult[](matchCount);
        remainingAmount = order.amount - order.filled;
        uint256 resultIndex = 0;

        for (uint256 i = 0; i < oppositeOrders.length && remainingAmount > 0; i++) {
            IOrderManager.Order memory oppositeOrder = orderManager.getOrder(oppositeOrders[i]);
            if (_canMatchOrders(order, oppositeOrder, isBuyOrder)) {
                uint256 matchAmount = _calculateMatchAmount(remainingAmount, oppositeOrder);
                uint256 executionPrice = oppositeOrder.price; // Maker price

                // Execute trade (will be done by RWA module)
                results[resultIndex] = MatchResult({
                    tradeId: 0, // Will be set by executeTrade
                    matchedAmount: matchAmount,
                    executionPrice: executionPrice,
                    fullyMatched: (remainingAmount == matchAmount)
                });

                remainingAmount -= matchAmount;
                resultIndex++;
            }
        }

        return results;
    }

    function _canMatchOrders(
        IOrderManager.Order memory order1,
        IOrderManager.Order memory order2,
        bool order1IsBuy
    ) internal pure returns (bool) {
        if (order2.status != IOrderManager.OrderStatus.ACTIVE) return false;
        if (order2.filled >= order2.amount) return false;

        // Price matching: buy price >= sell price
        if (order1IsBuy) {
            return order1.price >= order2.price;
        } else {
            return order1.price <= order2.price;
        }
    }

    function _calculateMatchAmount(uint256 remainingAmount, IOrderManager.Order memory oppositeOrder)
        internal
        pure
        returns (uint256)
    {
        uint256 oppositeRemaining = oppositeOrder.amount - oppositeOrder.filled;
        return remainingAmount < oppositeRemaining ? remainingAmount : oppositeRemaining;
    }

    function setMaxRecentTrades(uint256 max) external onlyOwner {
        MatchingStorage storage ms = _getStorage();
        ms.maxRecentTrades = max;
    }
}
