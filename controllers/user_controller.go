package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/skillswap/models"
	"github.com/your-username/skillswap/services"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	createdUser, err := services.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": createdUser})
}

