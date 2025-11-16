"""
AI Risk Engine API Service - Version 2
Enhanced with Phase 2 capabilities: Simulation, ML, Monitoring
"""

from fastapi import FastAPI, HTTPException, BackgroundTasks, Depends, WebSocket, WebSocketDisconnect
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse, HTMLResponse
from pydantic import BaseModel, Field
from typing import Dict, List, Optional, Any, Set
from datetime import datetime, timedelta
import asyncio
import logging
import os
import json
from contextlib import asynccontextmanager

# Import core modules
import sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from core.data_infrastructure import DataInfrastructure
from core.risk_calculator import RiskCalculator
from core.ml_data_pipeline import MLDataPipeline
from core.agent_simulation import GauntletSimulator, AgentType
from core.economic_model import GauntletEconomicModel, EconomicParameters
from core.ml_predictor import HybridMLPredictor
from core.realtime_monitor import GauntletRealTimeMonitor
from core.report_generator import GauntletReportGenerator

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# ============================================================
# PYDANTIC MODELS
# ============================================================

class SimulationRequest(BaseModel):
    """Request for running simulation"""
    scenarios: List[Dict[str, Any]] = Field(..., description="Simulation scenarios")
    num_agents: int = Field(10000, description="Number of agents")
    time_steps: int = Field(1000, description="Number of time steps")

class SimulationResponse(BaseModel):
    """Simulation results"""
    scenario_results: Dict[str, Any]
    risk_metrics: Dict[str, float]
    recommendations: List[str]
    execution_time: float

class MLTrainingRequest(BaseModel):
    """Request for ML model training"""
    training_data_path: Optional[str] = Field(None, description="Path to training data")
    model_type: str = Field("hybrid", description="Model type: xgboost, lstm, or hybrid")
    epochs: int = Field(100, description="Training epochs")

class MLPredictionRequest(BaseModel):
    """Request for ML prediction"""
    position_data: Dict[str, Any]
    market_data: Dict[str, Any]
    historical_data: Optional[List[float]] = None

class OptimizationRequest(BaseModel):
    """Request for parameter optimization"""
    current_parameters: Dict[str, float]
    market_conditions: Dict[str, Any]
    constraints: Optional[Dict[str, float]] = None

class ReportRequest(BaseModel):
    """Request for report generation"""
    report_type: str = Field("full", description="Report type: full, summary, or technical")
    include_visualizations: bool = Field(True)
    export_format: str = Field("json", description="Format: json, html, or pdf")

# ============================================================
# GLOBAL INSTANCES
# ============================================================

data_infrastructure: Optional[DataInfrastructure] = None
risk_calculator: Optional[RiskCalculator] = None
ml_pipeline: Optional[MLDataPipeline] = None
simulator: Optional[GauntletSimulator] = None
economic_model: Optional[GauntletEconomicModel] = None
ml_predictor: Optional[HybridMLPredictor] = None
monitor: Optional[GauntletRealTimeMonitor] = None
report_generator: Optional[GauntletReportGenerator] = None

# WebSocket connections
websocket_connections: Set[WebSocket] = set()

# ============================================================
# LIFESPAN MANAGEMENT
# ============================================================

