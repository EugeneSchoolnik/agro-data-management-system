package repository

import (
    "agro-data-management-system/internal/models"
    "github.com/jmoiron/sqlx"
)

type UserRepository interface {
    GetByEmail(email string) (*models.User, error)
}

type userPostgres struct{
    db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) UserRepository {
    return &userPostgres{db: db}
}

func (r *userPostgres) GetByEmail(email string) (*models.User, error) {
    var u models.User
    query := `SELECT id, email, password_hash, role, created_at FROM users WHERE lower(email) = lower($1) LIMIT 1`
    if err := r.db.Get(&u, query, email); err != nil {
        return nil, err
    }
    return &u, nil
}
