package repository

import (
	"agro-data-management-system/internal/config"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Драйвер для PostgreSQL
)

func NewPostgresDB(cfg config.DBConfig) (*sqlx.DB, error) {
	// Формуємо рядок підключення (DSN)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	// Відкриваємо з'єднання
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// Перевіряємо з'єднання (Ping)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return db, nil
}