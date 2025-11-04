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

// RWAYieldDistributorClaim is an auto generated low-level Go binding around an user-defined struct.
type RWAYieldDistributorClaim struct {
	Amount    *big.Int
	Claimed   bool
	ClaimedAt *big.Int
}

// RWAYieldDistributorDistribution is an auto generated low-level Go binding around an user-defined struct.
type RWAYieldDistributorDistribution struct {
	DistributionId *big.Int
	AssetId        *big.Int
	PaymentToken   common.Address
	TotalAmount    *big.Int
	TotalSupply    *big.Int
	DistributedAt  *big.Int
	ClaimDeadline  *big.Int
	TotalClaimed   *big.Int
	Finalized      bool
}

// RWAYieldDistributorMetaData contains all meta data concerning the RWAYieldDistributor contract.
var RWAYieldDistributorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"distributionId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalClaimed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unclaimed\",\"type\":\"uint256\"}],\"name\":\"DistributionFinalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"distributionId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"UnclaimedYieldReclaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"distributionId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"YieldClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"distributionId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"paymentToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"claimDeadline\",\"type\":\"uint256\"}],\"name\":\"YieldDeposited\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_CLAIM_PERIOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DISTRIBUTOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"assetFactory\",\"outputs\":[{\"internalType\":\"contractIRWAAsset\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"distributionIds\",\"type\":\"uint256[]\"}],\"name\":\"batchClaimYield\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"distributionId\",\"type\":\"uint256\"}],\"name\":\"claimYield\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"paymentToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"claimPeriod\",\"type\":\"uint256\"}],\"name\":\"depositYield\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"distributionId\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"emergencyWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"distributionId\",\"type\":\"uint256\"}],\"name\":\"finalizeDistribution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getAssetDistributions\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"distributionId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getClaimableAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"distributionId\",\"type\":\"uint256\"}],\"name\":\"getDistribution\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"distributionId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"paymentToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"distributedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"claimDeadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalClaimed\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"finalized\",\"type\":\"bool\"}],\"internalType\":\"structRWAYieldDistributor.Distribution\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getTotalClaimable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalClaimable\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"distributionId\",\"type\":\"uint256\"}],\"name\":\"getUnclaimedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"distributionId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserClaim\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"claimed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"claimedAt\",\"type\":\"uint256\"}],\"internalType\":\"structRWAYieldDistributor.Claim\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalDistributions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"totalYieldByToken\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalYieldDistributed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// RWAYieldDistributorABI is the input ABI used to generate the binding from.
// Deprecated: Use RWAYieldDistributorMetaData.ABI instead.
var RWAYieldDistributorABI = RWAYieldDistributorMetaData.ABI

// RWAYieldDistributor is an auto generated Go binding around an Ethereum contract.
type RWAYieldDistributor struct {
	RWAYieldDistributorCaller     // Read-only binding to the contract
	RWAYieldDistributorTransactor // Write-only binding to the contract
	RWAYieldDistributorFilterer   // Log filterer for contract events
}

// RWAYieldDistributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type RWAYieldDistributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAYieldDistributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RWAYieldDistributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAYieldDistributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RWAYieldDistributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAYieldDistributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RWAYieldDistributorSession struct {
	Contract     *RWAYieldDistributor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// RWAYieldDistributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RWAYieldDistributorCallerSession struct {
	Contract *RWAYieldDistributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// RWAYieldDistributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RWAYieldDistributorTransactorSession struct {
	Contract     *RWAYieldDistributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// RWAYieldDistributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type RWAYieldDistributorRaw struct {
	Contract *RWAYieldDistributor // Generic contract binding to access the raw methods on
}

// RWAYieldDistributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RWAYieldDistributorCallerRaw struct {
	Contract *RWAYieldDistributorCaller // Generic read-only contract binding to access the raw methods on
}

// RWAYieldDistributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RWAYieldDistributorTransactorRaw struct {
	Contract *RWAYieldDistributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRWAYieldDistributor creates a new instance of RWAYieldDistributor, bound to a specific deployed contract.
func NewRWAYieldDistributor(address common.Address, backend bind.ContractBackend) (*RWAYieldDistributor, error) {
	contract, err := bindRWAYieldDistributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributor{RWAYieldDistributorCaller: RWAYieldDistributorCaller{contract: contract}, RWAYieldDistributorTransactor: RWAYieldDistributorTransactor{contract: contract}, RWAYieldDistributorFilterer: RWAYieldDistributorFilterer{contract: contract}}, nil
}

// NewRWAYieldDistributorCaller creates a new read-only instance of RWAYieldDistributor, bound to a specific deployed contract.
func NewRWAYieldDistributorCaller(address common.Address, caller bind.ContractCaller) (*RWAYieldDistributorCaller, error) {
	contract, err := bindRWAYieldDistributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributorCaller{contract: contract}, nil
}

// NewRWAYieldDistributorTransactor creates a new write-only instance of RWAYieldDistributor, bound to a specific deployed contract.
func NewRWAYieldDistributorTransactor(address common.Address, transactor bind.ContractTransactor) (*RWAYieldDistributorTransactor, error) {
	contract, err := bindRWAYieldDistributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributorTransactor{contract: contract}, nil
}

// NewRWAYieldDistributorFilterer creates a new log filterer instance of RWAYieldDistributor, bound to a specific deployed contract.
func NewRWAYieldDistributorFilterer(address common.Address, filterer bind.ContractFilterer) (*RWAYieldDistributorFilterer, error) {
	contract, err := bindRWAYieldDistributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributorFilterer{contract: contract}, nil
}

// bindRWAYieldDistributor binds a generic wrapper to an already deployed contract.
func bindRWAYieldDistributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RWAYieldDistributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RWAYieldDistributor *RWAYieldDistributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RWAYieldDistributor.Contract.RWAYieldDistributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RWAYieldDistributor *RWAYieldDistributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.RWAYieldDistributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RWAYieldDistributor *RWAYieldDistributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.RWAYieldDistributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RWAYieldDistributor *RWAYieldDistributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RWAYieldDistributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RWAYieldDistributor *RWAYieldDistributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RWAYieldDistributor *RWAYieldDistributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAYieldDistributor *RWAYieldDistributorSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RWAYieldDistributor.Contract.DEFAULTADMINROLE(&_RWAYieldDistributor.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RWAYieldDistributor.Contract.DEFAULTADMINROLE(&_RWAYieldDistributor.CallOpts)
}

