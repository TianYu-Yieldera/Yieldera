package blockchain

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

// EIP712Domain represents the EIP-712 domain separator
type EIP712Domain struct {
	Name              string
	Version           string
	ChainID           *big.Int
	VerifyingContract common.Address
}

// HashTypedData hashes data according to EIP-712
func HashTypedData(domainSeparator []byte, structHash []byte) []byte {
	message := []byte("\x19\x01")
	message = append(message, domainSeparator...)
	message = append(message, structHash...)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(message)
	return hash.Sum(nil)
}

// HashStruct hashes a struct according to EIP-712
func HashStruct(typeHash []byte, data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(typeHash)
	hash.Write(data)
	return hash.Sum(nil)
}

// Hash domain separator
func HashDomain(domain EIP712Domain) []byte {
	// EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)
	typeHash := crypto.Keccak256([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))

	hash := sha3.NewLegacyKeccak256()
	hash.Write(typeHash)
	hash.Write(crypto.Keccak256([]byte(domain.Name)))
	hash.Write(crypto.Keccak256([]byte(domain.Version)))
	hash.Write(common.LeftPadBytes(domain.ChainID.Bytes(), 32))
	hash.Write(common.LeftPadBytes(domain.VerifyingContract.Bytes(), 32))

	return hash.Sum(nil)
}

// VerifyEIP712Signature verifies an EIP-712 signature
func VerifyEIP712Signature(signature string, messageHash []byte, expectedAddress string) (bool, error) {
	// Clean up signature
	signature = strings.TrimPrefix(signature, "0x")

	// Decode hex signature
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %w", err)
	}

	if len(sigBytes) != 65 {
		return false, fmt.Errorf("invalid signature length: %d", len(sigBytes))
	}

	// Transform yellow paper V from 27/28 to 0/1
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	// Recover public key
	pubKey, err := crypto.SigToPub(messageHash, sigBytes)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	// Get address from public key
	recoveredAddress := crypto.PubkeyToAddress(*pubKey)

	// Parse expected address
	expected := common.HexToAddress(expectedAddress)

	// Compare addresses
	return recoveredAddress == expected, nil
}

// PadBytes32 pads data to 32 bytes
func PadBytes32(data []byte) []byte {
	return common.LeftPadBytes(data, 32)
}

// Keccak256 calculates keccak256 hash
func Keccak256(data ...[]byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	for _, d := range data {
		hash.Write(d)
	}
	return hash.Sum(nil)
}
