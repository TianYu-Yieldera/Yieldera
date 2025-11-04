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

// LoyaltyUSDL1MetaData contains all meta data concerning the LoyaltyUSDL1 contract.
var LoyaltyUSDL1MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"l2TxHash\",\"type\":\"bytes32\"}],\"name\":\"BridgeBurn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"l2TxHash\",\"type\":\"bytes32\"}],\"name\":\"BridgeMint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"burner\",\"type\":\"address\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"DailyMintLimitUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pauser\",\"type\":\"address\"}],\"name\":\"EmergencyPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"unpauser\",\"type\":\"address\"}],\"name\":\"EmergencyUnpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldBridge\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newBridge\",\"type\":\"address\"}],\"name\":\"L2BridgeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BRIDGE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BURNER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_MINT_PER_TX\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"l2TxHash\",\"type\":\"bytes32\"}],\"name\":\"bridgeBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"l2TxHash\",\"type\":\"bytes32\"}],\"name\":\"bridgeMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"canBurn\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"canMint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dailyMintLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dailyMintedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRemainingDailyMintCapacity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTimeUntilLimitReset\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2Bridge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMintResetTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newLimit\",\"type\":\"uint256\"}],\"name\":\"setDailyMintLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l2Bridge\",\"type\":\"address\"}],\"name\":\"setL2Bridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupplyFormatted\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// LoyaltyUSDL1ABI is the input ABI used to generate the binding from.
// Deprecated: Use LoyaltyUSDL1MetaData.ABI instead.
var LoyaltyUSDL1ABI = LoyaltyUSDL1MetaData.ABI

// LoyaltyUSDL1 is an auto generated Go binding around an Ethereum contract.
type LoyaltyUSDL1 struct {
	LoyaltyUSDL1Caller     // Read-only binding to the contract
	LoyaltyUSDL1Transactor // Write-only binding to the contract
	LoyaltyUSDL1Filterer   // Log filterer for contract events
}

// LoyaltyUSDL1Caller is an auto generated read-only Go binding around an Ethereum contract.
type LoyaltyUSDL1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LoyaltyUSDL1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type LoyaltyUSDL1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LoyaltyUSDL1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LoyaltyUSDL1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LoyaltyUSDL1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LoyaltyUSDL1Session struct {
	Contract     *LoyaltyUSDL1     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LoyaltyUSDL1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LoyaltyUSDL1CallerSession struct {
	Contract *LoyaltyUSDL1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// LoyaltyUSDL1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LoyaltyUSDL1TransactorSession struct {
	Contract     *LoyaltyUSDL1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// LoyaltyUSDL1Raw is an auto generated low-level Go binding around an Ethereum contract.
type LoyaltyUSDL1Raw struct {
	Contract *LoyaltyUSDL1 // Generic contract binding to access the raw methods on
}

// LoyaltyUSDL1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LoyaltyUSDL1CallerRaw struct {
	Contract *LoyaltyUSDL1Caller // Generic read-only contract binding to access the raw methods on
}

// LoyaltyUSDL1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LoyaltyUSDL1TransactorRaw struct {
	Contract *LoyaltyUSDL1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewLoyaltyUSDL1 creates a new instance of LoyaltyUSDL1, bound to a specific deployed contract.
func NewLoyaltyUSDL1(address common.Address, backend bind.ContractBackend) (*LoyaltyUSDL1, error) {
	contract, err := bindLoyaltyUSDL1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1{LoyaltyUSDL1Caller: LoyaltyUSDL1Caller{contract: contract}, LoyaltyUSDL1Transactor: LoyaltyUSDL1Transactor{contract: contract}, LoyaltyUSDL1Filterer: LoyaltyUSDL1Filterer{contract: contract}}, nil
}

// NewLoyaltyUSDL1Caller creates a new read-only instance of LoyaltyUSDL1, bound to a specific deployed contract.
func NewLoyaltyUSDL1Caller(address common.Address, caller bind.ContractCaller) (*LoyaltyUSDL1Caller, error) {
	contract, err := bindLoyaltyUSDL1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1Caller{contract: contract}, nil
}

// NewLoyaltyUSDL1Transactor creates a new write-only instance of LoyaltyUSDL1, bound to a specific deployed contract.
func NewLoyaltyUSDL1Transactor(address common.Address, transactor bind.ContractTransactor) (*LoyaltyUSDL1Transactor, error) {
	contract, err := bindLoyaltyUSDL1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1Transactor{contract: contract}, nil
}

// NewLoyaltyUSDL1Filterer creates a new log filterer instance of LoyaltyUSDL1, bound to a specific deployed contract.
func NewLoyaltyUSDL1Filterer(address common.Address, filterer bind.ContractFilterer) (*LoyaltyUSDL1Filterer, error) {
	contract, err := bindLoyaltyUSDL1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1Filterer{contract: contract}, nil
}

// bindLoyaltyUSDL1 binds a generic wrapper to an already deployed contract.
func bindLoyaltyUSDL1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LoyaltyUSDL1MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LoyaltyUSDL1 *LoyaltyUSDL1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LoyaltyUSDL1.Contract.LoyaltyUSDL1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LoyaltyUSDL1 *LoyaltyUSDL1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.LoyaltyUSDL1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LoyaltyUSDL1 *LoyaltyUSDL1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.LoyaltyUSDL1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LoyaltyUSDL1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.contract.Transact(opts, method, params...)
}

// BRIDGEROLE is a free data retrieval call binding the contract method 0xb5bfddea.
//
// Solidity: function BRIDGE_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) BRIDGEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "BRIDGE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BRIDGEROLE is a free data retrieval call binding the contract method 0xb5bfddea.
//
// Solidity: function BRIDGE_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) BRIDGEROLE() ([32]byte, error) {
	return _LoyaltyUSDL1.Contract.BRIDGEROLE(&_LoyaltyUSDL1.CallOpts)
}

