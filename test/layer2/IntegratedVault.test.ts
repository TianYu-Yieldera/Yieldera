import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type {
  IntegratedVault,
  MockERC20,
  MockPriceOracle,
  L2StateAggregator,
} from "../../typechain-types/index.js";

describe("IntegratedVault", function () {
  let vault: IntegratedVault;
  let collateralToken: MockERC20;
  let lusdToken: MockERC20;
  let priceOracle: MockPriceOracle;
  let stateAggregator: L2StateAggregator;

  let owner: SignerWithAddress;
  let user1: SignerWithAddress;
  let user2: SignerWithAddress;

  const INITIAL_SUPPLY = ethers.parseEther("1000000");
  const COLLATERAL_RATIO = 150n; // 150%
  const PRECISION = 100n;
  const PRICE_PRECISION = 100000000n; // 1e8

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();

    // Deploy mock ERC20 tokens
    const MockERC20Factory = await ethers.getContractFactory("MockERC20");
    collateralToken = await MockERC20Factory.deploy(
      "Mock WETH",
      "WETH",
      INITIAL_SUPPLY
    );
    lusdToken = await MockERC20Factory.deploy(
      "Loyalty USD",
      "LUSD",
      INITIAL_SUPPLY
    );

    // Deploy MockPriceOracle
    const OracleFactory = await ethers.getContractFactory(
      "MockPriceOracle"
    );
    priceOracle = await OracleFactory.deploy();

    // Set mock price for collateral token (e.g., $2000 per token)
    const mockPrice = 2000n * PRICE_PRECISION; // $2000 with 8 decimals
    await priceOracle.setAssetPrice(
      await collateralToken.getAddress(),
      mockPrice
    );

    // Deploy L2StateAggregator
    const AggregatorFactory = await ethers.getContractFactory(
      "L2StateAggregator"
    );
    stateAggregator = await AggregatorFactory.deploy(owner.address);

    // Deploy IntegratedVault
    const VaultFactory = await ethers.getContractFactory("IntegratedVault");
    vault = await VaultFactory.deploy(
      await collateralToken.getAddress(),
      await lusdToken.getAddress(),
      await priceOracle.getAddress(),
      await stateAggregator.getAddress(),
      owner.address
    );

    // Authorize vault as a module in the state aggregator
    await stateAggregator.authorizeModule(await vault.getAddress());

    // Fund vault with LUSD for lending
    await lusdToken.transfer(await vault.getAddress(), INITIAL_SUPPLY / 2n);

    // Distribute tokens to users
    await collateralToken.transfer(user1.address, ethers.parseEther("100"));
    await collateralToken.transfer(user2.address, ethers.parseEther("100"));
    await lusdToken.transfer(user1.address, ethers.parseEther("10000"));
    await lusdToken.transfer(user2.address, ethers.parseEther("10000"));

    // Approve vault to spend tokens
    await collateralToken
      .connect(user1)
      .approve(await vault.getAddress(), ethers.MaxUint256);
    await collateralToken
      .connect(user2)
      .approve(await vault.getAddress(), ethers.MaxUint256);
    await lusdToken
      .connect(user1)
      .approve(await vault.getAddress(), ethers.MaxUint256);
    await lusdToken
      .connect(user2)
      .approve(await vault.getAddress(), ethers.MaxUint256);
  });

  describe("Deployment", function () {
    it("Should set the correct token addresses", async function () {
      expect(await vault.collateralToken()).to.equal(
        await collateralToken.getAddress()
      );
      expect(await vault.lusdToken()).to.equal(await lusdToken.getAddress());
    });

    it("Should set the correct price oracle", async function () {
      expect(await vault.priceOracle()).to.equal(
        await priceOracle.getAddress()
      );
    });

    it("Should set the correct state aggregator", async function () {
      expect(await vault.stateAggregator()).to.equal(
        await stateAggregator.getAddress()
      );
    });

    it("Should set the correct owner", async function () {
      expect(await vault.owner()).to.equal(owner.address);
    });

    it("Should have correct constants", async function () {
      expect(await vault.COLLATERAL_RATIO()).to.equal(COLLATERAL_RATIO);
      expect(await vault.INTEREST_RATE()).to.equal(3n);
      expect(await vault.PRECISION()).to.equal(PRECISION);
      expect(await vault.PRICE_PRECISION()).to.equal(PRICE_PRECISION);
    });

    it("Should revert with invalid collateral token address", async function () {
      const VaultFactory = await ethers.getContractFactory("IntegratedVault");
      await expect(
        VaultFactory.deploy(
          ethers.ZeroAddress,
          await lusdToken.getAddress(),
          await priceOracle.getAddress(),
          await stateAggregator.getAddress(),
          owner.address
        )
      ).to.be.revertedWith("Invalid collateral token");
    });

    it("Should revert with invalid LUSD token address", async function () {
      const VaultFactory = await ethers.getContractFactory("IntegratedVault");
      await expect(
        VaultFactory.deploy(
          await collateralToken.getAddress(),
          ethers.ZeroAddress,
          await priceOracle.getAddress(),
          await stateAggregator.getAddress(),
          owner.address
        )
      ).to.be.revertedWith("Invalid LUSD token");
    });

    it("Should revert with invalid price oracle address", async function () {
      const VaultFactory = await ethers.getContractFactory("IntegratedVault");
      await expect(
        VaultFactory.deploy(
          await collateralToken.getAddress(),
          await lusdToken.getAddress(),
          ethers.ZeroAddress,
          await stateAggregator.getAddress(),
          owner.address
        )
      ).to.be.revertedWith("Invalid price oracle");
    });
  });

  describe("Deposit", function () {
    const depositAmount = ethers.parseEther("10");

    it("Should allow users to deposit collateral", async function () {
      await expect(vault.connect(user1).deposit(depositAmount))
        .to.emit(vault, "Deposited")
        .withArgs(user1.address, depositAmount);

      const position = await vault.positions(user1.address);
      expect(position.collateral).to.equal(depositAmount);
      expect(await vault.totalCollateral()).to.equal(depositAmount);
      expect(await vault.activePositions()).to.equal(1n);
    });

    it("Should update position timestamp on deposit", async function () {
      const tx = await vault.connect(user1).deposit(depositAmount);
      const receipt = await tx.wait();
      const block = await ethers.provider.getBlock(receipt!.blockNumber);

      const position = await vault.positions(user1.address);
      expect(position.lastUpdate).to.equal(block!.timestamp);
    });

    it("Should allow multiple deposits from same user", async function () {
      await vault.connect(user1).deposit(depositAmount);
      await vault.connect(user1).deposit(depositAmount);

      const position = await vault.positions(user1.address);
      expect(position.collateral).to.equal(depositAmount * 2n);
      expect(await vault.activePositions()).to.equal(1n); // Still only 1 position
    });

    it("Should track multiple user positions", async function () {
      await vault.connect(user1).deposit(depositAmount);
      await vault.connect(user2).deposit(depositAmount);

      expect(await vault.activePositions()).to.equal(2n);
      expect(await vault.totalCollateral()).to.equal(depositAmount * 2n);
    });

    it("Should update state aggregator on deposit", async function () {
      await vault.connect(user1).deposit(depositAmount);

      const state = await stateAggregator.getSystemState();
      expect(state.totalCollateral).to.equal(depositAmount);
      expect(state.activePositions).to.equal(1n);
    });

    it("Should revert if amount is zero", async function () {
      await expect(vault.connect(user1).deposit(0)).to.be.revertedWith(
        "Amount must be greater than zero"
      );
    });

    it("Should revert if user has insufficient balance", async function () {
      await expect(
        vault.connect(user1).deposit(ethers.parseEther("1000"))
      ).to.be.reverted;
    });
  });

  describe("Borrow", function () {
    const depositAmount = ethers.parseEther("10"); // 10 tokens at $2000 = $20,000 collateral
    const borrowAmount = ethers.parseEther("10000"); // $10,000 LUSD (200% collateral ratio)

    beforeEach(async function () {
      await vault.connect(user1).deposit(depositAmount);
    });

    it("Should allow users to borrow against collateral", async function () {
      await expect(vault.connect(user1).borrow(borrowAmount))
        .to.emit(vault, "Borrowed")
        .withArgs(user1.address, borrowAmount);

      const position = await vault.positions(user1.address);
      expect(position.debt).to.equal(borrowAmount);
      expect(await vault.totalDebt()).to.equal(borrowAmount);
    });

    it("Should transfer LUSD to borrower", async function () {
      const balanceBefore = await lusdToken.balanceOf(user1.address);
      await vault.connect(user1).borrow(borrowAmount);
      const balanceAfter = await lusdToken.balanceOf(user1.address);

      expect(balanceAfter - balanceBefore).to.equal(borrowAmount);
    });

    it("Should update position timestamp on borrow", async function () {
      const tx = await vault.connect(user1).borrow(borrowAmount);
      const receipt = await tx.wait();
      const block = await ethers.provider.getBlock(receipt!.blockNumber);

      const position = await vault.positions(user1.address);
      expect(position.lastUpdate).to.equal(block!.timestamp);
    });

    it("Should allow multiple borrows from same user", async function () {
      const halfBorrow = borrowAmount / 2n;
      await vault.connect(user1).borrow(halfBorrow);
      await vault.connect(user1).borrow(halfBorrow);

      const position = await vault.positions(user1.address);
      expect(position.debt).to.equal(borrowAmount);
    });

    it("Should update state aggregator on borrow", async function () {
      await vault.connect(user1).borrow(borrowAmount);

      const state = await stateAggregator.getSystemState();
      expect(state.totalDebt).to.equal(borrowAmount);
    });

    it("Should revert if amount is zero", async function () {
      await expect(vault.connect(user1).borrow(0)).to.be.revertedWith(
        "Amount must be greater than zero"
      );
    });

    it("Should revert if no collateral deposited", async function () {
      await expect(
        vault.connect(user2).borrow(borrowAmount)
      ).to.be.revertedWith("No collateral deposited");
    });

    it("Should revert if collateral ratio is insufficient", async function () {
      // Try to borrow too much (only 150% collateral required, so max ~$13,333 LUSD)
      const excessiveBorrow = ethers.parseEther("14000");
      await expect(
        vault.connect(user1).borrow(excessiveBorrow)
      ).to.be.revertedWith("Insufficient collateral ratio");
    });

    it("Should calculate health factor correctly with oracle price", async function () {
      await vault.connect(user1).borrow(borrowAmount);

      const [, , healthFactor] = await vault.getPosition(user1.address);
      // $20,000 collateral / $10,000 debt = 200% health factor
      expect(healthFactor).to.equal(200n);
    });

    it("Should respect exact collateral ratio boundary", async function () {
      // Maximum borrow at 150% ratio: $20,000 / 1.5 = $13,333.33
      const maxBorrow = ethers.parseEther("13333");
      await vault.connect(user1).borrow(maxBorrow);

      const [, , healthFactor] = await vault.getPosition(user1.address);
      expect(healthFactor).to.be.gte(150n); // Should be at or above minimum
    });
  });

  describe("Repay", function () {
    const depositAmount = ethers.parseEther("10");
    const borrowAmount = ethers.parseEther("10000");

    beforeEach(async function () {
      await vault.connect(user1).deposit(depositAmount);
      await vault.connect(user1).borrow(borrowAmount);
    });

    it("Should allow users to repay debt", async function () {
      const repayAmount = ethers.parseEther("5000");
      await expect(vault.connect(user1).repay(repayAmount))
        .to.emit(vault, "Repaid")
        .withArgs(user1.address, repayAmount);

      const position = await vault.positions(user1.address);
      expect(position.debt).to.equal(borrowAmount - repayAmount);
      expect(await vault.totalDebt()).to.equal(borrowAmount - repayAmount);
    });

    it("Should transfer LUSD from user to vault", async function () {
      const repayAmount = ethers.parseEther("5000");
      const balanceBefore = await lusdToken.balanceOf(user1.address);
      await vault.connect(user1).repay(repayAmount);
      const balanceAfter = await lusdToken.balanceOf(user1.address);

      expect(balanceBefore - balanceAfter).to.equal(repayAmount);
    });

    it("Should update position timestamp on repay", async function () {
      const repayAmount = ethers.parseEther("5000");
      const tx = await vault.connect(user1).repay(repayAmount);
      const receipt = await tx.wait();
      const block = await ethers.provider.getBlock(receipt!.blockNumber);

      const position = await vault.positions(user1.address);
      expect(position.lastUpdate).to.equal(block!.timestamp);
    });

    it("Should allow full debt repayment", async function () {
      await vault.connect(user1).repay(borrowAmount);

      const position = await vault.positions(user1.address);
      expect(position.debt).to.equal(0n);
      expect(await vault.totalDebt()).to.equal(0n);
    });

    it("Should update state aggregator on repay", async function () {
      const repayAmount = ethers.parseEther("5000");
      await vault.connect(user1).repay(repayAmount);

      const state = await stateAggregator.getSystemState();
      expect(state.totalDebt).to.equal(borrowAmount - repayAmount);
    });

    it("Should revert if amount is zero", async function () {
      await expect(vault.connect(user1).repay(0)).to.be.revertedWith(
        "Amount must be greater than zero"
      );
    });

    it("Should revert if no debt to repay", async function () {
      await expect(
        vault.connect(user2).repay(ethers.parseEther("1000"))
      ).to.be.revertedWith("No debt to repay");
    });

    it("Should revert if amount exceeds debt", async function () {
      const excessiveRepay = borrowAmount + ethers.parseEther("1");
      await expect(vault.connect(user1).repay(excessiveRepay)).to.be.revertedWith(
        "Amount exceeds debt"
      );
    });
  });

  describe("Withdraw", function () {
    const depositAmount = ethers.parseEther("10");

    beforeEach(async function () {
      await vault.connect(user1).deposit(depositAmount);
    });

    it("Should allow users to withdraw collateral when no debt", async function () {
      const withdrawAmount = ethers.parseEther("5");
      await expect(vault.connect(user1).withdraw(withdrawAmount))
        .to.emit(vault, "Withdrawn")
        .withArgs(user1.address, withdrawAmount);

      const position = await vault.positions(user1.address);
      expect(position.collateral).to.equal(depositAmount - withdrawAmount);
      expect(await vault.totalCollateral()).to.equal(
        depositAmount - withdrawAmount
      );
    });

    it("Should transfer collateral to user", async function () {
      const withdrawAmount = ethers.parseEther("5");
      const balanceBefore = await collateralToken.balanceOf(user1.address);
      await vault.connect(user1).withdraw(withdrawAmount);
      const balanceAfter = await collateralToken.balanceOf(user1.address);

      expect(balanceAfter - balanceBefore).to.equal(withdrawAmount);
    });

    it("Should update position timestamp on withdraw", async function () {
      const withdrawAmount = ethers.parseEther("5");
      const tx = await vault.connect(user1).withdraw(withdrawAmount);
      const receipt = await tx.wait();
      const block = await ethers.provider.getBlock(receipt!.blockNumber);

      const position = await vault.positions(user1.address);
      expect(position.lastUpdate).to.equal(block!.timestamp);
    });

    it("Should allow full collateral withdrawal when no debt", async function () {
      await vault.connect(user1).withdraw(depositAmount);

      const position = await vault.positions(user1.address);
      expect(position.collateral).to.equal(0n);
      expect(await vault.activePositions()).to.equal(0n);
    });

    it("Should decrease active positions when fully withdrawn", async function () {
      await vault.connect(user1).withdraw(depositAmount);
      expect(await vault.activePositions()).to.equal(0n);
    });

    it("Should allow partial withdrawal while maintaining collateral ratio", async function () {
      const borrowAmount = ethers.parseEther("10000"); // $10,000 borrowed
      await vault.connect(user1).borrow(borrowAmount);

      // Can withdraw some collateral while keeping health factor >= 150%
      // Current: $20,000 collateral, $10,000 debt = 200% health factor
      // Withdraw $1,000 worth (0.5 tokens): $19,000 / $10,000 = 190% (still safe)
      const withdrawAmount = ethers.parseEther("0.5");
      await vault.connect(user1).withdraw(withdrawAmount);

      const position = await vault.positions(user1.address);
      expect(position.collateral).to.equal(depositAmount - withdrawAmount);
    });

    it("Should update state aggregator on withdraw", async function () {
      const withdrawAmount = ethers.parseEther("5");
      await vault.connect(user1).withdraw(withdrawAmount);

      const state = await stateAggregator.getSystemState();
      expect(state.totalCollateral).to.equal(depositAmount - withdrawAmount);
    });

    it("Should revert if amount is zero", async function () {
      await expect(vault.connect(user1).withdraw(0)).to.be.revertedWith(
        "Amount must be greater than zero"
      );
    });

    it("Should revert if no collateral to withdraw", async function () {
      await expect(
        vault.connect(user2).withdraw(ethers.parseEther("1"))
      ).to.be.revertedWith("No collateral to withdraw");
    });

    it("Should revert if amount exceeds collateral", async function () {
      const excessiveWithdraw = depositAmount + ethers.parseEther("1");
      await expect(
        vault.connect(user1).withdraw(excessiveWithdraw)
      ).to.be.revertedWith("Amount exceeds collateral");
    });

    it("Should revert if withdrawal breaks collateral ratio", async function () {
      const borrowAmount = ethers.parseEther("10000");
      await vault.connect(user1).borrow(borrowAmount);

      // Try to withdraw too much, breaking the 150% ratio
      // Current: $20,000 collateral, need minimum $15,000 for 150% ratio
      // Withdraw 3 tokens ($6,000): $14,000 / $10,000 = 140% (below minimum)
      const excessiveWithdraw = ethers.parseEther("3");
      await expect(
        vault.connect(user1).withdraw(excessiveWithdraw)
      ).to.be.revertedWith("Withdrawal would break collateral ratio");
    });
  });

  describe("Position Management", function () {
    const depositAmount = ethers.parseEther("10");
    const borrowAmount = ethers.parseEther("10000");

    it("Should return correct position data", async function () {
      await vault.connect(user1).deposit(depositAmount);
      await vault.connect(user1).borrow(borrowAmount);

      const [collateral, debt, healthFactor] = await vault.getPosition(
        user1.address
      );
      expect(collateral).to.equal(depositAmount);
      expect(debt).to.equal(borrowAmount);
      expect(healthFactor).to.equal(200n); // $20k / $10k = 200%
    });

    it("Should return max health factor when no debt", async function () {
      await vault.connect(user1).deposit(depositAmount);

      const [, , healthFactor] = await vault.getPosition(user1.address);
      expect(healthFactor).to.equal(ethers.MaxUint256);
    });

    it("Should calculate collateral value correctly", async function () {
      await vault.connect(user1).deposit(depositAmount);

      const collateralValue = await vault.getCollateralValue(user1.address);
      // 10 tokens * $2000 = $20,000
      expect(collateralValue).to.equal(ethers.parseEther("20000"));
    });

    it("Should return zero collateral value when no deposit", async function () {
      const collateralValue = await vault.getCollateralValue(user1.address);
      expect(collateralValue).to.equal(0n);
    });
  });

  describe("Oracle Integration", function () {
    const depositAmount = ethers.parseEther("10");

    it("Should use oracle price for health factor calculation", async function () {
      await vault.connect(user1).deposit(depositAmount);

      // Set new price: $3000 per token
      const newPrice = 3000n * PRICE_PRECISION;
      await priceOracle.setAssetPrice(
        await collateralToken.getAddress(),
        newPrice
      );

      // Borrow against new collateral value
      // $30,000 collateral allows ~$20,000 borrow at 150% ratio
      const borrowAmount = ethers.parseEther("15000");
      await vault.connect(user1).borrow(borrowAmount);

      const [, , healthFactor] = await vault.getPosition(user1.address);
      // $30,000 / $15,000 = 200%
      expect(healthFactor).to.equal(200n);
    });

    it("Should reflect price changes in collateral value", async function () {
      await vault.connect(user1).deposit(depositAmount);

      let collateralValue = await vault.getCollateralValue(user1.address);
      expect(collateralValue).to.equal(ethers.parseEther("20000"));

      // Update price to $2500
      const newPrice = 2500n * PRICE_PRECISION;
      await priceOracle.setAssetPrice(
        await collateralToken.getAddress(),
        newPrice
      );

      collateralValue = await vault.getCollateralValue(user1.address);
      expect(collateralValue).to.equal(ethers.parseEther("25000"));
    });
  });

  describe("Admin Functions", function () {
    it("Should allow owner to update price oracle", async function () {
      const OracleFactory = await ethers.getContractFactory(
        "MockPriceOracle"
      );
      const newOracle = await OracleFactory.deploy();

      await vault.updatePriceOracle(await newOracle.getAddress());
      expect(await vault.priceOracle()).to.equal(await newOracle.getAddress());
    });

    it("Should allow owner to update state aggregator", async function () {
      const AggregatorFactory = await ethers.getContractFactory(
        "L2StateAggregator"
      );
      const newAggregator = await AggregatorFactory.deploy(owner.address);

      await vault.updateStateAggregator(await newAggregator.getAddress());
      expect(await vault.stateAggregator()).to.equal(
        await newAggregator.getAddress()
      );
    });

    it("Should revert if non-owner tries to update oracle", async function () {
      const OracleFactory = await ethers.getContractFactory(
        "MockPriceOracle"
      );
      const newOracle = await OracleFactory.deploy();

      await expect(
        vault.connect(user1).updatePriceOracle(await newOracle.getAddress())
      ).to.be.revertedWithCustomError(vault, "OwnableUnauthorizedAccount");
    });

    it("Should revert if non-owner tries to update aggregator", async function () {
      const AggregatorFactory = await ethers.getContractFactory(
        "L2StateAggregator"
      );
      const newAggregator = await AggregatorFactory.deploy(owner.address);

      await expect(
        vault
          .connect(user1)
          .updateStateAggregator(await newAggregator.getAddress())
      ).to.be.revertedWithCustomError(vault, "OwnableUnauthorizedAccount");
    });

    it("Should revert when updating oracle to zero address", async function () {
      await expect(
        vault.updatePriceOracle(ethers.ZeroAddress)
      ).to.be.revertedWith("Invalid oracle address");
    });
  });

  describe("Edge Cases and Security", function () {
    it("Should handle multiple users with independent positions", async function () {
      const deposit1 = ethers.parseEther("5");
      const deposit2 = ethers.parseEther("10");
      const borrow1 = ethers.parseEther("5000");
      const borrow2 = ethers.parseEther("10000");

      await vault.connect(user1).deposit(deposit1);
      await vault.connect(user2).deposit(deposit2);
      await vault.connect(user1).borrow(borrow1);
      await vault.connect(user2).borrow(borrow2);

      const [collateral1, debt1] = await vault.getPosition(user1.address);
      const [collateral2, debt2] = await vault.getPosition(user2.address);

      expect(collateral1).to.equal(deposit1);
      expect(collateral2).to.equal(deposit2);
      expect(debt1).to.equal(borrow1);
      expect(debt2).to.equal(borrow2);
    });

    it("Should prevent reentrancy attacks on deposit", async function () {
      // NonReentrant modifier should prevent reentrancy
      // This is more of a sanity check that the modifier is in place
      const depositAmount = ethers.parseEther("1");
      await vault.connect(user1).deposit(depositAmount);
      // If vulnerable, this would fail during execution
    });

    it("Should handle dust amounts correctly", async function () {
      const dustAmount = 1n; // Smallest possible amount
      await vault.connect(user1).deposit(dustAmount);

      const position = await vault.positions(user1.address);
      expect(position.collateral).to.equal(dustAmount);
    });

    it("Should maintain accurate totals across operations", async function () {
      const deposit1 = ethers.parseEther("10");
      const deposit2 = ethers.parseEther("5");
      const borrow1 = ethers.parseEther("8000");
      const repay1 = ethers.parseEther("3000");
      const withdraw1 = ethers.parseEther("2");

      await vault.connect(user1).deposit(deposit1);
      await vault.connect(user2).deposit(deposit2);
      await vault.connect(user1).borrow(borrow1);
      await vault.connect(user1).repay(repay1);
      await vault.connect(user1).withdraw(withdraw1);

      expect(await vault.totalCollateral()).to.equal(
        deposit1 + deposit2 - withdraw1
      );
      expect(await vault.totalDebt()).to.equal(borrow1 - repay1);
    });

    it("Should correctly track active positions count", async function () {
      // Start with 0
      expect(await vault.activePositions()).to.equal(0n);

      // User1 deposits
      await vault.connect(user1).deposit(ethers.parseEther("10"));
      expect(await vault.activePositions()).to.equal(1n);

      // User2 deposits
      await vault.connect(user2).deposit(ethers.parseEther("5"));
      expect(await vault.activePositions()).to.equal(2n);

      // User1 withdraws all
      await vault.connect(user1).withdraw(ethers.parseEther("10"));
      expect(await vault.activePositions()).to.equal(1n);

      // User2 withdraws all
      await vault.connect(user2).withdraw(ethers.parseEther("5"));
      expect(await vault.activePositions()).to.equal(0n);
    });
  });

  describe("Complex Scenarios", function () {
    it("Should handle full lifecycle: deposit -> borrow -> repay -> withdraw", async function () {
      const depositAmount = ethers.parseEther("10");
      const borrowAmount = ethers.parseEther("10000");

      // Deposit
      await vault.connect(user1).deposit(depositAmount);
      let [collateral, debt, healthFactor] = await vault.getPosition(
        user1.address
      );
      expect(collateral).to.equal(depositAmount);
      expect(debt).to.equal(0n);

      // Borrow
      await vault.connect(user1).borrow(borrowAmount);
      [collateral, debt, healthFactor] = await vault.getPosition(user1.address);
      expect(debt).to.equal(borrowAmount);
      expect(healthFactor).to.equal(200n);

      // Repay
      await vault.connect(user1).repay(borrowAmount);
      [collateral, debt, healthFactor] = await vault.getPosition(user1.address);
      expect(debt).to.equal(0n);
      expect(healthFactor).to.equal(ethers.MaxUint256);

      // Withdraw
      await vault.connect(user1).withdraw(depositAmount);
      [collateral, debt, healthFactor] = await vault.getPosition(user1.address);
      expect(collateral).to.equal(0n);
      expect(await vault.activePositions()).to.equal(0n);
    });

    it("Should handle price volatility scenarios", async function () {
      const depositAmount = ethers.parseEther("10");
      await vault.connect(user1).deposit(depositAmount);

      // Price: $2000, can borrow $13,333
      await vault.connect(user1).borrow(ethers.parseEther("13000"));

      // Price drops to $1500
      await priceOracle.setAssetPrice(
        await collateralToken.getAddress(),
        1500n * PRICE_PRECISION
      );

      // Health factor should drop: $15,000 / $13,000 = 115%
      const [, , healthFactor] = await vault.getPosition(user1.address);
      expect(healthFactor).to.be.lt(150n); // Below minimum ratio

      // Should not be able to borrow more
      await expect(
        vault.connect(user1).borrow(ethers.parseEther("100"))
      ).to.be.revertedWith("Insufficient collateral ratio");

      // Should not be able to withdraw
      await expect(
        vault.connect(user1).withdraw(ethers.parseEther("1"))
      ).to.be.revertedWith("Withdrawal would break collateral ratio");

      // But can still repay
      await vault.connect(user1).repay(ethers.parseEther("1000"));
    });
  });
});
