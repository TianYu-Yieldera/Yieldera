// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/modules/vault/IInterestCalculator.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title InterestCalculator
 * @notice Calculates and manages interest accrual on debt
 * @dev Uses compound interest calculation with stability fee
 */
contract InterestCalculator is IInterestCalculator, Ownable {
    // ============ Constants ============

    uint256 private constant SECONDS_PER_YEAR = 365 days;
    uint256 private constant BASIS_POINTS = 10000; // 100% = 10000 basis points
    uint256 private constant PRECISION = 1e18;

    // ============ State Variables ============

    address public vaultModule; // Main coordinator contract

    // Stability fee in basis points (e.g., 200 = 2% annual)
    uint256 public stabilityFee;

    // Storage for accrued interest tracking
    bytes32 private constant INTEREST_STORAGE_POSITION = keccak256("interest.calculator.storage");

    struct InterestStorage {
        mapping(address => uint256) accruedInterest;
        mapping(address => uint256) lastAccrual;
    }

    // ============ Modifiers ============

    modifier onlyVaultModule() {
        require(msg.sender == vaultModule, "Only vault module");
        _;
    }

    // ============ Constructor ============

    constructor(uint256 _initialStabilityFee) {
        stabilityFee = _initialStabilityFee;
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

    function _getStorage() private pure returns (InterestStorage storage is_) {
        bytes32 position = INTEREST_STORAGE_POSITION;
        assembly {
            is_.slot := position
        }
    }

    // ============ IInterestCalculator Implementation ============

    /**
     * @notice Calculate accrued interest for a user
     * @param user User address (for future use in custom rates)
     * @param principal Principal debt amount
     * @param lastUpdate Last interest update timestamp
     * @return Accrued interest amount
     */
    function calculateInterest(address user, uint256 principal, uint256 lastUpdate)
        external
        view
        override
        returns (uint256)
    {
        // Avoid unused parameter warning
        user;

        if (principal == 0 || lastUpdate >= block.timestamp) {
            return 0;
        }

        uint256 duration = block.timestamp - lastUpdate;
        return calculateInterestForPeriod(principal, duration);
    }

    /**
     * @notice Accrue interest for a user
     * @param user User address
     * @param principal Principal debt amount
     * @param lastUpdate Last interest update timestamp
     * @return newDebt New total debt including interest
     */
    function accrueInterest(address user, uint256 principal, uint256 lastUpdate)
        external
        override
        onlyVaultModule
        returns (uint256 newDebt)
    {
        if (principal == 0 || lastUpdate >= block.timestamp) {
            return principal;
        }

        uint256 duration = block.timestamp - lastUpdate;
        uint256 interest = calculateInterestForPeriod(principal, duration);

        InterestStorage storage is_ = _getStorage();
        is_.accruedInterest[user] += interest;
        is_.lastAccrual[user] = block.timestamp;

        newDebt = principal + interest;

        emit InterestAccrued(user, interest, newDebt);

        return newDebt;
    }

    /**
     * @notice Get current stability fee (annual interest rate)
     * @return Stability fee in basis points
     */
    function getStabilityFee() external view override returns (uint256) {
        return stabilityFee;
    }

    /**
     * @notice Set stability fee
     * @param newFee New fee in basis points
     */
    function setStabilityFee(uint256 newFee) external override onlyOwner {
        require(newFee <= 5000, "Fee too high"); // Max 50% annual
        uint256 oldFee = stabilityFee;
        stabilityFee = newFee;
        emit StabilityFeeUpdated(oldFee, newFee);
    }

    /**
     * @notice Calculate interest for a time period
     * @param principal Principal amount
     * @param duration Duration in seconds
     * @return Interest amount
     */
    function calculateInterestForPeriod(uint256 principal, uint256 duration)
        public
        view
        override
        returns (uint256)
    {
        if (principal == 0 || duration == 0 || stabilityFee == 0) {
            return 0;
        }

        // Simple interest calculation: interest = principal * rate * time
        // rate = stabilityFee / BASIS_POINTS / SECONDS_PER_YEAR
        // interest = principal * (stabilityFee / BASIS_POINTS) * (duration / SECONDS_PER_YEAR)

        uint256 interest = (principal * stabilityFee * duration) / (BASIS_POINTS * SECONDS_PER_YEAR);

        return interest;
    }

    /**
     * @notice Get compounded amount
     * @dev For small time periods, uses linear approximation
     * @param principal Principal amount
     * @param duration Duration in seconds
     * @return Compounded amount
     */
    function getCompoundedAmount(uint256 principal, uint256 duration)
        external
        view
        override
        returns (uint256)
    {
        if (principal == 0 || duration == 0 || stabilityFee == 0) {
            return principal;
        }

        // For simplicity, use linear calculation
        // In production, this could use more sophisticated compounding
        uint256 interest = calculateInterestForPeriod(principal, duration);
        return principal + interest;
    }

    /**
     * @notice Get accrued interest for a user
     * @param user User address
     * @return Total accrued interest
     */
    function getAccruedInterest(address user) external view returns (uint256) {
        InterestStorage storage is_ = _getStorage();
        return is_.accruedInterest[user];
    }

    /**
     * @notice Get last accrual timestamp for a user
     * @param user User address
     * @return Last accrual timestamp
     */
    function getLastAccrual(address user) external view returns (uint256) {
        InterestStorage storage is_ = _getStorage();
        return is_.lastAccrual[user];
    }
}
