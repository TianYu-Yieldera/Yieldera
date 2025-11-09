import { BaseListener } from '../BaseListener';
import { ethers } from 'ethers';

// UniswapV3Adapter ABI (只包含需要的事件)
const UNISWAP_ADAPTER_ABI = [
  'event Swapped(address indexed user, address indexed tokenIn, address indexed tokenOut, uint256 amountIn, uint256 amountOut, uint24 fee)',
  'event MultiHopSwap(address indexed user, uint256 amountIn, uint256 amountOut)',
];

/**
 * UniswapListener - Uniswap V3 适配器事件监听
 *
 * 监控指标:
 * - Swap交易量和频率
 * - 滑点计算
 * - Gas消耗
 * - 异常交易检测
 */
export class UniswapListener extends BaseListener {
  private swapStats = {
    totalSwaps: 0,
    totalVolume: BigInt(0),
    totalSlippage: 0,
    lastUpdateTime: 0,
  };

  constructor(wsUrl: string, contractAddress: string) {
    super(wsUrl, contractAddress, UNISWAP_ADAPTER_ABI, 'UniswapListener');
  }

  /**
   * 注册Uniswap事件监听
   */
  protected async registerEventListeners(): Promise<void> {
    // 监听 Swapped 事件
    this.contract.on(
      'Swapped',
      async (user, tokenIn, tokenOut, amountIn, amountOut, fee, event) => {
        await this.handleSwapped(user, tokenIn, tokenOut, amountIn, amountOut, fee, event);
      }
    );

    // 监听 MultiHopSwap 事件
    this.contract.on(
      'MultiHopSwap',
      async (user, amountIn, amountOut, event) => {
        await this.handleMultiHopSwap(user, amountIn, amountOut, event);
      }
    );

    console.log(`[${this.listenerName}] Event listeners registered`);
  }

  /**
   * 处理Swapped事件
   */
  private async handleSwapped(
    user: string,
    tokenIn: string,
    tokenOut: string,
    amountIn: bigint,
    amountOut: bigint,
    fee: number,
    event: ethers.Log
  ): Promise<void> {
    try {
      const block = await event.getBlock();
      const transaction = await event.getTransaction();

      // 计算滑点 (简化版，实际需要价格预言机)
      const slippage = this.calculateSlippage(amountIn, amountOut);

      const swapData = {
        eventType: 'Swapped',
        user,
        tokenIn,
        tokenOut,
        amountIn: amountIn.toString(),
        amountOut: amountOut.toString(),
        fee: fee / 10000, // 转换为百分比
        slippage,
        blockNumber: event.blockNumber,
        transactionHash: event.transactionHash,
        timestamp: block.timestamp,
        gasUsed: transaction?.gasLimit.toString(),
      };

      // 更新统计
      this.updateStats(amountIn, slippage);

      // 检查异常
      this.checkAnomalies(swapData);

      // 发出事件
      this.emit('swap', swapData);

      console.log(`[${this.listenerName}] Swap detected:`, {
        user: user.substring(0, 10) + '...',
        amountIn: ethers.formatEther(amountIn),
        amountOut: ethers.formatEther(amountOut),
        slippage: `${(slippage * 100).toFixed(2)}%`,
      });
    } catch (error) {
      console.error(`[${this.listenerName}] Error handling Swapped event:`, error);
      this.emit('error', error);
    }
  }

  /**
   * 处理MultiHopSwap事件
   */
  private async handleMultiHopSwap(
    user: string,
    amountIn: bigint,
    amountOut: bigint,
    event: ethers.Log
  ): Promise<void> {
    try {
      const block = await event.getBlock();

      const swapData = {
        eventType: 'MultiHopSwap',
        user,
        amountIn: amountIn.toString(),
        amountOut: amountOut.toString(),
        blockNumber: event.blockNumber,
        transactionHash: event.transactionHash,
        timestamp: block.timestamp,
      };

      this.emit('multiHopSwap', swapData);

      console.log(`[${this.listenerName}] Multi-hop swap detected:`, {
        user: user.substring(0, 10) + '...',
        amountIn: ethers.formatEther(amountIn),
        amountOut: ethers.formatEther(amountOut),
      });
    } catch (error) {
      console.error(`[${this.listenerName}] Error handling MultiHopSwap event:`, error);
    }
  }

  /**
   * 计算滑点 (简化版)
   */
  private calculateSlippage(amountIn: bigint, amountOut: bigint): number {
    // 实际应该基于预言机价格，这里简化为固定比例
    // slippage = |expectedAmount - actualAmount| / expectedAmount
    const expectedRatio = 1.0; // 假设1:1
    const actualRatio = Number(amountOut) / Number(amountIn);
    return Math.abs(expectedRatio - actualRatio) / expectedRatio;
  }

  /**
   * 更新统计数据
   */
  private updateStats(amountIn: bigint, slippage: number): void {
    this.swapStats.totalSwaps++;
    this.swapStats.totalVolume += amountIn;
    this.swapStats.totalSlippage += slippage;
    this.swapStats.lastUpdateTime = Date.now();
  }

  /**
   * 检查异常交易
   */
  private checkAnomalies(swapData: any): void {
    // 高滑点告警
    if (swapData.slippage > 0.02) { // 2%
      this.emit('alert', {
        severity: 'warning',
        type: 'HIGH_SLIPPAGE',
        message: `High slippage detected: ${(swapData.slippage * 100).toFixed(2)}%`,
        data: swapData,
      });
    }

    // 大额交易告警
    const amountInEth = Number(ethers.formatEther(swapData.amountIn));
    if (amountInEth > 100) { // 假设100 ETH为大额
      this.emit('alert', {
        severity: 'info',
        type: 'LARGE_SWAP',
        message: `Large swap detected: ${amountInEth.toFixed(2)} ETH`,
        data: swapData,
      });
    }
  }

  /**
   * 获取统计数据
   */
  getStats() {
    const avgSlippage = this.swapStats.totalSwaps > 0
      ? this.swapStats.totalSlippage / this.swapStats.totalSwaps
      : 0;

    return {
      totalSwaps: this.swapStats.totalSwaps,
      totalVolume: ethers.formatEther(this.swapStats.totalVolume),
      avgSlippage: `${(avgSlippage * 100).toFixed(4)}%`,
      lastUpdateTime: new Date(this.swapStats.lastUpdateTime).toISOString(),
    };
  }

  /**
   * 重置统计
   */
  resetStats(): void {
    this.swapStats = {
      totalSwaps: 0,
      totalVolume: BigInt(0),
      totalSlippage: 0,
      lastUpdateTime: Date.now(),
    };
  }
}
