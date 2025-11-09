/**
 * GMX Position Listener
 *
 * 监听 GMXV2Adapter 合约的仓位事件:
 * - PositionOpened - 仓位开仓
 * - PositionClosed - 仓位平仓
 * - EmergencyHedgeExecuted - 紧急对冲执行
 *
 * 核心功能:
 * 1. 实时监控仓位风险
 * 2. 生成风险建议 (Advisory Mode)
 * 3. 追踪对冲效果
 * 4. Slack 告警通知
 */

import { ethers } from 'ethers';
import { BaseListener } from '../BaseListener';
import { db } from '../../config/database';

// 仓位统计
interface PositionStats {
  totalLongSize: bigint;
  totalShortSize: bigint;
  totalPositions: number;
  hedgePositions: number;
  highLeveragePositions: number;
  liquidationWarnings: number;
}

// 仓位数据
interface PositionData {
  user: string;
  orderKey: string;
  market: string;
  collateralToken: string;
  isLong: boolean;
  sizeInUsd: bigint;
  collateralAmount: bigint;
  leverage: bigint;
  isHedge: boolean;
}

// 风险建议
interface RiskRecommendation {
  level: 'INFO' | 'WARNING' | 'CRITICAL';
  type: string;
  message: string;
  action: string;
  priority: 'LOW' | 'MEDIUM' | 'HIGH';
  reason: string;
  expectedOutcome: string;
  userDecision: boolean; // 需要用户确认
}

export class GMXPositionListener extends BaseListener {
  private stats: PositionStats = {
    totalLongSize: BigInt(0),
    totalShortSize: BigInt(0),
    totalPositions: 0,
    hedgePositions: 0,
    highLeveragePositions: 0,
    liquidationWarnings: 0,
  };

  // 风险阈值 (可配置)
  private readonly LEVERAGE_WARNING = 30;    // 杠杆警告线
  private readonly LEVERAGE_CRITICAL = 40;   // 杠杆危险线
  private readonly LARGE_POSITION_USD = ethers.parseEther('50000'); // 大额仓位阈值

  constructor(wsUrl: string, contractAddress: string) {
    super(wsUrl, contractAddress, 'GMXPosition');
  }

  /**
   * 设置事件监听器
   */
  async setupEventListeners(): Promise<void> {
    // 1. 监听 PositionOpened 事件
    this.contract.on(
      'PositionOpened',
      async (
        user: string,
        orderKey: string,
        market: string,
        collateralToken: string,
        isLong: boolean,
        sizeInUsd: bigint,
        collateralAmount: bigint,
        leverage: bigint,
        isHedge: boolean,
        event: any
      ) => {
        try {
          await this.handlePositionOpened({
            user,
            orderKey,
            market,
            collateralToken,
            isLong,
            sizeInUsd,
            collateralAmount,
            leverage,
            isHedge,
          });
        } catch (error) {
          this.emit('error', { event: 'PositionOpened', error });
        }
      }
    );

    // 2. 监听 PositionClosed 事件
    this.contract.on(
      'PositionClosed',
      async (
        user: string,
        orderKey: string,
        market: string,
        sizeInUsd: bigint,
        pnl: bigint,
        event: any
      ) => {
        try {
          await this.handlePositionClosed({
            user,
            orderKey,
            market,
            sizeInUsd,
            pnl,
          });
        } catch (error) {
          this.emit('error', { event: 'PositionClosed', error });
        }
      }
    );

    // 3. 监听 EmergencyHedgeExecuted 事件 (最重要!)
    this.contract.on(
      'EmergencyHedgeExecuted',
      async (
        user: string,
        market: string,
        hedgeSize: bigint,
        reason: string,
        orderKey: string,
        event: any
      ) => {
        try {
          await this.handleEmergencyHedge({
            user,
            market,
            hedgeSize,
            reason,
            orderKey,
          });
        } catch (error) {
          this.emit('error', { event: 'EmergencyHedgeExecuted', error });
        }
      }
    );

    console.log(`✅ GMXPositionListener: Listening to ${this.contractAddress}`);
  }

