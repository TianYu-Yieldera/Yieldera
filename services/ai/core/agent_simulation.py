"""
Agent-Based Simulation Module
Gauntlet-standard market simulation with 10,000+ agents
Phase 2: Core Implementation
"""

import numpy as np
import asyncio
from enum import Enum
from typing import Dict, List, Optional, Tuple, Any
from dataclasses import dataclass, field
import logging
from datetime import datetime, timedelta
import random
from collections import defaultdict
import json

logger = logging.getLogger(__name__)

# ============================================================
# AGENT TYPES & MODELS
# ============================================================

class AgentType(Enum):
    """Gauntlet-standard agent classifications"""
    WHALE = "whale"                    # $1M-$100M positions
    ARBITRAGEUR = "arbitrageur"        # High-frequency traders
    LIQUIDATOR = "liquidator"          # Professional liquidators
    YIELD_FARMER = "yield_farmer"      # Yield optimizers
    RETAIL = "retail"                  # $1k-$100k positions
    PROTOCOL = "protocol"              # Protocol treasuries
    MARKET_MAKER = "market_maker"      # Liquidity providers
    INSTITUTIONAL = "institutional"    # Large institutions

@dataclass
class Agent:
    """Individual agent in the simulation"""
    id: str
    agent_type: AgentType
    wallet_balance: float
    positions: Dict[str, Any] = field(default_factory=dict)
    risk_appetite: float = 0.5  # 0=conservative, 1=aggressive
    leverage_preference: float = 1.0
    reaction_speed: float = 1.0  # How quickly agent reacts to market
    liquidation_history: List[Dict] = field(default_factory=list)
    pnl_history: List[float] = field(default_factory=list)
    strategy_params: Dict[str, Any] = field(default_factory=dict)

    def __post_init__(self):
        # Set agent-specific parameters based on type
        self._initialize_strategy()

    def _initialize_strategy(self):
        """Initialize strategy parameters based on agent type"""
        if self.agent_type == AgentType.WHALE:
            self.wallet_balance = np.random.uniform(1e6, 1e8)
            self.risk_appetite = np.random.uniform(0.3, 0.6)
            self.leverage_preference = np.random.uniform(1.0, 2.0)
            self.reaction_speed = np.random.uniform(0.5, 0.8)

        elif self.agent_type == AgentType.ARBITRAGEUR:
            self.wallet_balance = np.random.uniform(1e5, 1e6)
            self.risk_appetite = np.random.uniform(0.6, 0.9)
            self.leverage_preference = np.random.uniform(2.0, 5.0)
            self.reaction_speed = np.random.uniform(0.9, 1.0)

        elif self.agent_type == AgentType.LIQUIDATOR:
            self.wallet_balance = np.random.uniform(1e5, 5e6)
            self.risk_appetite = np.random.uniform(0.4, 0.7)
            self.leverage_preference = 1.0
            self.reaction_speed = 1.0
            self.strategy_params['min_profit'] = 0.05  # 5% minimum profit

        elif self.agent_type == AgentType.YIELD_FARMER:
            self.wallet_balance = np.random.uniform(1e4, 5e5)
            self.risk_appetite = np.random.uniform(0.5, 0.8)
            self.leverage_preference = np.random.uniform(1.5, 3.0)
            self.reaction_speed = np.random.uniform(0.3, 0.7)

        elif self.agent_type == AgentType.RETAIL:
            self.wallet_balance = np.random.uniform(1e3, 1e5)
            self.risk_appetite = np.random.uniform(0.2, 0.9)
            self.leverage_preference = np.random.uniform(1.0, 2.0)
            self.reaction_speed = np.random.uniform(0.1, 0.5)

        elif self.agent_type == AgentType.PROTOCOL:
            self.wallet_balance = np.random.uniform(1e6, 1e8)
            self.risk_appetite = np.random.uniform(0.1, 0.3)
            self.leverage_preference = 1.0
            self.reaction_speed = np.random.uniform(0.2, 0.4)

        elif self.agent_type == AgentType.MARKET_MAKER:
            self.wallet_balance = np.random.uniform(5e5, 1e7)
            self.risk_appetite = np.random.uniform(0.3, 0.5)
            self.leverage_preference = np.random.uniform(1.0, 1.5)
            self.reaction_speed = np.random.uniform(0.8, 1.0)
            self.strategy_params['spread'] = np.random.uniform(0.001, 0.005)

        elif self.agent_type == AgentType.INSTITUTIONAL:
            self.wallet_balance = np.random.uniform(1e7, 1e9)
            self.risk_appetite = np.random.uniform(0.2, 0.4)
            self.leverage_preference = np.random.uniform(1.0, 1.5)
            self.reaction_speed = np.random.uniform(0.3, 0.6)

