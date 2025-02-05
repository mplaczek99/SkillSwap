package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/skillswap/models"
	"github.com/your-username/skillswap/services"
)

// AddSkill handles adding a new skill
func AddSkill(c *gin.Context) {
	var skill models.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newSkill, err := services.CreateSkill(&skill)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newSkill)
}

// GetSkills retrieves all skills
func GetSkills(c *gin.Context) {
	skills, err := services.GetAllSkills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, skills)
}

