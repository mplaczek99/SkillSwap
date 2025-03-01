package routes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/routes"
)

// MockAuthService is a mock implementation of the AuthServiceInterface
type MockAuthService struct{}

// Register is a mock implementation of Register with the correct signature
func (m *MockAuthService) Register(user *models.User) (string, error) {
	return "mock-token", nil
}

// Login is a mock implementation of Login
func (m *MockAuthService) Login(email, password string) (string, error) {
	if email == "test@example.com" && password == "password" {
		return "mock-login-token", nil
	}
	return "", errors.New("invalid email or password")
}

func TestRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Routes Configuration", func(t *testing.T) {
		// Create a new instance of the real AuthController with a mock service
		mockService := &MockAuthService{}
		authController := controllers.NewAuthController(mockService)

		// Setup router with the real controller
		router := gin.New()
		routes.SetupRoutes(router, authController)

		// Test auth endpoints
		testAuthEndpoints(t, router)

		// Test search endpoint
		testSearchEndpoint(t, router)

		// Test protected routes
		testProtectedRoutes(t, router)

		// Test static file routing
		testStaticFileRouting(t, router)
	})
}

func testAuthEndpoints(t *testing.T, router *gin.Engine) {
	// Test register endpoint
	t.Run("Register Endpoint", func(t *testing.T) {
		reqBody := `{"name": "Test User", "email": "new@example.com", "password": "password123"}`
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Response: %s", http.StatusCreated, w.Code, w.Body.String())
		}
	})

	// Test login endpoint
	t.Run("Login Endpoint", func(t *testing.T) {
		reqBody := `{"email": "test@example.com", "password": "password"}`
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if _, exists := response["token"]; !exists {
			t.Error("Expected token in response")
		}
	})
}

func testSearchEndpoint(t *testing.T, router *gin.Engine) {
	t.Run("Search Without Query", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/search", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Search With Query", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/search?q=test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func testProtectedRoutes(t *testing.T, router *gin.Engine) {
	// Test protected endpoint without token
	t.Run("Protected Endpoint Without Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/protected", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	// Test admin endpoint without token
	t.Run("Admin Endpoint Without Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/dashboard", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func testStaticFileRouting(t *testing.T, router *gin.Engine) {
	// Test accessing a static file path
	req, _ := http.NewRequest("GET", "/uploads/test.txt", nil)
	w := httptest.NewRecorder()

	// Note: This won't actually find a file, but should route correctly
	router.ServeHTTP(w, req)

	// We don't expect a 501 Not Implemented because the route exists
	if w.Code == http.StatusNotImplemented {
		t.Errorf("Static file route not implemented")
	}
}
