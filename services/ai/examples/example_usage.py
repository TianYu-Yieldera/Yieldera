#!/usr/bin/env python3
"""
AI Risk Engine - Example Usage Scripts
Demonstrates all Phase 2 capabilities
"""

import asyncio
import json
import numpy as np
import pandas as pd
from datetime import datetime, timedelta
import sys
import os

# Add parent directory to path
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from core.data_infrastructure import DataInfrastructure
from core.risk_calculator import RiskCalculator
from core.agent_simulation import GauntletSimulator, AgentType
from core.economic_model import GauntletEconomicModel, EconomicParameters
from core.ml_predictor import HybridMLPredictor
from core.realtime_monitor import GauntletRealTimeMonitor
from core.report_generator import GauntletReportGenerator

# ============================================================
# EXAMPLE 1: AGENT-BASED SIMULATION
# ============================================================

async def example_simulation():
    """Run market simulation with 10,000 agents"""
    print("\n" + "="*60)
    print("EXAMPLE 1: Agent-Based Simulation")
    print("="*60)

    # Initialize simulator
    simulator = GauntletSimulator(num_agents=10000, seed=42)

    # Define scenarios
    scenarios = [
        {
            'name': 'Market Crash',
            'market_shock': -0.30,
            'volatility_multiplier': 2.0
        },
        {
            'name': 'Liquidity Crisis',
            'liquidity_shock': -0.50,
            'utilization_spike': 0.95
        },
        {
            'name': 'Black Swan',
            'market_shock': -0.60,
            'volatility_multiplier': 5.0,
            'correlation': 0.9
        }
    ]

    # Run simulation
    print(f"\n📊 Running simulation with {simulator.num_agents} agents...")
    results = await simulator.run_comprehensive_simulation(
        scenarios=scenarios,
        time_steps=100  # Reduced for demo
    )

    # Display results
    print("\n📈 Simulation Results:")
    for scenario_name, result in results['results'].items():
        print(f"\n  Scenario: {scenario_name}")
        print(f"  • Systemic Risk: {result.risk_metrics['systemic_risk']:.2f}")
        print(f"  • Total Liquidations: {result.risk_metrics['total_liquidations']}")
        print(f"  • Avg Volatility: {result.risk_metrics['avg_volatility']:.2%}")
        print(f"  • Execution Time: {result.execution_time:.2f}s")

    print("\n💡 Recommendations:")
    for i, rec in enumerate(results['recommendations'][:3], 1):
        print(f"  {i}. {rec}")

    return results

# ============================================================
# EXAMPLE 2: ML PREDICTION
# ============================================================

async def example_ml_prediction():
    """Train and use ML models for liquidation prediction"""
    print("\n" + "="*60)
    print("EXAMPLE 2: Machine Learning Prediction")
    print("="*60)

    # Initialize predictor
    predictor = HybridMLPredictor()

    # Generate sample training data
    print("\n🔧 Generating training data...")
    n_samples = 5000
    n_features = 20

    X_structured = np.random.randn(n_samples, n_features)
    X_sequences = np.random.randn(n_samples, 60, 10)  # 60 timesteps, 10 features

    # Create labels (binary: liquidated or not)
    # Make it somewhat correlated with first feature for realism
    y = (X_structured[:, 0] + np.random.randn(n_samples) * 0.5 > 1).astype(int)

    training_data = {
        'X_structured': X_structured,
        'X_sequences': X_sequences,
        'y': y
    }

    # Train models
    print(f"🎓 Training models on {n_samples} samples...")
    performance = predictor.train(training_data)

    print(f"\n📊 Model Performance:")
    print(f"  • Accuracy: {performance.accuracy:.2%}")
    print(f"  • Precision: {performance.precision:.2%}")
    print(f"  • Recall: {performance.recall:.2%}")
    print(f"  • F1 Score: {performance.f1_score:.2%}")
    print(f"  • AUC-ROC: {performance.auc_roc:.2%}")

    # Make prediction
    print("\n🔮 Making prediction for sample position...")

    position_data = {
        'health_factor': 1.3,
        'ltv': 0.75,
        'leverage': 3.5,
        'collateral_value_usd': 50000,
        'debt_value_usd': 37500,
        'protocol': 'aave'
    }

    market_data = {
        'volatility': 0.05,
        'utilization_rate': 0.8,
        'liquidity_ratio': 0.5
    }

    prediction = predictor.predict(position_data, market_data)

    print(f"\n⚠️  Prediction Results:")
    print(f"  • Liquidation Probability: {prediction.liquidation_probability:.1%}")
    print(f"  • Confidence Score: {prediction.confidence_score:.1%}")
    print(f"  • Time to Liquidation: {prediction.time_to_liquidation or 'N/A'} hours")

    print(f"\n🔍 Risk Factors:")
    for factor, score in prediction.risk_factors.items():
        print(f"  • {factor}: {score:.2f}")

    print(f"\n💡 Recommendations:")
    for i, rec in enumerate(prediction.recommendations[:3], 1):
        print(f"  {i}. {rec}")

    return predictor, prediction

