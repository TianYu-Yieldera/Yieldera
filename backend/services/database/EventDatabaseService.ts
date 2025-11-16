/**
 * Event Database Service
 *
 * 持久化监控事件到PostgreSQL数据库
 */

import { Pool, PoolClient } from 'pg';

export interface EventRecord {
  id?: number;
  eventType: string;
  contractName: string;
  contractAddress: string;
  blockNumber: number;
  transactionHash: string;
  eventData: any;
  timestamp: Date;
  processed: boolean;
}

export interface AlertRecord {
  id?: number;
  level: 'CRITICAL' | 'WARNING' | 'INFO';
  type: string;
  message: string;
  data?: any;
  timestamp: Date;
  acknowledged: boolean;
}

export class EventDatabaseService {
  private pool: Pool;
  private enabled: boolean;

  constructor(config: {
    host: string;
    port: number;
    database: string;
    user: string;
    password: string;
    enabled?: boolean;
  }) {
    this.enabled = config.enabled !== false;

    if (this.enabled) {
      this.pool = new Pool({
        host: config.host,
        port: config.port,
        database: config.database,
        user: config.user,
        password: config.password,
        max: 20, // 最大连接数
        idleTimeoutMillis: 30000,
        connectionTimeoutMillis: 2000,
      });

      // 监听错误
      this.pool.on('error', (err) => {
        console.error('[DB] Unexpected error on idle client', err);
      });
    }
  }

  /**
   * 初始化数据库表
   */
  async initialize(): Promise<void> {
    if (!this.enabled) {
      console.log('[DB] Database persistence disabled');
      return;
    }

    const client = await this.pool.connect();

    try {
      // 创建 events 表
      await client.query(`
        CREATE TABLE IF NOT EXISTS events (
          id SERIAL PRIMARY KEY,
          event_type VARCHAR(100) NOT NULL,
          contract_name VARCHAR(100) NOT NULL,
          contract_address VARCHAR(42) NOT NULL,
          block_number BIGINT NOT NULL,
          transaction_hash VARCHAR(66) NOT NULL,
          event_data JSONB NOT NULL,
          timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
          processed BOOLEAN DEFAULT FALSE,
          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
      `);

      // 创建索引
      await client.query(`
        CREATE INDEX IF NOT EXISTS idx_events_contract_name ON events(contract_name);
        CREATE INDEX IF NOT EXISTS idx_events_event_type ON events(event_type);
        CREATE INDEX IF NOT EXISTS idx_events_block_number ON events(block_number);
        CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
        CREATE INDEX IF NOT EXISTS idx_events_tx_hash ON events(transaction_hash);
      `);

      // 创建 alerts 表
      await client.query(`
        CREATE TABLE IF NOT EXISTS alerts (
          id SERIAL PRIMARY KEY,
          level VARCHAR(20) NOT NULL,
          type VARCHAR(100) NOT NULL,
          message TEXT NOT NULL,
          data JSONB,
          timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
          acknowledged BOOLEAN DEFAULT FALSE,
          acknowledged_at TIMESTAMP WITH TIME ZONE,
          acknowledged_by VARCHAR(100),
          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
      `);

      // 创建索引
      await client.query(`
        CREATE INDEX IF NOT EXISTS idx_alerts_level ON alerts(level);
        CREATE INDEX IF NOT EXISTS idx_alerts_type ON alerts(type);
        CREATE INDEX IF NOT EXISTS idx_alerts_timestamp ON alerts(timestamp);
        CREATE INDEX IF NOT EXISTS idx_alerts_acknowledged ON alerts(acknowledged);
      `);

      // 创建统计表
      await client.query(`
        CREATE TABLE IF NOT EXISTS statistics (
          id SERIAL PRIMARY KEY,
          contract_name VARCHAR(100) NOT NULL,
          metric_name VARCHAR(100) NOT NULL,
          metric_value NUMERIC NOT NULL,
          metadata JSONB,
          timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
      `);

      await client.query(`
        CREATE INDEX IF NOT EXISTS idx_stats_contract_metric ON statistics(contract_name, metric_name);
        CREATE INDEX IF NOT EXISTS idx_stats_timestamp ON statistics(timestamp);
      `);

      console.log('[DB] Database tables initialized successfully');
    } catch (error: any) {
      console.error('[DB] Failed to initialize tables:', error.message);
      throw error;
    } finally {
      client.release();
    }
  }

