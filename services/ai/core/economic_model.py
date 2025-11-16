"""
Economic Model & Risk Metrics Module
Implements Gauntlet-standard economic modeling and parameter optimization
Phase 2: Core Implementation
"""

import numpy as np
import pandas as pd
from typing import Dict, List, Tuple, Optional, Any
from dataclasses import dataclass
from datetime import datetime, timedelta
from scipy import stats, optimize
from scipy.optimize import differential_evolution, minimize
import logging

logger = logging.getLogger(__name__)

# ============================================================
# DATA MODELS
# ============================================================

@dataclass
class EconomicParameters:
    """Protocol economic parameters"""
    ltv: float  # Loan-to-Value ratio
    liquidation_threshold: float
    liquidation_penalty: float
    base_rate: float  # Interest rate model
    slope1: float
    slope2: float
    optimal_utilization: float = 0.8
    reserve_factor: float = 0.1

@dataclass
class RiskMetrics:
    """Comprehensive risk metrics"""
    var_95: float  # Value at Risk at 95% confidence
    var_99: float  # Value at Risk at 99% confidence
    cvar_95: float  # Conditional VaR at 95%
    cvar_99: float  # Conditional VaR at 99%
    max_drawdown: float
    sharpe_ratio: float
    sortino_ratio: float
    calmar_ratio: float
    beta: float
    alpha: float
    correlation_risk: float
    liquidity_risk: float
    concentration_risk: float

@dataclass
class OptimizationResult:
    """Parameter optimization result"""
    current_params: EconomicParameters
    optimal_params: EconomicParameters
    expected_improvement: Dict[str, float]
    risk_reduction: float
    revenue_increase: float
    capital_efficiency_gain: float
    implementation_risk: str
    confidence_score: float
    backtest_results: Dict[str, Any]

# ============================================================
# GAUNTLET ECONOMIC MODEL
# ============================================================

