# AI Forecast Engine

A lightweight FastAPI microservice for predicting Sunn pest (`E. integriceps`) damage risk using sensor metrics.

## Structure

- `main.py` - API endpoints and server initialization
- `model_engine.py` - model loading, inference, and fallback expert logic
- `train_model.py` - synthetic training data generation and model export
- `Dockerfile` - containerization recipe
- `requirements.txt` - Python dependencies

## Install

```bash
python -m venv .venv
source .venv/Scripts/activate  # Windows
pip install -r requirements.txt
```

## Train the model

Run the training script to generate `pest_model.pkl`:

```bash
python train_model.py
```

## Start the API

```bash
uvicorn main:app --host 0.0.0.0 --port 8000
```

## Predict

POST `/predict`

Request body:

```json
{
  "crop_name": "wheat",
  "variety": "spring",
  "pest_name": "E. integriceps",
  "metrics": [22.5, 68.0, 34.0, 3]
}
```

Response body:

```json
{
  "probability": 0.58,
  "recommendation": "Moderate pest risk: increase field monitoring and consider preventive actions."
}
```

## Docker

Build image:

```bash
docker build -t ai-forecast-engine .
```

Run container:

```bash
docker run --rm -p 8000:8000 ai-forecast-engine
```
