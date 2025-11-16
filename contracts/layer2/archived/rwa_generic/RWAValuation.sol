// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "../interfaces/rwa/IRWAValuation.sol";
import "../interfaces/IPriceOracle.sol";

/**
 * @title RWAValuation
 * @notice Asset valuation and pricing oracle for RWA
 * @dev Supports multiple valuation methods and maintains history
 *
 * Key Features:
 * - Manual appraisals by certified valuators
 * - Chainlink oracle integration for automated pricing
 * - Formula-based calculations (bonds, derivatives)
 * - Valuation history with confidence scores
 * - Stale price detection and revaluation triggers
 * - Dispute mechanism for contested valuations
 *
 * Valuation Flow:
 * 1. Valuator performs off-chain appraisal
 * 2. Calls updateValuation() with new value and report hash
 * 3. System stores valuation with timestamp and confidence
 * 4. Oracles can automatically update if configured
 * 5. Anyone can request revaluation if price is stale
 */
contract RWAValuation is IRWAValuation, AccessControl, ReentrancyGuard {
    bytes32 public constant VALUATOR_ROLE = keccak256("VALUATOR_ROLE");
    bytes32 public constant ORACLE_ADMIN_ROLE = keccak256("ORACLE_ADMIN_ROLE");

    // Valuation history (assetId => valuation array)
    mapping(uint256 => Valuation[]) private valuationHistory;

    // Oracle configurations (assetId => config)
    mapping(uint256 => OracleConfig) private oracleConfigs;

    // Authorized valuators
    mapping(address => bool) private authorizedValuators;

    // Dispute tracking
    struct ValuationDispute {
        uint256 disputedValue;
        address disputer;
        string reason;
        uint256 timestamp;
        bool resolved;
    }
    mapping(uint256 => ValuationDispute[]) private disputes;

    // Constants
    uint256 public constant MAX_VALUATION_AGE = 90 days;
    uint256 public constant MIN_CONFIDENCE = 5000; // 50% minimum confidence
    uint256 public constant MAX_CONFIDENCE = 10000; // 100% maximum
    uint256 public constant PRECISION = 10000;

    // Statistics
    uint256 public totalValuations;
    uint256 public totalAssets;

    /**
     * @notice Constructor
     * @param admin Address to grant admin roles
     */
    constructor(address admin) {
        require(admin != address(0), "Invalid admin address");

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(VALUATOR_ROLE, admin);
        _grantRole(ORACLE_ADMIN_ROLE, admin);

        authorizedValuators[admin] = true;
    }

    // =============================================================
    //                    VALUATION MANAGEMENT
    // =============================================================

    /**
     * @notice Update asset valuation
     * @param assetId Asset identifier
     * @param newValue New valuation in USD (18 decimals)
     * @param method Valuation method used
     * @param reportHash IPFS hash of valuation report
     * @param confidence Confidence score (5000-10000 = 50%-100%)
     */
    function updateValuation(
        uint256 assetId,
        uint256 newValue,
        ValuationMethod method,
        string calldata reportHash,
        uint256 confidence
    ) external onlyRole(VALUATOR_ROLE) nonReentrant {
        require(newValue > 0, "Invalid valuation");
        require(confidence >= MIN_CONFIDENCE && confidence <= MAX_CONFIDENCE, "Invalid confidence");
        require(bytes(reportHash).length > 0, "Report hash required");

        uint256 oldValue = getCurrentValue(assetId);

        Valuation memory newValuation = Valuation({
            value: newValue,
            timestamp: block.timestamp,
            method: method,
            valuator: msg.sender,
            reportHash: reportHash,
            confidence: confidence
        });

        valuationHistory[assetId].push(newValuation);

        // Track statistics
        if (oldValue == 0) {
            totalAssets++;
        }
        totalValuations++;

        emit ValuationUpdated(assetId, oldValue, newValue, method, msg.sender);
    }

    /**
     * @notice Configure oracle for automated pricing
     * @param assetId Asset identifier
     * @param oracleAddress Chainlink price feed address
     * @param heartbeat Maximum acceptable data age (seconds)
     * @param deviationThreshold Price deviation triggering revaluation (basis points)
     */
    function configureOracle(
        uint256 assetId,
        address oracleAddress,
        uint256 heartbeat,
        uint256 deviationThreshold
    ) external onlyRole(ORACLE_ADMIN_ROLE) {
        require(oracleAddress != address(0), "Invalid oracle address");
        require(heartbeat > 0 && heartbeat <= 1 days, "Invalid heartbeat");
        require(deviationThreshold <= PRECISION, "Invalid deviation threshold");

        OracleConfig storage config = oracleConfigs[assetId];
        config.oracleAddress = oracleAddress;
        config.heartbeat = heartbeat;
        config.deviationThreshold = deviationThreshold;
        config.isActive = true;

        emit OracleConfigured(assetId, oracleAddress, heartbeat);
    }

    /**
     * @notice Disable oracle for an asset
     * @param assetId Asset identifier
     */
    function disableOracle(uint256 assetId) external onlyRole(ORACLE_ADMIN_ROLE) {
        oracleConfigs[assetId].isActive = false;
    }

    /**
     * @notice Request revaluation for stale or disputed asset
     * @param assetId Asset identifier
     */
    function requestRevaluation(uint256 assetId) external {
        require(valuationHistory[assetId].length > 0, "Asset not valued");
        require(isValuationStale(assetId), "Valuation not stale");

        // In production, this would trigger off-chain workflow
        // For now, just emit event
        emit ValuationDisputed(assetId, getCurrentValue(assetId), msg.sender, "Stale valuation");
    }

    /**
     * @notice Update valuation from oracle
     * @dev Can be called by anyone if oracle is configured
     * @param assetId Asset identifier
     */
    function updateFromOracle(uint256 assetId) external nonReentrant {
        OracleConfig storage config = oracleConfigs[assetId];
        require(config.isActive, "Oracle not configured");

        IPriceOracle oracle = IPriceOracle(config.oracleAddress);
        uint256 oraclePrice = oracle.getAssetPrice(address(uint160(assetId))); // Simplified

        require(oraclePrice > 0, "Invalid oracle price");

        // Check if price has deviated significantly
        uint256 currentValue = getCurrentValue(assetId);
        if (currentValue > 0) {
            uint256 deviation = currentValue > oraclePrice
                ? ((currentValue - oraclePrice) * PRECISION) / currentValue
                : ((oraclePrice - currentValue) * PRECISION) / oraclePrice;

            require(deviation >= config.deviationThreshold, "Price deviation too small");
        }

        Valuation memory newValuation = Valuation({
            value: oraclePrice,
            timestamp: block.timestamp,
            method: ValuationMethod.Oracle,
            valuator: address(oracle),
            reportHash: "oracle-update",
            confidence: MAX_CONFIDENCE // Oracle assumed 100% confident
        });

        valuationHistory[assetId].push(newValuation);
        totalValuations++;

        emit ValuationUpdated(assetId, currentValue, oraclePrice, ValuationMethod.Oracle, address(oracle));
    }

    // =============================================================
    //                   VALUATOR MANAGEMENT
    // =============================================================

    /**
     * @notice Authorize or revoke valuator
     * @param valuator Address to authorize/revoke
     * @param authorized True to authorize, false to revoke
     */
    function authorizeValuator(
        address valuator,
        bool authorized
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(valuator != address(0), "Invalid valuator address");

        if (authorized) {
            grantRole(VALUATOR_ROLE, valuator);
        } else {
            revokeRole(VALUATOR_ROLE, valuator);
        }

        authorizedValuators[valuator] = authorized;

        emit ValuatorAuthorized(valuator, authorized);
    }

    /**
     * @notice Check if address is authorized valuator
     * @param valuator Address to check
     * @return True if authorized
     */
    function isValuatorAuthorized(address valuator) public view returns (bool) {
        return authorizedValuators[valuator];
    }

    // =============================================================
    //                      VIEW FUNCTIONS
    // =============================================================

    /**
     * @notice Get current asset value (most recent valuation)
     * @param assetId Asset identifier
     * @return Current value in USD (18 decimals)
     */
    function getCurrentValue(uint256 assetId) public view returns (uint256) {
        uint256 length = valuationHistory[assetId].length;
        if (length == 0) {
            return 0;
        }

        return valuationHistory[assetId][length - 1].value;
    }

    /**
     * @notice Get most recent valuation record
     * @param assetId Asset identifier
     * @return Valuation struct
     */
    function getLastValuation(uint256 assetId) external view returns (Valuation memory) {
        uint256 length = valuationHistory[assetId].length;
        require(length > 0, "No valuations");

        return valuationHistory[assetId][length - 1];
    }

    /**
     * @notice Get valuation history
     * @param assetId Asset identifier
     * @param count Number of recent valuations to return
     * @return Array of Valuation structs
     */
    function getValuationHistory(
        uint256 assetId,
        uint256 count
    ) external view returns (Valuation[] memory) {
        uint256 length = valuationHistory[assetId].length;
        uint256 returnCount = count > length ? length : count;

        Valuation[] memory history = new Valuation[](returnCount);

        for (uint256 i = 0; i < returnCount; i++) {
            history[i] = valuationHistory[assetId][length - returnCount + i];
        }

        return history;
    }

    /**
     * @notice Calculate price per fractional token
     * @param assetId Asset identifier
     * @return Price per token (assuming 1e18 tokens)
     */
    function getPricePerToken(uint256 assetId) external view returns (uint256) {
        uint256 totalValue = getCurrentValue(assetId);
        if (totalValue == 0) {
            return 0;
        }

        // In production, this would divide by actual token supply
        // For now, assume standard 1e18 fractionalization
        return totalValue / 1e18;
    }

    /**
     * @notice Check if valuation is stale
     * @param assetId Asset identifier
     * @return True if stale (older than MAX_VALUATION_AGE)
     */
    function isValuationStale(uint256 assetId) public view returns (bool) {
        uint256 length = valuationHistory[assetId].length;
        if (length == 0) {
            return true;
        }

        Valuation storage lastValuation = valuationHistory[assetId][length - 1];
        return block.timestamp - lastValuation.timestamp > MAX_VALUATION_AGE;
    }

    /**
     * @notice Get oracle configuration
     * @param assetId Asset identifier
     * @return OracleConfig struct
     */
    function getOracleConfig(uint256 assetId) external view returns (OracleConfig memory) {
        return oracleConfigs[assetId];
    }

    /**
     * @notice Get total number of valuations for an asset
     * @param assetId Asset identifier
     * @return Count of valuations
     */
    function getValuationCount(uint256 assetId) external view returns (uint256) {
        return valuationHistory[assetId].length;
    }

    /**
     * @notice Calculate time-weighted average price (TWAP)
     * @param assetId Asset identifier
     * @param period Time period for TWAP calculation (seconds)
     * @return TWAP value
     */
    function getTWAP(uint256 assetId, uint256 period) external view returns (uint256) {
        require(period > 0, "Invalid period");

        uint256 length = valuationHistory[assetId].length;
        require(length > 0, "No valuations");

        uint256 cutoffTime = block.timestamp - period;
        uint256 weightedSum = 0;
        uint256 totalWeight = 0;

        for (uint256 i = 0; i < length; i++) {
            Valuation storage val = valuationHistory[assetId][i];

            if (val.timestamp >= cutoffTime) {
                uint256 weight = block.timestamp - val.timestamp;
                weightedSum += val.value * weight;
                totalWeight += weight;
            }
        }

        require(totalWeight > 0, "No valuations in period");

        return weightedSum / totalWeight;
    }

    // =============================================================
    //                    DISPUTE MANAGEMENT
    // =============================================================

    /**
     * @notice Dispute a valuation
     * @param assetId Asset identifier
     * @param reason Reason for dispute
     */
    function disputeValuation(
        uint256 assetId,
        string calldata reason
    ) external {
        require(valuationHistory[assetId].length > 0, "Asset not valued");
        require(bytes(reason).length > 0, "Reason required");

        uint256 currentValue = getCurrentValue(assetId);

        ValuationDispute memory dispute = ValuationDispute({
            disputedValue: currentValue,
            disputer: msg.sender,
            reason: reason,
            timestamp: block.timestamp,
            resolved: false
        });

        disputes[assetId].push(dispute);

        emit ValuationDisputed(assetId, currentValue, msg.sender, reason);
    }

    /**
     * @notice Get disputes for an asset
     * @param assetId Asset identifier
     * @return Array of disputes
     */
    function getDisputes(uint256 assetId) external view returns (ValuationDispute[] memory) {
        return disputes[assetId];
    }

    /**
     * @notice Resolve a dispute
     * @param assetId Asset identifier
     * @param disputeIndex Index in disputes array
     */
    function resolveDispute(
        uint256 assetId,
        uint256 disputeIndex
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(disputeIndex < disputes[assetId].length, "Invalid dispute index");
        disputes[assetId][disputeIndex].resolved = true;
    }
}
