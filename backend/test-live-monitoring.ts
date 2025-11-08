/**
 * 实时监控测试 - 只启动Treasury监听器
 */

import dotenv from 'dotenv';
import { MarketplaceListener } from './listeners/treasury/MarketplaceListener';
import { AssetFactoryListener } from './listeners/treasury/AssetFactoryListener';

dotenv.config();

console.log('🚀 Starting Treasury Monitoring Test\n');

const wsUrl = process.env.ARBITRUM_SEPOLIA_WS!;
const marketplaceAddress = process.env.TREASURY_MARKETPLACE_ADDRESS!;
const assetFactoryAddress = process.env.TREASURY_ASSET_FACTORY_ADDRESS!;

console.log(`WebSocket URL: ${wsUrl.substring(0, 50)}...`);
console.log(`Marketplace: ${marketplaceAddress}`);
console.log(`AssetFactory: ${assetFactoryAddress}\n`);

async function main() {
  // 启动Marketplace监听器
  const marketplaceListener = new MarketplaceListener(wsUrl, marketplaceAddress);

  // 监听事件
  marketplaceListener.on('orderCreated', (data) => {
    console.log('📝 Order Created:', {
      orderId: data.orderId,
      seller: data.seller.substring(0, 10) + '...',
      assetId: data.assetId,
      amount: data.amount,
    });
  });

  marketplaceListener.on('orderFilled', (data) => {
    console.log('✅ Order Filled:', {
      orderId: data.orderId,
      buyer: data.buyer.substring(0, 10) + '...',
      totalPrice: data.totalPrice,
    });
  });

  marketplaceListener.on('alert', (alert) => {
    const emoji = alert.severity === 'critical' ? '🚨' :
                  alert.severity === 'warning' ? '⚠️' : 'ℹ️';
    console.log(`${emoji} ALERT [${alert.severity}]: ${alert.message}`);
  });

  marketplaceListener.on('error', (error) => {
    console.error('❌ Marketplace error:', error.message);
  });

  // 启动AssetFactory监听器 (延迟5秒避免速率限制)
  setTimeout(async () => {
    const assetFactoryListener = new AssetFactoryListener(wsUrl, assetFactoryAddress);

    assetFactoryListener.on('assetCreated', (data) => {
      console.log('🆕 Asset Created:', {
        assetId: data.assetId,
        symbol: data.symbol,
        totalValue: data.totalValue,
      });
    });

    assetFactoryListener.on('assetVerified', (data) => {
      console.log('✔️ Asset Verified:', {
        assetId: data.assetId,
        verifier: data.verifier.substring(0, 10) + '...',
      });
    });

    assetFactoryListener.on('alert', (alert) => {
      const emoji = alert.severity === 'critical' ? '🚨' :
                    alert.severity === 'warning' ? '⚠️' : 'ℹ️';
      console.log(`${emoji} ALERT [${alert.severity}]: ${alert.message}`);
    });

    assetFactoryListener.on('error', (error) => {
      console.error('❌ AssetFactory error:', error.message);
    });

    await assetFactoryListener.start();
    console.log('✅ AssetFactory listener started\n');

  }, 5000);

  // 启动监听
  await marketplaceListener.start();
  console.log('✅ Marketplace listener started');
  console.log('⏳ Starting AssetFactory listener in 5 seconds...\n');

  // 每30秒输出统计
  setInterval(() => {
    console.log('\n📊 === Stats Report ===');
    console.log('Marketplace:', marketplaceListener.getStats());
    console.log('======================\n');
  }, 30000);

  // 保持运行
  console.log('👀 Monitoring for events... (Press Ctrl+C to stop)\n');
}

main().catch(console.error);

// 优雅退出
process.on('SIGINT', () => {
  console.log('\n\n🛑 Shutting down...');
  process.exit(0);
});
