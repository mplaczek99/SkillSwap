package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/skillswap/config"
	"github.com/your-username/skillswap/routes"
)

func main() {
	// Load environment configuration
	config.LoadConfig()

	// Initialize database connection (with pooling) in config/initDB.go (not shown here)
	// e.g., db := config.InitDB()

	// Create a Gin router with production settings
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	// Attach security middlewares (e.g., CORS, JWT validation)

	// Set up routes
	routes.SetupRoutes(router)

	// Start the server on the configured port
	if err := router.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

