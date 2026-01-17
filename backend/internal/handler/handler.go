package handler

import (
	"agro-data-management-system/internal/service"
	"net/http"

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
	router.Use(h.corsMiddleware())

	api := router.Group("/api/v1")
	{
		crops := api.Group("/crops")
		{
			crops.POST("/", h.createCrop)
			crops.GET("/", h.getAllCrops)
			crops.GET("/:id", h.getCropById)
			crops.PUT("/:id", h.updateCrop)
			crops.DELETE("/:id", h.deleteCrop)
		}

		fields := api.Group("/fields")
		{
			fields.POST("/", h.createField)
			fields.GET("/", h.getAllFields)
			fields.GET("/:id", h.getFieldById)
			fields.PUT("/:id", h.updateField)
			fields.DELETE("/:id", h.deleteField)

			fields.GET("/:id/sensors", h.getSensorsByField)
		}

		sensors := api.Group("/sensors")
		{
			sensors.POST("/", h.registerSensor)
			sensors.GET("/:id", h.getSensorById)
			sensors.PATCH("/:id/status", h.updateSensorStatus)
			sensors.DELETE("/:id", h.deleteSensor)

			// metrics
			sensors.GET("/:id/metrics/latest", h.getLatestMetric)
			sensors.GET("/:id/metrics/history", h.getMetricHistory)
		}

		metrics := api.Group("/metrics")
		{
			metrics.POST("/", h.saveMetric)
		}

		pests := api.Group("/pests")
		{
			pests.POST("/", h.createPest)
			pests.GET("/", h.getAllPests)
			pests.GET("/:id", h.getPestById)
			pests.PUT("/:id", h.updatePest)
			pests.DELETE("/:id", h.deletePest)
		}
	}

	return router
}

func (h *Handler) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Для розробки можна *, для продакшну — конкретний URL
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Допоміжні методи для відповідей
func (h *Handler) newErrorResponse(c *gin.Context, statusCode int, message string) {
	h.log.Error("API Error", zap.Int("status", statusCode), zap.String("message", message))
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
