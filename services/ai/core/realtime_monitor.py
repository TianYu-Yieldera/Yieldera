"""
Real-time Monitoring System
24/7 risk monitoring with alerts and automated responses
Phase 2: Core Implementation
"""

import asyncio
import logging
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any, Set
from dataclasses import dataclass, field
from enum import Enum
import json
from collections import defaultdict, deque
import aiohttp
import websockets

logger = logging.getLogger(__name__)

# ============================================================
# DATA MODELS
# ============================================================

class AlertLevel(Enum):
    """Alert severity levels"""
    INFO = "info"
    WARNING = "warning"
    DANGER = "danger"
    CRITICAL = "critical"

class AlertType(Enum):
    """Types of alerts"""
    LIQUIDATION_RISK = "liquidation_risk"
    SYSTEMIC_RISK = "systemic_risk"
    UNUSUAL_ACTIVITY = "unusual_activity"
    MARKET_CRASH = "market_crash"
    HIGH_GAS = "high_gas"
    PROTOCOL_ISSUE = "protocol_issue"
    ORACLE_DEVIATION = "oracle_deviation"
    LIQUIDITY_CRISIS = "liquidity_crisis"

@dataclass
class Alert:
    """Alert message"""
    id: str
    level: AlertLevel
    alert_type: AlertType
    title: str
    message: str
    position_id: Optional[str]
    user_address: Optional[str]
    timestamp: datetime
    data: Dict[str, Any]
    action_required: str
    auto_response: Optional[str] = None

@dataclass
class MonitoringRule:
    """Rule for monitoring triggers"""
    name: str
    condition: str  # Expression to evaluate
    threshold: float
    alert_level: AlertLevel
    alert_type: AlertType
    cooldown_minutes: int = 5
    auto_action: Optional[str] = None

@dataclass
class MonitoringMetrics:
    """Current monitoring metrics"""
    timestamp: datetime
    positions_monitored: int
    alerts_triggered: int
    alerts_by_level: Dict[str, int]
    avg_health_factor: float
    at_risk_positions: int
    total_value_at_risk: float
    system_health: str
    response_time_ms: float

# ============================================================
# REAL-TIME MONITOR
# ============================================================

