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

// L1GatewayMetaData contains all meta data concerning the L1Gateway contract.
var L1GatewayMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_collateralVault\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_loyaltyUSD\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_collateralToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_inbox\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_outbox\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"depositId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DepositFinalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"depositId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"ticketId\",\"type\":\"uint256\"}],\"name\":\"DepositInitiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxSubmissionCost\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxGas\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPriceBid\",\"type\":\"uint256\"}],\"name\":\"GasParametersUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldGateway\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newGateway\",\"type\":\"address\"}],\"name\":\"L2GatewayUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"withdrawalId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawalExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"withdrawalId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawalInitiated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"l2TxHash\",\"type\":\"bytes32\"}],\"name\":\"bridgeBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"l2TxHash\",\"type\":\"bytes32\"}],\"name\":\"bridgeMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateRequiredEth\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collateralToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collateralVault\",\"outputs\":[{\"internalType\":\"contractCollateralVaultL1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositToL2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"depositId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ticketId\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"deposits\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"l2TxHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"processed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"withdrawalId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"finalizeWithdrawal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasPriceBid\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"depositId\",\"type\":\"uint256\"}],\"name\":\"getDeposit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"l2TxHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"processed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"withdrawalId\",\"type\":\"bytes32\"}],\"name\":\"getWithdrawal\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"executed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"inbox\",\"outputs\":[{\"internalType\":\"contractIInbox\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2Gateway\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"loyaltyUSD\",\"outputs\":[{\"internalType\":\"contractLoyaltyUSDL1\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxGas\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxSubmissionCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"outbox\",\"outputs\":[{\"internalType\":\"contractIOutbox\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_maxSubmissionCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasPriceBid\",\"type\":\"uint256\"}],\"name\":\"setGasParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_inbox\",\"type\":\"address\"}],\"name\":\"setInbox\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l2Gateway\",\"type\":\"address\"}],\"name\":\"setL2Gateway\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_outbox\",\"type\":\"address\"}],\"name\":\"setOutbox\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"withdrawals\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"executed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// L1GatewayABI is the input ABI used to generate the binding from.
// Deprecated: Use L1GatewayMetaData.ABI instead.
var L1GatewayABI = L1GatewayMetaData.ABI

// L1Gateway is an auto generated Go binding around an Ethereum contract.
type L1Gateway struct {
	L1GatewayCaller     // Read-only binding to the contract
	L1GatewayTransactor // Write-only binding to the contract
	L1GatewayFilterer   // Log filterer for contract events
}

// L1GatewayCaller is an auto generated read-only Go binding around an Ethereum contract.
type L1GatewayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1GatewayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L1GatewayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1GatewayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L1GatewayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1GatewaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L1GatewaySession struct {
	Contract     *L1Gateway        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L1GatewayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L1GatewayCallerSession struct {
	Contract *L1GatewayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// L1GatewayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L1GatewayTransactorSession struct {
	Contract     *L1GatewayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// L1GatewayRaw is an auto generated low-level Go binding around an Ethereum contract.
type L1GatewayRaw struct {
	Contract *L1Gateway // Generic contract binding to access the raw methods on
}

// L1GatewayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L1GatewayCallerRaw struct {
	Contract *L1GatewayCaller // Generic read-only contract binding to access the raw methods on
}

// L1GatewayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L1GatewayTransactorRaw struct {
	Contract *L1GatewayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL1Gateway creates a new instance of L1Gateway, bound to a specific deployed contract.
func NewL1Gateway(address common.Address, backend bind.ContractBackend) (*L1Gateway, error) {
	contract, err := bindL1Gateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L1Gateway{L1GatewayCaller: L1GatewayCaller{contract: contract}, L1GatewayTransactor: L1GatewayTransactor{contract: contract}, L1GatewayFilterer: L1GatewayFilterer{contract: contract}}, nil
}

// NewL1GatewayCaller creates a new read-only instance of L1Gateway, bound to a specific deployed contract.
func NewL1GatewayCaller(address common.Address, caller bind.ContractCaller) (*L1GatewayCaller, error) {
	contract, err := bindL1Gateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L1GatewayCaller{contract: contract}, nil
}

// NewL1GatewayTransactor creates a new write-only instance of L1Gateway, bound to a specific deployed contract.
func NewL1GatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*L1GatewayTransactor, error) {
	contract, err := bindL1Gateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L1GatewayTransactor{contract: contract}, nil
}

// NewL1GatewayFilterer creates a new log filterer instance of L1Gateway, bound to a specific deployed contract.
func NewL1GatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*L1GatewayFilterer, error) {
	contract, err := bindL1Gateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L1GatewayFilterer{contract: contract}, nil
}

// bindL1Gateway binds a generic wrapper to an already deployed contract.
func bindL1Gateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L1GatewayMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1Gateway *L1GatewayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1Gateway.Contract.L1GatewayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1Gateway *L1GatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Gateway.Contract.L1GatewayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1Gateway *L1GatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1Gateway.Contract.L1GatewayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1Gateway *L1GatewayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1Gateway.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1Gateway *L1GatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Gateway.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1Gateway *L1GatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1Gateway.Contract.contract.Transact(opts, method, params...)
}

// CalculateRequiredEth is a free data retrieval call binding the contract method 0x3928e5a9.
//
// Solidity: function calculateRequiredEth() view returns(uint256)
func (_L1Gateway *L1GatewayCaller) CalculateRequiredEth(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "calculateRequiredEth")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateRequiredEth is a free data retrieval call binding the contract method 0x3928e5a9.
//
// Solidity: function calculateRequiredEth() view returns(uint256)
func (_L1Gateway *L1GatewaySession) CalculateRequiredEth() (*big.Int, error) {
	return _L1Gateway.Contract.CalculateRequiredEth(&_L1Gateway.CallOpts)
}

// CalculateRequiredEth is a free data retrieval call binding the contract method 0x3928e5a9.
//
// Solidity: function calculateRequiredEth() view returns(uint256)
func (_L1Gateway *L1GatewayCallerSession) CalculateRequiredEth() (*big.Int, error) {
	return _L1Gateway.Contract.CalculateRequiredEth(&_L1Gateway.CallOpts)
}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_L1Gateway *L1GatewayCaller) CollateralToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "collateralToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_L1Gateway *L1GatewaySession) CollateralToken() (common.Address, error) {
	return _L1Gateway.Contract.CollateralToken(&_L1Gateway.CallOpts)
}

// CollateralToken is a free data retrieval call binding the contract method 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (_L1Gateway *L1GatewayCallerSession) CollateralToken() (common.Address, error) {
	return _L1Gateway.Contract.CollateralToken(&_L1Gateway.CallOpts)
}

