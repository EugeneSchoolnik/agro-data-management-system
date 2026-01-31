package main

import (
	"agro-data-management-system/internal/config"
	"agro-data-management-system/internal/handler"
	"agro-data-management-system/internal/repository"
	"agro-data-management-system/internal/service"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// Ініціалізація логера (можна передати рівень з конфігу)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	cfg, err := config.LoadConfig("config/local.yaml")
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatal("Database connection error", zap.Error(err))
	}
	defer db.Close()

	repos := repository.NewRepositories(db)

	deps := service.Dependencies{
		Repos: repos,
		Log:   logger,
		AiURL: cfg.Forecasting.PythonServiceURL, // майбутній Python сервер
	}
	services := service.NewServices(deps)

	handlers := handler.NewHandler(services, logger)

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        handlers.InitRoutes(),
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	logger.Info("Agro-system backend started successfully",
		zap.String("port", cfg.Server.Port),
		zap.String("db_name", cfg.Database.Name),
	)

	go func() {
		logger.Info("Starting server on port 8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Listen error", zap.Error(err))
		}
	}()

	// 6. Graceful Shutdown (очікування сигналу вимкнення)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
