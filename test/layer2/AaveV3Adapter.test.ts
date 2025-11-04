import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type {
  AaveV3Adapter,
  MockAaveV3Pool,
  MockAToken,
  MockFlashLoanReceiver,
  MockERC20,
  L2StateAggregator,
} from "../../typechain-types";

describe("AaveV3Adapter", function () {
  let adapter: AaveV3Adapter;
  let aavePool: MockAaveV3Pool;
  let stateAggregator: L2StateAggregator;
  let token: MockERC20;
  let aToken: MockAToken;
  let collateralToken: MockERC20;
  let aCollateralToken: MockAToken;
  let flashLoanReceiver: MockFlashLoanReceiver;

  let owner: SignerWithAddress;
  let user1: SignerWithAddress;
  let user2: SignerWithAddress;

  const SUPPLY_AMOUNT = ethers.parseEther("1000");
  const BORROW_AMOUNT = ethers.parseEther("500");
  const INTEREST_RATE_MODE_STABLE = 1;
  const INTEREST_RATE_MODE_VARIABLE = 2;

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();

    // Deploy mock ERC20 tokens
    const MockERC20Factory = await ethers.getContractFactory("MockERC20");
    token = await MockERC20Factory.deploy(
      "Test Token",
      "TEST",
      ethers.parseEther("1000000")
    );

    collateralToken = await MockERC20Factory.deploy(
      "Collateral Token",
      "COLL",
      ethers.parseEther("1000000")
    );

    // Deploy mock Aave V3 Pool
    const MockAaveV3PoolFactory = await ethers.getContractFactory("MockAaveV3Pool");
    aavePool = await MockAaveV3PoolFactory.deploy();

    // Deploy mock aTokens
    const MockATokenFactory = await ethers.getContractFactory("MockAToken");
    aToken = await MockATokenFactory.deploy("aToken", "aTEST");
    aCollateralToken = await MockATokenFactory.deploy("aCollateral", "aCOLL");

    // Set pool in aTokens
    await aToken.setPool(await aavePool.getAddress());
    await aCollateralToken.setPool(await aavePool.getAddress());

    // Set aTokens in pool
    await aavePool.setAToken(await token.getAddress(), await aToken.getAddress());
    await aavePool.setAToken(
      await collateralToken.getAddress(),
      await aCollateralToken.getAddress()
    );

    // Deploy L2StateAggregator
    const L2StateAggregatorFactory = await ethers.getContractFactory("L2StateAggregator");
    stateAggregator = await L2StateAggregatorFactory.deploy(await owner.getAddress());

    // Deploy AaveV3Adapter
    const AaveV3AdapterFactory = await ethers.getContractFactory("AaveV3Adapter");
    adapter = await AaveV3AdapterFactory.deploy(
      await aavePool.getAddress(),
      await stateAggregator.getAddress(),
      await owner.getAddress()
    );

    // Authorize adapter in state aggregator
    await stateAggregator.authorizeModule(await adapter.getAddress());

    // Fund pool with tokens for lending
    const poolAddress = await aavePool.getAddress();
    await token.transfer(poolAddress, ethers.parseEther("10000"));
    await collateralToken.transfer(poolAddress, ethers.parseEther("10000"));

    // Fund users
    await token.transfer(await user1.getAddress(), ethers.parseEther("5000"));
    await token.transfer(await user2.getAddress(), ethers.parseEther("5000"));
    await collateralToken.transfer(await user1.getAddress(), ethers.parseEther("5000"));
    await collateralToken.transfer(await user2.getAddress(), ethers.parseEther("5000"));

    // Approve adapter to spend user tokens
    await token
      .connect(user1)
      .approve(await adapter.getAddress(), ethers.MaxUint256);
    await token
      .connect(user2)
      .approve(await adapter.getAddress(), ethers.MaxUint256);
    await collateralToken
      .connect(user1)
      .approve(await adapter.getAddress(), ethers.MaxUint256);
    await collateralToken
      .connect(user2)
      .approve(await adapter.getAddress(), ethers.MaxUint256);

    // Deploy flash loan receiver
    const MockFlashLoanReceiverFactory = await ethers.getContractFactory(
      "MockFlashLoanReceiver"
    );
    flashLoanReceiver = await MockFlashLoanReceiverFactory.deploy(
      await adapter.getAddress()
    );
  });

  describe("Deployment", function () {
    it("Should set the correct Aave Pool address", async function () {
      expect(await adapter.aavePool()).to.equal(await aavePool.getAddress());
    });

    it("Should set the correct state aggregator", async function () {
      expect(await adapter.stateAggregator()).to.equal(
        await stateAggregator.getAddress()
      );
    });

    it("Should set the correct owner", async function () {
      expect(await adapter.owner()).to.equal(await owner.getAddress());
    });

    it("Should have correct interest rate mode constants", async function () {
      expect(await adapter.INTEREST_RATE_MODE_STABLE()).to.equal(1);
      expect(await adapter.INTEREST_RATE_MODE_VARIABLE()).to.equal(2);
    });

    it("Should initialize with zero statistics", async function () {
      expect(await adapter.totalSupplied()).to.equal(0);
      expect(await adapter.totalBorrowed()).to.equal(0);
      expect(await adapter.activeUsers()).to.equal(0);
    });

    it("Should revert with zero Aave Pool address", async function () {
      const AaveV3AdapterFactory = await ethers.getContractFactory("AaveV3Adapter");
      await expect(
        AaveV3AdapterFactory.deploy(
          ethers.ZeroAddress,
          await stateAggregator.getAddress(),
          await owner.getAddress()
        )
      ).to.be.revertedWith("Invalid Aave Pool address");
    });
  });

  describe("Supply", function () {
    it("Should allow users to supply assets", async function () {
      await expect(
        adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT)
      )
        .to.emit(adapter, "Supplied")
        .withArgs(
          await user1.getAddress(),
          await token.getAddress(),
          SUPPLY_AMOUNT,
          await ethers.provider.getBlock("latest").then((b) => b!.timestamp + 1)
        );

      const position = await adapter.userPositions(await user1.getAddress());
      expect(position.totalSupplied).to.equal(SUPPLY_AMOUNT);
    });

    it("Should transfer tokens from user to pool", async function () {
      const user1Address = await user1.getAddress();
      const poolAddress = await aavePool.getAddress();
      const initialUserBalance = await token.balanceOf(user1Address);
      const initialPoolBalance = await token.balanceOf(poolAddress);

      await adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT);

      expect(await token.balanceOf(user1Address)).to.equal(
        initialUserBalance - SUPPLY_AMOUNT
      );
      expect(await token.balanceOf(poolAddress)).to.equal(
        initialPoolBalance + SUPPLY_AMOUNT
      );
    });

    it("Should mint aTokens to user", async function () {
      const user1Address = await user1.getAddress();

      await adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT);

      expect(await aToken.balanceOf(user1Address)).to.equal(SUPPLY_AMOUNT);
    });

    it("Should update global statistics", async function () {
      await adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT);

      expect(await adapter.totalSupplied()).to.equal(SUPPLY_AMOUNT);
      expect(await adapter.activeUsers()).to.equal(1);
    });

    it("Should update state aggregator", async function () {
      await adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT);

      const [totalCollateral, totalDebt, activePositions] = await stateAggregator.getSystemState();
      expect(totalCollateral).to.equal(SUPPLY_AMOUNT);
      expect(activePositions).to.equal(1);
    });

    it("Should track multiple supplies from same user", async function () {
      await adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT);
      await adapter
        .connect(user1)
        .supply(await token.getAddress(), SUPPLY_AMOUNT / 2n);

      const position = await adapter.userPositions(await user1.getAddress());
      expect(position.totalSupplied).to.equal(SUPPLY_AMOUNT + SUPPLY_AMOUNT / 2n);
      expect(await adapter.totalSupplied()).to.equal(
        SUPPLY_AMOUNT + SUPPLY_AMOUNT / 2n
      );
    });

    it("Should track supplies from multiple users", async function () {
      await adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT);
      await adapter.connect(user2).supply(await token.getAddress(), SUPPLY_AMOUNT);

      expect(await adapter.activeUsers()).to.equal(2);
      expect(await adapter.totalSupplied()).to.equal(SUPPLY_AMOUNT * 2n);
    });

    it("Should update timestamp", async function () {
      await adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT);

      const position = await adapter.userPositions(await user1.getAddress());
      const currentBlock = await ethers.provider.getBlock("latest");
      expect(position.lastUpdate).to.equal(currentBlock!.timestamp);
    });

    it("Should revert if amount is zero", async function () {
      await expect(
        adapter.connect(user1).supply(await token.getAddress(), 0)
      ).to.be.revertedWith("Amount must be greater than zero");
    });

    it("Should revert if user has insufficient balance", async function () {
      const largeAmount = ethers.parseEther("10000");
      await expect(
        adapter.connect(user1).supply(await token.getAddress(), largeAmount)
      ).to.be.reverted;
    });
  });

  describe("Withdraw", function () {
    beforeEach(async function () {
      // Supply first so we can withdraw
      await adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT);
    });

    it("Should allow users to withdraw supplied assets", async function () {
      const withdrawAmount = SUPPLY_AMOUNT / 2n;

      await expect(
        adapter.connect(user1).withdraw(await token.getAddress(), withdrawAmount)
      )
        .to.emit(adapter, "Withdrawn")
        .withArgs(
          await user1.getAddress(),
          await token.getAddress(),
          withdrawAmount,
          await ethers.provider.getBlock("latest").then((b) => b!.timestamp + 1)
        );

      const position = await adapter.userPositions(await user1.getAddress());
      expect(position.totalSupplied).to.equal(SUPPLY_AMOUNT - withdrawAmount);
    });

    it("Should transfer tokens back to user", async function () {
      const user1Address = await user1.getAddress();
      const balanceBefore = await token.balanceOf(user1Address);

      const withdrawAmount = SUPPLY_AMOUNT / 2n;
      await adapter.connect(user1).withdraw(await token.getAddress(), withdrawAmount);

      expect(await token.balanceOf(user1Address)).to.equal(
        balanceBefore + withdrawAmount
      );
    });

    it("Should burn aTokens from user", async function () {
      const user1Address = await user1.getAddress();
      const aTokenBalanceBefore = await aToken.balanceOf(user1Address);

      const withdrawAmount = SUPPLY_AMOUNT / 2n;
      await adapter.connect(user1).withdraw(await token.getAddress(), withdrawAmount);

      expect(await aToken.balanceOf(user1Address)).to.equal(
        aTokenBalanceBefore - withdrawAmount
      );
    });

    it("Should update global statistics", async function () {
      const withdrawAmount = SUPPLY_AMOUNT / 2n;
      await adapter.connect(user1).withdraw(await token.getAddress(), withdrawAmount);

      expect(await adapter.totalSupplied()).to.equal(SUPPLY_AMOUNT - withdrawAmount);
    });

    it("Should update state aggregator", async function () {
      const withdrawAmount = SUPPLY_AMOUNT / 2n;
      await adapter.connect(user1).withdraw(await token.getAddress(), withdrawAmount);

      const [totalCollateral] = await stateAggregator.getSystemState();
      expect(totalCollateral).to.equal(SUPPLY_AMOUNT - withdrawAmount);
    });

    it("Should allow full withdrawal with type(uint256).max", async function () {
      await adapter
        .connect(user1)
        .withdraw(await token.getAddress(), ethers.MaxUint256);

      const position = await adapter.userPositions(await user1.getAddress());
      expect(position.totalSupplied).to.equal(0);
    });

    it("Should deactivate user when fully withdrawn", async function () {
      await adapter
        .connect(user1)
        .withdraw(await token.getAddress(), ethers.MaxUint256);

      expect(await adapter.activeUsers()).to.equal(0);
    });

    it("Should update timestamp", async function () {
      await adapter.connect(user1).withdraw(await token.getAddress(), SUPPLY_AMOUNT);

      const position = await adapter.userPositions(await user1.getAddress());
      const currentBlock = await ethers.provider.getBlock("latest");
      expect(position.lastUpdate).to.equal(currentBlock!.timestamp);
    });

    it("Should revert if amount is zero", async function () {
      await expect(
        adapter.connect(user1).withdraw(await token.getAddress(), 0)
      ).to.be.revertedWith("Amount must be greater than zero");
    });

    it("Should revert if user has insufficient supply", async function () {
      const largeAmount = SUPPLY_AMOUNT * 2n;
      await expect(
        adapter.connect(user1).withdraw(await token.getAddress(), largeAmount)
      ).to.be.revertedWith("Insufficient supply");
    });
  });

  describe("Borrow", function () {
    beforeEach(async function () {
      // Supply collateral first
      await adapter
        .connect(user1)
        .supply(await collateralToken.getAddress(), SUPPLY_AMOUNT);
    });

    it("Should allow users to borrow with variable rate", async function () {
      await expect(
        adapter
          .connect(user1)
          .borrow(
            await token.getAddress(),
            BORROW_AMOUNT,
            INTEREST_RATE_MODE_VARIABLE
          )
      )
        .to.emit(adapter, "Borrowed")
        .withArgs(
          await user1.getAddress(),
          await token.getAddress(),
          BORROW_AMOUNT,
          INTEREST_RATE_MODE_VARIABLE,
          await ethers.provider.getBlock("latest").then((b) => b!.timestamp + 1)
        );

      const position = await adapter.userPositions(await user1.getAddress());
      expect(position.totalBorrowed).to.equal(BORROW_AMOUNT);
    });

    it("Should allow users to borrow with stable rate", async function () {
      await expect(
        adapter
          .connect(user1)
          .borrow(await token.getAddress(), BORROW_AMOUNT, INTEREST_RATE_MODE_STABLE)
      )
        .to.emit(adapter, "Borrowed")
        .withArgs(
          await user1.getAddress(),
          await token.getAddress(),
          BORROW_AMOUNT,
          INTEREST_RATE_MODE_STABLE,
          await ethers.provider.getBlock("latest").then((b) => b!.timestamp + 1)
        );
    });

    it("Should transfer borrowed tokens to user", async function () {
      const user1Address = await user1.getAddress();
      const balanceBefore = await token.balanceOf(user1Address);

      await adapter
        .connect(user1)
        .borrow(
          await token.getAddress(),
          BORROW_AMOUNT,
          INTEREST_RATE_MODE_VARIABLE
        );

      expect(await token.balanceOf(user1Address)).to.equal(
        balanceBefore + BORROW_AMOUNT
      );
    });

    it("Should update global statistics", async function () {
      await adapter
        .connect(user1)
        .borrow(
          await token.getAddress(),
          BORROW_AMOUNT,
          INTEREST_RATE_MODE_VARIABLE
        );

      expect(await adapter.totalBorrowed()).to.equal(BORROW_AMOUNT);
    });

    it("Should update state aggregator", async function () {
      await adapter
        .connect(user1)
        .borrow(
          await token.getAddress(),
          BORROW_AMOUNT,
          INTEREST_RATE_MODE_VARIABLE
        );

      const [, totalDebt] = await stateAggregator.getSystemState();
      expect(totalDebt).to.equal(BORROW_AMOUNT);
    });

    it("Should track multiple borrows", async function () {
      await adapter
        .connect(user1)
        .borrow(
          await token.getAddress(),
          BORROW_AMOUNT,
          INTEREST_RATE_MODE_VARIABLE
        );
      await adapter
        .connect(user1)
        .borrow(
          await token.getAddress(),
          BORROW_AMOUNT / 2n,
          INTEREST_RATE_MODE_VARIABLE
        );

      const position = await adapter.userPositions(await user1.getAddress());
      expect(position.totalBorrowed).to.equal(BORROW_AMOUNT + BORROW_AMOUNT / 2n);
    });

    it("Should update timestamp", async function () {
      await adapter
        .connect(user1)
        .borrow(
          await token.getAddress(),
          BORROW_AMOUNT,
          INTEREST_RATE_MODE_VARIABLE
        );

      const position = await adapter.userPositions(await user1.getAddress());
      const currentBlock = await ethers.provider.getBlock("latest");
      expect(position.lastUpdate).to.equal(currentBlock!.timestamp);
    });

    it("Should revert if amount is zero", async function () {
      await expect(
        adapter
          .connect(user1)
          .borrow(await token.getAddress(), 0, INTEREST_RATE_MODE_VARIABLE)
      ).to.be.revertedWith("Amount must be greater than zero");
    });

    it("Should revert with invalid interest rate mode", async function () {
      await expect(
        adapter.connect(user1).borrow(await token.getAddress(), BORROW_AMOUNT, 3)
      ).to.be.revertedWith("Invalid interest rate mode");
    });
  });

  describe("Repay", function () {
    beforeEach(async function () {
      // Supply collateral and borrow first
      await adapter
        .connect(user1)
        .supply(await collateralToken.getAddress(), SUPPLY_AMOUNT);
      await adapter
        .connect(user1)
        .borrow(
          await token.getAddress(),
          BORROW_AMOUNT,
          INTEREST_RATE_MODE_VARIABLE
        );
    });

    it("Should allow users to repay borrowed assets", async function () {
      const repayAmount = BORROW_AMOUNT / 2n;

      await expect(
        adapter
          .connect(user1)
          .repay(await token.getAddress(), repayAmount, INTEREST_RATE_MODE_VARIABLE)
      )
        .to.emit(adapter, "Repaid")
        .withArgs(
          await user1.getAddress(),
          await token.getAddress(),
          repayAmount,
          INTEREST_RATE_MODE_VARIABLE,
          await ethers.provider.getBlock("latest").then((b) => b!.timestamp + 1)
        );

      const position = await adapter.userPositions(await user1.getAddress());
      expect(position.totalBorrowed).to.equal(BORROW_AMOUNT - repayAmount);
    });

    it("Should transfer tokens from user to pool", async function () {
      const user1Address = await user1.getAddress();
      const balanceBefore = await token.balanceOf(user1Address);
      const repayAmount = BORROW_AMOUNT / 2n;

      await adapter
        .connect(user1)
        .repay(await token.getAddress(), repayAmount, INTEREST_RATE_MODE_VARIABLE);

      expect(await token.balanceOf(user1Address)).to.equal(
        balanceBefore - repayAmount
      );
    });

    it("Should update global statistics", async function () {
      const repayAmount = BORROW_AMOUNT / 2n;

      await adapter
        .connect(user1)
        .repay(await token.getAddress(), repayAmount, INTEREST_RATE_MODE_VARIABLE);

      expect(await adapter.totalBorrowed()).to.equal(BORROW_AMOUNT - repayAmount);
    });

    it("Should update state aggregator", async function () {
      const repayAmount = BORROW_AMOUNT / 2n;

      await adapter
        .connect(user1)
        .repay(await token.getAddress(), repayAmount, INTEREST_RATE_MODE_VARIABLE);

      const [, totalDebt] = await stateAggregator.getSystemState();
      expect(totalDebt).to.equal(BORROW_AMOUNT - repayAmount);
    });

    it("Should allow full repayment with exact borrowed amount", async function () {
      const positionBefore = await adapter.userPositions(await user1.getAddress());

      await adapter
        .connect(user1)
        .repay(
          await token.getAddress(),
          positionBefore.totalBorrowed,
          INTEREST_RATE_MODE_VARIABLE
        );

      const position = await adapter.userPositions(await user1.getAddress());
      expect(position.totalBorrowed).to.equal(0);
    });

    it("Should handle overpayment correctly", async function () {
      const user1Address = await user1.getAddress();
      const balanceBefore = await token.balanceOf(user1Address);
      const overpayAmount = BORROW_AMOUNT * 2n;

      await adapter
        .connect(user1)
        .repay(
          await token.getAddress(),
          overpayAmount,
          INTEREST_RATE_MODE_VARIABLE
        );

      const position = await adapter.userPositions(await user1.getAddress());
      expect(position.totalBorrowed).to.equal(0);

      // Should refund excess
      expect(await token.balanceOf(user1Address)).to.equal(
        balanceBefore - BORROW_AMOUNT
      );
    });

    it("Should deactivate user when fully repaid and no supply", async function () {
      // First withdraw collateral
      await adapter
        .connect(user1)
        .withdraw(await collateralToken.getAddress(), ethers.MaxUint256);

      // Then repay all (get actual borrowed amount)
      const position = await adapter.userPositions(await user1.getAddress());
      await adapter
        .connect(user1)
        .repay(
          await token.getAddress(),
          position.totalBorrowed,
          INTEREST_RATE_MODE_VARIABLE
        );

      expect(await adapter.activeUsers()).to.equal(0);
    });

    it("Should update timestamp", async function () {
      await adapter
        .connect(user1)
        .repay(await token.getAddress(), BORROW_AMOUNT, INTEREST_RATE_MODE_VARIABLE);

      const position = await adapter.userPositions(await user1.getAddress());
      const currentBlock = await ethers.provider.getBlock("latest");
      expect(position.lastUpdate).to.equal(currentBlock!.timestamp);
    });

    it("Should revert if amount is zero", async function () {
      await expect(
        adapter.connect(user1).repay(await token.getAddress(), 0, INTEREST_RATE_MODE_VARIABLE)
      ).to.be.revertedWith("Amount must be greater than zero");
    });
  });

  describe("Flash Loan Simple", function () {
    const FLASH_AMOUNT = ethers.parseEther("1000");
    const EXPECTED_PREMIUM = (FLASH_AMOUNT * 9n) / 10000n; // 0.09%

    beforeEach(async function () {
      // Fund flash loan receiver for repayment
      await token.transfer(
        await flashLoanReceiver.getAddress(),
        EXPECTED_PREMIUM
      );
    });

    it("Should execute simple flash loan successfully", async function () {
      const receiverAddress = await flashLoanReceiver.getAddress();
      const params = ethers.AbiCoder.defaultAbiCoder().encode(
        ["address"],
        [receiverAddress]
      );

      await expect(
        adapter
          .connect(user1)
          .flashLoanSimple(await token.getAddress(), FLASH_AMOUNT, params)
      )
        .to.emit(flashLoanReceiver, "SimpleFlashLoanReceived")
        .withArgs(await token.getAddress(), FLASH_AMOUNT, EXPECTED_PREMIUM);
    });

    it("Should emit FlashLoanExecuted event", async function () {
      const receiverAddress = await flashLoanReceiver.getAddress();
      const params = ethers.AbiCoder.defaultAbiCoder().encode(
        ["address"],
        [receiverAddress]
      );

      await expect(
        adapter
          .connect(user1)
          .flashLoanSimple(await token.getAddress(), FLASH_AMOUNT, params)
      )
        .to.emit(adapter, "FlashLoanExecuted")
        .withArgs(
          await user1.getAddress(),
          await token.getAddress(),
          FLASH_AMOUNT,
          EXPECTED_PREMIUM,
          await ethers.provider.getBlock("latest").then((b) => b!.timestamp + 1)
        );
    });

    it("Should calculate correct premium", async function () {
      const calculatedPremium = await adapter.calculateFlashLoanPremium(FLASH_AMOUNT);
      expect(calculatedPremium).to.equal(EXPECTED_PREMIUM);
    });

    it("Should revert if amount is zero", async function () {
      const receiverAddress = await flashLoanReceiver.getAddress();
      const params = ethers.AbiCoder.defaultAbiCoder().encode(
        ["address"],
        [receiverAddress]
      );

      await expect(
        adapter.connect(user1).flashLoanSimple(await token.getAddress(), 0, params)
      ).to.be.revertedWith("Amount must be greater than zero");
    });

    it("Should revert if receiver fails execution", async function () {
      await flashLoanReceiver.setShouldFail(true);

      const receiverAddress = await flashLoanReceiver.getAddress();
      const params = ethers.AbiCoder.defaultAbiCoder().encode(
        ["address"],
        [receiverAddress]
      );

      await expect(
        adapter
          .connect(user1)
          .flashLoanSimple(await token.getAddress(), FLASH_AMOUNT, params)
      ).to.be.revertedWith("Flash loan failed intentionally");
    });

    it("Should revert if receiver doesn't repay", async function () {
      await flashLoanReceiver.setShouldNotRepay(true);

      const receiverAddress = await flashLoanReceiver.getAddress();
      const params = ethers.AbiCoder.defaultAbiCoder().encode(
        ["address"],
        [receiverAddress]
      );

      await expect(
        adapter
          .connect(user1)
          .flashLoanSimple(await token.getAddress(), FLASH_AMOUNT, params)
      ).to.be.reverted;
    });
  });

  describe("Flash Loan (Multiple Assets)", function () {
    const FLASH_AMOUNT1 = ethers.parseEther("1000");
    const FLASH_AMOUNT2 = ethers.parseEther("500");
    const EXPECTED_PREMIUM1 = (FLASH_AMOUNT1 * 9n) / 10000n;
    const EXPECTED_PREMIUM2 = (FLASH_AMOUNT2 * 9n) / 10000n;

    beforeEach(async function () {
      // Fund flash loan receiver for repayment
      await token.transfer(
        await flashLoanReceiver.getAddress(),
        EXPECTED_PREMIUM1
      );
      await collateralToken.transfer(
        await flashLoanReceiver.getAddress(),
        EXPECTED_PREMIUM2
      );
    });

    it("Should execute multi-asset flash loan successfully", async function () {
      const receiverAddress = await flashLoanReceiver.getAddress();
      const params = ethers.AbiCoder.defaultAbiCoder().encode(
        ["address"],
        [receiverAddress]
      );

      const assets = [await token.getAddress(), await collateralToken.getAddress()];
      const amounts = [FLASH_AMOUNT1, FLASH_AMOUNT2];
      const modes = [0, 0]; // No debt conversion

      await expect(
        adapter.connect(user1).flashLoan(assets, amounts, modes, params)
      )
        .to.emit(flashLoanReceiver, "FlashLoanReceived");
    });

    it("Should emit FlashLoanExecuted events for all assets", async function () {
      const receiverAddress = await flashLoanReceiver.getAddress();
      const params = ethers.AbiCoder.defaultAbiCoder().encode(
        ["address"],
        [receiverAddress]
      );

      const assets = [await token.getAddress(), await collateralToken.getAddress()];
      const amounts = [FLASH_AMOUNT1, FLASH_AMOUNT2];
      const modes = [0, 0];

      const tx = await adapter
        .connect(user1)
        .flashLoan(assets, amounts, modes, params);
      const receipt = await tx.wait();

      const flashLoanEvents = receipt!.logs.filter(
        (log: any) => log.fragment && log.fragment.name === "FlashLoanExecuted"
      );
      expect(flashLoanEvents.length).to.equal(2);
    });

    it("Should revert with empty assets array", async function () {
      const receiverAddress = await flashLoanReceiver.getAddress();
      const params = ethers.AbiCoder.defaultAbiCoder().encode(
        ["address"],
        [receiverAddress]
      );

      await expect(
        adapter.connect(user1).flashLoan([], [], [], params)
      ).to.be.revertedWith("Must borrow at least one asset");
    });

    it("Should revert with mismatched array lengths", async function () {
      const receiverAddress = await flashLoanReceiver.getAddress();
      const params = ethers.AbiCoder.defaultAbiCoder().encode(
        ["address"],
        [receiverAddress]
      );

      const assets = [await token.getAddress()];
      const amounts = [FLASH_AMOUNT1, FLASH_AMOUNT2];
      const modes = [0];

      await expect(
        adapter.connect(user1).flashLoan(assets, amounts, modes, params)
      ).to.be.revertedWith("Array length mismatch");
    });
  });

  describe("Collateral Management", function () {
    beforeEach(async function () {
      await adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT);
    });

    it("Should allow enabling asset as collateral", async function () {
      await expect(
        adapter
          .connect(user1)
          .setUserUseReserveAsCollateral(await token.getAddress(), true)
      )
        .to.emit(adapter, "CollateralStatusChanged")
        .withArgs(await user1.getAddress(), await token.getAddress(), true);
    });

    it("Should allow disabling asset as collateral", async function () {
      await adapter
        .connect(user1)
        .setUserUseReserveAsCollateral(await token.getAddress(), true);

      await expect(
        adapter
          .connect(user1)
          .setUserUseReserveAsCollateral(await token.getAddress(), false)
      )
        .to.emit(adapter, "CollateralStatusChanged")
        .withArgs(await user1.getAddress(), await token.getAddress(), false);
    });

    it("Should update collateral status in pool", async function () {
      await adapter
        .connect(user1)
        .setUserUseReserveAsCollateral(await token.getAddress(), true);

      // Note: In the adapter pattern, the pool sees the adapter as the caller
      const isCollateral = await aavePool.userUseAsCollateral(
        await adapter.getAddress(),
        await token.getAddress()
      );
      expect(isCollateral).to.be.true;
    });
  });

  describe("Interest Rate Mode Swapping", function () {
    beforeEach(async function () {
      await adapter
        .connect(user1)
        .supply(await collateralToken.getAddress(), SUPPLY_AMOUNT);
      await adapter
        .connect(user1)
        .borrow(
          await token.getAddress(),
          BORROW_AMOUNT,
          INTEREST_RATE_MODE_VARIABLE
        );
    });

    it("Should allow swapping from variable to stable rate", async function () {
      await expect(
        adapter
          .connect(user1)
          .swapBorrowRateMode(await token.getAddress(), INTEREST_RATE_MODE_VARIABLE)
      )
        .to.emit(adapter, "InterestRateModeSwapped")
        .withArgs(
          await user1.getAddress(),
          await token.getAddress(),
          INTEREST_RATE_MODE_STABLE
        );
    });

    it("Should allow swapping from stable to variable rate", async function () {
      // First swap to stable
      await adapter
        .connect(user1)
        .swapBorrowRateMode(await token.getAddress(), INTEREST_RATE_MODE_VARIABLE);

      // Then swap back to variable
      await expect(
        adapter
          .connect(user1)
          .swapBorrowRateMode(await token.getAddress(), INTEREST_RATE_MODE_STABLE)
      )
        .to.emit(adapter, "InterestRateModeSwapped")
        .withArgs(
          await user1.getAddress(),
          await token.getAddress(),
          INTEREST_RATE_MODE_VARIABLE
        );
    });
  });

  describe("View Functions", function () {
    it("Should return user account data", async function () {
      const accountData = await adapter.getUserAccountData(await user1.getAddress());

      expect(accountData.totalCollateralBase).to.equal(1000e8);
      expect(accountData.totalDebtBase).to.equal(500e8);
      expect(accountData.healthFactor).to.equal(ethers.parseEther("2"));
    });

    it("Should return reserve data", async function () {
      const reserveData = await adapter.getReserveData(await token.getAddress());

      expect(reserveData.aTokenAddress).to.equal(await aToken.getAddress());
      expect(reserveData.liquidityIndex).to.equal(ethers.parseUnits("1", 27)); // 1e27 (ray format)
    });

    it("Should return user position", async function () {
      await adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT);

      const position = await adapter.getUserPosition(await user1.getAddress());

      expect(position.totalSupplied).to.equal(SUPPLY_AMOUNT);
      expect(position.totalBorrowed).to.equal(0);
    });

    it("Should calculate flash loan premium correctly", async function () {
      const amount = ethers.parseEther("10000");
      const expectedPremium = (amount * 9n) / 10000n; // 0.09%

      const premium = await adapter.calculateFlashLoanPremium(amount);

      expect(premium).to.equal(expectedPremium);
    });
  });

  describe("Admin Functions", function () {
    it("Should allow owner to update state aggregator", async function () {
      const newAggregator = await (
        await ethers.getContractFactory("L2StateAggregator")
      ).deploy(await owner.getAddress());

      await adapter.updateStateAggregator(await newAggregator.getAddress());

      expect(await adapter.stateAggregator()).to.equal(
        await newAggregator.getAddress()
      );
    });

    it("Should allow owner to recover tokens", async function () {
      const recoverAmount = ethers.parseEther("100");
      await token.transfer(await adapter.getAddress(), recoverAmount);

      const ownerAddress = await owner.getAddress();
      const balanceBefore = await token.balanceOf(ownerAddress);

      await adapter.recoverToken(await token.getAddress(), recoverAmount);

      expect(await token.balanceOf(ownerAddress)).to.equal(
        balanceBefore + recoverAmount
      );
    });

    it("Should revert if non-owner tries to update state aggregator", async function () {
      const newAggregator = await (
        await ethers.getContractFactory("L2StateAggregator")
      ).deploy(await owner.getAddress());

      await expect(
        adapter
          .connect(user1)
          .updateStateAggregator(await newAggregator.getAddress())
      ).to.be.reverted;
    });

    it("Should revert if non-owner tries to recover tokens", async function () {
      const recoverAmount = ethers.parseEther("100");
      await token.transfer(await adapter.getAddress(), recoverAmount);

      await expect(
        adapter
          .connect(user1)
          .recoverToken(await token.getAddress(), recoverAmount)
      ).to.be.reverted;
    });
  });

  describe("Complex Scenarios", function () {
    it("Should handle complete lending lifecycle", async function () {
      // 1. Supply collateral
      await adapter
        .connect(user1)
        .supply(await collateralToken.getAddress(), SUPPLY_AMOUNT);

      // 2. Borrow against collateral
      await adapter
        .connect(user1)
        .borrow(
          await token.getAddress(),
          BORROW_AMOUNT,
          INTEREST_RATE_MODE_VARIABLE
        );

      // 3. Repay half
      await adapter
        .connect(user1)
        .repay(
          await token.getAddress(),
          BORROW_AMOUNT / 2n,
          INTEREST_RATE_MODE_VARIABLE
        );

      // 4. Repay remaining (get actual remaining amount)
      const positionBeforeRepay = await adapter.userPositions(await user1.getAddress());
      await adapter
        .connect(user1)
        .repay(
          await token.getAddress(),
          positionBeforeRepay.totalBorrowed,
          INTEREST_RATE_MODE_VARIABLE
        );

      // 5. Withdraw collateral
      await adapter
        .connect(user1)
        .withdraw(await collateralToken.getAddress(), ethers.MaxUint256);

      const position = await adapter.userPositions(await user1.getAddress());
      expect(position.totalSupplied).to.equal(0);
      expect(position.totalBorrowed).to.equal(0);
      expect(await adapter.activeUsers()).to.equal(0);
    });

    it("Should maintain independent positions for multiple users", async function () {
      // User1 supplies
      await adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT);

      // User2 supplies different amount
      await adapter
        .connect(user2)
        .supply(await token.getAddress(), SUPPLY_AMOUNT * 2n);

      const pos1 = await adapter.userPositions(await user1.getAddress());
      const pos2 = await adapter.userPositions(await user2.getAddress());

      expect(pos1.totalSupplied).to.equal(SUPPLY_AMOUNT);
      expect(pos2.totalSupplied).to.equal(SUPPLY_AMOUNT * 2n);
      expect(await adapter.activeUsers()).to.equal(2);
    });

    it("Should track total amounts correctly across multiple operations", async function () {
      // User1 operations
      await adapter.connect(user1).supply(await token.getAddress(), SUPPLY_AMOUNT);
      await adapter
        .connect(user1)
        .supply(await collateralToken.getAddress(), SUPPLY_AMOUNT);
      await adapter
        .connect(user1)
        .borrow(
          await token.getAddress(),
          BORROW_AMOUNT,
          INTEREST_RATE_MODE_VARIABLE
        );

      // User2 operations
      await adapter
        .connect(user2)
        .supply(await token.getAddress(), SUPPLY_AMOUNT / 2n);
      await adapter
        .connect(user2)
        .supply(await collateralToken.getAddress(), SUPPLY_AMOUNT);
      await adapter
        .connect(user2)
        .borrow(
          await token.getAddress(),
          BORROW_AMOUNT / 2n,
          INTEREST_RATE_MODE_VARIABLE
        );

      expect(await adapter.totalSupplied()).to.equal(
        SUPPLY_AMOUNT + SUPPLY_AMOUNT / 2n + SUPPLY_AMOUNT * 2n
      );
      expect(await adapter.totalBorrowed()).to.equal(
        BORROW_AMOUNT + BORROW_AMOUNT / 2n
      );
      expect(await adapter.activeUsers()).to.equal(2);
    });
  });
});
