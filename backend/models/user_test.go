package models_test

import (
	"testing"

	"github.com/mplaczek99/SkillSwap/models"
	"gorm.io/gorm"
)

func TestUserBeforeSave(t *testing.T) {
	t.Run("Hash Plain Password", func(t *testing.T) {
		user := &models.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "plainpassword",
		}

		// Call BeforeSave
		err := user.BeforeSave(&gorm.DB{})
		if err != nil {
			t.Fatalf("BeforeSave failed: %v", err)
		}

		// Password should now be hashed
		if user.Password == "plainpassword" {
			t.Error("Expected password to be hashed, but it remains plain")
		}

		// Check that it looks like a bcrypt hash (starts with $2a$ or $2b$)
		if len(user.Password) < 60 || (user.Password[:4] != "$2a$" && user.Password[:4] != "$2b$") {
			t.Errorf("Expected bcrypt hash, got: %s", user.Password)
		}
	})

	t.Run("Already Hashed Password", func(t *testing.T) {
		// Create a user with a string that looks like a bcrypt hash
		hashedPassword := "$2a$10$abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ12"
		user := &models.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		// Call BeforeSave
		err := user.BeforeSave(&gorm.DB{})
		if err != nil {
			t.Fatalf("BeforeSave failed: %v", err)
		}

		// Password should remain the same
		if user.Password != hashedPassword {
			t.Errorf("Expected password to remain unchanged, got: %s", user.Password)
		}
	})

	t.Run("Set Default Role", func(t *testing.T) {
		user := &models.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
			// Role intentionally omitted
		}

		// Call BeforeSave
		err := user.BeforeSave(&gorm.DB{})
		if err != nil {
			t.Fatalf("BeforeSave failed: %v", err)
		}

		// Role should be set to default "User"
		if user.Role != "User" {
			t.Errorf("Expected Role to be 'User', got: %s", user.Role)
		}
	})

	t.Run("Keep Existing Role", func(t *testing.T) {
		user := &models.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
			Role:     "Admin", // Set a role
		}

		// Call BeforeSave
		err := user.BeforeSave(&gorm.DB{})
		if err != nil {
			t.Fatalf("BeforeSave failed: %v", err)
		}

		// Role should remain "Admin"
		if user.Role != "Admin" {
			t.Errorf("Expected Role to remain 'Admin', got: %s", user.Role)
		}
	})
}

func TestUserComparePassword(t *testing.T) {
	// Create a user with a known password
	user := &models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Call BeforeSave to hash the password
	err := user.BeforeSave(&gorm.DB{})
	if err != nil {
		t.Fatalf("BeforeSave failed: %v", err)
	}

	t.Run("Correct Password", func(t *testing.T) {
		// Check original password
		if !user.ComparePassword("password123") {
			t.Error("Expected ComparePassword to return true for correct password")
		}
	})

	t.Run("Incorrect Password", func(t *testing.T) {
		// Check wrong password
		if user.ComparePassword("wrongpassword") {
			t.Error("Expected ComparePassword to return false for incorrect password")
		}
	})
}
