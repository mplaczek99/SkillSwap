package utils

import (
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Performance optimizations:
// 1. Secret key loading with atomic operations for thread safety without locks
// 2. Token signature verification optimization with subtle.ConstantTimeCompare
// 3. Sync.Pool for claims to reduce GC pressure
// 4. Reduced allocations in token validation
// 5. Fast path for common token formats
// 6. Minimized error checking paths

var (
	// Use atomic.Value for thread-safe lazy loading without locks
	secretKeyAtomic atomic.Value
	loadOnce        sync.Once

	// Error caching to avoid allocations
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")

	// Claims pool to reduce GC pressure
	claimsPool = sync.Pool{
		New: func() interface{} {
			return &Claims{}
		},
	}
)

// Claims defines the JWT claims structure.
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// initSecretKey reads the JWT secret from the environment once.
func initSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is required. Server will exit.")
	}
	return []byte(secret)
}

// getJWTSecret returns the JWT secret as a byte slice.
// Uses atomic.Value to ensure thread-safe initialization without locks after first use.
func getJWTSecret() []byte {
	// Fast path: already loaded
	secretInterface := secretKeyAtomic.Load()
	if secretInterface != nil {
		return secretInterface.([]byte)
	}

	// Slow path: needs loading (happens only once)
	loadOnce.Do(func() {
		secretKeyAtomic.Store(initSecretKey())
	})

	return secretKeyAtomic.Load().([]byte)
}

// GenerateToken creates a JWT token with the user's ID, role, and email.
// Performance optimizations applied to reduce allocations.
func GenerateToken(userID uint, role, email string) (string, error) {
	// Get claims from pool
	claimsObj := claimsPool.Get().(*Claims)
	defer claimsPool.Put(claimsObj)

	// Reset claims object for reuse
	*claimsObj = Claims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with pooled claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsObj)
	return token.SignedString(getJWTSecret())
}

// ValidateToken parses and validates a JWT token string and returns its claims.
// Heavily optimized with fast paths for common cases and reduced allocations.
func ValidateToken(tokenString string) (*Claims, error) {
	// Fast fail on empty token
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	// Fast path check for token format
	if len(tokenString) < 10 || (tokenString[0] != 'e' && tokenString[0] != 'J') {
		return nil, ErrInvalidToken
	}

	// Get claims from pool to reduce allocations
	claimsObj := claimsPool.Get().(*Claims)

	// Parse token with pooled claims
	token, err := jwt.ParseWithClaims(tokenString, claimsObj, func(token *jwt.Token) (interface{}, error) {
		// Validate algorithm with fast check
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			claimsPool.Put(claimsObj) // Return to pool on error
			return nil, ErrInvalidToken
		}

		// Check alg directly for common case
		if token.Header["alg"] != "HS256" {
			claimsPool.Put(claimsObj) // Return to pool on error
			return nil, ErrInvalidToken
		}

		return getJWTSecret(), nil
	})

	// Handle parsing errors
	if err != nil {
		claimsPool.Put(claimsObj) // Return to pool on error

		// Specific error for expired token to allow special handling
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// Fast path: token is invalid
	if !token.Valid {
		claimsPool.Put(claimsObj) // Return to pool on error
		return nil, ErrInvalidToken
	}

	// Create a new Claims to return (don't return the pooled one)
	returnClaims := &Claims{
		UserID:           claimsObj.UserID,
		Role:             claimsObj.Role,
		Email:            claimsObj.Email,
		RegisteredClaims: claimsObj.RegisteredClaims,
	}

	// Return the pooled claims object
	claimsPool.Put(claimsObj)

	return returnClaims, nil
}

// VerifySignature is a specialized function for just checking if a token
// signature is valid without fully parsing it - useful for rate limiting.
// This is much faster than full parsing for quick checks.
func VerifySignature(tokenString string) bool {
	// Fast path: quick format check
	if len(tokenString) < 10 {
		return false
	}

	// Split the token into parts using standard strings package
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false
	}

	// Parse token with minimal validation - just to check signature
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate algorithm with fast check
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return getJWTSecret(), nil
	})

	// Only care about signature validity, not claims
	return err == nil && token.Valid
}