// DEFAULTCLAIMPERIOD is a free data retrieval call binding the contract method 0x4c1ab2f4.
//
// Solidity: function DEFAULT_CLAIM_PERIOD() view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) DEFAULTCLAIMPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "DEFAULT_CLAIM_PERIOD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEFAULTCLAIMPERIOD is a free data retrieval call binding the contract method 0x4c1ab2f4.
//
// Solidity: function DEFAULT_CLAIM_PERIOD() view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorSession) DEFAULTCLAIMPERIOD() (*big.Int, error) {
	return _RWAYieldDistributor.Contract.DEFAULTCLAIMPERIOD(&_RWAYieldDistributor.CallOpts)
}

// DEFAULTCLAIMPERIOD is a free data retrieval call binding the contract method 0x4c1ab2f4.
//
// Solidity: function DEFAULT_CLAIM_PERIOD() view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) DEFAULTCLAIMPERIOD() (*big.Int, error) {
	return _RWAYieldDistributor.Contract.DEFAULTCLAIMPERIOD(&_RWAYieldDistributor.CallOpts)
}

// DISTRIBUTORROLE is a free data retrieval call binding the contract method 0xf0bd87cc.
//
// Solidity: function DISTRIBUTOR_ROLE() view returns(bytes32)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) DISTRIBUTORROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "DISTRIBUTOR_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DISTRIBUTORROLE is a free data retrieval call binding the contract method 0xf0bd87cc.
//
// Solidity: function DISTRIBUTOR_ROLE() view returns(bytes32)
func (_RWAYieldDistributor *RWAYieldDistributorSession) DISTRIBUTORROLE() ([32]byte, error) {
	return _RWAYieldDistributor.Contract.DISTRIBUTORROLE(&_RWAYieldDistributor.CallOpts)
}

// DISTRIBUTORROLE is a free data retrieval call binding the contract method 0xf0bd87cc.
//
// Solidity: function DISTRIBUTOR_ROLE() view returns(bytes32)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) DISTRIBUTORROLE() ([32]byte, error) {
	return _RWAYieldDistributor.Contract.DISTRIBUTORROLE(&_RWAYieldDistributor.CallOpts)
}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) PRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorSession) PRECISION() (*big.Int, error) {
	return _RWAYieldDistributor.Contract.PRECISION(&_RWAYieldDistributor.CallOpts)
}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) PRECISION() (*big.Int, error) {
	return _RWAYieldDistributor.Contract.PRECISION(&_RWAYieldDistributor.CallOpts)
}

// AssetFactory is a free data retrieval call binding the contract method 0xc3acb4d1.
//
// Solidity: function assetFactory() view returns(address)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) AssetFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "assetFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AssetFactory is a free data retrieval call binding the contract method 0xc3acb4d1.
//
// Solidity: function assetFactory() view returns(address)
func (_RWAYieldDistributor *RWAYieldDistributorSession) AssetFactory() (common.Address, error) {
	return _RWAYieldDistributor.Contract.AssetFactory(&_RWAYieldDistributor.CallOpts)
}

// AssetFactory is a free data retrieval call binding the contract method 0xc3acb4d1.
//
// Solidity: function assetFactory() view returns(address)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) AssetFactory() (common.Address, error) {
	return _RWAYieldDistributor.Contract.AssetFactory(&_RWAYieldDistributor.CallOpts)
}

// GetAssetDistributions is a free data retrieval call binding the contract method 0x51a8d88b.
//
// Solidity: function getAssetDistributions(uint256 assetId) view returns(uint256[])
func (_RWAYieldDistributor *RWAYieldDistributorCaller) GetAssetDistributions(opts *bind.CallOpts, assetId *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "getAssetDistributions", assetId)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAssetDistributions is a free data retrieval call binding the contract method 0x51a8d88b.
//
// Solidity: function getAssetDistributions(uint256 assetId) view returns(uint256[])
func (_RWAYieldDistributor *RWAYieldDistributorSession) GetAssetDistributions(assetId *big.Int) ([]*big.Int, error) {
	return _RWAYieldDistributor.Contract.GetAssetDistributions(&_RWAYieldDistributor.CallOpts, assetId)
}

// GetAssetDistributions is a free data retrieval call binding the contract method 0x51a8d88b.
//
// Solidity: function getAssetDistributions(uint256 assetId) view returns(uint256[])
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) GetAssetDistributions(assetId *big.Int) ([]*big.Int, error) {
	return _RWAYieldDistributor.Contract.GetAssetDistributions(&_RWAYieldDistributor.CallOpts, assetId)
}

// GetClaimableAmount is a free data retrieval call binding the contract method 0x78c5195e.
//
// Solidity: function getClaimableAmount(uint256 distributionId, address user) view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) GetClaimableAmount(opts *bind.CallOpts, distributionId *big.Int, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "getClaimableAmount", distributionId, user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetClaimableAmount is a free data retrieval call binding the contract method 0x78c5195e.
//
// Solidity: function getClaimableAmount(uint256 distributionId, address user) view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorSession) GetClaimableAmount(distributionId *big.Int, user common.Address) (*big.Int, error) {
	return _RWAYieldDistributor.Contract.GetClaimableAmount(&_RWAYieldDistributor.CallOpts, distributionId, user)
}

// GetClaimableAmount is a free data retrieval call binding the contract method 0x78c5195e.
//
// Solidity: function getClaimableAmount(uint256 distributionId, address user) view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) GetClaimableAmount(distributionId *big.Int, user common.Address) (*big.Int, error) {
	return _RWAYieldDistributor.Contract.GetClaimableAmount(&_RWAYieldDistributor.CallOpts, distributionId, user)
}

