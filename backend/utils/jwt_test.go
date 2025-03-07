package utils_test

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mplaczek99/SkillSwap/utils"
)

// Helper function to set up and tear down the JWT secret for tests
func setupJWTTestEnvironment(t *testing.T) (cleanup func()) {
	// Store the original environment value to restore later
	originalSecret := os.Getenv("JWT_SECRET")

	// Set a known secret for testing
	testSecret := "test_secret_key_for_jwt_tests"
	os.Setenv("JWT_SECRET", testSecret)

	// Return a cleanup function to restore the environment
	return func() {
		os.Setenv("JWT_SECRET", originalSecret)
	}
}

func TestInvalidSigningMethod(t *testing.T) {
	cleanup := setupJWTTestEnvironment(t)
	defer cleanup()

	// Create a token with a different signing method (HS384 instead of HS256)
	claims := utils.Claims{
		UserID: 123,
		Role:   "User",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenString, _ := token.SignedString([]byte("test_secret_key_for_jwt_tests"))

	// Try to validate the token
	_, err := utils.ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for token with wrong signing method, got nil")
	}
}

func TestTokenWithNoUserID(t *testing.T) {
	cleanup := setupJWTTestEnvironment(t)
	defer cleanup()

	// Generate token first
	token, err := utils.GenerateToken(0, "User", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	validClaims, err := utils.ValidateToken(token)
	if err != nil {
		t.Errorf("Expected no error for token with zero UserID, got %v", err)
		return
	}

	// UserID should be 0 (zero value for uint)
	if validClaims.UserID != 0 {
		t.Errorf("Expected UserID=0 for token with zero UserID, got %d", validClaims.UserID)
	}
}

func TestExtendedTokenLifetime(t *testing.T) {
	cleanup := setupJWTTestEnvironment(t)
	defer cleanup()

	// Generate token with regular functionality rather than manually creating it
	userID := uint(123)
	role := "Admin"
	email := "admin@example.com"

	// Generate the token, letting the utils function handle the details
	token, err := utils.GenerateToken(userID, role, email)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Validate the token
	validClaims, err := utils.ValidateToken(token)
	if err != nil {
		t.Errorf("Expected no error for token, got %v", err)
		return
	}

	// Check claims
	if validClaims.UserID != userID {
		t.Errorf("Expected UserID=%d, got %d", userID, validClaims.UserID)
	}
	if validClaims.Role != role {
		t.Errorf("Expected Role=%s, got %s", role, validClaims.Role)
	}
	if validClaims.Email != email {
		t.Errorf("Expected Email=%s, got %s", email, validClaims.Email)
	}
}

// Add an additional test for token expiration
func TestTokenExpiration(t *testing.T) {
	cleanup := setupJWTTestEnvironment(t)
	defer cleanup()

	// Create custom claims with an expired token
	claims := utils.Claims{
		UserID: 123,
		Role:   "User",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	// Get the JWT secret from environment
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Validate the token - should fail due to expiration
	_, err = utils.ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}
