import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type { RWACompliance } from "../../typechain-types/index.js";

describe("RWACompliance", function () {
  let compliance: RWACompliance;

  let admin: SignerWithAddress;
  let verifier: SignerWithAddress;
  let complianceOfficer: SignerWithAddress;
  let investor1: SignerWithAddress;
  let investor2: SignerWithAddress;
  let investor3: SignerWithAddress;

  const COMPLIANCE_OFFICER_ROLE = ethers.keccak256(
    ethers.toUtf8Bytes("COMPLIANCE_OFFICER_ROLE")
  );
  const VERIFIER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("VERIFIER_ROLE"));

  const ONE_YEAR = 365 * 24 * 60 * 60;

  beforeEach(async function () {
    [admin, verifier, complianceOfficer, investor1, investor2, investor3] =
      await ethers.getSigners();

    // Deploy RWACompliance
    const ComplianceFactory = await ethers.getContractFactory("RWACompliance");
    compliance = await ComplianceFactory.deploy(admin.address);

    // Grant roles
    await compliance.grantRole(VERIFIER_ROLE, verifier.address);
    await compliance.grantRole(COMPLIANCE_OFFICER_ROLE, complianceOfficer.address);
  });

  describe("Deployment", function () {
    it("Should set the correct admin", async function () {
      expect(
        await compliance.hasRole(await compliance.DEFAULT_ADMIN_ROLE(), admin.address)
      ).to.be.true;
    });

    it("Should grant initial roles", async function () {
      expect(await compliance.hasRole(COMPLIANCE_OFFICER_ROLE, admin.address)).to.be.true;
      expect(await compliance.hasRole(VERIFIER_ROLE, admin.address)).to.be.true;
    });

    it("Should initialize counters to zero", async function () {
      expect(await compliance.totalVerifiedInvestors()).to.equal(0);
      expect(await compliance.totalAssetsWithCompliance()).to.equal(0);
    });

    it("Should revert with invalid admin", async function () {
      const ComplianceFactory = await ethers.getContractFactory("RWACompliance");
      await expect(
        ComplianceFactory.deploy(ethers.ZeroAddress)
      ).to.be.revertedWith("Invalid admin address");
    });
  });

  describe("Investor Verification", function () {
    it("Should verify investor successfully", async function () {
      const tx = await compliance
        .connect(verifier)
        .verifyInvestor(
          investor1.address,
          1, // AccreditationTier.Accredited
          "US",
          ONE_YEAR,
          ethers.keccak256(ethers.toUtf8Bytes("kyc-doc-1"))
        );

      await expect(tx).to.emit(compliance, "InvestorVerified");

      expect(await compliance.totalVerifiedInvestors()).to.equal(1);
    });

    it("Should store correct investor profile", async function () {
      await compliance
        .connect(verifier)
        .verifyInvestor(
          investor1.address,
          1, // Accredited
          "US",
          ONE_YEAR,
          ethers.keccak256(ethers.toUtf8Bytes("kyc-doc-1"))
        );

      const profile = await compliance.getInvestorProfile(investor1.address);

      expect(profile.status).to.equal(2); // VerificationStatus.Approved
      expect(profile.tier).to.equal(1); // Accredited
      expect(profile.jurisdiction).to.equal("US");
      expect(profile.verifier).to.equal(verifier.address);
    });

    it("Should check if investor is verified", async function () {
      await compliance
        .connect(verifier)
        .verifyInvestor(
          investor1.address,
          1,
          "US",
          ONE_YEAR,
          ethers.keccak256(ethers.toUtf8Bytes("kyc-doc-1"))
        );

      expect(await compliance.isInvestorVerified(investor1.address)).to.be.true;
    });

    it("Should support multiple accreditation tiers", async function () {
      // Retail = 0, Accredited = 1, Institutional = 2, QualifiedPurchaser = 3
      for (let tier = 0; tier <= 3; tier++) {
        const investor = [investor1, investor2, investor3, admin][tier];
        await compliance
          .connect(verifier)
          .verifyInvestor(
            investor.address,
            tier,
            "US",
            ONE_YEAR,
            ethers.keccak256(ethers.toUtf8Bytes(`kyc-${tier}`))
          );

        const profile = await compliance.getInvestorProfile(investor.address);
        expect(profile.tier).to.equal(tier);
      }
    });

    it("Should revert with invalid investor address", async function () {
      await expect(
        compliance
          .connect(verifier)
          .verifyInvestor(
            ethers.ZeroAddress,
            1,
            "US",
            ONE_YEAR,
            ethers.keccak256(ethers.toUtf8Bytes("kyc-doc"))
          )
      ).to.be.revertedWith("Invalid investor address");
    });

    it("Should revert with invalid jurisdiction code", async function () {
      await expect(
        compliance
          .connect(verifier)
          .verifyInvestor(
            investor1.address,
            1,
            "USA", // Should be 2 characters
            ONE_YEAR,
            ethers.keccak256(ethers.toUtf8Bytes("kyc-doc"))
          )
      ).to.be.revertedWith("Invalid jurisdiction code");
    });

    it("Should revert if validity period too long", async function () {
      await expect(
        compliance
          .connect(verifier)
          .verifyInvestor(
            investor1.address,
            1,
            "US",
            731 * 24 * 60 * 60, // > 730 days
            ethers.keccak256(ethers.toUtf8Bytes("kyc-doc"))
          )
      ).to.be.revertedWith("Validity period too long");
    });

    it("Should revert without KYC document hash", async function () {
      await expect(
        compliance
          .connect(verifier)
          .verifyInvestor(investor1.address, 1, "US", ONE_YEAR, ethers.ZeroHash)
      ).to.be.revertedWith("KYC document hash required");
    });

    it("Should revert if not verifier role", async function () {
      await expect(
        compliance
          .connect(investor1)
          .verifyInvestor(
            investor2.address,
            1,
            "US",
            ONE_YEAR,
            ethers.keccak256(ethers.toUtf8Bytes("kyc-doc"))
          )
      ).to.be.revertedWithCustomError(compliance, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Investor Status Management", function () {
    beforeEach(async function () {
      await compliance
        .connect(verifier)
        .verifyInvestor(
          investor1.address,
          1,
          "US",
          ONE_YEAR,
          ethers.keccak256(ethers.toUtf8Bytes("kyc-doc-1"))
        );
    });

    it("Should update investor status", async function () {
      const tx = await compliance
        .connect(complianceOfficer)
        .updateInvestorStatus(investor1.address, 5); // VerificationStatus.Suspended

      await expect(tx)
        .to.emit(compliance, "InvestorStatusChanged")
        .withArgs(investor1.address, 2, 5); // Approved -> Suspended

      const profile = await compliance.getInvestorProfile(investor1.address);
      expect(profile.status).to.equal(5);
    });

    it("Should update verified investors counter", async function () {
      expect(await compliance.totalVerifiedInvestors()).to.equal(1);

      await compliance
        .connect(complianceOfficer)
        .updateInvestorStatus(investor1.address, 5); // Suspended

      expect(await compliance.totalVerifiedInvestors()).to.equal(0);
    });

    it("Should revert if status unchanged", async function () {
      await expect(
        compliance
          .connect(complianceOfficer)
          .updateInvestorStatus(investor1.address, 2) // Already Approved
      ).to.be.revertedWith("Status unchanged");
    });

    it("Should revert if not compliance officer", async function () {
      await expect(
        compliance.connect(investor2).updateInvestorStatus(investor1.address, 5)
      ).to.be.revertedWithCustomError(compliance, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Asset Compliance Configuration", function () {
    const assetId = 1;

    it("Should set asset compliance requirements", async function () {
      const tx = await compliance
        .connect(complianceOfficer)
        .setAssetCompliance(
          assetId,
          1, // AccreditationTier.Accredited
          true, // requiresKYC
          ethers.parseEther("10000"), // minInvestmentAmount
          100 // maxInvestors
        );

      await expect(tx)
        .to.emit(compliance, "AssetComplianceSet")
        .withArgs(assetId, 1, true);

      expect(await compliance.totalAssetsWithCompliance()).to.equal(1);
    });

    it("Should check if asset requires KYC", async function () {
      await compliance
        .connect(complianceOfficer)
        .setAssetCompliance(assetId, 1, true, ethers.parseEther("10000"), 100);

      expect(await compliance.doesAssetRequireKYC(assetId)).to.be.true;
    });

    it("Should get asset minimum tier", async function () {
      await compliance
        .connect(complianceOfficer)
        .setAssetCompliance(assetId, 2, true, ethers.parseEther("10000"), 100); // Institutional

      expect(await compliance.getAssetMinTier(assetId)).to.equal(2);
    });

    it("Should revert if not compliance officer", async function () {
      await expect(
        compliance
          .connect(investor1)
          .setAssetCompliance(assetId, 1, true, ethers.parseEther("10000"), 100)
      ).to.be.revertedWithCustomError(compliance, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Jurisdiction Restrictions", function () {
    const assetId = 1;

    beforeEach(async function () {
      await compliance
        .connect(complianceOfficer)
        .setAssetCompliance(assetId, 1, true, 0, 0);
    });

    it("Should update jurisdiction restriction", async function () {
      const tx = await compliance
        .connect(complianceOfficer)
        .updateJurisdictionRestriction(assetId, "US", true); // Allow US

      await expect(tx).to.emit(compliance, "JurisdictionUpdated").withArgs(assetId, "US", true);
    });

    it("Should check if jurisdiction is allowed", async function () {
      await compliance
        .connect(complianceOfficer)
        .updateJurisdictionRestriction(assetId, "US", true);

      expect(await compliance.isJurisdictionAllowed(assetId, "US")).to.be.true;
    });

    it("Should block specific jurisdictions", async function () {
      await compliance
        .connect(complianceOfficer)
        .updateJurisdictionRestriction(assetId, "CN", false); // Block CN

      expect(await compliance.isJurisdictionAllowed(assetId, "CN")).to.be.false;
    });

    it("Should allow all jurisdictions by default", async function () {
      // Asset 2 has no restrictions set
      expect(await compliance.isJurisdictionAllowed(2, "US")).to.be.true;
      expect(await compliance.isJurisdictionAllowed(2, "UK")).to.be.true;
    });

    it("Should revert with invalid jurisdiction code", async function () {
      await expect(
        compliance
          .connect(complianceOfficer)
          .updateJurisdictionRestriction(assetId, "USA", true)
      ).to.be.revertedWith("Invalid jurisdiction code");
    });
  });

  describe("Compliance Checks", function () {
    const assetId = 1;

    beforeEach(async function () {
      // Set up compliance requirements
      await compliance
        .connect(complianceOfficer)
        .setAssetCompliance(
          assetId,
          1, // Accredited tier required
          true, // KYC required
          0,
          0
        );

      // Verify investor1 as Accredited in US
      await compliance
        .connect(verifier)
        .verifyInvestor(
          investor1.address,
          1, // Accredited
          "US",
          ONE_YEAR,
          ethers.keccak256(ethers.toUtf8Bytes("kyc-1"))
        );

      // Verify investor2 as Retail in UK
      await compliance
        .connect(verifier)
        .verifyInvestor(
          investor2.address,
          0, // Retail
          "UK",
          ONE_YEAR,
          ethers.keccak256(ethers.toUtf8Bytes("kyc-2"))
        );
    });

    it("Should allow compliant investor to invest", async function () {
      expect(await compliance.canInvestInAsset(investor1.address, assetId)).to.be.true;
    });

    it("Should deny investor with insufficient tier", async function () {
      // investor2 is Retail (0), but asset requires Accredited (1)
      expect(await compliance.canInvestInAsset(investor2.address, assetId)).to.be.false;
    });

    it("Should deny unverified investor", async function () {
      expect(await compliance.canInvestInAsset(investor3.address, assetId)).to.be.false;
    });

    it("Should deny investor from blocked jurisdiction", async function () {
      await compliance
        .connect(complianceOfficer)
        .updateJurisdictionRestriction(assetId, "UK", false); // Block UK

      expect(await compliance.canInvestInAsset(investor2.address, assetId)).to.be.false;
    });

    it("Should allow transfer between compliant investors", async function () {
      // Verify investor2 as Accredited too
      await compliance
        .connect(verifier)
        .verifyInvestor(
          investor2.address,
          1,
          "US",
          ONE_YEAR,
          ethers.keccak256(ethers.toUtf8Bytes("kyc-2-updated"))
        );

      expect(
        await compliance.canTransferTokens(
          investor1.address,
          investor2.address,
          assetId,
          ethers.parseEther("100")
        )
      ).to.be.true;
    });

    it("Should deny transfer to non-compliant investor", async function () {
      expect(
        await compliance.canTransferTokens(
          investor1.address,
          investor3.address, // Not verified
          assetId,
          ethers.parseEther("100")
        )
      ).to.be.false;
    });

    it("Should allow minting to compliant investor", async function () {
      expect(
        await compliance.canTransferTokens(
          ethers.ZeroAddress, // Minting
          investor1.address,
          assetId,
          ethers.parseEther("100")
        )
      ).to.be.true;
    });

    it("Should allow burning from any investor", async function () {
      expect(
        await compliance.canTransferTokens(
          investor1.address,
          ethers.ZeroAddress, // Burning
          assetId,
          ethers.parseEther("100")
        )
      ).to.be.true;
    });
  });

  describe("Max Investors Limit", function () {
    const assetId = 1;

    beforeEach(async function () {
      // Set max investors to 2
      await compliance
        .connect(complianceOfficer)
        .setAssetCompliance(
          assetId,
          1, // Accredited
          true, // KYC required
          0,
          2 // Max 2 investors
        );

      // Verify 3 investors
      for (const investor of [investor1, investor2, investor3]) {
        await compliance
          .connect(verifier)
          .verifyInvestor(
            investor.address,
            1,
            "US",
            ONE_YEAR,
            ethers.keccak256(ethers.toUtf8Bytes(`kyc-${investor.address}`))
          );
      }
    });

    it("Should allow investment within limit", async function () {
      expect(await compliance.canInvestInAsset(investor1.address, assetId)).to.be.true;
    });

    it("Should track investor participation", async function () {
      await compliance.recordInvestorParticipation(assetId, investor1.address);
      await compliance.recordInvestorParticipation(assetId, investor2.address);

      expect(await compliance.getAssetInvestorCount(assetId)).to.equal(2);
    });

    it("Should deny new investor when limit reached", async function () {
      await compliance.recordInvestorParticipation(assetId, investor1.address);
      await compliance.recordInvestorParticipation(assetId, investor2.address);

      // investor3 is new and limit is reached
      expect(await compliance.canInvestInAsset(investor3.address, assetId)).to.be.false;
    });

    it("Should allow existing investor even when limit reached", async function () {
      await compliance.recordInvestorParticipation(assetId, investor1.address);
      await compliance.recordInvestorParticipation(assetId, investor2.address);

      // investor1 is already in, so they can still invest
      expect(await compliance.canInvestInAsset(investor1.address, assetId)).to.be.true;
    });

    it("Should remove investor participation", async function () {
      await compliance.recordInvestorParticipation(assetId, investor1.address);
      await compliance.removeInvestorParticipation(assetId, investor1.address);

      expect(await compliance.getAssetInvestorCount(assetId)).to.equal(0);
    });

    it("Should not double-count investors", async function () {
      await compliance.recordInvestorParticipation(assetId, investor1.address);
      await compliance.recordInvestorParticipation(assetId, investor1.address); // Second time

      expect(await compliance.getAssetInvestorCount(assetId)).to.equal(1);
    });
  });

  describe("Verification Expiry", function () {
    it("Should expire after validity period", async function () {
      await compliance
        .connect(verifier)
        .verifyInvestor(
          investor1.address,
          1,
          "US",
          24 * 60 * 60, // 1 day
          ethers.keccak256(ethers.toUtf8Bytes("kyc-1"))
        );

      expect(await compliance.isInvestorVerified(investor1.address)).to.be.true;

      // Fast forward 2 days
      await ethers.provider.send("evm_increaseTime", [2 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      expect(await compliance.isInvestorVerified(investor1.address)).to.be.false;
    });

    it("Should not allow investment after expiry", async function () {
      const assetId = 1;
      await compliance
        .connect(complianceOfficer)
        .setAssetCompliance(assetId, 1, true, 0, 0);

      await compliance
        .connect(verifier)
        .verifyInvestor(
          investor1.address,
          1,
          "US",
          24 * 60 * 60,
          ethers.keccak256(ethers.toUtf8Bytes("kyc-1"))
        );

      // Fast forward
      await ethers.provider.send("evm_increaseTime", [2 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      expect(await compliance.canInvestInAsset(investor1.address, assetId)).to.be.false;
    });
  });

  describe("Edge Cases", function () {
    it("Should handle asset without compliance requirements", async function () {
      // Asset 99 has no compliance set
      expect(await compliance.canInvestInAsset(investor1.address, 99)).to.be.true;
    });

    it("Should handle zero max investors (unlimited)", async function () {
      const assetId = 1;
      await compliance
        .connect(complianceOfficer)
        .setAssetCompliance(assetId, 1, true, 0, 0); // maxInvestors = 0 means unlimited

      await compliance
        .connect(verifier)
        .verifyInvestor(
          investor1.address,
          1,
          "US",
          ONE_YEAR,
          ethers.keccak256(ethers.toUtf8Bytes("kyc-1"))
        );

      // Record many investors
      for (let i = 0; i < 100; i++) {
        await compliance.recordInvestorParticipation(assetId, investor1.address);
      }

      // Should still allow new investors
      expect(await compliance.canInvestInAsset(investor1.address, assetId)).to.be.true;
    });
  });
});
