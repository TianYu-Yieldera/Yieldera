"""
ML Data Pipeline Module
Prepares and processes data for machine learning models
Phase 1: Core Implementation
"""

import numpy as np
import pandas as pd
from typing import Dict, List, Tuple, Optional, Any
from datetime import datetime, timedelta
from sklearn.preprocessing import StandardScaler, MinMaxScaler
from sklearn.feature_selection import SelectKBest, f_classif
import logging

logger = logging.getLogger(__name__)

# ============================================================
# FEATURE ENGINEERING
# ============================================================

class FeatureEngineer:
    """Extract and engineer features for ML models"""

    def __init__(self):
        self.scaler = StandardScaler()
        self.feature_names = []

    def extract_position_features(self, position: Dict) -> np.ndarray:
        """
        Extract features from a position

        Args:
            position: Position data

        Returns:
            Feature vector
        """
        features = []
        names = []

        # Basic position features
        features.append(position.get('health_factor', 2.0))
        names.append('health_factor')

        features.append(position.get('ltv', 0.5))
        names.append('ltv')

        features.append(position.get('leverage', 1.0))
        names.append('leverage')

        # Value features
        collateral_value = position.get('collateral_value_usd', 0)
        debt_value = position.get('debt_value_usd', 0)

        features.append(np.log1p(collateral_value))  # Log transform for scale
        names.append('log_collateral_value')

        features.append(np.log1p(debt_value))
        names.append('log_debt_value')

        # Ratio features
        if collateral_value > 0:
            features.append(debt_value / collateral_value)
        else:
            features.append(0)
        names.append('debt_to_collateral_ratio')

        # Protocol encoding (one-hot)
        protocol = position.get('protocol', 'unknown')
        protocols = ['aave', 'compound', 'gmx']
        for p in protocols:
            features.append(1.0 if protocol == p else 0.0)
            names.append(f'protocol_{p}')

        # Time-based features
        position_age = position.get('position_age_days', 0)
        features.append(position_age)
        names.append('position_age_days')

        features.append(1.0 if position_age < 1 else 0.0)
        names.append('is_new_position')

        # Store feature names
        if not self.feature_names:
            self.feature_names = names

        return np.array(features)

    def extract_market_features(self, market_data: Dict) -> np.ndarray:
        """
        Extract market-level features

        Args:
            market_data: Market data

        Returns:
            Feature vector
        """
        features = []

        # Volatility features
        volatility = market_data.get('volatility_metrics', {})
        for asset in ['ETH', 'BTC', 'USDC']:
            features.append(volatility.get(asset, 0.1))

        # Market depth features
        market_depth = market_data.get('market_depth', {})
        avg_utilization = np.mean([
            m.get('utilization_rate', 0.5)
            for m in market_depth.values()
        ])
        features.append(avg_utilization)

        # User behavior features
        user_metrics = market_data.get('user_metrics', {})
        features.append(user_metrics.get('avg_risk_score', 50) / 100)
        features.append(user_metrics.get('avg_leverage', 1.0))

        return np.array(features)

    def extract_user_features(self, user_data: Dict) -> np.ndarray:
        """
        Extract user-specific features

        Args:
            user_data: User profile data

        Returns:
            Feature vector
        """
        features = []

        # Risk profile features
        features.append(user_data.get('risk_score', 50) / 100)
        features.append(user_data.get('avg_leverage', 1.0))
        features.append(user_data.get('liquidation_count', 0))
        features.append(np.log1p(user_data.get('total_positions', 0)))
        features.append(user_data.get('avg_health_factor', 2.0))
        features.append(np.log1p(user_data.get('total_volume_usd', 0)))

        # Activity features
        last_activity = user_data.get('last_activity')
        if last_activity:
            days_since_activity = (datetime.utcnow() - last_activity).days
            features.append(days_since_activity)
        else:
            features.append(30)  # Default

        return np.array(features)

    def create_time_series_features(
        self,
        price_history: np.ndarray,
        window_sizes: List[int] = [5, 10, 20, 50]
    ) -> Dict[str, float]:
        """
        Create time series features from price history

        Args:
            price_history: Historical prices
            window_sizes: Moving average window sizes

        Returns:
            Dictionary of time series features
        """
        features = {}

        if len(price_history) < 2:
            return features

        # Returns
        returns = np.diff(price_history) / price_history[:-1]
        features['return_mean'] = np.mean(returns)
        features['return_std'] = np.std(returns)
        features['return_skew'] = self._calculate_skew(returns)
        features['return_kurtosis'] = self._calculate_kurtosis(returns)

        # Moving averages
        for window in window_sizes:
            if len(price_history) >= window:
                ma = np.mean(price_history[-window:])
                features[f'ma_{window}'] = ma
                features[f'price_to_ma_{window}'] = price_history[-1] / ma if ma > 0 else 1

        # Technical indicators
        if len(price_history) >= 14:
            features['rsi'] = self._calculate_rsi(price_history, 14)

        if len(price_history) >= 26:
            features['macd'] = self._calculate_macd(price_history)

        # Volatility
        features['volatility_7d'] = np.std(returns[-7:]) if len(returns) >= 7 else np.std(returns)
        features['volatility_30d'] = np.std(returns[-30:]) if len(returns) >= 30 else np.std(returns)

        return features

    def _calculate_skew(self, returns: np.ndarray) -> float:
        """Calculate skewness of returns"""
        if len(returns) < 3:
            return 0
        mean = np.mean(returns)
        std = np.std(returns)
        if std == 0:
            return 0
        return np.mean(((returns - mean) / std) ** 3)

    def _calculate_kurtosis(self, returns: np.ndarray) -> float:
        """Calculate kurtosis of returns"""
        if len(returns) < 4:
            return 0
        mean = np.mean(returns)
        std = np.std(returns)
        if std == 0:
            return 0
        return np.mean(((returns - mean) / std) ** 4) - 3

    def _calculate_rsi(self, prices: np.ndarray, period: int = 14) -> float:
        """Calculate RSI (Relative Strength Index)"""
        deltas = np.diff(prices)
        gains = deltas.copy()
        losses = deltas.copy()
        gains[gains < 0] = 0
        losses[losses > 0] = 0
        losses = -losses

        avg_gain = np.mean(gains[-period:]) if len(gains) >= period else np.mean(gains)
        avg_loss = np.mean(losses[-period:]) if len(losses) >= period else np.mean(losses)

        if avg_loss == 0:
            return 100
        rs = avg_gain / avg_loss
        rsi = 100 - (100 / (1 + rs))
        return rsi

    def _calculate_macd(self, prices: np.ndarray) -> float:
        """Calculate MACD (Moving Average Convergence Divergence)"""
        exp1 = pd.Series(prices).ewm(span=12).mean().iloc[-1]
        exp2 = pd.Series(prices).ewm(span=26).mean().iloc[-1]
        return exp1 - exp2

