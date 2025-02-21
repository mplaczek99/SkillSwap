package services

import (
	"errors"

	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/repositories"
	"github.com/mplaczek99/SkillSwap/utils"
)

// AuthServiceInterface defines the contract for authentication services.
type AuthServiceInterface interface {
	Register(user *models.User) (string, error)
	Login(email, password string) (string, error)
}

type AuthService struct {
	UserRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) AuthServiceInterface {
	return &AuthService{UserRepo: userRepo}
}

// Register creates a new user and returns a token.
func (s *AuthService) Register(user *models.User) (string, error) {
	existingUser, _ := s.UserRepo.GetUserByEmail(user.Email)
	if existingUser != nil {
		return "", errors.New("email already in use")
	}

	if err := s.UserRepo.CreateUser(user); err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(user.ID, user.Role, user.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Login authenticates a user and returns a token.
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !user.ComparePassword(password) {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Role, user.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}
