package repository

import (
	"os"
	"path/filepath"
	"testing"

	"agro-data-management-system/internal/config"
	"agro-data-management-system/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

// getProjectRoot знаходить корінь проєкту (де лежить go.mod)
func getProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// initTestConfig завантажує конфігурацію один раз
func initTestConfig(t *testing.T) *config.Config {
	root := getProjectRoot()

	// Завантажуємо .env з кореня
	_ = godotenv.Load(filepath.Join(root, ".env"))

	logger.Init("info")

	// Шлях до твого local.yaml
	cfg, err := config.LoadConfig(filepath.Join(root, "config/local.yaml"))
	require.NoError(t, err, "failed to load test config")

	return cfg
}

// setupTestDB ініціалізує БД та повертає функцію для очищення таблиць
func setupTestDB(t *testing.T, tables ...string) *sqlx.DB {
	cfg := initTestConfig(t)

	db, err := NewPostgresDB(config.DBConfig(cfg.TestDatabase))
	require.NoError(t, err, "failed to connect to test db")

	// Очищуємо передані таблиці для ізоляції тесту
	if len(tables) > 0 {
		for _, table := range tables {
			_, err := db.Exec("TRUNCATE TABLE " + table + " CASCADE")
			require.NoError(t, err, "failed to clear the table: "+table)
		}
	}

	return db
}
