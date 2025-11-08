import { BaseListener } from '../BaseListener';
import { ethers } from 'ethers';

// TreasuryMarketplace ABI
const TREASURY_MARKETPLACE_ABI = [
  'event OrderCreated(uint256 indexed orderId, address indexed seller, uint256 indexed assetId, uint256 amount, uint256 pricePerToken, uint8 orderType)',
  'event OrderFilled(uint256 indexed orderId, address indexed buyer, address indexed seller, uint256 amount, uint256 totalPrice)',
  'event OrderCancelled(uint256 indexed orderId, address indexed seller)',
  'event OrderPartiallyFilled(uint256 indexed orderId, address indexed buyer, uint256 filledAmount, uint256 remainingAmount)',
];

/**
 * MarketplaceListener - Treasury市场事件监听
 *
 * 监控指标:
 * - 24小时交易量
 * - 订单创建/成交/取消
 * - 价格趋势
 * - 流动性深度
 */
export class MarketplaceListener extends BaseListener {
  private marketStats = {
    totalOrders: 0,
    totalVolume: BigInt(0),
    totalTrades: 0,
    ordersCancelled: 0,
    lastUpdateTime: 0,
  };

  constructor(wsUrl: string, contractAddress: string) {
    super(wsUrl, contractAddress, TREASURY_MARKETPLACE_ABI, 'MarketplaceListener');
  }

  /**
   * 注册Treasury Marketplace事件监听
   */
  protected async registerEventListeners(): Promise<void> {
    // 订单创建
    this.contract.on('OrderCreated', async (orderId, seller, assetId, amount, pricePerToken, orderType, event) => {
      await this.handleOrderCreated(orderId, seller, assetId, amount, pricePerToken, orderType, event);
    });

    // 订单成交
    this.contract.on('OrderFilled', async (orderId, buyer, seller, amount, totalPrice, event) => {
      await this.handleOrderFilled(orderId, buyer, seller, amount, totalPrice, event);
    });

    // 订单取消
    this.contract.on('OrderCancelled', async (orderId, seller, event) => {
      await this.handleOrderCancelled(orderId, seller, event);
    });

    // 部分成交
    this.contract.on('OrderPartiallyFilled', async (orderId, buyer, filledAmount, remainingAmount, event) => {
      await this.handleOrderPartiallyFilled(orderId, buyer, filledAmount, remainingAmount, event);
    });

    console.log(`[${this.listenerName}] Event listeners registered`);
  }

  /**
   * 处理订单创建
   */
  private async handleOrderCreated(
    orderId: bigint,
    seller: string,
    assetId: bigint,
    amount: bigint,
    pricePerToken: bigint,
    orderType: number,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'OrderCreated',
      orderId: orderId.toString(),
      seller,
      assetId: assetId.toString(),
      amount: amount.toString(),
      pricePerToken: pricePerToken.toString(),
      orderType,
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    this.marketStats.totalOrders++;
    this.updateTimestamp();

    this.emit('orderCreated', eventData);

    const orderTypeStr = orderType === 0 ? 'FIXED' : orderType === 1 ? 'AUCTION' : 'DUTCH';
    console.log(`[${this.listenerName}] Order Created #${orderId}: ${orderTypeStr} - ${amount} tokens @ ${ethers.formatUnits(pricePerToken, 6)} USDC`);
  }

  /**
   * 处理订单成交
   */
  private async handleOrderFilled(
    orderId: bigint,
    buyer: string,
    seller: string,
    amount: bigint,
    totalPrice: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'OrderFilled',
      orderId: orderId.toString(),
      buyer,
      seller,
      amount: amount.toString(),
      totalPrice: totalPrice.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    this.marketStats.totalVolume += totalPrice;
    this.marketStats.totalTrades++;
    this.updateTimestamp();

    this.emit('orderFilled', eventData);

    // 大额交易告警
    const volumeUSD = Number(ethers.formatUnits(totalPrice, 6));
    if (volumeUSD > 10000) {
      this.emit('alert', {
        severity: 'info',
        type: 'LARGE_TRADE',
        message: `Large trade: $${volumeUSD.toFixed(2)}`,
        data: eventData,
      });
    }

    console.log(`[${this.listenerName}] Order Filled #${orderId}: $${volumeUSD.toFixed(2)}`);
  }

  /**
   * 处理订单取消
   */
  private async handleOrderCancelled(
    orderId: bigint,
    seller: string,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'OrderCancelled',
      orderId: orderId.toString(),
      seller,
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    this.marketStats.ordersCancelled++;
    this.updateTimestamp();

    this.emit('orderCancelled', eventData);

    console.log(`[${this.listenerName}] Order Cancelled #${orderId}`);
  }

  /**
   * 处理部分成交
   */
  private async handleOrderPartiallyFilled(
    orderId: bigint,
    buyer: string,
    filledAmount: bigint,
    remainingAmount: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'OrderPartiallyFilled',
      orderId: orderId.toString(),
      buyer,
      filledAmount: filledAmount.toString(),
      remainingAmount: remainingAmount.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    this.emit('orderPartiallyFilled', eventData);

    console.log(`[${this.listenerName}] Order #${orderId} Partially Filled: ${filledAmount}/${BigInt(filledAmount) + BigInt(remainingAmount)}`);
  }

  /**
   * 更新时间戳
   */
  private updateTimestamp(): void {
    this.marketStats.lastUpdateTime = Date.now();
  }

  /**
   * 获取统计数据
   */
  getStats() {
    return {
      totalOrders: this.marketStats.totalOrders,
      totalVolume: `$${ethers.formatUnits(this.marketStats.totalVolume, 6)}`,
      totalTrades: this.marketStats.totalTrades,
      ordersCancelled: this.marketStats.ordersCancelled,
      cancellationRate: this.marketStats.totalOrders > 0
        ? `${((this.marketStats.ordersCancelled / this.marketStats.totalOrders) * 100).toFixed(2)}%`
        : '0%',
      lastUpdateTime: new Date(this.marketStats.lastUpdateTime).toISOString(),
    };
  }

  /**
   * 重置统计
   */
  resetStats(): void {
    this.marketStats = {
      totalOrders: 0,
      totalVolume: BigInt(0),
      totalTrades: 0,
      ordersCancelled: 0,
      lastUpdateTime: Date.now(),
    };
  }
}
