package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"
	"agro-data-management-system/internal/weatherapi"

	"go.uber.org/zap"
)

type WeatherService interface {
	ListStations() ([]models.WeatherStation, error)
	GetLatestObservations(stationID int) ([]models.WeatherObservation, error)
	GetParameterByParamID(paramID int) (models.WeatherParameter, error)
	GetStationByExternalID(externalID int) (models.WeatherStation, error)
	SyncStation(ctx context.Context, externalID int) ([]models.WeatherObservation, error)
	SyncField(ctx context.Context, fieldID int) ([]models.WeatherObservation, error)
}

type weatherService struct {
	repo   repository.WeatherRepository
	client *weatherapi.Client
	log    *zap.Logger
}

func NewWeatherService(repo repository.WeatherRepository, client *weatherapi.Client, log *zap.Logger) WeatherService {
	return &weatherService{repo: repo, client: client, log: log}
}

func (s *weatherService) ListStations() ([]models.WeatherStation, error) {
	return s.repo.GetAllStations()
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

func (s *weatherService) SyncStation(ctx context.Context, externalID int) ([]models.WeatherObservation, error) {
	station, err := s.repo.GetStationByExternalID(externalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	resp, err := s.client.StationLastDataGet(ctx, externalID)
	if err != nil {
		return nil, err
	}

	var saved []models.WeatherObservation
	for _, param := range resp.Res {
		weatherParam, err := s.repo.GetParameterByParamID(param.ParamID)
		if err != nil {
			s.log.Warn("Unknown weather parameter", zap.Int("param_id", param.ParamID))
			continue
		}

		recordedAt, err := parseMeteoTimestamp(param.Date)
		if err != nil {
			s.log.Warn("Unable to parse date", zap.String("date", param.Date), zap.Error(err))
			continue
		}

		obs := models.WeatherObservation{
			StationID:          station.ID,
			WeatherParameterID: weatherParam.ID,
			StationParam:       param.StationParam,
			Value:              param.Value,
			RecordedAt:         recordedAt,
		}

		created, err := s.repo.CreateObservation(obs)
		if err != nil {
			return nil, err
		}
		saved = append(saved, created)
	}

	return saved, nil
}

func (s *weatherService) SyncField(ctx context.Context, fieldID int) ([]models.WeatherObservation, error) {
	resp, err := s.client.FieldLastDataGet(ctx, fieldID)
	if err != nil {
		return nil, err
	}

	var saved []models.WeatherObservation
	for _, param := range resp.Res {
		weatherParam, err := s.repo.GetParameterByParamID(param.ParamID)
		if err != nil {
			s.log.Warn("Unknown weather parameter", zap.Int("param_id", param.ParamID))
			continue
		}

		for _, stationValue := range param.StationValue {
			station, err := s.repo.GetStationByExternalID(stationValue.StationID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					station = models.WeatherStation{
						ExternalID: stationValue.StationID,
						Name:       stationValue.StationName,
						Active:     true,
					}
					station, err = s.repo.CreateStation(station)
					if err != nil {
						return nil, err
					}
				} else {
					return nil, err
				}
			}

			if _, err := s.repo.GetStationParameter(station.ID, weatherParam.ID); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					_, err = s.repo.CreateStationParameter(models.WeatherStationParameter{
						StationID:          station.ID,
						WeatherParameterID: weatherParam.ID,
						StationParam:       stationValue.StationParam,
					})
					if err != nil {
						return nil, err
					}
				} else {
					return nil, err
				}
			}

			recordedAt, err := parseMeteoTimestamp(stationValue.Date)
			if err != nil {
				s.log.Warn("Unable to parse date", zap.String("date", stationValue.Date), zap.Error(err))
				continue
			}

			obs := models.WeatherObservation{
				StationID:          station.ID,
				WeatherParameterID: weatherParam.ID,
				StationParam:       stationValue.StationParam,
				Value:              stationValue.Value,
				RecordedAt:         recordedAt,
			}

			created, err := s.repo.CreateObservation(obs)
			if err != nil {
				return nil, err
			}
			saved = append(saved, created)
		}
	}

	return saved, nil
}

func parseMeteoTimestamp(value string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02T15:04:05", value)
}