  /**
   * 保存事件
   */
  async saveEvent(event: EventRecord): Promise<number | null> {
    if (!this.enabled) return null;

    try {
      const result = await this.pool.query(
        `INSERT INTO events (
          event_type, contract_name, contract_address,
          block_number, transaction_hash, event_data,
          timestamp, processed
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id`,
        [
          event.eventType,
          event.contractName,
          event.contractAddress,
          event.blockNumber,
          event.transactionHash,
          JSON.stringify(event.eventData),
          event.timestamp,
          event.processed || false,
        ]
      );

      return result.rows[0].id;
    } catch (error: any) {
      console.error('[DB] Failed to save event:', error.message);
      return null;
    }
  }

  /**
   * 批量保存事件
   */
  async saveEvents(events: EventRecord[]): Promise<number> {
    if (!this.enabled || events.length === 0) return 0;

    const client = await this.pool.connect();

    try {
      await client.query('BEGIN');

      let savedCount = 0;
      for (const event of events) {
        await client.query(
          `INSERT INTO events (
            event_type, contract_name, contract_address,
            block_number, transaction_hash, event_data,
            timestamp, processed
          ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
          [
            event.eventType,
            event.contractName,
            event.contractAddress,
            event.blockNumber,
            event.transactionHash,
            JSON.stringify(event.eventData),
            event.timestamp,
            event.processed || false,
          ]
        );
        savedCount++;
      }

      await client.query('COMMIT');
      return savedCount;
    } catch (error: any) {
      await client.query('ROLLBACK');
      console.error('[DB] Failed to save events batch:', error.message);
      return 0;
    } finally {
      client.release();
    }
  }

  /**
   * 保存告警
   */
  async saveAlert(alert: AlertRecord): Promise<number | null> {
    if (!this.enabled) return null;

    try {
      const result = await this.pool.query(
        `INSERT INTO alerts (
          level, type, message, data, timestamp, acknowledged
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`,
        [
          alert.level,
          alert.type,
          alert.message,
          alert.data ? JSON.stringify(alert.data) : null,
          alert.timestamp,
          alert.acknowledged || false,
        ]
      );

      return result.rows[0].id;
    } catch (error: any) {
      console.error('[DB] Failed to save alert:', error.message);
      return null;
    }
  }

  /**
   * 获取最近的事件
   */
  async getRecentEvents(limit: number = 100, contractName?: string): Promise<EventRecord[]> {
    if (!this.enabled) return [];

    try {
      let query = 'SELECT * FROM events';
      const params: any[] = [];

      if (contractName) {
        query += ' WHERE contract_name = $1';
        params.push(contractName);
      }

      query += ' ORDER BY timestamp DESC LIMIT $' + (params.length + 1);
      params.push(limit);

      const result = await this.pool.query(query, params);

      return result.rows.map(row => ({
        id: row.id,
        eventType: row.event_type,
        contractName: row.contract_name,
        contractAddress: row.contract_address,
        blockNumber: row.block_number,
        transactionHash: row.transaction_hash,
        eventData: row.event_data,
        timestamp: row.timestamp,
        processed: row.processed,
      }));
    } catch (error: any) {
      console.error('[DB] Failed to get recent events:', error.message);
      return [];
    }
  }

  /**
   * 获取未确认的告警
   */
  async getUnacknowledgedAlerts(limit: number = 50): Promise<AlertRecord[]> {
    if (!this.enabled) return [];

    try {
      const result = await this.pool.query(
        `SELECT * FROM alerts
         WHERE acknowledged = FALSE
         ORDER BY timestamp DESC
         LIMIT $1`,
        [limit]
      );

      return result.rows.map(row => ({
        id: row.id,
        level: row.level,
        type: row.type,
        message: row.message,
        data: row.data,
        timestamp: row.timestamp,
        acknowledged: row.acknowledged,
      }));
    } catch (error: any) {
      console.error('[DB] Failed to get unacknowledged alerts:', error.message);
      return [];
    }
  }

  /**
   * 确认告警
   */
  async acknowledgeAlert(alertId: number, acknowledgedBy: string): Promise<boolean> {
    if (!this.enabled) return false;

    try {
      await this.pool.query(
        `UPDATE alerts
         SET acknowledged = TRUE,
             acknowledged_at = CURRENT_TIMESTAMP,
             acknowledged_by = $2
         WHERE id = $1`,
        [alertId, acknowledgedBy]
      );

      return true;
    } catch (error: any) {
      console.error('[DB] Failed to acknowledge alert:', error.message);
      return false;
    }
  }

  /**
   * 保存统计数据
   */
  async saveStatistic(
    contractName: string,
    metricName: string,
    metricValue: number,
    metadata?: any
  ): Promise<boolean> {
    if (!this.enabled) return false;

    try {
      await this.pool.query(
        `INSERT INTO statistics (contract_name, metric_name, metric_value, metadata)
         VALUES ($1, $2, $3, $4)`,
        [contractName, metricName, metricValue, metadata ? JSON.stringify(metadata) : null]
      );

      return true;
    } catch (error: any) {
      console.error('[DB] Failed to save statistic:', error.message);
      return false;
    }
  }

  /**
   * 获取统计数据
   */
  async getStatistics(
    contractName: string,
    metricName: string,
    hours: number = 24
  ): Promise<any[]> {
    if (!this.enabled) return [];

    try {
      const result = await this.pool.query(
        `SELECT metric_value, metadata, timestamp
         FROM statistics
         WHERE contract_name = $1
           AND metric_name = $2
           AND timestamp >= NOW() - INTERVAL '${hours} hours'
         ORDER BY timestamp DESC`,
        [contractName, metricName]
      );

      return result.rows;
    } catch (error: any) {
      console.error('[DB] Failed to get statistics:', error.message);
      return [];
    }
  }

  /**
   * 清理旧数据
   */
  async cleanupOldData(daysToKeep: number = 30): Promise<number> {
    if (!this.enabled) return 0;

    try {
      const eventsResult = await this.pool.query(
        `DELETE FROM events
         WHERE timestamp < NOW() - INTERVAL '${daysToKeep} days'`,
      );

      const alertsResult = await this.pool.query(
        `DELETE FROM alerts
         WHERE acknowledged = TRUE
           AND timestamp < NOW() - INTERVAL '${daysToKeep} days'`,
      );

      const statsResult = await this.pool.query(
        `DELETE FROM statistics
         WHERE timestamp < NOW() - INTERVAL '${daysToKeep} days'`,
      );

      const totalDeleted =
        (eventsResult.rowCount || 0) +
        (alertsResult.rowCount || 0) +
        (statsResult.rowCount || 0);

      console.log(`[DB] Cleaned up ${totalDeleted} old records`);
      return totalDeleted;
    } catch (error: any) {
      console.error('[DB] Failed to cleanup old data:', error.message);
      return 0;
    }
  }

  /**
   * 获取数据库统计
   */
  async getDatabaseStats(): Promise<any> {
    if (!this.enabled) {
      return { enabled: false };
    }

    try {
      const eventsCount = await this.pool.query('SELECT COUNT(*) FROM events');
      const alertsCount = await this.pool.query('SELECT COUNT(*) FROM alerts');
      const unackAlertsCount = await this.pool.query(
        'SELECT COUNT(*) FROM alerts WHERE acknowledged = FALSE'
      );
      const statsCount = await this.pool.query('SELECT COUNT(*) FROM statistics');

      return {
        enabled: true,
        totalEvents: parseInt(eventsCount.rows[0].count),
        totalAlerts: parseInt(alertsCount.rows[0].count),
        unacknowledgedAlerts: parseInt(unackAlertsCount.rows[0].count),
        totalStatistics: parseInt(statsCount.rows[0].count),
        poolSize: this.pool.totalCount,
        idleConnections: this.pool.idleCount,
        waitingRequests: this.pool.waitingCount,
      };
    } catch (error: any) {
      console.error('[DB] Failed to get database stats:', error.message);
      return { enabled: true, error: error.message };
    }
  }

  /**
   * 测试连接
   */
  async testConnection(): Promise<boolean> {
    if (!this.enabled) {
      console.log('[DB] Database disabled');
      return false;
    }

    try {
      const result = await this.pool.query('SELECT NOW()');
      console.log('[DB] Connection test successful:', result.rows[0].now);
      return true;
    } catch (error: any) {
      console.error('[DB] Connection test failed:', error.message);
      return false;
    }
  }

  /**
   * 关闭连接池
   */
  async close(): Promise<void> {
    if (this.enabled) {
      await this.pool.end();
      console.log('[DB] Connection pool closed');
    }
  }
}
