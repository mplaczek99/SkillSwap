package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/controllers"
	"github.com/mplaczek99/SkillSwap/models"
)

func TestCreateSchedule(t *testing.T) {
	// Use test mode
	gin.SetMode(gin.TestMode)

	t.Run("Create Valid Schedule", func(t *testing.T) {
		// Set up router with the controller
		router := gin.New()
		router.POST("/schedule", controllers.CreateSchedule)

		// Create valid schedule data
		startTime := time.Now().Add(24 * time.Hour) // 1 day in future
		endTime := startTime.Add(2 * time.Hour)     // 2 hours duration

		schedule := models.Schedule{
			UserID:    1,
			SkillID:   2,
			StartTime: startTime,
			EndTime:   endTime,
		}

		reqBody, _ := json.Marshal(schedule)

		// Create request
		req, _ := http.NewRequest("POST", "/schedule", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201 for valid schedule, got %d: %s", w.Code, w.Body.String())
		}

		// Verify response contains created schedule
		var response models.Schedule
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		if response.ID == 0 {
			t.Error("Expected ID to be set in response")
		}
	})

	t.Run("Create Schedule With Past Start Time", func(t *testing.T) {
		// Set up router with the controller
		router := gin.New()
		router.POST("/schedule", controllers.CreateSchedule)

		// Create schedule data with past start time
		startTime := time.Now().Add(-24 * time.Hour) // 1 day in past
		endTime := startTime.Add(2 * time.Hour)      // 2 hours duration

		schedule := models.Schedule{
			UserID:    1,
			SkillID:   2,
			StartTime: startTime,
			EndTime:   endTime,
		}

		reqBody, _ := json.Marshal(schedule)

		// Create request
		req, _ := http.NewRequest("POST", "/schedule", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for past start time, got %d", w.Code)
		}
	})

	t.Run("Create Schedule With End Time Before Start Time", func(t *testing.T) {
		// Set up router with the controller
		router := gin.New()
		router.POST("/schedule", controllers.CreateSchedule)

		// Create schedule data with end time before start time
		startTime := time.Now().Add(24 * time.Hour) // 1 day in future
		endTime := startTime.Add(-1 * time.Hour)    // 1 hour before start time

		schedule := models.Schedule{
			UserID:    1,
			SkillID:   2,
			StartTime: startTime,
			EndTime:   endTime,
		}

		reqBody, _ := json.Marshal(schedule)

		// Create request
		req, _ := http.NewRequest("POST", "/schedule", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for end time before start time, got %d", w.Code)
		}
	})

	t.Run("Create Schedule With Invalid JSON", func(t *testing.T) {
		// Set up router with the controller
		router := gin.New()
		router.POST("/schedule", controllers.CreateSchedule)

		// Create invalid JSON
		invalidJSON := `{"user_id": 1, "skill_id": 2, "start_time": "invalid-date"}`

		// Create request
		req, _ := http.NewRequest("POST", "/schedule", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid JSON, got %d", w.Code)
		}
	})
}

func TestGetSchedules(t *testing.T) {
	// Use test mode
	gin.SetMode(gin.TestMode)

	t.Run("Get Schedules Without Authentication", func(t *testing.T) {
		// Set up router with the controller
		router := gin.New()
		router.GET("/schedule", controllers.GetSchedules)

		// Create request
		req, _ := http.NewRequest("GET", "/schedule", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response - should fail without user_id in context
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for unauthenticated request, got %d", w.Code)
		}
	})

	t.Run("Get Schedules With Authentication", func(t *testing.T) {
		// Set up router with the controller and middleware to set user_id
		router := gin.New()
		router.GET("/schedule", func(c *gin.Context) {
			// Set user_id in context to simulate authentication
			c.Set("user_id", uint(1))
			controllers.GetSchedules(c)
		})

		// Create request
		req, _ := http.NewRequest("GET", "/schedule", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Verify response
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for authenticated request, got %d", w.Code)
		}

		// Verify response contains schedules
		var schedules []models.Schedule
		if err := json.Unmarshal(w.Body.Bytes(), &schedules); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		if len(schedules) == 0 {
			t.Error("Expected at least one schedule in response")
		}
	})
}
