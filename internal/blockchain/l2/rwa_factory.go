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

// IRWAAssetAssetMetadata is an auto generated low-level Go binding around an user-defined struct.
type IRWAAssetAssetMetadata struct {
	Name                string
	Symbol              string
	AssetType           uint8
	TotalValue          *big.Int
	TotalSupply         *big.Int
	Issuer              common.Address
	IssuanceDate        *big.Int
	MaturityDate        *big.Int
	LegalDocumentHash   string
	ValuationReportHash string
}

// IRWAAssetYieldTerms is an auto generated low-level Go binding around an user-defined struct.
type IRWAAssetYieldTerms struct {
	AnnualYieldRate       *big.Int
	YieldPaymentFrequency *big.Int
	LastYieldPayment      *big.Int
	TotalYieldPaid        *big.Int
}

// RWAAssetFactoryMetaData contains all meta data concerning the RWAAssetFactory contract.
var RWAAssetFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"compliance\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"valuation\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumIRWAAsset.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalValue\",\"type\":\"uint256\"}],\"name\":\"AssetCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"fractionalToken\",\"type\":\"address\"}],\"name\":\"AssetFractionalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumIRWAAsset.AssetStatus\",\"name\":\"oldStatus\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumIRWAAsset.AssetStatus\",\"name\":\"newStatus\",\"type\":\"uint8\"}],\"name\":\"AssetStatusChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"AssetValuationUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"YieldDistributed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ASSET_MANAGER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"COMPLIANCE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"activateAsset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activeAssets\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"complianceContract\",\"outputs\":[{\"internalType\":\"contractIRWACompliance\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"enumIRWAAsset.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"totalValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maturityDate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"legalDocumentHash\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valuationReportHash\",\"type\":\"string\"}],\"name\":\"createAsset\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"distributeYield\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"}],\"name\":\"fractionalizeAsset\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveAssets\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAssets\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getAssetMetadata\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"enumIRWAAsset.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"totalValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"issuanceDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maturityDate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"legalDocumentHash\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"valuationReportHash\",\"type\":\"string\"}],\"internalType\":\"structIRWAAsset.AssetMetadata\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getAssetStatus\",\"outputs\":[{\"internalType\":\"enumIRWAAsset.AssetStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getAssetValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIRWAAsset.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"}],\"name\":\"getAssetsByType\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getFractionalToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"}],\"name\":\"getIssuerAssets\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getYieldTerms\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"annualYieldRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"yieldPaymentFrequency\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastYieldPayment\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalYieldPaid\",\"type\":\"uint256\"}],\"internalType\":\"structIRWAAsset.YieldTerms\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"isAssetActive\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"annualYieldRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"yieldPaymentFrequency\",\"type\":\"uint256\"}],\"name\":\"setYieldTerms\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalAssets\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalValueLocked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"enumIRWAAsset.AssetStatus\",\"name\":\"newStatus\",\"type\":\"uint8\"}],\"name\":\"updateAssetStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"valuationReportHash\",\"type\":\"string\"}],\"name\":\"updateAssetValuation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"valuationContract\",\"outputs\":[{\"internalType\":\"contractIRWAValuation\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// RWAAssetFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use RWAAssetFactoryMetaData.ABI instead.
var RWAAssetFactoryABI = RWAAssetFactoryMetaData.ABI

// RWAAssetFactory is an auto generated Go binding around an Ethereum contract.
type RWAAssetFactory struct {
	RWAAssetFactoryCaller     // Read-only binding to the contract
	RWAAssetFactoryTransactor // Write-only binding to the contract
	RWAAssetFactoryFilterer   // Log filterer for contract events
}

// RWAAssetFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RWAAssetFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAAssetFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RWAAssetFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAAssetFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RWAAssetFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAAssetFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RWAAssetFactorySession struct {
	Contract     *RWAAssetFactory  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RWAAssetFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RWAAssetFactoryCallerSession struct {
	Contract *RWAAssetFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// RWAAssetFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RWAAssetFactoryTransactorSession struct {
	Contract     *RWAAssetFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// RWAAssetFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RWAAssetFactoryRaw struct {
	Contract *RWAAssetFactory // Generic contract binding to access the raw methods on
}

// RWAAssetFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RWAAssetFactoryCallerRaw struct {
	Contract *RWAAssetFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// RWAAssetFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RWAAssetFactoryTransactorRaw struct {
	Contract *RWAAssetFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRWAAssetFactory creates a new instance of RWAAssetFactory, bound to a specific deployed contract.
func NewRWAAssetFactory(address common.Address, backend bind.ContractBackend) (*RWAAssetFactory, error) {
	contract, err := bindRWAAssetFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactory{RWAAssetFactoryCaller: RWAAssetFactoryCaller{contract: contract}, RWAAssetFactoryTransactor: RWAAssetFactoryTransactor{contract: contract}, RWAAssetFactoryFilterer: RWAAssetFactoryFilterer{contract: contract}}, nil
}

// NewRWAAssetFactoryCaller creates a new read-only instance of RWAAssetFactory, bound to a specific deployed contract.
func NewRWAAssetFactoryCaller(address common.Address, caller bind.ContractCaller) (*RWAAssetFactoryCaller, error) {
	contract, err := bindRWAAssetFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryCaller{contract: contract}, nil
}

// NewRWAAssetFactoryTransactor creates a new write-only instance of RWAAssetFactory, bound to a specific deployed contract.
func NewRWAAssetFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*RWAAssetFactoryTransactor, error) {
	contract, err := bindRWAAssetFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryTransactor{contract: contract}, nil
}

// NewRWAAssetFactoryFilterer creates a new log filterer instance of RWAAssetFactory, bound to a specific deployed contract.
func NewRWAAssetFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*RWAAssetFactoryFilterer, error) {
	contract, err := bindRWAAssetFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryFilterer{contract: contract}, nil
}

// bindRWAAssetFactory binds a generic wrapper to an already deployed contract.
func bindRWAAssetFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RWAAssetFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RWAAssetFactory *RWAAssetFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RWAAssetFactory.Contract.RWAAssetFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RWAAssetFactory *RWAAssetFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.RWAAssetFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RWAAssetFactory *RWAAssetFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.RWAAssetFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RWAAssetFactory *RWAAssetFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RWAAssetFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RWAAssetFactory *RWAAssetFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RWAAssetFactory *RWAAssetFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.contract.Transact(opts, method, params...)
}

// ASSETMANAGERROLE is a free data retrieval call binding the contract method 0xa4b32de8.
//
// Solidity: function ASSET_MANAGER_ROLE() view returns(bytes32)
func (_RWAAssetFactory *RWAAssetFactoryCaller) ASSETMANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "ASSET_MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ASSETMANAGERROLE is a free data retrieval call binding the contract method 0xa4b32de8.
//
// Solidity: function ASSET_MANAGER_ROLE() view returns(bytes32)
func (_RWAAssetFactory *RWAAssetFactorySession) ASSETMANAGERROLE() ([32]byte, error) {
	return _RWAAssetFactory.Contract.ASSETMANAGERROLE(&_RWAAssetFactory.CallOpts)
}

