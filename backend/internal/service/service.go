package service

import (
	"agro-data-management-system/internal/repository"

	"go.uber.org/zap"
)

// Services — це головна структура, яка містить усі інтерфейси сервісів
type Services struct {
	Crop     CropService
	Field    FieldService
	Sensor   SensorService
	Metric   MetricService
	Pest     PestService
	Forecast ForecastService
}

// Dependencies — допоміжна структура для ініціалізації
type Dependencies struct {
	Repos *repository.Repositories
	Log   *zap.Logger
	AiURL string
}

// NewServices ініціалізує всі сервіси та повертає їх як єдиний об'єкт
func NewServices(deps Dependencies) *Services {
	return &Services{
		Crop:     NewCropService(deps.Repos.Crop, deps.Log),
		Field:    NewFieldService(deps.Repos.Field, deps.Repos.Crop, deps.Log),
		Sensor:   NewSensorService(deps.Repos.Sensor, deps.Repos.Field, deps.Log),
		Metric:   NewMetricService(deps.Repos.Metric, deps.Repos.Sensor, deps.Log),
		Pest:     NewPestService(deps.Repos.Pest, deps.Log),
		Forecast: NewForecastService(deps.Repos.Forecast, nil, nil, nil, deps.AiURL, deps.Log),
		// Примітка: для ForecastService ми передамо інші сервіси пізніше
		// або оновимо його конструктор, щоб він приймав цілісний об'єкт
	}
}
