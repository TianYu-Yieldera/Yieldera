// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/modules/vault/ILiquidationEngine.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title LiquidationEngine
 * @notice Handles all liquidation logic and calculations
 * @dev Calculates liquidation parameters and executes liquidations
 */
contract LiquidationEngine is ILiquidationEngine, Ownable {
    // ============ Constants ============

    uint256 private constant PERCENTAGE_BASE = 100;

    // ============ State Variables ============

    address public vaultModule; // Main coordinator contract

    // Liquidation threshold (e.g., 120 = 120% collateral ratio)
    uint256 public liquidationThreshold;

    // Liquidation penalty (e.g., 10 = 10% penalty)
    uint256 public liquidationPenalty;

    // Liquidator reward percentage (e.g., 5 = 5% of penalty goes to liquidator)
    uint256 public liquidatorRewardPercentage;

    // Storage for liquidation tracking
    bytes32 private constant LIQUIDATION_STORAGE_POSITION =
        keccak256("liquidation.engine.storage");

    struct LiquidationStorage {
        mapping(address => uint256) totalLiquidations;
        mapping(address => uint256) lastLiquidation;
        uint256 totalLiquidationVolume;
    }

    // ============ Modifiers ============

    modifier onlyVaultModule() {
        require(msg.sender == vaultModule, "Only vault module");
        _;
    }

    // ============ Constructor ============

    constructor(uint256 _liquidationThreshold, uint256 _liquidationPenalty) {
        require(_liquidationThreshold > 0, "Invalid threshold");
        require(_liquidationPenalty <= 50, "Penalty too high"); // Max 50%

        liquidationThreshold = _liquidationThreshold;
        liquidationPenalty = _liquidationPenalty;
        liquidatorRewardPercentage = 50; // 50% of penalty goes to liquidator
    }

    // ============ Admin Functions ============

    /**
     * @notice Set vault module address
     * @param _vaultModule Vault module address
     */
    function setVaultModule(address _vaultModule) external onlyOwner {
        require(_vaultModule != address(0), "Invalid address");
        vaultModule = _vaultModule;
    }

    /**
     * @notice Set liquidator reward percentage
     * @param percentage Percentage of penalty (0-100)
     */
    function setLiquidatorRewardPercentage(uint256 percentage) external onlyOwner {
        require(percentage <= 100, "Invalid percentage");
        liquidatorRewardPercentage = percentage;
    }

    // ============ Internal Storage Functions ============

    function _getStorage() private pure returns (LiquidationStorage storage ls) {
        bytes32 position = LIQUIDATION_STORAGE_POSITION;
        assembly {
            ls.slot := position
        }
    }

    // ============ ILiquidationEngine Implementation ============

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
    ) external override onlyVaultModule returns (LiquidationParams memory params) {
        require(user != address(0), "Invalid user");
        require(liquidator != address(0), "Invalid liquidator");

        // Calculate liquidation parameters
        params = calculateLiquidation(debtToCover, userCollateral, userDebt);

        // Update statistics
        LiquidationStorage storage ls = _getStorage();
        ls.totalLiquidations[user]++;
        ls.lastLiquidation[user] = block.timestamp;
        ls.totalLiquidationVolume += params.debtToRepay;

        emit PositionLiquidated(
            user, liquidator, params.collateralToSeize, params.debtToRepay, params.penalty
        );

        return params;
    }

    /**
     * @notice Calculate liquidation parameters
     * @param debtToCover Amount of debt to cover
     * @param userCollateral User's collateral amount
     * @param userDebt User's debt amount
     * @return params Calculated liquidation parameters
     */
    function calculateLiquidation(uint256 debtToCover, uint256 userCollateral, uint256 userDebt)
        public
        view
        override
        returns (LiquidationParams memory params)
    {
        require(userDebt > 0, "No debt to liquidate");
        require(debtToCover > 0, "Invalid debt amount");

        // Cap debt to cover at user's total debt
        uint256 actualDebtToCover = debtToCover > userDebt ? userDebt : debtToCover;

        // Calculate base collateral to seize (proportional to debt being covered)
        uint256 baseCollateral = (userCollateral * actualDebtToCover) / userDebt;

        // Calculate penalty amount
        uint256 penaltyAmount = (baseCollateral * liquidationPenalty) / PERCENTAGE_BASE;

        // Total collateral to seize includes penalty
        uint256 totalCollateralToSeize = baseCollateral + penaltyAmount;

        // Ensure we don't seize more than available
        if (totalCollateralToSeize > userCollateral) {
            totalCollateralToSeize = userCollateral;
            // Recalculate actual values
            baseCollateral = (userCollateral * PERCENTAGE_BASE) / (PERCENTAGE_BASE + liquidationPenalty);
            penaltyAmount = userCollateral - baseCollateral;
            actualDebtToCover = (userDebt * baseCollateral) / userCollateral;
        }

        // Calculate liquidator reward (portion of penalty)
        uint256 liquidatorReward = (penaltyAmount * liquidatorRewardPercentage) / PERCENTAGE_BASE;

        params = LiquidationParams({
            collateralToSeize: totalCollateralToSeize,
            debtToRepay: actualDebtToCover,
            penalty: penaltyAmount,
            liquidatorReward: liquidatorReward
        });

        return params;
    }

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
        override
        returns (bool)
    {
        // Avoid unused parameter warning
        collateral;

        if (debt == 0) return false;
        return collateralRatio < liquidationThreshold;
    }

    /**
     * @notice Get liquidation threshold
     * @return Liquidation threshold percentage
     */
    function getLiquidationThreshold() external view override returns (uint256) {
        return liquidationThreshold;
    }

    /**
     * @notice Get liquidation penalty
     * @return Liquidation penalty percentage
     */
    function getLiquidationPenalty() external view override returns (uint256) {
        return liquidationPenalty;
    }

    /**
     * @notice Set liquidation threshold
     * @param newThreshold New threshold percentage
     */
    function setLiquidationThreshold(uint256 newThreshold) external override onlyOwner {
        require(newThreshold > 0 && newThreshold <= 200, "Invalid threshold");
        uint256 oldThreshold = liquidationThreshold;
        liquidationThreshold = newThreshold;
        emit LiquidationThresholdUpdated(oldThreshold, newThreshold);
    }

    /**
     * @notice Set liquidation penalty
     * @param newPenalty New penalty percentage
     */
    function setLiquidationPenalty(uint256 newPenalty) external override onlyOwner {
        require(newPenalty <= 50, "Penalty too high"); // Max 50%
        uint256 oldPenalty = liquidationPenalty;
        liquidationPenalty = newPenalty;
        emit LiquidationPenaltyUpdated(oldPenalty, newPenalty);
    }

    /**
     * @notice Calculate maximum liquidation amount
     * @param userCollateral User's collateral
     * @param userDebt User's debt
     * @return Maximum amount that can be liquidated
     */
    function getMaxLiquidationAmount(uint256 userCollateral, uint256 userDebt)
        external
        view
        override
        returns (uint256)
    {
        if (userDebt == 0) return 0;

        // In this implementation, we allow full liquidation
        // Some systems may implement partial liquidation limits
        return userDebt;
    }

    /**
     * @notice Get total liquidations for a user
     * @param user User address
     * @return Number of times liquidated
     */
    function getTotalLiquidations(address user) external view returns (uint256) {
        LiquidationStorage storage ls = _getStorage();
        return ls.totalLiquidations[user];
    }

    /**
     * @notice Get last liquidation timestamp
     * @param user User address
     * @return Timestamp of last liquidation
     */
    function getLastLiquidation(address user) external view returns (uint256) {
        LiquidationStorage storage ls = _getStorage();
        return ls.lastLiquidation[user];
    }

    /**
     * @notice Get total liquidation volume
     * @return Total debt liquidated across all users
     */
    function getTotalLiquidationVolume() external view returns (uint256) {
        LiquidationStorage storage ls = _getStorage();
        return ls.totalLiquidationVolume;
    }
}
