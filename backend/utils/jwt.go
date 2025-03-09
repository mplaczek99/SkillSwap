package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Static secret to prevent regeneration on each call
var (
	secretKey []byte
	once      sync.Once
)

// getJWTSecret reads the JWT secret from the environment or generates a random one if not set.
// Uses sync.Once to ensure the secret is only generated once during the application lifecycle.
func getJWTSecret() []byte {
	once.Do(func() {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Println("WARNING: JWT_SECRET not set! Using a randomly generated secret which will remain consistent until restart.")
			// Generate a random 32-byte secret
			randomBytes := make([]byte, 32)
			if _, err := rand.Read(randomBytes); err != nil {
				log.Println("Failed to generate random secret, using fallback secret")
				secretKey = []byte("temporary_fallback_secret_key_please_set_JWT_SECRET")
				return
			}
			secretKey = []byte(base64.StdEncoding.EncodeToString(randomBytes))
			return
		}
		secretKey = []byte(secret)
	})
	return secretKey
}

// Claims defines the JWT claims structure.
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token with the user's ID, role, and email.
func GenerateToken(userID uint, role, email string) (string, error) {
	secret := getJWTSecret()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours.
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})
	return token.SignedString(secret)
}

// ValidateToken parses and validates a JWT token string and returns its claims.
func ValidateToken(tokenString string) (*Claims, error) {
	secret := getJWTSecret()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
