import { BaseListener } from '../BaseListener';
import { ethers } from 'ethers';

// TreasuryAssetFactory ABI
const ASSET_FACTORY_ABI = [
  'event AssetCreated(uint256 indexed assetId, string symbol, string cusip, uint256 totalValue, uint256 maturityDate, uint256 couponRate)',
  'event AssetVerified(uint256 indexed assetId, address indexed verifier, uint256 timestamp)',
  'event AssetStatusUpdated(uint256 indexed assetId, uint8 newStatus)',
  'event AssetMatured(uint256 indexed assetId, uint256 finalValue)',
];

/**
 * AssetFactoryListener - Treasury资产工厂事件监听
 *
 * 监控指标:
 * - 资产创建频率
 * - 验证状态跟踪
 * - 总资产价值
 * - 到期资产监控
 */
export class AssetFactoryListener extends BaseListener {
  private factoryStats = {
    totalAssetsCreated: 0,
    verifiedAssets: 0,
    totalValue: BigInt(0),
    maturedAssets: 0,
    lastAssetId: 0,
    lastUpdateTime: 0,
  };

  // 资产状态映射
  private assetStatus: Map<number, {
    verified: boolean,
    status: number,
    value: bigint,
  }> = new Map();

  constructor(wsUrl: string, contractAddress: string) {
    super(wsUrl, contractAddress, ASSET_FACTORY_ABI, 'AssetFactoryListener');
  }

  /**
   * 注册AssetFactory事件监听
   */
  protected async registerEventListeners(): Promise<void> {
    // 资产创建
    this.contract.on('AssetCreated', async (assetId, symbol, cusip, totalValue, maturityDate, couponRate, event) => {
      await this.handleAssetCreated(assetId, symbol, cusip, totalValue, maturityDate, couponRate, event);
    });

    // 资产验证
    this.contract.on('AssetVerified', async (assetId, verifier, timestamp, event) => {
      await this.handleAssetVerified(assetId, verifier, timestamp, event);
    });

    // 状态更新
    this.contract.on('AssetStatusUpdated', async (assetId, newStatus, event) => {
      await this.handleStatusUpdated(assetId, newStatus, event);
    });

    // 资产到期
    this.contract.on('AssetMatured', async (assetId, finalValue, event) => {
      await this.handleAssetMatured(assetId, finalValue, event);
    });

    console.log(`[${this.listenerName}] Event listeners registered`);
  }

  /**
   * 处理资产创建
   */
  private async handleAssetCreated(
    assetId: bigint,
    symbol: string,
    cusip: string,
    totalValue: bigint,
    maturityDate: bigint,
    couponRate: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'AssetCreated',
      assetId: assetId.toString(),
      symbol,
      cusip,
      totalValue: totalValue.toString(),
      maturityDate: maturityDate.toString(),
      couponRate: couponRate.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    this.factoryStats.totalAssetsCreated++;
    this.factoryStats.totalValue += totalValue;
    this.factoryStats.lastAssetId = Number(assetId);
    this.updateTimestamp();

    // 初始化资产状态
    this.assetStatus.set(Number(assetId), {
      verified: false,
      status: 0,
      value: totalValue,
    });

    this.emit('assetCreated', eventData);

    const valueUSD = Number(ethers.formatUnits(totalValue, 6));
    console.log(`[${this.listenerName}] Asset Created #${assetId}: ${symbol} (${cusip}) - $${valueUSD.toLocaleString()}`);

    // 大额资产告警
    if (valueUSD > 1000000) { // $1M
      this.emit('alert', {
        severity: 'info',
        type: 'LARGE_ASSET_CREATED',
        message: `Large asset created: ${symbol} worth $${valueUSD.toLocaleString()}`,
        data: eventData,
      });
    }
  }

  /**
   * 处理资产验证
   */
  private async handleAssetVerified(
    assetId: bigint,
    verifier: string,
    timestamp: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'AssetVerified',
      assetId: assetId.toString(),
      verifier,
      timestamp: timestamp.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    this.factoryStats.verifiedAssets++;
    this.updateTimestamp();

    // 更新资产状态
    const asset = this.assetStatus.get(Number(assetId));
    if (asset) {
      asset.verified = true;
    }

    this.emit('assetVerified', eventData);

    console.log(`[${this.listenerName}] Asset Verified #${assetId} by ${verifier.substring(0, 10)}...`);
  }

  /**
   * 处理状态更新
   */
  private async handleStatusUpdated(
    assetId: bigint,
    newStatus: number,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'AssetStatusUpdated',
      assetId: assetId.toString(),
      newStatus,
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    // 更新资产状态
    const asset = this.assetStatus.get(Number(assetId));
    if (asset) {
      asset.status = newStatus;
    }

    this.emit('assetStatusUpdated', eventData);

    const statusName = this.getStatusName(newStatus);
    console.log(`[${this.listenerName}] Asset #${assetId} status: ${statusName}`);

    // 暂停或违约状态告警
    if (newStatus === 2 || newStatus === 3) { // PAUSED or DEFAULTED
      this.emit('alert', {
        severity: 'critical',
        type: 'ASSET_ISSUE',
        message: `Asset #${assetId} status changed to: ${statusName}`,
        data: eventData,
      });
    }
  }

  /**
   * 处理资产到期
   */
  private async handleAssetMatured(
    assetId: bigint,
    finalValue: bigint,
    event: ethers.Log
  ): Promise<void> {
    const eventData = {
      eventType: 'AssetMatured',
      assetId: assetId.toString(),
      finalValue: finalValue.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
    };

    this.factoryStats.maturedAssets++;
    this.updateTimestamp();

    this.emit('assetMatured', eventData);

    const valueUSD = Number(ethers.formatUnits(finalValue, 6));
    console.log(`[${this.listenerName}] Asset Matured #${assetId}: Final Value $${valueUSD.toLocaleString()}`);
  }

  /**
   * 获取状态名称
   */
  private getStatusName(status: number): string {
    const statuses = ['ACTIVE', 'VERIFIED', 'PAUSED', 'DEFAULTED', 'MATURED'];
    return statuses[status] || 'UNKNOWN';
  }

  /**
   * 获取未验证资产数量
   */
  getUnverifiedCount(): number {
    let count = 0;
    this.assetStatus.forEach(asset => {
      if (!asset.verified) count++;
    });
    return count;
  }

  /**
   * 更新时间戳
   */
  private updateTimestamp(): void {
    this.factoryStats.lastUpdateTime = Date.now();
  }

  /**
   * 获取统计数据
   */
  getStats() {
    return {
      totalAssetsCreated: this.factoryStats.totalAssetsCreated,
      verifiedAssets: this.factoryStats.verifiedAssets,
      unverifiedAssets: this.getUnverifiedCount(),
      totalValue: `$${ethers.formatUnits(this.factoryStats.totalValue, 6)}`,
      maturedAssets: this.factoryStats.maturedAssets,
      lastAssetId: this.factoryStats.lastAssetId,
      verificationRate: this.factoryStats.totalAssetsCreated > 0
        ? `${((this.factoryStats.verifiedAssets / this.factoryStats.totalAssetsCreated) * 100).toFixed(1)}%`
        : '0%',
      lastUpdateTime: new Date(this.factoryStats.lastUpdateTime).toISOString(),
    };
  }

  /**
   * 重置统计
   */
  resetStats(): void {
    this.factoryStats = {
      totalAssetsCreated: 0,
      verifiedAssets: 0,
      totalValue: BigInt(0),
      maturedAssets: 0,
      lastAssetId: 0,
      lastUpdateTime: Date.now(),
    };
    this.assetStatus.clear();
  }
}