// ASSETMANAGERROLE is a free data retrieval call binding the contract method 0xa4b32de8.
//
// Solidity: function ASSET_MANAGER_ROLE() view returns(bytes32)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) ASSETMANAGERROLE() ([32]byte, error) {
	return _RWAAssetFactory.Contract.ASSETMANAGERROLE(&_RWAAssetFactory.CallOpts)
}

// COMPLIANCEROLE is a free data retrieval call binding the contract method 0x062d3bd7.
//
// Solidity: function COMPLIANCE_ROLE() view returns(bytes32)
func (_RWAAssetFactory *RWAAssetFactoryCaller) COMPLIANCEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "COMPLIANCE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// COMPLIANCEROLE is a free data retrieval call binding the contract method 0x062d3bd7.
//
// Solidity: function COMPLIANCE_ROLE() view returns(bytes32)
func (_RWAAssetFactory *RWAAssetFactorySession) COMPLIANCEROLE() ([32]byte, error) {
	return _RWAAssetFactory.Contract.COMPLIANCEROLE(&_RWAAssetFactory.CallOpts)
}

// COMPLIANCEROLE is a free data retrieval call binding the contract method 0x062d3bd7.
//
// Solidity: function COMPLIANCE_ROLE() view returns(bytes32)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) COMPLIANCEROLE() ([32]byte, error) {
	return _RWAAssetFactory.Contract.COMPLIANCEROLE(&_RWAAssetFactory.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAAssetFactory *RWAAssetFactoryCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAAssetFactory *RWAAssetFactorySession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RWAAssetFactory.Contract.DEFAULTADMINROLE(&_RWAAssetFactory.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RWAAssetFactory.Contract.DEFAULTADMINROLE(&_RWAAssetFactory.CallOpts)
}

// ActiveAssets is a free data retrieval call binding the contract method 0x1c17b946.
//
// Solidity: function activeAssets() view returns(uint256)
func (_RWAAssetFactory *RWAAssetFactoryCaller) ActiveAssets(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "activeAssets")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveAssets is a free data retrieval call binding the contract method 0x1c17b946.
//
// Solidity: function activeAssets() view returns(uint256)
func (_RWAAssetFactory *RWAAssetFactorySession) ActiveAssets() (*big.Int, error) {
	return _RWAAssetFactory.Contract.ActiveAssets(&_RWAAssetFactory.CallOpts)
}

// ActiveAssets is a free data retrieval call binding the contract method 0x1c17b946.
//
// Solidity: function activeAssets() view returns(uint256)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) ActiveAssets() (*big.Int, error) {
	return _RWAAssetFactory.Contract.ActiveAssets(&_RWAAssetFactory.CallOpts)
}

// ComplianceContract is a free data retrieval call binding the contract method 0xb2a2a4e2.
//
// Solidity: function complianceContract() view returns(address)
func (_RWAAssetFactory *RWAAssetFactoryCaller) ComplianceContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "complianceContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ComplianceContract is a free data retrieval call binding the contract method 0xb2a2a4e2.
//
// Solidity: function complianceContract() view returns(address)
func (_RWAAssetFactory *RWAAssetFactorySession) ComplianceContract() (common.Address, error) {
	return _RWAAssetFactory.Contract.ComplianceContract(&_RWAAssetFactory.CallOpts)
}

// ComplianceContract is a free data retrieval call binding the contract method 0xb2a2a4e2.
//
// Solidity: function complianceContract() view returns(address)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) ComplianceContract() (common.Address, error) {
	return _RWAAssetFactory.Contract.ComplianceContract(&_RWAAssetFactory.CallOpts)
}

// GetActiveAssets is a free data retrieval call binding the contract method 0x0fd8aee8.
//
// Solidity: function getActiveAssets() view returns(uint256[])
func (_RWAAssetFactory *RWAAssetFactoryCaller) GetActiveAssets(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "getActiveAssets")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveAssets is a free data retrieval call binding the contract method 0x0fd8aee8.
//
// Solidity: function getActiveAssets() view returns(uint256[])
func (_RWAAssetFactory *RWAAssetFactorySession) GetActiveAssets() ([]*big.Int, error) {
	return _RWAAssetFactory.Contract.GetActiveAssets(&_RWAAssetFactory.CallOpts)
}

// GetActiveAssets is a free data retrieval call binding the contract method 0x0fd8aee8.
//
// Solidity: function getActiveAssets() view returns(uint256[])
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) GetActiveAssets() ([]*big.Int, error) {
	return _RWAAssetFactory.Contract.GetActiveAssets(&_RWAAssetFactory.CallOpts)
}

// GetAllAssets is a free data retrieval call binding the contract method 0x2acada4d.
//
// Solidity: function getAllAssets() view returns(uint256[])
func (_RWAAssetFactory *RWAAssetFactoryCaller) GetAllAssets(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "getAllAssets")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAllAssets is a free data retrieval call binding the contract method 0x2acada4d.
//
// Solidity: function getAllAssets() view returns(uint256[])
func (_RWAAssetFactory *RWAAssetFactorySession) GetAllAssets() ([]*big.Int, error) {
	return _RWAAssetFactory.Contract.GetAllAssets(&_RWAAssetFactory.CallOpts)
}

// GetAllAssets is a free data retrieval call binding the contract method 0x2acada4d.
//
// Solidity: function getAllAssets() view returns(uint256[])
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) GetAllAssets() ([]*big.Int, error) {
	return _RWAAssetFactory.Contract.GetAllAssets(&_RWAAssetFactory.CallOpts)
}

// GetAssetMetadata is a free data retrieval call binding the contract method 0xcf2acf84.
//
// Solidity: function getAssetMetadata(uint256 assetId) view returns((string,string,uint8,uint256,uint256,address,uint256,uint256,string,string))
func (_RWAAssetFactory *RWAAssetFactoryCaller) GetAssetMetadata(opts *bind.CallOpts, assetId *big.Int) (IRWAAssetAssetMetadata, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "getAssetMetadata", assetId)

	if err != nil {
		return *new(IRWAAssetAssetMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new(IRWAAssetAssetMetadata)).(*IRWAAssetAssetMetadata)

	return out0, err

}

// GetAssetMetadata is a free data retrieval call binding the contract method 0xcf2acf84.
//
// Solidity: function getAssetMetadata(uint256 assetId) view returns((string,string,uint8,uint256,uint256,address,uint256,uint256,string,string))
func (_RWAAssetFactory *RWAAssetFactorySession) GetAssetMetadata(assetId *big.Int) (IRWAAssetAssetMetadata, error) {
	return _RWAAssetFactory.Contract.GetAssetMetadata(&_RWAAssetFactory.CallOpts, assetId)
}

// GetAssetMetadata is a free data retrieval call binding the contract method 0xcf2acf84.
//
// Solidity: function getAssetMetadata(uint256 assetId) view returns((string,string,uint8,uint256,uint256,address,uint256,uint256,string,string))
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) GetAssetMetadata(assetId *big.Int) (IRWAAssetAssetMetadata, error) {
	return _RWAAssetFactory.Contract.GetAssetMetadata(&_RWAAssetFactory.CallOpts, assetId)
}

