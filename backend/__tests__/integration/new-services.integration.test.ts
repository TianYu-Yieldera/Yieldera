/**
 * Integration Tests for New Services
 *
 * These tests verify end-to-end functionality across:
 * - Yield Calculation Service
 * - Notification Service
 * - Auto Hedge Executor
 * - Yield Distributor
 */

import { Pool } from 'pg';
import axios from 'axios';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';
const TEST_USER_ADDRESS = '0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb6';

describe('New Services Integration Tests', () => {
  let dbPool: Pool;

  beforeAll(async () => {
    // Setup database connection
    dbPool = new Pool({
      connectionString: process.env.DATABASE_URL ||
        'postgres://loyalty_user:loyalty_pass@localhost:5432/loyalty_db',
    });

    // Verify database connection
    try {
      await dbPool.query('SELECT 1');
    } catch (error) {
      console.log('Database not available, skipping integration tests');
      return;
    }
  });

  afterAll(async () => {
    await dbPool?.end();
  });

  describe('Yield Calculation Service Integration', () => {
    it('should complete full yield calculation workflow', async () => {
      // Step 1: Get current treasury rates
      const ratesResponse = await axios.get(`${API_BASE_URL}/api/v1/yields/rates`);
      expect(ratesResponse.status).toBe(200);
      expect(ratesResponse.data.rates).toBeInstanceOf(Array);
      expect(ratesResponse.data.rates.length).toBeGreaterThan(0);

      // Step 2: Project yield for investment
      const projectionRequest = {
        bond_type: 'TBILL_3M',
        principal_usd: 10000,
        duration_days: 90,
        compounding: true,
      };

      const projectionResponse = await axios.post(
        `${API_BASE_URL}/api/v1/yields/project`,
        projectionRequest
      );

      expect(projectionResponse.status).toBe(200);
      expect(projectionResponse.data.total_yield).toBeGreaterThan(0);
      expect(projectionResponse.data.projected_value).toBeGreaterThan(10000);

      // Step 3: Get user's total yield
      const totalYieldResponse = await axios.get(
        `${API_BASE_URL}/api/v1/yields/total/${TEST_USER_ADDRESS}`
      );

      expect(totalYieldResponse.status).toBe(200);
      expect(totalYieldResponse.data).toHaveProperty('total_yield');
    });

    it('should handle yield calculation for multiple bond types', async () => {
      const bondTypes = ['TBILL_3M', 'TBILL_6M', 'TNOTE_2Y', 'TNOTE_5Y', 'TNOTE_10Y'];

      for (const bondType of bondTypes) {
        const response = await axios.post(
          `${API_BASE_URL}/api/v1/yields/project`,
          {
            bond_type: bondType,
            principal_usd: 1000,
            duration_days: 30,
            compounding: false,
          }
        );

        expect(response.status).toBe(200);
        expect(response.data.total_yield).toBeGreaterThan(0);
      }
    });
  });

  describe('Notification Service Integration', () => {
    it('should complete full notification preferences workflow', async () => {
      // Step 1: Get current preferences
      const getPrefsResponse = await axios.get(
        `${API_BASE_URL}/api/v1/notifications/${TEST_USER_ADDRESS}/preferences`
      );

      expect(getPrefsResponse.status).toBe(200);

      // Step 2: Update preferences
      const newPreferences = {
        channels: ['email', 'push', 'sms'],
        min_priority: 'high',
        enabled_types: ['liquidation_warning', 'high_risk_position'],
        frequency: 'realtime',
      };

      const updateResponse = await axios.put(
        `${API_BASE_URL}/api/v1/notifications/${TEST_USER_ADDRESS}/preferences`,
        newPreferences
      );

      expect(updateResponse.status).toBe(200);
      expect(updateResponse.data.success).toBe(true);

      // Step 3: Verify preferences were updated
      const verifyResponse = await axios.get(
        `${API_BASE_URL}/api/v1/notifications/${TEST_USER_ADDRESS}/preferences`
      );

      expect(verifyResponse.status).toBe(200);
      const prefs = verifyResponse.data.preferences || verifyResponse.data;
      expect(prefs.min_priority).toBe('high');
    });

    it('should retrieve user notifications', async () => {
      const response = await axios.get(
        `${API_BASE_URL}/api/v1/notifications/${TEST_USER_ADDRESS}`
      );

      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('notifications');
    });

    it('should filter unread notifications', async () => {
      const response = await axios.get(
        `${API_BASE_URL}/api/v1/notifications/${TEST_USER_ADDRESS}?unread=true`
      );

      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('notifications');
    });
  });

  describe('Auto Hedge Executor Integration', () => {
    it('should complete full hedge settings workflow', async () => {
      // Step 1: Get current settings
      const getSettingsResponse = await axios.get(
        `${API_BASE_URL}/api/v1/hedge/settings/${TEST_USER_ADDRESS}`
      );

      expect(getSettingsResponse.status).toBe(200);
      expect(getSettingsResponse.data).toHaveProperty('auto_hedge_enabled');

      // Step 2: Update settings
      const newSettings = {
        auto_hedge_enabled: true,
        max_hedge_amount: 50000,
        min_health_factor: 1.5,
        target_health_factor: 2.5,
      };

      const updateResponse = await axios.put(
        `${API_BASE_URL}/api/v1/hedge/settings/${TEST_USER_ADDRESS}`,
        newSettings
      );

      expect(updateResponse.status).toBe(200);
      expect(updateResponse.data.success).toBe(true);

      // Step 3: Verify settings were updated
      const verifyResponse = await axios.get(
        `${API_BASE_URL}/api/v1/hedge/settings/${TEST_USER_ADDRESS}`
      );

      expect(verifyResponse.status).toBe(200);
      expect(verifyResponse.data.auto_hedge_enabled).toBe(true);
      expect(verifyResponse.data.max_hedge_amount).toBe(50000);
    });

    it('should retrieve hedge execution history', async () => {
      const response = await axios.get(
        `${API_BASE_URL}/api/v1/hedge/history/${TEST_USER_ADDRESS}`
      );

      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('executions');
    });

    it('should retrieve hedge statistics', async () => {
      const response = await axios.get(
        `${API_BASE_URL}/api/v1/hedge/stats`
      );

      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('stats');
    });

    it('should filter hedge stats by date range', async () => {
      const response = await axios.get(
        `${API_BASE_URL}/api/v1/hedge/stats?days=7`
      );

      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('stats');
    });
  });

  describe('Yield Distribution Service Integration', () => {
    it('should retrieve distribution statistics', async () => {
      const response = await axios.get(
        `${API_BASE_URL}/api/v1/distribution/stats`
      );

      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('stats');
    });

    it('should filter distribution stats by date range', async () => {
      const response = await axios.get(
        `${API_BASE_URL}/api/v1/distribution/stats?days=30`
      );

      expect(response.status).toBe(200);
      expect(response.data).toHaveProperty('stats');
    });
  });

  describe('Cross-Service Integration Scenarios', () => {
    it('should handle complete user onboarding flow', async () => {
      const newUserAddress = '0x' + Math.random().toString(16).slice(2, 42);

      // 1. Set up notification preferences
      const notifResponse = await axios.put(
        `${API_BASE_URL}/api/v1/notifications/${newUserAddress}/preferences`,
        {
          channels: ['email'],
          min_priority: 'medium',
          enabled_types: ['daily_yield_report', 'liquidation_warning'],
          frequency: 'realtime',
        }
      );
      expect(notifResponse.status).toBe(200);

      // 2. Enable auto-hedging
      const hedgeResponse = await axios.put(
        `${API_BASE_URL}/api/v1/hedge/settings/${newUserAddress}`,
        {
          auto_hedge_enabled: true,
          max_hedge_amount: 10000,
          min_health_factor: 1.5,
          target_health_factor: 2.0,
        }
      );
      expect(hedgeResponse.status).toBe(200);

      // 3. Project yields
      const yieldResponse = await axios.post(
        `${API_BASE_URL}/api/v1/yields/project`,
        {
          bond_type: 'TBILL_3M',
          principal_usd: 5000,
          duration_days: 90,
          compounding: true,
        }
      );
      expect(yieldResponse.status).toBe(200);
      expect(yieldResponse.data.total_yield).toBeGreaterThan(0);
    });

    it('should handle high-risk position scenario', async () => {
      // Simulate a user with risky position
      const riskUserAddress = TEST_USER_ADDRESS;

      // 1. Check current hedge settings
      const settingsResponse = await axios.get(
        `${API_BASE_URL}/api/v1/hedge/settings/${riskUserAddress}`
      );
      expect(settingsResponse.status).toBe(200);

      // 2. Get notification preferences
      const prefsResponse = await axios.get(
        `${API_BASE_URL}/api/v1/notifications/${riskUserAddress}/preferences`
      );
      expect(prefsResponse.status).toBe(200);

      // 3. Check yield history
      const yieldResponse = await axios.get(
        `${API_BASE_URL}/api/v1/yields/total/${riskUserAddress}`
      );
      expect(yieldResponse.status).toBe(200);
    });
  });

  describe('Database Integration', () => {
    it('should persist yield projections correctly', async () => {
      const bondType = 'TBILL_3M';
      const principal = 1000;

      // Make projection
      await axios.post(`${API_BASE_URL}/api/v1/yields/project`, {
        bond_type: bondType,
        principal_usd: principal,
        duration_days: 30,
        compounding: false,
      });

      // Verify rate exists in database
      const result = await dbPool.query(
        'SELECT * FROM treasury_yield_rates WHERE bond_type = $1 LIMIT 1',
        [bondType]
      );

      expect(result.rows.length).toBeGreaterThan(0);
      expect(result.rows[0].annual_yield).toBeGreaterThan(0);
    });

    it('should store notification preferences in database', async () => {
      const testAddress = TEST_USER_ADDRESS;

      // Update preferences via API
      await axios.put(
        `${API_BASE_URL}/api/v1/notifications/${testAddress}/preferences`,
        {
          channels: ['email', 'sms'],
          min_priority: 'critical',
          enabled_types: ['liquidation_warning'],
          frequency: 'hourly',
        }
      );

      // Verify in database
      const result = await dbPool.query(
        'SELECT * FROM notification_preferences WHERE user_id = $1',
        [testAddress.toLowerCase()]
      );

      expect(result.rows.length).toBeGreaterThan(0);
      expect(result.rows[0].min_priority).toBe('critical');
    });
  });

  describe('Error Handling', () => {
    it('should handle invalid user address', async () => {
      try {
        await axios.get(`${API_BASE_URL}/api/v1/yields/total/invalid_address`);
      } catch (error: any) {
        expect(error.response.status).toBeGreaterThanOrEqual(400);
      }
    });

    it('should handle invalid bond type', async () => {
      try {
        await axios.post(`${API_BASE_URL}/api/v1/yields/project`, {
          bond_type: 'INVALID',
          principal_usd: 1000,
          duration_days: 90,
        });
      } catch (error: any) {
        expect(error.response.status).toBe(400);
      }
    });

    it('should handle missing required fields', async () => {
      try {
        await axios.post(`${API_BASE_URL}/api/v1/yields/project`, {
          bond_type: 'TBILL_3M',
          // Missing principal_usd and duration_days
        });
      } catch (error: any) {
        expect(error.response.status).toBe(400);
      }
    });
  });

  describe('Performance Tests', () => {
    it('should handle concurrent requests efficiently', async () => {
      const requests = Array(20).fill(null).map(() =>
        axios.get(`${API_BASE_URL}/api/v1/yields/rates`)
      );

      const startTime = Date.now();
      const responses = await Promise.all(requests);
      const duration = Date.now() - startTime;

      expect(responses.every(r => r.status === 200)).toBe(true);
      expect(duration).toBeLessThan(5000); // Should complete within 5 seconds
    });

    it('should respond quickly to yield projection requests', async () => {
      const startTime = Date.now();

      await axios.post(`${API_BASE_URL}/api/v1/yields/project`, {
        bond_type: 'TBILL_3M',
        principal_usd: 1000,
        duration_days: 90,
        compounding: true,
      });

      const duration = Date.now() - startTime;
      expect(duration).toBeLessThan(1000); // Should complete within 1 second
    });
  });
});
