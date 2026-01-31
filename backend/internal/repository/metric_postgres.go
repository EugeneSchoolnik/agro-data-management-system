package repository

import (
	"fmt"
	"time"

	"agro-data-management-system/internal/models"

	"github.com/jmoiron/sqlx"
)

type MetricRepository interface {
	Create(metric models.Metric) (models.Metric, error)
	GetLatestBySensor(sensorID int) (models.Metric, error)
	GetHistoryBySensor(sensorID int, from, to time.Time) ([]models.Metric, error)
	GetAggregatedMetricsByField(fieldID int, from, to time.Time) ([]models.MetricAggregate, error)
}

type MetricPostgres struct {
	db *sqlx.DB
}

func NewMetricPostgres(db *sqlx.DB) *MetricPostgres {
	return &MetricPostgres{db: db}
}

func (r *MetricPostgres) Create(m models.Metric) (models.Metric, error) {
	query := `INSERT INTO metrics (sensor_id, value, recorded_at) 
              VALUES ($1, $2, $3) RETURNING id`

	if m.RecordedAt.IsZero() {
		m.RecordedAt = time.Now()
	}

	err := r.db.QueryRow(query, m.SensorID, m.Value, m.RecordedAt).Scan(&m.ID)
	if err != nil {
		return m, fmt.Errorf("failed to save metric: %w", err)
	}
	return m, nil
}

func (r *MetricPostgres) GetLatestBySensor(sensorID int) (models.Metric, error) {
	var m models.Metric
	query := `SELECT id, sensor_id, value, recorded_at FROM metrics 
              WHERE sensor_id = $1 ORDER BY recorded_at DESC LIMIT 1`
	err := r.db.Get(&m, query, sensorID)
	return m, err
}

func (r *MetricPostgres) GetHistoryBySensor(sensorID int, from, to time.Time) ([]models.Metric, error) {
	metrics := []models.Metric{}
	query := `SELECT id, sensor_id, value, recorded_at FROM metrics 
              WHERE sensor_id = $1 AND recorded_at BETWEEN $2 AND $3 
              ORDER BY recorded_at ASC`
	err := r.db.Select(&metrics, query, sensorID, from, to)
	return metrics, err
}

func (r *MetricPostgres) GetAggregatedMetricsByField(fieldID int, from, to time.Time) ([]models.MetricAggregate, error) {
	aggregates := []models.MetricAggregate{}
	query := `
SELECT s.sensor_type,
       AVG(m.value) AS avg_value,
       MIN(m.value) AS min_value,
       MAX(m.value) AS max_value
FROM metrics m
JOIN sensors s ON m.sensor_id = s.id
WHERE s.field_id = $1
  AND s.sensor_type IN ('temperature', 'air_humidity')
  AND m.recorded_at BETWEEN $2 AND $3
GROUP BY s.sensor_type`
	err := r.db.Select(&aggregates, query, fieldID, from, to)
	return aggregates, err
}
