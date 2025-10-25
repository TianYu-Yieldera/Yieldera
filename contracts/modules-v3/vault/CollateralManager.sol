// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/modules/vault/ICollateralManager.sol";
import "../../storage/VaultModuleStorage.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title CollateralManager
 * @notice Manages collateral deposits and withdrawals
 * @dev Uses Diamond Storage for state management
 */
contract CollateralManager is ICollateralManager, Ownable {
    using SafeERC20 for IERC20;
    using VaultModuleStorage for VaultModuleStorage.VaultData;

    // ============ State Variables ============

    IERC20 public immutable collateralToken;
    address public vaultModule;  // Main coordinator contract

    // Storage for collateral balances (using Diamond Storage pattern)
    bytes32 private constant COLLATERAL_STORAGE_POSITION =
        keccak256("collateral.manager.storage");

    struct CollateralStorage {
        mapping(address => uint256) balances;
        uint256 totalCollateral;
    }

    // ============ Modifiers ============

    modifier onlyVaultModule() {
        require(msg.sender == vaultModule, "Only vault module");
        _;
    }

    // ============ Constructor ============

    constructor(address _collateralToken) {
        require(_collateralToken != address(0), "Invalid token");
        collateralToken = IERC20(_collateralToken);
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

    // ============ Internal Storage Functions ============

    function _getStorage() private pure returns (CollateralStorage storage cs) {
        bytes32 position = COLLATERAL_STORAGE_POSITION;
        assembly {
            cs.slot := position
        }
    }

    // ============ ICollateralManager Implementation ============

    /**
     * @notice Deposit collateral for a user
     * @param user User address
     * @param amount Amount to deposit
     */
    function deposit(address user, uint256 amount) external override onlyVaultModule {
        require(user != address(0), "Invalid user");
        require(amount > 0, "Amount must be > 0");

        // Transfer tokens from user
        collateralToken.safeTransferFrom(user, address(this), amount);

        // Update storage
        CollateralStorage storage cs = _getStorage();
        cs.balances[user] += amount;
        cs.totalCollateral += amount;

        uint256 newTotal = cs.balances[user];

        emit CollateralDeposited(user, amount, newTotal);
    }

    /**
     * @notice Withdraw collateral for a user
     * @param user User address
     * @param amount Amount to withdraw
     */
    function withdraw(address user, uint256 amount) external override onlyVaultModule {
        require(user != address(0), "Invalid user");
        require(amount > 0, "Amount must be > 0");

        CollateralStorage storage cs = _getStorage();
        require(cs.balances[user] >= amount, "Insufficient balance");

        // Update storage
        cs.balances[user] -= amount;
        cs.totalCollateral -= amount;

        // Transfer tokens to user
        collateralToken.safeTransfer(user, amount);

        uint256 remaining = cs.balances[user];

        emit CollateralWithdrawn(user, amount, remaining);
    }

    /**
     * @notice Get collateral balance for a user
     * @param user User address
     * @return Collateral amount
     */
    function getBalance(address user) external view override returns (uint256) {
        CollateralStorage storage cs = _getStorage();
        return cs.balances[user];
    }

    /**
     * @notice Get total collateral in the system
     * @return Total collateral amount
     */
    function getTotalCollateral() external view override returns (uint256) {
        CollateralStorage storage cs = _getStorage();
        return cs.totalCollateral;
    }

    /**
     * @notice Check if user has sufficient collateral
     * @param user User address
     * @param amount Required amount
     * @return True if user has sufficient collateral
     */
    function hasSufficientCollateral(address user, uint256 amount)
        external
        view
        override
        returns (bool)
    {
        CollateralStorage storage cs = _getStorage();
        return cs.balances[user] >= amount;
    }

    /**
     * @notice Get collateral token address
     * @return Collateral token address
     */
    function getCollateralToken() external view override returns (address) {
        return address(collateralToken);
    }
}