// GetDistribution is a free data retrieval call binding the contract method 0x3b345a87.
//
// Solidity: function getDistribution(uint256 distributionId) view returns((uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bool))
func (_RWAYieldDistributor *RWAYieldDistributorCaller) GetDistribution(opts *bind.CallOpts, distributionId *big.Int) (RWAYieldDistributorDistribution, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "getDistribution", distributionId)

	if err != nil {
		return *new(RWAYieldDistributorDistribution), err
	}

	out0 := *abi.ConvertType(out[0], new(RWAYieldDistributorDistribution)).(*RWAYieldDistributorDistribution)

	return out0, err

}

// GetDistribution is a free data retrieval call binding the contract method 0x3b345a87.
//
// Solidity: function getDistribution(uint256 distributionId) view returns((uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bool))
func (_RWAYieldDistributor *RWAYieldDistributorSession) GetDistribution(distributionId *big.Int) (RWAYieldDistributorDistribution, error) {
	return _RWAYieldDistributor.Contract.GetDistribution(&_RWAYieldDistributor.CallOpts, distributionId)
}

// GetDistribution is a free data retrieval call binding the contract method 0x3b345a87.
//
// Solidity: function getDistribution(uint256 distributionId) view returns((uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bool))
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) GetDistribution(distributionId *big.Int) (RWAYieldDistributorDistribution, error) {
	return _RWAYieldDistributor.Contract.GetDistribution(&_RWAYieldDistributor.CallOpts, distributionId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAYieldDistributor *RWAYieldDistributorSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RWAYieldDistributor.Contract.GetRoleAdmin(&_RWAYieldDistributor.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RWAYieldDistributor.Contract.GetRoleAdmin(&_RWAYieldDistributor.CallOpts, role)
}

// GetTotalClaimable is a free data retrieval call binding the contract method 0x3d0f2200.
//
// Solidity: function getTotalClaimable(address user, uint256 assetId) view returns(uint256 totalClaimable)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) GetTotalClaimable(opts *bind.CallOpts, user common.Address, assetId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "getTotalClaimable", user, assetId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalClaimable is a free data retrieval call binding the contract method 0x3d0f2200.
//
// Solidity: function getTotalClaimable(address user, uint256 assetId) view returns(uint256 totalClaimable)
func (_RWAYieldDistributor *RWAYieldDistributorSession) GetTotalClaimable(user common.Address, assetId *big.Int) (*big.Int, error) {
	return _RWAYieldDistributor.Contract.GetTotalClaimable(&_RWAYieldDistributor.CallOpts, user, assetId)
}

// GetTotalClaimable is a free data retrieval call binding the contract method 0x3d0f2200.
//
// Solidity: function getTotalClaimable(address user, uint256 assetId) view returns(uint256 totalClaimable)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) GetTotalClaimable(user common.Address, assetId *big.Int) (*big.Int, error) {
	return _RWAYieldDistributor.Contract.GetTotalClaimable(&_RWAYieldDistributor.CallOpts, user, assetId)
}

// GetUnclaimedAmount is a free data retrieval call binding the contract method 0x18d474e3.
//
// Solidity: function getUnclaimedAmount(uint256 distributionId) view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) GetUnclaimedAmount(opts *bind.CallOpts, distributionId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "getUnclaimedAmount", distributionId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUnclaimedAmount is a free data retrieval call binding the contract method 0x18d474e3.
//
// Solidity: function getUnclaimedAmount(uint256 distributionId) view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorSession) GetUnclaimedAmount(distributionId *big.Int) (*big.Int, error) {
	return _RWAYieldDistributor.Contract.GetUnclaimedAmount(&_RWAYieldDistributor.CallOpts, distributionId)
}

// GetUnclaimedAmount is a free data retrieval call binding the contract method 0x18d474e3.
//
// Solidity: function getUnclaimedAmount(uint256 distributionId) view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) GetUnclaimedAmount(distributionId *big.Int) (*big.Int, error) {
	return _RWAYieldDistributor.Contract.GetUnclaimedAmount(&_RWAYieldDistributor.CallOpts, distributionId)
}

// GetUserClaim is a free data retrieval call binding the contract method 0x07cad02e.
//
// Solidity: function getUserClaim(uint256 distributionId, address user) view returns((uint256,bool,uint256))
func (_RWAYieldDistributor *RWAYieldDistributorCaller) GetUserClaim(opts *bind.CallOpts, distributionId *big.Int, user common.Address) (RWAYieldDistributorClaim, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "getUserClaim", distributionId, user)

	if err != nil {
		return *new(RWAYieldDistributorClaim), err
	}

	out0 := *abi.ConvertType(out[0], new(RWAYieldDistributorClaim)).(*RWAYieldDistributorClaim)

	return out0, err

}

// GetUserClaim is a free data retrieval call binding the contract method 0x07cad02e.
//
// Solidity: function getUserClaim(uint256 distributionId, address user) view returns((uint256,bool,uint256))
func (_RWAYieldDistributor *RWAYieldDistributorSession) GetUserClaim(distributionId *big.Int, user common.Address) (RWAYieldDistributorClaim, error) {
	return _RWAYieldDistributor.Contract.GetUserClaim(&_RWAYieldDistributor.CallOpts, distributionId, user)
}

