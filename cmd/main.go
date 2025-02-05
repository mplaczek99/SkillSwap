package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/config"
	"github.com/mplaczek99/SkillSwap/routes"
)

func main() {
	// Load configuration from environment variables (and .env file if available)
	config.LoadConfig()

	// Initialize database connection (using GORM)
	config.InitDB()

	// Create a new Gin router with Logger and Recovery middleware
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup API routes
	routes.SetupRoutes(router)

	// Start the server on the configured port
	port := config.AppConfig.ServerPort
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

