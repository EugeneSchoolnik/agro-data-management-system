package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	Crop     CropRepository
	Field    FieldRepository
	Sensor   SensorRepository
	Metric   MetricRepository
	Pest     PestRepository
	Forecast ForecastRepository
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Crop:     NewCropPostgres(db),
		Field:    NewFieldPostgres(db),
		Sensor:   NewSensorPostgres(db),
		Metric:   NewMetricPostgres(db),
		Pest:     NewPestPostgres(db),
		Forecast: NewForecastPostgres(db),
	}
}
