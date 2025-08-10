package main

import (
	"log"
	"sit-iot-message-mng-api/config"
	"sit-iot-message-mng-api/database"
	"sit-iot-message-mng-api/internal/controllers"
	"sit-iot-message-mng-api/internal/repositories"
	"sit-iot-message-mng-api/internal/routes"
	"sit-iot-message-mng-api/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Firebase
	firebaseApp, err := database.InitFirebase(cfg.FirebaseCredentialsPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	firebaseAuth, err := firebaseApp.Auth(nil)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Auth: %v", err)
	}

	// Initialize repositories
	messageRepo := repositories.NewMessageRepository(db)

	// Initialize services
	messageService := services.NewMessageService(messageRepo, firebaseAuth)

	// Initialize controllers
	messageController := controllers.NewMessageController(messageService)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, messageController, cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
