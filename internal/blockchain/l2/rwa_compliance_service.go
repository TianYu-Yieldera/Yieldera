package l2

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// RWAComplianceService provides methods to interact with RWACompliance contract
type RWAComplianceService struct {
	client   *ethclient.Client
	Contract *RWACompliance
	address  common.Address
}

// NewRWAComplianceService creates a new RWACompliance service
func NewRWAComplianceService(client *ethclient.Client, contractAddress string) (*RWAComplianceService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewRWACompliance(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate RWACompliance contract: %w", err)
	}
	return &RWAComplianceService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *RWAComplianceService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *RWAComplianceService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}
