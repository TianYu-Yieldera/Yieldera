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

// L1StateRegistryStateSnapshot is an auto generated low-level Go binding around an user-defined struct.
type L1StateRegistryStateSnapshot struct {
	StateRoot         [32]byte
	L2BlockNumber     *big.Int
	Timestamp         *big.Int
	TotalCollateral   *big.Int
	TotalDebt         *big.Int
	CriticalCondition bool
}

// L1StateRegistryMetaData contains all meta data concerning the L1StateRegistry contract.
var L1StateRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l2StateAggregator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalDebt\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ratio\",\"type\":\"uint256\"}],\"name\":\"CriticalConditionDetected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"EmergencyExitInitiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"triggeredBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"EmergencyPauseTriggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldAggregator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAggregator\",\"type\":\"address\"}],\"name\":\"L2AggregatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"l2Block\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"StateRootReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minCollateralRatio\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"liquidationThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxDebtCeiling\",\"type\":\"uint256\"}],\"name\":\"ThresholdsUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MIN_SUBMISSION_INTERVAL\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"canSubmit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"emergencyExitAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"emergencyExitClaimed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"emergencyPaused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestState\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"l2Block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getSnapshot\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"l2BlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebt\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"criticalCondition\",\"type\":\"bool\"}],\"internalType\":\"structL1StateRegistry.StateSnapshot\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSnapshotCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"l2Block\",\"type\":\"uint256\"}],\"name\":\"getStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"initiateEmergencyExit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isStateFresh\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2StateAggregator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastSubmissionTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestL2Block\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"l2Block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebt\",\"type\":\"uint256\"}],\"name\":\"receiveStateRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resumeFromEmergency\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAggregator\",\"type\":\"address\"}],\"name\":\"setL2Aggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"snapshots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"l2BlockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebt\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"criticalCondition\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stateRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stateTimestamps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"thresholds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"minCollateralRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxDebtCeiling\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeSinceLastSubmission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"triggerEmergencyPause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minCollateralRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_liquidationThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxDebtCeiling\",\"type\":\"uint256\"}],\"name\":\"updateThresholds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"l2Block\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"expectedRoot\",\"type\":\"bytes32\"}],\"name\":\"verifyStateRoot\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// L1StateRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use L1StateRegistryMetaData.ABI instead.
var L1StateRegistryABI = L1StateRegistryMetaData.ABI

// L1StateRegistry is an auto generated Go binding around an Ethereum contract.
type L1StateRegistry struct {
	L1StateRegistryCaller     // Read-only binding to the contract
	L1StateRegistryTransactor // Write-only binding to the contract
	L1StateRegistryFilterer   // Log filterer for contract events
}

// L1StateRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type L1StateRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1StateRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L1StateRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1StateRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L1StateRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1StateRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L1StateRegistrySession struct {
	Contract     *L1StateRegistry  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L1StateRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L1StateRegistryCallerSession struct {
	Contract *L1StateRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// L1StateRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L1StateRegistryTransactorSession struct {
	Contract     *L1StateRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// L1StateRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type L1StateRegistryRaw struct {
	Contract *L1StateRegistry // Generic contract binding to access the raw methods on
}

// L1StateRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L1StateRegistryCallerRaw struct {
	Contract *L1StateRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// L1StateRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L1StateRegistryTransactorRaw struct {
	Contract *L1StateRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL1StateRegistry creates a new instance of L1StateRegistry, bound to a specific deployed contract.
func NewL1StateRegistry(address common.Address, backend bind.ContractBackend) (*L1StateRegistry, error) {
	contract, err := bindL1StateRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L1StateRegistry{L1StateRegistryCaller: L1StateRegistryCaller{contract: contract}, L1StateRegistryTransactor: L1StateRegistryTransactor{contract: contract}, L1StateRegistryFilterer: L1StateRegistryFilterer{contract: contract}}, nil
}

// NewL1StateRegistryCaller creates a new read-only instance of L1StateRegistry, bound to a specific deployed contract.
func NewL1StateRegistryCaller(address common.Address, caller bind.ContractCaller) (*L1StateRegistryCaller, error) {
	contract, err := bindL1StateRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L1StateRegistryCaller{contract: contract}, nil
}

// NewL1StateRegistryTransactor creates a new write-only instance of L1StateRegistry, bound to a specific deployed contract.
func NewL1StateRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*L1StateRegistryTransactor, error) {
	contract, err := bindL1StateRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L1StateRegistryTransactor{contract: contract}, nil
}

// NewL1StateRegistryFilterer creates a new log filterer instance of L1StateRegistry, bound to a specific deployed contract.
func NewL1StateRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*L1StateRegistryFilterer, error) {
	contract, err := bindL1StateRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L1StateRegistryFilterer{contract: contract}, nil
}

// bindL1StateRegistry binds a generic wrapper to an already deployed contract.
func bindL1StateRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L1StateRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1StateRegistry *L1StateRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1StateRegistry.Contract.L1StateRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1StateRegistry *L1StateRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.L1StateRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1StateRegistry *L1StateRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.L1StateRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1StateRegistry *L1StateRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1StateRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1StateRegistry *L1StateRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1StateRegistry *L1StateRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.contract.Transact(opts, method, params...)
}

// MINSUBMISSIONINTERVAL is a free data retrieval call binding the contract method 0xbfefb795.
//
// Solidity: function MIN_SUBMISSION_INTERVAL() view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCaller) MINSUBMISSIONINTERVAL(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "MIN_SUBMISSION_INTERVAL")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINSUBMISSIONINTERVAL is a free data retrieval call binding the contract method 0xbfefb795.
//
// Solidity: function MIN_SUBMISSION_INTERVAL() view returns(uint256)
func (_L1StateRegistry *L1StateRegistrySession) MINSUBMISSIONINTERVAL() (*big.Int, error) {
	return _L1StateRegistry.Contract.MINSUBMISSIONINTERVAL(&_L1StateRegistry.CallOpts)
}

// MINSUBMISSIONINTERVAL is a free data retrieval call binding the contract method 0xbfefb795.
//
// Solidity: function MIN_SUBMISSION_INTERVAL() view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCallerSession) MINSUBMISSIONINTERVAL() (*big.Int, error) {
	return _L1StateRegistry.Contract.MINSUBMISSIONINTERVAL(&_L1StateRegistry.CallOpts)
}

// CanSubmit is a free data retrieval call binding the contract method 0xffbc9bd0.
//
// Solidity: function canSubmit() view returns(bool)
func (_L1StateRegistry *L1StateRegistryCaller) CanSubmit(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "canSubmit")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanSubmit is a free data retrieval call binding the contract method 0xffbc9bd0.
//
// Solidity: function canSubmit() view returns(bool)
func (_L1StateRegistry *L1StateRegistrySession) CanSubmit() (bool, error) {
	return _L1StateRegistry.Contract.CanSubmit(&_L1StateRegistry.CallOpts)
}

// CanSubmit is a free data retrieval call binding the contract method 0xffbc9bd0.
//
// Solidity: function canSubmit() view returns(bool)
func (_L1StateRegistry *L1StateRegistryCallerSession) CanSubmit() (bool, error) {
	return _L1StateRegistry.Contract.CanSubmit(&_L1StateRegistry.CallOpts)
}

