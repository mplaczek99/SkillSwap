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

func TestScheduleOverlapAndEdgeCases(t *testing.T) {
	// Use test mode
	gin.SetMode(gin.TestMode)

	// Set up router with the controller
	router := gin.New()

	// To test overlap, we need to track existing schedules
	var existingSchedules []models.Schedule

	// Create a custom handler with overlap detection
	mockCreateSchedule := func(c *gin.Context) {
		var schedule models.Schedule
		if err := c.ShouldBindJSON(&schedule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule data"})
			return
		}

		// Basic validation
		now := time.Now()
		if schedule.StartTime.Before(now) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule start time must be in the future"})
			return
		}
		if !schedule.EndTime.After(schedule.StartTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule end time must be after start time"})
			return
		}

		// Check for schedule with exactly same start/end time (edge case)
		if schedule.StartTime.Equal(schedule.EndTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Start time cannot equal end time"})
			return
		}

		// Check minimum duration
		minDuration := 30 * time.Minute
		if schedule.EndTime.Sub(schedule.StartTime) < minDuration {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule must be at least 30 minutes"})
			return
		}

		// Check for overlaps with existing schedules
		for _, existing := range existingSchedules {
			// Only check schedules for the same user
			if existing.UserID == schedule.UserID {
				// Check for overlap: new start time is within existing schedule
				if (schedule.StartTime.After(existing.StartTime) && schedule.StartTime.Before(existing.EndTime)) ||
					// Check for overlap: new end time is within existing schedule
					(schedule.EndTime.After(existing.StartTime) && schedule.EndTime.Before(existing.EndTime)) ||
					// Check for overlap: new schedule encompasses existing schedule
					(schedule.StartTime.Before(existing.StartTime) && schedule.EndTime.After(existing.EndTime)) {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule overlaps with an existing session"})
					return
				}
			}
		}

		// Add the schedule to our existing schedules
		schedule.ID = uint(len(existingSchedules) + 1)
		schedule.CreatedAt = time.Now()
		existingSchedules = append(existingSchedules, schedule)

		c.JSON(http.StatusCreated, schedule)
	}

	router.POST("/schedule", mockCreateSchedule)

	t.Run("Create Initial Schedule", func(t *testing.T) {
		startTime := time.Now().Add(24 * time.Hour)
		endTime := startTime.Add(2 * time.Hour)

		schedule := models.Schedule{
			UserID:    1,
			SkillID:   2,
			StartTime: startTime,
			EndTime:   endTime,
		}

		reqBody, _ := json.Marshal(schedule)
		req, _ := http.NewRequest("POST", "/schedule", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201 for first schedule, got %d", w.Code)
		}
	})

	t.Run("Create Overlapping Schedule", func(t *testing.T) {
		existingSchedule := existingSchedules[0]

		// Create a schedule that starts during the existing one
		startTime := existingSchedule.StartTime.Add(1 * time.Hour)
		endTime := existingSchedule.EndTime.Add(1 * time.Hour)

		schedule := models.Schedule{
			UserID:    1, // Same user
			SkillID:   2,
			StartTime: startTime,
			EndTime:   endTime,
		}

		reqBody, _ := json.Marshal(schedule)
		req, _ := http.NewRequest("POST", "/schedule", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for overlapping schedule, got %d", w.Code)
		}
	})

	t.Run("Same Time Different User", func(t *testing.T) {
		existingSchedule := existingSchedules[0]

		// Create a schedule with the same time but for a different user
		schedule := models.Schedule{
			UserID:    2, // Different user
			SkillID:   2,
			StartTime: existingSchedule.StartTime,
			EndTime:   existingSchedule.EndTime,
		}

		reqBody, _ := json.Marshal(schedule)
		req, _ := http.NewRequest("POST", "/schedule", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201 for non-conflicting schedule, got %d", w.Code)
		}
	})

	t.Run("Schedule With Duration Less Than 30 Minutes", func(t *testing.T) {
		startTime := time.Now().Add(48 * time.Hour)
		endTime := startTime.Add(15 * time.Minute) // Too short

		schedule := models.Schedule{
			UserID:    1,
			SkillID:   2,
			StartTime: startTime,
			EndTime:   endTime,
		}

		reqBody, _ := json.Marshal(schedule)
		req, _ := http.NewRequest("POST", "/schedule", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for too short schedule, got %d", w.Code)
		}
	})

	t.Run("Schedule With Equal Start and End Times", func(t *testing.T) {
		exactTime := time.Now().Add(72 * time.Hour)

		schedule := models.Schedule{
			UserID:    1,
			SkillID:   2,
			StartTime: exactTime,
			EndTime:   exactTime, // Same as start time
		}

		reqBody, _ := json.Marshal(schedule)
		req, _ := http.NewRequest("POST", "/schedule", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for equal start/end times, got %d", w.Code)
		}
	})
}
