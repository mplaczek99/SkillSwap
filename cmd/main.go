package main

import (
	"flag"
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
	"github.com/mplaczek99/SkillSwap/setup" // <-- new package for Postgres setup

	_ "github.com/mplaczek99/SkillSwap/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// getenv is a small helper to read environment variables with a fallback.
func getenv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	// 1) Parse CLI flags
	setupDB := flag.Bool("setup-db", false, "Install and setup PostgreSQL before starting the server")
	flag.Parse()

	// 2) Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 3) If --setup-db is passed, run the Postgres setup routines
	if *setupDB {
		if err := setup.CheckAndInstallPostgres(); err != nil {
			log.Fatalf("Error installing PostgreSQL: %v", err)
		}
		if err := setup.InitializePostgres(); err != nil {
			log.Fatalf("Error initializing PostgreSQL: %v", err)
		}
		if err := setup.StartPostgres(); err != nil {
			log.Fatalf("Error starting PostgreSQL: %v", err)
		}
		if err := setup.SetupDatabase(); err != nil {
			log.Fatalf("Error setting up database: %v", err)
		}
	}

	// 4) Connect to the database (config/ConnectDB)
	db := config.ConnectDB()

	// 5) Run auto-migration for your models (config/Migrate)
	config.Migrate(db)

	// 6) Initialize repositories, services, and controllers
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// 7) Set up the Gin router
	router := gin.Default()

	// 8) Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 9) Setup app routes
	routes.SetupRoutes(router, authController)

	// 10) Determine which port to run on (from environment, default = 8080)
	port := getenv("SERVER_PORT", "8080")

	// 11) Start the server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
