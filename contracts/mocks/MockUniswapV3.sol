// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/**
 * @title MockSwapRouter
 * @notice Mock implementation of Uniswap V3 SwapRouter for testing
 * @dev Simulates swaps with configurable exchange rates
 */
contract MockSwapRouter {
    using SafeERC20 for IERC20;

    // Exchange rate: 1 tokenIn = exchangeRate tokenOut (with 18 decimals precision)
    mapping(address => mapping(address => uint256)) public exchangeRates;

    // Track swap statistics
    uint256 public totalSwaps;

    event SwapExecuted(
        address indexed tokenIn,
        address indexed tokenOut,
        uint256 amountIn,
        uint256 amountOut
    );

    /**
     * @notice Set exchange rate between two tokens
     * @param tokenIn Input token
     * @param tokenOut Output token
     * @param rate Exchange rate (18 decimals)
     */
    function setExchangeRate(
        address tokenIn,
        address tokenOut,
        uint256 rate
    ) external {
        exchangeRates[tokenIn][tokenOut] = rate;
    }

    /**
     * @notice Swap exact input for maximum output (single hop)
     */
    function exactInputSingle(ExactInputSingleParams calldata params)
        external
        payable
        returns (uint256 amountOut)
    {
        require(params.amountIn > 0, "Invalid amount");

        // Calculate output amount using exchange rate
        uint256 rate = exchangeRates[params.tokenIn][params.tokenOut];
        require(rate > 0, "Exchange rate not set");

        amountOut = (params.amountIn * rate) / 1e18;
        require(amountOut >= params.amountOutMinimum, "Insufficient output");

        // Transfer tokens
        IERC20(params.tokenIn).safeTransferFrom(msg.sender, address(this), params.amountIn);
        IERC20(params.tokenOut).safeTransfer(params.recipient, amountOut);

        totalSwaps++;
        emit SwapExecuted(params.tokenIn, params.tokenOut, params.amountIn, amountOut);
    }

    /**
     * @notice Swap exact input for maximum output (multi-hop)
     */
    function exactInput(ExactInputParams calldata params)
        external
        payable
        returns (uint256 amountOut)
    {
        // For testing, we'll simulate a simple 2-hop swap
        // Extract first and last tokens from path
        (address tokenIn, address tokenOut) = _decodePath(params.path);

        uint256 rate = exchangeRates[tokenIn][tokenOut];
        require(rate > 0, "Exchange rate not set");

        amountOut = (params.amountIn * rate) / 1e18;
        require(amountOut >= params.amountOutMinimum, "Insufficient output");

        // Transfer tokens
        IERC20(tokenIn).safeTransferFrom(msg.sender, address(this), params.amountIn);
        IERC20(tokenOut).safeTransfer(params.recipient, amountOut);

        totalSwaps++;
        emit SwapExecuted(tokenIn, tokenOut, params.amountIn, amountOut);
    }

    /**
     * @notice Swap for exact output (single hop)
     */
    function exactOutputSingle(ExactOutputSingleParams calldata params)
        external
        payable
        returns (uint256 amountIn)
    {
        require(params.amountOut > 0, "Invalid amount");

        // Calculate input amount needed
        uint256 rate = exchangeRates[params.tokenIn][params.tokenOut];
        require(rate > 0, "Exchange rate not set");

        amountIn = (params.amountOut * 1e18) / rate;
        require(amountIn <= params.amountInMaximum, "Excessive input required");

        // Transfer tokens
        IERC20(params.tokenIn).safeTransferFrom(msg.sender, address(this), amountIn);
        IERC20(params.tokenOut).safeTransfer(params.recipient, params.amountOut);

        totalSwaps++;
        emit SwapExecuted(params.tokenIn, params.tokenOut, amountIn, params.amountOut);
    }

    /**
     * @notice Decode path to get first and last tokens
     */
    function _decodePath(bytes memory path)
        internal
        pure
        returns (address tokenIn, address tokenOut)
    {
        require(path.length >= 43, "Invalid path"); // 20 + 3 + 20 = 43 bytes minimum

        assembly {
            // First address is at position 32 (skip length) + 0
            tokenIn := div(mload(add(path, 32)), 0x1000000000000000000000000)
            // Last address starts at: 32 (length) + path.length - 20
            let lastAddrPos := add(add(path, 32), sub(mload(path), 20))
            tokenOut := div(mload(lastAddrPos), 0x1000000000000000000000000)
        }
    }

    // Structs matching Uniswap V3 interface
    struct ExactInputSingleParams {
        address tokenIn;
        address tokenOut;
        uint24 fee;
        address recipient;
        uint256 deadline;
        uint256 amountIn;
        uint256 amountOutMinimum;
        uint160 sqrtPriceLimitX96;
    }

    struct ExactInputParams {
        bytes path;
        address recipient;
        uint256 deadline;
        uint256 amountIn;
        uint256 amountOutMinimum;
    }

    struct ExactOutputSingleParams {
        address tokenIn;
        address tokenOut;
        uint24 fee;
        address recipient;
        uint256 deadline;
        uint256 amountOut;
        uint256 amountInMaximum;
        uint160 sqrtPriceLimitX96;
    }
}

/**
 * @title MockUniswapV3Factory
 * @notice Mock implementation of Uniswap V3 Factory for testing
 */
contract MockUniswapV3Factory {
    mapping(address => mapping(address => mapping(uint24 => address))) public pools;

    event PoolCreated(
        address indexed token0,
        address indexed token1,
        uint24 indexed fee,
        address pool
    );

    /**
     * @notice Get pool address for token pair and fee
     */
    function getPool(
        address tokenA,
        address tokenB,
        uint24 fee
    ) external view returns (address pool) {
        (address token0, address token1) = tokenA < tokenB
            ? (tokenA, tokenB)
            : (tokenB, tokenA);
        return pools[token0][token1][fee];
    }

    /**
     * @notice Create a new pool
     */
    function createPool(
        address tokenA,
        address tokenB,
        uint24 fee
    ) external returns (address pool) {
        (address token0, address token1) = tokenA < tokenB
            ? (tokenA, tokenB)
            : (tokenB, tokenA);

        require(pools[token0][token1][fee] == address(0), "Pool exists");

        // For testing, use deterministic address
        pool = address(uint160(uint256(keccak256(abi.encodePacked(token0, token1, fee)))));
        pools[token0][token1][fee] = pool;

        emit PoolCreated(token0, token1, fee, pool);
    }

    /**
     * @notice Set pool address manually for testing
     */
    function setPool(
        address tokenA,
        address tokenB,
        uint24 fee,
        address pool
    ) external {
        (address token0, address token1) = tokenA < tokenB
            ? (tokenA, tokenB)
            : (tokenB, tokenA);
        pools[token0][token1][fee] = pool;
    }
}
