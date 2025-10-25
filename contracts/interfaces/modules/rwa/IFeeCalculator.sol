// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IOrderManager.sol";

/**
 * @title IFeeCalculator
 * @notice Interface for trading fee calculations
 * @dev Handles fee computation and collection
 */
interface IFeeCalculator {
    // ============ Structs ============

    struct FeeConfig {
        uint256 makerFee; // Fee for makers (basis points)
        uint256 takerFee; // Fee for takers (basis points)
        uint256 minimumFee; // Minimum fee amount
        bool feesEnabled;
    }

    struct FeeBreakdown {
        uint256 makerFee;
        uint256 takerFee;
        uint256 protocolFee;
        uint256 totalFee;
    }

    // ============ Events ============

    event FeeCollected(
        address indexed user, uint256 amount, bool isMaker, uint256 tradeAmount, uint256 feeRate
    );
    event FeeConfigUpdated(uint256 makerFee, uint256 takerFee, uint256 minimumFee);
    event FeesWithdrawn(address indexed recipient, uint256 amount);

    // ============ Functions ============

    /**
     * @notice Calculate trading fee
     * @param tradeAmount Trade amount
     * @param isMaker True if user is maker, false if taker
     * @return fee Fee amount
     */
    function calculateFee(uint256 tradeAmount, bool isMaker) external view returns (uint256 fee);

    /**
     * @notice Calculate fee breakdown for a trade
     * @param tradeAmount Trade amount
     * @param buyOrderId Buy order ID
     * @param sellOrderId Sell order ID
     * @return breakdown Fee breakdown
     */
    function calculateFeeBreakdown(uint256 tradeAmount, uint256 buyOrderId, uint256 sellOrderId)
        external
        view
        returns (FeeBreakdown memory breakdown);

    /**
     * @notice Collect fee from a trade
     * @param user User paying the fee
     * @param tradeAmount Trade amount
     * @param isMaker True if user is maker
     * @return feeAmount Fee collected
     */
    function collectFee(address user, uint256 tradeAmount, bool isMaker)
        external
        returns (uint256 feeAmount);

    /**
     * @notice Get current fee configuration
     * @return config Fee configuration
     */
    function getFeeConfig() external view returns (FeeConfig memory config);

    /**
     * @notice Update fee configuration
     * @param makerFee New maker fee (basis points)
     * @param takerFee New taker fee (basis points)
     * @param minimumFee New minimum fee
     */
    function updateFeeConfig(uint256 makerFee, uint256 takerFee, uint256 minimumFee) external;

    /**
     * @notice Get total fees collected
     * @return Total fees accumulated
     */
    function getTotalFeesCollected() external view returns (uint256);

    /**
     * @notice Get fees collected from a user
     * @param user User address
     * @return Total fees paid by user
     */
    function getUserFeesCollected(address user) external view returns (uint256);

    /**
     * @notice Check if user qualifies for fee discount
     * @param user User address
     * @return hasDiscount True if user has discount
     * @return discountRate Discount rate (basis points)
     */
    function getFeeDiscount(address user)
        external
        view
        returns (bool hasDiscount, uint256 discountRate);

    /**
     * @notice Withdraw collected fees
     * @param recipient Recipient address
     * @param amount Amount to withdraw
     */
    function withdrawFees(address recipient, uint256 amount) external;
}
