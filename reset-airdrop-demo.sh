#!/bin/bash

echo "=========================================="
echo "空投演示重置工具"
echo "=========================================="
echo ""
echo "请选择重置方式："
echo "1. 仅重置当前账户 (0x3c072346)"
echo "2. 重置所有领取记录"
echo "3. 完整重置（包括重新导入测试数据）"
echo ""
read -p "输入选项 (1/2/3): " choice

case $choice in
  1)
    echo "正在重置账户 0x3c072346 的领取记录..."
    docker exec -i loyalty-postgres psql -U loyalty_user -d loyalty_db << 'EOF'
DELETE FROM airdrop_claims WHERE user_address = '0x3c072346';
UPDATE airdrop_campaigns
SET claimed_amount = GREATEST(claimed_amount - (
  SELECT COALESCE(SUM(CAST(amount AS NUMERIC)), 0)
  FROM airdrop_allocations
  WHERE user_address = '0x3c072346' AND campaign_id = airdrop_campaigns.id
), 0),
participant_count = GREATEST(participant_count - 1, 0);
SELECT '✅ 已重置账户 0x3c072346，可以再次领取！';
EOF
    ;;

  2)
    echo "正在重置所有领取记录..."
    docker exec -i loyalty-postgres psql -U loyalty_user -d loyalty_db << 'EOF'
TRUNCATE airdrop_claims;
UPDATE airdrop_campaigns SET claimed_amount = '0', participant_count = 0;
SELECT '✅ 已重置所有领取记录，所有账户都可以重新领取！';
EOF
    ;;

  3)
    echo "正在完整重置并重新导入测试数据..."
    docker exec -i loyalty-postgres psql -U loyalty_user -d loyalty_db << 'EOF'
-- 清空所有空投数据
TRUNCATE airdrop_claims, airdrop_allocations, airdrop_campaigns, admin_whitelist CASCADE;

-- 重新导入admin
INSERT INTO admin_whitelist (address, name) VALUES
  ('0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266', 'Test Admin 1'),
  ('0x70997970c51812dc3a010c7d01b50e0d17dc79c8', 'Test Admin 2'),
  ('0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc', 'Test Admin 3'),
  ('0x3c072346', 'Current User')
ON CONFLICT DO NOTHING;

-- 重新创建测试活动
INSERT INTO airdrop_campaigns
  (name, description, asset_type, status, start_time, end_time, total_budget, is_demo, created_by)
VALUES
  ('Season 1 Early Birds', 'Reward for early testers', 'points', 'active',
   '2025-01-01', '2025-12-31', '100000', true, '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266'),
  ('Community Builders', 'Reward for active members', 'points', 'scheduled',
   '2026-01-01', '2026-06-30', '50000', true, '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266')
RETURNING id;

-- 重新导入分配
INSERT INTO airdrop_allocations (campaign_id, user_address, amount) VALUES
  (1, '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266', '5000'),
  (1, '0x70997970c51812dc3a010c7d01b50e0d17dc79c8', '3000'),
  (1, '0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc', '2000'),
  (1, '0x3c072346', '10000'),
  (2, '0x3c072346', '5000')
ON CONFLICT DO NOTHING;

SELECT '✅ 完整重置成功！所有数据已恢复到初始状态。';
EOF
    ;;

  *)
    echo "无效选项"
    exit 1
    ;;
esac

echo ""
echo "=========================================="
echo "✅ 重置完成！"
echo "=========================================="
echo ""
echo "现在您可以："
echo "1. 刷新浏览器页面"
echo "2. 重新领取空投"
echo ""
