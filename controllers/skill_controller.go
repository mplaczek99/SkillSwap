package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/services"
	"github.com/mplaczek99/SkillSwap/utils"
)

// AddSkill handles adding a new skill.
func AddSkill(c *gin.Context) {
	var skill models.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	newSkill, err := services.CreateSkill(&skill)
	if err != nil {
		if strings.Contains(err.Error(), "required") {
			utils.JSONError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, newSkill)
}

// GetSkills retrieves all skills.
func GetSkills(c *gin.Context) {
	skills, err := services.GetAllSkills()
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, skills)
}