// EmergencyExitAmount is a free data retrieval call binding the contract method 0x0590b966.
//
// Solidity: function emergencyExitAmount(address ) view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCaller) EmergencyExitAmount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "emergencyExitAmount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EmergencyExitAmount is a free data retrieval call binding the contract method 0x0590b966.
//
// Solidity: function emergencyExitAmount(address ) view returns(uint256)
func (_L1StateRegistry *L1StateRegistrySession) EmergencyExitAmount(arg0 common.Address) (*big.Int, error) {
	return _L1StateRegistry.Contract.EmergencyExitAmount(&_L1StateRegistry.CallOpts, arg0)
}

// EmergencyExitAmount is a free data retrieval call binding the contract method 0x0590b966.
//
// Solidity: function emergencyExitAmount(address ) view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCallerSession) EmergencyExitAmount(arg0 common.Address) (*big.Int, error) {
	return _L1StateRegistry.Contract.EmergencyExitAmount(&_L1StateRegistry.CallOpts, arg0)
}

// EmergencyExitClaimed is a free data retrieval call binding the contract method 0x06e01a42.
//
// Solidity: function emergencyExitClaimed(address ) view returns(bool)
func (_L1StateRegistry *L1StateRegistryCaller) EmergencyExitClaimed(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "emergencyExitClaimed", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// EmergencyExitClaimed is a free data retrieval call binding the contract method 0x06e01a42.
//
// Solidity: function emergencyExitClaimed(address ) view returns(bool)
func (_L1StateRegistry *L1StateRegistrySession) EmergencyExitClaimed(arg0 common.Address) (bool, error) {
	return _L1StateRegistry.Contract.EmergencyExitClaimed(&_L1StateRegistry.CallOpts, arg0)
}

// EmergencyExitClaimed is a free data retrieval call binding the contract method 0x06e01a42.
//
// Solidity: function emergencyExitClaimed(address ) view returns(bool)
func (_L1StateRegistry *L1StateRegistryCallerSession) EmergencyExitClaimed(arg0 common.Address) (bool, error) {
	return _L1StateRegistry.Contract.EmergencyExitClaimed(&_L1StateRegistry.CallOpts, arg0)
}

// EmergencyPaused is a free data retrieval call binding the contract method 0x27c830a9.
//
// Solidity: function emergencyPaused() view returns(bool)
func (_L1StateRegistry *L1StateRegistryCaller) EmergencyPaused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "emergencyPaused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// EmergencyPaused is a free data retrieval call binding the contract method 0x27c830a9.
//
// Solidity: function emergencyPaused() view returns(bool)
func (_L1StateRegistry *L1StateRegistrySession) EmergencyPaused() (bool, error) {
	return _L1StateRegistry.Contract.EmergencyPaused(&_L1StateRegistry.CallOpts)
}

// EmergencyPaused is a free data retrieval call binding the contract method 0x27c830a9.
//
// Solidity: function emergencyPaused() view returns(bool)
func (_L1StateRegistry *L1StateRegistryCallerSession) EmergencyPaused() (bool, error) {
	return _L1StateRegistry.Contract.EmergencyPaused(&_L1StateRegistry.CallOpts)
}

// GetLatestState is a free data retrieval call binding the contract method 0x2d904eb8.
//
// Solidity: function getLatestState() view returns(bytes32 stateRoot, uint256 l2Block, uint256 timestamp)
func (_L1StateRegistry *L1StateRegistryCaller) GetLatestState(opts *bind.CallOpts) (struct {
	StateRoot [32]byte
	L2Block   *big.Int
	Timestamp *big.Int
}, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "getLatestState")

	outstruct := new(struct {
		StateRoot [32]byte
		L2Block   *big.Int
		Timestamp *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StateRoot = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.L2Block = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetLatestState is a free data retrieval call binding the contract method 0x2d904eb8.
//
// Solidity: function getLatestState() view returns(bytes32 stateRoot, uint256 l2Block, uint256 timestamp)
func (_L1StateRegistry *L1StateRegistrySession) GetLatestState() (struct {
	StateRoot [32]byte
	L2Block   *big.Int
	Timestamp *big.Int
}, error) {
	return _L1StateRegistry.Contract.GetLatestState(&_L1StateRegistry.CallOpts)
}

// GetLatestState is a free data retrieval call binding the contract method 0x2d904eb8.
//
// Solidity: function getLatestState() view returns(bytes32 stateRoot, uint256 l2Block, uint256 timestamp)
func (_L1StateRegistry *L1StateRegistryCallerSession) GetLatestState() (struct {
	StateRoot [32]byte
	L2Block   *big.Int
	Timestamp *big.Int
}, error) {
	return _L1StateRegistry.Contract.GetLatestState(&_L1StateRegistry.CallOpts)
}

// GetSnapshot is a free data retrieval call binding the contract method 0x76f10ad0.
//
// Solidity: function getSnapshot(uint256 index) view returns((bytes32,uint256,uint256,uint256,uint256,bool))
func (_L1StateRegistry *L1StateRegistryCaller) GetSnapshot(opts *bind.CallOpts, index *big.Int) (L1StateRegistryStateSnapshot, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "getSnapshot", index)

	if err != nil {
		return *new(L1StateRegistryStateSnapshot), err
	}

	out0 := *abi.ConvertType(out[0], new(L1StateRegistryStateSnapshot)).(*L1StateRegistryStateSnapshot)

	return out0, err

}

// GetSnapshot is a free data retrieval call binding the contract method 0x76f10ad0.
//
// Solidity: function getSnapshot(uint256 index) view returns((bytes32,uint256,uint256,uint256,uint256,bool))
func (_L1StateRegistry *L1StateRegistrySession) GetSnapshot(index *big.Int) (L1StateRegistryStateSnapshot, error) {
	return _L1StateRegistry.Contract.GetSnapshot(&_L1StateRegistry.CallOpts, index)
}

// GetSnapshot is a free data retrieval call binding the contract method 0x76f10ad0.
//
// Solidity: function getSnapshot(uint256 index) view returns((bytes32,uint256,uint256,uint256,uint256,bool))
func (_L1StateRegistry *L1StateRegistryCallerSession) GetSnapshot(index *big.Int) (L1StateRegistryStateSnapshot, error) {
	return _L1StateRegistry.Contract.GetSnapshot(&_L1StateRegistry.CallOpts, index)
}

// GetSnapshotCount is a free data retrieval call binding the contract method 0xbcbcf148.
//
// Solidity: function getSnapshotCount() view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCaller) GetSnapshotCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "getSnapshotCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSnapshotCount is a free data retrieval call binding the contract method 0xbcbcf148.
//
// Solidity: function getSnapshotCount() view returns(uint256)
func (_L1StateRegistry *L1StateRegistrySession) GetSnapshotCount() (*big.Int, error) {
	return _L1StateRegistry.Contract.GetSnapshotCount(&_L1StateRegistry.CallOpts)
}

// GetSnapshotCount is a free data retrieval call binding the contract method 0xbcbcf148.
//
// Solidity: function getSnapshotCount() view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCallerSession) GetSnapshotCount() (*big.Int, error) {
	return _L1StateRegistry.Contract.GetSnapshotCount(&_L1StateRegistry.CallOpts)
}

// GetStateRoot is a free data retrieval call binding the contract method 0xc3801938.
//
// Solidity: function getStateRoot(uint256 l2Block) view returns(bytes32)
func (_L1StateRegistry *L1StateRegistryCaller) GetStateRoot(opts *bind.CallOpts, l2Block *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "getStateRoot", l2Block)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetStateRoot is a free data retrieval call binding the contract method 0xc3801938.
//
// Solidity: function getStateRoot(uint256 l2Block) view returns(bytes32)
func (_L1StateRegistry *L1StateRegistrySession) GetStateRoot(l2Block *big.Int) ([32]byte, error) {
	return _L1StateRegistry.Contract.GetStateRoot(&_L1StateRegistry.CallOpts, l2Block)
}

// GetStateRoot is a free data retrieval call binding the contract method 0xc3801938.
//
// Solidity: function getStateRoot(uint256 l2Block) view returns(bytes32)
func (_L1StateRegistry *L1StateRegistryCallerSession) GetStateRoot(l2Block *big.Int) ([32]byte, error) {
	return _L1StateRegistry.Contract.GetStateRoot(&_L1StateRegistry.CallOpts, l2Block)
}

// IsStateFresh is a free data retrieval call binding the contract method 0xce727c43.
//
// Solidity: function isStateFresh() view returns(bool)
func (_L1StateRegistry *L1StateRegistryCaller) IsStateFresh(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "isStateFresh")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsStateFresh is a free data retrieval call binding the contract method 0xce727c43.
//
// Solidity: function isStateFresh() view returns(bool)
func (_L1StateRegistry *L1StateRegistrySession) IsStateFresh() (bool, error) {
	return _L1StateRegistry.Contract.IsStateFresh(&_L1StateRegistry.CallOpts)
}

// IsStateFresh is a free data retrieval call binding the contract method 0xce727c43.
//
// Solidity: function isStateFresh() view returns(bool)
func (_L1StateRegistry *L1StateRegistryCallerSession) IsStateFresh() (bool, error) {
	return _L1StateRegistry.Contract.IsStateFresh(&_L1StateRegistry.CallOpts)
}

// L2StateAggregator is a free data retrieval call binding the contract method 0xdebe7b28.
//
// Solidity: function l2StateAggregator() view returns(address)
func (_L1StateRegistry *L1StateRegistryCaller) L2StateAggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "l2StateAggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2StateAggregator is a free data retrieval call binding the contract method 0xdebe7b28.
//
// Solidity: function l2StateAggregator() view returns(address)
func (_L1StateRegistry *L1StateRegistrySession) L2StateAggregator() (common.Address, error) {
	return _L1StateRegistry.Contract.L2StateAggregator(&_L1StateRegistry.CallOpts)
}

// L2StateAggregator is a free data retrieval call binding the contract method 0xdebe7b28.
//
// Solidity: function l2StateAggregator() view returns(address)
func (_L1StateRegistry *L1StateRegistryCallerSession) L2StateAggregator() (common.Address, error) {
	return _L1StateRegistry.Contract.L2StateAggregator(&_L1StateRegistry.CallOpts)
}

// LastSubmissionTime is a free data retrieval call binding the contract method 0x4f70104e.
//
// Solidity: function lastSubmissionTime() view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCaller) LastSubmissionTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "lastSubmissionTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastSubmissionTime is a free data retrieval call binding the contract method 0x4f70104e.
//
// Solidity: function lastSubmissionTime() view returns(uint256)
func (_L1StateRegistry *L1StateRegistrySession) LastSubmissionTime() (*big.Int, error) {
	return _L1StateRegistry.Contract.LastSubmissionTime(&_L1StateRegistry.CallOpts)
}

// LastSubmissionTime is a free data retrieval call binding the contract method 0x4f70104e.
//
// Solidity: function lastSubmissionTime() view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCallerSession) LastSubmissionTime() (*big.Int, error) {
	return _L1StateRegistry.Contract.LastSubmissionTime(&_L1StateRegistry.CallOpts)
}

