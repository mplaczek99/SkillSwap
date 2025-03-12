package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a registered user in the SkillSwap platform.
type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Email       string    `gorm:"unique" json:"email"`
	Password    string    `json:"-"` // omit from JSON responses
	Bio         string    `json:"bio"`
	Role        string    `json:"role"`                           // "User" or "Admin"
	SkillPoints int       `json:"skillPoints" gorm:"default:100"` // Default starting balance
	CreatedAt   time.Time `json:"created_at"`
}

// BeforeSave hashes the password and sets default role if empty.
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	// Only hash if the password does not already appear hashed.
	if len(u.Password) < 60 || (len(u.Password) >= 60 && u.Password[:4] != "$2a$" && u.Password[:4] != "$2b$") {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	// Set default role to "User" if none provided.
	if u.Role == "" {
		u.Role = "User"
	}
	return nil
}

// ComparePassword checks if the provided password matches the stored hash.
func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
