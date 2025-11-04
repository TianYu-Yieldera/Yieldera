// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package l2

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

// L2StateAggregatorMetaData contains all meta data concerning the L2StateAggregator contract.
var L2StateAggregatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l1StateRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldRegistry\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newRegistry\",\"type\":\"address\"}],\"name\":\"L1RegistryUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"moduleId\",\"type\":\"bytes32\"}],\"name\":\"ModuleDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"moduleId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"moduleAddress\",\"type\":\"address\"}],\"name\":\"ModuleRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"moduleId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"updatedBy\",\"type\":\"address\"}],\"name\":\"ModuleStateUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"StateRootComputed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"l2Block\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"l1TxId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalCollateral\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalDebt\",\"type\":\"uint256\"}],\"name\":\"StateSubmittedToL1\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MIN_SUBMISSION_INTERVAL\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"moduleAddress\",\"type\":\"address\"}],\"name\":\"authorizeModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"authorizedModules\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"canSubmitToL1\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"moduleId\",\"type\":\"bytes32\"}],\"name\":\"deactivateModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceSubmitToL1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllModuleIds\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getModuleCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"moduleId\",\"type\":\"bytes32\"}],\"name\":\"getModuleState\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdate\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStateHistoryLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getStateRootByIndex\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSystemState\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"activePositions\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalOrders\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1StateRegistry\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastSubmission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"moduleIds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"moduleStates\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdate\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"moduleId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"moduleAddress\",\"type\":\"address\"}],\"name\":\"registerModule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"moduleAddress\",\"type\":\"address\"}],\"name\":\"revokeModuleAuthorization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l1StateRegistry\",\"type\":\"address\"}],\"name\":\"setL1StateRegistry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stateHistory\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"stateRootToBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"submitToL1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"systemState\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"activePositions\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalOrders\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeUntilNextSubmission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"moduleId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"}],\"name\":\"updateModuleState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"totalCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"activePositions\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalOrders\",\"type\":\"uint256\"}],\"name\":\"updateSystemState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// L2StateAggregatorABI is the input ABI used to generate the binding from.
// Deprecated: Use L2StateAggregatorMetaData.ABI instead.
var L2StateAggregatorABI = L2StateAggregatorMetaData.ABI

// L2StateAggregator is an auto generated Go binding around an Ethereum contract.
type L2StateAggregator struct {
	L2StateAggregatorCaller     // Read-only binding to the contract
	L2StateAggregatorTransactor // Write-only binding to the contract
	L2StateAggregatorFilterer   // Log filterer for contract events
}

// L2StateAggregatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type L2StateAggregatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2StateAggregatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L2StateAggregatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2StateAggregatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L2StateAggregatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2StateAggregatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L2StateAggregatorSession struct {
	Contract     *L2StateAggregator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// L2StateAggregatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L2StateAggregatorCallerSession struct {
	Contract *L2StateAggregatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// L2StateAggregatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L2StateAggregatorTransactorSession struct {
	Contract     *L2StateAggregatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// L2StateAggregatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type L2StateAggregatorRaw struct {
	Contract *L2StateAggregator // Generic contract binding to access the raw methods on
}

// L2StateAggregatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L2StateAggregatorCallerRaw struct {
	Contract *L2StateAggregatorCaller // Generic read-only contract binding to access the raw methods on
}

// L2StateAggregatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L2StateAggregatorTransactorRaw struct {
	Contract *L2StateAggregatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL2StateAggregator creates a new instance of L2StateAggregator, bound to a specific deployed contract.
func NewL2StateAggregator(address common.Address, backend bind.ContractBackend) (*L2StateAggregator, error) {
	contract, err := bindL2StateAggregator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L2StateAggregator{L2StateAggregatorCaller: L2StateAggregatorCaller{contract: contract}, L2StateAggregatorTransactor: L2StateAggregatorTransactor{contract: contract}, L2StateAggregatorFilterer: L2StateAggregatorFilterer{contract: contract}}, nil
}

// NewL2StateAggregatorCaller creates a new read-only instance of L2StateAggregator, bound to a specific deployed contract.
func NewL2StateAggregatorCaller(address common.Address, caller bind.ContractCaller) (*L2StateAggregatorCaller, error) {
	contract, err := bindL2StateAggregator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L2StateAggregatorCaller{contract: contract}, nil
}

// NewL2StateAggregatorTransactor creates a new write-only instance of L2StateAggregator, bound to a specific deployed contract.
func NewL2StateAggregatorTransactor(address common.Address, transactor bind.ContractTransactor) (*L2StateAggregatorTransactor, error) {
	contract, err := bindL2StateAggregator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L2StateAggregatorTransactor{contract: contract}, nil
}

// NewL2StateAggregatorFilterer creates a new log filterer instance of L2StateAggregator, bound to a specific deployed contract.
func NewL2StateAggregatorFilterer(address common.Address, filterer bind.ContractFilterer) (*L2StateAggregatorFilterer, error) {
	contract, err := bindL2StateAggregator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L2StateAggregatorFilterer{contract: contract}, nil
}

// bindL2StateAggregator binds a generic wrapper to an already deployed contract.
func bindL2StateAggregator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L2StateAggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2StateAggregator *L2StateAggregatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2StateAggregator.Contract.L2StateAggregatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2StateAggregator *L2StateAggregatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.L2StateAggregatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2StateAggregator *L2StateAggregatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.L2StateAggregatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2StateAggregator *L2StateAggregatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2StateAggregator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2StateAggregator *L2StateAggregatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2StateAggregator *L2StateAggregatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.contract.Transact(opts, method, params...)
}

// MINSUBMISSIONINTERVAL is a free data retrieval call binding the contract method 0xbfefb795.
//
// Solidity: function MIN_SUBMISSION_INTERVAL() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorCaller) MINSUBMISSIONINTERVAL(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "MIN_SUBMISSION_INTERVAL")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINSUBMISSIONINTERVAL is a free data retrieval call binding the contract method 0xbfefb795.
//
// Solidity: function MIN_SUBMISSION_INTERVAL() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorSession) MINSUBMISSIONINTERVAL() (*big.Int, error) {
	return _L2StateAggregator.Contract.MINSUBMISSIONINTERVAL(&_L2StateAggregator.CallOpts)
}

// MINSUBMISSIONINTERVAL is a free data retrieval call binding the contract method 0xbfefb795.
//
// Solidity: function MIN_SUBMISSION_INTERVAL() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorCallerSession) MINSUBMISSIONINTERVAL() (*big.Int, error) {
	return _L2StateAggregator.Contract.MINSUBMISSIONINTERVAL(&_L2StateAggregator.CallOpts)
}

// AuthorizedModules is a free data retrieval call binding the contract method 0x17f6fe5b.
//
// Solidity: function authorizedModules(address ) view returns(bool)
func (_L2StateAggregator *L2StateAggregatorCaller) AuthorizedModules(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "authorizedModules", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AuthorizedModules is a free data retrieval call binding the contract method 0x17f6fe5b.
//
// Solidity: function authorizedModules(address ) view returns(bool)
func (_L2StateAggregator *L2StateAggregatorSession) AuthorizedModules(arg0 common.Address) (bool, error) {
	return _L2StateAggregator.Contract.AuthorizedModules(&_L2StateAggregator.CallOpts, arg0)
}

// AuthorizedModules is a free data retrieval call binding the contract method 0x17f6fe5b.
//
// Solidity: function authorizedModules(address ) view returns(bool)
func (_L2StateAggregator *L2StateAggregatorCallerSession) AuthorizedModules(arg0 common.Address) (bool, error) {
	return _L2StateAggregator.Contract.AuthorizedModules(&_L2StateAggregator.CallOpts, arg0)
}

// CalculateStateRoot is a free data retrieval call binding the contract method 0x0edae3cb.
//
// Solidity: function calculateStateRoot() view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorCaller) CalculateStateRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "calculateStateRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalculateStateRoot is a free data retrieval call binding the contract method 0x0edae3cb.
//
// Solidity: function calculateStateRoot() view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorSession) CalculateStateRoot() ([32]byte, error) {
	return _L2StateAggregator.Contract.CalculateStateRoot(&_L2StateAggregator.CallOpts)
}

