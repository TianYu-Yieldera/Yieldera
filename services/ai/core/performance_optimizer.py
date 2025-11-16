#!/usr/bin/env python3
"""
Performance Optimization and Caching Strategy for AI Risk Engine
Implements advanced caching, connection pooling, and batch processing
"""

import asyncio
import hashlib
import json
import pickle
import time
from concurrent.futures import ThreadPoolExecutor, ProcessPoolExecutor
from dataclasses import dataclass
from datetime import datetime, timedelta
from functools import lru_cache, wraps
from typing import Any, Dict, List, Optional, Tuple, Callable

import numpy as np
import pandas as pd
import redis
from redis import asyncio as aioredis
from sqlalchemy import create_engine
from sqlalchemy.pool import QueuePool
import asyncpg
import msgpack

# ============================================================
# CONFIGURATION
# ============================================================

@dataclass
class PerformanceConfig:
    """Performance optimization configuration"""

    # Database
    db_pool_size: int = 20
    db_max_overflow: int = 10
    db_pool_timeout: int = 30
    db_pool_recycle: int = 3600

    # Redis
    redis_pool_size: int = 50
    redis_decode_responses: bool = False
    redis_socket_keepalive: bool = True

    # Caching
    cache_ttl_short: int = 60  # 1 minute
    cache_ttl_medium: int = 300  # 5 minutes
    cache_ttl_long: int = 3600  # 1 hour
    cache_ttl_persistent: int = 86400  # 24 hours

    # Batch Processing
    batch_size_small: int = 100
    batch_size_medium: int = 1000
    batch_size_large: int = 10000
    max_workers: int = 8

    # Circuit Breaker
    circuit_breaker_threshold: int = 5
    circuit_breaker_timeout: int = 60

    # Rate Limiting
    rate_limit_requests: int = 100
    rate_limit_window: int = 60  # seconds

# ============================================================
# CACHING LAYER
# ============================================================

class AdvancedCache:
    """Advanced caching with multiple strategies"""

    def __init__(self, redis_url: str = "redis://localhost:6379"):
        self.redis_url = redis_url
        self.local_cache = {}  # L1 cache (in-memory)
        self.redis_client = None  # L2 cache (Redis)
        self.stats = {
            'hits': 0,
            'misses': 0,
            'evictions': 0
        }

    async def connect(self):
        """Initialize Redis connection"""
        self.redis_client = await aioredis.create_redis_pool(
            self.redis_url,
            minsize=5,
            maxsize=PerformanceConfig.redis_pool_size
        )

    async def disconnect(self):
        """Close Redis connection"""
        if self.redis_client:
            self.redis_client.close()
            await self.redis_client.wait_closed()

    def _generate_key(self, prefix: str, params: Dict) -> str:
        """Generate cache key from parameters"""
        param_str = json.dumps(params, sort_keys=True)
        hash_digest = hashlib.md5(param_str.encode()).hexdigest()
        return f"{prefix}:{hash_digest}"

    async def get(self, key: str) -> Optional[Any]:
        """Get value from cache (L1 -> L2)"""
        # Check L1 cache
        if key in self.local_cache:
            self.stats['hits'] += 1
            return self.local_cache[key]['value']

        # Check L2 cache (Redis)
        if self.redis_client:
            value = await self.redis_client.get(key)
            if value:
                self.stats['hits'] += 1
                # Deserialize and store in L1
                deserialized = msgpack.unpackb(value, raw=False)
                self.local_cache[key] = {
                    'value': deserialized,
                    'timestamp': time.time()
                }
                return deserialized

        self.stats['misses'] += 1
        return None

    async def set(self, key: str, value: Any, ttl: int = 300):
        """Set value in cache (L1 + L2)"""
        # Store in L1 cache
        self.local_cache[key] = {
            'value': value,
            'timestamp': time.time(),
            'ttl': ttl
        }

        # Store in L2 cache (Redis)
        if self.redis_client:
            serialized = msgpack.packb(value, use_bin_type=True)
            await self.redis_client.setex(key, ttl, serialized)

        # Evict old entries from L1 if too large
        if len(self.local_cache) > 1000:
            self._evict_lru()

    def _evict_lru(self):
        """Evict least recently used items from L1 cache"""
        # Sort by timestamp and remove oldest 10%
        sorted_items = sorted(
            self.local_cache.items(),
            key=lambda x: x[1]['timestamp']
        )

        evict_count = len(sorted_items) // 10
        for key, _ in sorted_items[:evict_count]:
            del self.local_cache[key]
            self.stats['evictions'] += 1

    async def invalidate(self, pattern: str):
        """Invalidate cache entries matching pattern"""
        # Clear L1 cache
        keys_to_remove = [k for k in self.local_cache if pattern in k]
        for key in keys_to_remove:
            del self.local_cache[key]

        # Clear L2 cache
        if self.redis_client:
            cursor = 0
            while True:
                cursor, keys = await self.redis_client.scan(
                    cursor, match=f"*{pattern}*"
                )
                if keys:
                    await self.redis_client.delete(*keys)
                if cursor == 0:
                    break

    def get_stats(self) -> Dict:
        """Get cache statistics"""
        total = self.stats['hits'] + self.stats['misses']
        hit_rate = self.stats['hits'] / total if total > 0 else 0

        return {
            'hits': self.stats['hits'],
            'misses': self.stats['misses'],
            'evictions': self.stats['evictions'],
            'hit_rate': hit_rate,
            'l1_size': len(self.local_cache)
        }

