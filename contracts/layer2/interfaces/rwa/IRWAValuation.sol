// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IRWAValuation
 * @notice Interface for RWA asset valuation and pricing
 * @dev Integrates with oracles and external valuators
 */
interface IRWAValuation {
    /**
     * @notice Valuation method types
     */
    enum ValuationMethod {
        Manual,          // Manual appraisal by certified valuator
        Oracle,          // Chainlink or other oracle feed
        Formula,         // Algorithmic calculation (e.g., bond pricing)
        Hybrid           // Combination of methods
    }

    /**
     * @notice Valuation data point
     */
    struct Valuation {
        uint256 value;              // Asset value in USD (18 decimals)
        uint256 timestamp;          // When valuation was performed
        ValuationMethod method;     // How value was determined
        address valuator;           // Who performed valuation
        string reportHash;          // IPFS hash of valuation report
        uint256 confidence;         // Confidence score (0-10000 basis points)
    }

    /**
     * @notice Oracle configuration for automated pricing
     */
    struct OracleConfig {
        address oracleAddress;      // Chainlink oracle address
        uint256 heartbeat;          // Maximum age of oracle data (seconds)
        uint256 deviationThreshold; // Max % deviation before revaluation
        bool isActive;              // Oracle currently in use
    }

    // Events
    event ValuationUpdated(
        uint256 indexed assetId,
        uint256 oldValue,
        uint256 newValue,
        ValuationMethod method,
        address indexed valuator
    );

    event ValuatorAuthorized(
        address indexed valuator,
        bool authorized
    );

    event OracleConfigured(
        uint256 indexed assetId,
        address indexed oracle,
        uint256 heartbeat
    );

    event ValuationDisputed(
        uint256 indexed assetId,
        uint256 disputedValue,
        address indexed disputer,
        string reason
    );

    // View functions
    function getCurrentValue(uint256 assetId) external view returns (uint256);
    function getLastValuation(uint256 assetId) external view returns (Valuation memory);
    function getValuationHistory(uint256 assetId, uint256 count) external view returns (Valuation[] memory);
    function getPricePerToken(uint256 assetId) external view returns (uint256);
    function isValuationStale(uint256 assetId) external view returns (bool);
    function isValuatorAuthorized(address valuator) external view returns (bool);

    // State-changing functions
    function updateValuation(
        uint256 assetId,
        uint256 newValue,
        ValuationMethod method,
        string calldata reportHash,
        uint256 confidence
    ) external;

    function configureOracle(
        uint256 assetId,
        address oracleAddress,
        uint256 heartbeat,
        uint256 deviationThreshold
    ) external;

    function authorizeValuator(address valuator, bool authorized) external;

    function requestRevaluation(uint256 assetId) external;
}
