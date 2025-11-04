package l1

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// L1GatewayService provides methods to interact with L1Gateway contract
type L1GatewayService struct {
	client   *ethclient.Client
	Contract *L1Gateway
	address  common.Address
}

// NewL1GatewayService creates a new L1Gateway service
func NewL1GatewayService(client *ethclient.Client, contractAddress string) (*L1GatewayService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewL1Gateway(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate L1Gateway contract: %w", err)
	}

	return &L1GatewayService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *L1GatewayService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *L1GatewayService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}

// Note: Use s.Contract to call contract methods directly, for example:
// - s.Contract.SendMessageToL2(auth, target, data)
// - s.Contract.ReceiveMessageFromL2(&bind.CallOpts{Context: ctx}, messageHash)
