// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "../interfaces/IPriceOracle.sol";
import "../aggregator/L2StateAggregator.sol";

/**
 * @title IntegratedVault
 * @notice Collateral vault with lending capabilities and price oracle integration
 * @dev Supports deposit, borrow, repay, and withdraw operations with real-time pricing
 */
contract IntegratedVault is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // State variables
    IERC20 public immutable collateralToken;
    IERC20 public immutable lusdToken;
    IPriceOracle public priceOracle;
    L2StateAggregator public stateAggregator;

    uint256 public constant COLLATERAL_RATIO = 150; // 150% collateralization
    uint256 public constant INTEREST_RATE = 3; // 3% annual interest
    uint256 public constant PRECISION = 100;
    uint256 public constant PRICE_PRECISION = 1e8; // Chainlink uses 8 decimals

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
     * @param _priceOracle Address of price oracle
     * @param _stateAggregator Address of L2 state aggregator
     * @param initialOwner Address of contract owner
     */
    constructor(
        address _collateralToken,
        address _lusdToken,
        address _priceOracle,
        address _stateAggregator,
        address initialOwner
    ) Ownable(initialOwner) {
        require(_collateralToken != address(0), "Invalid collateral token");
        require(_lusdToken != address(0), "Invalid LUSD token");
        require(_priceOracle != address(0), "Invalid price oracle");

        collateralToken = IERC20(_collateralToken);
        lusdToken = IERC20(_lusdToken);
        priceOracle = IPriceOracle(_priceOracle);
        stateAggregator = L2StateAggregator(_stateAggregator);
    }

    /**
     * @notice Update price oracle address
     * @param _newOracle Address of new price oracle
     */
    function updatePriceOracle(address _newOracle) external onlyOwner {
        require(_newOracle != address(0), "Invalid oracle address");
        priceOracle = IPriceOracle(_newOracle);
    }

    /**
     * @notice Update state aggregator address
     * @param _newAggregator Address of new state aggregator
     */
    function updateStateAggregator(address _newAggregator) external onlyOwner {
        stateAggregator = L2StateAggregator(_newAggregator);
    }

    /**
     * @notice Update system state in the aggregator
     */
    function _updateAggregator() internal {
        if (address(stateAggregator) != address(0)) {
            stateAggregator.updateSystemState(
                totalCollateral,
                totalDebt,
                activePositions,
                0 // orderCount (not applicable for vault)
            );
        }
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
     * @notice Calculate health factor for a position with price oracle
     * @param collateralAmount Amount of collateral tokens
     * @param debtAmount Amount of debt in LUSD
     * @return Health factor (scaled by 100)
     */
    function _calculateHealthFactor(
        uint256 collateralAmount,
        uint256 debtAmount
    ) internal view returns (uint256) {
        if (debtAmount == 0) return type(uint256).max;

        // Get collateral price from oracle (8 decimals)
        uint256 collateralPrice = priceOracle.getAssetPrice(address(collateralToken));

        // Calculate collateral value in USD
        // collateralAmount * price / 1e8 (to normalize price decimals)
        uint256 collateralValueUSD = (collateralAmount * collateralPrice) / PRICE_PRECISION;

        // Health factor = (collateral value * 100) / debt
        return (collateralValueUSD * PRECISION) / debtAmount;
    }

    /**
     * @notice Get collateral value in USD for a user
     * @param user Address of the user
     * @return valueUSD Collateral value in USD
     */
    function getCollateralValue(address user) external view returns (uint256 valueUSD) {
        Position memory pos = positions[user];
        if (pos.collateral == 0) return 0;

        uint256 collateralPrice = priceOracle.getAssetPrice(address(collateralToken));
        return (pos.collateral * collateralPrice) / PRICE_PRECISION;
    }

    /**
     * @notice Deposit collateral into the vault
     * @param amount Amount of collateral to deposit
     */
    function deposit(uint256 amount) external nonReentrant {
        require(amount > 0, "Amount must be greater than zero");

        Position storage pos = positions[msg.sender];

        if (pos.collateral == 0) {
            activePositions++;
        }

        collateralToken.safeTransferFrom(msg.sender, address(this), amount);

        pos.collateral += amount;
        pos.lastUpdate = block.timestamp;
        totalCollateral += amount;

        emit Deposited(msg.sender, amount);

        _updateAggregator();
    }

    /**
     * @notice Borrow LUSD against collateral
     * @param amount Amount of LUSD to borrow
     */
    function borrow(uint256 amount) external nonReentrant {
        require(amount > 0, "Amount must be greater than zero");

        Position storage pos = positions[msg.sender];
        require(pos.collateral > 0, "No collateral deposited");

        uint256 newDebt = pos.debt + amount;
        uint256 healthFactor = _calculateHealthFactor(pos.collateral, newDebt);

        require(
            healthFactor >= COLLATERAL_RATIO,
            "Insufficient collateral ratio"
        );

        pos.debt = newDebt;
        pos.lastUpdate = block.timestamp;
        totalDebt += amount;

        lusdToken.safeTransfer(msg.sender, amount);

        emit Borrowed(msg.sender, amount);

        _updateAggregator();
    }

    /**
     * @notice Repay LUSD debt
     * @param amount Amount of LUSD to repay
     */
    function repay(uint256 amount) external nonReentrant {
        require(amount > 0, "Amount must be greater than zero");

        Position storage pos = positions[msg.sender];
        require(pos.debt > 0, "No debt to repay");
        require(amount <= pos.debt, "Amount exceeds debt");

        lusdToken.safeTransferFrom(msg.sender, address(this), amount);

        pos.debt -= amount;
        pos.lastUpdate = block.timestamp;
        totalDebt -= amount;

        emit Repaid(msg.sender, amount);

        _updateAggregator();
    }

    /**
     * @notice Withdraw collateral from the vault
     * @param amount Amount of collateral to withdraw
     */
    function withdraw(uint256 amount) external nonReentrant {
        require(amount > 0, "Amount must be greater than zero");

        Position storage pos = positions[msg.sender];
        require(pos.collateral > 0, "No collateral to withdraw");
        require(amount <= pos.collateral, "Amount exceeds collateral");

        uint256 remainingCollateral = pos.collateral - amount;

        // If user has debt, check that withdrawal maintains healthy collateral ratio
        if (pos.debt > 0) {
            uint256 newHealthFactor = _calculateHealthFactor(remainingCollateral, pos.debt);
            require(
                newHealthFactor >= COLLATERAL_RATIO,
                "Withdrawal would break collateral ratio"
            );
        }

        pos.collateral -= amount;
        pos.lastUpdate = block.timestamp;
        totalCollateral -= amount;

        // If position is fully withdrawn, decrease active positions count
        if (pos.collateral == 0 && pos.debt == 0) {
            activePositions--;
        }

        collateralToken.safeTransfer(msg.sender, amount);

        emit Withdrawn(msg.sender, amount);

        _updateAggregator();
    }
}
