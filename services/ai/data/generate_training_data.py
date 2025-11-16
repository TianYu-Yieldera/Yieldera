#!/usr/bin/env python3
"""
Generate Training Data for ML Models
Creates realistic synthetic data for model training
"""

import numpy as np
import pandas as pd
from datetime import datetime, timedelta
import json
import os
import random
from typing import List, Dict, Tuple
import pickle

# Set random seeds for reproducibility
np.random.seed(42)
random.seed(42)

# ============================================================
# DATA GENERATION PARAMETERS
# ============================================================

class DataGenerator:
    """Generate synthetic training data for AI risk models"""

    def __init__(self):
        self.protocols = ['aave', 'compound', 'gmx', 'maker']
        self.assets = ['ETH', 'BTC', 'USDC', 'USDT', 'DAI', 'LINK', 'UNI']
        self.n_users = 10000
        self.n_positions = 50000
        self.n_liquidations = 5000
        self.time_range_days = 365

    def generate_all_datasets(self, output_dir: str = 'training_data'):
        """Generate all training datasets"""
        os.makedirs(output_dir, exist_ok=True)

        print("🔧 Generating training datasets...")

        # 1. Price history
        price_history = self.generate_price_history()
        self.save_dataset(price_history, f"{output_dir}/price_history.csv")
        print(f"  ✅ Price history: {len(price_history)} records")

        # 2. Liquidation history
        liquidation_history = self.generate_liquidation_history()
        self.save_dataset(liquidation_history, f"{output_dir}/liquidation_history.csv")
        print(f"  ✅ Liquidation history: {len(liquidation_history)} records")

        # 3. Position snapshots
        position_snapshots = self.generate_position_snapshots()
        self.save_dataset(position_snapshots, f"{output_dir}/position_snapshots.csv")
        print(f"  ✅ Position snapshots: {len(position_snapshots)} records")

        # 4. User profiles
        user_profiles = self.generate_user_profiles()
        self.save_dataset(user_profiles, f"{output_dir}/user_profiles.csv")
        print(f"  ✅ User profiles: {len(user_profiles)} records")

        # 5. ML training data
        X_train, y_train = self.generate_ml_training_data()
        self.save_ml_data(X_train, y_train, output_dir)
        print(f"  ✅ ML training data: {len(X_train)} samples")

        # 6. Time series data
        time_series = self.generate_time_series_data()
        self.save_time_series(time_series, output_dir)
        print(f"  ✅ Time series data: {len(time_series)} sequences")

        print(f"\n✅ All datasets saved to {output_dir}/")

    # ============================================================
    # PRICE HISTORY GENERATION
    # ============================================================

    def generate_price_history(self) -> pd.DataFrame:
        """Generate realistic price history with volatility"""
        records = []
        start_date = datetime.now() - timedelta(days=self.time_range_days)

        # Initial prices
        base_prices = {
            'ETH': 2000,
            'BTC': 40000,
            'USDC': 1.0,
            'USDT': 1.0,
            'DAI': 1.0,
            'LINK': 15,
            'UNI': 10
        }

        for asset, base_price in base_prices.items():
            # Generate price path using Geometric Brownian Motion
            prices = self._generate_gbm_prices(
                base_price,
                days=self.time_range_days,
                volatility=0.02 if asset in ['USDC', 'USDT', 'DAI'] else 0.05
            )

            for day in range(self.time_range_days):
                timestamp = start_date + timedelta(days=day)

                # Add multiple data points per day
                for hour in [0, 6, 12, 18]:
                    # Add some intraday volatility
                    intraday_factor = 1 + np.random.normal(0, 0.005)
                    price = prices[day] * intraday_factor

                    records.append({
                        'time': timestamp + timedelta(hours=hour),
                        'asset': asset,
                        'price': price,
                        'source': np.random.choice(['chainlink', 'uniswap_twap', 'coingecko']),
                        'volume': abs(np.random.normal(1e6, 5e5)) if asset not in ['USDC', 'USDT', 'DAI'] else abs(np.random.normal(1e7, 1e6)),
                        'market_cap': price * 1e8  # Simplified
                    })

        return pd.DataFrame(records)

    def _generate_gbm_prices(self, S0: float, days: int, volatility: float) -> np.ndarray:
        """Generate price path using Geometric Brownian Motion"""
        dt = 1  # Daily
        drift = 0.0001  # Small positive drift
        prices = np.zeros(days)
        prices[0] = S0

        for t in range(1, days):
            dW = np.random.normal(0, np.sqrt(dt))
            prices[t] = prices[t-1] * np.exp(
                (drift - 0.5 * volatility**2) * dt + volatility * dW
            )

        return prices

    # ============================================================
    # LIQUIDATION HISTORY GENERATION
    # ============================================================

    def generate_liquidation_history(self) -> pd.DataFrame:
        """Generate realistic liquidation events"""
        records = []
        start_date = datetime.now() - timedelta(days=self.time_range_days)

        for i in range(self.n_liquidations):
            # Random timestamp
            days_ago = random.randint(0, self.time_range_days)
            timestamp = start_date + timedelta(days=days_ago)

            # Generate correlated features (unhealthy positions more likely to liquidate)
            health_factor = np.random.beta(2, 5) * 1.5  # Most are low
            leverage = np.random.gamma(2, 2)  # Right-skewed

            records.append({
                'id': i,
                'time': timestamp,
                'protocol': random.choice(self.protocols),
                'user_address': f"0x{random.randint(0, 2**160-1):040x}",
                'collateral_asset': random.choice(['ETH', 'BTC']),
                'debt_asset': random.choice(['USDC', 'USDT', 'DAI']),
                'collateral_amount': abs(np.random.lognormal(2, 1)),
                'debt_amount': abs(np.random.lognormal(3, 1)),
                'liquidation_price': abs(np.random.normal(2000, 500)),
                'health_factor': health_factor,
                'gas_price': abs(np.random.gamma(2, 25)),
                'tx_hash': f"0x{random.randint(0, 2**256-1):064x}",
                'block_number': random.randint(1000000, 2000000)
            })

        return pd.DataFrame(records)

    # ============================================================
    # POSITION SNAPSHOTS GENERATION
    # ============================================================

    def generate_position_snapshots(self) -> pd.DataFrame:
        """Generate position snapshot data"""
        records = []
        start_date = datetime.now() - timedelta(days=self.time_range_days)

        for i in range(self.n_positions):
            # User and position characteristics
            user_address = f"0x{random.randint(0, 2**160-1):040x}"
            is_risky_user = random.random() < 0.2  # 20% are risky

            # Generate correlated features
            if is_risky_user:
                health_factor = np.random.uniform(1.0, 1.5)
                leverage = np.random.uniform(3, 10)
                will_liquidate = random.random() < 0.4
            else:
                health_factor = np.random.uniform(1.5, 3.0)
                leverage = np.random.uniform(1, 3)
                will_liquidate = random.random() < 0.05

            collateral_amount = abs(np.random.lognormal(3, 1))
            collateral_asset = random.choice(['ETH', 'BTC'])
            debt_asset = random.choice(['USDC', 'USDT', 'DAI'])

            # Price based on asset
            collateral_price = 2000 if collateral_asset == 'ETH' else 40000
            collateral_value = collateral_amount * collateral_price
            debt_value = collateral_value / health_factor if health_factor > 0 else 0

            records.append({
                'id': i,
                'time': start_date + timedelta(days=random.randint(0, self.time_range_days)),
                'user_address': user_address,
                'protocol': random.choice(self.protocols),
                'position_id': f"pos_{i}",
                'collateral_asset': collateral_asset,
                'collateral_amount': collateral_amount,
                'debt_asset': debt_asset,
                'debt_amount': debt_value,  # Simplified
                'health_factor': health_factor,
                'ltv': 1 / health_factor if health_factor > 0 else 1,
                'leverage': leverage,
                'liquidation_price': collateral_price * 0.8,  # Simplified
                'collateral_value_usd': collateral_value,
                'debt_value_usd': debt_value,
                'was_liquidated': will_liquidate,
                'days_until_liquidation': random.randint(1, 30) if will_liquidate else None
            })

        return pd.DataFrame(records)

    # ============================================================
    # USER PROFILES GENERATION
    # ============================================================

    def generate_user_profiles(self) -> pd.DataFrame:
        """Generate user risk profiles"""
        records = []

        for i in range(self.n_users):
            # User types
            user_type = np.random.choice(
                ['conservative', 'moderate', 'aggressive', 'degen'],
                p=[0.3, 0.4, 0.25, 0.05]
            )

            # Type-based characteristics
            if user_type == 'conservative':
                risk_score = np.random.uniform(10, 30)
                avg_leverage = np.random.uniform(1, 1.5)
                liquidation_count = np.random.poisson(0.1)
            elif user_type == 'moderate':
                risk_score = np.random.uniform(30, 60)
                avg_leverage = np.random.uniform(1.5, 2.5)
                liquidation_count = np.random.poisson(0.5)
            elif user_type == 'aggressive':
                risk_score = np.random.uniform(60, 85)
                avg_leverage = np.random.uniform(2.5, 5)
                liquidation_count = np.random.poisson(2)
            else:  # degen
                risk_score = np.random.uniform(85, 100)
                avg_leverage = np.random.uniform(5, 10)
                liquidation_count = np.random.poisson(5)

            records.append({
                'user_address': f"0x{random.randint(0, 2**160-1):040x}",
                'risk_score': risk_score,
                'avg_leverage': avg_leverage,
                'liquidation_count': liquidation_count,
                'total_positions': np.random.poisson(5),
                'max_position_size': abs(np.random.lognormal(4, 1.5)),
                'preferred_assets': json.dumps(random.sample(self.assets, k=random.randint(1, 3))),
                'preferred_protocols': json.dumps(random.sample(self.protocols, k=random.randint(1, 2))),
                'avg_health_factor': 3.0 - (risk_score / 50),  # Inverse correlation
                'total_volume_usd': abs(np.random.lognormal(5, 2)),
                'last_activity': datetime.now() - timedelta(days=random.randint(0, 30)),
                'created_at': datetime.now() - timedelta(days=random.randint(30, 365))
            })

        return pd.DataFrame(records)

    # ============================================================
    # ML TRAINING DATA GENERATION
    # ============================================================

    def generate_ml_training_data(self) -> Tuple[np.ndarray, np.ndarray]:
        """Generate features and labels for ML training"""
        n_samples = 10000
        n_features = 20

        # Generate features
        X = np.zeros((n_samples, n_features))

        for i in range(n_samples):
            # Position features
            X[i, 0] = np.random.uniform(0.8, 3.0)  # health_factor
            X[i, 1] = np.random.uniform(0.3, 0.9)  # ltv
            X[i, 2] = np.random.uniform(1, 10)     # leverage
            X[i, 3] = np.log1p(abs(np.random.lognormal(3, 1)))  # log_collateral_value
            X[i, 4] = np.log1p(abs(np.random.lognormal(3, 1)))  # log_debt_value

            # Market features
            X[i, 5] = np.random.uniform(0.01, 0.1)  # volatility
            X[i, 6] = np.random.uniform(0.2, 0.95)  # utilization_rate
            X[i, 7] = np.random.uniform(0.1, 2.0)   # liquidity_ratio

            # Protocol encoding (one-hot)
            protocol_idx = random.randint(8, 11)
            X[i, protocol_idx] = 1.0

            # Time features
            X[i, 12] = np.random.uniform(0, 1)  # hour_of_day / 24
            X[i, 13] = np.random.uniform(0, 1)  # day_of_week / 7

            # Risk indicators
            X[i, 14] = 1.0 if X[i, 0] < 1.5 else 0.0  # low_health_factor
            X[i, 15] = 1.0 if X[i, 2] > 3.0 else 0.0   # high_leverage

            # User features
            X[i, 16] = np.random.uniform(0, 100)  # risk_score
            X[i, 17] = np.random.poisson(2)       # liquidation_count
            X[i, 18] = np.random.uniform(1, 5)    # avg_leverage
            X[i, 19] = np.random.uniform(0, 30)   # days_since_activity

        # Generate labels (liquidation within 24h)
        # Make it correlated with features
        y = np.zeros(n_samples)

        for i in range(n_samples):
            # Probability based on health factor and leverage
            base_prob = 0.01
            if X[i, 0] < 1.1:  # Very low health factor
                prob = 0.9
            elif X[i, 0] < 1.3:  # Low health factor
                prob = 0.5
            elif X[i, 0] < 1.5:  # Medium health factor
                prob = 0.2
            else:
                prob = base_prob

            # Adjust for leverage
            if X[i, 2] > 5:  # High leverage
                prob *= 2

            # Adjust for volatility
            if X[i, 5] > 0.05:  # High volatility
                prob *= 1.5

            y[i] = 1 if random.random() < min(prob, 1.0) else 0

        return X, y

    # ============================================================
    # TIME SERIES DATA GENERATION
    # ============================================================

    def generate_time_series_data(self) -> List[np.ndarray]:
        """Generate time series sequences for LSTM training"""
        n_sequences = 5000
        sequence_length = 60
        n_features = 10

        sequences = []

        for _ in range(n_sequences):
            # Generate correlated time series
            sequence = np.zeros((sequence_length, n_features))

            # Base trends
            trend = np.random.choice([-1, 0, 1])  # Down, flat, up
            volatility = np.random.uniform(0.01, 0.05)

            for t in range(sequence_length):
                # Price feature (with trend)
                if t == 0:
                    sequence[t, 0] = 1.0
                else:
                    sequence[t, 0] = sequence[t-1, 0] * (1 + trend * 0.001 + np.random.normal(0, volatility))

                # Volume (correlated with volatility)
                sequence[t, 1] = abs(np.random.normal(1, volatility * 10))

                # Health factor (decreasing if trend is down)
                sequence[t, 2] = 2.0 - t * 0.01 * max(0, -trend) + np.random.normal(0, 0.1)

                # Utilization rate
                sequence[t, 3] = min(0.95, abs(np.random.normal(0.6, 0.2)))

                # Other features
                for f in range(4, n_features):
                    sequence[t, f] = np.random.normal(0, 1)

            sequences.append(sequence)

        return sequences

    # ============================================================
    # SAVE UTILITIES
    # ============================================================

    def save_dataset(self, df: pd.DataFrame, filepath: str):
        """Save DataFrame to CSV"""
        df.to_csv(filepath, index=False)

    def save_ml_data(self, X: np.ndarray, y: np.ndarray, output_dir: str):
        """Save ML training data"""
        # Save as numpy arrays
        np.save(f"{output_dir}/X_train.npy", X)
        np.save(f"{output_dir}/y_train.npy", y)

        # Also save as pickle for complete preservation
        with open(f"{output_dir}/ml_training_data.pkl", 'wb') as f:
            pickle.dump({'X': X, 'y': y}, f)

    def save_time_series(self, sequences: List[np.ndarray], output_dir: str):
        """Save time series data"""
        # Convert to numpy array
        sequences_array = np.array(sequences)
        np.save(f"{output_dir}/time_series_sequences.npy", sequences_array)

        # Save metadata
        metadata = {
            'n_sequences': len(sequences),
            'sequence_length': sequences[0].shape[0],
            'n_features': sequences[0].shape[1]
        }
        with open(f"{output_dir}/time_series_metadata.json", 'w') as f:
            json.dump(metadata, f, indent=2)

