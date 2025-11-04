package l1

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
)

// CollateralVaultService provides methods to interact with CollateralVaultL1 contract
type CollateralVaultService struct {
	client   *ethclient.Client
	contract *CollateralVaultL1
	address  common.Address
}

// NewCollateralVaultService creates a new CollateralVaultL1 service
func NewCollateralVaultService(client *ethclient.Client, contractAddress string) (*CollateralVaultService, error) {
	address := common.HexToAddress(contractAddress)
	contract, err := NewCollateralVaultL1(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate CollateralVaultL1 contract: %w", err)
	}

	return &CollateralVaultService{
		client:   client,
		contract: contract,
		address:  address,
	}, nil
}

// ============ Read Methods ============

// GetLockedCollateral returns the locked collateral amount for a specific user
func (s *CollateralVaultService) GetLockedCollateral(ctx context.Context, user common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	amount, err := s.contract.GetLockedCollateral(opts, user)
	if err != nil {
		return nil, fmt.Errorf("failed to get locked collateral: %w", err)
	}
	return amount, nil
}

// GetTotalLocked returns the total locked collateral in the vault
func (s *CollateralVaultService) GetTotalLocked(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	amount, err := s.contract.GetTotalLocked(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get total locked: %w", err)
	}
	return amount, nil
}

// VaultStats represents vault statistics
type VaultStats struct {
	TotalLocked     *big.Int
	ContractBalance *big.Int
}

// GetVaultStats returns comprehensive vault statistics
func (s *CollateralVaultService) GetVaultStats(ctx context.Context) (*VaultStats, error) {
	opts := &bind.CallOpts{Context: ctx}
	stats, err := s.contract.GetVaultStats(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get vault stats: %w", err)
	}

	return &VaultStats{
		TotalLocked:     stats.TotalLocked,
		ContractBalance: stats.ContractBalance,
	}, nil
}

// GetRemainingDailyLockCapacity returns how much more can be locked today
func (s *CollateralVaultService) GetRemainingDailyLockCapacity(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	capacity, err := s.contract.GetRemainingDailyLockCapacity(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get remaining capacity: %w", err)
	}
	return capacity, nil
}

// CanExecuteEmergencyWithdrawal checks if a user can execute emergency withdrawal
func (s *CollateralVaultService) CanExecuteEmergencyWithdrawal(ctx context.Context, user common.Address) (bool, error) {
	opts := &bind.CallOpts{Context: ctx}
	canExecute, err := s.contract.CanExecuteEmergencyWithdrawal(opts, user)
	if err != nil {
		return false, fmt.Errorf("failed to check emergency withdrawal: %w", err)
	}
	return canExecute, nil
}

// EmergencyWithdrawal represents emergency withdrawal details
type EmergencyWithdrawal struct {
	Amount      *big.Int
	RequestTime *big.Int
	Executed    bool
	UnlockTime  *big.Int
}

// GetEmergencyWithdrawal returns emergency withdrawal details for a user
func (s *CollateralVaultService) GetEmergencyWithdrawal(ctx context.Context, user common.Address) (*EmergencyWithdrawal, error) {
	opts := &bind.CallOpts{Context: ctx}
	withdrawal, err := s.contract.GetEmergencyWithdrawal(opts, user)
	if err != nil {
		return nil, fmt.Errorf("failed to get emergency withdrawal: %w", err)
	}

	return &EmergencyWithdrawal{
		Amount:      withdrawal.Amount,
		RequestTime: withdrawal.RequestTime,
		Executed:    withdrawal.Executed,
		UnlockTime:  withdrawal.UnlockTime,
	}, nil
}

// GetCollateralToken returns the collateral token address
func (s *CollateralVaultService) GetCollateralToken(ctx context.Context) (common.Address, error) {
	opts := &bind.CallOpts{Context: ctx}
	token, err := s.contract.CollateralToken(opts)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get collateral token: %w", err)
	}
	return token, nil
}

// GetL2Bridge returns the L2 bridge address
func (s *CollateralVaultService) GetL2Bridge(ctx context.Context) (common.Address, error) {
	opts := &bind.CallOpts{Context: ctx}
	bridge, err := s.contract.L2Bridge(opts)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get L2 bridge: %w", err)
	}
	return bridge, nil
}

// GetStateRegistry returns the state registry address
func (s *CollateralVaultService) GetStateRegistry(ctx context.Context) (common.Address, error) {
	opts := &bind.CallOpts{Context: ctx}
	registry, err := s.contract.StateRegistry(opts)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get state registry: %w", err)
	}
	return registry, nil
}

// IsEmergencyPaused checks if the vault is in emergency pause mode
func (s *CollateralVaultService) IsEmergencyPaused(ctx context.Context) (bool, error) {
	opts := &bind.CallOpts{Context: ctx}
	paused, err := s.contract.EmergencyPaused(opts)
	if err != nil {
		return false, fmt.Errorf("failed to check emergency pause: %w", err)
	}
	return paused, nil
}

// GetDailyLockLimit returns the daily lock limit
func (s *CollateralVaultService) GetDailyLockLimit(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	limit, err := s.contract.DailyLockLimit(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily lock limit: %w", err)
	}
	return limit, nil
}

