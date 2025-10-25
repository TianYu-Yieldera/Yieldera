// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Utils.sol";

/**
 * @title ModuleProxy
 * @notice UUPS Proxy for pluggable modules
 * @dev Extends OpenZeppelin's ERC1967Proxy with module-specific features
 *
 * Features:
 * - UUPS upgrade pattern (upgrade logic in implementation)
 * - ERC1967 storage slots for implementation and admin
 * - Minimal proxy overhead
 * - Transparent to users
 */
contract ModuleProxy is ERC1967Proxy {
    /**
     * @notice Emitted when proxy is deployed
     * @param implementation Initial implementation address
     * @param admin Proxy admin address
     */
    event ProxyDeployed(address indexed implementation, address indexed admin);

    /**
     * @notice Constructor
     * @param _logic Address of the initial implementation
     * @param _data Initialization data to call on implementation
     * @dev Sets up ERC1967 proxy with initial implementation
     */
    constructor(address _logic, bytes memory _data) ERC1967Proxy(_logic, _data) {
        // Get admin from implementation if it supports IUpgradeable
        address admin = msg.sender;

        emit ProxyDeployed(_logic, admin);
    }

    /**
     * @notice Returns the current implementation address
     * @return Implementation contract address
     * @dev Uses ERC1967 storage slot
     */
    function implementation() external view returns (address) {
        return _implementation();
    }

    /**
     * @notice Returns the current admin address
     * @return Admin address
     * @dev Uses ERC1967 storage slot
     */
    function admin() external view returns (address) {
        return ERC1967Utils.getAdmin();
    }

    /**
     * @notice Change the admin of the proxy
     * @param newAdmin Address of the new admin
     * @dev Only callable by current admin
     */
    function changeAdmin(address newAdmin) external {
        require(msg.sender == ERC1967Utils.getAdmin(), "Only admin");
        ERC1967Utils.changeAdmin(newAdmin);
    }
}
