#!/usr/bin/env python3
"""
Prometheus Metrics Integration for AI Risk Engine
Exposes metrics for monitoring and observability
"""

from datetime import datetime
from typing import Dict, Optional

import prometheus_client
from prometheus_client import (
    Counter, Gauge, Histogram, Summary, Info,
    CollectorRegistry, generate_latest, CONTENT_TYPE_LATEST,
    push_to_gateway, start_http_server
)
from prometheus_client.core import GaugeMetricFamily, CounterMetricFamily

# ============================================================
# METRIC DEFINITIONS
# ============================================================

# Create custom registry
registry = CollectorRegistry()

# ============================================================
# REQUEST METRICS
# ============================================================

# Request counter
request_count = Counter(
    'ai_risk_api_requests_total',
    'Total number of API requests',
    ['method', 'endpoint', 'status'],
    registry=registry
)

# Request duration histogram
request_duration = Histogram(
    'ai_risk_api_request_duration_seconds',
    'Request duration in seconds',
    ['method', 'endpoint'],
    buckets=(0.001, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0),
    registry=registry
)

# Active requests gauge
active_requests = Gauge(
    'ai_risk_api_active_requests',
    'Number of active requests',
    registry=registry
)

# ============================================================
# RISK METRICS
# ============================================================

# Risk score gauge
risk_score_gauge = Gauge(
    'ai_risk_score',
    'Current overall risk score',
    ['protocol', 'asset'],
    registry=registry
)

# Liquidation counter
liquidation_count = Counter(
    'ai_risk_liquidations_total',
    'Total number of liquidations',
    ['protocol', 'asset', 'severity'],
    registry=registry
)

# Positions at risk gauge
positions_at_risk = Gauge(
    'ai_risk_positions_at_risk',
    'Number of positions at risk',
    ['risk_level'],
    registry=registry
)

# Total value at risk gauge
total_value_at_risk = Gauge(
    'ai_risk_total_value_at_risk_usd',
    'Total value at risk in USD',
    ['protocol'],
    registry=registry
)

# ============================================================
# ML MODEL METRICS
# ============================================================

# Model prediction counter
ml_predictions = Counter(
    'ai_risk_ml_predictions_total',
    'Total number of ML predictions',
    ['model_type', 'prediction_type'],
    registry=registry
)

# Model prediction latency
ml_prediction_latency = Histogram(
    'ai_risk_ml_prediction_duration_seconds',
    'ML prediction latency',
    ['model_type'],
    buckets=(0.01, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0),
    registry=registry
)

# Model accuracy gauge
ml_model_accuracy = Gauge(
    'ai_risk_ml_model_accuracy',
    'ML model accuracy score',
    ['model_type', 'metric'],
    registry=registry
)

# Training job counter
ml_training_jobs = Counter(
    'ai_risk_ml_training_jobs_total',
    'Total number of ML training jobs',
    ['model_type', 'status'],
    registry=registry
)

# ============================================================
# SIMULATION METRICS
# ============================================================

# Simulation runs counter
simulation_runs = Counter(
    'ai_risk_simulation_runs_total',
    'Total number of simulation runs',
    ['scenario_type', 'status'],
    registry=registry
)

# Simulation agents gauge
simulation_agents = Gauge(
    'ai_risk_simulation_agents_total',
    'Total number of agents in simulation',
    ['agent_type'],
    registry=registry
)

# Simulation duration histogram
simulation_duration = Histogram(
    'ai_risk_simulation_duration_seconds',
    'Simulation execution time',
    ['scenario_type'],
    buckets=(1, 5, 10, 30, 60, 120, 300, 600),
    registry=registry
)

# ============================================================
# ALERT METRICS
# ============================================================

# Alert counter
alert_count = Counter(
    'ai_risk_alerts_triggered_total',
    'Total number of alerts triggered',
    ['alert_level', 'alert_type'],
    registry=registry
)

# Alert response time
alert_response_time = Histogram(
    'ai_risk_alert_response_time_seconds',
    'Time to generate and send alert',
    ['alert_level'],
    buckets=(0.01, 0.05, 0.1, 0.25, 0.5, 1.0),
    registry=registry
)

# ============================================================
# SYSTEM METRICS
# ============================================================

# Database connections gauge
db_connections = Gauge(
    'ai_risk_db_connections_active',
    'Number of active database connections',
    ['database'],
    registry=registry
)

# Cache metrics
cache_hits = Counter(
    'ai_risk_cache_hits_total',
    'Total number of cache hits',
    ['cache_type'],
    registry=registry
)

cache_misses = Counter(
    'ai_risk_cache_misses_total',
    'Total number of cache misses',
    ['cache_type'],
    registry=registry
)

