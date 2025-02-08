package models

import "time"

// User represents a registered user in the SkillSwap platform.
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"password"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
}
