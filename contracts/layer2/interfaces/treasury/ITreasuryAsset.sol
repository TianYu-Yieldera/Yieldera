// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title ITreasuryAsset
 * @notice Interface for US Treasury tokenized assets
 */
interface ITreasuryAsset {
    /// @notice Treasury asset types
    enum TreasuryType {
        T_BILL,    // Treasury Bills: < 1 year
        T_NOTE,    // Treasury Notes: 2-10 years
        T_BOND     // Treasury Bonds: 20-30 years
    }

    /// @notice Asset status
    enum AssetStatus {
        Active,
        Matured,
        Suspended
    }

    /// @notice Treasury asset metadata
    struct TreasuryAssetInfo {
        uint256 assetId;
        TreasuryType treasuryType;
        string maturityTerm;        // e.g., "13W", "2Y", "10Y", "30Y"
        string cusip;               // CUSIP identifier
        uint256 issueDate;          // Timestamp
        uint256 maturityDate;       // Timestamp
        uint256 faceValue;          // Face value in USD (18 decimals)
        uint256 couponRate;         // Annual coupon rate in basis points (425 = 4.25%)
        uint256 tokensIssued;       // Total tokens issued
        uint256 tokensOutstanding;  // Tokens in circulation
        address tokenAddress;       // ERC20 token contract
        AssetStatus status;
        uint256 createdAt;
    }

    /// @notice Events
    event TreasuryAssetCreated(
        uint256 indexed assetId,
        TreasuryType treasuryType,
        string cusip,
        address indexed tokenAddress,
        uint256 faceValue
    );

    event AssetStatusUpdated(
        uint256 indexed assetId,
        AssetStatus oldStatus,
        AssetStatus newStatus
    );

    event TokensMinted(
        uint256 indexed assetId,
        address indexed recipient,
        uint256 amount
    );

    event TokensBurned(
        uint256 indexed assetId,
        address indexed holder,
        uint256 amount
    );

    /// @notice Create new treasury asset
    function createTreasuryAsset(
        TreasuryType treasuryType,
        string memory maturityTerm,
        string memory cusip,
        uint256 issueDate,
        uint256 maturityDate,
        uint256 faceValue,
        uint256 couponRate
    ) external returns (uint256 assetId, address tokenAddress);

    /// @notice Mint treasury tokens
    function mintTokens(
        uint256 assetId,
        address recipient,
        uint256 amount
    ) external;

    /// @notice Burn treasury tokens (at maturity or redemption)
    function burnTokens(
        uint256 assetId,
        address holder,
        uint256 amount
    ) external;

    /// @notice Update asset status
    function updateAssetStatus(
        uint256 assetId,
        AssetStatus newStatus
    ) external;

    /// @notice Get asset information
    function getAssetInfo(uint256 assetId) external view returns (TreasuryAssetInfo memory);

    /// @notice Get token address for asset
    function getTokenAddress(uint256 assetId) external view returns (address);

    /// @notice Check if asset is active
    function isAssetActive(uint256 assetId) external view returns (bool);

    /// @notice Get total assets count
    function getTotalAssets() external view returns (uint256);
}
