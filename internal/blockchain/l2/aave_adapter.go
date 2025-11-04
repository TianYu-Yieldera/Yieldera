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

// AaveV3AdapterUserPosition is an auto generated low-level Go binding around an user-defined struct.
type AaveV3AdapterUserPosition struct {
	TotalSupplied *big.Int
	TotalBorrowed *big.Int
	LastUpdate    *big.Int
}

// IAaveV3PoolReserveConfigurationMap is an auto generated low-level Go binding around an user-defined struct.
type IAaveV3PoolReserveConfigurationMap struct {
	Data *big.Int
}

// IAaveV3PoolReserveData is an auto generated low-level Go binding around an user-defined struct.
type IAaveV3PoolReserveData struct {
	Configuration               IAaveV3PoolReserveConfigurationMap
	LiquidityIndex              *big.Int
	CurrentLiquidityRate        *big.Int
	VariableBorrowIndex         *big.Int
	CurrentVariableBorrowRate   *big.Int
	CurrentStableBorrowRate     *big.Int
	LastUpdateTimestamp         *big.Int
	Id                          uint16
	ATokenAddress               common.Address
	StableDebtTokenAddress      common.Address
	VariableDebtTokenAddress    common.Address
	InterestRateStrategyAddress common.Address
	AccruedToTreasury           *big.Int
	Unbacked                    *big.Int
	IsolationModeTotalDebt      *big.Int
}

// AaveV3AdapterMetaData contains all meta data concerning the AaveV3Adapter contract.
var AaveV3AdapterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aavePool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_stateAggregator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"interestRateMode\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"Borrowed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"useAsCollateral\",\"type\":\"bool\"}],\"name\":\"CollateralStatusChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"FlashLoanExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMode\",\"type\":\"uint256\"}],\"name\":\"InterestRateModeSwapped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"interestRateMode\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"Repaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"Supplied\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FLASH_LOAN_PREMIUM\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INTEREST_RATE_MODE_STABLE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INTEREST_RATE_MODE_VARIABLE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PREMIUM_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REFERRAL_CODE\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"aavePool\",\"outputs\":[{\"internalType\":\"contractIAaveV3Pool\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activeUsers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"authorizedFlashLoanReceivers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRateMode\",\"type\":\"uint256\"}],\"name\":\"borrow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"calculateFlashLoanPremium\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"premium\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"executeOperation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"premiums\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"executeOperation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"assets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"interestRateModes\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"flashLoan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"flashLoanSimple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"getReserveData\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"internalType\":\"structIAaveV3Pool.ReserveConfigurationMap\",\"name\":\"configuration\",\"type\":\"tuple\"},{\"internalType\":\"uint128\",\"name\":\"liquidityIndex\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"currentLiquidityRate\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"variableBorrowIndex\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"currentVariableBorrowRate\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"currentStableBorrowRate\",\"type\":\"uint128\"},{\"internalType\":\"uint40\",\"name\":\"lastUpdateTimestamp\",\"type\":\"uint40\"},{\"internalType\":\"uint16\",\"name\":\"id\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"aTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stableDebtTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"variableDebtTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"interestRateStrategyAddress\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"accruedToTreasury\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"unbacked\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"isolationModeTotalDebt\",\"type\":\"uint128\"}],\"internalType\":\"structIAaveV3Pool.ReserveData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserAccountData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalCollateralBase\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalDebtBase\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"availableBorrowsBase\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentLiquidationThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ltv\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"healthFactor\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserPosition\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"totalSupplied\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalBorrowed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdate\",\"type\":\"uint256\"}],\"internalType\":\"structAaveV3Adapter.UserPosition\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"recoverToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRateMode\",\"type\":\"uint256\"}],\"name\":\"repay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"useAsCollateral\",\"type\":\"bool\"}],\"name\":\"setUserUseReserveAsCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateAggregator\",\"outputs\":[{\"internalType\":\"contractL2StateAggregator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"supply\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"currentRateMode\",\"type\":\"uint256\"}],\"name\":\"swapBorrowRateMode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalBorrowed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupplied\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newAggregator\",\"type\":\"address\"}],\"name\":\"updateStateAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"userPositions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalSupplied\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalBorrowed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AaveV3AdapterABI is the input ABI used to generate the binding from.
// Deprecated: Use AaveV3AdapterMetaData.ABI instead.
var AaveV3AdapterABI = AaveV3AdapterMetaData.ABI

// AaveV3Adapter is an auto generated Go binding around an Ethereum contract.
type AaveV3Adapter struct {
	AaveV3AdapterCaller     // Read-only binding to the contract
	AaveV3AdapterTransactor // Write-only binding to the contract
	AaveV3AdapterFilterer   // Log filterer for contract events
}

// AaveV3AdapterCaller is an auto generated read-only Go binding around an Ethereum contract.
type AaveV3AdapterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AaveV3AdapterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AaveV3AdapterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AaveV3AdapterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AaveV3AdapterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AaveV3AdapterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AaveV3AdapterSession struct {
	Contract     *AaveV3Adapter    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AaveV3AdapterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AaveV3AdapterCallerSession struct {
	Contract *AaveV3AdapterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// AaveV3AdapterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AaveV3AdapterTransactorSession struct {
	Contract     *AaveV3AdapterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// AaveV3AdapterRaw is an auto generated low-level Go binding around an Ethereum contract.
type AaveV3AdapterRaw struct {
	Contract *AaveV3Adapter // Generic contract binding to access the raw methods on
}

// AaveV3AdapterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AaveV3AdapterCallerRaw struct {
	Contract *AaveV3AdapterCaller // Generic read-only contract binding to access the raw methods on
}

// AaveV3AdapterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AaveV3AdapterTransactorRaw struct {
	Contract *AaveV3AdapterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAaveV3Adapter creates a new instance of AaveV3Adapter, bound to a specific deployed contract.
func NewAaveV3Adapter(address common.Address, backend bind.ContractBackend) (*AaveV3Adapter, error) {
	contract, err := bindAaveV3Adapter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AaveV3Adapter{AaveV3AdapterCaller: AaveV3AdapterCaller{contract: contract}, AaveV3AdapterTransactor: AaveV3AdapterTransactor{contract: contract}, AaveV3AdapterFilterer: AaveV3AdapterFilterer{contract: contract}}, nil
}

// NewAaveV3AdapterCaller creates a new read-only instance of AaveV3Adapter, bound to a specific deployed contract.
func NewAaveV3AdapterCaller(address common.Address, caller bind.ContractCaller) (*AaveV3AdapterCaller, error) {
	contract, err := bindAaveV3Adapter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AaveV3AdapterCaller{contract: contract}, nil
}

// NewAaveV3AdapterTransactor creates a new write-only instance of AaveV3Adapter, bound to a specific deployed contract.
func NewAaveV3AdapterTransactor(address common.Address, transactor bind.ContractTransactor) (*AaveV3AdapterTransactor, error) {
	contract, err := bindAaveV3Adapter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AaveV3AdapterTransactor{contract: contract}, nil
}

// NewAaveV3AdapterFilterer creates a new log filterer instance of AaveV3Adapter, bound to a specific deployed contract.
func NewAaveV3AdapterFilterer(address common.Address, filterer bind.ContractFilterer) (*AaveV3AdapterFilterer, error) {
	contract, err := bindAaveV3Adapter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AaveV3AdapterFilterer{contract: contract}, nil
}

// bindAaveV3Adapter binds a generic wrapper to an already deployed contract.
func bindAaveV3Adapter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AaveV3AdapterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AaveV3Adapter *AaveV3AdapterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AaveV3Adapter.Contract.AaveV3AdapterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AaveV3Adapter *AaveV3AdapterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.AaveV3AdapterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AaveV3Adapter *AaveV3AdapterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.AaveV3AdapterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AaveV3Adapter *AaveV3AdapterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AaveV3Adapter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AaveV3Adapter *AaveV3AdapterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AaveV3Adapter *AaveV3AdapterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.contract.Transact(opts, method, params...)
}

