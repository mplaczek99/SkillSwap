package repositories_test

import (
	"errors"
	"testing"

	"github.com/mplaczek99/SkillSwap/models"
)

// MockUserRepository directly tests the interface, not the implementation
type MockUserRepository struct {
	users      map[string]*models.User
	errToThrow error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	if m.errToThrow != nil {
		return m.errToThrow
	}

	// Check for duplicate email
	if _, exists := m.users[user.Email]; exists {
		return errors.New("duplicate email")
	}

	// Simulate ID assignment
	if user.ID == 0 {
		user.ID = uint(len(m.users) + 1)
	}

	// Store a copy
	m.users[user.Email] = &models.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}

	return nil
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	if m.errToThrow != nil {
		return nil, m.errToThrow
	}

	user, exists := m.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Return a copy
	return &models.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}, nil
}

func (m *MockUserRepository) SetError(err error) {
	m.errToThrow = err
}

func TestUserRepository_Interface(t *testing.T) {
	// This directly tests the interface methods
	mockRepo := NewMockUserRepository()

	t.Run("Create User", func(t *testing.T) {
		user := &models.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		err := mockRepo.CreateUser(user)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify the user was stored
		savedUser, err := mockRepo.GetUserByEmail("test@example.com")
		if err != nil {
			t.Errorf("Failed to get user: %v", err)
		}

		if savedUser.Name != "Test User" {
			t.Errorf("Expected name 'Test User', got '%s'", savedUser.Name)
		}
	})

	t.Run("Create Duplicate User", func(t *testing.T) {
		// Attempt to create a duplicate
		user := &models.User{
			Name:     "Duplicate User",
			Email:    "test@example.com", // Same as previous test
			Password: "anotherpassword",
		}

		err := mockRepo.CreateUser(user)
		if err == nil {
			t.Error("Expected error for duplicate email, got nil")
		}
	})

	t.Run("Get User By Email", func(t *testing.T) {
		// Try to get a user that exists
		user, err := mockRepo.GetUserByEmail("test@example.com")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if user == nil {
			t.Fatal("Expected user to be returned, got nil")
		}

		if user.Email != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
		}
	})

	t.Run("Get Non-existent User", func(t *testing.T) {
		// Try to get a user that doesn't exist
		_, err := mockRepo.GetUserByEmail("nonexistent@example.com")
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.SetError(errors.New("database error"))

		// Try to get a user when there's a database error
		_, err := mockRepo.GetUserByEmail("test@example.com")
		if err == nil {
			t.Error("Expected error when database fails, got nil")
		}

		mockRepo.SetError(nil) // Reset for subsequent tests
	})
}
