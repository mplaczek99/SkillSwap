package config

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/mplaczek99/SkillSwap/models"
)

// ConnectDB opens a connection to your MySQL database using the DSN from environment variables.
func ConnectDB() *gorm.DB {
	dsn := os.Getenv("DB_SOURCE") // e.g., "root:password@tcp(db:3306)/skillswap_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

// Migrate runs AutoMigrate on your models.
func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}) // Add other models like &models.Skill{}, &models.Transaction{} as needed.
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
