package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/services"
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
// @Summary Register a new user
// @Description Create a new user and return a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration details"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input"})
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
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to register user"})
		return
	}

	ctx.JSON(http.StatusCreated, RegisterResponse{Token: token})
}

// Login godoc
// @Summary Login a user
// @Description Authenticate user and return a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid credentials"})
		return
	}

	token, err := c.AuthService.Login(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid email or password"})
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{Token: token})
}
