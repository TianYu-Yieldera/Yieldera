# 🪙 LoyaltyUSD Smart Contracts

Complete smart contract suite for the LoyaltyUSD stablecoin protocol.

---

## 📂 Directory Structure

```
contracts/
├── core/                    # Core protocol contracts
│   ├── LoyaltyUSD.sol              ✅ ERC-20 stablecoin
│   ├── CollateralVault.sol         ✅ Collateral management
│   ├── StabilityManager.sol        🚧 Mint/redeem logic
│   └── LiquidationEngine.sol       📝 Liquidation execution
├── liquidity/               # DEX integration
│   ├── LiquidityVault.sol          📝 LP token staking
│   ├── RewardDistributor.sol       📝 Reward distribution
│   └── DEXConnector.sol            📝 Uniswap V3 adapter
├── oracle/                  # Price feeds
│   ├── OracleAdapter.sol           📝 Chainlink integration
│   └── TWAPOracle.sol              📝 Time-weighted price
├── governance/              # Protocol governance
│   ├── Timelock.sol                📝 Delayed execution
│   └── GovernanceModule.sol        📝 Parameter management
└── test/                    # Test contracts
    ├── MockERC20.sol
    └── MockOracle.sol
```

Legend:
- ✅ Completed
- 🚧 In Progress
- 📝 To Do

---

## 🔑 Core Contracts

### 1. LoyaltyUSD.sol

**Purpose:** USD-pegged stablecoin token

**Features:**
- ERC-20 standard compliance
- 6 decimals (matches USDC)
- Role-based access control
- Emergency pause capability
- Minting/burning controlled by StabilityManager

**Key Functions:**
```solidity
function mint(address to, uint256 amount) external onlyRole(MINTER_ROLE);
function burn(address from, uint256 amount) public;
function pause() external onlyRole(PAUSER_ROLE);
function unpause() external onlyRole(PAUSER_ROLE);
```

**Roles:**
- `MINTER_ROLE`: Can mint new LUSD (StabilityManager)
- `BURNER_ROLE`: Can burn LUSD (StabilityManager)
- `PAUSER_ROLE`: Can pause/unpause (Admin)

---

### 2. CollateralVault.sol

**Purpose:** Manages Loyalty Point collateral and debt tracking

**Features:**
- Deposits/withdrawals of LP tokens
- Debt tracking per user
- Collateral ratio calculations
- Liquidation checks
- Position health monitoring

**Key Functions:**
```solidity
function depositCollateral(uint256 amount) external;
function withdrawCollateral(uint256 amount) external;
function getMaxMintable(address user) public view returns (uint256);
function getCollateralRatio(address user) public view returns (uint256);
function liquidate(address user, uint256 debtToCover) external onlyOwner;
```

**Parameters:**
- Min Collateral Ratio: 150%
- Liquidation Threshold: 120%
- Stability Fee: 2% annual
- Liquidation Bonus: 10%

---

### 3. StabilityManager.sol (To Implement)

**Purpose:** Handles LUSD minting and redemption

**Key Features:**
- Mint LUSD against collateral
- Redeem collateral by burning LUSD
- Fee collection (0.2% mint/redeem)
- Position validation
- Integration with CollateralVault

**Expected Functions:**
```solidity
function mintLUSD(uint256 amount) external nonReentrant;
function redeemLUSD(uint256 amount) external nonReentrant;
function getPositionInfo(address user) external view returns (...);
```

---

### 4. LiquidationEngine.sol (To Implement)

**Purpose:** Executes liquidations of undercollateralized positions

**Key Features:**
- Monitors position health
- Triggers liquidations at 120% threshold
- Distributes liquidation bonus
- Updates collateral vault
- Emits liquidation events

**Expected Functions:**
```solidity
function liquidatePosition(address user) external;
function canLiquidate(address user) public view returns (bool);
function calculateLiquidationBonus(uint256 debt) public pure returns (uint256);
```

---

## 💧 Liquidity Contracts

### 5. LiquidityVault.sol (To Implement)

**Purpose:** Manages liquidity pool token staking

**Key Features:**
- Deposit/withdraw LP tokens
- Calculate pool shares
- Track rewards
- Integrate with Uniswap V3

---

### 6. RewardDistributor.sol (To Implement)

**Purpose:** Distributes liquidity mining rewards

**Key Features:**
- Calculate rewards per block
- Distribute PFI tokens
- Handle reward claims
- Update reward rates

---

## 🔮 Oracle Contracts

### 7. OracleAdapter.sol (To Implement)

**Purpose:** Price feed integration with Chainlink

**Key Features:**
- Fetch LP/USD price
- Fetch LUSD/USD price
- Staleness checks
- Fallback to TWAP

---

## 🏛️ Governance Contracts

### 8. Timelock.sol (To Implement)

**Purpose:** Delayed execution of governance actions

**Key Features:**
- Queue transactions
- Execute after delay
- Cancel transactions
- Admin controls

---

### 9. GovernanceModule.sol (To Implement)

**Purpose:** Protocol parameter management

**Key Features:**
- Adjust collateral ratio
- Modify fees
- Update oracle sources
- Emergency actions

---

## 🧪 Testing

