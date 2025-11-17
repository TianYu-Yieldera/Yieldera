"""
Multi-Chain Manager
===================

This is the **AI Brain's unified interface** to all blockchains.

The AI never talks to chains directly - it always goes through this manager.
This implements the "Chain Isolation + AI Orchestration" architecture.

Architecture:
    AI Brain
        ↓
    MultiChainManager (YOU ARE HERE)
        ↓
    ├── ArbitrumAdapter → Aggressive DeFi
    └── BaseAdapter → Conservative RWA

Key Principles:
- AI reads from both chains
- AI never executes transactions
- AI only provides recommendations
- Users manually execute on each chain
"""

from typing import List, Dict, Optional, Any, Tuple
from datetime import datetime, timedelta
import asyncio

from .base_adapter import ChainAdapter, Position, AssetPrice, ProtocolMetrics, PositionType
from .arbitrum_adapter import ArbitrumAdapter
from .base_chain_adapter import BaseAdapter


class MultiChainManager:
    """
    Unified interface for AI to interact with multiple blockchains.

    This is the ONLY class that the AI risk engine should use.
    """

    def __init__(self, db_connection, network: str = 'testnet'):
        """
        Initialize multi-chain manager.

        Args:
            db_connection: Database connection
            network: 'testnet' or 'mainnet'
        """
        self.db = db_connection
        self.network = network

        # Initialize adapters for each chain
        self.adapters: Dict[int, ChainAdapter] = {}

        # Arbitrum (aggressive)
        arb_adapter = ArbitrumAdapter(network=network, db_connection=db_connection)
        self.adapters[arb_adapter.chain_id] = arb_adapter

        # Base (conservative)
        base_adapter = BaseAdapter(network=network, db_connection=db_connection)
        self.adapters[base_adapter.chain_id] = base_adapter

        print(f"✅ MultiChainManager initialized with {len(self.adapters)} chains:")
        for chain_id, adapter in self.adapters.items():
            print(f"   - {adapter.chain_name} ({chain_id}): {adapter.get_chain_characteristics()['risk_profile']}")

    # ========================================================================
    # Unified Position Queries (AI Brain's Main Interface)
    # ========================================================================

    async def get_all_user_positions(self, user_address: str) -> Dict[int, List[Position]]:
        """
        Get user positions across ALL chains.

        This is the AI brain's main entry point.

        Returns:
            {
                421614: [Position(...), Position(...)],  # Arbitrum positions
                84532: [Position(...)]                    # Base positions
            }
        """
        results = {}

        tasks = [
            self._get_positions_for_chain(chain_id, user_address)
            for chain_id in self.adapters.keys()
        ]

        chain_positions = await asyncio.gather(*tasks)

        for chain_id, positions in zip(self.adapters.keys(), chain_positions):
            results[chain_id] = positions

        return results

    async def _get_positions_for_chain(
        self,
        chain_id: int,
        user_address: str
    ) -> List[Position]:
        """Helper to get positions for one chain"""
        adapter = self.adapters.get(chain_id)
        if not adapter:
            return []

        try:
            return await adapter.get_user_positions(user_address)
        except Exception as e:
            print(f"Error fetching positions from {chain_id}: {e}")
            return []

    async def get_aggregated_portfolio(self, user_address: str) -> Dict[str, Any]:
        """
        Get aggregated portfolio view across all chains.

        Returns:
            {
                'user_address': '0x...',
                'total_value_usd': 100000.0,
                'chains': {
                    421614: {
                        'chain_name': 'Arbitrum Sepolia',
                        'value_usd': 60000.0,
                        'percentage': 60.0,
                        'position_count': 5,
                        'risk_profile': 'aggressive'
                    },
                    84532: {
                        'chain_name': 'Base Sepolia',
                        'value_usd': 40000.0,
                        'percentage': 40.0,
                        'position_count': 2,
                        'risk_profile': 'conservative'
                    }
                },
                'diversification_score': 85.0,  # 0-100, higher = better diversified
                'all_positions': [...],
                'timestamp': datetime.now()
            }
        """
        all_positions_by_chain = await self.get_all_user_positions(user_address)

        # Calculate per-chain values
        chains_data = {}
        total_value = 0.0
        all_positions = []

        for chain_id, positions in all_positions_by_chain.items():
            adapter = self.adapters[chain_id]
            chain_value = sum(p.collateral_value_usd for p in positions)
            total_value += chain_value

            chains_data[chain_id] = {
                'chain_name': adapter.chain_name,
                'value_usd': chain_value,
                'percentage': 0.0,  # Calculate after we know total
                'position_count': len(positions),
                'risk_profile': adapter.get_chain_characteristics()['risk_profile'],
                'positions': positions
            }

            all_positions.extend(positions)

        # Calculate percentages
        for chain_id in chains_data.keys():
            if total_value > 0:
                chains_data[chain_id]['percentage'] = (chains_data[chain_id]['value_usd'] / total_value) * 100

        # Calculate diversification score
        diversification_score = self._calculate_diversification_score(chains_data, total_value)

        return {
            'user_address': user_address,
            'total_value_usd': total_value,
            'chains': chains_data,
            'diversification_score': diversification_score,
            'all_positions': all_positions,
            'timestamp': datetime.now()
        }

    # ========================================================================
    # Cross-Chain Risk Analysis
    # ========================================================================

    async def calculate_cross_chain_risk(self, user_address: str) -> Dict[str, Any]:
        """
        Calculate comprehensive risk across all chains.

        This is what the AI brain uses for recommendations.

        Returns:
            {
                'user_address': '0x...',
                'total_risk_score': 55.0,  # 0-100
                'chain_risks': {
                    421614: {
                        'risk_score': 75.0,
                        'liquidation_probability_24h': 0.05,
                        'high_risk_positions': 2
                    },
                    84532: {
                        'risk_score': 10.0,
                        'liquidation_probability_24h': 0.0,
                        'high_risk_positions': 0
                    }
                },
                'correlation_risk': 0.3,  # How correlated are the two chains
                'diversification_benefit': 25.0,  # Risk reduction from diversification
                'recommendations': [...]
            }
        """
        portfolio = await self.get_aggregated_portfolio(user_address)

        chain_risks = {}
        weighted_risk_scores = []

        for chain_id, chain_data in portfolio['chains'].items():
            adapter = self.adapters[chain_id]
            positions = chain_data['positions']

            # Calculate risk for this chain
            chain_risk = await self._calculate_chain_risk(adapter, positions)
            chain_risks[chain_id] = chain_risk

            # Weight risk by allocation
            allocation_pct = chain_data['percentage'] / 100.0
            weighted_risk_scores.append(chain_risk['risk_score'] * allocation_pct)

        # Aggregate risk score
        total_risk_score = sum(weighted_risk_scores)

        # Calculate correlation (heuristic: Arbitrum and Base are negatively correlated)
        correlation = await self._estimate_correlation(portfolio)

        # Diversification benefit
        diversification_benefit = self._calculate_diversification_benefit(
            chain_risks,
            portfolio['chains'],
            correlation
        )

        return {
            'user_address': user_address,
            'total_risk_score': total_risk_score,
            'chain_risks': chain_risks,
            'correlation_risk': correlation,
            'diversification_benefit': diversification_benefit,
            'portfolio': portfolio,
            'timestamp': datetime.now()
        }

    async def _calculate_chain_risk(
        self,
        adapter: ChainAdapter,
        positions: List[Position]
    ) -> Dict[str, Any]:
        """Calculate risk metrics for one chain"""
        if not positions:
            return {
                'risk_score': 0.0,
                'liquidation_probability_24h': 0.0,
                'high_risk_positions': 0,
                'avg_health_factor': float('inf')
            }

        # Calculate average risk
        total_value = sum(p.collateral_value_usd for p in positions)
        health_factors = [p.health_factor for p in positions if p.health_factor and p.health_factor != float('inf')]

        avg_health_factor = sum(health_factors) / len(health_factors) if health_factors else float('inf')

        # Count high-risk positions (HF < 1.5)
        high_risk_count = sum(1 for p in positions if p.health_factor and p.health_factor < 1.5)

        # Estimate liquidation probability
        liquidation_prob = self._estimate_liquidation_probability(positions)

        # Risk score (0-100)
        characteristics = adapter.get_chain_characteristics()
        base_risk = 85.0 if characteristics['risk_profile'] == 'aggressive' else 10.0

        # Adjust based on health factors
        if avg_health_factor < 1.2:
            risk_score = min(base_risk + 30, 100.0)
        elif avg_health_factor < 1.5:
            risk_score = min(base_risk + 15, 100.0)
        elif avg_health_factor < 2.0:
            risk_score = base_risk
        else:
            risk_score = max(base_risk - 10, 0.0)

        return {
            'risk_score': risk_score,
            'liquidation_probability_24h': liquidation_prob,
            'high_risk_positions': high_risk_count,
            'avg_health_factor': avg_health_factor,
            'total_value_usd': total_value
        }

    def _estimate_liquidation_probability(self, positions: List[Position]) -> float:
        """Estimate 24h liquidation probability"""
        if not positions:
            return 0.0

        liquidatable_positions = [
            p for p in positions
            if p.health_factor and p.health_factor != float('inf') and p.health_factor < 2.0
        ]

        if not liquidatable_positions:
            return 0.0

        # Heuristic based on health factors
        total_prob = 0.0
        for p in liquidatable_positions:
            if p.health_factor < 1.05:
                total_prob += 0.8  # 80% chance
            elif p.health_factor < 1.1:
                total_prob += 0.5  # 50% chance
            elif p.health_factor < 1.2:
                total_prob += 0.2  # 20% chance
            elif p.health_factor < 1.5:
                total_prob += 0.05  # 5% chance

        return min(total_prob / len(positions), 1.0)

    async def _estimate_correlation(self, portfolio: Dict) -> float:
        """
        Estimate correlation between Arbitrum and Base.

        Arbitrum (DeFi) and Base (Treasury) are negatively correlated:
        - Market crash: Arbitrum suffers, Base is stable
        - Market boom: Arbitrum thrives, Base returns stay flat

        Returns value between -1 and 1.
        """
        # Heuristic: Arbitrum DeFi vs Base Treasury are negatively correlated
        # Real implementation would calculate from historical returns

        arb_allocation = 0.0
        base_allocation = 0.0

        for chain_id, data in portfolio['chains'].items():
            if chain_id in [421614, 42161]:  # Arbitrum
                arb_allocation = data['percentage'] / 100.0
            elif chain_id in [84532, 8453]:  # Base
                base_allocation = data['percentage'] / 100.0

        # If both have allocations, assume moderate negative correlation
        if arb_allocation > 0 and base_allocation > 0:
            return -0.3  # Slightly negative correlation

        # If only one chain, no correlation benefit
        return 0.0

    def _calculate_diversification_benefit(
        self,
        chain_risks: Dict,
        chains_data: Dict,
        correlation: float
    ) -> float:
        """
        Calculate diversification benefit (risk reduction from multi-chain).

        Formula (simplified portfolio theory):
        diversification_benefit = undiversified_risk - diversified_risk
        """
        if len(chain_risks) < 2:
            return 0.0

        # Get individual risks
        risks = []
        weights = []

        for chain_id, risk_data in chain_risks.items():
            risks.append(risk_data['risk_score'])
            weights.append(chains_data[chain_id]['percentage'] / 100.0)

        # Undiversified risk (worst case: highest risk)
        undiversified_risk = max(risks)

        # Diversified risk (weighted average with correlation adjustment)
        weighted_risk = sum(r * w for r, w in zip(risks, weights))

        # Correlation benefit (negative correlation reduces risk)
        correlation_benefit = abs(correlation) * 20.0  # Up to 20 points reduction

        diversified_risk = weighted_risk - correlation_benefit

        return max(undiversified_risk - diversified_risk, 0.0)

    def _calculate_diversification_score(
        self,
        chains_data: Dict,
        total_value: float
    ) -> float:
        """
        Calculate diversification score (0-100).

        100 = perfectly balanced (50/50)
        0 = all on one chain
        """
        if total_value == 0 or len(chains_data) < 2:
            return 0.0

        # Get allocation percentages
        allocations = [data['percentage'] for data in chains_data.values()]

        # Perfect balance is equal weight across all chains
        perfect_weight = 100.0 / len(chains_data)

        # Calculate deviation from perfect balance
        deviations = [abs(alloc - perfect_weight) for alloc in allocations]
        avg_deviation = sum(deviations) / len(deviations)

        # Convert to score (0-100)
        max_deviation = perfect_weight  # Worst case: 100% on one chain
        score = 100.0 * (1.0 - avg_deviation / max_deviation)

        return score

    # ========================================================================
    # Price Queries
    # ========================================================================

    async def get_asset_price_multi_chain(
        self,
        asset_symbol: str
    ) -> Dict[int, AssetPrice]:
        """Get price for an asset across all chains"""
        results = {}

        for chain_id, adapter in self.adapters.items():
            try:
                price = await adapter.get_asset_price(asset_symbol)
                results[chain_id] = price
            except Exception as e:
                print(f"Could not get {asset_symbol} price on {chain_id}: {e}")

        return results

    # ========================================================================
    # Chain Selection (AI Recommendations)
    # ========================================================================

    def recommend_chain_for_action(
        self,
        action_type: str,
        user_risk_tolerance: str = 'balanced'
    ) -> Tuple[int, str]:
        """
        Recommend which chain to use for a given action.

        Args:
            action_type: 'deposit', 'hedge', 'high_yield', 'safe_haven'
            user_risk_tolerance: 'conservative', 'balanced', 'aggressive'

        Returns:
            (chain_id, reason)
        """
        # Map action types to chains
        if action_type == 'safe_haven' or action_type == 'hedge':
            # Always recommend Base for hedging/safe haven
            base_adapter = self._get_base_adapter()
            return (
                base_adapter.chain_id,
                f"Base is recommended for {action_type}: stable US Treasury bonds with 4.5% APY and zero liquidation risk"
            )

        elif action_type == 'high_yield':
            # Recommend based on risk tolerance
            if user_risk_tolerance == 'conservative':
                base_adapter = self._get_base_adapter()
                return (
                    base_adapter.chain_id,
                    "Base is recommended: Conservative user + high yield → US Treasury bonds (4.5% APY, no risk)"
                )
            else:
                arb_adapter = self._get_arbitrum_adapter()
                return (
                    arb_adapter.chain_id,
                    "Arbitrum is recommended: Seek high yield → GMX perpetuals or Aave lending (10-50% APY, higher risk)"
                )

        elif action_type == 'deposit':
            # Balanced approach
            if user_risk_tolerance == 'aggressive':
                arb_adapter = self._get_arbitrum_adapter()
                return (arb_adapter.chain_id, "Aggressive user → Arbitrum DeFi for maximum returns")
            elif user_risk_tolerance == 'conservative':
                base_adapter = self._get_base_adapter()
                return (base_adapter.chain_id, "Conservative user → Base Treasury for stability")
            else:
                # Recommend diversification
                return (
                    0,  # Special code for "both"
                    "Balanced user → Recommend 50% Arbitrum + 50% Base for diversification"
                )

        # Default
        base_adapter = self._get_base_adapter()
        return (base_adapter.chain_id, "Default recommendation: Base (safest option)")

    def _get_arbitrum_adapter(self) -> ChainAdapter:
        """Get Arbitrum adapter"""
        for chain_id, adapter in self.adapters.items():
            if chain_id in [421614, 42161]:
                return adapter
        raise ValueError("Arbitrum adapter not found")

    def _get_base_adapter(self) -> ChainAdapter:
        """Get Base adapter"""
        for chain_id, adapter in self.adapters.items():
            if chain_id in [84532, 8453]:
                return adapter
        raise ValueError("Base adapter not found")

    # ========================================================================
    # Utility Methods
    # ========================================================================

    def get_adapter(self, chain_id: int) -> Optional[ChainAdapter]:
        """Get adapter for specific chain"""
        return self.adapters.get(chain_id)

    def get_all_chains(self) -> List[Dict[str, Any]]:
        """Get info about all supported chains"""
        return [
            {
                'chain_id': adapter.chain_id,
                'chain_name': adapter.chain_name,
                'characteristics': adapter.get_chain_characteristics()
            }
            for adapter in self.adapters.values()
        ]

    async def health_check_all_chains(self) -> Dict[int, bool]:
        """Check health of all chains"""
        results = {}
        for chain_id, adapter in self.adapters.items():
            results[chain_id] = await adapter.is_chain_healthy()
        return results

    def __repr__(self) -> str:
        chain_names = [adapter.chain_name for adapter in self.adapters.values()]
        return f"<MultiChainManager chains={chain_names}>"