  /**
   * 处理仓位开仓事件
   */
  private async handlePositionOpened(data: PositionData): Promise<void> {
    console.log(`\n🎯 GMX Position Opened`);
    console.log(`  User: ${data.user}`);
    console.log(`  Market: ${data.market}`);
    console.log(`  Direction: ${data.isLong ? 'LONG' : 'SHORT'}`);
    console.log(`  Size: ${ethers.formatEther(data.sizeInUsd)} USD`);
    console.log(`  Leverage: ${data.leverage}x`);
    console.log(`  Is Hedge: ${data.isHedge ? 'YES' : 'NO'}`);

    // 更新统计
    this.updateStats(data);

    // 保存到数据库
    await this.savePosition(data);

    // 风险评估 (建议式)
    const recommendation = this.assessRisk(data);
    if (recommendation) {
      this.emit('alert', {
        ...recommendation,
        user: data.user,
        position: data,
      });
    }

    // 发射事件给上层处理
    this.emit('positionOpened', {
      user: data.user,
      market: data.market,
      isLong: data.isLong,
      sizeUsd: ethers.formatEther(data.sizeInUsd),
      leverage: data.leverage.toString(),
      isHedge: data.isHedge,
      timestamp: new Date(),
    });
  }

  /**
   * 处理仓位平仓事件
   */
  private async handlePositionClosed(data: any): Promise<void> {
    console.log(`\n✅ GMX Position Closed`);
    console.log(`  User: ${data.user}`);
    console.log(`  Market: ${data.market}`);
    console.log(`  Size: ${ethers.formatEther(data.sizeInUsd)} USD`);
    console.log(`  PnL: ${ethers.formatEther(data.pnl)} USD`);

    // 更新统计
    this.stats.totalPositions--;

    // 更新数据库
    await this.updateClosedPosition(data);

    // 分析 PnL
    const isProfitable = data.pnl > 0;
    const pnlPercent = this.calculatePnLPercent(data);

    // 发射事件
    this.emit('positionClosed', {
      user: data.user,
      market: data.market,
      sizeUsd: ethers.formatEther(data.sizeInUsd),
      pnl: ethers.formatEther(data.pnl),
      profitable: isProfitable,
      pnlPercent: pnlPercent,
      timestamp: new Date(),
    });

    // 如果是大额亏损，发送告警
    if (!isProfitable && data.pnl < -this.LARGE_POSITION_USD) {
      this.emit('alert', {
        level: 'WARNING',
        type: 'LARGE_LOSS',
        message: `⚠️ 大额亏损: ${ethers.formatEther(-data.pnl)} USD`,
        user: data.user,
        recommendation: {
          action: 'REVIEW_STRATEGY',
          priority: 'MEDIUM',
          reason: '单次亏损超过 50k USD，建议复盘交易策略',
          expectedOutcome: '优化未来交易决策，降低风险',
          userDecision: true,
        },
      });
    }
  }

  /**
   * 处理紧急对冲事件 (风控核心)
   */
  private async handleEmergencyHedge(data: any): Promise<void> {
    console.log(`\n🚨 Emergency Hedge Executed!`);
    console.log(`  User: ${data.user}`);
    console.log(`  Market: ${data.market}`);
    console.log(`  Hedge Size: ${ethers.formatEther(data.hedgeSize)} USD`);
    console.log(`  Reason: ${data.reason}`);

    this.stats.hedgePositions++;

    // 保存对冲记录
    await this.saveHedgeRecord(data);

    // 发送高优先级告警 (这是自动风控触发的)
    this.emit('alert', {
      level: 'CRITICAL',
      type: 'EMERGENCY_HEDGE',
      message: `🚨 紧急对冲已执行`,
      user: data.user,
      details: {
        market: data.market,
        hedgeSize: ethers.formatEther(data.hedgeSize),
        reason: data.reason,
      },
      recommendation: {
        action: 'REVIEW_POSITION',
        priority: 'HIGH',
        reason: data.reason,
        expectedOutcome: '风险敞口已降低，请检查仓位状态',
        userDecision: false, // 已自动执行
      },
    });

    // 发射事件
    this.emit('emergencyHedge', {
      user: data.user,
      market: data.market,
      hedgeSize: ethers.formatEther(data.hedgeSize),
      reason: data.reason,
      timestamp: new Date(),
    });
  }