// FLASHLOANPREMIUM is a free data retrieval call binding the contract method 0xe549ff21.
//
// Solidity: function FLASH_LOAN_PREMIUM() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCaller) FLASHLOANPREMIUM(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "FLASH_LOAN_PREMIUM")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FLASHLOANPREMIUM is a free data retrieval call binding the contract method 0xe549ff21.
//
// Solidity: function FLASH_LOAN_PREMIUM() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterSession) FLASHLOANPREMIUM() (*big.Int, error) {
	return _AaveV3Adapter.Contract.FLASHLOANPREMIUM(&_AaveV3Adapter.CallOpts)
}

// FLASHLOANPREMIUM is a free data retrieval call binding the contract method 0xe549ff21.
//
// Solidity: function FLASH_LOAN_PREMIUM() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) FLASHLOANPREMIUM() (*big.Int, error) {
	return _AaveV3Adapter.Contract.FLASHLOANPREMIUM(&_AaveV3Adapter.CallOpts)
}

// INTERESTRATEMODESTABLE is a free data retrieval call binding the contract method 0x6591e3a2.
//
// Solidity: function INTEREST_RATE_MODE_STABLE() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCaller) INTERESTRATEMODESTABLE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "INTEREST_RATE_MODE_STABLE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INTERESTRATEMODESTABLE is a free data retrieval call binding the contract method 0x6591e3a2.
//
// Solidity: function INTEREST_RATE_MODE_STABLE() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterSession) INTERESTRATEMODESTABLE() (*big.Int, error) {
	return _AaveV3Adapter.Contract.INTERESTRATEMODESTABLE(&_AaveV3Adapter.CallOpts)
}

// INTERESTRATEMODESTABLE is a free data retrieval call binding the contract method 0x6591e3a2.
//
// Solidity: function INTEREST_RATE_MODE_STABLE() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) INTERESTRATEMODESTABLE() (*big.Int, error) {
	return _AaveV3Adapter.Contract.INTERESTRATEMODESTABLE(&_AaveV3Adapter.CallOpts)
}

// INTERESTRATEMODEVARIABLE is a free data retrieval call binding the contract method 0x6669e405.
//
// Solidity: function INTEREST_RATE_MODE_VARIABLE() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCaller) INTERESTRATEMODEVARIABLE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "INTEREST_RATE_MODE_VARIABLE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INTERESTRATEMODEVARIABLE is a free data retrieval call binding the contract method 0x6669e405.
//
// Solidity: function INTEREST_RATE_MODE_VARIABLE() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterSession) INTERESTRATEMODEVARIABLE() (*big.Int, error) {
	return _AaveV3Adapter.Contract.INTERESTRATEMODEVARIABLE(&_AaveV3Adapter.CallOpts)
}

// INTERESTRATEMODEVARIABLE is a free data retrieval call binding the contract method 0x6669e405.
//
// Solidity: function INTEREST_RATE_MODE_VARIABLE() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) INTERESTRATEMODEVARIABLE() (*big.Int, error) {
	return _AaveV3Adapter.Contract.INTERESTRATEMODEVARIABLE(&_AaveV3Adapter.CallOpts)
}

// PREMIUMPRECISION is a free data retrieval call binding the contract method 0xedec34e4.
//
// Solidity: function PREMIUM_PRECISION() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCaller) PREMIUMPRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "PREMIUM_PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PREMIUMPRECISION is a free data retrieval call binding the contract method 0xedec34e4.
//
// Solidity: function PREMIUM_PRECISION() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterSession) PREMIUMPRECISION() (*big.Int, error) {
	return _AaveV3Adapter.Contract.PREMIUMPRECISION(&_AaveV3Adapter.CallOpts)
}

// PREMIUMPRECISION is a free data retrieval call binding the contract method 0xedec34e4.
//
// Solidity: function PREMIUM_PRECISION() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) PREMIUMPRECISION() (*big.Int, error) {
	return _AaveV3Adapter.Contract.PREMIUMPRECISION(&_AaveV3Adapter.CallOpts)
}

// REFERRALCODE is a free data retrieval call binding the contract method 0x3583849a.
//
// Solidity: function REFERRAL_CODE() view returns(uint16)
func (_AaveV3Adapter *AaveV3AdapterCaller) REFERRALCODE(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "REFERRAL_CODE")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// REFERRALCODE is a free data retrieval call binding the contract method 0x3583849a.
//
// Solidity: function REFERRAL_CODE() view returns(uint16)
func (_AaveV3Adapter *AaveV3AdapterSession) REFERRALCODE() (uint16, error) {
	return _AaveV3Adapter.Contract.REFERRALCODE(&_AaveV3Adapter.CallOpts)
}

// REFERRALCODE is a free data retrieval call binding the contract method 0x3583849a.
//
// Solidity: function REFERRAL_CODE() view returns(uint16)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) REFERRALCODE() (uint16, error) {
	return _AaveV3Adapter.Contract.REFERRALCODE(&_AaveV3Adapter.CallOpts)
}

// AavePool is a free data retrieval call binding the contract method 0xa03e4bc3.
//
// Solidity: function aavePool() view returns(address)
func (_AaveV3Adapter *AaveV3AdapterCaller) AavePool(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "aavePool")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AavePool is a free data retrieval call binding the contract method 0xa03e4bc3.
//
// Solidity: function aavePool() view returns(address)
func (_AaveV3Adapter *AaveV3AdapterSession) AavePool() (common.Address, error) {
	return _AaveV3Adapter.Contract.AavePool(&_AaveV3Adapter.CallOpts)
}

// AavePool is a free data retrieval call binding the contract method 0xa03e4bc3.
//
// Solidity: function aavePool() view returns(address)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) AavePool() (common.Address, error) {
	return _AaveV3Adapter.Contract.AavePool(&_AaveV3Adapter.CallOpts)
}

// ActiveUsers is a free data retrieval call binding the contract method 0x74aa34de.
//
// Solidity: function activeUsers() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCaller) ActiveUsers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "activeUsers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveUsers is a free data retrieval call binding the contract method 0x74aa34de.
//
// Solidity: function activeUsers() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterSession) ActiveUsers() (*big.Int, error) {
	return _AaveV3Adapter.Contract.ActiveUsers(&_AaveV3Adapter.CallOpts)
}

// ActiveUsers is a free data retrieval call binding the contract method 0x74aa34de.
//
// Solidity: function activeUsers() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) ActiveUsers() (*big.Int, error) {
	return _AaveV3Adapter.Contract.ActiveUsers(&_AaveV3Adapter.CallOpts)
}

// AuthorizedFlashLoanReceivers is a free data retrieval call binding the contract method 0xe55bd4ea.
//
// Solidity: function authorizedFlashLoanReceivers(address ) view returns(bool)
func (_AaveV3Adapter *AaveV3AdapterCaller) AuthorizedFlashLoanReceivers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "authorizedFlashLoanReceivers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AuthorizedFlashLoanReceivers is a free data retrieval call binding the contract method 0xe55bd4ea.
//
// Solidity: function authorizedFlashLoanReceivers(address ) view returns(bool)
func (_AaveV3Adapter *AaveV3AdapterSession) AuthorizedFlashLoanReceivers(arg0 common.Address) (bool, error) {
	return _AaveV3Adapter.Contract.AuthorizedFlashLoanReceivers(&_AaveV3Adapter.CallOpts, arg0)
}