// GetContractAddress returns the contract address
func (s *CollateralVaultService) GetContractAddress() common.Address {
	return s.address
}

// ============ Write Methods (User) ============

// LockCollateral locks collateral tokens (requires prior token approval)
// user: the user address to lock for
// amount: amount to lock
// l2TxHash: corresponding L2 transaction hash
func (s *CollateralVaultService) LockCollateral(
	auth *bind.TransactOpts,
	user common.Address,
	amount *big.Int,
	l2TxHash [32]byte,
) (*types.Transaction, error) {
	tx, err := s.contract.LockCollateral(auth, user, amount, l2TxHash)
	if err != nil {
		return nil, fmt.Errorf("failed to lock collateral: %w", err)
	}
	return tx, nil
}

// UnlockCollateral unlocks collateral tokens
// user: the user address to unlock for
// amount: amount to unlock
// l2TxHash: corresponding L2 transaction hash
func (s *CollateralVaultService) UnlockCollateral(
	auth *bind.TransactOpts,
	user common.Address,
	amount *big.Int,
	l2TxHash [32]byte,
) (*types.Transaction, error) {
	tx, err := s.contract.UnlockCollateral(auth, user, amount, l2TxHash)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock collateral: %w", err)
	}
	return tx, nil
}

// RequestEmergencyWithdrawal requests an emergency withdrawal (starts cooldown)
func (s *CollateralVaultService) RequestEmergencyWithdrawal(auth *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	tx, err := s.contract.RequestEmergencyWithdrawal(auth, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to request emergency withdrawal: %w", err)
	}
	return tx, nil
}

// ExecuteEmergencyWithdrawal executes an emergency withdrawal (after cooldown)
func (s *CollateralVaultService) ExecuteEmergencyWithdrawal(auth *bind.TransactOpts) (*types.Transaction, error) {
	tx, err := s.contract.ExecuteEmergencyWithdrawal(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to execute emergency withdrawal: %w", err)
	}
	return tx, nil
}

// ============ Write Methods (Admin Only) ============

// SetDailyLockLimit sets the daily lock limit (admin only)
func (s *CollateralVaultService) SetDailyLockLimit(auth *bind.TransactOpts, limit *big.Int) (*types.Transaction, error) {
	tx, err := s.contract.SetDailyLockLimit(auth, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to set daily lock limit: %w", err)
	}
	return tx, nil
}

// SetL2Bridge sets the L2 bridge address (admin only)
func (s *CollateralVaultService) SetL2Bridge(auth *bind.TransactOpts, bridge common.Address) (*types.Transaction, error) {
	tx, err := s.contract.SetL2Bridge(auth, bridge)
	if err != nil {
		return nil, fmt.Errorf("failed to set L2 bridge: %w", err)
	}
	return tx, nil
}

// SetStateRegistry sets the state registry address (admin only)
func (s *CollateralVaultService) SetStateRegistry(auth *bind.TransactOpts, registry common.Address) (*types.Transaction, error) {
	tx, err := s.contract.SetStateRegistry(auth, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to set state registry: %w", err)
	}
	return tx, nil
}

// TriggerEmergencyPause pauses the vault for emergency (admin only)
func (s *CollateralVaultService) TriggerEmergencyPause(auth *bind.TransactOpts) (*types.Transaction, error) {
	tx, err := s.contract.TriggerEmergencyPause(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger emergency pause: %w", err)
	}
	return tx, nil
}

// ResumeFromEmergency resumes the vault after emergency (admin only)
func (s *CollateralVaultService) ResumeFromEmergency(auth *bind.TransactOpts) (*types.Transaction, error) {
	tx, err := s.contract.ResumeFromEmergency(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to resume from emergency: %w", err)
	}
	return tx, nil
}

// ============ Event Watchers ============

// WatchCollateralLocked watches for CollateralLocked events
func (s *CollateralVaultService) WatchCollateralLocked(
	ctx context.Context,
	sink chan<- *CollateralVaultL1CollateralLocked,
	user []common.Address,
	l2TxHash [][32]byte,
) (event.Subscription, error) {
	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.contract.WatchCollateralLocked(watchOpts, sink, user, l2TxHash)
	if err != nil {
		return nil, fmt.Errorf("failed to watch CollateralLocked events: %w", err)
	}
	return sub, nil
}

// WatchCollateralUnlocked watches for CollateralUnlocked events
func (s *CollateralVaultService) WatchCollateralUnlocked(
	ctx context.Context,
	sink chan<- *CollateralVaultL1CollateralUnlocked,
	user []common.Address,
	l2TxHash [][32]byte,
) (event.Subscription, error) {
	watchOpts := &bind.WatchOpts{Context: ctx}
	sub, err := s.contract.WatchCollateralUnlocked(watchOpts, sink, user, l2TxHash)
	if err != nil {
		return nil, fmt.Errorf("failed to watch CollateralUnlocked events: %w", err)
	}
	return sub, nil
}

// ============ Helper Methods ============

// WaitForTransaction waits for a transaction to be mined and checks status
func (s *CollateralVaultService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed with status 0")
	}
	return receipt, nil
}
