import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type {
  TreasuryYieldDistributor,
  TreasuryAssetFactory,
  TreasuryToken,
} from "../../typechain-types/index.js";

describe("TreasuryYieldDistributor", function () {
  let distributor: TreasuryYieldDistributor;
  let factory: TreasuryAssetFactory;
  let token: TreasuryToken;

  let admin: SignerWithAddress;
  let holder1: SignerWithAddress;
  let holder2: SignerWithAddress;
  let yieldSource: SignerWithAddress;

  const OPERATOR_ROLE = ethers.keccak256(ethers.toUtf8Bytes("OPERATOR_ROLE"));

  let assetId: number;
  const FACE_VALUE = ethers.parseEther("1000");
  const HOLDER1_BALANCE = ethers.parseEther("100");
  const HOLDER2_BALANCE = ethers.parseEther("200");
  const YIELD_AMOUNT = ethers.parseEther("15"); // $15 total yield

  beforeEach(async function () {
    [admin, holder1, holder2, yieldSource] = await ethers.getSigners();

    // Deploy TreasuryAssetFactory
    const FactoryContract = await ethers.getContractFactory("TreasuryAssetFactory");
    factory = await FactoryContract.deploy(admin.address);

    // Deploy TreasuryYieldDistributor
    const DistributorFactory = await ethers.getContractFactory("TreasuryYieldDistributor");
    distributor = await DistributorFactory.deploy(
      admin.address,
      await factory.getAddress()
    );

    // Grant roles
    await factory.grantRole(OPERATOR_ROLE, admin.address);
    await distributor.grantRole(OPERATOR_ROLE, admin.address);

    // Create Treasury asset
    const issueDate = Math.floor(Date.now() / 1000);
    const maturityDate = issueDate + 90 * 24 * 60 * 60;

    await factory.createTreasuryAsset(
      0,
      "13W",
      "912796YZ1",
      issueDate,
      maturityDate,
      FACE_VALUE,
      525
    );

    assetId = 1;

    // Mint tokens to holders
    await factory.mintTokens(assetId, holder1.address, HOLDER1_BALANCE);
    await factory.mintTokens(assetId, holder2.address, HOLDER2_BALANCE);

    // Get token contract
    const asset = await factory.getTreasuryAsset(assetId);
    const TreasuryToken = await ethers.getContractFactory("TreasuryToken");
    token = TreasuryToken.attach(asset.tokenAddress) as any;
  });

  describe("Yield Deposit", function () {
    it("Should deposit yield for asset", async function () {
      const tx = await distributor.connect(admin).depositYield(
        assetId,
        { value: YIELD_AMOUNT }
      );

      await expect(tx)
        .to.emit(distributor, "YieldDeposited")
        .withArgs(assetId, YIELD_AMOUNT);

      const totalYield = await distributor.getTotalYield(assetId);
      expect(totalYield).to.equal(YIELD_AMOUNT);
    });

    it("Should fail with zero yield", async function () {
      await expect(
        distributor.connect(admin).depositYield(assetId, { value: 0 })
      ).to.be.revertedWith("Zero yield amount");
    });

    it("Should accumulate multiple deposits", async function () {
      await distributor.connect(admin).depositYield(assetId, { value: YIELD_AMOUNT });
      await distributor.connect(admin).depositYield(assetId, { value: YIELD_AMOUNT });

      const totalYield = await distributor.getTotalYield(assetId);
      expect(totalYield).to.equal(YIELD_AMOUNT * 2n);
    });
  });

  describe("Yield Calculation", function () {
    beforeEach(async function () {
      await distributor.connect(admin).depositYield(assetId, { value: YIELD_AMOUNT });
    });

    it("Should calculate pending yield correctly", async function () {
      const totalSupply = HOLDER1_BALANCE + HOLDER2_BALANCE;
      const expectedYield1 = (YIELD_AMOUNT * HOLDER1_BALANCE) / totalSupply;
      const expectedYield2 = (YIELD_AMOUNT * HOLDER2_BALANCE) / totalSupply;

      const pending1 = await distributor.getPendingYield(assetId, holder1.address);
      const pending2 = await distributor.getPendingYield(assetId, holder2.address);

      expect(pending1).to.be.closeTo(expectedYield1, ethers.parseEther("0.01"));
      expect(pending2).to.be.closeTo(expectedYield2, ethers.parseEther("0.01"));
    });

    it("Should return zero for non-holders", async function () {
      const pending = await distributor.getPendingYield(assetId, admin.address);
      expect(pending).to.equal(0);
    });
  });

  describe("Yield Claiming", function () {
    beforeEach(async function () {
      await distributor.connect(admin).depositYield(assetId, { value: YIELD_AMOUNT });
    });

    it("Should allow holders to claim yield", async function () {
      const pendingBefore = await distributor.getPendingYield(assetId, holder1.address);
      const balanceBefore = await ethers.provider.getBalance(holder1.address);

      const tx = await distributor.connect(holder1).claimYield(assetId);
      const receipt = await tx.wait();
      const gasUsed = receipt!.gasUsed * receipt!.gasPrice;

      const balanceAfter = await ethers.provider.getBalance(holder1.address);
      const pendingAfter = await distributor.getPendingYield(assetId, holder1.address);

      expect(pendingAfter).to.equal(0);
      expect(balanceAfter).to.be.closeTo(
        balanceBefore + pendingBefore - gasUsed,
        ethers.parseEther("0.001")
      );
    });

    it("Should emit YieldClaimed event", async function () {
      const pending = await distributor.getPendingYield(assetId, holder1.address);

      await expect(distributor.connect(holder1).claimYield(assetId))
        .to.emit(distributor, "YieldClaimed")
        .withArgs(assetId, holder1.address, pending);
    });

    it("Should fail if no yield to claim", async function () {
      await distributor.connect(holder1).claimYield(assetId);

      await expect(
        distributor.connect(holder1).claimYield(assetId)
      ).to.be.revertedWith("No yield to claim");
    });

    it("Should handle multiple claims correctly", async function () {
      await distributor.connect(holder1).claimYield(assetId);
      await distributor.connect(holder2).claimYield(assetId);

      // Deposit more yield
      await distributor.connect(admin).depositYield(assetId, { value: YIELD_AMOUNT });

      // Both should have new pending yield
      const pending1 = await distributor.getPendingYield(assetId, holder1.address);
      const pending2 = await distributor.getPendingYield(assetId, holder2.address);

      expect(pending1).to.be.gt(0);
      expect(pending2).to.be.gt(0);
    });
  });

  describe("Yield History", function () {
    it("Should track claimed yield per user", async function () {
      await distributor.connect(admin).depositYield(assetId, { value: YIELD_AMOUNT });
      await distributor.connect(holder1).claimYield(assetId);

      const claimedAmount = await distributor.getTotalClaimedYield(assetId, holder1.address);
      expect(claimedAmount).to.be.gt(0);
    });

    it("Should track total distributed yield", async function () {
      await distributor.connect(admin).depositYield(assetId, { value: YIELD_AMOUNT });
      await distributor.connect(holder1).claimYield(assetId);
      await distributor.connect(holder2).claimYield(assetId);

      const totalDistributed = await distributor.getTotalDistributedYield(assetId);
      expect(totalDistributed).to.be.closeTo(YIELD_AMOUNT, ethers.parseEther("0.01"));
    });
  });

  describe("Access Control", function () {
    it("Should only allow operator to deposit yield", async function () {
      await expect(
        distributor.connect(holder1).depositYield(assetId, { value: YIELD_AMOUNT })
      ).to.be.reverted;
    });
  });
});
