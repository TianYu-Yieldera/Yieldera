// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/core/IVaultModule.sol";
import "../../interfaces/IUpgradeable.sol";
import "../../plugins/core/BaseModule.sol";
import "../../storage/VaultModuleStorage.sol";
import "../../core/CollateralVault.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/**
 * @title VaultModuleV2
 * @notice Upgradeable vault module using UUPS pattern and Diamond Storage
 * @dev Implements IVaultModule interface with upgrade capabilities
 *
 * Upgrade Features:
 * - UUPS upgradeable
 * - Diamond Storage (no storage collisions)
 * - Backward compatible with VaultModule
 * - Version tracking
 */
contract VaultModuleV2 is
    Initializable,
    IVaultModule,
    IUpgradeable,
    BaseModule,
    OwnableUpgradeable,
    UUPSUpgradeable
{
    using VaultModuleStorage for VaultModuleStorage.VaultData;

    // ============ Constants ============

    bytes32 public constant MODULE_ID = keccak256("VAULT_MODULE");
    string public constant MODULE_NAME = "VaultModule";
    string public constant MODULE_VERSION = "2.0.0";

    // ============ Upgrade History ============

    struct UpgradeRecord {
        address implementation;
        uint256 timestamp;
        string version;
    }

    UpgradeRecord[] private _upgradeHistory;
    mapping(address => bool) private _authorizedUpgraders;
    bool private _upgradesPaused;

    // ============ Events ============

    event ConfigUpdated(string param, uint256 oldValue, uint256 newValue);

    // ============ Constructor & Initializer ============

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initialize the upgradeable contract
     * @param _vault CollateralVault address
     * @param _collateralToken Collateral token address
     * @param _debtToken Debt token address
     */
    function initialize(
        address _vault,
        address _collateralToken,
        address _debtToken
    ) public initializer {
        __Ownable_init();
        __UUPSUpgradeable_init();

        // Initialize Diamond Storage
        VaultModuleStorage.initialize(_vault, _collateralToken, _debtToken);

        // Initialize module config
        CollateralVault vault = CollateralVault(_vault);
        VaultModuleStorage.setConfig(
            IVaultModule.VaultConfig({
                minCollateralRatio: vault.COLLATERAL_RATIO(),
                liquidationThreshold: vault.LIQUIDATION_THRESHOLD(),
                liquidationPenalty: 10,
                stabilityFee: vault.STABILITY_FEE(),
                debtCeiling: type(uint256).max,
                minDebtAmount: 0,
                isPaused: false
            })
        );

        // Record initial deployment
        _upgradeHistory.push(
            UpgradeRecord({
                implementation: address(this),
                timestamp: block.timestamp,
                version: MODULE_VERSION
            })
        );

        // Authorize owner as upgrader
        _authorizedUpgraders[msg.sender] = true;
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
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        return !config.isPaused;
    }

    function initialize(bytes calldata) external override {
        // This is handled by the initialize(address,address,address) function
        revert("Use initialize(address,address,address)");
    }

    function healthCheck()
        external
        view
        override(IModule, BaseModule)
        returns (bool healthy, string memory message)
    {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        CollateralVault vaultContract = CollateralVault(vault);

        // Check vault health
        (uint256 totalCollateral, uint256 totalDebt, ) = vaultContract.getVaultStats();

        if (totalDebt > 0 && totalCollateral == 0) {
            return (false, "Vault has debt but no collateral");
        }

        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        uint256 ratio = totalDebt > 0 ? (totalCollateral * 100) / totalDebt : type(uint256).max;

        if (ratio < config.minCollateralRatio) {
            return (false, "System undercollateralized");
        }

        return (true, "Vault healthy");
    }

    // ============ Collateral Management ============

    function depositCollateral(uint256 amount) external override {
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        require(!config.isPaused, "Paused");
        require(amount > 0, "Amount must be > 0");

        (address vault, , ) = VaultModuleStorage.getAddresses();
        CollateralVault(vault).depositCollateral(amount);

        VaultModuleStorage.addPosition(msg.sender);

        emit CollateralDeposited(
            msg.sender,
            amount,
            CollateralVault(vault).collateralDeposited(msg.sender)
        );
    }

    function withdrawCollateral(uint256 amount) external override {
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        require(!config.isPaused, "Paused");
        require(amount > 0, "Amount must be > 0");

        (address vault, , ) = VaultModuleStorage.getAddresses();
        CollateralVault(vault).withdrawCollateral(amount);

        emit CollateralWithdrawn(
            msg.sender,
            amount,
            CollateralVault(vault).collateralDeposited(msg.sender)
        );
    }

    function getCollateralBalance(address user) external view override returns (uint256) {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        return CollateralVault(vault).collateralDeposited(user);
    }

    function getTotalCollateral() external view override returns (uint256) {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        return CollateralVault(vault).totalCollateral();
    }

    // ============ Debt Management ============

    function increaseDebt(uint256 amount) external override {
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        require(!config.isPaused, "Paused");
        require(amount > 0, "Amount must be > 0");

        (address vault, , ) = VaultModuleStorage.getAddresses();
        CollateralVault vaultContract = CollateralVault(vault);

        uint256 newTotalDebt = vaultContract.totalDebt() + amount;
        require(newTotalDebt <= config.debtCeiling, "Exceeds debt ceiling");

        vaultContract.increaseDebt(msg.sender, amount);

        VaultModuleStorage.addPosition(msg.sender);

        emit DebtIncreased(msg.sender, amount, vaultContract.debtAmount(msg.sender));
    }

    function decreaseDebt(uint256 amount) external override {
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        require(!config.isPaused, "Paused");
        require(amount > 0, "Amount must be > 0");

        (address vault, , ) = VaultModuleStorage.getAddresses();
        CollateralVault vaultContract = CollateralVault(vault);

        vaultContract.decreaseDebt(msg.sender, amount);

        emit DebtDecreased(msg.sender, amount, vaultContract.debtAmount(msg.sender));
    }

    function getTotalDebt(address user) external view override returns (uint256) {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        return CollateralVault(vault).getTotalDebt(user);
    }

    function getSystemDebt() external view override returns (uint256) {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        return CollateralVault(vault).totalDebt();
    }

    function calculateAccruedInterest(address user)
        external
        view
        override
        returns (uint256)
    {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        return CollateralVault(vault).accruedInterest(user);
    }

    function getMaxMintable(address user) external view override returns (uint256) {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        return CollateralVault(vault).getMaxMintable(user);
    }

    // ============ Position Management ============

    function getPosition(address user) external view override returns (Position memory) {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        CollateralVault vaultContract = CollateralVault(vault);

        uint256 collateral = vaultContract.collateralDeposited(user);
        uint256 debt = vaultContract.debtAmount(user);
        uint256 interest = vaultContract.accruedInterest(user);
        uint256 lastUpdate = vaultContract.lastInterestUpdate(user);

        return
            Position({
                collateralAmount: collateral,
                debtAmount: debt,
                lastInterestUpdate: lastUpdate,
                accruedInterest: interest,
                isActive: debt > 0 || collateral > 0
            });
    }

    function getCollateralRatio(address user) external view override returns (uint256) {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        return CollateralVault(vault).getCollateralRatio(user);
    }

    function isPositionHealthy(address user) external view override returns (bool) {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        return CollateralVault(vault).isPositionHealthy(user);
    }

    function canLiquidate(address user) external view override returns (bool) {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        return CollateralVault(vault).canLiquidate(user);
    }

    // ============ Liquidation ============

    function liquidate(address user, uint256 debtToCover)
        external
        override
        returns (uint256 collateralSeized)
    {
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        require(!config.isPaused, "Paused");

        (address vault, , ) = VaultModuleStorage.getAddresses();
        collateralSeized = CollateralVault(vault).liquidate(user, debtToCover);

        emit PositionLiquidated(user, msg.sender, collateralSeized, debtToCover);

        return collateralSeized;
    }

    function calculateLiquidation(address user, uint256 debtToCover)
        external
        view
        override
        returns (uint256 collateralToSeize, uint256 liquidationPenalty)
    {
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        (address vault, , ) = VaultModuleStorage.getAddresses();

        collateralToSeize = (debtToCover * 110) / 100;

        uint256 userCollateral = CollateralVault(vault).collateralDeposited(user);
        if (collateralToSeize > userCollateral) {
            collateralToSeize = userCollateral;
        }

        liquidationPenalty = (debtToCover * config.liquidationPenalty) / 100;

        return (collateralToSeize, liquidationPenalty);
    }

    function getLiquidatablePositions()
        external
        view
        override
        returns (address[] memory users)
    {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        CollateralVault vaultContract = CollateralVault(vault);

        address[] storage activePositions = VaultModuleStorage.getActivePositions();
        uint256 count = 0;

        for (uint256 i = 0; i < activePositions.length; i++) {
            if (vaultContract.canLiquidate(activePositions[i])) {
                count++;
            }
        }

        users = new address[](count);
        uint256 index = 0;

        for (uint256 i = 0; i < activePositions.length; i++) {
            if (vaultContract.canLiquidate(activePositions[i])) {
                users[index] = activePositions[i];
                index++;
            }
        }

        return users;
    }

    // ============ Configuration ============

    function getVaultConfig() external view override returns (VaultConfig memory) {
        return VaultModuleStorage.getConfig();
    }

    function setMinCollateralRatio(uint256 newRatio) external override onlyOwner {
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        uint256 oldRatio = config.minCollateralRatio;
        config.minCollateralRatio = newRatio;
        emit ConfigUpdated("minCollateralRatio", oldRatio, newRatio);
    }

    function setLiquidationThreshold(uint256 newThreshold) external override onlyOwner {
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        uint256 oldThreshold = config.liquidationThreshold;
        config.liquidationThreshold = newThreshold;
        emit ConfigUpdated("liquidationThreshold", oldThreshold, newThreshold);
    }

    function setStabilityFee(uint256 newFee) external override onlyOwner {
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        uint256 oldFee = config.stabilityFee;
        config.stabilityFee = newFee;
        emit ConfigUpdated("stabilityFee", oldFee, newFee);
    }

    function setDebtCeiling(uint256 newCeiling) external override onlyOwner {
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        uint256 oldCeiling = config.debtCeiling;
        config.debtCeiling = newCeiling;
        emit ConfigUpdated("debtCeiling", oldCeiling, newCeiling);
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
        (address vault, , ) = VaultModuleStorage.getAddresses();
        CollateralVault vaultContract = CollateralVault(vault);

        (totalCollateral, totalDebt, averageCollateralRatio) = vaultContract
            .getVaultStats();

        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        utilizationRate = config.debtCeiling > 0
            ? (totalDebt * 10000) / config.debtCeiling
            : 0;

        return (totalCollateral, totalDebt, averageCollateralRatio, utilizationRate);
    }

    function getActivePositionCount() external view override returns (uint256) {
        (address vault, , ) = VaultModuleStorage.getAddresses();
        CollateralVault vaultContract = CollateralVault(vault);

        address[] storage activePositions = VaultModuleStorage.getActivePositions();
        uint256 count = 0;

        for (uint256 i = 0; i < activePositions.length; i++) {
            address user = activePositions[i];
            if (
                vaultContract.collateralDeposited(user) > 0 ||
                vaultContract.debtAmount(user) > 0
            ) {
                count++;
            }
        }

        return count;
    }

    function getCollateralToken() external view override returns (address) {
        (, address collateralToken, ) = VaultModuleStorage.getAddresses();
        return collateralToken;
    }

    function getDebtToken() external view override returns (address) {
        (, , address debtToken) = VaultModuleStorage.getAddresses();
        return debtToken;
    }

    // ============ Pause Functions ============

    function pause() external override(IModule, BaseModule) onlyOwner {
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        config.isPaused = true;
    }

    function unpause() external override(IModule, BaseModule) onlyOwner {
        VaultModuleStorage.VaultConfig storage config = VaultModuleStorage.getConfig();
        config.isPaused = false;
    }

    // ============ IUpgradeable Implementation ============

    function upgradeTo(address newImplementation) external override onlyOwner {
        _authorizeUpgrade(newImplementation);
        _upgradeToAndCallUUPS(newImplementation, new bytes(0), false);

        _upgradeHistory.push(
            UpgradeRecord({
                implementation: newImplementation,
                timestamp: block.timestamp,
                version: MODULE_VERSION
            })
        );

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

        _upgradeHistory.push(
            UpgradeRecord({
                implementation: newImplementation,
                timestamp: block.timestamp,
                version: MODULE_VERSION
            })
        );

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
        if (newImplementation == address(0)) {
            return (false, "Zero address");
        }

        if (newImplementation.code.length == 0) {
            return (false, "Not a contract");
        }

        if (_upgradesPaused) {
            return (false, "Upgrades paused");
        }

        return (true, "");
    }

    function beforeUpgrade(address) external override onlyOwner {
        // Pre-upgrade logic if needed
    }

    function afterUpgrade(address) external override onlyOwner {
        // Post-upgrade logic if needed
    }

    function getUpgradeHistory()
        external
        view
        override
        returns (
            address[] memory implementations,
            uint256[] memory timestamps,
            string[] memory versions
        )
    {
        uint256 length = _upgradeHistory.length;
        implementations = new address[](length);
        timestamps = new uint256[](length);
        versions = new string[](length);

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

    // ============ UUPS Override ============

    function _authorizeUpgrade(address newImplementation)
        internal
        override
        onlyOwner
    {
        require(!_upgradesPaused, "Upgrades paused");
        require(newImplementation != address(0), "Zero address");
        require(
            _authorizedUpgraders[msg.sender] || msg.sender == owner(),
            "Not authorized"
        );
    }

    /**
     * @dev See {IERC165-supportsInterface}.
     */
    function getModuleInfo() external view override returns (ModuleInfo memory) {
        return
            ModuleInfo({
                moduleId: MODULE_ID,
                name: MODULE_NAME,
                version: MODULE_VERSION,
                implementation: address(this),
                state: ModuleState.ACTIVE,
                installedAt: _upgradeHistory.length > 0
                    ? _upgradeHistory[0].timestamp
                    : block.timestamp,
                lastUpdated: _upgradeHistory.length > 0
                    ? _upgradeHistory[_upgradeHistory.length - 1].timestamp
                    : block.timestamp
            });
    }
}
