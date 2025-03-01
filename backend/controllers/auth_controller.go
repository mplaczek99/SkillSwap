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
// @Summary Register a new user
// @Description Create a new user account and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User Registration Info"
// @Success 201 {object} RegisterResponse "Successfully registered"
// @Failure 400 {object} map[string]string "Invalid registration data"
// @Failure 500 {object} map[string]string "Registration failed"
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "Invalid registration data")
		return
	}
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	token, err := c.AuthService.Register(user)
	if err != nil {
		utils.Error("Registration failed for " + req.Email + ": " + err.Error())
		utils.JSONError(ctx, http.StatusInternalServerError, "Registration failed")
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
// @Summary Log in a user
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User Login Credentials"
// @Success 200 {object} LoginResponse "Successfully logged in"
// @Failure 400 {object} map[string]string "Invalid login data"
// @Failure 401 {object} map[string]string "Invalid email or password"
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "Invalid login data")
		return
	}
	token, err := c.AuthService.Login(req.Email, req.Password)
	if err != nil {
		utils.Error("Login failed for " + req.Email + ": " + err.Error())
		utils.JSONError(ctx, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	ctx.JSON(http.StatusOK, LoginResponse{Token: token})
}