// GetUserClaim is a free data retrieval call binding the contract method 0x07cad02e.
//
// Solidity: function getUserClaim(uint256 distributionId, address user) view returns((uint256,bool,uint256))
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) GetUserClaim(distributionId *big.Int, user common.Address) (RWAYieldDistributorClaim, error) {
	return _RWAYieldDistributor.Contract.GetUserClaim(&_RWAYieldDistributor.CallOpts, distributionId, user)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAYieldDistributor *RWAYieldDistributorSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RWAYieldDistributor.Contract.HasRole(&_RWAYieldDistributor.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RWAYieldDistributor.Contract.HasRole(&_RWAYieldDistributor.CallOpts, role, account)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RWAYieldDistributor *RWAYieldDistributorSession) Paused() (bool, error) {
	return _RWAYieldDistributor.Contract.Paused(&_RWAYieldDistributor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) Paused() (bool, error) {
	return _RWAYieldDistributor.Contract.Paused(&_RWAYieldDistributor.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAYieldDistributor *RWAYieldDistributorSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RWAYieldDistributor.Contract.SupportsInterface(&_RWAYieldDistributor.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RWAYieldDistributor.Contract.SupportsInterface(&_RWAYieldDistributor.CallOpts, interfaceId)
}

// TotalDistributions is a free data retrieval call binding the contract method 0x163db71b.
//
// Solidity: function totalDistributions() view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) TotalDistributions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "totalDistributions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalDistributions is a free data retrieval call binding the contract method 0x163db71b.
//
// Solidity: function totalDistributions() view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorSession) TotalDistributions() (*big.Int, error) {
	return _RWAYieldDistributor.Contract.TotalDistributions(&_RWAYieldDistributor.CallOpts)
}

// TotalDistributions is a free data retrieval call binding the contract method 0x163db71b.
//
// Solidity: function totalDistributions() view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) TotalDistributions() (*big.Int, error) {
	return _RWAYieldDistributor.Contract.TotalDistributions(&_RWAYieldDistributor.CallOpts)
}

// TotalYieldByToken is a free data retrieval call binding the contract method 0xe6a31f99.
//
// Solidity: function totalYieldByToken(address ) view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) TotalYieldByToken(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "totalYieldByToken", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalYieldByToken is a free data retrieval call binding the contract method 0xe6a31f99.
//
// Solidity: function totalYieldByToken(address ) view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorSession) TotalYieldByToken(arg0 common.Address) (*big.Int, error) {
	return _RWAYieldDistributor.Contract.TotalYieldByToken(&_RWAYieldDistributor.CallOpts, arg0)
}

// TotalYieldByToken is a free data retrieval call binding the contract method 0xe6a31f99.
//
// Solidity: function totalYieldByToken(address ) view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) TotalYieldByToken(arg0 common.Address) (*big.Int, error) {
	return _RWAYieldDistributor.Contract.TotalYieldByToken(&_RWAYieldDistributor.CallOpts, arg0)
}

// TotalYieldDistributed is a free data retrieval call binding the contract method 0x97206bd5.
//
// Solidity: function totalYieldDistributed() view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCaller) TotalYieldDistributed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAYieldDistributor.contract.Call(opts, &out, "totalYieldDistributed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalYieldDistributed is a free data retrieval call binding the contract method 0x97206bd5.
//
// Solidity: function totalYieldDistributed() view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorSession) TotalYieldDistributed() (*big.Int, error) {
	return _RWAYieldDistributor.Contract.TotalYieldDistributed(&_RWAYieldDistributor.CallOpts)
}

// TotalYieldDistributed is a free data retrieval call binding the contract method 0x97206bd5.
//
// Solidity: function totalYieldDistributed() view returns(uint256)
func (_RWAYieldDistributor *RWAYieldDistributorCallerSession) TotalYieldDistributed() (*big.Int, error) {
	return _RWAYieldDistributor.Contract.TotalYieldDistributed(&_RWAYieldDistributor.CallOpts)
}

// BatchClaimYield is a paid mutator transaction binding the contract method 0xd7c436d3.
//
// Solidity: function batchClaimYield(uint256[] distributionIds) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactor) BatchClaimYield(opts *bind.TransactOpts, distributionIds []*big.Int) (*types.Transaction, error) {
	return _RWAYieldDistributor.contract.Transact(opts, "batchClaimYield", distributionIds)
}

// BatchClaimYield is a paid mutator transaction binding the contract method 0xd7c436d3.
//
// Solidity: function batchClaimYield(uint256[] distributionIds) returns()
func (_RWAYieldDistributor *RWAYieldDistributorSession) BatchClaimYield(distributionIds []*big.Int) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.BatchClaimYield(&_RWAYieldDistributor.TransactOpts, distributionIds)
}

// BatchClaimYield is a paid mutator transaction binding the contract method 0xd7c436d3.
//
// Solidity: function batchClaimYield(uint256[] distributionIds) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactorSession) BatchClaimYield(distributionIds []*big.Int) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.BatchClaimYield(&_RWAYieldDistributor.TransactOpts, distributionIds)
}

// ClaimYield is a paid mutator transaction binding the contract method 0x40bd2e23.
//
// Solidity: function claimYield(uint256 distributionId) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactor) ClaimYield(opts *bind.TransactOpts, distributionId *big.Int) (*types.Transaction, error) {
	return _RWAYieldDistributor.contract.Transact(opts, "claimYield", distributionId)
}

// ClaimYield is a paid mutator transaction binding the contract method 0x40bd2e23.
//
// Solidity: function claimYield(uint256 distributionId) returns()
func (_RWAYieldDistributor *RWAYieldDistributorSession) ClaimYield(distributionId *big.Int) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.ClaimYield(&_RWAYieldDistributor.TransactOpts, distributionId)
}

// ClaimYield is a paid mutator transaction binding the contract method 0x40bd2e23.
//
// Solidity: function claimYield(uint256 distributionId) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactorSession) ClaimYield(distributionId *big.Int) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.ClaimYield(&_RWAYieldDistributor.TransactOpts, distributionId)
}

// DepositYield is a paid mutator transaction binding the contract method 0x5a10cac2.
//
// Solidity: function depositYield(uint256 assetId, address paymentToken, uint256 amount, uint256 claimPeriod) payable returns(uint256 distributionId)
func (_RWAYieldDistributor *RWAYieldDistributorTransactor) DepositYield(opts *bind.TransactOpts, assetId *big.Int, paymentToken common.Address, amount *big.Int, claimPeriod *big.Int) (*types.Transaction, error) {
	return _RWAYieldDistributor.contract.Transact(opts, "depositYield", assetId, paymentToken, amount, claimPeriod)
}

