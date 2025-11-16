"""
Base Chain Adapter
==================

Concrete implementation for Base chain (conservative RWA/Treasury).

Supported protocols:
- US Treasury Bonds (tokenized)
- Aerodrome (Base's largest DEX)
- Aave V3 (conservative lending only)
- Backed Finance (bIB01 - tokenized US Treasury)
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


class BaseAdapter(ChainAdapter):
    """Base chain adapter - focuses on conservative RWA strategies"""

    CHAIN_IDS = {
        'testnet': 84532,   # Base Sepolia
        'mainnet': 8453     # Base Mainnet
    }

    SUPPORTED_PROTOCOLS = {
        'treasury': 'US Treasury Bonds',
        'backed_finance': 'Backed Finance bIB01',
        'aerodrome': 'Aerodrome DEX',
        'aave': 'Aave V3 (conservative)',
    }

    # Base focuses on stable assets
    STABLE_ASSETS = ['USDC', 'USDbC', 'DAI', 'USDT']
    TREASURY_APY = 4.5  # Current US Treasury yield ~4.5%

    def __init__(self, network: str = 'testnet', db_connection=None):
        chain_id = self.CHAIN_IDS[network]
        chain_name = f"Base {'Sepolia' if network == 'testnet' else 'Mainnet'}"
        super().__init__(chain_id, chain_name, db_connection)
        self.network = network

    # ========================================================================
    # Position Management
    # ========================================================================

    async def get_user_positions(self, user_address: str) -> List[Position]:
        """Get all positions for user on Base (treasury + DeFi)"""
        positions = []

        # Get treasury bonds
        treasury_positions = await self._get_treasury_positions(user_address)
        positions.extend(treasury_positions)

        # Get DeFi positions (Aerodrome, Aave, etc.)
        defi_positions = await self._get_defi_positions(user_address)
        positions.extend(defi_positions)

        return positions

    async def _get_treasury_positions(self, user_address: str) -> List[Position]:
        """Get US Treasury bond positions"""
        query = """
        SELECT
            user_id,
            bond_type,
            token_amount,
            principal_usd,
            purchase_date,
            last_yield_date,
            total_yield_earned,
            compounding_enabled,
            chain_id,
            created_at,
            updated_at
        FROM treasury_holdings
        WHERE user_id = $1
          AND chain_id = $2
        ORDER BY purchase_date DESC
        """

        rows = await self.db.fetch(query, user_address.lower(), self.chain_id)

        positions = []
        for row in rows:
            # Calculate current value (principal + accrued yield)
            days_held = (datetime.now() - row['purchase_date']).days
            accrued_yield = float(row['principal_usd']) * (self.TREASURY_APY / 100) * (days_held / 365)
            current_value = float(row['principal_usd']) + accrued_yield

            positions.append(Position(
                position_id=f"treasury_{row['user_id']}_{row['bond_type']}",
                user_address=row['user_id'],
                chain_id=row['chain_id'],
                protocol='treasury',
                position_type=PositionType.TREASURY_BOND,

                collateral_asset=row['bond_type'],  # e.g., 'T_BILL_1Y'
                collateral_amount=float(row['token_amount']),
                collateral_value_usd=current_value,

                debt_asset=None,
                debt_amount=None,
                debt_value_usd=0.0,

                health_factor=float('inf'),  # Treasury bonds can't be liquidated
                liquidation_price=None,
                leverage=1.0,  # No leverage

                apy=self.TREASURY_APY,
                earned_usd=float(row['total_yield_earned']) + accrued_yield,

                opened_at=row['purchase_date'],
                last_updated=row['updated_at'],

                metadata={
                    'bond_type': row['bond_type'],
                    'compounding': row['compounding_enabled'],
                    'principal': float(row['principal_usd']),
                    'days_held': days_held
                }
            ))

        return positions

    async def _get_defi_positions(self, user_address: str) -> List[Position]:
        """Get DeFi positions (Aerodrome LP, Aave lending, etc.)"""
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

                debt_asset=None,  # Base positions are mostly debt-free
                debt_amount=None,
                debt_value_usd=float(row['debt_value_usd'] or 0),

                health_factor=float(row['health_factor']) if row['health_factor'] else float('inf'),
                liquidation_price=None,
                leverage=1.0,  # Base is low-leverage

                apy=float(row['apy']) if row['apy'] else None,
                earned_usd=float(row['earned'] or 0),

                opened_at=row['created_at'],
                last_updated=row['updated_at'],

                metadata={'db_id': row['id']}
            ))

        return positions

    async def get_position_by_id(self, position_id: str) -> Optional[Position]:
        """Get specific position"""
        # Check if it's a treasury position
        if position_id.startswith('treasury_'):
            parts = position_id.split('_')
            if len(parts) >= 3:
                user_id = parts[1]
                positions = await self._get_treasury_positions(user_id)
                for p in positions:
                    if p.position_id == position_id:
                        return p

        # Check DeFi positions
        query = """
        SELECT * FROM vault_positions
        WHERE position_id = $1 AND chain_id = $2
        LIMIT 1
        """
        row = await self.db.fetchrow(query, position_id, self.chain_id)
        # ... convert to Position

        return None

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
        """Get asset price on Base"""
        # For treasury bonds, return par value (1.0)
        if asset_symbol.startswith('T_'):
            return AssetPrice(
                asset_symbol=asset_symbol,
                price_usd=1.0,  # Treasury bonds are priced at par
                chain_id=self.chain_id,
                source='treasury_oracle',
                timestamp=datetime.now(),
                volatility=0.01  # Very low volatility
            )

        # For other assets, query price history
        query = """
        SELECT asset, price, source, time, volume, market_cap
        FROM price_history
        WHERE asset = $1
          AND chain_id = $2
        ORDER BY time DESC
        LIMIT 1
        """

        row = await self.db.fetchrow(query, asset_symbol.upper(), self.chain_id)

        if not row:
            # Fallback to chain-agnostic
            query_fallback = """
            SELECT asset, price, source, time, volume, market_cap
            FROM price_history
            WHERE asset = $1
            ORDER BY time DESC
            LIMIT 1
            """
            row = await self.db.fetchrow(query_fallback, asset_symbol.upper())

        if not row:
            # For stablecoins, assume $1
            if asset_symbol.upper() in self.STABLE_ASSETS:
                return AssetPrice(
                    asset_symbol=asset_symbol.upper(),
                    price_usd=1.0,
                    chain_id=self.chain_id,
                    source='stable_peg',
                    timestamp=datetime.now(),
                    volatility=0.001  # 0.1% volatility
                )
            raise ValueError(f"No price data for {asset_symbol}")

        return AssetPrice(
            asset_symbol=row['asset'],
            price_usd=float(row['price']),
            chain_id=self.chain_id,
            source=row['source'],
            timestamp=row['time'],
            volume_24h=float(row['volume']) if row['volume'] else None,
            market_cap=float(row['market_cap']) if row['market_cap'] else None,
            volatility=None
        )

    async def get_historical_prices(
        self,
        asset_symbol: str,
        start_time: datetime,
        end_time: datetime
    ) -> List[AssetPrice]:
        """Get historical prices"""
        # Treasury bonds have stable prices
        if asset_symbol.startswith('T_'):
            # Return flat line at $1
            return [
                AssetPrice(
                    asset_symbol=asset_symbol,
                    price_usd=1.0,
                    chain_id=self.chain_id,
                    source='treasury_oracle',
                    timestamp=start_time,
                    volatility=0.01
                )
            ]

        # Query price history
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
                volume_24h=float(row['volume']) if row['volume'] else None
            )
            for row in rows
        ]

    # ========================================================================
    # Protocol Metrics
    # ========================================================================

    async def get_protocol_metrics(self, protocol: str) -> ProtocolMetrics:
        """Get metrics for Base protocol"""
        if protocol == 'treasury':
            return await self._get_treasury_metrics()

        # For other protocols, query vault_positions
        query = """
        SELECT
            protocol,
            chain_id,
            COUNT(*) as position_count,
            SUM(collateral_value_usd) as total_tvl,
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

        return ProtocolMetrics(
            protocol=row['protocol'],
            chain_id=row['chain_id'],
            tvl_usd=float(row['total_tvl'] or 0),
            total_borrowed_usd=0.0,  # Base is mostly lending, not borrowing
            utilization_rate=0.0,
            avg_apy=float(row['avg_apy']) if row['avg_apy'] else None,
            risk_score=self._calculate_protocol_risk(protocol),
            metadata={'position_count': row['position_count']}
        )

    async def _get_treasury_metrics(self) -> ProtocolMetrics:
        """Get US Treasury metrics"""
        query = """
        SELECT
            COUNT(*) as holder_count,
            SUM(principal_usd) as total_principal,
            SUM(total_yield_earned) as total_yield
        FROM treasury_holdings
        WHERE chain_id = $1
        """

        row = await self.db.fetchrow(query, self.chain_id)

        return ProtocolMetrics(
            protocol='treasury',
            chain_id=self.chain_id,
            tvl_usd=float(row['total_principal'] or 0),
            total_borrowed_usd=0.0,
            utilization_rate=0.0,
            avg_apy=self.TREASURY_APY,
            risk_score=5.0,  # US Treasury is very low risk
            metadata={
                'holder_count': row['holder_count'],
                'total_yield_distributed': float(row['total_yield'] or 0)
            }
        )

    async def get_all_protocols(self) -> List[str]:
        """List all protocols on Base"""
        return list(self.SUPPORTED_PROTOCOLS.keys())

    # ========================================================================
    # Chain Health
    # ========================================================================

    async def get_chain_metrics(self) -> ChainMetrics:
        """Get Base chain metrics"""
        # Treasury TVL
        treasury_query = """
        SELECT SUM(principal_usd) as treasury_tvl
        FROM treasury_holdings
        WHERE chain_id = $1
        """
        treasury_row = await self.db.fetchrow(treasury_query, self.chain_id)
        treasury_tvl = float(treasury_row['treasury_tvl'] or 0)

        # DeFi TVL
        defi_query = """
        SELECT
            COUNT(DISTINCT user_address) as unique_users,
            COUNT(*) as total_positions,
            SUM(collateral_value_usd) as defi_tvl
        FROM vault_positions
        WHERE chain_id = $1
          AND COALESCE(active, true) = true
        """
        defi_row = await self.db.fetchrow(defi_query, self.chain_id)

        total_tvl = treasury_tvl + float(defi_row['defi_tvl'] or 0)

        return ChainMetrics(
            chain_id=self.chain_id,
            chain_name=self.chain_name,
            total_value_locked_usd=total_tvl,
            total_positions=defi_row['total_positions'] or 0,
            avg_gas_price_gwei=0.05,  # Base is very cheap
            block_time_seconds=2.0,  # Base ~2s blocks
            is_healthy=True,
            last_updated=datetime.now()
        )

    async def is_chain_healthy(self) -> bool:
        """Base health check"""
        return True

    # ========================================================================
    # Risk Calculations
    # ========================================================================

    async def calculate_liquidation_price(self, position: Position) -> Optional[float]:
        """Calculate liquidation price (mostly N/A for Base)"""
        # Treasury bonds cannot be liquidated
        if position.protocol == 'treasury':
            return None

        # For Aave (if used conservatively with no borrowing)
        if position.protocol == 'aave' and position.debt_value_usd == 0:
            return None  # No debt = no liquidation

        # If there is debt, calculate (rare on Base)
        if position.debt_value_usd > 0:
            current_price_obj = await self.get_asset_price(position.collateral_asset)
            current_price = current_price_obj.price_usd

            if position.health_factor:
                return current_price * (1.0 / position.health_factor)

        return None

    async def calculate_health_factor(self, position: Position) -> Optional[float]:
        """Calculate health factor (usually infinite on Base)"""
        # Treasury bonds always healthy
        if position.protocol == 'treasury':
            return float('inf')

        # No debt positions are always healthy
        if position.debt_value_usd == 0:
            return float('inf')

        # If has debt (rare), calculate
        if position.debt_value_usd > 0:
            liquidation_threshold = 0.85  # Higher threshold for Base (conservative)
            hf = (position.collateral_value_usd * liquidation_threshold) / position.debt_value_usd
            return hf

        return float('inf')

    # ========================================================================
    # Simulation
    # ========================================================================

    async def simulate_price_impact(
        self,
        position: Position,
        price_change_pct: float
    ) -> Dict[str, Any]:
        """Simulate price change (minimal impact on Base)"""
        # Treasury bonds are stable
        if position.protocol == 'treasury':
            return {
                'original_collateral_value_usd': position.collateral_value_usd,
                'new_collateral_value_usd': position.collateral_value_usd,  # No change
                'original_health_factor': float('inf'),
                'new_health_factor': float('inf'),
                'liquidation_probability': 0.0,
                'will_be_liquidated': False,
                'price_change_pct': price_change_pct,
                'note': 'Treasury bonds are stable and cannot be liquidated'
            }

        # Stablecoins have minimal price impact
        if position.collateral_asset in self.STABLE_ASSETS:
            # Stablecoins typically move <1%
            actual_change = min(abs(price_change_pct), 1.0) * (1 if price_change_pct > 0 else -1)
            new_collateral_value = position.collateral_value_usd * (1 + actual_change / 100)
        else:
            new_collateral_value = position.collateral_value_usd * (1 + price_change_pct / 100)

        return {
            'original_collateral_value_usd': position.collateral_value_usd,
            'new_collateral_value_usd': new_collateral_value,
            'original_health_factor': position.health_factor,
            'new_health_factor': float('inf'),  # Base positions rarely liquidate
            'liquidation_probability': 0.01,  # <1% on Base
            'will_be_liquidated': False,
            'price_change_pct': price_change_pct,
            'note': 'Base positions are conservative and low-risk'
        }

    # ========================================================================
    # Chain Characteristics
    # ========================================================================

    def get_chain_characteristics(self) -> Dict[str, Any]:
        """Base chain characteristics for AI"""
        return {
            'risk_profile': 'conservative',
            'primary_use_cases': ['rwa', 'treasury', 'stable_yield', 'safe_haven'],
            'avg_apy_range': (3.0, 8.0),  # Lower but stable
            'liquidation_risk': 'very_low',
            'gas_cost_level': 'very_low',  # Base is cheapest
            'max_leverage': 1.5,  # Conservative, low leverage
            'supported_protocols': list(self.SUPPORTED_PROTOCOLS.keys()),
            'recommended_for': 'all users, especially risk-averse',
            'benefits': [
                'US Treasury bonds (4.5% APY)',
                'No liquidation risk',
                'Stable, predictable returns',
                'Coinbase integration (Base Pay, Smart Wallet)'
            ]
        }

    # ========================================================================
    # Helper Methods
    # ========================================================================

    def _infer_position_type(self, protocol: str, row: Dict) -> PositionType:
        """Infer position type"""
        if protocol == 'treasury':
            return PositionType.TREASURY_BOND
        elif protocol == 'aerodrome':
            return PositionType.LP_POSITION
        elif protocol == 'aave':
            return PositionType.LENDING
        return PositionType.LENDING

    def _get_collateral_asset(self, protocol: str) -> str:
        """Get collateral asset"""
        if protocol == 'treasury':
            return 'T_BILL_1Y'
        return 'USDC'  # Default to USDC

    def _calculate_protocol_risk(self, protocol: str) -> float:
        """Calculate risk score (0-100)"""
        risk_scores = {
            'treasury': 5.0,       # US Treasury is safest
            'backed_finance': 10.0, # Regulated RWA
            'aerodrome': 20.0,     # DEX LP has some risk
            'aave': 25.0,          # Lending (conservative use)
        }
        return risk_scores.get(protocol, 30.0)
