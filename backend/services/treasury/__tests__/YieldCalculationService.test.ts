import { YieldCalculationService } from '../YieldCalculationService';
import { Pool } from 'pg';

// Mock pg Pool
jest.mock('pg', () => {
  const mClient = {
    query: jest.fn(),
    release: jest.fn(),
  };
  const mPool = {
    connect: jest.fn(() => Promise.resolve(mClient)),
    query: jest.fn(),
    end: jest.fn(),
  };
  return { Pool: jest.fn(() => mPool) };
});

describe('YieldCalculationService', () => {
  let service: YieldCalculationService;
  let mockPool: jest.Mocked<Pool>;

  beforeEach(() => {
    mockPool = new Pool() as jest.Mocked<Pool>;
    service = new YieldCalculationService(mockPool);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('calculateDailyYield', () => {
    it('should calculate daily yield correctly for simple interest', () => {
      const principal = 10000;
      const annualRate = 0.045; // 4.5%
      const days = 1;

      const result = service['calculateDailyYield'](principal, annualRate, days, false);

      expect(result).toBeCloseTo(1.23, 2); // $10000 * 0.045 / 365 ≈ $1.23
    });

    it('should calculate daily yield correctly for compound interest', () => {
      const principal = 10000;
      const annualRate = 0.045;
      const days = 90;

      const resultSimple = service['calculateDailyYield'](principal, annualRate, days, false);
      const resultCompound = service['calculateDailyYield'](principal, annualRate, days, true);

      expect(resultCompound).toBeGreaterThan(resultSimple);
    });

    it('should handle zero principal', () => {
      const result = service['calculateDailyYield'](0, 0.045, 1, false);
      expect(result).toBe(0);
    });

    it('should handle zero interest rate', () => {
      const result = service['calculateDailyYield'](10000, 0, 1, false);
      expect(result).toBe(0);
    });
  });

  describe('projectYield', () => {
    beforeEach(() => {
      // Mock database query for treasury rates
      (mockPool.query as jest.Mock).mockResolvedValue({
        rows: [{ annual_yield: 0.045 }],
      });
    });

    it('should project yield correctly', async () => {
      const result = await service.projectYield('TBILL_3M', 1000, 90, true);

      expect(result).toHaveProperty('totalYield');
      expect(result).toHaveProperty('effectiveAPY');
      expect(result).toHaveProperty('dailyYield');
      expect(result).toHaveProperty('projectedValue');
      expect(result.totalYield).toBeGreaterThan(0);
      expect(result.projectedValue).toBeGreaterThan(1000);
    });

    it('should throw error for invalid bond type', async () => {
      (mockPool.query as jest.Mock).mockResolvedValue({ rows: [] });

      await expect(
        service.projectYield('INVALID_TYPE', 1000, 90, true)
      ).rejects.toThrow();
    });

    it('should handle negative principal', async () => {
      await expect(
        service.projectYield('TBILL_3M', -1000, 90, true)
      ).rejects.toThrow();
    });

    it('should handle zero duration', async () => {
      await expect(
        service.projectYield('TBILL_3M', 1000, 0, true)
      ).rejects.toThrow();
    });
  });

  describe('getUserTotalYield', () => {
    it('should return total yield by bond type', async () => {
      const mockData = [
        { bond_type: 'TBILL_3M', total_yield: 100.50 },
        { bond_type: 'TBILL_6M', total_yield: 200.75 },
      ];

      (mockPool.query as jest.Mock).mockResolvedValue({ rows: mockData });

      const result = await service.getUserTotalYield('0x123...');

      expect(result.totalYield).toBe(301.25);
      expect(result.yieldByType).toHaveProperty('TBILL_3M', 100.50);
      expect(result.yieldByType).toHaveProperty('TBILL_6M', 200.75);
    });

    it('should handle user with no yields', async () => {
      (mockPool.query as jest.Mock).mockResolvedValue({ rows: [] });

      const result = await service.getUserTotalYield('0x123...');

      expect(result.totalYield).toBe(0);
      expect(Object.keys(result.yieldByType).length).toBe(0);
    });
  });

  describe('getTreasuryRates', () => {
    it('should fetch all current treasury rates', async () => {
      const mockRates = [
        {
          bond_type: 'TBILL_3M',
          annual_yield: 0.045,
          effective_date: new Date(),
          source: 'Initial',
        },
        {
          bond_type: 'TBILL_6M',
          annual_yield: 0.047,
          effective_date: new Date(),
          source: 'Initial',
        },
      ];

      (mockPool.query as jest.Mock).mockResolvedValue({ rows: mockRates });

      const result = await service.getTreasuryRates();

      expect(result).toHaveLength(2);
      expect(result[0].bondType).toBe('TBILL_3M');
      expect(result[0].annualYield).toBe(0.045);
    });
  });

  describe('updateYieldRates', () => {
    it('should update rates from US Treasury API', async () => {
      // Mock successful fetch
      global.fetch = jest.fn(() =>
        Promise.resolve({
          ok: true,
          json: () =>
            Promise.resolve({
              data: [
                { rate: '4.50', maturity: '3-Month' },
                { rate: '4.70', maturity: '6-Month' },
              ],
            }),
        })
      ) as jest.Mock;

      (mockPool.query as jest.Mock).mockResolvedValue({ rows: [] });

      await service['updateYieldRates']();

      expect(mockPool.query).toHaveBeenCalled();
    });

    it('should handle API fetch failure gracefully', async () => {
      global.fetch = jest.fn(() =>
        Promise.resolve({
          ok: false,
          status: 500,
        })
      ) as jest.Mock;

      // Should not throw error, just log it
      await expect(service['updateYieldRates']()).resolves.not.toThrow();
    });
  });

  describe('Performance', () => {
    it('should calculate yields for multiple users efficiently', async () => {
      const mockHoldings = Array(100).fill({
        user_id: '0x123...',
        bond_type: 'TBILL_3M',
        principal_amount: 1000,
        last_update: new Date(Date.now() - 86400000), // 1 day ago
      });

      (mockPool.query as jest.Mock)
        .mockResolvedValueOnce({ rows: mockHoldings })
        .mockResolvedValue({ rows: [{ annual_yield: 0.045 }] });

      const startTime = Date.now();
      await service.calculateDailyYields();
      const duration = Date.now() - startTime;

      // Should complete within 5 seconds for 100 users
      expect(duration).toBeLessThan(5000);
    });
  });
});
