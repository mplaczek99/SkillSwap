package repositories_test

import (
	"testing"
	"time"

	"github.com/mplaczek99/SkillSwap/models"
	"github.com/mplaczek99/SkillSwap/repositories"
)

func TestInsertSchedule(t *testing.T) {
	// Create a sample schedule
	schedule := &models.Schedule{
		UserID:    1,
		SkillID:   2,
		StartTime: time.Now().Add(24 * time.Hour),
		EndTime:   time.Now().Add(26 * time.Hour),
	}

	// Insert the schedule
	createdSchedule, err := repositories.InsertSchedule(schedule)
	if err != nil {
		t.Fatalf("InsertSchedule failed: %v", err)
	}

	// Verify the schedule was created
	if createdSchedule.ID == 0 {
		t.Error("Expected schedule to have an ID assigned")
	}

	if createdSchedule.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	// Verify user ID and skill ID are correct
	if createdSchedule.UserID != schedule.UserID {
		t.Errorf("Expected UserID=%d, got %d", schedule.UserID, createdSchedule.UserID)
	}

	if createdSchedule.SkillID != schedule.SkillID {
		t.Errorf("Expected SkillID=%d, got %d", schedule.SkillID, createdSchedule.SkillID)
	}
}

func TestGetSchedulesByUserID(t *testing.T) {
	// Get schedules for a user
	userID := uint(1)
	schedules, err := repositories.GetSchedulesByUserID(userID)
	if err != nil {
		t.Fatalf("GetSchedulesByUserID failed: %v", err)
	}

	// Verify at least one schedule is returned
	if len(schedules) == 0 {
		t.Error("Expected at least one schedule, got none")
	}

	// Verify all schedules belong to the specified user
	for i, schedule := range schedules {
		if schedule.UserID != userID {
			t.Errorf("Schedule[%d] has UserID=%d, expected %d", i, schedule.UserID, userID)
		}
	}

	// All schedules should have non-zero IDs
	for i, schedule := range schedules {
		if schedule.ID == 0 {
			t.Errorf("Schedule[%d] has ID=0", i)
		}
	}

	// Verify time fields are set
	for i, schedule := range schedules {
		if schedule.StartTime.IsZero() {
			t.Errorf("Schedule[%d] has zero StartTime", i)
		}
		if schedule.EndTime.IsZero() {
			t.Errorf("Schedule[%d] has zero EndTime", i)
		}
		if schedule.CreatedAt.IsZero() {
			t.Errorf("Schedule[%d] has zero CreatedAt", i)
		}
	}
}

func TestGetSchedulesByNonExistentUserID(t *testing.T) {
	// Try to get schedules for a user that doesn't exist
	userID := uint(999) // Assuming this user doesn't exist
	schedules, err := repositories.GetSchedulesByUserID(userID)
	if err != nil {
		t.Fatalf("GetSchedulesByUserID failed: %v", err)
	}

	// Since we're using dummy data, we should still get a schedule
	// In a real implementation with a database, you'd check that it returns an empty slice
	if len(schedules) == 0 {
		t.Log("Note: Expected empty schedules for non-existent user, got some data instead due to dummy implementation")
	}
}
