package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// StringArray is a custom type to handle string arrays in PostgreSQL
type StringArray []string

// Value converts StringArray to a JSON string for storage
func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan converts a stored JSON string back to a StringArray
func (a *StringArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// Job represents a job posting in the SkillSwap platform
type Job struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	Title           string      `json:"title"`
	Company         string      `json:"company"`
	Location        string      `json:"location"`
	Description     string      `json:"description"`
	SkillsRequired  StringArray `gorm:"type:jsonb" json:"skillsRequired"`
	ExperienceLevel string      `json:"experienceLevel"` // Entry, Mid, Senior
	JobType         string      `json:"jobType"`         // Full-time, Part-time, Contract
	SalaryRange     string      `json:"salaryRange"`
	ContactEmail    string      `json:"contactEmail"`
	PostedByUserID  uint        `json:"postedByUserID"`
	PostedByName    string      `json:"postedByName"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}
