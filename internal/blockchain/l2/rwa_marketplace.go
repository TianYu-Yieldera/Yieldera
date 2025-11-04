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

// IRWAMarketplaceAuction is an auto generated low-level Go binding around an user-defined struct.
type IRWAMarketplaceAuction struct {
	StartPrice     *big.Int
	ReservePrice   *big.Int
	CurrentBid     *big.Int
	CurrentBidder  common.Address
	BidIncrement   *big.Int
	AuctionEndTime *big.Int
}

// IRWAMarketplaceListing is an auto generated low-level Go binding around an user-defined struct.
type IRWAMarketplaceListing struct {
	ListingId       *big.Int
	AssetId         *big.Int
	Seller          common.Address
	OrderType       uint8
	Status          uint8
	Amount          *big.Int
	Price           *big.Int
	MinPurchase     *big.Int
	Filled          *big.Int
	CreatedAt       *big.Int
	ExpiresAt       *big.Int
	IsPrimaryMarket bool
}

// IRWAMarketplaceTrade is an auto generated low-level Go binding around an user-defined struct.
type IRWAMarketplaceTrade struct {
	TradeId    *big.Int
	ListingId  *big.Int
	AssetId    *big.Int
	Buyer      common.Address
	Seller     common.Address
	Amount     *big.Int
	Price      *big.Int
	TotalValue *big.Int
	Fee        *big.Int
	Timestamp  *big.Int
}

// RWAMarketplaceMetaData contains all meta data concerning the RWAMarketplace contract.
var RWAMarketplaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"compliance\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"valuation\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeCollector_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"winner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"finalPrice\",\"type\":\"uint256\"}],\"name\":\"AuctionFinalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"bidder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"bidAmount\",\"type\":\"uint256\"}],\"name\":\"BidPlaced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"}],\"name\":\"ListingCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumIRWAMarketplace.OrderType\",\"name\":\"orderType\",\"type\":\"uint8\"}],\"name\":\"ListingCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tradeId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalValue\",\"type\":\"uint256\"}],\"name\":\"TradExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_COLLECTOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MARKETPLACE_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_FEE_RATE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"assetFactory\",\"outputs\":[{\"internalType\":\"contractIRWAAsset\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"buyTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"tradeId\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"calculateFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"canCreateListing\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"}],\"name\":\"cancelListing\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"collectedFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"complianceContract\",\"outputs\":[{\"internalType\":\"contractIRWACompliance\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"enumIRWAMarketplace.OrderType\",\"name\":\"orderType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minPurchase\",\"type\":\"uint256\"}],\"name\":\"createListing\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeCollector\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"}],\"name\":\"finalizeAuction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getActiveListings\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"}],\"name\":\"getAuction\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"startPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reservePrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentBid\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"currentBidder\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"bidIncrement\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"auctionEndTime\",\"type\":\"uint256\"}],\"internalType\":\"structIRWAMarketplace.Auction\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"}],\"name\":\"getListing\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"enumIRWAMarketplace.OrderType\",\"name\":\"orderType\",\"type\":\"uint8\"},{\"internalType\":\"enumIRWAMarketplace.OrderStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minPurchase\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"filled\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiresAt\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPrimaryMarket\",\"type\":\"bool\"}],\"internalType\":\"structIRWAMarketplace.Listing\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"getTradeHistory\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"tradeId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIRWAMarketplace.Trade[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bidAmount\",\"type\":\"uint256\"}],\"name\":\"placeBid\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"setFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newFeeRate\",\"type\":\"uint256\"}],\"name\":\"setFeeRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalListings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalTrades\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalVolume\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"valuationContract\",\"outputs\":[{\"internalType\":\"contractIRWAValuation\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"withdrawFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// RWAMarketplaceABI is the input ABI used to generate the binding from.
// Deprecated: Use RWAMarketplaceMetaData.ABI instead.
var RWAMarketplaceABI = RWAMarketplaceMetaData.ABI

// RWAMarketplace is an auto generated Go binding around an Ethereum contract.
type RWAMarketplace struct {
	RWAMarketplaceCaller     // Read-only binding to the contract
	RWAMarketplaceTransactor // Write-only binding to the contract
	RWAMarketplaceFilterer   // Log filterer for contract events
}

// RWAMarketplaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type RWAMarketplaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAMarketplaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RWAMarketplaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAMarketplaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RWAMarketplaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAMarketplaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RWAMarketplaceSession struct {
	Contract     *RWAMarketplace   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RWAMarketplaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RWAMarketplaceCallerSession struct {
	Contract *RWAMarketplaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// RWAMarketplaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RWAMarketplaceTransactorSession struct {
	Contract     *RWAMarketplaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// RWAMarketplaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type RWAMarketplaceRaw struct {
	Contract *RWAMarketplace // Generic contract binding to access the raw methods on
}

// RWAMarketplaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RWAMarketplaceCallerRaw struct {
	Contract *RWAMarketplaceCaller // Generic read-only contract binding to access the raw methods on
}

// RWAMarketplaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RWAMarketplaceTransactorRaw struct {
	Contract *RWAMarketplaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRWAMarketplace creates a new instance of RWAMarketplace, bound to a specific deployed contract.
