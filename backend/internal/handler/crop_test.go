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
	"go.uber.org/zap"
)

func TestHandler_createCrop(t *testing.T) {
	// 1. Налаштування
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	t.Run("Success", func(t *testing.T) {
		mockCropService := new(serviceMocks.CropService)
		services := &service.Services{Crop: mockCropService}
		h := NewHandler(services, log)

		// Налаштовуємо роутер для тесту
		r := gin.New()
		r.POST("/crops", h.createCrop)

		// Вхідні дані
		input := models.Crop{Name: "Пшениця", Variety: "Скарбниця"}
		mockCropService.On("Create", input).Return(1, nil).Once()

		// Формуємо запит
		jsonInput, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/crops", bytes.NewBuffer(jsonInput))
		w := httptest.NewRecorder()

		// Виконуємо запит
		r.ServeHTTP(w, req)

		// Перевірки
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"id":1`)
		mockCropService.AssertExpectations(t)
	})

	t.Run("Invalid_Input", func(t *testing.T) {
		services := &service.Services{} // сервіс не буде викликаний
		h := NewHandler(services, log)
		r := gin.New()
		r.POST("/crops", h.createCrop)

		// Відправляємо битий JSON
		req, _ := http.NewRequest("POST", "/crops", bytes.NewBufferString("{invalid}"))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid input data")
	})
}

func TestHandler_getCropById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	t.Run("Found", func(t *testing.T) {
		mockCropService := new(serviceMocks.CropService)
		services := &service.Services{Crop: mockCropService}
		h := NewHandler(services, log)

		r := gin.New()
		r.GET("/crops/:id", h.getCropById)

		expectedCrop := models.Crop{ID: 1, Name: "Соняшник"}
		mockCropService.On("GetByID", 1).Return(expectedCrop, nil).Once()

		req, _ := http.NewRequest("GET", "/crops/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"Соняшник"`)
	})

	t.Run("Not_Found", func(t *testing.T) {
		mockCropService := new(serviceMocks.CropService)
		services := &service.Services{Crop: mockCropService}
		h := NewHandler(services, log)

		r := gin.New()
		r.GET("/crops/:id", h.getCropById)

		mockCropService.On("GetByID", 404).Return(models.Crop{}, errors.New("not found")).Once()

		req, _ := http.NewRequest("GET", "/crops/404", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "crop not found")
	})
}
