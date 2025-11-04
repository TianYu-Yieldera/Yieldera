// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IPriceOracle
 * @notice Price oracle interface for asset valuation
 */
interface IPriceOracle {
    /**
     * @notice Get the latest price of an asset in USD
     * @param asset Address of the asset token
     * @return price Price with 8 decimals (e.g., $2000.00000000 for 1 ETH)
     */
    function getAssetPrice(address asset) external view returns (uint256 price);

    /**
     * @notice Get the number of decimals for price data
     * @return decimals Number of decimals (typically 8)
     */
    function decimals() external view returns (uint8 decimals);
}
