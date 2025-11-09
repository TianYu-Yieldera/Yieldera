package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client wraps Ethereum client with useful methods
type Client struct {
	L1Client *ethclient.Client
	L2Client *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
	FromAddress common.Address
}

// NewClient creates a new blockchain client
func NewClient() (*Client, error) {
	// Get RPC URLs from environment
	l1RPC := os.Getenv("L1_RPC_URL")
	l2RPC := os.Getenv("L2_RPC_URL")
	privateKeyHex := os.Getenv("PRIVATE_KEY")

	if l1RPC == "" || l2RPC == "" {
		return nil, fmt.Errorf("L1_RPC_URL and L2_RPC_URL must be set")
	}

	// Connect to L1
	l1Client, err := ethclient.Dial(l1RPC)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to L1: %w", err)
	}

	// Connect to L2
	l2Client, err := ethclient.Dial(l2RPC)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to L2: %w", err)
	}

	// Load private key if provided
	var privateKey *ecdsa.PrivateKey
	var fromAddress common.Address

	if privateKeyHex != "" {
		// Remove 0x prefix if present
		if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
			privateKeyHex = privateKeyHex[2:]
		}

		privateKey, err = crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("failed to cast public key to ECDSA")
		}

		fromAddress = crypto.PubkeyToAddress(*publicKeyECDSA)
	}

	return &Client{
		L1Client:    l1Client,
		L2Client:    l2Client,
		PrivateKey:  privateKey,
		FromAddress: fromAddress,
	}, nil
}

// GetTransactOpts creates transaction options for L1
func (c *Client) GetL1TransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	if c.PrivateKey == nil {
		return nil, fmt.Errorf("private key not configured")
	}

	chainID, err := c.L1Client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.PrivateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// Get suggested gas price
	gasPrice, err := c.L1Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	auth.GasPrice = gasPrice
	auth.Context = ctx

	return auth, nil
}

// GetL2TransactOpts creates transaction options for L2
func (c *Client) GetL2TransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	if c.PrivateKey == nil {
		return nil, fmt.Errorf("private key not configured")
	}

	chainID, err := c.L2Client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.PrivateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// Get suggested gas price
	gasPrice, err := c.L2Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	auth.GasPrice = gasPrice
	auth.Context = ctx

	return auth, nil
}

// WaitForTransaction waits for transaction to be mined on L1
func (c *Client) WaitForL1Transaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, c.L1Client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return receipt, fmt.Errorf("transaction failed with status %d", receipt.Status)
	}

	return receipt, nil
}

// WaitForL2Transaction waits for transaction to be mined on L2
func (c *Client) WaitForL2Transaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, c.L2Client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return receipt, fmt.Errorf("transaction failed with status %d", receipt.Status)
	}

	return receipt, nil
}

// Close closes both clients
func (c *Client) Close() {
	if c.L1Client != nil {
		c.L1Client.Close()
	}
	if c.L2Client != nil {
		c.L2Client.Close()
	}
}

// VerifySignature verifies EIP-712 signature
func VerifySignature(message []byte, signature []byte, expectedAddress common.Address) (bool, error) {
	if len(signature) != 65 {
		return false, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	// Transform yellow paper V from 27/28 to 0/1
	if signature[64] >= 27 {
		signature[64] -= 27
	}

	// Recover public key from signature
	pubKey, err := crypto.SigToPub(message, signature)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	// Get address from public key
	recoveredAddress := crypto.PubkeyToAddress(*pubKey)

	// Compare addresses
	return recoveredAddress == expectedAddress, nil
}

// ParseBigInt parses string to big.Int
func ParseBigInt(s string) (*big.Int, error) {
	value := new(big.Int)
	_, ok := value.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid big int: %s", s)
	}
	return value, nil
}
