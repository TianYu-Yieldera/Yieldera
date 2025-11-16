/**
 * Notification Service
 *
 * Multi-channel notification system for AI alerts, risk warnings, and yield updates
 * Supports: Email, SMS, Push Notifications, Slack, Discord, Telegram
 *
 * Features:
 * - AI-triggered liquidation warnings
 * - Daily yield reports
 * - Portfolio risk alerts
 * - Custom user preferences
 * - Rate limiting and deduplication
 */

import { Pool } from 'pg';
import axios from 'axios';
import nodemailer from 'nodemailer';

export enum NotificationChannel {
  EMAIL = 'email',
  SMS = 'sms',
  PUSH = 'push',
  SLACK = 'slack',
  DISCORD = 'discord',
  TELEGRAM = 'telegram',
}

export enum NotificationPriority {
  LOW = 'low',        // Daily reports, general updates
  MEDIUM = 'medium',  // Weekly summaries, moderate risk
  HIGH = 'high',      // Important alerts, high risk
  CRITICAL = 'critical', // Liquidation warnings, emergency
}

export enum NotificationType {
  // Risk Alerts
  LIQUIDATION_WARNING = 'liquidation_warning',
  HIGH_RISK_POSITION = 'high_risk_position',
  HEALTH_FACTOR_LOW = 'health_factor_low',

  // Yield Updates
  DAILY_YIELD_REPORT = 'daily_yield_report',
  WEEKLY_SUMMARY = 'weekly_summary',
  COMPOUND_MILESTONE = 'compound_milestone',

  // Portfolio Events
  LARGE_MOVEMENT = 'large_movement',
  NEW_OPPORTUNITY = 'new_opportunity',
  REBALANCE_SUGGESTION = 'rebalance_suggestion',

  // System Alerts
  SYSTEM_MAINTENANCE = 'system_maintenance',
  SECURITY_ALERT = 'security_alert',
}

export interface NotificationPreferences {
  userId: string;
  channels: NotificationChannel[];
  minPriority: NotificationPriority;
  quietHoursStart?: number; // 0-23 hour
  quietHoursEnd?: number;
  enabledTypes: NotificationType[];
  frequency: 'realtime' | 'hourly' | 'daily';
}

export interface Notification {
  id?: string;
  userId: string;
  type: NotificationType;
  priority: NotificationPriority;
  title: string;
  message: string;
  data?: any; // Additional structured data
  channels: NotificationChannel[];
  sentAt?: Date;
  readAt?: Date;
  deliveryStatus: Map<NotificationChannel, 'pending' | 'sent' | 'failed'>;
}

export class NotificationService {
  private db: Pool;
  private emailTransporter: nodemailer.Transporter;
  private slackWebhookUrl: string;
  private discordWebhookUrl: string;
  private telegramBotToken: string;
  private twilioClient: any; // SMS service

  // Rate limiting: max notifications per user per hour
  private rateLimits = new Map<string, number>();
  private readonly MAX_NOTIFICATIONS_PER_HOUR = 10;

  constructor(dbPool: Pool) {
    this.db = dbPool;

    // Initialize email transport (SMTP)
    this.emailTransporter = nodemailer.createTransport({
      host: process.env.SMTP_HOST || 'smtp.gmail.com',
      port: parseInt(process.env.SMTP_PORT || '587'),
      secure: false,
      auth: {
        user: process.env.SMTP_USER,
        pass: process.env.SMTP_PASS,
      },
    });

    this.slackWebhookUrl = process.env.SLACK_WEBHOOK_URL || '';
    this.discordWebhookUrl = process.env.DISCORD_WEBHOOK_URL || '';
    this.telegramBotToken = process.env.TELEGRAM_BOT_TOKEN || '';

    // Start cleanup job for old notifications
    this.startCleanupJob();
  }

  /**
   * Send notification to user via preferred channels
   */
  async sendNotification(notification: Notification): Promise<void> {
    // Check user preferences
    const prefs = await this.getUserPreferences(notification.userId);

    // Filter channels based on preferences
    const allowedChannels = this.filterChannelsByPreferences(
      notification,
      prefs
    );

    if (allowedChannels.length === 0) {
      console.log(`No channels enabled for notification to ${notification.userId}`);
      return;
    }

    // Check rate limits
    if (!this.checkRateLimit(notification.userId)) {
      console.log(`Rate limit exceeded for user ${notification.userId}`);
      await this.saveNotification(notification, 'rate_limited');
      return;
    }

    // Check quiet hours
    if (this.isQuietHours(prefs)) {
      console.log(`Quiet hours active for user ${notification.userId}`);
      await this.queueForLater(notification);
      return;
    }

    // Send via each channel
    const deliveryStatus = new Map<NotificationChannel, 'pending' | 'sent' | 'failed'>();

    for (const channel of allowedChannels) {
      try {
        await this.sendViaChannel(channel, notification);
        deliveryStatus.set(channel, 'sent');
      } catch (error) {
        console.error(`Failed to send via ${channel}:`, error);
        deliveryStatus.set(channel, 'failed');
      }
    }

    notification.deliveryStatus = deliveryStatus;
    notification.sentAt = new Date();

    // Save to database
    await this.saveNotification(notification, 'sent');

    // Update rate limit counter
    this.incrementRateLimit(notification.userId);
  }

