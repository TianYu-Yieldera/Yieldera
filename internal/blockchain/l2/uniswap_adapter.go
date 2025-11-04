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

// UniswapV3AdapterMetaData contains all meta data concerning the UniswapV3Adapter contract.
var UniswapV3AdapterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_swapRouter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"MultiHopSwap\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"}],\"name\":\"Swapped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_HIGH\",\"outputs\":[{\"internalType\":\"uint24\",\"name\":\"\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_LOW\",\"outputs\":[{\"internalType\":\"uint24\",\"name\":\"\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_LOWEST\",\"outputs\":[{\"internalType\":\"uint24\",\"name\":\"\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_MEDIUM\",\"outputs\":[{\"internalType\":\"uint24\",\"name\":\"\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MANAGER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emergencyWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"contractIUniswapV3Factory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"}],\"name\":\"getPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"}],\"name\":\"poolExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stateAggregator\",\"type\":\"address\"}],\"name\":\"setStateAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateAggregator\",\"outputs\":[{\"internalType\":\"contractL2StateAggregator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactInputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMaximum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactOutputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapRouter\",\"outputs\":[{\"internalType\":\"contractISwapRouter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSwapVolume\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSwaps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// UniswapV3AdapterABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapV3AdapterMetaData.ABI instead.
var UniswapV3AdapterABI = UniswapV3AdapterMetaData.ABI

// UniswapV3Adapter is an auto generated Go binding around an Ethereum contract.
type UniswapV3Adapter struct {
	UniswapV3AdapterCaller     // Read-only binding to the contract
	UniswapV3AdapterTransactor // Write-only binding to the contract
	UniswapV3AdapterFilterer   // Log filterer for contract events
}

// UniswapV3AdapterCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniswapV3AdapterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV3AdapterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniswapV3AdapterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV3AdapterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniswapV3AdapterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV3AdapterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniswapV3AdapterSession struct {
	Contract     *UniswapV3Adapter // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UniswapV3AdapterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniswapV3AdapterCallerSession struct {
	Contract *UniswapV3AdapterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// UniswapV3AdapterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniswapV3AdapterTransactorSession struct {
	Contract     *UniswapV3AdapterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// UniswapV3AdapterRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniswapV3AdapterRaw struct {
	Contract *UniswapV3Adapter // Generic contract binding to access the raw methods on
}

// UniswapV3AdapterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniswapV3AdapterCallerRaw struct {
	Contract *UniswapV3AdapterCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapV3AdapterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniswapV3AdapterTransactorRaw struct {
	Contract *UniswapV3AdapterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapV3Adapter creates a new instance of UniswapV3Adapter, bound to a specific deployed contract.
func NewUniswapV3Adapter(address common.Address, backend bind.ContractBackend) (*UniswapV3Adapter, error) {
	contract, err := bindUniswapV3Adapter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapV3Adapter{UniswapV3AdapterCaller: UniswapV3AdapterCaller{contract: contract}, UniswapV3AdapterTransactor: UniswapV3AdapterTransactor{contract: contract}, UniswapV3AdapterFilterer: UniswapV3AdapterFilterer{contract: contract}}, nil
}

// NewUniswapV3AdapterCaller creates a new read-only instance of UniswapV3Adapter, bound to a specific deployed contract.
func NewUniswapV3AdapterCaller(address common.Address, caller bind.ContractCaller) (*UniswapV3AdapterCaller, error) {
	contract, err := bindUniswapV3Adapter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV3AdapterCaller{contract: contract}, nil
}

// NewUniswapV3AdapterTransactor creates a new write-only instance of UniswapV3Adapter, bound to a specific deployed contract.
func NewUniswapV3AdapterTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapV3AdapterTransactor, error) {
	contract, err := bindUniswapV3Adapter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV3AdapterTransactor{contract: contract}, nil
}

// NewUniswapV3AdapterFilterer creates a new log filterer instance of UniswapV3Adapter, bound to a specific deployed contract.
func NewUniswapV3AdapterFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapV3AdapterFilterer, error) {
	contract, err := bindUniswapV3Adapter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapV3AdapterFilterer{contract: contract}, nil
}

// bindUniswapV3Adapter binds a generic wrapper to an already deployed contract.
func bindUniswapV3Adapter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniswapV3AdapterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV3Adapter *UniswapV3AdapterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV3Adapter.Contract.UniswapV3AdapterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV3Adapter *UniswapV3AdapterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.UniswapV3AdapterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV3Adapter *UniswapV3AdapterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.UniswapV3AdapterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV3Adapter *UniswapV3AdapterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV3Adapter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV3Adapter *UniswapV3AdapterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV3Adapter *UniswapV3AdapterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_UniswapV3Adapter *UniswapV3AdapterSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _UniswapV3Adapter.Contract.DEFAULTADMINROLE(&_UniswapV3Adapter.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _UniswapV3Adapter.Contract.DEFAULTADMINROLE(&_UniswapV3Adapter.CallOpts)
}

