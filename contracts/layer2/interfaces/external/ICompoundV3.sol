// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IComet
 * @notice Interface for Compound V3 (Comet) protocol
 * @dev Simplified interface for core lending/borrowing functionality
 *
 * Key Differences from Aave:
 * - Single base asset per market (e.g., USDC)
 * - Multiple collateral assets
 * - Simpler interest rate model
 * - No isolated lending pools
 */
interface IComet {
    /**
     * @notice Supply an asset to the protocol
     * @param asset Address of the asset to supply
     * @param amount Amount to supply
     */
    function supply(address asset, uint256 amount) external;

    /**
     * @notice Supply an asset on behalf of another address
     * @param from Address supplying the asset
     * @param dst Destination address receiving the balance
     * @param asset Address of the asset
     * @param amount Amount to supply
     */
    function supplyTo(address from, address dst, address asset, uint256 amount) external;

    /**
     * @notice Withdraw an asset from the protocol
     * @param asset Address of the asset to withdraw
     * @param amount Amount to withdraw
     */
    function withdraw(address asset, uint256 amount) external;

    /**
     * @notice Withdraw an asset to a specific address
     * @param to Destination address
     * @param asset Address of the asset
     * @param amount Amount to withdraw
     */
    function withdrawTo(address to, address asset, uint256 amount) external;

    /**
     * @notice Get collateral balance for an account
     * @param account User address
     * @param asset Collateral asset address
     * @return balance Collateral balance
     */
    function collateralBalanceOf(address account, address asset) external view returns (uint128);

    /**
     * @notice Get borrow balance for an account (base asset)
     * @param account User address
     * @return balance Borrow balance (negative if borrowing)
     */
    function borrowBalanceOf(address account) external view returns (uint256);

    /**
     * @notice Get supply balance for an account (base asset)
     * @param account User address
     * @return balance Supply balance
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @notice Get the base asset address
     * @return baseAsset Address of base token (e.g., USDC)
     */
    function baseToken() external view returns (address);

    /**
     * @notice Get asset information
     * @param asset Asset address
     * @return offset Asset configuration offset
     * @return asset Asset address
     * @return priceFeed Price feed address
     * @return scale Scale factor
     * @return borrowCollateralFactor Borrow collateral factor
     * @return liquidateCollateralFactor Liquidation collateral factor
     * @return liquidationFactor Liquidation bonus
     * @return supplyCap Supply cap
     */
    function getAssetInfo(uint8 offset) external view returns (
        uint8,
        address asset,
        address priceFeed,
        uint64 scale,
        uint64 borrowCollateralFactor,
        uint64 liquidateCollateralFactor,
        uint64 liquidationFactor,
        uint128 supplyCap
    );

    /**
     * @notice Get supply rate per second
     * @param utilization Current utilization rate
     * @return Supply rate per second (scaled by 1e18)
     */
    function getSupplyRate(uint256 utilization) external view returns (uint64);

    /**
     * @notice Get borrow rate per second
     * @param utilization Current utilization rate
     * @return Borrow rate per second (scaled by 1e18)
     */
    function getBorrowRate(uint256 utilization) external view returns (uint64);

    /**
     * @notice Get current utilization
     * @return Utilization rate (scaled by 1e18)
     */
    function getUtilization() external view returns (uint256);

    /**
     * @notice Absorb accounts with underwater positions
     * @param absorber Address performing absorption
     * @param accounts Array of accounts to absorb
     */
    function absorb(address absorber, address[] calldata accounts) external;

    /**
     * @notice Check if account is liquidatable
     * @param account Account to check
     * @return isLiquidatable True if account can be liquidated
     */
    function isLiquidatable(address account) external view returns (bool);

    /**
     * @notice Accrues interest and updates indexes
     */
    function accrueAccount(address account) external;
}

/**
 * @title ICometRewards
 * @notice Interface for Compound V3 rewards distribution
 */
interface ICometRewards {
    /**
     * @notice Claim COMP rewards
     * @param comet Comet instance
     * @param src Address to claim for
     * @param shouldAccrue Whether to accrue before claiming
     */
    function claim(address comet, address src, bool shouldAccrue) external;

    /**
     * @notice Get claimable rewards
     * @param comet Comet instance
     * @param account Account address
     * @return token Reward token address
     * @return owed Amount of COMP owed
     */
    function getRewardOwed(address comet, address account) external view returns (
        address token,
        uint256 owed
    );
}