// CalculateStateRoot is a free data retrieval call binding the contract method 0x0edae3cb.
//
// Solidity: function calculateStateRoot() view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorCallerSession) CalculateStateRoot() ([32]byte, error) {
	return _L2StateAggregator.Contract.CalculateStateRoot(&_L2StateAggregator.CallOpts)
}

// CanSubmitToL1 is a free data retrieval call binding the contract method 0x9978bcfa.
//
// Solidity: function canSubmitToL1() view returns(bool)
func (_L2StateAggregator *L2StateAggregatorCaller) CanSubmitToL1(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "canSubmitToL1")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanSubmitToL1 is a free data retrieval call binding the contract method 0x9978bcfa.
//
// Solidity: function canSubmitToL1() view returns(bool)
func (_L2StateAggregator *L2StateAggregatorSession) CanSubmitToL1() (bool, error) {
	return _L2StateAggregator.Contract.CanSubmitToL1(&_L2StateAggregator.CallOpts)
}

// CanSubmitToL1 is a free data retrieval call binding the contract method 0x9978bcfa.
//
// Solidity: function canSubmitToL1() view returns(bool)
func (_L2StateAggregator *L2StateAggregatorCallerSession) CanSubmitToL1() (bool, error) {
	return _L2StateAggregator.Contract.CanSubmitToL1(&_L2StateAggregator.CallOpts)
}

// CurrentStateRoot is a free data retrieval call binding the contract method 0xac2eba98.
//
// Solidity: function currentStateRoot() view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorCaller) CurrentStateRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "currentStateRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CurrentStateRoot is a free data retrieval call binding the contract method 0xac2eba98.
//
// Solidity: function currentStateRoot() view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorSession) CurrentStateRoot() ([32]byte, error) {
	return _L2StateAggregator.Contract.CurrentStateRoot(&_L2StateAggregator.CallOpts)
}

// CurrentStateRoot is a free data retrieval call binding the contract method 0xac2eba98.
//
// Solidity: function currentStateRoot() view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorCallerSession) CurrentStateRoot() ([32]byte, error) {
	return _L2StateAggregator.Contract.CurrentStateRoot(&_L2StateAggregator.CallOpts)
}

// GetAllModuleIds is a free data retrieval call binding the contract method 0xcd869033.
//
// Solidity: function getAllModuleIds() view returns(bytes32[])
func (_L2StateAggregator *L2StateAggregatorCaller) GetAllModuleIds(opts *bind.CallOpts) ([][32]byte, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "getAllModuleIds")

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetAllModuleIds is a free data retrieval call binding the contract method 0xcd869033.
//
// Solidity: function getAllModuleIds() view returns(bytes32[])
func (_L2StateAggregator *L2StateAggregatorSession) GetAllModuleIds() ([][32]byte, error) {
	return _L2StateAggregator.Contract.GetAllModuleIds(&_L2StateAggregator.CallOpts)
}

// GetAllModuleIds is a free data retrieval call binding the contract method 0xcd869033.
//
// Solidity: function getAllModuleIds() view returns(bytes32[])
func (_L2StateAggregator *L2StateAggregatorCallerSession) GetAllModuleIds() ([][32]byte, error) {
	return _L2StateAggregator.Contract.GetAllModuleIds(&_L2StateAggregator.CallOpts)
}

// GetModuleCount is a free data retrieval call binding the contract method 0x57d56267.
//
// Solidity: function getModuleCount() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorCaller) GetModuleCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "getModuleCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetModuleCount is a free data retrieval call binding the contract method 0x57d56267.
//
// Solidity: function getModuleCount() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorSession) GetModuleCount() (*big.Int, error) {
	return _L2StateAggregator.Contract.GetModuleCount(&_L2StateAggregator.CallOpts)
}

// GetModuleCount is a free data retrieval call binding the contract method 0x57d56267.
//
// Solidity: function getModuleCount() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorCallerSession) GetModuleCount() (*big.Int, error) {
	return _L2StateAggregator.Contract.GetModuleCount(&_L2StateAggregator.CallOpts)
}

// GetModuleState is a free data retrieval call binding the contract method 0x8fe08c29.
//
// Solidity: function getModuleState(bytes32 moduleId) view returns(bytes32 stateHash, uint256 lastUpdate, bool active)
func (_L2StateAggregator *L2StateAggregatorCaller) GetModuleState(opts *bind.CallOpts, moduleId [32]byte) (struct {
	StateHash  [32]byte
	LastUpdate *big.Int
	Active     bool
}, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "getModuleState", moduleId)

	outstruct := new(struct {
		StateHash  [32]byte
		LastUpdate *big.Int
		Active     bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StateHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.LastUpdate = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Active = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// GetModuleState is a free data retrieval call binding the contract method 0x8fe08c29.
//
// Solidity: function getModuleState(bytes32 moduleId) view returns(bytes32 stateHash, uint256 lastUpdate, bool active)
func (_L2StateAggregator *L2StateAggregatorSession) GetModuleState(moduleId [32]byte) (struct {
	StateHash  [32]byte
	LastUpdate *big.Int
	Active     bool
}, error) {
	return _L2StateAggregator.Contract.GetModuleState(&_L2StateAggregator.CallOpts, moduleId)
}

// GetModuleState is a free data retrieval call binding the contract method 0x8fe08c29.
//
// Solidity: function getModuleState(bytes32 moduleId) view returns(bytes32 stateHash, uint256 lastUpdate, bool active)
func (_L2StateAggregator *L2StateAggregatorCallerSession) GetModuleState(moduleId [32]byte) (struct {
	StateHash  [32]byte
	LastUpdate *big.Int
	Active     bool
}, error) {
	return _L2StateAggregator.Contract.GetModuleState(&_L2StateAggregator.CallOpts, moduleId)
}

// GetStateHistoryLength is a free data retrieval call binding the contract method 0x1a3472a9.
//
// Solidity: function getStateHistoryLength() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorCaller) GetStateHistoryLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "getStateHistoryLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStateHistoryLength is a free data retrieval call binding the contract method 0x1a3472a9.
//
// Solidity: function getStateHistoryLength() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorSession) GetStateHistoryLength() (*big.Int, error) {
	return _L2StateAggregator.Contract.GetStateHistoryLength(&_L2StateAggregator.CallOpts)
}

// GetStateHistoryLength is a free data retrieval call binding the contract method 0x1a3472a9.
//
// Solidity: function getStateHistoryLength() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorCallerSession) GetStateHistoryLength() (*big.Int, error) {
	return _L2StateAggregator.Contract.GetStateHistoryLength(&_L2StateAggregator.CallOpts)
}

// GetStateRootByIndex is a free data retrieval call binding the contract method 0x3f06a302.
//
// Solidity: function getStateRootByIndex(uint256 index) view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorCaller) GetStateRootByIndex(opts *bind.CallOpts, index *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "getStateRootByIndex", index)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetStateRootByIndex is a free data retrieval call binding the contract method 0x3f06a302.
//
// Solidity: function getStateRootByIndex(uint256 index) view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorSession) GetStateRootByIndex(index *big.Int) ([32]byte, error) {
	return _L2StateAggregator.Contract.GetStateRootByIndex(&_L2StateAggregator.CallOpts, index)
}

// GetStateRootByIndex is a free data retrieval call binding the contract method 0x3f06a302.
//
// Solidity: function getStateRootByIndex(uint256 index) view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorCallerSession) GetStateRootByIndex(index *big.Int) ([32]byte, error) {
	return _L2StateAggregator.Contract.GetStateRootByIndex(&_L2StateAggregator.CallOpts, index)
}

