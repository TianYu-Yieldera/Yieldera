// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title ILiquidationEngine
 * @notice Interface for liquidation engine module
 * @dev Handles all liquidation logic and calculations
 */
interface ILiquidationEngine {
    // ============ Structs ============

    struct LiquidationParams {
        uint256 collateralToSeize;
        uint256 debtToRepay;
        uint256 penalty;
        uint256 liquidatorReward;
    }

    // ============ Events ============

    event PositionLiquidated(
        address indexed user,
        address indexed liquidator,
        uint256 collateralSeized,
        uint256 debtRepaid,
        uint256 penalty
    );
    event LiquidationThresholdUpdated(uint256 oldThreshold, uint256 newThreshold);
    event LiquidationPenaltyUpdated(uint256 oldPenalty, uint256 newPenalty);

    // ============ Functions ============

    /**
     * @notice Execute liquidation
     * @param user User to liquidate
     * @param liquidator Liquidator address
     * @param debtToCover Amount of debt to cover
     * @param userCollateral User's collateral amount
     * @param userDebt User's debt amount
     * @return params Liquidation parameters
     */
    function liquidate(
        address user,
        address liquidator,
        uint256 debtToCover,
        uint256 userCollateral,
        uint256 userDebt
    ) external returns (LiquidationParams memory params);

    /**
     * @notice Calculate liquidation parameters
     * @param debtToCover Amount of debt to cover
     * @param userCollateral User's collateral amount
     * @param userDebt User's debt amount
     * @return params Calculated liquidation parameters
     */
    function calculateLiquidation(
        uint256 debtToCover,
        uint256 userCollateral,
        uint256 userDebt
    ) external view returns (LiquidationParams memory params);

    /**
     * @notice Check if position can be liquidated
     * @param collateral Collateral amount
     * @param debt Debt amount
     * @param collateralRatio Current collateral ratio
     * @return True if liquidatable
     */
    function canLiquidate(uint256 collateral, uint256 debt, uint256 collateralRatio)
        external
        view
        returns (bool);

    /**
     * @notice Get liquidation threshold
     * @return Liquidation threshold percentage
     */
    function getLiquidationThreshold() external view returns (uint256);

    /**
     * @notice Get liquidation penalty
     * @return Liquidation penalty percentage
     */
    function getLiquidationPenalty() external view returns (uint256);

    /**
     * @notice Set liquidation threshold
     * @param newThreshold New threshold percentage
     */
    function setLiquidationThreshold(uint256 newThreshold) external;

    /**
     * @notice Set liquidation penalty
     * @param newPenalty New penalty percentage
     */
    function setLiquidationPenalty(uint256 newPenalty) external;

    /**
     * @notice Calculate maximum liquidation amount
     * @param userCollateral User's collateral
     * @param userDebt User's debt
     * @return Maximum amount that can be liquidated
     */
    function getMaxLiquidationAmount(uint256 userCollateral, uint256 userDebt)
        external
        view
        returns (uint256);
}
