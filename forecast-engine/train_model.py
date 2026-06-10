import logging

import joblib
import numpy as np
from sklearn.ensemble import RandomForestRegressor
from sklearn.metrics import mean_squared_error
from sklearn.model_selection import train_test_split


def build_synthetic_dataset(sample_count: int = 900, seed: int = 42):
    rng = np.random.default_rng(seed)
    temperature = rng.uniform(5, 35, size=sample_count)
    humidity = rng.uniform(20, 95, size=sample_count)
    soil_moisture = rng.uniform(5, 85, size=sample_count)
    crop_phase = rng.integers(1, 5, size=sample_count)

    temp_norm = np.clip((temperature - 10.0) / 20.0, 0.0, 1.0)
    humidity_norm = humidity / 100.0
    moisture_norm = soil_moisture / 100.0
    phase_norm = (crop_phase - 1) / 3.0

    hydro_coefficient = (humidity_norm + moisture_norm) / np.clip(temp_norm, 0.1, 1.0)
    probability = (
        0.25 * humidity_norm
        + 0.20 * (1.0 - moisture_norm)
        + 0.25 * phase_norm
        + 0.20 * np.clip(hydro_coefficient / 2.0, 0.0, 1.0)
        + 0.10 * rng.normal(0.0, 0.05, size=sample_count)
    )
    probability = np.clip(probability, 0.0, 1.0)

    features = np.column_stack(
        [temperature, humidity, soil_moisture, crop_phase.astype(float)]
    )
    return features, probability


def build_weather_synthetic_dataset(sample_count: int = 1200, seed: int = 42):
    """
    Build synthetic weather dataset for temperature forecasting.
    Weather features: [temp, pressure, wind_speed, dew_point, 
                       precipitation, humidity, soil_temp, solar_radiation, hours_ahead]
    Target: next hour temperature
    """
    rng = np.random.default_rng(seed)
    
    # Current observations
    current_temp = rng.uniform(0, 35, size=sample_count)
    pressure = rng.uniform(990, 1030, size=sample_count)
    wind_speed = rng.uniform(0, 15, size=sample_count)
    dew_point = current_temp - rng.uniform(5, 25, size=sample_count)
    precipitation = rng.exponential(2, size=sample_count)  # Mostly low values
    humidity = rng.uniform(30, 95, size=sample_count)
    soil_temp = rng.uniform(-5, 35, size=sample_count)
    solar_radiation = rng.uniform(0, 1000, size=sample_count)
    hours_ahead = rng.integers(1, 25, size=sample_count)
    
    # Generate target temperature (next hour)
    # Factors: persistence, solar radiation effect, time of day effect
    pressure_anomaly = (pressure - 1013.0) / 20.0  # -1 to +1
    solar_effect = (solar_radiation / 1000.0) * 3.0  # 0-3°C effect
    wind_cooling = wind_speed * 0.2  # Cooling effect
    time_decay = 1.0 - (hours_ahead / 100.0)  # Decay over forecast period
    
    next_temp = (
        current_temp * (1.0 - (hours_ahead / 100.0))  # Persistence with decay
        + pressure_anomaly * 0.5  # Pressure influence
        + solar_effect * time_decay
        - wind_cooling * time_decay
        + rng.normal(0.0, 0.5, size=sample_count)  # Random noise
    )
    
    features = np.column_stack([
        current_temp, pressure, wind_speed, dew_point,
        precipitation, humidity, soil_temp, solar_radiation, hours_ahead
    ])
    
    return features, next_temp


def train_and_export_model(model_path: str = "pest_model.pkl") -> None:
    logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
    logging.info("Building synthetic training dataset")
    X, y = build_synthetic_dataset()

    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.18, random_state=1
    )

    model = RandomForestRegressor(
        n_estimators=120,
        max_depth=8,
        random_state=3,
    )
    model.fit(X_train, y_train)

    predictions = model.predict(X_test)
    mse = mean_squared_error(y_test, predictions)
    joblib.dump(model, model_path)

    logging.info("Model saved to %s", model_path)
    logging.info("Test MSE: %.5f", mse)
    logging.info("Sample prediction: %.4f", predictions[0])


def train_and_export_weather_model(model_path: str = "weather_model.pkl") -> None:
    """Train and export weather forecasting model"""
    logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
    logging.info("Building synthetic weather training dataset")
    X, y = build_weather_synthetic_dataset()

    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.2, random_state=42
    )

    model = RandomForestRegressor(
        n_estimators=100,
        max_depth=10,
        random_state=42,
        n_jobs=-1,
    )
    model.fit(X_train, y_train)

    predictions = model.predict(X_test)
    mse = mean_squared_error(y_test, predictions)
    joblib.dump(model, model_path)

    logging.info("Weather model saved to %s", model_path)
    logging.info("Test MSE: %.5f", mse)
    logging.info("Sample prediction: %.2f°C", predictions[0])


if __name__ == "__main__":
    train_and_export_model()
    train_and_export_weather_model()
