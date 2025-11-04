// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "../interfaces/treasury/ITreasuryMarketplace.sol";
import "../interfaces/treasury/ITreasuryAsset.sol";

/**
 * @title TreasuryMarketplace
 * @notice Secondary market for trading US Treasury tokens
 * @dev Order book based marketplace with automatic matching
 *
 * Key Features:
 * - Buy and sell orders for treasury tokens
 * - Order matching and execution
 * - Fee collection on trades
 * - Trade history and analytics
 * - Support for partial fills
 */
contract TreasuryMarketplace is ITreasuryMarketplace, AccessControl, ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;

    bytes32 public constant MARKETPLACE_ADMIN_ROLE = keccak256("MARKETPLACE_ADMIN_ROLE");
    bytes32 public constant FEE_COLLECTOR_ROLE = keccak256("FEE_COLLECTOR_ROLE");
    bytes32 public constant MATCHER_ROLE = keccak256("MATCHER_ROLE");

    /// @notice Order counter
    uint256 private orderIdCounter;
    uint256 private tradeIdCounter;

    /// @notice Storage
    mapping(uint256 => MarketOrder) private orders;
    mapping(uint256 => uint256[]) private assetBuyOrders;  // assetId => orderIds
    mapping(uint256 => uint256[]) private assetSellOrders; // assetId => orderIds
    mapping(address => uint256[]) private userOrders;      // user => orderIds
    Trade[] private tradeHistory;
    mapping(uint256 => uint256[]) private assetTrades;     // assetId => tradeIds

    /// @notice Integration contracts
    ITreasuryAsset public immutable assetFactory;
    address public immutable paymentToken; // USDC or similar stablecoin

    /// @notice Fee configuration
    uint256 public feeRate = 25; // 0.25% (basis points)
    uint256 public constant MAX_FEE_RATE = 500; // 5% maximum
    uint256 public constant FEE_PRECISION = 10000;
    address public feeCollector;

    /// @notice Collected fees
    mapping(address => uint256) public collectedFees; // token => amount

    /// @notice Statistics
    uint256 public totalOrders;
    uint256 public totalTrades;
    uint256 public totalVolume; // Total USD volume

    /**
     * @notice Constructor
     * @param admin Admin address
     * @param factory TreasuryAssetFactory address
     * @param paymentToken_ Payment token address (USDC)
     * @param feeCollector_ Fee collector address
     */
    constructor(
        address admin,
        address factory,
        address paymentToken_,
        address feeCollector_
    ) {
        require(admin != address(0), "Invalid admin");
        require(factory != address(0), "Invalid factory");
        require(paymentToken_ != address(0), "Invalid payment token");
        require(feeCollector_ != address(0), "Invalid fee collector");

        assetFactory = ITreasuryAsset(factory);
        paymentToken = paymentToken_;
        feeCollector = feeCollector_;

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(MARKETPLACE_ADMIN_ROLE, admin);
        _grantRole(FEE_COLLECTOR_ROLE, admin);
        _grantRole(MATCHER_ROLE, admin);
    }

    // =============================================================
    //                     ORDER CREATION
    // =============================================================

    /**
     * @notice Create buy order
     * @param assetId Treasury asset ID
     * @param tokenAmount Amount of tokens to buy
     * @param pricePerToken Price per token in USD (18 decimals)
     * @param duration Order duration in seconds
     * @return orderId Order identifier
     */
    function createBuyOrder(
        uint256 assetId,
        uint256 tokenAmount,
        uint256 pricePerToken,
        uint256 duration
    ) external payable override whenNotPaused nonReentrant returns (uint256 orderId) {
        require(assetFactory.isAssetActive(assetId), "Asset not active");
        require(tokenAmount > 0, "Invalid amount");
        require(pricePerToken > 0, "Invalid price");
        require(duration > 0 && duration <= 30 days, "Invalid duration");

        uint256 totalValue = (tokenAmount * pricePerToken) / 1e18;
        require(totalValue > 0, "Total value too small");

        // Transfer payment token to escrow
        IERC20(paymentToken).safeTransferFrom(msg.sender, address(this), totalValue);

        orderIdCounter++;
        orderId = orderIdCounter;

        MarketOrder storage order = orders[orderId];
        order.orderId = orderId;
        order.assetId = assetId;
        order.orderType = OrderType.BUY;
        order.user = msg.sender;
        order.tokenAmount = tokenAmount;
        order.pricePerToken = pricePerToken;
        order.totalValue = totalValue;
        order.filledAmount = 0;
        order.status = OrderStatus.Open;
        order.createdAt = block.timestamp;
        order.expiresAt = block.timestamp + duration;

        assetBuyOrders[assetId].push(orderId);
        userOrders[msg.sender].push(orderId);
        totalOrders++;

        emit OrderCreated(orderId, assetId, msg.sender, OrderType.BUY, tokenAmount, pricePerToken);
    }

    /**
     * @notice Create sell order
     * @param assetId Treasury asset ID
     * @param tokenAmount Amount of tokens to sell
     * @param pricePerToken Price per token in USD (18 decimals)
     * @param duration Order duration in seconds
     * @return orderId Order identifier
     */
    function createSellOrder(
        uint256 assetId,
        uint256 tokenAmount,
        uint256 pricePerToken,
        uint256 duration
    ) external override whenNotPaused nonReentrant returns (uint256 orderId) {
        require(assetFactory.isAssetActive(assetId), "Asset not active");
        require(tokenAmount > 0, "Invalid amount");
        require(pricePerToken > 0, "Invalid price");
        require(duration > 0 && duration <= 30 days, "Invalid duration");

        // Get treasury token
        address tokenAddress = assetFactory.getTokenAddress(assetId);
        require(tokenAddress != address(0), "Token not found");

        // Transfer tokens to escrow
        IERC20(tokenAddress).safeTransferFrom(msg.sender, address(this), tokenAmount);

        orderIdCounter++;
        orderId = orderIdCounter;

        uint256 totalValue = (tokenAmount * pricePerToken) / 1e18;

        MarketOrder storage order = orders[orderId];
        order.orderId = orderId;
        order.assetId = assetId;
        order.orderType = OrderType.SELL;
        order.user = msg.sender;
        order.tokenAmount = tokenAmount;
        order.pricePerToken = pricePerToken;
        order.totalValue = totalValue;
        order.filledAmount = 0;
        order.status = OrderStatus.Open;
        order.createdAt = block.timestamp;
        order.expiresAt = block.timestamp + duration;

        assetSellOrders[assetId].push(orderId);
        userOrders[msg.sender].push(orderId);
        totalOrders++;

        emit OrderCreated(orderId, assetId, msg.sender, OrderType.SELL, tokenAmount, pricePerToken);
    }

    /**
     * @notice Cancel open order
     * @param orderId Order to cancel
     */
    function cancelOrder(uint256 orderId) external override nonReentrant {
        MarketOrder storage order = orders[orderId];

        require(order.user == msg.sender, "Not order owner");
        require(
            order.status == OrderStatus.Open || order.status == OrderStatus.PartiallyFilled,
            "Order not cancellable"
        );

        uint256 remaining = order.tokenAmount - order.filledAmount;
        require(remaining > 0, "Order fully filled");

        // Return escrowed assets
        if (order.orderType == OrderType.BUY) {
            // Return payment tokens
            uint256 refundAmount = (remaining * order.pricePerToken) / 1e18;
            IERC20(paymentToken).safeTransfer(msg.sender, refundAmount);
        } else {
            // Return treasury tokens
            address tokenAddress = assetFactory.getTokenAddress(order.assetId);
            IERC20(tokenAddress).safeTransfer(msg.sender, remaining);
        }

        order.status = OrderStatus.Cancelled;

        emit OrderCancelled(orderId, msg.sender);
    }

    // =============================================================
    //                     ORDER MATCHING
    // =============================================================

    /**
     * @notice Match buy and sell orders
     * @param buyOrderId Buy order ID
     * @param sellOrderId Sell order ID
     * @param amount Amount to match
     * @return tradeId Trade identifier
     */
    function matchOrders(
        uint256 buyOrderId,
        uint256 sellOrderId,
        uint256 amount
    ) external override whenNotPaused nonReentrant returns (uint256 tradeId) {
        MarketOrder storage buyOrder = orders[buyOrderId];
        MarketOrder storage sellOrder = orders[sellOrderId];

        // Validation
        require(buyOrder.orderType == OrderType.BUY, "Not buy order");
        require(sellOrder.orderType == OrderType.SELL, "Not sell order");
        require(buyOrder.assetId == sellOrder.assetId, "Asset mismatch");
        require(buyOrder.status != OrderStatus.Filled, "Buy order filled");
        require(sellOrder.status != OrderStatus.Filled, "Sell order filled");
        require(buyOrder.status != OrderStatus.Cancelled, "Buy order cancelled");
        require(sellOrder.status != OrderStatus.Cancelled, "Sell order cancelled");
        require(block.timestamp < buyOrder.expiresAt, "Buy order expired");
        require(block.timestamp < sellOrder.expiresAt, "Sell order expired");
        require(buyOrder.pricePerToken >= sellOrder.pricePerToken, "Price mismatch");

        // Calculate tradeable amount
        uint256 buyRemaining = buyOrder.tokenAmount - buyOrder.filledAmount;
        uint256 sellRemaining = sellOrder.tokenAmount - sellOrder.filledAmount;
        uint256 tradeAmount = amount;

        if (tradeAmount > buyRemaining) tradeAmount = buyRemaining;
        if (tradeAmount > sellRemaining) tradeAmount = sellRemaining;

        require(tradeAmount > 0, "No amount to trade");

        // Use sell order price for execution
        uint256 executionPrice = sellOrder.pricePerToken;
        uint256 totalValue = (tradeAmount * executionPrice) / 1e18;
        uint256 fee = calculateFee(totalValue);
        uint256 sellerReceives = totalValue - fee;

        // Transfer treasury tokens to buyer
        address tokenAddress = assetFactory.getTokenAddress(buyOrder.assetId);
        IERC20(tokenAddress).safeTransfer(buyOrder.user, tradeAmount);

        // Transfer payment to seller
        IERC20(paymentToken).safeTransfer(sellOrder.user, sellerReceives);

        // Collect fee
        collectedFees[paymentToken] += fee;

        // Update buy order
        buyOrder.filledAmount += tradeAmount;
        if (buyOrder.filledAmount == buyOrder.tokenAmount) {
            buyOrder.status = OrderStatus.Filled;

            // Refund excess payment if buy price > sell price
            uint256 buyTotalValue = (tradeAmount * buyOrder.pricePerToken) / 1e18;
            if (buyTotalValue > totalValue) {
                uint256 refund = buyTotalValue - totalValue;
                IERC20(paymentToken).safeTransfer(buyOrder.user, refund);
            }
        } else {
            buyOrder.status = OrderStatus.PartiallyFilled;
            emit OrderPartiallyFilled(buyOrderId, tradeAmount, buyOrder.tokenAmount - buyOrder.filledAmount);
        }

        // Update sell order
        sellOrder.filledAmount += tradeAmount;
        if (sellOrder.filledAmount == sellOrder.tokenAmount) {
            sellOrder.status = OrderStatus.Filled;
        } else {
            sellOrder.status = OrderStatus.PartiallyFilled;
            emit OrderPartiallyFilled(sellOrderId, tradeAmount, sellOrder.tokenAmount - sellOrder.filledAmount);
        }

        // Record trade
        tradeIdCounter++;
        tradeId = tradeIdCounter;

        Trade memory trade = Trade({
            tradeId: tradeId,
            orderId: buyOrderId, // Primary order
            assetId: buyOrder.assetId,
            buyer: buyOrder.user,
            seller: sellOrder.user,
            amount: tradeAmount,
            price: executionPrice,
            totalValue: totalValue,
            fee: fee,
            timestamp: block.timestamp
        });

        tradeHistory.push(trade);
        assetTrades[buyOrder.assetId].push(tradeId);

        totalTrades++;
        totalVolume += totalValue;

        emit OrderMatched(tradeId, buyOrderId, sellOrderId, buyOrder.user, sellOrder.user, tradeAmount, executionPrice);
    }

    // =============================================================
    //                      VIEW FUNCTIONS
    // =============================================================

    /**
     * @notice Get order details
     * @param orderId Order identifier
     * @return Market order struct
     */
    function getOrder(uint256 orderId) external view override returns (MarketOrder memory) {
        return orders[orderId];
    }

    /**
     * @notice Get open orders for asset
     * @param assetId Asset identifier
     * @param orderType Order type (BUY/SELL)
     * @return Array of order IDs
     */
    function getOpenOrders(
        uint256 assetId,
        OrderType orderType
    ) external view override returns (uint256[] memory) {
        uint256[] storage sourceOrders = orderType == OrderType.BUY
            ? assetBuyOrders[assetId]
            : assetSellOrders[assetId];

        uint256 count = 0;
        for (uint256 i = 0; i < sourceOrders.length; i++) {
            MarketOrder storage order = orders[sourceOrders[i]];
            if ((order.status == OrderStatus.Open || order.status == OrderStatus.PartiallyFilled) &&
                block.timestamp < order.expiresAt) {
                count++;
            }
        }

        uint256[] memory openOrders = new uint256[](count);
        uint256 index = 0;

        for (uint256 i = 0; i < sourceOrders.length; i++) {
            MarketOrder storage order = orders[sourceOrders[i]];
            if ((order.status == OrderStatus.Open || order.status == OrderStatus.PartiallyFilled) &&
                block.timestamp < order.expiresAt) {
                openOrders[index] = sourceOrders[i];
                index++;
            }
        }

        return openOrders;
    }

    /**
     * @notice Get user's orders
     * @param user User address
     * @return Array of order IDs
     */
    function getUserOrders(address user) external view override returns (uint256[] memory) {
        return userOrders[user];
    }

    /**
     * @notice Get trade history for asset
     * @param assetId Asset identifier
     * @param count Number of recent trades
     * @return Array of trades
     */
    function getTradeHistory(
        uint256 assetId,
        uint256 count
    ) external view override returns (Trade[] memory) {
        uint256[] storage tradeIds = assetTrades[assetId];
        uint256 returnCount = count > tradeIds.length ? tradeIds.length : count;

        Trade[] memory trades = new Trade[](returnCount);

        for (uint256 i = 0; i < returnCount; i++) {
            uint256 tradeIndex = tradeIds[tradeIds.length - returnCount + i] - 1;
            trades[i] = tradeHistory[tradeIndex];
        }

        return trades;
    }

    /**
     * @notice Calculate trading fee
     * @param totalValue Total trade value
     * @return Fee amount
     */
    function calculateFee(uint256 totalValue) public view override returns (uint256) {
        return (totalValue * feeRate) / FEE_PRECISION;
    }

    /**
     * @notice Get market statistics
     * @return stats Market statistics
     */
    function getMarketStats() external view returns (
        uint256 totalOrders_,
        uint256 totalTrades_,
        uint256 totalVolume_
    ) {
        return (totalOrders, totalTrades, totalVolume);
    }

    // =============================================================
    //                     ADMIN FUNCTIONS
    // =============================================================

    /**
     * @notice Set fee rate
     * @param newFeeRate New fee rate in basis points
     */
    function setFeeRate(uint256 newFeeRate) external onlyRole(MARKETPLACE_ADMIN_ROLE) {
        require(newFeeRate <= MAX_FEE_RATE, "Fee too high");
        feeRate = newFeeRate;
    }

    /**
     * @notice Set fee collector
     * @param newFeeCollector New fee collector address
     */
    function setFeeCollector(address newFeeCollector) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newFeeCollector != address(0), "Invalid address");
        feeCollector = newFeeCollector;
    }

    /**
     * @notice Withdraw collected fees
     * @param token Token address
     */
    function withdrawFees(address token) external onlyRole(FEE_COLLECTOR_ROLE) {
        uint256 amount = collectedFees[token];
        require(amount > 0, "No fees to withdraw");

        collectedFees[token] = 0;
        IERC20(token).safeTransfer(feeCollector, amount);
    }

    /**
     * @notice Pause marketplace
     */
    function pause() external onlyRole(MARKETPLACE_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause marketplace
     */
    function unpause() external onlyRole(MARKETPLACE_ADMIN_ROLE) {
        _unpause();
    }
}
