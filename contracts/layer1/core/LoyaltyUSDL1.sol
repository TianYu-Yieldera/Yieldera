// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

/**
 * @title LoyaltyUSDL1
 * @notice Layer 1 optimized version - Only minting/burning authority
 * @dev Stripped down version for L1, business logic moved to L2
 *
 * L1 Responsibilities:
 * - Maintain minting/burning authority
 * - Final settlement of L2 state
 * - Emergency pause capability
 * - Cross-chain bridge integration
 *
 * L2 handles:
 * - All business logic
 * - Interest calculations
 * - Collateral ratio checks
 * - User interactions
 */
contract LoyaltyUSDL1 is ERC20, ERC20Burnable, AccessControl, Pausable {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant BURNER_ROLE = keccak256("BURNER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    bytes32 public constant BRIDGE_ROLE = keccak256("BRIDGE_ROLE");

    // L2 Bridge address
    address public l2Bridge;

    // Minting limits for security
    uint256 public constant MAX_MINT_PER_TX = 1_000_000 * 1e6; // 1M LUSD
    uint256 public dailyMintLimit = 10_000_000 * 1e6; // 10M LUSD per day
    uint256 public dailyMintedAmount;
    uint256 public lastMintResetTime;

    // Events
    event Minted(address indexed to, uint256 amount, address indexed minter);
    event Burned(address indexed from, uint256 amount, address indexed burner);
    event EmergencyPaused(address indexed pauser);
    event EmergencyUnpaused(address indexed unpauser);
    event L2BridgeUpdated(address indexed oldBridge, address indexed newBridge);
    event DailyMintLimitUpdated(uint256 oldLimit, uint256 newLimit);
    event BridgeMint(address indexed to, uint256 amount, bytes32 indexed l2TxHash);
    event BridgeBurn(address indexed from, uint256 amount, bytes32 indexed l2TxHash);

    /**
     * @notice Constructor initializes L1 token with admin
     */
    constructor() ERC20("LoyaltyUSD", "LUSD") {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(PAUSER_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);
        _grantRole(BURNER_ROLE, msg.sender);

        lastMintResetTime = block.timestamp;
    }

    /**
     * @notice Mint new LUSD tokens (L1 authority)
     * @param to Recipient address
     * @param amount Amount to mint (in 6 decimals)
     * @dev Only callable by MINTER_ROLE or BRIDGE_ROLE
     */
    function mint(address to, uint256 amount) external onlyRole(MINTER_ROLE) {
        require(to != address(0), "LUSD: mint to zero address");
        require(amount > 0, "LUSD: mint amount must be positive");
        require(amount <= MAX_MINT_PER_TX, "LUSD: exceeds max mint per tx");

        // Check daily limit
        _checkAndUpdateDailyLimit(amount);

        _mint(to, amount);
        emit Minted(to, amount, msg.sender);
    }

    /**
     * @notice Mint from L2 bridge (cross-chain)
     * @param to Recipient address
     * @param amount Amount to mint
     * @param l2TxHash L2 transaction hash for tracking
     */
    function bridgeMint(address to, uint256 amount, bytes32 l2TxHash)
        external
        onlyRole(BRIDGE_ROLE)
    {
        require(to != address(0), "LUSD: mint to zero address");
        require(amount > 0, "LUSD: amount must be positive");
        require(amount <= MAX_MINT_PER_TX, "LUSD: exceeds max mint per tx");

        // Check daily limit
        _checkAndUpdateDailyLimit(amount);

        _mint(to, amount);
        emit BridgeMint(to, amount, l2TxHash);
    }

    /**
     * @notice Burn LUSD tokens from an address
     * @param from Address to burn from
     * @param amount Amount to burn
     * @dev Only callable by BURNER_ROLE or token owner
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
     * @notice Burn from L2 bridge (cross-chain)
     * @param from Address to burn from
     * @param amount Amount to burn
     * @param l2TxHash L2 transaction hash for tracking
     */
    function bridgeBurn(address from, uint256 amount, bytes32 l2TxHash)
        external
        onlyRole(BRIDGE_ROLE)
    {
        require(amount > 0, "LUSD: amount must be positive");

        _burn(from, amount);
        emit BridgeBurn(from, amount, l2TxHash);
    }

    /**
     * @notice Check and update daily mint limit
     * @dev Internal function to prevent over-minting
     */
    function _checkAndUpdateDailyLimit(uint256 amount) internal {
        // Reset daily counter if 24 hours passed
        if (block.timestamp >= lastMintResetTime + 1 days) {
            dailyMintedAmount = 0;
            lastMintResetTime = block.timestamp;
        }

        require(
            dailyMintedAmount + amount <= dailyMintLimit,
            "LUSD: exceeds daily mint limit"
        );

        dailyMintedAmount += amount;
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
     * @dev Overridden to add pausable functionality (OpenZeppelin 5.x uses _update)
     */
    function _update(
        address from,
        address to,
        uint256 amount
    ) internal override whenNotPaused {
        super._update(from, to, amount);
    }

    /**
     * @notice Set L2 bridge address
     * @param _l2Bridge L2 bridge address
     */
    function setL2Bridge(address _l2Bridge) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(_l2Bridge != address(0), "Invalid bridge address");
        address oldBridge = l2Bridge;
        l2Bridge = _l2Bridge;

        // Grant bridge role
        _grantRole(BRIDGE_ROLE, _l2Bridge);

        // Revoke old bridge role if exists
        if (oldBridge != address(0)) {
            _revokeRole(BRIDGE_ROLE, oldBridge);
        }

        emit L2BridgeUpdated(oldBridge, _l2Bridge);
    }

    /**
     * @notice Update daily mint limit
     * @param newLimit New daily limit
     */
    function setDailyMintLimit(uint256 newLimit)
        external
        onlyRole(DEFAULT_ADMIN_ROLE)
    {
        uint256 oldLimit = dailyMintLimit;
        dailyMintLimit = newLimit;
        emit DailyMintLimitUpdated(oldLimit, newLimit);
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
        return hasRole(MINTER_ROLE, account) || hasRole(BRIDGE_ROLE, account);
    }

    /**
     * @notice Check if address has burning capability
     */
    function canBurn(address account) external view returns (bool) {
        return hasRole(BURNER_ROLE, account) || hasRole(BRIDGE_ROLE, account);
    }

    /**
     * @notice Get remaining daily mint capacity
     */
    function getRemainingDailyMintCapacity() external view returns (uint256) {
        // Check if we need to reset
        if (block.timestamp >= lastMintResetTime + 1 days) {
            return dailyMintLimit;
        }

        if (dailyMintedAmount >= dailyMintLimit) {
            return 0;
        }

        return dailyMintLimit - dailyMintedAmount;
    }

    /**
     * @notice Get time until daily limit resets
     */
    function getTimeUntilLimitReset() external view returns (uint256) {
        uint256 resetTime = lastMintResetTime + 1 days;
        if (block.timestamp >= resetTime) {
            return 0;
        }
        return resetTime - block.timestamp;
    }
}
