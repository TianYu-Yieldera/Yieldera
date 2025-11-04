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

// IntegratedVaultMetaData contains all meta data concerning the IntegratedVault contract.
var IntegratedVaultMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_collateralToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_lusdToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_priceOracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_stateAggregator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Borrowed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"interestAmount\",\"type\":\"uint256\"}],\"name\":\"InterestAccrued\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"liquidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"debtRepaid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"collateralSeized\",\"type\":\"uint256\"}],\"name\":\"Liquidated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Repaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COLLATERAL_RATIO\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INTEREST_RATE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LIQUIDATION_BONUS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LIQUIDATION_THRESHOLD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PRICE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SECONDS_PER_YEAR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"accrueInterest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"interestAccrued\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activePositions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"borrow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collateralToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getAccruedInterest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"interestAccrued\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getCollateralValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"valueUSD\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getMaxBorrowAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"maxBorrow\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getMaxWithdrawAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"maxWithdraw\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getPosition\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"collateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"debt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"healthFactor\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getTotalDebt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalDebtWithInterest\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserHealthFactor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"healthFactor\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"isLiquidatable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"canLiquidate\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastInterestUpdate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"debtToCover\",\"type\":\"uint256\"}],\"name\":\"liquidate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lusdToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"positions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"collateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"debt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"priceOracle\",\"outputs\":[{\"internalType\":\"contractIPriceOracle\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"repay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateAggregator\",\"outputs\":[{\"internalType\":\"contractL2StateAggregator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalDebt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOracle\",\"type\":\"address\"}],\"name\":\"updatePriceOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newAggregator\",\"type\":\"address\"}],\"name\":\"updateStateAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IntegratedVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use IntegratedVaultMetaData.ABI instead.
var IntegratedVaultABI = IntegratedVaultMetaData.ABI

// IntegratedVault is an auto generated Go binding around an Ethereum contract.
type IntegratedVault struct {
	IntegratedVaultCaller     // Read-only binding to the contract
	IntegratedVaultTransactor // Write-only binding to the contract
	IntegratedVaultFilterer   // Log filterer for contract events
}

// IntegratedVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type IntegratedVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IntegratedVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IntegratedVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IntegratedVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IntegratedVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IntegratedVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IntegratedVaultSession struct {
	Contract     *IntegratedVault  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IntegratedVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IntegratedVaultCallerSession struct {
	Contract *IntegratedVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IntegratedVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IntegratedVaultTransactorSession struct {
	Contract     *IntegratedVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IntegratedVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type IntegratedVaultRaw struct {
	Contract *IntegratedVault // Generic contract binding to access the raw methods on
}

// IntegratedVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IntegratedVaultCallerRaw struct {
	Contract *IntegratedVaultCaller // Generic read-only contract binding to access the raw methods on
}

// IntegratedVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IntegratedVaultTransactorRaw struct {
	Contract *IntegratedVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIntegratedVault creates a new instance of IntegratedVault, bound to a specific deployed contract.
func NewIntegratedVault(address common.Address, backend bind.ContractBackend) (*IntegratedVault, error) {
	contract, err := bindIntegratedVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IntegratedVault{IntegratedVaultCaller: IntegratedVaultCaller{contract: contract}, IntegratedVaultTransactor: IntegratedVaultTransactor{contract: contract}, IntegratedVaultFilterer: IntegratedVaultFilterer{contract: contract}}, nil
}

// NewIntegratedVaultCaller creates a new read-only instance of IntegratedVault, bound to a specific deployed contract.
func NewIntegratedVaultCaller(address common.Address, caller bind.ContractCaller) (*IntegratedVaultCaller, error) {
	contract, err := bindIntegratedVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IntegratedVaultCaller{contract: contract}, nil
}

// NewIntegratedVaultTransactor creates a new write-only instance of IntegratedVault, bound to a specific deployed contract.
func NewIntegratedVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*IntegratedVaultTransactor, error) {
	contract, err := bindIntegratedVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IntegratedVaultTransactor{contract: contract}, nil
}

// NewIntegratedVaultFilterer creates a new log filterer instance of IntegratedVault, bound to a specific deployed contract.
func NewIntegratedVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*IntegratedVaultFilterer, error) {
	contract, err := bindIntegratedVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IntegratedVaultFilterer{contract: contract}, nil
}

// bindIntegratedVault binds a generic wrapper to an already deployed contract.
func bindIntegratedVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IntegratedVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IntegratedVault *IntegratedVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IntegratedVault.Contract.IntegratedVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IntegratedVault *IntegratedVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IntegratedVault.Contract.IntegratedVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IntegratedVault *IntegratedVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IntegratedVault.Contract.IntegratedVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IntegratedVault *IntegratedVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IntegratedVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IntegratedVault *IntegratedVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IntegratedVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IntegratedVault *IntegratedVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IntegratedVault.Contract.contract.Transact(opts, method, params...)
}

// COLLATERALRATIO is a free data retrieval call binding the contract method 0xd9e69a05.
//
// Solidity: function COLLATERAL_RATIO() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCaller) COLLATERALRATIO(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "COLLATERAL_RATIO")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// COLLATERALRATIO is a free data retrieval call binding the contract method 0xd9e69a05.
//
// Solidity: function COLLATERAL_RATIO() view returns(uint256)
func (_IntegratedVault *IntegratedVaultSession) COLLATERALRATIO() (*big.Int, error) {
	return _IntegratedVault.Contract.COLLATERALRATIO(&_IntegratedVault.CallOpts)
}

// COLLATERALRATIO is a free data retrieval call binding the contract method 0xd9e69a05.
//
// Solidity: function COLLATERAL_RATIO() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCallerSession) COLLATERALRATIO() (*big.Int, error) {
	return _IntegratedVault.Contract.COLLATERALRATIO(&_IntegratedVault.CallOpts)
}

// INTERESTRATE is a free data retrieval call binding the contract method 0x5b72a33a.
//
// Solidity: function INTEREST_RATE() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCaller) INTERESTRATE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "INTEREST_RATE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INTERESTRATE is a free data retrieval call binding the contract method 0x5b72a33a.
//
// Solidity: function INTEREST_RATE() view returns(uint256)
func (_IntegratedVault *IntegratedVaultSession) INTERESTRATE() (*big.Int, error) {
	return _IntegratedVault.Contract.INTERESTRATE(&_IntegratedVault.CallOpts)
}

// INTERESTRATE is a free data retrieval call binding the contract method 0x5b72a33a.
//
// Solidity: function INTEREST_RATE() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCallerSession) INTERESTRATE() (*big.Int, error) {
	return _IntegratedVault.Contract.INTERESTRATE(&_IntegratedVault.CallOpts)
}

