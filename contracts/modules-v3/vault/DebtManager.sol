// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/modules/vault/IDebtManager.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title DebtManager
 * @notice Manages user debt and system debt tracking
 * @dev Uses Diamond Storage for state management
 */
contract DebtManager is IDebtManager, Ownable {
    // ============ State Variables ============

    address public vaultModule; // Main coordinator contract

    // Storage for debt balances (using Diamond Storage pattern)
    bytes32 private constant DEBT_STORAGE_POSITION = keccak256("debt.manager.storage");

    struct DebtStorage {
        mapping(address => uint256) debts;
        uint256 totalDebt;
        uint256 debtCeiling;
    }

    // ============ Modifiers ============

    modifier onlyVaultModule() {
        require(msg.sender == vaultModule, "Only vault module");
        _;
    }

    // ============ Constructor ============

    constructor() {}

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
     * @notice Set debt ceiling
     * @param newCeiling New debt ceiling
     */
    function setDebtCeiling(uint256 newCeiling) external onlyOwner {
        DebtStorage storage ds = _getStorage();
        ds.debtCeiling = newCeiling;
    }

    // ============ Internal Storage Functions ============

    function _getStorage() private pure returns (DebtStorage storage ds) {
        bytes32 position = DEBT_STORAGE_POSITION;
        assembly {
            ds.slot := position
        }
    }

    // ============ IDebtManager Implementation ============

    /**
     * @notice Increase debt for a user
     * @param user User address
     * @param amount Amount to increase
     */
    function increaseDebt(address user, uint256 amount) external override onlyVaultModule {
        require(user != address(0), "Invalid user");
        require(amount > 0, "Amount must be > 0");

        DebtStorage storage ds = _getStorage();

        // Check debt ceiling
        if (ds.debtCeiling > 0) {
            require(ds.totalDebt + amount <= ds.debtCeiling, "Debt ceiling exceeded");
        }

        // Update storage
        ds.debts[user] += amount;
        ds.totalDebt += amount;

        uint256 newDebt = ds.debts[user];

        emit DebtIncreased(user, amount, newDebt);
    }

    /**
     * @notice Decrease debt for a user
     * @param user User address
     * @param amount Amount to decrease
     */
    function decreaseDebt(address user, uint256 amount) external override onlyVaultModule {
        require(user != address(0), "Invalid user");
        require(amount > 0, "Amount must be > 0");

        DebtStorage storage ds = _getStorage();
        require(ds.debts[user] >= amount, "Insufficient debt");

        // Update storage
        ds.debts[user] -= amount;
        ds.totalDebt -= amount;

        uint256 remainingDebt = ds.debts[user];

        emit DebtDecreased(user, amount, remainingDebt);
    }

    /**
     * @notice Get current debt for a user
     * @param user User address
     * @return Current debt amount
     */
    function getDebt(address user) external view override returns (uint256) {
        DebtStorage storage ds = _getStorage();
        return ds.debts[user];
    }

    /**
     * @notice Get total system debt
     * @return Total debt amount
     */
    function getTotalDebt() external view override returns (uint256) {
        DebtStorage storage ds = _getStorage();
        return ds.totalDebt;
    }

    /**
     * @notice Get maximum debt a user can take
     * @param user User address (unused, kept for interface compatibility)
     * @param collateralAmount User's collateral amount
     * @param collateralRatio Required collateral ratio (percentage, e.g., 150 = 150%)
     * @return Maximum debt amount
     */
    function getMaxDebt(address user, uint256 collateralAmount, uint256 collateralRatio)
        external
        view
        override
        returns (uint256)
    {
        // Avoid unused parameter warning
        user;

        if (collateralRatio == 0) return 0;

        // Max debt = collateral * 100 / ratio
        // Example: 1000 collateral, 150% ratio = 1000 * 100 / 150 = 666.66
        return (collateralAmount * 100) / collateralRatio;
    }

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
    ) external view override returns (bool) {
        DebtStorage storage ds = _getStorage();

        uint256 currentDebt = ds.debts[user];
        uint256 newDebt = currentDebt + amount;

        // Check debt ceiling
        if (ds.debtCeiling > 0 && ds.totalDebt + amount > ds.debtCeiling) {
            return false;
        }

        // Check if new debt maintains minimum ratio
        if (newDebt == 0) return true;

        uint256 ratio = (collateralAmount * 100) / newDebt;
        return ratio >= minRatio;
    }

    /**
     * @notice Get debt ceiling
     * @return Debt ceiling amount
     */
    function getDebtCeiling() external view returns (uint256) {
        DebtStorage storage ds = _getStorage();
        return ds.debtCeiling;
    }
}
