package services

import (
	"errors"
	"time"

	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/repositories"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser processes new user registration by hashing the password and saving the user.
func CreateUser(user *models.User) (*models.User, error) {
	if user.Email == "" || user.Password == "" {
		return nil, errors.New("email and password are required")
	}
	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	return repositories.InsertUser(user)
}

// GetUserByID retrieves a user by their ID.
func GetUserByID(id string) (*models.User, error) {
	return repositories.GetUserByID(id)
}
