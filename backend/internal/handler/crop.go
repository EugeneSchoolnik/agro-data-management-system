package handler

import (
	"agro-data-management-system/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// createCrop — Створення нової культури
func (h *Handler) createCrop(c *gin.Context) {
	var input models.Crop
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid input data")
		return
	}

	crop, err := h.services.Crop.Create(input)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, crop)
}

// getAllCrops — Отримання списку всіх культур
func (h *Handler) getAllCrops(c *gin.Context) {
	crops, err := h.services.Crop.GetAll()
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, crops)
}

// getCropById — Отримання однієї культури за ID
func (h *Handler) getCropById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}

	crop, err := h.services.Crop.GetByID(id)
	if err != nil {
		h.newErrorResponse(c, http.StatusNotFound, "crop not found")
		return
	}

	h.newSuccessResponse(c, crop)
}

// updateCrop — Оновлення даних культури
func (h *Handler) updateCrop(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var input models.Crop
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.ID = id

	if err := h.services.Crop.Update(input); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, "updated")
}

// deleteCrop — Видалення культури
func (h *Handler) deleteCrop(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.services.Crop.Delete(id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, "deleted")
}
