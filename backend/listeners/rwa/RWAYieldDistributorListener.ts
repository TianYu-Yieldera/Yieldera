import { BaseListener } from '../BaseListener';
import { ethers } from 'ethers';

// RWAYieldDistributor ABI
const RWA_YIELD_DISTRIBUTOR_ABI = [
  'event YieldDeposited(uint256 indexed distributionId, uint256 indexed assetId, address indexed paymentToken, uint256 amount, uint256 claimDeadline)',
  'event YieldClaimed(uint256 indexed distributionId, address indexed user, uint256 amount)',
  'event DistributionFinalized(uint256 indexed distributionId, uint256 totalClaimed, uint256 unclaimed)',
  'event UnclaimedYieldReclaimed(uint256 indexed distributionId, uint256 amount, address recipient)',
];

/**
 * RWAYieldDistributorListener - RWA资产收益分配监听
 *
 * 监控指标:
 * - 收益分配创建
 * - 用户领取行为
 * - 领取期限追踪
 * - 未领取收益回收
 * - 分配完成率
 */
export class RWAYieldDistributorListener extends BaseListener {
  private yieldStats = {
    totalDistributions: 0,
    totalYieldDeposited: BigInt(0),
    totalClaimed: BigInt(0),
    totalUnclaimed: BigInt(0),
    totalFinalized: 0,
    totalReclaimed: BigInt(0),
    activeDistributions: 0,
    lastUpdateTime: 0,
  };

  // 活跃分配追踪
  private activeDistributions = new Map<string, {
    distributionId: string;
    assetId: string;
    amount: bigint;
    claimDeadline: number;
    claimed: bigint;
    status: 'ACTIVE' | 'FINALIZED';
  }>();

  // 资产分配追踪
  private assetDistributions = new Map<string, {
    totalDistributions: number;
    totalYield: bigint;
    totalClaimed: bigint;
    lastDistributionTime: number;
  }>();

  // 支付代币统计
  private paymentTokenStats = new Map<string, {
    totalAmount: bigint;
    distributionCount: number;
  }>();

  constructor(wsUrl: string, contractAddress: string) {
    super(wsUrl, contractAddress, RWA_YIELD_DISTRIBUTOR_ABI, 'RWAYieldDistributor');
  }

  /**
   * 注册RWAYieldDistributor事件监听
   */
  protected async registerEventListeners(): Promise<void> {
    // 收益存入事件
    this.contract.on(
      'YieldDeposited',
      async (distributionId, assetId, paymentToken, amount, claimDeadline, event) => {
        await this.handleYieldDeposited(distributionId, assetId, paymentToken, amount, claimDeadline, event);
      }
    );

    // 用户领取事件
    this.contract.on('YieldClaimed', async (distributionId, user, amount, event) => {
      await this.handleYieldClaimed(distributionId, user, amount, event);
    });

    // 分配完成事件
    this.contract.on('DistributionFinalized', async (distributionId, totalClaimed, unclaimed, event) => {
      await this.handleDistributionFinalized(distributionId, totalClaimed, unclaimed, event);
    });

    // 未领取收益回收事件
    this.contract.on('UnclaimedYieldReclaimed', async (distributionId, amount, recipient, event) => {
      await this.handleUnclaimedYieldReclaimed(distributionId, amount, recipient, event);
    });

    console.log(`[${this.listenerName}] Event listeners registered`);
  }

  /**
   * 处理收益存入事件
   */
  private async handleYieldDeposited(
    distributionId: bigint,
    assetId: bigint,
    paymentToken: string,
    amount: bigint,
    claimDeadline: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'YieldDeposited',
      distributionId: distributionId.toString(),
      assetId: assetId.toString(),
      paymentToken,
      amount: amount.toString(),
      claimDeadline: claimDeadline.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
      timestamp: Date.now(),
    };

    // 更新统计
    this.yieldStats.totalDistributions++;
    this.yieldStats.totalYieldDeposited += amount;
    this.yieldStats.activeDistributions++;

    // 添加到活跃分配
    const distId = distributionId.toString();
    this.activeDistributions.set(distId, {
      distributionId: distId,
      assetId: assetId.toString(),
      amount,
      claimDeadline: Number(claimDeadline) * 1000, // 转换为毫秒
      claimed: BigInt(0),
      status: 'ACTIVE',
    });