class GauntletEconomicModel:
    """
    Economic modeling engine based on Gauntlet methodology
    Implements advanced risk metrics and parameter optimization
    """

    def __init__(self, config: Optional[Dict] = None):
        """
        Initialize economic model

        Args:
            config: Model configuration
        """
        self.config = config or self._default_config()
        self.risk_free_rate = self.config.get('risk_free_rate', 0.02)
        self.confidence_levels = self.config.get('confidence_levels', [0.95, 0.99])

    def _default_config(self) -> Dict:
        """Default configuration"""
        return {
            'risk_free_rate': 0.02,  # 2% annual
            'confidence_levels': [0.95, 0.99],
            'optimization_constraints': {
                'min_ltv': 0.5,
                'max_ltv': 0.85,
                'min_liquidation_threshold': 0.6,
                'max_liquidation_threshold': 0.9,
                'min_liquidation_penalty': 0.05,
                'max_liquidation_penalty': 0.15,
                'min_base_rate': 0.0,
                'max_base_rate': 0.05,
                'min_slope1': 0.04,
                'max_slope1': 0.20,
                'min_slope2': 0.5,
                'max_slope2': 3.0
            }
        }

    # ============================================================
    # RISK METRICS CALCULATION
    # ============================================================

    def calculate_comprehensive_risk_metrics(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame,
        correlation_matrix: Optional[np.ndarray] = None
    ) -> RiskMetrics:
        """
        Calculate all risk metrics for portfolio

        Args:
            portfolio: Portfolio positions
            market_data: Historical market data
            correlation_matrix: Asset correlation matrix

        Returns:
            Comprehensive RiskMetrics
        """
        try:
            # Calculate VaR metrics
            var_95, cvar_95 = self._calculate_var_cvar(portfolio, market_data, 0.95)
            var_99, cvar_99 = self._calculate_var_cvar(portfolio, market_data, 0.99)

            # Calculate drawdown metrics
            max_drawdown = self._calculate_max_drawdown(portfolio, market_data)

            # Calculate risk-adjusted returns
            sharpe_ratio = self._calculate_sharpe_ratio(portfolio, market_data)
            sortino_ratio = self._calculate_sortino_ratio(portfolio, market_data)
            calmar_ratio = self._calculate_calmar_ratio(portfolio, market_data, max_drawdown)

            # Calculate market risk metrics
            beta = self._calculate_beta(portfolio, market_data)
            alpha = self._calculate_alpha(portfolio, market_data, beta)

            # Calculate portfolio-specific risks
            correlation_risk = self._calculate_correlation_risk(portfolio, correlation_matrix)
            liquidity_risk = self._calculate_liquidity_risk(portfolio, market_data)
            concentration_risk = self._calculate_concentration_risk(portfolio)

            return RiskMetrics(
                var_95=var_95,
                var_99=var_99,
                cvar_95=cvar_95,
                cvar_99=cvar_99,
                max_drawdown=max_drawdown,
                sharpe_ratio=sharpe_ratio,
                sortino_ratio=sortino_ratio,
                calmar_ratio=calmar_ratio,
                beta=beta,
                alpha=alpha,
                correlation_risk=correlation_risk,
                liquidity_risk=liquidity_risk,
                concentration_risk=concentration_risk
            )

        except Exception as e:
            logger.error(f"Failed to calculate risk metrics: {e}")
            raise

    def _calculate_var_cvar(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame,
        confidence: float
    ) -> Tuple[float, float]:
        """
        Calculate Value at Risk and Conditional Value at Risk

        Uses three methods and takes weighted average:
        1. Historical simulation
        2. Parametric (variance-covariance)
        3. Monte Carlo simulation
        """
        # Extract portfolio value and returns
        portfolio_value = sum(p.get('value_usd', 0) for p in portfolio.values())

        if portfolio_value == 0:
            return 0, 0

        # Method 1: Historical VaR
        historical_var = self._historical_var(portfolio, market_data, confidence)

        # Method 2: Parametric VaR
        parametric_var = self._parametric_var(portfolio, market_data, confidence)

        # Method 3: Monte Carlo VaR
        monte_carlo_var, monte_carlo_cvar = self._monte_carlo_var(
            portfolio, market_data, confidence, simulations=10000
        )

        # Weighted average (Gauntlet approach)
        var = (
            historical_var * 0.3 +
            parametric_var * 0.3 +
            monte_carlo_var * 0.4
        )

        # CVaR is typically 20-40% worse than VaR
        cvar = monte_carlo_cvar if monte_carlo_cvar > 0 else var * 1.3

        return var, cvar

    def _historical_var(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame,
        confidence: float
    ) -> float:
        """Historical simulation VaR"""
        # Calculate portfolio returns
        portfolio_returns = self._calculate_portfolio_returns(portfolio, market_data)

        if len(portfolio_returns) == 0:
            return 0

        # Calculate VaR as percentile
        var_percentile = (1 - confidence) * 100
        var = np.percentile(portfolio_returns, var_percentile)

        # Convert to positive loss value
        portfolio_value = sum(p.get('value_usd', 0) for p in portfolio.values())
        return abs(var * portfolio_value)

    def _parametric_var(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame,
        confidence: float
    ) -> float:
        """Parametric (variance-covariance) VaR"""
        # Calculate portfolio returns
        portfolio_returns = self._calculate_portfolio_returns(portfolio, market_data)

        if len(portfolio_returns) == 0:
            return 0

        # Calculate mean and standard deviation
        mean_return = np.mean(portfolio_returns)
        std_return = np.std(portfolio_returns)

        # Calculate VaR using normal distribution
        z_score = stats.norm.ppf(1 - confidence)
        var_return = mean_return + z_score * std_return

        # Convert to value
        portfolio_value = sum(p.get('value_usd', 0) for p in portfolio.values())
        return abs(var_return * portfolio_value)

    def _monte_carlo_var(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame,
        confidence: float,
        simulations: int = 10000
    ) -> Tuple[float, float]:
        """Monte Carlo simulation VaR and CVaR"""
        # Get portfolio statistics
        portfolio_returns = self._calculate_portfolio_returns(portfolio, market_data)

        if len(portfolio_returns) == 0:
            return 0, 0

        mean_return = np.mean(portfolio_returns)
        std_return = np.std(portfolio_returns)

        # Run Monte Carlo simulations
        simulated_returns = np.random.normal(mean_return, std_return, simulations)

        # Calculate VaR
        var_percentile = (1 - confidence) * 100
        var_return = np.percentile(simulated_returns, var_percentile)

        # Calculate CVaR (expected shortfall)
        cvar_return = np.mean(simulated_returns[simulated_returns <= var_return])

        # Convert to values
        portfolio_value = sum(p.get('value_usd', 0) for p in portfolio.values())
        var = abs(var_return * portfolio_value)
        cvar = abs(cvar_return * portfolio_value)

        return var, cvar

    def _calculate_portfolio_returns(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame
    ) -> np.ndarray:
        """Calculate historical portfolio returns"""
        returns = []

        # Simple return calculation
        # Real implementation would weight by position sizes
        for _, position in portfolio.items():
            asset = position.get('asset', 'ETH')
            if asset in market_data.columns:
                asset_returns = market_data[asset].pct_change().dropna()
                position_weight = position.get('value_usd', 0) / sum(
                    p.get('value_usd', 0) for p in portfolio.values()
                )
                if len(returns) == 0:
                    returns = asset_returns.values * position_weight
                else:
                    returns += asset_returns.values[:len(returns)] * position_weight

        return returns if isinstance(returns, np.ndarray) else np.array([])

    def _calculate_max_drawdown(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame
    ) -> float:
        """Calculate maximum drawdown"""
        # Calculate portfolio value series
        portfolio_returns = self._calculate_portfolio_returns(portfolio, market_data)

        if len(portfolio_returns) == 0:
            return 0

        # Calculate cumulative returns
        cumulative_returns = (1 + portfolio_returns).cumprod()

        # Calculate running maximum
        running_max = np.maximum.accumulate(cumulative_returns)

        # Calculate drawdown
        drawdown = (cumulative_returns - running_max) / running_max

        return abs(np.min(drawdown))

    def _calculate_sharpe_ratio(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame
    ) -> float:
        """Calculate Sharpe ratio"""
        portfolio_returns = self._calculate_portfolio_returns(portfolio, market_data)

        if len(portfolio_returns) < 2:
            return 0

        # Calculate excess returns
        excess_returns = portfolio_returns - self.risk_free_rate / 252  # Daily risk-free rate

        # Calculate Sharpe ratio
        if np.std(excess_returns) > 0:
            sharpe = np.mean(excess_returns) / np.std(excess_returns) * np.sqrt(252)
            return sharpe
        return 0

    def _calculate_sortino_ratio(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame
    ) -> float:
        """Calculate Sortino ratio (uses downside deviation)"""
        portfolio_returns = self._calculate_portfolio_returns(portfolio, market_data)

        if len(portfolio_returns) < 2:
            return 0

        # Calculate excess returns
        excess_returns = portfolio_returns - self.risk_free_rate / 252

        # Calculate downside deviation
        downside_returns = excess_returns[excess_returns < 0]
        if len(downside_returns) > 0:
            downside_deviation = np.std(downside_returns)
            if downside_deviation > 0:
                sortino = np.mean(excess_returns) / downside_deviation * np.sqrt(252)
                return sortino
        return 0

    def _calculate_calmar_ratio(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame,
        max_drawdown: float
    ) -> float:
        """Calculate Calmar ratio (return / max drawdown)"""
        if max_drawdown == 0:
            return 0

        portfolio_returns = self._calculate_portfolio_returns(portfolio, market_data)
        if len(portfolio_returns) == 0:
            return 0

        annual_return = np.mean(portfolio_returns) * 252
        return annual_return / max_drawdown

    def _calculate_beta(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame
    ) -> float:
        """Calculate portfolio beta relative to market"""
        portfolio_returns = self._calculate_portfolio_returns(portfolio, market_data)

        # Use first asset as market proxy (simplified)
        if 'ETH' in market_data.columns:
            market_returns = market_data['ETH'].pct_change().dropna()

            if len(portfolio_returns) > 0 and len(market_returns) > 0:
                # Align lengths
                min_len = min(len(portfolio_returns), len(market_returns))
                portfolio_returns = portfolio_returns[:min_len]
                market_returns = market_returns[:min_len]

                # Calculate beta
                covariance = np.cov(portfolio_returns, market_returns)[0, 1]
                market_variance = np.var(market_returns)

                if market_variance > 0:
                    return covariance / market_variance
        return 1.0

    def _calculate_alpha(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame,
        beta: float
    ) -> float:
        """Calculate portfolio alpha (excess return)"""
        portfolio_returns = self._calculate_portfolio_returns(portfolio, market_data)

        if len(portfolio_returns) == 0:
            return 0

        # CAPM alpha
        if 'ETH' in market_data.columns:
            market_returns = market_data['ETH'].pct_change().dropna()
            if len(market_returns) > 0:
                portfolio_return = np.mean(portfolio_returns) * 252
                market_return = np.mean(market_returns) * 252
                expected_return = self.risk_free_rate + beta * (market_return - self.risk_free_rate)
                return portfolio_return - expected_return
        return 0

    def _calculate_correlation_risk(
        self,
        portfolio: Dict,
        correlation_matrix: Optional[np.ndarray]
    ) -> float:
        """Calculate risk from asset correlations"""
        if correlation_matrix is None or len(portfolio) < 2:
            return 0

        # Calculate average correlation
        n = correlation_matrix.shape[0]
        if n > 1:
            # Extract upper triangle (excluding diagonal)
            upper_triangle = np.triu(correlation_matrix, k=1)
            avg_correlation = np.sum(upper_triangle) / (n * (n - 1) / 2)
            return abs(avg_correlation)
        return 0

    def _calculate_liquidity_risk(
        self,
        portfolio: Dict,
        market_data: pd.DataFrame
    ) -> float:
        """Calculate liquidity risk"""
        # Simplified: based on position size relative to market
        total_position = sum(p.get('value_usd', 0) for p in portfolio.values())

        # Assume market depth (would get from real data)
        market_depth = 1e8  # $100M

        # Liquidity risk increases non-linearly with position size
        liquidity_risk = min((total_position / market_depth) ** 2, 1.0)

        return liquidity_risk

    def _calculate_concentration_risk(self, portfolio: Dict) -> float:
        """Calculate concentration risk using Herfindahl index"""
        if not portfolio:
            return 0

        total_value = sum(p.get('value_usd', 0) for p in portfolio.values())
        if total_value == 0:
            return 0

        # Calculate Herfindahl-Hirschman Index
        hhi = sum(
            (p.get('value_usd', 0) / total_value) ** 2
            for p in portfolio.values()
        )

        # Normalize (HHI ranges from 1/n to 1)
        n = len(portfolio)
        if n > 1:
            normalized_hhi = (hhi - 1/n) / (1 - 1/n)
            return normalized_hhi
        return 1.0

    # ============================================================
    # PARAMETER OPTIMIZATION
    # ============================================================

    def optimize_protocol_parameters(
        self,
        current_params: EconomicParameters,
        market_conditions: Dict,
        historical_data: pd.DataFrame,
        constraints: Optional[Dict] = None
    ) -> OptimizationResult:
        """
        Optimize protocol parameters using Gauntlet methodology

        Args:
            current_params: Current parameter values
            market_conditions: Current market state
            historical_data: Historical market data
            constraints: Optimization constraints

        Returns:
            OptimizationResult with optimal parameters
        """
        try:
            logger.info("Starting parameter optimization...")

            # Use provided constraints or defaults
            constraints = constraints or self.config['optimization_constraints']

            # Define objective function
            def objective_function(params):
                return self._optimization_objective(
                    params,
                    market_conditions,
                    historical_data
                )

            # Set parameter bounds
            bounds = [
                (constraints['min_ltv'], constraints['max_ltv']),
                (constraints['min_liquidation_threshold'], constraints['max_liquidation_threshold']),
                (constraints['min_liquidation_penalty'], constraints['max_liquidation_penalty']),
                (constraints['min_base_rate'], constraints['max_base_rate']),
                (constraints['min_slope1'], constraints['max_slope1']),
                (constraints['min_slope2'], constraints['max_slope2'])
            ]

            # Current parameters as list
            current_params_list = [
                current_params.ltv,
                current_params.liquidation_threshold,
                current_params.liquidation_penalty,
                current_params.base_rate,
                current_params.slope1,
                current_params.slope2
            ]

            # Run differential evolution optimization
            result = differential_evolution(
                objective_function,
                bounds,
                seed=42,
                maxiter=100,
                popsize=15,
                mutation=(0.5, 1),
                recombination=0.7,
                strategy='best1bin',
                x0=current_params_list
            )

            # Extract optimal parameters
            optimal_params = EconomicParameters(
                ltv=result.x[0],
                liquidation_threshold=result.x[1],
                liquidation_penalty=result.x[2],
                base_rate=result.x[3],
                slope1=result.x[4],
                slope2=result.x[5]
            )

            # Run backtest
            backtest_results = self._backtest_parameters(
                optimal_params,
                historical_data,
                market_conditions
            )

            # Calculate expected improvements
            current_performance = self._evaluate_parameters(
                current_params,
                market_conditions,
                historical_data
            )

            optimal_performance = self._evaluate_parameters(
                optimal_params,
                market_conditions,
                historical_data
            )

            # Calculate improvements
            risk_reduction = (
                (current_performance['risk'] - optimal_performance['risk']) /
                current_performance['risk']
            )

            revenue_increase = (
                (optimal_performance['revenue'] - current_performance['revenue']) /
                current_performance['revenue']
            )

            efficiency_gain = (
                (optimal_performance['capital_efficiency'] - current_performance['capital_efficiency']) /
                current_performance['capital_efficiency']
            )

            # Assess implementation risk
            implementation_risk = self._assess_implementation_risk(
                current_params,
                optimal_params
            )

            # Calculate confidence score
            confidence_score = self._calculate_optimization_confidence(
                backtest_results,
                result.fun
            )

            return OptimizationResult(
                current_params=current_params,
                optimal_params=optimal_params,
                expected_improvement={
                    'risk_reduction': risk_reduction,
                    'revenue_increase': revenue_increase,
                    'capital_efficiency_gain': efficiency_gain,
                    'objective_improvement': -result.fun  # Negative because we minimize
                },
                risk_reduction=risk_reduction,
                revenue_increase=revenue_increase,
                capital_efficiency_gain=efficiency_gain,
                implementation_risk=implementation_risk,
                confidence_score=confidence_score,
                backtest_results=backtest_results
            )

        except Exception as e:
            logger.error(f"Parameter optimization failed: {e}")
            raise

    def _optimization_objective(
        self,
        params: np.ndarray,
        market_conditions: Dict,
        historical_data: pd.DataFrame
    ) -> float:
        """
        Multi-objective optimization function
        Minimizes risk while maximizing revenue and capital efficiency
        """
        # Unpack parameters
        ltv, liq_threshold, liq_penalty, base_rate, slope1, slope2 = params

        # Create parameter object
        test_params = EconomicParameters(
            ltv=ltv,
            liquidation_threshold=liq_threshold,
            liquidation_penalty=liq_penalty,
            base_rate=base_rate,
            slope1=slope1,
            slope2=slope2
        )

        # Evaluate performance
        performance = self._evaluate_parameters(
            test_params,
            market_conditions,
            historical_data
        )

        # Multi-objective function (Gauntlet approach)
        # Minimize negative of good things, minimize positive of bad things
        objective = (
            0.4 * performance['risk'] -           # Minimize risk
            0.3 * performance['revenue'] -         # Maximize revenue
            0.2 * performance['capital_efficiency'] -  # Maximize efficiency
            0.1 * performance['stability']         # Minimize instability
        )

        return objective

    def _evaluate_parameters(
        self,
        params: EconomicParameters,
        market_conditions: Dict,
        historical_data: pd.DataFrame
    ) -> Dict[str, float]:
        """Evaluate parameter performance"""
        # Simulate with given parameters
        # This is simplified; real implementation would run full simulation

        # Risk metric (lower is better)
        risk = self._estimate_risk(params, market_conditions, historical_data)

        # Revenue metric (higher is better)
        revenue = self._estimate_revenue(params, market_conditions)

        # Capital efficiency (higher is better)
        efficiency = self._estimate_capital_efficiency(params, market_conditions)

        # Stability metric (lower is better - measures parameter sensitivity)
        stability = self._estimate_stability(params, historical_data)

        return {
            'risk': risk,
            'revenue': revenue,
            'capital_efficiency': efficiency,
            'stability': stability
        }

    def _estimate_risk(
        self,
        params: EconomicParameters,
        market_conditions: Dict,
        historical_data: pd.DataFrame
    ) -> float:
        """Estimate risk with given parameters"""
        # Simplified risk estimation
        # Real implementation would use full simulation

        # Lower LTV = lower risk
        ltv_risk = params.ltv

        # Lower liquidation threshold = higher risk
        threshold_risk = 1 - params.liquidation_threshold

        # Market volatility factor
        volatility = market_conditions.get('volatility', 0.03)

        # Combined risk score
        risk = (ltv_risk * 0.5 + threshold_risk * 0.3) * (1 + volatility)

        return risk

    def _estimate_revenue(
        self,
        params: EconomicParameters,
        market_conditions: Dict
    ) -> float:
        """Estimate protocol revenue"""
        # Revenue from interest rates
        utilization = market_conditions.get('utilization', 0.5)

        # Interest rate based on utilization
        if utilization < params.optimal_utilization:
            rate = params.base_rate + utilization * params.slope1
        else:
            rate = params.base_rate + params.slope1 * params.optimal_utilization + \
                   (utilization - params.optimal_utilization) * params.slope2

        # Revenue from liquidations
        liquidation_revenue = params.liquidation_penalty * market_conditions.get('liquidation_rate', 0.01)

        # Total revenue estimate
        total_borrowed = market_conditions.get('total_borrowed', 1e8)
        interest_revenue = rate * total_borrowed

        return interest_revenue + liquidation_revenue

    def _estimate_capital_efficiency(
        self,
        params: EconomicParameters,
        market_conditions: Dict
    ) -> float:
        """Estimate capital efficiency"""
        # Higher LTV = better capital efficiency
        # But must balance with risk

        # Utilization factor
        utilization = market_conditions.get('utilization', 0.5)

        # Efficiency score
        efficiency = params.ltv * utilization * (1 - params.reserve_factor)

        return efficiency

    def _estimate_stability(
        self,
        params: EconomicParameters,
        historical_data: pd.DataFrame
    ) -> float:
        """Estimate parameter stability/sensitivity"""
        # Measure how sensitive the system is to parameter changes
        # Lower is better (more stable)

        # Parameter sensitivity analysis
        sensitivity = 0

        # LTV sensitivity
        ltv_range = 0.85 - 0.50
        ltv_sensitivity = (params.ltv - 0.50) / ltv_range

        # Threshold sensitivity
        threshold_range = 0.90 - 0.60
        threshold_sensitivity = (params.liquidation_threshold - 0.60) / threshold_range

        # Rate curve sensitivity
        rate_sensitivity = params.slope2 / 3.0  # Normalized to max slope

        # Combined sensitivity (instability)
        sensitivity = (
            ltv_sensitivity * 0.3 +
            threshold_sensitivity * 0.3 +
            rate_sensitivity * 0.4
        )

        return sensitivity

    def _backtest_parameters(
        self,
        params: EconomicParameters,
        historical_data: pd.DataFrame,
        market_conditions: Dict
    ) -> Dict[str, Any]:
        """Backtest parameters on historical data"""
        results = {
            'periods_tested': len(historical_data),
            'liquidations': 0,
            'total_revenue': 0,
            'max_drawdown': 0,
            'avg_utilization': 0,
            'risk_events': []
        }

        # Simplified backtest
        # Real implementation would replay historical scenarios

        for i in range(len(historical_data)):
            # Simulate period
            period_data = historical_data.iloc[i]

            # Check for liquidations
            if period_data.get('price_change', 0) < -0.1:  # 10% drop
                liquidation_probability = 1 - params.liquidation_threshold
                if np.random.random() < liquidation_probability:
                    results['liquidations'] += 1
                    results['risk_events'].append({
                        'period': i,
                        'type': 'liquidation',
                        'severity': 'high'
                    })

            # Calculate revenue
            period_revenue = self._estimate_revenue(params, {'utilization': 0.5})
            results['total_revenue'] += period_revenue

        results['avg_utilization'] = 0.65  # Placeholder
        results['success_rate'] = 1 - (results['liquidations'] / max(len(historical_data), 1))

        return results

    def _assess_implementation_risk(
        self,
        current: EconomicParameters,
        optimal: EconomicParameters
    ) -> str:
        """Assess risk of implementing parameter changes"""
        # Calculate parameter deltas
        ltv_change = abs(optimal.ltv - current.ltv)
        threshold_change = abs(optimal.liquidation_threshold - current.liquidation_threshold)
        penalty_change = abs(optimal.liquidation_penalty - current.liquidation_penalty)
        rate_change = abs(optimal.base_rate - current.base_rate)

        # Assess magnitude of changes
        max_change = max(ltv_change, threshold_change, penalty_change, rate_change)

        if max_change > 0.2:
            return "HIGH - Major parameter changes require gradual rollout"
        elif max_change > 0.1:
            return "MEDIUM - Moderate changes, monitor closely during implementation"
        else:
            return "LOW - Minor adjustments, safe to implement"

    def _calculate_optimization_confidence(
        self,
        backtest_results: Dict,
        objective_value: float
    ) -> float:
        """Calculate confidence in optimization results"""
        # Factors affecting confidence:
        # 1. Backtest success rate
        # 2. Number of periods tested
        # 3. Objective function value

        confidence = 0.5  # Base confidence

        # Backtest performance
        if backtest_results['success_rate'] > 0.95:
            confidence += 0.2
        elif backtest_results['success_rate'] > 0.90:
            confidence += 0.1

        # Data quality
        if backtest_results['periods_tested'] > 1000:
            confidence += 0.2
        elif backtest_results['periods_tested'] > 500:
            confidence += 0.1

        # Optimization quality (lower objective is better)
        if objective_value < -0.5:
            confidence += 0.1

        return min(confidence, 1.0)

    # ============================================================
    # STRESS TESTING
    # ============================================================

    def run_stress_test(
        self,
        params: EconomicParameters,
        scenarios: List[Dict],
        portfolio: Dict
    ) -> Dict[str, Any]:
        """
        Run stress tests on given parameters

        Args:
            params: Parameters to test
            scenarios: Stress scenarios
            portfolio: Portfolio to test

        Returns:
            Stress test results
        """
        results = {}

        for scenario in scenarios:
            scenario_name = scenario['name']
            logger.info(f"Running stress test: {scenario_name}")

            # Apply stress
            stressed_market = self._apply_stress_scenario(scenario)

            # Calculate metrics under stress
            risk_metrics = self._calculate_stressed_metrics(
                params,
                portfolio,
                stressed_market
            )

            results[scenario_name] = {
                'risk_metrics': risk_metrics,
                'survival_probability': self._calculate_survival_probability(
                    risk_metrics,
                    scenario
                ),
                'expected_loss': risk_metrics.get('var_99', 0),
                'recovery_time': self._estimate_recovery_time(scenario)
            }

        return results

    def _apply_stress_scenario(self, scenario: Dict) -> Dict:
        """Apply stress scenario to market conditions"""
        stressed_market = {
            'prices': {},
            'volatility': {},
            'liquidity': {},
            'correlation': 0.8  # High correlation in stress
        }

        # Apply shocks
        for asset in ['ETH', 'BTC', 'USDC']:
            price_shock = scenario.get('price_shock', -0.30)
            volatility_mult = scenario.get('volatility_multiplier', 2.0)
            liquidity_shock = scenario.get('liquidity_shock', -0.50)

            stressed_market['prices'][asset] = 1 + price_shock
            stressed_market['volatility'][asset] = 0.03 * volatility_mult
            stressed_market['liquidity'][asset] = 1 + liquidity_shock

        return stressed_market

    def _calculate_stressed_metrics(
        self,
        params: EconomicParameters,
        portfolio: Dict,
        stressed_market: Dict
    ) -> Dict:
        """Calculate risk metrics under stress"""
        # Simplified stress metrics
        # Real implementation would be more comprehensive

        # Calculate stressed VaR
        portfolio_value = sum(p.get('value_usd', 0) for p in portfolio.values())
        stressed_var = portfolio_value * abs(stressed_market['prices']['ETH'] - 1)

        # Calculate probability of liquidation
        liquidation_prob = 0
        for position in portfolio.values():
            if position.get('leverage', 1) > 1:
                # Check if position would be liquidated
                price_drop = abs(stressed_market['prices'].get(position['asset'], 1) - 1)
                if price_drop > (1 / position['leverage']):
                    liquidation_prob = 1
                    break

        return {
            'var_99': stressed_var,
            'liquidation_probability': liquidation_prob,
            'expected_recovery_days': 30 * abs(stressed_market['prices']['ETH'] - 1)
        }

    def _calculate_survival_probability(
        self,
        risk_metrics: Dict,
        scenario: Dict
    ) -> float:
        """Calculate probability of surviving stress scenario"""
        # Based on liquidation probability and severity
        liquidation_prob = risk_metrics.get('liquidation_probability', 0)
        severity = scenario.get('severity', 'medium')

        severity_factors = {
            'low': 0.9,
            'medium': 0.7,
            'high': 0.5,
            'extreme': 0.3
        }

        base_survival = 1 - liquidation_prob
        severity_factor = severity_factors.get(severity, 0.7)

        return base_survival * severity_factor

    def _estimate_recovery_time(self, scenario: Dict) -> int:
        """Estimate recovery time from stress event (days)"""
        severity = scenario.get('severity', 'medium')

        recovery_times = {
            'low': 7,
            'medium': 30,
            'high': 90,
            'extreme': 180
        }

        return recovery_times.get(severity, 30)