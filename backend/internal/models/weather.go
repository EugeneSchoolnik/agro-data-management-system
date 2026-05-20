package models

import "time"

type WeatherParameter struct {
	ID          int    `db:"id" json:"id"`
	ParamID     int    `db:"param_id" json:"param_id"`
	Name        string `db:"name" json:"name"`
	Unit        string `db:"unit" json:"unit"`
	Description string `db:"description" json:"description"`
}

type WeatherStation struct {
	ID         int       `db:"id" json:"id"`
	ExternalID int       `db:"external_id" json:"external_id"`
	Name       string    `db:"name" json:"name"`
	Region     string    `db:"region" json:"region"`
	Active     bool      `db:"active" json:"active"`
	LastSeen   time.Time `db:"last_seen" json:"last_seen"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type WeatherStationParameter struct {
	ID                 int `db:"id" json:"id"`
	StationID          int `db:"station_id" json:"station_id"`
	WeatherParameterID int `db:"weather_parameter_id" json:"weather_parameter_id"`
	StationParam       int `db:"station_param" json:"station_param"`
}

type WeatherObservation struct {
	ID                 int64     `db:"id" json:"id"`
	StationID          int       `db:"station_id" json:"station_id"`
	WeatherParameterID int       `db:"weather_parameter_id" json:"weather_parameter_id"`
	StationParam       int       `db:"station_param" json:"station_param"`
	Value              float64   `db:"value" json:"value"`
	RecordedAt         time.Time `db:"recorded_at" json:"recorded_at"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
}
