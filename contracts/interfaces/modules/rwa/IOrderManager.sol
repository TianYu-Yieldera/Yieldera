// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IOrderManager
 * @notice Interface for order creation and management
 * @dev Handles all order lifecycle operations
 */
interface IOrderManager {
    // ============ Enums ============

    enum OrderType {
        BUY,
        SELL
    }

    enum OrderStatus {
        ACTIVE,
        FILLED,
        CANCELLED,
        EXPIRED
    }

    // ============ Structs ============

    struct Order {
        uint256 orderId;
        address trader;
        OrderType orderType;
        uint256 amount;
        uint256 price;
        uint256 filled;
        uint256 timestamp;
        OrderStatus status;
    }

    // ============ Events ============

    event OrderCreated(
        uint256 indexed orderId,
        address indexed trader,
        OrderType orderType,
        uint256 amount,
        uint256 price
    );
    event OrderCancelled(uint256 indexed orderId, address indexed trader);
    event OrderFilled(uint256 indexed orderId, uint256 amount, uint256 remaining);
    event OrderExpired(uint256 indexed orderId);

    // ============ Functions ============

    /**
     * @notice Create a new order
     * @param trader Trader address
     * @param orderType Order type (BUY/SELL)
     * @param amount Order amount
     * @param price Order price
     * @return orderId Created order ID
     */
    function createOrder(address trader, OrderType orderType, uint256 amount, uint256 price)
        external
        returns (uint256 orderId);

    /**
     * @notice Cancel an order
     * @param orderId Order ID to cancel
     * @param trader Trader address (for authorization)
     */
    function cancelOrder(uint256 orderId, address trader) external;

    /**
     * @notice Fill an order (partial or full)
     * @param orderId Order ID
     * @param amount Amount to fill
     */
    function fillOrder(uint256 orderId, uint256 amount) external;

    /**
     * @notice Get order details
     * @param orderId Order ID
     * @return Order data
     */
    function getOrder(uint256 orderId) external view returns (Order memory);

    /**
     * @notice Get user's active orders
     * @param trader Trader address
     * @return Array of order IDs
     */
    function getUserOrders(address trader) external view returns (uint256[] memory);

    /**
     * @notice Get all active orders
     * @return Array of order IDs
     */
    function getActiveOrders() external view returns (uint256[] memory);

    /**
     * @notice Get orders by type
     * @param orderType Order type
     * @return Array of order IDs
     */
    function getOrdersByType(OrderType orderType) external view returns (uint256[] memory);

    /**
     * @notice Check if order exists and is active
     * @param orderId Order ID
     * @return True if order is active
     */
    function isOrderActive(uint256 orderId) external view returns (bool);

    /**
     * @notice Get next order ID
     * @return Next available order ID
     */
    function getNextOrderId() external view returns (uint256);
}
