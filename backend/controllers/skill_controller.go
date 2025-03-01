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
// @Summary Add a new skill
// @Description Create a new skill entry
// @Tags skills
// @Accept json
// @Produce json
// @Param skill body models.Skill true "Skill information"
// @Success 201 {object} models.Skill "Successfully created skill"
// @Failure 400 {object} map[string]string "Invalid skill data"
// @Failure 500 {object} map[string]string "Failed to add skill"
// @Security BearerAuth
// @Router /skills [post]
func AddSkill(c *gin.Context) {
	var skill models.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "Invalid skill data")
		return
	}

	newSkill, err := services.CreateSkill(&skill)
	if err != nil {
		if strings.Contains(err.Error(), "required") {
			utils.JSONError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.JSONError(c, http.StatusInternalServerError, "Failed to add skill")
		return
	}
	c.JSON(http.StatusCreated, newSkill)
}

// GetSkills retrieves all skills.
// @Summary Get all skills
// @Description Retrieve a list of all available skills
// @Tags skills
// @Accept json
// @Produce json
// @Success 200 {array} models.Skill "List of skills"
// @Failure 500 {object} map[string]string "Failed to retrieve skills"
// @Router /skills [get]
func GetSkills(c *gin.Context) {
	skills, err := services.GetAllSkills()
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "Failed to retrieve skills")
		return
	}
	c.JSON(http.StatusOK, skills)
}
