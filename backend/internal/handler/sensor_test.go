package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestHandler_registerSensor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	t.Run("Success", func(t *testing.T) {
		mockSrv := new(serviceMocks.SensorService)
		h := NewHandler(&service.Services{Sensor: mockSrv}, log)
		r := gin.New()
		r.POST("/sensors", h.registerSensor)

		input := models.Sensor{
			FieldID:    1,
			SensorType: "DHT22",
			Status:     models.StatusActive,
		}

		mockSrv.On("Register", mock.MatchedBy(func(s models.Sensor) bool {
			return s.FieldID == 1 && s.SensorType == "DHT22"
		})).Return(101, nil).Once()

		jsonInput, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/sensors", bytes.NewBuffer(jsonInput))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"id":101`)
	})

	t.Run("Field_Not_Found", func(t *testing.T) {
		mockSrv := new(serviceMocks.SensorService)
		h := NewHandler(&service.Services{Sensor: mockSrv}, log)
		r := gin.New()
		r.POST("/sensors", h.registerSensor)

		mockSrv.On("Register", mock.Anything).
			Return(0, errors.New("field 404 not found")).Once()

		body, _ := json.Marshal(models.Sensor{FieldID: 404, SensorType: "DHT22", Status: models.StatusActive})
		req, _ := http.NewRequest("POST", "/sensors", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "field 404 not found")
	})
}

func TestHandler_updateSensorStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	t.Run("Success_Status_Update", func(t *testing.T) {
		mockSrv := new(serviceMocks.SensorService)
		h := NewHandler(&service.Services{Sensor: mockSrv}, log)
		r := gin.New()
		r.PATCH("/sensors/:id/status", h.updateSensorStatus)

		sensorID := 1
		newStatus := models.StatusError

		mockSrv.On("UpdateStatus", sensorID, newStatus).Return(nil).Once()

		jsonInput, _ := json.Marshal(map[string]string{"status": string(newStatus)})
		req, _ := http.NewRequest("PATCH", "/sensors/1/status", bytes.NewBuffer(jsonInput))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "status updated")
	})
}

func TestHandler_getSensorsByField(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	t.Run("Success", func(t *testing.T) {
		mockSrv := new(serviceMocks.SensorService)
		h := NewHandler(&service.Services{Sensor: mockSrv}, log)
		r := gin.New()
		r.GET("/fields/:id/sensors", h.getSensorsByField)

		fieldID := 1
		mockSensors := []models.Sensor{
			{ID: 1, FieldID: fieldID, SensorType: "Moisture"},
			{ID: 2, FieldID: fieldID, SensorType: "Temp"},
		}

		mockSrv.On("GetByField", fieldID).Return(mockSensors, nil).Once()

		req, _ := http.NewRequest("GET", "/fields/1/sensors", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Перевіряємо, що в масиві 2 датчики
		var resp map[string][]models.Sensor
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp["data"], 2)
	})
}
