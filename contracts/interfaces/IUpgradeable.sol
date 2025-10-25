// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IUpgradeable
 * @notice Interface for upgradeable modules using UUPS pattern
 * @dev Extends ERC1967 UUPS standard with additional safety features
 */
interface IUpgradeable {
    /**
     * @notice Emitted when implementation is upgraded
     * @param oldImplementation Previous implementation address
     * @param newImplementation New implementation address
     * @param version New version string
     */
    event Upgraded(
        address indexed oldImplementation,
        address indexed newImplementation,
        string version
    );

    /**
     * @notice Emitted when upgrade is authorized
     * @param implementation Implementation address
     * @param authorizer Address that authorized the upgrade
     */
    event UpgradeAuthorized(address indexed implementation, address indexed authorizer);

    /**
     * @notice Emitted when upgrade authorization is revoked
     * @param implementation Implementation address
     */
    event UpgradeAuthorizationRevoked(address indexed implementation);

    /**
     * @notice Upgrade to a new implementation
     * @param newImplementation Address of the new implementation
     * @dev Can only be called by authorized upgraders
     */
    function upgradeTo(address newImplementation) external;

    /**
     * @notice Upgrade to a new implementation and call a function
     * @param newImplementation Address of the new implementation
     * @param data Function call data
     * @dev Can only be called by authorized upgraders
     */
    function upgradeToAndCall(address newImplementation, bytes memory data) external payable;

    /**
     * @notice Get the current implementation address
     * @return Implementation contract address
     */
    function getImplementation() external view returns (address);

    /**
     * @notice Get the current implementation version
     * @return Version string
     */
    function getImplementationVersion() external view returns (string memory);

    /**
     * @notice Check if an address is authorized to upgrade
     * @param account Address to check
     * @return True if authorized
     */
    function canUpgrade(address account) external view returns (bool);

    /**
     * @notice Authorize an address to perform upgrades
     * @param account Address to authorize
     * @dev Only callable by admin
     */
    function authorizeUpgrader(address account) external;

    /**
     * @notice Revoke upgrade authorization
     * @param account Address to revoke
     * @dev Only callable by admin
     */
    function revokeUpgradeAuthorization(address account) external;

    /**
     * @notice Validate a new implementation before upgrade
     * @param newImplementation Address to validate
     * @return valid True if implementation is valid
     * @return reason Reason if invalid
     */
    function validateUpgrade(address newImplementation)
        external
        view
        returns (bool valid, string memory reason);

    /**
     * @notice Perform pre-upgrade checks and preparation
     * @param newImplementation Address of new implementation
     * @dev Called before upgrade, can revert to prevent upgrade
     */
    function beforeUpgrade(address newImplementation) external;

    /**
     * @notice Perform post-upgrade initialization
     * @param oldImplementation Address of previous implementation
     * @dev Called after upgrade completes
     */
    function afterUpgrade(address oldImplementation) external;

    /**
     * @notice Get upgrade history
     * @return implementations Array of implementation addresses
     * @return timestamps Array of upgrade timestamps
     * @return versions Array of version strings
     */
    function getUpgradeHistory()
        external
        view
        returns (
            address[] memory implementations,
            uint256[] memory timestamps,
            string[] memory versions
        );

    /**
     * @notice Emergency pause upgrades
     * @dev Can only be called by admin
     */
    function pauseUpgrades() external;

    /**
     * @notice Resume upgrades
     * @dev Can only be called by admin
     */
    function resumeUpgrades() external;

    /**
     * @notice Check if upgrades are paused
     * @return True if paused
     */
    function upgradesPaused() external view returns (bool);
}
