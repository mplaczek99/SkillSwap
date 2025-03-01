package controllers_test

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
)

// MockAuthService is a mock implementation of the AuthServiceInterface
type MockAuthService struct {
	// Maps to store expected returns for specific inputs
	registerResponses map[string]registerResponse
	loginResponses    map[string]loginResponse
}

type registerResponse struct {
	token string
	err   error
}

type loginResponse struct {
	token string
	err   error
}

func NewMockAuthService() *MockAuthService {
	return &MockAuthService{
		registerResponses: make(map[string]registerResponse),
		loginResponses:    make(map[string]loginResponse),
	}
}

// Register mocks the Register method
func (m *MockAuthService) Register(user *models.User) (string, error) {
	resp, exists := m.registerResponses[user.Email]
	if !exists {
		return "", errors.New("unexpected email in test")
	}
	return resp.token, resp.err
}

// Login mocks the Login method
func (m *MockAuthService) Login(email, password string) (string, error) {
	key := email + ":" + password
	resp, exists := m.loginResponses[key]
	if !exists {
		return "", errors.New("unexpected credentials in test")
	}
	return resp.token, resp.err
}

// SetupRegisterResponse sets up an expected response for Register
func (m *MockAuthService) SetupRegisterResponse(email string, token string, err error) {
	m.registerResponses[email] = registerResponse{token: token, err: err}
}

// SetupLoginResponse sets up an expected response for Login
func (m *MockAuthService) SetupLoginResponse(email, password string, token string, err error) {
	key := email + ":" + password
	m.loginResponses[key] = loginResponse{token: token, err: err}
}

func TestAuthController_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := NewMockAuthService()

	// Set up expected responses
	mockService.SetupRegisterResponse("success@example.com", "valid-token", nil)
	mockService.SetupRegisterResponse("error@example.com", "", errors.New("registration failed"))

	controller := controllers.NewAuthController(mockService)

	t.Run("Successful Registration", func(t *testing.T) {
		router := gin.New()
		router.POST("/register", controller.Register)

		// Create request body
		reqBody := `{
			"name": "Test User",
			"email": "success@example.com",
			"password": "password123"
		}`

		req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		if token, exists := response["token"]; !exists || token != "valid-token" {
			t.Errorf("Expected token 'valid-token', got %s", token)
		}
	})

	t.Run("Failed Registration", func(t *testing.T) {
		router := gin.New()
		router.POST("/register", controller.Register)

		// Create request body
		reqBody := `{
			"name": "Error User",
			"email": "error@example.com",
			"password": "password123"
		}`

		req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		if errMsg, exists := response["error"]; !exists || errMsg != "Registration failed" {
			t.Errorf("Expected error message 'Registration failed', got %s", errMsg)
		}
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		router := gin.New()
		router.POST("/register", controller.Register)

		// Create invalid request body (missing required fields)
		reqBody := `{
			"name": "Invalid User"
		}`

		req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestAuthController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := NewMockAuthService()

	// Set up expected responses
	mockService.SetupLoginResponse("success@example.com", "password123", "valid-login-token", nil)
	mockService.SetupLoginResponse("error@example.com", "password123", "", errors.New("invalid credentials"))

	controller := controllers.NewAuthController(mockService)

	t.Run("Successful Login", func(t *testing.T) {
		router := gin.New()
		router.POST("/login", controller.Login)

		// Create request body
		reqBody := `{
			"email": "success@example.com",
			"password": "password123"
		}`

		req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		if token, exists := response["token"]; !exists || token != "valid-login-token" {
			t.Errorf("Expected token 'valid-login-token', got %s", token)
		}
	})

	t.Run("Failed Login", func(t *testing.T) {
		router := gin.New()
		router.POST("/login", controller.Login)

		// Create request body
		reqBody := `{
			"email": "error@example.com",
			"password": "password123"
		}`

		req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		if errMsg, exists := response["error"]; !exists || errMsg != "Invalid email or password" {
			t.Errorf("Expected error message 'Invalid email or password', got %s", errMsg)
		}
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		router := gin.New()
		router.POST("/login", controller.Login)

		// Create invalid request body (missing required fields)
		reqBody := `{
			"email": "missing@example.com"
		}`

		req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
