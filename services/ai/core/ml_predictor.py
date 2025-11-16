"""
ML Predictor Module
Hybrid LSTM + XGBoost model for liquidation prediction
Phase 2: Core Implementation
"""

import numpy as np
import pandas as pd
from typing import Dict, List, Tuple, Optional, Any
from dataclasses import dataclass
import logging
from datetime import datetime, timedelta
import pickle
import json

# Machine Learning imports
from sklearn.preprocessing import StandardScaler, MinMaxScaler
from sklearn.model_selection import train_test_split
from sklearn.metrics import (
    accuracy_score, precision_score, recall_score, f1_score,
    roc_auc_score, roc_curve, confusion_matrix, classification_report
)
import xgboost as xgb

# Deep Learning imports
import tensorflow as tf
from tensorflow import keras
from tensorflow.keras import layers, models, callbacks
from tensorflow.keras.optimizers import Adam
from tensorflow.keras.losses import BinaryCrossentropy
from tensorflow.keras.metrics import AUC, Precision, Recall

logger = logging.getLogger(__name__)

# ============================================================
# DATA MODELS
# ============================================================

@dataclass
class PredictionResult:
    """Result from ML prediction"""
    liquidation_probability: float
    confidence_score: float
    time_to_liquidation: Optional[int]  # Hours
    risk_factors: Dict[str, float]
    model_explanations: Dict[str, Any]
    recommendations: List[str]

@dataclass
class ModelPerformance:
    """Model performance metrics"""
    accuracy: float
    precision: float
    recall: float
    f1_score: float
    auc_roc: float
    confusion_matrix: np.ndarray
    feature_importance: Dict[str, float]

# ============================================================
# HYBRID ML PREDICTOR
# ============================================================