// LIQUIDATIONBONUS is a free data retrieval call binding the contract method 0x3574d4c4.
//
// Solidity: function LIQUIDATION_BONUS() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCaller) LIQUIDATIONBONUS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "LIQUIDATION_BONUS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LIQUIDATIONBONUS is a free data retrieval call binding the contract method 0x3574d4c4.
//
// Solidity: function LIQUIDATION_BONUS() view returns(uint256)
func (_IntegratedVault *IntegratedVaultSession) LIQUIDATIONBONUS() (*big.Int, error) {
	return _IntegratedVault.Contract.LIQUIDATIONBONUS(&_IntegratedVault.CallOpts)
}

// LIQUIDATIONBONUS is a free data retrieval call binding the contract method 0x3574d4c4.
//
// Solidity: function LIQUIDATION_BONUS() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCallerSession) LIQUIDATIONBONUS() (*big.Int, error) {
	return _IntegratedVault.Contract.LIQUIDATIONBONUS(&_IntegratedVault.CallOpts)
}

// LIQUIDATIONTHRESHOLD is a free data retrieval call binding the contract method 0x90a8ae9b.
//
// Solidity: function LIQUIDATION_THRESHOLD() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCaller) LIQUIDATIONTHRESHOLD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "LIQUIDATION_THRESHOLD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LIQUIDATIONTHRESHOLD is a free data retrieval call binding the contract method 0x90a8ae9b.
//
// Solidity: function LIQUIDATION_THRESHOLD() view returns(uint256)
func (_IntegratedVault *IntegratedVaultSession) LIQUIDATIONTHRESHOLD() (*big.Int, error) {
	return _IntegratedVault.Contract.LIQUIDATIONTHRESHOLD(&_IntegratedVault.CallOpts)
}

// LIQUIDATIONTHRESHOLD is a free data retrieval call binding the contract method 0x90a8ae9b.
//
// Solidity: function LIQUIDATION_THRESHOLD() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCallerSession) LIQUIDATIONTHRESHOLD() (*big.Int, error) {
	return _IntegratedVault.Contract.LIQUIDATIONTHRESHOLD(&_IntegratedVault.CallOpts)
}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCaller) PRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_IntegratedVault *IntegratedVaultSession) PRECISION() (*big.Int, error) {
	return _IntegratedVault.Contract.PRECISION(&_IntegratedVault.CallOpts)
}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCallerSession) PRECISION() (*big.Int, error) {
	return _IntegratedVault.Contract.PRECISION(&_IntegratedVault.CallOpts)
}

// PRICEPRECISION is a free data retrieval call binding the contract method 0x95082d25.
//
// Solidity: function PRICE_PRECISION() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCaller) PRICEPRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "PRICE_PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PRICEPRECISION is a free data retrieval call binding the contract method 0x95082d25.
//
// Solidity: function PRICE_PRECISION() view returns(uint256)
func (_IntegratedVault *IntegratedVaultSession) PRICEPRECISION() (*big.Int, error) {
	return _IntegratedVault.Contract.PRICEPRECISION(&_IntegratedVault.CallOpts)
}

// PRICEPRECISION is a free data retrieval call binding the contract method 0x95082d25.
//
// Solidity: function PRICE_PRECISION() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCallerSession) PRICEPRECISION() (*big.Int, error) {
	return _IntegratedVault.Contract.PRICEPRECISION(&_IntegratedVault.CallOpts)
}

// SECONDSPERYEAR is a free data retrieval call binding the contract method 0xe6a69ab8.
//
// Solidity: function SECONDS_PER_YEAR() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCaller) SECONDSPERYEAR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "SECONDS_PER_YEAR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SECONDSPERYEAR is a free data retrieval call binding the contract method 0xe6a69ab8.
//
// Solidity: function SECONDS_PER_YEAR() view returns(uint256)
func (_IntegratedVault *IntegratedVaultSession) SECONDSPERYEAR() (*big.Int, error) {
	return _IntegratedVault.Contract.SECONDSPERYEAR(&_IntegratedVault.CallOpts)
}

// SECONDSPERYEAR is a free data retrieval call binding the contract method 0xe6a69ab8.
//
// Solidity: function SECONDS_PER_YEAR() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCallerSession) SECONDSPERYEAR() (*big.Int, error) {
	return _IntegratedVault.Contract.SECONDSPERYEAR(&_IntegratedVault.CallOpts)
}

// ActivePositions is a free data retrieval call binding the contract method 0x297add9d.
//
// Solidity: function activePositions() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCaller) ActivePositions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "activePositions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActivePositions is a free data retrieval call binding the contract method 0x297add9d.
//
// Solidity: function activePositions() view returns(uint256)
func (_IntegratedVault *IntegratedVaultSession) ActivePositions() (*big.Int, error) {
	return _IntegratedVault.Contract.ActivePositions(&_IntegratedVault.CallOpts)
}

// ActivePositions is a free data retrieval call binding the contract method 0x297add9d.
//
// Solidity: function activePositions() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCallerSession) ActivePositions() (*big.Int, error) {
	return _IntegratedVault.Contract.ActivePositions(&_IntegratedVault.CallOpts)
}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_IntegratedVault *IntegratedVaultCaller) CollateralToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "collateralToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_IntegratedVault *IntegratedVaultSession) CollateralToken() (common.Address, error) {
	return _IntegratedVault.Contract.CollateralToken(&_IntegratedVault.CallOpts)
}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_IntegratedVault *IntegratedVaultCallerSession) CollateralToken() (common.Address, error) {
	return _IntegratedVault.Contract.CollateralToken(&_IntegratedVault.CallOpts)
}