// GetAssetStatus is a free data retrieval call binding the contract method 0x823a65f5.
//
// Solidity: function getAssetStatus(uint256 assetId) view returns(uint8)
func (_RWAAssetFactory *RWAAssetFactoryCaller) GetAssetStatus(opts *bind.CallOpts, assetId *big.Int) (uint8, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "getAssetStatus", assetId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetAssetStatus is a free data retrieval call binding the contract method 0x823a65f5.
//
// Solidity: function getAssetStatus(uint256 assetId) view returns(uint8)
func (_RWAAssetFactory *RWAAssetFactorySession) GetAssetStatus(assetId *big.Int) (uint8, error) {
	return _RWAAssetFactory.Contract.GetAssetStatus(&_RWAAssetFactory.CallOpts, assetId)
}

// GetAssetStatus is a free data retrieval call binding the contract method 0x823a65f5.
//
// Solidity: function getAssetStatus(uint256 assetId) view returns(uint8)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) GetAssetStatus(assetId *big.Int) (uint8, error) {
	return _RWAAssetFactory.Contract.GetAssetStatus(&_RWAAssetFactory.CallOpts, assetId)
}

// GetAssetValue is a free data retrieval call binding the contract method 0x05c12a46.
//
// Solidity: function getAssetValue(uint256 assetId) view returns(uint256)
func (_RWAAssetFactory *RWAAssetFactoryCaller) GetAssetValue(opts *bind.CallOpts, assetId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "getAssetValue", assetId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAssetValue is a free data retrieval call binding the contract method 0x05c12a46.
//
// Solidity: function getAssetValue(uint256 assetId) view returns(uint256)
func (_RWAAssetFactory *RWAAssetFactorySession) GetAssetValue(assetId *big.Int) (*big.Int, error) {
	return _RWAAssetFactory.Contract.GetAssetValue(&_RWAAssetFactory.CallOpts, assetId)
}

// GetAssetValue is a free data retrieval call binding the contract method 0x05c12a46.
//
// Solidity: function getAssetValue(uint256 assetId) view returns(uint256)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) GetAssetValue(assetId *big.Int) (*big.Int, error) {
	return _RWAAssetFactory.Contract.GetAssetValue(&_RWAAssetFactory.CallOpts, assetId)
}

// GetAssetsByType is a free data retrieval call binding the contract method 0xd1b2cc2e.
//
// Solidity: function getAssetsByType(uint8 assetType) view returns(uint256[])
func (_RWAAssetFactory *RWAAssetFactoryCaller) GetAssetsByType(opts *bind.CallOpts, assetType uint8) ([]*big.Int, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "getAssetsByType", assetType)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAssetsByType is a free data retrieval call binding the contract method 0xd1b2cc2e.
//
// Solidity: function getAssetsByType(uint8 assetType) view returns(uint256[])
func (_RWAAssetFactory *RWAAssetFactorySession) GetAssetsByType(assetType uint8) ([]*big.Int, error) {
	return _RWAAssetFactory.Contract.GetAssetsByType(&_RWAAssetFactory.CallOpts, assetType)
}

// GetAssetsByType is a free data retrieval call binding the contract method 0xd1b2cc2e.
//
// Solidity: function getAssetsByType(uint8 assetType) view returns(uint256[])
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) GetAssetsByType(assetType uint8) ([]*big.Int, error) {
	return _RWAAssetFactory.Contract.GetAssetsByType(&_RWAAssetFactory.CallOpts, assetType)
}

// GetFractionalToken is a free data retrieval call binding the contract method 0xf306ee05.
//
// Solidity: function getFractionalToken(uint256 assetId) view returns(address)
func (_RWAAssetFactory *RWAAssetFactoryCaller) GetFractionalToken(opts *bind.CallOpts, assetId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "getFractionalToken", assetId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFractionalToken is a free data retrieval call binding the contract method 0xf306ee05.
//
// Solidity: function getFractionalToken(uint256 assetId) view returns(address)
func (_RWAAssetFactory *RWAAssetFactorySession) GetFractionalToken(assetId *big.Int) (common.Address, error) {
	return _RWAAssetFactory.Contract.GetFractionalToken(&_RWAAssetFactory.CallOpts, assetId)
}

// GetFractionalToken is a free data retrieval call binding the contract method 0xf306ee05.
//
// Solidity: function getFractionalToken(uint256 assetId) view returns(address)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) GetFractionalToken(assetId *big.Int) (common.Address, error) {
	return _RWAAssetFactory.Contract.GetFractionalToken(&_RWAAssetFactory.CallOpts, assetId)
}

// GetIssuerAssets is a free data retrieval call binding the contract method 0x6429d103.
//
// Solidity: function getIssuerAssets(address issuer) view returns(uint256[])
func (_RWAAssetFactory *RWAAssetFactoryCaller) GetIssuerAssets(opts *bind.CallOpts, issuer common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "getIssuerAssets", issuer)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetIssuerAssets is a free data retrieval call binding the contract method 0x6429d103.
//
// Solidity: function getIssuerAssets(address issuer) view returns(uint256[])
func (_RWAAssetFactory *RWAAssetFactorySession) GetIssuerAssets(issuer common.Address) ([]*big.Int, error) {
	return _RWAAssetFactory.Contract.GetIssuerAssets(&_RWAAssetFactory.CallOpts, issuer)
}

// GetIssuerAssets is a free data retrieval call binding the contract method 0x6429d103.
//
// Solidity: function getIssuerAssets(address issuer) view returns(uint256[])
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) GetIssuerAssets(issuer common.Address) ([]*big.Int, error) {
	return _RWAAssetFactory.Contract.GetIssuerAssets(&_RWAAssetFactory.CallOpts, issuer)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAAssetFactory *RWAAssetFactoryCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAAssetFactory *RWAAssetFactorySession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RWAAssetFactory.Contract.GetRoleAdmin(&_RWAAssetFactory.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RWAAssetFactory.Contract.GetRoleAdmin(&_RWAAssetFactory.CallOpts, role)
}

// GetYieldTerms is a free data retrieval call binding the contract method 0x7212520c.
//
// Solidity: function getYieldTerms(uint256 assetId) view returns((uint256,uint256,uint256,uint256))
func (_RWAAssetFactory *RWAAssetFactoryCaller) GetYieldTerms(opts *bind.CallOpts, assetId *big.Int) (IRWAAssetYieldTerms, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "getYieldTerms", assetId)

	if err != nil {
		return *new(IRWAAssetYieldTerms), err
	}

	out0 := *abi.ConvertType(out[0], new(IRWAAssetYieldTerms)).(*IRWAAssetYieldTerms)

	return out0, err

}

// GetYieldTerms is a free data retrieval call binding the contract method 0x7212520c.
//
// Solidity: function getYieldTerms(uint256 assetId) view returns((uint256,uint256,uint256,uint256))
func (_RWAAssetFactory *RWAAssetFactorySession) GetYieldTerms(assetId *big.Int) (IRWAAssetYieldTerms, error) {
	return _RWAAssetFactory.Contract.GetYieldTerms(&_RWAAssetFactory.CallOpts, assetId)
}

