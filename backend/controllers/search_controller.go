package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/repositories"
	"github.com/mplaczek99/SkillSwap/utils"
	"gorm.io/gorm"
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

	// Initialize user repository to search users from the database
	db, exists := c.Get("db")
	if !exists {
		// If db isn't in the context, we can fallback to the database connection
		// or report an error
		utils.Error("Database connection not found in context")
		utils.JSONError(c, http.StatusInternalServerError, "Database connection error")
		return
	}

	userRepo := repositories.NewUserRepository(db.(*gorm.DB))

	// Search users from the database
	users, err := userRepo.SearchUsers(searchTerm)
	if err != nil {
		utils.Error("Failed to search users: " + err.Error())
		utils.JSONError(c, http.StatusInternalServerError, "Failed to search users")
		return
	}

	// Add users to results
	for _, user := range users {
		results = append(results, user)
	}

	// Search skills from the repository
	skills, err := repositories.SearchSkills(searchTerm)
	if err != nil {
		utils.Error("Failed to search skills: " + err.Error())
		utils.JSONError(c, http.StatusInternalServerError, "Failed to search skills")
		return
	}

	// Add skills to results
	for _, skill := range skills {
		results = append(results, skill)
	}

	c.JSON(http.StatusOK, results)
}
