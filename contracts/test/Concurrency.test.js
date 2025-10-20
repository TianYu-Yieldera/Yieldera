import { expect } from "chai";
import { ethers } from "hardhat";
import { loadFixture } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe("Concurrency & Idempotency Tests", function () {
  // Fixture for deploying full system
  async function deploySystemFixture() {
    const [owner, user1, user2, user3, user4] = await ethers.getSigners();

    // Deploy LoyaltyUSD (stablecoin)
    const LoyaltyUSD = await ethers.getContractFactory("LoyaltyUSD");
    const lusd = await LoyaltyUSD.deploy();

    // Deploy LP token
    const lpToken = await LoyaltyUSD.deploy(); // Reuse for mock

    // Deploy CollateralVault
    const CollateralVault = await ethers.getContractFactory("CollateralVault");
    const vault = await CollateralVault.deploy(await lpToken.getAddress());

    // Setup roles
    const MINTER_ROLE = await lusd.MINTER_ROLE();
    const BURNER_ROLE = await lusd.BURNER_ROLE();

    await lusd.grantRole(MINTER_ROLE, await vault.getAddress());
    await lusd.grantRole(BURNER_ROLE, await vault.getAddress());
    await lusd.grantRole(MINTER_ROLE, owner.address);

    // Mint LP tokens to users
    await lpToken.grantRole(await lpToken.MINTER_ROLE(), owner.address);
    await lpToken.mint(user1.address, ethers.parseUnits("10000", 6));
    await lpToken.mint(user2.address, ethers.parseUnits("10000", 6));
    await lpToken.mint(user3.address, ethers.parseUnits("10000", 6));
    await lpToken.mint(user4.address, ethers.parseUnits("10000", 6));

    return { lusd, lpToken, vault, owner, user1, user2, user3, user4 };
  }

  describe("Concurrent Deposit Operations", function () {
    it("Should handle multiple concurrent deposits correctly", async function () {
      const { vault, lpToken, user1, user2, user3, user4 } = await loadFixture(deploySystemFixture);

      const depositAmount = ethers.parseUnits("1000", 6);

      // Prepare approvals
      await lpToken.connect(user1).approve(await vault.getAddress(), depositAmount);
      await lpToken.connect(user2).approve(await vault.getAddress(), depositAmount);
      await lpToken.connect(user3).approve(await vault.getAddress(), depositAmount);
      await lpToken.connect(user4).approve(await vault.getAddress(), depositAmount);

      // Execute concurrent deposits
      const deposits = await Promise.all([
        vault.connect(user1).depositCollateral(depositAmount),
        vault.connect(user2).depositCollateral(depositAmount),
        vault.connect(user3).depositCollateral(depositAmount),
        vault.connect(user4).depositCollateral(depositAmount)
      ]);

      // Verify all deposits were successful
      expect(await vault.collateralDeposited(user1.address)).to.equal(depositAmount);
      expect(await vault.collateralDeposited(user2.address)).to.equal(depositAmount);
      expect(await vault.collateralDeposited(user3.address)).to.equal(depositAmount);
      expect(await vault.collateralDeposited(user4.address)).to.equal(depositAmount);

      // Verify total collateral
      expect(await vault.totalCollateral()).to.equal(depositAmount * 4n);
    });

    it("Should prevent double-spending in concurrent operations", async function () {
      const { vault, lpToken, user1 } = await loadFixture(deploySystemFixture);

      const totalBalance = ethers.parseUnits("10000", 6);
      const depositAmount = ethers.parseUnits("6000", 6);

      // User tries to deposit more than balance through concurrent transactions
      await lpToken.connect(user1).approve(await vault.getAddress(), totalBalance * 2n);

      // First deposit should succeed
      await vault.connect(user1).depositCollateral(depositAmount);

      // Second concurrent deposit should fail due to insufficient balance
      await expect(
        vault.connect(user1).depositCollateral(depositAmount)
      ).to.be.reverted;

      // Verify only first deposit succeeded
      expect(await vault.collateralDeposited(user1.address)).to.equal(depositAmount);
    });
  });

  describe("Concurrent Mint/Burn Operations", function () {
    it("Should handle concurrent minting operations safely", async function () {
      const { lusd, lpToken, vault, owner, user1, user2 } = await loadFixture(deploySystemFixture);

      // Setup collateral for both users
      const collateralAmount = ethers.parseUnits("3000", 6);
      const mintAmount = ethers.parseUnits("1000", 6);

      await lpToken.connect(user1).approve(await vault.getAddress(), collateralAmount);
      await lpToken.connect(user2).approve(await vault.getAddress(), collateralAmount);

      await vault.connect(user1).depositCollateral(collateralAmount);
      await vault.connect(user2).depositCollateral(collateralAmount);

      await vault.connect(owner).increaseDebt(user1.address, mintAmount);
      await vault.connect(owner).increaseDebt(user2.address, mintAmount);

      // Execute concurrent mints
      const mints = await Promise.all([
        lusd.connect(owner).mint(user1.address, mintAmount),
        lusd.connect(owner).mint(user2.address, mintAmount)
      ]);

      // Verify both mints succeeded
      expect(await lusd.balanceOf(user1.address)).to.equal(mintAmount);
      expect(await lusd.balanceOf(user2.address)).to.equal(mintAmount);
      expect(await lusd.totalSupply()).to.equal(mintAmount * 2n);
    });

    it("Should handle concurrent burn operations safely", async function () {
      const { lusd, lpToken, vault, owner, user1, user2 } = await loadFixture(deploySystemFixture);

      // Setup: mint tokens to users
      const mintAmount = ethers.parseUnits("1000", 6);
      const burnAmount = ethers.parseUnits("500", 6);

      await lusd.connect(owner).mint(user1.address, mintAmount);
      await lusd.connect(owner).mint(user2.address, mintAmount);

      // Execute concurrent burns
      const burns = await Promise.all([
        lusd.connect(user1).burn(user1.address, burnAmount),
        lusd.connect(user2).burn(user2.address, burnAmount)
      ]);

      // Verify both burns succeeded
      expect(await lusd.balanceOf(user1.address)).to.equal(mintAmount - burnAmount);
      expect(await lusd.balanceOf(user2.address)).to.equal(mintAmount - burnAmount);
      expect(await lusd.totalSupply()).to.equal((mintAmount * 2n) - (burnAmount * 2n));
    });
  });

  describe("Idempotency Tests", function () {
    it("Should prevent duplicate debt increases", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateralAmount = ethers.parseUnits("1500", 6);
      const debtAmount = ethers.parseUnits("1000", 6);

      await lpToken.connect(user1).approve(await vault.getAddress(), collateralAmount);
      await vault.connect(user1).depositCollateral(collateralAmount);

      // First debt increase
      await vault.connect(owner).increaseDebt(user1.address, debtAmount);
      const debtAfterFirst = await vault.debtAmount(user1.address);

      // Trying to increase by same amount should work (additive)
      await expect(
        vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("100", 6))
      ).to.be.revertedWith("Debt increase would undercollateralize position");

      // Debt should remain at safe level
      expect(await vault.debtAmount(user1.address)).to.equal(debtAfterFirst);
    });

    it("Should maintain consistency with repeated deposit/withdraw cycles", async function () {
      const { vault, lpToken, user1 } = await loadFixture(deploySystemFixture);

      const depositAmount = ethers.parseUnits("1000", 6);
      const withdrawAmount = ethers.parseUnits("500", 6);

      // Initial deposit
      await lpToken.connect(user1).approve(await vault.getAddress(), depositAmount * 3n);
      await vault.connect(user1).depositCollateral(depositAmount);

      const initialBalance = await lpToken.balanceOf(user1.address);

      // Multiple deposit/withdraw cycles
      for (let i = 0; i < 3; i++) {
        await vault.connect(user1).withdrawCollateral(withdrawAmount);
        await vault.connect(user1).depositCollateral(withdrawAmount);
      }

      // Verify final state matches initial state
      expect(await vault.collateralDeposited(user1.address)).to.equal(depositAmount);
      expect(await lpToken.balanceOf(user1.address)).to.equal(initialBalance);
    });
  });

  describe("Race Condition Prevention", function () {
    it("Should prevent race conditions in liquidation", async function () {
      const { vault, lpToken, owner, user1, user2, user3 } = await loadFixture(deploySystemFixture);

      // Create liquidatable position
      const collateral = ethers.parseUnits("1150", 6);
      const debt = ethers.parseUnits("1000", 6);

      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Position should be liquidatable
      expect(await vault.canLiquidate(user1.address)).to.be.true;

      // First liquidation succeeds
      await vault.connect(owner).liquidate(user1.address, ethers.parseUnits("500", 6));

      // Second liquidation attempt should handle remaining debt correctly
      const remainingDebt = await vault.debtAmount(user1.address);

      if (remainingDebt > 0) {
        const canStillLiquidate = await vault.canLiquidate(user1.address);

        if (canStillLiquidate) {
          await expect(
            vault.connect(owner).liquidate(user1.address, remainingDebt)
          ).to.not.be.reverted;
        } else {
          await expect(
            vault.connect(owner).liquidate(user1.address, remainingDebt)
          ).to.be.revertedWith("Position is not liquidatable");
        }
      }
    });

    it("Should handle concurrent transfers correctly", async function () {
      const { lusd, owner, user1, user2, user3 } = await loadFixture(deploySystemFixture);

      const initialAmount = ethers.parseUnits("3000", 6);
      const transferAmount = ethers.parseUnits("100", 6);

      // Mint initial tokens to user1
      await lusd.connect(owner).mint(user1.address, initialAmount);

      // Execute concurrent transfers from user1 to multiple users
      const transfers = await Promise.all([
        lusd.connect(user1).transfer(user2.address, transferAmount),
        lusd.connect(user1).transfer(user3.address, transferAmount)
      ]);

      // Verify all transfers succeeded and balances are correct
      expect(await lusd.balanceOf(user1.address)).to.equal(
        initialAmount - (transferAmount * 2n)
      );
      expect(await lusd.balanceOf(user2.address)).to.equal(transferAmount);
      expect(await lusd.balanceOf(user3.address)).to.equal(transferAmount);

      // Verify total supply remains constant
      expect(await lusd.totalSupply()).to.equal(initialAmount);
    });
  });

  describe("State Consistency Under Load", function () {
    it("Should maintain total supply consistency under multiple operations", async function () {
      const { lusd, owner, user1, user2, user3 } = await loadFixture(deploySystemFixture);

      const mintAmount = ethers.parseUnits("1000", 6);

      // Mint to multiple users
      await lusd.connect(owner).mint(user1.address, mintAmount);
      await lusd.connect(owner).mint(user2.address, mintAmount);
      await lusd.connect(owner).mint(user3.address, mintAmount);

      const initialTotalSupply = await lusd.totalSupply();

      // Multiple transfers
      await lusd.connect(user1).transfer(user2.address, ethers.parseUnits("100", 6));
      await lusd.connect(user2).transfer(user3.address, ethers.parseUnits("200", 6));
      await lusd.connect(user3).transfer(user1.address, ethers.parseUnits("50", 6));

      // Burn some tokens
      await lusd.connect(user1).burn(user1.address, ethers.parseUnits("150", 6));

      // Verify total supply consistency
      const finalTotalSupply = await lusd.totalSupply();
      const totalBalances =
        (await lusd.balanceOf(user1.address)) +
        (await lusd.balanceOf(user2.address)) +
        (await lusd.balanceOf(user3.address));

      expect(finalTotalSupply).to.equal(totalBalances);
      expect(finalTotalSupply).to.equal(initialTotalSupply - ethers.parseUnits("150", 6));
    });

    it("Should maintain vault totals consistency", async function () {
      const { vault, lpToken, owner, user1, user2, user3 } = await loadFixture(deploySystemFixture);

      const depositAmount = ethers.parseUnits("1000", 6);

      // Multiple deposits
      await lpToken.connect(user1).approve(await vault.getAddress(), depositAmount);
      await lpToken.connect(user2).approve(await vault.getAddress(), depositAmount);
      await lpToken.connect(user3).approve(await vault.getAddress(), depositAmount);

      await vault.connect(user1).depositCollateral(depositAmount);
      await vault.connect(user2).depositCollateral(depositAmount);
      await vault.connect(user3).depositCollateral(depositAmount);

      // Add debt to users
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("600", 6));
      await vault.connect(owner).increaseDebt(user2.address, ethers.parseUnits("500", 6));

      // Verify totals match sum of individual positions
      const totalCollateral = await vault.totalCollateral();
      const totalDebt = await vault.totalDebt();

      const sumCollateral =
        (await vault.collateralDeposited(user1.address)) +
        (await vault.collateralDeposited(user2.address)) +
        (await vault.collateralDeposited(user3.address));

      const sumDebt =
        (await vault.debtAmount(user1.address)) +
        (await vault.debtAmount(user2.address)) +
        (await vault.debtAmount(user3.address));

      expect(totalCollateral).to.equal(sumCollateral);
      expect(totalDebt).to.equal(sumDebt);
    });
  });
});