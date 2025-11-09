-- GMX V2 数据库表结构
-- 创建日期: 2025-11-09

-- GMX 仓位表
CREATE TABLE IF NOT EXISTS gmx_positions (
  id SERIAL PRIMARY KEY,
  user_address VARCHAR(42) NOT NULL,
  order_key VARCHAR(66) NOT NULL UNIQUE,
  market VARCHAR(42) NOT NULL,
  collateral_token VARCHAR(42) NOT NULL,
  is_long BOOLEAN NOT NULL,
  size_usd VARCHAR(78) NOT NULL,
  collateral_amount VARCHAR(78) NOT NULL,
  leverage VARCHAR(10) NOT NULL,
  is_hedge BOOLEAN DEFAULT false,
  status VARCHAR(20) DEFAULT 'open',
  closed_pnl VARCHAR(78),
  closed_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_gmx_positions_user ON gmx_positions(user_address);
CREATE INDEX IF NOT EXISTS idx_gmx_positions_status ON gmx_positions(status);
CREATE INDEX IF NOT EXISTS idx_gmx_positions_market ON gmx_positions(market);
CREATE INDEX IF NOT EXISTS idx_gmx_positions_created ON gmx_positions(created_at DESC);

-- GMX 对冲记录表
CREATE TABLE IF NOT EXISTS gmx_hedge_records (
  id SERIAL PRIMARY KEY,
  user_address VARCHAR(42) NOT NULL,
  market VARCHAR(42) NOT NULL,
  hedge_size VARCHAR(78) NOT NULL,
  reason TEXT NOT NULL,
  order_key VARCHAR(66) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_gmx_hedge_records_user ON gmx_hedge_records(user_address);
CREATE INDEX IF NOT EXISTS idx_gmx_hedge_records_created ON gmx_hedge_records(created_at DESC);

-- GMX 风险建议记录表 (用于追踪建议效果)
CREATE TABLE IF NOT EXISTS gmx_risk_recommendations (
  id SERIAL PRIMARY KEY,
  user_address VARCHAR(42) NOT NULL,
  position_id INTEGER REFERENCES gmx_positions(id),
  recommendation_type VARCHAR(50) NOT NULL,
  level VARCHAR(20) NOT NULL,
  action VARCHAR(100) NOT NULL,
  priority VARCHAR(20) NOT NULL,
  reason TEXT NOT NULL,
  expected_outcome TEXT,
  user_decision BOOLEAN NOT NULL,
  accepted BOOLEAN,
  result TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  responded_at TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_gmx_recommendations_user ON gmx_risk_recommendations(user_address);
CREATE INDEX IF NOT EXISTS idx_gmx_recommendations_position ON gmx_risk_recommendations(position_id);
CREATE INDEX IF NOT EXISTS idx_gmx_recommendations_created ON gmx_risk_recommendations(created_at DESC);

-- GMX 统计视图
CREATE OR REPLACE VIEW gmx_position_stats AS
SELECT
  COUNT(*) as total_positions,
  COUNT(*) FILTER (WHERE status = 'open') as open_positions,
  COUNT(*) FILTER (WHERE status = 'closed') as closed_positions,
  COUNT(*) FILTER (WHERE is_hedge = true) as hedge_positions,
  COUNT(*) FILTER (WHERE is_long = true) as long_positions,
  COUNT(*) FILTER (WHERE is_long = false) as short_positions,
  AVG(CAST(leverage AS NUMERIC)) as avg_leverage,
  COUNT(*) FILTER (WHERE CAST(leverage AS NUMERIC) >= 30) as high_leverage_count
FROM gmx_positions;

-- 注释
COMMENT ON TABLE gmx_positions IS 'GMX V2 仓位记录';
COMMENT ON TABLE gmx_hedge_records IS 'GMX V2 紧急对冲记录';
COMMENT ON TABLE gmx_risk_recommendations IS 'GMX 风险建议追踪表';
COMMENT ON VIEW gmx_position_stats IS 'GMX 仓位统计视图';
