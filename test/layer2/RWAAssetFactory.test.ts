import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type {
  RWAAssetFactory,
  RWACompliance,
  RWAValuation,
  FractionalRWAToken,
} from "../../typechain-types/index.js";

describe("RWAAssetFactory", function () {
  let factory: RWAAssetFactory;
  let compliance: RWACompliance;
  let valuation: RWAValuation;

  let admin: SignerWithAddress;
  let issuer1: SignerWithAddress;
  let issuer2: SignerWithAddress;
  let complianceOfficer: SignerWithAddress;
  let assetManager: SignerWithAddress;

  const ASSET_MANAGER_ROLE = ethers.keccak256(
    ethers.toUtf8Bytes("ASSET_MANAGER_ROLE")
  );
  const COMPLIANCE_ROLE = ethers.keccak256(ethers.toUtf8Bytes("COMPLIANCE_ROLE"));
  const VERIFIER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("VERIFIER_ROLE"));
  const VALUATOR_ROLE = ethers.keccak256(ethers.toUtf8Bytes("VALUATOR_ROLE"));

  beforeEach(async function () {
    [admin, issuer1, issuer2, complianceOfficer, assetManager] =
      await ethers.getSigners();

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

    // Grant roles
    await compliance.grantRole(VERIFIER_ROLE, complianceOfficer.address);
    await factory.grantRole(COMPLIANCE_ROLE, complianceOfficer.address);
    await factory.grantRole(ASSET_MANAGER_ROLE, assetManager.address);

    // Verify issuers in compliance system
    await compliance
      .connect(complianceOfficer)
      .verifyInvestor(
        issuer1.address,
        1, // AccreditationTier.Accredited
        "US",
        365 * 24 * 60 * 60, // 1 year
        ethers.keccak256(ethers.toUtf8Bytes("issuer1-kyc"))
      );

    await compliance
      .connect(complianceOfficer)
      .verifyInvestor(
        issuer2.address,
        1, // AccreditationTier.Accredited
        "UK",
        365 * 24 * 60 * 60,
        ethers.keccak256(ethers.toUtf8Bytes("issuer2-kyc"))
      );
  });

  describe("Deployment", function () {
    it("Should set the correct admin", async function () {
      expect(await factory.hasRole(await factory.DEFAULT_ADMIN_ROLE(), admin.address)).to.be.true;
    });

    it("Should set the correct compliance contract", async function () {
      expect(await factory.complianceContract()).to.equal(
        await compliance.getAddress()
      );
    });

    it("Should set the correct valuation contract", async function () {
      expect(await factory.valuationContract()).to.equal(
        await valuation.getAddress()
      );
    });

    it("Should initialize counters to zero", async function () {
      expect(await factory.totalAssets()).to.equal(0);
      expect(await factory.activeAssets()).to.equal(0);
      expect(await factory.totalValueLocked()).to.equal(0);
    });

    it("Should revert with invalid admin address", async function () {
      const FactoryContract = await ethers.getContractFactory("RWAAssetFactory");
      await expect(
        FactoryContract.deploy(
          ethers.ZeroAddress,
          await compliance.getAddress(),
          await valuation.getAddress()
        )
      ).to.be.revertedWith("Invalid admin");
    });

    it("Should revert with invalid compliance address", async function () {
      const FactoryContract = await ethers.getContractFactory("RWAAssetFactory");
      await expect(
        FactoryContract.deploy(
          admin.address,
          ethers.ZeroAddress,
          await valuation.getAddress()
        )
      ).to.be.revertedWith("Invalid compliance");
    });

    it("Should revert with invalid valuation address", async function () {
      const FactoryContract = await ethers.getContractFactory("RWAAssetFactory");
      await expect(
        FactoryContract.deploy(
          admin.address,
          await compliance.getAddress(),
          ethers.ZeroAddress
        )
      ).to.be.revertedWith("Invalid valuation");
    });
  });

  describe("Asset Creation", function () {
    const assetName = "Luxury Apartment NYC";
    const assetSymbol = "LAPT-NYC";
    const totalValue = ethers.parseEther("1000000"); // $1M
    const maturityDate = 0; // Perpetual
    const legalDocs = "QmXYZ123...";
    const valuationReport = "QmABC456...";

    it("Should create a new asset successfully", async function () {
      const tx = await factory
        .connect(issuer1)
        .createAsset(
          assetName,
          assetSymbol,
          0, // AssetType.RealEstate
          totalValue,
          maturityDate,
          legalDocs,
          valuationReport
        );

      await expect(tx)
        .to.emit(factory, "AssetCreated")
        .withArgs(1, issuer1.address, 0, totalValue);

      expect(await factory.totalAssets()).to.equal(1);
    });

    it("Should assign correct asset ID sequentially", async function () {
      await factory
        .connect(issuer1)
        .createAsset(
          assetName,
          assetSymbol,
          0,
          totalValue,
          maturityDate,
          legalDocs,
          valuationReport
        );

      await factory
        .connect(issuer2)
        .createAsset(
          "Corporate Bond",
          "CB-001",
          1, // AssetType.Bonds
          ethers.parseEther("500000"),
          maturityDate,
          legalDocs,
          valuationReport
        );

      const metadata1 = await factory.getAssetMetadata(1);
      const metadata2 = await factory.getAssetMetadata(2);

      expect(metadata1.issuer).to.equal(issuer1.address);
      expect(metadata2.issuer).to.equal(issuer2.address);
    });

    it("Should store correct metadata", async function () {
      await factory
        .connect(issuer1)
        .createAsset(
          assetName,
          assetSymbol,
          0,
          totalValue,
          maturityDate,
          legalDocs,
          valuationReport
        );

      const metadata = await factory.getAssetMetadata(1);

      expect(metadata.name).to.equal(assetName);
      expect(metadata.symbol).to.equal(assetSymbol);
      expect(metadata.assetType).to.equal(0); // RealEstate
      expect(metadata.totalValue).to.equal(totalValue);
      expect(metadata.issuer).to.equal(issuer1.address);
      expect(metadata.maturityDate).to.equal(maturityDate);
      expect(metadata.legalDocumentHash).to.equal(legalDocs);
      expect(metadata.valuationReportHash).to.equal(valuationReport);
    });

    it("Should set asset status to Pending", async function () {
      await factory
        .connect(issuer1)
        .createAsset(
          assetName,
          assetSymbol,
          0,
          totalValue,
          maturityDate,
          legalDocs,
          valuationReport
        );

      expect(await factory.getAssetStatus(1)).to.equal(0); // AssetStatus.Pending
    });

    it("Should track issuer's assets", async function () {
      await factory
        .connect(issuer1)
        .createAsset(
          assetName,
          assetSymbol,
          0,
          totalValue,
          maturityDate,
          legalDocs,
          valuationReport
        );

      await factory
        .connect(issuer1)
        .createAsset(
          "Second Property",
          "PROP-2",
          0,
          ethers.parseEther("500000"),
          maturityDate,
          legalDocs,
          valuationReport
        );

      const issuerAssets = await factory.getIssuerAssets(issuer1.address);
      expect(issuerAssets.length).to.equal(2);
      expect(issuerAssets[0]).to.equal(1);
      expect(issuerAssets[1]).to.equal(2);
    });

    it("Should revert with empty name", async function () {
      await expect(
        factory
          .connect(issuer1)
          .createAsset("", assetSymbol, 0, totalValue, maturityDate, legalDocs, valuationReport)
      ).to.be.revertedWith("Invalid name");
    });

    it("Should revert with empty symbol", async function () {
      await expect(
        factory
          .connect(issuer1)
          .createAsset(assetName, "", 0, totalValue, maturityDate, legalDocs, valuationReport)
      ).to.be.revertedWith("Invalid symbol");
    });

    it("Should revert with zero valuation", async function () {
      await expect(
        factory
          .connect(issuer1)
          .createAsset(assetName, assetSymbol, 0, 0, maturityDate, legalDocs, valuationReport)
      ).to.be.revertedWith("Invalid valuation");
    });

    it("Should revert without legal documents", async function () {
      await expect(
        factory
          .connect(issuer1)
          .createAsset(assetName, assetSymbol, 0, totalValue, maturityDate, "", valuationReport)
      ).to.be.revertedWith("Legal docs required");
    });

    it("Should revert without valuation report", async function () {
      await expect(
        factory
          .connect(issuer1)
          .createAsset(assetName, assetSymbol, 0, totalValue, maturityDate, legalDocs, "")
      ).to.be.revertedWith("Valuation report required");
    });

    it("Should revert if issuer is not KYC verified", async function () {
      const [, , , , , unverifiedIssuer] = await ethers.getSigners();

      await expect(
        factory
          .connect(unverifiedIssuer)
          .createAsset(
            assetName,
            assetSymbol,
            0,
            totalValue,
            maturityDate,
            legalDocs,
            valuationReport
          )
      ).to.be.revertedWith("Issuer not verified");
    });

    it("Should support all asset types", async function () {
      // RealEstate = 0, Bonds = 1, Equity = 2, Commodities = 3, ArtCollectible = 4, Invoice = 5

      for (let assetType = 0; assetType <= 5; assetType++) {
        await factory
          .connect(issuer1)
          .createAsset(
            `Asset Type ${assetType}`,
            `AT-${assetType}`,
            assetType,
            totalValue,
            maturityDate,
            legalDocs,
            valuationReport
          );

        const metadata = await factory.getAssetMetadata(assetType + 1);
        expect(metadata.assetType).to.equal(assetType);
      }
    });

    it("Should work when contract is not paused", async function () {
      await factory
        .connect(issuer1)
        .createAsset(
          assetName,
          assetSymbol,
          0,
          totalValue,
          maturityDate,
          legalDocs,
          valuationReport
        );

      expect(await factory.totalAssets()).to.equal(1);
    });

    it("Should revert when contract is paused", async function () {
      await factory.pause();

      await expect(
        factory
          .connect(issuer1)
          .createAsset(
            assetName,
            assetSymbol,
            0,
            totalValue,
            maturityDate,
            legalDocs,
            valuationReport
          )
      ).to.be.revertedWithCustomError(factory, "EnforcedPause");
    });
  });

  describe("Asset Fractionalization", function () {
    const totalSupply = ethers.parseEther("1000000"); // 1M tokens

    beforeEach(async function () {
      // Create an asset first
      await factory
        .connect(issuer1)
        .createAsset(
          "Test Asset",
          "TASSET",
          0,
          ethers.parseEther("1000000"),
          0,
          "QmXYZ",
          "QmABC"
        );
    });

    it("Should fractionalize asset successfully", async function () {
      const tx = await factory.connect(issuer1).fractionalizeAsset(1, totalSupply);

      await expect(tx).to.emit(factory, "AssetFractionalized");

      const tokenAddress = await factory.getFractionalToken(1);
      expect(tokenAddress).to.not.equal(ethers.ZeroAddress);
    });

    it("Should deploy FractionalRWAToken with correct parameters", async function () {
      await factory.connect(issuer1).fractionalizeAsset(1, totalSupply);

      const tokenAddress = await factory.getFractionalToken(1);
      const token = await ethers.getContractAt("FractionalRWAToken", tokenAddress);

      expect(await token.name()).to.equal("Test Asset");
      expect(await token.symbol()).to.equal("TASSET");
      expect(await token.assetId()).to.equal(1);
      expect(await token.supplyCap()).to.equal(totalSupply);
    });

    it("Should update asset metadata with total supply", async function () {
      await factory.connect(issuer1).fractionalizeAsset(1, totalSupply);

      const metadata = await factory.getAssetMetadata(1);
      expect(metadata.totalSupply).to.equal(totalSupply);
    });

    it("Should revert if not called by issuer", async function () {
      await expect(
        factory.connect(issuer2).fractionalizeAsset(1, totalSupply)
      ).to.be.revertedWith("Not asset issuer");
    });

    it("Should revert if asset is not pending", async function () {
      await factory.connect(issuer1).fractionalizeAsset(1, totalSupply);
      await factory.connect(complianceOfficer).activateAsset(1);

      await expect(
        factory.connect(issuer1).fractionalizeAsset(1, totalSupply)
      ).to.be.revertedWith("Asset not pending");
    });

    it("Should revert if already fractionalized", async function () {
      await factory.connect(issuer1).fractionalizeAsset(1, totalSupply);

      await expect(
        factory.connect(issuer1).fractionalizeAsset(1, totalSupply)
      ).to.be.revertedWith("Already fractionalized");
    });

    it("Should revert with zero supply", async function () {
      await expect(
        factory.connect(issuer1).fractionalizeAsset(1, 0)
      ).to.be.revertedWith("Invalid supply");
    });
  });

  describe("Asset Activation", function () {
    beforeEach(async function () {
      await factory
        .connect(issuer1)
        .createAsset(
          "Test Asset",
          "TASSET",
          0,
          ethers.parseEther("1000000"),
          0,
          "QmXYZ",
          "QmABC"
        );
      await factory
        .connect(issuer1)
        .fractionalizeAsset(1, ethers.parseEther("1000000"));
    });

    it("Should activate asset successfully", async function () {
      const tx = await factory.connect(complianceOfficer).activateAsset(1);

      await expect(tx)
        .to.emit(factory, "AssetStatusChanged")
        .withArgs(1, 0, 1); // Pending -> Active

      expect(await factory.getAssetStatus(1)).to.equal(1); // Active
    });

    it("Should update active assets counter", async function () {
      await factory.connect(complianceOfficer).activateAsset(1);

      expect(await factory.activeAssets()).to.equal(1);
    });

    it("Should update total value locked", async function () {
      await factory.connect(complianceOfficer).activateAsset(1);

      expect(await factory.totalValueLocked()).to.equal(ethers.parseEther("1000000"));
    });

    it("Should revert if not called by compliance role", async function () {
      await expect(
        factory.connect(issuer1).activateAsset(1)
      ).to.be.revertedWithCustomError(factory, "AccessControlUnauthorizedAccount");
    });

    it("Should revert if asset is not pending", async function () {
      await factory.connect(complianceOfficer).activateAsset(1);

      await expect(
        factory.connect(complianceOfficer).activateAsset(1)
      ).to.be.revertedWith("Not pending");
    });

    it("Should revert if not fractionalized", async function () {
      await factory
        .connect(issuer1)
        .createAsset(
          "Another Asset",
          "ANOTHER",
          0,
          ethers.parseEther("500000"),
          0,
          "QmDEF",
          "QmGHI"
        );

      await expect(
        factory.connect(complianceOfficer).activateAsset(2)
      ).to.be.revertedWith("Not fractionalized");
    });
  });

  describe("Asset Status Management", function () {
    beforeEach(async function () {
      await factory
        .connect(issuer1)
        .createAsset(
          "Test Asset",
          "TASSET",
          0,
          ethers.parseEther("1000000"),
          0,
          "QmXYZ",
          "QmABC"
        );
      await factory
        .connect(issuer1)
        .fractionalizeAsset(1, ethers.parseEther("1000000"));
      await factory.connect(complianceOfficer).activateAsset(1);
    });

    it("Should update asset status", async function () {
      const tx = await factory
        .connect(assetManager)
        .updateAssetStatus(1, 2); // Suspended

      await expect(tx)
        .to.emit(factory, "AssetStatusChanged")
        .withArgs(1, 1, 2); // Active -> Suspended

      expect(await factory.getAssetStatus(1)).to.equal(2);
    });

    it("Should update TVL when status changes from Active", async function () {
      const tvlBefore = await factory.totalValueLocked();

      await factory
        .connect(assetManager)
        .updateAssetStatus(1, 2); // Active -> Suspended

      expect(await factory.totalValueLocked()).to.equal(
        tvlBefore - ethers.parseEther("1000000")
      );
      expect(await factory.activeAssets()).to.equal(0);
    });

    it("Should update TVL when status changes to Active", async function () {
      await factory
        .connect(assetManager)
        .updateAssetStatus(1, 2); // Active -> Suspended

      await factory
        .connect(assetManager)
        .updateAssetStatus(1, 1); // Suspended -> Active

      expect(await factory.totalValueLocked()).to.equal(ethers.parseEther("1000000"));
      expect(await factory.activeAssets()).to.equal(1);
    });

    it("Should revert if status unchanged", async function () {
      await expect(
        factory.connect(assetManager).updateAssetStatus(1, 1) // Active -> Active
      ).to.be.revertedWith("Status unchanged");
    });

    it("Should revert if not called by asset manager", async function () {
      await expect(
        factory.connect(issuer1).updateAssetStatus(1, 2)
      ).to.be.revertedWithCustomError(factory, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Asset Valuation Updates", function () {
    beforeEach(async function () {
      await factory
        .connect(issuer1)
        .createAsset(
          "Test Asset",
          "TASSET",
          0,
          ethers.parseEther("1000000"),
          0,
          "QmXYZ",
          "QmABC"
        );
      await factory
        .connect(issuer1)
        .fractionalizeAsset(1, ethers.parseEther("1000000"));
      await factory.connect(complianceOfficer).activateAsset(1);
    });

    it("Should update asset valuation", async function () {
      const newValue = ethers.parseEther("1200000");
      const tx = await factory
        .connect(assetManager)
        .updateAssetValuation(1, newValue, "QmNewReport");

      await expect(tx)
        .to.emit(factory, "AssetValuationUpdated")
        .withArgs(1, ethers.parseEther("1000000"), newValue, await ethers.provider.getBlock("latest").then(b => b!.timestamp));

      expect(await factory.getAssetValue(1)).to.equal(newValue);
    });

    it("Should update TVL for active assets", async function () {
      const newValue = ethers.parseEther("1200000");
      await factory
        .connect(assetManager)
        .updateAssetValuation(1, newValue, "QmNewReport");

      expect(await factory.totalValueLocked()).to.equal(newValue);
    });

    it("Should not affect TVL for inactive assets", async function () {
      await factory
        .connect(assetManager)
        .updateAssetStatus(1, 2); // Suspend

      const newValue = ethers.parseEther("1200000");
      await factory
        .connect(assetManager)
        .updateAssetValuation(1, newValue, "QmNewReport");

      expect(await factory.totalValueLocked()).to.equal(0);
    });

    it("Should update valuation report hash", async function () {
      const newValue = ethers.parseEther("1200000");
      await factory
        .connect(assetManager)
        .updateAssetValuation(1, newValue, "QmNewReport");

      const metadata = await factory.getAssetMetadata(1);
      expect(metadata.valuationReportHash).to.equal("QmNewReport");
    });

    it("Should revert with zero valuation", async function () {
      await expect(
        factory.connect(assetManager).updateAssetValuation(1, 0, "QmNewReport")
      ).to.be.revertedWith("Invalid valuation");
    });

    it("Should revert if not called by asset manager", async function () {
      await expect(
        factory
          .connect(issuer1)
          .updateAssetValuation(1, ethers.parseEther("1200000"), "QmNewReport")
      ).to.be.revertedWithCustomError(factory, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Yield Terms", function () {
    beforeEach(async function () {
      await factory
        .connect(issuer1)
        .createAsset(
          "Bond Asset",
          "BOND",
          1, // Bonds
          ethers.parseEther("1000000"),
          0,
          "QmXYZ",
          "QmABC"
        );
      await factory
        .connect(issuer1)
        .fractionalizeAsset(1, ethers.parseEther("1000000"));
      await factory.connect(complianceOfficer).activateAsset(1);
    });

    it("Should set yield terms", async function () {
      const annualYield = 500; // 5%
      const frequency = 30 * 24 * 60 * 60; // 30 days

      await factory.connect(issuer1).setYieldTerms(1, annualYield, frequency);

      const terms = await factory.getYieldTerms(1);
      expect(terms.annualYieldRate).to.equal(annualYield);
      expect(terms.yieldPaymentFrequency).to.equal(frequency);
      expect(terms.totalYieldPaid).to.equal(0);
    });

    it("Should revert if not called by issuer", async function () {
      await expect(
        factory.connect(issuer2).setYieldTerms(1, 500, 30 * 24 * 60 * 60)
      ).to.be.revertedWith("Not asset issuer");
    });

    it("Should revert with invalid yield rate", async function () {
      await expect(
        factory.connect(issuer1).setYieldTerms(1, 10001, 30 * 24 * 60 * 60) // > 100%
      ).to.be.revertedWith("Invalid yield rate");
    });

    it("Should revert with zero frequency", async function () {
      await expect(
        factory.connect(issuer1).setYieldTerms(1, 500, 0)
      ).to.be.revertedWith("Invalid frequency");
    });
  });

  describe("Yield Distribution", function () {
    beforeEach(async function () {
      await factory
        .connect(issuer1)
        .createAsset(
          "Bond Asset",
          "BOND",
          1,
          ethers.parseEther("1000000"),
          0,
          "QmXYZ",
          "QmABC"
        );
      await factory
        .connect(issuer1)
        .fractionalizeAsset(1, ethers.parseEther("1000000"));
      await factory.connect(complianceOfficer).activateAsset(1);
      await factory.connect(issuer1).setYieldTerms(1, 500, 30 * 24 * 60 * 60);
    });

    it("Should distribute yield", async function () {
      const yieldAmount = ethers.parseEther("50000");
      const tx = await factory.connect(issuer1).distributeYield(1, yieldAmount);

      await expect(tx).to.emit(factory, "YieldDistributed").withArgs(1, yieldAmount, await ethers.provider.getBlock("latest").then(b => b!.timestamp));

      const terms = await factory.getYieldTerms(1);
      expect(terms.totalYieldPaid).to.equal(yieldAmount);
    });

    it("Should track cumulative yield", async function () {
      await factory.connect(issuer1).distributeYield(1, ethers.parseEther("50000"));
      await factory.connect(issuer1).distributeYield(1, ethers.parseEther("30000"));

      const terms = await factory.getYieldTerms(1);
      expect(terms.totalYieldPaid).to.equal(ethers.parseEther("80000"));
    });

    it("Should revert if not called by issuer", async function () {
      await expect(
        factory.connect(issuer2).distributeYield(1, ethers.parseEther("50000"))
      ).to.be.revertedWith("Not asset issuer");
    });

    it("Should revert if asset is not active", async function () {
      await factory
        .connect(assetManager)
        .updateAssetStatus(1, 2); // Suspend

      await expect(
        factory.connect(issuer1).distributeYield(1, ethers.parseEther("50000"))
      ).to.be.revertedWith("Asset not active");
    });
  });

  describe("View Functions", function () {
    beforeEach(async function () {
      // Create multiple assets
      for (let i = 0; i < 3; i++) {
        await factory
          .connect(issuer1)
          .createAsset(
            `Asset ${i}`,
            `AST${i}`,
            i % 2, // Alternate between RealEstate and Bonds
            ethers.parseEther("1000000"),
            0,
            "QmXYZ",
            "QmABC"
          );
      }

      await factory
        .connect(issuer1)
        .fractionalizeAsset(1, ethers.parseEther("1000000"));
      await factory.connect(complianceOfficer).activateAsset(1);
    });

    it("Should return all assets", async function () {
      const allAssets = await factory.getAllAssets();
      expect(allAssets.length).to.equal(3);
      expect(allAssets[0]).to.equal(1);
      expect(allAssets[1]).to.equal(2);
      expect(allAssets[2]).to.equal(3);
    });

    it("Should return active assets only", async function () {
      const activeAssets = await factory.getActiveAssets();
      expect(activeAssets.length).to.equal(1);
      expect(activeAssets[0]).to.equal(1);
    });

    it("Should filter assets by type", async function () {
      const realEstateAssets = await factory.getAssetsByType(0); // RealEstate
      const bondAssets = await factory.getAssetsByType(1); // Bonds

      expect(realEstateAssets.length).to.equal(2); // Assets 1 and 3
      expect(bondAssets.length).to.equal(1); // Asset 2
    });

    it("Should check if asset is active", async function () {
      expect(await factory.isAssetActive(1)).to.be.true;
      expect(await factory.isAssetActive(2)).to.be.false;
      expect(await factory.isAssetActive(3)).to.be.false;
    });
  });

  describe("Admin Functions", function () {
    it("Should pause the contract", async function () {
      await factory.pause();
      expect(await factory.paused()).to.be.true;
    });

    it("Should unpause the contract", async function () {
      await factory.pause();
      await factory.unpause();
      expect(await factory.paused()).to.be.false;
    });

    it("Should revert pause if not admin", async function () {
      await expect(
        factory.connect(issuer1).pause()
      ).to.be.revertedWithCustomError(factory, "AccessControlUnauthorizedAccount");
    });

    it("Should revert unpause if not admin", async function () {
      await factory.pause();
      await expect(
        factory.connect(issuer1).unpause()
      ).to.be.revertedWithCustomError(factory, "AccessControlUnauthorizedAccount");
    });
  });
});