# ============================================================
# DATA VALIDATION
# ============================================================

def validate_data(output_dir: str = 'training_data'):
    """Validate generated training data"""
    print("\n🔍 Validating generated data...")

    # Check files exist
    required_files = [
        'price_history.csv',
        'liquidation_history.csv',
        'position_snapshots.csv',
        'user_profiles.csv',
        'X_train.npy',
        'y_train.npy',
        'time_series_sequences.npy'
    ]

    all_valid = True
    for file in required_files:
        filepath = f"{output_dir}/{file}"
        if os.path.exists(filepath):
            size = os.path.getsize(filepath) / 1024 / 1024  # MB
            print(f"  ✅ {file}: {size:.2f} MB")
        else:
            print(f"  ❌ {file}: Not found")
            all_valid = False

    # Load and check ML data
    if all_valid:
        X = np.load(f"{output_dir}/X_train.npy")
        y = np.load(f"{output_dir}/y_train.npy")

        print(f"\n📊 ML Data Statistics:")
        print(f"  • Samples: {len(X)}")
        print(f"  • Features: {X.shape[1]}")
        print(f"  • Positive class: {np.mean(y):.2%}")
        print(f"  • Feature range: [{X.min():.2f}, {X.max():.2f}]")

        # Load time series
        sequences = np.load(f"{output_dir}/time_series_sequences.npy")
        print(f"\n📈 Time Series Data:")
        print(f"  • Sequences: {sequences.shape[0]}")
        print(f"  • Length: {sequences.shape[1]}")
        print(f"  • Features: {sequences.shape[2]}")

    return all_valid

# ============================================================
# MAIN EXECUTION
# ============================================================

def main():
    """Generate all training datasets"""
    print("="*60)
    print("AI RISK ENGINE - TRAINING DATA GENERATOR")
    print("="*60)

    # Create generator
    generator = DataGenerator()

    # Generate all datasets
    output_dir = 'training_data'
    generator.generate_all_datasets(output_dir)

    # Validate data
    if validate_data(output_dir):
        print("\n✅ All training data generated successfully!")
        print(f"📁 Data saved in: {os.path.abspath(output_dir)}/")

        print("\n📚 Next steps:")
        print("  1. Use this data to train ML models")
        print("  2. Run: python examples/example_usage.py")
        print("  3. Start API: python api/main_v2.py")
        print("  4. Import data to TimescaleDB for production")
    else:
        print("\n❌ Data validation failed!")

if __name__ == "__main__":
    main()