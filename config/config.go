package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/joho/godotenv"
)

// Configuration holds all configurable values
type Configuration struct {
	DBDriver   string
	DBSource   string
	ServerPort string
	JWTSecret  string
}

var (
	// AppConfig is the global configuration instance
	AppConfig *Configuration
	// DB is the global database connection (GORM)
	DB *gorm.DB
)

// LoadConfig loads configuration from environment variables (or .env file)
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}
	AppConfig = &Configuration{
		DBDriver:   getEnv("DB_DRIVER", "mysql"),
		DBSource:   getEnv("DB_SOURCE", "user:password@tcp(127.0.0.1:3306)/skillswap_db?parseTime=true"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecretkey"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// InitDB initializes the database connection using GORM
func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(AppConfig.DBSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	// Optionally auto-migrate your models here
}

// Automatically load the configuration when the package is imported.
func init() {
	LoadConfig()
}

