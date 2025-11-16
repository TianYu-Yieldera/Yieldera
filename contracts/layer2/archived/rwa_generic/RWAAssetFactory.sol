// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "../interfaces/rwa/IRWAAsset.sol";
import "../interfaces/rwa/IRWACompliance.sol";
import "../interfaces/rwa/IRWAValuation.sol";
import "./FractionalRWAToken.sol";

/**
 * @title RWAAssetFactory
 * @notice Factory for creating and managing tokenized real-world assets
 * @dev Integrates compliance, valuation, and fractionalization
 *
 * Key Features:
 * - Create RWA assets with comprehensive metadata
 * - Fractionalize assets into ERC20 tokens
 * - Integrate with compliance for investor verification
 * - Integrate with valuation for pricing
 * - Manage asset lifecycle (Pending → Active → Matured/Defaulted)
 * - Support yield-bearing assets (bonds, real estate)
 *
 * Asset Creation Flow:
 * 1. Issuer calls createAsset() with metadata and legal docs
 * 2. Asset created in Pending status
 * 3. Compliance officer reviews and approves
 * 4. Issuer fractionalizes asset into ERC20 tokens
 * 5. Asset becomes Active and tradeable
 * 6. Valuation and yield distribution can proceed
 */
contract RWAAssetFactory is IRWAAsset, AccessControl, ReentrancyGuard, Pausable {
    bytes32 public constant ASSET_MANAGER_ROLE = keccak256("ASSET_MANAGER_ROLE");
    bytes32 public constant COMPLIANCE_ROLE = keccak256("COMPLIANCE_ROLE");

    // Asset counter
    uint256 private assetIdCounter;

    // Asset metadata storage
    mapping(uint256 => AssetMetadata) private assetMetadata;
    mapping(uint256 => YieldTerms) private assetYieldTerms;
    mapping(uint256 => AssetStatus) private assetStatus;
    mapping(uint256 => address) private fractionalTokens;

    // Integration contracts
    IRWACompliance public immutable complianceContract;
    IRWAValuation public immutable valuationContract;

    // Asset registry
    mapping(address => uint256[]) private issuerAssets;
    uint256[] private allAssetIds;

    // Statistics
    uint256 public totalAssets;
    uint256 public activeAssets;
    uint256 public totalValueLocked; // Sum of all active asset values

    /**
     * @notice Constructor
     * @param admin Admin address
     * @param compliance Compliance contract
     * @param valuation Valuation contract
     */
    constructor(
        address admin,
        address compliance,
        address valuation
    ) {
        require(admin != address(0), "Invalid admin");
        require(compliance != address(0), "Invalid compliance");
        require(valuation != address(0), "Invalid valuation");

        complianceContract = IRWACompliance(compliance);
        valuationContract = IRWAValuation(valuation);

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(ASSET_MANAGER_ROLE, admin);
        _grantRole(COMPLIANCE_ROLE, admin);
    }

    // =============================================================
    //                     ASSET CREATION
    // =============================================================

    /**
     * @notice Create new RWA asset
     * @param name Asset name
     * @param symbol Token symbol for fractionalized tokens
     * @param assetType Type of RWA
     * @param totalValue Initial valuation in USD (18 decimals)
     * @param maturityDate Maturity timestamp (0 if perpetual)
     * @param legalDocumentHash IPFS hash of legal documents
     * @param valuationReportHash IPFS hash of initial valuation report
     * @return assetId Newly created asset ID
     */
    function createAsset(
        string calldata name,
        string calldata symbol,
        AssetType assetType,
        uint256 totalValue,
        uint256 maturityDate,
        string calldata legalDocumentHash,
        string calldata valuationReportHash
    ) external whenNotPaused nonReentrant returns (uint256 assetId) {
        require(bytes(name).length > 0, "Invalid name");
        require(bytes(symbol).length > 0, "Invalid symbol");
        require(totalValue > 0, "Invalid valuation");
        require(bytes(legalDocumentHash).length > 0, "Legal docs required");
        require(bytes(valuationReportHash).length > 0, "Valuation report required");

        // Check issuer is KYC verified
        require(
            complianceContract.isInvestorVerified(msg.sender),
            "Issuer not verified"
        );

        assetIdCounter++;
        assetId = assetIdCounter;

        AssetMetadata storage metadata = assetMetadata[assetId];
        metadata.name = name;
        metadata.symbol = symbol;
        metadata.assetType = assetType;
        metadata.totalValue = totalValue;
        metadata.totalSupply = 0; // Set during fractionalization
        metadata.issuer = msg.sender;
        metadata.issuanceDate = block.timestamp;
        metadata.maturityDate = maturityDate;
        metadata.legalDocumentHash = legalDocumentHash;
        metadata.valuationReportHash = valuationReportHash;

        assetStatus[assetId] = AssetStatus.Pending;

        // Track issuer's assets
        issuerAssets[msg.sender].push(assetId);
        allAssetIds.push(assetId);

        totalAssets++;

        emit AssetCreated(assetId, msg.sender, assetType, totalValue);
    }

    /**
     * @notice Fractionalize asset into ERC20 tokens
     * @param assetId Asset to fractionalize
     * @param totalSupply Total supply of fractional tokens
     * @return tokenAddress Address of deployed fractional token
     */
    function fractionalizeAsset(
        uint256 assetId,
        uint256 totalSupply
    ) external nonReentrant returns (address tokenAddress) {
        AssetMetadata storage metadata = assetMetadata[assetId];

        require(metadata.issuer == msg.sender, "Not asset issuer");
        require(assetStatus[assetId] == AssetStatus.Pending, "Asset not pending");
        require(fractionalTokens[assetId] == address(0), "Already fractionalized");
        require(totalSupply > 0, "Invalid supply");

        // Deploy fractional token
        FractionalRWAToken token = new FractionalRWAToken(
            metadata.name,
            metadata.symbol,
            assetId,
            totalSupply,
            address(complianceContract),
            msg.sender // Issuer is admin
        );

        tokenAddress = address(token);
        fractionalTokens[assetId] = tokenAddress;
        metadata.totalSupply = totalSupply;

        emit AssetFractionalized(assetId, totalSupply, tokenAddress);
    }

    /**
     * @notice Activate asset after compliance approval
     * @param assetId Asset to activate
     */
    function activateAsset(uint256 assetId) external onlyRole(COMPLIANCE_ROLE) {
        require(assetStatus[assetId] == AssetStatus.Pending, "Not pending");
        require(fractionalTokens[assetId] != address(0), "Not fractionalized");

        _updateAssetStatus(assetId, AssetStatus.Active);

        activeAssets++;
        totalValueLocked += assetMetadata[assetId].totalValue;
    }

    // =============================================================
    //                   ASSET MANAGEMENT
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
        _updateAssetStatus(assetId, newStatus);
    }

    /**
     * @notice Internal status update with event
     */
    function _updateAssetStatus(uint256 assetId, AssetStatus newStatus) internal {
        AssetStatus oldStatus = assetStatus[assetId];
        require(oldStatus != newStatus, "Status unchanged");

        assetStatus[assetId] = newStatus;

        // Update active asset counter
        if (oldStatus == AssetStatus.Active && newStatus != AssetStatus.Active) {
            activeAssets--;
            totalValueLocked -= assetMetadata[assetId].totalValue;
        } else if (oldStatus != AssetStatus.Active && newStatus == AssetStatus.Active) {
            activeAssets++;
            totalValueLocked += assetMetadata[assetId].totalValue;
        }

        emit AssetStatusChanged(assetId, oldStatus, newStatus);
    }

    /**
     * @notice Update asset valuation
     * @param assetId Asset identifier
     * @param newValue New valuation
     * @param valuationReportHash IPFS hash of report
     */
    function updateAssetValuation(
        uint256 assetId,
        uint256 newValue,
        string calldata valuationReportHash
    ) external override onlyRole(ASSET_MANAGER_ROLE) {
        require(newValue > 0, "Invalid valuation");

        AssetMetadata storage metadata = assetMetadata[assetId];
        uint256 oldValue = metadata.totalValue;

        metadata.totalValue = newValue;
        metadata.valuationReportHash = valuationReportHash;

        // Update TVL if active
        if (assetStatus[assetId] == AssetStatus.Active) {
            totalValueLocked = totalValueLocked - oldValue + newValue;
        }

        emit AssetValuationUpdated(assetId, oldValue, newValue, block.timestamp);
    }

    /**
     * @notice Set yield terms for asset
     * @param assetId Asset identifier
     * @param annualYieldRate Annual yield in basis points (e.g., 500 = 5%)
     * @param yieldPaymentFrequency Seconds between payments
     */
    function setYieldTerms(
        uint256 assetId,
        uint256 annualYieldRate,
        uint256 yieldPaymentFrequency
    ) external {
        AssetMetadata storage metadata = assetMetadata[assetId];
        require(metadata.issuer == msg.sender, "Not asset issuer");
        require(annualYieldRate <= 10000, "Invalid yield rate"); // Max 100%
        require(yieldPaymentFrequency > 0, "Invalid frequency");

        YieldTerms storage terms = assetYieldTerms[assetId];
        terms.annualYieldRate = annualYieldRate;
        terms.yieldPaymentFrequency = yieldPaymentFrequency;
        terms.lastYieldPayment = block.timestamp;
        terms.totalYieldPaid = 0;
    }

    /**
     * @notice Distribute yield to token holders
     * @param assetId Asset identifier
     * @param amount Yield amount to distribute
     */
    function distributeYield(
        uint256 assetId,
        uint256 amount
    ) external override {
        AssetMetadata storage metadata = assetMetadata[assetId];
        require(metadata.issuer == msg.sender, "Not asset issuer");
        require(isAssetActive(assetId), "Asset not active");

        YieldTerms storage terms = assetYieldTerms[assetId];
        terms.lastYieldPayment = block.timestamp;
        terms.totalYieldPaid += amount;

        emit YieldDistributed(assetId, amount, block.timestamp);
    }

    // =============================================================
    //                      VIEW FUNCTIONS
    // =============================================================

    /**
     * @notice Get asset metadata
     * @param assetId Asset identifier
     * @return AssetMetadata struct
     */
    function getAssetMetadata(
        uint256 assetId
    ) external view override returns (AssetMetadata memory) {
        return assetMetadata[assetId];
    }

    /**
     * @notice Get asset status
     * @param assetId Asset identifier
     * @return AssetStatus enum
     */
    function getAssetStatus(
        uint256 assetId
    ) external view override returns (AssetStatus) {
        return assetStatus[assetId];
    }

    /**
     * @notice Get asset value
     * @param assetId Asset identifier
     * @return Value in USD (18 decimals)
     */
    function getAssetValue(
        uint256 assetId
    ) external view override returns (uint256) {
        return assetMetadata[assetId].totalValue;
    }

    /**
     * @notice Get yield terms
     * @param assetId Asset identifier
     * @return YieldTerms struct
     */
    function getYieldTerms(
        uint256 assetId
    ) external view override returns (YieldTerms memory) {
        return assetYieldTerms[assetId];
    }

    /**
     * @notice Get fractional token address
     * @param assetId Asset identifier
     * @return Token contract address
     */
    function getFractionalToken(
        uint256 assetId
    ) external view override returns (address) {
        return fractionalTokens[assetId];
    }

    /**
     * @notice Check if asset is active
     * @param assetId Asset identifier
     * @return True if active
     */
    function isAssetActive(
        uint256 assetId
    ) public view override returns (bool) {
        return assetStatus[assetId] == AssetStatus.Active;
    }

    /**
     * @notice Get all assets by issuer
     * @param issuer Issuer address
     * @return Array of asset IDs
     */
    function getIssuerAssets(address issuer) external view returns (uint256[] memory) {
        return issuerAssets[issuer];
    }

    /**
     * @notice Get all asset IDs
     * @return Array of all asset IDs
     */
    function getAllAssets() external view returns (uint256[] memory) {
        return allAssetIds;
    }

    /**
     * @notice Get assets by type
     * @param assetType Asset type to filter
     * @return Array of matching asset IDs
     */
    function getAssetsByType(AssetType assetType) external view returns (uint256[] memory) {
        uint256 count = 0;

        // Count matching assets
        for (uint256 i = 0; i < allAssetIds.length; i++) {
            if (assetMetadata[allAssetIds[i]].assetType == assetType) {
                count++;
            }
        }

        // Build result array
        uint256[] memory result = new uint256[](count);
        uint256 index = 0;

        for (uint256 i = 0; i < allAssetIds.length; i++) {
            if (assetMetadata[allAssetIds[i]].assetType == assetType) {
                result[index] = allAssetIds[i];
                index++;
            }
        }

        return result;
    }

    /**
     * @notice Get active assets
     * @return Array of active asset IDs
     */
    function getActiveAssets() external view returns (uint256[] memory) {
        uint256 count = 0;

        for (uint256 i = 0; i < allAssetIds.length; i++) {
            if (assetStatus[allAssetIds[i]] == AssetStatus.Active) {
                count++;
            }
        }

        uint256[] memory result = new uint256[](count);
        uint256 index = 0;

        for (uint256 i = 0; i < allAssetIds.length; i++) {
            if (assetStatus[allAssetIds[i]] == AssetStatus.Active) {
                result[index] = allAssetIds[i];
                index++;
            }
        }

        return result;
    }

    // =============================================================
    //                    ADMIN FUNCTIONS
    // =============================================================

    /**
     * @notice Pause asset creation
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause asset creation
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }
}
