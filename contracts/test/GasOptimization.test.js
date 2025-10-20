import { expect } from "chai";
import { ethers } from "hardhat";
import { loadFixture } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe("Gas Optimization & Consumption Tests", function () {
  // Gas limit thresholds for operations (in gas units)
  const GAS_LIMITS = {
    DEPOSIT: 150000,
    WITHDRAW: 100000,
    MINT: 100000,
    BURN: 80000,
    TRANSFER: 60000,
    DEBT_INCREASE: 120000,
    DEBT_DECREASE: 100000,
    LIQUIDATION: 200000
  };

  // Fixture for deploying full system
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
    const MINTER_ROLE = await lusd.MINTER_ROLE();
    const BURNER_ROLE = await lusd.BURNER_ROLE();

    await lusd.grantRole(MINTER_ROLE, await vault.getAddress());
    await lusd.grantRole(BURNER_ROLE, await vault.getAddress());
    await lusd.grantRole(MINTER_ROLE, owner.address);

    // Mint LP tokens to users
    await lpToken.grantRole(await lpToken.MINTER_ROLE(), owner.address);
    await lpToken.mint(user1.address, ethers.parseUnits("100000", 6));
    await lpToken.mint(user2.address, ethers.parseUnits("100000", 6));

    return { lusd, lpToken, vault, owner, user1, user2 };
  }

  describe("CollateralVault Gas Consumption", function () {
    it("Should consume reasonable gas for collateral deposit", async function () {
      const { vault, lpToken, user1 } = await loadFixture(deploySystemFixture);

      const depositAmount = ethers.parseUnits("1000", 6);
      await lpToken.connect(user1).approve(await vault.getAddress(), depositAmount);

      const tx = await vault.connect(user1).depositCollateral(depositAmount);
      const receipt = await tx.wait();

      console.log(`      Gas used for deposit: ${receipt.gasUsed.toString()}`);
      expect(receipt.gasUsed).to.be.lessThan(GAS_LIMITS.DEPOSIT);
    });

    it("Should consume reasonable gas for collateral withdrawal", async function () {
      const { vault, lpToken, user1 } = await loadFixture(deploySystemFixture);

      // Setup: deposit first
      const depositAmount = ethers.parseUnits("2000", 6);
      await lpToken.connect(user1).approve(await vault.getAddress(), depositAmount);
      await vault.connect(user1).depositCollateral(depositAmount);

      // Test withdrawal gas
      const withdrawAmount = ethers.parseUnits("1000", 6);
      const tx = await vault.connect(user1).withdrawCollateral(withdrawAmount);
      const receipt = await tx.wait();

      console.log(`      Gas used for withdrawal: ${receipt.gasUsed.toString()}`);
      expect(receipt.gasUsed).to.be.lessThan(GAS_LIMITS.WITHDRAW);
    });

    it("Should consume reasonable gas for debt operations", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      // Setup collateral
      const collateralAmount = ethers.parseUnits("3000", 6);
      await lpToken.connect(user1).approve(await vault.getAddress(), collateralAmount);
      await vault.connect(user1).depositCollateral(collateralAmount);

      // Test debt increase gas
      const debtAmount = ethers.parseUnits("1000", 6);
      const tx1 = await vault.connect(owner).increaseDebt(user1.address, debtAmount);
      const receipt1 = await tx1.wait();

      console.log(`      Gas used for debt increase: ${receipt1.gasUsed.toString()}`);
      expect(receipt1.gasUsed).to.be.lessThan(GAS_LIMITS.DEBT_INCREASE);

      // Test debt decrease gas
      const repayAmount = ethers.parseUnits("500", 6);
      const tx2 = await vault.connect(owner).decreaseDebt(user1.address, repayAmount);
      const receipt2 = await tx2.wait();

      console.log(`      Gas used for debt decrease: ${receipt2.gasUsed.toString()}`);
      expect(receipt2.gasUsed).to.be.lessThan(GAS_LIMITS.DEBT_DECREASE);
    });

    it("Should consume reasonable gas for liquidation", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      // Create liquidatable position
      const collateral = ethers.parseUnits("1150", 6);
      const debt = ethers.parseUnits("1000", 6);

      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Test liquidation gas
      const liquidateAmount = ethers.parseUnits("500", 6);
      const tx = await vault.connect(owner).liquidate(user1.address, liquidateAmount);
      const receipt = await tx.wait();

      console.log(`      Gas used for liquidation: ${receipt.gasUsed.toString()}`);
      expect(receipt.gasUsed).to.be.lessThan(GAS_LIMITS.LIQUIDATION);
    });
  });

  describe("LoyaltyUSD Gas Consumption", function () {
    it("Should consume reasonable gas for minting", async function () {
      const { lusd, owner, user1 } = await loadFixture(deploySystemFixture);

      const mintAmount = ethers.parseUnits("1000", 6);
      const tx = await lusd.connect(owner).mint(user1.address, mintAmount);
      const receipt = await tx.wait();

      console.log(`      Gas used for mint: ${receipt.gasUsed.toString()}`);
      expect(receipt.gasUsed).to.be.lessThan(GAS_LIMITS.MINT);
    });

    it("Should consume reasonable gas for burning", async function () {
      const { lusd, owner, user1 } = await loadFixture(deploySystemFixture);

      // Setup: mint first
      const mintAmount = ethers.parseUnits("1000", 6);
      await lusd.connect(owner).mint(user1.address, mintAmount);

      // Test burn gas
      const burnAmount = ethers.parseUnits("500", 6);
      const tx = await lusd.connect(user1).burn(user1.address, burnAmount);
      const receipt = await tx.wait();

      console.log(`      Gas used for burn: ${receipt.gasUsed.toString()}`);
      expect(receipt.gasUsed).to.be.lessThan(GAS_LIMITS.BURN);
    });

    it("Should consume reasonable gas for transfers", async function () {
      const { lusd, owner, user1, user2 } = await loadFixture(deploySystemFixture);

      // Setup: mint to user1
      const mintAmount = ethers.parseUnits("1000", 6);
      await lusd.connect(owner).mint(user1.address, mintAmount);

      // Test transfer gas
      const transferAmount = ethers.parseUnits("500", 6);
      const tx = await lusd.connect(user1).transfer(user2.address, transferAmount);
      const receipt = await tx.wait();

      console.log(`      Gas used for transfer: ${receipt.gasUsed.toString()}`);
      expect(receipt.gasUsed).to.be.lessThan(GAS_LIMITS.TRANSFER);
    });

    it("Should consume reasonable gas for approve + transferFrom", async function () {
      const { lusd, owner, user1, user2 } = await loadFixture(deploySystemFixture);

      // Setup: mint to user1
      const mintAmount = ethers.parseUnits("1000", 6);
      await lusd.connect(owner).mint(user1.address, mintAmount);

      // Test approve gas
      const approveAmount = ethers.parseUnits("500", 6);
      const tx1 = await lusd.connect(user1).approve(user2.address, approveAmount);
      const receipt1 = await tx1.wait();

      console.log(`      Gas used for approve: ${receipt1.gasUsed.toString()}`);

      // Test transferFrom gas
      const tx2 = await lusd.connect(user2).transferFrom(
        user1.address,
        user2.address,
        approveAmount
      );
      const receipt2 = await tx2.wait();

      console.log(`      Gas used for transferFrom: ${receipt2.gasUsed.toString()}`);
      expect(receipt2.gasUsed).to.be.lessThan(GAS_LIMITS.TRANSFER * 1.5);
    });
  });

  describe("Batch Operations Gas Optimization", function () {
    it("Should show gas savings with batched vs individual operations", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      // Prepare large collateral for multiple operations
      const totalCollateral = ethers.parseUnits("10000", 6);
      await lpToken.connect(user1).approve(await vault.getAddress(), totalCollateral);

      // Test individual small deposits (5 x 200)
      let totalGasIndividual = 0n;
      for (let i = 0; i < 5; i++) {
        const tx = await vault.connect(user1).depositCollateral(ethers.parseUnits("200", 6));
        const receipt = await tx.wait();
        totalGasIndividual += receipt.gasUsed;
      }

      // Reset state for comparison
      const withdrawTotal = ethers.parseUnits("1000", 6);
      await vault.connect(user1).withdrawCollateral(withdrawTotal);

      // Test single large deposit
      const tx = await vault.connect(user1).depositCollateral(ethers.parseUnits("1000", 6));
      const receipt = await tx.wait();
      const gasSingle = receipt.gasUsed;

      console.log(`      Total gas for 5 small deposits: ${totalGasIndividual.toString()}`);
      console.log(`      Gas for 1 large deposit: ${gasSingle.toString()}`);
      console.log(`      Gas savings: ${(totalGasIndividual - gasSingle).toString()}`);

      // Single transaction should be more efficient
      expect(gasSingle).to.be.lessThan(totalGasIndividual);
    });

    it("Should measure gas for complex multi-step operations", async function () {
      const { lusd, vault, lpToken, owner, user1, user2 } = await loadFixture(deploySystemFixture);

      // Complex operation: deposit, mint, transfer, burn, withdraw
      const collateralAmount = ethers.parseUnits("3000", 6);
      const mintAmount = ethers.parseUnits("1000", 6);

      // Step 1: Deposit collateral
      await lpToken.connect(user1).approve(await vault.getAddress(), collateralAmount);
      const tx1 = await vault.connect(user1).depositCollateral(collateralAmount);
      const gas1 = (await tx1.wait()).gasUsed;

      // Step 2: Increase debt
      const tx2 = await vault.connect(owner).increaseDebt(user1.address, mintAmount);
      const gas2 = (await tx2.wait()).gasUsed;

      // Step 3: Mint LUSD
      const tx3 = await lusd.connect(owner).mint(user1.address, mintAmount);
      const gas3 = (await tx3.wait()).gasUsed;

      // Step 4: Transfer LUSD
      const tx4 = await lusd.connect(user1).transfer(user2.address, ethers.parseUnits("500", 6));
      const gas4 = (await tx4.wait()).gasUsed;

      // Step 5: Burn LUSD
      const tx5 = await lusd.connect(user1).burn(user1.address, ethers.parseUnits("300", 6));
      const gas5 = (await tx5.wait()).gasUsed;

      const totalGas = gas1 + gas2 + gas3 + gas4 + gas5;
      console.log(`      Total gas for complete flow: ${totalGas.toString()}`);
      console.log(`      Average gas per operation: ${(totalGas / 5n).toString()}`);

      // Ensure total is reasonable (less than 500k gas for 5 operations)
      expect(totalGas).to.be.lessThan(500000);
    });
  });

  describe("Storage Optimization Tests", function () {
    it("Should efficiently handle storage updates", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      const depositAmount = ethers.parseUnits("1000", 6);
      await lpToken.connect(user1).approve(await vault.getAddress(), depositAmount * 10n);

      // First deposit (cold storage)
      const tx1 = await vault.connect(user1).depositCollateral(depositAmount);
      const gas1 = (await tx1.wait()).gasUsed;

      // Second deposit (warm storage)
      const tx2 = await vault.connect(user1).depositCollateral(depositAmount);
      const gas2 = (await tx2.wait()).gasUsed;

      console.log(`      Gas for first deposit (cold): ${gas1.toString()}`);
      console.log(`      Gas for second deposit (warm): ${gas2.toString()}`);
      console.log(`      Gas reduction: ${(gas1 - gas2).toString()}`);

      // Second transaction should use less gas due to warm storage
      expect(gas2).to.be.lessThan(gas1);
    });

    it("Should measure gas for reading state variables", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      // Setup position
      const collateral = ethers.parseUnits("2000", 6);
      const debt = ethers.parseUnits("1000", 6);

      await lpToken.connect(user1).approve(await vault.getAddress(), collateral);
      await vault.connect(user1).depositCollateral(collateral);
      await vault.connect(owner).increaseDebt(user1.address, debt);

      // Measure gas for view functions (these don't consume gas in production, but we can estimate)
      const position = await vault.getPosition(user1.address);
      const ratio = await vault.getCollateralRatio(user1.address);
      const maxMintable = await vault.getMaxMintable(user1.address);
      const canLiquidate = await vault.canLiquidate(user1.address);
      const stats = await vault.getVaultStats();

      // These are view functions, so gas is 0, but we verify they work
      expect(position.collateral).to.equal(collateral);
      expect(ratio).to.be.greaterThan(0);
      expect(maxMintable).to.be.greaterThanOrEqual(0);
      expect(canLiquidate).to.be.false;
      expect(stats._totalCollateral).to.equal(collateral);
    });
  });

  describe("Gas Optimization Patterns", function () {
    it("Should demonstrate SSTORE optimization", async function () {
      const { vault, lpToken, owner, user1 } = await loadFixture(deploySystemFixture);

      // Multiple small deposits vs one large deposit
      const smallAmount = ethers.parseUnits("100", 6);
      const largeAmount = ethers.parseUnits("1000", 6);

      await lpToken.connect(user1).approve(await vault.getAddress(), largeAmount * 2n);

      // Measure 10 small deposits
      let totalSmallGas = 0n;
      for (let i = 0; i < 10; i++) {
        const tx = await vault.connect(user1).depositCollateral(smallAmount);
        totalSmallGas += (await tx.wait()).gasUsed;
      }

      // Withdraw all
      await vault.connect(user1).withdrawCollateral(largeAmount);

      // Measure 1 large deposit
      const largeTx = await vault.connect(user1).depositCollateral(largeAmount);
      const largeGas = (await largeTx.wait()).gasUsed;

      const avgSmallGas = totalSmallGas / 10n;
      console.log(`      Average gas per small deposit: ${avgSmallGas.toString()}`);
      console.log(`      Gas for large deposit: ${largeGas.toString()}`);
      console.log(`      Efficiency ratio: ${((largeGas * 100n) / avgSmallGas).toString()}%`);

      // Large deposit should be more efficient per unit
      expect(largeGas).to.be.lessThan(avgSmallGas * 10n);
    });

    it("Should demonstrate event emission gas cost", async function () {
      const { vault, lpToken, user1 } = await loadFixture(deploySystemFixture);

      const depositAmount = ethers.parseUnits("1000", 6);
      await lpToken.connect(user1).approve(await vault.getAddress(), depositAmount);

      const tx = await vault.connect(user1).depositCollateral(depositAmount);
      const receipt = await tx.wait();

      // Check that events were emitted (gas cost included in total)
      const events = receipt.logs;
      console.log(`      Number of events emitted: ${events.length}`);
      console.log(`      Total gas with events: ${receipt.gasUsed.toString()}`);

      expect(events.length).to.be.greaterThan(0);
      expect(receipt.gasUsed).to.be.lessThan(GAS_LIMITS.DEPOSIT);
    });
  });
});