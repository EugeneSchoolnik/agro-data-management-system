package service

import (
	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type MetricService interface {
	Save(metric models.Metric) (int64, error)
	GetLatest(sensorID int) (models.Metric, error)
	GetHistory(sensorID int, from, to time.Time) ([]models.Metric, error)
}

type metricService struct {
	metricRepo repository.MetricRepository
	sensorRepo repository.SensorRepository // Для перевірки статусу датчика
	log        *zap.Logger
	validate   *validator.Validate
}

func NewMetricService(mr repository.MetricRepository, sr repository.SensorRepository, log *zap.Logger) MetricService {
	return &metricService{
		metricRepo: mr,
		sensorRepo: sr,
		log:        log,
		validate:   validator.New(),
	}
}

func (s *metricService) Save(m models.Metric) (int64, error) {
	// 1. Валідація структури
	if err := s.validate.Struct(m); err != nil {
		return 0, fmt.Errorf("metric validation failed: %w", err)
	}

	// 2. ПЕРЕВІРКА ДАТЧИКА: Чи існує він і чи він ACTIVE?
	sensor, err := s.sensorRepo.GetByID(m.SensorID)
	if err != nil {
		s.log.Warn("Received metric from unknown sensor", zap.Int("sensor_id", m.SensorID))
		return 0, fmt.Errorf("sensor %d not found", m.SensorID)
	}

	if sensor.Status != models.StatusActive {
		s.log.Warn("Ignoring metric from inactive sensor",
			zap.Int("sensor_id", m.SensorID),
			zap.String("status", string(sensor.Status)))
		return 0, fmt.Errorf("sensor %d is not active (current status: %s)", m.SensorID, sensor.Status)
	}

	// 3. Збереження
	id, err := s.metricRepo.Create(m)
	if err != nil {
		s.log.Error("Failed to save metric", zap.Error(err))
		return 0, err
	}

	s.log.Debug("Metric saved", zap.Int64("id", id), zap.Float64("value", m.Value))
	return id, nil
}

func (s *metricService) GetLatest(sensorID int) (models.Metric, error) {
	return s.metricRepo.GetLatestBySensor(sensorID)
}

func (s *metricService) GetHistory(sensorID int, from, to time.Time) ([]models.Metric, error) {
	if from.After(to) {
		return nil, fmt.Errorf("'from' date cannot be after 'to' date")
	}
	return s.metricRepo.GetHistoryBySensor(sensorID, from, to)
}