# ============================================================
# CONNECTION POOLING
# ============================================================

class ConnectionPoolManager:
    """Manage database connection pools"""

    def __init__(self):
        self.postgres_pool = None
        self.timescale_engine = None
        self.redis_pool = None

    async def initialize(
        self,
        postgres_dsn: str,
        redis_url: str
    ):
        """Initialize all connection pools"""
        # PostgreSQL async pool
        self.postgres_pool = await asyncpg.create_pool(
            postgres_dsn,
            min_size=10,
            max_size=PerformanceConfig.db_pool_size,
            max_queries=50000,
            max_inactive_connection_lifetime=300
        )

        # TimescaleDB SQLAlchemy engine
        self.timescale_engine = create_engine(
            postgres_dsn,
            poolclass=QueuePool,
            pool_size=PerformanceConfig.db_pool_size,
            max_overflow=PerformanceConfig.db_max_overflow,
            pool_timeout=PerformanceConfig.db_pool_timeout,
            pool_recycle=PerformanceConfig.db_pool_recycle,
            pool_pre_ping=True
        )

        # Redis connection pool
        self.redis_pool = aioredis.ConnectionPool.from_url(
            redis_url,
            max_connections=PerformanceConfig.redis_pool_size,
            decode_responses=PerformanceConfig.redis_decode_responses,
            socket_keepalive=PerformanceConfig.redis_socket_keepalive
        )

    async def get_postgres_connection(self):
        """Get PostgreSQL connection from pool"""
        return await self.postgres_pool.acquire()

    def get_timescale_connection(self):
        """Get TimescaleDB connection from pool"""
        return self.timescale_engine.connect()

    async def get_redis_connection(self):
        """Get Redis connection from pool"""
        return aioredis.Redis(connection_pool=self.redis_pool)

    async def close_all(self):
        """Close all connection pools"""
        if self.postgres_pool:
            await self.postgres_pool.close()
        if self.timescale_engine:
            self.timescale_engine.dispose()
        if self.redis_pool:
            await self.redis_pool.disconnect()

# ============================================================
# BATCH PROCESSING
# ============================================================

