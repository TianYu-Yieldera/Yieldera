// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/**
 * @title LiquidityMiningStrategy
 * @notice Example strategy for Uniswap V3 liquidity mining
 * @dev This is a demonstration contract showing how to use UniswapV3Adapter
 *
 * WARNING: This is for educational purposes only!
 * Real liquidity mining strategies require:
 * - Dynamic range adjustment based on market conditions
 * - Impermanent loss calculation
 * - Gas-efficient rebalancing
 * - MEV protection
 * - Slippage management
 *
 * Liquidity Mining Concept:
 * - Provide liquidity to Uniswap V3 pool
 * - Earn trading fees (0.01%-1% per swap)
 * - Concentrated liquidity = higher capital efficiency
 * - Need to manage price ranges actively
 */
contract LiquidityMiningStrategy {
    using SafeERC20 for IERC20;

    address public owner;
    address public uniswapAdapter;

    // Strategy parameters
    struct StrategyParams {
        address token0;         // First token (e.g., USDC)
        address token1;         // Second token (e.g., ETH)
        uint24 fee;            // Fee tier (500 = 0.05%)
        int24 tickLower;       // Lower tick of range
        int24 tickUpper;       // Upper tick of range
        uint256 rebalanceThreshold; // When to rebalance (basis points from range)
    }

    // Position tracking
    struct Position {
        uint256 tokenId;       // Uniswap V3 NFT ID
        uint128 liquidity;     // Current liquidity
        uint256 feesCollected0; // Total fees collected (token0)
        uint256 feesCollected1; // Total fees collected (token1)
        uint256 lastRebalance; // Last rebalance timestamp
    }

    mapping(uint256 => Position) public positions;
    uint256 public positionCount;

    // Events
    event StrategyDeployed(
        uint256 indexed positionId,
        address token0,
        address token1,
        uint24 fee,
        int24 tickLower,
        int24 tickUpper
    );

    event FeesCompounded(
        uint256 indexed positionId,
        uint256 amount0,
        uint256 amount1,
        uint128 liquidityAdded
    );

    event PositionRebalanced(
        uint256 indexed positionId,
        int24 newTickLower,
        int24 newTickUpper
    );

    constructor(address _uniswapAdapter) {
        owner = msg.sender;
        uniswapAdapter = _uniswapAdapter;
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }

    /**
     * @notice Deploy liquidity mining strategy
     * @param params Strategy parameters
     * @param amount0 Initial amount of token0
     * @param amount1 Initial amount of token1
     * @return positionId Strategy position ID
     */
    function deployStrategy(
        StrategyParams memory params,
        uint256 amount0,
        uint256 amount1
    ) external onlyOwner returns (uint256 positionId) {
        require(amount0 > 0 && amount1 > 0, "Invalid amounts");

        // Transfer tokens from owner
        IERC20(params.token0).safeTransferFrom(msg.sender, address(this), amount0);
        IERC20(params.token1).safeTransferFrom(msg.sender, address(this), amount1);

        // Approve Uniswap Adapter
        IERC20(params.token0).forceApprove(uniswapAdapter, amount0);
        IERC20(params.token1).forceApprove(uniswapAdapter, amount1);

        // Add liquidity via adapter
        // In production, you would call:
        // (uint256 tokenId, uint128 liquidity,,) = IUniswapV3Adapter(uniswapAdapter).addLiquidity(...)

        // For this example, we simulate the position
        positionId = positionCount++;
        Position storage pos = positions[positionId];
        pos.tokenId = positionId; // In reality, this would be the NFT ID
        pos.liquidity = uint128(amount0 + amount1); // Simplified
        pos.lastRebalance = block.timestamp;

        emit StrategyDeployed(
            positionId,
            params.token0,
            params.token1,
            params.fee,
            params.tickLower,
            params.tickUpper
        );
    }

    /**
     * @notice Collect and compound fees
     * @param positionId Strategy position ID
     * @dev This is a simplified version for demonstration
     */
    function compoundFees(uint256 positionId) external {
        Position storage pos = positions[positionId];
        require(pos.liquidity > 0, "Invalid position");

        // In production, you would:
        // 1. Call collectFees() from adapter
        // 2. Swap half of fees to maintain ratio
        // 3. Add liquidity back to position

        // Example pseudocode:
        // (uint256 fees0, uint256 fees1) = adapter.collectFees(pos.tokenId);
        //
        // if (fees0 > fees1) {
        //     // Swap half of fees0 to fees1
        //     uint256 swapAmount = (fees0 - fees1) / 2;
        //     fees1 += adapter.swap(token0, token1, swapAmount);
        //     fees0 -= swapAmount;
        // }
        //
        // (uint128 liquidityAdded,,) = adapter.increaseLiquidity(
        //     pos.tokenId, fees0, fees1, 0, 0, deadline
        // );

        emit FeesCompounded(positionId, 0, 0, 0);
    }

    /**
     * @notice Rebalance position to new price range
     * @param positionId Strategy position ID
     * @param newTickLower New lower tick
     * @param newTickUpper New upper tick
     * @dev This is a simplified version for demonstration
     */
    function rebalancePosition(
        uint256 positionId,
        int24 newTickLower,
        int24 newTickUpper
    ) external onlyOwner {
        Position storage pos = positions[positionId];
        require(pos.liquidity > 0, "Invalid position");

        // In production, you would:
        // 1. Remove all liquidity from old position
        // 2. Collect remaining fees
        // 3. Create new position with new tick range
        // 4. Migrate all liquidity to new position

        // Example pseudocode:
        // (uint256 amount0, uint256 amount1) = adapter.removeLiquidity(
        //     pos.tokenId, pos.liquidity, 0, 0, deadline
        // );
        //
        // (uint256 collected0, uint256 collected1) = adapter.collectFees(pos.tokenId);
        //
        // (uint256 newTokenId, uint128 newLiquidity,,) = adapter.addLiquidity(
        //     token0, token1, fee,
        //     newTickLower, newTickUpper,
        //     amount0 + collected0, amount1 + collected1,
        //     0, 0, deadline
        // );

        pos.lastRebalance = block.timestamp;

        emit PositionRebalanced(positionId, newTickLower, newTickUpper);
    }

    /**
     * @notice Calculate current position value
     * @param positionId Strategy position ID
     * @return value0 Value in token0
     * @return value1 Value in token1
     * @dev Simplified calculation for demonstration
     */
    function getPositionValue(uint256 positionId)
        external
        view
        returns (uint256 value0, uint256 value1)
    {
        Position storage pos = positions[positionId];

        // In production, you would:
        // 1. Get current position from adapter
        // 2. Calculate token amounts based on current price
        // 3. Add uncollected fees

        return (uint256(pos.liquidity) / 2, uint256(pos.liquidity) / 2);
    }

    /**
     * @notice Check if position needs rebalancing
     * @param positionId Strategy position ID
     * @return needsRebalance True if outside optimal range
     */
    function checkRebalanceNeeded(uint256 positionId) external view returns (bool needsRebalance) {
        Position storage pos = positions[positionId];

        // In production, you would:
        // 1. Get current pool price
        // 2. Compare with position range
        // 3. Calculate distance from center
        // 4. Return true if beyond threshold

        // Simplified: rebalance every 7 days
        return block.timestamp - pos.lastRebalance > 7 days;
    }

    /**
     * @notice Estimate APY for position
     * @param positionId Strategy position ID
     * @return apy Estimated APY in basis points
     */
    function estimateAPY(uint256 positionId) external view returns (uint256 apy) {
        Position storage pos = positions[positionId];

        // In production, you would:
        // 1. Calculate total fees collected
        // 2. Calculate time elapsed
        // 3. Annualize the return
        // 4. Factor in IL (impermanent loss)

        // Example calculation:
        // uint256 totalFees = pos.feesCollected0 + pos.feesCollected1;
        // uint256 totalValue = currentValue0 + currentValue1;
        // uint256 timeElapsed = block.timestamp - deploymentTime;
        // apy = (totalFees * 365 days * 10000) / (totalValue * timeElapsed);

        return 500; // 5% APY (simplified)
    }

    /**
     * @notice Emergency withdrawal
     * @param positionId Strategy position ID
     */
    function emergencyWithdraw(uint256 positionId) external onlyOwner {
        Position storage pos = positions[positionId];
        require(pos.liquidity > 0, "Invalid position");

        // In production, you would:
        // 1. Remove all liquidity
        // 2. Collect all fees
        // 3. Transfer tokens back to owner
        // 4. Burn the NFT position

        delete positions[positionId];
    }

    /**
     * @notice Withdraw profits
     * @param token Token address
     * @param amount Amount to withdraw
     */
    function withdrawProfits(address token, uint256 amount) external onlyOwner {
        IERC20(token).safeTransfer(owner, amount);
    }

    /**
     * @notice Get strategy statistics
     * @return totalPositions Number of active positions
     * @return totalLiquidity Sum of all liquidity
     * @return totalFees Total fees collected
     */
    function getStrategyStats()
        external
        view
        returns (
            uint256 totalPositions,
            uint256 totalLiquidity,
            uint256 totalFees
        )
    {
        totalPositions = positionCount;

        for (uint256 i = 0; i < positionCount; i++) {
            Position storage pos = positions[i];
            totalLiquidity += pos.liquidity;
            totalFees += pos.feesCollected0 + pos.feesCollected1;
        }
    }
}
