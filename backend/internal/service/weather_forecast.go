package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"
	"go.uber.org/zap"
)

type WeatherForecastService interface {
	PredictWeather(stationID int, hoursAhead int) (WeatherForecastResult, error)
}

type WeatherForecastResult struct {
	Temperature    float64 `json:"temperature"`
	HoursAhead     int     `json:"hours_ahead"`
	Recommendation string  `json:"recommendation"`
}

type weatherForecastRequest struct {
	WeatherData []float64 `json:"weather_data"`
	HoursAhead  int       `json:"hours_ahead"`
}

type weatherForecastService struct {
	repo       repository.WeatherRepository
	aiURL      string // базовий URL forecast-engine
	log        *zap.Logger
	httpClient *http.Client
}

func NewWeatherForecastService(
	repo repository.WeatherRepository,
	aiURL string,
	log *zap.Logger,
) WeatherForecastService {
	return &weatherForecastService{
		repo:       repo,
		aiURL:      aiURL,
		log:        log,
		httpClient: &http.Client{},
	}
}

// PredictWeather — прогнозує температуру на основі останніх спостережень
// Параметри для прогнозу (у порядку):
// 1. Температура повітря (param_id=1)
// 2. Атмосферний тиск (param_id=2)
// 3. Швидкість вітру (param_id=3)
// 4. Точка роси (param_id=5)
// 5. Опади накопичені (param_id=8)
// 6. Відносна вологість (param_id=9)
// 7. Температура ґрунту (param_id=10)
// 8. Сонячна радіація (param_id=20)
func (s *weatherForecastService) PredictWeather(stationID int, hoursAhead int) (WeatherForecastResult, error) {
	result := WeatherForecastResult{}

	if hoursAhead < 1 || hoursAhead > 24 {
		return result, errors.New("hours_ahead must be between 1 and 24")
	}

	// Отримуємо потрібні параметри з БД
	requiredParamIDs := []int{1, 2, 3, 5, 8, 9, 10, 20}
	observations, err := s.repo.GetObservationsByParameterIDs(stationID, requiredParamIDs)
	if err != nil {
		s.log.Error("Failed to fetch weather observations", zap.Error(err))
		return result, err
	}

	if len(observations) == 0 {
		return result, errors.New("no weather observations found for station")
	}

	// Форматуємо дані для forecast-engine
	weatherData, err := s.formatWeatherData(observations)
	if err != nil {
		s.log.Error("Failed to format weather data", zap.Error(err))
		return result, err
	}

	if len(weatherData) != 8 {
		return result, fmt.Errorf("expected 8 weather parameters, got %d", len(weatherData))
	}

	// Відправляємо запит до forecast-engine
	forecastResult, err := s.callForecastEngine(weatherData, hoursAhead)
	if err != nil {
		s.log.Error("Failed to call forecast engine", zap.Error(err))
		return result, err
	}

	result.Temperature = forecastResult.Temperature
	result.HoursAhead = hoursAhead
	result.Recommendation = forecastResult.Recommendation

	s.log.Info(
		"Weather forecast computed",
		zap.Int("station_id", stationID),
		zap.Int("hours_ahead", hoursAhead),
		zap.Float64("temperature", forecastResult.Temperature),
	)

	return result, nil
}

// formatWeatherData — форматує спостереження в послідовність параметрів для прогнозу
// Порядок: [temp, pressure, wind_speed, dew_point, precipitation, humidity, soil_temp, solar_radiation]
func (s *weatherForecastService) formatWeatherData(observations []models.WeatherObservation) ([]float64, error) {
	// Карта для швидкого доступу до значень за param_id
	paramMap := make(map[int]float64)

	for _, obs := range observations {
		param, err := s.repo.GetParameterByID(obs.WeatherParameterID)
		if err != nil {
			s.log.Warn("Unknown weather parameter", zap.Int("param_id", obs.WeatherParameterID))
			continue
		}

		// Зберігаємо останнє значення для кожного параметра
		paramMap[param.ParamID] = obs.Value
	}

	// Порядок параметрів для forecast-engine
	paramOrder := []int{1, 2, 3, 5, 8, 9, 10, 20}
	weatherData := make([]float64, 0, len(paramOrder))

	for _, paramID := range paramOrder {
		value, exists := paramMap[paramID]
		if !exists {
			// Якщо параметр не знайдено, встановлюємо 0 (або можна повернути помилку)
			s.log.Warn("Missing weather parameter", zap.Int("param_id", paramID))
			value = 0.0
		}
		weatherData = append(weatherData, value)
	}

	return weatherData, nil
}

// callForecastEngine — відправляє запит до forecast-engine API
func (s *weatherForecastService) callForecastEngine(weatherData []float64, hoursAhead int) (WeatherForecastResult, error) {
	result := WeatherForecastResult{}

	reqBody := weatherForecastRequest{
		WeatherData: weatherData,
		HoursAhead:  hoursAhead,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return result, err
	}

	url := s.aiURL + "/predict-weather"
	resp, err := s.httpClient.Post(url, "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		return result, fmt.Errorf("failed to connect to forecast engine: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return result, fmt.Errorf("forecast engine error (status %d): %s", resp.StatusCode, string(body))
	}

	var forecastResp WeatherForecastResult
	if err := json.NewDecoder(resp.Body).Decode(&forecastResp); err != nil {
		return result, fmt.Errorf("failed to parse forecast response: %w", err)
	}

	return forecastResp, nil
}
