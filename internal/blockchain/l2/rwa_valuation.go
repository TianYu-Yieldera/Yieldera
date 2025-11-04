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

// IRWAValuationOracleConfig is an auto generated low-level Go binding around an user-defined struct.
type IRWAValuationOracleConfig struct {
	OracleAddress      common.Address
	Heartbeat          *big.Int
	DeviationThreshold *big.Int
	IsActive           bool
}

// IRWAValuationValuation is an auto generated low-level Go binding around an user-defined struct.
type IRWAValuationValuation struct {
	Value      *big.Int
	Timestamp  *big.Int
	Method     uint8
	Valuator   common.Address
	ReportHash string
	Confidence *big.Int
}

// RWAValuationValuationDispute is an auto generated low-level Go binding around an user-defined struct.
type RWAValuationValuationDispute struct {
	DisputedValue *big.Int
	Disputer      common.Address
	Reason        string
	Timestamp     *big.Int
	Resolved      bool
}

// RWAValuationMetaData contains all meta data concerning the RWAValuation contract.
var RWAValuationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"heartbeat\",\"type\":\"uint256\"}],\"name\":\"OracleConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"disputedValue\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"disputer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"ValuationDisputed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumIRWAValuation.ValuationMethod\",\"name\":\"method\",\"type\":\"uint8\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"valuator\",\"type\":\"address\"}],\"name\":\"ValuationUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"valuator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"}],\"name\":\"ValuatorAuthorized\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CONFIDENCE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_VALUATION_AGE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_CONFIDENCE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ORACLE_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"VALUATOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"valuator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"}],\"name\":\"authorizeValuator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"heartbeat\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deviationThreshold\",\"type\":\"uint256\"}],\"name\":\"configureOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"disableOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"disputeValuation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getCurrentValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getDisputes\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"disputedValue\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"disputer\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"resolved\",\"type\":\"bool\"}],\"internalType\":\"structRWAValuation.ValuationDispute[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getLastValuation\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"enumIRWAValuation.ValuationMethod\",\"name\":\"method\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"valuator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"reportHash\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"confidence\",\"type\":\"uint256\"}],\"internalType\":\"structIRWAValuation.Valuation\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getOracleConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"heartbeat\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deviationThreshold\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"}],\"internalType\":\"structIRWAValuation.OracleConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getPricePerToken\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"getTWAP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getValuationCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"getValuationHistory\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"enumIRWAValuation.ValuationMethod\",\"name\":\"method\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"valuator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"reportHash\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"confidence\",\"type\":\"uint256\"}],\"internalType\":\"structIRWAValuation.Valuation[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"isValuationStale\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"valuator\",\"type\":\"address\"}],\"name\":\"isValuatorAuthorized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"requestRevaluation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"disputeIndex\",\"type\":\"uint256\"}],\"name\":\"resolveDispute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalAssets\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalValuations\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"updateFromOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"},{\"internalType\":\"enumIRWAValuation.ValuationMethod\",\"name\":\"method\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"reportHash\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"confidence\",\"type\":\"uint256\"}],\"name\":\"updateValuation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RWAValuationABI is the input ABI used to generate the binding from.
// Deprecated: Use RWAValuationMetaData.ABI instead.
var RWAValuationABI = RWAValuationMetaData.ABI

// RWAValuation is an auto generated Go binding around an Ethereum contract.
type RWAValuation struct {
	RWAValuationCaller     // Read-only binding to the contract
	RWAValuationTransactor // Write-only binding to the contract
	RWAValuationFilterer   // Log filterer for contract events
}

