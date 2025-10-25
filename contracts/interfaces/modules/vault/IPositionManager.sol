// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IPositionManager
 * @notice Interface for position management module
 * @dev Manages user positions and health calculations
 */
interface IPositionManager {
    // ============ Structs ============

    struct Position {
        uint256 collateralAmount;
        uint256 debtAmount;
        uint256 lastInterestUpdate;
        uint256 accruedInterest;
        bool isActive;
    }

    struct PositionHealth {
        uint256 collateralRatio;
        bool isHealthy;
        bool canBeLiquidated;
        uint256 maxWithdrawable;
        uint256 maxMintable;
    }

    // ============ Events ============

    event PositionOpened(address indexed user);
    event PositionClosed(address indexed user);
    event PositionUpdated(address indexed user, uint256 collateral, uint256 debt);

    // ============ Functions ============

    /**
     * @notice Get user's position
     * @param user User address
     * @return Position data
     */
    function getPosition(address user) external view returns (Position memory);

    /**
     * @notice Update position
     * @param user User address
     * @param collateralAmount New collateral amount
     * @param debtAmount New debt amount
     */
    function updatePosition(address user, uint256 collateralAmount, uint256 debtAmount) external;

    /**
     * @notice Get position health
     * @param user User address
     * @param minCollateralRatio Minimum required ratio
     * @param liquidationThreshold Liquidation threshold
     * @return health Position health data
     */
    function getPositionHealth(
        address user,
        uint256 minCollateralRatio,
        uint256 liquidationThreshold
    ) external view returns (PositionHealth memory health);

    /**
     * @notice Calculate collateral ratio
     * @param collateral Collateral amount
     * @param debt Debt amount
     * @return Collateral ratio percentage
     */
    function calculateCollateralRatio(uint256 collateral, uint256 debt)
        external
        pure
        returns (uint256);

    /**
     * @notice Check if position is healthy
     * @param collateral Collateral amount
     * @param debt Debt amount
     * @param minRatio Minimum required ratio
     * @return True if healthy
     */
    function isPositionHealthy(uint256 collateral, uint256 debt, uint256 minRatio)
        external
        pure
        returns (bool);

    /**
     * @notice Get all active positions
     * @return Array of user addresses with active positions
     */
    function getActivePositions() external view returns (address[] memory);

    /**
     * @notice Get number of active positions
     * @return Count of active positions
     */
    function getActivePositionCount() external view returns (uint256);

    /**
     * @notice Check if user has an active position
     * @param user User address
     * @return True if user has active position
     */
    function hasActivePosition(address user) external view returns (bool);

    /**
     * @notice Open a new position for user
     * @param user User address
     */
    function openPosition(address user) external;

    /**
     * @notice Close user's position
     * @param user User address
     */
    function closePosition(address user) external;
}
