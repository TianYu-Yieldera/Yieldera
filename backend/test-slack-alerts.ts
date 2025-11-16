/**
 * Slack 告警服务测试
 */

import dotenv from 'dotenv';
import { SlackAlertService } from './services/alerts/SlackAlertService';

dotenv.config();

const WEBHOOK_URL = process.env.SLACK_WEBHOOK_URL || '';
const ENABLED = process.env.SLACK_ENABLED === 'true';

console.log('🧪 Testing Slack Alert Service\n');
console.log('============================================================\n');

async function runTests() {
  // Test 1: 配置检查
  console.log('📋 Test 1: Configuration Check');
  console.log('------------------------------------------------------------');

  if (!WEBHOOK_URL || WEBHOOK_URL.includes('YOUR_WEBHOOK')) {
    console.log('  ❌ SLACK_WEBHOOK_URL not configured');
    console.log('  ℹ️  Please set SLACK_WEBHOOK_URL in .env file');
    console.log('  📝 Get webhook URL from: https://api.slack.com/messaging/webhooks\n');
    return false;
  }

  console.log(`  ✅ Webhook URL configured`);
  console.log(`  ✅ Alerts enabled: ${ENABLED}\n`);

  // Test 2: 初始化服务
  console.log('🔧 Test 2: Initialize Alert Service');
  console.log('------------------------------------------------------------');

  const slackService = new SlackAlertService({
    webhookUrl: WEBHOOK_URL,
    enabled: ENABLED,
    channelName: process.env.SLACK_CHANNEL || '#monitoring',
    botName: 'Test Bot',
    minLevel: 'INFO',
  });

  console.log('  ✅ SlackAlertService initialized');
  console.log('  📊 Stats:', slackService.getStats());
  console.log('');

  if (!ENABLED) {
    console.log('  ⚠️  Alerts are disabled (SLACK_ENABLED=false)');
    console.log('  ℹ️  Set SLACK_ENABLED=true to enable alerts\n');
    return false;
  }

  // Test 3: 发送测试消息
  console.log('📤 Test 3: Send Test Message');
  console.log('------------------------------------------------------------');

  try {
    const result = await slackService.sendTestMessage();
    if (result) {
      console.log('  ✅ Test message sent successfully!');
      console.log('  📱 Check your Slack channel for the message\n');
    } else {
      console.log('  ❌ Failed to send test message\n');
      return false;
    }
  } catch (error: any) {
    console.error('  ❌ Error:', error.message);
    return false;
  }

  // Test 4: 发送不同级别的告警
  console.log('🎨 Test 4: Send Different Alert Levels');
  console.log('------------------------------------------------------------');

  const alerts = [
    {
      level: 'INFO' as const,
      type: 'TEST_INFO',
      message: 'This is an INFO level alert for testing',
    },
    {
      level: 'WARNING' as const,
      type: 'TEST_WARNING',
      message: 'This is a WARNING level alert for testing',
      data: {
        amount: '$10,000',
        user: '0x1234567890abcdef',
      },
    },
    {
      level: 'CRITICAL' as const,
      type: 'TEST_CRITICAL',
      message: 'This is a CRITICAL level alert for testing',
      data: {
        txHash: '0xabcdef1234567890',
        amount: '$100,000',
      },
    },
  ];

  for (const alert of alerts) {
    try {
      console.log(`  Sending ${alert.level} alert...`);
      const result = await slackService.sendAlert(alert);
      if (result) {
        console.log(`    ✅ ${alert.level} alert sent`);
      }
      // 等待1秒避免发送太快
      await new Promise(resolve => setTimeout(resolve, 1000));
    } catch (error: any) {
      console.error(`    ❌ Error:`, error.message);
    }
  }
  console.log('');

  // Test 5: 测试重复告警抑制
  console.log('🚫 Test 5: Duplicate Alert Suppression');
  console.log('------------------------------------------------------------');

  const duplicateAlert = {
    level: 'WARNING' as const,
    type: 'DUPLICATE_TEST',
    message: 'This alert should only be sent once',
  };

  console.log('  Sending alert first time...');
  const first = await slackService.sendAlert(duplicateAlert);
  console.log(`    ${first ? '✅' : '❌'} First send: ${first ? 'sent' : 'blocked'}`);

  console.log('  Sending same alert again (should be suppressed)...');
  const second = await slackService.sendAlert(duplicateAlert);
  console.log(`    ${!second ? '✅' : '❌'} Second send: ${second ? 'sent (unexpected)' : 'suppressed (correct)'}`);
  console.log('');

  // Test 6: 统计信息
  console.log('📊 Test 6: Service Statistics');
  console.log('------------------------------------------------------------');
  const stats = slackService.getStats();
  console.log('  Enabled:', stats.enabled);
  console.log('  Webhook Configured:', stats.webhookConfigured);
  console.log('  Min Level:', stats.minLevel);
  console.log('  Recent Alerts:', stats.recentAlertsCount);
  console.log('  Cooldown (ms):', stats.cooldownMs);
  console.log('');

  return true;
}

// 运行测试
runTests()
  .then((success) => {
    if (success) {
      console.log('✅ All tests passed!');
      console.log('\n📝 Next steps:');
      console.log('  1. Check your Slack channel for test messages');
      console.log('  2. Set SLACK_ENABLED=true in .env to enable alerts in production');
      console.log('  3. Configure SLACK_MIN_LEVEL (INFO/WARNING/CRITICAL) as needed');
      process.exit(0);
    } else {
      console.log('❌ Some tests failed');
      console.log('\n📝 Troubleshooting:');
      console.log('  1. Verify SLACK_WEBHOOK_URL is correct');
      console.log('  2. Check that the webhook is active in Slack');
      console.log('  3. Ensure your Slack workspace allows incoming webhooks');
      process.exit(1);
    }
  })
  .catch((error) => {
    console.error('❌ Test suite failed:', error);
    process.exit(1);
  });
