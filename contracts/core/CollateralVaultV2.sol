// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "./interfaces/IPriceOracle.sol";

/**
 * @title CollateralVaultV2
 * @notice Enhanced CollateralVault with price oracle integration
 * @dev Manages collateral with real-time price feeds for accurate valuations
 */
contract CollateralVaultV2 is ReentrancyGuard, Ownable {
    using SafeERC20 for IERC20;

    // Core contracts
    IERC20 public immutable loyaltyToken;
    IPriceOracle public priceOracle;
    address public stabilityManager;

    // User collateral balances
    mapping(address => uint256) public collateralDeposited;
    mapping(address => uint256) public debtAmount;
    mapping(address => uint256) public lastInterestUpdate;

    // Vault totals
    uint256 public totalCollateral;
    uint256 public totalDebt;

    // Parameters (in percentage, 150 = 150%)
    uint256 public constant COLLATERAL_RATIO = 150;
    uint256 public constant LIQUIDATION_THRESHOLD = 120;
    uint256 public constant STABILITY_FEE = 200; // 2% annual
    uint256 public constant SECONDS_PER_YEAR = 365 days;
    uint256 public constant LIQUIDATION_BONUS = 110; // 10% bonus

    // Price precision
    uint256 public constant PRICE_PRECISION = 1e8;
    uint256 public constant LUSD_DECIMALS = 1e6;

    // Emergency mode
    bool public emergencyMode = false;

    // Events
    event CollateralDeposited(address indexed user, uint256 amount, uint256 newTotal);
    event CollateralWithdrawn(address indexed user, uint256 amount, uint256 remaining);
    event DebtIncreased(address indexed user, uint256 amount, uint256 newDebt);
    event DebtDecreased(address indexed user, uint256 amount, uint256 remaining);
    event PositionLiquidated(address indexed user, uint256 collateralSeized, uint256 debtRepaid);
    event PriceOracleUpdated(address indexed newOracle);
    event StabilityManagerUpdated(address indexed newManager);
    event EmergencyModeActivated();
    event EmergencyModeDeactivated();

    modifier onlyStabilityManager() {
        require(msg.sender == stabilityManager || msg.sender == owner(), "Not authorized");
        _;
    }

    modifier notInEmergency() {
        require(!emergencyMode, "Emergency mode active");
        _;
    }

    constructor(address _loyaltyToken, address _priceOracle) {
        require(_loyaltyToken != address(0), "Invalid token address");
        require(_priceOracle != address(0), "Invalid oracle address");

        loyaltyToken = IERC20(_loyaltyToken);
        priceOracle = IPriceOracle(_priceOracle);
    }

    /**
     * @notice Set the price oracle address
     * @param _oracle New oracle address
     */
    function setPriceOracle(address _oracle) external onlyOwner {
        require(_oracle != address(0), "Invalid oracle address");
        priceOracle = IPriceOracle(_oracle);
        emit PriceOracleUpdated(_oracle);
    }

    /**
     * @notice Set the stability manager address
     * @param _manager New manager address
     */
    function setStabilityManager(address _manager) external onlyOwner {
        require(_manager != address(0), "Invalid manager address");
        stabilityManager = _manager;
        emit StabilityManagerUpdated(_manager);
    }

    /**
     * @notice Get collateral value in USD
     * @param amount Amount of collateral tokens
     * @return Value in USD (with LUSD_DECIMALS precision)
     */
    function getCollateralValueUSD(uint256 amount) public view returns (uint256) {
        if (amount == 0) return 0;

        try priceOracle.getLatestPrice(address(loyaltyToken)) returns (
            uint256 price,
            uint8 decimals,
            uint256 /* timestamp */
        ) {
            // Convert price to LUSD decimals (6)
            // Price has 8 decimals, amount has 6 decimals, need result with 6 decimals
            // (amount * price) / 10^8 = result with 6 decimals
            return (amount * price) / PRICE_PRECISION;
        } catch {
            // Fallback to 1:1 if oracle fails
            return amount;
        }
    }

    /**
     * @notice Get required collateral for a given debt amount
     * @param debtAmountUSD Debt amount in USD
     * @return Required collateral in tokens
     */
    function getRequiredCollateral(uint256 debtAmountUSD) public view returns (uint256) {
        if (debtAmountUSD == 0) return 0;

        try priceOracle.getLatestPrice(address(loyaltyToken)) returns (
            uint256 price,
            uint8 decimals,
            uint256 /* timestamp */
        ) {
            // Calculate required collateral with ratio
            uint256 requiredValueUSD = (debtAmountUSD * COLLATERAL_RATIO) / 100;
            // Convert USD value to token amount
            return (requiredValueUSD * PRICE_PRECISION) / price;
        } catch {
            // Fallback to simple ratio if oracle fails
            return (debtAmountUSD * COLLATERAL_RATIO) / 100;
        }
    }

    /**
     * @notice Deposit collateral
     * @param amount Amount to deposit
     */
    function depositCollateral(uint256 amount) external nonReentrant notInEmergency {
        require(amount > 0, "Amount must be > 0");

        loyaltyToken.safeTransferFrom(msg.sender, address(this), amount);

        collateralDeposited[msg.sender] += amount;
        totalCollateral += amount;

        emit CollateralDeposited(msg.sender, amount, collateralDeposited[msg.sender]);
    }

    /**
     * @notice Withdraw collateral
     * @param amount Amount to withdraw
     */
    function withdrawCollateral(uint256 amount) external nonReentrant notInEmergency {
        require(collateralDeposited[msg.sender] >= amount, "Insufficient collateral");

        uint256 newCollateral = collateralDeposited[msg.sender] - amount;
        uint256 currentDebt = getTotalDebt(msg.sender);

        if (currentDebt > 0) {
            uint256 collateralValueUSD = getCollateralValueUSD(newCollateral);
            require(
                (collateralValueUSD * 100) / currentDebt >= COLLATERAL_RATIO,
                "Withdrawal would undercollateralize position"
            );
        }

        collateralDeposited[msg.sender] = newCollateral;
        totalCollateral -= amount;

        loyaltyToken.safeTransfer(msg.sender, amount);

        emit CollateralWithdrawn(msg.sender, amount, newCollateral);
    }

    /**
     * @notice Calculate accrued interest
     * @param user User address
     * @return Interest amount
     */
    function accruedInterest(address user) public view returns (uint256) {
        uint256 debt = debtAmount[user];
        if (debt == 0) return 0;

        uint256 lastUpdate = lastInterestUpdate[user];
        if (lastUpdate == 0) return 0;

        uint256 timePassed = block.timestamp - lastUpdate;
        if (timePassed == 0) return 0;

        // Interest = principal * rate * time / year
        uint256 interest = (debt * STABILITY_FEE * timePassed) / (10000 * SECONDS_PER_YEAR);

        return interest;
    }

    /**
     * @notice Get total debt including interest
     * @param user User address
     * @return Total debt
     */
    function getTotalDebt(address user) public view returns (uint256) {
        return debtAmount[user] + accruedInterest(user);
    }

    /**
     * @notice Accrue and compound interest
     * @param user User address
     */
    function _accrueInterest(address user) internal {
        uint256 interest = accruedInterest(user);
        if (interest > 0) {
            debtAmount[user] += interest;
            totalDebt += interest;
        }
        lastInterestUpdate[user] = block.timestamp;
    }

    /**
     * @notice Increase debt (mint LUSD)
     * @param user User address
     * @param amount Amount to increase
     */
    function increaseDebt(address user, uint256 amount)
        external
        onlyStabilityManager
        notInEmergency
    {
        require(amount > 0, "Amount must be positive");

        _accrueInterest(user);

        uint256 newDebt = debtAmount[user] + amount;
        uint256 collateralValue = getCollateralValueUSD(collateralDeposited[user]);

        require(
            (collateralValue * 100) / newDebt >= COLLATERAL_RATIO,
            "Insufficient collateral"
        );

        debtAmount[user] = newDebt;
        totalDebt += amount;

        emit DebtIncreased(user, amount, newDebt);
    }

    /**
     * @notice Decrease debt (burn LUSD)
     * @param user User address
     * @param amount Amount to decrease
     */
    function decreaseDebt(address user, uint256 amount)
        external
        onlyStabilityManager
    {
        _accrueInterest(user);

        uint256 currentDebt = debtAmount[user];
        require(currentDebt >= amount, "Debt underflow");

        debtAmount[user] = currentDebt - amount;
        totalDebt -= amount;

        emit DebtDecreased(user, amount, debtAmount[user]);
    }

    /**
     * @notice Get collateral ratio with price oracle
     * @param user User address
     * @return Ratio in percentage
     */
    function getCollateralRatio(address user) public view returns (uint256) {
        uint256 debt = getTotalDebt(user);
        if (debt == 0) return type(uint256).max;

        uint256 collateralValue = getCollateralValueUSD(collateralDeposited[user]);
        return (collateralValue * 100) / debt;
    }

    /**
     * @notice Check if position can be liquidated
     * @param user User address
     * @return True if liquidatable
     */
    function canLiquidate(address user) external view returns (bool) {
        uint256 debt = getTotalDebt(user);
        if (debt == 0) return false;

        return getCollateralRatio(user) < LIQUIDATION_THRESHOLD;
    }

    /**
     * @notice Liquidate position with price-aware calculations
     * @param user User to liquidate
     * @param debtToCover Debt amount to cover
     * @return collateralSeized Amount of collateral seized
     */
    function liquidate(address user, uint256 debtToCover)
        external
        onlyOwner
        nonReentrant
        returns (uint256 collateralSeized)
    {
        require(this.canLiquidate(user), "Position is not liquidatable");

        _accrueInterest(user);

        uint256 userDebt = debtAmount[user];
        require(debtToCover <= userDebt, "Debt to cover exceeds user debt");

        // Calculate collateral to seize based on current price
        uint256 debtValueUSD = debtToCover;
        uint256 seizeValueUSD = (debtValueUSD * LIQUIDATION_BONUS) / 100;

        // Convert USD value to collateral tokens
        try priceOracle.getLatestPrice(address(loyaltyToken)) returns (
            uint256 price,
            uint8 decimals,
            uint256 /* timestamp */
        ) {
            collateralSeized = (seizeValueUSD * PRICE_PRECISION) / price;
        } catch {
            // Fallback calculation
            collateralSeized = (seizeValueUSD * LIQUIDATION_BONUS) / 100;
        }

        // Ensure we don't seize more than available
        uint256 userCollateral = collateralDeposited[user];
        if (collateralSeized > userCollateral) {
            collateralSeized = userCollateral;
        }

        // Update state
        collateralDeposited[user] -= collateralSeized;
        debtAmount[user] -= debtToCover;
        totalCollateral -= collateralSeized;
        totalDebt -= debtToCover;

        // Transfer seized collateral to liquidator
        loyaltyToken.safeTransfer(msg.sender, collateralSeized);

        emit PositionLiquidated(user, collateralSeized, debtToCover);

        return collateralSeized;
    }

    /**
     * @notice Get maximum mintable amount based on collateral
     * @param user User address
     * @return Maximum LUSD that can be minted
     */
    function getMaxMintable(address user) external view returns (uint256) {
        uint256 collateralValue = getCollateralValueUSD(collateralDeposited[user]);
        uint256 maxDebt = (collateralValue * 100) / COLLATERAL_RATIO;
        uint256 currentDebt = getTotalDebt(user);

        if (maxDebt <= currentDebt) {
            return 0;
        }

        return maxDebt - currentDebt;
    }

    /**
     * @notice Get position details with USD values
     * @param user User address
     */
    function getPosition(address user) external view returns (
        uint256 collateral,
        uint256 collateralValueUSD,
        uint256 debt,
        uint256 interest,
        uint256 totalDebtAmount,
        uint256 collateralRatio,
        uint256 maxMintable,
        bool healthy
    ) {
        collateral = collateralDeposited[user];
        collateralValueUSD = getCollateralValueUSD(collateral);
        debt = debtAmount[user];
        interest = accruedInterest(user);
        totalDebtAmount = debt + interest;
        collateralRatio = getCollateralRatio(user);
        maxMintable = this.getMaxMintable(user);
        healthy = collateralRatio >= COLLATERAL_RATIO;
    }

    /**
     * @notice Activate emergency mode (pause most operations)
     */
    function activateEmergencyMode() external onlyOwner {
        emergencyMode = true;
        emit EmergencyModeActivated();
    }

    /**
     * @notice Deactivate emergency mode
     */
    function deactivateEmergencyMode() external onlyOwner {
        emergencyMode = false;
        emit EmergencyModeDeactivated();
    }

    /**
     * @notice Emergency withdraw (only in emergency mode)
     * @param user User address
     */
    function emergencyWithdraw(address user) external nonReentrant {
        require(emergencyMode, "Not in emergency mode");
        require(msg.sender == user, "Can only withdraw own funds");

        uint256 amount = collateralDeposited[user];
        require(amount > 0, "No collateral to withdraw");

        collateralDeposited[user] = 0;
        totalCollateral -= amount;

        loyaltyToken.safeTransfer(user, amount);

        emit CollateralWithdrawn(user, amount, 0);
    }
}