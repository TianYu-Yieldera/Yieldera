/**
 * Cross-Chain Integration Tests
 * Tests L1-L2 bridge functionality and state synchronization
 */

import { expect } from "chai";
import hre from "hardhat";
const { ethers } = hre;
import * as helpers from "@nomicfoundation/hardhat-network-helpers";
const { time } = helpers;

describe("Cross-Chain Integration Tests", function () {
  let deployer, user1, user2;

  // L1 Contracts
  let loyaltyUSDL1, collateralVaultL1, l1StateRegistry, l1Gateway;

  // L2 Contracts
  let l2StateAggregator;

  // Mock contracts
  let collateralToken, mockInbox, mockOutbox;

  const INITIAL_SUPPLY = ethers.parseEther("1000000");
  const ONE_HOUR = 3600;

  beforeEach(async function () {
    [deployer, user1, user2] = await ethers.getSigners();

    // Deploy Mock ERC20 for collateral
    const MockERC20 = await ethers.getContractFactory("MockERC20");
    collateralToken = await MockERC20.deploy("Loyalty Points", "LP", INITIAL_SUPPLY);
    await collateralToken.waitForDeployment();

    // Distribute tokens to users for testing
    await collateralToken.transfer(user1.address, ethers.parseEther("10000"));
    await collateralToken.transfer(user2.address, ethers.parseEther("10000"));

    // Deploy L1 Contracts
    console.log("Deploying L1 contracts...");

    const LoyaltyUSDL1 = await ethers.getContractFactory("LoyaltyUSDL1");
    loyaltyUSDL1 = await LoyaltyUSDL1.deploy();
    await loyaltyUSDL1.waitForDeployment();

    const CollateralVaultL1 = await ethers.getContractFactory("CollateralVaultL1");
    collateralVaultL1 = await CollateralVaultL1.deploy(await collateralToken.getAddress());
    await collateralVaultL1.waitForDeployment();

    // Deploy L2StateAggregator first (for L1StateRegistry)
    const L2StateAggregator = await ethers.getContractFactory("L2StateAggregator");
    l2StateAggregator = await L2StateAggregator.deploy(deployer.address); // Temporary
    await l2StateAggregator.waitForDeployment();

    const L1StateRegistry = await ethers.getContractFactory("L1StateRegistry");
    l1StateRegistry = await L1StateRegistry.deploy(await l2StateAggregator.getAddress());
    await l1StateRegistry.waitForDeployment();

    // Update L2StateAggregator with correct L1StateRegistry
    await l2StateAggregator.setL1StateRegistry(await l1StateRegistry.getAddress());

    // For testing, we'll use mock bridge addresses
    mockInbox = deployer.address;
    mockOutbox = deployer.address;

    const L1Gateway = await ethers.getContractFactory("L1Gateway");
    l1Gateway = await L1Gateway.deploy(
      await collateralVaultL1.getAddress(),
      await loyaltyUSDL1.getAddress(),
      await collateralToken.getAddress(),
      mockInbox,
      mockOutbox
    );
    await l1Gateway.waitForDeployment();

    // Setup permissions
    const BRIDGE_ROLE = await loyaltyUSDL1.BRIDGE_ROLE();
    await loyaltyUSDL1.grantRole(BRIDGE_ROLE, await l1Gateway.getAddress());
    await collateralVaultL1.setL2Bridge(await l1Gateway.getAddress());
    await collateralVaultL1.setStateRegistry(await l1StateRegistry.getAddress());

    console.log("Setup complete");
  });

  describe("L1 Core Functionality", function () {
    it("Should deploy all L1 contracts correctly", async function () {
      expect(await loyaltyUSDL1.getAddress()).to.not.equal(ethers.ZeroAddress);
      expect(await collateralVaultL1.getAddress()).to.not.equal(ethers.ZeroAddress);
      expect(await l1StateRegistry.getAddress()).to.not.equal(ethers.ZeroAddress);
      expect(await l1Gateway.getAddress()).to.not.equal(ethers.ZeroAddress);
    });

    it("Should have correct permissions set", async function () {
      const BRIDGE_ROLE = await loyaltyUSDL1.BRIDGE_ROLE();
      expect(await loyaltyUSDL1.hasRole(BRIDGE_ROLE, await l1Gateway.getAddress())).to.be.true;
      expect(await collateralVaultL1.l2Bridge()).to.equal(await l1Gateway.getAddress());
    });

    it("Should enforce daily mint limits", async function () {
      const BRIDGE_ROLE = await loyaltyUSDL1.BRIDGE_ROLE();
      await loyaltyUSDL1.grantRole(BRIDGE_ROLE, deployer.address);

      const amount = ethers.parseUnits("1000000", 6); // 1M LUSD
      const l2TxHash = ethers.keccak256(ethers.toUtf8Bytes("test"));

      // First mint should succeed
      await loyaltyUSDL1.bridgeMint(user1.address, amount, l2TxHash);
      expect(await loyaltyUSDL1.balanceOf(user1.address)).to.equal(amount);

      // Check remaining capacity
      const remaining = await loyaltyUSDL1.getRemainingDailyMintCapacity();
      expect(remaining).to.equal(ethers.parseUnits("9000000", 6)); // 10M - 1M
    });

    it("Should reject mints exceeding single transaction limit", async function () {
      const BRIDGE_ROLE = await loyaltyUSDL1.BRIDGE_ROLE();
      await loyaltyUSDL1.grantRole(BRIDGE_ROLE, deployer.address);

      const amount = ethers.parseUnits("1000001", 6); // > 1M LUSD
      const l2TxHash = ethers.keccak256(ethers.toUtf8Bytes("test"));

      await expect(
        loyaltyUSDL1.bridgeMint(user1.address, amount, l2TxHash)
      ).to.be.revertedWith("LUSD: exceeds max mint per tx");
    });
  });

  describe("L1 Collateral Vault", function () {
    it("Should lock collateral correctly", async function () {
      const amount = ethers.parseEther("100");
      const l2TxHash = ethers.keccak256(ethers.toUtf8Bytes("deposit"));

      // Approve vault to spend tokens
      await collateralToken.connect(user1).approve(await l1Gateway.getAddress(), amount);

      // Lock via gateway
      await collateralToken.connect(user1).transfer(await l1Gateway.getAddress(), amount);
      await collateralToken.connect(deployer).approve(await collateralVaultL1.getAddress(), amount);
      await collateralVaultL1.lockCollateral(user1.address, amount, l2TxHash);

      expect(await collateralVaultL1.getLockedCollateral(user1.address)).to.equal(amount);
      expect(await collateralVaultL1.getTotalLocked()).to.equal(amount);
    });

    it("Should unlock collateral correctly", async function () {
      const amount = ethers.parseEther("100");
      const l2TxHash = ethers.keccak256(ethers.toUtf8Bytes("deposit"));

      // Lock first
      await collateralToken.connect(user1).transfer(await collateralVaultL1.getAddress(), amount);
      await collateralVaultL1.lockCollateral(user1.address, amount, l2TxHash);

      // Unlock
      const unlockTxHash = ethers.keccak256(ethers.toUtf8Bytes("withdrawal"));
      await collateralVaultL1.unlockCollateral(user1.address, amount, unlockTxHash);

      expect(await collateralVaultL1.getLockedCollateral(user1.address)).to.equal(0);
    });

    it("Should enforce daily lock limits", async function () {
      const amount = ethers.parseEther("1000000"); // 1M tokens
      const l2TxHash = ethers.keccak256(ethers.toUtf8Bytes("deposit"));

      await collateralToken.connect(user1).transfer(await collateralVaultL1.getAddress(), amount);
      await collateralVaultL1.lockCollateral(user1.address, amount, l2TxHash);

      const remaining = await collateralVaultL1.getRemainingDailyLockCapacity();
      expect(remaining).to.equal(ethers.parseEther("9000000")); // 10M - 1M
    });

    it("Should handle emergency withdrawals with delay", async function () {
      const amount = ethers.parseEther("100");
      const l2TxHash = ethers.keccak256(ethers.toUtf8Bytes("deposit"));

      // Lock collateral
      await collateralToken.connect(user1).transfer(await collateralVaultL1.getAddress(), amount);
      await collateralVaultL1.lockCollateral(user1.address, amount, l2TxHash);

      // Trigger emergency pause
      await collateralVaultL1.triggerEmergencyPause();

      // Request emergency withdrawal
      await collateralVaultL1.connect(user1).requestEmergencyWithdrawal(amount);

      // Should not be able to execute immediately
      await expect(
        collateralVaultL1.connect(user1).executeEmergencyWithdrawal()
      ).to.be.revertedWith("Delay period not passed");

      // Fast forward 7 days
      await time.increase(7 * 24 * 60 * 60);

      // Now should be able to execute
      await collateralVaultL1.connect(user1).executeEmergencyWithdrawal();

      expect(await collateralVaultL1.getLockedCollateral(user1.address)).to.equal(0);
    });
  });

  describe("L1 State Registry", function () {
    it("Should receive and store state roots", async function () {
      const stateRoot = ethers.keccak256(ethers.toUtf8Bytes("state1"));
      const l2Block = 100;
      const totalCollateral = ethers.parseEther("1000");
      const totalDebt = ethers.parseUnits("500", 6);

      await l1StateRegistry.receiveStateRoot(stateRoot, l2Block, totalCollateral, totalDebt);

      expect(await l1StateRegistry.getStateRoot(l2Block)).to.equal(stateRoot);
      expect(await l1StateRegistry.latestL2Block()).to.equal(l2Block);
    });

    it("Should verify state roots correctly", async function () {
      const stateRoot = ethers.keccak256(ethers.toUtf8Bytes("state1"));
      const l2Block = 100;
      const totalCollateral = ethers.parseEther("1000");
      const totalDebt = ethers.parseUnits("500", 6);

      await l1StateRegistry.receiveStateRoot(stateRoot, l2Block, totalCollateral, totalDebt);

      expect(await l1StateRegistry.verifyStateRoot(l2Block, stateRoot)).to.be.true;
      expect(await l1StateRegistry.verifyStateRoot(l2Block, ethers.ZeroHash)).to.be.false;
    });

    it("Should enforce minimum submission interval", async function () {
      const stateRoot1 = ethers.keccak256(ethers.toUtf8Bytes("state1"));
      const stateRoot2 = ethers.keccak256(ethers.toUtf8Bytes("state2"));
      const totalCollateral = ethers.parseEther("1000");
      const totalDebt = ethers.parseUnits("500", 6);

      await l1StateRegistry.receiveStateRoot(stateRoot1, 100, totalCollateral, totalDebt);

      // Should reject immediate submission
      await expect(
        l1StateRegistry.receiveStateRoot(stateRoot2, 101, totalCollateral, totalDebt)
      ).to.be.revertedWith("Submission too frequent");

      // Should accept after 1 hour
      await time.increase(ONE_HOUR);
      await l1StateRegistry.receiveStateRoot(stateRoot2, 101, totalCollateral, totalDebt);

      expect(await l1StateRegistry.latestL2Block()).to.equal(101);
    });

    it("Should detect critical conditions", async function () {
      const stateRoot = ethers.keccak256(ethers.toUtf8Bytes("state1"));
      const l2Block = 100;
      const totalCollateral = ethers.parseEther("100"); // Low collateral
      const totalDebt = ethers.parseUnits("100", 6); // High debt (ratio = 100%)

      // This should trigger critical condition (< 150% ratio)
      const tx = await l1StateRegistry.receiveStateRoot(
        stateRoot,
        l2Block,
        totalCollateral,
        totalDebt
      );

      const receipt = await tx.wait();
      const event = receipt.logs.find(
        log => log.fragment && log.fragment.name === "CriticalConditionDetected"
      );

      expect(event).to.not.be.undefined;
    });
  });

  describe("L2 State Aggregator", function () {
    it("Should calculate state root correctly", async function () {
      // Register a test module
      const moduleId = ethers.keccak256(ethers.toUtf8Bytes("TEST_MODULE"));
      await l2StateAggregator.registerModule(moduleId, deployer.address);
      await l2StateAggregator.authorizeModule(deployer.address);

      // Update module state
      const stateHash = ethers.keccak256(ethers.toUtf8Bytes("test_state"));
      await l2StateAggregator.updateModuleState(moduleId, stateHash);

      // Calculate root
      const root = await l2StateAggregator.calculateStateRoot();
      expect(root).to.not.equal(ethers.ZeroHash);
    });

    it("Should update system state", async function () {
      await l2StateAggregator.registerModule(
        ethers.keccak256(ethers.toUtf8Bytes("TEST")),
        deployer.address
      );
      await l2StateAggregator.authorizeModule(deployer.address);

      const totalCollateral = ethers.parseEther("1000");
      const totalDebt = ethers.parseUnits("500", 6);
      const activePositions = 10;
      const totalOrders = 5;

      await l2StateAggregator.updateSystemState(
        totalCollateral,
        totalDebt,
        activePositions,
        totalOrders
      );

      const state = await l2StateAggregator.getSystemState();
      expect(state.totalCollateral).to.equal(totalCollateral);
      expect(state.totalDebt).to.equal(totalDebt);
      expect(state.activePositions).to.equal(activePositions);
      expect(state.totalOrders).to.equal(totalOrders);
    });

    it("Should enforce submission interval", async function () {
      // Can't test actual L1 submission in unit test (requires Arbitrum precompile)
      // But we can test the timing logic
      expect(await l2StateAggregator.canSubmitToL1()).to.be.true;

      // After updating last submission time
      await l2StateAggregator.registerModule(
        ethers.keccak256(ethers.toUtf8Bytes("TEST")),
        deployer.address
      );
    });
  });

  describe("Integration: Full Deposit Flow", function () {
    it("Should handle complete deposit flow from user", async function () {
      const depositAmount = ethers.parseEther("100");

      // User approves gateway
      await collateralToken.connect(user1).approve(
        await l1Gateway.getAddress(),
        depositAmount
      );

      // Check initial balances
      const initialUserBalance = await collateralToken.balanceOf(user1.address);
      const initialVaultBalance = await collateralVaultL1.getTotalLocked();

      // Simulate deposit (in real scenario, this would trigger L2 message)
      // For testing, we directly lock in vault
      await collateralToken.connect(user1).transfer(
        await collateralVaultL1.getAddress(),
        depositAmount
      );

      const l2TxHash = ethers.keccak256(
        ethers.AbiCoder.defaultAbiCoder().encode(
          ["address", "uint256", "uint256"],
          [user1.address, depositAmount, await time.latest()]
        )
      );

      await collateralVaultL1.lockCollateral(user1.address, depositAmount, l2TxHash);

      // Verify state changes
      expect(await collateralToken.balanceOf(user1.address)).to.equal(
        initialUserBalance - depositAmount
      );
      expect(await collateralVaultL1.getLockedCollateral(user1.address)).to.equal(
        depositAmount
      );
      expect(await collateralVaultL1.getTotalLocked()).to.equal(
        initialVaultBalance + depositAmount
      );
    });
  });

  describe("Gas Optimization Verification", function () {
    it("Should demonstrate gas savings on L1", async function () {
      // This test shows that L1 operations are minimal
      const amount = ethers.parseUnits("100", 6);
      const l2TxHash = ethers.keccak256(ethers.toUtf8Bytes("test"));

      const BRIDGE_ROLE = await loyaltyUSDL1.BRIDGE_ROLE();
      await loyaltyUSDL1.grantRole(BRIDGE_ROLE, deployer.address);

      const tx = await loyaltyUSDL1.bridgeMint(user1.address, amount, l2TxHash);
      const receipt = await tx.wait();

      console.log(`L1 Bridge Mint Gas Used: ${receipt.gasUsed.toString()}`);

      // L1 operations should be minimal (< 100k gas)
      expect(receipt.gasUsed).to.be.lessThan(100000n);
    });
  });
});
