package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the loyalty points system
type Config struct {
	// Database
	DatabaseURL string

	// Kafka
	KafkaBrokers     string
	KafkaTopicRaw    string // Legacy topic (for backward compatibility)
	KafkaTopicL1     string
	KafkaTopicL2     string
	KafkaTopicBridge string

	// L1 Configuration (Ethereum)
	L1ChainID       int64
	L1RPCURL        string
	L1WSSURL        string
	L1Confirmations int

	// L2 Configuration (Arbitrum)
	L2ChainID       int64
	L2RPCURL        string
	L2WSSURL        string
	L2Confirmations int

	// L1 Contract Addresses
	L1CollateralVault string
	L1StateRegistry   string
	L1LoyaltyUSD      string
	L1Gateway         string

	// L2 Core Contract Addresses
	L2IntegratedVault string
	L2StateAggregator string

	// L2 DeFi Adapter Addresses
	L2AaveAdapter     string
	L2CompoundAdapter string
	L2UniswapAdapter  string

	// L2 RWA Contract Addresses
	L2RWAFactory          string
	L2RWAMarketplace      string
	L2RWAYieldDistributor string
	L2RWACompliance       string
	L2RWAValuation        string
	L2RWAGovernance       string

	// Arbitrum Bridge Addresses
	ArbitrumInbox  string
	ArbitrumOutbox string

	// API Configuration
	APIPort        string
	APIAllowOrigin string

	// Service URLs
	VaultServiceURL  string
	RWAServiceURL    string
	OracleServiceURL string

	// Points System Configuration
	PointsRate           float64
	SchedulerIntervalSec int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		// Database
		DatabaseURL: os.Getenv("DATABASE_URL"),

		// Kafka
		KafkaBrokers:     getEnvOrDefault("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopicRaw:    getEnvOrDefault("KAFKA_TOPIC_RAW", "events.raw"),
		KafkaTopicL1:     getEnvOrDefault("KAFKA_TOPIC_L1", "events.l1"),
		KafkaTopicL2:     getEnvOrDefault("KAFKA_TOPIC_L2", "events.l2"),
		KafkaTopicBridge: getEnvOrDefault("KAFKA_TOPIC_BRIDGE", "events.bridge"),

		// L1 Configuration
		L1ChainID:       getEnvInt64("L1_CHAIN_ID", 11155111), // Sepolia by default
		L1RPCURL:        os.Getenv("L1_RPC_URL"),
		L1WSSURL:        os.Getenv("L1_WSS_URL"),
		L1Confirmations: getEnvInt("L1_CONFIRMATIONS", 12),

		// L2 Configuration
		L2ChainID:       getEnvInt64("L2_CHAIN_ID", 421614), // Arbitrum Sepolia by default
		L2RPCURL:        os.Getenv("L2_RPC_URL"),
		L2WSSURL:        os.Getenv("L2_WSS_URL"),
		L2Confirmations: getEnvInt("L2_CONFIRMATIONS", 1),

		// L1 Contract Addresses
		L1CollateralVault: os.Getenv("L1_COLLATERAL_VAULT"),
		L1StateRegistry:   os.Getenv("L1_STATE_REGISTRY"),
		L1LoyaltyUSD:      os.Getenv("L1_LOYALTY_USD"),
		L1Gateway:         os.Getenv("L1_GATEWAY"),

		// L2 Core Contract Addresses
		L2IntegratedVault: os.Getenv("L2_INTEGRATED_VAULT"),
		L2StateAggregator: os.Getenv("L2_STATE_AGGREGATOR"),

		// L2 DeFi Adapter Addresses
		L2AaveAdapter:     os.Getenv("L2_AAVE_ADAPTER"),
		L2CompoundAdapter: os.Getenv("L2_COMPOUND_ADAPTER"),
		L2UniswapAdapter:  os.Getenv("L2_UNISWAP_ADAPTER"),

		// L2 RWA Contract Addresses
		L2RWAFactory:          os.Getenv("L2_RWA_FACTORY"),
		L2RWAMarketplace:      os.Getenv("L2_RWA_MARKETPLACE"),
		L2RWAYieldDistributor: os.Getenv("L2_RWA_YIELD_DISTRIBUTOR"),
		L2RWACompliance:       os.Getenv("L2_RWA_COMPLIANCE"),
		L2RWAValuation:        os.Getenv("L2_RWA_VALUATION"),
		L2RWAGovernance:       os.Getenv("L2_RWA_GOVERNANCE"),

		// Arbitrum Bridge
		ArbitrumInbox:  getEnvOrDefault("ARBITRUM_SEPOLIA_INBOX", "0xaAe29B0366299461418F5324a79Afc425BE5ae21"),
		ArbitrumOutbox: getEnvOrDefault("ARBITRUM_SEPOLIA_OUTBOX", "0x65f07C7D521164a4d5DaC6eB8Fac8DA067A3B78F"),

		// API Configuration
		APIPort:        getEnvOrDefault("API_PORT", "8080"),
		APIAllowOrigin: getEnvOrDefault("API_ALLOW_ORIGIN", "*"),

		// Service URLs
		VaultServiceURL:  getEnvOrDefault("VAULT_SERVICE_URL", "http://vault:8081"),
		RWAServiceURL:    getEnvOrDefault("RWA_SERVICE_URL", "http://rwa:8082"),
		OracleServiceURL: getEnvOrDefault("ORACLE_SERVICE_URL", "http://oracle:8083"),

		// Points System
		PointsRate:           getEnvFloat64("POINTS_RATE", 0.05),
		SchedulerIntervalSec: getEnvInt("SCHEDULER_INTERVAL_SEC", 60),
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if all required configuration values are set
func (c *Config) Validate() error {
	// Check database URL
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	// Check L1 RPC
	if c.L1RPCURL == "" {
		return fmt.Errorf("L1_RPC_URL is required")
	}

	// Check L2 RPC
	if c.L2RPCURL == "" {
		return fmt.Errorf("L2_RPC_URL is required")
	}

	// Check Kafka brokers
	if c.KafkaBrokers == "" {
		return fmt.Errorf("KAFKA_BROKERS is required")
	}

	return nil
}

// IsL1ContractConfigured checks if an L1 contract address is configured
func (c *Config) IsL1ContractConfigured() bool {
	return c.L1CollateralVault != "" &&
		c.L1StateRegistry != "" &&
		c.L1LoyaltyUSD != "" &&
		c.L1Gateway != ""
}

// IsL2ContractConfigured checks if all L2 contracts are configured
func (c *Config) IsL2ContractConfigured() bool {
	return c.L2IntegratedVault != "" &&
		c.L2StateAggregator != "" &&
		c.L2AaveAdapter != "" &&
		c.L2CompoundAdapter != "" &&
		c.L2UniswapAdapter != "" &&
		c.L2RWAFactory != "" &&
		c.L2RWAMarketplace != "" &&
		c.L2RWAYieldDistributor != "" &&
		c.L2RWACompliance != "" &&
		c.L2RWAValuation != "" &&
		c.L2RWAGovernance != ""
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.L1ChainID == 1 && c.L2ChainID == 42161 // Ethereum mainnet + Arbitrum One
}

// IsTestnet returns true if running in testnet mode
func (c *Config) IsTestnet() bool {
	return c.L1ChainID == 11155111 && c.L2ChainID == 421614 // Sepolia + Arbitrum Sepolia
}

// Helper functions

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvInt64(key string, defaultVal int64) int64 {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvFloat64(key string, defaultVal float64) float64 {
	if val := os.Getenv(key); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return defaultVal
}
