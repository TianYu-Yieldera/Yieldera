#!/usr/bin/env python3
"""
A/B Testing Framework for AI Risk Engine Parameter Optimization
Implements statistical testing for parameter tuning
"""

import hashlib
import json
import time
import uuid
from dataclasses import dataclass, asdict
from datetime import datetime, timedelta
from enum import Enum
from typing import Any, Dict, List, Optional, Tuple, Union

import numpy as np
import pandas as pd
from scipy import stats
from scipy.stats import chi2_contingency, ttest_ind, mannwhitneyu
import statsmodels.stats.power as smp

# ============================================================
# EXPERIMENT TYPES
# ============================================================

class ExperimentType(Enum):
    """Types of A/B experiments"""
    PARAMETER_OPTIMIZATION = "parameter_optimization"
    MODEL_COMPARISON = "model_comparison"
    ALGORITHM_TESTING = "algorithm_testing"
    UI_FEATURE = "ui_feature"
    RISK_THRESHOLD = "risk_threshold"

class ExperimentStatus(Enum):
    """Experiment status"""
    DRAFT = "draft"
    RUNNING = "running"
    PAUSED = "paused"
    COMPLETED = "completed"
    FAILED = "failed"

# ============================================================
# EXPERIMENT CONFIGURATION
# ============================================================

@dataclass
class ExperimentConfig:
    """Configuration for an A/B test experiment"""

    experiment_id: str
    name: str
    description: str
    experiment_type: ExperimentType
    start_date: datetime
    end_date: datetime

    # Variants
    control_variant: Dict[str, Any]
    test_variants: List[Dict[str, Any]]

    # Statistical parameters
    significance_level: float = 0.05  # Alpha
    power: float = 0.80  # 1 - Beta
    minimum_detectable_effect: float = 0.05  # MDE

    # Traffic allocation
    traffic_percentage: float = 1.0  # Percentage of traffic to include
    variant_weights: Optional[Dict[str, float]] = None  # Custom weights

    # Metrics
    primary_metric: str
    secondary_metrics: List[str] = None
    guardrail_metrics: List[str] = None  # Metrics that shouldn't degrade

    # Sample size
    required_sample_size: Optional[int] = None
    max_duration_days: int = 30

@dataclass
class ExperimentResult:
    """Results from an A/B test experiment"""

    experiment_id: str
    status: ExperimentStatus
    start_time: datetime
    end_time: Optional[datetime]

    # Sample sizes
    control_samples: int
    test_samples: Dict[str, int]

    # Metric results
    control_metrics: Dict[str, float]
    test_metrics: Dict[str, Dict[str, float]]

    # Statistical results
    p_values: Dict[str, Dict[str, float]]
    confidence_intervals: Dict[str, Dict[str, Tuple[float, float]]]
    effect_sizes: Dict[str, Dict[str, float]]

    # Winner determination
    winner: Optional[str]
    confidence_level: float
    recommendation: str

# ============================================================
# A/B TESTING ENGINE
# ============================================================

