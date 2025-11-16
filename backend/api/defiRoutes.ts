/**
 * DeFi Position API Routes
 * Exposes institutional-grade portfolio tracking to frontend
 */

import { Router, Request, Response } from 'express';
import { positionTracker } from '../services/defi/positionTracker';
import { Pool } from 'pg';

const router = Router();
const db = new Pool({ connectionString: process.env.DATABASE_URL });

/**
 * GET /api/defi/positions/:address
 * Get all DeFi positions for a user
 */
router.get('/positions/:address', async (req: Request, res: Response) => {
  try {
    const { address } = req.params;

    // Track current positions
    const portfolio = await positionTracker.trackUserPositions(address);

    // Get historical data
    const history = await getPositionHistory(address);

    res.json({
      success: true,
      data: {
        current: portfolio,
        history: history,
        lastUpdated: new Date(),
      },
    });
  } catch (error) {
    console.error('Error fetching positions:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to fetch positions',
    });
  }
});

/**
 * GET /api/defi/risk/:address
 * Get risk assessment for user portfolio
 */
router.get('/risk/:address', async (req: Request, res: Response) => {
  try {
    const { address } = req.params;

    // Get latest portfolio snapshot
    const result = await db.query(
      `SELECT
        overall_risk,
        average_health_factor,
        liquidation_risk,
        recommendations,
        ai_analysis,
        timestamp
       FROM portfolio_snapshots
       WHERE user_id = $1
       ORDER BY timestamp DESC
       LIMIT 1`,
      [address]
    );

    if (result.rows.length === 0) {
      // No cached data, fetch fresh
      const portfolio = await positionTracker.trackUserPositions(address);
      return res.json({
        success: true,
        data: {
          overallRisk: portfolio.overallRisk,
          healthFactor: portfolio.averageHealthFactor,
          recommendations: portfolio.recommendations,
          positions: portfolio.positions.map(p => ({
            protocol: p.protocol,
            risk: p.risk,
            healthFactor: p.healthFactor,
          })),
        },
      });
    }

    const risk = result.rows[0];
    res.json({
      success: true,
      data: {
        overallRisk: risk.overall_risk,
        healthFactor: risk.average_health_factor,
        liquidationRisk: risk.liquidation_risk,
        recommendations: risk.recommendations,
        aiAnalysis: risk.ai_analysis,
        lastUpdated: risk.timestamp,
      },
    });
  } catch (error) {
    console.error('Error fetching risk:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to fetch risk assessment',
    });
  }
});

/**
 * GET /api/defi/alerts/:address
 * Get active liquidation alerts
 */
router.get('/alerts/:address', async (req: Request, res: Response) => {
  try {
    const { address } = req.params;

    const result = await db.query(
      `SELECT
        id,
        position_id,
        alert_type,
        health_factor,
        risk_score,
        predicted_liquidation_time,
        confidence_score,
        recommended_action,
        required_collateral,
        created_at
       FROM liquidation_alerts
       WHERE user_id = $1 AND status = 'active'
       ORDER BY
         CASE alert_type
           WHEN 'critical' THEN 1
           WHEN 'warning' THEN 2
           ELSE 3
         END,
         created_at DESC`,
      [address]
    );

    res.json({
      success: true,
      data: result.rows,
    });
  } catch (error) {
    console.error('Error fetching alerts:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to fetch alerts',
    });
  }
});

/**
 * POST /api/defi/hedge
 * Initiate automated hedging
 */
router.post('/hedge', async (req: Request, res: Response) => {
  try {
    const { userId, positionId, hedgeType, amount } = req.body;

    // Validate request
    if (!userId || !positionId || !hedgeType) {
      return res.status(400).json({
        success: false,
        error: 'Missing required parameters',
      });
    }

    // Create hedging transaction record
    const result = await db.query(
      `INSERT INTO hedging_transactions
       (user_id, position_id, hedge_type, amount, status)
       VALUES ($1, $2, $3, $4, 'pending')
       RETURNING id`,
      [userId, positionId, hedgeType, amount || 0]
    );

    const hedgeId = result.rows[0].id;

    // TODO: Execute hedging transaction on-chain
    // This would interact with smart contracts

    res.json({
      success: true,
      data: {
        hedgeId,
        status: 'pending',
        message: 'Hedging transaction initiated',
      },
    });
  } catch (error) {
    console.error('Error initiating hedge:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to initiate hedging',
    });
  }
});

