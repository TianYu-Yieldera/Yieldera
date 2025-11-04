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

// CompoundV3AdapterMetaData contains all meta data concerning the CompoundV3Adapter contract.
var CompoundV3AdapterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_comet\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_cometRewards\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BaseSupplied\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BaseWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Borrowed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"CollateralSupplied\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"CollateralWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Repaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RewardsClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MANAGER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activePositions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseAsset\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"borrow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"shouldAccrue\",\"type\":\"bool\"}],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"comet\",\"outputs\":[{\"internalType\":\"contractIComet\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cometRewards\",\"outputs\":[{\"internalType\":\"contractICometRewards\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emergencyWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getPosition\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"baseSupplied\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseBorrowed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collateralValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRates\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"supplyRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowRate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"isPositionHealthy\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"positions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"baseSupplied\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseBorrowed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdate\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"repay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stateAggregator\",\"type\":\"address\"}],\"name\":\"setStateAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateAggregator\",\"outputs\":[{\"internalType\":\"contractL2StateAggregator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"supplyBase\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"supplyCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalBorrowed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupplied\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawBase\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// CompoundV3AdapterABI is the input ABI used to generate the binding from.
// Deprecated: Use CompoundV3AdapterMetaData.ABI instead.
var CompoundV3AdapterABI = CompoundV3AdapterMetaData.ABI

// CompoundV3Adapter is an auto generated Go binding around an Ethereum contract.
type CompoundV3Adapter struct {
	CompoundV3AdapterCaller     // Read-only binding to the contract
	CompoundV3AdapterTransactor // Write-only binding to the contract
	CompoundV3AdapterFilterer   // Log filterer for contract events
}

// CompoundV3AdapterCaller is an auto generated read-only Go binding around an Ethereum contract.
type CompoundV3AdapterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompoundV3AdapterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CompoundV3AdapterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompoundV3AdapterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CompoundV3AdapterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompoundV3AdapterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CompoundV3AdapterSession struct {
	Contract     *CompoundV3Adapter // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// CompoundV3AdapterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CompoundV3AdapterCallerSession struct {
	Contract *CompoundV3AdapterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// CompoundV3AdapterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CompoundV3AdapterTransactorSession struct {
	Contract     *CompoundV3AdapterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// CompoundV3AdapterRaw is an auto generated low-level Go binding around an Ethereum contract.
type CompoundV3AdapterRaw struct {
	Contract *CompoundV3Adapter // Generic contract binding to access the raw methods on
}

// CompoundV3AdapterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CompoundV3AdapterCallerRaw struct {
	Contract *CompoundV3AdapterCaller // Generic read-only contract binding to access the raw methods on
}

// CompoundV3AdapterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CompoundV3AdapterTransactorRaw struct {
	Contract *CompoundV3AdapterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCompoundV3Adapter creates a new instance of CompoundV3Adapter, bound to a specific deployed contract.
func NewCompoundV3Adapter(address common.Address, backend bind.ContractBackend) (*CompoundV3Adapter, error) {
	contract, err := bindCompoundV3Adapter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CompoundV3Adapter{CompoundV3AdapterCaller: CompoundV3AdapterCaller{contract: contract}, CompoundV3AdapterTransactor: CompoundV3AdapterTransactor{contract: contract}, CompoundV3AdapterFilterer: CompoundV3AdapterFilterer{contract: contract}}, nil
}

// NewCompoundV3AdapterCaller creates a new read-only instance of CompoundV3Adapter, bound to a specific deployed contract.
func NewCompoundV3AdapterCaller(address common.Address, caller bind.ContractCaller) (*CompoundV3AdapterCaller, error) {
	contract, err := bindCompoundV3Adapter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterCaller{contract: contract}, nil
}

// NewCompoundV3AdapterTransactor creates a new write-only instance of CompoundV3Adapter, bound to a specific deployed contract.
func NewCompoundV3AdapterTransactor(address common.Address, transactor bind.ContractTransactor) (*CompoundV3AdapterTransactor, error) {
	contract, err := bindCompoundV3Adapter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterTransactor{contract: contract}, nil
}

// NewCompoundV3AdapterFilterer creates a new log filterer instance of CompoundV3Adapter, bound to a specific deployed contract.
func NewCompoundV3AdapterFilterer(address common.Address, filterer bind.ContractFilterer) (*CompoundV3AdapterFilterer, error) {
	contract, err := bindCompoundV3Adapter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterFilterer{contract: contract}, nil
}

// bindCompoundV3Adapter binds a generic wrapper to an already deployed contract.
func bindCompoundV3Adapter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CompoundV3AdapterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CompoundV3Adapter *CompoundV3AdapterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CompoundV3Adapter.Contract.CompoundV3AdapterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CompoundV3Adapter *CompoundV3AdapterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.CompoundV3AdapterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CompoundV3Adapter *CompoundV3AdapterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.CompoundV3AdapterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CompoundV3Adapter *CompoundV3AdapterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CompoundV3Adapter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CompoundV3Adapter *CompoundV3AdapterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CompoundV3Adapter *CompoundV3AdapterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_CompoundV3Adapter *CompoundV3AdapterSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _CompoundV3Adapter.Contract.DEFAULTADMINROLE(&_CompoundV3Adapter.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _CompoundV3Adapter.Contract.DEFAULTADMINROLE(&_CompoundV3Adapter.CallOpts)
}

// MANAGERROLE is a free data retrieval call binding the contract method 0xec87621c.
//
// Solidity: function MANAGER_ROLE() view returns(bytes32)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) MANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MANAGERROLE is a free data retrieval call binding the contract method 0xec87621c.
//
// Solidity: function MANAGER_ROLE() view returns(bytes32)
func (_CompoundV3Adapter *CompoundV3AdapterSession) MANAGERROLE() ([32]byte, error) {
	return _CompoundV3Adapter.Contract.MANAGERROLE(&_CompoundV3Adapter.CallOpts)
}

// MANAGERROLE is a free data retrieval call binding the contract method 0xec87621c.
//
// Solidity: function MANAGER_ROLE() view returns(bytes32)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) MANAGERROLE() ([32]byte, error) {
	return _CompoundV3Adapter.Contract.MANAGERROLE(&_CompoundV3Adapter.CallOpts)
}

