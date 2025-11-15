package models

import "time"

type Sensor struct {
	ID         int       `db:"id" json:"id"`
	FieldID    int       `db:"field_id" json:"field_id"`
	SensorType string    `db:"sensor_type" json:"sensor_type"`
	Status     string    `db:"status" json:"status"`
	LastSync   time.Time `db:"last_sync" json:"last_sync"`
}

type Metric struct {
	ID         int64     `db:"id" json:"id"`
	SensorID   int       `db:"sensor_id" json:"sensor_id"`
	Value      float64   `db:"value" json:"value"`
	RecordedAt time.Time `db:"recorded_at" json:"recorded_at"`
}