/**
 * GET /api/defi/history/:address
 * Get portfolio history for charts
 */
router.get('/history/:address', async (req: Request, res: Response) => {
  try {
    const { address } = req.params;
    const { period = '7d' } = req.query;

    let interval = '1 hour';
    let limit = 168; // 7 days * 24 hours

    switch (period) {
      case '24h':
        interval = '15 minutes';
        limit = 96;
        break;
      case '7d':
        interval = '1 hour';
        limit = 168;
        break;
      case '30d':
        interval = '4 hours';
        limit = 180;
        break;
      case '90d':
        interval = '1 day';
        limit = 90;
        break;
    }

    const result = await db.query(
      `SELECT
        date_trunc($1, timestamp) as time,
        AVG(total_value) as value,
        AVG(overall_risk) as risk,
        MIN(average_health_factor) as min_health_factor
       FROM portfolio_snapshots
       WHERE user_id = $2 AND timestamp > NOW() - INTERVAL $3
       GROUP BY time
       ORDER BY time DESC
       LIMIT $4`,
      [interval, address, period, limit]
    );

    res.json({
      success: true,
      data: result.rows,
    });
  } catch (error) {
    console.error('Error fetching history:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to fetch history',
    });
  }
});

/**
 * GET /api/defi/protocols
 * Get protocol-level statistics
 */
router.get('/protocols', async (_req: Request, res: Response) => {
  try {
    const result = await db.query(
      `SELECT
        protocol,
        total_value_locked,
        user_count,
        utilization_rate,
        protocol_risk_score,
        timestamp
       FROM protocol_tvl
       WHERE timestamp = (
         SELECT MAX(timestamp) FROM protocol_tvl
       )
       ORDER BY total_value_locked DESC`
    );

    res.json({
      success: true,
      data: result.rows,
    });
  } catch (error) {
    console.error('Error fetching protocols:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to fetch protocol stats',
    });
  }
});

/**
 * WebSocket endpoint for real-time updates
 * WS /api/defi/stream/:address
 */
router.ws('/stream/:address', (ws, req) => {
  const { address } = req.params;

  // Subscribe to position updates
  const interval = setInterval(async () => {
    try {
      // Get latest risk metrics
      const result = await db.query(
        `SELECT
          overall_risk,
          average_health_factor,
          total_value
         FROM portfolio_snapshots
         WHERE user_id = $1
         ORDER BY timestamp DESC
         LIMIT 1`,
        [address]
      );

      if (result.rows.length > 0) {
        ws.send(JSON.stringify({
          type: 'risk_update',
          data: result.rows[0],
          timestamp: new Date(),
        }));
      }

      // Check for new alerts
      const alerts = await db.query(
        `SELECT * FROM liquidation_alerts
         WHERE user_id = $1
         AND status = 'active'
         AND created_at > NOW() - INTERVAL '1 minute'`,
        [address]
      );

      if (alerts.rows.length > 0) {
        ws.send(JSON.stringify({
          type: 'new_alerts',
          data: alerts.rows,
          timestamp: new Date(),
        }));
      }
    } catch (error) {
      console.error('WebSocket update error:', error);
    }
  }, 5000); // Update every 5 seconds

  ws.on('close', () => {
    clearInterval(interval);
  });
});

// Helper functions
async function getPositionHistory(address: string) {
  const result = await db.query(
    `SELECT * FROM defi_positions
     WHERE user_id = $1
     ORDER BY last_updated DESC
     LIMIT 100`,
    [address]
  );
  return result.rows;
}

export default router;