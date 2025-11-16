// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

/**
 * @title KYCRegistry
 * @notice On-chain KYC/AML compliance registry for RWA and Treasury operations
 * @dev Stores KYC verification status without exposing personal data
 */
contract KYCRegistry is AccessControl, Pausable {
    bytes32 public constant KYC_PROVIDER_ROLE = keccak256("KYC_PROVIDER_ROLE");
    bytes32 public constant COMPLIANCE_OFFICER_ROLE = keccak256("COMPLIANCE_OFFICER_ROLE");

    enum KYCLevel {
        None,
        Basic,
        Enhanced,
        Institutional
    }

    struct KYCData {
        KYCLevel level;
        uint256 verifiedAt;
        uint256 expiresAt;
        uint256 riskScore; // 0-100, lower is better
        bytes32 dataHash; // Hash of off-chain data for verification
        address provider;
        bool isSanctioned;
    }

    // User address => KYC data
    mapping(address => KYCData) public kycRecords;

    // Sanctioned addresses
    mapping(address => bool) public sanctionsList;

    // Country restrictions for certain operations
    mapping(string => bool) public restrictedCountries;

    // Compliance requirements for different operations
    mapping(string => KYCLevel) public operationRequirements;

    // Events
    event KYCVerified(
        address indexed user,
        KYCLevel level,
        uint256 riskScore,
        address indexed provider
    );

    event KYCUpdated(
        address indexed user,
        KYCLevel oldLevel,
        KYCLevel newLevel,
        address indexed provider
    );

    event KYCExpired(address indexed user);

    event SanctionAdded(address indexed user, address indexed addedBy);

    event SanctionRemoved(address indexed user, address indexed removedBy);

    event OperationRequirementSet(string operation, KYCLevel requiredLevel);

    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(COMPLIANCE_OFFICER_ROLE, msg.sender);

        // Set default operation requirements
        operationRequirements["treasury"] = KYCLevel.Enhanced;
        operationRequirements["rwa"] = KYCLevel.Basic;
        operationRequirements["defi"] = KYCLevel.None;
        operationRequirements["institutional"] = KYCLevel.Institutional;
    }

    /**
     * @notice Submit KYC verification for a user
     * @param user Address of the user
     * @param level KYC verification level
     * @param riskScore Risk assessment score (0-100)
     * @param validityPeriod How long the KYC is valid (in seconds)
     * @param dataHash Hash of off-chain KYC data
     */
    function submitKYC(
        address user,
        KYCLevel level,
        uint256 riskScore,
        uint256 validityPeriod,
        bytes32 dataHash
    ) external onlyRole(KYC_PROVIDER_ROLE) whenNotPaused {
        require(user != address(0), "Invalid user address");
        require(riskScore <= 100, "Invalid risk score");
        require(validityPeriod > 0, "Invalid validity period");

        KYCData storage kyc = kycRecords[user];
        KYCLevel oldLevel = kyc.level;

        kyc.level = level;
        kyc.verifiedAt = block.timestamp;
        kyc.expiresAt = block.timestamp + validityPeriod;
        kyc.riskScore = riskScore;
        kyc.dataHash = dataHash;
        kyc.provider = msg.sender;

        if (oldLevel == KYCLevel.None) {
            emit KYCVerified(user, level, riskScore, msg.sender);
        } else {
            emit KYCUpdated(user, oldLevel, level, msg.sender);
        }
    }

    /**
     * @notice Add address to sanctions list
     * @param user Address to sanction
     */
    function addSanction(address user) external onlyRole(COMPLIANCE_OFFICER_ROLE) {
        require(user != address(0), "Invalid address");
        require(!sanctionsList[user], "Already sanctioned");

        sanctionsList[user] = true;
        kycRecords[user].isSanctioned = true;

        emit SanctionAdded(user, msg.sender);
    }

    /**
     * @notice Remove address from sanctions list
     * @param user Address to remove from sanctions
     */
    function removeSanction(address user) external onlyRole(COMPLIANCE_OFFICER_ROLE) {
        require(sanctionsList[user], "Not sanctioned");

        sanctionsList[user] = false;
        kycRecords[user].isSanctioned = false;

        emit SanctionRemoved(user, msg.sender);
    }

    /**
     * @notice Check if user is compliant for specific operation
     * @param user User address to check
     * @param operation Operation type
     * @return compliant Whether user meets compliance requirements
     * @return reason Reason if not compliant
     */
    function checkCompliance(
        address user,
        string memory operation
    ) external view returns (bool compliant, string memory reason) {
        // Check sanctions first
        if (sanctionsList[user]) {
            return (false, "User is sanctioned");
        }

        KYCData memory kyc = kycRecords[user];

        // Check if KYC exists
        if (kyc.verifiedAt == 0) {
            return (false, "No KYC record found");
        }

        // Check if KYC is expired
        if (block.timestamp > kyc.expiresAt) {
            return (false, "KYC expired");
        }

        // Check if KYC level meets operation requirements
        KYCLevel requiredLevel = operationRequirements[operation];
        if (kyc.level < requiredLevel) {
            return (false, "Insufficient KYC level");
        }

        // Check risk score for high-risk operations
        if (requiredLevel >= KYCLevel.Enhanced && kyc.riskScore > 70) {
            return (false, "Risk score too high");
        }

        return (true, "Compliant");
    }

    /**
     * @notice Get user's KYC status
     * @param user User address
     * @return level KYC level
     * @return isValid Whether KYC is currently valid
     * @return riskScore Risk assessment score
     */
    function getUserKYC(address user) external view returns (
        KYCLevel level,
        bool isValid,
        uint256 riskScore
    ) {
        KYCData memory kyc = kycRecords[user];

        level = kyc.level;
        isValid = kyc.verifiedAt > 0 &&
                 block.timestamp <= kyc.expiresAt &&
                 !kyc.isSanctioned;
        riskScore = kyc.riskScore;
    }

    /**
     * @notice Set operation compliance requirements
     * @param operation Operation identifier
     * @param requiredLevel Required KYC level
     */
    function setOperationRequirement(
        string memory operation,
        KYCLevel requiredLevel
    ) external onlyRole(COMPLIANCE_OFFICER_ROLE) {
        operationRequirements[operation] = requiredLevel;
        emit OperationRequirementSet(operation, requiredLevel);
    }

    /**
     * @notice Batch check multiple addresses for sanctions
     * @param addresses Array of addresses to check
     * @return sanctioned Array of boolean values indicating sanction status
     */
    function batchCheckSanctions(
        address[] memory addresses
    ) external view returns (bool[] memory sanctioned) {
        sanctioned = new bool[](addresses.length);
        for (uint256 i = 0; i < addresses.length; i++) {
            sanctioned[i] = sanctionsList[addresses[i]];
        }
    }

    /**
     * @notice Emergency pause
     */
    function pause() external onlyRole(COMPLIANCE_OFFICER_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause
     */
    function unpause() external onlyRole(COMPLIANCE_OFFICER_ROLE) {
        _unpause();
    }

    /**
     * @notice Check if address has specific role
     */
    function hasRole(bytes32 role, address account) public view override returns (bool) {
        return super.hasRole(role, account);
    }
}