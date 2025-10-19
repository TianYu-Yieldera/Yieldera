package middleware

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	// JWT secret key - should be loaded from environment in production
	jwtSecret = []byte("your-secret-key-change-in-production")

	// Token expiration time
	tokenExpiration = 24 * time.Hour
)

// Claims represents the JWT claims
type Claims struct {
	Address string `json:"address"`
	jwt.RegisteredClaims
}

// AuthRequest represents the authentication request body
type AuthRequest struct {
	Address   string `json:"address" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

// WalletAuthMiddleware validates the JWT token from the Authorization header
func WalletAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set address in context for downstream handlers
		c.Set("userAddress", claims.Address)
		c.Next()
	}
}

// VerifySignature verifies an Ethereum signature
func VerifySignature(address, message, signature string) error {
	// Normalize address
	address = strings.ToLower(address)
	if !strings.HasPrefix(address, "0x") {
		return errors.New("address must start with 0x")
	}

	// Decode signature
	sigBytes, err := hexutil.Decode(signature)
	if err != nil {
		return fmt.Errorf("invalid signature format: %v", err)
	}

	if len(sigBytes) != 65 {
		return errors.New("signature must be 65 bytes")
	}

	// Adjust V value (Ethereum uses 27/28, we need 0/1)
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	// Create message hash (EIP-191)
	messageHash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))

	// Recover public key
	pubKey, err := crypto.SigToPub(messageHash.Bytes(), sigBytes)
	if err != nil {
		return fmt.Errorf("failed to recover public key: %v", err)
	}

	// Get address from public key
	recoveredAddr := crypto.PubkeyToAddress(*pubKey).Hex()
	recoveredAddr = strings.ToLower(recoveredAddr)

	// Compare addresses
	if recoveredAddr != address {
		return fmt.Errorf("signature verification failed: expected %s, got %s", address, recoveredAddr)
	}

	return nil
}

// GenerateToken generates a JWT token for the given address
func GenerateToken(address string) (string, error) {
	expirationTime := time.Now().Add(tokenExpiration)

	claims := &Claims{
		Address: address,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "loyalty-defi",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// AuthenticateHandler handles the wallet authentication request
func AuthenticateHandler(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Verify the signature
	if err := VerifySignature(req.Address, req.Message, req.Signature); err != nil {
		log.Printf("Signature verification failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// Generate JWT token
	token, err := GenerateToken(req.Address)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"address":   req.Address,
		"expiresIn": int(tokenExpiration.Seconds()),
	})
}

// GetAuthMessage generates a message for the user to sign
func GetAuthMessageHandler(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address parameter required"})
		return
	}

	// Validate address format
	if len(address) != 42 || !strings.HasPrefix(address, "0x") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Ethereum address"})
		return
	}

	// Generate a message with timestamp to prevent replay attacks
	timestamp := time.Now().Unix()
	message := fmt.Sprintf("Sign this message to authenticate with PointFi Protocol.\n\nAddress: %s\nTimestamp: %d", address, timestamp)

	c.JSON(http.StatusOK, gin.H{
		"message":   message,
		"timestamp": timestamp,
	})
}

// Helper function to get address from context
func GetUserAddress(c *gin.Context) (string, bool) {
	address, exists := c.Get("userAddress")
	if !exists {
		return "", false
	}
	addrStr, ok := address.(string)
	return addrStr, ok
}

// Optional: Helper to verify ECDSA signature (alternative implementation)
func verifyECDSASignature(message, signature string, publicKey *ecdsa.PublicKey) bool {
	messageHash := crypto.Keccak256Hash([]byte(message))
	sigBytes, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	if err != nil {
		return false
	}

	if len(sigBytes) != 65 {
		return false
	}

	r := sigBytes[:32]
	s := sigBytes[32:64]

	return ecdsa.Verify(publicKey, messageHash.Bytes(), new(big.Int).SetBytes(r), new(big.Int).SetBytes(s))
}
