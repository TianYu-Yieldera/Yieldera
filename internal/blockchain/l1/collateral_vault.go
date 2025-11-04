// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package l1

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// CollateralVaultL1MetaData contains all meta data concerning the CollateralVaultL1 contract.
var CollateralVaultL1MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_collateralToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalUserLocked\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"l2TxHash\",\"type\":\"bytes32\"}],\"name\":\"CollateralLocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remaining\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"l2TxHash\",\"type\":\"bytes32\"}],\"name\":\"CollateralUnlocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"triggeredBy\",\"type\":\"address\"}],\"name\":\"EmergencyPauseTriggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"resumedBy\",\"type\":\"address\"}],\"name\":\"EmergencyResumed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EmergencyWithdrawalExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"}],\"name\":\"EmergencyWithdrawalRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldBridge\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newBridge\",\"type\":\"address\"}],\"name\":\"L2BridgeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldRegistry\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newRegistry\",\"type\":\"address\"}],\"name\":\"StateRegistryUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"EMERGENCY_DELAY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_LOCK_PER_TX\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"canExecuteEmergencyWithdrawal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collateralToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dailyLockLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dailyLockedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emergencyPaused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"emergencyWithdrawals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestTime\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"executed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"executeEmergencyWithdrawal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getEmergencyWithdrawal\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestTime\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"executed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getLockedCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRemainingDailyLockCapacity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTimeUntilLimitReset\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalLocked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVaultStats\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_totalLocked\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"contractBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2Bridge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastLockResetTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"l2TxHash\",\"type\":\"bytes32\"}],\"name\":\"lockCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"lockedCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"requestEmergencyWithdrawal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resumeFromEmergency\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"setDailyLockLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l2Bridge\",\"type\":\"address\"}],\"name\":\"setL2Bridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stateRegistry\",\"type\":\"address\"}],\"name\":\"setStateRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateRegistry\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalLocked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"triggerEmergencyPause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"l2TxHash\",\"type\":\"bytes32\"}],\"name\":\"unlockCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// CollateralVaultL1ABI is the input ABI used to generate the binding from.
// Deprecated: Use CollateralVaultL1MetaData.ABI instead.
var CollateralVaultL1ABI = CollateralVaultL1MetaData.ABI

// CollateralVaultL1 is an auto generated Go binding around an Ethereum contract.
type CollateralVaultL1 struct {
	CollateralVaultL1Caller     // Read-only binding to the contract
	CollateralVaultL1Transactor // Write-only binding to the contract
	CollateralVaultL1Filterer   // Log filterer for contract events
}

// CollateralVaultL1Caller is an auto generated read-only Go binding around an Ethereum contract.
type CollateralVaultL1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralVaultL1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type CollateralVaultL1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralVaultL1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CollateralVaultL1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralVaultL1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CollateralVaultL1Session struct {
	Contract     *CollateralVaultL1 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// CollateralVaultL1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CollateralVaultL1CallerSession struct {
	Contract *CollateralVaultL1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// CollateralVaultL1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CollateralVaultL1TransactorSession struct {
	Contract     *CollateralVaultL1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// CollateralVaultL1Raw is an auto generated low-level Go binding around an Ethereum contract.
type CollateralVaultL1Raw struct {
	Contract *CollateralVaultL1 // Generic contract binding to access the raw methods on
}

// CollateralVaultL1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CollateralVaultL1CallerRaw struct {
	Contract *CollateralVaultL1Caller // Generic read-only contract binding to access the raw methods on
}

// CollateralVaultL1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CollateralVaultL1TransactorRaw struct {
	Contract *CollateralVaultL1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewCollateralVaultL1 creates a new instance of CollateralVaultL1, bound to a specific deployed contract.
func NewCollateralVaultL1(address common.Address, backend bind.ContractBackend) (*CollateralVaultL1, error) {
	contract, err := bindCollateralVaultL1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1{CollateralVaultL1Caller: CollateralVaultL1Caller{contract: contract}, CollateralVaultL1Transactor: CollateralVaultL1Transactor{contract: contract}, CollateralVaultL1Filterer: CollateralVaultL1Filterer{contract: contract}}, nil
}

// NewCollateralVaultL1Caller creates a new read-only instance of CollateralVaultL1, bound to a specific deployed contract.
func NewCollateralVaultL1Caller(address common.Address, caller bind.ContractCaller) (*CollateralVaultL1Caller, error) {
	contract, err := bindCollateralVaultL1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1Caller{contract: contract}, nil
}

// NewCollateralVaultL1Transactor creates a new write-only instance of CollateralVaultL1, bound to a specific deployed contract.
func NewCollateralVaultL1Transactor(address common.Address, transactor bind.ContractTransactor) (*CollateralVaultL1Transactor, error) {
	contract, err := bindCollateralVaultL1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1Transactor{contract: contract}, nil
}

// NewCollateralVaultL1Filterer creates a new log filterer instance of CollateralVaultL1, bound to a specific deployed contract.
func NewCollateralVaultL1Filterer(address common.Address, filterer bind.ContractFilterer) (*CollateralVaultL1Filterer, error) {
	contract, err := bindCollateralVaultL1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1Filterer{contract: contract}, nil
}

// bindCollateralVaultL1 binds a generic wrapper to an already deployed contract.
func bindCollateralVaultL1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CollateralVaultL1MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CollateralVaultL1 *CollateralVaultL1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CollateralVaultL1.Contract.CollateralVaultL1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CollateralVaultL1 *CollateralVaultL1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.CollateralVaultL1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CollateralVaultL1 *CollateralVaultL1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.CollateralVaultL1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CollateralVaultL1 *CollateralVaultL1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CollateralVaultL1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CollateralVaultL1 *CollateralVaultL1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CollateralVaultL1 *CollateralVaultL1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.contract.Transact(opts, method, params...)
}

// EMERGENCYDELAY is a free data retrieval call binding the contract method 0x82944e2d.
//
// Solidity: function EMERGENCY_DELAY() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Caller) EMERGENCYDELAY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "EMERGENCY_DELAY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EMERGENCYDELAY is a free data retrieval call binding the contract method 0x82944e2d.
//
// Solidity: function EMERGENCY_DELAY() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Session) EMERGENCYDELAY() (*big.Int, error) {
	return _CollateralVaultL1.Contract.EMERGENCYDELAY(&_CollateralVaultL1.CallOpts)
}

// EMERGENCYDELAY is a free data retrieval call binding the contract method 0x82944e2d.
//
// Solidity: function EMERGENCY_DELAY() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) EMERGENCYDELAY() (*big.Int, error) {
	return _CollateralVaultL1.Contract.EMERGENCYDELAY(&_CollateralVaultL1.CallOpts)
}

// MAXLOCKPERTX is a free data retrieval call binding the contract method 0xed846607.
//
// Solidity: function MAX_LOCK_PER_TX() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Caller) MAXLOCKPERTX(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "MAX_LOCK_PER_TX")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXLOCKPERTX is a free data retrieval call binding the contract method 0xed846607.
//
// Solidity: function MAX_LOCK_PER_TX() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Session) MAXLOCKPERTX() (*big.Int, error) {
	return _CollateralVaultL1.Contract.MAXLOCKPERTX(&_CollateralVaultL1.CallOpts)
}

// MAXLOCKPERTX is a free data retrieval call binding the contract method 0xed846607.
//
// Solidity: function MAX_LOCK_PER_TX() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) MAXLOCKPERTX() (*big.Int, error) {
	return _CollateralVaultL1.Contract.MAXLOCKPERTX(&_CollateralVaultL1.CallOpts)
}

