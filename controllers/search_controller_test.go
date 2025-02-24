package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/api/search", controllers.Search)
	return router
}

func TestSearchEndpoint_NoQuery(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/api/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing query, got %d", w.Code)
	}
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to unmarshal response")
	}
	if !strings.Contains(response["error"], "Query parameter") {
		t.Errorf("Expected error message about missing query, got %v", response["error"])
	}
}

func TestSearchEndpoint_MatchSkill(t *testing.T) {
	router := setupRouter()
	// "dummy" should match the dummy skill "Dummy Skill"
	req, _ := http.NewRequest("GET", "/api/search?q=dummy", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var results []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
		t.Fatal("Failed to unmarshal response")
	}

	found := false
	for _, r := range results {
		if name, ok := r["name"].(string); ok && strings.Contains(strings.ToLower(name), "dummy") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected to find a skill containing 'dummy' in response")
	}
}

func TestSearchEndpoint_MatchUser(t *testing.T) {
	router := setupRouter()
	// "test" should match the dummy user "Test User"
	req, _ := http.NewRequest("GET", "/api/search?q=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var results []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
		t.Fatal("Failed to unmarshal response")
	}

	found := false
	for _, r := range results {
		if email, ok := r["email"].(string); ok && email == "test@example.com" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected to find a user with email 'test@example.com' in response")
	}
}
