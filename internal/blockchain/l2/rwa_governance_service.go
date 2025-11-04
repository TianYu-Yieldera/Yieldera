package l2

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// RWAGovernanceService provides methods to interact with RWAGovernance contract
type RWAGovernanceService struct {
	client   *ethclient.Client
	Contract *RWAGovernance
	address  common.Address
}

// NewRWAGovernanceService creates a new RWAGovernance service
func NewRWAGovernanceService(client *ethclient.Client, contractAddress string) (*RWAGovernanceService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewRWAGovernance(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate RWAGovernance contract: %w", err)
	}
	return &RWAGovernanceService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *RWAGovernanceService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *RWAGovernanceService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}
