import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type {
  RWAYieldDistributor,
  RWAAssetFactory,
  RWACompliance,
  RWAValuation,
  FractionalRWAToken,
  MockERC20,
} from "../../typechain-types/index.js";

describe("RWAYieldDistributor", function () {
  let distributor: RWAYieldDistributor;
  let factory: RWAAssetFactory;
  let compliance: RWACompliance;
  let valuation: RWAValuation;
  let token: FractionalRWAToken;
  let paymentToken: MockERC20;

  let admin: SignerWithAddress;
  let issuer: SignerWithAddress;
  let holder1: SignerWithAddress;
  let holder2: SignerWithAddress;
  let holder3: SignerWithAddress;

  const DISTRIBUTOR_ROLE = ethers.keccak256(ethers.toUtf8Bytes("DISTRIBUTOR_ROLE"));
  const COMPLIANCE_ROLE = ethers.keccak256(ethers.toUtf8Bytes("COMPLIANCE_ROLE"));
  const VERIFIER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("VERIFIER_ROLE"));
  const MINTER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("MINTER_ROLE"));

  let assetId: number;
  const TOTAL_SUPPLY = ethers.parseEther("1000000"); // 1M tokens
  const YIELD_AMOUNT = ethers.parseEther("0.1"); // 0.1 ETH in yield (affordable for testing)

  beforeEach(async function () {
    [admin, issuer, holder1, holder2, holder3] = await ethers.getSigners();

    // Deploy RWACompliance
    const ComplianceFactory = await ethers.getContractFactory("RWACompliance");
    compliance = await ComplianceFactory.deploy(admin.address);

    // Deploy RWAValuation
    const ValuationFactory = await ethers.getContractFactory("RWAValuation");
    valuation = await ValuationFactory.deploy(admin.address);

    // Deploy RWAAssetFactory
    const FactoryContract = await ethers.getContractFactory("RWAAssetFactory");
    factory = await FactoryContract.deploy(
      admin.address,
      await compliance.getAddress(),
      await valuation.getAddress()
    );

    // Deploy RWAYieldDistributor
    const DistributorFactory = await ethers.getContractFactory("RWAYieldDistributor");
    distributor = await DistributorFactory.deploy(
      admin.address,
      await factory.getAddress()
    );

    // Deploy MockERC20 for payment token
    const MockERC20Factory = await ethers.getContractFactory("MockERC20");
    paymentToken = await MockERC20Factory.deploy(
      "USD Coin",
      "USDC",
      ethers.parseEther("1000000")
    );

    // Grant roles
    await compliance.grantRole(VERIFIER_ROLE, admin.address);
    await factory.grantRole(COMPLIANCE_ROLE, admin.address);

    // Verify all participants
    const oneYear = 365 * 24 * 60 * 60;
    for (const user of [issuer, holder1, holder2, holder3]) {
      await compliance.verifyInvestor(
        user.address,
        1, // AccreditationTier.Accredited
        "US",
        oneYear,
        ethers.keccak256(ethers.toUtf8Bytes(`${user.address}-kyc`))
      );
    }

    // Create and activate an asset
    await factory
      .connect(issuer)
      .createAsset(
        "Rental Property",
        "RENT",
        0, // RealEstate
        ethers.parseEther("1000000"), // $1M
        0,
        "QmXYZ",
        "QmABC"
      );

    assetId = 1;

    // Fractionalize and activate
    await factory.connect(issuer).fractionalizeAsset(assetId, TOTAL_SUPPLY);
    await factory.activateAsset(assetId);

    // Get token contract
    const tokenAddress = await factory.getFractionalToken(assetId);
    token = await ethers.getContractAt("FractionalRWAToken", tokenAddress);

    // Mint tokens to issuer and distribute to holders
    await token.connect(issuer).grantRole(MINTER_ROLE, issuer.address);
    await token.connect(issuer).mint(issuer.address, TOTAL_SUPPLY);

    // Distribute tokens: 50% holder1, 30% holder2, 20% holder3
    await token.connect(issuer).transfer(holder1.address, TOTAL_SUPPLY * 50n / 100n);
    await token.connect(issuer).transfer(holder2.address, TOTAL_SUPPLY * 30n / 100n);
    await token.connect(issuer).transfer(holder3.address, TOTAL_SUPPLY * 20n / 100n);

    // Give issuer payment tokens and ETH
    await paymentToken.transfer(issuer.address, ethers.parseEther("100000"));
    await paymentToken
      .connect(issuer)
      .approve(await distributor.getAddress(), ethers.MaxUint256);
  });

  describe("Deployment", function () {
    it("Should set the correct admin", async function () {
      expect(
        await distributor.hasRole(await distributor.DEFAULT_ADMIN_ROLE(), admin.address)
      ).to.be.true;
    });

    it("Should set the correct asset factory", async function () {
      expect(await distributor.assetFactory()).to.equal(await factory.getAddress());
    });

    it("Should initialize counters to zero", async function () {
      expect(await distributor.totalDistributions()).to.equal(0);
      expect(await distributor.totalYieldDistributed()).to.equal(0);
    });

    it("Should revert with invalid admin", async function () {
      const DistributorFactory = await ethers.getContractFactory("RWAYieldDistributor");
      await expect(
        DistributorFactory.deploy(ethers.ZeroAddress, await factory.getAddress())
      ).to.be.revertedWith("Invalid admin");
    });

    it("Should revert with invalid factory", async function () {
      const DistributorFactory = await ethers.getContractFactory("RWAYieldDistributor");
      await expect(
        DistributorFactory.deploy(admin.address, ethers.ZeroAddress)
      ).to.be.revertedWith("Invalid factory");
    });
  });

  describe("Yield Deposit - ETH", function () {
    it("Should deposit ETH yield successfully", async function () {
      const claimPeriod = 30 * 24 * 60 * 60; // 30 days

      const tx = await distributor
        .connect(issuer)
        .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, claimPeriod, {
          value: YIELD_AMOUNT,
        });

      await expect(tx)
        .to.emit(distributor, "YieldDeposited")
        .withArgs(1, assetId, ethers.ZeroAddress, YIELD_AMOUNT, await ethers.provider.getBlock("latest").then(b => b!.timestamp + claimPeriod));

      expect(await distributor.totalDistributions()).to.equal(1);
      expect(await distributor.totalYieldDistributed()).to.equal(YIELD_AMOUNT);
    });

    it("Should store correct distribution data", async function () {
      const claimPeriod = 30 * 24 * 60 * 60;

      await distributor
        .connect(issuer)
        .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, claimPeriod, {
          value: YIELD_AMOUNT,
        });

      const dist = await distributor.getDistribution(1);

      expect(dist.distributionId).to.equal(1);
      expect(dist.assetId).to.equal(assetId);
      expect(dist.paymentToken).to.equal(ethers.ZeroAddress);
      expect(dist.totalAmount).to.equal(YIELD_AMOUNT);
      expect(dist.totalSupply).to.equal(TOTAL_SUPPLY);
      expect(dist.totalClaimed).to.equal(0);
      expect(dist.finalized).to.be.false;
    });

    it("Should track asset distributions", async function () {
      const claimPeriod = 30 * 24 * 60 * 60;

      await distributor
        .connect(issuer)
        .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, claimPeriod, {
          value: YIELD_AMOUNT,
        });

      await distributor
        .connect(issuer)
        .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, claimPeriod, {
          value: YIELD_AMOUNT,
        });

      const distributions = await distributor.getAssetDistributions(assetId);
      expect(distributions.length).to.equal(2);
      expect(distributions[0]).to.equal(1);
      expect(distributions[1]).to.equal(2);
    });

    it("Should revert if asset not active", async function () {
      await factory.updateAssetStatus(assetId, 2); // Suspended

      await expect(
        distributor
          .connect(issuer)
          .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, 30 * 24 * 60 * 60, {
            value: YIELD_AMOUNT,
          })
      ).to.be.revertedWith("Asset not active");
    });

    it("Should revert with zero amount", async function () {
      await expect(
        distributor
          .connect(issuer)
          .depositYield(assetId, ethers.ZeroAddress, 0, 30 * 24 * 60 * 60, {
            value: 0,
          })
      ).to.be.revertedWith("Invalid amount");
    });

    it("Should revert with invalid claim period", async function () {
      await expect(
        distributor
          .connect(issuer)
          .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, 0, {
            value: YIELD_AMOUNT,
          })
      ).to.be.revertedWith("Invalid claim period");
    });

    it("Should revert if ETH amount incorrect", async function () {
      await expect(
        distributor
          .connect(issuer)
          .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, 30 * 24 * 60 * 60, {
            value: YIELD_AMOUNT - 1n,
          })
      ).to.be.revertedWith("Incorrect ETH amount");
    });
  });

  describe("Yield Deposit - ERC20", function () {
    it("Should deposit ERC20 yield successfully", async function () {
      const claimPeriod = 30 * 24 * 60 * 60;

      const tx = await distributor
        .connect(issuer)
        .depositYield(
          assetId,
          await paymentToken.getAddress(),
          YIELD_AMOUNT,
          claimPeriod
        );

      await expect(tx).to.emit(distributor, "YieldDeposited");

      const distributorBalance = await paymentToken.balanceOf(
        await distributor.getAddress()
      );
      expect(distributorBalance).to.equal(YIELD_AMOUNT);
    });

    it("Should revert if ETH sent with ERC20 payment", async function () {
      await expect(
        distributor
          .connect(issuer)
          .depositYield(
            assetId,
            await paymentToken.getAddress(),
            YIELD_AMOUNT,
            30 * 24 * 60 * 60,
            { value: ethers.parseEther("1") }
          )
      ).to.be.revertedWith("ETH not accepted for token payment");
    });
  });

  describe("Yield Claiming", function () {
    beforeEach(async function () {
      const claimPeriod = 30 * 24 * 60 * 60;
      await distributor
        .connect(issuer)
        .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, claimPeriod, {
          value: YIELD_AMOUNT,
        });
    });

    it("Should allow holder to claim yield", async function () {
      const holder1BalanceBefore = await ethers.provider.getBalance(holder1.address);

      const tx = await distributor.connect(holder1).claimYield(1);
      const receipt = await tx.wait();
      const gasUsed = receipt!.gasUsed * receipt!.gasPrice;

      const holder1BalanceAfter = await ethers.provider.getBalance(holder1.address);

      // Holder1 has 50% of tokens, should get 50% of yield
      const expectedYield = YIELD_AMOUNT * 50n / 100n;
      const received = holder1BalanceAfter - holder1BalanceBefore + gasUsed;

      expect(received).to.equal(expectedYield);
    });

    it("Should emit YieldClaimed event", async function () {
      const expectedYield = YIELD_AMOUNT * 50n / 100n;

      await expect(distributor.connect(holder1).claimYield(1))
        .to.emit(distributor, "YieldClaimed")
        .withArgs(1, holder1.address, expectedYield);
    });

    it("Should update claim status", async function () {
      await distributor.connect(holder1).claimYield(1);

      const claim = await distributor.getUserClaim(1, holder1.address);
      expect(claim.claimed).to.be.true;
      expect(claim.amount).to.equal(YIELD_AMOUNT * 50n / 100n);
    });

    it("Should update distribution claimed amount", async function () {
      await distributor.connect(holder1).claimYield(1);

      const dist = await distributor.getDistribution(1);
      expect(dist.totalClaimed).to.equal(YIELD_AMOUNT * 50n / 100n);
    });

    it("Should calculate correct pro-rata amounts", async function () {
      // Holder1: 50%, Holder2: 30%, Holder3: 20%
      const holder1Claimable = await distributor.getClaimableAmount(1, holder1.address);
      const holder2Claimable = await distributor.getClaimableAmount(1, holder2.address);
      const holder3Claimable = await distributor.getClaimableAmount(1, holder3.address);

      expect(holder1Claimable).to.equal(YIELD_AMOUNT * 50n / 100n);
      expect(holder2Claimable).to.equal(YIELD_AMOUNT * 30n / 100n);
      expect(holder3Claimable).to.equal(YIELD_AMOUNT * 20n / 100n);
    });

    it("Should revert if already claimed", async function () {
      await distributor.connect(holder1).claimYield(1);

      await expect(
        distributor.connect(holder1).claimYield(1)
      ).to.be.revertedWith("Already claimed");
    });

    it("Should revert if no tokens held", async function () {
      const [, , , , , , noTokensHolder] = await ethers.getSigners();

      await expect(
        distributor.connect(noTokensHolder).claimYield(1)
      ).to.be.revertedWith("No tokens held");
    });

    it("Should revert if claim period expired", async function () {
      // Fast forward past claim deadline
      await ethers.provider.send("evm_increaseTime", [31 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      await expect(
        distributor.connect(holder1).claimYield(1)
      ).to.be.revertedWith("Claim period expired");
    });

    it("Should revert if distribution finalized", async function () {
      // Fast forward and finalize
      await ethers.provider.send("evm_increaseTime", [31 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      await distributor.finalizeDistribution(1);

      await expect(
        distributor.connect(holder1).claimYield(1)
      ).to.be.revertedWith("Distribution finalized");
    });
  });

  describe("Batch Claiming", function () {
    beforeEach(async function () {
      const claimPeriod = 30 * 24 * 60 * 60;

      // Create 3 distributions
      for (let i = 0; i < 3; i++) {
        await distributor
          .connect(issuer)
          .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, claimPeriod, {
            value: YIELD_AMOUNT,
          });
      }
    });

    it("Should batch claim multiple distributions", async function () {
      const holder1BalanceBefore = await ethers.provider.getBalance(holder1.address);

      const tx = await distributor.connect(holder1).batchClaimYield([1, 2, 3]);
      const receipt = await tx.wait();
      const gasUsed = receipt!.gasUsed * receipt!.gasPrice;

      const holder1BalanceAfter = await ethers.provider.getBalance(holder1.address);

      // Should receive 50% of yield from 3 distributions
      const expectedTotal = (YIELD_AMOUNT * 50n / 100n) * 3n;
      const received = holder1BalanceAfter - holder1BalanceBefore + gasUsed;

      expect(received).to.equal(expectedTotal);
    });

    it("Should skip already claimed distributions", async function () {
      // Claim first distribution
      await distributor.connect(holder1).claimYield(1);

      const holder1BalanceBefore = await ethers.provider.getBalance(holder1.address);

      const tx = await distributor.connect(holder1).batchClaimYield([1, 2, 3]);
      const receipt = await tx.wait();
      const gasUsed = receipt!.gasUsed * receipt!.gasPrice;

      const holder1BalanceAfter = await ethers.provider.getBalance(holder1.address);

      // Should only receive from distributions 2 and 3
      const expectedTotal = (YIELD_AMOUNT * 50n / 100n) * 2n;
      const received = holder1BalanceAfter - holder1BalanceBefore + gasUsed;

      expect(received).to.equal(expectedTotal);
    });

    it("Should skip invalid distributions", async function () {
      // Should not revert, just skip
      await expect(
        distributor.connect(holder1).batchClaimYield([1, 2, 999])
      ).to.not.be.reverted;
    });
  });

  describe("Distribution Finalization", function () {
    beforeEach(async function () {
      const claimPeriod = 30 * 24 * 60 * 60;
      await distributor
        .connect(issuer)
        .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, claimPeriod, {
          value: YIELD_AMOUNT,
        });
    });

    it("Should finalize distribution after claim period", async function () {
      // Fast forward past claim deadline
      await ethers.provider.send("evm_increaseTime", [31 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      const tx = await distributor.finalizeDistribution(1);

      await expect(tx).to.emit(distributor, "DistributionFinalized");

      const dist = await distributor.getDistribution(1);
      expect(dist.finalized).to.be.true;
    });

    it("Should return unclaimed yield to issuer", async function () {
      // Only holder1 claims (50%)
      await distributor.connect(holder1).claimYield(1);

      // Fast forward and finalize
      await ethers.provider.send("evm_increaseTime", [31 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      const issuerBalanceBefore = await ethers.provider.getBalance(issuer.address);

      const tx = await distributor.finalizeDistribution(1);
      const receipt = await tx.wait();
      const gasUsed = receipt!.gasUsed * receipt!.gasPrice;

      const issuerBalanceAfter = await ethers.provider.getBalance(issuer.address);

      // Issuer should receive 50% back (unclaimed by holder2 and holder3)
      const unclaimed = YIELD_AMOUNT * 50n / 100n;
      const received = issuerBalanceAfter - issuerBalanceBefore + gasUsed;

      // Use closeTo to account for rounding in gas calculations
      expect(received).to.be.closeTo(unclaimed, ethers.parseEther("0.001"));
    });

    it("Should emit UnclaimedYieldReclaimed event", async function () {
      await ethers.provider.send("evm_increaseTime", [31 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      await expect(distributor.finalizeDistribution(1))
        .to.emit(distributor, "UnclaimedYieldReclaimed")
        .withArgs(1, YIELD_AMOUNT, issuer.address);
    });

    it("Should revert if claim period not ended", async function () {
      await expect(distributor.finalizeDistribution(1)).to.be.revertedWith(
        "Claim period not ended"
      );
    });

    it("Should revert if already finalized", async function () {
      await ethers.provider.send("evm_increaseTime", [31 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      await distributor.finalizeDistribution(1);

      await expect(distributor.finalizeDistribution(1)).to.be.revertedWith(
        "Already finalized"
      );
    });

    it("Should revert if not distributor role", async function () {
      await ethers.provider.send("evm_increaseTime", [31 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      await expect(
        distributor.connect(holder1).finalizeDistribution(1)
      ).to.be.revertedWithCustomError(distributor, "AccessControlUnauthorizedAccount");
    });
  });

  describe("View Functions", function () {
    beforeEach(async function () {
      const claimPeriod = 30 * 24 * 60 * 60;
      await distributor
        .connect(issuer)
        .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, claimPeriod, {
          value: YIELD_AMOUNT,
        });
    });

    it("Should return correct claimable amount", async function () {
      const claimable = await distributor.getClaimableAmount(1, holder1.address);
      expect(claimable).to.equal(YIELD_AMOUNT * 50n / 100n);
    });

    it("Should return zero for claimed distributions", async function () {
      await distributor.connect(holder1).claimYield(1);

      const claimable = await distributor.getClaimableAmount(1, holder1.address);
      expect(claimable).to.equal(0);
    });

    it("Should return zero for expired distributions", async function () {
      await ethers.provider.send("evm_increaseTime", [31 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      const claimable = await distributor.getClaimableAmount(1, holder1.address);
      expect(claimable).to.equal(0);
    });

    it("Should return unclaimed amount", async function () {
      await distributor.connect(holder1).claimYield(1);

      const unclaimed = await distributor.getUnclaimedAmount(1);
      expect(unclaimed).to.equal(YIELD_AMOUNT * 50n / 100n);
    });

    it("Should calculate total claimable across distributions", async function () {
      // Create another distribution
      await distributor
        .connect(issuer)
        .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, 30 * 24 * 60 * 60, {
          value: YIELD_AMOUNT,
        });

      const totalClaimable = await distributor.getTotalClaimable(
        holder1.address,
        assetId
      );

      // 50% of 2 distributions
      expect(totalClaimable).to.equal((YIELD_AMOUNT * 50n / 100n) * 2n);
    });
  });

  describe("Admin Functions", function () {
    it("Should pause and unpause", async function () {
      await distributor.pause();
      expect(await distributor.paused()).to.be.true;

      await distributor.unpause();
      expect(await distributor.paused()).to.be.false;
    });

    it("Should revert deposit when paused", async function () {
      await distributor.pause();

      await expect(
        distributor
          .connect(issuer)
          .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, 30 * 24 * 60 * 60, {
            value: YIELD_AMOUNT,
          })
      ).to.be.revertedWithCustomError(distributor, "EnforcedPause");
    });

    it("Should revert claim when paused", async function () {
      await distributor
        .connect(issuer)
        .depositYield(assetId, ethers.ZeroAddress, YIELD_AMOUNT, 30 * 24 * 60 * 60, {
          value: YIELD_AMOUNT,
        });

      await distributor.pause();

      await expect(
        distributor.connect(holder1).claimYield(1)
      ).to.be.revertedWithCustomError(distributor, "EnforcedPause");
    });

    it("Should emergency withdraw ETH", async function () {
      // Send some ETH to contract
      await admin.sendTransaction({
        to: await distributor.getAddress(),
        value: ethers.parseEther("1"),
      });

      const recipientBalanceBefore = await ethers.provider.getBalance(holder1.address);

      await distributor.emergencyWithdraw(
        ethers.ZeroAddress,
        ethers.parseEther("1"),
        holder1.address
      );

      const recipientBalanceAfter = await ethers.provider.getBalance(holder1.address);

      expect(recipientBalanceAfter - recipientBalanceBefore).to.equal(
        ethers.parseEther("1")
      );
    });

    it("Should emergency withdraw ERC20", async function () {
      // Send some tokens to contract
      await paymentToken.transfer(
        await distributor.getAddress(),
        ethers.parseEther("1000")
      );

      await distributor.emergencyWithdraw(
        await paymentToken.getAddress(),
        ethers.parseEther("1000"),
        holder1.address
      );

      const balance = await paymentToken.balanceOf(holder1.address);
      expect(balance).to.equal(ethers.parseEther("1000"));
    });
  });
});