// AuthorizedFlashLoanReceivers is a free data retrieval call binding the contract method 0xe55bd4ea.
//
// Solidity: function authorizedFlashLoanReceivers(address ) view returns(bool)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) AuthorizedFlashLoanReceivers(arg0 common.Address) (bool, error) {
	return _AaveV3Adapter.Contract.AuthorizedFlashLoanReceivers(&_AaveV3Adapter.CallOpts, arg0)
}

// CalculateFlashLoanPremium is a free data retrieval call binding the contract method 0x87095873.
//
// Solidity: function calculateFlashLoanPremium(uint256 amount) pure returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCaller) CalculateFlashLoanPremium(opts *bind.CallOpts, amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "calculateFlashLoanPremium", amount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateFlashLoanPremium is a free data retrieval call binding the contract method 0x87095873.
//
// Solidity: function calculateFlashLoanPremium(uint256 amount) pure returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterSession) CalculateFlashLoanPremium(amount *big.Int) (*big.Int, error) {
	return _AaveV3Adapter.Contract.CalculateFlashLoanPremium(&_AaveV3Adapter.CallOpts, amount)
}

// CalculateFlashLoanPremium is a free data retrieval call binding the contract method 0x87095873.
//
// Solidity: function calculateFlashLoanPremium(uint256 amount) pure returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) CalculateFlashLoanPremium(amount *big.Int) (*big.Int, error) {
	return _AaveV3Adapter.Contract.CalculateFlashLoanPremium(&_AaveV3Adapter.CallOpts, amount)
}

// GetReserveData is a free data retrieval call binding the contract method 0x35ea6a75.
//
// Solidity: function getReserveData(address asset) view returns(((uint256),uint128,uint128,uint128,uint128,uint128,uint40,uint16,address,address,address,address,uint128,uint128,uint128))
func (_AaveV3Adapter *AaveV3AdapterCaller) GetReserveData(opts *bind.CallOpts, asset common.Address) (IAaveV3PoolReserveData, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "getReserveData", asset)

	if err != nil {
		return *new(IAaveV3PoolReserveData), err
	}

	out0 := *abi.ConvertType(out[0], new(IAaveV3PoolReserveData)).(*IAaveV3PoolReserveData)

	return out0, err

}

// GetReserveData is a free data retrieval call binding the contract method 0x35ea6a75.
//
// Solidity: function getReserveData(address asset) view returns(((uint256),uint128,uint128,uint128,uint128,uint128,uint40,uint16,address,address,address,address,uint128,uint128,uint128))
func (_AaveV3Adapter *AaveV3AdapterSession) GetReserveData(asset common.Address) (IAaveV3PoolReserveData, error) {
	return _AaveV3Adapter.Contract.GetReserveData(&_AaveV3Adapter.CallOpts, asset)
}

// GetReserveData is a free data retrieval call binding the contract method 0x35ea6a75.
//
// Solidity: function getReserveData(address asset) view returns(((uint256),uint128,uint128,uint128,uint128,uint128,uint40,uint16,address,address,address,address,uint128,uint128,uint128))
func (_AaveV3Adapter *AaveV3AdapterCallerSession) GetReserveData(asset common.Address) (IAaveV3PoolReserveData, error) {
	return _AaveV3Adapter.Contract.GetReserveData(&_AaveV3Adapter.CallOpts, asset)
}

// GetUserAccountData is a free data retrieval call binding the contract method 0xbf92857c.
//
// Solidity: function getUserAccountData(address user) view returns(uint256 totalCollateralBase, uint256 totalDebtBase, uint256 availableBorrowsBase, uint256 currentLiquidationThreshold, uint256 ltv, uint256 healthFactor)
func (_AaveV3Adapter *AaveV3AdapterCaller) GetUserAccountData(opts *bind.CallOpts, user common.Address) (struct {
	TotalCollateralBase         *big.Int
	TotalDebtBase               *big.Int
	AvailableBorrowsBase        *big.Int
	CurrentLiquidationThreshold *big.Int
	Ltv                         *big.Int
	HealthFactor                *big.Int
}, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "getUserAccountData", user)

	outstruct := new(struct {
		TotalCollateralBase         *big.Int
		TotalDebtBase               *big.Int
		AvailableBorrowsBase        *big.Int
		CurrentLiquidationThreshold *big.Int
		Ltv                         *big.Int
		HealthFactor                *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalCollateralBase = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalDebtBase = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.AvailableBorrowsBase = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.CurrentLiquidationThreshold = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Ltv = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.HealthFactor = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetUserAccountData is a free data retrieval call binding the contract method 0xbf92857c.
//
// Solidity: function getUserAccountData(address user) view returns(uint256 totalCollateralBase, uint256 totalDebtBase, uint256 availableBorrowsBase, uint256 currentLiquidationThreshold, uint256 ltv, uint256 healthFactor)
func (_AaveV3Adapter *AaveV3AdapterSession) GetUserAccountData(user common.Address) (struct {
	TotalCollateralBase         *big.Int
	TotalDebtBase               *big.Int
	AvailableBorrowsBase        *big.Int
	CurrentLiquidationThreshold *big.Int
	Ltv                         *big.Int
	HealthFactor                *big.Int
}, error) {
	return _AaveV3Adapter.Contract.GetUserAccountData(&_AaveV3Adapter.CallOpts, user)
}

// GetUserAccountData is a free data retrieval call binding the contract method 0xbf92857c.
//
// Solidity: function getUserAccountData(address user) view returns(uint256 totalCollateralBase, uint256 totalDebtBase, uint256 availableBorrowsBase, uint256 currentLiquidationThreshold, uint256 ltv, uint256 healthFactor)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) GetUserAccountData(user common.Address) (struct {
	TotalCollateralBase         *big.Int
	TotalDebtBase               *big.Int
	AvailableBorrowsBase        *big.Int
	CurrentLiquidationThreshold *big.Int
	Ltv                         *big.Int
	HealthFactor                *big.Int
}, error) {
	return _AaveV3Adapter.Contract.GetUserAccountData(&_AaveV3Adapter.CallOpts, user)
}

// GetUserPosition is a free data retrieval call binding the contract method 0x5b7c2dad.
//
// Solidity: function getUserPosition(address user) view returns((uint256,uint256,uint256))
func (_AaveV3Adapter *AaveV3AdapterCaller) GetUserPosition(opts *bind.CallOpts, user common.Address) (AaveV3AdapterUserPosition, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "getUserPosition", user)

	if err != nil {
		return *new(AaveV3AdapterUserPosition), err
	}

	out0 := *abi.ConvertType(out[0], new(AaveV3AdapterUserPosition)).(*AaveV3AdapterUserPosition)

	return out0, err

}

// GetUserPosition is a free data retrieval call binding the contract method 0x5b7c2dad.
//
// Solidity: function getUserPosition(address user) view returns((uint256,uint256,uint256))
func (_AaveV3Adapter *AaveV3AdapterSession) GetUserPosition(user common.Address) (AaveV3AdapterUserPosition, error) {
	return _AaveV3Adapter.Contract.GetUserPosition(&_AaveV3Adapter.CallOpts, user)
}

