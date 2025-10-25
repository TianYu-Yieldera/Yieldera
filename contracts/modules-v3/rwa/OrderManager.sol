// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/modules/rwa/IOrderManager.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title OrderManager
 * @notice Manages order creation and lifecycle
 * @dev Uses Diamond Storage for state management
 */
contract OrderManager is IOrderManager, Ownable {
    // ============ State Variables ============

    address public rwaModule; // Main coordinator contract

    // Storage for orders (using Diamond Storage pattern)
    bytes32 private constant ORDER_STORAGE_POSITION = keccak256("order.manager.storage");

    struct OrderStorage {
        mapping(uint256 => Order) orders;
        mapping(address => uint256[]) userOrders;
        uint256[] activeOrders;
        uint256[] buyOrders;
        uint256[] sellOrders;
        uint256 nextOrderId;
    }

    // ============ Modifiers ============

    modifier onlyRWAModule() {
        require(msg.sender == rwaModule, "Only RWA module");
        _;
    }

    // ============ Constructor ============

    constructor() {}

    // ============ Admin Functions ============

    function setRWAModule(address _rwaModule) external onlyOwner {
        require(_rwaModule != address(0), "Invalid address");
        rwaModule = _rwaModule;
    }

    // ============ Internal Storage ============

    function _getStorage() private pure returns (OrderStorage storage os) {
        bytes32 position = ORDER_STORAGE_POSITION;
        assembly {
            os.slot := position
        }
    }

    // ============ IOrderManager Implementation ============

    function createOrder(address trader, OrderType orderType, uint256 amount, uint256 price)
        external
        override
        onlyRWAModule
        returns (uint256 orderId)
    {
        require(trader != address(0), "Invalid trader");
        require(amount > 0, "Invalid amount");
        require(price > 0, "Invalid price");

        OrderStorage storage os = _getStorage();
        orderId = os.nextOrderId++;

        Order storage order = os.orders[orderId];
        order.orderId = orderId;
        order.trader = trader;
        order.orderType = orderType;
        order.amount = amount;
        order.price = price;
        order.filled = 0;
        order.timestamp = block.timestamp;
        order.status = OrderStatus.ACTIVE;

        os.userOrders[trader].push(orderId);
        os.activeOrders.push(orderId);

        if (orderType == OrderType.BUY) {
            os.buyOrders.push(orderId);
        } else {
            os.sellOrders.push(orderId);
        }

        emit OrderCreated(orderId, trader, orderType, amount, price);

        return orderId;
    }

    function cancelOrder(uint256 orderId, address trader) external override onlyRWAModule {
        OrderStorage storage os = _getStorage();
        Order storage order = os.orders[orderId];

        require(order.trader == trader, "Not order owner");
        require(order.status == OrderStatus.ACTIVE, "Order not active");

        order.status = OrderStatus.CANCELLED;

        emit OrderCancelled(orderId, trader);
    }

    function fillOrder(uint256 orderId, uint256 amount) external override onlyRWAModule {
        OrderStorage storage os = _getStorage();
        Order storage order = os.orders[orderId];

        require(order.status == OrderStatus.ACTIVE, "Order not active");
        require(order.filled + amount <= order.amount, "Exceeds order amount");

        order.filled += amount;

        uint256 remaining = order.amount - order.filled;
        if (remaining == 0) {
            order.status = OrderStatus.FILLED;
        }

        emit OrderFilled(orderId, amount, remaining);
    }

    function getOrder(uint256 orderId) external view override returns (Order memory) {
        OrderStorage storage os = _getStorage();
        return os.orders[orderId];
    }

    function getUserOrders(address trader) external view override returns (uint256[] memory) {
        OrderStorage storage os = _getStorage();
        return os.userOrders[trader];
    }

    function getActiveOrders() external view override returns (uint256[] memory) {
        OrderStorage storage os = _getStorage();
        return os.activeOrders;
    }

    function getOrdersByType(OrderType orderType) external view override returns (uint256[] memory) {
        OrderStorage storage os = _getStorage();
        return orderType == OrderType.BUY ? os.buyOrders : os.sellOrders;
    }

    function isOrderActive(uint256 orderId) external view override returns (bool) {
        OrderStorage storage os = _getStorage();
        return os.orders[orderId].status == OrderStatus.ACTIVE;
    }

    function getNextOrderId() external view override returns (uint256) {
        OrderStorage storage os = _getStorage();
        return os.nextOrderId;
    }
}
