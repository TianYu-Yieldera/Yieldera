// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

/**
 * @title LoyaltyUSD
 * @notice USD-pegged stablecoin backed by Loyalty Points
 * @dev ERC-20 token with minting/burning controlled by StabilityManager
 *
 * Features:
 * - 1:1 peg with USD
 * - Minting requires MINTER_ROLE (StabilityManager contract)
 * - Burning requires BURNER_ROLE (StabilityManager or user)
 * - 6 decimals to match USDC
 * - Emergency pause capability
 *
 * Roles:
 * - DEFAULT_ADMIN_ROLE: Can grant/revoke other roles
 * - MINTER_ROLE: Can mint new LUSD
 * - BURNER_ROLE: Can burn LUSD
 * - PAUSER_ROLE: Can pause/unpause transfers
 */
contract LoyaltyUSD is ERC20, ERC20Burnable, AccessControl, Pausable {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant BURNER_ROLE = keccak256("BURNER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");

    // Events
    event Minted(address indexed to, uint256 amount, address indexed minter);
    event Burned(address indexed from, uint256 amount, address indexed burner);
    event EmergencyPaused(address indexed pauser);
    event EmergencyUnpaused(address indexed unpauser);

    /**
     * @notice Constructor initializes the token with admin
     */
    constructor() ERC20("LoyaltyUSD", "LUSD") {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(PAUSER_ROLE, msg.sender);
    }

    /**
     * @notice Mint new LUSD tokens
     * @param to Recipient address
     * @param amount Amount to mint (in 6 decimals)
     * @dev Only callable by MINTER_ROLE (StabilityManager)
     */
    function mint(address to, uint256 amount) external onlyRole(MINTER_ROLE) {
        require(to != address(0), "LUSD: mint to zero address");
        require(amount > 0, "LUSD: mint amount must be positive");

        _mint(to, amount);
        emit Minted(to, amount, msg.sender);
    }

    /**
     * @notice Burn LUSD tokens from an address
     * @param from Address to burn from
     * @param amount Amount to burn
     * @dev Only callable by BURNER_ROLE (StabilityManager) or token owner
     */
    function burn(address from, uint256 amount) public {
        require(
            msg.sender == from || hasRole(BURNER_ROLE, msg.sender),
            "LUSD: must have BURNER_ROLE or be token owner"
        );
        require(amount > 0, "LUSD: burn amount must be positive");

        _burn(from, amount);
        emit Burned(from, amount, msg.sender);
    }

    /**
     * @notice Returns 6 decimals to match USDC
     */
    function decimals() public pure override returns (uint8) {
        return 6;
    }

    /**
     * @notice Pause all token transfers
     * @dev Only callable by PAUSER_ROLE
     */
    function pause() external onlyRole(PAUSER_ROLE) {
        _pause();
        emit EmergencyPaused(msg.sender);
    }

    /**
     * @notice Unpause token transfers
     * @dev Only callable by PAUSER_ROLE
     */
    function unpause() external onlyRole(PAUSER_ROLE) {
        _unpause();
        emit EmergencyUnpaused(msg.sender);
    }

    /**
     * @notice Hook that is called before any transfer of tokens
     * @dev Overridden to add pausable functionality
     */
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal override whenNotPaused {
        super._beforeTokenTransfer(from, to, amount);
    }

    /**
     * @notice Get total supply in human-readable format
     * @return Total supply with decimals applied
     */
    function totalSupplyFormatted() external view returns (uint256) {
        return totalSupply() / (10 ** decimals());
    }

    /**
     * @notice Check if address has minting capability
     */
    function canMint(address account) external view returns (bool) {
        return hasRole(MINTER_ROLE, account);
    }

    /**
     * @notice Check if address has burning capability
     */
    function canBurn(address account) external view returns (bool) {
        return hasRole(BURNER_ROLE, account);
    }
}
