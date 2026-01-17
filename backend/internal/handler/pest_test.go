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

// Допоміжна функція для ініціалізації тестового оточення
func setupPestTest() (*gin.Engine, *serviceMocks.PestService, *Handler) {
	gin.SetMode(gin.TestMode)
	mockSrv := new(serviceMocks.PestService)
	h := NewHandler(&service.Services{Pest: mockSrv}, zap.NewNop())

	r := gin.New()
	// Маршрути мають збігатися з InitRoutes
	api := r.Group("/api/v1/pests")
	{
		api.POST("/", h.createPest)
		api.GET("/", h.getAllPests)
		api.GET("/:id", h.getPestById)
		api.PUT("/:id", h.updatePest)
		api.DELETE("/:id", h.deletePest)
	}
	return r, mockSrv, h
}

func TestPestHandler_Create(t *testing.T) {
	r, mockSrv, _ := setupPestTest()

	t.Run("Success", func(t *testing.T) {
		input := models.Pest{Name: "Попелиця", ScientificName: "Aphidoidea"}
		mockSrv.On("Create", input).Return(1, nil).Once()

		body, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/api/v1/pests/", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"id":1`)
	})

	t.Run("Invalid_JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/pests/", bytes.NewBufferString("{invalid}"))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid input data")
	})
}

func TestPestHandler_GetAll(t *testing.T) {
	r, mockSrv, _ := setupPestTest()

	t.Run("Success", func(t *testing.T) {
		pests := []models.Pest{
			{ID: 1, Name: "Pest 1"},
			{ID: 2, Name: "Pest 2"},
		}
		mockSrv.On("GetAll").Return(pests, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/pests/", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Pest 1")
		assert.Contains(t, w.Body.String(), "Pest 2")
	})
}

func TestPestHandler_GetByID(t *testing.T) {
	r, mockSrv, _ := setupPestTest()

	t.Run("Found", func(t *testing.T) {
		pest := models.Pest{ID: 1, Name: "Жук"}
		mockSrv.On("GetByID", 1).Return(pest, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/pests/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Жук")
	})

	t.Run("Not_Found", func(t *testing.T) {
		mockSrv.On("GetByID", 404).Return(models.Pest{}, errors.New("not found")).Once()

		req, _ := http.NewRequest("GET", "/api/v1/pests/404", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestPestHandler_Update(t *testing.T) {
	r, mockSrv, _ := setupPestTest()

	t.Run("Success", func(t *testing.T) {
		input := models.Pest{Name: "Оновлена назва", ScientificName: "New scientific"}
		// Очікуємо, що в сервіс прийде об'єкт з ID=1
		expectedInput := input
		expectedInput.ID = 1

		mockSrv.On("Update", expectedInput).Return(nil).Once()

		body, _ := json.Marshal(input)
		req, _ := http.NewRequest("PUT", "/api/v1/pests/1", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "updated")
	})
}

func TestPestHandler_Delete(t *testing.T) {
	r, mockSrv, _ := setupPestTest()

	t.Run("Success", func(t *testing.T) {
		mockSrv.On("Delete", 1).Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/api/v1/pests/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "deleted")
	})

	t.Run("Server_Error", func(t *testing.T) {
		mockSrv.On("Delete", 500).Return(errors.New("db error")).Once()

		req, _ := http.NewRequest("DELETE", "/api/v1/pests/500", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
