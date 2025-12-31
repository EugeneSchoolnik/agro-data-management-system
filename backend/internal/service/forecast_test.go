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

		// Повертаємо фейкову відповідь від "нейромережі"
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
	pestSrv := new(serviceMocks.PestService)
	log := zap.NewNop()

	// Передаємо URL нашого фейкового сервера як aiURL
	srv := NewForecastService(fRepo, mSrv, fieldSrv, pestSrv, mockAI.URL, log)

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

		mSrv.On("GetHistory", fieldID, mock.Anything, mock.Anything).Return([]models.Metric{
			{Value: 22.5}, {Value: 23.0},
		}, nil).Once()

		// Очікуємо, що сервіс збереже результат у БД
		fRepo.On("Create", mock.MatchedBy(func(f models.Forecast) bool {
			return f.Probability == 0.85 && f.FieldID == fieldID
		})).Return(100, nil).Once()

		// 4. Викликаємо метод, який тестуємо
		result, err := srv.Predict(fieldID, pestID)

		// 5. Перевіряємо результати
		assert.NoError(t, err)
		assert.Equal(t, 100, result.ID)
		assert.Equal(t, 0.85, result.Probability)
		assert.Contains(t, result.Recommendation, "високий ризик")

		// Перевіряємо, що всі моки були викликані
		mock.AssertExpectationsForObjects(t, fRepo, mSrv, fieldSrv, pestSrv)
	})
}
