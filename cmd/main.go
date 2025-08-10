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

	log.Printf("Using database provider: %s", cfg.DatabaseProvider)

	// Initialize databases based on configuration
	dbClients, err := database.InitDatabases(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}

	// Initialize Firebase for authentication
	firebaseApp, err := database.InitFirebase(cfg.FirebaseCredentialsPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	firebaseAuth, err := firebaseApp.Auth(nil)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Auth: %v", err)
	}

	// Initialize repository factory
	repoFactory := repositories.NewRepositoryFactory(cfg)

	// Create message repository using the factory
	messageRepo, err := repoFactory.CreateMessageRepository(dbClients.MongoDB, dbClients.Firestore)
	if err != nil {
		log.Fatalf("Failed to create message repository: %v", err)
	}

	log.Printf("Message repository initialized for %s", repoFactory.GetDatabaseProvider())

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
