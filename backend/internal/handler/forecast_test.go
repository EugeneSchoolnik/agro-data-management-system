package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/service"
	serviceMocks "agro-data-management-system/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestHandler_predict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	t.Run("Success_AI_Prediction", func(t *testing.T) {
		mockSrv := new(serviceMocks.ForecastService)
		h := NewHandler(&service.Services{Forecast: mockSrv}, log)
		r := gin.New()
		r.POST("/predict", h.predict)

		expectedForecast := models.Forecast{
			ID:             1,
			Probability:    0.92,
			Recommendation: "Високий ризик. Рекомендована обробка.",
		}

		mockSrv.On("Predict", 1, 2).Return(expectedForecast, nil).Once()

		body, _ := json.Marshal(map[string]int{"field_id": 1, "pest_id": 2})
		req, _ := http.NewRequest("POST", "/predict", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"probability":0.92`)
	})

	t.Run("AI_Engine_Error", func(t *testing.T) {
		mockSrv := new(serviceMocks.ForecastService)
		h := NewHandler(&service.Services{Forecast: mockSrv}, log)
		r := gin.New()
		r.POST("/predict", h.predict)

		mockSrv.On("Predict", mock.Anything, mock.Anything).
			Return(models.Forecast{}, assert.AnError).Once()

		body, _ := json.Marshal(map[string]int{"field_id": 1, "pest_id": 2})
		req, _ := http.NewRequest("POST", "/predict", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
