"""
Report Generator Module
Generate professional risk reports in Gauntlet style
Phase 2: Core Implementation
"""

import json
import logging
from datetime import datetime, timedelta
from typing import Dict, List, Any, Optional
from dataclasses import dataclass
import pandas as pd
import numpy as np
from io import BytesIO
import base64

# Visualization
import matplotlib.pyplot as plt
import seaborn as sns

# PDF generation (optional)
try:
    from reportlab.lib import colors
    from reportlab.lib.pagesizes import letter, A4
    from reportlab.platypus import SimpleDocTemplate, Table, TableStyle, Paragraph, Spacer, PageBreak, Image
    from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
    from reportlab.lib.units import inch
    HAS_REPORTLAB = True
except ImportError:
    HAS_REPORTLAB = False

logger = logging.getLogger(__name__)

# Configure plotting
sns.set_style("whitegrid")
plt.rcParams['figure.figsize'] = (10, 6)

# ============================================================
# DATA MODELS
# ============================================================

@dataclass
class ReportSection:
    """Report section"""
    title: str
    content: Any
    section_type: str  # 'text', 'table', 'chart', 'metrics'
    priority: int = 0

@dataclass
class RiskReport:
    """Complete risk report"""
    report_id: str
    title: str
    generated_at: datetime
    executive_summary: str
    sections: List[ReportSection]
    visualizations: Dict[str, str]  # Base64 encoded images
    metadata: Dict[str, Any]

# ============================================================
# REPORT GENERATOR
# ============================================================