// GetAccruedInterest is a free data retrieval call binding the contract method 0xbbf14e74.
//
// Solidity: function getAccruedInterest(address user) view returns(uint256 interestAccrued)
func (_IntegratedVault *IntegratedVaultCaller) GetAccruedInterest(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "getAccruedInterest", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAccruedInterest is a free data retrieval call binding the contract method 0xbbf14e74.
//
// Solidity: function getAccruedInterest(address user) view returns(uint256 interestAccrued)
func (_IntegratedVault *IntegratedVaultSession) GetAccruedInterest(user common.Address) (*big.Int, error) {
	return _IntegratedVault.Contract.GetAccruedInterest(&_IntegratedVault.CallOpts, user)
}

// GetAccruedInterest is a free data retrieval call binding the contract method 0xbbf14e74.
//
// Solidity: function getAccruedInterest(address user) view returns(uint256 interestAccrued)
func (_IntegratedVault *IntegratedVaultCallerSession) GetAccruedInterest(user common.Address) (*big.Int, error) {
	return _IntegratedVault.Contract.GetAccruedInterest(&_IntegratedVault.CallOpts, user)
}

// GetCollateralValue is a free data retrieval call binding the contract method 0x97904e42.
//
// Solidity: function getCollateralValue(address user) view returns(uint256 valueUSD)
func (_IntegratedVault *IntegratedVaultCaller) GetCollateralValue(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "getCollateralValue", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCollateralValue is a free data retrieval call binding the contract method 0x97904e42.
//
// Solidity: function getCollateralValue(address user) view returns(uint256 valueUSD)
func (_IntegratedVault *IntegratedVaultSession) GetCollateralValue(user common.Address) (*big.Int, error) {
	return _IntegratedVault.Contract.GetCollateralValue(&_IntegratedVault.CallOpts, user)
}

// GetCollateralValue is a free data retrieval call binding the contract method 0x97904e42.
//
// Solidity: function getCollateralValue(address user) view returns(uint256 valueUSD)
func (_IntegratedVault *IntegratedVaultCallerSession) GetCollateralValue(user common.Address) (*big.Int, error) {
	return _IntegratedVault.Contract.GetCollateralValue(&_IntegratedVault.CallOpts, user)
}

// GetMaxBorrowAmount is a free data retrieval call binding the contract method 0x4b79cb0f.
//
// Solidity: function getMaxBorrowAmount(address user) view returns(uint256 maxBorrow)
func (_IntegratedVault *IntegratedVaultCaller) GetMaxBorrowAmount(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "getMaxBorrowAmount", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxBorrowAmount is a free data retrieval call binding the contract method 0x4b79cb0f.
//
// Solidity: function getMaxBorrowAmount(address user) view returns(uint256 maxBorrow)
func (_IntegratedVault *IntegratedVaultSession) GetMaxBorrowAmount(user common.Address) (*big.Int, error) {
	return _IntegratedVault.Contract.GetMaxBorrowAmount(&_IntegratedVault.CallOpts, user)
}

// GetMaxBorrowAmount is a free data retrieval call binding the contract method 0x4b79cb0f.
//
// Solidity: function getMaxBorrowAmount(address user) view returns(uint256 maxBorrow)
func (_IntegratedVault *IntegratedVaultCallerSession) GetMaxBorrowAmount(user common.Address) (*big.Int, error) {
	return _IntegratedVault.Contract.GetMaxBorrowAmount(&_IntegratedVault.CallOpts, user)
}

// GetMaxWithdrawAmount is a free data retrieval call binding the contract method 0x21ff22d6.
//
// Solidity: function getMaxWithdrawAmount(address user) view returns(uint256 maxWithdraw)
func (_IntegratedVault *IntegratedVaultCaller) GetMaxWithdrawAmount(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "getMaxWithdrawAmount", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxWithdrawAmount is a free data retrieval call binding the contract method 0x21ff22d6.
//
// Solidity: function getMaxWithdrawAmount(address user) view returns(uint256 maxWithdraw)
func (_IntegratedVault *IntegratedVaultSession) GetMaxWithdrawAmount(user common.Address) (*big.Int, error) {
	return _IntegratedVault.Contract.GetMaxWithdrawAmount(&_IntegratedVault.CallOpts, user)
}

// GetMaxWithdrawAmount is a free data retrieval call binding the contract method 0x21ff22d6.
//
// Solidity: function getMaxWithdrawAmount(address user) view returns(uint256 maxWithdraw)
func (_IntegratedVault *IntegratedVaultCallerSession) GetMaxWithdrawAmount(user common.Address) (*big.Int, error) {
	return _IntegratedVault.Contract.GetMaxWithdrawAmount(&_IntegratedVault.CallOpts, user)
}

// GetPosition is a free data retrieval call binding the contract method 0x16c19739.
//
// Solidity: function getPosition(address user) view returns(uint256 collateral, uint256 debt, uint256 healthFactor)
func (_IntegratedVault *IntegratedVaultCaller) GetPosition(opts *bind.CallOpts, user common.Address) (struct {
	Collateral   *big.Int
	Debt         *big.Int
	HealthFactor *big.Int
}, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "getPosition", user)

	outstruct := new(struct {
		Collateral   *big.Int
		Debt         *big.Int
		HealthFactor *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Collateral = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Debt = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.HealthFactor = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetPosition is a free data retrieval call binding the contract method 0x16c19739.
//
// Solidity: function getPosition(address user) view returns(uint256 collateral, uint256 debt, uint256 healthFactor)
func (_IntegratedVault *IntegratedVaultSession) GetPosition(user common.Address) (struct {
	Collateral   *big.Int
	Debt         *big.Int
	HealthFactor *big.Int
}, error) {
	return _IntegratedVault.Contract.GetPosition(&_IntegratedVault.CallOpts, user)
}

// GetPosition is a free data retrieval call binding the contract method 0x16c19739.
//
// Solidity: function getPosition(address user) view returns(uint256 collateral, uint256 debt, uint256 healthFactor)
func (_IntegratedVault *IntegratedVaultCallerSession) GetPosition(user common.Address) (struct {
	Collateral   *big.Int
	Debt         *big.Int
	HealthFactor *big.Int
}, error) {
	return _IntegratedVault.Contract.GetPosition(&_IntegratedVault.CallOpts, user)
}

// GetTotalDebt is a free data retrieval call binding the contract method 0x4d44ac4f.
//
// Solidity: function getTotalDebt(address user) view returns(uint256 totalDebtWithInterest)
func (_IntegratedVault *IntegratedVaultCaller) GetTotalDebt(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "getTotalDebt", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalDebt is a free data retrieval call binding the contract method 0x4d44ac4f.
//
// Solidity: function getTotalDebt(address user) view returns(uint256 totalDebtWithInterest)
func (_IntegratedVault *IntegratedVaultSession) GetTotalDebt(user common.Address) (*big.Int, error) {
	return _IntegratedVault.Contract.GetTotalDebt(&_IntegratedVault.CallOpts, user)
}

// GetTotalDebt is a free data retrieval call binding the contract method 0x4d44ac4f.
//
// Solidity: function getTotalDebt(address user) view returns(uint256 totalDebtWithInterest)
func (_IntegratedVault *IntegratedVaultCallerSession) GetTotalDebt(user common.Address) (*big.Int, error) {
	return _IntegratedVault.Contract.GetTotalDebt(&_IntegratedVault.CallOpts, user)
}

// GetUserHealthFactor is a free data retrieval call binding the contract method 0x71cbfc98.
//
// Solidity: function getUserHealthFactor(address user) view returns(uint256 healthFactor)
func (_IntegratedVault *IntegratedVaultCaller) GetUserHealthFactor(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "getUserHealthFactor", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserHealthFactor is a free data retrieval call binding the contract method 0x71cbfc98.
//
// Solidity: function getUserHealthFactor(address user) view returns(uint256 healthFactor)
func (_IntegratedVault *IntegratedVaultSession) GetUserHealthFactor(user common.Address) (*big.Int, error) {
	return _IntegratedVault.Contract.GetUserHealthFactor(&_IntegratedVault.CallOpts, user)
}

// GetUserHealthFactor is a free data retrieval call binding the contract method 0x71cbfc98.
//
// Solidity: function getUserHealthFactor(address user) view returns(uint256 healthFactor)
func (_IntegratedVault *IntegratedVaultCallerSession) GetUserHealthFactor(user common.Address) (*big.Int, error) {
	return _IntegratedVault.Contract.GetUserHealthFactor(&_IntegratedVault.CallOpts, user)
}

// IsLiquidatable is a free data retrieval call binding the contract method 0x042e02cf.
//
// Solidity: function isLiquidatable(address user) view returns(bool canLiquidate)
func (_IntegratedVault *IntegratedVaultCaller) IsLiquidatable(opts *bind.CallOpts, user common.Address) (bool, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "isLiquidatable", user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsLiquidatable is a free data retrieval call binding the contract method 0x042e02cf.
//
// Solidity: function isLiquidatable(address user) view returns(bool canLiquidate)
func (_IntegratedVault *IntegratedVaultSession) IsLiquidatable(user common.Address) (bool, error) {
	return _IntegratedVault.Contract.IsLiquidatable(&_IntegratedVault.CallOpts, user)
}

// IsLiquidatable is a free data retrieval call binding the contract method 0x042e02cf.
//
// Solidity: function isLiquidatable(address user) view returns(bool canLiquidate)
func (_IntegratedVault *IntegratedVaultCallerSession) IsLiquidatable(user common.Address) (bool, error) {
	return _IntegratedVault.Contract.IsLiquidatable(&_IntegratedVault.CallOpts, user)
}

// LastInterestUpdate is a free data retrieval call binding the contract method 0x74b52142.
//
// Solidity: function lastInterestUpdate() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCaller) LastInterestUpdate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "lastInterestUpdate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastInterestUpdate is a free data retrieval call binding the contract method 0x74b52142.
//
// Solidity: function lastInterestUpdate() view returns(uint256)
func (_IntegratedVault *IntegratedVaultSession) LastInterestUpdate() (*big.Int, error) {
	return _IntegratedVault.Contract.LastInterestUpdate(&_IntegratedVault.CallOpts)
}

// LastInterestUpdate is a free data retrieval call binding the contract method 0x74b52142.
//
// Solidity: function lastInterestUpdate() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCallerSession) LastInterestUpdate() (*big.Int, error) {
	return _IntegratedVault.Contract.LastInterestUpdate(&_IntegratedVault.CallOpts)
}

// LusdToken is a free data retrieval call binding the contract method 0xb83f91a2.
//
// Solidity: function lusdToken() view returns(address)
func (_IntegratedVault *IntegratedVaultCaller) LusdToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "lusdToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LusdToken is a free data retrieval call binding the contract method 0xb83f91a2.
//
// Solidity: function lusdToken() view returns(address)
func (_IntegratedVault *IntegratedVaultSession) LusdToken() (common.Address, error) {
	return _IntegratedVault.Contract.LusdToken(&_IntegratedVault.CallOpts)
}

// LusdToken is a free data retrieval call binding the contract method 0xb83f91a2.
//
// Solidity: function lusdToken() view returns(address)
func (_IntegratedVault *IntegratedVaultCallerSession) LusdToken() (common.Address, error) {
	return _IntegratedVault.Contract.LusdToken(&_IntegratedVault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IntegratedVault *IntegratedVaultCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IntegratedVault *IntegratedVaultSession) Owner() (common.Address, error) {
	return _IntegratedVault.Contract.Owner(&_IntegratedVault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IntegratedVault *IntegratedVaultCallerSession) Owner() (common.Address, error) {
	return _IntegratedVault.Contract.Owner(&_IntegratedVault.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_IntegratedVault *IntegratedVaultCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_IntegratedVault *IntegratedVaultSession) Paused() (bool, error) {
	return _IntegratedVault.Contract.Paused(&_IntegratedVault.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_IntegratedVault *IntegratedVaultCallerSession) Paused() (bool, error) {
	return _IntegratedVault.Contract.Paused(&_IntegratedVault.CallOpts)
}

// Positions is a free data retrieval call binding the contract method 0x55f57510.
//
// Solidity: function positions(address ) view returns(uint256 collateral, uint256 debt, uint256 lastUpdate)
func (_IntegratedVault *IntegratedVaultCaller) Positions(opts *bind.CallOpts, arg0 common.Address) (struct {
	Collateral *big.Int
	Debt       *big.Int
	LastUpdate *big.Int
}, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "positions", arg0)

	outstruct := new(struct {
		Collateral *big.Int
		Debt       *big.Int
		LastUpdate *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Collateral = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Debt = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.LastUpdate = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Positions is a free data retrieval call binding the contract method 0x55f57510.
//
// Solidity: function positions(address ) view returns(uint256 collateral, uint256 debt, uint256 lastUpdate)
func (_IntegratedVault *IntegratedVaultSession) Positions(arg0 common.Address) (struct {
	Collateral *big.Int
	Debt       *big.Int
	LastUpdate *big.Int
}, error) {
	return _IntegratedVault.Contract.Positions(&_IntegratedVault.CallOpts, arg0)
}

// Positions is a free data retrieval call binding the contract method 0x55f57510.
//
// Solidity: function positions(address ) view returns(uint256 collateral, uint256 debt, uint256 lastUpdate)
func (_IntegratedVault *IntegratedVaultCallerSession) Positions(arg0 common.Address) (struct {
	Collateral *big.Int
	Debt       *big.Int
	LastUpdate *big.Int
}, error) {
	return _IntegratedVault.Contract.Positions(&_IntegratedVault.CallOpts, arg0)
}

// PriceOracle is a free data retrieval call binding the contract method 0x2630c12f.
//
// Solidity: function priceOracle() view returns(address)
func (_IntegratedVault *IntegratedVaultCaller) PriceOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "priceOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PriceOracle is a free data retrieval call binding the contract method 0x2630c12f.
//
// Solidity: function priceOracle() view returns(address)
func (_IntegratedVault *IntegratedVaultSession) PriceOracle() (common.Address, error) {
	return _IntegratedVault.Contract.PriceOracle(&_IntegratedVault.CallOpts)
}

// PriceOracle is a free data retrieval call binding the contract method 0x2630c12f.
//
// Solidity: function priceOracle() view returns(address)
func (_IntegratedVault *IntegratedVaultCallerSession) PriceOracle() (common.Address, error) {
	return _IntegratedVault.Contract.PriceOracle(&_IntegratedVault.CallOpts)
}

// StateAggregator is a free data retrieval call binding the contract method 0x54ef2920.
//
// Solidity: function stateAggregator() view returns(address)
func (_IntegratedVault *IntegratedVaultCaller) StateAggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "stateAggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StateAggregator is a free data retrieval call binding the contract method 0x54ef2920.
//
// Solidity: function stateAggregator() view returns(address)
func (_IntegratedVault *IntegratedVaultSession) StateAggregator() (common.Address, error) {
	return _IntegratedVault.Contract.StateAggregator(&_IntegratedVault.CallOpts)
}

// StateAggregator is a free data retrieval call binding the contract method 0x54ef2920.
//
// Solidity: function stateAggregator() view returns(address)
func (_IntegratedVault *IntegratedVaultCallerSession) StateAggregator() (common.Address, error) {
	return _IntegratedVault.Contract.StateAggregator(&_IntegratedVault.CallOpts)
}

// TotalCollateral is a free data retrieval call binding the contract method 0x4ac8eb5f.
//
// Solidity: function totalCollateral() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCaller) TotalCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "totalCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalCollateral is a free data retrieval call binding the contract method 0x4ac8eb5f.
//
// Solidity: function totalCollateral() view returns(uint256)
func (_IntegratedVault *IntegratedVaultSession) TotalCollateral() (*big.Int, error) {
	return _IntegratedVault.Contract.TotalCollateral(&_IntegratedVault.CallOpts)
}

// TotalCollateral is a free data retrieval call binding the contract method 0x4ac8eb5f.
//
// Solidity: function totalCollateral() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCallerSession) TotalCollateral() (*big.Int, error) {
	return _IntegratedVault.Contract.TotalCollateral(&_IntegratedVault.CallOpts)
}

// TotalDebt is a free data retrieval call binding the contract method 0xfc7b9c18.
//
// Solidity: function totalDebt() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCaller) TotalDebt(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IntegratedVault.contract.Call(opts, &out, "totalDebt")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalDebt is a free data retrieval call binding the contract method 0xfc7b9c18.
//
// Solidity: function totalDebt() view returns(uint256)
func (_IntegratedVault *IntegratedVaultSession) TotalDebt() (*big.Int, error) {
	return _IntegratedVault.Contract.TotalDebt(&_IntegratedVault.CallOpts)
}

// TotalDebt is a free data retrieval call binding the contract method 0xfc7b9c18.
//
// Solidity: function totalDebt() view returns(uint256)
func (_IntegratedVault *IntegratedVaultCallerSession) TotalDebt() (*big.Int, error) {
	return _IntegratedVault.Contract.TotalDebt(&_IntegratedVault.CallOpts)
}

// AccrueInterest is a paid mutator transaction binding the contract method 0x9198e515.
//
// Solidity: function accrueInterest(address user) returns(uint256 interestAccrued)
func (_IntegratedVault *IntegratedVaultTransactor) AccrueInterest(opts *bind.TransactOpts, user common.Address) (*types.Transaction, error) {
	return _IntegratedVault.contract.Transact(opts, "accrueInterest", user)
}

// AccrueInterest is a paid mutator transaction binding the contract method 0x9198e515.
//
// Solidity: function accrueInterest(address user) returns(uint256 interestAccrued)
func (_IntegratedVault *IntegratedVaultSession) AccrueInterest(user common.Address) (*types.Transaction, error) {
	return _IntegratedVault.Contract.AccrueInterest(&_IntegratedVault.TransactOpts, user)
}

// AccrueInterest is a paid mutator transaction binding the contract method 0x9198e515.
//
// Solidity: function accrueInterest(address user) returns(uint256 interestAccrued)
func (_IntegratedVault *IntegratedVaultTransactorSession) AccrueInterest(user common.Address) (*types.Transaction, error) {
	return _IntegratedVault.Contract.AccrueInterest(&_IntegratedVault.TransactOpts, user)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_IntegratedVault *IntegratedVaultTransactor) Borrow(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.contract.Transact(opts, "borrow", amount)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_IntegratedVault *IntegratedVaultSession) Borrow(amount *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.Contract.Borrow(&_IntegratedVault.TransactOpts, amount)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_IntegratedVault *IntegratedVaultTransactorSession) Borrow(amount *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.Contract.Borrow(&_IntegratedVault.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_IntegratedVault *IntegratedVaultTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.contract.Transact(opts, "deposit", amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_IntegratedVault *IntegratedVaultSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.Contract.Deposit(&_IntegratedVault.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_IntegratedVault *IntegratedVaultTransactorSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.Contract.Deposit(&_IntegratedVault.TransactOpts, amount)
}

// Liquidate is a paid mutator transaction binding the contract method 0xbcbaf487.
//
// Solidity: function liquidate(address user, uint256 debtToCover) returns()
func (_IntegratedVault *IntegratedVaultTransactor) Liquidate(opts *bind.TransactOpts, user common.Address, debtToCover *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.contract.Transact(opts, "liquidate", user, debtToCover)
}

// Liquidate is a paid mutator transaction binding the contract method 0xbcbaf487.
//
// Solidity: function liquidate(address user, uint256 debtToCover) returns()
func (_IntegratedVault *IntegratedVaultSession) Liquidate(user common.Address, debtToCover *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.Contract.Liquidate(&_IntegratedVault.TransactOpts, user, debtToCover)
}

// Liquidate is a paid mutator transaction binding the contract method 0xbcbaf487.
//
// Solidity: function liquidate(address user, uint256 debtToCover) returns()
func (_IntegratedVault *IntegratedVaultTransactorSession) Liquidate(user common.Address, debtToCover *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.Contract.Liquidate(&_IntegratedVault.TransactOpts, user, debtToCover)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_IntegratedVault *IntegratedVaultTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IntegratedVault.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_IntegratedVault *IntegratedVaultSession) Pause() (*types.Transaction, error) {
	return _IntegratedVault.Contract.Pause(&_IntegratedVault.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_IntegratedVault *IntegratedVaultTransactorSession) Pause() (*types.Transaction, error) {
	return _IntegratedVault.Contract.Pause(&_IntegratedVault.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IntegratedVault *IntegratedVaultTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IntegratedVault.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IntegratedVault *IntegratedVaultSession) RenounceOwnership() (*types.Transaction, error) {
	return _IntegratedVault.Contract.RenounceOwnership(&_IntegratedVault.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IntegratedVault *IntegratedVaultTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _IntegratedVault.Contract.RenounceOwnership(&_IntegratedVault.TransactOpts)
}

// Repay is a paid mutator transaction binding the contract method 0x371fd8e6.
//
// Solidity: function repay(uint256 amount) returns()
func (_IntegratedVault *IntegratedVaultTransactor) Repay(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.contract.Transact(opts, "repay", amount)
}

// Repay is a paid mutator transaction binding the contract method 0x371fd8e6.
//
// Solidity: function repay(uint256 amount) returns()
func (_IntegratedVault *IntegratedVaultSession) Repay(amount *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.Contract.Repay(&_IntegratedVault.TransactOpts, amount)
}

// Repay is a paid mutator transaction binding the contract method 0x371fd8e6.
//
// Solidity: function repay(uint256 amount) returns()
func (_IntegratedVault *IntegratedVaultTransactorSession) Repay(amount *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.Contract.Repay(&_IntegratedVault.TransactOpts, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IntegratedVault *IntegratedVaultTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _IntegratedVault.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IntegratedVault *IntegratedVaultSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IntegratedVault.Contract.TransferOwnership(&_IntegratedVault.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IntegratedVault *IntegratedVaultTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IntegratedVault.Contract.TransferOwnership(&_IntegratedVault.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IntegratedVault *IntegratedVaultTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IntegratedVault.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IntegratedVault *IntegratedVaultSession) Unpause() (*types.Transaction, error) {
	return _IntegratedVault.Contract.Unpause(&_IntegratedVault.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IntegratedVault *IntegratedVaultTransactorSession) Unpause() (*types.Transaction, error) {
	return _IntegratedVault.Contract.Unpause(&_IntegratedVault.TransactOpts)
}

// UpdatePriceOracle is a paid mutator transaction binding the contract method 0x39d1fc82.
//
// Solidity: function updatePriceOracle(address _newOracle) returns()
func (_IntegratedVault *IntegratedVaultTransactor) UpdatePriceOracle(opts *bind.TransactOpts, _newOracle common.Address) (*types.Transaction, error) {
	return _IntegratedVault.contract.Transact(opts, "updatePriceOracle", _newOracle)
}

// UpdatePriceOracle is a paid mutator transaction binding the contract method 0x39d1fc82.
//
// Solidity: function updatePriceOracle(address _newOracle) returns()
func (_IntegratedVault *IntegratedVaultSession) UpdatePriceOracle(_newOracle common.Address) (*types.Transaction, error) {
	return _IntegratedVault.Contract.UpdatePriceOracle(&_IntegratedVault.TransactOpts, _newOracle)
}

// UpdatePriceOracle is a paid mutator transaction binding the contract method 0x39d1fc82.
//
// Solidity: function updatePriceOracle(address _newOracle) returns()
func (_IntegratedVault *IntegratedVaultTransactorSession) UpdatePriceOracle(_newOracle common.Address) (*types.Transaction, error) {
	return _IntegratedVault.Contract.UpdatePriceOracle(&_IntegratedVault.TransactOpts, _newOracle)
}

// UpdateStateAggregator is a paid mutator transaction binding the contract method 0x5ec71e77.
//
// Solidity: function updateStateAggregator(address _newAggregator) returns()
func (_IntegratedVault *IntegratedVaultTransactor) UpdateStateAggregator(opts *bind.TransactOpts, _newAggregator common.Address) (*types.Transaction, error) {
	return _IntegratedVault.contract.Transact(opts, "updateStateAggregator", _newAggregator)
}

// UpdateStateAggregator is a paid mutator transaction binding the contract method 0x5ec71e77.
//
// Solidity: function updateStateAggregator(address _newAggregator) returns()
func (_IntegratedVault *IntegratedVaultSession) UpdateStateAggregator(_newAggregator common.Address) (*types.Transaction, error) {
	return _IntegratedVault.Contract.UpdateStateAggregator(&_IntegratedVault.TransactOpts, _newAggregator)
}

// UpdateStateAggregator is a paid mutator transaction binding the contract method 0x5ec71e77.
//
// Solidity: function updateStateAggregator(address _newAggregator) returns()
func (_IntegratedVault *IntegratedVaultTransactorSession) UpdateStateAggregator(_newAggregator common.Address) (*types.Transaction, error) {
	return _IntegratedVault.Contract.UpdateStateAggregator(&_IntegratedVault.TransactOpts, _newAggregator)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_IntegratedVault *IntegratedVaultTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_IntegratedVault *IntegratedVaultSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.Contract.Withdraw(&_IntegratedVault.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_IntegratedVault *IntegratedVaultTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _IntegratedVault.Contract.Withdraw(&_IntegratedVault.TransactOpts, amount)
}

// IntegratedVaultBorrowedIterator is returned from FilterBorrowed and is used to iterate over the raw logs and unpacked data for Borrowed events raised by the IntegratedVault contract.
type IntegratedVaultBorrowedIterator struct {
	Event *IntegratedVaultBorrowed // Event containing the contract specifics and raw log

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
func (it *IntegratedVaultBorrowedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntegratedVaultBorrowed)
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
		it.Event = new(IntegratedVaultBorrowed)
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
func (it *IntegratedVaultBorrowedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntegratedVaultBorrowedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntegratedVaultBorrowed represents a Borrowed event raised by the IntegratedVault contract.
type IntegratedVaultBorrowed struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBorrowed is a free log retrieval operation binding the contract event 0xac59582e5396aca512fa873a2047e7f4c80f8f55d4a06cb34a78a0187f62719f.
//
// Solidity: event Borrowed(address indexed user, uint256 amount)
func (_IntegratedVault *IntegratedVaultFilterer) FilterBorrowed(opts *bind.FilterOpts, user []common.Address) (*IntegratedVaultBorrowedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _IntegratedVault.contract.FilterLogs(opts, "Borrowed", userRule)
	if err != nil {
		return nil, err
	}
	return &IntegratedVaultBorrowedIterator{contract: _IntegratedVault.contract, event: "Borrowed", logs: logs, sub: sub}, nil
}

// WatchBorrowed is a free log subscription operation binding the contract event 0xac59582e5396aca512fa873a2047e7f4c80f8f55d4a06cb34a78a0187f62719f.
//
// Solidity: event Borrowed(address indexed user, uint256 amount)
func (_IntegratedVault *IntegratedVaultFilterer) WatchBorrowed(opts *bind.WatchOpts, sink chan<- *IntegratedVaultBorrowed, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _IntegratedVault.contract.WatchLogs(opts, "Borrowed", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntegratedVaultBorrowed)
				if err := _IntegratedVault.contract.UnpackLog(event, "Borrowed", log); err != nil {
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
func (_IntegratedVault *IntegratedVaultFilterer) ParseBorrowed(log types.Log) (*IntegratedVaultBorrowed, error) {
	event := new(IntegratedVaultBorrowed)
	if err := _IntegratedVault.contract.UnpackLog(event, "Borrowed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntegratedVaultDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the IntegratedVault contract.
type IntegratedVaultDepositedIterator struct {
	Event *IntegratedVaultDeposited // Event containing the contract specifics and raw log

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
func (it *IntegratedVaultDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntegratedVaultDeposited)
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
		it.Event = new(IntegratedVaultDeposited)
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
func (it *IntegratedVaultDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntegratedVaultDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntegratedVaultDeposited represents a Deposited event raised by the IntegratedVault contract.
type IntegratedVaultDeposited struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed user, uint256 amount)
func (_IntegratedVault *IntegratedVaultFilterer) FilterDeposited(opts *bind.FilterOpts, user []common.Address) (*IntegratedVaultDepositedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _IntegratedVault.contract.FilterLogs(opts, "Deposited", userRule)
	if err != nil {
		return nil, err
	}
	return &IntegratedVaultDepositedIterator{contract: _IntegratedVault.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed user, uint256 amount)
func (_IntegratedVault *IntegratedVaultFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *IntegratedVaultDeposited, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _IntegratedVault.contract.WatchLogs(opts, "Deposited", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntegratedVaultDeposited)
				if err := _IntegratedVault.contract.UnpackLog(event, "Deposited", log); err != nil {
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

// ParseDeposited is a log parse operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed user, uint256 amount)
func (_IntegratedVault *IntegratedVaultFilterer) ParseDeposited(log types.Log) (*IntegratedVaultDeposited, error) {
	event := new(IntegratedVaultDeposited)
	if err := _IntegratedVault.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntegratedVaultInterestAccruedIterator is returned from FilterInterestAccrued and is used to iterate over the raw logs and unpacked data for InterestAccrued events raised by the IntegratedVault contract.
type IntegratedVaultInterestAccruedIterator struct {
	Event *IntegratedVaultInterestAccrued // Event containing the contract specifics and raw log

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
func (it *IntegratedVaultInterestAccruedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntegratedVaultInterestAccrued)
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
		it.Event = new(IntegratedVaultInterestAccrued)
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
func (it *IntegratedVaultInterestAccruedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntegratedVaultInterestAccruedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntegratedVaultInterestAccrued represents a InterestAccrued event raised by the IntegratedVault contract.
type IntegratedVaultInterestAccrued struct {
	User           common.Address
	InterestAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterInterestAccrued is a free log retrieval operation binding the contract event 0x5e804d42ae3b860f881d11cb44a4bb1f2f0d5b3d081f5539a32d6f97b629d978.
//
// Solidity: event InterestAccrued(address indexed user, uint256 interestAmount)
func (_IntegratedVault *IntegratedVaultFilterer) FilterInterestAccrued(opts *bind.FilterOpts, user []common.Address) (*IntegratedVaultInterestAccruedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _IntegratedVault.contract.FilterLogs(opts, "InterestAccrued", userRule)
	if err != nil {
		return nil, err
	}
	return &IntegratedVaultInterestAccruedIterator{contract: _IntegratedVault.contract, event: "InterestAccrued", logs: logs, sub: sub}, nil
}

// WatchInterestAccrued is a free log subscription operation binding the contract event 0x5e804d42ae3b860f881d11cb44a4bb1f2f0d5b3d081f5539a32d6f97b629d978.
//
// Solidity: event InterestAccrued(address indexed user, uint256 interestAmount)
func (_IntegratedVault *IntegratedVaultFilterer) WatchInterestAccrued(opts *bind.WatchOpts, sink chan<- *IntegratedVaultInterestAccrued, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _IntegratedVault.contract.WatchLogs(opts, "InterestAccrued", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntegratedVaultInterestAccrued)
				if err := _IntegratedVault.contract.UnpackLog(event, "InterestAccrued", log); err != nil {
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

// ParseInterestAccrued is a log parse operation binding the contract event 0x5e804d42ae3b860f881d11cb44a4bb1f2f0d5b3d081f5539a32d6f97b629d978.
//
// Solidity: event InterestAccrued(address indexed user, uint256 interestAmount)
func (_IntegratedVault *IntegratedVaultFilterer) ParseInterestAccrued(log types.Log) (*IntegratedVaultInterestAccrued, error) {
	event := new(IntegratedVaultInterestAccrued)
	if err := _IntegratedVault.contract.UnpackLog(event, "InterestAccrued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntegratedVaultLiquidatedIterator is returned from FilterLiquidated and is used to iterate over the raw logs and unpacked data for Liquidated events raised by the IntegratedVault contract.
type IntegratedVaultLiquidatedIterator struct {
	Event *IntegratedVaultLiquidated // Event containing the contract specifics and raw log

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
func (it *IntegratedVaultLiquidatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntegratedVaultLiquidated)
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
		it.Event = new(IntegratedVaultLiquidated)
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
func (it *IntegratedVaultLiquidatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntegratedVaultLiquidatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntegratedVaultLiquidated represents a Liquidated event raised by the IntegratedVault contract.
type IntegratedVaultLiquidated struct {
	User             common.Address
	Liquidator       common.Address
	DebtRepaid       *big.Int
	CollateralSeized *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLiquidated is a free log retrieval operation binding the contract event 0x1f0c6615429d1cdae0dfa233abf91d3b31cdbdd82c8081389832a61e1072f1ea.
//
// Solidity: event Liquidated(address indexed user, address indexed liquidator, uint256 debtRepaid, uint256 collateralSeized)
func (_IntegratedVault *IntegratedVaultFilterer) FilterLiquidated(opts *bind.FilterOpts, user []common.Address, liquidator []common.Address) (*IntegratedVaultLiquidatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var liquidatorRule []interface{}
	for _, liquidatorItem := range liquidator {
		liquidatorRule = append(liquidatorRule, liquidatorItem)
	}

	logs, sub, err := _IntegratedVault.contract.FilterLogs(opts, "Liquidated", userRule, liquidatorRule)
	if err != nil {
		return nil, err
	}
	return &IntegratedVaultLiquidatedIterator{contract: _IntegratedVault.contract, event: "Liquidated", logs: logs, sub: sub}, nil
}

// WatchLiquidated is a free log subscription operation binding the contract event 0x1f0c6615429d1cdae0dfa233abf91d3b31cdbdd82c8081389832a61e1072f1ea.
//
// Solidity: event Liquidated(address indexed user, address indexed liquidator, uint256 debtRepaid, uint256 collateralSeized)
func (_IntegratedVault *IntegratedVaultFilterer) WatchLiquidated(opts *bind.WatchOpts, sink chan<- *IntegratedVaultLiquidated, user []common.Address, liquidator []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var liquidatorRule []interface{}
	for _, liquidatorItem := range liquidator {
		liquidatorRule = append(liquidatorRule, liquidatorItem)
	}

	logs, sub, err := _IntegratedVault.contract.WatchLogs(opts, "Liquidated", userRule, liquidatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntegratedVaultLiquidated)
				if err := _IntegratedVault.contract.UnpackLog(event, "Liquidated", log); err != nil {
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

// ParseLiquidated is a log parse operation binding the contract event 0x1f0c6615429d1cdae0dfa233abf91d3b31cdbdd82c8081389832a61e1072f1ea.
//
// Solidity: event Liquidated(address indexed user, address indexed liquidator, uint256 debtRepaid, uint256 collateralSeized)
func (_IntegratedVault *IntegratedVaultFilterer) ParseLiquidated(log types.Log) (*IntegratedVaultLiquidated, error) {
	event := new(IntegratedVaultLiquidated)
	if err := _IntegratedVault.contract.UnpackLog(event, "Liquidated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntegratedVaultOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the IntegratedVault contract.
type IntegratedVaultOwnershipTransferredIterator struct {
	Event *IntegratedVaultOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *IntegratedVaultOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntegratedVaultOwnershipTransferred)
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
		it.Event = new(IntegratedVaultOwnershipTransferred)
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
func (it *IntegratedVaultOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntegratedVaultOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntegratedVaultOwnershipTransferred represents a OwnershipTransferred event raised by the IntegratedVault contract.
type IntegratedVaultOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IntegratedVault *IntegratedVaultFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IntegratedVaultOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IntegratedVault.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IntegratedVaultOwnershipTransferredIterator{contract: _IntegratedVault.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IntegratedVault *IntegratedVaultFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IntegratedVaultOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IntegratedVault.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntegratedVaultOwnershipTransferred)
				if err := _IntegratedVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_IntegratedVault *IntegratedVaultFilterer) ParseOwnershipTransferred(log types.Log) (*IntegratedVaultOwnershipTransferred, error) {
	event := new(IntegratedVaultOwnershipTransferred)
	if err := _IntegratedVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntegratedVaultPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the IntegratedVault contract.
type IntegratedVaultPausedIterator struct {
	Event *IntegratedVaultPaused // Event containing the contract specifics and raw log

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
func (it *IntegratedVaultPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntegratedVaultPaused)
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
		it.Event = new(IntegratedVaultPaused)
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
func (it *IntegratedVaultPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntegratedVaultPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntegratedVaultPaused represents a Paused event raised by the IntegratedVault contract.
type IntegratedVaultPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_IntegratedVault *IntegratedVaultFilterer) FilterPaused(opts *bind.FilterOpts) (*IntegratedVaultPausedIterator, error) {

	logs, sub, err := _IntegratedVault.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &IntegratedVaultPausedIterator{contract: _IntegratedVault.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_IntegratedVault *IntegratedVaultFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *IntegratedVaultPaused) (event.Subscription, error) {

	logs, sub, err := _IntegratedVault.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntegratedVaultPaused)
				if err := _IntegratedVault.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_IntegratedVault *IntegratedVaultFilterer) ParsePaused(log types.Log) (*IntegratedVaultPaused, error) {
	event := new(IntegratedVaultPaused)
	if err := _IntegratedVault.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntegratedVaultRepaidIterator is returned from FilterRepaid and is used to iterate over the raw logs and unpacked data for Repaid events raised by the IntegratedVault contract.
type IntegratedVaultRepaidIterator struct {
	Event *IntegratedVaultRepaid // Event containing the contract specifics and raw log

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
func (it *IntegratedVaultRepaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntegratedVaultRepaid)
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
		it.Event = new(IntegratedVaultRepaid)
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
func (it *IntegratedVaultRepaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntegratedVaultRepaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntegratedVaultRepaid represents a Repaid event raised by the IntegratedVault contract.
type IntegratedVaultRepaid struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRepaid is a free log retrieval operation binding the contract event 0x0516911bcc3a0a7412a44601057c0a0a1ec628bde049a84284bc428866534488.
//
// Solidity: event Repaid(address indexed user, uint256 amount)
func (_IntegratedVault *IntegratedVaultFilterer) FilterRepaid(opts *bind.FilterOpts, user []common.Address) (*IntegratedVaultRepaidIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _IntegratedVault.contract.FilterLogs(opts, "Repaid", userRule)
	if err != nil {
		return nil, err
	}
	return &IntegratedVaultRepaidIterator{contract: _IntegratedVault.contract, event: "Repaid", logs: logs, sub: sub}, nil
}

// WatchRepaid is a free log subscription operation binding the contract event 0x0516911bcc3a0a7412a44601057c0a0a1ec628bde049a84284bc428866534488.
//
// Solidity: event Repaid(address indexed user, uint256 amount)
func (_IntegratedVault *IntegratedVaultFilterer) WatchRepaid(opts *bind.WatchOpts, sink chan<- *IntegratedVaultRepaid, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _IntegratedVault.contract.WatchLogs(opts, "Repaid", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntegratedVaultRepaid)
				if err := _IntegratedVault.contract.UnpackLog(event, "Repaid", log); err != nil {
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
func (_IntegratedVault *IntegratedVaultFilterer) ParseRepaid(log types.Log) (*IntegratedVaultRepaid, error) {
	event := new(IntegratedVaultRepaid)
	if err := _IntegratedVault.contract.UnpackLog(event, "Repaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntegratedVaultUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the IntegratedVault contract.
type IntegratedVaultUnpausedIterator struct {
	Event *IntegratedVaultUnpaused // Event containing the contract specifics and raw log

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
func (it *IntegratedVaultUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntegratedVaultUnpaused)
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
		it.Event = new(IntegratedVaultUnpaused)
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
func (it *IntegratedVaultUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntegratedVaultUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntegratedVaultUnpaused represents a Unpaused event raised by the IntegratedVault contract.
type IntegratedVaultUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_IntegratedVault *IntegratedVaultFilterer) FilterUnpaused(opts *bind.FilterOpts) (*IntegratedVaultUnpausedIterator, error) {

	logs, sub, err := _IntegratedVault.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &IntegratedVaultUnpausedIterator{contract: _IntegratedVault.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_IntegratedVault *IntegratedVaultFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IntegratedVaultUnpaused) (event.Subscription, error) {

	logs, sub, err := _IntegratedVault.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntegratedVaultUnpaused)
				if err := _IntegratedVault.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_IntegratedVault *IntegratedVaultFilterer) ParseUnpaused(log types.Log) (*IntegratedVaultUnpaused, error) {
	event := new(IntegratedVaultUnpaused)
	if err := _IntegratedVault.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntegratedVaultWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the IntegratedVault contract.
type IntegratedVaultWithdrawnIterator struct {
	Event *IntegratedVaultWithdrawn // Event containing the contract specifics and raw log

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
func (it *IntegratedVaultWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntegratedVaultWithdrawn)
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
		it.Event = new(IntegratedVaultWithdrawn)
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
func (it *IntegratedVaultWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntegratedVaultWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntegratedVaultWithdrawn represents a Withdrawn event raised by the IntegratedVault contract.
type IntegratedVaultWithdrawn struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed user, uint256 amount)
func (_IntegratedVault *IntegratedVaultFilterer) FilterWithdrawn(opts *bind.FilterOpts, user []common.Address) (*IntegratedVaultWithdrawnIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _IntegratedVault.contract.FilterLogs(opts, "Withdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return &IntegratedVaultWithdrawnIterator{contract: _IntegratedVault.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed user, uint256 amount)
func (_IntegratedVault *IntegratedVaultFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *IntegratedVaultWithdrawn, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _IntegratedVault.contract.WatchLogs(opts, "Withdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntegratedVaultWithdrawn)
				if err := _IntegratedVault.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed user, uint256 amount)
func (_IntegratedVault *IntegratedVaultFilterer) ParseWithdrawn(log types.Log) (*IntegratedVaultWithdrawn, error) {
	event := new(IntegratedVaultWithdrawn)
	if err := _IntegratedVault.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