// ActivePositions is a free data retrieval call binding the contract method 0x297add9d.
//
// Solidity: function activePositions() view returns(uint256)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) ActivePositions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "activePositions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActivePositions is a free data retrieval call binding the contract method 0x297add9d.
//
// Solidity: function activePositions() view returns(uint256)
func (_CompoundV3Adapter *CompoundV3AdapterSession) ActivePositions() (*big.Int, error) {
	return _CompoundV3Adapter.Contract.ActivePositions(&_CompoundV3Adapter.CallOpts)
}

// ActivePositions is a free data retrieval call binding the contract method 0x297add9d.
//
// Solidity: function activePositions() view returns(uint256)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) ActivePositions() (*big.Int, error) {
	return _CompoundV3Adapter.Contract.ActivePositions(&_CompoundV3Adapter.CallOpts)
}

// BaseAsset is a free data retrieval call binding the contract method 0xcdf456e1.
//
// Solidity: function baseAsset() view returns(address)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) BaseAsset(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "baseAsset")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BaseAsset is a free data retrieval call binding the contract method 0xcdf456e1.
//
// Solidity: function baseAsset() view returns(address)
func (_CompoundV3Adapter *CompoundV3AdapterSession) BaseAsset() (common.Address, error) {
	return _CompoundV3Adapter.Contract.BaseAsset(&_CompoundV3Adapter.CallOpts)
}

// BaseAsset is a free data retrieval call binding the contract method 0xcdf456e1.
//
// Solidity: function baseAsset() view returns(address)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) BaseAsset() (common.Address, error) {
	return _CompoundV3Adapter.Contract.BaseAsset(&_CompoundV3Adapter.CallOpts)
}

// Comet is a free data retrieval call binding the contract method 0xba3e9c12.
//
// Solidity: function comet() view returns(address)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) Comet(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "comet")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Comet is a free data retrieval call binding the contract method 0xba3e9c12.
//
// Solidity: function comet() view returns(address)
func (_CompoundV3Adapter *CompoundV3AdapterSession) Comet() (common.Address, error) {
	return _CompoundV3Adapter.Contract.Comet(&_CompoundV3Adapter.CallOpts)
}

// Comet is a free data retrieval call binding the contract method 0xba3e9c12.
//
// Solidity: function comet() view returns(address)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) Comet() (common.Address, error) {
	return _CompoundV3Adapter.Contract.Comet(&_CompoundV3Adapter.CallOpts)
}

// CometRewards is a free data retrieval call binding the contract method 0x32315972.
//
// Solidity: function cometRewards() view returns(address)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) CometRewards(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "cometRewards")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CometRewards is a free data retrieval call binding the contract method 0x32315972.
//
// Solidity: function cometRewards() view returns(address)
func (_CompoundV3Adapter *CompoundV3AdapterSession) CometRewards() (common.Address, error) {
	return _CompoundV3Adapter.Contract.CometRewards(&_CompoundV3Adapter.CallOpts)
}

// CometRewards is a free data retrieval call binding the contract method 0x32315972.
//
// Solidity: function cometRewards() view returns(address)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) CometRewards() (common.Address, error) {
	return _CompoundV3Adapter.Contract.CometRewards(&_CompoundV3Adapter.CallOpts)
}

// GetPosition is a free data retrieval call binding the contract method 0x16c19739.
//
// Solidity: function getPosition(address user) view returns(uint256 baseSupplied, uint256 baseBorrowed, uint256 collateralValue, uint256 lastUpdate)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) GetPosition(opts *bind.CallOpts, user common.Address) (struct {
	BaseSupplied    *big.Int
	BaseBorrowed    *big.Int
	CollateralValue *big.Int
	LastUpdate      *big.Int
}, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "getPosition", user)

	outstruct := new(struct {
		BaseSupplied    *big.Int
		BaseBorrowed    *big.Int
		CollateralValue *big.Int
		LastUpdate      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BaseSupplied = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BaseBorrowed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.CollateralValue = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.LastUpdate = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetPosition is a free data retrieval call binding the contract method 0x16c19739.
//
// Solidity: function getPosition(address user) view returns(uint256 baseSupplied, uint256 baseBorrowed, uint256 collateralValue, uint256 lastUpdate)
func (_CompoundV3Adapter *CompoundV3AdapterSession) GetPosition(user common.Address) (struct {
	BaseSupplied    *big.Int
	BaseBorrowed    *big.Int
	CollateralValue *big.Int
	LastUpdate      *big.Int
}, error) {
	return _CompoundV3Adapter.Contract.GetPosition(&_CompoundV3Adapter.CallOpts, user)
}

// GetPosition is a free data retrieval call binding the contract method 0x16c19739.
//
// Solidity: function getPosition(address user) view returns(uint256 baseSupplied, uint256 baseBorrowed, uint256 collateralValue, uint256 lastUpdate)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) GetPosition(user common.Address) (struct {
	BaseSupplied    *big.Int
	BaseBorrowed    *big.Int
	CollateralValue *big.Int
	LastUpdate      *big.Int
}, error) {
	return _CompoundV3Adapter.Contract.GetPosition(&_CompoundV3Adapter.CallOpts, user)
}