// RWAValuationCaller is an auto generated read-only Go binding around an Ethereum contract.
type RWAValuationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAValuationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RWAValuationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAValuationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RWAValuationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAValuationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RWAValuationSession struct {
	Contract     *RWAValuation     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RWAValuationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RWAValuationCallerSession struct {
	Contract *RWAValuationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// RWAValuationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RWAValuationTransactorSession struct {
	Contract     *RWAValuationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// RWAValuationRaw is an auto generated low-level Go binding around an Ethereum contract.
type RWAValuationRaw struct {
	Contract *RWAValuation // Generic contract binding to access the raw methods on
}

// RWAValuationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RWAValuationCallerRaw struct {
	Contract *RWAValuationCaller // Generic read-only contract binding to access the raw methods on
}

// RWAValuationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RWAValuationTransactorRaw struct {
	Contract *RWAValuationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRWAValuation creates a new instance of RWAValuation, bound to a specific deployed contract.
func NewRWAValuation(address common.Address, backend bind.ContractBackend) (*RWAValuation, error) {
	contract, err := bindRWAValuation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RWAValuation{RWAValuationCaller: RWAValuationCaller{contract: contract}, RWAValuationTransactor: RWAValuationTransactor{contract: contract}, RWAValuationFilterer: RWAValuationFilterer{contract: contract}}, nil
}

// NewRWAValuationCaller creates a new read-only instance of RWAValuation, bound to a specific deployed contract.
func NewRWAValuationCaller(address common.Address, caller bind.ContractCaller) (*RWAValuationCaller, error) {
	contract, err := bindRWAValuation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RWAValuationCaller{contract: contract}, nil
}

// NewRWAValuationTransactor creates a new write-only instance of RWAValuation, bound to a specific deployed contract.
func NewRWAValuationTransactor(address common.Address, transactor bind.ContractTransactor) (*RWAValuationTransactor, error) {
	contract, err := bindRWAValuation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RWAValuationTransactor{contract: contract}, nil
}

// NewRWAValuationFilterer creates a new log filterer instance of RWAValuation, bound to a specific deployed contract.
func NewRWAValuationFilterer(address common.Address, filterer bind.ContractFilterer) (*RWAValuationFilterer, error) {
	contract, err := bindRWAValuation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RWAValuationFilterer{contract: contract}, nil
}

// bindRWAValuation binds a generic wrapper to an already deployed contract.
func bindRWAValuation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RWAValuationMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RWAValuation *RWAValuationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RWAValuation.Contract.RWAValuationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RWAValuation *RWAValuationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAValuation.Contract.RWAValuationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RWAValuation *RWAValuationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RWAValuation.Contract.RWAValuationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RWAValuation *RWAValuationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RWAValuation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RWAValuation *RWAValuationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAValuation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RWAValuation *RWAValuationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RWAValuation.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAValuation *RWAValuationCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAValuation *RWAValuationSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RWAValuation.Contract.DEFAULTADMINROLE(&_RWAValuation.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAValuation *RWAValuationCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RWAValuation.Contract.DEFAULTADMINROLE(&_RWAValuation.CallOpts)
}

// MAXCONFIDENCE is a free data retrieval call binding the contract method 0x6c78449c.
//
// Solidity: function MAX_CONFIDENCE() view returns(uint256)
func (_RWAValuation *RWAValuationCaller) MAXCONFIDENCE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "MAX_CONFIDENCE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXCONFIDENCE is a free data retrieval call binding the contract method 0x6c78449c.
//
// Solidity: function MAX_CONFIDENCE() view returns(uint256)
func (_RWAValuation *RWAValuationSession) MAXCONFIDENCE() (*big.Int, error) {
	return _RWAValuation.Contract.MAXCONFIDENCE(&_RWAValuation.CallOpts)
}

// MAXCONFIDENCE is a free data retrieval call binding the contract method 0x6c78449c.
//
// Solidity: function MAX_CONFIDENCE() view returns(uint256)
func (_RWAValuation *RWAValuationCallerSession) MAXCONFIDENCE() (*big.Int, error) {
	return _RWAValuation.Contract.MAXCONFIDENCE(&_RWAValuation.CallOpts)
}

// MAXVALUATIONAGE is a free data retrieval call binding the contract method 0x25dbd5a4.
//
// Solidity: function MAX_VALUATION_AGE() view returns(uint256)
func (_RWAValuation *RWAValuationCaller) MAXVALUATIONAGE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "MAX_VALUATION_AGE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXVALUATIONAGE is a free data retrieval call binding the contract method 0x25dbd5a4.
//
// Solidity: function MAX_VALUATION_AGE() view returns(uint256)
func (_RWAValuation *RWAValuationSession) MAXVALUATIONAGE() (*big.Int, error) {
	return _RWAValuation.Contract.MAXVALUATIONAGE(&_RWAValuation.CallOpts)
}

// MAXVALUATIONAGE is a free data retrieval call binding the contract method 0x25dbd5a4.
//
// Solidity: function MAX_VALUATION_AGE() view returns(uint256)
func (_RWAValuation *RWAValuationCallerSession) MAXVALUATIONAGE() (*big.Int, error) {
	return _RWAValuation.Contract.MAXVALUATIONAGE(&_RWAValuation.CallOpts)
}

// MINCONFIDENCE is a free data retrieval call binding the contract method 0x41f93757.
//
// Solidity: function MIN_CONFIDENCE() view returns(uint256)
func (_RWAValuation *RWAValuationCaller) MINCONFIDENCE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "MIN_CONFIDENCE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINCONFIDENCE is a free data retrieval call binding the contract method 0x41f93757.
//
// Solidity: function MIN_CONFIDENCE() view returns(uint256)
func (_RWAValuation *RWAValuationSession) MINCONFIDENCE() (*big.Int, error) {
	return _RWAValuation.Contract.MINCONFIDENCE(&_RWAValuation.CallOpts)
}

// MINCONFIDENCE is a free data retrieval call binding the contract method 0x41f93757.
//
// Solidity: function MIN_CONFIDENCE() view returns(uint256)
func (_RWAValuation *RWAValuationCallerSession) MINCONFIDENCE() (*big.Int, error) {
	return _RWAValuation.Contract.MINCONFIDENCE(&_RWAValuation.CallOpts)
}

// ORACLEADMINROLE is a free data retrieval call binding the contract method 0x8003a94f.
//
// Solidity: function ORACLE_ADMIN_ROLE() view returns(bytes32)
func (_RWAValuation *RWAValuationCaller) ORACLEADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "ORACLE_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ORACLEADMINROLE is a free data retrieval call binding the contract method 0x8003a94f.
//
// Solidity: function ORACLE_ADMIN_ROLE() view returns(bytes32)
func (_RWAValuation *RWAValuationSession) ORACLEADMINROLE() ([32]byte, error) {
	return _RWAValuation.Contract.ORACLEADMINROLE(&_RWAValuation.CallOpts)
}

// ORACLEADMINROLE is a free data retrieval call binding the contract method 0x8003a94f.
//
// Solidity: function ORACLE_ADMIN_ROLE() view returns(bytes32)
func (_RWAValuation *RWAValuationCallerSession) ORACLEADMINROLE() ([32]byte, error) {
	return _RWAValuation.Contract.ORACLEADMINROLE(&_RWAValuation.CallOpts)
}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_RWAValuation *RWAValuationCaller) PRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_RWAValuation *RWAValuationSession) PRECISION() (*big.Int, error) {
	return _RWAValuation.Contract.PRECISION(&_RWAValuation.CallOpts)
}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_RWAValuation *RWAValuationCallerSession) PRECISION() (*big.Int, error) {
	return _RWAValuation.Contract.PRECISION(&_RWAValuation.CallOpts)
}

// VALUATORROLE is a free data retrieval call binding the contract method 0x7ef0c851.
//
// Solidity: function VALUATOR_ROLE() view returns(bytes32)
func (_RWAValuation *RWAValuationCaller) VALUATORROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "VALUATOR_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// VALUATORROLE is a free data retrieval call binding the contract method 0x7ef0c851.
//
// Solidity: function VALUATOR_ROLE() view returns(bytes32)
func (_RWAValuation *RWAValuationSession) VALUATORROLE() ([32]byte, error) {
	return _RWAValuation.Contract.VALUATORROLE(&_RWAValuation.CallOpts)
}

// VALUATORROLE is a free data retrieval call binding the contract method 0x7ef0c851.
//
// Solidity: function VALUATOR_ROLE() view returns(bytes32)
func (_RWAValuation *RWAValuationCallerSession) VALUATORROLE() ([32]byte, error) {
	return _RWAValuation.Contract.VALUATORROLE(&_RWAValuation.CallOpts)
}

// GetCurrentValue is a free data retrieval call binding the contract method 0x3fcad964.
//
// Solidity: function getCurrentValue(uint256 assetId) view returns(uint256)
func (_RWAValuation *RWAValuationCaller) GetCurrentValue(opts *bind.CallOpts, assetId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "getCurrentValue", assetId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentValue is a free data retrieval call binding the contract method 0x3fcad964.
//
// Solidity: function getCurrentValue(uint256 assetId) view returns(uint256)
func (_RWAValuation *RWAValuationSession) GetCurrentValue(assetId *big.Int) (*big.Int, error) {
	return _RWAValuation.Contract.GetCurrentValue(&_RWAValuation.CallOpts, assetId)
}

// GetCurrentValue is a free data retrieval call binding the contract method 0x3fcad964.
//
// Solidity: function getCurrentValue(uint256 assetId) view returns(uint256)
func (_RWAValuation *RWAValuationCallerSession) GetCurrentValue(assetId *big.Int) (*big.Int, error) {
	return _RWAValuation.Contract.GetCurrentValue(&_RWAValuation.CallOpts, assetId)
}

// GetDisputes is a free data retrieval call binding the contract method 0x623d1e38.
//
// Solidity: function getDisputes(uint256 assetId) view returns((uint256,address,string,uint256,bool)[])
func (_RWAValuation *RWAValuationCaller) GetDisputes(opts *bind.CallOpts, assetId *big.Int) ([]RWAValuationValuationDispute, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "getDisputes", assetId)

	if err != nil {
		return *new([]RWAValuationValuationDispute), err
	}

	out0 := *abi.ConvertType(out[0], new([]RWAValuationValuationDispute)).(*[]RWAValuationValuationDispute)

	return out0, err

}

