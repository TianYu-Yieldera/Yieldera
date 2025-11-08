import { BaseListener } from '../BaseListener';
import { ethers } from 'ethers';

// CompoundV3Adapter ABI
const COMPOUND_ADAPTER_ABI = [
  'event Supplied(address indexed user, address indexed asset, uint256 amount, uint256 timestamp)',
  'event Withdrawn(address indexed user, address indexed asset, uint256 amount, uint256 timestamp)',
  'event SupplyRateUpdated(uint256 newRate, uint256 timestamp)',
  'event BorrowRateUpdated(uint256 newRate, uint256 timestamp)',
];

/**
 * CompoundListener - Compound V3 适配器事件监听
 *
 * 监控指标:
 * - 总存款/借款
 * - 供给率和借款率
 * - 利率变化检测
 * - 活跃用户
 */
export class CompoundListener extends BaseListener {
  private protocolStats = {
    totalSupplied: BigInt(0),
    totalWithdrawn: BigInt(0),
    activeUsers: new Set<string>(),
    currentSupplyRate: 0,
    currentBorrowRate: 0,
    lastRateUpdate: 0,
    lastUpdateTime: 0,
  };

  constructor(wsUrl: string, contractAddress: string) {
    super(wsUrl, contractAddress, COMPOUND_ADAPTER_ABI, 'CompoundListener');
  }

  /**
   * 注册Compound事件监听
   */
  protected async registerEventListeners(): Promise<void> {
    // 存款事件
    this.contract.on('Supplied', async (user, asset, amount, timestamp, event) => {
      await this.handleSupplied(user, asset, amount, timestamp, event);
    });

    // 取款事件
    this.contract.on('Withdrawn', async (user, asset, amount, timestamp, event) => {
      await this.handleWithdrawn(user, asset, amount, timestamp, event);
    });

    // 供给率更新
    this.contract.on('SupplyRateUpdated', async (newRate, timestamp, event) => {
      await this.handleSupplyRateUpdated(newRate, timestamp, event);
    });

    // 借款率更新
    this.contract.on('BorrowRateUpdated', async (newRate, timestamp, event) => {
      await this.handleBorrowRateUpdated(newRate, timestamp, event);
    });

    console.log(`[${this.listenerName}] Event listeners registered`);
  }

  /**
   * 处理存款事件
   */
  private async handleSupplied(
    user: string,
    asset: string,
    amount: bigint,
    timestamp: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'Supplied',
      user,
      asset,
      amount: amount.toString(),
      timestamp: timestamp.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    this.protocolStats.totalSupplied += amount;
    this.protocolStats.activeUsers.add(user);
    this.updateTimestamp();

    this.emit('supply', eventData);

    console.log(`[${this.listenerName}] Supply: ${ethers.formatEther(amount)} by ${user.substring(0, 10)}...`);
  }

  /**
   * 处理取款事件
   */
  private async handleWithdrawn(
    user: string,
    asset: string,
    amount: bigint,
    timestamp: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'Withdrawn',
      user,
      asset,
      amount: amount.toString(),
      timestamp: timestamp.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    this.protocolStats.totalWithdrawn += amount;
    this.updateTimestamp();

    this.emit('withdraw', eventData);

    // 大额取款告警
    const amountEth = Number(ethers.formatEther(amount));
    if (amountEth > 50) {
      this.emit('alert', {
        severity: 'warning',
        type: 'LARGE_WITHDRAWAL',
        message: `Large withdrawal from Compound: ${amountEth.toFixed(2)} ETH`,
        data: eventData,
      });
    }

    console.log(`[${this.listenerName}] Withdraw: ${ethers.formatEther(amount)} by ${user.substring(0, 10)}...`);
  }

  /**
   * 处理供给率更新
   */
  private async handleSupplyRateUpdated(
    newRate: bigint,
    timestamp: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'SupplyRateUpdated',
      newRate: newRate.toString(),
      timestamp: timestamp.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    const oldRate = this.protocolStats.currentSupplyRate;
    const newRatePercent = Number(newRate) / 100;
    this.protocolStats.currentSupplyRate = newRatePercent;
    this.protocolStats.lastRateUpdate = Date.now();

    this.emit('supplyRateUpdated', eventData);

    // 利率剧烈变化告警
    if (oldRate > 0) {
      const changePercent = Math.abs(newRatePercent - oldRate) / oldRate;
      if (changePercent > 0.5) { // 50%变化
        this.emit('alert', {
          severity: 'warning',
          type: 'SUPPLY_RATE_SPIKE',
          message: `Compound supply rate changed by ${(changePercent * 100).toFixed(1)}%: ${oldRate.toFixed(2)}% → ${newRatePercent.toFixed(2)}%`,
          data: eventData,
        });
      }
    }

    console.log(`[${this.listenerName}] Supply Rate Updated: ${newRatePercent.toFixed(2)}%`);
  }

  /**
   * 处理借款率更新
   */
  private async handleBorrowRateUpdated(
    newRate: bigint,
    timestamp: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'BorrowRateUpdated',
      newRate: newRate.toString(),
      timestamp: timestamp.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    const oldRate = this.protocolStats.currentBorrowRate;
    const newRatePercent = Number(newRate) / 100;
    this.protocolStats.currentBorrowRate = newRatePercent;

    this.emit('borrowRateUpdated', eventData);

    // 利率剧烈变化告警
    if (oldRate > 0) {
      const changePercent = Math.abs(newRatePercent - oldRate) / oldRate;
      if (changePercent > 0.5) { // 50%变化
        this.emit('alert', {
          severity: 'warning',
          type: 'BORROW_RATE_SPIKE',
          message: `Compound borrow rate changed by ${(changePercent * 100).toFixed(1)}%: ${oldRate.toFixed(2)}% → ${newRatePercent.toFixed(2)}%`,
          data: eventData,
        });
      }
    }

    console.log(`[${this.listenerName}] Borrow Rate Updated: ${newRatePercent.toFixed(2)}%`);
  }

  /**
   * 更新时间戳
   */
  private updateTimestamp(): void {
    this.protocolStats.lastUpdateTime = Date.now();
  }

  /**
   * 获取净存款
   */
  getNetSupply(): bigint {
    return this.protocolStats.totalSupplied - this.protocolStats.totalWithdrawn;
  }

  /**
   * 获取统计数据
   */
  getStats() {
    return {
      totalSupplied: ethers.formatEther(this.protocolStats.totalSupplied),
      totalWithdrawn: ethers.formatEther(this.protocolStats.totalWithdrawn),
      netSupply: ethers.formatEther(this.getNetSupply()),
      activeUsers: this.protocolStats.activeUsers.size,
      currentSupplyRate: `${this.protocolStats.currentSupplyRate.toFixed(2)}%`,
      currentBorrowRate: `${this.protocolStats.currentBorrowRate.toFixed(2)}%`,
      lastRateUpdate: this.protocolStats.lastRateUpdate > 0
        ? new Date(this.protocolStats.lastRateUpdate).toISOString()
        : 'N/A',
      lastUpdateTime: new Date(this.protocolStats.lastUpdateTime).toISOString(),
    };
  }

  /**
   * 重置统计
   */
  resetStats(): void {
    this.protocolStats = {
      totalSupplied: BigInt(0),
      totalWithdrawn: BigInt(0),
      activeUsers: new Set<string>(),
      currentSupplyRate: 0,
      currentBorrowRate: 0,
      lastRateUpdate: 0,
      lastUpdateTime: Date.now(),
    };
  }
}