// CanExecuteEmergencyWithdrawal is a free data retrieval call binding the contract method 0x3a944126.
//
// Solidity: function canExecuteEmergencyWithdrawal(address user) view returns(bool)
func (_CollateralVaultL1 *CollateralVaultL1Caller) CanExecuteEmergencyWithdrawal(opts *bind.CallOpts, user common.Address) (bool, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "canExecuteEmergencyWithdrawal", user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanExecuteEmergencyWithdrawal is a free data retrieval call binding the contract method 0x3a944126.
//
// Solidity: function canExecuteEmergencyWithdrawal(address user) view returns(bool)
func (_CollateralVaultL1 *CollateralVaultL1Session) CanExecuteEmergencyWithdrawal(user common.Address) (bool, error) {
	return _CollateralVaultL1.Contract.CanExecuteEmergencyWithdrawal(&_CollateralVaultL1.CallOpts, user)
}

// CanExecuteEmergencyWithdrawal is a free data retrieval call binding the contract method 0x3a944126.
//
// Solidity: function canExecuteEmergencyWithdrawal(address user) view returns(bool)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) CanExecuteEmergencyWithdrawal(user common.Address) (bool, error) {
	return _CollateralVaultL1.Contract.CanExecuteEmergencyWithdrawal(&_CollateralVaultL1.CallOpts, user)
}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_CollateralVaultL1 *CollateralVaultL1Caller) CollateralToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "collateralToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_CollateralVaultL1 *CollateralVaultL1Session) CollateralToken() (common.Address, error) {
	return _CollateralVaultL1.Contract.CollateralToken(&_CollateralVaultL1.CallOpts)
}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) CollateralToken() (common.Address, error) {
	return _CollateralVaultL1.Contract.CollateralToken(&_CollateralVaultL1.CallOpts)
}

// DailyLockLimit is a free data retrieval call binding the contract method 0xa73497df.
//
// Solidity: function dailyLockLimit() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Caller) DailyLockLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "dailyLockLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DailyLockLimit is a free data retrieval call binding the contract method 0xa73497df.
//
// Solidity: function dailyLockLimit() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Session) DailyLockLimit() (*big.Int, error) {
	return _CollateralVaultL1.Contract.DailyLockLimit(&_CollateralVaultL1.CallOpts)
}

// DailyLockLimit is a free data retrieval call binding the contract method 0xa73497df.
//
// Solidity: function dailyLockLimit() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) DailyLockLimit() (*big.Int, error) {
	return _CollateralVaultL1.Contract.DailyLockLimit(&_CollateralVaultL1.CallOpts)
}

// DailyLockedAmount is a free data retrieval call binding the contract method 0x729d0715.
//
// Solidity: function dailyLockedAmount() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Caller) DailyLockedAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "dailyLockedAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DailyLockedAmount is a free data retrieval call binding the contract method 0x729d0715.
//
// Solidity: function dailyLockedAmount() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Session) DailyLockedAmount() (*big.Int, error) {
	return _CollateralVaultL1.Contract.DailyLockedAmount(&_CollateralVaultL1.CallOpts)
}

// DailyLockedAmount is a free data retrieval call binding the contract method 0x729d0715.
//
// Solidity: function dailyLockedAmount() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) DailyLockedAmount() (*big.Int, error) {
	return _CollateralVaultL1.Contract.DailyLockedAmount(&_CollateralVaultL1.CallOpts)
}

// EmergencyPaused is a free data retrieval call binding the contract method 0x27c830a9.
//
// Solidity: function emergencyPaused() view returns(bool)
func (_CollateralVaultL1 *CollateralVaultL1Caller) EmergencyPaused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "emergencyPaused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// EmergencyPaused is a free data retrieval call binding the contract method 0x27c830a9.
//
// Solidity: function emergencyPaused() view returns(bool)
func (_CollateralVaultL1 *CollateralVaultL1Session) EmergencyPaused() (bool, error) {
	return _CollateralVaultL1.Contract.EmergencyPaused(&_CollateralVaultL1.CallOpts)
}

// EmergencyPaused is a free data retrieval call binding the contract method 0x27c830a9.
//
// Solidity: function emergencyPaused() view returns(bool)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) EmergencyPaused() (bool, error) {
	return _CollateralVaultL1.Contract.EmergencyPaused(&_CollateralVaultL1.CallOpts)
}

// EmergencyWithdrawals is a free data retrieval call binding the contract method 0xdc6bae38.
//
// Solidity: function emergencyWithdrawals(address ) view returns(uint256 amount, uint256 requestTime, bool executed)
func (_CollateralVaultL1 *CollateralVaultL1Caller) EmergencyWithdrawals(opts *bind.CallOpts, arg0 common.Address) (struct {
	Amount      *big.Int
	RequestTime *big.Int
	Executed    bool
}, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "emergencyWithdrawals", arg0)

	outstruct := new(struct {
		Amount      *big.Int
		RequestTime *big.Int
		Executed    bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.RequestTime = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Executed = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// EmergencyWithdrawals is a free data retrieval call binding the contract method 0xdc6bae38.
//
// Solidity: function emergencyWithdrawals(address ) view returns(uint256 amount, uint256 requestTime, bool executed)
func (_CollateralVaultL1 *CollateralVaultL1Session) EmergencyWithdrawals(arg0 common.Address) (struct {
	Amount      *big.Int
	RequestTime *big.Int
	Executed    bool
}, error) {
	return _CollateralVaultL1.Contract.EmergencyWithdrawals(&_CollateralVaultL1.CallOpts, arg0)
}

// EmergencyWithdrawals is a free data retrieval call binding the contract method 0xdc6bae38.
//
// Solidity: function emergencyWithdrawals(address ) view returns(uint256 amount, uint256 requestTime, bool executed)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) EmergencyWithdrawals(arg0 common.Address) (struct {
	Amount      *big.Int
	RequestTime *big.Int
	Executed    bool
}, error) {
	return _CollateralVaultL1.Contract.EmergencyWithdrawals(&_CollateralVaultL1.CallOpts, arg0)
}

// GetEmergencyWithdrawal is a free data retrieval call binding the contract method 0x9ac406a5.
//
// Solidity: function getEmergencyWithdrawal(address user) view returns(uint256 amount, uint256 requestTime, bool executed, uint256 unlockTime)
func (_CollateralVaultL1 *CollateralVaultL1Caller) GetEmergencyWithdrawal(opts *bind.CallOpts, user common.Address) (struct {
	Amount      *big.Int
	RequestTime *big.Int
	Executed    bool
	UnlockTime  *big.Int
}, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "getEmergencyWithdrawal", user)

	outstruct := new(struct {
		Amount      *big.Int
		RequestTime *big.Int
		Executed    bool
		UnlockTime  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.RequestTime = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Executed = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.UnlockTime = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetEmergencyWithdrawal is a free data retrieval call binding the contract method 0x9ac406a5.
//
// Solidity: function getEmergencyWithdrawal(address user) view returns(uint256 amount, uint256 requestTime, bool executed, uint256 unlockTime)
func (_CollateralVaultL1 *CollateralVaultL1Session) GetEmergencyWithdrawal(user common.Address) (struct {
	Amount      *big.Int
	RequestTime *big.Int
	Executed    bool
	UnlockTime  *big.Int
}, error) {
	return _CollateralVaultL1.Contract.GetEmergencyWithdrawal(&_CollateralVaultL1.CallOpts, user)
}

// GetEmergencyWithdrawal is a free data retrieval call binding the contract method 0x9ac406a5.
//
// Solidity: function getEmergencyWithdrawal(address user) view returns(uint256 amount, uint256 requestTime, bool executed, uint256 unlockTime)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) GetEmergencyWithdrawal(user common.Address) (struct {
	Amount      *big.Int
	RequestTime *big.Int
	Executed    bool
	UnlockTime  *big.Int
}, error) {
	return _CollateralVaultL1.Contract.GetEmergencyWithdrawal(&_CollateralVaultL1.CallOpts, user)
}