// BRIDGEROLE is a free data retrieval call binding the contract method 0xb5bfddea.
//
// Solidity: function BRIDGE_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) BRIDGEROLE() ([32]byte, error) {
	return _LoyaltyUSDL1.Contract.BRIDGEROLE(&_LoyaltyUSDL1.CallOpts)
}

// BURNERROLE is a free data retrieval call binding the contract method 0x282c51f3.
//
// Solidity: function BURNER_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) BURNERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "BURNER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BURNERROLE is a free data retrieval call binding the contract method 0x282c51f3.
//
// Solidity: function BURNER_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) BURNERROLE() ([32]byte, error) {
	return _LoyaltyUSDL1.Contract.BURNERROLE(&_LoyaltyUSDL1.CallOpts)
}

// BURNERROLE is a free data retrieval call binding the contract method 0x282c51f3.
//
// Solidity: function BURNER_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) BURNERROLE() ([32]byte, error) {
	return _LoyaltyUSDL1.Contract.BURNERROLE(&_LoyaltyUSDL1.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) DEFAULTADMINROLE() ([32]byte, error) {
	return _LoyaltyUSDL1.Contract.DEFAULTADMINROLE(&_LoyaltyUSDL1.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _LoyaltyUSDL1.Contract.DEFAULTADMINROLE(&_LoyaltyUSDL1.CallOpts)
}

// MAXMINTPERTX is a free data retrieval call binding the contract method 0x8ecad721.
//
// Solidity: function MAX_MINT_PER_TX() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) MAXMINTPERTX(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "MAX_MINT_PER_TX")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXMINTPERTX is a free data retrieval call binding the contract method 0x8ecad721.
//
// Solidity: function MAX_MINT_PER_TX() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) MAXMINTPERTX() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.MAXMINTPERTX(&_LoyaltyUSDL1.CallOpts)
}

// MAXMINTPERTX is a free data retrieval call binding the contract method 0x8ecad721.
//
// Solidity: function MAX_MINT_PER_TX() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) MAXMINTPERTX() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.MAXMINTPERTX(&_LoyaltyUSDL1.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) MINTERROLE() ([32]byte, error) {
	return _LoyaltyUSDL1.Contract.MINTERROLE(&_LoyaltyUSDL1.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) MINTERROLE() ([32]byte, error) {
	return _LoyaltyUSDL1.Contract.MINTERROLE(&_LoyaltyUSDL1.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) PAUSERROLE() ([32]byte, error) {
	return _LoyaltyUSDL1.Contract.PAUSERROLE(&_LoyaltyUSDL1.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) PAUSERROLE() ([32]byte, error) {
	return _LoyaltyUSDL1.Contract.PAUSERROLE(&_LoyaltyUSDL1.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.Allowance(&_LoyaltyUSDL1.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.Allowance(&_LoyaltyUSDL1.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.BalanceOf(&_LoyaltyUSDL1.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.BalanceOf(&_LoyaltyUSDL1.CallOpts, account)
}

// CanBurn is a free data retrieval call binding the contract method 0x3820a686.
//
// Solidity: function canBurn(address account) view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) CanBurn(opts *bind.CallOpts, account common.Address) (bool, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "canBurn", account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanBurn is a free data retrieval call binding the contract method 0x3820a686.
//
// Solidity: function canBurn(address account) view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) CanBurn(account common.Address) (bool, error) {
	return _LoyaltyUSDL1.Contract.CanBurn(&_LoyaltyUSDL1.CallOpts, account)
}

// CanBurn is a free data retrieval call binding the contract method 0x3820a686.
//
// Solidity: function canBurn(address account) view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) CanBurn(account common.Address) (bool, error) {
	return _LoyaltyUSDL1.Contract.CanBurn(&_LoyaltyUSDL1.CallOpts, account)
}

// CanMint is a free data retrieval call binding the contract method 0xc2ba4744.
//
// Solidity: function canMint(address account) view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) CanMint(opts *bind.CallOpts, account common.Address) (bool, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "canMint", account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanMint is a free data retrieval call binding the contract method 0xc2ba4744.
//
// Solidity: function canMint(address account) view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) CanMint(account common.Address) (bool, error) {
	return _LoyaltyUSDL1.Contract.CanMint(&_LoyaltyUSDL1.CallOpts, account)
}

// CanMint is a free data retrieval call binding the contract method 0xc2ba4744.
//
// Solidity: function canMint(address account) view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) CanMint(account common.Address) (bool, error) {
	return _LoyaltyUSDL1.Contract.CanMint(&_LoyaltyUSDL1.CallOpts, account)
}

// DailyMintLimit is a free data retrieval call binding the contract method 0x62680e4b.
//
// Solidity: function dailyMintLimit() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) DailyMintLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "dailyMintLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DailyMintLimit is a free data retrieval call binding the contract method 0x62680e4b.
//
// Solidity: function dailyMintLimit() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) DailyMintLimit() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.DailyMintLimit(&_LoyaltyUSDL1.CallOpts)
}

// DailyMintLimit is a free data retrieval call binding the contract method 0x62680e4b.
//
// Solidity: function dailyMintLimit() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) DailyMintLimit() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.DailyMintLimit(&_LoyaltyUSDL1.CallOpts)
}

// DailyMintedAmount is a free data retrieval call binding the contract method 0x71f110ff.
//
// Solidity: function dailyMintedAmount() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) DailyMintedAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "dailyMintedAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DailyMintedAmount is a free data retrieval call binding the contract method 0x71f110ff.
//
// Solidity: function dailyMintedAmount() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) DailyMintedAmount() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.DailyMintedAmount(&_LoyaltyUSDL1.CallOpts)
}

// DailyMintedAmount is a free data retrieval call binding the contract method 0x71f110ff.
//
// Solidity: function dailyMintedAmount() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) DailyMintedAmount() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.DailyMintedAmount(&_LoyaltyUSDL1.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) Decimals() (uint8, error) {
	return _LoyaltyUSDL1.Contract.Decimals(&_LoyaltyUSDL1.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) Decimals() (uint8, error) {
	return _LoyaltyUSDL1.Contract.Decimals(&_LoyaltyUSDL1.CallOpts)
}

