package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/services"
	"github.com/mplaczek99/SkillSwap/utils"
)

type AuthController struct {
	AuthService services.AuthServiceInterface
}

func NewAuthController(authService services.AuthServiceInterface) *AuthController {
	return &AuthController{AuthService: authService}
}

// RegisterRequest defines the required fields for registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterResponse defines the response after successful registration.
type RegisterResponse struct {
	Token string `json:"token"`
}

// Register handles user registration.
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Enhanced error handling with field validation details
		errorMsg := "Invalid registration data"
		if gin.IsDebugging() {
			errorMsg = fmt.Sprintf("Invalid registration data: %v", err)
		}

		utils.Error(fmt.Sprintf("Registration validation failed: %v", err))
		utils.JSONError(ctx, http.StatusBadRequest, errorMsg)
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	token, err := c.AuthService.Register(user)
	if err != nil {
		// Log detailed error for server logs
		utils.Error(fmt.Sprintf("Registration failed for %s: %v", req.Email, err))

		// Prepare appropriate client response based on error type
		statusCode := http.StatusInternalServerError
		errorMsg := "Registration failed"

		// Provide more context in development mode or for specific errors
		if isDevelopmentMode() || isClientVisibleError(err) {
			// Check for specific error types for better context
			if strings.Contains(err.Error(), "email already in use") {
				statusCode = http.StatusConflict
				errorMsg = "Email address is already in use"
			} else if strings.Contains(err.Error(), "database") {
				errorMsg = "Registration failed due to database error"
			} else if strings.Contains(err.Error(), "validation") {
				statusCode = http.StatusBadRequest
				errorMsg = err.Error()
			}
		}

		utils.JSONError(ctx, statusCode, errorMsg)
		return
	}

	ctx.JSON(http.StatusCreated, RegisterResponse{Token: token})
}

// LoginRequest defines the required fields for login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse defines the response after successful login.
type LoginResponse struct {
	Token string `json:"token"`
}

// Login handles user login.
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Enhanced error with validation details
		errorMsg := "Invalid login data"
		if gin.IsDebugging() {
			errorMsg = fmt.Sprintf("Invalid login data: %v", err)
		}

		utils.Error(fmt.Sprintf("Login validation failed: %v", err))
		utils.JSONError(ctx, http.StatusBadRequest, errorMsg)
		return
	}

	token, err := c.AuthService.Login(req.Email, req.Password)
	if err != nil {
		// Log detailed error for server logs
		utils.Error(fmt.Sprintf("Login failed for %s: %v", req.Email, err))

		// Prepare contextual error message
		statusCode := http.StatusUnauthorized
		errorMsg := "Invalid email or password"

		// Provide more specific errors in development mode
		if isDevelopmentMode() {
			if strings.Contains(err.Error(), "not found") {
				errorMsg = "User not found with provided email"
			} else if strings.Contains(err.Error(), "password") {
				errorMsg = "Incorrect password"
			} else if strings.Contains(err.Error(), "database") {
				statusCode = http.StatusInternalServerError
				errorMsg = "Authentication failed due to database error"
			}
		}

		utils.JSONError(ctx, statusCode, errorMsg)
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{Token: token})
}

// Helper functions

// isDevelopmentMode returns true if the application is running in development mode
func isDevelopmentMode() bool {
	return os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "dev" || os.Getenv("GIN_MODE") != "release"
}

// isClientVisibleError determines if an error should be exposed to clients
// even in production (like validation errors)
func isClientVisibleError(err error) bool {
	if err == nil {
		return false
	}

	// List of error types that are safe to show to clients
	safeErrorPrefixes := []string{
		"validation:",
		"input:",
		"email already in use",
	}

	errMsg := err.Error()
	for _, prefix := range safeErrorPrefixes {
		if strings.HasPrefix(errMsg, prefix) || strings.Contains(errMsg, prefix) {
			return true
		}
	}

	return false
}
