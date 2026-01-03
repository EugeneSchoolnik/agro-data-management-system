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
		// Тут будуть інші групи: fields, sensors, metrics...
	}

	return router
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