@dataclass
class MarketState:
    """Current state of the simulated market"""
    timestamp: datetime
    prices: Dict[str, float]
    volatility: Dict[str, float]
    liquidity: Dict[str, float]
    total_supply: Dict[str, float]
    total_borrowed: Dict[str, float]
    utilization_rate: Dict[str, float]
    liquidations_24h: int = 0
    volume_24h: float = 0
    gas_price: float = 50  # Gwei

@dataclass
class SimulationResult:
    """Results from a simulation run"""
    scenario_name: str
    time_steps: int
    num_agents: int
    market_states: List[MarketState]
    liquidation_events: List[Dict]
    agent_pnl: Dict[str, float]
    cascade_analysis: Dict[str, Any]
    risk_metrics: Dict[str, float]
    confidence_interval: Tuple[float, float]
    execution_time: float

# ============================================================
# MAIN SIMULATION ENGINE
# ============================================================

class GauntletSimulator:
    """
    Gauntlet-level market simulator
    Simulates market dynamics with thousands of interacting agents
    """

    def __init__(self, num_agents: int = 10000, seed: Optional[int] = None):
        """
        Initialize simulator

        Args:
            num_agents: Number of agents to simulate
            seed: Random seed for reproducibility
        """
        if seed:
            np.random.seed(seed)
            random.seed(seed)

        self.num_agents = num_agents
        self.agents = []
        self.market_state = None
        self.simulation_history = []
        self.liquidation_history = []

        # Agent distribution (Gauntlet-calibrated)
        self.agent_distribution = {
            AgentType.WHALE: 0.01,          # 1% - 100 whales
            AgentType.ARBITRAGEUR: 0.05,    # 5% - 500 arbitrageurs
            AgentType.LIQUIDATOR: 0.02,     # 2% - 200 liquidators
            AgentType.YIELD_FARMER: 0.30,   # 30% - 3000 farmers
            AgentType.RETAIL: 0.55,         # 55% - 5500 retail
            AgentType.PROTOCOL: 0.02,       # 2% - 200 protocols
            AgentType.MARKET_MAKER: 0.03,   # 3% - 300 market makers
            AgentType.INSTITUTIONAL: 0.02   # 2% - 200 institutions
        }

        # Initialize agents
        self._initialize_agents()

        # Initialize market state
        self._initialize_market()

        logger.info(f"Initialized simulator with {num_agents} agents")

    def _initialize_agents(self):
        """Create and initialize all agents"""
        agent_id = 0

        for agent_type, proportion in self.agent_distribution.items():
            count = int(self.num_agents * proportion)
            for _ in range(count):
                agent = Agent(
                    id=f"agent_{agent_id}",
                    agent_type=agent_type,
                    wallet_balance=0  # Will be set in __post_init__
                )
                self.agents.append(agent)
                agent_id += 1

        # Ensure we have exactly num_agents
        while len(self.agents) < self.num_agents:
            agent = Agent(
                id=f"agent_{agent_id}",
                agent_type=AgentType.RETAIL,
                wallet_balance=0
            )
            self.agents.append(agent)
            agent_id += 1

    def _initialize_market(self):
        """Initialize market state"""
        self.market_state = MarketState(
            timestamp=datetime.utcnow(),
            prices={
                'ETH': 2000.0,
                'BTC': 40000.0,
                'USDC': 1.0,
                'USDT': 1.0
            },
            volatility={
                'ETH': 0.03,
                'BTC': 0.025,
                'USDC': 0.001,
                'USDT': 0.001
            },
            liquidity={
                'ETH': 1e8,
                'BTC': 5e7,
                'USDC': 1e9,
                'USDT': 1e9
            },
            total_supply={
                'ETH': 1e7,
                'BTC': 1e6,
                'USDC': 1e9,
                'USDT': 1e9
            },
            total_borrowed={
                'ETH': 5e6,
                'BTC': 5e5,
                'USDC': 5e8,
                'USDT': 5e8
            },
            utilization_rate={
                'ETH': 0.5,
                'BTC': 0.5,
                'USDC': 0.5,
                'USDT': 0.5
            }
        )

    # ============================================================
    # SIMULATION EXECUTION
    # ============================================================

    async def run_comprehensive_simulation(
        self,
        scenarios: List[Dict],
        time_steps: int = 1000
    ) -> Dict[str, SimulationResult]:
        """
        Run comprehensive multi-scenario simulation

        Args:
            scenarios: List of scenarios to simulate
            time_steps: Number of time steps per scenario

        Returns:
            Dictionary of results by scenario name
        """
        results = {}

        for scenario in scenarios:
            logger.info(f"Running scenario: {scenario['name']}")
            result = await self.simulate_scenario(
                scenario_config=scenario,
                time_steps=time_steps
            )
            results[scenario['name']] = result

        # Run comparative analysis
        analysis = self._comparative_analysis(results)

        return {
            'results': results,
            'analysis': analysis,
            'recommendations': self._generate_recommendations(analysis)
        }

    async def simulate_scenario(
        self,
        scenario_config: Dict,
        time_steps: int = 1000
    ) -> SimulationResult:
        """
        Simulate a specific market scenario

        Args:
            scenario_config: Scenario configuration
            time_steps: Number of simulation steps

        Returns:
            SimulationResult with detailed metrics
        """
        start_time = datetime.utcnow()

        # Reset market state for scenario
        self._apply_scenario(scenario_config)

        # Track states and events
        market_states = []
        liquidation_events = []

        # Run simulation steps
        for step in range(time_steps):
            # 1. Update market conditions
            self._update_market_conditions(step, scenario_config)

            # 2. Agents make decisions
            agent_actions = self._agent_decision_phase()

            # 3. Execute actions and update market
            self._execute_actions(agent_actions)

            # 4. Check for liquidations
            liquidations = self._check_liquidations()
            if liquidations:
                liquidation_events.extend(liquidations)

                # 5. Process liquidation cascade
                cascade = self._process_liquidation_cascade(liquidations)
                if cascade['rounds'] > 1:
                    logger.warning(f"Liquidation cascade detected: {cascade['rounds']} rounds")

            # 6. Record state
            market_states.append(self._snapshot_market_state())

            # 7. Update agent states
            self._update_agent_states()

        # Calculate final metrics
        agent_pnl = self._calculate_agent_pnl()
        cascade_analysis = self._analyze_cascades(liquidation_events)
        risk_metrics = self._calculate_risk_metrics(market_states, liquidation_events)
        confidence = self._calculate_confidence_interval(risk_metrics)

        execution_time = (datetime.utcnow() - start_time).total_seconds()

        return SimulationResult(
            scenario_name=scenario_config['name'],
            time_steps=time_steps,
            num_agents=self.num_agents,
            market_states=market_states,
            liquidation_events=liquidation_events,
            agent_pnl=agent_pnl,
            cascade_analysis=cascade_analysis,
            risk_metrics=risk_metrics,
            confidence_interval=confidence,
            execution_time=execution_time
        )

    # ============================================================
    # MARKET DYNAMICS
    # ============================================================

    def _apply_scenario(self, scenario: Dict):
        """Apply scenario parameters to market"""
        if 'market_shock' in scenario:
            shock = scenario['market_shock']
            for asset in self.market_state.prices:
                self.market_state.prices[asset] *= (1 + shock)

        if 'volatility_multiplier' in scenario:
            mult = scenario['volatility_multiplier']
            for asset in self.market_state.volatility:
                self.market_state.volatility[asset] *= mult

        if 'liquidity_shock' in scenario:
            shock = scenario['liquidity_shock']
            for asset in self.market_state.liquidity:
                self.market_state.liquidity[asset] *= (1 + shock)

    def _update_market_conditions(self, step: int, scenario: Dict):
        """Update market conditions for current step"""
        # Price evolution (Geometric Brownian Motion)
        for asset in self.market_state.prices:
            volatility = self.market_state.volatility[asset]
            drift = scenario.get('drift', 0.0)

            # GBM: dS = μS dt + σS dW
            dt = 1 / 365  # Daily steps
            dW = np.random.normal(0, np.sqrt(dt))

            price_change = drift * dt + volatility * dW
            self.market_state.prices[asset] *= (1 + price_change)

        # Update utilization based on agent actions
        for asset in self.market_state.utilization_rate:
            utilization = self.market_state.total_borrowed[asset] / self.market_state.total_supply[asset]
            self.market_state.utilization_rate[asset] = min(utilization, 0.95)

        # Update gas price (simulate network congestion)
        if step % 10 == 0:
            congestion = len([a for a in self.agents if len(a.positions) > 0]) / len(self.agents)
            self.market_state.gas_price = 50 * (1 + congestion * 2)

    # ============================================================
    # AGENT BEHAVIOR
    # ============================================================

    def _agent_decision_phase(self) -> List[Dict]:
        """Agents make decisions based on market state"""
        actions = []

        for agent in self.agents:
            # Skip if agent has insufficient funds
            if agent.wallet_balance < 100:
                continue

            # Agent-specific decision logic
            if agent.agent_type == AgentType.ARBITRAGEUR:
                action = self._arbitrageur_strategy(agent)
            elif agent.agent_type == AgentType.LIQUIDATOR:
                action = self._liquidator_strategy(agent)
            elif agent.agent_type == AgentType.YIELD_FARMER:
                action = self._yield_farmer_strategy(agent)
            elif agent.agent_type == AgentType.WHALE:
                action = self._whale_strategy(agent)
            elif agent.agent_type == AgentType.MARKET_MAKER:
                action = self._market_maker_strategy(agent)
            else:
                action = self._default_strategy(agent)

            if action:
                actions.append(action)

        return actions

    def _arbitrageur_strategy(self, agent: Agent) -> Optional[Dict]:
        """Arbitrageur looks for price discrepancies"""
        # Simplified arbitrage logic
        eth_price = self.market_state.prices['ETH']

        # Random decision with bias based on market conditions
        if np.random.random() < agent.reaction_speed * 0.1:
            return {
                'agent_id': agent.id,
                'action': 'arbitrage',
                'asset': 'ETH',
                'amount': min(agent.wallet_balance * 0.5, self.market_state.liquidity['ETH'] * 0.01),
                'leverage': agent.leverage_preference
            }
        return None

    def _liquidator_strategy(self, agent: Agent) -> Optional[Dict]:
        """Liquidator monitors unhealthy positions"""
        # Check for liquidatable positions
        liquidatable = []

        for other_agent in self.agents:
            if other_agent.id == agent.id:
                continue

            for position_id, position in other_agent.positions.items():
                health_factor = self._calculate_health_factor(position)
                if health_factor < 1.05:  # Liquidation threshold
                    expected_profit = position['collateral_value'] * agent.strategy_params.get('min_profit', 0.05)
                    if expected_profit > self.market_state.gas_price * 10:  # Profitable after gas
                        liquidatable.append({
                            'target_agent': other_agent.id,
                            'position': position_id,
                            'expected_profit': expected_profit
                        })

        if liquidatable:
            # Choose most profitable liquidation
            best = max(liquidatable, key=lambda x: x['expected_profit'])
            return {
                'agent_id': agent.id,
                'action': 'liquidate',
                'target': best['target_agent'],
                'position': best['position']
            }
        return None

    def _yield_farmer_strategy(self, agent: Agent) -> Optional[Dict]:
        """Yield farmer seeks best APY opportunities"""
        # Find best yield opportunity
        best_asset = None
        best_apy = 0

        for asset in self.market_state.prices:
            # Calculate APY based on utilization
            utilization = self.market_state.utilization_rate[asset]
            base_rate = 0.02
            apy = base_rate + utilization * 0.15  # Simple rate curve

            if apy > best_apy:
                best_apy = apy
                best_asset = asset

        if best_asset and np.random.random() < agent.reaction_speed * 0.2:
            return {
                'agent_id': agent.id,
                'action': 'supply',
                'asset': best_asset,
                'amount': agent.wallet_balance * agent.risk_appetite * 0.3,
                'leverage': agent.leverage_preference
            }
        return None

    def _whale_strategy(self, agent: Agent) -> Optional[Dict]:
        """Whale makes large, market-moving trades"""
        # Whales move more slowly but with larger amounts
        if np.random.random() < agent.reaction_speed * 0.05:
            asset = np.random.choice(list(self.market_state.prices.keys()))

            # Large position that can move the market
            amount = agent.wallet_balance * agent.risk_appetite * 0.1

            # Check if position is too large for market
            if amount > self.market_state.liquidity[asset] * 0.05:
                amount = self.market_state.liquidity[asset] * 0.05

            action_type = np.random.choice(['supply', 'borrow'])

            return {
                'agent_id': agent.id,
                'action': action_type,
                'asset': asset,
                'amount': amount,
                'leverage': agent.leverage_preference
            }
        return None

    def _market_maker_strategy(self, agent: Agent) -> Optional[Dict]:
        """Market maker provides liquidity"""
        # Provide liquidity on both sides
        asset = np.random.choice(list(self.market_state.prices.keys()))
        spread = agent.strategy_params.get('spread', 0.002)

        if np.random.random() < agent.reaction_speed * 0.3:
            return {
                'agent_id': agent.id,
                'action': 'provide_liquidity',
                'asset': asset,
                'amount': agent.wallet_balance * 0.2,
                'spread': spread
            }
        return None

    def _default_strategy(self, agent: Agent) -> Optional[Dict]:
        """Default strategy for retail and other agents"""
        if np.random.random() < agent.reaction_speed * 0.1:
            action_type = np.random.choice(['supply', 'borrow'], p=[0.6, 0.4])
            asset = np.random.choice(list(self.market_state.prices.keys()))

            return {
                'agent_id': agent.id,
                'action': action_type,
                'asset': asset,
                'amount': agent.wallet_balance * agent.risk_appetite * 0.1,
                'leverage': min(agent.leverage_preference, 2.0)
            }
        return None

    # ============================================================
    # ACTION EXECUTION
    # ============================================================

    def _execute_actions(self, actions: List[Dict]):
        """Execute agent actions and update market"""
        for action in actions:
            agent = next(a for a in self.agents if a.id == action['agent_id'])

            if action['action'] == 'supply':
                self._execute_supply(agent, action)
            elif action['action'] == 'borrow':
                self._execute_borrow(agent, action)
            elif action['action'] == 'liquidate':
                self._execute_liquidation(agent, action)
            elif action['action'] == 'provide_liquidity':
                self._execute_liquidity_provision(agent, action)
            elif action['action'] == 'arbitrage':
                self._execute_arbitrage(agent, action)

    def _execute_supply(self, agent: Agent, action: Dict):
        """Execute supply action"""
        asset = action['asset']
        amount = action['amount']

        # Update market state
        self.market_state.total_supply[asset] += amount
        self.market_state.liquidity[asset] += amount

        # Update agent position
        position_id = f"{agent.id}_{asset}_supply_{len(agent.positions)}"
        agent.positions[position_id] = {
            'type': 'supply',
            'asset': asset,
            'amount': amount,
            'entry_price': self.market_state.prices[asset],
            'timestamp': self.market_state.timestamp
        }

        # Deduct from wallet
        agent.wallet_balance -= amount * self.market_state.prices[asset]

    def _execute_borrow(self, agent: Agent, action: Dict):
        """Execute borrow action"""
        asset = action['asset']
        amount = action['amount']
        leverage = action.get('leverage', 1.0)

        # Calculate collateral required
        collateral_factor = 1.5  # 150% collateralization
        collateral_required = amount * self.market_state.prices[asset] * collateral_factor / leverage

        if agent.wallet_balance < collateral_required:
            return  # Insufficient collateral

        # Update market state
        self.market_state.total_borrowed[asset] += amount
        self.market_state.liquidity[asset] -= amount

        # Update agent position
        position_id = f"{agent.id}_{asset}_borrow_{len(agent.positions)}"
        agent.positions[position_id] = {
            'type': 'borrow',
            'asset': asset,
            'amount': amount,
            'collateral': collateral_required,
            'entry_price': self.market_state.prices[asset],
            'leverage': leverage,
            'health_factor': collateral_factor,
            'timestamp': self.market_state.timestamp
        }

        # Add borrowed amount to wallet, deduct collateral
        agent.wallet_balance += amount * self.market_state.prices[asset] - collateral_required

    def _execute_liquidation(self, agent: Agent, action: Dict):
        """Execute liquidation"""
        target_agent = next(a for a in self.agents if a.id == action['target'])
        position = target_agent.positions.get(action['position'])

        if not position:
            return

        # Calculate liquidation bonus (usually 5-10%)
        liquidation_bonus = 0.05
        collateral_seized = position.get('collateral', 0) * (1 + liquidation_bonus)

        # Transfer collateral to liquidator
        agent.wallet_balance += collateral_seized

        # Remove position from target
        del target_agent.positions[action['position']]

        # Record liquidation
        target_agent.liquidation_history.append({
            'timestamp': self.market_state.timestamp,
            'position': action['position'],
            'liquidator': agent.id,
            'collateral_lost': collateral_seized
        })

        # Update market state
        self.market_state.liquidations_24h += 1

    def _execute_arbitrage(self, agent: Agent, action: Dict):
        """Execute arbitrage trade"""
        # Simplified arbitrage execution
        asset = action['asset']
        amount = action['amount']

        # Simulate profit from arbitrage (1-3% typically)
        profit_rate = np.random.uniform(0.01, 0.03)
        profit = amount * profit_rate

        agent.wallet_balance += profit
        agent.pnl_history.append(profit)

        # Add some market impact
        self.market_state.volume_24h += amount

    def _execute_liquidity_provision(self, agent: Agent, action: Dict):
        """Execute liquidity provision"""
        asset = action['asset']
        amount = action['amount']

        # Add to liquidity pool
        self.market_state.liquidity[asset] += amount

        # Create LP position
        position_id = f"{agent.id}_{asset}_lp_{len(agent.positions)}"
        agent.positions[position_id] = {
            'type': 'liquidity',
            'asset': asset,
            'amount': amount,
            'entry_price': self.market_state.prices[asset],
            'timestamp': self.market_state.timestamp
        }

        agent.wallet_balance -= amount * self.market_state.prices[asset]

    # ============================================================
    # LIQUIDATION CASCADE ANALYSIS
    # ============================================================

    def _check_liquidations(self) -> List[Dict]:
        """Check all positions for liquidation"""
        liquidations = []

        for agent in self.agents:
            for position_id, position in list(agent.positions.items()):
                if position['type'] == 'borrow':
                    health_factor = self._calculate_health_factor(position)

                    if health_factor < 1.05:
                        liquidations.append({
                            'agent_id': agent.id,
                            'position_id': position_id,
                            'health_factor': health_factor,
                            'collateral': position.get('collateral', 0),
                            'debt': position['amount'] * self.market_state.prices[position['asset']]
                        })

        return liquidations

    def _calculate_health_factor(self, position: Dict) -> float:
        """Calculate health factor for a position"""
        if position['type'] != 'borrow':
            return float('inf')

        collateral_value = position.get('collateral', 0)
        debt_value = position['amount'] * self.market_state.prices[position['asset']]

        if debt_value == 0:
            return float('inf')

        return collateral_value / debt_value

    def _process_liquidation_cascade(self, initial_liquidations: List[Dict]) -> Dict:
        """
        Process liquidation cascade
        Core Gauntlet methodology for systemic risk analysis
        """
        cascade_rounds = []
        remaining_liquidations = initial_liquidations.copy()
        total_losses = 0
        round_num = 0

        while remaining_liquidations and round_num < 10:
            round_num += 1

            # Calculate market impact of this round
            market_impact = self._calculate_market_impact(remaining_liquidations)

            # Apply market impact to prices
            for asset in self.market_state.prices:
                self.market_state.prices[asset] *= (1 - market_impact)

            # Find next round of liquidations
            next_round_liquidations = []
            for agent in self.agents:
                # Skip already liquidated agents
                if any(liq['agent_id'] == agent.id for liq in remaining_liquidations):
                    continue

                # Check if market impact triggers new liquidations
                for position_id, position in agent.positions.items():
                    if position['type'] == 'borrow':
                        health_factor = self._calculate_health_factor(position)
                        if health_factor < 1.05:
                            next_round_liquidations.append({
                                'agent_id': agent.id,
                                'position_id': position_id,
                                'health_factor': health_factor,
                                'collateral': position.get('collateral', 0)
                            })

            # Record round
            round_losses = sum(liq['collateral'] for liq in remaining_liquidations)
            total_losses += round_losses

            cascade_rounds.append({
                'round': round_num,
                'liquidations': len(remaining_liquidations),
                'losses': round_losses,
                'market_impact': market_impact,
                'cumulative_losses': total_losses
            })

            # Setup next round
            remaining_liquidations = next_round_liquidations

        return {
            'rounds': round_num,
            'total_liquidations': sum(r['liquidations'] for r in cascade_rounds),
            'total_losses': total_losses,
            'cascade_depth': round_num,
            'systemic_risk': self._calculate_systemic_risk(cascade_rounds, total_losses),
            'details': cascade_rounds
        }

    def _calculate_market_impact(self, liquidations: List[Dict]) -> float:
        """Calculate market impact from liquidations"""
        if not liquidations:
            return 0

        # Total value being liquidated
        total_liquidated = sum(liq['collateral'] for liq in liquidations)

        # Average market liquidity
        avg_liquidity = np.mean(list(self.market_state.liquidity.values()))

        # Impact formula (simplified)
        # Real Gauntlet uses more sophisticated models
        impact = min(total_liquidated / (avg_liquidity * 10), 0.1)  # Max 10% impact

        return impact

    def _calculate_systemic_risk(self, cascade_rounds: List[Dict], total_losses: float) -> float:
        """
        Calculate systemic risk score
        Gauntlet's proprietary metric
        """
        if not cascade_rounds:
            return 0

        # Factors:
        # 1. Number of cascade rounds (depth)
        depth_factor = min(len(cascade_rounds) / 5, 1.0)  # Normalized to 5 rounds

        # 2. Total losses relative to market
        total_market_value = sum(
            agent.wallet_balance + sum(
                pos.get('collateral', 0) for pos in agent.positions.values()
            )
            for agent in self.agents
        )
        loss_factor = min(total_losses / total_market_value, 1.0) if total_market_value > 0 else 0

        # 3. Contagion speed
        if len(cascade_rounds) > 1:
            contagion_factor = cascade_rounds[1]['liquidations'] / cascade_rounds[0]['liquidations']
        else:
            contagion_factor = 0

        # Weighted combination
        systemic_risk = (
            depth_factor * 0.4 +
            loss_factor * 0.4 +
            min(contagion_factor, 1.0) * 0.2
        )

        return systemic_risk

    # ============================================================
    # METRICS AND ANALYSIS
    # ============================================================

    def _snapshot_market_state(self) -> MarketState:
        """Take snapshot of current market state"""
        return MarketState(
            timestamp=self.market_state.timestamp,
            prices=self.market_state.prices.copy(),
            volatility=self.market_state.volatility.copy(),
            liquidity=self.market_state.liquidity.copy(),
            total_supply=self.market_state.total_supply.copy(),
            total_borrowed=self.market_state.total_borrowed.copy(),
            utilization_rate=self.market_state.utilization_rate.copy(),
            liquidations_24h=self.market_state.liquidations_24h,
            volume_24h=self.market_state.volume_24h,
            gas_price=self.market_state.gas_price
        )

    def _update_agent_states(self):
        """Update agent states after each step"""
        for agent in self.agents:
            # Update position values
            for position in agent.positions.values():
                if 'asset' in position:
                    current_price = self.market_state.prices[position['asset']]
                    entry_price = position.get('entry_price', current_price)

                    # Calculate unrealized PnL
                    if position['type'] == 'supply':
                        pnl = (current_price - entry_price) * position['amount']
                    elif position['type'] == 'borrow':
                        pnl = (entry_price - current_price) * position['amount']
                    else:
                        pnl = 0

                    position['unrealized_pnl'] = pnl
                    position['current_value'] = position['amount'] * current_price

    def _calculate_agent_pnl(self) -> Dict[str, float]:
        """Calculate P&L for all agents"""
        pnl_by_type = defaultdict(list)

        for agent in self.agents:
            total_pnl = sum(agent.pnl_history)

            # Add unrealized PnL
            for position in agent.positions.values():
                total_pnl += position.get('unrealized_pnl', 0)

            pnl_by_type[agent.agent_type.value].append(total_pnl)

        # Calculate statistics by agent type
        result = {}
        for agent_type, pnls in pnl_by_type.items():
            result[agent_type] = {
                'mean': np.mean(pnls),
                'std': np.std(pnls),
                'min': np.min(pnls),
                'max': np.max(pnls),
                'total': np.sum(pnls)
            }

        return result

    def _analyze_cascades(self, liquidation_events: List[Dict]) -> Dict:
        """Analyze liquidation cascade patterns"""
        if not liquidation_events:
            return {
                'cascade_detected': False,
                'max_cascade_depth': 0,
                'total_liquidations': 0
            }

        # Group liquidations by time proximity
        cascades = []
        current_cascade = []

        for event in liquidation_events:
            if not current_cascade:
                current_cascade.append(event)
            else:
                # If events are close in time, consider them part of same cascade
                # This is simplified; real implementation would be more sophisticated
                current_cascade.append(event)

                # Check if cascade is complete (simplified)
                if len(current_cascade) > 10:
                    cascades.append(current_cascade)
                    current_cascade = []

        if current_cascade:
            cascades.append(current_cascade)

        return {
            'cascade_detected': len(cascades) > 0 and any(len(c) > 3 for c in cascades),
            'num_cascades': len(cascades),
            'max_cascade_depth': max(len(c) for c in cascades) if cascades else 0,
            'total_liquidations': len(liquidation_events),
            'cascade_patterns': [
                {
                    'size': len(cascade),
                    'agents_affected': len(set(e['agent_id'] for e in cascade))
                }
                for cascade in cascades
            ]
        }

    def _calculate_risk_metrics(
        self,
        market_states: List[MarketState],
        liquidation_events: List[Dict]
    ) -> Dict[str, float]:
        """Calculate comprehensive risk metrics"""

        # Price volatility
        price_series = {
            asset: [state.prices[asset] for state in market_states]
            for asset in market_states[0].prices
        }

        volatilities = {}
        for asset, prices in price_series.items():
            returns = np.diff(prices) / prices[:-1]
            volatilities[asset] = np.std(returns) * np.sqrt(365)  # Annualized

        # Market metrics
        avg_utilization = np.mean([
            np.mean(list(state.utilization_rate.values()))
            for state in market_states
        ])

        # Liquidation metrics
        liquidation_rate = len(liquidation_events) / (len(market_states) * self.num_agents)

        # Systemic risk (simplified)
        systemic_risk = min(liquidation_rate * 100, 1.0)  # Normalized

        return {
            'avg_volatility': np.mean(list(volatilities.values())),
            'max_volatility': np.max(list(volatilities.values())),
            'avg_utilization': avg_utilization,
            'liquidation_rate': liquidation_rate,
            'systemic_risk': systemic_risk,
            'total_liquidations': len(liquidation_events),
            'volatilities_by_asset': volatilities
        }

    def _calculate_confidence_interval(
        self,
        risk_metrics: Dict[str, float]
    ) -> Tuple[float, float]:
        """Calculate confidence interval for risk metrics"""
        # Simplified confidence calculation
        # Real implementation would use bootstrap or analytical methods

        central_estimate = risk_metrics['systemic_risk']

        # Width based on data quality and variance
        if risk_metrics['total_liquidations'] > 100:
            width = 0.05
        elif risk_metrics['total_liquidations'] > 50:
            width = 0.10
        else:
            width = 0.20

        lower = max(0, central_estimate - width)
        upper = min(1, central_estimate + width)

        return (lower, upper)

    def _comparative_analysis(self, results: Dict[str, SimulationResult]) -> Dict:
        """Compare results across scenarios"""
        comparison = {
            'scenarios': list(results.keys()),
            'risk_comparison': {},
            'liquidation_comparison': {},
            'pnl_comparison': {}
        }

        for scenario_name, result in results.items():
            comparison['risk_comparison'][scenario_name] = result.risk_metrics['systemic_risk']
            comparison['liquidation_comparison'][scenario_name] = result.risk_metrics['total_liquidations']

            # Average PnL across all agent types
            avg_pnl = np.mean([
                stats['mean']
                for stats in result.agent_pnl.values()
            ])
            comparison['pnl_comparison'][scenario_name] = avg_pnl

        # Identify worst-case scenario
        comparison['worst_case'] = max(
            comparison['risk_comparison'].items(),
            key=lambda x: x[1]
        )[0]

        return comparison

    def _generate_recommendations(self, analysis: Dict) -> List[str]:
        """Generate recommendations based on simulation results"""
        recommendations = []

        # Check systemic risk levels
        worst_risk = max(analysis['risk_comparison'].values())
        if worst_risk > 0.7:
            recommendations.append(
                "CRITICAL: High systemic risk detected. "
                "Immediately reduce leverage limits and increase collateral requirements."
            )
        elif worst_risk > 0.5:
            recommendations.append(
                "WARNING: Moderate systemic risk. "
                "Consider tightening risk parameters gradually."
            )

        # Check liquidation rates
        max_liquidations = max(analysis['liquidation_comparison'].values())
        if max_liquidations > 100:
            recommendations.append(
                "High liquidation activity expected. "
                "Increase liquidation incentives and ensure sufficient liquidator capacity."
            )

        # Check PnL distribution
        pnl_variance = np.std(list(analysis['pnl_comparison'].values()))
        if pnl_variance > 10000:
            recommendations.append(
                "High variance in outcomes across scenarios. "
                "Implement robust hedging strategies."
            )

        # Scenario-specific recommendations
        if analysis['worst_case']:
            recommendations.append(
                f"Worst-case scenario: {analysis['worst_case']}. "
                f"Develop specific contingency plans for this scenario."
            )

        # General recommendations
        recommendations.extend([
            "Maintain real-time monitoring of all risk metrics",
            "Ensure adequate capital reserves for stress scenarios",
            "Regular parameter optimization based on market conditions"
        ])

        return recommendations