class ABTestingEngine:
    """Core A/B testing engine with statistical analysis"""

    def __init__(self):
        self.experiments = {}
        self.active_experiments = set()
        self.experiment_data = {}

    def create_experiment(self, config: ExperimentConfig) -> str:
        """Create a new experiment"""
        experiment_id = config.experiment_id

        # Store experiment
        self.experiments[experiment_id] = {
            'config': config,
            'status': ExperimentStatus.DRAFT,
            'created_at': datetime.utcnow(),
            'data': {
                'control': [],
                'test_variants': {
                    f"variant_{i}": []
                    for i in range(len(config.test_variants))
                }
            }
        }

        # Calculate required sample size
        if not config.required_sample_size:
            config.required_sample_size = self._calculate_sample_size(config)

        return experiment_id

    def start_experiment(self, experiment_id: str):
        """Start an experiment"""
        if experiment_id not in self.experiments:
            raise ValueError(f"Experiment {experiment_id} not found")

        experiment = self.experiments[experiment_id]
        experiment['status'] = ExperimentStatus.RUNNING
        experiment['start_time'] = datetime.utcnow()

        self.active_experiments.add(experiment_id)

    def stop_experiment(self, experiment_id: str):
        """Stop an experiment"""
        if experiment_id in self.active_experiments:
            self.active_experiments.remove(experiment_id)

        if experiment_id in self.experiments:
            experiment = self.experiments[experiment_id]
            experiment['status'] = ExperimentStatus.COMPLETED
            experiment['end_time'] = datetime.utcnow()

    def get_variant(
        self,
        experiment_id: str,
        user_id: str
    ) -> Tuple[str, Dict[str, Any]]:
        """
        Determine which variant a user should see

        Returns:
            Tuple of (variant_name, variant_config)
        """
        if experiment_id not in self.active_experiments:
            return None, None

        config = self.experiments[experiment_id]['config']

        # Check traffic percentage
        if not self._should_include_user(user_id, config.traffic_percentage):
            return None, None

        # Deterministic assignment based on user ID
        variant = self._assign_variant(user_id, experiment_id, config)

        return variant

    def record_observation(
        self,
        experiment_id: str,
        user_id: str,
        variant: str,
        metrics: Dict[str, float]
    ):
        """Record an observation for the experiment"""
        if experiment_id not in self.experiments:
            return

        experiment = self.experiments[experiment_id]

        # Store observation
        observation = {
            'user_id': user_id,
            'variant': variant,
            'metrics': metrics,
            'timestamp': datetime.utcnow()
        }

        if variant == 'control':
            experiment['data']['control'].append(observation)
        else:
            if variant in experiment['data']['test_variants']:
                experiment['data']['test_variants'][variant].append(observation)

    def analyze_experiment(
        self,
        experiment_id: str
    ) -> ExperimentResult:
        """Analyze experiment results with statistical testing"""
        if experiment_id not in self.experiments:
            raise ValueError(f"Experiment {experiment_id} not found")

        experiment = self.experiments[experiment_id]
        config = experiment['config']
        data = experiment['data']

        # Extract metrics
        control_metrics = self._extract_metrics(data['control'])
        test_metrics = {
            variant: self._extract_metrics(observations)
            for variant, observations in data['test_variants'].items()
        }

        # Perform statistical tests
        p_values = {}
        confidence_intervals = {}
        effect_sizes = {}

        for variant, variant_data in data['test_variants'].items():
            if not variant_data:
                continue

            variant_results = {}
            ci_results = {}
            effect_results = {}

            # Test each metric
            for metric in [config.primary_metric] + (config.secondary_metrics or []):
                control_values = [
                    obs['metrics'].get(metric, 0)
                    for obs in data['control']
                ]
                test_values = [
                    obs['metrics'].get(metric, 0)
                    for obs in variant_data
                ]

                if control_values and test_values:
                    # Perform t-test
                    p_value = self._perform_statistical_test(
                        control_values,
                        test_values,
                        metric
                    )
                    variant_results[metric] = p_value

                    # Calculate confidence interval
                    ci = self._calculate_confidence_interval(
                        control_values,
                        test_values,
                        config.significance_level
                    )
                    ci_results[metric] = ci

                    # Calculate effect size
                    effect = self._calculate_effect_size(
                        control_values,
                        test_values
                    )
                    effect_results[metric] = effect

            p_values[variant] = variant_results
            confidence_intervals[variant] = ci_results
            effect_sizes[variant] = effect_results

        # Determine winner
        winner = self._determine_winner(
            p_values,
            effect_sizes,
            config
        )

        # Generate recommendation
        recommendation = self._generate_recommendation(
            winner,
            p_values,
            effect_sizes,
            config
        )

        return ExperimentResult(
            experiment_id=experiment_id,
            status=experiment['status'],
            start_time=experiment.get('start_time'),
            end_time=experiment.get('end_time'),
            control_samples=len(data['control']),
            test_samples={
                variant: len(observations)
                for variant, observations in data['test_variants'].items()
            },
            control_metrics=control_metrics,
            test_metrics=test_metrics,
            p_values=p_values,
            confidence_intervals=confidence_intervals,
            effect_sizes=effect_sizes,
            winner=winner,
            confidence_level=1 - config.significance_level,
            recommendation=recommendation
        )

    def _calculate_sample_size(self, config: ExperimentConfig) -> int:
        """Calculate required sample size for experiment"""
        # For simplicity, using formula for comparing two proportions
        # In practice, would adjust based on metric type

        effect_size = config.minimum_detectable_effect
        alpha = config.significance_level
        power = config.power

        # Cohen's d for effect size
        cohens_d = effect_size / 0.5  # Assuming pooled std of 0.5

        # Calculate sample size
        sample_size = smp.tt_solve_power(
            effect_size=cohens_d,
            alpha=alpha,
            power=power,
            ratio=1.0,
            alternative='two-sided'
        )

        # Multiply by number of variants
        total_variants = 1 + len(config.test_variants)
        total_sample_size = int(sample_size * total_variants)

        return total_sample_size

    def _should_include_user(self, user_id: str, traffic_percentage: float) -> bool:
        """Determine if user should be included in experiment"""
        # Hash user ID for consistent assignment
        hash_value = int(hashlib.md5(user_id.encode()).hexdigest(), 16)
        return (hash_value % 100) < (traffic_percentage * 100)

    def _assign_variant(
        self,
        user_id: str,
        experiment_id: str,
        config: ExperimentConfig
    ) -> Tuple[str, Dict[str, Any]]:
        """Assign user to a variant deterministically"""
        # Create unique hash for user + experiment
        combined = f"{user_id}:{experiment_id}"
        hash_value = int(hashlib.md5(combined.encode()).hexdigest(), 16)

        # Get weights
        if config.variant_weights:
            weights = config.variant_weights
        else:
            # Equal weights
            num_variants = 1 + len(config.test_variants)
            weights = {
                'control': 1.0 / num_variants,
                **{
                    f"variant_{i}": 1.0 / num_variants
                    for i in range(len(config.test_variants))
                }
            }

        # Assign based on hash
        cumulative = 0
        bucket = (hash_value % 1000) / 1000

        for variant, weight in weights.items():
            cumulative += weight
            if bucket < cumulative:
                if variant == 'control':
                    return 'control', config.control_variant
                else:
                    variant_idx = int(variant.split('_')[1])
                    return variant, config.test_variants[variant_idx]

        # Fallback to control
        return 'control', config.control_variant

    def _extract_metrics(self, observations: List[Dict]) -> Dict[str, float]:
        """Extract aggregated metrics from observations"""
        if not observations:
            return {}

        metrics = {}
        metric_names = set()

        # Collect all metric names
        for obs in observations:
            metric_names.update(obs['metrics'].keys())

        # Calculate mean for each metric
        for metric in metric_names:
            values = [
                obs['metrics'].get(metric, 0)
                for obs in observations
            ]
            metrics[metric] = np.mean(values) if values else 0

        return metrics

    def _perform_statistical_test(
        self,
        control_values: List[float],
        test_values: List[float],
        metric: str
    ) -> float:
        """Perform appropriate statistical test"""
        # Check if data is normally distributed
        _, p_control = stats.shapiro(control_values[:min(5000, len(control_values))])
        _, p_test = stats.shapiro(test_values[:min(5000, len(test_values))])

        if p_control > 0.05 and p_test > 0.05:
            # Use t-test for normally distributed data
            _, p_value = ttest_ind(control_values, test_values)
        else:
            # Use Mann-Whitney U test for non-normal data
            _, p_value = mannwhitneyu(control_values, test_values, alternative='two-sided')

        return p_value

    def _calculate_confidence_interval(
        self,
        control_values: List[float],
        test_values: List[float],
        alpha: float
    ) -> Tuple[float, float]:
        """Calculate confidence interval for difference in means"""
        control_mean = np.mean(control_values)
        test_mean = np.mean(test_values)
        diff = test_mean - control_mean

        # Pooled standard error
        control_var = np.var(control_values, ddof=1)
        test_var = np.var(test_values, ddof=1)

        pooled_se = np.sqrt(
            control_var / len(control_values) + test_var / len(test_values)
        )

        # Critical value
        z_critical = stats.norm.ppf(1 - alpha / 2)

        # Confidence interval
        ci_lower = diff - z_critical * pooled_se
        ci_upper = diff + z_critical * pooled_se

        return (ci_lower, ci_upper)

    def _calculate_effect_size(
        self,
        control_values: List[float],
        test_values: List[float]
    ) -> float:
        """Calculate Cohen's d effect size"""
        control_mean = np.mean(control_values)
        test_mean = np.mean(test_values)

        # Pooled standard deviation
        control_std = np.std(control_values, ddof=1)
        test_std = np.std(test_values, ddof=1)

        pooled_std = np.sqrt(
            ((len(control_values) - 1) * control_std**2 +
             (len(test_values) - 1) * test_std**2) /
            (len(control_values) + len(test_values) - 2)
        )

        if pooled_std == 0:
            return 0

        # Cohen's d
        effect_size = (test_mean - control_mean) / pooled_std

        return effect_size

    def _determine_winner(
        self,
        p_values: Dict,
        effect_sizes: Dict,
        config: ExperimentConfig
    ) -> Optional[str]:
        """Determine the winning variant"""
        winners = []

        for variant, variant_p_values in p_values.items():
            primary_p_value = variant_p_values.get(config.primary_metric)
            primary_effect = effect_sizes[variant].get(config.primary_metric, 0)

            if primary_p_value and primary_p_value < config.significance_level:
                # Check if effect is positive
                if primary_effect > config.minimum_detectable_effect:
                    # Check guardrail metrics
                    guardrails_ok = True

                    if config.guardrail_metrics:
                        for guardrail in config.guardrail_metrics:
                            guardrail_effect = effect_sizes[variant].get(guardrail, 0)
                            if guardrail_effect < -config.minimum_detectable_effect:
                                guardrails_ok = False
                                break

                    if guardrails_ok:
                        winners.append((variant, primary_effect))

        if winners:
            # Return variant with largest effect size
            winners.sort(key=lambda x: x[1], reverse=True)
            return winners[0][0]

        return None

    def _generate_recommendation(
        self,
        winner: Optional[str],
        p_values: Dict,
        effect_sizes: Dict,
        config: ExperimentConfig
    ) -> str:
        """Generate recommendation based on results"""
        if winner:
            primary_p = p_values[winner].get(config.primary_metric, 1)
            primary_effect = effect_sizes[winner].get(config.primary_metric, 0)

            return (
                f"Recommend implementing {winner}. "
                f"Primary metric shows {primary_effect:.2%} improvement "
                f"with p-value {primary_p:.4f}. "
                f"This result is statistically significant at {config.significance_level} level."
            )

        return (
            "No clear winner identified. "
            "Continue experiment to gather more data or "
            "consider larger effect sizes for practical significance."
        )