// GetSystemState is a free data retrieval call binding the contract method 0x28fa777f.
//
// Solidity: function getSystemState() view returns(uint256 totalCollateral, uint256 totalDebt, uint256 activePositions, uint256 totalOrders, uint256 blockNumber, uint256 timestamp)
func (_L2StateAggregator *L2StateAggregatorCaller) GetSystemState(opts *bind.CallOpts) (struct {
	TotalCollateral *big.Int
	TotalDebt       *big.Int
	ActivePositions *big.Int
	TotalOrders     *big.Int
	BlockNumber     *big.Int
	Timestamp       *big.Int
}, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "getSystemState")

	outstruct := new(struct {
		TotalCollateral *big.Int
		TotalDebt       *big.Int
		ActivePositions *big.Int
		TotalOrders     *big.Int
		BlockNumber     *big.Int
		Timestamp       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalCollateral = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalDebt = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ActivePositions = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TotalOrders = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.BlockNumber = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetSystemState is a free data retrieval call binding the contract method 0x28fa777f.
//
// Solidity: function getSystemState() view returns(uint256 totalCollateral, uint256 totalDebt, uint256 activePositions, uint256 totalOrders, uint256 blockNumber, uint256 timestamp)
func (_L2StateAggregator *L2StateAggregatorSession) GetSystemState() (struct {
	TotalCollateral *big.Int
	TotalDebt       *big.Int
	ActivePositions *big.Int
	TotalOrders     *big.Int
	BlockNumber     *big.Int
	Timestamp       *big.Int
}, error) {
	return _L2StateAggregator.Contract.GetSystemState(&_L2StateAggregator.CallOpts)
}

// GetSystemState is a free data retrieval call binding the contract method 0x28fa777f.
//
// Solidity: function getSystemState() view returns(uint256 totalCollateral, uint256 totalDebt, uint256 activePositions, uint256 totalOrders, uint256 blockNumber, uint256 timestamp)
func (_L2StateAggregator *L2StateAggregatorCallerSession) GetSystemState() (struct {
	TotalCollateral *big.Int
	TotalDebt       *big.Int
	ActivePositions *big.Int
	TotalOrders     *big.Int
	BlockNumber     *big.Int
	Timestamp       *big.Int
}, error) {
	return _L2StateAggregator.Contract.GetSystemState(&_L2StateAggregator.CallOpts)
}

// L1StateRegistry is a free data retrieval call binding the contract method 0xfbd824e7.
//
// Solidity: function l1StateRegistry() view returns(address)
func (_L2StateAggregator *L2StateAggregatorCaller) L1StateRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "l1StateRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L1StateRegistry is a free data retrieval call binding the contract method 0xfbd824e7.
//
// Solidity: function l1StateRegistry() view returns(address)
func (_L2StateAggregator *L2StateAggregatorSession) L1StateRegistry() (common.Address, error) {
	return _L2StateAggregator.Contract.L1StateRegistry(&_L2StateAggregator.CallOpts)
}

// L1StateRegistry is a free data retrieval call binding the contract method 0xfbd824e7.
//
// Solidity: function l1StateRegistry() view returns(address)
func (_L2StateAggregator *L2StateAggregatorCallerSession) L1StateRegistry() (common.Address, error) {
	return _L2StateAggregator.Contract.L1StateRegistry(&_L2StateAggregator.CallOpts)
}

// LastSubmission is a free data retrieval call binding the contract method 0x8f60482d.
//
// Solidity: function lastSubmission() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorCaller) LastSubmission(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "lastSubmission")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastSubmission is a free data retrieval call binding the contract method 0x8f60482d.
//
// Solidity: function lastSubmission() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorSession) LastSubmission() (*big.Int, error) {
	return _L2StateAggregator.Contract.LastSubmission(&_L2StateAggregator.CallOpts)
}

// LastSubmission is a free data retrieval call binding the contract method 0x8f60482d.
//
// Solidity: function lastSubmission() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorCallerSession) LastSubmission() (*big.Int, error) {
	return _L2StateAggregator.Contract.LastSubmission(&_L2StateAggregator.CallOpts)
}

// ModuleIds is a free data retrieval call binding the contract method 0x5eb66191.
//
// Solidity: function moduleIds(uint256 ) view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorCaller) ModuleIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "moduleIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ModuleIds is a free data retrieval call binding the contract method 0x5eb66191.
//
// Solidity: function moduleIds(uint256 ) view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorSession) ModuleIds(arg0 *big.Int) ([32]byte, error) {
	return _L2StateAggregator.Contract.ModuleIds(&_L2StateAggregator.CallOpts, arg0)
}

// ModuleIds is a free data retrieval call binding the contract method 0x5eb66191.
//
// Solidity: function moduleIds(uint256 ) view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorCallerSession) ModuleIds(arg0 *big.Int) ([32]byte, error) {
	return _L2StateAggregator.Contract.ModuleIds(&_L2StateAggregator.CallOpts, arg0)
}

