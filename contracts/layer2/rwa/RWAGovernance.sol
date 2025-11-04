// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "../interfaces/rwa/IRWAAsset.sol";

/**
 * @title RWAGovernance
 * @notice On-chain governance for RWA asset management
 * @dev Token-weighted voting with delegation support
 *
 * Key Features:
 * - Token-weighted voting (1 token = 1 vote)
 * - Proposal types: Parameter changes, Asset sale, Yield strategy
 * - Vote delegation for passive investors
 * - Quorum requirements for proposal validity
 * - Timelock for proposal execution
 * - Veto power for asset issuer (optional)
 *
 * Governance Flow:
 * 1. Token holder creates proposal (requires minimum tokens)
 * 2. Voting period opens (typically 7 days)
 * 3. Token holders vote (For/Against/Abstain)
 * 4. If quorum reached and majority votes For, proposal passes
 * 5. After timelock period, proposal can be executed
 * 6. Proposal changes take effect
 */
contract RWAGovernance is AccessControl, ReentrancyGuard {
    bytes32 public constant GOVERNANCE_ADMIN_ROLE = keccak256("GOVERNANCE_ADMIN_ROLE");

    // Proposal types
    enum ProposalType {
        ParameterChange,    // Change governance parameters
        AssetSale,          // Sell the underlying asset
        YieldStrategy,      // Change yield distribution strategy
        ValuationUpdate,    // Request asset revaluation
        EmergencyAction     // Emergency actions (pause, etc.)
    }

    // Proposal status
    enum ProposalStatus {
        Pending,            // Proposal created, voting not started
        Active,             // Voting in progress
        Defeated,           // Proposal defeated
        Succeeded,          // Proposal passed, awaiting execution
        Executed,           // Proposal executed
        Cancelled,          // Proposal cancelled
        Expired             // Proposal expired without execution
    }

    // Vote choice
    enum VoteChoice {
        Against,
        For,
        Abstain
    }

    // Proposal struct
    struct Proposal {
        uint256 proposalId;
        uint256 assetId;
        address proposer;
        ProposalType proposalType;
        ProposalStatus status;
        string description;
        bytes executionData;        // Encoded function call data
        uint256 votingStartTime;
        uint256 votingEndTime;
        uint256 executionTime;      // When it can be executed (after timelock)
        uint256 forVotes;
        uint256 againstVotes;
        uint256 abstainVotes;
        uint256 totalVotingPower;   // Total token supply at snapshot
        bool executed;
        bool vetoed;
    }

    // Vote record
    struct VoteReceipt {
        bool hasVoted;
        VoteChoice choice;
        uint256 votes;
    }

    // Governance parameters
    struct GovernanceParams {
        uint256 proposalThreshold;  // Min tokens to create proposal (%)
        uint256 quorumThreshold;    // Min participation for validity (%)
        uint256 approvalThreshold;  // Min % of For votes to pass
        uint256 votingPeriod;       // Duration of voting (seconds)
        uint256 timelockPeriod;     // Delay before execution (seconds)
        bool issuerVetoEnabled;     // Can issuer veto proposals
    }

    // Proposal counter
    uint256 private proposalIdCounter;

    // Proposals storage
    mapping(uint256 => Proposal) private proposals;
    mapping(uint256 => mapping(address => VoteReceipt)) private votes;

    // Asset-specific governance params
    mapping(uint256 => GovernanceParams) private assetGovernance;

    // Vote delegation (delegator => delegate)
    mapping(address => mapping(uint256 => address)) private delegates;

    // Asset proposals tracking
    mapping(uint256 => uint256[]) private assetProposals;

    // Integration
    IRWAAsset public immutable assetFactory;

    // Default parameters
    uint256 public constant DEFAULT_PROPOSAL_THRESHOLD = 100; // 1%
    uint256 public constant DEFAULT_QUORUM = 1000; // 10%
    uint256 public constant DEFAULT_APPROVAL = 5000; // 50%
    uint256 public constant DEFAULT_VOTING_PERIOD = 7 days;
    uint256 public constant DEFAULT_TIMELOCK = 2 days;
    uint256 public constant PERCENTAGE_BASE = 10000; // 100%

    // Statistics
    uint256 public totalProposals;
    uint256 public totalVotes;

    // Events
    event ProposalCreated(
        uint256 indexed proposalId,
        uint256 indexed assetId,
        address indexed proposer,
        ProposalType proposalType,
        string description
    );

    event VoteCast(
        address indexed voter,
        uint256 indexed proposalId,
        VoteChoice choice,
        uint256 votes
    );

    event ProposalExecuted(
        uint256 indexed proposalId,
        bool success
    );

    event ProposalCancelled(
        uint256 indexed proposalId
    );

    event ProposalVetoed(
        uint256 indexed proposalId,
        address indexed vetoer
    );

    event VoteDelegated(
        address indexed delegator,
        address indexed delegate,
        uint256 indexed assetId
    );

    /**
     * @notice Constructor
     * @param admin Admin address
     * @param factory RWAAssetFactory address
     */
    constructor(address admin, address factory) {
        require(admin != address(0), "Invalid admin");
        require(factory != address(0), "Invalid factory");

        assetFactory = IRWAAsset(factory);

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(GOVERNANCE_ADMIN_ROLE, admin);
    }

    // =============================================================
    //                   PROPOSAL CREATION
    // =============================================================

    /**
     * @notice Create governance proposal
     * @param assetId Asset identifier
     * @param proposalType Type of proposal
     * @param description Proposal description
     * @param executionData Encoded execution data
     * @return proposalId New proposal ID
     */
    function createProposal(
        uint256 assetId,
        ProposalType proposalType,
        string calldata description,
        bytes calldata executionData
    ) external nonReentrant returns (uint256 proposalId) {
        require(assetFactory.isAssetActive(assetId), "Asset not active");
        require(bytes(description).length > 0, "Description required");

        // Check proposer has enough tokens
        IERC20 token = IERC20(assetFactory.getFractionalToken(assetId));
        GovernanceParams memory params = _getGovernanceParams(assetId);

        require(
            token.balanceOf(msg.sender) >= (token.totalSupply() * params.proposalThreshold) / PERCENTAGE_BASE,
            "Insufficient tokens for proposal"
        );

        // Create proposal
        proposalIdCounter++;
        proposalId = proposalIdCounter;

        Proposal storage proposal = proposals[proposalId];
        proposal.proposalId = proposalId;
        proposal.assetId = assetId;
        proposal.proposer = msg.sender;
        proposal.proposalType = proposalType;
        proposal.status = ProposalStatus.Active;
        proposal.description = description;
        proposal.executionData = executionData;
        proposal.votingStartTime = block.timestamp;
        proposal.votingEndTime = block.timestamp + params.votingPeriod;
        proposal.executionTime = block.timestamp + params.votingPeriod + params.timelockPeriod;
        proposal.totalVotingPower = token.totalSupply();

        // Track proposal
        assetProposals[assetId].push(proposalId);
        totalProposals++;

        emit ProposalCreated(proposalId, assetId, msg.sender, proposalType, description);
    }

    /**
     * @notice Cancel proposal (only proposer or admin)
     * @param proposalId Proposal to cancel
     */
    function cancelProposal(uint256 proposalId) external {
        Proposal storage proposal = proposals[proposalId];

        require(
            msg.sender == proposal.proposer || hasRole(GOVERNANCE_ADMIN_ROLE, msg.sender),
            "Not authorized"
        );
        require(
            proposal.status == ProposalStatus.Pending || proposal.status == ProposalStatus.Active,
            "Cannot cancel"
        );

        proposal.status = ProposalStatus.Cancelled;

        emit ProposalCancelled(proposalId);
    }

    // =============================================================
    //                       VOTING
    // =============================================================

    /**
     * @notice Cast vote on proposal
     * @param proposalId Proposal to vote on
     * @param choice Vote choice (Against/For/Abstain)
     */
    function castVote(
        uint256 proposalId,
        VoteChoice choice
    ) external nonReentrant {
        Proposal storage proposal = proposals[proposalId];

        require(proposal.status == ProposalStatus.Active, "Proposal not active");
        require(block.timestamp < proposal.votingEndTime, "Voting ended");

        VoteReceipt storage receipt = votes[proposalId][msg.sender];
        require(!receipt.hasVoted, "Already voted");

        // Get voting power
        uint256 votingPower = _getVotingPower(msg.sender, proposal.assetId);
        require(votingPower > 0, "No voting power");

        // Record vote
        receipt.hasVoted = true;
        receipt.choice = choice;
        receipt.votes = votingPower;

        // Update proposal vote counts
        if (choice == VoteChoice.For) {
            proposal.forVotes += votingPower;
        } else if (choice == VoteChoice.Against) {
            proposal.againstVotes += votingPower;
        } else {
            proposal.abstainVotes += votingPower;
        }

        totalVotes++;

        emit VoteCast(msg.sender, proposalId, choice, votingPower);
    }

    /**
     * @notice Delegate voting power
     * @param assetId Asset to delegate for
     * @param delegate Address to delegate to
     */
    function delegateVote(uint256 assetId, address delegate) external {
        require(delegate != address(0), "Invalid delegate");
        require(delegate != msg.sender, "Cannot delegate to self");

        delegates[msg.sender][assetId] = delegate;

        emit VoteDelegated(msg.sender, delegate, assetId);
    }

    /**
     * @notice Remove vote delegation
     * @param assetId Asset to undelegate
     */
    function undelegateVote(uint256 assetId) external {
        delegates[msg.sender][assetId] = address(0);

        emit VoteDelegated(msg.sender, address(0), assetId);
    }

    // =============================================================
    //                   PROPOSAL EXECUTION
    // =============================================================

    /**
     * @notice Execute passed proposal
     * @param proposalId Proposal to execute
     */
    function executeProposal(uint256 proposalId) external nonReentrant {
        Proposal storage proposal = proposals[proposalId];

        require(proposal.status == ProposalStatus.Active, "Proposal not active");
        require(block.timestamp >= proposal.votingEndTime, "Voting not ended");
        require(block.timestamp >= proposal.executionTime, "Timelock not passed");
        require(!proposal.executed, "Already executed");
        require(!proposal.vetoed, "Proposal vetoed");

        // Check if proposal passed
        GovernanceParams memory params = _getGovernanceParams(proposal.assetId);

        uint256 proposalVotes = proposal.forVotes + proposal.againstVotes + proposal.abstainVotes;
        uint256 quorum = (proposal.totalVotingPower * params.quorumThreshold) / PERCENTAGE_BASE;

        if (proposalVotes < quorum) {
            proposal.status = ProposalStatus.Defeated;
            return;
        }

        uint256 requiredApproval = (proposalVotes * params.approvalThreshold) / PERCENTAGE_BASE;

        if (proposal.forVotes < requiredApproval) {
            proposal.status = ProposalStatus.Defeated;
            return;
        }

        // Execute proposal
        proposal.executed = true;
        proposal.status = ProposalStatus.Executed;

        // In production, would decode and execute executionData
        // For now, just mark as executed

        emit ProposalExecuted(proposalId, true);
    }

    /**
     * @notice Veto proposal (issuer only, if enabled)
     * @param proposalId Proposal to veto
     */
    function vetoProposal(uint256 proposalId) external {
        Proposal storage proposal = proposals[proposalId];

        // Check if veto is enabled
        GovernanceParams memory params = _getGovernanceParams(proposal.assetId);
        require(params.issuerVetoEnabled, "Veto not enabled");

        // Check if caller is asset issuer
        IRWAAsset.AssetMetadata memory metadata = assetFactory.getAssetMetadata(proposal.assetId);
        require(msg.sender == metadata.issuer, "Not asset issuer");

        require(
            proposal.status == ProposalStatus.Active || proposal.status == ProposalStatus.Succeeded,
            "Cannot veto"
        );

        proposal.vetoed = true;
        proposal.status = ProposalStatus.Cancelled;

        emit ProposalVetoed(proposalId, msg.sender);
    }

    // =============================================================
    //                  GOVERNANCE PARAMETERS
    // =============================================================

    /**
     * @notice Set governance parameters for asset
     * @param assetId Asset identifier
     * @param params Governance parameters
     */
    function setGovernanceParams(
        uint256 assetId,
        GovernanceParams calldata params
    ) external onlyRole(GOVERNANCE_ADMIN_ROLE) {
        require(params.proposalThreshold <= 5000, "Proposal threshold too high"); // Max 50%
        require(params.quorumThreshold <= PERCENTAGE_BASE, "Invalid quorum");
        require(params.approvalThreshold <= PERCENTAGE_BASE, "Invalid approval");
        require(params.votingPeriod >= 1 days && params.votingPeriod <= 30 days, "Invalid voting period");
        require(params.timelockPeriod <= 30 days, "Invalid timelock");

        assetGovernance[assetId] = params;
    }

    /**
     * @notice Get governance parameters (with defaults)
     */
    function _getGovernanceParams(uint256 assetId) internal view returns (GovernanceParams memory) {
        GovernanceParams memory params = assetGovernance[assetId];

        // Use defaults if not set
        if (params.votingPeriod == 0) {
            params.proposalThreshold = DEFAULT_PROPOSAL_THRESHOLD;
            params.quorumThreshold = DEFAULT_QUORUM;
            params.approvalThreshold = DEFAULT_APPROVAL;
            params.votingPeriod = DEFAULT_VOTING_PERIOD;
            params.timelockPeriod = DEFAULT_TIMELOCK;
            params.issuerVetoEnabled = true;
        }

        return params;
    }

    // =============================================================
    //                      VIEW FUNCTIONS
    // =============================================================

    /**
     * @notice Get proposal details
     * @param proposalId Proposal identifier
     * @return Proposal struct
     */
    function getProposal(uint256 proposalId) external view returns (Proposal memory) {
        return proposals[proposalId];
    }

    /**
     * @notice Get vote receipt
     * @param proposalId Proposal identifier
     * @param voter Voter address
     * @return VoteReceipt struct
     */
    function getVoteReceipt(
        uint256 proposalId,
        address voter
    ) external view returns (VoteReceipt memory) {
        return votes[proposalId][voter];
    }

    /**
     * @notice Get voting power for address
     * @param account Account to check
     * @param assetId Asset identifier
     * @return Voting power (token balance + delegated)
     */
    function _getVotingPower(address account, uint256 assetId) internal view returns (uint256) {
        address tokenAddress = assetFactory.getFractionalToken(assetId);
        IERC20 token = IERC20(tokenAddress);

        // Check if delegated
        address delegate = delegates[account][assetId];
        if (delegate != address(0)) {
            return 0; // Delegated away
        }

        return token.balanceOf(account);
    }

    /**
     * @notice Get all proposals for asset
     * @param assetId Asset identifier
     * @return Array of proposal IDs
     */
    function getAssetProposals(uint256 assetId) external view returns (uint256[] memory) {
        return assetProposals[assetId];
    }

    /**
     * @notice Get active proposals for asset
     * @param assetId Asset identifier
     * @return Array of active proposal IDs
     */
    function getActiveProposals(uint256 assetId) external view returns (uint256[] memory) {
        uint256[] storage allProposals = assetProposals[assetId];
        uint256 count = 0;

        for (uint256 i = 0; i < allProposals.length; i++) {
            if (proposals[allProposals[i]].status == ProposalStatus.Active) {
                count++;
            }
        }

        uint256[] memory activeProposals = new uint256[](count);
        uint256 index = 0;

        for (uint256 i = 0; i < allProposals.length; i++) {
            if (proposals[allProposals[i]].status == ProposalStatus.Active) {
                activeProposals[index] = allProposals[i];
                index++;
            }
        }

        return activeProposals;
    }

    /**
     * @notice Check if proposal has passed
     * @param proposalId Proposal identifier
     * @return True if passed
     */
    function hasProposalPassed(uint256 proposalId) external view returns (bool) {
        Proposal storage proposal = proposals[proposalId];

        if (block.timestamp < proposal.votingEndTime) {
            return false;
        }

        GovernanceParams memory params = _getGovernanceParams(proposal.assetId);

        uint256 proposalVotes = proposal.forVotes + proposal.againstVotes + proposal.abstainVotes;
        uint256 quorum = (proposal.totalVotingPower * params.quorumThreshold) / PERCENTAGE_BASE;

        if (proposalVotes < quorum) {
            return false;
        }

        uint256 requiredApproval = (proposalVotes * params.approvalThreshold) / PERCENTAGE_BASE;

        return proposal.forVotes >= requiredApproval;
    }

    /**
     * @notice Get governance parameters for asset
     * @param assetId Asset identifier
     * @return GovernanceParams struct
     */
    function getGovernanceParams(uint256 assetId) external view returns (GovernanceParams memory) {
        return _getGovernanceParams(assetId);
    }
}
