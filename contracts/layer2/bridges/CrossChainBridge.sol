// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

/**
 * @title CrossChainBridge
 * @notice Bridge contract for transferring assets between L1 and L2
 * @dev Supports Ethereum mainnet <-> Arbitrum/Optimism bridging
 */
contract CrossChainBridge is ReentrancyGuard, AccessControl, Pausable {
    using ECDSA for bytes32;

    bytes32 public constant RELAYER_ROLE = keccak256("RELAYER_ROLE");
    bytes32 public constant VALIDATOR_ROLE = keccak256("VALIDATOR_ROLE");

    struct BridgeRequest {
        address token;
        address from;
        address to;
        uint256 amount;
        uint256 nonce;
        uint256 targetChainId;
        uint256 timestamp;
        bool processed;
    }

    // Chain IDs
    uint256 public constant ETHEREUM_MAINNET = 1;
    uint256 public constant ARBITRUM_ONE = 42161;
    uint256 public constant OPTIMISM = 10;
    uint256 public constant BASE = 8453;

    uint256 public immutable chainId;
    uint256 public bridgeFee = 0.001 ether; // Bridge fee
    uint256 public minConfirmations = 3; // Required validator signatures

    // Supported tokens on this chain
    mapping(address => bool) public supportedTokens;
    mapping(address => uint256) public tokenLimits; // Max amount per transfer

    // Bridge requests
    mapping(bytes32 => BridgeRequest) public bridgeRequests;
    mapping(address => uint256) public userNonces;

    // Cross-chain token mappings (local token => remote chain => remote token)
    mapping(address => mapping(uint256 => address)) public tokenMappings;

    // Validator signatures for requests
    mapping(bytes32 => mapping(address => bool)) public validatorSignatures;
    mapping(bytes32 => uint256) public signatureCount;

    // Treasury for collected fees
    address public treasury;

    // Events
    event BridgeInitiated(
        bytes32 indexed requestId,
        address indexed from,
        address indexed token,
        uint256 amount,
        uint256 targetChainId
    );

    event BridgeCompleted(
        bytes32 indexed requestId,
        address indexed to,
        address indexed token,
        uint256 amount
    );

    event ValidatorSignature(
        bytes32 indexed requestId,
        address indexed validator
    );

    event TokenMappingSet(
        address localToken,
        uint256 remoteChainId,
        address remoteToken
    );

    event FeeUpdated(uint256 newFee);

    constructor(uint256 _chainId, address _treasury) {
        chainId = _chainId;
        treasury = _treasury;
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
    }

    /**
     * @notice Initiate bridge transfer to another chain
     * @param token Token to bridge
     * @param amount Amount to bridge
     * @param targetChainId Target chain ID
     * @param recipient Recipient address on target chain
     */
    function bridge(
        address token,
        uint256 amount,
        uint256 targetChainId,
        address recipient
    ) external payable nonReentrant whenNotPaused {
        require(msg.value >= bridgeFee, "Insufficient bridge fee");
        require(supportedTokens[token], "Token not supported");
        require(amount > 0 && amount <= tokenLimits[token], "Invalid amount");
        require(targetChainId != chainId, "Same chain transfer");
        require(recipient != address(0), "Invalid recipient");

        // Check if target chain is supported
        require(
            targetChainId == ETHEREUM_MAINNET ||
            targetChainId == ARBITRUM_ONE ||
            targetChainId == OPTIMISM ||
            targetChainId == BASE,
            "Unsupported target chain"
        );

        // Transfer tokens to bridge
        IERC20(token).transferFrom(msg.sender, address(this), amount);

        // Create bridge request
        uint256 nonce = userNonces[msg.sender]++;
        bytes32 requestId = keccak256(
            abi.encodePacked(
                token,
                msg.sender,
                recipient,
                amount,
                nonce,
                targetChainId,
                block.timestamp,
                chainId
            )
        );

        bridgeRequests[requestId] = BridgeRequest({
            token: token,
            from: msg.sender,
            to: recipient,
            amount: amount,
            nonce: nonce,
            targetChainId: targetChainId,
            timestamp: block.timestamp,
            processed: false
        });

        // Transfer fee to treasury
        if (msg.value > 0) {
            (bool success, ) = treasury.call{value: msg.value}("");
            require(success, "Fee transfer failed");
        }

        emit BridgeInitiated(requestId, msg.sender, token, amount, targetChainId);
    }

    /**
     * @notice Complete bridge transfer from another chain
     * @param requestId Bridge request ID from source chain
     * @param token Local token address
     * @param to Recipient address
     * @param amount Amount to transfer
     * @param sourceChainId Source chain ID
     * @param signatures Validator signatures
     */
    function completeBridge(
        bytes32 requestId,
        address token,
        address to,
        uint256 amount,
        uint256 sourceChainId,
        bytes[] memory signatures
    ) external nonReentrant whenNotPaused {
        require(!bridgeRequests[requestId].processed, "Already processed");
        require(supportedTokens[token], "Token not supported");
        require(signatures.length >= minConfirmations, "Insufficient signatures");

        // Verify signatures
        bytes32 messageHash = keccak256(
            abi.encodePacked(requestId, token, to, amount, sourceChainId, chainId)
        );
        bytes32 ethSignedMessageHash = messageHash.toEthSignedMessageHash();

        uint256 validSignatures = 0;
        for (uint256 i = 0; i < signatures.length; i++) {
            address signer = ethSignedMessageHash.recover(signatures[i]);
            if (hasRole(VALIDATOR_ROLE, signer) && !validatorSignatures[requestId][signer]) {
                validatorSignatures[requestId][signer] = true;
                validSignatures++;
            }
        }

        require(validSignatures >= minConfirmations, "Invalid signatures");

        // Mark as processed
        bridgeRequests[requestId].processed = true;

        // Transfer tokens to recipient
        IERC20(token).transfer(to, amount);

        emit BridgeCompleted(requestId, to, token, amount);
    }

    /**
     * @notice Emergency withdraw for stuck tokens
     * @param token Token address
     * @param amount Amount to withdraw
     */
    function emergencyWithdraw(
        address token,
        uint256 amount
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        IERC20(token).transfer(msg.sender, amount);
    }

    /**
     * @notice Set supported token
     * @param token Token address
     * @param supported Whether token is supported
     * @param limit Max transfer limit
     */
    function setSupportedToken(
        address token,
        bool supported,
        uint256 limit
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        supportedTokens[token] = supported;
        tokenLimits[token] = limit;
    }

    /**
     * @notice Set token mapping for cross-chain
     * @param localToken Local token address
     * @param remoteChainId Remote chain ID
     * @param remoteToken Remote token address
     */
    function setTokenMapping(
        address localToken,
        uint256 remoteChainId,
        address remoteToken
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        tokenMappings[localToken][remoteChainId] = remoteToken;
        emit TokenMappingSet(localToken, remoteChainId, remoteToken);
    }

    /**
     * @notice Update bridge fee
     * @param newFee New bridge fee
     */
    function updateBridgeFee(uint256 newFee) external onlyRole(DEFAULT_ADMIN_ROLE) {
        bridgeFee = newFee;
        emit FeeUpdated(newFee);
    }

    /**
     * @notice Update minimum confirmations required
     * @param newMin New minimum confirmations
     */
    function updateMinConfirmations(uint256 newMin) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newMin > 0, "Invalid confirmation count");
        minConfirmations = newMin;
    }

    /**
     * @notice Pause bridge operations
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause bridge operations
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }

    /**
     * @notice Get bridge request details
     * @param requestId Request ID
     */
    function getBridgeRequest(bytes32 requestId) external view returns (BridgeRequest memory) {
        return bridgeRequests[requestId];
    }

    /**
     * @notice Check if request has enough signatures
     * @param requestId Request ID
     */
    function hasEnoughSignatures(bytes32 requestId) external view returns (bool) {
        return signatureCount[requestId] >= minConfirmations;
    }
}