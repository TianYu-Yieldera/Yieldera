/**
 * Slack Alert Service
 *
 * 发送监控告警到 Slack 频道
 */

import https from 'https';

export interface Alert {
  level: 'CRITICAL' | 'WARNING' | 'INFO';
  type: string;
  message: string;
  data?: any;
  timestamp?: number;
}

export class SlackAlertService {
  private webhookUrl: string;
  private enabled: boolean;
  private channelName: string;
  private botName: string;
  private minLevel: 'CRITICAL' | 'WARNING' | 'INFO';

  // 防止告警轰炸
  private recentAlerts: Map<string, number> = new Map();
  private alertCooldown: number = 300000; // 5分钟冷却

  constructor(config: {
    webhookUrl: string;
    enabled?: boolean;
    channelName?: string;
    botName?: string;
    minLevel?: 'CRITICAL' | 'WARNING' | 'INFO';
  }) {
    this.webhookUrl = config.webhookUrl;
    this.enabled = config.enabled !== false;
    this.channelName = config.channelName || '#monitoring';
    this.botName = config.botName || 'Loyalty Points Monitor';
    this.minLevel = config.minLevel || 'WARNING';
  }

  /**
   * 发送告警
   */
  async sendAlert(alert: Alert): Promise<boolean> {
    if (!this.enabled) {
      console.log('[Slack] Alerts disabled, skipping...');
      return false;
    }

    if (!this.webhookUrl || this.webhookUrl.includes('YOUR_WEBHOOK')) {
      console.log('[Slack] Webhook URL not configured');
      return false;
    }

    // 检查告警级别
    if (!this.shouldSendAlert(alert.level)) {
      console.log(`[Slack] Alert level ${alert.level} below minimum ${this.minLevel}, skipping...`);
      return false;
    }

    // 防止重复告警
    const alertKey = `${alert.type}:${alert.message}`;
    if (this.isRecentAlert(alertKey)) {
      console.log(`[Slack] Duplicate alert suppressed: ${alert.type}`);
      return false;
    }

    try {
      const payload = this.buildSlackMessage(alert);
      await this.postToSlack(payload);

      // 记录告警时间
      this.recentAlerts.set(alertKey, Date.now());

      console.log(`[Slack] Alert sent: ${alert.type} (${alert.level})`);
      return true;
    } catch (error: any) {
      console.error('[Slack] Failed to send alert:', error.message);
      return false;
    }
  }

  /**
   * 检查是否应发送告警
   */
  private shouldSendAlert(level: string): boolean {
    const levels = ['INFO', 'WARNING', 'CRITICAL'];
    const minIndex = levels.indexOf(this.minLevel);
    const currentIndex = levels.indexOf(level);
    return currentIndex >= minIndex;
  }

  /**
   * 检查是否为近期重复告警
   */
  private isRecentAlert(alertKey: string): boolean {
    const lastTime = this.recentAlerts.get(alertKey);
    if (!lastTime) return false;

    const elapsed = Date.now() - lastTime;
    return elapsed < this.alertCooldown;
  }

  /**
   * 构建 Slack 消息
   */
  private buildSlackMessage(alert: Alert): any {
    const emoji = this.getEmojiForLevel(alert.level);
    const color = this.getColorForLevel(alert.level);

    const timestamp = alert.timestamp || Date.now();
    const timeStr = new Date(timestamp).toISOString();

    // 构建附件字段
    const fields: any[] = [
      {
        title: 'Alert Type',
        value: alert.type,
        short: true,
      },
      {
        title: 'Level',
        value: alert.level,
        short: true,
      },
      {
        title: 'Time',
        value: timeStr,
        short: false,
      },
    ];

    // 如果有数据，添加关键信息
    if (alert.data) {
      if (alert.data.txHash) {
        fields.push({
          title: 'Transaction',
          value: `\`${alert.data.txHash}\``,
          short: false,
        });
      }
      if (alert.data.amount) {
        fields.push({
          title: 'Amount',
          value: alert.data.amount,
          short: true,
        });
      }
      if (alert.data.user) {
        fields.push({
          title: 'User',
          value: `\`${alert.data.user}\``,
          short: true,
        });
      }
    }

    return {
      username: this.botName,
      channel: this.channelName,
      icon_emoji: emoji,
      attachments: [
        {
          color,
          title: `${emoji} ${alert.level} Alert`,
          text: alert.message,
          fields,
          footer: 'Loyalty Points Monitoring System',
          footer_icon: 'https://platform.slack-edge.com/img/default_application_icon.png',
          ts: Math.floor(timestamp / 1000),
        },
      ],
    };
  }

