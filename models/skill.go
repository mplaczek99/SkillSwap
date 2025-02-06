package models

import "time"

// Skill represents a skill that a user can offer or request.
type Skill struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id"` // ID of the user offering the skill
	CreatedAt   time.Time `json:"created_at"`
}