func NewRWAMarketplace(address common.Address, backend bind.ContractBackend) (*RWAMarketplace, error) {
	contract, err := bindRWAMarketplace(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RWAMarketplace{RWAMarketplaceCaller: RWAMarketplaceCaller{contract: contract}, RWAMarketplaceTransactor: RWAMarketplaceTransactor{contract: contract}, RWAMarketplaceFilterer: RWAMarketplaceFilterer{contract: contract}}, nil
}

// NewRWAMarketplaceCaller creates a new read-only instance of RWAMarketplace, bound to a specific deployed contract.
func NewRWAMarketplaceCaller(address common.Address, caller bind.ContractCaller) (*RWAMarketplaceCaller, error) {
	contract, err := bindRWAMarketplace(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RWAMarketplaceCaller{contract: contract}, nil
}

// NewRWAMarketplaceTransactor creates a new write-only instance of RWAMarketplace, bound to a specific deployed contract.
func NewRWAMarketplaceTransactor(address common.Address, transactor bind.ContractTransactor) (*RWAMarketplaceTransactor, error) {
	contract, err := bindRWAMarketplace(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RWAMarketplaceTransactor{contract: contract}, nil
}

// NewRWAMarketplaceFilterer creates a new log filterer instance of RWAMarketplace, bound to a specific deployed contract.
func NewRWAMarketplaceFilterer(address common.Address, filterer bind.ContractFilterer) (*RWAMarketplaceFilterer, error) {
	contract, err := bindRWAMarketplace(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RWAMarketplaceFilterer{contract: contract}, nil
}

// bindRWAMarketplace binds a generic wrapper to an already deployed contract.
func bindRWAMarketplace(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RWAMarketplaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RWAMarketplace *RWAMarketplaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RWAMarketplace.Contract.RWAMarketplaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RWAMarketplace *RWAMarketplaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.RWAMarketplaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RWAMarketplace *RWAMarketplaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.RWAMarketplaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RWAMarketplace *RWAMarketplaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RWAMarketplace.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RWAMarketplace *RWAMarketplaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RWAMarketplace *RWAMarketplaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAMarketplace *RWAMarketplaceCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAMarketplace *RWAMarketplaceSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RWAMarketplace.Contract.DEFAULTADMINROLE(&_RWAMarketplace.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAMarketplace *RWAMarketplaceCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RWAMarketplace.Contract.DEFAULTADMINROLE(&_RWAMarketplace.CallOpts)
}

// FEECOLLECTORROLE is a free data retrieval call binding the contract method 0x62a2a47c.
//
// Solidity: function FEE_COLLECTOR_ROLE() view returns(bytes32)
func (_RWAMarketplace *RWAMarketplaceCaller) FEECOLLECTORROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "FEE_COLLECTOR_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FEECOLLECTORROLE is a free data retrieval call binding the contract method 0x62a2a47c.
//
// Solidity: function FEE_COLLECTOR_ROLE() view returns(bytes32)
func (_RWAMarketplace *RWAMarketplaceSession) FEECOLLECTORROLE() ([32]byte, error) {
	return _RWAMarketplace.Contract.FEECOLLECTORROLE(&_RWAMarketplace.CallOpts)
}

// FEECOLLECTORROLE is a free data retrieval call binding the contract method 0x62a2a47c.
//
// Solidity: function FEE_COLLECTOR_ROLE() view returns(bytes32)
func (_RWAMarketplace *RWAMarketplaceCallerSession) FEECOLLECTORROLE() ([32]byte, error) {
	return _RWAMarketplace.Contract.FEECOLLECTORROLE(&_RWAMarketplace.CallOpts)
}

// FEEPRECISION is a free data retrieval call binding the contract method 0xe63a391f.
//
// Solidity: function FEE_PRECISION() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCaller) FEEPRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "FEE_PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FEEPRECISION is a free data retrieval call binding the contract method 0xe63a391f.
//
// Solidity: function FEE_PRECISION() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceSession) FEEPRECISION() (*big.Int, error) {
	return _RWAMarketplace.Contract.FEEPRECISION(&_RWAMarketplace.CallOpts)
}

// FEEPRECISION is a free data retrieval call binding the contract method 0xe63a391f.
//
// Solidity: function FEE_PRECISION() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCallerSession) FEEPRECISION() (*big.Int, error) {
	return _RWAMarketplace.Contract.FEEPRECISION(&_RWAMarketplace.CallOpts)
}

// MARKETPLACEADMINROLE is a free data retrieval call binding the contract method 0xb2a4eea0.
//
// Solidity: function MARKETPLACE_ADMIN_ROLE() view returns(bytes32)
func (_RWAMarketplace *RWAMarketplaceCaller) MARKETPLACEADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "MARKETPLACE_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MARKETPLACEADMINROLE is a free data retrieval call binding the contract method 0xb2a4eea0.
//
// Solidity: function MARKETPLACE_ADMIN_ROLE() view returns(bytes32)
func (_RWAMarketplace *RWAMarketplaceSession) MARKETPLACEADMINROLE() ([32]byte, error) {
	return _RWAMarketplace.Contract.MARKETPLACEADMINROLE(&_RWAMarketplace.CallOpts)
}

// MARKETPLACEADMINROLE is a free data retrieval call binding the contract method 0xb2a4eea0.
//
// Solidity: function MARKETPLACE_ADMIN_ROLE() view returns(bytes32)
func (_RWAMarketplace *RWAMarketplaceCallerSession) MARKETPLACEADMINROLE() ([32]byte, error) {
	return _RWAMarketplace.Contract.MARKETPLACEADMINROLE(&_RWAMarketplace.CallOpts)
}

// MAXFEERATE is a free data retrieval call binding the contract method 0x92f6576e.
//
// Solidity: function MAX_FEE_RATE() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCaller) MAXFEERATE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "MAX_FEE_RATE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXFEERATE is a free data retrieval call binding the contract method 0x92f6576e.
//
// Solidity: function MAX_FEE_RATE() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceSession) MAXFEERATE() (*big.Int, error) {
	return _RWAMarketplace.Contract.MAXFEERATE(&_RWAMarketplace.CallOpts)
}

// MAXFEERATE is a free data retrieval call binding the contract method 0x92f6576e.
//
// Solidity: function MAX_FEE_RATE() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCallerSession) MAXFEERATE() (*big.Int, error) {
	return _RWAMarketplace.Contract.MAXFEERATE(&_RWAMarketplace.CallOpts)
}

// AssetFactory is a free data retrieval call binding the contract method 0xc3acb4d1.
//
// Solidity: function assetFactory() view returns(address)
func (_RWAMarketplace *RWAMarketplaceCaller) AssetFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "assetFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AssetFactory is a free data retrieval call binding the contract method 0xc3acb4d1.
//
// Solidity: function assetFactory() view returns(address)
func (_RWAMarketplace *RWAMarketplaceSession) AssetFactory() (common.Address, error) {
	return _RWAMarketplace.Contract.AssetFactory(&_RWAMarketplace.CallOpts)
}

// AssetFactory is a free data retrieval call binding the contract method 0xc3acb4d1.
//
// Solidity: function assetFactory() view returns(address)
func (_RWAMarketplace *RWAMarketplaceCallerSession) AssetFactory() (common.Address, error) {
	return _RWAMarketplace.Contract.AssetFactory(&_RWAMarketplace.CallOpts)
}

// CalculateFee is a free data retrieval call binding the contract method 0x34e73122.
//
// Solidity: function calculateFee(uint256 amount, uint256 price) view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCaller) CalculateFee(opts *bind.CallOpts, amount *big.Int, price *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "calculateFee", amount, price)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateFee is a free data retrieval call binding the contract method 0x34e73122.
//
// Solidity: function calculateFee(uint256 amount, uint256 price) view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceSession) CalculateFee(amount *big.Int, price *big.Int) (*big.Int, error) {
	return _RWAMarketplace.Contract.CalculateFee(&_RWAMarketplace.CallOpts, amount, price)
}

// CalculateFee is a free data retrieval call binding the contract method 0x34e73122.
//
// Solidity: function calculateFee(uint256 amount, uint256 price) view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCallerSession) CalculateFee(amount *big.Int, price *big.Int) (*big.Int, error) {
	return _RWAMarketplace.Contract.CalculateFee(&_RWAMarketplace.CallOpts, amount, price)
}

// CanCreateListing is a free data retrieval call binding the contract method 0xde3b5ea9.
//
// Solidity: function canCreateListing(address seller, uint256 assetId, uint256 amount) view returns(bool)
func (_RWAMarketplace *RWAMarketplaceCaller) CanCreateListing(opts *bind.CallOpts, seller common.Address, assetId *big.Int, amount *big.Int) (bool, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "canCreateListing", seller, assetId, amount)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanCreateListing is a free data retrieval call binding the contract method 0xde3b5ea9.
//
// Solidity: function canCreateListing(address seller, uint256 assetId, uint256 amount) view returns(bool)
func (_RWAMarketplace *RWAMarketplaceSession) CanCreateListing(seller common.Address, assetId *big.Int, amount *big.Int) (bool, error) {
	return _RWAMarketplace.Contract.CanCreateListing(&_RWAMarketplace.CallOpts, seller, assetId, amount)
}

// CanCreateListing is a free data retrieval call binding the contract method 0xde3b5ea9.
//
// Solidity: function canCreateListing(address seller, uint256 assetId, uint256 amount) view returns(bool)
func (_RWAMarketplace *RWAMarketplaceCallerSession) CanCreateListing(seller common.Address, assetId *big.Int, amount *big.Int) (bool, error) {
	return _RWAMarketplace.Contract.CanCreateListing(&_RWAMarketplace.CallOpts, seller, assetId, amount)
}

// CollectedFees is a free data retrieval call binding the contract method 0x1cead9a7.
//
// Solidity: function collectedFees(address ) view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCaller) CollectedFees(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "collectedFees", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CollectedFees is a free data retrieval call binding the contract method 0x1cead9a7.
//
// Solidity: function collectedFees(address ) view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceSession) CollectedFees(arg0 common.Address) (*big.Int, error) {
	return _RWAMarketplace.Contract.CollectedFees(&_RWAMarketplace.CallOpts, arg0)
}

// CollectedFees is a free data retrieval call binding the contract method 0x1cead9a7.
//
// Solidity: function collectedFees(address ) view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCallerSession) CollectedFees(arg0 common.Address) (*big.Int, error) {
	return _RWAMarketplace.Contract.CollectedFees(&_RWAMarketplace.CallOpts, arg0)
}