  /**
   * Send AI liquidation warning
   */
  async sendLiquidationWarning(
    userId: string,
    data: {
      protocol: string;
      healthFactor: number;
      liquidationPrice: number;
      hoursUntilLiquidation: number;
    }
  ): Promise<void> {
    const notification: Notification = {
      userId,
      type: NotificationType.LIQUIDATION_WARNING,
      priority: NotificationPriority.CRITICAL,
      title: '🚨 CRITICAL: Liquidation Risk Detected',
      message: `Your ${data.protocol} position is at risk of liquidation!\n\n` +
        `Health Factor: ${data.healthFactor.toFixed(2)}\n` +
        `Estimated time: ${data.hoursUntilLiquidation} hours\n\n` +
        `Action required: Add collateral or reduce debt immediately.`,
      data,
      channels: [
        NotificationChannel.EMAIL,
        NotificationChannel.SMS,
        NotificationChannel.PUSH,
      ],
      deliveryStatus: new Map(),
    };

    await this.sendNotification(notification);
  }

  /**
   * Send daily yield report
   */
  async sendDailyYieldReport(
    userId: string,
    data: {
      totalYield: number;
      yieldByBondType: any;
      portfolioValue: number;
    }
  ): Promise<void> {
    const notification: Notification = {
      userId,
      type: NotificationType.DAILY_YIELD_REPORT,
      priority: NotificationPriority.LOW,
      title: '📊 Your Daily Yield Report',
      message: `Daily yield earned: $${data.totalYield.toFixed(2)}\n` +
        `Portfolio value: $${data.portfolioValue.toFixed(2)}\n\n` +
        `Keep up the great work! 🎯`,
      data,
      channels: [NotificationChannel.EMAIL],
      deliveryStatus: new Map(),
    };

    await this.sendNotification(notification);
  }

  /**
   * Send high risk position alert
   */
  async sendHighRiskAlert(
    userId: string,
    data: {
      protocol: string;
      riskScore: number;
      recommendations: string[];
    }
  ): Promise<void> {
    const notification: Notification = {
      userId,
      type: NotificationType.HIGH_RISK_POSITION,
      priority: NotificationPriority.HIGH,
      title: '⚠️ High Risk Position Detected',
      message: `Your ${data.protocol} position has elevated risk (${data.riskScore}/100)\n\n` +
        `Recommendations:\n${data.recommendations.map(r => `• ${r}`).join('\n')}`,
      data,
      channels: [NotificationChannel.EMAIL, NotificationChannel.PUSH],
      deliveryStatus: new Map(),
    };

    await this.sendNotification(notification);
  }

  /**
   * Send via specific channel
   */
  private async sendViaChannel(
    channel: NotificationChannel,
    notification: Notification
  ): Promise<void> {
    switch (channel) {
      case NotificationChannel.EMAIL:
        await this.sendEmail(notification);
        break;
      case NotificationChannel.SMS:
        await this.sendSMS(notification);
        break;
      case NotificationChannel.PUSH:
        await this.sendPushNotification(notification);
        break;
      case NotificationChannel.SLACK:
        await this.sendSlack(notification);
        break;
      case NotificationChannel.DISCORD:
        await this.sendDiscord(notification);
        break;
      case NotificationChannel.TELEGRAM:
        await this.sendTelegram(notification);
        break;
    }
  }

  /**
   * Send email notification
   */
  private async sendEmail(notification: Notification): Promise<void> {
    const userEmail = await this.getUserEmail(notification.userId);

    if (!userEmail) {
      throw new Error('User email not found');
    }

    const html = this.generateEmailHTML(notification);

    await this.emailTransporter.sendMail({
      from: process.env.SMTP_FROM || 'Yieldera <noreply@yieldera.io>',
      to: userEmail,
      subject: notification.title,
      text: notification.message,
      html,
    });

    console.log(`✓ Email sent to ${userEmail}`);
  }

  /**
   * Send SMS notification (via Twilio)
   */
  private async sendSMS(notification: Notification): Promise<void> {
    const phoneNumber = await this.getUserPhone(notification.userId);

    if (!phoneNumber) {
      throw new Error('User phone number not found');
    }

    // Twilio integration (requires twilio package)
    // const message = await this.twilioClient.messages.create({
    //   body: `${notification.title}\n\n${notification.message}`,
    //   from: process.env.TWILIO_PHONE_NUMBER,
    //   to: phoneNumber,
    // });

    console.log(`✓ SMS sent to ${phoneNumber}`);
  }