// GetRemainingDailyMintCapacity is a free data retrieval call binding the contract method 0xd8fe4e9f.
//
// Solidity: function getRemainingDailyMintCapacity() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) GetRemainingDailyMintCapacity(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "getRemainingDailyMintCapacity")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRemainingDailyMintCapacity is a free data retrieval call binding the contract method 0xd8fe4e9f.
//
// Solidity: function getRemainingDailyMintCapacity() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) GetRemainingDailyMintCapacity() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.GetRemainingDailyMintCapacity(&_LoyaltyUSDL1.CallOpts)
}

// GetRemainingDailyMintCapacity is a free data retrieval call binding the contract method 0xd8fe4e9f.
//
// Solidity: function getRemainingDailyMintCapacity() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) GetRemainingDailyMintCapacity() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.GetRemainingDailyMintCapacity(&_LoyaltyUSDL1.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _LoyaltyUSDL1.Contract.GetRoleAdmin(&_LoyaltyUSDL1.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _LoyaltyUSDL1.Contract.GetRoleAdmin(&_LoyaltyUSDL1.CallOpts, role)
}

// GetTimeUntilLimitReset is a free data retrieval call binding the contract method 0xb3aac012.
//
// Solidity: function getTimeUntilLimitReset() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) GetTimeUntilLimitReset(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "getTimeUntilLimitReset")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTimeUntilLimitReset is a free data retrieval call binding the contract method 0xb3aac012.
//
// Solidity: function getTimeUntilLimitReset() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) GetTimeUntilLimitReset() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.GetTimeUntilLimitReset(&_LoyaltyUSDL1.CallOpts)
}

// GetTimeUntilLimitReset is a free data retrieval call binding the contract method 0xb3aac012.
//
// Solidity: function getTimeUntilLimitReset() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) GetTimeUntilLimitReset() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.GetTimeUntilLimitReset(&_LoyaltyUSDL1.CallOpts)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _LoyaltyUSDL1.Contract.HasRole(&_LoyaltyUSDL1.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _LoyaltyUSDL1.Contract.HasRole(&_LoyaltyUSDL1.CallOpts, role, account)
}

// L2Bridge is a free data retrieval call binding the contract method 0xae1f6aaf.
//
// Solidity: function l2Bridge() view returns(address)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) L2Bridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "l2Bridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2Bridge is a free data retrieval call binding the contract method 0xae1f6aaf.
//
// Solidity: function l2Bridge() view returns(address)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) L2Bridge() (common.Address, error) {
	return _LoyaltyUSDL1.Contract.L2Bridge(&_LoyaltyUSDL1.CallOpts)
}

// L2Bridge is a free data retrieval call binding the contract method 0xae1f6aaf.
//
// Solidity: function l2Bridge() view returns(address)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) L2Bridge() (common.Address, error) {
	return _LoyaltyUSDL1.Contract.L2Bridge(&_LoyaltyUSDL1.CallOpts)
}

// LastMintResetTime is a free data retrieval call binding the contract method 0x9f2d8b80.
//
// Solidity: function lastMintResetTime() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) LastMintResetTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "lastMintResetTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastMintResetTime is a free data retrieval call binding the contract method 0x9f2d8b80.
//
// Solidity: function lastMintResetTime() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) LastMintResetTime() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.LastMintResetTime(&_LoyaltyUSDL1.CallOpts)
}

// LastMintResetTime is a free data retrieval call binding the contract method 0x9f2d8b80.
//
// Solidity: function lastMintResetTime() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) LastMintResetTime() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.LastMintResetTime(&_LoyaltyUSDL1.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) Name() (string, error) {
	return _LoyaltyUSDL1.Contract.Name(&_LoyaltyUSDL1.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) Name() (string, error) {
	return _LoyaltyUSDL1.Contract.Name(&_LoyaltyUSDL1.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) Paused() (bool, error) {
	return _LoyaltyUSDL1.Contract.Paused(&_LoyaltyUSDL1.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) Paused() (bool, error) {
	return _LoyaltyUSDL1.Contract.Paused(&_LoyaltyUSDL1.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LoyaltyUSDL1.Contract.SupportsInterface(&_LoyaltyUSDL1.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LoyaltyUSDL1.Contract.SupportsInterface(&_LoyaltyUSDL1.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) Symbol() (string, error) {
	return _LoyaltyUSDL1.Contract.Symbol(&_LoyaltyUSDL1.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) Symbol() (string, error) {
	return _LoyaltyUSDL1.Contract.Symbol(&_LoyaltyUSDL1.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) TotalSupply() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.TotalSupply(&_LoyaltyUSDL1.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) TotalSupply() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.TotalSupply(&_LoyaltyUSDL1.CallOpts)
}

// TotalSupplyFormatted is a free data retrieval call binding the contract method 0x936aeef3.
//
// Solidity: function totalSupplyFormatted() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Caller) TotalSupplyFormatted(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoyaltyUSDL1.contract.Call(opts, &out, "totalSupplyFormatted")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupplyFormatted is a free data retrieval call binding the contract method 0x936aeef3.
//
// Solidity: function totalSupplyFormatted() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) TotalSupplyFormatted() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.TotalSupplyFormatted(&_LoyaltyUSDL1.CallOpts)
}

// TotalSupplyFormatted is a free data retrieval call binding the contract method 0x936aeef3.
//
// Solidity: function totalSupplyFormatted() view returns(uint256)
func (_LoyaltyUSDL1 *LoyaltyUSDL1CallerSession) TotalSupplyFormatted() (*big.Int, error) {
	return _LoyaltyUSDL1.Contract.TotalSupplyFormatted(&_LoyaltyUSDL1.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Approve(&_LoyaltyUSDL1.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Approve(&_LoyaltyUSDL1.TransactOpts, spender, value)
}

// BridgeBurn is a paid mutator transaction binding the contract method 0x87e04042.
//
// Solidity: function bridgeBurn(address from, uint256 amount, bytes32 l2TxHash) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) BridgeBurn(opts *bind.TransactOpts, from common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "bridgeBurn", from, amount, l2TxHash)
}

// BridgeBurn is a paid mutator transaction binding the contract method 0x87e04042.
//
// Solidity: function bridgeBurn(address from, uint256 amount, bytes32 l2TxHash) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) BridgeBurn(from common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.BridgeBurn(&_LoyaltyUSDL1.TransactOpts, from, amount, l2TxHash)
}

// BridgeBurn is a paid mutator transaction binding the contract method 0x87e04042.
//
// Solidity: function bridgeBurn(address from, uint256 amount, bytes32 l2TxHash) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) BridgeBurn(from common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.BridgeBurn(&_LoyaltyUSDL1.TransactOpts, from, amount, l2TxHash)
}

// BridgeMint is a paid mutator transaction binding the contract method 0xf661e631.
//
// Solidity: function bridgeMint(address to, uint256 amount, bytes32 l2TxHash) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) BridgeMint(opts *bind.TransactOpts, to common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "bridgeMint", to, amount, l2TxHash)
}

// BridgeMint is a paid mutator transaction binding the contract method 0xf661e631.
//
// Solidity: function bridgeMint(address to, uint256 amount, bytes32 l2TxHash) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) BridgeMint(to common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.BridgeMint(&_LoyaltyUSDL1.TransactOpts, to, amount, l2TxHash)
}

// BridgeMint is a paid mutator transaction binding the contract method 0xf661e631.
//
// Solidity: function bridgeMint(address to, uint256 amount, bytes32 l2TxHash) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) BridgeMint(to common.Address, amount *big.Int, l2TxHash [32]byte) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.BridgeMint(&_LoyaltyUSDL1.TransactOpts, to, amount, l2TxHash)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) Burn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "burn", value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) Burn(value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Burn(&_LoyaltyUSDL1.TransactOpts, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Burn(&_LoyaltyUSDL1.TransactOpts, value)
}

// Burn0 is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) Burn0(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "burn0", from, amount)
}

