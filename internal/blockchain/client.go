package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client wraps an Ethereum client with additional utilities
type Client struct {
	ethClient *ethclient.Client
	chainID   *big.Int
	rpcURL    string
}

// Config holds configuration for blockchain client
type Config struct {
	RPCURL  string
	ChainID int64
}

// NewClient creates a new blockchain client
func NewClient(cfg Config) (*Client, error) {
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	chainID := big.NewInt(cfg.ChainID)

	log.Printf("✅ Connected to Ethereum node: %s (ChainID: %d)", cfg.RPCURL, cfg.ChainID)

	return &Client{
		ethClient: client,
		chainID:   chainID,
		rpcURL:    cfg.RPCURL,
	}, nil
}

// Close closes the Ethereum client connection
func (c *Client) Close() {
	c.ethClient.Close()
}

// GetClient returns the underlying eth client
func (c *Client) GetClient() *ethclient.Client {
	return c.ethClient
}

// GetChainID returns the chain ID
func (c *Client) GetChainID() *big.Int {
	return c.chainID
}

// GetBlockNumber returns the latest block number
func (c *Client) GetBlockNumber(ctx context.Context) (uint64, error) {
	return c.ethClient.BlockNumber(ctx)
}

// GetBalance returns the balance of an address
func (c *Client) GetBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	return c.ethClient.BalanceAt(ctx, address, nil)
}

// EstimateGas estimates gas for a transaction
func (c *Client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return c.ethClient.EstimateGas(ctx, msg)
}

// SuggestGasPrice suggests a gas price
func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.ethClient.SuggestGasPrice(ctx)
}

// WaitForTransaction waits for a transaction to be mined
func (c *Client) WaitForTransaction(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	log.Printf("⏳ Waiting for transaction: %s", txHash.Hex())

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for transaction: %s", txHash.Hex())
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			receipt, err := c.ethClient.TransactionReceipt(ctx, txHash)
			if err != nil {
				continue
			}

			if receipt.Status == types.ReceiptStatusSuccessful {
				log.Printf("✅ Transaction successful: %s (Block: %d, Gas: %d)",
					txHash.Hex(), receipt.BlockNumber.Uint64(), receipt.GasUsed)
				return receipt, nil
			}
			return receipt, fmt.Errorf("transaction failed")
		}
	}
}

// GetTransactionOpts returns transaction options for signing
func (c *Client) GetTransactionOpts(ctx context.Context, privateKey *ecdsa.PrivateKey, value *big.Int) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, c.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := c.ethClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = uint64(0)
	auth.GasPrice = gasPrice
	auth.Context = ctx

	return auth, nil
}

// GetCallOpts returns call options for read-only calls
func (c *Client) GetCallOpts(ctx context.Context) *bind.CallOpts {
	return &bind.CallOpts{
		Context: ctx,
		Pending: false,
	}
}

// IsContract checks if an address is a contract
func (c *Client) IsContract(ctx context.Context, address common.Address) (bool, error) {
	code, err := c.ethClient.CodeAt(ctx, address, nil)
	if err != nil {
		return false, err
	}
	return len(code) > 0, nil
}

// TransactionStatus represents the status of a transaction
type TransactionStatus struct {
	Hash        common.Hash
	Status      string
	BlockNumber uint64
	GasUsed     uint64
	Error       error
}

// GetTransactionStatus returns the status of a transaction
func (c *Client) GetTransactionStatus(ctx context.Context, txHash common.Hash) (*TransactionStatus, error) {
	receipt, err := c.ethClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		_, isPending, err := c.ethClient.TransactionByHash(ctx, txHash)
		if err != nil {
			return nil, fmt.Errorf("transaction not found: %w", err)
		}

		if isPending {
			return &TransactionStatus{
				Hash:   txHash,
				Status: "pending",
			}, nil
		}

		return nil, fmt.Errorf("transaction not found")
	}

	status := &TransactionStatus{
		Hash:        txHash,
		BlockNumber: receipt.BlockNumber.Uint64(),
		GasUsed:     receipt.GasUsed,
	}

	if receipt.Status == types.ReceiptStatusSuccessful {
		status.Status = "success"
	} else {
		status.Status = "failed"
		status.Error = fmt.Errorf("transaction reverted")
	}

	return status, nil
}

// GasSettings contains gas pricing information
type GasSettings struct {
	GasPrice             *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
}

// GetGasSettings returns recommended gas settings
func (c *Client) GetGasSettings(ctx context.Context) (*GasSettings, error) {
	baseFee, err := c.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	maxPriorityFeePerGas := big.NewInt(2000000000)

	maxFeePerGas := new(big.Int).Add(
		new(big.Int).Mul(baseFee, big.NewInt(2)),
		maxPriorityFeePerGas,
	)

	return &GasSettings{
		GasPrice:             baseFee,
		MaxFeePerGas:         maxFeePerGas,
		MaxPriorityFeePerGas: maxPriorityFeePerGas,
	}, nil
}
