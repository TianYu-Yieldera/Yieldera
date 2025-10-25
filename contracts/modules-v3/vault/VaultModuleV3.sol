// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/core/IVaultModule.sol";
import "../../interfaces/IUpgradeable.sol";
import "../../interfaces/modules/vault/ICollateralManager.sol";
import "../../interfaces/modules/vault/IPositionManager.sol";
import "../../interfaces/modules/vault/IDebtManager.sol";
import "../../interfaces/modules/vault/IInterestCalculator.sol";
import "../../interfaces/modules/vault/ILiquidationEngine.sol";
import "../../plugins/core/BaseModule.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/**
 * @title VaultModuleV3
 * @notice Modularized vault module - Phase 4 implementation
 * @dev Coordinates multiple sub-modules for vault operations
 *
 * Architecture:
 * - Uses composition pattern instead of inheritance
 * - Each business logic is in a separate module
 * - This contract acts as a coordinator
 * - Maintains upgrade capability via UUPS
 */
contract VaultModuleV3 is
    Initializable,
    IVaultModule,
    IUpgradeable,
    BaseModule,
    OwnableUpgradeable,
    UUPSUpgradeable
{
    // ============ Constants ============

    bytes32 public constant MODULE_ID = keccak256("VAULT_MODULE");
    string public constant MODULE_NAME = "VaultModule";
    string public constant MODULE_VERSION = "3.0.0";

    // ============ Sub-Modules ============

    ICollateralManager public collateralManager;
    IPositionManager public positionManager;
    IDebtManager public debtManager;
    IInterestCalculator public interestCalculator;
    ILiquidationEngine public liquidationEngine;

    // ============ Configuration ============

    VaultConfig private _config;

    // Legacy vault for backward compatibility
    address public legacyVault;
    address public debtToken;

    // Upgrade tracking
    struct UpgradeRecord {
        address implementation;
        uint256 timestamp;
        string version;
    }

    UpgradeRecord[] private _upgradeHistory;
    mapping(address => bool) private _authorizedUpgraders;
    bool private _upgradesPaused;

    // ============ Events ============

    event SubModuleUpdated(string indexed moduleName, address oldAddress, address newAddress);
    event ConfigUpdated(string param, uint256 oldValue, uint256 newValue);

    // ============ Constructor & Initializer ============

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initialize the modular vault
     * @param _collateralManager Collateral manager address
     * @param _positionManager Position manager address
     * @param _debtManager Debt manager address
     * @param _interestCalculator Interest calculator address
     * @param _liquidationEngine Liquidation engine address
     * @param _legacyVault Legacy vault address (for backward compatibility)
     * @param _debtToken Debt token address
     */
    function initialize(
        address _collateralManager,
        address _positionManager,
        address _debtManager,
        address _interestCalculator,
        address _liquidationEngine,
        address _legacyVault,
        address _debtToken
    ) public initializer {
        __Ownable_init();
        __UUPSUpgradeable_init();

        require(_collateralManager != address(0), "Invalid collateral manager");
        require(_positionManager != address(0), "Invalid position manager");
        require(_debtManager != address(0), "Invalid debt manager");
        require(_interestCalculator != address(0), "Invalid interest calculator");
        require(_liquidationEngine != address(0), "Invalid liquidation engine");
        require(_legacyVault != address(0), "Invalid legacy vault");
        require(_debtToken != address(0), "Invalid debt token");

        collateralManager = ICollateralManager(_collateralManager);
        positionManager = IPositionManager(_positionManager);
        debtManager = IDebtManager(_debtManager);
        interestCalculator = IInterestCalculator(_interestCalculator);
        liquidationEngine = ILiquidationEngine(_liquidationEngine);
        legacyVault = _legacyVault;
        debtToken = _debtToken;

        // Initialize default config
        _config = VaultConfig({
            minCollateralRatio: 150,
            liquidationThreshold: 120,
            liquidationPenalty: 10,
            stabilityFee: 200,
            debtCeiling: type(uint256).max,
            minDebtAmount: 0,
            isPaused: false
        });

        // Record deployment
        _upgradeHistory.push(
            UpgradeRecord({
                implementation: address(this),
                timestamp: block.timestamp,
                version: MODULE_VERSION
            })
        );

        _authorizedUpgraders[msg.sender] = true;
    }

    // ============ Sub-Module Management ============

    /**
     * @notice Update collateral manager
     * @param newManager New manager address
     */
    function setCollateralManager(address newManager) external onlyOwner {
        require(newManager != address(0), "Invalid address");
        address oldManager = address(collateralManager);
        collateralManager = ICollateralManager(newManager);
        emit SubModuleUpdated("CollateralManager", oldManager, newManager);
    }

    /**
     * @notice Update position manager
     * @param newManager New manager address
     */
    function setPositionManager(address newManager) external onlyOwner {
        require(newManager != address(0), "Invalid address");
        address oldManager = address(positionManager);
        positionManager = IPositionManager(newManager);
        emit SubModuleUpdated("PositionManager", oldManager, newManager);
    }

    /**
     * @notice Update debt manager
     * @param newManager New manager address
     */
    function setDebtManager(address newManager) external onlyOwner {
        require(newManager != address(0), "Invalid address");
        address oldManager = address(debtManager);
        debtManager = IDebtManager(newManager);
        emit SubModuleUpdated("DebtManager", oldManager, newManager);
    }

    /**
     * @notice Update interest calculator
     * @param newCalculator New calculator address
     */
    function setInterestCalculator(address newCalculator) external onlyOwner {
        require(newCalculator != address(0), "Invalid address");
        address oldCalculator = address(interestCalculator);
        interestCalculator = IInterestCalculator(newCalculator);
        emit SubModuleUpdated("InterestCalculator", oldCalculator, newCalculator);
    }

    /**
     * @notice Update liquidation engine
     * @param newEngine New engine address
     */
    function setLiquidationEngine(address newEngine) external onlyOwner {
        require(newEngine != address(0), "Invalid address");
        address oldEngine = address(liquidationEngine);
        liquidationEngine = ILiquidationEngine(newEngine);
        emit SubModuleUpdated("LiquidationEngine", oldEngine, newEngine);
    }

    // ============ BaseModule Overrides ============

    function getModuleId() external pure override returns (bytes32) {
        return MODULE_ID;
    }

    function getVersion() external pure override returns (string memory) {
        return MODULE_VERSION;
    }

    function getDependencies() external pure override returns (bytes32[] memory) {
        bytes32[] memory deps = new bytes32[](2);
        deps[0] = keccak256("PRICE_ORACLE_MODULE");
        deps[1] = keccak256("AUDIT_MODULE");
        return deps;
    }

    function isActive() external view override returns (bool) {
        return !_config.isPaused;
    }

    function initialize(bytes calldata) external override {
        revert("Use initialize(address,address,address,address)");
    }

    function healthCheck()
        external
        view
        override(IModule, BaseModule)
        returns (bool healthy, string memory message)
    {
        // Check if sub-modules are set
        if (address(collateralManager) == address(0)) {
            return (false, "CollateralManager not set");
        }
        if (address(positionManager) == address(0)) {
            return (false, "PositionManager not set");
        }
        if (address(debtManager) == address(0)) {
            return (false, "DebtManager not set");
        }
        if (address(interestCalculator) == address(0)) {
            return (false, "InterestCalculator not set");
        }
        if (address(liquidationEngine) == address(0)) {
            return (false, "LiquidationEngine not set");
        }

        // Check total collateral vs debt
        uint256 totalCollateral = collateralManager.getTotalCollateral();
        uint256 totalDebt = debtManager.getTotalDebt();

        // Verify system health: total collateral should exceed total debt
        if (totalDebt > 0) {
            uint256 systemRatio = (totalCollateral * 100) / totalDebt;
            if (systemRatio < _config.minCollateralRatio) {
                return (false, "System under-collateralized");
            }
        }

        return (true, "Vault healthy");
    }

    // ============ Collateral Management ============

    function depositCollateral(uint256 amount) external override {
        require(!_config.isPaused, "Paused");
        require(amount > 0, "Amount must be > 0");

        // Delegate to CollateralManager
        collateralManager.deposit(msg.sender, amount);

        // Update position
        uint256 currentCollateral = collateralManager.getBalance(msg.sender);
        positionManager.updatePosition(msg.sender, currentCollateral, 0);

        emit CollateralDeposited(msg.sender, amount, currentCollateral);
    }

    function withdrawCollateral(uint256 amount) external override {
        require(!_config.isPaused, "Paused");
        require(amount > 0, "Amount must be > 0");

        uint256 currentCollateral = collateralManager.getBalance(msg.sender);
        require(currentCollateral >= amount, "Insufficient collateral");

        // TODO: Check position health before withdrawal
        // For now, allow withdrawal

        // Delegate to CollateralManager
        collateralManager.withdraw(msg.sender, amount);

        // Update position
        uint256 newCollateral = collateralManager.getBalance(msg.sender);
        positionManager.updatePosition(msg.sender, newCollateral, 0);

        emit CollateralWithdrawn(msg.sender, amount, newCollateral);
    }

    function getCollateralBalance(address user) external view override returns (uint256) {
        return collateralManager.getBalance(user);
    }

    function getTotalCollateral() external view override returns (uint256) {
        return collateralManager.getTotalCollateral();
    }

    // ============ Debt Management ============

    function increaseDebt(uint256 amount) external override {
        require(!_config.isPaused, "Paused");
        require(amount > 0, "Amount must be > 0");

        // Get current position
        uint256 currentCollateral = collateralManager.getBalance(msg.sender);
        uint256 currentDebt = debtManager.getDebt(msg.sender);

        // Check if debt can be increased
        require(
            debtManager.canIncreaseDebt(msg.sender, amount, currentCollateral, _config.minCollateralRatio),
            "Insufficient collateral ratio"
        );

        // Delegate to DebtManager
        debtManager.increaseDebt(msg.sender, amount);

        // Update position
        uint256 newDebt = debtManager.getDebt(msg.sender);
        positionManager.updatePosition(msg.sender, currentCollateral, newDebt);

        // Mint debt tokens to user
        IERC20(debtToken).transfer(msg.sender, amount);

        emit DebtIncreased(msg.sender, amount, newDebt);
    }

    function decreaseDebt(uint256 amount) external override {
        require(!_config.isPaused, "Paused");
        require(amount > 0, "Amount must be > 0");

        uint256 currentDebt = debtManager.getDebt(msg.sender);
        require(currentDebt >= amount, "Insufficient debt");

        // Burn debt tokens from user
        IERC20(debtToken).transferFrom(msg.sender, address(this), amount);

        // Delegate to DebtManager
        debtManager.decreaseDebt(msg.sender, amount);

        // Update position
        uint256 currentCollateral = collateralManager.getBalance(msg.sender);
        uint256 newDebt = debtManager.getDebt(msg.sender);
        positionManager.updatePosition(msg.sender, currentCollateral, newDebt);

        emit DebtDecreased(msg.sender, amount, newDebt);
    }

    function getTotalDebt(address user) external view override returns (uint256) {
        uint256 principalDebt = debtManager.getDebt(user);
        if (principalDebt == 0) return 0;

        // Get position to find last update time
        IPositionManager.Position memory pos = positionManager.getPosition(user);

        // Calculate accrued interest
        uint256 interest = interestCalculator.calculateInterest(
            user,
            principalDebt,
            pos.lastInterestUpdate
        );

        return principalDebt + interest;
    }

    function getSystemDebt() external view override returns (uint256) {
        return debtManager.getTotalDebt();
    }

    function calculateAccruedInterest(address user) external view override returns (uint256) {
        uint256 principalDebt = debtManager.getDebt(user);
        if (principalDebt == 0) return 0;

        IPositionManager.Position memory pos = positionManager.getPosition(user);
        return interestCalculator.calculateInterest(user, principalDebt, pos.lastInterestUpdate);
    }

    function getMaxMintable(address user) external view override returns (uint256) {
        uint256 collateral = collateralManager.getBalance(user);
        return debtManager.getMaxDebt(user, collateral, _config.minCollateralRatio);
    }

    // ============ Position Management ============

    function getPosition(address user) external view override returns (Position memory) {
        IPositionManager.Position memory pos = positionManager.getPosition(user);
        return Position({
            collateralAmount: pos.collateralAmount,
            debtAmount: pos.debtAmount,
            lastInterestUpdate: pos.lastInterestUpdate,
            accruedInterest: pos.accruedInterest,
            isActive: pos.isActive
        });
    }

    function getCollateralRatio(address user) external view override returns (uint256) {
        IPositionManager.Position memory pos = positionManager.getPosition(user);
        return positionManager.calculateCollateralRatio(pos.collateralAmount, pos.debtAmount);
    }

    function isPositionHealthy(address user) external view override returns (bool) {
        IPositionManager.Position memory pos = positionManager.getPosition(user);
        return positionManager.isPositionHealthy(
            pos.collateralAmount,
            pos.debtAmount,
            _config.minCollateralRatio
        );
    }

    function canLiquidate(address user) external view override returns (bool) {
        IPositionManager.Position memory pos = positionManager.getPosition(user);
        uint256 ratio = positionManager.calculateCollateralRatio(
            pos.collateralAmount,
            pos.debtAmount
        );
        return ratio < _config.liquidationThreshold && pos.debtAmount > 0;
    }

    // ============ Liquidation ============

    function liquidate(address user, uint256 debtAmount)
        external
        override
        returns (uint256 collateralSeized)
    {
        require(!_config.isPaused, "Paused");
        require(user != address(0), "Invalid user");
        require(debtAmount > 0, "Invalid amount");

        // Get user position
        IPositionManager.Position memory pos = positionManager.getPosition(user);
        require(pos.isActive, "No active position");

        // Check if position can be liquidated
        uint256 ratio = positionManager.calculateCollateralRatio(pos.collateralAmount, pos.debtAmount);
        require(
            liquidationEngine.canLiquidate(pos.collateralAmount, pos.debtAmount, ratio),
            "Position cannot be liquidated"
        );

        // Execute liquidation through LiquidationEngine
        ILiquidationEngine.LiquidationParams memory params = liquidationEngine.liquidate(
            user, msg.sender, debtAmount, pos.collateralAmount, pos.debtAmount
        );

        // Transfer debt tokens from liquidator
        IERC20(debtToken).transferFrom(msg.sender, address(this), params.debtToRepay);

        // Decrease user's debt
        debtManager.decreaseDebt(user, params.debtToRepay);

        // Transfer collateral to liquidator
        collateralManager.withdraw(user, params.collateralToSeize);
        IERC20(collateralManager.getCollateralToken()).transfer(msg.sender, params.liquidatorReward);

        // Protocol keeps the remaining penalty
        uint256 protocolFee = params.penalty - params.liquidatorReward;
        if (protocolFee > 0) {
            IERC20(collateralManager.getCollateralToken()).transfer(owner(), protocolFee);
        }

        // Update position
        uint256 newCollateral = pos.collateralAmount - params.collateralToSeize;
        uint256 newDebt = pos.debtAmount - params.debtToRepay;
        positionManager.updatePosition(user, newCollateral, newDebt);

        emit PositionLiquidated(
            user,
            msg.sender,
            params.collateralToSeize,
            params.debtToRepay,
            params.penalty
        );

        return params.collateralToSeize;
    }

    function calculateLiquidation(address user, uint256 debtAmount)
        external
        view
        override
        returns (uint256 collateralToSeize, uint256 penalty)
    {
        IPositionManager.Position memory pos = positionManager.getPosition(user);

        ILiquidationEngine.LiquidationParams memory params =
            liquidationEngine.calculateLiquidation(debtAmount, pos.collateralAmount, pos.debtAmount);

        return (params.collateralToSeize, params.penalty);
    }

    function getLiquidatablePositions() external view override returns (address[] memory) {
        address[] memory activePositions = positionManager.getActivePositions();
        uint256 count = 0;

        // First pass: count liquidatable positions
        for (uint256 i = 0; i < activePositions.length; i++) {
            IPositionManager.Position memory pos = positionManager.getPosition(activePositions[i]);
            if (pos.debtAmount > 0) {
                uint256 ratio = positionManager.calculateCollateralRatio(pos.collateralAmount, pos.debtAmount);
                if (liquidationEngine.canLiquidate(pos.collateralAmount, pos.debtAmount, ratio)) {
                    count++;
                }
            }
        }

        // Second pass: populate array
        address[] memory liquidatable = new address[](count);
        uint256 index = 0;
        for (uint256 i = 0; i < activePositions.length; i++) {
            IPositionManager.Position memory pos = positionManager.getPosition(activePositions[i]);
            if (pos.debtAmount > 0) {
                uint256 ratio = positionManager.calculateCollateralRatio(pos.collateralAmount, pos.debtAmount);
                if (liquidationEngine.canLiquidate(pos.collateralAmount, pos.debtAmount, ratio)) {
                    liquidatable[index++] = activePositions[i];
                }
            }
        }

        return liquidatable;
    }

    // ============ Configuration ============

    function getVaultConfig() external view override returns (VaultConfig memory) {
        return _config;
    }

    function setMinCollateralRatio(uint256 newRatio) external override onlyOwner {
        uint256 oldRatio = _config.minCollateralRatio;
        _config.minCollateralRatio = newRatio;
        emit ConfigUpdated("minCollateralRatio", oldRatio, newRatio);
    }

    function setLiquidationThreshold(uint256 newThreshold) external override onlyOwner {
        uint256 oldThreshold = _config.liquidationThreshold;
        _config.liquidationThreshold = newThreshold;
        emit ConfigUpdated("liquidationThreshold", oldThreshold, newThreshold);
    }

    function setStabilityFee(uint256 newFee) external override onlyOwner {
        uint256 oldFee = _config.stabilityFee;
        _config.stabilityFee = newFee;
        emit ConfigUpdated("stabilityFee", oldFee, newFee);
    }

    function setDebtCeiling(uint256 newCeiling) external override onlyOwner {
        uint256 oldCeiling = _config.debtCeiling;
        _config.debtCeiling = newCeiling;
        emit ConfigUpdated("debtCeiling", oldCeiling, newCeiling);
    }

    // ============ Statistics ============

    function getVaultStats()
        external
        view
        override
        returns (uint256 totalCollateral, uint256 totalDebt, uint256 avgRatio, uint256 utilization)
    {
        totalCollateral = collateralManager.getTotalCollateral();
        totalDebt = 0; // Simplified
        avgRatio = 0; // Simplified
        utilization = 0; // Simplified
    }

    function getActivePositionCount() external view override returns (uint256) {
        return positionManager.getActivePositionCount();
    }

    function getCollateralToken() external view override returns (address) {
        return collateralManager.getCollateralToken();
    }

    function getDebtToken() external view override returns (address) {
        return debtToken;
    }

    // ============ Pause ============

    function pause() external override(IModule, BaseModule) onlyOwner {
        _config.isPaused = true;
    }

    function unpause() external override(IModule, BaseModule) onlyOwner {
        _config.isPaused = false;
    }

    // ============ IUpgradeable Implementation ============
    // (Reusing implementation from V2 - omitted for brevity)

    function upgradeTo(address newImplementation) external override onlyOwner {
        _authorizeUpgrade(newImplementation);
        _upgradeToAndCallUUPS(newImplementation, new bytes(0), false);
        _upgradeHistory.push(UpgradeRecord({
            implementation: newImplementation,
            timestamp: block.timestamp,
            version: MODULE_VERSION
        }));
        emit Upgraded(address(this), newImplementation, MODULE_VERSION);
    }

    function upgradeToAndCall(address newImplementation, bytes memory data)
        external
        payable
        override
        onlyOwner
    {
        _authorizeUpgrade(newImplementation);
        _upgradeToAndCallUUPS(newImplementation, data, true);
        _upgradeHistory.push(UpgradeRecord({
            implementation: newImplementation,
            timestamp: block.timestamp,
            version: MODULE_VERSION
        }));
        emit Upgraded(address(this), newImplementation, MODULE_VERSION);
    }

    function getImplementation() external view override returns (address) {
        return _getImplementation();
    }

    function getImplementationVersion() external pure override returns (string memory) {
        return MODULE_VERSION;
    }

    function canUpgrade(address account) external view override returns (bool) {
        return _authorizedUpgraders[account] || account == owner();
    }

    function authorizeUpgrader(address account) external override onlyOwner {
        _authorizedUpgraders[account] = true;
        emit UpgradeAuthorized(account, msg.sender);
    }

    function revokeUpgradeAuthorization(address account) external override onlyOwner {
        _authorizedUpgraders[account] = false;
        emit UpgradeAuthorizationRevoked(account);
    }

    function validateUpgrade(address newImplementation)
        external
        view
        override
        returns (bool valid, string memory reason)
    {
        if (newImplementation == address(0)) return (false, "Zero address");
        if (newImplementation.code.length == 0) return (false, "Not a contract");
        if (_upgradesPaused) return (false, "Upgrades paused");
        return (true, "");
    }

    function beforeUpgrade(address) external override onlyOwner {}
    function afterUpgrade(address) external override onlyOwner {}

    function getUpgradeHistory()
        external
        view
        override
        returns (address[] memory, uint256[] memory, string[] memory)
    {
        uint256 length = _upgradeHistory.length;
        address[] memory implementations = new address[](length);
        uint256[] memory timestamps = new uint256[](length);
        string[] memory versions = new string[](length);

        for (uint256 i = 0; i < length; i++) {
            implementations[i] = _upgradeHistory[i].implementation;
            timestamps[i] = _upgradeHistory[i].timestamp;
            versions[i] = _upgradeHistory[i].version;
        }

        return (implementations, timestamps, versions);
    }

    function pauseUpgrades() external override onlyOwner {
        _upgradesPaused = true;
    }

    function resumeUpgrades() external override onlyOwner {
        _upgradesPaused = false;
    }

    function upgradesPaused() external view override returns (bool) {
        return _upgradesPaused;
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {
        require(!_upgradesPaused, "Upgrades paused");
        require(newImplementation != address(0), "Zero address");
        require(
            _authorizedUpgraders[msg.sender] || msg.sender == owner(),
            "Not authorized"
        );
    }

    function getModuleInfo() external view override returns (ModuleInfo memory) {
        return ModuleInfo({
            moduleId: MODULE_ID,
            name: MODULE_NAME,
            version: MODULE_VERSION,
            implementation: address(this),
            state: _config.isPaused ? ModuleState.PAUSED : ModuleState.ACTIVE,
            installedAt: _upgradeHistory.length > 0
                ? _upgradeHistory[0].timestamp
                : block.timestamp,
            lastUpdated: _upgradeHistory.length > 0
                ? _upgradeHistory[_upgradeHistory.length - 1].timestamp
                : block.timestamp
        });
    }
}
