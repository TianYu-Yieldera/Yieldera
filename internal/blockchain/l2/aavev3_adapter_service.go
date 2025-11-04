package l2

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// AaveV3AdapterService provides methods to interact with AaveV3Adapter contract
type AaveV3AdapterService struct {
	client   *ethclient.Client
	Contract *AaveV3Adapter
	address  common.Address
}

// NewAaveV3AdapterService creates a new AaveV3Adapter service
func NewAaveV3AdapterService(client *ethclient.Client, contractAddress string) (*AaveV3AdapterService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewAaveV3Adapter(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate AaveV3Adapter contract: %w", err)
	}
	return &AaveV3AdapterService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *AaveV3AdapterService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *AaveV3AdapterService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}
