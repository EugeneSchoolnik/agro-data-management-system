package service

import (
	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type PestService interface {
	Create(pest models.Pest) (models.Pest, error)
	GetByID(id int) (models.Pest, error)
	GetAll() ([]models.Pest, error)
	Update(pest models.Pest) error
	Delete(id int) error
}

type pestService struct {
	repo     repository.PestRepository
	log      *zap.Logger
	validate *validator.Validate
}

func NewPestService(repo repository.PestRepository, log *zap.Logger) PestService {
	return &pestService{
		repo:     repo,
		log:      log,
		validate: validator.New(),
	}
}

func (s *pestService) Create(p models.Pest) (models.Pest, error) {
	if err := s.validate.Struct(p); err != nil {
		return p, fmt.Errorf("pest validation failed: %w", err)
	}

	pest, err := s.repo.Create(p)
	if err != nil {
		s.log.Error("Failed to add new pest to dictionary", zap.Error(err), zap.String("pest", p.Name))
		return p, err
	}

	s.log.Info("New pest added to system", zap.Int("id", pest.ID), zap.String("latin", pest.ScientificName))
	return pest, nil
}

func (s *pestService) GetByID(id int) (models.Pest, error) {
	return s.repo.GetByID(id)
}

func (s *pestService) GetAll() ([]models.Pest, error) {
	return s.repo.GetAll()
}

func (s *pestService) Update(p models.Pest) error {
	if err := s.validate.Struct(p); err != nil {
		return err
	}
	return s.repo.Update(p)
}

func (s *pestService) Delete(id int) error {
	s.log.Warn("Deleting pest from dictionary", zap.Int("id", id))
	return s.repo.Delete(id)
}