### Running Tests

```bash
# Install dependencies
npm install

# Compile contracts
npx hardhat compile

# Run all tests
npx hardhat test

# Run with gas reporting
REPORT_GAS=true npx hardhat test

# Run coverage
npx hardhat coverage
```

### Test Structure

```
test/
├── unit/
│   ├── LoyaltyUSD.test.js
│   ├── CollateralVault.test.js
│   ├── StabilityManager.test.js
│   └── LiquidationEngine.test.js
├── integration/
│   ├── mint-redeem.test.js
│   ├── liquidation.test.js
│   └── liquidity-mining.test.js
└── fuzzing/
    └── collateral-ratio.test.js
```

---

## 🚀 Deployment

### Sepolia Testnet

```bash
# 1. Set environment variables
cp .env.example .env
# Edit .env with your keys

# 2. Deploy contracts
npx hardhat run scripts/deploy-stablecoin.js --network sepolia

# 3. Verify on Etherscan
npx hardhat verify --network sepolia <CONTRACT_ADDRESS>
```

### Deployment Script

```javascript
// scripts/deploy-stablecoin.js

async function main() {
  // 1. Deploy LoyaltyToken (if not exists)
  const LoyaltyToken = await ethers.getContractFactory("LoyaltyToken");
  const loyaltyToken = await LoyaltyToken.deploy();

  // 2. Deploy LoyaltyUSD
  const LoyaltyUSD = await ethers.getContractFactory("LoyaltyUSD");
  const lusd = await LoyaltyUSD.deploy();

  // 3. Deploy CollateralVault
  const CollateralVault = await ethers.getContractFactory("CollateralVault");
  const vault = await CollateralVault.deploy(loyaltyToken.address);

  // 4. Deploy StabilityManager
  const StabilityManager = await ethers.getContractFactory("StabilityManager");
  const manager = await StabilityManager.deploy(
    lusd.address,
    vault.address,
    feeCollector
  );

  // 5. Grant roles
  await lusd.grantRole(MINTER_ROLE, manager.address);
  await lusd.grantRole(BURNER_ROLE, manager.address);
  await vault.transferOwnership(manager.address);

  console.log("Deployment complete!");
}
```

---

## 🔒 Security

### Audit Checklist

- [ ] Reentrancy protection
- [ ] Integer overflow/underflow checks
- [ ] Access control verification
- [ ] Front-running prevention
- [ ] Flash loan attack resistance
- [ ] Oracle manipulation protection
- [ ] Emergency pause functionality

### Known Risks

1. **Oracle Risk**: Price feed manipulation or staleness
2. **Liquidation Risk**: Insufficient liquidators during market crashes
3. **Smart Contract Risk**: Bugs or vulnerabilities
4. **Governance Risk**: Malicious parameter changes

### Mitigation

- Multi-sig admin wallet
- Timelock for governance actions
- Circuit breakers for extreme conditions
- Insurance fund for bad debt
- Regular security audits

---

## 📊 Gas Optimization

### Estimated Gas Costs (Sepolia)

| Operation | Gas Cost | USD (at 50 gwei, $3000 ETH) |
|-----------|----------|------------------------------|
| Deposit Collateral | ~80,000 | $12.00 |
| Mint LUSD | ~120,000 | $18.00 |
| Redeem LUSD | ~110,000 | $16.50 |
| Withdraw Collateral | ~70,000 | $10.50 |
| Liquidate Position | ~150,000 | $22.50 |

### Optimization Techniques

- Packed storage variables
- Batch operations where possible
- Efficient loop iterations
- Minimal SLOAD operations
- Event emission optimization

---

## 📚 Additional Resources

- [OpenZeppelin Contracts](https://docs.openzeppelin.com/contracts/)
- [Hardhat Documentation](https://hardhat.org/docs)
- [Uniswap V3 Docs](https://docs.uniswap.org/protocol/introduction)
- [Chainlink Price Feeds](https://docs.chain.link/data-feeds)

---

## 🔗 Contract Addresses

### Sepolia Testnet

```json
{
  "network": "sepolia",
  "chainId": 11155111,
  "contracts": {
    "LoyaltyToken": "0x... (TBD)",
    "LoyaltyUSD": "0x... (TBD)",
    "CollateralVault": "0x... (TBD)",
    "StabilityManager": "0x... (TBD)",
    "LiquidationEngine": "0x... (TBD)",
    "OracleAdapter": "0x... (TBD)",
    "LiquidityVault": "0x... (TBD)",
    "RewardDistributor": "0x... (TBD)"
  },
  "uniswapV3": {
    "factory": "0x0227628f3F023bb0B980b67D528571c95c6DaC1c",
    "router": "0x3bFA4769FB09eefC5a80d6E87c3B9C650f7Ae48E",
    "quoter": "0xEd1f6473345F45b75F8179591dd5bA1888cf2FB3"
  },
  "chainlink": {
    "priceFeed": "0x... (TBD)"
  }
}
```

---

## 📝 License

MIT License - See LICENSE file for details

---

**Status:** 🚧 Active Development
**Version:** v0.1.0 (Pre-Alpha)
**Last Updated:** 2025-10-16
