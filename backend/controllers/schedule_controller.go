package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/repositories"
)

// CreateSchedule handles scheduling a new session.
func CreateSchedule(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule data"})
		return
	}

	// Validate that the session is scheduled for the future.
	now := time.Now()
	if schedule.StartTime.Before(now) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule start time must be in the future"})
		return
	}
	// Validate that the end time is after the start time.
	if !schedule.EndTime.After(schedule.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule end time must be after start time"})
		return
	}

	created, err := repositories.InsertSchedule(&schedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to schedule session"})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// GetSchedules retrieves scheduled sessions for the authenticated user.
func GetSchedules(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	schedules, err := repositories.GetSchedulesByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve schedules"})
		return
	}
	c.JSON(http.StatusOK, schedules)
}
