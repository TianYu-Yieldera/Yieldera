// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../layer2/interfaces/external/IGMXV2.sol";

/**
 * @title MockGMXV2ExchangeRouter
 * @notice GMX V2 交易路由的 Mock 实现 - 用于测试
 */
contract MockGMXV2ExchangeRouter is IGMXV2ExchangeRouter {
    uint256 private orderCounter;
    mapping(bytes32 => bool) public orders;

    event OrderCreated(bytes32 indexed orderKey, address indexed user);
    event OrderCancelled(bytes32 indexed orderKey);

    function createOrder(CreateOrderParams calldata params)
        external
        payable
        override
        returns (bytes32 orderKey)
    {
        orderCounter++;
        orderKey = keccak256(abi.encodePacked(msg.sender, orderCounter, block.timestamp));
        orders[orderKey] = true;

        emit OrderCreated(orderKey, msg.sender);
        return orderKey;
    }

    function cancelOrder(bytes32 key) external override {
        require(orders[key], "Order not found");
        orders[key] = false;
        emit OrderCancelled(key);
    }
}

/**
 * @title MockGMXV2Reader
 * @notice GMX V2 Reader 的 Mock 实现
 */
contract MockGMXV2Reader is IGMXV2Reader {
    mapping(bytes32 => PositionInfo) private positions;

    function getPosition(
        address dataStore,
        address account,
        address market,
        address collateralToken,
        bool isLong
    ) external view override returns (PositionInfo memory) {
        bytes32 key = keccak256(abi.encodePacked(account, market, collateralToken, isLong));
        return positions[key];
    }

    function getMarketTokenPrice(
        address dataStore,
        address market,
        uint256 longTokenPrice,
        uint256 shortTokenPrice,
        uint256 indexTokenPrice,
        bytes32 pnlFactorType,
        bool maximize
    ) external view override returns (int256, MarketPrices memory) {
        MarketPrices memory prices = MarketPrices({
            indexTokenPrice: indexTokenPrice,
            longTokenPrice: longTokenPrice,
            shortTokenPrice: shortTokenPrice
        });
        return (0, prices);
    }

    // 测试辅助函数
    function setPosition(
        address account,
        address market,
        address collateralToken,
        bool isLong,
        uint256 sizeInUsd,
        uint256 collateralAmount
    ) external {
        bytes32 key = keccak256(abi.encodePacked(account, market, collateralToken, isLong));
        positions[key] = PositionInfo({
            account: account,
            market: market,
            collateralToken: collateralToken,
            isLong: isLong,
            sizeInUsd: sizeInUsd,
            sizeInTokens: 0,
            collateralAmount: collateralAmount,
            borrowingFactor: 0,
            fundingFeeAmountPerSize: 0,
            longTokenClaimableFundingAmountPerSize: 0,
            shortTokenClaimableFundingAmountPerSize: 0
        });
    }
}

/**
 * @title MockGMXV2DataStore
 * @notice GMX V2 DataStore 的 Mock 实现
 */
contract MockGMXV2DataStore is IGMXV2DataStore {
    mapping(bytes32 => uint256) private uintValues;
    mapping(bytes32 => int256) private intValues;
    mapping(bytes32 => address) private addressValues;
    mapping(bytes32 => bool) private boolValues;

    function getUint(bytes32 key) external view override returns (uint256) {
        return uintValues[key];
    }

    function getInt(bytes32 key) external view override returns (int256) {
        return intValues[key];
    }

    function getAddress(bytes32 key) external view override returns (address) {
        return addressValues[key];
    }

    function getBool(bytes32 key) external view override returns (bool) {
        return boolValues[key];
    }

    // 测试辅助函数
    function setUint(bytes32 key, uint256 value) external {
        uintValues[key] = value;
    }

    function setInt(bytes32 key, int256 value) external {
        intValues[key] = value;
    }

    function setAddress(bytes32 key, address value) external {
        addressValues[key] = value;
    }

    function setBool(bytes32 key, bool value) external {
        boolValues[key] = value;
    }
}
