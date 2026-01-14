package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/service"
	serviceMocks "agro-data-management-system/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestHandler_saveMetric(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	t.Run("Success", func(t *testing.T) {
		mockSrv := new(serviceMocks.MetricService)
		h := NewHandler(&service.Services{Metric: mockSrv}, log)
		r := gin.New()
		r.POST("/metrics", h.saveMetric)

		input := models.Metric{
			SensorID: 1,
			Value:    24.5,
		}

		mockSrv.On("Save", mock.MatchedBy(func(m models.Metric) bool {
			return m.SensorID == 1 && m.Value == 24.5
		})).Return(int64(1), nil).Once()

		jsonInput, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/metrics", bytes.NewBuffer(jsonInput))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"id":1`)
	})

	t.Run("Inactive_Sensor_Error", func(t *testing.T) {
		mockSrv := new(serviceMocks.MetricService)
		h := NewHandler(&service.Services{Metric: mockSrv}, log)
		r := gin.New()
		r.POST("/metrics", h.saveMetric)

		// Емулюємо помилку від сервісу (наприклад, датчик вимкнено)
		mockSrv.On("Save", mock.Anything).
			Return(int64(0), assert.AnError).Once()

		body, _ := json.Marshal(models.Metric{SensorID: 2, Value: 10.0})
		req, _ := http.NewRequest("POST", "/metrics", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_getMetricHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	t.Run("Valid_History_Request", func(t *testing.T) {
		mockSrv := new(serviceMocks.MetricService)
		h := NewHandler(&service.Services{Metric: mockSrv}, log)
		r := gin.New()
		r.GET("/sensors/:id/metrics/history", h.getMetricHistory)

		sensorID := 1
		from := time.Now().Add(-1 * time.Hour).Truncate(time.Second).UTC()
		to := time.Now().Truncate(time.Second).UTC()

		mockHistory := []models.Metric{
			{ID: 1, SensorID: sensorID, Value: 20.1, RecordedAt: from},
		}

		mockSrv.On("GetHistory", sensorID, mock.Anything, mock.Anything).Return(mockHistory, nil).Once()

		// Формуємо URL з query-параметрами у форматі RFC3339
		url := "/sensors/1/metrics/history?from=" + from.Format(time.RFC3339) + "&to=" + to.Format(time.RFC3339)
		req, _ := http.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"value":20.1`)
	})

	t.Run("Invalid_Time_Format", func(t *testing.T) {
		h := NewHandler(&service.Services{}, log)
		r := gin.New()
		r.GET("/sensors/:id/metrics/history", h.getMetricHistory)

		req, _ := http.NewRequest("GET", "/sensors/1/metrics/history?from=2026-01-01&to=now", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid time format")
	})
}
