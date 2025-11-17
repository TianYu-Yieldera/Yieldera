// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "../interfaces/external/IAerodromeRouter.sol";
import "../interfaces/external/IAerodromePair.sol";
import "../interfaces/external/IAerodromeGauge.sol";
import "../aggregator/L2StateAggregator.sol";

/**
 * @title AerodromeAdapter
 * @notice Adapter for integrating with Aerodrome Finance on Base
 * @dev Provides liquidity provision, swapping, and staking functionality
 *
 * Key Features:
 * - Add/remove liquidity to Aerodrome pools
 * - Swap tokens through Aerodrome Router
 * - Stake LP tokens in gauges to earn AERO rewards
 * - Support for both stable and volatile pools
 * - Integration with L2StateAggregator for position tracking
 *
 * Aerodrome is Base's native DEX based on Solidly/Velodrome design:
 * - Volatile pools: Traditional x*y=k AMM for uncorrelated assets
 * - Stable pools: Optimized for correlated assets (stablecoins)
 * - veAERO: Vote-escrowed model for governance and fee distribution
 */
contract AerodromeAdapter is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ============ State Variables ============

    /// @notice Aerodrome Router
    IAerodromeRouter public immutable router;

    /// @notice L2 State Aggregator for tracking system state
    L2StateAggregator public stateAggregator;

    /// @notice AERO token address
    address public immutable aeroToken;

    /// @notice WETH address on Base
    address public constant WETH = 0x4200000000000000000000000000000000000006;

    // Position tracking
    struct LiquidityPosition {
        address pair;
        uint256 lpTokenAmount;
        uint256 token0Amount;
        uint256 token1Amount;
        uint256 stakedAmount;
        address gauge;
        uint256 lastUpdate;
    }

    struct SwapRecord {
        address tokenIn;
        address tokenOut;
        uint256 amountIn;
        uint256 amountOut;
        uint256 timestamp;
    }

    mapping(address => mapping(address => LiquidityPosition)) public userPositions; // user => pair => position
    mapping(address => uint256) public userPositionCount;
    mapping(address => SwapRecord[]) public userSwapHistory;

    // Statistics
    uint256 public totalLiquidityProvided;
    uint256 public totalSwapVolume;
    uint256 public totalAeroEarned;
    uint256 public activeUsers;

    // Slippage protection (in basis points, 10000 = 100%)
    uint256 public maxSlippage = 100; // 1% default
    uint256 public constant SLIPPAGE_PRECISION = 10000;

    // ============ Events ============

    event LiquidityAdded(
        address indexed user,
        address indexed pair,
        uint256 amount0,
        uint256 amount1,
        uint256 liquidity,
        uint256 timestamp
    );

    event LiquidityRemoved(
        address indexed user,
        address indexed pair,
        uint256 amount0,
        uint256 amount1,
        uint256 liquidity,
        uint256 timestamp
    );

    event TokensSwapped(
        address indexed user,
        address indexed tokenIn,
        address indexed tokenOut,
        uint256 amountIn,
        uint256 amountOut,
        uint256 timestamp
    );

    event LPTokensStaked(
        address indexed user,
        address indexed gauge,
        uint256 amount,
        uint256 timestamp
    );

    event LPTokensUnstaked(
        address indexed user,
        address indexed gauge,
        uint256 amount,
        uint256 timestamp
    );

    event RewardsClaimed(
        address indexed user,
        address indexed gauge,
        uint256 aeroAmount,
        uint256 timestamp
    );

    event StateAggregatorUpdated(address indexed oldAggregator, address indexed newAggregator);
    event MaxSlippageUpdated(uint256 oldSlippage, uint256 newSlippage);

    // ============ Constructor ============

    /**
     * @notice Initialize Aerodrome adapter
     * @param _router Aerodrome router address
     * @param _aeroToken AERO token address
     * @param _stateAggregator L2 state aggregator address
     */
    constructor(
        address _router,
        address _aeroToken,
        address _stateAggregator
    ) Ownable(msg.sender) {
        require(_router != address(0), "Invalid router");
        require(_aeroToken != address(0), "Invalid AERO token");
        require(_stateAggregator != address(0), "Invalid aggregator");

        router = IAerodromeRouter(_router);
        aeroToken = _aeroToken;
        stateAggregator = L2StateAggregator(_stateAggregator);
    }

    // ============ Liquidity Management ============

    /**
     * @notice Add liquidity to Aerodrome pool
     * @param tokenA First token address
     * @param tokenB Second token address
     * @param stable True for stable pool, false for volatile
     * @param amountADesired Amount of tokenA to add
     * @param amountBDesired Amount of tokenB to add
     * @param amountAMin Minimum amount of tokenA (slippage protection)
     * @param amountBMin Minimum amount of tokenB (slippage protection)
     * @param deadline Transaction deadline
     * @return amountA Actual amount of tokenA added
     * @return amountB Actual amount of tokenB added
     * @return liquidity LP tokens received
     */
    function addLiquidity(
        address tokenA,
        address tokenB,
        bool stable,
        uint256 amountADesired,
        uint256 amountBDesired,
        uint256 amountAMin,
        uint256 amountBMin,
        uint256 deadline
    ) external nonReentrant returns (
        uint256 amountA,
        uint256 amountB,
        uint256 liquidity
    ) {
        require(tokenA != address(0) && tokenB != address(0), "Invalid tokens");
        require(amountADesired > 0 && amountBDesired > 0, "Invalid amounts");
        require(deadline >= block.timestamp, "Deadline expired");

        // Transfer and approve
        IERC20(tokenA).safeTransferFrom(msg.sender, address(this), amountADesired);
        IERC20(tokenB).safeTransferFrom(msg.sender, address(this), amountBDesired);
        IERC20(tokenA).forceApprove(address(router), amountADesired);
        IERC20(tokenB).forceApprove(address(router), amountBDesired);

        // Add liquidity
        (amountA, amountB, liquidity) = router.addLiquidity(
            tokenA, tokenB, stable, amountADesired, amountBDesired,
            amountAMin, amountBMin, msg.sender, deadline
        );

        // Update position
        _updateLiquidityPosition(router.pairFor(tokenA, tokenB, stable), liquidity, amountA, amountB, true);

        // Return unused tokens
        if (amountADesired > amountA) IERC20(tokenA).safeTransfer(msg.sender, amountADesired - amountA);
        if (amountBDesired > amountB) IERC20(tokenB).safeTransfer(msg.sender, amountBDesired - amountB);

        emit LiquidityAdded(msg.sender, router.pairFor(tokenA, tokenB, stable), amountA, amountB, liquidity, block.timestamp);
    }

    function _updateLiquidityPosition(address pair, uint256 liquidity, uint256 amt0, uint256 amt1, bool isAdd) private {
        LiquidityPosition storage pos = userPositions[msg.sender][pair];
        if (isAdd) {
            if (pos.lpTokenAmount == 0) {
                if (++userPositionCount[msg.sender] == 1) activeUsers++;
            }
            pos.pair = pair;
            pos.lpTokenAmount += liquidity;
            pos.token0Amount += amt0;
            pos.token1Amount += amt1;
            totalLiquidityProvided += liquidity;
        } else {
            pos.lpTokenAmount -= liquidity;
            pos.token0Amount -= amt0;
            pos.token1Amount -= amt1;
            if (pos.lpTokenAmount == 0 && --userPositionCount[msg.sender] == 0) activeUsers--;
        }
        pos.lastUpdate = block.timestamp;
    }

    /**
     * @notice Remove liquidity from Aerodrome pool
     * @param tokenA First token address
     * @param tokenB Second token address
     * @param stable Pool type
     * @param liquidity Amount of LP tokens to burn
     * @param amountAMin Minimum tokenA to receive
     * @param amountBMin Minimum tokenB to receive
     * @param deadline Transaction deadline
     * @return amountA Amount of tokenA received
     * @return amountB Amount of tokenB received
     */
    function removeLiquidity(
        address tokenA,
        address tokenB,
        bool stable,
        uint256 liquidity,
        uint256 amountAMin,
        uint256 amountBMin,
        uint256 deadline
    ) external nonReentrant returns (uint256 amountA, uint256 amountB) {
        require(liquidity > 0, "Invalid liquidity");
        require(deadline >= block.timestamp, "Expired");

        address pair = router.pairFor(tokenA, tokenB, stable);
        require(userPositions[msg.sender][pair].lpTokenAmount >= liquidity, "Insufficient LP");

        // Transfer and approve
        IERC20(pair).safeTransferFrom(msg.sender, address(this), liquidity);
        IERC20(pair).forceApprove(address(router), liquidity);

        // Remove liquidity
        (amountA, amountB) = router.removeLiquidity(
            tokenA, tokenB, stable, liquidity,
            amountAMin, amountBMin, msg.sender, deadline
        );

        // Update position
        _updateLiquidityPosition(pair, liquidity, amountA, amountB, false);

        emit LiquidityRemoved(msg.sender, pair, amountA, amountB, liquidity, block.timestamp);
    }

    // ============ Swapping ============

    /**
     * @notice Swap exact tokens for tokens
     * @param amountIn Amount of input tokens
     * @param amountOutMin Minimum output amount (slippage protection)
     * @param routes Array of swap routes
     * @param deadline Transaction deadline
     * @return amounts Array of amounts for each swap step
     */
    function swapExactTokensForTokens(
        uint256 amountIn,
        uint256 amountOutMin,
        IAerodromeRouter.Route[] calldata routes,
        uint256 deadline
    ) external nonReentrant returns (uint256[] memory amounts) {
        require(amountIn > 0, "Invalid amount");
        require(routes.length > 0, "No routes");
        require(deadline >= block.timestamp, "Deadline expired");

        address tokenIn = routes[0].from;
        address tokenOut = routes[routes.length - 1].to;

        // Transfer input tokens from user
        IERC20(tokenIn).safeTransferFrom(msg.sender, address(this), amountIn);

        // Approve router
        IERC20(tokenIn).forceApprove(address(router), amountIn);

        // Execute swap
        amounts = router.swapExactTokensForTokens(
            amountIn,
            amountOutMin,
            routes,
            msg.sender,
            deadline
        );

        uint256 amountOut = amounts[amounts.length - 1];

        // Record swap
        userSwapHistory[msg.sender].push(SwapRecord({
            tokenIn: tokenIn,
            tokenOut: tokenOut,
            amountIn: amountIn,
            amountOut: amountOut,
            timestamp: block.timestamp
        }));

        totalSwapVolume += amountIn;

        emit TokensSwapped(msg.sender, tokenIn, tokenOut, amountIn, amountOut, block.timestamp);
    }

    // ============ Staking ============

    /**
     * @notice Stake LP tokens in gauge to earn AERO rewards
     * @param gauge Gauge address
     * @param amount Amount of LP tokens to stake
     */
    function stakeLPTokens(
        address gauge,
        uint256 amount
    ) external nonReentrant {
        require(gauge != address(0), "Invalid gauge");
        require(amount > 0, "Invalid amount");

        IAerodromeGauge gaugeContract = IAerodromeGauge(gauge);
        address stakingToken = gaugeContract.stakingToken();

        // Transfer LP tokens from user
        IERC20(stakingToken).safeTransferFrom(msg.sender, address(this), amount);

        // Approve gauge
        IERC20(stakingToken).forceApprove(gauge, amount);

        // Stake in gauge
        gaugeContract.deposit(amount);

        // Update position
        LiquidityPosition storage position = userPositions[msg.sender][stakingToken];
        position.stakedAmount += amount;
        position.gauge = gauge;
        position.lastUpdate = block.timestamp;

        emit LPTokensStaked(msg.sender, gauge, amount, block.timestamp);
    }

    /**
     * @notice Unstake LP tokens from gauge
     * @param gauge Gauge address
     * @param amount Amount to unstake
     */
    function unstakeLPTokens(
        address gauge,
        uint256 amount
    ) external nonReentrant {
        require(amount > 0, "Invalid amount");

        IAerodromeGauge gaugeContract = IAerodromeGauge(gauge);
        address stakingToken = gaugeContract.stakingToken();

        LiquidityPosition storage position = userPositions[msg.sender][stakingToken];
        require(position.stakedAmount >= amount, "Insufficient staked amount");

        // Withdraw from gauge
        gaugeContract.withdraw(amount);

        // Transfer LP tokens to user
        IERC20(stakingToken).safeTransfer(msg.sender, amount);

        // Update position
        position.stakedAmount -= amount;
        position.lastUpdate = block.timestamp;

        emit LPTokensUnstaked(msg.sender, gauge, amount, block.timestamp);
    }

    /**
     * @notice Claim AERO rewards from gauge
     * @param gauge Gauge address
     * @return reward Amount of AERO claimed
     */
    function claimRewards(address gauge) external nonReentrant returns (uint256 reward) {
        IAerodromeGauge gaugeContract = IAerodromeGauge(gauge);

        // Claim rewards
        gaugeContract.getReward(address(this));

        // Get AERO balance
        reward = IERC20(aeroToken).balanceOf(address(this));

        if (reward > 0) {
            // Transfer AERO to user
            IERC20(aeroToken).safeTransfer(msg.sender, reward);
            totalAeroEarned += reward;

            emit RewardsClaimed(msg.sender, gauge, reward, block.timestamp);
        }
    }

    // ============ View Functions ============

    /**
     * @notice Get optimal swap route quote
     * @param amountIn Input amount
     * @param routes Swap routes
     * @return amounts Expected output amounts
     */
    function getAmountsOut(
        uint256 amountIn,
        IAerodromeRouter.Route[] memory routes
    ) external view returns (uint256[] memory amounts) {
        return router.getAmountsOut(amountIn, routes);
    }

    /**
     * @notice Get user's liquidity position
     * @param user User address
     * @param pair Pair address
     * @return position Liquidity position details
     */
    function getUserPosition(
        address user,
        address pair
    ) external view returns (LiquidityPosition memory position) {
        return userPositions[user][pair];
    }

    /**
     * @notice Get user's swap history
     * @param user User address
     * @return swaps Array of swap records
     */
    function getUserSwapHistory(
        address user
    ) external view returns (SwapRecord[] memory swaps) {
        return userSwapHistory[user];
    }

    /**
     * @notice Get pending AERO rewards
     * @param gauge Gauge address
     * @param user User address
     * @return reward Pending rewards
     */
    function getPendingRewards(
        address gauge,
        address user
    ) external view returns (uint256 reward) {
        return IAerodromeGauge(gauge).earned(user);
    }

    // ============ Admin Functions ============

    /**
     * @notice Update state aggregator
     * @param newAggregator New aggregator address
     */
    function setStateAggregator(address newAggregator) external onlyOwner {
        require(newAggregator != address(0), "Invalid aggregator");
        address oldAggregator = address(stateAggregator);
        stateAggregator = L2StateAggregator(newAggregator);
        emit StateAggregatorUpdated(oldAggregator, newAggregator);
    }

    /**
     * @notice Update max slippage tolerance
     * @param newSlippage New slippage in basis points
     */
    function setMaxSlippage(uint256 newSlippage) external onlyOwner {
        require(newSlippage <= 1000, "Slippage too high"); // Max 10%
        uint256 oldSlippage = maxSlippage;
        maxSlippage = newSlippage;
        emit MaxSlippageUpdated(oldSlippage, newSlippage);
    }

    /**
     * @notice Emergency token rescue
     * @param token Token to rescue
     * @param amount Amount to rescue
     */
    function rescueTokens(address token, uint256 amount) external onlyOwner {
        IERC20(token).safeTransfer(owner(), amount);
    }
}
