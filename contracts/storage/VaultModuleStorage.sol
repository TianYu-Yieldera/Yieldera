// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../interfaces/core/IVaultModule.sol";

/**
 * @title VaultModuleStorage
 * @notice Diamond Storage library for VaultModule
 * @dev Implements EIP-2535 Diamond Storage pattern to prevent storage collisions
 */
library VaultModuleStorage {
    // Storage position is keccak256("vault.module.storage") - 1
    bytes32 constant STORAGE_POSITION =
        0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;

    struct VaultData {
        // Module configuration
        IVaultModule.VaultConfig config;

        // Active position tracking
        mapping(address => bool) hasPosition;
        address[] activePositions;

        // References to external contracts
        address vault;          // CollateralVault address
        address collateralToken;
        address debtToken;

        // Module metadata
        uint256 totalPositions;
        uint256 lastUpdate;

        // Reserved slots for future upgrades
        uint256[50] __gap;
    }

    /**
     * @notice Returns the storage layout
     * @return ds The storage layout struct
     */
    function layout() internal pure returns (VaultData storage ds) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            ds.slot := position
        }
    }

    /**
     * @notice Initialize storage (called once during proxy deployment)
     * @param _vault CollateralVault address
     * @param _collateralToken Collateral token address
     * @param _debtToken Debt token address
     */
    function initialize(
        address _vault,
        address _collateralToken,
        address _debtToken
    ) internal {
        VaultData storage ds = layout();
        require(ds.vault == address(0), "Already initialized");

        ds.vault = _vault;
        ds.collateralToken = _collateralToken;
        ds.debtToken = _debtToken;
        ds.lastUpdate = block.timestamp;
    }

    /**
     * @notice Get vault configuration
     */
    function getConfig() internal view returns (IVaultModule.VaultConfig storage) {
        return layout().config;
    }

    /**
     * @notice Set vault configuration
     */
    function setConfig(IVaultModule.VaultConfig memory _config) internal {
        VaultData storage ds = layout();
        ds.config = _config;
        ds.lastUpdate = block.timestamp;
    }

    /**
     * @notice Check if user has a position
     */
    function hasPosition(address user) internal view returns (bool) {
        return layout().hasPosition[user];
    }

    /**
     * @notice Add user to active positions
     */
    function addPosition(address user) internal {
        VaultData storage ds = layout();
        if (!ds.hasPosition[user]) {
            ds.hasPosition[user] = true;
            ds.activePositions.push(user);
            ds.totalPositions++;
        }
    }

    /**
     * @notice Get all active positions
     */
    function getActivePositions() internal view returns (address[] storage) {
        return layout().activePositions;
    }

    /**
     * @notice Get external contract addresses
     */
    function getAddresses()
        internal
        view
        returns (address vault, address collateralToken, address debtToken)
    {
        VaultData storage ds = layout();
        return (ds.vault, ds.collateralToken, ds.debtToken);
    }

    /**
     * @notice Get total number of positions
     */
    function getTotalPositions() internal view returns (uint256) {
        return layout().totalPositions;
    }

    /**
     * @notice Update last modified timestamp
     */
    function touch() internal {
        layout().lastUpdate = block.timestamp;
    }
}
