package models_test

import (
	"errors"
	"testing"
	"time"

	"github.com/mplaczek99/SkillSwap/models"
)

// Add these validation methods to your models package or create them in the test
// For testing purposes, we're defining them here

// Validate validates Skill fields
func validateSkill(s *models.Skill) error {
	if s.Name == "" {
		return errors.New("name is required")
	}
	if s.Description == "" {
		return errors.New("description is required")
	}
	if s.UserID == 0 {
		return errors.New("user_id is required")
	}
	return nil
}

// Validate validates Schedule fields
func validateSchedule(s *models.Schedule) error {
	if s.UserID == 0 {
		return errors.New("user_id is required")
	}
	if s.SkillID == 0 {
		return errors.New("skill_id is required")
	}
	if s.StartTime.IsZero() {
		return errors.New("start_time is required")
	}
	if s.EndTime.IsZero() {
		return errors.New("end_time is required")
	}
	if !s.EndTime.After(s.StartTime) {
		return errors.New("end_time must be after start_time")
	}
	if s.StartTime.Before(time.Now()) {
		return errors.New("start_time must be in the future")
	}
	// Validate duration (e.g., minimum 30 minutes, maximum 4 hours)
	minDuration := 30 * time.Minute
	maxDuration := 4 * time.Hour
	duration := s.EndTime.Sub(s.StartTime)
	if duration < minDuration {
		return errors.New("schedule duration must be at least 30 minutes")
	}
	if duration > maxDuration {
		return errors.New("schedule duration cannot exceed 4 hours")
	}
	return nil
}

func TestSkillValidation(t *testing.T) {
	testCases := []struct {
		name        string
		skill       models.Skill
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid Skill",
			skill: models.Skill{
				Name:        "Test Skill",
				Description: "Test description",
				UserID:      1,
			},
			expectError: false,
		},
		{
			name: "Missing Name",
			skill: models.Skill{
				Description: "Test description",
				UserID:      1,
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "Missing Description",
			skill: models.Skill{
				Name:   "Test Skill",
				UserID: 1,
			},
			expectError: true,
			errorMsg:    "description is required",
		},
		{
			name: "Missing UserID",
			skill: models.Skill{
				Name:        "Test Skill",
				Description: "Test description",
			},
			expectError: true,
			errorMsg:    "user_id is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateSkill(&tc.skill)
			if tc.expectError && err == nil {
				t.Errorf("Expected validation error, got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no validation error, got: %v", err)
			}
			if tc.expectError && err != nil && err.Error() != tc.errorMsg {
				t.Errorf("Expected error message '%s', got '%s'", tc.errorMsg, err.Error())
			}
		})
	}
}

func TestScheduleValidation(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)

	testCases := []struct {
		name        string
		schedule    models.Schedule
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid Schedule",
			schedule: models.Schedule{
				UserID:    1,
				SkillID:   2,
				StartTime: future,
				EndTime:   future.Add(1 * time.Hour),
			},
			expectError: false,
		},
		{
			name: "Missing UserID",
			schedule: models.Schedule{
				SkillID:   2,
				StartTime: future,
				EndTime:   future.Add(1 * time.Hour),
			},
			expectError: true,
			errorMsg:    "user_id is required",
		},
		{
			name: "Missing SkillID",
			schedule: models.Schedule{
				UserID:    1,
				StartTime: future,
				EndTime:   future.Add(1 * time.Hour),
			},
			expectError: true,
			errorMsg:    "skill_id is required",
		},
		{
			name: "EndTime Before StartTime",
			schedule: models.Schedule{
				UserID:    1,
				SkillID:   2,
				StartTime: future.Add(2 * time.Hour),
				EndTime:   future.Add(1 * time.Hour),
			},
			expectError: true,
			errorMsg:    "end_time must be after start_time",
		},
		{
			name: "StartTime In Past",
			schedule: models.Schedule{
				UserID:    1,
				SkillID:   2,
				StartTime: past,
				EndTime:   future,
			},
			expectError: true,
			errorMsg:    "start_time must be in the future",
		},
		{
			name: "Duration Too Short",
			schedule: models.Schedule{
				UserID:    1,
				SkillID:   2,
				StartTime: future,
				EndTime:   future.Add(15 * time.Minute),
			},
			expectError: true,
			errorMsg:    "schedule duration must be at least 30 minutes",
		},
		{
			name: "Duration Too Long",
			schedule: models.Schedule{
				UserID:    1,
				SkillID:   2,
				StartTime: future,
				EndTime:   future.Add(5 * time.Hour),
			},
			expectError: true,
			errorMsg:    "schedule duration cannot exceed 4 hours",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateSchedule(&tc.schedule)
			if tc.expectError && err == nil {
				t.Errorf("Expected validation error, got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no validation error, got: %v", err)
			}
			if tc.expectError && err != nil && err.Error() != tc.errorMsg {
				t.Errorf("Expected error message '%s', got '%s'", tc.errorMsg, err.Error())
			}
		})
	}
}

