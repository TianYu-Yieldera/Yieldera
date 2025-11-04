// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "../interfaces/treasury/ITreasuryAsset.sol";
import "./TreasuryToken.sol";

/**
 * @title TreasuryAssetFactory
 * @notice Factory contract for creating and managing US Treasury tokenized assets
 * @dev Creates ERC20 tokens representing fractional ownership of US Treasuries
 *
 * Key Features:
 * - Issue T-Bills, T-Notes, and T-Bonds as ERC20 tokens
 * - Track asset metadata (CUSIP, maturity, coupon rate)
 * - Mint and burn tokens for issuance and redemption
 * - Asset lifecycle management (Active → Matured)
 */
contract TreasuryAssetFactory is ITreasuryAsset, AccessControl, ReentrancyGuard, Pausable {
    bytes32 public constant ISSUER_ROLE = keccak256("ISSUER_ROLE");
    bytes32 public constant ASSET_MANAGER_ROLE = keccak256("ASSET_MANAGER_ROLE");

    /// @notice Asset counter
    uint256 private assetIdCounter;

    /// @notice Asset storage
    mapping(uint256 => TreasuryAssetInfo) private assets;
    mapping(string => uint256) private cusipToAssetId; // CUSIP → assetId
    mapping(uint256 => address) private assetTokens;   // assetId → token address

    /// @notice Asset arrays for iteration
    uint256[] private allAssetIds;
    mapping(TreasuryType => uint256[]) private assetsByType;

    /**
     * @notice Constructor
     * @param admin Admin address
     */
    constructor(address admin) {
        require(admin != address(0), "Invalid admin");

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(ISSUER_ROLE, admin);
        _grantRole(ASSET_MANAGER_ROLE, admin);
    }

    // =============================================================
    //                     ASSET CREATION
    // =============================================================

    /**
     * @notice Create new treasury asset
     * @param treasuryType Type of treasury (T-Bill, T-Note, T-Bond)
     * @param maturityTerm Maturity term (e.g., "13W", "2Y", "10Y")
     * @param cusip CUSIP identifier
     * @param issueDate Issue timestamp
     * @param maturityDate Maturity timestamp
     * @param faceValue Face value in USD (18 decimals)
     * @param couponRate Annual coupon rate in basis points
     * @return assetId Asset identifier
     * @return tokenAddress Token contract address
     */
    function createTreasuryAsset(
        TreasuryType treasuryType,
        string memory maturityTerm,
        string memory cusip,
        uint256 issueDate,
        uint256 maturityDate,
        uint256 faceValue,
        uint256 couponRate
    ) external override onlyRole(ISSUER_ROLE) whenNotPaused nonReentrant
        returns (uint256 assetId, address tokenAddress)
    {
        require(bytes(cusip).length > 0, "Invalid CUSIP");
        require(cusipToAssetId[cusip] == 0, "CUSIP already exists");
        require(maturityDate > issueDate, "Invalid maturity date");
        require(maturityDate > block.timestamp, "Maturity in past");
        require(faceValue > 0, "Invalid face value");
        require(couponRate <= 10000, "Coupon rate > 100%");

        // Validate treasury type and maturity
        _validateTreasuryType(treasuryType, maturityTerm, maturityDate - issueDate);

        // Create asset ID
        assetIdCounter++;
        assetId = assetIdCounter;

        // Generate token name and symbol
        (string memory name, string memory symbol) = _generateTokenMetadata(
            treasuryType,
            maturityTerm,
            cusip
        );

        // Deploy token contract
        TreasuryToken token = new TreasuryToken(
            name,
            symbol,
            assetId,
            cusip,
            maturityTerm,
            maturityDate,
            couponRate,
            address(this)
        );

        tokenAddress = address(token);

        // Store asset info
        TreasuryAssetInfo storage asset = assets[assetId];
        asset.assetId = assetId;
        asset.treasuryType = treasuryType;
        asset.maturityTerm = maturityTerm;
        asset.cusip = cusip;
        asset.issueDate = issueDate;
        asset.maturityDate = maturityDate;
        asset.faceValue = faceValue;
        asset.couponRate = couponRate;
        asset.tokensIssued = 0;
        asset.tokensOutstanding = 0;
        asset.tokenAddress = tokenAddress;
        asset.status = AssetStatus.Active;
        asset.createdAt = block.timestamp;

        // Store mappings
        cusipToAssetId[cusip] = assetId;
        assetTokens[assetId] = tokenAddress;
        allAssetIds.push(assetId);
        assetsByType[treasuryType].push(assetId);

        emit TreasuryAssetCreated(assetId, treasuryType, cusip, tokenAddress, faceValue);
    }

    /**
     * @notice Validate treasury type matches maturity period
     */
    function _validateTreasuryType(
        TreasuryType treasuryType,
        string memory maturityTerm,
        uint256 maturityPeriod
    ) private pure {
        if (treasuryType == TreasuryType.T_BILL) {
            // T-Bills: 4 weeks to 52 weeks
            require(maturityPeriod <= 365 days, "T-Bill maturity > 1 year");
        } else if (treasuryType == TreasuryType.T_NOTE) {
            // T-Notes: 2 to 10 years
            require(
                maturityPeriod >= 365 days * 2 && maturityPeriod <= 365 days * 10,
                "T-Note maturity invalid"
            );
        } else if (treasuryType == TreasuryType.T_BOND) {
            // T-Bonds: 20 to 30 years
            require(
                maturityPeriod >= 365 days * 20 && maturityPeriod <= 365 days * 30,
                "T-Bond maturity invalid"
            );
        }
    }

    /**
     * @notice Generate token name and symbol
     */
    function _generateTokenMetadata(
        TreasuryType treasuryType,
        string memory maturityTerm,
        string memory cusip
    ) private pure returns (string memory name, string memory symbol) {
        string memory typeStr;

        if (treasuryType == TreasuryType.T_BILL) {
            typeStr = "Bill";
        } else if (treasuryType == TreasuryType.T_NOTE) {
            typeStr = "Note";
        } else {
            typeStr = "Bond";
        }

        name = string(abi.encodePacked("US Treasury ", maturityTerm, " ", typeStr));
        symbol = string(abi.encodePacked("UST-", maturityTerm));
    }

    // =============================================================
    //                     TOKEN OPERATIONS
    // =============================================================

    /**
     * @notice Mint treasury tokens
     * @param assetId Asset identifier
     * @param recipient Token recipient
     * @param amount Amount to mint
     */
    function mintTokens(
        uint256 assetId,
        address recipient,
        uint256 amount
    ) external override onlyRole(ISSUER_ROLE) whenNotPaused nonReentrant {
        require(assets[assetId].assetId != 0, "Asset not found");
        require(assets[assetId].status == AssetStatus.Active, "Asset not active");
        require(recipient != address(0), "Invalid recipient");
        require(amount > 0, "Invalid amount");

        TreasuryAssetInfo storage asset = assets[assetId];

        // Mint tokens
        TreasuryToken token = TreasuryToken(asset.tokenAddress);
        token.mint(recipient, amount);

        // Update counters
        asset.tokensIssued += amount;
        asset.tokensOutstanding += amount;

        emit TokensMinted(assetId, recipient, amount);
    }

    /**
     * @notice Burn treasury tokens (redemption at maturity)
     * @param assetId Asset identifier
     * @param holder Token holder
     * @param amount Amount to burn
     */
    function burnTokens(
        uint256 assetId,
        address holder,
        uint256 amount
    ) external override onlyRole(ASSET_MANAGER_ROLE) nonReentrant {
        require(assets[assetId].assetId != 0, "Asset not found");
        require(holder != address(0), "Invalid holder");
        require(amount > 0, "Invalid amount");

        TreasuryAssetInfo storage asset = assets[assetId];

        // Burn tokens
        TreasuryToken token = TreasuryToken(asset.tokenAddress);
        token.burnFrom(holder, amount);

        // Update counter
        asset.tokensOutstanding -= amount;

        emit TokensBurned(assetId, holder, amount);
    }

    // =============================================================
    //                     ASSET MANAGEMENT
    // =============================================================

    /**
     * @notice Update asset status
     * @param assetId Asset identifier
     * @param newStatus New status
     */
    function updateAssetStatus(
        uint256 assetId,
        AssetStatus newStatus
    ) external override onlyRole(ASSET_MANAGER_ROLE) {
        require(assets[assetId].assetId != 0, "Asset not found");

        TreasuryAssetInfo storage asset = assets[assetId];
        AssetStatus oldStatus = asset.status;

        require(oldStatus != newStatus, "Status unchanged");

        asset.status = newStatus;

        emit AssetStatusUpdated(assetId, oldStatus, newStatus);
    }

    // =============================================================
    //                      VIEW FUNCTIONS
    // =============================================================

    /**
     * @notice Get asset information
     * @param assetId Asset identifier
     * @return Asset info struct
     */
    function getAssetInfo(uint256 assetId)
        external
        view
        override
        returns (TreasuryAssetInfo memory)
    {
        require(assets[assetId].assetId != 0, "Asset not found");
        return assets[assetId];
    }

    /**
     * @notice Get token address for asset
     * @param assetId Asset identifier
     * @return Token contract address
     */
    function getTokenAddress(uint256 assetId) external view override returns (address) {
        return assetTokens[assetId];
    }

    /**
     * @notice Check if asset is active
     * @param assetId Asset identifier
     * @return True if active
     */
    function isAssetActive(uint256 assetId) external view override returns (bool) {
        return assets[assetId].status == AssetStatus.Active;
    }

    /**
     * @notice Get total assets count
     * @return Total number of assets
     */
    function getTotalAssets() external view override returns (uint256) {
        return assetIdCounter;
    }

    /**
     * @notice Get asset ID by CUSIP
     * @param cusip CUSIP identifier
     * @return Asset ID
     */
    function getAssetByCUSIP(string memory cusip) external view returns (uint256) {
        return cusipToAssetId[cusip];
    }

    /**
     * @notice Get all assets of a specific type
     * @param treasuryType Treasury type
     * @return Array of asset IDs
     */
    function getAssetsByType(TreasuryType treasuryType)
        external
        view
        returns (uint256[] memory)
    {
        return assetsByType[treasuryType];
    }

    /**
     * @notice Get all asset IDs
     * @return Array of all asset IDs
     */
    function getAllAssets() external view returns (uint256[] memory) {
        return allAssetIds;
    }

    // =============================================================
    //                     ADMIN FUNCTIONS
    // =============================================================

    /**
     * @notice Pause contract
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause contract
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }
}
