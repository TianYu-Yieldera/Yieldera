import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type { RWAValuation } from "../../typechain-types/index.js";

describe("RWAValuation", function () {
  let valuation: RWAValuation;

  let admin: SignerWithAddress;
  let valuator1: SignerWithAddress;
  let valuator2: SignerWithAddress;
  let oracleAdmin: SignerWithAddress;
  let user: SignerWithAddress;

  const VALUATOR_ROLE = ethers.keccak256(ethers.toUtf8Bytes("VALUATOR_ROLE"));
  const ORACLE_ADMIN_ROLE = ethers.keccak256(ethers.toUtf8Bytes("ORACLE_ADMIN_ROLE"));

  const assetId = 1;

  beforeEach(async function () {
    [admin, valuator1, valuator2, oracleAdmin, user] = await ethers.getSigners();

    // Deploy RWAValuation
    const ValuationFactory = await ethers.getContractFactory("RWAValuation");
    valuation = await ValuationFactory.deploy(admin.address);

    // Grant roles
    await valuation.grantRole(VALUATOR_ROLE, valuator1.address);
    await valuation.grantRole(ORACLE_ADMIN_ROLE, oracleAdmin.address);
  });

  describe("Deployment", function () {
    it("Should set the correct admin", async function () {
      expect(
        await valuation.hasRole(await valuation.DEFAULT_ADMIN_ROLE(), admin.address)
      ).to.be.true;
    });

    it("Should grant initial roles", async function () {
      expect(await valuation.hasRole(VALUATOR_ROLE, admin.address)).to.be.true;
      expect(await valuation.hasRole(ORACLE_ADMIN_ROLE, admin.address)).to.be.true;
    });

    it("Should mark admin as authorized valuator", async function () {
      expect(await valuation.isValuatorAuthorized(admin.address)).to.be.true;
    });

    it("Should initialize counters to zero", async function () {
      expect(await valuation.totalValuations()).to.equal(0);
      expect(await valuation.totalAssets()).to.equal(0);
    });

    it("Should revert with invalid admin", async function () {
      const ValuationFactory = await ethers.getContractFactory("RWAValuation");
      await expect(
        ValuationFactory.deploy(ethers.ZeroAddress)
      ).to.be.revertedWith("Invalid admin address");
    });
  });

  describe("Valuation Updates", function () {
    const value = ethers.parseEther("1000000"); // $1M
    const confidence = 9000; // 90%

    it("Should update valuation successfully", async function () {
      const tx = await valuation
        .connect(valuator1)
        .updateValuation(
          assetId,
          value,
          0, // ValuationMethod.Manual
          "QmValuationReport123",
          confidence
        );

      await expect(tx)
        .to.emit(valuation, "ValuationUpdated")
        .withArgs(assetId, 0, value, 0, valuator1.address);

      expect(await valuation.totalValuations()).to.equal(1);
      expect(await valuation.totalAssets()).to.equal(1);
    });

    it("Should store valuation data correctly", async function () {
      await valuation
        .connect(valuator1)
        .updateValuation(assetId, value, 0, "QmValuationReport123", confidence);

      const lastValuation = await valuation.getLastValuation(assetId);

      expect(lastValuation.value).to.equal(value);
      expect(lastValuation.method).to.equal(0); // Manual
      expect(lastValuation.valuator).to.equal(valuator1.address);
      expect(lastValuation.reportHash).to.equal("QmValuationReport123");
      expect(lastValuation.confidence).to.equal(confidence);
    });

    it("Should get current value", async function () {
      await valuation
        .connect(valuator1)
        .updateValuation(assetId, value, 0, "QmReport1", confidence);

      expect(await valuation.getCurrentValue(assetId)).to.equal(value);
    });

    it("Should support all valuation methods", async function () {
      // Manual = 0, Oracle = 1, Formula = 2, Hybrid = 3
      for (let method = 0; method <= 3; method++) {
        await valuation
          .connect(valuator1)
          .updateValuation(
            method + 1, // Different asset ID for each
            value,
            method,
            `QmReport${method}`,
            confidence
          );

        const lastVal = await valuation.getLastValuation(method + 1);
        expect(lastVal.method).to.equal(method);
      }
    });

    it("Should track multiple valuations for same asset", async function () {
      await valuation
        .connect(valuator1)
        .updateValuation(assetId, value, 0, "QmReport1", confidence);

      await valuation
        .connect(valuator1)
        .updateValuation(assetId, value * 2n, 0, "QmReport2", confidence);

      expect(await valuation.getValuationCount(assetId)).to.equal(2);
    });

    it("Should revert with zero valuation", async function () {
      await expect(
        valuation
          .connect(valuator1)
          .updateValuation(assetId, 0, 0, "QmReport", confidence)
      ).to.be.revertedWith("Invalid valuation");
    });

    it("Should revert with low confidence", async function () {
      await expect(
        valuation
          .connect(valuator1)
          .updateValuation(assetId, value, 0, "QmReport", 4999) // < 5000
      ).to.be.revertedWith("Invalid confidence");
    });

    it("Should revert with high confidence", async function () {
      await expect(
        valuation
          .connect(valuator1)
          .updateValuation(assetId, value, 0, "QmReport", 10001) // > 10000
      ).to.be.revertedWith("Invalid confidence");
    });

    it("Should revert without report hash", async function () {
      await expect(
        valuation.connect(valuator1).updateValuation(assetId, value, 0, "", confidence)
      ).to.be.revertedWith("Report hash required");
    });

    it("Should revert if not valuator role", async function () {
      await expect(
        valuation.connect(user).updateValuation(assetId, value, 0, "QmReport", confidence)
      ).to.be.revertedWithCustomError(valuation, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Valuator Management", function () {
    it("Should authorize new valuator", async function () {
      const tx = await valuation.authorizeValuator(valuator2.address, true);

      await expect(tx)
        .to.emit(valuation, "ValuatorAuthorized")
        .withArgs(valuator2.address, true);

      expect(await valuation.isValuatorAuthorized(valuator2.address)).to.be.true;
      expect(await valuation.hasRole(VALUATOR_ROLE, valuator2.address)).to.be.true;
    });

    it("Should revoke valuator", async function () {
      await valuation.authorizeValuator(valuator2.address, true);
      await valuation.authorizeValuator(valuator2.address, false);

      expect(await valuation.isValuatorAuthorized(valuator2.address)).to.be.false;
      expect(await valuation.hasRole(VALUATOR_ROLE, valuator2.address)).to.be.false;
    });

    it("Should revert with invalid valuator address", async function () {
      await expect(
        valuation.authorizeValuator(ethers.ZeroAddress, true)
      ).to.be.revertedWith("Invalid valuator address");
    });

    it("Should revert if not admin", async function () {
      await expect(
        valuation.connect(user).authorizeValuator(valuator2.address, true)
      ).to.be.revertedWithCustomError(valuation, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Valuation History", function () {
    beforeEach(async function () {
      // Create 5 valuations
      for (let i = 0; i < 5; i++) {
        await valuation
          .connect(valuator1)
          .updateValuation(
            assetId,
            ethers.parseEther(`${1000000 + i * 10000}`),
            0,
            `QmReport${i}`,
            9000
          );
      }
    });

    it("Should get valuation history", async function () {
      const history = await valuation.getValuationHistory(assetId, 3);

      expect(history.length).to.equal(3);
      // Should return most recent 3
      expect(history[2].value).to.equal(ethers.parseEther("1040000"));
    });

    it("Should return all valuations if count exceeds available", async function () {
      const history = await valuation.getValuationHistory(assetId, 100);

      expect(history.length).to.equal(5);
    });

    it("Should get valuation count", async function () {
      expect(await valuation.getValuationCount(assetId)).to.equal(5);
    });

    it("Should get most recent valuation", async function () {
      const lastVal = await valuation.getLastValuation(assetId);

      expect(lastVal.value).to.equal(ethers.parseEther("1040000"));
    });

    it("Should revert getting last valuation for unvalued asset", async function () {
      await expect(valuation.getLastValuation(999)).to.be.revertedWith("No valuations");
    });
  });

  describe("Price Calculations", function () {
    beforeEach(async function () {
      await valuation
        .connect(valuator1)
        .updateValuation(
          assetId,
          ethers.parseEther("1000000"),
          0,
          "QmReport",
          9000
        );
    });

    it("Should calculate price per token", async function () {
      const pricePerToken = await valuation.getPricePerToken(assetId);

      // $1M (in wei: 1000000 * 10^18) / 1e18 tokens = 1000000 (raw value)
      expect(pricePerToken).to.equal(1000000);
    });

    it("Should return zero for unvalued asset", async function () {
      expect(await valuation.getCurrentValue(999)).to.equal(0);
      expect(await valuation.getPricePerToken(999)).to.equal(0);
    });
  });

  describe("Stale Valuation Detection", function () {
    it("Should detect fresh valuation", async function () {
      await valuation
        .connect(valuator1)
        .updateValuation(
          assetId,
          ethers.parseEther("1000000"),
          0,
          "QmReport",
          9000
        );

      expect(await valuation.isValuationStale(assetId)).to.be.false;
    });

    it("Should detect stale valuation after 90 days", async function () {
      await valuation
        .connect(valuator1)
        .updateValuation(
          assetId,
          ethers.parseEther("1000000"),
          0,
          "QmReport",
          9000
        );

      // Fast forward 91 days
      await ethers.provider.send("evm_increaseTime", [91 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      expect(await valuation.isValuationStale(assetId)).to.be.true;
    });

    it("Should consider unvalued asset as stale", async function () {
      expect(await valuation.isValuationStale(999)).to.be.true;
    });
  });

  describe("Revaluation Requests", function () {
    beforeEach(async function () {
      await valuation
        .connect(valuator1)
        .updateValuation(
          assetId,
          ethers.parseEther("1000000"),
          0,
          "QmReport",
          9000
        );
    });

    it("Should request revaluation for stale asset", async function () {
      // Fast forward to make it stale
      await ethers.provider.send("evm_increaseTime", [91 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      const tx = await valuation.connect(user).requestRevaluation(assetId);

      await expect(tx).to.emit(valuation, "ValuationDisputed");
    });

    it("Should revert requesting revaluation for fresh asset", async function () {
      await expect(
        valuation.connect(user).requestRevaluation(assetId)
      ).to.be.revertedWith("Valuation not stale");
    });

    it("Should revert for unvalued asset", async function () {
      await expect(
        valuation.connect(user).requestRevaluation(999)
      ).to.be.revertedWith("Asset not valued");
    });
  });

  describe("Oracle Configuration", function () {
    const mockOracle = "0x1111111111111111111111111111111111111111";
    const heartbeat = 3600; // 1 hour
    const deviationThreshold = 500; // 5%

    it("Should configure oracle", async function () {
      const tx = await valuation
        .connect(oracleAdmin)
        .configureOracle(assetId, mockOracle, heartbeat, deviationThreshold);

      await expect(tx)
        .to.emit(valuation, "OracleConfigured")
        .withArgs(assetId, mockOracle, heartbeat);
    });

    it("Should store oracle configuration", async function () {
      await valuation
        .connect(oracleAdmin)
        .configureOracle(assetId, mockOracle, heartbeat, deviationThreshold);

      const config = await valuation.getOracleConfig(assetId);

      expect(config.oracleAddress).to.equal(mockOracle);
      expect(config.heartbeat).to.equal(heartbeat);
      expect(config.deviationThreshold).to.equal(deviationThreshold);
      expect(config.isActive).to.be.true;
    });

    it("Should disable oracle", async function () {
      await valuation
        .connect(oracleAdmin)
        .configureOracle(assetId, mockOracle, heartbeat, deviationThreshold);

      await valuation.connect(oracleAdmin).disableOracle(assetId);

      const config = await valuation.getOracleConfig(assetId);
      expect(config.isActive).to.be.false;
    });

    it("Should revert with invalid oracle address", async function () {
      await expect(
        valuation
          .connect(oracleAdmin)
          .configureOracle(assetId, ethers.ZeroAddress, heartbeat, deviationThreshold)
      ).to.be.revertedWith("Invalid oracle address");
    });

    it("Should revert with zero heartbeat", async function () {
      await expect(
        valuation
          .connect(oracleAdmin)
          .configureOracle(assetId, mockOracle, 0, deviationThreshold)
      ).to.be.revertedWith("Invalid heartbeat");
    });

    it("Should revert with heartbeat too long", async function () {
      await expect(
        valuation
          .connect(oracleAdmin)
          .configureOracle(
            assetId,
            mockOracle,
            2 * 24 * 60 * 60, // 2 days
            deviationThreshold
          )
      ).to.be.revertedWith("Invalid heartbeat");
    });

    it("Should revert with invalid deviation threshold", async function () {
      await expect(
        valuation
          .connect(oracleAdmin)
          .configureOracle(assetId, mockOracle, heartbeat, 10001) // > 10000
      ).to.be.revertedWith("Invalid deviation threshold");
    });

    it("Should revert if not oracle admin", async function () {
      await expect(
        valuation
          .connect(user)
          .configureOracle(assetId, mockOracle, heartbeat, deviationThreshold)
      ).to.be.revertedWithCustomError(valuation, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Dispute Management", function () {
    beforeEach(async function () {
      await valuation
        .connect(valuator1)
        .updateValuation(
          assetId,
          ethers.parseEther("1000000"),
          0,
          "QmReport",
          9000
        );
    });

    it("Should create dispute", async function () {
      const tx = await valuation
        .connect(user)
        .disputeValuation(assetId, "Valuation seems too high");

      await expect(tx).to.emit(valuation, "ValuationDisputed");
    });

    it("Should store dispute data", async function () {
      await valuation
        .connect(user)
        .disputeValuation(assetId, "Valuation seems too high");

      const disputes = await valuation.getDisputes(assetId);

      expect(disputes.length).to.equal(1);
      expect(disputes[0].disputer).to.equal(user.address);
      expect(disputes[0].reason).to.equal("Valuation seems too high");
      expect(disputes[0].resolved).to.be.false;
      expect(disputes[0].disputedValue).to.equal(ethers.parseEther("1000000"));
    });

    it("Should allow multiple disputes", async function () {
      await valuation.connect(user).disputeValuation(assetId, "Reason 1");
      await valuation.connect(valuator2).disputeValuation(assetId, "Reason 2");

      const disputes = await valuation.getDisputes(assetId);
      expect(disputes.length).to.equal(2);
    });

    it("Should resolve dispute", async function () {
      await valuation.connect(user).disputeValuation(assetId, "Dispute reason");

      await valuation.resolveDispute(assetId, 0);

      const disputes = await valuation.getDisputes(assetId);
      expect(disputes[0].resolved).to.be.true;
    });

    it("Should revert dispute for unvalued asset", async function () {
      await expect(
        valuation.connect(user).disputeValuation(999, "Reason")
      ).to.be.revertedWith("Asset not valued");
    });

    it("Should revert dispute without reason", async function () {
      await expect(
        valuation.connect(user).disputeValuation(assetId, "")
      ).to.be.revertedWith("Reason required");
    });

    it("Should revert resolving invalid dispute", async function () {
      await expect(valuation.resolveDispute(assetId, 999)).to.be.revertedWith(
        "Invalid dispute index"
      );
    });

    it("Should revert resolving if not admin", async function () {
      await valuation.connect(user).disputeValuation(assetId, "Reason");

      await expect(
        valuation.connect(user).resolveDispute(assetId, 0)
      ).to.be.revertedWithCustomError(valuation, "AccessControlUnauthorizedAccount");
    });
  });

  describe("TWAP Calculation", function () {
    beforeEach(async function () {
      // Create valuations at different times
      await valuation
        .connect(valuator1)
        .updateValuation(
          assetId,
          ethers.parseEther("1000000"),
          0,
          "QmReport1",
          9000
        );

      await ethers.provider.send("evm_increaseTime", [3600]); // 1 hour
      await ethers.provider.send("evm_mine", []);

      await valuation
        .connect(valuator1)
        .updateValuation(
          assetId,
          ethers.parseEther("1100000"),
          0,
          "QmReport2",
          9000
        );

      await ethers.provider.send("evm_increaseTime", [3600]); // 1 hour
      await ethers.provider.send("evm_mine", []);

      await valuation
        .connect(valuator1)
        .updateValuation(
          assetId,
          ethers.parseEther("1200000"),
          0,
          "QmReport3",
          9000
        );
    });

    it("Should calculate TWAP", async function () {
      const twap = await valuation.getTWAP(assetId, 10000); // Last ~3 hours

      // TWAP should be between min and max values
      expect(twap).to.be.gte(ethers.parseEther("1000000"));
      expect(twap).to.be.lte(ethers.parseEther("1200000"));
    });

    it("Should revert TWAP with zero period", async function () {
      await expect(valuation.getTWAP(assetId, 0)).to.be.revertedWith("Invalid period");
    });

    it("Should revert TWAP for unvalued asset", async function () {
      await expect(valuation.getTWAP(999, 3600)).to.be.revertedWith("No valuations");
    });
  });

  describe("Edge Cases", function () {
    it("Should handle asset with no valuations", async function () {
      expect(await valuation.getCurrentValue(999)).to.equal(0);
      expect(await valuation.getValuationCount(999)).to.equal(0);
    });

    it("Should update same asset multiple times", async function () {
      for (let i = 0; i < 10; i++) {
        await valuation
          .connect(valuator1)
          .updateValuation(
            assetId,
            ethers.parseEther(`${1000000 + i * 1000}`),
            0,
            `QmReport${i}`,
            9000
          );
      }

      expect(await valuation.getValuationCount(assetId)).to.equal(10);
      expect(await valuation.getCurrentValue(assetId)).to.equal(
        ethers.parseEther("1009000")
      );
    });

    it("Should track statistics across multiple assets", async function () {
      for (let asset = 1; asset <= 5; asset++) {
        await valuation
          .connect(valuator1)
          .updateValuation(
            asset,
            ethers.parseEther("1000000"),
            0,
            "QmReport",
            9000
          );
      }

      expect(await valuation.totalAssets()).to.equal(5);
      expect(await valuation.totalValuations()).to.equal(5);
    });
  });
});
