package repository

import (
	"agro-data-management-system/internal/models"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type SensorRepository interface {
	Create(sensor models.Sensor) (models.Sensor, error)
	GetByID(id int) (models.Sensor, error)
	GetByFieldID(fieldID int) ([]models.Sensor, error)
	UpdateStatus(id int, status models.SensorStatus) error
	Delete(id int) error
}

type SensorPostgres struct {
	db *sqlx.DB
}

func NewSensorPostgres(db *sqlx.DB) *SensorPostgres {
	return &SensorPostgres{db: db}
}

func (r *SensorPostgres) Create(s models.Sensor) (models.Sensor, error) {
	// Валідація "enum" перед записом
	if err := s.Status.IsValid(); err != nil {
		return s, err
	}

	query := `INSERT INTO sensors (field_id, sensor_type, status, last_sync) 
              VALUES ($1, $2, $3, $4) RETURNING id`

	s.LastSync = time.Now()
	row := r.db.QueryRow(query, s.FieldID, s.SensorType, s.Status, s.LastSync)
	if err := row.Scan(&s.ID); err != nil {
		return s, fmt.Errorf("failed to create sensor: %w", err)
	}
	return s, nil
}

func (r *SensorPostgres) GetByID(id int) (models.Sensor, error) {
	var sensor models.Sensor
	query := `SELECT * FROM sensors WHERE id = $1`
	err := r.db.Get(&sensor, query, id)
	return sensor, err
}

func (r *SensorPostgres) GetByFieldID(fieldID int) ([]models.Sensor, error) {
	var sensors []models.Sensor
	query := `SELECT * FROM sensors WHERE field_id = $1 ORDER BY id ASC`
	err := r.db.Select(&sensors, query, fieldID)
	return sensors, err
}

func (r *SensorPostgres) UpdateStatus(id int, status models.SensorStatus) error {
	if err := status.IsValid(); err != nil {
		return err
	}

	query := `UPDATE sensors SET status = $1, last_sync = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *SensorPostgres) Delete(id int) error {
	query := `DELETE FROM sensors WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