// GetRates is a free data retrieval call binding the contract method 0x9accab55.
//
// Solidity: function getRates() view returns(uint256 supplyRate, uint256 borrowRate)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) GetRates(opts *bind.CallOpts) (struct {
	SupplyRate *big.Int
	BorrowRate *big.Int
}, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "getRates")

	outstruct := new(struct {
		SupplyRate *big.Int
		BorrowRate *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SupplyRate = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BorrowRate = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRates is a free data retrieval call binding the contract method 0x9accab55.
//
// Solidity: function getRates() view returns(uint256 supplyRate, uint256 borrowRate)
func (_CompoundV3Adapter *CompoundV3AdapterSession) GetRates() (struct {
	SupplyRate *big.Int
	BorrowRate *big.Int
}, error) {
	return _CompoundV3Adapter.Contract.GetRates(&_CompoundV3Adapter.CallOpts)
}

// GetRates is a free data retrieval call binding the contract method 0x9accab55.
//
// Solidity: function getRates() view returns(uint256 supplyRate, uint256 borrowRate)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) GetRates() (struct {
	SupplyRate *big.Int
	BorrowRate *big.Int
}, error) {
	return _CompoundV3Adapter.Contract.GetRates(&_CompoundV3Adapter.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_CompoundV3Adapter *CompoundV3AdapterSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _CompoundV3Adapter.Contract.GetRoleAdmin(&_CompoundV3Adapter.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _CompoundV3Adapter.Contract.GetRoleAdmin(&_CompoundV3Adapter.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_CompoundV3Adapter *CompoundV3AdapterSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _CompoundV3Adapter.Contract.HasRole(&_CompoundV3Adapter.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _CompoundV3Adapter.Contract.HasRole(&_CompoundV3Adapter.CallOpts, role, account)
}

// IsPositionHealthy is a free data retrieval call binding the contract method 0xbafca75e.
//
// Solidity: function isPositionHealthy(address user) view returns(bool)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) IsPositionHealthy(opts *bind.CallOpts, user common.Address) (bool, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "isPositionHealthy", user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPositionHealthy is a free data retrieval call binding the contract method 0xbafca75e.
//
// Solidity: function isPositionHealthy(address user) view returns(bool)
func (_CompoundV3Adapter *CompoundV3AdapterSession) IsPositionHealthy(user common.Address) (bool, error) {
	return _CompoundV3Adapter.Contract.IsPositionHealthy(&_CompoundV3Adapter.CallOpts, user)
}

// IsPositionHealthy is a free data retrieval call binding the contract method 0xbafca75e.
//
// Solidity: function isPositionHealthy(address user) view returns(bool)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) IsPositionHealthy(user common.Address) (bool, error) {
	return _CompoundV3Adapter.Contract.IsPositionHealthy(&_CompoundV3Adapter.CallOpts, user)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CompoundV3Adapter *CompoundV3AdapterSession) Paused() (bool, error) {
	return _CompoundV3Adapter.Contract.Paused(&_CompoundV3Adapter.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) Paused() (bool, error) {
	return _CompoundV3Adapter.Contract.Paused(&_CompoundV3Adapter.CallOpts)
}

// Positions is a free data retrieval call binding the contract method 0x55f57510.
//
// Solidity: function positions(address ) view returns(uint256 baseSupplied, uint256 baseBorrowed, uint256 lastUpdate, bool active)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) Positions(opts *bind.CallOpts, arg0 common.Address) (struct {
	BaseSupplied *big.Int
	BaseBorrowed *big.Int
	LastUpdate   *big.Int
	Active       bool
}, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "positions", arg0)

	outstruct := new(struct {
		BaseSupplied *big.Int
		BaseBorrowed *big.Int
		LastUpdate   *big.Int
		Active       bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BaseSupplied = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BaseBorrowed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.LastUpdate = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Active = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// Positions is a free data retrieval call binding the contract method 0x55f57510.
//
// Solidity: function positions(address ) view returns(uint256 baseSupplied, uint256 baseBorrowed, uint256 lastUpdate, bool active)
func (_CompoundV3Adapter *CompoundV3AdapterSession) Positions(arg0 common.Address) (struct {
	BaseSupplied *big.Int
	BaseBorrowed *big.Int
	LastUpdate   *big.Int
	Active       bool
}, error) {
	return _CompoundV3Adapter.Contract.Positions(&_CompoundV3Adapter.CallOpts, arg0)
}

// Positions is a free data retrieval call binding the contract method 0x55f57510.
//
// Solidity: function positions(address ) view returns(uint256 baseSupplied, uint256 baseBorrowed, uint256 lastUpdate, bool active)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) Positions(arg0 common.Address) (struct {
	BaseSupplied *big.Int
	BaseBorrowed *big.Int
	LastUpdate   *big.Int
	Active       bool
}, error) {
	return _CompoundV3Adapter.Contract.Positions(&_CompoundV3Adapter.CallOpts, arg0)
}

// StateAggregator is a free data retrieval call binding the contract method 0x54ef2920.
//
// Solidity: function stateAggregator() view returns(address)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) StateAggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "stateAggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StateAggregator is a free data retrieval call binding the contract method 0x54ef2920.
//
// Solidity: function stateAggregator() view returns(address)
func (_CompoundV3Adapter *CompoundV3AdapterSession) StateAggregator() (common.Address, error) {
	return _CompoundV3Adapter.Contract.StateAggregator(&_CompoundV3Adapter.CallOpts)
}

// StateAggregator is a free data retrieval call binding the contract method 0x54ef2920.
//
// Solidity: function stateAggregator() view returns(address)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) StateAggregator() (common.Address, error) {
	return _CompoundV3Adapter.Contract.StateAggregator(&_CompoundV3Adapter.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CompoundV3Adapter *CompoundV3AdapterSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CompoundV3Adapter.Contract.SupportsInterface(&_CompoundV3Adapter.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CompoundV3Adapter.Contract.SupportsInterface(&_CompoundV3Adapter.CallOpts, interfaceId)
}

// TotalBorrowed is a free data retrieval call binding the contract method 0x4c19386c.
//
// Solidity: function totalBorrowed() view returns(uint256)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) TotalBorrowed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "totalBorrowed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalBorrowed is a free data retrieval call binding the contract method 0x4c19386c.
//
// Solidity: function totalBorrowed() view returns(uint256)
func (_CompoundV3Adapter *CompoundV3AdapterSession) TotalBorrowed() (*big.Int, error) {
	return _CompoundV3Adapter.Contract.TotalBorrowed(&_CompoundV3Adapter.CallOpts)
}

// TotalBorrowed is a free data retrieval call binding the contract method 0x4c19386c.
//
// Solidity: function totalBorrowed() view returns(uint256)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) TotalBorrowed() (*big.Int, error) {
	return _CompoundV3Adapter.Contract.TotalBorrowed(&_CompoundV3Adapter.CallOpts)
}

// TotalCollateral is a free data retrieval call binding the contract method 0x4ac8eb5f.
//
// Solidity: function totalCollateral() view returns(uint256)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) TotalCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "totalCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalCollateral is a free data retrieval call binding the contract method 0x4ac8eb5f.
//
// Solidity: function totalCollateral() view returns(uint256)
func (_CompoundV3Adapter *CompoundV3AdapterSession) TotalCollateral() (*big.Int, error) {
	return _CompoundV3Adapter.Contract.TotalCollateral(&_CompoundV3Adapter.CallOpts)
}

// TotalCollateral is a free data retrieval call binding the contract method 0x4ac8eb5f.
//
// Solidity: function totalCollateral() view returns(uint256)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) TotalCollateral() (*big.Int, error) {
	return _CompoundV3Adapter.Contract.TotalCollateral(&_CompoundV3Adapter.CallOpts)
}

// TotalSupplied is a free data retrieval call binding the contract method 0x630fd0ac.
//
// Solidity: function totalSupplied() view returns(uint256)
func (_CompoundV3Adapter *CompoundV3AdapterCaller) TotalSupplied(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CompoundV3Adapter.contract.Call(opts, &out, "totalSupplied")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupplied is a free data retrieval call binding the contract method 0x630fd0ac.
//
// Solidity: function totalSupplied() view returns(uint256)
func (_CompoundV3Adapter *CompoundV3AdapterSession) TotalSupplied() (*big.Int, error) {
	return _CompoundV3Adapter.Contract.TotalSupplied(&_CompoundV3Adapter.CallOpts)
}

// TotalSupplied is a free data retrieval call binding the contract method 0x630fd0ac.
//
// Solidity: function totalSupplied() view returns(uint256)
func (_CompoundV3Adapter *CompoundV3AdapterCallerSession) TotalSupplied() (*big.Int, error) {
	return _CompoundV3Adapter.Contract.TotalSupplied(&_CompoundV3Adapter.CallOpts)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) Borrow(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "borrow", amount)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) Borrow(amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.Borrow(&_CompoundV3Adapter.TransactOpts, amount)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) Borrow(amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.Borrow(&_CompoundV3Adapter.TransactOpts, amount)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x0e6878a3.
//
// Solidity: function claimRewards(bool shouldAccrue) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) ClaimRewards(opts *bind.TransactOpts, shouldAccrue bool) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "claimRewards", shouldAccrue)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x0e6878a3.
//
// Solidity: function claimRewards(bool shouldAccrue) returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) ClaimRewards(shouldAccrue bool) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.ClaimRewards(&_CompoundV3Adapter.TransactOpts, shouldAccrue)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x0e6878a3.
//
// Solidity: function claimRewards(bool shouldAccrue) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) ClaimRewards(shouldAccrue bool) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.ClaimRewards(&_CompoundV3Adapter.TransactOpts, shouldAccrue)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x95ccea67.
//
// Solidity: function emergencyWithdraw(address token, uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) EmergencyWithdraw(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "emergencyWithdraw", token, amount)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x95ccea67.
//
// Solidity: function emergencyWithdraw(address token, uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) EmergencyWithdraw(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.EmergencyWithdraw(&_CompoundV3Adapter.TransactOpts, token, amount)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x95ccea67.
//
// Solidity: function emergencyWithdraw(address token, uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) EmergencyWithdraw(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.EmergencyWithdraw(&_CompoundV3Adapter.TransactOpts, token, amount)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.GrantRole(&_CompoundV3Adapter.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.GrantRole(&_CompoundV3Adapter.TransactOpts, role, account)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) Pause() (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.Pause(&_CompoundV3Adapter.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) Pause() (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.Pause(&_CompoundV3Adapter.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.RenounceRole(&_CompoundV3Adapter.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.RenounceRole(&_CompoundV3Adapter.TransactOpts, role, callerConfirmation)
}

// Repay is a paid mutator transaction binding the contract method 0x371fd8e6.
//
// Solidity: function repay(uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) Repay(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "repay", amount)
}

