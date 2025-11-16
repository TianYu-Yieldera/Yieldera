"""
AI Risk Engine API Service
FastAPI implementation for AI risk assessment endpoints
Phase 1: Core API Implementation
"""

from fastapi import FastAPI, HTTPException, BackgroundTasks, Depends
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from pydantic import BaseModel, Field
from typing import Dict, List, Optional, Any
from datetime import datetime
import asyncio
import logging
import os
from contextlib import asynccontextmanager

# Import core modules
import sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from core.data_infrastructure import DataInfrastructure
from core.risk_calculator import RiskCalculator
from core.ml_data_pipeline import MLDataPipeline

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# ============================================================
# PYDANTIC MODELS
# ============================================================

class HealthCheck(BaseModel):
    status: str
    timestamp: datetime
    version: str
    services: Dict[str, str]

class PositionRequest(BaseModel):
    user_address: str = Field(..., description="User wallet address")
    protocol: str = Field(..., description="Protocol name (aave, compound, gmx)")
    collateral_asset: str = Field(..., description="Collateral asset address")
    collateral_amount: float = Field(..., description="Collateral amount")
    debt_asset: Optional[str] = Field(None, description="Debt asset address")
    debt_amount: Optional[float] = Field(0, description="Debt amount")
    health_factor: Optional[float] = Field(None, description="Current health factor")
    ltv: Optional[float] = Field(None, description="Loan-to-value ratio")
    collateral_value_usd: Optional[float] = Field(None, description="Collateral value in USD")
    debt_value_usd: Optional[float] = Field(None, description="Debt value in USD")

class RiskAssessmentResponse(BaseModel):
    position_id: str
    user_address: str
    protocol: str
    risk_score: float = Field(..., description="Overall risk score 0-100")
    liquidation_probability: float = Field(..., description="Liquidation probability 0-1")
    alert_level: str = Field(..., description="Alert level: safe, warning, danger, critical")
    risk_metrics: Dict[str, float]
    recommendations: List[str]
    confidence_level: float = Field(..., description="Model confidence 0-1")
    timestamp: datetime

class MarketDataResponse(BaseModel):
    timestamp: datetime
    price_data: Dict[str, Any]
    volatility_metrics: Dict[str, float]
    market_depth: Dict[str, Any]
    liquidation_stats: Dict[str, Any]
    user_metrics: Dict[str, Any]

class UserRiskProfileResponse(BaseModel):
    user_address: str
    risk_score: float
    risk_trend: str  # improving, worsening, stable
    liquidation_count: int
    total_positions: int
    avg_health_factor: float
    avg_leverage: float
    risk_ranking: Optional[int]
    risk_percentile: Optional[float]
    last_updated: datetime

class SimulationRequest(BaseModel):
    scenario: str = Field(..., description="Simulation scenario type")
    parameters: Dict[str, Any] = Field(..., description="Scenario parameters")
    time_steps: Optional[int] = Field(100, description="Number of simulation steps")
    num_agents: Optional[int] = Field(1000, description="Number of agents in simulation")

class SimulationResponse(BaseModel):
    scenario: str
    results: Dict[str, Any]
    risk_metrics: Dict[str, float]
    recommendations: List[str]
    confidence_interval: Dict[str, float]
    execution_time: float

# ============================================================
# GLOBAL INSTANCES
# ============================================================

data_infrastructure: Optional[DataInfrastructure] = None
risk_calculator: Optional[RiskCalculator] = None
ml_pipeline: Optional[MLDataPipeline] = None

# ============================================================
# LIFESPAN MANAGEMENT
# ============================================================

@asynccontextmanager
async def lifespan(app: FastAPI):
    """Manage application lifecycle"""
    global data_infrastructure, risk_calculator, ml_pipeline

    # Startup
    logger.info("Starting AI Risk Engine API...")

    # Initialize components
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

    logger.info("AI Risk Engine API started successfully")

    yield

    # Shutdown
    logger.info("Shutting down AI Risk Engine API...")
    if data_infrastructure:
        await data_infrastructure.close()
    logger.info("AI Risk Engine API shut down complete")

