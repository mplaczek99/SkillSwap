package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/repositories"
)

// Search handles GET requests to search for skills and users.
// It expects a query parameter "q" and returns matching results.
func Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	var results []interface{}

	// Search dummy skills from the repository.
	skills, err := repositories.GetAllSkills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch skills"})
		return
	}

	// Filter skills by name or description.
	for _, skill := range skills {
		if strings.Contains(strings.ToLower(skill.Name), strings.ToLower(q)) ||
			strings.Contains(strings.ToLower(skill.Description), strings.ToLower(q)) {
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
	if strings.Contains(strings.ToLower(dummyUser.Name), strings.ToLower(q)) ||
		strings.Contains(strings.ToLower(dummyUser.Email), strings.ToLower(q)) {
		results = append(results, dummyUser)
	}

	c.JSON(http.StatusOK, results)
}

