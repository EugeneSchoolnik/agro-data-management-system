import logging
import os
from typing import List, Tuple

import joblib
import numpy as np

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
