package handler

import (
	"agro-data-management-system/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createPest(c *gin.Context) {
	var input models.Pest
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid input data")
		return
	}

	id, err := h.services.Pest.Create(input)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, gin.H{"id": id})
}

func (h *Handler) getAllPests(c *gin.Context) {
	pests, err := h.services.Pest.GetAll()
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, pests)
}

func (h *Handler) getPestById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	pest, err := h.services.Pest.GetByID(id)
	if err != nil {
		h.newErrorResponse(c, http.StatusNotFound, "pest not found")
		return
	}

	h.newSuccessResponse(c, pest)
}

func (h *Handler) updatePest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input models.Pest
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.ID = id

	if err := h.services.Pest.Update(input); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, "updated")
}

func (h *Handler) deletePest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.services.Pest.Delete(id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, "deleted")
}
