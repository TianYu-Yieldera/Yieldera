// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/**
 * @title MockComet
 * @notice Mock implementation of Compound V3 (Comet) for testing
 * @dev Simulates basic lending/borrowing functionality
 */
contract MockComet {
    using SafeERC20 for IERC20;

    address public immutable baseToken;

    // User balances
    mapping(address => uint256) public baseBalances;       // Base asset supplied
    mapping(address => uint256) public borrowBalances;     // Base asset borrowed
    mapping(address => mapping(address => uint128)) public collateralBalances; // Collateral supplied

    // Interest rates (scaled by 1e18)
    uint64 public supplyRate = 2e16;  // 2% APR
    uint64 public borrowRate = 5e16;  // 5% APR

    // Utilization (scaled by 1e18)
    uint256 public utilization = 7e17; // 70%

    // Events
    event Supply(address indexed from, address indexed dst, address indexed asset, uint256 amount);
    event Withdraw(address indexed src, address indexed to, address indexed asset, uint256 amount);

    constructor(address _baseToken) {
        require(_baseToken != address(0), "Invalid base token");
        baseToken = _baseToken;
    }

    /**
     * @notice Supply an asset
     */
    function supply(address asset, uint256 amount) external {
        require(amount > 0, "Invalid amount");

        if (asset == baseToken) {
            // Supply base asset
            IERC20(baseToken).safeTransferFrom(msg.sender, address(this), amount);
            baseBalances[msg.sender] += amount;
        } else {
            // Supply collateral
            IERC20(asset).safeTransferFrom(msg.sender, address(this), amount);
            collateralBalances[msg.sender][asset] += uint128(amount);
        }

        emit Supply(msg.sender, msg.sender, asset, amount);
    }

    /**
     * @notice Supply to a specific address
     */
    function supplyTo(address from, address dst, address asset, uint256 amount) external {
        require(amount > 0, "Invalid amount");

        if (asset == baseToken) {
            IERC20(baseToken).safeTransferFrom(from, address(this), amount);
            baseBalances[dst] += amount;
        } else {
            IERC20(asset).safeTransferFrom(from, address(this), amount);
            collateralBalances[dst][asset] += uint128(amount);
        }

        emit Supply(from, dst, asset, amount);
    }

    /**
     * @notice Withdraw an asset
     */
    function withdraw(address asset, uint256 amount) external {
        withdrawTo(msg.sender, asset, amount);
    }

    /**
     * @notice Withdraw to a specific address
     */
    function withdrawTo(address to, address asset, uint256 amount) public {
        require(amount > 0, "Invalid amount");

        if (asset == baseToken) {
            // Withdraw base asset or borrow
            if (baseBalances[msg.sender] >= amount) {
                // Normal withdrawal
                baseBalances[msg.sender] -= amount;
                IERC20(baseToken).safeTransfer(to, amount);
            } else {
                // Borrowing (withdrawal when no supply)
                borrowBalances[msg.sender] += amount;
                IERC20(baseToken).safeTransfer(to, amount);
            }
        } else {
            // Withdraw collateral
            require(collateralBalances[msg.sender][asset] >= amount, "Insufficient collateral");
            collateralBalances[msg.sender][asset] -= uint128(amount);
            IERC20(asset).safeTransfer(to, amount);
        }

        emit Withdraw(msg.sender, to, asset, amount);
    }

    /**
     * @notice Get collateral balance
     */
    function collateralBalanceOf(address account, address asset) external view returns (uint128) {
        return collateralBalances[account][asset];
    }

    /**
     * @notice Get borrow balance
     */
    function borrowBalanceOf(address account) external view returns (uint256) {
        return borrowBalances[account];
    }

    /**
     * @notice Get base balance (supply)
     */
    function balanceOf(address account) external view returns (uint256) {
        return baseBalances[account];
    }

    /**
     * @notice Get supply rate
     */
    function getSupplyRate(uint256 /* utilization */) external view returns (uint64) {
        return supplyRate;
    }

    /**
     * @notice Get borrow rate
     */
    function getBorrowRate(uint256 /* utilization */) external view returns (uint64) {
        return borrowRate;
    }

    /**
     * @notice Get utilization
     */
    function getUtilization() external view returns (uint256) {
        return utilization;
    }

    /**
     * @notice Check if account is liquidatable (simplified)
     */
    function isLiquidatable(address account) external view returns (bool) {
        // Simplified: account is liquidatable if borrow > collateral value
        // In reality, this would involve price feeds and collateral factors
        return borrowBalances[account] > baseBalances[account];
    }

    /**
     * @notice Accrue account (no-op in mock)
     */
    function accrueAccount(address /* account */) external pure {
        // No-op in mock
    }

    /**
     * @notice Absorb underwater accounts (no-op in mock)
     */
    function absorb(address /* absorber */, address[] calldata /* accounts */) external pure {
        // No-op in mock
    }

    /**
     * @notice Get asset info (stub implementation)
     */
    function getAssetInfo(uint8 /* offset */)
        external
        view
        returns (
            uint8,
            address asset,
            address priceFeed,
            uint64 scale,
            uint64 borrowCollateralFactor,
            uint64 liquidateCollateralFactor,
            uint64 liquidationFactor,
            uint128 supplyCap
        )
    {
        return (0, address(0), address(0), 0, 0, 0, 0, 0);
    }

    /**
     * @notice Set supply rate (testing only)
     */
    function setSupplyRate(uint64 rate) external {
        supplyRate = rate;
    }

    /**
     * @notice Set borrow rate (testing only)
     */
    function setBorrowRate(uint64 rate) external {
        borrowRate = rate;
    }

    /**
     * @notice Set utilization (testing only)
     */
    function setUtilization(uint256 util) external {
        utilization = util;
    }
}

/**
 * @title MockCometRewards
 * @notice Mock implementation of Compound V3 rewards
 */
contract MockCometRewards {
    address public immutable rewardToken;

    mapping(address => mapping(address => uint256)) public rewardsOwed;

    event RewardClaimed(address indexed src, address indexed recipient, address indexed token, uint256 amount);

    constructor(address _rewardToken) {
        rewardToken = _rewardToken;
    }

    /**
     * @notice Claim rewards
     */
    function claim(address /* comet */, address src, bool /* shouldAccrue */) external {
        uint256 owed = rewardsOwed[msg.sender][src];
        if (owed > 0) {
            rewardsOwed[msg.sender][src] = 0;
            IERC20(rewardToken).transfer(src, owed);
            emit RewardClaimed(src, src, rewardToken, owed);
        }
    }

    /**
     * @notice Get reward owed
     */
    function getRewardOwed(address comet, address account)
        external
        view
        returns (address token, uint256 owed)
    {
        return (rewardToken, rewardsOwed[comet][account]);
    }

    /**
     * @notice Set rewards owed (testing only)
     */
    function setRewardsOwed(address comet, address account, uint256 amount) external {
        rewardsOwed[comet][account] = amount;
    }
}