// GetLockedCollateral is a free data retrieval call binding the contract method 0xd49ab22c.
//
// Solidity: function getLockedCollateral(address user) view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Caller) GetLockedCollateral(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "getLockedCollateral", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockedCollateral is a free data retrieval call binding the contract method 0xd49ab22c.
//
// Solidity: function getLockedCollateral(address user) view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Session) GetLockedCollateral(user common.Address) (*big.Int, error) {
	return _CollateralVaultL1.Contract.GetLockedCollateral(&_CollateralVaultL1.CallOpts, user)
}

// GetLockedCollateral is a free data retrieval call binding the contract method 0xd49ab22c.
//
// Solidity: function getLockedCollateral(address user) view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) GetLockedCollateral(user common.Address) (*big.Int, error) {
	return _CollateralVaultL1.Contract.GetLockedCollateral(&_CollateralVaultL1.CallOpts, user)
}

// GetRemainingDailyLockCapacity is a free data retrieval call binding the contract method 0x034bfb81.
//
// Solidity: function getRemainingDailyLockCapacity() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Caller) GetRemainingDailyLockCapacity(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "getRemainingDailyLockCapacity")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRemainingDailyLockCapacity is a free data retrieval call binding the contract method 0x034bfb81.
//
// Solidity: function getRemainingDailyLockCapacity() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Session) GetRemainingDailyLockCapacity() (*big.Int, error) {
	return _CollateralVaultL1.Contract.GetRemainingDailyLockCapacity(&_CollateralVaultL1.CallOpts)
}

// GetRemainingDailyLockCapacity is a free data retrieval call binding the contract method 0x034bfb81.
//
// Solidity: function getRemainingDailyLockCapacity() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) GetRemainingDailyLockCapacity() (*big.Int, error) {
	return _CollateralVaultL1.Contract.GetRemainingDailyLockCapacity(&_CollateralVaultL1.CallOpts)
}

// GetTimeUntilLimitReset is a free data retrieval call binding the contract method 0xb3aac012.
//
// Solidity: function getTimeUntilLimitReset() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Caller) GetTimeUntilLimitReset(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "getTimeUntilLimitReset")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimeUntilLimitReset is a free data retrieval call binding the contract method 0xb3aac012.
//
// Solidity: function getTimeUntilLimitReset() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Session) GetTimeUntilLimitReset() (*big.Int, error) {
	return _CollateralVaultL1.Contract.GetTimeUntilLimitReset(&_CollateralVaultL1.CallOpts)
}

// GetTimeUntilLimitReset is a free data retrieval call binding the contract method 0xb3aac012.
//
// Solidity: function getTimeUntilLimitReset() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) GetTimeUntilLimitReset() (*big.Int, error) {
	return _CollateralVaultL1.Contract.GetTimeUntilLimitReset(&_CollateralVaultL1.CallOpts)
}

// GetTotalLocked is a free data retrieval call binding the contract method 0xf4732da6.
//
// Solidity: function getTotalLocked() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Caller) GetTotalLocked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "getTotalLocked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalLocked is a free data retrieval call binding the contract method 0xf4732da6.
//
// Solidity: function getTotalLocked() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Session) GetTotalLocked() (*big.Int, error) {
	return _CollateralVaultL1.Contract.GetTotalLocked(&_CollateralVaultL1.CallOpts)
}

// GetTotalLocked is a free data retrieval call binding the contract method 0xf4732da6.
//
// Solidity: function getTotalLocked() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) GetTotalLocked() (*big.Int, error) {
	return _CollateralVaultL1.Contract.GetTotalLocked(&_CollateralVaultL1.CallOpts)
}

// GetVaultStats is a free data retrieval call binding the contract method 0xa59aa5a6.
//
// Solidity: function getVaultStats() view returns(uint256 _totalLocked, uint256 contractBalance)
func (_CollateralVaultL1 *CollateralVaultL1Caller) GetVaultStats(opts *bind.CallOpts) (struct {
	TotalLocked     *big.Int
	ContractBalance *big.Int
}, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "getVaultStats")

	outstruct := new(struct {
		TotalLocked     *big.Int
		ContractBalance *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalLocked = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ContractBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetVaultStats is a free data retrieval call binding the contract method 0xa59aa5a6.
//
// Solidity: function getVaultStats() view returns(uint256 _totalLocked, uint256 contractBalance)
func (_CollateralVaultL1 *CollateralVaultL1Session) GetVaultStats() (struct {
	TotalLocked     *big.Int
	ContractBalance *big.Int
}, error) {
	return _CollateralVaultL1.Contract.GetVaultStats(&_CollateralVaultL1.CallOpts)
}

// GetVaultStats is a free data retrieval call binding the contract method 0xa59aa5a6.
//
// Solidity: function getVaultStats() view returns(uint256 _totalLocked, uint256 contractBalance)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) GetVaultStats() (struct {
	TotalLocked     *big.Int
	ContractBalance *big.Int
}, error) {
	return _CollateralVaultL1.Contract.GetVaultStats(&_CollateralVaultL1.CallOpts)
}

// L2Bridge is a free data retrieval call binding the contract method 0xae1f6aaf.
//
// Solidity: function l2Bridge() view returns(address)
func (_CollateralVaultL1 *CollateralVaultL1Caller) L2Bridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "l2Bridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2Bridge is a free data retrieval call binding the contract method 0xae1f6aaf.
//
// Solidity: function l2Bridge() view returns(address)
func (_CollateralVaultL1 *CollateralVaultL1Session) L2Bridge() (common.Address, error) {
	return _CollateralVaultL1.Contract.L2Bridge(&_CollateralVaultL1.CallOpts)
}

// L2Bridge is a free data retrieval call binding the contract method 0xae1f6aaf.
//
// Solidity: function l2Bridge() view returns(address)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) L2Bridge() (common.Address, error) {
	return _CollateralVaultL1.Contract.L2Bridge(&_CollateralVaultL1.CallOpts)
}

// LastLockResetTime is a free data retrieval call binding the contract method 0xe39e5e0b.
//
// Solidity: function lastLockResetTime() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Caller) LastLockResetTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "lastLockResetTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastLockResetTime is a free data retrieval call binding the contract method 0xe39e5e0b.
//
// Solidity: function lastLockResetTime() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Session) LastLockResetTime() (*big.Int, error) {
	return _CollateralVaultL1.Contract.LastLockResetTime(&_CollateralVaultL1.CallOpts)
}

// LastLockResetTime is a free data retrieval call binding the contract method 0xe39e5e0b.
//
// Solidity: function lastLockResetTime() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) LastLockResetTime() (*big.Int, error) {
	return _CollateralVaultL1.Contract.LastLockResetTime(&_CollateralVaultL1.CallOpts)
}

// LockedCollateral is a free data retrieval call binding the contract method 0x92bdf9ba.
//
// Solidity: function lockedCollateral(address ) view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Caller) LockedCollateral(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "lockedCollateral", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockedCollateral is a free data retrieval call binding the contract method 0x92bdf9ba.
//
// Solidity: function lockedCollateral(address ) view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Session) LockedCollateral(arg0 common.Address) (*big.Int, error) {
	return _CollateralVaultL1.Contract.LockedCollateral(&_CollateralVaultL1.CallOpts, arg0)
}

