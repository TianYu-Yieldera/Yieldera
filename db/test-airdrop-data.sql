-- Test Data for Airdrop Feature
-- Run this script to set up test data for airdrop functionality

-- 1. Add test admin addresses
INSERT INTO admin_whitelist (address, name) VALUES
  -- Common test addresses from MetaMask/Hardhat
  ('0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266', 'Test Admin 1 (Hardhat Account 0)'),
  ('0x70997970c51812dc3a010c7d01b50e0d17dc79c8', 'Test Admin 2 (Hardhat Account 1)'),
  ('0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc', 'Test Admin 3 (Hardhat Account 2)'),
  -- Add your personal address here (lowercase)
  ('0x742d35cc6634c0532925a3b844bc9e7595f0beb1', 'Personal Test Admin')
ON CONFLICT (address) DO NOTHING;

-- 2. Create a test campaign
INSERT INTO airdrop_campaigns
  (name, description, asset_type, status, start_time, end_time, total_budget, claimed_amount, participant_count, is_demo, created_by)
VALUES
  (
    'Season 1 Early Birds',
    'Reward for early testers and contributors',
    'points',
    'active',
    '2025-01-01 00:00:00',
    '2025-12-31 23:59:59',
    '100000',
    '0',
    0,
    true,
    '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266'
  )
ON CONFLICT DO NOTHING
RETURNING id;

-- 3. Add test allocations (whitelist)
-- NOTE: Replace campaign_id=1 with the actual ID from above if different
INSERT INTO airdrop_allocations (campaign_id, user_address, amount) VALUES
  (1, '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266', '5000'),
  (1, '0x70997970c51812dc3a010c7d01b50e0d17dc79c8', '3000'),
  (1, '0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc', '2000'),
  (1, '0x90f79bf6eb2c4f870365e785982e1f101e93b906', '1000'),
  (1, '0x15d34aaf54267db7d7c367839aaf71a00a2c6a65', '1500'),
  (1, '0x9965507d1a55bcc2695c58ba16fb37d819b0a4dc', '2500'),
  (1, '0x976ea74026e726554db657fa54763abd0c3a0aa9', '4000'),
  (1, '0x14dc79964da2c08b23698b3d3cc7ca32193d9955', '3500'),
  (1, '0x23618e81e3f5cdf7f54c3d65f7fbc0abf5b21e8f', '1200'),
  (1, '0xa0ee7a142d267c1f36714e4a8f75612f20a79720', '1800'),
  -- Add your personal address here too if you want to test claiming
  (1, '0x742d35cc6634c0532925a3b844bc9e7595f0beb1', '10000')
ON CONFLICT (campaign_id, user_address) DO NOTHING;

-- 4. Create a second campaign (scheduled for future)
INSERT INTO airdrop_campaigns
  (name, description, asset_type, status, start_time, end_time, total_budget, claimed_amount, participant_count, is_demo, created_by)
VALUES
  (
    'Community Builders',
    'Reward for active community members',
    'points',
    'scheduled',
    '2026-01-01 00:00:00',
    '2026-06-30 23:59:59',
    '50000',
    '0',
    0,
    true,
    '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266'
  )
ON CONFLICT DO NOTHING;

-- 5. Verify data
SELECT 'Admin Whitelist Count:' as info, COUNT(*) as count FROM admin_whitelist
UNION ALL
SELECT 'Campaigns Count:', COUNT(*) FROM airdrop_campaigns
UNION ALL
SELECT 'Allocations Count:', COUNT(*) FROM airdrop_allocations
UNION ALL
SELECT 'Claims Count:', COUNT(*) FROM airdrop_claims;

-- Show created campaigns
SELECT id, name, status, total_budget, created_by
FROM airdrop_campaigns
ORDER BY id;
