package controllers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/models"
)

// DummyAuthService is a mock implementation of the AuthServiceInterface.
type DummyAuthService struct{}

func (d *DummyAuthService) Register(user *models.User) (string, error) {
	if user.Email == "" || user.Password == "" {
		return "", errors.New("invalid input")
	}
	return "dummy_register_token", nil
}

func (d *DummyAuthService) Login(email, password string) (string, error) {
	if email == "test@example.com" && password == "somepassword" {
		return "dummy_login_token", nil
	}
	return "", errors.New("invalid email or password")
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Create a dummy AuthController with DummyAuthService
	authController := controllers.NewAuthController(&DummyAuthService{})

	router := gin.Default()
	router.POST("/login", authController.Login)

	tests := []struct {
		name           string
		requestBody    string
		wantStatusCode int
		wantContains   string
	}{
		{
			name:           "Valid Credentials",
			requestBody:    `{"email":"test@example.com","password":"somepassword"}`,
			wantStatusCode: http.StatusOK,
			wantContains:   `"token":"dummy_login_token"`,
		},
		{
			name:           "Missing Fields",
			requestBody:    `{"email":""}`,
			wantStatusCode: http.StatusBadRequest,
			wantContains:   `"error"`,
		},
		{
			name:           "Invalid Email Format",
			requestBody:    `{"email":"notanemail","password":"secret"}`,
			wantStatusCode: http.StatusBadRequest,
			wantContains:   `"error"`,
		},
		{
			name:           "User not found",
			requestBody:    `{"email":"nope@example.com","password":"secret"}`,
			wantStatusCode: http.StatusUnauthorized,
			wantContains:   `"error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(tt.requestBody)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("[%s] expected status %d, got %d", tt.name, tt.wantStatusCode, w.Code)
			}
			if tt.wantContains != "" && !bytes.Contains(w.Body.Bytes(), []byte(tt.wantContains)) {
				t.Errorf("[%s] response body does not contain %q. Body: %s", tt.name, tt.wantContains, w.Body.String())
			}
		})
	}
}
