// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "../interfaces/external/IGMXV2.sol";

/**
 * @title GMXV2Adapter
 * @notice GMX V2 永续合约适配器 - 用于风险对冲和衍生品交易
 * @dev 支持开仓、平仓、紧急对冲等功能
 *
 * 核心功能:
 * 1. 开仓/平仓 - 做多或做空永续合约
 * 2. 紧急对冲 - 风控系统触发的自动对冲
 * 3. 仓位查询 - 实时查询用户仓位信息
 * 4. 风险管理 - 杠杆限制、滑点保护
 */
contract GMXV2Adapter is AccessControl, ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;

    // ============ 角色定义 ============
    bytes32 public constant RISK_MANAGER_ROLE = keccak256("RISK_MANAGER_ROLE");
    bytes32 public constant OPERATOR_ROLE = keccak256("OPERATOR_ROLE");

    // ============ GMX V2 合约接口 ============
    IGMXV2ExchangeRouter public immutable exchangeRouter;
    IGMXV2Reader public immutable reader;
    IGMXV2DataStore public immutable dataStore;

    // ============ 配置参数 ============
    uint256 public constant MAX_LEVERAGE = 50; // 最大杠杆 50x
    uint256 public constant MIN_EXECUTION_FEE = 0.0001 ether; // 最小执行费用
    uint256 public constant MAX_SLIPPAGE_BPS = 200; // 最大滑点 2%

    // ============ 支持的市场和代币 ============
    mapping(address => bool) public supportedMarkets;      // 支持的交易市场
    mapping(address => bool) public supportedCollateral;   // 支持的抵押品

    // ============ 用户仓位追踪 ============
    struct UserPosition {
        address market;           // 市场地址
        address collateralToken;  // 抵押品代币
        bool isLong;              // 是否做多
        uint256 sizeInUsd;        // 仓位大小 (USD)
        uint256 collateralAmount; // 抵押品数量
        uint256 leverage;         // 杠杆倍数
        uint256 openTimestamp;    // 开仓时间
        bool isHedge;             // 是否为对冲仓位
    }

    mapping(address => UserPosition[]) public userPositions;
    mapping(bytes32 => address) public orderToUser; // orderKey => user

    // ============ 统计数据 ============
    struct Statistics {
        uint256 totalOrders;
        uint256 totalHedges;
        uint256 totalVolume;        // 总交易量 (USD)
        uint256 successfulOrders;
        uint256 failedOrders;
    }
    Statistics public stats;

    // ============ 事件 ============
    event PositionOpened(
        address indexed user,
        bytes32 indexed orderKey,
        address market,
        address collateralToken,
        bool isLong,
        uint256 sizeInUsd,
        uint256 collateralAmount,
        uint256 leverage,
        bool isHedge
    );

    event PositionClosed(
        address indexed user,
        bytes32 indexed orderKey,
        address market,
        uint256 sizeInUsd,
        int256 pnl
    );

    event EmergencyHedgeExecuted(
        address indexed user,
        address indexed market,
        uint256 hedgeSize,
        string reason,
        bytes32 orderKey
    );

    event OrderCancelled(
        address indexed user,
        bytes32 indexed orderKey,
        string reason
    );

    event MarketAdded(address indexed market);
    event CollateralAdded(address indexed token);

    // ============ 错误定义 ============
    error UnsupportedMarket(address market);
    error UnsupportedCollateral(address token);
    error InvalidLeverage(uint256 leverage);
    error InsufficientExecutionFee(uint256 provided, uint256 required);
    error SlippageTooHigh(uint256 slippage);
    error NoPositionFound();
    error Unauthorized();

    // ============ 构造函数 ============
    /**
     * @param _exchangeRouter GMX V2 交易路由地址
     * @param _reader GMX V2 Reader 地址
     * @param _dataStore GMX V2 DataStore 地址
     */
    constructor(
        address _exchangeRouter,
        address _reader,
        address _dataStore
    ) {
        require(_exchangeRouter != address(0), "Invalid exchange router");
        require(_reader != address(0), "Invalid reader");
        require(_dataStore != address(0), "Invalid data store");

        exchangeRouter = IGMXV2ExchangeRouter(_exchangeRouter);
        reader = IGMXV2Reader(_reader);
        dataStore = IGMXV2DataStore(_dataStore);

        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(RISK_MANAGER_ROLE, msg.sender);
        _grantRole(OPERATOR_ROLE, msg.sender);
    }

    // ============ 核心交易功能 ============

    /**
     * @notice 开仓 - 创建做多或做空仓位
     * @param market GMX 市场地址
     * @param collateralToken 抵押品代币地址
     * @param collateralAmount 抵押品数量
     * @param sizeInUsd 仓位大小 (USD, 18 decimals)
     * @param isLong 是否做多
     * @param acceptablePrice 可接受的执行价格
     * @param executionFee 执行费用 (ETH)
     * @return orderKey 订单唯一标识
     */
    function openPosition(
        address market,
        address collateralToken,
        uint256 collateralAmount,
        uint256 sizeInUsd,
        bool isLong,
        uint256 acceptablePrice,
        uint256 executionFee
    ) external payable nonReentrant whenNotPaused returns (bytes32 orderKey) {
        // 验证参数
        if (!supportedMarkets[market]) revert UnsupportedMarket(market);
        if (!supportedCollateral[collateralToken]) revert UnsupportedCollateral(collateralToken);
        if (executionFee < MIN_EXECUTION_FEE) {
            revert InsufficientExecutionFee(executionFee, MIN_EXECUTION_FEE);
        }

        // 验证杠杆并返回杠杆值
        uint256 leverage = _validateAndCalculateLeverage(collateralToken, collateralAmount, sizeInUsd);

        // 转移抵押品到合约
        IERC20(collateralToken).safeTransferFrom(msg.sender, address(this), collateralAmount);

        // 批准 GMX Router 使用抵押品
        IERC20(collateralToken).forceApprove(address(exchangeRouter), collateralAmount);

        // 创建订单
        orderKey = exchangeRouter.createOrder{value: executionFee}(
            IGMXV2ExchangeRouter.CreateOrderParams({
                addresses: _buildAddresses(market, collateralToken, msg.sender),
                numbers: _buildNumbers(sizeInUsd, collateralAmount, acceptablePrice, executionFee),
                orderType: 0, // 0 = MarketIncrease (市价开仓)
                isLong: isLong,
                shouldUnwrapNativeToken: false
            })
        );

        // 记录用户仓位
        userPositions[msg.sender].push(UserPosition({
            market: market,
            collateralToken: collateralToken,
            isLong: isLong,
            sizeInUsd: sizeInUsd,
            collateralAmount: collateralAmount,
            leverage: leverage,
            openTimestamp: block.timestamp,
            isHedge: false
        }));

        orderToUser[orderKey] = msg.sender;

        // 更新统计
        stats.totalOrders++;
        stats.totalVolume += sizeInUsd;

        emit PositionOpened(
            msg.sender,
            orderKey,
            market,
            collateralToken,
            isLong,
            sizeInUsd,
            collateralAmount,
            leverage,
            false
        );

        return orderKey;
    }

    /**
     * @notice 平仓 - 关闭现有仓位
     * @param market 市场地址
     * @param collateralToken 抵押品代币
     * @param sizeInUsd 平仓大小 (USD)
     * @param isLong 是否做多
     * @param acceptablePrice 可接受价格
     * @param executionFee 执行费用
     * @return orderKey 订单 key
     */
    function closePosition(
        address market,
        address collateralToken,
        uint256 sizeInUsd,
        bool isLong,
        uint256 acceptablePrice,
        uint256 executionFee
    ) external payable nonReentrant whenNotPaused returns (bytes32 orderKey) {
        if (!supportedMarkets[market]) revert UnsupportedMarket(market);
        if (executionFee < MIN_EXECUTION_FEE) {
            revert InsufficientExecutionFee(executionFee, MIN_EXECUTION_FEE);
        }

        // 验证用户有仓位
        _verifyUserPosition(msg.sender, market, isLong);

        // 构建平仓订单
        IGMXV2ExchangeRouter.CreateOrderParams memory params = IGMXV2ExchangeRouter.CreateOrderParams({
            addresses: _buildAddresses(market, collateralToken, msg.sender),
            numbers: _buildNumbers(sizeInUsd, 0, acceptablePrice, executionFee),
            orderType: 2, // 2 = MarketDecrease (市价平仓)
            isLong: isLong,
            shouldUnwrapNativeToken: false
        });

        orderKey = exchangeRouter.createOrder{value: executionFee}(params);
        orderToUser[orderKey] = msg.sender;

        stats.totalOrders++;

        emit PositionClosed(msg.sender, orderKey, market, sizeInUsd, 0);

        return orderKey;
    }

    /**
     * @notice 紧急对冲 - 由风控系统触发
     * @dev 只有 RISK_MANAGER 角色可以调用
     * @param user 用户地址
     * @param market 市场地址
     * @param collateralToken 抵押品
     * @param hedgeSize 对冲规模 (USD)
     * @param reason 对冲原因
     * @return orderKey 订单 key
     */
    function emergencyHedge(
        address user,
        address market,
        address collateralToken,
        uint256 hedgeSize,
        string calldata reason
    ) external payable onlyRole(RISK_MANAGER_ROLE) nonReentrant returns (bytes32) {
        if (!supportedMarkets[market]) revert UnsupportedMarket(market);

        // 计算并转移抵押品
        uint256 collateralNeeded = _calculateCollateral(collateralToken, hedgeSize);
        IERC20(collateralToken).safeTransferFrom(user, address(this), collateralNeeded);
        IERC20(collateralToken).forceApprove(address(exchangeRouter), collateralNeeded);

        // 创建对冲订单并记录
        bytes32 orderKey = _createHedgeOrder(user, market, collateralToken, hedgeSize, collateralNeeded);

        stats.totalHedges++;
        emit EmergencyHedgeExecuted(user, market, hedgeSize, reason, orderKey);

        return orderKey;
    }

    /**
     * @notice 取消订单
     * @param orderKey 订单 key
     */
    function cancelOrder(bytes32 orderKey) external nonReentrant {
        address user = orderToUser[orderKey];
        if (user != msg.sender && !hasRole(OPERATOR_ROLE, msg.sender)) {
            revert Unauthorized();
        }

        exchangeRouter.cancelOrder(orderKey);

        emit OrderCancelled(user, orderKey, "User cancelled");
    }

    // ============ 查询功能 ============

    /**
     * @notice 获取用户在指定市场的仓位信息
     * @param user 用户地址
     * @param market 市场地址
     * @param collateralToken 抵押品代币
     * @param isLong 是否做多
     * @return position 仓位信息
     */
    function getPosition(
        address user,
        address market,
        address collateralToken,
        bool isLong
    ) external view returns (IGMXV2Reader.PositionInfo memory position) {
        return reader.getPosition(
            address(dataStore),
            user,
            market,
            collateralToken,
            isLong
        );
    }

    /**
     * @notice 获取用户所有仓位
     * @param user 用户地址
     * @return positions 仓位数组
     */
    function getUserPositions(address user) external view returns (UserPosition[] memory) {
        return userPositions[user];
    }

    /**
     * @notice 获取统计数据
     */
    function getStatistics() external view returns (Statistics memory) {
        return stats;
    }

    // ============ 管理功能 ============

    /**
     * @notice 添加支持的市场
     * @param market 市场地址
     */
    function addMarket(address market) external onlyRole(DEFAULT_ADMIN_ROLE) {
        supportedMarkets[market] = true;
        emit MarketAdded(market);
    }

    /**
     * @notice 添加支持的抵押品
     * @param token 代币地址
     */
    function addCollateral(address token) external onlyRole(DEFAULT_ADMIN_ROLE) {
        supportedCollateral[token] = true;
        emit CollateralAdded(token);
    }

    /**
     * @notice 暂停合约
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @notice 恢复合约
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }

    /**
     * @notice 紧急提取代币
     * @param token 代币地址
     * @param amount 数量
     */
    function emergencyWithdraw(
        address token,
        uint256 amount
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        IERC20(token).safeTransfer(msg.sender, amount);
    }

    // ============ 内部辅助函数 ============

    /**
     * @notice 构建 GMX 订单地址数组
     */
    function _buildAddresses(
        address market,
        address collateralToken,
        address receiver
    ) private pure returns (address[] memory) {
        address[] memory addresses = new address[](6);
        addresses[0] = receiver;           // receiver
        addresses[1] = address(0);         // callbackContract
        addresses[2] = address(0);         // uiFeeReceiver
        addresses[3] = market;             // market
        addresses[4] = collateralToken;    // initialCollateralToken
        addresses[5] = address(0);         // swapPath
        return addresses;
    }

    /**
     * @notice 构建 GMX 订单数字数组
     */
    function _buildNumbers(
        uint256 sizeInUsd,
        uint256 collateralAmount,
        uint256 acceptablePrice,
        uint256 executionFee
    ) private pure returns (uint256[] memory) {
        uint256[] memory numbers = new uint256[](7);
        numbers[0] = sizeInUsd;              // sizeDeltaUsd
        numbers[1] = collateralAmount;       // initialCollateralDeltaAmount
        numbers[2] = 0;                      // triggerPrice (0 for market orders)
        numbers[3] = acceptablePrice;        // acceptablePrice
        numbers[4] = executionFee;           // executionFee
        numbers[5] = 0;                      // callbackGasLimit
        numbers[6] = 0;                      // minOutputAmount
        return numbers;
    }

    /**
     * @notice 验证用户是否有指定仓位
     */
    function _verifyUserPosition(
        address user,
        address market,
        bool isLong
    ) private view {
        UserPosition[] memory positions = userPositions[user];
        bool found = false;

        for (uint256 i = 0; i < positions.length; i++) {
            if (positions[i].market == market && positions[i].isLong == isLong) {
                found = true;
                break;
            }
        }

        if (!found) revert NoPositionFound();
    }

    /**
     * @notice 获取用户在市场的主要方向
     */
    function _getUserMainDirection(address user, address market) private view returns (bool) {
        UserPosition[] memory positions = userPositions[user];
        uint256 longSize = 0;
        uint256 shortSize = 0;

        for (uint256 i = 0; i < positions.length; i++) {
            if (positions[i].market == market && !positions[i].isHedge) {
                if (positions[i].isLong) {
                    longSize += positions[i].sizeInUsd;
                } else {
                    shortSize += positions[i].sizeInUsd;
                }
            }
        }

        return longSize > shortSize; // true = mainly long, false = mainly short
    }

    /**
     * @notice 验证并计算杠杆
     * @dev 将抵押品金额转换为 USD 价值并计算杠杆
     */
    function _validateAndCalculateLeverage(
        address collateralToken,
        uint256 collateralAmount,
        uint256 sizeInUsd
    ) private view returns (uint256 leverage) {
        uint256 collateralDecimals = _getTokenDecimals(collateralToken);
        uint256 collateralInUsd = collateralAmount * (10 ** (18 - collateralDecimals));
        leverage = sizeInUsd / collateralInUsd;
        if (leverage > MAX_LEVERAGE) {
            revert InvalidLeverage(leverage);
        }
    }

    /**
     * @notice 创建对冲订单
     */
    function _createHedgeOrder(
        address user,
        address market,
        address collateralToken,
        uint256 hedgeSize,
        uint256 collateralAmount
    ) private returns (bytes32 orderKey) {
        bool isLong = !_getUserMainDirection(user, market);

        orderKey = exchangeRouter.createOrder{value: msg.value}(
            IGMXV2ExchangeRouter.CreateOrderParams({
                addresses: _buildAddresses(market, collateralToken, user),
                numbers: _buildNumbers(hedgeSize, collateralAmount, 0, msg.value),
                orderType: 0,
                isLong: isLong,
                shouldUnwrapNativeToken: false
            })
        );

        userPositions[user].push(UserPosition({
            market: market,
            collateralToken: collateralToken,
            isLong: isLong,
            sizeInUsd: hedgeSize,
            collateralAmount: collateralAmount,
            leverage: 10,
            openTimestamp: block.timestamp,
            isHedge: true
        }));

        orderToUser[orderKey] = user;
    }

    /**
     * @notice 计算所需抵押品 (10x 杠杆)
     * @dev 将 USD 价值转换为抵押品代币精度
     */
    function _calculateCollateral(
        address collateralToken,
        uint256 sizeInUsd
    ) private view returns (uint256) {
        uint256 decimals = _getTokenDecimals(collateralToken);
        return (sizeInUsd / 10) / (10 ** (18 - decimals));
    }

    /**
     * @notice 获取代币精度
     * @dev 尝试调用 decimals()，如果失败则假设为 18
     */
    function _getTokenDecimals(address token) private view returns (uint256) {
        try IERC20Metadata(token).decimals() returns (uint8 decimals) {
            return uint256(decimals);
        } catch {
            return 18; // 默认为 18 decimals
        }
    }

    /**
     * @notice 接收 ETH（用于执行费用）
     */
    receive() external payable {}
}
