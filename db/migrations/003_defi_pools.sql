-- Migration 003: DeFi Pools and Stablecoin Tables
-- Created: 2025-10-16
-- Purpose: Add tables for DeFi protocol pools and stablecoin functionality

-- 1. DeFi 协议池配置表
CREATE TABLE IF NOT EXISTS defi_pools (
    id VARCHAR(50) PRIMARY KEY,
    protocol VARCHAR(100) NOT NULL,
    name VARCHAR(200) NOT NULL,
    pool_type VARCHAR(50) NOT NULL,
    apr VARCHAR(20),
    tvl VARCHAR(50),
    risk_level VARCHAR(20),
    icon VARCHAR(10),
    color VARCHAR(20),
    points_required INTEGER DEFAULT 0,
    points_multiplier DECIMAL(3, 1) DEFAULT 1.0,
    features JSONB,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 2. 用户 DeFi 仓位表
CREATE TABLE IF NOT EXISTS user_defi_positions (
    id SERIAL PRIMARY KEY,
    user_address VARCHAR(42) NOT NULL,
    pool_id VARCHAR(50) NOT NULL,
    deposited VARCHAR(78) DEFAULT '0',
    earned VARCHAR(78) DEFAULT '0',
    points_earned VARCHAR(78) DEFAULT '0',
    last_updated TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_address, pool_id)
);

-- 3. DeFi 交易历史表
CREATE TABLE IF NOT EXISTS defi_transactions (
    id SERIAL PRIMARY KEY,
    user_address VARCHAR(42) NOT NULL,
    pool_id VARCHAR(50) NOT NULL,
    tx_type VARCHAR(20) NOT NULL, -- deposit, withdraw, claim
    amount VARCHAR(78) NOT NULL,
    tx_hash VARCHAR(66),
    block_number BIGINT,
    status VARCHAR(20) DEFAULT 'pending', -- pending, confirmed, failed
    timestamp TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_defi_tx_user ON defi_transactions(user_address);
CREATE INDEX IF NOT EXISTS idx_defi_tx_pool ON defi_transactions(pool_id);
CREATE INDEX IF NOT EXISTS idx_defi_tx_type ON defi_transactions(tx_type);
CREATE INDEX IF NOT EXISTS idx_defi_tx_time ON defi_transactions(timestamp);

-- 4. 稳定币仓位表
CREATE TABLE IF NOT EXISTS stablecoin_positions (
    id SERIAL PRIMARY KEY,
    user_address VARCHAR(42) NOT NULL UNIQUE,
    collateral_amount VARCHAR(78) DEFAULT '0',
    debt_amount VARCHAR(78) DEFAULT '0',
    collateral_ratio DECIMAL(10, 2) DEFAULT 0,
    health_status VARCHAR(20) DEFAULT 'safe', -- safe, warning, danger
    last_updated TIMESTAMP DEFAULT NOW()
);

-- 5. 稳定币交易历史
CREATE TABLE IF NOT EXISTS stablecoin_transactions (
    id SERIAL PRIMARY KEY,
    user_address VARCHAR(42) NOT NULL,
    tx_type VARCHAR(20) NOT NULL, -- mint, redeem
    amount VARCHAR(78) NOT NULL,
    fee VARCHAR(78),
    collateral_ratio_before DECIMAL(10, 2),
    collateral_ratio_after DECIMAL(10, 2),
    tx_hash VARCHAR(66),
    status VARCHAR(20) DEFAULT 'pending',
    timestamp TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_stable_tx_user ON stablecoin_transactions(user_address);
CREATE INDEX IF NOT EXISTS idx_stable_tx_type ON stablecoin_transactions(tx_type);
CREATE INDEX IF NOT EXISTS idx_stable_tx_time ON stablecoin_transactions(timestamp);

-- 插入初始 DeFi 协议池数据
INSERT INTO defi_pools (id, protocol, name, pool_type, apr, tvl, risk_level, icon, color, points_required, points_multiplier, features, metadata) VALUES
(
    'uniswap',
    'Uniswap V3',
    'PFI-USDC 流动性池',
    '流动性挖矿',
    '45.2%',
    '$2,450,000',
    '中等',
    '🦄',
    '#FF007A',
    1000,
    1.5,
    '["交易手续费收入", "积分奖励加成", "无常损失保护", "自动复投"]'::jsonb,
    '{"tradingFeeApr": "12.8%", "pointsRewardApr": "32.4%", "baseApr": "12.8%"}'::jsonb
),
(
    'aave',
    'Aave V3',
    'USDC 存款池',
    '借贷协议',
    '8.5%',
    '$8,920,000',
    '低',
    '👻',
    '#B6509E',
    500,
    1.2,
    '["稳定利息收入", "积分加速收益", "即时取款", "本金保护"]'::jsonb,
    '{"baseApr": "5.2%", "pointsBoostApr": "3.3%"}'::jsonb
),
(
    'stablecoin',
    'LoyaltyUSD',
    'LUSD 稳定币协议',
    '超额抵押稳定币',
    'N/A',
    '$18,450,000',
    '低',
    '💵',
    '#10b981',
    0,
    1.0,
    '["1:1 美元挂钩", "超额抵押保障", "铸造手续费 0.2%", "随时赎回抵押物"]'::jsonb,
    '{"collateralRatio": "150%", "liquidationThreshold": "120%", "mintFee": "0.2%", "redeemFee": "0.2%"}'::jsonb
),
(
    'staking',
    'PointFi Protocol',
    'PFI 质押池',
    '质押挖矿',
    '125%',
    '$5,680,000',
    '低',
    '🔒',
    '#6366F1',
    0,
    1.0,
    '["高额积分奖励", "质押即挖矿", "随时解除质押", "无锁定期"]'::jsonb,
    '{"stakingApr": "125%", "pointsRewardApr": "125%"}'::jsonb
)
ON CONFLICT (id) DO NOTHING;

-- 完成
COMMENT ON TABLE defi_pools IS 'DeFi protocol pool configurations';
COMMENT ON TABLE user_defi_positions IS 'User positions in DeFi pools';
COMMENT ON TABLE defi_transactions IS 'DeFi transaction history';
COMMENT ON TABLE stablecoin_positions IS 'User stablecoin positions (LUSD)';
COMMENT ON TABLE stablecoin_transactions IS 'Stablecoin mint/redeem history';
