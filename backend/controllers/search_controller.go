package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/repositories"
	"github.com/mplaczek99/SkillSwap/utils"
)

// Search handles GET requests to search for skills and users.
// It expects a query parameter "q" and returns matching results.
func Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		utils.JSONError(c, http.StatusBadRequest, "Query parameter 'q' is required")
		return
	}

	searchTerm := strings.ToLower(q)
	var results []interface{}

	// Search skills from the repository.
	skills, err := repositories.GetAllSkills()
	if err != nil {
		utils.Error("Failed to fetch skills: " + err.Error())
		utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch skills")
		return
	}

	// Filter skills by name or description.
	for _, skill := range skills {
		if strings.Contains(strings.ToLower(skill.Name), searchTerm) ||
			strings.Contains(strings.ToLower(skill.Description), searchTerm) {
			results = append(results, skill)
		}
	}

	// Also include a dummy user if the query matches.
	dummyUser := models.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
		Bio:   "This is a dummy user",
	}
	if strings.Contains(strings.ToLower(dummyUser.Name), searchTerm) ||
		strings.Contains(strings.ToLower(dummyUser.Email), searchTerm) {
		results = append(results, dummyUser)
	}

	c.JSON(http.StatusOK, results)
}
