package repository

import (
	"agro-data-management-system/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type WeatherRepository interface {
	CreateParameter(param models.WeatherParameter) (models.WeatherParameter, error)
	GetParameterByParamID(paramID int) (models.WeatherParameter, error)
	GetParameterByID(id int) (models.WeatherParameter, error)
	CreateStation(station models.WeatherStation) (models.WeatherStation, error)
	GetStationByExternalID(externalID int) (models.WeatherStation, error)
	GetAllStations() ([]models.WeatherStation, error)
	GetStationParameter(stationID, weatherParameterID int) (models.WeatherStationParameter, error)
	CreateStationParameter(mapping models.WeatherStationParameter) (models.WeatherStationParameter, error)
	CreateObservation(obs models.WeatherObservation) (models.WeatherObservation, error)
	GetLatestObservationsByStation(stationID int) ([]models.WeatherObservation, error)
}

type WeatherPostgres struct {
	db *sqlx.DB
}

func NewWeatherPostgres(db *sqlx.DB) *WeatherPostgres {
	return &WeatherPostgres{db: db}
}

func (r *WeatherPostgres) CreateParameter(param models.WeatherParameter) (models.WeatherParameter, error) {
	query := `INSERT INTO weather_parameters (param_id, name, unit, description) VALUES ($1, $2, $3, $4) RETURNING id`
	row := r.db.QueryRow(query, param.ParamID, param.Name, param.Unit, param.Description)
	if err := row.Scan(&param.ID); err != nil {
		return param, fmt.Errorf("failed to create weather parameter: %w", err)
	}
	return param, nil
}

func (r *WeatherPostgres) GetParameterByParamID(paramID int) (models.WeatherParameter, error) {
	var param models.WeatherParameter
	query := `SELECT * FROM weather_parameters WHERE param_id = $1`
	err := r.db.Get(&param, query, paramID)
	return param, err
}

func (r *WeatherPostgres) GetParameterByID(id int) (models.WeatherParameter, error) {
	var param models.WeatherParameter
	query := `SELECT * FROM weather_parameters WHERE id = $1`
	err := r.db.Get(&param, query, id)
	return param, err
}

func (r *WeatherPostgres) CreateStation(station models.WeatherStation) (models.WeatherStation, error) {
	query := `INSERT INTO weather_stations (external_id, name, region, active, last_seen) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	row := r.db.QueryRow(query, station.ExternalID, station.Name, station.Region, station.Active, station.LastSeen)
	if err := row.Scan(&station.ID, &station.CreatedAt); err != nil {
		return station, fmt.Errorf("failed to create weather station: %w", err)
	}
	return station, nil
}

func (r *WeatherPostgres) GetStationByExternalID(externalID int) (models.WeatherStation, error) {
	var station models.WeatherStation
	query := `SELECT * FROM weather_stations WHERE external_id = $1`
	err := r.db.Get(&station, query, externalID)
	return station, err
}

func (r *WeatherPostgres) GetAllStations() ([]models.WeatherStation, error) {
	var stations []models.WeatherStation
	query := `SELECT * FROM weather_stations ORDER BY external_id`
	err := r.db.Select(&stations, query)
	return stations, err
}

func (r *WeatherPostgres) GetStationParameter(stationID, weatherParameterID int) (models.WeatherStationParameter, error) {
	var mapping models.WeatherStationParameter
	query := `SELECT * FROM weather_station_parameters WHERE station_id = $1 AND weather_parameter_id = $2`
	err := r.db.Get(&mapping, query, stationID, weatherParameterID)
	return mapping, err
}

func (r *WeatherPostgres) CreateStationParameter(mapping models.WeatherStationParameter) (models.WeatherStationParameter, error) {
	query := `INSERT INTO weather_station_parameters (station_id, weather_parameter_id, station_param) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRow(query, mapping.StationID, mapping.WeatherParameterID, mapping.StationParam)
	if err := row.Scan(&mapping.ID); err != nil {
		return mapping, fmt.Errorf("failed to map station parameter: %w", err)
	}
	return mapping, nil
}

func (r *WeatherPostgres) CreateObservation(obs models.WeatherObservation) (models.WeatherObservation, error) {
	query := `INSERT INTO weather_observations (station_id, weather_parameter_id, station_param, value, recorded_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	row := r.db.QueryRow(query, obs.StationID, obs.WeatherParameterID, obs.StationParam, obs.Value, obs.RecordedAt)
	if err := row.Scan(&obs.ID, &obs.CreatedAt); err != nil {
		return obs, fmt.Errorf("failed to create weather observation: %w", err)
	}
	return obs, nil
}

func (r *WeatherPostgres) GetLatestObservationsByStation(stationID int) ([]models.WeatherObservation, error) {
	var observations []models.WeatherObservation
	query := `SELECT * FROM weather_observations WHERE station_id = $1 ORDER BY recorded_at DESC LIMIT 100`
	err := r.db.Select(&observations, query, stationID)
	return observations, err
}
