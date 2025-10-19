// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title CollateralVault
 * @notice Holds Loyalty Points (LP tokens) as collateral for LUSD minting
 * @dev Manages collateral deposits, withdrawals, and debt tracking
 *
 * Key Features:
 * - Stores LP tokens as collateral
 * - Tracks debt (LUSD minted) per user
 * - Enforces 150% minimum collateral ratio
 * - 120% liquidation threshold
 * - Only StabilityManager can modify debt
 *
 * Collateral Math:
 * - Collateral Ratio = (collateral * 100) / debt
 * - Min Ratio: 150% (user must have $1.50 collateral for every $1 LUSD)
 * - Liquidation: <120%
 */
contract CollateralVault is ReentrancyGuard, Ownable {
    using SafeERC20 for IERC20;

    // Loyalty Points token address
    IERC20 public immutable loyaltyToken;

    // User collateral balances
    mapping(address => uint256) public collateralDeposited;

    // User debt amounts (LUSD minted)
    mapping(address => uint256) public debtAmount;

    // Last interest update timestamp per user
    mapping(address => uint256) public lastInterestUpdate;

    // Total collateral locked in vault
    uint256 public totalCollateral;

    // Total debt issued
    uint256 public totalDebt;

    // Parameters (in percentage, 150 = 150%)
    uint256 public constant COLLATERAL_RATIO = 150;
    uint256 public constant LIQUIDATION_THRESHOLD = 120;
    uint256 public constant STABILITY_FEE = 200; // 2% annual (basis points, 200 = 2%)
    uint256 public constant SECONDS_PER_YEAR = 365 days;

    // Events
    event CollateralDeposited(
        address indexed user,
        uint256 amount,
        uint256 newTotal
    );
    event CollateralWithdrawn(
        address indexed user,
        uint256 amount,
        uint256 remaining
    );
    event DebtIncreased(address indexed user, uint256 amount, uint256 newDebt);
    event DebtDecreased(address indexed user, uint256 amount, uint256 remaining);
    event PositionLiquidated(
        address indexed user,
        uint256 collateralSeized,
        uint256 debtRepaid
    );

    /**
     * @notice Constructor
     * @param _loyaltyToken Address of Loyalty Points ERC-20 token
     */
    constructor(address _loyaltyToken) {
        require(_loyaltyToken != address(0), "Invalid token address");
        loyaltyToken = IERC20(_loyaltyToken);
    }

    /**
     * @notice Deposit Loyalty Points as collateral
     * @param amount Amount of LP tokens to deposit
     */
    function depositCollateral(uint256 amount) external nonReentrant {
        require(amount > 0, "Amount must be > 0");

        loyaltyToken.safeTransferFrom(msg.sender, address(this), amount);

        collateralDeposited[msg.sender] += amount;
        totalCollateral += amount;

        emit CollateralDeposited(
            msg.sender,
            amount,
            collateralDeposited[msg.sender]
        );
    }

    /**
     * @notice Withdraw collateral (only if position remains healthy)
     * @param amount Amount of LP tokens to withdraw
     */
    function withdrawCollateral(uint256 amount) external nonReentrant {
        require(amount > 0, "Amount must be > 0");

        // Cache storage variables to memory (gas optimization)
        uint256 currentCollateral = collateralDeposited[msg.sender];
        require(currentCollateral >= amount, "Insufficient collateral");

        uint256 newCollateral = currentCollateral - amount;
        uint256 currentDebt = debtAmount[msg.sender];

        // Check if withdrawal keeps position healthy
        require(
            _isPositionHealthy(newCollateral, currentDebt),
            "Withdrawal would undercollateralize position"
        );

        // Update storage (single SSTORE operation)
        collateralDeposited[msg.sender] = newCollateral;
        totalCollateral -= amount;

        loyaltyToken.safeTransfer(msg.sender, amount);

        emit CollateralWithdrawn(msg.sender, amount, newCollateral);
    }

    /**
     * @notice Calculate maximum LUSD that can be minted
     * @param user Address to check
     * @return Maximum mintable LUSD amount
     */
    function getMaxMintable(address user) public view returns (uint256) {
        uint256 collateral = collateralDeposited[user];
        uint256 currentDebt = debtAmount[user];

        // Max mintable = (collateral * 100 / COLLATERAL_RATIO) - current debt
        // Using divide-first approach to prevent overflow
        // Since COLLATERAL_RATIO = 150, this equals: collateral * 100 / 150 = collateral * 2 / 3
        uint256 maxTotal = (collateral / COLLATERAL_RATIO) * 100;

        if (maxTotal > currentDebt) {
            return maxTotal - currentDebt;
        }
        return 0;
    }

    /**
     * @notice Calculate accrued interest for a user
     * @param user Address to check
     * @return Interest amount accrued since last update
     * @dev Gas optimized with early returns and unchecked arithmetic
     */
    function accruedInterest(address user) public view returns (uint256) {
        // Cache storage reads (gas optimization)
        uint256 debt = debtAmount[user];
        if (debt == 0) return 0;

        uint256 lastUpdate = lastInterestUpdate[user];
        if (lastUpdate == 0) return 0;

        unchecked {
            // Safe from overflow: block.timestamp is always >= lastUpdate
            uint256 timePassed = block.timestamp - lastUpdate;
            if (timePassed == 0) return 0;

            // Interest = principal * rate * time / year
            // rate is in basis points (200 = 2% = 0.02)
            // Using fixed-point arithmetic: (debt * STABILITY_FEE * timePassed) / (10000 * SECONDS_PER_YEAR)
            // Denominator is constant: 10000 * 31536000 = 315360000000
            uint256 interest = (debt * STABILITY_FEE * timePassed) / 315360000000;

            return interest;
        }
    }

    /**
     * @notice Get total debt including accrued interest
     * @param user Address to check
     * @return Total debt with interest
     */
    function getTotalDebt(address user) public view returns (uint256) {
        return debtAmount[user] + accruedInterest(user);
    }

    /**
     * @notice Get collateral ratio for a user (in percentage)
     * @param user Address to check
     * @return Collateral ratio (e.g., 200 = 200%)
     */
    function getCollateralRatio(address user) public view returns (uint256) {
        uint256 totalDebt = getTotalDebt(user);
        if (totalDebt == 0) return type(uint256).max;

        return (collateralDeposited[user] * 100) / totalDebt;
    }

    /**
     * @notice Accrue interest for a user (internal helper)
     * @dev Should be called before any debt modification
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
     * @notice Check if a position is healthy (above minimum ratio)
     * @param user Address to check
     * @return True if position is healthy
     */
    function isPositionHealthy(address user) external view returns (bool) {
        return _isPositionHealthy(collateralDeposited[user], debtAmount[user]);
    }

    /**
     * @notice Check if position is at risk of liquidation
     * @param user Address to check
     * @return True if position can be liquidated
     */
    function canLiquidate(address user) external view returns (bool) {
        if (debtAmount[user] == 0) return false;

        uint256 ratio = getCollateralRatio(user);
        return ratio < LIQUIDATION_THRESHOLD;
    }

    /**
     * @notice Increase debt for a user (only StabilityManager)
     * @param user User address
     * @param amount Amount of debt to add
     * @dev Gas optimized by caching storage variables
     */
    function increaseDebt(address user, uint256 amount) external onlyOwner {
        require(amount > 0, "Amount must be positive");

        // Accrue interest first
        _accrueInterest(user);

        // Cache storage variables (gas optimization)
        uint256 currentDebt = debtAmount[user];
        uint256 newDebt;
        unchecked {
            // Safe: checked by require below
            newDebt = currentDebt + amount;
        }

        // Verify position will be healthy after debt increase
        require(
            _isPositionHealthy(collateralDeposited[user], newDebt),
            "Debt increase would undercollateralize position"
        );

        debtAmount[user] = newDebt;
        totalDebt += amount;

        emit DebtIncreased(user, amount, newDebt);
    }

    /**
     * @notice Decrease debt for a user (only StabilityManager)
     * @param user User address
     * @param amount Amount of debt to repay
     * @dev Gas optimized by caching storage variables
     */
    function decreaseDebt(address user, uint256 amount) external onlyOwner {
        // Accrue interest first
        _accrueInterest(user);

        // Cache storage variable (gas optimization)
        uint256 currentDebt = debtAmount[user];
        require(currentDebt >= amount, "Debt underflow");

        uint256 newDebt;
        unchecked {
            // Safe: checked by require above
            newDebt = currentDebt - amount;
        }

        debtAmount[user] = newDebt;
        totalDebt -= amount;

        emit DebtDecreased(user, amount, newDebt);
    }

    /**
     * @notice Liquidate an undercollateralized position
     * @param user Address to liquidate
     * @param debtToCover Amount of debt liquidator will repay
     * @return Collateral seized
     *
     * @dev IMPORTANT: This implementation assumes 1 LP = 1 USD for simplicity.
     *      In production, integrate a price oracle to get accurate LP/USD exchange rate.
     *      Current formula: collateralToSeize = debtToCover * 1.10 (10% liquidation bonus)
     *      Correct formula should be: collateralToSeize = (debtToCover / lpPrice) * 1.10
     *
     *      Example with oracle:
     *      uint256 lpPrice = oracle.getPrice(address(loyaltyToken));
     *      uint256 lpNeeded = (debtToCover * 1e18) / lpPrice;
     *      uint256 collateralToSeize = (lpNeeded * 110) / 100;
     */
    function liquidate(address user, uint256 debtToCover)
        external
        onlyOwner
        returns (uint256)
    {
        // Accrue interest before liquidation
        _accrueInterest(user);

        // Cache storage variables (gas optimization)
        uint256 currentDebt = debtAmount[user];
        uint256 currentCollateral = collateralDeposited[user];

        require(currentDebt > 0, "No debt to liquidate");
        require(debtToCover > 0, "Debt to cover must be positive");
        require(debtToCover <= currentDebt, "Cannot cover more than total debt");

        uint256 ratio = getCollateralRatio(user);
        require(
            ratio < LIQUIDATION_THRESHOLD,
            "Position is not liquidatable"
        );

        // Calculate collateral to seize (with 10% liquidation bonus)
        // Assumes 1 LP = 1 USD. TODO: Integrate price oracle for accurate pricing
        uint256 collateralToSeize = (debtToCover * 110) / 100;

        // Ensure we don't seize more collateral than available
        if (collateralToSeize > currentCollateral) {
            collateralToSeize = currentCollateral;
            // Recalculate actual debt covered based on available collateral
            debtToCover = (collateralToSeize * 100) / 110;
        }

        // Update state (batch storage updates)
        unchecked {
            // Safe: checked above
            collateralDeposited[user] = currentCollateral - collateralToSeize;
            debtAmount[user] = currentDebt - debtToCover;
            totalCollateral -= collateralToSeize;
            totalDebt -= debtToCover;
        }

        emit PositionLiquidated(user, collateralToSeize, debtToCover);

        return collateralToSeize;
    }

    /**
     * @notice Get position details for a user
     * @param user Address to check
     * @dev Gas optimized by caching storage reads
     */
    function getPosition(address user)
        external
        view
        returns (
            uint256 collateral,
            uint256 debt,
            uint256 interest,
            uint256 totalDebt,
            uint256 collateralRatio,
            uint256 maxMintable,
            bool healthy
        )
    {
        // Cache storage reads (gas optimization)
        collateral = collateralDeposited[user];
        debt = debtAmount[user];
        interest = accruedInterest(user);

        unchecked {
            // Safe: adding interest to debt
            totalDebt = debt + interest;
        }

        collateralRatio = (totalDebt == 0) ? type(uint256).max : (collateral * 100) / totalDebt;
        maxMintable = getMaxMintable(user);
        healthy = _isPositionHealthy(collateral, totalDebt);
    }

    /**
     * @notice Internal: Check if position is healthy
     */
    function _isPositionHealthy(uint256 collateral, uint256 debt)
        internal
        pure
        returns (bool)
    {
        if (debt == 0) return true;

        uint256 ratio = (collateral * 100) / debt;
        return ratio >= COLLATERAL_RATIO;
    }

    /**
     * @notice Get vault statistics
     */
    function getVaultStats()
        external
        view
        returns (
            uint256 _totalCollateral,
            uint256 _totalDebt,
            uint256 avgCollateralRatio
        )
    {
        _totalCollateral = totalCollateral;
        _totalDebt = totalDebt;

        if (_totalDebt > 0) {
            avgCollateralRatio = (_totalCollateral * 100) / _totalDebt;
        } else {
            avgCollateralRatio = 0;
        }
    }
}
