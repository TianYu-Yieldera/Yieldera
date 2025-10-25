// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/modules/rwa/ILiquidityAnalyzer.sol";
import "../../interfaces/modules/rwa/IOrderManager.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title LiquidityAnalyzer
 * @notice Analyzes market liquidity and health
 * @dev Provides liquidity metrics and impact calculations
 */
contract LiquidityAnalyzer is ILiquidityAnalyzer, Ownable {
    // ============ Constants ============

    uint256 private constant BASIS_POINTS = 10000;
    uint256 private constant MAX_DEPTH_SCORE = 10000;

    // ============ State Variables ============

    address public rwaModule; // Main coordinator contract
    IOrderManager public orderManager;

    // Storage for liquidity data (using Diamond Storage pattern)
    bytes32 private constant LIQUIDITY_STORAGE_POSITION = keccak256("liquidity.analyzer.storage");

    struct LiquidityStorage {
        LiquidityMetrics currentMetrics;
        mapping(uint256 => LiquidityLevel) priceLevels;
        uint256[] trackedPriceLevels;
        uint256 minLiquidityRatio; // Minimum healthy ratio (basis points)
        uint256 minDepthScore; // Minimum healthy depth score
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

        LiquidityStorage storage ls = _getStorage();
        ls.minLiquidityRatio = 8000; // 80% - reasonable balance
        ls.minDepthScore = 5000; // 50% - minimum acceptable depth
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

    function setHealthThresholds(uint256 minRatio, uint256 minDepth) external onlyOwner {
        require(minRatio <= BASIS_POINTS, "Invalid ratio");
        require(minDepth <= MAX_DEPTH_SCORE, "Invalid depth");

        LiquidityStorage storage ls = _getStorage();
        ls.minLiquidityRatio = minRatio;
        ls.minDepthScore = minDepth;
    }

    // ============ Internal Storage ============

    function _getStorage() private pure returns (LiquidityStorage storage ls) {
        bytes32 position = LIQUIDITY_STORAGE_POSITION;
        assembly {
            ls.slot := position
        }
    }

    // ============ ILiquidityAnalyzer Implementation ============

    function analyzeLiquidity() external view override returns (LiquidityMetrics memory metrics) {
        LiquidityStorage storage ls = _getStorage();
        return ls.currentMetrics;
    }

    function getLiquidityDistribution(uint256 levels)
        external
        view
        override
        returns (LiquidityLevel[] memory distribution)
    {
        LiquidityStorage storage ls = _getStorage();

        uint256 count = ls.trackedPriceLevels.length;
        if (levels > count) levels = count;

        distribution = new LiquidityLevel[](levels);
        for (uint256 i = 0; i < levels; i++) {
            uint256 priceLevel = ls.trackedPriceLevels[i];
            distribution[i] = ls.priceLevels[priceLevel];
        }

        return distribution;
    }

    function calculateMarketImpact(bool isBuy, uint256 amount)
        external
        view
        override
        returns (uint256 priceImpact, uint256 averagePrice)
    {
        uint256[] memory orders =
            orderManager.getOrdersByType(isBuy ? IOrderManager.OrderType.SELL : IOrderManager.OrderType.BUY);

        if (orders.length == 0) return (BASIS_POINTS, 0); // 100% impact if no liquidity

        uint256 remainingAmount = amount;
        uint256 totalCost = 0;
        uint256 startPrice = 0;
        bool foundStart = false;

        // Simulate filling orders
        for (uint256 i = 0; i < orders.length && remainingAmount > 0; i++) {
            IOrderManager.Order memory order = orderManager.getOrder(orders[i]);

            if (order.status != IOrderManager.OrderStatus.ACTIVE) continue;

            uint256 availableAmount = order.amount - order.filled;
            if (availableAmount == 0) continue;

            if (!foundStart) {
                startPrice = order.price;
                foundStart = true;
            }

            uint256 fillAmount = remainingAmount < availableAmount ? remainingAmount : availableAmount;
            totalCost += fillAmount * order.price;
            remainingAmount -= fillAmount;
            averagePrice = totalCost / (amount - remainingAmount);
        }

        if (!foundStart || remainingAmount > 0) {
            // Not enough liquidity
            return (BASIS_POINTS, 0);
        }

        // Calculate price impact
        if (startPrice > 0) {
            if (averagePrice > startPrice) {
                priceImpact = ((averagePrice - startPrice) * BASIS_POINTS) / startPrice;
            } else {
                priceImpact = ((startPrice - averagePrice) * BASIS_POINTS) / startPrice;
            }
        }

        return (priceImpact, averagePrice);
    }

    function getLiquidityAtPrice(uint256 price, bool isBuy)
        external
        view
        override
        returns (uint256 volume)
    {
        LiquidityStorage storage ls = _getStorage();
        LiquidityLevel storage level = ls.priceLevels[price];

        return isBuy ? level.buyVolume : level.sellVolume;
    }

    function getLiquidityInRange(uint256 minPrice, uint256 maxPrice)
        external
        view
        override
        returns (uint256 buyVolume, uint256 sellVolume)
    {
        LiquidityStorage storage ls = _getStorage();

        for (uint256 i = 0; i < ls.trackedPriceLevels.length; i++) {
            uint256 price = ls.trackedPriceLevels[i];

            if (price >= minPrice && price <= maxPrice) {
                LiquidityLevel storage level = ls.priceLevels[price];
                buyVolume += level.buyVolume;
                sellVolume += level.sellVolume;
            }
        }

        return (buyVolume, sellVolume);
    }

    function hasSufficientLiquidity(bool isBuy, uint256 amount, uint256 maxSlippage)
        external
        view
        override
        returns (bool)
    {
        uint256[] memory orders =
            orderManager.getOrdersByType(isBuy ? IOrderManager.OrderType.SELL : IOrderManager.OrderType.BUY);

        uint256 availableLiquidity = 0;
        uint256 startPrice = 0;
        bool foundStart = false;

        for (uint256 i = 0; i < orders.length; i++) {
            IOrderManager.Order memory order = orderManager.getOrder(orders[i]);

            if (order.status != IOrderManager.OrderStatus.ACTIVE) continue;

            if (!foundStart) {
                startPrice = order.price;
                foundStart = true;
            }

            // Check if this order is within acceptable slippage
            uint256 priceDeviation;
            if (order.price > startPrice) {
                priceDeviation = ((order.price - startPrice) * BASIS_POINTS) / startPrice;
            } else {
                priceDeviation = ((startPrice - order.price) * BASIS_POINTS) / startPrice;
            }

            if (priceDeviation <= maxSlippage) {
                availableLiquidity += (order.amount - order.filled);
            }

            if (availableLiquidity >= amount) {
                return true;
            }
        }

        return false;
    }

    function estimateSlippage(bool isBuy, uint256 amount)
        external
        view
        override
        returns (uint256 slippage)
    {
        uint256[] memory orders =
            orderManager.getOrdersByType(isBuy ? IOrderManager.OrderType.SELL : IOrderManager.OrderType.BUY);

        if (orders.length == 0) return BASIS_POINTS;

        uint256 remainingAmount = amount;
        uint256 startPrice = 0;
        uint256 endPrice = 0;
        bool foundStart = false;

        for (uint256 i = 0; i < orders.length && remainingAmount > 0; i++) {
            IOrderManager.Order memory order = orderManager.getOrder(orders[i]);

            if (order.status != IOrderManager.OrderStatus.ACTIVE) continue;

            uint256 availableAmount = order.amount - order.filled;
            if (availableAmount == 0) continue;

            if (!foundStart) {
                startPrice = order.price;
                foundStart = true;
            }

            uint256 fillAmount = remainingAmount < availableAmount ? remainingAmount : availableAmount;
            remainingAmount -= fillAmount;
            endPrice = order.price;
        }

        if (!foundStart || remainingAmount > 0) {
            return BASIS_POINTS; // 100% slippage if can't fill
        }

        if (startPrice == 0) return 0;

        if (endPrice > startPrice) {
            slippage = ((endPrice - startPrice) * BASIS_POINTS) / startPrice;
        } else {
            slippage = ((startPrice - endPrice) * BASIS_POINTS) / startPrice;
        }

        return slippage;
    }

    function updateLiquidityMetrics(uint256 buyLiquidity, uint256 sellLiquidity)
        external
        override
        onlyRWAModule
    {
        LiquidityStorage storage ls = _getStorage();

        ls.currentMetrics.totalBuyLiquidity = buyLiquidity;
        ls.currentMetrics.totalSellLiquidity = sellLiquidity;

        // Calculate liquidity ratio
        if (sellLiquidity > 0) {
            ls.currentMetrics.liquidityRatio = (buyLiquidity * BASIS_POINTS) / sellLiquidity;
        } else if (buyLiquidity > 0) {
            ls.currentMetrics.liquidityRatio = BASIS_POINTS * 2; // Heavily skewed
        } else {
            ls.currentMetrics.liquidityRatio = BASIS_POINTS; // No liquidity = balanced?
        }

        // Calculate depth score
        ls.currentMetrics.depthScore = this.calculateDepthScore();

        // Determine health
        (bool healthy, ) = this.checkLiquidityHealth();
        ls.currentMetrics.isHealthy = healthy;

        emit LiquidityUpdated(buyLiquidity, sellLiquidity, ls.currentMetrics.liquidityRatio);
        emit DepthScoreUpdated(ls.currentMetrics.depthScore);
    }

    function calculateDepthScore() external view override returns (uint256 score) {
        LiquidityStorage storage ls = _getStorage();

        uint256 totalLiquidity = ls.currentMetrics.totalBuyLiquidity + ls.currentMetrics.totalSellLiquidity;

        if (totalLiquidity == 0) return 0;

        // Score based on multiple factors
        uint256 volumeScore = 0;
        uint256 balanceScore = 0;
        uint256 distributionScore = 0;

        // Volume score (0-4000): more liquidity = higher score
        if (totalLiquidity > 1000000 ether) {
            volumeScore = 4000;
        } else {
            volumeScore = (totalLiquidity * 4000) / (1000000 ether);
        }

        // Balance score (0-3000): closer to 50/50 = higher score
        uint256 ratio = ls.currentMetrics.liquidityRatio;
        if (ratio > BASIS_POINTS) {
            balanceScore = (BASIS_POINTS * 3000) / ratio;
        } else {
            balanceScore = (ratio * 3000) / BASIS_POINTS;
        }

        // Distribution score (0-3000): more price levels = better distribution
        uint256 numLevels = ls.trackedPriceLevels.length;
        if (numLevels >= 20) {
            distributionScore = 3000;
        } else {
            distributionScore = (numLevels * 3000) / 20;
        }

        score = volumeScore + balanceScore + distributionScore;
        if (score > MAX_DEPTH_SCORE) score = MAX_DEPTH_SCORE;

        return score;
    }

    function checkLiquidityHealth()
        external
        view
        override
        returns (bool isHealthy, string memory reason)
    {
        LiquidityStorage storage ls = _getStorage();

        // Check liquidity ratio
        uint256 ratio = ls.currentMetrics.liquidityRatio;
        if (ratio < ls.minLiquidityRatio || ratio > (BASIS_POINTS * 2 - ls.minLiquidityRatio)) {
            emit LiquidityWarning("Imbalanced liquidity", ratio);
            return (false, "Liquidity imbalance");
        }

        // Check depth score
        uint256 depthScore = this.calculateDepthScore();
        if (depthScore < ls.minDepthScore) {
            emit LiquidityWarning("Low depth score", depthScore);
            return (false, "Insufficient depth");
        }

        // Check total liquidity
        uint256 totalLiquidity = ls.currentMetrics.totalBuyLiquidity + ls.currentMetrics.totalSellLiquidity;
        if (totalLiquidity < 100 ether) {
            emit LiquidityWarning("Low total liquidity", totalLiquidity);
            return (false, "Low liquidity");
        }

        return (true, "");
    }

    // ============ Additional Functions ============

    /**
     * @notice Update price level data
     * @param priceLevel Price level
     * @param buyVolume Buy volume at this level
     * @param sellVolume Sell volume at this level
     * @param buyOrders Number of buy orders
     * @param sellOrders Number of sell orders
     */
    function updatePriceLevel(
        uint256 priceLevel,
        uint256 buyVolume,
        uint256 sellVolume,
        uint256 buyOrders,
        uint256 sellOrders
    ) external onlyRWAModule {
        LiquidityStorage storage ls = _getStorage();

        LiquidityLevel storage level = ls.priceLevels[priceLevel];
        bool isNew = level.priceLevel == 0;

        level.priceLevel = priceLevel;
        level.buyVolume = buyVolume;
        level.sellVolume = sellVolume;
        level.buyOrders = buyOrders;
        level.sellOrders = sellOrders;

        if (isNew && (buyVolume > 0 || sellVolume > 0)) {
            ls.trackedPriceLevels.push(priceLevel);
        }
    }

    /**
     * @notice Get current health thresholds
     * @return minRatio Minimum liquidity ratio
     * @return minDepth Minimum depth score
     */
    function getHealthThresholds() external view returns (uint256 minRatio, uint256 minDepth) {
        LiquidityStorage storage ls = _getStorage();
        return (ls.minLiquidityRatio, ls.minDepthScore);
    }
}
