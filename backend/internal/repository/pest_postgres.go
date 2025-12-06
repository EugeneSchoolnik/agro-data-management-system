package repository

import (
	"fmt"
	"agro-data-management-system/internal/models"
	"github.com/jmoiron/sqlx"
)

type PestRepository interface {
	Create(pest models.Pest) (int, error)
	GetByID(id int) (models.Pest, error)
	GetAll() ([]models.Pest, error)
	Update(pest models.Pest) error
	Delete(id int) error
}

type PestPostgres struct {
	db *sqlx.DB
}

func NewPestPostgres(db *sqlx.DB) *PestPostgres {
	return &PestPostgres{db: db}
}

func (r *PestPostgres) Create(p models.Pest) (int, error) {
	var id int
	query := `INSERT INTO pests (name, scientific_name) VALUES ($1, $2) RETURNING id`
	
	err := r.db.QueryRow(query, p.Name, p.ScientificName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create pest: %w", err)
	}
	return id, nil
}

func (r *PestPostgres) GetByID(id int) (models.Pest, error) {
	var p models.Pest
	// Явно вказуємо поля, бо в моделі Go немає description, який є в БД
	query := `SELECT id, name, scientific_name FROM pests WHERE id = $1`
	err := r.db.Get(&p, query, id)
	return p, err
}

func (r *PestPostgres) GetAll() ([]models.Pest, error) {
	var pests []models.Pest
	query := `SELECT id, name, scientific_name FROM pests ORDER BY name ASC`
	err := r.db.Select(&pests, query)
	return pests, err
}

func (r *PestPostgres) Update(p models.Pest) error {
	query := `UPDATE pests SET name = $1, scientific_name = $2 WHERE id = $3`
	_, err := r.db.Exec(query, p.Name, p.ScientificName, p.ID)
	return err
}

func (r *PestPostgres) Delete(id int) error {
	query := `DELETE FROM pests WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}