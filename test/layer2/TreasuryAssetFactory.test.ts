import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type {
  TreasuryAssetFactory,
  TreasuryToken,
} from "../../typechain-types/index.js";

describe("TreasuryAssetFactory", function () {
  let factory: TreasuryAssetFactory;

  let admin: SignerWithAddress;
  let operator: SignerWithAddress;
  let user1: SignerWithAddress;
  let user2: SignerWithAddress;

  const OPERATOR_ROLE = ethers.keccak256(ethers.toUtf8Bytes("OPERATOR_ROLE"));
  const MINTER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("MINTER_ROLE"));

  beforeEach(async function () {
    [admin, operator, user1, user2] = await ethers.getSigners();

    // Deploy TreasuryAssetFactory
    const FactoryContract = await ethers.getContractFactory("TreasuryAssetFactory");
    factory = await FactoryContract.deploy(admin.address);

    // Grant operator role
    await factory.grantRole(OPERATOR_ROLE, operator.address);
  });

  describe("Deployment", function () {
    it("Should set the correct admin", async function () {
      const DEFAULT_ADMIN_ROLE = await factory.DEFAULT_ADMIN_ROLE();
      expect(await factory.hasRole(DEFAULT_ADMIN_ROLE, admin.address)).to.be.true;
    });

    it("Should have correct initial state", async function () {
      expect(await factory.assetCounter()).to.equal(0);
    });
  });

  describe("Asset Creation", function () {
    const CUSIP = "912796YZ1";
    const FACE_VALUE = ethers.parseEther("1000"); // $1000
    const COUPON_RATE = 525; // 5.25% (basis points: 525 / 10000)
    const MATURITY_TERM = "13W";

    it("Should create a new Treasury asset", async function () {
      const issueDate = Math.floor(Date.now() / 1000);
      const maturityDate = issueDate + 90 * 24 * 60 * 60; // 90 days

      const tx = await factory.connect(operator).createTreasuryAsset(
        0, // T-BILL
        MATURITY_TERM,
        CUSIP,
        issueDate,
        maturityDate,
        FACE_VALUE,
        COUPON_RATE
      );

      const receipt = await tx.wait();
      const assetId = 1;

      // Check event emission
      await expect(tx)
        .to.emit(factory, "TreasuryAssetCreated")
        .withArgs(assetId, CUSIP, await factory.getAddress());

      // Verify asset details
      const asset = await factory.getTreasuryAsset(assetId);
      expect(asset.cusip).to.equal(CUSIP);
      expect(asset.faceValue).to.equal(FACE_VALUE);
      expect(asset.couponRate).to.equal(COUPON_RATE);
      expect(asset.status).to.equal(0); // Active
    });

    it("Should fail if not called by operator", async function () {
      const issueDate = Math.floor(Date.now() / 1000);
      const maturityDate = issueDate + 90 * 24 * 60 * 60;

      await expect(
        factory.connect(user1).createTreasuryAsset(
          0,
          MATURITY_TERM,
          CUSIP,
          issueDate,
          maturityDate,
          FACE_VALUE,
          COUPON_RATE
        )
      ).to.be.reverted;
    });

    it("Should fail with invalid maturity date", async function () {
      const issueDate = Math.floor(Date.now() / 1000);
      const maturityDate = issueDate - 1; // Invalid: before issue date

      await expect(
        factory.connect(operator).createTreasuryAsset(
          0,
          MATURITY_TERM,
          CUSIP,
          issueDate,
          maturityDate,
          FACE_VALUE,
          COUPON_RATE
        )
      ).to.be.revertedWith("Invalid maturity date");
    });

    it("Should create token contract", async function () {
      const issueDate = Math.floor(Date.now() / 1000);
      const maturityDate = issueDate + 90 * 24 * 60 * 60;

      await factory.connect(operator).createTreasuryAsset(
        0,
        MATURITY_TERM,
        CUSIP,
        issueDate,
        maturityDate,
        FACE_VALUE,
        COUPON_RATE
      );

      const assetId = 1;
      const asset = await factory.getTreasuryAsset(assetId);

      expect(asset.tokenAddress).to.not.equal(ethers.ZeroAddress);
    });
  });

  describe("Asset Management", function () {
    let assetId: number;

    beforeEach(async function () {
      const issueDate = Math.floor(Date.now() / 1000);
      const maturityDate = issueDate + 90 * 24 * 60 * 60;

      await factory.connect(operator).createTreasuryAsset(
        0,
        "13W",
        "912796YZ1",
        issueDate,
        maturityDate,
        ethers.parseEther("1000"),
        525
      );

      assetId = 1;
    });

    it("Should mint tokens to specified address", async function () {
      const mintAmount = ethers.parseEther("100");

      await factory.connect(operator).mintTokens(assetId, user1.address, mintAmount);

      const asset = await factory.getTreasuryAsset(assetId);
      const TreasuryToken = await ethers.getContractFactory("TreasuryToken");
      const token = TreasuryToken.attach(asset.tokenAddress) as any;

      expect(await token.balanceOf(user1.address)).to.equal(mintAmount);
    });

    it("Should burn tokens from address", async function () {
      const mintAmount = ethers.parseEther("100");
      const burnAmount = ethers.parseEther("30");

      await factory.connect(operator).mintTokens(assetId, user1.address, mintAmount);
      await factory.connect(operator).burnTokens(assetId, user1.address, burnAmount);

      const asset = await factory.getTreasuryAsset(assetId);
      const TreasuryToken = await ethers.getContractFactory("TreasuryToken");
      const token = TreasuryToken.attach(asset.tokenAddress) as any;

      expect(await token.balanceOf(user1.address)).to.equal(mintAmount - burnAmount);
    });

    it("Should update asset status", async function () {
      await factory.connect(operator).updateAssetStatus(assetId, 2); // Matured

      const asset = await factory.getTreasuryAsset(assetId);
      expect(asset.status).to.equal(2);
    });

    it("Should get asset by CUSIP", async function () {
      const asset = await factory.getAssetByCUSIP("912796YZ1");
      expect(asset.cusip).to.equal("912796YZ1");
    });

    it("Should list all assets", async function () {
      // Create another asset
      const issueDate = Math.floor(Date.now() / 1000);
      const maturityDate = issueDate + 180 * 24 * 60 * 60;

      await factory.connect(operator).createTreasuryAsset(
        1, // T-NOTE
        "26W",
        "912796AA1",
        issueDate,
        maturityDate,
        ethers.parseEther("1000"),
        450
      );

      const assets = await factory.getAllAssets();
      expect(assets.length).to.equal(2);
    });
  });

  describe("Access Control", function () {
    it("Should allow admin to grant roles", async function () {
      await factory.connect(admin).grantRole(OPERATOR_ROLE, user1.address);
      expect(await factory.hasRole(OPERATOR_ROLE, user1.address)).to.be.true;
    });

    it("Should prevent non-admin from granting roles", async function () {
      await expect(
        factory.connect(user1).grantRole(OPERATOR_ROLE, user2.address)
      ).to.be.reverted;
    });
  });

  describe("Edge Cases", function () {
    it("Should handle zero face value", async function () {
      const issueDate = Math.floor(Date.now() / 1000);
      const maturityDate = issueDate + 90 * 24 * 60 * 60;

      await expect(
        factory.connect(operator).createTreasuryAsset(
          0,
          "13W",
          "ZERO123",
          issueDate,
          maturityDate,
          0, // Zero face value
          525
        )
      ).to.be.revertedWith("Invalid face value");
    });

    it("Should handle very high coupon rate", async function () {
      const issueDate = Math.floor(Date.now() / 1000);
      const maturityDate = issueDate + 90 * 24 * 60 * 60;

      await expect(
        factory.connect(operator).createTreasuryAsset(
          0,
          "13W",
          "HIGH123",
          issueDate,
          maturityDate,
          ethers.parseEther("1000"),
          10001 // > 100%
        )
      ).to.be.revertedWith("Invalid coupon rate");
    });
  });
});
