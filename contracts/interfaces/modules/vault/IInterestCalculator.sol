// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IInterestCalculator
 * @notice Interface for interest calculation module
 * @dev Calculates and manages interest accrual on debt
 */
interface IInterestCalculator {
    // ============ Events ============

    event InterestAccrued(address indexed user, uint256 interestAmount, uint256 newDebt);
    event StabilityFeeUpdated(uint256 oldFee, uint256 newFee);

    // ============ Functions ============

    /**
     * @notice Calculate accrued interest for a user
     * @param user User address
     * @param principal Principal debt amount
     * @param lastUpdate Last interest update timestamp
     * @return Accrued interest amount
     */
    function calculateInterest(address user, uint256 principal, uint256 lastUpdate)
        external
        view
        returns (uint256);

    /**
     * @notice Accrue interest for a user
     * @param user User address
     * @param principal Principal debt amount
     * @param lastUpdate Last interest update timestamp
     * @return newDebt New total debt including interest
     */
    function accrueInterest(address user, uint256 principal, uint256 lastUpdate)
        external
        returns (uint256 newDebt);

    /**
     * @notice Get current stability fee (annual interest rate)
     * @return Stability fee in basis points
     */
    function getStabilityFee() external view returns (uint256);

    /**
     * @notice Set stability fee
     * @param newFee New fee in basis points
     */
    function setStabilityFee(uint256 newFee) external;

    /**
     * @notice Calculate interest for a time period
     * @param principal Principal amount
     * @param duration Duration in seconds
     * @return Interest amount
     */
    function calculateInterestForPeriod(uint256 principal, uint256 duration)
        external
        view
        returns (uint256);

    /**
     * @notice Get compounded interest rate
     * @param principal Principal amount
     * @param duration Duration in seconds
     * @return Compounded amount
     */
    function getCompoundedAmount(uint256 principal, uint256 duration)
        external
        view
        returns (uint256);
}
