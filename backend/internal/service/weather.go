package service

import (
	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"

	"go.uber.org/zap"
)

type WeatherService interface {
	RecordObservation(obs models.WeatherObservation) (models.WeatherObservation, error)
	GetLatestObservations(stationID int) ([]models.WeatherObservation, error)
	GetParameterByParamID(paramID int) (models.WeatherParameter, error)
	GetStationByExternalID(externalID int) (models.WeatherStation, error)
}

type weatherService struct {
	repo repository.WeatherRepository
	log  *zap.Logger
}

func NewWeatherService(repo repository.WeatherRepository, log *zap.Logger) WeatherService {
	return &weatherService{repo: repo, log: log}
}

func (s *weatherService) RecordObservation(obs models.WeatherObservation) (models.WeatherObservation, error) {
	return s.repo.CreateObservation(obs)
}

func (s *weatherService) GetLatestObservations(stationID int) ([]models.WeatherObservation, error) {
	return s.repo.GetLatestObservationsByStation(stationID)
}

func (s *weatherService) GetParameterByParamID(paramID int) (models.WeatherParameter, error) {
	return s.repo.GetParameterByParamID(paramID)
}

func (s *weatherService) GetStationByExternalID(externalID int) (models.WeatherStation, error) {
	return s.repo.GetStationByExternalID(externalID)
}
