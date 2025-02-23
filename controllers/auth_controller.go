package controllers

import (
	"net/http"

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

// RegisterRequest defines the required fields for registration
type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
}

// RegisterResponse defines the response after successful registration
type RegisterResponse struct {
	Token string `json:"token" example:"your.jwt.token"`
}

// LoginRequest defines the required fields for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

// LoginResponse defines the response after successful login
type LoginResponse struct {
	Token string `json:"token" example:"your.jwt.token"`
}

// ErrorResponse defines the error structure
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}

// Register godoc
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // Will be hashed in service layer
		Role:     "User",       // Default role
	}

	token, err := c.AuthService.Register(user)
	if err != nil {
		utils.JSONError(ctx, http.StatusInternalServerError, "Failed to register user")
		return
	}

	ctx.JSON(http.StatusCreated, RegisterResponse{Token: token})
}

// Login godoc
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "Invalid credentials")
		return
	}

	token, err := c.AuthService.Login(req.Email, req.Password)
	if err != nil {
		utils.JSONError(ctx, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{Token: token})
}
