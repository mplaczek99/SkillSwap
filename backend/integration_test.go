package integration_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/middleware"
	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/utils"
)

// This test file contains integration tests that test the interaction between
// different components of the system. These tests focus on API endpoints
// that require multiple layers to work together.

// TestMain sets up environment for tests in the main package
func TestMain(m *testing.M) {
	// Set JWT_SECRET for tests
	os.Setenv("JWT_SECRET", "test-secret-for-integration-tests")

	// Run the tests
	code := m.Run()

	// Exit with the appropriate code
	os.Exit(code)
}

// Helper function to set up test router
func setupTestRouter() *gin.Engine {
	router := gin.New()

	// Add basic group for API
	api := router.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	{
		// Create a new auth controller with a real auth service
		authService := &mockAuthService{}
		authController := controllers.NewAuthController(authService)

		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "You are authenticated"})
		})

		// Schedule routes
		protected.POST("/schedule", controllers.CreateSchedule)
		protected.GET("/schedule", controllers.GetSchedules)
	}

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Admin dashboard"})
		})
	}

	return router
}

// Mock auth service for testing
type mockAuthService struct{}

func (m *mockAuthService) Register(user *models.User) (string, error) {
	// For testing, just generate a token
	return utils.GenerateToken(1, user.Role, user.Email)
}

func (m *mockAuthService) Login(email, password string) (string, error) {
	// For testing, generate a token if password is correct
	if password == "password123" {
		return utils.GenerateToken(1, "User", email)
	}
	// Return an error for wrong password
	return "", errors.New("invalid email or password")
}

func TestAuthAndProtectedEndpoints(t *testing.T) {
	// Use test mode
	gin.SetMode(gin.TestMode)

	// Set up router with controllers and middleware
	router := setupTestRouter()

	// Create a test user
	user := models.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123", // Will be hashed by BeforeSave
		Role:     "User",
	}

	// Call BeforeSave to hash the password
	err := user.BeforeSave(nil)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	t.Run("Access Protected Route Without Token", func(t *testing.T) {
		// Create request to a protected route
		req, _ := http.NewRequest("GET", "/api/protected", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for protected route without token, got %d", w.Code)
		}
	})

	t.Run("Access Protected Route With Invalid Token", func(t *testing.T) {
		// Create request with invalid token
		req, _ := http.NewRequest("GET", "/api/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for protected route with invalid token, got %d", w.Code)
		}
	})

	t.Run("Access Protected Route With Valid Token", func(t *testing.T) {
		// Generate token for the test user
		token, err := utils.GenerateToken(user.ID, user.Role, user.Email)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Create request with valid token
		req, _ := http.NewRequest("GET", "/api/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for protected route with valid token, got %d", w.Code)
		}
	})

	t.Run("Access Admin Route With User Token", func(t *testing.T) {
		// Generate token for a regular user
		token, err := utils.GenerateToken(user.ID, "User", user.Email)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Create request to admin route with user token
		req, _ := http.NewRequest("GET", "/api/admin/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 for admin route with user token, got %d", w.Code)
		}
	})

	t.Run("Access Admin Route With Admin Token", func(t *testing.T) {
		// Generate token for an admin user
		token, err := utils.GenerateToken(user.ID, "Admin", user.Email)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Create request to admin route with admin token
		req, _ := http.NewRequest("GET", "/api/admin/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for admin route with admin token, got %d", w.Code)
		}
	})
}

func TestAuthentication_Register_Login_Flow(t *testing.T) {
	// Use test mode
	gin.SetMode(gin.TestMode)

	// Set up router with auth controller and services
	router := setupTestRouter()

	// Email to use for the test
	testEmail := "newuser@example.com"
	testPassword := "password123"

	t.Run("Register New User", func(t *testing.T) {
		// Create registration data
		regData := map[string]string{
			"name":     "New User",
			"email":    testEmail,
			"password": testPassword,
		}
		reqBody, _ := json.Marshal(regData)

		// Create request
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201 for registration, got %d: %s", w.Code, w.Body.String())
		}

		// Check response contains token
		var response map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if token, exists := response["token"]; !exists || token == "" {
			t.Error("Expected token in response")
		}
	})

	t.Run("Login With New User", func(t *testing.T) {
		// Create login data
		loginData := map[string]string{
			"email":    testEmail,
			"password": testPassword,
		}
		reqBody, _ := json.Marshal(loginData)

		// Create request
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for login, got %d: %s", w.Code, w.Body.String())
		}

		// Check response contains token
		var response map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if token, exists := response["token"]; !exists || token == "" {
			t.Error("Expected token in response")
		}
	})

	t.Run("Login With Wrong Password", func(t *testing.T) {
		// Create login data with wrong password
		loginData := map[string]string{
			"email":    testEmail,
			"password": "wrongpassword",
		}
		reqBody, _ := json.Marshal(loginData)

		// Create request
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for login with wrong password, got %d", w.Code)
		}
	})
}

func TestScheduleWithAuthentication(t *testing.T) {
	// Use test mode
	gin.SetMode(gin.TestMode)

	// Set up router with auth and schedule controllers
	router := setupTestRouter()

	// Generate token for a test user
	token, err := utils.GenerateToken(1, "User", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	t.Run("Create Schedule With Authentication", func(t *testing.T) {
		// Create valid schedule data
		startTime := time.Now().Add(24 * time.Hour) // 1 day in future
		endTime := startTime.Add(2 * time.Hour)     // 2 hours duration

		schedule := models.Schedule{
			UserID:    1,
			SkillID:   2,
			StartTime: startTime,
			EndTime:   endTime,
		}

		reqBody, _ := json.Marshal(schedule)

		// Create request with token
		req, _ := http.NewRequest("POST", "/api/schedule", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201 for schedule creation with auth, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Get Schedules With Authentication", func(t *testing.T) {
		// Create request with token
		req, _ := http.NewRequest("GET", "/api/schedule", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for get schedules with auth, got %d", w.Code)
		}

		// Verify response contains schedules
		var schedules []models.Schedule
		if err := json.Unmarshal(w.Body.Bytes(), &schedules); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		if len(schedules) == 0 {
			t.Error("Expected at least one schedule in response")
		}
	})
}
