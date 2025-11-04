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

// RWAGovernanceGovernanceParams is an auto generated low-level Go binding around an user-defined struct.
type RWAGovernanceGovernanceParams struct {
	ProposalThreshold *big.Int
	QuorumThreshold   *big.Int
	ApprovalThreshold *big.Int
	VotingPeriod      *big.Int
	TimelockPeriod    *big.Int
	IssuerVetoEnabled bool
}

// RWAGovernanceProposal is an auto generated low-level Go binding around an user-defined struct.
type RWAGovernanceProposal struct {
	ProposalId       *big.Int
	AssetId          *big.Int
	Proposer         common.Address
	ProposalType     uint8
	Status           uint8
	Description      string
	ExecutionData    []byte
	VotingStartTime  *big.Int
	VotingEndTime    *big.Int
	ExecutionTime    *big.Int
	ForVotes         *big.Int
	AgainstVotes     *big.Int
	AbstainVotes     *big.Int
	TotalVotingPower *big.Int
	Executed         bool
	Vetoed           bool
}

// RWAGovernanceVoteReceipt is an auto generated low-level Go binding around an user-defined struct.
type RWAGovernanceVoteReceipt struct {
	HasVoted bool
	Choice   uint8
	Votes    *big.Int
}

// RWAGovernanceMetaData contains all meta data concerning the RWAGovernance contract.
var RWAGovernanceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"ProposalCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumRWAGovernance.ProposalType\",\"name\":\"proposalType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"name\":\"ProposalCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"ProposalExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"vetoer\",\"type\":\"address\"}],\"name\":\"ProposalVetoed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumRWAGovernance.VoteChoice\",\"name\":\"choice\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"}],\"name\":\"VoteCast\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"VoteDelegated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_APPROVAL\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_PROPOSAL_THRESHOLD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_QUORUM\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_TIMELOCK\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_VOTING_PERIOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GOVERNANCE_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PERCENTAGE_BASE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"assetFactory\",\"outputs\":[{\"internalType\":\"contractIRWAAsset\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"cancelProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"enumRWAGovernance.VoteChoice\",\"name\":\"choice\",\"type\":\"uint8\"}],\"name\":\"castVote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"enumRWAGovernance.ProposalType\",\"name\":\"proposalType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"executionData\",\"type\":\"bytes\"}],\"name\":\"createProposal\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"}],\"name\":\"delegateVote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"executeProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getActiveProposals\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getAssetProposals\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"getGovernanceParams\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"proposalThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"quorumThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"approvalThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"votingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timelockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"issuerVetoEnabled\",\"type\":\"bool\"}],\"internalType\":\"structRWAGovernance.GovernanceParams\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"getProposal\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"enumRWAGovernance.ProposalType\",\"name\":\"proposalType\",\"type\":\"uint8\"},{\"internalType\":\"enumRWAGovernance.ProposalStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"executionData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"votingStartTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"votingEndTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"executionTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"forVotes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"againstVotes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"abstainVotes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalVotingPower\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"executed\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"vetoed\",\"type\":\"bool\"}],\"internalType\":\"structRWAGovernance.Proposal\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"}],\"name\":\"getVoteReceipt\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"hasVoted\",\"type\":\"bool\"},{\"internalType\":\"enumRWAGovernance.VoteChoice\",\"name\":\"choice\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"votes\",\"type\":\"uint256\"}],\"internalType\":\"structRWAGovernance.VoteReceipt\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"hasProposalPassed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"proposalThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"quorumThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"approvalThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"votingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timelockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"issuerVetoEnabled\",\"type\":\"bool\"}],\"internalType\":\"structRWAGovernance.GovernanceParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"setGovernanceParams\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalVotes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetId\",\"type\":\"uint256\"}],\"name\":\"undelegateVote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"vetoProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RWAGovernanceABI is the input ABI used to generate the binding from.
// Deprecated: Use RWAGovernanceMetaData.ABI instead.
var RWAGovernanceABI = RWAGovernanceMetaData.ABI

// RWAGovernance is an auto generated Go binding around an Ethereum contract.
type RWAGovernance struct {
	RWAGovernanceCaller     // Read-only binding to the contract
	RWAGovernanceTransactor // Write-only binding to the contract
	RWAGovernanceFilterer   // Log filterer for contract events
}

// RWAGovernanceCaller is an auto generated read-only Go binding around an Ethereum contract.
type RWAGovernanceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAGovernanceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RWAGovernanceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAGovernanceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RWAGovernanceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RWAGovernanceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RWAGovernanceSession struct {
	Contract     *RWAGovernance    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RWAGovernanceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RWAGovernanceCallerSession struct {
	Contract *RWAGovernanceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// RWAGovernanceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RWAGovernanceTransactorSession struct {
	Contract     *RWAGovernanceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// RWAGovernanceRaw is an auto generated low-level Go binding around an Ethereum contract.
type RWAGovernanceRaw struct {
	Contract *RWAGovernance // Generic contract binding to access the raw methods on
}

// RWAGovernanceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RWAGovernanceCallerRaw struct {
	Contract *RWAGovernanceCaller // Generic read-only contract binding to access the raw methods on
}

// RWAGovernanceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RWAGovernanceTransactorRaw struct {
	Contract *RWAGovernanceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRWAGovernance creates a new instance of RWAGovernance, bound to a specific deployed contract.
func NewRWAGovernance(address common.Address, backend bind.ContractBackend) (*RWAGovernance, error) {
	contract, err := bindRWAGovernance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RWAGovernance{RWAGovernanceCaller: RWAGovernanceCaller{contract: contract}, RWAGovernanceTransactor: RWAGovernanceTransactor{contract: contract}, RWAGovernanceFilterer: RWAGovernanceFilterer{contract: contract}}, nil
}

// NewRWAGovernanceCaller creates a new read-only instance of RWAGovernance, bound to a specific deployed contract.
func NewRWAGovernanceCaller(address common.Address, caller bind.ContractCaller) (*RWAGovernanceCaller, error) {
	contract, err := bindRWAGovernance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RWAGovernanceCaller{contract: contract}, nil
}

// NewRWAGovernanceTransactor creates a new write-only instance of RWAGovernance, bound to a specific deployed contract.
func NewRWAGovernanceTransactor(address common.Address, transactor bind.ContractTransactor) (*RWAGovernanceTransactor, error) {
	contract, err := bindRWAGovernance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RWAGovernanceTransactor{contract: contract}, nil
}

// NewRWAGovernanceFilterer creates a new log filterer instance of RWAGovernance, bound to a specific deployed contract.
func NewRWAGovernanceFilterer(address common.Address, filterer bind.ContractFilterer) (*RWAGovernanceFilterer, error) {
	contract, err := bindRWAGovernance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RWAGovernanceFilterer{contract: contract}, nil
}

// bindRWAGovernance binds a generic wrapper to an already deployed contract.
func bindRWAGovernance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RWAGovernanceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RWAGovernance *RWAGovernanceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RWAGovernance.Contract.RWAGovernanceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RWAGovernance *RWAGovernanceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAGovernance.Contract.RWAGovernanceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RWAGovernance *RWAGovernanceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RWAGovernance.Contract.RWAGovernanceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RWAGovernance *RWAGovernanceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RWAGovernance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RWAGovernance *RWAGovernanceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RWAGovernance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RWAGovernance *RWAGovernanceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RWAGovernance.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAGovernance *RWAGovernanceCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAGovernance *RWAGovernanceSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RWAGovernance.Contract.DEFAULTADMINROLE(&_RWAGovernance.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_RWAGovernance *RWAGovernanceCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _RWAGovernance.Contract.DEFAULTADMINROLE(&_RWAGovernance.CallOpts)
}

// DEFAULTAPPROVAL is a free data retrieval call binding the contract method 0xdcba46d9.
//
// Solidity: function DEFAULT_APPROVAL() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCaller) DEFAULTAPPROVAL(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "DEFAULT_APPROVAL")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEFAULTAPPROVAL is a free data retrieval call binding the contract method 0xdcba46d9.
//
// Solidity: function DEFAULT_APPROVAL() view returns(uint256)
func (_RWAGovernance *RWAGovernanceSession) DEFAULTAPPROVAL() (*big.Int, error) {
	return _RWAGovernance.Contract.DEFAULTAPPROVAL(&_RWAGovernance.CallOpts)
}

// DEFAULTAPPROVAL is a free data retrieval call binding the contract method 0xdcba46d9.
//
// Solidity: function DEFAULT_APPROVAL() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCallerSession) DEFAULTAPPROVAL() (*big.Int, error) {
	return _RWAGovernance.Contract.DEFAULTAPPROVAL(&_RWAGovernance.CallOpts)
}

