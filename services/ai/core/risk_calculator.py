"""
Risk Calculator Module
Implements risk scoring and assessment algorithms
Phase 1: Core Implementation
"""

import numpy as np
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass
from datetime import datetime, timedelta
import logging

logger = logging.getLogger(__name__)

# ============================================================
# RISK MODELS
# ============================================================

@dataclass
class RiskMetrics:
    """Comprehensive risk metrics for a position or user"""
    overall_risk_score: float  # 0-100
    liquidation_probability: float  # 0-1
    market_risk: float  # 0-100
    liquidity_risk: float  # 0-100
    credit_risk: float  # 0-100
    operational_risk: float  # 0-100
    value_at_risk: float  # VaR in USD
    expected_shortfall: float  # CVaR in USD
    max_drawdown: float  # Percentage
    sharpe_ratio: float
    health_factor: float
    leverage: float
    confidence_level: float  # Model confidence 0-1

@dataclass
class PositionRisk:
    """Risk assessment for a specific position"""
    position_id: str
    user_address: str
    protocol: str
    risk_metrics: RiskMetrics
    recommended_actions: List[str]
    alert_level: str  # 'safe', 'warning', 'danger', 'critical'

# ============================================================
# RISK CALCULATOR CLASS
# ============================================================

class RiskCalculator:
    """
    Advanced risk calculation engine
    Implements multiple risk models and scoring algorithms
    """

    def __init__(self, config: Optional[Dict] = None):
        """
        Initialize risk calculator

        Args:
            config: Configuration parameters for risk models
        """
        self.config = config or self._default_config()
        self.weights = self.config['risk_weights']
        self.thresholds = self.config['risk_thresholds']

    def _default_config(self) -> Dict:
        """Default configuration for risk models"""
        return {
            'risk_weights': {
                'market_risk': 0.30,
                'liquidity_risk': 0.25,
                'credit_risk': 0.25,
                'operational_risk': 0.20
            },
            'risk_thresholds': {
                'safe': 30,
                'warning': 50,
                'danger': 70,
                'critical': 85
            },
            'var_confidence': 0.99,
            'var_horizon_days': 1,
            'liquidation_thresholds': {
                'aave': 1.05,
                'compound': 1.08,
                'gmx': 1.00
            }
        }

    # ============================================================
    # MAIN CALCULATION METHODS
    # ============================================================

    def calculate_position_risk(
        self,
        position_data: Dict,
        market_data: Dict,
        historical_data: Optional[np.ndarray] = None
    ) -> PositionRisk:
        """
        Calculate comprehensive risk metrics for a position

        Args:
            position_data: Position details
            market_data: Current market conditions
            historical_data: Historical price data for VaR calculation

        Returns:
            PositionRisk with all metrics and recommendations
        """
        try:
            # Extract position details
            collateral_value = position_data.get('collateral_value_usd', 0)
            debt_value = position_data.get('debt_value_usd', 0)
            health_factor = position_data.get('health_factor', 2.0)
            protocol = position_data.get('protocol', 'unknown')

            # Calculate individual risk components
            market_risk = self._calculate_market_risk(position_data, market_data)
            liquidity_risk = self._calculate_liquidity_risk(position_data, market_data)
            credit_risk = self._calculate_credit_risk(position_data)
            operational_risk = self._calculate_operational_risk(position_data)

            # Calculate overall risk score
            overall_risk = self._calculate_overall_risk({
                'market_risk': market_risk,
                'liquidity_risk': liquidity_risk,
                'credit_risk': credit_risk,
                'operational_risk': operational_risk
            })

            # Calculate VaR and CVaR
            var, cvar = self._calculate_var_cvar(
                collateral_value,
                historical_data
            ) if historical_data is not None else (0, 0)

            # Calculate other metrics
            liquidation_prob = self._calculate_liquidation_probability(
                health_factor,
                protocol
            )
            leverage = debt_value / collateral_value if collateral_value > 0 else 0
            max_drawdown = self._calculate_max_drawdown(historical_data) if historical_data is not None else 0
            sharpe_ratio = self._calculate_sharpe_ratio(historical_data) if historical_data is not None else 0

            # Determine confidence level
            confidence = self._calculate_confidence_level(position_data, market_data)

            # Create risk metrics
            risk_metrics = RiskMetrics(
                overall_risk_score=overall_risk,
                liquidation_probability=liquidation_prob,
                market_risk=market_risk,
                liquidity_risk=liquidity_risk,
                credit_risk=credit_risk,
                operational_risk=operational_risk,
                value_at_risk=var,
                expected_shortfall=cvar,
                max_drawdown=max_drawdown,
                sharpe_ratio=sharpe_ratio,
                health_factor=health_factor,
                leverage=leverage,
                confidence_level=confidence
            )

            # Generate recommendations
            recommendations = self._generate_recommendations(risk_metrics, position_data)

            # Determine alert level
            alert_level = self._determine_alert_level(overall_risk)

            return PositionRisk(
                position_id=position_data.get('position_id', 'unknown'),
                user_address=position_data.get('user_address', 'unknown'),
                protocol=protocol,
                risk_metrics=risk_metrics,
                recommended_actions=recommendations,
                alert_level=alert_level
            )

        except Exception as e:
            logger.error(f"Failed to calculate position risk: {e}")
            raise

    def calculate_portfolio_risk(
        self,
        positions: List[Dict],
        market_data: Dict,
        correlation_matrix: Optional[np.ndarray] = None
    ) -> RiskMetrics:
        """
        Calculate risk metrics for entire portfolio

        Args:
            positions: List of positions
            market_data: Current market conditions
            correlation_matrix: Asset correlation matrix

        Returns:
            Portfolio-level RiskMetrics
        """
        if not positions:
            return self._empty_risk_metrics()

        # Calculate individual position risks
        position_risks = []
        for position in positions:
            risk = self.calculate_position_risk(position, market_data)
            position_risks.append(risk)

        # Aggregate metrics
        total_value = sum(p.get('collateral_value_usd', 0) for p in positions)
        total_debt = sum(p.get('debt_value_usd', 0) for p in positions)

        # Weight risks by position value
        weighted_risks = {
            'market_risk': 0,
            'liquidity_risk': 0,
            'credit_risk': 0,
            'operational_risk': 0
        }

        for i, risk in enumerate(position_risks):
            weight = positions[i].get('collateral_value_usd', 0) / total_value if total_value > 0 else 0
            weighted_risks['market_risk'] += risk.risk_metrics.market_risk * weight
            weighted_risks['liquidity_risk'] += risk.risk_metrics.liquidity_risk * weight
            weighted_risks['credit_risk'] += risk.risk_metrics.credit_risk * weight
            weighted_risks['operational_risk'] += risk.risk_metrics.operational_risk * weight

        # Calculate portfolio VaR with correlation
        portfolio_var = self._calculate_portfolio_var(
            positions,
            correlation_matrix
        ) if correlation_matrix is not None else sum(r.risk_metrics.value_at_risk for r in position_risks)

        # Calculate overall metrics
        overall_risk = self._calculate_overall_risk(weighted_risks)
        avg_health_factor = np.mean([p.get('health_factor', 2.0) for p in positions])
        portfolio_leverage = total_debt / total_value if total_value > 0 else 0

        return RiskMetrics(
            overall_risk_score=overall_risk,
            liquidation_probability=np.mean([r.risk_metrics.liquidation_probability for r in position_risks]),
            market_risk=weighted_risks['market_risk'],
            liquidity_risk=weighted_risks['liquidity_risk'],
            credit_risk=weighted_risks['credit_risk'],
            operational_risk=weighted_risks['operational_risk'],
            value_at_risk=portfolio_var,
            expected_shortfall=portfolio_var * 1.2,  # Approximate
            max_drawdown=max(r.risk_metrics.max_drawdown for r in position_risks),
            sharpe_ratio=np.mean([r.risk_metrics.sharpe_ratio for r in position_risks]),
            health_factor=avg_health_factor,
            leverage=portfolio_leverage,
            confidence_level=np.mean([r.risk_metrics.confidence_level for r in position_risks])
        )

    # ============================================================
    # RISK COMPONENT CALCULATIONS
    # ============================================================

    def _calculate_market_risk(self, position: Dict, market: Dict) -> float:
        """Calculate market risk component (0-100)"""
        try:
            # Factors: volatility, price trend, correlation
            volatility = market.get('volatility', {}).get(position.get('collateral_asset', ''), 0.1)
            price_change = market.get('price_change_24h', {}).get(position.get('collateral_asset', ''), 0)

            # Volatility contributes 60%
            volatility_risk = min(volatility * 100, 60)

            # Price trend contributes 40%
            trend_risk = min(abs(price_change) * 2, 40)

            return volatility_risk + trend_risk

        except Exception as e:
            logger.error(f"Error calculating market risk: {e}")
            return 50.0  # Default medium risk

    def _calculate_liquidity_risk(self, position: Dict, market: Dict) -> float:
        """Calculate liquidity risk component (0-100)"""
        try:
            # Factors: available liquidity, utilization rate, position size
            market_key = f"{position.get('protocol')}_{position.get('market', 'default')}"
            market_depth = market.get('market_depth', {}).get(market_key, {})

            available_liquidity = market_depth.get('available_liquidity', float('inf'))
            utilization_rate = market_depth.get('utilization_rate', 0.5)
            position_size = position.get('collateral_value_usd', 0)

            # Utilization contributes 40%
            utilization_risk = utilization_rate * 40

            # Relative position size contributes 60%
            if available_liquidity > 0:
                size_risk = min((position_size / available_liquidity) * 60, 60)
            else:
                size_risk = 30

            return utilization_risk + size_risk

        except Exception as e:
            logger.error(f"Error calculating liquidity risk: {e}")
            return 50.0

    def _calculate_credit_risk(self, position: Dict) -> float:
        """Calculate credit risk component (0-100)"""
        try:
            # Factors: health factor, LTV, collateral quality
            health_factor = position.get('health_factor', 2.0)
            ltv = position.get('ltv', 0.5)

            # Handle None values
            if health_factor is None:
                health_factor = 2.0
            if ltv is None:
                ltv = 0.5

            # Health factor contributes 70%
            if health_factor < 1.0:
                health_risk = 70
            elif health_factor < 1.5:
                health_risk = 50
            elif health_factor < 2.0:
                health_risk = 30
            else:
                health_risk = 10

            # LTV contributes 30%
            ltv_risk = ltv * 30

            return health_risk + ltv_risk

        except Exception as e:
            logger.error(f"Error calculating credit risk: {e}")
            return 50.0

    def _calculate_operational_risk(self, position: Dict) -> float:
        """Calculate operational risk component (0-100)"""
        try:
            # Factors: protocol risk, smart contract age, audit status
            protocol = position.get('protocol', 'unknown')

            # Protocol risk mapping (simplified)
            protocol_risks = {
                'aave': 10,
                'compound': 15,
                'gmx': 30,
                'unknown': 50
            }

            base_risk = protocol_risks.get(protocol, 50)

            # Add risk for new positions
            position_age_days = position.get('position_age_days', 0)
            if position_age_days < 1:
                base_risk += 20
            elif position_age_days < 7:
                base_risk += 10

            return min(base_risk, 100)

        except Exception as e:
            logger.error(f"Error calculating operational risk: {e}")
            return 50.0

    # ============================================================
    # VAR AND STATISTICAL CALCULATIONS
    # ============================================================

    def _calculate_var_cvar(
        self,
        portfolio_value: float,
        historical_returns: np.ndarray,
        confidence: float = 0.99
    ) -> Tuple[float, float]:
        """
        Calculate Value at Risk and Conditional Value at Risk

        Args:
            portfolio_value: Current portfolio value
            historical_returns: Historical return data
            confidence: Confidence level (default 99%)

        Returns:
            Tuple of (VaR, CVaR) in USD
        """
        try:
            if historical_returns is None or len(historical_returns) == 0:
                return 0, 0

            # Calculate returns if prices provided
            if historical_returns.ndim == 1:
                returns = np.diff(historical_returns) / historical_returns[:-1]
            else:
                returns = historical_returns

            # Calculate VaR
            var_percentile = (1 - confidence) * 100
            var_return = np.percentile(returns, var_percentile)
            var_usd = abs(var_return * portfolio_value)

            # Calculate CVaR (Expected Shortfall)
            tail_returns = returns[returns <= var_return]
            if len(tail_returns) > 0:
                cvar_return = np.mean(tail_returns)
                cvar_usd = abs(cvar_return * portfolio_value)
            else:
                cvar_usd = var_usd * 1.2  # Approximate

            return var_usd, cvar_usd

        except Exception as e:
            logger.error(f"Error calculating VaR/CVaR: {e}")
            return 0, 0

    def _calculate_portfolio_var(
        self,
        positions: List[Dict],
        correlation_matrix: np.ndarray
    ) -> float:
        """Calculate portfolio VaR considering correlations"""
        try:
            if correlation_matrix is None:
                # Simple sum if no correlation data
                return sum(p.get('var', 0) for p in positions)

            # Extract position VaRs
            vars = np.array([p.get('var', 0) for p in positions])

            # Calculate portfolio VaR with correlation
            portfolio_var = np.sqrt(vars @ correlation_matrix @ vars.T)

            return float(portfolio_var)

        except Exception as e:
            logger.error(f"Error calculating portfolio VaR: {e}")
            return sum(p.get('var', 0) for p in positions)

    def _calculate_max_drawdown(self, price_series: np.ndarray) -> float:
        """Calculate maximum drawdown from price series"""
        try:
            if price_series is None or len(price_series) < 2:
                return 0

            # Calculate running maximum
            running_max = np.maximum.accumulate(price_series)

            # Calculate drawdown
            drawdown = (price_series - running_max) / running_max

            return abs(float(np.min(drawdown)))

        except Exception as e:
            logger.error(f"Error calculating max drawdown: {e}")
            return 0

    def _calculate_sharpe_ratio(
        self,
        returns: np.ndarray,
        risk_free_rate: float = 0.02
    ) -> float:
        """Calculate Sharpe ratio"""
        try:
            if returns is None or len(returns) < 2:
                return 0

            # Calculate excess returns
            excess_returns = returns - risk_free_rate / 365  # Daily risk-free rate

            # Calculate Sharpe ratio
            if np.std(excess_returns) > 0:
                sharpe = np.mean(excess_returns) / np.std(excess_returns) * np.sqrt(365)
                return float(sharpe)
            else:
                return 0

        except Exception as e:
            logger.error(f"Error calculating Sharpe ratio: {e}")
            return 0

    # ============================================================
    # PROBABILITY AND SCORING
    # ============================================================

    def _calculate_liquidation_probability(
        self,
        health_factor: float,
        protocol: str
    ) -> float:
        """
        Calculate probability of liquidation

        Args:
            health_factor: Current health factor
            protocol: Protocol name

        Returns:
            Probability between 0 and 1
        """
        try:
            # Get liquidation threshold for protocol
            threshold = self.config['liquidation_thresholds'].get(protocol, 1.0)

            # Handle None or invalid health factor
            if health_factor is None or not isinstance(health_factor, (int, float)):
                logger.warning(f"Invalid health factor: {health_factor}, using default")
                return 0.5

            if health_factor <= threshold:
                return 1.0  # Already liquidatable

            # Exponential decay model
            buffer = health_factor - threshold
            probability = np.exp(-2 * buffer)  # Decay factor of 2

            return float(min(probability, 1.0))

        except Exception as e:
            logger.error(f"Error calculating liquidation probability: {e}")
            return 0.5

    def _calculate_overall_risk(self, risk_components: Dict[str, float]) -> float:
        """
        Calculate overall risk score from components

        Args:
            risk_components: Individual risk scores

        Returns:
            Overall risk score (0-100)
        """
        overall = 0
        for component, score in risk_components.items():
            weight = self.weights.get(component, 0.25)
            overall += score * weight

        return min(overall, 100)

    def _calculate_confidence_level(
        self,
        position: Dict,
        market: Dict
    ) -> float:
        """Calculate model confidence level"""
        # Factors affecting confidence:
        # - Data quality/completeness
        # - Position age
        # - Market conditions

        confidence = 1.0

        # Reduce confidence for incomplete data
        required_fields = ['health_factor', 'collateral_value_usd', 'debt_value_usd']
        for field in required_fields:
            if field not in position or position[field] is None:
                confidence *= 0.9

        # Reduce confidence for new positions
        position_age = position.get('position_age_days', 0)
        if position_age < 1:
            confidence *= 0.8
        elif position_age < 7:
            confidence *= 0.9

        # Reduce confidence in high volatility
        volatility = market.get('volatility', {}).get(position.get('collateral_asset', ''), 0)
        if volatility is not None and isinstance(volatility, (int, float)) and volatility > 0.5:
            confidence *= 0.85

        return max(confidence, 0.5)  # Minimum 50% confidence

    # ============================================================
    # RECOMMENDATIONS AND ALERTS
    # ============================================================

    def _generate_recommendations(
        self,
        risk_metrics: RiskMetrics,
        position: Dict
    ) -> List[str]:
        """Generate actionable recommendations based on risk metrics"""
        recommendations = []

        # Health factor recommendations
        if risk_metrics.health_factor < 1.2:
            recommendations.append("URGENT: Add collateral immediately to avoid liquidation")
            recommendations.append("Consider closing position to prevent losses")
        elif risk_metrics.health_factor < 1.5:
            recommendations.append("Add collateral to improve health factor above 1.5")
            recommendations.append("Reduce debt by repaying part of the loan")
        elif risk_metrics.health_factor < 2.0:
            recommendations.append("Monitor position closely - approaching risk zone")

        # Leverage recommendations
        if risk_metrics.leverage > 3:
            recommendations.append("High leverage detected - consider reducing exposure")
        elif risk_metrics.leverage > 5:
            recommendations.append("CRITICAL: Extremely high leverage - immediate action required")

        # VaR recommendations
        if risk_metrics.value_at_risk > position.get('collateral_value_usd', 0) * 0.1:
            recommendations.append(f"High VaR: potential 1-day loss of ${risk_metrics.value_at_risk:.2f}")

        # Market risk recommendations
        if risk_metrics.market_risk > 70:
            recommendations.append("High market volatility - consider hedging strategies")

        # Liquidity risk recommendations
        if risk_metrics.liquidity_risk > 70:
            recommendations.append("Low market liquidity - gradual position reduction recommended")

        return recommendations

    def _determine_alert_level(self, risk_score: float) -> str:
        """Determine alert level based on risk score"""
        if risk_score < self.thresholds['safe']:
            return 'safe'
        elif risk_score < self.thresholds['warning']:
            return 'warning'
        elif risk_score < self.thresholds['danger']:
            return 'danger'
        else:
            return 'critical'

    def _empty_risk_metrics(self) -> RiskMetrics:
        """Return empty risk metrics"""
        return RiskMetrics(
            overall_risk_score=0,
            liquidation_probability=0,
            market_risk=0,
            liquidity_risk=0,
            credit_risk=0,
            operational_risk=0,
            value_at_risk=0,
            expected_shortfall=0,
            max_drawdown=0,
            sharpe_ratio=0,
            health_factor=float('inf'),
            leverage=0,
            confidence_level=0
        )