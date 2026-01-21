package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// predict — ініціація розрахунку ШІ
func (h *Handler) predict(c *gin.Context) {
	var input struct {
		FieldID int `json:"field_id" binding:"required"`
		PestID  int `json:"pest_id" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "field_id and pest_id are required")
		return
	}

	forecast, err := h.services.Forecast.Predict(input.FieldID, input.PestID)
	if err != nil {
		// Тут може бути помилка відсутності даних або недоступності AI Engine
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, forecast)
}

// getLatestForecast — отримати останній збережений прогноз для поля
func (h *Handler) getLatestForecast(c *gin.Context) {
	fieldID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid field id")
		return
	}

	forecast, err := h.services.Forecast.GetLatest(fieldID)
	if err != nil {
		h.newErrorResponse(c, http.StatusNotFound, "no forecasts found for this field")
		return
	}

	h.newSuccessResponse(c, forecast)
}
