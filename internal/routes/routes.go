package routes

import (
	"sit-iot-message-mng-api/config"
	"sit-iot-message-mng-api/internal/controllers"
	"sit-iot-message-mng-api/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, messageController *controllers.MessageController, cfg *config.Config) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:5173", "http://127.0.0.1:5173", "https://console.sit-iot.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Range"},
		ExposeHeaders:    []string{"Content-Length", "Content-Range"},
		AllowCredentials: true,
	}))

	api := router.Group("/api", middleware.IdentityPlatformMiddleware(cfg))
	{
		// Message routes
		api.GET("/message/:id", messageController.GetMessage)

		// Device-specific message routes
		api.GET("/message/device/:deviceId", messageController.ListMessagesByDevice)

		// Aggregated data for device (for graphing max, min, avg)
		api.GET("/message/aggregations/device/:deviceId", messageController.GetAggregatedDataByDevice)
	}
}