# ============================================================
# PARAMETER OPTIMIZATION A/B TESTER
# ============================================================

class ParameterOptimizationABTester:
    """Specialized A/B tester for parameter optimization"""

    def __init__(self, engine: ABTestingEngine):
        self.engine = engine

    def test_risk_parameters(
        self,
        current_params: Dict[str, float],
        test_params_list: List[Dict[str, float]],
        duration_days: int = 7
    ) -> str:
        """Test different risk parameter configurations"""

        config = ExperimentConfig(
            experiment_id=f"param_opt_{int(time.time())}",
            name="Risk Parameter Optimization",
            description="Testing optimized risk parameters",
            experiment_type=ExperimentType.PARAMETER_OPTIMIZATION,
            start_date=datetime.utcnow(),
            end_date=datetime.utcnow() + timedelta(days=duration_days),
            control_variant=current_params,
            test_variants=test_params_list,
            primary_metric="risk_adjusted_return",
            secondary_metrics=["liquidation_rate", "capital_efficiency"],
            guardrail_metrics=["max_drawdown", "systemic_risk"],
            minimum_detectable_effect=0.02,
            significance_level=0.05,
            power=0.80
        )

        experiment_id = self.engine.create_experiment(config)
        self.engine.start_experiment(experiment_id)

        return experiment_id

    def test_model_variants(
        self,
        models: List[Dict[str, Any]],
        duration_days: int = 14
    ) -> str:
        """Test different ML model variants"""

        control_model = models[0]
        test_models = models[1:]

        config = ExperimentConfig(
            experiment_id=f"model_test_{int(time.time())}",
            name="ML Model Comparison",
            description="Comparing ML model performance",
            experiment_type=ExperimentType.MODEL_COMPARISON,
            start_date=datetime.utcnow(),
            end_date=datetime.utcnow() + timedelta(days=duration_days),
            control_variant={"model": control_model},
            test_variants=[{"model": m} for m in test_models],
            primary_metric="prediction_accuracy",
            secondary_metrics=["precision", "recall", "f1_score"],
            guardrail_metrics=["prediction_latency"],
            minimum_detectable_effect=0.01,
            significance_level=0.01,  # More stringent for model changes
            power=0.90
        )

        experiment_id = self.engine.create_experiment(config)
        self.engine.start_experiment(experiment_id)

        return experiment_id

