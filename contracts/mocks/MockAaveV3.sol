// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "../layer2/interfaces/external/IAaveV3Pool.sol";
import "../layer2/interfaces/external/IFlashLoanReceiver.sol";

/**
 * @title MockAaveV3Pool
 * @notice Mock implementation of Aave V3 Pool for testing
 * @dev Simplified implementation for unit tests only
 */
contract MockAaveV3Pool is IAaveV3Pool {
    // Supply tracking
    mapping(address => mapping(address => uint256)) public userSupplies; // user => asset => amount
    mapping(address => mapping(address => uint256)) public userBorrows; // user => asset => amount
    mapping(address => mapping(address => bool)) public userUseAsCollateral; // user => asset => enabled

    // Reserve data
    mapping(address => ReserveData) private reserveData;
    mapping(address => address) public aTokens; // asset => aToken

    // Interest rates (simplified, in ray format: 1e27 = 100%)
    uint128 public constant DEFAULT_LIQUIDITY_RATE = 3e25; // 3% APR
    uint128 public constant DEFAULT_VARIABLE_BORROW_RATE = 5e25; // 5% APR
    uint128 public constant DEFAULT_STABLE_BORROW_RATE = 6e25; // 6% APR

    // Flash loan premium: 0.09% = 9 basis points
    uint256 public constant FLASHLOAN_PREMIUM_TOTAL = 9;

    // Events
    event SupplyExecuted(address indexed user, address indexed asset, uint256 amount);
    event WithdrawExecuted(address indexed user, address indexed asset, uint256 amount);
    event BorrowExecuted(address indexed user, address indexed asset, uint256 amount, uint256 rateMode);
    event RepayExecuted(address indexed user, address indexed asset, uint256 amount, uint256 rateMode);

    /**
     * @notice Initialize mock aToken for an asset
     * @param asset The underlying asset
     * @param aToken The aToken address
     */
    function setAToken(address asset, address aToken) external {
        aTokens[asset] = aToken;
        _initReserveData(asset);
    }

    /**
     * @notice Supply assets to the pool
     */
    function supply(
        address asset,
        uint256 amount,
        address onBehalfOf,
        uint16 /* referralCode */
    ) external override {
        require(amount > 0, "Amount must be > 0");

        // Transfer tokens to pool
        IERC20(asset).transferFrom(msg.sender, address(this), amount);

        // Update supply balance
        userSupplies[onBehalfOf][asset] += amount;

        // Mint aTokens to user
        if (aTokens[asset] != address(0)) {
            MockAToken(aTokens[asset]).mint(onBehalfOf, amount);
        }

        emit SupplyExecuted(onBehalfOf, asset, amount);
    }

    /**
     * @notice Withdraw assets from the pool
     */
    function withdraw(
        address asset,
        uint256 amount,
        address to
    ) external override returns (uint256) {
        require(amount > 0, "Amount must be > 0");

        // In real Aave, the caller must have aTokens. Check aToken balance.
        uint256 aTokenBalance = aTokens[asset] != address(0)
            ? MockAToken(aTokens[asset]).balanceOf(to)
            : userSupplies[to][asset];

        uint256 withdrawAmount = amount == type(uint256).max ? aTokenBalance : amount;

        require(withdrawAmount <= aTokenBalance, "Insufficient supply");

        // Update supply balance
        userSupplies[to][asset] -= withdrawAmount;

        // Burn aTokens from the actual token holder (to address)
        if (aTokens[asset] != address(0)) {
            MockAToken(aTokens[asset]).burn(to, withdrawAmount);
        }

        // Transfer underlying tokens to the "to" address
        IERC20(asset).transfer(to, withdrawAmount);

        emit WithdrawExecuted(to, asset, withdrawAmount);

        return withdrawAmount;
    }

    /**
     * @notice Borrow assets from the pool
     */
    function borrow(
        address asset,
        uint256 amount,
        uint256 interestRateMode,
        uint16 /* referralCode */,
        address onBehalfOf
    ) external override {
        require(amount > 0, "Amount must be > 0");
        require(
            interestRateMode == 1 || interestRateMode == 2,
            "Invalid rate mode"
        );

        // Update borrow balance for the user (onBehalfOf)
        userBorrows[onBehalfOf][asset] += amount;

        // In real Aave, tokens go to msg.sender, but for simplicity in tests,
        // we send directly to onBehalfOf (the user)
        IERC20(asset).transfer(onBehalfOf, amount);

        emit BorrowExecuted(onBehalfOf, asset, amount, interestRateMode);
    }

    /**
     * @notice Repay borrowed assets
     */
    function repay(
        address asset,
        uint256 amount,
        uint256 interestRateMode,
        address onBehalfOf
    ) external override returns (uint256) {
        require(amount > 0, "Amount must be > 0");

        uint256 userDebt = userBorrows[onBehalfOf][asset];
        uint256 repayAmount = amount > userDebt ? userDebt : amount;

        // Update borrow balance
        userBorrows[onBehalfOf][asset] -= repayAmount;

        // Transfer tokens from user
        IERC20(asset).transferFrom(msg.sender, address(this), repayAmount);

        emit RepayExecuted(onBehalfOf, asset, repayAmount, interestRateMode);

        return repayAmount;
    }

    /**
     * @notice Set asset as collateral
     */
    function setUserUseReserveAsCollateral(
        address asset,
        bool useAsCollateral
    ) external override {
        userUseAsCollateral[msg.sender][asset] = useAsCollateral;
    }

    /**
     * @notice Execute flash loan
     */
    function flashLoan(
        address receiverAddress,
        address[] calldata assets,
        uint256[] calldata amounts,
        uint256[] calldata interestRateModes,
        address /* onBehalfOf */,
        bytes calldata params,
        uint16 /* referralCode */
    ) external override {
        require(assets.length > 0, "Empty arrays");
        require(
            assets.length == amounts.length && assets.length == interestRateModes.length,
            "Inconsistent arrays"
        );

        uint256[] memory premiums = new uint256[](assets.length);

        // Transfer assets to receiver
        for (uint256 i = 0; i < assets.length; i++) {
            premiums[i] = (amounts[i] * FLASHLOAN_PREMIUM_TOTAL) / 10000;
            IERC20(assets[i]).transfer(receiverAddress, amounts[i]);
        }

        // Call receiver
        require(
            IFlashLoanReceiver(receiverAddress).executeOperation(
                assets,
                amounts,
                premiums,
                msg.sender,
                params
            ),
            "Flash loan execution failed"
        );

        // Pull back assets + premium
        for (uint256 i = 0; i < assets.length; i++) {
            uint256 amountPlusPremium = amounts[i] + premiums[i];
            IERC20(assets[i]).transferFrom(
                receiverAddress,
                address(this),
                amountPlusPremium
            );
        }
    }

    /**
     * @notice Execute simple flash loan (single asset)
     */
    function flashLoanSimple(
        address receiverAddress,
        address asset,
        uint256 amount,
        bytes calldata params,
        uint16 /* referralCode */
    ) external override {
        require(amount > 0, "Amount must be > 0");

        uint256 premium = (amount * FLASHLOAN_PREMIUM_TOTAL) / 10000;

        // Transfer asset to receiver
        IERC20(asset).transfer(receiverAddress, amount);

        // Call receiver
        require(
            IFlashLoanSimpleReceiver(receiverAddress).executeOperation(
                asset,
                amount,
                premium,
                msg.sender,
                params
            ),
            "Flash loan execution failed"
        );

        // Pull back asset + premium
        uint256 amountPlusPremium = amount + premium;
        IERC20(asset).transferFrom(
            receiverAddress,
            address(this),
            amountPlusPremium
        );
    }

    /**
     * @notice Get user account data
     */
    function getUserAccountData(address /* user */)
        external
        pure
        override
        returns (
            uint256 totalCollateralBase,
            uint256 totalDebtBase,
            uint256 availableBorrowsBase,
            uint256 currentLiquidationThreshold,
            uint256 ltv,
            uint256 healthFactor
        )
    {
        // Simplified: return mock values
        // In real Aave, this would aggregate across all assets
        totalCollateralBase = 1000e8; // $1000 in base currency (8 decimals)
        totalDebtBase = 500e8; // $500 debt
        availableBorrowsBase = 300e8; // $300 available
        currentLiquidationThreshold = 8000; // 80%
        ltv = 7500; // 75%
        healthFactor = 2e18; // Health factor of 2.0 (safe)
    }

    /**
     * @notice Swap borrow rate mode
     */
    function swapBorrowRateMode(
        address /* asset */,
        uint256 /* interestRateMode */
    ) external override {
        // Mock implementation: no-op
    }

    /**
     * @notice Get reserve data
     */
    function getReserveData(address asset)
        external
        view
        override
        returns (ReserveData memory)
    {
        return reserveData[asset];
    }

    /**
     * @notice Initialize reserve data for testing
     */
    function _initReserveData(address asset) internal {
        ReserveData storage data = reserveData[asset];
        data.liquidityIndex = 1e27; // 1.0 in ray
        data.currentLiquidityRate = DEFAULT_LIQUIDITY_RATE;
        data.variableBorrowIndex = 1e27;
        data.currentVariableBorrowRate = DEFAULT_VARIABLE_BORROW_RATE;
        data.currentStableBorrowRate = DEFAULT_STABLE_BORROW_RATE;
        data.lastUpdateTimestamp = uint40(block.timestamp);
        data.aTokenAddress = aTokens[asset];
    }

    /**
     * @notice Helper to get user supply
     */
    function getUserSupply(address user, address asset) external view returns (uint256) {
        return userSupplies[user][asset];
    }

    /**
     * @notice Helper to get user borrow
     */
    function getUserBorrow(address user, address asset) external view returns (uint256) {
        return userBorrows[user][asset];
    }
}

