package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server      HTTPServer `yaml:"server"`
	Database    DBConfig   `yaml:"database"`
	Forecasting AIConfig   `yaml:"forecasting"`
}

type HTTPServer struct {
	Port    string        `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"10s"`
}

type DBConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-required:"true"`
	Password string `yaml:"pass" env:"DB_PASSWORD" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
	SSLMode string `yaml:"ssl_mode" env-default:"disable"`
}

type AIConfig struct {
	PythonServiceURL string `yaml:"python_service_url"`
	APIKey           string `yaml:"api_key" env:"AI_API_KEY"`
}

func LoadConfig(path string) (*Config, error) {
	// 1. Завантажуємо .env у змінні оточення (якщо файл існує)
	if err := godotenv.Load(); err != nil {
		// Не фатально, бо на сервері змінні можуть бути прописані в ОС
		log.Println("Info: .env file not found, using system environment variables")
	}

	var cfg Config

	// 2. Читаємо YAML, але теги `env` автоматично перекриють значення
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}