package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID           int        `json:"id,omitempty" gorm:"primaryKey"`
	TelegramID   string     `json:"telegram_id" gorm:"uniqueIndex;not null"`
	Username     string     `json:"username,omitempty"`
	FirstName    string     `json:"first_name" gorm:"not null"`
	LastName     string     `json:"last_name,omitempty"`
	RegisteredAt *time.Time `json:"registered_at,omitempty" gorm:"default:CURRENT_TIMESTAMP"`
	LastActive   *time.Time `json:"last_active,omitempty" gorm:"default:CURRENT_TIMESTAMP"`
}

// UnmarshalJSON handles multiple timestamp formats
func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		RegisteredAt string `json:"registered_at"`
		LastActive   string `json:"last_active"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse timestamps - handle both RFC3339 and Postgres format
	if aux.RegisteredAt != "" {
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05.999999",
			"2006-01-02T15:04:05",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, aux.RegisteredAt); err == nil {
				u.RegisteredAt = &t
				break
			}
		}
	}

	if aux.LastActive != "" {
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05.999999",
			"2006-01-02T15:04:05",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, aux.LastActive); err == nil {
				u.LastActive = &t
				break
			}
		}
	}

	return nil
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}
