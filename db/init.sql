-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  address TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS balances (
  id SERIAL PRIMARY KEY,
  user_address TEXT UNIQUE NOT NULL,
  balance NUMERIC(78, 18) NOT NULL DEFAULT 0,
  is_demo BOOLEAN DEFAULT FALSE,
  updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS balance_events (
  id BIGSERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL,
  event_type TEXT NOT NULL,
  tx_hash TEXT,
  chain TEXT,
  block_number BIGINT,
  confirmed BOOLEAN DEFAULT TRUE,
  is_demo BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_balance_events_user ON balance_events(user_address);
CREATE TABLE IF NOT EXISTS points (
  id SERIAL PRIMARY KEY,
  user_address TEXT UNIQUE NOT NULL,
  points NUMERIC(78, 18) NOT NULL DEFAULT 0,
  is_demo BOOLEAN DEFAULT FALSE,
  updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS points_events (
  id BIGSERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  points_delta NUMERIC(78, 18) NOT NULL,
  reason TEXT,
  is_demo BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS badges (
  id SERIAL PRIMARY KEY,
  user_address TEXT NOT NULL,
  badge_code TEXT NOT NULL,
  is_demo BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Airdrop Feature Tables

-- Admin whitelist for campaign management
CREATE TABLE IF NOT EXISTS admin_whitelist (
  id SERIAL PRIMARY KEY,
  address TEXT UNIQUE NOT NULL,
  name TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Airdrop campaigns
CREATE TABLE IF NOT EXISTS airdrop_campaigns (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  asset_type TEXT DEFAULT 'points' CHECK (asset_type IN ('points', 'tokens', 'native')),
  status TEXT DEFAULT 'draft' CHECK (status IN ('draft', 'scheduled', 'active', 'claimable', 'closed', 'archived')),
  start_time TIMESTAMP,
  end_time TIMESTAMP,
  total_budget NUMERIC(78, 18) NOT NULL DEFAULT 0,
  claimed_amount NUMERIC(78, 18) NOT NULL DEFAULT 0,
  participant_count INT DEFAULT 0,
  is_demo BOOLEAN DEFAULT FALSE,
  created_by TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_airdrop_campaigns_status ON airdrop_campaigns(status);
CREATE INDEX IF NOT EXISTS idx_airdrop_campaigns_created_by ON airdrop_campaigns(created_by);

-- Airdrop allocations (whitelist)
CREATE TABLE IF NOT EXISTS airdrop_allocations (
  id SERIAL PRIMARY KEY,
  campaign_id INT NOT NULL REFERENCES airdrop_campaigns(id) ON DELETE CASCADE,
  user_address TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(campaign_id, user_address)
);
CREATE INDEX IF NOT EXISTS idx_airdrop_allocations_campaign ON airdrop_allocations(campaign_id);
CREATE INDEX IF NOT EXISTS idx_airdrop_allocations_user ON airdrop_allocations(user_address);

-- Airdrop claims (claim records)
CREATE TABLE IF NOT EXISTS airdrop_claims (
  id SERIAL PRIMARY KEY,
  campaign_id INT NOT NULL REFERENCES airdrop_campaigns(id) ON DELETE CASCADE,
  user_address TEXT NOT NULL,
  amount NUMERIC(78, 18) NOT NULL,
  nonce TEXT NOT NULL,
  signature TEXT,
  claimed_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(campaign_id, user_address)
);
CREATE INDEX IF NOT EXISTS idx_airdrop_claims_campaign ON airdrop_claims(campaign_id);
CREATE INDEX IF NOT EXISTS idx_airdrop_claims_user ON airdrop_claims(user_address);
