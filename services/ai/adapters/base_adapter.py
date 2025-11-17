"""
Base Abstract Chain Adapter
============================

Defines the interface that all chain adapters must implement.
This allows the AI brain to work with any blockchain in a unified way.
"""

from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import List, Dict, Optional, Any
from datetime import datetime
from enum import Enum


class ChainType(Enum):
    """Supported chain types"""
    ARBITRUM = "arbitrum"
    BASE = "base"


class PositionType(Enum):
    """Types of positions across chains"""
    LENDING = "lending"              # Aave, Compound
    BORROWING = "borrowing"          # Aave, Compound
    PERPETUAL_LONG = "perpetual_long"    # GMX long
    PERPETUAL_SHORT = "perpetual_short"  # GMX short
    TREASURY_BOND = "treasury_bond"      # US Treasury (Base only)
    LP_POSITION = "lp_position"          # Aerodrome, Uniswap
    STAKING = "staking"                   # General staking


@dataclass
class AssetPrice:
    """Unified asset price structure"""
    asset_symbol: str
    price_usd: float
    chain_id: int
    source: str
    timestamp: datetime
    volume_24h: Optional[float] = None
    market_cap: Optional[float] = None
    volatility: Optional[float] = None  # Annualized volatility


@dataclass
class Position:
    """Unified position structure across all chains and protocols"""
    # Identity
    position_id: str
    user_address: str
    chain_id: int
    protocol: str
    position_type: PositionType

    # Value
    collateral_asset: str
    collateral_amount: float
    collateral_value_usd: float

    debt_asset: Optional[str] = None
    debt_amount: Optional[float] = None
    debt_value_usd: Optional[float] = 0.0

    # Risk metrics
    health_factor: Optional[float] = None
    liquidation_price: Optional[float] = None
    leverage: Optional[float] = 1.0

    # Yield info (for treasury bonds, lending, etc.)
    apy: Optional[float] = None
    earned_usd: Optional[float] = 0.0

    # Timestamps
    opened_at: datetime = None
    last_updated: datetime = None

    # Additional metadata
    metadata: Dict[str, Any] = None


@dataclass
class ProtocolMetrics:
    """Protocol-level metrics for risk assessment"""
    protocol: str
    chain_id: int
    tvl_usd: float
    total_borrowed_usd: Optional[float] = None
    utilization_rate: Optional[float] = None
    avg_apy: Optional[float] = None
    risk_score: Optional[float] = None  # 0-100

    # Protocol-specific
    metadata: Dict[str, Any] = None


@dataclass
class ChainMetrics:
    """Chain-level metrics"""
    chain_id: int
    chain_name: str
    total_value_locked_usd: float
    total_positions: int
    avg_gas_price_gwei: float
    block_time_seconds: float
    is_healthy: bool
    last_updated: datetime


class ChainAdapter(ABC):
    """
    Abstract base class for chain adapters.

    The AI brain uses this interface without knowing which chain it's talking to.
    Each concrete adapter (ArbitrumAdapter, BaseAdapter) implements these methods.
    """

    def __init__(self, chain_id: int, chain_name: str, db_connection):
        self.chain_id = chain_id
        self.chain_name = chain_name
        self.db = db_connection

    # ========================================================================
    # Position Management
    # ========================================================================

    @abstractmethod
    async def get_user_positions(self, user_address: str) -> List[Position]:
        """
        Get all positions for a user on this chain.

        Args:
            user_address: User's wallet address

        Returns:
            List of positions in unified format
        """
        pass

    @abstractmethod
    async def get_position_by_id(self, position_id: str) -> Optional[Position]:
        """Get a specific position by ID"""
        pass

    @abstractmethod
    async def get_positions_by_protocol(
        self,
        user_address: str,
        protocol: str
    ) -> List[Position]:
        """Get user positions in a specific protocol"""
        pass

    # ========================================================================
    # Price Data
    # ========================================================================

    @abstractmethod
    async def get_asset_price(self, asset_symbol: str) -> AssetPrice:
        """
        Get current price for an asset on this chain.

        Args:
            asset_symbol: Asset symbol (e.g., 'ETH', 'USDC')

        Returns:
            Asset price with metadata
        """
        pass

    @abstractmethod
    async def get_historical_prices(
        self,
        asset_symbol: str,
        start_time: datetime,
        end_time: datetime
    ) -> List[AssetPrice]:
        """Get historical price data"""
        pass

    # ========================================================================
    # Protocol Metrics
    # ========================================================================

    @abstractmethod
    async def get_protocol_metrics(self, protocol: str) -> ProtocolMetrics:
        """Get metrics for a specific protocol on this chain"""
        pass

    @abstractmethod
    async def get_all_protocols(self) -> List[str]:
        """List all protocols available on this chain"""
        pass

    # ========================================================================
    # Chain Health
    # ========================================================================

    @abstractmethod
    async def get_chain_metrics(self) -> ChainMetrics:
        """Get overall chain health and metrics"""
        pass

    @abstractmethod
    async def is_chain_healthy(self) -> bool:
        """Quick health check"""
        pass

    # ========================================================================
    # Risk Calculations (Chain-Specific Logic)
    # ========================================================================

    @abstractmethod
    async def calculate_liquidation_price(self, position: Position) -> Optional[float]:
        """
        Calculate liquidation price for a position.
        Implementation varies by protocol and chain.
        """
        pass

    @abstractmethod
    async def calculate_health_factor(self, position: Position) -> Optional[float]:
        """
        Calculate health factor for a position.
        Implementation varies by protocol.
        """
        pass

    # ========================================================================
    # Simulation (for AI predictions)
    # ========================================================================

    @abstractmethod
    async def simulate_price_impact(
        self,
        position: Position,
        price_change_pct: float
    ) -> Dict[str, Any]:
        """
        Simulate impact of price change on a position.

        Args:
            position: The position to simulate
            price_change_pct: Price change percentage (e.g., -10.0 for 10% drop)

        Returns:
            Simulation results with new health_factor, liquidation probability, etc.
        """
        pass

    # ========================================================================
    # Chain Characteristics (for AI strategy)
    # ========================================================================

    @abstractmethod
    def get_chain_characteristics(self) -> Dict[str, Any]:
        """
        Return chain characteristics for AI decision making.

        Example:
            {
                'risk_profile': 'aggressive' | 'conservative',
                'primary_use_cases': ['defi', 'lending'] | ['rwa', 'treasury'],
                'avg_apy_range': (5.0, 50.0),
                'liquidation_risk': 'high' | 'low',
                'gas_cost_level': 'low' | 'medium' | 'high'
            }
        """
        pass

    # ========================================================================
    # Utility Methods
    # ========================================================================

    def get_chain_id(self) -> int:
        """Get chain ID"""
        return self.chain_id

    def get_chain_name(self) -> str:
        """Get chain name"""
        return self.chain_name

    def __repr__(self) -> str:
        return f"<{self.__class__.__name__} chain_id={self.chain_id} name={self.chain_name}>"