// GetDisputes is a free data retrieval call binding the contract method 0x623d1e38.
//
// Solidity: function getDisputes(uint256 assetId) view returns((uint256,address,string,uint256,bool)[])
func (_RWAValuation *RWAValuationSession) GetDisputes(assetId *big.Int) ([]RWAValuationValuationDispute, error) {
	return _RWAValuation.Contract.GetDisputes(&_RWAValuation.CallOpts, assetId)
}

// GetDisputes is a free data retrieval call binding the contract method 0x623d1e38.
//
// Solidity: function getDisputes(uint256 assetId) view returns((uint256,address,string,uint256,bool)[])
func (_RWAValuation *RWAValuationCallerSession) GetDisputes(assetId *big.Int) ([]RWAValuationValuationDispute, error) {
	return _RWAValuation.Contract.GetDisputes(&_RWAValuation.CallOpts, assetId)
}

// GetLastValuation is a free data retrieval call binding the contract method 0xbd6ff20f.
//
// Solidity: function getLastValuation(uint256 assetId) view returns((uint256,uint256,uint8,address,string,uint256))
func (_RWAValuation *RWAValuationCaller) GetLastValuation(opts *bind.CallOpts, assetId *big.Int) (IRWAValuationValuation, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "getLastValuation", assetId)

	if err != nil {
		return *new(IRWAValuationValuation), err
	}

	out0 := *abi.ConvertType(out[0], new(IRWAValuationValuation)).(*IRWAValuationValuation)

	return out0, err

}

// GetLastValuation is a free data retrieval call binding the contract method 0xbd6ff20f.
//
// Solidity: function getLastValuation(uint256 assetId) view returns((uint256,uint256,uint8,address,string,uint256))
func (_RWAValuation *RWAValuationSession) GetLastValuation(assetId *big.Int) (IRWAValuationValuation, error) {
	return _RWAValuation.Contract.GetLastValuation(&_RWAValuation.CallOpts, assetId)
}

// GetLastValuation is a free data retrieval call binding the contract method 0xbd6ff20f.
//
// Solidity: function getLastValuation(uint256 assetId) view returns((uint256,uint256,uint8,address,string,uint256))
func (_RWAValuation *RWAValuationCallerSession) GetLastValuation(assetId *big.Int) (IRWAValuationValuation, error) {
	return _RWAValuation.Contract.GetLastValuation(&_RWAValuation.CallOpts, assetId)
}

// GetOracleConfig is a free data retrieval call binding the contract method 0x3e1cf1d8.
//
// Solidity: function getOracleConfig(uint256 assetId) view returns((address,uint256,uint256,bool))
func (_RWAValuation *RWAValuationCaller) GetOracleConfig(opts *bind.CallOpts, assetId *big.Int) (IRWAValuationOracleConfig, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "getOracleConfig", assetId)

	if err != nil {
		return *new(IRWAValuationOracleConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IRWAValuationOracleConfig)).(*IRWAValuationOracleConfig)

	return out0, err

}

// GetOracleConfig is a free data retrieval call binding the contract method 0x3e1cf1d8.
//
// Solidity: function getOracleConfig(uint256 assetId) view returns((address,uint256,uint256,bool))
func (_RWAValuation *RWAValuationSession) GetOracleConfig(assetId *big.Int) (IRWAValuationOracleConfig, error) {
	return _RWAValuation.Contract.GetOracleConfig(&_RWAValuation.CallOpts, assetId)
}

// GetOracleConfig is a free data retrieval call binding the contract method 0x3e1cf1d8.
//
// Solidity: function getOracleConfig(uint256 assetId) view returns((address,uint256,uint256,bool))
func (_RWAValuation *RWAValuationCallerSession) GetOracleConfig(assetId *big.Int) (IRWAValuationOracleConfig, error) {
	return _RWAValuation.Contract.GetOracleConfig(&_RWAValuation.CallOpts, assetId)
}

