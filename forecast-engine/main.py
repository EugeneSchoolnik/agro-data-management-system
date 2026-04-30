import logging
from typing import List

from fastapi import FastAPI, HTTPException, Request, status
from fastapi.responses import JSONResponse
from pydantic import BaseModel, ValidationError, validator

from model_engine import PestModelEngine

logger = logging.getLogger("forecast-engine")


class PredictionRequest(BaseModel):
    crop_name: str
    variety: str
    pest_name: str
    metrics: List[float]

    @validator("metrics")
    def validate_metrics(cls, value):
        if len(value) != 4:
            raise ValueError("metrics must contain exactly four values")
        return value


class PredictionResponse(BaseModel):
    probability: float
    recommendation: str


def configure_logging() -> None:
    handler = logging.StreamHandler()
    formatter = logging.Formatter(
        "%(asctime)s %(levelname)s %(name)s %(message)s"
    )
    handler.setFormatter(formatter)
    root_logger = logging.getLogger()
    root_logger.setLevel(logging.INFO)
    if not root_logger.handlers:
        root_logger.addHandler(handler)

    logger.info("Logging configured for AI Forecast Engine")


configure_logging()
app = FastAPI(
    title="AI Forecast Engine",
    description="Predicts Sunn pest damage risk from sensor metrics.",
    version="0.1.0",
)
engine = PestModelEngine()


@app.exception_handler(ValidationError)
async def validation_exception_handler(request: Request, exc: ValidationError):
    logger.warning("Validation error: %s", exc)
    return JSONResponse(
        status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
        content={"detail": exc.errors()},
    )


@app.exception_handler(Exception)
async def generic_exception_handler(request: Request, exc: Exception):
    logger.exception("Unexpected error: %s", exc)
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={"detail": "Internal server error"},
    )


@app.get("/health")
async def health_check():
    return {"status": "ok"}


@app.post("/predict", response_model=PredictionResponse)
async def predict(payload: PredictionRequest):
    try:
        probability, recommendation = engine.predict(payload.metrics)
        logger.info(
            "Prediction requested for crop=%s pest=%s variety=%s probability=%f",
            payload.crop_name,
            payload.pest_name,
            payload.variety,
            probability,
        )
        return PredictionResponse(
            probability=probability,
            recommendation=recommendation,
        )
    except ValueError as exc:
        logger.warning("Invalid request data: %s", exc)
        raise HTTPException(status_code=400, detail=str(exc))
