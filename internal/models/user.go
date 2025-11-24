package models

import "time"

type User struct {
	ID           int       `json:"id,omitempty" gorm:"primaryKey"`
	TelegramID   string    `json:"telegram_id" gorm:"uniqueIndex;not null"`
	Username     string    `json:"username,omitempty"`
	FirstName    string    `json:"first_name" gorm:"not null"`
	LastName     string    `json:"last_name,omitempty"`
	RegisteredAt time.Time `json:"registered_at" gorm:"default:CURRENT_TIMESTAMP"`
	LastActive   time.Time `json:"last_active" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}
