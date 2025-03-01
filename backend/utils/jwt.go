package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret reads the JWT secret from the environment or returns a default (for development).
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_key" // Warning: use a strong secret in production!
	}
	return secret
}

var secretKey = []byte(getJWTSecret())

// Claims defines the JWT claims structure.
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token with the user's ID, role, and email.
func GenerateToken(userID uint, role, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours.
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	return token.SignedString(secretKey)
}

// ValidateToken parses and validates a JWT token string and returns its claims.
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
