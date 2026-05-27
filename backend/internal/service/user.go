package service

import (
    "agro-data-management-system/internal/models"
    "agro-data-management-system/internal/repository"
    "go.uber.org/zap"
)

type UserService interface {
    GetByEmail(email string) (*models.User, error)
}

type userService struct{
    repo repository.UserRepository
    log  *zap.Logger
}

func NewUserService(repo repository.UserRepository, log *zap.Logger) UserService {
    return &userService{repo: repo, log: log}
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
    return s.repo.GetByEmail(email)
}
