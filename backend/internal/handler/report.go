package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getFieldReport(c *gin.Context) {
	fieldID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid field id")
		return
	}

	fromRaw := c.Query("from")
	toRaw := c.Query("to")
	if fromRaw == "" || toRaw == "" {
		h.newErrorResponse(c, http.StatusBadRequest, "from and to query parameters are required in RFC3339 format")
		return
	}

	from, err := time.Parse(time.RFC3339, fromRaw)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid from date format, use RFC3339")
		return
	}

	to, err := time.Parse(time.RFC3339, toRaw)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid to date format, use RFC3339")
		return
	}

	report, err := h.services.Report.GenerateFieldReport(fieldID, from, to)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, report)
}