/**
 * @title MockAToken
 * @notice Mock aToken (interest-bearing token)
 */
contract MockAToken is ERC20 {
    address public pool;

    constructor(string memory name, string memory symbol) ERC20(name, symbol) {
        // Don't set pool in constructor - will be set by setPool
    }

    function setPool(address _pool) external {
        require(pool == address(0), "Pool already set");
        pool = _pool;
    }

    function mint(address to, uint256 amount) external {
        require(msg.sender == pool, "Only pool can mint");
        _mint(to, amount);
    }

    function burn(address from, uint256 amount) external {
        require(msg.sender == pool, "Only pool can burn");
        _burn(from, amount);
    }
}

/**
 * @title MockFlashLoanReceiver
 * @notice Mock flash loan receiver for testing
 */
contract MockFlashLoanReceiver is IFlashLoanReceiver, IFlashLoanSimpleReceiver {
    address public adapter;
    bool public shouldFail;
    bool public shouldNotRepay;

    event FlashLoanReceived(address[] assets, uint256[] amounts, uint256[] premiums);
    event SimpleFlashLoanReceived(address asset, uint256 amount, uint256 premium);

    constructor(address _adapter) {
        adapter = _adapter;
    }

    function setShouldFail(bool _shouldFail) external {
        shouldFail = _shouldFail;
    }

    function setShouldNotRepay(bool _shouldNotRepay) external {
        shouldNotRepay = _shouldNotRepay;
    }

    /**
     * @notice Handle flash loan (multiple assets)
     */
    function executeOperation(
        address[] calldata assets,
        uint256[] calldata amounts,
        uint256[] calldata premiums,
        address /* initiator */,
        bytes calldata /* params */
    ) external override returns (bool) {
        require(!shouldFail, "Flash loan failed intentionally");

        emit FlashLoanReceived(assets, amounts, premiums);

        if (!shouldNotRepay) {
            // Transfer tokens back to adapter (msg.sender)
            for (uint256 i = 0; i < assets.length; i++) {
                uint256 amountOwing = amounts[i] + premiums[i];
                IERC20(assets[i]).transfer(msg.sender, amountOwing);
            }
        }

        return true;
    }

    /**
     * @notice Handle simple flash loan (single asset)
     */
    function executeOperation(
        address asset,
        uint256 amount,
        uint256 premium,
        address /* initiator */,
        bytes calldata /* params */
    ) external override returns (bool) {
        require(!shouldFail, "Flash loan failed intentionally");

        emit SimpleFlashLoanReceived(asset, amount, premium);

        if (!shouldNotRepay) {
            // Transfer tokens back to adapter (msg.sender)
            uint256 amountOwing = amount + premium;
            IERC20(asset).transfer(msg.sender, amountOwing);
        }

        return true;
    }

    /**
     * @notice Fund receiver with tokens for repayment
     */
    function fundReceiver(address token, uint256 amount) external {
        IERC20(token).transferFrom(msg.sender, address(this), amount);
    }
}
