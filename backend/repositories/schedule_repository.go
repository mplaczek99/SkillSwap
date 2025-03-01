package repositories

import (
	"time"

	"github.com/mplaczek99/SkillSwap/models"
)

// InsertSchedule saves a new schedule and returns it.
func InsertSchedule(schedule *models.Schedule) (*models.Schedule, error) {
	// In a real implementation, you would insert into your database.
	// For now, assign a dummy ID and current timestamp.
	schedule.ID = 1
	schedule.CreatedAt = time.Now()
	return schedule, nil
}

// GetSchedulesByUserID returns all schedules for a given user.
func GetSchedulesByUserID(userID uint) ([]models.Schedule, error) {
	// Dummy schedule for demonstration purposes.
	dummySchedule := models.Schedule{
		ID:        1,
		UserID:    userID,
		SkillID:   1,
		StartTime: time.Now().Add(48 * time.Hour),
		EndTime:   time.Now().Add(50 * time.Hour),
		CreatedAt: time.Now(),
	}
	return []models.Schedule{dummySchedule}, nil
}