// Repay is a paid mutator transaction binding the contract method 0x371fd8e6.
//
// Solidity: function repay(uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) Repay(amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.Repay(&_CompoundV3Adapter.TransactOpts, amount)
}

// Repay is a paid mutator transaction binding the contract method 0x371fd8e6.
//
// Solidity: function repay(uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) Repay(amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.Repay(&_CompoundV3Adapter.TransactOpts, amount)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.RevokeRole(&_CompoundV3Adapter.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.RevokeRole(&_CompoundV3Adapter.TransactOpts, role, account)
}

// SetStateAggregator is a paid mutator transaction binding the contract method 0xc0804b46.
//
// Solidity: function setStateAggregator(address _stateAggregator) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) SetStateAggregator(opts *bind.TransactOpts, _stateAggregator common.Address) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "setStateAggregator", _stateAggregator)
}

// SetStateAggregator is a paid mutator transaction binding the contract method 0xc0804b46.
//
// Solidity: function setStateAggregator(address _stateAggregator) returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) SetStateAggregator(_stateAggregator common.Address) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.SetStateAggregator(&_CompoundV3Adapter.TransactOpts, _stateAggregator)
}

// SetStateAggregator is a paid mutator transaction binding the contract method 0xc0804b46.
//
// Solidity: function setStateAggregator(address _stateAggregator) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) SetStateAggregator(_stateAggregator common.Address) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.SetStateAggregator(&_CompoundV3Adapter.TransactOpts, _stateAggregator)
}

// SupplyBase is a paid mutator transaction binding the contract method 0x5993a8b3.
//
// Solidity: function supplyBase(uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) SupplyBase(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "supplyBase", amount)
}

// SupplyBase is a paid mutator transaction binding the contract method 0x5993a8b3.
//
// Solidity: function supplyBase(uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) SupplyBase(amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.SupplyBase(&_CompoundV3Adapter.TransactOpts, amount)
}

// SupplyBase is a paid mutator transaction binding the contract method 0x5993a8b3.
//
// Solidity: function supplyBase(uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) SupplyBase(amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.SupplyBase(&_CompoundV3Adapter.TransactOpts, amount)
}

// SupplyCollateral is a paid mutator transaction binding the contract method 0xd2a8607b.
//
// Solidity: function supplyCollateral(address asset, uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) SupplyCollateral(opts *bind.TransactOpts, asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "supplyCollateral", asset, amount)
}

