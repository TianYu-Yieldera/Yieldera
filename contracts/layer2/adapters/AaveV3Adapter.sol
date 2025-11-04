// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "../interfaces/external/IAaveV3Pool.sol";
import "../interfaces/external/IFlashLoanReceiver.sol";
import "../aggregator/L2StateAggregator.sol";

/**
 * @title AaveV3Adapter
 * @notice Adapter for integrating with Aave V3 Protocol
 * @dev Provides lending, borrowing, and flash loan functionality
 *
 * Key Features:
 * - Supply/withdraw assets to earn interest
 * - Borrow assets with collateral
 * - Flash loans for advanced strategies
 * - Interest rate mode switching
 * - Integration with L2StateAggregator
 */
contract AaveV3Adapter is Ownable, ReentrancyGuard, IFlashLoanReceiver, IFlashLoanSimpleReceiver {
    using SafeERC20 for IERC20;

    // ============ State Variables ============

    /// @notice Aave V3 Pool contract
    IAaveV3Pool public immutable aavePool;

    /// @notice L2 State Aggregator for tracking system state
    L2StateAggregator public stateAggregator;

    /// @notice Interest rate modes
    uint256 public constant INTEREST_RATE_MODE_STABLE = 1;
    uint256 public constant INTEREST_RATE_MODE_VARIABLE = 2;

    /// @notice Referral code (can be used for tracking)
    uint16 public constant REFERRAL_CODE = 0;

    /// @notice Flash loan premium (0.09% = 9 basis points)
    uint256 public constant FLASH_LOAN_PREMIUM = 9; // 0.09%
    uint256 public constant PREMIUM_PRECISION = 10000;

    // User tracking
    struct UserPosition {
        uint256 totalSupplied;
        uint256 totalBorrowed;
        uint256 lastUpdate;
    }

    mapping(address => UserPosition) public userPositions;
    uint256 public totalSupplied;
    uint256 public totalBorrowed;
    uint256 public activeUsers;

    // Flash loan tracking
    mapping(address => bool) public authorizedFlashLoanReceivers;
    address private currentFlashLoanInitiator;

    // ============ Events ============

    event Supplied(
        address indexed user,
        address indexed asset,
        uint256 amount,
        uint256 timestamp
    );

    event Withdrawn(
        address indexed user,
        address indexed asset,
        uint256 amount,
        uint256 timestamp
    );

    event Borrowed(
        address indexed user,
        address indexed asset,
        uint256 amount,
        uint256 interestRateMode,
        uint256 timestamp
    );

    event Repaid(
        address indexed user,
        address indexed asset,
        uint256 amount,
        uint256 interestRateMode,
        uint256 timestamp
    );

    event FlashLoanExecuted(
        address indexed initiator,
        address indexed asset,
        uint256 amount,
        uint256 premium,
        uint256 timestamp
    );

    event CollateralStatusChanged(
        address indexed user,
        address indexed asset,
        bool useAsCollateral
    );

    event InterestRateModeSwapped(
        address indexed user,
        address indexed asset,
        uint256 newMode
    );

    // ============ Constructor ============

    /**
     * @notice Initialize the Aave V3 Adapter
     * @param _aavePool Address of Aave V3 Pool contract
     * @param _stateAggregator Address of L2 State Aggregator
     * @param initialOwner Address of contract owner
     */
    constructor(
        address _aavePool,
        address _stateAggregator,
        address initialOwner
    ) Ownable(initialOwner) {
        require(_aavePool != address(0), "Invalid Aave Pool address");

        aavePool = IAaveV3Pool(_aavePool);
        stateAggregator = L2StateAggregator(_stateAggregator);
    }

    // ============ Core Lending Functions ============

    /**
     * @notice Supply assets to Aave to earn interest
     * @param asset The address of the underlying asset to supply
     * @param amount The amount to be supplied
     */
    function supply(
        address asset,
        uint256 amount
    ) external nonReentrant {
        require(amount > 0, "Amount must be greater than zero");

        // Transfer tokens from user to this contract
        IERC20(asset).safeTransferFrom(msg.sender, address(this), amount);

        // Approve Aave Pool to spend tokens
        IERC20(asset).forceApprove(address(aavePool), amount);

        // Supply to Aave on behalf of user
        aavePool.supply(asset, amount, msg.sender, REFERRAL_CODE);

        // Update user position
        UserPosition storage pos = userPositions[msg.sender];
        if (pos.totalSupplied == 0) {
            activeUsers++;
        }
        pos.totalSupplied += amount;
        pos.lastUpdate = block.timestamp;
        totalSupplied += amount;

        emit Supplied(msg.sender, asset, amount, block.timestamp);

        _updateAggregator();
    }

    /**
     * @notice Withdraw supplied assets from Aave
     * @param asset The address of the underlying asset to withdraw
     * @param amount The amount to withdraw (use type(uint256).max for full withdrawal)
     * @return The actual amount withdrawn
     */
    function withdraw(
        address asset,
        uint256 amount
    ) external nonReentrant returns (uint256) {
        require(amount > 0, "Amount must be greater than zero");

        // Withdraw from Aave (aTokens are burned, underlying is transferred to msg.sender)
        uint256 withdrawn = aavePool.withdraw(asset, amount, msg.sender);

        // Update user position
        UserPosition storage pos = userPositions[msg.sender];
        pos.totalSupplied -= withdrawn;
        pos.lastUpdate = block.timestamp;
        totalSupplied -= withdrawn;

        if (pos.totalSupplied == 0 && pos.totalBorrowed == 0) {
            activeUsers--;
        }

        emit Withdrawn(msg.sender, asset, withdrawn, block.timestamp);

        _updateAggregator();

        return withdrawn;
    }

    /**
     * @notice Borrow assets from Aave
     * @param asset The address of the underlying asset to borrow
     * @param amount The amount to be borrowed
     * @param interestRateMode The interest rate mode (1 = Stable, 2 = Variable)
     */
    function borrow(
        address asset,
        uint256 amount,
        uint256 interestRateMode
    ) external nonReentrant {
        require(amount > 0, "Amount must be greater than zero");
        require(
            interestRateMode == INTEREST_RATE_MODE_STABLE ||
            interestRateMode == INTEREST_RATE_MODE_VARIABLE,
            "Invalid interest rate mode"
        );

        // Borrow from Aave (tokens are transferred to msg.sender)
        aavePool.borrow(
            asset,
            amount,
            interestRateMode,
            REFERRAL_CODE,
            msg.sender
        );

        // Update user position
        UserPosition storage pos = userPositions[msg.sender];
        if (pos.totalBorrowed == 0 && pos.totalSupplied == 0) {
            activeUsers++;
        }
        pos.totalBorrowed += amount;
        pos.lastUpdate = block.timestamp;
        totalBorrowed += amount;

        emit Borrowed(msg.sender, asset, amount, interestRateMode, block.timestamp);

        _updateAggregator();
    }

    /**
     * @notice Repay borrowed assets to Aave
     * @param asset The address of the borrowed underlying asset
     * @param amount The amount to repay (use type(uint256).max for full repayment)
     * @param interestRateMode The interest rate mode (1 = Stable, 2 = Variable)
     * @return The actual amount repaid
     */
    function repay(
        address asset,
        uint256 amount,
        uint256 interestRateMode
    ) external nonReentrant returns (uint256) {
        require(amount > 0, "Amount must be greater than zero");

        // Transfer tokens from user to this contract
        IERC20(asset).safeTransferFrom(msg.sender, address(this), amount);

        // Approve Aave Pool to spend tokens
        IERC20(asset).forceApprove(address(aavePool), amount);

        // Repay to Aave on behalf of user
        uint256 repaid = aavePool.repay(
            asset,
            amount,
            interestRateMode,
            msg.sender
        );

        // Refund any excess tokens
        if (repaid < amount) {
            IERC20(asset).safeTransfer(msg.sender, amount - repaid);
        }

        // Update user position
        UserPosition storage pos = userPositions[msg.sender];
        pos.totalBorrowed -= repaid;
        pos.lastUpdate = block.timestamp;
        totalBorrowed -= repaid;

        if (pos.totalBorrowed == 0 && pos.totalSupplied == 0) {
            activeUsers--;
        }

        emit Repaid(msg.sender, asset, repaid, interestRateMode, block.timestamp);

        _updateAggregator();

        return repaid;
    }

    // ============ Flash Loan Functions ============

    /**
     * @notice Execute a flash loan (single asset)
     * @param asset The address of the asset to flash borrow
     * @param amount The amount to flash borrow
     * @param params Encoded parameters for the flash loan strategy
     */
    function flashLoanSimple(
        address asset,
        uint256 amount,
        bytes calldata params
    ) external nonReentrant {
        require(amount > 0, "Amount must be greater than zero");

        currentFlashLoanInitiator = msg.sender;

        // Execute flash loan
        aavePool.flashLoanSimple(
            address(this),
            asset,
            amount,
            params,
            REFERRAL_CODE
        );

        currentFlashLoanInitiator = address(0);
    }

    /**
     * @notice Execute a flash loan (multiple assets)
     * @param assets Array of asset addresses to flash borrow
     * @param amounts Array of amounts to flash borrow
     * @param interestRateModes Types of debt to open if not returned (0 = no debt, 1/2 = debt)
     * @param params Encoded parameters for the flash loan strategy
     */
    function flashLoan(
        address[] calldata assets,
        uint256[] calldata amounts,
        uint256[] calldata interestRateModes,
        bytes calldata params
    ) external nonReentrant {
        require(assets.length > 0, "Must borrow at least one asset");
        require(
            assets.length == amounts.length &&
            assets.length == interestRateModes.length,
            "Array length mismatch"
        );

        currentFlashLoanInitiator = msg.sender;

        // Execute flash loan
        aavePool.flashLoan(
            address(this),
            assets,
            amounts,
            interestRateModes,
            msg.sender,
            params,
            REFERRAL_CODE
        );

        currentFlashLoanInitiator = address(0);
    }

    /**
     * @notice Callback function for flash loan execution
     * @dev This function is called by Aave Pool after transferring the flash-borrowed assets
     * @param assets The addresses of the flash-borrowed assets
     * @param amounts The amounts of the flash-borrowed assets
     * @param premiums The fees of the flash-borrowed assets
     * @param initiator The address that initiated the flash loan
     * @param params The byte-encoded params passed when initiating the flash loan
     * @return True if the execution succeeds
     */
    function executeOperation(
        address[] calldata assets,
        uint256[] calldata amounts,
        uint256[] calldata premiums,
        address initiator,
        bytes calldata params
    ) external override returns (bool) {
        require(msg.sender == address(aavePool), "Caller must be Aave Pool");
        require(initiator == address(this), "Initiator must be this contract");
        require(currentFlashLoanInitiator != address(0), "No active flash loan");

        // Decode params to get user's strategy contract
        address userStrategy = abi.decode(params, (address));

        // Transfer borrowed assets to user's strategy contract
        for (uint256 i = 0; i < assets.length; i++) {
            IERC20(assets[i]).safeTransfer(userStrategy, amounts[i]);
        }

        // Call user's strategy
        IFlashLoanReceiver(userStrategy).executeOperation(
            assets,
            amounts,
            premiums,
            currentFlashLoanInitiator,
            params
        );

        // Approve Aave Pool to pull back the borrowed amount + premium
        for (uint256 i = 0; i < assets.length; i++) {
            uint256 amountOwing = amounts[i] + premiums[i];
            IERC20(assets[i]).forceApprove(address(aavePool), amountOwing);

            emit FlashLoanExecuted(
                currentFlashLoanInitiator,
                assets[i],
                amounts[i],
                premiums[i],
                block.timestamp
            );
        }

        return true;
    }

    /**
     * @notice Callback function for simple flash loan execution
     * @dev This function is called by Aave Pool after transferring the flash-borrowed asset
     * @param asset The address of the flash-borrowed asset
     * @param amount The amount of the flash-borrowed asset
     * @param premium The fee of the flash-borrowed asset
     * @param initiator The address that initiated the flash loan
     * @param params The byte-encoded params passed when initiating the flash loan
     * @return True if the execution succeeds
     */
    function executeOperation(
        address asset,
        uint256 amount,
        uint256 premium,
        address initiator,
        bytes calldata params
    ) external override returns (bool) {
        require(msg.sender == address(aavePool), "Caller must be Aave Pool");
        require(initiator == address(this), "Initiator must be this contract");
        require(currentFlashLoanInitiator != address(0), "No active flash loan");

        // Decode params to get user's strategy contract
        address userStrategy = abi.decode(params, (address));

        // Transfer borrowed asset to user's strategy contract
        IERC20(asset).safeTransfer(userStrategy, amount);

        // Call user's strategy
        IFlashLoanSimpleReceiver(userStrategy).executeOperation(
            asset,
            amount,
            premium,
            currentFlashLoanInitiator,
            params
        );

        // Approve Aave Pool to pull back the borrowed amount + premium
        uint256 amountOwing = amount + premium;
        IERC20(asset).forceApprove(address(aavePool), amountOwing);

        emit FlashLoanExecuted(
            currentFlashLoanInitiator,
            asset,
            amount,
            premium,
            block.timestamp
        );

        return true;
    }

    // ============ Collateral & Interest Rate Management ============

    /**
     * @notice Enable/disable an asset as collateral
     * @param asset The address of the underlying asset
     * @param useAsCollateral True to use as collateral, false otherwise
     */
    function setUserUseReserveAsCollateral(
        address asset,
        bool useAsCollateral
    ) external {
        aavePool.setUserUseReserveAsCollateral(asset, useAsCollateral);

        emit CollateralStatusChanged(msg.sender, asset, useAsCollateral);
    }

    /**
     * @notice Swap borrow rate mode between stable and variable
     * @param asset The address of the underlying asset borrowed
     * @param currentRateMode The current interest rate mode
     */
    function swapBorrowRateMode(
        address asset,
        uint256 currentRateMode
    ) external {
        aavePool.swapBorrowRateMode(asset, currentRateMode);

        uint256 newMode = currentRateMode == INTEREST_RATE_MODE_STABLE
            ? INTEREST_RATE_MODE_VARIABLE
            : INTEREST_RATE_MODE_STABLE;

        emit InterestRateModeSwapped(msg.sender, asset, newMode);
    }

    // ============ View Functions ============

    /**
     * @notice Get user account data from Aave
     * @param user The address of the user
     * @return totalCollateralBase Total collateral in base currency
     * @return totalDebtBase Total debt in base currency
     * @return availableBorrowsBase Available borrowing power in base currency
     * @return currentLiquidationThreshold Current liquidation threshold
     * @return ltv Loan to value ratio
     * @return healthFactor Current health factor
     */
    function getUserAccountData(address user)
        external
        view
        returns (
            uint256 totalCollateralBase,
            uint256 totalDebtBase,
            uint256 availableBorrowsBase,
            uint256 currentLiquidationThreshold,
            uint256 ltv,
            uint256 healthFactor
        )
    {
        return aavePool.getUserAccountData(user);
    }

    /**
     * @notice Get reserve data for an asset
     * @param asset The address of the underlying asset
     * @return Reserve configuration and rates
     */
    function getReserveData(address asset)
        external
        view
        returns (IAaveV3Pool.ReserveData memory)
    {
        return aavePool.getReserveData(asset);
    }

    /**
     * @notice Get user position tracked by this adapter
     * @param user The address of the user
     * @return Position data
     */
    function getUserPosition(address user)
        external
        view
        returns (UserPosition memory)
    {
        return userPositions[user];
    }

    /**
     * @notice Calculate flash loan premium
     * @param amount The flash loan amount
     * @return The premium amount
     */
    function calculateFlashLoanPremium(uint256 amount)
        external
        pure
        returns (uint256)
    {
        return (amount * FLASH_LOAN_PREMIUM) / PREMIUM_PRECISION;
    }

    // ============ Admin Functions ============

    /**
     * @notice Update state aggregator address
     * @param _newAggregator New state aggregator address
     */
    function updateStateAggregator(address _newAggregator) external onlyOwner {
        stateAggregator = L2StateAggregator(_newAggregator);
    }

    /**
     * @notice Emergency token recovery (only for tokens accidentally sent to contract)
     * @param token Token address
     * @param amount Amount to recover
     */
    function recoverToken(address token, uint256 amount) external onlyOwner {
        IERC20(token).safeTransfer(owner(), amount);
    }

    // ============ Internal Functions ============

    /**
     * @notice Update system state in aggregator
     */
    function _updateAggregator() internal {
        if (address(stateAggregator) != address(0)) {
            stateAggregator.updateSystemState(
                totalSupplied,
                totalBorrowed,
                activeUsers,
                0 // No orders in Aave adapter
            );
        }
    }
}
