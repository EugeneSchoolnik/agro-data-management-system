package repository

import (
	"fmt"
	"time"

	"agro-data-management-system/internal/models"

	"github.com/jmoiron/sqlx"
)

type ForecastRepository interface {
	Create(f models.Forecast) (int, error)
	GetByID(id int) (models.Forecast, error)
	GetLatestByField(fieldID int) (models.Forecast, error)
	GetHistoryByField(fieldID int) ([]models.Forecast, error)
}

type ForecastPostgres struct {
	db *sqlx.DB
}

func NewForecastPostgres(db *sqlx.DB) *ForecastPostgres {
	return &ForecastPostgres{db: db}
}

func (r *ForecastPostgres) Create(f models.Forecast) (int, error) {
	var id int
	query := `INSERT INTO forecasts (field_id, pest_id, probability, recommendation, created_at) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`

	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now().UTC()
	}

	err := r.db.QueryRow(query, f.FieldID, f.PestID, f.Probability, f.Recommendation, f.CreatedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to save forecast: %w", err)
	}
	return id, nil
}

func (r *ForecastPostgres) GetByID(id int) (models.Forecast, error) {
	var f models.Forecast
	query := `SELECT * FROM forecasts WHERE id = $1`
	err := r.db.Get(&f, query, id)
	return f, err
}

func (r *ForecastPostgres) GetLatestByField(fieldID int) (models.Forecast, error) {
	var f models.Forecast
	query := `SELECT * FROM forecasts WHERE field_id = $1 ORDER BY created_at DESC LIMIT 1`
	err := r.db.Get(&f, query, fieldID)
	return f, err
}

func (r *ForecastPostgres) GetHistoryByField(fieldID int) ([]models.Forecast, error) {
	var forecasts []models.Forecast
	query := `SELECT * FROM forecasts WHERE field_id = $1 ORDER BY created_at DESC`
	err := r.db.Select(&forecasts, query, fieldID)
	return forecasts, err
}
