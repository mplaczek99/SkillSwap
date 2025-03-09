package models

import "time"

// Transaction records the exchange of SkillPoints between users.
type Transaction struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Amount     int       `json:"amount"`
	Note       string    `json:"note,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}
