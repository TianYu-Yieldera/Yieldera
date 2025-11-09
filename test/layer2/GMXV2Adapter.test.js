const { expect } = require("chai");
const { ethers } = require("hardhat");
const { loadFixture } = require("@nomicfoundation/hardhat-network-helpers");

describe("GMXV2Adapter", function () {
  // ============ Fixture ============
  async function deployGMXV2AdapterFixture() {
    const [owner, user1, user2, riskManager] = await ethers.getSigners();

    // 部署 Mock GMX V2 合约
    const MockGMXV2 = await ethers.getContractFactory("MockGMXV2ExchangeRouter");
    const mockExchangeRouter = await MockGMXV2.deploy();

    const MockReader = await ethers.getContractFactory("MockGMXV2Reader");
    const mockReader = await MockReader.deploy();

    const MockDataStore = await ethers.getContractFactory("MockGMXV2DataStore");
    const mockDataStore = await MockDataStore.deploy();

    // 部署 GMXV2Adapter
    const GMXV2Adapter = await ethers.getContractFactory("GMXV2Adapter");
    const adapter = await GMXV2Adapter.deploy(
      await mockExchangeRouter.getAddress(),
      await mockReader.getAddress(),
      await mockDataStore.getAddress()
    );

    // 部署测试代币
    const MockERC20 = await ethers.getContractFactory("MockERC20");
    const usdc = await MockERC20.deploy("USD Coin", "USDC", 6);
    const weth = await MockERC20.deploy("Wrapped ETH", "WETH", 18);

    // 创建测试市场地址
    const mockMarket = "0x70d95587d40A2caf56bd97485aB3Eec10Bee6336"; // 示例市场地址

    // 配置 Adapter
    await adapter.addMarket(mockMarket);
    await adapter.addCollateral(await usdc.getAddress());
    await adapter.addCollateral(await weth.getAddress());

    // 授予 riskManager 角色
    const RISK_MANAGER_ROLE = await adapter.RISK_MANAGER_ROLE();
    await adapter.grantRole(RISK_MANAGER_ROLE, riskManager.address);

    // 给用户铸造代币
    await usdc.mint(user1.address, ethers.parseUnits("100000", 6)); // 100k USDC
    await weth.mint(user1.address, ethers.parseEther("100")); // 100 WETH

    return {
      adapter,
      mockExchangeRouter,
      mockReader,
      mockDataStore,
      usdc,
      weth,
      mockMarket,
      owner,
      user1,
      user2,
      riskManager,
    };
  }

  // ============ 部署测试 ============
  describe("Deployment", function () {
    it("Should set the correct GMX V2 contracts", async function () {
      const { adapter, mockExchangeRouter, mockReader, mockDataStore } =
        await loadFixture(deployGMXV2AdapterFixture);

      expect(await adapter.exchangeRouter()).to.equal(
        await mockExchangeRouter.getAddress()
      );
      expect(await adapter.reader()).to.equal(await mockReader.getAddress());
      expect(await adapter.dataStore()).to.equal(
        await mockDataStore.getAddress()
      );
    });

    it("Should grant admin role to deployer", async function () {
      const { adapter, owner } = await loadFixture(deployGMXV2AdapterFixture);

      const DEFAULT_ADMIN_ROLE = await adapter.DEFAULT_ADMIN_ROLE();
      expect(await adapter.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be
        .true;
    });

    it("Should initialize with correct constants", async function () {
      const { adapter } = await loadFixture(deployGMXV2AdapterFixture);

      expect(await adapter.MAX_LEVERAGE()).to.equal(50);
      expect(await adapter.MIN_EXECUTION_FEE()).to.equal(
        ethers.parseEther("0.0001")
      );
      expect(await adapter.MAX_SLIPPAGE_BPS()).to.equal(200);
    });
  });

  // ============ 市场和抵押品管理 ============
  describe("Market and Collateral Management", function () {
    it("Should allow admin to add markets", async function () {
      const { adapter, owner } = await loadFixture(deployGMXV2AdapterFixture);

      const newMarket = "0x1234567890123456789012345678901234567890";
      await expect(adapter.addMarket(newMarket))
        .to.emit(adapter, "MarketAdded")
        .withArgs(newMarket);

      expect(await adapter.supportedMarkets(newMarket)).to.be.true;
    });

    it("Should allow admin to add collateral", async function () {
      const { adapter, usdc } = await loadFixture(deployGMXV2AdapterFixture);

      const newToken = "0x9876543210987654321098765432109876543210";
      await expect(adapter.addCollateral(newToken))
        .to.emit(adapter, "CollateralAdded")
        .withArgs(newToken);

      expect(await adapter.supportedCollateral(newToken)).to.be.true;
    });

    it("Should reject market addition by non-admin", async function () {
      const { adapter, user1 } = await loadFixture(deployGMXV2AdapterFixture);

      const newMarket = "0x1234567890123456789012345678901234567890";
      await expect(
        adapter.connect(user1).addMarket(newMarket)
      ).to.be.revertedWithCustomError(adapter, "AccessControlUnauthorizedAccount");
    });
  });

  // ============ 开仓功能 ============
  describe("Open Position", function () {
    it("Should open a long position successfully", async function () {
      const { adapter, usdc, mockMarket, user1 } = await loadFixture(
        deployGMXV2AdapterFixture
      );

      const collateralAmount = ethers.parseUnits("1000", 6); // 1000 USDC
      const sizeInUsd = ethers.parseEther("10000"); // 10k USD
      const executionFee = ethers.parseEther("0.001");

      // 批准 adapter 使用 USDC
      await usdc.connect(user1).approve(await adapter.getAddress(), collateralAmount);

      // 开多单
      const tx = await adapter.connect(user1).openPosition(
        mockMarket,
        await usdc.getAddress(),
        collateralAmount,
        sizeInUsd,
        true, // isLong
        ethers.parseEther("2000"), // acceptablePrice
        executionFee,
        { value: executionFee }
      );

      // 验证事件
      await expect(tx)
        .to.emit(adapter, "PositionOpened")
        .withArgs(
          user1.address,
          ethers.AnyValue, // orderKey
          mockMarket,
          await usdc.getAddress(),
          true, // isLong
          sizeInUsd,
          collateralAmount,
          10, // leverage (10x)
          false // isHedge
        );

      // 验证统计数据
      const stats = await adapter.getStatistics();
      expect(stats.totalOrders).to.equal(1);
      expect(stats.totalVolume).to.equal(sizeInUsd);
    });

    it("Should open a short position successfully", async function () {
      const { adapter, usdc, mockMarket, user1 } = await loadFixture(
        deployGMXV2AdapterFixture
      );

      const collateralAmount = ethers.parseUnits("2000", 6); // 2000 USDC
      const sizeInUsd = ethers.parseEther("20000"); // 20k USD
      const executionFee = ethers.parseEther("0.001");

      await usdc.connect(user1).approve(await adapter.getAddress(), collateralAmount);

      const tx = await adapter.connect(user1).openPosition(
        mockMarket,
        await usdc.getAddress(),
        collateralAmount,
        sizeInUsd,
        false, // isLong = false (空单)
        ethers.parseEther("1800"),
        executionFee,
        { value: executionFee }
      );

      await expect(tx).to.emit(adapter, "PositionOpened");

      // 验证用户仓位记录
      const positions = await adapter.getUserPositions(user1.address);
      expect(positions.length).to.equal(1);
      expect(positions[0].isLong).to.be.false;
      expect(positions[0].leverage).to.equal(10);
    });

    it("Should reject position with excessive leverage", async function () {
      const { adapter, usdc, mockMarket, user1 } = await loadFixture(
        deployGMXV2AdapterFixture
      );

      const collateralAmount = ethers.parseUnits("100", 6); // 100 USDC
      const sizeInUsd = ethers.parseEther("10000"); // 10k USD (100x 杠杆)
      const executionFee = ethers.parseEther("0.001");

      await usdc.connect(user1).approve(await adapter.getAddress(), collateralAmount);

      await expect(
        adapter.connect(user1).openPosition(
          mockMarket,
          await usdc.getAddress(),
          collateralAmount,
          sizeInUsd,
          true,
          ethers.parseEther("2000"),
          executionFee,
          { value: executionFee }
        )
      ).to.be.revertedWithCustomError(adapter, "InvalidLeverage");
    });

    it("Should reject unsupported market", async function () {
      const { adapter, usdc, user1 } = await loadFixture(
        deployGMXV2AdapterFixture
      );

      const unsupportedMarket = "0x0000000000000000000000000000000000000001";
      const collateralAmount = ethers.parseUnits("1000", 6);
      const executionFee = ethers.parseEther("0.001");

      await usdc.connect(user1).approve(await adapter.getAddress(), collateralAmount);

      await expect(
        adapter.connect(user1).openPosition(
          unsupportedMarket,
          await usdc.getAddress(),
          collateralAmount,
          ethers.parseEther("10000"),
          true,
          ethers.parseEther("2000"),
          executionFee,
          { value: executionFee }
        )
      ).to.be.revertedWithCustomError(adapter, "UnsupportedMarket");
    });

    it("Should reject insufficient execution fee", async function () {
      const { adapter, usdc, mockMarket, user1 } = await loadFixture(
        deployGMXV2AdapterFixture
      );

      const collateralAmount = ethers.parseUnits("1000", 6);
      const lowExecutionFee = ethers.parseEther("0.00001"); // 太低

      await usdc.connect(user1).approve(await adapter.getAddress(), collateralAmount);

      await expect(
        adapter.connect(user1).openPosition(
          mockMarket,
          await usdc.getAddress(),
          collateralAmount,
          ethers.parseEther("10000"),
          true,
          ethers.parseEther("2000"),
          lowExecutionFee,
          { value: lowExecutionFee }
        )
      ).to.be.revertedWithCustomError(adapter, "InsufficientExecutionFee");
    });
  });

  // ============ 平仓功能 ============
  describe("Close Position", function () {
    it("Should close position successfully", async function () {
      const { adapter, usdc, mockMarket, user1 } = await loadFixture(
        deployGMXV2AdapterFixture
      );

      // 先开仓
      const collateralAmount = ethers.parseUnits("1000", 6);
      const sizeInUsd = ethers.parseEther("10000");
      const executionFee = ethers.parseEther("0.001");

      await usdc.connect(user1).approve(await adapter.getAddress(), collateralAmount);

      await adapter.connect(user1).openPosition(
        mockMarket,
        await usdc.getAddress(),
        collateralAmount,
        sizeInUsd,
        true,
        ethers.parseEther("2000"),
        executionFee,
        { value: executionFee }
      );

      // 平仓
      const closeTx = await adapter.connect(user1).closePosition(
        mockMarket,
        await usdc.getAddress(),
        sizeInUsd,
        true, // isLong
        ethers.parseEther("2100"),
        executionFee,
        { value: executionFee }
      );

      await expect(closeTx)
        .to.emit(adapter, "PositionClosed")
        .withArgs(
          user1.address,
          ethers.AnyValue,
          mockMarket,
          sizeInUsd,
          0
        );
    });

    it("Should reject closing non-existent position", async function () {
      const { adapter, usdc, mockMarket, user1 } = await loadFixture(
        deployGMXV2AdapterFixture
      );

      const executionFee = ethers.parseEther("0.001");

      await expect(
        adapter.connect(user1).closePosition(
          mockMarket,
          await usdc.getAddress(),
          ethers.parseEther("10000"),
          true,
          ethers.parseEther("2000"),
          executionFee,
          { value: executionFee }
        )
      ).to.be.revertedWithCustomError(adapter, "NoPositionFound");
    });
  });

  // ============ 紧急对冲功能 ============
  describe("Emergency Hedge", function () {
    it("Should execute emergency hedge by risk manager", async function () {
      const { adapter, usdc, mockMarket, user1, riskManager } =
        await loadFixture(deployGMXV2AdapterFixture);

      // 先让用户开一个多单
      const collateralAmount = ethers.parseUnits("2000", 6);
      const sizeInUsd = ethers.parseEther("20000");
      const executionFee = ethers.parseEther("0.001");

      await usdc.connect(user1).approve(await adapter.getAddress(), collateralAmount * 2n);

      await adapter.connect(user1).openPosition(
        mockMarket,
        await usdc.getAddress(),
        collateralAmount,
        sizeInUsd,
        true, // 多单
        ethers.parseEther("2000"),
        executionFee,
        { value: executionFee }
      );

      // Risk Manager 执行对冲（开空单）
      const hedgeSize = ethers.parseEther("10000"); // 对冲一半
      const hedgeCollateral = hedgeSize / 10n; // 10x 杠杆

      const hedgeTx = await adapter.connect(riskManager).emergencyHedge(
        user1.address,
        mockMarket,
        await usdc.getAddress(),
        hedgeSize,
        "Liquidation protection",
        { value: executionFee }
      );

      await expect(hedgeTx)
        .to.emit(adapter, "EmergencyHedgeExecuted")
        .withArgs(
          user1.address,
          mockMarket,
          hedgeSize,
          "Liquidation protection",
          ethers.AnyValue
        );

      // 验证统计
      const stats = await adapter.getStatistics();
      expect(stats.totalHedges).to.equal(1);

      // 验证对冲仓位
      const positions = await adapter.getUserPositions(user1.address);
      expect(positions.length).to.equal(2); // 原仓位 + 对冲仓位
      expect(positions[1].isHedge).to.be.true;
      expect(positions[1].isLong).to.be.false; // 对冲是空单
    });

    it("Should reject hedge by unauthorized user", async function () {
      const { adapter, usdc, mockMarket, user1, user2 } =
        await loadFixture(deployGMXV2AdapterFixture);

      const executionFee = ethers.parseEther("0.001");

      await expect(
        adapter.connect(user2).emergencyHedge(
          user1.address,
          mockMarket,
          await usdc.getAddress(),
          ethers.parseEther("10000"),
          "Unauthorized hedge",
          { value: executionFee }
        )
      ).to.be.revertedWithCustomError(adapter, "AccessControlUnauthorizedAccount");
    });
  });

  // ============ 查询功能 ============
  describe("Query Functions", function () {
    it("Should return user positions", async function () {
      const { adapter, usdc, mockMarket, user1 } = await loadFixture(
        deployGMXV2AdapterFixture
      );

      // 开两个仓位
      const collateralAmount = ethers.parseUnits("1000", 6);
      const executionFee = ethers.parseEther("0.001");

      await usdc.connect(user1).approve(await adapter.getAddress(), collateralAmount * 2n);

      await adapter.connect(user1).openPosition(
        mockMarket,
        await usdc.getAddress(),
        collateralAmount,
        ethers.parseEther("10000"),
        true,
        ethers.parseEther("2000"),
        executionFee,
        { value: executionFee }
      );

      await adapter.connect(user1).openPosition(
        mockMarket,
        await usdc.getAddress(),
        collateralAmount,
        ethers.parseEther("5000"),
        false,
        ethers.parseEther("1800"),
        executionFee,
        { value: executionFee }
      );

      const positions = await adapter.getUserPositions(user1.address);
      expect(positions.length).to.equal(2);
      expect(positions[0].isLong).to.be.true;
      expect(positions[1].isLong).to.be.false;
    });

    it("Should return statistics", async function () {
      const { adapter, usdc, mockMarket, user1 } = await loadFixture(
        deployGMXV2AdapterFixture
      );

      const collateralAmount = ethers.parseUnits("1000", 6);
      const sizeInUsd = ethers.parseEther("10000");
      const executionFee = ethers.parseEther("0.001");

      await usdc.connect(user1).approve(await adapter.getAddress(), collateralAmount);

      await adapter.connect(user1).openPosition(
        mockMarket,
        await usdc.getAddress(),
        collateralAmount,
        sizeInUsd,
        true,
        ethers.parseEther("2000"),
        executionFee,
        { value: executionFee }
      );

      const stats = await adapter.getStatistics();
      expect(stats.totalOrders).to.equal(1);
      expect(stats.totalVolume).to.equal(sizeInUsd);
      expect(stats.totalHedges).to.equal(0);
    });
  });

  // ============ 暂停功能 ============
  describe("Pause Functionality", function () {
    it("Should pause and unpause by admin", async function () {
      const { adapter, owner } = await loadFixture(deployGMXV2AdapterFixture);

      await adapter.pause();
      expect(await adapter.paused()).to.be.true;

      await adapter.unpause();
      expect(await adapter.paused()).to.be.false;
    });

    it("Should reject operations when paused", async function () {
      const { adapter, usdc, mockMarket, user1 } = await loadFixture(
        deployGMXV2AdapterFixture
      );

      await adapter.pause();

      const collateralAmount = ethers.parseUnits("1000", 6);
      const executionFee = ethers.parseEther("0.001");

      await usdc.connect(user1).approve(await adapter.getAddress(), collateralAmount);

      await expect(
        adapter.connect(user1).openPosition(
          mockMarket,
          await usdc.getAddress(),
          collateralAmount,
          ethers.parseEther("10000"),
          true,
          ethers.parseEther("2000"),
          executionFee,
          { value: executionFee }
        )
      ).to.be.revertedWithCustomError(adapter, "EnforcedPause");
    });
  });
});