// DepositYield is a paid mutator transaction binding the contract method 0x5a10cac2.
//
// Solidity: function depositYield(uint256 assetId, address paymentToken, uint256 amount, uint256 claimPeriod) payable returns(uint256 distributionId)
func (_RWAYieldDistributor *RWAYieldDistributorSession) DepositYield(assetId *big.Int, paymentToken common.Address, amount *big.Int, claimPeriod *big.Int) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.DepositYield(&_RWAYieldDistributor.TransactOpts, assetId, paymentToken, amount, claimPeriod)
}

// DepositYield is a paid mutator transaction binding the contract method 0x5a10cac2.
//
// Solidity: function depositYield(uint256 assetId, address paymentToken, uint256 amount, uint256 claimPeriod) payable returns(uint256 distributionId)
func (_RWAYieldDistributor *RWAYieldDistributorTransactorSession) DepositYield(assetId *big.Int, paymentToken common.Address, amount *big.Int, claimPeriod *big.Int) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.DepositYield(&_RWAYieldDistributor.TransactOpts, assetId, paymentToken, amount, claimPeriod)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x551512de.
//
// Solidity: function emergencyWithdraw(address token, uint256 amount, address recipient) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactor) EmergencyWithdraw(opts *bind.TransactOpts, token common.Address, amount *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _RWAYieldDistributor.contract.Transact(opts, "emergencyWithdraw", token, amount, recipient)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x551512de.
//
// Solidity: function emergencyWithdraw(address token, uint256 amount, address recipient) returns()
func (_RWAYieldDistributor *RWAYieldDistributorSession) EmergencyWithdraw(token common.Address, amount *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.EmergencyWithdraw(&_RWAYieldDistributor.TransactOpts, token, amount, recipient)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x551512de.
//
// Solidity: function emergencyWithdraw(address token, uint256 amount, address recipient) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactorSession) EmergencyWithdraw(token common.Address, amount *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.EmergencyWithdraw(&_RWAYieldDistributor.TransactOpts, token, amount, recipient)
}

// FinalizeDistribution is a paid mutator transaction binding the contract method 0x95f52f4e.
//
// Solidity: function finalizeDistribution(uint256 distributionId) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactor) FinalizeDistribution(opts *bind.TransactOpts, distributionId *big.Int) (*types.Transaction, error) {
	return _RWAYieldDistributor.contract.Transact(opts, "finalizeDistribution", distributionId)
}

// FinalizeDistribution is a paid mutator transaction binding the contract method 0x95f52f4e.
//
// Solidity: function finalizeDistribution(uint256 distributionId) returns()
func (_RWAYieldDistributor *RWAYieldDistributorSession) FinalizeDistribution(distributionId *big.Int) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.FinalizeDistribution(&_RWAYieldDistributor.TransactOpts, distributionId)
}

// FinalizeDistribution is a paid mutator transaction binding the contract method 0x95f52f4e.
//
// Solidity: function finalizeDistribution(uint256 distributionId) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactorSession) FinalizeDistribution(distributionId *big.Int) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.FinalizeDistribution(&_RWAYieldDistributor.TransactOpts, distributionId)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAYieldDistributor.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAYieldDistributor *RWAYieldDistributorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.GrantRole(&_RWAYieldDistributor.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.GrantRole(&_RWAYieldDistributor.TransactOpts, role, account)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAYieldDistributor.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RWAYieldDistributor *RWAYieldDistributorSession) Pause() (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.Pause(&_RWAYieldDistributor.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactorSession) Pause() (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.Pause(&_RWAYieldDistributor.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAYieldDistributor.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAYieldDistributor *RWAYieldDistributorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.RenounceRole(&_RWAYieldDistributor.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.RenounceRole(&_RWAYieldDistributor.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAYieldDistributor.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAYieldDistributor *RWAYieldDistributorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.RevokeRole(&_RWAYieldDistributor.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.RevokeRole(&_RWAYieldDistributor.TransactOpts, role, account)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAYieldDistributor.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RWAYieldDistributor *RWAYieldDistributorSession) Unpause() (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.Unpause(&_RWAYieldDistributor.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactorSession) Unpause() (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.Unpause(&_RWAYieldDistributor.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAYieldDistributor.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_RWAYieldDistributor *RWAYieldDistributorSession) Receive() (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.Receive(&_RWAYieldDistributor.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_RWAYieldDistributor *RWAYieldDistributorTransactorSession) Receive() (*types.Transaction, error) {
	return _RWAYieldDistributor.Contract.Receive(&_RWAYieldDistributor.TransactOpts)
}