// DEFAULTPROPOSALTHRESHOLD is a free data retrieval call binding the contract method 0x66efd6ca.
//
// Solidity: function DEFAULT_PROPOSAL_THRESHOLD() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCaller) DEFAULTPROPOSALTHRESHOLD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "DEFAULT_PROPOSAL_THRESHOLD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEFAULTPROPOSALTHRESHOLD is a free data retrieval call binding the contract method 0x66efd6ca.
//
// Solidity: function DEFAULT_PROPOSAL_THRESHOLD() view returns(uint256)
func (_RWAGovernance *RWAGovernanceSession) DEFAULTPROPOSALTHRESHOLD() (*big.Int, error) {
	return _RWAGovernance.Contract.DEFAULTPROPOSALTHRESHOLD(&_RWAGovernance.CallOpts)
}

// DEFAULTPROPOSALTHRESHOLD is a free data retrieval call binding the contract method 0x66efd6ca.
//
// Solidity: function DEFAULT_PROPOSAL_THRESHOLD() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCallerSession) DEFAULTPROPOSALTHRESHOLD() (*big.Int, error) {
	return _RWAGovernance.Contract.DEFAULTPROPOSALTHRESHOLD(&_RWAGovernance.CallOpts)
}

// DEFAULTQUORUM is a free data retrieval call binding the contract method 0x9a341b9f.
//
// Solidity: function DEFAULT_QUORUM() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCaller) DEFAULTQUORUM(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "DEFAULT_QUORUM")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEFAULTQUORUM is a free data retrieval call binding the contract method 0x9a341b9f.
//
// Solidity: function DEFAULT_QUORUM() view returns(uint256)
func (_RWAGovernance *RWAGovernanceSession) DEFAULTQUORUM() (*big.Int, error) {
	return _RWAGovernance.Contract.DEFAULTQUORUM(&_RWAGovernance.CallOpts)
}

// DEFAULTQUORUM is a free data retrieval call binding the contract method 0x9a341b9f.
//
// Solidity: function DEFAULT_QUORUM() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCallerSession) DEFAULTQUORUM() (*big.Int, error) {
	return _RWAGovernance.Contract.DEFAULTQUORUM(&_RWAGovernance.CallOpts)
}

// DEFAULTTIMELOCK is a free data retrieval call binding the contract method 0xc1c7fca3.
//
// Solidity: function DEFAULT_TIMELOCK() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCaller) DEFAULTTIMELOCK(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "DEFAULT_TIMELOCK")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEFAULTTIMELOCK is a free data retrieval call binding the contract method 0xc1c7fca3.
//
// Solidity: function DEFAULT_TIMELOCK() view returns(uint256)
func (_RWAGovernance *RWAGovernanceSession) DEFAULTTIMELOCK() (*big.Int, error) {
	return _RWAGovernance.Contract.DEFAULTTIMELOCK(&_RWAGovernance.CallOpts)
}

// DEFAULTTIMELOCK is a free data retrieval call binding the contract method 0xc1c7fca3.
//
// Solidity: function DEFAULT_TIMELOCK() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCallerSession) DEFAULTTIMELOCK() (*big.Int, error) {
	return _RWAGovernance.Contract.DEFAULTTIMELOCK(&_RWAGovernance.CallOpts)
}

// DEFAULTVOTINGPERIOD is a free data retrieval call binding the contract method 0xaedfe53f.
//
// Solidity: function DEFAULT_VOTING_PERIOD() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCaller) DEFAULTVOTINGPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "DEFAULT_VOTING_PERIOD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEFAULTVOTINGPERIOD is a free data retrieval call binding the contract method 0xaedfe53f.
//
// Solidity: function DEFAULT_VOTING_PERIOD() view returns(uint256)
func (_RWAGovernance *RWAGovernanceSession) DEFAULTVOTINGPERIOD() (*big.Int, error) {
	return _RWAGovernance.Contract.DEFAULTVOTINGPERIOD(&_RWAGovernance.CallOpts)
}

// DEFAULTVOTINGPERIOD is a free data retrieval call binding the contract method 0xaedfe53f.
//
// Solidity: function DEFAULT_VOTING_PERIOD() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCallerSession) DEFAULTVOTINGPERIOD() (*big.Int, error) {
	return _RWAGovernance.Contract.DEFAULTVOTINGPERIOD(&_RWAGovernance.CallOpts)
}

// GOVERNANCEADMINROLE is a free data retrieval call binding the contract method 0x77fb8d7a.
//
// Solidity: function GOVERNANCE_ADMIN_ROLE() view returns(bytes32)
func (_RWAGovernance *RWAGovernanceCaller) GOVERNANCEADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "GOVERNANCE_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GOVERNANCEADMINROLE is a free data retrieval call binding the contract method 0x77fb8d7a.
//
// Solidity: function GOVERNANCE_ADMIN_ROLE() view returns(bytes32)
func (_RWAGovernance *RWAGovernanceSession) GOVERNANCEADMINROLE() ([32]byte, error) {
	return _RWAGovernance.Contract.GOVERNANCEADMINROLE(&_RWAGovernance.CallOpts)
}

// GOVERNANCEADMINROLE is a free data retrieval call binding the contract method 0x77fb8d7a.
//
// Solidity: function GOVERNANCE_ADMIN_ROLE() view returns(bytes32)
func (_RWAGovernance *RWAGovernanceCallerSession) GOVERNANCEADMINROLE() ([32]byte, error) {
	return _RWAGovernance.Contract.GOVERNANCEADMINROLE(&_RWAGovernance.CallOpts)
}

// PERCENTAGEBASE is a free data retrieval call binding the contract method 0x87c13943.
//
// Solidity: function PERCENTAGE_BASE() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCaller) PERCENTAGEBASE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "PERCENTAGE_BASE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PERCENTAGEBASE is a free data retrieval call binding the contract method 0x87c13943.
//
// Solidity: function PERCENTAGE_BASE() view returns(uint256)
func (_RWAGovernance *RWAGovernanceSession) PERCENTAGEBASE() (*big.Int, error) {
	return _RWAGovernance.Contract.PERCENTAGEBASE(&_RWAGovernance.CallOpts)
}

// PERCENTAGEBASE is a free data retrieval call binding the contract method 0x87c13943.
//
// Solidity: function PERCENTAGE_BASE() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCallerSession) PERCENTAGEBASE() (*big.Int, error) {
	return _RWAGovernance.Contract.PERCENTAGEBASE(&_RWAGovernance.CallOpts)
}

