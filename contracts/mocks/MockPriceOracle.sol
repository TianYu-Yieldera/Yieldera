// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../layer2/interfaces/IPriceOracle.sol";

/**
 * @title MockPriceOracle
 * @notice Mock price oracle for testing
 * @dev Allows setting arbitrary prices for testing purposes
 */
contract MockPriceOracle is IPriceOracle {
    // Mapping from asset to price
    mapping(address => uint256) public prices;

    /**
     * @notice Set price for an asset (testing only)
     * @param asset Asset address
     * @param price Price with 8 decimals
     */
    function setAssetPrice(address asset, uint256 price) external {
        prices[asset] = price;
    }

    /**
     * @notice Get asset price
     * @param asset Asset address
     * @return price Price with 8 decimals
     */
    function getAssetPrice(address asset) external view override returns (uint256) {
        uint256 price = prices[asset];
        require(price > 0, "Price not set");
        return price;
    }

    /**
     * @notice Get decimals
     * @return Always returns 8
     */
    function decimals() external pure override returns (uint8) {
        return 8;
    }
}
