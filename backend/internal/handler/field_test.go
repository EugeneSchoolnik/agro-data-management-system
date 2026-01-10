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

func TestHandler_createField(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	t.Run("Success", func(t *testing.T) {
		mockFieldService := new(serviceMocks.FieldService)
		services := &service.Services{Field: mockFieldService}
		h := NewHandler(services, log)

		r := gin.New()
		r.POST("/fields", h.createField)

		input := models.Field{
			Name:     "Північне поле",
			Area:     10.5,
			Location: "50.1, 30.2",
			CropID:   1,
		}

		mockFieldService.On("Create", mock.MatchedBy(func(f models.Field) bool {
			return f.Name == input.Name && f.CropID == input.CropID
		})).Return(1, nil).Once()

		jsonInput, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/fields", bytes.NewBuffer(jsonInput))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"id":1`)
	})

	t.Run("Crop_Not_Found_Error", func(t *testing.T) {
		mockFieldService := new(serviceMocks.FieldService)
		services := &service.Services{Field: mockFieldService}
		h := NewHandler(services, log)

		r := gin.New()
		r.POST("/fields", h.createField)

		input := models.Field{Name: "Поле X", Area: 5, Location: "0,0", CropID: 999}

		// Емулюємо помилку бізнес-логіки (культури не існує)
		mockFieldService.On("Create", mock.Anything).
			Return(0, errors.New("cannot create field: crop with id 999 not found")).Once()

		jsonInput, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/fields", bytes.NewBuffer(jsonInput))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "crop with id 999 not found")
	})
}

func TestHandler_getAllFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := zap.NewNop()

	t.Run("Success", func(t *testing.T) {
		mockFieldService := new(serviceMocks.FieldService)
		services := &service.Services{Field: mockFieldService}
		h := NewHandler(services, log)

		r := gin.New()
		r.GET("/fields", h.getAllFields)

		// Очікуємо дані зі зв'язаною культурою
		mockData := []models.FieldWithCrop{
			{
				Field:    models.Field{ID: 1, Name: "Поле 1"},
				CropName: "Кукурудза",
			},
		}
		mockFieldService.On("GetAll").Return(mockData, nil).Once()

		req, _ := http.NewRequest("GET", "/fields", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"crop_name":"Кукурудза"`)
	})
}