// ModuleStates is a free data retrieval call binding the contract method 0x6477768b.
//
// Solidity: function moduleStates(bytes32 ) view returns(bytes32 stateHash, uint256 lastUpdate, bool active)
func (_L2StateAggregator *L2StateAggregatorCaller) ModuleStates(opts *bind.CallOpts, arg0 [32]byte) (struct {
	StateHash  [32]byte
	LastUpdate *big.Int
	Active     bool
}, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "moduleStates", arg0)

	outstruct := new(struct {
		StateHash  [32]byte
		LastUpdate *big.Int
		Active     bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StateHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.LastUpdate = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Active = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// ModuleStates is a free data retrieval call binding the contract method 0x6477768b.
//
// Solidity: function moduleStates(bytes32 ) view returns(bytes32 stateHash, uint256 lastUpdate, bool active)
func (_L2StateAggregator *L2StateAggregatorSession) ModuleStates(arg0 [32]byte) (struct {
	StateHash  [32]byte
	LastUpdate *big.Int
	Active     bool
}, error) {
	return _L2StateAggregator.Contract.ModuleStates(&_L2StateAggregator.CallOpts, arg0)
}

// ModuleStates is a free data retrieval call binding the contract method 0x6477768b.
//
// Solidity: function moduleStates(bytes32 ) view returns(bytes32 stateHash, uint256 lastUpdate, bool active)
func (_L2StateAggregator *L2StateAggregatorCallerSession) ModuleStates(arg0 [32]byte) (struct {
	StateHash  [32]byte
	LastUpdate *big.Int
	Active     bool
}, error) {
	return _L2StateAggregator.Contract.ModuleStates(&_L2StateAggregator.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2StateAggregator *L2StateAggregatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2StateAggregator *L2StateAggregatorSession) Owner() (common.Address, error) {
	return _L2StateAggregator.Contract.Owner(&_L2StateAggregator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2StateAggregator *L2StateAggregatorCallerSession) Owner() (common.Address, error) {
	return _L2StateAggregator.Contract.Owner(&_L2StateAggregator.CallOpts)
}

// StateHistory is a free data retrieval call binding the contract method 0x90544984.
//
// Solidity: function stateHistory(uint256 ) view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorCaller) StateHistory(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "stateHistory", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// StateHistory is a free data retrieval call binding the contract method 0x90544984.
//
// Solidity: function stateHistory(uint256 ) view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorSession) StateHistory(arg0 *big.Int) ([32]byte, error) {
	return _L2StateAggregator.Contract.StateHistory(&_L2StateAggregator.CallOpts, arg0)
}

// StateHistory is a free data retrieval call binding the contract method 0x90544984.
//
// Solidity: function stateHistory(uint256 ) view returns(bytes32)
func (_L2StateAggregator *L2StateAggregatorCallerSession) StateHistory(arg0 *big.Int) ([32]byte, error) {
	return _L2StateAggregator.Contract.StateHistory(&_L2StateAggregator.CallOpts, arg0)
}

// StateRootToBlock is a free data retrieval call binding the contract method 0xc8e2469e.
//
// Solidity: function stateRootToBlock(bytes32 ) view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorCaller) StateRootToBlock(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "stateRootToBlock", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StateRootToBlock is a free data retrieval call binding the contract method 0xc8e2469e.
//
// Solidity: function stateRootToBlock(bytes32 ) view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorSession) StateRootToBlock(arg0 [32]byte) (*big.Int, error) {
	return _L2StateAggregator.Contract.StateRootToBlock(&_L2StateAggregator.CallOpts, arg0)
}

// StateRootToBlock is a free data retrieval call binding the contract method 0xc8e2469e.
//
// Solidity: function stateRootToBlock(bytes32 ) view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorCallerSession) StateRootToBlock(arg0 [32]byte) (*big.Int, error) {
	return _L2StateAggregator.Contract.StateRootToBlock(&_L2StateAggregator.CallOpts, arg0)
}

// SystemState is a free data retrieval call binding the contract method 0x991292e3.
//
// Solidity: function systemState() view returns(uint256 totalCollateral, uint256 totalDebt, uint256 activePositions, uint256 totalOrders, uint256 blockNumber, uint256 timestamp)
func (_L2StateAggregator *L2StateAggregatorCaller) SystemState(opts *bind.CallOpts) (struct {
	TotalCollateral *big.Int
	TotalDebt       *big.Int
	ActivePositions *big.Int
	TotalOrders     *big.Int
	BlockNumber     *big.Int
	Timestamp       *big.Int
}, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "systemState")

	outstruct := new(struct {
		TotalCollateral *big.Int
		TotalDebt       *big.Int
		ActivePositions *big.Int
		TotalOrders     *big.Int
		BlockNumber     *big.Int
		Timestamp       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalCollateral = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalDebt = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ActivePositions = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TotalOrders = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.BlockNumber = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// SystemState is a free data retrieval call binding the contract method 0x991292e3.
//
// Solidity: function systemState() view returns(uint256 totalCollateral, uint256 totalDebt, uint256 activePositions, uint256 totalOrders, uint256 blockNumber, uint256 timestamp)
func (_L2StateAggregator *L2StateAggregatorSession) SystemState() (struct {
	TotalCollateral *big.Int
	TotalDebt       *big.Int
	ActivePositions *big.Int
	TotalOrders     *big.Int
	BlockNumber     *big.Int
	Timestamp       *big.Int
}, error) {
	return _L2StateAggregator.Contract.SystemState(&_L2StateAggregator.CallOpts)
}

// SystemState is a free data retrieval call binding the contract method 0x991292e3.
//
// Solidity: function systemState() view returns(uint256 totalCollateral, uint256 totalDebt, uint256 activePositions, uint256 totalOrders, uint256 blockNumber, uint256 timestamp)
func (_L2StateAggregator *L2StateAggregatorCallerSession) SystemState() (struct {
	TotalCollateral *big.Int
	TotalDebt       *big.Int
	ActivePositions *big.Int
	TotalOrders     *big.Int
	BlockNumber     *big.Int
	Timestamp       *big.Int
}, error) {
	return _L2StateAggregator.Contract.SystemState(&_L2StateAggregator.CallOpts)
}

// TimeUntilNextSubmission is a free data retrieval call binding the contract method 0x34624ca3.
//
// Solidity: function timeUntilNextSubmission() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorCaller) TimeUntilNextSubmission(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L2StateAggregator.contract.Call(opts, &out, "timeUntilNextSubmission")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TimeUntilNextSubmission is a free data retrieval call binding the contract method 0x34624ca3.
//
// Solidity: function timeUntilNextSubmission() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorSession) TimeUntilNextSubmission() (*big.Int, error) {
	return _L2StateAggregator.Contract.TimeUntilNextSubmission(&_L2StateAggregator.CallOpts)
}

// TimeUntilNextSubmission is a free data retrieval call binding the contract method 0x34624ca3.
//
// Solidity: function timeUntilNextSubmission() view returns(uint256)
func (_L2StateAggregator *L2StateAggregatorCallerSession) TimeUntilNextSubmission() (*big.Int, error) {
	return _L2StateAggregator.Contract.TimeUntilNextSubmission(&_L2StateAggregator.CallOpts)
}

// AuthorizeModule is a paid mutator transaction binding the contract method 0xa30a8e71.
//
// Solidity: function authorizeModule(address moduleAddress) returns()
func (_L2StateAggregator *L2StateAggregatorTransactor) AuthorizeModule(opts *bind.TransactOpts, moduleAddress common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.contract.Transact(opts, "authorizeModule", moduleAddress)
}

// AuthorizeModule is a paid mutator transaction binding the contract method 0xa30a8e71.
//
// Solidity: function authorizeModule(address moduleAddress) returns()
func (_L2StateAggregator *L2StateAggregatorSession) AuthorizeModule(moduleAddress common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.AuthorizeModule(&_L2StateAggregator.TransactOpts, moduleAddress)
}

// AuthorizeModule is a paid mutator transaction binding the contract method 0xa30a8e71.
//
// Solidity: function authorizeModule(address moduleAddress) returns()
func (_L2StateAggregator *L2StateAggregatorTransactorSession) AuthorizeModule(moduleAddress common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.AuthorizeModule(&_L2StateAggregator.TransactOpts, moduleAddress)
}

// DeactivateModule is a paid mutator transaction binding the contract method 0xa446f568.
//
// Solidity: function deactivateModule(bytes32 moduleId) returns()
func (_L2StateAggregator *L2StateAggregatorTransactor) DeactivateModule(opts *bind.TransactOpts, moduleId [32]byte) (*types.Transaction, error) {
	return _L2StateAggregator.contract.Transact(opts, "deactivateModule", moduleId)
}

// DeactivateModule is a paid mutator transaction binding the contract method 0xa446f568.
//
// Solidity: function deactivateModule(bytes32 moduleId) returns()
func (_L2StateAggregator *L2StateAggregatorSession) DeactivateModule(moduleId [32]byte) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.DeactivateModule(&_L2StateAggregator.TransactOpts, moduleId)
}

// DeactivateModule is a paid mutator transaction binding the contract method 0xa446f568.
//
// Solidity: function deactivateModule(bytes32 moduleId) returns()
func (_L2StateAggregator *L2StateAggregatorTransactorSession) DeactivateModule(moduleId [32]byte) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.DeactivateModule(&_L2StateAggregator.TransactOpts, moduleId)
}

// ForceSubmitToL1 is a paid mutator transaction binding the contract method 0x8ab94eaf.
//
// Solidity: function forceSubmitToL1() returns()
func (_L2StateAggregator *L2StateAggregatorTransactor) ForceSubmitToL1(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2StateAggregator.contract.Transact(opts, "forceSubmitToL1")
}

// ForceSubmitToL1 is a paid mutator transaction binding the contract method 0x8ab94eaf.
//
// Solidity: function forceSubmitToL1() returns()
func (_L2StateAggregator *L2StateAggregatorSession) ForceSubmitToL1() (*types.Transaction, error) {
	return _L2StateAggregator.Contract.ForceSubmitToL1(&_L2StateAggregator.TransactOpts)
}

// ForceSubmitToL1 is a paid mutator transaction binding the contract method 0x8ab94eaf.
//
// Solidity: function forceSubmitToL1() returns()
func (_L2StateAggregator *L2StateAggregatorTransactorSession) ForceSubmitToL1() (*types.Transaction, error) {
	return _L2StateAggregator.Contract.ForceSubmitToL1(&_L2StateAggregator.TransactOpts)
}

// RegisterModule is a paid mutator transaction binding the contract method 0xa78e922b.
//
// Solidity: function registerModule(bytes32 moduleId, address moduleAddress) returns()
func (_L2StateAggregator *L2StateAggregatorTransactor) RegisterModule(opts *bind.TransactOpts, moduleId [32]byte, moduleAddress common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.contract.Transact(opts, "registerModule", moduleId, moduleAddress)
}

