package service

import (
	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type FieldService interface {
	Create(field models.Field) (int, error)
	GetByID(id int) (models.FieldWithCrop, error)
	GetAll() ([]models.FieldWithCrop, error)
	Update(field models.Field) error
	Delete(id int) error
}

type fieldService struct {
	fieldRepo repository.FieldRepository
	cropRepo  repository.CropRepository
	log       *zap.Logger
	validate  *validator.Validate
}

func NewFieldService(fr repository.FieldRepository, cr repository.CropRepository, log *zap.Logger) FieldService {
	return &fieldService{
		fieldRepo: fr,
		cropRepo:  cr,
		log:       log,
		validate:  validator.New(),
	}
}

func (s *fieldService) Create(f models.Field) (int, error) {
	// 1. Валідація структури
	if err := s.validate.Struct(f); err != nil {
		return 0, fmt.Errorf("field validation failed: %w", err)
	}

	// 2. Перевірка існування культури
	_, err := s.cropRepo.GetByID(f.CropID)
	if err != nil {
		s.log.Warn("Crop check failed for field creation", zap.Int("crop_id", f.CropID), zap.Error(err))
		return 0, fmt.Errorf("cannot create field: crop with id %d not found", f.CropID)
	}

	// 3. Збереження
	id, err := s.fieldRepo.Create(f)
	if err != nil {
		s.log.Error("Failed to save field", zap.Error(err))
		return 0, err
	}

	s.log.Info("Field created successfully", zap.Int("id", id), zap.String("name", f.Name))
	return id, nil
}

func (s *fieldService) GetByID(id int) (models.FieldWithCrop, error) {
	return s.fieldRepo.GetByIDWithCrop(id)
}

func (s *fieldService) GetAll() ([]models.FieldWithCrop, error) {
	return s.fieldRepo.GetAll()
}

func (s *fieldService) Update(f models.Field) error {
	if err := s.validate.Struct(f); err != nil {
		return err
	}

	// Також перевіряємо культуру при оновленні (якщо її змінили)
	if _, err := s.cropRepo.GetByID(f.CropID); err != nil {
		return fmt.Errorf("invalid crop_id: %w", err)
	}

	return s.fieldRepo.Update(f)
}

func (s *fieldService) Delete(id int) error {
	// В майбутньому тут можна додати перевірку, чи немає на полі активних датчиків
	return s.fieldRepo.Delete(id)
}