// LockedCollateral is a free data retrieval call binding the contract method 0x92bdf9ba.
//
// Solidity: function lockedCollateral(address ) view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) LockedCollateral(arg0 common.Address) (*big.Int, error) {
	return _CollateralVaultL1.Contract.LockedCollateral(&_CollateralVaultL1.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CollateralVaultL1 *CollateralVaultL1Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CollateralVaultL1 *CollateralVaultL1Session) Owner() (common.Address, error) {
	return _CollateralVaultL1.Contract.Owner(&_CollateralVaultL1.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) Owner() (common.Address, error) {
	return _CollateralVaultL1.Contract.Owner(&_CollateralVaultL1.CallOpts)
}

// StateRegistry is a free data retrieval call binding the contract method 0x85ff240e.
//
// Solidity: function stateRegistry() view returns(address)
func (_CollateralVaultL1 *CollateralVaultL1Caller) StateRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "stateRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StateRegistry is a free data retrieval call binding the contract method 0x85ff240e.
//
// Solidity: function stateRegistry() view returns(address)
func (_CollateralVaultL1 *CollateralVaultL1Session) StateRegistry() (common.Address, error) {
	return _CollateralVaultL1.Contract.StateRegistry(&_CollateralVaultL1.CallOpts)
}

// StateRegistry is a free data retrieval call binding the contract method 0x85ff240e.
//
// Solidity: function stateRegistry() view returns(address)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) StateRegistry() (common.Address, error) {
	return _CollateralVaultL1.Contract.StateRegistry(&_CollateralVaultL1.CallOpts)
}

// TotalLocked is a free data retrieval call binding the contract method 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Caller) TotalLocked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CollateralVaultL1.contract.Call(opts, &out, "totalLocked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalLocked is a free data retrieval call binding the contract method 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1Session) TotalLocked() (*big.Int, error) {
	return _CollateralVaultL1.Contract.TotalLocked(&_CollateralVaultL1.CallOpts)
}

// TotalLocked is a free data retrieval call binding the contract method 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (_CollateralVaultL1 *CollateralVaultL1CallerSession) TotalLocked() (*big.Int, error) {
	return _CollateralVaultL1.Contract.TotalLocked(&_CollateralVaultL1.CallOpts)
}

// ExecuteEmergencyWithdrawal is a paid mutator transaction binding the contract method 0xf4993bbd.
//
// Solidity: function executeEmergencyWithdrawal() returns()
func (_CollateralVaultL1 *CollateralVaultL1Transactor) ExecuteEmergencyWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralVaultL1.contract.Transact(opts, "executeEmergencyWithdrawal")
}

// ExecuteEmergencyWithdrawal is a paid mutator transaction binding the contract method 0xf4993bbd.
//
// Solidity: function executeEmergencyWithdrawal() returns()
func (_CollateralVaultL1 *CollateralVaultL1Session) ExecuteEmergencyWithdrawal() (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.ExecuteEmergencyWithdrawal(&_CollateralVaultL1.TransactOpts)
}

// ExecuteEmergencyWithdrawal is a paid mutator transaction binding the contract method 0xf4993bbd.
//
// Solidity: function executeEmergencyWithdrawal() returns()
func (_CollateralVaultL1 *CollateralVaultL1TransactorSession) ExecuteEmergencyWithdrawal() (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.ExecuteEmergencyWithdrawal(&_CollateralVaultL1.TransactOpts)
}

// LockCollateral is a paid mutator transaction binding the contract method 0x6b7c6d20.
//
// Solidity: function lockCollateral(address user, uint256 amount, bytes32 l2TxHash) returns()
func (_CollateralVaultL1 *CollateralVaultL1Transactor) LockCollateral(opts *bind.TransactOpts, user common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _CollateralVaultL1.contract.Transact(opts, "lockCollateral", user, amount, l2TxHash)
}

// LockCollateral is a paid mutator transaction binding the contract method 0x6b7c6d20.
//
// Solidity: function lockCollateral(address user, uint256 amount, bytes32 l2TxHash) returns()
func (_CollateralVaultL1 *CollateralVaultL1Session) LockCollateral(user common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.LockCollateral(&_CollateralVaultL1.TransactOpts, user, amount, l2TxHash)
}

// LockCollateral is a paid mutator transaction binding the contract method 0x6b7c6d20.
//
// Solidity: function lockCollateral(address user, uint256 amount, bytes32 l2TxHash) returns()
func (_CollateralVaultL1 *CollateralVaultL1TransactorSession) LockCollateral(user common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.LockCollateral(&_CollateralVaultL1.TransactOpts, user, amount, l2TxHash)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CollateralVaultL1 *CollateralVaultL1Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralVaultL1.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CollateralVaultL1 *CollateralVaultL1Session) RenounceOwnership() (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.RenounceOwnership(&_CollateralVaultL1.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CollateralVaultL1 *CollateralVaultL1TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.RenounceOwnership(&_CollateralVaultL1.TransactOpts)
}

// RequestEmergencyWithdrawal is a paid mutator transaction binding the contract method 0x23c7b49c.
//
// Solidity: function requestEmergencyWithdrawal(uint256 amount) returns()
func (_CollateralVaultL1 *CollateralVaultL1Transactor) RequestEmergencyWithdrawal(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _CollateralVaultL1.contract.Transact(opts, "requestEmergencyWithdrawal", amount)
}

// RequestEmergencyWithdrawal is a paid mutator transaction binding the contract method 0x23c7b49c.
//
// Solidity: function requestEmergencyWithdrawal(uint256 amount) returns()
func (_CollateralVaultL1 *CollateralVaultL1Session) RequestEmergencyWithdrawal(amount *big.Int) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.RequestEmergencyWithdrawal(&_CollateralVaultL1.TransactOpts, amount)
}

// RequestEmergencyWithdrawal is a paid mutator transaction binding the contract method 0x23c7b49c.
//
// Solidity: function requestEmergencyWithdrawal(uint256 amount) returns()
func (_CollateralVaultL1 *CollateralVaultL1TransactorSession) RequestEmergencyWithdrawal(amount *big.Int) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.RequestEmergencyWithdrawal(&_CollateralVaultL1.TransactOpts, amount)
}

// ResumeFromEmergency is a paid mutator transaction binding the contract method 0x4f206d7f.
//
// Solidity: function resumeFromEmergency() returns()
func (_CollateralVaultL1 *CollateralVaultL1Transactor) ResumeFromEmergency(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralVaultL1.contract.Transact(opts, "resumeFromEmergency")
}

// ResumeFromEmergency is a paid mutator transaction binding the contract method 0x4f206d7f.
//
// Solidity: function resumeFromEmergency() returns()
func (_CollateralVaultL1 *CollateralVaultL1Session) ResumeFromEmergency() (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.ResumeFromEmergency(&_CollateralVaultL1.TransactOpts)
}

// ResumeFromEmergency is a paid mutator transaction binding the contract method 0x4f206d7f.
//
// Solidity: function resumeFromEmergency() returns()
func (_CollateralVaultL1 *CollateralVaultL1TransactorSession) ResumeFromEmergency() (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.ResumeFromEmergency(&_CollateralVaultL1.TransactOpts)
}

// SetDailyLockLimit is a paid mutator transaction binding the contract method 0x0619cc15.
//
// Solidity: function setDailyLockLimit(uint256 newLimit) returns()
func (_CollateralVaultL1 *CollateralVaultL1Transactor) SetDailyLockLimit(opts *bind.TransactOpts, newLimit *big.Int) (*types.Transaction, error) {
	return _CollateralVaultL1.contract.Transact(opts, "setDailyLockLimit", newLimit)
}

// SetDailyLockLimit is a paid mutator transaction binding the contract method 0x0619cc15.
//
// Solidity: function setDailyLockLimit(uint256 newLimit) returns()
func (_CollateralVaultL1 *CollateralVaultL1Session) SetDailyLockLimit(newLimit *big.Int) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.SetDailyLockLimit(&_CollateralVaultL1.TransactOpts, newLimit)
}

// SetDailyLockLimit is a paid mutator transaction binding the contract method 0x0619cc15.
//
// Solidity: function setDailyLockLimit(uint256 newLimit) returns()
func (_CollateralVaultL1 *CollateralVaultL1TransactorSession) SetDailyLockLimit(newLimit *big.Int) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.SetDailyLockLimit(&_CollateralVaultL1.TransactOpts, newLimit)
}

