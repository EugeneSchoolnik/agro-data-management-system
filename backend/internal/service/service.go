package service

import (
	"agro-data-management-system/internal/repository"
	"agro-data-management-system/internal/weatherapi"

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
	Weather  WeatherService
	Report   ReportService
}

// Dependencies — допоміжна структура для ініціалізації
type Dependencies struct {
	Repos              *repository.Repositories
	Log                *zap.Logger
	AiURL              string
	WeatherAPIURL      string
	WeatherAPILogin    string
	WeatherAPIPassword string
}

// NewServices ініціалізує всі сервіси та повертає їх як єдиний об'єкт
func NewServices(deps Dependencies) *Services {
	// Спершу створюємо незалежні сервіси
	cropSrv := NewCropService(deps.Repos.Crop, deps.Log)
	fieldSrv := NewFieldService(deps.Repos.Field, deps.Repos.Crop, deps.Log)
	sensorSrv := NewSensorService(deps.Repos.Sensor, deps.Repos.Field, deps.Log)
	pestSrv := NewPestService(deps.Repos.Pest, deps.Log)
	metricSrv := NewMetricService(deps.Repos.Metric, deps.Repos.Sensor, deps.Log)
	weatherAPIClient := weatherapi.NewClient(deps.WeatherAPIURL, deps.WeatherAPILogin, deps.WeatherAPIPassword)
	weatherSrv := NewWeatherService(deps.Repos.Weather, weatherAPIClient, deps.Log)

	reportSrv := NewReportService(deps.Repos.Field, deps.Repos.Metric, deps.Repos.Forecast, deps.Log)

	forecastSrv := NewForecastService(
		deps.Repos.Forecast,
		metricSrv,
		fieldSrv,
		sensorSrv,
		pestSrv,
		deps.AiURL,
		deps.Log,
	)

	return &Services{
		Crop:     cropSrv,
		Field:    fieldSrv,
		Sensor:   sensorSrv,
		Metric:   metricSrv,
		Pest:     pestSrv,
		Forecast: forecastSrv,
		Weather:  weatherSrv,
		Report:   reportSrv,
	}
}
