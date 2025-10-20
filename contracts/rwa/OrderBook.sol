// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

/**
 * @title OrderBook
 * @notice Limit order matching engine for RWA token trading
 * @dev Implements price-time priority matching with partial fills
 */
contract OrderBook is ReentrancyGuard, Ownable, Pausable {
    using SafeERC20 for IERC20;

    // Order types
    enum OrderType { BUY, SELL }
    enum OrderStatus { OPEN, PARTIALLY_FILLED, FILLED, CANCELLED }

    // Order structure
    struct Order {
        uint256 id;
        address trader;
        OrderType orderType;
        uint256 price; // Price per token in quote currency (6 decimals)
        uint256 amount; // Amount of tokens
        uint256 filled; // Amount filled
        uint256 timestamp;
        OrderStatus status;
    }

    // Trading pair
    IERC20 public immutable baseToken; // RWA Token
    IERC20 public immutable quoteToken; // LUSD or USDC
    uint256 public constant PRICE_DECIMALS = 1e6; // 6 decimals for price

    // Order storage
    uint256 public nextOrderId = 1;
    mapping(uint256 => Order) public orders;
    mapping(address => uint256[]) public userOrders;

    // Order books (sorted by price)
    uint256[] public buyOrders; // Sorted descending (highest first)
    uint256[] public sellOrders; // Sorted ascending (lowest first)

    // Trading limits
    uint256 public minOrderSize = 1e18; // 1 token minimum
    uint256 public maxOrderSize = 100000e18; // 100k tokens maximum
    uint256 public minPrice = 1e4; // 0.01 USD minimum
    uint256 public maxPrice = 1000000e6; // 1M USD maximum

    // Fees
    uint256 public makerFee = 10; // 0.1% (10 basis points)
    uint256 public takerFee = 30; // 0.3% (30 basis points)
    address public feeCollector;

    // Statistics
    uint256 public totalVolume;
    uint256 public lastPrice;
    uint256 public highPrice24h;
    uint256 public lowPrice24h;
    uint256 public volume24h;
    uint256 public lastTradeTimestamp;

    // Events
    event OrderPlaced(
        uint256 indexed orderId,
        address indexed trader,
        OrderType orderType,
        uint256 price,
        uint256 amount
    );
    event OrderCancelled(uint256 indexed orderId, address indexed trader);
    event OrderMatched(
        uint256 indexed buyOrderId,
        uint256 indexed sellOrderId,
        uint256 price,
        uint256 amount
    );
    event Trade(
        address indexed buyer,
        address indexed seller,
        uint256 price,
        uint256 amount,
        uint256 timestamp
    );
    event FeesUpdated(uint256 makerFee, uint256 takerFee);
    event LimitsUpdated(uint256 minSize, uint256 maxSize, uint256 minPrice, uint256 maxPrice);

    constructor(
        address _baseToken,
        address _quoteToken,
        address _feeCollector
    ) {
        require(_baseToken != address(0), "Invalid base token");
        require(_quoteToken != address(0), "Invalid quote token");
        require(_feeCollector != address(0), "Invalid fee collector");

        baseToken = IERC20(_baseToken);
        quoteToken = IERC20(_quoteToken);
        feeCollector = _feeCollector;
    }

    /**
     * @notice Place a limit order
     * @param orderType Buy or Sell
     * @param price Price per token
     * @param amount Amount of tokens
     */
    function placeOrder(
        OrderType orderType,
        uint256 price,
        uint256 amount
    ) external nonReentrant whenNotPaused returns (uint256 orderId) {
        // Validate order
        require(amount >= minOrderSize && amount <= maxOrderSize, "Invalid order size");
        require(price >= minPrice && price <= maxPrice, "Invalid price");

        orderId = nextOrderId++;

        // Create order
        orders[orderId] = Order({
            id: orderId,
            trader: msg.sender,
            orderType: orderType,
            price: price,
            amount: amount,
            filled: 0,
            timestamp: block.timestamp,
            status: OrderStatus.OPEN
        });

        userOrders[msg.sender].push(orderId);

        // Lock funds
        if (orderType == OrderType.BUY) {
            uint256 totalCost = (amount * price) / PRICE_DECIMALS;
            quoteToken.safeTransferFrom(msg.sender, address(this), totalCost);
        } else {
            baseToken.safeTransferFrom(msg.sender, address(this), amount);
        }

        // Add to order book
        _addToOrderBook(orderId);

        emit OrderPlaced(orderId, msg.sender, orderType, price, amount);

        // Try to match order
        _matchOrders();

        return orderId;
    }

    /**
     * @notice Cancel an open order
     * @param orderId Order ID to cancel
     */
    function cancelOrder(uint256 orderId) external nonReentrant {
        Order storage order = orders[orderId];
        require(order.trader == msg.sender, "Not order owner");
        require(order.status == OrderStatus.OPEN || order.status == OrderStatus.PARTIALLY_FILLED, "Order not cancellable");

        uint256 remainingAmount = order.amount - order.filled;

        // Update status
        order.status = OrderStatus.CANCELLED;

        // Remove from order book
        _removeFromOrderBook(orderId);

        // Refund remaining funds
        if (order.orderType == OrderType.BUY) {
            uint256 refund = (remainingAmount * order.price) / PRICE_DECIMALS;
            quoteToken.safeTransfer(msg.sender, refund);
        } else {
            baseToken.safeTransfer(msg.sender, remainingAmount);
        }

        emit OrderCancelled(orderId, msg.sender);
    }

    /**
     * @notice Match orders in the order book
     */
    function _matchOrders() internal {
        while (buyOrders.length > 0 && sellOrders.length > 0) {
            uint256 buyOrderId = buyOrders[0];
            uint256 sellOrderId = sellOrders[0];

            Order storage buyOrder = orders[buyOrderId];
            Order storage sellOrder = orders[sellOrderId];

            // Check if orders can match
            if (buyOrder.price < sellOrder.price) {
                break; // No match possible
            }

            // Calculate match amount
            uint256 buyRemaining = buyOrder.amount - buyOrder.filled;
            uint256 sellRemaining = sellOrder.amount - sellOrder.filled;
            uint256 matchAmount = buyRemaining < sellRemaining ? buyRemaining : sellRemaining;

            // Use sell price (price-time priority)
            uint256 matchPrice = sellOrder.price;

            // Update orders
            buyOrder.filled += matchAmount;
            sellOrder.filled += matchAmount;

            // Update order status
            if (buyOrder.filled == buyOrder.amount) {
                buyOrder.status = OrderStatus.FILLED;
                _removeFromOrderBook(buyOrderId);
            } else {
                buyOrder.status = OrderStatus.PARTIALLY_FILLED;
            }

            if (sellOrder.filled == sellOrder.amount) {
                sellOrder.status = OrderStatus.FILLED;
                _removeFromOrderBook(sellOrderId);
            } else {
                sellOrder.status = OrderStatus.PARTIALLY_FILLED;
            }

            // Execute trade
            _executeTrade(
                buyOrder.trader,
                sellOrder.trader,
                matchAmount,
                matchPrice
            );

            // Update statistics
            _updateStatistics(matchPrice, matchAmount);

            emit OrderMatched(buyOrderId, sellOrderId, matchPrice, matchAmount);
            emit Trade(buyOrder.trader, sellOrder.trader, matchPrice, matchAmount, block.timestamp);
        }
    }

    /**
     * @notice Execute a trade between buyer and seller
     */
    function _executeTrade(
        address buyer,
        address seller,
        uint256 amount,
        uint256 price
    ) internal {
        uint256 totalCost = (amount * price) / PRICE_DECIMALS;

        // Calculate fees
        uint256 buyerFee = (amount * takerFee) / 10000;
        uint256 sellerFee = (totalCost * makerFee) / 10000;

        // Transfer tokens
        baseToken.safeTransfer(buyer, amount - buyerFee);
        baseToken.safeTransfer(feeCollector, buyerFee);

        quoteToken.safeTransfer(seller, totalCost - sellerFee);
        quoteToken.safeTransfer(feeCollector, sellerFee);
    }

    /**
     * @notice Add order to order book maintaining price priority
     */
    function _addToOrderBook(uint256 orderId) internal {
        Order memory order = orders[orderId];

        if (order.orderType == OrderType.BUY) {
            // Insert in descending order (highest price first)
            uint256 insertIndex = 0;
            for (uint256 i = 0; i < buyOrders.length; i++) {
                if (order.price > orders[buyOrders[i]].price) {
                    break;
                }
                insertIndex = i + 1;
            }
            _insertAt(buyOrders, orderId, insertIndex);
        } else {
            // Insert in ascending order (lowest price first)
            uint256 insertIndex = 0;
            for (uint256 i = 0; i < sellOrders.length; i++) {
                if (order.price < orders[sellOrders[i]].price) {
                    break;
                }
                insertIndex = i + 1;
            }
            _insertAt(sellOrders, orderId, insertIndex);
        }
    }

    /**
     * @notice Remove order from order book
     */
    function _removeFromOrderBook(uint256 orderId) internal {
        Order memory order = orders[orderId];
        uint256[] storage orderBook = order.orderType == OrderType.BUY ? buyOrders : sellOrders;

        for (uint256 i = 0; i < orderBook.length; i++) {
            if (orderBook[i] == orderId) {
                orderBook[i] = orderBook[orderBook.length - 1];
                orderBook.pop();
                break;
            }
        }
    }

    /**
     * @notice Insert element at specific index in array
     */
    function _insertAt(uint256[] storage array, uint256 element, uint256 index) internal {
        array.push(0); // Extend array
        for (uint256 i = array.length - 1; i > index; i--) {
            array[i] = array[i - 1];
        }
        array[index] = element;
    }

    /**
     * @notice Update trading statistics
     */
    function _updateStatistics(uint256 price, uint256 amount) internal {
        lastPrice = price;
        lastTradeTimestamp = block.timestamp;
        totalVolume += amount;

        // Update 24h stats (simplified - in production would use time-weighted)
        if (block.timestamp - lastTradeTimestamp < 86400) {
            volume24h += amount;
            if (price > highPrice24h) highPrice24h = price;
            if (price < lowPrice24h || lowPrice24h == 0) lowPrice24h = price;
        } else {
            // Reset 24h stats
            volume24h = amount;
            highPrice24h = price;
            lowPrice24h = price;
        }
    }

    /**
     * @notice Get order book depth
     * @param orderType Buy or Sell
     * @param limit Number of orders to return
     */
    function getOrderBook(OrderType orderType, uint256 limit)
        external
        view
        returns (Order[] memory)
    {
        uint256[] storage orderBook = orderType == OrderType.BUY ? buyOrders : sellOrders;
        uint256 length = orderBook.length < limit ? orderBook.length : limit;

        Order[] memory result = new Order[](length);
        for (uint256 i = 0; i < length; i++) {
            result[i] = orders[orderBook[i]];
        }

        return result;
    }

    /**
     * @notice Get user's open orders
     */
    function getUserOpenOrders(address user) external view returns (Order[] memory) {
        uint256[] memory orderIds = userOrders[user];
        uint256 count = 0;

        // Count open orders
        for (uint256 i = 0; i < orderIds.length; i++) {
            Order memory order = orders[orderIds[i]];
            if (order.status == OrderStatus.OPEN || order.status == OrderStatus.PARTIALLY_FILLED) {
                count++;
            }
        }

        // Collect open orders
        Order[] memory result = new Order[](count);
        uint256 index = 0;

        for (uint256 i = 0; i < orderIds.length && index < count; i++) {
            Order memory order = orders[orderIds[i]];
            if (order.status == OrderStatus.OPEN || order.status == OrderStatus.PARTIALLY_FILLED) {
                result[index++] = order;
            }
        }

        return result;
    }

    /**
     * @notice Get market statistics
     */
    function getMarketStats() external view returns (
        uint256 _lastPrice,
        uint256 _highPrice24h,
        uint256 _lowPrice24h,
        uint256 _volume24h,
        uint256 _totalVolume,
        uint256 _buyOrdersCount,
        uint256 _sellOrdersCount
    ) {
        return (
            lastPrice,
            highPrice24h,
            lowPrice24h,
            volume24h,
            totalVolume,
            buyOrders.length,
            sellOrders.length
        );
    }

    /**
     * @notice Update trading fees
     */
    function updateFees(uint256 _makerFee, uint256 _takerFee) external onlyOwner {
        require(_makerFee <= 100, "Maker fee too high"); // Max 1%
        require(_takerFee <= 100, "Taker fee too high"); // Max 1%

        makerFee = _makerFee;
        takerFee = _takerFee;

        emit FeesUpdated(_makerFee, _takerFee);
    }

    /**
     * @notice Update order limits
     */
    function updateLimits(
        uint256 _minSize,
        uint256 _maxSize,
        uint256 _minPrice,
        uint256 _maxPrice
    ) external onlyOwner {
        minOrderSize = _minSize;
        maxOrderSize = _maxSize;
        minPrice = _minPrice;
        maxPrice = _maxPrice;

        emit LimitsUpdated(_minSize, _maxSize, _minPrice, _maxPrice);
    }

    /**
     * @notice Update fee collector
     */
    function updateFeeCollector(address _feeCollector) external onlyOwner {
        require(_feeCollector != address(0), "Invalid fee collector");
        feeCollector = _feeCollector;
    }

    /**
     * @notice Pause trading
     */
    function pause() external onlyOwner {
        _pause();
    }

    /**
     * @notice Resume trading
     */
    function unpause() external onlyOwner {
        _unpause();
    }
}