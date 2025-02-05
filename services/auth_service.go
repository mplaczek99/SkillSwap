package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/repositories"
	"github.com/mplaczek99/SkillSwap/utils"
)

// AuthenticateUser verifies user credentials and returns a JWT token upon success.
func AuthenticateUser(email, password string) (string, error) {
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	// Compare the hashed password with the provided one
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}

