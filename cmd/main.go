package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mplaczek99/SkillSwap/config"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/repositories"
	"github.com/mplaczek99/SkillSwap/routes"
	"github.com/mplaczek99/SkillSwap/services"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/mplaczek99/SkillSwap/docs"
)

// @title SkillSwap API
// @version 1.0
// @description This is a SkillSwap API server.
// @host localhost:8080
// @BasePath /
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

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

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup routes
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