// GetUserPosition is a free data retrieval call binding the contract method 0x5b7c2dad.
//
// Solidity: function getUserPosition(address user) view returns((uint256,uint256,uint256))
func (_AaveV3Adapter *AaveV3AdapterCallerSession) GetUserPosition(user common.Address) (AaveV3AdapterUserPosition, error) {
	return _AaveV3Adapter.Contract.GetUserPosition(&_AaveV3Adapter.CallOpts, user)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AaveV3Adapter *AaveV3AdapterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AaveV3Adapter *AaveV3AdapterSession) Owner() (common.Address, error) {
	return _AaveV3Adapter.Contract.Owner(&_AaveV3Adapter.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) Owner() (common.Address, error) {
	return _AaveV3Adapter.Contract.Owner(&_AaveV3Adapter.CallOpts)
}

// StateAggregator is a free data retrieval call binding the contract method 0x54ef2920.
//
// Solidity: function stateAggregator() view returns(address)
func (_AaveV3Adapter *AaveV3AdapterCaller) StateAggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "stateAggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StateAggregator is a free data retrieval call binding the contract method 0x54ef2920.
//
// Solidity: function stateAggregator() view returns(address)
func (_AaveV3Adapter *AaveV3AdapterSession) StateAggregator() (common.Address, error) {
	return _AaveV3Adapter.Contract.StateAggregator(&_AaveV3Adapter.CallOpts)
}

// StateAggregator is a free data retrieval call binding the contract method 0x54ef2920.
//
// Solidity: function stateAggregator() view returns(address)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) StateAggregator() (common.Address, error) {
	return _AaveV3Adapter.Contract.StateAggregator(&_AaveV3Adapter.CallOpts)
}

// TotalBorrowed is a free data retrieval call binding the contract method 0x4c19386c.
//
// Solidity: function totalBorrowed() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCaller) TotalBorrowed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "totalBorrowed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalBorrowed is a free data retrieval call binding the contract method 0x4c19386c.
//
// Solidity: function totalBorrowed() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterSession) TotalBorrowed() (*big.Int, error) {
	return _AaveV3Adapter.Contract.TotalBorrowed(&_AaveV3Adapter.CallOpts)
}

// TotalBorrowed is a free data retrieval call binding the contract method 0x4c19386c.
//
// Solidity: function totalBorrowed() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) TotalBorrowed() (*big.Int, error) {
	return _AaveV3Adapter.Contract.TotalBorrowed(&_AaveV3Adapter.CallOpts)
}

// TotalSupplied is a free data retrieval call binding the contract method 0x630fd0ac.
//
// Solidity: function totalSupplied() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCaller) TotalSupplied(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "totalSupplied")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupplied is a free data retrieval call binding the contract method 0x630fd0ac.
//
// Solidity: function totalSupplied() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterSession) TotalSupplied() (*big.Int, error) {
	return _AaveV3Adapter.Contract.TotalSupplied(&_AaveV3Adapter.CallOpts)
}

// TotalSupplied is a free data retrieval call binding the contract method 0x630fd0ac.
//
// Solidity: function totalSupplied() view returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) TotalSupplied() (*big.Int, error) {
	return _AaveV3Adapter.Contract.TotalSupplied(&_AaveV3Adapter.CallOpts)
}

// UserPositions is a free data retrieval call binding the contract method 0x613cf420.
//
// Solidity: function userPositions(address ) view returns(uint256 totalSupplied, uint256 totalBorrowed, uint256 lastUpdate)
func (_AaveV3Adapter *AaveV3AdapterCaller) UserPositions(opts *bind.CallOpts, arg0 common.Address) (struct {
	TotalSupplied *big.Int
	TotalBorrowed *big.Int
	LastUpdate    *big.Int
}, error) {
	var out []interface{}
	err := _AaveV3Adapter.contract.Call(opts, &out, "userPositions", arg0)

	outstruct := new(struct {
		TotalSupplied *big.Int
		TotalBorrowed *big.Int
		LastUpdate    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalSupplied = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalBorrowed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.LastUpdate = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// UserPositions is a free data retrieval call binding the contract method 0x613cf420.
//
// Solidity: function userPositions(address ) view returns(uint256 totalSupplied, uint256 totalBorrowed, uint256 lastUpdate)
func (_AaveV3Adapter *AaveV3AdapterSession) UserPositions(arg0 common.Address) (struct {
	TotalSupplied *big.Int
	TotalBorrowed *big.Int
	LastUpdate    *big.Int
}, error) {
	return _AaveV3Adapter.Contract.UserPositions(&_AaveV3Adapter.CallOpts, arg0)
}

// UserPositions is a free data retrieval call binding the contract method 0x613cf420.
//
// Solidity: function userPositions(address ) view returns(uint256 totalSupplied, uint256 totalBorrowed, uint256 lastUpdate)
func (_AaveV3Adapter *AaveV3AdapterCallerSession) UserPositions(arg0 common.Address) (struct {
	TotalSupplied *big.Int
	TotalBorrowed *big.Int
	LastUpdate    *big.Int
}, error) {
	return _AaveV3Adapter.Contract.UserPositions(&_AaveV3Adapter.CallOpts, arg0)
}

// Borrow is a paid mutator transaction binding the contract method 0xc1bce0b7.
//
// Solidity: function borrow(address asset, uint256 amount, uint256 interestRateMode) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactor) Borrow(opts *bind.TransactOpts, asset common.Address, amount *big.Int, interestRateMode *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "borrow", asset, amount, interestRateMode)
}

// Borrow is a paid mutator transaction binding the contract method 0xc1bce0b7.
//
// Solidity: function borrow(address asset, uint256 amount, uint256 interestRateMode) returns()
func (_AaveV3Adapter *AaveV3AdapterSession) Borrow(asset common.Address, amount *big.Int, interestRateMode *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.Borrow(&_AaveV3Adapter.TransactOpts, asset, amount, interestRateMode)
}

// Borrow is a paid mutator transaction binding the contract method 0xc1bce0b7.
//
// Solidity: function borrow(address asset, uint256 amount, uint256 interestRateMode) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) Borrow(asset common.Address, amount *big.Int, interestRateMode *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.Borrow(&_AaveV3Adapter.TransactOpts, asset, amount, interestRateMode)
}

// ExecuteOperation is a paid mutator transaction binding the contract method 0x1b11d0ff.
//
// Solidity: function executeOperation(address asset, uint256 amount, uint256 premium, address initiator, bytes params) returns(bool)
func (_AaveV3Adapter *AaveV3AdapterTransactor) ExecuteOperation(opts *bind.TransactOpts, asset common.Address, amount *big.Int, premium *big.Int, initiator common.Address, params []byte) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "executeOperation", asset, amount, premium, initiator, params)
}

// ExecuteOperation is a paid mutator transaction binding the contract method 0x1b11d0ff.
//
// Solidity: function executeOperation(address asset, uint256 amount, uint256 premium, address initiator, bytes params) returns(bool)
func (_AaveV3Adapter *AaveV3AdapterSession) ExecuteOperation(asset common.Address, amount *big.Int, premium *big.Int, initiator common.Address, params []byte) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.ExecuteOperation(&_AaveV3Adapter.TransactOpts, asset, amount, premium, initiator, params)
}

// ExecuteOperation is a paid mutator transaction binding the contract method 0x1b11d0ff.
//
// Solidity: function executeOperation(address asset, uint256 amount, uint256 premium, address initiator, bytes params) returns(bool)
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) ExecuteOperation(asset common.Address, amount *big.Int, premium *big.Int, initiator common.Address, params []byte) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.ExecuteOperation(&_AaveV3Adapter.TransactOpts, asset, amount, premium, initiator, params)
}

// ExecuteOperation0 is a paid mutator transaction binding the contract method 0x920f5c84.
//
// Solidity: function executeOperation(address[] assets, uint256[] amounts, uint256[] premiums, address initiator, bytes params) returns(bool)
func (_AaveV3Adapter *AaveV3AdapterTransactor) ExecuteOperation0(opts *bind.TransactOpts, assets []common.Address, amounts []*big.Int, premiums []*big.Int, initiator common.Address, params []byte) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "executeOperation0", assets, amounts, premiums, initiator, params)
}

