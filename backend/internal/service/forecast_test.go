package service

import (
	"agro-data-management-system/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// Імпортуємо пакет з моками сервісів
	repoMocks "agro-data-management-system/internal/repository/mocks"
	serviceMocks "agro-data-management-system/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestForecastService_Predict(t *testing.T) {
	// 1. Створюємо фіктивний ШІ-сервер (Mock AI Engine)
	mockAI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Перевіряємо, чи правильний шлях запиту
		assert.Equal(t, "/predict", r.URL.Path)

		var req models.ForecastRequest
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		assert.Equal(t, "Пшениця", req.CropName)
		assert.Equal(t, "Скарбниця", req.Variety)
		assert.Equal(t, "Eurygaster integriceps", req.PestName)
		assert.Equal(t, []float64{22.0, 60.0, 30.0, 3.0}, req.Metrics)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ForecastResponse{
			Probability:    0.85,
			Recommendation: "Термінова обробка: високий ризик E. integriceps",
		})
	}))
	defer mockAI.Close()

	// 2. Ініціалізуємо моки залежностей
	fRepo := new(repoMocks.ForecastRepository)
	mSrv := new(serviceMocks.MetricService)
	fieldSrv := new(serviceMocks.FieldService)
	sensorSrv := new(serviceMocks.SensorService)
	pestSrv := new(serviceMocks.PestService)
	log := zap.NewNop()

	// Передаємо URL нашого фейкового сервера як aiURL
	srv := NewForecastService(fRepo, mSrv, fieldSrv, sensorSrv, pestSrv, mockAI.URL, log)

	t.Run("Full_Success_Flow", func(t *testing.T) {
		fieldID, pestID := 1, 2

		// 3. Налаштовуємо очікування для сервісів (Mock Expectations)
		fieldSrv.On("GetByID", fieldID).Return(models.FieldWithCrop{
			Field:       models.Field{ID: fieldID},
			CropName:    "Пшениця",
			CropVariety: "Скарбниця",
		}, nil).Once()

		pestSrv.On("GetByID", pestID).Return(models.Pest{
			ID:             pestID,
			ScientificName: "Eurygaster integriceps",
		}, nil).Once()

		sensorSrv.On("GetByField", fieldID).Return([]models.Sensor{
			{ID: 10, SensorType: "temperature"},
			{ID: 11, SensorType: "air_humidity"},
			{ID: 12, SensorType: "soil_moisture"},
			{ID: 13, SensorType: "crop_phase"},
		}, nil).Once()

		mSrv.On("GetHistory", 10, mock.Anything, mock.Anything).Return([]models.Metric{
			{Value: 20.0},
			{Value: 24.0},
		}, nil).Once()
		mSrv.On("GetHistory", 11, mock.Anything, mock.Anything).Return([]models.Metric{{Value: 60.0}}, nil).Once()
		mSrv.On("GetHistory", 12, mock.Anything, mock.Anything).Return([]models.Metric{{Value: 30.0}}, nil).Once()
		mSrv.On("GetHistory", 13, mock.Anything, mock.Anything).Return([]models.Metric{{Value: 3.0}}, nil).Once()

		// Очікуємо, що сервіс збереже результат у БД
		fRepo.On("Create", mock.MatchedBy(func(f models.Forecast) bool {
			return f.Probability == 0.85 && f.FieldID == fieldID
		})).Return(models.Forecast{ID: 100, Probability: 0.85, FieldID: fieldID, Recommendation: "Термінова обробка: високий ризик E. integriceps"}, nil).Once()

		// 4. Викликаємо метод, який тестуємо
		result, err := srv.Predict(fieldID, pestID)

		// 5. Перевіряємо результати
		assert.NoError(t, err)
		assert.Equal(t, 100, result.ID)
		assert.Equal(t, 0.85, result.Probability)
		assert.Contains(t, result.Recommendation, "високий ризик")

		// Перевіряємо, що всі моки були викликані
		mock.AssertExpectationsForObjects(t, fRepo, mSrv, fieldSrv, sensorSrv, pestSrv)
	})
}
