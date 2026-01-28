package repository

import (
	"agro-data-management-system/internal/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type FieldRepository interface {
	Create(field models.Field) (models.Field, error)
	GetByID(id int) (models.Field, error)
	GetByIDWithCrop(id int) (models.FieldWithCrop, error)
	GetAll() ([]models.FieldWithCrop, error)
	Update(field models.Field) error
	Delete(id int) error
}

type FieldPostgres struct {
	db *sqlx.DB
}

func NewFieldPostgres(db *sqlx.DB) *FieldPostgres {
	return &FieldPostgres{db: db}
}

func (r *FieldPostgres) Create(f models.Field) (models.Field, error) {
	query := `INSERT INTO fields (name, area, location, crop_id) 
              VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	row := r.db.QueryRow(query, f.Name, f.Area, f.Location, f.CropID)
	if err := row.Scan(&f.ID, &f.CreatedAt); err != nil {
		return f, fmt.Errorf("failed to create field: %w", err)
	}
	return f, nil
}

func (r *FieldPostgres) GetByID(id int) (models.Field, error) {
	var field models.Field
	query := `SELECT id, name, area, location, crop_id, created_at FROM fields WHERE id = $1`
	err := r.db.Get(&field, query, id)
	if err == sql.ErrNoRows {
		return field, fmt.Errorf("field not found")
	}
	return field, err
}

func (r *FieldPostgres) GetByIDWithCrop(id int) (models.FieldWithCrop, error) {
	var res models.FieldWithCrop
	query := `
		SELECT f.*, c.name as crop_name, c.variety as crop_variety 
		FROM fields f 
		LEFT JOIN crops c ON f.crop_id = c.id 
		WHERE f.id = $1`
	err := r.db.Get(&res, query, id)
	return res, err
}

func (r *FieldPostgres) GetAll() ([]models.FieldWithCrop, error) {
	var fields []models.FieldWithCrop
	query := `
		SELECT f.*, c.name as crop_name, c.variety as crop_variety 
		FROM fields f 
		LEFT JOIN crops c ON f.crop_id = c.id 
		ORDER BY f.created_at DESC`
	err := r.db.Select(&fields, query)
	return fields, err
}

func (r *FieldPostgres) Update(f models.Field) error {
	query := `UPDATE fields SET name=$1, area=$2, location=$3, crop_id=$4 WHERE id=$5`
	_, err := r.db.Exec(query, f.Name, f.Area, f.Location, f.CropID, f.ID)
	return err
}

func (r *FieldPostgres) Delete(id int) error {
	query := `DELETE FROM fields WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
