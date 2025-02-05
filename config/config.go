package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver   string
	DBSource   string
	ServerPort string
	JWTSecret  string
}

var AppConfig *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading config from ENV")
	}
	AppConfig = &Config{
		DBDriver:   getEnv("DB_DRIVER", "mysql"),
		DBSource:   getEnv("DB_SOURCE", "username:password@tcp(127.0.0.1:3306)/skillswap_db?parseTime=true"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecretkey"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