// SupplyCollateral is a paid mutator transaction binding the contract method 0xd2a8607b.
//
// Solidity: function supplyCollateral(address asset, uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) SupplyCollateral(asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.SupplyCollateral(&_CompoundV3Adapter.TransactOpts, asset, amount)
}

// SupplyCollateral is a paid mutator transaction binding the contract method 0xd2a8607b.
//
// Solidity: function supplyCollateral(address asset, uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) SupplyCollateral(asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.SupplyCollateral(&_CompoundV3Adapter.TransactOpts, asset, amount)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) Unpause() (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.Unpause(&_CompoundV3Adapter.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) Unpause() (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.Unpause(&_CompoundV3Adapter.TransactOpts)
}

// WithdrawBase is a paid mutator transaction binding the contract method 0xf98bea15.
//
// Solidity: function withdrawBase(uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) WithdrawBase(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "withdrawBase", amount)
}

// WithdrawBase is a paid mutator transaction binding the contract method 0xf98bea15.
//
// Solidity: function withdrawBase(uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) WithdrawBase(amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.WithdrawBase(&_CompoundV3Adapter.TransactOpts, amount)
}

// WithdrawBase is a paid mutator transaction binding the contract method 0xf98bea15.
//
// Solidity: function withdrawBase(uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) WithdrawBase(amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.WithdrawBase(&_CompoundV3Adapter.TransactOpts, amount)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x350c35e9.
//
// Solidity: function withdrawCollateral(address asset, uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) WithdrawCollateral(opts *bind.TransactOpts, asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.Transact(opts, "withdrawCollateral", asset, amount)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x350c35e9.
//
// Solidity: function withdrawCollateral(address asset, uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) WithdrawCollateral(asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.WithdrawCollateral(&_CompoundV3Adapter.TransactOpts, asset, amount)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x350c35e9.
//
// Solidity: function withdrawCollateral(address asset, uint256 amount) returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) WithdrawCollateral(asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.WithdrawCollateral(&_CompoundV3Adapter.TransactOpts, asset, amount)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompoundV3Adapter.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_CompoundV3Adapter *CompoundV3AdapterSession) Receive() (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.Receive(&_CompoundV3Adapter.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_CompoundV3Adapter *CompoundV3AdapterTransactorSession) Receive() (*types.Transaction, error) {
	return _CompoundV3Adapter.Contract.Receive(&_CompoundV3Adapter.TransactOpts)
}

