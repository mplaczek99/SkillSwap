package config

import (
	"log"
	"os"

	"github.com/mplaczek99/SkillSwap/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB opens a connection to your PostgreSQL database using the DSN from environment variables.
func ConnectDB() *gorm.DB {
	// Example DSN for PostgreSQL:
	// "host=localhost user=postgres password=your_postgres_password dbname=skillswap_db port=5432 sslmode=disable"
	dsn := os.Getenv("DB_SOURCE")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
