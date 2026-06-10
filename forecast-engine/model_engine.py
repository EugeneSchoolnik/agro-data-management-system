import logging
import os
from typing import List, Tuple

import joblib
import numpy as np
from sklearn.ensemble import RandomForestRegressor

logger = logging.getLogger("forecast-engine")


class PestModelEngine:
    def __init__(self, model_path: str = "pest_model.pkl"):
        self.model_path = model_path
        self.model = self._load_model()

    def _load_model(self):
        if os.path.isfile(self.model_path):
            try:
                model = joblib.load(self.model_path)
                logger.info("Loaded model from %s", self.model_path)
                return model
            except Exception as exc:
                logger.exception("Failed to load model %s", self.model_path)
                return None

        logger.warning("Model file %s not found; using rule-based fallback", self.model_path)
        return None

    def predict(self, metrics: List[float]) -> Tuple[float, str]:
        if len(metrics) != 4:
            raise ValueError("metrics must contain exactly 4 values")

        try:
            values = [float(x) for x in metrics]
        except (TypeError, ValueError) as exc:
            raise ValueError("metrics must be numeric values") from exc

        if self.model is not None:
            try:
                probability = float(self.model.predict(np.array(values).reshape(1, -1))[0])
                probability = self._clamp(probability)
                logger.info("Model inference success: probability=%f, metrics=%s", probability, values)
                return probability, self._format_recommendation(probability)
            except Exception:
                logger.exception("Model inference failed; falling back to rule-based logic")

        probability = self._rule_based_probability(values)
        recommendation = self._format_recommendation(probability, fallback=True)
        return probability, recommendation

    def _clamp(self, value: float) -> float:
        return max(0.0, min(1.0, value))

    def _rule_based_probability(self, metrics: List[float]) -> float:
        temp, humidity, soil_moisture, phase = metrics

        temp_score = max(0.0, min(1.0, (temp - 7.0) / 22.0))
        humidity_score = max(0.0, min(1.0, humidity / 100.0))
        moisture_score = max(0.0, min(1.0, soil_moisture / 100.0))
        phase_weight = min(max(int(round(phase)), 1), 5) / 5.0

        hydro_coefficient = (humidity_score + moisture_score) / max(0.1, temp_score)
        dryness_penalty = 1.0 - moisture_score

        base_risk = (
            0.30 * humidity_score
            + 0.20 * dryness_penalty
            + 0.25 * phase_weight
            + 0.15 * min(1.0, hydro_coefficient / 2.0)
            + 0.10 * temp_score
        )

        probability = self._clamp(base_risk)
        logger.info("Rule-based prediction: probability=%f, metrics=%s", probability, metrics)
        return probability

    def _format_recommendation(self, probability: float, fallback: bool = False) -> str:
        if probability >= 0.75:
            message = (
                "High pest damage risk: deploy targeted monitoring and control measures immediately."
            )
        elif probability >= 0.45:
            message = (
                "Moderate pest risk: increase field monitoring and consider preventive actions."
            )
        else:
            message = (
                "Low risk: maintain regular monitoring and favorable crop management."
            )

        if fallback:
            message += " (fallback expert system used)"

        return message


class WeatherModelEngine:
    def __init__(self, model_path: str = "weather_model.pkl"):
        self.model_path = model_path
        self.model = self._load_model()
        # Weather parameters: temp, pressure, wind_speed, dew_point, 
        # precipitation, humidity, soil_temp, solar_radiation
        self.feature_count = 8

    def _load_model(self) -> RandomForestRegressor:
        if os.path.isfile(self.model_path):
            try:
                model = joblib.load(self.model_path)
                logger.info("Loaded weather model from %s", self.model_path)
                return model
            except Exception as exc:
                logger.exception("Failed to load weather model %s", self.model_path)
                return None

        logger.warning("Weather model file %s not found; using simple fallback", self.model_path)
        return None

    def predict(self, weather_data: List[float], hours_ahead: int = 3) -> Tuple[float, str]:
        """
        Predict temperature ahead using weather parameters.
        
        weather_data: [temp, pressure, wind_speed, dew_point, 
                       precipitation, humidity, soil_temp, solar_radiation]
        hours_ahead: forecast hours (1-24)
        """
        if len(weather_data) != self.feature_count:
            raise ValueError(f"weather_data must contain exactly {self.feature_count} values")

        if not (1 <= hours_ahead <= 24):
            raise ValueError("hours_ahead must be between 1 and 24")

        try:
            values = np.array([float(x) for x in weather_data])
        except (TypeError, ValueError) as exc:
            raise ValueError("weather_data must be numeric values") from exc

        if self.model is not None:
            try:
                # Add hours_ahead as a feature for time context
                features = np.concatenate([values, [hours_ahead]])
                predicted_temp = float(self.model.predict(features.reshape(1, -1))[0])
                predicted_temp = self._clamp(predicted_temp, -50, 50)
                
                logger.info(
                    "Weather model inference success: temp_forecast=%.2f°C, hours=%d",
                    predicted_temp, hours_ahead
                )
                return predicted_temp, self._format_weather_recommendation(predicted_temp, hours_ahead)
            except Exception:
                logger.exception("Weather model inference failed; falling back to persistence")

        # Simple persistence + trend fallback
        predicted_temp = self._persistence_forecast(values, hours_ahead)
        recommendation = self._format_weather_recommendation(predicted_temp, hours_ahead, fallback=True)
        return predicted_temp, recommendation

    def _clamp(self, value: float, min_val: float = -50, max_val: float = 50) -> float:
        return max(min_val, min(max_val, value))

    def _persistence_forecast(self, weather_data: np.ndarray, hours_ahead: int) -> float:
        """Simple persistence + decay forecast: current temp with small trend"""
        current_temp = weather_data[0]
        humidity = weather_data[5]
        
        # Small trend based on conditions
        trend = 0.1 * (humidity / 100.0 - 0.5)  # Humid tends warmer
        decay = 1.0 - (hours_ahead / 100.0)  # Trend decreases over time
        
        predicted_temp = current_temp + (trend * decay)
        return self._clamp(predicted_temp)

    def _format_weather_recommendation(self, temp: float, hours_ahead: int, fallback: bool = False) -> str:
        """Generate recommendation based on predicted temperature"""
        confidence = "high" if not fallback else "low (fallback)"
        
        if temp < -5:
            status = "Frost risk"
            advice = "Protect sensitive crops from freezing"
        elif temp < 0:
            status = "Below freezing"
            advice = "Monitor frost conditions closely"
        elif temp < 10:
            status = "Cold"
            advice = "Suitable for cool-season crops"
        elif temp < 25:
            status = "Mild"
            advice = "Optimal growing conditions expected"
        elif temp < 35:
            status = "Warm"
            advice = "Monitor irrigation and heat stress"
        else:
            status = "Hot"
            advice = "Heat stress risk - ensure adequate irrigation"

        message = f"({hours_ahead}h) {status}: {advice} [confidence: {confidence}]"
        return message