// GetYieldTerms is a free data retrieval call binding the contract method 0x7212520c.
//
// Solidity: function getYieldTerms(uint256 assetId) view returns((uint256,uint256,uint256,uint256))
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) GetYieldTerms(assetId *big.Int) (IRWAAssetYieldTerms, error) {
	return _RWAAssetFactory.Contract.GetYieldTerms(&_RWAAssetFactory.CallOpts, assetId)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAAssetFactory *RWAAssetFactoryCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAAssetFactory *RWAAssetFactorySession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RWAAssetFactory.Contract.HasRole(&_RWAAssetFactory.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RWAAssetFactory.Contract.HasRole(&_RWAAssetFactory.CallOpts, role, account)
}

// IsAssetActive is a free data retrieval call binding the contract method 0x22e900c2.
//
// Solidity: function isAssetActive(uint256 assetId) view returns(bool)
func (_RWAAssetFactory *RWAAssetFactoryCaller) IsAssetActive(opts *bind.CallOpts, assetId *big.Int) (bool, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "isAssetActive", assetId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAssetActive is a free data retrieval call binding the contract method 0x22e900c2.
//
// Solidity: function isAssetActive(uint256 assetId) view returns(bool)
func (_RWAAssetFactory *RWAAssetFactorySession) IsAssetActive(assetId *big.Int) (bool, error) {
	return _RWAAssetFactory.Contract.IsAssetActive(&_RWAAssetFactory.CallOpts, assetId)
}

// IsAssetActive is a free data retrieval call binding the contract method 0x22e900c2.
//
// Solidity: function isAssetActive(uint256 assetId) view returns(bool)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) IsAssetActive(assetId *big.Int) (bool, error) {
	return _RWAAssetFactory.Contract.IsAssetActive(&_RWAAssetFactory.CallOpts, assetId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RWAAssetFactory *RWAAssetFactoryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RWAAssetFactory *RWAAssetFactorySession) Paused() (bool, error) {
	return _RWAAssetFactory.Contract.Paused(&_RWAAssetFactory.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) Paused() (bool, error) {
	return _RWAAssetFactory.Contract.Paused(&_RWAAssetFactory.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAAssetFactory *RWAAssetFactoryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAAssetFactory *RWAAssetFactorySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RWAAssetFactory.Contract.SupportsInterface(&_RWAAssetFactory.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RWAAssetFactory.Contract.SupportsInterface(&_RWAAssetFactory.CallOpts, interfaceId)
}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_RWAAssetFactory *RWAAssetFactoryCaller) TotalAssets(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "totalAssets")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_RWAAssetFactory *RWAAssetFactorySession) TotalAssets() (*big.Int, error) {
	return _RWAAssetFactory.Contract.TotalAssets(&_RWAAssetFactory.CallOpts)
}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) TotalAssets() (*big.Int, error) {
	return _RWAAssetFactory.Contract.TotalAssets(&_RWAAssetFactory.CallOpts)
}

// TotalValueLocked is a free data retrieval call binding the contract method 0xec18154e.
//
// Solidity: function totalValueLocked() view returns(uint256)
func (_RWAAssetFactory *RWAAssetFactoryCaller) TotalValueLocked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "totalValueLocked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalValueLocked is a free data retrieval call binding the contract method 0xec18154e.
//
// Solidity: function totalValueLocked() view returns(uint256)
func (_RWAAssetFactory *RWAAssetFactorySession) TotalValueLocked() (*big.Int, error) {
	return _RWAAssetFactory.Contract.TotalValueLocked(&_RWAAssetFactory.CallOpts)
}

// TotalValueLocked is a free data retrieval call binding the contract method 0xec18154e.
//
// Solidity: function totalValueLocked() view returns(uint256)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) TotalValueLocked() (*big.Int, error) {
	return _RWAAssetFactory.Contract.TotalValueLocked(&_RWAAssetFactory.CallOpts)
}

// ValuationContract is a free data retrieval call binding the contract method 0x274dc7ca.
//
// Solidity: function valuationContract() view returns(address)
func (_RWAAssetFactory *RWAAssetFactoryCaller) ValuationContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RWAAssetFactory.contract.Call(opts, &out, "valuationContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ValuationContract is a free data retrieval call binding the contract method 0x274dc7ca.
//
// Solidity: function valuationContract() view returns(address)
func (_RWAAssetFactory *RWAAssetFactorySession) ValuationContract() (common.Address, error) {
	return _RWAAssetFactory.Contract.ValuationContract(&_RWAAssetFactory.CallOpts)
}

// ValuationContract is a free data retrieval call binding the contract method 0x274dc7ca.
//
// Solidity: function valuationContract() view returns(address)
func (_RWAAssetFactory *RWAAssetFactoryCallerSession) ValuationContract() (common.Address, error) {
	return _RWAAssetFactory.Contract.ValuationContract(&_RWAAssetFactory.CallOpts)
}

// ActivateAsset is a paid mutator transaction binding the contract method 0xffba205f.
//
// Solidity: function activateAsset(uint256 assetId) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactor) ActivateAsset(opts *bind.TransactOpts, assetId *big.Int) (*types.Transaction, error) {
	return _RWAAssetFactory.contract.Transact(opts, "activateAsset", assetId)
}

// ActivateAsset is a paid mutator transaction binding the contract method 0xffba205f.
//
// Solidity: function activateAsset(uint256 assetId) returns()
func (_RWAAssetFactory *RWAAssetFactorySession) ActivateAsset(assetId *big.Int) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.ActivateAsset(&_RWAAssetFactory.TransactOpts, assetId)
}

// ActivateAsset is a paid mutator transaction binding the contract method 0xffba205f.
//
// Solidity: function activateAsset(uint256 assetId) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactorSession) ActivateAsset(assetId *big.Int) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.ActivateAsset(&_RWAAssetFactory.TransactOpts, assetId)
}

// CreateAsset is a paid mutator transaction binding the contract method 0xc183a1a6.
//
// Solidity: function createAsset(string name, string symbol, uint8 assetType, uint256 totalValue, uint256 maturityDate, string legalDocumentHash, string valuationReportHash) returns(uint256 assetId)
func (_RWAAssetFactory *RWAAssetFactoryTransactor) CreateAsset(opts *bind.TransactOpts, name string, symbol string, assetType uint8, totalValue *big.Int, maturityDate *big.Int, legalDocumentHash string, valuationReportHash string) (*types.Transaction, error) {
	return _RWAAssetFactory.contract.Transact(opts, "createAsset", name, symbol, assetType, totalValue, maturityDate, legalDocumentHash, valuationReportHash)
}

// CreateAsset is a paid mutator transaction binding the contract method 0xc183a1a6.
//
// Solidity: function createAsset(string name, string symbol, uint8 assetType, uint256 totalValue, uint256 maturityDate, string legalDocumentHash, string valuationReportHash) returns(uint256 assetId)
func (_RWAAssetFactory *RWAAssetFactorySession) CreateAsset(name string, symbol string, assetType uint8, totalValue *big.Int, maturityDate *big.Int, legalDocumentHash string, valuationReportHash string) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.CreateAsset(&_RWAAssetFactory.TransactOpts, name, symbol, assetType, totalValue, maturityDate, legalDocumentHash, valuationReportHash)
}