// AssetFactory is a free data retrieval call binding the contract method 0xc3acb4d1.
//
// Solidity: function assetFactory() view returns(address)
func (_RWAGovernance *RWAGovernanceCaller) AssetFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "assetFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AssetFactory is a free data retrieval call binding the contract method 0xc3acb4d1.
//
// Solidity: function assetFactory() view returns(address)
func (_RWAGovernance *RWAGovernanceSession) AssetFactory() (common.Address, error) {
	return _RWAGovernance.Contract.AssetFactory(&_RWAGovernance.CallOpts)
}

// AssetFactory is a free data retrieval call binding the contract method 0xc3acb4d1.
//
// Solidity: function assetFactory() view returns(address)
func (_RWAGovernance *RWAGovernanceCallerSession) AssetFactory() (common.Address, error) {
	return _RWAGovernance.Contract.AssetFactory(&_RWAGovernance.CallOpts)
}

// GetActiveProposals is a free data retrieval call binding the contract method 0x6213e5ea.
//
// Solidity: function getActiveProposals(uint256 assetId) view returns(uint256[])
func (_RWAGovernance *RWAGovernanceCaller) GetActiveProposals(opts *bind.CallOpts, assetId *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "getActiveProposals", assetId)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetActiveProposals is a free data retrieval call binding the contract method 0x6213e5ea.
//
// Solidity: function getActiveProposals(uint256 assetId) view returns(uint256[])
func (_RWAGovernance *RWAGovernanceSession) GetActiveProposals(assetId *big.Int) ([]*big.Int, error) {
	return _RWAGovernance.Contract.GetActiveProposals(&_RWAGovernance.CallOpts, assetId)
}

// GetActiveProposals is a free data retrieval call binding the contract method 0x6213e5ea.
//
// Solidity: function getActiveProposals(uint256 assetId) view returns(uint256[])
func (_RWAGovernance *RWAGovernanceCallerSession) GetActiveProposals(assetId *big.Int) ([]*big.Int, error) {
	return _RWAGovernance.Contract.GetActiveProposals(&_RWAGovernance.CallOpts, assetId)
}

// GetAssetProposals is a free data retrieval call binding the contract method 0x0443e428.
//
// Solidity: function getAssetProposals(uint256 assetId) view returns(uint256[])
func (_RWAGovernance *RWAGovernanceCaller) GetAssetProposals(opts *bind.CallOpts, assetId *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "getAssetProposals", assetId)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAssetProposals is a free data retrieval call binding the contract method 0x0443e428.
//
// Solidity: function getAssetProposals(uint256 assetId) view returns(uint256[])
func (_RWAGovernance *RWAGovernanceSession) GetAssetProposals(assetId *big.Int) ([]*big.Int, error) {
	return _RWAGovernance.Contract.GetAssetProposals(&_RWAGovernance.CallOpts, assetId)
}

// GetAssetProposals is a free data retrieval call binding the contract method 0x0443e428.
//
// Solidity: function getAssetProposals(uint256 assetId) view returns(uint256[])
func (_RWAGovernance *RWAGovernanceCallerSession) GetAssetProposals(assetId *big.Int) ([]*big.Int, error) {
	return _RWAGovernance.Contract.GetAssetProposals(&_RWAGovernance.CallOpts, assetId)
}

// GetGovernanceParams is a free data retrieval call binding the contract method 0xa1debd13.
//
// Solidity: function getGovernanceParams(uint256 assetId) view returns((uint256,uint256,uint256,uint256,uint256,bool))
func (_RWAGovernance *RWAGovernanceCaller) GetGovernanceParams(opts *bind.CallOpts, assetId *big.Int) (RWAGovernanceGovernanceParams, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "getGovernanceParams", assetId)

	if err != nil {
		return *new(RWAGovernanceGovernanceParams), err
	}

	out0 := *abi.ConvertType(out[0], new(RWAGovernanceGovernanceParams)).(*RWAGovernanceGovernanceParams)

	return out0, err

}

// GetGovernanceParams is a free data retrieval call binding the contract method 0xa1debd13.
//
// Solidity: function getGovernanceParams(uint256 assetId) view returns((uint256,uint256,uint256,uint256,uint256,bool))
func (_RWAGovernance *RWAGovernanceSession) GetGovernanceParams(assetId *big.Int) (RWAGovernanceGovernanceParams, error) {
	return _RWAGovernance.Contract.GetGovernanceParams(&_RWAGovernance.CallOpts, assetId)
}

// GetGovernanceParams is a free data retrieval call binding the contract method 0xa1debd13.
//
// Solidity: function getGovernanceParams(uint256 assetId) view returns((uint256,uint256,uint256,uint256,uint256,bool))
func (_RWAGovernance *RWAGovernanceCallerSession) GetGovernanceParams(assetId *big.Int) (RWAGovernanceGovernanceParams, error) {
	return _RWAGovernance.Contract.GetGovernanceParams(&_RWAGovernance.CallOpts, assetId)
}

// GetProposal is a free data retrieval call binding the contract method 0xc7f758a8.
//
// Solidity: function getProposal(uint256 proposalId) view returns((uint256,uint256,address,uint8,uint8,string,bytes,uint256,uint256,uint256,uint256,uint256,uint256,uint256,bool,bool))
func (_RWAGovernance *RWAGovernanceCaller) GetProposal(opts *bind.CallOpts, proposalId *big.Int) (RWAGovernanceProposal, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "getProposal", proposalId)

	if err != nil {
		return *new(RWAGovernanceProposal), err
	}

	out0 := *abi.ConvertType(out[0], new(RWAGovernanceProposal)).(*RWAGovernanceProposal)

	return out0, err

}

// GetProposal is a free data retrieval call binding the contract method 0xc7f758a8.
//
// Solidity: function getProposal(uint256 proposalId) view returns((uint256,uint256,address,uint8,uint8,string,bytes,uint256,uint256,uint256,uint256,uint256,uint256,uint256,bool,bool))
func (_RWAGovernance *RWAGovernanceSession) GetProposal(proposalId *big.Int) (RWAGovernanceProposal, error) {
	return _RWAGovernance.Contract.GetProposal(&_RWAGovernance.CallOpts, proposalId)
}

