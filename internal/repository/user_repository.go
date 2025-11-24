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
	_, err := r.client.From("users").Insert(user, false, "", "*", "").ExecuteTo(&results)
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
	_, err := r.client.From("users").
		Select("*", "", false).
		Eq("telegram_id", telegramID).
		ExecuteTo(&results)
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
	_, err := r.client.From("users").
		Update(updates, "*", "").
		Eq("telegram_id", telegramID).
		ExecuteTo(&results)
	return err
}