class BatchProcessor:
    """High-performance batch processing"""

    def __init__(self):
        self.thread_executor = ThreadPoolExecutor(
            max_workers=PerformanceConfig.max_workers
        )
        self.process_executor = ProcessPoolExecutor(
            max_workers=PerformanceConfig.max_workers // 2
        )

    async def process_in_batches(
        self,
        items: List[Any],
        processor: Callable,
        batch_size: int = None,
        use_multiprocess: bool = False
    ) -> List[Any]:
        """Process items in batches with parallelization"""
        if not items:
            return []

        batch_size = batch_size or PerformanceConfig.batch_size_medium
        results = []

        # Split into batches
        batches = [
            items[i:i + batch_size]
            for i in range(0, len(items), batch_size)
        ]

        # Choose executor
        executor = self.process_executor if use_multiprocess else self.thread_executor

        # Process batches in parallel
        loop = asyncio.get_event_loop()
        futures = []

        for batch in batches:
            future = loop.run_in_executor(executor, processor, batch)
            futures.append(future)

        # Wait for all batches to complete
        batch_results = await asyncio.gather(*futures)

        # Flatten results
        for batch_result in batch_results:
            results.extend(batch_result)

        return results

    def shutdown(self):
        """Shutdown executors"""
        self.thread_executor.shutdown(wait=True)
        self.process_executor.shutdown(wait=True)

# ============================================================
# QUERY OPTIMIZATION
# ============================================================

class QueryOptimizer:
    """Optimize database queries"""

    def __init__(self, connection_pool: ConnectionPoolManager):
        self.pool = connection_pool
        self.query_cache = {}
        self.prepared_statements = {}

    async def prepare_statements(self):
        """Prepare frequently used statements"""
        statements = {
            'get_positions': """
                SELECT * FROM positions
                WHERE user_address = $1 AND protocol = $2
                ORDER BY created_at DESC
                LIMIT $3
            """,
            'get_liquidations': """
                SELECT * FROM liquidations
                WHERE timestamp >= $1 AND timestamp < $2
                ORDER BY timestamp DESC
            """,
            'get_market_data': """
                SELECT * FROM market_data_1h
                WHERE asset_id = $1 AND bucket >= $2
                ORDER BY bucket DESC
                LIMIT $3
            """
        }

        async with await self.pool.get_postgres_connection() as conn:
            for name, query in statements.items():
                self.prepared_statements[name] = await conn.prepare(query)

    async def execute_optimized(
        self,
        query_name: str,
        params: Tuple,
        cache_ttl: int = 60
    ) -> List[Dict]:
        """Execute optimized query with caching"""
        # Generate cache key
        cache_key = f"query:{query_name}:{hash(params)}"

        # Check cache
        if cache_key in self.query_cache:
            cached = self.query_cache[cache_key]
            if time.time() - cached['timestamp'] < cache_ttl:
                return cached['data']

        # Execute query
        async with await self.pool.get_postgres_connection() as conn:
            if query_name in self.prepared_statements:
                # Use prepared statement
                stmt = self.prepared_statements[query_name]
                rows = await stmt.fetch(*params)
            else:
                # Fallback to regular query
                rows = await conn.fetch(query_name, *params)

        # Convert to dict
        result = [dict(row) for row in rows]

        # Cache result
        self.query_cache[cache_key] = {
            'data': result,
            'timestamp': time.time()
        }

        return result

    async def bulk_insert(
        self,
        table: str,
        records: List[Dict],
        on_conflict: str = None
    ):
        """Optimized bulk insert"""
        if not records:
            return

        # Build query
        columns = list(records[0].keys())
        values_template = ','.join([f'${i+1}' for i in range(len(columns))])

        query = f"""
            INSERT INTO {table} ({','.join(columns)})
            VALUES ({values_template})
        """

        if on_conflict:
            query += f" ON CONFLICT {on_conflict}"

        # Prepare data
        values_list = [
            [record[col] for col in columns]
            for record in records
        ]

        # Execute in batches
        async with await self.pool.get_postgres_connection() as conn:
            async with conn.transaction():
                stmt = await conn.prepare(query)
                await stmt.executemany(values_list)

