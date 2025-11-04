package l2

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// RWAAssetFactoryService provides methods to interact with RWAAssetFactory contract
type RWAAssetFactoryService struct {
	client   *ethclient.Client
	Contract *RWAAssetFactory
	address  common.Address
}

// NewRWAAssetFactoryService creates a new RWAAssetFactory service
func NewRWAAssetFactoryService(client *ethclient.Client, contractAddress string) (*RWAAssetFactoryService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewRWAAssetFactory(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate RWAAssetFactory contract: %w", err)
	}
	return &RWAAssetFactoryService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *RWAAssetFactoryService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *RWAAssetFactoryService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}
