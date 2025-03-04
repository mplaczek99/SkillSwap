package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret reads the JWT secret from the environment or generates a random one if not set.
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("WARNING: JWT_SECRET not set! Using a randomly generated secret which will change on restart.")
		// Generate a random 32-byte secret
		randomBytes := make([]byte, 32)
		if _, err := rand.Read(randomBytes); err != nil {
			log.Println("Failed to generate random secret, using fallback secret")
			return "temporary_fallback_secret_key_please_set_JWT_SECRET"
		}
		return base64.StdEncoding.EncodeToString(randomBytes)
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
