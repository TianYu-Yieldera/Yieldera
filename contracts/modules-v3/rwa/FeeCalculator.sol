// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/modules/rwa/IFeeCalculator.sol";
import "../../interfaces/modules/rwa/IOrderManager.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/**
 * @title FeeCalculator
 * @notice Manages trading fee calculations and collection
 * @dev Supports maker/taker fees with configurable discounts
 */
contract FeeCalculator is IFeeCalculator, Ownable {
    using SafeERC20 for IERC20;

    // ============ Constants ============

    uint256 private constant BASIS_POINTS = 10000;
    uint256 private constant MAX_FEE = 1000; // Max 10%

    // ============ State Variables ============

    address public rwaModule; // Main coordinator contract
    IOrderManager public orderManager;
    IERC20 public feeToken; // Token used for fee payment

    // Storage for fees (using Diamond Storage pattern)
    bytes32 private constant FEE_STORAGE_POSITION = keccak256("fee.calculator.storage");

    struct FeeStorage {
        FeeConfig config;
        mapping(address => uint256) userFeesCollected;
        mapping(address => uint256) userFeeDiscount; // Discount in basis points
        uint256 totalFeesCollected;
        uint256 protocolFeeShare; // Percentage of fees for protocol (basis points)
    }

    // ============ Modifiers ============

    modifier onlyRWAModule() {
        require(msg.sender == rwaModule, "Only RWA module");
        _;
    }

    // ============ Constructor ============

    constructor(address _orderManager, address _feeToken, uint256 _makerFee, uint256 _takerFee) {
        require(_orderManager != address(0), "Invalid order manager");
        require(_feeToken != address(0), "Invalid fee token");
        require(_makerFee <= MAX_FEE, "Maker fee too high");
        require(_takerFee <= MAX_FEE, "Taker fee too high");

        orderManager = IOrderManager(_orderManager);
        feeToken = IERC20(_feeToken);

        FeeStorage storage fs = _getStorage();
        fs.config = FeeConfig({
            makerFee: _makerFee,
            takerFee: _takerFee,
            minimumFee: 0,
            feesEnabled: true
        });
        fs.protocolFeeShare = 5000; // 50% to protocol by default
    }

    // ============ Admin Functions ============

    function setRWAModule(address _rwaModule) external onlyOwner {
        require(_rwaModule != address(0), "Invalid address");
        rwaModule = _rwaModule;
    }

    function setOrderManager(address _orderManager) external onlyOwner {
        require(_orderManager != address(0), "Invalid address");
        orderManager = IOrderManager(_orderManager);
    }

    function setProtocolFeeShare(uint256 share) external onlyOwner {
        require(share <= BASIS_POINTS, "Invalid share");
        FeeStorage storage fs = _getStorage();
        fs.protocolFeeShare = share;
    }

    // ============ Internal Storage ============

    function _getStorage() private pure returns (FeeStorage storage fs) {
        bytes32 position = FEE_STORAGE_POSITION;
        assembly {
            fs.slot := position
        }
    }

    // ============ IFeeCalculator Implementation ============

    function calculateFee(uint256 tradeAmount, bool isMaker)
        external
        view
        override
        returns (uint256 fee)
    {
        FeeStorage storage fs = _getStorage();

        if (!fs.config.feesEnabled) return 0;

        uint256 feeRate = isMaker ? fs.config.makerFee : fs.config.takerFee;
        fee = (tradeAmount * feeRate) / BASIS_POINTS;

        // Apply minimum fee
        if (fee < fs.config.minimumFee) {
            fee = fs.config.minimumFee;
        }

        return fee;
    }

    function calculateFeeBreakdown(uint256 tradeAmount, uint256 buyOrderId, uint256 sellOrderId)
        external
        view
        override
        returns (FeeBreakdown memory breakdown)
    {
        FeeStorage storage fs = _getStorage();

        if (!fs.config.feesEnabled) {
            return FeeBreakdown({makerFee: 0, takerFee: 0, protocolFee: 0, totalFee: 0});
        }

        IOrderManager.Order memory buyOrder = orderManager.getOrder(buyOrderId);
        IOrderManager.Order memory sellOrder = orderManager.getOrder(sellOrderId);

        // Determine which is maker (earlier order)
        bool buyIsMaker = buyOrder.timestamp < sellOrder.timestamp;

        uint256 makerFee = (tradeAmount * fs.config.makerFee) / BASIS_POINTS;
        uint256 takerFee = (tradeAmount * fs.config.takerFee) / BASIS_POINTS;

        // Apply discounts
        address makerAddress = buyIsMaker ? buyOrder.trader : sellOrder.trader;
        address takerAddress = buyIsMaker ? sellOrder.trader : buyOrder.trader;

        if (fs.userFeeDiscount[makerAddress] > 0) {
            makerFee = (makerFee * (BASIS_POINTS - fs.userFeeDiscount[makerAddress])) / BASIS_POINTS;
        }
        if (fs.userFeeDiscount[takerAddress] > 0) {
            takerFee = (takerFee * (BASIS_POINTS - fs.userFeeDiscount[takerAddress])) / BASIS_POINTS;
        }

        // Apply minimum fees
        if (makerFee < fs.config.minimumFee) makerFee = fs.config.minimumFee;
        if (takerFee < fs.config.minimumFee) takerFee = fs.config.minimumFee;

        uint256 totalFee = makerFee + takerFee;
        uint256 protocolFee = (totalFee * fs.protocolFeeShare) / BASIS_POINTS;

        breakdown = FeeBreakdown({
            makerFee: makerFee,
            takerFee: takerFee,
            protocolFee: protocolFee,
            totalFee: totalFee
        });

        return breakdown;
    }

    function collectFee(address user, uint256 tradeAmount, bool isMaker)
        external
        override
        onlyRWAModule
        returns (uint256 feeAmount)
    {
        FeeStorage storage fs = _getStorage();

        if (!fs.config.feesEnabled) return 0;

        uint256 feeRate = isMaker ? fs.config.makerFee : fs.config.takerFee;
        feeAmount = (tradeAmount * feeRate) / BASIS_POINTS;

        // Apply discount
        if (fs.userFeeDiscount[user] > 0) {
            feeAmount = (feeAmount * (BASIS_POINTS - fs.userFeeDiscount[user])) / BASIS_POINTS;
        }

        // Apply minimum fee
        if (feeAmount < fs.config.minimumFee) {
            feeAmount = fs.config.minimumFee;
        }

        // Collect fee from user
        feeToken.safeTransferFrom(user, address(this), feeAmount);

        // Update statistics
        fs.userFeesCollected[user] += feeAmount;
        fs.totalFeesCollected += feeAmount;

        emit FeeCollected(user, feeAmount, isMaker, tradeAmount, feeRate);

        return feeAmount;
    }

    function getFeeConfig() external view override returns (FeeConfig memory config) {
        FeeStorage storage fs = _getStorage();
        return fs.config;
    }

    function updateFeeConfig(uint256 makerFee, uint256 takerFee, uint256 minimumFee)
        external
        override
        onlyOwner
    {
        require(makerFee <= MAX_FEE, "Maker fee too high");
        require(takerFee <= MAX_FEE, "Taker fee too high");

        FeeStorage storage fs = _getStorage();
        fs.config.makerFee = makerFee;
        fs.config.takerFee = takerFee;
        fs.config.minimumFee = minimumFee;

        emit FeeConfigUpdated(makerFee, takerFee, minimumFee);
    }

    function getTotalFeesCollected() external view override returns (uint256) {
        FeeStorage storage fs = _getStorage();
        return fs.totalFeesCollected;
    }

    function getUserFeesCollected(address user) external view override returns (uint256) {
        FeeStorage storage fs = _getStorage();
        return fs.userFeesCollected[user];
    }

    function getFeeDiscount(address user)
        external
        view
        override
        returns (bool hasDiscount, uint256 discountRate)
    {
        FeeStorage storage fs = _getStorage();
        discountRate = fs.userFeeDiscount[user];
        hasDiscount = discountRate > 0;
        return (hasDiscount, discountRate);
    }

    function withdrawFees(address recipient, uint256 amount) external override onlyOwner {
        require(recipient != address(0), "Invalid recipient");
        require(amount > 0, "Invalid amount");

        uint256 balance = feeToken.balanceOf(address(this));
        require(balance >= amount, "Insufficient balance");

        feeToken.safeTransfer(recipient, amount);

        emit FeesWithdrawn(recipient, amount);
    }

    // ============ Additional Functions ============

    /**
     * @notice Set fee discount for a user
     * @param user User address
     * @param discount Discount in basis points (e.g., 500 = 5% discount)
     */
    function setUserFeeDiscount(address user, uint256 discount) external onlyOwner {
        require(discount <= BASIS_POINTS, "Invalid discount");
        FeeStorage storage fs = _getStorage();
        fs.userFeeDiscount[user] = discount;
    }

    /**
     * @notice Enable or disable fees
     * @param enabled True to enable fees
     */
    function setFeesEnabled(bool enabled) external onlyOwner {
        FeeStorage storage fs = _getStorage();
        fs.config.feesEnabled = enabled;
    }

    /**
     * @notice Calculate effective fee rate for user (after discount)
     * @param user User address
     * @param isMaker True if user is maker
     * @return effectiveFeeRate Effective fee rate in basis points
     */
    function getEffectiveFeeRate(address user, bool isMaker) external view returns (uint256 effectiveFeeRate) {
        FeeStorage storage fs = _getStorage();

        uint256 baseFee = isMaker ? fs.config.makerFee : fs.config.takerFee;

        if (fs.userFeeDiscount[user] > 0) {
            effectiveFeeRate = (baseFee * (BASIS_POINTS - fs.userFeeDiscount[user])) / BASIS_POINTS;
        } else {
            effectiveFeeRate = baseFee;
        }

        return effectiveFeeRate;
    }

    /**
     * @notice Get fee token address
     * @return Fee token address
     */
    function getFeeToken() external view returns (address) {
        return address(feeToken);
    }

    /**
     * @notice Get protocol fee share
     * @return Protocol fee share in basis points
     */
    function getProtocolFeeShare() external view returns (uint256) {
        FeeStorage storage fs = _getStorage();
        return fs.protocolFeeShare;
    }
}