// CreateAsset is a paid mutator transaction binding the contract method 0xc183a1a6.
//
// Solidity: function createAsset(string name, string symbol, uint8 assetType, uint256 totalValue, uint256 maturityDate, string legalDocumentHash, string valuationReportHash) returns(uint256 assetId)
func (_RWAAssetFactory *RWAAssetFactoryTransactorSession) CreateAsset(name string, symbol string, assetType uint8, totalValue *big.Int, maturityDate *big.Int, legalDocumentHash string, valuationReportHash string) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.CreateAsset(&_RWAAssetFactory.TransactOpts, name, symbol, assetType, totalValue, maturityDate, legalDocumentHash, valuationReportHash)
}

// DistributeYield is a paid mutator transaction binding the contract method 0xf007e9d5.
//
// Solidity: function distributeYield(uint256 assetId, uint256 amount) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactor) DistributeYield(opts *bind.TransactOpts, assetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RWAAssetFactory.contract.Transact(opts, "distributeYield", assetId, amount)
}

// DistributeYield is a paid mutator transaction binding the contract method 0xf007e9d5.
//
// Solidity: function distributeYield(uint256 assetId, uint256 amount) returns()
func (_RWAAssetFactory *RWAAssetFactorySession) DistributeYield(assetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.DistributeYield(&_RWAAssetFactory.TransactOpts, assetId, amount)
}

// DistributeYield is a paid mutator transaction binding the contract method 0xf007e9d5.
//
// Solidity: function distributeYield(uint256 assetId, uint256 amount) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactorSession) DistributeYield(assetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.DistributeYield(&_RWAAssetFactory.TransactOpts, assetId, amount)
}

// FractionalizeAsset is a paid mutator transaction binding the contract method 0x394b2122.
//
// Solidity: function fractionalizeAsset(uint256 assetId, uint256 totalSupply) returns(address tokenAddress)
func (_RWAAssetFactory *RWAAssetFactoryTransactor) FractionalizeAsset(opts *bind.TransactOpts, assetId *big.Int, totalSupply *big.Int) (*types.Transaction, error) {
	return _RWAAssetFactory.contract.Transact(opts, "fractionalizeAsset", assetId, totalSupply)
}

// FractionalizeAsset is a paid mutator transaction binding the contract method 0x394b2122.
//
// Solidity: function fractionalizeAsset(uint256 assetId, uint256 totalSupply) returns(address tokenAddress)
func (_RWAAssetFactory *RWAAssetFactorySession) FractionalizeAsset(assetId *big.Int, totalSupply *big.Int) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.FractionalizeAsset(&_RWAAssetFactory.TransactOpts, assetId, totalSupply)
}

// FractionalizeAsset is a paid mutator transaction binding the contract method 0x394b2122.
//
// Solidity: function fractionalizeAsset(uint256 assetId, uint256 totalSupply) returns(address tokenAddress)
func (_RWAAssetFactory *RWAAssetFactoryTransactorSession) FractionalizeAsset(assetId *big.Int, totalSupply *big.Int) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.FractionalizeAsset(&_RWAAssetFactory.TransactOpts, assetId, totalSupply)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAAssetFactory.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAAssetFactory *RWAAssetFactorySession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.GrantRole(&_RWAAssetFactory.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.GrantRole(&_RWAAssetFactory.TransactOpts, role, account)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAAssetFactory.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RWAAssetFactory *RWAAssetFactorySession) Pause() (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.Pause(&_RWAAssetFactory.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactorSession) Pause() (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.Pause(&_RWAAssetFactory.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAAssetFactory.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAAssetFactory *RWAAssetFactorySession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.RenounceRole(&_RWAAssetFactory.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.RenounceRole(&_RWAAssetFactory.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAAssetFactory.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAAssetFactory *RWAAssetFactorySession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.RevokeRole(&_RWAAssetFactory.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.RevokeRole(&_RWAAssetFactory.TransactOpts, role, account)
}

// SetYieldTerms is a paid mutator transaction binding the contract method 0xe4d3c825.
//
// Solidity: function setYieldTerms(uint256 assetId, uint256 annualYieldRate, uint256 yieldPaymentFrequency) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactor) SetYieldTerms(opts *bind.TransactOpts, assetId *big.Int, annualYieldRate *big.Int, yieldPaymentFrequency *big.Int) (*types.Transaction, error) {
	return _RWAAssetFactory.contract.Transact(opts, "setYieldTerms", assetId, annualYieldRate, yieldPaymentFrequency)
}

// SetYieldTerms is a paid mutator transaction binding the contract method 0xe4d3c825.
//
// Solidity: function setYieldTerms(uint256 assetId, uint256 annualYieldRate, uint256 yieldPaymentFrequency) returns()
func (_RWAAssetFactory *RWAAssetFactorySession) SetYieldTerms(assetId *big.Int, annualYieldRate *big.Int, yieldPaymentFrequency *big.Int) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.SetYieldTerms(&_RWAAssetFactory.TransactOpts, assetId, annualYieldRate, yieldPaymentFrequency)
}

// SetYieldTerms is a paid mutator transaction binding the contract method 0xe4d3c825.
//
// Solidity: function setYieldTerms(uint256 assetId, uint256 annualYieldRate, uint256 yieldPaymentFrequency) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactorSession) SetYieldTerms(assetId *big.Int, annualYieldRate *big.Int, yieldPaymentFrequency *big.Int) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.SetYieldTerms(&_RWAAssetFactory.TransactOpts, assetId, annualYieldRate, yieldPaymentFrequency)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAAssetFactory.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RWAAssetFactory *RWAAssetFactorySession) Unpause() (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.Unpause(&_RWAAssetFactory.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactorSession) Unpause() (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.Unpause(&_RWAAssetFactory.TransactOpts)
}

// UpdateAssetStatus is a paid mutator transaction binding the contract method 0x4554f3cf.
//
// Solidity: function updateAssetStatus(uint256 assetId, uint8 newStatus) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactor) UpdateAssetStatus(opts *bind.TransactOpts, assetId *big.Int, newStatus uint8) (*types.Transaction, error) {
	return _RWAAssetFactory.contract.Transact(opts, "updateAssetStatus", assetId, newStatus)
}

// UpdateAssetStatus is a paid mutator transaction binding the contract method 0x4554f3cf.
//
// Solidity: function updateAssetStatus(uint256 assetId, uint8 newStatus) returns()
func (_RWAAssetFactory *RWAAssetFactorySession) UpdateAssetStatus(assetId *big.Int, newStatus uint8) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.UpdateAssetStatus(&_RWAAssetFactory.TransactOpts, assetId, newStatus)
}

// UpdateAssetStatus is a paid mutator transaction binding the contract method 0x4554f3cf.
//
// Solidity: function updateAssetStatus(uint256 assetId, uint8 newStatus) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactorSession) UpdateAssetStatus(assetId *big.Int, newStatus uint8) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.UpdateAssetStatus(&_RWAAssetFactory.TransactOpts, assetId, newStatus)
}

// UpdateAssetValuation is a paid mutator transaction binding the contract method 0xa1a5062b.
//
// Solidity: function updateAssetValuation(uint256 assetId, uint256 newValue, string valuationReportHash) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactor) UpdateAssetValuation(opts *bind.TransactOpts, assetId *big.Int, newValue *big.Int, valuationReportHash string) (*types.Transaction, error) {
	return _RWAAssetFactory.contract.Transact(opts, "updateAssetValuation", assetId, newValue, valuationReportHash)
}

