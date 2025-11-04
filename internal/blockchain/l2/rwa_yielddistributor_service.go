package l2

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// RWAYieldDistributorDistributorService provides methods to interact with RWAYieldDistributor contract
type RWAYieldDistributorDistributorService struct {
	client   *ethclient.Client
	Contract *RWAYieldDistributor
	address  common.Address
}

// NewRWAYieldDistributorDistributorService creates a new RWAYieldDistributor service
func NewRWAYieldDistributorDistributorService(client *ethclient.Client, contractAddress string) (*RWAYieldDistributorDistributorService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewRWAYieldDistributor(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate RWAYieldDistributor contract: %w", err)
	}
	return &RWAYieldDistributorDistributorService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *RWAYieldDistributorDistributorService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *RWAYieldDistributorDistributorService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}