# ============================================================
# FASTAPI APP
# ============================================================

app = FastAPI(
    title="AI Risk Engine API",
    description="Gauntlet-standard AI risk assessment for DeFi positions",
    version="1.0.0",
    lifespan=lifespan
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure appropriately for production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include frontend adapter routes (provides frontend-friendly API paths)
try:
    from api.frontend_adapter import router as frontend_router
    app.include_router(frontend_router)
    logger.info("✅ Frontend adapter routes registered")
except ImportError as e:
    logger.warning(f"Frontend adapter not loaded: {e}")

# ============================================================
# HEALTH & STATUS ENDPOINTS
# ============================================================

@app.get("/", response_model=HealthCheck)
async def health_check():
    """Health check endpoint"""
    return HealthCheck(
        status="healthy",
        timestamp=datetime.utcnow(),
        version="1.0.0",
        services={
            "database": "connected" if data_infrastructure and data_infrastructure._initialized else "disconnected",
            "redis": "connected" if data_infrastructure and data_infrastructure.redis else "disconnected",
            "risk_calculator": "ready" if risk_calculator else "not initialized",
            "ml_pipeline": "ready" if ml_pipeline else "not initialized"
        }
    )

@app.get("/api/status")
async def api_status():
    """Detailed API status"""
    return {
        "status": "operational",
        "uptime": "N/A",  # Would track actual uptime
        "requests_processed": 0,  # Would track actual requests
        "average_response_time": 0,  # Would track actual response times
        "last_data_update": datetime.utcnow().isoformat(),
        "models_loaded": True
    }

# ============================================================
# RISK ASSESSMENT ENDPOINTS
# ============================================================

@app.post("/api/risk/assess", response_model=RiskAssessmentResponse)
async def assess_position_risk(position: PositionRequest):
    """
    Assess risk for a specific position

    Returns comprehensive risk metrics and recommendations
    """
    try:
        # Get current market data
        market_data = await data_infrastructure.collect_comprehensive_data()

        # Convert request to dict
        position_data = position.dict()

        # Calculate risk
        risk_assessment = risk_calculator.calculate_position_risk(
            position_data,
            {
                'volatility': market_data.volatility_metrics,
                'market_depth': market_data.market_depth,
                'price_change_24h': {}  # Would calculate from price_data
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
                'credit_risk': risk_assessment.risk_metrics.credit_risk,
                'operational_risk': risk_assessment.risk_metrics.operational_risk,
                'value_at_risk': risk_assessment.risk_metrics.value_at_risk,
                'expected_shortfall': risk_assessment.risk_metrics.expected_shortfall,
                'max_drawdown': risk_assessment.risk_metrics.max_drawdown,
                'sharpe_ratio': risk_assessment.risk_metrics.sharpe_ratio,
                'health_factor': risk_assessment.risk_metrics.health_factor,
                'leverage': risk_assessment.risk_metrics.leverage
            },
            recommendations=risk_assessment.recommended_actions,
            confidence_level=risk_assessment.risk_metrics.confidence_level,
            timestamp=datetime.utcnow()
        )

    except Exception as e:
        logger.error(f"Failed to assess position risk: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/risk/portfolio")
async def assess_portfolio_risk(positions: List[PositionRequest]):
    """
    Assess risk for entire portfolio

    Returns portfolio-level risk metrics
    """
    try:
        # Get current market data
        market_data = await data_infrastructure.collect_comprehensive_data()

        # Convert requests to dicts
        positions_data = [p.dict() for p in positions]

        # Calculate portfolio risk
        portfolio_risk = risk_calculator.calculate_portfolio_risk(
            positions_data,
            {
                'volatility': market_data.volatility_metrics,
                'market_depth': market_data.market_depth
            },
            market_data.correlation_matrix
        )

        return {
            'overall_risk_score': portfolio_risk.overall_risk_score,
            'risk_metrics': {
                'liquidation_probability': portfolio_risk.liquidation_probability,
                'market_risk': portfolio_risk.market_risk,
                'liquidity_risk': portfolio_risk.liquidity_risk,
                'credit_risk': portfolio_risk.credit_risk,
                'operational_risk': portfolio_risk.operational_risk,
                'value_at_risk': portfolio_risk.value_at_risk,
                'expected_shortfall': portfolio_risk.expected_shortfall,
                'max_drawdown': portfolio_risk.max_drawdown,
                'sharpe_ratio': portfolio_risk.sharpe_ratio,
                'avg_health_factor': portfolio_risk.health_factor,
                'portfolio_leverage': portfolio_risk.leverage
            },
            'confidence_level': portfolio_risk.confidence_level,
            'timestamp': datetime.utcnow()
        }

    except Exception as e:
        logger.error(f"Failed to assess portfolio risk: {e}")
        raise HTTPException(status_code=500, detail=str(e))

# ============================================================
# MARKET DATA ENDPOINTS
# ============================================================

@app.get("/api/market/data", response_model=MarketDataResponse)
async def get_market_data():
    """
    Get current market data and analytics

    Returns comprehensive market snapshot
    """
    try:
        # Check cache first
        cached = await data_infrastructure.get_cached_snapshot()
        if cached:
            market_data = cached
        else:
            market_data = await data_infrastructure.collect_comprehensive_data()

        # Get liquidation stats
        liquidation_stats = await data_infrastructure.collect_liquidation_history(days=1)

        return MarketDataResponse(
            timestamp=market_data.timestamp,
            price_data=market_data.price_data,
            volatility_metrics=market_data.volatility_metrics,
            market_depth=market_data.market_depth,
            liquidation_stats={
                'count_24h': len(liquidation_stats),
                'recent_liquidations': liquidation_stats[:5]
            },
            user_metrics=market_data.user_metrics
        )

    except Exception as e:
        logger.error(f"Failed to get market data: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/market/volatility/{asset}")
async def get_asset_volatility(asset: str, days: int = 30):
    """Get historical volatility for an asset"""
    try:
        volatility = await data_infrastructure.calculate_volatility_metrics()

        return {
            'asset': asset,
            'volatility': volatility.get(asset, 0),
            'period_days': days,
            'timestamp': datetime.utcnow()
        }

    except Exception as e:
        logger.error(f"Failed to get volatility: {e}")
        raise HTTPException(status_code=500, detail=str(e))

# ============================================================
# USER ENDPOINTS
# ============================================================

@app.get("/api/user/{user_address}/profile", response_model=UserRiskProfileResponse)
async def get_user_risk_profile(user_address: str):
    """
    Get user risk profile

    Returns user's risk metrics and history
    """
    try:
        # Get user risk score
        risk_score = await data_infrastructure.get_user_risk_score(user_address)

        if risk_score is None:
            raise HTTPException(status_code=404, detail="User not found")

        # Get user metrics from database
        async with data_infrastructure.db_pool.acquire() as conn:
            row = await conn.fetchrow(
                """
                SELECT * FROM user_risk_profiles
                WHERE user_address = $1
                """,
                user_address
            )

            if not row:
                raise HTTPException(status_code=404, detail="User not found")

            # Get user ranking
            ranking = await conn.fetchrow(
                """
                SELECT risk_rank, risk_percentile
                FROM user_risk_rankings
                WHERE user_address = $1
                """,
                user_address
            )

        # Determine trend (simplified)
        trend = "stable"  # Would calculate from historical data

        return UserRiskProfileResponse(
            user_address=user_address,
            risk_score=float(row['risk_score']),
            risk_trend=trend,
            liquidation_count=row['liquidation_count'],
            total_positions=row['total_positions'],
            avg_health_factor=float(row['avg_health_factor']) if row['avg_health_factor'] else 2.0,
            avg_leverage=float(row['avg_leverage']) if row['avg_leverage'] else 1.0,
            risk_ranking=ranking['risk_rank'] if ranking else None,
            risk_percentile=float(ranking['risk_percentile']) if ranking else None,
            last_updated=row['last_updated']
        )

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Failed to get user profile: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/user/{user_address}/positions")
async def get_user_positions(user_address: str):
    """Get all positions for a user"""
    try:
        async with data_infrastructure.db_pool.acquire() as conn:
            rows = await conn.fetch(
                """
                SELECT * FROM position_snapshots
                WHERE user_address = $1
                AND time > NOW() - INTERVAL '24 hours'
                ORDER BY time DESC
                LIMIT 10
                """,
                user_address
            )

            positions = []
            for row in rows:
                positions.append({
                    'protocol': row['protocol'],
                    'collateral_asset': row['collateral_asset'],
                    'collateral_amount': float(row['collateral_amount']),
                    'debt_asset': row['debt_asset'],
                    'debt_amount': float(row['debt_amount']) if row['debt_amount'] else 0,
                    'health_factor': float(row['health_factor']) if row['health_factor'] else None,
                    'ltv': float(row['ltv']) if row['ltv'] else None,
                    'leverage': float(row['leverage']) if row['leverage'] else None,
                    'timestamp': row['time'].isoformat()
                })

            return {
                'user_address': user_address,
                'positions': positions,
                'count': len(positions)
            }

    except Exception as e:
        logger.error(f"Failed to get user positions: {e}")
        raise HTTPException(status_code=500, detail=str(e))

# ============================================================
# MONITORING ENDPOINTS
# ============================================================

@app.get("/api/monitor/high-risk")
async def get_high_risk_positions(limit: int = 50):
    """Get current high-risk positions across all protocols"""
    try:
        positions = await data_infrastructure.get_high_risk_positions(limit)

        return {
            'count': len(positions),
            'positions': positions,
            'timestamp': datetime.utcnow()
        }

    except Exception as e:
        logger.error(f"Failed to get high-risk positions: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/monitor/liquidations/recent")
async def get_recent_liquidations(hours: int = 24):
    """Get recent liquidation events"""
    try:
        liquidations = await data_infrastructure.collect_liquidation_history(days=hours/24)

        return {
            'count': len(liquidations),
            'liquidations': liquidations[:100],  # Limit response size
            'period_hours': hours,
            'timestamp': datetime.utcnow()
        }

    except Exception as e:
        logger.error(f"Failed to get recent liquidations: {e}")
        raise HTTPException(status_code=500, detail=str(e))

# ============================================================
# SIMULATION ENDPOINTS (Placeholder)
# ============================================================

@app.post("/api/simulate", response_model=SimulationResponse)
async def run_simulation(request: SimulationRequest):
    """
    Run risk simulation (placeholder for Phase 2)

    This endpoint will run agent-based simulations in Phase 2
    """
    return SimulationResponse(
        scenario=request.scenario,
        results={
            "message": "Simulation endpoint will be implemented in Phase 2",
            "phase": "1",
            "status": "not_implemented"
        },
        risk_metrics={},
        recommendations=["Full simulation capabilities coming in Phase 2"],
        confidence_interval={},
        execution_time=0.0
    )

# ============================================================
# ERROR HANDLERS
# ============================================================

@app.exception_handler(Exception)
async def general_exception_handler(request, exc):
    logger.error(f"Unhandled exception: {exc}")
    return JSONResponse(
        status_code=500,
        content={
            "error": "Internal server error",
            "detail": str(exc) if os.getenv("DEBUG") == "true" else "An error occurred",
            "timestamp": datetime.utcnow().isoformat()
        }
    )

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