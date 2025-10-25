// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title ILiquidityAnalyzer
 * @notice Interface for liquidity analysis
 * @dev Provides liquidity metrics and health indicators
 */
interface ILiquidityAnalyzer {
    // ============ Structs ============

    struct LiquidityMetrics {
        uint256 totalBuyLiquidity;
        uint256 totalSellLiquidity;
        uint256 liquidityRatio; // Buy/Sell ratio (in basis points)
        uint256 averageSpread;
        uint256 depthScore; // Liquidity depth score (0-10000)
        bool isHealthy;
    }

    struct LiquidityLevel {
        uint256 priceLevel;
        uint256 buyVolume;
        uint256 sellVolume;
        uint256 buyOrders;
        uint256 sellOrders;
    }

    // ============ Events ============

    event LiquidityUpdated(uint256 buyLiquidity, uint256 sellLiquidity, uint256 ratio);
    event LiquidityWarning(string reason, uint256 metric);
    event DepthScoreUpdated(uint256 newScore);

    // ============ Functions ============

    /**
     * @notice Analyze current liquidity
     * @return metrics Current liquidity metrics
     */
    function analyzeLiquidity() external view returns (LiquidityMetrics memory metrics);

    /**
     * @notice Get liquidity distribution across price levels
     * @param levels Number of price levels to analyze
     * @return distribution Array of liquidity levels
     */
    function getLiquidityDistribution(uint256 levels)
        external
        view
        returns (LiquidityLevel[] memory distribution);

    /**
     * @notice Calculate market impact for a trade
     * @param isBuy True for buy order, false for sell
     * @param amount Trade amount
     * @return priceImpact Expected price impact (in basis points)
     * @return averagePrice Average execution price
     */
    function calculateMarketImpact(bool isBuy, uint256 amount)
        external
        view
        returns (uint256 priceImpact, uint256 averagePrice);

    /**
     * @notice Get available liquidity at price
     * @param price Price level
     * @param isBuy True for buy side, false for sell side
     * @return volume Available volume at price
     */
    function getLiquidityAtPrice(uint256 price, bool isBuy) external view returns (uint256 volume);

    /**
     * @notice Get liquidity within price range
     * @param minPrice Minimum price
     * @param maxPrice Maximum price
     * @return buyVolume Total buy volume in range
     * @return sellVolume Total sell volume in range
     */
    function getLiquidityInRange(uint256 minPrice, uint256 maxPrice)
        external
        view
        returns (uint256 buyVolume, uint256 sellVolume);

    /**
     * @notice Check if market has sufficient liquidity for trade
     * @param isBuy True for buy order
     * @param amount Trade amount
     * @param maxSlippage Maximum acceptable slippage (basis points)
     * @return hasSufficientLiquidity True if liquidity is sufficient
     */
    function hasSufficientLiquidity(bool isBuy, uint256 amount, uint256 maxSlippage)
        external
        view
        returns (bool hasSufficientLiquidity);

    /**
     * @notice Get slippage estimate for trade
     * @param isBuy True for buy order
     * @param amount Trade amount
     * @return slippage Expected slippage (basis points)
     */
    function estimateSlippage(bool isBuy, uint256 amount) external view returns (uint256 slippage);

    /**
     * @notice Update liquidity metrics after order/trade event
     * @param buyLiquidity Total buy-side liquidity
     * @param sellLiquidity Total sell-side liquidity
     */
    function updateLiquidityMetrics(uint256 buyLiquidity, uint256 sellLiquidity) external;

    /**
     * @notice Calculate depth score (0-10000)
     * @dev Higher score indicates better liquidity depth
     * @return score Liquidity depth score
     */
    function calculateDepthScore() external view returns (uint256 score);

    /**
     * @notice Check if liquidity is healthy
     * @return isHealthy True if liquidity meets health criteria
     * @return reason Reason if unhealthy
     */
    function checkLiquidityHealth() external view returns (bool isHealthy, string memory reason);
}
