package service

import (
	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type SensorService interface {
	Register(sensor models.Sensor) (models.Sensor, error)
	GetByID(id int) (models.Sensor, error)
	GetByField(fieldID int) ([]models.Sensor, error)
	UpdateStatus(id int, status models.SensorStatus) error
	Delete(id int) error
}

type sensorService struct {
	sensorRepo repository.SensorRepository
	fieldRepo  repository.FieldRepository // Перевірка прив'язки
	log        *zap.Logger
	validate   *validator.Validate
}

func NewSensorService(sr repository.SensorRepository, fr repository.FieldRepository, log *zap.Logger) SensorService {
	return &sensorService{
		sensorRepo: sr,
		fieldRepo:  fr,
		log:        log,
		validate:   validator.New(),
	}
}

func (s *sensorService) Register(sensor models.Sensor) (models.Sensor, error) {
	// 1. Валідація структури
	if err := s.validate.Struct(sensor); err != nil {
		return sensor, fmt.Errorf("sensor validation failed: %w", err)
	}

	// 2. Перевірка статусу через метод моделі
	if err := sensor.Status.IsValid(); err != nil {
		return sensor, err
	}

	// 3. БІЗНЕС-ЛОГІКА: Чи існує поле, куди ми ставимо датчик?
	_, err := s.fieldRepo.GetByIDWithCrop(sensor.FieldID)
	if err != nil {
		s.log.Warn("Attempt to register sensor for unknown field", zap.Int("field_id", sensor.FieldID))
		return sensor, fmt.Errorf("field %d not found", sensor.FieldID)
	}

	// 4. Збереження
	sensorCreated, err := s.sensorRepo.Create(sensor)
	if err != nil {
		s.log.Error("Failed to register sensor", zap.Error(err))
		return sensor, err
	}

	s.log.Info("Sensor registered", zap.Int("id", sensorCreated.ID), zap.String("type", sensorCreated.SensorType))
	return sensorCreated, nil
}

func (s *sensorService) GetByID(id int) (models.Sensor, error) {
	return s.sensorRepo.GetByID(id)
}

func (s *sensorService) GetByField(fieldID int) ([]models.Sensor, error) {
	return s.sensorRepo.GetByFieldID(fieldID)
}

func (s *sensorService) UpdateStatus(id int, status models.SensorStatus) error {
	if err := status.IsValid(); err != nil {
		return err
	}
	s.log.Info("Updating sensor status", zap.Int("id", id), zap.String("status", string(status)))
	return s.sensorRepo.UpdateStatus(id, status)
}

func (s *sensorService) Delete(id int) error {
	return s.sensorRepo.Delete(id)
}
