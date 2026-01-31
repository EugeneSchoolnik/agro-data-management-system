package service

import (
	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ForecastService interface {
	Predict(fieldID, pestID int) (models.Forecast, error)
	GetLatest(fieldID int) (models.Forecast, error)
}

type forecastService struct {
	repo          repository.ForecastRepository
	metricService MetricService
	fieldService  FieldService
	sensorService SensorService
	pestService   PestService
	aiServiceURL  string // URL вашого Python-сервісу
	log           *zap.Logger
}

func NewForecastService(
	repo repository.ForecastRepository,
	ms MetricService,
	fs FieldService,
	ss SensorService,
	ps PestService,
	aiURL string,
	log *zap.Logger,
) ForecastService {
	return &forecastService{
		repo:          repo,
		metricService: ms,
		fieldService:  fs,
		sensorService: ss,
		pestService:   ps,
		aiServiceURL:  aiURL,
		log:           log,
	}
}

func (s *forecastService) Predict(fieldID, pestID int) (models.Forecast, error) {
	// 1. Збираємо дані про поле та культуру
	field, err := s.fieldService.GetByID(fieldID)
	if err != nil {
		return models.Forecast{}, fmt.Errorf("failed to get field info: %w", err)
	}

	// 2. Отримуємо дані про шкідника
	pest, err := s.pestService.GetByID(pestID)
	if err != nil {
		return models.Forecast{}, fmt.Errorf("failed to get pest info: %w", err)
	}

	// 3. Збираємо сенсори поля і дані метрик за останні 7 днів
	to := time.Now().UTC()
	from := to.AddDate(0, 0, -7)

	sensors, err := s.sensorService.GetByField(fieldID)
	if err != nil {
		return models.Forecast{}, fmt.Errorf("failed to get field sensors: %w", err)
	}

	sensorIDs := map[string]int{}
	for _, sensor := range sensors {
		sensorIDs[sensor.SensorType] = sensor.ID
	}

	required := []string{"temperature", "air_humidity", "soil_moisture", "crop_phase"}
	for _, sensorType := range required {
		if _, ok := sensorIDs[sensorType]; !ok {
			return models.Forecast{}, fmt.Errorf("missing required sensor: %s", sensorType)
		}
	}

	temperatureHistory, err := s.metricService.GetHistory(sensorIDs["temperature"], from, to)
	if err != nil {
		return models.Forecast{}, fmt.Errorf("failed to get temperature history: %w", err)
	}
	if len(temperatureHistory) == 0 {
		return models.Forecast{}, fmt.Errorf("no temperature metrics available")
	}
	temperature := averageMetrics(temperatureHistory)

	airHumidityHistory, err := s.metricService.GetHistory(sensorIDs["air_humidity"], from, to)
	if err != nil {
		return models.Forecast{}, fmt.Errorf("failed to get air humidity history: %w", err)
	}
	if len(airHumidityHistory) == 0 {
		return models.Forecast{}, fmt.Errorf("no air humidity metrics available")
	}
	airHumidity := averageMetrics(airHumidityHistory)

	soilMoistureHistory, err := s.metricService.GetHistory(sensorIDs["soil_moisture"], from, to)
	if err != nil {
		return models.Forecast{}, fmt.Errorf("failed to get soil moisture history: %w", err)
	}
	if len(soilMoistureHistory) == 0 {
		return models.Forecast{}, fmt.Errorf("no soil moisture metrics available")
	}
	soilMoisture := averageMetrics(soilMoistureHistory)

	cropPhaseHistory, err := s.metricService.GetHistory(sensorIDs["crop_phase"], from, to)
	if err != nil {
		return models.Forecast{}, fmt.Errorf("failed to get crop phase history: %w", err)
	}
	if len(cropPhaseHistory) == 0 {
		return models.Forecast{}, fmt.Errorf("no crop phase metrics available")
	}
	cropPhase := math.Round(latestMetric(cropPhaseHistory))

	metrics := []float64{temperature, airHumidity, soilMoisture, cropPhase}

	// 4. Формуємо запит до ШІ
	reqBody := models.ForecastRequest{
		CropName: field.CropName,
		Variety:  field.CropVariety,
		Metrics:  metrics,
		PestName: pest.ScientificName,
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(s.aiServiceURL+"/predict", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		s.log.Error("AI Engine connection error", zap.Error(err))
		return models.Forecast{}, fmt.Errorf("AI engine unavailable")
	}
	defer resp.Body.Close()

	var aiResp models.ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return models.Forecast{}, fmt.Errorf("failed to decode AI response")
	}

	// 5. Зберігаємо результат у БД
	forecast := models.Forecast{
		FieldID:        fieldID,
		PestID:         pestID,
		Probability:    aiResp.Probability,
		Recommendation: aiResp.Recommendation,
		CreatedAt:      time.Now().UTC(),
	}

	createdForecast, err := s.repo.Create(forecast)
	if err != nil {
		return models.Forecast{}, err
	}

	s.log.Info("Forecast generated", zap.Int("field_id", fieldID), zap.Float64("prob", aiResp.Probability))
	return createdForecast, nil
}

func averageMetrics(history []models.Metric) float64 {
	sum := 0.0
	for _, item := range history {
		sum += item.Value
	}
	return sum / float64(len(history))
}

func latestMetric(history []models.Metric) float64 {
	return history[len(history)-1].Value
}

func (s *forecastService) GetLatest(fieldID int) (models.Forecast, error) {
	return s.repo.GetLatestByField(fieldID)
}
