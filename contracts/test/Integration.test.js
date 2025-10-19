import { expect } from "chai";
import { ethers } from "hardhat";
import { loadFixture } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe("LoyaltyUSD + CollateralVault Integration", function () {
  // Fixture for deploying full system
  async function deploySystemFixture() {
    const [owner, user1, user2, liquidator] = await ethers.getSigners();

    // Deploy LoyaltyUSD (stablecoin)
    const LoyaltyUSD = await ethers.getContractFactory("LoyaltyUSD");
    const lusd = await LoyaltyUSD.deploy();

    // Deploy a mock LP token to use as collateral
    const LPToken = await ethers.getContractFactory("LoyaltyUSD"); // Reuse for mock
    const lpToken = await LPToken.deploy();

    // Deploy CollateralVault
    const CollateralVault = await ethers.getContractFactory("CollateralVault");
    const vault = await CollateralVault.deploy(await lpToken.getAddress());

    // Setup roles
    const MINTER_ROLE = await lusd.MINTER_ROLE();
    const BURNER_ROLE = await lusd.BURNER_ROLE();

    await lusd.grantRole(MINTER_ROLE, await vault.getAddress());
    await lusd.grantRole(BURNER_ROLE, await vault.getAddress());

    // Mint LP tokens to users for collateral
    await lpToken.grantRole(await lpToken.MINTER_ROLE(), owner.address);
    await lpToken.mint(user1.address, ethers.parseUnits("10000", 6));
    await lpToken.mint(user2.address, ethers.parseUnits("5000", 6));
    await lpToken.mint(liquidator.address, ethers.parseUnits("10000", 6));

    return { lusd, lpToken, vault, owner, user1, user2, liquidator };
  }

  describe("System Setup", function () {
    it("Should deploy all contracts successfully", async function () {
      const { lusd, lpToken, vault } = await loadFixture(deploySystemFixture);

      expect(await lusd.getAddress()).to.not.equal(ethers.ZeroAddress);
      expect(await lpToken.getAddress()).to.not.equal(ethers.ZeroAddress);
      expect(await vault.getAddress()).to.not.equal(ethers.ZeroAddress);
    });

    it("Should grant vault MINTER_ROLE on LUSD", async function () {
      const { lusd, vault } = await loadFixture(deploySystemFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();

      expect(await lusd.hasRole(MINTER_ROLE, await vault.getAddress())).to.be.true;
    });

    it("Should grant vault BURNER_ROLE on LUSD", async function () {
      const { lusd, vault } = await loadFixture(deploySystemFixture);
      const BURNER_ROLE = await lusd.BURNER_ROLE();

      expect(await lusd.hasRole(BURNER_ROLE, await vault.getAddress())).to.be.true;
    });
  });

  describe("Mint & Redeem Flow", function () {
    it("Should complete full mint workflow", async function () {
      const { lusd, lpToken, vault, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateralAmount = ethers.parseUnits("1500", 6);
      const mintAmount = ethers.parseUnits("1000", 6);

      // Step 1: User deposits collateral
      await lpToken.connect(user1).approve(await vault.getAddress(), collateralAmount);
      await vault.connect(user1).depositCollateral(collateralAmount);

      // Step 2: Vault increases debt (simulating LUSD mint)
      await vault.connect(owner).increaseDebt(user1.address, mintAmount);

      // Step 3: Vault mints LUSD to user (in real impl, this happens automatically)
      await lusd.connect(vault.runner).mint(user1.address, mintAmount);

      // Verify state
      expect(await lusd.balanceOf(user1.address)).to.equal(mintAmount);
      expect(await vault.debtAmount(user1.address)).to.equal(mintAmount);
      expect(await vault.collateralDeposited(user1.address)).to.equal(collateralAmount);
      expect(await vault.getCollateralRatio(user1.address)).to.equal(150);
    });

    it("Should complete full redeem workflow", async function () {
      const { lusd, lpToken, vault, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateralAmount = ethers.parseUnits("1500", 6);
      const mintAmount = ethers.parseUnits("1000", 6);
      const redeemAmount = ethers.parseUnits("500", 6);

      // Mint flow
      await lpToken.connect(user1).approve(await vault.getAddress(), collateralAmount);
      await vault.connect(user1).depositCollateral(collateralAmount);
      await vault.connect(owner).increaseDebt(user1.address, mintAmount);
      await lusd.connect(vault.runner).mint(user1.address, mintAmount);

      // Redeem flow
      // Step 1: Burn LUSD
      await lusd.connect(user1).burn(user1.address, redeemAmount);

      // Step 2: Decrease debt in vault
      await vault.connect(owner).decreaseDebt(user1.address, redeemAmount);

      // Verify state
      expect(await lusd.balanceOf(user1.address)).to.equal(mintAmount - redeemAmount);
      expect(await vault.debtAmount(user1.address)).to.equal(mintAmount - redeemAmount);
      expect(await vault.getCollateralRatio(user1.address)).to.equal(300); // Improved ratio
    });

    it("Should allow withdrawing collateral after full redemption", async function () {
      const { lusd, lpToken, vault, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateralAmount = ethers.parseUnits("1500", 6);
      const mintAmount = ethers.parseUnits("1000", 6);

      // Mint
      await lpToken.connect(user1).approve(await vault.getAddress(), collateralAmount);
      await vault.connect(user1).depositCollateral(collateralAmount);
      await vault.connect(owner).increaseDebt(user1.address, mintAmount);
      await lusd.connect(vault.runner).mint(user1.address, mintAmount);

      // Redeem all
      await lusd.connect(user1).burn(user1.address, mintAmount);
      await vault.connect(owner).decreaseDebt(user1.address, mintAmount);

      // Withdraw all collateral
      await vault.connect(user1).withdrawCollateral(collateralAmount);

      expect(await lpToken.balanceOf(user1.address)).to.equal(ethers.parseUnits("10000", 6));
      expect(await vault.collateralDeposited(user1.address)).to.equal(0);
      expect(await vault.debtAmount(user1.address)).to.equal(0);
    });
  });

  describe("Collateral Management", function () {
    it("Should maintain proper collateral ratio after minting", async function () {
      const { lusd, lpToken, vault, owner, user1 } = await loadFixture(deploySystemFixture);

      // Deposit 3000 collateral
      await lpToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("3000", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("3000", 6));

      // Mint 1000 LUSD (300% ratio)
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));
      await lusd.connect(vault.runner).mint(user1.address, ethers.parseUnits("1000", 6));

      // Mint another 500 LUSD (200% ratio)
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("500", 6));
      await lusd.connect(vault.runner).mint(user1.address, ethers.parseUnits("500", 6));

      expect(await vault.getCollateralRatio(user1.address)).to.equal(200);
      expect(await lusd.balanceOf(user1.address)).to.equal(ethers.parseUnits("1500", 6));
    });

    it("Should allow adding collateral to improve health", async function () {
      const { lusd, lpToken, vault, owner, user1 } = await loadFixture(deploySystemFixture);

      // Initial position: 1500 collateral, 1000 debt (150% ratio)
      await lpToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("3000", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("1500", 6));
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));

      // Add more collateral
      await vault.connect(user1).depositCollateral(ethers.parseUnits("500", 6));

      expect(await vault.getCollateralRatio(user1.address)).to.equal(200); // Improved to 200%
    });

    it("Should prevent withdrawal that would undercollateralize", async function () {
      const { lusd, lpToken, vault, owner, user1 } = await loadFixture(deploySystemFixture);

      await lpToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("1500", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("1500", 6));
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));

      // Try to withdraw collateral that would drop ratio below 150%
      await expect(
        vault.connect(user1).withdrawCollateral(ethers.parseUnits("100", 6))
      ).to.be.revertedWith("Withdrawal would undercollateralize position");
    });
  });

  describe("Multi-User Scenarios", function () {
    it("Should handle multiple users with independent positions", async function () {
      const { lusd, lpToken, vault, owner, user1, user2 } = await loadFixture(deploySystemFixture);

      // User1: 1500 collateral, 1000 debt
      await lpToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("1500", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("1500", 6));
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));
      await lusd.connect(vault.runner).mint(user1.address, ethers.parseUnits("1000", 6));

      // User2: 2000 collateral, 1000 debt
      await lpToken.connect(user2).approve(await vault.getAddress(), ethers.parseUnits("2000", 6));
      await vault.connect(user2).depositCollateral(ethers.parseUnits("2000", 6));
      await vault.connect(owner).increaseDebt(user2.address, ethers.parseUnits("1000", 6));
      await lusd.connect(vault.runner).mint(user2.address, ethers.parseUnits("1000", 6));

      // Verify independent positions
      expect(await vault.getCollateralRatio(user1.address)).to.equal(150);
      expect(await vault.getCollateralRatio(user2.address)).to.equal(200);
      expect(await lusd.totalSupply()).to.equal(ethers.parseUnits("2000", 6));
    });

    it("Should handle transfers between users", async function () {
      const { lusd, lpToken, vault, owner, user1, user2 } = await loadFixture(deploySystemFixture);

      // User1 mints LUSD
      await lpToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("1500", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("1500", 6));
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));
      await lusd.connect(vault.runner).mint(user1.address, ethers.parseUnits("1000", 6));

      // User1 transfers LUSD to User2
      await lusd.connect(user1).transfer(user2.address, ethers.parseUnits("300", 6));

      // Verify balances (debt stays with user1)
      expect(await lusd.balanceOf(user1.address)).to.equal(ethers.parseUnits("700", 6));
      expect(await lusd.balanceOf(user2.address)).to.equal(ethers.parseUnits("300", 6));
      expect(await vault.debtAmount(user1.address)).to.equal(ethers.parseUnits("1000", 6));
      expect(await vault.debtAmount(user2.address)).to.equal(0);
    });
  });

  describe("Liquidation Scenarios", function () {
    it("Should liquidate undercollateralized position", async function () {
      const { lusd, lpToken, vault, owner, user1, liquidator } = await loadFixture(deploySystemFixture);

      // Create position at liquidation threshold (115% ratio)
      const collateral = ethers.parseUnits("1150", 6);
      const debt = ethers.parseUnits("1000", 6);

      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);
      await lusd.connect(vault.runner).mint(user1.address, debt);

      // Position should be liquidatable
      expect(await vault.canLiquidate(user1.address)).to.be.true;

      // Liquidate half the debt
      const liquidateAmount = ethers.parseUnits("500", 6);
      const expectedSeized = ethers.parseUnits("550", 6); // 500 * 1.10

      await vault.connect(owner).liquidate(user1.address, liquidateAmount);

      // Verify liquidation results
      expect(await vault.debtAmount(user1.address)).to.equal(debt - liquidateAmount);
      expect(await vault.collateralDeposited(user1.address)).to.equal(collateral - expectedSeized);
    });

    it("Should prevent liquidation of healthy positions", async function () {
      const { lusd, lpToken, vault, owner, user1 } = await loadFixture(deploySystemFixture);

      await lpToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("1500", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("1500", 6));
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));

      expect(await vault.canLiquidate(user1.address)).to.be.false;

      await expect(
        vault.connect(owner).liquidate(user1.address, ethers.parseUnits("500", 6))
      ).to.be.revertedWith("Position is not liquidatable");
    });
  });

  describe("Pause Functionality", function () {
    it("Should prevent minting when LUSD is paused", async function () {
      const { lusd, lpToken, vault, owner, user1 } = await loadFixture(deploySystemFixture);

      await lpToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("1500", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("1500", 6));
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));

      // Pause LUSD
      await lusd.pause();

      // Try to mint
      await expect(
        lusd.connect(vault.runner).mint(user1.address, ethers.parseUnits("1000", 6))
      ).to.be.reverted;
    });

    it("Should allow operations after unpausing", async function () {
      const { lusd, lpToken, vault, owner, user1 } = await loadFixture(deploySystemFixture);

      await lpToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("1500", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("1500", 6));
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));

      // Pause and unpause
      await lusd.pause();
      await lusd.unpause();

      // Should work now
      await expect(
        lusd.connect(vault.runner).mint(user1.address, ethers.parseUnits("1000", 6))
      ).to.not.be.reverted;
    });
  });

  describe("Edge Cases & Stress Tests", function () {
    it("Should handle maximum mintable amount correctly", async function () {
      const { lusd, lpToken, vault, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("1500", 6);
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);

      // Max mintable = 1500 * 100 / 150 = 1000
      const maxMintable = await vault.getMaxMintable(user1.address);
      expect(maxMintable).to.equal(ethers.parseUnits("1000", 6));

      // Mint exactly max amount
      await vault.connect(owner).increaseDebt(user1.address, maxMintable);
      await lusd.connect(vault.runner).mint(user1.address, maxMintable);

      // Should be at exactly 150% ratio
      expect(await vault.getCollateralRatio(user1.address)).to.equal(150);
      expect(await vault.getMaxMintable(user1.address)).to.equal(0);
    });

    it("Should handle rapid mint and redeem cycles", async function () {
      const { lusd, lpToken, vault, owner, user1 } = await loadFixture(deploySystemFixture);

      const collateral = ethers.parseUnits("3000", 6);
      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);

      // Cycle 1: Mint 1000
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));
      await lusd.connect(vault.runner).mint(user1.address, ethers.parseUnits("1000", 6));

      // Cycle 2: Redeem 500
      await lusd.connect(user1).burn(user1.address, ethers.parseUnits("500", 6));
      await vault.connect(owner).decreaseDebt(user1.address, ethers.parseUnits("500", 6));

      // Cycle 3: Mint 800
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("800", 6));
      await lusd.connect(vault.runner).mint(user1.address, ethers.parseUnits("800", 6));

      // Final state
      expect(await lusd.balanceOf(user1.address)).to.equal(ethers.parseUnits("1300", 6));
      expect(await vault.debtAmount(user1.address)).to.equal(ethers.parseUnits("1300", 6));
    });

    it("Should maintain system totals across multiple operations", async function () {
      const { lusd, lpToken, vault, owner, user1, user2 } = await loadFixture(deploySystemFixture);

      // User1 operations
      await lpToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("1500", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("1500", 6));
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));
      await lusd.connect(vault.runner).mint(user1.address, ethers.parseUnits("1000", 6));

      // User2 operations
      await lpToken.connect(user2).approve(await vault.getAddress(), ethers.parseUnits("3000", 6));
      await vault.connect(user2).depositCollateral(ethers.parseUnits("3000", 6));
      await vault.connect(owner).increaseDebt(user2.address, ethers.parseUnits("2000", 6));
      await lusd.connect(vault.runner).mint(user2.address, ethers.parseUnits("2000", 6));

      // Verify totals
      expect(await vault.totalCollateral()).to.equal(ethers.parseUnits("4500", 6));
      expect(await vault.totalDebt()).to.equal(ethers.parseUnits("3000", 6));
      expect(await lusd.totalSupply()).to.equal(ethers.parseUnits("3000", 6));

      // User1 redeems
      await lusd.connect(user1).burn(user1.address, ethers.parseUnits("500", 6));
      await vault.connect(owner).decreaseDebt(user1.address, ethers.parseUnits("500", 6));

      // Verify totals updated
      expect(await vault.totalDebt()).to.equal(ethers.parseUnits("2500", 6));
      expect(await lusd.totalSupply()).to.equal(ethers.parseUnits("2500", 6));
    });
  });

  describe("System Statistics", function () {
    it("Should track vault statistics correctly", async function () {
      const { lusd, lpToken, vault, owner, user1, user2 } = await loadFixture(deploySystemFixture);

      // Setup two positions
      await lpToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("1500", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("1500", 6));
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));

      await lpToken.connect(user2).approve(await vault.getAddress(), ethers.parseUnits("2000", 6));
      await vault.connect(user2).depositCollateral(ethers.parseUnits("2000", 6));
      await vault.connect(owner).increaseDebt(user2.address, ethers.parseUnits("1000", 6));

      const stats = await vault.getVaultStats();

      expect(stats._totalCollateral).to.equal(ethers.parseUnits("3500", 6));
      expect(stats._totalDebt).to.equal(ethers.parseUnits("2000", 6));
      expect(stats.avgCollateralRatio).to.equal(175); // 3500/2000 = 175%
    });

    it("Should track LUSD total supply correctly", async function () {
      const { lusd, lpToken, vault, owner, user1 } = await loadFixture(deploySystemFixture);

      await lpToken.connect(user1).approve(await vault.getAddress(), ethers.parseUnits("3000", 6));
      await vault.connect(user1).depositCollateral(ethers.parseUnits("3000", 6));

      // Mint 1000
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("1000", 6));
      await lusd.connect(vault.runner).mint(user1.address, ethers.parseUnits("1000", 6));
      expect(await lusd.totalSupply()).to.equal(ethers.parseUnits("1000", 6));

      // Mint another 500
      await vault.connect(owner).increaseDebt(user1.address, ethers.parseUnits("500", 6));
      await lusd.connect(vault.runner).mint(user1.address, ethers.parseUnits("500", 6));
      expect(await lusd.totalSupply()).to.equal(ethers.parseUnits("1500", 6));

      // Burn 300
      await lusd.connect(user1).burn(user1.address, ethers.parseUnits("300", 6));
      expect(await lusd.totalSupply()).to.equal(ethers.parseUnits("1200", 6));
    });
  });
});
