import { BaseListener } from '../BaseListener';
import { ethers } from 'ethers';

// AaveV3Adapter ABI
const AAVE_ADAPTER_ABI = [
  'event Supplied(address indexed user, address indexed asset, uint256 amount, uint256 timestamp)',
  'event Withdrawn(address indexed user, address indexed asset, uint256 amount, uint256 timestamp)',
  'event Borrowed(address indexed user, address indexed asset, uint256 amount, uint256 interestRateMode, uint256 timestamp)',
  'event Repaid(address indexed user, address indexed asset, uint256 amount, uint256 interestRateMode, uint256 timestamp)',
  'event FlashLoanExecuted(address indexed initiator, address indexed asset, uint256 amount, uint256 premium, uint256 timestamp)',
];

/**
 * AaveListener - Aave V3 适配器事件监听
 *
 * 监控指标:
 * - 总存款/借款
 * - 利用率
 * - 活跃用户
 * - 闪电贷监控
 * - 健康因子预警
 */
export class AaveListener extends BaseListener {
  private protocolStats = {
    totalSupplied: BigInt(0),
    totalBorrowed: BigInt(0),
    activeUsers: new Set<string>(),
    flashLoanCount: 0,
    lastUpdateTime: 0,
  };

  constructor(wsUrl: string, contractAddress: string) {
    super(wsUrl, contractAddress, AAVE_ADAPTER_ABI, 'AaveListener');
  }

  /**
   * 注册Aave事件监听
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

    // 借款事件
    this.contract.on('Borrowed', async (user, asset, amount, interestRateMode, timestamp, event) => {
      await this.handleBorrowed(user, asset, amount, interestRateMode, timestamp, event);
    });

    // 还款事件
    this.contract.on('Repaid', async (user, asset, amount, interestRateMode, timestamp, event) => {
      await this.handleRepaid(user, asset, amount, interestRateMode, timestamp, event);
    });

    // 闪电贷事件
    this.contract.on('FlashLoanExecuted', async (initiator, asset, amount, premium, timestamp, event) => {
      await this.handleFlashLoan(initiator, asset, amount, premium, timestamp, event);
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
    this.checkUtilization();

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

    this.protocolStats.totalSupplied -= amount;
    this.updateTimestamp();

    this.emit('withdraw', eventData);
    this.checkUtilization();

    // 大额取款告警
    const amountEth = Number(ethers.formatEther(amount));
    if (amountEth > 50) {
      this.emit('alert', {
        severity: 'warning',
        type: 'LARGE_WITHDRAWAL',
        message: `Large withdrawal: ${amountEth.toFixed(2)} ETH`,
        data: eventData,
      });
    }

    console.log(`[${this.listenerName}] Withdraw: ${ethers.formatEther(amount)} by ${user.substring(0, 10)}...`);
  }

  /**
   * 处理借款事件
   */
  private async handleBorrowed(
    user: string,
    asset: string,
    amount: bigint,
    interestRateMode: bigint,
    timestamp: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'Borrowed',
      user,
      asset,
      amount: amount.toString(),
      interestRateMode: interestRateMode.toString(),
      timestamp: timestamp.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    this.protocolStats.totalBorrowed += amount;
    this.protocolStats.activeUsers.add(user);
    this.updateTimestamp();

    this.emit('borrow', eventData);
    this.checkUtilization();

    console.log(`[${this.listenerName}] Borrow: ${ethers.formatEther(amount)} by ${user.substring(0, 10)}...`);
  }

  /**
   * 处理还款事件
   */
  private async handleRepaid(
    user: string,
    asset: string,
    amount: bigint,
    interestRateMode: bigint,
    timestamp: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'Repaid',
      user,
      asset,
      amount: amount.toString(),
      interestRateMode: interestRateMode.toString(),
      timestamp: timestamp.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    this.protocolStats.totalBorrowed -= amount;
    this.updateTimestamp();

    this.emit('repay', eventData);
    this.checkUtilization();

    console.log(`[${this.listenerName}] Repay: ${ethers.formatEther(amount)} by ${user.substring(0, 10)}...`);
  }

  /**
   * 处理闪电贷事件
   */
  private async handleFlashLoan(
    initiator: string,
    asset: string,
    amount: bigint,
    premium: bigint,
    timestamp: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'FlashLoanExecuted',
      initiator,
      asset,
      amount: amount.toString(),
      premium: premium.toString(),
      timestamp: timestamp.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    this.protocolStats.flashLoanCount++;

    this.emit('flashLoan', eventData);

    // 闪电贷告警（潜在攻击）
    const amountEth = Number(ethers.formatEther(amount));
    if (amountEth > 100) {
      this.emit('alert', {
        severity: 'critical',
        type: 'LARGE_FLASH_LOAN',
        message: `Large flash loan: ${amountEth.toFixed(2)} ETH`,
        data: eventData,
      });
    }

    console.log(`[${this.listenerName}] Flash Loan: ${ethers.formatEther(amount)} by ${initiator.substring(0, 10)}...`);
  }

  /**
   * 检查利用率
   */
  private checkUtilization(): void {
    const utilizationRate = this.getUtilizationRate();

    if (utilizationRate > 0.9) {
      this.emit('alert', {
        severity: 'critical',
        type: 'HIGH_UTILIZATION',
        message: `Critical utilization rate: ${(utilizationRate * 100).toFixed(2)}%`,
        data: this.getStats(),
      });
    } else if (utilizationRate > 0.8) {
      this.emit('alert', {
        severity: 'warning',
        type: 'HIGH_UTILIZATION',
        message: `High utilization rate: ${(utilizationRate * 100).toFixed(2)}%`,
        data: this.getStats(),
      });
    }
  }

  /**
   * 计算利用率
   */
  private getUtilizationRate(): number {
    if (this.protocolStats.totalSupplied === BigInt(0)) {
      return 0;
    }
    return Number(this.protocolStats.totalBorrowed) / Number(this.protocolStats.totalSupplied);
  }

  /**
   * 更新时间戳
   */
  private updateTimestamp(): void {
    this.protocolStats.lastUpdateTime = Date.now();
  }

  /**
   * 获取统计数据
   */
  getStats() {
    const utilizationRate = this.getUtilizationRate();

    return {
      totalSupplied: ethers.formatEther(this.protocolStats.totalSupplied),
      totalBorrowed: ethers.formatEther(this.protocolStats.totalBorrowed),
      utilizationRate: `${(utilizationRate * 100).toFixed(2)}%`,
      activeUsers: this.protocolStats.activeUsers.size,
      flashLoanCount: this.protocolStats.flashLoanCount,
      lastUpdateTime: new Date(this.protocolStats.lastUpdateTime).toISOString(),
    };
  }

  /**
   * 重置统计
   */
  resetStats(): void {
    this.protocolStats = {
      totalSupplied: BigInt(0),
      totalBorrowed: BigInt(0),
      activeUsers: new Set<string>(),
      flashLoanCount: 0,
      lastUpdateTime: Date.now(),
    };
  }
}