# Memory usage gauge
memory_usage = Gauge(
    'ai_risk_memory_usage_bytes',
    'Memory usage in bytes',
    ['component'],
    registry=registry
)

# CPU usage gauge
cpu_usage = Gauge(
    'ai_risk_cpu_usage_percent',
    'CPU usage percentage',
    ['component'],
    registry=registry
)

# ============================================================
# BATCH PROCESSING METRICS
# ============================================================

# Batch job counter
batch_jobs = Counter(
    'ai_risk_batch_jobs_total',
    'Total number of batch jobs',
    ['job_type', 'status'],
    registry=registry
)

# Batch job duration
batch_job_duration = Histogram(
    'ai_risk_batch_job_duration_seconds',
    'Batch job execution time',
    ['job_type'],
    buckets=(10, 30, 60, 120, 300, 600, 1200, 3600),
    registry=registry
)

# Batch items processed
batch_items_processed = Counter(
    'ai_risk_batch_items_processed_total',
    'Total number of items processed in batch',
    ['job_type'],
    registry=registry
)

# ============================================================
# WEBSOCKET METRICS
# ============================================================

# WebSocket connections gauge
websocket_connections = Gauge(
    'ai_risk_websocket_connections_active',
    'Number of active WebSocket connections',
    registry=registry
)

# WebSocket messages counter
websocket_messages = Counter(
    'ai_risk_websocket_messages_total',
    'Total number of WebSocket messages',
    ['direction', 'message_type'],
    registry=registry
)

# ============================================================
# CUSTOM COLLECTORS
# ============================================================

class RiskEngineCollector:
    """Custom collector for complex metrics"""

    def __init__(self, risk_engine):
        self.risk_engine = risk_engine

    def collect(self):
        """Collect custom metrics"""
        # Collect portfolio metrics
        portfolio_metrics = self.risk_engine.get_portfolio_metrics()

        yield GaugeMetricFamily(
            'ai_risk_portfolio_health_factor',
            'Average portfolio health factor',
            value=portfolio_metrics.get('avg_health_factor', 0)
        )

        yield GaugeMetricFamily(
            'ai_risk_portfolio_total_value_usd',
            'Total portfolio value in USD',
            value=portfolio_metrics.get('total_value_usd', 0)
        )

        # Collect market metrics
        market_metrics = self.risk_engine.get_market_metrics()

        yield GaugeMetricFamily(
            'ai_risk_market_volatility',
            'Current market volatility',
            value=market_metrics.get('volatility', 0)
        )

        yield GaugeMetricFamily(
            'ai_risk_market_liquidity_ratio',
            'Market liquidity ratio',
            value=market_metrics.get('liquidity_ratio', 0)
        )

# ============================================================
# METRIC HELPERS
# ============================================================

