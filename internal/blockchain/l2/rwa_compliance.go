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

// IRWAComplianceInvestorProfile is an auto generated low-level Go binding around an user-defined struct.
type IRWAComplianceInvestorProfile struct {
	Status          uint8
	Tier            uint8
	VerifiedAt      *big.Int
	ExpiresAt       *big.Int
	Jurisdiction    string
	KycDocumentHash [32]byte
	Verifier        common.Address
}

// RWAComplianceMetaData contains all meta data concerning the RWACompliance contract.
var RWAComplianceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumIRWACompliance.AccreditationTier\",\"name\":\"minTier\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"requiresKYC\",\"type\":\"bool\"}],\"name\":\"AssetComplianceSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"ComplianceViolation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumIRWACompliance.VerificationStatus\",\"name\":\"oldStatus\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumIRWACompliance.VerificationStatus\",\"name\":\"newStatus\",\"type\":\"uint8\"}],\"name\":\"InvestorStatusChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumIRWACompliance.VerificationStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumIRWACompliance.AccreditationTier\",\"name\":\"tier\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiresAt\",\"type\":\"uint256\"}],\"name\":\"InvestorVerified\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"jurisdiction\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"JurisdictionUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COMPLIANCE_OFFICER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_VERIFICATION_PERIOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_VERIFICATION_PERIOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"VERIFIER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"canInvestInAsset\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"canTransferTokens\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"doesAssetRequireKYC\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getAssetInvestorCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getAssetMinTier\",\"outputs\":[{\"internalType\":\"enumIRWACompliance.AccreditationTier\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"}],\"name\":\"getInvestorProfile\",\"outputs\":[{\"components\":[{\"internalType\":\"enumIRWACompliance.VerificationStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"enumIRWACompliance.AccreditationTier\",\"name\":\"tier\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"verifiedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiresAt\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"jurisdiction\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"kycDocumentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"}],\"internalType\":\"structIRWACompliance.InvestorProfile\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"}],\"name\":\"isInvestorVerified\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"jurisdiction\",\"type\":\"string\"}],\"name\":\"isJurisdictionAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"}],\"name\":\"recordInvestorParticipation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"}],\"name\":\"removeInvestorParticipation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"enumIRWACompliance.AccreditationTier\",\"name\":\"minTier\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"requiresKYC\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"minInvestmentAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxInvestors\",\"type\":\"uint256\"}],\"name\":\"setAssetCompliance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalAssetsWithCompliance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalVerifiedInvestors\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"},{\"internalType\":\"enumIRWACompliance.VerificationStatus\",\"name\":\"newStatus\",\"type\":\"uint8\"}],\"name\":\"updateInvestorStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"jurisdiction\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"updateJurisdictionRestriction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"investor\",\"type\":\"address\"},{\"internalType\":\"enumIRWACompliance.AccreditationTier\",\"name\":\"tier\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"jurisdiction\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"validityPeriod\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"kycDocumentHash\",\"type\":\"bytes32\"}],\"name\":\"verifyInvestor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RWAComplianceABI is the input ABI used to generate the binding from.
// Deprecated: Use RWAComplianceMetaData.ABI instead.
var RWAComplianceABI = RWAComplianceMetaData.ABI

// RWACompliance is an auto generated Go binding around an Ethereum contract.
type RWACompliance struct {
	RWAComplianceCaller     // Read-only binding to the contract
	RWAComplianceTransactor // Write-only binding to the contract
	RWAComplianceFilterer   // Log filterer for contract events
}

// RWAComplianceCaller is an auto generated read-only Go binding around an Ethereum contract.
type RWAComplianceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAComplianceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RWAComplianceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAComplianceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RWAComplianceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAComplianceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RWAComplianceSession struct {
	Contract     *RWACompliance    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RWAComplianceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RWAComplianceCallerSession struct {
	Contract *RWAComplianceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// RWAComplianceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RWAComplianceTransactorSession struct {
	Contract     *RWAComplianceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// RWAComplianceRaw is an auto generated low-level Go binding around an Ethereum contract.
type RWAComplianceRaw struct {
	Contract *RWACompliance // Generic contract binding to access the raw methods on
}

// RWAComplianceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RWAComplianceCallerRaw struct {
	Contract *RWAComplianceCaller // Generic read-only contract binding to access the raw methods on
}

// RWAComplianceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RWAComplianceTransactorRaw struct {
	Contract *RWAComplianceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRWACompliance creates a new instance of RWACompliance, bound to a specific deployed contract.
func NewRWACompliance(address common.Address, backend bind.ContractBackend) (*RWACompliance, error) {
	contract, err := bindRWACompliance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RWACompliance{RWAComplianceCaller: RWAComplianceCaller{contract: contract}, RWAComplianceTransactor: RWAComplianceTransactor{contract: contract}, RWAComplianceFilterer: RWAComplianceFilterer{contract: contract}}, nil
}

// NewRWAComplianceCaller creates a new read-only instance of RWACompliance, bound to a specific deployed contract.
func NewRWAComplianceCaller(address common.Address, caller bind.ContractCaller) (*RWAComplianceCaller, error) {
	contract, err := bindRWACompliance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RWAComplianceCaller{contract: contract}, nil
}

// NewRWAComplianceTransactor creates a new write-only instance of RWACompliance, bound to a specific deployed contract.
func NewRWAComplianceTransactor(address common.Address, transactor bind.ContractTransactor) (*RWAComplianceTransactor, error) {
	contract, err := bindRWACompliance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RWAComplianceTransactor{contract: contract}, nil
}

// NewRWAComplianceFilterer creates a new log filterer instance of RWACompliance, bound to a specific deployed contract.
func NewRWAComplianceFilterer(address common.Address, filterer bind.ContractFilterer) (*RWAComplianceFilterer, error) {
	contract, err := bindRWACompliance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RWAComplianceFilterer{contract: contract}, nil
}

// bindRWACompliance binds a generic wrapper to an already deployed contract.
func bindRWACompliance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RWAComplianceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RWACompliance *RWAComplianceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RWACompliance.Contract.RWAComplianceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RWACompliance *RWAComplianceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWACompliance.Contract.RWAComplianceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RWACompliance *RWAComplianceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RWACompliance.Contract.RWAComplianceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RWACompliance *RWAComplianceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RWACompliance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RWACompliance *RWAComplianceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWACompliance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RWACompliance *RWAComplianceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RWACompliance.Contract.contract.Transact(opts, method, params...)
}

// COMPLIANCEOFFICERROLE is a free data retrieval call binding the contract method 0x198596b5.
//
// Solidity: function COMPLIANCE_OFFICER_ROLE() view returns(bytes32)
func (_RWACompliance *RWAComplianceCaller) COMPLIANCEOFFICERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "COMPLIANCE_OFFICER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// COMPLIANCEOFFICERROLE is a free data retrieval call binding the contract method 0x198596b5.
//
// Solidity: function COMPLIANCE_OFFICER_ROLE() view returns(bytes32)
func (_RWACompliance *RWAComplianceSession) COMPLIANCEOFFICERROLE() ([32]byte, error) {
	return _RWACompliance.Contract.COMPLIANCEOFFICERROLE(&_RWACompliance.CallOpts)
}

// COMPLIANCEOFFICERROLE is a free data retrieval call binding the contract method 0x198596b5.
//
// Solidity: function COMPLIANCE_OFFICER_ROLE() view returns(bytes32)
func (_RWACompliance *RWAComplianceCallerSession) COMPLIANCEOFFICERROLE() ([32]byte, error) {
	return _RWACompliance.Contract.COMPLIANCEOFFICERROLE(&_RWACompliance.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWACompliance *RWAComplianceCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWACompliance *RWAComplianceSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RWACompliance.Contract.DEFAULTADMINROLE(&_RWACompliance.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWACompliance *RWAComplianceCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RWACompliance.Contract.DEFAULTADMINROLE(&_RWACompliance.CallOpts)
}

// DEFAULTVERIFICATIONPERIOD is a free data retrieval call binding the contract method 0x79168430.
//
// Solidity: function DEFAULT_VERIFICATION_PERIOD() view returns(uint256)
func (_RWACompliance *RWAComplianceCaller) DEFAULTVERIFICATIONPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "DEFAULT_VERIFICATION_PERIOD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEFAULTVERIFICATIONPERIOD is a free data retrieval call binding the contract method 0x79168430.
//
// Solidity: function DEFAULT_VERIFICATION_PERIOD() view returns(uint256)
func (_RWACompliance *RWAComplianceSession) DEFAULTVERIFICATIONPERIOD() (*big.Int, error) {
	return _RWACompliance.Contract.DEFAULTVERIFICATIONPERIOD(&_RWACompliance.CallOpts)
}

// DEFAULTVERIFICATIONPERIOD is a free data retrieval call binding the contract method 0x79168430.
//
// Solidity: function DEFAULT_VERIFICATION_PERIOD() view returns(uint256)
func (_RWACompliance *RWAComplianceCallerSession) DEFAULTVERIFICATIONPERIOD() (*big.Int, error) {
	return _RWACompliance.Contract.DEFAULTVERIFICATIONPERIOD(&_RWACompliance.CallOpts)
}

// MAXVERIFICATIONPERIOD is a free data retrieval call binding the contract method 0x449b9290.
//
// Solidity: function MAX_VERIFICATION_PERIOD() view returns(uint256)
func (_RWACompliance *RWAComplianceCaller) MAXVERIFICATIONPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "MAX_VERIFICATION_PERIOD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXVERIFICATIONPERIOD is a free data retrieval call binding the contract method 0x449b9290.
//
// Solidity: function MAX_VERIFICATION_PERIOD() view returns(uint256)
func (_RWACompliance *RWAComplianceSession) MAXVERIFICATIONPERIOD() (*big.Int, error) {
	return _RWACompliance.Contract.MAXVERIFICATIONPERIOD(&_RWACompliance.CallOpts)
}

// MAXVERIFICATIONPERIOD is a free data retrieval call binding the contract method 0x449b9290.
//
// Solidity: function MAX_VERIFICATION_PERIOD() view returns(uint256)
func (_RWACompliance *RWAComplianceCallerSession) MAXVERIFICATIONPERIOD() (*big.Int, error) {
	return _RWACompliance.Contract.MAXVERIFICATIONPERIOD(&_RWACompliance.CallOpts)
}

// VERIFIERROLE is a free data retrieval call binding the contract method 0xe7705db6.
//
// Solidity: function VERIFIER_ROLE() view returns(bytes32)
func (_RWACompliance *RWAComplianceCaller) VERIFIERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "VERIFIER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// VERIFIERROLE is a free data retrieval call binding the contract method 0xe7705db6.
//
// Solidity: function VERIFIER_ROLE() view returns(bytes32)
func (_RWACompliance *RWAComplianceSession) VERIFIERROLE() ([32]byte, error) {
	return _RWACompliance.Contract.VERIFIERROLE(&_RWACompliance.CallOpts)
}

// VERIFIERROLE is a free data retrieval call binding the contract method 0xe7705db6.
//
// Solidity: function VERIFIER_ROLE() view returns(bytes32)
func (_RWACompliance *RWAComplianceCallerSession) VERIFIERROLE() ([32]byte, error) {
	return _RWACompliance.Contract.VERIFIERROLE(&_RWACompliance.CallOpts)
}

// CanInvestInAsset is a free data retrieval call binding the contract method 0xc8198ebe.
//
// Solidity: function canInvestInAsset(address investor, uint256 assetId) view returns(bool)
func (_RWACompliance *RWAComplianceCaller) CanInvestInAsset(opts *bind.CallOpts, investor common.Address, assetId *big.Int) (bool, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "canInvestInAsset", investor, assetId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanInvestInAsset is a free data retrieval call binding the contract method 0xc8198ebe.
//
// Solidity: function canInvestInAsset(address investor, uint256 assetId) view returns(bool)
func (_RWACompliance *RWAComplianceSession) CanInvestInAsset(investor common.Address, assetId *big.Int) (bool, error) {
	return _RWACompliance.Contract.CanInvestInAsset(&_RWACompliance.CallOpts, investor, assetId)
}

// CanInvestInAsset is a free data retrieval call binding the contract method 0xc8198ebe.
//
// Solidity: function canInvestInAsset(address investor, uint256 assetId) view returns(bool)
func (_RWACompliance *RWAComplianceCallerSession) CanInvestInAsset(investor common.Address, assetId *big.Int) (bool, error) {
	return _RWACompliance.Contract.CanInvestInAsset(&_RWACompliance.CallOpts, investor, assetId)
}

// CanTransferTokens is a free data retrieval call binding the contract method 0x4ac6165f.
//
// Solidity: function canTransferTokens(address from, address to, uint256 assetId, uint256 amount) view returns(bool)
func (_RWACompliance *RWAComplianceCaller) CanTransferTokens(opts *bind.CallOpts, from common.Address, to common.Address, assetId *big.Int, amount *big.Int) (bool, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "canTransferTokens", from, to, assetId, amount)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanTransferTokens is a free data retrieval call binding the contract method 0x4ac6165f.
//
// Solidity: function canTransferTokens(address from, address to, uint256 assetId, uint256 amount) view returns(bool)
func (_RWACompliance *RWAComplianceSession) CanTransferTokens(from common.Address, to common.Address, assetId *big.Int, amount *big.Int) (bool, error) {
	return _RWACompliance.Contract.CanTransferTokens(&_RWACompliance.CallOpts, from, to, assetId, amount)
}

// CanTransferTokens is a free data retrieval call binding the contract method 0x4ac6165f.
//
// Solidity: function canTransferTokens(address from, address to, uint256 assetId, uint256 amount) view returns(bool)
func (_RWACompliance *RWAComplianceCallerSession) CanTransferTokens(from common.Address, to common.Address, assetId *big.Int, amount *big.Int) (bool, error) {
	return _RWACompliance.Contract.CanTransferTokens(&_RWACompliance.CallOpts, from, to, assetId, amount)
}

// DoesAssetRequireKYC is a free data retrieval call binding the contract method 0xcd9be9e7.
//
// Solidity: function doesAssetRequireKYC(uint256 assetId) view returns(bool)
func (_RWACompliance *RWAComplianceCaller) DoesAssetRequireKYC(opts *bind.CallOpts, assetId *big.Int) (bool, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "doesAssetRequireKYC", assetId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DoesAssetRequireKYC is a free data retrieval call binding the contract method 0xcd9be9e7.
//
// Solidity: function doesAssetRequireKYC(uint256 assetId) view returns(bool)
func (_RWACompliance *RWAComplianceSession) DoesAssetRequireKYC(assetId *big.Int) (bool, error) {
	return _RWACompliance.Contract.DoesAssetRequireKYC(&_RWACompliance.CallOpts, assetId)
}

// DoesAssetRequireKYC is a free data retrieval call binding the contract method 0xcd9be9e7.
//
// Solidity: function doesAssetRequireKYC(uint256 assetId) view returns(bool)
func (_RWACompliance *RWAComplianceCallerSession) DoesAssetRequireKYC(assetId *big.Int) (bool, error) {
	return _RWACompliance.Contract.DoesAssetRequireKYC(&_RWACompliance.CallOpts, assetId)
}

// GetAssetInvestorCount is a free data retrieval call binding the contract method 0xca444eda.
//
// Solidity: function getAssetInvestorCount(uint256 assetId) view returns(uint256)
func (_RWACompliance *RWAComplianceCaller) GetAssetInvestorCount(opts *bind.CallOpts, assetId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "getAssetInvestorCount", assetId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAssetInvestorCount is a free data retrieval call binding the contract method 0xca444eda.
//
// Solidity: function getAssetInvestorCount(uint256 assetId) view returns(uint256)
func (_RWACompliance *RWAComplianceSession) GetAssetInvestorCount(assetId *big.Int) (*big.Int, error) {
	return _RWACompliance.Contract.GetAssetInvestorCount(&_RWACompliance.CallOpts, assetId)
}

// GetAssetInvestorCount is a free data retrieval call binding the contract method 0xca444eda.
//
// Solidity: function getAssetInvestorCount(uint256 assetId) view returns(uint256)
func (_RWACompliance *RWAComplianceCallerSession) GetAssetInvestorCount(assetId *big.Int) (*big.Int, error) {
	return _RWACompliance.Contract.GetAssetInvestorCount(&_RWACompliance.CallOpts, assetId)
}

// GetAssetMinTier is a free data retrieval call binding the contract method 0x74f205c4.
//
// Solidity: function getAssetMinTier(uint256 assetId) view returns(uint8)
func (_RWACompliance *RWAComplianceCaller) GetAssetMinTier(opts *bind.CallOpts, assetId *big.Int) (uint8, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "getAssetMinTier", assetId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetAssetMinTier is a free data retrieval call binding the contract method 0x74f205c4.
//
// Solidity: function getAssetMinTier(uint256 assetId) view returns(uint8)
func (_RWACompliance *RWAComplianceSession) GetAssetMinTier(assetId *big.Int) (uint8, error) {
	return _RWACompliance.Contract.GetAssetMinTier(&_RWACompliance.CallOpts, assetId)
}

// GetAssetMinTier is a free data retrieval call binding the contract method 0x74f205c4.
//
// Solidity: function getAssetMinTier(uint256 assetId) view returns(uint8)
func (_RWACompliance *RWAComplianceCallerSession) GetAssetMinTier(assetId *big.Int) (uint8, error) {
	return _RWACompliance.Contract.GetAssetMinTier(&_RWACompliance.CallOpts, assetId)
}

// GetInvestorProfile is a free data retrieval call binding the contract method 0xbcad0083.
//
// Solidity: function getInvestorProfile(address investor) view returns((uint8,uint8,uint256,uint256,string,bytes32,address))
func (_RWACompliance *RWAComplianceCaller) GetInvestorProfile(opts *bind.CallOpts, investor common.Address) (IRWAComplianceInvestorProfile, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "getInvestorProfile", investor)

	if err != nil {
		return *new(IRWAComplianceInvestorProfile), err
	}

	out0 := *abi.ConvertType(out[0], new(IRWAComplianceInvestorProfile)).(*IRWAComplianceInvestorProfile)

	return out0, err

}

// GetInvestorProfile is a free data retrieval call binding the contract method 0xbcad0083.
//
// Solidity: function getInvestorProfile(address investor) view returns((uint8,uint8,uint256,uint256,string,bytes32,address))
func (_RWACompliance *RWAComplianceSession) GetInvestorProfile(investor common.Address) (IRWAComplianceInvestorProfile, error) {
	return _RWACompliance.Contract.GetInvestorProfile(&_RWACompliance.CallOpts, investor)
}

// GetInvestorProfile is a free data retrieval call binding the contract method 0xbcad0083.
//
// Solidity: function getInvestorProfile(address investor) view returns((uint8,uint8,uint256,uint256,string,bytes32,address))
func (_RWACompliance *RWAComplianceCallerSession) GetInvestorProfile(investor common.Address) (IRWAComplianceInvestorProfile, error) {
	return _RWACompliance.Contract.GetInvestorProfile(&_RWACompliance.CallOpts, investor)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWACompliance *RWAComplianceCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWACompliance *RWAComplianceSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RWACompliance.Contract.GetRoleAdmin(&_RWACompliance.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWACompliance *RWAComplianceCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RWACompliance.Contract.GetRoleAdmin(&_RWACompliance.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWACompliance *RWAComplianceCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWACompliance *RWAComplianceSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RWACompliance.Contract.HasRole(&_RWACompliance.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWACompliance *RWAComplianceCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RWACompliance.Contract.HasRole(&_RWACompliance.CallOpts, role, account)
}

// IsInvestorVerified is a free data retrieval call binding the contract method 0xcfccef12.
//
// Solidity: function isInvestorVerified(address investor) view returns(bool)
func (_RWACompliance *RWAComplianceCaller) IsInvestorVerified(opts *bind.CallOpts, investor common.Address) (bool, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "isInvestorVerified", investor)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsInvestorVerified is a free data retrieval call binding the contract method 0xcfccef12.
//
// Solidity: function isInvestorVerified(address investor) view returns(bool)
func (_RWACompliance *RWAComplianceSession) IsInvestorVerified(investor common.Address) (bool, error) {
	return _RWACompliance.Contract.IsInvestorVerified(&_RWACompliance.CallOpts, investor)
}

// IsInvestorVerified is a free data retrieval call binding the contract method 0xcfccef12.
//
// Solidity: function isInvestorVerified(address investor) view returns(bool)
func (_RWACompliance *RWAComplianceCallerSession) IsInvestorVerified(investor common.Address) (bool, error) {
	return _RWACompliance.Contract.IsInvestorVerified(&_RWACompliance.CallOpts, investor)
}

// IsJurisdictionAllowed is a free data retrieval call binding the contract method 0x4171892c.
//
// Solidity: function isJurisdictionAllowed(uint256 assetId, string jurisdiction) view returns(bool)
func (_RWACompliance *RWAComplianceCaller) IsJurisdictionAllowed(opts *bind.CallOpts, assetId *big.Int, jurisdiction string) (bool, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "isJurisdictionAllowed", assetId, jurisdiction)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsJurisdictionAllowed is a free data retrieval call binding the contract method 0x4171892c.
//
// Solidity: function isJurisdictionAllowed(uint256 assetId, string jurisdiction) view returns(bool)
func (_RWACompliance *RWAComplianceSession) IsJurisdictionAllowed(assetId *big.Int, jurisdiction string) (bool, error) {
	return _RWACompliance.Contract.IsJurisdictionAllowed(&_RWACompliance.CallOpts, assetId, jurisdiction)
}

// IsJurisdictionAllowed is a free data retrieval call binding the contract method 0x4171892c.
//
// Solidity: function isJurisdictionAllowed(uint256 assetId, string jurisdiction) view returns(bool)
func (_RWACompliance *RWAComplianceCallerSession) IsJurisdictionAllowed(assetId *big.Int, jurisdiction string) (bool, error) {
	return _RWACompliance.Contract.IsJurisdictionAllowed(&_RWACompliance.CallOpts, assetId, jurisdiction)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWACompliance *RWAComplianceCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWACompliance *RWAComplianceSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RWACompliance.Contract.SupportsInterface(&_RWACompliance.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWACompliance *RWAComplianceCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RWACompliance.Contract.SupportsInterface(&_RWACompliance.CallOpts, interfaceId)
}

// TotalAssetsWithCompliance is a free data retrieval call binding the contract method 0xffd0f9dc.
//
// Solidity: function totalAssetsWithCompliance() view returns(uint256)
func (_RWACompliance *RWAComplianceCaller) TotalAssetsWithCompliance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "totalAssetsWithCompliance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAssetsWithCompliance is a free data retrieval call binding the contract method 0xffd0f9dc.
//
// Solidity: function totalAssetsWithCompliance() view returns(uint256)
func (_RWACompliance *RWAComplianceSession) TotalAssetsWithCompliance() (*big.Int, error) {
	return _RWACompliance.Contract.TotalAssetsWithCompliance(&_RWACompliance.CallOpts)
}

// TotalAssetsWithCompliance is a free data retrieval call binding the contract method 0xffd0f9dc.
//
// Solidity: function totalAssetsWithCompliance() view returns(uint256)
func (_RWACompliance *RWAComplianceCallerSession) TotalAssetsWithCompliance() (*big.Int, error) {
	return _RWACompliance.Contract.TotalAssetsWithCompliance(&_RWACompliance.CallOpts)
}

// TotalVerifiedInvestors is a free data retrieval call binding the contract method 0xf22bf8d0.
//
// Solidity: function totalVerifiedInvestors() view returns(uint256)
func (_RWACompliance *RWAComplianceCaller) TotalVerifiedInvestors(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWACompliance.contract.Call(opts, &out, "totalVerifiedInvestors")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalVerifiedInvestors is a free data retrieval call binding the contract method 0xf22bf8d0.
//
// Solidity: function totalVerifiedInvestors() view returns(uint256)
func (_RWACompliance *RWAComplianceSession) TotalVerifiedInvestors() (*big.Int, error) {
	return _RWACompliance.Contract.TotalVerifiedInvestors(&_RWACompliance.CallOpts)
}

// TotalVerifiedInvestors is a free data retrieval call binding the contract method 0xf22bf8d0.
//
// Solidity: function totalVerifiedInvestors() view returns(uint256)
func (_RWACompliance *RWAComplianceCallerSession) TotalVerifiedInvestors() (*big.Int, error) {
	return _RWACompliance.Contract.TotalVerifiedInvestors(&_RWACompliance.CallOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWACompliance *RWAComplianceTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWACompliance.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWACompliance *RWAComplianceSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWACompliance.Contract.GrantRole(&_RWACompliance.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWACompliance *RWAComplianceTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWACompliance.Contract.GrantRole(&_RWACompliance.TransactOpts, role, account)
}

// RecordInvestorParticipation is a paid mutator transaction binding the contract method 0xd23506b9.
//
// Solidity: function recordInvestorParticipation(uint256 assetId, address investor) returns()
func (_RWACompliance *RWAComplianceTransactor) RecordInvestorParticipation(opts *bind.TransactOpts, assetId *big.Int, investor common.Address) (*types.Transaction, error) {
	return _RWACompliance.contract.Transact(opts, "recordInvestorParticipation", assetId, investor)
}

// RecordInvestorParticipation is a paid mutator transaction binding the contract method 0xd23506b9.
//
// Solidity: function recordInvestorParticipation(uint256 assetId, address investor) returns()
func (_RWACompliance *RWAComplianceSession) RecordInvestorParticipation(assetId *big.Int, investor common.Address) (*types.Transaction, error) {
	return _RWACompliance.Contract.RecordInvestorParticipation(&_RWACompliance.TransactOpts, assetId, investor)
}

// RecordInvestorParticipation is a paid mutator transaction binding the contract method 0xd23506b9.
//
// Solidity: function recordInvestorParticipation(uint256 assetId, address investor) returns()
func (_RWACompliance *RWAComplianceTransactorSession) RecordInvestorParticipation(assetId *big.Int, investor common.Address) (*types.Transaction, error) {
	return _RWACompliance.Contract.RecordInvestorParticipation(&_RWACompliance.TransactOpts, assetId, investor)
}

// RemoveInvestorParticipation is a paid mutator transaction binding the contract method 0x6f04e96c.
//
// Solidity: function removeInvestorParticipation(uint256 assetId, address investor) returns()
func (_RWACompliance *RWAComplianceTransactor) RemoveInvestorParticipation(opts *bind.TransactOpts, assetId *big.Int, investor common.Address) (*types.Transaction, error) {
	return _RWACompliance.contract.Transact(opts, "removeInvestorParticipation", assetId, investor)
}

// RemoveInvestorParticipation is a paid mutator transaction binding the contract method 0x6f04e96c.
//
// Solidity: function removeInvestorParticipation(uint256 assetId, address investor) returns()
func (_RWACompliance *RWAComplianceSession) RemoveInvestorParticipation(assetId *big.Int, investor common.Address) (*types.Transaction, error) {
	return _RWACompliance.Contract.RemoveInvestorParticipation(&_RWACompliance.TransactOpts, assetId, investor)
}

// RemoveInvestorParticipation is a paid mutator transaction binding the contract method 0x6f04e96c.
//
// Solidity: function removeInvestorParticipation(uint256 assetId, address investor) returns()
func (_RWACompliance *RWAComplianceTransactorSession) RemoveInvestorParticipation(assetId *big.Int, investor common.Address) (*types.Transaction, error) {
	return _RWACompliance.Contract.RemoveInvestorParticipation(&_RWACompliance.TransactOpts, assetId, investor)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWACompliance *RWAComplianceTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWACompliance.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWACompliance *RWAComplianceSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWACompliance.Contract.RenounceRole(&_RWACompliance.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWACompliance *RWAComplianceTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWACompliance.Contract.RenounceRole(&_RWACompliance.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWACompliance *RWAComplianceTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWACompliance.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWACompliance *RWAComplianceSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWACompliance.Contract.RevokeRole(&_RWACompliance.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWACompliance *RWAComplianceTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWACompliance.Contract.RevokeRole(&_RWACompliance.TransactOpts, role, account)
}

// SetAssetCompliance is a paid mutator transaction binding the contract method 0x1e4ac0ab.
//
// Solidity: function setAssetCompliance(uint256 assetId, uint8 minTier, bool requiresKYC, uint256 minInvestmentAmount, uint256 maxInvestors) returns()
func (_RWACompliance *RWAComplianceTransactor) SetAssetCompliance(opts *bind.TransactOpts, assetId *big.Int, minTier uint8, requiresKYC bool, minInvestmentAmount *big.Int, maxInvestors *big.Int) (*types.Transaction, error) {
	return _RWACompliance.contract.Transact(opts, "setAssetCompliance", assetId, minTier, requiresKYC, minInvestmentAmount, maxInvestors)
}

// SetAssetCompliance is a paid mutator transaction binding the contract method 0x1e4ac0ab.
//
// Solidity: function setAssetCompliance(uint256 assetId, uint8 minTier, bool requiresKYC, uint256 minInvestmentAmount, uint256 maxInvestors) returns()
func (_RWACompliance *RWAComplianceSession) SetAssetCompliance(assetId *big.Int, minTier uint8, requiresKYC bool, minInvestmentAmount *big.Int, maxInvestors *big.Int) (*types.Transaction, error) {
	return _RWACompliance.Contract.SetAssetCompliance(&_RWACompliance.TransactOpts, assetId, minTier, requiresKYC, minInvestmentAmount, maxInvestors)
}

// SetAssetCompliance is a paid mutator transaction binding the contract method 0x1e4ac0ab.
//
// Solidity: function setAssetCompliance(uint256 assetId, uint8 minTier, bool requiresKYC, uint256 minInvestmentAmount, uint256 maxInvestors) returns()
func (_RWACompliance *RWAComplianceTransactorSession) SetAssetCompliance(assetId *big.Int, minTier uint8, requiresKYC bool, minInvestmentAmount *big.Int, maxInvestors *big.Int) (*types.Transaction, error) {
	return _RWACompliance.Contract.SetAssetCompliance(&_RWACompliance.TransactOpts, assetId, minTier, requiresKYC, minInvestmentAmount, maxInvestors)
}

// UpdateInvestorStatus is a paid mutator transaction binding the contract method 0x06aadaf1.
//
// Solidity: function updateInvestorStatus(address investor, uint8 newStatus) returns()
func (_RWACompliance *RWAComplianceTransactor) UpdateInvestorStatus(opts *bind.TransactOpts, investor common.Address, newStatus uint8) (*types.Transaction, error) {
	return _RWACompliance.contract.Transact(opts, "updateInvestorStatus", investor, newStatus)
}

// UpdateInvestorStatus is a paid mutator transaction binding the contract method 0x06aadaf1.
//
// Solidity: function updateInvestorStatus(address investor, uint8 newStatus) returns()
func (_RWACompliance *RWAComplianceSession) UpdateInvestorStatus(investor common.Address, newStatus uint8) (*types.Transaction, error) {
	return _RWACompliance.Contract.UpdateInvestorStatus(&_RWACompliance.TransactOpts, investor, newStatus)
}

// UpdateInvestorStatus is a paid mutator transaction binding the contract method 0x06aadaf1.
//
// Solidity: function updateInvestorStatus(address investor, uint8 newStatus) returns()
func (_RWACompliance *RWAComplianceTransactorSession) UpdateInvestorStatus(investor common.Address, newStatus uint8) (*types.Transaction, error) {
	return _RWACompliance.Contract.UpdateInvestorStatus(&_RWACompliance.TransactOpts, investor, newStatus)
}

// UpdateJurisdictionRestriction is a paid mutator transaction binding the contract method 0x43c1e23a.
//
// Solidity: function updateJurisdictionRestriction(uint256 assetId, string jurisdiction, bool allowed) returns()
func (_RWACompliance *RWAComplianceTransactor) UpdateJurisdictionRestriction(opts *bind.TransactOpts, assetId *big.Int, jurisdiction string, allowed bool) (*types.Transaction, error) {
	return _RWACompliance.contract.Transact(opts, "updateJurisdictionRestriction", assetId, jurisdiction, allowed)
}

// UpdateJurisdictionRestriction is a paid mutator transaction binding the contract method 0x43c1e23a.
//
// Solidity: function updateJurisdictionRestriction(uint256 assetId, string jurisdiction, bool allowed) returns()
func (_RWACompliance *RWAComplianceSession) UpdateJurisdictionRestriction(assetId *big.Int, jurisdiction string, allowed bool) (*types.Transaction, error) {
	return _RWACompliance.Contract.UpdateJurisdictionRestriction(&_RWACompliance.TransactOpts, assetId, jurisdiction, allowed)
}

// UpdateJurisdictionRestriction is a paid mutator transaction binding the contract method 0x43c1e23a.
//
// Solidity: function updateJurisdictionRestriction(uint256 assetId, string jurisdiction, bool allowed) returns()
func (_RWACompliance *RWAComplianceTransactorSession) UpdateJurisdictionRestriction(assetId *big.Int, jurisdiction string, allowed bool) (*types.Transaction, error) {
	return _RWACompliance.Contract.UpdateJurisdictionRestriction(&_RWACompliance.TransactOpts, assetId, jurisdiction, allowed)
}

// VerifyInvestor is a paid mutator transaction binding the contract method 0x826642cb.
//
// Solidity: function verifyInvestor(address investor, uint8 tier, string jurisdiction, uint256 validityPeriod, bytes32 kycDocumentHash) returns()
func (_RWACompliance *RWAComplianceTransactor) VerifyInvestor(opts *bind.TransactOpts, investor common.Address, tier uint8, jurisdiction string, validityPeriod *big.Int, kycDocumentHash [32]byte) (*types.Transaction, error) {
	return _RWACompliance.contract.Transact(opts, "verifyInvestor", investor, tier, jurisdiction, validityPeriod, kycDocumentHash)
}

// VerifyInvestor is a paid mutator transaction binding the contract method 0x826642cb.
//
// Solidity: function verifyInvestor(address investor, uint8 tier, string jurisdiction, uint256 validityPeriod, bytes32 kycDocumentHash) returns()
func (_RWACompliance *RWAComplianceSession) VerifyInvestor(investor common.Address, tier uint8, jurisdiction string, validityPeriod *big.Int, kycDocumentHash [32]byte) (*types.Transaction, error) {
	return _RWACompliance.Contract.VerifyInvestor(&_RWACompliance.TransactOpts, investor, tier, jurisdiction, validityPeriod, kycDocumentHash)
}

// VerifyInvestor is a paid mutator transaction binding the contract method 0x826642cb.
//
// Solidity: function verifyInvestor(address investor, uint8 tier, string jurisdiction, uint256 validityPeriod, bytes32 kycDocumentHash) returns()
func (_RWACompliance *RWAComplianceTransactorSession) VerifyInvestor(investor common.Address, tier uint8, jurisdiction string, validityPeriod *big.Int, kycDocumentHash [32]byte) (*types.Transaction, error) {
	return _RWACompliance.Contract.VerifyInvestor(&_RWACompliance.TransactOpts, investor, tier, jurisdiction, validityPeriod, kycDocumentHash)
}

// RWAComplianceAssetComplianceSetIterator is returned from FilterAssetComplianceSet and is used to iterate over the raw logs and unpacked data for AssetComplianceSet events raised by the RWACompliance contract.
type RWAComplianceAssetComplianceSetIterator struct {
	Event *RWAComplianceAssetComplianceSet // Event containing the contract specifics and raw log

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
func (it *RWAComplianceAssetComplianceSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAComplianceAssetComplianceSet)
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
		it.Event = new(RWAComplianceAssetComplianceSet)
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
func (it *RWAComplianceAssetComplianceSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAComplianceAssetComplianceSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAComplianceAssetComplianceSet represents a AssetComplianceSet event raised by the RWACompliance contract.
type RWAComplianceAssetComplianceSet struct {
	AssetId     *big.Int
	MinTier     uint8
	RequiresKYC bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssetComplianceSet is a free log retrieval operation binding the contract event 0xbe20a26a4383a126e361a15c2af181ee462fcd824c5eb87609a20caf62eb99aa.
//
// Solidity: event AssetComplianceSet(uint256 indexed assetId, uint8 minTier, bool requiresKYC)
func (_RWACompliance *RWAComplianceFilterer) FilterAssetComplianceSet(opts *bind.FilterOpts, assetId []*big.Int) (*RWAComplianceAssetComplianceSetIterator, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWACompliance.contract.FilterLogs(opts, "AssetComplianceSet", assetIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAComplianceAssetComplianceSetIterator{contract: _RWACompliance.contract, event: "AssetComplianceSet", logs: logs, sub: sub}, nil
}

// WatchAssetComplianceSet is a free log subscription operation binding the contract event 0xbe20a26a4383a126e361a15c2af181ee462fcd824c5eb87609a20caf62eb99aa.
//
// Solidity: event AssetComplianceSet(uint256 indexed assetId, uint8 minTier, bool requiresKYC)
func (_RWACompliance *RWAComplianceFilterer) WatchAssetComplianceSet(opts *bind.WatchOpts, sink chan<- *RWAComplianceAssetComplianceSet, assetId []*big.Int) (event.Subscription, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWACompliance.contract.WatchLogs(opts, "AssetComplianceSet", assetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAComplianceAssetComplianceSet)
				if err := _RWACompliance.contract.UnpackLog(event, "AssetComplianceSet", log); err != nil {
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

// ParseAssetComplianceSet is a log parse operation binding the contract event 0xbe20a26a4383a126e361a15c2af181ee462fcd824c5eb87609a20caf62eb99aa.
//
// Solidity: event AssetComplianceSet(uint256 indexed assetId, uint8 minTier, bool requiresKYC)
func (_RWACompliance *RWAComplianceFilterer) ParseAssetComplianceSet(log types.Log) (*RWAComplianceAssetComplianceSet, error) {
	event := new(RWAComplianceAssetComplianceSet)
	if err := _RWACompliance.contract.UnpackLog(event, "AssetComplianceSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAComplianceComplianceViolationIterator is returned from FilterComplianceViolation and is used to iterate over the raw logs and unpacked data for ComplianceViolation events raised by the RWACompliance contract.
type RWAComplianceComplianceViolationIterator struct {
	Event *RWAComplianceComplianceViolation // Event containing the contract specifics and raw log

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
func (it *RWAComplianceComplianceViolationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAComplianceComplianceViolation)
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
		it.Event = new(RWAComplianceComplianceViolation)
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
func (it *RWAComplianceComplianceViolationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAComplianceComplianceViolationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAComplianceComplianceViolation represents a ComplianceViolation event raised by the RWACompliance contract.
type RWAComplianceComplianceViolation struct {
	Investor common.Address
	AssetId  *big.Int
	Reason   string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterComplianceViolation is a free log retrieval operation binding the contract event 0xa1aefd70770f07d9883de74fb449e58b34fa313dcf0c452aef82a7081d848400.
//
// Solidity: event ComplianceViolation(address indexed investor, uint256 indexed assetId, string reason)
func (_RWACompliance *RWAComplianceFilterer) FilterComplianceViolation(opts *bind.FilterOpts, investor []common.Address, assetId []*big.Int) (*RWAComplianceComplianceViolationIterator, error) {

	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}
	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWACompliance.contract.FilterLogs(opts, "ComplianceViolation", investorRule, assetIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAComplianceComplianceViolationIterator{contract: _RWACompliance.contract, event: "ComplianceViolation", logs: logs, sub: sub}, nil
}

// WatchComplianceViolation is a free log subscription operation binding the contract event 0xa1aefd70770f07d9883de74fb449e58b34fa313dcf0c452aef82a7081d848400.
//
// Solidity: event ComplianceViolation(address indexed investor, uint256 indexed assetId, string reason)
func (_RWACompliance *RWAComplianceFilterer) WatchComplianceViolation(opts *bind.WatchOpts, sink chan<- *RWAComplianceComplianceViolation, investor []common.Address, assetId []*big.Int) (event.Subscription, error) {

	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}
	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWACompliance.contract.WatchLogs(opts, "ComplianceViolation", investorRule, assetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAComplianceComplianceViolation)
				if err := _RWACompliance.contract.UnpackLog(event, "ComplianceViolation", log); err != nil {
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

// ParseComplianceViolation is a log parse operation binding the contract event 0xa1aefd70770f07d9883de74fb449e58b34fa313dcf0c452aef82a7081d848400.
//
// Solidity: event ComplianceViolation(address indexed investor, uint256 indexed assetId, string reason)
func (_RWACompliance *RWAComplianceFilterer) ParseComplianceViolation(log types.Log) (*RWAComplianceComplianceViolation, error) {
	event := new(RWAComplianceComplianceViolation)
	if err := _RWACompliance.contract.UnpackLog(event, "ComplianceViolation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAComplianceInvestorStatusChangedIterator is returned from FilterInvestorStatusChanged and is used to iterate over the raw logs and unpacked data for InvestorStatusChanged events raised by the RWACompliance contract.
type RWAComplianceInvestorStatusChangedIterator struct {
	Event *RWAComplianceInvestorStatusChanged // Event containing the contract specifics and raw log

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
func (it *RWAComplianceInvestorStatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAComplianceInvestorStatusChanged)
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
		it.Event = new(RWAComplianceInvestorStatusChanged)
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
func (it *RWAComplianceInvestorStatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAComplianceInvestorStatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAComplianceInvestorStatusChanged represents a InvestorStatusChanged event raised by the RWACompliance contract.
type RWAComplianceInvestorStatusChanged struct {
	Investor  common.Address
	OldStatus uint8
	NewStatus uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterInvestorStatusChanged is a free log retrieval operation binding the contract event 0xe92b9476c821a6af7ae217fe20447a76cde894b52ffc32c531818e17f3c8a026.
//
// Solidity: event InvestorStatusChanged(address indexed investor, uint8 oldStatus, uint8 newStatus)
func (_RWACompliance *RWAComplianceFilterer) FilterInvestorStatusChanged(opts *bind.FilterOpts, investor []common.Address) (*RWAComplianceInvestorStatusChangedIterator, error) {

	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}

	logs, sub, err := _RWACompliance.contract.FilterLogs(opts, "InvestorStatusChanged", investorRule)
	if err != nil {
		return nil, err
	}
	return &RWAComplianceInvestorStatusChangedIterator{contract: _RWACompliance.contract, event: "InvestorStatusChanged", logs: logs, sub: sub}, nil
}

// WatchInvestorStatusChanged is a free log subscription operation binding the contract event 0xe92b9476c821a6af7ae217fe20447a76cde894b52ffc32c531818e17f3c8a026.
//
// Solidity: event InvestorStatusChanged(address indexed investor, uint8 oldStatus, uint8 newStatus)
func (_RWACompliance *RWAComplianceFilterer) WatchInvestorStatusChanged(opts *bind.WatchOpts, sink chan<- *RWAComplianceInvestorStatusChanged, investor []common.Address) (event.Subscription, error) {

	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}

	logs, sub, err := _RWACompliance.contract.WatchLogs(opts, "InvestorStatusChanged", investorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAComplianceInvestorStatusChanged)
				if err := _RWACompliance.contract.UnpackLog(event, "InvestorStatusChanged", log); err != nil {
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

// ParseInvestorStatusChanged is a log parse operation binding the contract event 0xe92b9476c821a6af7ae217fe20447a76cde894b52ffc32c531818e17f3c8a026.
//
// Solidity: event InvestorStatusChanged(address indexed investor, uint8 oldStatus, uint8 newStatus)
func (_RWACompliance *RWAComplianceFilterer) ParseInvestorStatusChanged(log types.Log) (*RWAComplianceInvestorStatusChanged, error) {
	event := new(RWAComplianceInvestorStatusChanged)
	if err := _RWACompliance.contract.UnpackLog(event, "InvestorStatusChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAComplianceInvestorVerifiedIterator is returned from FilterInvestorVerified and is used to iterate over the raw logs and unpacked data for InvestorVerified events raised by the RWACompliance contract.
type RWAComplianceInvestorVerifiedIterator struct {
	Event *RWAComplianceInvestorVerified // Event containing the contract specifics and raw log

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
func (it *RWAComplianceInvestorVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAComplianceInvestorVerified)
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
		it.Event = new(RWAComplianceInvestorVerified)
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
func (it *RWAComplianceInvestorVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAComplianceInvestorVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAComplianceInvestorVerified represents a InvestorVerified event raised by the RWACompliance contract.
type RWAComplianceInvestorVerified struct {
	Investor  common.Address
	Status    uint8
	Tier      uint8
	ExpiresAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterInvestorVerified is a free log retrieval operation binding the contract event 0x1b9207a4cac66996dd26eb5db81506e1bf928134d9c3c55b384689d604030a4e.
//
// Solidity: event InvestorVerified(address indexed investor, uint8 status, uint8 tier, uint256 expiresAt)
func (_RWACompliance *RWAComplianceFilterer) FilterInvestorVerified(opts *bind.FilterOpts, investor []common.Address) (*RWAComplianceInvestorVerifiedIterator, error) {

	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}

	logs, sub, err := _RWACompliance.contract.FilterLogs(opts, "InvestorVerified", investorRule)
	if err != nil {
		return nil, err
	}
	return &RWAComplianceInvestorVerifiedIterator{contract: _RWACompliance.contract, event: "InvestorVerified", logs: logs, sub: sub}, nil
}

// WatchInvestorVerified is a free log subscription operation binding the contract event 0x1b9207a4cac66996dd26eb5db81506e1bf928134d9c3c55b384689d604030a4e.
//
// Solidity: event InvestorVerified(address indexed investor, uint8 status, uint8 tier, uint256 expiresAt)
func (_RWACompliance *RWAComplianceFilterer) WatchInvestorVerified(opts *bind.WatchOpts, sink chan<- *RWAComplianceInvestorVerified, investor []common.Address) (event.Subscription, error) {

	var investorRule []interface{}
	for _, investorItem := range investor {
		investorRule = append(investorRule, investorItem)
	}

	logs, sub, err := _RWACompliance.contract.WatchLogs(opts, "InvestorVerified", investorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAComplianceInvestorVerified)
				if err := _RWACompliance.contract.UnpackLog(event, "InvestorVerified", log); err != nil {
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

// ParseInvestorVerified is a log parse operation binding the contract event 0x1b9207a4cac66996dd26eb5db81506e1bf928134d9c3c55b384689d604030a4e.
//
// Solidity: event InvestorVerified(address indexed investor, uint8 status, uint8 tier, uint256 expiresAt)
func (_RWACompliance *RWAComplianceFilterer) ParseInvestorVerified(log types.Log) (*RWAComplianceInvestorVerified, error) {
	event := new(RWAComplianceInvestorVerified)
	if err := _RWACompliance.contract.UnpackLog(event, "InvestorVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAComplianceJurisdictionUpdatedIterator is returned from FilterJurisdictionUpdated and is used to iterate over the raw logs and unpacked data for JurisdictionUpdated events raised by the RWACompliance contract.
type RWAComplianceJurisdictionUpdatedIterator struct {
	Event *RWAComplianceJurisdictionUpdated // Event containing the contract specifics and raw log

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
func (it *RWAComplianceJurisdictionUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAComplianceJurisdictionUpdated)
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
		it.Event = new(RWAComplianceJurisdictionUpdated)
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
func (it *RWAComplianceJurisdictionUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAComplianceJurisdictionUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAComplianceJurisdictionUpdated represents a JurisdictionUpdated event raised by the RWACompliance contract.
type RWAComplianceJurisdictionUpdated struct {
	AssetId      *big.Int
	Jurisdiction string
	Allowed      bool
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterJurisdictionUpdated is a free log retrieval operation binding the contract event 0xc1d095a8ade79740c62541f6a4715710774ca4de188904762c842fe18e2fae4c.
//
// Solidity: event JurisdictionUpdated(uint256 indexed assetId, string jurisdiction, bool allowed)
func (_RWACompliance *RWAComplianceFilterer) FilterJurisdictionUpdated(opts *bind.FilterOpts, assetId []*big.Int) (*RWAComplianceJurisdictionUpdatedIterator, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWACompliance.contract.FilterLogs(opts, "JurisdictionUpdated", assetIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAComplianceJurisdictionUpdatedIterator{contract: _RWACompliance.contract, event: "JurisdictionUpdated", logs: logs, sub: sub}, nil
}

// WatchJurisdictionUpdated is a free log subscription operation binding the contract event 0xc1d095a8ade79740c62541f6a4715710774ca4de188904762c842fe18e2fae4c.
//
// Solidity: event JurisdictionUpdated(uint256 indexed assetId, string jurisdiction, bool allowed)
func (_RWACompliance *RWAComplianceFilterer) WatchJurisdictionUpdated(opts *bind.WatchOpts, sink chan<- *RWAComplianceJurisdictionUpdated, assetId []*big.Int) (event.Subscription, error) {

	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWACompliance.contract.WatchLogs(opts, "JurisdictionUpdated", assetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAComplianceJurisdictionUpdated)
				if err := _RWACompliance.contract.UnpackLog(event, "JurisdictionUpdated", log); err != nil {
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

// ParseJurisdictionUpdated is a log parse operation binding the contract event 0xc1d095a8ade79740c62541f6a4715710774ca4de188904762c842fe18e2fae4c.
//
// Solidity: event JurisdictionUpdated(uint256 indexed assetId, string jurisdiction, bool allowed)
func (_RWACompliance *RWAComplianceFilterer) ParseJurisdictionUpdated(log types.Log) (*RWAComplianceJurisdictionUpdated, error) {
	event := new(RWAComplianceJurisdictionUpdated)
	if err := _RWACompliance.contract.UnpackLog(event, "JurisdictionUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAComplianceRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the RWACompliance contract.
type RWAComplianceRoleAdminChangedIterator struct {
	Event *RWAComplianceRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *RWAComplianceRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAComplianceRoleAdminChanged)
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
		it.Event = new(RWAComplianceRoleAdminChanged)
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
func (it *RWAComplianceRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAComplianceRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAComplianceRoleAdminChanged represents a RoleAdminChanged event raised by the RWACompliance contract.
type RWAComplianceRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RWACompliance *RWAComplianceFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*RWAComplianceRoleAdminChangedIterator, error) {

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

	logs, sub, err := _RWACompliance.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &RWAComplianceRoleAdminChangedIterator{contract: _RWACompliance.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RWACompliance *RWAComplianceFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *RWAComplianceRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _RWACompliance.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAComplianceRoleAdminChanged)
				if err := _RWACompliance.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_RWACompliance *RWAComplianceFilterer) ParseRoleAdminChanged(log types.Log) (*RWAComplianceRoleAdminChanged, error) {
	event := new(RWAComplianceRoleAdminChanged)
	if err := _RWACompliance.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAComplianceRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the RWACompliance contract.
type RWAComplianceRoleGrantedIterator struct {
	Event *RWAComplianceRoleGranted // Event containing the contract specifics and raw log

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
func (it *RWAComplianceRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAComplianceRoleGranted)
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
		it.Event = new(RWAComplianceRoleGranted)
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
func (it *RWAComplianceRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAComplianceRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAComplianceRoleGranted represents a RoleGranted event raised by the RWACompliance contract.
type RWAComplianceRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWACompliance *RWAComplianceFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RWAComplianceRoleGrantedIterator, error) {

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

	logs, sub, err := _RWACompliance.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RWAComplianceRoleGrantedIterator{contract: _RWACompliance.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWACompliance *RWAComplianceFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *RWAComplianceRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RWACompliance.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAComplianceRoleGranted)
				if err := _RWACompliance.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_RWACompliance *RWAComplianceFilterer) ParseRoleGranted(log types.Log) (*RWAComplianceRoleGranted, error) {
	event := new(RWAComplianceRoleGranted)
	if err := _RWACompliance.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAComplianceRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the RWACompliance contract.
type RWAComplianceRoleRevokedIterator struct {
	Event *RWAComplianceRoleRevoked // Event containing the contract specifics and raw log

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
func (it *RWAComplianceRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAComplianceRoleRevoked)
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
		it.Event = new(RWAComplianceRoleRevoked)
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
func (it *RWAComplianceRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAComplianceRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAComplianceRoleRevoked represents a RoleRevoked event raised by the RWACompliance contract.
type RWAComplianceRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWACompliance *RWAComplianceFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RWAComplianceRoleRevokedIterator, error) {

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

	logs, sub, err := _RWACompliance.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RWAComplianceRoleRevokedIterator{contract: _RWACompliance.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWACompliance *RWAComplianceFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *RWAComplianceRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RWACompliance.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAComplianceRoleRevoked)
				if err := _RWACompliance.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_RWACompliance *RWAComplianceFilterer) ParseRoleRevoked(log types.Log) (*RWAComplianceRoleRevoked, error) {
	event := new(RWAComplianceRoleRevoked)
	if err := _RWACompliance.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
