package services_test

import (
	"errors"
	"os"
	"testing"

	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/services"
	"golang.org/x/crypto/bcrypt"
)

// UserRepositoryInterface is a copy of the interface that auth service expects
type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

// MockUserRepository implements the UserRepositoryInterface
type MockUserRepository struct {
	users map[string]*models.User
}

// NewMockUserRepository creates a new mock repository with predefined test users
func NewMockUserRepository() *MockUserRepository {
	// Create a hash of "password123" for our test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Create a mock repository with one existing user
	mockRepo := &MockUserRepository{
		users: make(map[string]*models.User),
	}

	mockRepo.users["existing@example.com"] = &models.User{
		ID:       1,
		Name:     "Existing User",
		Email:    "existing@example.com",
		Password: string(hashedPassword),
		Role:     "User",
	}

	return mockRepo
}

// CreateUser implements the repository interface
func (m *MockUserRepository) CreateUser(user *models.User) error {
	// Check if user already exists
	if _, exists := m.users[user.Email]; exists {
		return errors.New("email already in use")
	}

	// Set an ID
	user.ID = uint(len(m.users) + 1)

	// Hash the password if needed
	if len(user.Password) > 0 && (len(user.Password) < 60 || user.Password[0] != '$') {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	// Store a copy of the user
	m.users[user.Email] = &models.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}

	return nil
}

// GetUserByEmail implements the repository interface
func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Return a copy to prevent modifications
	return &models.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}, nil
}

// Create a custom AuthService for testing that accepts our interface
func newTestAuthService(repo UserRepositoryInterface) *services.AuthService {
	// For now, use direct type assertion to make the mock compatible with the service
	// This requires us to know the internal structure of AuthService
	// A better approach would be to update the real AuthService to work with
	// a repository interface instead of a concrete type

	// Let's assume AuthService accepts an interface in its constructor
	return services.NewAuthService(repo)
}

func TestAuthService_Register(t *testing.T) {
	// Set a consistent JWT secret for tests
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test_secret_key")
	defer os.Setenv("JWT_SECRET", originalSecret)

	mockRepo := NewMockUserRepository()

	// Create our auth service with the interface-compatible mock
	authService := newTestAuthService(mockRepo)

	t.Run("Register New User", func(t *testing.T) {
		user := &models.User{
			Name:     "New User",
			Email:    "new@example.com",
			Password: "newpassword",
		}

		token, err := authService.Register(user)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if token == "" {
			t.Error("Expected non-empty token")
		}

		// Check that user was added to repository
		savedUser, err := mockRepo.GetUserByEmail("new@example.com")
		if err != nil {
			t.Errorf("Failed to retrieve registered user: %v", err)
		}

		if savedUser == nil {
			t.Fatal("Expected saved user, got nil")
		}

		if savedUser.Name != "New User" {
			t.Errorf("Expected name 'New User', got '%s'", savedUser.Name)
		}
	})

	t.Run("Register With Existing Email", func(t *testing.T) {
		user := &models.User{
			Name:     "Duplicate User",
			Email:    "existing@example.com", // This email already exists
			Password: "duplicatepassword",
		}

		_, err := authService.Register(user)
		if err == nil {
			t.Error("Expected error for duplicate email, got nil")
		}
	})
}

func TestAuthService_Login(t *testing.T) {
	// Set a consistent JWT secret for tests
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test_secret_key")
	defer os.Setenv("JWT_SECRET", originalSecret)

	mockRepo := NewMockUserRepository()
	authService := newTestAuthService(mockRepo)

	t.Run("Login With Valid Credentials", func(t *testing.T) {
		token, err := authService.Login("existing@example.com", "password123")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if token == "" {
			t.Error("Expected non-empty token")
		}
	})

	t.Run("Login With Wrong Password", func(t *testing.T) {
		_, err := authService.Login("existing@example.com", "wrongpassword")
		if err == nil {
			t.Error("Expected error for wrong password, got nil")
		}
	})

	t.Run("Login With Non-existent Email", func(t *testing.T) {
		_, err := authService.Login("nonexistent@example.com", "password123")
		if err == nil {
			t.Error("Expected error for non-existent email, got nil")
		}
	})
}