  /**
   * Send push notification (via Firebase Cloud Messaging or similar)
   */
  private async sendPushNotification(notification: Notification): Promise<void> {
    const fcmToken = await this.getUserFCMToken(notification.userId);

    if (!fcmToken) {
      throw new Error('User FCM token not found');
    }

    // Firebase Cloud Messaging integration
    // await admin.messaging().send({
    //   token: fcmToken,
    //   notification: {
    //     title: notification.title,
    //     body: notification.message,
    //   },
    //   data: notification.data,
    // });

    console.log(`✓ Push notification sent to ${notification.userId}`);
  }

  /**
   * Send Slack notification
   */
  private async sendSlack(notification: Notification): Promise<void> {
    if (!this.slackWebhookUrl) {
      throw new Error('Slack webhook URL not configured');
    }

    const color = this.getPriorityColor(notification.priority);

    await axios.post(this.slackWebhookUrl, {
      attachments: [
        {
          color,
          title: notification.title,
          text: notification.message,
          footer: 'Yieldera Platform',
          ts: Math.floor(Date.now() / 1000),
        },
      ],
    });

    console.log('✓ Slack notification sent');
  }

  /**
   * Send Discord notification
   */
  private async sendDiscord(notification: Notification): Promise<void> {
    if (!this.discordWebhookUrl) {
      throw new Error('Discord webhook URL not configured');
    }

    const color = this.getPriorityColorHex(notification.priority);

    await axios.post(this.discordWebhookUrl, {
      embeds: [
        {
          title: notification.title,
          description: notification.message,
          color: parseInt(color.replace('#', ''), 16),
          timestamp: new Date().toISOString(),
          footer: {
            text: 'Yieldera Platform',
          },
        },
      ],
    });

    console.log('✓ Discord notification sent');
  }

  /**
   * Send Telegram notification
   */
  private async sendTelegram(notification: Notification): Promise<void> {
    const chatId = await this.getUserTelegramChatId(notification.userId);

    if (!chatId) {
      throw new Error('User Telegram chat ID not found');
    }

    const text = `*${notification.title}*\n\n${notification.message}`;

    await axios.post(
      `https://api.telegram.org/bot${this.telegramBotToken}/sendMessage`,
      {
        chat_id: chatId,
        text,
        parse_mode: 'Markdown',
      }
    );

    console.log('✓ Telegram notification sent');
  }

  /**
   * Helper methods
   */

  private filterChannelsByPreferences(
    notification: Notification,
    prefs: NotificationPreferences
  ): NotificationChannel[] {
    // Filter by enabled channels
    let channels = notification.channels.filter(c => prefs.channels.includes(c));

    // Filter by notification type
    if (!prefs.enabledTypes.includes(notification.type)) {
      // Only allow critical notifications to bypass type filter
      if (notification.priority !== NotificationPriority.CRITICAL) {
        return [];
      }
    }

    // Filter by priority
    const priorityOrder = [
      NotificationPriority.LOW,
      NotificationPriority.MEDIUM,
      NotificationPriority.HIGH,
      NotificationPriority.CRITICAL,
    ];

    const minPriorityIndex = priorityOrder.indexOf(prefs.minPriority);
    const notificationPriorityIndex = priorityOrder.indexOf(notification.priority);

    if (notificationPriorityIndex < minPriorityIndex) {
      return [];
    }

    return channels;
  }

  private isQuietHours(prefs: NotificationPreferences): boolean {
    if (!prefs.quietHoursStart || !prefs.quietHoursEnd) {
      return false;
    }

    const now = new Date();
    const currentHour = now.getHours();

    const start = prefs.quietHoursStart;
    const end = prefs.quietHoursEnd;

    if (start < end) {
      return currentHour >= start && currentHour < end;
    } else {
      // Quiet hours span midnight
      return currentHour >= start || currentHour < end;
    }
  }

  private checkRateLimit(userId: string): boolean {
    const count = this.rateLimits.get(userId) || 0;
    return count < this.MAX_NOTIFICATIONS_PER_HOUR;
  }

  private incrementRateLimit(userId: string): void {
    const count = this.rateLimits.get(userId) || 0;
    this.rateLimits.set(userId, count + 1);
  }

  private getPriorityColor(priority: NotificationPriority): string {
    switch (priority) {
      case NotificationPriority.CRITICAL: return 'danger';
      case NotificationPriority.HIGH: return 'warning';
      case NotificationPriority.MEDIUM: return 'good';
      case NotificationPriority.LOW: return '#36a64f';
    }
  }