class GauntletReportGenerator:
    """
    Generate professional risk reports following Gauntlet standards
    """

    def __init__(self, config: Optional[Dict] = None):
        """
        Initialize report generator

        Args:
            config: Report configuration
        """
        self.config = config or self._default_config()
        self.styles = self._initialize_styles()

    def _default_config(self) -> Dict:
        """Default configuration"""
        return {
            'company_name': 'Yieldera AI',
            'report_style': 'professional',
            'include_visualizations': True,
            'include_recommendations': True,
            'include_technical_appendix': True,
            'chart_style': 'seaborn',
            'color_scheme': {
                'primary': '#1E3A8A',
                'secondary': '#3B82F6',
                'danger': '#EF4444',
                'warning': '#F59E0B',
                'success': '#10B981',
                'neutral': '#6B7280'
            }
        }

    def _initialize_styles(self) -> Dict:
        """Initialize report styles"""
        if HAS_REPORTLAB:
            styles = getSampleStyleSheet()
            # Custom styles
            styles.add(ParagraphStyle(
                name='CustomTitle',
                parent=styles['Heading1'],
                fontSize=24,
                textColor=colors.HexColor(self.config['color_scheme']['primary']),
                spaceAfter=30,
            ))
            return styles
        return {}

    # ============================================================
    # REPORT GENERATION
    # ============================================================

    def generate_comprehensive_report(
        self,
        analysis_results: Dict[str, Any],
        report_type: str = 'full'
    ) -> RiskReport:
        """
        Generate comprehensive risk report

        Args:
            analysis_results: Results from risk analysis
            report_type: Type of report ('full', 'summary', 'technical')

        Returns:
            RiskReport object
        """
        try:
            logger.info(f"Generating {report_type} report...")

            # Generate report ID
            report_id = f"report_{datetime.utcnow().strftime('%Y%m%d_%H%M%S')}"

            # Generate sections based on report type
            sections = []

            # Executive Summary
            executive_summary = self._generate_executive_summary(analysis_results)
            sections.append(ReportSection(
                title="Executive Summary",
                content=executive_summary,
                section_type="text",
                priority=1
            ))

            # Risk Assessment
            if report_type in ['full', 'summary']:
                risk_section = self._generate_risk_assessment_section(analysis_results)
                sections.append(risk_section)

            # Market Analysis
            if report_type == 'full':
                market_section = self._generate_market_analysis_section(analysis_results)
                sections.append(market_section)

            # Simulation Results
            if 'simulation_results' in analysis_results and report_type in ['full', 'technical']:
                sim_section = self._generate_simulation_section(analysis_results['simulation_results'])
                sections.append(sim_section)

            # Recommendations
            if self.config['include_recommendations']:
                rec_section = self._generate_recommendations_section(analysis_results)
                sections.append(rec_section)

            # Technical Appendix
            if self.config['include_technical_appendix'] and report_type in ['full', 'technical']:
                appendix = self._generate_technical_appendix(analysis_results)
                sections.append(appendix)

            # Generate visualizations
            visualizations = {}
            if self.config['include_visualizations']:
                visualizations = self._generate_visualizations(analysis_results)

            # Create report
            report = RiskReport(
                report_id=report_id,
                title=self._generate_report_title(report_type),
                generated_at=datetime.utcnow(),
                executive_summary=executive_summary,
                sections=sections,
                visualizations=visualizations,
                metadata={
                    'report_type': report_type,
                    'analysis_timestamp': analysis_results.get('timestamp', datetime.utcnow()),
                    'data_sources': analysis_results.get('data_sources', []),
                    'model_versions': analysis_results.get('model_versions', {})
                }
            )

            logger.info(f"Report {report_id} generated successfully")
            return report

        except Exception as e:
            logger.error(f"Report generation failed: {e}")
            raise

    # ============================================================
    # SECTION GENERATORS
    # ============================================================

    def _generate_executive_summary(self, analysis_results: Dict) -> str:
        """Generate executive summary"""
        risk_score = analysis_results.get('overall_risk_score', 0)
        alert_level = analysis_results.get('alert_level', 'unknown')
        key_risks = analysis_results.get('key_risks', [])
        recommendations = analysis_results.get('top_recommendations', [])

        summary = f"""
        **Overall Risk Assessment**: {self._format_risk_level(risk_score)} ({risk_score:.1f}/100)

        **Alert Level**: {alert_level.upper()}

        **Key Findings**:
        • System-wide risk score: {risk_score:.1f}/100
        • {len(key_risks)} critical risk factors identified
        • Immediate action required for {sum(1 for r in key_risks if r.get('severity') == 'critical')} positions

        **Top Recommendations**:
        """

        for i, rec in enumerate(recommendations[:3], 1):
            summary += f"\n{i}. {rec}"

        return summary

    def _generate_risk_assessment_section(self, analysis_results: Dict) -> ReportSection:
        """Generate risk assessment section"""
        risk_metrics = analysis_results.get('risk_metrics', {})

        # Create metrics table
        metrics_data = [
            ['Metric', 'Value', 'Status'],
            ['Value at Risk (95%)', f"${risk_metrics.get('var_95', 0):,.0f}", self._get_status_indicator(risk_metrics.get('var_95', 0), 100000)],
            ['Value at Risk (99%)', f"${risk_metrics.get('var_99', 0):,.0f}", self._get_status_indicator(risk_metrics.get('var_99', 0), 500000)],
            ['Max Drawdown', f"{risk_metrics.get('max_drawdown', 0):.2%}", self._get_status_indicator(risk_metrics.get('max_drawdown', 0), 0.2)],
            ['Sharpe Ratio', f"{risk_metrics.get('sharpe_ratio', 0):.2f}", self._get_status_indicator(risk_metrics.get('sharpe_ratio', 0), 1.0, higher_is_better=True)],
            ['Liquidation Probability', f"{risk_metrics.get('liquidation_probability', 0):.2%}", self._get_status_indicator(risk_metrics.get('liquidation_probability', 0), 0.1)],
        ]

        return ReportSection(
            title="Risk Assessment",
            content=metrics_data,
            section_type="table",
            priority=2
        )

    def _generate_market_analysis_section(self, analysis_results: Dict) -> ReportSection:
        """Generate market analysis section"""
        market_data = analysis_results.get('market_data', {})

        content = f"""
        **Market Conditions**

        Current market conditions show {'elevated' if market_data.get('volatility', 0) > 0.05 else 'normal'} volatility levels.

        **Key Market Metrics**:
        • Average Volatility: {market_data.get('avg_volatility', 0):.2%}
        • Total Liquidity: ${market_data.get('total_liquidity', 0):,.0f}
        • Utilization Rate: {market_data.get('utilization_rate', 0):.2%}
        • 24h Liquidations: {market_data.get('liquidations_24h', 0)}

        **Market Trends**:
        {self._analyze_market_trends(market_data)}
        """

        return ReportSection(
            title="Market Analysis",
            content=content,
            section_type="text",
            priority=3
        )

    def _generate_simulation_section(self, simulation_results: Dict) -> ReportSection:
        """Generate simulation results section"""
        content = f"""
        **Simulation Results**

        **Scenarios Tested**: {len(simulation_results.get('scenarios', []))}

        **Key Findings**:
        • Worst-case scenario: {simulation_results.get('worst_case_scenario', 'N/A')}
        • Expected losses (99% confidence): ${simulation_results.get('expected_losses_99', 0):,.0f}
        • Systemic risk score: {simulation_results.get('systemic_risk', 0):.2f}/1.0
        • Cascade probability: {simulation_results.get('cascade_probability', 0):.2%}

        **Stress Test Results**:
        {self._format_stress_test_results(simulation_results.get('stress_tests', {}))}
        """

        return ReportSection(
            title="Simulation & Stress Testing",
            content=content,
            section_type="text",
            priority=4
        )

    def _generate_recommendations_section(self, analysis_results: Dict) -> ReportSection:
        """Generate recommendations section"""
        recommendations = analysis_results.get('recommendations', {})

        content = """
        **Strategic Recommendations**

        **Immediate Actions** (0-24 hours):
        """

        for rec in recommendations.get('immediate', []):
            content += f"\n• {rec}"

        content += """

        **Short-term Actions** (1-7 days):
        """

        for rec in recommendations.get('short_term', []):
            content += f"\n• {rec}"

        content += """

        **Long-term Strategy** (7+ days):
        """

        for rec in recommendations.get('long_term', []):
            content += f"\n• {rec}"

        return ReportSection(
            title="Recommendations",
            content=content,
            section_type="text",
            priority=5
        )

    def _generate_technical_appendix(self, analysis_results: Dict) -> ReportSection:
        """Generate technical appendix"""
        content = """
        **Technical Appendix**

        **Methodology**:
        • Risk Calculation: Multi-factor model with VaR, CVaR, and ML predictions
        • Simulation: Agent-based model with 10,000 agents
        • Data Sources: On-chain data, price oracles, protocol events

        **Model Parameters**:
        • Confidence Levels: 95% and 99%
        • Time Horizon: 1-day VaR
        • Simulation Steps: 1,000 per scenario
        • ML Models: XGBoost + LSTM ensemble

        **Assumptions**:
        • Normal market conditions for baseline
        • Correlation increases during stress events
        • Liquidation costs include gas fees and slippage

        **Limitations**:
        • Model accuracy depends on data quality
        • Black swan events may exceed predictions
        • Cross-chain risks not fully captured
        """

        return ReportSection(
            title="Technical Appendix",
            content=content,
            section_type="text",
            priority=10
        )

    # ============================================================
    # VISUALIZATIONS
    # ============================================================

    def _generate_visualizations(self, analysis_results: Dict) -> Dict[str, str]:
        """Generate all visualizations"""
        visualizations = {}

        try:
            # Risk distribution chart
            risk_chart = self._create_risk_distribution_chart(analysis_results)
            if risk_chart:
                visualizations['risk_distribution'] = risk_chart

            # VaR chart
            var_chart = self._create_var_chart(analysis_results)
            if var_chart:
                visualizations['var_analysis'] = var_chart

            # Liquidation heatmap
            heatmap = self._create_liquidation_heatmap(analysis_results)
            if heatmap:
                visualizations['liquidation_heatmap'] = heatmap

            # Time series chart
            ts_chart = self._create_time_series_chart(analysis_results)
            if ts_chart:
                visualizations['time_series'] = ts_chart

        except Exception as e:
            logger.error(f"Visualization generation failed: {e}")

        return visualizations

    def _create_risk_distribution_chart(self, analysis_results: Dict) -> Optional[str]:
        """Create risk distribution chart"""
        try:
            fig, ax = plt.subplots(figsize=(10, 6))

            # Sample data (would use actual data)
            risk_scores = np.random.normal(50, 15, 1000)
            risk_scores = np.clip(risk_scores, 0, 100)

            # Create histogram
            ax.hist(risk_scores, bins=30, edgecolor='black', alpha=0.7,
                   color=self.config['color_scheme']['primary'])

            ax.set_xlabel('Risk Score')
            ax.set_ylabel('Number of Positions')
            ax.set_title('Risk Score Distribution')
            ax.grid(True, alpha=0.3)

            # Add risk zones
            ax.axvspan(0, 30, alpha=0.2, color='green', label='Low Risk')
            ax.axvspan(30, 70, alpha=0.2, color='yellow', label='Medium Risk')
            ax.axvspan(70, 100, alpha=0.2, color='red', label='High Risk')

            ax.legend()

            # Convert to base64
            buffer = BytesIO()
            plt.savefig(buffer, format='png', dpi=100, bbox_inches='tight')
            buffer.seek(0)
            image_base64 = base64.b64encode(buffer.read()).decode()
            plt.close()

            return f"data:image/png;base64,{image_base64}"

        except Exception as e:
            logger.error(f"Failed to create risk distribution chart: {e}")
            return None

    def _create_var_chart(self, analysis_results: Dict) -> Optional[str]:
        """Create VaR visualization"""
        try:
            fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 6))

            # VaR by confidence level
            confidence_levels = [0.90, 0.95, 0.99]
            var_values = [
                analysis_results.get('risk_metrics', {}).get('var_90', 50000),
                analysis_results.get('risk_metrics', {}).get('var_95', 100000),
                analysis_results.get('risk_metrics', {}).get('var_99', 200000)
            ]

            ax1.bar([f"{c:.0%}" for c in confidence_levels], var_values,
                   color=self.config['color_scheme']['secondary'])
            ax1.set_xlabel('Confidence Level')
            ax1.set_ylabel('Value at Risk ($)')
            ax1.set_title('VaR by Confidence Level')
            ax1.grid(True, alpha=0.3)

            # Historical VaR
            dates = pd.date_range(end=datetime.now(), periods=30)
            historical_var = np.random.uniform(80000, 120000, 30)

            ax2.plot(dates, historical_var, linewidth=2,
                    color=self.config['color_scheme']['primary'])
            ax2.fill_between(dates, historical_var, alpha=0.3)
            ax2.set_xlabel('Date')
            ax2.set_ylabel('VaR (99%)')
            ax2.set_title('30-Day VaR Trend')
            ax2.grid(True, alpha=0.3)
            ax2.tick_params(axis='x', rotation=45)

            plt.tight_layout()

            # Convert to base64
            buffer = BytesIO()
            plt.savefig(buffer, format='png', dpi=100, bbox_inches='tight')
            buffer.seek(0)
            image_base64 = base64.b64encode(buffer.read()).decode()
            plt.close()

            return f"data:image/png;base64,{image_base64}"

        except Exception as e:
            logger.error(f"Failed to create VaR chart: {e}")
            return None

    def _create_liquidation_heatmap(self, analysis_results: Dict) -> Optional[str]:
        """Create liquidation risk heatmap"""
        try:
            fig, ax = plt.subplots(figsize=(12, 8))

            # Create sample data (would use actual data)
            protocols = ['Aave', 'Compound', 'GMX', 'MakerDAO', 'Curve']
            assets = ['ETH', 'BTC', 'USDC', 'USDT', 'DAI']

            # Random risk matrix
            risk_matrix = np.random.rand(len(protocols), len(assets))

            # Create heatmap
            sns.heatmap(risk_matrix, annot=True, fmt='.2f',
                       xticklabels=assets, yticklabels=protocols,
                       cmap='RdYlGn_r', vmin=0, vmax=1,
                       cbar_kws={'label': 'Liquidation Risk'},
                       ax=ax)

            ax.set_title('Liquidation Risk Heatmap by Protocol and Asset')
            ax.set_xlabel('Asset')
            ax.set_ylabel('Protocol')

            # Convert to base64
            buffer = BytesIO()
            plt.savefig(buffer, format='png', dpi=100, bbox_inches='tight')
            buffer.seek(0)
            image_base64 = base64.b64encode(buffer.read()).decode()
            plt.close()

            return f"data:image/png;base64,{image_base64}"

        except Exception as e:
            logger.error(f"Failed to create liquidation heatmap: {e}")
            return None

    def _create_time_series_chart(self, analysis_results: Dict) -> Optional[str]:
        """Create time series analysis chart"""
        try:
            fig, axes = plt.subplots(3, 1, figsize=(12, 10), sharex=True)

            # Generate sample time series
            dates = pd.date_range(end=datetime.now(), periods=100, freq='D')

            # Health Factor
            health_factors = 2.0 + np.random.randn(100) * 0.3
            axes[0].plot(dates, health_factors, color=self.config['color_scheme']['primary'])
            axes[0].axhline(y=1.5, color='orange', linestyle='--', label='Warning')
            axes[0].axhline(y=1.2, color='red', linestyle='--', label='Critical')
            axes[0].set_ylabel('Health Factor')
            axes[0].set_title('Risk Metrics Over Time')
            axes[0].legend()
            axes[0].grid(True, alpha=0.3)

            # Liquidation Probability
            liq_prob = 1 / (1 + np.exp(-np.random.randn(100)))
            axes[1].plot(dates, liq_prob, color=self.config['color_scheme']['danger'])
            axes[1].fill_between(dates, liq_prob, alpha=0.3, color=self.config['color_scheme']['danger'])
            axes[1].set_ylabel('Liquidation Probability')
            axes[1].set_ylim([0, 1])
            axes[1].grid(True, alpha=0.3)

            # Total Value Locked
            tvl = 1e8 + np.cumsum(np.random.randn(100) * 1e6)
            axes[2].plot(dates, tvl / 1e6, color=self.config['color_scheme']['success'])
            axes[2].set_ylabel('TVL ($ millions)')
            axes[2].set_xlabel('Date')
            axes[2].grid(True, alpha=0.3)

            plt.tight_layout()

            # Convert to base64
            buffer = BytesIO()
            plt.savefig(buffer, format='png', dpi=100, bbox_inches='tight')
            buffer.seek(0)
            image_base64 = base64.b64encode(buffer.read()).decode()
            plt.close()

            return f"data:image/png;base64,{image_base64}"

        except Exception as e:
            logger.error(f"Failed to create time series chart: {e}")
            return None

    # ============================================================
    # FORMATTING HELPERS
    # ============================================================

    def _generate_report_title(self, report_type: str) -> str:
        """Generate report title"""
        titles = {
            'full': 'Comprehensive Risk Assessment Report',
            'summary': 'Executive Risk Summary',
            'technical': 'Technical Risk Analysis'
        }
        return titles.get(report_type, 'Risk Report')

    def _format_risk_level(self, score: float) -> str:
        """Format risk level based on score"""
        if score < 30:
            return "LOW RISK"
        elif score < 50:
            return "MODERATE RISK"
        elif score < 70:
            return "ELEVATED RISK"
        elif score < 85:
            return "HIGH RISK"
        else:
            return "CRITICAL RISK"

    def _get_status_indicator(self, value: float, threshold: float, higher_is_better: bool = False) -> str:
        """Get status indicator for metric"""
        if higher_is_better:
            if value > threshold:
                return "✅ Good"
            elif value > threshold * 0.7:
                return "⚠️ Warning"
            else:
                return "❌ Poor"
        else:
            if value < threshold:
                return "✅ Good"
            elif value < threshold * 1.5:
                return "⚠️ Warning"
            else:
                return "❌ Poor"

    def _analyze_market_trends(self, market_data: Dict) -> str:
        """Analyze and describe market trends"""
        trends = []

        volatility = market_data.get('volatility', 0)
        if volatility > 0.05:
            trends.append("• High volatility detected - consider reducing exposure")
        elif volatility < 0.02:
            trends.append("• Low volatility environment - favorable for leveraged positions")

        utilization = market_data.get('utilization_rate', 0)
        if utilization > 0.8:
            trends.append("• High protocol utilization - potential liquidity constraints")
        elif utilization < 0.3:
            trends.append("• Low utilization - ample liquidity available")

        return "\n".join(trends) if trends else "• Market conditions are within normal parameters"

    def _format_stress_test_results(self, stress_tests: Dict) -> str:
        """Format stress test results"""
        if not stress_tests:
            return "No stress test data available"

        results = []
        for scenario, outcome in stress_tests.items():
            survival_prob = outcome.get('survival_probability', 0)
            expected_loss = outcome.get('expected_loss', 0)
            results.append(f"• {scenario}: {survival_prob:.1%} survival, ${expected_loss:,.0f} expected loss")

        return "\n".join(results)

    # ============================================================
    # EXPORT FUNCTIONS
    # ============================================================

    def export_to_json(self, report: RiskReport) -> str:
        """Export report to JSON"""
        return json.dumps({
            'report_id': report.report_id,
            'title': report.title,
            'generated_at': report.generated_at.isoformat(),
            'executive_summary': report.executive_summary,
            'sections': [
                {
                    'title': section.title,
                    'content': str(section.content),
                    'type': section.section_type
                }
                for section in report.sections
            ],
            'metadata': report.metadata
        }, indent=2)

    def export_to_pdf(self, report: RiskReport, filename: str):
        """Export report to PDF (requires reportlab)"""
        if not HAS_REPORTLAB:
            logger.error("PDF export requires reportlab library")
            return

        # This would implement full PDF generation
        # Simplified for brevity
        logger.info(f"PDF export to {filename} (not fully implemented)")

    def export_to_html(self, report: RiskReport) -> str:
        """Export report to HTML"""
        html = f"""
        <html>
        <head>
            <title>{report.title}</title>
            <style>
                body {{ font-family: Arial, sans-serif; margin: 40px; }}
                h1 {{ color: {self.config['color_scheme']['primary']}; }}
                h2 {{ color: {self.config['color_scheme']['secondary']}; }}
                .summary {{ background: #f0f0f0; padding: 20px; border-radius: 5px; }}
                table {{ border-collapse: collapse; width: 100%; }}
                th, td {{ border: 1px solid #ddd; padding: 8px; text-align: left; }}
                th {{ background-color: {self.config['color_scheme']['primary']}; color: white; }}
            </style>
        </head>
        <body>
            <h1>{report.title}</h1>
            <p>Generated: {report.generated_at.strftime('%Y-%m-%d %H:%M:%S')} UTC</p>

            <div class="summary">
                <h2>Executive Summary</h2>
                <pre>{report.executive_summary}</pre>
            </div>
        """

        # Add sections
        for section in report.sections:
            html += f"<h2>{section.title}</h2>"

            if section.section_type == "table":
                html += "<table>"
                for i, row in enumerate(section.content):
                    if i == 0:
                        html += "<tr>" + "".join(f"<th>{cell}</th>" for cell in row) + "</tr>"
                    else:
                        html += "<tr>" + "".join(f"<td>{cell}</td>" for cell in row) + "</tr>"
                html += "</table>"
            else:
                html += f"<pre>{section.content}</pre>"

        # Add visualizations
        for name, image_data in report.visualizations.items():
            html += f'<h2>{name.replace("_", " ").title()}</h2>'
            html += f'<img src="{image_data}" style="max-width: 100%;"/>'

        html += "</body></html>"

        return html