// Burn0 is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) Burn0(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Burn0(&_LoyaltyUSDL1.TransactOpts, from, amount)
}

// Burn0 is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) Burn0(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Burn0(&_LoyaltyUSDL1.TransactOpts, from, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) BurnFrom(opts *bind.TransactOpts, account common.Address, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "burnFrom", account, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) BurnFrom(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.BurnFrom(&_LoyaltyUSDL1.TransactOpts, account, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) BurnFrom(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.BurnFrom(&_LoyaltyUSDL1.TransactOpts, account, value)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.GrantRole(&_LoyaltyUSDL1.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.GrantRole(&_LoyaltyUSDL1.TransactOpts, role, account)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Mint(&_LoyaltyUSDL1.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Mint(&_LoyaltyUSDL1.TransactOpts, to, amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) Pause() (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Pause(&_LoyaltyUSDL1.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) Pause() (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Pause(&_LoyaltyUSDL1.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.RenounceRole(&_LoyaltyUSDL1.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.RenounceRole(&_LoyaltyUSDL1.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.RevokeRole(&_LoyaltyUSDL1.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.RevokeRole(&_LoyaltyUSDL1.TransactOpts, role, account)
}

// SetDailyMintLimit is a paid mutator transaction binding the contract method 0xb2d52d27.
//
// Solidity: function setDailyMintLimit(uint256 newLimit) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) SetDailyMintLimit(opts *bind.TransactOpts, newLimit *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "setDailyMintLimit", newLimit)
}

// SetDailyMintLimit is a paid mutator transaction binding the contract method 0xb2d52d27.
//
// Solidity: function setDailyMintLimit(uint256 newLimit) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) SetDailyMintLimit(newLimit *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.SetDailyMintLimit(&_LoyaltyUSDL1.TransactOpts, newLimit)
}

// SetDailyMintLimit is a paid mutator transaction binding the contract method 0xb2d52d27.
//
// Solidity: function setDailyMintLimit(uint256 newLimit) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) SetDailyMintLimit(newLimit *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.SetDailyMintLimit(&_LoyaltyUSDL1.TransactOpts, newLimit)
}

// SetL2Bridge is a paid mutator transaction binding the contract method 0x3d36d971.
//
// Solidity: function setL2Bridge(address _l2Bridge) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) SetL2Bridge(opts *bind.TransactOpts, _l2Bridge common.Address) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "setL2Bridge", _l2Bridge)
}

// SetL2Bridge is a paid mutator transaction binding the contract method 0x3d36d971.
//
// Solidity: function setL2Bridge(address _l2Bridge) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) SetL2Bridge(_l2Bridge common.Address) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.SetL2Bridge(&_LoyaltyUSDL1.TransactOpts, _l2Bridge)
}

// SetL2Bridge is a paid mutator transaction binding the contract method 0x3d36d971.
//
// Solidity: function setL2Bridge(address _l2Bridge) returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) SetL2Bridge(_l2Bridge common.Address) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.SetL2Bridge(&_LoyaltyUSDL1.TransactOpts, _l2Bridge)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Transfer(&_LoyaltyUSDL1.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Transfer(&_LoyaltyUSDL1.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.TransferFrom(&_LoyaltyUSDL1.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.TransferFrom(&_LoyaltyUSDL1.TransactOpts, from, to, value)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoyaltyUSDL1.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1Session) Unpause() (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Unpause(&_LoyaltyUSDL1.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_LoyaltyUSDL1 *LoyaltyUSDL1TransactorSession) Unpause() (*types.Transaction, error) {
	return _LoyaltyUSDL1.Contract.Unpause(&_LoyaltyUSDL1.TransactOpts)
}

