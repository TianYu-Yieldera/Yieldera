/**
 * FastAPI Risk Service Client
 * 直接连接到FastAPI Python服务 (http://localhost:8000)
 * 提供VaR/CVaR、Agent模拟、ML预测等高级风险功能
 */

const FASTAPI_BASE = import.meta.env.VITE_FASTAPI_URL || 'http://localhost:8000';

class FastAPIRiskService {
  constructor() {
    this.baseURL = FASTAPI_BASE;
    this.healthCheckCache = null;
    this.cacheExpiry = 0;
  }

  /**
   * 检查FastAPI服务健康状态 (带缓存)
   */
  async checkHealth() {
    const now = Date.now();

    // 使用5秒缓存
    if (this.healthCheckCache && now < this.cacheExpiry) {
      return this.healthCheckCache;
    }

    try {
      const response = await fetch(`${this.baseURL}/health`, {
        signal: AbortSignal.timeout(3000), // 3秒超时
      });

      if (!response.ok) {
        throw new Error(`Health check failed: ${response.status}`);
      }

      const data = await response.json();

      this.healthCheckCache = {
        available: data.status === 'healthy',
        modules: data.modules_loaded || [],
        version: data.version,
        pythonVersion: data.python_version,
      };

      this.cacheExpiry = now + 5000; // 5秒缓存

      return this.healthCheckCache;
    } catch (error) {
      console.warn('[FastAPI] Service unavailable:', error.message);

      // 返回离线状态
      this.healthCheckCache = {
        available: false,
        modules: [],
        error: error.message,
      };

      this.cacheExpiry = now + 5000;

      return this.healthCheckCache;
    }
  }

  /**
   * 计算高级风险指标 (VaR/CVaR/Sharpe)
   * @param {Object} position - 仓位数据
   * @param {Object} marketData - 市场数据
   * @param {Array} historicalPrices - 历史价格数据
   */
  async calculateRisk(position, marketData, historicalPrices = []) {
    try {
      const response = await fetch(`${this.baseURL}/api/calculate_risk`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          position: this.serializePosition(position),
          market_data: this.serializeMarketData(marketData),
          historical_data: historicalPrices,
        }),
        signal: AbortSignal.timeout(30000), // 30秒超时
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.detail || `HTTP ${response.status}`);
      }

      const data = await response.json();

      if (!data.success) {
        throw new Error(data.error || 'Risk calculation failed');
      }

      return data.risk_metrics;
    } catch (error) {
      console.error('[FastAPI] Risk calculation error:', error);
      throw error;
    }
  }

  /**
   * 运行Agent模拟回测
   * @param {Object} params - 模拟参数
   */
  async runAgentSimulation(params) {
    try {
      const response = await fetch(`${this.baseURL}/api/agent_simulation`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(params),
        signal: AbortSignal.timeout(60000), // 60秒超时
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.detail || `HTTP ${response.status}`);
      }

      const data = await response.json();

      if (!data.success) {
        throw new Error(data.error || 'Simulation failed');
      }

      return {
        scenario: data.scenario,
        timestamp: new Date(data.timestamp),
        summary: data.summary,
        agentStats: data.agent_stats,
        timeSeries: data.time_series,
      };
    } catch (error) {
      console.error('[FastAPI] Agent simulation error:', error);
      throw error;
    }
  }

  /**
   * ML风险预测 (24小时)
   * @param {Object} position - 仓位数据
   * @param {Object} marketData - 市场数据
   * @param {number} timeHorizon - 预测时间范围 (小时)
   */
  async predictRisk(position, marketData, timeHorizon = 24) {
    try {
      const response = await fetch(`${this.baseURL}/api/predict_risk`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          position: this.serializePosition(position),
          market_data: this.serializeMarketData(marketData),
          time_horizon: timeHorizon,
        }),
        signal: AbortSignal.timeout(30000),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.detail || `HTTP ${response.status}`);
      }

      const data = await response.json();

      if (!data.success) {
        throw new Error(data.error || 'Prediction failed');
      }

      return data.prediction || data;
    } catch (error) {
      console.error('[FastAPI] Risk prediction error:', error);
      throw error;
    }
  }

  /**
   * 压力测试
   * @param {Array} positions - 仓位列表
   * @param {Array} scenarios - 压力测试场景
   */
  async runStressTest(positions, scenarios) {
    try {
      const response = await fetch(`${this.baseURL}/api/stress_test`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          positions: positions.map(p => this.serializePosition(p)),
          scenarios,
        }),
        signal: AbortSignal.timeout(45000),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.detail || `HTTP ${response.status}`);
      }

      const data = await response.json();

      if (!data.success) {
        throw new Error(data.error || 'Stress test failed');
      }

      return data.stress_test_results;
    } catch (error) {
      console.error('[FastAPI] Stress test error:', error);
      throw error;
    }
  }

  /**
   * 序列化仓位数据为Python格式
   */
  serializePosition(position) {
    return {
      position_id: position.id || position.position_id || 'unknown',
      protocol: (position.protocol || 'aave').toLowerCase(),
      chain: (position.chain || 'arbitrum').toLowerCase(),
      collateral_asset: position.collateralAsset || position.collateral_asset || 'ETH',
      collateral_value_usd: position.collateralValueUSD || position.collateral_value_usd || 0,
      debt_value_usd: position.borrowValueUSD || position.debt_value_usd || 0,
      health_factor: position.healthFactor || position.health_factor || 2.0,
      leverage: position.leverage || 1.0,
      ltv: position.ltv || (position.borrowValueUSD && position.collateralValueUSD
        ? position.borrowValueUSD / position.collateralValueUSD
        : 0),
      current_price: position.currentPrice || position.current_price || 0,
      liquidation_price: position.liquidationPrice || position.liquidation_price || 0,
      position_age_days: position.position_age_days || 0,
    };
  }

  /**
   * 序列化市场数据为Python格式
   */
  serializeMarketData(marketData) {
    return {
      asset: marketData.asset || 'ETH',
      price: marketData.price || 0,
      price_change_24h: marketData.priceChange24h || marketData.price_change_24h || 0,
      volume_24h: marketData.volume24h || marketData.volume_24h || 0,
      volatility: marketData.volatility || 0.1,
      liquidity: marketData.liquidity || 0,
      timestamp: marketData.timestamp || new Date().toISOString(),
    };
  }

  /**
   * 批量计算风险 (用于仪表板)
   * @param {Array} positions - 仓位列表
   * @param {Object} marketData - 市场数据映射
   */
  async batchCalculateRisk(positions, marketData) {
    const results = [];

    for (const position of positions) {
      try {
        const asset = position.collateralAsset || 'ETH';
        const assetMarketData = marketData[asset] || {};

        const risk = await this.calculateRisk(position, assetMarketData);

        results.push({
          position_id: position.id,
          success: true,
          risk,
        });
      } catch (error) {
        results.push({
          position_id: position.id,
          success: false,
          error: error.message,
        });
      }
    }

    return results;
  }
}

// Export singleton instance
const fastAPIRiskService = new FastAPIRiskService();
export default fastAPIRiskService;