// SetL2Bridge is a paid mutator transaction binding the contract method 0x3d36d971.
//
// Solidity: function setL2Bridge(address _l2Bridge) returns()
func (_CollateralVaultL1 *CollateralVaultL1Transactor) SetL2Bridge(opts *bind.TransactOpts, _l2Bridge common.Address) (*types.Transaction, error) {
	return _CollateralVaultL1.contract.Transact(opts, "setL2Bridge", _l2Bridge)
}

// SetL2Bridge is a paid mutator transaction binding the contract method 0x3d36d971.
//
// Solidity: function setL2Bridge(address _l2Bridge) returns()
func (_CollateralVaultL1 *CollateralVaultL1Session) SetL2Bridge(_l2Bridge common.Address) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.SetL2Bridge(&_CollateralVaultL1.TransactOpts, _l2Bridge)
}

// SetL2Bridge is a paid mutator transaction binding the contract method 0x3d36d971.
//
// Solidity: function setL2Bridge(address _l2Bridge) returns()
func (_CollateralVaultL1 *CollateralVaultL1TransactorSession) SetL2Bridge(_l2Bridge common.Address) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.SetL2Bridge(&_CollateralVaultL1.TransactOpts, _l2Bridge)
}

// SetStateRegistry is a paid mutator transaction binding the contract method 0xa396602a.
//
// Solidity: function setStateRegistry(address _stateRegistry) returns()
func (_CollateralVaultL1 *CollateralVaultL1Transactor) SetStateRegistry(opts *bind.TransactOpts, _stateRegistry common.Address) (*types.Transaction, error) {
	return _CollateralVaultL1.contract.Transact(opts, "setStateRegistry", _stateRegistry)
}

// SetStateRegistry is a paid mutator transaction binding the contract method 0xa396602a.
//
// Solidity: function setStateRegistry(address _stateRegistry) returns()
func (_CollateralVaultL1 *CollateralVaultL1Session) SetStateRegistry(_stateRegistry common.Address) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.SetStateRegistry(&_CollateralVaultL1.TransactOpts, _stateRegistry)
}

// SetStateRegistry is a paid mutator transaction binding the contract method 0xa396602a.
//
// Solidity: function setStateRegistry(address _stateRegistry) returns()
func (_CollateralVaultL1 *CollateralVaultL1TransactorSession) SetStateRegistry(_stateRegistry common.Address) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.SetStateRegistry(&_CollateralVaultL1.TransactOpts, _stateRegistry)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CollateralVaultL1 *CollateralVaultL1Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CollateralVaultL1.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CollateralVaultL1 *CollateralVaultL1Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.TransferOwnership(&_CollateralVaultL1.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CollateralVaultL1 *CollateralVaultL1TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.TransferOwnership(&_CollateralVaultL1.TransactOpts, newOwner)
}

// TriggerEmergencyPause is a paid mutator transaction binding the contract method 0xab523298.
//
// Solidity: function triggerEmergencyPause() returns()
func (_CollateralVaultL1 *CollateralVaultL1Transactor) TriggerEmergencyPause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralVaultL1.contract.Transact(opts, "triggerEmergencyPause")
}

// TriggerEmergencyPause is a paid mutator transaction binding the contract method 0xab523298.
//
// Solidity: function triggerEmergencyPause() returns()
func (_CollateralVaultL1 *CollateralVaultL1Session) TriggerEmergencyPause() (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.TriggerEmergencyPause(&_CollateralVaultL1.TransactOpts)
}

// TriggerEmergencyPause is a paid mutator transaction binding the contract method 0xab523298.
//
// Solidity: function triggerEmergencyPause() returns()
func (_CollateralVaultL1 *CollateralVaultL1TransactorSession) TriggerEmergencyPause() (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.TriggerEmergencyPause(&_CollateralVaultL1.TransactOpts)
}

// UnlockCollateral is a paid mutator transaction binding the contract method 0x139a61ff.
//
// Solidity: function unlockCollateral(address user, uint256 amount, bytes32 l2TxHash) returns()
func (_CollateralVaultL1 *CollateralVaultL1Transactor) UnlockCollateral(opts *bind.TransactOpts, user common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _CollateralVaultL1.contract.Transact(opts, "unlockCollateral", user, amount, l2TxHash)
}

// UnlockCollateral is a paid mutator transaction binding the contract method 0x139a61ff.
//
// Solidity: function unlockCollateral(address user, uint256 amount, bytes32 l2TxHash) returns()
func (_CollateralVaultL1 *CollateralVaultL1Session) UnlockCollateral(user common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.UnlockCollateral(&_CollateralVaultL1.TransactOpts, user, amount, l2TxHash)
}

// UnlockCollateral is a paid mutator transaction binding the contract method 0x139a61ff.
//
// Solidity: function unlockCollateral(address user, uint256 amount, bytes32 l2TxHash) returns()
func (_CollateralVaultL1 *CollateralVaultL1TransactorSession) UnlockCollateral(user common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _CollateralVaultL1.Contract.UnlockCollateral(&_CollateralVaultL1.TransactOpts, user, amount, l2TxHash)
}