# ============================================================
# EXAMPLE 3: PARAMETER OPTIMIZATION
# ============================================================

async def example_parameter_optimization():
    """Optimize protocol parameters"""
    print("\n" + "="*60)
    print("EXAMPLE 3: Economic Parameter Optimization")
    print("="*60)

    # Initialize economic model
    model = GauntletEconomicModel()

    # Current parameters
    current_params = EconomicParameters(
        ltv=0.70,
        liquidation_threshold=0.75,
        liquidation_penalty=0.10,
        base_rate=0.02,
        slope1=0.08,
        slope2=1.0
    )

    # Market conditions
    market_conditions = {
        'utilization': 0.65,
        'volatility': 0.04,
        'liquidation_rate': 0.02,
        'total_borrowed': 100_000_000
    }

    # Historical data (simplified)
    historical_data = pd.DataFrame({
        'ETH': np.random.randn(100) * 0.05 + 1,
        'BTC': np.random.randn(100) * 0.04 + 1,
        'price_change': np.random.randn(100) * 0.03
    })

    print("\n📊 Current Parameters:")
    print(f"  • LTV: {current_params.ltv:.2%}")
    print(f"  • Liquidation Threshold: {current_params.liquidation_threshold:.2%}")
    print(f"  • Liquidation Penalty: {current_params.liquidation_penalty:.2%}")
    print(f"  • Base Rate: {current_params.base_rate:.2%}")

    print("\n🔧 Optimizing parameters...")
    result = model.optimize_protocol_parameters(
        current_params=current_params,
        market_conditions=market_conditions,
        historical_data=historical_data
    )

    print("\n✅ Optimal Parameters:")
    print(f"  • LTV: {result.optimal_params.ltv:.2%}")
    print(f"  • Liquidation Threshold: {result.optimal_params.liquidation_threshold:.2%}")
    print(f"  • Liquidation Penalty: {result.optimal_params.liquidation_penalty:.2%}")
    print(f"  • Base Rate: {result.optimal_params.base_rate:.2%}")

    print("\n📈 Expected Improvements:")
    print(f"  • Risk Reduction: {result.risk_reduction:.1%}")
    print(f"  • Revenue Increase: {result.revenue_increase:.1%}")
    print(f"  • Capital Efficiency Gain: {result.capital_efficiency_gain:.1%}")

    print(f"\n⚠️  Implementation Risk: {result.implementation_risk}")
    print(f"🎯 Confidence Score: {result.confidence_score:.1%}")

    return result

# ============================================================
# EXAMPLE 4: REAL-TIME MONITORING
# ============================================================

