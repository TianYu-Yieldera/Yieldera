// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "../interfaces/external/IUniswapV3.sol";
import "../aggregator/L2StateAggregator.sol";

/**
 * @title UniswapV3Adapter
 * @notice Adapter for Uniswap V3 DEX - Token swaps
 * @dev Integrates with Uniswap V3 SwapRouter for optimal trading
 *
 * Key Features:
 * - Single-hop token swaps
 * - Multi-hop routing for best prices
 * - Slippage protection
 * - Deadline management
 * - Volume tracking
 *
 * Note: This MVP focuses on swap functionality.
 * Liquidity management (concentrated liquidity, NFT positions) can be added later
 * as it requires more complex stack management and is not critical for initial DeFi aggregator.
 */
contract UniswapV3Adapter is AccessControl, ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;

    bytes32 public constant MANAGER_ROLE = keccak256("MANAGER_ROLE");

    // Uniswap V3 contracts
    ISwapRouter public immutable swapRouter;
    IUniswapV3Factory public immutable factory;

    // State aggregator integration
    L2StateAggregator public stateAggregator;

    // Statistics
    uint256 public totalSwapVolume;
    uint256 public totalSwaps;

    // Fee tiers
    uint24 public constant FEE_LOWEST = 100;     // 0.01%
    uint24 public constant FEE_LOW = 500;        // 0.05%
    uint24 public constant FEE_MEDIUM = 3000;    // 0.30%
    uint24 public constant FEE_HIGH = 10000;     // 1%

    // Events
    event Swapped(
        address indexed user,
        address indexed tokenIn,
        address indexed tokenOut,
        uint256 amountIn,
        uint256 amountOut,
        uint24 fee
    );

    event MultiHopSwap(
        address indexed user,
        uint256 amountIn,
        uint256 amountOut
    );

    /**
     * @notice Constructor
     * @param _swapRouter Uniswap V3 SwapRouter address
     * @param _factory Uniswap V3 Factory address
     * @param _admin Admin address
     */
    constructor(
        address _swapRouter,
        address _factory,
        address _admin
    ) {
        require(_swapRouter != address(0), "Invalid swap router");
        require(_factory != address(0), "Invalid factory");
        require(_admin != address(0), "Invalid admin");

        swapRouter = ISwapRouter(_swapRouter);
        factory = IUniswapV3Factory(_factory);

        _grantRole(DEFAULT_ADMIN_ROLE, _admin);
        _grantRole(MANAGER_ROLE, _admin);
    }

    // =============================================================
    //                        SWAP FUNCTIONS
    // =============================================================

    /**
     * @notice Swap exact input for maximum output (single hop)
     * @param tokenIn Input token address
     * @param tokenOut Output token address
     * @param fee Pool fee tier
     * @param amountIn Exact input amount
     * @param amountOutMinimum Minimum output amount (slippage protection)
     * @param deadline Transaction deadline
     * @return amountOut Actual output amount
     */
    function swapExactInputSingle(
        address tokenIn,
        address tokenOut,
        uint24 fee,
        uint256 amountIn,
        uint256 amountOutMinimum,
        uint256 deadline
    ) external nonReentrant whenNotPaused returns (uint256 amountOut) {
        require(amountIn > 0, "Invalid input amount");
        require(deadline >= block.timestamp, "Deadline passed");

        // Transfer tokens from user
        IERC20(tokenIn).safeTransferFrom(msg.sender, address(this), amountIn);

        // Approve SwapRouter
        IERC20(tokenIn).forceApprove(address(swapRouter), amountIn);

        // Execute swap
        ISwapRouter.ExactInputSingleParams memory params = ISwapRouter.ExactInputSingleParams({
            tokenIn: tokenIn,
            tokenOut: tokenOut,
            fee: fee,
            recipient: msg.sender,
            deadline: deadline,
            amountIn: amountIn,
            amountOutMinimum: amountOutMinimum,
            sqrtPriceLimitX96: 0
        });

        amountOut = swapRouter.exactInputSingle(params);

        // Update statistics
        totalSwapVolume += amountIn;
        totalSwaps++;
        _updateAggregator();

        emit Swapped(msg.sender, tokenIn, tokenOut, amountIn, amountOut, fee);
    }

    /**
     * @notice Swap exact input for maximum output (multi-hop)
     * @param path Encoded path (token0, fee0, token1, fee1, token2...)
     * @param amountIn Exact input amount
     * @param amountOutMinimum Minimum output amount
     * @param deadline Transaction deadline
     * @return amountOut Actual output amount
     */
    function swapExactInput(
        bytes memory path,
        uint256 amountIn,
        uint256 amountOutMinimum,
        uint256 deadline
    ) external nonReentrant whenNotPaused returns (uint256 amountOut) {
        require(amountIn > 0, "Invalid input amount");
        require(deadline >= block.timestamp, "Deadline passed");

        // Extract first token from path
        address tokenIn = _getFirstTokenFromPath(path);

        // Transfer tokens from user
        IERC20(tokenIn).safeTransferFrom(msg.sender, address(this), amountIn);

        // Approve SwapRouter
        IERC20(tokenIn).forceApprove(address(swapRouter), amountIn);

        // Execute multi-hop swap
        ISwapRouter.ExactInputParams memory params = ISwapRouter.ExactInputParams({
            path: path,
            recipient: msg.sender,
            deadline: deadline,
            amountIn: amountIn,
            amountOutMinimum: amountOutMinimum
        });

        amountOut = swapRouter.exactInput(params);

        // Update statistics
        totalSwapVolume += amountIn;
        totalSwaps++;
        _updateAggregator();

        emit MultiHopSwap(msg.sender, amountIn, amountOut);
    }

    /**
     * @notice Swap for exact output (single hop)
     * @param tokenIn Input token address
     * @param tokenOut Output token address
     * @param fee Pool fee tier
     * @param amountOut Exact output amount desired
     * @param amountInMaximum Maximum input amount willing to spend
     * @param deadline Transaction deadline
     * @return amountIn Actual input amount spent
     */
    function swapExactOutputSingle(
        address tokenIn,
        address tokenOut,
        uint24 fee,
        uint256 amountOut,
        uint256 amountInMaximum,
        uint256 deadline
    ) external nonReentrant whenNotPaused returns (uint256 amountIn) {
        require(amountOut > 0, "Invalid output amount");
        require(deadline >= block.timestamp, "Deadline passed");

        // Transfer max tokens from user
        IERC20(tokenIn).safeTransferFrom(msg.sender, address(this), amountInMaximum);

        // Approve SwapRouter
        IERC20(tokenIn).forceApprove(address(swapRouter), amountInMaximum);

        // Execute swap
        ISwapRouter.ExactOutputSingleParams memory params = ISwapRouter.ExactOutputSingleParams({
            tokenIn: tokenIn,
            tokenOut: tokenOut,
            fee: fee,
            recipient: msg.sender,
            deadline: deadline,
            amountOut: amountOut,
            amountInMaximum: amountInMaximum,
            sqrtPriceLimitX96: 0
        });

        amountIn = swapRouter.exactOutputSingle(params);

        // Refund unused tokens
        if (amountInMaximum > amountIn) {
            IERC20(tokenIn).safeTransfer(msg.sender, amountInMaximum - amountIn);
        }

        // Update statistics
        totalSwapVolume += amountIn;
        totalSwaps++;
        _updateAggregator();

        emit Swapped(msg.sender, tokenIn, tokenOut, amountIn, amountOut, fee);
    }

    // =============================================================
    //                      VIEW FUNCTIONS
    // =============================================================

    /**
     * @notice Get pool address for token pair
     * @param tokenA First token
     * @param tokenB Second token
     * @param fee Fee tier
     * @return pool Pool address
     */
    function getPool(
        address tokenA,
        address tokenB,
        uint24 fee
    ) external view returns (address pool) {
        return factory.getPool(tokenA, tokenB, fee);
    }

    /**
     * @notice Check if pool exists
     * @param tokenA First token
     * @param tokenB Second token
     * @param fee Fee tier
     * @return exists True if pool exists
     */
    function poolExists(
        address tokenA,
        address tokenB,
        uint24 fee
    ) external view returns (bool exists) {
        address pool = factory.getPool(tokenA, tokenB, fee);
        return pool != address(0);
    }

    // =============================================================
    //                    INTERNAL FUNCTIONS
    // =============================================================

    /**
     * @notice Extract first token address from encoded path
     * @param path Encoded path
     * @return tokenIn First token address
     */
    function _getFirstTokenFromPath(bytes memory path) internal pure returns (address tokenIn) {
        require(path.length >= 20, "Invalid path");
        assembly {
            tokenIn := mload(add(path, 20))
        }
    }

    /**
     * @notice Update state aggregator
     */
    function _updateAggregator() internal {
        if (address(stateAggregator) != address(0)) {
            stateAggregator.updateSystemState(0, 0, totalSwaps, totalSwapVolume);
        }
    }

    // =============================================================
    //                     ADMIN FUNCTIONS
    // =============================================================

    /**
     * @notice Set state aggregator
     * @param _stateAggregator State aggregator address
     */
    function setStateAggregator(address _stateAggregator) external onlyRole(DEFAULT_ADMIN_ROLE) {
        stateAggregator = L2StateAggregator(_stateAggregator);
    }

    /**
     * @notice Pause contract
     */
    function pause() external onlyRole(MANAGER_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause contract
     */
    function unpause() external onlyRole(MANAGER_ROLE) {
        _unpause();
    }

    /**
     * @notice Emergency token withdrawal
     * @param token Token address
     * @param amount Amount to withdraw
     */
    function emergencyWithdraw(address token, uint256 amount) external onlyRole(DEFAULT_ADMIN_ROLE) {
        IERC20(token).safeTransfer(msg.sender, amount);
    }

    /**
     * @notice Receive ETH
     */
    receive() external payable {}
}
