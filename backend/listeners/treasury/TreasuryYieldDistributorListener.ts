import { BaseListener } from '../BaseListener';
import { ethers } from 'ethers';

// TreasuryYieldDistributor ABI
const TREASURY_YIELD_DISTRIBUTOR_ABI = [
  'event YieldDeposited(uint256 indexed distributionId, uint256 indexed assetId, uint256 totalYield, uint256 yieldPerToken, string distributionType)',
  'event YieldClaimed(address indexed user, uint256 indexed assetId, uint256 amount, uint256 distributionId)',
  'event BatchDistributed(uint256 indexed distributionId, uint256 indexed assetId, uint256 recipientsCount, uint256 totalAmount)',
];

/**
 * TreasuryYieldDistributorListener - Treasury资产收益分配监听
 *
 * 监控指标:
 * - 票息支付 (COUPON)
 * - 到期赎回 (MATURITY)
 * - 用户领取行为
 * - 批量分配效率
 * - 累计分配金额
 */
export class TreasuryYieldDistributorListener extends BaseListener {
  private yieldStats = {
    totalDistributions: 0,
    totalYieldDistributed: BigInt(0),
    totalClaims: 0,
    totalClaimAmount: BigInt(0),
    batchDistributions: 0,
    couponPayments: 0,
    maturityPayments: 0,
    lastUpdateTime: 0,
  };

  // 分配类型统计
  private distributionsByType = new Map<string, {
    count: number;
    totalAmount: bigint;
  }>();

  // 资产分配追踪
  private assetYields = new Map<string, {
    distributions: number;
    totalYield: bigint;
    lastDistributionTime: number;
  }>();

  constructor(wsUrl: string, contractAddress: string) {
    super(wsUrl, contractAddress, TREASURY_YIELD_DISTRIBUTOR_ABI, 'TreasuryYieldDistributor');
  }

  /**
   * 注册TreasuryYieldDistributor事件监听
   */
  protected async registerEventListeners(): Promise<void> {
    // 收益存入事件
    this.contract.on(
      'YieldDeposited',
      async (distributionId, assetId, totalYield, yieldPerToken, distributionType, event) => {
        await this.handleYieldDeposited(
          distributionId,
          assetId,
          totalYield,
          yieldPerToken,
          distributionType,
          event
        );
      }
    );

    // 用户领取事件
    this.contract.on('YieldClaimed', async (user, assetId, amount, distributionId, event) => {
      await this.handleYieldClaimed(user, assetId, amount, distributionId, event);
    });

    // 批量分配事件
    this.contract.on(
      'BatchDistributed',
      async (distributionId, assetId, recipientsCount, totalAmount, event) => {
        await this.handleBatchDistributed(distributionId, assetId, recipientsCount, totalAmount, event);
      }
    );

    console.log(`[${this.listenerName}] Event listeners registered`);
  }

  /**
   * 处理收益存入事件
   */
  private async handleYieldDeposited(
    distributionId: bigint,
    assetId: bigint,
    totalYield: bigint,
    yieldPerToken: bigint,
    distributionType: string,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'YieldDeposited',
      distributionId: distributionId.toString(),
      assetId: assetId.toString(),
      totalYield: totalYield.toString(),
      yieldPerToken: yieldPerToken.toString(),
      distributionType,
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
      timestamp: Date.now(),
    };

    // 更新统计
    this.yieldStats.totalDistributions++;
    this.yieldStats.totalYieldDistributed += totalYield;

    // 更新分配类型统计
    if (distributionType === 'COUPON') {
      this.yieldStats.couponPayments++;
    } else if (distributionType === 'MATURITY') {
      this.yieldStats.maturityPayments++;
    }

    const typeStats = this.distributionsByType.get(distributionType) || {
      count: 0,
      totalAmount: BigInt(0),
    };
    typeStats.count++;
    typeStats.totalAmount += totalYield;
    this.distributionsByType.set(distributionType, typeStats);

    // 更新资产收益追踪
    const assetIdStr = assetId.toString();
    const assetYield = this.assetYields.get(assetIdStr) || {
      distributions: 0,
      totalYield: BigInt(0),
      lastDistributionTime: 0,
    };
    assetYield.distributions++;
    assetYield.totalYield += totalYield;
    assetYield.lastDistributionTime = Date.now();
    this.assetYields.set(assetIdStr, assetYield);

    this.updateTimestamp();
    this.emit('yieldDeposited', eventData);