async def example_monitoring():
    """Demonstrate real-time monitoring capabilities"""
    print("\n" + "="*60)
    print("EXAMPLE 4: Real-Time Monitoring")
    print("="*60)

    # Note: This requires database connection
    # Using mock data for demonstration

    print("\n🔍 Monitoring System Features:")
    print("  • 24/7 position surveillance")
    print("  • Multi-level alert system")
    print("  • WebSocket real-time updates")
    print("  • Automated response capabilities")

    # Simulate monitoring metrics
    mock_metrics = {
        'positions_monitored': 247,
        'alerts_triggered': 12,
        'alerts_by_level': {
            'critical': 1,
            'danger': 3,
            'warning': 5,
            'info': 3
        },
        'at_risk_positions': 8,
        'total_value_at_risk': 1_500_000,
        'avg_health_factor': 1.85,
        'system_health': 'healthy',
        'response_time_ms': 85
    }

    print("\n📊 Current Monitoring Metrics:")
    print(f"  • Positions Monitored: {mock_metrics['positions_monitored']}")
    print(f"  • Alerts Triggered: {mock_metrics['alerts_triggered']}")
    print(f"  • At-Risk Positions: {mock_metrics['at_risk_positions']}")
    print(f"  • Value at Risk: ${mock_metrics['total_value_at_risk']:,.0f}")
    print(f"  • System Health: {mock_metrics['system_health']}")
    print(f"  • Response Time: {mock_metrics['response_time_ms']}ms")

    print("\n🚨 Alert Distribution:")
    for level, count in mock_metrics['alerts_by_level'].items():
        print(f"  • {level.upper()}: {count}")

    print("\n💡 To see real-time monitoring:")
    print("  1. Start the API: python api/main_v2.py")
    print("  2. Open WebSocket client: websocket_client.html")
    print("  3. Watch live alerts and metrics flow!")

    return mock_metrics

# ============================================================
# EXAMPLE 5: REPORT GENERATION
# ============================================================

async def example_report_generation():
    """Generate professional risk report"""
    print("\n" + "="*60)
    print("EXAMPLE 5: Report Generation")
    print("="*60)

    # Initialize report generator
    generator = GauntletReportGenerator()

    # Sample analysis results
    analysis_results = {
        'overall_risk_score': 68.5,
        'alert_level': 'warning',
        'timestamp': datetime.utcnow(),
        'key_risks': [
            {'factor': 'High market volatility', 'severity': 'high'},
            {'factor': 'Elevated leverage', 'severity': 'medium'},
            {'factor': 'Liquidity constraints', 'severity': 'low'}
        ],
        'top_recommendations': [
            'Reduce leverage in high-risk positions',
            'Increase monitoring frequency to 30-second intervals',
            'Consider hedging with options or futures'
        ],
        'risk_metrics': {
            'var_95': 150000,
            'var_99': 300000,
            'max_drawdown': 0.18,
            'sharpe_ratio': 1.45,
            'liquidation_probability': 0.35
        },
        'market_data': {
            'avg_volatility': 0.045,
            'total_liquidity': 2_500_000_000,
            'utilization_rate': 0.72,
            'liquidations_24h': 89
        },
        'simulation_results': {
            'scenarios': ['Market Crash', 'Liquidity Crisis'],
            'worst_case_scenario': 'Market Crash',
            'expected_losses_99': 450000,
            'systemic_risk': 0.65,
            'cascade_probability': 0.23
        }
    }

    print("\n📄 Generating comprehensive risk report...")
    report = generator.generate_comprehensive_report(
        analysis_results,
        report_type='full'
    )

    print(f"\n✅ Report Generated:")
    print(f"  • Report ID: {report.report_id}")
    print(f"  • Title: {report.title}")
    print(f"  • Generated At: {report.generated_at}")
    print(f"  • Sections: {len(report.sections)}")
    print(f"  • Visualizations: {len(report.visualizations)}")

    print(f"\n📑 Report Sections:")
    for section in report.sections[:5]:
        print(f"  • {section.title}")

    # Export to different formats
    print(f"\n💾 Export Options:")
    print("  • JSON: Complete data export")
    print("  • HTML: Interactive web report")
    print("  • PDF: Professional document")

    # Save JSON report
    json_report = generator.export_to_json(report)
    output_file = "example_report.json"

    with open(output_file, 'w') as f:
        f.write(json_report)

    print(f"\n✅ Report saved to: {output_file}")

    return report