# ============================================================
# CIRCUIT BREAKER
# ============================================================

class CircuitBreaker:
    """Circuit breaker pattern for fault tolerance"""

    def __init__(
        self,
        failure_threshold: int = 5,
        recovery_timeout: int = 60
    ):
        self.failure_threshold = failure_threshold
        self.recovery_timeout = recovery_timeout
        self.failure_count = 0
        self.last_failure_time = None
        self.state = 'CLOSED'  # CLOSED, OPEN, HALF_OPEN

    async def call(self, func: Callable, *args, **kwargs):
        """Execute function with circuit breaker protection"""
        # Check if circuit is open
        if self.state == 'OPEN':
            if self.last_failure_time:
                if time.time() - self.last_failure_time > self.recovery_timeout:
                    self.state = 'HALF_OPEN'
                else:
                    raise Exception("Circuit breaker is OPEN")

        try:
            # Execute function
            result = await func(*args, **kwargs)

            # Reset on success
            if self.state == 'HALF_OPEN':
                self.state = 'CLOSED'
                self.failure_count = 0

            return result

        except Exception as e:
            self.failure_count += 1
            self.last_failure_time = time.time()

            if self.failure_count >= self.failure_threshold:
                self.state = 'OPEN'

            raise e

# ============================================================
# PERFORMANCE MONITOR
# ============================================================

