// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "../interfaces/IPriceOracle.sol";
import "../aggregator/L2StateAggregator.sol";

/**
 * @title IntegratedVault
 * @notice Collateral vault with lending capabilities and price oracle integration
 * @dev Supports deposit, borrow, repay, and withdraw operations with real-time pricing
 * @dev Includes emergency pause functionality and liquidation mechanism
 */
contract IntegratedVault is Ownable, ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;

    // State variables
    IERC20 public immutable collateralToken;
    IERC20 public immutable lusdToken;
    IPriceOracle public priceOracle;
    L2StateAggregator public stateAggregator;

    uint256 public constant COLLATERAL_RATIO = 150; // 150% collateralization
    uint256 public constant LIQUIDATION_THRESHOLD = 130; // 130% - liquidation trigger
    uint256 public constant LIQUIDATION_BONUS = 5; // 5% bonus for liquidators
    uint256 public constant INTEREST_RATE = 3; // 3% annual interest (in percentage)
    uint256 public constant PRECISION = 100;
    uint256 public constant PRICE_PRECISION = 1e8; // Chainlink uses 8 decimals
    uint256 public constant SECONDS_PER_YEAR = 365 days;

    uint256 public totalCollateral;
    uint256 public totalDebt;
    uint256 public activePositions;
    uint256 public lastInterestUpdate; // Timestamp of last global interest update

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
    event Liquidated(
        address indexed user,
        address indexed liquidator,
        uint256 debtRepaid,
        uint256 collateralSeized
    );
    event InterestAccrued(address indexed user, uint256 interestAmount);

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
        lastInterestUpdate = block.timestamp;
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
     * @notice Get maximum amount a user can borrow based on their collateral
     * @param user Address of the user
     * @return maxBorrow Maximum LUSD amount that can be borrowed
     */
    function getMaxBorrowAmount(address user) external view returns (uint256 maxBorrow) {
        Position memory pos = positions[user];
        if (pos.collateral == 0) return 0;

        // Get collateral value in USD
        uint256 collateralPrice = priceOracle.getAssetPrice(address(collateralToken));
        uint256 collateralValueUSD = (pos.collateral * collateralPrice) / PRICE_PRECISION;

        // Max borrow = (collateral value * 100) / collateral ratio - current debt
        // This ensures health factor stays at exactly COLLATERAL_RATIO
        uint256 maxTotalDebt = (collateralValueUSD * PRECISION) / COLLATERAL_RATIO;

        if (maxTotalDebt <= pos.debt) {
            return 0;
        }

        return maxTotalDebt - pos.debt;
    }

    /**
     * @notice Get maximum amount a user can withdraw while maintaining collateral ratio
     * @param user Address of the user
     * @return maxWithdraw Maximum collateral amount that can be withdrawn
     */
    function getMaxWithdrawAmount(address user) external view returns (uint256 maxWithdraw) {
        Position memory pos = positions[user];
        if (pos.collateral == 0) return 0;

        // If no debt, can withdraw all collateral
        if (pos.debt == 0) {
            return pos.collateral;
        }

        // Calculate minimum collateral needed to maintain COLLATERAL_RATIO
        // minCollateralValue = (debt * COLLATERAL_RATIO) / 100
        uint256 minCollateralValueUSD = (pos.debt * COLLATERAL_RATIO) / PRECISION;

        // Convert USD value to collateral tokens
        uint256 collateralPrice = priceOracle.getAssetPrice(address(collateralToken));
        uint256 minCollateral = (minCollateralValueUSD * PRICE_PRECISION) / collateralPrice;

        if (pos.collateral <= minCollateral) {
            return 0;
        }

        return pos.collateral - minCollateral;
    }

    /**
     * @notice Get user's current health factor
     * @param user Address of the user
     * @return healthFactor Current health factor (scaled by 100)
     * @dev Returns type(uint256).max if user has no debt
     */
    function getUserHealthFactor(address user) external view returns (uint256 healthFactor) {
        Position memory pos = positions[user];
        return _calculateHealthFactor(pos.collateral, pos.debt);
    }

    /**
     * @notice Deposit collateral into the vault
     * @param amount Amount of collateral to deposit
     */
    function deposit(uint256 amount) external nonReentrant whenNotPaused {
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
    function borrow(uint256 amount) external nonReentrant whenNotPaused {
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
    function withdraw(uint256 amount) external nonReentrant whenNotPaused {
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

    /**
     * @notice Accrue interest for a user's debt
     * @param user Address of the user
     * @return interestAccrued Amount of interest accrued
     * @dev Interest is calculated as: debt * rate * time / (100 * SECONDS_PER_YEAR)
     */
    function accrueInterest(address user) public returns (uint256 interestAccrued) {
        Position storage pos = positions[user];

        if (pos.debt == 0) {
            return 0;
        }

        uint256 timeElapsed = block.timestamp - pos.lastUpdate;
        if (timeElapsed == 0) {
            return 0;
        }

        // Calculate interest: debt * rate * time / (100 * SECONDS_PER_YEAR)
        // Interest = principal * rate% * time / year
        interestAccrued = (pos.debt * INTEREST_RATE * timeElapsed) / (PRECISION * SECONDS_PER_YEAR);

        if (interestAccrued > 0) {
            pos.debt += interestAccrued;
            totalDebt += interestAccrued;
            pos.lastUpdate = block.timestamp;

            emit InterestAccrued(user, interestAccrued);
        }

        return interestAccrued;
    }

    /**
     * @notice Get accrued interest for a user without updating state
     * @param user Address of the user
     * @return interestAccrued Amount of interest that would be accrued
     */
    function getAccruedInterest(address user) external view returns (uint256 interestAccrued) {
        Position memory pos = positions[user];

        if (pos.debt == 0) {
            return 0;
        }

        uint256 timeElapsed = block.timestamp - pos.lastUpdate;
        if (timeElapsed == 0) {
            return 0;
        }

        return (pos.debt * INTEREST_RATE * timeElapsed) / (PRECISION * SECONDS_PER_YEAR);
    }

    /**
     * @notice Liquidate an undercollateralized position
     * @param user Address of the position to liquidate
     * @param debtToCover Amount of debt to repay (in LUSD)
     * @dev Liquidator must repay user's debt and receives collateral + bonus
     * @dev Can only liquidate if health factor < LIQUIDATION_THRESHOLD
     */
    function liquidate(address user, uint256 debtToCover) external nonReentrant {
        require(user != msg.sender, "Cannot liquidate own position");
        require(debtToCover > 0, "Must cover some debt");

        // Accrue interest before liquidation
        accrueInterest(user);

        Position storage pos = positions[user];
        require(pos.debt > 0, "No debt to liquidate");

        // Check if position is liquidatable
        uint256 healthFactor = _calculateHealthFactor(pos.collateral, pos.debt);
        require(
            healthFactor < LIQUIDATION_THRESHOLD,
            "Position is healthy, cannot liquidate"
        );

        // Cannot cover more debt than exists
        if (debtToCover > pos.debt) {
            debtToCover = pos.debt;
        }

        // Calculate collateral to seize
        // collateralValue = debtToCover * (100 + LIQUIDATION_BONUS) / 100
        uint256 collateralPrice = priceOracle.getAssetPrice(address(collateralToken));
        uint256 collateralValueToSeize = (debtToCover * (PRECISION + LIQUIDATION_BONUS)) / PRECISION;
        uint256 collateralToSeize = (collateralValueToSeize * PRICE_PRECISION) / collateralPrice;

        // Ensure we don't seize more collateral than available
        if (collateralToSeize > pos.collateral) {
            collateralToSeize = pos.collateral;
        }

        // Update position
        pos.debt -= debtToCover;
        pos.collateral -= collateralToSeize;
        pos.lastUpdate = block.timestamp;

        // Update global state
        totalDebt -= debtToCover;
        totalCollateral -= collateralToSeize;

        // If position is fully liquidated, decrease active positions
        if (pos.collateral == 0 && pos.debt == 0) {
            activePositions--;
        }

        // Transfer debt payment from liquidator
        lusdToken.safeTransferFrom(msg.sender, address(this), debtToCover);

        // Transfer seized collateral to liquidator
        collateralToken.safeTransfer(msg.sender, collateralToSeize);

        emit Liquidated(user, msg.sender, debtToCover, collateralToSeize);

        _updateAggregator();
    }

    /**
     * @notice Check if a position can be liquidated
     * @param user Address of the user
     * @return canLiquidate True if position health factor < LIQUIDATION_THRESHOLD
     */
    function isLiquidatable(address user) external view returns (bool canLiquidate) {
        Position memory pos = positions[user];
        if (pos.debt == 0) {
            return false;
        }

        uint256 healthFactor = _calculateHealthFactor(pos.collateral, pos.debt);
        return healthFactor < LIQUIDATION_THRESHOLD;
    }

    /**
     * @notice Get the total debt including accrued interest for a user
     * @param user Address of the user
     * @return totalDebtWithInterest Total debt including pending interest
     */
    function getTotalDebt(address user) external view returns (uint256 totalDebtWithInterest) {
        Position memory pos = positions[user];

        if (pos.debt == 0) {
            return 0;
        }

        uint256 timeElapsed = block.timestamp - pos.lastUpdate;
        if (timeElapsed == 0) {
            return pos.debt;
        }

        uint256 interest = (pos.debt * INTEREST_RATE * timeElapsed) / (PRECISION * SECONDS_PER_YEAR);
        return pos.debt + interest;
    }

    /**
     * @notice Pause the vault in case of emergency
     * @dev Only owner can pause. Pausing prevents deposit, borrow, and withdraw
     * @dev Liquidations and repayments remain active to protect system health
     */
    function pause() external onlyOwner {
        _pause();
    }

    /**
     * @notice Unpause the vault
     * @dev Only owner can unpause
     */
    function unpause() external onlyOwner {
        _unpause();
    }
}
