import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type {
  TreasuryPriceOracle,
  TreasuryAssetFactory,
} from "../../typechain-types/index.js";

describe("TreasuryPriceOracle", function () {
  let oracle: TreasuryPriceOracle;
  let factory: TreasuryAssetFactory;

  let admin: SignerWithAddress;
  let priceFeeder: SignerWithAddress;
  let user: SignerWithAddress;

  const OPERATOR_ROLE = ethers.keccak256(ethers.toUtf8Bytes("OPERATOR_ROLE"));
  const PRICE_FEEDER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("PRICE_FEEDER_ROLE"));

  let assetId: number;
  const INITIAL_PRICE = ethers.parseEther("995"); // $995
  const INITIAL_YIELD = 527; // 5.27% in basis points

  beforeEach(async function () {
    [admin, priceFeeder, user] = await ethers.getSigners();

    // Deploy TreasuryAssetFactory
    const FactoryContract = await ethers.getContractFactory("TreasuryAssetFactory");
    factory = await FactoryContract.deploy(admin.address);

    // Deploy TreasuryPriceOracle
    const OracleFactory = await ethers.getContractFactory("TreasuryPriceOracle");
    oracle = await OracleFactory.deploy(
      admin.address,
      await factory.getAddress()
    );

    // Grant roles
    await factory.grantRole(OPERATOR_ROLE, admin.address);
    await oracle.grantRole(PRICE_FEEDER_ROLE, priceFeeder.address);

    // Create Treasury asset
    const issueDate = Math.floor(Date.now() / 1000);
    const maturityDate = issueDate + 90 * 24 * 60 * 60;

    await factory.createTreasuryAsset(
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

  describe("Price Updates", function () {
    it("Should update asset price", async function () {
      const tx = await oracle.connect(priceFeeder).updatePrice(
        assetId,
        INITIAL_PRICE,
        INITIAL_YIELD
      );

      await expect(tx)
        .to.emit(oracle, "PriceUpdated")
        .withArgs(assetId, INITIAL_PRICE, INITIAL_YIELD);

      const price = await oracle.getPrice(assetId);
      expect(price).to.equal(INITIAL_PRICE);
    });

    it("Should update yield rate", async function () {
      await oracle.connect(priceFeeder).updatePrice(assetId, INITIAL_PRICE, INITIAL_YIELD);

      const yieldRate = await oracle.getYield(assetId);
      expect(yieldRate).to.equal(INITIAL_YIELD);
    });

    it("Should track price history", async function () {
      await oracle.connect(priceFeeder).updatePrice(assetId, INITIAL_PRICE, INITIAL_YIELD);

      const newPrice = ethers.parseEther("998");
      const newYield = 520;
      await oracle.connect(priceFeeder).updatePrice(assetId, newPrice, newYield);

      const history = await oracle.getPriceHistory(assetId, 10);
      expect(history.length).to.equal(2);
      expect(history[0].price).to.equal(newPrice); // Most recent first
      expect(history[1].price).to.equal(INITIAL_PRICE);
    });

    it("Should only allow price feeder to update", async function () {
      await expect(
        oracle.connect(user).updatePrice(assetId, INITIAL_PRICE, INITIAL_YIELD)
      ).to.be.reverted;
    });

    it("Should fail with invalid price", async function () {
      await expect(
        oracle.connect(priceFeeder).updatePrice(assetId, 0, INITIAL_YIELD)
      ).to.be.revertedWith("Invalid price");
    });

    it("Should fail with invalid yield", async function () {
      await expect(
        oracle.connect(priceFeeder).updatePrice(assetId, INITIAL_PRICE, 10001) // > 100%
      ).to.be.revertedWith("Invalid yield");
    });
  });

  describe("Price Queries", function () {
    beforeEach(async function () {
      await oracle.connect(priceFeeder).updatePrice(assetId, INITIAL_PRICE, INITIAL_YIELD);
    });

    it("Should get latest price", async function () {
      const price = await oracle.getPrice(assetId);
      expect(price).to.equal(INITIAL_PRICE);
    });

    it("Should get latest yield", async function () {
      const yieldRate = await oracle.getYield(assetId);
      expect(yieldRate).to.equal(INITIAL_YIELD);
    });

    it("Should return price data with timestamp", async function () {
      const data = await oracle.getPriceData(assetId);
      expect(data.price).to.equal(INITIAL_PRICE);
      expect(data.yieldRate).to.equal(INITIAL_YIELD);
      expect(data.timestamp).to.be.gt(0);
    });

    it("Should handle non-existent asset", async function () {
      const nonExistentId = 999;
      await expect(
        oracle.getPrice(nonExistentId)
      ).to.be.revertedWith("No price data");
    });
  });

  describe("Price Staleness", function () {
    it("Should detect stale prices", async function () {
      await oracle.connect(priceFeeder).updatePrice(assetId, INITIAL_PRICE, INITIAL_YIELD);

      // Set staleness threshold to 1 hour
      await oracle.connect(admin).setStalenessThreshold(3600);

      // Initially not stale
      expect(await oracle.isPriceStale(assetId)).to.be.false;

      // Fast forward time
      await ethers.provider.send("evm_increaseTime", [3601]);
      await ethers.provider.send("evm_mine", []);

      // Now stale
      expect(await oracle.isPriceStale(assetId)).to.be.true;
    });

    it("Should allow admin to set staleness threshold", async function () {
      await oracle.connect(admin).setStalenessThreshold(7200);
      expect(await oracle.stalenessThreshold()).to.equal(7200);
    });
  });

  describe("Batch Updates", function () {
    it("Should update multiple assets at once", async function () {
      // Create another asset
      const issueDate = Math.floor(Date.now() / 1000);
      const maturityDate = issueDate + 180 * 24 * 60 * 60;

      await factory.connect(admin).createTreasuryAsset(
        1, // T-NOTE
        "26W",
        "912796AA1",
        issueDate,
        maturityDate,
        ethers.parseEther("1000"),
        450
      );

      const assetId2 = 2;

      const assetIds = [assetId, assetId2];
      const prices = [INITIAL_PRICE, ethers.parseEther("997")];
      const yields = [INITIAL_YIELD, 530];

      await oracle.connect(priceFeeder).batchUpdatePrices(assetIds, prices, yields);

      expect(await oracle.getPrice(assetId)).to.equal(prices[0]);
      expect(await oracle.getPrice(assetId2)).to.equal(prices[1]);
    });

    it("Should fail if array lengths mismatch", async function () {
      const assetIds = [assetId];
      const prices = [INITIAL_PRICE, ethers.parseEther("997")];
      const yields = [INITIAL_YIELD];

      await expect(
        oracle.connect(priceFeeder).batchUpdatePrices(assetIds, prices, yields)
      ).to.be.revertedWith("Array length mismatch");
    });
  });

  describe("Price Deviation", function () {
    beforeEach(async function () {
      await oracle.connect(priceFeeder).updatePrice(assetId, INITIAL_PRICE, INITIAL_YIELD);
    });

    it("Should calculate price change percentage", async function () {
      const newPrice = ethers.parseEther("1005"); // +1% from 995
      await oracle.connect(priceFeeder).updatePrice(assetId, newPrice, INITIAL_YIELD);

      const deviation = await oracle.getPriceChangePercent(assetId);
      // Should be approximately 100 basis points (1%)
      expect(deviation).to.be.closeTo(100, 10);
    });

    it("Should handle price increases", async function () {
      const newPrice = ethers.parseEther("1045"); // +5% from 995
      await oracle.connect(priceFeeder).updatePrice(assetId, newPrice, INITIAL_YIELD);

      const deviation = await oracle.getPriceChangePercent(assetId);
      expect(deviation).to.be.gt(0);
    });

    it("Should handle price decreases", async function () {
      const newPrice = ethers.parseEther("945"); // -5% from 995
      await oracle.connect(priceFeeder).updatePrice(assetId, newPrice, INITIAL_YIELD);

      const deviation = await oracle.getPriceChangePercent(assetId);
      expect(deviation).to.be.lt(0);
    });
  });

  describe("Access Control", function () {
    it("Should allow admin to grant price feeder role", async function () {
      await oracle.connect(admin).grantRole(PRICE_FEEDER_ROLE, user.address);
      expect(await oracle.hasRole(PRICE_FEEDER_ROLE, user.address)).to.be.true;
    });

    it("Should prevent non-admin from granting roles", async function () {
      await expect(
        oracle.connect(user).grantRole(PRICE_FEEDER_ROLE, user.address)
      ).to.be.reverted;
    });
  });
});