class PerformanceMonitor:
    """Monitor and track performance metrics"""

    def __init__(self):
        self.metrics = {
            'request_count': 0,
            'total_latency': 0,
            'error_count': 0,
            'cache_hits': 0,
            'cache_misses': 0,
            'db_queries': 0,
            'db_latency': 0
        }
        self.latency_histogram = []

    def record_request(self, latency: float, success: bool = True):
        """Record request metrics"""
        self.metrics['request_count'] += 1
        self.metrics['total_latency'] += latency
        self.latency_histogram.append(latency)

        if not success:
            self.metrics['error_count'] += 1

        # Keep histogram size manageable
        if len(self.latency_histogram) > 10000:
            self.latency_histogram = self.latency_histogram[-5000:]

    def record_cache(self, hit: bool):
        """Record cache metrics"""
        if hit:
            self.metrics['cache_hits'] += 1
        else:
            self.metrics['cache_misses'] += 1

    def record_db_query(self, latency: float):
        """Record database query metrics"""
        self.metrics['db_queries'] += 1
        self.metrics['db_latency'] += latency

    def get_stats(self) -> Dict:
        """Get performance statistics"""
        if self.metrics['request_count'] == 0:
            return self.metrics

        # Calculate percentiles
        if self.latency_histogram:
            sorted_latencies = sorted(self.latency_histogram)
            p50 = sorted_latencies[len(sorted_latencies) // 2]
            p95 = sorted_latencies[int(len(sorted_latencies) * 0.95)]
            p99 = sorted_latencies[int(len(sorted_latencies) * 0.99)]
        else:
            p50 = p95 = p99 = 0

        return {
            **self.metrics,
            'avg_latency': self.metrics['total_latency'] / self.metrics['request_count'],
            'error_rate': self.metrics['error_count'] / self.metrics['request_count'],
            'cache_hit_rate': self.metrics['cache_hits'] / (self.metrics['cache_hits'] + self.metrics['cache_misses'])
            if (self.metrics['cache_hits'] + self.metrics['cache_misses']) > 0 else 0,
            'p50_latency': p50,
            'p95_latency': p95,
            'p99_latency': p99,
            'avg_db_latency': self.metrics['db_latency'] / self.metrics['db_queries']
            if self.metrics['db_queries'] > 0 else 0
        }

# ============================================================
# DECORATOR UTILITIES
# ============================================================

def cached(ttl: int = 300):
    """Decorator for caching function results"""
    def decorator(func):
        cache = {}

        @wraps(func)
        async def wrapper(*args, **kwargs):
            # Generate cache key
            key = f"{func.__name__}:{str(args)}:{str(kwargs)}"

            # Check cache
            if key in cache:
                cached_value, cached_time = cache[key]
                if time.time() - cached_time < ttl:
                    return cached_value

            # Execute function
            result = await func(*args, **kwargs)

            # Store in cache
            cache[key] = (result, time.time())

            return result

        return wrapper
    return decorator

def rate_limited(max_calls: int = 10, time_window: int = 60):
    """Decorator for rate limiting"""
    def decorator(func):
        calls = []

        @wraps(func)
        async def wrapper(*args, **kwargs):
            now = time.time()

            # Remove old calls
            calls[:] = [t for t in calls if now - t < time_window]

            # Check rate limit
            if len(calls) >= max_calls:
                raise Exception(f"Rate limit exceeded: {max_calls} calls per {time_window} seconds")

            # Record call
            calls.append(now)

            # Execute function
            return await func(*args, **kwargs)

        return wrapper
    return decorator

def timed():
    """Decorator to measure execution time"""
    def decorator(func):
        @wraps(func)
        async def wrapper(*args, **kwargs):
            start = time.time()
            try:
                result = await func(*args, **kwargs)
                return result
            finally:
                elapsed = time.time() - start
                print(f"{func.__name__} took {elapsed:.3f} seconds")

        return wrapper
    return decorator

# ============================================================
# MAIN PERFORMANCE OPTIMIZER
# ============================================================

class PerformanceOptimizer:
    """Main performance optimization orchestrator"""

    def __init__(
        self,
        postgres_dsn: str,
        redis_url: str
    ):
        self.cache = AdvancedCache(redis_url)
        self.connection_pool = ConnectionPoolManager()
        self.batch_processor = BatchProcessor()
        self.query_optimizer = None
        self.circuit_breaker = CircuitBreaker()
        self.performance_monitor = PerformanceMonitor()

        self.postgres_dsn = postgres_dsn
        self.redis_url = redis_url

    async def initialize(self):
        """Initialize all components"""
        await self.cache.connect()
        await self.connection_pool.initialize(
            self.postgres_dsn,
            self.redis_url
        )
        self.query_optimizer = QueryOptimizer(self.connection_pool)
        await self.query_optimizer.prepare_statements()

    async def shutdown(self):
        """Shutdown all components"""
        await self.cache.disconnect()
        await self.connection_pool.close_all()
        self.batch_processor.shutdown()

    def get_performance_report(self) -> Dict:
        """Get comprehensive performance report"""
        return {
            'cache_stats': self.cache.get_stats(),
            'performance_metrics': self.performance_monitor.get_stats(),
            'circuit_breaker_state': self.circuit_breaker.state,
            'timestamp': datetime.utcnow().isoformat()
        }

# ============================================================
# EXAMPLE USAGE
# ============================================================

async def example_usage():
    """Example of using the performance optimizer"""
    optimizer = PerformanceOptimizer(
        postgres_dsn="postgresql://user:pass@localhost/dbname",
        redis_url="redis://localhost:6379"
    )

    await optimizer.initialize()

    try:
        # Use caching
        @cached(ttl=60)
        async def expensive_calculation(n: int) -> int:
            await asyncio.sleep(1)  # Simulate expensive operation
            return n * n

        # First call - will be slow
        result1 = await expensive_calculation(42)

        # Second call - will be fast (cached)
        result2 = await expensive_calculation(42)

        # Batch processing
        items = list(range(10000))
        results = await optimizer.batch_processor.process_in_batches(
            items,
            lambda batch: [x * 2 for x in batch],
            batch_size=1000
        )

        # Get performance report
        report = optimizer.get_performance_report()
        print(json.dumps(report, indent=2))

    finally:
        await optimizer.shutdown()

if __name__ == "__main__":
    asyncio.run(example_usage())