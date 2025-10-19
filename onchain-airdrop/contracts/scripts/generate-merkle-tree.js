const { MerkleTree } = require("merkletreejs");
const keccak256 = require("keccak256");
const fs = require("fs");
const path = require("path");
const { ethers } = require("hardhat");

/**
 * 从 CSV 文件生成 Merkle Tree
 *
 * CSV 格式:
 * address,amount
 * 0x1234...,1000000000000000000
 * 0x5678...,2000000000000000000
 *
 * 用法:
 * node scripts/generate-merkle-tree.js <csv_file_path>
 *
 * 示例:
 * node scripts/generate-merkle-tree.js data/whitelist.csv
 */

// 生成叶子节点（Solidity keccak256(abi.encodePacked(address, amount))）
function hashLeaf(address, amount) {
  return Buffer.from(
    ethers.solidityPackedKeccak256(
      ["address", "uint256"],
      [address, amount]
    ).slice(2),
    "hex"
  );
}

// 解析 CSV 文件
function parseCSV(filePath) {
  const content = fs.readFileSync(filePath, "utf-8");
  const lines = content.trim().split("\n");

  // 跳过标题行
  const data = lines.slice(1).map((line) => {
    const [address, amount] = line.split(",").map((s) => s.trim());

    // 验证地址
    if (!ethers.isAddress(address)) {
      throw new Error(`无效的地址: ${address}`);
    }

    // 验证金额
    const parsedAmount = BigInt(amount);
    if (parsedAmount <= 0n) {
      throw new Error(`无效的金额: ${amount}`);
    }

    return {
      address: ethers.getAddress(address), // 标准化地址
      amount: amount,
    };
  });

  return data;
}

// 生成示例 CSV（如果没有提供文件）
function generateExampleCSV() {
  const examplePath = path.join(__dirname, "../data/whitelist-example.csv");
  const dataDir = path.dirname(examplePath);

  if (!fs.existsSync(dataDir)) {
    fs.mkdirSync(dataDir, { recursive: true });
  }

  const content = `address,amount
0x70997970C51812dc3A010C7d01b50e0d17dc79C8,1000000000000000000000
0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC,2000000000000000000000
0x90F79bf6EB2c4f870365E785982E1f101E93b906,500000000000000000000
0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65,750000000000000000000
0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc,1500000000000000000000`;

  fs.writeFileSync(examplePath, content);
  console.log("✅ 示例 CSV 已生成:", examplePath, "\n");
  return examplePath;
}

async function main() {
  console.log("🌳 开始生成 Merkle Tree...\n");

  // 获取 CSV 文件路径
  let csvPath = process.argv[2];

  if (!csvPath) {
    console.log("⚠️  未提供 CSV 文件，使用示例数据...");
    csvPath = generateExampleCSV();
  }

  if (!fs.existsSync(csvPath)) {
    console.error("❌ 错误: CSV 文件不存在:", csvPath);
    process.exit(1);
  }

  console.log("📄 读取 CSV 文件:", csvPath);

  // 解析 CSV
  const whitelist = parseCSV(csvPath);
  console.log(`✅ 成功解析 ${whitelist.length} 个地址\n`);

  // 生成叶子节点
  const leaves = whitelist.map((entry) =>
    hashLeaf(entry.address, entry.amount)
  );

  // 构建 Merkle Tree
  const tree = new MerkleTree(leaves, keccak256, { sortPairs: true });
  const root = tree.getHexRoot();

  console.log("🌲 Merkle Tree 信息:");
  console.log("   根哈希 (Merkle Root):", root);
  console.log("   叶子节点数量:", leaves.length);
  console.log("   树的深度:", tree.getDepth(), "\n");

  // 生成每个地址的 Merkle Proof
  const proofs = {};
  whitelist.forEach((entry, index) => {
    const leaf = hashLeaf(entry.address, entry.amount);
    const proof = tree.getHexProof(leaf);

    proofs[entry.address] = {
      address: entry.address,
      amount: entry.amount,
      amountFormatted: ethers.formatEther(entry.amount) + " YLD",
      proof: proof,
      leaf: ethers.hexlify(leaf),
    };
  });

  // 保存结果
  const outputDir = path.join(__dirname, "../data");
  if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
  }

  // 1. 保存完整数据（包含所有 proofs）
  const fullOutputPath = path.join(
    outputDir,
    `merkle-tree-${Date.now()}.json`
  );
  const fullData = {
    merkleRoot: root,
    totalEntries: whitelist.length,
    timestamp: new Date().toISOString(),
    proofs: proofs,
  };
  fs.writeFileSync(fullOutputPath, JSON.stringify(fullData, null, 2));
  console.log("📦 完整数据已保存:", fullOutputPath);

  // 2. 保存简化数据（仅 root + 总结）
  const summaryPath = path.join(outputDir, "merkle-summary.json");
  const summary = {
    merkleRoot: root,
    totalEntries: whitelist.length,
    treeDepth: tree.getDepth(),
    timestamp: new Date().toISOString(),
  };
  fs.writeFileSync(summaryPath, JSON.stringify(summary, null, 2));
  console.log("📄 摘要数据已保存:", summaryPath, "\n");

  // 打印前 3 个地址的 proof（用于测试）
  console.log("📋 前 3 个地址的 Merkle Proof:");
  whitelist.slice(0, 3).forEach((entry) => {
    const data = proofs[entry.address];
    console.log(`\n   地址: ${data.address}`);
    console.log(`   金额: ${data.amountFormatted}`);
    console.log(`   Proof: [${data.proof.length} 个哈希]`);
    console.log(`   ${data.proof.join(", ")}`);
  });

  console.log("\n" + "=".repeat(60));
  console.log("🎉 Merkle Tree 生成完成！");
  console.log("=".repeat(60));
  console.log("Merkle Root:", root);
  console.log("总地址数:", whitelist.length);
  console.log("\n📝 下一步:");
  console.log("   1. 使用 Merkle Root 创建空投活动");
  console.log("   2. 将 proof 数据部署到前端或 IPFS");
  console.log("   3. 用户使用 proof 领取空投");
  console.log("=".repeat(60));
  console.log("");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("❌ 错误:", error.message);
    process.exit(1);
  });
