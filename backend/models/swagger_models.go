package models

import "time"

// ErrorResponse represents a standard API error response
// @Description Standard error response format
type ErrorResponse struct {
	Error string `json:"error" example:"Something went wrong"`
}

// User represents a registered user in the SkillSwap platform.
// @Description User account information
type SwaggerUser struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john@example.com"`
	Bio       string    `json:"bio" example:"Experienced developer with 5 years of experience"`
	Role      string    `json:"role" example:"User"` // "User" or "Admin"
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// Skill represents a skill that a user can offer or request.
// @Description Skill information
type SwaggerSkill struct {
	ID          uint      `json:"id" example:"1"`
	Name        string    `json:"name" example:"JavaScript Programming"`
	Description string    `json:"description" example:"Learn modern JavaScript from basics to advanced concepts"`
	UserID      uint      `json:"user_id" example:"1"` // ID of the user offering the skill
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// Schedule represents a scheduled skill exchange session.
// @Description Scheduled session information
type SwaggerSchedule struct {
	ID        uint      `json:"id" example:"1"`
	UserID    uint      `json:"user_id" example:"1"`
	SkillID   uint      `json:"skill_id" example:"1"`
	StartTime time.Time `json:"start_time" example:"2023-01-15T14:00:00Z"`
	EndTime   time.Time `json:"end_time" example:"2023-01-15T15:00:00Z"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// Transaction records the exchange of SkillPoints between users.
// @Description Transaction record for SkillPoints
type SwaggerTransaction struct {
	ID         uint      `json:"id" example:"1"`
	SenderID   uint      `json:"sender_id" example:"1"`
	ReceiverID uint      `json:"receiver_id" example:"2"`
	Amount     int       `json:"amount" example:"10"`
	CreatedAt  time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// Video represents a video uploaded by a user
// @Description Video information
type SwaggerVideo struct {
	ID           string    `json:"id" example:"video1.mp4"`
	Name         string    `json:"name" example:"tutorial.mp4"`
	Size         int64     `json:"size" example:"12345678"`
	UploadedAt   time.Time `json:"uploaded_at" example:"2023-01-01T00:00:00Z"`
	HasThumbnail bool      `json:"has_thumbnail" example:"true"`
}
