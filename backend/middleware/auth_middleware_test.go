package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mplaczek99/SkillSwap/middleware"
	"github.com/mplaczek99/SkillSwap/utils"
)

func TestAuthMiddleware(t *testing.T) {
	// Set a consistent JWT secret for tests
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test_secret_key")
	defer os.Setenv("JWT_SECRET", originalSecret)

	gin.SetMode(gin.TestMode)

	t.Run("No Authorization Header", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.AuthMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "You are authenticated")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for missing token, got %d", w.Code)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.AuthMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "You are authenticated")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for invalid token, got %d", w.Code)
		}
	})

	t.Run("Valid Token", func(t *testing.T) {
		// Generate a valid token for testing
		token, err := utils.GenerateToken(123, "User", "test@example.com")
		if err != nil {
			t.Fatalf("Failed to generate token for testing: %v", err)
		}

		router := gin.New()
		router.Use(middleware.AuthMiddleware())
		router.GET("/test", func(c *gin.Context) {
			// Check that user data was set in context correctly
			userID, exists := c.Get("user_id")
			if !exists || userID.(uint) != 123 {
				t.Errorf("Expected user_id to be 123, got %v", userID)
			}

			role, exists := c.Get("role")
			if !exists || role.(string) != "User" {
				t.Errorf("Expected role to be 'User', got %v", role)
			}

			email, exists := c.Get("email")
			if !exists || email.(string) != "test@example.com" {
				t.Errorf("Expected email to be 'test@example.com', got %v", email)
			}

			c.String(http.StatusOK, "You are authenticated")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for valid token, got %d", w.Code)
		}
	})

	t.Run("Expired Token", func(t *testing.T) {
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

		expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := expiredToken.SignedString([]byte("test_secret_key"))

		router := gin.New()
		router.Use(middleware.AuthMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "You are authenticated")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for expired token, got %d", w.Code)
		}
	})

	t.Run("Token Without Bearer Prefix", func(t *testing.T) {
		token, err := utils.GenerateToken(123, "User", "test@example.com")
		if err != nil {
			t.Fatalf("Failed to generate token for testing: %v", err)
		}

		router := gin.New()
		router.Use(middleware.AuthMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "You are authenticated")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", token) // No "Bearer " prefix
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for token without Bearer prefix, got %d", w.Code)
		}
	})
}