# ============================================================
# MULTI-ARMED BANDIT
# ============================================================

class MultiArmedBandit:
    """Multi-armed bandit for continuous optimization"""

    def __init__(self, arms: List[str], epsilon: float = 0.1):
        """
        Initialize multi-armed bandit

        Args:
            arms: List of arm names (variants)
            epsilon: Exploration rate for epsilon-greedy
        """
        self.arms = arms
        self.epsilon = epsilon
        self.counts = {arm: 0 for arm in arms}
        self.values = {arm: 0.0 for arm in arms}
        self.total_counts = 0

    def select_arm(self) -> str:
        """Select an arm using epsilon-greedy strategy"""
        if np.random.random() < self.epsilon:
            # Exploration: random selection
            return np.random.choice(self.arms)
        else:
            # Exploitation: select best performing
            return max(self.arms, key=lambda x: self.values[x])

    def update(self, arm: str, reward: float):
        """Update arm statistics"""
        self.counts[arm] += 1
        self.total_counts += 1

        # Update average reward (incremental update)
        n = self.counts[arm]
        value = self.values[arm]
        new_value = ((n - 1) * value + reward) / n
        self.values[arm] = new_value

    def get_statistics(self) -> Dict:
        """Get bandit statistics"""
        return {
            'counts': self.counts.copy(),
            'values': self.values.copy(),
            'total_counts': self.total_counts,
            'best_arm': max(self.arms, key=lambda x: self.values[x])
        }

