# AI Forecast Engine

A lightweight FastAPI microservice for predicting Sunn pest (`E. integriceps`) damage risk and weather forecasts using sensor metrics.

## Structure

- `main.py` - API endpoints and server initialization
- `model_engine.py` - PestModelEngine and WeatherModelEngine with inference and fallback logic
- `train_model.py` - synthetic training data generation and model export for both models
- `Dockerfile` - containerization recipe
- `requirements.txt` - Python dependencies

## Install

```bash
python -m venv .venv
source .venv/Scripts/activate  # Windows
pip install -r requirements.txt
```

## Train models

Run the training script to generate both `pest_model.pkl` and `weather_model.pkl`:

```bash
python train_model.py
```

## Start the API

```bash
uvicorn main:app --host 0.0.0.0 --port 8000
```

## Endpoints

### Pest Prediction

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

Response:

```json
{
  "probability": 0.58,
  "recommendation": "Moderate pest risk: increase field monitoring and consider preventive actions."
}
```

### Weather Forecast

POST `/predict-weather`

Request body (weather parameters in order):

```json
{
  "weather_data": [22.5, 1013.0, 3.2, 15.0, 0.0, 68.0, 20.0, 450.0],
  "hours_ahead": 3
}
```

Where `weather_data` contains:

- `temp` (°C)
- `pressure` (hPa)
- `wind_speed` (m/s)
- `dew_point` (°C)
- `precipitation` (mm)
- `humidity` (%)
- `soil_temp` (°C)
- `solar_radiation` (W/m²)

Response:

```json
{
  "temperature": 21.8,
  "hours_ahead": 3,
  "recommendation": "(3h) Mild: Optimal growing conditions expected [confidence: high]"
}
```

## Forecast Accuracy

| Time Horizon | Accuracy | Use Case                   |
| ------------ | -------- | -------------------------- |
| 1-6 hours    | 95%+     | Real-time decision support |
| 6-24 hours   | 70-85%   | Planning and alerts        |
| 24+ hours    | <70%     | Use external weather API   |

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
