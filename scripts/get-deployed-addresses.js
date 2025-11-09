/**
 * 获取已部署的合约地址
 * 用于配置监控系统
 */

const fs = require('fs');
const path = require('path');

async function getDeployedAddresses() {
  console.log('📋 Deployed Contract Addresses\n');

  // 读取Treasury部署记录
  const treasuryDeployment = path.join(__dirname, '../deployments/treasury-arbitrumSepolia-1762266915826.json');

  if (fs.existsSync(treasuryDeployment)) {
    const data = JSON.parse(fs.readFileSync(treasuryDeployment, 'utf8'));

    console.log('✅ Treasury Contracts (Arbitrum Sepolia):');
    console.log(`  Network: ${data.network} (Chain ID: ${data.chainId})`);
    console.log(`  Deployed: ${data.timestamp}\n`);

    console.log('  Contracts:');
    for (const [name, address] of Object.entries(data.contracts)) {
      console.log(`    ${name}: ${address}`);
    }
    console.log('');
  }

  // 检查L2适配器部署
  const l2DeploymentPattern = /adapters.*\.json$/;
  const deploymentsDir = path.join(__dirname, '../deployments');

  if (fs.existsSync(deploymentsDir)) {
    const files = fs.readdirSync(deploymentsDir);
    const l2Files = files.filter(f => l2DeploymentPattern.test(f));

    if (l2Files.length > 0) {
      console.log('✅ DeFi Adapters (Arbitrum Sepolia):');
      l2Files.forEach(file => {
        const data = JSON.parse(fs.readFileSync(path.join(deploymentsDir, file), 'utf8'));
        console.log(`  File: ${file}`);
        for (const [name, address] of Object.entries(data.contracts || data)) {
          console.log(`    ${name}: ${address}`);
        }
      });
    } else {
      console.log('⚠️  DeFi Adapters not yet deployed');
      console.log('   Run: npx hardhat run scripts/layer2/deploy-l2.js --network arbitrumSepolia\n');
    }
  }

  // 生成.env配置
  console.log('\n📝 Add to backend/.env:');
  console.log('─'.repeat(60));

  if (fs.existsSync(treasuryDeployment)) {
    const data = JSON.parse(fs.readFileSync(treasuryDeployment, 'utf8'));
    console.log('# Treasury Contracts');
    for (const [name, address] of Object.entries(data.contracts)) {
      const envName = name.toUpperCase().replace(/([A-Z])/g, '_$1').replace(/^_/, '');
      console.log(`${envName}_ADDRESS=${address}`);
    }
  }
}

getDeployedAddresses().catch(console.error);
