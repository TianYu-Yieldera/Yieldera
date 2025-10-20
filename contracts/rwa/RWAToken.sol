// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

/**
 * @title RWAToken
 * @notice ERC20 token representing tokenized Real World Assets
 * @dev Supports fractional ownership of real-world assets with compliance features
 */
contract RWAToken is ERC20, ERC20Burnable, AccessControl, Pausable {
    // Roles
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant COMPLIANCE_ROLE = keccak256("COMPLIANCE_ROLE");
    bytes32 public constant TRANSFER_ROLE = keccak256("TRANSFER_ROLE");

    // Asset information
    struct AssetInfo {
        string assetType; // e.g., "Real Estate", "Commodity", "Security"
        string assetId; // Unique identifier
        uint256 totalValue; // Total value in USD
        uint256 fractionSize; // Size of each token fraction
        string documentHash; // IPFS hash of legal documents
        bool isActive;
    }

    AssetInfo public assetInfo;

    // Compliance
    mapping(address => bool) public whitelist;
    mapping(address => bool) public blacklist;
    bool public whitelistEnabled = true;

    // Transfer restrictions
    uint256 public minTransferAmount = 1e18; // 1 token minimum
    uint256 public maxTransferAmount = 1000000e18; // 1M tokens maximum
    uint256 public dailyTransferLimit = 10000000e18; // 10M tokens daily
    mapping(address => uint256) public lastTransferDay;
    mapping(address => uint256) public dailyTransferred;

    // Events
    event AssetInfoUpdated(string assetType, string assetId, uint256 totalValue);
    event AddressWhitelisted(address indexed account);
    event AddressBlacklisted(address indexed account);
    event ComplianceCheckUpdated(bool whitelistEnabled);
    event TransferLimitsUpdated(uint256 min, uint256 max, uint256 daily);

    constructor(
        string memory name,
        string memory symbol,
        string memory _assetType,
        string memory _assetId,
        uint256 _totalValue
    ) ERC20(name, symbol) {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);
        _grantRole(COMPLIANCE_ROLE, msg.sender);

        assetInfo = AssetInfo({
            assetType: _assetType,
            assetId: _assetId,
            totalValue: _totalValue,
            fractionSize: 1e18, // 1 token = 1 fraction
            documentHash: "",
            isActive: true
        });
    }

    /**
     * @notice Mint new RWA tokens
     * @param to Recipient address
     * @param amount Amount to mint
     */
    function mint(address to, uint256 amount) external onlyRole(MINTER_ROLE) {
        require(assetInfo.isActive, "Asset not active");
        _mint(to, amount);
    }

    /**
     * @notice Update asset information
     */
    function updateAssetInfo(
        string memory _assetType,
        string memory _assetId,
        uint256 _totalValue,
        string memory _documentHash
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        assetInfo.assetType = _assetType;
        assetInfo.assetId = _assetId;
        assetInfo.totalValue = _totalValue;
        assetInfo.documentHash = _documentHash;

        emit AssetInfoUpdated(_assetType, _assetId, _totalValue);
    }

    /**
     * @notice Add address to whitelist
     */
    function addToWhitelist(address account) external onlyRole(COMPLIANCE_ROLE) {
        whitelist[account] = true;
        emit AddressWhitelisted(account);
    }

    /**
     * @notice Remove from whitelist
     */
    function removeFromWhitelist(address account) external onlyRole(COMPLIANCE_ROLE) {
        whitelist[account] = false;
    }

    /**
     * @notice Add to blacklist
     */
    function addToBlacklist(address account) external onlyRole(COMPLIANCE_ROLE) {
        blacklist[account] = true;
        emit AddressBlacklisted(account);
    }

    /**
     * @notice Remove from blacklist
     */
    function removeFromBlacklist(address account) external onlyRole(COMPLIANCE_ROLE) {
        blacklist[account] = false;
    }

    /**
     * @notice Set transfer limits
     */
    function setTransferLimits(
        uint256 _min,
        uint256 _max,
        uint256 _daily
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        minTransferAmount = _min;
        maxTransferAmount = _max;
        dailyTransferLimit = _daily;

        emit TransferLimitsUpdated(_min, _max, _daily);
    }

    /**
     * @notice Check if transfer is compliant
     */
    function isTransferCompliant(
        address from,
        address to,
        uint256 amount
    ) public view returns (bool) {
        // Check blacklist
        if (blacklist[from] || blacklist[to]) {
            return false;
        }

        // Check whitelist if enabled
        if (whitelistEnabled && (!whitelist[from] || !whitelist[to])) {
            return false;
        }

        // Check transfer limits
        if (amount < minTransferAmount || amount > maxTransferAmount) {
            return false;
        }

        // Check daily limit
        uint256 currentDay = block.timestamp / 86400;
        if (lastTransferDay[from] == currentDay) {
            if (dailyTransferred[from] + amount > dailyTransferLimit) {
                return false;
            }
        }

        return true;
    }

    /**
     * @notice Override transfer to include compliance checks
     */
    function _transfer(
        address from,
        address to,
        uint256 amount
    ) internal override whenNotPaused {
        require(isTransferCompliant(from, to, amount), "Transfer not compliant");

        // Update daily transfer tracking
        uint256 currentDay = block.timestamp / 86400;
        if (lastTransferDay[from] != currentDay) {
            lastTransferDay[from] = currentDay;
            dailyTransferred[from] = 0;
        }
        dailyTransferred[from] += amount;

        super._transfer(from, to, amount);
    }

    /**
     * @notice Get token value in USD
     */
    function getTokenValueUSD() external view returns (uint256) {
        if (totalSupply() == 0) return 0;
        return assetInfo.totalValue / totalSupply();
    }

    /**
     * @notice Pause transfers
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause transfers
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }

    /**
     * @notice Enable/disable whitelist
     */
    function setWhitelistEnabled(bool enabled) external onlyRole(COMPLIANCE_ROLE) {
        whitelistEnabled = enabled;
        emit ComplianceCheckUpdated(enabled);
    }
}