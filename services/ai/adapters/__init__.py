"""
AI Multi-Chain Adapter Layer
=============================

This module provides abstract adapters for the AI brain to interact with
multiple blockchains in a unified way.

Architecture:
    ChainAdapter (Abstract Base)
        ├── ArbitrumAdapter (Aggressive DeFi)
        │   ├── GMXAdapter
        │   ├── AaveAdapter
        │   └── CompoundAdapter
        │
        └── BaseAdapter (Conservative RWA)
            ├── TreasuryAdapter
            ├── AerodromeAdapter
            └── RWAAdapter

The AI brain uses these adapters without knowing which chain it's talking to.
"""

from .base_adapter import ChainAdapter, Position, AssetPrice, ProtocolMetrics
from .arbitrum_adapter import ArbitrumAdapter
from .base_chain_adapter import BaseAdapter
from .multi_chain_manager import MultiChainManager

__all__ = [
    'ChainAdapter',
    'Position',
    'AssetPrice',
    'ProtocolMetrics',
    'ArbitrumAdapter',
    'BaseAdapter',
    'MultiChainManager',
]
