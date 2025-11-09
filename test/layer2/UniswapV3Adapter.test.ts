import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type {
  UniswapV3Adapter,
  MockSwapRouter,
  MockUniswapV3Factory,
  MockERC20,
  L2StateAggregator,
} from "../../typechain-types/index.js";

describe("UniswapV3Adapter", function () {
  let adapter: UniswapV3Adapter;
  let swapRouter: MockSwapRouter;
  let factory: MockUniswapV3Factory;
  let tokenA: MockERC20;
  let tokenB: MockERC20;
  let tokenC: MockERC20;
  let stateAggregator: L2StateAggregator;

  let owner: SignerWithAddress;
  let user1: SignerWithAddress;
  let user2: SignerWithAddress;

  const INITIAL_SUPPLY = ethers.parseEther("1000000");
  const FEE_LOW = 500; // 0.05%
  const FEE_MEDIUM = 3000; // 0.30%

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();

    // Deploy mock ERC20 tokens
    const MockERC20Factory = await ethers.getContractFactory("MockERC20");
    tokenA = await MockERC20Factory.deploy("Token A", "TKA", INITIAL_SUPPLY);
    tokenB = await MockERC20Factory.deploy("Token B", "TKB", INITIAL_SUPPLY);
    tokenC = await MockERC20Factory.deploy("Token C", "TKC", INITIAL_SUPPLY);

    // Deploy mock Uniswap V3 contracts
    const SwapRouterFactory = await ethers.getContractFactory("MockSwapRouter");
    swapRouter = await SwapRouterFactory.deploy();

    const FactoryFactory = await ethers.getContractFactory(
      "MockUniswapV3Factory"
    );
    factory = await FactoryFactory.deploy();

    // Deploy L2StateAggregator
    const AggregatorFactory = await ethers.getContractFactory(
      "L2StateAggregator"
    );
    stateAggregator = await AggregatorFactory.deploy(owner.address);

    // Deploy UniswapV3Adapter
    const AdapterFactory = await ethers.getContractFactory("UniswapV3Adapter");
    adapter = await AdapterFactory.deploy(
      await swapRouter.getAddress(),
      await factory.getAddress(),
      owner.address
    );

    // Set state aggregator
    await adapter.setStateAggregator(await stateAggregator.getAddress());

    // Authorize adapter as a module in the state aggregator
    await stateAggregator.authorizeModule(await adapter.getAddress());

    // Setup exchange rates in mock router
    // 1 TKA = 2 TKB (rate = 2e18)
    await swapRouter.setExchangeRate(
      await tokenA.getAddress(),
      await tokenB.getAddress(),
      ethers.parseEther("2")
    );
    // 1 TKB = 0.5 TKA (inverse rate)
    await swapRouter.setExchangeRate(
      await tokenB.getAddress(),
      await tokenA.getAddress(),
      ethers.parseEther("0.5")
    );
    // 1 TKA = 3 TKC (for multi-hop testing)
    await swapRouter.setExchangeRate(
      await tokenA.getAddress(),
      await tokenC.getAddress(),
      ethers.parseEther("3")
    );

    // Create mock pools
    await factory.createPool(
      await tokenA.getAddress(),
      await tokenB.getAddress(),
      FEE_MEDIUM
    );

    // Fund swap router with liquidity
    await tokenB.transfer(await swapRouter.getAddress(), INITIAL_SUPPLY / 2n);
    await tokenC.transfer(await swapRouter.getAddress(), INITIAL_SUPPLY / 2n);

    // Distribute tokens to users
    await tokenA.transfer(user1.address, ethers.parseEther("1000"));
    await tokenA.transfer(user2.address, ethers.parseEther("1000"));
    await tokenB.transfer(user1.address, ethers.parseEther("1000"));
    await tokenB.transfer(user2.address, ethers.parseEther("1000"));

    // Approve adapter to spend tokens
    await tokenA
      .connect(user1)
      .approve(await adapter.getAddress(), ethers.MaxUint256);
    await tokenB
      .connect(user1)
      .approve(await adapter.getAddress(), ethers.MaxUint256);
    await tokenA
      .connect(user2)
      .approve(await adapter.getAddress(), ethers.MaxUint256);
    await tokenB
      .connect(user2)
      .approve(await adapter.getAddress(), ethers.MaxUint256);
  });

  describe("Deployment", function () {
    it("Should set the correct swap router", async function () {
      expect(await adapter.swapRouter()).to.equal(
        await swapRouter.getAddress()
      );
    });

    it("Should set the correct factory", async function () {
      expect(await adapter.factory()).to.equal(await factory.getAddress());
    });

    it("Should set the correct admin roles", async function () {
      const DEFAULT_ADMIN_ROLE = ethers.ZeroHash;
      const MANAGER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("MANAGER_ROLE"));

      expect(await adapter.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be
        .true;
      expect(await adapter.hasRole(MANAGER_ROLE, owner.address)).to.be.true;
    });

    it("Should have correct fee tier constants", async function () {
      expect(await adapter.FEE_LOWEST()).to.equal(100);
      expect(await adapter.FEE_LOW()).to.equal(500);
      expect(await adapter.FEE_MEDIUM()).to.equal(3000);
      expect(await adapter.FEE_HIGH()).to.equal(10000);
    });

    it("Should initialize with zero statistics", async function () {
      expect(await adapter.totalSwapVolume()).to.equal(0);
      expect(await adapter.totalSwaps()).to.equal(0);
    });

    it("Should revert with invalid swap router address", async function () {
      const AdapterFactory = await ethers.getContractFactory("UniswapV3Adapter");
      await expect(
        AdapterFactory.deploy(
          ethers.ZeroAddress,
          await factory.getAddress(),
          owner.address
        )
      ).to.be.revertedWith("Invalid swap router");
    });

    it("Should revert with invalid factory address", async function () {
      const AdapterFactory = await ethers.getContractFactory("UniswapV3Adapter");
      await expect(
        AdapterFactory.deploy(
          await swapRouter.getAddress(),
          ethers.ZeroAddress,
          owner.address
        )
      ).to.be.revertedWith("Invalid factory");
    });

    it("Should revert with invalid admin address", async function () {
      const AdapterFactory = await ethers.getContractFactory("UniswapV3Adapter");
      await expect(
        AdapterFactory.deploy(
          await swapRouter.getAddress(),
          await factory.getAddress(),
          ethers.ZeroAddress
        )
      ).to.be.revertedWith("Invalid admin");
    });
  });

  describe("Exact Input Single Swap", function () {
    const swapAmount = ethers.parseEther("10");
    const expectedOutput = ethers.parseEther("20"); // 10 * 2 = 20

    it("Should swap exact input for output (single hop)", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      await expect(
        adapter
          .connect(user1)
          .swapExactInputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            swapAmount,
            expectedOutput,
            deadline
          )
      )
        .to.emit(adapter, "Swapped")
        .withArgs(
          user1.address,
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          swapAmount,
          expectedOutput,
          FEE_MEDIUM
        );
    });

    it("Should transfer tokens correctly", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      const tokenABalanceBefore = await tokenA.balanceOf(user1.address);
      const tokenBBalanceBefore = await tokenB.balanceOf(user1.address);

      await adapter
        .connect(user1)
        .swapExactInputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_MEDIUM,
          swapAmount,
          expectedOutput,
          deadline
        );

      const tokenABalanceAfter = await tokenA.balanceOf(user1.address);
      const tokenBBalanceAfter = await tokenB.balanceOf(user1.address);

      expect(tokenABalanceBefore - tokenABalanceAfter).to.equal(swapAmount);
      expect(tokenBBalanceAfter - tokenBBalanceBefore).to.equal(expectedOutput);
    });

    it("Should update statistics correctly", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      await adapter
        .connect(user1)
        .swapExactInputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_MEDIUM,
          swapAmount,
          expectedOutput,
          deadline
        );

      expect(await adapter.totalSwapVolume()).to.equal(swapAmount);
      expect(await adapter.totalSwaps()).to.equal(1);
    });

    it("Should update state aggregator", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      await adapter
        .connect(user1)
        .swapExactInputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_MEDIUM,
          swapAmount,
          expectedOutput,
          deadline
        );

      const state = await stateAggregator.getSystemState();
      // Adapter maps totalSwaps -> activePositions and totalSwapVolume -> totalOrders
      expect(state.activePositions).to.equal(1);
      expect(state.totalOrders).to.equal(swapAmount);
    });

    it("Should track multiple swaps correctly", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      // First swap
      await adapter
        .connect(user1)
        .swapExactInputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_MEDIUM,
          swapAmount,
          expectedOutput,
          deadline
        );

      // Second swap
      await adapter
        .connect(user2)
        .swapExactInputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_MEDIUM,
          swapAmount,
          expectedOutput,
          deadline
        );

      expect(await adapter.totalSwapVolume()).to.equal(swapAmount * 2n);
      expect(await adapter.totalSwaps()).to.equal(2);
    });

    it("Should revert if amount is zero", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      await expect(
        adapter
          .connect(user1)
          .swapExactInputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            0,
            0,
            deadline
          )
      ).to.be.revertedWith("Invalid input amount");
    });

    it("Should revert if deadline has passed", async function () {
      const pastDeadline = (await ethers.provider.getBlock("latest"))!.timestamp - 1;

      await expect(
        adapter
          .connect(user1)
          .swapExactInputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            swapAmount,
            expectedOutput,
            pastDeadline
          )
      ).to.be.revertedWith("Deadline passed");
    });

    it("Should revert if slippage tolerance exceeded", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;
      const excessiveMinOutput = ethers.parseEther("30"); // Expecting more than possible

      await expect(
        adapter
          .connect(user1)
          .swapExactInputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            swapAmount,
            excessiveMinOutput,
            deadline
          )
      ).to.be.revertedWith("Insufficient output");
    });

    it("Should revert if user has insufficient balance", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;
      const largeAmount = ethers.parseEther("10000");

      await expect(
        adapter
          .connect(user1)
          .swapExactInputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            largeAmount,
            0,
            deadline
          )
      ).to.be.reverted;
    });

    it("Should revert when contract is paused", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      await adapter.pause();

      await expect(
        adapter
          .connect(user1)
          .swapExactInputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            swapAmount,
            expectedOutput,
            deadline
          )
      ).to.be.revertedWithCustomError(adapter, "EnforcedPause");
    });
  });

  describe("Exact Input Multi-Hop Swap", function () {
    const swapAmount = ethers.parseEther("10");
    const expectedOutput = ethers.parseEther("30"); // 10 * 3 = 30 (via multi-hop)

    it("Should swap with multi-hop path", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      // Encode path: tokenA -> (fee) -> tokenB -> (fee) -> tokenC
      const path = ethers.solidityPacked(
        ["address", "uint24", "address", "uint24", "address"],
        [
          await tokenA.getAddress(),
          FEE_MEDIUM,
          await tokenB.getAddress(),
          FEE_MEDIUM,
          await tokenC.getAddress(),
        ]
      );

      await expect(
        adapter
          .connect(user1)
          .swapExactInput(path, swapAmount, expectedOutput, deadline)
      )
        .to.emit(adapter, "MultiHopSwap")
        .withArgs(user1.address, swapAmount, expectedOutput);
    });

    it("Should transfer tokens correctly in multi-hop", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      const path = ethers.solidityPacked(
        ["address", "uint24", "address", "uint24", "address"],
        [
          await tokenA.getAddress(),
          FEE_MEDIUM,
          await tokenB.getAddress(),
          FEE_MEDIUM,
          await tokenC.getAddress(),
        ]
      );

      const tokenABalanceBefore = await tokenA.balanceOf(user1.address);
      const tokenCBalanceBefore = await tokenC.balanceOf(user1.address);

      await adapter
        .connect(user1)
        .swapExactInput(path, swapAmount, expectedOutput, deadline);

      const tokenABalanceAfter = await tokenA.balanceOf(user1.address);
      const tokenCBalanceAfter = await tokenC.balanceOf(user1.address);

      expect(tokenABalanceBefore - tokenABalanceAfter).to.equal(swapAmount);
      expect(tokenCBalanceAfter - tokenCBalanceBefore).to.equal(expectedOutput);
    });

    it("Should update statistics for multi-hop swap", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      const path = ethers.solidityPacked(
        ["address", "uint24", "address", "uint24", "address"],
        [
          await tokenA.getAddress(),
          FEE_MEDIUM,
          await tokenB.getAddress(),
          FEE_MEDIUM,
          await tokenC.getAddress(),
        ]
      );

      await adapter
        .connect(user1)
        .swapExactInput(path, swapAmount, expectedOutput, deadline);

      expect(await adapter.totalSwapVolume()).to.equal(swapAmount);
      expect(await adapter.totalSwaps()).to.equal(1);
    });

    it("Should revert with invalid path", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      // Path too short
      const invalidPath = "0x1234";

      await expect(
        adapter
          .connect(user1)
          .swapExactInput(invalidPath, swapAmount, expectedOutput, deadline)
      ).to.be.revertedWith("Invalid path");
    });

    it("Should revert if amount is zero", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      const path = ethers.solidityPacked(
        ["address", "uint24", "address"],
        [await tokenA.getAddress(), FEE_MEDIUM, await tokenB.getAddress()]
      );

      await expect(
        adapter.connect(user1).swapExactInput(path, 0, 0, deadline)
      ).to.be.revertedWith("Invalid input amount");
    });

    it("Should revert if deadline has passed", async function () {
      const pastDeadline = (await ethers.provider.getBlock("latest"))!.timestamp - 1;

      const path = ethers.solidityPacked(
        ["address", "uint24", "address"],
        [await tokenA.getAddress(), FEE_MEDIUM, await tokenB.getAddress()]
      );

      await expect(
        adapter
          .connect(user1)
          .swapExactInput(path, swapAmount, expectedOutput, pastDeadline)
      ).to.be.revertedWith("Deadline passed");
    });
  });

  describe("Exact Output Single Swap", function () {
    const desiredOutput = ethers.parseEther("20");
    const maxInput = ethers.parseEther("15"); // Should need ~10 tokens

    it("Should swap for exact output amount", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      await expect(
        adapter
          .connect(user1)
          .swapExactOutputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            desiredOutput,
            maxInput,
            deadline
          )
      )
        .to.emit(adapter, "Swapped")
        .withArgs(
          user1.address,
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          ethers.parseEther("10"), // Actual input needed
          desiredOutput,
          FEE_MEDIUM
        );
    });

    it("Should refund unused tokens", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      const tokenABalanceBefore = await tokenA.balanceOf(user1.address);

      await adapter
        .connect(user1)
        .swapExactOutputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_MEDIUM,
          desiredOutput,
          maxInput,
          deadline
        );

      const tokenABalanceAfter = await tokenA.balanceOf(user1.address);

      // Should only spend 10 tokens, not the full maxInput of 15
      expect(tokenABalanceBefore - tokenABalanceAfter).to.equal(
        ethers.parseEther("10")
      );
    });

    it("Should update statistics correctly", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      await adapter
        .connect(user1)
        .swapExactOutputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_MEDIUM,
          desiredOutput,
          maxInput,
          deadline
        );

      expect(await adapter.totalSwapVolume()).to.equal(ethers.parseEther("10"));
      expect(await adapter.totalSwaps()).to.equal(1);
    });

    it("Should revert if amount is zero", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      await expect(
        adapter
          .connect(user1)
          .swapExactOutputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            0,
            maxInput,
            deadline
          )
      ).to.be.revertedWith("Invalid output amount");
    });

    it("Should revert if deadline has passed", async function () {
      const pastDeadline = (await ethers.provider.getBlock("latest"))!.timestamp - 1;

      await expect(
        adapter
          .connect(user1)
          .swapExactOutputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            desiredOutput,
            maxInput,
            pastDeadline
          )
      ).to.be.revertedWith("Deadline passed");
    });

    it("Should revert if max input exceeded", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;
      const insufficientMaxInput = ethers.parseEther("5"); // Not enough

      await expect(
        adapter
          .connect(user1)
          .swapExactOutputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            desiredOutput,
            insufficientMaxInput,
            deadline
          )
      ).to.be.revertedWith("Excessive input required");
    });
  });

  describe("Pool Queries", function () {
    it("Should get pool address for token pair", async function () {
      const pool = await adapter.getPool(
        await tokenA.getAddress(),
        await tokenB.getAddress(),
        FEE_MEDIUM
      );

      expect(pool).to.not.equal(ethers.ZeroAddress);
    });

    it("Should check if pool exists", async function () {
      const exists = await adapter.poolExists(
        await tokenA.getAddress(),
        await tokenB.getAddress(),
        FEE_MEDIUM
      );

      expect(exists).to.be.true;
    });

    it("Should return false for non-existent pool", async function () {
      const exists = await adapter.poolExists(
        await tokenA.getAddress(),
        await tokenC.getAddress(),
        FEE_MEDIUM
      );

      expect(exists).to.be.false;
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
      // Contract should be paused - test by trying a swap
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;
      await expect(
        adapter
          .connect(user1)
          .swapExactInputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            ethers.parseEther("1"),
            0,
            deadline
          )
      ).to.be.revertedWithCustomError(adapter, "EnforcedPause");
    });

    it("Should allow manager to unpause contract", async function () {
      await adapter.pause();
      await adapter.unpause();

      // Should be able to swap again
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;
      await expect(
        adapter
          .connect(user1)
          .swapExactInputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            ethers.parseEther("1"),
            ethers.parseEther("1"),
            deadline
          )
      ).to.not.be.reverted;
    });

    it("Should allow admin to emergency withdraw tokens", async function () {
      // Send some tokens to adapter by mistake
      await tokenA.transfer(await adapter.getAddress(), ethers.parseEther("100"));

      const balanceBefore = await tokenA.balanceOf(owner.address);
      await adapter.emergencyWithdraw(
        await tokenA.getAddress(),
        ethers.parseEther("100")
      );
      const balanceAfter = await tokenA.balanceOf(owner.address);

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
          .emergencyWithdraw(await tokenA.getAddress(), ethers.parseEther("1"))
      ).to.be.reverted;
    });
  });

  describe("Edge Cases and Security", function () {
    it("Should handle dust amounts correctly", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;
      const dustAmount = 1n;

      await adapter
        .connect(user1)
        .swapExactInputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_MEDIUM,
          dustAmount,
          0,
          deadline
        );

      expect(await adapter.totalSwapVolume()).to.equal(dustAmount);
    });

    it("Should prevent reentrancy attacks", async function () {
      // NonReentrant modifier should prevent reentrancy
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;
      await adapter
        .connect(user1)
        .swapExactInputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_MEDIUM,
          ethers.parseEther("1"),
          0,
          deadline
        );
      // If vulnerable, this would fail during execution
    });

    it("Should handle multiple concurrent swaps", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;
      const swapAmount = ethers.parseEther("5");

      // Multiple users swap at same time
      await Promise.all([
        adapter
          .connect(user1)
          .swapExactInputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            swapAmount,
            0,
            deadline
          ),
        adapter
          .connect(user2)
          .swapExactInputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            swapAmount,
            0,
            deadline
          ),
      ]);

      expect(await adapter.totalSwapVolume()).to.equal(swapAmount * 2n);
      expect(await adapter.totalSwaps()).to.equal(2);
    });

    it("Should maintain accurate statistics across different swap types", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      // Single hop swap
      await adapter
        .connect(user1)
        .swapExactInputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_MEDIUM,
          ethers.parseEther("10"),
          0,
          deadline
        );

      // Multi-hop swap
      const path = ethers.solidityPacked(
        ["address", "uint24", "address"],
        [await tokenA.getAddress(), FEE_MEDIUM, await tokenB.getAddress()]
      );
      await adapter
        .connect(user1)
        .swapExactInput(path, ethers.parseEther("5"), 0, deadline);

      // Exact output swap
      await adapter
        .connect(user1)
        .swapExactOutputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_MEDIUM,
          ethers.parseEther("10"),
          ethers.parseEther("10"),
          deadline
        );

      expect(await adapter.totalSwaps()).to.equal(3);
      expect(await adapter.totalSwapVolume()).to.equal(ethers.parseEther("20")); // 10 + 5 + 5
    });

    it("Should handle different fee tiers", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      // Create pool with different fee
      await factory.createPool(
        await tokenA.getAddress(),
        await tokenB.getAddress(),
        FEE_LOW
      );

      await adapter
        .connect(user1)
        .swapExactInputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_LOW,
          ethers.parseEther("1"),
          0,
          deadline
        );

      expect(await adapter.totalSwaps()).to.equal(1);
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

  describe("Integration with State Aggregator", function () {
    it("Should update aggregator even when aggregator is not set", async function () {
      // Deploy new adapter without aggregator
      const AdapterFactory = await ethers.getContractFactory("UniswapV3Adapter");
      const newAdapter = await AdapterFactory.deploy(
        await swapRouter.getAddress(),
        await factory.getAddress(),
        owner.address
      );

      // Approve tokens for new adapter
      await tokenA
        .connect(user1)
        .approve(await newAdapter.getAddress(), ethers.MaxUint256);

      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      // Should not revert even without aggregator
      await expect(
        newAdapter
          .connect(user1)
          .swapExactInputSingle(
            await tokenA.getAddress(),
            await tokenB.getAddress(),
            FEE_MEDIUM,
            ethers.parseEther("1"),
            0,
            deadline
          )
      ).to.not.be.reverted;
    });

    it("Should correctly update aggregator state", async function () {
      const deadline = (await ethers.provider.getBlock("latest"))!.timestamp + 3600;

      await adapter
        .connect(user1)
        .swapExactInputSingle(
          await tokenA.getAddress(),
          await tokenB.getAddress(),
          FEE_MEDIUM,
          ethers.parseEther("100"),
          0,
          deadline
        );

      const state = await stateAggregator.getSystemState();
      // Adapter maps totalSwaps -> activePositions and totalSwapVolume -> totalOrders
      expect(state.activePositions).to.equal(1);
      expect(state.totalOrders).to.equal(ethers.parseEther("100"));
    });
  });
});