// LatestL2Block is a free data retrieval call binding the contract method 0x651c11a5.
//
// Solidity: function latestL2Block() view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCaller) LatestL2Block(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "latestL2Block")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestL2Block is a free data retrieval call binding the contract method 0x651c11a5.
//
// Solidity: function latestL2Block() view returns(uint256)
func (_L1StateRegistry *L1StateRegistrySession) LatestL2Block() (*big.Int, error) {
	return _L1StateRegistry.Contract.LatestL2Block(&_L1StateRegistry.CallOpts)
}

// LatestL2Block is a free data retrieval call binding the contract method 0x651c11a5.
//
// Solidity: function latestL2Block() view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCallerSession) LatestL2Block() (*big.Int, error) {
	return _L1StateRegistry.Contract.LatestL2Block(&_L1StateRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1StateRegistry *L1StateRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1StateRegistry *L1StateRegistrySession) Owner() (common.Address, error) {
	return _L1StateRegistry.Contract.Owner(&_L1StateRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1StateRegistry *L1StateRegistryCallerSession) Owner() (common.Address, error) {
	return _L1StateRegistry.Contract.Owner(&_L1StateRegistry.CallOpts)
}

// Snapshots is a free data retrieval call binding the contract method 0xd6565a2d.
//
// Solidity: function snapshots(uint256 ) view returns(bytes32 stateRoot, uint256 l2BlockNumber, uint256 timestamp, uint256 totalCollateral, uint256 totalDebt, bool criticalCondition)
func (_L1StateRegistry *L1StateRegistryCaller) Snapshots(opts *bind.CallOpts, arg0 *big.Int) (struct {
	StateRoot         [32]byte
	L2BlockNumber     *big.Int
	Timestamp         *big.Int
	TotalCollateral   *big.Int
	TotalDebt         *big.Int
	CriticalCondition bool
}, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "snapshots", arg0)

	outstruct := new(struct {
		StateRoot         [32]byte
		L2BlockNumber     *big.Int
		Timestamp         *big.Int
		TotalCollateral   *big.Int
		TotalDebt         *big.Int
		CriticalCondition bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StateRoot = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.L2BlockNumber = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TotalCollateral = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.TotalDebt = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.CriticalCondition = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

// Snapshots is a free data retrieval call binding the contract method 0xd6565a2d.
//
// Solidity: function snapshots(uint256 ) view returns(bytes32 stateRoot, uint256 l2BlockNumber, uint256 timestamp, uint256 totalCollateral, uint256 totalDebt, bool criticalCondition)
func (_L1StateRegistry *L1StateRegistrySession) Snapshots(arg0 *big.Int) (struct {
	StateRoot         [32]byte
	L2BlockNumber     *big.Int
	Timestamp         *big.Int
	TotalCollateral   *big.Int
	TotalDebt         *big.Int
	CriticalCondition bool
}, error) {
	return _L1StateRegistry.Contract.Snapshots(&_L1StateRegistry.CallOpts, arg0)
}

// Snapshots is a free data retrieval call binding the contract method 0xd6565a2d.
//
// Solidity: function snapshots(uint256 ) view returns(bytes32 stateRoot, uint256 l2BlockNumber, uint256 timestamp, uint256 totalCollateral, uint256 totalDebt, bool criticalCondition)
func (_L1StateRegistry *L1StateRegistryCallerSession) Snapshots(arg0 *big.Int) (struct {
	StateRoot         [32]byte
	L2BlockNumber     *big.Int
	Timestamp         *big.Int
	TotalCollateral   *big.Int
	TotalDebt         *big.Int
	CriticalCondition bool
}, error) {
	return _L1StateRegistry.Contract.Snapshots(&_L1StateRegistry.CallOpts, arg0)
}

// StateRoots is a free data retrieval call binding the contract method 0xf4cac30d.
//
// Solidity: function stateRoots(uint256 ) view returns(bytes32)
func (_L1StateRegistry *L1StateRegistryCaller) StateRoots(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "stateRoots", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// StateRoots is a free data retrieval call binding the contract method 0xf4cac30d.
//
// Solidity: function stateRoots(uint256 ) view returns(bytes32)
func (_L1StateRegistry *L1StateRegistrySession) StateRoots(arg0 *big.Int) ([32]byte, error) {
	return _L1StateRegistry.Contract.StateRoots(&_L1StateRegistry.CallOpts, arg0)
}

// StateRoots is a free data retrieval call binding the contract method 0xf4cac30d.
//
// Solidity: function stateRoots(uint256 ) view returns(bytes32)
func (_L1StateRegistry *L1StateRegistryCallerSession) StateRoots(arg0 *big.Int) ([32]byte, error) {
	return _L1StateRegistry.Contract.StateRoots(&_L1StateRegistry.CallOpts, arg0)
}

// StateTimestamps is a free data retrieval call binding the contract method 0x32edf24e.
//
// Solidity: function stateTimestamps(uint256 ) view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCaller) StateTimestamps(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "stateTimestamps", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StateTimestamps is a free data retrieval call binding the contract method 0x32edf24e.
//
// Solidity: function stateTimestamps(uint256 ) view returns(uint256)
func (_L1StateRegistry *L1StateRegistrySession) StateTimestamps(arg0 *big.Int) (*big.Int, error) {
	return _L1StateRegistry.Contract.StateTimestamps(&_L1StateRegistry.CallOpts, arg0)
}

// StateTimestamps is a free data retrieval call binding the contract method 0x32edf24e.
//
// Solidity: function stateTimestamps(uint256 ) view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCallerSession) StateTimestamps(arg0 *big.Int) (*big.Int, error) {
	return _L1StateRegistry.Contract.StateTimestamps(&_L1StateRegistry.CallOpts, arg0)
}

// Thresholds is a free data retrieval call binding the contract method 0xcda43b3a.
//
// Solidity: function thresholds() view returns(uint256 minCollateralRatio, uint256 liquidationThreshold, uint256 maxDebtCeiling)
func (_L1StateRegistry *L1StateRegistryCaller) Thresholds(opts *bind.CallOpts) (struct {
	MinCollateralRatio   *big.Int
	LiquidationThreshold *big.Int
	MaxDebtCeiling       *big.Int
}, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "thresholds")

	outstruct := new(struct {
		MinCollateralRatio   *big.Int
		LiquidationThreshold *big.Int
		MaxDebtCeiling       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MinCollateralRatio = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.LiquidationThreshold = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.MaxDebtCeiling = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Thresholds is a free data retrieval call binding the contract method 0xcda43b3a.
//
// Solidity: function thresholds() view returns(uint256 minCollateralRatio, uint256 liquidationThreshold, uint256 maxDebtCeiling)
func (_L1StateRegistry *L1StateRegistrySession) Thresholds() (struct {
	MinCollateralRatio   *big.Int
	LiquidationThreshold *big.Int
	MaxDebtCeiling       *big.Int
}, error) {
	return _L1StateRegistry.Contract.Thresholds(&_L1StateRegistry.CallOpts)
}

// Thresholds is a free data retrieval call binding the contract method 0xcda43b3a.
//
// Solidity: function thresholds() view returns(uint256 minCollateralRatio, uint256 liquidationThreshold, uint256 maxDebtCeiling)
func (_L1StateRegistry *L1StateRegistryCallerSession) Thresholds() (struct {
	MinCollateralRatio   *big.Int
	LiquidationThreshold *big.Int
	MaxDebtCeiling       *big.Int
}, error) {
	return _L1StateRegistry.Contract.Thresholds(&_L1StateRegistry.CallOpts)
}

// TimeSinceLastSubmission is a free data retrieval call binding the contract method 0xd676c062.
//
// Solidity: function timeSinceLastSubmission() view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCaller) TimeSinceLastSubmission(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "timeSinceLastSubmission")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TimeSinceLastSubmission is a free data retrieval call binding the contract method 0xd676c062.
//
// Solidity: function timeSinceLastSubmission() view returns(uint256)
func (_L1StateRegistry *L1StateRegistrySession) TimeSinceLastSubmission() (*big.Int, error) {
	return _L1StateRegistry.Contract.TimeSinceLastSubmission(&_L1StateRegistry.CallOpts)
}

// TimeSinceLastSubmission is a free data retrieval call binding the contract method 0xd676c062.
//
// Solidity: function timeSinceLastSubmission() view returns(uint256)
func (_L1StateRegistry *L1StateRegistryCallerSession) TimeSinceLastSubmission() (*big.Int, error) {
	return _L1StateRegistry.Contract.TimeSinceLastSubmission(&_L1StateRegistry.CallOpts)
}

// VerifyStateRoot is a free data retrieval call binding the contract method 0x01bcc117.
//
// Solidity: function verifyStateRoot(uint256 l2Block, bytes32 expectedRoot) view returns(bool)
func (_L1StateRegistry *L1StateRegistryCaller) VerifyStateRoot(opts *bind.CallOpts, l2Block *big.Int, expectedRoot [32]byte) (bool, error) {
	var out []interface{}
	err := _L1StateRegistry.contract.Call(opts, &out, "verifyStateRoot", l2Block, expectedRoot)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyStateRoot is a free data retrieval call binding the contract method 0x01bcc117.
//
// Solidity: function verifyStateRoot(uint256 l2Block, bytes32 expectedRoot) view returns(bool)
func (_L1StateRegistry *L1StateRegistrySession) VerifyStateRoot(l2Block *big.Int, expectedRoot [32]byte) (bool, error) {
	return _L1StateRegistry.Contract.VerifyStateRoot(&_L1StateRegistry.CallOpts, l2Block, expectedRoot)
}

// VerifyStateRoot is a free data retrieval call binding the contract method 0x01bcc117.
//
// Solidity: function verifyStateRoot(uint256 l2Block, bytes32 expectedRoot) view returns(bool)
func (_L1StateRegistry *L1StateRegistryCallerSession) VerifyStateRoot(l2Block *big.Int, expectedRoot [32]byte) (bool, error) {
	return _L1StateRegistry.Contract.VerifyStateRoot(&_L1StateRegistry.CallOpts, l2Block, expectedRoot)
}

// InitiateEmergencyExit is a paid mutator transaction binding the contract method 0xf23ae8a1.
//
// Solidity: function initiateEmergencyExit(address user, uint256 amount) returns()
func (_L1StateRegistry *L1StateRegistryTransactor) InitiateEmergencyExit(opts *bind.TransactOpts, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L1StateRegistry.contract.Transact(opts, "initiateEmergencyExit", user, amount)
}

// InitiateEmergencyExit is a paid mutator transaction binding the contract method 0xf23ae8a1.
//
// Solidity: function initiateEmergencyExit(address user, uint256 amount) returns()
func (_L1StateRegistry *L1StateRegistrySession) InitiateEmergencyExit(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.InitiateEmergencyExit(&_L1StateRegistry.TransactOpts, user, amount)
}

// InitiateEmergencyExit is a paid mutator transaction binding the contract method 0xf23ae8a1.
//
// Solidity: function initiateEmergencyExit(address user, uint256 amount) returns()
func (_L1StateRegistry *L1StateRegistryTransactorSession) InitiateEmergencyExit(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.InitiateEmergencyExit(&_L1StateRegistry.TransactOpts, user, amount)
}

// ReceiveStateRoot is a paid mutator transaction binding the contract method 0x74394df7.
//
// Solidity: function receiveStateRoot(bytes32 stateRoot, uint256 l2Block, uint256 totalCollateral, uint256 totalDebt) returns()
func (_L1StateRegistry *L1StateRegistryTransactor) ReceiveStateRoot(opts *bind.TransactOpts, stateRoot [32]byte, l2Block *big.Int, totalCollateral *big.Int, totalDebt *big.Int) (*types.Transaction, error) {
	return _L1StateRegistry.contract.Transact(opts, "receiveStateRoot", stateRoot, l2Block, totalCollateral, totalDebt)
}

// ReceiveStateRoot is a paid mutator transaction binding the contract method 0x74394df7.
//
// Solidity: function receiveStateRoot(bytes32 stateRoot, uint256 l2Block, uint256 totalCollateral, uint256 totalDebt) returns()
func (_L1StateRegistry *L1StateRegistrySession) ReceiveStateRoot(stateRoot [32]byte, l2Block *big.Int, totalCollateral *big.Int, totalDebt *big.Int) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.ReceiveStateRoot(&_L1StateRegistry.TransactOpts, stateRoot, l2Block, totalCollateral, totalDebt)
}

// ReceiveStateRoot is a paid mutator transaction binding the contract method 0x74394df7.
//
// Solidity: function receiveStateRoot(bytes32 stateRoot, uint256 l2Block, uint256 totalCollateral, uint256 totalDebt) returns()
func (_L1StateRegistry *L1StateRegistryTransactorSession) ReceiveStateRoot(stateRoot [32]byte, l2Block *big.Int, totalCollateral *big.Int, totalDebt *big.Int) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.ReceiveStateRoot(&_L1StateRegistry.TransactOpts, stateRoot, l2Block, totalCollateral, totalDebt)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L1StateRegistry *L1StateRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1StateRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L1StateRegistry *L1StateRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _L1StateRegistry.Contract.RenounceOwnership(&_L1StateRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L1StateRegistry *L1StateRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _L1StateRegistry.Contract.RenounceOwnership(&_L1StateRegistry.TransactOpts)
}

// ResumeFromEmergency is a paid mutator transaction binding the contract method 0x4f206d7f.
//
// Solidity: function resumeFromEmergency() returns()
func (_L1StateRegistry *L1StateRegistryTransactor) ResumeFromEmergency(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1StateRegistry.contract.Transact(opts, "resumeFromEmergency")
}

// ResumeFromEmergency is a paid mutator transaction binding the contract method 0x4f206d7f.
//
// Solidity: function resumeFromEmergency() returns()
func (_L1StateRegistry *L1StateRegistrySession) ResumeFromEmergency() (*types.Transaction, error) {
	return _L1StateRegistry.Contract.ResumeFromEmergency(&_L1StateRegistry.TransactOpts)
}

// ResumeFromEmergency is a paid mutator transaction binding the contract method 0x4f206d7f.
//
// Solidity: function resumeFromEmergency() returns()
func (_L1StateRegistry *L1StateRegistryTransactorSession) ResumeFromEmergency() (*types.Transaction, error) {
	return _L1StateRegistry.Contract.ResumeFromEmergency(&_L1StateRegistry.TransactOpts)
}

// SetL2Aggregator is a paid mutator transaction binding the contract method 0x615c765a.
//
// Solidity: function setL2Aggregator(address newAggregator) returns()
func (_L1StateRegistry *L1StateRegistryTransactor) SetL2Aggregator(opts *bind.TransactOpts, newAggregator common.Address) (*types.Transaction, error) {
	return _L1StateRegistry.contract.Transact(opts, "setL2Aggregator", newAggregator)
}

// SetL2Aggregator is a paid mutator transaction binding the contract method 0x615c765a.
//
// Solidity: function setL2Aggregator(address newAggregator) returns()
func (_L1StateRegistry *L1StateRegistrySession) SetL2Aggregator(newAggregator common.Address) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.SetL2Aggregator(&_L1StateRegistry.TransactOpts, newAggregator)
}

// SetL2Aggregator is a paid mutator transaction binding the contract method 0x615c765a.
//
// Solidity: function setL2Aggregator(address newAggregator) returns()
func (_L1StateRegistry *L1StateRegistryTransactorSession) SetL2Aggregator(newAggregator common.Address) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.SetL2Aggregator(&_L1StateRegistry.TransactOpts, newAggregator)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1StateRegistry *L1StateRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _L1StateRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1StateRegistry *L1StateRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.TransferOwnership(&_L1StateRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1StateRegistry *L1StateRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.TransferOwnership(&_L1StateRegistry.TransactOpts, newOwner)
}

// TriggerEmergencyPause is a paid mutator transaction binding the contract method 0xbc4afd37.
//
// Solidity: function triggerEmergencyPause(string reason) returns()
func (_L1StateRegistry *L1StateRegistryTransactor) TriggerEmergencyPause(opts *bind.TransactOpts, reason string) (*types.Transaction, error) {
	return _L1StateRegistry.contract.Transact(opts, "triggerEmergencyPause", reason)
}

// TriggerEmergencyPause is a paid mutator transaction binding the contract method 0xbc4afd37.
//
// Solidity: function triggerEmergencyPause(string reason) returns()
func (_L1StateRegistry *L1StateRegistrySession) TriggerEmergencyPause(reason string) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.TriggerEmergencyPause(&_L1StateRegistry.TransactOpts, reason)
}

// TriggerEmergencyPause is a paid mutator transaction binding the contract method 0xbc4afd37.
//
// Solidity: function triggerEmergencyPause(string reason) returns()
func (_L1StateRegistry *L1StateRegistryTransactorSession) TriggerEmergencyPause(reason string) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.TriggerEmergencyPause(&_L1StateRegistry.TransactOpts, reason)
}

// UpdateThresholds is a paid mutator transaction binding the contract method 0xfd19e10c.
//
// Solidity: function updateThresholds(uint256 _minCollateralRatio, uint256 _liquidationThreshold, uint256 _maxDebtCeiling) returns()
func (_L1StateRegistry *L1StateRegistryTransactor) UpdateThresholds(opts *bind.TransactOpts, _minCollateralRatio *big.Int, _liquidationThreshold *big.Int, _maxDebtCeiling *big.Int) (*types.Transaction, error) {
	return _L1StateRegistry.contract.Transact(opts, "updateThresholds", _minCollateralRatio, _liquidationThreshold, _maxDebtCeiling)
}

// UpdateThresholds is a paid mutator transaction binding the contract method 0xfd19e10c.
//
// Solidity: function updateThresholds(uint256 _minCollateralRatio, uint256 _liquidationThreshold, uint256 _maxDebtCeiling) returns()
func (_L1StateRegistry *L1StateRegistrySession) UpdateThresholds(_minCollateralRatio *big.Int, _liquidationThreshold *big.Int, _maxDebtCeiling *big.Int) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.UpdateThresholds(&_L1StateRegistry.TransactOpts, _minCollateralRatio, _liquidationThreshold, _maxDebtCeiling)
}

// UpdateThresholds is a paid mutator transaction binding the contract method 0xfd19e10c.
//
// Solidity: function updateThresholds(uint256 _minCollateralRatio, uint256 _liquidationThreshold, uint256 _maxDebtCeiling) returns()
func (_L1StateRegistry *L1StateRegistryTransactorSession) UpdateThresholds(_minCollateralRatio *big.Int, _liquidationThreshold *big.Int, _maxDebtCeiling *big.Int) (*types.Transaction, error) {
	return _L1StateRegistry.Contract.UpdateThresholds(&_L1StateRegistry.TransactOpts, _minCollateralRatio, _liquidationThreshold, _maxDebtCeiling)
}

// L1StateRegistryCriticalConditionDetectedIterator is returned from FilterCriticalConditionDetected and is used to iterate over the raw logs and unpacked data for CriticalConditionDetected events raised by the L1StateRegistry contract.
type L1StateRegistryCriticalConditionDetectedIterator struct {
	Event *L1StateRegistryCriticalConditionDetected // Event containing the contract specifics and raw log

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
func (it *L1StateRegistryCriticalConditionDetectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1StateRegistryCriticalConditionDetected)
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
		it.Event = new(L1StateRegistryCriticalConditionDetected)
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
func (it *L1StateRegistryCriticalConditionDetectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1StateRegistryCriticalConditionDetectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1StateRegistryCriticalConditionDetected represents a CriticalConditionDetected event raised by the L1StateRegistry contract.
type L1StateRegistryCriticalConditionDetected struct {
	StateRoot       [32]byte
	TotalCollateral *big.Int
	TotalDebt       *big.Int
	Ratio           *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCriticalConditionDetected is a free log retrieval operation binding the contract event 0xf44c388b67f9789716cda6ea316b888f14bf9b794fbe1c885a856db83e4a2d5e.
//
// Solidity: event CriticalConditionDetected(bytes32 stateRoot, uint256 totalCollateral, uint256 totalDebt, uint256 ratio)
func (_L1StateRegistry *L1StateRegistryFilterer) FilterCriticalConditionDetected(opts *bind.FilterOpts) (*L1StateRegistryCriticalConditionDetectedIterator, error) {

	logs, sub, err := _L1StateRegistry.contract.FilterLogs(opts, "CriticalConditionDetected")
	if err != nil {
		return nil, err
	}
	return &L1StateRegistryCriticalConditionDetectedIterator{contract: _L1StateRegistry.contract, event: "CriticalConditionDetected", logs: logs, sub: sub}, nil
}

// WatchCriticalConditionDetected is a free log subscription operation binding the contract event 0xf44c388b67f9789716cda6ea316b888f14bf9b794fbe1c885a856db83e4a2d5e.
//
// Solidity: event CriticalConditionDetected(bytes32 stateRoot, uint256 totalCollateral, uint256 totalDebt, uint256 ratio)
func (_L1StateRegistry *L1StateRegistryFilterer) WatchCriticalConditionDetected(opts *bind.WatchOpts, sink chan<- *L1StateRegistryCriticalConditionDetected) (event.Subscription, error) {

	logs, sub, err := _L1StateRegistry.contract.WatchLogs(opts, "CriticalConditionDetected")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1StateRegistryCriticalConditionDetected)
				if err := _L1StateRegistry.contract.UnpackLog(event, "CriticalConditionDetected", log); err != nil {
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

// ParseCriticalConditionDetected is a log parse operation binding the contract event 0xf44c388b67f9789716cda6ea316b888f14bf9b794fbe1c885a856db83e4a2d5e.
//
// Solidity: event CriticalConditionDetected(bytes32 stateRoot, uint256 totalCollateral, uint256 totalDebt, uint256 ratio)
func (_L1StateRegistry *L1StateRegistryFilterer) ParseCriticalConditionDetected(log types.Log) (*L1StateRegistryCriticalConditionDetected, error) {
	event := new(L1StateRegistryCriticalConditionDetected)
	if err := _L1StateRegistry.contract.UnpackLog(event, "CriticalConditionDetected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1StateRegistryEmergencyExitInitiatedIterator is returned from FilterEmergencyExitInitiated and is used to iterate over the raw logs and unpacked data for EmergencyExitInitiated events raised by the L1StateRegistry contract.
type L1StateRegistryEmergencyExitInitiatedIterator struct {
	Event *L1StateRegistryEmergencyExitInitiated // Event containing the contract specifics and raw log

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
func (it *L1StateRegistryEmergencyExitInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1StateRegistryEmergencyExitInitiated)
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
		it.Event = new(L1StateRegistryEmergencyExitInitiated)
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
func (it *L1StateRegistryEmergencyExitInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1StateRegistryEmergencyExitInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1StateRegistryEmergencyExitInitiated represents a EmergencyExitInitiated event raised by the L1StateRegistry contract.
type L1StateRegistryEmergencyExitInitiated struct {
	User        common.Address
	Amount      *big.Int
	BlockNumber *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterEmergencyExitInitiated is a free log retrieval operation binding the contract event 0xe8d4f864bd0e9c4978f9c4e397acbebf923f90abbb7c142e650b5e4639e97c69.
//
// Solidity: event EmergencyExitInitiated(address indexed user, uint256 amount, uint256 blockNumber)
func (_L1StateRegistry *L1StateRegistryFilterer) FilterEmergencyExitInitiated(opts *bind.FilterOpts, user []common.Address) (*L1StateRegistryEmergencyExitInitiatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _L1StateRegistry.contract.FilterLogs(opts, "EmergencyExitInitiated", userRule)
	if err != nil {
		return nil, err
	}
	return &L1StateRegistryEmergencyExitInitiatedIterator{contract: _L1StateRegistry.contract, event: "EmergencyExitInitiated", logs: logs, sub: sub}, nil
}

// WatchEmergencyExitInitiated is a free log subscription operation binding the contract event 0xe8d4f864bd0e9c4978f9c4e397acbebf923f90abbb7c142e650b5e4639e97c69.
//
// Solidity: event EmergencyExitInitiated(address indexed user, uint256 amount, uint256 blockNumber)
func (_L1StateRegistry *L1StateRegistryFilterer) WatchEmergencyExitInitiated(opts *bind.WatchOpts, sink chan<- *L1StateRegistryEmergencyExitInitiated, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _L1StateRegistry.contract.WatchLogs(opts, "EmergencyExitInitiated", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1StateRegistryEmergencyExitInitiated)
				if err := _L1StateRegistry.contract.UnpackLog(event, "EmergencyExitInitiated", log); err != nil {
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

// ParseEmergencyExitInitiated is a log parse operation binding the contract event 0xe8d4f864bd0e9c4978f9c4e397acbebf923f90abbb7c142e650b5e4639e97c69.
//
// Solidity: event EmergencyExitInitiated(address indexed user, uint256 amount, uint256 blockNumber)
func (_L1StateRegistry *L1StateRegistryFilterer) ParseEmergencyExitInitiated(log types.Log) (*L1StateRegistryEmergencyExitInitiated, error) {
	event := new(L1StateRegistryEmergencyExitInitiated)
	if err := _L1StateRegistry.contract.UnpackLog(event, "EmergencyExitInitiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1StateRegistryEmergencyPauseTriggeredIterator is returned from FilterEmergencyPauseTriggered and is used to iterate over the raw logs and unpacked data for EmergencyPauseTriggered events raised by the L1StateRegistry contract.
type L1StateRegistryEmergencyPauseTriggeredIterator struct {
	Event *L1StateRegistryEmergencyPauseTriggered // Event containing the contract specifics and raw log

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
func (it *L1StateRegistryEmergencyPauseTriggeredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1StateRegistryEmergencyPauseTriggered)
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
		it.Event = new(L1StateRegistryEmergencyPauseTriggered)
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
func (it *L1StateRegistryEmergencyPauseTriggeredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1StateRegistryEmergencyPauseTriggeredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1StateRegistryEmergencyPauseTriggered represents a EmergencyPauseTriggered event raised by the L1StateRegistry contract.
type L1StateRegistryEmergencyPauseTriggered struct {
	TriggeredBy common.Address
	Reason      string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterEmergencyPauseTriggered is a free log retrieval operation binding the contract event 0x0d5facded7fa89357ce603b9b1f5da16956ec198bd324b4b40594541b3d7945d.
//
// Solidity: event EmergencyPauseTriggered(address indexed triggeredBy, string reason)
func (_L1StateRegistry *L1StateRegistryFilterer) FilterEmergencyPauseTriggered(opts *bind.FilterOpts, triggeredBy []common.Address) (*L1StateRegistryEmergencyPauseTriggeredIterator, error) {

	var triggeredByRule []interface{}
	for _, triggeredByItem := range triggeredBy {
		triggeredByRule = append(triggeredByRule, triggeredByItem)
	}

	logs, sub, err := _L1StateRegistry.contract.FilterLogs(opts, "EmergencyPauseTriggered", triggeredByRule)
	if err != nil {
		return nil, err
	}
	return &L1StateRegistryEmergencyPauseTriggeredIterator{contract: _L1StateRegistry.contract, event: "EmergencyPauseTriggered", logs: logs, sub: sub}, nil
}

// WatchEmergencyPauseTriggered is a free log subscription operation binding the contract event 0x0d5facded7fa89357ce603b9b1f5da16956ec198bd324b4b40594541b3d7945d.
//
// Solidity: event EmergencyPauseTriggered(address indexed triggeredBy, string reason)
func (_L1StateRegistry *L1StateRegistryFilterer) WatchEmergencyPauseTriggered(opts *bind.WatchOpts, sink chan<- *L1StateRegistryEmergencyPauseTriggered, triggeredBy []common.Address) (event.Subscription, error) {

	var triggeredByRule []interface{}
	for _, triggeredByItem := range triggeredBy {
		triggeredByRule = append(triggeredByRule, triggeredByItem)
	}

	logs, sub, err := _L1StateRegistry.contract.WatchLogs(opts, "EmergencyPauseTriggered", triggeredByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1StateRegistryEmergencyPauseTriggered)
				if err := _L1StateRegistry.contract.UnpackLog(event, "EmergencyPauseTriggered", log); err != nil {
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

// ParseEmergencyPauseTriggered is a log parse operation binding the contract event 0x0d5facded7fa89357ce603b9b1f5da16956ec198bd324b4b40594541b3d7945d.
//
// Solidity: event EmergencyPauseTriggered(address indexed triggeredBy, string reason)
func (_L1StateRegistry *L1StateRegistryFilterer) ParseEmergencyPauseTriggered(log types.Log) (*L1StateRegistryEmergencyPauseTriggered, error) {
	event := new(L1StateRegistryEmergencyPauseTriggered)
	if err := _L1StateRegistry.contract.UnpackLog(event, "EmergencyPauseTriggered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1StateRegistryL2AggregatorUpdatedIterator is returned from FilterL2AggregatorUpdated and is used to iterate over the raw logs and unpacked data for L2AggregatorUpdated events raised by the L1StateRegistry contract.
type L1StateRegistryL2AggregatorUpdatedIterator struct {
	Event *L1StateRegistryL2AggregatorUpdated // Event containing the contract specifics and raw log

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
func (it *L1StateRegistryL2AggregatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1StateRegistryL2AggregatorUpdated)
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
		it.Event = new(L1StateRegistryL2AggregatorUpdated)
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
func (it *L1StateRegistryL2AggregatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1StateRegistryL2AggregatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1StateRegistryL2AggregatorUpdated represents a L2AggregatorUpdated event raised by the L1StateRegistry contract.
type L1StateRegistryL2AggregatorUpdated struct {
	OldAggregator common.Address
	NewAggregator common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterL2AggregatorUpdated is a free log retrieval operation binding the contract event 0x6ef6075bc0848b3988d8dc62c5c46c42343ba651284ae88c91f181d57478a1cb.
//
// Solidity: event L2AggregatorUpdated(address indexed oldAggregator, address indexed newAggregator)
func (_L1StateRegistry *L1StateRegistryFilterer) FilterL2AggregatorUpdated(opts *bind.FilterOpts, oldAggregator []common.Address, newAggregator []common.Address) (*L1StateRegistryL2AggregatorUpdatedIterator, error) {

	var oldAggregatorRule []interface{}
	for _, oldAggregatorItem := range oldAggregator {
		oldAggregatorRule = append(oldAggregatorRule, oldAggregatorItem)
	}
	var newAggregatorRule []interface{}
	for _, newAggregatorItem := range newAggregator {
		newAggregatorRule = append(newAggregatorRule, newAggregatorItem)
	}

	logs, sub, err := _L1StateRegistry.contract.FilterLogs(opts, "L2AggregatorUpdated", oldAggregatorRule, newAggregatorRule)
	if err != nil {
		return nil, err
	}
	return &L1StateRegistryL2AggregatorUpdatedIterator{contract: _L1StateRegistry.contract, event: "L2AggregatorUpdated", logs: logs, sub: sub}, nil
}

// WatchL2AggregatorUpdated is a free log subscription operation binding the contract event 0x6ef6075bc0848b3988d8dc62c5c46c42343ba651284ae88c91f181d57478a1cb.
//
// Solidity: event L2AggregatorUpdated(address indexed oldAggregator, address indexed newAggregator)
func (_L1StateRegistry *L1StateRegistryFilterer) WatchL2AggregatorUpdated(opts *bind.WatchOpts, sink chan<- *L1StateRegistryL2AggregatorUpdated, oldAggregator []common.Address, newAggregator []common.Address) (event.Subscription, error) {

	var oldAggregatorRule []interface{}
	for _, oldAggregatorItem := range oldAggregator {
		oldAggregatorRule = append(oldAggregatorRule, oldAggregatorItem)
	}
	var newAggregatorRule []interface{}
	for _, newAggregatorItem := range newAggregator {
		newAggregatorRule = append(newAggregatorRule, newAggregatorItem)
	}

	logs, sub, err := _L1StateRegistry.contract.WatchLogs(opts, "L2AggregatorUpdated", oldAggregatorRule, newAggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1StateRegistryL2AggregatorUpdated)
				if err := _L1StateRegistry.contract.UnpackLog(event, "L2AggregatorUpdated", log); err != nil {
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

// ParseL2AggregatorUpdated is a log parse operation binding the contract event 0x6ef6075bc0848b3988d8dc62c5c46c42343ba651284ae88c91f181d57478a1cb.
//
// Solidity: event L2AggregatorUpdated(address indexed oldAggregator, address indexed newAggregator)
func (_L1StateRegistry *L1StateRegistryFilterer) ParseL2AggregatorUpdated(log types.Log) (*L1StateRegistryL2AggregatorUpdated, error) {
	event := new(L1StateRegistryL2AggregatorUpdated)
	if err := _L1StateRegistry.contract.UnpackLog(event, "L2AggregatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1StateRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the L1StateRegistry contract.
type L1StateRegistryOwnershipTransferredIterator struct {
	Event *L1StateRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *L1StateRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1StateRegistryOwnershipTransferred)
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
		it.Event = new(L1StateRegistryOwnershipTransferred)
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
func (it *L1StateRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1StateRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1StateRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the L1StateRegistry contract.
type L1StateRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1StateRegistry *L1StateRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*L1StateRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1StateRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &L1StateRegistryOwnershipTransferredIterator{contract: _L1StateRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1StateRegistry *L1StateRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *L1StateRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1StateRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1StateRegistryOwnershipTransferred)
				if err := _L1StateRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_L1StateRegistry *L1StateRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*L1StateRegistryOwnershipTransferred, error) {
	event := new(L1StateRegistryOwnershipTransferred)
	if err := _L1StateRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1StateRegistryStateRootReceivedIterator is returned from FilterStateRootReceived and is used to iterate over the raw logs and unpacked data for StateRootReceived events raised by the L1StateRegistry contract.
type L1StateRegistryStateRootReceivedIterator struct {
	Event *L1StateRegistryStateRootReceived // Event containing the contract specifics and raw log

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
func (it *L1StateRegistryStateRootReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1StateRegistryStateRootReceived)
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
		it.Event = new(L1StateRegistryStateRootReceived)
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
func (it *L1StateRegistryStateRootReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1StateRegistryStateRootReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1StateRegistryStateRootReceived represents a StateRootReceived event raised by the L1StateRegistry contract.
type L1StateRegistryStateRootReceived struct {
	StateRoot [32]byte
	L2Block   *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStateRootReceived is a free log retrieval operation binding the contract event 0xcaa33a8c8a0a8bad27ea33f2c94181fe631de8b420f54a94615f5f822ed0f461.
//
// Solidity: event StateRootReceived(bytes32 indexed stateRoot, uint256 indexed l2Block, uint256 timestamp)
func (_L1StateRegistry *L1StateRegistryFilterer) FilterStateRootReceived(opts *bind.FilterOpts, stateRoot [][32]byte, l2Block []*big.Int) (*L1StateRegistryStateRootReceivedIterator, error) {

	var stateRootRule []interface{}
	for _, stateRootItem := range stateRoot {
		stateRootRule = append(stateRootRule, stateRootItem)
	}
	var l2BlockRule []interface{}
	for _, l2BlockItem := range l2Block {
		l2BlockRule = append(l2BlockRule, l2BlockItem)
	}

	logs, sub, err := _L1StateRegistry.contract.FilterLogs(opts, "StateRootReceived", stateRootRule, l2BlockRule)
	if err != nil {
		return nil, err
	}
	return &L1StateRegistryStateRootReceivedIterator{contract: _L1StateRegistry.contract, event: "StateRootReceived", logs: logs, sub: sub}, nil
}

// WatchStateRootReceived is a free log subscription operation binding the contract event 0xcaa33a8c8a0a8bad27ea33f2c94181fe631de8b420f54a94615f5f822ed0f461.
//
// Solidity: event StateRootReceived(bytes32 indexed stateRoot, uint256 indexed l2Block, uint256 timestamp)
func (_L1StateRegistry *L1StateRegistryFilterer) WatchStateRootReceived(opts *bind.WatchOpts, sink chan<- *L1StateRegistryStateRootReceived, stateRoot [][32]byte, l2Block []*big.Int) (event.Subscription, error) {

	var stateRootRule []interface{}
	for _, stateRootItem := range stateRoot {
		stateRootRule = append(stateRootRule, stateRootItem)
	}
	var l2BlockRule []interface{}
	for _, l2BlockItem := range l2Block {
		l2BlockRule = append(l2BlockRule, l2BlockItem)
	}

	logs, sub, err := _L1StateRegistry.contract.WatchLogs(opts, "StateRootReceived", stateRootRule, l2BlockRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1StateRegistryStateRootReceived)
				if err := _L1StateRegistry.contract.UnpackLog(event, "StateRootReceived", log); err != nil {
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

// ParseStateRootReceived is a log parse operation binding the contract event 0xcaa33a8c8a0a8bad27ea33f2c94181fe631de8b420f54a94615f5f822ed0f461.
//
// Solidity: event StateRootReceived(bytes32 indexed stateRoot, uint256 indexed l2Block, uint256 timestamp)
func (_L1StateRegistry *L1StateRegistryFilterer) ParseStateRootReceived(log types.Log) (*L1StateRegistryStateRootReceived, error) {
	event := new(L1StateRegistryStateRootReceived)
	if err := _L1StateRegistry.contract.UnpackLog(event, "StateRootReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1StateRegistryThresholdsUpdatedIterator is returned from FilterThresholdsUpdated and is used to iterate over the raw logs and unpacked data for ThresholdsUpdated events raised by the L1StateRegistry contract.
type L1StateRegistryThresholdsUpdatedIterator struct {
	Event *L1StateRegistryThresholdsUpdated // Event containing the contract specifics and raw log

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
func (it *L1StateRegistryThresholdsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1StateRegistryThresholdsUpdated)
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
		it.Event = new(L1StateRegistryThresholdsUpdated)
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
func (it *L1StateRegistryThresholdsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1StateRegistryThresholdsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1StateRegistryThresholdsUpdated represents a ThresholdsUpdated event raised by the L1StateRegistry contract.
type L1StateRegistryThresholdsUpdated struct {
	MinCollateralRatio   *big.Int
	LiquidationThreshold *big.Int
	MaxDebtCeiling       *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterThresholdsUpdated is a free log retrieval operation binding the contract event 0x5c18dc8d95da80ea715f2473abe6f01199d4bfa87d2ed1ef051058a60dcce258.
//
// Solidity: event ThresholdsUpdated(uint256 minCollateralRatio, uint256 liquidationThreshold, uint256 maxDebtCeiling)
func (_L1StateRegistry *L1StateRegistryFilterer) FilterThresholdsUpdated(opts *bind.FilterOpts) (*L1StateRegistryThresholdsUpdatedIterator, error) {

	logs, sub, err := _L1StateRegistry.contract.FilterLogs(opts, "ThresholdsUpdated")
	if err != nil {
		return nil, err
	}
	return &L1StateRegistryThresholdsUpdatedIterator{contract: _L1StateRegistry.contract, event: "ThresholdsUpdated", logs: logs, sub: sub}, nil
}

// WatchThresholdsUpdated is a free log subscription operation binding the contract event 0x5c18dc8d95da80ea715f2473abe6f01199d4bfa87d2ed1ef051058a60dcce258.
//
// Solidity: event ThresholdsUpdated(uint256 minCollateralRatio, uint256 liquidationThreshold, uint256 maxDebtCeiling)
func (_L1StateRegistry *L1StateRegistryFilterer) WatchThresholdsUpdated(opts *bind.WatchOpts, sink chan<- *L1StateRegistryThresholdsUpdated) (event.Subscription, error) {

	logs, sub, err := _L1StateRegistry.contract.WatchLogs(opts, "ThresholdsUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1StateRegistryThresholdsUpdated)
				if err := _L1StateRegistry.contract.UnpackLog(event, "ThresholdsUpdated", log); err != nil {
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

// ParseThresholdsUpdated is a log parse operation binding the contract event 0x5c18dc8d95da80ea715f2473abe6f01199d4bfa87d2ed1ef051058a60dcce258.
//
// Solidity: event ThresholdsUpdated(uint256 minCollateralRatio, uint256 liquidationThreshold, uint256 maxDebtCeiling)
func (_L1StateRegistry *L1StateRegistryFilterer) ParseThresholdsUpdated(log types.Log) (*L1StateRegistryThresholdsUpdated, error) {
	event := new(L1StateRegistryThresholdsUpdated)
	if err := _L1StateRegistry.contract.UnpackLog(event, "ThresholdsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
