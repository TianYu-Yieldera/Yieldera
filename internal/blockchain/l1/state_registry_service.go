package l1

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// L1StateRegistryService provides methods to interact with L1StateRegistry contract
type L1StateRegistryService struct {
	client   *ethclient.Client
	Contract *L1StateRegistry
	address  common.Address
}

// NewL1StateRegistryService creates a new L1StateRegistry service
func NewL1StateRegistryService(client *ethclient.Client, contractAddress string) (*L1StateRegistryService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewL1StateRegistry(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate L1StateRegistry contract: %w", err)
	}

	return &L1StateRegistryService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *L1StateRegistryService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *L1StateRegistryService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}

// Note: Use s.Contract to call contract methods directly, for example:
// - s.Contract.GetLatestState(&bind.CallOpts{Context: ctx})
// - s.Contract.ReceiveStateRoot(auth, stateRoot, l2Block, timestamp, blockNumber)
// - s.Contract.VerifyStateRoot(&bind.CallOpts{Context: ctx}, l2Block, stateRoot)