// CollateralVault is a free data retrieval call binding the contract method 0x0bece79c.
//
// Solidity: function collateralVault() view returns(address)
func (_L1Gateway *L1GatewayCaller) CollateralVault(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "collateralVault")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CollateralVault is a free data retrieval call binding the contract method 0x0bece79c.
//
// Solidity: function collateralVault() view returns(address)
func (_L1Gateway *L1GatewaySession) CollateralVault() (common.Address, error) {
	return _L1Gateway.Contract.CollateralVault(&_L1Gateway.CallOpts)
}

// CollateralVault is a free data retrieval call binding the contract method 0x0bece79c.
//
// Solidity: function collateralVault() view returns(address)
func (_L1Gateway *L1GatewayCallerSession) CollateralVault() (common.Address, error) {
	return _L1Gateway.Contract.CollateralVault(&_L1Gateway.CallOpts)
}

// DepositNonce is a free data retrieval call binding the contract method 0xde35f5cb.
//
// Solidity: function depositNonce() view returns(uint256)
func (_L1Gateway *L1GatewayCaller) DepositNonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "depositNonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositNonce is a free data retrieval call binding the contract method 0xde35f5cb.
//
// Solidity: function depositNonce() view returns(uint256)
func (_L1Gateway *L1GatewaySession) DepositNonce() (*big.Int, error) {
	return _L1Gateway.Contract.DepositNonce(&_L1Gateway.CallOpts)
}

// DepositNonce is a free data retrieval call binding the contract method 0xde35f5cb.
//
// Solidity: function depositNonce() view returns(uint256)
func (_L1Gateway *L1GatewayCallerSession) DepositNonce() (*big.Int, error) {
	return _L1Gateway.Contract.DepositNonce(&_L1Gateway.CallOpts)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) view returns(address user, uint256 amount, uint256 timestamp, bytes32 l2TxHash, bool processed)