  /**
   * 获取告警级别对应的 emoji
   */
  private getEmojiForLevel(level: string): string {
    switch (level) {
      case 'CRITICAL':
        return ':rotating_light:';
      case 'WARNING':
        return ':warning:';
      case 'INFO':
        return ':information_source:';
      default:
        return ':bell:';
    }
  }

  /**
   * 获取告警级别对应的颜色
   */
  private getColorForLevel(level: string): string {
    switch (level) {
      case 'CRITICAL':
        return '#FF0000'; // 红色
      case 'WARNING':
        return '#FFA500'; // 橙色
      case 'INFO':
        return '#36A64F'; // 绿色
      default:
        return '#808080'; // 灰色
    }
  }

  /**
   * POST 到 Slack Webhook
   */
  private postToSlack(payload: any): Promise<void> {
    return new Promise((resolve, reject) => {
      const data = JSON.stringify(payload);
      const url = new URL(this.webhookUrl);

      const options = {
        hostname: url.hostname,
        port: 443,
        path: url.pathname + url.search,
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Content-Length': Buffer.byteLength(data),
        },
      };

      const req = https.request(options, (res) => {
        let responseData = '';

        res.on('data', (chunk) => {
          responseData += chunk;
        });

        res.on('end', () => {
          if (res.statusCode === 200) {
            resolve();
          } else {
            reject(new Error(`Slack API returned ${res.statusCode}: ${responseData}`));
          }
        });
      });

      req.on('error', (error) => {
        reject(error);
      });

      req.write(data);
      req.end();
    });
  }

  /**
   * 发送测试消息
   */
  async sendTestMessage(): Promise<boolean> {
    const testAlert: Alert = {
      level: 'INFO',
      type: 'TEST',
      message: '✅ Slack integration test successful! The monitoring system is now connected.',
      timestamp: Date.now(),
    };

    return this.sendAlert(testAlert);
  }

  /**
   * 发送系统启动通知
   */
  async sendStartupNotification(): Promise<boolean> {
    const alert: Alert = {
      level: 'INFO',
      type: 'SYSTEM_STARTUP',
      message: '🚀 Loyalty Points Monitoring System started successfully!',
      timestamp: Date.now(),
    };

    return this.sendAlert(alert);
  }

  /**
   * 发送系统关闭通知
   */
  async sendShutdownNotification(): Promise<boolean> {
    const alert: Alert = {
      level: 'INFO',
      type: 'SYSTEM_SHUTDOWN',
      message: '🛑 Loyalty Points Monitoring System is shutting down.',
      timestamp: Date.now(),
    };

    return this.sendAlert(alert);
  }

  /**
   * 清理旧的告警记录
   */
  cleanupOldAlerts(): void {
    const now = Date.now();
    for (const [key, time] of this.recentAlerts.entries()) {
      if (now - time > this.alertCooldown * 2) {
        this.recentAlerts.delete(key);
      }
    }
  }

  /**
   * 获取统计信息
   */
  getStats() {
    return {
      enabled: this.enabled,
      webhookConfigured: !this.webhookUrl.includes('YOUR_WEBHOOK'),
      minLevel: this.minLevel,
      recentAlertsCount: this.recentAlerts.size,
      cooldownMs: this.alertCooldown,
    };
  }
}