  /**
   * 风险评估 (生成建议)
   */
  private assessRisk(data: PositionData): RiskRecommendation | null {
    const leverage = Number(data.leverage);
    const sizeInUsd = data.sizeInUsd;

    // 1. 杠杆风险评估
    if (leverage >= this.LEVERAGE_CRITICAL) {
      return {
        level: 'CRITICAL',
        type: 'EXTREME_LEVERAGE',
        message: `🚨 极高杠杆风险: ${leverage}x`,
        action: 'REDUCE_LEVERAGE_URGENT',
        priority: 'HIGH',
        reason: `当前杠杆 ${leverage}x 超过危险线 (${this.LEVERAGE_CRITICAL}x)，极易被清算`,
        expectedOutcome: '建议立即平仓 50% 或增加抵押品，将杠杆降至 20x 以下',
        userDecision: true,
      };
    } else if (leverage >= this.LEVERAGE_WARNING) {
      return {
        level: 'WARNING',
        type: 'HIGH_LEVERAGE',
        message: `⚠️ 高杠杆警告: ${leverage}x`,
        action: 'REDUCE_LEVERAGE',
        priority: 'MEDIUM',
        reason: `当前杠杆 ${leverage}x 接近警告线 (${this.LEVERAGE_WARNING}x)`,
        expectedOutcome: '建议降低杠杆至 25x 以下，增加安全边际',
        userDecision: true,
      };
    }

    // 2. 大额仓位提醒
    if (sizeInUsd >= this.LARGE_POSITION_USD) {
      return {
        level: 'INFO',
        type: 'LARGE_POSITION',
        message: `📊 大额仓位: ${ethers.formatEther(sizeInUsd)} USD`,
        action: 'MONITOR_CLOSELY',
        priority: 'LOW',
        reason: '仓位规模较大，建议密切监控',
        expectedOutcome: '及时关注市场波动，避免大额损失',
        userDecision: false,
      };
    }

    return null;
  }

  /**
   * 更新统计数据
   */
  private updateStats(data: PositionData): void {
    this.stats.totalPositions++;

    if (data.isLong) {
      this.stats.totalLongSize += data.sizeInUsd;
    } else {
      this.stats.totalShortSize += data.sizeInUsd;
    }

    if (data.isHedge) {
      this.stats.hedgePositions++;
    }

    if (Number(data.leverage) >= this.LEVERAGE_WARNING) {
      this.stats.highLeveragePositions++;
    }
  }

  /**
   * 保存仓位到数据库
   */
  private async savePosition(data: PositionData): Promise<void> {
    try {
      await db.query(`
        INSERT INTO gmx_positions (
          user_address,
          order_key,
          market,
          collateral_token,
          is_long,
          size_usd,
          collateral_amount,
          leverage,
          is_hedge,
          status,
          created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
      `, [
        data.user,
        data.orderKey,
        data.market,
        data.collateralToken,
        data.isLong,
        data.sizeInUsd.toString(),
        data.collateralAmount.toString(),
        data.leverage.toString(),
        data.isHedge,
        'open',
      ]);
    } catch (error) {
      console.error('Error saving position to database:', error);
    }
  }

  /**
   * 更新平仓仓位
   */
  private async updateClosedPosition(data: any): Promise<void> {
    try {
      await db.query(`
        UPDATE gmx_positions
        SET status = 'closed',
            closed_pnl = $1,
            closed_at = NOW()
        WHERE order_key = $2
      `, [
        data.pnl.toString(),
        data.orderKey,
      ]);
    } catch (error) {
      console.error('Error updating closed position:', error);
    }
  }

  /**
   * 保存对冲记录
   */
  private async saveHedgeRecord(data: any): Promise<void> {
    try {
      await db.query(`
        INSERT INTO gmx_hedge_records (
          user_address,
          market,
          hedge_size,
          reason,
          order_key,
          created_at
        ) VALUES ($1, $2, $3, $4, $5, NOW())
      `, [
        data.user,
        data.market,
        data.hedgeSize.toString(),
        data.reason,
        data.orderKey,
      ]);
    } catch (error) {
      console.error('Error saving hedge record:', error);
    }
  }

  /**
   * 计算 PnL 百分比
   */
  private calculatePnLPercent(data: any): number {
    // 简化计算: PnL / Size * 100
    if (data.sizeInUsd === BigInt(0)) return 0;
    return Number((data.pnl * BigInt(10000)) / data.sizeInUsd) / 100;
  }

  /**
   * 获取统计数据
   */
  getStats(): PositionStats & { stats: any } {
    return {
      ...this.stats,
      stats: {
        totalLongUsd: ethers.formatEther(this.stats.totalLongSize),
        totalShortUsd: ethers.formatEther(this.stats.totalShortSize),
        totalPositions: this.stats.totalPositions,
        hedgePositions: this.stats.hedgePositions,
        highLeveragePositions: this.stats.highLeveragePositions,
      },
    };
  }

  /**
   * 获取监听器类型
   */
  getType(): string {
    return 'GMX Position Listener';
  }
}
