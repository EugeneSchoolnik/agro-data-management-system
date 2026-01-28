package handler

import (
	"agro-data-management-system/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// createField — створення нового поля
func (h *Handler) createField(c *gin.Context) {
	var input models.Field
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid input data")
		return
	}

	field, err := h.services.Field.Create(input)
	if err != nil {
		// Якщо сервіс повернув помилку "crop not found", це 400 Bad Request
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.newSuccessResponse(c, field)
}

// getAllFields — список усіх полів з назвами культур
func (h *Handler) getAllFields(c *gin.Context) {
	fields, err := h.services.Field.GetAll()
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, fields)
}

// getFieldById — детальне інфо про поле
func (h *Handler) getFieldById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}

	field, err := h.services.Field.GetByID(id)
	if err != nil {
		h.newErrorResponse(c, http.StatusNotFound, "field not found")
		return
	}

	h.newSuccessResponse(c, field)
}

// updateField — оновлення даних поля
func (h *Handler) updateField(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var input models.Field
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.ID = id

	if err := h.services.Field.Update(input); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, "updated")
}

// deleteField — видалення поля
func (h *Handler) deleteField(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.services.Field.Delete(id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, "deleted")
}