// ExecuteOperation0 is a paid mutator transaction binding the contract method 0x920f5c84.
//
// Solidity: function executeOperation(address[] assets, uint256[] amounts, uint256[] premiums, address initiator, bytes params) returns(bool)
func (_AaveV3Adapter *AaveV3AdapterSession) ExecuteOperation0(assets []common.Address, amounts []*big.Int, premiums []*big.Int, initiator common.Address, params []byte) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.ExecuteOperation0(&_AaveV3Adapter.TransactOpts, assets, amounts, premiums, initiator, params)
}

// ExecuteOperation0 is a paid mutator transaction binding the contract method 0x920f5c84.
//
// Solidity: function executeOperation(address[] assets, uint256[] amounts, uint256[] premiums, address initiator, bytes params) returns(bool)
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) ExecuteOperation0(assets []common.Address, amounts []*big.Int, premiums []*big.Int, initiator common.Address, params []byte) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.ExecuteOperation0(&_AaveV3Adapter.TransactOpts, assets, amounts, premiums, initiator, params)
}

// FlashLoan is a paid mutator transaction binding the contract method 0x54296154.
//
// Solidity: function flashLoan(address[] assets, uint256[] amounts, uint256[] interestRateModes, bytes params) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactor) FlashLoan(opts *bind.TransactOpts, assets []common.Address, amounts []*big.Int, interestRateModes []*big.Int, params []byte) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "flashLoan", assets, amounts, interestRateModes, params)
}

// FlashLoan is a paid mutator transaction binding the contract method 0x54296154.
//
// Solidity: function flashLoan(address[] assets, uint256[] amounts, uint256[] interestRateModes, bytes params) returns()
func (_AaveV3Adapter *AaveV3AdapterSession) FlashLoan(assets []common.Address, amounts []*big.Int, interestRateModes []*big.Int, params []byte) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.FlashLoan(&_AaveV3Adapter.TransactOpts, assets, amounts, interestRateModes, params)
}

// FlashLoan is a paid mutator transaction binding the contract method 0x54296154.
//
// Solidity: function flashLoan(address[] assets, uint256[] amounts, uint256[] interestRateModes, bytes params) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) FlashLoan(assets []common.Address, amounts []*big.Int, interestRateModes []*big.Int, params []byte) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.FlashLoan(&_AaveV3Adapter.TransactOpts, assets, amounts, interestRateModes, params)
}

// FlashLoanSimple is a paid mutator transaction binding the contract method 0x00666eea.
//
// Solidity: function flashLoanSimple(address asset, uint256 amount, bytes params) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactor) FlashLoanSimple(opts *bind.TransactOpts, asset common.Address, amount *big.Int, params []byte) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "flashLoanSimple", asset, amount, params)
}

// FlashLoanSimple is a paid mutator transaction binding the contract method 0x00666eea.
//
// Solidity: function flashLoanSimple(address asset, uint256 amount, bytes params) returns()
func (_AaveV3Adapter *AaveV3AdapterSession) FlashLoanSimple(asset common.Address, amount *big.Int, params []byte) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.FlashLoanSimple(&_AaveV3Adapter.TransactOpts, asset, amount, params)
}

// FlashLoanSimple is a paid mutator transaction binding the contract method 0x00666eea.
//
// Solidity: function flashLoanSimple(address asset, uint256 amount, bytes params) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) FlashLoanSimple(asset common.Address, amount *big.Int, params []byte) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.FlashLoanSimple(&_AaveV3Adapter.TransactOpts, asset, amount, params)
}

// RecoverToken is a paid mutator transaction binding the contract method 0xb29a8140.
//
// Solidity: function recoverToken(address token, uint256 amount) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactor) RecoverToken(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "recoverToken", token, amount)
}

// RecoverToken is a paid mutator transaction binding the contract method 0xb29a8140.
//
// Solidity: function recoverToken(address token, uint256 amount) returns()
func (_AaveV3Adapter *AaveV3AdapterSession) RecoverToken(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.RecoverToken(&_AaveV3Adapter.TransactOpts, token, amount)
}

// RecoverToken is a paid mutator transaction binding the contract method 0xb29a8140.
//
// Solidity: function recoverToken(address token, uint256 amount) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) RecoverToken(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.RecoverToken(&_AaveV3Adapter.TransactOpts, token, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AaveV3Adapter *AaveV3AdapterTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AaveV3Adapter *AaveV3AdapterSession) RenounceOwnership() (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.RenounceOwnership(&_AaveV3Adapter.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.RenounceOwnership(&_AaveV3Adapter.TransactOpts)
}

// Repay is a paid mutator transaction binding the contract method 0x8cd2e0c7.
//
// Solidity: function repay(address asset, uint256 amount, uint256 interestRateMode) returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterTransactor) Repay(opts *bind.TransactOpts, asset common.Address, amount *big.Int, interestRateMode *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "repay", asset, amount, interestRateMode)
}

// Repay is a paid mutator transaction binding the contract method 0x8cd2e0c7.
//
// Solidity: function repay(address asset, uint256 amount, uint256 interestRateMode) returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterSession) Repay(asset common.Address, amount *big.Int, interestRateMode *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.Repay(&_AaveV3Adapter.TransactOpts, asset, amount, interestRateMode)
}

// Repay is a paid mutator transaction binding the contract method 0x8cd2e0c7.
//
// Solidity: function repay(address asset, uint256 amount, uint256 interestRateMode) returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) Repay(asset common.Address, amount *big.Int, interestRateMode *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.Repay(&_AaveV3Adapter.TransactOpts, asset, amount, interestRateMode)
}

// SetUserUseReserveAsCollateral is a paid mutator transaction binding the contract method 0x5a3b74b9.
//
// Solidity: function setUserUseReserveAsCollateral(address asset, bool useAsCollateral) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactor) SetUserUseReserveAsCollateral(opts *bind.TransactOpts, asset common.Address, useAsCollateral bool) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "setUserUseReserveAsCollateral", asset, useAsCollateral)
}

// SetUserUseReserveAsCollateral is a paid mutator transaction binding the contract method 0x5a3b74b9.
//
// Solidity: function setUserUseReserveAsCollateral(address asset, bool useAsCollateral) returns()
func (_AaveV3Adapter *AaveV3AdapterSession) SetUserUseReserveAsCollateral(asset common.Address, useAsCollateral bool) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.SetUserUseReserveAsCollateral(&_AaveV3Adapter.TransactOpts, asset, useAsCollateral)
}

// SetUserUseReserveAsCollateral is a paid mutator transaction binding the contract method 0x5a3b74b9.
//
// Solidity: function setUserUseReserveAsCollateral(address asset, bool useAsCollateral) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) SetUserUseReserveAsCollateral(asset common.Address, useAsCollateral bool) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.SetUserUseReserveAsCollateral(&_AaveV3Adapter.TransactOpts, asset, useAsCollateral)
}

// Supply is a paid mutator transaction binding the contract method 0xf2b9fdb8.
//
// Solidity: function supply(address asset, uint256 amount) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactor) Supply(opts *bind.TransactOpts, asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "supply", asset, amount)
}

// Supply is a paid mutator transaction binding the contract method 0xf2b9fdb8.
//
// Solidity: function supply(address asset, uint256 amount) returns()
func (_AaveV3Adapter *AaveV3AdapterSession) Supply(asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.Supply(&_AaveV3Adapter.TransactOpts, asset, amount)
}

// Supply is a paid mutator transaction binding the contract method 0xf2b9fdb8.
//
// Solidity: function supply(address asset, uint256 amount) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) Supply(asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.Supply(&_AaveV3Adapter.TransactOpts, asset, amount)
}

