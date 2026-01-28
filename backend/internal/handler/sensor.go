package handler

import (
	"agro-data-management-system/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) registerSensor(c *gin.Context) {
	var input models.Sensor
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}

	sensor, err := h.services.Sensor.Register(input)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.newSuccessResponse(c, sensor)
}

func (h *Handler) getSensorById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	sensor, err := h.services.Sensor.GetByID(id)
	if err != nil {
		h.newErrorResponse(c, http.StatusNotFound, "sensor not found")
		return
	}
	h.newSuccessResponse(c, sensor)
}

func (h *Handler) getSensorsByField(c *gin.Context) {
	fieldID, _ := strconv.Atoi(c.Param("id"))
	sensors, err := h.services.Sensor.GetByField(fieldID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.newSuccessResponse(c, sensors)
}

func (h *Handler) updateSensorStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Status models.SensorStatus `json:"status" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "status is required")
		return
	}

	if err := h.services.Sensor.UpdateStatus(id, input.Status); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, "status updated")
}

func (h *Handler) deleteSensor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.services.Sensor.Delete(id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.newSuccessResponse(c, "deleted")
}
