// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "../core/CollateralVaultL1.sol";
import "../core/LoyaltyUSDL1.sol";

/**
 * @title L1Gateway
 * @notice Layer 1 bridge gateway for cross-chain operations
 * @dev Integrates with Arbitrum's native bridge for L1<->L2 communication
 *
 * Architecture:
 * - Uses Arbitrum's native Inbox/Outbox for message passing
 * - Handles collateral deposits to L2
 * - Processes L2 state updates
 * - Manages emergency exits
 *
 * Integration with Arbitrum:
 * - Inbox: 0x4Dbd4fc535Ac27206064B68FfCf827b0A60BAB3f (Mainnet)
 * - Uses retryable tickets for guaranteed L2 execution
 */

// Arbitrum Inbox interface
interface IInbox {
    function createRetryableTicket(
        address to,
        uint256 l2CallValue,
        uint256 maxSubmissionCost,
        address excessFeeRefundAddress,
        address callValueRefundAddress,
        uint256 gasLimit,
        uint256 maxFeePerGas,
        bytes calldata data
    ) external payable returns (uint256);
}

// Arbitrum Outbox interface
interface IOutbox {
    function l2ToL1Sender() external view returns (address);
}

contract L1Gateway is ReentrancyGuard, Ownable {
    using SafeERC20 for IERC20;

    // ============ State Variables ============

    // Core L1 contracts
    CollateralVaultL1 public immutable collateralVault;
    LoyaltyUSDL1 public immutable loyaltyUSD;
    IERC20 public immutable collateralToken;

    // Arbitrum bridge contracts
    IInbox public inbox;
    IOutbox public outbox;

    // L2 counterpart address
    address public l2Gateway;

    // Gas parameters for L2 calls
    uint256 public maxSubmissionCost = 0.01 ether;
    uint256 public maxGas = 300000;
    uint256 public gasPriceBid = 1 gwei;

    // Deposit tracking
    struct Deposit {
        address user;
        uint256 amount;
        uint256 timestamp;
        bytes32 l2TxHash;
        bool processed;
    }

    mapping(uint256 => Deposit) public deposits;
    uint256 public depositNonce;

    // Withdrawal tracking
    struct Withdrawal {
        address user;
        uint256 amount;
        uint256 timestamp;
        bool executed;
    }

    mapping(bytes32 => Withdrawal) public withdrawals;

    // Emergency pause
    bool public paused;

    // ============ Events ============

    event DepositInitiated(
        uint256 indexed depositId,
        address indexed user,
        uint256 amount,
        uint256 indexed ticketId
    );

    event DepositFinalized(
        uint256 indexed depositId,
        address indexed user,
        uint256 amount
    );

    event WithdrawalInitiated(
        bytes32 indexed withdrawalId,
        address indexed user,
        uint256 amount
    );

    event WithdrawalExecuted(
        bytes32 indexed withdrawalId,
        address indexed user,
        uint256 amount
    );

    event L2GatewayUpdated(address indexed oldGateway, address indexed newGateway);
    event GasParametersUpdated(uint256 maxSubmissionCost, uint256 maxGas, uint256 gasPriceBid);

    // ============ Modifiers ============

    modifier whenNotPaused() {
        require(!paused, "L1Gateway: paused");
        _;
    }

    modifier onlyL2() {
        require(msg.sender == address(outbox), "L1Gateway: not from outbox");
        require(outbox.l2ToL1Sender() == l2Gateway, "L1Gateway: not from L2 gateway");
        _;
    }

    // ============ Constructor ============

    constructor(
        address _collateralVault,
        address _loyaltyUSD,
        address _collateralToken,
        address _inbox,
        address _outbox
    ) Ownable(msg.sender) {
        require(_collateralVault != address(0), "Invalid vault");
        require(_loyaltyUSD != address(0), "Invalid LUSD");
        require(_collateralToken != address(0), "Invalid collateral");
        require(_inbox != address(0), "Invalid inbox");
        require(_outbox != address(0), "Invalid outbox");

        collateralVault = CollateralVaultL1(_collateralVault);
        loyaltyUSD = LoyaltyUSDL1(_loyaltyUSD);
        collateralToken = IERC20(_collateralToken);
        inbox = IInbox(_inbox);
        outbox = IOutbox(_outbox);
    }

    // ============ Deposit Functions ============

    /**
     * @notice Deposit collateral to L2
     * @param amount Amount of collateral to deposit
     */
    function depositToL2(uint256 amount)
        external
        payable
        whenNotPaused
        nonReentrant
        returns (uint256 depositId, uint256 ticketId)
    {
        require(amount > 0, "Amount must be > 0");
        require(l2Gateway != address(0), "L2 gateway not set");
        require(msg.value >= maxSubmissionCost, "Insufficient ETH for submission");

        // Increment nonce
        depositId = depositNonce++;

        // Transfer collateral from user to this contract
        collateralToken.safeTransferFrom(msg.sender, address(this), amount);

        // Approve vault to spend (using forceApprove for OZ 5.x)
        collateralToken.forceApprove(address(collateralVault), amount);

        // Lock in vault
        bytes32 l2TxHash = keccak256(abi.encodePacked(depositId, msg.sender, amount, block.timestamp));
        collateralVault.lockCollateral(msg.sender, amount, l2TxHash);

        // Create retryable ticket to L2
        bytes memory data = abi.encodeWithSignature(
            "finalizeDeposit(uint256,address,uint256)",
            depositId,
            msg.sender,
            amount
        );

        ticketId = inbox.createRetryableTicket{value: msg.value}(
            l2Gateway,
            0, // l2CallValue
            maxSubmissionCost,
            msg.sender, // excessFeeRefundAddress
            msg.sender, // callValueRefundAddress
            maxGas,
            gasPriceBid,
            data
        );

        // Store deposit info
        deposits[depositId] = Deposit({
            user: msg.sender,
            amount: amount,
            timestamp: block.timestamp,
            l2TxHash: l2TxHash,
            processed: false
        });

        emit DepositInitiated(depositId, msg.sender, amount, ticketId);
    }

    // ============ Withdrawal Functions ============

    /**
     * @notice Finalize withdrawal from L2
     * @param withdrawalId Withdrawal identifier
     * @param user User address
     * @param amount Amount to withdraw
     * @dev Called by L2 via Arbitrum's outbox
     */
    function finalizeWithdrawal(
        bytes32 withdrawalId,
        address user,
        uint256 amount
    ) external onlyL2 nonReentrant {
        require(user != address(0), "Invalid user");
        require(amount > 0, "Invalid amount");
        require(!withdrawals[withdrawalId].executed, "Already executed");

        // Record withdrawal
        withdrawals[withdrawalId] = Withdrawal({
            user: user,
            amount: amount,
            timestamp: block.timestamp,
            executed: true
        });

        // Unlock from vault
        bytes32 l2TxHash = withdrawalId;
        collateralVault.unlockCollateral(user, amount, l2TxHash);

        emit WithdrawalExecuted(withdrawalId, user, amount);
    }

    // ============ LUSD Minting/Burning (Bridge Operations) ============

    /**
     * @notice Mint LUSD based on L2 state
     * @param user User address
     * @param amount Amount to mint
     * @param l2TxHash L2 transaction hash
     * @dev Called by L2 via Arbitrum's outbox
     */
    function bridgeMint(address user, uint256 amount, bytes32 l2TxHash)
        external
        onlyL2
        nonReentrant
    {
        require(user != address(0), "Invalid user");
        require(amount > 0, "Invalid amount");

        loyaltyUSD.bridgeMint(user, amount, l2TxHash);
    }

    /**
     * @notice Burn LUSD based on L2 state
     * @param user User address
     * @param amount Amount to burn
     * @param l2TxHash L2 transaction hash
     * @dev Called by L2 via Arbitrum's outbox
     */
    function bridgeBurn(address user, uint256 amount, bytes32 l2TxHash)
        external
        onlyL2
        nonReentrant
    {
        require(user != address(0), "Invalid user");
        require(amount > 0, "Invalid amount");

        loyaltyUSD.bridgeBurn(user, amount, l2TxHash);
    }

    // ============ Admin Functions ============

    /**
     * @notice Set L2 gateway address
     * @param _l2Gateway L2 gateway address
     */
    function setL2Gateway(address _l2Gateway) external onlyOwner {
        require(_l2Gateway != address(0), "Invalid gateway");
        address oldGateway = l2Gateway;
        l2Gateway = _l2Gateway;
        emit L2GatewayUpdated(oldGateway, _l2Gateway);
    }

    /**
     * @notice Update gas parameters for L2 calls
     */
    function setGasParameters(
        uint256 _maxSubmissionCost,
        uint256 _maxGas,
        uint256 _gasPriceBid
    ) external onlyOwner {
        maxSubmissionCost = _maxSubmissionCost;
        maxGas = _maxGas;
        gasPriceBid = _gasPriceBid;
        emit GasParametersUpdated(_maxSubmissionCost, _maxGas, _gasPriceBid);
    }

    /**
     * @notice Update Arbitrum Inbox address
     */
    function setInbox(address _inbox) external onlyOwner {
        require(_inbox != address(0), "Invalid inbox");
        inbox = IInbox(_inbox);
    }

    /**
     * @notice Update Arbitrum Outbox address
     */
    function setOutbox(address _outbox) external onlyOwner {
        require(_outbox != address(0), "Invalid outbox");
        outbox = IOutbox(_outbox);
    }

    /**
     * @notice Pause deposits
     */
    function pause() external onlyOwner {
        paused = true;
    }

    /**
     * @notice Resume deposits
     */
    function unpause() external onlyOwner {
        paused = false;
    }

    // ============ View Functions ============

    /**
     * @notice Get deposit info
     */
    function getDeposit(uint256 depositId)
        external
        view
        returns (
            address user,
            uint256 amount,
            uint256 timestamp,
            bytes32 l2TxHash,
            bool processed
        )
    {
        Deposit memory deposit = deposits[depositId];
        return (
            deposit.user,
            deposit.amount,
            deposit.timestamp,
            deposit.l2TxHash,
            deposit.processed
        );
    }

    /**
     * @notice Get withdrawal info
     */
    function getWithdrawal(bytes32 withdrawalId)
        external
        view
        returns (address user, uint256 amount, uint256 timestamp, bool executed)
    {
        Withdrawal memory withdrawal = withdrawals[withdrawalId];
        return (
            withdrawal.user,
            withdrawal.amount,
            withdrawal.timestamp,
            withdrawal.executed
        );
    }

    /**
     * @notice Calculate required ETH for deposit
     */
    function calculateRequiredEth() external view returns (uint256) {
        return maxSubmissionCost + (maxGas * gasPriceBid);
    }
}