// CollateralVaultL1CollateralLockedIterator is returned from FilterCollateralLocked and is used to iterate over the raw logs and unpacked data for CollateralLocked events raised by the CollateralVaultL1 contract.
type CollateralVaultL1CollateralLockedIterator struct {
	Event *CollateralVaultL1CollateralLocked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CollateralVaultL1CollateralLockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralVaultL1CollateralLocked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CollateralVaultL1CollateralLocked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CollateralVaultL1CollateralLockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralVaultL1CollateralLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralVaultL1CollateralLocked represents a CollateralLocked event raised by the CollateralVaultL1 contract.
type CollateralVaultL1CollateralLocked struct {
	User            common.Address
	Amount          *big.Int
	TotalUserLocked *big.Int
	L2TxHash        [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCollateralLocked is a free log retrieval operation binding the contract event 0x6f55f454f279aaa9e234e6da6843a322398a7ee40229ee9db56109504a079a3c.
//
// Solidity: event CollateralLocked(address indexed user, uint256 amount, uint256 totalUserLocked, bytes32 indexed l2TxHash)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) FilterCollateralLocked(opts *bind.FilterOpts, user []common.Address, l2TxHash [][32]byte) (*CollateralVaultL1CollateralLockedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var l2TxHashRule []interface{}
	for _, l2TxHashItem := range l2TxHash {
		l2TxHashRule = append(l2TxHashRule, l2TxHashItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.FilterLogs(opts, "CollateralLocked", userRule, l2TxHashRule)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1CollateralLockedIterator{contract: _CollateralVaultL1.contract, event: "CollateralLocked", logs: logs, sub: sub}, nil
}

// WatchCollateralLocked is a free log subscription operation binding the contract event 0x6f55f454f279aaa9e234e6da6843a322398a7ee40229ee9db56109504a079a3c.
//
// Solidity: event CollateralLocked(address indexed user, uint256 amount, uint256 totalUserLocked, bytes32 indexed l2TxHash)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) WatchCollateralLocked(opts *bind.WatchOpts, sink chan<- *CollateralVaultL1CollateralLocked, user []common.Address, l2TxHash [][32]byte) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var l2TxHashRule []interface{}
	for _, l2TxHashItem := range l2TxHash {
		l2TxHashRule = append(l2TxHashRule, l2TxHashItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.WatchLogs(opts, "CollateralLocked", userRule, l2TxHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralVaultL1CollateralLocked)
				if err := _CollateralVaultL1.contract.UnpackLog(event, "CollateralLocked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCollateralLocked is a log parse operation binding the contract event 0x6f55f454f279aaa9e234e6da6843a322398a7ee40229ee9db56109504a079a3c.
//
// Solidity: event CollateralLocked(address indexed user, uint256 amount, uint256 totalUserLocked, bytes32 indexed l2TxHash)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) ParseCollateralLocked(log types.Log) (*CollateralVaultL1CollateralLocked, error) {
	event := new(CollateralVaultL1CollateralLocked)
	if err := _CollateralVaultL1.contract.UnpackLog(event, "CollateralLocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralVaultL1CollateralUnlockedIterator is returned from FilterCollateralUnlocked and is used to iterate over the raw logs and unpacked data for CollateralUnlocked events raised by the CollateralVaultL1 contract.
type CollateralVaultL1CollateralUnlockedIterator struct {
	Event *CollateralVaultL1CollateralUnlocked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CollateralVaultL1CollateralUnlockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralVaultL1CollateralUnlocked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CollateralVaultL1CollateralUnlocked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CollateralVaultL1CollateralUnlockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralVaultL1CollateralUnlockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralVaultL1CollateralUnlocked represents a CollateralUnlocked event raised by the CollateralVaultL1 contract.
type CollateralVaultL1CollateralUnlocked struct {
	User      common.Address
	Amount    *big.Int
	Remaining *big.Int
	L2TxHash  [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCollateralUnlocked is a free log retrieval operation binding the contract event 0xe6fea673ec5891fd03be95f2657436619789223c68c840fd2fdc88bb3bbe35e9.
//
// Solidity: event CollateralUnlocked(address indexed user, uint256 amount, uint256 remaining, bytes32 indexed l2TxHash)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) FilterCollateralUnlocked(opts *bind.FilterOpts, user []common.Address, l2TxHash [][32]byte) (*CollateralVaultL1CollateralUnlockedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var l2TxHashRule []interface{}
	for _, l2TxHashItem := range l2TxHash {
		l2TxHashRule = append(l2TxHashRule, l2TxHashItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.FilterLogs(opts, "CollateralUnlocked", userRule, l2TxHashRule)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1CollateralUnlockedIterator{contract: _CollateralVaultL1.contract, event: "CollateralUnlocked", logs: logs, sub: sub}, nil
}

// WatchCollateralUnlocked is a free log subscription operation binding the contract event 0xe6fea673ec5891fd03be95f2657436619789223c68c840fd2fdc88bb3bbe35e9.
//
// Solidity: event CollateralUnlocked(address indexed user, uint256 amount, uint256 remaining, bytes32 indexed l2TxHash)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) WatchCollateralUnlocked(opts *bind.WatchOpts, sink chan<- *CollateralVaultL1CollateralUnlocked, user []common.Address, l2TxHash [][32]byte) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var l2TxHashRule []interface{}
	for _, l2TxHashItem := range l2TxHash {
		l2TxHashRule = append(l2TxHashRule, l2TxHashItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.WatchLogs(opts, "CollateralUnlocked", userRule, l2TxHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralVaultL1CollateralUnlocked)
				if err := _CollateralVaultL1.contract.UnpackLog(event, "CollateralUnlocked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCollateralUnlocked is a log parse operation binding the contract event 0xe6fea673ec5891fd03be95f2657436619789223c68c840fd2fdc88bb3bbe35e9.
//
// Solidity: event CollateralUnlocked(address indexed user, uint256 amount, uint256 remaining, bytes32 indexed l2TxHash)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) ParseCollateralUnlocked(log types.Log) (*CollateralVaultL1CollateralUnlocked, error) {
	event := new(CollateralVaultL1CollateralUnlocked)
	if err := _CollateralVaultL1.contract.UnpackLog(event, "CollateralUnlocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralVaultL1EmergencyPauseTriggeredIterator is returned from FilterEmergencyPauseTriggered and is used to iterate over the raw logs and unpacked data for EmergencyPauseTriggered events raised by the CollateralVaultL1 contract.
type CollateralVaultL1EmergencyPauseTriggeredIterator struct {
	Event *CollateralVaultL1EmergencyPauseTriggered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CollateralVaultL1EmergencyPauseTriggeredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralVaultL1EmergencyPauseTriggered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CollateralVaultL1EmergencyPauseTriggered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CollateralVaultL1EmergencyPauseTriggeredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralVaultL1EmergencyPauseTriggeredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralVaultL1EmergencyPauseTriggered represents a EmergencyPauseTriggered event raised by the CollateralVaultL1 contract.
type CollateralVaultL1EmergencyPauseTriggered struct {
	TriggeredBy common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterEmergencyPauseTriggered is a free log retrieval operation binding the contract event 0x7fc2170197efe2877ae0df3c17bce1285ce7a01eadade9718118a394495e7222.
//
// Solidity: event EmergencyPauseTriggered(address indexed triggeredBy)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) FilterEmergencyPauseTriggered(opts *bind.FilterOpts, triggeredBy []common.Address) (*CollateralVaultL1EmergencyPauseTriggeredIterator, error) {

	var triggeredByRule []interface{}
	for _, triggeredByItem := range triggeredBy {
		triggeredByRule = append(triggeredByRule, triggeredByItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.FilterLogs(opts, "EmergencyPauseTriggered", triggeredByRule)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1EmergencyPauseTriggeredIterator{contract: _CollateralVaultL1.contract, event: "EmergencyPauseTriggered", logs: logs, sub: sub}, nil
}

// WatchEmergencyPauseTriggered is a free log subscription operation binding the contract event 0x7fc2170197efe2877ae0df3c17bce1285ce7a01eadade9718118a394495e7222.
//
// Solidity: event EmergencyPauseTriggered(address indexed triggeredBy)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) WatchEmergencyPauseTriggered(opts *bind.WatchOpts, sink chan<- *CollateralVaultL1EmergencyPauseTriggered, triggeredBy []common.Address) (event.Subscription, error) {

	var triggeredByRule []interface{}
	for _, triggeredByItem := range triggeredBy {
		triggeredByRule = append(triggeredByRule, triggeredByItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.WatchLogs(opts, "EmergencyPauseTriggered", triggeredByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralVaultL1EmergencyPauseTriggered)
				if err := _CollateralVaultL1.contract.UnpackLog(event, "EmergencyPauseTriggered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEmergencyPauseTriggered is a log parse operation binding the contract event 0x7fc2170197efe2877ae0df3c17bce1285ce7a01eadade9718118a394495e7222.
//
// Solidity: event EmergencyPauseTriggered(address indexed triggeredBy)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) ParseEmergencyPauseTriggered(log types.Log) (*CollateralVaultL1EmergencyPauseTriggered, error) {
	event := new(CollateralVaultL1EmergencyPauseTriggered)
	if err := _CollateralVaultL1.contract.UnpackLog(event, "EmergencyPauseTriggered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralVaultL1EmergencyResumedIterator is returned from FilterEmergencyResumed and is used to iterate over the raw logs and unpacked data for EmergencyResumed events raised by the CollateralVaultL1 contract.
type CollateralVaultL1EmergencyResumedIterator struct {
	Event *CollateralVaultL1EmergencyResumed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CollateralVaultL1EmergencyResumedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralVaultL1EmergencyResumed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CollateralVaultL1EmergencyResumed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CollateralVaultL1EmergencyResumedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralVaultL1EmergencyResumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralVaultL1EmergencyResumed represents a EmergencyResumed event raised by the CollateralVaultL1 contract.
type CollateralVaultL1EmergencyResumed struct {
	ResumedBy common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEmergencyResumed is a free log retrieval operation binding the contract event 0x1e9c7db3ca9c0db24e98d98c1065c14eb6e9b0eeaacb3326360b12789ccebbb7.
//
// Solidity: event EmergencyResumed(address indexed resumedBy)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) FilterEmergencyResumed(opts *bind.FilterOpts, resumedBy []common.Address) (*CollateralVaultL1EmergencyResumedIterator, error) {

	var resumedByRule []interface{}
	for _, resumedByItem := range resumedBy {
		resumedByRule = append(resumedByRule, resumedByItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.FilterLogs(opts, "EmergencyResumed", resumedByRule)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1EmergencyResumedIterator{contract: _CollateralVaultL1.contract, event: "EmergencyResumed", logs: logs, sub: sub}, nil
}

// WatchEmergencyResumed is a free log subscription operation binding the contract event 0x1e9c7db3ca9c0db24e98d98c1065c14eb6e9b0eeaacb3326360b12789ccebbb7.
//
// Solidity: event EmergencyResumed(address indexed resumedBy)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) WatchEmergencyResumed(opts *bind.WatchOpts, sink chan<- *CollateralVaultL1EmergencyResumed, resumedBy []common.Address) (event.Subscription, error) {

	var resumedByRule []interface{}
	for _, resumedByItem := range resumedBy {
		resumedByRule = append(resumedByRule, resumedByItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.WatchLogs(opts, "EmergencyResumed", resumedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralVaultL1EmergencyResumed)
				if err := _CollateralVaultL1.contract.UnpackLog(event, "EmergencyResumed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEmergencyResumed is a log parse operation binding the contract event 0x1e9c7db3ca9c0db24e98d98c1065c14eb6e9b0eeaacb3326360b12789ccebbb7.
//
// Solidity: event EmergencyResumed(address indexed resumedBy)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) ParseEmergencyResumed(log types.Log) (*CollateralVaultL1EmergencyResumed, error) {
	event := new(CollateralVaultL1EmergencyResumed)
	if err := _CollateralVaultL1.contract.UnpackLog(event, "EmergencyResumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralVaultL1EmergencyWithdrawalExecutedIterator is returned from FilterEmergencyWithdrawalExecuted and is used to iterate over the raw logs and unpacked data for EmergencyWithdrawalExecuted events raised by the CollateralVaultL1 contract.
type CollateralVaultL1EmergencyWithdrawalExecutedIterator struct {
	Event *CollateralVaultL1EmergencyWithdrawalExecuted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CollateralVaultL1EmergencyWithdrawalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralVaultL1EmergencyWithdrawalExecuted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CollateralVaultL1EmergencyWithdrawalExecuted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CollateralVaultL1EmergencyWithdrawalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralVaultL1EmergencyWithdrawalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralVaultL1EmergencyWithdrawalExecuted represents a EmergencyWithdrawalExecuted event raised by the CollateralVaultL1 contract.
type CollateralVaultL1EmergencyWithdrawalExecuted struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEmergencyWithdrawalExecuted is a free log retrieval operation binding the contract event 0xaf6e7d750e5a786a0b4d0e00806a5f778d09769d9c8c0ee377cc6df72af07115.
//
// Solidity: event EmergencyWithdrawalExecuted(address indexed user, uint256 amount)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) FilterEmergencyWithdrawalExecuted(opts *bind.FilterOpts, user []common.Address) (*CollateralVaultL1EmergencyWithdrawalExecutedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.FilterLogs(opts, "EmergencyWithdrawalExecuted", userRule)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1EmergencyWithdrawalExecutedIterator{contract: _CollateralVaultL1.contract, event: "EmergencyWithdrawalExecuted", logs: logs, sub: sub}, nil
}

// WatchEmergencyWithdrawalExecuted is a free log subscription operation binding the contract event 0xaf6e7d750e5a786a0b4d0e00806a5f778d09769d9c8c0ee377cc6df72af07115.
//
// Solidity: event EmergencyWithdrawalExecuted(address indexed user, uint256 amount)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) WatchEmergencyWithdrawalExecuted(opts *bind.WatchOpts, sink chan<- *CollateralVaultL1EmergencyWithdrawalExecuted, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.WatchLogs(opts, "EmergencyWithdrawalExecuted", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralVaultL1EmergencyWithdrawalExecuted)
				if err := _CollateralVaultL1.contract.UnpackLog(event, "EmergencyWithdrawalExecuted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEmergencyWithdrawalExecuted is a log parse operation binding the contract event 0xaf6e7d750e5a786a0b4d0e00806a5f778d09769d9c8c0ee377cc6df72af07115.
//
// Solidity: event EmergencyWithdrawalExecuted(address indexed user, uint256 amount)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) ParseEmergencyWithdrawalExecuted(log types.Log) (*CollateralVaultL1EmergencyWithdrawalExecuted, error) {
	event := new(CollateralVaultL1EmergencyWithdrawalExecuted)
	if err := _CollateralVaultL1.contract.UnpackLog(event, "EmergencyWithdrawalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralVaultL1EmergencyWithdrawalRequestedIterator is returned from FilterEmergencyWithdrawalRequested and is used to iterate over the raw logs and unpacked data for EmergencyWithdrawalRequested events raised by the CollateralVaultL1 contract.
type CollateralVaultL1EmergencyWithdrawalRequestedIterator struct {
	Event *CollateralVaultL1EmergencyWithdrawalRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CollateralVaultL1EmergencyWithdrawalRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralVaultL1EmergencyWithdrawalRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CollateralVaultL1EmergencyWithdrawalRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CollateralVaultL1EmergencyWithdrawalRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralVaultL1EmergencyWithdrawalRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralVaultL1EmergencyWithdrawalRequested represents a EmergencyWithdrawalRequested event raised by the CollateralVaultL1 contract.
type CollateralVaultL1EmergencyWithdrawalRequested struct {
	User       common.Address
	Amount     *big.Int
	UnlockTime *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterEmergencyWithdrawalRequested is a free log retrieval operation binding the contract event 0x6055b9e008637eb8e9ab4faa67bf39b6e34013ab63b9df2935891398a66b213c.
//
// Solidity: event EmergencyWithdrawalRequested(address indexed user, uint256 amount, uint256 unlockTime)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) FilterEmergencyWithdrawalRequested(opts *bind.FilterOpts, user []common.Address) (*CollateralVaultL1EmergencyWithdrawalRequestedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.FilterLogs(opts, "EmergencyWithdrawalRequested", userRule)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1EmergencyWithdrawalRequestedIterator{contract: _CollateralVaultL1.contract, event: "EmergencyWithdrawalRequested", logs: logs, sub: sub}, nil
}

// WatchEmergencyWithdrawalRequested is a free log subscription operation binding the contract event 0x6055b9e008637eb8e9ab4faa67bf39b6e34013ab63b9df2935891398a66b213c.
//
// Solidity: event EmergencyWithdrawalRequested(address indexed user, uint256 amount, uint256 unlockTime)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) WatchEmergencyWithdrawalRequested(opts *bind.WatchOpts, sink chan<- *CollateralVaultL1EmergencyWithdrawalRequested, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.WatchLogs(opts, "EmergencyWithdrawalRequested", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralVaultL1EmergencyWithdrawalRequested)
				if err := _CollateralVaultL1.contract.UnpackLog(event, "EmergencyWithdrawalRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEmergencyWithdrawalRequested is a log parse operation binding the contract event 0x6055b9e008637eb8e9ab4faa67bf39b6e34013ab63b9df2935891398a66b213c.
//
// Solidity: event EmergencyWithdrawalRequested(address indexed user, uint256 amount, uint256 unlockTime)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) ParseEmergencyWithdrawalRequested(log types.Log) (*CollateralVaultL1EmergencyWithdrawalRequested, error) {
	event := new(CollateralVaultL1EmergencyWithdrawalRequested)
	if err := _CollateralVaultL1.contract.UnpackLog(event, "EmergencyWithdrawalRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralVaultL1L2BridgeUpdatedIterator is returned from FilterL2BridgeUpdated and is used to iterate over the raw logs and unpacked data for L2BridgeUpdated events raised by the CollateralVaultL1 contract.
type CollateralVaultL1L2BridgeUpdatedIterator struct {
	Event *CollateralVaultL1L2BridgeUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CollateralVaultL1L2BridgeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralVaultL1L2BridgeUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CollateralVaultL1L2BridgeUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CollateralVaultL1L2BridgeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralVaultL1L2BridgeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralVaultL1L2BridgeUpdated represents a L2BridgeUpdated event raised by the CollateralVaultL1 contract.
type CollateralVaultL1L2BridgeUpdated struct {
	OldBridge common.Address
	NewBridge common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterL2BridgeUpdated is a free log retrieval operation binding the contract event 0xfbeaf70015f6ad12833bc6726f03af9bb91a7ee52827c6413c6c250c77854c95.
//
// Solidity: event L2BridgeUpdated(address indexed oldBridge, address indexed newBridge)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) FilterL2BridgeUpdated(opts *bind.FilterOpts, oldBridge []common.Address, newBridge []common.Address) (*CollateralVaultL1L2BridgeUpdatedIterator, error) {

	var oldBridgeRule []interface{}
	for _, oldBridgeItem := range oldBridge {
		oldBridgeRule = append(oldBridgeRule, oldBridgeItem)
	}
	var newBridgeRule []interface{}
	for _, newBridgeItem := range newBridge {
		newBridgeRule = append(newBridgeRule, newBridgeItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.FilterLogs(opts, "L2BridgeUpdated", oldBridgeRule, newBridgeRule)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1L2BridgeUpdatedIterator{contract: _CollateralVaultL1.contract, event: "L2BridgeUpdated", logs: logs, sub: sub}, nil
}

// WatchL2BridgeUpdated is a free log subscription operation binding the contract event 0xfbeaf70015f6ad12833bc6726f03af9bb91a7ee52827c6413c6c250c77854c95.
//
// Solidity: event L2BridgeUpdated(address indexed oldBridge, address indexed newBridge)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) WatchL2BridgeUpdated(opts *bind.WatchOpts, sink chan<- *CollateralVaultL1L2BridgeUpdated, oldBridge []common.Address, newBridge []common.Address) (event.Subscription, error) {

	var oldBridgeRule []interface{}
	for _, oldBridgeItem := range oldBridge {
		oldBridgeRule = append(oldBridgeRule, oldBridgeItem)
	}
	var newBridgeRule []interface{}
	for _, newBridgeItem := range newBridge {
		newBridgeRule = append(newBridgeRule, newBridgeItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.WatchLogs(opts, "L2BridgeUpdated", oldBridgeRule, newBridgeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralVaultL1L2BridgeUpdated)
				if err := _CollateralVaultL1.contract.UnpackLog(event, "L2BridgeUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseL2BridgeUpdated is a log parse operation binding the contract event 0xfbeaf70015f6ad12833bc6726f03af9bb91a7ee52827c6413c6c250c77854c95.
//
// Solidity: event L2BridgeUpdated(address indexed oldBridge, address indexed newBridge)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) ParseL2BridgeUpdated(log types.Log) (*CollateralVaultL1L2BridgeUpdated, error) {
	event := new(CollateralVaultL1L2BridgeUpdated)
	if err := _CollateralVaultL1.contract.UnpackLog(event, "L2BridgeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralVaultL1OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CollateralVaultL1 contract.
type CollateralVaultL1OwnershipTransferredIterator struct {
	Event *CollateralVaultL1OwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CollateralVaultL1OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralVaultL1OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CollateralVaultL1OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CollateralVaultL1OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralVaultL1OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralVaultL1OwnershipTransferred represents a OwnershipTransferred event raised by the CollateralVaultL1 contract.
type CollateralVaultL1OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CollateralVaultL1OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1OwnershipTransferredIterator{contract: _CollateralVaultL1.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CollateralVaultL1OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralVaultL1OwnershipTransferred)
				if err := _CollateralVaultL1.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) ParseOwnershipTransferred(log types.Log) (*CollateralVaultL1OwnershipTransferred, error) {
	event := new(CollateralVaultL1OwnershipTransferred)
	if err := _CollateralVaultL1.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralVaultL1StateRegistryUpdatedIterator is returned from FilterStateRegistryUpdated and is used to iterate over the raw logs and unpacked data for StateRegistryUpdated events raised by the CollateralVaultL1 contract.
type CollateralVaultL1StateRegistryUpdatedIterator struct {
	Event *CollateralVaultL1StateRegistryUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CollateralVaultL1StateRegistryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralVaultL1StateRegistryUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CollateralVaultL1StateRegistryUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CollateralVaultL1StateRegistryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralVaultL1StateRegistryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralVaultL1StateRegistryUpdated represents a StateRegistryUpdated event raised by the CollateralVaultL1 contract.
type CollateralVaultL1StateRegistryUpdated struct {
	OldRegistry common.Address
	NewRegistry common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterStateRegistryUpdated is a free log retrieval operation binding the contract event 0x343bc515543dc5468d9b7d9d18f0ba746a1d4d51b8fbaf8dbcc3a810f8d0e82d.
//
// Solidity: event StateRegistryUpdated(address indexed oldRegistry, address indexed newRegistry)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) FilterStateRegistryUpdated(opts *bind.FilterOpts, oldRegistry []common.Address, newRegistry []common.Address) (*CollateralVaultL1StateRegistryUpdatedIterator, error) {

	var oldRegistryRule []interface{}
	for _, oldRegistryItem := range oldRegistry {
		oldRegistryRule = append(oldRegistryRule, oldRegistryItem)
	}
	var newRegistryRule []interface{}
	for _, newRegistryItem := range newRegistry {
		newRegistryRule = append(newRegistryRule, newRegistryItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.FilterLogs(opts, "StateRegistryUpdated", oldRegistryRule, newRegistryRule)
	if err != nil {
		return nil, err
	}
	return &CollateralVaultL1StateRegistryUpdatedIterator{contract: _CollateralVaultL1.contract, event: "StateRegistryUpdated", logs: logs, sub: sub}, nil
}

// WatchStateRegistryUpdated is a free log subscription operation binding the contract event 0x343bc515543dc5468d9b7d9d18f0ba746a1d4d51b8fbaf8dbcc3a810f8d0e82d.
//
// Solidity: event StateRegistryUpdated(address indexed oldRegistry, address indexed newRegistry)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) WatchStateRegistryUpdated(opts *bind.WatchOpts, sink chan<- *CollateralVaultL1StateRegistryUpdated, oldRegistry []common.Address, newRegistry []common.Address) (event.Subscription, error) {

	var oldRegistryRule []interface{}
	for _, oldRegistryItem := range oldRegistry {
		oldRegistryRule = append(oldRegistryRule, oldRegistryItem)
	}
	var newRegistryRule []interface{}
	for _, newRegistryItem := range newRegistry {
		newRegistryRule = append(newRegistryRule, newRegistryItem)
	}

	logs, sub, err := _CollateralVaultL1.contract.WatchLogs(opts, "StateRegistryUpdated", oldRegistryRule, newRegistryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralVaultL1StateRegistryUpdated)
				if err := _CollateralVaultL1.contract.UnpackLog(event, "StateRegistryUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStateRegistryUpdated is a log parse operation binding the contract event 0x343bc515543dc5468d9b7d9d18f0ba746a1d4d51b8fbaf8dbcc3a810f8d0e82d.
//
// Solidity: event StateRegistryUpdated(address indexed oldRegistry, address indexed newRegistry)
func (_CollateralVaultL1 *CollateralVaultL1Filterer) ParseStateRegistryUpdated(log types.Log) (*CollateralVaultL1StateRegistryUpdated, error) {
	event := new(CollateralVaultL1StateRegistryUpdated)
	if err := _CollateralVaultL1.contract.UnpackLog(event, "StateRegistryUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