// ComplianceContract is a free data retrieval call binding the contract method 0xb2a2a4e2.
//
// Solidity: function complianceContract() view returns(address)
func (_RWAMarketplace *RWAMarketplaceCaller) ComplianceContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "complianceContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ComplianceContract is a free data retrieval call binding the contract method 0xb2a2a4e2.
//
// Solidity: function complianceContract() view returns(address)
func (_RWAMarketplace *RWAMarketplaceSession) ComplianceContract() (common.Address, error) {
	return _RWAMarketplace.Contract.ComplianceContract(&_RWAMarketplace.CallOpts)
}

// ComplianceContract is a free data retrieval call binding the contract method 0xb2a2a4e2.
//
// Solidity: function complianceContract() view returns(address)
func (_RWAMarketplace *RWAMarketplaceCallerSession) ComplianceContract() (common.Address, error) {
	return _RWAMarketplace.Contract.ComplianceContract(&_RWAMarketplace.CallOpts)
}

// FeeCollector is a free data retrieval call binding the contract method 0xc415b95c.
//
// Solidity: function feeCollector() view returns(address)
func (_RWAMarketplace *RWAMarketplaceCaller) FeeCollector(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "feeCollector")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeCollector is a free data retrieval call binding the contract method 0xc415b95c.
//
// Solidity: function feeCollector() view returns(address)
func (_RWAMarketplace *RWAMarketplaceSession) FeeCollector() (common.Address, error) {
	return _RWAMarketplace.Contract.FeeCollector(&_RWAMarketplace.CallOpts)
}

// FeeCollector is a free data retrieval call binding the contract method 0xc415b95c.
//
// Solidity: function feeCollector() view returns(address)
func (_RWAMarketplace *RWAMarketplaceCallerSession) FeeCollector() (common.Address, error) {
	return _RWAMarketplace.Contract.FeeCollector(&_RWAMarketplace.CallOpts)
}

// FeeRate is a free data retrieval call binding the contract method 0x978bbdb9.
//
// Solidity: function feeRate() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCaller) FeeRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "feeRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeRate is a free data retrieval call binding the contract method 0x978bbdb9.
//
// Solidity: function feeRate() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceSession) FeeRate() (*big.Int, error) {
	return _RWAMarketplace.Contract.FeeRate(&_RWAMarketplace.CallOpts)
}

// FeeRate is a free data retrieval call binding the contract method 0x978bbdb9.
//
// Solidity: function feeRate() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCallerSession) FeeRate() (*big.Int, error) {
	return _RWAMarketplace.Contract.FeeRate(&_RWAMarketplace.CallOpts)
}

// GetActiveListings is a free data retrieval call binding the contract method 0xe6987501.
//
// Solidity: function getActiveListings(uint256 assetId) view returns(uint256[])
func (_RWAMarketplace *RWAMarketplaceCaller) GetActiveListings(opts *bind.CallOpts, assetId *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "getActiveListings", assetId)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveListings is a free data retrieval call binding the contract method 0xe6987501.
//
// Solidity: function getActiveListings(uint256 assetId) view returns(uint256[])
func (_RWAMarketplace *RWAMarketplaceSession) GetActiveListings(assetId *big.Int) ([]*big.Int, error) {
	return _RWAMarketplace.Contract.GetActiveListings(&_RWAMarketplace.CallOpts, assetId)
}

// GetActiveListings is a free data retrieval call binding the contract method 0xe6987501.
//
// Solidity: function getActiveListings(uint256 assetId) view returns(uint256[])
func (_RWAMarketplace *RWAMarketplaceCallerSession) GetActiveListings(assetId *big.Int) ([]*big.Int, error) {
	return _RWAMarketplace.Contract.GetActiveListings(&_RWAMarketplace.CallOpts, assetId)
}

// GetAuction is a free data retrieval call binding the contract method 0x78bd7935.
//
// Solidity: function getAuction(uint256 listingId) view returns((uint256,uint256,uint256,address,uint256,uint256))
func (_RWAMarketplace *RWAMarketplaceCaller) GetAuction(opts *bind.CallOpts, listingId *big.Int) (IRWAMarketplaceAuction, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "getAuction", listingId)

	if err != nil {
		return *new(IRWAMarketplaceAuction), err
	}

	out0 := *abi.ConvertType(out[0], new(IRWAMarketplaceAuction)).(*IRWAMarketplaceAuction)

	return out0, err

}

// GetAuction is a free data retrieval call binding the contract method 0x78bd7935.
//
// Solidity: function getAuction(uint256 listingId) view returns((uint256,uint256,uint256,address,uint256,uint256))
func (_RWAMarketplace *RWAMarketplaceSession) GetAuction(listingId *big.Int) (IRWAMarketplaceAuction, error) {
	return _RWAMarketplace.Contract.GetAuction(&_RWAMarketplace.CallOpts, listingId)
}

// GetAuction is a free data retrieval call binding the contract method 0x78bd7935.
//
// Solidity: function getAuction(uint256 listingId) view returns((uint256,uint256,uint256,address,uint256,uint256))
func (_RWAMarketplace *RWAMarketplaceCallerSession) GetAuction(listingId *big.Int) (IRWAMarketplaceAuction, error) {
	return _RWAMarketplace.Contract.GetAuction(&_RWAMarketplace.CallOpts, listingId)
}

// GetListing is a free data retrieval call binding the contract method 0x107a274a.
//
// Solidity: function getListing(uint256 listingId) view returns((uint256,uint256,address,uint8,uint8,uint256,uint256,uint256,uint256,uint256,uint256,bool))
func (_RWAMarketplace *RWAMarketplaceCaller) GetListing(opts *bind.CallOpts, listingId *big.Int) (IRWAMarketplaceListing, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "getListing", listingId)

	if err != nil {
		return *new(IRWAMarketplaceListing), err
	}

	out0 := *abi.ConvertType(out[0], new(IRWAMarketplaceListing)).(*IRWAMarketplaceListing)

	return out0, err

}

// GetListing is a free data retrieval call binding the contract method 0x107a274a.
//
// Solidity: function getListing(uint256 listingId) view returns((uint256,uint256,address,uint8,uint8,uint256,uint256,uint256,uint256,uint256,uint256,bool))
func (_RWAMarketplace *RWAMarketplaceSession) GetListing(listingId *big.Int) (IRWAMarketplaceListing, error) {
	return _RWAMarketplace.Contract.GetListing(&_RWAMarketplace.CallOpts, listingId)
}