// RegisterModule is a paid mutator transaction binding the contract method 0xa78e922b.
//
// Solidity: function registerModule(bytes32 moduleId, address moduleAddress) returns()
func (_L2StateAggregator *L2StateAggregatorSession) RegisterModule(moduleId [32]byte, moduleAddress common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.RegisterModule(&_L2StateAggregator.TransactOpts, moduleId, moduleAddress)
}

// RegisterModule is a paid mutator transaction binding the contract method 0xa78e922b.
//
// Solidity: function registerModule(bytes32 moduleId, address moduleAddress) returns()
func (_L2StateAggregator *L2StateAggregatorTransactorSession) RegisterModule(moduleId [32]byte, moduleAddress common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.RegisterModule(&_L2StateAggregator.TransactOpts, moduleId, moduleAddress)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L2StateAggregator *L2StateAggregatorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2StateAggregator.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L2StateAggregator *L2StateAggregatorSession) RenounceOwnership() (*types.Transaction, error) {
	return _L2StateAggregator.Contract.RenounceOwnership(&_L2StateAggregator.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L2StateAggregator *L2StateAggregatorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _L2StateAggregator.Contract.RenounceOwnership(&_L2StateAggregator.TransactOpts)
}

// RevokeModuleAuthorization is a paid mutator transaction binding the contract method 0x8e4df2f1.
//
// Solidity: function revokeModuleAuthorization(address moduleAddress) returns()
func (_L2StateAggregator *L2StateAggregatorTransactor) RevokeModuleAuthorization(opts *bind.TransactOpts, moduleAddress common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.contract.Transact(opts, "revokeModuleAuthorization", moduleAddress)
}

// RevokeModuleAuthorization is a paid mutator transaction binding the contract method 0x8e4df2f1.
//
// Solidity: function revokeModuleAuthorization(address moduleAddress) returns()
func (_L2StateAggregator *L2StateAggregatorSession) RevokeModuleAuthorization(moduleAddress common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.RevokeModuleAuthorization(&_L2StateAggregator.TransactOpts, moduleAddress)
}

// RevokeModuleAuthorization is a paid mutator transaction binding the contract method 0x8e4df2f1.
//
// Solidity: function revokeModuleAuthorization(address moduleAddress) returns()
func (_L2StateAggregator *L2StateAggregatorTransactorSession) RevokeModuleAuthorization(moduleAddress common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.RevokeModuleAuthorization(&_L2StateAggregator.TransactOpts, moduleAddress)
}

// SetL1StateRegistry is a paid mutator transaction binding the contract method 0x902e0dd5.
//
// Solidity: function setL1StateRegistry(address _l1StateRegistry) returns()
func (_L2StateAggregator *L2StateAggregatorTransactor) SetL1StateRegistry(opts *bind.TransactOpts, _l1StateRegistry common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.contract.Transact(opts, "setL1StateRegistry", _l1StateRegistry)
}

// SetL1StateRegistry is a paid mutator transaction binding the contract method 0x902e0dd5.
//
// Solidity: function setL1StateRegistry(address _l1StateRegistry) returns()
func (_L2StateAggregator *L2StateAggregatorSession) SetL1StateRegistry(_l1StateRegistry common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.SetL1StateRegistry(&_L2StateAggregator.TransactOpts, _l1StateRegistry)
}

// SetL1StateRegistry is a paid mutator transaction binding the contract method 0x902e0dd5.
//
// Solidity: function setL1StateRegistry(address _l1StateRegistry) returns()
func (_L2StateAggregator *L2StateAggregatorTransactorSession) SetL1StateRegistry(_l1StateRegistry common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.SetL1StateRegistry(&_L2StateAggregator.TransactOpts, _l1StateRegistry)
}

// SubmitToL1 is a paid mutator transaction binding the contract method 0x8c1799d1.
//
// Solidity: function submitToL1() returns()
func (_L2StateAggregator *L2StateAggregatorTransactor) SubmitToL1(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2StateAggregator.contract.Transact(opts, "submitToL1")
}

// SubmitToL1 is a paid mutator transaction binding the contract method 0x8c1799d1.
//
// Solidity: function submitToL1() returns()
func (_L2StateAggregator *L2StateAggregatorSession) SubmitToL1() (*types.Transaction, error) {
	return _L2StateAggregator.Contract.SubmitToL1(&_L2StateAggregator.TransactOpts)
}

// SubmitToL1 is a paid mutator transaction binding the contract method 0x8c1799d1.
//
// Solidity: function submitToL1() returns()
func (_L2StateAggregator *L2StateAggregatorTransactorSession) SubmitToL1() (*types.Transaction, error) {
	return _L2StateAggregator.Contract.SubmitToL1(&_L2StateAggregator.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2StateAggregator *L2StateAggregatorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2StateAggregator *L2StateAggregatorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.TransferOwnership(&_L2StateAggregator.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2StateAggregator *L2StateAggregatorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.TransferOwnership(&_L2StateAggregator.TransactOpts, newOwner)
}

// UpdateModuleState is a paid mutator transaction binding the contract method 0xd87c3376.
//
// Solidity: function updateModuleState(bytes32 moduleId, bytes32 stateHash) returns()
func (_L2StateAggregator *L2StateAggregatorTransactor) UpdateModuleState(opts *bind.TransactOpts, moduleId [32]byte, stateHash [32]byte) (*types.Transaction, error) {
	return _L2StateAggregator.contract.Transact(opts, "updateModuleState", moduleId, stateHash)
}

// UpdateModuleState is a paid mutator transaction binding the contract method 0xd87c3376.
//
// Solidity: function updateModuleState(bytes32 moduleId, bytes32 stateHash) returns()
func (_L2StateAggregator *L2StateAggregatorSession) UpdateModuleState(moduleId [32]byte, stateHash [32]byte) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.UpdateModuleState(&_L2StateAggregator.TransactOpts, moduleId, stateHash)
}

// UpdateModuleState is a paid mutator transaction binding the contract method 0xd87c3376.
//
// Solidity: function updateModuleState(bytes32 moduleId, bytes32 stateHash) returns()
func (_L2StateAggregator *L2StateAggregatorTransactorSession) UpdateModuleState(moduleId [32]byte, stateHash [32]byte) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.UpdateModuleState(&_L2StateAggregator.TransactOpts, moduleId, stateHash)
}

// UpdateSystemState is a paid mutator transaction binding the contract method 0x16a222e4.
//
// Solidity: function updateSystemState(uint256 totalCollateral, uint256 totalDebt, uint256 activePositions, uint256 totalOrders) returns()
func (_L2StateAggregator *L2StateAggregatorTransactor) UpdateSystemState(opts *bind.TransactOpts, totalCollateral *big.Int, totalDebt *big.Int, activePositions *big.Int, totalOrders *big.Int) (*types.Transaction, error) {
	return _L2StateAggregator.contract.Transact(opts, "updateSystemState", totalCollateral, totalDebt, activePositions, totalOrders)
}

// UpdateSystemState is a paid mutator transaction binding the contract method 0x16a222e4.
//
// Solidity: function updateSystemState(uint256 totalCollateral, uint256 totalDebt, uint256 activePositions, uint256 totalOrders) returns()
func (_L2StateAggregator *L2StateAggregatorSession) UpdateSystemState(totalCollateral *big.Int, totalDebt *big.Int, activePositions *big.Int, totalOrders *big.Int) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.UpdateSystemState(&_L2StateAggregator.TransactOpts, totalCollateral, totalDebt, activePositions, totalOrders)
}

// UpdateSystemState is a paid mutator transaction binding the contract method 0x16a222e4.
//
// Solidity: function updateSystemState(uint256 totalCollateral, uint256 totalDebt, uint256 activePositions, uint256 totalOrders) returns()
func (_L2StateAggregator *L2StateAggregatorTransactorSession) UpdateSystemState(totalCollateral *big.Int, totalDebt *big.Int, activePositions *big.Int, totalOrders *big.Int) (*types.Transaction, error) {
	return _L2StateAggregator.Contract.UpdateSystemState(&_L2StateAggregator.TransactOpts, totalCollateral, totalDebt, activePositions, totalOrders)
}

