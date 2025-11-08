package main

import (
	"agro-data-management-system/internal/config"
	"agro-data-management-system/internal/logger"
	"agro-data-management-system/internal/repository"

	"go.uber.org/zap"
)

func main() {
	// 1. Ініціалізація логера (можна передати рівень з конфігу)
	logger.Init("info")
	defer logger.Log.Sync() // Очищення буфера перед виходом

	cfg, err := config.LoadConfig("config/local.yaml")
	if err != nil {
		logger.Log.Fatal("Failed to load config", zap.Error(err))
	}

	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Log.Fatal("Database connection error", zap.Error(err))
	}
	_ = db

	logger.Log.Info("Agro-system backend started successfully",
		zap.String("port", cfg.Server.Port),
		zap.String("db_name", cfg.Database.Name),
	)
}