class MetricsManager:
    """Manage Prometheus metrics"""

    def __init__(self, port: int = 9090):
        self.port = port
        self.registry = registry

    def start_metrics_server(self):
        """Start Prometheus metrics HTTP server"""
        start_http_server(self.port, registry=self.registry)
        print(f"Metrics server started on port {self.port}")

    def push_to_gateway(self, gateway_url: str, job: str = 'ai_risk_engine'):
        """Push metrics to Prometheus Pushgateway"""
        push_to_gateway(gateway_url, job=job, registry=self.registry)

    @staticmethod
    def record_request(method: str, endpoint: str, status: int, duration: float):
        """Record API request metrics"""
        request_count.labels(
            method=method,
            endpoint=endpoint,
            status=str(status)
        ).inc()

        request_duration.labels(
            method=method,
            endpoint=endpoint
        ).observe(duration)

    @staticmethod
    def record_risk_metrics(risk_data: Dict):
        """Record risk-related metrics"""
        # Update risk scores
        for asset, score in risk_data.get('risk_scores', {}).items():
            risk_score_gauge.labels(
                protocol=risk_data.get('protocol', 'unknown'),
                asset=asset
            ).set(score)

        # Update positions at risk
        for level, count in risk_data.get('positions_at_risk', {}).items():
            positions_at_risk.labels(risk_level=level).set(count)

        # Update total value at risk
        total_value_at_risk.labels(
            protocol=risk_data.get('protocol', 'unknown')
        ).set(risk_data.get('total_value_at_risk', 0))

    @staticmethod
    def record_ml_prediction(
        model_type: str,
        prediction_type: str,
        duration: float,
        success: bool = True
    ):
        """Record ML prediction metrics"""
        ml_predictions.labels(
            model_type=model_type,
            prediction_type=prediction_type
        ).inc()

        ml_prediction_latency.labels(
            model_type=model_type
        ).observe(duration)

    @staticmethod
    def record_simulation(
        scenario_type: str,
        duration: float,
        num_agents: int,
        success: bool = True
    ):
        """Record simulation metrics"""
        simulation_runs.labels(
            scenario_type=scenario_type,
            status='success' if success else 'failure'
        ).inc()

        simulation_duration.labels(
            scenario_type=scenario_type
        ).observe(duration)

        simulation_agents.labels(
            agent_type='all'
        ).set(num_agents)

    @staticmethod
    def record_alert(alert_level: str, alert_type: str, response_time: float):
        """Record alert metrics"""
        alert_count.labels(
            alert_level=alert_level,
            alert_type=alert_type
        ).inc()

        alert_response_time.labels(
            alert_level=alert_level
        ).observe(response_time)

    @staticmethod
    def record_cache_operation(cache_type: str, hit: bool):
        """Record cache operation"""
        if hit:
            cache_hits.labels(cache_type=cache_type).inc()
        else:
            cache_misses.labels(cache_type=cache_type).inc()

    @staticmethod
    def record_batch_job(
        job_type: str,
        status: str,
        duration: float,
        items_processed: int
    ):
        """Record batch job metrics"""
        batch_jobs.labels(
            job_type=job_type,
            status=status
        ).inc()

        if status == 'completed':
            batch_job_duration.labels(job_type=job_type).observe(duration)
            batch_items_processed.labels(job_type=job_type).inc(items_processed)

    @staticmethod
    def update_system_metrics(metrics: Dict):
        """Update system metrics"""
        # Update database connections
        db_connections.labels(
            database='postgresql'
        ).set(metrics.get('db_connections', 0))

        # Update memory usage
        memory_usage.labels(
            component='api'
        ).set(metrics.get('memory_bytes', 0))

        # Update CPU usage
        cpu_usage.labels(
            component='api'
        ).set(metrics.get('cpu_percent', 0))

        # Update WebSocket connections
        websocket_connections.set(metrics.get('websocket_connections', 0))

# ============================================================
# METRIC DECORATORS
# ============================================================

def track_request_metrics(endpoint: str):
    """Decorator to track request metrics"""
    def decorator(func):
        async def wrapper(*args, **kwargs):
            import time
            start_time = time.time()

            active_requests.inc()
            try:
                result = await func(*args, **kwargs)
                status = 200
            except Exception as e:
                status = 500
                raise e
            finally:
                duration = time.time() - start_time
                active_requests.dec()

                MetricsManager.record_request(
                    method=kwargs.get('method', 'GET'),
                    endpoint=endpoint,
                    status=status,
                    duration=duration
                )

            return result
        return wrapper
    return decorator

def track_ml_metrics(model_type: str):
    """Decorator to track ML model metrics"""
    def decorator(func):
        async def wrapper(*args, **kwargs):
            import time
            start_time = time.time()

            try:
                result = await func(*args, **kwargs)
                success = True
            except Exception as e:
                success = False
                raise e
            finally:
                duration = time.time() - start_time

                MetricsManager.record_ml_prediction(
                    model_type=model_type,
                    prediction_type=kwargs.get('prediction_type', 'unknown'),
                    duration=duration,
                    success=success
                )

            return result
        return wrapper
    return decorator

# ============================================================
# EXAMPLE USAGE
# ============================================================

def example_usage():
    """Example of using Prometheus metrics"""

    # Initialize metrics manager
    metrics_manager = MetricsManager(port=9090)

    # Start metrics server
    metrics_manager.start_metrics_server()

    # Record some metrics
    MetricsManager.record_request(
        method='GET',
        endpoint='/api/risk/assess',
        status=200,
        duration=0.125
    )

    MetricsManager.record_risk_metrics({
        'protocol': 'aave',
        'risk_scores': {'ETH': 65.5, 'BTC': 72.3},
        'positions_at_risk': {'high': 5, 'medium': 12, 'low': 30},
        'total_value_at_risk': 1500000
    })

    MetricsManager.record_ml_prediction(
        model_type='lstm',
        prediction_type='liquidation',
        duration=0.085,
        success=True
    )

    MetricsManager.record_simulation(
        scenario_type='market_crash',
        duration=8.5,
        num_agents=10000,
        success=True
    )

    MetricsManager.record_alert(
        alert_level='critical',
        alert_type='liquidation_cascade',
        response_time=0.045
    )

    print("Metrics recorded successfully!")
    print(f"View metrics at http://localhost:{metrics_manager.port}/metrics")

if __name__ == "__main__":
    example_usage()