package main

import (
	"bank-management-system/internal/config"
	"bank-management-system/internal/database"
	"bank-management-system/internal/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	if err := database.ConnectDB(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	// Setup routes
	router := routes.SetupRoutes()

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