// CompoundV3AdapterBaseSuppliedIterator is returned from FilterBaseSupplied and is used to iterate over the raw logs and unpacked data for BaseSupplied events raised by the CompoundV3Adapter contract.
type CompoundV3AdapterBaseSuppliedIterator struct {
	Event *CompoundV3AdapterBaseSupplied // Event containing the contract specifics and raw log

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
func (it *CompoundV3AdapterBaseSuppliedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundV3AdapterBaseSupplied)
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
		it.Event = new(CompoundV3AdapterBaseSupplied)
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
func (it *CompoundV3AdapterBaseSuppliedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundV3AdapterBaseSuppliedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundV3AdapterBaseSupplied represents a BaseSupplied event raised by the CompoundV3Adapter contract.
type CompoundV3AdapterBaseSupplied struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBaseSupplied is a free log retrieval operation binding the contract event 0x404c0cabb6935529d6ab61c3b44d8100ead2c25cfcfe1c3bcff4f39b9b90c207.
//
// Solidity: event BaseSupplied(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) FilterBaseSupplied(opts *bind.FilterOpts, user []common.Address) (*CompoundV3AdapterBaseSuppliedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.FilterLogs(opts, "BaseSupplied", userRule)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterBaseSuppliedIterator{contract: _CompoundV3Adapter.contract, event: "BaseSupplied", logs: logs, sub: sub}, nil
}

// WatchBaseSupplied is a free log subscription operation binding the contract event 0x404c0cabb6935529d6ab61c3b44d8100ead2c25cfcfe1c3bcff4f39b9b90c207.
//
// Solidity: event BaseSupplied(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) WatchBaseSupplied(opts *bind.WatchOpts, sink chan<- *CompoundV3AdapterBaseSupplied, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.WatchLogs(opts, "BaseSupplied", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundV3AdapterBaseSupplied)
				if err := _CompoundV3Adapter.contract.UnpackLog(event, "BaseSupplied", log); err != nil {
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

// ParseBaseSupplied is a log parse operation binding the contract event 0x404c0cabb6935529d6ab61c3b44d8100ead2c25cfcfe1c3bcff4f39b9b90c207.
//
// Solidity: event BaseSupplied(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) ParseBaseSupplied(log types.Log) (*CompoundV3AdapterBaseSupplied, error) {
	event := new(CompoundV3AdapterBaseSupplied)
	if err := _CompoundV3Adapter.contract.UnpackLog(event, "BaseSupplied", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompoundV3AdapterBaseWithdrawnIterator is returned from FilterBaseWithdrawn and is used to iterate over the raw logs and unpacked data for BaseWithdrawn events raised by the CompoundV3Adapter contract.
type CompoundV3AdapterBaseWithdrawnIterator struct {
	Event *CompoundV3AdapterBaseWithdrawn // Event containing the contract specifics and raw log

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
func (it *CompoundV3AdapterBaseWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundV3AdapterBaseWithdrawn)
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
		it.Event = new(CompoundV3AdapterBaseWithdrawn)
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
func (it *CompoundV3AdapterBaseWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundV3AdapterBaseWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundV3AdapterBaseWithdrawn represents a BaseWithdrawn event raised by the CompoundV3Adapter contract.
type CompoundV3AdapterBaseWithdrawn struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBaseWithdrawn is a free log retrieval operation binding the contract event 0x24bf52d87309c5594d8e64b5b23596bffbe07396bec5b6b60aabd458a248062c.
//
// Solidity: event BaseWithdrawn(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) FilterBaseWithdrawn(opts *bind.FilterOpts, user []common.Address) (*CompoundV3AdapterBaseWithdrawnIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.FilterLogs(opts, "BaseWithdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterBaseWithdrawnIterator{contract: _CompoundV3Adapter.contract, event: "BaseWithdrawn", logs: logs, sub: sub}, nil
}

// WatchBaseWithdrawn is a free log subscription operation binding the contract event 0x24bf52d87309c5594d8e64b5b23596bffbe07396bec5b6b60aabd458a248062c.
//
// Solidity: event BaseWithdrawn(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) WatchBaseWithdrawn(opts *bind.WatchOpts, sink chan<- *CompoundV3AdapterBaseWithdrawn, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.WatchLogs(opts, "BaseWithdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundV3AdapterBaseWithdrawn)
				if err := _CompoundV3Adapter.contract.UnpackLog(event, "BaseWithdrawn", log); err != nil {
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

// ParseBaseWithdrawn is a log parse operation binding the contract event 0x24bf52d87309c5594d8e64b5b23596bffbe07396bec5b6b60aabd458a248062c.
//
// Solidity: event BaseWithdrawn(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) ParseBaseWithdrawn(log types.Log) (*CompoundV3AdapterBaseWithdrawn, error) {
	event := new(CompoundV3AdapterBaseWithdrawn)
	if err := _CompoundV3Adapter.contract.UnpackLog(event, "BaseWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompoundV3AdapterBorrowedIterator is returned from FilterBorrowed and is used to iterate over the raw logs and unpacked data for Borrowed events raised by the CompoundV3Adapter contract.
type CompoundV3AdapterBorrowedIterator struct {
	Event *CompoundV3AdapterBorrowed // Event containing the contract specifics and raw log

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
func (it *CompoundV3AdapterBorrowedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundV3AdapterBorrowed)
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
		it.Event = new(CompoundV3AdapterBorrowed)
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
func (it *CompoundV3AdapterBorrowedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundV3AdapterBorrowedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundV3AdapterBorrowed represents a Borrowed event raised by the CompoundV3Adapter contract.
type CompoundV3AdapterBorrowed struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBorrowed is a free log retrieval operation binding the contract event 0xac59582e5396aca512fa873a2047e7f4c80f8f55d4a06cb34a78a0187f62719f.
//
// Solidity: event Borrowed(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) FilterBorrowed(opts *bind.FilterOpts, user []common.Address) (*CompoundV3AdapterBorrowedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.FilterLogs(opts, "Borrowed", userRule)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterBorrowedIterator{contract: _CompoundV3Adapter.contract, event: "Borrowed", logs: logs, sub: sub}, nil
}

// WatchBorrowed is a free log subscription operation binding the contract event 0xac59582e5396aca512fa873a2047e7f4c80f8f55d4a06cb34a78a0187f62719f.
//
// Solidity: event Borrowed(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) WatchBorrowed(opts *bind.WatchOpts, sink chan<- *CompoundV3AdapterBorrowed, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.WatchLogs(opts, "Borrowed", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundV3AdapterBorrowed)
				if err := _CompoundV3Adapter.contract.UnpackLog(event, "Borrowed", log); err != nil {
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

// ParseBorrowed is a log parse operation binding the contract event 0xac59582e5396aca512fa873a2047e7f4c80f8f55d4a06cb34a78a0187f62719f.
//
// Solidity: event Borrowed(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) ParseBorrowed(log types.Log) (*CompoundV3AdapterBorrowed, error) {
	event := new(CompoundV3AdapterBorrowed)
	if err := _CompoundV3Adapter.contract.UnpackLog(event, "Borrowed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompoundV3AdapterCollateralSuppliedIterator is returned from FilterCollateralSupplied and is used to iterate over the raw logs and unpacked data for CollateralSupplied events raised by the CompoundV3Adapter contract.
type CompoundV3AdapterCollateralSuppliedIterator struct {
	Event *CompoundV3AdapterCollateralSupplied // Event containing the contract specifics and raw log

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
func (it *CompoundV3AdapterCollateralSuppliedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundV3AdapterCollateralSupplied)
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
		it.Event = new(CompoundV3AdapterCollateralSupplied)
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
func (it *CompoundV3AdapterCollateralSuppliedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundV3AdapterCollateralSuppliedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundV3AdapterCollateralSupplied represents a CollateralSupplied event raised by the CompoundV3Adapter contract.
type CompoundV3AdapterCollateralSupplied struct {
	User   common.Address
	Asset  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCollateralSupplied is a free log retrieval operation binding the contract event 0xa1d3b269caaa0f650745881dd7077600f8c21e0f1ee4c7292c13543f4d89c7b1.
//
// Solidity: event CollateralSupplied(address indexed user, address indexed asset, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) FilterCollateralSupplied(opts *bind.FilterOpts, user []common.Address, asset []common.Address) (*CompoundV3AdapterCollateralSuppliedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.FilterLogs(opts, "CollateralSupplied", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterCollateralSuppliedIterator{contract: _CompoundV3Adapter.contract, event: "CollateralSupplied", logs: logs, sub: sub}, nil
}

// WatchCollateralSupplied is a free log subscription operation binding the contract event 0xa1d3b269caaa0f650745881dd7077600f8c21e0f1ee4c7292c13543f4d89c7b1.
//
// Solidity: event CollateralSupplied(address indexed user, address indexed asset, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) WatchCollateralSupplied(opts *bind.WatchOpts, sink chan<- *CompoundV3AdapterCollateralSupplied, user []common.Address, asset []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.WatchLogs(opts, "CollateralSupplied", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundV3AdapterCollateralSupplied)
				if err := _CompoundV3Adapter.contract.UnpackLog(event, "CollateralSupplied", log); err != nil {
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

// ParseCollateralSupplied is a log parse operation binding the contract event 0xa1d3b269caaa0f650745881dd7077600f8c21e0f1ee4c7292c13543f4d89c7b1.
//
// Solidity: event CollateralSupplied(address indexed user, address indexed asset, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) ParseCollateralSupplied(log types.Log) (*CompoundV3AdapterCollateralSupplied, error) {
	event := new(CompoundV3AdapterCollateralSupplied)
	if err := _CompoundV3Adapter.contract.UnpackLog(event, "CollateralSupplied", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompoundV3AdapterCollateralWithdrawnIterator is returned from FilterCollateralWithdrawn and is used to iterate over the raw logs and unpacked data for CollateralWithdrawn events raised by the CompoundV3Adapter contract.
type CompoundV3AdapterCollateralWithdrawnIterator struct {
	Event *CompoundV3AdapterCollateralWithdrawn // Event containing the contract specifics and raw log

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
func (it *CompoundV3AdapterCollateralWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundV3AdapterCollateralWithdrawn)
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
		it.Event = new(CompoundV3AdapterCollateralWithdrawn)
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
func (it *CompoundV3AdapterCollateralWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundV3AdapterCollateralWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundV3AdapterCollateralWithdrawn represents a CollateralWithdrawn event raised by the CompoundV3Adapter contract.
type CompoundV3AdapterCollateralWithdrawn struct {
	User   common.Address
	Asset  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCollateralWithdrawn is a free log retrieval operation binding the contract event 0x45892a46e6cef329bb642da6d69846d324db43d19008edc141ed82382eda1bee.
//
// Solidity: event CollateralWithdrawn(address indexed user, address indexed asset, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) FilterCollateralWithdrawn(opts *bind.FilterOpts, user []common.Address, asset []common.Address) (*CompoundV3AdapterCollateralWithdrawnIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.FilterLogs(opts, "CollateralWithdrawn", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterCollateralWithdrawnIterator{contract: _CompoundV3Adapter.contract, event: "CollateralWithdrawn", logs: logs, sub: sub}, nil
}

// WatchCollateralWithdrawn is a free log subscription operation binding the contract event 0x45892a46e6cef329bb642da6d69846d324db43d19008edc141ed82382eda1bee.
//
// Solidity: event CollateralWithdrawn(address indexed user, address indexed asset, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) WatchCollateralWithdrawn(opts *bind.WatchOpts, sink chan<- *CompoundV3AdapterCollateralWithdrawn, user []common.Address, asset []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.WatchLogs(opts, "CollateralWithdrawn", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundV3AdapterCollateralWithdrawn)
				if err := _CompoundV3Adapter.contract.UnpackLog(event, "CollateralWithdrawn", log); err != nil {
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

// ParseCollateralWithdrawn is a log parse operation binding the contract event 0x45892a46e6cef329bb642da6d69846d324db43d19008edc141ed82382eda1bee.
//
// Solidity: event CollateralWithdrawn(address indexed user, address indexed asset, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) ParseCollateralWithdrawn(log types.Log) (*CompoundV3AdapterCollateralWithdrawn, error) {
	event := new(CompoundV3AdapterCollateralWithdrawn)
	if err := _CompoundV3Adapter.contract.UnpackLog(event, "CollateralWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompoundV3AdapterPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the CompoundV3Adapter contract.
type CompoundV3AdapterPausedIterator struct {
	Event *CompoundV3AdapterPaused // Event containing the contract specifics and raw log

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
func (it *CompoundV3AdapterPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundV3AdapterPaused)
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
		it.Event = new(CompoundV3AdapterPaused)
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
func (it *CompoundV3AdapterPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundV3AdapterPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundV3AdapterPaused represents a Paused event raised by the CompoundV3Adapter contract.
type CompoundV3AdapterPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) FilterPaused(opts *bind.FilterOpts) (*CompoundV3AdapterPausedIterator, error) {

	logs, sub, err := _CompoundV3Adapter.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterPausedIterator{contract: _CompoundV3Adapter.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *CompoundV3AdapterPaused) (event.Subscription, error) {

	logs, sub, err := _CompoundV3Adapter.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundV3AdapterPaused)
				if err := _CompoundV3Adapter.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) ParsePaused(log types.Log) (*CompoundV3AdapterPaused, error) {
	event := new(CompoundV3AdapterPaused)
	if err := _CompoundV3Adapter.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompoundV3AdapterRepaidIterator is returned from FilterRepaid and is used to iterate over the raw logs and unpacked data for Repaid events raised by the CompoundV3Adapter contract.
type CompoundV3AdapterRepaidIterator struct {
	Event *CompoundV3AdapterRepaid // Event containing the contract specifics and raw log

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
func (it *CompoundV3AdapterRepaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundV3AdapterRepaid)
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
		it.Event = new(CompoundV3AdapterRepaid)
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
func (it *CompoundV3AdapterRepaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundV3AdapterRepaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundV3AdapterRepaid represents a Repaid event raised by the CompoundV3Adapter contract.
type CompoundV3AdapterRepaid struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRepaid is a free log retrieval operation binding the contract event 0x0516911bcc3a0a7412a44601057c0a0a1ec628bde049a84284bc428866534488.
//
// Solidity: event Repaid(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) FilterRepaid(opts *bind.FilterOpts, user []common.Address) (*CompoundV3AdapterRepaidIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.FilterLogs(opts, "Repaid", userRule)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterRepaidIterator{contract: _CompoundV3Adapter.contract, event: "Repaid", logs: logs, sub: sub}, nil
}

// WatchRepaid is a free log subscription operation binding the contract event 0x0516911bcc3a0a7412a44601057c0a0a1ec628bde049a84284bc428866534488.
//
// Solidity: event Repaid(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) WatchRepaid(opts *bind.WatchOpts, sink chan<- *CompoundV3AdapterRepaid, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.WatchLogs(opts, "Repaid", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundV3AdapterRepaid)
				if err := _CompoundV3Adapter.contract.UnpackLog(event, "Repaid", log); err != nil {
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

// ParseRepaid is a log parse operation binding the contract event 0x0516911bcc3a0a7412a44601057c0a0a1ec628bde049a84284bc428866534488.
//
// Solidity: event Repaid(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) ParseRepaid(log types.Log) (*CompoundV3AdapterRepaid, error) {
	event := new(CompoundV3AdapterRepaid)
	if err := _CompoundV3Adapter.contract.UnpackLog(event, "Repaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompoundV3AdapterRewardsClaimedIterator is returned from FilterRewardsClaimed and is used to iterate over the raw logs and unpacked data for RewardsClaimed events raised by the CompoundV3Adapter contract.
type CompoundV3AdapterRewardsClaimedIterator struct {
	Event *CompoundV3AdapterRewardsClaimed // Event containing the contract specifics and raw log

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
func (it *CompoundV3AdapterRewardsClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundV3AdapterRewardsClaimed)
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
		it.Event = new(CompoundV3AdapterRewardsClaimed)
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
func (it *CompoundV3AdapterRewardsClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundV3AdapterRewardsClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundV3AdapterRewardsClaimed represents a RewardsClaimed event raised by the CompoundV3Adapter contract.
type CompoundV3AdapterRewardsClaimed struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardsClaimed is a free log retrieval operation binding the contract event 0xfc30cddea38e2bf4d6ea7d3f9ed3b6ad7f176419f4963bd81318067a4aee73fe.
//
// Solidity: event RewardsClaimed(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) FilterRewardsClaimed(opts *bind.FilterOpts, user []common.Address) (*CompoundV3AdapterRewardsClaimedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.FilterLogs(opts, "RewardsClaimed", userRule)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterRewardsClaimedIterator{contract: _CompoundV3Adapter.contract, event: "RewardsClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardsClaimed is a free log subscription operation binding the contract event 0xfc30cddea38e2bf4d6ea7d3f9ed3b6ad7f176419f4963bd81318067a4aee73fe.
//
// Solidity: event RewardsClaimed(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) WatchRewardsClaimed(opts *bind.WatchOpts, sink chan<- *CompoundV3AdapterRewardsClaimed, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.WatchLogs(opts, "RewardsClaimed", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundV3AdapterRewardsClaimed)
				if err := _CompoundV3Adapter.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
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

// ParseRewardsClaimed is a log parse operation binding the contract event 0xfc30cddea38e2bf4d6ea7d3f9ed3b6ad7f176419f4963bd81318067a4aee73fe.
//
// Solidity: event RewardsClaimed(address indexed user, uint256 amount)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) ParseRewardsClaimed(log types.Log) (*CompoundV3AdapterRewardsClaimed, error) {
	event := new(CompoundV3AdapterRewardsClaimed)
	if err := _CompoundV3Adapter.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompoundV3AdapterRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the CompoundV3Adapter contract.
type CompoundV3AdapterRoleAdminChangedIterator struct {
	Event *CompoundV3AdapterRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *CompoundV3AdapterRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundV3AdapterRoleAdminChanged)
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
		it.Event = new(CompoundV3AdapterRoleAdminChanged)
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
func (it *CompoundV3AdapterRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundV3AdapterRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundV3AdapterRoleAdminChanged represents a RoleAdminChanged event raised by the CompoundV3Adapter contract.
type CompoundV3AdapterRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*CompoundV3AdapterRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterRoleAdminChangedIterator{contract: _CompoundV3Adapter.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *CompoundV3AdapterRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundV3AdapterRoleAdminChanged)
				if err := _CompoundV3Adapter.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) ParseRoleAdminChanged(log types.Log) (*CompoundV3AdapterRoleAdminChanged, error) {
	event := new(CompoundV3AdapterRoleAdminChanged)
	if err := _CompoundV3Adapter.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompoundV3AdapterRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the CompoundV3Adapter contract.
type CompoundV3AdapterRoleGrantedIterator struct {
	Event *CompoundV3AdapterRoleGranted // Event containing the contract specifics and raw log

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
func (it *CompoundV3AdapterRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundV3AdapterRoleGranted)
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
		it.Event = new(CompoundV3AdapterRoleGranted)
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
func (it *CompoundV3AdapterRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundV3AdapterRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundV3AdapterRoleGranted represents a RoleGranted event raised by the CompoundV3Adapter contract.
type CompoundV3AdapterRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*CompoundV3AdapterRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterRoleGrantedIterator{contract: _CompoundV3Adapter.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *CompoundV3AdapterRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundV3AdapterRoleGranted)
				if err := _CompoundV3Adapter.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) ParseRoleGranted(log types.Log) (*CompoundV3AdapterRoleGranted, error) {
	event := new(CompoundV3AdapterRoleGranted)
	if err := _CompoundV3Adapter.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompoundV3AdapterRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the CompoundV3Adapter contract.
type CompoundV3AdapterRoleRevokedIterator struct {
	Event *CompoundV3AdapterRoleRevoked // Event containing the contract specifics and raw log

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
func (it *CompoundV3AdapterRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundV3AdapterRoleRevoked)
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
		it.Event = new(CompoundV3AdapterRoleRevoked)
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
func (it *CompoundV3AdapterRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundV3AdapterRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundV3AdapterRoleRevoked represents a RoleRevoked event raised by the CompoundV3Adapter contract.
type CompoundV3AdapterRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*CompoundV3AdapterRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterRoleRevokedIterator{contract: _CompoundV3Adapter.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *CompoundV3AdapterRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _CompoundV3Adapter.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundV3AdapterRoleRevoked)
				if err := _CompoundV3Adapter.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) ParseRoleRevoked(log types.Log) (*CompoundV3AdapterRoleRevoked, error) {
	event := new(CompoundV3AdapterRoleRevoked)
	if err := _CompoundV3Adapter.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CompoundV3AdapterUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the CompoundV3Adapter contract.
type CompoundV3AdapterUnpausedIterator struct {
	Event *CompoundV3AdapterUnpaused // Event containing the contract specifics and raw log

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
func (it *CompoundV3AdapterUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompoundV3AdapterUnpaused)
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
		it.Event = new(CompoundV3AdapterUnpaused)
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
func (it *CompoundV3AdapterUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompoundV3AdapterUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompoundV3AdapterUnpaused represents a Unpaused event raised by the CompoundV3Adapter contract.
type CompoundV3AdapterUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) FilterUnpaused(opts *bind.FilterOpts) (*CompoundV3AdapterUnpausedIterator, error) {

	logs, sub, err := _CompoundV3Adapter.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &CompoundV3AdapterUnpausedIterator{contract: _CompoundV3Adapter.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *CompoundV3AdapterUnpaused) (event.Subscription, error) {

	logs, sub, err := _CompoundV3Adapter.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompoundV3AdapterUnpaused)
				if err := _CompoundV3Adapter.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_CompoundV3Adapter *CompoundV3AdapterFilterer) ParseUnpaused(log types.Log) (*CompoundV3AdapterUnpaused, error) {
	event := new(CompoundV3AdapterUnpaused)
	if err := _CompoundV3Adapter.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