@asynccontextmanager
async def lifespan(app: FastAPI):
    """Manage application lifecycle"""
    global data_infrastructure, risk_calculator, ml_pipeline
    global simulator, economic_model, ml_predictor, monitor, report_generator

    # Startup
    logger.info("Starting AI Risk Engine API v2...")

    # Initialize Phase 1 components
    db_config = {
        'host': os.getenv('DB_HOST', 'localhost'),
        'port': int(os.getenv('DB_PORT', 5432)),
        'user': os.getenv('DB_USER', 'postgres'),
        'password': os.getenv('DB_PASSWORD', 'postgres'),
        'database': os.getenv('DB_NAME', 'loyalty_points')
    }

    redis_config = {
        'host': os.getenv('REDIS_HOST', 'localhost'),
        'port': int(os.getenv('REDIS_PORT', 6379)),
        'password': os.getenv('REDIS_PASSWORD')
    }

    data_infrastructure = DataInfrastructure(db_config, redis_config)
    await data_infrastructure.initialize()

    risk_calculator = RiskCalculator()
    ml_pipeline = MLDataPipeline()

    # Initialize Phase 2 components
    simulator = GauntletSimulator(num_agents=10000)
    economic_model = GauntletEconomicModel()
    ml_predictor = HybridMLPredictor()

    # Load pre-trained models if available
    model_path = os.getenv('MODEL_PATH', '/models')
    if os.path.exists(f"{model_path}/xgboost_model.bin"):
        try:
            ml_predictor.load_models(model_path)
            logger.info("Pre-trained models loaded successfully")
        except Exception as e:
            logger.warning(f"Could not load pre-trained models: {e}")

    # Initialize monitoring
    monitor = GauntletRealTimeMonitor(
        data_infrastructure,
        risk_calculator,
        ml_predictor
    )

    # Start monitoring in background
    asyncio.create_task(monitor.start_monitoring())

    # Initialize report generator
    report_generator = GauntletReportGenerator()

    logger.info("AI Risk Engine API v2 started successfully")

    yield

    # Shutdown
    logger.info("Shutting down AI Risk Engine API v2...")

    # Stop monitoring
    if monitor:
        await monitor.stop_monitoring()

    # Close connections
    if data_infrastructure:
        await data_infrastructure.close()

    logger.info("AI Risk Engine API v2 shut down complete")

# ============================================================
# FASTAPI APP
# ============================================================