class GauntletRealTimeMonitor:
    """
    Real-time monitoring system for 24/7 risk surveillance
    """

    def __init__(
        self,
        data_infrastructure,
        risk_calculator,
        ml_predictor,
        config: Optional[Dict] = None
    ):
        """
        Initialize monitor

        Args:
            data_infrastructure: Data access layer
            risk_calculator: Risk calculation engine
            ml_predictor: ML prediction model
            config: Monitoring configuration
        """
        self.data_infrastructure = data_infrastructure
        self.risk_calculator = risk_calculator
        self.ml_predictor = ml_predictor
        self.config = config or self._default_config()

        # Monitoring state
        self.monitoring_active = False
        self.monitored_positions: Set[str] = set()
        self.monitored_users: Set[str] = set()
        self.alert_history: deque = deque(maxlen=1000)
        self.alert_cooldowns: Dict[str, datetime] = {}
        self.metrics_history: deque = deque(maxlen=100)

        # Alert channels
        self.alert_channels = []
        self.websocket_connections = set()

        # Monitoring rules
        self.rules = self._initialize_rules()

        # Performance tracking
        self.performance_metrics = defaultdict(list)

    def _default_config(self) -> Dict:
        """Default monitoring configuration"""
        return {
            'monitoring_interval': 1,  # seconds
            'batch_size': 100,
            'alert_channels': ['log', 'webhook', 'websocket'],
            'webhook_url': None,
            'thresholds': {
                'critical_health_factor': 1.05,
                'warning_health_factor': 1.2,
                'high_liquidation_probability': 0.7,
                'systemic_risk_threshold': 0.6,
                'unusual_activity_zscore': 3.0,
                'high_gas_threshold': 200,  # Gwei
                'oracle_deviation_threshold': 0.05  # 5%
            },
            'auto_actions': {
                'enable_auto_deleveraging': False,
                'enable_auto_hedging': False,
                'enable_emergency_shutdown': False
            },
            'monitoring_scope': {
                'protocols': ['aave', 'compound', 'gmx'],
                'min_position_value': 1000,  # USD
                'max_positions': 10000
            }
        }

    def _initialize_rules(self) -> List[MonitoringRule]:
        """Initialize monitoring rules"""
        thresholds = self.config['thresholds']

        return [
            MonitoringRule(
                name="Critical Health Factor",
                condition="health_factor < threshold",
                threshold=thresholds['critical_health_factor'],
                alert_level=AlertLevel.CRITICAL,
                alert_type=AlertType.LIQUIDATION_RISK,
                cooldown_minutes=1,
                auto_action="emergency_deleverage" if self.config['auto_actions']['enable_auto_deleveraging'] else None
            ),
            MonitoringRule(
                name="Warning Health Factor",
                condition="health_factor < threshold",
                threshold=thresholds['warning_health_factor'],
                alert_level=AlertLevel.WARNING,
                alert_type=AlertType.LIQUIDATION_RISK,
                cooldown_minutes=5
            ),
            MonitoringRule(
                name="High Liquidation Probability",
                condition="liquidation_probability > threshold",
                threshold=thresholds['high_liquidation_probability'],
                alert_level=AlertLevel.DANGER,
                alert_type=AlertType.LIQUIDATION_RISK,
                cooldown_minutes=5
            ),
            MonitoringRule(
                name="Systemic Risk",
                condition="systemic_risk > threshold",
                threshold=thresholds['systemic_risk_threshold'],
                alert_level=AlertLevel.DANGER,
                alert_type=AlertType.SYSTEMIC_RISK,
                cooldown_minutes=15
            ),
            MonitoringRule(
                name="Unusual Activity",
                condition="activity_zscore > threshold",
                threshold=thresholds['unusual_activity_zscore'],
                alert_level=AlertLevel.WARNING,
                alert_type=AlertType.UNUSUAL_ACTIVITY,
                cooldown_minutes=10
            ),
            MonitoringRule(
                name="High Gas Price",
                condition="gas_price > threshold",
                threshold=thresholds['high_gas_threshold'],
                alert_level=AlertLevel.INFO,
                alert_type=AlertType.HIGH_GAS,
                cooldown_minutes=30
            ),
            MonitoringRule(
                name="Oracle Price Deviation",
                condition="oracle_deviation > threshold",
                threshold=thresholds['oracle_deviation_threshold'],
                alert_level=AlertLevel.DANGER,
                alert_type=AlertType.ORACLE_DEVIATION,
                cooldown_minutes=5
            )
        ]

    # ============================================================
    # MONITORING LIFECYCLE
    # ============================================================

    async def start_monitoring(self):
        """Start real-time monitoring"""
        if self.monitoring_active:
            logger.warning("Monitoring already active")
            return

        self.monitoring_active = True
        logger.info("Starting real-time monitoring system...")

        # Start monitoring tasks
        tasks = [
            asyncio.create_task(self._monitoring_loop()),
            asyncio.create_task(self._health_check_loop()),
            asyncio.create_task(self._metrics_aggregation_loop())
        ]

        try:
            await asyncio.gather(*tasks)
        except Exception as e:
            logger.error(f"Monitoring error: {e}")
        finally:
            self.monitoring_active = False

    async def stop_monitoring(self):
        """Stop monitoring"""
        logger.info("Stopping monitoring system...")
        self.monitoring_active = False

        # Close WebSocket connections
        for ws in self.websocket_connections:
            await ws.close()

    async def _monitoring_loop(self):
        """Main monitoring loop"""
        interval = self.config['monitoring_interval']

        while self.monitoring_active:
            try:
                start_time = datetime.utcnow()

                # Monitor positions
                await self._monitor_positions()

                # Monitor market conditions
                await self._monitor_market()

                # Monitor protocol health
                await self._monitor_protocols()

                # Track performance
                elapsed = (datetime.utcnow() - start_time).total_seconds() * 1000
                self.performance_metrics['response_time'].append(elapsed)

                # Sleep for remaining interval
                sleep_time = max(0, interval - elapsed / 1000)
                await asyncio.sleep(sleep_time)

            except Exception as e:
                logger.error(f"Monitoring loop error: {e}")
                await asyncio.sleep(interval)

    async def _health_check_loop(self):
        """System health check loop"""
        while self.monitoring_active:
            try:
                # Check data source health
                data_health = await self._check_data_health()

                # Check model health
                model_health = await self._check_model_health()

                # Check alert system health
                alert_health = await self._check_alert_health()

                # Overall system health
                if all([data_health, model_health, alert_health]):
                    system_health = "healthy"
                elif any([not data_health, not model_health, not alert_health]):
                    system_health = "degraded"
                else:
                    system_health = "critical"

                # Log health status
                if system_health != "healthy":
                    logger.warning(f"System health: {system_health}")

                await asyncio.sleep(60)  # Check every minute

            except Exception as e:
                logger.error(f"Health check error: {e}")
                await asyncio.sleep(60)

    async def _metrics_aggregation_loop(self):
        """Aggregate and store metrics"""
        while self.monitoring_active:
            try:
                metrics = await self._calculate_metrics()
                self.metrics_history.append(metrics)

                # Broadcast metrics to WebSocket clients
                await self._broadcast_metrics(metrics)

                await asyncio.sleep(10)  # Update every 10 seconds

            except Exception as e:
                logger.error(f"Metrics aggregation error: {e}")
                await asyncio.sleep(10)

    # ============================================================
    # MONITORING FUNCTIONS
    # ============================================================

    async def _monitor_positions(self):
        """Monitor all tracked positions"""
        # Get high-risk positions
        positions = await self.data_infrastructure.get_high_risk_positions(
            limit=self.config['monitoring_scope']['max_positions']
        )

        monitored_count = 0
        alerts_triggered = []

        for position in positions:
            try:
                # Skip small positions
                if position.get('collateral_value_usd', 0) < self.config['monitoring_scope']['min_position_value']:
                    continue

                # Check position risk
                alerts = await self._check_position_risk(position)
                if alerts:
                    alerts_triggered.extend(alerts)

                monitored_count += 1

                # Add to monitored set
                self.monitored_positions.add(position.get('position_id', ''))
                self.monitored_users.add(position.get('user_address', ''))

            except Exception as e:
                logger.error(f"Error monitoring position {position.get('position_id')}: {e}")

        # Process alerts
        for alert in alerts_triggered:
            await self._process_alert(alert)

        logger.debug(f"Monitored {monitored_count} positions, triggered {len(alerts_triggered)} alerts")

    async def _monitor_market(self):
        """Monitor market conditions"""
        try:
            # Get market data
            market_data = await self.data_infrastructure.collect_comprehensive_data()

            # Check market conditions
            alerts = []

            # Check for market crash
            for asset, volatility in market_data.volatility_metrics.items():
                if volatility > 0.1:  # 10% volatility
                    alert = self._create_alert(
                        level=AlertLevel.DANGER,
                        alert_type=AlertType.MARKET_CRASH,
                        title=f"High volatility detected for {asset}",
                        message=f"Volatility: {volatility:.2%}",
                        data={'asset': asset, 'volatility': volatility}
                    )
                    alerts.append(alert)

            # Check liquidity
            for market_key, depth in market_data.market_depth.items():
                if depth.get('available_liquidity', float('inf')) < 1e6:  # Less than $1M
                    alert = self._create_alert(
                        level=AlertLevel.WARNING,
                        alert_type=AlertType.LIQUIDITY_CRISIS,
                        title=f"Low liquidity in {market_key}",
                        message=f"Available: ${depth['available_liquidity']:,.0f}",
                        data={'market': market_key, 'liquidity': depth['available_liquidity']}
                    )
                    alerts.append(alert)

            # Process alerts
            for alert in alerts:
                await self._process_alert(alert)

        except Exception as e:
            logger.error(f"Market monitoring error: {e}")

    async def _monitor_protocols(self):
        """Monitor protocol health"""
        for protocol in self.config['monitoring_scope']['protocols']:
            try:
                # Check protocol-specific metrics
                # This would integrate with protocol-specific monitoring
                pass
            except Exception as e:
                logger.error(f"Protocol monitoring error for {protocol}: {e}")

    async def _check_position_risk(self, position: Dict) -> List[Alert]:
        """Check risk for a single position"""
        alerts = []

        # Calculate risk metrics
        risk_assessment = self.risk_calculator.calculate_position_risk(
            position,
            {'volatility': {}, 'market_depth': {}}  # Simplified
        )

        # ML prediction
        if self.ml_predictor and self.ml_predictor.is_trained:
            prediction = self.ml_predictor.predict(
                position,
                {'volatility': 0.03, 'utilization_rate': 0.5}
            )
            liquidation_probability = prediction.liquidation_probability
        else:
            liquidation_probability = risk_assessment.risk_metrics.liquidation_probability

        # Check rules
        for rule in self.rules:
            if self._should_trigger_rule(rule, position, risk_assessment, liquidation_probability):
                alert = self._create_position_alert(
                    rule,
                    position,
                    risk_assessment,
                    liquidation_probability
                )
                if alert:
                    alerts.append(alert)

        return alerts

    def _should_trigger_rule(
        self,
        rule: MonitoringRule,
        position: Dict,
        risk_assessment,
        liquidation_probability: float
    ) -> bool:
        """Check if rule should trigger"""
        # Check cooldown
        cooldown_key = f"{rule.name}_{position.get('position_id', '')}"
        if cooldown_key in self.alert_cooldowns:
            if datetime.utcnow() < self.alert_cooldowns[cooldown_key]:
                return False

        # Evaluate condition
        if rule.name == "Critical Health Factor":
            return risk_assessment.risk_metrics.health_factor < rule.threshold
        elif rule.name == "Warning Health Factor":
            return risk_assessment.risk_metrics.health_factor < rule.threshold
        elif rule.name == "High Liquidation Probability":
            return liquidation_probability > rule.threshold
        elif rule.name == "High Gas Price":
            return False  # Would check actual gas price

        return False

    def _create_position_alert(
        self,
        rule: MonitoringRule,
        position: Dict,
        risk_assessment,
        liquidation_probability: float
    ) -> Optional[Alert]:
        """Create alert for position"""
        alert = self._create_alert(
            level=rule.alert_level,
            alert_type=rule.alert_type,
            title=f"{rule.name} triggered for position",
            message=self._format_alert_message(rule, position, risk_assessment),
            position_id=position.get('position_id'),
            user_address=position.get('user_address'),
            data={
                'health_factor': risk_assessment.risk_metrics.health_factor,
                'liquidation_probability': liquidation_probability,
                'collateral_value': position.get('collateral_value_usd', 0),
                'debt_value': position.get('debt_value_usd', 0)
            }
        )

        # Set action required
        if rule.alert_level == AlertLevel.CRITICAL:
            alert.action_required = "Immediate action required - add collateral or close position"
        elif rule.alert_level == AlertLevel.DANGER:
            alert.action_required = "Urgent - reduce risk within 1 hour"
        elif rule.alert_level == AlertLevel.WARNING:
            alert.action_required = "Monitor closely and consider risk reduction"
        else:
            alert.action_required = "No immediate action required"

        # Set auto response if configured
        if rule.auto_action:
            alert.auto_response = rule.auto_action

        # Update cooldown
        cooldown_key = f"{rule.name}_{position.get('position_id', '')}"
        self.alert_cooldowns[cooldown_key] = datetime.utcnow() + timedelta(minutes=rule.cooldown_minutes)

        return alert

    def _create_alert(
        self,
        level: AlertLevel,
        alert_type: AlertType,
        title: str,
        message: str,
        position_id: Optional[str] = None,
        user_address: Optional[str] = None,
        data: Optional[Dict] = None
    ) -> Alert:
        """Create alert object"""
        alert_id = f"alert_{datetime.utcnow().timestamp()}_{len(self.alert_history)}"

        return Alert(
            id=alert_id,
            level=level,
            alert_type=alert_type,
            title=title,
            message=message,
            position_id=position_id,
            user_address=user_address,
            timestamp=datetime.utcnow(),
            data=data or {},
            action_required=""
        )

    def _format_alert_message(
        self,
        rule: MonitoringRule,
        position: Dict,
        risk_assessment
    ) -> str:
        """Format alert message"""
        protocol = position.get('protocol', 'Unknown')
        health_factor = risk_assessment.risk_metrics.health_factor
        collateral = position.get('collateral_value_usd', 0)

        return (
            f"Protocol: {protocol} | "
            f"Health Factor: {health_factor:.2f} | "
            f"Collateral: ${collateral:,.0f}"
        )

    # ============================================================
    # ALERT PROCESSING
    # ============================================================

    async def _process_alert(self, alert: Alert):
        """Process and distribute alert"""
        # Add to history
        self.alert_history.append(alert)

        # Log alert
        await self._log_alert(alert)

        # Send to channels
        for channel in self.config['alert_channels']:
            try:
                if channel == 'webhook':
                    await self._send_webhook_alert(alert)
                elif channel == 'websocket':
                    await self._broadcast_alert(alert)
            except Exception as e:
                logger.error(f"Failed to send alert via {channel}: {e}")

        # Execute auto response if configured
        if alert.auto_response:
            await self._execute_auto_response(alert)

    async def _log_alert(self, alert: Alert):
        """Log alert to system"""
        log_level = {
            AlertLevel.INFO: logging.INFO,
            AlertLevel.WARNING: logging.WARNING,
            AlertLevel.DANGER: logging.ERROR,
            AlertLevel.CRITICAL: logging.CRITICAL
        }.get(alert.level, logging.INFO)

        logger.log(
            log_level,
            f"[{alert.level.value.upper()}] {alert.title}: {alert.message}"
        )

    async def _send_webhook_alert(self, alert: Alert):
        """Send alert via webhook"""
        if not self.config.get('webhook_url'):
            return

        async with aiohttp.ClientSession() as session:
            payload = {
                'id': alert.id,
                'level': alert.level.value,
                'type': alert.alert_type.value,
                'title': alert.title,
                'message': alert.message,
                'timestamp': alert.timestamp.isoformat(),
                'data': alert.data,
                'action_required': alert.action_required
            }

            async with session.post(
                self.config['webhook_url'],
                json=payload,
                timeout=aiohttp.ClientTimeout(total=5)
            ) as response:
                if response.status != 200:
                    logger.error(f"Webhook alert failed: {response.status}")

    async def _broadcast_alert(self, alert: Alert):
        """Broadcast alert to WebSocket clients"""
        if not self.websocket_connections:
            return

        message = json.dumps({
            'type': 'alert',
            'alert': {
                'id': alert.id,
                'level': alert.level.value,
                'alert_type': alert.alert_type.value,
                'title': alert.title,
                'message': alert.message,
                'timestamp': alert.timestamp.isoformat(),
                'data': alert.data
            }
        })

        # Send to all connected clients
        disconnected = set()
        for ws in self.websocket_connections:
            try:
                await ws.send(message)
            except Exception:
                disconnected.add(ws)

        # Remove disconnected clients
        self.websocket_connections -= disconnected

    async def _broadcast_metrics(self, metrics: MonitoringMetrics):
        """Broadcast metrics to WebSocket clients"""
        if not self.websocket_connections:
            return

        message = json.dumps({
            'type': 'metrics',
            'metrics': {
                'timestamp': metrics.timestamp.isoformat(),
                'positions_monitored': metrics.positions_monitored,
                'alerts_triggered': metrics.alerts_triggered,
                'at_risk_positions': metrics.at_risk_positions,
                'total_value_at_risk': metrics.total_value_at_risk,
                'system_health': metrics.system_health
            }
        })

        # Send to all connected clients
        disconnected = set()
        for ws in self.websocket_connections:
            try:
                await ws.send(message)
            except Exception:
                disconnected.add(ws)

        self.websocket_connections -= disconnected

    async def _execute_auto_response(self, alert: Alert):
        """Execute automated response to alert"""
        logger.info(f"Executing auto response: {alert.auto_response}")

        if alert.auto_response == "emergency_deleverage":
            # Would trigger emergency deleveraging
            logger.warning("Emergency deleveraging triggered (simulation)")
        elif alert.auto_response == "hedge_position":
            # Would create hedging position
            logger.info("Hedging position created (simulation)")

    # ============================================================
    # HEALTH CHECKS
    # ============================================================

    async def _check_data_health(self) -> bool:
        """Check data source health"""
        try:
            # Check database connection
            test_data = await self.data_infrastructure.get_high_risk_positions(limit=1)
            return True
        except Exception:
            return False

    async def _check_model_health(self) -> bool:
        """Check ML model health"""
        if self.ml_predictor:
            return self.ml_predictor.is_trained
        return True  # No model is OK

    async def _check_alert_health(self) -> bool:
        """Check alert system health"""
        # Check if alerts are being processed
        recent_alerts = [
            alert for alert in self.alert_history
            if (datetime.utcnow() - alert.timestamp).seconds < 300
        ]
        return True  # Simplified

    # ============================================================
    # METRICS
    # ============================================================

    async def _calculate_metrics(self) -> MonitoringMetrics:
        """Calculate current monitoring metrics"""
        # Count alerts by level
        alerts_by_level = defaultdict(int)
        for alert in self.alert_history:
            if (datetime.utcnow() - alert.timestamp).seconds < 3600:  # Last hour
                alerts_by_level[alert.level.value] += 1

        # Calculate average response time
        avg_response_time = (
            sum(self.performance_metrics['response_time'][-100:]) /
            len(self.performance_metrics['response_time'][-100:])
            if self.performance_metrics['response_time'] else 0
        )

        # Determine system health
        if alerts_by_level.get('critical', 0) > 0:
            system_health = "critical"
        elif alerts_by_level.get('danger', 0) > 5:
            system_health = "degraded"
        else:
            system_health = "healthy"

        return MonitoringMetrics(
            timestamp=datetime.utcnow(),
            positions_monitored=len(self.monitored_positions),
            alerts_triggered=len(self.alert_history),
            alerts_by_level=dict(alerts_by_level),
            avg_health_factor=2.0,  # Would calculate from actual data
            at_risk_positions=sum(1 for a in self.alert_history if a.alert_type == AlertType.LIQUIDATION_RISK),
            total_value_at_risk=0,  # Would calculate from actual data
            system_health=system_health,
            response_time_ms=avg_response_time
        )

    # ============================================================
    # WEBSOCKET MANAGEMENT
    # ============================================================

    async def add_websocket_connection(self, websocket):
        """Add WebSocket connection for real-time updates"""
        self.websocket_connections.add(websocket)
        logger.info(f"WebSocket client connected. Total: {len(self.websocket_connections)}")

        # Send current metrics
        metrics = await self._calculate_metrics()
        await self._broadcast_metrics(metrics)

    async def remove_websocket_connection(self, websocket):
        """Remove WebSocket connection"""
        self.websocket_connections.discard(websocket)
        logger.info(f"WebSocket client disconnected. Total: {len(self.websocket_connections)}")