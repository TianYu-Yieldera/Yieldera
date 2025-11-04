package l2

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// CompoundV3AdapterService provides methods to interact with CompoundV3Adapter contract
type CompoundV3AdapterService struct {
	client   *ethclient.Client
	Contract *CompoundV3Adapter
	address  common.Address
}

// NewCompoundV3AdapterService creates a new CompoundV3Adapter service
func NewCompoundV3AdapterService(client *ethclient.Client, contractAddress string) (*CompoundV3AdapterService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewCompoundV3Adapter(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate CompoundV3Adapter contract: %w", err)
	}
	return &CompoundV3AdapterService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *CompoundV3AdapterService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *CompoundV3AdapterService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}