func TestScheduleConcurrentValidation(t *testing.T) {
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)

	// Create a function that would validate if a new schedule conflicts with existing ones
	validateNoConflict := func(newSchedule models.Schedule, existingSchedules []models.Schedule) error {
		for _, existing := range existingSchedules {
			// Only check conflicts for the same user
			if existing.UserID == newSchedule.UserID {
				// Check if new schedule starts during existing schedule
				if newSchedule.StartTime.After(existing.StartTime) &&
					newSchedule.StartTime.Before(existing.EndTime) {
					return errors.New("schedule conflicts with an existing schedule")
				}

				// Check if new schedule ends during existing schedule
				if newSchedule.EndTime.After(existing.StartTime) &&
					newSchedule.EndTime.Before(existing.EndTime) {
					return errors.New("schedule conflicts with an existing schedule")
				}

				// Check if new schedule encompasses existing schedule
				if newSchedule.StartTime.Before(existing.StartTime) &&
					newSchedule.EndTime.After(existing.EndTime) {
					return errors.New("schedule conflicts with an existing schedule")
				}
			}
		}
		return nil
	}

	// Create existing schedules
	existingSchedules := []models.Schedule{
		{
			ID:        1,
			UserID:    1,
			SkillID:   2,
			StartTime: tomorrow.Add(10 * time.Hour), // 10 AM
			EndTime:   tomorrow.Add(12 * time.Hour), // 12 PM
		},
		{
			ID:        2,
			UserID:    1,
			SkillID:   3,
			StartTime: tomorrow.Add(14 * time.Hour), // 2 PM
			EndTime:   tomorrow.Add(16 * time.Hour), // 4 PM
		},
		{
			ID:        3,
			UserID:    2, // Different user
			SkillID:   2,
			StartTime: tomorrow.Add(10 * time.Hour), // 10 AM (same time as user 1's schedule)
			EndTime:   tomorrow.Add(12 * time.Hour), // 12 PM
		},
	}

	testCases := []struct {
		name        string
		schedule    models.Schedule
		expectError bool
	}{
		{
			name: "Non-conflicting Schedule",
			schedule: models.Schedule{
				UserID:    1,
				SkillID:   2,
				StartTime: tomorrow.Add(8 * time.Hour), // 8 AM
				EndTime:   tomorrow.Add(9 * time.Hour), // 9 AM
			},
			expectError: false,
		},
		{
			name: "Conflicting Start Time",
			schedule: models.Schedule{
				UserID:    1,
				SkillID:   2,
				StartTime: tomorrow.Add(11 * time.Hour), // 11 AM (during first schedule)
				EndTime:   tomorrow.Add(13 * time.Hour), // 1 PM
			},
			expectError: true,
		},
		{
			name: "Conflicting End Time",
			schedule: models.Schedule{
				UserID:    1,
				SkillID:   2,
				StartTime: tomorrow.Add(9 * time.Hour),  // 9 AM
				EndTime:   tomorrow.Add(11 * time.Hour), // 11 AM (during first schedule)
			},
			expectError: true,
		},
		{
			name: "Encompassing Existing Schedule",
			schedule: models.Schedule{
				UserID:    1,
				SkillID:   2,
				StartTime: tomorrow.Add(9 * time.Hour),  // 9 AM
				EndTime:   tomorrow.Add(13 * time.Hour), // 1 PM (encompasses first schedule)
			},
			expectError: true,
		},
		{
			name: "Same Time but Different User",
			schedule: models.Schedule{
				UserID:    3, // Different user
				SkillID:   2,
				StartTime: tomorrow.Add(10 * time.Hour), // 10 AM (same as first schedule)
				EndTime:   tomorrow.Add(12 * time.Hour), // 12 PM
			},
			expectError: false, // No conflict because it's a different user
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateNoConflict(tc.schedule, existingSchedules)
			if tc.expectError && err == nil {
				t.Errorf("Expected conflict error, got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no conflict error, got: %v", err)
			}
		})
	}
}
