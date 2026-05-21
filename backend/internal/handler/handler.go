package handler

import (
	"agro-data-management-system/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	services *service.Services
	log      *zap.Logger
}

func NewHandler(services *service.Services, log *zap.Logger) *Handler {
	return &Handler{
		services: services,
		log:      log,
	}
}

// InitRoutes налаштовує всі шляхи API
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Додаємо стандартні Middleware
	router.Use(gin.Recovery())
	router.Use(h.loggingMiddleware()) // Наш кастомний логер запитів

	// CORS конфігурація
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * 24 * 3600,
	}))

	api := router.Group("/api/v1")

	{
		crops := api.Group("/crops")

		{
			crops.POST("", h.createCrop)
			crops.GET("", h.getAllCrops)
			crops.GET("/:id", h.getCropById)
			crops.PUT("/:id", h.updateCrop)
			crops.DELETE("/:id", h.deleteCrop)
		}

		fields := api.Group("/fields")
		{
			fields.POST("", h.createField)
			fields.GET("", h.getAllFields)
			fields.GET("/:id", h.getFieldById)
			fields.PUT("/:id", h.updateField)
			fields.DELETE("/:id", h.deleteField)

			// sensor
			fields.GET("/:id/sensors", h.getSensorsByField)

			// forecast
			fields.GET("/:id/forecast/latest", h.getLatestForecast)
		}
		reports := api.Group("/reports")
		{
			reports.GET("/fields/:id", h.getFieldReport)
		}
		sensors := api.Group("/sensors")
		{
			sensors.POST("", h.registerSensor)
			sensors.GET("/:id", h.getSensorById)
			sensors.PATCH("/:id/status", h.updateSensorStatus)
			sensors.DELETE("/:id", h.deleteSensor)

			// metrics
			sensors.GET("/:id/metrics/latest", h.getLatestMetric)
			sensors.GET("/:id/metrics/history", h.getMetricHistory)
		}

		metrics := api.Group("/metrics")
		{
			metrics.POST("", h.saveMetric)
		}

		weather := api.Group("/weather")
		{
			weather.GET("/stations", h.getWeatherStations)
			weather.GET("/stations/:external_id/observations", h.getWeatherStationObservations)
			weather.POST("/sync/station/:external_id", h.syncWeatherStation)
			weather.POST("/sync/field/:field_id", h.syncWeatherField)
		}

		pests := api.Group("/pests")
		{
			pests.POST("", h.createPest)
			pests.GET("", h.getAllPests)
			pests.GET("/:id", h.getPestById)
			pests.PUT("/:id", h.updatePest)
			pests.DELETE("/:id", h.deletePest)
		}

		forecasts := api.Group("/forecasts")
		{
			// Запуск розрахунку прогнозу
			forecasts.POST("/predict", h.predict)
		}
	}

	return router
}

// Допоміжні методи для відповідей
func (h *Handler) newErrorResponse(c *gin.Context, statusCode int, message string) {
	if statusCode >= 500 {
		h.log.Error("API Error", zap.Int("status", statusCode), zap.String("message", message))
	} else {
		h.log.Info("API Client Error", zap.Int("status", statusCode), zap.String("message", message))
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}

func (h *Handler) newSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.log.Info("Request started",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		c.Next()
	}
}

func (h *Handler) getWeatherStations(c *gin.Context) {
	stations, err := h.services.Weather.ListStations()
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to list weather stations")
		return
	}
	h.newSuccessResponse(c, stations)
}

func (h *Handler) getWeatherStationObservations(c *gin.Context) {
	externalID, err := strconv.Atoi(c.Param("external_id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid station external_id")
		return
	}

	station, err := h.services.Weather.GetStationByExternalID(externalID)
	if err != nil {
		h.newErrorResponse(c, http.StatusNotFound, "station not found")
		return
	}

	observations, err := h.services.Weather.GetLatestObservations(station.ID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to get observations")
		return
	}

	h.newSuccessResponse(c, observations)
}

func (h *Handler) syncWeatherStation(c *gin.Context) {
	externalID, err := strconv.Atoi(c.Param("external_id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid station external_id")
		return
	}

	observations, err := h.services.Weather.SyncStation(c.Request.Context(), externalID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "weather station sync failed")
		return
	}

	h.newSuccessResponse(c, observations)
}

func (h *Handler) syncWeatherField(c *gin.Context) {
	fieldID, err := strconv.Atoi(c.Param("field_id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid field_id")
		return
	}

	observations, err := h.services.Weather.SyncField(c.Request.Context(), fieldID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "weather field sync failed")
		return
	}

	h.newSuccessResponse(c, observations)
}
