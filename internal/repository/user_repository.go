package repository

import (
	"fmt"

	"github.com/supabase-community/supabase-go"
	"github.com/yourusername/pokemon-chatbot-api/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByTelegramID(telegramID string) (*models.User, error)
	UpdateLastActive(telegramID string) error
}

type userRepository struct {
	client *supabase.Client
}

func NewUserRepository(client *supabase.Client) UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) Create(user *models.User) error {
	var results []models.User
	err := r.client.DB.From("users").Insert(user).Execute(&results)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	if len(results) > 0 {
		*user = results[0]
	}
	return nil
}

func (r *userRepository) FindByTelegramID(telegramID string) (*models.User, error) {
	var results []models.User
	err := r.client.DB.From("users").
		Select("*").
		Eq("telegram_id", telegramID).
		Execute(&results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return &results[0], nil
}

func (r *userRepository) UpdateLastActive(telegramID string) error {
	var results []models.User
	updates := map[string]interface{}{
		"last_active": "now()",
	}
	err := r.client.DB.From("users").
		Update(updates).
		Eq("telegram_id", telegramID).
		Execute(&results)
	return err
}