// RWAYieldDistributorDistributionFinalizedIterator is returned from FilterDistributionFinalized and is used to iterate over the raw logs and unpacked data for DistributionFinalized events raised by the RWAYieldDistributor contract.
type RWAYieldDistributorDistributionFinalizedIterator struct {
	Event *RWAYieldDistributorDistributionFinalized // Event containing the contract specifics and raw log

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
func (it *RWAYieldDistributorDistributionFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAYieldDistributorDistributionFinalized)
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
		it.Event = new(RWAYieldDistributorDistributionFinalized)
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
func (it *RWAYieldDistributorDistributionFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAYieldDistributorDistributionFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAYieldDistributorDistributionFinalized represents a DistributionFinalized event raised by the RWAYieldDistributor contract.
type RWAYieldDistributorDistributionFinalized struct {
	DistributionId *big.Int
	TotalClaimed   *big.Int
	Unclaimed      *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDistributionFinalized is a free log retrieval operation binding the contract event 0x89e1062634f0f216b5cdbe358b7a151c30e0f2d0179633e2dc56637c0942475f.
//
// Solidity: event DistributionFinalized(uint256 indexed distributionId, uint256 totalClaimed, uint256 unclaimed)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) FilterDistributionFinalized(opts *bind.FilterOpts, distributionId []*big.Int) (*RWAYieldDistributorDistributionFinalizedIterator, error) {

	var distributionIdRule []interface{}
	for _, distributionIdItem := range distributionId {
		distributionIdRule = append(distributionIdRule, distributionIdItem)
	}

	logs, sub, err := _RWAYieldDistributor.contract.FilterLogs(opts, "DistributionFinalized", distributionIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributorDistributionFinalizedIterator{contract: _RWAYieldDistributor.contract, event: "DistributionFinalized", logs: logs, sub: sub}, nil
}

// WatchDistributionFinalized is a free log subscription operation binding the contract event 0x89e1062634f0f216b5cdbe358b7a151c30e0f2d0179633e2dc56637c0942475f.
//
// Solidity: event DistributionFinalized(uint256 indexed distributionId, uint256 totalClaimed, uint256 unclaimed)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) WatchDistributionFinalized(opts *bind.WatchOpts, sink chan<- *RWAYieldDistributorDistributionFinalized, distributionId []*big.Int) (event.Subscription, error) {

	var distributionIdRule []interface{}
	for _, distributionIdItem := range distributionId {
		distributionIdRule = append(distributionIdRule, distributionIdItem)
	}

	logs, sub, err := _RWAYieldDistributor.contract.WatchLogs(opts, "DistributionFinalized", distributionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAYieldDistributorDistributionFinalized)
				if err := _RWAYieldDistributor.contract.UnpackLog(event, "DistributionFinalized", log); err != nil {
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

// ParseDistributionFinalized is a log parse operation binding the contract event 0x89e1062634f0f216b5cdbe358b7a151c30e0f2d0179633e2dc56637c0942475f.
//
// Solidity: event DistributionFinalized(uint256 indexed distributionId, uint256 totalClaimed, uint256 unclaimed)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) ParseDistributionFinalized(log types.Log) (*RWAYieldDistributorDistributionFinalized, error) {
	event := new(RWAYieldDistributorDistributionFinalized)
	if err := _RWAYieldDistributor.contract.UnpackLog(event, "DistributionFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAYieldDistributorPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the RWAYieldDistributor contract.
type RWAYieldDistributorPausedIterator struct {
	Event *RWAYieldDistributorPaused // Event containing the contract specifics and raw log

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
func (it *RWAYieldDistributorPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAYieldDistributorPaused)
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
		it.Event = new(RWAYieldDistributorPaused)
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
func (it *RWAYieldDistributorPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAYieldDistributorPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAYieldDistributorPaused represents a Paused event raised by the RWAYieldDistributor contract.
type RWAYieldDistributorPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) FilterPaused(opts *bind.FilterOpts) (*RWAYieldDistributorPausedIterator, error) {

	logs, sub, err := _RWAYieldDistributor.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributorPausedIterator{contract: _RWAYieldDistributor.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *RWAYieldDistributorPaused) (event.Subscription, error) {

	logs, sub, err := _RWAYieldDistributor.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAYieldDistributorPaused)
				if err := _RWAYieldDistributor.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) ParsePaused(log types.Log) (*RWAYieldDistributorPaused, error) {
	event := new(RWAYieldDistributorPaused)
	if err := _RWAYieldDistributor.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAYieldDistributorRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the RWAYieldDistributor contract.
type RWAYieldDistributorRoleAdminChangedIterator struct {
	Event *RWAYieldDistributorRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *RWAYieldDistributorRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAYieldDistributorRoleAdminChanged)
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
		it.Event = new(RWAYieldDistributorRoleAdminChanged)
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
func (it *RWAYieldDistributorRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAYieldDistributorRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAYieldDistributorRoleAdminChanged represents a RoleAdminChanged event raised by the RWAYieldDistributor contract.
type RWAYieldDistributorRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*RWAYieldDistributorRoleAdminChangedIterator, error) {

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

	logs, sub, err := _RWAYieldDistributor.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributorRoleAdminChangedIterator{contract: _RWAYieldDistributor.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *RWAYieldDistributorRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _RWAYieldDistributor.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAYieldDistributorRoleAdminChanged)
				if err := _RWAYieldDistributor.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) ParseRoleAdminChanged(log types.Log) (*RWAYieldDistributorRoleAdminChanged, error) {
	event := new(RWAYieldDistributorRoleAdminChanged)
	if err := _RWAYieldDistributor.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAYieldDistributorRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the RWAYieldDistributor contract.
type RWAYieldDistributorRoleGrantedIterator struct {
	Event *RWAYieldDistributorRoleGranted // Event containing the contract specifics and raw log

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
func (it *RWAYieldDistributorRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAYieldDistributorRoleGranted)
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
		it.Event = new(RWAYieldDistributorRoleGranted)
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
func (it *RWAYieldDistributorRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAYieldDistributorRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAYieldDistributorRoleGranted represents a RoleGranted event raised by the RWAYieldDistributor contract.
type RWAYieldDistributorRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RWAYieldDistributorRoleGrantedIterator, error) {

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

	logs, sub, err := _RWAYieldDistributor.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributorRoleGrantedIterator{contract: _RWAYieldDistributor.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *RWAYieldDistributorRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RWAYieldDistributor.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAYieldDistributorRoleGranted)
				if err := _RWAYieldDistributor.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) ParseRoleGranted(log types.Log) (*RWAYieldDistributorRoleGranted, error) {
	event := new(RWAYieldDistributorRoleGranted)
	if err := _RWAYieldDistributor.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAYieldDistributorRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the RWAYieldDistributor contract.
type RWAYieldDistributorRoleRevokedIterator struct {
	Event *RWAYieldDistributorRoleRevoked // Event containing the contract specifics and raw log

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
func (it *RWAYieldDistributorRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAYieldDistributorRoleRevoked)
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
		it.Event = new(RWAYieldDistributorRoleRevoked)
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
func (it *RWAYieldDistributorRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAYieldDistributorRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAYieldDistributorRoleRevoked represents a RoleRevoked event raised by the RWAYieldDistributor contract.
type RWAYieldDistributorRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RWAYieldDistributorRoleRevokedIterator, error) {

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

	logs, sub, err := _RWAYieldDistributor.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributorRoleRevokedIterator{contract: _RWAYieldDistributor.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *RWAYieldDistributorRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RWAYieldDistributor.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAYieldDistributorRoleRevoked)
				if err := _RWAYieldDistributor.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) ParseRoleRevoked(log types.Log) (*RWAYieldDistributorRoleRevoked, error) {
	event := new(RWAYieldDistributorRoleRevoked)
	if err := _RWAYieldDistributor.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAYieldDistributorUnclaimedYieldReclaimedIterator is returned from FilterUnclaimedYieldReclaimed and is used to iterate over the raw logs and unpacked data for UnclaimedYieldReclaimed events raised by the RWAYieldDistributor contract.
type RWAYieldDistributorUnclaimedYieldReclaimedIterator struct {
	Event *RWAYieldDistributorUnclaimedYieldReclaimed // Event containing the contract specifics and raw log

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
func (it *RWAYieldDistributorUnclaimedYieldReclaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAYieldDistributorUnclaimedYieldReclaimed)
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
		it.Event = new(RWAYieldDistributorUnclaimedYieldReclaimed)
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
func (it *RWAYieldDistributorUnclaimedYieldReclaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAYieldDistributorUnclaimedYieldReclaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAYieldDistributorUnclaimedYieldReclaimed represents a UnclaimedYieldReclaimed event raised by the RWAYieldDistributor contract.
type RWAYieldDistributorUnclaimedYieldReclaimed struct {
	DistributionId *big.Int
	Amount         *big.Int
	Recipient      common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUnclaimedYieldReclaimed is a free log retrieval operation binding the contract event 0xf744ff5ba0db05729f6d246a15dfcd16be09e279a0d035c285d08b0af0d8a2b7.
//
// Solidity: event UnclaimedYieldReclaimed(uint256 indexed distributionId, uint256 amount, address recipient)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) FilterUnclaimedYieldReclaimed(opts *bind.FilterOpts, distributionId []*big.Int) (*RWAYieldDistributorUnclaimedYieldReclaimedIterator, error) {

	var distributionIdRule []interface{}
	for _, distributionIdItem := range distributionId {
		distributionIdRule = append(distributionIdRule, distributionIdItem)
	}

	logs, sub, err := _RWAYieldDistributor.contract.FilterLogs(opts, "UnclaimedYieldReclaimed", distributionIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributorUnclaimedYieldReclaimedIterator{contract: _RWAYieldDistributor.contract, event: "UnclaimedYieldReclaimed", logs: logs, sub: sub}, nil
}

// WatchUnclaimedYieldReclaimed is a free log subscription operation binding the contract event 0xf744ff5ba0db05729f6d246a15dfcd16be09e279a0d035c285d08b0af0d8a2b7.
//
// Solidity: event UnclaimedYieldReclaimed(uint256 indexed distributionId, uint256 amount, address recipient)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) WatchUnclaimedYieldReclaimed(opts *bind.WatchOpts, sink chan<- *RWAYieldDistributorUnclaimedYieldReclaimed, distributionId []*big.Int) (event.Subscription, error) {

	var distributionIdRule []interface{}
	for _, distributionIdItem := range distributionId {
		distributionIdRule = append(distributionIdRule, distributionIdItem)
	}

	logs, sub, err := _RWAYieldDistributor.contract.WatchLogs(opts, "UnclaimedYieldReclaimed", distributionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAYieldDistributorUnclaimedYieldReclaimed)
				if err := _RWAYieldDistributor.contract.UnpackLog(event, "UnclaimedYieldReclaimed", log); err != nil {
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

// ParseUnclaimedYieldReclaimed is a log parse operation binding the contract event 0xf744ff5ba0db05729f6d246a15dfcd16be09e279a0d035c285d08b0af0d8a2b7.
//
// Solidity: event UnclaimedYieldReclaimed(uint256 indexed distributionId, uint256 amount, address recipient)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) ParseUnclaimedYieldReclaimed(log types.Log) (*RWAYieldDistributorUnclaimedYieldReclaimed, error) {
	event := new(RWAYieldDistributorUnclaimedYieldReclaimed)
	if err := _RWAYieldDistributor.contract.UnpackLog(event, "UnclaimedYieldReclaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAYieldDistributorUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the RWAYieldDistributor contract.
type RWAYieldDistributorUnpausedIterator struct {
	Event *RWAYieldDistributorUnpaused // Event containing the contract specifics and raw log

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
func (it *RWAYieldDistributorUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAYieldDistributorUnpaused)
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
		it.Event = new(RWAYieldDistributorUnpaused)
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
func (it *RWAYieldDistributorUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAYieldDistributorUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAYieldDistributorUnpaused represents a Unpaused event raised by the RWAYieldDistributor contract.
type RWAYieldDistributorUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) FilterUnpaused(opts *bind.FilterOpts) (*RWAYieldDistributorUnpausedIterator, error) {

	logs, sub, err := _RWAYieldDistributor.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributorUnpausedIterator{contract: _RWAYieldDistributor.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *RWAYieldDistributorUnpaused) (event.Subscription, error) {

	logs, sub, err := _RWAYieldDistributor.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAYieldDistributorUnpaused)
				if err := _RWAYieldDistributor.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) ParseUnpaused(log types.Log) (*RWAYieldDistributorUnpaused, error) {
	event := new(RWAYieldDistributorUnpaused)
	if err := _RWAYieldDistributor.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAYieldDistributorYieldClaimedIterator is returned from FilterYieldClaimed and is used to iterate over the raw logs and unpacked data for YieldClaimed events raised by the RWAYieldDistributor contract.
type RWAYieldDistributorYieldClaimedIterator struct {
	Event *RWAYieldDistributorYieldClaimed // Event containing the contract specifics and raw log

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
func (it *RWAYieldDistributorYieldClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAYieldDistributorYieldClaimed)
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
		it.Event = new(RWAYieldDistributorYieldClaimed)
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
func (it *RWAYieldDistributorYieldClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAYieldDistributorYieldClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAYieldDistributorYieldClaimed represents a YieldClaimed event raised by the RWAYieldDistributor contract.
type RWAYieldDistributorYieldClaimed struct {
	DistributionId *big.Int
	User           common.Address
	Amount         *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterYieldClaimed is a free log retrieval operation binding the contract event 0x1448cb68839238fd008abe5bad73d996f9849e020313466fa9f3476d98006df2.
//
// Solidity: event YieldClaimed(uint256 indexed distributionId, address indexed user, uint256 amount)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) FilterYieldClaimed(opts *bind.FilterOpts, distributionId []*big.Int, user []common.Address) (*RWAYieldDistributorYieldClaimedIterator, error) {

	var distributionIdRule []interface{}
	for _, distributionIdItem := range distributionId {
		distributionIdRule = append(distributionIdRule, distributionIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _RWAYieldDistributor.contract.FilterLogs(opts, "YieldClaimed", distributionIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributorYieldClaimedIterator{contract: _RWAYieldDistributor.contract, event: "YieldClaimed", logs: logs, sub: sub}, nil
}

// WatchYieldClaimed is a free log subscription operation binding the contract event 0x1448cb68839238fd008abe5bad73d996f9849e020313466fa9f3476d98006df2.
//
// Solidity: event YieldClaimed(uint256 indexed distributionId, address indexed user, uint256 amount)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) WatchYieldClaimed(opts *bind.WatchOpts, sink chan<- *RWAYieldDistributorYieldClaimed, distributionId []*big.Int, user []common.Address) (event.Subscription, error) {

	var distributionIdRule []interface{}
	for _, distributionIdItem := range distributionId {
		distributionIdRule = append(distributionIdRule, distributionIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _RWAYieldDistributor.contract.WatchLogs(opts, "YieldClaimed", distributionIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAYieldDistributorYieldClaimed)
				if err := _RWAYieldDistributor.contract.UnpackLog(event, "YieldClaimed", log); err != nil {
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

// ParseYieldClaimed is a log parse operation binding the contract event 0x1448cb68839238fd008abe5bad73d996f9849e020313466fa9f3476d98006df2.
//
// Solidity: event YieldClaimed(uint256 indexed distributionId, address indexed user, uint256 amount)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) ParseYieldClaimed(log types.Log) (*RWAYieldDistributorYieldClaimed, error) {
	event := new(RWAYieldDistributorYieldClaimed)
	if err := _RWAYieldDistributor.contract.UnpackLog(event, "YieldClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAYieldDistributorYieldDepositedIterator is returned from FilterYieldDeposited and is used to iterate over the raw logs and unpacked data for YieldDeposited events raised by the RWAYieldDistributor contract.
type RWAYieldDistributorYieldDepositedIterator struct {
	Event *RWAYieldDistributorYieldDeposited // Event containing the contract specifics and raw log

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
func (it *RWAYieldDistributorYieldDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAYieldDistributorYieldDeposited)
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
		it.Event = new(RWAYieldDistributorYieldDeposited)
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
func (it *RWAYieldDistributorYieldDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAYieldDistributorYieldDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAYieldDistributorYieldDeposited represents a YieldDeposited event raised by the RWAYieldDistributor contract.
type RWAYieldDistributorYieldDeposited struct {
	DistributionId *big.Int
	AssetId        *big.Int
	PaymentToken   common.Address
	Amount         *big.Int
	ClaimDeadline  *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterYieldDeposited is a free log retrieval operation binding the contract event 0xccd0710800fceb311886b4326a31fc8191a754a02d6621b70133234cb3713228.
//
// Solidity: event YieldDeposited(uint256 indexed distributionId, uint256 indexed assetId, address indexed paymentToken, uint256 amount, uint256 claimDeadline)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) FilterYieldDeposited(opts *bind.FilterOpts, distributionId []*big.Int, assetId []*big.Int, paymentToken []common.Address) (*RWAYieldDistributorYieldDepositedIterator, error) {

	var distributionIdRule []interface{}
	for _, distributionIdItem := range distributionId {
		distributionIdRule = append(distributionIdRule, distributionIdItem)
	}
	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var paymentTokenRule []interface{}
	for _, paymentTokenItem := range paymentToken {
		paymentTokenRule = append(paymentTokenRule, paymentTokenItem)
	}

	logs, sub, err := _RWAYieldDistributor.contract.FilterLogs(opts, "YieldDeposited", distributionIdRule, assetIdRule, paymentTokenRule)
	if err != nil {
		return nil, err
	}
	return &RWAYieldDistributorYieldDepositedIterator{contract: _RWAYieldDistributor.contract, event: "YieldDeposited", logs: logs, sub: sub}, nil
}

// WatchYieldDeposited is a free log subscription operation binding the contract event 0xccd0710800fceb311886b4326a31fc8191a754a02d6621b70133234cb3713228.
//
// Solidity: event YieldDeposited(uint256 indexed distributionId, uint256 indexed assetId, address indexed paymentToken, uint256 amount, uint256 claimDeadline)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) WatchYieldDeposited(opts *bind.WatchOpts, sink chan<- *RWAYieldDistributorYieldDeposited, distributionId []*big.Int, assetId []*big.Int, paymentToken []common.Address) (event.Subscription, error) {

	var distributionIdRule []interface{}
	for _, distributionIdItem := range distributionId {
		distributionIdRule = append(distributionIdRule, distributionIdItem)
	}
	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var paymentTokenRule []interface{}
	for _, paymentTokenItem := range paymentToken {
		paymentTokenRule = append(paymentTokenRule, paymentTokenItem)
	}

	logs, sub, err := _RWAYieldDistributor.contract.WatchLogs(opts, "YieldDeposited", distributionIdRule, assetIdRule, paymentTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAYieldDistributorYieldDeposited)
				if err := _RWAYieldDistributor.contract.UnpackLog(event, "YieldDeposited", log); err != nil {
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

// ParseYieldDeposited is a log parse operation binding the contract event 0xccd0710800fceb311886b4326a31fc8191a754a02d6621b70133234cb3713228.
//
// Solidity: event YieldDeposited(uint256 indexed distributionId, uint256 indexed assetId, address indexed paymentToken, uint256 amount, uint256 claimDeadline)
func (_RWAYieldDistributor *RWAYieldDistributorFilterer) ParseYieldDeposited(log types.Log) (*RWAYieldDistributorYieldDeposited, error) {
	event := new(RWAYieldDistributorYieldDeposited)
	if err := _RWAYieldDistributor.contract.UnpackLog(event, "YieldDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
