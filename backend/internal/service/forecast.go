package service

import (
	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"
	"bytes"
	"encoding/json"
	"fmt"
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
	pestService   PestService
	aiServiceURL  string // URL вашого Python-сервісу
	log           *zap.Logger
}

func NewForecastService(
	repo repository.ForecastRepository,
	ms MetricService,
	fs FieldService,
	ps PestService,
	aiURL string,
	log *zap.Logger,
) ForecastService {
	return &forecastService{
		repo:          repo,
		metricService: ms,
		fieldService:  fs,
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

	// 3. Збираємо історію метрик (наприклад, за останні 7 днів)
	// Для спрощення беремо метрики одного основного датчика поля
	// У реальній системі тут може бути агрегація по всіх датчиках поля
	to := time.Now().UTC()
	from := to.AddDate(0, 0, -7)

	// Припустимо, ми знаємо ID датчика або беремо історію по полю (потрібен метод у MetricService)
	// Тут для прикладу використовуємо заглушку отримання метрик
	history, err := s.metricService.GetHistory(fieldID, from, to)
	if err != nil {
		return models.Forecast{}, fmt.Errorf("failed to get metrics history: %w", err)
	}

	values := make([]float64, len(history))
	for i, m := range history {
		values[i] = m.Value
	}

	// 4. Формуємо запит до ШІ
	reqBody := models.ForecastRequest{
		CropName: field.CropName,
		Variety:  field.CropVariety,
		Metrics:  values,
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

	id, err := s.repo.Create(forecast)
	if err != nil {
		return models.Forecast{}, err
	}
	forecast.ID = id

	s.log.Info("Forecast generated", zap.Int("field_id", fieldID), zap.Float64("prob", aiResp.Probability))
	return forecast, nil
}

func (s *forecastService) GetLatest(fieldID int) (models.Forecast, error) {
	return s.repo.GetLatestByField(fieldID)
}
