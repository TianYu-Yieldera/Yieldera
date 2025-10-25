// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/core/IRWAModule.sol";
import "../../interfaces/IUpgradeable.sol";
import "../../interfaces/modules/rwa/IOrderManager.sol";
import "../../interfaces/modules/rwa/IMatchingEngine.sol";
import "../../interfaces/modules/rwa/IMarketDataProvider.sol";
import "../../interfaces/modules/rwa/IFeeCalculator.sol";
import "../../interfaces/modules/rwa/ILiquidityAnalyzer.sol";
import "../../plugins/core/BaseModule.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/**
 * @title RWAModuleV3
 * @notice Modularized RWA trading module - Phase 4 implementation
 * @dev Coordinates multiple sub-modules for RWA trading operations
 *
 * Architecture:
 * - Uses composition pattern instead of inheritance
 * - Each business logic is in a separate module
 * - This contract acts as a coordinator
 * - Maintains upgrade capability via UUPS
 */
contract RWAModuleV3 is
    Initializable,
    IRWAModule,
    IUpgradeable,
    BaseModule,
    OwnableUpgradeable,
    UUPSUpgradeable
{
    using SafeERC20 for IERC20;

    // ============ Constants ============

    bytes32 public constant MODULE_ID = keccak256("RWA_MODULE");
    string public constant MODULE_NAME = "RWAModule";
    string public constant MODULE_VERSION = "3.0.0";

    // ============ Sub-Modules ============

    IOrderManager public orderManager;
    IMatchingEngine public matchingEngine;
    IMarketDataProvider public marketDataProvider;
    IFeeCalculator public feeCalculator;
    ILiquidityAnalyzer public liquidityAnalyzer;

    // ============ Configuration ============

    TradingPair private _tradingPair;
    address public feeCollector;
    address public legacyOrderBook; // For backward compatibility

    bool private _isPaused;

    // Upgrade tracking
    struct UpgradeRecord {
        address implementation;
        uint256 timestamp;
        string version;
    }

    UpgradeRecord[] private _upgradeHistory;
    mapping(address => bool) private _authorizedUpgraders;
    bool private _upgradesPaused;

    // ============ Events ============

    event SubModuleUpdated(string indexed moduleName, address oldAddress, address newAddress);

    // ============ Constructor & Initializer ============

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initialize the modular RWA module
     */
    function initialize(
        address _orderManager,
        address _matchingEngine,
        address _marketDataProvider,
        address _feeCalculator,
        address _liquidityAnalyzer,
        address baseToken,
        address quoteToken,
        address _legacyOrderBook
    ) public initializer {
        __Ownable_init();
        __UUPSUpgradeable_init();

        require(_orderManager != address(0), "Invalid order manager");
        require(_matchingEngine != address(0), "Invalid matching engine");
        require(_marketDataProvider != address(0), "Invalid market data provider");
        require(_feeCalculator != address(0), "Invalid fee calculator");
        require(_liquidityAnalyzer != address(0), "Invalid liquidity analyzer");
        require(baseToken != address(0), "Invalid base token");
        require(quoteToken != address(0), "Invalid quote token");

        orderManager = IOrderManager(_orderManager);
        matchingEngine = IMatchingEngine(_matchingEngine);
        marketDataProvider = IMarketDataProvider(_marketDataProvider);
        feeCalculator = IFeeCalculator(_feeCalculator);
        liquidityAnalyzer = ILiquidityAnalyzer(_liquidityAnalyzer);
        legacyOrderBook = _legacyOrderBook;

        // Initialize trading pair
        _tradingPair = TradingPair({
            baseToken: baseToken,
            quoteToken: quoteToken,
            minOrderSize: 1 ether,
            maxOrderSize: 1000000 ether,
            minPrice: 0,
            maxPrice: type(uint256).max,
            makerFee: 25, // 0.25%
            takerFee: 50, // 0.50%
            isActive: true
        });

        feeCollector = msg.sender;

        // Record deployment
        _upgradeHistory.push(
            UpgradeRecord({
                implementation: address(this),
                timestamp: block.timestamp,
                version: MODULE_VERSION
            })
        );

        _authorizedUpgraders[msg.sender] = true;
    }

    // ============ Sub-Module Management ============

    function setOrderManager(address newManager) external onlyOwner {
        require(newManager != address(0), "Invalid address");
        address oldManager = address(orderManager);
        orderManager = IOrderManager(newManager);
        emit SubModuleUpdated("OrderManager", oldManager, newManager);
    }

    function setMatchingEngine(address newEngine) external onlyOwner {
        require(newEngine != address(0), "Invalid address");
        address oldEngine = address(matchingEngine);
        matchingEngine = IMatchingEngine(newEngine);
        emit SubModuleUpdated("MatchingEngine", oldEngine, newEngine);
    }

    function setMarketDataProvider(address newProvider) external onlyOwner {
        require(newProvider != address(0), "Invalid address");
        address oldProvider = address(marketDataProvider);
        marketDataProvider = IMarketDataProvider(newProvider);
        emit SubModuleUpdated("MarketDataProvider", oldProvider, newProvider);
    }

    function setFeeCalculator(address newCalculator) external onlyOwner {
        require(newCalculator != address(0), "Invalid address");
        address oldCalculator = address(feeCalculator);
        feeCalculator = IFeeCalculator(newCalculator);
        emit SubModuleUpdated("FeeCalculator", oldCalculator, newCalculator);
    }

    function setLiquidityAnalyzer(address newAnalyzer) external onlyOwner {
        require(newAnalyzer != address(0), "Invalid address");
        address oldAnalyzer = address(liquidityAnalyzer);
        liquidityAnalyzer = ILiquidityAnalyzer(newAnalyzer);
        emit SubModuleUpdated("LiquidityAnalyzer", oldAnalyzer, newAnalyzer);
    }

    // ============ BaseModule Overrides ============

    function getModuleId() external pure override returns (bytes32) {
        return MODULE_ID;
    }

    function getVersion() external pure override returns (string memory) {
        return MODULE_VERSION;
    }

    function getDependencies() external pure override returns (bytes32[] memory) {
        bytes32[] memory deps = new bytes32[](2);
        deps[0] = keccak256("PRICE_ORACLE_MODULE");
        deps[1] = keccak256("AUDIT_MODULE");
        return deps;
    }

    function isActive() external view override returns (bool) {
        return !_isPaused && _tradingPair.isActive;
    }

    function initialize(bytes calldata) external override {
        revert("Use initialize with all parameters");
    }

    function healthCheck()
        external
        view
        override(IModule, BaseModule)
        returns (bool healthy, string memory message)
    {
        // Check if sub-modules are set
        if (address(orderManager) == address(0)) {
            return (false, "OrderManager not set");
        }
        if (address(matchingEngine) == address(0)) {
            return (false, "MatchingEngine not set");
        }
        if (address(marketDataProvider) == address(0)) {
            return (false, "MarketDataProvider not set");
        }
        if (address(feeCalculator) == address(0)) {
            return (false, "FeeCalculator not set");
        }
        if (address(liquidityAnalyzer) == address(0)) {
            return (false, "LiquidityAnalyzer not set");
        }

        // Check liquidity health
        (bool liquidityHealthy, string memory reason) = liquidityAnalyzer.checkLiquidityHealth();
        if (!liquidityHealthy) {
            return (false, reason);
        }

        return (true, "RWA module healthy");
    }

    // ============ Order Management ============

    function placeOrder(OrderType orderType, uint256 price, uint256 amount)
        external
        override
        returns (uint256 orderId)
    {
        require(!_isPaused, "Trading paused");
        require(_tradingPair.isActive, "Trading pair inactive");
        require(amount >= _tradingPair.minOrderSize, "Amount too small");
        require(amount <= _tradingPair.maxOrderSize, "Amount too large");
        require(price >= _tradingPair.minPrice, "Price too low");
        require(price <= _tradingPair.maxPrice, "Price too high");

        // Lock funds
        if (orderType == OrderType.BUY) {
            // Lock quote token (payment)
            uint256 totalCost = (price * amount) / 1 ether;
            IERC20(_tradingPair.quoteToken).safeTransferFrom(msg.sender, address(this), totalCost);
        } else {
            // Lock base token (RWA)
            IERC20(_tradingPair.baseToken).safeTransferFrom(msg.sender, address(this), amount);
        }

        // Create order through OrderManager
        IOrderManager.OrderType orderManagerType = orderType == OrderType.BUY
            ? IOrderManager.OrderType.BUY
            : IOrderManager.OrderType.SELL;

        orderId = orderManager.createOrder(msg.sender, orderManagerType, amount, price);

        emit OrderPlaced(orderId, msg.sender, orderType, price, amount);

        // Try to match order
        this.matchOrders();

        return orderId;
    }

    function cancelOrder(uint256 orderId) external override {
        IOrderManager.Order memory order = orderManager.getOrder(orderId);
        require(order.trader == msg.sender, "Not order owner");

        uint256 remainingAmount = order.amount - order.filled;
        require(remainingAmount > 0, "Order fully filled");

        // Cancel through OrderManager
        orderManager.cancelOrder(orderId, msg.sender);

        // Refund locked funds
        if (order.orderType == IOrderManager.OrderType.BUY) {
            uint256 refundAmount = (order.price * remainingAmount) / 1 ether;
            IERC20(_tradingPair.quoteToken).safeTransfer(msg.sender, refundAmount);
        } else {
            IERC20(_tradingPair.baseToken).safeTransfer(msg.sender, remainingAmount);
        }

        emit OrderCancelled(orderId, msg.sender, remainingAmount);
    }

    function getOrder(uint256 orderId) external view override returns (Order memory) {
        IOrderManager.Order memory managerOrder = orderManager.getOrder(orderId);

        return Order({
            orderId: managerOrder.orderId,
            trader: managerOrder.trader,
            orderType: managerOrder.orderType == IOrderManager.OrderType.BUY ? OrderType.BUY : OrderType.SELL,
            price: managerOrder.price,
            amount: managerOrder.amount,
            filled: managerOrder.filled,
            timestamp: managerOrder.timestamp,
            status: _convertOrderStatus(managerOrder.status)
        });
    }

    function getUserOpenOrders(address user) external view override returns (Order[] memory) {
        uint256[] memory orderIds = orderManager.getUserOrders(user);

        // Count active orders
        uint256 activeCount = 0;
        for (uint256 i = 0; i < orderIds.length; i++) {
            IOrderManager.Order memory order = orderManager.getOrder(orderIds[i]);
            if (order.status == IOrderManager.OrderStatus.ACTIVE) {
                activeCount++;
            }
        }

        // Populate result
        Order[] memory orders = new Order[](activeCount);
        uint256 index = 0;
        for (uint256 i = 0; i < orderIds.length; i++) {
            IOrderManager.Order memory order = orderManager.getOrder(orderIds[i]);
            if (order.status == IOrderManager.OrderStatus.ACTIVE) {
                orders[index++] = Order({
                    orderId: order.orderId,
                    trader: order.trader,
                    orderType: order.orderType == IOrderManager.OrderType.BUY ? OrderType.BUY : OrderType.SELL,
                    price: order.price,
                    amount: order.amount,
                    filled: order.filled,
                    timestamp: order.timestamp,
                    status: _convertOrderStatus(order.status)
                });
            }
        }

        return orders;
    }

    function getUserOrderHistory(address user, uint256 offset, uint256 limit)
        external
        view
        override
        returns (Order[] memory)
    {
        uint256[] memory orderIds = orderManager.getUserOrders(user);

        if (offset >= orderIds.length) return new Order[](0);

        uint256 end = offset + limit;
        if (end > orderIds.length) end = orderIds.length;

        Order[] memory orders = new Order[](end - offset);
        for (uint256 i = offset; i < end; i++) {
            IOrderManager.Order memory order = orderManager.getOrder(orderIds[i]);
            orders[i - offset] = Order({
                orderId: order.orderId,
                trader: order.trader,
                orderType: order.orderType == IOrderManager.OrderType.BUY ? OrderType.BUY : OrderType.SELL,
                price: order.price,
                amount: order.amount,
                filled: order.filled,
                timestamp: order.timestamp,
                status: _convertOrderStatus(order.status)
            });
        }

        return orders;
    }

    // ============ Trading (Simplified implementations) ============

    function placeMarketOrder(OrderType, uint256)
        external
        pure
        override
        returns (uint256, uint256)
    {
        revert("Market orders not implemented in V3");
    }

    function matchOrders() external override {
        // Get active buy and sell orders
        uint256[] memory buyOrders = orderManager.getOrdersByType(IOrderManager.OrderType.BUY);
        uint256[] memory sellOrders = orderManager.getOrdersByType(IOrderManager.OrderType.SELL);

        // Match orders (simplified - in production would be more sophisticated)
        for (uint256 i = 0; i < buyOrders.length && i < 10; i++) {
            // Limit iterations for gas
            IOrderManager.Order memory buyOrder = orderManager.getOrder(buyOrders[i]);
            if (buyOrder.status != IOrderManager.OrderStatus.ACTIVE) continue;

            for (uint256 j = 0; j < sellOrders.length && j < 10; j++) {
                IOrderManager.Order memory sellOrder = orderManager.getOrder(sellOrders[j]);
                if (sellOrder.status != IOrderManager.OrderStatus.ACTIVE) continue;

                // Check if orders can match
                (bool canMatch, uint256 matchAmount) =
                    matchingEngine.canMatch(buyOrders[i], sellOrders[j]);

                if (canMatch && matchAmount > 0) {
                    // Execute trade
                    _executeTrade(buyOrders[i], sellOrders[j], matchAmount, sellOrder.price);
                }
            }
        }
    }

    // ============ Order Book (Simplified) ============

    function getOrderBookDepth(OrderType, uint256)
        external
        pure
        override
        returns (uint256[] memory, uint256[] memory)
    {
        // Simplified - return empty arrays
        return (new uint256[](0), new uint256[](0));
    }

    function getBestBid() external view override returns (uint256 price, uint256 amount) {
        price = matchingEngine.getBestPrice(IOrderManager.OrderType.BUY);
        amount = 0; // Simplified
        return (price, amount);
    }

    function getBestAsk() external view override returns (uint256 price, uint256 amount) {
        price = matchingEngine.getBestPrice(IOrderManager.OrderType.SELL);
        amount = 0; // Simplified
        return (price, amount);
    }

    function getSpread() external view override returns (uint256 spread) {
        return marketDataProvider.getSpread();
    }

    function getMidPrice() external view override returns (uint256) {
        (uint256 bestBid,) = this.getBestBid();
        (uint256 bestAsk,) = this.getBestAsk();
        if (bestBid == 0 || bestAsk == 0) return 0;
        return (bestBid + bestAsk) / 2;
    }

    // ============ Trade History ============

    function getTradeHistory(uint256, uint256)
        external
        view
        override
        returns (Trade[] memory)
    {
        // Simplified
        return new Trade[](0);
    }

    function getUserTrades(address user, uint256, uint256)
        external
        view
        override
        returns (Trade[] memory)
    {
        uint256[] memory tradeIds = matchingEngine.getUserTrades(user);

        Trade[] memory trades = new Trade[](tradeIds.length);
        for (uint256 i = 0; i < tradeIds.length; i++) {
            IMatchingEngine.Trade memory trade = matchingEngine.getTrade(tradeIds[i]);
            trades[i] = Trade({
                tradeId: trade.tradeId,
                buyOrderId: trade.buyOrderId,
                sellOrderId: trade.sellOrderId,
                buyer: trade.buyer,
                seller: trade.seller,
                price: trade.price,
                amount: trade.amount,
                timestamp: trade.timestamp
            });
        }

        return trades;
    }

    // ============ Market Statistics ============

    function getMarketStats() external view override returns (MarketStats memory) {
        IMarketDataProvider.MarketStats memory stats = marketDataProvider.getMarketStats();

        return MarketStats({
            lastPrice: stats.lastPrice,
            highPrice24h: stats.highPrice24h,
            lowPrice24h: stats.lowPrice24h,
            volume24h: stats.volumeLast24h,
            totalVolume: stats.totalVolume,
            openBuyOrders: 0, // Simplified
            openSellOrders: 0, // Simplified
            lastTradeTime: block.timestamp
        });
    }

    function getTradingPair() external view override returns (TradingPair memory) {
        return _tradingPair;
    }

    function get24hVolume() external view override returns (uint256) {
        return marketDataProvider.get24hVolume();
    }

    function getTotalVolume() external view override returns (uint256) {
        IMarketDataProvider.MarketStats memory stats = marketDataProvider.getMarketStats();
        return stats.totalVolume;
    }

    function getLastPrice() external view override returns (uint256) {
        IMarketDataProvider.MarketStats memory stats = marketDataProvider.getMarketStats();
        return stats.lastPrice;
    }

    // ============ Configuration ============

    function updateFees(uint256 makerFee, uint256 takerFee) external override onlyOwner {
        _tradingPair.makerFee = makerFee;
        _tradingPair.takerFee = takerFee;
        emit FeesUpdated(makerFee, takerFee);
    }

    function updateOrderLimits(uint256 minSize, uint256 maxSize) external override onlyOwner {
        _tradingPair.minOrderSize = minSize;
        _tradingPair.maxOrderSize = maxSize;
        emit OrderLimitsUpdated(minSize, maxSize);
    }

    function updatePriceLimits(uint256 minPrice, uint256 maxPrice) external override onlyOwner {
        _tradingPair.minPrice = minPrice;
        _tradingPair.maxPrice = maxPrice;
    }

    function setFeeCollector(address _feeCollector) external override onlyOwner {
        require(_feeCollector != address(0), "Invalid address");
        feeCollector = _feeCollector;
    }

    // ============ Liquidity ============

    function getTotalLiquidity()
        external
        view
        override
        returns (uint256 buyLiquidity, uint256 sellLiquidity)
    {
        ILiquidityAnalyzer.LiquidityMetrics memory metrics = liquidityAnalyzer.analyzeLiquidity();
        return (metrics.totalBuyLiquidity, metrics.totalSellLiquidity);
    }

    function checkLiquidity(OrderType orderType, uint256 amount, uint256 maxSlippage)
        external
        view
        override
        returns (bool sufficient, uint256 estimatedPrice)
    {
        bool isBuy = (orderType == OrderType.BUY);
        sufficient = liquidityAnalyzer.hasSufficientLiquidity(isBuy, amount, maxSlippage);

        (uint256 slippage, uint256 avgPrice) = liquidityAnalyzer.calculateMarketImpact(isBuy, amount);

        estimatedPrice = slippage <= maxSlippage ? avgPrice : 0;

        return (sufficient, estimatedPrice);
    }

    // ============ Token Information ============

    function getBaseToken() external view override returns (address) {
        return _tradingPair.baseToken;
    }

    function getQuoteToken() external view override returns (address) {
        return _tradingPair.quoteToken;
    }

    // ============ Pause ============

    function pause() external override(IModule, BaseModule) onlyOwner {
        _isPaused = true;
    }

    function unpause() external override(IModule, BaseModule) onlyOwner {
        _isPaused = false;
    }

    // ============ Internal Helpers ============

    function _executeTrade(uint256 buyOrderId, uint256 sellOrderId, uint256 amount, uint256 price)
        internal
    {
        // Execute trade through MatchingEngine
        uint256 tradeId = matchingEngine.executeTrade(buyOrderId, sellOrderId, amount, price);

        // Fill orders
        orderManager.fillOrder(buyOrderId, amount);
        orderManager.fillOrder(sellOrderId, amount);

        // Transfer tokens
        IOrderManager.Order memory buyOrder = orderManager.getOrder(buyOrderId);
        IOrderManager.Order memory sellOrder = orderManager.getOrder(sellOrderId);

        IERC20(_tradingPair.baseToken).safeTransfer(buyOrder.trader, amount);
        uint256 payment = (price * amount) / 1 ether;
        IERC20(_tradingPair.quoteToken).safeTransfer(sellOrder.trader, payment);

        // Collect fees (simplified - should use FeeCalculator)
        // Update market data
        marketDataProvider.updateStats(price, amount);

        emit TradeExecuted(tradeId, buyOrder.trader, sellOrder.trader, price, amount);
        emit OrderMatched(buyOrderId, sellOrderId, price, amount);
    }

    function _convertOrderStatus(IOrderManager.OrderStatus status)
        internal
        pure
        returns (OrderStatus)
    {
        if (status == IOrderManager.OrderStatus.ACTIVE) return OrderStatus.OPEN;
        if (status == IOrderManager.OrderStatus.FILLED) return OrderStatus.FILLED;
        if (status == IOrderManager.OrderStatus.CANCELLED) return OrderStatus.CANCELLED;
        return OrderStatus.OPEN;
    }

    // ============ IUpgradeable Implementation ============

    function upgradeTo(address newImplementation) external override onlyOwner {
        _authorizeUpgrade(newImplementation);
        _upgradeToAndCallUUPS(newImplementation, new bytes(0), false);
        _upgradeHistory.push(
            UpgradeRecord({
                implementation: newImplementation,
                timestamp: block.timestamp,
                version: MODULE_VERSION
            })
        );
        emit Upgraded(address(this), newImplementation, MODULE_VERSION);
    }

    function upgradeToAndCall(address newImplementation, bytes memory data)
        external
        payable
        override
        onlyOwner
    {
        _authorizeUpgrade(newImplementation);
        _upgradeToAndCallUUPS(newImplementation, data, true);
        _upgradeHistory.push(
            UpgradeRecord({
                implementation: newImplementation,
                timestamp: block.timestamp,
                version: MODULE_VERSION
            })
        );
        emit Upgraded(address(this), newImplementation, MODULE_VERSION);
    }

    function getImplementation() external view override returns (address) {
        return _getImplementation();
    }

    function getImplementationVersion() external pure override returns (string memory) {
        return MODULE_VERSION;
    }

    function canUpgrade(address account) external view override returns (bool) {
        return _authorizedUpgraders[account] || account == owner();
    }

    function authorizeUpgrader(address account) external override onlyOwner {
        _authorizedUpgraders[account] = true;
        emit UpgradeAuthorized(account, msg.sender);
    }

    function revokeUpgradeAuthorization(address account) external override onlyOwner {
        _authorizedUpgraders[account] = false;
        emit UpgradeAuthorizationRevoked(account);
    }

    function validateUpgrade(address newImplementation)
        external
        view
        override
        returns (bool valid, string memory reason)
    {
        if (newImplementation == address(0)) return (false, "Zero address");
        if (newImplementation.code.length == 0) return (false, "Not a contract");
        if (_upgradesPaused) return (false, "Upgrades paused");
        return (true, "");
    }

    function beforeUpgrade(address) external override onlyOwner {}
    function afterUpgrade(address) external override onlyOwner {}

    function getUpgradeHistory()
        external
        view
        override
        returns (address[] memory, uint256[] memory, string[] memory)
    {
        uint256 length = _upgradeHistory.length;
        address[] memory implementations = new address[](length);
        uint256[] memory timestamps = new uint256[](length);
        string[] memory versions = new string[](length);

        for (uint256 i = 0; i < length; i++) {
            implementations[i] = _upgradeHistory[i].implementation;
            timestamps[i] = _upgradeHistory[i].timestamp;
            versions[i] = _upgradeHistory[i].version;
        }

        return (implementations, timestamps, versions);
    }

    function pauseUpgrades() external override onlyOwner {
        _upgradesPaused = true;
    }

    function resumeUpgrades() external override onlyOwner {
        _upgradesPaused = false;
    }

    function upgradesPaused() external view override returns (bool) {
        return _upgradesPaused;
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {
        require(!_upgradesPaused, "Upgrades paused");
        require(newImplementation != address(0), "Zero address");
        require(
            _authorizedUpgraders[msg.sender] || msg.sender == owner(), "Not authorized"
        );
    }

    function getModuleInfo() external view override returns (ModuleInfo memory) {
        return ModuleInfo({
            moduleId: MODULE_ID,
            name: MODULE_NAME,
            version: MODULE_VERSION,
            implementation: address(this),
            state: _isPaused ? ModuleState.PAUSED : ModuleState.ACTIVE,
            installedAt: _upgradeHistory.length > 0
                ? _upgradeHistory[0].timestamp
                : block.timestamp,
            lastUpdated: _upgradeHistory.length > 0
                ? _upgradeHistory[_upgradeHistory.length - 1].timestamp
                : block.timestamp
        });
    }
}
