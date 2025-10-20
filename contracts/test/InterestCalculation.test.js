import { expect } from "chai";
import { ethers } from "hardhat";
import { loadFixture, time } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe("Interest Calculation Tests", function () {
  // Constants from contract
  const STABILITY_FEE = 200; // 2% annual (200 basis points)
  const SECONDS_PER_YEAR = 365 * 24 * 60 * 60;

  // Fixture for deploying system
  async function deploySystemFixture() {
    const [owner, user1, user2] = await ethers.getSigners();

    // Deploy LoyaltyUSD
    const LoyaltyUSD = await ethers.getContractFactory("LoyaltyUSD");
    const lusd = await LoyaltyUSD.deploy();

    // Deploy LP token
    const lpToken = await LoyaltyUSD.deploy();

    // Deploy CollateralVault
    const CollateralVault = await ethers.getContractFactory("CollateralVault");
    const vault = await CollateralVault.deploy(await lpToken.getAddress());

    // Setup roles
    await lusd.grantRole(await lusd.MINTER_ROLE(), owner.address);
    await lpToken.grantRole(await lpToken.MINTER_ROLE(), owner.address);

    // Mint LP tokens
    await lpToken.mint(user1.address, ethers.parseUnits("100000", 6));
    await lpToken.mint(user2.address, ethers.parseUnits("100000", 6));

    return { lusd, lpToken, vault, owner, user1, user2 };
  }

  describe("Basic Interest Calculation", function () {
    it("Should calculate interest correctly for 1 year", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("1000", 6);

      // Setup position
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Move time forward by 1 year
      await time.increase(SECONDS_PER_YEAR);

      // Calculate expected interest: 1000 * 0.02 = 20 LUSD
      const expectedInterest = ethers.parseUnits("20", 6);
      const accruedInterest = await vault.accruedInterest(user1.address);

      // Allow small rounding difference
      expect(accruedInterest).to.be.closeTo(expectedInterest, ethers.parseUnits("0.1", 6));
    });

    it("Should calculate interest correctly for 6 months", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("1000", 6);

      // Setup position
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Move time forward by 6 months
      await time.increase(SECONDS_PER_YEAR / 2);

      // Calculate expected interest: 1000 * 0.02 * 0.5 = 10 LUSD
      const expectedInterest = ethers.parseUnits("10", 6);
      const accruedInterest = await vault.accruedInterest(user1.address);

      expect(accruedInterest).to.be.closeTo(expectedInterest, ethers.parseUnits("0.1", 6));
    });

    it("Should calculate interest correctly for 30 days", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("1000", 6);

      // Setup position
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Move time forward by 30 days
      const thirtyDays = 30 * 24 * 60 * 60;
      await time.increase(thirtyDays);

      // Calculate expected interest: 1000 * 0.02 * (30/365) ≈ 1.64 LUSD
      const expectedInterest = (debt * BigInt(STABILITY_FEE) * BigInt(thirtyDays)) /
        BigInt(10000 * SECONDS_PER_YEAR);
      const accruedInterest = await vault.accruedInterest(user1.address);

      expect(accruedInterest).to.equal(expectedInterest);
    });

    it("Should return zero interest for zero debt", async function () {
      const { vault, lpToken, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("1000", 6);

      // Deposit collateral but no debt
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);

      // Move time forward
      await time.increase(SECONDS_PER_YEAR);

      const accruedInterest = await vault.accruedInterest(user1.address);
      expect(accruedInterest).to.equal(0);
    });
  });

  describe("Compound Interest Effects", function () {
    it("Should compound interest when debt is increased", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("3000", 6);
      const initialDebt = ethers.parseUnits("1000", 6);

      // Setup position
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, initialDebt);

      // Move time forward by 6 months
      await time.increase(SECONDS_PER_YEAR / 2);

      // Interest should be accrued (1000 * 0.02 * 0.5 = 10)
      const midInterest = await vault.accruedInterest(user1.address);
      expect(midInterest).to.be.closeTo(ethers.parseUnits("10", 6), ethers.parseUnits("0.1", 6));

      // Increase debt (this should compound the interest)
      const additionalDebt = ethers.parseUnits("500", 6);
      await vault.connect(owner).increaseDebt(user1.address, additionalDebt);

      // Check that interest was compounded into debt
      const debtAfterIncrease = await vault.debtAmount(user1.address);
      expect(debtAfterIncrease).to.be.closeTo(
        initialDebt + midInterest + additionalDebt,
        ethers.parseUnits("0.1", 6)
      );

      // Move time forward another 6 months
      await time.increase(SECONDS_PER_YEAR / 2);

      // New interest should be calculated on compounded amount
      const finalInterest = await vault.accruedInterest(user1.address);
      // (1510) * 0.02 * 0.5 ≈ 15.1 LUSD
      expect(finalInterest).to.be.closeTo(
        ethers.parseUnits("15.1", 6),
        ethers.parseUnits("0.5", 6)
      );
    });

    it("Should reset interest calculation after debt decrease", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("2000", 6);
      const debt = ethers.parseUnits("1000", 6);

      // Setup position
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Move time forward
      await time.increase(SECONDS_PER_YEAR / 4); // 3 months

      // Accrue some interest
      const interestBefore = await vault.accruedInterest(user1.address);
      expect(interestBefore).to.be.greaterThan(0);

      // Decrease debt (this should reset interest calculation)
      const repayAmount = ethers.parseUnits("500", 6);
      await vault.connect(owner).decreaseDebt(user1.address, repayAmount);

      // Interest should be zero immediately after debt change
      const interestAfter = await vault.accruedInterest(user1.address);
      expect(interestAfter).to.equal(0);

      // Move time forward again
      await time.increase(SECONDS_PER_YEAR / 4);

      // Interest should accrue on remaining debt
      const newInterest = await vault.accruedInterest(user1.address);
      // Remaining debt ≈ 505 (500 + compounded interest), interest ≈ 505 * 0.02 * 0.25 ≈ 2.52
      expect(newInterest).to.be.closeTo(
        ethers.parseUnits("2.52", 6),
        ethers.parseUnits("0.5", 6)
      );
    });
  });

  describe("Interest Impact on Position Health", function () {
    it("Should reduce collateral ratio as interest accrues", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("1000", 6);

      // Setup position at exactly 150% ratio
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Initial ratio should be 150%
      let ratio = await vault.getCollateralRatio(user1.address);
      expect(ratio).to.equal(150);

      // Move time forward by 1 year
      await time.increase(SECONDS_PER_YEAR);

      // Ratio should decrease due to interest (debt increased by ~2%)
      ratio = await vault.getCollateralRatio(user1.address);
      // New debt ≈ 1020, ratio ≈ 1500/1020 ≈ 147%
      expect(ratio).to.be.closeTo(147, 1);
    });

    it("Should make position liquidatable after enough interest accrual", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      // Setup position just above liquidation threshold (121%)
      const collateral = ethers.parseUnits("1210", 6);
      const debt = ethers.parseUnits("1000", 6);

      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Initially not liquidatable
      let canLiquidate = await vault.canLiquidate(user1.address);
      expect(canLiquidate).to.be.false;

      // Move time forward to accrue interest
      // Need debt to grow from 1000 to > 1008.33 (1210/120*100)
      // Interest needed: 8.33+, which is ~0.833% of 1000
      // Time needed: 0.833/2 * year ≈ 0.42 years ≈ 152 days
      await time.increase(152 * 24 * 60 * 60);

      // Should now be liquidatable
      canLiquidate = await vault.canLiquidate(user1.address);
      expect(canLiquidate).to.be.true;
    });
  });

  describe("getTotalDebt Function", function () {
    it("Should return debt plus accrued interest", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("2000", 6);
      const debt = ethers.parseUnits("1000", 6);

      // Setup position
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Move time forward
      await time.increase(SECONDS_PER_YEAR / 2); // 6 months

      // Get total debt
      const totalDebt = await vault.getTotalDebt(user1.address);
      const expectedTotal = debt + ethers.parseUnits("10", 6); // 1000 + (1000 * 0.02 * 0.5)

      expect(totalDebt).to.be.closeTo(expectedTotal, ethers.parseUnits("0.1", 6));
    });

    it("Should match sum of debt and accrued interest", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("2000", 6);
      const debt = ethers.parseUnits("1234", 6);

      // Setup position
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Move time forward
      const days90 = 90 * 24 * 60 * 60;
      await time.increase(days90);

      // Calculate components
      const baseDebt = await vault.debtAmount(user1.address);
      const interest = await vault.accruedInterest(user1.address);
      const totalDebt = await vault.getTotalDebt(user1.address);

      // Total should equal sum
      expect(totalDebt).to.equal(baseDebt + interest);
    });
  });

  describe("Multi-User Interest Tracking", function () {
    it("Should track interest independently for each user", async function () {
      const { vault, lpToken, owner, user1, user2 } = await loadFixture(deploySystemFixture);

      // Setup different positions at different times
      const collateral1 = ethers.parseUnits("2000", 6);
      const debt1 = ethers.parseUnits("1000", 6);

      const collateral2 = ethers.parseUnits("3000", 6);
      const debt2 = ethers.parseUnits("1500", 6);

      // User1 position
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral1);
      await vault.connect(user1).depositCollateral(collateral1);
      await vault.connect(owner).increaseDebt(user1.address, debt1);

      // Wait 30 days
      await time.increase(30 * 24 * 60 * 60);

      // User2 position (30 days after user1)
      await lpToken.connect(user2).approve(await vault.getAddress(), collateral2);
      await vault.connect(user2).depositCollateral(collateral2);
      await vault.connect(owner).increaseDebt(user2.address, debt2);

      // Wait another 30 days
      await time.increase(30 * 24 * 60 * 60);

      // User1 should have 60 days of interest
      const interest1 = await vault.accruedInterest(user1.address);
      const expected1 = (debt1 * BigInt(STABILITY_FEE) * BigInt(60 * 24 * 60 * 60)) /
        BigInt(10000 * SECONDS_PER_YEAR);

      // User2 should have 30 days of interest
      const interest2 = await vault.accruedInterest(user2.address);
      const expected2 = (debt2 * BigInt(STABILITY_FEE) * BigInt(30 * 24 * 60 * 60)) /
        BigInt(10000 * SECONDS_PER_YEAR);

      expect(interest1).to.equal(expected1);
      expect(interest2).to.equal(expected2);
      expect(interest1).to.be.greaterThan(interest2 * debt1 / debt2); // Adjusted for debt difference
    });
  });

  describe("Edge Cases", function () {
    it("Should handle zero time passed", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("1500", 6);
      const debt = ethers.parseUnits("1000", 6);

      // Setup position
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // No time passed
      const interest = await vault.accruedInterest(user1.address);
      expect(interest).to.equal(0);
    });

    it("Should handle very small debt amounts", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("100", 6);
      const debt = ethers.parseUnits("1", 6); // 1 LUSD

      // Setup position
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Move time forward by 1 year
      await time.increase(SECONDS_PER_YEAR);

      // Interest should be 0.02 LUSD
      const interest = await vault.accruedInterest(user1.address);
      const expected = debt * BigInt(STABILITY_FEE) / BigInt(10000);

      expect(interest).to.equal(expected);
    });

    it("Should handle interest calculation after full debt repayment", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("2000", 6);
      const debt = ethers.parseUnits("1000", 6);

      // Setup position
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Move time forward
      await time.increase(SECONDS_PER_YEAR / 4);

      // Full repayment
      const totalDebt = await vault.getTotalDebt(user1.address);
      await vault.connect(owner).decreaseDebt(user1.address, totalDebt);

      // No interest should accrue on zero debt
      await time.increase(SECONDS_PER_YEAR);
      const interest = await vault.accruedInterest(user1.address);
      expect(interest).to.equal(0);
    });
  });
});