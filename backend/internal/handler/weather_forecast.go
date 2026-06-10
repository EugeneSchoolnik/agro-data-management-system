package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type weatherForecastRequest struct {
	StationID  int `json:"station_id" binding:"required,min=1"`
	HoursAhead int `json:"hours_ahead" binding:"required,min=1,max=24"`
}

// predictWeatherForecast — прогноз погоди на основі останніх спостережень зі станції
// POST /api/v1/weather/forecast/predict
func (h *Handler) predictWeatherForecast(c *gin.Context) {
	var input weatherForecastRequest

	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "station_id and hours_ahead are required (1-24 hours)")
		return
	}

	forecast, err := h.services.WeatherForecast.PredictWeather(input.StationID, input.HoursAhead)
	if err != nil {
		h.log.Warn("Weather forecast error", zap.Error(err), zap.Int("station_id", input.StationID))
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.newSuccessResponse(c, forecast)
}
