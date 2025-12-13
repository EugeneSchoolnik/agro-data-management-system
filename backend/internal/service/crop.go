package service

import (
	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type CropService interface {
	Create(crop models.Crop) (int, error)
	GetByID(id int) (models.Crop, error)
	GetAll() ([]models.Crop, error)
	Update(crop models.Crop) error
	Delete(id int) error
}

type cropService struct {
	repo     repository.CropRepository
	log      *zap.Logger
	validate *validator.Validate
}

func NewCropService(repo repository.CropRepository, log *zap.Logger) CropService {
	return &cropService{
		repo:     repo,
		log:      log,
		validate: validator.New(),
	}
}

func (s *cropService) Create(c models.Crop) (int, error) {
	// 1. Валідація даних
	if err := s.validate.Struct(c); err != nil {
		s.log.Warn("Validation failed for Crop creation", zap.Error(err))
		return 0, fmt.Errorf("invalid input: %w", err)
	}

	// 2. Виклик репозиторію
	id, err := s.repo.Create(c)
	if err != nil {
		s.log.Error("Failed to create crop in database", zap.Error(err), zap.String("crop_name", c.Name))
		return 0, err
	}

	s.log.Info("New crop created successfully", zap.Int("id", id), zap.String("name", c.Name))
	return id, nil
}

func (s *cropService) GetByID(id int) (models.Crop, error) {
	crop, err := s.repo.GetByID(id)
	if err != nil {
		s.log.Warn("Crop not found", zap.Int("id", id))
		return crop, err
	}
	return crop, nil
}

func (s *cropService) GetAll() ([]models.Crop, error) {
	return s.repo.GetAll()
}

func (s *cropService) Update(c models.Crop) error {
	if err := s.validate.Struct(c); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	err := s.repo.Update(c)
	if err != nil {
		s.log.Error("Failed to update crop", zap.Int("id", c.ID), zap.Error(err))
		return err
	}

	s.log.Info("Crop updated", zap.Int("id", c.ID), zap.String("name", c.Name))
	return nil
}

func (s *cropService) Delete(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		s.log.Error("Failed to delete crop", zap.Int("id", id), zap.Error(err))
		return err
	}

	s.log.Info("Crop deleted from system", zap.Int("id", id))
	return nil
}