// UpdateAssetValuation is a paid mutator transaction binding the contract method 0xa1a5062b.
//
// Solidity: function updateAssetValuation(uint256 assetId, uint256 newValue, string valuationReportHash) returns()
func (_RWAAssetFactory *RWAAssetFactorySession) UpdateAssetValuation(assetId *big.Int, newValue *big.Int, valuationReportHash string) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.UpdateAssetValuation(&_RWAAssetFactory.TransactOpts, assetId, newValue, valuationReportHash)
}

// UpdateAssetValuation is a paid mutator transaction binding the contract method 0xa1a5062b.
//
// Solidity: function updateAssetValuation(uint256 assetId, uint256 newValue, string valuationReportHash) returns()
func (_RWAAssetFactory *RWAAssetFactoryTransactorSession) UpdateAssetValuation(assetId *big.Int, newValue *big.Int, valuationReportHash string) (*types.Transaction, error) {
	return _RWAAssetFactory.Contract.UpdateAssetValuation(&_RWAAssetFactory.TransactOpts, assetId, newValue, valuationReportHash)
}

// RWAAssetFactoryAssetCreatedIterator is returned from FilterAssetCreated and is used to iterate over the raw logs and unpacked data for AssetCreated events raised by the RWAAssetFactory contract.
type RWAAssetFactoryAssetCreatedIterator struct {
	Event *RWAAssetFactoryAssetCreated // Event containing the contract specifics and raw log

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
func (it *RWAAssetFactoryAssetCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAAssetFactoryAssetCreated)
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
		it.Event = new(RWAAssetFactoryAssetCreated)
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
func (it *RWAAssetFactoryAssetCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAAssetFactoryAssetCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAAssetFactoryAssetCreated represents a AssetCreated event raised by the RWAAssetFactory contract.
type RWAAssetFactoryAssetCreated struct {
	AssetId    *big.Int
	Issuer     common.Address
	AssetType  uint8
	TotalValue *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAssetCreated is a free log retrieval operation binding the contract event 0x879557c82af2c8be83d28c5d82d6a216bccb4f00d96c23d68117050d86945aa6.
//
// Solidity: event AssetCreated(uint256 indexed assetId, address indexed issuer, uint8 assetType, uint256 totalValue)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) FilterAssetCreated(opts *bind.FilterOpts, assetId []*big.Int, issuer []common.Address) (*RWAAssetFactoryAssetCreatedIterator, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}

	logs, sub, err := _RWAAssetFactory.contract.FilterLogs(opts, "AssetCreated", assetIdRule, issuerRule)
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryAssetCreatedIterator{contract: _RWAAssetFactory.contract, event: "AssetCreated", logs: logs, sub: sub}, nil
}

// WatchAssetCreated is a free log subscription operation binding the contract event 0x879557c82af2c8be83d28c5d82d6a216bccb4f00d96c23d68117050d86945aa6.
//
// Solidity: event AssetCreated(uint256 indexed assetId, address indexed issuer, uint8 assetType, uint256 totalValue)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) WatchAssetCreated(opts *bind.WatchOpts, sink chan<- *RWAAssetFactoryAssetCreated, assetId []*big.Int, issuer []common.Address) (event.Subscription, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}

	logs, sub, err := _RWAAssetFactory.contract.WatchLogs(opts, "AssetCreated", assetIdRule, issuerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAAssetFactoryAssetCreated)
				if err := _RWAAssetFactory.contract.UnpackLog(event, "AssetCreated", log); err != nil {
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

// ParseAssetCreated is a log parse operation binding the contract event 0x879557c82af2c8be83d28c5d82d6a216bccb4f00d96c23d68117050d86945aa6.
//
// Solidity: event AssetCreated(uint256 indexed assetId, address indexed issuer, uint8 assetType, uint256 totalValue)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) ParseAssetCreated(log types.Log) (*RWAAssetFactoryAssetCreated, error) {
	event := new(RWAAssetFactoryAssetCreated)
	if err := _RWAAssetFactory.contract.UnpackLog(event, "AssetCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAAssetFactoryAssetFractionalizedIterator is returned from FilterAssetFractionalized and is used to iterate over the raw logs and unpacked data for AssetFractionalized events raised by the RWAAssetFactory contract.
type RWAAssetFactoryAssetFractionalizedIterator struct {
	Event *RWAAssetFactoryAssetFractionalized // Event containing the contract specifics and raw log

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
func (it *RWAAssetFactoryAssetFractionalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAAssetFactoryAssetFractionalized)
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
		it.Event = new(RWAAssetFactoryAssetFractionalized)
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
func (it *RWAAssetFactoryAssetFractionalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAAssetFactoryAssetFractionalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAAssetFactoryAssetFractionalized represents a AssetFractionalized event raised by the RWAAssetFactory contract.
type RWAAssetFactoryAssetFractionalized struct {
	AssetId         *big.Int
	TotalSupply     *big.Int
	FractionalToken common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterAssetFractionalized is a free log retrieval operation binding the contract event 0x5d1d18878d47a9f0ef56fd4e9baf222467ab793967f7aaa3d3a7c8c2db482bb1.
//
// Solidity: event AssetFractionalized(uint256 indexed assetId, uint256 totalSupply, address fractionalToken)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) FilterAssetFractionalized(opts *bind.FilterOpts, assetId []*big.Int) (*RWAAssetFactoryAssetFractionalizedIterator, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWAAssetFactory.contract.FilterLogs(opts, "AssetFractionalized", assetIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryAssetFractionalizedIterator{contract: _RWAAssetFactory.contract, event: "AssetFractionalized", logs: logs, sub: sub}, nil
}

// WatchAssetFractionalized is a free log subscription operation binding the contract event 0x5d1d18878d47a9f0ef56fd4e9baf222467ab793967f7aaa3d3a7c8c2db482bb1.
//
// Solidity: event AssetFractionalized(uint256 indexed assetId, uint256 totalSupply, address fractionalToken)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) WatchAssetFractionalized(opts *bind.WatchOpts, sink chan<- *RWAAssetFactoryAssetFractionalized, assetId []*big.Int) (event.Subscription, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWAAssetFactory.contract.WatchLogs(opts, "AssetFractionalized", assetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAAssetFactoryAssetFractionalized)
				if err := _RWAAssetFactory.contract.UnpackLog(event, "AssetFractionalized", log); err != nil {
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

// ParseAssetFractionalized is a log parse operation binding the contract event 0x5d1d18878d47a9f0ef56fd4e9baf222467ab793967f7aaa3d3a7c8c2db482bb1.
//
// Solidity: event AssetFractionalized(uint256 indexed assetId, uint256 totalSupply, address fractionalToken)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) ParseAssetFractionalized(log types.Log) (*RWAAssetFactoryAssetFractionalized, error) {
	event := new(RWAAssetFactoryAssetFractionalized)
	if err := _RWAAssetFactory.contract.UnpackLog(event, "AssetFractionalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAAssetFactoryAssetStatusChangedIterator is returned from FilterAssetStatusChanged and is used to iterate over the raw logs and unpacked data for AssetStatusChanged events raised by the RWAAssetFactory contract.
type RWAAssetFactoryAssetStatusChangedIterator struct {
	Event *RWAAssetFactoryAssetStatusChanged // Event containing the contract specifics and raw log

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
func (it *RWAAssetFactoryAssetStatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAAssetFactoryAssetStatusChanged)
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
		it.Event = new(RWAAssetFactoryAssetStatusChanged)
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
func (it *RWAAssetFactoryAssetStatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAAssetFactoryAssetStatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAAssetFactoryAssetStatusChanged represents a AssetStatusChanged event raised by the RWAAssetFactory contract.
type RWAAssetFactoryAssetStatusChanged struct {
	AssetId   *big.Int
	OldStatus uint8
	NewStatus uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAssetStatusChanged is a free log retrieval operation binding the contract event 0xadf241d656a852a4cb14a0ae6771aaa2cd30d046fb5f646395777ffd80f5a859.
//
// Solidity: event AssetStatusChanged(uint256 indexed assetId, uint8 oldStatus, uint8 newStatus)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) FilterAssetStatusChanged(opts *bind.FilterOpts, assetId []*big.Int) (*RWAAssetFactoryAssetStatusChangedIterator, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWAAssetFactory.contract.FilterLogs(opts, "AssetStatusChanged", assetIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryAssetStatusChangedIterator{contract: _RWAAssetFactory.contract, event: "AssetStatusChanged", logs: logs, sub: sub}, nil
}

// WatchAssetStatusChanged is a free log subscription operation binding the contract event 0xadf241d656a852a4cb14a0ae6771aaa2cd30d046fb5f646395777ffd80f5a859.
//
// Solidity: event AssetStatusChanged(uint256 indexed assetId, uint8 oldStatus, uint8 newStatus)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) WatchAssetStatusChanged(opts *bind.WatchOpts, sink chan<- *RWAAssetFactoryAssetStatusChanged, assetId []*big.Int) (event.Subscription, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWAAssetFactory.contract.WatchLogs(opts, "AssetStatusChanged", assetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAAssetFactoryAssetStatusChanged)
				if err := _RWAAssetFactory.contract.UnpackLog(event, "AssetStatusChanged", log); err != nil {
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

// ParseAssetStatusChanged is a log parse operation binding the contract event 0xadf241d656a852a4cb14a0ae6771aaa2cd30d046fb5f646395777ffd80f5a859.
//
// Solidity: event AssetStatusChanged(uint256 indexed assetId, uint8 oldStatus, uint8 newStatus)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) ParseAssetStatusChanged(log types.Log) (*RWAAssetFactoryAssetStatusChanged, error) {
	event := new(RWAAssetFactoryAssetStatusChanged)
	if err := _RWAAssetFactory.contract.UnpackLog(event, "AssetStatusChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAAssetFactoryAssetValuationUpdatedIterator is returned from FilterAssetValuationUpdated and is used to iterate over the raw logs and unpacked data for AssetValuationUpdated events raised by the RWAAssetFactory contract.
type RWAAssetFactoryAssetValuationUpdatedIterator struct {
	Event *RWAAssetFactoryAssetValuationUpdated // Event containing the contract specifics and raw log

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
func (it *RWAAssetFactoryAssetValuationUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAAssetFactoryAssetValuationUpdated)
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
		it.Event = new(RWAAssetFactoryAssetValuationUpdated)
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
func (it *RWAAssetFactoryAssetValuationUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAAssetFactoryAssetValuationUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAAssetFactoryAssetValuationUpdated represents a AssetValuationUpdated event raised by the RWAAssetFactory contract.
type RWAAssetFactoryAssetValuationUpdated struct {
	AssetId   *big.Int
	OldValue  *big.Int
	NewValue  *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAssetValuationUpdated is a free log retrieval operation binding the contract event 0xcf08dea81404ce127cec13182d2dc4831f8ff76c67fa3043fbdba3cf13a221e1.
//
// Solidity: event AssetValuationUpdated(uint256 indexed assetId, uint256 oldValue, uint256 newValue, uint256 timestamp)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) FilterAssetValuationUpdated(opts *bind.FilterOpts, assetId []*big.Int) (*RWAAssetFactoryAssetValuationUpdatedIterator, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWAAssetFactory.contract.FilterLogs(opts, "AssetValuationUpdated", assetIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryAssetValuationUpdatedIterator{contract: _RWAAssetFactory.contract, event: "AssetValuationUpdated", logs: logs, sub: sub}, nil
}

// WatchAssetValuationUpdated is a free log subscription operation binding the contract event 0xcf08dea81404ce127cec13182d2dc4831f8ff76c67fa3043fbdba3cf13a221e1.
//
// Solidity: event AssetValuationUpdated(uint256 indexed assetId, uint256 oldValue, uint256 newValue, uint256 timestamp)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) WatchAssetValuationUpdated(opts *bind.WatchOpts, sink chan<- *RWAAssetFactoryAssetValuationUpdated, assetId []*big.Int) (event.Subscription, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWAAssetFactory.contract.WatchLogs(opts, "AssetValuationUpdated", assetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAAssetFactoryAssetValuationUpdated)
				if err := _RWAAssetFactory.contract.UnpackLog(event, "AssetValuationUpdated", log); err != nil {
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

// ParseAssetValuationUpdated is a log parse operation binding the contract event 0xcf08dea81404ce127cec13182d2dc4831f8ff76c67fa3043fbdba3cf13a221e1.
//
// Solidity: event AssetValuationUpdated(uint256 indexed assetId, uint256 oldValue, uint256 newValue, uint256 timestamp)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) ParseAssetValuationUpdated(log types.Log) (*RWAAssetFactoryAssetValuationUpdated, error) {
	event := new(RWAAssetFactoryAssetValuationUpdated)
	if err := _RWAAssetFactory.contract.UnpackLog(event, "AssetValuationUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAAssetFactoryPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the RWAAssetFactory contract.
type RWAAssetFactoryPausedIterator struct {
	Event *RWAAssetFactoryPaused // Event containing the contract specifics and raw log

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
func (it *RWAAssetFactoryPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAAssetFactoryPaused)
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
		it.Event = new(RWAAssetFactoryPaused)
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
func (it *RWAAssetFactoryPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAAssetFactoryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAAssetFactoryPaused represents a Paused event raised by the RWAAssetFactory contract.
type RWAAssetFactoryPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) FilterPaused(opts *bind.FilterOpts) (*RWAAssetFactoryPausedIterator, error) {

	logs, sub, err := _RWAAssetFactory.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryPausedIterator{contract: _RWAAssetFactory.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *RWAAssetFactoryPaused) (event.Subscription, error) {

	logs, sub, err := _RWAAssetFactory.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAAssetFactoryPaused)
				if err := _RWAAssetFactory.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_RWAAssetFactory *RWAAssetFactoryFilterer) ParsePaused(log types.Log) (*RWAAssetFactoryPaused, error) {
	event := new(RWAAssetFactoryPaused)
	if err := _RWAAssetFactory.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAAssetFactoryRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the RWAAssetFactory contract.
type RWAAssetFactoryRoleAdminChangedIterator struct {
	Event *RWAAssetFactoryRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *RWAAssetFactoryRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAAssetFactoryRoleAdminChanged)
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
		it.Event = new(RWAAssetFactoryRoleAdminChanged)
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
func (it *RWAAssetFactoryRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAAssetFactoryRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAAssetFactoryRoleAdminChanged represents a RoleAdminChanged event raised by the RWAAssetFactory contract.
type RWAAssetFactoryRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*RWAAssetFactoryRoleAdminChangedIterator, error) {

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

	logs, sub, err := _RWAAssetFactory.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryRoleAdminChangedIterator{contract: _RWAAssetFactory.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *RWAAssetFactoryRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _RWAAssetFactory.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAAssetFactoryRoleAdminChanged)
				if err := _RWAAssetFactory.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_RWAAssetFactory *RWAAssetFactoryFilterer) ParseRoleAdminChanged(log types.Log) (*RWAAssetFactoryRoleAdminChanged, error) {
	event := new(RWAAssetFactoryRoleAdminChanged)
	if err := _RWAAssetFactory.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAAssetFactoryRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the RWAAssetFactory contract.
type RWAAssetFactoryRoleGrantedIterator struct {
	Event *RWAAssetFactoryRoleGranted // Event containing the contract specifics and raw log

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
func (it *RWAAssetFactoryRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAAssetFactoryRoleGranted)
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
		it.Event = new(RWAAssetFactoryRoleGranted)
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
func (it *RWAAssetFactoryRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAAssetFactoryRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAAssetFactoryRoleGranted represents a RoleGranted event raised by the RWAAssetFactory contract.
type RWAAssetFactoryRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RWAAssetFactoryRoleGrantedIterator, error) {

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

	logs, sub, err := _RWAAssetFactory.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryRoleGrantedIterator{contract: _RWAAssetFactory.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *RWAAssetFactoryRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RWAAssetFactory.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAAssetFactoryRoleGranted)
				if err := _RWAAssetFactory.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_RWAAssetFactory *RWAAssetFactoryFilterer) ParseRoleGranted(log types.Log) (*RWAAssetFactoryRoleGranted, error) {
	event := new(RWAAssetFactoryRoleGranted)
	if err := _RWAAssetFactory.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAAssetFactoryRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the RWAAssetFactory contract.
type RWAAssetFactoryRoleRevokedIterator struct {
	Event *RWAAssetFactoryRoleRevoked // Event containing the contract specifics and raw log

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
func (it *RWAAssetFactoryRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAAssetFactoryRoleRevoked)
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
		it.Event = new(RWAAssetFactoryRoleRevoked)
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
func (it *RWAAssetFactoryRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAAssetFactoryRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAAssetFactoryRoleRevoked represents a RoleRevoked event raised by the RWAAssetFactory contract.
type RWAAssetFactoryRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RWAAssetFactoryRoleRevokedIterator, error) {

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

	logs, sub, err := _RWAAssetFactory.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryRoleRevokedIterator{contract: _RWAAssetFactory.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *RWAAssetFactoryRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RWAAssetFactory.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAAssetFactoryRoleRevoked)
				if err := _RWAAssetFactory.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_RWAAssetFactory *RWAAssetFactoryFilterer) ParseRoleRevoked(log types.Log) (*RWAAssetFactoryRoleRevoked, error) {
	event := new(RWAAssetFactoryRoleRevoked)
	if err := _RWAAssetFactory.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAAssetFactoryUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the RWAAssetFactory contract.
type RWAAssetFactoryUnpausedIterator struct {
	Event *RWAAssetFactoryUnpaused // Event containing the contract specifics and raw log

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
func (it *RWAAssetFactoryUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAAssetFactoryUnpaused)
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
		it.Event = new(RWAAssetFactoryUnpaused)
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
func (it *RWAAssetFactoryUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAAssetFactoryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAAssetFactoryUnpaused represents a Unpaused event raised by the RWAAssetFactory contract.
type RWAAssetFactoryUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*RWAAssetFactoryUnpausedIterator, error) {

	logs, sub, err := _RWAAssetFactory.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryUnpausedIterator{contract: _RWAAssetFactory.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *RWAAssetFactoryUnpaused) (event.Subscription, error) {

	logs, sub, err := _RWAAssetFactory.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAAssetFactoryUnpaused)
				if err := _RWAAssetFactory.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_RWAAssetFactory *RWAAssetFactoryFilterer) ParseUnpaused(log types.Log) (*RWAAssetFactoryUnpaused, error) {
	event := new(RWAAssetFactoryUnpaused)
	if err := _RWAAssetFactory.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAAssetFactoryYieldDistributedIterator is returned from FilterYieldDistributed and is used to iterate over the raw logs and unpacked data for YieldDistributed events raised by the RWAAssetFactory contract.
type RWAAssetFactoryYieldDistributedIterator struct {
	Event *RWAAssetFactoryYieldDistributed // Event containing the contract specifics and raw log

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
func (it *RWAAssetFactoryYieldDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAAssetFactoryYieldDistributed)
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
		it.Event = new(RWAAssetFactoryYieldDistributed)
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
func (it *RWAAssetFactoryYieldDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAAssetFactoryYieldDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAAssetFactoryYieldDistributed represents a YieldDistributed event raised by the RWAAssetFactory contract.
type RWAAssetFactoryYieldDistributed struct {
	AssetId   *big.Int
	Amount    *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterYieldDistributed is a free log retrieval operation binding the contract event 0xf571bb88b0ec5dc25398fc62c69d5970d50ea313a0a6596b0ca7f0197d37383f.
//
// Solidity: event YieldDistributed(uint256 indexed assetId, uint256 amount, uint256 timestamp)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) FilterYieldDistributed(opts *bind.FilterOpts, assetId []*big.Int) (*RWAAssetFactoryYieldDistributedIterator, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWAAssetFactory.contract.FilterLogs(opts, "YieldDistributed", assetIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAAssetFactoryYieldDistributedIterator{contract: _RWAAssetFactory.contract, event: "YieldDistributed", logs: logs, sub: sub}, nil
}

// WatchYieldDistributed is a free log subscription operation binding the contract event 0xf571bb88b0ec5dc25398fc62c69d5970d50ea313a0a6596b0ca7f0197d37383f.
//
// Solidity: event YieldDistributed(uint256 indexed assetId, uint256 amount, uint256 timestamp)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) WatchYieldDistributed(opts *bind.WatchOpts, sink chan<- *RWAAssetFactoryYieldDistributed, assetId []*big.Int) (event.Subscription, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWAAssetFactory.contract.WatchLogs(opts, "YieldDistributed", assetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAAssetFactoryYieldDistributed)
				if err := _RWAAssetFactory.contract.UnpackLog(event, "YieldDistributed", log); err != nil {
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

// ParseYieldDistributed is a log parse operation binding the contract event 0xf571bb88b0ec5dc25398fc62c69d5970d50ea313a0a6596b0ca7f0197d37383f.
//
// Solidity: event YieldDistributed(uint256 indexed assetId, uint256 amount, uint256 timestamp)
func (_RWAAssetFactory *RWAAssetFactoryFilterer) ParseYieldDistributed(log types.Log) (*RWAAssetFactoryYieldDistributed, error) {
	event := new(RWAAssetFactoryYieldDistributed)
	if err := _RWAAssetFactory.contract.UnpackLog(event, "YieldDistributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
