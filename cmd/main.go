package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/config"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/repositories"
	"github.com/mplaczek99/SkillSwap/routes"
	"github.com/mplaczek99/SkillSwap/services"
)

func main() {
	// Connect to the database
	db := config.ConnectDB()

	// Run auto-migration for your models
	config.Migrate(db)

	// Initialize repositories, services, and controllers
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// Initialize the Gin router
	router := gin.Default()

	// Setup routes (this will include public, authenticated, and admin endpoints)
	routes.SetupRoutes(router, authController)

	// Retrieve port from environment variable, defaulting to 8080 if not set
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Start the Gin server on the specified port
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