# ============================================================
# EXAMPLE USAGE
# ============================================================

def example_usage():
    """Example of using the A/B testing framework"""

    # Initialize engine
    engine = ABTestingEngine()

    # Create an experiment
    config = ExperimentConfig(
        experiment_id="test_001",
        name="LTV Optimization Test",
        description="Testing new LTV parameters",
        experiment_type=ExperimentType.PARAMETER_OPTIMIZATION,
        start_date=datetime.utcnow(),
        end_date=datetime.utcnow() + timedelta(days=7),
        control_variant={"ltv": 0.70, "liquidation_threshold": 0.75},
        test_variants=[
            {"ltv": 0.72, "liquidation_threshold": 0.77},
            {"ltv": 0.68, "liquidation_threshold": 0.73}
        ],
        primary_metric="revenue_per_user",
        secondary_metrics=["liquidation_rate", "user_retention"],
        guardrail_metrics=["systemic_risk"],
        significance_level=0.05,
        power=0.80
    )

    experiment_id = engine.create_experiment(config)
    engine.start_experiment(experiment_id)

    # Simulate some observations
    np.random.seed(42)

    for i in range(1000):
        user_id = f"user_{i}"

        # Get variant for user
        variant, params = engine.get_variant(experiment_id, user_id)

        if variant:
            # Simulate metrics based on variant
            if variant == "control":
                revenue = np.random.normal(100, 20)
                liquidation = np.random.binomial(1, 0.05)
            else:
                # Test variants might perform differently
                revenue = np.random.normal(105, 20)
                liquidation = np.random.binomial(1, 0.04)

            # Record observation
            engine.record_observation(
                experiment_id,
                user_id,
                variant,
                {
                    "revenue_per_user": revenue,
                    "liquidation_rate": liquidation,
                    "user_retention": np.random.random(),
                    "systemic_risk": np.random.random() * 0.1
                }
            )

    # Analyze results
    results = engine.analyze_experiment(experiment_id)

    print(f"Experiment: {results.experiment_id}")
    print(f"Status: {results.status}")
    print(f"Control samples: {results.control_samples}")
    print(f"Test samples: {results.test_samples}")
    print(f"\nWinner: {results.winner}")
    print(f"Confidence: {results.confidence_level:.1%}")
    print(f"\nRecommendation: {results.recommendation}")

    # Multi-armed bandit example
    print("\n" + "="*50)
    print("Multi-Armed Bandit Example")
    print("="*50)

    bandit = MultiArmedBandit(["param_set_a", "param_set_b", "param_set_c"])

    for _ in range(1000):
        arm = bandit.select_arm()

        # Simulate reward
        if arm == "param_set_a":
            reward = np.random.normal(10, 2)
        elif arm == "param_set_b":
            reward = np.random.normal(12, 3)  # Best arm
        else:
            reward = np.random.normal(8, 2)

        bandit.update(arm, reward)

    stats = bandit.get_statistics()
    print(f"Best arm: {stats['best_arm']}")
    print(f"Arm values: {stats['values']}")
    print(f"Arm counts: {stats['counts']}")

if __name__ == "__main__":
    example_usage()