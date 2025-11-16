"""
Data Infrastructure Module
Handles all data collection, storage, and retrieval for AI risk engine
Phase 1: Core Implementation
"""

import asyncio
import logging
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any
import asyncpg
import redis.asyncio as aioredis
import numpy as np
from dataclasses import dataclass, asdict
import json

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# ============================================================
# DATA MODELS
# ============================================================

@dataclass
class MarketSnapshot:
    """Complete market snapshot for AI analysis"""
    timestamp: datetime
    price_data: Dict[str, float]
    liquidations: List[Dict]
    market_depth: Dict[str, Any]
    user_metrics: Dict[str, Any]
    volatility_metrics: Dict[str, float]
    correlation_matrix: Optional[np.ndarray] = None

@dataclass
class PriceData:
    """Price data point"""
    asset: str
    price: float
    source: str
    volume: Optional[float] = None
    market_cap: Optional[float] = None
    timestamp: datetime = None

@dataclass
class LiquidationData:
    """Liquidation event data"""
    protocol: str
    user_address: str
    collateral_asset: str
    collateral_amount: float
    debt_amount: float
    health_factor: float
    timestamp: datetime
    tx_hash: str

# ============================================================
# DATA INFRASTRUCTURE CLASS
# ============================================================

class DataInfrastructure:
    """
    Unified data infrastructure for AI risk engine
    Integrates with TimescaleDB and provides fast data access
    """

    def __init__(self, db_config: Dict, redis_config: Dict):
        """
        Initialize data infrastructure

        Args:
            db_config: PostgreSQL/TimescaleDB connection config
            redis_config: Redis connection config
        """
        self.db_config = db_config
        self.redis_config = redis_config
        self.db_pool = None
        self.redis = None
        self._initialized = False

    async def initialize(self):
        """Initialize database connections"""
        if self._initialized:
            return

        try:
            # Create PostgreSQL connection pool
            self.db_pool = await asyncpg.create_pool(
                host=self.db_config.get('host', 'localhost'),
                port=self.db_config.get('port', 5432),
                user=self.db_config.get('user', 'postgres'),
                password=self.db_config.get('password', 'postgres'),
                database=self.db_config.get('database', 'loyalty_points'),
                min_size=10,
                max_size=20,
                command_timeout=60
            )

            # Create Redis connection
            redis_password = self.redis_config.get('password')
            redis_url = f"redis://{self.redis_config.get('host', 'localhost')}:{self.redis_config.get('port', 6379)}"
            if redis_password:
                redis_url = f"redis://:{redis_password}@{self.redis_config.get('host', 'localhost')}:{self.redis_config.get('port', 6379)}"

            self.redis = await aioredis.from_url(
                redis_url,
                encoding="utf-8",
                decode_responses=True
            )

            self._initialized = True
            logger.info("Data infrastructure initialized successfully")

        except Exception as e:
            logger.error(f"Failed to initialize data infrastructure: {e}")
            raise

    async def close(self):
        """Close all connections"""
        if self.db_pool:
            await self.db_pool.close()
        if self.redis:
            await self.redis.close()

    # ============================================================
    # DATA COLLECTION METHODS
    # ============================================================

    async def collect_comprehensive_data(self) -> MarketSnapshot:
        """
        Collect comprehensive market data for AI analysis

        Returns:
            MarketSnapshot containing all relevant data
        """
        try:
            # Collect all data types in parallel
            tasks = [
                self.collect_price_history(hours=24),
                self.collect_liquidation_history(days=7),
                self.collect_market_depth(),
                self.collect_user_behavior(),
                self.calculate_volatility_metrics(),
                self.calculate_correlation_matrix()
            ]

            results = await asyncio.gather(*tasks, return_exceptions=True)

            # Handle any exceptions
            for i, result in enumerate(results):
                if isinstance(result, Exception):
                    logger.error(f"Task {i} failed: {result}")
                    results[i] = {}  # Use empty dict as fallback

            # Create market snapshot
            snapshot = MarketSnapshot(
                timestamp=datetime.utcnow(),
                price_data=results[0],
                liquidations=results[1],
                market_depth=results[2],
                user_metrics=results[3],
                volatility_metrics=results[4],
                correlation_matrix=results[5]
            )

            # Cache the snapshot
            await self._cache_snapshot(snapshot)

            return snapshot

        except Exception as e:
            logger.error(f"Failed to collect comprehensive data: {e}")
            raise

    async def collect_price_history(self, hours: int = 24) -> Dict[str, List[Dict]]:
        """
        Collect historical price data

        Args:
            hours: Number of hours of history to collect

        Returns:
            Price history by asset
        """
        async with self.db_pool.acquire() as conn:
            query = """
                SELECT
                    asset,
                    time,
                    price,
                    source,
                    volume,
                    market_cap
                FROM price_history
                WHERE time > NOW() - INTERVAL '%s hours'
                ORDER BY asset, time DESC
            """

            rows = await conn.fetch(query, hours)

            # Group by asset
            price_history = {}
            for row in rows:
                asset = row['asset']
                if asset not in price_history:
                    price_history[asset] = []

                price_history[asset].append({
                    'time': row['time'].isoformat(),
                    'price': float(row['price']),
                    'source': row['source'],
                    'volume': float(row['volume']) if row['volume'] else None,
                    'market_cap': float(row['market_cap']) if row['market_cap'] else None
                })

            return price_history

    async def collect_liquidation_history(self, days: int = 7) -> List[Dict]:
        """
        Collect liquidation history

        Args:
            days: Number of days of history to collect

        Returns:
            List of liquidation events
        """
        async with self.db_pool.acquire() as conn:
            query = """
                SELECT
                    time,
                    protocol,
                    user_address,
                    collateral_asset,
                    debt_asset,
                    collateral_amount,
                    debt_amount,
                    health_factor,
                    tx_hash
                FROM liquidation_history
                WHERE time > NOW() - INTERVAL '%s days'
                ORDER BY time DESC
                LIMIT 1000
            """

            rows = await conn.fetch(query, days)

            liquidations = []
            for row in rows:
                liquidations.append({
                    'time': row['time'].isoformat(),
                    'protocol': row['protocol'],
                    'user_address': row['user_address'],
                    'collateral_asset': row['collateral_asset'],
                    'debt_asset': row['debt_asset'],
                    'collateral_amount': float(row['collateral_amount']),
                    'debt_amount': float(row['debt_amount']) if row['debt_amount'] else 0,
                    'health_factor': float(row['health_factor']) if row['health_factor'] else None,
                    'tx_hash': row['tx_hash']
                })

            return liquidations

    async def collect_market_depth(self) -> Dict[str, Any]:
        """
        Collect current market depth data

        Returns:
            Market depth metrics by protocol and market
        """
        async with self.db_pool.acquire() as conn:
            query = """
                SELECT
                    protocol,
                    market,
                    time,
                    total_supply,
                    total_borrow,
                    utilization_rate,
                    supply_apy,
                    borrow_apy,
                    available_liquidity
                FROM market_depth_snapshots
                WHERE time > NOW() - INTERVAL '1 hour'
                ORDER BY time DESC
            """

            rows = await conn.fetch(query)

            market_depth = {}
            for row in rows:
                key = f"{row['protocol']}_{row['market']}"
                if key not in market_depth:
                    market_depth[key] = {
                        'protocol': row['protocol'],
                        'market': row['market'],
                        'total_supply': float(row['total_supply']) if row['total_supply'] else 0,
                        'total_borrow': float(row['total_borrow']) if row['total_borrow'] else 0,
                        'utilization_rate': float(row['utilization_rate']) if row['utilization_rate'] else 0,
                        'supply_apy': float(row['supply_apy']) if row['supply_apy'] else 0,
                        'borrow_apy': float(row['borrow_apy']) if row['borrow_apy'] else 0,
                        'available_liquidity': float(row['available_liquidity']) if row['available_liquidity'] else 0,
                        'timestamp': row['time'].isoformat()
                    }

            return market_depth

    async def collect_user_behavior(self) -> Dict[str, Any]:
        """
        Collect user behavior metrics

        Returns:
            User behavior statistics
        """
        async with self.db_pool.acquire() as conn:
            # Get aggregate user metrics
            query = """
                SELECT
                    COUNT(DISTINCT user_address) as active_users,
                    AVG(risk_score) as avg_risk_score,
                    AVG(avg_leverage) as avg_leverage,
                    SUM(liquidation_count) as total_liquidations,
                    AVG(avg_health_factor) as avg_health_factor
                FROM user_risk_profiles
                WHERE last_activity > NOW() - INTERVAL '30 days'
            """

            row = await conn.fetchrow(query)

            # Get risk distribution
            distribution_query = """
                SELECT
                    CASE
                        WHEN risk_score < 20 THEN 'very_low'
                        WHEN risk_score < 40 THEN 'low'
                        WHEN risk_score < 60 THEN 'medium'
                        WHEN risk_score < 80 THEN 'high'
                        ELSE 'very_high'
                    END as risk_level,
                    COUNT(*) as count
                FROM user_risk_profiles
                GROUP BY risk_level
            """

            distribution = await conn.fetch(distribution_query)

            return {
                'active_users': row['active_users'] or 0,
                'avg_risk_score': float(row['avg_risk_score']) if row['avg_risk_score'] else 50.0,
                'avg_leverage': float(row['avg_leverage']) if row['avg_leverage'] else 1.0,
                'total_liquidations': row['total_liquidations'] or 0,
                'avg_health_factor': float(row['avg_health_factor']) if row['avg_health_factor'] else 2.0,
                'risk_distribution': {
                    r['risk_level']: r['count']
                    for r in distribution
                }
            }

    async def calculate_volatility_metrics(self) -> Dict[str, float]:
        """
        Calculate volatility metrics for major assets

        Returns:
            Volatility metrics by asset
        """
        async with self.db_pool.acquire() as conn:
            # Use the built-in volatility function
            assets = ['ETH', 'BTC', 'USDC', 'USDT']
            volatility = {}

            for asset in assets:
                # This would normally map to actual asset addresses
                query = "SELECT calculate_volatility($1, 30) as vol"
                result = await conn.fetchval(query, asset)
                volatility[asset] = float(result) if result else 0.0

            return volatility

    async def calculate_correlation_matrix(self) -> Optional[np.ndarray]:
        """
        Calculate correlation matrix for major assets

        Returns:
            Correlation matrix as numpy array
        """
        try:
            # Get price data for correlation calculation
            price_data = await self.collect_price_history(hours=24*7)  # 7 days

            if len(price_data) < 2:
                return None

            # Extract price series
            assets = list(price_data.keys())[:5]  # Top 5 assets
            price_series = []

            for asset in assets:
                prices = [p['price'] for p in price_data[asset]]
                if prices:
                    price_series.append(prices[:100])  # Use last 100 data points

            if len(price_series) < 2:
                return None

            # Calculate correlation matrix
            correlation_matrix = np.corrcoef(price_series)

            return correlation_matrix

        except Exception as e:
            logger.error(f"Failed to calculate correlation matrix: {e}")
            return None

    # ============================================================
    # DATA STORAGE METHODS
    # ============================================================

    async def store_price_data(self, price_data: List[PriceData]):
        """Store price data points"""
        async with self.db_pool.acquire() as conn:
            # Prepare batch insert
            records = [
                (
                    p.timestamp or datetime.utcnow(),
                    p.asset,
                    p.price,
                    p.source,
                    p.volume,
                    p.market_cap
                )
                for p in price_data
            ]

            # Batch insert
            await conn.executemany(
                """
                INSERT INTO price_history (time, asset, price, source, volume, market_cap)
                VALUES ($1, $2, $3, $4, $5, $6)
                ON CONFLICT (time, asset, source) DO UPDATE
                SET price = EXCLUDED.price,
                    volume = EXCLUDED.volume,
                    market_cap = EXCLUDED.market_cap
                """,
                records
            )

            logger.info(f"Stored {len(records)} price data points")

    async def store_liquidation_event(self, liquidation: LiquidationData):
        """Store liquidation event"""
        async with self.db_pool.acquire() as conn:
            await conn.execute(
                """
                INSERT INTO liquidation_history (
                    time, protocol, user_address, collateral_asset,
                    debt_asset, collateral_amount, debt_amount,
                    health_factor, tx_hash, block_number
                ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                """,
                liquidation.timestamp,
                liquidation.protocol,
                liquidation.user_address,
                liquidation.collateral_asset,
                liquidation.debt_asset,
                liquidation.collateral_amount,
                liquidation.debt_amount,
                liquidation.health_factor,
                liquidation.tx_hash,
                0  # block_number placeholder
            )

            # Update user risk profile
            await self.update_user_risk_profile(
                liquidation.user_address,
                {'liquidation_count': 1}  # Increment
            )

    async def update_user_risk_profile(self, user_address: str, updates: Dict):
        """Update user risk profile"""
        async with self.db_pool.acquire() as conn:
            # Build update query dynamically
            set_clauses = []
            values = []
            i = 2

            for key, value in updates.items():
                if key == 'liquidation_count':
                    set_clauses.append(f"liquidation_count = liquidation_count + ${i}")
                else:
                    set_clauses.append(f"{key} = ${i}")
                values.append(value)
                i += 1

            values.append(datetime.utcnow())
            set_clauses.append(f"last_updated = ${i}")

            query = f"""
                INSERT INTO user_risk_profiles (user_address, {', '.join(updates.keys())}, last_updated)
                VALUES ($1, {', '.join([f'${j}' for j in range(2, i)])}, ${i})
                ON CONFLICT (user_address) DO UPDATE
                SET {', '.join(set_clauses)}
            """

            await conn.execute(query, user_address, *values)

    # ============================================================
    # CACHING METHODS
    # ============================================================

    async def _cache_snapshot(self, snapshot: MarketSnapshot):
        """Cache market snapshot in Redis"""
        if not self.redis:
            return

        try:
            # Serialize snapshot
            cache_data = {
                'timestamp': snapshot.timestamp.isoformat(),
                'price_data': snapshot.price_data,
                'liquidations': snapshot.liquidations[:10],  # Cache only recent
                'market_depth': snapshot.market_depth,
                'user_metrics': snapshot.user_metrics,
                'volatility_metrics': snapshot.volatility_metrics
            }

            # Store with 5 minute TTL
            await self.redis.setex(
                'market_snapshot:latest',
                300,  # 5 minutes
                json.dumps(cache_data)
            )

        except Exception as e:
            logger.error(f"Failed to cache snapshot: {e}")

    async def get_cached_snapshot(self) -> Optional[MarketSnapshot]:
        """Get cached market snapshot"""
        if not self.redis:
            return None

        try:
            data = await self.redis.get('market_snapshot:latest')
            if not data:
                return None

            cache_data = json.loads(data)

            return MarketSnapshot(
                timestamp=datetime.fromisoformat(cache_data['timestamp']),
                price_data=cache_data['price_data'],
                liquidations=cache_data['liquidations'],
                market_depth=cache_data['market_depth'],
                user_metrics=cache_data['user_metrics'],
                volatility_metrics=cache_data['volatility_metrics'],
                correlation_matrix=None  # Don't cache numpy arrays
            )

        except Exception as e:
            logger.error(f"Failed to get cached snapshot: {e}")
            return None

    # ============================================================
    # QUERY METHODS
    # ============================================================

    async def get_high_risk_positions(self, limit: int = 100) -> List[Dict]:
        """Get current high risk positions"""
        async with self.db_pool.acquire() as conn:
            query = """
                SELECT * FROM high_risk_positions
                LIMIT $1
            """

            rows = await conn.fetch(query, limit)

            return [dict(row) for row in rows]

    async def get_liquidation_probability(self, protocol: str, health_factor: float) -> float:
        """Get liquidation probability for given health factor"""
        async with self.db_pool.acquire() as conn:
            result = await conn.fetchval(
                "SELECT get_liquidation_probability($1, $2)",
                protocol,
                health_factor
            )

            return float(result) if result else 0.0

    async def get_user_risk_score(self, user_address: str) -> Optional[float]:
        """Get user's current risk score"""
        async with self.db_pool.acquire() as conn:
            result = await conn.fetchval(
                "SELECT risk_score FROM user_risk_profiles WHERE user_address = $1",
                user_address
            )

            return float(result) if result else None