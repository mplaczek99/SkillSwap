package repositories

import (
	"errors"
	"github.com/mplaczek99/SkillSwap/models"
  "golang.org/x/crypto/bcrypt"
)

// InsertUser returns the user with a dummy ID.
func InsertUser(user *models.User) (*models.User, error) {
	user.ID = 1
	return user, nil
}

// GetUserByEmail returns a dummy user for a specific email.
func GetUserByEmail(email string) (*models.User, error) {
	// Return a dummy user if the email matches; otherwise, return an error.
	if email == "test@example.com" {
		// Generate a valid bcrypt hash for "somepassword"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("somepassword"), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		return &models.User{
			ID:       1,
			Email:    email,
			Password: string(hashedPassword),
			Name:     "Dummy User",
		}, nil
	}
	return nil, errors.New("user not found")
}

// GetUserByID returns a dummy user.
func GetUserByID(id string) (*models.User, error) {
	if id != "1" {
		return nil, errors.New("user not found")
	}
	return &models.User{
		ID:    1,
		Email: "test@example.com",
		Name:  "Dummy User",
	}, nil
}
