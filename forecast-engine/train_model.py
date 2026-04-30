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


if __name__ == "__main__":
    train_and_export_model()
