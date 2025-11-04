package l1

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// LoyaltyUSDService provides methods to interact with LoyaltyUSDL1 contract
type LoyaltyUSDService struct {
	client   *ethclient.Client
	Contract *LoyaltyUSDL1
	address  common.Address
}

// NewLoyaltyUSDService creates a new LoyaltyUSDL1 service
func NewLoyaltyUSDService(client *ethclient.Client, contractAddress string) (*LoyaltyUSDService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewLoyaltyUSDL1(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate LoyaltyUSDL1 contract: %w", err)
	}

	return &LoyaltyUSDService{
		client:   client,
		Contract: contract,
		address:  address,
	}, nil
}

// GetContractAddress returns the contract address
func (s *LoyaltyUSDService) GetContractAddress() common.Address {
	return s.address
}

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *LoyaltyUSDService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
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
// - s.Contract.BalanceOf(&bind.CallOpts{Context: ctx}, account)
// - s.Contract.Transfer(auth, recipient, amount)
// - s.Contract.Mint(auth, to, amount)