// GetListing is a free data retrieval call binding the contract method 0x107a274a.
//
// Solidity: function getListing(uint256 listingId) view returns((uint256,uint256,address,uint8,uint8,uint256,uint256,uint256,uint256,uint256,uint256,bool))
func (_RWAMarketplace *RWAMarketplaceCallerSession) GetListing(listingId *big.Int) (IRWAMarketplaceListing, error) {
	return _RWAMarketplace.Contract.GetListing(&_RWAMarketplace.CallOpts, listingId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAMarketplace *RWAMarketplaceCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAMarketplace *RWAMarketplaceSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RWAMarketplace.Contract.GetRoleAdmin(&_RWAMarketplace.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAMarketplace *RWAMarketplaceCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RWAMarketplace.Contract.GetRoleAdmin(&_RWAMarketplace.CallOpts, role)
}

// GetTradeHistory is a free data retrieval call binding the contract method 0xf52bbdc0.
//
// Solidity: function getTradeHistory(uint256 assetId, uint256 count) view returns((uint256,uint256,uint256,address,address,uint256,uint256,uint256,uint256,uint256)[])
func (_RWAMarketplace *RWAMarketplaceCaller) GetTradeHistory(opts *bind.CallOpts, assetId *big.Int, count *big.Int) ([]IRWAMarketplaceTrade, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "getTradeHistory", assetId, count)

	if err != nil {
		return *new([]IRWAMarketplaceTrade), err
	}

	out0 := *abi.ConvertType(out[0], new([]IRWAMarketplaceTrade)).(*[]IRWAMarketplaceTrade)

	return out0, err

}

// GetTradeHistory is a free data retrieval call binding the contract method 0xf52bbdc0.
//
// Solidity: function getTradeHistory(uint256 assetId, uint256 count) view returns((uint256,uint256,uint256,address,address,uint256,uint256,uint256,uint256,uint256)[])
func (_RWAMarketplace *RWAMarketplaceSession) GetTradeHistory(assetId *big.Int, count *big.Int) ([]IRWAMarketplaceTrade, error) {
	return _RWAMarketplace.Contract.GetTradeHistory(&_RWAMarketplace.CallOpts, assetId, count)
}

// GetTradeHistory is a free data retrieval call binding the contract method 0xf52bbdc0.
//
// Solidity: function getTradeHistory(uint256 assetId, uint256 count) view returns((uint256,uint256,uint256,address,address,uint256,uint256,uint256,uint256,uint256)[])
func (_RWAMarketplace *RWAMarketplaceCallerSession) GetTradeHistory(assetId *big.Int, count *big.Int) ([]IRWAMarketplaceTrade, error) {
	return _RWAMarketplace.Contract.GetTradeHistory(&_RWAMarketplace.CallOpts, assetId, count)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAMarketplace *RWAMarketplaceCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAMarketplace *RWAMarketplaceSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RWAMarketplace.Contract.HasRole(&_RWAMarketplace.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAMarketplace *RWAMarketplaceCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RWAMarketplace.Contract.HasRole(&_RWAMarketplace.CallOpts, role, account)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RWAMarketplace *RWAMarketplaceCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RWAMarketplace *RWAMarketplaceSession) Paused() (bool, error) {
	return _RWAMarketplace.Contract.Paused(&_RWAMarketplace.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RWAMarketplace *RWAMarketplaceCallerSession) Paused() (bool, error) {
	return _RWAMarketplace.Contract.Paused(&_RWAMarketplace.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAMarketplace *RWAMarketplaceCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAMarketplace *RWAMarketplaceSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RWAMarketplace.Contract.SupportsInterface(&_RWAMarketplace.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAMarketplace *RWAMarketplaceCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RWAMarketplace.Contract.SupportsInterface(&_RWAMarketplace.CallOpts, interfaceId)
}

// TotalListings is a free data retrieval call binding the contract method 0xc78b616c.
//
// Solidity: function totalListings() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCaller) TotalListings(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "totalListings")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalListings is a free data retrieval call binding the contract method 0xc78b616c.
//
// Solidity: function totalListings() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceSession) TotalListings() (*big.Int, error) {
	return _RWAMarketplace.Contract.TotalListings(&_RWAMarketplace.CallOpts)
}

// TotalListings is a free data retrieval call binding the contract method 0xc78b616c.
//
// Solidity: function totalListings() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCallerSession) TotalListings() (*big.Int, error) {
	return _RWAMarketplace.Contract.TotalListings(&_RWAMarketplace.CallOpts)
}

// TotalTrades is a free data retrieval call binding the contract method 0xe275c997.
//
// Solidity: function totalTrades() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCaller) TotalTrades(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "totalTrades")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalTrades is a free data retrieval call binding the contract method 0xe275c997.
//
// Solidity: function totalTrades() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceSession) TotalTrades() (*big.Int, error) {
	return _RWAMarketplace.Contract.TotalTrades(&_RWAMarketplace.CallOpts)
}

// TotalTrades is a free data retrieval call binding the contract method 0xe275c997.
//
// Solidity: function totalTrades() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCallerSession) TotalTrades() (*big.Int, error) {
	return _RWAMarketplace.Contract.TotalTrades(&_RWAMarketplace.CallOpts)
}

// TotalVolume is a free data retrieval call binding the contract method 0x5f81a57c.
//
// Solidity: function totalVolume() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCaller) TotalVolume(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "totalVolume")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalVolume is a free data retrieval call binding the contract method 0x5f81a57c.
//
// Solidity: function totalVolume() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceSession) TotalVolume() (*big.Int, error) {
	return _RWAMarketplace.Contract.TotalVolume(&_RWAMarketplace.CallOpts)
}

// TotalVolume is a free data retrieval call binding the contract method 0x5f81a57c.
//
// Solidity: function totalVolume() view returns(uint256)
func (_RWAMarketplace *RWAMarketplaceCallerSession) TotalVolume() (*big.Int, error) {
	return _RWAMarketplace.Contract.TotalVolume(&_RWAMarketplace.CallOpts)
}

// ValuationContract is a free data retrieval call binding the contract method 0x274dc7ca.
//
// Solidity: function valuationContract() view returns(address)
func (_RWAMarketplace *RWAMarketplaceCaller) ValuationContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RWAMarketplace.contract.Call(opts, &out, "valuationContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ValuationContract is a free data retrieval call binding the contract method 0x274dc7ca.
//
// Solidity: function valuationContract() view returns(address)
func (_RWAMarketplace *RWAMarketplaceSession) ValuationContract() (common.Address, error) {
	return _RWAMarketplace.Contract.ValuationContract(&_RWAMarketplace.CallOpts)
}

// ValuationContract is a free data retrieval call binding the contract method 0x274dc7ca.
//
// Solidity: function valuationContract() view returns(address)
func (_RWAMarketplace *RWAMarketplaceCallerSession) ValuationContract() (common.Address, error) {
	return _RWAMarketplace.Contract.ValuationContract(&_RWAMarketplace.CallOpts)
}

// BuyTokens is a paid mutator transaction binding the contract method 0x7975ce28.
//
// Solidity: function buyTokens(uint256 listingId, uint256 amount) payable returns(uint256 tradeId)
func (_RWAMarketplace *RWAMarketplaceTransactor) BuyTokens(opts *bind.TransactOpts, listingId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "buyTokens", listingId, amount)
}

// BuyTokens is a paid mutator transaction binding the contract method 0x7975ce28.
//
// Solidity: function buyTokens(uint256 listingId, uint256 amount) payable returns(uint256 tradeId)
func (_RWAMarketplace *RWAMarketplaceSession) BuyTokens(listingId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.BuyTokens(&_RWAMarketplace.TransactOpts, listingId, amount)
}

// BuyTokens is a paid mutator transaction binding the contract method 0x7975ce28.
//
// Solidity: function buyTokens(uint256 listingId, uint256 amount) payable returns(uint256 tradeId)
func (_RWAMarketplace *RWAMarketplaceTransactorSession) BuyTokens(listingId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.BuyTokens(&_RWAMarketplace.TransactOpts, listingId, amount)
}

// CancelListing is a paid mutator transaction binding the contract method 0x305a67a8.
//
// Solidity: function cancelListing(uint256 listingId) returns()
func (_RWAMarketplace *RWAMarketplaceTransactor) CancelListing(opts *bind.TransactOpts, listingId *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "cancelListing", listingId)
}

// CancelListing is a paid mutator transaction binding the contract method 0x305a67a8.
//
// Solidity: function cancelListing(uint256 listingId) returns()
func (_RWAMarketplace *RWAMarketplaceSession) CancelListing(listingId *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.CancelListing(&_RWAMarketplace.TransactOpts, listingId)
}

// CancelListing is a paid mutator transaction binding the contract method 0x305a67a8.
//
// Solidity: function cancelListing(uint256 listingId) returns()
func (_RWAMarketplace *RWAMarketplaceTransactorSession) CancelListing(listingId *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.CancelListing(&_RWAMarketplace.TransactOpts, listingId)
}