// FEEHIGH is a free data retrieval call binding the contract method 0x63075eb9.
//
// Solidity: function FEE_HIGH() view returns(uint24)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) FEEHIGH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "FEE_HIGH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FEEHIGH is a free data retrieval call binding the contract method 0x63075eb9.
//
// Solidity: function FEE_HIGH() view returns(uint24)
func (_UniswapV3Adapter *UniswapV3AdapterSession) FEEHIGH() (*big.Int, error) {
	return _UniswapV3Adapter.Contract.FEEHIGH(&_UniswapV3Adapter.CallOpts)
}

// FEEHIGH is a free data retrieval call binding the contract method 0x63075eb9.
//
// Solidity: function FEE_HIGH() view returns(uint24)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) FEEHIGH() (*big.Int, error) {
	return _UniswapV3Adapter.Contract.FEEHIGH(&_UniswapV3Adapter.CallOpts)
}

// FEELOW is a free data retrieval call binding the contract method 0x3add5c05.
//
// Solidity: function FEE_LOW() view returns(uint24)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) FEELOW(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "FEE_LOW")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FEELOW is a free data retrieval call binding the contract method 0x3add5c05.
//
// Solidity: function FEE_LOW() view returns(uint24)
func (_UniswapV3Adapter *UniswapV3AdapterSession) FEELOW() (*big.Int, error) {
	return _UniswapV3Adapter.Contract.FEELOW(&_UniswapV3Adapter.CallOpts)
}

// FEELOW is a free data retrieval call binding the contract method 0x3add5c05.
//
// Solidity: function FEE_LOW() view returns(uint24)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) FEELOW() (*big.Int, error) {
	return _UniswapV3Adapter.Contract.FEELOW(&_UniswapV3Adapter.CallOpts)
}

// FEELOWEST is a free data retrieval call binding the contract method 0x79dc9b32.
//
// Solidity: function FEE_LOWEST() view returns(uint24)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) FEELOWEST(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "FEE_LOWEST")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FEELOWEST is a free data retrieval call binding the contract method 0x79dc9b32.
//
// Solidity: function FEE_LOWEST() view returns(uint24)
func (_UniswapV3Adapter *UniswapV3AdapterSession) FEELOWEST() (*big.Int, error) {
	return _UniswapV3Adapter.Contract.FEELOWEST(&_UniswapV3Adapter.CallOpts)
}

// FEELOWEST is a free data retrieval call binding the contract method 0x79dc9b32.
//
// Solidity: function FEE_LOWEST() view returns(uint24)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) FEELOWEST() (*big.Int, error) {
	return _UniswapV3Adapter.Contract.FEELOWEST(&_UniswapV3Adapter.CallOpts)
}

// FEEMEDIUM is a free data retrieval call binding the contract method 0xca216247.
//
// Solidity: function FEE_MEDIUM() view returns(uint24)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) FEEMEDIUM(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "FEE_MEDIUM")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FEEMEDIUM is a free data retrieval call binding the contract method 0xca216247.
//
// Solidity: function FEE_MEDIUM() view returns(uint24)
func (_UniswapV3Adapter *UniswapV3AdapterSession) FEEMEDIUM() (*big.Int, error) {
	return _UniswapV3Adapter.Contract.FEEMEDIUM(&_UniswapV3Adapter.CallOpts)
}

// FEEMEDIUM is a free data retrieval call binding the contract method 0xca216247.
//
// Solidity: function FEE_MEDIUM() view returns(uint24)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) FEEMEDIUM() (*big.Int, error) {
	return _UniswapV3Adapter.Contract.FEEMEDIUM(&_UniswapV3Adapter.CallOpts)
}

// MANAGERROLE is a free data retrieval call binding the contract method 0xec87621c.
//
// Solidity: function MANAGER_ROLE() view returns(bytes32)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) MANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MANAGERROLE is a free data retrieval call binding the contract method 0xec87621c.
//
// Solidity: function MANAGER_ROLE() view returns(bytes32)
func (_UniswapV3Adapter *UniswapV3AdapterSession) MANAGERROLE() ([32]byte, error) {
	return _UniswapV3Adapter.Contract.MANAGERROLE(&_UniswapV3Adapter.CallOpts)
}

