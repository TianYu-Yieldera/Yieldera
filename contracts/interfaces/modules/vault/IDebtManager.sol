// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IDebtManager
 * @notice Interface for debt management module
 * @dev Handles all debt increase and decrease operations
 */
interface IDebtManager {
    // ============ Events ============

    event DebtIncreased(address indexed user, uint256 amount, uint256 newDebt);
    event DebtDecreased(address indexed user, uint256 amount, uint256 remainingDebt);

    // ============ Functions ============

    /**
     * @notice Increase debt for a user
     * @param user User address
     * @param amount Amount to increase
     */
    function increaseDebt(address user, uint256 amount) external;

    /**
     * @notice Decrease debt for a user
     * @param user User address
     * @param amount Amount to decrease
     */
    function decreaseDebt(address user, uint256 amount) external;

    /**
     * @notice Get current debt for a user
     * @param user User address
     * @return Current debt amount
     */
    function getDebt(address user) external view returns (uint256);

    /**
     * @notice Get total system debt
     * @return Total debt amount
     */
    function getTotalDebt() external view returns (uint256);

    /**
     * @notice Get maximum debt a user can take
     * @param user User address
     * @param collateralAmount User's collateral amount
     * @param collateralRatio Required collateral ratio
     * @return Maximum debt amount
     */
    function getMaxDebt(address user, uint256 collateralAmount, uint256 collateralRatio)
        external
        view
        returns (uint256);

    /**
     * @notice Check if user can increase debt
     * @param user User address
     * @param amount Amount to increase
     * @param collateralAmount User's collateral
     * @param minRatio Minimum collateral ratio
     * @return True if debt can be increased
     */
    function canIncreaseDebt(
        address user,
        uint256 amount,
        uint256 collateralAmount,
        uint256 minRatio
    ) external view returns (bool);
}
