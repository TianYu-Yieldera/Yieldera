package l2

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// RWAValuationService provides methods to interact with RWAValuation contract
type RWAValuationService struct {
	client   *ethclient.Client
	Contract *RWAValuation
	address  common.Address
}

// NewRWAValuationService creates a new RWAValuation service
func NewRWAValuationService(client *ethclient.Client, contractAddress string) (*RWAValuationService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewRWAValuation(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate RWAValuation contract: %w", err)
	}
	return &RWAValuationService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *RWAValuationService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *RWAValuationService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}
