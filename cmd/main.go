package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mplaczek99/SkillSwap/config"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/repositories"
	"github.com/mplaczek99/SkillSwap/routes"
	"github.com/mplaczek99/SkillSwap/services"

	_ "github.com/mplaczek99/SkillSwap/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
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

	// 3.1) Seed a default test user if not already present (for development/testing)
	testEmail := "test@example.com"
	testPassword := "somepassword" // Plain text; will be hashed by the BeforeSave hook

	var user models.User
	if err := db.Where("email = ?", testEmail).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser := models.User{
				Name:     "Test User",
				Email:    testEmail,
				Password: testPassword,
			}
			if err := db.Create(&newUser).Error; err != nil {
				log.Printf("Failed to create test user: %v", err)
			} else {
				log.Printf("Test user created: %s / %s", testEmail, testPassword)
			}
		} else {
			log.Printf("Error checking for test user: %v", err)
		}
	} else {
		log.Printf("Test user already exists: %s", testEmail)
	}

	// 4) Initialize repositories, services, and controllers
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// 5) Set up the Gin router
	router := gin.Default()

	// Enable CORS middleware (if needed for frontend)
	router.Use(cors.Default())

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
