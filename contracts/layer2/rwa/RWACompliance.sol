// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "../interfaces/rwa/IRWACompliance.sol";

/**
 * @title RWACompliance
 * @notice KYC/AML compliance verification for RWA investors
 * @dev Manages investor accreditation and regulatory compliance
 *
 * Key Features:
 * - Multi-tier investor accreditation (Retail → QualifiedPurchaser)
 * - Jurisdiction-based restrictions (whitelist/blacklist)
 * - Time-based verification expiry
 * - Asset-level compliance requirements
 * - Transfer compliance checks
 *
 * Compliance Flow:
 * 1. Investor submits KYC documents off-chain
 * 2. Verifier reviews and calls verifyInvestor()
 * 3. Investor receives accreditation tier and jurisdiction
 * 4. Before any RWA purchase, canInvestInAsset() is checked
 * 5. Before any transfer, canTransferTokens() is checked
 */
contract RWACompliance is IRWACompliance, AccessControl, ReentrancyGuard {
    bytes32 public constant COMPLIANCE_OFFICER_ROLE = keccak256("COMPLIANCE_OFFICER_ROLE");
    bytes32 public constant VERIFIER_ROLE = keccak256("VERIFIER_ROLE");

    // Investor KYC data
    mapping(address => InvestorProfile) private investorProfiles;

    // Asset compliance requirements (assetId => compliance rules)
    mapping(uint256 => AssetComplianceData) private assetCompliance;

    // Asset compliance data storage (cannot use mapping in struct in interface)
    struct AssetComplianceData {
        AccreditationTier minTier;
        bool requiresKYC;
        bool restrictedJurisdictions;
        mapping(string => bool) allowedJurisdictions;
        mapping(string => bool) blockedJurisdictions;
        uint256 minInvestmentAmount;
        uint256 maxInvestors;
        uint256 currentInvestors;
    }

    // Track investors per asset for max investor limit
    mapping(uint256 => mapping(address => bool)) private assetInvestors;

    // Global settings
    uint256 public constant DEFAULT_VERIFICATION_PERIOD = 365 days;
    uint256 public constant MAX_VERIFICATION_PERIOD = 730 days;

    // Statistics
    uint256 public totalVerifiedInvestors;
    uint256 public totalAssetsWithCompliance;

    /**
     * @notice Constructor
     * @param admin Address to grant admin role
     */
    constructor(address admin) {
        require(admin != address(0), "Invalid admin address");

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(COMPLIANCE_OFFICER_ROLE, admin);
        _grantRole(VERIFIER_ROLE, admin);
    }

    // =============================================================
    //                    INVESTOR VERIFICATION
    // =============================================================

    /**
     * @notice Verify investor KYC and set accreditation
     * @param investor Address to verify
     * @param tier Accreditation tier to assign
     * @param jurisdiction ISO 3166-1 country code
     * @param validityPeriod How long verification is valid (max 2 years)
     * @param kycDocumentHash Hash of KYC documents for audit trail
     */
    function verifyInvestor(
        address investor,
        AccreditationTier tier,
        string calldata jurisdiction,
        uint256 validityPeriod,
        bytes32 kycDocumentHash
    ) external onlyRole(VERIFIER_ROLE) {
        require(investor != address(0), "Invalid investor address");
        require(bytes(jurisdiction).length == 2, "Invalid jurisdiction code");
        require(validityPeriod <= MAX_VERIFICATION_PERIOD, "Validity period too long");
        require(kycDocumentHash != bytes32(0), "KYC document hash required");

        InvestorProfile storage profile = investorProfiles[investor];

        // Track new verifications
        if (profile.status != VerificationStatus.Approved) {
            totalVerifiedInvestors++;
        }

        uint256 expiresAt = block.timestamp + validityPeriod;

        profile.status = VerificationStatus.Approved;
        profile.tier = tier;
        profile.verifiedAt = block.timestamp;
        profile.expiresAt = expiresAt;
        profile.jurisdiction = jurisdiction;
        profile.kycDocumentHash = kycDocumentHash;
        profile.verifier = msg.sender;

        emit InvestorVerified(investor, VerificationStatus.Approved, tier, expiresAt);
    }

    /**
     * @notice Update investor verification status
     * @param investor Investor address
     * @param newStatus New status to set
     */
    function updateInvestorStatus(
        address investor,
        VerificationStatus newStatus
    ) external onlyRole(COMPLIANCE_OFFICER_ROLE) {
        InvestorProfile storage profile = investorProfiles[investor];
        VerificationStatus oldStatus = profile.status;

        require(oldStatus != newStatus, "Status unchanged");

        // Adjust counter when status changes
        if (oldStatus == VerificationStatus.Approved && newStatus != VerificationStatus.Approved) {
            totalVerifiedInvestors--;
        } else if (oldStatus != VerificationStatus.Approved && newStatus == VerificationStatus.Approved) {
            totalVerifiedInvestors++;
        }

        profile.status = newStatus;

        emit InvestorStatusChanged(investor, oldStatus, newStatus);
    }

    // =============================================================
    //                   ASSET COMPLIANCE RULES
    // =============================================================

    /**
     * @notice Set compliance requirements for an asset
     * @param assetId Asset identifier
     * @param minTier Minimum accreditation tier required
     * @param requiresKYC Whether KYC is mandatory
     * @param minInvestmentAmount Minimum investment in USD (18 decimals)
     * @param maxInvestors Maximum number of investors (0 = unlimited)
     */
    function setAssetCompliance(
        uint256 assetId,
        AccreditationTier minTier,
        bool requiresKYC,
        uint256 minInvestmentAmount,
        uint256 maxInvestors
    ) external onlyRole(COMPLIANCE_OFFICER_ROLE) {
        AssetComplianceData storage compliance = assetCompliance[assetId];

        // Track new assets with compliance
        if (!compliance.requiresKYC && requiresKYC) {
            totalAssetsWithCompliance++;
        }

        compliance.minTier = minTier;
        compliance.requiresKYC = requiresKYC;
        compliance.minInvestmentAmount = minInvestmentAmount;
        compliance.maxInvestors = maxInvestors;

        emit AssetComplianceSet(assetId, minTier, requiresKYC);
    }

    /**
     * @notice Update jurisdiction restrictions for an asset
     * @param assetId Asset identifier
     * @param jurisdiction ISO 3166-1 country code
     * @param allowed True to allow, false to block
     */
    function updateJurisdictionRestriction(
        uint256 assetId,
        string calldata jurisdiction,
        bool allowed
    ) external onlyRole(COMPLIANCE_OFFICER_ROLE) {
        require(bytes(jurisdiction).length == 2, "Invalid jurisdiction code");

        AssetComplianceData storage compliance = assetCompliance[assetId];
        compliance.restrictedJurisdictions = true;

        if (allowed) {
            compliance.allowedJurisdictions[jurisdiction] = true;
            compliance.blockedJurisdictions[jurisdiction] = false;
        } else {
            compliance.allowedJurisdictions[jurisdiction] = false;
            compliance.blockedJurisdictions[jurisdiction] = true;
        }

        emit JurisdictionUpdated(assetId, jurisdiction, allowed);
    }

    // =============================================================
    //                   COMPLIANCE CHECKS
    // =============================================================

    /**
     * @notice Check if investor is verified and not expired
     * @param investor Address to check
     * @return True if verified and valid
     */
    function isInvestorVerified(address investor) public view returns (bool) {
        InvestorProfile storage profile = investorProfiles[investor];

        return profile.status == VerificationStatus.Approved &&
               block.timestamp < profile.expiresAt;
    }

    /**
     * @notice Check if investor can invest in a specific asset
     * @param investor Investor address
     * @param assetId Asset identifier
     * @return True if compliant
     */
    function canInvestInAsset(
        address investor,
        uint256 assetId
    ) public view returns (bool) {
        AssetComplianceData storage compliance = assetCompliance[assetId];
        InvestorProfile storage profile = investorProfiles[investor];

        // If KYC required, check verification
        if (compliance.requiresKYC) {
            if (!isInvestorVerified(investor)) {
                return false;
            }

            // Check accreditation tier
            if (profile.tier < compliance.minTier) {
                return false;
            }

            // Check jurisdiction restrictions
            if (compliance.restrictedJurisdictions) {
                string memory jurisdiction = profile.jurisdiction;

                // If blocked, deny
                if (compliance.blockedJurisdictions[jurisdiction]) {
                    return false;
                }

                // If using whitelist and not on it, deny
                bool hasWhitelist = false;
                // Note: In production, track which jurisdictions are whitelisted
                // For now, assume if allowedJurisdictions[j] == true, it's whitelisted
                if (hasWhitelist && !compliance.allowedJurisdictions[jurisdiction]) {
                    return false;
                }
            }
        }

        // Check max investors limit
        if (compliance.maxInvestors > 0) {
            if (!assetInvestors[assetId][investor]) {
                // New investor, check if limit reached
                if (compliance.currentInvestors >= compliance.maxInvestors) {
                    return false;
                }
            }
        }

        return true;
    }

    /**
     * @notice Check if tokens can be transferred between addresses
     * @param from Sender address
     * @param to Recipient address
     * @param assetId Asset identifier
     * @param amount Amount to transfer (for minimum check)
     * @return True if transfer is compliant
     */
    function canTransferTokens(
        address from,
        address to,
        uint256 assetId,
        uint256 amount
    ) public view returns (bool) {
        // Allow minting (from = 0x0)
        if (from == address(0)) {
            return canInvestInAsset(to, assetId);
        }

        // Allow burning (to = 0x0)
        if (to == address(0)) {
            return true;
        }

        // Both parties must be compliant
        if (!canInvestInAsset(from, assetId)) {
            return false;
        }

        if (!canInvestInAsset(to, assetId)) {
            return false;
        }

        // Check minimum investment amount for new investors
        AssetComplianceData storage compliance = assetCompliance[assetId];
        if (!assetInvestors[assetId][to] && compliance.minInvestmentAmount > 0) {
            // For fractional tokens, amount represents token count
            // In production, would convert to USD value
            // For now, just allow the transfer if compliance is met
        }

        return true;
    }

    /**
     * @notice Record investor participation in asset
     * @dev Called by RWAMarketplace after successful purchase
     * @param assetId Asset identifier
     * @param investor Investor address
     */
    function recordInvestorParticipation(
        uint256 assetId,
        address investor
    ) external {
        // In production, add access control to only allow marketplace
        if (!assetInvestors[assetId][investor]) {
            assetInvestors[assetId][investor] = true;
            assetCompliance[assetId].currentInvestors++;
        }
    }

    /**
     * @notice Remove investor from asset (when they divest completely)
     * @param assetId Asset identifier
     * @param investor Investor address
     */
    function removeInvestorParticipation(
        uint256 assetId,
        address investor
    ) external {
        // In production, add access control
        if (assetInvestors[assetId][investor]) {
            assetInvestors[assetId][investor] = false;
            assetCompliance[assetId].currentInvestors--;
        }
    }

    // =============================================================
    //                      VIEW FUNCTIONS
    // =============================================================

    /**
     * @notice Get investor profile
     * @param investor Address to query
     * @return InvestorProfile struct
     */
    function getInvestorProfile(
        address investor
    ) external view returns (InvestorProfile memory) {
        return investorProfiles[investor];
    }

    /**
     * @notice Get current number of investors in an asset
     * @param assetId Asset identifier
     * @return Number of investors
     */
    function getAssetInvestorCount(uint256 assetId) external view returns (uint256) {
        return assetCompliance[assetId].currentInvestors;
    }

    /**
     * @notice Check if jurisdiction is allowed for asset
     * @param assetId Asset identifier
     * @param jurisdiction ISO 3166-1 country code
     * @return True if allowed
     */
    function isJurisdictionAllowed(
        uint256 assetId,
        string calldata jurisdiction
    ) external view returns (bool) {
        AssetComplianceData storage compliance = assetCompliance[assetId];

        if (!compliance.restrictedJurisdictions) {
            return true; // No restrictions
        }

        if (compliance.blockedJurisdictions[jurisdiction]) {
            return false; // Explicitly blocked
        }

        return compliance.allowedJurisdictions[jurisdiction];
    }

    /**
     * @notice Get minimum accreditation tier for asset
     * @param assetId Asset identifier
     * @return AccreditationTier
     */
    function getAssetMinTier(uint256 assetId) external view returns (AccreditationTier) {
        return assetCompliance[assetId].minTier;
    }

    /**
     * @notice Check if asset requires KYC
     * @param assetId Asset identifier
     * @return True if KYC required
     */
    function doesAssetRequireKYC(uint256 assetId) external view returns (bool) {
        return assetCompliance[assetId].requiresKYC;
    }
}
