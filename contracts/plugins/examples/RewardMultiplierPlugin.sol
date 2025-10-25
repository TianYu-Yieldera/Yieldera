// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../core/BasePlugin.sol";

/**
 * @title RewardMultiplierPlugin
 * @notice Plugin that applies multipliers to user rewards based on behavior
 * @dev Example plugin demonstrating the plugin system capabilities
 *
 * Features:
 * - Dynamic reward multipliers based on user activity
 * - Tiered multiplier system (Bronze, Silver, Gold, Platinum)
 * - Decay mechanism to encourage continued engagement
 * - Admin controls for multiplier configuration
 */
contract RewardMultiplierPlugin is BasePlugin {
    // ============ Constants ============

    bytes32 public constant override PLUGIN_ID = keccak256("REWARD_MULTIPLIER_PLUGIN");
    string public constant PLUGIN_NAME = "Reward Multiplier Plugin";
    string public constant PLUGIN_VERSION = "1.0.0";

    uint256 private constant BASIS_POINTS = 10000;

    // ============ Enums ============

    enum Tier {
        NONE,
        BRONZE,
        SILVER,
        GOLD,
        PLATINUM
    }

    // ============ Structs ============

    struct UserMultiplier {
        Tier tier;
        uint256 multiplier; // In basis points (10000 = 1x, 15000 = 1.5x)
        uint256 activityScore;
        uint256 lastActivity;
        uint256 totalRewardsEarned;
    }

    struct TierConfig {
        uint256 minActivityScore;
        uint256 baseMultiplier;
        uint256 decayRate; // Per day in basis points
    }

    // ============ Storage ============

    mapping(address => UserMultiplier) private _userMultipliers;
    mapping(Tier => TierConfig) private _tierConfigs;

    uint256 private _decayPeriod = 1 days;
    uint256 private _maxMultiplier = 30000; // 3x max
    uint256 private _minMultiplier = 10000; // 1x min

    // ============ Events ============

    event MultiplierUpdated(address indexed user, Tier tier, uint256 multiplier);
    event TierUpgraded(address indexed user, Tier oldTier, Tier newTier);
    event ActivityRecorded(address indexed user, uint256 activityScore, uint256 newTotal);
    event RewardCalculated(address indexed user, uint256 baseReward, uint256 finalReward);

    // ============ Constructor ============

    constructor(address author) BasePlugin(author) {
        _initializeTierConfigs();
    }

    // ============ Plugin Metadata ============

    function getPluginId() external pure override returns (bytes32) {
        return PLUGIN_ID;
    }

    function getPluginName() external pure override returns (string memory) {
        return PLUGIN_NAME;
    }

    function getPluginVersion() external pure override returns (string memory) {
        return PLUGIN_VERSION;
    }

    function getPluginType() external pure override returns (PluginType) {
        return PluginType.REWARDS;
    }

    function getRequiredPermissions() external pure override returns (bytes32[] memory) {
        bytes32[] memory permissions = new bytes32[](2);
        permissions[0] = keccak256("READ_VAULT");
        permissions[1] = keccak256("READ_REWARDS");
        return permissions;
    }

    function getDependencies() external pure override returns (bytes32[] memory) {
        bytes32[] memory deps = new bytes32[](1);
        deps[0] = keccak256("VAULT_MODULE");
        return deps;
    }

    function getDescription() external pure override returns (string memory) {
        return "Applies dynamic reward multipliers based on user activity and tier status";
    }

    function getDocumentationURI() external pure override returns (string memory) {
        return "https://docs.example.com/plugins/reward-multiplier";
    }

    function getSourceURI() external pure override returns (string memory) {
        return "https://github.com/example/plugins/reward-multiplier";
    }

    // ============ Base Plugin Overrides ============

    function _onInitialize(bytes calldata data) internal override {
        // Decode initialization data if any
        if (data.length > 0) {
            // Can configure initial settings here
        }
    }

    function _execute(bytes calldata data) internal override returns (bytes memory result) {
        // Decode the operation type
        (string memory operation, bytes memory params) = abi.decode(data, (string, bytes));

        if (keccak256(bytes(operation)) == keccak256("calculateReward")) {
            return _handleCalculateReward(params);
        } else if (keccak256(bytes(operation)) == keccak256("recordActivity")) {
            return _handleRecordActivity(params);
        } else if (keccak256(bytes(operation)) == keccak256("getMultiplier")) {
            return _handleGetMultiplier(params);
        } else if (keccak256(bytes(operation)) == keccak256("getUserInfo")) {
            return _handleGetUserInfo(params);
        } else {
            revert("Unknown operation");
        }
    }

    function _healthCheck() internal view override returns (bool, string memory) {
        // Check if tier configs are properly set
        if (_tierConfigs[Tier.BRONZE].baseMultiplier == 0) {
            return (false, "Tier configs not initialized");
        }
        return (true, "Plugin healthy");
    }

    // ============ Core Functionality ============

    /**
     * @notice Calculate reward with multiplier applied
     * @param params ABI encoded (address user, uint256 baseReward)
     * @return result ABI encoded uint256 finalReward
     */
    function _handleCalculateReward(bytes memory params) internal returns (bytes memory result) {
        (address user, uint256 baseReward) = abi.decode(params, (address, uint256));

        // Update user's multiplier (apply decay)
        _updateMultiplier(user);

        // Get current multiplier
        uint256 multiplier = _userMultipliers[user].multiplier;
        if (multiplier == 0) {
            multiplier = BASIS_POINTS; // Default 1x
        }

        // Calculate final reward
        uint256 finalReward = (baseReward * multiplier) / BASIS_POINTS;

        // Update stats
        _userMultipliers[user].totalRewardsEarned += finalReward;
        _userMultipliers[user].lastActivity = block.timestamp;

        emit RewardCalculated(user, baseReward, finalReward);

        return abi.encode(finalReward);
    }

    /**
     * @notice Record user activity and update scores
     * @param params ABI encoded (address user, uint256 activityPoints)
     */
    function _handleRecordActivity(bytes memory params) internal returns (bytes memory) {
        (address user, uint256 activityPoints) = abi.decode(params, (address, uint256));

        UserMultiplier storage userMult = _userMultipliers[user];

        // Update activity score
        userMult.activityScore += activityPoints;
        userMult.lastActivity = block.timestamp;

        // Check for tier upgrade
        Tier newTier = _calculateTier(userMult.activityScore);
        if (newTier != userMult.tier) {
            Tier oldTier = userMult.tier;
            userMult.tier = newTier;
            emit TierUpgraded(user, oldTier, newTier);
        }

        // Update multiplier based on new tier
        _updateMultiplier(user);

        emit ActivityRecorded(user, activityPoints, userMult.activityScore);

        return abi.encode(userMult.activityScore, userMult.multiplier);
    }

    /**
     * @notice Get user's current multiplier
     * @param params ABI encoded address user
     */
    function _handleGetMultiplier(bytes memory params) internal view returns (bytes memory) {
        address user = abi.decode(params, (address));
        uint256 multiplier = _userMultipliers[user].multiplier;
        if (multiplier == 0) multiplier = BASIS_POINTS;
        return abi.encode(multiplier);
    }

    /**
     * @notice Get complete user information
     * @param params ABI encoded address user
     */
    function _handleGetUserInfo(bytes memory params) internal view returns (bytes memory) {
        address user = abi.decode(params, (address));
        return abi.encode(_userMultipliers[user]);
    }

    // ============ Internal Functions ============

    function _initializeTierConfigs() internal {
        _tierConfigs[Tier.BRONZE] = TierConfig({
            minActivityScore: 100,
            baseMultiplier: 11000, // 1.1x
            decayRate: 50 // 0.5% per day
        });

        _tierConfigs[Tier.SILVER] = TierConfig({
            minActivityScore: 500,
            baseMultiplier: 12500, // 1.25x
            decayRate: 40 // 0.4% per day
        });

        _tierConfigs[Tier.GOLD] = TierConfig({
            minActivityScore: 2000,
            baseMultiplier: 15000, // 1.5x
            decayRate: 30 // 0.3% per day
        });

        _tierConfigs[Tier.PLATINUM] = TierConfig({
            minActivityScore: 10000,
            baseMultiplier: 20000, // 2x
            decayRate: 20 // 0.2% per day
        });
    }

    function _updateMultiplier(address user) internal {
        UserMultiplier storage userMult = _userMultipliers[user];

        // Calculate decay if time has passed
        if (userMult.lastActivity > 0) {
            uint256 timePassed = block.timestamp - userMult.lastActivity;
            uint256 periodsPassed = timePassed / _decayPeriod;

            if (periodsPassed > 0 && userMult.multiplier > _minMultiplier) {
                TierConfig memory config = _tierConfigs[userMult.tier];
                uint256 decay = (userMult.multiplier * config.decayRate * periodsPassed) / BASIS_POINTS;

                if (userMult.multiplier > decay + _minMultiplier) {
                    userMult.multiplier -= decay;
                } else {
                    userMult.multiplier = _minMultiplier;
                }
            }
        }

        // Set multiplier based on tier if not set
        if (userMult.multiplier == 0 && userMult.tier != Tier.NONE) {
            userMult.multiplier = _tierConfigs[userMult.tier].baseMultiplier;
        }

        // Ensure within bounds
        if (userMult.multiplier > _maxMultiplier) {
            userMult.multiplier = _maxMultiplier;
        }
        if (userMult.multiplier < _minMultiplier && userMult.tier != Tier.NONE) {
            userMult.multiplier = _minMultiplier;
        }

        emit MultiplierUpdated(user, userMult.tier, userMult.multiplier);
    }

    function _calculateTier(uint256 activityScore) internal view returns (Tier) {
        if (activityScore >= _tierConfigs[Tier.PLATINUM].minActivityScore) {
            return Tier.PLATINUM;
        } else if (activityScore >= _tierConfigs[Tier.GOLD].minActivityScore) {
            return Tier.GOLD;
        } else if (activityScore >= _tierConfigs[Tier.SILVER].minActivityScore) {
            return Tier.SILVER;
        } else if (activityScore >= _tierConfigs[Tier.BRONZE].minActivityScore) {
            return Tier.BRONZE;
        } else {
            return Tier.NONE;
        }
    }

    // ============ Admin Functions ============

    function configureTier(
        Tier tier,
        uint256 minActivityScore,
        uint256 baseMultiplier,
        uint256 decayRate
    ) external onlyOwner {
        require(tier != Tier.NONE, "Cannot configure NONE tier");
        require(baseMultiplier >= BASIS_POINTS, "Multiplier too low");
        require(baseMultiplier <= _maxMultiplier, "Multiplier too high");

        _tierConfigs[tier] = TierConfig({
            minActivityScore: minActivityScore,
            baseMultiplier: baseMultiplier,
            decayRate: decayRate
        });
    }

    function setDecayPeriod(uint256 period) external onlyOwner {
        require(period > 0, "Invalid period");
        _decayPeriod = period;
    }

    function setMultiplierBounds(uint256 min, uint256 max) external onlyOwner {
        require(min >= BASIS_POINTS, "Min too low");
        require(max > min, "Max must be greater than min");
        _minMultiplier = min;
        _maxMultiplier = max;
    }

    // ============ View Functions ============

    function getTierConfig(Tier tier) external view returns (TierConfig memory) {
        return _tierConfigs[tier];
    }

    function getUserMultiplier(address user) external view returns (UserMultiplier memory) {
        return _userMultipliers[user];
    }

    function getDecayPeriod() external view returns (uint256) {
        return _decayPeriod;
    }

    function getMultiplierBounds() external view returns (uint256 min, uint256 max) {
        return (_minMultiplier, _maxMultiplier);
    }
}
