package models

import (
	"fmt"
	"time"
)

// SensorStatus — наш власний тип для "enum"
type SensorStatus string

const (
	StatusActive   SensorStatus = "active"
	StatusInactive SensorStatus = "inactive"
	StatusError    SensorStatus = "error"
	StatusTesting  SensorStatus = "testing"
)

// IsValid перевіряє, чи входить статус у дозволений список
func (s SensorStatus) IsValid() error {
	switch s {
	case StatusActive, StatusInactive, StatusError, StatusTesting:
		return nil
	}
	return fmt.Errorf("invalid sensor status: %s", s)
}

type Sensor struct {
	ID      int `db:"id" json:"id"`
	FieldID int `db:"field_id" json:"field_id" validate:"required,gt=0"`
	// Тип датчика (напр. "humidity", "temperature", "camera")
	SensorType string       `db:"sensor_type" json:"sensor_type" validate:"required,min=3"`
	Status     SensorStatus `db:"status" json:"status" validate:"required"`
	LastSync   time.Time    `db:"last_sync" json:"last_sync"`
}

type Metric struct {
	ID         int64     `db:"id" json:"id"`
	SensorID   int       `db:"sensor_id" json:"sensor_id"`
	Value      float64   `db:"value" json:"value"`
	RecordedAt time.Time `db:"recorded_at" json:"recorded_at"`
}