    const yieldUSD = ethers.formatUnits(totalYield, 6);
    const yieldPerTokenFormatted = ethers.formatUnits(yieldPerToken, 18);
    console.log(
      `[${this.listenerName}] 💰 Yield Deposited #${distributionId}: ${distributionType} - $${yieldUSD} (${yieldPerTokenFormatted} per token) for Asset #${assetId}`
    );

    // 大额分配告警
    if (totalYield > ethers.parseUnits('100000', 6)) {
      // > $100,000
      this.emit('alert', {
        level: 'WARNING',
        type: 'LARGE_YIELD_DEPOSIT',
        message: `Large yield deposit detected: $${yieldUSD} for Asset #${assetId}`,
        data: eventData,
      });
    }
  }

  /**
   * 处理用户领取事件
   */
  private async handleYieldClaimed(
    user: string,
    assetId: bigint,
    amount: bigint,
    distributionId: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'YieldClaimed',
      user,
      assetId: assetId.toString(),
      amount: amount.toString(),
      distributionId: distributionId.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
      timestamp: Date.now(),
    };

    // 更新统计
    this.yieldStats.totalClaims++;
    this.yieldStats.totalClaimAmount += amount;
    this.updateTimestamp();

    this.emit('yieldClaimed', eventData);

    const amountUSD = ethers.formatUnits(amount, 6);
    console.log(
      `[${this.listenerName}] 🎁 Yield Claimed: $${amountUSD} by ${user.substring(0, 8)}... (Distribution #${distributionId})`
    );
  }

  /**
   * 处理批量分配事件
   */
  private async handleBatchDistributed(
    distributionId: bigint,
    assetId: bigint,
    recipientsCount: bigint,
    totalAmount: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'BatchDistributed',
      distributionId: distributionId.toString(),
      assetId: assetId.toString(),
      recipientsCount: recipientsCount.toString(),
      totalAmount: totalAmount.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
      timestamp: Date.now(),
    };

    // 更新统计
    this.yieldStats.batchDistributions++;
    this.updateTimestamp();

    this.emit('batchDistributed', eventData);

    const amountUSD = ethers.formatUnits(totalAmount, 6);
    console.log(
      `[${this.listenerName}] 📦 Batch Distributed #${distributionId}: $${amountUSD} to ${recipientsCount} recipients`
    );

    // 大规模批量分配告警
    if (recipientsCount > BigInt(100)) {
      this.emit('alert', {
        level: 'INFO',
        type: 'LARGE_BATCH_DISTRIBUTION',
        message: `Large batch distribution: ${recipientsCount} recipients received $${amountUSD}`,
        data: eventData,
      });
    }
  }

  /**
   * 更新时间戳
   */
  private updateTimestamp(): void {
    this.yieldStats.lastUpdateTime = Date.now();
  }

  /**
   * 获取统计数据
   */
  public getStats() {
    return {
      ...this.yieldStats,
      totalYieldDistributed: this.yieldStats.totalYieldDistributed.toString(),
      totalClaimAmount: this.yieldStats.totalClaimAmount.toString(),
      distributionsByType: Array.from(this.distributionsByType.entries()).map(([type, stats]) => ({
        type,
        count: stats.count,
        totalAmount: stats.totalAmount.toString(),
      })),
      assetYields: Array.from(this.assetYields.entries()).map(([assetId, stats]) => ({
        assetId,
        distributions: stats.distributions,
        totalYield: stats.totalYield.toString(),
        lastDistributionTime: stats.lastDistributionTime,
      })),
      avgYieldPerDistribution:
        this.yieldStats.totalDistributions > 0
          ? (this.yieldStats.totalYieldDistributed / BigInt(this.yieldStats.totalDistributions)).toString()
          : '0',
      avgClaimAmount:
        this.yieldStats.totalClaims > 0
          ? (this.yieldStats.totalClaimAmount / BigInt(this.yieldStats.totalClaims)).toString()
          : '0',
      claimRate:
        this.yieldStats.totalDistributions > 0
          ? ((this.yieldStats.totalClaims / this.yieldStats.totalDistributions) * 100).toFixed(2) + '%'
          : '0%',
    };
  }

  /**
   * 获取资产收益历史
   */
  public getAssetYieldHistory(assetId: string) {
    return this.assetYields.get(assetId) || null;
  }

  /**
   * 重置统计数据
   */
  public resetStats(): void {
    this.yieldStats = {
      totalDistributions: 0,
      totalYieldDistributed: BigInt(0),
      totalClaims: 0,
      totalClaimAmount: BigInt(0),
      batchDistributions: 0,
      couponPayments: 0,
      maturityPayments: 0,
      lastUpdateTime: 0,
    };
    this.distributionsByType.clear();
    this.assetYields.clear();
  }
}
