// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IModule.sol";

/**
 * @title IRWAModule
 * @notice Standard interface for Real World Asset trading modules
 * @dev Extends IModule with RWA-specific functionality
 *
 * Purpose:
 * - Manage RWA token listings
 * - Handle order book operations
 * - Execute trades
 * - Track market statistics
 */
interface IRWAModule is IModule {
    /**
     * @notice Order types
     */
    enum OrderType { BUY, SELL }

    /**
     * @notice Order status
     */
    enum OrderStatus { OPEN, PARTIALLY_FILLED, FILLED, CANCELLED }

    /**
     * @notice Order structure
     */
    struct Order {
        uint256 orderId;            // Unique order identifier
        address trader;             // Trader address
        OrderType orderType;        // Buy or sell
        uint256 price;              // Price per token
        uint256 amount;             // Total amount
        uint256 filled;             // Amount filled
        uint256 timestamp;          // Order creation time
        OrderStatus status;         // Current status
    }

    /**
     * @notice Trading pair information
     */
    struct TradingPair {
        address baseToken;          // RWA token address
        address quoteToken;         // Quote currency (LUSD, USDC, etc.)
        uint256 minOrderSize;       // Minimum order size
        uint256 maxOrderSize;       // Maximum order size
        uint256 minPrice;           // Minimum price
        uint256 maxPrice;           // Maximum price
        uint256 makerFee;           // Maker fee in basis points
        uint256 takerFee;           // Taker fee in basis points
        bool isActive;              // Trading pair status
    }

    /**
     * @notice Market statistics
     */
    struct MarketStats {
        uint256 lastPrice;          // Last trade price
        uint256 highPrice24h;       // 24h high
        uint256 lowPrice24h;        // 24h low
        uint256 volume24h;          // 24h volume
        uint256 totalVolume;        // All-time volume
        uint256 openBuyOrders;      // Number of open buy orders
        uint256 openSellOrders;     // Number of open sell orders
        uint256 lastTradeTime;      // Timestamp of last trade
    }

    /**
     * @notice Trade execution data
     */
    struct Trade {
        uint256 tradeId;            // Unique trade identifier
        uint256 buyOrderId;         // Buy order ID
        uint256 sellOrderId;        // Sell order ID
        address buyer;              // Buyer address
        address seller;             // Seller address
        uint256 price;              // Execution price
        uint256 amount;             // Trade amount
        uint256 timestamp;          // Trade timestamp
    }

    // ============ Events ============

    event OrderPlaced(
        uint256 indexed orderId,
        address indexed trader,
        OrderType orderType,
        uint256 price,
        uint256 amount
    );

    event OrderCancelled(
        uint256 indexed orderId,
        address indexed trader,
        uint256 remainingAmount
    );

    event OrderMatched(
        uint256 indexed buyOrderId,
        uint256 indexed sellOrderId,
        uint256 price,
        uint256 amount
    );

    event TradeExecuted(
        uint256 indexed tradeId,
        address indexed buyer,
        address indexed seller,
        uint256 price,
        uint256 amount
    );

    event TradingPairListed(
        address indexed baseToken,
        address indexed quoteToken,
        string name
    );

    event TradingPairDelisted(
        address indexed baseToken,
        address indexed quoteToken
    );

    event FeesUpdated(uint256 makerFee, uint256 takerFee);
    event OrderLimitsUpdated(uint256 minSize, uint256 maxSize);

    // ============ Order Management ============

    /**
     * @notice Place a limit order
     * @param orderType Buy or sell
     * @param price Price per token
     * @param amount Amount of tokens
     * @return orderId Unique order identifier
     * @dev Locks required funds (quote currency for buy, base token for sell)
     */
    function placeOrder(
        OrderType orderType,
        uint256 price,
        uint256 amount
    ) external returns (uint256 orderId);

    /**
     * @notice Cancel an open order
     * @param orderId Order identifier
     * @dev Refunds remaining locked funds
     */
    function cancelOrder(uint256 orderId) external;

    /**
     * @notice Get order details
     * @param orderId Order identifier
     * @return Order structure
     */
    function getOrder(uint256 orderId) external view returns (Order memory);

    /**
     * @notice Get user's open orders
     * @param user User address
     * @return Array of open orders
     */
    function getUserOpenOrders(address user) external view returns (Order[] memory);

    /**
     * @notice Get user's order history
     * @param user User address
     * @param offset Starting index
     * @param limit Number of orders to return
     * @return Array of historical orders
     */
    function getUserOrderHistory(
        address user,
        uint256 offset,
        uint256 limit
    ) external view returns (Order[] memory);

    // ============ Order Book ============