// MANAGERROLE is a free data retrieval call binding the contract method 0xec87621c.
//
// Solidity: function MANAGER_ROLE() view returns(bytes32)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) MANAGERROLE() ([32]byte, error) {
	return _UniswapV3Adapter.Contract.MANAGERROLE(&_UniswapV3Adapter.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_UniswapV3Adapter *UniswapV3AdapterSession) Factory() (common.Address, error) {
	return _UniswapV3Adapter.Contract.Factory(&_UniswapV3Adapter.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) Factory() (common.Address, error) {
	return _UniswapV3Adapter.Contract.Factory(&_UniswapV3Adapter.CallOpts)
}

// GetPool is a free data retrieval call binding the contract method 0x1698ee82.
//
// Solidity: function getPool(address tokenA, address tokenB, uint24 fee) view returns(address pool)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) GetPool(opts *bind.CallOpts, tokenA common.Address, tokenB common.Address, fee *big.Int) (common.Address, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "getPool", tokenA, tokenB, fee)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPool is a free data retrieval call binding the contract method 0x1698ee82.
//
// Solidity: function getPool(address tokenA, address tokenB, uint24 fee) view returns(address pool)
func (_UniswapV3Adapter *UniswapV3AdapterSession) GetPool(tokenA common.Address, tokenB common.Address, fee *big.Int) (common.Address, error) {
	return _UniswapV3Adapter.Contract.GetPool(&_UniswapV3Adapter.CallOpts, tokenA, tokenB, fee)
}

// GetPool is a free data retrieval call binding the contract method 0x1698ee82.
//
// Solidity: function getPool(address tokenA, address tokenB, uint24 fee) view returns(address pool)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) GetPool(tokenA common.Address, tokenB common.Address, fee *big.Int) (common.Address, error) {
	return _UniswapV3Adapter.Contract.GetPool(&_UniswapV3Adapter.CallOpts, tokenA, tokenB, fee)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_UniswapV3Adapter *UniswapV3AdapterSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _UniswapV3Adapter.Contract.GetRoleAdmin(&_UniswapV3Adapter.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _UniswapV3Adapter.Contract.GetRoleAdmin(&_UniswapV3Adapter.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_UniswapV3Adapter *UniswapV3AdapterSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _UniswapV3Adapter.Contract.HasRole(&_UniswapV3Adapter.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _UniswapV3Adapter.Contract.HasRole(&_UniswapV3Adapter.CallOpts, role, account)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_UniswapV3Adapter *UniswapV3AdapterSession) Paused() (bool, error) {
	return _UniswapV3Adapter.Contract.Paused(&_UniswapV3Adapter.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) Paused() (bool, error) {
	return _UniswapV3Adapter.Contract.Paused(&_UniswapV3Adapter.CallOpts)
}

// PoolExists is a free data retrieval call binding the contract method 0x806eacf8.
//
// Solidity: function poolExists(address tokenA, address tokenB, uint24 fee) view returns(bool exists)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) PoolExists(opts *bind.CallOpts, tokenA common.Address, tokenB common.Address, fee *big.Int) (bool, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "poolExists", tokenA, tokenB, fee)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// PoolExists is a free data retrieval call binding the contract method 0x806eacf8.
//
// Solidity: function poolExists(address tokenA, address tokenB, uint24 fee) view returns(bool exists)
func (_UniswapV3Adapter *UniswapV3AdapterSession) PoolExists(tokenA common.Address, tokenB common.Address, fee *big.Int) (bool, error) {
	return _UniswapV3Adapter.Contract.PoolExists(&_UniswapV3Adapter.CallOpts, tokenA, tokenB, fee)
}

// PoolExists is a free data retrieval call binding the contract method 0x806eacf8.
//
// Solidity: function poolExists(address tokenA, address tokenB, uint24 fee) view returns(bool exists)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) PoolExists(tokenA common.Address, tokenB common.Address, fee *big.Int) (bool, error) {
	return _UniswapV3Adapter.Contract.PoolExists(&_UniswapV3Adapter.CallOpts, tokenA, tokenB, fee)
}

// StateAggregator is a free data retrieval call binding the contract method 0x54ef2920.
//
// Solidity: function stateAggregator() view returns(address)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) StateAggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "stateAggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StateAggregator is a free data retrieval call binding the contract method 0x54ef2920.
//
// Solidity: function stateAggregator() view returns(address)
func (_UniswapV3Adapter *UniswapV3AdapterSession) StateAggregator() (common.Address, error) {
	return _UniswapV3Adapter.Contract.StateAggregator(&_UniswapV3Adapter.CallOpts)
}

// StateAggregator is a free data retrieval call binding the contract method 0x54ef2920.
//
// Solidity: function stateAggregator() view returns(address)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) StateAggregator() (common.Address, error) {
	return _UniswapV3Adapter.Contract.StateAggregator(&_UniswapV3Adapter.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_UniswapV3Adapter *UniswapV3AdapterSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _UniswapV3Adapter.Contract.SupportsInterface(&_UniswapV3Adapter.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _UniswapV3Adapter.Contract.SupportsInterface(&_UniswapV3Adapter.CallOpts, interfaceId)
}

// SwapRouter is a free data retrieval call binding the contract method 0xc31c9c07.
//
// Solidity: function swapRouter() view returns(address)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) SwapRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "swapRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SwapRouter is a free data retrieval call binding the contract method 0xc31c9c07.
//
// Solidity: function swapRouter() view returns(address)
func (_UniswapV3Adapter *UniswapV3AdapterSession) SwapRouter() (common.Address, error) {
	return _UniswapV3Adapter.Contract.SwapRouter(&_UniswapV3Adapter.CallOpts)
}

// SwapRouter is a free data retrieval call binding the contract method 0xc31c9c07.
//
// Solidity: function swapRouter() view returns(address)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) SwapRouter() (common.Address, error) {
	return _UniswapV3Adapter.Contract.SwapRouter(&_UniswapV3Adapter.CallOpts)
}

// TotalSwapVolume is a free data retrieval call binding the contract method 0x01b84065.
//
// Solidity: function totalSwapVolume() view returns(uint256)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) TotalSwapVolume(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "totalSwapVolume")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSwapVolume is a free data retrieval call binding the contract method 0x01b84065.
//
// Solidity: function totalSwapVolume() view returns(uint256)
func (_UniswapV3Adapter *UniswapV3AdapterSession) TotalSwapVolume() (*big.Int, error) {
	return _UniswapV3Adapter.Contract.TotalSwapVolume(&_UniswapV3Adapter.CallOpts)
}

// TotalSwapVolume is a free data retrieval call binding the contract method 0x01b84065.
//
// Solidity: function totalSwapVolume() view returns(uint256)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) TotalSwapVolume() (*big.Int, error) {
	return _UniswapV3Adapter.Contract.TotalSwapVolume(&_UniswapV3Adapter.CallOpts)
}

// TotalSwaps is a free data retrieval call binding the contract method 0xb4a800ce.
//
// Solidity: function totalSwaps() view returns(uint256)
func (_UniswapV3Adapter *UniswapV3AdapterCaller) TotalSwaps(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Adapter.contract.Call(opts, &out, "totalSwaps")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSwaps is a free data retrieval call binding the contract method 0xb4a800ce.
//
// Solidity: function totalSwaps() view returns(uint256)
func (_UniswapV3Adapter *UniswapV3AdapterSession) TotalSwaps() (*big.Int, error) {
	return _UniswapV3Adapter.Contract.TotalSwaps(&_UniswapV3Adapter.CallOpts)
}

// TotalSwaps is a free data retrieval call binding the contract method 0xb4a800ce.
//
// Solidity: function totalSwaps() view returns(uint256)
func (_UniswapV3Adapter *UniswapV3AdapterCallerSession) TotalSwaps() (*big.Int, error) {
	return _UniswapV3Adapter.Contract.TotalSwaps(&_UniswapV3Adapter.CallOpts)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x95ccea67.
//
// Solidity: function emergencyWithdraw(address token, uint256 amount) returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactor) EmergencyWithdraw(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _UniswapV3Adapter.contract.Transact(opts, "emergencyWithdraw", token, amount)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x95ccea67.
//
// Solidity: function emergencyWithdraw(address token, uint256 amount) returns()
func (_UniswapV3Adapter *UniswapV3AdapterSession) EmergencyWithdraw(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.EmergencyWithdraw(&_UniswapV3Adapter.TransactOpts, token, amount)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x95ccea67.
//
// Solidity: function emergencyWithdraw(address token, uint256 amount) returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactorSession) EmergencyWithdraw(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.EmergencyWithdraw(&_UniswapV3Adapter.TransactOpts, token, amount)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _UniswapV3Adapter.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_UniswapV3Adapter *UniswapV3AdapterSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.GrantRole(&_UniswapV3Adapter.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.GrantRole(&_UniswapV3Adapter.TransactOpts, role, account)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Adapter.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_UniswapV3Adapter *UniswapV3AdapterSession) Pause() (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.Pause(&_UniswapV3Adapter.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactorSession) Pause() (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.Pause(&_UniswapV3Adapter.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _UniswapV3Adapter.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_UniswapV3Adapter *UniswapV3AdapterSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.RenounceRole(&_UniswapV3Adapter.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.RenounceRole(&_UniswapV3Adapter.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _UniswapV3Adapter.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_UniswapV3Adapter *UniswapV3AdapterSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.RevokeRole(&_UniswapV3Adapter.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.RevokeRole(&_UniswapV3Adapter.TransactOpts, role, account)
}

// SetStateAggregator is a paid mutator transaction binding the contract method 0xc0804b46.
//
// Solidity: function setStateAggregator(address _stateAggregator) returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactor) SetStateAggregator(opts *bind.TransactOpts, _stateAggregator common.Address) (*types.Transaction, error) {
	return _UniswapV3Adapter.contract.Transact(opts, "setStateAggregator", _stateAggregator)
}

// SetStateAggregator is a paid mutator transaction binding the contract method 0xc0804b46.
//
// Solidity: function setStateAggregator(address _stateAggregator) returns()
func (_UniswapV3Adapter *UniswapV3AdapterSession) SetStateAggregator(_stateAggregator common.Address) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.SetStateAggregator(&_UniswapV3Adapter.TransactOpts, _stateAggregator)
}

// SetStateAggregator is a paid mutator transaction binding the contract method 0xc0804b46.
//
// Solidity: function setStateAggregator(address _stateAggregator) returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactorSession) SetStateAggregator(_stateAggregator common.Address) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.SetStateAggregator(&_UniswapV3Adapter.TransactOpts, _stateAggregator)
}

// SwapExactInput is a paid mutator transaction binding the contract method 0x1d56ec33.
//
// Solidity: function swapExactInput(bytes path, uint256 amountIn, uint256 amountOutMinimum, uint256 deadline) returns(uint256 amountOut)
func (_UniswapV3Adapter *UniswapV3AdapterTransactor) SwapExactInput(opts *bind.TransactOpts, path []byte, amountIn *big.Int, amountOutMinimum *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV3Adapter.contract.Transact(opts, "swapExactInput", path, amountIn, amountOutMinimum, deadline)
}

// SwapExactInput is a paid mutator transaction binding the contract method 0x1d56ec33.
//
// Solidity: function swapExactInput(bytes path, uint256 amountIn, uint256 amountOutMinimum, uint256 deadline) returns(uint256 amountOut)
func (_UniswapV3Adapter *UniswapV3AdapterSession) SwapExactInput(path []byte, amountIn *big.Int, amountOutMinimum *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.SwapExactInput(&_UniswapV3Adapter.TransactOpts, path, amountIn, amountOutMinimum, deadline)
}

// SwapExactInput is a paid mutator transaction binding the contract method 0x1d56ec33.
//
// Solidity: function swapExactInput(bytes path, uint256 amountIn, uint256 amountOutMinimum, uint256 deadline) returns(uint256 amountOut)
func (_UniswapV3Adapter *UniswapV3AdapterTransactorSession) SwapExactInput(path []byte, amountIn *big.Int, amountOutMinimum *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.SwapExactInput(&_UniswapV3Adapter.TransactOpts, path, amountIn, amountOutMinimum, deadline)
}

// SwapExactInputSingle is a paid mutator transaction binding the contract method 0x51e820e7.
//
// Solidity: function swapExactInputSingle(address tokenIn, address tokenOut, uint24 fee, uint256 amountIn, uint256 amountOutMinimum, uint256 deadline) returns(uint256 amountOut)
func (_UniswapV3Adapter *UniswapV3AdapterTransactor) SwapExactInputSingle(opts *bind.TransactOpts, tokenIn common.Address, tokenOut common.Address, fee *big.Int, amountIn *big.Int, amountOutMinimum *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV3Adapter.contract.Transact(opts, "swapExactInputSingle", tokenIn, tokenOut, fee, amountIn, amountOutMinimum, deadline)
}

// SwapExactInputSingle is a paid mutator transaction binding the contract method 0x51e820e7.
//
// Solidity: function swapExactInputSingle(address tokenIn, address tokenOut, uint24 fee, uint256 amountIn, uint256 amountOutMinimum, uint256 deadline) returns(uint256 amountOut)
func (_UniswapV3Adapter *UniswapV3AdapterSession) SwapExactInputSingle(tokenIn common.Address, tokenOut common.Address, fee *big.Int, amountIn *big.Int, amountOutMinimum *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.SwapExactInputSingle(&_UniswapV3Adapter.TransactOpts, tokenIn, tokenOut, fee, amountIn, amountOutMinimum, deadline)
}

// SwapExactInputSingle is a paid mutator transaction binding the contract method 0x51e820e7.
//
// Solidity: function swapExactInputSingle(address tokenIn, address tokenOut, uint24 fee, uint256 amountIn, uint256 amountOutMinimum, uint256 deadline) returns(uint256 amountOut)
func (_UniswapV3Adapter *UniswapV3AdapterTransactorSession) SwapExactInputSingle(tokenIn common.Address, tokenOut common.Address, fee *big.Int, amountIn *big.Int, amountOutMinimum *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.SwapExactInputSingle(&_UniswapV3Adapter.TransactOpts, tokenIn, tokenOut, fee, amountIn, amountOutMinimum, deadline)
}

// SwapExactOutputSingle is a paid mutator transaction binding the contract method 0x023884b7.
//
// Solidity: function swapExactOutputSingle(address tokenIn, address tokenOut, uint24 fee, uint256 amountOut, uint256 amountInMaximum, uint256 deadline) returns(uint256 amountIn)
func (_UniswapV3Adapter *UniswapV3AdapterTransactor) SwapExactOutputSingle(opts *bind.TransactOpts, tokenIn common.Address, tokenOut common.Address, fee *big.Int, amountOut *big.Int, amountInMaximum *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV3Adapter.contract.Transact(opts, "swapExactOutputSingle", tokenIn, tokenOut, fee, amountOut, amountInMaximum, deadline)
}

// SwapExactOutputSingle is a paid mutator transaction binding the contract method 0x023884b7.
//
// Solidity: function swapExactOutputSingle(address tokenIn, address tokenOut, uint24 fee, uint256 amountOut, uint256 amountInMaximum, uint256 deadline) returns(uint256 amountIn)
func (_UniswapV3Adapter *UniswapV3AdapterSession) SwapExactOutputSingle(tokenIn common.Address, tokenOut common.Address, fee *big.Int, amountOut *big.Int, amountInMaximum *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.SwapExactOutputSingle(&_UniswapV3Adapter.TransactOpts, tokenIn, tokenOut, fee, amountOut, amountInMaximum, deadline)
}

// SwapExactOutputSingle is a paid mutator transaction binding the contract method 0x023884b7.
//
// Solidity: function swapExactOutputSingle(address tokenIn, address tokenOut, uint24 fee, uint256 amountOut, uint256 amountInMaximum, uint256 deadline) returns(uint256 amountIn)
func (_UniswapV3Adapter *UniswapV3AdapterTransactorSession) SwapExactOutputSingle(tokenIn common.Address, tokenOut common.Address, fee *big.Int, amountOut *big.Int, amountInMaximum *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.SwapExactOutputSingle(&_UniswapV3Adapter.TransactOpts, tokenIn, tokenOut, fee, amountOut, amountInMaximum, deadline)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Adapter.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_UniswapV3Adapter *UniswapV3AdapterSession) Unpause() (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.Unpause(&_UniswapV3Adapter.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactorSession) Unpause() (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.Unpause(&_UniswapV3Adapter.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Adapter.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_UniswapV3Adapter *UniswapV3AdapterSession) Receive() (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.Receive(&_UniswapV3Adapter.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_UniswapV3Adapter *UniswapV3AdapterTransactorSession) Receive() (*types.Transaction, error) {
	return _UniswapV3Adapter.Contract.Receive(&_UniswapV3Adapter.TransactOpts)
}

// UniswapV3AdapterMultiHopSwapIterator is returned from FilterMultiHopSwap and is used to iterate over the raw logs and unpacked data for MultiHopSwap events raised by the UniswapV3Adapter contract.
type UniswapV3AdapterMultiHopSwapIterator struct {
	Event *UniswapV3AdapterMultiHopSwap // Event containing the contract specifics and raw log

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
func (it *UniswapV3AdapterMultiHopSwapIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3AdapterMultiHopSwap)
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
		it.Event = new(UniswapV3AdapterMultiHopSwap)
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
func (it *UniswapV3AdapterMultiHopSwapIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3AdapterMultiHopSwapIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3AdapterMultiHopSwap represents a MultiHopSwap event raised by the UniswapV3Adapter contract.
type UniswapV3AdapterMultiHopSwap struct {
	User      common.Address
	AmountIn  *big.Int
	AmountOut *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterMultiHopSwap is a free log retrieval operation binding the contract event 0x785d843068b6b9579cfe6da7d04a95c5ee9543c8adc6dcf0fd7c8bc6841686a8.
//
// Solidity: event MultiHopSwap(address indexed user, uint256 amountIn, uint256 amountOut)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) FilterMultiHopSwap(opts *bind.FilterOpts, user []common.Address) (*UniswapV3AdapterMultiHopSwapIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _UniswapV3Adapter.contract.FilterLogs(opts, "MultiHopSwap", userRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3AdapterMultiHopSwapIterator{contract: _UniswapV3Adapter.contract, event: "MultiHopSwap", logs: logs, sub: sub}, nil
}

// WatchMultiHopSwap is a free log subscription operation binding the contract event 0x785d843068b6b9579cfe6da7d04a95c5ee9543c8adc6dcf0fd7c8bc6841686a8.
//
// Solidity: event MultiHopSwap(address indexed user, uint256 amountIn, uint256 amountOut)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) WatchMultiHopSwap(opts *bind.WatchOpts, sink chan<- *UniswapV3AdapterMultiHopSwap, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _UniswapV3Adapter.contract.WatchLogs(opts, "MultiHopSwap", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3AdapterMultiHopSwap)
				if err := _UniswapV3Adapter.contract.UnpackLog(event, "MultiHopSwap", log); err != nil {
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

// ParseMultiHopSwap is a log parse operation binding the contract event 0x785d843068b6b9579cfe6da7d04a95c5ee9543c8adc6dcf0fd7c8bc6841686a8.
//
// Solidity: event MultiHopSwap(address indexed user, uint256 amountIn, uint256 amountOut)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) ParseMultiHopSwap(log types.Log) (*UniswapV3AdapterMultiHopSwap, error) {
	event := new(UniswapV3AdapterMultiHopSwap)
	if err := _UniswapV3Adapter.contract.UnpackLog(event, "MultiHopSwap", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3AdapterPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the UniswapV3Adapter contract.
type UniswapV3AdapterPausedIterator struct {
	Event *UniswapV3AdapterPaused // Event containing the contract specifics and raw log

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
func (it *UniswapV3AdapterPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3AdapterPaused)
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
		it.Event = new(UniswapV3AdapterPaused)
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
func (it *UniswapV3AdapterPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3AdapterPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3AdapterPaused represents a Paused event raised by the UniswapV3Adapter contract.
type UniswapV3AdapterPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) FilterPaused(opts *bind.FilterOpts) (*UniswapV3AdapterPausedIterator, error) {

	logs, sub, err := _UniswapV3Adapter.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &UniswapV3AdapterPausedIterator{contract: _UniswapV3Adapter.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *UniswapV3AdapterPaused) (event.Subscription, error) {

	logs, sub, err := _UniswapV3Adapter.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3AdapterPaused)
				if err := _UniswapV3Adapter.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) ParsePaused(log types.Log) (*UniswapV3AdapterPaused, error) {
	event := new(UniswapV3AdapterPaused)
	if err := _UniswapV3Adapter.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3AdapterRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the UniswapV3Adapter contract.
type UniswapV3AdapterRoleAdminChangedIterator struct {
	Event *UniswapV3AdapterRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *UniswapV3AdapterRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3AdapterRoleAdminChanged)
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
		it.Event = new(UniswapV3AdapterRoleAdminChanged)
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
func (it *UniswapV3AdapterRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3AdapterRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3AdapterRoleAdminChanged represents a RoleAdminChanged event raised by the UniswapV3Adapter contract.
type UniswapV3AdapterRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*UniswapV3AdapterRoleAdminChangedIterator, error) {

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

	logs, sub, err := _UniswapV3Adapter.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3AdapterRoleAdminChangedIterator{contract: _UniswapV3Adapter.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *UniswapV3AdapterRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _UniswapV3Adapter.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3AdapterRoleAdminChanged)
				if err := _UniswapV3Adapter.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) ParseRoleAdminChanged(log types.Log) (*UniswapV3AdapterRoleAdminChanged, error) {
	event := new(UniswapV3AdapterRoleAdminChanged)
	if err := _UniswapV3Adapter.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3AdapterRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the UniswapV3Adapter contract.
type UniswapV3AdapterRoleGrantedIterator struct {
	Event *UniswapV3AdapterRoleGranted // Event containing the contract specifics and raw log

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
func (it *UniswapV3AdapterRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3AdapterRoleGranted)
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
		it.Event = new(UniswapV3AdapterRoleGranted)
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
func (it *UniswapV3AdapterRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3AdapterRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3AdapterRoleGranted represents a RoleGranted event raised by the UniswapV3Adapter contract.
type UniswapV3AdapterRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*UniswapV3AdapterRoleGrantedIterator, error) {

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

	logs, sub, err := _UniswapV3Adapter.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3AdapterRoleGrantedIterator{contract: _UniswapV3Adapter.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *UniswapV3AdapterRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _UniswapV3Adapter.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3AdapterRoleGranted)
				if err := _UniswapV3Adapter.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) ParseRoleGranted(log types.Log) (*UniswapV3AdapterRoleGranted, error) {
	event := new(UniswapV3AdapterRoleGranted)
	if err := _UniswapV3Adapter.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3AdapterRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the UniswapV3Adapter contract.
type UniswapV3AdapterRoleRevokedIterator struct {
	Event *UniswapV3AdapterRoleRevoked // Event containing the contract specifics and raw log

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
func (it *UniswapV3AdapterRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3AdapterRoleRevoked)
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
		it.Event = new(UniswapV3AdapterRoleRevoked)
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
func (it *UniswapV3AdapterRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3AdapterRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3AdapterRoleRevoked represents a RoleRevoked event raised by the UniswapV3Adapter contract.
type UniswapV3AdapterRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*UniswapV3AdapterRoleRevokedIterator, error) {

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

	logs, sub, err := _UniswapV3Adapter.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3AdapterRoleRevokedIterator{contract: _UniswapV3Adapter.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *UniswapV3AdapterRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _UniswapV3Adapter.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3AdapterRoleRevoked)
				if err := _UniswapV3Adapter.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) ParseRoleRevoked(log types.Log) (*UniswapV3AdapterRoleRevoked, error) {
	event := new(UniswapV3AdapterRoleRevoked)
	if err := _UniswapV3Adapter.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3AdapterSwappedIterator is returned from FilterSwapped and is used to iterate over the raw logs and unpacked data for Swapped events raised by the UniswapV3Adapter contract.
type UniswapV3AdapterSwappedIterator struct {
	Event *UniswapV3AdapterSwapped // Event containing the contract specifics and raw log

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
func (it *UniswapV3AdapterSwappedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3AdapterSwapped)
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
		it.Event = new(UniswapV3AdapterSwapped)
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
func (it *UniswapV3AdapterSwappedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3AdapterSwappedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3AdapterSwapped represents a Swapped event raised by the UniswapV3Adapter contract.
type UniswapV3AdapterSwapped struct {
	User      common.Address
	TokenIn   common.Address
	TokenOut  common.Address
	AmountIn  *big.Int
	AmountOut *big.Int
	Fee       *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSwapped is a free log retrieval operation binding the contract event 0x62a422cc36234492f287aeb712657b805d79c0ae49f3c5f2d23570342854656a.
//
// Solidity: event Swapped(address indexed user, address indexed tokenIn, address indexed tokenOut, uint256 amountIn, uint256 amountOut, uint24 fee)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) FilterSwapped(opts *bind.FilterOpts, user []common.Address, tokenIn []common.Address, tokenOut []common.Address) (*UniswapV3AdapterSwappedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenInRule []interface{}
	for _, tokenInItem := range tokenIn {
		tokenInRule = append(tokenInRule, tokenInItem)
	}
	var tokenOutRule []interface{}
	for _, tokenOutItem := range tokenOut {
		tokenOutRule = append(tokenOutRule, tokenOutItem)
	}

	logs, sub, err := _UniswapV3Adapter.contract.FilterLogs(opts, "Swapped", userRule, tokenInRule, tokenOutRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3AdapterSwappedIterator{contract: _UniswapV3Adapter.contract, event: "Swapped", logs: logs, sub: sub}, nil
}

// WatchSwapped is a free log subscription operation binding the contract event 0x62a422cc36234492f287aeb712657b805d79c0ae49f3c5f2d23570342854656a.
//
// Solidity: event Swapped(address indexed user, address indexed tokenIn, address indexed tokenOut, uint256 amountIn, uint256 amountOut, uint24 fee)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) WatchSwapped(opts *bind.WatchOpts, sink chan<- *UniswapV3AdapterSwapped, user []common.Address, tokenIn []common.Address, tokenOut []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenInRule []interface{}
	for _, tokenInItem := range tokenIn {
		tokenInRule = append(tokenInRule, tokenInItem)
	}
	var tokenOutRule []interface{}
	for _, tokenOutItem := range tokenOut {
		tokenOutRule = append(tokenOutRule, tokenOutItem)
	}

	logs, sub, err := _UniswapV3Adapter.contract.WatchLogs(opts, "Swapped", userRule, tokenInRule, tokenOutRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3AdapterSwapped)
				if err := _UniswapV3Adapter.contract.UnpackLog(event, "Swapped", log); err != nil {
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

// ParseSwapped is a log parse operation binding the contract event 0x62a422cc36234492f287aeb712657b805d79c0ae49f3c5f2d23570342854656a.
//
// Solidity: event Swapped(address indexed user, address indexed tokenIn, address indexed tokenOut, uint256 amountIn, uint256 amountOut, uint24 fee)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) ParseSwapped(log types.Log) (*UniswapV3AdapterSwapped, error) {
	event := new(UniswapV3AdapterSwapped)
	if err := _UniswapV3Adapter.contract.UnpackLog(event, "Swapped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3AdapterUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the UniswapV3Adapter contract.
type UniswapV3AdapterUnpausedIterator struct {
	Event *UniswapV3AdapterUnpaused // Event containing the contract specifics and raw log

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
func (it *UniswapV3AdapterUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3AdapterUnpaused)
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
		it.Event = new(UniswapV3AdapterUnpaused)
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
func (it *UniswapV3AdapterUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3AdapterUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3AdapterUnpaused represents a Unpaused event raised by the UniswapV3Adapter contract.
type UniswapV3AdapterUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) FilterUnpaused(opts *bind.FilterOpts) (*UniswapV3AdapterUnpausedIterator, error) {

	logs, sub, err := _UniswapV3Adapter.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &UniswapV3AdapterUnpausedIterator{contract: _UniswapV3Adapter.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *UniswapV3AdapterUnpaused) (event.Subscription, error) {

	logs, sub, err := _UniswapV3Adapter.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3AdapterUnpaused)
				if err := _UniswapV3Adapter.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_UniswapV3Adapter *UniswapV3AdapterFilterer) ParseUnpaused(log types.Log) (*UniswapV3AdapterUnpaused, error) {
	event := new(UniswapV3AdapterUnpaused)
	if err := _UniswapV3Adapter.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
