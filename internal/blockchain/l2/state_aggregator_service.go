package l2

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// L2StateAggregatorService provides methods to interact with L2StateAggregator contract
type L2StateAggregatorService struct {
	client   *ethclient.Client
	Contract *L2StateAggregator
	address  common.Address
}

// NewL2StateAggregatorService creates a new L2StateAggregator service
func NewL2StateAggregatorService(client *ethclient.Client, contractAddress string) (*L2StateAggregatorService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewL2StateAggregator(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate L2StateAggregator contract: %w", err)
	}
	return &L2StateAggregatorService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *L2StateAggregatorService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *L2StateAggregatorService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}
