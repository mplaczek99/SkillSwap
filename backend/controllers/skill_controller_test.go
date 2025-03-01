package controllers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
)

func TestAddSkill(t *testing.T) {
	router := gin.Default()
	router.POST("/skills", controllers.AddSkill)

	tests := []struct {
		name           string
		requestBody    string
		wantStatusCode int
		wantContains   string
	}{
		{
			name:           "Valid Skill",
			requestBody:    `{"name":"Go Programming","description":"Basics of Go","user_id":1}`,
			wantStatusCode: http.StatusCreated,
			wantContains:   `"name":"Go Programming"`,
		},
		{
			name:           "Missing Skill Name",
			requestBody:    `{"description":"Basics of Go","user_id":1}`,
			wantStatusCode: http.StatusBadRequest,
			wantContains:   `"error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/skills", bytes.NewBuffer([]byte(tt.requestBody)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("[%s] expected status %d, got %d", tt.name, tt.wantStatusCode, w.Code)
			}
			if tt.wantContains != "" && !bytes.Contains(w.Body.Bytes(), []byte(tt.wantContains)) {
				t.Errorf("[%s] body does not contain %q. Body: %s", tt.name, tt.wantContains, w.Body.String())
			}
		})
	}
}

func TestGetSkills(t *testing.T) {
	router := gin.Default()
	router.GET("/skills", controllers.GetSkills)

	req, _ := http.NewRequest("GET", "/skills", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Dummy Skill")) {
		t.Errorf("response does not contain dummy skill. Body: %s", w.Body.String())
	}
}
