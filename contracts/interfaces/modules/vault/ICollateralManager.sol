// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title ICollateralManager
 * @notice Interface for collateral management module
 * @dev Handles all collateral deposit and withdrawal operations
 */
interface ICollateralManager {
    // ============ Events ============

    event CollateralDeposited(address indexed user, uint256 amount, uint256 newTotal);
    event CollateralWithdrawn(address indexed user, uint256 amount, uint256 remaining);

    // ============ Functions ============

    /**
     * @notice Deposit collateral for a user
     * @param user User address
     * @param amount Amount to deposit
     */
    function deposit(address user, uint256 amount) external;

    /**
     * @notice Withdraw collateral for a user
     * @param user User address
     * @param amount Amount to withdraw
     */
    function withdraw(address user, uint256 amount) external;

    /**
     * @notice Get collateral balance for a user
     * @param user User address
     * @return Collateral amount
     */
    function getBalance(address user) external view returns (uint256);

    /**
     * @notice Get total collateral in the system
     * @return Total collateral amount
     */
    function getTotalCollateral() external view returns (uint256);

    /**
     * @notice Check if user has sufficient collateral
     * @param user User address
     * @param amount Required amount
     * @return True if user has sufficient collateral
     */
    function hasSufficientCollateral(address user, uint256 amount) external view returns (bool);

    /**
     * @notice Get collateral token address
     * @return Collateral token address
     */
    function getCollateralToken() external view returns (address);
}
