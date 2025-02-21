package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mplaczek99/SkillSwap/config"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/repositories"
	"github.com/mplaczek99/SkillSwap/routes"
	"github.com/mplaczek99/SkillSwap/services"

	_ "github.com/mplaczek99/SkillSwap/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// getenv is a small helper to read environment variables or a fallback.
func getenv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	// 1) Load environment variables (optional if .env doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	// 2) Connect to the database
	db := config.ConnectDB()

	// 3) Run migrations
	config.Migrate(db)

	// 4) Initialize repositories, services, and controllers
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// 5) Set up the Gin router
	router := gin.Default()

	// 6) Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 7) Setup routes
	routes.SetupRoutes(router, authController)

	// 8) Determine which port to run on (from environment, default = 8080)
	port := getenv("SERVER_PORT", "8080")
	addr := fmt.Sprintf(":%s", port)

	// 9) Start the server
	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
