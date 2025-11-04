// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IRWACompliance
 * @notice Interface for KYC/AML compliance verification
 * @dev Manages investor accreditation and regulatory compliance
 */
interface IRWACompliance {
    /**
     * @notice Investor verification status
     */
    enum VerificationStatus {
        NotVerified,     // No verification submitted
        Pending,         // Under review
        Approved,        // KYC approved
        Rejected,        // KYC rejected
        Expired,         // Verification expired
        Suspended        // Temporarily suspended
    }

    /**
     * @notice Investor accreditation tiers
     */
    enum AccreditationTier {
        Retail,          // Standard retail investor
        Accredited,      // Accredited investor (US SEC definition)
        Institutional,   // Institutional investor
        QualifiedPurchaser // Qualified purchaser (higher threshold)
    }

    /**
     * @notice Investor compliance data
     */
    struct InvestorProfile {
        VerificationStatus status;
        AccreditationTier tier;
        uint256 verifiedAt;
        uint256 expiresAt;
        string jurisdiction;        // ISO 3166-1 country code
        bytes32 kycDocumentHash;    // Hash of KYC documents
        address verifier;           // Address that verified KYC
    }

    /**
     * @notice Asset compliance requirements
     */
    struct AssetCompliance {
        AccreditationTier minTier;  // Minimum tier required
        bool requiresKYC;           // KYC mandatory
        bool restrictedJurisdictions; // Has jurisdiction restrictions
        mapping(string => bool) allowedJurisdictions; // Whitelist
        mapping(string => bool) blockedJurisdictions; // Blacklist
        uint256 minInvestmentAmount; // Minimum investment in USD
        uint256 maxInvestors;       // Maximum number of investors
        uint256 currentInvestors;   // Current investor count
    }

    // Events
    event InvestorVerified(
        address indexed investor,
        VerificationStatus status,
        AccreditationTier tier,
        uint256 expiresAt
    );

    event InvestorStatusChanged(
        address indexed investor,
        VerificationStatus oldStatus,
        VerificationStatus newStatus
    );

    event AssetComplianceSet(
        uint256 indexed assetId,
        AccreditationTier minTier,
        bool requiresKYC
    );

    event JurisdictionUpdated(
        uint256 indexed assetId,
        string jurisdiction,
        bool allowed
    );

    event ComplianceViolation(
        address indexed investor,
        uint256 indexed assetId,
        string reason
    );

    // View functions
    function isInvestorVerified(address investor) external view returns (bool);
    function getInvestorProfile(address investor) external view returns (InvestorProfile memory);
    function canInvestInAsset(address investor, uint256 assetId) external view returns (bool);
    function canTransferTokens(address from, address to, uint256 assetId, uint256 amount) external view returns (bool);
    function getAssetInvestorCount(uint256 assetId) external view returns (uint256);

    // State-changing functions
    function verifyInvestor(
        address investor,
        AccreditationTier tier,
        string calldata jurisdiction,
        uint256 validityPeriod,
        bytes32 kycDocumentHash
    ) external;

    function updateInvestorStatus(address investor, VerificationStatus newStatus) external;

    function setAssetCompliance(
        uint256 assetId,
        AccreditationTier minTier,
        bool requiresKYC,
        uint256 minInvestmentAmount,
        uint256 maxInvestors
    ) external;

    function updateJurisdictionRestriction(
        uint256 assetId,
        string calldata jurisdiction,
        bool allowed
    ) external;
}
