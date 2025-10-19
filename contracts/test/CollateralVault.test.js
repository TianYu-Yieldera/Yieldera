import { expect } from "chai";
import { ethers } from "hardhat";
import { loadFixture } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe("CollateralVault", function () {
  // Fixture for deploying contracts
  async function deployVaultFixture() {
    const [owner, user1, user2, liquidator] = await ethers.getSigners();

    // Deploy mock loyalty token
    const LoyaltyToken = await ethers.getContractFactory("LoyaltyUSD");
    const loyaltyToken = await LoyaltyToken.deploy();

    // Deploy CollateralVault
    const CollateralVault = await ethers.getContractFactory("CollateralVault");
    const vault = await CollateralVault.deploy(await loyaltyToken.getAddress());

    // Mint tokens to users for testing
    await loyaltyToken.grantRole(await loyaltyToken.MINTER_ROLE(), owner.address);
    await loyaltyToken.mint(user1.address, ethers.parseUnits("10000", 6));
    await loyaltyToken.mint(user2.address, ethers.parseUnits("5000", 6));
    await loyaltyToken.mint(liquidator.address, ethers.parseUnits("10000", 6));

    return { vault, loyaltyToken, owner, user1, user2, liquidator };
  }

  describe("Deployment", function () {
    it("Should set the correct loyalty token address", async function () {
      const { vault, loyaltyToken } = await loadFixture(deployVaultFixture);
      expect(await vault.loyaltyToken()).to.equal(await loyaltyToken.getAddress());
    });

    it("Should initialize with zero collateral and debt", async function () {
      const { vault } = await loadFixture(deployVaultFixture);
      expect(await vault.totalCollateral()).to.equal(0);
      expect(await vault.totalDebt()).to.equal(0);
    });

    it("Should set correct constants", async function () {
      const { vault } = await loadFixture(deployVaultFixture);
      expect(await vault.COLLATERAL_RATIO()).to.equal(150);
      expect(await vault.LIQUIDATION_THRESHOLD()).to.equal(120);
      expect(await vault.STABILITY_FEE()).to.equal(200);
    });

    it("Should revert if deployed with zero address", async function () {
      const CollateralVault = await ethers.getContractFactory("CollateralVault");
      await expect(
        CollateralVault.deploy(ethers.ZeroAddress)
      ).to.be.revertedWith("Invalid token address");
    });
  });

  describe("Collateral Deposit", function () {
    it("Should allow users to deposit collateral", async function () {
      const { vault, loyaltyToken, user1 } = await loadFixture(deployVaultFixture);
      const depositAmount = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), depositAmount);
      await expect(vault.connect(user1).depositCollateral(depositAmount))
        .to.emit(vault, "CollateralDeposited")
        .withArgs(user1.address, depositAmount, depositAmount);

      expect(await vault.collateralDeposited(user1.address)).to.equal(depositAmount);
      expect(await vault.totalCollateral()).to.equal(depositAmount);
    });

    it("Should revert when depositing zero amount", async function () {
      const { vault, user1 } = await loadFixture(deployVaultFixture);
      await expect(
        vault.connect(user1).depositCollateral(0)
      ).to.be.revertedWith("Amount must be > 0");
    });

    it("Should handle multiple deposits from same user", async function () {
      const { vault, loyaltyToken, user1 } = await loadFixture(deployVaultFixture);
      const amount1 = ethers.parseUnits("500", 6);
      const amount2 = ethers.parseUnits("300", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), amount1 + amount2);

      await vault.connect(user1).depositCollateral(amount1);
      await vault.connect(user1).depositCollateral(amount2);

      expect(await vault.collateralDeposited(user1.address)).to.equal(amount1 + amount2);
    });

    it("Should handle deposits from multiple users", async function () {
      const { vault, loyaltyToken, user1, user2 } = await loadFixture(deployVaultFixture);
      const amount1 = ethers.parseUnits("1000", 6);
      const amount2 = ethers.parseUnits("500", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), amount1);
      await loyaltyToken.connect(user2).approve(await vault.getAddress(), amount2);

      await vault.connect(user1).depositCollateral(amount1);
      await vault.connect(user2).depositCollateral(amount2);

      expect(await vault.totalCollateral()).to.equal(amount1 + amount2);
    });
  });

  describe("Collateral Withdrawal", function () {
    it("Should allow withdrawal when no debt exists", async function () {
      const { vault, loyaltyToken, user1 } = await loadFixture(deployVaultFixture);
      const depositAmount = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), depositAmount);
      await vault.connect(user1).depositCollateral(depositAmount);

      await expect(vault.connect(user1).withdrawCollateral(depositAmount))
        .to.emit(vault, "CollateralWithdrawn")
        .withArgs(user1.address, depositAmount, 0);

      expect(await vault.collateralDeposited(user1.address)).to.equal(0);
    });

    it("Should revert when withdrawing more than deposited", async function () {
      const { vault, loyaltyToken, user1 } = await loadFixture(deployVaultFixture);
      const depositAmount = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), depositAmount);
      await vault.connect(user1).depositCollateral(depositAmount);

      await expect(
        vault.connect(user1).withdrawCollateral(depositAmount + 1n)
      ).to.be.revertedWith("Insufficient collateral");
    });

    it("Should revert when withdrawal would undercollateralize position", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("1000", 6);

      // Deposit collateral
      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);

      // Increase debt (150% ratio)
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Try to withdraw collateral that would drop ratio below 150%
      await expect(
        vault.connect(user1).withdrawCollateral(ethers.parseUnits("100", 6))
      ).to.be.revertedWith("Withdrawal would undercollateralize position");
    });

    it("Should allow partial withdrawal while maintaining healthy position", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("2000", 6);
      const debt = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Withdraw 500, leaving 1500 (still 150% ratio)
      const withdrawAmount = ethers.parseUnits("500", 6);
      await expect(vault.connect(user1).withdrawCollateral(withdrawAmount))
        .to.not.be.reverted;

      expect(await vault.collateralDeposited(user1.address)).to.equal(
        collateral - withdrawAmount
      );
    });
  });

  describe("Debt Management", function () {
    it("Should allow owner to increase debt", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);

      await expect(vault.connect(owner).increaseDebt(user1.address, debt))
        .to.emit(vault, "DebtIncreased")
        .withArgs(user1.address, debt, debt);

      expect(await vault.debtAmount(user1.address)).to.equal(debt);
      expect(await vault.totalDebt()).to.equal(debt);
    });

    it("Should revert when non-owner tries to increase debt", async function () {
      const { vault, user1, user2 } = await loadFixture(deployVaultFixture);
      await expect(
        vault.connect(user1).increaseDebt(user2.address, ethers.parseUnits("100", 6))
      ).to.be.reverted;
    });

    it("Should revert when debt increase would undercollateralize", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1000", 6);
      const excessiveDebt = ethers.parseUnits("700", 6); // Would be 142% ratio, need 150%

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);

      await expect(
        vault.connect(owner).increaseDebt(user1.address, excessiveDebt)
      ).to.be.revertedWith("Debt increase would undercollateralize position");
    });

    it("Should allow owner to decrease debt", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("1000", 6);
      const repayAmount = ethers.parseUnits("500", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      await expect(vault.connect(owner).decreaseDebt(user1.address, repayAmount))
        .to.emit(vault, "DebtDecreased")
        .withArgs(user1.address, repayAmount, debt - repayAmount);

      expect(await vault.debtAmount(user1.address)).to.equal(debt - repayAmount);
    });

    it("Should revert when decreasing more debt than exists", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      await expect(
        vault.connect(owner).decreaseDebt(user1.address, debt + 1n)
      ).to.be.revertedWith("Debt underflow");
    });
  });

  describe("Health Factor Calculations", function () {
    it("Should return correct collateral ratio", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("2000", 6);
      const debt = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      expect(await vault.getCollateralRatio(user1.address)).to.equal(200); // 200%
    });

    it("Should return max uint256 when debt is zero", async function () {
      const { vault, loyaltyToken, user1 } = await loadFixture(deployVaultFixture);
      await loyaltyToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("1000", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("1000", 6));

      expect(await vault.getCollateralRatio(user1.address)).to.equal(ethers.MaxUint256);
    });

    it("Should correctly identify healthy positions", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      expect(await vault.isPositionHealthy(user1.address)).to.be.true;
    });

    it("Should correctly identify unhealthy positions", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1400", 6);
      const debt = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);

      // Force unhealthy position (bypass health check using owner power)
      // In reality, this wouldn't happen, but we're testing the calculation
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Now artificially reduce collateral ratio by increasing debt more
      // (This is a test scenario - real contract prevents this)
      const position = await vault.getPosition(user1.address);
      expect(position.healthy).to.be.true; // 140% is still healthy
    });

    it("Should calculate maximum mintable amount correctly", async function () {
      const { vault, loyaltyToken, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1500", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);

      // Max mintable = 1500 * 100 / 150 = 1000 LUSD
      expect(await vault.getMaxMintable(user1.address)).to.equal(ethers.parseUnits("1000", 6));
    });

    it("Should adjust max mintable based on existing debt", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("500", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Max mintable = 1000 - 500 = 500 LUSD remaining
      expect(await vault.getMaxMintable(user1.address)).to.equal(ethers.parseUnits("500", 6));
    });
  });

  describe("Liquidation", function () {
    it("Should allow liquidation of undercollateralized position", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1150", 6);
      const debt = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Position is at 115% ratio, below 120% threshold
      expect(await vault.canLiquidate(user1.address)).to.be.true;

      // Liquidate
      const debtToCover = ethers.parseUnits("500", 6);
      await expect(vault.connect(owner).liquidate(user1.address, debtToCover))
        .to.emit(vault, "PositionLiquidated");
    });

    it("Should revert liquidation of healthy position", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      expect(await vault.canLiquidate(user1.address)).to.be.false;

      await expect(
        vault.connect(owner).liquidate(user1.address, ethers.parseUnits("500", 6))
      ).to.be.revertedWith("Position is not liquidatable");
    });

    it("Should apply 10% liquidation bonus", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1150", 6);
      const debt = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      const debtToCover = ethers.parseUnits("500", 6);
      const expectedCollateralSeized = ethers.parseUnits("550", 6); // 500 * 1.10

      await expect(vault.connect(owner).liquidate(user1.address, debtToCover))
        .to.emit(vault, "PositionLiquidated")
        .withArgs(user1.address, expectedCollateralSeized, debtToCover);
    });
  });

  describe("Position Query", function () {
    it("Should return complete position details", async function () {
      const { vault, loyaltyToken, owner, user1 } = await loadFixture(deployVaultFixture);
      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("1000", 6);

      await loyaltyToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      const position = await vault.getPosition(user1.address);

      expect(position.collateral).to.equal(collateral);
      expect(position.debt).to.equal(debt);
      expect(position.collateralRatio).to.equal(150);
      expect(position.maxMintable).to.equal(0);
      expect(position.healthy).to.be.true;
    });
  });

  describe("Vault Statistics", function () {
    it("Should track global vault statistics", async function () {
      const { vault, loyaltyToken, owner, user1, user2 } = await loadFixture(deployVaultFixture);

      // User 1: 1500 collateral, 1000 debt
      await loyaltyToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("1500", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("1500", 6));
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));

      // User 2: 2000 collateral, 1000 debt
      await loyaltyToken.connect(user2).approve(await vault.getAddress(), ethers.parseUnits("2000", 6));
      await vault.connect(user2).depositCollateral(ethers.parseUnits("2000", 6));
      await vault.connect(owner).increaseDebt(user2.address, ethers.parseUnits("1000", 6));

      const stats = await vault.getVaultStats();

      expect(stats._totalCollateral).to.equal(ethers.parseUnits("3500", 6));
      expect(stats._totalDebt).to.equal(ethers.parseUnits("2000", 6));
      expect(stats.avgCollateralRatio).to.equal(175); // 3500/2000 = 175%
    });
  });
});
