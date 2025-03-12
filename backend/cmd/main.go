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

func main() {
	// 1) Load environment variables (optional if .env doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	// 2) Load application configuration
	appConfig := config.LoadConfig()

	// 3) Connect to the database
	db := config.ConnectDB()

	// 4) Run migrations
	config.Migrate(db)

	// 5) Seed test users if not already present (for development/testing)
	seedTestUsers(db)

	// 6) Initialize repositories, services, and controllers
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// 7) Set up the Gin router
	router := gin.Default()

	// 8) Enable CORS middleware with configuration from appConfig
	corsConfig := cors.DefaultConfig()

	if appConfig.CORSAllowAll {
		log.Println("CORS: Allowing all origins")
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = appConfig.CORSAllowedOrigins
		log.Printf("CORS: Allowing specific origins: %v", corsConfig.AllowOrigins)
	}

	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type", "Origin", "Accept", "X-Requested-With"}
	corsConfig.AllowCredentials = true
	corsConfig.ExposeHeaders = []string{"Content-Length", "Content-Type"}
	corsConfig.MaxAge = appConfig.CORSMaxAge

	router.Use(cors.New(corsConfig))

	// 9) Add database to the gin context for controllers
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// 10) Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 11) Setup routes
	routes.SetupRoutes(router, authController)

	// 12) Create uploads directory if it doesn't exist
	os.MkdirAll("./uploads", os.ModePerm)

	// 13) Start the server
	addr := fmt.Sprintf(":%s", appConfig.ServerPort)
	log.Printf("Server starting on port %s\n", appConfig.ServerPort)
	log.Printf("CORS configuration: AllowAllOrigins=%v, AllowedOrigins=%v",
		corsConfig.AllowAllOrigins, corsConfig.AllowOrigins)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}

// seedTestUsers creates test users if they don't already exist
func seedTestUsers(db *gorm.DB) {
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
}