// SwapBorrowRateMode is a paid mutator transaction binding the contract method 0x94ba89a2.
//
// Solidity: function swapBorrowRateMode(address asset, uint256 currentRateMode) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactor) SwapBorrowRateMode(opts *bind.TransactOpts, asset common.Address, currentRateMode *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "swapBorrowRateMode", asset, currentRateMode)
}

// SwapBorrowRateMode is a paid mutator transaction binding the contract method 0x94ba89a2.
//
// Solidity: function swapBorrowRateMode(address asset, uint256 currentRateMode) returns()
func (_AaveV3Adapter *AaveV3AdapterSession) SwapBorrowRateMode(asset common.Address, currentRateMode *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.SwapBorrowRateMode(&_AaveV3Adapter.TransactOpts, asset, currentRateMode)
}

// SwapBorrowRateMode is a paid mutator transaction binding the contract method 0x94ba89a2.
//
// Solidity: function swapBorrowRateMode(address asset, uint256 currentRateMode) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) SwapBorrowRateMode(asset common.Address, currentRateMode *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.SwapBorrowRateMode(&_AaveV3Adapter.TransactOpts, asset, currentRateMode)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AaveV3Adapter *AaveV3AdapterSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.TransferOwnership(&_AaveV3Adapter.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.TransferOwnership(&_AaveV3Adapter.TransactOpts, newOwner)
}

// UpdateStateAggregator is a paid mutator transaction binding the contract method 0x5ec71e77.
//
// Solidity: function updateStateAggregator(address _newAggregator) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactor) UpdateStateAggregator(opts *bind.TransactOpts, _newAggregator common.Address) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "updateStateAggregator", _newAggregator)
}

// UpdateStateAggregator is a paid mutator transaction binding the contract method 0x5ec71e77.
//
// Solidity: function updateStateAggregator(address _newAggregator) returns()
func (_AaveV3Adapter *AaveV3AdapterSession) UpdateStateAggregator(_newAggregator common.Address) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.UpdateStateAggregator(&_AaveV3Adapter.TransactOpts, _newAggregator)
}