// CreateListing is a paid mutator transaction binding the contract method 0x1a2ddc77.
//
// Solidity: function createListing(uint256 assetId, uint256 amount, uint256 price, uint8 orderType, uint256 duration, uint256 minPurchase) returns(uint256 listingId)
func (_RWAMarketplace *RWAMarketplaceTransactor) CreateListing(opts *bind.TransactOpts, assetId *big.Int, amount *big.Int, price *big.Int, orderType uint8, duration *big.Int, minPurchase *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "createListing", assetId, amount, price, orderType, duration, minPurchase)
}

// CreateListing is a paid mutator transaction binding the contract method 0x1a2ddc77.
//
// Solidity: function createListing(uint256 assetId, uint256 amount, uint256 price, uint8 orderType, uint256 duration, uint256 minPurchase) returns(uint256 listingId)
func (_RWAMarketplace *RWAMarketplaceSession) CreateListing(assetId *big.Int, amount *big.Int, price *big.Int, orderType uint8, duration *big.Int, minPurchase *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.CreateListing(&_RWAMarketplace.TransactOpts, assetId, amount, price, orderType, duration, minPurchase)
}

// CreateListing is a paid mutator transaction binding the contract method 0x1a2ddc77.
//
// Solidity: function createListing(uint256 assetId, uint256 amount, uint256 price, uint8 orderType, uint256 duration, uint256 minPurchase) returns(uint256 listingId)
func (_RWAMarketplace *RWAMarketplaceTransactorSession) CreateListing(assetId *big.Int, amount *big.Int, price *big.Int, orderType uint8, duration *big.Int, minPurchase *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.CreateListing(&_RWAMarketplace.TransactOpts, assetId, amount, price, orderType, duration, minPurchase)
}

// FinalizeAuction is a paid mutator transaction binding the contract method 0xe8083863.
//
// Solidity: function finalizeAuction(uint256 listingId) returns()
func (_RWAMarketplace *RWAMarketplaceTransactor) FinalizeAuction(opts *bind.TransactOpts, listingId *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "finalizeAuction", listingId)
}

// FinalizeAuction is a paid mutator transaction binding the contract method 0xe8083863.
//
// Solidity: function finalizeAuction(uint256 listingId) returns()
func (_RWAMarketplace *RWAMarketplaceSession) FinalizeAuction(listingId *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.FinalizeAuction(&_RWAMarketplace.TransactOpts, listingId)
}