func (_L1Gateway *L1GatewayCaller) Deposits(opts *bind.CallOpts, arg0 *big.Int) (struct {
	User      common.Address
	Amount    *big.Int
	Timestamp *big.Int
	L2TxHash  [32]byte
	Processed bool
}, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "deposits", arg0)

	outstruct := new(struct {
		User      common.Address
		Amount    *big.Int
		Timestamp *big.Int
		L2TxHash  [32]byte
		Processed bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.User = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.L2TxHash = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.Processed = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) view returns(address user, uint256 amount, uint256 timestamp, bytes32 l2TxHash, bool processed)
func (_L1Gateway *L1GatewaySession) Deposits(arg0 *big.Int) (struct {
	User      common.Address
	Amount    *big.Int
	Timestamp *big.Int
	L2TxHash  [32]byte
	Processed bool
}, error) {
	return _L1Gateway.Contract.Deposits(&_L1Gateway.CallOpts, arg0)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) view returns(address user, uint256 amount, uint256 timestamp, bytes32 l2TxHash, bool processed)
func (_L1Gateway *L1GatewayCallerSession) Deposits(arg0 *big.Int) (struct {
	User      common.Address
	Amount    *big.Int
	Timestamp *big.Int
	L2TxHash  [32]byte
	Processed bool
}, error) {
	return _L1Gateway.Contract.Deposits(&_L1Gateway.CallOpts, arg0)
}

// GasPriceBid is a free data retrieval call binding the contract method 0x5d942ac1.
//
// Solidity: function gasPriceBid() view returns(uint256)
func (_L1Gateway *L1GatewayCaller) GasPriceBid(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "gasPriceBid")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GasPriceBid is a free data retrieval call binding the contract method 0x5d942ac1.
//
// Solidity: function gasPriceBid() view returns(uint256)
func (_L1Gateway *L1GatewaySession) GasPriceBid() (*big.Int, error) {
	return _L1Gateway.Contract.GasPriceBid(&_L1Gateway.CallOpts)
}

// GasPriceBid is a free data retrieval call binding the contract method 0x5d942ac1.
//
// Solidity: function gasPriceBid() view returns(uint256)
func (_L1Gateway *L1GatewayCallerSession) GasPriceBid() (*big.Int, error) {
	return _L1Gateway.Contract.GasPriceBid(&_L1Gateway.CallOpts)
}

// GetDeposit is a free data retrieval call binding the contract method 0x9f9fb968.
//
// Solidity: function getDeposit(uint256 depositId) view returns(address user, uint256 amount, uint256 timestamp, bytes32 l2TxHash, bool processed)
func (_L1Gateway *L1GatewayCaller) GetDeposit(opts *bind.CallOpts, depositId *big.Int) (struct {
	User      common.Address
	Amount    *big.Int
	Timestamp *big.Int
	L2TxHash  [32]byte
	Processed bool
}, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "getDeposit", depositId)

	outstruct := new(struct {
		User      common.Address
		Amount    *big.Int
		Timestamp *big.Int
		L2TxHash  [32]byte
		Processed bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.User = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.L2TxHash = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.Processed = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// GetDeposit is a free data retrieval call binding the contract method 0x9f9fb968.
//
// Solidity: function getDeposit(uint256 depositId) view returns(address user, uint256 amount, uint256 timestamp, bytes32 l2TxHash, bool processed)
func (_L1Gateway *L1GatewaySession) GetDeposit(depositId *big.Int) (struct {
	User      common.Address
	Amount    *big.Int
	Timestamp *big.Int
	L2TxHash  [32]byte
	Processed bool
}, error) {
	return _L1Gateway.Contract.GetDeposit(&_L1Gateway.CallOpts, depositId)
}

// GetDeposit is a free data retrieval call binding the contract method 0x9f9fb968.
//
// Solidity: function getDeposit(uint256 depositId) view returns(address user, uint256 amount, uint256 timestamp, bytes32 l2TxHash, bool processed)
func (_L1Gateway *L1GatewayCallerSession) GetDeposit(depositId *big.Int) (struct {
	User      common.Address
	Amount    *big.Int
	Timestamp *big.Int
	L2TxHash  [32]byte
	Processed bool
}, error) {
	return _L1Gateway.Contract.GetDeposit(&_L1Gateway.CallOpts, depositId)
}

// GetWithdrawal is a free data retrieval call binding the contract method 0x98efa653.
//
// Solidity: function getWithdrawal(bytes32 withdrawalId) view returns(address user, uint256 amount, uint256 timestamp, bool executed)
func (_L1Gateway *L1GatewayCaller) GetWithdrawal(opts *bind.CallOpts, withdrawalId [32]byte) (struct {
	User      common.Address
	Amount    *big.Int
	Timestamp *big.Int
	Executed  bool
}, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "getWithdrawal", withdrawalId)

	outstruct := new(struct {
		User      common.Address
		Amount    *big.Int
		Timestamp *big.Int
		Executed  bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.User = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Executed = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// GetWithdrawal is a free data retrieval call binding the contract method 0x98efa653.
//
// Solidity: function getWithdrawal(bytes32 withdrawalId) view returns(address user, uint256 amount, uint256 timestamp, bool executed)
func (_L1Gateway *L1GatewaySession) GetWithdrawal(withdrawalId [32]byte) (struct {
	User      common.Address
	Amount    *big.Int
	Timestamp *big.Int
	Executed  bool
}, error) {
	return _L1Gateway.Contract.GetWithdrawal(&_L1Gateway.CallOpts, withdrawalId)
}

// GetWithdrawal is a free data retrieval call binding the contract method 0x98efa653.
//
// Solidity: function getWithdrawal(bytes32 withdrawalId) view returns(address user, uint256 amount, uint256 timestamp, bool executed)
func (_L1Gateway *L1GatewayCallerSession) GetWithdrawal(withdrawalId [32]byte) (struct {
	User      common.Address
	Amount    *big.Int
	Timestamp *big.Int
	Executed  bool
}, error) {
	return _L1Gateway.Contract.GetWithdrawal(&_L1Gateway.CallOpts, withdrawalId)
}

// Inbox is a free data retrieval call binding the contract method 0xfb0e722b.
//
// Solidity: function inbox() view returns(address)
func (_L1Gateway *L1GatewayCaller) Inbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "inbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Inbox is a free data retrieval call binding the contract method 0xfb0e722b.
//
// Solidity: function inbox() view returns(address)
func (_L1Gateway *L1GatewaySession) Inbox() (common.Address, error) {
	return _L1Gateway.Contract.Inbox(&_L1Gateway.CallOpts)
}

// Inbox is a free data retrieval call binding the contract method 0xfb0e722b.
//
// Solidity: function inbox() view returns(address)
func (_L1Gateway *L1GatewayCallerSession) Inbox() (common.Address, error) {
	return _L1Gateway.Contract.Inbox(&_L1Gateway.CallOpts)
}

// L2Gateway is a free data retrieval call binding the contract method 0x8fa74a0e.
//
// Solidity: function l2Gateway() view returns(address)
func (_L1Gateway *L1GatewayCaller) L2Gateway(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "l2Gateway")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2Gateway is a free data retrieval call binding the contract method 0x8fa74a0e.
//
// Solidity: function l2Gateway() view returns(address)
func (_L1Gateway *L1GatewaySession) L2Gateway() (common.Address, error) {
	return _L1Gateway.Contract.L2Gateway(&_L1Gateway.CallOpts)
}

// L2Gateway is a free data retrieval call binding the contract method 0x8fa74a0e.
//
// Solidity: function l2Gateway() view returns(address)
func (_L1Gateway *L1GatewayCallerSession) L2Gateway() (common.Address, error) {
	return _L1Gateway.Contract.L2Gateway(&_L1Gateway.CallOpts)
}

// LoyaltyUSD is a free data retrieval call binding the contract method 0xa2e950cc.
//
// Solidity: function loyaltyUSD() view returns(address)
func (_L1Gateway *L1GatewayCaller) LoyaltyUSD(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "loyaltyUSD")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LoyaltyUSD is a free data retrieval call binding the contract method 0xa2e950cc.
//
// Solidity: function loyaltyUSD() view returns(address)
func (_L1Gateway *L1GatewaySession) LoyaltyUSD() (common.Address, error) {
	return _L1Gateway.Contract.LoyaltyUSD(&_L1Gateway.CallOpts)
}

// LoyaltyUSD is a free data retrieval call binding the contract method 0xa2e950cc.
//
// Solidity: function loyaltyUSD() view returns(address)
func (_L1Gateway *L1GatewayCallerSession) LoyaltyUSD() (common.Address, error) {
	return _L1Gateway.Contract.LoyaltyUSD(&_L1Gateway.CallOpts)
}

// MaxGas is a free data retrieval call binding the contract method 0x501d815c.
//
// Solidity: function maxGas() view returns(uint256)
func (_L1Gateway *L1GatewayCaller) MaxGas(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "maxGas")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxGas is a free data retrieval call binding the contract method 0x501d815c.
//
// Solidity: function maxGas() view returns(uint256)
func (_L1Gateway *L1GatewaySession) MaxGas() (*big.Int, error) {
	return _L1Gateway.Contract.MaxGas(&_L1Gateway.CallOpts)
}

// MaxGas is a free data retrieval call binding the contract method 0x501d815c.
//
// Solidity: function maxGas() view returns(uint256)
func (_L1Gateway *L1GatewayCallerSession) MaxGas() (*big.Int, error) {
	return _L1Gateway.Contract.MaxGas(&_L1Gateway.CallOpts)
}

// MaxSubmissionCost is a free data retrieval call binding the contract method 0x70123fee.
//
// Solidity: function maxSubmissionCost() view returns(uint256)
func (_L1Gateway *L1GatewayCaller) MaxSubmissionCost(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "maxSubmissionCost")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxSubmissionCost is a free data retrieval call binding the contract method 0x70123fee.
//
// Solidity: function maxSubmissionCost() view returns(uint256)
func (_L1Gateway *L1GatewaySession) MaxSubmissionCost() (*big.Int, error) {
	return _L1Gateway.Contract.MaxSubmissionCost(&_L1Gateway.CallOpts)
}

// MaxSubmissionCost is a free data retrieval call binding the contract method 0x70123fee.
//
// Solidity: function maxSubmissionCost() view returns(uint256)
func (_L1Gateway *L1GatewayCallerSession) MaxSubmissionCost() (*big.Int, error) {
	return _L1Gateway.Contract.MaxSubmissionCost(&_L1Gateway.CallOpts)
}

// Outbox is a free data retrieval call binding the contract method 0xce11e6ab.
//
// Solidity: function outbox() view returns(address)
func (_L1Gateway *L1GatewayCaller) Outbox(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "outbox")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Outbox is a free data retrieval call binding the contract method 0xce11e6ab.
//
// Solidity: function outbox() view returns(address)
func (_L1Gateway *L1GatewaySession) Outbox() (common.Address, error) {
	return _L1Gateway.Contract.Outbox(&_L1Gateway.CallOpts)
}

// Outbox is a free data retrieval call binding the contract method 0xce11e6ab.
//
// Solidity: function outbox() view returns(address)
func (_L1Gateway *L1GatewayCallerSession) Outbox() (common.Address, error) {
	return _L1Gateway.Contract.Outbox(&_L1Gateway.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1Gateway *L1GatewayCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1Gateway *L1GatewaySession) Owner() (common.Address, error) {
	return _L1Gateway.Contract.Owner(&_L1Gateway.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1Gateway *L1GatewayCallerSession) Owner() (common.Address, error) {
	return _L1Gateway.Contract.Owner(&_L1Gateway.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_L1Gateway *L1GatewayCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_L1Gateway *L1GatewaySession) Paused() (bool, error) {
	return _L1Gateway.Contract.Paused(&_L1Gateway.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_L1Gateway *L1GatewayCallerSession) Paused() (bool, error) {
	return _L1Gateway.Contract.Paused(&_L1Gateway.CallOpts)
}

// Withdrawals is a free data retrieval call binding the contract method 0xefbf64a7.
//
// Solidity: function withdrawals(bytes32 ) view returns(address user, uint256 amount, uint256 timestamp, bool executed)
func (_L1Gateway *L1GatewayCaller) Withdrawals(opts *bind.CallOpts, arg0 [32]byte) (struct {
	User      common.Address
	Amount    *big.Int
	Timestamp *big.Int
	Executed  bool
}, error) {
	var out []interface{}
	err := _L1Gateway.contract.Call(opts, &out, "withdrawals", arg0)

	outstruct := new(struct {
		User      common.Address
		Amount    *big.Int
		Timestamp *big.Int
		Executed  bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.User = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Executed = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// Withdrawals is a free data retrieval call binding the contract method 0xefbf64a7.
//
// Solidity: function withdrawals(bytes32 ) view returns(address user, uint256 amount, uint256 timestamp, bool executed)
func (_L1Gateway *L1GatewaySession) Withdrawals(arg0 [32]byte) (struct {
	User      common.Address
	Amount    *big.Int
	Timestamp *big.Int
	Executed  bool
}, error) {
	return _L1Gateway.Contract.Withdrawals(&_L1Gateway.CallOpts, arg0)
}

// Withdrawals is a free data retrieval call binding the contract method 0xefbf64a7.
//
// Solidity: function withdrawals(bytes32 ) view returns(address user, uint256 amount, uint256 timestamp, bool executed)
func (_L1Gateway *L1GatewayCallerSession) Withdrawals(arg0 [32]byte) (struct {
	User      common.Address
	Amount    *big.Int
	Timestamp *big.Int
	Executed  bool
}, error) {
	return _L1Gateway.Contract.Withdrawals(&_L1Gateway.CallOpts, arg0)
}

// BridgeBurn is a paid mutator transaction binding the contract method 0x87e04042.
//
// Solidity: function bridgeBurn(address user, uint256 amount, bytes32 l2TxHash) returns()
func (_L1Gateway *L1GatewayTransactor) BridgeBurn(opts *bind.TransactOpts, user common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _L1Gateway.contract.Transact(opts, "bridgeBurn", user, amount, l2TxHash)
}

// BridgeBurn is a paid mutator transaction binding the contract method 0x87e04042.
//
// Solidity: function bridgeBurn(address user, uint256 amount, bytes32 l2TxHash) returns()
func (_L1Gateway *L1GatewaySession) BridgeBurn(user common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _L1Gateway.Contract.BridgeBurn(&_L1Gateway.TransactOpts, user, amount, l2TxHash)
}

// BridgeBurn is a paid mutator transaction binding the contract method 0x87e04042.
//
// Solidity: function bridgeBurn(address user, uint256 amount, bytes32 l2TxHash) returns()
func (_L1Gateway *L1GatewayTransactorSession) BridgeBurn(user common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _L1Gateway.Contract.BridgeBurn(&_L1Gateway.TransactOpts, user, amount, l2TxHash)
}

// BridgeMint is a paid mutator transaction binding the contract method 0xf661e631.
//
// Solidity: function bridgeMint(address user, uint256 amount, bytes32 l2TxHash) returns()
func (_L1Gateway *L1GatewayTransactor) BridgeMint(opts *bind.TransactOpts, user common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _L1Gateway.contract.Transact(opts, "bridgeMint", user, amount, l2TxHash)
}

// BridgeMint is a paid mutator transaction binding the contract method 0xf661e631.
//
// Solidity: function bridgeMint(address user, uint256 amount, bytes32 l2TxHash) returns()
func (_L1Gateway *L1GatewaySession) BridgeMint(user common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _L1Gateway.Contract.BridgeMint(&_L1Gateway.TransactOpts, user, amount, l2TxHash)
}

// BridgeMint is a paid mutator transaction binding the contract method 0xf661e631.
//
// Solidity: function bridgeMint(address user, uint256 amount, bytes32 l2TxHash) returns()
func (_L1Gateway *L1GatewayTransactorSession) BridgeMint(user common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _L1Gateway.Contract.BridgeMint(&_L1Gateway.TransactOpts, user, amount, l2TxHash)
}

// DepositToL2 is a paid mutator transaction binding the contract method 0xe09addc4.
//
// Solidity: function depositToL2(uint256 amount) payable returns(uint256 depositId, uint256 ticketId)
func (_L1Gateway *L1GatewayTransactor) DepositToL2(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _L1Gateway.contract.Transact(opts, "depositToL2", amount)
}

// DepositToL2 is a paid mutator transaction binding the contract method 0xe09addc4.
//
// Solidity: function depositToL2(uint256 amount) payable returns(uint256 depositId, uint256 ticketId)
func (_L1Gateway *L1GatewaySession) DepositToL2(amount *big.Int) (*types.Transaction, error) {
	return _L1Gateway.Contract.DepositToL2(&_L1Gateway.TransactOpts, amount)
}

// DepositToL2 is a paid mutator transaction binding the contract method 0xe09addc4.
//
// Solidity: function depositToL2(uint256 amount) payable returns(uint256 depositId, uint256 ticketId)
func (_L1Gateway *L1GatewayTransactorSession) DepositToL2(amount *big.Int) (*types.Transaction, error) {
	return _L1Gateway.Contract.DepositToL2(&_L1Gateway.TransactOpts, amount)
}

// FinalizeWithdrawal is a paid mutator transaction binding the contract method 0x185b550c.
//
// Solidity: function finalizeWithdrawal(bytes32 withdrawalId, address user, uint256 amount) returns()
func (_L1Gateway *L1GatewayTransactor) FinalizeWithdrawal(opts *bind.TransactOpts, withdrawalId [32]byte, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L1Gateway.contract.Transact(opts, "finalizeWithdrawal", withdrawalId, user, amount)
}

// FinalizeWithdrawal is a paid mutator transaction binding the contract method 0x185b550c.
//
// Solidity: function finalizeWithdrawal(bytes32 withdrawalId, address user, uint256 amount) returns()
func (_L1Gateway *L1GatewaySession) FinalizeWithdrawal(withdrawalId [32]byte, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L1Gateway.Contract.FinalizeWithdrawal(&_L1Gateway.TransactOpts, withdrawalId, user, amount)
}

// FinalizeWithdrawal is a paid mutator transaction binding the contract method 0x185b550c.
//
// Solidity: function finalizeWithdrawal(bytes32 withdrawalId, address user, uint256 amount) returns()
func (_L1Gateway *L1GatewayTransactorSession) FinalizeWithdrawal(withdrawalId [32]byte, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _L1Gateway.Contract.FinalizeWithdrawal(&_L1Gateway.TransactOpts, withdrawalId, user, amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_L1Gateway *L1GatewayTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Gateway.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_L1Gateway *L1GatewaySession) Pause() (*types.Transaction, error) {
	return _L1Gateway.Contract.Pause(&_L1Gateway.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_L1Gateway *L1GatewayTransactorSession) Pause() (*types.Transaction, error) {
	return _L1Gateway.Contract.Pause(&_L1Gateway.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L1Gateway *L1GatewayTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Gateway.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L1Gateway *L1GatewaySession) RenounceOwnership() (*types.Transaction, error) {
	return _L1Gateway.Contract.RenounceOwnership(&_L1Gateway.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L1Gateway *L1GatewayTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _L1Gateway.Contract.RenounceOwnership(&_L1Gateway.TransactOpts)
}

// SetGasParameters is a paid mutator transaction binding the contract method 0x0629dc5e.
//
// Solidity: function setGasParameters(uint256 _maxSubmissionCost, uint256 _maxGas, uint256 _gasPriceBid) returns()
func (_L1Gateway *L1GatewayTransactor) SetGasParameters(opts *bind.TransactOpts, _maxSubmissionCost *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int) (*types.Transaction, error) {
	return _L1Gateway.contract.Transact(opts, "setGasParameters", _maxSubmissionCost, _maxGas, _gasPriceBid)
}

// SetGasParameters is a paid mutator transaction binding the contract method 0x0629dc5e.
//
// Solidity: function setGasParameters(uint256 _maxSubmissionCost, uint256 _maxGas, uint256 _gasPriceBid) returns()
func (_L1Gateway *L1GatewaySession) SetGasParameters(_maxSubmissionCost *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int) (*types.Transaction, error) {
	return _L1Gateway.Contract.SetGasParameters(&_L1Gateway.TransactOpts, _maxSubmissionCost, _maxGas, _gasPriceBid)
}

// SetGasParameters is a paid mutator transaction binding the contract method 0x0629dc5e.
//
// Solidity: function setGasParameters(uint256 _maxSubmissionCost, uint256 _maxGas, uint256 _gasPriceBid) returns()
func (_L1Gateway *L1GatewayTransactorSession) SetGasParameters(_maxSubmissionCost *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int) (*types.Transaction, error) {
	return _L1Gateway.Contract.SetGasParameters(&_L1Gateway.TransactOpts, _maxSubmissionCost, _maxGas, _gasPriceBid)
}

// SetInbox is a paid mutator transaction binding the contract method 0x53b60c4a.
//
// Solidity: function setInbox(address _inbox) returns()
func (_L1Gateway *L1GatewayTransactor) SetInbox(opts *bind.TransactOpts, _inbox common.Address) (*types.Transaction, error) {
	return _L1Gateway.contract.Transact(opts, "setInbox", _inbox)
}

// SetInbox is a paid mutator transaction binding the contract method 0x53b60c4a.
//
// Solidity: function setInbox(address _inbox) returns()
func (_L1Gateway *L1GatewaySession) SetInbox(_inbox common.Address) (*types.Transaction, error) {
	return _L1Gateway.Contract.SetInbox(&_L1Gateway.TransactOpts, _inbox)
}

// SetInbox is a paid mutator transaction binding the contract method 0x53b60c4a.
//
// Solidity: function setInbox(address _inbox) returns()
func (_L1Gateway *L1GatewayTransactorSession) SetInbox(_inbox common.Address) (*types.Transaction, error) {
	return _L1Gateway.Contract.SetInbox(&_L1Gateway.TransactOpts, _inbox)
}

// SetL2Gateway is a paid mutator transaction binding the contract method 0xc80a0f84.
//
// Solidity: function setL2Gateway(address _l2Gateway) returns()
func (_L1Gateway *L1GatewayTransactor) SetL2Gateway(opts *bind.TransactOpts, _l2Gateway common.Address) (*types.Transaction, error) {
	return _L1Gateway.contract.Transact(opts, "setL2Gateway", _l2Gateway)
}

// SetL2Gateway is a paid mutator transaction binding the contract method 0xc80a0f84.
//
// Solidity: function setL2Gateway(address _l2Gateway) returns()
func (_L1Gateway *L1GatewaySession) SetL2Gateway(_l2Gateway common.Address) (*types.Transaction, error) {
	return _L1Gateway.Contract.SetL2Gateway(&_L1Gateway.TransactOpts, _l2Gateway)
}

// SetL2Gateway is a paid mutator transaction binding the contract method 0xc80a0f84.
//
// Solidity: function setL2Gateway(address _l2Gateway) returns()
func (_L1Gateway *L1GatewayTransactorSession) SetL2Gateway(_l2Gateway common.Address) (*types.Transaction, error) {
	return _L1Gateway.Contract.SetL2Gateway(&_L1Gateway.TransactOpts, _l2Gateway)
}

// SetOutbox is a paid mutator transaction binding the contract method 0xff204f3b.
//
// Solidity: function setOutbox(address _outbox) returns()
func (_L1Gateway *L1GatewayTransactor) SetOutbox(opts *bind.TransactOpts, _outbox common.Address) (*types.Transaction, error) {
	return _L1Gateway.contract.Transact(opts, "setOutbox", _outbox)
}

// SetOutbox is a paid mutator transaction binding the contract method 0xff204f3b.
//
// Solidity: function setOutbox(address _outbox) returns()
func (_L1Gateway *L1GatewaySession) SetOutbox(_outbox common.Address) (*types.Transaction, error) {
	return _L1Gateway.Contract.SetOutbox(&_L1Gateway.TransactOpts, _outbox)
}

// SetOutbox is a paid mutator transaction binding the contract method 0xff204f3b.
//
// Solidity: function setOutbox(address _outbox) returns()
func (_L1Gateway *L1GatewayTransactorSession) SetOutbox(_outbox common.Address) (*types.Transaction, error) {
	return _L1Gateway.Contract.SetOutbox(&_L1Gateway.TransactOpts, _outbox)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1Gateway *L1GatewayTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _L1Gateway.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1Gateway *L1GatewaySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L1Gateway.Contract.TransferOwnership(&_L1Gateway.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1Gateway *L1GatewayTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L1Gateway.Contract.TransferOwnership(&_L1Gateway.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_L1Gateway *L1GatewayTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1Gateway.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_L1Gateway *L1GatewaySession) Unpause() (*types.Transaction, error) {
	return _L1Gateway.Contract.Unpause(&_L1Gateway.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_L1Gateway *L1GatewayTransactorSession) Unpause() (*types.Transaction, error) {
	return _L1Gateway.Contract.Unpause(&_L1Gateway.TransactOpts)
}

// L1GatewayDepositFinalizedIterator is returned from FilterDepositFinalized and is used to iterate over the raw logs and unpacked data for DepositFinalized events raised by the L1Gateway contract.
type L1GatewayDepositFinalizedIterator struct {
	Event *L1GatewayDepositFinalized // Event containing the contract specifics and raw log

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
func (it *L1GatewayDepositFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1GatewayDepositFinalized)
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
		it.Event = new(L1GatewayDepositFinalized)
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
func (it *L1GatewayDepositFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1GatewayDepositFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1GatewayDepositFinalized represents a DepositFinalized event raised by the L1Gateway contract.
type L1GatewayDepositFinalized struct {
	DepositId *big.Int
	User      common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDepositFinalized is a free log retrieval operation binding the contract event 0x502ddd372bcee89768aa5cf92094111e7cb768dfccde24817da356d1022fc311.
//
// Solidity: event DepositFinalized(uint256 indexed depositId, address indexed user, uint256 amount)
func (_L1Gateway *L1GatewayFilterer) FilterDepositFinalized(opts *bind.FilterOpts, depositId []*big.Int, user []common.Address) (*L1GatewayDepositFinalizedIterator, error) {

	var depositIdRule []interface{}
	for _, depositIdItem := range depositId {
		depositIdRule = append(depositIdRule, depositIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _L1Gateway.contract.FilterLogs(opts, "DepositFinalized", depositIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &L1GatewayDepositFinalizedIterator{contract: _L1Gateway.contract, event: "DepositFinalized", logs: logs, sub: sub}, nil
}

// WatchDepositFinalized is a free log subscription operation binding the contract event 0x502ddd372bcee89768aa5cf92094111e7cb768dfccde24817da356d1022fc311.
//
// Solidity: event DepositFinalized(uint256 indexed depositId, address indexed user, uint256 amount)
func (_L1Gateway *L1GatewayFilterer) WatchDepositFinalized(opts *bind.WatchOpts, sink chan<- *L1GatewayDepositFinalized, depositId []*big.Int, user []common.Address) (event.Subscription, error) {

	var depositIdRule []interface{}
	for _, depositIdItem := range depositId {
		depositIdRule = append(depositIdRule, depositIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _L1Gateway.contract.WatchLogs(opts, "DepositFinalized", depositIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1GatewayDepositFinalized)
				if err := _L1Gateway.contract.UnpackLog(event, "DepositFinalized", log); err != nil {
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

// ParseDepositFinalized is a log parse operation binding the contract event 0x502ddd372bcee89768aa5cf92094111e7cb768dfccde24817da356d1022fc311.
//
// Solidity: event DepositFinalized(uint256 indexed depositId, address indexed user, uint256 amount)
func (_L1Gateway *L1GatewayFilterer) ParseDepositFinalized(log types.Log) (*L1GatewayDepositFinalized, error) {
	event := new(L1GatewayDepositFinalized)
	if err := _L1Gateway.contract.UnpackLog(event, "DepositFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1GatewayDepositInitiatedIterator is returned from FilterDepositInitiated and is used to iterate over the raw logs and unpacked data for DepositInitiated events raised by the L1Gateway contract.
type L1GatewayDepositInitiatedIterator struct {
	Event *L1GatewayDepositInitiated // Event containing the contract specifics and raw log

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
func (it *L1GatewayDepositInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1GatewayDepositInitiated)
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
		it.Event = new(L1GatewayDepositInitiated)
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
func (it *L1GatewayDepositInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1GatewayDepositInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1GatewayDepositInitiated represents a DepositInitiated event raised by the L1Gateway contract.
type L1GatewayDepositInitiated struct {
	DepositId *big.Int
	User      common.Address
	Amount    *big.Int
	TicketId  *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDepositInitiated is a free log retrieval operation binding the contract event 0x8bf3c565a3ccbbed4c759fe404c97a3f3a900345d2224cd2a9ace35403641728.
//
// Solidity: event DepositInitiated(uint256 indexed depositId, address indexed user, uint256 amount, uint256 indexed ticketId)
func (_L1Gateway *L1GatewayFilterer) FilterDepositInitiated(opts *bind.FilterOpts, depositId []*big.Int, user []common.Address, ticketId []*big.Int) (*L1GatewayDepositInitiatedIterator, error) {

	var depositIdRule []interface{}
	for _, depositIdItem := range depositId {
		depositIdRule = append(depositIdRule, depositIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var ticketIdRule []interface{}
	for _, ticketIdItem := range ticketId {
		ticketIdRule = append(ticketIdRule, ticketIdItem)
	}

	logs, sub, err := _L1Gateway.contract.FilterLogs(opts, "DepositInitiated", depositIdRule, userRule, ticketIdRule)
	if err != nil {
		return nil, err
	}
	return &L1GatewayDepositInitiatedIterator{contract: _L1Gateway.contract, event: "DepositInitiated", logs: logs, sub: sub}, nil
}

// WatchDepositInitiated is a free log subscription operation binding the contract event 0x8bf3c565a3ccbbed4c759fe404c97a3f3a900345d2224cd2a9ace35403641728.
//
// Solidity: event DepositInitiated(uint256 indexed depositId, address indexed user, uint256 amount, uint256 indexed ticketId)
func (_L1Gateway *L1GatewayFilterer) WatchDepositInitiated(opts *bind.WatchOpts, sink chan<- *L1GatewayDepositInitiated, depositId []*big.Int, user []common.Address, ticketId []*big.Int) (event.Subscription, error) {

	var depositIdRule []interface{}
	for _, depositIdItem := range depositId {
		depositIdRule = append(depositIdRule, depositIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var ticketIdRule []interface{}
	for _, ticketIdItem := range ticketId {
		ticketIdRule = append(ticketIdRule, ticketIdItem)
	}

	logs, sub, err := _L1Gateway.contract.WatchLogs(opts, "DepositInitiated", depositIdRule, userRule, ticketIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1GatewayDepositInitiated)
				if err := _L1Gateway.contract.UnpackLog(event, "DepositInitiated", log); err != nil {
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

// ParseDepositInitiated is a log parse operation binding the contract event 0x8bf3c565a3ccbbed4c759fe404c97a3f3a900345d2224cd2a9ace35403641728.
//
// Solidity: event DepositInitiated(uint256 indexed depositId, address indexed user, uint256 amount, uint256 indexed ticketId)
func (_L1Gateway *L1GatewayFilterer) ParseDepositInitiated(log types.Log) (*L1GatewayDepositInitiated, error) {
	event := new(L1GatewayDepositInitiated)
	if err := _L1Gateway.contract.UnpackLog(event, "DepositInitiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1GatewayGasParametersUpdatedIterator is returned from FilterGasParametersUpdated and is used to iterate over the raw logs and unpacked data for GasParametersUpdated events raised by the L1Gateway contract.
type L1GatewayGasParametersUpdatedIterator struct {
	Event *L1GatewayGasParametersUpdated // Event containing the contract specifics and raw log

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
func (it *L1GatewayGasParametersUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1GatewayGasParametersUpdated)
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
		it.Event = new(L1GatewayGasParametersUpdated)
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
func (it *L1GatewayGasParametersUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1GatewayGasParametersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1GatewayGasParametersUpdated represents a GasParametersUpdated event raised by the L1Gateway contract.
type L1GatewayGasParametersUpdated struct {
	MaxSubmissionCost *big.Int
	MaxGas            *big.Int
	GasPriceBid       *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterGasParametersUpdated is a free log retrieval operation binding the contract event 0xa5801b007372d0623488a3aa478f88d8063ba2317686061e02cea1bba6e1eb2c.
//
// Solidity: event GasParametersUpdated(uint256 maxSubmissionCost, uint256 maxGas, uint256 gasPriceBid)
func (_L1Gateway *L1GatewayFilterer) FilterGasParametersUpdated(opts *bind.FilterOpts) (*L1GatewayGasParametersUpdatedIterator, error) {

	logs, sub, err := _L1Gateway.contract.FilterLogs(opts, "GasParametersUpdated")
	if err != nil {
		return nil, err
	}
	return &L1GatewayGasParametersUpdatedIterator{contract: _L1Gateway.contract, event: "GasParametersUpdated", logs: logs, sub: sub}, nil
}

// WatchGasParametersUpdated is a free log subscription operation binding the contract event 0xa5801b007372d0623488a3aa478f88d8063ba2317686061e02cea1bba6e1eb2c.
//
// Solidity: event GasParametersUpdated(uint256 maxSubmissionCost, uint256 maxGas, uint256 gasPriceBid)
func (_L1Gateway *L1GatewayFilterer) WatchGasParametersUpdated(opts *bind.WatchOpts, sink chan<- *L1GatewayGasParametersUpdated) (event.Subscription, error) {

	logs, sub, err := _L1Gateway.contract.WatchLogs(opts, "GasParametersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1GatewayGasParametersUpdated)
				if err := _L1Gateway.contract.UnpackLog(event, "GasParametersUpdated", log); err != nil {
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

// ParseGasParametersUpdated is a log parse operation binding the contract event 0xa5801b007372d0623488a3aa478f88d8063ba2317686061e02cea1bba6e1eb2c.
//
// Solidity: event GasParametersUpdated(uint256 maxSubmissionCost, uint256 maxGas, uint256 gasPriceBid)
func (_L1Gateway *L1GatewayFilterer) ParseGasParametersUpdated(log types.Log) (*L1GatewayGasParametersUpdated, error) {
	event := new(L1GatewayGasParametersUpdated)
	if err := _L1Gateway.contract.UnpackLog(event, "GasParametersUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1GatewayL2GatewayUpdatedIterator is returned from FilterL2GatewayUpdated and is used to iterate over the raw logs and unpacked data for L2GatewayUpdated events raised by the L1Gateway contract.
type L1GatewayL2GatewayUpdatedIterator struct {
	Event *L1GatewayL2GatewayUpdated // Event containing the contract specifics and raw log

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
func (it *L1GatewayL2GatewayUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1GatewayL2GatewayUpdated)
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
		it.Event = new(L1GatewayL2GatewayUpdated)
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
func (it *L1GatewayL2GatewayUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1GatewayL2GatewayUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1GatewayL2GatewayUpdated represents a L2GatewayUpdated event raised by the L1Gateway contract.
type L1GatewayL2GatewayUpdated struct {
	OldGateway common.Address
	NewGateway common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterL2GatewayUpdated is a free log retrieval operation binding the contract event 0x113186743aa005710d6b5693b32a224cb33334370ee31a6cc864a7b04f7d6005.
//
// Solidity: event L2GatewayUpdated(address indexed oldGateway, address indexed newGateway)
func (_L1Gateway *L1GatewayFilterer) FilterL2GatewayUpdated(opts *bind.FilterOpts, oldGateway []common.Address, newGateway []common.Address) (*L1GatewayL2GatewayUpdatedIterator, error) {

	var oldGatewayRule []interface{}
	for _, oldGatewayItem := range oldGateway {
		oldGatewayRule = append(oldGatewayRule, oldGatewayItem)
	}
	var newGatewayRule []interface{}
	for _, newGatewayItem := range newGateway {
		newGatewayRule = append(newGatewayRule, newGatewayItem)
	}

	logs, sub, err := _L1Gateway.contract.FilterLogs(opts, "L2GatewayUpdated", oldGatewayRule, newGatewayRule)
	if err != nil {
		return nil, err
	}
	return &L1GatewayL2GatewayUpdatedIterator{contract: _L1Gateway.contract, event: "L2GatewayUpdated", logs: logs, sub: sub}, nil
}

// WatchL2GatewayUpdated is a free log subscription operation binding the contract event 0x113186743aa005710d6b5693b32a224cb33334370ee31a6cc864a7b04f7d6005.
//
// Solidity: event L2GatewayUpdated(address indexed oldGateway, address indexed newGateway)
func (_L1Gateway *L1GatewayFilterer) WatchL2GatewayUpdated(opts *bind.WatchOpts, sink chan<- *L1GatewayL2GatewayUpdated, oldGateway []common.Address, newGateway []common.Address) (event.Subscription, error) {

	var oldGatewayRule []interface{}
	for _, oldGatewayItem := range oldGateway {
		oldGatewayRule = append(oldGatewayRule, oldGatewayItem)
	}
	var newGatewayRule []interface{}
	for _, newGatewayItem := range newGateway {
		newGatewayRule = append(newGatewayRule, newGatewayItem)
	}

	logs, sub, err := _L1Gateway.contract.WatchLogs(opts, "L2GatewayUpdated", oldGatewayRule, newGatewayRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1GatewayL2GatewayUpdated)
				if err := _L1Gateway.contract.UnpackLog(event, "L2GatewayUpdated", log); err != nil {
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

// ParseL2GatewayUpdated is a log parse operation binding the contract event 0x113186743aa005710d6b5693b32a224cb33334370ee31a6cc864a7b04f7d6005.
//
// Solidity: event L2GatewayUpdated(address indexed oldGateway, address indexed newGateway)
func (_L1Gateway *L1GatewayFilterer) ParseL2GatewayUpdated(log types.Log) (*L1GatewayL2GatewayUpdated, error) {
	event := new(L1GatewayL2GatewayUpdated)
	if err := _L1Gateway.contract.UnpackLog(event, "L2GatewayUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1GatewayOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the L1Gateway contract.
type L1GatewayOwnershipTransferredIterator struct {
	Event *L1GatewayOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *L1GatewayOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1GatewayOwnershipTransferred)
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
		it.Event = new(L1GatewayOwnershipTransferred)
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
func (it *L1GatewayOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1GatewayOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1GatewayOwnershipTransferred represents a OwnershipTransferred event raised by the L1Gateway contract.
type L1GatewayOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1Gateway *L1GatewayFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*L1GatewayOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1Gateway.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &L1GatewayOwnershipTransferredIterator{contract: _L1Gateway.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1Gateway *L1GatewayFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *L1GatewayOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1Gateway.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1GatewayOwnershipTransferred)
				if err := _L1Gateway.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_L1Gateway *L1GatewayFilterer) ParseOwnershipTransferred(log types.Log) (*L1GatewayOwnershipTransferred, error) {
	event := new(L1GatewayOwnershipTransferred)
	if err := _L1Gateway.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1GatewayWithdrawalExecutedIterator is returned from FilterWithdrawalExecuted and is used to iterate over the raw logs and unpacked data for WithdrawalExecuted events raised by the L1Gateway contract.
type L1GatewayWithdrawalExecutedIterator struct {
	Event *L1GatewayWithdrawalExecuted // Event containing the contract specifics and raw log

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
func (it *L1GatewayWithdrawalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1GatewayWithdrawalExecuted)
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
		it.Event = new(L1GatewayWithdrawalExecuted)
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
func (it *L1GatewayWithdrawalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1GatewayWithdrawalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1GatewayWithdrawalExecuted represents a WithdrawalExecuted event raised by the L1Gateway contract.
type L1GatewayWithdrawalExecuted struct {
	WithdrawalId [32]byte
	User         common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalExecuted is a free log retrieval operation binding the contract event 0x9b500c765f0aa508908d3e165168183a5cef5acd68573acc126ce58c80ac6b3c.
//
// Solidity: event WithdrawalExecuted(bytes32 indexed withdrawalId, address indexed user, uint256 amount)
func (_L1Gateway *L1GatewayFilterer) FilterWithdrawalExecuted(opts *bind.FilterOpts, withdrawalId [][32]byte, user []common.Address) (*L1GatewayWithdrawalExecutedIterator, error) {

	var withdrawalIdRule []interface{}
	for _, withdrawalIdItem := range withdrawalId {
		withdrawalIdRule = append(withdrawalIdRule, withdrawalIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _L1Gateway.contract.FilterLogs(opts, "WithdrawalExecuted", withdrawalIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &L1GatewayWithdrawalExecutedIterator{contract: _L1Gateway.contract, event: "WithdrawalExecuted", logs: logs, sub: sub}, nil
}

// WatchWithdrawalExecuted is a free log subscription operation binding the contract event 0x9b500c765f0aa508908d3e165168183a5cef5acd68573acc126ce58c80ac6b3c.
//
// Solidity: event WithdrawalExecuted(bytes32 indexed withdrawalId, address indexed user, uint256 amount)
func (_L1Gateway *L1GatewayFilterer) WatchWithdrawalExecuted(opts *bind.WatchOpts, sink chan<- *L1GatewayWithdrawalExecuted, withdrawalId [][32]byte, user []common.Address) (event.Subscription, error) {

	var withdrawalIdRule []interface{}
	for _, withdrawalIdItem := range withdrawalId {
		withdrawalIdRule = append(withdrawalIdRule, withdrawalIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _L1Gateway.contract.WatchLogs(opts, "WithdrawalExecuted", withdrawalIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1GatewayWithdrawalExecuted)
				if err := _L1Gateway.contract.UnpackLog(event, "WithdrawalExecuted", log); err != nil {
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

// ParseWithdrawalExecuted is a log parse operation binding the contract event 0x9b500c765f0aa508908d3e165168183a5cef5acd68573acc126ce58c80ac6b3c.
//
// Solidity: event WithdrawalExecuted(bytes32 indexed withdrawalId, address indexed user, uint256 amount)
func (_L1Gateway *L1GatewayFilterer) ParseWithdrawalExecuted(log types.Log) (*L1GatewayWithdrawalExecuted, error) {
	event := new(L1GatewayWithdrawalExecuted)
	if err := _L1Gateway.contract.UnpackLog(event, "WithdrawalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1GatewayWithdrawalInitiatedIterator is returned from FilterWithdrawalInitiated and is used to iterate over the raw logs and unpacked data for WithdrawalInitiated events raised by the L1Gateway contract.
type L1GatewayWithdrawalInitiatedIterator struct {
	Event *L1GatewayWithdrawalInitiated // Event containing the contract specifics and raw log

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
func (it *L1GatewayWithdrawalInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1GatewayWithdrawalInitiated)
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
		it.Event = new(L1GatewayWithdrawalInitiated)
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
func (it *L1GatewayWithdrawalInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1GatewayWithdrawalInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1GatewayWithdrawalInitiated represents a WithdrawalInitiated event raised by the L1Gateway contract.
type L1GatewayWithdrawalInitiated struct {
	WithdrawalId [32]byte
	User         common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalInitiated is a free log retrieval operation binding the contract event 0xf966226666bdb2bc2f79ad476eb39c82d8668fabd37e3d27998e5225d4d2a88a.
//
// Solidity: event WithdrawalInitiated(bytes32 indexed withdrawalId, address indexed user, uint256 amount)
func (_L1Gateway *L1GatewayFilterer) FilterWithdrawalInitiated(opts *bind.FilterOpts, withdrawalId [][32]byte, user []common.Address) (*L1GatewayWithdrawalInitiatedIterator, error) {

	var withdrawalIdRule []interface{}
	for _, withdrawalIdItem := range withdrawalId {
		withdrawalIdRule = append(withdrawalIdRule, withdrawalIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _L1Gateway.contract.FilterLogs(opts, "WithdrawalInitiated", withdrawalIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &L1GatewayWithdrawalInitiatedIterator{contract: _L1Gateway.contract, event: "WithdrawalInitiated", logs: logs, sub: sub}, nil
}

// WatchWithdrawalInitiated is a free log subscription operation binding the contract event 0xf966226666bdb2bc2f79ad476eb39c82d8668fabd37e3d27998e5225d4d2a88a.
//
// Solidity: event WithdrawalInitiated(bytes32 indexed withdrawalId, address indexed user, uint256 amount)
func (_L1Gateway *L1GatewayFilterer) WatchWithdrawalInitiated(opts *bind.WatchOpts, sink chan<- *L1GatewayWithdrawalInitiated, withdrawalId [][32]byte, user []common.Address) (event.Subscription, error) {

	var withdrawalIdRule []interface{}
	for _, withdrawalIdItem := range withdrawalId {
		withdrawalIdRule = append(withdrawalIdRule, withdrawalIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _L1Gateway.contract.WatchLogs(opts, "WithdrawalInitiated", withdrawalIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1GatewayWithdrawalInitiated)
				if err := _L1Gateway.contract.UnpackLog(event, "WithdrawalInitiated", log); err != nil {
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

// ParseWithdrawalInitiated is a log parse operation binding the contract event 0xf966226666bdb2bc2f79ad476eb39c82d8668fabd37e3d27998e5225d4d2a88a.
//
// Solidity: event WithdrawalInitiated(bytes32 indexed withdrawalId, address indexed user, uint256 amount)
func (_L1Gateway *L1GatewayFilterer) ParseWithdrawalInitiated(log types.Log) (*L1GatewayWithdrawalInitiated, error) {
	event := new(L1GatewayWithdrawalInitiated)
	if err := _L1Gateway.contract.UnpackLog(event, "WithdrawalInitiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
