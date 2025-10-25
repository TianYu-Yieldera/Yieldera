// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IModule.sol";

/**
 * @title IVaultModule
 * @notice Standard interface for vault/collateral management modules
 * @dev Extends IModule with vault-specific functionality
 *
 * Purpose:
 * - Manage user collateral deposits
 * - Track debt positions
 * - Handle liquidations
 * - Calculate collateral ratios
 */
interface IVaultModule is IModule {
    /**
     * @notice User position data
     */
    struct Position {
        uint256 collateralAmount;   // Amount of collateral deposited
        uint256 debtAmount;          // Amount of debt owed
        uint256 lastInterestUpdate;  // Last interest calculation timestamp
        uint256 accruedInterest;     // Accumulated interest
        bool isActive;               // Position status
    }

    /**
     * @notice Vault configuration parameters
     */
    struct VaultConfig {
        uint256 minCollateralRatio;     // Minimum collateral ratio (e.g., 150 = 150%)
        uint256 liquidationThreshold;   // Liquidation threshold (e.g., 120 = 120%)
        uint256 liquidationPenalty;     // Liquidation penalty (e.g., 10 = 10%)
        uint256 stabilityFee;           // Annual interest rate in basis points
        uint256 debtCeiling;            // Maximum total debt allowed
        uint256 minDebtAmount;          // Minimum debt per position
        bool isPaused;                  // Vault pause status
    }

    // ============ Events ============

    event CollateralDeposited(address indexed user, uint256 amount, uint256 newTotal);
    event CollateralWithdrawn(address indexed user, uint256 amount, uint256 remaining);
    event DebtIncreased(address indexed user, uint256 amount, uint256 newDebt);
    event DebtDecreased(address indexed user, uint256 amount, uint256 remainingDebt);
    event PositionLiquidated(
        address indexed user,
        address indexed liquidator,
        uint256 collateralSeized,
        uint256 debtRepaid
    );
    event InterestAccrued(address indexed user, uint256 interestAmount, uint256 newDebt);
    event VaultConfigUpdated(string parameter, uint256 oldValue, uint256 newValue);

    // ============ Collateral Management ============

    /**
     * @notice Deposit collateral into vault
     * @param amount Amount of collateral to deposit
     * @dev Transfers tokens from user to vault
     */
    function depositCollateral(uint256 amount) external;

    /**
     * @notice Withdraw collateral from vault
     * @param amount Amount of collateral to withdraw
     * @dev Ensures position remains healthy after withdrawal
     */
    function withdrawCollateral(uint256 amount) external;

    /**
     * @notice Get user's collateral balance
     * @param user User address
     * @return Collateral amount
     */
    function getCollateralBalance(address user) external view returns (uint256);

    /**
     * @notice Get total collateral locked in vault
     * @return Total collateral amount
     */
    function getTotalCollateral() external view returns (uint256);

    // ============ Debt Management ============

    /**
     * @notice Increase user's debt (mint stablecoin)
     * @param amount Amount of debt to add
     * @dev Requires sufficient collateral
     */
    function increaseDebt(uint256 amount) external;

    /**
     * @notice Decrease user's debt (repay/burn stablecoin)
     * @param amount Amount of debt to repay
     */
    function decreaseDebt(uint256 amount) external;

    /**
     * @notice Get user's total debt including interest
     * @param user User address
     * @return Total debt amount
     */
    function getTotalDebt(address user) external view returns (uint256);

    /**
     * @notice Get total debt across all users
     * @return Total system debt
     */
    function getSystemDebt() external view returns (uint256);

    /**
     * @notice Calculate accrued interest for user
     * @param user User address
     * @return Interest amount
     */
    function calculateAccruedInterest(address user) external view returns (uint256);

    /**
     * @notice Get maximum debt user can take
     * @param user User address
     * @return Maximum mintable amount
     */
    function getMaxMintable(address user) external view returns (uint256);

    // ============ Position Management ============

    /**
     * @notice Get user's complete position
     * @param user User address
     * @return Position structure
     */
    function getPosition(address user) external view returns (Position memory);

    /**
     * @notice Calculate collateral ratio for user
     * @param user User address
     * @return Collateral ratio (e.g., 200 = 200%)
     */
    function getCollateralRatio(address user) external view returns (uint256);

    /**
     * @notice Check if position is healthy (above minimum ratio)
     * @param user User address
     * @return True if position is healthy
     */
    function isPositionHealthy(address user) external view returns (bool);

    /**
     * @notice Check if position can be liquidated
     * @param user User address
     * @return True if position is liquidatable
     */
    function canLiquidate(address user) external view returns (bool);

    // ============ Liquidation ============

    /**
     * @notice Liquidate an undercollateralized position
     * @param user User to liquidate
     * @param debtToCover Amount of debt liquidator will repay
     * @return collateralSeized Amount of collateral seized
     * @dev Caller must provide stablecoin to cover debt
     */
    function liquidate(address user, uint256 debtToCover) external returns (uint256 collateralSeized);

    /**
     * @notice Calculate liquidation parameters
     * @param user User address
     * @param debtToCover Amount of debt to cover
     * @return collateralToSeize Amount of collateral to seize
     * @return liquidationPenalty Penalty amount
     */
    function calculateLiquidation(address user, uint256 debtToCover)
        external
        view
        returns (uint256 collateralToSeize, uint256 liquidationPenalty);

    /**
     * @notice Get all liquidatable positions
     * @return users Array of user addresses that can be liquidated
     */
    function getLiquidatablePositions() external view returns (address[] memory users);

    // ============ Configuration ============

    /**
     * @notice Get vault configuration
     * @return Vault configuration structure
     */
    function getVaultConfig() external view returns (VaultConfig memory);

    /**
     * @notice Update minimum collateral ratio
     * @param newRatio New minimum ratio
     * @dev Only callable by governance
     */
    function setMinCollateralRatio(uint256 newRatio) external;

    /**
     * @notice Update liquidation threshold
     * @param newThreshold New threshold
     * @dev Only callable by governance
     */
    function setLiquidationThreshold(uint256 newThreshold) external;

    /**
     * @notice Update stability fee (interest rate)
     * @param newFee New fee in basis points
     * @dev Only callable by governance
     */
    function setStabilityFee(uint256 newFee) external;

    /**
     * @notice Update debt ceiling
     * @param newCeiling New debt ceiling
     * @dev Only callable by governance
     */
    function setDebtCeiling(uint256 newCeiling) external;

    // ============ Statistics ============

    /**
     * @notice Get vault statistics
     * @return totalCollateral Total collateral locked
     * @return totalDebt Total debt issued
     * @return averageCollateralRatio System-wide average ratio
     * @return utilizationRate Debt ceiling utilization
     */
    function getVaultStats()
        external
        view
        returns (
            uint256 totalCollateral,
            uint256 totalDebt,
            uint256 averageCollateralRatio,
            uint256 utilizationRate
        );

    /**
     * @notice Get number of active positions
     * @return Number of positions with debt > 0
     */
    function getActivePositionCount() external view returns (uint256);

    /**
     * @notice Get collateral token address
     * @return Address of collateral token
     */
    function getCollateralToken() external view returns (address);

    /**
     * @notice Get debt token address
     * @return Address of debt/stablecoin token
     */
    function getDebtToken() external view returns (address);
}
