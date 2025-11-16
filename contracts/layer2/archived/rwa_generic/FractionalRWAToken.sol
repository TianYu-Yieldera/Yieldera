// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "../interfaces/rwa/IRWACompliance.sol";

/**
 * @title FractionalRWAToken
 * @notice ERC20 token representing fractional ownership of RWA
 * @dev Compliance-checked transfers, pausable, burnable
 *
 * Key Features:
 * - Represents fractional ownership of real-world asset
 * - All transfers checked against RWACompliance contract
 * - Can be paused for emergency situations
 * - Burnable for asset redemptions
 * - Metadata links to parent RWA asset
 */
contract FractionalRWAToken is ERC20, ERC20Burnable, Pausable, AccessControl {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");

    // Parent asset information
    uint256 public immutable assetId;
    address public immutable assetFactory;

    // Compliance integration
    IRWACompliance public immutable complianceContract;

    // Supply cap (based on asset fractionalization)
    uint256 public immutable supplyCap;

    /**
     * @notice Constructor
     * @param name_ Token name
     * @param symbol_ Token symbol
     * @param assetId_ Parent RWA asset ID
     * @param supplyCap_ Maximum token supply
     * @param compliance_ Compliance contract address
     * @param admin_ Admin address
     */
    constructor(
        string memory name_,
        string memory symbol_,
        uint256 assetId_,
        uint256 supplyCap_,
        address compliance_,
        address admin_
    ) ERC20(name_, symbol_) {
        require(assetId_ > 0, "Invalid asset ID");
        require(supplyCap_ > 0, "Invalid supply cap");
        require(compliance_ != address(0), "Invalid compliance address");
        require(admin_ != address(0), "Invalid admin address");

        assetId = assetId_;
        assetFactory = msg.sender; // Factory is deployer
        supplyCap = supplyCap_;
        complianceContract = IRWACompliance(compliance_);

        _grantRole(DEFAULT_ADMIN_ROLE, admin_);
        _grantRole(MINTER_ROLE, assetFactory); // Factory can mint
        _grantRole(PAUSER_ROLE, admin_);
    }

    /**
     * @notice Mint new tokens (fractionalization)
     * @param to Recipient address
     * @param amount Amount to mint
     */
    function mint(address to, uint256 amount) external onlyRole(MINTER_ROLE) whenNotPaused {
        require(totalSupply() + amount <= supplyCap, "Exceeds supply cap");
        _mint(to, amount);
    }

    /**
     * @notice Pause token transfers
     */
    function pause() external onlyRole(PAUSER_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause token transfers
     */
    function unpause() external onlyRole(PAUSER_ROLE) {
        _unpause();
    }

    /**
     * @notice Override transfer to add compliance checks
     * @param from Sender
     * @param to Recipient
     * @param amount Transfer amount
     */
    function _update(
        address from,
        address to,
        uint256 amount
    ) internal virtual override whenNotPaused {
        // Check compliance for non-mint/burn transfers
        if (from != address(0) && to != address(0)) {
            require(
                complianceContract.canTransferTokens(from, to, assetId, amount),
                "Transfer not compliant"
            );
        }

        super._update(from, to, amount);
    }

    /**
     * @notice Get asset metadata
     * @return assetId, factory address
     */
    function getAssetInfo() external view returns (uint256, address) {
        return (assetId, assetFactory);
    }
}
