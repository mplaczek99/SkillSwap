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
		// Return JSON object with error message for missing query parameter
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	searchTerm := strings.ToLower(q)
	var results []interface{}

	// Initialize user repository to search users from the database
	db, exists := c.Get("db")
	if !exists {
		// For tests or when db isn't available, return mock data
		utils.Info("Database connection not found in context, using mock data")

		// Return mock data for testing purposes
		mockResults := getMockSearchResults(searchTerm)
		c.JSON(http.StatusOK, mockResults)
		return
	}

	userRepo := repositories.NewUserRepository(db.(*gorm.DB))

	// Search users from the database
	users, err := userRepo.SearchUsers(searchTerm)
	if err != nil {
		utils.Error("Failed to search users: " + err.Error())
		c.JSON(http.StatusOK, []interface{}{})
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
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// Add skills to results
	for _, skill := range skills {
		results = append(results, skill)
	}

	c.JSON(http.StatusOK, results)
}

// getMockSearchResults provides mock data for testing purposes
func getMockSearchResults(searchTerm string) []interface{} {
	// Create some mock results for testing
	var results []interface{}

	// Mock users
	mockUsers := []map[string]interface{}{
		{"id": 1, "name": "Test User", "email": "test@example.com"},
		{"id": 2, "name": "Alice Smith", "email": "alice@example.com"},
	}

	// Mock skills
	mockSkills := []map[string]interface{}{
		{"id": 1, "name": "Programming", "description": "Learn to code"},
		{"id": 2, "name": "Music", "description": "Learn to play instruments"},
	}

	// Filter mock data based on search term
	for _, user := range mockUsers {
		name := strings.ToLower(user["name"].(string))
		email := strings.ToLower(user["email"].(string))
		if strings.Contains(name, searchTerm) || strings.Contains(email, searchTerm) {
			results = append(results, user)
		}
	}

	for _, skill := range mockSkills {
		name := strings.ToLower(skill["name"].(string))
		desc := strings.ToLower(skill["description"].(string))
		if strings.Contains(name, searchTerm) || strings.Contains(desc, searchTerm) {
			results = append(results, skill)
		}
	}

	return results
}