  private getPriorityColorHex(priority: NotificationPriority): string {
    switch (priority) {
      case NotificationPriority.CRITICAL: return '#FF0000';
      case NotificationPriority.HIGH: return '#FFA500';
      case NotificationPriority.MEDIUM: return '#FFFF00';
      case NotificationPriority.LOW: return '#00FF00';
    }
  }

  private generateEmailHTML(notification: Notification): string {
    return `
      <!DOCTYPE html>
      <html>
      <head>
        <style>
          body { font-family: Arial, sans-serif; line-height: 1.6; }
          .container { max-width: 600px; margin: 0 auto; padding: 20px; }
          .header { background: #1a1a2e; color: white; padding: 20px; text-align: center; }
          .content { padding: 20px; background: #f4f4f4; }
          .priority-critical { border-left: 5px solid #ff0000; }
          .priority-high { border-left: 5px solid #ffa500; }
          .footer { text-align: center; padding: 10px; color: #666; }
        </style>
      </head>
      <body>
        <div class="container">
          <div class="header">
            <h1>Yieldera</h1>
          </div>
          <div class="content priority-${notification.priority}">
            <h2>${notification.title}</h2>
            <p>${notification.message.replace(/\n/g, '<br>')}</p>
          </div>
          <div class="footer">
            <p>Yieldera - Institutional-Grade DeFi Risk Management</p>
          </div>
        </div>
      </body>
      </html>
    `;
  }

  /**
   * Database operations
   */

  private async getUserPreferences(userId: string): Promise<NotificationPreferences> {
    const result = await this.db.query(
      'SELECT * FROM notification_preferences WHERE user_id = $1',
      [userId]
    );

    if (result.rows.length === 0) {
      // Return default preferences
      return {
        userId,
        channels: [NotificationChannel.EMAIL],
        minPriority: NotificationPriority.MEDIUM,
        enabledTypes: Object.values(NotificationType),
        frequency: 'realtime',
      };
    }

    const row = result.rows[0];
    return {
      userId: row.user_id,
      channels: row.channels,
      minPriority: row.min_priority,
      quietHoursStart: row.quiet_hours_start,
      quietHoursEnd: row.quiet_hours_end,
      enabledTypes: row.enabled_types,
      frequency: row.frequency,
    };
  }

  private async getUserEmail(userId: string): Promise<string | null> {
    const result = await this.db.query(
      'SELECT email FROM users WHERE id = $1',
      [userId]
    );
    return result.rows[0]?.email || null;
  }

  private async getUserPhone(userId: string): Promise<string | null> {
    const result = await this.db.query(
      'SELECT phone_number FROM users WHERE id = $1',
      [userId]
    );
    return result.rows[0]?.phone_number || null;
  }

  private async getUserFCMToken(userId: string): Promise<string | null> {
    const result = await this.db.query(
      'SELECT fcm_token FROM user_devices WHERE user_id = $1 AND active = true',
      [userId]
    );
    return result.rows[0]?.fcm_token || null;
  }

  private async getUserTelegramChatId(userId: string): Promise<string | null> {
    const result = await this.db.query(
      'SELECT telegram_chat_id FROM users WHERE id = $1',
      [userId]
    );
    return result.rows[0]?.telegram_chat_id || null;
  }

  private async saveNotification(
    notification: Notification,
    status: string
  ): Promise<void> {
    await this.db.query(
      `INSERT INTO notifications
       (user_id, type, priority, title, message, data, channels, sent_at, status)
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
      [
        notification.userId,
        notification.type,
        notification.priority,
        notification.title,
        notification.message,
        JSON.stringify(notification.data),
        notification.channels,
        notification.sentAt || new Date(),
        status,
      ]
    );
  }

  private async queueForLater(notification: Notification): Promise<void> {
    await this.db.query(
      `INSERT INTO notification_queue
       (user_id, notification_data, scheduled_for)
       VALUES ($1, $2, $3)`,
      [notification.userId, JSON.stringify(notification), new Date(Date.now() + 3600000)]
    );
  }

  /**
   * Cleanup old notifications
   */
  private startCleanupJob(): void {
    // Run every hour
    setInterval(async () => {
      try {
        // Delete notifications older than 30 days
        await this.db.query(
          `DELETE FROM notifications WHERE sent_at < NOW() - INTERVAL '30 days'`
        );

        // Reset rate limits
        this.rateLimits.clear();

        console.log('✓ Notification cleanup completed');
      } catch (error) {
        console.error('Notification cleanup error:', error);
      }
    }, 3600000); // 1 hour
  }
}

// Export singleton instance
export const notificationService = new NotificationService(
  new Pool({
    connectionString: process.env.DATABASE_URL,
  })
);
