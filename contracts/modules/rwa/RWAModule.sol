// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/core/IRWAModule.sol";
import "../../plugins/core/BaseModule.sol";
import "../../rwa/OrderBook.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title RWAModule
 * @notice Pluggable RWA trading module - adapter for existing OrderBook
 * @dev Wraps OrderBook and implements IRWAModule interface
 */
contract RWAModule is IRWAModule, BaseModule, Ownable {
    // ============ Constants ============

    bytes32 public constant MODULE_ID = keccak256("RWA_MODULE");
    string public constant MODULE_NAME = "RWAModule";
    string public constant MODULE_VERSION = "1.0.0";

    // ============ State Variables ============

    OrderBook public immutable orderBook;
    address public immutable baseToken;
    address public immutable quoteToken;

    // Trade history tracking
    Trade[] private _tradeHistory;
    mapping(address => uint256[]) private _userTradeIndices;

    // ============ Constructor ============

    constructor(address _orderBook)
        BaseModule(MODULE_ID, MODULE_NAME, MODULE_VERSION)
    {
        require(_orderBook != address(0), "Invalid order book");

        orderBook = OrderBook(_orderBook);
        baseToken = address(orderBook.baseToken());
        quoteToken = address(orderBook.quoteToken());
    }

    // ============ BaseModule Overrides ============

    function getDependencies() external pure override returns (bytes32[] memory) {
        bytes32[] memory deps = new bytes32[](2);
        deps[0] = keccak256("PRICE_ORACLE_MODULE");
        deps[1] = keccak256("AUDIT_MODULE");
        return deps;
    }

    function healthCheck()
        external
        view
        override(IModule, BaseModule)
        returns (bool healthy, string memory message)
    {
        (bool baseHealthy, string memory baseMessage) = BaseModule.healthCheck();
        if (!baseHealthy) {
            return (false, baseMessage);
        }

        // Check if order book is functional
        try orderBook.paused() returns (bool isPaused) {
            if (isPaused) {
                return (false, "OrderBook is paused");
            }
        } catch {
            return (false, "OrderBook not accessible");
        }

        return (true, "RWA module healthy");
    }

    // ============ Order Management ============

    function placeOrder(
        OrderType orderType,
        uint256 price,
        uint256 amount
    ) external override whenNotPaused returns (uint256 orderId) {
        _requireActive();

        OrderBook.OrderType obOrderType = orderType == OrderType.BUY
            ? OrderBook.OrderType.BUY
            : OrderBook.OrderType.SELL;

        orderId = orderBook.placeOrder(obOrderType, price, amount);

        emit OrderPlaced(orderId, msg.sender, orderType, price, amount);

        return orderId;
    }

    function cancelOrder(uint256 orderId) external override whenNotPaused {
        _requireActive();

        (,,,,, uint256 filled,, OrderBook.OrderStatus status) = orderBook.orders(orderId);

        orderBook.cancelOrder(orderId);

        uint256 remainingAmount = 0; // Would calculate from order
        emit OrderCancelled(orderId, msg.sender, remainingAmount);
    }

    function getOrder(uint256 orderId) external view override returns (Order memory) {
        (
            uint256 id,
            address trader,
            OrderBook.OrderType obOrderType,
            uint256 price,
            uint256 amount,
            uint256 filled,
            uint256 timestamp,
            OrderBook.OrderStatus obStatus
        ) = orderBook.orders(orderId);

        OrderType orderType = obOrderType == OrderBook.OrderType.BUY
            ? OrderType.BUY
            : OrderType.SELL;

        OrderStatus status;
        if (obStatus == OrderBook.OrderStatus.OPEN) {
            status = OrderStatus.OPEN;
        } else if (obStatus == OrderBook.OrderStatus.PARTIALLY_FILLED) {
            status = OrderStatus.PARTIALLY_FILLED;
        } else if (obStatus == OrderBook.OrderStatus.FILLED) {
            status = OrderStatus.FILLED;
        } else {
            status = OrderStatus.CANCELLED;
        }

        return Order({
            orderId: id,
            trader: trader,
            orderType: orderType,
            price: price,
            amount: amount,
            filled: filled,
            timestamp: timestamp,
            status: status
        });
    }

    function getUserOpenOrders(address user)
        external
        view
        override
        returns (Order[] memory)
    {
        OrderBook.Order[] memory obOrders = orderBook.getUserOpenOrders(user);
        Order[] memory orders = new Order[](obOrders.length);

        for (uint256 i = 0; i < obOrders.length; i++) {
            orders[i] = _convertOrder(obOrders[i]);
        }

        return orders;
    }

    function getUserOrderHistory(address user, uint256 offset, uint256 limit)
        external
        view
        override
        returns (Order[] memory)
    {
        // Stub: would need to track order history
        return new Order[](0);
    }

    // ============ Order Book ============

    function getOrderBookDepth(OrderType orderType, uint256 depth)
        external
        view
        override
        returns (uint256[] memory prices, uint256[] memory amounts)
    {
        OrderBook.OrderType obOrderType = orderType == OrderType.BUY
            ? OrderBook.OrderType.BUY
            : OrderBook.OrderType.SELL;

        OrderBook.Order[] memory orders = orderBook.getOrderBook(obOrderType, depth);

        prices = new uint256[](orders.length);
        amounts = new uint256[](orders.length);

        for (uint256 i = 0; i < orders.length; i++) {
            prices[i] = orders[i].price;
            amounts[i] = orders[i].amount - orders[i].filled;
        }

        return (prices, amounts);
    }

    function getBestBid() external view override returns (uint256 price, uint256 amount) {
        OrderBook.Order[] memory buyOrders = orderBook.getOrderBook(OrderBook.OrderType.BUY, 1);

        if (buyOrders.length > 0) {
            price = buyOrders[0].price;
            amount = buyOrders[0].amount - buyOrders[0].filled;
        }

        return (price, amount);
    }

    function getBestAsk() external view override returns (uint256 price, uint256 amount) {
        OrderBook.Order[] memory sellOrders = orderBook.getOrderBook(OrderBook.OrderType.SELL, 1);

        if (sellOrders.length > 0) {
            price = sellOrders[0].price;
            amount = sellOrders[0].amount - sellOrders[0].filled;
        }

        return (price, amount);
    }

    function getSpread() external view override returns (uint256 spread) {
        (uint256 bestBidPrice,) = this.getBestBid();
        (uint256 bestAskPrice,) = this.getBestAsk();

        if (bestBidPrice > 0 && bestAskPrice > 0) {
            spread = bestAskPrice - bestBidPrice;
        }

        return spread;
    }

    function getMidPrice() external view override returns (uint256) {
        (uint256 bestBidPrice,) = this.getBestBid();
        (uint256 bestAskPrice,) = this.getBestAsk();

        if (bestBidPrice > 0 && bestAskPrice > 0) {
            return (bestBidPrice + bestAskPrice) / 2;
        }

        return 0;
    }

    // ============ Trading ============

    function placeMarketOrder(OrderType orderType, uint256 amount)
        external
        override
        whenNotPaused
        returns (uint256 executedAmount, uint256 averagePrice)
    {
        _requireActive();

        // For market orders, use best available price
        (uint256 targetPrice,) = orderType == OrderType.BUY
            ? this.getBestAsk()
            : this.getBestBid();

        require(targetPrice > 0, "No liquidity available");

        OrderBook.OrderType obOrderType = orderType == OrderType.BUY
            ? OrderBook.OrderType.BUY
            : OrderBook.OrderType.SELL;

        // Place limit order at best price (acts as market order)
        orderBook.placeOrder(obOrderType, targetPrice, amount);

        // Return executed amount and price (simplified)
        return (amount, targetPrice);
    }

    function matchOrders() external override {
        // OrderBook handles matching internally
        // This is a no-op for the adapter
    }

    function getTradeHistory(uint256 offset, uint256 limit)
        external
        view
        override
        returns (Trade[] memory)
    {
        uint256 length = _tradeHistory.length;

        if (offset >= length) {
            return new Trade[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }

        Trade[] memory trades = new Trade[](end - offset);

        for (uint256 i = offset; i < end; i++) {
            trades[i - offset] = _tradeHistory[i];
        }

        return trades;
    }

    function getUserTrades(address user, uint256 offset, uint256 limit)
        external
        view
        override
        returns (Trade[] memory)
    {
        uint256[] memory indices = _userTradeIndices[user];
        uint256 length = indices.length;

        if (offset >= length) {
            return new Trade[](0);
        }

        uint256 end = offset + limit;
        if (end > length) {
            end = length;
        }

        Trade[] memory trades = new Trade[](end - offset);

        for (uint256 i = offset; i < end; i++) {
            trades[i - offset] = _tradeHistory[indices[i]];
        }

        return trades;
    }

    // ============ Market Statistics ============

    function getMarketStats() external view override returns (MarketStats memory) {
        (
            uint256 lastPrice,
            uint256 highPrice24h,
            uint256 lowPrice24h,
            uint256 volume24h,
            uint256 totalVolume,
            uint256 buyOrdersCount,
            uint256 sellOrdersCount
        ) = orderBook.getMarketStats();

        return MarketStats({
            lastPrice: lastPrice,
            highPrice24h: highPrice24h,
            lowPrice24h: lowPrice24h,
            volume24h: volume24h,
            totalVolume: totalVolume,
            openBuyOrders: buyOrdersCount,
            openSellOrders: sellOrdersCount,
            lastTradeTime: orderBook.lastTradeTimestamp()
        });
    }

    function getTradingPair() external view override returns (TradingPair memory) {
        return TradingPair({
            baseToken: baseToken,
            quoteToken: quoteToken,
            minOrderSize: orderBook.minOrderSize(),
            maxOrderSize: orderBook.maxOrderSize(),
            minPrice: orderBook.minPrice(),
            maxPrice: orderBook.maxPrice(),
            makerFee: orderBook.makerFee(),
            takerFee: orderBook.takerFee(),
            isActive: !orderBook.paused()
        });
    }

    function get24hVolume() external view override returns (uint256) {
        (, , , uint256 volume24h, , , ) = orderBook.getMarketStats();
        return volume24h;
    }

    function getTotalVolume() external view override returns (uint256) {
        return orderBook.totalVolume();
    }

    function getLastPrice() external view override returns (uint256) {
        return orderBook.lastPrice();
    }

    // ============ Configuration ============

    function updateFees(uint256 makerFee, uint256 takerFee)
        external
        override
        onlyOwner
    {
        orderBook.updateFees(makerFee, takerFee);
        emit FeesUpdated(makerFee, takerFee);
    }

    function updateOrderLimits(uint256 minSize, uint256 maxSize)
        external
        override
        onlyOwner
    {
        orderBook.updateLimits(minSize, maxSize, orderBook.minPrice(), orderBook.maxPrice());
        emit OrderLimitsUpdated(minSize, maxSize);
    }

    function updatePriceLimits(uint256 minPrice, uint256 maxPrice)
        external
        override
        onlyOwner
    {
        orderBook.updateLimits(
            orderBook.minOrderSize(),
            orderBook.maxOrderSize(),
            minPrice,
            maxPrice
        );
    }

    function setFeeCollector(address feeCollector)
        external
        override
        onlyOwner
    {
        orderBook.updateFeeCollector(feeCollector);
    }

    // ============ Liquidity ============

    function getTotalLiquidity()
        external
        view
        override
        returns (uint256 buyLiquidity, uint256 sellLiquidity)
    {
        // Calculate total liquidity from order books
        OrderBook.Order[] memory buyOrders = orderBook.getOrderBook(
            OrderBook.OrderType.BUY,
            1000
        );
        OrderBook.Order[] memory sellOrders = orderBook.getOrderBook(
            OrderBook.OrderType.SELL,
            1000
        );

        for (uint256 i = 0; i < buyOrders.length; i++) {
            buyLiquidity += (buyOrders[i].amount - buyOrders[i].filled);
        }

        for (uint256 i = 0; i < sellOrders.length; i++) {
            sellLiquidity += (sellOrders[i].amount - sellOrders[i].filled);
        }

        return (buyLiquidity, sellLiquidity);
    }

    function checkLiquidity(
        OrderType orderType,
        uint256 amount,
        uint256 maxSlippage
    ) external view override returns (bool sufficient, uint256 estimatedPrice) {
        OrderBook.OrderType obOrderType = orderType == OrderType.BUY
            ? OrderBook.OrderType.BUY
            : OrderBook.OrderType.SELL;

        OrderBook.Order[] memory orders = orderBook.getOrderBook(obOrderType, 100);

        uint256 availableLiquidity = 0;
        uint256 totalCost = 0;

        for (uint256 i = 0; i < orders.length; i++) {
            uint256 orderAmount = orders[i].amount - orders[i].filled;
            availableLiquidity += orderAmount;
            totalCost += orderAmount * orders[i].price;

            if (availableLiquidity >= amount) {
                break;
            }
        }

        sufficient = availableLiquidity >= amount;
        estimatedPrice = sufficient ? totalCost / amount : 0;

        return (sufficient, estimatedPrice);
    }

    // ============ Token Information ============

    function getBaseToken() external view override returns (address) {
        return baseToken;
    }

    function getQuoteToken() external view override returns (address) {
        return quoteToken;
    }

    // ============ Internal Helper Functions ============

    function _convertOrder(OrderBook.Order memory obOrder)
        internal
        pure
        returns (Order memory)
    {
        OrderType orderType = obOrder.orderType == OrderBook.OrderType.BUY
            ? OrderType.BUY
            : OrderType.SELL;

        OrderStatus status;
        if (obOrder.status == OrderBook.OrderStatus.OPEN) {
            status = OrderStatus.OPEN;
        } else if (obOrder.status == OrderBook.OrderStatus.PARTIALLY_FILLED) {
            status = OrderStatus.PARTIALLY_FILLED;
        } else if (obOrder.status == OrderBook.OrderStatus.FILLED) {
            status = OrderStatus.FILLED;
        } else {
            status = OrderStatus.CANCELLED;
        }

        return Order({
            orderId: obOrder.id,
            trader: obOrder.trader,
            orderType: orderType,
            price: obOrder.price,
            amount: obOrder.amount,
            filled: obOrder.filled,
            timestamp: obOrder.timestamp,
            status: status
        });
    }

    // ============ Override Required Functions ============

    function pause() external override(IModule, BaseModule) onlyOwner {
        BaseModule.pause();
    }

    function unpause() external override(IModule, BaseModule) onlyOwner {
        BaseModule.unpause();
    }
}