# ============================================================
# ML DATA PIPELINE CLASS
# ============================================================

class MLDataPipeline:
    """
    End-to-end data pipeline for ML models
    Handles data preparation, feature engineering, and preprocessing
    """

    def __init__(self):
        self.feature_engineer = FeatureEngineer()
        self.scaler = StandardScaler()
        self.is_fitted = False

    def prepare_training_data(
        self,
        positions: List[Dict],
        market_data: Dict,
        labels: Optional[List[int]] = None
    ) -> Tuple[np.ndarray, Optional[np.ndarray]]:
        """
        Prepare training data for ML models

        Args:
            positions: List of historical positions
            market_data: Market conditions
            labels: Target labels (e.g., liquidated or not)

        Returns:
            Tuple of (features, labels)
        """
        try:
            feature_list = []

            for position in positions:
                # Extract position features
                pos_features = self.feature_engineer.extract_position_features(position)

                # Extract market features
                market_features = self.feature_engineer.extract_market_features(market_data)

                # Combine features
                combined = np.concatenate([pos_features, market_features])
                feature_list.append(combined)

            # Convert to numpy array
            X = np.array(feature_list)

            # Fit and transform features
            if not self.is_fitted:
                X = self.scaler.fit_transform(X)
                self.is_fitted = True
            else:
                X = self.scaler.transform(X)

            # Convert labels if provided
            y = np.array(labels) if labels else None

            return X, y

        except Exception as e:
            logger.error(f"Failed to prepare training data: {e}")
            raise

    def prepare_prediction_data(
        self,
        position: Dict,
        market_data: Dict,
        user_data: Optional[Dict] = None,
        price_history: Optional[np.ndarray] = None
    ) -> np.ndarray:
        """
        Prepare data for prediction

        Args:
            position: Current position
            market_data: Current market conditions
            user_data: User profile data
            price_history: Historical prices

        Returns:
            Feature vector ready for prediction
        """
        try:
            features = []

            # Position features
            pos_features = self.feature_engineer.extract_position_features(position)
            features.extend(pos_features)

            # Market features
            market_features = self.feature_engineer.extract_market_features(market_data)
            features.extend(market_features)

            # User features
            if user_data:
                user_features = self.feature_engineer.extract_user_features(user_data)
                features.extend(user_features)

            # Time series features
            if price_history is not None and len(price_history) > 0:
                ts_features = self.feature_engineer.create_time_series_features(price_history)
                features.extend(ts_features.values())

            # Convert to numpy array and reshape
            X = np.array(features).reshape(1, -1)

            # Scale features
            if self.is_fitted:
                X = self.scaler.transform(X)

            return X

        except Exception as e:
            logger.error(f"Failed to prepare prediction data: {e}")
            raise

    def create_sequences(
        self,
        data: np.ndarray,
        sequence_length: int = 60,
        target_column: int = -1
    ) -> Tuple[np.ndarray, np.ndarray]:
        """
        Create sequences for LSTM models

        Args:
            data: Time series data
            sequence_length: Length of each sequence
            target_column: Column to predict

        Returns:
            Tuple of (sequences, targets)
        """
        sequences = []
        targets = []

        for i in range(len(data) - sequence_length):
            seq = data[i:i + sequence_length]
            target = data[i + sequence_length, target_column]
            sequences.append(seq)
            targets.append(target)

        return np.array(sequences), np.array(targets)

    def augment_data(
        self,
        X: np.ndarray,
        y: np.ndarray,
        augmentation_factor: float = 0.1
    ) -> Tuple[np.ndarray, np.ndarray]:
        """
        Augment training data with synthetic examples

        Args:
            X: Feature matrix
            y: Labels
            augmentation_factor: Fraction of synthetic data to add

        Returns:
            Augmented data
        """
        n_synthetic = int(len(X) * augmentation_factor)
        synthetic_X = []
        synthetic_y = []

        for _ in range(n_synthetic):
            # Select random sample
            idx = np.random.randint(0, len(X))
            sample = X[idx].copy()
            label = y[idx]

            # Add noise
            noise = np.random.normal(0, 0.05, sample.shape)
            sample += noise

            synthetic_X.append(sample)
            synthetic_y.append(label)

        # Combine original and synthetic
        X_augmented = np.vstack([X, np.array(synthetic_X)])
        y_augmented = np.hstack([y, np.array(synthetic_y)])

        return X_augmented, y_augmented

    def balance_dataset(
        self,
        X: np.ndarray,
        y: np.ndarray
    ) -> Tuple[np.ndarray, np.ndarray]:
        """
        Balance dataset by oversampling minority class

        Args:
            X: Feature matrix
            y: Labels

        Returns:
            Balanced dataset
        """
        from sklearn.utils import resample

        # Separate classes
        class_0_idx = np.where(y == 0)[0]
        class_1_idx = np.where(y == 1)[0]

        # Determine minority and majority
        if len(class_0_idx) < len(class_1_idx):
            minority_idx = class_0_idx
            majority_idx = class_1_idx
        else:
            minority_idx = class_1_idx
            majority_idx = class_0_idx

        # Oversample minority
        minority_oversampled = resample(
            minority_idx,
            n_samples=len(majority_idx),
            random_state=42
        )

        # Combine
        balanced_idx = np.hstack([majority_idx, minority_oversampled])
        np.random.shuffle(balanced_idx)

        return X[balanced_idx], y[balanced_idx]

    def select_features(
        self,
        X: np.ndarray,
        y: np.ndarray,
        k: int = 20
    ) -> Tuple[np.ndarray, List[int]]:
        """
        Select top k features

        Args:
            X: Feature matrix
            y: Labels
            k: Number of features to select

        Returns:
            Tuple of (selected features, indices)
        """
        selector = SelectKBest(f_classif, k=k)
        X_selected = selector.fit_transform(X, y)
        selected_indices = selector.get_support(indices=True)

        return X_selected, selected_indices.tolist()

    def split_data(
        self,
        X: np.ndarray,
        y: np.ndarray,
        test_size: float = 0.2,
        val_size: float = 0.1
    ) -> Dict[str, np.ndarray]:
        """
        Split data into train, validation, and test sets

        Args:
            X: Feature matrix
            y: Labels
            test_size: Test set size
            val_size: Validation set size

        Returns:
            Dictionary with train, val, test splits
        """
        from sklearn.model_selection import train_test_split

        # Split off test set
        X_temp, X_test, y_temp, y_test = train_test_split(
            X, y, test_size=test_size, random_state=42, stratify=y
        )

        # Split remaining into train and validation
        val_size_adjusted = val_size / (1 - test_size)
        X_train, X_val, y_train, y_val = train_test_split(
            X_temp, y_temp, test_size=val_size_adjusted, random_state=42, stratify=y_temp
        )

        return {
            'X_train': X_train,
            'y_train': y_train,
            'X_val': X_val,
            'y_val': y_val,
            'X_test': X_test,
            'y_test': y_test
        }

    def create_labels(
        self,
        positions: List[Dict],
        label_type: str = 'liquidation'
    ) -> np.ndarray:
        """
        Create labels for supervised learning

        Args:
            positions: List of positions
            label_type: Type of label to create

        Returns:
            Label array
        """
        labels = []

        for position in positions:
            if label_type == 'liquidation':
                # Binary: was liquidated or not
                label = 1 if position.get('was_liquidated', False) else 0
            elif label_type == 'risk_level':
                # Multi-class: risk level
                health = position.get('health_factor', 2.0)
                if health < 1.2:
                    label = 3  # Critical
                elif health < 1.5:
                    label = 2  # High
                elif health < 2.0:
                    label = 1  # Medium
                else:
                    label = 0  # Low
            elif label_type == 'days_to_liquidation':
                # Regression: days until liquidation
                label = position.get('days_until_liquidation', 999)
            else:
                label = 0

            labels.append(label)

        return np.array(labels)