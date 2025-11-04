import { expect } from "chai";
import hre from "hardhat";
import type { SignerWithAddress } from "@nomicfoundation/hardhat-ethers/signers.js";

const { ethers } = hre;
import type {
  RWAGovernance,
  RWAAssetFactory,
  RWACompliance,
  RWAValuation,
  FractionalRWAToken,
} from "../../typechain-types/index.js";

describe("RWAGovernance", function () {
  let governance: RWAGovernance;
  let factory: RWAAssetFactory;
  let compliance: RWACompliance;
  let valuation: RWAValuation;
  let token: FractionalRWAToken;

  let admin: SignerWithAddress;
  let issuer: SignerWithAddress;
  let holder1: SignerWithAddress;
  let holder2: SignerWithAddress;
  let holder3: SignerWithAddress;

  const GOVERNANCE_ADMIN_ROLE = ethers.keccak256(
    ethers.toUtf8Bytes("GOVERNANCE_ADMIN_ROLE")
  );
  const COMPLIANCE_ROLE = ethers.keccak256(ethers.toUtf8Bytes("COMPLIANCE_ROLE"));
  const VERIFIER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("VERIFIER_ROLE"));
  const MINTER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("MINTER_ROLE"));

  let assetId: number;
  const TOTAL_SUPPLY = ethers.parseEther("1000000"); // 1M tokens

  beforeEach(async function () {
    [admin, issuer, holder1, holder2, holder3] = await ethers.getSigners();

    // Deploy dependencies
    const ComplianceFactory = await ethers.getContractFactory("RWACompliance");
    compliance = await ComplianceFactory.deploy(admin.address);

    const ValuationFactory = await ethers.getContractFactory("RWAValuation");
    valuation = await ValuationFactory.deploy(admin.address);

    const FactoryContract = await ethers.getContractFactory("RWAAssetFactory");
    factory = await FactoryContract.deploy(
      admin.address,
      await compliance.getAddress(),
      await valuation.getAddress()
    );

    // Deploy RWAGovernance
    const GovernanceFactory = await ethers.getContractFactory("RWAGovernance");
    governance = await GovernanceFactory.deploy(
      admin.address,
      await factory.getAddress()
    );

    // Setup roles and verification
    await compliance.grantRole(VERIFIER_ROLE, admin.address);
    await factory.grantRole(COMPLIANCE_ROLE, admin.address);

    const oneYear = 365 * 24 * 60 * 60;
    for (const user of [issuer, holder1, holder2, holder3]) {
      await compliance.verifyInvestor(
        user.address,
        1,
        "US",
        oneYear,
        ethers.keccak256(ethers.toUtf8Bytes(`${user.address}-kyc`))
      );
    }

    // Create and activate asset
    await factory
      .connect(issuer)
      .createAsset("Gov Asset", "GOVAST", 0, ethers.parseEther("1000000"), 0, "QmXYZ", "QmABC");

    assetId = 1;

    await factory.connect(issuer).fractionalizeAsset(assetId, TOTAL_SUPPLY);
    await factory.activateAsset(assetId);

    // Get and distribute tokens
    const tokenAddress = await factory.getFractionalToken(assetId);
    token = await ethers.getContractAt("FractionalRWAToken", tokenAddress);

    await token.connect(issuer).grantRole(MINTER_ROLE, issuer.address);
    await token.connect(issuer).mint(issuer.address, TOTAL_SUPPLY);

    // Distribute: 40% holder1, 30% holder2, 20% holder3, 10% issuer
    await token.connect(issuer).transfer(holder1.address, TOTAL_SUPPLY * 40n / 100n);
    await token.connect(issuer).transfer(holder2.address, TOTAL_SUPPLY * 30n / 100n);
    await token.connect(issuer).transfer(holder3.address, TOTAL_SUPPLY * 20n / 100n);
    // Issuer keeps 10%
  });

  describe("Deployment", function () {
    it("Should set the correct admin", async function () {
      expect(
        await governance.hasRole(await governance.DEFAULT_ADMIN_ROLE(), admin.address)
      ).to.be.true;
    });

    it("Should set the correct factory", async function () {
      expect(await governance.assetFactory()).to.equal(await factory.getAddress());
    });

    it("Should initialize counters to zero", async function () {
      expect(await governance.totalProposals()).to.equal(0);
      expect(await governance.totalVotes()).to.equal(0);
    });

    it("Should revert with invalid admin", async function () {
      const GovernanceFactory = await ethers.getContractFactory("RWAGovernance");
      await expect(
        GovernanceFactory.deploy(ethers.ZeroAddress, await factory.getAddress())
      ).to.be.revertedWith("Invalid admin");
    });

    it("Should revert with invalid factory", async function () {
      const GovernanceFactory = await ethers.getContractFactory("RWAGovernance");
      await expect(
        GovernanceFactory.deploy(admin.address, ethers.ZeroAddress)
      ).to.be.revertedWith("Invalid factory");
    });
  });

  describe("Proposal Creation", function () {
    it("Should create proposal successfully", async function () {
      const tx = await governance
        .connect(holder1)
        .createProposal(
          assetId,
          0, // ProposalType.ParameterChange
          "Change voting period to 14 days",
          "0x"
        );

      await expect(tx)
        .to.emit(governance, "ProposalCreated")
        .withArgs(1, assetId, holder1.address, 0, "Change voting period to 14 days");

      expect(await governance.totalProposals()).to.equal(1);
    });

    it("Should store proposal data correctly", async function () {
      await governance
        .connect(holder1)
        .createProposal(assetId, 0, "Test proposal", "0x1234");

      const proposal = await governance.getProposal(1);

      expect(proposal.proposalId).to.equal(1);
      expect(proposal.assetId).to.equal(assetId);
      expect(proposal.proposer).to.equal(holder1.address);
      expect(proposal.proposalType).to.equal(0);
      expect(proposal.status).to.equal(1); // ProposalStatus.Active
      expect(proposal.description).to.equal("Test proposal");
      expect(proposal.executionData).to.equal("0x1234");
      expect(proposal.totalVotingPower).to.equal(TOTAL_SUPPLY);
    });

    it("Should support all proposal types", async function () {
      // ParameterChange=0, AssetSale=1, YieldStrategy=2, ValuationUpdate=3, EmergencyAction=4
      for (let proposalType = 0; proposalType <= 4; proposalType++) {
        await governance
          .connect(holder1)
          .createProposal(assetId, proposalType, `Proposal type ${proposalType}`, "0x");

        const proposal = await governance.getProposal(proposalType + 1);
        expect(proposal.proposalType).to.equal(proposalType);
      }
    });

    it("Should track asset proposals", async function () {
      await governance.connect(holder1).createProposal(assetId, 0, "Proposal 1", "0x");
      await governance.connect(holder1).createProposal(assetId, 1, "Proposal 2", "0x");

      const assetProposals = await governance.getAssetProposals(assetId);
      expect(assetProposals.length).to.equal(2);
      expect(assetProposals[0]).to.equal(1);
      expect(assetProposals[1]).to.equal(2);
    });

    it("Should revert if asset not active", async function () {
      await factory.updateAssetStatus(assetId, 2); // Suspend

      await expect(
        governance.connect(holder1).createProposal(assetId, 0, "Test", "0x")
      ).to.be.revertedWith("Asset not active");
    });

    it("Should revert without description", async function () {
      await expect(
        governance.connect(holder1).createProposal(assetId, 0, "", "0x")
      ).to.be.revertedWith("Description required");
    });

    it("Should revert if insufficient tokens", async function () {
      // holder3 has 20% which is above default 1% threshold, but let's use someone with 0%
      const [, , , , , , noTokensUser] = await ethers.getSigners();

      await expect(
        governance.connect(noTokensUser).createProposal(assetId, 0, "Test", "0x")
      ).to.be.revertedWith("Insufficient tokens for proposal");
    });
  });

  describe("Voting", function () {
    let proposalId: number;

    beforeEach(async function () {
      await governance
        .connect(holder1)
        .createProposal(assetId, 0, "Test proposal", "0x");
      proposalId = 1;
    });

    it("Should cast vote successfully", async function () {
      const tx = await governance.connect(holder2).castVote(proposalId, 1); // VoteChoice.For

      const holder2Balance = await token.balanceOf(holder2.address);
      await expect(tx)
        .to.emit(governance, "VoteCast")
        .withArgs(holder2.address, proposalId, 1, holder2Balance);

      expect(await governance.totalVotes()).to.equal(1);
    });

    it("Should record vote receipt", async function () {
      await governance.connect(holder2).castVote(proposalId, 1);

      const receipt = await governance.getVoteReceipt(proposalId, holder2.address);

      expect(receipt.hasVoted).to.be.true;
      expect(receipt.choice).to.equal(1); // For
      expect(receipt.votes).to.equal(TOTAL_SUPPLY * 30n / 100n); // holder2 has 30%
    });

    it("Should update proposal vote counts", async function () {
      await governance.connect(holder1).castVote(proposalId, 1); // For (40%)
      await governance.connect(holder2).castVote(proposalId, 0); // Against (30%)
      await governance.connect(holder3).castVote(proposalId, 2); // Abstain (20%)

      const proposal = await governance.getProposal(proposalId);

      expect(proposal.forVotes).to.equal(TOTAL_SUPPLY * 40n / 100n);
      expect(proposal.againstVotes).to.equal(TOTAL_SUPPLY * 30n / 100n);
      expect(proposal.abstainVotes).to.equal(TOTAL_SUPPLY * 20n / 100n);
    });

    it("Should revert if already voted", async function () {
      await governance.connect(holder1).castVote(proposalId, 1);

      await expect(
        governance.connect(holder1).castVote(proposalId, 1)
      ).to.be.revertedWith("Already voted");
    });

    it("Should revert if no voting power", async function () {
      const [, , , , , , noTokensUser] = await ethers.getSigners();

      await expect(
        governance.connect(noTokensUser).castVote(proposalId, 1)
      ).to.be.revertedWith("No voting power");
    });

    it("Should revert if voting ended", async function () {
      // Fast forward past voting period (7 days)
      await ethers.provider.send("evm_increaseTime", [8 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      await expect(
        governance.connect(holder1).castVote(proposalId, 1)
      ).to.be.revertedWith("Voting ended");
    });
  });

  describe("Vote Delegation", function () {
    it("Should delegate vote successfully", async function () {
      const tx = await governance.connect(holder1).delegateVote(assetId, holder2.address);

      await expect(tx)
        .to.emit(governance, "VoteDelegated")
        .withArgs(holder1.address, holder2.address, assetId);
    });

    it("Should reduce delegator's voting power to zero", async function () {
      await governance.connect(holder1).delegateVote(assetId, holder2.address);

      // Create proposal and try to vote
      await governance.connect(holder2).createProposal(assetId, 0, "Test", "0x");

      await expect(
        governance.connect(holder1).castVote(1, 1)
      ).to.be.revertedWith("No voting power");
    });

    it("Should undelegate vote", async function () {
      await governance.connect(holder1).delegateVote(assetId, holder2.address);

      const tx = await governance.connect(holder1).undelegateVote(assetId);

      await expect(tx)
        .to.emit(governance, "VoteDelegated")
        .withArgs(holder1.address, ethers.ZeroAddress, assetId);
    });

    it("Should restore voting power after undelegation", async function () {
      await governance.connect(holder1).delegateVote(assetId, holder2.address);
      await governance.connect(holder1).undelegateVote(assetId);

      // Should be able to vote again
      await governance.connect(holder2).createProposal(assetId, 0, "Test", "0x");

      await expect(governance.connect(holder1).castVote(1, 1)).to.not.be.reverted;
    });

    it("Should revert delegating to zero address", async function () {
      await expect(
        governance.connect(holder1).delegateVote(assetId, ethers.ZeroAddress)
      ).to.be.revertedWith("Invalid delegate");
    });

    it("Should revert delegating to self", async function () {
      await expect(
        governance.connect(holder1).delegateVote(assetId, holder1.address)
      ).to.be.revertedWith("Cannot delegate to self");
    });
  });

  describe("Proposal Execution", function () {
    let proposalId: number;

    beforeEach(async function () {
      await governance
        .connect(holder1)
        .createProposal(assetId, 0, "Test proposal", "0x");
      proposalId = 1;
    });

    it("Should execute passed proposal", async function () {
      // Vote For (70% approval)
      await governance.connect(holder1).castVote(proposalId, 1); // 40%
      await governance.connect(holder2).castVote(proposalId, 1); // 30%

      // Fast forward past voting + timelock period (7 + 2 days)
      await ethers.provider.send("evm_increaseTime", [10 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      const tx = await governance.executeProposal(proposalId);

      await expect(tx).to.emit(governance, "ProposalExecuted").withArgs(proposalId, true);

      const proposal = await governance.getProposal(proposalId);
      expect(proposal.status).to.equal(4); // ProposalStatus.Executed
      expect(proposal.executed).to.be.true;
    });

    it("Should defeat proposal if quorum not reached", async function () {
      // Only 10% votes (issuer), below 10% quorum... wait, 10% = quorum
      // Let's make sure we're below by not voting at all
      // Actually default quorum is 10%, so we need less than 10%

      // Fast forward
      await ethers.provider.send("evm_increaseTime", [10 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      await governance.executeProposal(proposalId);

      const proposal = await governance.getProposal(proposalId);
      expect(proposal.status).to.equal(2); // ProposalStatus.Defeated
    });

    it("Should defeat proposal if approval threshold not met", async function () {
      // Vote 40% For, 50% Against
      await governance.connect(holder1).castVote(proposalId, 1); // For 40%
      await governance.connect(holder2).castVote(proposalId, 0); // Against 30%
      await governance.connect(holder3).castVote(proposalId, 0); // Against 20%

      // Fast forward
      await ethers.provider.send("evm_increaseTime", [10 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      await governance.executeProposal(proposalId);

      const proposal = await governance.getProposal(proposalId);
      expect(proposal.status).to.equal(2); // Defeated (40% < 50% approval)
    });

    it("Should revert if voting not ended", async function () {
      await governance.connect(holder1).castVote(proposalId, 1);

      await expect(governance.executeProposal(proposalId)).to.be.revertedWith(
        "Voting not ended"
      );
    });

    it("Should revert if timelock not passed", async function () {
      await governance.connect(holder1).castVote(proposalId, 1);

      // Fast forward only past voting (7 days), not timelock (+ 2 days)
      await ethers.provider.send("evm_increaseTime", [7 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      await expect(governance.executeProposal(proposalId)).to.be.revertedWith(
        "Timelock not passed"
      );
    });

    it("Should check if proposal has passed", async function () {
      await governance.connect(holder1).castVote(proposalId, 1); // 40%
      await governance.connect(holder2).castVote(proposalId, 1); // 30%

      // Before voting ends
      expect(await governance.hasProposalPassed(proposalId)).to.be.false;

      // After voting ends
      await ethers.provider.send("evm_increaseTime", [8 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      expect(await governance.hasProposalPassed(proposalId)).to.be.true;
    });
  });

  describe("Proposal Cancellation", function () {
    let proposalId: number;

    beforeEach(async function () {
      await governance
        .connect(holder1)
        .createProposal(assetId, 0, "Test proposal", "0x");
      proposalId = 1;
    });

    it("Should allow proposer to cancel", async function () {
      const tx = await governance.connect(holder1).cancelProposal(proposalId);

      await expect(tx).to.emit(governance, "ProposalCancelled").withArgs(proposalId);

      const proposal = await governance.getProposal(proposalId);
      expect(proposal.status).to.equal(5); // ProposalStatus.Cancelled
    });

    it("Should allow admin to cancel", async function () {
      await expect(governance.connect(admin).cancelProposal(proposalId)).to.not.be
        .reverted;
    });

    it("Should revert if not authorized", async function () {
      await expect(
        governance.connect(holder2).cancelProposal(proposalId)
      ).to.be.revertedWith("Not authorized");
    });
  });

  describe("Proposal Veto", function () {
    let proposalId: number;

    beforeEach(async function () {
      await governance
        .connect(holder1)
        .createProposal(assetId, 0, "Test proposal", "0x");
      proposalId = 1;
    });

    it("Should allow issuer to veto", async function () {
      const tx = await governance.connect(issuer).vetoProposal(proposalId);

      await expect(tx)
        .to.emit(governance, "ProposalVetoed")
        .withArgs(proposalId, issuer.address);

      const proposal = await governance.getProposal(proposalId);
      expect(proposal.vetoed).to.be.true;
      expect(proposal.status).to.equal(5); // Cancelled
    });

    it("Should revert if not issuer", async function () {
      await expect(
        governance.connect(holder1).vetoProposal(proposalId)
      ).to.be.revertedWith("Not asset issuer");
    });

    it("Should revert executing vetoed proposal", async function () {
      await governance.connect(issuer).vetoProposal(proposalId);

      await ethers.provider.send("evm_increaseTime", [10 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine", []);

      // After veto, status is Cancelled, so execution reverts with "Proposal not active"
      await expect(governance.executeProposal(proposalId)).to.be.revertedWith(
        "Proposal not active"
      );
    });
  });

  describe("Governance Parameters", function () {
    it("Should use default parameters", async function () {
      const params = await governance.getGovernanceParams(assetId);

      expect(params.proposalThreshold).to.equal(100); // 1%
      expect(params.quorumThreshold).to.equal(1000); // 10%
      expect(params.approvalThreshold).to.equal(5000); // 50%
      expect(params.votingPeriod).to.equal(7 * 24 * 60 * 60);
      expect(params.timelockPeriod).to.equal(2 * 24 * 60 * 60);
      expect(params.issuerVetoEnabled).to.be.true;
    });

    it("Should set custom governance parameters", async function () {
      const customParams = {
        proposalThreshold: 500, // 5%
        quorumThreshold: 2000, // 20%
        approvalThreshold: 6000, // 60%
        votingPeriod: 14 * 24 * 60 * 60, // 14 days
        timelockPeriod: 3 * 24 * 60 * 60, // 3 days
        issuerVetoEnabled: false,
      };

      await governance.setGovernanceParams(assetId, customParams);

      const params = await governance.getGovernanceParams(assetId);

      expect(params.proposalThreshold).to.equal(500);
      expect(params.quorumThreshold).to.equal(2000);
      expect(params.approvalThreshold).to.equal(6000);
      expect(params.votingPeriod).to.equal(14 * 24 * 60 * 60);
      expect(params.timelockPeriod).to.equal(3 * 24 * 60 * 60);
      expect(params.issuerVetoEnabled).to.be.false;
    });

    it("Should revert with high proposal threshold", async function () {
      const params = {
        proposalThreshold: 5001, // > 50%
        quorumThreshold: 1000,
        approvalThreshold: 5000,
        votingPeriod: 7 * 24 * 60 * 60,
        timelockPeriod: 2 * 24 * 60 * 60,
        issuerVetoEnabled: true,
      };

      await expect(
        governance.setGovernanceParams(assetId, params)
      ).to.be.revertedWith("Proposal threshold too high");
    });

    it("Should revert with invalid voting period", async function () {
      const params = {
        proposalThreshold: 100,
        quorumThreshold: 1000,
        approvalThreshold: 5000,
        votingPeriod: 31 * 24 * 60 * 60, // > 30 days
        timelockPeriod: 2 * 24 * 60 * 60,
        issuerVetoEnabled: true,
      };

      await expect(
        governance.setGovernanceParams(assetId, params)
      ).to.be.revertedWith("Invalid voting period");
    });

    it("Should revert if not governance admin", async function () {
      const params = {
        proposalThreshold: 100,
        quorumThreshold: 1000,
        approvalThreshold: 5000,
        votingPeriod: 7 * 24 * 60 * 60,
        timelockPeriod: 2 * 24 * 60 * 60,
        issuerVetoEnabled: true,
      };

      await expect(
        governance.connect(holder1).setGovernanceParams(assetId, params)
      ).to.be.revertedWithCustomError(governance, "AccessControlUnauthorizedAccount");
    });
  });

  describe("View Functions", function () {
    beforeEach(async function () {
      // Create multiple proposals
      await governance.connect(holder1).createProposal(assetId, 0, "Proposal 1", "0x");
      await governance.connect(holder1).createProposal(assetId, 1, "Proposal 2", "0x");
      await governance.connect(holder1).createProposal(assetId, 2, "Proposal 3", "0x");

      // Cancel one
      await governance.connect(holder1).cancelProposal(2);
    });

    it("Should get all asset proposals", async function () {
      const proposals = await governance.getAssetProposals(assetId);

      expect(proposals.length).to.equal(3);
    });

    it("Should get only active proposals", async function () {
      const activeProposals = await governance.getActiveProposals(assetId);

      expect(activeProposals.length).to.equal(2); // 1 and 2, not 3 (cancelled)
      expect(activeProposals[0]).to.equal(1);
      expect(activeProposals[1]).to.equal(3);
    });

    it("Should get proposal details", async function () {
      const proposal = await governance.getProposal(1);

      expect(proposal.proposalId).to.equal(1);
      expect(proposal.description).to.equal("Proposal 1");
    });
  });
});
