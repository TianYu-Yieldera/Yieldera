// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IRWAAsset
 * @notice Interface for Real World Asset tokens
 * @dev Combines ERC721 (whole asset) and ERC20 (fractionalized) representations
 */
interface IRWAAsset {
    /**
     * @notice Asset types supported by the system
     */
    enum AssetType {
        RealEstate,      // Physical property
        Bonds,           // Debt instruments
        Equity,          // Company shares
        Commodities,     // Gold, oil, etc.
        ArtCollectible,  // Art and collectibles
        Invoice          // Receivables financing
    }

    /**
     * @notice Asset status lifecycle
     */
    enum AssetStatus {
        Pending,         // Awaiting compliance approval
        Active,          // Approved and tradeable
        Suspended,       // Temporarily halted
        Matured,         // Reached maturity (for bonds)
        Defaulted        // Failed to meet obligations
    }

    /**
     * @notice Core asset metadata
     */
    struct AssetMetadata {
        string name;                  // Asset name
        string symbol;                // Token symbol
        AssetType assetType;          // Type of RWA
        uint256 totalValue;           // Total valuation in USD (18 decimals)
        uint256 totalSupply;          // Total fractional tokens
        address issuer;               // Asset issuer/originator
        uint256 issuanceDate;         // Timestamp of issuance
        uint256 maturityDate;         // Maturity date (0 if perpetual)
        string legalDocumentHash;     // IPFS hash of legal documents
        string valuationReportHash;   // IPFS hash of valuation report
    }

    /**
     * @notice Financial terms for yield-bearing assets
     */
    struct YieldTerms {
        uint256 annualYieldRate;      // Annual yield in basis points (e.g., 500 = 5%)
        uint256 yieldPaymentFrequency; // Seconds between payments
        uint256 lastYieldPayment;     // Timestamp of last payment
        uint256 totalYieldPaid;       // Cumulative yield distributed
    }

    // Events
    event AssetCreated(
        uint256 indexed assetId,
        address indexed issuer,
        AssetType assetType,
        uint256 totalValue
    );

    event AssetFractionalized(
        uint256 indexed assetId,
        uint256 totalSupply,
        address fractionalToken
    );

    event AssetStatusChanged(
        uint256 indexed assetId,
        AssetStatus oldStatus,
        AssetStatus newStatus
    );

    event AssetValuationUpdated(
        uint256 indexed assetId,
        uint256 oldValue,
        uint256 newValue,
        uint256 timestamp
    );

    event YieldDistributed(
        uint256 indexed assetId,
        uint256 amount,
        uint256 timestamp
    );

    // View functions
    function getAssetMetadata(uint256 assetId) external view returns (AssetMetadata memory);
    function getAssetStatus(uint256 assetId) external view returns (AssetStatus);
    function getAssetValue(uint256 assetId) external view returns (uint256);
    function getYieldTerms(uint256 assetId) external view returns (YieldTerms memory);
    function getFractionalToken(uint256 assetId) external view returns (address);
    function isAssetActive(uint256 assetId) external view returns (bool);

    // State-changing functions
    function updateAssetStatus(uint256 assetId, AssetStatus newStatus) external;
    function updateAssetValuation(uint256 assetId, uint256 newValue, string calldata valuationReportHash) external;
    function distributeYield(uint256 assetId, uint256 amount) external;
}