    /**
     * @notice Get order book depth
     * @param orderType Buy or sell side
     * @param depth Number of price levels to return
     * @return prices Array of prices
     * @return amounts Array of amounts at each price
     */
    function getOrderBookDepth(OrderType orderType, uint256 depth)
        external
        view
        returns (uint256[] memory prices, uint256[] memory amounts);

    /**
     * @notice Get best bid (highest buy price)
     * @return price Best bid price
     * @return amount Total amount at best bid
     */
    function getBestBid() external view returns (uint256 price, uint256 amount);

    /**
     * @notice Get best ask (lowest sell price)
     * @return price Best ask price
     * @return amount Total amount at best ask
     */
    function getBestAsk() external view returns (uint256 price, uint256 amount);

    /**
     * @notice Get bid-ask spread
     * @return spread Difference between best ask and best bid
     */
    function getSpread() external view returns (uint256 spread);

    /**
     * @notice Get mid price (average of best bid and ask)
     * @return Mid market price
     */
    function getMidPrice() external view returns (uint256);

    // ============ Trading ============

    /**
     * @notice Place market order (immediate execution at best price)
     * @param orderType Buy or sell
     * @param amount Amount of tokens
     * @return executedAmount Amount actually traded
     * @return averagePrice Average execution price
     * @dev May result in partial fill if insufficient liquidity
     */
    function placeMarketOrder(OrderType orderType, uint256 amount)
        external
        returns (uint256 executedAmount, uint256 averagePrice);

    /**
     * @notice Execute pending order matches
     * @dev Public function to trigger order matching (can be called by anyone)
     */
    function matchOrders() external;

    /**
     * @notice Get trade history
     * @param offset Starting index
     * @param limit Number of trades to return
     * @return Array of recent trades
     */
    function getTradeHistory(uint256 offset, uint256 limit)
        external
        view
        returns (Trade[] memory);

    /**
     * @notice Get user's trade history
     * @param user User address
     * @param offset Starting index
     * @param limit Number of trades to return
     * @return Array of user's trades
     */
    function getUserTrades(address user, uint256 offset, uint256 limit)
        external
        view
        returns (Trade[] memory);

    // ============ Market Statistics ============

    /**
     * @notice Get current market statistics
     * @return Market statistics structure
     */
    function getMarketStats() external view returns (MarketStats memory);

    /**
     * @notice Get trading pair information
     * @return Trading pair structure
     */
    function getTradingPair() external view returns (TradingPair memory);

    /**
     * @notice Get 24h volume
     * @return Volume in last 24 hours
     */
    function get24hVolume() external view returns (uint256);

    /**
     * @notice Get all-time volume
     * @return Total volume since inception
     */
    function getTotalVolume() external view returns (uint256);

    /**
     * @notice Get last trade price
     * @return Last execution price
     */
    function getLastPrice() external view returns (uint256);

    // ============ Configuration ============

    /**
     * @notice Update trading fees
     * @param makerFee New maker fee in basis points
     * @param takerFee New taker fee in basis points
     * @dev Only callable by governance
     */
    function updateFees(uint256 makerFee, uint256 takerFee) external;

    /**
     * @notice Update order size limits
     * @param minSize Minimum order size
     * @param maxSize Maximum order size
     * @dev Only callable by governance
     */
    function updateOrderLimits(uint256 minSize, uint256 maxSize) external;

    /**
     * @notice Update price limits
     * @param minPrice Minimum allowed price
     * @param maxPrice Maximum allowed price
     * @dev Only callable by governance
     */
    function updatePriceLimits(uint256 minPrice, uint256 maxPrice) external;

    /**
     * @notice Set fee collector address
     * @param feeCollector Address to receive trading fees
     * @dev Only callable by governance
     */
    function setFeeCollector(address feeCollector) external;

    // ============ Liquidity ============

    /**
     * @notice Get total liquidity (sum of all orders)
     * @return buyLiquidity Total buy-side liquidity
     * @return sellLiquidity Total sell-side liquidity
     */
    function getTotalLiquidity()
        external
        view
        returns (uint256 buyLiquidity, uint256 sellLiquidity);

    /**
     * @notice Check if there is sufficient liquidity for an order
     * @param orderType Buy or sell
     * @param amount Desired amount
     * @param maxSlippage Maximum acceptable slippage (basis points)
     * @return sufficient True if liquidity is sufficient
     * @return estimatedPrice Estimated execution price
     */
    function checkLiquidity(
        OrderType orderType,
        uint256 amount,
        uint256 maxSlippage
    ) external view returns (bool sufficient, uint256 estimatedPrice);

    // ============ Token Information ============

    /**
     * @notice Get base token (RWA token) address
     * @return Base token address
     */
    function getBaseToken() external view returns (address);

    /**
     * @notice Get quote token (payment currency) address
     * @return Quote token address
     */
    function getQuoteToken() external view returns (address);
}