class HybridMLPredictor:
    """
    Hybrid machine learning predictor
    Combines XGBoost for structured features and LSTM for time series
    """

    def __init__(self, config: Optional[Dict] = None):
        """
        Initialize predictor

        Args:
            config: Model configuration
        """
        self.config = config or self._default_config()
        self.xgboost_model = None
        self.lstm_model = None
        self.scaler = StandardScaler()
        self.sequence_scaler = MinMaxScaler()
        self.is_trained = False
        self.feature_names = []
        self.ensemble_weights = {
            'xgboost': 0.4,
            'lstm': 0.6
        }

    def _default_config(self) -> Dict:
        """Default model configuration"""
        return {
            'xgboost': {
                'n_estimators': 500,
                'max_depth': 8,
                'learning_rate': 0.01,
                'objective': 'binary:logistic',
                'subsample': 0.8,
                'colsample_bytree': 0.8,
                'scale_pos_weight': 2,  # Handle imbalanced data
                'eval_metric': 'auc',
                'early_stopping_rounds': 50
            },
            'lstm': {
                'sequence_length': 60,
                'n_features': 10,
                'lstm_units': [128, 64, 32],
                'dropout_rate': 0.2,
                'learning_rate': 0.001,
                'batch_size': 32,
                'epochs': 100
            },
            'training': {
                'test_size': 0.2,
                'val_size': 0.1,
                'random_state': 42
            }
        }

    # ============================================================
    # TRAINING
    # ============================================================

    def train(
        self,
        training_data: Dict[str, np.ndarray],
        save_path: Optional[str] = None
    ) -> ModelPerformance:
        """
        Train both models on provided data

        Args:
            training_data: Dictionary with X_train, y_train, etc.
            save_path: Path to save trained models

        Returns:
            ModelPerformance metrics
        """
        try:
            logger.info("Starting hybrid model training...")

            # Extract data
            X_structured = training_data['X_structured']
            X_sequences = training_data.get('X_sequences')
            y = training_data['y']

            # Split data
            splits = self._split_data(X_structured, y)

            # Train XGBoost
            logger.info("Training XGBoost model...")
            xgb_performance = self._train_xgboost(
                splits['X_train'],
                splits['y_train'],
                splits['X_val'],
                splits['y_val']
            )

            # Train LSTM if sequence data provided
            if X_sequences is not None:
                logger.info("Training LSTM model...")
                lstm_performance = self._train_lstm(
                    X_sequences,
                    y,
                    splits['indices']
                )
            else:
                lstm_performance = None

            # Combine performances
            performance = self._combine_performances(xgb_performance, lstm_performance)

            # Save models
            if save_path:
                self._save_models(save_path)

            self.is_trained = True
            logger.info("Model training completed successfully")

            return performance

        except Exception as e:
            logger.error(f"Training failed: {e}")
            raise

    def _train_xgboost(
        self,
        X_train: np.ndarray,
        y_train: np.ndarray,
        X_val: np.ndarray,
        y_val: np.ndarray
    ) -> ModelPerformance:
        """Train XGBoost model"""
        # Scale features
        X_train_scaled = self.scaler.fit_transform(X_train)
        X_val_scaled = self.scaler.transform(X_val)

        # Create DMatrix
        dtrain = xgb.DMatrix(X_train_scaled, label=y_train)
        dval = xgb.DMatrix(X_val_scaled, label=y_val)

        # Training parameters
        params = self.config['xgboost'].copy()
        num_rounds = params.pop('n_estimators', 500)
        early_stopping = params.pop('early_stopping_rounds', 50)

        # Train model
        evals = [(dtrain, 'train'), (dval, 'val')]
        self.xgboost_model = xgb.train(
            params,
            dtrain,
            num_rounds,
            evals,
            early_stopping_rounds=early_stopping,
            verbose_eval=False
        )

        # Evaluate
        y_pred_proba = self.xgboost_model.predict(dval)
        y_pred = (y_pred_proba > 0.5).astype(int)

        # Calculate metrics
        performance = ModelPerformance(
            accuracy=accuracy_score(y_val, y_pred),
            precision=precision_score(y_val, y_pred, zero_division=0),
            recall=recall_score(y_val, y_pred, zero_division=0),
            f1_score=f1_score(y_val, y_pred, zero_division=0),
            auc_roc=roc_auc_score(y_val, y_pred_proba) if len(np.unique(y_val)) > 1 else 0,
            confusion_matrix=confusion_matrix(y_val, y_pred),
            feature_importance=self._get_xgboost_importance()
        )

        logger.info(f"XGBoost Performance - AUC: {performance.auc_roc:.4f}, F1: {performance.f1_score:.4f}")

        return performance

    def _train_lstm(
        self,
        X_sequences: np.ndarray,
        y: np.ndarray,
        indices: Dict[str, np.ndarray]
    ) -> ModelPerformance:
        """Train LSTM model"""
        # Prepare sequences
        train_idx = indices['train']
        val_idx = indices['val']

        X_train = X_sequences[train_idx]
        y_train = y[train_idx]
        X_val = X_sequences[val_idx]
        y_val = y[val_idx]

        # Scale sequences
        n_samples, n_timesteps, n_features = X_train.shape
        X_train_reshaped = X_train.reshape(-1, n_features)
        X_train_scaled = self.sequence_scaler.fit_transform(X_train_reshaped)
        X_train_scaled = X_train_scaled.reshape(n_samples, n_timesteps, n_features)

        X_val_reshaped = X_val.reshape(-1, n_features)
        X_val_scaled = self.sequence_scaler.transform(X_val_reshaped)
        X_val_scaled = X_val_scaled.reshape(X_val.shape[0], n_timesteps, n_features)

        # Build LSTM model
        self.lstm_model = self._build_lstm_model(
            n_timesteps,
            n_features
        )

        # Callbacks
        callbacks_list = [
            callbacks.EarlyStopping(
                monitor='val_loss',
                patience=10,
                restore_best_weights=True
            ),
            callbacks.ReduceLROnPlateau(
                monitor='val_loss',
                factor=0.5,
                patience=5,
                min_lr=1e-6
            )
        ]

        # Train
        history = self.lstm_model.fit(
            X_train_scaled,
            y_train,
            validation_data=(X_val_scaled, y_val),
            epochs=self.config['lstm']['epochs'],
            batch_size=self.config['lstm']['batch_size'],
            callbacks=callbacks_list,
            verbose=0
        )

        # Evaluate
        y_pred_proba = self.lstm_model.predict(X_val_scaled, verbose=0).squeeze()
        y_pred = (y_pred_proba > 0.5).astype(int)

        # Calculate metrics
        performance = ModelPerformance(
            accuracy=accuracy_score(y_val, y_pred),
            precision=precision_score(y_val, y_pred, zero_division=0),
            recall=recall_score(y_val, y_pred, zero_division=0),
            f1_score=f1_score(y_val, y_pred, zero_division=0),
            auc_roc=roc_auc_score(y_val, y_pred_proba) if len(np.unique(y_val)) > 1 else 0,
            confusion_matrix=confusion_matrix(y_val, y_pred),
            feature_importance={}  # LSTM doesn't have simple feature importance
        )

        logger.info(f"LSTM Performance - AUC: {performance.auc_roc:.4f}, F1: {performance.f1_score:.4f}")

        return performance

    def _build_lstm_model(self, sequence_length: int, n_features: int) -> keras.Model:
        """Build LSTM architecture"""
        model = models.Sequential()

        # LSTM layers
        lstm_units = self.config['lstm']['lstm_units']
        dropout_rate = self.config['lstm']['dropout_rate']

        # First LSTM layer
        model.add(layers.LSTM(
            lstm_units[0],
            return_sequences=True,
            input_shape=(sequence_length, n_features)
        ))
        model.add(layers.BatchNormalization())
        model.add(layers.Dropout(dropout_rate))

        # Second LSTM layer
        if len(lstm_units) > 1:
            model.add(layers.LSTM(lstm_units[1], return_sequences=True))
            model.add(layers.BatchNormalization())
            model.add(layers.Dropout(dropout_rate))

        # Third LSTM layer
        if len(lstm_units) > 2:
            model.add(layers.LSTM(lstm_units[2], return_sequences=False))
            model.add(layers.BatchNormalization())
            model.add(layers.Dropout(dropout_rate))

        # Dense layers
        model.add(layers.Dense(32, activation='relu'))
        model.add(layers.Dropout(dropout_rate))
        model.add(layers.Dense(16, activation='relu'))
        model.add(layers.Dense(1, activation='sigmoid'))

        # Compile
        model.compile(
            optimizer=Adam(learning_rate=self.config['lstm']['learning_rate']),
            loss=BinaryCrossentropy(),
            metrics=[
                'accuracy',
                AUC(name='auc'),
                Precision(name='precision'),
                Recall(name='recall')
            ]
        )

        return model

    # ============================================================
    # PREDICTION
    # ============================================================

    def predict(
        self,
        position_data: Dict,
        market_data: Dict,
        historical_data: Optional[np.ndarray] = None
    ) -> PredictionResult:
        """
        Make liquidation prediction for a position

        Args:
            position_data: Current position information
            market_data: Current market conditions
            historical_data: Historical price/volume data

        Returns:
            PredictionResult with probability and explanations
        """
        try:
            if not self.is_trained:
                raise ValueError("Models must be trained before prediction")

            # Extract features
            structured_features = self._extract_structured_features(
                position_data,
                market_data
            )

            # XGBoost prediction
            xgb_prob = self._predict_xgboost(structured_features)

            # LSTM prediction if historical data available
            if historical_data is not None and self.lstm_model is not None:
                sequence_features = self._prepare_sequence(historical_data)
                lstm_prob = self._predict_lstm(sequence_features)
            else:
                lstm_prob = xgb_prob  # Fallback to XGBoost only

            # Ensemble prediction
            final_probability = self._ensemble_prediction(xgb_prob, lstm_prob)

            # Calculate confidence
            confidence = self._calculate_confidence(xgb_prob, lstm_prob)

            # Estimate time to liquidation
            time_to_liquidation = self._estimate_time_to_liquidation(
                final_probability,
                position_data
            )

            # Get risk factors
            risk_factors = self._identify_risk_factors(
                structured_features,
                position_data
            )

            # Generate explanations
            explanations = self._generate_explanations(
                final_probability,
                risk_factors,
                xgb_prob,
                lstm_prob
            )

            # Generate recommendations
            recommendations = self._generate_recommendations(
                final_probability,
                risk_factors,
                position_data
            )

            return PredictionResult(
                liquidation_probability=final_probability,
                confidence_score=confidence,
                time_to_liquidation=time_to_liquidation,
                risk_factors=risk_factors,
                model_explanations=explanations,
                recommendations=recommendations
            )

        except Exception as e:
            logger.error(f"Prediction failed: {e}")
            raise

    def _predict_xgboost(self, features: np.ndarray) -> float:
        """Make prediction with XGBoost"""
        features_scaled = self.scaler.transform(features.reshape(1, -1))
        dmatrix = xgb.DMatrix(features_scaled)
        probability = self.xgboost_model.predict(dmatrix)[0]
        return float(probability)

    def _predict_lstm(self, sequence: np.ndarray) -> float:
        """Make prediction with LSTM"""
        # Scale sequence
        n_timesteps, n_features = sequence.shape
        sequence_reshaped = sequence.reshape(-1, n_features)
        sequence_scaled = self.sequence_scaler.transform(sequence_reshaped)
        sequence_scaled = sequence_scaled.reshape(1, n_timesteps, n_features)

        # Predict
        probability = self.lstm_model.predict(sequence_scaled, verbose=0)[0, 0]
        return float(probability)

    def _ensemble_prediction(self, xgb_prob: float, lstm_prob: float) -> float:
        """Combine predictions from both models"""
        ensemble_prob = (
            self.ensemble_weights['xgboost'] * xgb_prob +
            self.ensemble_weights['lstm'] * lstm_prob
        )
        return ensemble_prob

    def _calculate_confidence(self, xgb_prob: float, lstm_prob: float) -> float:
        """Calculate prediction confidence"""
        # Confidence based on model agreement
        agreement = 1 - abs(xgb_prob - lstm_prob)

        # Confidence based on probability extremity
        extremity = abs(0.5 - ((xgb_prob + lstm_prob) / 2)) * 2

        # Combined confidence
        confidence = (agreement * 0.6 + extremity * 0.4)

        return confidence

    def _estimate_time_to_liquidation(
        self,
        probability: float,
        position_data: Dict
    ) -> Optional[int]:
        """Estimate hours until liquidation"""
        if probability < 0.3:
            return None  # Low risk

        health_factor = position_data.get('health_factor', 2.0)

        # Exponential decay model
        if health_factor < 1.1:
            hours = 1
        elif health_factor < 1.2:
            hours = 6
        elif health_factor < 1.5:
            hours = 24
        elif health_factor < 2.0:
            hours = 72
        else:
            hours = None

        # Adjust based on probability
        if hours and probability > 0.7:
            hours = int(hours * 0.5)  # High probability = sooner

        return hours

    # ============================================================
    # FEATURE EXTRACTION
    # ============================================================

    def _extract_structured_features(
        self,
        position_data: Dict,
        market_data: Dict
    ) -> np.ndarray:
        """Extract features for XGBoost"""
        features = []

        # Position features
        features.append(position_data.get('health_factor', 2.0))
        features.append(position_data.get('ltv', 0.5))
        features.append(position_data.get('leverage', 1.0))
        features.append(np.log1p(position_data.get('collateral_value_usd', 0)))
        features.append(np.log1p(position_data.get('debt_value_usd', 0)))

        # Market features
        features.append(market_data.get('volatility', 0.03))
        features.append(market_data.get('utilization_rate', 0.5))
        features.append(market_data.get('liquidity_ratio', 1.0))

        # Protocol features
        protocol = position_data.get('protocol', 'unknown')
        protocols = ['aave', 'compound', 'gmx']
        for p in protocols:
            features.append(1.0 if protocol == p else 0.0)

        # Time features
        hour_of_day = datetime.now().hour
        day_of_week = datetime.now().weekday()
        features.append(hour_of_day / 24)
        features.append(day_of_week / 7)

        # Risk indicators
        features.append(1.0 if position_data.get('health_factor', 2.0) < 1.5 else 0.0)
        features.append(1.0 if position_data.get('leverage', 1.0) > 3.0 else 0.0)

        return np.array(features)

    def _prepare_sequence(self, historical_data: np.ndarray) -> np.ndarray:
        """Prepare sequence for LSTM"""
        sequence_length = self.config['lstm']['sequence_length']

        # Ensure we have enough data
        if len(historical_data) < sequence_length:
            # Pad with zeros if needed
            padding = np.zeros((sequence_length - len(historical_data), historical_data.shape[1]))
            historical_data = np.vstack([padding, historical_data])

        # Take last sequence_length timesteps
        sequence = historical_data[-sequence_length:]

        return sequence

    # ============================================================
    # EXPLAINABILITY
    # ============================================================

    def _identify_risk_factors(
        self,
        features: np.ndarray,
        position_data: Dict
    ) -> Dict[str, float]:
        """Identify key risk factors"""
        risk_factors = {}

        # Health factor risk
        health_factor = position_data.get('health_factor', 2.0)
        if health_factor < 1.2:
            risk_factors['critical_health_factor'] = 1.0
        elif health_factor < 1.5:
            risk_factors['low_health_factor'] = 0.7
        elif health_factor < 2.0:
            risk_factors['moderate_health_factor'] = 0.4

        # Leverage risk
        leverage = position_data.get('leverage', 1.0)
        if leverage > 5:
            risk_factors['extreme_leverage'] = 1.0
        elif leverage > 3:
            risk_factors['high_leverage'] = 0.6
        elif leverage > 2:
            risk_factors['moderate_leverage'] = 0.3

        # Market risk
        if 'volatility' in position_data:
            volatility = position_data['volatility']
            if volatility > 0.05:
                risk_factors['high_volatility'] = min(volatility * 10, 1.0)

        # Position size risk
        position_size = position_data.get('collateral_value_usd', 0)
        if position_size > 1e6:
            risk_factors['large_position'] = 0.5

        return risk_factors

    def _generate_explanations(
        self,
        probability: float,
        risk_factors: Dict[str, float],
        xgb_prob: float,
        lstm_prob: float
    ) -> Dict[str, Any]:
        """Generate model explanations"""
        explanations = {
            'prediction_breakdown': {
                'xgboost_probability': xgb_prob,
                'lstm_probability': lstm_prob,
                'ensemble_probability': probability,
                'model_weights': self.ensemble_weights
            },
            'risk_factors': risk_factors,
            'feature_importance': self._get_top_features(),
            'confidence_factors': {
                'model_agreement': 1 - abs(xgb_prob - lstm_prob),
                'probability_certainty': abs(0.5 - probability) * 2
            }
        }

        # Add risk level interpretation
        if probability < 0.2:
            explanations['risk_level'] = 'LOW'
            explanations['interpretation'] = 'Position is relatively safe'
        elif probability < 0.5:
            explanations['risk_level'] = 'MEDIUM'
            explanations['interpretation'] = 'Position requires monitoring'
        elif probability < 0.8:
            explanations['risk_level'] = 'HIGH'
            explanations['interpretation'] = 'Position at significant risk'
        else:
            explanations['risk_level'] = 'CRITICAL'
            explanations['interpretation'] = 'Immediate action required'

        return explanations

    def _generate_recommendations(
        self,
        probability: float,
        risk_factors: Dict[str, float],
        position_data: Dict
    ) -> List[str]:
        """Generate actionable recommendations"""
        recommendations = []

        # Probability-based recommendations
        if probability > 0.8:
            recommendations.append("🚨 URGENT: Add collateral immediately or close position")
            recommendations.append("Consider using flashloan to deleverage")
        elif probability > 0.5:
            recommendations.append("⚠️ Add collateral to improve health factor above 1.5")
            recommendations.append("Reduce position size by 30-50%")
        elif probability > 0.3:
            recommendations.append("Monitor position closely for next 24 hours")
            recommendations.append("Set up alerts for price movements > 5%")

        # Risk factor specific recommendations
        if 'critical_health_factor' in risk_factors:
            recommendations.append("Health factor critical - liquidation imminent")

        if 'extreme_leverage' in risk_factors:
            recommendations.append("Leverage too high - reduce to below 3x")

        if 'high_volatility' in risk_factors:
            recommendations.append("Market volatility elevated - consider hedging")

        if 'large_position' in risk_factors:
            recommendations.append("Large position size - consider splitting across protocols")

        # Protocol specific recommendations
        protocol = position_data.get('protocol', '')
        if protocol == 'gmx':
            recommendations.append("GMX positions have no grace period - act immediately")
        elif protocol == 'aave':
            recommendations.append("Use Aave's flash loans for efficient deleveraging")

        return recommendations[:5]  # Limit to top 5 recommendations

    # ============================================================
    # FEATURE IMPORTANCE & INTERPRETABILITY
    # ============================================================

    def _get_xgboost_importance(self) -> Dict[str, float]:
        """Get XGBoost feature importance"""
        if not self.xgboost_model:
            return {}

        importance = self.xgboost_model.get_score(importance_type='gain')

        # Normalize
        total = sum(importance.values())
        if total > 0:
            importance = {k: v/total for k, v in importance.items()}

        return importance

    def _get_top_features(self, n: int = 10) -> List[Tuple[str, float]]:
        """Get top n most important features"""
        importance = self._get_xgboost_importance()

        # Sort by importance
        sorted_features = sorted(
            importance.items(),
            key=lambda x: x[1],
            reverse=True
        )

        return sorted_features[:n]

    # ============================================================
    # MODEL PERSISTENCE
    # ============================================================

    def _save_models(self, path: str):
        """Save trained models"""
        # Save XGBoost
        if self.xgboost_model:
            self.xgboost_model.save_model(f"{path}/xgboost_model.bin")

        # Save LSTM
        if self.lstm_model:
            self.lstm_model.save(f"{path}/lstm_model.h5")

        # Save scalers
        with open(f"{path}/scalers.pkl", 'wb') as f:
            pickle.dump({
                'structured_scaler': self.scaler,
                'sequence_scaler': self.sequence_scaler
            }, f)

        # Save config
        with open(f"{path}/config.json", 'w') as f:
            json.dump(self.config, f)

        logger.info(f"Models saved to {path}")

    def load_models(self, path: str):
        """Load trained models"""
        # Load XGBoost
        self.xgboost_model = xgb.Booster()
        self.xgboost_model.load_model(f"{path}/xgboost_model.bin")

        # Load LSTM
        self.lstm_model = keras.models.load_model(f"{path}/lstm_model.h5")

        # Load scalers
        with open(f"{path}/scalers.pkl", 'rb') as f:
            scalers = pickle.load(f)
            self.scaler = scalers['structured_scaler']
            self.sequence_scaler = scalers['sequence_scaler']

        # Load config
        with open(f"{path}/config.json", 'r') as f:
            self.config = json.load(f)

        self.is_trained = True
        logger.info(f"Models loaded from {path}")

    # ============================================================
    # UTILITIES
    # ============================================================

    def _split_data(
        self,
        X: np.ndarray,
        y: np.ndarray
    ) -> Dict[str, np.ndarray]:
        """Split data into train/val/test sets"""
        test_size = self.config['training']['test_size']
        val_size = self.config['training']['val_size']
        random_state = self.config['training']['random_state']

        # First split: train+val vs test
        X_temp, X_test, y_temp, y_test = train_test_split(
            X, y,
            test_size=test_size,
            random_state=random_state,
            stratify=y
        )

        # Second split: train vs val
        val_size_adjusted = val_size / (1 - test_size)
        X_train, X_val, y_train, y_val = train_test_split(
            X_temp, y_temp,
            test_size=val_size_adjusted,
            random_state=random_state,
            stratify=y_temp
        )

        # Create indices for sequence data
        n_total = len(X)
        indices = np.arange(n_total)
        train_idx = indices[:len(X_train)]
        val_idx = indices[len(X_train):len(X_train)+len(X_val)]
        test_idx = indices[len(X_train)+len(X_val):]

        return {
            'X_train': X_train,
            'X_val': X_val,
            'X_test': X_test,
            'y_train': y_train,
            'y_val': y_val,
            'y_test': y_test,
            'indices': {
                'train': train_idx,
                'val': val_idx,
                'test': test_idx
            }
        }

    def _combine_performances(
        self,
        xgb_perf: ModelPerformance,
        lstm_perf: Optional[ModelPerformance]
    ) -> ModelPerformance:
        """Combine performance metrics from both models"""
        if lstm_perf is None:
            return xgb_perf

        # Weighted average of metrics
        weights = self.ensemble_weights

        combined = ModelPerformance(
            accuracy=(
                xgb_perf.accuracy * weights['xgboost'] +
                lstm_perf.accuracy * weights['lstm']
            ),
            precision=(
                xgb_perf.precision * weights['xgboost'] +
                lstm_perf.precision * weights['lstm']
            ),
            recall=(
                xgb_perf.recall * weights['xgboost'] +
                lstm_perf.recall * weights['lstm']
            ),
            f1_score=(
                xgb_perf.f1_score * weights['xgboost'] +
                lstm_perf.f1_score * weights['lstm']
            ),
            auc_roc=(
                xgb_perf.auc_roc * weights['xgboost'] +
                lstm_perf.auc_roc * weights['lstm']
            ),
            confusion_matrix=xgb_perf.confusion_matrix,  # Use XGBoost as primary
            feature_importance=xgb_perf.feature_importance
        )

        return combined