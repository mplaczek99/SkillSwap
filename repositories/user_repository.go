package repositories

import (
	"github.com/mplaczek99/SkillSwap/config"
	"github.com/mplaczek99/SkillSwap/models"
)

// InsertUser saves a new user record to the database.
func InsertUser(user *models.User) (*models.User, error) {
	result := config.DB.Create(user)
	return user, result.Error
}

// GetUserByEmail finds a user by their email.
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	return &user, result.Error
}

// GetUserByID retrieves a user by ID.
// Note: In a production system, you might want to use a numeric ID type.
func GetUserByID(id string) (*models.User, error) {
	var user models.User
	result := config.DB.Where("id = ?", id).First(&user)
	return &user, result.Error
}

