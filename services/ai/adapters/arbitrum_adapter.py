"""
Arbitrum Chain Adapter
======================

Concrete implementation for Arbitrum chain (aggressive DeFi).

Supported protocols:
- GMX V2 (perpetuals)
- Aave V3 (lending/borrowing)
- Compound V3 (lending/borrowing)
- Uniswap V3 (DEX)
"""

from typing import List, Dict, Optional, Any
from datetime import datetime
from .base_adapter import (
    ChainAdapter,
    Position,
    PositionType,
    AssetPrice,
    ProtocolMetrics,
    ChainMetrics
)


class ArbitrumAdapter(ChainAdapter):
    """Arbitrum chain adapter - focuses on aggressive DeFi strategies"""

    CHAIN_IDS = {
        'testnet': 421614,  # Arbitrum Sepolia
        'mainnet': 42161    # Arbitrum One
    }

    SUPPORTED_PROTOCOLS = {
        'gmx': 'GMX V2 Perpetuals',
        'aave': 'Aave V3',
        'compound': 'Compound V3',
        'uniswap': 'Uniswap V3'
    }

    def __init__(self, network: str = 'testnet', db_connection=None):
        chain_id = self.CHAIN_IDS[network]
        chain_name = f"Arbitrum {'Sepolia' if network == 'testnet' else 'One'}"
        super().__init__(chain_id, chain_name, db_connection)
        self.network = network

    # ========================================================================
    # Position Management
    # ========================================================================

    async def get_user_positions(self, user_address: str) -> List[Position]:
        """Get all DeFi positions for user on Arbitrum"""
        query = """
        SELECT
            id,
            user_address,
            chain_id,
            protocol,
            amount,
            earned,
            apy,
            collateral_value_usd,
            debt_value_usd,
            health_factor,
            position_id,
            created_at,
            updated_at
        FROM vault_positions
        WHERE user_address = $1
          AND chain_id = $2
          AND COALESCE(active, true) = true
        ORDER BY updated_at DESC
        """

        rows = await self.db.fetch(query, user_address.lower(), self.chain_id)

        positions = []
        for row in rows:
            # Determine position type based on protocol
            position_type = self._infer_position_type(row['protocol'], row)

            positions.append(Position(
                position_id=row['position_id'] or f"{row['id']}",
                user_address=row['user_address'],
                chain_id=row['chain_id'],
                protocol=row['protocol'],
                position_type=position_type,

                collateral_asset=self._get_collateral_asset(row['protocol']),
                collateral_amount=float(row['amount'] or 0),
                collateral_value_usd=float(row['collateral_value_usd'] or row['amount'] or 0),

                debt_asset=self._get_debt_asset(row['protocol']),
                debt_amount=None,  # TODO: fetch from protocol-specific table
                debt_value_usd=float(row['debt_value_usd'] or 0),

                health_factor=float(row['health_factor']) if row['health_factor'] else None,
                liquidation_price=None,  # Calculate based on health factor
                leverage=self._calculate_leverage(row),

                apy=float(row['apy']) if row['apy'] else None,
                earned_usd=float(row['earned'] or 0),

                opened_at=row['created_at'],
                last_updated=row['updated_at'],

                metadata={'db_id': row['id']}
            ))

        return positions

    async def get_position_by_id(self, position_id: str) -> Optional[Position]:
        """Get specific position by ID"""
        query = """
        SELECT * FROM vault_positions
        WHERE position_id = $1 AND chain_id = $2
        LIMIT 1
        """
        row = await self.db.fetchrow(query, position_id, self.chain_id)

        if not row:
            return None

        # Convert to Position object (similar to above)
        # ... (implement similar to get_user_positions)

    async def get_positions_by_protocol(
        self,
        user_address: str,
        protocol: str
    ) -> List[Position]:
        """Get user positions in specific protocol"""
        all_positions = await self.get_user_positions(user_address)
        return [p for p in all_positions if p.protocol.lower() == protocol.lower()]

    # ========================================================================
    # Price Data
    # ========================================================================

    async def get_asset_price(self, asset_symbol: str) -> AssetPrice:
        """Get current asset price on Arbitrum"""
        query = """
        SELECT
            asset as asset_symbol,
            price,
            source,
            time as timestamp,
            volume,
            market_cap
        FROM price_history
        WHERE asset = $1
          AND chain_id = $2
        ORDER BY time DESC
        LIMIT 1
        """

        row = await self.db.fetchrow(query, asset_symbol.upper(), self.chain_id)

        if not row:
            # Fallback to chain-agnostic price if no chain-specific price
            query_fallback = """
            SELECT asset, price, source, time, volume, market_cap
            FROM price_history
            WHERE asset = $1
            ORDER BY time DESC
            LIMIT 1
            """
            row = await self.db.fetchrow(query_fallback, asset_symbol.upper())

        if not row:
            raise ValueError(f"No price data for {asset_symbol}")

        return AssetPrice(
            asset_symbol=row['asset_symbol'],
            price_usd=float(row['price']),
            chain_id=self.chain_id,
            source=row['source'],
            timestamp=row['timestamp'],
            volume_24h=float(row['volume']) if row['volume'] else None,
            market_cap=float(row['market_cap']) if row['market_cap'] else None,
            volatility=None  # TODO: calculate from historical data
        )

    async def get_historical_prices(
        self,
        asset_symbol: str,
        start_time: datetime,
        end_time: datetime
    ) -> List[AssetPrice]:
        """Get historical prices for volatility calculation"""
        query = """
        SELECT asset, price, source, time, volume, market_cap
        FROM price_history
        WHERE asset = $1
          AND time BETWEEN $2 AND $3
        ORDER BY time ASC
        """

        rows = await self.db.fetch(query, asset_symbol.upper(), start_time, end_time)

        return [
            AssetPrice(
                asset_symbol=row['asset'],
                price_usd=float(row['price']),
                chain_id=self.chain_id,
                source=row['source'],
                timestamp=row['time'],
                volume_24h=float(row['volume']) if row['volume'] else None,
                market_cap=float(row['market_cap']) if row['market_cap'] else None
            )
            for row in rows
        ]

    # ========================================================================
    # Protocol Metrics
    # ========================================================================

    async def get_protocol_metrics(self, protocol: str) -> ProtocolMetrics:
        """Get metrics for Arbitrum DeFi protocol"""
        # Query aggregated protocol data
        query = """
        SELECT
            protocol,
            chain_id,
            COUNT(*) as position_count,
            SUM(collateral_value_usd) as total_tvl,
            SUM(debt_value_usd) as total_borrowed,
            AVG(apy) as avg_apy
        FROM vault_positions
        WHERE protocol = $1
          AND chain_id = $2
          AND COALESCE(active, true) = true
        GROUP BY protocol, chain_id
        """

        row = await self.db.fetchrow(query, protocol.lower(), self.chain_id)

        if not row:
            return ProtocolMetrics(
                protocol=protocol,
                chain_id=self.chain_id,
                tvl_usd=0.0
            )

        utilization = None
        if row['total_tvl'] and row['total_borrowed']:
            utilization = float(row['total_borrowed']) / float(row['total_tvl'])

        return ProtocolMetrics(
            protocol=row['protocol'],
            chain_id=row['chain_id'],
            tvl_usd=float(row['total_tvl'] or 0),
            total_borrowed_usd=float(row['total_borrowed'] or 0),
            utilization_rate=utilization,
            avg_apy=float(row['avg_apy']) if row['avg_apy'] else None,
            risk_score=self._calculate_protocol_risk(protocol, row),
            metadata={'position_count': row['position_count']}
        )

    async def get_all_protocols(self) -> List[str]:
        """List all active protocols on Arbitrum"""
        return list(self.SUPPORTED_PROTOCOLS.keys())

    # ========================================================================
    # Chain Health
    # ========================================================================

    async def get_chain_metrics(self) -> ChainMetrics:
        """Get Arbitrum chain health metrics"""
        query = """
        SELECT
            COUNT(DISTINCT user_address) as unique_users,
            COUNT(*) as total_positions,
            SUM(collateral_value_usd) as total_tvl
        FROM vault_positions
        WHERE chain_id = $1
          AND COALESCE(active, true) = true
        """

        row = await self.db.fetchrow(query, self.chain_id)

        # TODO: Get actual gas price from RPC
        avg_gas_price = 0.1  # gwei (Arbitrum is cheap)

        return ChainMetrics(
            chain_id=self.chain_id,
            chain_name=self.chain_name,
            total_value_locked_usd=float(row['total_tvl'] or 0),
            total_positions=row['total_positions'] or 0,
            avg_gas_price_gwei=avg_gas_price,
            block_time_seconds=0.25,  # Arbitrum ~250ms
            is_healthy=True,  # TODO: implement health check
            last_updated=datetime.now()
        )

    async def is_chain_healthy(self) -> bool:
        """Quick Arbitrum health check"""
        # TODO: Check RPC connectivity, recent blocks, etc.
        return True

    # ========================================================================
    # Risk Calculations
    # ========================================================================

    async def calculate_liquidation_price(self, position: Position) -> Optional[float]:
        """Calculate liquidation price for Arbitrum position"""
        if not position.health_factor:
            return None

        # For lending protocols (Aave, Compound)
        if position.protocol in ['aave', 'compound']:
            # Simplified: liquidation when health factor < 1
            # liquidation_price = current_price * health_factor
            current_price_obj = await self.get_asset_price(position.collateral_asset)
            current_price = current_price_obj.price_usd

            # Rough approximation
            liquidation_price = current_price * (1.0 / position.health_factor)
            return liquidation_price

        # For GMX perpetuals
        elif position.protocol == 'gmx':
            # GMX has different liquidation logic
            # TODO: Implement GMX-specific calculation
            return None

        return None

    async def calculate_health_factor(self, position: Position) -> Optional[float]:
        """Calculate/update health factor"""
        # If already calculated, return it
        if position.health_factor:
            return position.health_factor

        # Otherwise calculate based on protocol
        if position.protocol in ['aave', 'compound']:
            if position.debt_value_usd == 0:
                return float('inf')  # No debt = infinite health

            # Simplified HF = (collateral * liquidation_threshold) / debt
            liquidation_threshold = 0.8  # 80% for most assets on Aave
            hf = (position.collateral_value_usd * liquidation_threshold) / position.debt_value_usd
            return hf

        return None

    # ========================================================================
    # Simulation
    # ========================================================================

    async def simulate_price_impact(
        self,
        position: Position,
        price_change_pct: float
    ) -> Dict[str, Any]:
        """Simulate price drop impact on Arbitrum position"""

        # Calculate new collateral value
        new_collateral_value = position.collateral_value_usd * (1 + price_change_pct / 100)

        # Debt stays the same
        new_debt_value = position.debt_value_usd

        # Calculate new health factor
        if position.protocol in ['aave', 'compound']:
            liquidation_threshold = 0.8
            new_hf = (new_collateral_value * liquidation_threshold) / new_debt_value if new_debt_value > 0 else float('inf')
        else:
            new_hf = None

        # Liquidation probability (heuristic)
        if new_hf is not None:
            if new_hf < 1.0:
                liq_prob = 1.0  # 100% will be liquidated
            elif new_hf < 1.1:
                liq_prob = 0.8  # 80% chance
            elif new_hf < 1.2:
                liq_prob = 0.5  # 50% chance
            elif new_hf < 1.5:
                liq_prob = 0.2  # 20% chance
            else:
                liq_prob = 0.01  # 1% chance
        else:
            liq_prob = 0.0

        return {
            'original_collateral_value_usd': position.collateral_value_usd,
            'new_collateral_value_usd': new_collateral_value,
            'original_health_factor': position.health_factor,
            'new_health_factor': new_hf,
            'liquidation_probability': liq_prob,
            'will_be_liquidated': new_hf < 1.0 if new_hf else False,
            'price_change_pct': price_change_pct
        }

    # ========================================================================
    # Chain Characteristics
    # ========================================================================

    def get_chain_characteristics(self) -> Dict[str, Any]:
        """Arbitrum characteristics for AI strategy"""
        return {
            'risk_profile': 'aggressive',
            'primary_use_cases': ['defi', 'lending', 'perpetuals', 'high_leverage'],
            'avg_apy_range': (5.0, 100.0),  # Can go very high with leverage
            'liquidation_risk': 'high',
            'gas_cost_level': 'low',  # Arbitrum has low gas
            'max_leverage': 50.0,  # GMX supports up to 50x
            'supported_protocols': list(self.SUPPORTED_PROTOCOLS.keys()),
            'recommended_for': 'experienced users seeking high returns',
            'warning': 'High risk of liquidation, requires active monitoring'
        }

    # ========================================================================
    # Helper Methods
    # ========================================================================

    def _infer_position_type(self, protocol: str, row: Dict) -> PositionType:
        """Infer position type from protocol and data"""
        if protocol == 'gmx':
            # TODO: Check if long or short from metadata
            return PositionType.PERPETUAL_LONG
        elif protocol in ['aave', 'compound']:
            if row.get('debt_value_usd', 0) > 0:
                return PositionType.BORROWING
            return PositionType.LENDING
        elif protocol == 'uniswap':
            return PositionType.LP_POSITION
        return PositionType.LENDING  # Default

    def _get_collateral_asset(self, protocol: str) -> str:
        """Get default collateral asset for protocol"""
        # TODO: Get from position metadata
        return 'USDC'  # Default

    def _get_debt_asset(self, protocol: str) -> Optional[str]:
        """Get debt asset if any"""
        if protocol in ['aave', 'compound']:
            return 'USDC'  # TODO: get from metadata
        return None

    def _calculate_leverage(self, row: Dict) -> float:
        """Calculate effective leverage"""
        collateral = float(row.get('collateral_value_usd') or row.get('amount') or 0)
        debt = float(row.get('debt_value_usd') or 0)

        if collateral == 0:
            return 1.0

        return (collateral + debt) / collateral

    def _calculate_protocol_risk(self, protocol: str, metrics: Dict) -> float:
        """Calculate risk score for protocol"""
        # Heuristic risk scoring (0-100)
        risk_scores = {
            'gmx': 85.0,      # High risk (leveraged perpetuals)
            'aave': 40.0,     # Medium risk (lending)
            'compound': 35.0, # Medium-low risk
            'uniswap': 30.0   # Low risk (just swaps)
        }

        base_risk = risk_scores.get(protocol, 50.0)

        # Adjust based on utilization
        if metrics.get('utilization_rate'):
            util = metrics['utilization_rate']
            if util > 0.9:  # >90% utilization is risky
                base_risk += 20
            elif util > 0.8:
                base_risk += 10

        return min(base_risk, 100.0)
