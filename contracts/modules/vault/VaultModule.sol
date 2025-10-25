// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/core/IVaultModule.sol";
import "../../plugins/core/BaseModule.sol";
import "../../core/CollateralVault.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title VaultModule
 * @notice Pluggable vault module - adapter for existing CollateralVault
 * @dev Wraps CollateralVault and implements IVaultModule interface
 */
contract VaultModule is IVaultModule, BaseModule, Ownable {
    // ============ Constants ============

    bytes32 public constant MODULE_ID = keccak256("VAULT_MODULE");
    string public constant MODULE_NAME = "VaultModule";
    string public constant MODULE_VERSION = "1.0.0";

    // ============ State Variables ============

    CollateralVault public immutable vault;
    address public immutable collateralToken;
    address public immutable debtToken;

    // Module configuration
    VaultConfig private _config;

    // Active position tracking
    mapping(address => bool) private _hasPosition;
    address[] private _activePositions;

    // ============ Constructor ============

    constructor(
        address _vault,
        address _collateralToken,
        address _debtToken
    ) BaseModule(MODULE_ID, MODULE_NAME, MODULE_VERSION) {
        require(_vault != address(0), "Invalid vault");
        require(_collateralToken != address(0), "Invalid collateral token");
        require(_debtToken != address(0), "Invalid debt token");

        vault = CollateralVault(_vault);
        collateralToken = _collateralToken;
        debtToken = _debtToken;

        // Initialize config from vault constants
        _config = VaultConfig({
            minCollateralRatio: vault.COLLATERAL_RATIO(),
            liquidationThreshold: vault.LIQUIDATION_THRESHOLD(),
            liquidationPenalty: 10, // 10%
            stabilityFee: vault.STABILITY_FEE(),
            debtCeiling: type(uint256).max, // No ceiling initially
            minDebtAmount: 0,
            isPaused: false
        });
    }

    // ============ BaseModule Overrides ============

    function getDependencies() external pure override returns (bytes32[] memory) {
        bytes32[] memory deps = new bytes32[](2);
        deps[0] = keccak256("PRICE_ORACLE_MODULE");
        deps[1] = keccak256("AUDIT_MODULE");
        return deps;
    }

    function healthCheck()
        external
        view
        override(IModule, BaseModule)
        returns (bool healthy, string memory message)
    {
        (bool baseHealthy, string memory baseMessage) = BaseModule.healthCheck();
        if (!baseHealthy) {
            return (false, baseMessage);
        }

        // Check vault health
        (uint256 totalCollateral, uint256 totalDebt, ) = vault.getVaultStats();

        if (totalDebt > 0 && totalCollateral == 0) {
            return (false, "Vault has debt but no collateral");
        }

        uint256 ratio = totalDebt > 0 ? (totalCollateral * 100) / totalDebt : type(uint256).max;
        if (ratio < _config.minCollateralRatio) {
            return (false, "System undercollateralized");
        }

        return (true, "Vault healthy");
    }

    // ============ Collateral Management ============

    function depositCollateral(uint256 amount) external override whenNotPaused {
        _requireActive();
        require(amount > 0, "Amount must be > 0");

        vault.depositCollateral(amount);

        if (!_hasPosition[msg.sender]) {
            _hasPosition[msg.sender] = true;
            _activePositions.push(msg.sender);
        }

        emit CollateralDeposited(msg.sender, amount, vault.collateralDeposited(msg.sender));
    }

    function withdrawCollateral(uint256 amount) external override whenNotPaused {
        _requireActive();
        require(amount > 0, "Amount must be > 0");

        vault.withdrawCollateral(amount);

        emit CollateralWithdrawn(msg.sender, amount, vault.collateralDeposited(msg.sender));
    }

    function getCollateralBalance(address user) external view override returns (uint256) {
        return vault.collateralDeposited(user);
    }

    function getTotalCollateral() external view override returns (uint256) {
        return vault.totalCollateral();
    }

    // ============ Debt Management ============

    function increaseDebt(uint256 amount) external override whenNotPaused {
        _requireActive();
        require(amount > 0, "Amount must be > 0");

        // Check debt ceiling
        uint256 newTotalDebt = vault.totalDebt() + amount;
        require(newTotalDebt <= _config.debtCeiling, "Exceeds debt ceiling");

        vault.increaseDebt(msg.sender, amount);

        if (!_hasPosition[msg.sender]) {
            _hasPosition[msg.sender] = true;
            _activePositions.push(msg.sender);
        }

        emit DebtIncreased(msg.sender, amount, vault.debtAmount(msg.sender));
    }

    function decreaseDebt(uint256 amount) external override whenNotPaused {
        _requireActive();
        require(amount > 0, "Amount must be > 0");

        vault.decreaseDebt(msg.sender, amount);

        emit DebtDecreased(msg.sender, amount, vault.debtAmount(msg.sender));
    }

    function getTotalDebt(address user) external view override returns (uint256) {
        return vault.getTotalDebt(user);
    }

    function getSystemDebt() external view override returns (uint256) {
        return vault.totalDebt();
    }

    function calculateAccruedInterest(address user) external view override returns (uint256) {
        return vault.accruedInterest(user);
    }

    function getMaxMintable(address user) external view override returns (uint256) {
        return vault.getMaxMintable(user);
    }

    // ============ Position Management ============

    function getPosition(address user) external view override returns (Position memory) {
        uint256 collateral = vault.collateralDeposited(user);
        uint256 debt = vault.debtAmount(user);
        uint256 interest = vault.accruedInterest(user);
        uint256 lastUpdate = vault.lastInterestUpdate(user);

        return Position({
            collateralAmount: collateral,
            debtAmount: debt,
            lastInterestUpdate: lastUpdate,
            accruedInterest: interest,
            isActive: debt > 0 || collateral > 0
        });
    }

    function getCollateralRatio(address user) external view override returns (uint256) {
        return vault.getCollateralRatio(user);
    }

    function isPositionHealthy(address user) external view override returns (bool) {
        return vault.isPositionHealthy(user);
    }

    function canLiquidate(address user) external view override returns (bool) {
        return vault.canLiquidate(user);
    }

    // ============ Liquidation ============

    function liquidate(address user, uint256 debtToCover)
        external
        override
        whenNotPaused
        returns (uint256 collateralSeized)
    {
        _requireActive();

        collateralSeized = vault.liquidate(user, debtToCover);

        emit PositionLiquidated(user, msg.sender, collateralSeized, debtToCover);

        return collateralSeized;
    }

    function calculateLiquidation(address user, uint256 debtToCover)
        external
        view
        override
        returns (uint256 collateralToSeize, uint256 liquidationPenalty)
    {
        // Simple calculation: assumes 1 LP = 1 USD
        collateralToSeize = (debtToCover * 110) / 100; // 10% penalty

        uint256 userCollateral = vault.collateralDeposited(user);
        if (collateralToSeize > userCollateral) {
            collateralToSeize = userCollateral;
        }

        liquidationPenalty = (debtToCover * _config.liquidationPenalty) / 100;

        return (collateralToSeize, liquidationPenalty);
    }

    function getLiquidatablePositions() external view override returns (address[] memory users) {
        uint256 count = 0;

        // Count liquidatable positions
        for (uint256 i = 0; i < _activePositions.length; i++) {
            if (vault.canLiquidate(_activePositions[i])) {
                count++;
            }
        }

        // Collect liquidatable addresses
        users = new address[](count);
        uint256 index = 0;

        for (uint256 i = 0; i < _activePositions.length; i++) {
            if (vault.canLiquidate(_activePositions[i])) {
                users[index] = _activePositions[i];
                index++;
            }
        }

        return users;
    }

    // ============ Configuration ============

    function getVaultConfig() external view override returns (VaultConfig memory) {
        return _config;
    }

    function setMinCollateralRatio(uint256 newRatio) external override onlyOwner {
        uint256 oldRatio = _config.minCollateralRatio;
        _config.minCollateralRatio = newRatio;
        emit VaultConfigUpdated("minCollateralRatio", oldRatio, newRatio);
    }

    function setLiquidationThreshold(uint256 newThreshold) external override onlyOwner {
        uint256 oldThreshold = _config.liquidationThreshold;
        _config.liquidationThreshold = newThreshold;
        emit VaultConfigUpdated("liquidationThreshold", oldThreshold, newThreshold);
    }

    function setStabilityFee(uint256 newFee) external override onlyOwner {
        uint256 oldFee = _config.stabilityFee;
        _config.stabilityFee = newFee;
        emit VaultConfigUpdated("stabilityFee", oldFee, newFee);
    }

    function setDebtCeiling(uint256 newCeiling) external override onlyOwner {
        uint256 oldCeiling = _config.debtCeiling;
        _config.debtCeiling = newCeiling;
        emit VaultConfigUpdated("debtCeiling", oldCeiling, newCeiling);
    }

    // ============ Statistics ============

    function getVaultStats()
        external
        view
        override
        returns (
            uint256 totalCollateral,
            uint256 totalDebt,
            uint256 averageCollateralRatio,
            uint256 utilizationRate
        )
    {
        (totalCollateral, totalDebt, averageCollateralRatio) = vault.getVaultStats();

        utilizationRate = _config.debtCeiling > 0
            ? (totalDebt * 10000) / _config.debtCeiling
            : 0;

        return (totalCollateral, totalDebt, averageCollateralRatio, utilizationRate);
    }

    function getActivePositionCount() external view override returns (uint256) {
        uint256 count = 0;

        for (uint256 i = 0; i < _activePositions.length; i++) {
            address user = _activePositions[i];
            if (vault.collateralDeposited(user) > 0 || vault.debtAmount(user) > 0) {
                count++;
            }
        }

        return count;
    }

    function getCollateralToken() external view override returns (address) {
        return collateralToken;
    }

    function getDebtToken() external view override returns (address) {
        return debtToken;
    }

    // ============ Override Required Functions ============

    function pause() external override(IModule, BaseModule) onlyOwner {
        _config.isPaused = true;
        BaseModule.pause();
    }

    function unpause() external override(IModule, BaseModule) onlyOwner {
        _config.isPaused = false;
        BaseModule.unpause();
    }
}