// FinalizeAuction is a paid mutator transaction binding the contract method 0xe8083863.
//
// Solidity: function finalizeAuction(uint256 listingId) returns()
func (_RWAMarketplace *RWAMarketplaceTransactorSession) FinalizeAuction(listingId *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.FinalizeAuction(&_RWAMarketplace.TransactOpts, listingId)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAMarketplace *RWAMarketplaceTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAMarketplace *RWAMarketplaceSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.GrantRole(&_RWAMarketplace.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAMarketplace *RWAMarketplaceTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.GrantRole(&_RWAMarketplace.TransactOpts, role, account)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RWAMarketplace *RWAMarketplaceTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RWAMarketplace *RWAMarketplaceSession) Pause() (*types.Transaction, error) {
	return _RWAMarketplace.Contract.Pause(&_RWAMarketplace.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RWAMarketplace *RWAMarketplaceTransactorSession) Pause() (*types.Transaction, error) {
	return _RWAMarketplace.Contract.Pause(&_RWAMarketplace.TransactOpts)
}

// PlaceBid is a paid mutator transaction binding the contract method 0x57c90de5.
//
// Solidity: function placeBid(uint256 listingId, uint256 bidAmount) payable returns()
func (_RWAMarketplace *RWAMarketplaceTransactor) PlaceBid(opts *bind.TransactOpts, listingId *big.Int, bidAmount *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "placeBid", listingId, bidAmount)
}

// PlaceBid is a paid mutator transaction binding the contract method 0x57c90de5.
//
// Solidity: function placeBid(uint256 listingId, uint256 bidAmount) payable returns()
func (_RWAMarketplace *RWAMarketplaceSession) PlaceBid(listingId *big.Int, bidAmount *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.PlaceBid(&_RWAMarketplace.TransactOpts, listingId, bidAmount)
}

// PlaceBid is a paid mutator transaction binding the contract method 0x57c90de5.
//
// Solidity: function placeBid(uint256 listingId, uint256 bidAmount) payable returns()
func (_RWAMarketplace *RWAMarketplaceTransactorSession) PlaceBid(listingId *big.Int, bidAmount *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.PlaceBid(&_RWAMarketplace.TransactOpts, listingId, bidAmount)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAMarketplace *RWAMarketplaceTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAMarketplace *RWAMarketplaceSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.RenounceRole(&_RWAMarketplace.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAMarketplace *RWAMarketplaceTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.RenounceRole(&_RWAMarketplace.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAMarketplace *RWAMarketplaceTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAMarketplace *RWAMarketplaceSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.RevokeRole(&_RWAMarketplace.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAMarketplace *RWAMarketplaceTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.RevokeRole(&_RWAMarketplace.TransactOpts, role, account)
}

// SetFeeCollector is a paid mutator transaction binding the contract method 0xa42dce80.
//
// Solidity: function setFeeCollector(address newFeeCollector) returns()
func (_RWAMarketplace *RWAMarketplaceTransactor) SetFeeCollector(opts *bind.TransactOpts, newFeeCollector common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "setFeeCollector", newFeeCollector)
}

// SetFeeCollector is a paid mutator transaction binding the contract method 0xa42dce80.
//
// Solidity: function setFeeCollector(address newFeeCollector) returns()
func (_RWAMarketplace *RWAMarketplaceSession) SetFeeCollector(newFeeCollector common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.SetFeeCollector(&_RWAMarketplace.TransactOpts, newFeeCollector)
}

// SetFeeCollector is a paid mutator transaction binding the contract method 0xa42dce80.
//
// Solidity: function setFeeCollector(address newFeeCollector) returns()
func (_RWAMarketplace *RWAMarketplaceTransactorSession) SetFeeCollector(newFeeCollector common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.SetFeeCollector(&_RWAMarketplace.TransactOpts, newFeeCollector)
}

// SetFeeRate is a paid mutator transaction binding the contract method 0x45596e2e.
//
// Solidity: function setFeeRate(uint256 newFeeRate) returns()
func (_RWAMarketplace *RWAMarketplaceTransactor) SetFeeRate(opts *bind.TransactOpts, newFeeRate *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "setFeeRate", newFeeRate)
}

// SetFeeRate is a paid mutator transaction binding the contract method 0x45596e2e.
//
// Solidity: function setFeeRate(uint256 newFeeRate) returns()
func (_RWAMarketplace *RWAMarketplaceSession) SetFeeRate(newFeeRate *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.SetFeeRate(&_RWAMarketplace.TransactOpts, newFeeRate)
}

// SetFeeRate is a paid mutator transaction binding the contract method 0x45596e2e.
//
// Solidity: function setFeeRate(uint256 newFeeRate) returns()
func (_RWAMarketplace *RWAMarketplaceTransactorSession) SetFeeRate(newFeeRate *big.Int) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.SetFeeRate(&_RWAMarketplace.TransactOpts, newFeeRate)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RWAMarketplace *RWAMarketplaceTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RWAMarketplace *RWAMarketplaceSession) Unpause() (*types.Transaction, error) {
	return _RWAMarketplace.Contract.Unpause(&_RWAMarketplace.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RWAMarketplace *RWAMarketplaceTransactorSession) Unpause() (*types.Transaction, error) {
	return _RWAMarketplace.Contract.Unpause(&_RWAMarketplace.TransactOpts)
}

// WithdrawFees is a paid mutator transaction binding the contract method 0x164e68de.
//
// Solidity: function withdrawFees(address token) returns()
func (_RWAMarketplace *RWAMarketplaceTransactor) WithdrawFees(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.contract.Transact(opts, "withdrawFees", token)
}

// WithdrawFees is a paid mutator transaction binding the contract method 0x164e68de.
//
// Solidity: function withdrawFees(address token) returns()
func (_RWAMarketplace *RWAMarketplaceSession) WithdrawFees(token common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.WithdrawFees(&_RWAMarketplace.TransactOpts, token)
}

// WithdrawFees is a paid mutator transaction binding the contract method 0x164e68de.
//
// Solidity: function withdrawFees(address token) returns()
func (_RWAMarketplace *RWAMarketplaceTransactorSession) WithdrawFees(token common.Address) (*types.Transaction, error) {
	return _RWAMarketplace.Contract.WithdrawFees(&_RWAMarketplace.TransactOpts, token)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_RWAMarketplace *RWAMarketplaceTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAMarketplace.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_RWAMarketplace *RWAMarketplaceSession) Receive() (*types.Transaction, error) {
	return _RWAMarketplace.Contract.Receive(&_RWAMarketplace.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_RWAMarketplace *RWAMarketplaceTransactorSession) Receive() (*types.Transaction, error) {
	return _RWAMarketplace.Contract.Receive(&_RWAMarketplace.TransactOpts)
}

// RWAMarketplaceAuctionFinalizedIterator is returned from FilterAuctionFinalized and is used to iterate over the raw logs and unpacked data for AuctionFinalized events raised by the RWAMarketplace contract.
type RWAMarketplaceAuctionFinalizedIterator struct {
	Event *RWAMarketplaceAuctionFinalized // Event containing the contract specifics and raw log

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
func (it *RWAMarketplaceAuctionFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAMarketplaceAuctionFinalized)
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
		it.Event = new(RWAMarketplaceAuctionFinalized)
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
func (it *RWAMarketplaceAuctionFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAMarketplaceAuctionFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAMarketplaceAuctionFinalized represents a AuctionFinalized event raised by the RWAMarketplace contract.
type RWAMarketplaceAuctionFinalized struct {
	ListingId  *big.Int
	Winner     common.Address
	FinalPrice *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAuctionFinalized is a free log retrieval operation binding the contract event 0x4d9113a1377d665eaa1f9168a9c9080f2e488cb820b10149de3d6d2e0f2780c7.
//
// Solidity: event AuctionFinalized(uint256 indexed listingId, address indexed winner, uint256 finalPrice)
func (_RWAMarketplace *RWAMarketplaceFilterer) FilterAuctionFinalized(opts *bind.FilterOpts, listingId []*big.Int, winner []common.Address) (*RWAMarketplaceAuctionFinalizedIterator, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var winnerRule []interface{}
	for _, winnerItem := range winner {
		winnerRule = append(winnerRule, winnerItem)
	}

	logs, sub, err := _RWAMarketplace.contract.FilterLogs(opts, "AuctionFinalized", listingIdRule, winnerRule)
	if err != nil {
		return nil, err
	}
	return &RWAMarketplaceAuctionFinalizedIterator{contract: _RWAMarketplace.contract, event: "AuctionFinalized", logs: logs, sub: sub}, nil
}

// WatchAuctionFinalized is a free log subscription operation binding the contract event 0x4d9113a1377d665eaa1f9168a9c9080f2e488cb820b10149de3d6d2e0f2780c7.
//
// Solidity: event AuctionFinalized(uint256 indexed listingId, address indexed winner, uint256 finalPrice)
func (_RWAMarketplace *RWAMarketplaceFilterer) WatchAuctionFinalized(opts *bind.WatchOpts, sink chan<- *RWAMarketplaceAuctionFinalized, listingId []*big.Int, winner []common.Address) (event.Subscription, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var winnerRule []interface{}
	for _, winnerItem := range winner {
		winnerRule = append(winnerRule, winnerItem)
	}

	logs, sub, err := _RWAMarketplace.contract.WatchLogs(opts, "AuctionFinalized", listingIdRule, winnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAMarketplaceAuctionFinalized)
				if err := _RWAMarketplace.contract.UnpackLog(event, "AuctionFinalized", log); err != nil {
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

// ParseAuctionFinalized is a log parse operation binding the contract event 0x4d9113a1377d665eaa1f9168a9c9080f2e488cb820b10149de3d6d2e0f2780c7.
//
// Solidity: event AuctionFinalized(uint256 indexed listingId, address indexed winner, uint256 finalPrice)
func (_RWAMarketplace *RWAMarketplaceFilterer) ParseAuctionFinalized(log types.Log) (*RWAMarketplaceAuctionFinalized, error) {
	event := new(RWAMarketplaceAuctionFinalized)
	if err := _RWAMarketplace.contract.UnpackLog(event, "AuctionFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAMarketplaceBidPlacedIterator is returned from FilterBidPlaced and is used to iterate over the raw logs and unpacked data for BidPlaced events raised by the RWAMarketplace contract.
type RWAMarketplaceBidPlacedIterator struct {
	Event *RWAMarketplaceBidPlaced // Event containing the contract specifics and raw log

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
func (it *RWAMarketplaceBidPlacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAMarketplaceBidPlaced)
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
		it.Event = new(RWAMarketplaceBidPlaced)
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
func (it *RWAMarketplaceBidPlacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAMarketplaceBidPlacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAMarketplaceBidPlaced represents a BidPlaced event raised by the RWAMarketplace contract.
type RWAMarketplaceBidPlaced struct {
	ListingId *big.Int
	Bidder    common.Address
	BidAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBidPlaced is a free log retrieval operation binding the contract event 0x0e54eff26401bf69b81b26f60bd85ef47f5d85275c1d268d84f68d6897431c47.
//
// Solidity: event BidPlaced(uint256 indexed listingId, address indexed bidder, uint256 bidAmount)
func (_RWAMarketplace *RWAMarketplaceFilterer) FilterBidPlaced(opts *bind.FilterOpts, listingId []*big.Int, bidder []common.Address) (*RWAMarketplaceBidPlacedIterator, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}

	logs, sub, err := _RWAMarketplace.contract.FilterLogs(opts, "BidPlaced", listingIdRule, bidderRule)
	if err != nil {
		return nil, err
	}
	return &RWAMarketplaceBidPlacedIterator{contract: _RWAMarketplace.contract, event: "BidPlaced", logs: logs, sub: sub}, nil
}

// WatchBidPlaced is a free log subscription operation binding the contract event 0x0e54eff26401bf69b81b26f60bd85ef47f5d85275c1d268d84f68d6897431c47.
//
// Solidity: event BidPlaced(uint256 indexed listingId, address indexed bidder, uint256 bidAmount)
func (_RWAMarketplace *RWAMarketplaceFilterer) WatchBidPlaced(opts *bind.WatchOpts, sink chan<- *RWAMarketplaceBidPlaced, listingId []*big.Int, bidder []common.Address) (event.Subscription, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}

	logs, sub, err := _RWAMarketplace.contract.WatchLogs(opts, "BidPlaced", listingIdRule, bidderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAMarketplaceBidPlaced)
				if err := _RWAMarketplace.contract.UnpackLog(event, "BidPlaced", log); err != nil {
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

// ParseBidPlaced is a log parse operation binding the contract event 0x0e54eff26401bf69b81b26f60bd85ef47f5d85275c1d268d84f68d6897431c47.
//
// Solidity: event BidPlaced(uint256 indexed listingId, address indexed bidder, uint256 bidAmount)
func (_RWAMarketplace *RWAMarketplaceFilterer) ParseBidPlaced(log types.Log) (*RWAMarketplaceBidPlaced, error) {
	event := new(RWAMarketplaceBidPlaced)
	if err := _RWAMarketplace.contract.UnpackLog(event, "BidPlaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAMarketplaceListingCancelledIterator is returned from FilterListingCancelled and is used to iterate over the raw logs and unpacked data for ListingCancelled events raised by the RWAMarketplace contract.
type RWAMarketplaceListingCancelledIterator struct {
	Event *RWAMarketplaceListingCancelled // Event containing the contract specifics and raw log

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
func (it *RWAMarketplaceListingCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAMarketplaceListingCancelled)
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
		it.Event = new(RWAMarketplaceListingCancelled)
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
func (it *RWAMarketplaceListingCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAMarketplaceListingCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAMarketplaceListingCancelled represents a ListingCancelled event raised by the RWAMarketplace contract.
type RWAMarketplaceListingCancelled struct {
	ListingId *big.Int
	Seller    common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterListingCancelled is a free log retrieval operation binding the contract event 0x8e25282255ab31897df2b0456bb993ac7f84d376861aefd84901d2d63a7428a2.
//
// Solidity: event ListingCancelled(uint256 indexed listingId, address indexed seller)
func (_RWAMarketplace *RWAMarketplaceFilterer) FilterListingCancelled(opts *bind.FilterOpts, listingId []*big.Int, seller []common.Address) (*RWAMarketplaceListingCancelledIterator, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _RWAMarketplace.contract.FilterLogs(opts, "ListingCancelled", listingIdRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return &RWAMarketplaceListingCancelledIterator{contract: _RWAMarketplace.contract, event: "ListingCancelled", logs: logs, sub: sub}, nil
}

// WatchListingCancelled is a free log subscription operation binding the contract event 0x8e25282255ab31897df2b0456bb993ac7f84d376861aefd84901d2d63a7428a2.
//
// Solidity: event ListingCancelled(uint256 indexed listingId, address indexed seller)
func (_RWAMarketplace *RWAMarketplaceFilterer) WatchListingCancelled(opts *bind.WatchOpts, sink chan<- *RWAMarketplaceListingCancelled, listingId []*big.Int, seller []common.Address) (event.Subscription, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _RWAMarketplace.contract.WatchLogs(opts, "ListingCancelled", listingIdRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAMarketplaceListingCancelled)
				if err := _RWAMarketplace.contract.UnpackLog(event, "ListingCancelled", log); err != nil {
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

// ParseListingCancelled is a log parse operation binding the contract event 0x8e25282255ab31897df2b0456bb993ac7f84d376861aefd84901d2d63a7428a2.
//
// Solidity: event ListingCancelled(uint256 indexed listingId, address indexed seller)
func (_RWAMarketplace *RWAMarketplaceFilterer) ParseListingCancelled(log types.Log) (*RWAMarketplaceListingCancelled, error) {
	event := new(RWAMarketplaceListingCancelled)
	if err := _RWAMarketplace.contract.UnpackLog(event, "ListingCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAMarketplaceListingCreatedIterator is returned from FilterListingCreated and is used to iterate over the raw logs and unpacked data for ListingCreated events raised by the RWAMarketplace contract.
type RWAMarketplaceListingCreatedIterator struct {
	Event *RWAMarketplaceListingCreated // Event containing the contract specifics and raw log

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
func (it *RWAMarketplaceListingCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAMarketplaceListingCreated)
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
		it.Event = new(RWAMarketplaceListingCreated)
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
func (it *RWAMarketplaceListingCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAMarketplaceListingCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAMarketplaceListingCreated represents a ListingCreated event raised by the RWAMarketplace contract.
type RWAMarketplaceListingCreated struct {
	ListingId *big.Int
	AssetId   *big.Int
	Seller    common.Address
	Amount    *big.Int
	Price     *big.Int
	OrderType uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterListingCreated is a free log retrieval operation binding the contract event 0x8c402037981c4d37438b522257265c63606f2178505000e17f0cdaa1c1b14bde.
//
// Solidity: event ListingCreated(uint256 indexed listingId, uint256 indexed assetId, address indexed seller, uint256 amount, uint256 price, uint8 orderType)
func (_RWAMarketplace *RWAMarketplaceFilterer) FilterListingCreated(opts *bind.FilterOpts, listingId []*big.Int, assetId []*big.Int, seller []common.Address) (*RWAMarketplaceListingCreatedIterator, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _RWAMarketplace.contract.FilterLogs(opts, "ListingCreated", listingIdRule, assetIdRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return &RWAMarketplaceListingCreatedIterator{contract: _RWAMarketplace.contract, event: "ListingCreated", logs: logs, sub: sub}, nil
}

// WatchListingCreated is a free log subscription operation binding the contract event 0x8c402037981c4d37438b522257265c63606f2178505000e17f0cdaa1c1b14bde.
//
// Solidity: event ListingCreated(uint256 indexed listingId, uint256 indexed assetId, address indexed seller, uint256 amount, uint256 price, uint8 orderType)
func (_RWAMarketplace *RWAMarketplaceFilterer) WatchListingCreated(opts *bind.WatchOpts, sink chan<- *RWAMarketplaceListingCreated, listingId []*big.Int, assetId []*big.Int, seller []common.Address) (event.Subscription, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _RWAMarketplace.contract.WatchLogs(opts, "ListingCreated", listingIdRule, assetIdRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAMarketplaceListingCreated)
				if err := _RWAMarketplace.contract.UnpackLog(event, "ListingCreated", log); err != nil {
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

// ParseListingCreated is a log parse operation binding the contract event 0x8c402037981c4d37438b522257265c63606f2178505000e17f0cdaa1c1b14bde.
//
// Solidity: event ListingCreated(uint256 indexed listingId, uint256 indexed assetId, address indexed seller, uint256 amount, uint256 price, uint8 orderType)
func (_RWAMarketplace *RWAMarketplaceFilterer) ParseListingCreated(log types.Log) (*RWAMarketplaceListingCreated, error) {
	event := new(RWAMarketplaceListingCreated)
	if err := _RWAMarketplace.contract.UnpackLog(event, "ListingCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAMarketplacePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the RWAMarketplace contract.
type RWAMarketplacePausedIterator struct {
	Event *RWAMarketplacePaused // Event containing the contract specifics and raw log

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
func (it *RWAMarketplacePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAMarketplacePaused)
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
		it.Event = new(RWAMarketplacePaused)
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
func (it *RWAMarketplacePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAMarketplacePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAMarketplacePaused represents a Paused event raised by the RWAMarketplace contract.
type RWAMarketplacePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RWAMarketplace *RWAMarketplaceFilterer) FilterPaused(opts *bind.FilterOpts) (*RWAMarketplacePausedIterator, error) {

	logs, sub, err := _RWAMarketplace.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &RWAMarketplacePausedIterator{contract: _RWAMarketplace.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RWAMarketplace *RWAMarketplaceFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *RWAMarketplacePaused) (event.Subscription, error) {

	logs, sub, err := _RWAMarketplace.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAMarketplacePaused)
				if err := _RWAMarketplace.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_RWAMarketplace *RWAMarketplaceFilterer) ParsePaused(log types.Log) (*RWAMarketplacePaused, error) {
	event := new(RWAMarketplacePaused)
	if err := _RWAMarketplace.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAMarketplaceRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the RWAMarketplace contract.
type RWAMarketplaceRoleAdminChangedIterator struct {
	Event *RWAMarketplaceRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *RWAMarketplaceRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAMarketplaceRoleAdminChanged)
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
		it.Event = new(RWAMarketplaceRoleAdminChanged)
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
func (it *RWAMarketplaceRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAMarketplaceRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAMarketplaceRoleAdminChanged represents a RoleAdminChanged event raised by the RWAMarketplace contract.
type RWAMarketplaceRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RWAMarketplace *RWAMarketplaceFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*RWAMarketplaceRoleAdminChangedIterator, error) {

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

	logs, sub, err := _RWAMarketplace.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &RWAMarketplaceRoleAdminChangedIterator{contract: _RWAMarketplace.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RWAMarketplace *RWAMarketplaceFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *RWAMarketplaceRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _RWAMarketplace.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAMarketplaceRoleAdminChanged)
				if err := _RWAMarketplace.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_RWAMarketplace *RWAMarketplaceFilterer) ParseRoleAdminChanged(log types.Log) (*RWAMarketplaceRoleAdminChanged, error) {
	event := new(RWAMarketplaceRoleAdminChanged)
	if err := _RWAMarketplace.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAMarketplaceRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the RWAMarketplace contract.
type RWAMarketplaceRoleGrantedIterator struct {
	Event *RWAMarketplaceRoleGranted // Event containing the contract specifics and raw log

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
func (it *RWAMarketplaceRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAMarketplaceRoleGranted)
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
		it.Event = new(RWAMarketplaceRoleGranted)
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
func (it *RWAMarketplaceRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAMarketplaceRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAMarketplaceRoleGranted represents a RoleGranted event raised by the RWAMarketplace contract.
type RWAMarketplaceRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAMarketplace *RWAMarketplaceFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RWAMarketplaceRoleGrantedIterator, error) {

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

	logs, sub, err := _RWAMarketplace.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RWAMarketplaceRoleGrantedIterator{contract: _RWAMarketplace.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAMarketplace *RWAMarketplaceFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *RWAMarketplaceRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RWAMarketplace.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAMarketplaceRoleGranted)
				if err := _RWAMarketplace.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_RWAMarketplace *RWAMarketplaceFilterer) ParseRoleGranted(log types.Log) (*RWAMarketplaceRoleGranted, error) {
	event := new(RWAMarketplaceRoleGranted)
	if err := _RWAMarketplace.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAMarketplaceRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the RWAMarketplace contract.
type RWAMarketplaceRoleRevokedIterator struct {
	Event *RWAMarketplaceRoleRevoked // Event containing the contract specifics and raw log

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
func (it *RWAMarketplaceRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAMarketplaceRoleRevoked)
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
		it.Event = new(RWAMarketplaceRoleRevoked)
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
func (it *RWAMarketplaceRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAMarketplaceRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAMarketplaceRoleRevoked represents a RoleRevoked event raised by the RWAMarketplace contract.
type RWAMarketplaceRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAMarketplace *RWAMarketplaceFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RWAMarketplaceRoleRevokedIterator, error) {

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

	logs, sub, err := _RWAMarketplace.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RWAMarketplaceRoleRevokedIterator{contract: _RWAMarketplace.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAMarketplace *RWAMarketplaceFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *RWAMarketplaceRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RWAMarketplace.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAMarketplaceRoleRevoked)
				if err := _RWAMarketplace.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_RWAMarketplace *RWAMarketplaceFilterer) ParseRoleRevoked(log types.Log) (*RWAMarketplaceRoleRevoked, error) {
	event := new(RWAMarketplaceRoleRevoked)
	if err := _RWAMarketplace.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAMarketplaceTradExecutedIterator is returned from FilterTradExecuted and is used to iterate over the raw logs and unpacked data for TradExecuted events raised by the RWAMarketplace contract.
type RWAMarketplaceTradExecutedIterator struct {
	Event *RWAMarketplaceTradExecuted // Event containing the contract specifics and raw log

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
func (it *RWAMarketplaceTradExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAMarketplaceTradExecuted)
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
		it.Event = new(RWAMarketplaceTradExecuted)
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
func (it *RWAMarketplaceTradExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAMarketplaceTradExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAMarketplaceTradExecuted represents a TradExecuted event raised by the RWAMarketplace contract.
type RWAMarketplaceTradExecuted struct {
	TradeId    *big.Int
	ListingId  *big.Int
	Buyer      common.Address
	Seller     common.Address
	Amount     *big.Int
	Price      *big.Int
	TotalValue *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTradExecuted is a free log retrieval operation binding the contract event 0x32db0875308bdbe0d95a6d58fcdba95156d11db2de1045f86fcbd7eb3acd310f.
//
// Solidity: event TradExecuted(uint256 indexed tradeId, uint256 indexed listingId, address indexed buyer, address seller, uint256 amount, uint256 price, uint256 totalValue)
func (_RWAMarketplace *RWAMarketplaceFilterer) FilterTradExecuted(opts *bind.FilterOpts, tradeId []*big.Int, listingId []*big.Int, buyer []common.Address) (*RWAMarketplaceTradExecutedIterator, error) {

	var tradeIdRule []interface{}
	for _, tradeIdItem := range tradeId {
		tradeIdRule = append(tradeIdRule, tradeIdItem)
	}
	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _RWAMarketplace.contract.FilterLogs(opts, "TradExecuted", tradeIdRule, listingIdRule, buyerRule)
	if err != nil {
		return nil, err
	}
	return &RWAMarketplaceTradExecutedIterator{contract: _RWAMarketplace.contract, event: "TradExecuted", logs: logs, sub: sub}, nil
}

// WatchTradExecuted is a free log subscription operation binding the contract event 0x32db0875308bdbe0d95a6d58fcdba95156d11db2de1045f86fcbd7eb3acd310f.
//
// Solidity: event TradExecuted(uint256 indexed tradeId, uint256 indexed listingId, address indexed buyer, address seller, uint256 amount, uint256 price, uint256 totalValue)
func (_RWAMarketplace *RWAMarketplaceFilterer) WatchTradExecuted(opts *bind.WatchOpts, sink chan<- *RWAMarketplaceTradExecuted, tradeId []*big.Int, listingId []*big.Int, buyer []common.Address) (event.Subscription, error) {

	var tradeIdRule []interface{}
	for _, tradeIdItem := range tradeId {
		tradeIdRule = append(tradeIdRule, tradeIdItem)
	}
	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _RWAMarketplace.contract.WatchLogs(opts, "TradExecuted", tradeIdRule, listingIdRule, buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAMarketplaceTradExecuted)
				if err := _RWAMarketplace.contract.UnpackLog(event, "TradExecuted", log); err != nil {
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

// ParseTradExecuted is a log parse operation binding the contract event 0x32db0875308bdbe0d95a6d58fcdba95156d11db2de1045f86fcbd7eb3acd310f.
//
// Solidity: event TradExecuted(uint256 indexed tradeId, uint256 indexed listingId, address indexed buyer, address seller, uint256 amount, uint256 price, uint256 totalValue)
func (_RWAMarketplace *RWAMarketplaceFilterer) ParseTradExecuted(log types.Log) (*RWAMarketplaceTradExecuted, error) {
	event := new(RWAMarketplaceTradExecuted)
	if err := _RWAMarketplace.contract.UnpackLog(event, "TradExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAMarketplaceUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the RWAMarketplace contract.
type RWAMarketplaceUnpausedIterator struct {
	Event *RWAMarketplaceUnpaused // Event containing the contract specifics and raw log

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
func (it *RWAMarketplaceUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAMarketplaceUnpaused)
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
		it.Event = new(RWAMarketplaceUnpaused)
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
func (it *RWAMarketplaceUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAMarketplaceUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAMarketplaceUnpaused represents a Unpaused event raised by the RWAMarketplace contract.
type RWAMarketplaceUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RWAMarketplace *RWAMarketplaceFilterer) FilterUnpaused(opts *bind.FilterOpts) (*RWAMarketplaceUnpausedIterator, error) {

	logs, sub, err := _RWAMarketplace.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &RWAMarketplaceUnpausedIterator{contract: _RWAMarketplace.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RWAMarketplace *RWAMarketplaceFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *RWAMarketplaceUnpaused) (event.Subscription, error) {

	logs, sub, err := _RWAMarketplace.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAMarketplaceUnpaused)
				if err := _RWAMarketplace.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_RWAMarketplace *RWAMarketplaceFilterer) ParseUnpaused(log types.Log) (*RWAMarketplaceUnpaused, error) {
	event := new(RWAMarketplaceUnpaused)
	if err := _RWAMarketplace.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
