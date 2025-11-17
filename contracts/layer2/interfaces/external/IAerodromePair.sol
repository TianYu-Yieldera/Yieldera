// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IAerodromePair
 * @notice Interface for Aerodrome liquidity pair
 */
interface IAerodromePair {
    function token0() external view returns (address);
    function token1() external view returns (address);
    function stable() external view returns (bool);

    function getReserves() external view returns (
        uint256 reserve0,
        uint256 reserve1,
        uint256 blockTimestampLast
    );

    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);

    function mint(address to) external returns (uint256 liquidity);
    function burn(address to) external returns (uint256 amount0, uint256 amount1);

    function swap(
        uint256 amount0Out,
        uint256 amount1Out,
        address to,
        bytes calldata data
    ) external;
}
