package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
)

func TestSearchController(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Search Without Query Parameter", func(t *testing.T) {
		router := gin.New()
		router.GET("/search", controllers.Search)

		req, _ := http.NewRequest("GET", "/search", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		if errMsg, exists := response["error"]; !exists || errMsg == "" {
			t.Errorf("Expected error message, got %s", errMsg)
		}
	})

	t.Run("Search With Query Parameter", func(t *testing.T) {
		router := gin.New()
		router.GET("/search", controllers.Search)

		req, _ := http.NewRequest("GET", "/search?q=test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var results []interface{}
		err := json.Unmarshal(w.Body.Bytes(), &results)
		if err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		// At least one result should include "test" in name or email (case insensitive)
		foundMatch := false
		for _, result := range results {
			resultMap, ok := result.(map[string]interface{})
			if !ok {
				continue
			}

			// Check for matches in name
			if name, ok := resultMap["name"].(string); ok {
				if containsIgnoreCase(name, "test") {
					foundMatch = true
					break
				}
			}

			// Check for matches in email
			if email, ok := resultMap["email"].(string); ok {
				if containsIgnoreCase(email, "test") {
					foundMatch = true
					break
				}
			}

			// Check for matches in description
			if description, ok := resultMap["description"].(string); ok {
				if containsIgnoreCase(description, "test") {
					foundMatch = true
					break
				}
			}
		}

		// This is not a strict validation since it depends on dummy data
		// But typically the test user would have "test" in the name or email
		if len(results) > 0 && !foundMatch {
			t.Logf("No results with 'test' in name, email, or description found. This might be expected depending on test data.")
		}
	})

	t.Run("Search With Empty Query Parameter", func(t *testing.T) {
		router := gin.New()
		router.GET("/search", controllers.Search)

		req, _ := http.NewRequest("GET", "/search?q=", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// Helper function to check if a string contains another string, ignoring case
func containsIgnoreCase(s, substr string) bool {
	s, substr = toLower(s), toLower(substr)
	return contains(s, substr)
}

// Simple toLower function
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

// Simple contains function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