// LoyaltyUSDL1ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1ApprovalIterator struct {
	Event *LoyaltyUSDL1Approval // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1Approval)
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
		it.Event = new(LoyaltyUSDL1Approval)
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
func (it *LoyaltyUSDL1ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1Approval represents a Approval event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*LoyaltyUSDL1ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1ApprovalIterator{contract: _LoyaltyUSDL1.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1Approval)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseApproval(log types.Log) (*LoyaltyUSDL1Approval, error) {
	event := new(LoyaltyUSDL1Approval)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1BridgeBurnIterator is returned from FilterBridgeBurn and is used to iterate over the raw logs and unpacked data for BridgeBurn events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1BridgeBurnIterator struct {
	Event *LoyaltyUSDL1BridgeBurn // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1BridgeBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1BridgeBurn)
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
		it.Event = new(LoyaltyUSDL1BridgeBurn)
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
func (it *LoyaltyUSDL1BridgeBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1BridgeBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1BridgeBurn represents a BridgeBurn event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1BridgeBurn struct {
	From     common.Address
	Amount   *big.Int
	L2TxHash [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterBridgeBurn is a free log retrieval operation binding the contract event 0x7b7c7eab98a871dfe9496db9802be5d6175a2c454530fa7f467b98cc06895101.
//
// Solidity: event BridgeBurn(address indexed from, uint256 amount, bytes32 indexed l2TxHash)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterBridgeBurn(opts *bind.FilterOpts, from []common.Address, l2TxHash [][32]byte) (*LoyaltyUSDL1BridgeBurnIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	var l2TxHashRule []interface{}
	for _, l2TxHashItem := range l2TxHash {
		l2TxHashRule = append(l2TxHashRule, l2TxHashItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "BridgeBurn", fromRule, l2TxHashRule)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1BridgeBurnIterator{contract: _LoyaltyUSDL1.contract, event: "BridgeBurn", logs: logs, sub: sub}, nil
}

// WatchBridgeBurn is a free log subscription operation binding the contract event 0x7b7c7eab98a871dfe9496db9802be5d6175a2c454530fa7f467b98cc06895101.
//
// Solidity: event BridgeBurn(address indexed from, uint256 amount, bytes32 indexed l2TxHash)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchBridgeBurn(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1BridgeBurn, from []common.Address, l2TxHash [][32]byte) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	var l2TxHashRule []interface{}
	for _, l2TxHashItem := range l2TxHash {
		l2TxHashRule = append(l2TxHashRule, l2TxHashItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "BridgeBurn", fromRule, l2TxHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1BridgeBurn)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "BridgeBurn", log); err != nil {
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

// ParseBridgeBurn is a log parse operation binding the contract event 0x7b7c7eab98a871dfe9496db9802be5d6175a2c454530fa7f467b98cc06895101.
//
// Solidity: event BridgeBurn(address indexed from, uint256 amount, bytes32 indexed l2TxHash)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseBridgeBurn(log types.Log) (*LoyaltyUSDL1BridgeBurn, error) {
	event := new(LoyaltyUSDL1BridgeBurn)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "BridgeBurn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1BridgeMintIterator is returned from FilterBridgeMint and is used to iterate over the raw logs and unpacked data for BridgeMint events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1BridgeMintIterator struct {
	Event *LoyaltyUSDL1BridgeMint // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1BridgeMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1BridgeMint)
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
		it.Event = new(LoyaltyUSDL1BridgeMint)
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
func (it *LoyaltyUSDL1BridgeMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1BridgeMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1BridgeMint represents a BridgeMint event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1BridgeMint struct {
	To       common.Address
	Amount   *big.Int
	L2TxHash [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterBridgeMint is a free log retrieval operation binding the contract event 0x096d591342eea46377f5d86fed8a91ab1e9be20a7c3f5d698cecc1d4f698680a.
//
// Solidity: event BridgeMint(address indexed to, uint256 amount, bytes32 indexed l2TxHash)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterBridgeMint(opts *bind.FilterOpts, to []common.Address, l2TxHash [][32]byte) (*LoyaltyUSDL1BridgeMintIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	var l2TxHashRule []interface{}
	for _, l2TxHashItem := range l2TxHash {
		l2TxHashRule = append(l2TxHashRule, l2TxHashItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "BridgeMint", toRule, l2TxHashRule)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1BridgeMintIterator{contract: _LoyaltyUSDL1.contract, event: "BridgeMint", logs: logs, sub: sub}, nil
}

// WatchBridgeMint is a free log subscription operation binding the contract event 0x096d591342eea46377f5d86fed8a91ab1e9be20a7c3f5d698cecc1d4f698680a.
//
// Solidity: event BridgeMint(address indexed to, uint256 amount, bytes32 indexed l2TxHash)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchBridgeMint(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1BridgeMint, to []common.Address, l2TxHash [][32]byte) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	var l2TxHashRule []interface{}
	for _, l2TxHashItem := range l2TxHash {
		l2TxHashRule = append(l2TxHashRule, l2TxHashItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "BridgeMint", toRule, l2TxHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1BridgeMint)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "BridgeMint", log); err != nil {
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

// ParseBridgeMint is a log parse operation binding the contract event 0x096d591342eea46377f5d86fed8a91ab1e9be20a7c3f5d698cecc1d4f698680a.
//
// Solidity: event BridgeMint(address indexed to, uint256 amount, bytes32 indexed l2TxHash)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseBridgeMint(log types.Log) (*LoyaltyUSDL1BridgeMint, error) {
	event := new(LoyaltyUSDL1BridgeMint)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "BridgeMint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1BurnedIterator is returned from FilterBurned and is used to iterate over the raw logs and unpacked data for Burned events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1BurnedIterator struct {
	Event *LoyaltyUSDL1Burned // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1BurnedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1Burned)
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
		it.Event = new(LoyaltyUSDL1Burned)
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
func (it *LoyaltyUSDL1BurnedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1BurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1Burned represents a Burned event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1Burned struct {
	From   common.Address
	Amount *big.Int
	Burner common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBurned is a free log retrieval operation binding the contract event 0x4bec313b95a1f7373890b5d9dce4c4737945f27093ca8f3e357bec1219299eaf.
//
// Solidity: event Burned(address indexed from, uint256 amount, address indexed burner)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterBurned(opts *bind.FilterOpts, from []common.Address, burner []common.Address) (*LoyaltyUSDL1BurnedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "Burned", fromRule, burnerRule)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1BurnedIterator{contract: _LoyaltyUSDL1.contract, event: "Burned", logs: logs, sub: sub}, nil
}

// WatchBurned is a free log subscription operation binding the contract event 0x4bec313b95a1f7373890b5d9dce4c4737945f27093ca8f3e357bec1219299eaf.
//
// Solidity: event Burned(address indexed from, uint256 amount, address indexed burner)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1Burned, from []common.Address, burner []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	var burnerRule []interface{}
	for _, burnerItem := range burner {
		burnerRule = append(burnerRule, burnerItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "Burned", fromRule, burnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1Burned)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "Burned", log); err != nil {
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

// ParseBurned is a log parse operation binding the contract event 0x4bec313b95a1f7373890b5d9dce4c4737945f27093ca8f3e357bec1219299eaf.
//
// Solidity: event Burned(address indexed from, uint256 amount, address indexed burner)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseBurned(log types.Log) (*LoyaltyUSDL1Burned, error) {
	event := new(LoyaltyUSDL1Burned)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1DailyMintLimitUpdatedIterator is returned from FilterDailyMintLimitUpdated and is used to iterate over the raw logs and unpacked data for DailyMintLimitUpdated events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1DailyMintLimitUpdatedIterator struct {
	Event *LoyaltyUSDL1DailyMintLimitUpdated // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1DailyMintLimitUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1DailyMintLimitUpdated)
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
		it.Event = new(LoyaltyUSDL1DailyMintLimitUpdated)
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
func (it *LoyaltyUSDL1DailyMintLimitUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1DailyMintLimitUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1DailyMintLimitUpdated represents a DailyMintLimitUpdated event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1DailyMintLimitUpdated struct {
	OldLimit *big.Int
	NewLimit *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDailyMintLimitUpdated is a free log retrieval operation binding the contract event 0x8cf975fdba36321e100760d71c0717a9286810fcc2b7b1154ee84ef8450444b3.
//
// Solidity: event DailyMintLimitUpdated(uint256 oldLimit, uint256 newLimit)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterDailyMintLimitUpdated(opts *bind.FilterOpts) (*LoyaltyUSDL1DailyMintLimitUpdatedIterator, error) {

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "DailyMintLimitUpdated")
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1DailyMintLimitUpdatedIterator{contract: _LoyaltyUSDL1.contract, event: "DailyMintLimitUpdated", logs: logs, sub: sub}, nil
}

// WatchDailyMintLimitUpdated is a free log subscription operation binding the contract event 0x8cf975fdba36321e100760d71c0717a9286810fcc2b7b1154ee84ef8450444b3.
//
// Solidity: event DailyMintLimitUpdated(uint256 oldLimit, uint256 newLimit)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchDailyMintLimitUpdated(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1DailyMintLimitUpdated) (event.Subscription, error) {

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "DailyMintLimitUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1DailyMintLimitUpdated)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "DailyMintLimitUpdated", log); err != nil {
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

// ParseDailyMintLimitUpdated is a log parse operation binding the contract event 0x8cf975fdba36321e100760d71c0717a9286810fcc2b7b1154ee84ef8450444b3.
//
// Solidity: event DailyMintLimitUpdated(uint256 oldLimit, uint256 newLimit)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseDailyMintLimitUpdated(log types.Log) (*LoyaltyUSDL1DailyMintLimitUpdated, error) {
	event := new(LoyaltyUSDL1DailyMintLimitUpdated)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "DailyMintLimitUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1EmergencyPausedIterator is returned from FilterEmergencyPaused and is used to iterate over the raw logs and unpacked data for EmergencyPaused events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1EmergencyPausedIterator struct {
	Event *LoyaltyUSDL1EmergencyPaused // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1EmergencyPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1EmergencyPaused)
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
		it.Event = new(LoyaltyUSDL1EmergencyPaused)
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
func (it *LoyaltyUSDL1EmergencyPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1EmergencyPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1EmergencyPaused represents a EmergencyPaused event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1EmergencyPaused struct {
	Pauser common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEmergencyPaused is a free log retrieval operation binding the contract event 0xb8fad2fa0ed7a383e747c309ef2c4391d7b65592a48893e57ccc1fab70791456.
//
// Solidity: event EmergencyPaused(address indexed pauser)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterEmergencyPaused(opts *bind.FilterOpts, pauser []common.Address) (*LoyaltyUSDL1EmergencyPausedIterator, error) {

	var pauserRule []interface{}
	for _, pauserItem := range pauser {
		pauserRule = append(pauserRule, pauserItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "EmergencyPaused", pauserRule)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1EmergencyPausedIterator{contract: _LoyaltyUSDL1.contract, event: "EmergencyPaused", logs: logs, sub: sub}, nil
}

// WatchEmergencyPaused is a free log subscription operation binding the contract event 0xb8fad2fa0ed7a383e747c309ef2c4391d7b65592a48893e57ccc1fab70791456.
//
// Solidity: event EmergencyPaused(address indexed pauser)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchEmergencyPaused(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1EmergencyPaused, pauser []common.Address) (event.Subscription, error) {

	var pauserRule []interface{}
	for _, pauserItem := range pauser {
		pauserRule = append(pauserRule, pauserItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "EmergencyPaused", pauserRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1EmergencyPaused)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "EmergencyPaused", log); err != nil {
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

// ParseEmergencyPaused is a log parse operation binding the contract event 0xb8fad2fa0ed7a383e747c309ef2c4391d7b65592a48893e57ccc1fab70791456.
//
// Solidity: event EmergencyPaused(address indexed pauser)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseEmergencyPaused(log types.Log) (*LoyaltyUSDL1EmergencyPaused, error) {
	event := new(LoyaltyUSDL1EmergencyPaused)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "EmergencyPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1EmergencyUnpausedIterator is returned from FilterEmergencyUnpaused and is used to iterate over the raw logs and unpacked data for EmergencyUnpaused events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1EmergencyUnpausedIterator struct {
	Event *LoyaltyUSDL1EmergencyUnpaused // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1EmergencyUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1EmergencyUnpaused)
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
		it.Event = new(LoyaltyUSDL1EmergencyUnpaused)
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
func (it *LoyaltyUSDL1EmergencyUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1EmergencyUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1EmergencyUnpaused represents a EmergencyUnpaused event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1EmergencyUnpaused struct {
	Unpauser common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterEmergencyUnpaused is a free log retrieval operation binding the contract event 0xf5cbf596165cc457b2cd92e8d8450827ee314968160a5696402d75766fc52caf.
//
// Solidity: event EmergencyUnpaused(address indexed unpauser)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterEmergencyUnpaused(opts *bind.FilterOpts, unpauser []common.Address) (*LoyaltyUSDL1EmergencyUnpausedIterator, error) {

	var unpauserRule []interface{}
	for _, unpauserItem := range unpauser {
		unpauserRule = append(unpauserRule, unpauserItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "EmergencyUnpaused", unpauserRule)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1EmergencyUnpausedIterator{contract: _LoyaltyUSDL1.contract, event: "EmergencyUnpaused", logs: logs, sub: sub}, nil
}

// WatchEmergencyUnpaused is a free log subscription operation binding the contract event 0xf5cbf596165cc457b2cd92e8d8450827ee314968160a5696402d75766fc52caf.
//
// Solidity: event EmergencyUnpaused(address indexed unpauser)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchEmergencyUnpaused(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1EmergencyUnpaused, unpauser []common.Address) (event.Subscription, error) {

	var unpauserRule []interface{}
	for _, unpauserItem := range unpauser {
		unpauserRule = append(unpauserRule, unpauserItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "EmergencyUnpaused", unpauserRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1EmergencyUnpaused)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "EmergencyUnpaused", log); err != nil {
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

// ParseEmergencyUnpaused is a log parse operation binding the contract event 0xf5cbf596165cc457b2cd92e8d8450827ee314968160a5696402d75766fc52caf.
//
// Solidity: event EmergencyUnpaused(address indexed unpauser)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseEmergencyUnpaused(log types.Log) (*LoyaltyUSDL1EmergencyUnpaused, error) {
	event := new(LoyaltyUSDL1EmergencyUnpaused)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "EmergencyUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1L2BridgeUpdatedIterator is returned from FilterL2BridgeUpdated and is used to iterate over the raw logs and unpacked data for L2BridgeUpdated events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1L2BridgeUpdatedIterator struct {
	Event *LoyaltyUSDL1L2BridgeUpdated // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1L2BridgeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1L2BridgeUpdated)
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
		it.Event = new(LoyaltyUSDL1L2BridgeUpdated)
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
func (it *LoyaltyUSDL1L2BridgeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1L2BridgeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1L2BridgeUpdated represents a L2BridgeUpdated event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1L2BridgeUpdated struct {
	OldBridge common.Address
	NewBridge common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterL2BridgeUpdated is a free log retrieval operation binding the contract event 0xfbeaf70015f6ad12833bc6726f03af9bb91a7ee52827c6413c6c250c77854c95.
//
// Solidity: event L2BridgeUpdated(address indexed oldBridge, address indexed newBridge)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterL2BridgeUpdated(opts *bind.FilterOpts, oldBridge []common.Address, newBridge []common.Address) (*LoyaltyUSDL1L2BridgeUpdatedIterator, error) {

	var oldBridgeRule []interface{}
	for _, oldBridgeItem := range oldBridge {
		oldBridgeRule = append(oldBridgeRule, oldBridgeItem)
	}
	var newBridgeRule []interface{}
	for _, newBridgeItem := range newBridge {
		newBridgeRule = append(newBridgeRule, newBridgeItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "L2BridgeUpdated", oldBridgeRule, newBridgeRule)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1L2BridgeUpdatedIterator{contract: _LoyaltyUSDL1.contract, event: "L2BridgeUpdated", logs: logs, sub: sub}, nil
}

// WatchL2BridgeUpdated is a free log subscription operation binding the contract event 0xfbeaf70015f6ad12833bc6726f03af9bb91a7ee52827c6413c6c250c77854c95.
//
// Solidity: event L2BridgeUpdated(address indexed oldBridge, address indexed newBridge)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchL2BridgeUpdated(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1L2BridgeUpdated, oldBridge []common.Address, newBridge []common.Address) (event.Subscription, error) {

	var oldBridgeRule []interface{}
	for _, oldBridgeItem := range oldBridge {
		oldBridgeRule = append(oldBridgeRule, oldBridgeItem)
	}
	var newBridgeRule []interface{}
	for _, newBridgeItem := range newBridge {
		newBridgeRule = append(newBridgeRule, newBridgeItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "L2BridgeUpdated", oldBridgeRule, newBridgeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1L2BridgeUpdated)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "L2BridgeUpdated", log); err != nil {
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

// ParseL2BridgeUpdated is a log parse operation binding the contract event 0xfbeaf70015f6ad12833bc6726f03af9bb91a7ee52827c6413c6c250c77854c95.
//
// Solidity: event L2BridgeUpdated(address indexed oldBridge, address indexed newBridge)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseL2BridgeUpdated(log types.Log) (*LoyaltyUSDL1L2BridgeUpdated, error) {
	event := new(LoyaltyUSDL1L2BridgeUpdated)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "L2BridgeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1MintedIterator is returned from FilterMinted and is used to iterate over the raw logs and unpacked data for Minted events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1MintedIterator struct {
	Event *LoyaltyUSDL1Minted // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1MintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1Minted)
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
		it.Event = new(LoyaltyUSDL1Minted)
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
func (it *LoyaltyUSDL1MintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1MintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1Minted represents a Minted event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1Minted struct {
	To     common.Address
	Amount *big.Int
	Minter common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinted is a free log retrieval operation binding the contract event 0xfd1ba5519feb014ad2f631cbd516bd3cec96be935749391554d122972593e7bf.
//
// Solidity: event Minted(address indexed to, uint256 amount, address indexed minter)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterMinted(opts *bind.FilterOpts, to []common.Address, minter []common.Address) (*LoyaltyUSDL1MintedIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "Minted", toRule, minterRule)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1MintedIterator{contract: _LoyaltyUSDL1.contract, event: "Minted", logs: logs, sub: sub}, nil
}

// WatchMinted is a free log subscription operation binding the contract event 0xfd1ba5519feb014ad2f631cbd516bd3cec96be935749391554d122972593e7bf.
//
// Solidity: event Minted(address indexed to, uint256 amount, address indexed minter)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1Minted, to []common.Address, minter []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	var minterRule []interface{}
	for _, minterItem := range minter {
		minterRule = append(minterRule, minterItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "Minted", toRule, minterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1Minted)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "Minted", log); err != nil {
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

// ParseMinted is a log parse operation binding the contract event 0xfd1ba5519feb014ad2f631cbd516bd3cec96be935749391554d122972593e7bf.
//
// Solidity: event Minted(address indexed to, uint256 amount, address indexed minter)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseMinted(log types.Log) (*LoyaltyUSDL1Minted, error) {
	event := new(LoyaltyUSDL1Minted)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1PausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1PausedIterator struct {
	Event *LoyaltyUSDL1Paused // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1PausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1Paused)
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
		it.Event = new(LoyaltyUSDL1Paused)
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
func (it *LoyaltyUSDL1PausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1Paused represents a Paused event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1Paused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterPaused(opts *bind.FilterOpts) (*LoyaltyUSDL1PausedIterator, error) {

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1PausedIterator{contract: _LoyaltyUSDL1.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1Paused) (event.Subscription, error) {

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1Paused)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParsePaused(log types.Log) (*LoyaltyUSDL1Paused, error) {
	event := new(LoyaltyUSDL1Paused)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1RoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1RoleAdminChangedIterator struct {
	Event *LoyaltyUSDL1RoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1RoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1RoleAdminChanged)
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
		it.Event = new(LoyaltyUSDL1RoleAdminChanged)
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
func (it *LoyaltyUSDL1RoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1RoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1RoleAdminChanged represents a RoleAdminChanged event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1RoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*LoyaltyUSDL1RoleAdminChangedIterator, error) {

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

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1RoleAdminChangedIterator{contract: _LoyaltyUSDL1.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1RoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1RoleAdminChanged)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseRoleAdminChanged(log types.Log) (*LoyaltyUSDL1RoleAdminChanged, error) {
	event := new(LoyaltyUSDL1RoleAdminChanged)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1RoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1RoleGrantedIterator struct {
	Event *LoyaltyUSDL1RoleGranted // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1RoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1RoleGranted)
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
		it.Event = new(LoyaltyUSDL1RoleGranted)
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
func (it *LoyaltyUSDL1RoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1RoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1RoleGranted represents a RoleGranted event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1RoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*LoyaltyUSDL1RoleGrantedIterator, error) {

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

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1RoleGrantedIterator{contract: _LoyaltyUSDL1.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1RoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1RoleGranted)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseRoleGranted(log types.Log) (*LoyaltyUSDL1RoleGranted, error) {
	event := new(LoyaltyUSDL1RoleGranted)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1RoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1RoleRevokedIterator struct {
	Event *LoyaltyUSDL1RoleRevoked // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1RoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1RoleRevoked)
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
		it.Event = new(LoyaltyUSDL1RoleRevoked)
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
func (it *LoyaltyUSDL1RoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1RoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1RoleRevoked represents a RoleRevoked event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1RoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*LoyaltyUSDL1RoleRevokedIterator, error) {

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

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1RoleRevokedIterator{contract: _LoyaltyUSDL1.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1RoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1RoleRevoked)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseRoleRevoked(log types.Log) (*LoyaltyUSDL1RoleRevoked, error) {
	event := new(LoyaltyUSDL1RoleRevoked)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1TransferIterator struct {
	Event *LoyaltyUSDL1Transfer // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1Transfer)
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
		it.Event = new(LoyaltyUSDL1Transfer)
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
func (it *LoyaltyUSDL1TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1Transfer represents a Transfer event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LoyaltyUSDL1TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1TransferIterator{contract: _LoyaltyUSDL1.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1Transfer)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseTransfer(log types.Log) (*LoyaltyUSDL1Transfer, error) {
	event := new(LoyaltyUSDL1Transfer)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoyaltyUSDL1UnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1UnpausedIterator struct {
	Event *LoyaltyUSDL1Unpaused // Event containing the contract specifics and raw log

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
func (it *LoyaltyUSDL1UnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoyaltyUSDL1Unpaused)
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
		it.Event = new(LoyaltyUSDL1Unpaused)
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
func (it *LoyaltyUSDL1UnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoyaltyUSDL1UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoyaltyUSDL1Unpaused represents a Unpaused event raised by the LoyaltyUSDL1 contract.
type LoyaltyUSDL1Unpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) FilterUnpaused(opts *bind.FilterOpts) (*LoyaltyUSDL1UnpausedIterator, error) {

	logs, sub, err := _LoyaltyUSDL1.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &LoyaltyUSDL1UnpausedIterator{contract: _LoyaltyUSDL1.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *LoyaltyUSDL1Unpaused) (event.Subscription, error) {

	logs, sub, err := _LoyaltyUSDL1.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoyaltyUSDL1Unpaused)
				if err := _LoyaltyUSDL1.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_LoyaltyUSDL1 *LoyaltyUSDL1Filterer) ParseUnpaused(log types.Log) (*LoyaltyUSDL1Unpaused, error) {
	event := new(LoyaltyUSDL1Unpaused)
	if err := _LoyaltyUSDL1.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
