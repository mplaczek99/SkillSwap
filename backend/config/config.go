package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mplaczek99/SkillSwap/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// AppConfig holds all application configuration with default values
type AppConfig struct {
	// Server settings
	ServerPort string

	// Database settings
	DBDriver string
	DBSource string

	// Security settings
	JWTSecret string

	// CORS settings
	CORSAllowedOrigins []string
	CORSAllowAll       bool
	CORSMaxAge         time.Duration
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *AppConfig {
	config := &AppConfig{
		// Default values
		ServerPort:         "8080",
		DBDriver:           "postgres",
		DBSource:           "host=db port=5432 user=techie password=techiestrongpassword dbname=skillswap_db sslmode=disable",
		CORSAllowedOrigins: []string{"http://localhost:8081", "http://frontend:80"},
		CORSAllowAll:       false,
		CORSMaxAge:         12 * time.Hour,
	}

	// Override with environment variables if they exist
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.ServerPort = port
	}

	if driver := os.Getenv("DB_DRIVER"); driver != "" {
		config.DBDriver = driver
	}

	if source := os.Getenv("DB_SOURCE"); source != "" {
		config.DBSource = source
	}

	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		config.JWTSecret = secret
	} else {
		// JWT secret is critical - log a warning if not set
		log.Println("WARNING: JWT_SECRET not set. Using an insecure default value!")
		config.JWTSecret = "your_secret_key_should_be_long_and_secure_in_production"
	}

	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		config.CORSAllowedOrigins = strings.Split(origins, ",")
	}

	if allowAll := os.Getenv("CORS_ALLOW_ALL"); allowAll != "" {
		config.CORSAllowAll, _ = strconv.ParseBool(allowAll)
	}

	// Validate critical configuration
	if config.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return config
}

// ConnectDB opens a connection to your PostgreSQL database using the DSN from environment variables.
func ConnectDB() *gorm.DB {
	config := LoadConfig()

	db, err := gorm.Open(postgres.Open(config.DBSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

// Migrate runs AutoMigrate on your models.
func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Skill{},
		&models.Transaction{},
		&models.Schedule{},
		&models.Job{}, // Add Job model to migrations
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
