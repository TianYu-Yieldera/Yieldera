// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title IGMXV2ExchangeRouter
 * @notice GMX V2 交易路由接口
 * @dev 基于 GMX V2 实际合约简化版本
 */
interface IGMXV2ExchangeRouter {
    struct CreateOrderParams {
        address[] addresses;  // [receiver, callbackContract, uiFeeReceiver, market, initialCollateralToken, swapPath]
        uint256[] numbers;    // [sizeDeltaUsd, initialCollateralDeltaAmount, triggerPrice, acceptablePrice, executionFee, callbackGasLimit, minOutputAmount]
        uint8 orderType;      // 0: MarketIncrease, 1: LimitIncrease, 2: MarketDecrease, 3: LimitDecrease, etc.
        bool isLong;
        bool shouldUnwrapNativeToken;
    }

    /**
     * @notice 创建订单
     * @param params 订单参数
     * @return orderKey 订单唯一标识
     */
    function createOrder(CreateOrderParams calldata params) external payable returns (bytes32 orderKey);

    /**
     * @notice 取消订单
     * @param key 订单 key
     */
    function cancelOrder(bytes32 key) external;
}

/**
 * @title IGMXV2Reader
 * @notice GMX V2 数据读取接口
 */
interface IGMXV2Reader {
    struct MarketPrices {
        uint256 indexTokenPrice;
        uint256 longTokenPrice;
        uint256 shortTokenPrice;
    }

    struct PositionInfo {
        address account;
        address market;
        address collateralToken;
        bool isLong;
        uint256 sizeInUsd;
        uint256 sizeInTokens;
        uint256 collateralAmount;
        uint256 borrowingFactor;
        uint256 fundingFeeAmountPerSize;
        uint256 longTokenClaimableFundingAmountPerSize;
        uint256 shortTokenClaimableFundingAmountPerSize;
    }

    /**
     * @notice 获取账户仓位
     * @param dataStore 数据存储合约地址
     * @param account 用户地址
     * @param market 市场地址
     * @param collateralToken 抵押品代币
     * @param isLong 是否做多
     */
    function getPosition(
        address dataStore,
        address account,
        address market,
        address collateralToken,
        bool isLong
    ) external view returns (PositionInfo memory);

    /**
     * @notice 获取市场价格
     */
    function getMarketTokenPrice(
        address dataStore,
        address market,
        uint256 longTokenPrice,
        uint256 shortTokenPrice,
        uint256 indexTokenPrice,
        bytes32 pnlFactorType,
        bool maximize
    ) external view returns (int256, MarketPrices memory);
}

/**
 * @title IGMXV2DataStore
 * @notice GMX V2 数据存储接口
 */
interface IGMXV2DataStore {
    function getUint(bytes32 key) external view returns (uint256);
    function getInt(bytes32 key) external view returns (int256);
    function getAddress(bytes32 key) external view returns (address);
    function getBool(bytes32 key) external view returns (bool);
}