app = FastAPI(
    title="AI Risk Engine API v2",
    description="Gauntlet-standard AI risk assessment with simulation, ML, and real-time monitoring",
    version="2.0.0",
    lifespan=lifespan
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Note: Phase 1 basic API endpoints (/api/risk/*, /api/user/*, etc.)
# are defined below after Phase 2 endpoints for better organization

# ============================================================
# SIMULATION ENDPOINTS
# ============================================================

@app.post("/api/v2/simulate", response_model=SimulationResponse)
async def run_simulation(request: SimulationRequest):
    """
    Run comprehensive market simulation

    Simulates market conditions with thousands of agents
    """
    try:
        start_time = datetime.utcnow()

        # Run simulation
        results = await simulator.run_comprehensive_simulation(
            scenarios=request.scenarios,
            time_steps=request.time_steps
        )

        # Extract key metrics
        risk_metrics = {}
        for scenario_name, result in results['results'].items():
            risk_metrics[scenario_name] = {
                'systemic_risk': result.risk_metrics['systemic_risk'],
                'total_liquidations': result.risk_metrics['total_liquidations'],
                'avg_volatility': result.risk_metrics['avg_volatility']
            }

        execution_time = (datetime.utcnow() - start_time).total_seconds()

        return SimulationResponse(
            scenario_results=results['results'],
            risk_metrics=risk_metrics,
            recommendations=results['recommendations'],
            execution_time=execution_time
        )

    except Exception as e:
        logger.error(f"Simulation failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/v2/simulate/scenarios")
async def get_simulation_scenarios():
    """Get predefined simulation scenarios"""
    return {
        "scenarios": [
            {
                "name": "Market Crash",
                "description": "30% market decline",
                "market_shock": -0.30,
                "volatility_multiplier": 2.0
            },
            {
                "name": "Flash Crash",
                "description": "Sudden 50% drop and recovery",
                "market_shock": -0.50,
                "volatility_multiplier": 5.0,
                "recovery_time": 100
            },
            {
                "name": "Liquidity Crisis",
                "description": "50% liquidity reduction",
                "liquidity_shock": -0.50,
                "utilization_spike": 0.95
            },
            {
                "name": "Black Swan",
                "description": "Extreme tail event",
                "market_shock": -0.70,
                "volatility_multiplier": 10.0,
                "correlation": 0.95
            }
        ]
    }

# ============================================================
# MACHINE LEARNING ENDPOINTS
# ============================================================

@app.post("/api/v2/ml/train")
async def train_ml_model(request: MLTrainingRequest, background_tasks: BackgroundTasks):
    """
    Train machine learning models

    Trains XGBoost and/or LSTM models on historical data
    """
    try:
        # This would typically load data from database or file
        # For demo, using generated data
        import numpy as np

        # Generate sample training data
        n_samples = 10000
        n_features = 20

        X_structured = np.random.randn(n_samples, n_features)
        X_sequences = np.random.randn(n_samples, 60, 10)  # 60 timesteps, 10 features
        y = (np.random.randn(n_samples) > 0.5).astype(int)  # Binary classification

        training_data = {
            'X_structured': X_structured,
            'X_sequences': X_sequences if request.model_type in ['lstm', 'hybrid'] else None,
            'y': y
        }

        # Train in background
        background_tasks.add_task(
            ml_predictor.train,
            training_data,
            save_path="/models"
        )

        return {
            "status": "training_started",
            "model_type": request.model_type,
            "samples": n_samples,
            "message": "Model training initiated in background"
        }

    except Exception as e:
        logger.error(f"ML training failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/v2/ml/predict")
async def predict_liquidation(request: MLPredictionRequest):
    """
    Predict liquidation probability using ML models

    Returns probability, confidence, and recommendations
    """
    try:
        if not ml_predictor.is_trained:
            raise HTTPException(
                status_code=400,
                detail="Models not trained. Please train models first."
            )

        # Convert historical data if provided
        historical_array = None
        if request.historical_data:
            historical_array = np.array(request.historical_data).reshape(-1, 1)

        # Make prediction
        prediction = ml_predictor.predict(
            position_data=request.position_data,
            market_data=request.market_data,
            historical_data=historical_array
        )

        return {
            "liquidation_probability": prediction.liquidation_probability,
            "confidence_score": prediction.confidence_score,
            "time_to_liquidation": prediction.time_to_liquidation,
            "risk_factors": prediction.risk_factors,
            "recommendations": prediction.recommendations,
            "model_explanations": prediction.model_explanations
        }

    except Exception as e:
        logger.error(f"ML prediction failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/v2/ml/performance")
async def get_model_performance():
    """Get ML model performance metrics"""
    if not ml_predictor.is_trained:
        return {"status": "not_trained"}

    # This would return actual performance metrics
    return {
        "status": "trained",
        "metrics": {
            "accuracy": 0.85,
            "precision": 0.82,
            "recall": 0.78,
            "f1_score": 0.80,
            "auc_roc": 0.88
        },
        "feature_importance": ml_predictor._get_top_features(10)
    }

# ============================================================
# ECONOMIC MODEL ENDPOINTS
# ============================================================

@app.post("/api/v2/optimize/parameters")
async def optimize_parameters(request: OptimizationRequest):
    """
    Optimize protocol parameters

    Uses differential evolution to find optimal risk parameters
    """
    try:
        # Convert request to EconomicParameters
        current_params = EconomicParameters(
            ltv=request.current_parameters.get('ltv', 0.7),
            liquidation_threshold=request.current_parameters.get('liquidation_threshold', 0.8),
            liquidation_penalty=request.current_parameters.get('liquidation_penalty', 0.05),
            base_rate=request.current_parameters.get('base_rate', 0.02),
            slope1=request.current_parameters.get('slope1', 0.08),
            slope2=request.current_parameters.get('slope2', 1.0)
        )

        # Get historical data (simplified)
        import pandas as pd
        historical_data = pd.DataFrame({
            'ETH': np.random.randn(100) * 0.05 + 1,
            'BTC': np.random.randn(100) * 0.04 + 1
        })

        # Optimize
        result = economic_model.optimize_protocol_parameters(
            current_params=current_params,
            market_conditions=request.market_conditions,
            historical_data=historical_data,
            constraints=request.constraints
        )

        return {
            "current_parameters": {
                "ltv": result.current_params.ltv,
                "liquidation_threshold": result.current_params.liquidation_threshold,
                "liquidation_penalty": result.current_params.liquidation_penalty,
                "base_rate": result.current_params.base_rate,
                "slope1": result.current_params.slope1,
                "slope2": result.current_params.slope2
            },
            "optimal_parameters": {
                "ltv": result.optimal_params.ltv,
                "liquidation_threshold": result.optimal_params.liquidation_threshold,
                "liquidation_penalty": result.optimal_params.liquidation_penalty,
                "base_rate": result.optimal_params.base_rate,
                "slope1": result.optimal_params.slope1,
                "slope2": result.optimal_params.slope2
            },
            "expected_improvement": result.expected_improvement,
            "implementation_risk": result.implementation_risk,
            "confidence_score": result.confidence_score
        }

    except Exception as e:
        logger.error(f"Parameter optimization failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/v2/stress-test")
async def run_stress_test(scenarios: List[Dict[str, Any]]):
    """Run stress tests on current parameters"""
    try:
        # Get current parameters (would fetch from database)
        current_params = EconomicParameters(
            ltv=0.7,
            liquidation_threshold=0.8,
            liquidation_penalty=0.05,
            base_rate=0.02,
            slope1=0.08,
            slope2=1.0
        )

        # Sample portfolio
        portfolio = {
            'position_1': {'value_usd': 100000, 'asset': 'ETH', 'leverage': 2.0}
        }

        # Run stress test
        results = economic_model.run_stress_test(
            params=current_params,
            scenarios=scenarios,
            portfolio=portfolio
        )

        return {
            "stress_test_results": results,
            "summary": {
                "scenarios_tested": len(scenarios),
                "worst_case_loss": max(r.get('expected_loss', 0) for r in results.values()),
                "average_survival_probability": np.mean([r.get('survival_probability', 0) for r in results.values()])
            }
        }

    except Exception as e:
        logger.error(f"Stress test failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))

# ============================================================
# MONITORING ENDPOINTS
# ============================================================

@app.get("/api/v2/monitor/status")
async def get_monitoring_status():
    """Get real-time monitoring status"""
    if not monitor:
        return {"status": "not_initialized"}

    metrics = await monitor._calculate_metrics()

    return {
        "status": "active" if monitor.monitoring_active else "inactive",
        "positions_monitored": metrics.positions_monitored,
        "alerts_triggered": metrics.alerts_triggered,
        "alerts_by_level": metrics.alerts_by_level,
        "at_risk_positions": metrics.at_risk_positions,
        "total_value_at_risk": metrics.total_value_at_risk,
        "system_health": metrics.system_health,
        "response_time_ms": metrics.response_time_ms
    }

@app.get("/api/v2/monitor/alerts")
async def get_recent_alerts(limit: int = 50):
    """Get recent monitoring alerts"""
    if not monitor:
        return {"alerts": []}

    recent_alerts = list(monitor.alert_history)[-limit:]

    return {
        "alerts": [
            {
                "id": alert.id,
                "level": alert.level.value,
                "type": alert.alert_type.value,
                "title": alert.title,
                "message": alert.message,
                "timestamp": alert.timestamp.isoformat(),
                "action_required": alert.action_required
            }
            for alert in recent_alerts
        ],
        "total_count": len(monitor.alert_history)
    }

@app.post("/api/v2/monitor/rules")
async def update_monitoring_rules(rules: List[Dict[str, Any]]):
    """Update monitoring rules"""
    # This would update the monitoring rules
    return {
        "status": "rules_updated",
        "rules_count": len(rules)
    }

# ============================================================
# REPORT ENDPOINTS
# ============================================================

@app.post("/api/v2/reports/generate")
async def generate_report(request: ReportRequest):
    """
    Generate comprehensive risk report

    Creates professional reports with visualizations
    """
    try:
        # Gather analysis results
        market_data = await data_infrastructure.collect_comprehensive_data()

        # Simulate some analysis results
        analysis_results = {
            'overall_risk_score': 65.5,
            'alert_level': 'warning',
            'key_risks': [
                {'factor': 'High leverage', 'severity': 'high'},
                {'factor': 'Market volatility', 'severity': 'medium'}
            ],
            'top_recommendations': [
                'Reduce leverage across high-risk positions',
                'Increase monitoring frequency',
                'Consider hedging strategies'
            ],
            'risk_metrics': {
                'var_95': 100000,
                'var_99': 250000,
                'max_drawdown': 0.15,
                'sharpe_ratio': 1.23,
                'liquidation_probability': 0.23
            },
            'market_data': {
                'volatility': 0.04,
                'total_liquidity': 1e9,
                'utilization_rate': 0.65,
                'liquidations_24h': 42
            },
            'timestamp': datetime.utcnow()
        }

        # Generate report
        report = report_generator.generate_comprehensive_report(
            analysis_results,
            report_type=request.report_type
        )

        # Export based on format
        if request.export_format == 'json':
            content = report_generator.export_to_json(report)
        elif request.export_format == 'html':
            content = report_generator.export_to_html(report)
        else:
            content = "PDF export not fully implemented"

        return {
            "report_id": report.report_id,
            "title": report.title,
            "generated_at": report.generated_at.isoformat(),
            "content": content,
            "format": request.export_format
        }

    except Exception as e:
        logger.error(f"Report generation failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/v2/reports/templates")
async def get_report_templates():
    """Get available report templates"""
    return {
        "templates": [
            {
                "type": "full",
                "name": "Comprehensive Risk Report",
                "description": "Complete analysis with all sections"
            },
            {
                "type": "summary",
                "name": "Executive Summary",
                "description": "High-level overview for decision makers"
            },
            {
                "type": "technical",
                "name": "Technical Analysis",
                "description": "Detailed technical metrics and simulations"
            }
        ]
    }

# ============================================================
# WEBSOCKET ENDPOINT
# ============================================================

@app.websocket("/ws/monitor")
async def websocket_monitor(websocket: WebSocket):
    """
    WebSocket endpoint for real-time monitoring

    Streams alerts and metrics to connected clients
    """
    await websocket.accept()
    websocket_connections.add(websocket)

    # Add to monitor's connections
    if monitor:
        await monitor.add_websocket_connection(websocket)

    try:
        # Send initial status
        await websocket.send_json({
            "type": "connection",
            "status": "connected",
            "timestamp": datetime.utcnow().isoformat()
        })

        # Keep connection alive and handle messages
        while True:
            try:
                # Receive message from client
                data = await websocket.receive_json()

                # Handle different message types
                if data.get("type") == "subscribe":
                    # Subscribe to specific alerts
                    await websocket.send_json({
                        "type": "subscription",
                        "status": "subscribed",
                        "topics": data.get("topics", [])
                    })

                elif data.get("type") == "ping":
                    # Respond to ping
                    await websocket.send_json({
                        "type": "pong",
                        "timestamp": datetime.utcnow().isoformat()
                    })

            except WebSocketDisconnect:
                break
            except Exception as e:
                logger.error(f"WebSocket error: {e}")
                break

    finally:
        websocket_connections.discard(websocket)
        if monitor:
            await monitor.remove_websocket_connection(websocket)

# ============================================================
# HEALTH & INFO ENDPOINTS
# ============================================================

@app.get("/")
async def root():
    """API root with version info"""
    return {
        "name": "AI Risk Engine API",
        "version": "2.0.0",
        "status": "operational",
        "features": [
            "risk_assessment",
            "agent_simulation",
            "ml_prediction",
            "parameter_optimization",
            "real_time_monitoring",
            "report_generation"
        ],
        "docs": "/docs",
        "websocket": "/ws/monitor"
    }

@app.get("/api/v2/health")
async def health_check():
    """Comprehensive health check"""
    health_status = {
        "api": "healthy",
        "database": "healthy" if data_infrastructure and data_infrastructure._initialized else "unhealthy",
        "redis": "healthy" if data_infrastructure and data_infrastructure.redis else "unhealthy",
        "simulator": "ready" if simulator else "not_initialized",
        "ml_models": "trained" if ml_predictor and ml_predictor.is_trained else "not_trained",
        "monitor": "active" if monitor and monitor.monitoring_active else "inactive",
        "websockets": len(websocket_connections)
    }

    overall_status = "healthy" if all(
        v in ["healthy", "ready", "trained", "active"] or isinstance(v, int)
        for v in health_status.values()
    ) else "degraded"

    return {
        "status": overall_status,
        "components": health_status,
        "timestamp": datetime.utcnow().isoformat()
    }

# ============================================================
# PHASE 1 BASIC ENDPOINTS (Backward Compatibility)
# ============================================================

# Pydantic models for Phase 1 endpoints
class PositionRequest(BaseModel):
    user_address: str
    protocol: str
    collateral_asset: str
    collateral_amount: float
    debt_asset: Optional[str] = None
    debt_amount: Optional[float] = 0
    health_factor: Optional[float] = None
    ltv: Optional[float] = None

class RiskAssessmentResponse(BaseModel):
    position_id: str
    user_address: str
    protocol: str
    risk_score: float
    liquidation_probability: float
    alert_level: str
    risk_metrics: Dict[str, float]
    recommendations: List[str]
    confidence_level: float
    timestamp: datetime

@app.post("/api/risk/assess", response_model=RiskAssessmentResponse)
async def assess_position_risk(position: PositionRequest):
    """Basic risk assessment endpoint (Phase 1 compatibility)"""
    try:
        market_data = await data_infrastructure.collect_comprehensive_data()

        # Convert position data to include USD values
        # Assuming ETH = $2000, other prices would come from market data in production
        position_data = position.dict()

        # Add calculated USD values if not present
        if 'collateral_value_usd' not in position_data:
            # Simple price estimation (in production, fetch from oracle)
            price_map = {
                'ETH': 2000,
                'BTC': 40000,
                'USDC': 1,
                'USDT': 1,
                'DAI': 1
            }
            collateral_price = price_map.get(position.collateral_asset, 1)
            position_data['collateral_value_usd'] = position.collateral_amount * collateral_price

        if 'debt_value_usd' not in position_data and position.debt_amount:
            debt_price = price_map.get(position.debt_asset, 1) if position.debt_asset else 1
            position_data['debt_value_usd'] = position.debt_amount * debt_price
        else:
            position_data['debt_value_usd'] = 0

        risk_assessment = risk_calculator.calculate_position_risk(
            position_data,
            {
                'volatility': market_data.volatility_metrics if market_data.volatility_metrics else {},
                'market_depth': market_data.market_depth if market_data.market_depth else {},
                'price_change_24h': {}
            }
        )

        return RiskAssessmentResponse(
            position_id=risk_assessment.position_id,
            user_address=risk_assessment.user_address,
            protocol=risk_assessment.protocol,
            risk_score=risk_assessment.risk_metrics.overall_risk_score,
            liquidation_probability=risk_assessment.risk_metrics.liquidation_probability,
            alert_level=risk_assessment.alert_level,
            risk_metrics={
                'market_risk': risk_assessment.risk_metrics.market_risk,
                'liquidity_risk': risk_assessment.risk_metrics.liquidity_risk,
                'value_at_risk': risk_assessment.risk_metrics.value_at_risk,
            },
            recommendations=risk_assessment.recommended_actions,
            confidence_level=risk_assessment.risk_metrics.confidence_level,
            timestamp=datetime.utcnow()
        )
    except Exception as e:
        logger.error(f"Risk assessment failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/user/{user_address}/profile")
async def get_user_risk_profile(user_address: str):
    """Get user risk profile"""
    try:
        risk_score = await data_infrastructure.get_user_risk_score(user_address)
        if risk_score is None:
            raise HTTPException(status_code=404, detail="User not found")

        return {
            "user_address": user_address,
            "risk_score": risk_score,
            "risk_trend": "stable",
            "last_updated": datetime.utcnow().isoformat()
        }
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Failed to get user profile: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/market/data")
async def get_market_data():
    """Get market data"""
    try:
        market_data = await data_infrastructure.collect_comprehensive_data()
        return {
            "timestamp": market_data.timestamp.isoformat(),
            "price_data": market_data.price_data,
            "volatility_metrics": market_data.volatility_metrics,
            "market_depth": market_data.market_depth
        }
    except Exception as e:
        logger.error(f"Failed to get market data: {e}")
        raise HTTPException(status_code=500, detail=str(e))

# ============================================================
# FRONTEND-FRIENDLY ENDPOINTS (/api/ai/*)
# ============================================================

@app.get("/api/ai/health")
async def ai_health_check():
    """Frontend health check"""
    return {"status": "healthy", "service": "AI Risk Engine", "version": "2.0.0"}

@app.post("/api/ai/risk/calculate")
async def calculate_risk_frontend(position: PositionRequest):
    """Frontend-friendly risk calculation"""
    return await assess_position_risk(position)

@app.get("/api/ai/portfolio/{user_address}/risk")
async def get_portfolio_risk_frontend(user_address: str):
    """Frontend-friendly portfolio risk"""
    return await get_user_risk_profile(user_address)

@app.get("/api/ai/market/risk")
async def get_market_risk_frontend():
    """Frontend-friendly market risk"""
    return await get_market_data()

@app.get("/api/ai/monitor/alerts")
async def get_alerts_frontend(limit: int = 50):
    """Frontend-friendly monitoring alerts (alias for /api/v2/monitor/alerts)"""
    return await get_recent_alerts(limit)

@app.post("/api/ai/simulation/run")
async def run_simulation_frontend(request: SimulationRequest):
    """Frontend-friendly simulation"""
    return await run_simulation(request)

# ============================================================
# MAIN
# ============================================================

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=8084,
        log_level="info"
    )