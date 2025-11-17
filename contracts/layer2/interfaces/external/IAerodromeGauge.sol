// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IAerodromeGauge
 * @notice Interface for Aerodrome liquidity gauge (staking)
 */
interface IAerodromeGauge {
    function stakingToken() external view returns (address);
    function rewardToken() external view returns (address);

    function deposit(uint256 amount) external;
    function withdraw(uint256 amount) external;

    function getReward(address account) external;
    function earned(address account) external view returns (uint256);

    function balanceOf(address account) external view returns (uint256);
    function totalSupply() external view returns (uint256);

    function rewardRate() external view returns (uint256);
    function periodFinish() external view returns (uint256);
}