# ============================================================
# EXAMPLE 6: STRESS TESTING
# ============================================================

async def example_stress_testing():
    """Run comprehensive stress tests"""
    print("\n" + "="*60)
    print("EXAMPLE 6: Stress Testing")
    print("="*60)

    model = GauntletEconomicModel()

    # Define stress scenarios
    stress_scenarios = [
        {
            'name': 'Moderate Stress',
            'price_shock': -0.20,
            'volatility_multiplier': 1.5,
            'liquidity_shock': -0.30,
            'severity': 'medium'
        },
        {
            'name': 'Severe Stress',
            'price_shock': -0.40,
            'volatility_multiplier': 3.0,
            'liquidity_shock': -0.60,
            'severity': 'high'
        },
        {
            'name': 'Extreme Stress',
            'price_shock': -0.60,
            'volatility_multiplier': 5.0,
            'liquidity_shock': -0.80,
            'severity': 'extreme'
        }
    ]

    # Current parameters
    params = EconomicParameters(
        ltv=0.75,
        liquidation_threshold=0.80,
        liquidation_penalty=0.05,
        base_rate=0.02,
        slope1=0.08,
        slope2=1.0
    )

    # Sample portfolio
    portfolio = {
        'position_1': {'value_usd': 100000, 'asset': 'ETH', 'leverage': 2.0},
        'position_2': {'value_usd': 50000, 'asset': 'BTC', 'leverage': 1.5},
        'position_3': {'value_usd': 75000, 'asset': 'ETH', 'leverage': 3.0}
    }

    print(f"\n📊 Portfolio Under Test:")
    total_value = sum(p['value_usd'] for p in portfolio.values())
    print(f"  • Total Value: ${total_value:,.0f}")
    print(f"  • Positions: {len(portfolio)}")
    print(f"  • Avg Leverage: {np.mean([p['leverage'] for p in portfolio.values()]):.1f}x")

    print(f"\n🔥 Running stress tests...")
    results = model.run_stress_test(params, stress_scenarios, portfolio)

    print(f"\n📈 Stress Test Results:")
    for scenario_name, outcome in results.items():
        print(f"\n  {scenario_name}:")
        print(f"    • Survival Probability: {outcome['survival_probability']:.1%}")
        print(f"    • Expected Loss: ${outcome['expected_loss']:,.0f}")
        print(f"    • Recovery Time: {outcome['recovery_time']} days")

    # Calculate overall resilience
    avg_survival = np.mean([r['survival_probability'] for r in results.values()])
    max_loss = max(r['expected_loss'] for r in results.values())

    print(f"\n🎯 Overall Assessment:")
    print(f"  • Average Survival Rate: {avg_survival:.1%}")
    print(f"  • Maximum Expected Loss: ${max_loss:,.0f}")

    if avg_survival > 0.7:
        print("  • ✅ Portfolio shows good resilience")
    elif avg_survival > 0.5:
        print("  • ⚠️ Portfolio has moderate risk")
    else:
        print("  • ❌ Portfolio is high risk")

    return results

# ============================================================
# MAIN EXECUTION
# ============================================================

async def main():
    """Run all examples"""
    print("\n" + "🚀 "*20)
    print(" AI RISK ENGINE - COMPREHENSIVE EXAMPLES")
    print("🚀 "*20)

    try:
        # Run examples
        await example_simulation()
        await example_ml_prediction()
        await example_parameter_optimization()
        await example_monitoring()
        await example_report_generation()
        await example_stress_testing()

        print("\n" + "="*60)
        print("✅ ALL EXAMPLES COMPLETED SUCCESSFULLY!")
        print("="*60)

        print("\n📚 Next Steps:")
        print("  1. Start the API server: python api/main_v2.py")
        print("  2. Open WebSocket client: websocket_client.html")
        print("  3. Use the API endpoints for production integration")
        print("  4. Train ML models with your actual data")

    except Exception as e:
        print(f"\n❌ Error running examples: {e}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    # Run async main
    asyncio.run(main())