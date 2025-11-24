package models

import "time"

type User struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	TelegramID   string    `json:"telegram_id" gorm:"uniqueIndex;not null"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name" gorm:"not null"`
	LastName     string    `json:"last_name"`
	RegisteredAt time.Time `json:"registered_at" gorm:"default:CURRENT_TIMESTAMP"`
	LastActive   time.Time `json:"last_active" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}
