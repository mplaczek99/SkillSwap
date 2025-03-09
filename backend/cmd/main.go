package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

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

	// 3.1) Seed test users if not already present (for development/testing)
	testUsers := []struct {
		name     string
		email    string
		password string
	}{
		{"Test User", "test@example.com", "somepassword"},
		{"Test User 2", "test2@example.com", "somepassword2"},
	}

	// Create each test user if they don't exist
	for _, testUser := range testUsers {
		var user models.User
		if err := db.Where("email = ?", testUser.email).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newUser := models.User{
					Name:     testUser.name,
					Email:    testUser.email,
					Password: testUser.password,
				}
				if err := db.Create(&newUser).Error; err != nil {
					log.Printf("Failed to create test user: %v", err)
				} else {
					log.Printf("Test user created: %s / %s", testUser.email, testUser.password)
				}
			} else {
				log.Printf("Error checking for test user: %v", err)
			}
		} else {
			log.Printf("Test user already exists: %s", testUser.email)
		}
	}

	// 4) Initialize repositories, services, and controllers
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// 5) Set up the Gin router
	router := gin.Default()

	// Enable CORS middleware with updated configuration
	corsConfig := cors.DefaultConfig()

	// Get allowed origins from environment, defaulting to localhost
	allowedOrigins := getenv("CORS_ALLOWED_ORIGINS", "http://localhost:8081,http://frontend:80")

	// Check if we should allow all origins
	if os.Getenv("CORS_ALLOW_ALL") == "true" {
		corsConfig.AllowAllOrigins = true
	} else {
		// Only set specific origins if we're not allowing all
		corsConfig.AllowOrigins = strings.Split(allowedOrigins, ",")
	}

	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type", "Origin"}
	corsConfig.AllowCredentials = true
	corsConfig.ExposeHeaders = []string{"Content-Length"}

	router.Use(cors.New(corsConfig))

	// Add database to the gin context so it can be used in controllers
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// 6) Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 7) Setup routes
	routes.SetupRoutes(router, authController)

	// Create uploads directory if it doesn't exist
	os.MkdirAll("./uploads", os.ModePerm)

	// 8) Determine which port to run on (from environment, default = 8080)
	port := getenv("SERVER_PORT", "8080")
	addr := fmt.Sprintf(":%s", port)

	// 9) Start the server
	log.Printf("Server starting on port %s\n", port)
	log.Printf("CORS configuration: AllowAllOrigins=%v, AllowedOrigins=%v",
		corsConfig.AllowAllOrigins, corsConfig.AllowOrigins)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