    // 更新资产统计
    const assetIdStr = assetId.toString();
    const assetStats = this.assetDistributions.get(assetIdStr) || {
      totalDistributions: 0,
      totalYield: BigInt(0),
      totalClaimed: BigInt(0),
      lastDistributionTime: 0,
    };
    assetStats.totalDistributions++;
    assetStats.totalYield += amount;
    assetStats.lastDistributionTime = Date.now();
    this.assetDistributions.set(assetIdStr, assetStats);

    // 更新支付代币统计
    const tokenStats = this.paymentTokenStats.get(paymentToken) || {
      totalAmount: BigInt(0),
      distributionCount: 0,
    };
    tokenStats.totalAmount += amount;
    tokenStats.distributionCount++;
    this.paymentTokenStats.set(paymentToken, tokenStats);

    this.updateTimestamp();
    this.emit('yieldDeposited', eventData);

    const tokenName = this.getTokenName(paymentToken);
    const amountFormatted = this.formatAmount(amount, paymentToken);
    const deadlineDate = new Date(Number(claimDeadline) * 1000);
    console.log(
      `[${this.listenerName}] 💰 Yield Deposited #${distributionId}: ${amountFormatted} ${tokenName} for Asset #${assetId} (Deadline: ${deadlineDate.toISOString()})`
    );

    // 检查领取期限是否过短
    const claimPeriod = Number(claimDeadline) - Math.floor(Date.now() / 1000);
    if (claimPeriod < 7 * 24 * 60 * 60) {
      // < 7 days
      this.emit('alert', {
        level: 'WARNING',
        type: 'SHORT_CLAIM_PERIOD',
        message: `Short claim period detected: ${Math.floor(claimPeriod / 86400)} days for Distribution #${distributionId}`,
        data: eventData,
      });
    }