// GetProposal is a free data retrieval call binding the contract method 0xc7f758a8.
//
// Solidity: function getProposal(uint256 proposalId) view returns((uint256,uint256,address,uint8,uint8,string,bytes,uint256,uint256,uint256,uint256,uint256,uint256,uint256,bool,bool))
func (_RWAGovernance *RWAGovernanceCallerSession) GetProposal(proposalId *big.Int) (RWAGovernanceProposal, error) {
	return _RWAGovernance.Contract.GetProposal(&_RWAGovernance.CallOpts, proposalId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAGovernance *RWAGovernanceCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAGovernance *RWAGovernanceSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RWAGovernance.Contract.GetRoleAdmin(&_RWAGovernance.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_RWAGovernance *RWAGovernanceCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _RWAGovernance.Contract.GetRoleAdmin(&_RWAGovernance.CallOpts, role)
}

// GetVoteReceipt is a free data retrieval call binding the contract method 0x920f80eb.
//
// Solidity: function getVoteReceipt(uint256 proposalId, address voter) view returns((bool,uint8,uint256))
func (_RWAGovernance *RWAGovernanceCaller) GetVoteReceipt(opts *bind.CallOpts, proposalId *big.Int, voter common.Address) (RWAGovernanceVoteReceipt, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "getVoteReceipt", proposalId, voter)

	if err != nil {
		return *new(RWAGovernanceVoteReceipt), err
	}

	out0 := *abi.ConvertType(out[0], new(RWAGovernanceVoteReceipt)).(*RWAGovernanceVoteReceipt)

	return out0, err

}

// GetVoteReceipt is a free data retrieval call binding the contract method 0x920f80eb.
//
// Solidity: function getVoteReceipt(uint256 proposalId, address voter) view returns((bool,uint8,uint256))
func (_RWAGovernance *RWAGovernanceSession) GetVoteReceipt(proposalId *big.Int, voter common.Address) (RWAGovernanceVoteReceipt, error) {
	return _RWAGovernance.Contract.GetVoteReceipt(&_RWAGovernance.CallOpts, proposalId, voter)
}

// GetVoteReceipt is a free data retrieval call binding the contract method 0x920f80eb.
//
// Solidity: function getVoteReceipt(uint256 proposalId, address voter) view returns((bool,uint8,uint256))
func (_RWAGovernance *RWAGovernanceCallerSession) GetVoteReceipt(proposalId *big.Int, voter common.Address) (RWAGovernanceVoteReceipt, error) {
	return _RWAGovernance.Contract.GetVoteReceipt(&_RWAGovernance.CallOpts, proposalId, voter)
}

// HasProposalPassed is a free data retrieval call binding the contract method 0x6414b6a1.
//
// Solidity: function hasProposalPassed(uint256 proposalId) view returns(bool)
func (_RWAGovernance *RWAGovernanceCaller) HasProposalPassed(opts *bind.CallOpts, proposalId *big.Int) (bool, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "hasProposalPassed", proposalId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasProposalPassed is a free data retrieval call binding the contract method 0x6414b6a1.
//
// Solidity: function hasProposalPassed(uint256 proposalId) view returns(bool)
func (_RWAGovernance *RWAGovernanceSession) HasProposalPassed(proposalId *big.Int) (bool, error) {
	return _RWAGovernance.Contract.HasProposalPassed(&_RWAGovernance.CallOpts, proposalId)
}

// HasProposalPassed is a free data retrieval call binding the contract method 0x6414b6a1.
//
// Solidity: function hasProposalPassed(uint256 proposalId) view returns(bool)
func (_RWAGovernance *RWAGovernanceCallerSession) HasProposalPassed(proposalId *big.Int) (bool, error) {
	return _RWAGovernance.Contract.HasProposalPassed(&_RWAGovernance.CallOpts, proposalId)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAGovernance *RWAGovernanceCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAGovernance *RWAGovernanceSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RWAGovernance.Contract.HasRole(&_RWAGovernance.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_RWAGovernance *RWAGovernanceCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _RWAGovernance.Contract.HasRole(&_RWAGovernance.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAGovernance *RWAGovernanceCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAGovernance *RWAGovernanceSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RWAGovernance.Contract.SupportsInterface(&_RWAGovernance.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_RWAGovernance *RWAGovernanceCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _RWAGovernance.Contract.SupportsInterface(&_RWAGovernance.CallOpts, interfaceId)
}

// TotalProposals is a free data retrieval call binding the contract method 0xa78d80fc.
//
// Solidity: function totalProposals() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCaller) TotalProposals(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "totalProposals")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalProposals is a free data retrieval call binding the contract method 0xa78d80fc.
//
// Solidity: function totalProposals() view returns(uint256)
func (_RWAGovernance *RWAGovernanceSession) TotalProposals() (*big.Int, error) {
	return _RWAGovernance.Contract.TotalProposals(&_RWAGovernance.CallOpts)
}

// TotalProposals is a free data retrieval call binding the contract method 0xa78d80fc.
//
// Solidity: function totalProposals() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCallerSession) TotalProposals() (*big.Int, error) {
	return _RWAGovernance.Contract.TotalProposals(&_RWAGovernance.CallOpts)
}

// TotalVotes is a free data retrieval call binding the contract method 0x0d15fd77.
//
// Solidity: function totalVotes() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCaller) TotalVotes(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RWAGovernance.contract.Call(opts, &out, "totalVotes")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalVotes is a free data retrieval call binding the contract method 0x0d15fd77.
//
// Solidity: function totalVotes() view returns(uint256)
func (_RWAGovernance *RWAGovernanceSession) TotalVotes() (*big.Int, error) {
	return _RWAGovernance.Contract.TotalVotes(&_RWAGovernance.CallOpts)
}

// TotalVotes is a free data retrieval call binding the contract method 0x0d15fd77.
//
// Solidity: function totalVotes() view returns(uint256)
func (_RWAGovernance *RWAGovernanceCallerSession) TotalVotes() (*big.Int, error) {
	return _RWAGovernance.Contract.TotalVotes(&_RWAGovernance.CallOpts)
}

// CancelProposal is a paid mutator transaction binding the contract method 0xe0a8f6f5.
//
// Solidity: function cancelProposal(uint256 proposalId) returns()
func (_RWAGovernance *RWAGovernanceTransactor) CancelProposal(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _RWAGovernance.contract.Transact(opts, "cancelProposal", proposalId)
}

// CancelProposal is a paid mutator transaction binding the contract method 0xe0a8f6f5.
//
// Solidity: function cancelProposal(uint256 proposalId) returns()
func (_RWAGovernance *RWAGovernanceSession) CancelProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _RWAGovernance.Contract.CancelProposal(&_RWAGovernance.TransactOpts, proposalId)
}

// CancelProposal is a paid mutator transaction binding the contract method 0xe0a8f6f5.
//
// Solidity: function cancelProposal(uint256 proposalId) returns()
func (_RWAGovernance *RWAGovernanceTransactorSession) CancelProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _RWAGovernance.Contract.CancelProposal(&_RWAGovernance.TransactOpts, proposalId)
}

// CastVote is a paid mutator transaction binding the contract method 0x56781388.
//
// Solidity: function castVote(uint256 proposalId, uint8 choice) returns()
func (_RWAGovernance *RWAGovernanceTransactor) CastVote(opts *bind.TransactOpts, proposalId *big.Int, choice uint8) (*types.Transaction, error) {
	return _RWAGovernance.contract.Transact(opts, "castVote", proposalId, choice)
}

// CastVote is a paid mutator transaction binding the contract method 0x56781388.
//
// Solidity: function castVote(uint256 proposalId, uint8 choice) returns()
func (_RWAGovernance *RWAGovernanceSession) CastVote(proposalId *big.Int, choice uint8) (*types.Transaction, error) {
	return _RWAGovernance.Contract.CastVote(&_RWAGovernance.TransactOpts, proposalId, choice)
}

// CastVote is a paid mutator transaction binding the contract method 0x56781388.
//
// Solidity: function castVote(uint256 proposalId, uint8 choice) returns()
func (_RWAGovernance *RWAGovernanceTransactorSession) CastVote(proposalId *big.Int, choice uint8) (*types.Transaction, error) {
	return _RWAGovernance.Contract.CastVote(&_RWAGovernance.TransactOpts, proposalId, choice)
}

// CreateProposal is a paid mutator transaction binding the contract method 0x1b4aad32.
//
// Solidity: function createProposal(uint256 assetId, uint8 proposalType, string description, bytes executionData) returns(uint256 proposalId)
func (_RWAGovernance *RWAGovernanceTransactor) CreateProposal(opts *bind.TransactOpts, assetId *big.Int, proposalType uint8, description string, executionData []byte) (*types.Transaction, error) {
	return _RWAGovernance.contract.Transact(opts, "createProposal", assetId, proposalType, description, executionData)
}

// CreateProposal is a paid mutator transaction binding the contract method 0x1b4aad32.
//
// Solidity: function createProposal(uint256 assetId, uint8 proposalType, string description, bytes executionData) returns(uint256 proposalId)
func (_RWAGovernance *RWAGovernanceSession) CreateProposal(assetId *big.Int, proposalType uint8, description string, executionData []byte) (*types.Transaction, error) {
	return _RWAGovernance.Contract.CreateProposal(&_RWAGovernance.TransactOpts, assetId, proposalType, description, executionData)
}

// CreateProposal is a paid mutator transaction binding the contract method 0x1b4aad32.
//
// Solidity: function createProposal(uint256 assetId, uint8 proposalType, string description, bytes executionData) returns(uint256 proposalId)
func (_RWAGovernance *RWAGovernanceTransactorSession) CreateProposal(assetId *big.Int, proposalType uint8, description string, executionData []byte) (*types.Transaction, error) {
	return _RWAGovernance.Contract.CreateProposal(&_RWAGovernance.TransactOpts, assetId, proposalType, description, executionData)
}

// DelegateVote is a paid mutator transaction binding the contract method 0xac71fe18.
//
// Solidity: function delegateVote(uint256 assetId, address delegate) returns()
func (_RWAGovernance *RWAGovernanceTransactor) DelegateVote(opts *bind.TransactOpts, assetId *big.Int, delegate common.Address) (*types.Transaction, error) {
	return _RWAGovernance.contract.Transact(opts, "delegateVote", assetId, delegate)
}

// DelegateVote is a paid mutator transaction binding the contract method 0xac71fe18.
//
// Solidity: function delegateVote(uint256 assetId, address delegate) returns()
func (_RWAGovernance *RWAGovernanceSession) DelegateVote(assetId *big.Int, delegate common.Address) (*types.Transaction, error) {
	return _RWAGovernance.Contract.DelegateVote(&_RWAGovernance.TransactOpts, assetId, delegate)
}

// DelegateVote is a paid mutator transaction binding the contract method 0xac71fe18.
//
// Solidity: function delegateVote(uint256 assetId, address delegate) returns()
func (_RWAGovernance *RWAGovernanceTransactorSession) DelegateVote(assetId *big.Int, delegate common.Address) (*types.Transaction, error) {
	return _RWAGovernance.Contract.DelegateVote(&_RWAGovernance.TransactOpts, assetId, delegate)
}

// ExecuteProposal is a paid mutator transaction binding the contract method 0x0d61b519.
//
// Solidity: function executeProposal(uint256 proposalId) returns()
func (_RWAGovernance *RWAGovernanceTransactor) ExecuteProposal(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _RWAGovernance.contract.Transact(opts, "executeProposal", proposalId)
}

// ExecuteProposal is a paid mutator transaction binding the contract method 0x0d61b519.
//
// Solidity: function executeProposal(uint256 proposalId) returns()
func (_RWAGovernance *RWAGovernanceSession) ExecuteProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _RWAGovernance.Contract.ExecuteProposal(&_RWAGovernance.TransactOpts, proposalId)
}

// ExecuteProposal is a paid mutator transaction binding the contract method 0x0d61b519.
//
// Solidity: function executeProposal(uint256 proposalId) returns()
func (_RWAGovernance *RWAGovernanceTransactorSession) ExecuteProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _RWAGovernance.Contract.ExecuteProposal(&_RWAGovernance.TransactOpts, proposalId)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAGovernance *RWAGovernanceTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAGovernance.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAGovernance *RWAGovernanceSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAGovernance.Contract.GrantRole(&_RWAGovernance.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_RWAGovernance *RWAGovernanceTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAGovernance.Contract.GrantRole(&_RWAGovernance.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAGovernance *RWAGovernanceTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAGovernance.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAGovernance *RWAGovernanceSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAGovernance.Contract.RenounceRole(&_RWAGovernance.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_RWAGovernance *RWAGovernanceTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _RWAGovernance.Contract.RenounceRole(&_RWAGovernance.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAGovernance *RWAGovernanceTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAGovernance.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAGovernance *RWAGovernanceSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAGovernance.Contract.RevokeRole(&_RWAGovernance.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_RWAGovernance *RWAGovernanceTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _RWAGovernance.Contract.RevokeRole(&_RWAGovernance.TransactOpts, role, account)
}

// SetGovernanceParams is a paid mutator transaction binding the contract method 0x833a1f7e.
//
// Solidity: function setGovernanceParams(uint256 assetId, (uint256,uint256,uint256,uint256,uint256,bool) params) returns()
func (_RWAGovernance *RWAGovernanceTransactor) SetGovernanceParams(opts *bind.TransactOpts, assetId *big.Int, params RWAGovernanceGovernanceParams) (*types.Transaction, error) {
	return _RWAGovernance.contract.Transact(opts, "setGovernanceParams", assetId, params)
}

// SetGovernanceParams is a paid mutator transaction binding the contract method 0x833a1f7e.
//
// Solidity: function setGovernanceParams(uint256 assetId, (uint256,uint256,uint256,uint256,uint256,bool) params) returns()
func (_RWAGovernance *RWAGovernanceSession) SetGovernanceParams(assetId *big.Int, params RWAGovernanceGovernanceParams) (*types.Transaction, error) {
	return _RWAGovernance.Contract.SetGovernanceParams(&_RWAGovernance.TransactOpts, assetId, params)
}

// SetGovernanceParams is a paid mutator transaction binding the contract method 0x833a1f7e.
//
// Solidity: function setGovernanceParams(uint256 assetId, (uint256,uint256,uint256,uint256,uint256,bool) params) returns()
func (_RWAGovernance *RWAGovernanceTransactorSession) SetGovernanceParams(assetId *big.Int, params RWAGovernanceGovernanceParams) (*types.Transaction, error) {
	return _RWAGovernance.Contract.SetGovernanceParams(&_RWAGovernance.TransactOpts, assetId, params)
}

// UndelegateVote is a paid mutator transaction binding the contract method 0x9d02f577.
//
// Solidity: function undelegateVote(uint256 assetId) returns()
func (_RWAGovernance *RWAGovernanceTransactor) UndelegateVote(opts *bind.TransactOpts, assetId *big.Int) (*types.Transaction, error) {
	return _RWAGovernance.contract.Transact(opts, "undelegateVote", assetId)
}

// UndelegateVote is a paid mutator transaction binding the contract method 0x9d02f577.
//
// Solidity: function undelegateVote(uint256 assetId) returns()
func (_RWAGovernance *RWAGovernanceSession) UndelegateVote(assetId *big.Int) (*types.Transaction, error) {
	return _RWAGovernance.Contract.UndelegateVote(&_RWAGovernance.TransactOpts, assetId)
}

// UndelegateVote is a paid mutator transaction binding the contract method 0x9d02f577.
//
// Solidity: function undelegateVote(uint256 assetId) returns()
func (_RWAGovernance *RWAGovernanceTransactorSession) UndelegateVote(assetId *big.Int) (*types.Transaction, error) {
	return _RWAGovernance.Contract.UndelegateVote(&_RWAGovernance.TransactOpts, assetId)
}

// VetoProposal is a paid mutator transaction binding the contract method 0x6f65108c.
//
// Solidity: function vetoProposal(uint256 proposalId) returns()
func (_RWAGovernance *RWAGovernanceTransactor) VetoProposal(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _RWAGovernance.contract.Transact(opts, "vetoProposal", proposalId)
}

// VetoProposal is a paid mutator transaction binding the contract method 0x6f65108c.
//
// Solidity: function vetoProposal(uint256 proposalId) returns()
func (_RWAGovernance *RWAGovernanceSession) VetoProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _RWAGovernance.Contract.VetoProposal(&_RWAGovernance.TransactOpts, proposalId)
}

// VetoProposal is a paid mutator transaction binding the contract method 0x6f65108c.
//
// Solidity: function vetoProposal(uint256 proposalId) returns()
func (_RWAGovernance *RWAGovernanceTransactorSession) VetoProposal(proposalId *big.Int) (*types.Transaction, error) {
	return _RWAGovernance.Contract.VetoProposal(&_RWAGovernance.TransactOpts, proposalId)
}

// RWAGovernanceProposalCancelledIterator is returned from FilterProposalCancelled and is used to iterate over the raw logs and unpacked data for ProposalCancelled events raised by the RWAGovernance contract.
type RWAGovernanceProposalCancelledIterator struct {
	Event *RWAGovernanceProposalCancelled // Event containing the contract specifics and raw log

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
func (it *RWAGovernanceProposalCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAGovernanceProposalCancelled)
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
		it.Event = new(RWAGovernanceProposalCancelled)
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
func (it *RWAGovernanceProposalCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAGovernanceProposalCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAGovernanceProposalCancelled represents a ProposalCancelled event raised by the RWAGovernance contract.
type RWAGovernanceProposalCancelled struct {
	ProposalId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalCancelled is a free log retrieval operation binding the contract event 0x416e669c63d9a3a5e36ee7cc7e2104b8db28ccd286aa18966e98fa230c73b08c.
//
// Solidity: event ProposalCancelled(uint256 indexed proposalId)
func (_RWAGovernance *RWAGovernanceFilterer) FilterProposalCancelled(opts *bind.FilterOpts, proposalId []*big.Int) (*RWAGovernanceProposalCancelledIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _RWAGovernance.contract.FilterLogs(opts, "ProposalCancelled", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAGovernanceProposalCancelledIterator{contract: _RWAGovernance.contract, event: "ProposalCancelled", logs: logs, sub: sub}, nil
}

// WatchProposalCancelled is a free log subscription operation binding the contract event 0x416e669c63d9a3a5e36ee7cc7e2104b8db28ccd286aa18966e98fa230c73b08c.
//
// Solidity: event ProposalCancelled(uint256 indexed proposalId)
func (_RWAGovernance *RWAGovernanceFilterer) WatchProposalCancelled(opts *bind.WatchOpts, sink chan<- *RWAGovernanceProposalCancelled, proposalId []*big.Int) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _RWAGovernance.contract.WatchLogs(opts, "ProposalCancelled", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAGovernanceProposalCancelled)
				if err := _RWAGovernance.contract.UnpackLog(event, "ProposalCancelled", log); err != nil {
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

// ParseProposalCancelled is a log parse operation binding the contract event 0x416e669c63d9a3a5e36ee7cc7e2104b8db28ccd286aa18966e98fa230c73b08c.
//
// Solidity: event ProposalCancelled(uint256 indexed proposalId)
func (_RWAGovernance *RWAGovernanceFilterer) ParseProposalCancelled(log types.Log) (*RWAGovernanceProposalCancelled, error) {
	event := new(RWAGovernanceProposalCancelled)
	if err := _RWAGovernance.contract.UnpackLog(event, "ProposalCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAGovernanceProposalCreatedIterator is returned from FilterProposalCreated and is used to iterate over the raw logs and unpacked data for ProposalCreated events raised by the RWAGovernance contract.
type RWAGovernanceProposalCreatedIterator struct {
	Event *RWAGovernanceProposalCreated // Event containing the contract specifics and raw log

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
func (it *RWAGovernanceProposalCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAGovernanceProposalCreated)
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
		it.Event = new(RWAGovernanceProposalCreated)
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
func (it *RWAGovernanceProposalCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAGovernanceProposalCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAGovernanceProposalCreated represents a ProposalCreated event raised by the RWAGovernance contract.
type RWAGovernanceProposalCreated struct {
	ProposalId   *big.Int
	AssetId      *big.Int
	Proposer     common.Address
	ProposalType uint8
	Description  string
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterProposalCreated is a free log retrieval operation binding the contract event 0xb1d569dc5621a0a51229b7ac0eff81f3e6df70a3f72fe608aa58f74bb69d7bb6.
//
// Solidity: event ProposalCreated(uint256 indexed proposalId, uint256 indexed assetId, address indexed proposer, uint8 proposalType, string description)
func (_RWAGovernance *RWAGovernanceFilterer) FilterProposalCreated(opts *bind.FilterOpts, proposalId []*big.Int, assetId []*big.Int, proposer []common.Address) (*RWAGovernanceProposalCreatedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}
	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _RWAGovernance.contract.FilterLogs(opts, "ProposalCreated", proposalIdRule, assetIdRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return &RWAGovernanceProposalCreatedIterator{contract: _RWAGovernance.contract, event: "ProposalCreated", logs: logs, sub: sub}, nil
}

// WatchProposalCreated is a free log subscription operation binding the contract event 0xb1d569dc5621a0a51229b7ac0eff81f3e6df70a3f72fe608aa58f74bb69d7bb6.
//
// Solidity: event ProposalCreated(uint256 indexed proposalId, uint256 indexed assetId, address indexed proposer, uint8 proposalType, string description)
func (_RWAGovernance *RWAGovernanceFilterer) WatchProposalCreated(opts *bind.WatchOpts, sink chan<- *RWAGovernanceProposalCreated, proposalId []*big.Int, assetId []*big.Int, proposer []common.Address) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}
	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _RWAGovernance.contract.WatchLogs(opts, "ProposalCreated", proposalIdRule, assetIdRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAGovernanceProposalCreated)
				if err := _RWAGovernance.contract.UnpackLog(event, "ProposalCreated", log); err != nil {
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

// ParseProposalCreated is a log parse operation binding the contract event 0xb1d569dc5621a0a51229b7ac0eff81f3e6df70a3f72fe608aa58f74bb69d7bb6.
//
// Solidity: event ProposalCreated(uint256 indexed proposalId, uint256 indexed assetId, address indexed proposer, uint8 proposalType, string description)
func (_RWAGovernance *RWAGovernanceFilterer) ParseProposalCreated(log types.Log) (*RWAGovernanceProposalCreated, error) {
	event := new(RWAGovernanceProposalCreated)
	if err := _RWAGovernance.contract.UnpackLog(event, "ProposalCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAGovernanceProposalExecutedIterator is returned from FilterProposalExecuted and is used to iterate over the raw logs and unpacked data for ProposalExecuted events raised by the RWAGovernance contract.
type RWAGovernanceProposalExecutedIterator struct {
	Event *RWAGovernanceProposalExecuted // Event containing the contract specifics and raw log

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
func (it *RWAGovernanceProposalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAGovernanceProposalExecuted)
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
		it.Event = new(RWAGovernanceProposalExecuted)
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
func (it *RWAGovernanceProposalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAGovernanceProposalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAGovernanceProposalExecuted represents a ProposalExecuted event raised by the RWAGovernance contract.
type RWAGovernanceProposalExecuted struct {
	ProposalId *big.Int
	Success    bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalExecuted is a free log retrieval operation binding the contract event 0x948f4a9cd986f1118c3fbd459f7a22b23c0693e1ca3ef06a6a8be5aa7d39cc03.
//
// Solidity: event ProposalExecuted(uint256 indexed proposalId, bool success)
func (_RWAGovernance *RWAGovernanceFilterer) FilterProposalExecuted(opts *bind.FilterOpts, proposalId []*big.Int) (*RWAGovernanceProposalExecutedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _RWAGovernance.contract.FilterLogs(opts, "ProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAGovernanceProposalExecutedIterator{contract: _RWAGovernance.contract, event: "ProposalExecuted", logs: logs, sub: sub}, nil
}

// WatchProposalExecuted is a free log subscription operation binding the contract event 0x948f4a9cd986f1118c3fbd459f7a22b23c0693e1ca3ef06a6a8be5aa7d39cc03.
//
// Solidity: event ProposalExecuted(uint256 indexed proposalId, bool success)
func (_RWAGovernance *RWAGovernanceFilterer) WatchProposalExecuted(opts *bind.WatchOpts, sink chan<- *RWAGovernanceProposalExecuted, proposalId []*big.Int) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _RWAGovernance.contract.WatchLogs(opts, "ProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAGovernanceProposalExecuted)
				if err := _RWAGovernance.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
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

// ParseProposalExecuted is a log parse operation binding the contract event 0x948f4a9cd986f1118c3fbd459f7a22b23c0693e1ca3ef06a6a8be5aa7d39cc03.
//
// Solidity: event ProposalExecuted(uint256 indexed proposalId, bool success)
func (_RWAGovernance *RWAGovernanceFilterer) ParseProposalExecuted(log types.Log) (*RWAGovernanceProposalExecuted, error) {
	event := new(RWAGovernanceProposalExecuted)
	if err := _RWAGovernance.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAGovernanceProposalVetoedIterator is returned from FilterProposalVetoed and is used to iterate over the raw logs and unpacked data for ProposalVetoed events raised by the RWAGovernance contract.
type RWAGovernanceProposalVetoedIterator struct {
	Event *RWAGovernanceProposalVetoed // Event containing the contract specifics and raw log

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
func (it *RWAGovernanceProposalVetoedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAGovernanceProposalVetoed)
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
		it.Event = new(RWAGovernanceProposalVetoed)
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
func (it *RWAGovernanceProposalVetoedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAGovernanceProposalVetoedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAGovernanceProposalVetoed represents a ProposalVetoed event raised by the RWAGovernance contract.
type RWAGovernanceProposalVetoed struct {
	ProposalId *big.Int
	Vetoer     common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalVetoed is a free log retrieval operation binding the contract event 0xc2131db10d833d6b93ae553fa450ea0f6c0c2dfb0160cb3309468bf72718f3eb.
//
// Solidity: event ProposalVetoed(uint256 indexed proposalId, address indexed vetoer)
func (_RWAGovernance *RWAGovernanceFilterer) FilterProposalVetoed(opts *bind.FilterOpts, proposalId []*big.Int, vetoer []common.Address) (*RWAGovernanceProposalVetoedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}
	var vetoerRule []interface{}
	for _, vetoerItem := range vetoer {
		vetoerRule = append(vetoerRule, vetoerItem)
	}

	logs, sub, err := _RWAGovernance.contract.FilterLogs(opts, "ProposalVetoed", proposalIdRule, vetoerRule)
	if err != nil {
		return nil, err
	}
	return &RWAGovernanceProposalVetoedIterator{contract: _RWAGovernance.contract, event: "ProposalVetoed", logs: logs, sub: sub}, nil
}

// WatchProposalVetoed is a free log subscription operation binding the contract event 0xc2131db10d833d6b93ae553fa450ea0f6c0c2dfb0160cb3309468bf72718f3eb.
//
// Solidity: event ProposalVetoed(uint256 indexed proposalId, address indexed vetoer)
func (_RWAGovernance *RWAGovernanceFilterer) WatchProposalVetoed(opts *bind.WatchOpts, sink chan<- *RWAGovernanceProposalVetoed, proposalId []*big.Int, vetoer []common.Address) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}
	var vetoerRule []interface{}
	for _, vetoerItem := range vetoer {
		vetoerRule = append(vetoerRule, vetoerItem)
	}

	logs, sub, err := _RWAGovernance.contract.WatchLogs(opts, "ProposalVetoed", proposalIdRule, vetoerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAGovernanceProposalVetoed)
				if err := _RWAGovernance.contract.UnpackLog(event, "ProposalVetoed", log); err != nil {
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

// ParseProposalVetoed is a log parse operation binding the contract event 0xc2131db10d833d6b93ae553fa450ea0f6c0c2dfb0160cb3309468bf72718f3eb.
//
// Solidity: event ProposalVetoed(uint256 indexed proposalId, address indexed vetoer)
func (_RWAGovernance *RWAGovernanceFilterer) ParseProposalVetoed(log types.Log) (*RWAGovernanceProposalVetoed, error) {
	event := new(RWAGovernanceProposalVetoed)
	if err := _RWAGovernance.contract.UnpackLog(event, "ProposalVetoed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAGovernanceRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the RWAGovernance contract.
type RWAGovernanceRoleAdminChangedIterator struct {
	Event *RWAGovernanceRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *RWAGovernanceRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAGovernanceRoleAdminChanged)
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
		it.Event = new(RWAGovernanceRoleAdminChanged)
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
func (it *RWAGovernanceRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAGovernanceRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAGovernanceRoleAdminChanged represents a RoleAdminChanged event raised by the RWAGovernance contract.
type RWAGovernanceRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RWAGovernance *RWAGovernanceFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*RWAGovernanceRoleAdminChangedIterator, error) {

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

	logs, sub, err := _RWAGovernance.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &RWAGovernanceRoleAdminChangedIterator{contract: _RWAGovernance.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_RWAGovernance *RWAGovernanceFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *RWAGovernanceRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _RWAGovernance.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAGovernanceRoleAdminChanged)
				if err := _RWAGovernance.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_RWAGovernance *RWAGovernanceFilterer) ParseRoleAdminChanged(log types.Log) (*RWAGovernanceRoleAdminChanged, error) {
	event := new(RWAGovernanceRoleAdminChanged)
	if err := _RWAGovernance.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAGovernanceRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the RWAGovernance contract.
type RWAGovernanceRoleGrantedIterator struct {
	Event *RWAGovernanceRoleGranted // Event containing the contract specifics and raw log

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
func (it *RWAGovernanceRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAGovernanceRoleGranted)
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
		it.Event = new(RWAGovernanceRoleGranted)
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
func (it *RWAGovernanceRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAGovernanceRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAGovernanceRoleGranted represents a RoleGranted event raised by the RWAGovernance contract.
type RWAGovernanceRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAGovernance *RWAGovernanceFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RWAGovernanceRoleGrantedIterator, error) {

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

	logs, sub, err := _RWAGovernance.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RWAGovernanceRoleGrantedIterator{contract: _RWAGovernance.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAGovernance *RWAGovernanceFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *RWAGovernanceRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RWAGovernance.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAGovernanceRoleGranted)
				if err := _RWAGovernance.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_RWAGovernance *RWAGovernanceFilterer) ParseRoleGranted(log types.Log) (*RWAGovernanceRoleGranted, error) {
	event := new(RWAGovernanceRoleGranted)
	if err := _RWAGovernance.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAGovernanceRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the RWAGovernance contract.
type RWAGovernanceRoleRevokedIterator struct {
	Event *RWAGovernanceRoleRevoked // Event containing the contract specifics and raw log

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
func (it *RWAGovernanceRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAGovernanceRoleRevoked)
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
		it.Event = new(RWAGovernanceRoleRevoked)
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
func (it *RWAGovernanceRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAGovernanceRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAGovernanceRoleRevoked represents a RoleRevoked event raised by the RWAGovernance contract.
type RWAGovernanceRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAGovernance *RWAGovernanceFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*RWAGovernanceRoleRevokedIterator, error) {

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

	logs, sub, err := _RWAGovernance.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &RWAGovernanceRoleRevokedIterator{contract: _RWAGovernance.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_RWAGovernance *RWAGovernanceFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *RWAGovernanceRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _RWAGovernance.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAGovernanceRoleRevoked)
				if err := _RWAGovernance.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_RWAGovernance *RWAGovernanceFilterer) ParseRoleRevoked(log types.Log) (*RWAGovernanceRoleRevoked, error) {
	event := new(RWAGovernanceRoleRevoked)
	if err := _RWAGovernance.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAGovernanceVoteCastIterator is returned from FilterVoteCast and is used to iterate over the raw logs and unpacked data for VoteCast events raised by the RWAGovernance contract.
type RWAGovernanceVoteCastIterator struct {
	Event *RWAGovernanceVoteCast // Event containing the contract specifics and raw log

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
func (it *RWAGovernanceVoteCastIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAGovernanceVoteCast)
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
		it.Event = new(RWAGovernanceVoteCast)
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
func (it *RWAGovernanceVoteCastIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAGovernanceVoteCastIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAGovernanceVoteCast represents a VoteCast event raised by the RWAGovernance contract.
type RWAGovernanceVoteCast struct {
	Voter      common.Address
	ProposalId *big.Int
	Choice     uint8
	Votes      *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVoteCast is a free log retrieval operation binding the contract event 0x2c9deb38f462962eadbd85a9d3a4120503ee091f1582eaaa10aa8c6797651d29.
//
// Solidity: event VoteCast(address indexed voter, uint256 indexed proposalId, uint8 choice, uint256 votes)
func (_RWAGovernance *RWAGovernanceFilterer) FilterVoteCast(opts *bind.FilterOpts, voter []common.Address, proposalId []*big.Int) (*RWAGovernanceVoteCastIterator, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}
	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _RWAGovernance.contract.FilterLogs(opts, "VoteCast", voterRule, proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAGovernanceVoteCastIterator{contract: _RWAGovernance.contract, event: "VoteCast", logs: logs, sub: sub}, nil
}

// WatchVoteCast is a free log subscription operation binding the contract event 0x2c9deb38f462962eadbd85a9d3a4120503ee091f1582eaaa10aa8c6797651d29.
//
// Solidity: event VoteCast(address indexed voter, uint256 indexed proposalId, uint8 choice, uint256 votes)
func (_RWAGovernance *RWAGovernanceFilterer) WatchVoteCast(opts *bind.WatchOpts, sink chan<- *RWAGovernanceVoteCast, voter []common.Address, proposalId []*big.Int) (event.Subscription, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}
	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _RWAGovernance.contract.WatchLogs(opts, "VoteCast", voterRule, proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAGovernanceVoteCast)
				if err := _RWAGovernance.contract.UnpackLog(event, "VoteCast", log); err != nil {
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

// ParseVoteCast is a log parse operation binding the contract event 0x2c9deb38f462962eadbd85a9d3a4120503ee091f1582eaaa10aa8c6797651d29.
//
// Solidity: event VoteCast(address indexed voter, uint256 indexed proposalId, uint8 choice, uint256 votes)
func (_RWAGovernance *RWAGovernanceFilterer) ParseVoteCast(log types.Log) (*RWAGovernanceVoteCast, error) {
	event := new(RWAGovernanceVoteCast)
	if err := _RWAGovernance.contract.UnpackLog(event, "VoteCast", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RWAGovernanceVoteDelegatedIterator is returned from FilterVoteDelegated and is used to iterate over the raw logs and unpacked data for VoteDelegated events raised by the RWAGovernance contract.
type RWAGovernanceVoteDelegatedIterator struct {
	Event *RWAGovernanceVoteDelegated // Event containing the contract specifics and raw log

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
func (it *RWAGovernanceVoteDelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RWAGovernanceVoteDelegated)
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
		it.Event = new(RWAGovernanceVoteDelegated)
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
func (it *RWAGovernanceVoteDelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RWAGovernanceVoteDelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RWAGovernanceVoteDelegated represents a VoteDelegated event raised by the RWAGovernance contract.
type RWAGovernanceVoteDelegated struct {
	Delegator common.Address
	Delegate  common.Address
	AssetId   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVoteDelegated is a free log retrieval operation binding the contract event 0xf5c2b450c3b6c4c3cc8c45e612ffa5ce51cda974e8b2ecaf58abda86d8965847.
//
// Solidity: event VoteDelegated(address indexed delegator, address indexed delegate, uint256 indexed assetId)
func (_RWAGovernance *RWAGovernanceFilterer) FilterVoteDelegated(opts *bind.FilterOpts, delegator []common.Address, delegate []common.Address, assetId []*big.Int) (*RWAGovernanceVoteDelegatedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var delegateRule []interface{}
	for _, delegateItem := range delegate {
		delegateRule = append(delegateRule, delegateItem)
	}
	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWAGovernance.contract.FilterLogs(opts, "VoteDelegated", delegatorRule, delegateRule, assetIdRule)
	if err != nil {
		return nil, err
	}
	return &RWAGovernanceVoteDelegatedIterator{contract: _RWAGovernance.contract, event: "VoteDelegated", logs: logs, sub: sub}, nil
}

// WatchVoteDelegated is a free log subscription operation binding the contract event 0xf5c2b450c3b6c4c3cc8c45e612ffa5ce51cda974e8b2ecaf58abda86d8965847.
//
// Solidity: event VoteDelegated(address indexed delegator, address indexed delegate, uint256 indexed assetId)
func (_RWAGovernance *RWAGovernanceFilterer) WatchVoteDelegated(opts *bind.WatchOpts, sink chan<- *RWAGovernanceVoteDelegated, delegator []common.Address, delegate []common.Address, assetId []*big.Int) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var delegateRule []interface{}
	for _, delegateItem := range delegate {
		delegateRule = append(delegateRule, delegateItem)
	}
	var assetIdRule []interface{}
	for _, assetIdItem := range assetId {
		assetIdRule = append(assetIdRule, assetIdItem)
	}

	logs, sub, err := _RWAGovernance.contract.WatchLogs(opts, "VoteDelegated", delegatorRule, delegateRule, assetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RWAGovernanceVoteDelegated)
				if err := _RWAGovernance.contract.UnpackLog(event, "VoteDelegated", log); err != nil {
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

// ParseVoteDelegated is a log parse operation binding the contract event 0xf5c2b450c3b6c4c3cc8c45e612ffa5ce51cda974e8b2ecaf58abda86d8965847.
//
// Solidity: event VoteDelegated(address indexed delegator, address indexed delegate, uint256 indexed assetId)
func (_RWAGovernance *RWAGovernanceFilterer) ParseVoteDelegated(log types.Log) (*RWAGovernanceVoteDelegated, error) {
	event := new(RWAGovernanceVoteDelegated)
	if err := _RWAGovernance.contract.UnpackLog(event, "VoteDelegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
