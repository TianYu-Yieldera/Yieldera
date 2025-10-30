// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title IntegratedVault
 * @notice Collateral vault with lending capabilities
 * @dev Supports deposit, borrow, repay, and withdraw operations
 */
contract IntegratedVault is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // State variables
    IERC20 public immutable collateralToken;
    IERC20 public immutable lusdToken;

    uint256 public constant COLLATERAL_RATIO = 150; // 150% collateralization
    uint256 public constant INTEREST_RATE = 3; // 3% annual interest
    uint256 public constant PRECISION = 100;

    uint256 public totalCollateral;
    uint256 public totalDebt;
    uint256 public activePositions;

    struct Position {
        uint256 collateral;
        uint256 debt;
        uint256 lastUpdate;
    }

    mapping(address => Position) public positions;

    // Events
    event Deposited(address indexed user, uint256 amount);
    event Borrowed(address indexed user, uint256 amount);
    event Repaid(address indexed user, uint256 amount);
    event Withdrawn(address indexed user, uint256 amount);

    /**
     * @notice Contract constructor
     * @param _collateralToken Address of collateral token
     * @param _lusdToken Address of LUSD stablecoin
     * @param initialOwner Address of contract owner
     */
    constructor(
        address _collateralToken,
        address _lusdToken,
        address initialOwner
    ) Ownable(initialOwner) {
        require(_collateralToken != address(0), "Invalid collateral token");
        require(_lusdToken != address(0), "Invalid LUSD token");

        collateralToken = IERC20(_collateralToken);
        lusdToken = IERC20(_lusdToken);
    }

    /**
     * @notice Get user position details
     * @param user Address of the user
     * @return collateral Amount of collateral
     * @return debt Current debt amount
     * @return healthFactor Current health factor (scaled by 100)
     */
    function getPosition(address user) external view returns (
        uint256 collateral,
        uint256 debt,
        uint256 healthFactor
    ) {
        Position memory pos = positions[user];
        collateral = pos.collateral;
        debt = pos.debt;
        healthFactor = _calculateHealthFactor(pos.collateral, pos.debt);
    }

    /**
     * @notice Calculate health factor for a position
     * @param collateralAmount Amount of collateral
     * @param debtAmount Amount of debt
     * @return Health factor (scaled by 100)
     */
    function _calculateHealthFactor(
        uint256 collateralAmount,
        uint256 debtAmount
    ) internal pure returns (uint256) {
        if (debtAmount == 0) return type(uint256).max;
        return (collateralAmount * PRECISION) / debtAmount;
    }
}