// UpdateStateAggregator is a paid mutator transaction binding the contract method 0x5ec71e77.
//
// Solidity: function updateStateAggregator(address _newAggregator) returns()
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) UpdateStateAggregator(_newAggregator common.Address) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.UpdateStateAggregator(&_AaveV3Adapter.TransactOpts, _newAggregator)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address asset, uint256 amount) returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterTransactor) Withdraw(opts *bind.TransactOpts, asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.contract.Transact(opts, "withdraw", asset, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address asset, uint256 amount) returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterSession) Withdraw(asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.Withdraw(&_AaveV3Adapter.TransactOpts, asset, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address asset, uint256 amount) returns(uint256)
func (_AaveV3Adapter *AaveV3AdapterTransactorSession) Withdraw(asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AaveV3Adapter.Contract.Withdraw(&_AaveV3Adapter.TransactOpts, asset, amount)
}

// AaveV3AdapterBorrowedIterator is returned from FilterBorrowed and is used to iterate over the raw logs and unpacked data for Borrowed events raised by the AaveV3Adapter contract.
type AaveV3AdapterBorrowedIterator struct {
	Event *AaveV3AdapterBorrowed // Event containing the contract specifics and raw log

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
func (it *AaveV3AdapterBorrowedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveV3AdapterBorrowed)
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
		it.Event = new(AaveV3AdapterBorrowed)
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
func (it *AaveV3AdapterBorrowedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveV3AdapterBorrowedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveV3AdapterBorrowed represents a Borrowed event raised by the AaveV3Adapter contract.
type AaveV3AdapterBorrowed struct {
	User             common.Address
	Asset            common.Address
	Amount           *big.Int
	InterestRateMode *big.Int
	Timestamp        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterBorrowed is a free log retrieval operation binding the contract event 0xc1cba78646fef030830d099fc25cb498953709c9d47d883848f81fd207174c9f.
//
// Solidity: event Borrowed(address indexed user, address indexed asset, uint256 amount, uint256 interestRateMode, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) FilterBorrowed(opts *bind.FilterOpts, user []common.Address, asset []common.Address) (*AaveV3AdapterBorrowedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.FilterLogs(opts, "Borrowed", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &AaveV3AdapterBorrowedIterator{contract: _AaveV3Adapter.contract, event: "Borrowed", logs: logs, sub: sub}, nil
}

// WatchBorrowed is a free log subscription operation binding the contract event 0xc1cba78646fef030830d099fc25cb498953709c9d47d883848f81fd207174c9f.
//
// Solidity: event Borrowed(address indexed user, address indexed asset, uint256 amount, uint256 interestRateMode, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) WatchBorrowed(opts *bind.WatchOpts, sink chan<- *AaveV3AdapterBorrowed, user []common.Address, asset []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.WatchLogs(opts, "Borrowed", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveV3AdapterBorrowed)
				if err := _AaveV3Adapter.contract.UnpackLog(event, "Borrowed", log); err != nil {
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

// ParseBorrowed is a log parse operation binding the contract event 0xc1cba78646fef030830d099fc25cb498953709c9d47d883848f81fd207174c9f.
//
// Solidity: event Borrowed(address indexed user, address indexed asset, uint256 amount, uint256 interestRateMode, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) ParseBorrowed(log types.Log) (*AaveV3AdapterBorrowed, error) {
	event := new(AaveV3AdapterBorrowed)
	if err := _AaveV3Adapter.contract.UnpackLog(event, "Borrowed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AaveV3AdapterCollateralStatusChangedIterator is returned from FilterCollateralStatusChanged and is used to iterate over the raw logs and unpacked data for CollateralStatusChanged events raised by the AaveV3Adapter contract.
type AaveV3AdapterCollateralStatusChangedIterator struct {
	Event *AaveV3AdapterCollateralStatusChanged // Event containing the contract specifics and raw log

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
func (it *AaveV3AdapterCollateralStatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveV3AdapterCollateralStatusChanged)
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
		it.Event = new(AaveV3AdapterCollateralStatusChanged)
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
func (it *AaveV3AdapterCollateralStatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveV3AdapterCollateralStatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveV3AdapterCollateralStatusChanged represents a CollateralStatusChanged event raised by the AaveV3Adapter contract.
type AaveV3AdapterCollateralStatusChanged struct {
	User            common.Address
	Asset           common.Address
	UseAsCollateral bool
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCollateralStatusChanged is a free log retrieval operation binding the contract event 0x540342942401f36944b5df8d514c51499b07194c3655ab20df5e3dd2f1af6cb4.
//
// Solidity: event CollateralStatusChanged(address indexed user, address indexed asset, bool useAsCollateral)
func (_AaveV3Adapter *AaveV3AdapterFilterer) FilterCollateralStatusChanged(opts *bind.FilterOpts, user []common.Address, asset []common.Address) (*AaveV3AdapterCollateralStatusChangedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.FilterLogs(opts, "CollateralStatusChanged", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &AaveV3AdapterCollateralStatusChangedIterator{contract: _AaveV3Adapter.contract, event: "CollateralStatusChanged", logs: logs, sub: sub}, nil
}

// WatchCollateralStatusChanged is a free log subscription operation binding the contract event 0x540342942401f36944b5df8d514c51499b07194c3655ab20df5e3dd2f1af6cb4.
//
// Solidity: event CollateralStatusChanged(address indexed user, address indexed asset, bool useAsCollateral)
func (_AaveV3Adapter *AaveV3AdapterFilterer) WatchCollateralStatusChanged(opts *bind.WatchOpts, sink chan<- *AaveV3AdapterCollateralStatusChanged, user []common.Address, asset []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.WatchLogs(opts, "CollateralStatusChanged", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveV3AdapterCollateralStatusChanged)
				if err := _AaveV3Adapter.contract.UnpackLog(event, "CollateralStatusChanged", log); err != nil {
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

// ParseCollateralStatusChanged is a log parse operation binding the contract event 0x540342942401f36944b5df8d514c51499b07194c3655ab20df5e3dd2f1af6cb4.
//
// Solidity: event CollateralStatusChanged(address indexed user, address indexed asset, bool useAsCollateral)
func (_AaveV3Adapter *AaveV3AdapterFilterer) ParseCollateralStatusChanged(log types.Log) (*AaveV3AdapterCollateralStatusChanged, error) {
	event := new(AaveV3AdapterCollateralStatusChanged)
	if err := _AaveV3Adapter.contract.UnpackLog(event, "CollateralStatusChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AaveV3AdapterFlashLoanExecutedIterator is returned from FilterFlashLoanExecuted and is used to iterate over the raw logs and unpacked data for FlashLoanExecuted events raised by the AaveV3Adapter contract.
type AaveV3AdapterFlashLoanExecutedIterator struct {
	Event *AaveV3AdapterFlashLoanExecuted // Event containing the contract specifics and raw log

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
func (it *AaveV3AdapterFlashLoanExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveV3AdapterFlashLoanExecuted)
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
		it.Event = new(AaveV3AdapterFlashLoanExecuted)
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
func (it *AaveV3AdapterFlashLoanExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveV3AdapterFlashLoanExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveV3AdapterFlashLoanExecuted represents a FlashLoanExecuted event raised by the AaveV3Adapter contract.
type AaveV3AdapterFlashLoanExecuted struct {
	Initiator common.Address
	Asset     common.Address
	Amount    *big.Int
	Premium   *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFlashLoanExecuted is a free log retrieval operation binding the contract event 0x7496e4be6af6dfb89358a476c6b4e25416b1223f4d0a29999fefbaf4adb8d515.
//
// Solidity: event FlashLoanExecuted(address indexed initiator, address indexed asset, uint256 amount, uint256 premium, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) FilterFlashLoanExecuted(opts *bind.FilterOpts, initiator []common.Address, asset []common.Address) (*AaveV3AdapterFlashLoanExecutedIterator, error) {

	var initiatorRule []interface{}
	for _, initiatorItem := range initiator {
		initiatorRule = append(initiatorRule, initiatorItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.FilterLogs(opts, "FlashLoanExecuted", initiatorRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &AaveV3AdapterFlashLoanExecutedIterator{contract: _AaveV3Adapter.contract, event: "FlashLoanExecuted", logs: logs, sub: sub}, nil
}

// WatchFlashLoanExecuted is a free log subscription operation binding the contract event 0x7496e4be6af6dfb89358a476c6b4e25416b1223f4d0a29999fefbaf4adb8d515.
//
// Solidity: event FlashLoanExecuted(address indexed initiator, address indexed asset, uint256 amount, uint256 premium, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) WatchFlashLoanExecuted(opts *bind.WatchOpts, sink chan<- *AaveV3AdapterFlashLoanExecuted, initiator []common.Address, asset []common.Address) (event.Subscription, error) {

	var initiatorRule []interface{}
	for _, initiatorItem := range initiator {
		initiatorRule = append(initiatorRule, initiatorItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.WatchLogs(opts, "FlashLoanExecuted", initiatorRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveV3AdapterFlashLoanExecuted)
				if err := _AaveV3Adapter.contract.UnpackLog(event, "FlashLoanExecuted", log); err != nil {
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

// ParseFlashLoanExecuted is a log parse operation binding the contract event 0x7496e4be6af6dfb89358a476c6b4e25416b1223f4d0a29999fefbaf4adb8d515.
//
// Solidity: event FlashLoanExecuted(address indexed initiator, address indexed asset, uint256 amount, uint256 premium, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) ParseFlashLoanExecuted(log types.Log) (*AaveV3AdapterFlashLoanExecuted, error) {
	event := new(AaveV3AdapterFlashLoanExecuted)
	if err := _AaveV3Adapter.contract.UnpackLog(event, "FlashLoanExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AaveV3AdapterInterestRateModeSwappedIterator is returned from FilterInterestRateModeSwapped and is used to iterate over the raw logs and unpacked data for InterestRateModeSwapped events raised by the AaveV3Adapter contract.
type AaveV3AdapterInterestRateModeSwappedIterator struct {
	Event *AaveV3AdapterInterestRateModeSwapped // Event containing the contract specifics and raw log

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
func (it *AaveV3AdapterInterestRateModeSwappedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveV3AdapterInterestRateModeSwapped)
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
		it.Event = new(AaveV3AdapterInterestRateModeSwapped)
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
func (it *AaveV3AdapterInterestRateModeSwappedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveV3AdapterInterestRateModeSwappedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveV3AdapterInterestRateModeSwapped represents a InterestRateModeSwapped event raised by the AaveV3Adapter contract.
type AaveV3AdapterInterestRateModeSwapped struct {
	User    common.Address
	Asset   common.Address
	NewMode *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInterestRateModeSwapped is a free log retrieval operation binding the contract event 0x7a1499fe284581882902f8fabb44bfcdcc0b84f606cd58eeff27c5d7e9f38cc5.
//
// Solidity: event InterestRateModeSwapped(address indexed user, address indexed asset, uint256 newMode)
func (_AaveV3Adapter *AaveV3AdapterFilterer) FilterInterestRateModeSwapped(opts *bind.FilterOpts, user []common.Address, asset []common.Address) (*AaveV3AdapterInterestRateModeSwappedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.FilterLogs(opts, "InterestRateModeSwapped", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &AaveV3AdapterInterestRateModeSwappedIterator{contract: _AaveV3Adapter.contract, event: "InterestRateModeSwapped", logs: logs, sub: sub}, nil
}

// WatchInterestRateModeSwapped is a free log subscription operation binding the contract event 0x7a1499fe284581882902f8fabb44bfcdcc0b84f606cd58eeff27c5d7e9f38cc5.
//
// Solidity: event InterestRateModeSwapped(address indexed user, address indexed asset, uint256 newMode)
func (_AaveV3Adapter *AaveV3AdapterFilterer) WatchInterestRateModeSwapped(opts *bind.WatchOpts, sink chan<- *AaveV3AdapterInterestRateModeSwapped, user []common.Address, asset []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.WatchLogs(opts, "InterestRateModeSwapped", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveV3AdapterInterestRateModeSwapped)
				if err := _AaveV3Adapter.contract.UnpackLog(event, "InterestRateModeSwapped", log); err != nil {
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

// ParseInterestRateModeSwapped is a log parse operation binding the contract event 0x7a1499fe284581882902f8fabb44bfcdcc0b84f606cd58eeff27c5d7e9f38cc5.
//
// Solidity: event InterestRateModeSwapped(address indexed user, address indexed asset, uint256 newMode)
func (_AaveV3Adapter *AaveV3AdapterFilterer) ParseInterestRateModeSwapped(log types.Log) (*AaveV3AdapterInterestRateModeSwapped, error) {
	event := new(AaveV3AdapterInterestRateModeSwapped)
	if err := _AaveV3Adapter.contract.UnpackLog(event, "InterestRateModeSwapped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AaveV3AdapterOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AaveV3Adapter contract.
type AaveV3AdapterOwnershipTransferredIterator struct {
	Event *AaveV3AdapterOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AaveV3AdapterOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveV3AdapterOwnershipTransferred)
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
		it.Event = new(AaveV3AdapterOwnershipTransferred)
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
func (it *AaveV3AdapterOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveV3AdapterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveV3AdapterOwnershipTransferred represents a OwnershipTransferred event raised by the AaveV3Adapter contract.
type AaveV3AdapterOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AaveV3Adapter *AaveV3AdapterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AaveV3AdapterOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AaveV3AdapterOwnershipTransferredIterator{contract: _AaveV3Adapter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AaveV3Adapter *AaveV3AdapterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AaveV3AdapterOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveV3AdapterOwnershipTransferred)
				if err := _AaveV3Adapter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_AaveV3Adapter *AaveV3AdapterFilterer) ParseOwnershipTransferred(log types.Log) (*AaveV3AdapterOwnershipTransferred, error) {
	event := new(AaveV3AdapterOwnershipTransferred)
	if err := _AaveV3Adapter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AaveV3AdapterRepaidIterator is returned from FilterRepaid and is used to iterate over the raw logs and unpacked data for Repaid events raised by the AaveV3Adapter contract.
type AaveV3AdapterRepaidIterator struct {
	Event *AaveV3AdapterRepaid // Event containing the contract specifics and raw log

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
func (it *AaveV3AdapterRepaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveV3AdapterRepaid)
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
		it.Event = new(AaveV3AdapterRepaid)
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
func (it *AaveV3AdapterRepaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveV3AdapterRepaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveV3AdapterRepaid represents a Repaid event raised by the AaveV3Adapter contract.
type AaveV3AdapterRepaid struct {
	User             common.Address
	Asset            common.Address
	Amount           *big.Int
	InterestRateMode *big.Int
	Timestamp        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRepaid is a free log retrieval operation binding the contract event 0x244fed5b3268b6a17fff986a36f3d87e163ac678d91d552866eb1e74f010f3a6.
//
// Solidity: event Repaid(address indexed user, address indexed asset, uint256 amount, uint256 interestRateMode, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) FilterRepaid(opts *bind.FilterOpts, user []common.Address, asset []common.Address) (*AaveV3AdapterRepaidIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.FilterLogs(opts, "Repaid", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &AaveV3AdapterRepaidIterator{contract: _AaveV3Adapter.contract, event: "Repaid", logs: logs, sub: sub}, nil
}

// WatchRepaid is a free log subscription operation binding the contract event 0x244fed5b3268b6a17fff986a36f3d87e163ac678d91d552866eb1e74f010f3a6.
//
// Solidity: event Repaid(address indexed user, address indexed asset, uint256 amount, uint256 interestRateMode, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) WatchRepaid(opts *bind.WatchOpts, sink chan<- *AaveV3AdapterRepaid, user []common.Address, asset []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.WatchLogs(opts, "Repaid", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveV3AdapterRepaid)
				if err := _AaveV3Adapter.contract.UnpackLog(event, "Repaid", log); err != nil {
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

// ParseRepaid is a log parse operation binding the contract event 0x244fed5b3268b6a17fff986a36f3d87e163ac678d91d552866eb1e74f010f3a6.
//
// Solidity: event Repaid(address indexed user, address indexed asset, uint256 amount, uint256 interestRateMode, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) ParseRepaid(log types.Log) (*AaveV3AdapterRepaid, error) {
	event := new(AaveV3AdapterRepaid)
	if err := _AaveV3Adapter.contract.UnpackLog(event, "Repaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AaveV3AdapterSuppliedIterator is returned from FilterSupplied and is used to iterate over the raw logs and unpacked data for Supplied events raised by the AaveV3Adapter contract.
type AaveV3AdapterSuppliedIterator struct {
	Event *AaveV3AdapterSupplied // Event containing the contract specifics and raw log

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
func (it *AaveV3AdapterSuppliedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveV3AdapterSupplied)
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
		it.Event = new(AaveV3AdapterSupplied)
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
func (it *AaveV3AdapterSuppliedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveV3AdapterSuppliedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveV3AdapterSupplied represents a Supplied event raised by the AaveV3Adapter contract.
type AaveV3AdapterSupplied struct {
	User      common.Address
	Asset     common.Address
	Amount    *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSupplied is a free log retrieval operation binding the contract event 0x2b650c94fbe31b46ae7946bf007a115a23c6e4532d75ad025261819304334598.
//
// Solidity: event Supplied(address indexed user, address indexed asset, uint256 amount, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) FilterSupplied(opts *bind.FilterOpts, user []common.Address, asset []common.Address) (*AaveV3AdapterSuppliedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.FilterLogs(opts, "Supplied", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &AaveV3AdapterSuppliedIterator{contract: _AaveV3Adapter.contract, event: "Supplied", logs: logs, sub: sub}, nil
}

// WatchSupplied is a free log subscription operation binding the contract event 0x2b650c94fbe31b46ae7946bf007a115a23c6e4532d75ad025261819304334598.
//
// Solidity: event Supplied(address indexed user, address indexed asset, uint256 amount, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) WatchSupplied(opts *bind.WatchOpts, sink chan<- *AaveV3AdapterSupplied, user []common.Address, asset []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.WatchLogs(opts, "Supplied", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveV3AdapterSupplied)
				if err := _AaveV3Adapter.contract.UnpackLog(event, "Supplied", log); err != nil {
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

// ParseSupplied is a log parse operation binding the contract event 0x2b650c94fbe31b46ae7946bf007a115a23c6e4532d75ad025261819304334598.
//
// Solidity: event Supplied(address indexed user, address indexed asset, uint256 amount, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) ParseSupplied(log types.Log) (*AaveV3AdapterSupplied, error) {
	event := new(AaveV3AdapterSupplied)
	if err := _AaveV3Adapter.contract.UnpackLog(event, "Supplied", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AaveV3AdapterWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the AaveV3Adapter contract.
type AaveV3AdapterWithdrawnIterator struct {
	Event *AaveV3AdapterWithdrawn // Event containing the contract specifics and raw log

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
func (it *AaveV3AdapterWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AaveV3AdapterWithdrawn)
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
		it.Event = new(AaveV3AdapterWithdrawn)
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
func (it *AaveV3AdapterWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AaveV3AdapterWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AaveV3AdapterWithdrawn represents a Withdrawn event raised by the AaveV3Adapter contract.
type AaveV3AdapterWithdrawn struct {
	User      common.Address
	Asset     common.Address
	Amount    *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x91fb9d98b786c57d74c099ccd2beca1739e9f6a81fb49001ca465c4b7591bbe2.
//
// Solidity: event Withdrawn(address indexed user, address indexed asset, uint256 amount, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) FilterWithdrawn(opts *bind.FilterOpts, user []common.Address, asset []common.Address) (*AaveV3AdapterWithdrawnIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.FilterLogs(opts, "Withdrawn", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &AaveV3AdapterWithdrawnIterator{contract: _AaveV3Adapter.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x91fb9d98b786c57d74c099ccd2beca1739e9f6a81fb49001ca465c4b7591bbe2.
//
// Solidity: event Withdrawn(address indexed user, address indexed asset, uint256 amount, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *AaveV3AdapterWithdrawn, user []common.Address, asset []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _AaveV3Adapter.contract.WatchLogs(opts, "Withdrawn", userRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AaveV3AdapterWithdrawn)
				if err := _AaveV3Adapter.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0x91fb9d98b786c57d74c099ccd2beca1739e9f6a81fb49001ca465c4b7591bbe2.
//
// Solidity: event Withdrawn(address indexed user, address indexed asset, uint256 amount, uint256 timestamp)
func (_AaveV3Adapter *AaveV3AdapterFilterer) ParseWithdrawn(log types.Log) (*AaveV3AdapterWithdrawn, error) {
	event := new(AaveV3AdapterWithdrawn)
	if err := _AaveV3Adapter.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
