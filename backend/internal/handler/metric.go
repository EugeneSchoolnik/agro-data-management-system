package handler

import (
	"agro-data-management-system/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// saveMetric — приймає дані від датчика
func (h *Handler) saveMetric(c *gin.Context) {
	var input models.Metric
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid metric data")
		return
	}

	// Встановлюємо час запису, якщо він не прийшов від датчика
	if input.RecordedAt.IsZero() {
		input.RecordedAt = time.Now().UTC()
	}

	metric, err := h.services.Metric.Save(input)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.newSuccessResponse(c, metric)
}

// getLatestMetric — останнє значення конкретного датчика
func (h *Handler) getLatestMetric(c *gin.Context) {
	sensorID, _ := strconv.Atoi(c.Param("id"))

	metric, err := h.services.Metric.GetLatest(sensorID)
	if err != nil {
		h.newErrorResponse(c, http.StatusNotFound, "no metrics found for this sensor")
		return
	}

	h.newSuccessResponse(c, metric)
}

// getMetricHistory — дані за період для графіків
func (h *Handler) getMetricHistory(c *gin.Context) {
	sensorID, _ := strconv.Atoi(c.Param("id"))

	// Очікуємо формат RFC3339: 2026-04-11T15:00:00Z
	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, errFrom := time.Parse(time.RFC3339, fromStr)
	to, errTo := time.Parse(time.RFC3339, toStr)

	if errFrom != nil || errTo != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid time format. Use RFC3339")
		return
	}

	history, err := h.services.Metric.GetHistory(sensorID, from, to)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, history)
}
