import { expect } from "chai";
import { ethers } from "hardhat";
import { loadFixture } from "@nomicfoundation/hardhat-toolbox/network-helpers";

describe("LoyaltyUSD", function () {
  // Fixture for deploying contract
  async function deployLUSDFixture() {
    const [owner, minter, burner, pauser, user1, user2] = await ethers.getSigners();

    const LoyaltyUSD = await ethers.getContractFactory("LoyaltyUSD");
    const lusd = await LoyaltyUSD.deploy();

    return { lusd, owner, minter, burner, pauser, user1, user2 };
  }

  describe("Deployment", function () {
    it("Should set the correct name and symbol", async function () {
      const { lusd } = await loadFixture(deployLUSDFixture);
      expect(await lusd.name()).to.equal("LoyaltyUSD");
      expect(await lusd.symbol()).to.equal("LUSD");
    });

    it("Should set 6 decimals", async function () {
      const { lusd } = await loadFixture(deployLUSDFixture);
      expect(await lusd.decimals()).to.equal(6);
    });

    it("Should grant DEFAULT_ADMIN_ROLE to deployer", async function () {
      const { lusd, owner } = await loadFixture(deployLUSDFixture);
      const DEFAULT_ADMIN_ROLE = await lusd.DEFAULT_ADMIN_ROLE();
      expect(await lusd.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be.true;
    });

    it("Should grant PAUSER_ROLE to deployer", async function () {
      const { lusd, owner } = await loadFixture(deployLUSDFixture);
      const PAUSER_ROLE = await lusd.PAUSER_ROLE();
      expect(await lusd.hasRole(PAUSER_ROLE, owner.address)).to.be.true;
    });

    it("Should NOT grant MINTER_ROLE to deployer", async function () {
      const { lusd, owner } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      expect(await lusd.hasRole(MINTER_ROLE, owner.address)).to.be.false;
    });

    it("Should start with zero total supply", async function () {
      const { lusd } = await loadFixture(deployLUSDFixture);
      expect(await lusd.totalSupply()).to.equal(0);
    });
  });

  describe("Role Management", function () {
    it("Should allow admin to grant MINTER_ROLE", async function () {
      const { lusd, owner, minter } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();

      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);
      expect(await lusd.hasRole(MINTER_ROLE, minter.address)).to.be.true;
    });

    it("Should allow admin to grant BURNER_ROLE", async function () {
      const { lusd, owner, burner } = await loadFixture(deployLUSDFixture);
      const BURNER_ROLE = await lusd.BURNER_ROLE();

      await lusd.connect(owner).grantRole(BURNER_ROLE, burner.address);
      expect(await lusd.hasRole(BURNER_ROLE, burner.address)).to.be.true;
    });

    it("Should allow admin to revoke roles", async function () {
      const { lusd, owner, minter } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();

      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);
      await lusd.connect(owner).revokeRole(MINTER_ROLE, minter.address);

      expect(await lusd.hasRole(MINTER_ROLE, minter.address)).to.be.false;
    });

    it("Should prevent non-admin from granting roles", async function () {
      const { lusd, user1, minter } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();

      await expect(
        lusd.connect(user1).grantRole(MINTER_ROLE, minter.address)
      ).to.be.reverted;
    });

    it("Should check canMint correctly", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();

      expect(await lusd.canMint(minter.address)).to.be.false;

      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);
      expect(await lusd.canMint(minter.address)).to.be.true;
      expect(await lusd.canMint(user1.address)).to.be.false;
    });

    it("Should check canBurn correctly", async function () {
      const { lusd, owner, burner, user1 } = await loadFixture(deployLUSDFixture);
      const BURNER_ROLE = await lusd.BURNER_ROLE();

      expect(await lusd.canBurn(burner.address)).to.be.false;

      await lusd.connect(owner).grantRole(BURNER_ROLE, burner.address);
      expect(await lusd.canBurn(burner.address)).to.be.true;
      expect(await lusd.canBurn(user1.address)).to.be.false;
    });
  });

  describe("Minting", function () {
    it("Should allow MINTER_ROLE to mint tokens", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const mintAmount = ethers.parseUnits("1000", 6);

      await expect(lusd.connect(minter).mint(user1.address, mintAmount))
        .to.emit(lusd, "Minted")
        .withArgs(user1.address, mintAmount, minter.address);

      expect(await lusd.balanceOf(user1.address)).to.equal(mintAmount);
      expect(await lusd.totalSupply()).to.equal(mintAmount);
    });

    it("Should prevent non-MINTER from minting", async function () {
      const { lusd, user1, user2 } = await loadFixture(deployLUSDFixture);
      const mintAmount = ethers.parseUnits("1000", 6);

      await expect(
        lusd.connect(user1).mint(user2.address, mintAmount)
      ).to.be.reverted;
    });

    it("Should revert when minting to zero address", async function () {
      const { lusd, owner, minter } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      await expect(
        lusd.connect(minter).mint(ethers.ZeroAddress, ethers.parseUnits("1000", 6))
      ).to.be.revertedWith("LUSD: mint to zero address");
    });

    it("Should revert when minting zero amount", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      await expect(
        lusd.connect(minter).mint(user1.address, 0)
      ).to.be.revertedWith("LUSD: mint amount must be positive");
    });

    it("Should handle multiple mints", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const amount1 = ethers.parseUnits("500", 6);
      const amount2 = ethers.parseUnits("300", 6);

      await lusd.connect(minter).mint(user1.address, amount1);
      await lusd.connect(minter).mint(user1.address, amount2);

      expect(await lusd.balanceOf(user1.address)).to.equal(amount1 + amount2);
    });

    it("Should mint to multiple users", async function () {
      const { lusd, owner, minter, user1, user2 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const amount1 = ethers.parseUnits("1000", 6);
      const amount2 = ethers.parseUnits("500", 6);

      await lusd.connect(minter).mint(user1.address, amount1);
      await lusd.connect(minter).mint(user2.address, amount2);

      expect(await lusd.totalSupply()).to.equal(amount1 + amount2);
    });
  });

  describe("Burning", function () {
    it("Should allow token owner to burn their own tokens", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const mintAmount = ethers.parseUnits("1000", 6);
      const burnAmount = ethers.parseUnits("300", 6);

      await lusd.connect(minter).mint(user1.address, mintAmount);

      await expect(lusd.connect(user1).burn(user1.address, burnAmount))
        .to.emit(lusd, "Burned")
        .withArgs(user1.address, burnAmount, user1.address);

      expect(await lusd.balanceOf(user1.address)).to.equal(mintAmount - burnAmount);
    });

    it("Should allow BURNER_ROLE to burn from any address", async function () {
      const { lusd, owner, minter, burner, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      const BURNER_ROLE = await lusd.BURNER_ROLE();

      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);
      await lusd.connect(owner).grantRole(BURNER_ROLE, burner.address);

      const mintAmount = ethers.parseUnits("1000", 6);
      const burnAmount = ethers.parseUnits("400", 6);

      await lusd.connect(minter).mint(user1.address, mintAmount);
      await lusd.connect(burner).burn(user1.address, burnAmount);

      expect(await lusd.balanceOf(user1.address)).to.equal(mintAmount - burnAmount);
    });

    it("Should prevent unauthorized burning from other addresses", async function () {
      const { lusd, owner, minter, user1, user2 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      await lusd.connect(minter).mint(user1.address, ethers.parseUnits("1000", 6));

      await expect(
        lusd.connect(user2).burn(user1.address, ethers.parseUnits("100", 6))
      ).to.be.revertedWith("LUSD: must have BURNER_ROLE or be token owner");
    });

    it("Should revert when burning zero amount", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      await lusd.connect(minter).mint(user1.address, ethers.parseUnits("1000", 6));

      await expect(
        lusd.connect(user1).burn(user1.address, 0)
      ).to.be.revertedWith("LUSD: burn amount must be positive");
    });

    it("Should revert when burning more than balance", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const mintAmount = ethers.parseUnits("1000", 6);
      await lusd.connect(minter).mint(user1.address, mintAmount);

      await expect(
        lusd.connect(user1).burn(user1.address, mintAmount + 1n)
      ).to.be.reverted;
    });

    it("Should decrease total supply when burning", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const mintAmount = ethers.parseUnits("1000", 6);
      const burnAmount = ethers.parseUnits("600", 6);

      await lusd.connect(minter).mint(user1.address, mintAmount);
      await lusd.connect(user1).burn(user1.address, burnAmount);

      expect(await lusd.totalSupply()).to.equal(mintAmount - burnAmount);
    });
  });

  describe("Pausable Functionality", function () {
    it("Should allow PAUSER_ROLE to pause", async function () {
      const { lusd, owner } = await loadFixture(deployLUSDFixture);

      await expect(lusd.connect(owner).pause())
        .to.emit(lusd, "EmergencyPaused")
        .withArgs(owner.address);

      expect(await lusd.paused()).to.be.true;
    });

    it("Should prevent non-PAUSER from pausing", async function () {
      const { lusd, user1 } = await loadFixture(deployLUSDFixture);

      await expect(lusd.connect(user1).pause()).to.be.reverted;
    });

    it("Should prevent transfers when paused", async function () {
      const { lusd, owner, minter, user1, user2 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      await lusd.connect(minter).mint(user1.address, ethers.parseUnits("1000", 6));
      await lusd.connect(owner).pause();

      await expect(
        lusd.connect(user1).transfer(user2.address, ethers.parseUnits("100", 6))
      ).to.be.reverted;
    });

    it("Should prevent minting when paused", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      await lusd.connect(owner).pause();

      await expect(
        lusd.connect(minter).mint(user1.address, ethers.parseUnits("1000", 6))
      ).to.be.reverted;
    });

    it("Should prevent burning when paused", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      await lusd.connect(minter).mint(user1.address, ethers.parseUnits("1000", 6));
      await lusd.connect(owner).pause();

      await expect(
        lusd.connect(user1).burn(user1.address, ethers.parseUnits("100", 6))
      ).to.be.reverted;
    });

    it("Should allow unpausing", async function () {
      const { lusd, owner } = await loadFixture(deployLUSDFixture);

      await lusd.connect(owner).pause();

      await expect(lusd.connect(owner).unpause())
        .to.emit(lusd, "EmergencyUnpaused")
        .withArgs(owner.address);

      expect(await lusd.paused()).to.be.false;
    });

    it("Should allow transfers after unpausing", async function () {
      const { lusd, owner, minter, user1, user2 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const mintAmount = ethers.parseUnits("1000", 6);
      await lusd.connect(minter).mint(user1.address, mintAmount);

      await lusd.connect(owner).pause();
      await lusd.connect(owner).unpause();

      const transferAmount = ethers.parseUnits("100", 6);
      await expect(
        lusd.connect(user1).transfer(user2.address, transferAmount)
      ).to.not.be.reverted;

      expect(await lusd.balanceOf(user2.address)).to.equal(transferAmount);
    });
  });

  describe("ERC20 Standard Functions", function () {
    it("Should support transfer", async function () {
      const { lusd, owner, minter, user1, user2 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const mintAmount = ethers.parseUnits("1000", 6);
      const transferAmount = ethers.parseUnits("300", 6);

      await lusd.connect(minter).mint(user1.address, mintAmount);
      await lusd.connect(user1).transfer(user2.address, transferAmount);

      expect(await lusd.balanceOf(user1.address)).to.equal(mintAmount - transferAmount);
      expect(await lusd.balanceOf(user2.address)).to.equal(transferAmount);
    });

    it("Should support approve and transferFrom", async function () {
      const { lusd, owner, minter, user1, user2 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const mintAmount = ethers.parseUnits("1000", 6);
      const approveAmount = ethers.parseUnits("500", 6);
      const transferAmount = ethers.parseUnits("300", 6);

      await lusd.connect(minter).mint(user1.address, mintAmount);
      await lusd.connect(user1).approve(user2.address, approveAmount);

      expect(await lusd.allowance(user1.address, user2.address)).to.equal(approveAmount);

      await lusd.connect(user2).transferFrom(user1.address, user2.address, transferAmount);

      expect(await lusd.balanceOf(user2.address)).to.equal(transferAmount);
      expect(await lusd.allowance(user1.address, user2.address)).to.equal(
        approveAmount - transferAmount
      );
    });

    it("Should revert transferFrom without approval", async function () {
      const { lusd, owner, minter, user1, user2 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      await lusd.connect(minter).mint(user1.address, ethers.parseUnits("1000", 6));

      await expect(
        lusd.connect(user2).transferFrom(
          user1.address,
          user2.address,
          ethers.parseUnits("100", 6)
        )
      ).to.be.reverted;
    });
  });

  describe("Utility Functions", function () {
    it("Should return formatted total supply", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const mintAmount = ethers.parseUnits("1000", 6); // 1000 LUSD
      await lusd.connect(minter).mint(user1.address, mintAmount);

      expect(await lusd.totalSupplyFormatted()).to.equal(1000);
    });

    it("Should handle formatted supply with decimals", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      // 1,234.56 LUSD = 1234560000 (6 decimals)
      const mintAmount = 1234560000n;
      await lusd.connect(minter).mint(user1.address, mintAmount);

      expect(await lusd.totalSupplyFormatted()).to.equal(1234n);
    });
  });

  describe("Edge Cases", function () {
    it("Should handle maximum uint256 approval", async function () {
      const { lusd, user1, user2 } = await loadFixture(deployLUSDFixture);

      await lusd.connect(user1).approve(user2.address, ethers.MaxUint256);
      expect(await lusd.allowance(user1.address, user2.address)).to.equal(ethers.MaxUint256);
    });

    it("Should handle zero transfers", async function () {
      const { lusd, owner, minter, user1, user2 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      await lusd.connect(minter).mint(user1.address, ethers.parseUnits("1000", 6));

      await expect(lusd.connect(user1).transfer(user2.address, 0)).to.not.be.reverted;
    });

    it("Should handle self-transfer", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const mintAmount = ethers.parseUnits("1000", 6);
      await lusd.connect(minter).mint(user1.address, mintAmount);

      await lusd.connect(user1).transfer(user1.address, ethers.parseUnits("100", 6));
      expect(await lusd.balanceOf(user1.address)).to.equal(mintAmount);
    });

    it("Should emit Transfer event on mint", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const mintAmount = ethers.parseUnits("1000", 6);

      await expect(lusd.connect(minter).mint(user1.address, mintAmount))
        .to.emit(lusd, "Transfer")
        .withArgs(ethers.ZeroAddress, user1.address, mintAmount);
    });

    it("Should emit Transfer event on burn", async function () {
      const { lusd, owner, minter, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);

      const mintAmount = ethers.parseUnits("1000", 6);
      const burnAmount = ethers.parseUnits("500", 6);

      await lusd.connect(minter).mint(user1.address, mintAmount);

      await expect(lusd.connect(user1).burn(user1.address, burnAmount))
        .to.emit(lusd, "Transfer")
        .withArgs(user1.address, ethers.ZeroAddress, burnAmount);
    });
  });

  describe("Integration Scenarios", function () {
    it("Should handle complete lifecycle: mint, transfer, burn", async function () {
      const { lusd, owner, minter, burner, user1, user2 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      const BURNER_ROLE = await lusd.BURNER_ROLE();

      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);
      await lusd.connect(owner).grantRole(BURNER_ROLE, burner.address);

      // Mint
      const mintAmount = ethers.parseUnits("1000", 6);
      await lusd.connect(minter).mint(user1.address, mintAmount);

      // Transfer
      const transferAmount = ethers.parseUnits("400", 6);
      await lusd.connect(user1).transfer(user2.address, transferAmount);

      // Burn from user1
      const burnAmount1 = ethers.parseUnits("200", 6);
      await lusd.connect(burner).burn(user1.address, burnAmount1);

      // Burn from user2
      const burnAmount2 = ethers.parseUnits("100", 6);
      await lusd.connect(burner).burn(user2.address, burnAmount2);

      expect(await lusd.balanceOf(user1.address)).to.equal(
        mintAmount - transferAmount - burnAmount1
      );
      expect(await lusd.balanceOf(user2.address)).to.equal(
        transferAmount - burnAmount2
      );
      expect(await lusd.totalSupply()).to.equal(
        mintAmount - burnAmount1 - burnAmount2
      );
    });

    it("Should handle multiple minters and burners", async function () {
      const { lusd, owner, minter, burner, pauser, user1 } = await loadFixture(deployLUSDFixture);
      const MINTER_ROLE = await lusd.MINTER_ROLE();
      const BURNER_ROLE = await lusd.BURNER_ROLE();

      // Grant roles to multiple addresses
      await lusd.connect(owner).grantRole(MINTER_ROLE, minter.address);
      await lusd.connect(owner).grantRole(MINTER_ROLE, pauser.address);
      await lusd.connect(owner).grantRole(BURNER_ROLE, burner.address);

      // Mint from different minters
      await lusd.connect(minter).mint(user1.address, ethers.parseUnits("500", 6));
      await lusd.connect(pauser).mint(user1.address, ethers.parseUnits("500", 6));

      // Burn
      await lusd.connect(burner).burn(user1.address, ethers.parseUnits("300", 6));

      expect(await lusd.balanceOf(user1.address)).to.equal(ethers.parseUnits("700", 6));
    });
  });
});
