package repositories_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mplaczek99/SkillSwap/models"
)

// MockUserRepository directly tests the interface, not the implementation
type MockUserRepository struct {
	users      map[string]*models.User
	errToThrow error
	// Add these methods for testing
	UpdateUser             func(user *models.User) error
	GetUsersWithPagination func(limit, offset int) ([]*models.User, int, error)
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
		// Initialize with stub implementations
		UpdateUser: func(user *models.User) error {
			return errors.New("method not implemented")
		},
		GetUsersWithPagination: func(limit, offset int) ([]*models.User, int, error) {
			return nil, 0, errors.New("method not implemented")
		},
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

	// Store a copy - make sure to include Bio field
	m.users[user.Email] = &models.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Bio:       user.Bio, // Important: Include Bio field
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
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

	// Return a copy - make sure to include Bio field
	return &models.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Bio:       user.Bio, // Important: Include Bio field
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
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
			t.Errorf("Failed to retrieve registered user: %v", err)
		}

		if savedUser == nil {
			t.Fatal("Expected saved user, got nil")
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

func TestUserRepository_UpdateAndPagination(t *testing.T) {
	mockRepo := NewMockUserRepository()

	// Create several test users
	for i := 1; i <= 10; i++ {
		user := &models.User{
			Name:     fmt.Sprintf("User %d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Password: "password",
			Role:     "User",
			Bio:      "", // Initialize with empty bio
		}
		if err := mockRepo.CreateUser(user); err != nil {
			t.Fatalf("Failed to create user %d: %v", i, err)
		}
	}

	// Add an UpdateUser method to the mock repository for testing
	mockRepo.UpdateUser = func(user *models.User) error {
		// Find the user by email
		existingUser, exists := mockRepo.users[user.Email]
		if !exists {
			return errors.New("user not found")
		}

		// Update fields
		if user.Name != "" {
			existingUser.Name = user.Name
		}
		if user.Bio != "" {
			existingUser.Bio = user.Bio // Make sure this is copying correctly
		}
		if user.Role != "" && user.Role != existingUser.Role {
			existingUser.Role = user.Role
		}

		// Store updated user back in the map
		mockRepo.users[user.Email] = existingUser
		return nil
	}

	// Add a GetUsersWithPagination method to the mock repository
	mockRepo.GetUsersWithPagination = func(limit, offset int) ([]*models.User, int, error) {
		// Get all users
		users := make([]*models.User, 0, len(mockRepo.users))
		for _, user := range mockRepo.users {
			users = append(users, user)
		}

		totalCount := len(users)

		// Apply pagination
		if offset >= totalCount {
			return []*models.User{}, totalCount, nil
		}

		end := offset + limit
		if end > totalCount {
			end = totalCount
		}

		return users[offset:end], totalCount, nil
	}

	// Test updating a user's profile
	t.Run("Update User Profile", func(t *testing.T) {
		// Get the first user
		firstUserEmail := "user1@example.com"
		user, err := mockRepo.GetUserByEmail(firstUserEmail)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		// Update profile
		user.Name = "Updated Name"
		user.Bio = "This is my updated bio"
		err = mockRepo.UpdateUser(user)
		if err != nil {
			t.Errorf("UpdateUser returned error: %v", err)
		}

		// Verify update
		updated, err := mockRepo.GetUserByEmail(firstUserEmail)
		if err != nil {
			t.Fatalf("Failed to get updated user: %v", err)
		}

		if updated.Name != "Updated Name" {
			t.Errorf("Expected name to be 'Updated Name', got '%s'", updated.Name)
		}
		if updated.Bio != "This is my updated bio" {
			t.Errorf("Expected bio update, got: '%s'", updated.Bio)
		}
	})

	// Test updating user's role
	t.Run("Update User Role", func(t *testing.T) {
		// Get the second user
		secondUserEmail := "user2@example.com"
		user, err := mockRepo.GetUserByEmail(secondUserEmail)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		// Update to admin
		user.Role = "Admin"
		err = mockRepo.UpdateUser(user)
		if err != nil {
			t.Errorf("UpdateUser returned error: %v", err)
		}

		// Verify role update
		updated, err := mockRepo.GetUserByEmail(secondUserEmail)
		if err != nil {
			t.Fatalf("Failed to get updated user: %v", err)
		}

		if updated.Role != "Admin" {
			t.Errorf("Expected role to be 'Admin', got '%s'", updated.Role)
		}
	})

	// Test getting users with pagination
	t.Run("Get Users With Pagination", func(t *testing.T) {
		// Get first page (3 users)
		users, total, err := mockRepo.GetUsersWithPagination(3, 0)
		if err != nil {
			t.Errorf("GetUsersWithPagination returned error: %v", err)
		}

		if total != 10 {
			t.Errorf("Expected total 10 users, got %d", total)
		}

		if len(users) != 3 {
			t.Errorf("Expected 3 users in first page, got %d", len(users))
		}

		// Get second page (3 users)
		users, total, err = mockRepo.GetUsersWithPagination(3, 3)
		if err != nil {
			t.Errorf("GetUsersWithPagination returned error: %v", err)
		}

		if len(users) != 3 {
			t.Errorf("Expected 3 users in second page, got %d", len(users))
		}

		// Get last page (1 user)
		users, total, err = mockRepo.GetUsersWithPagination(3, 9)
		if err != nil {
			t.Errorf("GetUsersWithPagination returned error: %v", err)
		}

		if len(users) != 1 {
			t.Errorf("Expected 1 user in last page, got %d", len(users))
		}
	})
}
