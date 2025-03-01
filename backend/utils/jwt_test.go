package utils_test

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mplaczek99/SkillSwap/utils"
)

func TestGenerateAndValidateToken(t *testing.T) {
	// Set a consistent JWT secret for tests
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test_secret_key")
	defer os.Setenv("JWT_SECRET", originalSecret)

	t.Run("Generate Valid Token", func(t *testing.T) {
		token, err := utils.GenerateToken(123, "User", "test@example.com")
		if err != nil {
			t.Fatalf("GenerateToken failed: %v", err)
		}
		if token == "" {
			t.Error("Expected non-empty token string")
		}

		// Validate the token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			t.Fatalf("ValidateToken failed: %v", err)
		}

		// Check claims
		if claims.UserID != 123 {
			t.Errorf("Expected UserID=123, got %d", claims.UserID)
		}
		if claims.Role != "User" {
			t.Errorf("Expected Role=User, got %s", claims.Role)
		}
		if claims.Email != "test@example.com" {
			t.Errorf("Expected Email=test@example.com, got %s", claims.Email)
		}
	})

	t.Run("Validate Expired Token", func(t *testing.T) {
		// Create an expired token
		claims := utils.Claims{
			UserID: 123,
			Role:   "User",
			Email:  "test@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("test_secret_key"))
		if err != nil {
			t.Fatalf("Failed to create expired token: %v", err)
		}

		// Try to validate the expired token
		_, err = utils.ValidateToken(tokenString)
		if err == nil {
			t.Error("Expected error for expired token, got nil")
		}
	})

	t.Run("Validate Invalid Token", func(t *testing.T) {
		// Try to validate an invalid token
		_, err := utils.ValidateToken("invalid.token.string")
		if err == nil {
			t.Error("Expected error for invalid token, got nil")
		}
	})

	t.Run("Validate Token With Wrong Signature", func(t *testing.T) {
		// Create a token with a different secret
		claims := utils.Claims{
			UserID: 123,
			Role:   "User",
			Email:  "test@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("different_secret_key"))
		if err != nil {
			t.Fatalf("Failed to create token with wrong signature: %v", err)
		}

		// Try to validate with the wrong signature
		_, err = utils.ValidateToken(tokenString)
		if err == nil {
			t.Error("Expected error for token with wrong signature, got nil")
		}
	})
}