// GetPricePerToken is a free data retrieval call binding the contract method 0xb8bfdf9f.
//
// Solidity: function getPricePerToken(uint256 assetId) view returns(uint256)
func (_RWAValuation *RWAValuationCaller) GetPricePerToken(opts *bind.CallOpts, assetId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "getPricePerToken", assetId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPricePerToken is a free data retrieval call binding the contract method 0xb8bfdf9f.
//
// Solidity: function getPricePerToken(uint256 assetId) view returns(uint256)
func (_RWAValuation *RWAValuationSession) GetPricePerToken(assetId *big.Int) (*big.Int, error) {
	return _RWAValuation.Contract.GetPricePerToken(&_RWAValuation.CallOpts, assetId)
}

// GetPricePerToken is a free data retrieval call binding the contract method 0xb8bfdf9f.
//
// Solidity: function getPricePerToken(uint256 assetId) view returns(uint256)
func (_RWAValuation *RWAValuationCallerSession) GetPricePerToken(assetId *big.Int) (*big.Int, error) {
	return _RWAValuation.Contract.GetPricePerToken(&_RWAValuation.CallOpts, assetId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAValuation *RWAValuationCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAValuation *RWAValuationSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RWAValuation.Contract.GetRoleAdmin(&_RWAValuation.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAValuation *RWAValuationCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RWAValuation.Contract.GetRoleAdmin(&_RWAValuation.CallOpts, role)
}

// GetTWAP is a free data retrieval call binding the contract method 0xeeca29b5.
//
// Solidity: function getTWAP(uint256 assetId, uint256 period) view returns(uint256)
func (_RWAValuation *RWAValuationCaller) GetTWAP(opts *bind.CallOpts, assetId *big.Int, period *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "getTWAP", assetId, period)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTWAP is a free data retrieval call binding the contract method 0xeeca29b5.
//
// Solidity: function getTWAP(uint256 assetId, uint256 period) view returns(uint256)
func (_RWAValuation *RWAValuationSession) GetTWAP(assetId *big.Int, period *big.Int) (*big.Int, error) {
	return _RWAValuation.Contract.GetTWAP(&_RWAValuation.CallOpts, assetId, period)
}

// GetTWAP is a free data retrieval call binding the contract method 0xeeca29b5.
//
// Solidity: function getTWAP(uint256 assetId, uint256 period) view returns(uint256)
func (_RWAValuation *RWAValuationCallerSession) GetTWAP(assetId *big.Int, period *big.Int) (*big.Int, error) {
	return _RWAValuation.Contract.GetTWAP(&_RWAValuation.CallOpts, assetId, period)
}

// GetValuationCount is a free data retrieval call binding the contract method 0xc960578f.
//
// Solidity: function getValuationCount(uint256 assetId) view returns(uint256)
func (_RWAValuation *RWAValuationCaller) GetValuationCount(opts *bind.CallOpts, assetId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "getValuationCount", assetId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetValuationCount is a free data retrieval call binding the contract method 0xc960578f.
//
// Solidity: function getValuationCount(uint256 assetId) view returns(uint256)
func (_RWAValuation *RWAValuationSession) GetValuationCount(assetId *big.Int) (*big.Int, error) {
	return _RWAValuation.Contract.GetValuationCount(&_RWAValuation.CallOpts, assetId)
}

// GetValuationCount is a free data retrieval call binding the contract method 0xc960578f.
//
// Solidity: function getValuationCount(uint256 assetId) view returns(uint256)
func (_RWAValuation *RWAValuationCallerSession) GetValuationCount(assetId *big.Int) (*big.Int, error) {
	return _RWAValuation.Contract.GetValuationCount(&_RWAValuation.CallOpts, assetId)
}

// GetValuationHistory is a free data retrieval call binding the contract method 0x9bdcf183.
//
// Solidity: function getValuationHistory(uint256 assetId, uint256 count) view returns((uint256,uint256,uint8,address,string,uint256)[])
func (_RWAValuation *RWAValuationCaller) GetValuationHistory(opts *bind.CallOpts, assetId *big.Int, count *big.Int) ([]IRWAValuationValuation, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "getValuationHistory", assetId, count)

	if err != nil {
		return *new([]IRWAValuationValuation), err
	}

	out0 := *abi.ConvertType(out[0], new([]IRWAValuationValuation)).(*[]IRWAValuationValuation)

	return out0, err

}

// GetValuationHistory is a free data retrieval call binding the contract method 0x9bdcf183.
//
// Solidity: function getValuationHistory(uint256 assetId, uint256 count) view returns((uint256,uint256,uint8,address,string,uint256)[])
func (_RWAValuation *RWAValuationSession) GetValuationHistory(assetId *big.Int, count *big.Int) ([]IRWAValuationValuation, error) {
	return _RWAValuation.Contract.GetValuationHistory(&_RWAValuation.CallOpts, assetId, count)
}

// GetValuationHistory is a free data retrieval call binding the contract method 0x9bdcf183.
//
// Solidity: function getValuationHistory(uint256 assetId, uint256 count) view returns((uint256,uint256,uint8,address,string,uint256)[])
func (_RWAValuation *RWAValuationCallerSession) GetValuationHistory(assetId *big.Int, count *big.Int) ([]IRWAValuationValuation, error) {
	return _RWAValuation.Contract.GetValuationHistory(&_RWAValuation.CallOpts, assetId, count)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAValuation *RWAValuationCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAValuation *RWAValuationSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RWAValuation.Contract.HasRole(&_RWAValuation.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAValuation *RWAValuationCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RWAValuation.Contract.HasRole(&_RWAValuation.CallOpts, role, account)
}

// IsValuationStale is a free data retrieval call binding the contract method 0xef274540.
//
// Solidity: function isValuationStale(uint256 assetId) view returns(bool)
func (_RWAValuation *RWAValuationCaller) IsValuationStale(opts *bind.CallOpts, assetId *big.Int) (bool, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "isValuationStale", assetId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValuationStale is a free data retrieval call binding the contract method 0xef274540.
//
// Solidity: function isValuationStale(uint256 assetId) view returns(bool)
func (_RWAValuation *RWAValuationSession) IsValuationStale(assetId *big.Int) (bool, error) {
	return _RWAValuation.Contract.IsValuationStale(&_RWAValuation.CallOpts, assetId)
}

// IsValuationStale is a free data retrieval call binding the contract method 0xef274540.
//
// Solidity: function isValuationStale(uint256 assetId) view returns(bool)
func (_RWAValuation *RWAValuationCallerSession) IsValuationStale(assetId *big.Int) (bool, error) {
	return _RWAValuation.Contract.IsValuationStale(&_RWAValuation.CallOpts, assetId)
}

// IsValuatorAuthorized is a free data retrieval call binding the contract method 0xb067ee18.
//
// Solidity: function isValuatorAuthorized(address valuator) view returns(bool)
func (_RWAValuation *RWAValuationCaller) IsValuatorAuthorized(opts *bind.CallOpts, valuator common.Address) (bool, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "isValuatorAuthorized", valuator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValuatorAuthorized is a free data retrieval call binding the contract method 0xb067ee18.
//
// Solidity: function isValuatorAuthorized(address valuator) view returns(bool)
func (_RWAValuation *RWAValuationSession) IsValuatorAuthorized(valuator common.Address) (bool, error) {
	return _RWAValuation.Contract.IsValuatorAuthorized(&_RWAValuation.CallOpts, valuator)
}

// IsValuatorAuthorized is a free data retrieval call binding the contract method 0xb067ee18.
//
// Solidity: function isValuatorAuthorized(address valuator) view returns(bool)
func (_RWAValuation *RWAValuationCallerSession) IsValuatorAuthorized(valuator common.Address) (bool, error) {
	return _RWAValuation.Contract.IsValuatorAuthorized(&_RWAValuation.CallOpts, valuator)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAValuation *RWAValuationCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAValuation *RWAValuationSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RWAValuation.Contract.SupportsInterface(&_RWAValuation.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAValuation *RWAValuationCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RWAValuation.Contract.SupportsInterface(&_RWAValuation.CallOpts, interfaceId)
}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_RWAValuation *RWAValuationCaller) TotalAssets(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "totalAssets")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_RWAValuation *RWAValuationSession) TotalAssets() (*big.Int, error) {
	return _RWAValuation.Contract.TotalAssets(&_RWAValuation.CallOpts)
}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_RWAValuation *RWAValuationCallerSession) TotalAssets() (*big.Int, error) {
	return _RWAValuation.Contract.TotalAssets(&_RWAValuation.CallOpts)
}

// TotalValuations is a free data retrieval call binding the contract method 0xc3888220.
//
// Solidity: function totalValuations() view returns(uint256)
func (_RWAValuation *RWAValuationCaller) TotalValuations(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAValuation.contract.Call(opts, &out, "totalValuations")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalValuations is a free data retrieval call binding the contract method 0xc3888220.
//
// Solidity: function totalValuations() view returns(uint256)
func (_RWAValuation *RWAValuationSession) TotalValuations() (*big.Int, error) {
	return _RWAValuation.Contract.TotalValuations(&_RWAValuation.CallOpts)
}

// TotalValuations is a free data retrieval call binding the contract method 0xc3888220.
//
// Solidity: function totalValuations() view returns(uint256)
func (_RWAValuation *RWAValuationCallerSession) TotalValuations() (*big.Int, error) {
	return _RWAValuation.Contract.TotalValuations(&_RWAValuation.CallOpts)
}

// AuthorizeValuator is a paid mutator transaction binding the contract method 0x05208bff.
//
// Solidity: function authorizeValuator(address valuator, bool authorized) returns()
func (_RWAValuation *RWAValuationTransactor) AuthorizeValuator(opts *bind.TransactOpts, valuator common.Address, authorized bool) (*types.Transaction, error) {
	return _RWAValuation.contract.Transact(opts, "authorizeValuator", valuator, authorized)
}

// AuthorizeValuator is a paid mutator transaction binding the contract method 0x05208bff.
//
// Solidity: function authorizeValuator(address valuator, bool authorized) returns()
func (_RWAValuation *RWAValuationSession) AuthorizeValuator(valuator common.Address, authorized bool) (*types.Transaction, error) {
	return _RWAValuation.Contract.AuthorizeValuator(&_RWAValuation.TransactOpts, valuator, authorized)
}

// AuthorizeValuator is a paid mutator transaction binding the contract method 0x05208bff.
//
// Solidity: function authorizeValuator(address valuator, bool authorized) returns()
func (_RWAValuation *RWAValuationTransactorSession) AuthorizeValuator(valuator common.Address, authorized bool) (*types.Transaction, error) {
	return _RWAValuation.Contract.AuthorizeValuator(&_RWAValuation.TransactOpts, valuator, authorized)
}

// ConfigureOracle is a paid mutator transaction binding the contract method 0xc30b1806.
//
// Solidity: function configureOracle(uint256 assetId, address oracleAddress, uint256 heartbeat, uint256 deviationThreshold) returns()
func (_RWAValuation *RWAValuationTransactor) ConfigureOracle(opts *bind.TransactOpts, assetId *big.Int, oracleAddress common.Address, heartbeat *big.Int, deviationThreshold *big.Int) (*types.Transaction, error) {
	return _RWAValuation.contract.Transact(opts, "configureOracle", assetId, oracleAddress, heartbeat, deviationThreshold)
}

// ConfigureOracle is a paid mutator transaction binding the contract method 0xc30b1806.
//
// Solidity: function configureOracle(uint256 assetId, address oracleAddress, uint256 heartbeat, uint256 deviationThreshold) returns()
func (_RWAValuation *RWAValuationSession) ConfigureOracle(assetId *big.Int, oracleAddress common.Address, heartbeat *big.Int, deviationThreshold *big.Int) (*types.Transaction, error) {
	return _RWAValuation.Contract.ConfigureOracle(&_RWAValuation.TransactOpts, assetId, oracleAddress, heartbeat, deviationThreshold)
}

// ConfigureOracle is a paid mutator transaction binding the contract method 0xc30b1806.
//
// Solidity: function configureOracle(uint256 assetId, address oracleAddress, uint256 heartbeat, uint256 deviationThreshold) returns()
func (_RWAValuation *RWAValuationTransactorSession) ConfigureOracle(assetId *big.Int, oracleAddress common.Address, heartbeat *big.Int, deviationThreshold *big.Int) (*types.Transaction, error) {
	return _RWAValuation.Contract.ConfigureOracle(&_RWAValuation.TransactOpts, assetId, oracleAddress, heartbeat, deviationThreshold)
}

// DisableOracle is a paid mutator transaction binding the contract method 0x7d1b5d9c.
//
// Solidity: function disableOracle(uint256 assetId) returns()
func (_RWAValuation *RWAValuationTransactor) DisableOracle(opts *bind.TransactOpts, assetId *big.Int) (*types.Transaction, error) {
	return _RWAValuation.contract.Transact(opts, "disableOracle", assetId)
}

// DisableOracle is a paid mutator transaction binding the contract method 0x7d1b5d9c.
//
// Solidity: function disableOracle(uint256 assetId) returns()
func (_RWAValuation *RWAValuationSession) DisableOracle(assetId *big.Int) (*types.Transaction, error) {
	return _RWAValuation.Contract.DisableOracle(&_RWAValuation.TransactOpts, assetId)
}

// DisableOracle is a paid mutator transaction binding the contract method 0x7d1b5d9c.
//
// Solidity: function disableOracle(uint256 assetId) returns()
func (_RWAValuation *RWAValuationTransactorSession) DisableOracle(assetId *big.Int) (*types.Transaction, error) {
	return _RWAValuation.Contract.DisableOracle(&_RWAValuation.TransactOpts, assetId)
}

// DisputeValuation is a paid mutator transaction binding the contract method 0x13761246.
//
// Solidity: function disputeValuation(uint256 assetId, string reason) returns()
func (_RWAValuation *RWAValuationTransactor) DisputeValuation(opts *bind.TransactOpts, assetId *big.Int, reason string) (*types.Transaction, error) {
	return _RWAValuation.contract.Transact(opts, "disputeValuation", assetId, reason)
}

// DisputeValuation is a paid mutator transaction binding the contract method 0x13761246.
//
// Solidity: function disputeValuation(uint256 assetId, string reason) returns()
func (_RWAValuation *RWAValuationSession) DisputeValuation(assetId *big.Int, reason string) (*types.Transaction, error) {
	return _RWAValuation.Contract.DisputeValuation(&_RWAValuation.TransactOpts, assetId, reason)
}

// DisputeValuation is a paid mutator transaction binding the contract method 0x13761246.
//
// Solidity: function disputeValuation(uint256 assetId, string reason) returns()
func (_RWAValuation *RWAValuationTransactorSession) DisputeValuation(assetId *big.Int, reason string) (*types.Transaction, error) {
	return _RWAValuation.Contract.DisputeValuation(&_RWAValuation.TransactOpts, assetId, reason)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAValuation *RWAValuationTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAValuation.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAValuation *RWAValuationSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAValuation.Contract.GrantRole(&_RWAValuation.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAValuation *RWAValuationTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAValuation.Contract.GrantRole(&_RWAValuation.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAValuation *RWAValuationTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAValuation.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAValuation *RWAValuationSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAValuation.Contract.RenounceRole(&_RWAValuation.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAValuation *RWAValuationTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAValuation.Contract.RenounceRole(&_RWAValuation.TransactOpts, role, callerConfirmation)
}

// RequestRevaluation is a paid mutator transaction binding the contract method 0xf02f1b22.
//
// Solidity: function requestRevaluation(uint256 assetId) returns()
func (_RWAValuation *RWAValuationTransactor) RequestRevaluation(opts *bind.TransactOpts, assetId *big.Int) (*types.Transaction, error) {
	return _RWAValuation.contract.Transact(opts, "requestRevaluation", assetId)
}

// RequestRevaluation is a paid mutator transaction binding the contract method 0xf02f1b22.
//
// Solidity: function requestRevaluation(uint256 assetId) returns()
func (_RWAValuation *RWAValuationSession) RequestRevaluation(assetId *big.Int) (*types.Transaction, error) {
	return _RWAValuation.Contract.RequestRevaluation(&_RWAValuation.TransactOpts, assetId)
}

// RequestRevaluation is a paid mutator transaction binding the contract method 0xf02f1b22.
//
// Solidity: function requestRevaluation(uint256 assetId) returns()
func (_RWAValuation *RWAValuationTransactorSession) RequestRevaluation(assetId *big.Int) (*types.Transaction, error) {
	return _RWAValuation.Contract.RequestRevaluation(&_RWAValuation.TransactOpts, assetId)
}

// ResolveDispute is a paid mutator transaction binding the contract method 0xbdc84ac3.
//
// Solidity: function resolveDispute(uint256 assetId, uint256 disputeIndex) returns()
func (_RWAValuation *RWAValuationTransactor) ResolveDispute(opts *bind.TransactOpts, assetId *big.Int, disputeIndex *big.Int) (*types.Transaction, error) {
	return _RWAValuation.contract.Transact(opts, "resolveDispute", assetId, disputeIndex)
}

// ResolveDispute is a paid mutator transaction binding the contract method 0xbdc84ac3.
//
// Solidity: function resolveDispute(uint256 assetId, uint256 disputeIndex) returns()
func (_RWAValuation *RWAValuationSession) ResolveDispute(assetId *big.Int, disputeIndex *big.Int) (*types.Transaction, error) {
	return _RWAValuation.Contract.ResolveDispute(&_RWAValuation.TransactOpts, assetId, disputeIndex)
}

// ResolveDispute is a paid mutator transaction binding the contract method 0xbdc84ac3.
//
// Solidity: function resolveDispute(uint256 assetId, uint256 disputeIndex) returns()
func (_RWAValuation *RWAValuationTransactorSession) ResolveDispute(assetId *big.Int, disputeIndex *big.Int) (*types.Transaction, error) {
	return _RWAValuation.Contract.ResolveDispute(&_RWAValuation.TransactOpts, assetId, disputeIndex)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAValuation *RWAValuationTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAValuation.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAValuation *RWAValuationSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAValuation.Contract.RevokeRole(&_RWAValuation.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAValuation *RWAValuationTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAValuation.Contract.RevokeRole(&_RWAValuation.TransactOpts, role, account)
}

// UpdateFromOracle is a paid mutator transaction binding the contract method 0x815d89eb.
//
// Solidity: function updateFromOracle(uint256 assetId) returns()
func (_RWAValuation *RWAValuationTransactor) UpdateFromOracle(opts *bind.TransactOpts, assetId *big.Int) (*types.Transaction, error) {
	return _RWAValuation.contract.Transact(opts, "updateFromOracle", assetId)
}

// UpdateFromOracle is a paid mutator transaction binding the contract method 0x815d89eb.
//
// Solidity: function updateFromOracle(uint256 assetId) returns()
func (_RWAValuation *RWAValuationSession) UpdateFromOracle(assetId *big.Int) (*types.Transaction, error) {
	return _RWAValuation.Contract.UpdateFromOracle(&_RWAValuation.TransactOpts, assetId)
}

// UpdateFromOracle is a paid mutator transaction binding the contract method 0x815d89eb.
//
// Solidity: function updateFromOracle(uint256 assetId) returns()
func (_RWAValuation *RWAValuationTransactorSession) UpdateFromOracle(assetId *big.Int) (*types.Transaction, error) {
	return _RWAValuation.Contract.UpdateFromOracle(&_RWAValuation.TransactOpts, assetId)
}

// UpdateValuation is a paid mutator transaction binding the contract method 0xc685a228.
//
// Solidity: function updateValuation(uint256 assetId, uint256 newValue, uint8 method, string reportHash, uint256 confidence) returns()
func (_RWAValuation *RWAValuationTransactor) UpdateValuation(opts *bind.TransactOpts, assetId *big.Int, newValue *big.Int, method uint8, reportHash string, confidence *big.Int) (*types.Transaction, error) {
	return _RWAValuation.contract.Transact(opts, "updateValuation", assetId, newValue, method, reportHash, confidence)
}

// UpdateValuation is a paid mutator transaction binding the contract method 0xc685a228.
//
// Solidity: function updateValuation(uint256 assetId, uint256 newValue, uint8 method, string reportHash, uint256 confidence) returns()
func (_RWAValuation *RWAValuationSession) UpdateValuation(assetId *big.Int, newValue *big.Int, method uint8, reportHash string, confidence *big.Int) (*types.Transaction, error) {
	return _RWAValuation.Contract.UpdateValuation(&_RWAValuation.TransactOpts, assetId, newValue, method, reportHash, confidence)
}

// UpdateValuation is a paid mutator transaction binding the contract method 0xc685a228.
//
// Solidity: function updateValuation(uint256 assetId, uint256 newValue, uint8 method, string reportHash, uint256 confidence) returns()
func (_RWAValuation *RWAValuationTransactorSession) UpdateValuation(assetId *big.Int, newValue *big.Int, method uint8, reportHash string, confidence *big.Int) (*types.Transaction, error) {
	return _RWAValuation.Contract.UpdateValuation(&_RWAValuation.TransactOpts, assetId, newValue, method, reportHash, confidence)
}

// RWAValuationOracleConfiguredIterator is returned from FilterOracleConfigured and is used to iterate over the raw logs and unpacked data for OracleConfigured events raised by the RWAValuation contract.
type RWAValuationOracleConfiguredIterator struct {
	Event *RWAValuationOracleConfigured // Event containing the contract specifics and raw log

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
func (it *RWAValuationOracleConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAValuationOracleConfigured)
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
		it.Event = new(RWAValuationOracleConfigured)
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
func (it *RWAValuationOracleConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAValuationOracleConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAValuationOracleConfigured represents a OracleConfigured event raised by the RWAValuation contract.
type RWAValuationOracleConfigured struct {
	AssetId   *big.Int
	Oracle    common.Address
	Heartbeat *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOracleConfigured is a free log retrieval operation binding the contract event 0x64ba9e72c40f67c528def016e092f1b8aab971f13f117c660c841cefd14663ba.
//
// Solidity: event OracleConfigured(uint256 indexed assetId, address indexed oracle, uint256 heartbeat)
func (_RWAValuation *RWAValuationFilterer) FilterOracleConfigured(opts *bind.FilterOpts, assetId []*big.Int, oracle []common.Address) (*RWAValuationOracleConfiguredIterator, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _RWAValuation.contract.FilterLogs(opts, "OracleConfigured", assetIdRule, oracleRule)
	if err != nil {
		return nil, err
	}
	return &RWAValuationOracleConfiguredIterator{contract: _RWAValuation.contract, event: "OracleConfigured", logs: logs, sub: sub}, nil
}

// WatchOracleConfigured is a free log subscription operation binding the contract event 0x64ba9e72c40f67c528def016e092f1b8aab971f13f117c660c841cefd14663ba.
//
// Solidity: event OracleConfigured(uint256 indexed assetId, address indexed oracle, uint256 heartbeat)
func (_RWAValuation *RWAValuationFilterer) WatchOracleConfigured(opts *bind.WatchOpts, sink chan<- *RWAValuationOracleConfigured, assetId []*big.Int, oracle []common.Address) (event.Subscription, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _RWAValuation.contract.WatchLogs(opts, "OracleConfigured", assetIdRule, oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAValuationOracleConfigured)
				if err := _RWAValuation.contract.UnpackLog(event, "OracleConfigured", log); err != nil {
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

// ParseOracleConfigured is a log parse operation binding the contract event 0x64ba9e72c40f67c528def016e092f1b8aab971f13f117c660c841cefd14663ba.
//
// Solidity: event OracleConfigured(uint256 indexed assetId, address indexed oracle, uint256 heartbeat)
func (_RWAValuation *RWAValuationFilterer) ParseOracleConfigured(log types.Log) (*RWAValuationOracleConfigured, error) {
	event := new(RWAValuationOracleConfigured)
	if err := _RWAValuation.contract.UnpackLog(event, "OracleConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAValuationRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the RWAValuation contract.
type RWAValuationRoleAdminChangedIterator struct {
	Event *RWAValuationRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *RWAValuationRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAValuationRoleAdminChanged)
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
		it.Event = new(RWAValuationRoleAdminChanged)
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
func (it *RWAValuationRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAValuationRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAValuationRoleAdminChanged represents a RoleAdminChanged event raised by the RWAValuation contract.
type RWAValuationRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RWAValuation *RWAValuationFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*RWAValuationRoleAdminChangedIterator, error) {

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

	logs, sub, err := _RWAValuation.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &RWAValuationRoleAdminChangedIterator{contract: _RWAValuation.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RWAValuation *RWAValuationFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *RWAValuationRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _RWAValuation.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAValuationRoleAdminChanged)
				if err := _RWAValuation.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_RWAValuation *RWAValuationFilterer) ParseRoleAdminChanged(log types.Log) (*RWAValuationRoleAdminChanged, error) {
	event := new(RWAValuationRoleAdminChanged)
	if err := _RWAValuation.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAValuationRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the RWAValuation contract.
type RWAValuationRoleGrantedIterator struct {
	Event *RWAValuationRoleGranted // Event containing the contract specifics and raw log

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
func (it *RWAValuationRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAValuationRoleGranted)
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
		it.Event = new(RWAValuationRoleGranted)
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
func (it *RWAValuationRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAValuationRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAValuationRoleGranted represents a RoleGranted event raised by the RWAValuation contract.
type RWAValuationRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAValuation *RWAValuationFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RWAValuationRoleGrantedIterator, error) {

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

	logs, sub, err := _RWAValuation.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RWAValuationRoleGrantedIterator{contract: _RWAValuation.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAValuation *RWAValuationFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *RWAValuationRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RWAValuation.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAValuationRoleGranted)
				if err := _RWAValuation.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_RWAValuation *RWAValuationFilterer) ParseRoleGranted(log types.Log) (*RWAValuationRoleGranted, error) {
	event := new(RWAValuationRoleGranted)
	if err := _RWAValuation.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAValuationRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the RWAValuation contract.
type RWAValuationRoleRevokedIterator struct {
	Event *RWAValuationRoleRevoked // Event containing the contract specifics and raw log

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
func (it *RWAValuationRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAValuationRoleRevoked)
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
		it.Event = new(RWAValuationRoleRevoked)
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
func (it *RWAValuationRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAValuationRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAValuationRoleRevoked represents a RoleRevoked event raised by the RWAValuation contract.
type RWAValuationRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAValuation *RWAValuationFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RWAValuationRoleRevokedIterator, error) {

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

	logs, sub, err := _RWAValuation.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RWAValuationRoleRevokedIterator{contract: _RWAValuation.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAValuation *RWAValuationFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *RWAValuationRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RWAValuation.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAValuationRoleRevoked)
				if err := _RWAValuation.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_RWAValuation *RWAValuationFilterer) ParseRoleRevoked(log types.Log) (*RWAValuationRoleRevoked, error) {
	event := new(RWAValuationRoleRevoked)
	if err := _RWAValuation.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAValuationValuationDisputedIterator is returned from FilterValuationDisputed and is used to iterate over the raw logs and unpacked data for ValuationDisputed events raised by the RWAValuation contract.
type RWAValuationValuationDisputedIterator struct {
	Event *RWAValuationValuationDisputed // Event containing the contract specifics and raw log

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
func (it *RWAValuationValuationDisputedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAValuationValuationDisputed)
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
		it.Event = new(RWAValuationValuationDisputed)
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
func (it *RWAValuationValuationDisputedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAValuationValuationDisputedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAValuationValuationDisputed represents a ValuationDisputed event raised by the RWAValuation contract.
type RWAValuationValuationDisputed struct {
	AssetId       *big.Int
	DisputedValue *big.Int
	Disputer      common.Address
	Reason        string
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterValuationDisputed is a free log retrieval operation binding the contract event 0xfff46f6fa1de49158e4d2f75b5bcbf5dc2a41965a143cc98f1e2399a0d40b0a1.
//
// Solidity: event ValuationDisputed(uint256 indexed assetId, uint256 disputedValue, address indexed disputer, string reason)
func (_RWAValuation *RWAValuationFilterer) FilterValuationDisputed(opts *bind.FilterOpts, assetId []*big.Int, disputer []common.Address) (*RWAValuationValuationDisputedIterator, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	var disputerRule []interface{}
	for _, disputerItem := range disputer {
		disputerRule = append(disputerRule, disputerItem)
	}

	logs, sub, err := _RWAValuation.contract.FilterLogs(opts, "ValuationDisputed", assetIdRule, disputerRule)
	if err != nil {
		return nil, err
	}
	return &RWAValuationValuationDisputedIterator{contract: _RWAValuation.contract, event: "ValuationDisputed", logs: logs, sub: sub}, nil
}

// WatchValuationDisputed is a free log subscription operation binding the contract event 0xfff46f6fa1de49158e4d2f75b5bcbf5dc2a41965a143cc98f1e2399a0d40b0a1.
//
// Solidity: event ValuationDisputed(uint256 indexed assetId, uint256 disputedValue, address indexed disputer, string reason)
func (_RWAValuation *RWAValuationFilterer) WatchValuationDisputed(opts *bind.WatchOpts, sink chan<- *RWAValuationValuationDisputed, assetId []*big.Int, disputer []common.Address) (event.Subscription, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	var disputerRule []interface{}
	for _, disputerItem := range disputer {
		disputerRule = append(disputerRule, disputerItem)
	}

	logs, sub, err := _RWAValuation.contract.WatchLogs(opts, "ValuationDisputed", assetIdRule, disputerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAValuationValuationDisputed)
				if err := _RWAValuation.contract.UnpackLog(event, "ValuationDisputed", log); err != nil {
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

// ParseValuationDisputed is a log parse operation binding the contract event 0xfff46f6fa1de49158e4d2f75b5bcbf5dc2a41965a143cc98f1e2399a0d40b0a1.
//
// Solidity: event ValuationDisputed(uint256 indexed assetId, uint256 disputedValue, address indexed disputer, string reason)
func (_RWAValuation *RWAValuationFilterer) ParseValuationDisputed(log types.Log) (*RWAValuationValuationDisputed, error) {
	event := new(RWAValuationValuationDisputed)
	if err := _RWAValuation.contract.UnpackLog(event, "ValuationDisputed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAValuationValuationUpdatedIterator is returned from FilterValuationUpdated and is used to iterate over the raw logs and unpacked data for ValuationUpdated events raised by the RWAValuation contract.
type RWAValuationValuationUpdatedIterator struct {
	Event *RWAValuationValuationUpdated // Event containing the contract specifics and raw log

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
func (it *RWAValuationValuationUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAValuationValuationUpdated)
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
		it.Event = new(RWAValuationValuationUpdated)
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
func (it *RWAValuationValuationUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAValuationValuationUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAValuationValuationUpdated represents a ValuationUpdated event raised by the RWAValuation contract.
type RWAValuationValuationUpdated struct {
	AssetId  *big.Int
	OldValue *big.Int
	NewValue *big.Int
	Method   uint8
	Valuator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterValuationUpdated is a free log retrieval operation binding the contract event 0x358f42fffcc0ebe94fe24a39bdb7991745d494092cbf78077150d819f2c6b6e6.
//
// Solidity: event ValuationUpdated(uint256 indexed assetId, uint256 oldValue, uint256 newValue, uint8 method, address indexed valuator)
func (_RWAValuation *RWAValuationFilterer) FilterValuationUpdated(opts *bind.FilterOpts, assetId []*big.Int, valuator []common.Address) (*RWAValuationValuationUpdatedIterator, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	var valuatorRule []interface{}
	for _, valuatorItem := range valuator {
		valuatorRule = append(valuatorRule, valuatorItem)
	}

	logs, sub, err := _RWAValuation.contract.FilterLogs(opts, "ValuationUpdated", assetIdRule, valuatorRule)
	if err != nil {
		return nil, err
	}
	return &RWAValuationValuationUpdatedIterator{contract: _RWAValuation.contract, event: "ValuationUpdated", logs: logs, sub: sub}, nil
}

// WatchValuationUpdated is a free log subscription operation binding the contract event 0x358f42fffcc0ebe94fe24a39bdb7991745d494092cbf78077150d819f2c6b6e6.
//
// Solidity: event ValuationUpdated(uint256 indexed assetId, uint256 oldValue, uint256 newValue, uint8 method, address indexed valuator)
func (_RWAValuation *RWAValuationFilterer) WatchValuationUpdated(opts *bind.WatchOpts, sink chan<- *RWAValuationValuationUpdated, assetId []*big.Int, valuator []common.Address) (event.Subscription, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	var valuatorRule []interface{}
	for _, valuatorItem := range valuator {
		valuatorRule = append(valuatorRule, valuatorItem)
	}

	logs, sub, err := _RWAValuation.contract.WatchLogs(opts, "ValuationUpdated", assetIdRule, valuatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAValuationValuationUpdated)
				if err := _RWAValuation.contract.UnpackLog(event, "ValuationUpdated", log); err != nil {
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

// ParseValuationUpdated is a log parse operation binding the contract event 0x358f42fffcc0ebe94fe24a39bdb7991745d494092cbf78077150d819f2c6b6e6.
//
// Solidity: event ValuationUpdated(uint256 indexed assetId, uint256 oldValue, uint256 newValue, uint8 method, address indexed valuator)
func (_RWAValuation *RWAValuationFilterer) ParseValuationUpdated(log types.Log) (*RWAValuationValuationUpdated, error) {
	event := new(RWAValuationValuationUpdated)
	if err := _RWAValuation.contract.UnpackLog(event, "ValuationUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAValuationValuatorAuthorizedIterator is returned from FilterValuatorAuthorized and is used to iterate over the raw logs and unpacked data for ValuatorAuthorized events raised by the RWAValuation contract.
type RWAValuationValuatorAuthorizedIterator struct {
	Event *RWAValuationValuatorAuthorized // Event containing the contract specifics and raw log

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
func (it *RWAValuationValuatorAuthorizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAValuationValuatorAuthorized)
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
		it.Event = new(RWAValuationValuatorAuthorized)
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
func (it *RWAValuationValuatorAuthorizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAValuationValuatorAuthorizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAValuationValuatorAuthorized represents a ValuatorAuthorized event raised by the RWAValuation contract.
type RWAValuationValuatorAuthorized struct {
	Valuator   common.Address
	Authorized bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterValuatorAuthorized is a free log retrieval operation binding the contract event 0x8834f82d05d9a65903a0ef9ca145fd65cec65f2f4451a67ff1d1d090181ecf86.
//
// Solidity: event ValuatorAuthorized(address indexed valuator, bool authorized)
func (_RWAValuation *RWAValuationFilterer) FilterValuatorAuthorized(opts *bind.FilterOpts, valuator []common.Address) (*RWAValuationValuatorAuthorizedIterator, error) {

	var valuatorRule []interface{}
	for _, valuatorItem := range valuator {
		valuatorRule = append(valuatorRule, valuatorItem)
	}

	logs, sub, err := _RWAValuation.contract.FilterLogs(opts, "ValuatorAuthorized", valuatorRule)
	if err != nil {
		return nil, err
	}
	return &RWAValuationValuatorAuthorizedIterator{contract: _RWAValuation.contract, event: "ValuatorAuthorized", logs: logs, sub: sub}, nil
}

// WatchValuatorAuthorized is a free log subscription operation binding the contract event 0x8834f82d05d9a65903a0ef9ca145fd65cec65f2f4451a67ff1d1d090181ecf86.
//
// Solidity: event ValuatorAuthorized(address indexed valuator, bool authorized)
func (_RWAValuation *RWAValuationFilterer) WatchValuatorAuthorized(opts *bind.WatchOpts, sink chan<- *RWAValuationValuatorAuthorized, valuator []common.Address) (event.Subscription, error) {

	var valuatorRule []interface{}
	for _, valuatorItem := range valuator {
		valuatorRule = append(valuatorRule, valuatorItem)
	}

	logs, sub, err := _RWAValuation.contract.WatchLogs(opts, "ValuatorAuthorized", valuatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAValuationValuatorAuthorized)
				if err := _RWAValuation.contract.UnpackLog(event, "ValuatorAuthorized", log); err != nil {
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

// ParseValuatorAuthorized is a log parse operation binding the contract event 0x8834f82d05d9a65903a0ef9ca145fd65cec65f2f4451a67ff1d1d090181ecf86.
//
// Solidity: event ValuatorAuthorized(address indexed valuator, bool authorized)
func (_RWAValuation *RWAValuationFilterer) ParseValuatorAuthorized(log types.Log) (*RWAValuationValuatorAuthorized, error) {
	event := new(RWAValuationValuatorAuthorized)
	if err := _RWAValuation.contract.UnpackLog(event, "ValuatorAuthorized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
