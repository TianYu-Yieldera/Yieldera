package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Event processing metrics
	EventsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loyalty_events_processed_total",
			Help: "Total number of blockchain events processed",
		},
		[]string{"layer", "event_type", "status"},
	)

	EventProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loyalty_event_processing_duration_seconds",
			Help:    "Event processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"layer", "event_type"},
	)

	// Kafka metrics
	KafkaLag = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "loyalty_kafka_consumer_lag",
			Help: "Current Kafka consumer lag",
		},
		[]string{"topic", "partition"},
	)

	KafkaMessagesConsumed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loyalty_kafka_messages_consumed_total",
			Help: "Total number of Kafka messages consumed",
		},
		[]string{"topic"},
	)

	// Bridge metrics
	BridgeMessagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loyalty_bridge_messages_total",
			Help: "Total number of bridge messages",
		},
		[]string{"direction", "status"},
	)

	BridgeConfirmationTime = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "loyalty_bridge_confirmation_seconds",
			Help:    "Bridge message confirmation time in seconds",
			Buckets: []float64{60, 300, 600, 900, 1800, 3600},
		},
	)

	BridgePendingMessages = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "loyalty_bridge_pending_messages",
			Help: "Number of pending bridge messages",
		},
	)

	BridgeRetryCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loyalty_bridge_retry_total",
			Help: "Total number of bridge message retries",
		},
		[]string{"direction"},
	)

	// API metrics
	APIRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loyalty_api_requests_total",
			Help: "Total number of API requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	APIRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loyalty_api_request_duration_seconds",
			Help:    "API request duration in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1.0, 2.0, 5.0},
		},
		[]string{"method", "endpoint"},
	)

	// Database metrics
	DBConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "loyalty_db_connections_active",
			Help: "Number of active database connections",
		},
	)

	DBConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "loyalty_db_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loyalty_db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{0.001, 0.01, 0.05, 0.1, 0.5, 1.0},
		},
		[]string{"operation"},
	)

	// Treasury module metrics
	TreasuryAssetsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "loyalty_treasury_assets_total",
			Help: "Total number of treasury assets",
		},
	)

	TreasuryTVL = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "loyalty_treasury_tvl_usd",
			Help: "Total value locked in treasury assets (USD)",
		},
	)

	TreasuryOrdersTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loyalty_treasury_orders_total",
			Help: "Total number of treasury market orders",
		},
		[]string{"order_type", "status"},
	)

	TreasuryYieldDistributed = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "loyalty_treasury_yield_distributed_usd",
			Help: "Total yield distributed to users (USD)",
		},
	)

	// System metrics
	SystemErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loyalty_system_errors_total",
			Help: "Total number of system errors",
		},
		[]string{"service", "error_type"},
	)

	// DeFi adapter metrics
	DeFiAdapterCalls = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loyalty_defi_adapter_calls_total",
			Help: "Total number of DeFi adapter calls",
		},
		[]string{"protocol", "operation", "status"},
	)

	DeFiAdapterDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "loyalty_defi_adapter_duration_seconds",
			Help:    "DeFi adapter call duration in seconds",
			Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0, 10.0},
		},
		[]string{"protocol", "operation"},
	)
)
