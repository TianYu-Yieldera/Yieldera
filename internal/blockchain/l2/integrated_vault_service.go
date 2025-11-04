package l2

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// IntegratedVaultService provides methods to interact with IntegratedVault contract
type IntegratedVaultService struct {
	client   *ethclient.Client
	Contract *IntegratedVault
	address  common.Address
}

// NewIntegratedVaultService creates a new IntegratedVault service
func NewIntegratedVaultService(client *ethclient.Client, contractAddress string) (*IntegratedVaultService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewIntegratedVault(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate IntegratedVault contract: %w", err)
	}
	return &IntegratedVaultService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *IntegratedVaultService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *IntegratedVaultService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}
