import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type {
  CompoundV3Adapter,
  MockComet,
  MockCometRewards,
  MockERC20,
  L2StateAggregator,
} from "../../typechain-types/index.js";

describe("CompoundV3Adapter", function () {
  let adapter: CompoundV3Adapter;
  let comet: MockComet;
  let cometRewards: MockCometRewards;
  let baseToken: MockERC20;
  let collateralToken: MockERC20;
  let rewardToken: MockERC20;
  let stateAggregator: L2StateAggregator;

  let owner: SignerWithAddress;
  let user1: SignerWithAddress;
  let user2: SignerWithAddress;

  const INITIAL_SUPPLY = ethers.parseEther("1000000");
  const SUPPLY_AMOUNT = ethers.parseEther("1000");
  const BORROW_AMOUNT = ethers.parseEther("500");

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();

    // Deploy mock ERC20 tokens
    const MockERC20Factory = await ethers.getContractFactory("MockERC20");
    baseToken = await MockERC20Factory.deploy("USDC", "USDC", INITIAL_SUPPLY);
    collateralToken = await MockERC20Factory.deploy(
      "Wrapped ETH",
      "WETH",
      INITIAL_SUPPLY
    );
    rewardToken = await MockERC20Factory.deploy("COMP", "COMP", INITIAL_SUPPLY);

    // Deploy mock Compound V3 contracts
    const MockCometFactory = await ethers.getContractFactory("MockComet");
    comet = await MockCometFactory.deploy(await baseToken.getAddress());

    const MockCometRewardsFactory = await ethers.getContractFactory(
      "MockCometRewards"
    );
    cometRewards = await MockCometRewardsFactory.deploy(
      await rewardToken.getAddress()
    );

    // Deploy L2StateAggregator
    const AggregatorFactory = await ethers.getContractFactory(
      "L2StateAggregator"
    );
    stateAggregator = await AggregatorFactory.deploy(owner.address);

    // Deploy CompoundV3Adapter
    const AdapterFactory = await ethers.getContractFactory("CompoundV3Adapter");
    adapter = await AdapterFactory.deploy(
      await comet.getAddress(),
      await cometRewards.getAddress(),
      owner.address
    );

    // Set state aggregator
    await adapter.setStateAggregator(await stateAggregator.getAddress());

    // Authorize adapter as a module in the state aggregator
    await stateAggregator.authorizeModule(await adapter.getAddress());

    // Fund comet with liquidity
    await baseToken.transfer(await comet.getAddress(), INITIAL_SUPPLY / 2n);

    // Distribute tokens to users
    await baseToken.transfer(user1.address, ethers.parseEther("10000"));
    await baseToken.transfer(user2.address, ethers.parseEther("10000"));
    await collateralToken.transfer(user1.address, ethers.parseEther("100"));
    await collateralToken.transfer(user2.address, ethers.parseEther("100"));

    // Approve adapter to spend tokens
    await baseToken
      .connect(user1)
      .approve(await adapter.getAddress(), ethers.MaxUint256);
    await baseToken
      .connect(user2)
      .approve(await adapter.getAddress(), ethers.MaxUint256);
    await collateralToken
      .connect(user1)
      .approve(await adapter.getAddress(), ethers.MaxUint256);
    await collateralToken
      .connect(user2)
      .approve(await adapter.getAddress(), ethers.MaxUint256);
  });

  describe("Deployment", function () {
    it("Should set the correct comet address", async function () {
      expect(await adapter.comet()).to.equal(await comet.getAddress());
    });

    it("Should set the correct rewards contract", async function () {
      expect(await adapter.cometRewards()).to.equal(
        await cometRewards.getAddress()
      );
    });

    it("Should set the correct base asset", async function () {
      expect(await adapter.baseAsset()).to.equal(await baseToken.getAddress());
    });

    it("Should set the correct admin roles", async function () {
      const DEFAULT_ADMIN_ROLE = ethers.ZeroHash;
      const MANAGER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("MANAGER_ROLE"));

      expect(await adapter.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be
        .true;
      expect(await adapter.hasRole(MANAGER_ROLE, owner.address)).to.be.true;
    });

    it("Should initialize with zero statistics", async function () {
      expect(await adapter.totalSupplied()).to.equal(0);
      expect(await adapter.totalBorrowed()).to.equal(0);
      expect(await adapter.totalCollateral()).to.equal(0);
      expect(await adapter.activePositions()).to.equal(0);
    });

    it("Should revert with invalid comet address", async function () {
      const AdapterFactory = await ethers.getContractFactory("CompoundV3Adapter");
      await expect(
        AdapterFactory.deploy(
          ethers.ZeroAddress,
          await cometRewards.getAddress(),
          owner.address
        )
      ).to.be.revertedWith("Invalid comet");
    });

    it("Should revert with invalid admin address", async function () {
      const AdapterFactory = await ethers.getContractFactory("CompoundV3Adapter");
      await expect(
        AdapterFactory.deploy(
          await comet.getAddress(),
          await cometRewards.getAddress(),
          ethers.ZeroAddress
        )
      ).to.be.revertedWith("Invalid admin");
    });
  });

  describe("Supply Base Asset", function () {
    it("Should allow users to supply base asset", async function () {
      await expect(adapter.connect(user1).supplyBase(SUPPLY_AMOUNT))
        .to.emit(adapter, "BaseSupplied")
        .withArgs(user1.address, SUPPLY_AMOUNT);

      const position = await adapter.positions(user1.address);
      expect(position.baseSupplied).to.equal(SUPPLY_AMOUNT);
      expect(position.active).to.be.true;
    });

    it("Should transfer tokens correctly", async function () {
      const balanceBefore = await baseToken.balanceOf(user1.address);

      await adapter.connect(user1).supplyBase(SUPPLY_AMOUNT);

      const balanceAfter = await baseToken.balanceOf(user1.address);
      expect(balanceBefore - balanceAfter).to.equal(SUPPLY_AMOUNT);
    });

    it("Should update statistics correctly", async function () {
      await adapter.connect(user1).supplyBase(SUPPLY_AMOUNT);

      expect(await adapter.totalSupplied()).to.equal(SUPPLY_AMOUNT);
      expect(await adapter.activePositions()).to.equal(1);
    });

    it("Should update state aggregator", async function () {
      await adapter.connect(user1).supplyBase(SUPPLY_AMOUNT);

      const state = await stateAggregator.getSystemState();
      expect(state.activePositions).to.equal(1);
      expect(state.totalOrders).to.equal(SUPPLY_AMOUNT); // totalSupplied maps to totalOrders
    });

    it("Should track multiple supplies from same user", async function () {
      await adapter.connect(user1).supplyBase(SUPPLY_AMOUNT);
      await adapter.connect(user1).supplyBase(SUPPLY_AMOUNT);

      const position = await adapter.positions(user1.address);
      expect(position.baseSupplied).to.equal(SUPPLY_AMOUNT * 2n);
      expect(await adapter.activePositions()).to.equal(1);
    });

    it("Should track multiple users", async function () {
      await adapter.connect(user1).supplyBase(SUPPLY_AMOUNT);
      await adapter.connect(user2).supplyBase(SUPPLY_AMOUNT);

      expect(await adapter.activePositions()).to.equal(2);
      expect(await adapter.totalSupplied()).to.equal(SUPPLY_AMOUNT * 2n);
    });

    it("Should update timestamp on supply", async function () {
      const tx = await adapter.connect(user1).supplyBase(SUPPLY_AMOUNT);
      const receipt = await tx.wait();
      const block = await ethers.provider.getBlock(receipt!.blockNumber);

      const position = await adapter.positions(user1.address);
      expect(position.lastUpdate).to.equal(block!.timestamp);
    });

    it("Should revert if amount is zero", async function () {
      await expect(
        adapter.connect(user1).supplyBase(0)
      ).to.be.revertedWith("Amount must be greater than zero");
    });

    it("Should revert if user has insufficient balance", async function () {
      await expect(
        adapter.connect(user1).supplyBase(ethers.parseEther("100000"))
      ).to.be.reverted;
    });

    it("Should revert when contract is paused", async function () {
      await adapter.pause();

      await expect(
        adapter.connect(user1).supplyBase(SUPPLY_AMOUNT)
      ).to.be.revertedWithCustomError(adapter, "EnforcedPause");
    });
  });

  describe("Supply Collateral", function () {
    it("Should allow users to supply collateral", async function () {
      const collateralAmount = ethers.parseEther("10");

      await expect(
        adapter
          .connect(user1)
          .supplyCollateral(await collateralToken.getAddress(), collateralAmount)
      )
        .to.emit(adapter, "CollateralSupplied")
        .withArgs(
          user1.address,
          await collateralToken.getAddress(),
          collateralAmount
        );
    });

    it("Should transfer collateral tokens correctly", async function () {
      const collateralAmount = ethers.parseEther("10");
      const balanceBefore = await collateralToken.balanceOf(user1.address);

      await adapter
        .connect(user1)
        .supplyCollateral(await collateralToken.getAddress(), collateralAmount);

      const balanceAfter = await collateralToken.balanceOf(user1.address);
      expect(balanceBefore - balanceAfter).to.equal(collateralAmount);
    });

    it("Should update position active status", async function () {
      const collateralAmount = ethers.parseEther("10");

      await adapter
        .connect(user1)
        .supplyCollateral(await collateralToken.getAddress(), collateralAmount);

      const position = await adapter.positions(user1.address);
      expect(position.active).to.be.true;
      expect(await adapter.activePositions()).to.equal(1);
    });

    it("Should revert if amount is zero", async function () {
      await expect(
        adapter
          .connect(user1)
          .supplyCollateral(await collateralToken.getAddress(), 0)
      ).to.be.revertedWith("Amount must be greater than zero");
    });

    it("Should revert if supplying base asset as collateral", async function () {
      await expect(
        adapter
          .connect(user1)
          .supplyCollateral(await baseToken.getAddress(), SUPPLY_AMOUNT)
      ).to.be.revertedWith("Use supplyBase for base asset");
    });
  });

  describe("Withdraw Base Asset", function () {
    beforeEach(async function () {
      await adapter.connect(user1).supplyBase(SUPPLY_AMOUNT);
    });

    it("Should allow users to withdraw base asset", async function () {
      const withdrawAmount = SUPPLY_AMOUNT / 2n;

      await expect(adapter.connect(user1).withdrawBase(withdrawAmount))
        .to.emit(adapter, "BaseWithdrawn")
        .withArgs(user1.address, withdrawAmount);

      const position = await adapter.positions(user1.address);
      expect(position.baseSupplied).to.equal(SUPPLY_AMOUNT - withdrawAmount);
    });

    it("Should transfer tokens correctly", async function () {
      const withdrawAmount = SUPPLY_AMOUNT / 2n;
      const balanceBefore = await baseToken.balanceOf(user1.address);

      await adapter.connect(user1).withdrawBase(withdrawAmount);

      const balanceAfter = await baseToken.balanceOf(user1.address);
      expect(balanceAfter - balanceBefore).to.equal(withdrawAmount);
    });

    it("Should update statistics correctly", async function () {
      const withdrawAmount = SUPPLY_AMOUNT / 2n;

      await adapter.connect(user1).withdrawBase(withdrawAmount);

      expect(await adapter.totalSupplied()).to.equal(
        SUPPLY_AMOUNT - withdrawAmount
      );
    });

    it("Should deactivate position when fully withdrawn", async function () {
      await adapter.connect(user1).withdrawBase(SUPPLY_AMOUNT);

      const position = await adapter.positions(user1.address);
      expect(position.baseSupplied).to.equal(0);
      expect(position.active).to.be.false;
      expect(await adapter.activePositions()).to.equal(0);
    });

    it("Should update timestamp on withdrawal", async function () {
      const tx = await adapter
        .connect(user1)
        .withdrawBase(SUPPLY_AMOUNT / 2n);
      const receipt = await tx.wait();
      const block = await ethers.provider.getBlock(receipt!.blockNumber);

      const position = await adapter.positions(user1.address);
      expect(position.lastUpdate).to.equal(block!.timestamp);
    });

    it("Should revert if amount is zero", async function () {
      await expect(
        adapter.connect(user1).withdrawBase(0)
      ).to.be.revertedWith("Amount must be greater than zero");
    });

    it("Should revert if amount exceeds supplied balance", async function () {
      await expect(
        adapter.connect(user1).withdrawBase(SUPPLY_AMOUNT + ethers.parseEther("1"))
      ).to.be.revertedWith("Insufficient supplied balance");
    });

    it("Should revert when contract is paused", async function () {
      await adapter.pause();

      await expect(
        adapter.connect(user1).withdrawBase(SUPPLY_AMOUNT / 2n)
      ).to.be.revertedWithCustomError(adapter, "EnforcedPause");
    });
  });

  describe("Withdraw Collateral", function () {
    const collateralAmount = ethers.parseEther("10");

    beforeEach(async function () {
      await adapter
        .connect(user1)
        .supplyCollateral(await collateralToken.getAddress(), collateralAmount);
    });

    it("Should allow users to withdraw collateral", async function () {
      const withdrawAmount = collateralAmount / 2n;

      await expect(
        adapter
          .connect(user1)
          .withdrawCollateral(await collateralToken.getAddress(), withdrawAmount)
      )
        .to.emit(adapter, "CollateralWithdrawn")
        .withArgs(user1.address, await collateralToken.getAddress(), withdrawAmount);
    });

    it("Should transfer tokens correctly", async function () {
      const withdrawAmount = collateralAmount / 2n;
      const balanceBefore = await collateralToken.balanceOf(user1.address);

      await adapter
        .connect(user1)
        .withdrawCollateral(await collateralToken.getAddress(), withdrawAmount);

      const balanceAfter = await collateralToken.balanceOf(user1.address);
      expect(balanceAfter - balanceBefore).to.equal(withdrawAmount);
    });

    it("Should revert if amount is zero", async function () {
      await expect(
        adapter
          .connect(user1)
          .withdrawCollateral(await collateralToken.getAddress(), 0)
      ).to.be.revertedWith("Amount must be greater than zero");
    });

    it("Should revert if using base asset", async function () {
      await expect(
        adapter
          .connect(user1)
          .withdrawCollateral(await baseToken.getAddress(), collateralAmount)
      ).to.be.revertedWith("Use withdrawBase for base asset");
    });
  });

  describe("Borrow", function () {
    beforeEach(async function () {
      // Supply collateral first
      await adapter
        .connect(user1)
        .supplyCollateral(
          await collateralToken.getAddress(),
          ethers.parseEther("10")
        );
    });

    it("Should allow users to borrow base asset", async function () {
      await expect(adapter.connect(user1).borrow(BORROW_AMOUNT))
        .to.emit(adapter, "Borrowed")
        .withArgs(user1.address, BORROW_AMOUNT);

      const position = await adapter.positions(user1.address);
      expect(position.baseBorrowed).to.equal(BORROW_AMOUNT);
    });

    it("Should transfer tokens correctly", async function () {
      const balanceBefore = await baseToken.balanceOf(user1.address);

      await adapter.connect(user1).borrow(BORROW_AMOUNT);

      const balanceAfter = await baseToken.balanceOf(user1.address);
      expect(balanceAfter - balanceBefore).to.equal(BORROW_AMOUNT);
    });

    it("Should update statistics correctly", async function () {
      await adapter.connect(user1).borrow(BORROW_AMOUNT);

      expect(await adapter.totalBorrowed()).to.equal(BORROW_AMOUNT);
    });

    it("Should update state aggregator", async function () {
      await adapter.connect(user1).borrow(BORROW_AMOUNT);

      const state = await stateAggregator.getSystemState();
      expect(state.totalDebt).to.equal(BORROW_AMOUNT);
    });

    it("Should track multiple borrows from same user", async function () {
      await adapter.connect(user1).borrow(BORROW_AMOUNT);
      await adapter.connect(user1).borrow(BORROW_AMOUNT);

      const position = await adapter.positions(user1.address);
      expect(position.baseBorrowed).to.equal(BORROW_AMOUNT * 2n);
    });

    it("Should revert if amount is zero", async function () {
      await expect(adapter.connect(user1).borrow(0)).to.be.revertedWith(
        "Amount must be greater than zero"
      );
    });

    it("Should revert when contract is paused", async function () {
      await adapter.pause();

      await expect(
        adapter.connect(user1).borrow(BORROW_AMOUNT)
      ).to.be.revertedWithCustomError(adapter, "EnforcedPause");
    });
  });

  describe("Repay", function () {
    beforeEach(async function () {
      // Supply collateral and borrow
      await adapter
        .connect(user1)
        .supplyCollateral(
          await collateralToken.getAddress(),
          ethers.parseEther("10")
        );
      await adapter.connect(user1).borrow(BORROW_AMOUNT);
    });

    it("Should allow users to repay debt", async function () {
      const repayAmount = BORROW_AMOUNT / 2n;

      await expect(adapter.connect(user1).repay(repayAmount))
        .to.emit(adapter, "Repaid")
        .withArgs(user1.address, repayAmount);

      const position = await adapter.positions(user1.address);
      expect(position.baseBorrowed).to.equal(BORROW_AMOUNT - repayAmount);
    });

    it("Should transfer tokens correctly", async function () {
      const repayAmount = BORROW_AMOUNT / 2n;
      const balanceBefore = await baseToken.balanceOf(user1.address);

      await adapter.connect(user1).repay(repayAmount);

      const balanceAfter = await baseToken.balanceOf(user1.address);
      expect(balanceBefore - balanceAfter).to.equal(repayAmount);
    });

    it("Should update statistics correctly", async function () {
      const repayAmount = BORROW_AMOUNT / 2n;

      await adapter.connect(user1).repay(repayAmount);

      expect(await adapter.totalBorrowed()).to.equal(BORROW_AMOUNT - repayAmount);
    });

    it("Should allow full repayment", async function () {
      await adapter.connect(user1).repay(BORROW_AMOUNT);

      const position = await adapter.positions(user1.address);
      expect(position.baseBorrowed).to.equal(0);
    });

    it("Should handle overpayment correctly", async function () {
      const overpayment = BORROW_AMOUNT + ethers.parseEther("100");

      await adapter.connect(user1).repay(overpayment);

      const position = await adapter.positions(user1.address);
      expect(position.baseBorrowed).to.equal(0);
    });

    it("Should deactivate position when fully repaid with no supply", async function () {
      await adapter.connect(user1).repay(BORROW_AMOUNT);

      const position = await adapter.positions(user1.address);
      expect(position.active).to.be.false;
      expect(await adapter.activePositions()).to.equal(0);
    });

    it("Should revert if amount is zero", async function () {
      await expect(adapter.connect(user1).repay(0)).to.be.revertedWith(
        "Amount must be greater than zero"
      );
    });

    it("Should revert if no debt to repay", async function () {
      await expect(
        adapter.connect(user2).repay(BORROW_AMOUNT)
      ).to.be.revertedWith("No debt to repay");
    });

    it("Should revert when contract is paused", async function () {
      await adapter.pause();

      await expect(
        adapter.connect(user1).repay(BORROW_AMOUNT / 2n)
      ).to.be.revertedWithCustomError(adapter, "EnforcedPause");
    });
  });

  describe("View Functions", function () {
    beforeEach(async function () {
      await adapter.connect(user1).supplyBase(SUPPLY_AMOUNT);
      await adapter
        .connect(user1)
        .supplyCollateral(
          await collateralToken.getAddress(),
          ethers.parseEther("10")
        );
      await adapter.connect(user1).borrow(BORROW_AMOUNT);
    });

    it("Should return correct position data", async function () {
      const [baseSupplied, baseBorrowed, collateralValue, lastUpdate] =
        await adapter.getPosition(user1.address);

      expect(baseSupplied).to.equal(SUPPLY_AMOUNT);
      expect(baseBorrowed).to.equal(BORROW_AMOUNT);
      expect(collateralValue).to.equal(0); // Not implemented in MVP
      expect(lastUpdate).to.be.gt(0);
    });

    it("Should return current rates", async function () {
      const [supplyRate, borrowRate] = await adapter.getRates();

      expect(supplyRate).to.be.gt(0);
      expect(borrowRate).to.be.gt(0);
      expect(borrowRate).to.be.gt(supplyRate);
    });

    it("Should check position health correctly", async function () {
      const isHealthy = await adapter.isPositionHealthy(user1.address);
      expect(isHealthy).to.be.true;
    });
  });

  describe("Admin Functions", function () {
    it("Should allow admin to set state aggregator", async function () {
      const AggregatorFactory = await ethers.getContractFactory(
        "L2StateAggregator"
      );
      const newAggregator = await AggregatorFactory.deploy(owner.address);

      await adapter.setStateAggregator(await newAggregator.getAddress());
      expect(await adapter.stateAggregator()).to.equal(
        await newAggregator.getAddress()
      );
    });

    it("Should allow manager to pause contract", async function () {
      await adapter.pause();

      await expect(
        adapter.connect(user1).supplyBase(SUPPLY_AMOUNT)
      ).to.be.revertedWithCustomError(adapter, "EnforcedPause");
    });

    it("Should allow manager to unpause contract", async function () {
      await adapter.pause();
      await adapter.unpause();

      await expect(adapter.connect(user1).supplyBase(SUPPLY_AMOUNT)).to.not.be
        .reverted;
    });

    it("Should allow admin to emergency withdraw tokens", async function () {
      await baseToken.transfer(await adapter.getAddress(), ethers.parseEther("100"));

      const balanceBefore = await baseToken.balanceOf(owner.address);
      await adapter.emergencyWithdraw(
        await baseToken.getAddress(),
        ethers.parseEther("100")
      );
      const balanceAfter = await baseToken.balanceOf(owner.address);

      expect(balanceAfter - balanceBefore).to.equal(ethers.parseEther("100"));
    });

    it("Should revert if non-admin tries to set aggregator", async function () {
      const AggregatorFactory = await ethers.getContractFactory(
        "L2StateAggregator"
      );
      const newAggregator = await AggregatorFactory.deploy(owner.address);

      await expect(
        adapter.connect(user1).setStateAggregator(await newAggregator.getAddress())
      ).to.be.reverted;
    });

    it("Should revert if non-manager tries to pause", async function () {
      await expect(adapter.connect(user1).pause()).to.be.reverted;
    });

    it("Should revert if non-admin tries to emergency withdraw", async function () {
      await expect(
        adapter
          .connect(user1)
          .emergencyWithdraw(await baseToken.getAddress(), ethers.parseEther("1"))
      ).to.be.reverted;
    });
  });

  describe("Complex Scenarios", function () {
    it("Should handle full lifecycle: supply -> borrow -> repay -> withdraw", async function () {
      // Supply base
      await adapter.connect(user1).supplyBase(SUPPLY_AMOUNT);
      let [baseSupplied, baseBorrowed] = await adapter.getPosition(user1.address);
      expect(baseSupplied).to.equal(SUPPLY_AMOUNT);
      expect(baseBorrowed).to.equal(0);

      // Supply collateral
      await adapter
        .connect(user1)
        .supplyCollateral(
          await collateralToken.getAddress(),
          ethers.parseEther("10")
        );

      // Borrow
      await adapter.connect(user1).borrow(BORROW_AMOUNT);
      [baseSupplied, baseBorrowed] = await adapter.getPosition(user1.address);
      expect(baseBorrowed).to.equal(BORROW_AMOUNT);

      // Repay
      await adapter.connect(user1).repay(BORROW_AMOUNT);
      [baseSupplied, baseBorrowed] = await adapter.getPosition(user1.address);
      expect(baseBorrowed).to.equal(0);

      // Withdraw base
      await adapter.connect(user1).withdrawBase(SUPPLY_AMOUNT);
      [baseSupplied, baseBorrowed] = await adapter.getPosition(user1.address);
      expect(baseSupplied).to.equal(0);
    });

    it("Should handle multiple users with independent positions", async function () {
      const supply1 = ethers.parseEther("1000");
      const supply2 = ethers.parseEther("2000");
      const borrow1 = ethers.parseEther("500");
      const borrow2 = ethers.parseEther("800");

      await adapter.connect(user1).supplyBase(supply1);
      await adapter.connect(user2).supplyBase(supply2);

      await adapter
        .connect(user1)
        .supplyCollateral(
          await collateralToken.getAddress(),
          ethers.parseEther("10")
        );
      await adapter
        .connect(user2)
        .supplyCollateral(
          await collateralToken.getAddress(),
          ethers.parseEther("20")
        );

      await adapter.connect(user1).borrow(borrow1);
      await adapter.connect(user2).borrow(borrow2);

      const [baseSupplied1, baseBorrowed1] = await adapter.getPosition(
        user1.address
      );
      const [baseSupplied2, baseBorrowed2] = await adapter.getPosition(
        user2.address
      );

      expect(baseSupplied1).to.equal(supply1);
      expect(baseBorrowed1).to.equal(borrow1);
      expect(baseSupplied2).to.equal(supply2);
      expect(baseBorrowed2).to.equal(borrow2);
    });

    it("Should maintain accurate totals across operations", async function () {
      const supply1 = ethers.parseEther("1000");
      const supply2 = ethers.parseEther("500");
      const borrow1 = ethers.parseEther("300");
      const repay1 = ethers.parseEther("100");
      const withdraw1 = ethers.parseEther("200");

      await adapter.connect(user1).supplyBase(supply1);
      await adapter.connect(user2).supplyBase(supply2);
      await adapter
        .connect(user1)
        .supplyCollateral(
          await collateralToken.getAddress(),
          ethers.parseEther("10")
        );
      await adapter.connect(user1).borrow(borrow1);
      await adapter.connect(user1).repay(repay1);
      await adapter.connect(user1).withdrawBase(withdraw1);

      expect(await adapter.totalSupplied()).to.equal(supply1 + supply2 - withdraw1);
      expect(await adapter.totalBorrowed()).to.equal(borrow1 - repay1);
    });

    it("Should receive ETH correctly", async function () {
      const [signer] = await ethers.getSigners();
      await signer.sendTransaction({
        to: await adapter.getAddress(),
        value: ethers.parseEther("1"),
      });

      const balance = await ethers.provider.getBalance(await adapter.getAddress());
      expect(balance).to.equal(ethers.parseEther("1"));
    });
  });
});
