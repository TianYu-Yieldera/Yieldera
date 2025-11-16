// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title GMXPerformanceTracker
 * @notice Track and compare GMX V2 performance across Arbitrum and Base
 * @dev Collects metrics for AI-driven performance analysis and optimization
 *
 * Key Metrics Tracked:
 * - Execution prices and slippage
 * - Funding rates
 * - Liquidity depth
 * - Gas costs
 * - Trade execution time
 * - PnL performance
 */
contract GMXPerformanceTracker is Ownable {
    // ============ Enums ============

    enum Chain {
        ARBITRUM,
        BASE
    }

    enum TradeType {
        OPEN_LONG,
        OPEN_SHORT,
        CLOSE_LONG,
        CLOSE_SHORT
    }

    // ============ Structs ============

    struct TradeMetrics {
        Chain chain;
        TradeType tradeType;
        address market;
        uint256 sizeInUsd;
        uint256 executionPrice;      // Actual execution price
        uint256 expectedPrice;       // Expected price at order time
        uint256 slippageBps;         // Slippage in basis points
        uint256 fundingRate;         // Funding rate at trade time
        uint256 gasUsed;             // Gas consumed
        uint256 executionTimeMs;     // Time from order to execution
        uint256 timestamp;
        bytes32 orderKey;
    }

    struct PerformanceSummary {
        uint256 totalTrades;
        uint256 totalVolume;
        uint256 avgSlippageBps;
        uint256 avgFundingRate;
        uint256 avgGasUsed;
        uint256 avgExecutionTimeMs;
        uint256 successRate;         // Percentage of successful trades
        int256 totalPnl;             // Total profit/loss
    }

    struct MarketLiquidity {
        Chain chain;
        address market;
        uint256 longOpenInterest;
        uint256 shortOpenInterest;
        uint256 availableLiquidity;
        uint256 utilizationRate;     // Percentage
        uint256 timestamp;
    }

    // ============ Storage ============

    /// @notice All trade metrics
    TradeMetrics[] public allTrades;

    /// @notice Trades by chain
    mapping(Chain => TradeMetrics[]) public tradesByChain;

    /// @notice Performance summary by chain
    mapping(Chain => PerformanceSummary) public performanceSummary;

    /// @notice Liquidity snapshots
    MarketLiquidity[] public liquiditySnapshots;

    /// @notice Authorized reporters (GMX adapters)
    mapping(address => bool) public authorizedReporters;

    /// @notice Last update timestamp
    uint256 public lastUpdateTime;

    // ============ Events ============

    event TradeRecorded(
        Chain indexed chain,
        TradeType indexed tradeType,
        address indexed market,
        uint256 sizeInUsd,
        uint256 slippageBps,
        bytes32 orderKey
    );

    event LiquiditySnapshotRecorded(
        Chain indexed chain,
        address indexed market,
        uint256 availableLiquidity,
        uint256 utilizationRate
    );

    event PerformanceSummaryUpdated(
        Chain indexed chain,
        uint256 totalTrades,
        uint256 avgSlippageBps,
        int256 totalPnl
    );

    event ReporterAuthorized(address indexed reporter, bool authorized);

    // ============ Modifiers ============

    modifier onlyAuthorized() {
        require(authorizedReporters[msg.sender], "Not authorized");
        _;
    }

    // ============ Constructor ============

    constructor() Ownable(msg.sender) {
        // Grant owner as initial reporter
        authorizedReporters[msg.sender] = true;
    }

    // ============ Trade Recording ============

    /**
     * @notice Record a trade execution
     * @param chain Chain where trade executed
     * @param tradeType Type of trade
     * @param market Market address
     * @param sizeInUsd Trade size in USD
     * @param executionPrice Actual execution price
     * @param expectedPrice Expected price
     * @param fundingRate Current funding rate
     * @param gasUsed Gas consumed
     * @param executionTimeMs Execution time in ms
     * @param orderKey GMX order key
     */
    function recordTrade(
        Chain chain,
        TradeType tradeType,
        address market,
        uint256 sizeInUsd,
        uint256 executionPrice,
        uint256 expectedPrice,
        uint256 fundingRate,
        uint256 gasUsed,
        uint256 executionTimeMs,
        bytes32 orderKey
    ) external onlyAuthorized {
        // Calculate slippage
        uint256 slippageBps = _calculateSlippage(executionPrice, expectedPrice);

        TradeMetrics memory trade = TradeMetrics({
            chain: chain,
            tradeType: tradeType,
            market: market,
            sizeInUsd: sizeInUsd,
            executionPrice: executionPrice,
            expectedPrice: expectedPrice,
            slippageBps: slippageBps,
            fundingRate: fundingRate,
            gasUsed: gasUsed,
            executionTimeMs: executionTimeMs,
            timestamp: block.timestamp,
            orderKey: orderKey
        });

        allTrades.push(trade);
        tradesByChain[chain].push(trade);

        // Update performance summary
        _updatePerformanceSummary(chain, trade);

        emit TradeRecorded(chain, tradeType, market, sizeInUsd, slippageBps, orderKey);
    }

    /**
     * @notice Record liquidity snapshot
     * @param chain Chain
     * @param market Market address
     * @param longOI Long open interest
     * @param shortOI Short open interest
     * @param availableLiq Available liquidity
     */
    function recordLiquiditySnapshot(
        Chain chain,
        address market,
        uint256 longOI,
        uint256 shortOI,
        uint256 availableLiq
    ) external onlyAuthorized {
        uint256 totalOI = longOI + shortOI;
        uint256 utilization = totalOI > 0 ? (totalOI * 10000) / (totalOI + availableLiq) : 0;

        MarketLiquidity memory snapshot = MarketLiquidity({
            chain: chain,
            market: market,
            longOpenInterest: longOI,
            shortOpenInterest: shortOI,
            availableLiquidity: availableLiq,
            utilizationRate: utilization,
            timestamp: block.timestamp
        });

        liquiditySnapshots.push(snapshot);

        emit LiquiditySnapshotRecorded(chain, market, availableLiq, utilization);
    }

    // ============ Internal Functions ============

    /**
     * @notice Calculate slippage in basis points
     */
    function _calculateSlippage(
        uint256 executionPrice,
        uint256 expectedPrice
    ) private pure returns (uint256) {
        if (expectedPrice == 0) return 0;

        uint256 priceDiff = executionPrice > expectedPrice
            ? executionPrice - expectedPrice
            : expectedPrice - executionPrice;

        return (priceDiff * 10000) / expectedPrice;
    }

    /**
     * @notice Update performance summary for a chain
     */
    function _updatePerformanceSummary(Chain chain, TradeMetrics memory trade) private {
        PerformanceSummary storage summary = performanceSummary[chain];

        uint256 prevTrades = summary.totalTrades;
        summary.totalTrades++;
        summary.totalVolume += trade.sizeInUsd;

        // Update running averages
        summary.avgSlippageBps = _updateAverage(
            summary.avgSlippageBps,
            trade.slippageBps,
            prevTrades
        );

        summary.avgFundingRate = _updateAverage(
            summary.avgFundingRate,
            trade.fundingRate,
            prevTrades
        );

        summary.avgGasUsed = _updateAverage(
            summary.avgGasUsed,
            trade.gasUsed,
            prevTrades
        );

        summary.avgExecutionTimeMs = _updateAverage(
            summary.avgExecutionTimeMs,
            trade.executionTimeMs,
            prevTrades
        );

        lastUpdateTime = block.timestamp;

        emit PerformanceSummaryUpdated(
            chain,
            summary.totalTrades,
            summary.avgSlippageBps,
            summary.totalPnl
        );
    }

    /**
     * @notice Update running average
     */
    function _updateAverage(
        uint256 currentAvg,
        uint256 newValue,
        uint256 count
    ) private pure returns (uint256) {
        if (count == 0) return newValue;
        return ((currentAvg * count) + newValue) / (count + 1);
    }

    // ============ View Functions ============

    /**
     * @notice Get total trades count
     */
    function getTotalTrades() external view returns (uint256) {
        return allTrades.length;
    }

    /**
     * @notice Get trades by chain
     */
    function getTradesByChain(Chain chain) external view returns (TradeMetrics[] memory) {
        return tradesByChain[chain];
    }

    /**
     * @notice Get recent trades
     */
    function getRecentTrades(uint256 count) external view returns (TradeMetrics[] memory) {
        uint256 total = allTrades.length;
        uint256 returnCount = count > total ? total : count;

        TradeMetrics[] memory recent = new TradeMetrics[](returnCount);
        for (uint256 i = 0; i < returnCount; i++) {
            recent[i] = allTrades[total - returnCount + i];
        }

        return recent;
    }

    /**
     * @notice Get liquidity snapshots for a market
     */
    function getLiquidityHistory(
        Chain chain,
        address market,
        uint256 count
    ) external view returns (MarketLiquidity[] memory) {
        uint256 matchCount = 0;

        // Count matching snapshots
        for (uint256 i = 0; i < liquiditySnapshots.length; i++) {
            if (liquiditySnapshots[i].chain == chain && liquiditySnapshots[i].market == market) {
                matchCount++;
            }
        }

        uint256 returnCount = count > matchCount ? matchCount : count;
        MarketLiquidity[] memory history = new MarketLiquidity[](returnCount);
        uint256 index = 0;

        // Get most recent matching snapshots
        for (uint256 i = liquiditySnapshots.length; i > 0 && index < returnCount; i--) {
            MarketLiquidity memory snapshot = liquiditySnapshots[i - 1];
            if (snapshot.chain == chain && snapshot.market == market) {
                history[index] = snapshot;
                index++;
            }
        }

        return history;
    }

    /**
     * @notice Compare performance between chains
     */
    function comparePerformance() external view returns (
        PerformanceSummary memory arbitrumPerf,
        PerformanceSummary memory basePerf,
        string memory recommendation
    ) {
        arbitrumPerf = performanceSummary[Chain.ARBITRUM];
        basePerf = performanceSummary[Chain.BASE];

        // Simple recommendation logic
        if (basePerf.avgSlippageBps < arbitrumPerf.avgSlippageBps &&
            basePerf.avgGasUsed < arbitrumPerf.avgGasUsed) {
            recommendation = "Base offers better execution";
        } else if (arbitrumPerf.totalVolume > basePerf.totalVolume * 2) {
            recommendation = "Arbitrum has deeper liquidity";
        } else {
            recommendation = "Performance is comparable";
        }
    }

    // ============ Admin Functions ============

    /**
     * @notice Authorize/deauthorize a reporter
     */
    function setReporterAuthorization(address reporter, bool authorized) external onlyOwner {
        authorizedReporters[reporter] = authorized;
        emit ReporterAuthorized(reporter, authorized);
    }

    /**
     * @notice Update PnL for a chain (called by off-chain system)
     */
    function updatePnL(Chain chain, int256 pnlDelta) external onlyAuthorized {
        performanceSummary[chain].totalPnl += pnlDelta;

        emit PerformanceSummaryUpdated(
            chain,
            performanceSummary[chain].totalTrades,
            performanceSummary[chain].avgSlippageBps,
            performanceSummary[chain].totalPnl
        );
    }
}
