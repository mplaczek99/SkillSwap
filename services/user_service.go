package services

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/your-username/skillswap/models"
	"github.com/your-username/skillswap/repositories"
)

func CreateUser(user *models.User) (*models.User, error) {
	// Validate input data (e.g., required fields)

	// Hash the password securely
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	// Save user using repository layer
	return repositories.InsertUser(user)
}

