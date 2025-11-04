package l2

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// RWAMarketplaceService provides methods to interact with RWAMarketplace contract
type RWAMarketplaceService struct {
	client   *ethclient.Client
	Contract *RWAMarketplace
	address  common.Address
}

// NewRWAMarketplaceService creates a new RWAMarketplace service
func NewRWAMarketplaceService(client *ethclient.Client, contractAddress string) (*RWAMarketplaceService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewRWAMarketplace(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate RWAMarketplace contract: %w", err)
	}
	return &RWAMarketplaceService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *RWAMarketplaceService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *RWAMarketplaceService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}
