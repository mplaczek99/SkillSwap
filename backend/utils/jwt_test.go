package utils_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mplaczek99/SkillSwap/utils"
)

func TestGenerateAndValidateToken(t *testing.T) {
	tokenString, err := utils.GenerateToken(123, "User", "test@example.com")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if tokenString == "" {
		t.Error("expected a token string, got empty")
	}

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims == nil {
		t.Fatal("expected valid claims, got nil")
	}
	if claims.UserID != 123 || claims.Email != "test@example.com" || claims.Role != "User" {
		t.Errorf("claims do not match expected values, got %+v", claims)
	}
}

func TestValidateTokenExpired(t *testing.T) {
	// Create an expired token manually.
	expiredClaims := utils.Claims{
		UserID: 999,
		Role:   "User",
		Email:  "expired@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Already expired.
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	// Use the same default secret as in the utils package.
	tokenStr, err := token.SignedString([]byte("default_secret_key"))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	_, err = utils.ValidateToken(tokenStr)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}