    // 大额分配告警
    if (amount > ethers.parseUnits('50000', 6)) {
      // > $50,000 USDC equivalent
      this.emit('alert', {
        level: 'INFO',
        type: 'LARGE_YIELD_DEPOSIT',
        message: `Large yield deposit: ${amountFormatted} ${tokenName} for Asset #${assetId}`,
        data: eventData,
      });
    }
  }

  /**
   * 处理用户领取事件
   */
  private async handleYieldClaimed(
    distributionId: bigint,
    user: string,
    amount: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'YieldClaimed',
      distributionId: distributionId.toString(),
      user,
      amount: amount.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
      timestamp: Date.now(),
    };

    // 更新统计
    this.yieldStats.totalClaimed += amount;

    // 更新活跃分配
    const distId = distributionId.toString();
    const dist = this.activeDistributions.get(distId);
    if (dist) {
      dist.claimed += amount;

      // 更新资产统计
      const assetStats = this.assetDistributions.get(dist.assetId);
      if (assetStats) {
        assetStats.totalClaimed += amount;
      }
    }

    this.updateTimestamp();
    this.emit('yieldClaimed', eventData);

    console.log(
      `[${this.listenerName}] 🎁 Yield Claimed: Distribution #${distributionId} - ${user.substring(0, 8)}... claimed ${amount.toString()}`
    );
  }

  /**
   * 处理分配完成事件
   */
  private async handleDistributionFinalized(
    distributionId: bigint,
    totalClaimed: bigint,
    unclaimed: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'DistributionFinalized',
      distributionId: distributionId.toString(),
      totalClaimed: totalClaimed.toString(),
      unclaimed: unclaimed.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
      timestamp: Date.now(),
    };

    // 更新统计
    this.yieldStats.totalFinalized++;
    this.yieldStats.totalUnclaimed += unclaimed;
    this.yieldStats.activeDistributions--;

    // 更新活跃分配状态
    const distId = distributionId.toString();
    const dist = this.activeDistributions.get(distId);
    if (dist) {
      dist.status = 'FINALIZED';
    }

    this.updateTimestamp();
    this.emit('distributionFinalized', eventData);

    const claimRate = totalClaimed > 0
      ? ((Number(totalClaimed) / (Number(totalClaimed) + Number(unclaimed))) * 100).toFixed(2)
      : '0.00';

    console.log(
      `[${this.listenerName}] ✅ Distribution Finalized #${distributionId}: Claimed: ${totalClaimed.toString()}, Unclaimed: ${unclaimed.toString()} (${claimRate}% claimed)`
    );

    // 低领取率告警
    if (parseFloat(claimRate) < 50 && totalClaimed + unclaimed > BigInt(0)) {
      this.emit('alert', {
        level: 'WARNING',
        type: 'LOW_CLAIM_RATE',
        message: `Low claim rate detected: ${claimRate}% for Distribution #${distributionId}`,
        data: eventData,
      });
    }
  }

  /**
   * 处理未领取收益回收事件
   */
  private async handleUnclaimedYieldReclaimed(
    distributionId: bigint,
    amount: bigint,
    recipient: string,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'UnclaimedYieldReclaimed',
      distributionId: distributionId.toString(),
      amount: amount.toString(),
      recipient,
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
      timestamp: Date.now(),
    };

    // 更新统计
    this.yieldStats.totalReclaimed += amount;
    this.updateTimestamp();

    this.emit('unclaimedYieldReclaimed', eventData);

    console.log(
      `[${this.listenerName}] 🔙 Unclaimed Yield Reclaimed #${distributionId}: ${amount.toString()} returned to ${recipient.substring(0, 8)}...`
    );
  }

  /**
   * 更新时间戳
   */
  private updateTimestamp(): void {
    this.yieldStats.lastUpdateTime = Date.now();
  }

  /**
   * 格式化金额
   */
  private formatAmount(amount: bigint, token: string): string {
    const decimals = token === ethers.ZeroAddress ? 18 : 6; // ETH: 18, USDC: 6
    return ethers.formatUnits(amount, decimals);
  }

  /**
   * 获取代币名称
   */
  private getTokenName(token: string): string {
    if (token === ethers.ZeroAddress) {
      return 'ETH';
    }
    // 简化版，可以扩展为合约调用
    return 'USDC';
  }

  /**
   * 获取统计数据
   */
  public getStats() {
    return {
      ...this.yieldStats,
      totalYieldDeposited: this.yieldStats.totalYieldDeposited.toString(),
      totalClaimed: this.yieldStats.totalClaimed.toString(),
      totalUnclaimed: this.yieldStats.totalUnclaimed.toString(),
      totalReclaimed: this.yieldStats.totalReclaimed.toString(),
      paymentTokens: Array.from(this.paymentTokenStats.entries()).map(([token, stats]) => ({
        token,
        tokenName: this.getTokenName(token),
        totalAmount: stats.totalAmount.toString(),
        distributionCount: stats.distributionCount,
      })),
      assetDistributions: Array.from(this.assetDistributions.entries()).map(([assetId, stats]) => ({
        assetId,
        totalDistributions: stats.totalDistributions,
        totalYield: stats.totalYield.toString(),
        totalClaimed: stats.totalClaimed.toString(),
        lastDistributionTime: stats.lastDistributionTime,
        claimRate:
          stats.totalYield > 0
            ? ((Number(stats.totalClaimed) / Number(stats.totalYield)) * 100).toFixed(2) + '%'
            : '0%',
      })),
      overallClaimRate:
        this.yieldStats.totalYieldDeposited > 0
          ? (
              (Number(this.yieldStats.totalClaimed) / Number(this.yieldStats.totalYieldDeposited)) *
              100
            ).toFixed(2) + '%'
          : '0%',
      avgYieldPerDistribution:
        this.yieldStats.totalDistributions > 0
          ? (this.yieldStats.totalYieldDeposited / BigInt(this.yieldStats.totalDistributions)).toString()
          : '0',
    };
  }

  /**
   * 获取即将到期的分配
   */
  public getExpiringDistributions(withinHours: number = 24): Array<any> {
    const now = Date.now();
    const threshold = now + withinHours * 60 * 60 * 1000;

    return Array.from(this.activeDistributions.values())
      .filter((dist) => dist.status === 'ACTIVE' && dist.claimDeadline <= threshold && dist.claimDeadline > now)
      .map((dist) => ({
        ...dist,
        amount: dist.amount.toString(),
        claimed: dist.claimed.toString(),
        hoursRemaining: ((dist.claimDeadline - now) / (60 * 60 * 1000)).toFixed(2),
      }));
  }

  /**
   * 获取资产分配历史
   */
  public getAssetDistributionHistory(assetId: string) {
    return this.assetDistributions.get(assetId) || null;
  }

  /**
   * 重置统计数据
   */
  public resetStats(): void {
    this.yieldStats = {
      totalDistributions: 0,
      totalYieldDeposited: BigInt(0),
      totalClaimed: BigInt(0),
      totalUnclaimed: BigInt(0),
      totalFinalized: 0,
      totalReclaimed: BigInt(0),
      activeDistributions: 0,
      lastUpdateTime: 0,
    };
    this.activeDistributions.clear();
    this.assetDistributions.clear();
    this.paymentTokenStats.clear();
  }
}
