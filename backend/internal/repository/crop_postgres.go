package repository

import (
	"agro-data-management-system/internal/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CropRepository interface {
	Create(crop models.Crop) (int, error)
	GetByID(id int) (models.Crop, error)
	GetAll() ([]models.Crop, error)
	Update(crop models.Crop) error
	Delete(id int) error
}

type CropPostgres struct {
	db *sqlx.DB
}

func NewCropPostgres(db *sqlx.DB) *CropPostgres {
	return &CropPostgres{db: db}
}

func (r *CropPostgres) Create(c models.Crop) (int, error) {
	var id int
	query := `INSERT INTO crops (name, variety, description) 
              VALUES ($1, $2, $3) RETURNING id`

	row := r.db.QueryRow(query, c.Name, c.Variety, c.Description)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to create crop: %w", err)
	}
	return id, nil
}

func (r *CropPostgres) GetByID(id int) (models.Crop, error) {
	var crop models.Crop
	query := `SELECT id, name, variety, description FROM crops WHERE id = $1`
	err := r.db.Get(&crop, query, id)
	if err == sql.ErrNoRows {
		return crop, fmt.Errorf("crop not found")
	}
	return crop, err
}

func (r *CropPostgres) GetAll() ([]models.Crop, error) {
	var crops []models.Crop
	query := `SELECT id, name, variety, description FROM crops ORDER BY name ASC`
	err := r.db.Select(&crops, query)
	return crops, err
}

func (r *CropPostgres) Update(c models.Crop) error {
	query := `UPDATE crops SET name=$1, variety=$2, description=$3 WHERE id=$4`
	res, err := r.db.Exec(query, c.Name, c.Variety, c.Description, c.ID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no crop found to update")
	}
	return nil
}

func (r *CropPostgres) Delete(id int) error {
	query := `DELETE FROM crops WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
