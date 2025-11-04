// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "../interfaces/treasury/ITreasuryAsset.sol";

/**
 * @title TreasuryPriceOracle
 * @notice Price oracle for US Treasury securities
 * @dev Stores and provides price and yield data for treasury assets
 *
 * Key Features:
 * - Store current price and yield for each asset
 * - Historical price tracking
 * - Price update validation
 * - Support for multiple price sources
 * - Price staleness checks
 */
contract TreasuryPriceOracle is AccessControl, Pausable {
    bytes32 public constant ORACLE_ROLE = keccak256("ORACLE_ROLE");
    bytes32 public constant PRICE_UPDATER_ROLE = keccak256("PRICE_UPDATER_ROLE");

    /// @notice Price data structure
    struct PriceData {
        uint256 price;              // Current price in USD (18 decimals)
        uint256 yield;              // Current yield in basis points (425 = 4.25%)
        uint256 timestamp;          // Last update timestamp
        uint256 updateCount;        // Number of updates
        string source;              // Data source identifier
    }

    /// @notice Price history entry
    struct PriceHistory {
        uint256 price;
        uint256 yield;
        uint256 timestamp;
        string source;
    }

    /// @notice Price validation parameters
    uint256 public constant MAX_PRICE_DEVIATION = 1000; // 10% max deviation (basis points)
    uint256 public constant PRICE_FRESHNESS = 24 hours;  // Price considered stale after 24h
    uint256 public constant MAX_YIELD = 10000;           // 100% max yield (basis points)

    /// @notice Storage
    mapping(uint256 => PriceData) public assetPrices;     // assetId => current price
    mapping(uint256 => PriceHistory[]) public priceHistory; // assetId => price history
    mapping(uint256 => uint256) public lastPriceUpdate;   // assetId => timestamp

    /// @notice Integration contracts
    ITreasuryAsset public immutable assetFactory;

    /// @notice Statistics
    uint256 public totalPriceUpdates;

    /// @notice Events
    event PriceUpdated(
        uint256 indexed assetId,
        uint256 price,
        uint256 yield,
        uint256 timestamp,
        string source
    );

    event PriceHistoryRecorded(
        uint256 indexed assetId,
        uint256 price,
        uint256 yield,
        uint256 timestamp
    );

    event PriceSourceChanged(
        uint256 indexed assetId,
        string oldSource,
        string newSource
    );

    /**
     * @notice Constructor
     * @param admin Admin address
     * @param factory TreasuryAssetFactory address
     */
    constructor(address admin, address factory) {
        require(admin != address(0), "Invalid admin");
        require(factory != address(0), "Invalid factory");

        assetFactory = ITreasuryAsset(factory);

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(ORACLE_ROLE, admin);
        _grantRole(PRICE_UPDATER_ROLE, admin);
    }

    // =============================================================
    //                     PRICE UPDATES
    // =============================================================

    /**
     * @notice Update price for treasury asset
     * @param assetId Asset identifier
     * @param newPrice New price in USD (18 decimals)
     * @param newYield New yield in basis points
     * @param source Data source identifier
     */
    function updatePrice(
        uint256 assetId,
        uint256 newPrice,
        uint256 newYield,
        string memory source
    ) external onlyRole(PRICE_UPDATER_ROLE) whenNotPaused {
        require(assetFactory.isAssetActive(assetId), "Asset not active");
        require(newPrice > 0, "Invalid price");
        require(newYield <= MAX_YIELD, "Yield too high");
        require(bytes(source).length > 0, "Invalid source");

        PriceData storage priceData = assetPrices[assetId];

        // Validate price deviation if not first update
        if (priceData.timestamp > 0) {
            uint256 priceDiff = newPrice > priceData.price
                ? newPrice - priceData.price
                : priceData.price - newPrice;

            uint256 deviationBps = (priceDiff * 10000) / priceData.price;
            require(deviationBps <= MAX_PRICE_DEVIATION, "Price deviation too high");
        }

        // Record in history
        priceHistory[assetId].push(PriceHistory({
            price: newPrice,
            yield: newYield,
            timestamp: block.timestamp,
            source: source
        }));

        // Update current price
        priceData.price = newPrice;
        priceData.yield = newYield;
        priceData.timestamp = block.timestamp;
        priceData.updateCount++;
        priceData.source = source;

        lastPriceUpdate[assetId] = block.timestamp;
        totalPriceUpdates++;

        emit PriceUpdated(assetId, newPrice, newYield, block.timestamp, source);
        emit PriceHistoryRecorded(assetId, newPrice, newYield, block.timestamp);
    }

    /**
     * @notice Batch update prices
     * @param assetIds Array of asset IDs
     * @param prices Array of prices
     * @param yields Array of yields
     * @param source Data source
     */
    function batchUpdatePrices(
        uint256[] calldata assetIds,
        uint256[] calldata prices,
        uint256[] calldata yields,
        string memory source
    ) external onlyRole(PRICE_UPDATER_ROLE) whenNotPaused {
        require(assetIds.length == prices.length, "Length mismatch");
        require(assetIds.length == yields.length, "Length mismatch");

        for (uint256 i = 0; i < assetIds.length; i++) {
            _updatePriceInternal(assetIds[i], prices[i], yields[i], source);
        }
    }

    /**
     * @notice Internal price update (skip external checks)
     */
    function _updatePriceInternal(
        uint256 assetId,
        uint256 newPrice,
        uint256 newYield,
        string memory source
    ) private {
        require(newPrice > 0, "Invalid price");
        require(newYield <= MAX_YIELD, "Yield too high");

        PriceData storage priceData = assetPrices[assetId];

        // Record in history
        priceHistory[assetId].push(PriceHistory({
            price: newPrice,
            yield: newYield,
            timestamp: block.timestamp,
            source: source
        }));

        // Update current price
        priceData.price = newPrice;
        priceData.yield = newYield;
        priceData.timestamp = block.timestamp;
        priceData.updateCount++;
        priceData.source = source;

        lastPriceUpdate[assetId] = block.timestamp;
        totalPriceUpdates++;

        emit PriceUpdated(assetId, newPrice, newYield, block.timestamp, source);
    }

    // =============================================================
    //                      VIEW FUNCTIONS
    // =============================================================

    /**
     * @notice Get latest price and yield
     * @param assetId Asset identifier
     * @return price Current price
     * @return yield Current yield
     * @return timestamp Last update timestamp
     */
    function getLatestPrice(uint256 assetId)
        external
        view
        returns (
            uint256 price,
            uint256 yield,
            uint256 timestamp
        )
    {
        PriceData storage priceData = assetPrices[assetId];
        return (priceData.price, priceData.yield, priceData.timestamp);
    }

    /**
     * @notice Get full price data
     * @param assetId Asset identifier
     * @return Price data struct
     */
    function getPriceData(uint256 assetId) external view returns (PriceData memory) {
        return assetPrices[assetId];
    }

    /**
     * @notice Get price history
     * @param assetId Asset identifier
     * @param count Number of recent entries to return
     * @return Array of price history entries
     */
    function getPriceHistory(uint256 assetId, uint256 count)
        external
        view
        returns (PriceHistory[] memory)
    {
        PriceHistory[] storage history = priceHistory[assetId];
        uint256 returnCount = count > history.length ? history.length : count;

        PriceHistory[] memory result = new PriceHistory[](returnCount);

        for (uint256 i = 0; i < returnCount; i++) {
            result[i] = history[history.length - returnCount + i];
        }

        return result;
    }

    /**
     * @notice Check if price is fresh (updated within freshness period)
     * @param assetId Asset identifier
     * @return True if fresh
     */
    function isPriceFresh(uint256 assetId) external view returns (bool) {
        PriceData storage priceData = assetPrices[assetId];
        if (priceData.timestamp == 0) return false;

        return (block.timestamp - priceData.timestamp) <= PRICE_FRESHNESS;
    }

    /**
     * @notice Get price age in seconds
     * @param assetId Asset identifier
     * @return Age in seconds
     */
    function getPriceAge(uint256 assetId) external view returns (uint256) {
        PriceData storage priceData = assetPrices[assetId];
        if (priceData.timestamp == 0) return type(uint256).max;

        return block.timestamp - priceData.timestamp;
    }

    /**
     * @notice Get multiple asset prices
     * @param assetIds Array of asset IDs
     * @return prices Array of prices
     * @return yields Array of yields
     * @return timestamps Array of timestamps
     */
    function getBatchPrices(uint256[] calldata assetIds)
        external
        view
        returns (
            uint256[] memory prices,
            uint256[] memory yields,
            uint256[] memory timestamps
        )
    {
        prices = new uint256[](assetIds.length);
        yields = new uint256[](assetIds.length);
        timestamps = new uint256[](assetIds.length);

        for (uint256 i = 0; i < assetIds.length; i++) {
            PriceData storage priceData = assetPrices[assetIds[i]];
            prices[i] = priceData.price;
            yields[i] = priceData.yield;
            timestamps[i] = priceData.timestamp;
        }

        return (prices, yields, timestamps);
    }

    /**
     * @notice Calculate market value for token amount
     * @param assetId Asset identifier
     * @param tokenAmount Token amount
     * @return Market value in USD
     */
    function calculateMarketValue(uint256 assetId, uint256 tokenAmount)
        external
        view
        returns (uint256)
    {
        PriceData storage priceData = assetPrices[assetId];
        require(priceData.price > 0, "Price not set");

        return (tokenAmount * priceData.price) / 1e18;
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

    /**
     * @notice Emergency price update (skip validation)
     * @param assetId Asset identifier
     * @param price New price
     * @param yield New yield
     */
    function emergencyPriceUpdate(
        uint256 assetId,
        uint256 price,
        uint256 yield
    ) external onlyRole(ORACLE_ROLE) {
        require(price > 0, "Invalid price");
        require(yield <= MAX_YIELD, "Yield too high");

        PriceData storage priceData = assetPrices[assetId];
        priceData.price = price;
        priceData.yield = yield;
        priceData.timestamp = block.timestamp;
        priceData.updateCount++;
        priceData.source = "EMERGENCY_UPDATE";

        emit PriceUpdated(assetId, price, yield, block.timestamp, "EMERGENCY_UPDATE");
    }
}
