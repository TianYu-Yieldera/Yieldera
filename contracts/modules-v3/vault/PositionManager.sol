// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../../interfaces/modules/vault/IPositionManager.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title PositionManager
 * @notice Manages user positions and health calculations
 * @dev Tracks position state and calculates health metrics
 */
contract PositionManager is IPositionManager, Ownable {
    // ============ Storage ============

    bytes32 private constant POSITION_STORAGE_POSITION =
        keccak256("position.manager.storage");

    struct PositionStorage {
        mapping(address => Position) positions;
        address[] activePositions;
        mapping(address => bool) hasPosition;
        mapping(address => uint256) positionIndex; // Index in activePositions array
    }

    address public vaultModule;

    // ============ Modifiers ============

    modifier onlyVaultModule() {
        require(msg.sender == vaultModule, "Only vault module");
        _;
    }

    // ============ Constructor ============

    constructor() {}

    // ============ Admin Functions ============

    function setVaultModule(address _vaultModule) external onlyOwner {
        require(_vaultModule != address(0), "Invalid address");
        vaultModule = _vaultModule;
    }

    // ============ Internal Storage ============

    function _getStorage() private pure returns (PositionStorage storage ps) {
        bytes32 position = POSITION_STORAGE_POSITION;
        assembly {
            ps.slot := position
        }
    }

    // ============ IPositionManager Implementation ============

    function getPosition(address user) external view override returns (Position memory) {
        PositionStorage storage ps = _getStorage();
        return ps.positions[user];
    }

    function updatePosition(address user, uint256 collateralAmount, uint256 debtAmount)
        external
        override
        onlyVaultModule
    {
        PositionStorage storage ps = _getStorage();
        Position storage pos = ps.positions[user];

        pos.collateralAmount = collateralAmount;
        pos.debtAmount = debtAmount;
        pos.lastInterestUpdate = block.timestamp;
        pos.isActive = (collateralAmount > 0 || debtAmount > 0);

        // Add to active positions if not already there
        if (pos.isActive && !ps.hasPosition[user]) {
            ps.hasPosition[user] = true;
            ps.positionIndex[user] = ps.activePositions.length;
            ps.activePositions.push(user);
        }

        emit PositionUpdated(user, collateralAmount, debtAmount);
    }

    function getPositionHealth(
        address user,
        uint256 minCollateralRatio,
        uint256 liquidationThreshold
    ) external view override returns (PositionHealth memory health) {
        PositionStorage storage ps = _getStorage();
        Position storage pos = ps.positions[user];

        if (pos.debtAmount == 0) {
            return PositionHealth({
                collateralRatio: type(uint256).max,
                isHealthy: true,
                canBeLiquidated: false,
                maxWithdrawable: pos.collateralAmount,
                maxMintable: type(uint256).max
            });
        }

        uint256 ratio = calculateCollateralRatio(pos.collateralAmount, pos.debtAmount);

        health.collateralRatio = ratio;
        health.isHealthy = ratio >= minCollateralRatio;
        health.canBeLiquidated = ratio < liquidationThreshold;

        // Calculate max withdrawable (maintaining min ratio)
        if (ratio > minCollateralRatio) {
            uint256 requiredCollateral = (pos.debtAmount * minCollateralRatio) / 100;
            health.maxWithdrawable = pos.collateralAmount > requiredCollateral
                ? pos.collateralAmount - requiredCollateral
                : 0;
        } else {
            health.maxWithdrawable = 0;
        }

        // Calculate max mintable
        uint256 maxDebt = (pos.collateralAmount * 100) / minCollateralRatio;
        health.maxMintable = maxDebt > pos.debtAmount ? maxDebt - pos.debtAmount : 0;

        return health;
    }

    function calculateCollateralRatio(uint256 collateral, uint256 debt)
        public
        pure
        override
        returns (uint256)
    {
        if (debt == 0) return type(uint256).max;
        return (collateral * 100) / debt;
    }

    function isPositionHealthy(uint256 collateral, uint256 debt, uint256 minRatio)
        public
        pure
        override
        returns (bool)
    {
        if (debt == 0) return true;
        return calculateCollateralRatio(collateral, debt) >= minRatio;
    }

    function getActivePositions() external view override returns (address[] memory) {
        PositionStorage storage ps = _getStorage();
        return ps.activePositions;
    }

    function getActivePositionCount() external view override returns (uint256) {
        PositionStorage storage ps = _getStorage();

        uint256 count = 0;
        for (uint256 i = 0; i < ps.activePositions.length; i++) {
            Position storage pos = ps.positions[ps.activePositions[i]];
            if (pos.isActive) {
                count++;
            }
        }
        return count;
    }

    function hasActivePosition(address user) external view override returns (bool) {
        PositionStorage storage ps = _getStorage();
        return ps.positions[user].isActive;
    }

    function openPosition(address user) external override onlyVaultModule {
        PositionStorage storage ps = _getStorage();

        if (!ps.hasPosition[user]) {
            ps.hasPosition[user] = true;
            ps.positionIndex[user] = ps.activePositions.length;
            ps.activePositions.push(user);

            ps.positions[user] = Position({
                collateralAmount: 0,
                debtAmount: 0,
                lastInterestUpdate: block.timestamp,
                accruedInterest: 0,
                isActive: true
            });

            emit PositionOpened(user);
        }
    }

    function closePosition(address user) external override onlyVaultModule {
        PositionStorage storage ps = _getStorage();
        Position storage pos = ps.positions[user];

        require(pos.collateralAmount == 0 && pos.debtAmount == 0, "Position not empty");

        pos.isActive = false;

        emit PositionClosed(user);
    }
}
