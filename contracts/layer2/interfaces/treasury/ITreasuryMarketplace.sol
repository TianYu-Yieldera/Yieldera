// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title ITreasuryMarketplace
 * @notice Interface for secondary market trading of US Treasury tokens
 */
interface ITreasuryMarketplace {
    /// @notice Order types
    enum OrderType {
        BUY,
        SELL
    }

    /// @notice Order status
    enum OrderStatus {
        Open,
        PartiallyFilled,
        Filled,
        Cancelled
    }

    /// @notice Market order structure
    struct MarketOrder {
        uint256 orderId;
        uint256 assetId;
        OrderType orderType;
        address user;
        uint256 tokenAmount;
        uint256 pricePerToken;      // Price in USD (18 decimals)
        uint256 totalValue;
        uint256 filledAmount;
        OrderStatus status;
        uint256 createdAt;
        uint256 expiresAt;
    }

    /// @notice Trade execution record
    struct Trade {
        uint256 tradeId;
        uint256 orderId;
        uint256 assetId;
        address buyer;
        address seller;
        uint256 amount;
        uint256 price;
        uint256 totalValue;
        uint256 fee;
        uint256 timestamp;
    }

    /// @notice Events
    event OrderCreated(
        uint256 indexed orderId,
        uint256 indexed assetId,
        address indexed user,
        OrderType orderType,
        uint256 tokenAmount,
        uint256 pricePerToken
    );

    event OrderMatched(
        uint256 indexed tradeId,
        uint256 indexed buyOrderId,
        uint256 indexed sellOrderId,
        address buyer,
        address seller,
        uint256 amount,
        uint256 price
    );

    event OrderCancelled(
        uint256 indexed orderId,
        address indexed user
    );

    event OrderPartiallyFilled(
        uint256 indexed orderId,
        uint256 filledAmount,
        uint256 remainingAmount
    );

    /// @notice Create buy order
    function createBuyOrder(
        uint256 assetId,
        uint256 tokenAmount,
        uint256 pricePerToken,
        uint256 duration
    ) external payable returns (uint256 orderId);

    /// @notice Create sell order
    function createSellOrder(
        uint256 assetId,
        uint256 tokenAmount,
        uint256 pricePerToken,
        uint256 duration
    ) external returns (uint256 orderId);

    /// @notice Cancel order
    function cancelOrder(uint256 orderId) external;

    /// @notice Match orders (called by matcher or anyone)
    function matchOrders(
        uint256 buyOrderId,
        uint256 sellOrderId,
        uint256 amount
    ) external returns (uint256 tradeId);

    /// @notice Get order details
    function getOrder(uint256 orderId) external view returns (MarketOrder memory);

    /// @notice Get open orders for asset
    function getOpenOrders(
        uint256 assetId,
        OrderType orderType
    ) external view returns (uint256[] memory);

    /// @notice Get user orders
    function getUserOrders(address user) external view returns (uint256[] memory);

    /// @notice Get trade history
    function getTradeHistory(
        uint256 assetId,
        uint256 count
    ) external view returns (Trade[] memory);

    /// @notice Calculate trading fee
    function calculateFee(uint256 totalValue) external view returns (uint256);
}
