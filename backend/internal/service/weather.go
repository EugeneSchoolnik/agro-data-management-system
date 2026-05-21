package service

import (
	"context"
	"database/sql"
	"errors"
	"sort"
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
	GetParameterByID(id int) (models.WeatherParameter, error)
	GetStationWeatherSummary(externalID int) (models.WeatherStationSummary, error)
	GetStationByExternalID(externalID int) (models.WeatherStation, error)
	SyncStation(ctx context.Context, externalID int) ([]models.WeatherObservation, error)
	SyncField(ctx context.Context, fieldID int) ([]models.WeatherObservation, error)
	StartPeriodicSync(ctx context.Context, interval time.Duration)
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

func (s *weatherService) GetParameterByID(id int) (models.WeatherParameter, error) {
	return s.repo.GetParameterByID(id)
}

func (s *weatherService) GetStationByExternalID(externalID int) (models.WeatherStation, error) {
	return s.repo.GetStationByExternalID(externalID)
}

func (s *weatherService) GetStationWeatherSummary(externalID int) (models.WeatherStationSummary, error) {
	station, err := s.repo.GetStationByExternalID(externalID)
	if err != nil {
		return models.WeatherStationSummary{}, err
	}

	observations, err := s.repo.GetLatestObservationsByStation(station.ID)
	if err != nil {
		return models.WeatherStationSummary{}, err
	}

	latestByParam := make(map[int]models.WeatherObservation)
	for _, obs := range observations {
		if _, ok := latestByParam[obs.WeatherParameterID]; ok {
			continue
		}
		latestByParam[obs.WeatherParameterID] = obs
	}

	published := make([]models.WeatherParameterSummary, 0, len(latestByParam))
	var updatedAt time.Time

	for _, obs := range latestByParam {
		param, err := s.repo.GetParameterByID(obs.WeatherParameterID)
		if err != nil {
			s.log.Warn("Unknown weather parameter", zap.Int("weather_parameter_id", obs.WeatherParameterID), zap.Error(err))
			continue
		}

		if obs.RecordedAt.After(updatedAt) {
			updatedAt = obs.RecordedAt
		}

		published = append(published, models.WeatherParameterSummary{
			Parameter:    param,
			Value:        obs.Value,
			StationParam: obs.StationParam,
			RecordedAt:   obs.RecordedAt,
		})
	}

	dayEnd := time.Now().UTC()
	dayStart := dayEnd.Add(-24 * time.Hour)

	dailyObservations, err := s.repo.GetObservationsByStationAndRecordedAtRange(station.ID, dayStart, dayEnd)
	if err != nil {
		return models.WeatherStationSummary{}, err
	}

	type aggregate struct {
		sum   float64
		min   float64
		max   float64
		count int
	}

	dailyByParam := make(map[int]*aggregate)
	trendByParam := make(map[int]map[time.Time]*aggregate)
	trendStart := dayEnd.Add(-6 * time.Hour)

	for _, obs := range dailyObservations {
		agg, ok := dailyByParam[obs.WeatherParameterID]
		if !ok {
			agg = &aggregate{sum: obs.Value, min: obs.Value, max: obs.Value, count: 1}
			dailyByParam[obs.WeatherParameterID] = agg
		} else {
			agg.sum += obs.Value
			if obs.Value < agg.min {
				agg.min = obs.Value
			}
			if obs.Value > agg.max {
				agg.max = obs.Value
			}
			agg.count++
		}

		if obs.RecordedAt.After(trendStart) {
			hour := obs.RecordedAt.Truncate(time.Hour)
			hours, ok := trendByParam[obs.WeatherParameterID]
			if !ok {
				hours = make(map[time.Time]*aggregate)
				trendByParam[obs.WeatherParameterID] = hours
			}

			hourAgg, ok := hours[hour]
			if !ok {
				hourAgg = &aggregate{sum: obs.Value, min: obs.Value, max: obs.Value, count: 1}
				hours[hour] = hourAgg
			} else {
				hourAgg.sum += obs.Value
				if obs.Value < hourAgg.min {
					hourAgg.min = obs.Value
				}
				if obs.Value > hourAgg.max {
					hourAgg.max = obs.Value
				}
				hourAgg.count++
			}
		}
	}

	daily := make([]models.WeatherParameterAggregate, 0, len(dailyByParam))
	for paramID, agg := range dailyByParam {
		param, err := s.repo.GetParameterByID(paramID)
		if err != nil {
			s.log.Warn("Unknown weather parameter for daily aggregate", zap.Int("weather_parameter_id", paramID), zap.Error(err))
			continue
		}
		daily = append(daily, models.WeatherParameterAggregate{
			Parameter: param,
			Average:   agg.sum / float64(agg.count),
			Min:       agg.min,
			Max:       agg.max,
			Count:     agg.count,
		})
	}

	sort.Slice(daily, func(i, j int) bool {
		return daily[i].Parameter.ParamID < daily[j].Parameter.ParamID
	})

	hourlyTrend := make([]models.WeatherParameterTrend, 0, len(trendByParam))
	for paramID, hours := range trendByParam {
		param, err := s.repo.GetParameterByID(paramID)
		if err != nil {
			s.log.Warn("Unknown weather parameter for hourly trend", zap.Int("weather_parameter_id", paramID), zap.Error(err))
			continue
		}

		points := make([]models.HourlyTrendPoint, 0, len(hours))
		for hour, agg := range hours {
			points = append(points, models.HourlyTrendPoint{
				Hour:  hour,
				Value: agg.sum / float64(agg.count),
			})
		}

		sort.Slice(points, func(i, j int) bool {
			return points[i].Hour.Before(points[j].Hour)
		})

		hourlyTrend = append(hourlyTrend, models.WeatherParameterTrend{
			Parameter: param,
			Points:    points,
		})
	}

	sort.Slice(hourlyTrend, func(i, j int) bool {
		return hourlyTrend[i].Parameter.ParamID < hourlyTrend[j].Parameter.ParamID
	})

	sort.Slice(published, func(i, j int) bool {
		return published[i].Parameter.ParamID < published[j].Parameter.ParamID
	})

	return models.WeatherStationSummary{
		Station:     station,
		Latest:      published,
		Daily:       daily,
		HourlyTrend: hourlyTrend,
		UpdatedAt:   updatedAt,
	}, nil
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

func (s *weatherService) StartPeriodicSync(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = 30 * time.Minute
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	s.log.Info("Weather periodic sync started", zap.Duration("interval", interval))
	if err := s.syncAllStations(ctx); err != nil {
		s.log.Error("Weather periodic sync failed", zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			s.log.Info("Weather periodic sync stopped")
			return
		case <-ticker.C:
			if err := s.syncAllStations(ctx); err != nil {
				s.log.Error("Weather periodic sync failed", zap.Error(err))
			}
		}
	}
}

func (s *weatherService) syncAllStations(ctx context.Context) error {
	stations, err := s.repo.GetAllStations()
	if err != nil {
		return err
	}

	for _, station := range stations {
		if _, err := s.SyncStation(ctx, station.ExternalID); err != nil {
			s.log.Error("Failed to sync weather station", zap.Int("external_id", station.ExternalID), zap.Error(err))
		}
	}

	return nil
}

func parseMeteoTimestamp(value string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02T15:04:05", value)
}