// L2StateAggregatorL1RegistryUpdatedIterator is returned from FilterL1RegistryUpdated and is used to iterate over the raw logs and unpacked data for L1RegistryUpdated events raised by the L2StateAggregator contract.
type L2StateAggregatorL1RegistryUpdatedIterator struct {
	Event *L2StateAggregatorL1RegistryUpdated // Event containing the contract specifics and raw log

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
func (it *L2StateAggregatorL1RegistryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2StateAggregatorL1RegistryUpdated)
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
		it.Event = new(L2StateAggregatorL1RegistryUpdated)
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
func (it *L2StateAggregatorL1RegistryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2StateAggregatorL1RegistryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2StateAggregatorL1RegistryUpdated represents a L1RegistryUpdated event raised by the L2StateAggregator contract.
type L2StateAggregatorL1RegistryUpdated struct {
	OldRegistry common.Address
	NewRegistry common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterL1RegistryUpdated is a free log retrieval operation binding the contract event 0x452f789ea74007bd52f942eca0f35f65ac2366985a059a268cab05aa3b4ecbf6.
//
// Solidity: event L1RegistryUpdated(address indexed oldRegistry, address indexed newRegistry)
func (_L2StateAggregator *L2StateAggregatorFilterer) FilterL1RegistryUpdated(opts *bind.FilterOpts, oldRegistry []common.Address, newRegistry []common.Address) (*L2StateAggregatorL1RegistryUpdatedIterator, error) {

	var oldRegistryRule []interface{}
	for _, oldRegistryItem := range oldRegistry {
		oldRegistryRule = append(oldRegistryRule, oldRegistryItem)
	}
	var newRegistryRule []interface{}
	for _, newRegistryItem := range newRegistry {
		newRegistryRule = append(newRegistryRule, newRegistryItem)
	}

	logs, sub, err := _L2StateAggregator.contract.FilterLogs(opts, "L1RegistryUpdated", oldRegistryRule, newRegistryRule)
	if err != nil {
		return nil, err
	}
	return &L2StateAggregatorL1RegistryUpdatedIterator{contract: _L2StateAggregator.contract, event: "L1RegistryUpdated", logs: logs, sub: sub}, nil
}

// WatchL1RegistryUpdated is a free log subscription operation binding the contract event 0x452f789ea74007bd52f942eca0f35f65ac2366985a059a268cab05aa3b4ecbf6.
//
// Solidity: event L1RegistryUpdated(address indexed oldRegistry, address indexed newRegistry)
func (_L2StateAggregator *L2StateAggregatorFilterer) WatchL1RegistryUpdated(opts *bind.WatchOpts, sink chan<- *L2StateAggregatorL1RegistryUpdated, oldRegistry []common.Address, newRegistry []common.Address) (event.Subscription, error) {

	var oldRegistryRule []interface{}
	for _, oldRegistryItem := range oldRegistry {
		oldRegistryRule = append(oldRegistryRule, oldRegistryItem)
	}
	var newRegistryRule []interface{}
	for _, newRegistryItem := range newRegistry {
		newRegistryRule = append(newRegistryRule, newRegistryItem)
	}

	logs, sub, err := _L2StateAggregator.contract.WatchLogs(opts, "L1RegistryUpdated", oldRegistryRule, newRegistryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2StateAggregatorL1RegistryUpdated)
				if err := _L2StateAggregator.contract.UnpackLog(event, "L1RegistryUpdated", log); err != nil {
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

// ParseL1RegistryUpdated is a log parse operation binding the contract event 0x452f789ea74007bd52f942eca0f35f65ac2366985a059a268cab05aa3b4ecbf6.
//
// Solidity: event L1RegistryUpdated(address indexed oldRegistry, address indexed newRegistry)
func (_L2StateAggregator *L2StateAggregatorFilterer) ParseL1RegistryUpdated(log types.Log) (*L2StateAggregatorL1RegistryUpdated, error) {
	event := new(L2StateAggregatorL1RegistryUpdated)
	if err := _L2StateAggregator.contract.UnpackLog(event, "L1RegistryUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2StateAggregatorModuleDeactivatedIterator is returned from FilterModuleDeactivated and is used to iterate over the raw logs and unpacked data for ModuleDeactivated events raised by the L2StateAggregator contract.
type L2StateAggregatorModuleDeactivatedIterator struct {
	Event *L2StateAggregatorModuleDeactivated // Event containing the contract specifics and raw log

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
func (it *L2StateAggregatorModuleDeactivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2StateAggregatorModuleDeactivated)
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
		it.Event = new(L2StateAggregatorModuleDeactivated)
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
func (it *L2StateAggregatorModuleDeactivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2StateAggregatorModuleDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2StateAggregatorModuleDeactivated represents a ModuleDeactivated event raised by the L2StateAggregator contract.
type L2StateAggregatorModuleDeactivated struct {
	ModuleId [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterModuleDeactivated is a free log retrieval operation binding the contract event 0x1b3395b7c125da23d86f4eeff6f9e82c3e1e159933a183b553fb0ceca0bd0be5.
//
// Solidity: event ModuleDeactivated(bytes32 indexed moduleId)
func (_L2StateAggregator *L2StateAggregatorFilterer) FilterModuleDeactivated(opts *bind.FilterOpts, moduleId [][32]byte) (*L2StateAggregatorModuleDeactivatedIterator, error) {

	var moduleIdRule []interface{}
	for _, moduleIdItem := range moduleId {
		moduleIdRule = append(moduleIdRule, moduleIdItem)
	}

	logs, sub, err := _L2StateAggregator.contract.FilterLogs(opts, "ModuleDeactivated", moduleIdRule)
	if err != nil {
		return nil, err
	}
	return &L2StateAggregatorModuleDeactivatedIterator{contract: _L2StateAggregator.contract, event: "ModuleDeactivated", logs: logs, sub: sub}, nil
}

// WatchModuleDeactivated is a free log subscription operation binding the contract event 0x1b3395b7c125da23d86f4eeff6f9e82c3e1e159933a183b553fb0ceca0bd0be5.
//
// Solidity: event ModuleDeactivated(bytes32 indexed moduleId)
func (_L2StateAggregator *L2StateAggregatorFilterer) WatchModuleDeactivated(opts *bind.WatchOpts, sink chan<- *L2StateAggregatorModuleDeactivated, moduleId [][32]byte) (event.Subscription, error) {

	var moduleIdRule []interface{}
	for _, moduleIdItem := range moduleId {
		moduleIdRule = append(moduleIdRule, moduleIdItem)
	}

	logs, sub, err := _L2StateAggregator.contract.WatchLogs(opts, "ModuleDeactivated", moduleIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2StateAggregatorModuleDeactivated)
				if err := _L2StateAggregator.contract.UnpackLog(event, "ModuleDeactivated", log); err != nil {
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

// ParseModuleDeactivated is a log parse operation binding the contract event 0x1b3395b7c125da23d86f4eeff6f9e82c3e1e159933a183b553fb0ceca0bd0be5.
//
// Solidity: event ModuleDeactivated(bytes32 indexed moduleId)
func (_L2StateAggregator *L2StateAggregatorFilterer) ParseModuleDeactivated(log types.Log) (*L2StateAggregatorModuleDeactivated, error) {
	event := new(L2StateAggregatorModuleDeactivated)
	if err := _L2StateAggregator.contract.UnpackLog(event, "ModuleDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2StateAggregatorModuleRegisteredIterator is returned from FilterModuleRegistered and is used to iterate over the raw logs and unpacked data for ModuleRegistered events raised by the L2StateAggregator contract.
type L2StateAggregatorModuleRegisteredIterator struct {
	Event *L2StateAggregatorModuleRegistered // Event containing the contract specifics and raw log

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
func (it *L2StateAggregatorModuleRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2StateAggregatorModuleRegistered)
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
		it.Event = new(L2StateAggregatorModuleRegistered)
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
func (it *L2StateAggregatorModuleRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2StateAggregatorModuleRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2StateAggregatorModuleRegistered represents a ModuleRegistered event raised by the L2StateAggregator contract.
type L2StateAggregatorModuleRegistered struct {
	ModuleId      [32]byte
	ModuleAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterModuleRegistered is a free log retrieval operation binding the contract event 0xd63be02155b46636309fb0a4a79647c60971aecaad53cbc83aad90cd75fd9d54.
//
// Solidity: event ModuleRegistered(bytes32 indexed moduleId, address moduleAddress)
func (_L2StateAggregator *L2StateAggregatorFilterer) FilterModuleRegistered(opts *bind.FilterOpts, moduleId [][32]byte) (*L2StateAggregatorModuleRegisteredIterator, error) {

	var moduleIdRule []interface{}
	for _, moduleIdItem := range moduleId {
		moduleIdRule = append(moduleIdRule, moduleIdItem)
	}

	logs, sub, err := _L2StateAggregator.contract.FilterLogs(opts, "ModuleRegistered", moduleIdRule)
	if err != nil {
		return nil, err
	}
	return &L2StateAggregatorModuleRegisteredIterator{contract: _L2StateAggregator.contract, event: "ModuleRegistered", logs: logs, sub: sub}, nil
}

// WatchModuleRegistered is a free log subscription operation binding the contract event 0xd63be02155b46636309fb0a4a79647c60971aecaad53cbc83aad90cd75fd9d54.
//
// Solidity: event ModuleRegistered(bytes32 indexed moduleId, address moduleAddress)
func (_L2StateAggregator *L2StateAggregatorFilterer) WatchModuleRegistered(opts *bind.WatchOpts, sink chan<- *L2StateAggregatorModuleRegistered, moduleId [][32]byte) (event.Subscription, error) {

	var moduleIdRule []interface{}
	for _, moduleIdItem := range moduleId {
		moduleIdRule = append(moduleIdRule, moduleIdItem)
	}

	logs, sub, err := _L2StateAggregator.contract.WatchLogs(opts, "ModuleRegistered", moduleIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2StateAggregatorModuleRegistered)
				if err := _L2StateAggregator.contract.UnpackLog(event, "ModuleRegistered", log); err != nil {
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

// ParseModuleRegistered is a log parse operation binding the contract event 0xd63be02155b46636309fb0a4a79647c60971aecaad53cbc83aad90cd75fd9d54.
//
// Solidity: event ModuleRegistered(bytes32 indexed moduleId, address moduleAddress)
func (_L2StateAggregator *L2StateAggregatorFilterer) ParseModuleRegistered(log types.Log) (*L2StateAggregatorModuleRegistered, error) {
	event := new(L2StateAggregatorModuleRegistered)
	if err := _L2StateAggregator.contract.UnpackLog(event, "ModuleRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2StateAggregatorModuleStateUpdatedIterator is returned from FilterModuleStateUpdated and is used to iterate over the raw logs and unpacked data for ModuleStateUpdated events raised by the L2StateAggregator contract.
type L2StateAggregatorModuleStateUpdatedIterator struct {
	Event *L2StateAggregatorModuleStateUpdated // Event containing the contract specifics and raw log

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
func (it *L2StateAggregatorModuleStateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2StateAggregatorModuleStateUpdated)
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
		it.Event = new(L2StateAggregatorModuleStateUpdated)
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
func (it *L2StateAggregatorModuleStateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2StateAggregatorModuleStateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2StateAggregatorModuleStateUpdated represents a ModuleStateUpdated event raised by the L2StateAggregator contract.
type L2StateAggregatorModuleStateUpdated struct {
	ModuleId  [32]byte
	StateHash [32]byte
	UpdatedBy common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterModuleStateUpdated is a free log retrieval operation binding the contract event 0xd99ea162e2f693d07eb615544a07553554baf290ce277fe70379dbbd1f1b2761.
//
// Solidity: event ModuleStateUpdated(bytes32 indexed moduleId, bytes32 stateHash, address updatedBy)
func (_L2StateAggregator *L2StateAggregatorFilterer) FilterModuleStateUpdated(opts *bind.FilterOpts, moduleId [][32]byte) (*L2StateAggregatorModuleStateUpdatedIterator, error) {

	var moduleIdRule []interface{}
	for _, moduleIdItem := range moduleId {
		moduleIdRule = append(moduleIdRule, moduleIdItem)
	}

	logs, sub, err := _L2StateAggregator.contract.FilterLogs(opts, "ModuleStateUpdated", moduleIdRule)
	if err != nil {
		return nil, err
	}
	return &L2StateAggregatorModuleStateUpdatedIterator{contract: _L2StateAggregator.contract, event: "ModuleStateUpdated", logs: logs, sub: sub}, nil
}

// WatchModuleStateUpdated is a free log subscription operation binding the contract event 0xd99ea162e2f693d07eb615544a07553554baf290ce277fe70379dbbd1f1b2761.
//
// Solidity: event ModuleStateUpdated(bytes32 indexed moduleId, bytes32 stateHash, address updatedBy)
func (_L2StateAggregator *L2StateAggregatorFilterer) WatchModuleStateUpdated(opts *bind.WatchOpts, sink chan<- *L2StateAggregatorModuleStateUpdated, moduleId [][32]byte) (event.Subscription, error) {

	var moduleIdRule []interface{}
	for _, moduleIdItem := range moduleId {
		moduleIdRule = append(moduleIdRule, moduleIdItem)
	}

	logs, sub, err := _L2StateAggregator.contract.WatchLogs(opts, "ModuleStateUpdated", moduleIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2StateAggregatorModuleStateUpdated)
				if err := _L2StateAggregator.contract.UnpackLog(event, "ModuleStateUpdated", log); err != nil {
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

// ParseModuleStateUpdated is a log parse operation binding the contract event 0xd99ea162e2f693d07eb615544a07553554baf290ce277fe70379dbbd1f1b2761.
//
// Solidity: event ModuleStateUpdated(bytes32 indexed moduleId, bytes32 stateHash, address updatedBy)
func (_L2StateAggregator *L2StateAggregatorFilterer) ParseModuleStateUpdated(log types.Log) (*L2StateAggregatorModuleStateUpdated, error) {
	event := new(L2StateAggregatorModuleStateUpdated)
	if err := _L2StateAggregator.contract.UnpackLog(event, "ModuleStateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2StateAggregatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the L2StateAggregator contract.
type L2StateAggregatorOwnershipTransferredIterator struct {
	Event *L2StateAggregatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *L2StateAggregatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2StateAggregatorOwnershipTransferred)
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
		it.Event = new(L2StateAggregatorOwnershipTransferred)
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
func (it *L2StateAggregatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2StateAggregatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2StateAggregatorOwnershipTransferred represents a OwnershipTransferred event raised by the L2StateAggregator contract.
type L2StateAggregatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L2StateAggregator *L2StateAggregatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*L2StateAggregatorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L2StateAggregator.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &L2StateAggregatorOwnershipTransferredIterator{contract: _L2StateAggregator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L2StateAggregator *L2StateAggregatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *L2StateAggregatorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L2StateAggregator.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2StateAggregatorOwnershipTransferred)
				if err := _L2StateAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_L2StateAggregator *L2StateAggregatorFilterer) ParseOwnershipTransferred(log types.Log) (*L2StateAggregatorOwnershipTransferred, error) {
	event := new(L2StateAggregatorOwnershipTransferred)
	if err := _L2StateAggregator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2StateAggregatorStateRootComputedIterator is returned from FilterStateRootComputed and is used to iterate over the raw logs and unpacked data for StateRootComputed events raised by the L2StateAggregator contract.
type L2StateAggregatorStateRootComputedIterator struct {
	Event *L2StateAggregatorStateRootComputed // Event containing the contract specifics and raw log

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
func (it *L2StateAggregatorStateRootComputedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2StateAggregatorStateRootComputed)
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
		it.Event = new(L2StateAggregatorStateRootComputed)
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
func (it *L2StateAggregatorStateRootComputedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2StateAggregatorStateRootComputedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2StateAggregatorStateRootComputed represents a StateRootComputed event raised by the L2StateAggregator contract.
type L2StateAggregatorStateRootComputed struct {
	StateRoot   [32]byte
	BlockNumber *big.Int
	Timestamp   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterStateRootComputed is a free log retrieval operation binding the contract event 0xa938aa98582150744c56367318d3a5f5cd92bde822f25ff5a31eeb88e36b4121.
//
// Solidity: event StateRootComputed(bytes32 indexed stateRoot, uint256 blockNumber, uint256 timestamp)
func (_L2StateAggregator *L2StateAggregatorFilterer) FilterStateRootComputed(opts *bind.FilterOpts, stateRoot [][32]byte) (*L2StateAggregatorStateRootComputedIterator, error) {

	var stateRootRule []interface{}
	for _, stateRootItem := range stateRoot {
		stateRootRule = append(stateRootRule, stateRootItem)
	}

	logs, sub, err := _L2StateAggregator.contract.FilterLogs(opts, "StateRootComputed", stateRootRule)
	if err != nil {
		return nil, err
	}
	return &L2StateAggregatorStateRootComputedIterator{contract: _L2StateAggregator.contract, event: "StateRootComputed", logs: logs, sub: sub}, nil
}

// WatchStateRootComputed is a free log subscription operation binding the contract event 0xa938aa98582150744c56367318d3a5f5cd92bde822f25ff5a31eeb88e36b4121.
//
// Solidity: event StateRootComputed(bytes32 indexed stateRoot, uint256 blockNumber, uint256 timestamp)
func (_L2StateAggregator *L2StateAggregatorFilterer) WatchStateRootComputed(opts *bind.WatchOpts, sink chan<- *L2StateAggregatorStateRootComputed, stateRoot [][32]byte) (event.Subscription, error) {

	var stateRootRule []interface{}
	for _, stateRootItem := range stateRoot {
		stateRootRule = append(stateRootRule, stateRootItem)
	}

	logs, sub, err := _L2StateAggregator.contract.WatchLogs(opts, "StateRootComputed", stateRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2StateAggregatorStateRootComputed)
				if err := _L2StateAggregator.contract.UnpackLog(event, "StateRootComputed", log); err != nil {
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

// ParseStateRootComputed is a log parse operation binding the contract event 0xa938aa98582150744c56367318d3a5f5cd92bde822f25ff5a31eeb88e36b4121.
//
// Solidity: event StateRootComputed(bytes32 indexed stateRoot, uint256 blockNumber, uint256 timestamp)
func (_L2StateAggregator *L2StateAggregatorFilterer) ParseStateRootComputed(log types.Log) (*L2StateAggregatorStateRootComputed, error) {
	event := new(L2StateAggregatorStateRootComputed)
	if err := _L2StateAggregator.contract.UnpackLog(event, "StateRootComputed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2StateAggregatorStateSubmittedToL1Iterator is returned from FilterStateSubmittedToL1 and is used to iterate over the raw logs and unpacked data for StateSubmittedToL1 events raised by the L2StateAggregator contract.
type L2StateAggregatorStateSubmittedToL1Iterator struct {
	Event *L2StateAggregatorStateSubmittedToL1 // Event containing the contract specifics and raw log

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
func (it *L2StateAggregatorStateSubmittedToL1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2StateAggregatorStateSubmittedToL1)
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
		it.Event = new(L2StateAggregatorStateSubmittedToL1)
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
func (it *L2StateAggregatorStateSubmittedToL1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2StateAggregatorStateSubmittedToL1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2StateAggregatorStateSubmittedToL1 represents a StateSubmittedToL1 event raised by the L2StateAggregator contract.
type L2StateAggregatorStateSubmittedToL1 struct {
	StateRoot       [32]byte
	L2Block         *big.Int
	L1TxId          *big.Int
	TotalCollateral *big.Int
	TotalDebt       *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterStateSubmittedToL1 is a free log retrieval operation binding the contract event 0xc8281a000eab8bc6f4321b42029109cbb1e783ff071b93128f7a2407727bc466.
//
// Solidity: event StateSubmittedToL1(bytes32 indexed stateRoot, uint256 indexed l2Block, uint256 indexed l1TxId, uint256 totalCollateral, uint256 totalDebt)
func (_L2StateAggregator *L2StateAggregatorFilterer) FilterStateSubmittedToL1(opts *bind.FilterOpts, stateRoot [][32]byte, l2Block []*big.Int, l1TxId []*big.Int) (*L2StateAggregatorStateSubmittedToL1Iterator, error) {

	var stateRootRule []interface{}
	for _, stateRootItem := range stateRoot {
		stateRootRule = append(stateRootRule, stateRootItem)
	}
	var l2BlockRule []interface{}
	for _, l2BlockItem := range l2Block {
		l2BlockRule = append(l2BlockRule, l2BlockItem)
	}
	var l1TxIdRule []interface{}
	for _, l1TxIdItem := range l1TxId {
		l1TxIdRule = append(l1TxIdRule, l1TxIdItem)
	}

	logs, sub, err := _L2StateAggregator.contract.FilterLogs(opts, "StateSubmittedToL1", stateRootRule, l2BlockRule, l1TxIdRule)
	if err != nil {
		return nil, err
	}
	return &L2StateAggregatorStateSubmittedToL1Iterator{contract: _L2StateAggregator.contract, event: "StateSubmittedToL1", logs: logs, sub: sub}, nil
}

// WatchStateSubmittedToL1 is a free log subscription operation binding the contract event 0xc8281a000eab8bc6f4321b42029109cbb1e783ff071b93128f7a2407727bc466.
//
// Solidity: event StateSubmittedToL1(bytes32 indexed stateRoot, uint256 indexed l2Block, uint256 indexed l1TxId, uint256 totalCollateral, uint256 totalDebt)
func (_L2StateAggregator *L2StateAggregatorFilterer) WatchStateSubmittedToL1(opts *bind.WatchOpts, sink chan<- *L2StateAggregatorStateSubmittedToL1, stateRoot [][32]byte, l2Block []*big.Int, l1TxId []*big.Int) (event.Subscription, error) {

	var stateRootRule []interface{}
	for _, stateRootItem := range stateRoot {
		stateRootRule = append(stateRootRule, stateRootItem)
	}
	var l2BlockRule []interface{}
	for _, l2BlockItem := range l2Block {
		l2BlockRule = append(l2BlockRule, l2BlockItem)
	}
	var l1TxIdRule []interface{}
	for _, l1TxIdItem := range l1TxId {
		l1TxIdRule = append(l1TxIdRule, l1TxIdItem)
	}

	logs, sub, err := _L2StateAggregator.contract.WatchLogs(opts, "StateSubmittedToL1", stateRootRule, l2BlockRule, l1TxIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2StateAggregatorStateSubmittedToL1)
				if err := _L2StateAggregator.contract.UnpackLog(event, "StateSubmittedToL1", log); err != nil {
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

// ParseStateSubmittedToL1 is a log parse operation binding the contract event 0xc8281a000eab8bc6f4321b42029109cbb1e783ff071b93128f7a2407727bc466.
//
// Solidity: event StateSubmittedToL1(bytes32 indexed stateRoot, uint256 indexed l2Block, uint256 indexed l1TxId, uint256 totalCollateral, uint256 totalDebt)
func (_L2StateAggregator *L2StateAggregatorFilterer) ParseStateSubmittedToL1(log types.Log) (*L2StateAggregatorStateSubmittedToL1, error) {
	event := new(L2StateAggregatorStateSubmittedToL1)
	if err := _L2StateAggregator.contract.UnpackLog(event, "StateSubmittedToL1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
